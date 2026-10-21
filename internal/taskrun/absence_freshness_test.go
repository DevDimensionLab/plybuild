package taskrun

import (
	"encoding/json"
	"testing"
	"time"
)

// Herdr is modeled by the existing adapter. The clock is deterministic; no
// current process or provider liveness is inferred from historical transport.
func TestAL20AbsenceSeparatesHistoricalTransportAndFreshInventory(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	s, err := workflowRead(f.R.WorkspaceRoot, f.Owner.RunID)
	if err != nil {
		t.Fatal(err)
	}
	old := s.Result.Transport
	old.State, old.Observation, old.ObservedAt = "working", "fresh", "2026-10-21T09:32:37Z"
	s.Result.Transport = old
	now := time.Date(2026, 10, 21, 16, 6, 25, 0, time.UTC)
	f.D.Now = func() time.Time { return now }
	o := (taskWriters{f.D, s}).herdr("")
	if o.err != nil || o.state != taskWritersAbsent {
		t.Fatalf("absence: %+v", o)
	}
	if o.absence.Transport.Observation != "cached" || o.absence.Transport.ObservedAt != old.ObservedAt || o.absence.Transport.State != "working" {
		t.Fatalf("historical transport misrepresented: %+v", o.absence.Transport)
	}
	if s.Result.Transport != old {
		t.Fatal("changed stored transport")
	}
	raw, err := json.Marshal(o.absence)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	if fields["inventory_observed_at"] != now.Format(time.RFC3339Nano) || o.absence.PaneInventorySHA == "" || o.absence.AgentInventorySHA == "" || o.absence.ProcessLiveness != "unknown" {
		t.Fatalf("fresh inventory lacks its own source/time or invented process inspection: %s", raw)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": nil})
	unknown := (taskWriters{f.D, s}).herdr("")
	if unknown.state != taskWritersUnknown || unknown.err == nil || unknown.absence.InventoryObservedAt != "" {
		t.Fatalf("incomplete inspection gained completed inventory timestamp: %+v", unknown)
	}
}

// This is an actual candidate-CLI retirement/readback with modeled Herdr. The
// unresumed delivery path deliberately has no local-process inspection claim.
func TestAL20CLIRetirementReceiptAndHistoricalReadback(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := retireFixture(t)
	old := o.Transport
	old.State, old.Observation, old.ObservedAt = "working", "fresh", "2026-10-21T09:32:37Z"
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Transport = old
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	f.D.CWD = func() (string, error) { return f.Parent, nil }
	out := al14CLIOK(t, binary, f.workflowFixture, "execute", "retire", o.RunID, "--observer", "AL20 modeled coordinator", "--reason", "Retain historical transport and fresh inventory provenance")
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	r, err := deliveryRetirement(s)
	if err != nil || r == nil {
		t.Fatalf("receipt: %v", err)
	}
	if r.Absence.Transport.Observation != "cached" || r.Absence.Transport.State != old.State || r.Absence.Transport.ObservedAt != old.ObservedAt || r.Absence.InventoryObservedAt == "" || r.Absence.PaneInventorySHA == "" || r.Absence.AgentInventorySHA == "" || r.Absence.ProcessLiveness != "unknown" {
		t.Fatalf("misrepresented receipt: %+v", r.Absence)
	}
	if _, err := time.Parse(time.RFC3339Nano, r.Absence.InventoryObservedAt); err != nil {
		t.Fatal(err)
	}
	event := out.Delivery.Events[len(out.Delivery.Events)-1]
	before := hashFileTest(t, event.Binding.Locator)
	read := al14CLIOK(t, binary, f.workflowFixture, "execute", "show", o.RunID)
	if len(read.Delivery.Events) != len(out.Delivery.Events) || before != hashFileTest(t, event.Binding.Locator) {
		t.Fatal("readback changed retirement")
	}
	// A synthetic old receipt has the exact pre-field shape and misleading fresh
	// label. The current reader must preserve those historical claims unchanged.
	r.Absence.Transport = old
	r.Absence.InventoryObservedAt = ""
	path := writeAny(t, f.R.WorkspaceRoot, "historical-retirement.json", r)
	original := hashFileTest(t, path)
	historical := s
	historical.Result.Delivery = &DeliveryState{Events: []DeliveryEventRecord{{Kind: "owner_retirement", Binding: FileBinding{path, original}}}}
	prior, err := deliveryRetirement(historical)
	if err != nil || prior == nil {
		t.Fatalf("historical receipt rejected: %v", err)
	}
	if prior.Absence.Transport != old || prior.Absence.InventoryObservedAt != "" || original != hashFileTest(t, path) {
		t.Fatal("historical receipt rewritten or refreshed")
	}
}
