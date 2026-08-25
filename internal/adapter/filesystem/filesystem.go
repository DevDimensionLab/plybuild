// Package filesystem isolates the filesystem operations used by migrated callers.
package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ErrNoFilesystem reports a zero-value dependency without performing an
// operating-system filesystem operation.
var ErrNoFilesystem = errors.New("filesystem dependency is not configured")

// File is a created file that can receive bytes and be closed by its caller.
type File interface {
	io.Writer
	io.Closer
}

// FileSystem performs the filesystem operations used by a migrated flow.
type FileSystem interface {
	ReadFile(string) ([]byte, error)
	ReadDir(string) ([]fs.FileInfo, error)
	Stat(string) (fs.FileInfo, error)
	Mkdir(string, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
	WriteFile(string, []byte, fs.FileMode) error
	OpenFile(string, int, fs.FileMode) (*os.File, error)
	Remove(string) error
	RemoveAll(string) error
	Rename(string, string) error
	Glob(string) ([]string, error)
	Walk(string, filepath.WalkFunc) error
	Create(string) (File, error)
	Copy(File, io.Reader) (int64, error)
}

// Dependencies contains the filesystem effect used by a caller. Its zero
// value returns ErrNoFilesystem without mutating the filesystem; production
// callers must select System explicitly.
type Dependencies struct {
	FileSystem FileSystem
}

// ReadFile passes the complete source path to the configured dependency.
func ReadFile(dependencies Dependencies, path string) ([]byte, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.ReadFile(path)
}

// ReadDir passes the complete directory path to the configured dependency.
func ReadDir(dependencies Dependencies, path string) ([]fs.FileInfo, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.ReadDir(path)
}

// Stat passes the complete path to the configured dependency.
func Stat(dependencies Dependencies, path string) (fs.FileInfo, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.Stat(path)
}

// Exists preserves the legacy existence rule: only an IsNotExist result means
// absent. A zero-value dependency is absent so callers reach a safe operation
// that returns ErrNoFilesystem before mutation.
func Exists(dependencies Dependencies, path string) bool {
	if dependencies.FileSystem == nil {
		return false
	}
	_, err := dependencies.FileSystem.Stat(path)
	return !os.IsNotExist(err)
}

// Mkdir passes the complete directory path and mode to the configured dependency.
func Mkdir(dependencies Dependencies, path string, mode fs.FileMode) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.Mkdir(path, mode)
}

// MkdirAll passes the complete directory path and mode to the configured dependency.
func MkdirAll(dependencies Dependencies, path string, mode fs.FileMode) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.MkdirAll(path, mode)
}

// WriteFile passes the complete destination path, bytes, and mode to the configured dependency.
func WriteFile(dependencies Dependencies, path string, data []byte, mode fs.FileMode) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.WriteFile(path, data, mode)
}

// OpenFile passes the complete path, flags, and mode to the configured dependency.
func OpenFile(dependencies Dependencies, path string, flags int, mode fs.FileMode) (*os.File, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.OpenFile(path, flags, mode)
}

// Remove passes the complete path to the configured dependency.
func Remove(dependencies Dependencies, path string) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.Remove(path)
}

// RemoveAll passes the complete path to the configured dependency.
func RemoveAll(dependencies Dependencies, path string) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.RemoveAll(path)
}

// Rename passes the complete source and destination paths to the configured dependency.
func Rename(dependencies Dependencies, source, destination string) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.Rename(source, destination)
}

// Glob passes the complete pattern to the configured dependency.
func Glob(dependencies Dependencies, pattern string) ([]string, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.Glob(pattern)
}

// Walk passes the complete root and callback to the configured dependency.
func Walk(dependencies Dependencies, root string, callback filepath.WalkFunc) error {
	if dependencies.FileSystem == nil {
		return ErrNoFilesystem
	}
	return dependencies.FileSystem.Walk(root, callback)
}

// Create passes the complete destination path to the configured dependency.
func Create(dependencies Dependencies, path string) (File, error) {
	if dependencies.FileSystem == nil {
		return nil, ErrNoFilesystem
	}
	return dependencies.FileSystem.Create(path)
}

// Copy passes the created destination and response body to the configured dependency.
func Copy(dependencies Dependencies, destination File, source io.Reader) (int64, error) {
	if dependencies.FileSystem == nil {
		return 0, ErrNoFilesystem
	}
	return dependencies.FileSystem.Copy(destination, source)
}

// System returns the production dependency that uses the operating-system filesystem.
func System() Dependencies {
	return Dependencies{FileSystem: systemFilesystem{}}
}

type systemFilesystem struct{}

func (systemFilesystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (systemFilesystem) ReadDir(path string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(path)
}

func (systemFilesystem) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (systemFilesystem) Mkdir(path string, mode fs.FileMode) error {
	return os.Mkdir(path, mode)
}

func (systemFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (systemFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	return os.WriteFile(path, data, mode)
}

func (systemFilesystem) OpenFile(path string, flags int, mode fs.FileMode) (*os.File, error) {
	return os.OpenFile(path, flags, mode)
}

func (systemFilesystem) Remove(path string) error {
	return os.Remove(path)
}

func (systemFilesystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (systemFilesystem) Rename(source, destination string) error {
	return os.Rename(source, destination)
}

func (systemFilesystem) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (systemFilesystem) Walk(root string, callback filepath.WalkFunc) error {
	return filepath.Walk(root, callback)
}

func (systemFilesystem) Create(path string) (File, error) {
	return os.Create(path)
}

func (systemFilesystem) Copy(destination File, source io.Reader) (int64, error) {
	return io.Copy(destination, source)
}
