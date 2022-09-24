package backend

import (
	"time"
)

type Backend interface {
	Exists(path string) (bool, error)
	Get(path string) (interface{}, error)
	Write(path string, object interface{}) error
	Delete(path string) error
	List(path string) ([]string, error)
	GetLastModified(path string) (time.Time, error)
}
