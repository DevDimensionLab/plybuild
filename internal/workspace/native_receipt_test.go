package workspace

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestNativeReceiptSingleCompletionSource(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	ready, err := CreateTaskWorktree(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.WorktreeOperations) != 0 {
		t.Fatalf("native completion retained %d operation duplicates", len(r.WorktreeOperations))
	}
	for _, state := range []WorkItemState{WorkItemReady, WorkItemRetired} {
		r.Tasks[0].WorktreeState = state
		if err := validateWorkItemRegistry(r); err != nil {
			t.Fatalf("%s: %v", state, err)
		}
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil || !show.WorktreeReady || !show.ReadyForHandoff || show.Operation != nil {
		t.Fatalf("readback: %#v, %v", show, err)
	}
	before, _ := os.ReadFile(workItemsPath(root))
	// An exact replay needs neither a Git reader nor a Git writer.
	f.dependencies.WorkGit = receiptReplayGit{}
	retry, err := CreateTaskWorktree(f.dependencies, in)
	if err != nil || retry.Outcome != "already_ready" || retry.Task.Worktree.ID != ready.Task.Worktree.ID {
		t.Fatalf("replay: %#v, %v", retry, err)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if string(before) != string(after) {
		t.Fatal("replay published registry")
	}
}

type receiptReplayGit struct{ WorkItemGit }

func (receiptReplayGit) ValidateBranch(string) error { return nil }

func TestNativeReceiptRejectsInvalidBindings(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	mutations := map[string]func(*WorkItemRegistry){
		"unknown version": func(r *WorkItemRegistry) { r.Tasks[0].Worktree.NativeReceiptVersion = 2 },
		"adoption":        func(r *WorkItemRegistry) { r.Tasks[0].Worktree.Origin = "adopted" },
		"creating":        func(r *WorkItemRegistry) { r.Tasks[0].WorktreeState = WorkItemCreating },
		"reconciliation":  func(r *WorkItemRegistry) { r.Tasks[0].WorktreeState = WorkItemReconciliationRequired },
		"missing marker":  func(r *WorkItemRegistry) { r.Tasks[0].Worktree.NativeReceiptVersion = 0 },
		"common dir":      func(r *WorkItemRegistry) { r.Tasks[0].Worktree.GitCommonDir += "-wrong" },
		"ref":             func(r *WorkItemRegistry) { r.Tasks[0].Worktree.Ref += "-wrong" },
		"path":            func(r *WorkItemRegistry) { r.Tasks[0].Worktree.Locator += "-wrong" },
		"start":           func(r *WorkItemRegistry) { r.Tasks[0].Worktree.OID = strings.Repeat("e", 40) },
		"parent":          func(r *WorkItemRegistry) { r.Tasks[0].Worktree.ParentOID = strings.Repeat("e", 40) },
		"parent ref":      func(r *WorkItemRegistry) { r.Tasks[0].Worktree.ParentRef += "-wrong" },
		"digest":          func(r *WorkItemRegistry) { r.Tasks[0].Worktree.IntentDigest = "sha256:" + strings.Repeat("e", 64) },
		"worktree id":     func(r *WorkItemRegistry) { r.Tasks[0].Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID },
		"ready duplicate": func(r *WorkItemRegistry) {
			// Legacy fixture has an actual, valid ready operation; only its Task identity
			// is rebound so both completion sources describe this native Task.
			w := r.Tasks[0].Worktree
			ep := r.Epics[0].RepoBindings[0].Worktree
			repo := RepoRecord{Locator: filepath.Join(f.wrapper, "main"), GitCommonDir: w.GitCommonDir}
			op := WorktreeOperationRecord{ID: w.OperationID, TaskID: r.Tasks[0].ID, WorktreeID: w.ID, IntentDigest: w.IntentDigest, State: "ready", ProjectID: r.Tasks[0].ProjectID, RepoID: r.Tasks[0].RepoID, GitCommonDir: w.GitCommonDir, ParentEpicID: w.ParentEpicID, ParentWorktreeID: w.ParentWorktreeID, ParentRef: w.ParentRef, ParentOID: w.ParentOID, ParentTree: w.ParentTree, SourceRef: w.Ref, TargetLocator: w.Locator}
			_, op.LastObservation = observeTaskCreate(f.dependencies, repo, EpicRepoBinding{GitCommonDir: w.GitCommonDir, Worktree: ep}, op)
			r.WorktreeOperations = []WorktreeOperationRecord{op}
		},
		"resource conflict": func(r *WorkItemRegistry) {
			other := r.Tasks[0]
			other.ID = "other"
			w := *other.Worktree
			other.Worktree = &w
			w.OwnerID = other.ID
			w.ID = WorktreeID("wt_" + strings.Repeat("e", 32))
			w.OperationID = WorktreeOperationID("wop_" + strings.Repeat("e", 32))
			w.IntentDigest, _ = nativeIntentDigest(w.GitCommonDir, other.ParentEpicID, w.ParentOID, w.ParentRef, w.ParentTree, w.ParentWorktreeID, other.ProjectID, other.RepoID, w.Ref, w.Locator, other.ID)
			r.Tasks = append([]TaskRecord{other}, r.Tasks...)
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			b, e := encodeWorkItemRegistry(r)
			if e != nil {
				t.Fatal(e)
			}
			copy, e := decodeWorkItemRegistry(b)
			if e != nil {
				t.Fatal(e)
			}
			mutate(&copy)
			if e = validateWorkItemRegistry(copy); e == nil {
				t.Fatal("invalid receipt accepted")
			}
		})
	}
	for _, state := range []WorkItemState{WorkItemReady, WorkItemRetired} {
		r.Tasks[0].WorktreeState = state
		for _, resource := range []string{"ref", "path"} {
			w := r.Tasks[0].Worktree
			ref, path := "refs/heads/free", filepath.Join(f.wrapper, "free")
			if resource == "ref" {
				ref = w.Ref
			} else {
				path = w.Locator
			}
			if nativeReceiptReservation(r, w.GitCommonDir, ref, path) == nil {
				t.Fatalf("%s did not reserve %s", state, resource)
			}
		}
	}
}

func TestNativeReceiptWriterBytesAndFreshReadback(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(workItemsPath(filepath.Dir(f.wrapper)))
	if err != nil {
		t.Fatal(err)
	}
	if dest := os.Getenv("PLY_NATIVE_RECEIPT_WRITER_BYTES"); dest != "" {
		if err = os.WriteFile(dest, b, 0600); err != nil {
			t.Fatal(err)
		}
		// Export the actual physical envelope and its native referenced objects,
		// allowing the old store reader to consume unchanged writer bytes.
		root := filepath.Dir(f.wrapper)
		err = filepath.WalkDir(filepath.Join(root, MarkerDirectory), func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			target := filepath.Join(filepath.Dir(dest), "writer-fixture", rel)
			if entry.IsDir() {
				return os.MkdirAll(target, 0700)
			}
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if err = os.WriteFile(target, bytes, 0600); err != nil {
				return err
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			return os.Chmod(target, info.Mode().Perm())
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if !strings.Contains(string(b), "native_receipt_version: 1") {
		t.Fatal("writer omitted contract marker")
	}
	if err = os.WriteFile(filepath.Join(in.Path, "dirty"), []byte("dirty\n"), 0600); err != nil {
		t.Fatal(err)
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil || show.ReadyForHandoff || !slices.Contains(show.Reasons, "target_worktree_dirty") {
		t.Fatalf("fresh readback: %#v %v", show, err)
	}
	json, err := MarshalTaskReadback(show)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(json), `"operation":null`) || !strings.Contains(string(json), `"native_receipt_version":1`) {
		t.Fatalf("API lacks receipt/null operation: %s", json)
	}
}

func TestNativeReceiptAtomicReadyPublicationAndLostReply(t *testing.T) {
	for _, stage := range []string{"temp-open", "write", "file-sync", "replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			calls := 0
			f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
				if point == stage {
					calls++
					if calls == 2 {
						return errInjected
					}
				}
				return nil
			}}}
			in, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
			if _, err := CreateTaskWorktree(f.dependencies, in); err == nil {
				t.Fatal("ready publication fault not observed")
			}
			f.dependencies.WorkItems = newSystemWorkItemStore()
			root := filepath.Dir(f.wrapper)
			r, err := f.dependencies.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			var id WorktreeID
			if stage == "directory-sync" {
				if r.Tasks[0].Worktree == nil || len(r.WorktreeOperations) != 0 {
					t.Fatal("lost reply left duplicate or incomplete receipt")
				}
				id = r.Tasks[0].Worktree.ID
				f.dependencies.WorkGit = receiptReplayGit{}
			} else {
				if r.Tasks[0].Worktree != nil || len(r.WorktreeOperations) != 1 || r.Tasks[0].WorktreeState != WorkItemCreating {
					t.Fatal("precommit fault lost recovery intent")
				}
				id = r.WorktreeOperations[0].WorktreeID
				// Recovery may observe the effect but cannot execute Git again.
				f.dependencies.WorkGit = receiptNoEffectGit{f.dependencies.WorkGit}
			}
			ready, err := CreateTaskWorktree(f.dependencies, in)
			if err != nil || ready.Task.Worktree.ID != id {
				t.Fatalf("recovery %#v %v", ready, err)
			}
			r, err = f.dependencies.WorkItems.Snapshot(root)
			if err != nil || len(r.WorktreeOperations) != 0 || r.Tasks[0].Worktree == nil {
				t.Fatalf("completion %#v %v", r, err)
			}
		})
	}
}

