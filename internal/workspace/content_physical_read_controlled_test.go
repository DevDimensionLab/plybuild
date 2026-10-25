package workspace

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type wr01ReadFile struct {
	trace                 *[]string
	data                  []byte
	readErr, closeErr     error
	stats                 []os.FileInfo
	statErrors            []error
	statCalls, closeCalls int
}

func (f *wr01ReadFile) Read(p []byte) (int, error) {
	*f.trace = append(*f.trace, "read")
	n := copy(p, f.data)
	f.data = f.data[n:]
	if f.readErr != nil {
		return n, f.readErr
	}
	if len(f.data) == 0 {
		return n, io.EOF
	}
	return n, nil
}
func (f *wr01ReadFile) Stat() (os.FileInfo, error) {
	i := f.statCalls
	f.statCalls++
	*f.trace = append(*f.trace, []string{"stat-open", "stat-after"}[i])
	return f.stats[i], f.statErrors[i]
}
func (f *wr01ReadFile) Close() error {
	f.closeCalls++
	*f.trace = append(*f.trace, "close")
	return f.closeErr
}

type wr01AfterInfo struct {
	os.FileInfo
	size int64
	mod  time.Time
}

func (i wr01AfterInfo) Size() int64        { return i.size }
func (i wr01AfterInfo) ModTime() time.Time { return i.mod }

