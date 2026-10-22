package taskexecute

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// Change only the injected snapshot at the existing before-reservation seam.
// No registry, lock, request or real workspace file is changed.
type startConflictSnapshot struct {
	workflowhandoff.WorkspaceObserver
	changed bool
}

func (s *startConflictSnapshot) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	v, err := s.WorkspaceObserver.ObserveContaining()
	if s.changed {
		v.Observation.MarkerSHA256 = digestBytes([]byte("AL50 injected changed snapshot"))
	}
	return v, err
}

func TestExecutePreReservationBindingConflict(t *testing.T) {
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
		t.Fatal(err)
	}
	d.StartupProcessTree = func(int) ([]taskrun.StartupProcess, error) {
		t.Fatal("AL50 must not inspect OS processes")
		return nil, nil
	}
	snapshot := &startConflictSnapshot{WorkspaceObserver: d.Workflow.Workspace}
	d.Workflow.Workspace = snapshot
	reservations := 0
	d.Fault = func(point string) error {
		switch point {
		case "workflow_before_reservation":
			snapshot.changed = true
		case "workflow_after_reservation":
			reservations++
		}
		return nil
	}
	out, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "start bindings changed before reservation") {
		t.Fatalf("AL50 did not reach the locked confirmation conflict: %+v %v", out, err)
	}
	if out.State != "start_pending" || out.Run != nil {
		t.Fatalf("AL50 known conflict misdiagnosed: state=%s run=%+v error=%v", out.State, out.Run, err)
	}
	var startErr *taskrun.Error
	if !errors.As(err, &startErr) || startErr.Code != taskrun.WorkflowStartBindingConflict || startErr.Exit != 4 {
		t.Fatalf("AL50 lost typed conflict or exit: %v", err)
	}
	for _, text := range []string{"No Run", "Task", "preparation", "request", "workflow run start --file " + taskrun.ShellQuote(out.RequestPath) + " --check --format json"} {
		if !strings.Contains(out.NextAction, text) {
			t.Fatalf("AL50 next action missing %q: %s", text, out.NextAction)
		}
	}
	requestBefore, err := os.ReadFile(out.RequestPath)
	if err != nil {
		t.Fatal(err)
	}
	var request taskrun.WorkflowRequest
	if err = json.Unmarshal(requestBefore, &request); err != nil {
		t.Fatal(err)
	}
	snapshot.changed = false
	preview, err := taskrun.WorkflowPreviewStart(d, out.RequestPath)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := preview.(taskrun.WorkflowPreview)
	if !ok || p.Confirmation == nil || reservations != 0 {
		t.Fatalf("AL50 request already reserved: %+v reservations=%d", preview, reservations)
	}
	// Frozen acceptance supplies the exact built CLI. Exercise the recommended
	// read-only command against the preserved fixture, beyond checking its text.
	if binary := os.Getenv("PLY_CANDIDATE_BINARY"); binary != "" {
		command := exec.Command(binary, "workflow", "run", "start", "--file", out.RequestPath, "--check", "--format", "json")
		var stderr bytes.Buffer
		command.Stderr = &stderr
		raw, err := command.Output()
		if err != nil {
			t.Fatalf("AL50 supported inspection failed: %v %s", err, stderr.String())
		}
		var inspected taskrun.WorkflowPreview
		if err := json.Unmarshal(raw, &inspected); err != nil || inspected.RunID != p.RunID || inspected.Confirmation == nil || *inspected.Confirmation != *p.Confirmation {
			t.Fatalf("AL50 inspection did not return the same fresh preview: %s %v", raw, err)
		}
	}
	for _, path := range []string{p.Paths.RunRoot, filepath.Join(root, ".ply", "workflow-runs", "requests", p.RunID+".json"), filepath.Join(root, "bin", "recovery-model.json")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("AL50 unexpected reservation/provider effect at %s: %v", path, err)
		}
	}
	registry, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out.Goal.WorktreePath); err != nil {
		t.Fatal("AL50 lost prepared Task worktree:", err)
	}
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			reservations++
		}
		return nil
	}
	started, err := Execute(d, in)
	if err != nil || started.State != "started" || started.Run == nil || started.Run.RunID != p.RunID || reservations != 1 {
		t.Fatalf("AL50 fresh same-request start: %+v %v reservations=%d", started, err, reservations)
	}
	if after, err := os.ReadFile(out.RequestPath); err != nil || !bytes.Equal(requestBefore, after) {
		t.Fatal("AL50 changed immutable request")
	}
	if after, err := d.Workspace.WorkItems.Snapshot(root); err != nil || !reflect.DeepEqual(registry.TaskPreparations, after.TaskPreparations) {
		t.Fatal("AL50 changed preparation")
	}
	modelPath := filepath.Join(root, "bin", "recovery-model.json")
	before, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	var model struct{ Tabs, Starts, Prompts int }
	if err := json.Unmarshal(before, &model); err != nil || model.Tabs != 1 || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("AL50 expected exactly one injected owner: %+v %v", model, err)
	}
	again, err := Execute(d, in)
	if err != nil || again.State != "existing" || again.Run.RunID != p.RunID {
		t.Fatalf("AL50 repeat: %+v %v", again, err)
	}
	if after, err := os.ReadFile(modelPath); err != nil || !bytes.Equal(before, after) {
		t.Fatal("AL50 repeat touched injected transport")
	}
}

func TestExecutePartialReservationStaysUnknown(t *testing.T) {
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	reservations := 0
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			reservations++
			return errors.New("start bindings changed before reservation") // Same words are not proof.
		}
		return nil
	}
	out, err := Execute(d, in)
	if err == nil || out.State != "start_unknown" || out.Run == nil || out.Run.RunID == "" || reservations != 1 {
		t.Fatalf("AL50 partial reservation was treated as safe: %+v %v", out, err)
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	again, err := Execute(d, in)
	if err != nil || again.State != "existing" || again.Run.RunID != out.Run.RunID || reservations != 1 {
		t.Fatalf("AL50 replayed unknown start: %+v %v", again, err)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("AL50 ordinary repeat changed partial reservation")
	}
}
