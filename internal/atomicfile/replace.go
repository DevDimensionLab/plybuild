// Package atomicfile replaces files with complete bytes; see README.md for guarantees.
package atomicfile

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Error records the failed operation and whether replacement returned success.
// Replaced=false does not establish absence of external effects; see README.md.
type Error struct {
	Op       string
	Path     string
	Replaced bool
	Err      error
}

func (e *Error) Error() string { return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err) }
func (e *Error) Unwrap() error { return e.Err }

// Replace writes, chmods, syncs and closes a same-directory temp, replaces the
// destination, then syncs its directory (a no-op on Windows). Failures are *Error.
func Replace(destination string, contents []byte, mode fs.FileMode) error {
	return replace(destination, contents, mode, systemOperations())
}

// operations is a private, per-call fault seam; production has no mutable hooks.
type operations struct {
	createTemp    func(string, string) (*os.File, error)
	write         func(*os.File, []byte) (int, error)
	chmod         func(*os.File, fs.FileMode) error
	sync          func(*os.File) error
	close         func(*os.File) error
	replace       func(string, string) error
	syncDirectory func(string) error
	remove        func(string) error
}

func systemOperations() operations {
	return operations{
		createTemp: os.CreateTemp,
		write:      func(f *os.File, b []byte) (int, error) { return f.Write(b) },
		chmod:      func(f *os.File, m fs.FileMode) error { return f.Chmod(m) },
		sync:       func(f *os.File) error { return f.Sync() },
		close:      func(f *os.File) error { return f.Close() },
		replace:    platformReplace, syncDirectory: syncDirectory, remove: os.Remove,
	}
}

func replace(destination string, contents []byte, mode fs.FileMode, ops operations) error {
	directory := filepath.Dir(destination)
	temporary, err := ops.createTemp(directory, ".atomicfile-*.tmp")
	if err != nil {
		return &Error{Op: "open temp", Path: directory, Err: err}
	}
	path := temporary.Name()
	closed := false
	abort := func(op, errorPath string, primary error) error {
		if !closed {
			if err := ops.close(temporary); err != nil && !errors.Is(err, os.ErrClosed) {
				primary = fmt.Errorf("%w; close temp: %v", primary, err)
			}
		}
		if err := ops.remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			primary = fmt.Errorf("%w; cleanup %s: %v", primary, path, err)
		}
		return &Error{Op: op, Path: errorPath, Err: primary}
	}
	written, err := ops.write(temporary, contents)
	if err != nil {
		return abort("write", path, err)
	}
	if written != len(contents) {
		return abort("write", path, io.ErrShortWrite)
	}
	if err := ops.chmod(temporary, mode); err != nil {
		return abort("chmod", path, err)
	}
	if err := ops.sync(temporary); err != nil {
		return abort("sync", path, err)
	}
	if err := ops.close(temporary); err != nil {
		return abort("close", path, err)
	}
	closed = true
	if err := ops.replace(path, destination); err != nil {
		return abort("replace", destination, err)
	}
	if err := ops.syncDirectory(directory); err != nil {
		return &Error{Op: "sync directory", Path: directory, Replaced: true, Err: err}
	}
	return nil
}
