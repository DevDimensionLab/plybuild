package workspacedef

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

const sample = `{
  "format_version": 1,
  "projects": [
    {
      "id": "ply",
      "name": "Ply",
      "wrapper": "ply",
      "data_paths": {"backlog": "planning/backlog"},
      "repos": [{"id": "ply", "path": "ply/main"}],
      "epics": [{"id": "ply-master", "title": "Master", "repo": "ply", "path": "ply/main", "ref": "refs/heads/master"}]
    },
    {"id": "app", "name": "App", "wrapper": "app", "repos": [{"id": "app", "path": "app/main"}], "epics": []}
  ]
}
`

func writeRoot(t *testing.T, contents string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".ply"), 0o755); err != nil {
		t.Fatal(err)
	}
	if contents != "" {
		if err := os.WriteFile(filepath.Join(root, File), []byte(contents), 0o640); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

type fileState struct {
	raw   []byte
	mtime time.Time
}

func stateOf(t *testing.T, root string) fileState {
	t.Helper()
	path := filepath.Join(root, File)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return fileState{raw: raw, mtime: info.ModTime()}
}

func (s fileState) requireSame(t *testing.T, root string) {
	t.Helper()
	now := stateOf(t, root)
	if !bytes.Equal(now.raw, s.raw) || !now.mtime.Equal(s.mtime) {
		t.Fatalf("definitions changed:\n%s", now.raw)
	}
}

func requireKind(t *testing.T, err error, kind Kind) {
	t.Helper()
	var typed *Error
	if !errors.As(err, &typed) || typed.Kind != kind {
		t.Fatalf("error %v, want kind %s", err, kind)
	}
}

func basisOf(t *testing.T, root string) string {
	t.Helper()
	_, digest, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func addRepo(project, id, path string) func(*Definitions) (bool, error) {
	return func(d *Definitions) (bool, error) { return AddRepo(d, project, id, path) }
}

func accept(Definitions) error { return nil }

func TestReadAbsentAndPresent(t *testing.T) {
	defs, digest, err := Read(writeRoot(t, ""))
	if defs != nil || digest != Absent || err != nil {
		t.Fatalf("absent: %v %q %v", defs, digest, err)
	}
	root := writeRoot(t, sample)
	defs, digest, err = Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if digest != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(sample))) {
		t.Fatalf("digest %s", digest)
	}
	if len(defs.Projects) != 2 || defs.Projects[0].DataPaths["backlog"] != "planning/backlog" || defs.Projects[0].Epics[0].Ref != "refs/heads/master" || defs.Projects[1].Repos[0].Path != "app/main" {
		t.Fatalf("decoded %+v", defs)
	}
}

func TestReadRejectsInOrder(t *testing.T) {
	for name, tc := range map[string]struct {
		contents string
		kind     Kind
	}{
		"invalid JSON":        {`{"format_version": 1, "projects": [}`, KindInvalid},
		"duplicate key":       {`{"format_version": 1, "format_version": 1, "projects": []}`, KindInvalid},
		"unknown field":       {`{"format_version": 1, "projects": [], "extra": true}`, KindSchema},
		"unknown nested":      {`{"format_version": 1, "projects": [{"id": "a", "name": "A", "wrapper": "a", "repos": [{"id": "a", "path": "a", "url": "x"}], "epics": []}]}`, KindSchema},
		"wrong type":          {`{"format_version": "1", "projects": []}`, KindSchema},
		"format 2":            {`{"format_version": 2, "projects": []}`, KindFormat},
		"projects null":       {`{"format_version": 1, "projects": null}`, KindFormat},
		"projects missing":    {`{"format_version": 1}`, KindFormat},
		"unknown then format": {`{"format_version": 2, "projects": [], "extra": 1}`, KindSchema},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := Read(writeRoot(t, tc.contents))
			requireKind(t, err, tc.kind)
		})
	}
}

