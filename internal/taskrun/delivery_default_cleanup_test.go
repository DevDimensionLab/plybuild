package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// CL02: a new automatic local delivery removes its released Task worktree by
// default. The owner's own integrate leaves cleanup pending while it is still
// present in the source; the coordinator's ordinary follow finishes the same
// closeout after the owner has positively left. Old frozen keep is retained.

const legacyAutomaticReleaseReason = "Automatic local integration, queue closure and Epic base update are observed. Release source writing ownership for native closeout while retaining its worktree and branch."

func cleanupCLIFixture(t *testing.T, binary string) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic CL02 orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	acceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance)
	deliveryTestCandidate(t, f, "default cleanup candidate")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, deliveryTestReview(t, f, "cleanup"), binary)...)
	if o.Delivery.Phase != "automatic_acceptance_passed" {
		t.Fatalf("fixture candidate did not qualify: %+v", o.Delivery.Phase)
	}
	return f, o
}

// ownerPaneStandin lets the synthetic Herdr keep the bound owner pane open
// until it is closed; then the pane, its managed agent and every Task writer
// are positively absent. Call it before the Run starts.
func ownerPaneStandin(t *testing.T, f *deliveryFixture) {
	t.Helper()
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
}

// cleanupOwnerLeaves models the released owner's actual exit: its bound pane,
// managed agent and Task processes are gone (requires ownerPaneStandin).
func cleanupOwnerLeaves(t *testing.T, f deliveryFixture) {
	t.Helper()
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
}

// cleanupOwnerHandoff asserts the owner's integrate ended in the CL02 handoff
// with the integration observed and the source still in place.
func cleanupOwnerHandoff(t *testing.T, f deliveryFixture, o WorkflowRun, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), "owner_handoff:") || o.Delivery == nil || o.Delivery.Phase != "closing" || o.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("owner integrate did not end in the pending-cleanup handoff: %v", err)
	}
	if _, statErr := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); statErr != nil {
		t.Fatal("active owner's worktree was removed")
	}
}

// cleanupCLIAt runs the public CLI from an explicit cwd, such as the
// coordinator's registered return checkout outside the Task source.
func cleanupCLIAt(binary, cwd string, args ...string) (WorkflowRun, []byte, error) {
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var o WorkflowRun
	if decodeErr := json.Unmarshal(stdout.Bytes(), &o); decodeErr != nil && err == nil {
		err = decodeErr
	}
	return o, append(stdout.Bytes(), stderr.Bytes()...), err
}

// The owner integrates from its own live worktree. Integration is observed,
// but the source stays in place with an accurate pending diagnosis.
func cleanupOwnerIntegrates(t *testing.T, binary string, f deliveryFixture, o WorkflowRun) (WorkflowRun, string) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	pending, output, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if err == nil || pending.Delivery == nil || pending.Delivery.Phase != "closing" || pending.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("owner integrate did not leave observed integration with pending cleanup: %v\n%s", err, output)
	}
	if !strings.Contains(pending.NextAction.Message, "follow") || !strings.Contains(pending.NextAction.Message, "cleanup") {
		t.Fatalf("pending cleanup has no concrete next action: %+v", pending.NextAction)
	}
	if _, err = os.Stat(source); err != nil {
		t.Fatalf("active owner's worktree was removed under it: %v", err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("automatic delivery did not integrate the exact candidate")
	}
	if closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID); err != nil || closeout != nil && closeout.State == "complete" {
		t.Fatalf("closeout completed before the owner left: %+v %v", closeout, err)
	}
	return pending, gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
}

func assertCleanupRemoved(t *testing.T, f deliveryFixture, o WorkflowRun, reflog string) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if o.Delivery.Phase != "completed" || o.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("delivery not completed: %+v", o.DeliveryStatus)
	}
	if _, err := os.Lstat(source); !os.IsNotExist(err) {
		t.Fatalf("released Task worktree directory remains: %v", err)
	}
	if strings.Contains(gitOutput(t, f.Parent, "worktree", "list", "--porcelain"), source) {
		t.Fatal("Git still lists the removed Task worktree")
	}
	if err := exec.Command(testGit(t), "-C", f.Parent, "rev-parse", "--verify", "--quiet", c.TaskResult.SourceRef).Run(); err == nil {
		t.Fatal("source branch was not handled by the existing closeout ref policy")
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || !closeout.LifecycleCompleted || closeout.ResourceState != "retired" || !closeout.WorktreeRemoved || !closeout.BranchRemoved || closeout.Retention == nil || closeout.Plan.Keep {
		t.Fatalf("native closeout did not record physical removal: %+v %v", closeout, err)
	}
	if gitOutput(t, f.Parent, "rev-parse", closeout.Plan.ArchiveRef) != c.OID {
		t.Fatal("delivered candidate history is not preserved")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("cleanup repeated or changed the Git integration")
	}
	if _, err = os.Stat(f.Parent); err != nil {
		t.Fatal("target worktree was not preserved")
	}
}

