// Package taskexecute composes goal selection and the existing delivery transport.
// The launcher owns durable intent and control files, not provider permissions.
package taskexecute

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func readPreservedJSON(path string, value any) ([]byte, error) {
	if err := physicalPath(path); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > 16<<20 {
		return nil, fmt.Errorf("preserved execute input is not a bounded regular file: %s", path)
	}
	b, err := io.ReadAll(io.LimitReader(f, (16<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(b) > 16<<20 {
		return nil, fmt.Errorf("preserved execute input exceeds 16 MiB")
	}
	if _, err = canonicaljson.DecodeStrict(b); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(value); err != nil {
		return nil, err
	}
	return b, nil
}

func digestBytes(b []byte) string {
	s := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(s[:])
}

// physicalPath checks every existing path component. In particular an absent
// leaf beneath a symlink is not an acceptable place for a control executable.
func physicalPath(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("execute path must be absolute and clean: %s", path)
	}
	for p := path; ; p = filepath.Dir(p) {
		i, err := os.Lstat(p)
		if err == nil {
			if i.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("execute path contains a symlink: %s", p)
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if p == filepath.Dir(p) {
			return nil
		}
	}
}

func privateDir(path string) error {
	if err := physicalPath(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0700)
}

// writeOnce publishes without replacing a previous intent. A lost successful
// reply may reuse identical bytes; conflicting bytes always require inspection.
func writeOnce(path string, data []byte, mode os.FileMode) error {
	if err := physicalPath(path); err != nil {
		return err
	}
	if old, err := os.ReadFile(path); err == nil {
		if bytes.Equal(old, data) {
			return nil
		}
		return fmt.Errorf("preserved execute file differs: %s", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := privateDir(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".execute-publish-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err = os.Link(f.Name(), path); os.IsExist(err) {
		old, readErr := os.ReadFile(path)
		if readErr == nil && bytes.Equal(old, data) {
			return nil
		}
		return fmt.Errorf("execute publication conflicted: %s", path)
	}
	if err != nil {
		return err
	}
	d, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

func bindExecutable(path string) (taskrun.Executable, error) {
	var out taskrun.Executable
	absolute, err := filepath.Abs(path)
	if err != nil {
		return out, err
	}
	path, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return out, err
	}
	f, err := os.Open(path)
	if err != nil {
		return out, err
	}
	defer f.Close()
	i, err := f.Stat()
	if err != nil {
		return out, err
	}
	if !i.Mode().IsRegular() || i.Mode().Perm()&0111 == 0 {
		return out, fmt.Errorf("not an executable regular file: %s", path)
	}
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return out, err
	}
	return taskrun.Executable{Path: path, SHA256: "sha256:" + hex.EncodeToString(h.Sum(nil))}, nil
}

// PreserveControl copies the running binary into the execution's private area.
// A hard link to the installed binary would not protect against in-place writes.
func PreserveControl(source taskrun.Executable, directory string) (taskrun.Executable, error) {
	out := taskrun.Executable{Path: filepath.Join(directory, "ply-control"), SHA256: source.SHA256}
	b, err := os.ReadFile(source.Path)
	if err != nil {
		return out, err
	}
	if digestBytes(b) != source.SHA256 {
		return out, fmt.Errorf("control source changed before preservation")
	}
	if err = writeOnce(out.Path, b, 0500); err != nil {
		return out, err
	}
	actual, err := bindExecutable(out.Path)
	if err != nil || actual != out {
		return out, fmt.Errorf("preserved control executable is not intact: %v", err)
	}
	return out, nil
}
