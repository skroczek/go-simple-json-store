package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/skroczek/acme-restful/internal/helper"
	"github.com/skroczek/acme-restful/pkg/backend"
	"github.com/skroczek/acme-restful/pkg/router"
)

type Server struct {
	Backend       backend.Backend
	routerOptions []router.Option
}

func (s *Server) GetHandler(c *gin.Context) {
	path := c.Params.ByName("path")
	if strings.HasSuffix(path, "__all.json") {
		s.getAllHandler(c)
		return
	}
	if strings.HasSuffix(path, "__list.json") {
		s.getListHandler(c)
		return
	}
	data, err := helper.FromJSON(s.Backend.Get(path))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	modTime, _ := s.Backend.GetLastModified(path)
	c.Header("Last-Modified", modTime.Format(time.RFC1123))
	c.JSON(http.StatusOK, data)
}

func (s *Server) getAllHandler(c *gin.Context) {
	//data, err := cos.ReadAll(root, c.Params.ByName("path"))
	path := c.Params.ByName("path")[0 : len(c.Params.ByName("path"))-len("__all.json")]
	list, err := s.Backend.List(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	data := make([]interface{}, len(list))
	for i := range list {
		obj, err := s.Backend.Get(filepath.Join(path, list[i]))
		if err != nil {
			log.Panicf("Error: %+v", err)
		}
		data[i] = obj
	}
	c.JSON(http.StatusOK, data)
}

func (s *Server) getListHandler(c *gin.Context) {
	//data, err := cos.List(root, c.Params.ByName("path"))
	data, err := s.Backend.List(c.Params.ByName("path"))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.JSON(http.StatusOK, data)
}

func (s *Server) PostHandler(c *gin.Context) {
	s.PutHandler(c)
}

func (s *Server) PutHandler(c *gin.Context) {
	data, err := helper.FromJSON(io.ReadAll(c.Request.Body))
	if err != nil {
		log.Printf("Error: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = s.Backend.Write(c.Params.ByName("path"), helper.ToJSON(data))
	if err != nil {
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusCreated)
}

func (s *Server) DeleteHandler(c *gin.Context) {
	err := s.Backend.Delete(c.Params.ByName("path"))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) PatchHandler(c *gin.Context) {
	path := c.Params.ByName("path")
	object, err := helper.FromJSON(s.Backend.Get(path))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	patchData, _ := helper.FromJSON(io.ReadAll(c.Request.Body))
	if patchDataMap, ok := patchData.(map[string]interface{}); ok {
		if dataMap, ok := object.(map[string]interface{}); ok {
			object = helper.MergeMap(dataMap, patchDataMap)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else if patchDataSlice, ok := patchData.([]interface{}); ok {
		if dataSlice, ok := object.([]interface{}); ok {
			object = append(dataSlice, patchDataSlice...)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = s.Backend.Write(path, helper.ToJSON(object))
	if err != nil {
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusCreated)
}

func (s *Server) HeadHandler(c *gin.Context) {
	s.GetHandler(c)
}

func (s *Server) OptionsHandler(c *gin.Context) {
	isFile, err := s.Backend.Exists(c.Params.ByName("path"))
	if err != nil {
		log.Panicf("Error: %+v", err)
	}
	if isFile {
		c.Header("Allow", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

type Options func(*Server)

func WithBackend(be backend.Backend) Options {
	return func(s *Server) {
		if pbe, ok := be.(backend.Proxy); ok {
			pbe.SetBackend(s.Backend)
		}
		s.Backend = be
	}
}

func WithRouterOptions(opts ...router.Option) Options {
	return func(s *Server) {
		s.routerOptions = opts
	}
}

func NewServer(opts ...Options) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Run() {
	r := router.DefaultRouter(s.routerOptions...)

	r.GET("/*path", s.GetHandler)
	r.POST("/*path", s.PostHandler)
	r.PUT("/*path", s.PutHandler)
	r.DELETE("/*path", s.DeleteHandler)
	r.PATCH("/*path", s.PatchHandler)
	r.HEAD("/*path", s.HeadHandler)
	r.OPTIONS("/*path", s.OptionsHandler)

	r.Run()
}
