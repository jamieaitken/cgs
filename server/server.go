package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultAddr         = ":8080"
	defaultReadTimeout  = time.Second * 30
	defaultWriteTimeout = time.Second * 30
)

type Server struct {
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	handler      http.Handler
	tlsConfig    *TLSConfig
	httpServer   *http.Server
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

type Option func(*Server)

func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		addr:         defaultAddr,
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
		handler:      handler,
	}

	s.add(opts...)

	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      s.handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

func (s *Server) add(opts ...Option) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) ReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *Server) WriteTimeout() time.Duration {
	return s.writeTimeout
}

func (s *Server) Start(ctx context.Context) error {
	go func(ctx context.Context) {
		osSignals := make(chan os.Signal, 1)
		defer close(osSignals)

		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPROF)
		<-osSignals

		s.Stop(ctx)
	}(ctx)

	if s.tlsConfig == nil {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	}

	err := s.httpServer.ListenAndServeTLS(s.tlsConfig.CertFile, s.tlsConfig.KeyFile)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c := make(chan error)
	defer close(c)

	go func(ctx context.Context) {
		c <- s.httpServer.Shutdown(ctx)
	}(ctx)

	return <-c
}
