// Package filesystem isolates the filesystem operations used by migrated callers.
package filesystem

import (
	"errors"
	"io/fs"
	"os"
)

// ErrNoFilesystem reports a zero-value dependency without performing an
// operating-system filesystem operation.
var ErrNoFilesystem = errors.New("filesystem dependency is not configured")

// FileSystem performs the filesystem operations used by a migrated flow.
type FileSystem interface {
	ReadFile(string) ([]byte, error)
	Stat(string) (fs.FileInfo, error)
	MkdirAll(string, fs.FileMode) error
	WriteFile(string, []byte, fs.FileMode) error
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
