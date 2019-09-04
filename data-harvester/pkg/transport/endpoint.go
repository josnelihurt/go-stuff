package transport

import (
	"context"
	"encoding/json"
	"net/http"

	service "github.com/josnelihurt/go-stuff/data-harvester/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

// makeStatusEndpoint create a caller that will redirect the call into the service status
func makeStatusEndpoint(srv service.DataHarvestService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(StatusRequest) // we really just need the request, we don't use any value from it
		s, err := srv.Status(ctx)
		if err != nil {
			return StatusResponse{s}, err
		}
		return StatusResponse{s}, nil
	}
}

// Status endpoint mapping
func (e Endpoints) Status(ctx context.Context) (bool, error) {
	req := StatusRequest{}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	statusResp := resp.(StatusResponse)
	return statusResp.Status, nil
}

//DecodeStatusRequest parse data in
func DecodeStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StatusRequest
	return req, nil
}

//EncodeResponse pase data out
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
