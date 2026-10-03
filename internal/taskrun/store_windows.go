//go:build windows

package taskrun

import (
	"golang.org/x/sys/windows"
	"os"
)

func openRead(p string) (*os.File, error) { return physicalOpenFile(p, os.O_RDONLY, 0) }
func lockFile(p string) (*os.File, error) {
	f, e := physicalOpenFile(p, os.O_RDWR|os.O_CREATE, 0600)
	if e != nil {
		return nil, e
	}
	e = windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &windows.Overlapped{})
	if e != nil {
		f.Close()
		return nil, e
	}
	return f, nil
}
func unlockFile(f *os.File) error {
	e := windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, &windows.Overlapped{})
	c := f.Close()
	if e != nil {
		return e
	}
	return c
}
