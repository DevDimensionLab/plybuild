package workflowhandoff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type fixedSnapshotStore struct {
	Store
	snapshot Snapshot
}

func (store fixedSnapshotStore) ReadByID(string, HandoffID) (Snapshot, error) {
	return store.snapshot, nil
}

func (store fixedSnapshotStore) ReadByLocator(string) (Snapshot, error) {
	return store.snapshot, nil
}

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

func TestInspectionRedactsCapabilitySecretAndProofsWhileRawDocumentsStayExact(t *testing.T) {
	dependencies, created, terminal := completeRoundTrip(t)
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	proofValue := func(document canonicaljson.Object) string {
		value, _ := objectMember(document, "capability_proof")
		proof := value.(canonicaljson.Object)
		return objectString(proof, "value")
	}
	startProof := proofValue(snapshot.Start.Value)
	terminalProof := proofValue(snapshot.Terminal.Value)
	for phase, document := range map[string]canonicaljson.Object{"start": snapshot.Start.Value, "terminal": snapshot.Terminal.Value} {
		payload := removeObjectMember(document, "capability_proof")
		payloadBytes, err := canonicaljson.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		want, err := capabilityProof(snapshot.Handoff.ReplySecret, payloadBytes)
		if err != nil {
			t.Fatal(err)
		}
		if proofValue(document) != want {
			t.Fatalf("%s proof = %s, want %s", phase, proofValue(document), want)
		}
	}
	tampered := replaceObjectMember(snapshot.Terminal.Value, "summary", "Tampered summary.")
	if err := validateStoredAccepted(tampered, "terminal", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret); err == nil {
		t.Fatal("tampered terminal proof was accepted")
	}

	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	for name, sensitive := range map[string]string{"secret": snapshot.Handoff.ReplySecret, "start proof": startProof, "terminal proof": terminalProof} {
		if bytes.Contains(inspection.Bytes, []byte(sensitive)) {
			t.Fatalf("redacted inspection leaked %s", name)
		}
	}
	if bytes.Count(inspection.Bytes, []byte(`"value":"[REDACTED]"`)) != 2 || !bytes.Contains(inspection.Bytes, []byte(`"secret":"[REDACTED]"`)) {
		t.Fatalf("redacted inspection = %s", inspection.Bytes)
	}
	for rawKind, want := range map[string][]byte{"start": snapshot.Start.Bytes, "result": snapshot.Terminal.Bytes} {
		raw, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Raw: rawKind})
		if err != nil {
			t.Fatal(err)
		}
		if raw.AppendLF || !bytes.Equal(raw.Bytes, want) {
			t.Fatalf("raw %s is not byte-exact", rawKind)
		}
	}
	if terminal.SHA256 != digestBytes(snapshot.Terminal.Bytes) {
		t.Fatalf("terminal digest = %s", terminal.SHA256)
	}
}

func TestBlockedBudgetAndUnknownOutcomesHaveExactClosedProjections(t *testing.T) {
	tests := []struct {
		outcome, reason, status, reportedMeaning, meaning string
		rounds                                            int64
	}{
		{outcome: "blocked", reason: "product_decision_required", rounds: 1, status: "The recipient stopped because the bounded work is blocked.", reportedMeaning: "The bounded task needs a product decision.", meaning: "The bounded task needs a product decision. The preserved result is bound to the expected run, but work must not continue under this handoff."},
		{outcome: "budget_exhausted", reason: "round_budget_exhausted", rounds: 2, status: "The recipient stopped after using the allowed round budget.", reportedMeaning: "The bounded task used its allowed rounds.", meaning: "The bounded task used its allowed rounds. The preserved result is bound to the expected run, but work must not continue under this handoff."},
		{outcome: "unknown", reason: "unknown_or_partial_effect", rounds: 1, status: "The recipient reported that the target effect or result is unknown.", reportedMeaning: "The target effect cannot be established.", meaning: "The target effect cannot be established. No continuation is authorized under this handoff."},
	}
	for _, test := range tests {
		t.Run(test.outcome, func(t *testing.T) {
			dependencies, created, snapshot, principal := createTestHandoff(t)
			startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, test.outcome+"-start.json", minimalStartDraftValue(t, snapshot, principal))
			start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
			if err != nil {
				t.Fatal(err)
			}
			terminal := minimalTerminalDraftValue(snapshot, principal, start)
			terminal = replaceObjectMember(terminal, "reported_outcome", test.outcome)
			terminal = replaceObjectMember(terminal, "rounds_used", test.rounds)
			terminal = replaceObjectMember(terminal, "summary", "The bounded task stopped safely.")
			terminal = replaceObjectMember(terminal, "meaning", test.reportedMeaning)
			terminal = replaceObjectMember(terminal, "stop_reasons", []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: test.reason}, {Name: "detail", Value: "The normative stop condition was reached."}, {Name: "primary", Value: true}}})
			resultPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, test.outcome+"-result.json", terminal)
			if _, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath}); err != nil {
				t.Fatal(err)
			}
			show, err := Show(dependencies, created.HandoffID)
			if err != nil {
				t.Fatal(err)
			}
			if show.Status != test.status || show.Result != "The bounded task stopped safely." || show.Meaning != test.meaning || !strings.Contains(show.NextAction, "ply workflow handoff inspect") {
				t.Fatalf("show = %#v", show)
			}
			inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(inspection.Bytes, []byte(`"derived_state":"`+test.outcome+`"`)) || !bytes.Contains(inspection.Bytes, []byte(`"next_transition_authorized":false`)) {
				t.Fatalf("inspection = %s", inspection.Bytes)
			}
		})
	}
}

