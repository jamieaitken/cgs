package server

import (
	"time"
)

func WithAddr(addr string) Option {
	return func(server *Server) {
		server.addr = addr
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(server *Server) {
		server.writeTimeout = timeout
	}
}

func WithTLSConfig(conf *TLSConfig) Option {
	return func(server *Server) {
		server.tlsConfig = conf
	}
}
