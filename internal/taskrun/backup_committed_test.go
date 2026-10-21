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

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func backupCommittedFixture(t *testing.T, f deliveryFixture) string {
	t.Helper()
	path := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	for name, content := range map[string][]byte{"useful.txt": []byte("Useful source implementation\n"), "binary.dat": {0, 1, 2, 255, 0, 34}} {
		if err := os.WriteFile(filepath.Join(path, name), content, 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(filepath.Join(path, "useful.txt"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(path, "README.md")); err != nil {
		t.Fatal(err)
	}
	runGit(t, path, "add", "-A")
	runGit(t, path, "commit", "-m", "Useful committed work")
	return path
}

func TestBackupCommittedDryRun(t *testing.T) {
	f, o := backupRetiredFixture(t)
	path := backupCommittedFixture(t, f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", PreserveCommitted: true}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatalf("explicit committed dry-run preview: %v", err)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("preview wrote state")
	}
	if p.Committed == nil || p.Committed.SourceOID != gitOutput(t, path, "rev-parse", "HEAD") {
		t.Fatal("missing source readback")
	}
	a, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	patch, err := os.ReadFile(a.Committed.Locator)
	if err != nil || hash(patch) != a.Committed.SHA256 {
		t.Fatalf("preserved bytes: %v", err)
	}
	// The product only preserves input. Actual application is test-only in a
	// disposable checkout at the exact original base.
	roundtrip := filepath.Join(t.TempDir(), "roundtrip")
	runGit(t, path, "worktree", "add", "--detach", roundtrip, p.Committed.BaseOID)
	c := exec.Command(testGit(t), "-C", roundtrip, "apply", "--index", "--binary", "-")
	c.Stdin = bytes.NewReader(patch)
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("Git patch roundtrip: %v %s", err, out)
	}
	if gitOutput(t, roundtrip, "write-tree") != p.Committed.SourceTree {
		t.Fatal("roundtrip changed actual source tree")
	}
	for _, name := range []string{"useful.txt", "binary.dat"} {
		want, _ := os.ReadFile(filepath.Join(path, name))
		got, err := os.ReadFile(filepath.Join(roundtrip, name))
		if err != nil || !bytes.Equal(want, got) {
			t.Fatalf("roundtrip content %s: %v", name, err)
		}
	}
	if _, err = os.Stat(filepath.Join(roundtrip, "README.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("deletion lost")
	}
	if info, err := os.Stat(filepath.Join(roundtrip, "useful.txt")); err != nil || info.Mode()&0111 == 0 {
		t.Fatal("executable mode lost")
	}
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil || retry.Goal != a.Goal || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatalf("retry changed native publications: %v", err)
	}
	runGit(t, roundtrip, "worktree", "remove", path)
	// Native fresh preparation and handoff consume the preserved CAS document
	// with no old Task worktree. The fresh code is still exactly its own base.
	plan, err := workspace.PreviewTaskGoalExecution(f.D.Workspace, workspace.TaskGoalExecuteInput{SpecID: in.SpecID, Target: workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, input := range plan.Goal.RequiredInputs {
		if input.SHA256 == a.Committed.SHA256 {
			found = true
			got, e := os.ReadFile(input.Locator)
			if e != nil || !bytes.Equal(got, patch) {
				t.Fatal("native required input unavailable")
			}
		}
	}
	if !found {
		t.Fatal("not a native required input")
	}
	prepared, err := workspace.PrepareTaskGoalExecution(f.D.Workspace, plan.Input, plan.Confirmation, workspace.QueueHumanDecision{ActorClaim: "modeled coordinator", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Prepare isolated backup fixture."})
	if err != nil {
		t.Fatal(err)
	}
	prep := prepared.Preparation.Preparation
	draft, err := workflowhandoff.BuildDeliveryHandoffDraft(f.D.Workflow, prep.ID, "modeled coordinator", "fixture owner", "fixture/backup-dryrun", prepared.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{Agreement: prepared.Delivery})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(draft, []byte(a.Committed.SHA256)) || !bytes.Contains(draft, []byte(a.Committed.Locator)) {
		t.Fatal("fresh handoff omitted preserved input")
	}
	if gitOutput(t, prep.Plan.WorktreePath, "rev-parse", "HEAD^{tree}") != plan.ParentTree {
		t.Fatal("product automatically applied source")
	}
}

func TestBackupCommittedPublicationRecovery(t *testing.T) {
	f, o := backupRetiredFixture(t)
	backupCommittedFixture(t, f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", PreserveCommitted: true}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	// Interruption after preservation/Task/problem publication reuses the same
	// CAS bytes and native publication keys, including an unrelated current slot.
	integrationCLIPrepareSuccessor(t, f)
	queue, err := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if err != nil {
		t.Fatal(err)
	}
	store := f.D.Workspace.WorkItems
	remaining := 2
	f.D.Workspace.WorkItems = backupFailStore{store, &remaining}
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("lost ack hidden")
	}
	if raw, e := os.ReadFile(p.Committed.Locator); e != nil || hash(raw) != p.Committed.SHA256 {
		t.Fatal("artifact preservation lost")
	}
	f.D.Workspace.WorkItems = store
	// Lose the acknowledgement of the next Goal/queue publication as well.
	// The same native keys recover the already-published queue append.
	remaining = 2
	f.D.Workspace.WorkItems = backupFailStore{store, &remaining}
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("lost Goal/queue acknowledgement hidden")
	}
	f.D.Workspace.WorkItems = store
	a, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(queue.Current, a.Queue.Current) || len(a.Queue.Pending) != len(queue.Pending)+1 || !equal(queue.Pending, a.Queue.Pending[:len(queue.Pending)]) {
		t.Fatal("unrelated queue changed")
	}
	f.D.Now = func() time.Time { return time.Date(2026, 10, 22, 2, 0, 0, 0, time.UTC) }
	_, _ = WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 300*time.Millisecond)
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil || retry.Goal != a.Goal || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatalf("recovery after follow: %v", err)
	}
}

func TestBackupCommittedGuards(t *testing.T) {
	f, o := backupRetiredFixture(t)
	path := backupCommittedFixture(t, f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", PreserveCommitted: true}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	// Existing default remains pre-work only.
	ordinary := in
	ordinary.PreserveCommitted = false
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, ordinary, false, ""); err == nil {
		t.Fatal("default accepted committed source")
	}
	for _, name := range []string{"untracked", "dirty", "ancestry", "drift", "input-limit", "candidate", "released", "ownership", "collision", "unknown"} {
		t.Run(name, func(t *testing.T) {
			original, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			saved, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := workflowSave(f.D, saved); err != nil {
					t.Fatal(err)
				}
			}()
			changed := in
			switch name {
			case "untracked":
				if err = os.WriteFile(filepath.Join(path, "stray"), []byte("stray"), 0644); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(filepath.Join(path, "stray"))
			case "dirty":
				if err = os.WriteFile(filepath.Join(path, "useful.txt"), []byte("dirty"), 0755); err != nil {
					t.Fatal(err)
				}
				defer runGit(t, path, "checkout", "--", "useful.txt")
			case "ancestry":
				runGit(t, path, "checkout", "--orphan", "unrelated")
				runGit(t, path, "commit", "-m", "unrelated source")
				alien := gitOutput(t, path, "rev-parse", "HEAD")
				runGit(t, path, "checkout", strings.TrimPrefix(original.Observed.Target.Ref, "refs/heads/"))
				runGit(t, path, "reset", "--hard", alien)
				defer runGit(t, path, "reset", "--hard", p.Committed.SourceOID)
			case "drift":
				runGit(t, path, "commit", "--allow-empty", "-m", "new source identity")
				defer runGit(t, path, "reset", "--hard", p.Committed.SourceOID)
			case "input-limit":
				if err = os.WriteFile(filepath.Join(path, "large.txt"), bytes.Repeat([]byte("large input\n"), 400000), 0644); err != nil {
					t.Fatal(err)
				}
				runGit(t, path, "add", "large.txt")
				runGit(t, path, "commit", "-m", "oversized input")
				defer runGit(t, path, "reset", "--hard", p.Committed.SourceOID)
			case "candidate":
				original.Result.Delivery.Candidates = append(original.Result.Delivery.Candidates, DeliveryCandidate{TaskResult: workspace.TaskResultRecord{Artifacts: []workspace.TaskArtifactRecord{}, VerifierResults: []workspace.TaskVerifierResultRecord{}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{}, Review: workspace.TaskReviewRecord{Findings: []workspace.TaskReviewEntry{}, Fixes: []workspace.TaskReviewEntry{}, OpenActionableFindings: []workspace.TaskReviewEntry{}}}})
			case "released":
				original.Result.Delivery.Phase = "completed"
			case "ownership":
				original.Result.Delivery.OwnershipRelease = &FileBinding{}
			case "unknown":
				proof := writeAny(t, f.R.WorkspaceRoot, "unknown-proof.json", map[string]string{"source": "modeled unknown effect"})
				if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "unknown-committed", Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "modeled coordinator", Reason: "Unconfirmed action", Action: "product action", EffectStatus: DeliveryEffectUnknown, EvidencePath: proof}); err != nil {
					t.Fatal(err)
				}
			case "collision":
				changed.TaskID = "task"
			}
			if name == "candidate" || name == "released" || name == "ownership" {
				if err = workflowSave(f.D, original); err != nil {
					t.Fatal(err)
				}
			}
			before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			_, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, changed, true, p.Confirmation)
			if err == nil {
				t.Fatal("unsafe source/input accepted")
			}
			if name == "candidate" && !strings.Contains(err.Error(), "backup_goal_post_work_unsupported") {
				t.Fatalf("candidate missed intended source guard: %v", err)
			}
			t.Logf("%s: %v", name, err)
			if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
				t.Fatal("guard wrote state")
			}
		})
	}
}

