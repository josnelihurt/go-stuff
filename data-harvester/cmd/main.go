package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/josnelihurt/go-stuff/data-harvester/cmd/server"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/endpoint"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := service.NewPingService()
	// mapping endpoints
	endpoints := endpoint.Endpoints{
		StatusEndpoint: endpoint.MakeStatusEndpoint(srv),
	}
	// Interrupt handler.
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	// HTTP transport
	go func() {
		log.Println("service is listening on port:", *httpAddr)
		handler := server.NewHTTP(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)

}
