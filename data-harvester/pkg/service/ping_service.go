package service

import (
	"context"
	"log"
)

type pingService struct{}

//NewPingService returns an implementation of a DataHarvestService for ping recollector
func newPingService() DataHarvestService {
	log.Println("Creating new NewPingService")
	return &pingService{}
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

//Collect implemented from Service
func (context *pingService) Collect(ctx context.Context, param DataHarvestServiceParam) (DataHarvestServiceResult, error) {
	return DataHarvestServiceResult{}, nil
}

func (context *pingService) Status(ctx context.Context) (bool, error) {
	return false, nil
}
