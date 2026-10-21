package taskexecute

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Native registry/Git/execute composition, with the existing synthetic Herdr
// and injected process observations. No provider or OS liveness claim.
func backupStartFixture(t *testing.T) (taskrun.Dependencies, string, string, taskrun.BackupGoalInput, RuntimeOptions) {
	return backupStartSourceFixture(t, true)
}

func backupStartSourceFixture(t *testing.T, retire bool) (taskrun.Dependencies, string, string, taskrun.BackupGoalInput, RuntimeOptions) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "modeled coordinator", TimeoutSeconds: 30}}
	d, in, root, _ := launcherFixtureForProvider(t, "codex", a)
	script := strings.Replace(strings.Replace(recoveryLauncherHerdr, `elif cmd == ["pane", "list"]:`, `elif cmd == ["agent", "list"]:
    result = {"agents":[]}
elif cmd == ["pane", "list"]:`, 1), `elif cmd == ["pane", "get"] and model.get("closed"):`, `elif cmd == ["pane", "close"]:
    model["closed"] = True
    model["starts"] = 0
    model["prompts"] = 0
    result = {"ok":True}
elif cmd == ["pane", "get"] and model.get("closed"):`, 1)
	script = strings.Replace(script, `"value":"replacement-fixture-session"`, `"value":"fixture-session-" + str(model.get("tabs",0))`, 1)
	if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	d.HerdrTimeout = 2 * time.Second
	d.LocalProcesses = func(string, string) ([]taskrun.LocalProcess, error) { return nil, nil }
	d.StartupProcessTree = func(pid int) ([]taskrun.StartupProcess, error) {
		return []taskrun.StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	started, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	if !retire {
		proof := filepath.Join(t.TempDir(), "prework-stop.json")
		if err := os.WriteFile(proof, []byte(`{"source":"modeled coordinator not_started observation"}`), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := taskrun.WorkflowObserveDeliveryStop(d, root, started.Run.RunID, taskrun.DeliveryStopInput{EventID: "prework-stop", Provider: "codex", NativeSession: started.Run.Transport.AgentSessionID, Stage: "work", Observer: "modeled coordinator", Reason: "Known first action not started", Action: "first product action", EffectStatus: taskrun.DeliveryEffectNotStarted, EvidencePath: proof})
		if err != nil {
			t.Fatal(err)
		}
		return d, root, started.Run.RunID, taskrun.BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "codex", Retire: &taskrun.DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Selected fresh backup"}}, in.Runtime
	}
	retired, err := taskrun.WorkflowRetireDelivery(d, root, started.Run.RunID, taskrun.DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Fresh isolated backup"})
	if err != nil {
		t.Fatal(err)
	}
	return d, root, retired.RunID, taskrun.BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "codex"}, in.Runtime
}

func TestBackupStartNativeRetry(t *testing.T) {
	d, root, id, in, runtime := backupStartFixture(t)
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	p, err := BackupGoal(d, root, id, in, runtime, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.State != "ready" || p.Runtime == nil || !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("preview readiness or writes")
	}
	started, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !started.Registration.Applied || started.Execution == nil || started.Execution.Run == nil {
		t.Fatalf("missing fresh execution: %+v", started)
	}
	run := started.Execution.Run
	if run.RunID == id || run.DeliveryStatus == nil || run.DeliveryStatus.QualifiedCandidate != nil {
		t.Fatal("inherited source identity or acceptance")
	}
	after := recoveryFiles(t, filepath.Join(root, ".ply"))
	retry, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if retry.Execution.Run.RunID != run.RunID || !reflect.DeepEqual(retry.Execution.Run.Transport, run.Transport) || !reflect.DeepEqual(after, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("reply-loss retry duplicated execution")
	}
	changed := runtime
	changed.HerdrWorkspace = "other"
	if _, err = BackupGoal(d, root, id, in, changed, true, p.Confirmation); err == nil {
		t.Fatal("changed runtime accepted")
	}
}

func TestBackupStartPreservesRegisteredBoundaryAndUnknownLaunch(t *testing.T) {
	d, root, id, in, runtime := backupStartFixture(t)
	p, err := BackupGoal(d, root, id, in, runtime, false, "")
	if err != nil {
		t.Fatal(err)
	}
	// A genuine failure after publication must not hide the exact Goal refs.
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return fmt.Errorf("modeled startup blocked after registration")
		}
		return nil
	}
	blocked, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err == nil || blocked.State != "registered_start_blocked" || !blocked.Registration.Applied || blocked.Registration.Goal.Spec.ManifestSHA256 == "" {
		t.Fatalf("lost registered boundary: %+v %v", blocked, err)
	}
	// The next native attempt stops after reserving tab input. Native
	// reservation state owns that uncertainty and an identical retry reads it.
	d.Fault = func(point string) error {
		if point == "workflow_before_tab_send" {
			return fmt.Errorf("modeled interruption with reserved tab effect")
		}
		return nil
	}
	uncertain, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err == nil || uncertain.Execution == nil || uncertain.Execution.Run == nil {
		t.Fatalf("missing preserved uncertain Run: %+v %v", uncertain, err)
	}
	d.Fault = nil
	modelPath := filepath.Join(root, "bin", "recovery-model.json")
	before, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	retry, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil || retry.State != "existing" || retry.Execution.Run.RunID != uncertain.Execution.Run.RunID {
		t.Fatalf("unknown readback: %+v %v", retry, err)
	}
	after, _ := os.ReadFile(modelPath)
	if !bytes.Equal(before, after) {
		t.Fatal("unknown startup replayed transport")
	}
}

