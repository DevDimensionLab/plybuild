//go:build !windows

package workflownotification

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// Open every component relative to a pinned directory descriptor. O_NOFOLLOW on
// the leaf alone would still allow a parent symlink to redirect an effect.
func openDirectory(path string) (*os.File, error) {
	if !validPath(path) {
		return nil, invalid()
	}
	fd, err := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	for _, part := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		if part == "" {
			continue
		}
		next, e := unix.Openat(fd, part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		unix.Close(fd)
		if e != nil {
			return nil, e
		}
		fd = next
	}
	return os.NewFile(uintptr(fd), path), nil
}
func privateInfo(info os.FileInfo) bool {
	s, ok := info.Sys().(*syscall.Stat_t)
	return ok && s.Uid == uint32(os.Getuid()) && info.Mode().Perm()&0o077 == 0
}
func readAt(dir *os.File, name string, limit int64, private bool) ([]byte, error) {
	fd, err := unix.Openat(int(dir.Fd()), name, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_CLOEXEC|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), name)
	defer f.Close()
	before, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() || before.Size() > limit || (private && !privateInfo(before)) {
		return nil, invalid()
	}
	data, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil {
		return nil, e
	}
	after, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if int64(len(data)) > limit || before.Size() != after.Size() || !before.ModTime().Equal(after.ModTime()) {
		return nil, invalid()
	}
	var leaf unix.Stat_t
	if e = unix.Fstatat(int(dir.Fd()), name, &leaf, unix.AT_SYMLINK_NOFOLLOW); e != nil {
		return nil, e
	}
	old := after.Sys().(*syscall.Stat_t)
	if uint64(old.Ino) != uint64(leaf.Ino) || uint64(old.Dev) != uint64(leaf.Dev) || leaf.Mode&unix.S_IFMT != unix.S_IFREG {
		if private {
			return nil, fail(4, "state_changed", "Notification state changed during readback; inspect it again.")
		}
		return nil, invalid()
	}
	return data, nil
}
func readFile(path string, limit int64, private bool) ([]byte, error) {
	if !validPath(path) {
		return nil, invalid()
	}
	parent, e := openDirectory(filepath.Dir(path))
	if e != nil {
		return nil, e
	}
	defer parent.Close()
	b, e := readAt(parent, filepath.Base(path), limit, private)
	if e != nil {
		return nil, e
	}
	current, e := openDirectory(filepath.Dir(path))
	if e != nil {
		return nil, e
	}
	defer current.Close()
	a, _ := parent.Stat()
	c, _ := current.Stat()
	if a == nil || c == nil || !os.SameFile(a, c) {
		return nil, invalid()
	}
	return b, nil
}
func stateDirectory(path string, create bool) (*os.File, error) {
	d, e := openDirectory(path)
	if errors.Is(e, os.ErrNotExist) {
		parent, err := openDirectory(filepath.Dir(path))
		if err != nil {
			return nil, fail(2, "invalid_state_parent", "State root requires an existing physical parent directory.")
		}
		defer parent.Close()
		if !create {
			return nil, os.ErrNotExist
		}
		if err = unix.Mkdirat(int(parent.Fd()), filepath.Base(path), 0o700); err != nil && !errors.Is(err, syscall.EEXIST) {
			return nil, err
		}
		if err = parent.Sync(); err != nil {
			return nil, err
		}
		d, e = openDirectory(path)
	}
	if e != nil {
		return nil, e
	}
	i, e := d.Stat()
	if e != nil || !privateInfo(i) {
		d.Close()
		return nil, fail(2, "unsafe_state", "Notification state must be private and owned by the current user.")
	}
	return d, nil
}
func lockState(dir *os.File) (func(), error) {
	fd, e := unix.Openat(int(dir.Fd()), ".lock", unix.O_RDWR|unix.O_CREAT|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if e != nil {
		if errors.Is(e, os.ErrNotExist) {
			return nil, fail(4, "state_changed", "The state root or lock changed during reservation; inspect it again.")
		}
		return nil, e
	}
	f := os.NewFile(uintptr(fd), ".lock")
	i, e := f.Stat()
	if e != nil || !i.Mode().IsRegular() || !privateInfo(i) {
		f.Close()
		return nil, invalid()
	}
	if e = unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); e != nil {
		f.Close()
		return nil, fail(4, "busy", "A notification operation is already holding this state root.")
	}
	return func() { _ = unix.Flock(fd, unix.LOCK_UN); _ = f.Close() }, nil
}
func atomicState(dir *os.File, data []byte) error {
	// A unique exclusive temporary inode; the committed name is never opened for write.
	name := ".pending-" + strings.TrimPrefix(digestValue([]any{os.Getpid(), timeNonce()}), "sha256:")
	fd, e := unix.Openat(int(dir.Fd()), name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if e != nil {
		return e
	}
	f := os.NewFile(uintptr(fd), name)
	defer f.Close()
	defer unix.Unlinkat(int(dir.Fd()), name, 0)
	if _, e = f.Write(data); e != nil {
		return e
	}
	if e = f.Sync(); e != nil {
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	if e = unix.Renameat(int(dir.Fd()), name, int(dir.Fd()), "state.json"); e != nil {
		return e
	}
	return dir.Sync()
}

// bindRoot is a separate durable tombstone: deleting the database must not make
// an already used root look like a new empty destination for another attempt.
func bindRoot(dir *os.File, name, sha string) error {
	data := []byte(name + "\n" + sha + "\n")
	fd, e := unix.Openat(int(dir.Fd()), "binding", unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if errors.Is(e, syscall.EEXIST) {
		b, err := readAt(dir, "binding", 1024, true)
		if err != nil {
			return err
		}
		if string(b) != string(data) {
			return fail(1, "state_corrupt", "The state root binding does not match.")
		}
		return nil
	}
	if e != nil {
		return e
	}
	f := os.NewFile(uintptr(fd), "binding")
	defer f.Close()
	if _, e = f.Write(data); e != nil {
		return e
	}
	if e = f.Sync(); e != nil {
		return e
	}
	if e = f.Close(); e != nil {
		return e
	}
	return dir.Sync()
}
