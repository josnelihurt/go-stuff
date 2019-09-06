package service

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grandcat/zeroconf"
)

type pingService struct {
	dataCache   *DataHarvestServiceResult
	atomicCache atomic.Value
	writersSync sync.Mutex // used only by writers
	started     bool
}

func getLocalAddresses() (result []string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				//	ip = v.IP
				continue
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			result = append(result, addr.(*net.IPNet).String())
		}
	}
	return
}
func (service *pingService) backgroundScanner() {
	currentResult := &DataHarvestServiceResult{}
	currentResult.HostIPs = getLocalAddresses()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			currentResult.DiscoveredElements = append(currentResult.DiscoveredElements, Entry{
				HostName:    entry.HostName,
				IP:          entry.AddrIPv4[0],
				ServicePort: entry.Port,
			})
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	service.safeWriteCache(currentResult)
	service.started = true
	time.AfterFunc(30*time.Second, func() { service.backgroundScanner() })

}
func (service *pingService) safeReadCache() (result DataHarvestServiceResult) {
	// read function can be used to read the data without further synchronization
	result = *service.atomicCache.Load().(*DataHarvestServiceResult) // I wonder if a copy is created by returning the instance
	log.Printf(" %v", result)
	return result
}
func (service *pingService) safeWriteCache(replaceWith *DataHarvestServiceResult) {
	service.writersSync.Lock() // synchronize with other potential writers
	defer service.writersSync.Unlock()
	_ = service.atomicCache.Load().(*DataHarvestServiceResult) // load current value of the data structure
	service.atomicCache.Store(replaceWith)                     // atomically replace the current object with the new one
	// At this point all new readers start working with the new version.
	// The old version will be garbage collected once the existing readers
	// (if any) are done with it.
}

//NewPingService returns an implementation of a DataHarvestService for ping recollector
func newPingService() DataHarvestService {
	log.Println("Creating new NewPingService")
	service := &pingService{}
	data := &DataHarvestServiceResult{}
	data.HostIPs = getLocalAddresses()
	service.atomicCache.Store(data)
	go service.backgroundScanner()
	return service
}

// New returns a DataHarvestService with all of the expected middleware wired in.
func New(middleware []Middleware) DataHarvestService {
	log.Println("Creating new service")
	var svc = newPingService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

//Collect method implemented from Service
func (service *pingService) Collect(ctx context.Context, param DataHarvestServiceParam) (result DataHarvestServiceResult, err error) {
	result = service.safeReadCache()
	log.Printf(" %v", result)
	return result, nil
}

//Status method implemented from Service
func (service *pingService) Status(ctx context.Context) (bool, error) {
	return service.started, nil
}
