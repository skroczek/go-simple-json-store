package backend

import (
	"io/fs"
	goos "os"
	"path/filepath"
	"syscall"
	"time"
)

type FilesystemBackend struct {
	Root string
}

func (f FilesystemBackend) Exists(path string) (bool, error) {
	stat, err := goos.Stat(filepath.Join(f.Root, path))
	if err != nil {
		return false, err
	}
	return !stat.IsDir(), nil
}

func (f FilesystemBackend) Get(path string) ([]byte, error) {
	return goos.ReadFile(filepath.Join(f.Root, path))
}

func (f FilesystemBackend) Write(path string, data []byte) error {
	fullPath := filepath.Join(f.Root, path)
	dir := filepath.Dir(fullPath)
	if err := goos.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return goos.WriteFile(fullPath, data, 0644)
}

func (f FilesystemBackend) Delete(path string) error {
	fullPath := filepath.Join(f.Root, path)
	err := goos.Remove(fullPath)
	if err != nil {
		if err, ok := err.(*goos.PathError); ok {
			// TODO: we need some windows specific code here
			if err.Err == syscall.ENOTEMPTY {
				// ignore not empty error
				return nil
			}
		}
		return err
	}
	parentPath := filepath.Dir(fullPath)
	if parentPath != f.Root {
		return f.Delete(parentPath[len(f.Root)+1:])
	}
	return nil
}

func (f FilesystemBackend) List(path string) ([]string, error) {
	path = filepath.Dir(filepath.Join(f.Root, path))
	var files []string
	err := filepath.WalkDir(path, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".json" {
			files = append(files, d.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (f FilesystemBackend) GetLastModified(path string) (time.Time, error) {
	info, err := goos.Stat(filepath.Join(f.Root, path))
	return info.ModTime(), err
}

func NewFilesystemBackend(root string) *FilesystemBackend {
	return &FilesystemBackend{Root: root}
}
