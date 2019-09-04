package server

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	kitLog "github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/service"
	"github.com/josnelihurt/go-stuff/data-harvester/pkg/transport"
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
	endpoints   transport.Endpoints
}

var logger kitLog.Logger

func getServiceMiddleware(logger kitLog.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = append(mw, service.LoggingMiddleware(logger))
	// Append your middleware here

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
	srv := service.New(getServiceMiddleware(logger))
	server.endpoints = transport.MakeEndpoints(srv)
	return server
}

//GetErrorChannel returns the main error channel
func (context *Server) GetErrorChannel() chan error {
	return context.errorChan
}

//Run the internal components
func (context *Server) Run() {
	context.listentAndServeHTTP()
}

// NewHTTP is a good little server
func (context *Server) NewHTTP() http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware) // @see https://stackoverflow.com/a/51456342

	r.Methods("GET").Path("/status").Handler(httptransport.NewServer(
		context.endpoints.StatusEndpoint,
		transport.DecodeStatusRequest,
		transport.EncodeResponse,
	))

	return r
}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func (context *Server) listentAndServeHTTP() {
	log.Println("service is listening on port:", *context.httpAddress)
	handler := context.NewHTTP()
	context.errorChan <- http.ListenAndServe(*context.httpAddress, handler)
}
