package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"github.com/skroczek/acme-restful/router"
	"net/http"
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

func (s *Server) prepareEngine() *gin.Engine {
	r := gin.New()

	// ToDo: make logger and recovery middleware configurable
	r.Use(gin.Logger(), gin.Recovery())

	for _, option := range s.routerOptions {
		option(r)
	}

	r.Use(func(context *gin.Context) {
		switch context.Request.Method {
		case "GET":
			s.GetHandler(context)
		case "POST":
			s.PostHandler(context)
		case "PUT":
			s.PutHandler(context)
		case "DELETE":
			s.DeleteHandler(context)
		case "PATCH":
			s.PatchHandler(context)
		case "HEAD":
			s.HeadHandler(context)
		case "OPTIONS":
			s.OptionsHandler(context)
		default:
			context.AbortWithStatus(http.StatusMethodNotAllowed)
		}
	})

	return r
}

func (s *Server) Run(addr ...string) {
	_ = s.prepareEngine().Run(addr...)
}

func (s *Server) RunUnix(path string) {
	_ = s.prepareEngine().RunUnix(path)
}