func TestBackupCommittedCandidateCLI(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires make-built CLI")
	}
	f, o := backupRetiredFixture(t)
	backupCommittedFixture(t, f)
	evidence := os.Getenv("PLY_AL30_EVIDENCE_DIR")
	if evidence == "" {
		evidence = t.TempDir()
	}
	args := []string{"workflow", "execute", "backup-goal", o.RunID, "--task", "backup", "--spec", "backup-goal", "--provider", "claude", "--preserve-committed", "--format", "json"}
	sequence := 0
	cli := func(argv ...string) []byte {
		sequence++
		c := exec.Command(binary, argv...)
		c.Dir = f.Parent
		var stdout, stderr bytes.Buffer
		c.Stdout = &stdout
		c.Stderr = &stderr
		err := c.Run()
		data, _ := json.Marshal(map[string]any{"argv": c.Args, "cwd": c.Dir, "stdout": stdout.String(), "stderr": stderr.String(), "error": fmt.Sprint(err)})
		if e := os.WriteFile(filepath.Join(evidence, fmt.Sprintf("committed-cli-%02d.json", sequence)), data, 0600); e != nil {
			t.Fatal(e)
		}
		if err != nil {
			t.Fatalf("CLI %v: %v %s", argv, err, stderr.Bytes())
		}
		return stdout.Bytes()
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	var p BackupGoalResult
	if err := json.Unmarshal(cli(args...), &p); err != nil {
		t.Fatal(err)
	}
	if p.Committed == nil || before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("preview missing artifact or wrote state")
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", p.Confirmation)
	var a BackupGoalResult
	if err := json.Unmarshal(cli(apply...), &a); err != nil {
		t.Fatal(err)
	}
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	cli(apply...)
	if after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("CLI retry changed publications")
	}
	raw := cli("workspace", "task", "spec", "show", "backup", "--spec", "backup-goal", "--revision", "1", "--format", "json")
	if !bytes.Contains(raw, []byte(a.Committed.SHA256)) || !bytes.Contains(raw, []byte("required_inputs")) {
		t.Fatal("CLI missing preserved native input")
	}
}

