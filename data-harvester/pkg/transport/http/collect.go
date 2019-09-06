package http

import (
	"context"
	"encoding/json"
	"net/http"

	service "github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
)

//CollectRequest request package
type CollectRequest struct {
	Param service.DataHarvestServiceParam
}

//CollectResponse response package
type CollectResponse struct {
	Internal service.DataHarvestServiceResult `json:"response"`
}

//DecodeCollectRequest parse data in
func DecodeCollectRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StatusRequest
	return req, nil
}

//EncodeCollectResponse pase data out
func EncodeCollectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
