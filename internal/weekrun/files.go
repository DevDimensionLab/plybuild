package weekrun

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func hash(b []byte) string { v := sha256.Sum256(b); return "sha256:" + hex.EncodeToString(v[:]) }
func canonical(v any) []byte {
	b, err := taskrun.Canonical(v)
	if err != nil {
		panic(err)
	}
	return b
}
func physical(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("week path must be absolute and clean: %s", path)
	}
	for p := path; ; p = filepath.Dir(p) {
		i, e := os.Lstat(p)
		if e == nil && i.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("week path contains a symlink: %s", p)
		}
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if p == filepath.Dir(p) {
			return nil
		}
	}
}
func directory(path string) error {
	if err := physical(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0700)
}
func readJSON(path string, v any) error {
	if err := physical(path); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Size() > 4<<20 {
		return fmt.Errorf("week record is not a bounded regular file: %s", path)
	}
	b, err := io.ReadAll(io.LimitReader(f, 4<<20+1))
	if err != nil {
		return err
	}
	if len(b) > 4<<20 {
		return fmt.Errorf("week record too large")
	}
	if _, err = canonicaljson.DecodeStrict(b); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func publish(path string, v any, replace bool) error {
	return publishBytes(path, canonical(v), 0600, replace)
}
func publishBytes(path string, b []byte, mode os.FileMode, replace bool) error {
	if err := physical(path); err != nil {
		return err
	}
	if err := directory(filepath.Dir(path)); err != nil {
		return err
	}
	if !replace {
		if old, e := os.ReadFile(path); e == nil {
			if bytes.Equal(old, b) {
				return nil
			}
			return fmt.Errorf("immutable week record differs: %s", path)
		} else if !os.IsNotExist(e) {
			return e
		}
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".week-publish-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = f.Sync()
	}
	ce := f.Close()
	if err == nil {
		err = ce
	}
	if err != nil {
		return err
	}
	if replace {
		err = os.Rename(f.Name(), path)
	} else {
		err = os.Link(f.Name(), path)
		if os.IsExist(err) {
			old, e := os.ReadFile(path)
			if e == nil && bytes.Equal(old, b) {
				return nil
			}
		}
	}
	if err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}
func runDir(root, id string) (string, error) {
	if !slug.MatchString(id) {
		return "", fmt.Errorf("invalid week ID")
	}
	p := filepath.Join(root, ".ply", "week-runs", id)
	return p, physical(p)
}
func contractID(c Contract) string {
	return "week_" + strings.TrimPrefix(hash(canonical([]string{c.Workspace, c.Name})), "sha256:")[:32]
}
func readBinding(root, id string) (Binding, string, error) {
	dir, err := runDir(root, id)
	if err != nil {
		return Binding{}, "", err
	}
	var b Binding
	if err = readJSON(filepath.Join(dir, "binding.json"), &b); err != nil {
		return b, dir, err
	}
	if err = Validate(b.Contract); err != nil {
		return b, dir, err
	}
	if b.ID != id || b.Contract.Workspace != root || contractID(b.Contract) != id || b.SHA256 != hash(canonical(b.Contract)) {
		return b, dir, fmt.Errorf("week immutable binding differs")
	}
	if b.Contract.RuntimeControl != nil {
		if !bytes.Equal(canonical(b.ControlContract), canonical(taskrun.NewControlContract("week", id, b.SHA256))) {
			return b, dir, fmt.Errorf("week control contract differs")
		}
	} else if b.ControlContract != nil {
		return b, dir, fmt.Errorf("old week cannot acquire bound control")
	}
	return b, dir, nil
}
func readState(b Binding, dir string) (State, error) {
	var s State
	if err := readJSON(filepath.Join(dir, "state.json"), &s); err != nil {
		return s, err
	}
	if s.Kind != "ply.workflow.week" || s.SchemaVersion != 1 || s.ID != b.ID || s.ContractSHA256 != b.SHA256 || s.StartedAt != b.StartedAt || s.Deadline != b.Deadline || len(s.Entries) != len(b.Contract.Selected) {
		return s, fmt.Errorf("week observation does not match immutable binding")
	}
	for i, e := range s.Entries {
		if e.TaskID != b.Contract.Selected[i].TaskID {
			return s, fmt.Errorf("week selected Task order differs")
		}
	}
	return s, nil
}
func executable(path string) (taskrun.Executable, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return taskrun.Executable{}, err
	}
	p, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return taskrun.Executable{}, err
	}
	f, err := os.Open(p)
	if err != nil {
		return taskrun.Executable{}, err
	}
	defer f.Close()
	i, err := f.Stat()
	if err != nil {
		return taskrun.Executable{}, err
	}
	if !i.Mode().IsRegular() || i.Mode().Perm()&0111 == 0 {
		return taskrun.Executable{}, fmt.Errorf("week control is not an executable regular file")
	}
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return taskrun.Executable{}, err
	}
	return taskrun.Executable{Path: p, SHA256: "sha256:" + hex.EncodeToString(h.Sum(nil))}, nil
}
