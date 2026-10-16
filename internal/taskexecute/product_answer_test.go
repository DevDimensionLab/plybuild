package taskexecute

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func sl01ExecuteFixture(t *testing.T) (taskrun.Dependencies, Input, string) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("worker requires macOS runtime adapter")
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "SL01 native execute regression", TimeoutSeconds: 30}}
	d, in, root, _ := launcherFixtureForProvider(t, "codex", a)
	t.Setenv("HERDR_ENV", "")
	t.Setenv("HERDR_WORKSPACE_ID", "")
	in.Runtime.HerdrWorkspace = ""
	in.Worker = &taskrun.WorkerContract{SchemaVersion: 2, MaxAttempts: 1, AgentTimeoutSeconds: 60, MaxProviderRetries: 0, ProductQuestions: &taskrun.ProductQuestionContract{SchemaVersion: 1, Mode: taskrun.BoundAnswerMode, MaxContinuationTurns: 1}}
	in.ControlContract = taskrun.NewControlContract("week", "week_"+strings.Repeat("a", 32), "sha256:"+strings.Repeat("b", 64))
	return d, in, root
}

func TestSL01ExecuteAcceptsAndPreservesNewWeekQuestionContract(t *testing.T) {
	d, in, root := sl01ExecuteFixture(t)
	in.Check = true
	before := fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain")
	preview, err := Execute(d, in)
	if err != nil || preview.State != "ready" || preview.Runtime == nil || preview.Runtime.Runtime.Mode != "headless" {
		t.Fatalf("new explicitly opted-in week cannot reach native execution preview: state=%s err=%v", preview.State, err)
	}
	if _, err := os.Lstat(filepath.Join(root, ".ply", "task-executions")); !os.IsNotExist(err) {
		t.Fatalf("preview wrote execution state: %v", err)
	}
	if before != fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain") {
		t.Fatal("preview created a Task worktree")
	}

	// Exercise the actual effect route through immutable intent publication,
	// then stop before preparation/Run/provider start. No fake successful Run.
	stop := errors.New("SL01 test stops after the real launcher intent")
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return stop
		}
		return nil
	}
	in.Check = false
	started, err := Execute(d, in)
	if !errors.Is(err, stop) || started.Run != nil {
		t.Fatalf("new worker did not reach the preserved intent boundary: %+v %v", started.Run, err)
	}
	var intent launchIntent
	if _, err := readPreservedJSON(filepath.Join(filepath.Dir(started.RequestPath), "intent.json"), &intent); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(intent.Worker, in.Worker) || !reflect.DeepEqual(intent.ControlContract, in.ControlContract) || intent.Runtime.Runtime.Mode != "headless" {
		t.Fatal("launcher changed the opted-in worker, zero retry limit, or parent binding")
	}
	if before != fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain") {
		t.Fatal("test crossed its before-preparation boundary")
	}
}

func TestSL01ExecuteRejectsUnboundOrMalformedQuestionWorkersBeforeWrites(t *testing.T) {
	for _, name := range []string{"legacy_upgrade", "missing_opt_in", "missing_parent", "standalone_parent", "wrong_mode", "wrong_question_schema", "extra_continuation"} {
		t.Run(name, func(t *testing.T) {
			d, in, root := sl01ExecuteFixture(t)
			switch name {
			case "legacy_upgrade":
				in.Worker.SchemaVersion = 1
			case "missing_opt_in":
				in.Worker.ProductQuestions = nil
			case "missing_parent":
				in.ControlContract = nil
			case "standalone_parent":
				in.ControlContract = taskrun.NewControlContract("standalone", "", "")
			case "wrong_mode":
				in.Worker.ProductQuestions.Mode = "unsupported"
			case "wrong_question_schema":
				in.Worker.ProductQuestions.SchemaVersion = 2
			case "extra_continuation":
				in.Worker.ProductQuestions.MaxContinuationTurns = 2
			}
			before := fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain")
			for _, check := range []bool{true, false} {
				in.Check = check
				if _, err := Execute(d, in); err == nil {
					t.Fatalf("invalid question authority passed with check=%v", check)
				}
			}
			if _, err := os.Lstat(filepath.Join(root, ".ply", "task-executions")); !os.IsNotExist(err) {
				t.Fatalf("rejected worker wrote execution state: %v", err)
			}
			if before != fixtureGit(t, filepath.Join(root, "main"), "worktree", "list", "--porcelain") {
				t.Fatal("rejected worker created a Task worktree")
			}
		})
	}
}
