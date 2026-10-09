package taskexecute

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestExecuteHeadlessCheckAndRepeatPreserveOneNativeRun(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("worker requires macOS runtime adapter")
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic worker test", TimeoutSeconds: 30}}
	d, in, root, _ := launcherFixtureForProvider(t, "codex", a)
	t.Setenv("HERDR_ENV", "")
	t.Setenv("HERDR_WORKSPACE_ID", "")
	in.Runtime.HerdrWorkspace = ""
	in.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 60}
	in.Check = true
	before := fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain")
	preview, err := Execute(d, in)
	if err != nil || preview.State != "ready" {
		t.Fatalf("headless preview failed: %+v %v", preview, err)
	}
	if _, err = os.Lstat(filepath.Join(root, ".ply", "task-executions")); !os.IsNotExist(err) {
		t.Fatalf("check wrote execution state: %v", err)
	}
	if before != fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain") {
		t.Fatal("check created a worktree")
	}
	in.Check = false
	first, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	if first.Run == nil || first.Run.Worker == nil || first.Run.Transport.AgentSessionID != "" || first.Run.Transport.TabID != "" {
		t.Fatalf("worker acquired fabricated transport: %+v", first.Run)
	}
	second, err := Execute(d, in)
	if err != nil || second.Run == nil || second.Run.RunID != first.Run.RunID {
		t.Fatalf("repeated start did not preserve Run: %+v %v", second, err)
	}
	var req taskrun.WorkflowRequest
	if _, err = readPreservedJSON(first.RequestPath, &req); err != nil {
		t.Fatal(err)
	}
	if req.Runtime.Mode != "headless" || req.Runtime.Model != "fixture-model" || req.Delivery.ReasoningEffort != "medium" || req.Delivery.AllowLocalInstall {
		t.Fatalf("worker changed assigned runtime/authority: %+v", req)
	}
	in.Worker = nil
	if _, err = Execute(d, in); err == nil {
		t.Fatal("interactive repeat acquired a worker")
	}
}

func TestExecuteHeadlessRejectsHistoricalHumanAgreementBeforeWrites(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("worker requires macOS runtime adapter")
	}
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	in.Runtime.HerdrWorkspace = ""
	in.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 2, AgentTimeoutSeconds: 60}
	if _, err := Execute(d, in); err == nil {
		t.Fatal("historical human contract silently became automatic")
	}
	if _, err := os.Lstat(filepath.Join(root, ".ply", "task-executions")); !os.IsNotExist(err) {
		t.Fatalf("invalid headless start wrote state: %v", err)
	}
}
