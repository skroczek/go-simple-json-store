package backend

import (
	"time"
)

type Backend interface {
	Exists(path string) (bool, error)
	Get(path string) ([]byte, error)
	Write(path string, data []byte) error
	Delete(path string) error
	List(path string) ([]string, error)
	GetLastModified(path string) (time.Time, error)
}

type Proxy interface {
	SetBackend(backend Backend)
}
