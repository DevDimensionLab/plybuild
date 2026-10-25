package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func wr01PhysicalDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func wr01ContentError(t *testing.T, err error, code, detail string) {
	t.Helper()
	var content *TaskContentError
	if !errors.As(err, &content) || content.Code != code || content.Detail != detail || err.Error() != code+": "+detail || errors.Unwrap(err) != nil {
		t.Fatalf("error = %#v, want %s: %s without a cause", err, code, detail)
	}
}

func TestWR01PhysicalPaths(t *testing.T) {
	dir := wr01PhysicalDir(t)
	file := filepath.Join(dir, "content")
	if err := os.WriteFile(file, []byte("bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{dir, file} {
		if err := contentPhysical(path); err != nil {
			t.Fatalf("physical %s: %v", path, err)
		}
	}
	for _, path := range []string{"relative", dir + "/./content", dir + "/sub/../content"} {
		wr01ContentError(t, contentPhysical(path), "task_content_invalid_input", "path must be absolute and clean")
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(file, link); err != nil {
		t.Fatal(err)
	}
	dirLink := filepath.Join(dir, "dir-link")
	if err := os.Symlink(dir, dirLink); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{link, filepath.Join(dirLink, "content")} {
		wr01ContentError(t, contentPhysical(path), "task_content_integrity_conflict", "path contains a symlink: "+path)
		data, err := contentReadBounded(path, 10, false)
		if data != nil {
			t.Fatalf("symlink data = %q", data)
		}
		wr01ContentError(t, err, "task_content_integrity_conflict", "path contains a symlink: "+path)
	}
	missing := filepath.Join(dir, "missing")
	err := contentPhysical(missing)
	var pathError *os.PathError
	if !errors.Is(err, os.ErrNotExist) || !errors.As(err, &pathError) || pathError.Op != "lstat" || pathError.Path != missing {
		t.Fatalf("missing path error = %#v", err)
	}
}

func TestWR01ReadRealFiles(t *testing.T) {
	for _, tc := range []struct {
		name         string
		data         []byte
		mode         os.FileMode
		limit        int
		immutable    bool
		code, detail string
	}{
		{name: "exact limit", data: []byte{0, 255, '\n', 'x'}, mode: 0600, limit: 4},
		{name: "limit plus one", data: []byte("12345"), mode: 0600, limit: 4, code: "task_content_invalid_input", detail: "file exceeds 4 bytes"},
		{name: "empty", data: []byte{}, mode: 0600, limit: 0},
		{name: "zero limit", data: []byte("x"), mode: 0600, limit: 0, code: "task_content_invalid_input", detail: "file exceeds 0 bytes"},
		{name: "negative limit unchanged", data: []byte{}, mode: 0600, limit: -1, code: "task_content_invalid_input", detail: "file exceeds -1 bytes"},
		{name: "immutable", data: []byte("abc"), mode: 0400, limit: 3, immutable: true},
		{name: "optional immutable mode", data: []byte("abc"), mode: 0400, limit: 3},
		{name: "immutable wrong mode", data: []byte("abc"), mode: 0600, limit: 3, immutable: true, code: "task_content_integrity_conflict"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := wr01PhysicalDir(t)
			path := filepath.Join(dir, "content")
			if err := os.WriteFile(path, tc.data, tc.mode); err != nil {
				t.Fatal(err)
			}
			before, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			data, err := contentReadBounded(path, tc.limit, tc.immutable)
			if tc.code != "" {
				detail := tc.detail
				if detail == "" {
					detail = "unexpected file type or immutable mode: " + path
				}
				wr01ContentError(t, err, tc.code, detail)
				if data != nil {
					t.Fatalf("failure data = %q", data)
				}
			} else if err != nil || data == nil || !bytes.Equal(data, tc.data) {
				t.Fatalf("read = %q, %v; want %q", data, err, tc.data)
			}
			after, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(before, after) || before.Mode() != after.Mode() || before.Size() != after.Size() || before.ModTime() != after.ModTime() || !bytes.Equal(raw, tc.data) || len(entries) != 1 || entries[0].Name() != "content" {
				t.Fatal("read changed file identity, bytes, mode, metadata or directory entries")
			}
		})
	}
	dir := wr01PhysicalDir(t)
	data, err := contentReadBounded(dir, 10, false)
	if data != nil {
		t.Fatalf("directory data = %q", data)
	}
	wr01ContentError(t, err, "task_content_integrity_conflict", "unexpected file type or immutable mode: "+dir)
}

// Keep the existing private caller signatures checked as part of this contract.
var _ func(string) error = contentPhysical
var _ func(string, int, bool) ([]byte, error) = contentReadBounded

// Used by the deterministic operation trace tests below.
func wr01Trace(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("operations = %v, want %v", got, want)
	}
}
