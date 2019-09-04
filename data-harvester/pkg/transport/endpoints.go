package transport

import (
	"github.com/go-kit/kit/endpoint"
	service "github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
)

// Endpoints are exposed this are the connection methods for remote execution
type Endpoints struct {
	CollectEndpoint endpoint.Endpoint
	StatusEndpoint  endpoint.Endpoint
}

//MakeEndpoints exposes creation for al endpoints in the service
func MakeEndpoints(srv service.DataHarvestService) Endpoints {
	return Endpoints{
		StatusEndpoint: makeStatusEndpoint(srv),
	}
}
