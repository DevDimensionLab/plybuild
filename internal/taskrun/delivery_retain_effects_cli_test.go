package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Actual binary invocations use only help, show and pre-inspection refusals.
// The disposable native stores and Herdr are modeled. No CLI success here
// claims physical writer absence or permission for OS process inspection.
func TestAL40PublicCLIBoundary(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("requires actual make-built candidate via PLY_CANDIDATE_BINARY")
	}
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectUnknown)
	retireAsCoordinator(&f)
	run := func(args ...string) ([]byte, []byte, error) {
		t.Helper()
		c := exec.Command(binary, args...)
		c.Dir = f.Parent
		stdout, stderr, err := al14Command(t, c)
		exit := 0
		if e, ok := err.(*exec.ExitError); ok {
			exit = e.ExitCode()
		}
		t.Logf("ACTUAL CLI argv=%q exit=%d\nstdout=%s\nstderr=%s", args, exit, stdout, stderr)
		return stdout, stderr, err
	}
	for _, command := range []string{"execute", "run"} {
		stdout, _, err := run("workflow", command, "retire", "--help")
		if err != nil || !strings.Contains(string(stdout), "--retain-effects") {
			t.Fatal("public help", err)
		}
	}
	stdout, _, err := run("workflow", "execute", "show", o.RunID, "--format", "json")
	var shown WorkflowRun
	if err != nil || json.Unmarshal(stdout, &shown) != nil || shown.DeliveryStatus.ResourceRetirement.State != "not_retired" {
		t.Fatal("public show", err, string(stdout))
	}
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	// Bound pane remains live: refusal must precede local process inspection.
	stdout, stderr, err := run("workflow", "execute", "retire", o.RunID, "--retain-effects", "--observer", "AL40 CLI fixture", "--reason", "Retain only", "--format", "json")
	if err == nil || !strings.Contains(string(stderr), "owner_active") {
		t.Fatal("live pane accepted", err, string(stderr))
	}
	if json.Unmarshal(stdout, &shown) != nil {
		t.Fatal("negative stdout not JSON", string(stdout))
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(before.Result.Delivery.Events, after.Result.Delivery.Events) || retireQueueCurrent(t, f) == nil || workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("public refusal had resource effects")
	}
	// Alias routes to the same mode and refuses a known native attempt before
	// any pane or process inspection, even when the owner is modeled absent.
	if err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "verification", State: "attempted"}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	count := workflowTestCalls(t, f.workflowFixture, "pane get")
	_, stderr, err = run("workflow", "run", "retire", o.RunID, "--retain-effects", "--observer", "AL40 CLI fixture", "--reason", "Retain only", "--format", "json")
	if err == nil || !strings.Contains(string(stderr), "unresolved native attempt") || workflowTestCalls(t, f.workflowFixture, "pane get") != count {
		t.Fatal("native attempt did not block before inspection", err, string(stderr))
	}
}
