package endpoint

import (
	"context"

	gokitEndpoint "github.com/go-kit/kit/endpoint"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/http"
)

// Collect endpoint mapping
func (e Endpoints) Collect(ctx context.Context) (service.DataHarvestServiceResult, error) {
	req := http.CollectRequest{}
	resp, err := e.CollectEndpoint(ctx, req)
	if err != nil {
		return service.DataHarvestServiceResult{}, err
	}
	response := resp.(http.CollectResponse)
	return response.Internal, nil
}

//makeCollectEndpoint creates a caller that will redirect the call into the service
func makeCollectEndpoint(srv service.DataHarvestService) gokitEndpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := http.CollectRequest{} // we really just need the request, we don't use any value from it
		s, err := srv.Collect(ctx, req.Param)
		if err != nil {
			return http.CollectResponse{s}, err
		}
		return http.CollectResponse{s}, nil
	}
}
