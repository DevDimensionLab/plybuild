package workflowhandoff

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"testing"
)

// NotificationFixture is exported only in the test variant. Native notification
// tests use the existing canonical fixture builders and real store validators.
func NotificationFixture(t *testing.T, version int, blocked bool) []byte {
	t.Helper()
	var d Dependencies
	var snapshot Snapshot
	var principal canonicaljson.Object
	var start SubmitResult
	var locator string
	var err error
	if version == 2 {
		var input SubmitInput
		d, _, _, input, _ = prepareBoundTaskStart(t)
		locator = input.HandoffLocator
		snapshot, err = d.Store.ReadByLocator(locator)
		if err != nil {
			t.Fatal(err)
		}
		principal = taskSpecFixtureValue(t, map[string]any{"expected_principal_id": "codex-delivery-agent", "human_start_principal": "fixture human, not actual approval", "start_surface": "isolated service fixture", "session_id": "fixture-session", "runtime_id": "fixture", "model_id": "fixture", "started_at_utc": "2026-09-29T12:00:00Z"}).(canonicaljson.Object)
		start, err = SubmitStart(d, input)
	} else if version == 0 {
		root, target, ref, oid := prepareServiceWorkspace(t)
		d = SystemDependencies()
		draft := append(replaceObjectMember(minimalHandoffDraftValue(target, ref, oid), "schema_version", int64(2)), canonicaljson.Member{Name: "task_spec_binding", Value: nil})
		path := writeCanonicalTestFile(t, root, "notification-general-v2.json", draft)
		created, e := Create(d, CreateInput{DraftPath: path})
		if e != nil {
			t.Fatal(e)
		}
		locator = created.Locator
		snapshot, err = d.Store.ReadByLocator(locator)
		if err != nil {
			t.Fatal(err)
		}
		principal = canonicaljson.Object{{Name: "expected_principal_id", Value: "codex-delivery-agent"}, {Name: "human_start_principal", Value: "fixture human"}, {Name: "start_surface", Value: "isolated fixture"}, {Name: "session_id", Value: "fixture-session"}, {Name: "runtime_id", Value: "fixture"}, {Name: "model_id", Value: "fixture"}, {Name: "started_at_utc", Value: "2026-09-29T12:00:00Z"}}
		startValue := append(replaceObjectMember(minimalStartDraftValue(t, snapshot, principal), "schema_version", int64(2)), canonicaljson.Member{Name: "task_spec_binding", Value: nil})
		sandbox, _ := objectMember(startValue, "sandbox")
		startValue = replaceObjectMember(startValue, "sandbox", replaceObjectMember(sandbox.(canonicaljson.Object), "read_roots", []canonicaljson.Value{root}))
		startPath := writeCanonicalTestFile(t, root, "notification-general-start.json", startValue)
		start, err = SubmitStart(d, SubmitInput{HandoffLocator: locator, DraftPath: startPath})
	} else {
		var created CreateResult
		d, created, snapshot, principal = createTestHandoff(t)
		locator = created.Locator
		path := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "notification-start.json", minimalStartDraftValue(t, snapshot, principal))
		start, err = SubmitStart(d, SubmitInput{HandoffLocator: locator, DraftPath: path})
	}
	if err != nil {
		t.Fatal(err)
	}
	terminal := minimalTerminalDraftValue(snapshot, principal, start)
	if blocked {
		terminal = replaceObjectMember(terminal, "reported_outcome", "blocked")
		terminal = replaceObjectMember(terminal, "summary", "The bounded task stopped safely.")
		terminal = replaceObjectMember(terminal, "meaning", "The bounded task needs a product decision.")
		terminal = replaceObjectMember(terminal, "stop_reasons", []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: "product_decision_required"}, {Name: "detail", Value: "Fixture product decision remains open."}, {Name: "primary", Value: true}}})
	}
	path := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "notification-terminal.json", terminal)
	if _, err = SubmitResultDocument(d, SubmitInput{HandoffLocator: locator, DraftPath: path}); err != nil {
		t.Fatal(err)
	}
	inspected, err := Inspect(d, InspectInput{HandoffLocator: locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	return inspected.Bytes
}
