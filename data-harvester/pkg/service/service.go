package service

import "context"

//DataHarvestServiceParam represents an operation parameters
type DataHarvestServiceParam struct {
}

//DataHarvestServiceResult represents an operation response
type DataHarvestServiceResult struct {
}

//DataHarvestService describes how data harvester service should behave
type DataHarvestService interface {
	Collect(ctx context.Context, param DataHarvestServiceParam) (DataHarvestServiceResult, error)
	Status(ctx context.Context) (bool, error)
}
