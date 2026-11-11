//go:build windows

package inbox

import (
	"errors"
	"io/fs"
	"os"

	"golang.org/x/sys/windows"
)

func openLockFile(path string) (*os.File, error) {
	if info, err := os.Lstat(path); err == nil && !info.Mode().IsRegular() {
		return nil, errNotRegular
	}
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, fileMode)
}

func tryLock(file *os.File) (bool, error) {
	var overlapped windows.Overlapped
	err := windows.LockFileEx(windows.Handle(file.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &overlapped)
	if errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
		return false, nil
	}
	return err == nil, err
}

func unlockFile(file *os.File) error {
	var overlapped windows.Overlapped
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, 1, 0, &overlapped)
}

// openNoFollow opens path read-only; the path must not be a symlink.
func openNoFollow(path string) (*os.File, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fs.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errNotRegular
	}
	return os.Open(path)
}

func syncDirectory(string) error { return nil }
