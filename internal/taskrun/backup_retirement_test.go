package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectedBackupRetirementRetry(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Selected pre-work replacement"}}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Retirement == nil || p.Retirement.Completed || p.Retirement.Stop.Input.NativeSession != o.Transport.AgentSessionID || before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("preview changed source or lost identity")
	}
	// Simulate lost reply immediately after native retirement, before registration.
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, *in.Retire); err != nil {
		t.Fatal(err)
	}
	a, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Applied || !a.Retirement.Completed {
		t.Fatal("missing boundaries")
	}
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
	if err != nil || retry.Goal != a.Goal || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatalf("retry changed effects: %v", err)
	}
	in.Retire.Reason = "changed choice"
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("changed retirement accepted")
	}
}

func TestSelectedBackupRetirementGuards(t *testing.T) {
	for _, kind := range []string{"rule", "unknown", "no-stop", "bad-coordinator", "collision", "changed-source"} {
		t.Run(kind, func(t *testing.T) {
			f, o := retireFixture(t)
			if kind != "no-stop" {
				effect := DeliveryEffectNotStarted
				if kind == "unknown" {
					effect = DeliveryEffectUnknown
				}
				o = retireStopFixture(t, f, o, effect)
			}
			if kind == "rule" {
				proof := writeAny(t, f.R.WorkspaceRoot, "rule.json", map[string]string{"source": "modeled policy refusal"})
				var err error
				o, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "policy", Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "coordinator", Reason: "Actual modeled refusal", EvidencePath: proof, Action: "denied action", EffectStatus: DeliveryEffectNotStarted, Rule: "policy"})
				if err != nil {
					t.Fatal(err)
				}
			}
			retireAsCoordinator(&f)
			in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}}
			if kind == "bad-coordinator" {
				in.Retire.Observer = ""
			}
			if kind == "collision" {
				in.TaskID = "task"
			}
			if kind == "changed-source" {
				runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "commit", "--allow-empty", "-m", "modeled product work")
			}
			before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			_, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
			if err == nil {
				t.Fatal("unsafe preview accepted")
			}
			if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
				t.Fatal("guard mutated state")
			}
			if strings.Contains(err.Error(), "no such") {
				t.Fatal(err)
			}
		})
	}
}

func TestSelectedBackupRetirementPreservesFailedInspection(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"retire_inventory_denied": true})
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("failed inspection accepted")
	}
	calls := workflowTestCalls(t, f.workflowFixture, "pane list")
	if calls != 1 {
		t.Fatalf("denied observation replayed: %d", calls)
	}
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
		t.Fatal("preserved runtime stop bypassed")
	}
	if workflowTestCalls(t, f.workflowFixture, "pane list") != calls {
		t.Fatal("retry repeated denied inspection")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	stops, err := deliveryUnresolvedStops(s)
	if err != nil || len(stops) != 2 || stops[0].Input.EventID != "refusal-not_started" || stops[1].Input.Stage != "delivery" {
		t.Fatalf("runtime stop/history missing: %+v %v", stops, err)
	}
	if retireQueueCurrent(t, f) == nil {
		t.Fatal("failure released queue")
	}
}

