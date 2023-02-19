package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/pkg/ext/jwt"
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

func WithOICD(o oidc.Oicd) Option {
	return func(r *gin.Engine) {
		r.Use(o.Middleware())
	}
}

func WithJWTAuth() Option {
	return func(r *gin.Engine) {
		r.Use(jwt.Protect)
	}
}

func WithDefaultCors(allowCredentials bool) Option {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = allowCredentials
	if allowCredentials {
		config.AddAllowHeaders("Authorization")
	}
	return func(r *gin.Engine) {
		r.Use(cors.New(config))
	}
}