func TestDefaultCleanupPublicFollowRemovesReleasedSourceAfterOwnerLeaves(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	done, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err != nil {
		t.Fatalf("coordinator follow did not finish cleanup: %v\n%s", err, output)
	}
	assertCleanupRemoved(t, f, done, reflog)
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
		t.Fatal("follow did not close exactly the released owner's settled pane")
	}
	again, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "5")
	if err != nil {
		t.Fatalf("repeated follow after cleanup failed: %v\n%s", err, output)
	}
	assertCleanupRemoved(t, f, again, reflog)
	shown, output, err := cleanupCLIAt(binary, f.Parent, "execute", "show", o.RunID)
	if err != nil || shown.Delivery.Phase != "completed" {
		t.Fatalf("readback after cleanup failed: %v\n%s", err, output)
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
		t.Fatal("repetition closed a pane again")
	}
}

// Review r001-f01: the coordinator's ordinary follow is usually already
// running when the owner integrates. That same follow observes the change to
// closing and finishes cleanup once the released owner settles.
func TestDefaultCleanupFollowAlreadyRunningFinishesAfterOwnerIntegrates(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "working"})
	observed := workflowTestCalls(t, f.workflowFixture, "agent get")
	follow := exec.Command(binary, "workflow", "execute", "follow", o.RunID, "--timeout", "60", "--format", "json")
	follow.Dir = f.Parent
	var stdout, stderr bytes.Buffer
	follow.Stdout, follow.Stderr = &stdout, &stderr
	if err := follow.Start(); err != nil {
		t.Fatal(err)
	}
	waited := make(chan error, 1)
	go func() { waited <- follow.Wait() }()
	deadline := time.Now().Add(15 * time.Second)
	for workflowTestCalls(t, f.workflowFixture, "agent get") <= observed {
		if time.Now().After(deadline) {
			t.Fatal("follow did not start observing")
		}
		time.Sleep(50 * time.Millisecond)
	}
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	select {
	case err := <-waited:
		if err != nil {
			t.Fatalf("running follow did not finish cleanup: %v\n%s%s", err, stdout.Bytes(), stderr.Bytes())
		}
	case <-time.After(90 * time.Second):
		follow.Process.Kill()
		t.Fatal("running follow did not return")
	}
	var done WorkflowRun
	if err := json.Unmarshal(stdout.Bytes(), &done); err != nil {
		t.Fatal(err)
	}
	assertCleanupRemoved(t, f, done, reflog)
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
		t.Fatal("running follow did not close exactly the released owner's settled pane")
	}
}

// Review r001-f02: the owner's normal integrate ends in a distinct handoff,
// not a stop. Its writing is finished, the coordinator acts next, and a
// repeated integrate is harmless: no new event and no new Git effect.
func TestDefaultCleanupOwnerIntegrateIsHandoffNotStop(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	pending, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	for attempt := 0; attempt < 2; attempt++ {
		again, output, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 5 || !bytes.Contains(output, []byte("owner_handoff:")) {
			t.Fatalf("owner integrate is not the distinct pending-cleanup handoff: %v\n%s", err, output)
		}
		if again.NextAction.Actor != "coordinator" || pending.NextAction.Actor != "coordinator" {
			t.Fatalf("owner is still named as the next actor: %+v", again.NextAction)
		}
		if len(again.Delivery.Events) != len(pending.Delivery.Events) || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
			t.Fatal("repeated owner integrate changed native history or Git")
		}
	}
	if _, err := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("owner handoff removed the active owner's worktree")
	}
}

func TestDefaultCleanupActiveOwnerBlocksOnlyRemoval(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "working"})
	waiting, _, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "2")
	if err == nil || waiting.Delivery == nil || waiting.Delivery.Phase != "closing" {
		t.Fatalf("follow treated a working owner as released: %v %+v", err, waiting.Delivery)
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("follow closed a working owner's pane")
	}
	if _, err = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("active owner's worktree was removed")
	}
	// The owner finishes; the same recorded closeout continues.
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	done, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err != nil {
		t.Fatalf("follow did not continue after the owner settled: %v\n%s", err, output)
	}
	assertCleanupRemoved(t, f, done, reflog)
}

