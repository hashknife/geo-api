package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/briandowns/hashknife/common/middleware"
	"github.com/briandowns/hashknife/common/services"
	"github.com/briandowns/hashknife/geo-api/config"
	"github.com/briandowns/hashknife/geo-api/endpoints"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	component = "geo-api"
	prefix    = "/" + component
	version   = "/v1"
)

// HTTPListenerParams holds parameters for building routers
type HTTPListenerParams struct {
	Logger  kitlog.Logger
	Root    context.Context
	ErrChan chan error
	Config  *config.Config
}

// StartApplicationHTTPListener creates a goroutine that has an HTTP listener for the application endpoints
func StartApplicationHTTPListener(hlp *HTTPListenerParams, t38 services.Tile38er) {
	go func() {
		ctx, cancel := context.WithCancel(hlp.Root)
		defer cancel()
		hlp.Logger.Log("HTTPAddress", *hlp.Config.HTTPAddress, "transport", "HTTP/JSON")
		cs := endpoints.NewCourierService(hlp.Config, t38)
		router := createApplicationRouter(ctx, hlp.Logger, hlp.Config, cs)
		//listenerMetrics := metrics.NewSimpleMetricsMiddleware(component, "app_http_listener", *hlp.Config.StatsdReportingIntervalSeconds, *hlp.Config.StatsdAddress)
		//hlp.ErrChan <- http.ListenAndServe(*hlp.Config.HTTPAddress, handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(kitlog.NewStdlibAdapter(hlp.Logger), listenerMetrics.Annotate(router))))
		hlp.ErrChan <- http.ListenAndServe(*hlp.Config.HTTPAddress, handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(kitlog.NewStdlibAdapter(hlp.Logger), router)))
	}()
}

// createApplicationRouter sets up the router that will handle all of the application routes
func createApplicationRouter(ctx context.Context, l kitlog.Logger, conf *config.Config, cs endpoints.CourierServicer) *mux.Router {
	router := mux.NewRouter().PathPrefix(prefix + version).Subrouter()
	auth := middleware.NewHashknifeRequestAuthenticator(conf.HashknifeAuthToken)
	router.Handle(
		"/courier/{account_id}/{courier_id}",
		kithttp.NewServer(
			endpoint.Chain(auth.EndpointAuthenticate())(cs.Location),
			decodeCourierHTTPRequest,
			encodeResponse,
			kithttp.ServerBefore(middleware.KitServerBefore),
		)).Methods(http.MethodGet)
	router.Handle(
		"/courier/{account_id}/{courier_id}",
		kithttp.NewServer(
			endpoint.Chain(auth.EndpointAuthenticate())(cs.UpdateLocation),
			decodeCourierUpdateRequest,
			encodeResponse,
			kithttp.ServerBefore(middleware.KitServerBefore),
		)).Methods(http.MethodPost)
	return router
}

// StartHealthCheckHTTPListener creates a goroutine that has an HTTP listener for the healthcheck endpoint
func StartHealthCheckHTTPListener(p *HTTPListenerParams, gs string) {
	go func() {
		ctx, cancel := context.WithCancel(p.Root)
		defer cancel()
		p.Logger.Log("HealthCheckAddress", *p.Config.HealthCheckAddress, "transport", "HTTP/JSON")
		router := createHealthCheckRouter(ctx, p.Logger, endpoints.NewHealthCheckEndpoint(gs))
		p.ErrChan <- http.ListenAndServe(*p.Config.HealthCheckAddress, handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(kitlog.NewStdlibAdapter(p.Logger), router)))
	}()
}

// createHealthCheckRouter setups up the router that provides the health checking functionality
func createHealthCheckRouter(ctx context.Context, l kitlog.Logger, h endpoints.HealthCheckServicer) *mux.Router {
	router := mux.NewRouter().PathPrefix(prefix).Subrouter()
	router.Handle(
		"/healthcheck",
		kithttp.NewServer(
			h.Run,
			noOpDecodeRequest,
			encodeHealthCheckHTTPResponse,
		)).Methods(http.MethodGet)
	return router
}

// encodeHealthCheckHTTPResponse encodes the response given by the healthCheck Endpoint
func encodeHealthCheckHTTPResponse(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(i.(*endpoints.HealthCheckResponse))
}

// encodeResponse is used by the application routes to return their data
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	addPoweredHeaders(w)
	return json.NewEncoder(w).Encode(response)
}

// decodeCourierHTTPRequest
func decodeCourierHTTPRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var courierRequest endpoints.CourierServiceLocationRequest
	params := mux.Vars(r)
	if params["account_id"] != "" {
		courierRequest.AccountID = params["account_id"]
	}
	courierRequest.CourierID = params["courier_id"]
	if params["courier_id"] != "" {
		courierRequest.CourierID = params["courier_id"]
	}
	return &courierRequest, nil
}

// decodeCourierUpdateRequest
func decodeCourierUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Body == nil {
		return nil, errors.New("error: empty body")
	}
	defer r.Body.Close()
	var ur endpoints.CourierServiceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		return nil, err
	}
	params := mux.Vars(r)
	if params["account_id"] != "" {
		ur.AccountID = params["account_id"]
	}
	if params["courier_id"] != "" {
		ur.CourierID = params["courier_id"]
	}
	if params["latitude"] != "" {
		ur.CourierID = params["latitude"]
	}
	if params["longitude"] != "" {
		ur.CourierID = params["longitude"]
	}
	return &ur, nil
}

// noOpDecodeRequest is of type GoKit http.DecodeRequestFunc and is used in place of writing out
// an empty function (func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },)
// directly in the NewServer call.
func noOpDecodeRequest(context.Context, *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

// addPoweredHeaders
func addPoweredHeaders(w http.ResponseWriter) {
	w.Header().Add("X-Powered-By", "hashknife"+component)
}
