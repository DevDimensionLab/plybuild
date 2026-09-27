package workflowhandoff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestShowAndInspectCompleteRoundTrip(t *testing.T) {
	dependencies, created, terminal := completeRoundTrip(t)
	show, err := Show(dependencies, created.HandoffID)
	if err != nil {
		t.Fatal(err)
	}
	if show.Status != "The recipient reported a complete result for the expected run." ||
		show.Result != "The bounded change is complete." ||
		show.Meaning != "Technical checks were reported as passing. Ply validated the file, identity, sequence, and reported policy fields; it has not independently verified the work and no QA or integration is authorized." ||
		!strings.Contains(show.NextAction, "ply workflow handoff inspect") {
		t.Fatalf("show = %#v", show)
	}

	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !inspection.AppendLF {
		t.Fatal("inspection does not request one LF")
	}
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(inspection.Bytes, []byte(snapshot.Handoff.ReplySecret)) || !bytes.Contains(inspection.Bytes, []byte(`"secret":"[REDACTED]"`)) || !bytes.Contains(inspection.Bytes, []byte(`"next_transition_authorized":false`)) {
		t.Fatalf("inspection redaction/authority = %s", inspection.Bytes)
	}
	if _, err := canonicaljson.DecodeStrict(inspection.Bytes); err != nil {
		t.Fatalf("inspection is not canonical JSON: %v", err)
	}

	raw, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Raw: "result"})
	if err != nil {
		t.Fatal(err)
	}
	if raw.AppendLF || digestBytes(raw.Bytes) != terminal.SHA256 {
		t.Fatalf("raw result = appendLF %t digest %s", raw.AppendLF, digestBytes(raw.Bytes))
	}
	if _, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Raw: "handoff"}); !IsClass(err, ErrorInvalidArguments) {
		t.Fatalf("raw handoff without ack error = %v", err)
	}
}

func TestOverBudgetResultIsPreservedAndInspectedAsConflict(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "start.json", minimalStartDraftValue(t, snapshot, principal))
	start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil {
		t.Fatal(err)
	}
	terminal := replaceObjectMember(minimalTerminalDraftValue(snapshot, principal, start), "rounds_used", int64(3))
	resultPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "result.json", terminal)
	if _, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath}); err != nil {
		t.Fatalf("over-budget bound report was not preserved: %v", err)
	}
	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(inspection.Bytes, []byte(`"derived_state":"conflict"`)) || !bytes.Contains(inspection.Bytes, []byte(`rounds_used exceeds handoff budget`)) {
		t.Fatalf("inspection = %s", inspection.Bytes)
	}
}

func TestReportedCompleteWithForbiddenEffectIsPreservedAsConflict(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "start.json", minimalStartDraftValue(t, snapshot, principal))
	start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil {
		t.Fatal(err)
	}
	terminal := minimalTerminalDraftValue(snapshot, principal, start)
	terminal = replaceObjectMember(terminal, "forbidden_effects_observed", []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: "network"}, {Name: "detail", Value: "A forbidden network effect was reported."}, {Name: "artifact_ids", Value: []canonicaljson.Value{}}, {Name: "effect_ids", Value: []canonicaljson.Value{}}}})
	resultPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "result.json", terminal)
	if _, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath}); err != nil {
		t.Fatalf("bound conflict report was not preserved: %v", err)
	}
	show, err := Show(dependencies, created.HandoffID)
	if err != nil {
		t.Fatal(err)
	}
	if show.Status != "The recipient result or preserved evidence conflicts with the handoff." {
		t.Fatalf("show = %#v", show)
	}
	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(inspection.Bytes, []byte(`"derived_state":"conflict"`)) || !bytes.Contains(inspection.Bytes, []byte(`forbidden effects were reported`)) {
		t.Fatalf("inspection = %s", inspection.Bytes)
	}
}

func completeRoundTrip(t *testing.T) (Dependencies, CreateResult, SubmitResult) {
	t.Helper()
	workspaceRoot, target, ref, oid := prepareServiceWorkspace(t)
	draftBytes, _ := canonicaljson.Marshal(minimalHandoffDraftValue(target, ref, oid))
	draftPath := filepath.Join(workspaceRoot, "handoff-draft.json")
	if err := os.WriteFile(draftPath, draftBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	dependencies := SystemDependencies()
	created, err := Create(dependencies, CreateInput{DraftPath: draftPath})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principal := canonicaljson.Object{{Name: "expected_principal_id", Value: "codex-delivery-agent"}, {Name: "human_start_principal", Value: "product-owner"}, {Name: "start_surface", Value: "Codex chat"}, {Name: "session_id", Value: "session-test"}, {Name: "runtime_id", Value: "codex"}, {Name: "model_id", Value: "test-model"}, {Name: "started_at_utc", Value: "2026-09-27T00:00:00Z"}}
	startBytes, _ := canonicaljson.Marshal(minimalStartDraftValue(t, snapshot, principal))
	startPath := filepath.Join(workspaceRoot, "start.json")
	if err := os.WriteFile(startPath, startBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil {
		t.Fatal(err)
	}
	resultBytes, _ := canonicaljson.Marshal(minimalTerminalDraftValue(snapshot, principal, start))
	resultPath := filepath.Join(workspaceRoot, "result.json")
	if err := os.WriteFile(resultPath, resultBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	terminal, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath})
	if err != nil {
		t.Fatal(err)
	}
	return dependencies, created, terminal
}
