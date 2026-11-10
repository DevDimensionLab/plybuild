//go:build !windows

package filelock

import (
	"os"

	"golang.org/x/sys/unix"
)

func platformLock(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_EX)
}

func platformUnlock(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
