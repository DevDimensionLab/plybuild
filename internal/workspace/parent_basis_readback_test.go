package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func advancedReadbackFixture(t *testing.T) (workItemJourneyFixture, string, string) {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	in, err := ParseTaskWorktreeCreateInput("task", "old-task", filepath.Join(f.wrapper, "old-task"), f.oid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	epic := filepath.Join(f.wrapper, "epic")
	if err = os.WriteFile(filepath.Join(epic, "advance.txt"), []byte("B\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, epic, "add", "advance.txt")
	runLocalGit(t, epic, "commit", "-m", "advance base")
	oid := runLocalGit(t, epic, "rev-parse", "HEAD")
	tree := runLocalGit(t, epic, "rev-parse", "HEAD^{tree}")
	preview, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil {
		t.Fatal(err)
	}
	upgrade := ""
	if preview.Plan.RegistryUpgrade != nil {
		upgrade = preview.Plan.RegistryUpgrade.RegistrySHA256
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Upgrade: upgrade, Apply: true, Confirmation: *preview.Confirmation}); err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTask(f.dependencies, TaskCreateInput{TaskID: "new-task", Title: "New", Description: "New task", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	return f, oid, tree
}

func assertNullableParent(t *testing.T, r TaskReadbackResult, oid, tree string) {
	t.Helper()
	if r.Integration == nil {
		t.Fatal("missing nullable integration")
	}
	b, err := MarshalTaskIntegrationReadback(*r.Integration)
	if err != nil {
		t.Fatal(err)
	}
	var v struct {
		Epic IntegrationPlanEpic `json:"epic"`
		Plan any                 `json:"plan"`
	}
	if err = json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	// Every bound Task gets the same fresh integration check; only an unbound
	// Task has no plan.
	if (v.Plan != nil) != (r.Task.Worktree != nil) || v.Epic.ExpectedParentOID != oid || v.Epic.ExpectedParentTree != tree {
		t.Errorf("nullable expected parent = %+v, want %s/%s; plan=%v", v.Epic, oid, tree, v.Plan)
	}
}

func TestParentBasisReadbackUsesFrozenTaskAndCurrentUnboundBase(t *testing.T) {
	f, b, tree := advancedReadbackFixture(t)
	unbound, err := ShowTask(f.dependencies, "new-task")
	if err != nil {
		t.Fatal(err)
	}
	if unbound.ParentFreshness.OID != "fresh" || unbound.ParentFreshness.Tree != "fresh" || len(unbound.Reasons) != 1 || !containsReason(unbound.Reasons, "task_worktree_unbound") || unbound.ReadyForHandoff {
		t.Errorf("unbound parent freshness=%+v reasons=%v ready=%v", unbound.ParentFreshness, unbound.Reasons, unbound.ReadyForHandoff)
	}
	assertNullableParent(t, unbound, b, tree)
	in, err := ParseTaskWorktreeCreateInput("new-task", "new-task", filepath.Join(f.wrapper, "new-task"), b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	fresh, err := ShowTask(f.dependencies, "new-task")
	if err != nil {
		t.Fatal(err)
	}
	if fresh.ParentFreshness.OID != "fresh" || fresh.ParentFreshness.Tree != "fresh" || len(fresh.Reasons) != 0 || !fresh.ReadyForHandoff {
		t.Errorf("new bound freshness=%+v reasons=%v ready=%v", fresh.ParentFreshness, fresh.Reasons, fresh.ReadyForHandoff)
	}
	assertNullableParent(t, fresh, b, tree)
	old, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if old.ParentFreshness.OID != "stale" || old.ParentFreshness.Tree != "stale" || !containsReason(old.Reasons, "parent_ref_stale") || old.ReadyForHandoff {
		t.Errorf("old frozen parent freshness=%+v reasons=%v ready=%v", old.ParentFreshness, old.Reasons, old.ReadyForHandoff)
	}
	assertNullableParent(t, old, f.oid, old.Task.Worktree.ParentTree)
	// Rewind only the private fixture, never the candidate or a live Epic.
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "reset", "--hard", f.oid)
	moved, err := ShowTask(f.dependencies, "new-task")
	if err != nil {
		t.Fatal(err)
	}
	if moved.ParentFreshness.OID != "stale" || moved.ParentFreshness.Tree != "stale" || !containsReason(moved.Reasons, "parent_ref_stale") || moved.ReadyForHandoff {
		t.Errorf("rewound parent freshness=%+v reasons=%v ready=%v", moved.ParentFreshness, moved.Reasons, moved.ReadyForHandoff)
	}
	assertNullableParent(t, moved, b, tree)
}

func TestParentBasisReadbackEpicUsesOneCurrentObservation(t *testing.T) {
	f, b, tree := advancedReadbackFixture(t)
	show, err := ShowEpic(f.dependencies, "epic")
	if err != nil {
		t.Fatal(err)
	}
	if !show.Ready || len(show.Reasons) != 0 || show.Freshness[0].OID != "fresh" || show.Freshness[0].Tree != "fresh" || show.Base[0].Freshness != "fresh" || show.Base[0].OID != b || show.Base[0].Tree != tree {
		t.Errorf("current Epic = %+v", show)
	}
	if show.Epic.RepoBindings[0].Worktree.OID != f.oid {
		t.Error("historical adoption anchor was replaced")
	}
	_, repos, err := f.dependencies.Projects.Snapshot(filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	f.dependencies.Projects = projectSnapshotOverride{ProjectStore: f.dependencies.Projects, repositories: repos}
	show, err = ShowEpic(f.dependencies, "epic")
	if err != nil {
		t.Fatal(err)
	}
	if show.Ready || !containsReason(show.Reasons, "project_binding_stale") || show.Base[0].Freshness != "unknown" {
		t.Errorf("project conflict Epic = %+v", show)
	}
}

func TestParentBasisReadbackPendingIntentRetainsHistoricalBase(t *testing.T) {
	for _, reconcile := range []bool{false, true} {
		t.Run(map[bool]string{false: "creating", true: "reconciliation"}[reconcile], func(t *testing.T) {
			f, b, tree := advancedReadbackFixture(t)
			baseGit := f.dependencies.WorkGit
			if reconcile {
				f.dependencies.WorkGit = wrongBranchGit{baseGit}
			} else {
				f.dependencies.WorkGit = noEffectGit{baseGit}
			}
			in, err := ParseTaskWorktreeCreateInput("new-task", "new-task", filepath.Join(f.wrapper, "new-task"), b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = CreateTaskWorktree(f.dependencies, in); err == nil {
				t.Fatal("fault did not preserve intent")
			}
			f.dependencies.WorkGit = baseGit
			// Move only the private parent after the intent froze B; pending intents block base publication.
			epic := filepath.Join(f.wrapper, "epic")
			runLocalGit(t, epic, "commit", "--allow-empty", "-m", "base C")
			show, err := ShowTask(f.dependencies, "new-task")
			if err != nil {
				t.Fatal(err)
			}
			reason := "task_worktree_creation_in_progress"
			if reconcile {
				reason = "task_worktree_reconciliation_required"
			}
			if show.Operation == nil || show.Operation.ParentOID != b || show.ParentFreshness.OID != "stale" || !containsReason(show.Reasons, "parent_ref_stale") || !containsReason(show.Reasons, reason) || show.ReadyForHandoff {
				t.Fatalf("pending frozen parent = %+v", show)
			}
			assertNullableParent(t, show, b, tree)
			runLocalGit(t, epic, "reset", "--hard", b)
			show, err = ShowTask(f.dependencies, "new-task")
			if err != nil {
				t.Fatal(err)
			}
			if show.ParentFreshness.OID != "fresh" || show.ParentFreshness.Tree != "fresh" || containsReason(show.Reasons, "parent_ref_stale") || show.ReadyForHandoff {
				t.Fatalf("pending B at B freshness=%+v reasons=%v", show.ParentFreshness, show.Reasons)
			}
		})
	}
}

type parentReadbackGit struct {
	WorkItemGit
	mode  string
	calls int
}

func (g *parentReadbackGit) ObserveWorktree(path string) (GitWorktreeObservation, error) {
	g.calls++
	o, err := g.WorkItemGit.ObserveWorktree(path)
	if err != nil {
		return o, err
	}
	switch g.mode {
	case "ref":
		o.Ref += "-foreign"
	case "locator":
		o.Locator += "-foreign"
	case "common":
		o.GitCommonDir += "-foreign"
	case "oid":
		o.OID = "foreign"
	case "tree":
		o.Tree = "foreign"
	case "dirty":
		o.Clean = false
	case "inventory":
		o.InventoryMatch = false
	case "unknown":
		return GitWorktreeObservation{}, os.ErrNotExist
	}
	return o, nil
}

func TestParentBasisReadbackPreservesObservationGuards(t *testing.T) {
	f, _, _ := advancedReadbackFixture(t)
	for _, tt := range []struct{ mode, epicReason, parentReason, field string }{
		{"ref", "epic_ref_stale", "parent_ref_stale", "Ref"},
		{"oid", "epic_ref_stale", "parent_ref_stale", "OID"},
		{"tree", "epic_ref_stale", "parent_ref_stale", "Tree"},
		{"locator", "epic_worktree_stale", "parent_worktree_stale", "Locator"},
		{"common", "epic_worktree_stale", "parent_worktree_stale", "GitCommonDir"},
		{"dirty", "epic_worktree_dirty", "parent_worktree_dirty", "Clean"},
		{"inventory", "worktree_inventory_stale", "worktree_inventory_stale", "InventoryMatch"},
		{"unknown", "epic_observation_unknown", "parent_observation_unknown", ""},
		{"nil", "epic_observation_unknown", "parent_observation_unknown", ""},
	} {
		t.Run(tt.mode, func(t *testing.T) {
			d := f.dependencies
			git := &parentReadbackGit{WorkItemGit: d.WorkGit, mode: tt.mode}
			d.WorkGit = git
			if tt.mode == "nil" {
				d.WorkGit = nil
			}
			show, err := ShowEpic(d, "epic")
			if err != nil {
				t.Fatal(err)
			}
			expected := "stale"
			if tt.field == "" {
				expected = "unknown"
			}
			if show.Ready || !containsReason(show.Reasons, tt.epicReason) || show.Base[0].Freshness != expected {
				t.Errorf("guarded Epic = %+v", show)
			}
			if tt.mode != "nil" && git.calls != 1 {
				t.Errorf("Epic observed %d times, want one coherent observation", git.calls)
			}
			task, err := ShowTask(d, "new-task")
			if err != nil {
				t.Fatal(err)
			}
			if task.ReadyForHandoff || !containsReason(task.Reasons, tt.parentReason) {
				t.Errorf("guarded Task reasons=%v ready=%v", task.Reasons, task.ReadyForHandoff)
			}
			if tt.field == "" && (task.Parent.OID != nil || task.ParentFreshness.OID != "unknown") {
				t.Errorf("unknown parent became observed: %+v", task.ParentFreshness)
			}
		})
	}
}

type parentReadbackStore struct {
	WorkItemStore
	mutate func(*WorkItemRegistry)
}

func (s parentReadbackStore) Snapshot(root string) (WorkItemRegistry, error) {
	r, err := s.WorkItemStore.Snapshot(root)
	if err == nil {
		s.mutate(&r)
	}
	return r, err
}

func TestParentBasisReadbackRejectsUnknownHistoricalParent(t *testing.T) {
	f, _, _ := advancedReadbackFixture(t)
	for _, operation := range []bool{false, true} {
		d := f.dependencies
		d.WorkItems = parentReadbackStore{d.WorkItems, func(r *WorkItemRegistry) {
			// A malformed dependency must not silently fall back to the current base.
			if operation {
				task, _ := findTask(*r, "new-task")
				task.WorktreeState = WorkItemCreating
				r.WorktreeOperations = append(r.WorktreeOperations, WorktreeOperationRecord{TaskID: "new-task", ParentOID: "unknown", ParentTree: "unknown"})
			} else {
				task, _ := findTask(*r, "new-task")
				task.Worktree = &TaskWorktreeBinding{ParentOID: "unknown", ParentTree: "unknown"}
			}
		}}
		if _, err := ShowTask(d, "new-task"); err == nil {
			t.Fatalf("unknown history accepted (operation=%v)", operation)
		}
	}
}

func TestParentBasisReadbackLegacyFallback(t *testing.T) {
	f := newWorkItemJourneyFixtureVersion(t, true)
	show, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if show.ParentFreshness.OID != "fresh" || containsReason(show.Reasons, "parent_ref_stale") {
		t.Fatalf("legacy unbound = %+v", show)
	}
	in, err := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	show, err = ShowTask(f.dependencies, "task")
	if err != nil || !show.ReadyForHandoff {
		t.Fatalf("legacy bound = %+v, %v", show, err)
	}
	epic, err := ShowEpic(f.dependencies, "epic")
	if err != nil || !epic.Ready || len(epic.Base) != 0 || len(epic.Reasons) != 0 {
		t.Fatalf("legacy Epic = %+v, %v", epic, err)
	}
}
