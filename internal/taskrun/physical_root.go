package taskrun

import (
	"os"
	"path/filepath"
	"strings"
)

// Open each component relative to an already open directory. Checking a path
// and subsequently opening its full spelling would allow an ancestor to be
// replaced by a symlink between those operations.
func physicalRoot(path string, create bool) (*os.Root, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, invalid("directory must be absolute and clean")
	}
	prefix := filepath.VolumeName(path) + string(filepath.Separator)
	root, e := os.OpenRoot(prefix)
	if e != nil {
		return nil, e
	}
	for _, part := range strings.Split(strings.TrimPrefix(path, prefix), string(filepath.Separator)) {
		if part == "" {
			continue
		}
		before, err := root.Lstat(part)
		if os.IsNotExist(err) && create {
			err = root.Mkdir(part, 0700)
			if err == nil {
				dir, de := root.Open(".")
				if de == nil {
					de = dir.Sync()
					dir.Close()
				}
				if de != nil {
					root.Close()
					return nil, de
				}
			}
			if err == nil || os.IsExist(err) {
				before, err = root.Lstat(part)
			}
		}
		if err != nil {
			root.Close()
			return nil, err
		}
		if !before.IsDir() || before.Mode()&os.ModeSymlink != 0 {
			root.Close()
			return nil, conflict("directory component is not a physical directory")
		}
		next, err := root.OpenRoot(part)
		if err != nil {
			root.Close()
			return nil, err
		}
		actual, err := next.Stat(".")
		after, ae := root.Lstat(part)
		if err != nil || ae != nil || !os.SameFile(before, actual) || !os.SameFile(actual, after) || after.Mode()&os.ModeSymlink != 0 {
			next.Close()
			root.Close()
			return nil, conflict("directory changed while opening")
		}
		root.Close()
		root = next
	}
	return root, nil
}
func physicalOpenFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	root, e := physicalRoot(filepath.Dir(path), false)
	if e != nil {
		return nil, e
	}
	defer root.Close()
	name := filepath.Base(path)
	before, be := root.Lstat(name)
	if be != nil && !(os.IsNotExist(be) && flags&os.O_CREATE != 0) {
		return nil, be
	}
	if be == nil && (!before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0) {
		return nil, conflict("file is not a physical regular file")
	}
	f, e := root.OpenFile(name, flags, mode)
	if e != nil {
		return nil, e
	}
	actual, e := f.Stat()
	after, ae := root.Lstat(name)
	if e != nil || ae != nil || !actual.Mode().IsRegular() || after.Mode()&os.ModeSymlink != 0 || !os.SameFile(actual, after) || be == nil && !os.SameFile(before, actual) {
		f.Close()
		return nil, conflict("file changed while opening")
	}
	return f, nil
}
