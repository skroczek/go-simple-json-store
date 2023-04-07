package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skroczek/acme-restful/backend/fs"
	"github.com/skroczek/acme-restful/errors"
	"github.com/skroczek/acme-restful/helper"
	"io"
	"net/http"
	"os"
	"time"
)

// GetHandler handles GET requests
func (s *Server) GetHandler(c *gin.Context) {
	path := c.Request.URL.Path
	data, err := helper.FromJSON(s.Backend.Get(c, path))
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if errors.IsClientError(err) {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	modTime, _ := s.Backend.GetLastModified(c, path)
	c.Header("Last-Modified", modTime.Format(time.RFC1123))
	c.JSON(http.StatusOK, data)
}

// PostHandler handles POST requests
func (s *Server) PostHandler(c *gin.Context) {
	s.PutHandler(c)
}

// PutHandler handles PUT requests
func (s *Server) PutHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	data, err := helper.FromJSON(io.ReadAll(c.Request.Body))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = s.Backend.Write(c, urlPath, helper.ToJSON(data))
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if errors.IsClientError(err) {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusCreated)
}

// DeleteHandler handles DELETE requests
func (s *Server) DeleteHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	err := s.Backend.Delete(c, urlPath)
	if err != nil {
		if _, ok := err.(*fs.DeleteDirectoryError); ok {
			_ = c.AbortWithError(http.StatusMethodNotAllowed, err)
			return
		}
		if os.IsNotExist(err) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if errors.IsClientError(err) {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// PatchHandler handles PATCH requests
func (s *Server) PatchHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	object, err := helper.FromJSON(s.Backend.Get(c, urlPath))
	if err != nil {
		if os.IsNotExist(err) {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	patchData, _ := helper.FromJSON(io.ReadAll(c.Request.Body))
	if patchDataMap, ok := patchData.(map[string]interface{}); ok {
		if dataMap, ok := object.(map[string]interface{}); ok {
			object = helper.MergeMap(dataMap, patchDataMap)
		} else {
			// TODO: maybe replace original object with patchDataMap?
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unable to merge map with %T", object))
			return
		}
	} else if patchDataSlice, ok := patchData.([]interface{}); ok {
		if dataSlice, ok := object.([]interface{}); ok {
			object = append(dataSlice, patchDataSlice...)
		} else {
			// TODO: maybe replace original object with patchDataSlice?
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unable to merge slice with %T", object))
			return
		}
	} else {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("unable to merge %T with %T", object, patchData))
		return
	}
	err = s.Backend.Write(c, urlPath, helper.ToJSON(object))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusCreated)
}

// HeadHandler handles HEAD requests
func (s *Server) HeadHandler(c *gin.Context) {
	s.GetHandler(c)
}

// OptionsHandler handles OPTIONS requests
func (s *Server) OptionsHandler(c *gin.Context) {
	urlPath := c.Request.URL.Path
	isFile, err := s.Backend.Exists(c, urlPath)
	if err != nil {
		if errors.IsClientError(err) {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if isFile {
		c.Header("Allow", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
		c.Status(http.StatusOK)
		return
	}
	_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("path %s does not exist", urlPath))
}
