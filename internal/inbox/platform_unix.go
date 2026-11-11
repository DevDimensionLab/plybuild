//go:build !windows

package inbox

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

// openLockFile opens or creates the lock file without following a symlink.
func openLockFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|unix.O_NOFOLLOW|unix.O_CLOEXEC, fileMode)
}

// tryLock takes the exclusive lock without waiting; false means another holder has it.
func tryLock(file *os.File) (bool, error) {
	err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if errors.Is(err, unix.EWOULDBLOCK) {
		return false, nil
	}
	return err == nil, err
}

func unlockFile(file *os.File) error { return unix.Flock(int(file.Fd()), unix.LOCK_UN) }

// openNoFollow opens path read-only without following a final symlink or blocking on a FIFO.
func openNoFollow(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if errors.Is(err, unix.ELOOP) || errors.Is(err, unix.ENXIO) || errors.Is(err, unix.EOPNOTSUPP) {
		return nil, errNotRegular
	}
	return file, err
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