func TestReadRejectsNonRegular(t *testing.T) {
	directory := writeRoot(t, "")
	if err := os.Mkdir(filepath.Join(directory, File), 0o755); err != nil {
		t.Fatal(err)
	}
	_, _, err := Read(directory)
	requireKind(t, err, KindNotRegular)

	target := writeRoot(t, sample)
	link := writeRoot(t, "")
	if err := os.Symlink(filepath.Join(target, File), filepath.Join(link, File)); err != nil {
		t.Fatal(err)
	}
	_, _, err = Read(link)
	requireKind(t, err, KindNotRegular)
}

func TestReadReturnsRawInspectError(t *testing.T) {
	root := writeRoot(t, "")
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := Read(file)
	var typed *Error
	if err == nil || errors.As(err, &typed) || errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("error %v", err)
	}
}

func TestAddRepo(t *testing.T) {
	defs := Definitions{FormatVersion: 1, Projects: []Project{{ID: "ply", Repos: []Repo{{ID: "ply", Path: "ply/main"}}}}}
	if _, err := AddRepo(&defs, "other", "x", "x"); err == nil {
		t.Fatal("missing Project accepted")
	} else {
		requireKind(t, err, KindNotFound)
	}
	if changed, err := AddRepo(&defs, "ply", "ply", "ply/main"); changed || err != nil {
		t.Fatalf("identical: %v %v", changed, err)
	}
	if _, err := AddRepo(&defs, "ply", "ply", "ply/other"); err == nil {
		t.Fatal("different path accepted")
	} else {
		requireKind(t, err, KindConflict)
	}
	if changed, err := AddRepo(&defs, "ply", "map", "map"); !changed || err != nil {
		t.Fatalf("new: %v %v", changed, err)
	}
	if want := []Repo{{ID: "ply", Path: "ply/main"}, {ID: "map", Path: "map"}}; !reflect.DeepEqual(defs.Projects[0].Repos, want) {
		t.Fatalf("repos %+v", defs.Projects[0].Repos)
	}
}

func TestRemoveRepo(t *testing.T) {
	defs := Definitions{FormatVersion: 1, Projects: []Project{{ID: "ply", Repos: []Repo{{ID: "ply", Path: "ply/main"}, {ID: "map", Path: "map/main"}}}}}
	if _, err := RemoveRepo(&defs, "other", "map", "map/main"); err == nil {
		t.Fatal("missing Project accepted")
	} else {
		requireKind(t, err, KindNotFound)
	}
	if _, err := RemoveRepo(&defs, "ply", "map", "map/other"); err == nil {
		t.Fatal("different path accepted")
	} else {
		requireKind(t, err, KindConflict)
	}
	if changed, err := RemoveRepo(&defs, "ply", "absent", "absent/main"); changed || err != nil {
		t.Fatalf("absent: %v %v", changed, err)
	}
	if changed, err := RemoveRepo(&defs, "ply", "map", "map/main"); !changed || err != nil {
		t.Fatalf("exact: %v %v", changed, err)
	}
	if want := []Repo{{ID: "ply", Path: "ply/main"}}; !reflect.DeepEqual(defs.Projects[0].Repos, want) {
		t.Fatalf("repos %+v", defs.Projects[0].Repos)
	}
}

func TestUpdateAddsRepoAndPreservesSemantics(t *testing.T) {
	root := writeRoot(t, sample)
	before, basis, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Update(root, basis, addRepo("ply", "ply-code-map", "ply-code-map"), accept)
	if err != nil || result != (Result{Effect: Completed, Changed: true}) {
		t.Fatalf("result %+v, error %v", result, err)
	}
	after, _, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	before.Projects[0].Repos = append(before.Projects[0].Repos, Repo{ID: "ply-code-map", Path: "ply-code-map"})
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("semantics changed:\n%+v\n%+v", before, after)
	}
	raw, _ := os.ReadFile(filepath.Join(root, File))
	want, _ := json.MarshalIndent(after, "", "  ")
	if !bytes.Equal(raw, append(want, '\n')) {
		t.Fatalf("not indented JSON with newline:\n%s", raw)
	}
	info, _ := os.Stat(filepath.Join(root, File))
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("mode %v", info.Mode())
	}
}

