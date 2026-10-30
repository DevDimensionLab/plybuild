package workflowhandoff

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ValidateDeliveryOwnerContinuity(d Dependencies, locator string) error {
	snapshot, err := d.Store.ReadByLocator(locator)
	if err != nil {
		return err
	}
	if snapshot.Start == nil || snapshot.Start.Outcome != "started" || snapshot.Terminal != nil {
		return fmt.Errorf("the original run is not active")
	}
	return nil
}

func deliveryReadFile(files FileSystem, path string, limit int64) ([]byte, error) {
	physical, e := physicalRegularFile(files, path)
	if e != nil {
		return nil, e
	}
	if physical != path {
		return nil, fmt.Errorf("evidence path must be physical")
	}
	st, e := files.Lstat(path)
	if e != nil {
		return nil, e
	}
	if st.Size() > limit {
		return nil, fmt.Errorf("evidence exceeds size limit")
	}
	b, e := files.ReadFile(path)
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("evidence exceeds size limit")
	}
	return b, e
}

func deliveryEnsureDir(files FileSystem, path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("delivery directory must be absolute and clean")
	}
	// Check the existing ancestor before MkdirAll so a pre-existing staging
	// symlink cannot create a directory outside the managed reply area.
	for ancestor := path; ; ancestor = filepath.Dir(ancestor) {
		if _, err := files.Lstat(ancestor); err == nil {
			physical, err := files.EvalSymlinks(ancestor)
			if err != nil {
				return err
			}
			if physical != ancestor {
				return fmt.Errorf("delivery directory ancestor must be physical")
			}
			break
		} else if !os.IsNotExist(err) {
			return err
		}
		if filepath.Dir(ancestor) == ancestor {
			return fmt.Errorf("delivery directory ancestor is missing")
		}
	}
	if e := files.MkdirAll(path, 0700); e != nil {
		return e
	}
	p, e := files.EvalSymlinks(path)
	if e != nil {
		return e
	}
	if p != path {
		return fmt.Errorf("delivery directory must be physical")
	}
	return nil
}

func deliveryWriteOnce(files FileSystem, path string, b []byte) error {
	if _, e := files.Lstat(path); e == nil {
		old, e := deliveryReadFile(files, path, int64(len(b))+1)
		if e != nil {
			return e
		}
		if !bytes.Equal(old, b) {
			return fmt.Errorf("immutable delivery artifact differs: %s", path)
		}
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	p, e := files.EvalSymlinks(filepath.Dir(path))
	if e != nil {
		return e
	}
	if p != filepath.Dir(path) {
		return fmt.Errorf("delivery artifact parent must be physical")
	}
	f, e := files.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fs.FileMode(0600))
	if e != nil {
		return e
	}
	n, e := f.Write(b)
	if e == nil && n != len(b) {
		e = fmt.Errorf("short delivery artifact write")
	}
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
