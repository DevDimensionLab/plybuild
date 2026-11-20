package taskwork

import (
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkItemLockAdaptersPreserveErrorParity(t *testing.T) {
	t.Parallel()
	const p = "/ws/.ply/x.lock"
	e, u, c := errors.New("cause"), errors.New("unlock"), errors.New("close")
	lockErr := func(op string, injected bool, err error) *filelock.Error {
		return &filelock.Error{Op: op, Injected: injected, Path: p, Err: err}
	}
	work := func(class workspace.WorkItemErrorClass, detail string, err error) error {
		return &workspace.WorkItemError{Class: class, Detail: detail, Err: err}
	}
	tests := []struct {
		name      string
		got, want error
	}{
		{"work inspect", workItemLockError(p, lockErr("inspect", false, e)), work(workspace.ErrorWorkStoreConflict, "cannot inspect work-item lock "+p, e)},
		{"work not-regular", workItemLockError(p, lockErr("not-regular", false, nil)), work(workspace.ErrorWorkStoreConflict, "work-item lock "+p+" is not a regular file", nil)},
		{"work injected open", workItemLockError(p, lockErr("open", true, e)), work(workspace.ErrorWorkIO, "open work-item lock", e)},
		{"work open", workItemLockError(p, lockErr("open", false, e)), work(workspace.ErrorWorkIO, "open work-item lock "+p, e)},
		{"work stat-opened", workItemLockError(p, lockErr("stat-opened", false, e)), work(workspace.ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", e)},
		{"work lstat-path", workItemLockError(p, lockErr("lstat-path", false, e)), work(workspace.ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", e)},
		{"work changed", workItemLockError(p, lockErr("changed", false, nil)), work(workspace.ErrorWorkStoreConflict, "work-item lock "+p+" changed or is not a regular file", nil)},
		{"work injected lock", workItemLockError(p, lockErr("lock", true, e)), work(workspace.ErrorWorkIO, "lock work-item registry", e)},
		{"work lock", workItemLockError(p, lockErr("lock", false, e)), work(workspace.ErrorWorkIO, "lock work-item registry", e)},
		{"work release both", workItemReleaseError(&filelock.ReleaseError{Path: p, Unlock: u, Close: c}), work(workspace.ErrorWorkIO, "unlock work-item lock "+p, u)},
		{"work release close", workItemReleaseError(&filelock.ReleaseError{Path: p, Close: c, CloseInjected: true}), work(workspace.ErrorWorkIO, "close work-item lock "+p, c)},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.got, test.want) {
			t.Errorf("%s: got %#v, want %#v", test.name, test.got, test.want)
		}
	}
}
