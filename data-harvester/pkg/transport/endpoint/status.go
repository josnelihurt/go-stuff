package endpoint

import (
	"context"

	gokitEndpoint "github.com/go-kit/kit/endpoint"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/http"
)

// Status endpoint mapping
func (e Endpoints) Status(ctx context.Context) (bool, error) {
	req := http.StatusRequest{}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	statusResp := resp.(http.StatusResponse)
	return statusResp.Status, nil
}

// makeStatusEndpoint creates a caller that will redirect the call into the service status
func makeStatusEndpoint(srv service.DataHarvestService) gokitEndpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(http.StatusRequest) // we really just need the request, we don't use any value from it
		s, err := srv.Status(ctx)
		if err != nil {
			return http.StatusResponse{s}, err
		}
		return http.StatusResponse{s}, nil
	}
}