func TestUpdateUnchangedValidatesWithoutWriting(t *testing.T) {
	root := writeRoot(t, sample)
	state := stateOf(t, root)
	checked := 0
	check := func(Definitions) error { checked++; return nil }
	result, err := Update(root, basisOf(t, root), addRepo("ply", "ply", "ply/main"), check)
	if err != nil || result != (Result{Effect: Completed}) || checked != 1 {
		t.Fatalf("result %+v, error %v, %d checks", result, err, checked)
	}
	state.requireSame(t, root)

	failure := errors.New("current invalid")
	result, err = Update(root, basisOf(t, root), addRepo("ply", "ply", "ply/main"), func(Definitions) error { return failure })
	if !errors.Is(err, failure) || result.Effect != NotAttempted {
		t.Fatalf("result %+v, error %v", result, err)
	}
	state.requireSame(t, root)
}

func TestUpdateRequiresCallbacksBeforeLock(t *testing.T) {
	root := t.TempDir() // no .ply: acquiring the lock would fail first
	for _, tc := range []struct {
		mutate func(*Definitions) (bool, error)
		check  func(Definitions) error
	}{{nil, accept}, {addRepo("ply", "x", "x"), nil}} {
		result, err := Update(root, Absent, tc.mutate, tc.check)
		requireKind(t, err, KindUsage)
		if result.Effect != NotAttempted {
			t.Fatalf("result %+v", result)
		}
	}
	if _, err := os.Lstat(filepath.Join(root, ".ply")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("lock side effect: %v", err)
	}
}

func TestUpdateRejectsMissingAndChangedBasis(t *testing.T) {
	missing := writeRoot(t, "")
	result, err := Update(missing, Absent, addRepo("ply", "x", "x"), accept)
	requireKind(t, err, KindMissing)
	if result.Effect != NotAttempted {
		t.Fatalf("result %+v", result)
	}
	if _, err := os.Lstat(filepath.Join(missing, File)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("definitions created: %v", err)
	}

	root := writeRoot(t, sample)
	state := stateOf(t, root)
	result, err = Update(root, "sha256:stale", addRepo("ply", "x", "x"), accept)
	requireKind(t, err, KindConflict)
	if result.Effect != NotAttempted {
		t.Fatalf("result %+v", result)
	}
	state.requireSame(t, root)
}

func TestUpdateReturnsReadErrorsUnderLock(t *testing.T) {
	root := writeRoot(t, `{"format_version": 2, "projects": []}`)
	state := stateOf(t, root)
	result, err := Update(root, "sha256:any", addRepo("ply", "x", "x"), accept)
	requireKind(t, err, KindFormat)
	if result.Effect != NotAttempted {
		t.Fatalf("result %+v", result)
	}
	state.requireSame(t, root)
}

func TestUpdateCandidateCheckPreventsWrite(t *testing.T) {
	root := writeRoot(t, sample)
	state := stateOf(t, root)
	failure := errors.New("candidate invalid")
	calls := 0
	check := func(d Definitions) error {
		calls++
		if len(d.Projects[0].Repos) == 2 {
			return failure
		}
		return nil
	}
	result, err := Update(root, basisOf(t, root), addRepo("ply", "x", "x"), check)
	if !errors.Is(err, failure) || result.Effect != NotAttempted || calls != 2 {
		t.Fatalf("result %+v, error %v, %d checks", result, err, calls)
	}
	state.requireSame(t, root)
}

func TestUpdateMutatesACopy(t *testing.T) {
	root := writeRoot(t, sample)
	var current Definitions
	calls := 0
	check := func(d Definitions) error {
		calls++
		if calls == 1 {
			current = d
		}
		return nil
	}
	mutate := func(d *Definitions) (bool, error) {
		d.Projects[0].Repos[0].Path = "changed"
		d.Projects[0].DataPaths["backlog"] = "changed"
		d.Projects[0].Epics[0].Title = "changed"
		return false, errors.New("stop")
	}
	if _, err := Update(root, basisOf(t, root), mutate, check); err == nil {
		t.Fatal("mutate error ignored")
	}
	if current.Projects[0].Repos[0].Path != "ply/main" || current.Projects[0].DataPaths["backlog"] != "planning/backlog" || current.Projects[0].Epics[0].Title != "Master" {
		t.Fatalf("current definitions were mutated: %+v", current.Projects[0])
	}
}

