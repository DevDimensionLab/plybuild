package taskrun

import "testing"

func TestInventoryStatusFollowupVeto(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s := deliveryInventoryState(t, f, "status-veto")
	deliveryInventoryReport(t, f, &s, "working", "working")
	row := deliveryInventoryRows(t, f)[s.Result.RunID]
	if row.StatusFollowupVeto != "clear" {
		t.Fatal("working native coverage should be clear", row.StatusFollowupVeto)
	}
	proof := writeAny(t, f.R.WorkspaceRoot, "status-refusal.json", map[string]string{"observed": "synthetic refusal"})
	evidence, _, err := deliveryReadBinding(proof, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	stop := DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: s.Result.RunID, Input: DeliveryStopInput{EventID: "status-refusal", Provider: "codex", NativeSession: s.Result.Transport.AgentSessionID, Stage: "work", Observer: "fixture coordinator", Reason: "Synthetic refusal", EvidencePath: proof, Action: "refused action", EffectStatus: DeliveryEffectNotStarted}, Source: evidence, Evidence: evidence, ObservedAtUTC: "2026-10-23T13:00:00Z"}
	if _, err = deliveryAppendEvent(f.D, &s, "status-refusal", "coordinator_stop", stop); err != nil {
		t.Fatal(err)
	}
	if e := workflowSave(f.D, s); e != nil {
		t.Fatal(e)
	}
	row = deliveryInventoryRows(t, f)[s.Result.RunID]
	if row.StatusFollowupVeto == "clear" {
		t.Fatal("active refusal lost through historical adapter")
	}
}
