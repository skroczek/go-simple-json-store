package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Option func(r *gin.Engine)

func DefaultRouter(options ...Option) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	for _, option := range options {
		option(router)
	}

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
