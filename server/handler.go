package server

import (
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend/fs"
	"github.com/skroczek/acme-restful/errors"
	helper2 "github.com/skroczek/acme-restful/helper"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (s *Server) GetHandler(c *gin.Context) {
	path := c.Request.URL.Path
	data, err := helper2.FromJSON(s.Backend.Get(c, path))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.IsClientError(err) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	modTime, _ := s.Backend.GetLastModified(c, path)
	c.Header("Last-Modified", modTime.Format(time.RFC1123))
	c.JSON(http.StatusOK, data)
}

func (s *Server) PostHandler(c *gin.Context) {
	s.PutHandler(c)
}

func (s *Server) PutHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	data, err := helper2.FromJSON(io.ReadAll(c.Request.Body))
	if err != nil {
		log.Printf("Error: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = s.Backend.Write(c, urlPath, helper2.ToJSON(data))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.IsClientError(err) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusCreated)
}

func (s *Server) DeleteHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	err := s.Backend.Delete(c, urlPath)
	if err != nil {
		if _, ok := err.(*fs.DeleteDirectoryError); ok {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if errors.IsClientError(err) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) PatchHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	object, err := helper2.FromJSON(s.Backend.Get(c, urlPath))
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	patchData, _ := helper2.FromJSON(io.ReadAll(c.Request.Body))
	if patchDataMap, ok := patchData.(map[string]interface{}); ok {
		if dataMap, ok := object.(map[string]interface{}); ok {
			object = helper2.MergeMap(dataMap, patchDataMap)
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
	err = s.Backend.Write(c, urlPath, helper2.ToJSON(object))
	if err != nil {
		log.Panicf("Error: %+v", err)
	}
	c.Status(http.StatusCreated)
}

func (s *Server) HeadHandler(c *gin.Context) {
	s.GetHandler(c)
}

func (s *Server) OptionsHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	isFile, err := s.Backend.Exists(c, urlPath)
	if err != nil {
		if errors.IsClientError(err) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Panicf("Error: %+v", err)
	}
	if isFile {
		c.Header("Allow", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
		c.Status(http.StatusOK)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}
