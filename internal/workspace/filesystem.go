package workspace

import (
	"io/fs"
	"os"
	"path/filepath"
)

// File is the smallest file handle needed while publishing a workspace marker.
type File interface {
	Write([]byte) (int, error)
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

// FileSystem isolates all filesystem effects used by Init.
type FileSystem interface {
	Getwd() (string, error)
	EvalSymlinks(string) (string, error)
	Lstat(string) (fs.FileInfo, error)
	Stat(string) (fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
	MkdirTemp(string, string) (string, error)
	OpenFile(string, int, fs.FileMode) (File, error)
	Chmod(string, fs.FileMode) error
	Rename(string, string) error
	RemoveAll(string) error
}

type systemFileSystem struct{}

func (systemFileSystem) Getwd() (string, error) { return os.Getwd() }

func (systemFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (systemFileSystem) Lstat(path string) (fs.FileInfo, error) { return os.Lstat(path) }

func (systemFileSystem) Stat(path string) (fs.FileInfo, error) { return os.Stat(path) }

func (systemFileSystem) ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }

func (systemFileSystem) MkdirTemp(directory, pattern string) (string, error) {
	return os.MkdirTemp(directory, pattern)
}

func (systemFileSystem) OpenFile(path string, flag int, mode fs.FileMode) (File, error) {
	return os.OpenFile(path, flag, mode)
}

func (systemFileSystem) Chmod(path string, mode fs.FileMode) error {
	return os.Chmod(path, mode)
}

func (systemFileSystem) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (systemFileSystem) RemoveAll(path string) error { return os.RemoveAll(path) }
