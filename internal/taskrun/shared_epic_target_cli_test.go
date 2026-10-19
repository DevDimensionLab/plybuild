package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// EB01 native journey across the public CLI. Two logical Epics share one clean
// integration base. Provider, runtime permission and review are synthetic
// stand-ins; acceptance, Git integration and native closure are real.
func TestSharedEpicTargetCLIJourneyIntegratesOnlyTheExactCandidate(t *testing.T) {
	binary := automaticCLIBinary(t)
	var siblingTask workspace.TaskWorktreeMutationResult
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic shared-target orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: a, BeforePrepare: func(w workspace.Dependencies, root, epic string, _ workspace.TaskGoalRef) {
		oid := gitOutput(t, epic, "rev-parse", "HEAD")
		if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "sibling", Title: "Sibling Epic", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: oid}); err != nil {
			t.Fatalf("sibling Epic cannot share the integration base: %v", err)
		}
		if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "sibling-task", Title: "Sibling Task", Description: "Separate Task copy on the shared base.", ParentEpicID: "sibling", ProjectID: "ply", RepoID: "ply"}); err != nil {
			t.Fatal(err)
		}
		var err error
		siblingTask, err = workspace.CreateTaskWorktree(w, workspace.TaskWorktreeCreateInput{TaskID: "sibling-task", Branch: "sibling-task", Path: filepath.Join(filepath.Dir(epic), "sibling-task"), ExpectedParentOID: oid})
		if err != nil {
			t.Fatal(err)
		}
	}})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o)))

	plan := f.Prepared.Preparation.Preparation.Plan
	if plan.Target.EpicID != "epic" || plan.Target.ParentLocator != f.Parent || plan.Target.ParentRef != "refs/heads/epic" {
		t.Fatalf("prepared target = %+v", plan.Target)
	}
	if plan.WorktreePath == siblingTask.Task.Worktree.Locator || plan.Branch == siblingTask.Task.Worktree.Ref || plan.WorktreePath == f.Parent {
		t.Fatalf("Tasks must use separate copies: %+v / %+v", plan, siblingTask.Task.Worktree)
	}
	base := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	deliveryTestCandidate(t, f, "shared target")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'shared target'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, deliveryTestReview(t, f, "shared-target"), binary)...)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != base || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("shared base changed before the agreed integration")
	}
	assertAutomaticCLIClosure(t, binary, f, o)
	candidate := o.Delivery.Candidates[len(o.Delivery.Candidates)-1].OID

	// The changed target is visible to the sibling Epic and requires a fresh
	// base; its old basis gives no new Task copy and no Git effect.
	reflog = gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	siblingBase, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "sibling", RepoID: "ply"})
	if err != nil || siblingBase.Freshness != "stale" || deliveryCLILatestBase(siblingBase).OID != base {
		t.Fatalf("sibling Epic base = %+v, %v", siblingBase, err)
	}
	if _, err := workspace.CreateTask(f.D.Workspace, workspace.TaskCreateInput{TaskID: "late-task", Title: "Late Task", Description: "Uses an old basis.", ParentEpicID: "sibling", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	late := filepath.Join(filepath.Dir(f.Parent), "late-task")
	if _, err := workspace.CreateTaskWorktree(f.D.Workspace, workspace.TaskWorktreeCreateInput{TaskID: "late-task", Branch: "late-task", Path: late, ExpectedParentOID: base}); err == nil {
		t.Fatal("stale shared basis created a Task copy")
	}
	if _, err := os.Stat(late); !os.IsNotExist(err) {
		t.Fatalf("refused Task copy left a worktree: %v", err)
	}
	preview, err := workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "sibling", RepoID: "ply"})
	if err != nil || preview.State != "preview" || preview.Plan == nil || preview.Plan.NextOID != candidate {
		t.Fatalf("sibling fresh-base preview = %+v, %v", preview, err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("sibling refusal or preview changed the shared base")
	}
	if gitOutput(t, siblingTask.Task.Worktree.Locator, "rev-parse", "HEAD") != base {
		t.Fatal("integration touched the sibling Task copy")
	}
	sibling, err := workspace.ShowEpic(f.D.Workspace, "sibling")
	if err != nil || len(sibling.Epic.RepoBindings) != 1 || sibling.Epic.RepoBindings[0].Worktree.OID != base {
		t.Fatalf("sibling historical binding changed: %+v, %v", sibling, err)
	}
}

// From a directory shared by two Epics, a goal Spec pending in only one of
// them selects that Epic; without it the caller gets explicit choices.
func TestSharedEpicTargetUniqueGoalSpecSelectsEpicFromSharedDirectory(t *testing.T) {
	checked := false
	newDeliveryFixture(t, "codex", deliveryFixtureOptions{BeforePrepare: func(w workspace.Dependencies, root, epic string, _ workspace.TaskGoalRef) {
		if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "sibling", Title: "Sibling Epic", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: gitOutput(t, epic, "rev-parse", "HEAD")}); err != nil {
			t.Fatal(err)
		}
		t.Chdir(epic)
		acceptance := filepath.Join(root, "acceptance.sh")
		preview, err := workspace.PreviewTaskGoalExecution(w, workspace.TaskGoalExecuteInput{SpecID: "solution", AcceptancePath: acceptance})
		if err != nil || preview.Target.EpicID != "epic" {
			t.Fatalf("unique goal Spec did not select its Epic: %+v, %v", preview.Target, err)
		}
		_, err = workspace.PreviewTaskGoalExecution(w, workspace.TaskGoalExecuteInput{Next: true, AcceptancePath: acceptance})
		if err == nil || !strings.Contains(err.Error(), "--project ply --repo ply --epic epic") || !strings.Contains(err.Error(), "--project ply --repo ply --epic sibling") {
			t.Fatalf("ambiguous shared directory = %v", err)
		}
		checked = true
	}})
	if !checked {
		t.Fatal("shared directory selection was not exercised")
	}
}
