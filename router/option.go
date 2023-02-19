package router

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/ext/jwt"
	"github.com/skroczek/acme-restful/ext/oidc"
)

type Option func(r *gin.Engine)

func WithBasicAuth(auth gin.Accounts) Option {
	return func(r *gin.Engine) {
		r.Use(gin.BasicAuth(auth))
	}
}

func WithOIDC(o oidc.Oidc) Option {
	return func(r *gin.Engine) {
		r.Use(o.Middleware())
	}
}

func WithJWTAuth() Option {
	return func(r *gin.Engine) {
		r.Use(jwt.Protect)
	}
}

// WithTrustedProxies set trusted proxies
func WithTrustedProxies(proxies []string) Option {
	return func(r *gin.Engine) {
		err := r.SetTrustedProxies(proxies)
		if err != nil {
			panic(fmt.Errorf("error setting trusted proxies: %w", err))
		}
	}
}

// WithoutTrustedProxies set trusted proxies to nil
// @see https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
func WithoutTrustedProxies() Option {
	return WithTrustedProxies(nil)
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
