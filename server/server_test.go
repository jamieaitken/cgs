package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/jamieaitken/cgs"

	"github.com/google/go-cmp/cmp"
	"github.com/jamieaitken/cgs/server"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                 string
		givenOpts            []server.Option
		expectedAddr         string
		expectedReadTimeout  time.Duration
		expectedWriteTimeout time.Duration
	}{
		{
			name:                 "given custom read timeout, expect defaults for the rest",
			givenOpts:            []server.Option{server.WithReadTimeout(time.Second * 60)},
			expectedAddr:         ":8080",
			expectedReadTimeout:  time.Second * 60,
			expectedWriteTimeout: time.Second * 30,
		},
		{
			name:                 "given custom write timeout, expect defaults for the rest",
			givenOpts:            []server.Option{server.WithWriteTimeout(time.Second * 60)},
			expectedAddr:         ":8080",
			expectedReadTimeout:  time.Second * 30,
			expectedWriteTimeout: time.Second * 60,
		},
		{
			name:                 "given custom addr, expect defaults for the rest",
			givenOpts:            []server.Option{server.WithAddr(":3030")},
			expectedAddr:         ":3030",
			expectedReadTimeout:  time.Second * 30,
			expectedWriteTimeout: time.Second * 30,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := server.New(nil, test.givenOpts...)

			if !cmp.Equal(actual.Addr(), test.expectedAddr) {
				t.Fatalf(cmp.Diff(actual.Addr(), test.expectedAddr))
			}

			if !cmp.Equal(actual.ReadTimeout(), test.expectedReadTimeout) {
				t.Fatalf(cmp.Diff(actual.ReadTimeout(), test.expectedReadTimeout))
			}

			if !cmp.Equal(actual.WriteTimeout(), test.expectedWriteTimeout) {
				t.Fatalf(cmp.Diff(actual.WriteTimeout(), test.expectedWriteTimeout))
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	tests := []struct {
		name      string
		givenOpts []server.Option
	}{
		{
			name:      "given absent tls config, expect non-tls server started",
			givenOpts: []server.Option{},
		},
		{
			name: "given tls config, expect tls server started",
			givenOpts: []server.Option{
				server.WithTLSConfig(&server.TLSConfig{
					CertFile: "./server_test.client.chain.crt",
					KeyFile:  "./server_test_test.client.key",
				}),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New(
				cgs.WithRouter(),
				cgs.WithServer(test.givenOpts...),
			)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			chn := make(chan error, 1)

			go func() {
				chn <- app.Server().Start(context.Background())
			}()

			go func() {
				time.Sleep(1 * time.Second)

				err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				if err != nil {
					t.Error(err)
				}
			}()

			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

			rr := httptest.NewRecorder()

			app.Router().Mux().ServeHTTP(rr, req)

			resp := rr.Result()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %v", resp.StatusCode)
			}
		})
	}
}
