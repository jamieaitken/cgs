package router

import (
	"github.com/gorilla/handlers"
	"github.com/jamieaitken/requestid"
	"net/http"
)

func WithRoute(route Route) Option {
	return func(router *Router) {
		router.Mux().Handle(route.Path, buildHandler(router, route))
	}
}

func WithTracer(opts ...requestid.Option) Option {
	return func(router *Router) {
		router.tracer = requestid.New(opts...)
	}
}

func buildHandler(router *Router, route Route) http.Handler {
	h := handlers.MethodHandler{}

	for method, handler := range route.HandlerFuncs {
		h[method] = router.tracer.Trace(router.instrumentation.HandleFor(handler))
	}

	return h
}
