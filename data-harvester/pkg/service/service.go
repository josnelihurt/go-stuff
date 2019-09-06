package service

import (
	"context"
	"net"
)

//DataHarvestServiceParam represents an operation parameters
type DataHarvestServiceParam struct {
}

//Entry represents the data collected
type Entry struct {
	IP          net.IP `json:"ip"`
	HostName    string `json:"hostname"`
	ServicePort int    `json:"port"`
}

//DataHarvestServiceResult represents an operation response
type DataHarvestServiceResult struct {
	HostIPs            []string `json:"hostIps"`
	DiscoveredElements []Entry  `json:"items"`
}

//DataHarvestService describes how data harvester service should behave
type DataHarvestService interface {
	Collect(ctx context.Context, param DataHarvestServiceParam) (DataHarvestServiceResult, error)
	Status(ctx context.Context) (bool, error)
}
