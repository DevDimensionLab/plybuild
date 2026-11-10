package atomicfile

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReplaceCompleteBytesAndMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	for _, contents := range []string{"old", "complete replacement", ""} {
		if err := Replace(path, []byte(contents), 0o640); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != contents {
			t.Fatalf("bytes = %q, %v", got, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0o640 {
			t.Fatalf("mode = %v", info.Mode())
		}
		entries, err := os.ReadDir(filepath.Dir(path))
		if err != nil || len(entries) != 1 {
			t.Fatalf("entries = %v, %v", entries, err)
		}
	}
}

func TestReplaceFailureOutcomes(t *testing.T) {
	for _, stage := range []string{"temp", "write", "short write", "chmod", "sync", "close", "replace", "directory sync", "cleanup", "abort close"} {
		t.Run(stage, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "state")
			if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
				t.Fatal(err)
			}
			ops := systemOperations()
			injected := errors.New("injected " + stage)
			secondary := errors.New("secondary failure")
			var handle *os.File
			create := ops.createTemp
			ops.createTemp = func(dir, pattern string) (*os.File, error) {
				if stage == "temp" {
					return nil, injected
				}
				f, err := create(dir, pattern)
				handle = f
				return f, err
			}
			switch stage {
			case "write", "cleanup", "abort close":
				ops.write = func(*os.File, []byte) (int, error) { return 1, injected }
			case "short write":
				ops.write = func(*os.File, []byte) (int, error) { return 1, nil }
			case "chmod":
				ops.chmod = func(*os.File, fs.FileMode) error { return injected }
			case "sync":
				ops.sync = func(*os.File) error { return injected }
			case "close":
				ops.close = func(f *os.File) error { _ = f.Close(); return injected }
			case "replace":
				ops.replace = func(string, string) error { return injected }
			case "directory sync":
				ops.syncDirectory = func(string) error { return injected }
			}
			if stage == "cleanup" {
				ops.remove = func(string) error { return secondary }
			}
			if stage == "abort close" {
				ops.close = func(f *os.File) error { _ = f.Close(); return secondary }
			}
			err := replace(path, []byte("complete new"), 0o640, ops)
			var failure *Error
			replaced := stage == "directory sync"
			if !errors.As(err, &failure) || failure.Replaced != replaced {
				t.Fatalf("outcome = %#v, %v", failure, err)
			}
			cause := injected
			if stage == "short write" {
				cause = io.ErrShortWrite
			}
			if !errors.Is(err, cause) {
				t.Fatalf("cause lost: %v", err)
			}
			if (stage == "cleanup" || stage == "abort close") && !strings.Contains(err.Error(), secondary.Error()) {
				t.Fatalf("secondary failure lost: %v", err)
			}
			wantOp := map[string]string{"temp": "open temp", "short write": "write", "cleanup": "write", "abort close": "write", "directory sync": "sync directory"}[stage]
			if wantOp == "" {
				wantOp = stage
			}
			if failure.Op != wantOp || failure.Path == "" {
				t.Fatalf("context = %#v", failure)
			}
			got, readErr := os.ReadFile(path)
			want := "old"
			if replaced {
				want = "complete new"
			}
			if readErr != nil || string(got) != want {
				t.Fatalf("destination = %q, %v", got, readErr)
			}
			entries, err := os.ReadDir(dir)
			wantEntries := 1
			if stage == "cleanup" {
				wantEntries = 2
			}
			if err != nil || len(entries) != wantEntries {
				t.Fatalf("cleanup: %v, %v", entries, err)
			}
			if handle != nil {
				if _, err := handle.Stat(); !errors.Is(err, os.ErrClosed) {
					t.Fatalf("handle not closed: %v", err)
				}
			}
		})
	}
}

func TestReplaceMissingDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent", "state")
	var failure *Error
	err := Replace(path, []byte("new"), 0o600)
	if !errors.As(err, &failure) || failure.Replaced || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("error = %v", err)
	}
}

func TestReplaceErrorDoesNotProveAbsenceOfEffect(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	ops := systemOperations()
	original := ops.replace
	injected := errors.New("unconfirmed replace outcome")
	ops.replace = func(source, destination string) error {
		if err := original(source, destination); err != nil {
			return err
		}
		return injected
	}
	err := replace(path, []byte("published"), 0o600, ops)
	var failure *Error
	if !errors.As(err, &failure) || failure.Replaced || !errors.Is(err, injected) {
		t.Fatalf("unobserved success = %#v, %v", failure, err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != "published" {
		t.Fatalf("effect = %q, %v", got, err)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil || len(entries) != 1 {
		t.Fatalf("cleanup = %v, %v", entries, err)
	}
}
