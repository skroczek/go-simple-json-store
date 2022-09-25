package pkg

import (
	"github.com/gin-gonic/gin"
)

func DefaultRouter(server *Server) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/*path", server.GetHandler)
	router.POST("/*path", server.PostHandler)
	router.PUT("/*path", server.PutHandler)
	router.DELETE("/*path", server.DeleteHandler)
	router.PATCH("/*path", server.PatchHandler)
	router.HEAD("/*path", server.HeadHandler)
	router.OPTIONS("/*path", server.OptionsHandler)
	return router
}
