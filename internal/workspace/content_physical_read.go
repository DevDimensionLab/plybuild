package workspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func contentPhysical(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return contentError("task_content_invalid_input", "path must be absolute and clean", nil)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	if resolved != path {
		return contentError("task_content_integrity_conflict", "path contains a symlink: "+path, nil)
	}
	return nil
}

// contentReadFile and contentReadOps expose only the operations needed by one
// bounded read. Production and deterministic tests share the same read loop.
type contentReadFile interface {
	io.Reader
	Stat() (os.FileInfo, error)
	Close() error
}

type contentReadOps struct {
	lstat func(string) (os.FileInfo, error)
	open  func(string) (contentReadFile, error)
}

func contentReadBounded(path string, limit int, immutable bool) ([]byte, error) {
	return contentReadBoundedWith(path, limit, immutable, contentReadOps{
		lstat: os.Lstat,
		open:  func(path string) (contentReadFile, error) { return os.Open(path) },
	})
}

func contentReadBoundedWith(path string, limit int, immutable bool, ops contentReadOps) (data []byte, err error) {
	if err = contentPhysical(path); err != nil {
		return nil, err
	}
	before, err := ops.lstat(path)
	if err != nil {
		return nil, err
	}
	if !before.Mode().IsRegular() || immutable && before.Mode().Perm() != 0o400 {
		return nil, contentError("task_content_integrity_conflict", "unexpected file type or immutable mode: "+path, nil)
	}
	f, err := ops.open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if e := f.Close(); err == nil {
			err = e
		}
	}()
	opened, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(before, opened) {
		return nil, contentError("task_content_integrity_conflict", "file changed before read", nil)
	}
	data, err = io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > limit {
		return nil, contentError("task_content_invalid_input", fmt.Sprintf("file exceeds %d bytes", limit), nil)
	}
	after, e := f.Stat()
	current, p := ops.lstat(path)
	if e != nil || p != nil {
		return nil, firstError(e, p)
	}
	if !current.Mode().IsRegular() || !os.SameFile(opened, current) || opened.Size() != int64(len(data)) || after.Size() != opened.Size() || after.ModTime() != opened.ModTime() {
		return nil, contentError("task_content_integrity_conflict", "file changed or was short during read", nil)
	}
	return data, nil
}