func TestBackupStartRejectsOriginRuntimeAndCollision(t *testing.T) {
	d, root, id, in, runtime := backupStartFixture(t)
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	wrong := runtime
	wrong.PermissionProfile = "missing-profile"
	if _, err := BackupGoal(d, root, id, in, wrong, false, ""); err == nil {
		t.Fatal("invalid profile accepted")
	}
	origin := d.CWD
	d.CWD = func() (string, error) { return filepath.Join(root, "main"), nil }
	if _, err := BackupGoal(d, root, id, in, runtime, false, ""); err == nil {
		t.Fatal("wrong return accepted")
	}
	d.CWD = origin
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("negative preview wrote metadata")
	}
	p, err := BackupGoal(d, root, id, in, runtime, false, "")
	if err != nil {
		t.Fatal(err)
	}
	// Prepare this exact Goal through native execute, then select a different
	// fresh identity: AL27 provenance must reject it before duplicate effects.
	started, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	before = recoveryFiles(t, filepath.Join(root, ".ply"))
	in.TaskID = "other-backup"
	if _, err = BackupGoal(d, root, id, in, runtime, true, p.Confirmation); err == nil {
		t.Fatal("colliding source identity accepted")
	}
	if started.Execution.Run == nil || !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("collision changed existing execution")
	}
}

func TestBackupStartRejectsChangedRuntimeBeforeRunPublication(t *testing.T) {
	for _, boundary := range []string{"execute_before_preparation", "workflow_before_reservation"} {
		for _, choice := range []string{"destination", "permission-profile"} {
			t.Run(boundary+"/"+choice, func(t *testing.T) {
				d, root, id, in, runtime := backupStartFixture(t)
				p, err := BackupGoal(d, root, id, in, runtime, false, "")
				if err != nil {
					t.Fatal(err)
				}
				d.Fault = func(point string) error {
					if point == boundary {
						return fmt.Errorf("modeled interruption before Run publication")
					}
					return nil
				}
				interrupted, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
				if err == nil || !interrupted.Registration.Applied || interrupted.Execution == nil || (interrupted.Execution.Run != nil && interrupted.Execution.Run.RunID != "") {
					t.Fatalf("missing pre-Run boundary: %+v %v", interrupted, err)
				}
				d.Fault = nil
				before := recoveryFiles(t, filepath.Join(root, ".ply"))
				transportBefore, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
				if err != nil {
					t.Fatal(err)
				}
				changed := runtime
				if choice == "destination" {
					changed.HerdrWorkspace = "changed-destination"
				} else {
					changed.PermissionProfile = ":read-only"
				}
				// A fresh confirmation cannot legitimize different choices for the
				// already preserved native intent/request, even with no Run yet.
				preview, err := BackupGoal(d, root, id, in, changed, false, "")
				if err == nil || preview.State == "ready" || preview.Confirmation != "" {
					t.Errorf("changed %s advertised before existing intent validation: state=%s confirmation=%s error=%v", choice, preview.State, preview.Confirmation, err)
				}
				changedConfirmation := boundHash(struct {
					Registration string
					Runtime      RuntimeOptions
				}{p.Registration.Confirmation, changed})
				rejected, err := BackupGoal(d, root, id, in, changed, true, changedConfirmation)
				if err == nil || !strings.Contains(err.Error(), "execute_runtime_conflict") || rejected.Confirmation != "" {
					t.Fatalf("changed apply bypassed native intent validation: %+v %v", rejected, err)
				}
				if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
					t.Fatal("changed runtime preview/apply mutated partial state")
				}

				transportUnchanged, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
				if !bytes.Equal(transportBefore, transportUnchanged) {
					t.Fatal("changed runtime preview invoked transport")
				}
				// Same-input preview and recovery retain the original durable identity.
				same, err := BackupGoal(d, root, id, in, runtime, false, "")
				if err != nil {
					t.Fatal(err)
				}
				if same.Confirmation != p.Confirmation {
					t.Fatal("same-input confirmation changed")
				}
				recovered, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
				if err != nil {
					t.Fatal(err)
				}
				if recovered.Execution == nil || recovered.Execution.Run == nil || recovered.Execution.Run.Transport.WorkspaceID != runtime.HerdrWorkspace || recovered.Registration.Goal != interrupted.Registration.Goal {
					t.Fatal("same-input recovery lost original identity")
				}
				var request taskrun.WorkflowRequest
				if _, err = readPreservedJSON(recovered.Execution.RequestPath, &request); err != nil {
					t.Fatal(err)
				}
				if request.Runtime.PermissionBinding.ProfileID != "fixture" {
					t.Fatal("same-input recovery changed permission profile")
				}
				transportAfter, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
				if bytes.Equal(transportBefore, transportAfter) {
					t.Fatal("same-input recovery never started modeled transport")
				}
			})
		}
	}
}

