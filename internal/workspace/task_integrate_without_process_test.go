package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Integration is one Git operation the human chooses to run. A Goal's old
// acceptance policy, a missing Goal or missing process records do not gate it.

func taskWithoutGoalFixture(t *testing.T) goalDeliveryFixture {
	t.Helper()
	base := newWorkItemJourneyFixture(t)
	taskPath := filepath.Join(base.wrapper, "task")
	input, _ := ParseTaskWorktreeCreateInput("task", "task", taskPath, base.oid)
	if _, err := CreateTaskWorktree(base.dependencies, input); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "result.txt"), []byte("result\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, taskPath, "add", "result.txt")
	runLocalGit(t, taskPath, "commit", "-m", "plain result")
	result := runLocalGit(t, taskPath, "rev-parse", "HEAD")
	base.dependencies.TaskLifecycleIDs = &sequenceTaskLifecycleIDs{}
	return goalDeliveryFixture{workItemJourneyFixture: base, root: filepath.Dir(base.wrapper), taskPath: taskPath, epicPath: filepath.Join(base.wrapper, "epic"), resultOID: result, input: TaskIntegrationInput{TaskID: "task", ExpectedResultOID: result, ExpectedParentOID: base.oid}}
}

// The bytes `goal register --acceptance human` produced before DP01.
func humanAcceptanceGoal(a map[string]any) {
	a["schema_version"] = 1
	delete(a, "acceptance")
}

func expectIntegratedOnce(t *testing.T, f goalDeliveryFixture) {
	t.Helper()
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	if plan := integrationPlanFromReadback(t, checked.Readback); checked.Readback.Classification != "ready" || !strings.Contains(plan, `"reasons":[]`) {
		t.Fatalf("check = %s %s", checked.Readback.Classification, plan)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Readback.Classification != "exact_effect" || applied.Readback.RecoveryStatus != "complete" || git.count() != 1 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.resultOID {
		t.Fatalf("apply = %s %s merges=%d", applied.Readback.Classification, applied.Readback.RecoveryStatus, git.count())
	}
}

func TestIntegrateDeliversGoalRegisteredWithHumanAcceptance(t *testing.T) {
	f := newGoalDeliveryFixture(t, humanAcceptanceGoal)
	expectIntegratedOnce(t, f)
}

func TestIntegrateDeliversTaskWithoutGoalWithSameCommand(t *testing.T) {
	expectIntegratedOnce(t, taskWithoutGoalFixture(t))
}

func TestIntegrateSeesParentAlreadyMovedToResultByGit(t *testing.T) {
	f := newGoalDeliveryFixture(t, humanAcceptanceGoal)
	runLocalGit(t, f.epicPath, "merge", "--ff-only", f.resultOID)
	out, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	if out.Readback.Classification != "already_integrated" || strings.Contains(integrationPlanFromReadback(t, out.Readback), "goal_requires_human_acceptance") {
		t.Fatalf("check = %s %s", out.Readback.Classification, integrationPlanFromReadback(t, out.Readback))
	}
}

func TestGoalRegistersWithoutDeliveryPolicyAndJoinsQueue(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	in := FirstGoalInput{TaskID: "fresh", SpecID: "fresh-v1", Title: "Plain Goal", Objective: "Register without a delivery policy", Criteria: []string{"The Goal is in the queue."}, Target: queueTargetFixture()}
	preview, err := RegisterFirstGoal(f.dependencies, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	out, err := RegisterFirstGoal(f.dependencies, in, true, preview.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Applied || out.Contract == nil || len(out.Queue.Pending) != 1 || out.Queue.Pending[0].TaskID != "fresh" {
		t.Fatalf("not registered and queued: %#v", out)
	}
}
