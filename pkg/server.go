package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/internal/helper"
	"github.com/skroczek/acme-restful/pkg/backend"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	Backend backend.Backend
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
	data, err := s.Backend.Get(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNotFound)
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
			c.Status(http.StatusNotFound)
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
			c.Status(http.StatusNotFound)
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
		c.Status(http.StatusBadRequest)
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
			c.Status(http.StatusNotFound)
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
			c.Status(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	patchData, _ := helper.FromJSON(io.ReadAll(c.Request.Body))
	if patchDataMap, ok := patchData.(map[string]interface{}); ok {
		if dataMap, ok := object.(map[string]interface{}); ok {
			object = helper.MergeMap(dataMap, patchDataMap)
		} else {
			c.Status(http.StatusBadRequest)
			return
		}
	} else if patchDataSlice, ok := patchData.([]interface{}); ok {
		if dataSlice, ok := object.([]interface{}); ok {
			object = append(dataSlice, patchDataSlice...)
		} else {
			c.Status(http.StatusBadRequest)
			return
		}
	} else {
		c.Status(http.StatusBadRequest)
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
	c.Status(http.StatusNotFound)
}

func NewServer(backend backend.Backend) *Server {
	return &Server{
		Backend: backend,
	}
}