func TestBackupStartPreRunReadbackShowsPreservedEffectiveProfile(t *testing.T) {
	for _, boundary := range []string{"execute_before_preparation", "workflow_before_reservation"} {
		t.Run(boundary, func(t *testing.T) {
			d, root, id, in, runtime := backupStartFixture(t)
			runtime.PermissionProfile = ":read-only"
			p, err := BackupGoal(d, root, id, in, runtime, false, "")
			if err != nil {
				t.Fatal(err)
			}
			d.Fault = func(point string) error {
				if point == boundary {
					return fmt.Errorf("modeled pre-Run interruption")
				}
				return nil
			}
			interrupted, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
			if err == nil || interrupted.Execution == nil || (interrupted.Execution.Run != nil && interrupted.Execution.Run.RunID != "") {
				t.Fatal("missing preserved pre-Run fixture")
			}
			d.Fault = nil
			before := recoveryFiles(t, filepath.Join(root, ".ply"))
			// Native execute's omitted choice retains its existing launch. It must
			// show the preserved effective profile, not today's default 'fixture'.
			runtime.PermissionProfile = ""
			readback, err := BackupGoal(d, root, id, in, runtime, false, "")
			if err != nil {
				t.Fatal(err)
			}
			if readback.Runtime == nil || readback.Runtime.Runtime.PermissionBinding.ProfileID != ":read-only" {
				t.Fatalf("omitted choice misrepresents preserved profile: %+v", readback.Runtime)
			}
			if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
				t.Fatal("profile readback rewrote partial state")
			}
		})
	}
}

func TestSelectedBackupStartPhaseRecovery(t *testing.T) {
	testSelectedBackupStartPhaseRecovery(t, false)
}

func TestBackupCommittedStartRecovery(t *testing.T) { testSelectedBackupStartPhaseRecovery(t, true) }

