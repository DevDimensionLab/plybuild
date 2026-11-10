package filelock

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcquireCreatesMissingAndOpensExistingLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	lock, err := Acquire(path, Options{})
	if err != nil || lock.Path() != path {
		t.Fatalf("lock = %v, err = %v", lock, err)
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("created lock = %v, err = %v", info, err)
	}
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("kept"), 0o600); err != nil {
		t.Fatal(err)
	}
	lock, err = Acquire(path, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Release()
	if contents, _ := os.ReadFile(path); string(contents) != "kept" {
		t.Fatalf("existing lock contents = %q", contents)
	}
}

func TestAcquireRejectsSymlinkAndDirectory(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.lock")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{link, filepath.Join(dir, "sub")} {
		if path != link {
			if err := os.Mkdir(path, 0o700); err != nil {
				t.Fatal(err)
			}
		}
		_, err := Acquire(path, Options{})
		assertError(t, err, "not-regular", false, path, nil)
	}
}

func TestAcquireIdentityChecksAfterOpen(t *testing.T) {
	tests := []struct {
		op      string
		mutate  func(*testing.T, *os.File, string)
		wantErr bool
	}{
		{"stat-opened", func(t *testing.T, file *os.File, _ string) { _ = file.Close() }, true},
		{"lstat-path", func(t *testing.T, _ *os.File, path string) { remove(t, path) }, true},
		{"changed", func(t *testing.T, _ *os.File, path string) {
			remove(t, path)
			if err := os.WriteFile(path, nil, 0o600); err != nil {
				t.Fatal(err)
			}
		}, false},
	}
	for _, test := range tests {
		t.Run(test.op, func(t *testing.T) {
			for _, existing := range []bool{false, true} {
				path := filepath.Join(t.TempDir(), "x.lock")
				if err := os.WriteFile(path, nil, 0o600); err != nil {
					t.Fatal(err)
				}
				var opened *os.File
				hook := func(file *os.File) { opened = file; test.mutate(t, file, path) }
				var err error
				if existing {
					_, err = acquireExisting(path, hook)
				} else {
					_, err = Acquire(path, Options{beforeIdentity: hook})
				}
				var lockErr *Error
				if !errors.As(err, &lockErr) || lockErr.Op != test.op || lockErr.Injected || (lockErr.Err != nil) != test.wantErr {
					t.Fatalf("existing = %t, err = %#v", existing, err)
				}
				if closeErr := opened.Close(); !errors.Is(closeErr, os.ErrClosed) {
					t.Fatalf("opened file was not closed: %v", closeErr)
				}
			}
		})
	}
}

func TestAcquireInjectedAndActualFailures(t *testing.T) {
	fault := errors.New("fault")
	for _, test := range []struct{ fail, op string }{{"lock-open", "open"}, {"lock", "lock"}} {
		path := filepath.Join(t.TempDir(), "x.lock")
		_, err := Acquire(path, Options{Fail: failAt(test.fail, fault)})
		assertError(t, err, test.op, true, path, fault)
		lock, err := Acquire(path, Options{})
		if err != nil {
			t.Fatalf("%s: lock was retained after failure: %v", test.fail, err)
		}
		lock.Release()
	}
	missingParent := filepath.Join(t.TempDir(), "missing", "x.lock")
	_, err := Acquire(missingParent, Options{Fail: failAt("lock", fault)})
	var lockErr *Error
	if !errors.As(err, &lockErr) || lockErr.Op != "open" || lockErr.Injected || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("actual open err = %#v", err)
	}
}

