package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const listAllSuffix = "__list.json"

func getListHandler(c *gin.Context, be backend.Backend) {
	urlPath := c.Request.URL.Path
	data, err := be.List(c, urlPath[0:len(urlPath)-len(listAllSuffix)])
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	if _, ok := c.GetQuery("withoutExtension"); ok {
		for i, v := range data {
			data[i] = strings.TrimSuffix(v, filepath.Ext(v))
		}
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
