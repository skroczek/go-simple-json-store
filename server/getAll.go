package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"github.com/skroczek/acme-restful/helper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const getAllSuffix = "/__all.json"

func getAllHandler(c *gin.Context, be backend.Backend) {
	urlPath := c.Request.URL.Path
	path := urlPath[0 : len(urlPath)-len(getAllSuffix)]
	list, err := be.List(c, path)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	data := make([]interface{}, len(list))
	ch := make(chan interface{}, len(list))
	for i := range list {
		go func(k int) {
			obj, err := helper.FromJSON(be.Get(c, filepath.Join(path, list[k])))
			if err != nil {
				log.Panicf("Error: %+v", err)
			}
			ch <- obj
		}(i)
	}
	for i := range list {
		data[i] = <-ch
	}
	c.AbortWithStatusJSON(http.StatusOK, data)
}

func WithGetAll() Options {
	return func(s *Server) {
		s.AddRouterOption(func(r *gin.Engine) {
			r.Use(func(c *gin.Context) {
				if c.Request.Method == http.MethodGet && strings.HasSuffix(c.Request.URL.Path, getAllSuffix) {
					getAllHandler(c, s.Backend)
					return
				}
				c.Next()
			})
		})
	}
}
