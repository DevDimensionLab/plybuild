package cmd

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestExecuteAcceptanceRequiresExactPrivateInputsBeforeWorkspaceAccess(t *testing.T) {
	for _, args := range [][]string{
		{}, {"wfr_fixture", "extra"}, {"wfr_fixture"},
		{"wfr_fixture", "--context", "/private/context.json"},
		{"wfr_fixture", "--file", "/private/choice.json"},
		{"wfr_fixture", "--context", "context.json", "--file", "/private/choice.json"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "choice.json"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "/private/../choice.json"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "/private/choice.json", "--format", "yaml"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "/private/choice.json", "--force"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "/private/choice.json", "--outcome", "pass"},
		{"wfr_fixture", "--context", "/private/context.json", "--file", "/private/choice.json", "--target", "refs/heads/main"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			// Nil workspace dependencies panic if invalid input reaches a native
			// operation. Invalid choices must stop before any observation/effect.
			c := newWorkflowExecuteAcceptanceCommand(taskrun.Dependencies{})
			c.SetOut(&bytes.Buffer{})
			c.SetErr(&bytes.Buffer{})
			c.SetArgs(args)
			if err := c.Execute(); err == nil {
				t.Fatal("invalid acceptance choice reached execution")
			}
		})
	}
}

func TestExecuteAcceptanceExistingRunUsesSelectedGatesAndPreservesFrozenGoal(t *testing.T) {
	original := workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic, SourceRef: "refs/heads/task", TargetRef: "refs/heads/epic", TargetWorktree: "/fixture/epic"}
	effective := original
	effective.SchemaVersion = 3
	effective.Acceptance = &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}
	for _, selected := range []bool{true, false} {
		name, gates := "selected", "Delivery gates: exact candidate automatic pass and preserved review, actual runtime authority"
		if !selected {
			name, gates = "original", "Delivery gates: exact candidate human pass, actual runtime authority"
		}
		t.Run(name, func(t *testing.T) {
			status := &taskrun.DeliveryStatus{}
			if selected {
				status.AcceptanceSelection = &taskrun.DeliveryAcceptanceSelectionStatus{EffectiveAgreement: effective}
				status.Acceptance = &taskrun.DeliveryAcceptanceStatus{DeliveryAcceptanceDecision: workspace.DeliveryAcceptanceDecision{Mode: "automatic", Outcome: "blocked", Reason: "The policy choice still requires candidate verification."}}
			}
			result := taskexecute.Result{State: "existing", Goal: &workspace.TaskGoalExecutePreview{Delivery: &original}, Run: &taskrun.WorkflowRun{RunID: "wfr_fixture", DeliveryStatus: status}}
			c := newWorkflowExecuteCommand(taskrun.Dependencies{})
			var output bytes.Buffer
			c.SetOut(&output)
			if err := writeExecuteResult(c, "text", result); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), gates) || selected && strings.Contains(output.String(), "Delivery gates: exact candidate human pass") {
				t.Fatalf("existing-run text contradicts its validated current acceptance policy:\n%s", output.String())
			}
			output.Reset()
			if err := writeExecuteResult(c, "json", result); err != nil {
				t.Fatal(err)
			}
			var preserved taskexecute.Result
			if err := json.Unmarshal(output.Bytes(), &preserved); err != nil {
				t.Fatal(err)
			}
			if preserved.Goal == nil || !reflect.DeepEqual(preserved.Goal.Delivery, &original) || preserved.Goal.Delivery.SchemaVersion != 1 {
				t.Fatal("rendering rewrote the frozen Goal agreement")
			}
			if selected && (preserved.Run == nil || preserved.Run.DeliveryStatus.AcceptanceSelection == nil || !reflect.DeepEqual(preserved.Run.DeliveryStatus.AcceptanceSelection.EffectiveAgreement, effective)) {
				t.Fatal("JSON lost the separately validated effective agreement")
			}
		})
	}
}

func TestExecuteAcceptancePublicHelpAndCapabilitiesPreserveSeparateMeaning(t *testing.T) {
	c, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"acceptance"})
	if err != nil || c.Name() != "acceptance" {
		t.Fatalf("existing-run policy choice is not discoverable: %v", err)
	}
	for _, meaning := range []string{"DeliveryAcceptanceChoice@1", "no test result, human QA or runtime permission", "workflow execute continue"} {
		if !strings.Contains(c.Long, meaning) {
			t.Fatalf("acceptance help omitted %q", meaning)
		}
	}
	catalog := coreCapabilities(capabilityBuild{})
	if len(catalog.Operations) != 5 || catalog.Coverage != "workspace-core-read" {
		t.Fatal("acceptance changed the historical core read catalog")
	}
	for _, operation := range catalog.WorkflowExtensions {
		if operation.ID != "workflow.execute.acceptance" {
			continue
		}
		if !reflect.DeepEqual(operation.Command, []string{"workflow", "execute", "acceptance"}) || !reflect.DeepEqual(operation.Selectors, []string{"--context", "--file"}) || operation.Effect != "acceptance_policy_selection" || operation.Mode != "explicit_policy_selection" || len(operation.ResultSchemas) != 1 || operation.ResultSchemas[0].Kind != "ply.workflow.run" {
			t.Fatalf("capabilities mislabeled the policy-choice operation: %+v", operation)
		}
		return
	}
	t.Fatal("capabilities omitted the supported existing-run acceptance choice")
}
