package backend

import (
	"context"
	"github.com/skroczek/go-simple-json-store/helper"
	"time"
)

type Encrypted struct {
	Backend Backend
}

func NewEncrypted(backend Backend) *Encrypted {
	return &Encrypted{Backend: backend}
}

func (e *Encrypted) SetBackend(backend Backend) {
	e.Backend = backend
}

func (e *Encrypted) Exists(ctx context.Context, path string) (bool, error) {
	return e.Backend.Exists(ctx, path)
}

func (e *Encrypted) Get(ctx context.Context, path string) ([]byte, error) {
	// TODO: decrypt data
	return helper.Decrypt(e.Backend.Get(ctx, path))
}

func (e *Encrypted) Write(ctx context.Context, path string, data []byte) error {
	// TODO: encrypt data
	return e.Backend.Write(ctx, path, helper.Encrypt(data))
}

func (e *Encrypted) Delete(ctx context.Context, path string) error {
	return e.Backend.Delete(ctx, path)
}

func (e *Encrypted) List(ctx context.Context, path string) ([]string, error) {
	return e.Backend.List(ctx, path)
}

func (e *Encrypted) GetLastModified(ctx context.Context, path string) (time.Time, error) {
	return e.Backend.GetLastModified(ctx, path)
}