func TestWithheldSensitiveAndTooLargeArtifactsRemainDescriptorsOnly(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "withheld-start.json", minimalStartDraftValue(t, snapshot, principal))
	start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil {
		t.Fatal(err)
	}
	terminal := minimalTerminalDraftValue(snapshot, principal, start)
	terminal = replaceObjectMember(terminal, "artifacts", []canonicaljson.Value{
		canonicaljson.Object{{Name: "artifact_id", Value: "sensitive"}, {Name: "kind", Value: "withheld"}, {Name: "description", Value: "Sensitive verifier output was withheld."}, {Name: "reason", Value: "sensitive"}, {Name: "media_type", Value: nil}, {Name: "size_bytes", Value: int64(123)}, {Name: "sha256", Value: nil}},
		canonicaljson.Object{{Name: "artifact_id", Value: "too-large"}, {Name: "kind", Value: "withheld"}, {Name: "description", Value: "Oversized verifier output was withheld."}, {Name: "reason", Value: "too_large"}, {Name: "media_type", Value: "application/octet-stream"}, {Name: "size_bytes", Value: int64(64<<20 + 1)}, {Name: "sha256", Value: digestBytes([]byte("descriptor only"))}},
	})
	terminal = replaceObjectMember(terminal, "evidence_gaps", []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: "withheld_artifacts"}, {Name: "detail", Value: "Sensitive or oversized evidence bytes were not copied."}, {Name: "artifact_ids", Value: []canonicaljson.Value{"sensitive", "too-large"}}, {Name: "effect_ids", Value: []canonicaljson.Value{}}}})
	resultPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "withheld-result.json", terminal)
	if _, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath}); err != nil {
		t.Fatal(err)
	}
	snapshot, err = dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	artifactEntries, err := os.ReadDir(filepath.Join(snapshot.Handoff.ReplyRoot, "artifacts", "sha256"))
	if err != nil {
		t.Fatal(err)
	}
	if len(artifactEntries) != 0 {
		t.Fatalf("withheld artifact bytes were copied: %#v", artifactEntries)
	}
	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(inspection.Bytes, []byte(`"kind":"withheld"`)) || !bytes.Contains(inspection.Bytes, []byte(`"evidence_coverage":{"checked":true,"reasons":["terminal result contains withheld evidence","terminal result reports evidence gaps"],"valid":false}`)) {
		t.Fatalf("inspection evidence coverage = %s", inspection.Bytes)
	}
	for name, sensitive := range map[string]string{"reply secret": snapshot.Handoff.ReplySecret, "terminal proof": objectString(func() canonicaljson.Object {
		value, _ := objectMember(snapshot.Terminal.Value, "capability_proof")
		return value.(canonicaljson.Object)
	}(), "value")} {
		if bytes.Contains(inspection.Bytes, []byte(sensitive)) {
			t.Fatalf("inspection leaked %s", name)
		}
	}
	show, err := Show(dependencies, created.HandoffID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(show.Status+show.Result+show.Meaning+show.NextAction, snapshot.Handoff.ReplySecret) {
		t.Fatal("human output leaked reply secret")
	}
}

