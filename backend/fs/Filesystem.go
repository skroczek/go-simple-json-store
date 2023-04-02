package fs

import (
	"fmt"
	"io/fs"
	goos "os"
	"path/filepath"
	"syscall"
	"time"
)

type DeleteDirectoryError struct {
	Path string
}

func NewDeleteDirectoryError(path string) *DeleteDirectoryError {
	return &DeleteDirectoryError{Path: path}
}

func (d *DeleteDirectoryError) Error() string {
	return "cannot delete directory " + d.Path
}

type FilesystemOption func(*FilesystemBackend)

func WithDeleteEmptyDirs() FilesystemOption {
	return func(f *FilesystemBackend) {
		f.options |= deleteEmptyDirs
	}
}

func WithCreateDirs() FilesystemOption {
	return func(f *FilesystemBackend) {
		f.options |= createDirs
	}
}

type filesystemOption uint8

const (
	// FilesystemBackendConfig is the default configuration for the filesystem backend
	createDirs filesystemOption = 1 << iota
	deleteEmptyDirs
)

type FilesystemBackend struct {
	Root    string
	options filesystemOption
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
	if f.options&createDirs != 0 {
		dir := filepath.Dir(fullPath)
		if err := goos.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return goos.WriteFile(fullPath, data, 0644)
}

func (f FilesystemBackend) Delete(path string) error {
	if path == "" {
		return fmt.Errorf("cannot delete root")
	}
	fullPath := filepath.Join(f.Root, path)
	fileInfo, err := goos.Stat(fullPath)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() && f.options&deleteEmptyDirs == 0 {
		return NewDeleteDirectoryError(path)
	}
	if err := goos.Remove(fullPath); err != nil {
		if _, ok := err.(*DeleteDirectoryError); ok {
			return nil
		}
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
	path = filepath.Join(f.Root, path)
	var fileNames []string
	files, err := goos.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Type() == fs.ModeDir {
			continue
		}
		if filepath.Ext(file.Name()) == ".json" {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func (f FilesystemBackend) ListTypes(path string, mode fs.FileMode) ([]string, error) {
	path = filepath.Join(f.Root, path)
	list := make([]string, 0)
	files, err := goos.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Type() != mode {
			continue
		}
		list = append(list, file.Name())
	}
	return list, nil
}

func (f FilesystemBackend) GetLastModified(path string) (time.Time, error) {
	info, err := goos.Stat(filepath.Join(f.Root, path))
	return info.ModTime(), err
}

func NewFilesystemBackend(root string, options ...FilesystemOption) *FilesystemBackend {
	b := &FilesystemBackend{Root: root}
	for _, option := range options {
		option(b)
	}
	return b
}
