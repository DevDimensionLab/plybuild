package privatestore

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func physicalTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestPrivateWriteOncePublishesOnceAndRejectsDifferentBytes(t *testing.T) {
	path := filepath.Join(physicalTempDir(t), "store", "fact.json")
	if err := PrivateWriteOnce(path, []byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	if err := PrivateWriteOnce(path, []byte(`{"a":1}`)); err != nil {
		t.Fatalf("identical bytes: %v", err)
	}
	var typed *Error
	if err := PrivateWriteOnce(path, []byte(`{"a":2}`)); !errors.As(err, &typed) || typed.Code != "task_run_conflict" || typed.Exit != 3 {
		t.Fatalf("different bytes: %v", err)
	}
	if raw, err := PrivateRead(path, 1<<10); err != nil || string(raw) != `{"a":1}` {
		t.Fatalf("read = %q %v", raw, err)
	}
	if info, err := os.Stat(filepath.Dir(path)); err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("store directory not private: %v %v", info.Mode(), err)
	}
}

func TestPrivateReadRefusesSharedModeAndSymlinks(t *testing.T) {
	dir := physicalTempDir(t)
	shared := filepath.Join(dir, "shared.json")
	if err := os.WriteFile(shared, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := PrivateRead(shared, 1<<10); err == nil || !strings.Contains(err.Error(), "owner-only") || !strings.Contains(err.Error(), "mode 0644") {
		t.Fatalf("shared file read: %v", err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}
	if err := PhysicalPath(filepath.Join(link, "shared.json")); err == nil {
		t.Fatal("symlinked path accepted")
	}
	if err := PhysicalPath("relative"); err == nil {
		t.Fatal("relative path accepted")
	}
}

func TestStrictDecodeRequiresExactMembers(t *testing.T) {
	type doc struct {
		A string  `json:"a"`
		B *string `json:"b"`
	}
	var out doc
	if err := StrictDecode([]byte(`{"a":"x","b":null}`), 64, &out); err != nil || out.A != "x" {
		t.Fatalf("exact = %+v %v", out, err)
	}
	for _, raw := range []string{`{"a":"x"}`, `{"a":"x","b":null,"c":1}`, `{"a":"x","a":"y","b":null}`, `{"a":"x","b":null} {}`} {
		if err := StrictDecode([]byte(raw), 64, &doc{}); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
	if err := StrictDecode([]byte(`{"a":"x","b":null}`), 4, &doc{}); err == nil {
		t.Error("accepted oversized document")
	}
}

func TestWithPrivateLockHoldsStoreLock(t *testing.T) {
	dir := filepath.Join(physicalTempDir(t), "locked")
	called := false
	if err := WithPrivateLock(dir, func() error { called = true; return nil }); err != nil || !called {
		t.Fatalf("lock: %v", err)
	}
	if info, err := os.Stat(filepath.Join(dir, "store.lock")); err != nil || info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("lock file: %v", err)
	}
}

// The injected boundary controls read order only; writes use the real file.
func TestReadRefusesFileChangedWhileReading(t *testing.T) {
	for _, change := range []string{"replace inode", "grow beyond limit"} {
		t.Run(change, func(t *testing.T) {
			path := filepath.Join(physicalTempDir(t), "input")
			if err := os.WriteFile(path, []byte("12345678"), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := readFileAsWithRead(path, "input file", 8, noSharedWriteAccess, func(r io.Reader) ([]byte, error) {
				if change == "replace inode" {
					next := path + ".next"
					if err := os.WriteFile(next, []byte("abcdefgh"), 0o644); err != nil {
						return nil, err
					}
					if err := os.Rename(next, path); err != nil {
						return nil, err
					}
				} else if err := os.WriteFile(path, []byte("123456789"), 0o644); err != nil {
					return nil, err
				}
				return io.ReadAll(r)
			})
			if err == nil || !strings.Contains(err.Error(), "file changed while reading") || !strings.Contains(err.Error(), path) {
				t.Fatalf("instability: %v", err)
			}
		})
	}
}

func TestCanonicalSortsMembers(t *testing.T) {
	raw, err := Canonical(map[string]int{"b": 1, "a": 2})
	if err != nil || string(raw) != `{"a":2,"b":1}` {
		t.Fatalf("canonical = %s %v", raw, err)
	}
}
