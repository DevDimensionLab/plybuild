package workspace

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

func TestProjectLockAdaptersPreserveErrorParity(t *testing.T) {
	t.Parallel()
	const p = "/ws/.ply/x.lock"
	e, u, c := errors.New("cause"), errors.New("unlock"), errors.New("close")
	lockErr := func(op string, injected bool, err error) *filelock.Error {
		return &filelock.Error{Op: op, Injected: injected, Path: p, Err: err}
	}
	project := func(class ErrorClass, detail string, err error) error {
		return &ProjectError{Class: class, Detail: detail, Err: err}
	}
	tests := []struct {
		name      string
		got, want error
	}{
		{"project inspect", projectLockError(p, lockErr("inspect", false, e)), project(ErrorProjectStoreConflict, p+": cause", e)},
		{"project not-regular", projectLockError(p, lockErr("not-regular", false, nil)), project(ErrorProjectStoreConflict, p+": project lock is not a regular file", nil)},
		{"project stat-opened", projectLockError(p, lockErr("stat-opened", false, e)), project(ErrorProjectStoreConflict, p+": cannot validate opened project lock", e)},
		{"project lstat-path", projectLockError(p, lockErr("lstat-path", false, e)), project(ErrorProjectStoreConflict, p+": cannot validate project lock path", e)},
		{"project changed", projectLockError(p, lockErr("changed", false, nil)), project(ErrorProjectStoreConflict, p+": project lock changed or is not a regular file", nil)},
		{"project injected open", projectLockError(p, lockErr("open", true, e)), project(ErrorProjectIO, "open lock "+p+": cause", e)},
		{"project open", projectLockError(p, lockErr("open", false, e)), project(ErrorProjectIO, "open lock "+p+": cause", e)},
		{"project injected lock", projectLockError(p, lockErr("lock", true, e)), project(ErrorProjectIO, "lock "+p+": cause", e)},
		{"project lock", projectLockError(p, lockErr("lock", false, e)), project(ErrorProjectIO, "lock "+p+": cause", e)},
		{"project release both", projectReleaseError(&filelock.ReleaseError{Path: p, Unlock: u, Close: c}), project(ErrorProjectIO, "unlock "+p+": unlock; close: close", fmt.Errorf("%w; close: %v", u, c))},
		{"project release unlock", projectReleaseError(&filelock.ReleaseError{Path: p, Unlock: u, UnlockInjected: true}), project(ErrorProjectIO, "unlock "+p+": unlock", u)},
		{"project release close", projectReleaseError(&filelock.ReleaseError{Path: p, Close: c}), project(ErrorProjectIO, "close lock "+p+": close", c)},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.got, test.want) {
			t.Errorf("%s: got %#v, want %#v", test.name, test.got, test.want)
		}
	}
}

func TestProjectsLockContentionBetweenFilelockAndConsumers(t *testing.T) {
	t.Parallel()
	root := createProjectWorkspace(t)
	path := filepath.Join(root, MarkerDirectory, projectsLock)
	store := newSystemProjectStore()
	consumers := []func() error{
		func() error { return store.WithSnapshotLock(root, func(ProjectSnapshot) error { return nil }) },
	}
	for _, consumer := range consumers {
		held, err := filelock.Acquire(path, filelock.Options{})
		if err != nil {
			t.Fatal(err)
		}
		assertBlocksUntilRelease(t, consumer, func() { held.Release() })
	}
	holders := []func() (func(), error){
		func() (func(), error) { lock, err := store.acquire(root); return func() { store.release(lock) }, err },
	}
	for _, holder := range holders {
		release, err := holder()
		if err != nil {
			t.Fatal(err)
		}
		assertBlocksUntilRelease(t, func() error {
			lock, err := filelock.Acquire(path, filelock.Options{})
			if err == nil {
				lock.Release()
			}
			return err
		}, release)
	}
}

func assertBlocksUntilRelease(t *testing.T, acquire func() error, release func()) {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- acquire() }()
	select {
	case err := <-done:
		t.Fatalf("acquire did not block while the lock was held: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	release()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("acquire remained blocked after release")
	}
}