// The post-close check uses the same source/hash confirmation. A writer's last
// commit at the modeled absence boundary cannot be silently preserved instead.
func TestBackupCommittedRetirementBoundary(t *testing.T) {
	f, o := retireFixture(t)
	// Every process observation in this test is synthetic; the original-owner
	// retirement path currently obtains absence through the modeled Herdr seam.
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	path := backupCommittedFixture(t, f)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", PreserveCommitted: true, Retire: &DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Selected committed backup"}}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}

	changed := false
	f.D.Workspace.IntegrationGit = backupBoundaryGit{TaskIntegrationGit: f.D.Workspace.IntegrationGit, before: func(observedPath string) {
		if observedPath == path && !changed && workflowTestCalls(t, f.workflowFixture, "pane close") == 1 {
			changed = true
			runGit(t, path, "commit", "--allow-empty", "-m", "last writer commit before absence")
		}
	}}
	blocked, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err == nil || !strings.Contains(err.Error(), "confirmation") || !changed {
		t.Fatalf("post-close drift not diagnosed: %v", err)
	}
	if blocked.Applied || blocked.Retirement.Completed || retireQueueCurrent(t, f) == nil {
		t.Fatal("drift released or registered source")
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 1 {
		t.Fatal("fixture did not model the close boundary")
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil || before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("stale source retry performed effects")
	}
}

// Wrap the existing Git observer only; all close and absence observations still
// pass through the fixture's existing native Herdr adapter.
type backupBoundaryGit struct {
	workspace.TaskIntegrationGit
	before func(string)
}

func (g backupBoundaryGit) ObserveIntegrationWorktree(path, ref string) (workspace.IntegrationWorktreeObservation, error) {
	g.before(path)
	return g.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
}
func (g backupBoundaryGit) ObserveBackupCommitted(path, common, ref, base, baseTree, source, sourceTree string) (*workspace.BackupCommittedWork, error) {
	return g.TaskIntegrationGit.(workspace.BackupCommittedGit).ObserveBackupCommitted(path, common, ref, base, baseTree, source, sourceTree)
}
