package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/pkg/backend"
	"log"
	"net/http"
	"os"
	"strings"
)

const listAllSuffix = "__list.json"

func getListHandler(c *gin.Context, be backend.Backend) {
	urlPath := c.Request.URL.Path
	data, err := be.List(urlPath[0 : len(urlPath)-len(listAllSuffix)])
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.AbortWithStatusJSON(http.StatusOK, data)
}

func WithListAll() Options {
	return func(s *Server) {
		s.AddRouterOption(func(r *gin.Engine) {
			r.Use(func(c *gin.Context) {
				if c.Request.Method == http.MethodGet && strings.HasSuffix(c.Request.URL.Path, listAllSuffix) {
					getListHandler(c, s.Backend)
					return
				}
				c.Next()
			})
		})
	}
}
