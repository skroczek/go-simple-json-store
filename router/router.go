package router

import (
	"github.com/gin-gonic/gin"
)

func DefaultRouter(options ...Option) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	for _, option := range options {
		option(router)
	}

	return router
}
