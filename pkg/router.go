package pkg

import (
	"github.com/gin-gonic/gin"
)

func DefaultRouter(server *Server, handlers ...gin.HandlerFunc) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	group := router.Group("/", handlers...)

	group.GET("/*path", server.GetHandler)
	group.POST("/*path", server.PostHandler)
	group.PUT("/*path", server.PutHandler)
	group.DELETE("/*path", server.DeleteHandler)
	group.PATCH("/*path", server.PatchHandler)
	group.HEAD("/*path", server.HeadHandler)
	group.OPTIONS("/*path", server.OptionsHandler)

	return router
}
