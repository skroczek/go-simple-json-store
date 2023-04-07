package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const listDirSuffix = "__dir.json"

func getListDirHandler(c *gin.Context, be backend.FileBackend) {
	urlPath := c.Request.URL.Path
	data, err := be.ListTypes(c, urlPath[0:len(urlPath)-len(listDirSuffix)], fs.ModeDir)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.AbortWithStatusJSON(http.StatusOK, data)
}

func WithListDir() Options {
	return func(s *Server) {
		if b, ok := s.Backend.(backend.FileBackend); ok {
			s.AddRouterOption(func(r *gin.Engine) {
				r.Use(func(c *gin.Context) {
					if c.Request.Method == http.MethodGet && strings.HasSuffix(c.Request.URL.Path, listDirSuffix) {
						getListDirHandler(c, b)
						return
					}
					c.Next()
				})
			})
		} else {
			log.Panicf("Error: backend does not implement backend.FileBackend")
		}
	}
}
