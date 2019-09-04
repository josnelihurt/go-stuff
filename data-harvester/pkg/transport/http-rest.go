package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

type StatusRequest struct{}

type StatusResponse struct {
	Status bool `json:"status"`
}

func decodeStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StatusRequest
	return req, nil
}

// Last but not least, we have the encoder for the response output
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
