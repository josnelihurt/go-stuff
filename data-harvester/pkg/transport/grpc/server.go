package grpc

import (
	"github.com/go-kit/kit/transport/grpc"
	serviceEndpoint "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/endpoint"
	serviceProto "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/grpc/proto"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	hnd grpc.Handler
}

//NewGRPCServer creates a new gRPC Server
func NewGRPCServer(endpoints serviceEndpoint.Endpoints, options map[string][]grpc.ServerOption) serviceProto.DataHarvestServiceServer {
	return &grpcServer{hnd: makeStatusHandler(endpoints, options["Status"])}
}
