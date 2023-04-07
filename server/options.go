package server

import (
	"github.com/skroczek/go-simple-json-store/backend"
	"github.com/skroczek/go-simple-json-store/router"
)

type Options func(*Server)

func WithBackend(be backend.Backend) Options {
	return func(s *Server) {
		if pbe, ok := be.(backend.Proxy); ok {
			pbe.SetBackend(s.Backend)
		}
		s.Backend = be
	}
}

func WithRouterOptions(opts ...router.Option) Options {
	return func(s *Server) {
		s.AddRouterOption(opts...)
	}
}
