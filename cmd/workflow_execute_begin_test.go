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