func TestDefaultCleanupDirtySourceStaysPendingThenContinues(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	stray := filepath.Join(source, "late-untracked.txt")
	if err := os.WriteFile(stray, []byte("late owner output\n"), 0600); err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	blocked, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err == nil || blocked.Delivery == nil || blocked.Delivery.Phase != "closing" || blocked.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("dirty source did not stop only cleanup: %v\n%s", err, output)
	}
	if !strings.Contains(blocked.NextAction.Message, "untracked") {
		t.Fatalf("pending cleanup lacks its concrete cause: %+v", blocked.NextAction)
	}
	if _, err = os.Stat(stray); err != nil {
		t.Fatal("dirty source content was removed")
	}
	if err = os.Remove(stray); err != nil {
		t.Fatal(err)
	}
	done, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err != nil {
		t.Fatalf("same closeout did not continue after the source was clean: %v\n%s", err, output)
	}
	assertCleanupRemoved(t, f, done, reflog)
}

// legacyRetainDraft reproduces a mandate frozen before the cleanup effect
// existed: cleanup is a forbidden effect of the local delivery contract.
func legacyRetainDraft(t *testing.T, raw json.RawMessage) json.RawMessage {
	t.Helper()
	var draft map[string]any
	if err := json.Unmarshal(raw, &draft); err != nil {
		t.Fatal(err)
	}
	authority := draft["authority"].(map[string]any)
	allowed := []any{}
	for _, item := range authority["allowed_effects"].([]any) {
		if effect := item.(map[string]any); effect["type"] != "cleanup" {
			effect["sequence"] = len(allowed) + 1
			allowed = append(allowed, effect)
		}
	}
	if len(allowed) == len(authority["allowed_effects"].([]any)) {
		t.Fatal("new automatic mandate does not grant source cleanup")
	}
	authority["allowed_effects"] = allowed
	authority["forbidden_effects"] = append([]any{map[string]any{"type": "cleanup", "reason": "Outside the selected local delivery contract."}}, authority["forbidden_effects"].([]any)...)
	out, err := Canonical(draft)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// A Run started under an older mandate that forbids cleanup and has no owner
// release yet keeps that meaning when a newer CLI continues its delivery.
func TestDefaultCleanupOldFrozenRunWithoutReleaseKeepsSource(t *testing.T) {
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	f.R.HandoffDraft = legacyRetainDraft(t, f.R.HandoffDraft)
	f.File = writeAny(t, f.R.WorkspaceRoot, "legacy-delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "legacy retained candidate")
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "legacy"), binary)
	if err != nil || o.Delivery.OwnershipRelease != nil {
		t.Fatalf("legacy fixture did not qualify before release: %v", err)
	}
	c := o.Delivery.Candidates[0]
	done, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || done.Delivery.Phase != "completed" {
		t.Fatalf("legacy retained delivery did not complete: %+v %v", done.Delivery, err)
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.ResourceState != "kept" || !closeout.Plan.Keep || closeout.WorktreeRemoved || closeout.BranchRemoved {
		t.Fatalf("old frozen keep became automatic removal: %+v %v", closeout, err)
	}
	if _, err = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("old retained worktree was removed")
	}
	if gitOutput(t, f.Parent, "rev-parse", c.TaskResult.SourceRef) != c.OID {
		t.Fatal("old retained branch was removed")
	}
	raw, err := workflowBound(*done.Delivery.OwnershipRelease, 1<<20)
	if err != nil || !bytes.Contains(raw, []byte(legacyAutomaticReleaseReason)) {
		t.Fatalf("release does not state the retained contract: %v", err)
	}
}

type cleanupCwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f cleanupCwdFiles) Getwd() (string, error) { return f.cwd, nil }

// cleanupAt models the process cwd that native closeout observes.
func cleanupAt(d Dependencies, cwd string) Dependencies {
	w := d.Workspace
	w.Files = cleanupCwdFiles{FileSystem: w.Files, cwd: cwd}
	d.Workspace = w
	d.Workflow.TaskWorkspace = &w
	return d
}

