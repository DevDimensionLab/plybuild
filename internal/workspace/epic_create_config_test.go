package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

const epicCreateConfig = `{"format_version":1,"projects":[{"id":"sample","name":"Sample","wrapper":".","data_paths":{"backlog":"work/backlog"},"repos":[{"id":"code","path":"repo"}],"epics":[{"id":"feature","title":"Feature","repo":"code","path":"feature","ref":"refs/heads/feature"}]}]}`

func epicCreateConfigFixture(t *testing.T) (string, Dependencies, EpicCreateInput) {
	t.Helper()
	root, repo, d := durableFixture(t)
	runLocalGit(t, repo, "worktree", "add", "-b", "feature", filepath.Join(root, "feature"))
	if err := os.WriteFile(filepath.Join(root, DefinitionsFile), []byte(epicCreateConfig), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	upgradeEpicCreateFixture(t, d, root)
	in := EpicCreateInput{EpicID: "new-epic", Title: "New Epic", ProjectID: "sample", RepoID: "code", ParentRef: "refs/heads/main", ExpectedParentOID: runLocalGit(t, repo, "rev-parse", "HEAD"), Branch: "new-epic", Path: filepath.Join(root, "new-epic")}
	return root, d, in
}

func readConfigFixture(t *testing.T, root string) *workspacedef.Definitions {
	t.Helper()
	defs, _, err := readDefinitions(root)
	if err != nil {
		t.Fatal(err)
	}
	return defs
}

func expectedConfigWithNewEpic(t *testing.T, root string) *workspacedef.Definitions {
	t.Helper()
	defs, _, err := readDefinitions(root)
	if err != nil {
		t.Fatal(err)
	}
	defs.Projects[0].Epics = append(defs.Projects[0].Epics, workspacedef.Epic{ID: "new-epic", Title: "New Epic", Repo: "code", Path: "new-epic", Ref: "refs/heads/new-epic"})
	return defs
}

func TestEpicCreateWithActiveConfigRecordsDefinition(t *testing.T) {
	root, d, in := epicCreateConfigFixture(t)
	want := expectedConfigWithNewEpic(t, root)
	before, err := ShowEpic(d, "feature")
	if err != nil {
		t.Fatal(err)
	}
	got, err := CreateEpic(d, in)
	if err != nil || got.Outcome != "created" || got.Effects != (EpicCreateEffects{Worktree: "completed", Config: "completed"}) {
		t.Fatalf("%+v %v", got, err)
	}
	if defs := readConfigFixture(t, root); !reflect.DeepEqual(defs, want) {
		t.Fatalf("config %+v", defs)
	}
	epics, err := ListEpics(d)
	if err != nil || len(epics.Epics) != 2 {
		t.Fatalf("%+v %v", epics, err)
	}
	after, err := ShowEpic(d, "feature")
	if err != nil || !reflect.DeepEqual(before.Epic, after.Epic) {
		t.Fatalf("unrelated Epic changed: %+v %v", after, err)
	}
	raw, err := os.ReadFile(filepath.Join(root, DefinitionsFile))
	if err != nil {
		t.Fatal(err)
	}
	again, err := CreateEpic(d, in)
	if err != nil || again.Outcome != "already_created" || again.Effects != (EpicCreateEffects{Worktree: "not_attempted", Config: "not_attempted"}) {
		t.Fatalf("%+v %v", again, err)
	}
	if unchanged, _ := os.ReadFile(filepath.Join(root, DefinitionsFile)); string(unchanged) != string(raw) {
		t.Fatal("exact retry rewrote config")
	}
}

type epicPublishCrashStore struct{ WorkItemStore }
type epicPublishCrashSession struct{ WorkItemStoreSession }

func (s epicPublishCrashStore) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session WorkItemStoreSession) error { return fn(epicPublishCrashSession{session}) })
}
func (s epicPublishCrashSession) Publish(r WorkItemRegistry) error {
	if epic, _ := findEpic(r, "new-epic"); epic != nil {
		return errors.New("interrupted before runtime Epic")
	}
	return s.WorkItemStoreSession.Publish(r)
}

