//go:build !windows

package weekrun

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

var errOwned = errors.New("week controller is already owned")

func lock(path string) (*os.File, error) {
	if err := physical(path); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|unix.O_NOFOLLOW, 0600)
	if err != nil {
		return nil, err
	}
	i, err := f.Stat()
	if err != nil || !i.Mode().IsRegular() || i.Mode().Perm()&0077 != 0 {
		f.Close()
		return nil, fmt.Errorf("invalid week lock")
	}
	err = unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		f.Close()
		if errors.Is(err, unix.EWOULDBLOCK) {
			return nil, errOwned
		}
		return nil, err
	}
	return f, nil
}
func ownerState(dir string) string {
	f, err := lock(dir + "/owner.lock")
	if err == nil {
		f.Close()
		return "unowned"
	}
	if errors.Is(err, errOwned) {
		return "lock_held"
	}
	return "unknown"
}
