//go:build !windows

package workspace

import "os"

func PlatformReplace(source, destination string) error {
	return os.Rename(source, destination)
}

func PlatformSyncDirectory(path string) error {
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
