package taskwork

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
)

// Fixed inputs deliberately exercise every extracted binding field, both null
// and bound Worktree, and nil versus explicit empty repository slices.
func wi01Models() ([]EpicRecord, []TaskRecord, WorktreeOperationRecord) {
	epic := EpicRecord{ID: "epic", Title: "Epic", ProjectID: "project", RepoBindings: []EpicRepoBinding{{RepoID: "repo", GitCommonDir: "/wi01/main/.git", Worktree: EpicWorktreeBinding{ID: "wt_11111111111111111111111111111111", OwnerKind: "epic", OwnerID: "epic", Origin: "adopted", Locator: "/wi01/epic", Ref: "refs/heads/epic", OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), GitCommonDir: "/wi01/main/.git"}}}}
	task := TaskRecord{ID: "task", Title: "Task", Description: "Description", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo", GitCommonDir: "/wi01/main/.git", WorktreeState: WorkItemUnbound}
	operation := WorktreeOperationRecord{ID: "wop_33333333333333333333333333333333", TaskID: "task", WorktreeID: "wt_22222222222222222222222222222222", State: "creating", ProjectID: task.ProjectID, RepoID: task.RepoID, GitCommonDir: task.GitCommonDir, ParentEpicID: epic.ID, ParentWorktreeID: epic.RepoBindings[0].Worktree.ID, ParentRef: "refs/heads/epic", ParentOID: strings.Repeat("a", 40), ParentTree: strings.Repeat("b", 40), SourceRef: "refs/heads/task", TargetLocator: "/wi01/task", LastObservation: emptyCreateObservation("not_started")}
	operation.IntentDigest, _ = intentDigest(operation)
	bound := task
	bound.ID, bound.WorktreeState = "bound", WorkItemReady
	bound.Worktree = &TaskWorktreeBinding{ID: operation.WorktreeID, OwnerKind: "task", OwnerID: bound.ID, Origin: "ply_created", Locator: operation.TargetLocator, Ref: operation.SourceRef, OID: operation.ParentOID, Tree: operation.ParentTree, GitCommonDir: operation.GitCommonDir, ParentEpicID: operation.ParentEpicID, ParentWorktreeID: operation.ParentWorktreeID, ParentRef: operation.ParentRef, ParentOID: operation.ParentOID, ParentTree: operation.ParentTree, OperationID: operation.ID, IntentDigest: operation.IntentDigest}
	return []EpicRecord{epic, {ID: "nil", RepoBindings: nil}, {ID: "empty", RepoBindings: []EpicRepoBinding{}}}, []TaskRecord{task, bound}, operation
}

func wi01Registry() WorkItemRegistry {
	epics, tasks, _ := wi01Models()
	r := emptyWorkItemRegistry()
	r.Epics, r.Tasks = epics[:1], tasks[:1]
	r.EpicBaseVersions = []EpicBaseVersion{adoptedEpicBase(epics[0], epics[0].RepoBindings[0])}
	return r
}

func wi01Wire(t *testing.T) map[string][]byte {
	t.Helper()
	epics, tasks, op := wi01Models()
	model := struct {
		Epics []EpicRecord `json:"epics" yaml:"epics"`
		Tasks []TaskRecord `json:"tasks" yaml:"tasks"`
	}{epics, tasks}
	j, err := json.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	y, err := yaml.Marshal(model)
	if err != nil {
		t.Fatal(err)
	}
	result := map[string][]byte{"models.json": append(j, '\n'), "models.yaml": y, "intent.txt": []byte(op.IntentDigest + "\n")}
	b, err := encodeWorkItemRegistry(wi01Registry())
	if err != nil {
		t.Fatal(err)
	}
	result["registry.yaml"] = b
	return result
}