func TestEpicCreateWithActiveConfigResumesAfterInterruption(t *testing.T) {
	root, d, in := epicCreateConfigFixture(t)
	want := expectedConfigWithNewEpic(t, root)
	crashing := d
	crashing.WorkItems = epicPublishCrashStore{d.WorkItems}
	got, err := CreateEpic(crashing, in)
	if err == nil || got.Effects != (EpicCreateEffects{Worktree: "completed", Config: "completed"}) {
		t.Fatalf("%+v %v", got, err)
	}
	if epics, err := ListEpics(d); err != nil || len(epics.Epics) != 1 {
		t.Fatalf("pending create broke ordinary use: %+v %v", epics, err)
	}
	got, err = CreateEpic(d, in)
	if err != nil || got.Outcome != "created" || got.Effects != (EpicCreateEffects{Worktree: "completed", Config: "completed"}) {
		t.Fatalf("%+v %v", got, err)
	}
	if defs := readConfigFixture(t, root); !reflect.DeepEqual(defs, want) {
		t.Fatalf("duplicate or lost definition %+v", defs)
	}
	if epics, err := ListEpics(d); err != nil || len(epics.Epics) != 2 {
		t.Fatalf("%+v %v", epics, err)
	}
}

type epicCreateConfigFaultGit struct {
	WorkItemGit
	fault func()
}

func (g *epicCreateConfigFaultGit) CreateBranchAndWorktree(in GitCreateInput) GitCommandOutcome {
	out := g.WorkItemGit.CreateBranchAndWorktree(in)
	if g.fault != nil {
		g.fault()
		g.fault = nil
	}
	return out
}

func TestEpicCreateWithActiveConfigReportsConfigEffect(t *testing.T) {
	for _, test := range []struct{ name, config string }{{"concurrent", "not_attempted"}, {"unwritable", "no_effect"}} {
		t.Run(test.name, func(t *testing.T) {
			root, d, in := epicCreateConfigFixture(t)
			path := filepath.Join(root, DefinitionsFile)
			want := expectedConfigWithNewEpic(t, root)
			changed := epicCreateConfig + "\n"
			git := &epicCreateConfigFaultGit{WorkItemGit: d.WorkGit}
			git.fault = func() {
				if test.name == "concurrent" {
					_ = os.WriteFile(path, []byte(changed), 0644)
				} else {
					_ = os.Chmod(root, 0555)
				}
			}
			t.Cleanup(func() { _ = os.Chmod(root, 0755) })
			d.WorkGit = git
			got, err := CreateEpic(d, in)
			if err == nil || !strings.Contains(err.Error(), path) || got.Effects != (EpicCreateEffects{Worktree: "completed", Config: test.config}) || got.Outcome != "pending" {
				t.Fatalf("%+v %v", got, err)
			}
			raw, _ := os.ReadFile(path)
			if test.name == "concurrent" && string(raw) != changed || test.name != "concurrent" && string(raw) != epicCreateConfig {
				t.Fatal("config overwritten or changed by a failed write")
			}
			if err = os.Chmod(root, 0755); err != nil {
				t.Fatal(err)
			}
			got, err = CreateEpic(d, in)
			if err != nil || got.Outcome != "created" {
				t.Fatalf("retry %+v %v", got, err)
			}
			if defs := readConfigFixture(t, root); !reflect.DeepEqual(defs, want) {
				t.Fatalf("config %+v", defs)
			}
		})
	}
}

func TestEpicCreateWithActiveConfigRejectsWithPath(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*EpicCreateInput)
		reason string
	}{
		{"defined-epic", func(in *EpicCreateInput) { in.EpicID = "feature" }, "already defined differently"},
		{"undefined-repo", func(in *EpicCreateInput) { in.RepoID = "missing" }, "is not defined"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root, d, in := epicCreateConfigFixture(t)
			path := filepath.Join(root, DefinitionsFile)
			test.change(&in)
			got, err := CreateEpic(d, in)
			if err == nil || !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), test.reason) || got.Effects != (EpicCreateEffects{Worktree: "not_attempted", Config: "not_attempted"}) {
				t.Fatalf("%+v %v", got, err)
			}
			if raw, _ := os.ReadFile(path); string(raw) != epicCreateConfig {
				t.Fatal("rejection changed config")
			}
			if _, err = os.Lstat(in.Path); !os.IsNotExist(err) {
				t.Fatal("rejection created worktree")
			}
		})
	}
}

