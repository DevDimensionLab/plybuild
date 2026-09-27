//go:build !windows

package workspace

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

func platformReplace(source, destination string) error {
	return os.Rename(source, destination)
}

func platformSyncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := directory.Sync(); err != nil {
		_ = directory.Close()
		return err
	}
	return directory.Close()
}
