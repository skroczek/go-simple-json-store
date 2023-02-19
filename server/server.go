package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"github.com/skroczek/acme-restful/router"
)

type Server struct {
	Backend       backend.Backend
	routerOptions []router.Option
}

func (s *Server) AddRouterOption(option ...router.Option) {
	s.routerOptions = append(s.routerOptions, option...)
}

func NewServer(opts ...Options) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) PrepareEngine() *gin.Engine {
	r := router.DefaultRouter(s.routerOptions...)
	r.GET("/*path", s.GetHandler)
	r.POST("/*path", s.PostHandler)
	r.PUT("/*path", s.PutHandler)
	r.DELETE("/*path", s.DeleteHandler)
	r.PATCH("/*path", s.PatchHandler)
	r.HEAD("/*path", s.HeadHandler)
	r.OPTIONS("/*path", s.OptionsHandler)
	return r
}

func (s *Server) Run(addr ...string) {
	_ = s.PrepareEngine().Run(addr...)
}

func (s *Server) RunUnix(path string) {
	_ = s.PrepareEngine().RunUnix(path)
}
