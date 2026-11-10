//go:build !windows

package taskjournal

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Walk physical components through pinned descriptors. O_NOFOLLOW applies to
// every component, including ancestors, so a path race cannot redirect writes.
func openDir(path string, create bool) (*os.File, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, invalid("path must be absolute and clean")
	}
	fd, e := unix.Open("/", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if e != nil {
		return nil, e
	}
	for _, part := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		if part == "" {
			continue
		}
		next, err := unix.Openat(fd, part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err == unix.ENOENT && create {
			if err = unix.Mkdirat(fd, part, 0700); err != nil && err != unix.EEXIST {
				unix.Close(fd)
				return nil, err
			}
			next, err = unix.Openat(fd, part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		}
		unix.Close(fd)
		if err != nil {
			return nil, &os.PathError{Op: "open physical directory", Path: path, Err: err}
		}
		fd = next
	}
	return os.NewFile(uintptr(fd), path), nil
}

// Only journal-owned components can be created. Every parent is synced before
// acknowledging a record; a missing workspace ancestor is never recreated.
func createStore(root, task string) (*os.File, error) {
	d, e := openDir(filepath.Join(root, ".ply"), false)
	if e != nil {
		return nil, e
	}
	for _, part := range []string{"task-process", "v1", task} {
		fd, err := unix.Openat(int(d.Fd()), part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err == unix.ENOENT {
			err = unix.Mkdirat(int(d.Fd()), part, 0700)
			if err == nil || err == unix.EEXIST {
				err = d.Sync()
			}
			if err == nil {
				fd, err = unix.Openat(int(d.Fd()), part, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
			}
		}
		if err != nil {
			d.Close()
			return nil, &os.PathError{Op: "open journal directory", Path: part, Err: err}
		}
		next := os.NewFile(uintptr(fd), filepath.Join(d.Name(), part))
		d.Close()
		d = next
		st, err := d.Stat()
		if err != nil || st.Mode().Perm() != 0700 {
			d.Close()
			return nil, conflict("journal directories must be private")
		}
	}
	return d, nil
}
func openAt(dir *os.File, name string, flags int, mode uint32) (*os.File, error) {
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name {
		return nil, invalid("invalid local filename")
	}
	fd, e := unix.Openat(int(dir.Fd()), name, flags|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, mode)
	if e != nil {
		return nil, &os.PathError{Op: "open", Path: filepath.Join(dir.Name(), name), Err: e}
	}
	f := os.NewFile(uintptr(fd), name)
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() {
		f.Close()
		return nil, invalid("file is not regular")
	}
	return f, nil
}
func readOpened(f *os.File, limit int) ([]byte, error) {
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	b, e := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if e != nil {
		return nil, e
	}
	after, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if len(b) > limit {
		return nil, invalid("file exceeds size limit")
	}
	if st.Size() != int64(len(b)) || st.Size() != after.Size() || st.ModTime() != after.ModTime() {
		return nil, conflict("file changed during read")
	}
	return b, nil
}
func readFile(path string, limit int) ([]byte, error) {
	d, e := openDir(filepath.Dir(path), false)
	if e != nil {
		return nil, e
	}
	defer d.Close()
	f, e := openAt(d, filepath.Base(path), unix.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return readOpened(f, limit)
}
func sourceHash(path string) (string, error) {
	d, e := openDir(filepath.Dir(path), false)
	if e != nil {
		return "", e
	}
	defer d.Close()
	f, e := openAt(d, filepath.Base(path), unix.O_RDONLY, 0)
	if e != nil {
		return "", e
	}
	defer f.Close()
	return hashStream(f)
}
func lock(dir *os.File) (*os.File, error) {
	// Separate creation from opening an existing lock. This also avoids a
	// competing O_CREAT lookup racing the first publication on Darwin.
	f, e := openAt(dir, "append.lock", unix.O_RDWR|unix.O_CREAT|unix.O_EXCL, 0600)
	if os.IsExist(e) {
		f, e = openAt(dir, "append.lock", unix.O_RDWR, 0)
	}
	if e != nil {
		return nil, e
	}
	st, e := f.Stat()
	if e != nil || st.Mode().Perm() != 0600 {
		f.Close()
		return nil, conflict("append lock must be private")
	}
	if e = unix.Flock(int(f.Fd()), unix.LOCK_EX); e != nil {
		f.Close()
		return nil, e
	}
	return f, nil
}
func unlock(f *os.File) error {
	e := unix.Flock(int(f.Fd()), unix.LOCK_UN)
	c := f.Close()
	if e != nil {
		return e
	}
	return c
}
func linkAt(dir *os.File, from, to string) error {
	return unix.Linkat(int(dir.Fd()), from, int(dir.Fd()), to, 0)
}
func removeAt(dir *os.File, name string) error { return unix.Unlinkat(int(dir.Fd()), name, 0) }
func writeTemp(dir *os.File, name string, b []byte) error {
	f, e := openAt(dir, name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL, 0600)
	if e != nil {
		return e
	}
	n, e := f.Write(b)
	if e == nil && n != len(b) {
		e = io.ErrShortWrite
	}
	if e == nil {
		e = f.Sync()
	}
	c := f.Close()
	if e != nil {
		return e
	}
	return c
}
func regularMode(path string) error {
	d, e := openDir(path, false)
	if e != nil {
		return e
	}
	defer d.Close()
	st, e := d.Stat()
	if e != nil {
		return e
	}
	if st.Mode().Perm() != 0700 {
		return fmt.Errorf("journal directory must be mode 0700")
	}
	return nil
}