func TestUpdateRejectsChangeBeforeReplace(t *testing.T) {
	root := writeRoot(t, sample)
	path := filepath.Join(root, File)
	edited := sample + " "
	hooks := systemHooks()
	hooks.beforeFresh = func() {
		if err := os.WriteFile(path, []byte(edited), 0o640); err != nil {
			t.Fatal(err)
		}
	}
	replaced := false
	hooks.replace = func(string, []byte, fs.FileMode) error { replaced = true; return nil }
	result, err := update(root, basisOf(t, root), addRepo("ply", "x", "x"), accept, hooks)
	requireKind(t, err, KindConflict)
	if result.Effect != NotAttempted || replaced {
		t.Fatalf("result %+v, replaced %v", result, replaced)
	}
	if raw, _ := os.ReadFile(path); string(raw) != edited {
		t.Fatalf("edit lost: %s", raw)
	}
}

func TestUpdatePublicationOutcomes(t *testing.T) {
	failure := errors.New("replace failed")
	for name, tc := range map[string]struct {
		replace func(t *testing.T, path string, raw []byte) error
		want    Result
		kind    Kind
	}{
		"failed with basis readback": {
			replace: func(*testing.T, string, []byte) error { return failure },
			want:    Result{Effect: NoEffect}, kind: KindPublish,
		},
		"failed after candidate was published": {
			replace: func(t *testing.T, path string, raw []byte) error {
				if err := os.WriteFile(path, raw, 0o640); err != nil {
					t.Fatal(err)
				}
				return failure
			},
			want: Result{Effect: Completed, Changed: true},
		},
		"other content after failure": {
			replace: func(t *testing.T, path string, raw []byte) error {
				if err := os.WriteFile(path, []byte(`{"format_version": 1, "projects": []}`), 0o640); err != nil {
					t.Fatal(err)
				}
				return failure
			},
			want: Result{Effect: Unknown}, kind: KindReadback,
		},
		"unreadable after success": {
			replace: func(t *testing.T, path string, raw []byte) error {
				if err := os.WriteFile(path, []byte("{"), 0o640); err != nil {
					t.Fatal(err)
				}
				return nil
			},
			want: Result{Effect: Unknown}, kind: KindReadback,
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := writeRoot(t, sample)
			hooks := systemHooks()
			hooks.replace = func(path string, raw []byte, mode fs.FileMode) error {
				if mode != 0o640 {
					t.Fatalf("mode %v", mode)
				}
				return tc.replace(t, path, raw)
			}
			result, err := update(root, basisOf(t, root), addRepo("ply", "x", "x"), accept, hooks)
			if result != tc.want {
				t.Fatalf("result %+v, error %v", result, err)
			}
			if tc.kind == "" {
				if err != nil {
					t.Fatalf("error %v", err)
				}
				return
			}
			requireKind(t, err, tc.kind)
			if !errors.Is(err, failure) && name != "unreadable after success" {
				t.Fatalf("replace error not kept: %v", err)
			}
		})
	}
}

func TestConcurrentUpdatesWithSameBasis(t *testing.T) {
	root := writeRoot(t, sample)
	basis := basisOf(t, root)
	results := make([]Result, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i, id := range []string{"first", "second"} {
		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()
			results[i], errs[i] = Update(root, basis, addRepo("ply", id, id), accept)
		}(i, id)
	}
	wg.Wait()
	completed, conflicted := -1, -1
	for i := range results {
		var typed *Error
		switch {
		case errs[i] == nil && results[i] == (Result{Effect: Completed, Changed: true}):
			completed = i
		case errors.As(errs[i], &typed) && typed.Kind == KindConflict && results[i].Effect == NotAttempted:
			conflicted = i
		}
	}
	if completed < 0 || conflicted < 0 {
		t.Fatalf("results %+v, errors %v", results, errs)
	}
	retry := []string{"first", "second"}[conflicted]
	if result, err := Update(root, basisOf(t, root), addRepo("ply", retry, retry), accept); err != nil || !result.Changed {
		t.Fatalf("retry %+v %v", result, err)
	}
	defs, _, _ := Read(root)
	if len(defs.Projects[0].Repos) != 3 {
		t.Fatalf("repos %+v", defs.Projects[0].Repos)
	}
}