func TestSelectedBackupRetirementCandidateCLI(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires explicit make-built CLI")
	}
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	args := []string{"workflow", "execute", "backup-goal", o.RunID, "--task", "backup", "--spec", "backup-goal", "--provider", "claude", "--retire", "--coordinator", "modeled coordinator", "--reason", "Selected pre-work replacement", "--format", "json"}
	sequence := 0
	cli := func(extra ...string) ([]byte, error) {
		sequence++
		c := exec.Command(binary, append(append([]string{}, args...), extra...)...)
		c.Dir = f.Parent
		raw, err := c.CombinedOutput()
		if path := os.Getenv("PLY_AL29_EVIDENCE_DIR"); path != "" {
			data, _ := json.Marshal(map[string]any{"argv": c.Args, "output": string(raw), "error": fmt.Sprint(err)})
			if e := os.WriteFile(filepath.Join(path, fmt.Sprintf("selected-cli-%02d.json", sequence)), data, 0600); e != nil {
				t.Fatal(e)
			}
		}
		return raw, err
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	raw, err := cli()
	if err != nil {
		t.Fatalf("preview: %s %v", raw, err)
	}
	var p BackupGoalResult
	if err = json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	if p.Retirement == nil || p.Retirement.Completed || p.Retirement.Stop.Input.Action == "" {
		t.Fatalf("lost readback: %s", raw)
	}
	if raw, err = cli("--apply", "--confirmation", "wrong"); err == nil || !strings.Contains(string(raw), "confirmation") {
		t.Fatalf("confirmation guard: %s %v", raw, err)
	}
	// Local stand-ins must never be called by --start preview/guards.
	bin := t.TempDir()
	for _, name := range []string{"herdr", "claude"} {
		if e := os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\nexit 97\n"), 0700); e != nil {
			t.Fatal(e)
		}
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	raw, err = cli("--start", "--herdr-workspace", "w-modeled-preview", "--permission-profile", "auto")
	if err != nil {
		t.Fatalf("start preview: %s %v", raw, err)
	}
	var start map[string]any
	if json.Unmarshal(raw, &start) != nil || start["runtime"] == nil || start["state"] != "ready" {
		t.Fatalf("runtime readback: %s", raw)
	}
	if raw, err = cli("--start", "--herdr-workspace", "w-modeled-preview", "--permission-profile", "missing-profile"); err == nil {
		t.Fatalf("profile accepted: %s", raw)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("CLI preview/guard wrote state")
	}
	// Positive close/absence uses the injected native adapter. The real CLI
	// recovers that phase and publishes with exactly the original confirmation.
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Selected pre-work replacement"}); err != nil {
		t.Fatal(err)
	}
	raw, err = cli("--apply", "--confirmation", p.Confirmation)
	if err != nil {
		t.Fatalf("registration recovery: %s %v", raw, err)
	}
	var applied BackupGoalResult
	if json.Unmarshal(raw, &applied) != nil || !applied.Applied || !applied.Retirement.Completed {
		t.Fatalf("phase readback: %s", raw)
	}
	after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	raw, err = cli("--apply", "--confirmation", p.Confirmation)
	if err != nil || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatalf("CLI retry: %s %v", raw, err)
	}
}

func TestSelectedBackupRetirementFailureRecovery(t *testing.T) {
	for _, kind := range []string{"transient-read", "unknown-close"} {
		t.Run(kind, func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
			retireAsCoordinator(&f)
			in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}}
			p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
			if err != nil {
				t.Fatal(err)
			}
			flag := "retire_inventory_transient"
			if kind == "unknown-close" {
				flag = "retire_close_unknown"
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{flag: true})
			if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil {
				t.Fatal("failed operation succeeded")
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{flag: false})
			before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
			if kind == "transient-read" {
				if err != nil || !retry.Applied {
					t.Fatalf("transient read became permanent stop: %+v %v", retry, err)
				}
			} else {
				if err == nil || before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
					t.Fatalf("unknown close replayed: %v", err)
				}
			}
		})
	}
}

func TestSelectedBackupRetirementRejectsNewStopBeforeClose(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}}
	p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
	if err != nil {
		t.Fatal(err)
	}
	proof := writeAny(t, f.R.WorkspaceRoot, "new-stop.json", map[string]string{"source": "new sourced fixture stop"})
	if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "new-stop", Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "coordinator", Reason: "A new actual stop", Action: "different first action", EffectStatus: DeliveryEffectNotStarted, EvidencePath: proof}); err != nil {
		t.Fatal(err)
	}
	before := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
	if _, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation); err == nil || !strings.Contains(err.Error(), "confirmation") {
		t.Fatalf("new stop accepted: %v", err)
	}
	if before != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
		t.Fatal("stale confirmation closed or wrote")
	}
}
