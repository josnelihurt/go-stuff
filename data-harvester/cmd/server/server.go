package server

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/gorilla/mux"
	service "github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
	serviceEnpoint "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/endpoint"
	serviceGrpc "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/grpc"
	serviceProto "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/grpc/proto"
	serviceHttp "github.com/josnelihurt/go-stuff/data-harvester/pkg/transport/http"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	grpc "google.golang.org/grpc"

	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	gokithttp "github.com/go-kit/kit/transport/http"
)

const (
	httpProto = "http"
	httpPort  = ":8080"
	httpMsg   = "http listen address"
)

//Server represents an abstraction for a server
type Server struct {
	httpAddress *string
	errorChan   chan error
	ctx         context.Context
	endpoints   serviceEnpoint.Endpoints
	service     service.DataHarvestService
}

var tracer opentracinggo.Tracer
var logger kitLog.Logger

func getServiceMiddleware(logger kitLog.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = append(mw, service.LoggingMiddleware(logger))
	return
}

//NewServer initialize all components required to start the server, it will include differents transport layers in the future
func NewServer() *Server {
	server := &Server{}
	server.errorChan = make(chan error)
	server.httpAddress = flag.String("http", ":8080", "http listen address")
	flag.Parse()
	server.ctx = context.Background()
	// Create a single logger, which we'll use and give to other components.
	logger = kitLog.NewLogfmtLogger(os.Stderr)
	logger = kitLog.With(logger, "ts", kitLog.DefaultTimestampUTC)
	logger = kitLog.With(logger, "caller", kitLog.DefaultCaller)
	tracer = opentracinggo.GlobalTracer()

	server.service = service.New(getServiceMiddleware(logger))
	server.endpoints = serviceEnpoint.MakeEndpoints(server.service)
	return server
}

//GetErrorChannel returns the main error channel
func (server *Server) GetErrorChannel() chan error {
	return server.errorChan
}

//Run the internal components
func (server *Server) Run() {
	server.newGRPC()
	server.listenAndServeHTTP()
}

func (server *Server) newGRPC() {
	options := map[string][]gokitgrpc.ServerOption{"Status": {gokitgrpc.ServerErrorLogger(logger), gokitgrpc.ServerBefore(opentracing.GRPCToContext(tracer, "Status", logger))}}
	g := &group.Group{}
	grpcListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", ":8081")
		baseServer := grpc.NewServer()
		srv := serviceGrpc.NewGRPCServer(server.endpoints, options)
		serviceProto.RegisterDataHarvestServiceServer(baseServer, srv)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})
	g.Run()

	//initMetricsEndpoint(eps)
	//initCancelInterrupt(g)
}

// NewHTTP is a good little server
func (server *Server) newHTTP() http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware) // @see https://stackoverflow.com/a/51456342

	r.Methods("GET").Path("/status").Handler(gokithttp.NewServer(
		server.endpoints.StatusEndpoint,
		serviceHttp.DecodeStatusRequest,
		serviceHttp.EncodeResponse,
	))
	r.Methods("GET").Path("/discover").Handler(gokithttp.NewServer(
		server.endpoints.CollectEndpoint,
		serviceHttp.DecodeCollectRequest,
		serviceHttp.EncodeCollectResponse,
	))

	return r
}
func (server *Server) listenAndServeHTTP() {
	log.Println("service is listening on port:", *server.httpAddress)
	handler := server.newHTTP()
	server.errorChan <- http.ListenAndServe(*server.httpAddress, handler)
}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
