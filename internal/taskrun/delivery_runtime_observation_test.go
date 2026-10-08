package taskrun

import (
	"testing"
	"time"
)

// This helper is a synthetic recipient observation, never actual runtime or
// human QA evidence. Production owners must inspect their current policy.
func deliveryTestRuntimeObservation(t *testing.T, f deliveryFixture, o WorkflowRun) string {
	t.Helper()
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	var acceptance DeliveryAcceptance
	if err = decode(raw, 1<<20, &acceptance); err != nil {
		t.Fatal(err)
	}
	observation := DeliveryRuntimeObservation{Envelope: Envelope{"ply.workflow.delivery-runtime-observation", 1}, RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, Context: FileBinding{o.Paths.Context, s.ContextSHA256}, Acceptance: *s.Acceptance, RuntimeClaim: acceptance.RuntimeClaim, DeliveryPermission: acceptance.DeliveryPermission, ObservedAtUTC: f.D.Now().UTC().Format(time.RFC3339Nano)}
	return writeAny(t, f.R.WorkspaceRoot, "synthetic-current-runtime.json", observation)
}
