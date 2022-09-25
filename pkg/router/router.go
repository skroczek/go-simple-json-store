package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/pkg/server"
)

type Option func(r *gin.Engine)

func DefaultRouter(server *server.Server, options ...Option) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	for _, option := range options {
		option(router)
	}

	router.GET("/*path", server.GetHandler)
	router.POST("/*path", server.PostHandler)
	router.PUT("/*path", server.PutHandler)
	router.DELETE("/*path", server.DeleteHandler)
	router.PATCH("/*path", server.PatchHandler)
	router.HEAD("/*path", server.HeadHandler)
	router.OPTIONS("/*path", server.OptionsHandler)

	return router
}

func WithBasicAuth(auth gin.Accounts) Option {
	return func(r *gin.Engine) {
		r.Use(gin.BasicAuth(auth))
	}
}

func WithDefaultCors(allowCredentials bool) Option {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = allowCredentials
	return func(r *gin.Engine) {
		r.Use(cors.New(config))
	}
}