type receiptNoEffectGit struct{ WorkItemGit }

func (receiptNoEffectGit) CreateBranchAndWorktree(GitCreateInput) GitCommandOutcome {
	panic("recovery started a second Git effect")
}
func (receiptNoEffectGit) AddWorktreeForExistingBranch(GitCreateInput) GitCommandOutcome {
	panic("recovery added another worktree")
}

func TestNativeReceiptReplayAndReservationsAfterCheckoutRemoval(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	other, _ := ParseTaskCreateInput("other", "Other", "Other task", "epic", "ply", "ply")
	if _, err := CreateTask(f.dependencies, other); err != nil {
		t.Fatal(err)
	}
	in, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(f.wrapper, "main")
	runLocalGit(t, repo, "worktree", "remove", in.Path)
	runLocalGit(t, repo, "branch", "-D", in.Branch)
	before, _ := os.ReadFile(workItemsPath(root))
	deps := f.dependencies
	deps.WorkGit = receiptReplayGit{}
	replay, err := CreateTaskWorktree(deps, in)
	if err != nil || replay.Outcome != "already_ready" || replay.Operation != nil {
		t.Fatalf("missing checkout replay: %#v %v", replay, err)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if string(before) != string(after) {
		t.Fatal("replay rewrote registry")
	}
	for _, resource := range []string{"ref", "path"} {
		branch, path := "other", filepath.Join(f.wrapper, "other")
		if resource == "ref" {
			branch = in.Branch
		} else {
			path = in.Path
		}
		conflicting, _ := ParseTaskWorktreeCreateInput("other", branch, path, f.oid)
		if _, err = CreateTaskWorktree(f.dependencies, conflicting); err == nil || !hasWorkErrorClass(err, ErrorWorkIdentityConflict) {
			t.Fatalf("lost %s reservation: %v", resource, err)
		}
	}
}
