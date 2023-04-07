package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const listAllSuffix = "__list.json"
const optionWithoutExtension = "withoutExtension"

func getListHandler(c *gin.Context, be backend.Backend) {
	urlPath := c.Request.URL.Path
	data, err := be.List(urlPath[0 : len(urlPath)-len(listAllSuffix)])
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	if _, ok := c.GetQuery(optionWithoutExtension); ok {
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
				if strings.HasSuffix(c.Request.URL.Path, listAllSuffix) {
					if c.Request.Method != http.MethodGet {
						_ = c.AbortWithError(http.StatusMethodNotAllowed, errMethodNotAllowed)
						return
					}
					getListHandler(c, s.Backend)
					return
				}
				c.Next()
			})
		})
	}
}
