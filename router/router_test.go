package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jamieaitken/cgs"
	"github.com/jamieaitken/cgs/router"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		givenOpts []router.Option
	}{
		{
			givenOpts: []router.Option{router.WithRoute(router.Route{
				Path: "/v1/music",
				HandlerFuncs: map[string]http.HandlerFunc{
					http.MethodGet: func(writer http.ResponseWriter, request *http.Request) {
						writer.WriteHeader(http.StatusOK)
					},
				},
			})},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New(
				cgs.WithRouter(test.givenOpts...),
			)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/v1/music", nil)
			rr := httptest.NewRecorder()

			app.Router().Mux().ServeHTTP(rr, req)

			resp := rr.Result()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %v", resp.StatusCode)
			}
		})
	}
}
