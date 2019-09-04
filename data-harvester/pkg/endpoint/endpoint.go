package endpoint

import (
	"context"
	"encoding/json"
	"net/http"

	service "github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/transport"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints are exposed
type Endpoints struct {
	CollectEndpoint endpoint.Endpoint
	StatusEndpoint  endpoint.Endpoint
}

// MakeStatusEndpoint returns the response from our service "status"
func MakeStatusEndpoint(srv service.DataHarvestService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(transport.StatusRequest) // we really just need the request, we don't use any value from it
		s, err := srv.Status(ctx)
		if err != nil {
			return transport.StatusResponse{s}, err
		}
		return transport.StatusResponse{s}, nil
	}
}

// Status endpoint mapping
func (e Endpoints) Status(ctx context.Context) (bool, error) {
	req := transport.StatusRequest{}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	statusResp := resp.(transport.StatusResponse)
	return statusResp.Status, nil
}
func DecodeStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.StatusRequest
	return req, nil
}

// Last but not least, we have the encoder for the response output
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
