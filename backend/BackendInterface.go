package backend

import (
	"context"
	"io/fs"
	"time"
)

type Backend interface {
	Exists(ctx context.Context, path string) (bool, error)
	Get(ctx context.Context, path string) ([]byte, error)
	Write(ctx context.Context, path string, data []byte) error
	Delete(ctx context.Context, path string) error
	List(ctx context.Context, path string) ([]string, error)
	GetLastModified(ctx context.Context, path string) (time.Time, error)
}

type FileBackend interface {
	Backend
	ListTypes(ctx context.Context, path string, mode fs.FileMode) ([]string, error)
}

type Proxy interface {
	SetBackend(backend Backend)
}
