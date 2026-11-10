//go:build !windows

package backlog

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// openSource opens a regular file read-only without following a final symlink
// or blocking on a FIFO.
func openSource(source string) (*os.File, error) {
	fd, err := unix.Open(source, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if errors.Is(err, unix.ELOOP) || errors.Is(err, unix.ENXIO) || errors.Is(err, unix.EOPNOTSUPP) {
		return nil, errNotRegular
	}
	if err != nil {
		return nil, openError(err)
	}
	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		_ = unix.Close(fd)
		return nil, openError(err)
	}
	if stat.Mode&unix.S_IFMT != unix.S_IFREG {
		_ = unix.Close(fd)
		return nil, errNotRegular
	}
	if err := unix.SetNonblock(fd, false); err != nil {
		_ = unix.Close(fd)
		return nil, openError(err)
	}
	return os.NewFile(uintptr(fd), source), nil
}

func syncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	err = directory.Sync()
	if closeErr := directory.Close(); err == nil {
		err = closeErr
	}
	return err
}
