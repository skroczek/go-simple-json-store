package fs

import (
	"fmt"
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
		return nil, fmt.Errorf("invalid path")
	}
	if !strings.HasSuffix(path, ".json") {
		return nil, fmt.Errorf("path must end with .json")
	}
	parts := strings.Split(path, "/")
	tree := m.tree
	var ok bool
	for i := 0; i < (len(parts) - 1); i++ {
		if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("invalid path")
		}
	}
	if blob, ok := tree[parts[len(parts)-1]].(*Blob); ok {
		return blob, nil
	}
	return nil, fmt.Errorf("invalid path")
}

func (m *Memory) Exists(path string) (bool, error) {
	path = strings.Trim(path, "/")
	if len(path) < 6 {
		return false, fmt.Errorf("invalid path")
	}
	if !strings.HasSuffix(path, ".json") {
		return false, fmt.Errorf("path must end with .json")
	}
	blob, err := m.getBlob(path)
	if err != nil {
		return false, nil
	}
	return blob != nil, nil
}

func (m *Memory) Get(path string) ([]byte, error) {
	blob, err := m.getBlob(path)
	if err != nil {
		return nil, err
	}
	return blob.Content, nil
}

func (m *Memory) Write(path string, data []byte) error {
	path = strings.Trim(path, "/")
	if len(path) < 6 {
		return fmt.Errorf("invalid path")
	}
	if !strings.HasSuffix(path, ".json") {
		return fmt.Errorf("path must end with .json")
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

func (m *Memory) Delete(path string) error {
	parts := strings.Split(path, "/")
	tree := m.tree
	var ok bool
	for i := 0; i < (len(parts) - 1); i++ {
		if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
			return fmt.Errorf("invalid path")
		}
	}
	if _, ok := tree[parts[len(parts)-1]].(*Blob); ok {
		delete(tree, parts[len(parts)-1])
		return nil
	}
	return fmt.Errorf("invalid path")
}

func (m *Memory) List(path string) ([]string, error) {
	path = strings.Trim(path, "/")
	if strings.HasSuffix(path, ".json") {
		return nil, fmt.Errorf("path must not end with .json")
	}
	tree := m.tree
	if path != "" {
		parts := strings.Split(path, "/")
		var ok bool
		for i := 0; i < (len(parts) - 1); i++ {
			if tree, ok = tree[parts[i]].(map[string]interface{}); !ok {
				return nil, fmt.Errorf("path does not exist")
			}
		}
		if tree, ok = tree[parts[len(parts)-1]].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("path does not exist")
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

func (m *Memory) GetLastModified(path string) (time.Time, error) {
	blob, err := m.getBlob(path)
	if err != nil {
		return time.Time{}, err
	}
	return blob.ModTime, nil
}