func TestReleasePreservesUnlockAndCloseErrors(t *testing.T) {
	fault := errors.New("fault")
	both := func(op string) error { return fault }
	tests := []struct {
		name                string
		fail                func(string) error
		closeFirst          bool
		unlockInj, closeInj bool
	}{
		{name: "injected unlock", fail: failAt("unlock", fault), unlockInj: true},
		{name: "injected close", fail: failAt("lock-close", fault), closeInj: true},
		{name: "injected both", fail: both, unlockInj: true, closeInj: true},
		{name: "actual both overrides injection", fail: both, closeFirst: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "x.lock")
			lock, err := Acquire(path, Options{Fail: func(op string) error {
				if op == "unlock" || op == "lock-close" {
					return test.fail(op)
				}
				return nil
			}})
			if err != nil {
				t.Fatal(err)
			}
			if test.closeFirst {
				_ = lock.file.Close()
			}
			releaseErr := lock.Release()
			if releaseErr == nil || releaseErr.Path != path || releaseErr.UnlockInjected != test.unlockInj || releaseErr.CloseInjected != test.closeInj {
				t.Fatalf("release err = %#v", releaseErr)
			}
			wantUnlock := test.unlockInj || test.closeFirst
			wantClose := test.closeInj || test.closeFirst
			if (releaseErr.Unlock != nil) != wantUnlock || (releaseErr.Close != nil) != wantClose {
				t.Fatalf("release err = %#v", releaseErr)
			}
			if test.unlockInj && releaseErr.Unlock != fault || test.closeInj && releaseErr.Close != fault {
				t.Fatalf("injected faults were not preserved: %#v", releaseErr)
			}
			if test.closeFirst && !errors.Is(releaseErr.Close, os.ErrClosed) {
				t.Fatalf("actual close err = %v", releaseErr.Close)
			}
			if wrapped := releaseErr.Unwrap(); wantUnlock && wantClose && (len(wrapped) != 2 || wrapped[0] != releaseErr.Unlock || wrapped[1] != releaseErr.Close) {
				t.Fatalf("unwrap order = %v", wrapped)
			}
		})
	}
}

func TestRawCloseReleasesLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	lock, err := Acquire(path, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := lock.file.Close(); err != nil {
		t.Fatal(err)
	}
	next := acquireAsync(func() (*Lock, error) { return Acquire(path, Options{}) })
	select {
	case result := <-next:
		if result.err != nil {
			t.Fatal(result.err)
		}
		result.lock.Release()
	case <-time.After(10 * time.Second):
		t.Fatal("acquire remained blocked after the lock handle closed")
	}
}

func TestHeldLockBlocksAcquireAndAcquireExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	for _, second := range []func() (*Lock, error){
		func() (*Lock, error) { return Acquire(path, Options{}) },
		func() (*Lock, error) { return AcquireExisting(path) },
	} {
		held, err := Acquire(path, Options{})
		if err != nil {
			t.Fatal(err)
		}
		next := acquireAsync(second)
		select {
		case <-next:
			t.Fatal("second acquire did not block while the lock was held")
		case <-time.After(100 * time.Millisecond):
		}
		if err := held.Release(); err != nil {
			t.Fatal(err)
		}
		select {
		case result := <-next:
			if result.err != nil {
				t.Fatal(result.err)
			}
			result.lock.Release()
		case <-time.After(10 * time.Second):
			t.Fatal("second acquire remained blocked after release")
		}
	}
}

func TestAcquireExistingMissingReturnsRawOpenErrorWithoutCreating(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.lock")
	_, err := AcquireExisting(path)
	var lockErr *Error
	var pathErr *fs.PathError
	if errors.As(err, &lockErr) || !errors.As(err, &pathErr) || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("err = %#v", err)
	}
	if _, statErr := os.Lstat(path); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("lock was created: %v", statErr)
	}
}

type acquired struct {
	lock *Lock
	err  error
}

func acquireAsync(acquire func() (*Lock, error)) <-chan acquired {
	done := make(chan acquired, 1)
	go func() {
		lock, err := acquire()
		done <- acquired{lock, err}
	}()
	return done
}

func failAt(op string, fault error) func(string) error {
	return func(current string) error {
		if current == op {
			return fault
		}
		return nil
	}
}

func assertError(t *testing.T, err error, op string, injected bool, path string, cause error) {
	t.Helper()
	var lockErr *Error
	if !errors.As(err, &lockErr) || lockErr.Op != op || lockErr.Injected != injected || lockErr.Path != path || lockErr.Err != cause {
		t.Fatalf("err = %#v, want op %s injected %t", err, op, injected)
	}
}

func remove(t *testing.T, path string) {
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}
