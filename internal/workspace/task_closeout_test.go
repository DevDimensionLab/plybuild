package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func closeoutFixture(t *testing.T, integrated, legacy bool) (integrationJourneyFixture, TaskCloseoutInput) {
	t.Helper()
	f := newIntegrationJourneyFixtureVersion(t, legacy)
	b := []byte(`{"fixture":"preserved native evidence"}`)
	sha := digestTaskBytes(b)
	for _, name := range []string{"handoff", "start", "terminal", "owner-release"} {
		if err := os.WriteFile(filepath.Join(f.root, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		tr := findTaskResult(r, f.resultID)
		tr.HandoffSHA256, tr.StartReceiptSHA256, tr.TerminalResultSHA256 = sha, sha, sha
		artifact, err := os.ReadFile(filepath.Join(f.taskPath, "result.txt"))
		if err != nil {
			return err
		}
		tr.Artifacts = []TaskArtifactRecord{{ArtifactID: "review-source", Role: "review", Locator: filepath.Join(f.taskPath, "result.txt"), SHA256: digestTaskBytes(artifact), SizeBytes: int64(len(artifact))}}
		f.evidence.evidence.Artifacts = tr.Artifacts
		f.evidence.evidence.HandoffSHA256, f.evidence.evidence.StartReceiptSHA256, f.evidence.evidence.TerminalResultSHA256 = sha, sha, sha
		for i := range r.TaskResultSpecBindings {
			r.TaskResultSpecBindings[i].HandoffSHA256, r.TaskResultSpecBindings[i].StartReceiptSHA256 = sha, sha
		}
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	in := TaskCloseoutInput{OperationID: "integration-fixture", TaskID: "task", TaskResultID: f.resultID, ExpectedResultOID: f.resultOID, ExpectedResultTree: f.resultTree, TargetRef: "refs/heads/epic", Ownership: TaskCloseoutOwnership{Released: true, TaskID: "task", TaskResultID: f.resultID, ResultOID: f.resultOID, ResultTree: f.resultTree, EvidenceLocator: filepath.Join(f.root, "owner-release"), EvidenceSHA256: sha}}
	if integrated {
		closeoutIntegrateFixture(t, f)
	}
	return f, in
}

func closeoutIntegrateFixture(t *testing.T, f integrationJourneyFixture) {
	t.Helper()
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	p, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	if p.Readback.Classification != "ready" {
		t.Fatalf("integration is not ready: %s %+v", p.Readback.Classification, p.Readback.NextAction)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, p.Readback)
	if _, err := ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
}

func TestTaskCloseoutPreviewAndApplyRequireDifferentIntegrationPhases(t *testing.T) {
	f, in := closeoutFixture(t, false, false)
	registryBefore, _ := os.ReadFile(workItemsPath(f.root))
	preview, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !preview.Ready || preview.IntegrationObserved {
		t.Fatalf("preview = %+v, %v", preview, err)
	}
	if _, err := os.Stat(filepath.Join(f.root, MarkerDirectory, taskCloseoutsFile)); !os.IsNotExist(err) {
		t.Fatal("preview created closeout metadata")
	}
	registryAfter, _ := os.ReadFile(workItemsPath(f.root))
	if string(registryAfter) != string(registryBefore) {
		t.Fatal("preview changed registry")
	}
	if _, err := ApplyTaskCloseout(f.dependencies, in, preview.PlanSHA256); err == nil || !strings.Contains(err.Error(), "integration_not_observed") {
		t.Fatalf("apply without integration = %v", err)
	}
	closeoutIntegrateFixture(t, f)
	previewAfter, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !previewAfter.Ready || !previewAfter.IntegrationObserved || previewAfter.PlanSHA256 != preview.PlanSHA256 {
		t.Fatalf("integration changed approved resource plan: %+v %v", previewAfter, err)
	}
	receipt, err := ApplyTaskCloseout(f.dependencies, in, preview.PlanSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.State != "complete" || receipt.ResourceState != "retired" || !receipt.BranchRemoved || !receipt.WorktreeRemoved || !receipt.LifecycleCompleted {
		t.Fatalf("receipt = %+v", receipt)
	}
	if _, err := os.Stat(f.taskPath); !os.IsNotExist(err) {
		t.Fatal("source still exists")
	}
	for _, name := range []string{"main", "epic"} {
		if _, err := os.Stat(filepath.Join(f.wrapper, name)); err != nil {
			t.Fatal("protected worktree removed", err)
		}
	}
	if _, err := execGitOutput(filepath.Join(f.wrapper, "main"), "show-ref", "--verify", "refs/heads/task"); err == nil {
		t.Fatal("task branch still exists")
	}
	readback, err := ShowTask(f.dependencies, "task")
	if err != nil || readback.Task.WorktreeState != WorkItemRetired || readback.Lifecycle != LifecycleCompleted || readback.Target.Kind != "retired" {
		t.Fatalf("retired readback %+v %v", readback, err)
	}
	for _, reason := range readback.Reasons {
		if strings.Contains(reason, "stale") || strings.Contains(reason, "unknown") {
			t.Fatalf("retired source misclassified: %v", readback.Reasons)
		}
	}
	if receipt.Retention == nil || len(receipt.Retention.Entries) < 7 {
		t.Fatalf("incomplete retention: %+v", receipt.Retention)
	}
	retained, err := ReadTaskRetainedEvidence(f.root, "task", filepath.Join(f.taskPath, "result.txt"), digestTaskBytes([]byte("result\n")))
	if err != nil || string(retained) != "result\n" {
		t.Fatalf("retained source evidence unavailable: %q %v", retained, err)
	}
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "gc", "--prune=now")
	if err := verifyCloseoutRetention(f.dependencies, f.root, receipt.Plan, *receipt.Retention); err != nil {
		t.Fatal(err)
	}
	storeBefore, _ := os.ReadFile(filepath.Join(f.root, MarkerDirectory, taskCloseoutsFile))
	in.OperationID = "second-delivery-same-effect"
	repeated, err := ApplyTaskCloseout(f.dependencies, in, preview.PlanSHA256)
	if err != nil || repeated.OperationID != receipt.OperationID {
		t.Fatalf("repeat: %+v %v", repeated, err)
	}
	storeAfter, _ := os.ReadFile(filepath.Join(f.root, MarkerDirectory, taskCloseoutsFile))
	if string(storeBefore) != string(storeAfter) {
		t.Fatal("completed retry rewrote receipt")
	}
}

func TestTaskCloseoutRejectsUnsafeSources(t *testing.T) {
	for _, scenario := range []string{"tracked", "assume-unchanged", "skip-worktree", "untracked", "ignored", "nested", "locked", "active", "unknown-owner", "inside-source", "symlink", "changed-ref", "symbolic-ref", "missing", "submodule"} {
		t.Run(scenario, func(t *testing.T) {
			f, in := closeoutFixture(t, false, true)
			switch scenario {
			case "tracked":
				os.WriteFile(filepath.Join(f.taskPath, "result.txt"), []byte("new\n"), 0600)
			case "assume-unchanged", "skip-worktree":
				runLocalGit(t, f.taskPath, "update-index", "--"+scenario, "README.md")
				if err := os.WriteFile(filepath.Join(f.taskPath, "README.md"), []byte("hidden changed bytes\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "untracked":
				os.WriteFile(filepath.Join(f.taskPath, "unknown"), []byte("unknown\n"), 0600)
			case "ignored":
				common := runLocalGit(t, f.taskPath, "rev-parse", "--git-common-dir")
				os.WriteFile(filepath.Join(common, "info", "exclude"), []byte("ignored\n"), 0600)
				os.WriteFile(filepath.Join(f.taskPath, "ignored"), []byte("not disposable\n"), 0600)
			case "nested":
				os.Mkdir(filepath.Join(f.taskPath, "nested"), 0700)
				runLocalGit(t, filepath.Join(f.taskPath, "nested"), "init")
			case "locked":
				runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "lock", f.taskPath)
			case "active":
				in.Ownership.Released = false
			case "unknown-owner":
				in.Ownership.EvidenceSHA256 = "sha256:" + strings.Repeat("0", 64)
			case "inside-source":
				f.dependencies.Files = projectCwdFileSystem{FileSystem: f.dependencies.Files, cwd: f.taskPath}
			case "symlink":
				if err := os.Rename(f.taskPath, f.taskPath+"-moved"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(f.taskPath+"-moved", f.taskPath); err != nil {
					t.Fatal(err)
				}
			case "changed-ref":
				runLocalGit(t, f.taskPath, "switch", "-c", "unowned")
			case "symbolic-ref":
				runLocalGit(t, f.taskPath, "branch", "unowned", f.resultOID)
				runLocalGit(t, f.taskPath, "symbolic-ref", "refs/heads/task", "refs/heads/unowned")
			case "missing":
				if err := os.Rename(f.taskPath, f.taskPath+"-moved"); err != nil {
					t.Fatal(err)
				}
			case "submodule":
				runLocalGit(t, f.taskPath, "update-index", "--add", "--cacheinfo", "160000,"+f.resultOID+",submodule")
			}
			p, err := PreviewTaskCloseout(f.dependencies, in)
			if err == nil && p.Ready {
				t.Fatalf("unsafe source accepted: %+v", p)
			}
			if _, err := os.Stat(filepath.Join(f.root, MarkerDirectory, taskCloseoutsFile)); !os.IsNotExist(err) {
				t.Fatal("unsafe preview mutated state")
			}
		})
	}
}

func TestTaskCloseoutResumesEachDestructiveBoundary(t *testing.T) {
	for _, boundary := range []string{"before-worktree-remove", "after-worktree-remove", "before-branch-remove", "after-branch-remove", "before-lifecycle", "after-lifecycle"} {
		t.Run(boundary, func(t *testing.T) {
			f, in := closeoutFixture(t, true, true)
			p, err := PreviewTaskCloseout(f.dependencies, in)
			if err != nil || !p.Ready {
				t.Fatalf("preview %+v %v", p, err)
			}
			_, err = applyTaskCloseout(f.dependencies, in, p.PlanSHA256, func(stage string) error {
				if stage == boundary {
					return errors.New("injected interruption")
				}
				return nil
			})
			if err == nil {
				t.Fatal("fault was not reached")
			}
			before, err := ReadTaskCloseout(f.dependencies, "task")
			if err != nil || before == nil || before.State == "complete" {
				t.Fatalf("partial receipt %+v %v", before, err)
			}
			if boundary == "after-worktree-remove" {
				show, err := ShowTask(f.dependencies, "task")
				if err != nil || show.Target.Kind != "removal_pending" || show.Integration == nil || show.Integration.NextAction.Kind != "resume_closeout" || show.Closeout.WorktreeRemoved {
					t.Fatalf("reserved removal lost truthful recovery readback: %+v %v", show, err)
				}
			}
			actual, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256)
			if err != nil || actual.State != "complete" || actual.OperationID != in.OperationID {
				t.Fatalf("recovery %+v %v", actual, err)
			}
			lifecycle, err := readWorkItemLifecycle(f.root)
			if err != nil || len(lifecycle.Events) != 1 {
				t.Fatalf("lifecycle repeated: %+v %v", lifecycle, err)
			}
		})
	}
}

func TestTaskCloseoutKeepAndConcurrentIdempotence(t *testing.T) {
	f, in := closeoutFixture(t, true, true)
	in.Keep = true
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview %+v %v", p, err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256); errs <- err }()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(f.taskPath); err != nil {
		t.Fatal("keep removed source", err)
	}
	r, err := ReadTaskCloseout(f.dependencies, "task")
	if err != nil || r.ResourceState != "kept" || r.WorktreeRemoved || r.BranchRemoved {
		t.Fatalf("keep receipt %+v %v", r, err)
	}
	changed := in
	changed.Keep = false
	if _, err := ApplyTaskCloseout(f.dependencies, changed, p.PlanSHA256); err == nil {
		t.Fatal("changed deletion scope reused confirmation")
	}
}

func TestTaskCloseoutRetentionFailureLeavesSourceAndReconcileRejectsAmbiguity(t *testing.T) {
	f, in := closeoutFixture(t, true, true)
	p, err := PreviewTaskReconcile(f.dependencies, "task", in.Ownership, false)
	if err != nil || !p.Ready {
		t.Fatalf("legacy preview %+v %v", p, err)
	}
	if err := os.Remove(filepath.Join(f.root, "handoff")); err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256); err == nil || !strings.Contains(err.Error(), "evidence_retention_failed") {
		t.Fatalf("missing evidence allowed: %v", err)
	}
	if _, err := os.Stat(f.taskPath); err != nil {
		t.Fatal("source removed without evidence", err)
	}
	err = f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		second := r.TaskResults[0]
		second.ID = "trs_22222222222222222222222222222222"
		second.PublicationKey = "second/result"
		second.TerminalResultID = "res_22222222222222222222222222222222"
		r.TaskResults = append(r.TaskResults, second)
		sortWorkRegistry(&r)
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := PreviewTaskReconcile(f.dependencies, "task", in.Ownership, false); err == nil || !strings.Contains(err.Error(), "legacy_result_ambiguous") {
		t.Fatalf("ambiguous reconciliation accepted: %v", err)
	}
}

func TestTaskCloseoutSquashProofRetainsOtherwiseUnreachableCandidate(t *testing.T) {
	f, in := closeoutFixture(t, false, true)
	in.TargetRef = "refs/heads/master"
	proof := []byte(`{"fixture":"exact observed external squash merge"}`)
	path := filepath.Join(f.root, "pr-observation.json")
	if err := os.WriteFile(path, proof, 0600); err != nil {
		t.Fatal(err)
	}
	in.ObservedPR = &TaskCloseoutPRIntegration{Repository: "fixture/repository", Number: 42, HeadOID: f.resultOID, BaseRef: in.TargetRef, MergeOID: strings.Repeat("2", 40), URL: "https://example.invalid/fixture/pull/42", MergedAtUTC: "2026-09-28T12:03:00Z", EvidenceLocator: path, EvidenceSHA256: digestTaskBytes(proof)}
	in.RequireIntegration = true
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready || !p.IntegrationObserved {
		t.Fatalf("squash preview %+v %v", p, err)
	}
	receipt, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256)
	if err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "reflog", "expire", "--expire=now", "--all")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "gc", "--prune=now")
	if err := verifyCloseoutRetention(f.dependencies, f.root, receipt.Plan, *receipt.Retention); err != nil {
		t.Fatal(err)
	}
	if got := runLocalGit(t, filepath.Join(f.wrapper, "main"), "rev-parse", "refs/heads/epic"); got != f.oid {
		t.Fatal("PR closeout changed local base")
	}
}
