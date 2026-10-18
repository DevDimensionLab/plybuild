package cmd

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestExecuteBeginAcceptsObservationsWithoutHandwrittenAcceptance(t *testing.T) {
	c, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"begin"})
	if err != nil || c.Name() != "begin" {
		t.Fatalf("native recipient startup is missing: %v", err)
	}
	for _, flag := range []string{"context", "provider", "native-session", "policy", "read-root", "write-root", "permission-confirmed"} {
		if c.Flags().Lookup(flag) == nil {
			t.Errorf("missing actual observation flag %s", flag)
		}
	}
	if c.Flags().Lookup("file") != nil {
		t.Fatal("normal startup still requires a handwritten acceptance document")
	}
}

func TestStartupCapabilitiesExposeDifferentAuthorityEffects(t *testing.T) {
	want := map[string]string{"workflow.execute.begin": "recipient_runtime_acceptance", "workflow.execute.observe-stop": "coordinator_observation", "workflow.execute.resolve-stop": "owner_stop_resolution"}
	for _, op := range coreCapabilities(capabilityBuild{}).WorkflowExtensions {
		if effect, ok := want[op.ID]; ok {
			if op.Effect != effect {
				t.Errorf("wrong effect for %s: %s", op.ID, op.Effect)
			}
			delete(want, op.ID)
		}
	}
	if len(want) != 0 {
		t.Fatalf("startup capabilities absent: %v", want)
	}
}

func TestExecuteCheckInHasNoBlanketPermissionFlag(t *testing.T) {
	c, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"check-in"})
	if err != nil || c.Name() != "check-in" {
		t.Fatalf("missing neutral receipt: %v", err)
	}
	for _, flag := range []string{"context", "provider", "native-session", "policy", "read-root", "write-root"} {
		if c.Flags().Lookup(flag) == nil {
			t.Errorf("missing receipt selector %s", flag)
		}
	}
	for _, flag := range []string{"permission-confirmed", "file"} {
		if c.Flags().Lookup(flag) != nil {
			t.Errorf("receipt requires authority attestation: %s", flag)
		}
	}
}

func TestExecuteResolveStopNamesChangedActionWithoutPermissionClaim(t *testing.T) {
	c, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"resolve-stop"})
	if err != nil || c.Name() != "resolve-stop" {
		t.Fatalf("missing continuation: %v", err)
	}
	for _, flag := range []string{"context", "stop", "resolution", "action"} {
		if c.Flags().Lookup(flag) == nil {
			t.Errorf("missing continuation selector %s", flag)
		}
	}
	if c.Flags().Lookup("permission-confirmed") != nil {
		t.Error("continuation still asks the owner to attest future permission")
	}
}

func TestExecuteRetireIsOneCoordinatorOperation(t *testing.T) {
	c, _, err := newWorkflowExecuteCommand(taskrun.Dependencies{}).Find([]string{"retire"})
	if err != nil || c.Name() != "retire" {
		t.Fatalf("missing coordinator retirement: %v", err)
	}
	for _, flag := range []string{"observer", "reason", "format"} {
		if c.Flags().Lookup(flag) == nil {
			t.Errorf("missing retirement selector %s", flag)
		}
	}
	if c.Flags().Lookup("context") != nil {
		t.Error("retirement must not need the retired owner's private context")
	}
	found := false
	for _, op := range coreCapabilities(capabilityBuild{}).WorkflowExtensions {
		switch op.ID {
		case "workflow.execute.retire":
			found = op.Effect == "coordinator_owner_retirement"
		case "workflow.execute.resolve-stop":
			for _, s := range op.Selectors {
				if s == "--permission-confirmed" {
					t.Error("capabilities still advertise permission confirmation for continuation")
				}
			}
		}
	}
	if !found {
		t.Error("retirement capability missing")
	}
}
