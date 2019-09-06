package http

import (
	"context"
	"encoding/json"
	"net/http"
)

//StatusRequest request package
type StatusRequest struct{
}

//StatusResponse response package
type StatusResponse struct {
	Status bool `json:"status"`
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
