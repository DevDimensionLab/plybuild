// Package privatestore holds the private-file primitives shared by sibling
// stores, such as human input and agent work, so every store keeps one
// publication and symlink policy.
package privatestore

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func Canonical(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	x, e := canonicaljson.DecodeStrict(b)
	if e != nil {
		return nil, e
	}
	return canonicaljson.Marshal(x)
}

func decode(b []byte, max int, out any) error {
	if len(b) > max {
		return invalid("document exceeds size limit")
	}
	if _, e := canonicaljson.DecodeStrict(b); e != nil {
		return invalid(e.Error())
	}
	var shape any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if e := dec.Decode(&shape); e != nil {
		return invalid(e.Error())
	}
	if e := exactShape(shape, reflect.TypeOf(out).Elem()); e != nil {
		return invalid(e.Error())
	}
	dec = json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if e := dec.Decode(out); e != nil {
		return invalid(e.Error())
	}
	return nil
}

// exactShape requires all members (including explicit nullable members), at every
// typed level. Raw WF objects are checked by the authoritative WF validators.
func exactShape(v any, t reflect.Type) error {
	if t == reflect.TypeOf(json.RawMessage{}) {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		if v == nil {
			return nil
		}
		return exactShape(v, t.Elem())
	}
	if v == nil {
		return fmt.Errorf("unexpected null for %s", t)
	}
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		fields := map[string]reflect.Type{}
		optional := map[string]bool{}
		var add func(reflect.Type)
		add = func(t reflect.Type) {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Anonymous {
					add(f.Type)
					continue
				}
				name := strings.Split(f.Tag.Get("json"), ",")[0]
				if name != "-" {
					fields[name] = f.Type
					optional[name] = strings.Contains(f.Tag.Get("json"), ",omitempty")
				}
			}
		}
		add(t)
		if len(m) > len(fields) {
			return fmt.Errorf("missing or unknown fields for %s", t)
		}
		for n, ft := range fields {
			x, ok := m[n]
			if !ok {
				if optional[n] {
					continue
				}
				return fmt.Errorf("missing %s", n)
			}
			if e := exactShape(x, ft); e != nil {
				return fmt.Errorf("%s: %w", n, e)
			}
		}
	case reflect.Slice:
		a, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected array")
		}
		for _, x := range a {
			if e := exactShape(x, t.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}

func physical(path string, missing bool) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || strings.ContainsRune(path, 0) {
		return invalid("path must be absolute and clean")
	}
	for p := path; ; p = filepath.Dir(p) {
		info, e := os.Lstat(p)
		if e != nil {
			if missing && os.IsNotExist(e) {
			} else {
				return e
			}
		} else if info.Mode()&os.ModeSymlink != 0 {
			return conflict("symlink path rejected")
		}
		if p == filepath.Dir(p) {
			break
		}
	}
	return nil
}

// fileAccess is the mode requirement of a file role.
type fileAccess int

const (
	anyAccess fileAccess = iota
	// Native stores, controls and authority-bearing inputs.
	ownerOnlyAccess
	// Ordinary hash-bound inputs: group/other read is allowed, write is not.
	// Mode bits say nothing about who owns or wrote the file.
	noSharedWriteAccess
)

func readFile(path string, max int, private bool) ([]byte, error) {
	if private {
		return readFileAs(path, "private file", max, ownerOnlyAccess)
	}
	return readFileAs(path, "input file", max, anyAccess)
}

// readFileAs names the role, path, failing requirement, observed metadata and
// supported handling in every retained refusal, never the file's contents.
func readFileAs(path, role string, max int, access fileAccess) ([]byte, error) {
	return readFileAsWithRead(path, role, max, access, io.ReadAll)
}

