//go:build !windows

package taskrun

import (
	"golang.org/x/sys/unix"
	"os"
)

func openRead(p string) (*os.File, error) { return os.OpenFile(p, os.O_RDONLY|unix.O_NOFOLLOW, 0) }
func lockFile(p string) (*os.File, error) {
	f, e := os.OpenFile(p, os.O_RDWR|os.O_CREATE|unix.O_NOFOLLOW, 0600)
	if e != nil {
		return nil, e
	}
	info, e := f.Stat()
	if e != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 {
		f.Close()
		return nil, integrity("invalid lock file")
	}
	if e = unix.Flock(int(f.Fd()), unix.LOCK_EX); e != nil {
		f.Close()
		return nil, e
	}
	return f, nil
}
func unlockFile(f *os.File) error {
	e := unix.Flock(int(f.Fd()), unix.LOCK_UN)
	c := f.Close()
	if e != nil {
		return e
	}
	return c
}
