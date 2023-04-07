package fs

import (
	"context"
	"github.com/skroczek/acme-restful/errors"
	"os"
	"strings"
	"time"
)

type Blob struct {
	Content []byte
	ModTime time.Time
}

type Memory struct {
	tree map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{
		tree: make(map[string]interface{}),
	}
}

func (m *Memory) getBlob(path string) (*Blob, error) {
	path = strings.Trim(path, "/")
	if len(path) < 6 {
		return nil, errors.ErrorInvalidPath
	}
	if !strings.HasSuffix(path, ".json") {
		return nil, errors.ErrorMissingExtension
	}
	parts := strings.Split(path, "/")
	tree := m.tree
	var ok bool
	for i := 0; i < (len(parts) - 1); i++ {
		if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
			return nil, os.ErrNotExist
		}
	}
	if blob, ok := tree[parts[len(parts)-1]].(*Blob); ok {
		return blob, nil
	}
	return nil, os.ErrNotExist
}

func (m *Memory) Exists(ctx context.Context, path string) (bool, error) {
	path = strings.Trim(path, "/")
	if len(path) < 6 {
		return false, errors.ErrorInvalidPath
	}
	if !strings.HasSuffix(path, ".json") {
		return false, errors.ErrorMissingExtension
	}
	blob, err := m.getBlob(path)
	if err != nil {
		return false, nil
	}
	return blob != nil, nil
}

func (m *Memory) Get(ctx context.Context, path string) ([]byte, error) {
	blob, err := m.getBlob(path)
	if err != nil {
		return nil, err
	}
	return blob.Content, nil
}

func (m *Memory) Write(ctx context.Context, path string, data []byte) error {
	path = strings.Trim(path, "/")
	if len(path) < 6 {
		return errors.ErrorInvalidPath
	}
	if !strings.HasSuffix(path, ".json") {
		return errors.ErrorMissingExtension
	}
	parts := strings.Split(path, "/")
	tree := m.tree
	var ok bool
	for i := 0; i < (len(parts) - 1); i++ {
		if _, ok = tree[parts[i]].(map[string]interface{}); !ok {
			tree[parts[i]] = make(map[string]interface{})
			tree, _ = tree[parts[i]].(map[string]interface{})
		} else {
			tree, _ = tree[parts[i]].(map[string]interface{})
		}
	}
	tree[parts[len(parts)-1]] = &Blob{Content: data, ModTime: time.Now()}
	return nil
}

func (m *Memory) Delete(ctx context.Context, path string) error {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	tree := m.tree
	var ok bool
	for i := 0; i < (len(parts) - 1); i++ {
		if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
			return os.ErrNotExist
		}
	}
	if _, ok := tree[parts[len(parts)-1]].(*Blob); ok {
		delete(tree, parts[len(parts)-1])
		return nil
	}
	return os.ErrNotExist
}

func (m *Memory) List(ctx context.Context, path string) ([]string, error) {
	path = strings.Trim(path, "/")
	tree := m.tree
	if path != "" {
		parts := strings.Split(path, "/")
		var ok bool
		for i := 0; i < (len(parts) - 1); i++ {
			if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
				return nil, os.ErrNotExist
			}
		}
		if tree, ok = tree[parts[len(parts)-1]].(map[string]interface{}); !ok {
			return nil, os.ErrNotExist
		}
	}
	var result []string
	for k, v := range tree {
		if _, ok := v.(*Blob); ok {
			result = append(result, k)
		}
	}
	return result, nil
}

func (m *Memory) GetLastModified(ctx context.Context, path string) (time.Time, error) {
	blob, err := m.getBlob(path)
	if err != nil {
		return time.Time{}, err
	}
	return blob.ModTime, nil
}
