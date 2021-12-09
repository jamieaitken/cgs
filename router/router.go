package router

import (
	"github.com/jamieaitken/promred/handler"
	"github.com/jamieaitken/requestid"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	mux             *http.ServeMux
	instrumentation handler.Handler
	tracer          *requestid.Tracer
}

type Route struct {
	Path         string
	HandlerFuncs map[string]http.HandlerFunc
}

type Option func(*Router)

func New(health healthcheck.Handler, opts ...Option) *Router {
	mux := http.NewServeMux()

	instrHandler := handler.New()

	tracer := requestid.New()

	mux.Handle("/metrics", handlers.MethodHandler{
		http.MethodGet: promhttp.Handler(),
	})
	mux.Handle("/live", handlers.MethodHandler{
		http.MethodGet: http.HandlerFunc(instrHandler.HandleFor(health.LiveEndpoint)),
	})
	mux.Handle("/ready", handlers.MethodHandler{
		http.MethodGet: http.HandlerFunc(instrHandler.HandleFor(health.ReadyEndpoint)),
	})

	r := &Router{
		mux:             mux,
		instrumentation: instrHandler,
		tracer:          tracer,
	}

	r.Add(opts...)

	return r
}

func (r *Router) Mux() *http.ServeMux {
	return r.mux
}

func (r *Router) Add(options ...Option) {
	for _, opt := range options {
		opt(r)
	}
}
