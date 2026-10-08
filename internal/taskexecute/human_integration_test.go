package taskexecute

import (
	"errors"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestExecuteFreezesExplicitHumanIntegrationOwner(t *testing.T) {
	a := workspace.DeliveryAgreement{SchemaVersion: 2, IntegrationOwner: workspace.IntegrationOwnerHuman, Mode: workspace.DeliveryLocalEpic}
	d, in, _, _ := launcherFixtureForProvider(t, "claude", a)
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			return errors.New("fixture stop before provider startup")
		}
		return nil
	}
	out, err := Execute(d, in)
	if err == nil || out.RequestPath == "" {
		t.Fatalf("human owner not reserved: %+v %v", out, err)
	}
	var request taskrun.WorkflowRequest
	if _, err = readPreservedJSON(out.RequestPath, &request); err != nil {
		t.Fatal(err)
	}
	if request.Delivery == nil || request.Delivery.Agreement == nil || !request.Delivery.Agreement.HumanOwnedIntegration() || request.Delivery.LocalIntegration != "none" || len(request.Delivery.Agreement.AllowedEffects()) != 0 {
		t.Fatal("human integration owner or developer stop boundary lost during execution")
	}
	if err = workflowhandoff.ValidateDeliveryMandate(request.HandoffDraft, request.Delivery.Agreement); err != nil {
		t.Fatal(err)
	}
}
