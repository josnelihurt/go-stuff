package grpc

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	serviceEndpoint "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/endpoint"
	serviceProto "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/grpc/proto"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/http"
)

// makeStatusHandler creates the handler logic
func makeStatusHandler(endpoints serviceEndpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.StatusEndpoint, decodeStatusRequest, encodeStatusResponse, options...)
}
func decodeStatusRequest(_ context.Context, grpcRequest interface{}) (businessRequest interface{}, err error) {
	return http.StatusRequest{}, nil
}
func encodeStatusResponse(_ context.Context, businessResponse interface{}) (grpcResponse interface{}, err error) {
	return &serviceProto.DataHarvestServiceResponse{Status: businessResponse.(http.StatusResponse).Status}, nil
}
func (g *grpcServer) Status(ctx context.Context, req *serviceProto.DataHarvestServiceRequest) (*serviceProto.DataHarvestServiceResponse, error) {
	_, rep, err := g.hnd.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*serviceProto.DataHarvestServiceResponse), nil
}
