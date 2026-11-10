// Package filelock holds one exclusive advisory lock on a physical file.
package filelock

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Options.Fail is a test seam; ops are lock-open, lock, unlock and lock-close.
type Options struct {
	Fail func(op string) error

	beforeIdentity func(*os.File) // package-test seam between open and identity reads
}

type Lock struct {
	file *os.File
	fail func(string) error
}

// Error records the failed step: inspect, not-regular, open, stat-opened,
// lstat-path, changed or lock. Injected marks an Options.Fail fault.
type Error struct {
	Op       string
	Injected bool
	Path     string
	Err      error
}

func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s %s", e.Op, e.Path)
	}
	return fmt.Sprintf("%s %s: %v", e.Op, e.Path, e.Err)
}
func (e *Error) Unwrap() error { return e.Err }

// ReleaseError keeps both release steps' failures in execution order.
type ReleaseError struct {
	Path                          string
	Unlock, Close                 error
	UnlockInjected, CloseInjected bool
}

func (e *ReleaseError) Error() string {
	return fmt.Sprintf("release %s: unlock: %v; close: %v", e.Path, e.Unlock, e.Close)
}

func (e *ReleaseError) Unwrap() []error {
	var errs []error
	for _, err := range []error{e.Unlock, e.Close} {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// Acquire opens or creates path as a regular file and locks it exclusively.
func Acquire(path string, o Options) (*Lock, error) {
	fail := o.Fail
	if fail == nil {
		fail = func(string) error { return nil }
	}
	for {
		info, err := os.Lstat(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, &Error{Op: "inspect", Path: path, Err: err}
		}
		if err == nil && (!info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0) {
			return nil, &Error{Op: "not-regular", Path: path}
		}
		flags := os.O_RDWR
		if errors.Is(err, fs.ErrNotExist) {
			flags |= os.O_CREATE | os.O_EXCL
		}
		if fault := fail("lock-open"); fault != nil {
			return nil, &Error{Op: "open", Injected: true, Path: path, Err: fault}
		}
		file, err := os.OpenFile(path, flags, 0o600)
		if err != nil && flags&os.O_EXCL != 0 && errors.Is(err, fs.ErrExist) {
			continue
		}
		if err != nil {
			return nil, &Error{Op: "open", Path: path, Err: err}
		}
		if err := lockOpened(file, path, o.beforeIdentity, fail); err != nil {
			return nil, err
		}
		return &Lock{file: file, fail: fail}, nil
	}
}

// AcquireExisting locks an existing regular file without creating it. The
// os.Open error is returned unwrapped.
func AcquireExisting(path string) (*Lock, error) {
	return acquireExisting(path, nil)
}

func acquireExisting(path string, beforeIdentity func(*os.File)) (*Lock, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	noFault := func(string) error { return nil }
	if err := lockOpened(file, path, beforeIdentity, noFault); err != nil {
		return nil, err
	}
	return &Lock{file: file, fail: noFault}, nil
}

// lockOpened verifies that file is still the regular file at path, then locks
// it. The file is closed on every failure.
func lockOpened(file *os.File, path string, beforeIdentity func(*os.File), fail func(string) error) *Error {
	lockErr := func() *Error {
		if beforeIdentity != nil {
			beforeIdentity(file)
		}
		opened, err := file.Stat()
		if err != nil {
			return &Error{Op: "stat-opened", Path: path, Err: err}
		}
		current, err := os.Lstat(path)
		if err != nil {
			return &Error{Op: "lstat-path", Path: path, Err: err}
		}
		if !opened.Mode().IsRegular() || !current.Mode().IsRegular() || current.Mode()&fs.ModeSymlink != 0 || !os.SameFile(opened, current) {
			return &Error{Op: "changed", Path: path}
		}
		if fault := fail("lock"); fault != nil {
			return &Error{Op: "lock", Injected: true, Path: path, Err: fault}
		}
		if err := platformLock(file); err != nil {
			return &Error{Op: "lock", Path: path, Err: err}
		}
		return nil
	}()
	if lockErr != nil {
		_ = file.Close()
	}
	return lockErr
}

// Path is the path the lock was acquired with.
func (l *Lock) Path() string { return l.file.Name() }

// Release always unlocks and then closes. A fault replaces only a successful step.
func (l *Lock) Release() *ReleaseError {
	result := ReleaseError{Path: l.file.Name(), Unlock: platformUnlock(l.file)}
	if fault := l.fail("unlock"); fault != nil && result.Unlock == nil {
		result.Unlock, result.UnlockInjected = fault, true
	}
	result.Close = l.file.Close()
	if fault := l.fail("lock-close"); fault != nil && result.Close == nil {
		result.Close, result.CloseInjected = fault, true
	}
	if result.Unlock == nil && result.Close == nil {
		return nil
	}
	return &result
}
