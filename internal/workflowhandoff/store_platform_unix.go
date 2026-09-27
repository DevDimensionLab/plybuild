//go:build !windows

package workflowhandoff

import (
	"os"

	"golang.org/x/sys/unix"
)

func storePlatformLock(file *os.File) error                 { return unix.Flock(int(file.Fd()), unix.LOCK_EX) }
func storePlatformUnlock(file *os.File) error               { return unix.Flock(int(file.Fd()), unix.LOCK_UN) }
func storePlatformReplace(source, destination string) error { return os.Rename(source, destination) }
func storeSyncDirectory(path string) error {
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