func TestUpdateWaitsForProjectsLock(t *testing.T) {
	root := writeRoot(t, sample)
	basis := basisOf(t, root)
	lock, err := filelock.Acquire(filepath.Join(root, ".ply", "projects.lock"), filelock.Options{})
	if err != nil {
		t.Fatal(err)
	}
	state := stateOf(t, root)
	done := make(chan error, 1)
	go func() {
		_, err := Update(root, basis, addRepo("ply", "x", "x"), accept)
		done <- err
	}()
	select {
	case err := <-done:
		t.Fatalf("update did not wait for the lock: %v", err)
	case <-time.After(200 * time.Millisecond):
	}
	state.requireSame(t, root)
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestUpdateReportsLockFailure(t *testing.T) {
	root := writeRoot(t, sample)
	lockPath := filepath.Join(root, ".ply", "projects.lock")
	if err := os.Mkdir(lockPath, 0o755); err != nil {
		t.Fatal(err)
	}
	state := stateOf(t, root)
	result, err := Update(root, basisOf(t, root), addRepo("ply", "x", "x"), accept)
	var lockErr *filelock.Error
	if !errors.As(err, &lockErr) || result.Effect != NotAttempted {
		t.Fatalf("result %+v, error %v", result, err)
	}
	state.requireSame(t, root)
}

func TestWithLockReturnsCallbackErrorUnchanged(t *testing.T) {
	root := writeRoot(t, sample)
	cause := errors.New("callback failed")
	if err := WithLock(root, func() error { return cause }); err != cause {
		t.Fatalf("error %v", err)
	}
	if err := WithLock(root, func() error { return nil }); err != nil {
		t.Fatal(err)
	}
}

func TestWithLockRequiresCallbackAndLock(t *testing.T) {
	root := t.TempDir() // no .ply: acquiring the lock fails
	called := false
	err := WithLock(root, func() error { called = true; return nil })
	var lockErr *filelock.Error
	if !errors.As(err, &lockErr) || called {
		t.Fatalf("error %v, called %v", err, called)
	}
	requireKind(t, WithLock(root, nil), KindUsage)
	withPly := writeRoot(t, sample)
	requireKind(t, WithLock(withPly, nil), KindUsage)
	if _, err := os.Lstat(filepath.Join(withPly, ".ply", "projects.lock")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("lock side effect: %v", err)
	}
}

func TestWithLockSerializesWithUpdate(t *testing.T) {
	root := writeRoot(t, sample)
	basis := basisOf(t, root)
	state := stateOf(t, root)
	entered, leave := make(chan struct{}), make(chan struct{})
	held := make(chan error, 1)
	go func() {
		held <- WithLock(root, func() error {
			close(entered)
			<-leave
			return nil
		})
	}()
	<-entered
	started := make(chan struct{})
	done := make(chan Result, 1)
	go func() {
		close(started)
		result, err := Update(root, basis, addRepo("ply", "x", "x"), accept)
		if err != nil {
			t.Error(err)
		}
		done <- result
	}()
	<-started
	select {
	case result := <-done:
		t.Fatalf("update did not wait for WithLock: %+v", result)
	case <-time.After(200 * time.Millisecond):
	}
	state.requireSame(t, root)
	close(leave)
	if err := <-held; err != nil {
		t.Fatal(err)
	}
	select {
	case result := <-done:
		if result != (Result{Effect: Completed, Changed: true}) {
			t.Fatalf("result %+v", result)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("update did not finish after release")
	}
}

func TestErrorMessage(t *testing.T) {
	cause := errors.New("cause")
	if got := (&Error{Kind: KindPublish, Path: "/w/" + File, Detail: "not replaced", Err: cause}).Error(); got != "publish /w/ply-workspace.json: not replaced: cause" {
		t.Fatalf("message %q", got)
	}
	if got := (&Error{Kind: KindNotFound, Detail: "Project x is not defined"}).Error(); got != "not_found: Project x is not defined" {
		t.Fatalf("message %q", got)
	}
}
