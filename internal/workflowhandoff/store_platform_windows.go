//go:build windows

package workflowhandoff

import (
	"os"

	"golang.org/x/sys/windows"
)

func storePlatformLock(file *os.File) error {
	var overlapped windows.Overlapped
	return windows.LockFileEx(windows.Handle(file.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &overlapped)
}
func storePlatformUnlock(file *os.File) error {
	var overlapped windows.Overlapped
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, 1, 0, &overlapped)
}
func storeSyncDirectory(string) error { return nil }
