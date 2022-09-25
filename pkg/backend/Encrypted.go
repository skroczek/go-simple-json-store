package backend

import (
	"github.com/skroczek/acme-restful/internal/helper"
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

func (e *Encrypted) Exists(path string) (bool, error) {
	return e.Backend.Exists(path)
}

func (e *Encrypted) Get(path string) ([]byte, error) {
	// TODO: decrypt data
	return helper.Decrypt(e.Backend.Get(path))
}

func (e *Encrypted) Write(path string, data []byte) error {
	// TODO: encrypt data
	return e.Backend.Write(path, helper.Encrypt(data))
}

func (e *Encrypted) Delete(path string) error {
	return e.Backend.Delete(path)
}

func (e *Encrypted) List(path string) ([]string, error) {
	return e.Backend.List(path)
}

func (e *Encrypted) GetLastModified(path string) (time.Time, error) {
	return e.Backend.GetLastModified(path)
}
