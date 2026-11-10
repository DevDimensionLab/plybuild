//go:build windows

package backlog

import "os"

// openSource opens a regular file read-only; the path must not end in a symlink.
func openSource(source string) (*os.File, error) {
	info, err := os.Lstat(source)
	if err != nil {
		return nil, openError(err)
	}
	if !info.Mode().IsRegular() {
		return nil, errNotRegular
	}
	file, err := os.Open(source)
	if err != nil {
		return nil, openError(err)
	}
	if info, err := file.Stat(); err != nil || !info.Mode().IsRegular() {
		_ = file.Close()
		if err != nil {
			return nil, openError(err)
		}
		return nil, errNotRegular
	}
	return file, nil
}

func syncDirectory(string) error { return nil }
