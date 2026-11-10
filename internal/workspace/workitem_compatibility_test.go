package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func wi01Registry(version int) WorkItemRegistry {
	epics, tasks, _ := wi01Models()
	r := WorkItemRegistry{FormatVersion: version, Epics: epics[:1], Tasks: tasks[:1], WorktreeOperations: []WorktreeOperationRecord{}}
	if version == 2 {
		r.TaskResults = []HistoricalRecord{}
		r.HumanQARecords = []HistoricalRecord{}
		r.IntegrationAuthorities = []IntegrationAuthority{}
		r.IntegrationIntents = []IntegrationIntent{}
		r.IntegrationAttempts = []IntegrationAttempt{}
		r.IntegrationResults = []IntegrationResult{}
	}
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
	for _, version := range []int{1, 2} {
		b, err := encodeWorkItemRegistry(wi01Registry(version))
		if err != nil {
			t.Fatal(err)
		}
		result[fmt.Sprintf("registry-v%d.yaml", version)] = b
	}
	return result
}

func TestWI01CompatibilityWire(t *testing.T) {
	for name, got := range wi01Wire(t) {
		want, err := os.ReadFile(filepath.Join("testdata", "wi01", name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("%s changed:\n%s\nwant:\n%s", name, got, want)
		}
	}
	for _, version := range []int{1, 2} {
		b, err := os.ReadFile(filepath.Join("testdata", "wi01", fmt.Sprintf("registry-v%d.yaml", version)))
		if err != nil {
			t.Fatal(err)
		}
		r, err := decodeWorkItemRegistry(b)
		if err != nil {
			t.Fatal(err)
		}
		got, err := encodeWorkItemRegistry(r)
		if err != nil || string(got) != string(b) {
			t.Fatalf("v%d reencoding changed: %s %v", version, got, err)
		}
		// The store's digest binds the original bytes, including insignificant comments.
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Dir(workItemsPath(root)), 0700); err != nil {
			t.Fatal(err)
		}
		raw := append([]byte("# WI01 original bytes\n"), b...)
		if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
			t.Fatal(err)
		}
		stored, err := newSystemWorkItemStore().Snapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		hash, err := os.ReadFile(filepath.Join("testdata", "wi01", fmt.Sprintf("raw-v%d.txt", version)))
		if err != nil {
			t.Fatal(err)
		}
		if stored.RawSHA256 != strings.TrimSpace(string(hash)) {
			t.Fatalf("v%d raw digest = %s", version, stored.RawSHA256)
		}
	}
}

func TestWI01CompatibilityValidationOrderAndErrors(t *testing.T) {
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
			r := wi01Registry(1)
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
			if !errors.As(err, &diagnostic) || diagnostic.Code != "registry_invalid_content" || diagnostic.Detail != tc.want || diagnostic.FormatVersion != 1 || errors.Unwrap(err) != nil {
				t.Fatalf("decode error = %#v", err)
			}
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Dir(workItemsPath(root)), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(workItemsPath(root), raw, 0600); err != nil {
				t.Fatal(err)
			}
			_, err = newSystemWorkItemStore().Snapshot(root)
			var work *WorkItemError
			var wrapped *RegistryDiagnostic
			if !errors.As(err, &work) || work.Class != ErrorWorkStoreConflict || !errors.As(err, &wrapped) || errors.Unwrap(work) != wrapped || wrapped.Detail != tc.want {
				t.Fatalf("store error = %#v", err)
			}
			want := "workspace_work_store_conflict: invalid work-item registry " + workItemsPath(root) + ": " + diagnostic.Error()
			if err.Error() != want {
				t.Fatalf("store error = %s, want %s", err, want)
			}
		})
	}
	_, err := decodeWorkItemRegistry([]byte("format_version: 999\nepics: []\ntasks: []\nworktree_operations: []\n"))
	var d *RegistryDiagnostic
	if !errors.As(err, &d) || d.Code != "registry_unsupported_format" || d.Detail != "declared format 999 exceeds this reader" || d.FormatVersion != 999 || errors.Unwrap(err) != nil {
		t.Fatalf("format error = %#v", err)
	}
}
