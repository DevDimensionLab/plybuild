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

func TestLockAdaptersPreserveErrorParity(t *testing.T) {
	const p = "/ws/.ply/x.lock"
	e, u, c := errors.New("cause"), errors.New("unlock"), errors.New("close")
	lockErr := func(op string, injected bool, err error) *filelock.Error {
		return &filelock.Error{Op: op, Injected: injected, Path: p, Err: err}
	}
	project := func(class ErrorClass, detail string, err error) error {
		return &ProjectError{Class: class, Detail: detail, Err: err}
	}
	work := func(class WorkItemErrorClass, detail string, err error) error {
		return &WorkItemError{Class: class, Detail: detail, Err: err}
	}
	content := func(detail string, err error) error {
		return &TaskContentError{Code: "task_spec_lock_unavailable", Detail: detail, Err: err}
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
		{"work inspect", workItemLockError(p, lockErr("inspect", false, e)), work(ErrorWorkStoreConflict, "cannot inspect work-item lock "+p, e)},
		{"work not-regular", workItemLockError(p, lockErr("not-regular", false, nil)), work(ErrorWorkStoreConflict, "work-item lock "+p+" is not a regular file", nil)},
		{"work injected open", workItemLockError(p, lockErr("open", true, e)), work(ErrorWorkIO, "open work-item lock", e)},
		{"work open", workItemLockError(p, lockErr("open", false, e)), work(ErrorWorkIO, "open work-item lock "+p, e)},
		{"work stat-opened", workItemLockError(p, lockErr("stat-opened", false, e)), work(ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", e)},
		{"work lstat-path", workItemLockError(p, lockErr("lstat-path", false, e)), work(ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", e)},
		{"work changed", workItemLockError(p, lockErr("changed", false, nil)), work(ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", nil)},
		{"work injected lock", workItemLockError(p, lockErr("lock", true, e)), work(ErrorWorkIO, "lock work-item registry", e)},
		{"work lock", workItemLockError(p, lockErr("lock", false, e)), work(ErrorWorkIO, "lock work-item registry", e)},
		{"work release both", workItemReleaseError(&filelock.ReleaseError{Path: p, Unlock: u, Close: c}), work(ErrorWorkIO, "unlock work-item lock "+p, u)},
		{"work release close", workItemReleaseError(&filelock.ReleaseError{Path: p, Close: c, CloseInjected: true}), work(ErrorWorkIO, "close work-item lock "+p, c)},
		{"content stat-opened", contentLockError(lockErr("stat-opened", false, e)), content("existing lock identity differs", e)},
		{"content lstat-path", contentLockError(lockErr("lstat-path", false, e)), content("existing lock identity differs", e)},
		{"content changed", contentLockError(lockErr("changed", false, nil)), content("existing lock identity differs", nil)},
		{"content lock", contentLockError(lockErr("lock", false, e)), content("cannot acquire existing lock", e)},
		{"content release both", contentReleaseError(&filelock.ReleaseError{Path: p, Unlock: u, Close: c}), u},
		{"content release close", contentReleaseError(&filelock.ReleaseError{Path: p, Close: c}), c},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.got, test.want) {
			t.Errorf("%s: got %#v, want %#v", test.name, test.got, test.want)
		}
	}
}

func TestProjectsLockContentionBetweenFilelockAndConsumers(t *testing.T) {
	root, err := filepath.EvalSymlinks(emptyStoreRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, MarkerDirectory, projectsLock)
	store := newSystemProjectStore()
	project, repositories := storeFixture(root, "one", "One", "one")
	consumers := []func() error{
		func() error { _, _, _, err := store.Add(root, project, repositories); return err },
		func() error { lock, err := acquireExistingContentLock(path); releaseContent(lock); return err },
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
		func() (func(), error) {
			lock, err := acquireExistingContentLock(path)
			return func() { releaseContent(lock) }, err
		},
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

func releaseContent(lock *filelock.Lock) {
	if lock != nil {
		_ = releaseContentLock(lock)
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
