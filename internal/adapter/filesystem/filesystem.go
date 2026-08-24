// Package filesystem isolates the filesystem operations used by migrated callers.
package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"
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
	Stat(string) (fs.FileInfo, error)
	MkdirAll(string, fs.FileMode) error
	WriteFile(string, []byte, fs.FileMode) error
	OpenFile(string, int, fs.FileMode) (*os.File, error)
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

func (systemFilesystem) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
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

func (systemFilesystem) Create(path string) (File, error) {
	return os.Create(path)
}

func (systemFilesystem) Copy(destination File, source io.Reader) (int64, error) {
	return io.Copy(destination, source)
}
