package taskexecute

import (
	"bytes"
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