func TestAllTwelveLifecycleStatesHaveExactClosedHumanAndInspectionProjections(t *testing.T) {
	dependencies, created, _ := completeRoundTrip(t)
	completeSnapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		state, status string
		mutate        func(*Snapshot)
	}{
		{state: "ready", status: "The agent handoff is ready for the expected recipient.", mutate: func(snapshot *Snapshot) { snapshot.Start, snapshot.Terminal, snapshot.HeadState = nil, nil, "ready" }},
		{state: "started", status: "The expected recipient reported a successful start.", mutate: func(snapshot *Snapshot) { snapshot.Terminal, snapshot.HeadState = nil, "ready" }},
		{state: "start_failed", status: "The recipient reported that work could not start.", mutate: func(snapshot *Snapshot) {
			start := *snapshot.Start
			start.Outcome = "rejected"
			snapshot.Start, snapshot.Terminal, snapshot.HeadState = &start, nil, "ready"
		}},
		{state: "complete", status: "The recipient reported a complete result for the expected run.", mutate: func(snapshot *Snapshot) {}},
		{state: "blocked", status: "The recipient stopped because the bounded work is blocked.", mutate: func(snapshot *Snapshot) {
			terminal := *snapshot.Terminal
			terminal.Outcome = "blocked"
			snapshot.Terminal = &terminal
		}},
		{state: "budget_exhausted", status: "The recipient stopped after using the allowed round budget.", mutate: func(snapshot *Snapshot) {
			terminal := *snapshot.Terminal
			terminal.Outcome = "budget_exhausted"
			snapshot.Terminal = &terminal
		}},
		{state: "unknown", status: "The recipient reported that the target effect or result is unknown.", mutate: func(snapshot *Snapshot) {
			terminal := *snapshot.Terminal
			terminal.Outcome = "unknown"
			snapshot.Terminal = &terminal
		}},
		{state: "conflict", status: "The recipient result or preserved evidence conflicts with the handoff.", mutate: func(snapshot *Snapshot) { snapshot.Handoff.MaxRounds = 0 }},
		{state: "cancelled", status: "The unstarted agent handoff was cancelled.", mutate: func(snapshot *Snapshot) {
			snapshot.Start, snapshot.Terminal, snapshot.HeadState = nil, nil, "cancelled"
		}},
		{state: "superseded", status: "The agent handoff was superseded before start.", mutate: func(snapshot *Snapshot) {
			snapshot.Start, snapshot.Terminal, snapshot.HeadState = nil, nil, "superseded"
		}},
		{state: "abandoned_unknown", status: "The started agent handoff was abandoned with target effects unknown.", mutate: func(snapshot *Snapshot) { snapshot.Terminal, snapshot.HeadState = nil, "abandoned_unknown" }},
		{state: "integrity_conflict", status: "Ply found conflicting or incomplete handoff evidence.", mutate: func(snapshot *Snapshot) { snapshot.IntegrityReasons = []string{"injected integrity conflict"} }},
	}
	workspaceSnapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		t.Run(test.state, func(t *testing.T) {
			snapshot := completeSnapshot
			test.mutate(&snapshot)
			local := dependencies
			local.Store = fixedSnapshotStore{snapshot: snapshot}
			local.Workspace = fixedWorkspaceObserver{snapshot: workspaceSnapshot}
			show, err := Show(local, created.HandoffID)
			if err != nil {
				t.Fatal(err)
			}
			if show.Status != test.status {
				t.Fatalf("status = %q, want %q", show.Status, test.status)
			}
			for name, line := range map[string]string{"status": show.Status, "result": show.Result, "meaning": show.Meaning, "next action": show.NextAction} {
				if line == "" || strings.ContainsAny(line, "\r\n") {
					t.Fatalf("%s is not one exact line: %q", name, line)
				}
			}
			inspection, err := Inspect(local, InspectInput{HandoffLocator: snapshot.Handoff.Locator, Format: "json"})
			if err != nil {
				t.Fatal(err)
			}
			if !inspection.AppendLF || !bytes.Contains(inspection.Bytes, []byte(`"derived_state":"`+test.state+`"`)) || !bytes.Contains(inspection.Bytes, []byte(`"next_transition_authorized":false`)) {
				t.Fatalf("inspection for %s = %s", test.state, inspection.Bytes)
			}
		})
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