func testSelectedBackupStartPhaseRecovery(t *testing.T, committed bool) {
	t.Helper()
	d, root, id, in, runtime := backupStartSourceFixture(t, false)
	if committed {
		registry, err := d.Workspace.WorkItems.Snapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		path := registry.Tasks[0].Worktree.Locator
		if err := os.WriteFile(filepath.Join(path, "useful.txt"), []byte("useful source work\n"), 0644); err != nil {
			t.Fatal(err)
		}
		fixtureGit(t, path, "add", "useful.txt")
		fixtureGit(t, path, "commit", "-m", "useful committed source")
		in.PreserveCommitted = true
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	p, err := BackupGoal(d, root, id, in, runtime, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Runtime == nil || p.Registration.Retirement == nil || !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("preview mutated or omitted prefix")
	}
	// Reject changed destination before closing the old owner.
	changed := runtime
	changed.HerdrWorkspace = "other"
	if _, err = BackupGoal(d, root, id, in, changed, true, p.Confirmation); err == nil {
		t.Fatal("changed destination accepted")
	}
	d.Fault = func(point string) error {
		if point == "backup_after_retirement" {
			return fmt.Errorf("modeled lost response after retirement")
		}
		return nil
	}
	stopped, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err == nil || !stopped.Registration.Retirement.Completed || stopped.Registration.Applied {
		t.Fatalf("lost retirement boundary: %+v %v", stopped, err)
	}
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return fmt.Errorf("modeled lost response after registration")
		}
		return nil
	}
	registered, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err == nil || !registered.Registration.Applied {
		t.Fatal("lost registration boundary", err)
	}
	// AL28 still fences a changed profile even when the native intent has no Run.
	changed = runtime
	changed.PermissionProfile = ":read-only"
	if _, err = BackupGoal(d, root, id, in, changed, false, ""); err == nil {
		t.Fatal("changed profile accepted after registration")
	}
	d.Fault = nil
	started, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil || started.Execution == nil || started.Execution.Run == nil {
		t.Fatalf("start: %+v %v", started, err)
	}
	if committed {
		artifact := started.Registration.Committed
		raw, err := os.ReadFile(started.Execution.Run.Paths.Context)
		if err != nil {
			t.Fatal(err)
		}
		// Context links the exact handoff; its preserved input closure is native.
		if !bytes.Contains(raw, []byte(started.Execution.Run.Handoff.SHA256)) {
			t.Fatal("context lost handoff identity")
		}
		request, err := os.ReadFile(started.Execution.Run.Handoff.Locator)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(request, []byte(artifact.SHA256)) {
			t.Fatal("fresh native context missing committed input")
		}
	}
	after := recoveryFiles(t, filepath.Join(root, ".ply"))
	retry, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil || retry.Execution.Run.RunID != started.Execution.Run.RunID || !reflect.DeepEqual(after, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatalf("duplicate launch: %v", err)
	}
	if committed {
		// A second fresh backup retains the earlier immutable dry-run document
		// alongside its own Run/hash-scoped input. Native retirement stays modeled.
		oldArtifact := started.Registration.Committed
		source := started.Execution.Run
		if _, err = taskrun.WorkflowRetireDelivery(d, root, source.RunID, taskrun.DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Second isolated backup"}); err != nil {
			t.Fatal(err)
		}
		registry, err := d.Workspace.WorkItems.Snapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		path := ""
		for _, task := range registry.Tasks {
			if task.ID == in.TaskID {
				path = task.Worktree.Locator
			}
		}
		if path == "" {
			t.Fatal("fresh source missing")
		}
		if err = os.WriteFile(filepath.Join(path, "second.txt"), []byte("second useful source\n"), 0644); err != nil {
			t.Fatal(err)
		}
		fixtureGit(t, path, "add", "second.txt")
		fixtureGit(t, path, "commit", "-m", "second committed source")
		next := taskrun.BackupGoalInput{TaskID: "backup-two", SpecID: "backup-two-goal", Provider: "codex", PreserveCommitted: true}
		preview, err := taskrun.WorkflowBackupGoal(d, root, source.RunID, next, false, "")
		if err != nil {
			t.Fatal(err)
		}
		registered, err := taskrun.WorkflowBackupGoal(d, root, source.RunID, next, true, preview.Confirmation)
		if err != nil {
			t.Fatal(err)
		}
		shown, err := workspace.ShowTaskSpec(d.Workspace, workspace.TaskContentQuery{TaskID: next.TaskID, SpecID: next.SpecID, Revision: 1})
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := json.Marshal(shown.Value)
		if !bytes.Contains(raw, []byte(oldArtifact.SHA256)) || !bytes.Contains(raw, []byte(registered.Committed.SHA256)) || oldArtifact.DocumentID() == registered.Committed.DocumentID() {
			t.Fatal("carried source input was lost or collided")
		}
		oldBytes, err := os.ReadFile(oldArtifact.Locator)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("sha256:%x", sha256.Sum256(oldBytes)) != oldArtifact.SHA256 {
			t.Fatal("carried document bytes changed")
		}
	}

}

func TestSelectedBackupStartRecoversAfterNativeFollow(t *testing.T) {
	d, root, id, in, runtime := backupStartSourceFixture(t, false)
	p, err := BackupGoal(d, root, id, in, runtime, false, "")
	if err != nil {
		t.Fatal(err)
	}
	started, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil || started.Execution == nil || started.Execution.Run == nil {
		t.Fatalf("start: %+v %v", started, err)
	}
	// Observe the retired source using the public service. Its modeled agent
	// absence changes readback, but starts and stops no provider.
	shown, followErr := taskrun.WorkflowFollow(d, root, id, 300*time.Millisecond)
	t.Logf("retired source follow: %+v %v", shown.Transport, followErr)
	if reflect.DeepEqual(shown.Transport, p.Registration.Retirement.Owner) {
		t.Fatal("fixture did not change transport observation")
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	retry, err := BackupGoal(d, root, id, in, runtime, true, p.Confirmation)
	if err != nil {
		t.Fatalf("same-input recovery after source follow: %v", err)
	}
	if retry.Execution == nil || retry.Execution.Run == nil || retry.Execution.Run.RunID != started.Execution.Run.RunID || !reflect.DeepEqual(retry.Execution.Run.Transport, started.Execution.Run.Transport) || retry.Registration.Goal != started.Registration.Goal || !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("follow recovery changed refs or duplicated launch")
	}
	if !reflect.DeepEqual(retry.Registration.Retirement.Owner, shown.Transport) {
		t.Fatal("full source observation lost from readback")
	}
}