func TestEpicCreateConcurrentDefinitionWritersLoseNoCompletedDefinition(t *testing.T) {
	root, _, in := epicCreateConfigFixture(t)
	_, basis, err := readDefinitions(root)
	if err != nil {
		t.Fatal(err)
	}
	other := in
	other.EpicID, other.Branch, other.Path = "other-epic", "other-epic", filepath.Join(root, "other-epic")
	type result struct {
		effect string
		err    error
	}
	inputs := []EpicCreateInput{in, other}
	results := make([]result, len(inputs))
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range inputs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			effect, err := recordEpicDefinition(root, basis, inputs[i])
			results[i] = result{effect, err}
		}(i)
	}
	close(start)
	wg.Wait()
	// One writer publishes from the shared basis; the other sees it changed.
	completed, lost := results[0], results[1]
	if completed.effect != "completed" {
		completed, lost = lost, completed
	}
	if completed.effect != "completed" || completed.err != nil || lost.effect != "not_attempted" {
		t.Fatalf("want one completed and one not_attempted: %+v", results)
	}
	assertWorkErrorClass(t, lost.err, ErrorWorkIdentityConflict)
	defs := readConfigFixture(t, root)
	for i, r := range results {
		defined := false
		for _, e := range defs.Projects[0].Epics {
			defined = defined || e.ID == string(inputs[i].EpicID)
		}
		if r.effect == "completed" && !defined {
			t.Fatalf("%s reported completed but is absent: %v", inputs[i].EpicID, r.err)
		}
	}
}

type secondRoundFailGit struct {
	WorkItemGit
	calls int
}

func (g *secondRoundFailGit) ValidateBranch(branch string) error {
	if g.calls++; g.calls == 2 {
		return errors.New("second round failed early")
	}
	return g.WorkItemGit.ValidateBranch(branch)
}

func TestEpicCreateKeepsObservedEffectsWhenFinalRoundFailsEarly(t *testing.T) {
	_, d, in := epicCreateConfigFixture(t)
	d.WorkGit = &secondRoundFailGit{WorkItemGit: d.WorkGit}
	got, err := CreateEpic(d, in)
	if err == nil || got.Operation.ID == "" || got.Outcome != "pending" || got.Effects != (EpicCreateEffects{Worktree: "completed", Config: "completed"}) || !strings.Contains(got.NextAction, "second round failed early") {
		t.Fatalf("%+v %v", got, err)
	}
}

type configChangingLocks struct {
	ProjectSnapshotLocker
	change func()
}

func (l *configChangingLocks) WithSnapshotLock(root string, fn func(ProjectSnapshot) error) error {
	if l.change != nil {
		l.change()
		l.change = nil
	}
	return l.ProjectSnapshotLocker.WithSnapshotLock(root, fn)
}

func TestEpicCreateRecordsAttemptWhenConfigChangedBeforeRuntimeLock(t *testing.T) {
	root, d, in := epicCreateConfigFixture(t)
	path := filepath.Join(root, DefinitionsFile)
	d.ProjectLocks = &configChangingLocks{ProjectSnapshotLocker: d.ProjectLocks, change: func() { _ = os.WriteFile(path, []byte(epicCreateConfig+"\n"), 0644) }}
	if _, err := CreateEpic(d, in); err == nil || !strings.Contains(err.Error(), path) {
		t.Fatal(err)
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if op, _ := findEpicCreate(r, in.EpicID); op == nil || len(op.Attempts) != 1 || op.Observations[len(op.Observations)-1].Classification != "exact_effect" {
		t.Fatalf("attempt or observation lost: %+v", op)
	}
}
