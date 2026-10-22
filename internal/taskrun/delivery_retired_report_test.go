package taskrun

import (
	"os"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type al40CountSource struct {
	workspace.TaskIntegrationGit
	observations int
}

func (g *al40CountSource) ObserveIntegrationWorktree(path, ref string) (workspace.IntegrationWorktreeObservation, error) {
	g.observations++
	return g.TaskIntegrationGit.ObserveIntegrationWorktree(path, ref)
}

func TestAL40RetiredReportsRejectReplayBeforeOwnerValidation(t *testing.T) {
	f, o := retireFixture(t)
	owner := f.D
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "prior-ordinary", Phase: "working", Summary: "Modeled work", Meaning: "Preserve ordinary report history", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	ordinary := writeAny(t, f.R.WorkspaceRoot, "ordinary.json", report)
	var err error
	o, err = WorkflowDeliveryReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, ordinary)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowDeliveryReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, ordinary); err != nil {
		t.Fatal("unretired ordinary replay", err)
	}
	report.EventID = "prior-incomplete"
	report.PreviousEventSHA256 = o.Delivery.LastEventSHA256
	report.Phase = "stopped"
	incomplete := writeAny(t, f.R.WorkspaceRoot, "incomplete.json", report)
	o, err = WorkflowDeliveryIncompleteReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, incomplete)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowDeliveryIncompleteReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, incomplete); err != nil {
		t.Fatal("unretired incomplete replay", err)
	}
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	retireAsCoordinator(&f)
	o, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 review fixture", Reason: "Retain original effects", RetainEffects: true})
	if err != nil {
		t.Fatal(err)
	}
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	report.EventID = "fresh-incomplete"
	report.PreviousEventSHA256 = o.Delivery.LastEventSHA256
	fresh := writeAny(t, f.R.WorkspaceRoot, "fresh-incomplete.json", report)
	cwdCalls := 0
	owner.CWD = func() (string, error) { cwdCalls++; return before.Observed.Target.WorktreeLocator, nil }
	source := &al40CountSource{TaskIntegrationGit: owner.Workspace.IntegrationGit}
	owner.Workspace.IntegrationGit = source
	// Historical show is allowed its ordinary source projection. Additional
	// callback source validation must never precede retired-owner rejection.
	if _, err = WorkflowShow(owner, f.R.WorkspaceRoot, o.RunID); err != nil {
		t.Fatal(err)
	}
	readbackObservations := source.observations
	for _, tc := range []struct {
		name, path string
		incomplete bool
	}{{"ordinary-replay", ordinary, false}, {"incomplete-fresh", fresh, true}, {"incomplete-replay", incomplete, true}} {
		t.Run(tc.name, func(t *testing.T) {
			source.observations = 0
			cwdCalls = 0
			calls := workflowTestCalls(t, f.workflowFixture, "agent get") + workflowTestCalls(t, f.workflowFixture, "pane get")
			var err error
			if tc.incomplete {
				_, err = WorkflowDeliveryIncompleteReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, tc.path)
			} else {
				_, err = WorkflowDeliveryReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, tc.path)
			}
			if err == nil || !strings.Contains(err.Error(), "owner_retired") {
				t.Errorf("retired report not rejected as owner_retired: %v", err)
			}
			if cwdCalls != 0 || source.observations != readbackObservations {
				t.Errorf("owner/source validation before denial: cwd=%d source=%d readback=%d", cwdCalls, source.observations, readbackObservations)
			}
			if after := workflowTestCalls(t, f.workflowFixture, "agent get") + workflowTestCalls(t, f.workflowFixture, "pane get"); after != calls {
				t.Error("liveness inspected retired owner")
			}
			after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if !equal(before.Result.Delivery, after.Result.Delivery) {
				t.Error("retired callback changed history")
			}
		})
	}
}

func TestAL40RetirementRequiresCallerOutsideTaskSource(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	// Model exec of the inspecting CLI inside the source; its own PID is
	// excluded by liveness, so caller location needs its explicit scope guard.
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 source caller", Reason: "Resources only", RetainEffects: true}); err == nil || !strings.Contains(err.Error(), "outside the Task source") {
		t.Fatal("source caller not refused", err)
	}
	if retireQueueCurrent(t, f) == nil {
		t.Fatal("source caller released queue")
	}
}