// A headless integrate attempt runs in the Task worktree. Its controller
// finishes the same closeout after the attempt process has ended, outside
// the source, without a second integration.
func TestDefaultCleanupHeadlessWorkerRemovesSourceAfterAttemptExit(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	workerTestReserveNative(t, f, o.RunID, "verify")
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary); err != nil {
		t.Fatal(err)
	}
	workerTestReserveNative(t, f, o.RunID, "integrate")
	if _, err := WorkflowDeliveryIntegrate(cleanupAt(f.D, source), f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("attempt inside the source reported completed cleanup")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	c := s.Result.Delivery.Candidates[0]
	if s.Result.Delivery.Phase != "closing" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatalf("attempt did not leave observed integration pending cleanup: %s", s.Result.Delivery.Phase)
	}
	if _, err = os.Stat(source); err != nil {
		t.Fatal("source removed under the live attempt process")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	stream := filepath.Join(s.Result.Worker.Job.Directory, "native-output.txt")
	if err = f.D.writeOnce(stream, nil); err != nil {
		t.Fatal(err)
	}
	result := workerAttemptResult{RunID: o.RunID, JobID: s.Result.Worker.Job.ID, IntentSHA256: s.Result.Worker.Job.IntentHash, FinishedAt: "2026-10-19T00:00:01Z", Error: "native closeout pending", NativeCompleted: false, Stdout: FileBinding{stream, hash(nil)}, Stderr: FileBinding{stream, hash(nil)}}
	if _, err = workflowKeep(f.D, filepath.Join(s.Result.Worker.Job.Directory, "result.json"), result); err != nil {
		t.Fatal(err)
	}
	probe := filepath.Join(workerDirectory(s), "policy", fmt.Sprintf("probe-%06d.json", s.Result.Worker.Generation+1))
	if _, err = workflowKeep(f.D, probe, map[string]any{"passed": true, "fixture": "synthetic runtime for native closeout boundary"}); err != nil {
		t.Fatal(err)
	}
	f.D.WorkerOwner.lock.Close()
	if err = ServeWorker(f.D, f.R.WorkspaceRoot, o.RunID); err != nil {
		t.Fatal(err)
	}
	done, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	assertCleanupRemoved(t, f, done, reflog)
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.IntegrationResults) != 1 {
		t.Fatalf("headless cleanup repeated integration: %v", err)
	}
}

// The detached worker controller must not hold its cwd inside the Task source
// it later removes; only its integrate attempt runs there.
func TestDefaultCleanupWorkerControllerRunsOutsideSource(t *testing.T) {
	var recorded string
	f, o := newWorkerDeliveryFixture(t, func(f *deliveryFixture) {
		ply := filepath.Join(f.R.WorkspaceRoot, "recording-ply")
		recorded = filepath.Join(f.R.WorkspaceRoot, "supervise-cwd")
		if err := os.WriteFile(ply, []byte("#!/bin/sh\nif [ \"$1 $2 $3\" = 'workflow worker _supervise' ]; then pwd -P > "+ShellQuote(recorded)+"; fi\nexit 0\n"), 0700); err != nil {
			t.Fatal(err)
		}
		f.R.Runtime.PlyExecutable = Executable{Path: ply, SHA256: hashFileTest(t, ply)}
		f.D.Executable = func() (string, error) { return ply, nil }
	})
	if _, err := StartWorker(f.D, f.R.WorkspaceRoot, o.RunID, false); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		b, err := os.ReadFile(recorded)
		if err == nil && len(b) > 0 {
			got := strings.TrimSpace(string(b))
			source := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			if pathWithinTest(source, got) || got != f.R.WorkspaceRoot {
				t.Fatalf("worker controller cwd %s is not the workspace control directory", got)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("worker controller was not started")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func pathWithinTest(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func TestAR03CleanupUnknownInspectionPreservesDeliveredIntegration(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	_, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	cleanupOwnerLeaves(t, f)
	workflowTestModel(t, f.workflowFixture, map[string]any{"retire_inventory_denied": true})
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	calls := workflowTestCalls(t, f.workflowFixture, "pane list")
	blocked, _, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "5")
	if err == nil || blocked.DeliveryStatus.FinalDelivery.State != "delivered" || blocked.Delivery.Phase != "closing" {
		t.Fatalf("unknown inspection lost delivered state: %+v %v", blocked, err)
	}
	if n := workflowTestCalls(t, f.workflowFixture, "pane list") - calls; n != 1 {
		t.Fatalf("replayed denied inspection %d times", n)
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(before.Result.Delivery.Events, after.Result.Delivery.Events) {
		t.Fatal("failed cleanup rewrote history")
	}
	if _, err = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("failed cleanup removed source")
	}
	if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("failed cleanup repeated integration")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"retire_inventory_denied": false})
	done, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err != nil {
		t.Fatalf("same closeout did not resume: %v\n%s", err, output)
	}
	assertCleanupRemoved(t, f, done, reflog)
}