func TestWR01ControlledRead(t *testing.T) {
	dir := wr01PhysicalDir(t)
	path := filepath.Join(dir, "content")
	otherPath := filepath.Join(dir, "other")
	for _, p := range []string{path, otherPath} {
		if err := os.WriteFile(p, []byte("abc"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	other, err := os.Lstat(otherPath)
	if err != nil {
		t.Fatal(err)
	}
	directory, err := os.Lstat(dir)
	if err != nil {
		t.Fatal(err)
	}
	// Real FileInfo values preserve os.SameFile's identity checks. Only the
	// final descriptor metadata (not used for SameFile) needs an override.
	sizeChanged := wr01AfterInfo{info, 4, info.ModTime()}
	modChanged := wr01AfterInfo{info, info.Size(), info.ModTime().Add(time.Second)}
	cause := errors.New("WR01 injected OS failure")
	failure := func(op string) error { return &os.PathError{Op: op, Path: path, Err: cause} }
	readErr, closeErr, statErr, lstatErr, openErr := failure("read"), failure("close"), failure("stat"), failure("lstat"), failure("open")
	full := []string{"lstat-before", "open", "stat-open", "read", "stat-after", "lstat-current", "close"}
	preRead := []string{"lstat-before", "open", "stat-open", "close"}
	readOnly := []string{"lstat-before", "open", "stat-open", "read", "close"}
	for _, tc := range []struct {
		name                                                                                string
		readData                                                                            string
		readErr, closeErr, firstStatErr, afterStatErr, firstLstatErr, lastLstatErr, openErr error
		before, opened, after, current                                                      os.FileInfo
		limit                                                                               int
		wantData                                                                            []byte
		wantErr                                                                             error
		code, detail                                                                        string
		trace                                                                               []string
	}{
		{name: "success", readData: "abc", limit: 3, wantData: []byte("abc"), trace: full},
		{name: "close keeps bytes", readData: "abc", limit: 3, closeErr: closeErr, wantData: []byte("abc"), wantErr: closeErr, trace: full},
		{name: "partial read beats close", readData: "ab", limit: 3, readErr: readErr, closeErr: closeErr, wantErr: readErr, trace: readOnly},
		{name: "read without bytes", limit: 3, readErr: readErr, wantErr: readErr, trace: readOnly},
		{name: "first stat beats close", firstStatErr: statErr, closeErr: closeErr, wantErr: statErr, trace: preRead},
		{name: "replacement before read beats close", opened: other, closeErr: closeErr, code: "task_content_integrity_conflict", detail: "file changed before read", trace: preRead},
		{name: "limit beats close", readData: "abcd", limit: 3, closeErr: closeErr, code: "task_content_invalid_input", detail: "file exceeds 3 bytes", trace: readOnly},
		{name: "short read beats close", readData: "ab", limit: 3, closeErr: closeErr, code: "task_content_integrity_conflict", trace: full},
		{name: "descriptor size changed", readData: "abc", limit: 3, after: sizeChanged, code: "task_content_integrity_conflict", trace: full},
		{name: "descriptor modtime changed", readData: "abc", limit: 3, after: modChanged, code: "task_content_integrity_conflict", trace: full},
		{name: "path replaced after read", readData: "abc", limit: 3, current: other, code: "task_content_integrity_conflict", trace: full},
		{name: "path is now directory", readData: "abc", limit: 3, current: directory, code: "task_content_integrity_conflict", trace: full},
		{name: "post stat beats close", readData: "abc", limit: 3, afterStatErr: statErr, closeErr: closeErr, wantErr: statErr, trace: full},
		{name: "post lstat beats close", readData: "abc", limit: 3, lastLstatErr: lstatErr, closeErr: closeErr, wantErr: lstatErr, trace: full},
		{name: "both post observations attempted stat wins", readData: "abc", limit: 3, afterStatErr: statErr, lastLstatErr: lstatErr, closeErr: closeErr, wantErr: statErr, trace: full},
		{name: "first lstat no open", firstLstatErr: lstatErr, wantErr: lstatErr, trace: []string{"lstat-before"}},
		{name: "file type no open", before: directory, code: "task_content_integrity_conflict", detail: "unexpected file type or immutable mode: " + path, trace: []string{"lstat-before"}},
		{name: "open failure no close", openErr: openErr, wantErr: openErr, trace: []string{"lstat-before", "open"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var trace []string
			choose := func(value os.FileInfo) os.FileInfo {
				if value != nil {
					return value
				}
				return info
			}
			file := &wr01ReadFile{trace: &trace, data: []byte(tc.readData), readErr: tc.readErr, closeErr: tc.closeErr,
				stats: []os.FileInfo{choose(tc.opened), choose(tc.after)}, statErrors: []error{tc.firstStatErr, tc.afterStatErr}}
			lstatCalls, openCalls := 0, 0
			ops := contentReadOps{
				lstat: func(got string) (os.FileInfo, error) {
					if got != path {
						t.Fatalf("Lstat path = %s", got)
					}
					lstatCalls++
					if lstatCalls == 1 {
						trace = append(trace, "lstat-before")
						return choose(tc.before), tc.firstLstatErr
					}
					trace = append(trace, "lstat-current")
					return choose(tc.current), tc.lastLstatErr
				},
				open: func(got string) (contentReadFile, error) {
					if got != path {
						t.Fatalf("Open path = %s", got)
					}
					openCalls++
					trace = append(trace, "open")
					if tc.openErr != nil {
						return nil, tc.openErr
					}
					return file, nil
				},
			}
			data, err := contentReadBoundedWith(path, tc.limit, false, ops)
			if (data == nil) != (tc.wantData == nil) || !bytes.Equal(data, tc.wantData) {
				t.Fatalf("data = %q, want %q", data, tc.wantData)
			}
			if tc.code != "" {
				detail := tc.detail
				if detail == "" {
					detail = "file changed or was short during read"
				}
				wr01ContentError(t, err, tc.code, detail)
			} else if err != tc.wantErr {
				t.Fatalf("error = %#v, want original %#v", err, tc.wantErr)
			} else if err != nil {
				var pathError *os.PathError
				if !errors.Is(err, cause) || !errors.As(err, &pathError) || errors.Unwrap(err) != cause || err.Error() != tc.wantErr.Error() {
					t.Fatalf("OS error wrapping changed: %v", err)
				}
			}
			wr01Trace(t, trace, tc.trace)
			wantClose := openCalls
			if tc.openErr != nil {
				wantClose = 0
			}
			if file.closeCalls != wantClose {
				t.Fatalf("close count = %d, want %d", file.closeCalls, wantClose)
			}
		})
	}
	t.Run("physical validation before operations", func(t *testing.T) {
		data, err := contentReadBoundedWith("relative", 3, false, contentReadOps{
			lstat: func(string) (os.FileInfo, error) { t.Fatal("Lstat before path validation"); return nil, nil },
			open:  func(string) (contentReadFile, error) { t.Fatal("Open before path validation"); return nil, nil },
		})
		if data != nil {
			t.Fatalf("data = %q", data)
		}
		wr01ContentError(t, err, "task_content_invalid_input", "path must be absolute and clean")
	})
}

type wr01ReplacingFile struct {
	*os.File
	afterRead func()
	closes    int
}

func (f *wr01ReplacingFile) Read(p []byte) (int, error) {
	n, err := f.File.Read(p)
	if f.afterRead != nil {
		change := f.afterRead
		f.afterRead = nil
		change()
	}
	return n, err
}
func (f *wr01ReplacingFile) Close() error { f.closes++; return f.File.Close() }

func TestWR01RealReplacementAtReadBoundaries(t *testing.T) {
	for _, boundary := range []string{"before-open", "after-read"} {
		t.Run(boundary, func(t *testing.T) {
			dir := wr01PhysicalDir(t)
			path, replacement := filepath.Join(dir, "content"), filepath.Join(dir, "replacement")
			for _, p := range []string{path, replacement} {
				if err := os.WriteFile(p, []byte("abc"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			replace := func() {
				if err := os.Rename(replacement, path); err != nil {
					t.Fatal(err)
				}
			}
			var file *wr01ReplacingFile
			data, err := contentReadBoundedWith(path, 3, false, contentReadOps{
				lstat: os.Lstat,
				open: func(p string) (contentReadFile, error) {
					if boundary == "before-open" {
						replace()
					}
					opened, err := os.Open(p)
					if err != nil {
						return nil, err
					}
					file = &wr01ReplacingFile{File: opened}
					if boundary == "after-read" {
						file.afterRead = replace
					}
					return file, nil
				},
			})
			if data != nil || file == nil || file.closes != 1 {
				t.Fatalf("data = %q, file = %#v", data, file)
			}
			detail := "file changed before read"
			if boundary == "after-read" {
				detail = "file changed or was short during read"
			}
			wr01ContentError(t, err, "task_content_integrity_conflict", detail)
		})
	}
}