// readFileAsWithRead keeps the bounded read local so tests can deterministically
// interleave real writes. Identity checks do not establish an atomic snapshot:
// same-inode writes may yield mixed bytes, which callers bind by their digest.
func readFileAsWithRead(path, role string, max int, access fileAccess, read func(io.Reader) ([]byte, error)) ([]byte, error) {
	if e := physical(path, false); e != nil {
		if t, ok := e.(*Error); ok {
			return nil, failure(t.Code, t.Exit, fmt.Sprintf("%s %s: %s; use the absolute physical path of a regular file, without symlinks", role, path, t.Detail))
		}
		return nil, e
	}
	before, e := os.Lstat(path)
	if e != nil {
		return nil, e
	}
	refuse := func(requirement, handling string) error {
		return invalid(fmt.Sprintf("%s %s: %s; observed %s, mode %04o, %d bytes; %s", role, path, requirement, fileKind(before.Mode()), before.Mode().Perm(), before.Size(), handling))
	}
	perm := before.Mode().Perm()
	switch {
	case !before.Mode().IsRegular():
		return nil, refuse("requires a regular file", "submit a regular file at this path")
	case before.Size() > int64(max):
		return nil, refuse(fmt.Sprintf("requires at most %d bytes", max), "keep this input within the limit and move bulk output to a separate bounded file")
	case access == ownerOnlyAccess && perm&0077 != 0:
		return nil, refuse("requires owner-only permissions because it holds native or authority-bearing data", "check that no other user changed it, then restore owner-only mode (chmod go-rwx)")
	case access == noSharedWriteAccess && perm&0022 != 0:
		return nil, refuse("must not be writable by group or other users, who could change it at any time before Ply reads it", "remove that write access (chmod go-w) and submit the same path again; read access for others is accepted")
	}
	f, e := openRead(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(before, actual) {
		return nil, conflict(role + " " + path + ": file changed while opening; submit it again once its writer has finished")
	}
	b, e := read(io.LimitReader(f, int64(max)+1))
	if e != nil {
		return nil, e
	}
	after, e := os.Lstat(path)
	if e != nil || !os.SameFile(actual, after) || len(b) > max {
		return nil, conflict(role + " " + path + ": file changed while reading; submit it again once its writer has finished")
	}
	return b, nil
}

func fileKind(m os.FileMode) string {
	switch {
	case m.IsRegular():
		return "regular file"
	case m.IsDir():
		return "directory"
	case m&os.ModeSymlink != 0:
		return "symlink"
	case m&os.ModeNamedPipe != 0:
		return "named pipe"
	case m&os.ModeSocket != 0:
		return "socket"
	case m&os.ModeDevice != 0:
		return "device"
	}
	return "irregular file"
}

func privateDir(path string) error {
	root, e := physicalRoot(path, true)
	if e != nil {
		return e
	}
	defer root.Close()
	info, e := root.Stat(".")
	if e != nil {
		return e
	}
	if !info.IsDir() || info.Mode().Perm()&0077 != 0 {
		return integrity("store directory must be private")
	}
	return nil
}

// writeOnce publishes only complete, synced bytes. Link is atomic and cannot
// clobber a concurrent writer; the temporary is owned exclusively by this call.
func writeOnce(path string, b []byte) error { return publishFile(path, b, false) }

func publishFile(path string, b []byte, replace bool) error {
	if e := physical(path, true); e != nil {
		return e
	}
	if !replace {
		// Idempotent publication must also work for bounded caller inputs such
		// as private executable snapshots, which can exceed the JSON limit.
		limit := 8 << 20
		if len(b) > limit {
			limit = len(b)
		}
		if old, e := readFile(path, limit, true); e == nil {
			if bytes.Equal(old, b) {
				return nil
			}
			return conflict("immutable bytes differ: " + filepath.Base(path))
		} else if !os.IsNotExist(e) {
			return e
		}
	}
	parent := filepath.Dir(path)
	if e := privateDir(parent); e != nil {
		return e
	}
	before, e := os.Lstat(parent)
	if e != nil {
		return e
	}
	root, e := physicalRoot(parent, false)
	if e != nil {
		return e
	}
	defer root.Close()
	actual, e := root.Stat(".")
	if e != nil || !os.SameFile(before, actual) {
		return conflict("publication directory changed")
	}
	nonce := make([]byte, 16)
	if _, e = rand.Read(nonce); e != nil {
		return e
	}
	name := ".publish-" + hex.EncodeToString(nonce)
	f, e := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	defer root.Remove(name)
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	c := f.Close()
	if e != nil {
		return e
	}
	if c != nil {
		return c
	}
	destination := filepath.Base(path)
	if replace {
		e = root.Rename(name, destination)
	} else {
		e = root.Link(name, destination)
	}
	if e != nil {
		if !replace && os.IsExist(e) {
			old, err := root.ReadFile(destination)
			if err == nil && bytes.Equal(old, b) {
				return nil
			}
		}
		return e
	}
	dir, e := root.Open(".")
	if e != nil {
		return e
	}
	defer dir.Close()
	return dir.Sync()
}

type Error struct {
	Code, Detail string
	Exit         int
}

func (e *Error) Error() string { return e.Code + ": " + e.Detail }

func failure(code string, exit int, detail string) error { return &Error{code, detail, exit} }

func invalid(detail string) error { return failure("task_run_schema_invalid", 2, detail) }

func conflict(detail string) error { return failure("task_run_conflict", 3, detail) }

func integrity(detail string) error { return failure("task_run_integrity", 3, detail) }