func TestWI01CompatibilityWire(t *testing.T) {
	t.Parallel()
	for name, got := range wi01Wire(t) {
		want, err := os.ReadFile(filepath.Join("testdata", "wi01", name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("%s changed:\n%s\nwant:\n%s", name, got, want)
		}
	}
	b, err := os.ReadFile(filepath.Join("testdata", "wi01", "registry.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	r, err := decodeWorkItemRegistry(b)
	if err != nil {
		t.Fatal(err)
	}
	got, err := encodeWorkItemRegistry(r)
	if err != nil || string(got) != string(b) {
		t.Fatalf("reencoding changed: %s %v", got, err)
	}
	// The store's digest binds the original bytes, including insignificant comments.
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(testWork(t, root).Registry()), 0700); err != nil {
		t.Fatal(err)
	}
	raw := append([]byte("# WI01 original bytes\n"), physicalRegistry(b)...)
	if err := os.WriteFile(testWork(t, root).Registry(), raw, 0600); err != nil {
		t.Fatal(err)
	}
	stored, err := newSystemWorkItemStore().Snapshot(testWork(t, root))
	if err != nil {
		t.Fatal(err)
	}
	hash, err := os.ReadFile(filepath.Join("testdata", "wi01", "raw.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if stored.RawSHA256 != strings.TrimSpace(string(hash)) {
		t.Fatalf("raw digest = %s", stored.RawSHA256)
	}
}

func TestWI01CompatibilityValidationOrderAndErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, want string
		mutate     func(*WorkItemRegistry)
	}{
		{"parent before repo and state", "Task task has invalid Epic/Project reference", func(r *WorkItemRegistry) {
			r.Tasks[0].ParentEpicID = "missing"
			r.Tasks[0].RepoID = "missing"
			r.Tasks[0].WorktreeState = "invalid"
		}},
		{"empty missing parent", "Task task has invalid Epic/Project reference", func(r *WorkItemRegistry) {
			r.Tasks[0].ParentEpicID = ""
			r.Tasks[0].ProjectID = ""
		}},
		{"project before repo", "Task task has invalid Epic/Project reference", func(r *WorkItemRegistry) { r.Tasks[0].ProjectID = "other"; r.Tasks[0].RepoID = "missing" }},
		{"repo before state", "Task task has invalid repository binding", func(r *WorkItemRegistry) { r.Tasks[0].RepoID = "missing"; r.Tasks[0].WorktreeState = "invalid" }},
		{"common dir", "Task task has invalid repository binding", func(r *WorkItemRegistry) { r.Tasks[0].GitCommonDir = "/other/.git" }},
		{"operation before parent", "invalid or duplicate operation ID invalid", func(r *WorkItemRegistry) {
			r.WorktreeOperations = []WorktreeOperationRecord{{ID: "invalid"}}
			r.Tasks[0].ParentEpicID = "missing"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := wi01Registry()
			tc.mutate(&r)
			if err := validateWorkItemRegistry(r); err == nil || err.Error() != tc.want {
				t.Fatalf("validation = %v, want %q", err, tc.want)
			}
			// Marshal bypasses validation to test the real decode/store error boundary.
			raw, err := yaml.Marshal(r)
			if err != nil {
				t.Fatal(err)
			}
			_, err = decodeWorkItemRegistry(raw)
			var diagnostic *RegistryDiagnostic
			if !errors.As(err, &diagnostic) || diagnostic.Code != "registry_invalid_content" || diagnostic.Detail != tc.want || diagnostic.FormatVersion != RegistryFormat || errors.Unwrap(err) != nil {
				t.Fatalf("decode error = %#v", err)
			}
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Dir(testWork(t, root).Registry()), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(testWork(t, root).Registry(), physicalRegistry(raw), 0600); err != nil {
				t.Fatal(err)
			}
			_, err = newSystemWorkItemStore().Snapshot(testWork(t, root))
			var work *workspace.WorkItemError
			var wrapped *RegistryDiagnostic
			if !errors.As(err, &work) || work.Class != workspace.ErrorWorkStoreConflict || !errors.As(err, &wrapped) || errors.Unwrap(work) != wrapped || wrapped.Detail != tc.want {
				t.Fatalf("store error = %#v", err)
			}
			want := "workspace_work_store_conflict: invalid work-item registry " + testWork(t, root).Registry() + ": " + diagnostic.Error()
			if err.Error() != want {
				t.Fatalf("store error = %s, want %s", err, want)
			}
		})
	}
	_, err := newSystemWorkItemStore().decodeStoredRegistry(testWork(t, t.TempDir()), []byte("format_version: 999\nepics: []\ntasks: []\nworktree_operations: []\n"))
	var d *RegistryDiagnostic
	if !errors.As(err, &d) || d.Code != "registry_unsupported_format" || d.Detail != "declared format 999 is not format 11" || d.FormatVersion != 999 || errors.Unwrap(err) != nil {
		t.Fatalf("format error = %#v", err)
	}
}
