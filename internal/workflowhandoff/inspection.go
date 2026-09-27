package workflowhandoff

import (
	"fmt"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func Show(dependencies Dependencies, id HandoffID) (ShowResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return ShowResult{}, err
	}
	if _, err := ParseHandoffID(string(id)); err != nil {
		return ShowResult{}, err
	}
	workspaceSnapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		return ShowResult{}, err
	}
	snapshot, err := dependencies.Store.ReadByID(workspaceSnapshot.Observation.Root, id)
	if err != nil {
		return ShowResult{}, err
	}
	state, _, _ := deriveState(snapshot)
	inspectAction := fmt.Sprintf("Start the separate result control with `ply workflow handoff inspect --handoff %s --format json`.", snapshot.Handoff.Locator)
	result := ShowResult{Result: "No terminal result has been received."}
	switch state {
	case "ready":
		result.Status = "The agent handoff is ready for the expected recipient."
		result.Meaning = "The immutable contract is published, but no recipient start has been received."
		result.NextAction = fmt.Sprintf("Open a fresh recipient agent in %s and tell it: \"Read and execute the handoff at %s.\"", snapshot.Handoff.Target.Worktree, snapshot.Handoff.Locator)
	case "started":
		result.Status = "The expected recipient reported a successful start."
		result.Meaning = "Ply validated the bound start receipt; it has not received or verified a terminal result."
		result.NextAction = "No action now; wait for the recipient result."
	case "start_failed":
		result.Status = "The recipient reported that work could not start."
		result.Meaning = "The bound start receipt did not authorize target effects; no continuation is authorized under this handoff."
		result.NextAction = inspectAction
	case "complete":
		result.Status = "The recipient reported a complete result for the expected run."
		result.Result = snapshot.Terminal.Summary
		result.Meaning = snapshot.Terminal.Meaning + " Ply validated the file, identity, sequence, and reported policy fields; it has not independently verified the work and no QA or integration is authorized."
		result.NextAction = inspectAction
	case "blocked":
		result.Status = "The recipient stopped because the bounded work is blocked."
		result.Result = snapshot.Terminal.Summary
		result.Meaning = snapshot.Terminal.Meaning + " The preserved result is bound to the expected run, but work must not continue under this handoff."
		result.NextAction = inspectAction
	case "budget_exhausted":
		result.Status = "The recipient stopped after using the allowed round budget."
		result.Result = snapshot.Terminal.Summary
		result.Meaning = snapshot.Terminal.Meaning + " The preserved result is bound to the expected run, but work must not continue under this handoff."
		result.NextAction = inspectAction
	case "unknown":
		result.Status = "The recipient reported that the target effect or result is unknown."
		if snapshot.Terminal != nil {
			result.Result = snapshot.Terminal.Summary
			result.Meaning = snapshot.Terminal.Meaning
		} else {
			result.Meaning = "Target effects may have occurred and Ply cannot infer a safe retry."
		}
		result.Meaning += " No continuation is authorized under this handoff."
		result.NextAction = inspectAction
	case "conflict":
		result.Status = "The recipient result or preserved evidence conflicts with the handoff."
		if snapshot.Terminal != nil {
			result.Result = snapshot.Terminal.Summary
			result.Meaning = snapshot.Terminal.Meaning
		} else {
			result.Meaning = "Ply preserved conflicting evidence; no continuation or target retry is authorized."
		}
		result.Meaning += " No next transition is authorized."
		result.NextAction = inspectAction
	case "cancelled":
		result.Status = "The unstarted agent handoff was cancelled."
		result.Meaning = "The immutable handoff remains readable, but recipient start and result submission are closed."
		result.NextAction = "No action now; the handoff is closed."
	case "superseded":
		result.Status = "The agent handoff was superseded before start."
		result.Meaning = "This immutable handoff is closed and a replacement run exists under the same activity."
		result.NextAction = inspectAction
	case "abandoned_unknown":
		result.Status = "The started agent handoff was abandoned with target effects unknown."
		result.Meaning = "Target effects may have occurred; Ply will not infer completion or recommend a retry."
		result.NextAction = inspectAction
	default:
		result.Status = "Ply found conflicting or incomplete handoff evidence."
		if snapshot.Terminal != nil {
			result.Result = snapshot.Terminal.Summary
		}
		result.Meaning = "No lifecycle transition is safe until a separate control resolves the preserved evidence."
		result.NextAction = inspectAction
	}
	return result, nil
}

func Inspect(dependencies Dependencies, input InspectInput) (InspectResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return InspectResult{}, err
	}
	if input.HandoffLocator == "" {
		return InspectResult{}, InvalidArguments("--handoff is required")
	}
	if (input.Format == "") == (input.Raw == "") {
		return InspectResult{}, InvalidArguments("exactly one of --format or --raw is required")
	}
	if input.Format != "" && input.Format != "json" {
		return InspectResult{}, InvalidArguments("--format must be json")
	}
	if input.Raw != "" && !setOf("handoff", "start", "result")[input.Raw] {
		return InspectResult{}, InvalidArguments("--raw must be handoff, start, or result")
	}
	if input.AcknowledgeSecret && (input.Raw != "handoff") {
		return InspectResult{}, InvalidArguments("--acknowledge-secret-exposure is only valid with --raw handoff")
	}
	snapshot, err := dependencies.Store.ReadByLocator(input.HandoffLocator)
	if err != nil {
		return InspectResult{}, err
	}
	if input.Raw != "" {
		switch input.Raw {
		case "handoff":
			if !input.AcknowledgeSecret {
				return InspectResult{}, InvalidArguments("--raw handoff requires --acknowledge-secret-exposure")
			}
			return InspectResult{Bytes: append([]byte(nil), snapshot.Handoff.Bytes...)}, nil
		case "start":
			if snapshot.Start == nil {
				return InspectResult{}, classified(ErrorNotFound, "no accepted start receipt exists", nil)
			}
			return InspectResult{Bytes: append([]byte(nil), snapshot.Start.Bytes...)}, nil
		case "result":
			if snapshot.Terminal == nil {
				return InspectResult{}, classified(ErrorNotFound, "no accepted terminal result exists", nil)
			}
			return InspectResult{Bytes: append([]byte(nil), snapshot.Terminal.Bytes...)}, nil
		}
	}
	state, policyReasons, evidenceReasons := deriveState(snapshot)
	redacted := redactHandoff(snapshot.Handoff.Value)
	documents := canonicaljson.Object{{Name: "handoff", Value: documentView(snapshot.Handoff.Locator, snapshot.Handoff.SHA256, snapshot.Handoff.Bytes, redacted)}, {Name: "start_receipt", Value: nil}, {Name: "terminal_result", Value: nil}}
	if snapshot.Start != nil {
		documents = replaceObjectMember(documents, "start_receipt", documentView(snapshot.Start.Locator, snapshot.Start.SHA256, snapshot.Start.Bytes, redactCapabilityProof(snapshot.Start.Value)))
	}
	if snapshot.Terminal != nil {
		documents = replaceObjectMember(documents, "terminal_result", documentView(snapshot.Terminal.Locator, snapshot.Terminal.SHA256, snapshot.Terminal.Bytes, redactCapabilityProof(snapshot.Terminal.Value)))
	}
	valid := func(checked bool, reasons []string) canonicaljson.Object {
		return canonicaljson.Object{{Name: "checked", Value: checked}, {Name: "valid", Value: checked && len(reasons) == 0}, {Name: "reasons", Value: stringValues(reasons)}}
	}
	validations := canonicaljson.Object{{Name: "schema", Value: valid(true, nil)}, {Name: "digest", Value: valid(true, nil)}, {Name: "lifecycle", Value: valid(true, nil)}, {Name: "binding", Value: valid(snapshot.Start != nil, nil)}, {Name: "capability", Value: valid(snapshot.Start != nil, nil)}, {Name: "principal_session", Value: valid(snapshot.Terminal != nil, nil)}, {Name: "policy", Value: valid(snapshot.Terminal != nil, policyReasons)}, {Name: "evidence_coverage", Value: valid(snapshot.Terminal != nil, evidenceReasons)}}
	attempts := make([]canonicaljson.Value, 0, len(snapshot.Attempts))
	sort.Slice(snapshot.Attempts, func(i, j int) bool {
		if snapshot.Attempts[i].Phase == snapshot.Attempts[j].Phase {
			return snapshot.Attempts[i].Sequence < snapshot.Attempts[j].Sequence
		}
		return snapshot.Attempts[i].Phase < snapshot.Attempts[j].Phase
	})
	for _, attempt := range snapshot.Attempts {
		attempts = append(attempts, canonicaljson.Object{{Name: "phase", Value: attempt.Phase}, {Name: "classification", Value: attempt.Classification}, {Name: "document_id", Value: nullable(attempt.DocumentID)}, {Name: "digest", Value: nullable(attempt.Digest)}, {Name: "size_bytes", Value: attempt.SizeBytes}, {Name: "locator", Value: nullable(attempt.Locator)}, {Name: "error_class", Value: attempt.ErrorClass}, {Name: "detail", Value: attempt.Detail}})
	}
	action, reason := recommendation(state)
	inspection := envelope("ply.workflow.handoff-inspection", canonicaljson.Member{Name: "derived_state", Value: state}, canonicaljson.Member{Name: "integrity", Value: canonicaljson.Object{{Name: "valid", Value: len(snapshot.IntegrityReasons) == 0}, {Name: "reasons", Value: stringValues(snapshot.IntegrityReasons)}}}, canonicaljson.Member{Name: "identities", Value: identityValue(snapshot.Handoff.Identity)}, canonicaljson.Member{Name: "documents", Value: documents}, canonicaljson.Member{Name: "validations", Value: validations}, canonicaljson.Member{Name: "attempts", Value: attempts}, canonicaljson.Member{Name: "artifacts", Value: snapshotArtifacts(snapshot)}, canonicaljson.Member{Name: "next_transition_authorized", Value: false}, canonicaljson.Member{Name: "control_recommendation", Value: canonicaljson.Object{{Name: "action", Value: action}, {Name: "reason", Value: reason}, {Name: "argv", Value: []canonicaljson.Value{}}}})
	bytes, err := canonicaljson.Marshal(inspection)
	if err != nil {
		return InspectResult{}, schemaError("inspection", err)
	}
	return InspectResult{Bytes: bytes, AppendLF: true}, nil
}

func deriveState(snapshot Snapshot) (string, []string, []string) {
	if len(snapshot.IntegrityReasons) > 0 {
		return "integrity_conflict", nil, nil
	}
	switch snapshot.HeadState {
	case "cancelled", "superseded", "abandoned_unknown":
		return snapshot.HeadState, nil, nil
	}
	if snapshot.Terminal != nil {
		policy, evidence := validateReportedPolicy(snapshot)
		if len(policy) > 0 || len(evidence) > 0 {
			return "conflict", policy, evidence
		}
		return snapshot.Terminal.Outcome, policy, evidence
	}
	if snapshot.Start != nil {
		if snapshot.Start.Outcome == "started" {
			return "started", nil, nil
		}
		return "start_failed", nil, nil
	}
	return "ready", nil, nil
}
func validateReportedPolicy(snapshot Snapshot) ([]string, []string) {
	if snapshot.Terminal == nil {
		return nil, nil
	}
	policy := []string{}
	evidence := []string{}
	rounds := int64(0)
	if value, ok := objectMember(snapshot.Terminal.Value, "rounds_used"); ok {
		rounds, _ = value.(int64)
	}
	if rounds > snapshot.Handoff.MaxRounds {
		policy = append(policy, "rounds_used exceeds handoff budget")
	}
	if reviewValue, ok := objectMember(snapshot.Terminal.Value, "review"); ok {
		review, _ := reviewValue.(canonicaljson.Object)
		if open, found := objectMember(review, "open_actionable_findings"); found {
			if list, ok := open.([]canonicaljson.Value); ok && len(list) > 0 {
				policy = append(policy, "review has open actionable findings")
			}
		}
	}
	if forbidden, ok := objectMember(snapshot.Terminal.Value, "forbidden_effects_observed"); ok {
		if list, ok := forbidden.([]canonicaljson.Value); ok && len(list) > 0 {
			policy = append(policy, "forbidden effects were reported")
		}
	}
	if gaps, ok := objectMember(snapshot.Terminal.Value, "evidence_gaps"); ok {
		if list, ok := gaps.([]canonicaljson.Value); ok && len(list) > 0 {
			evidence = append(evidence, "terminal result reports evidence gaps")
		}
	}
	type expectedVerifier struct {
		exit                         int64
		captureStdout, captureStderr bool
	}
	expected := map[string]expectedVerifier{}
	if verifiers, ok := objectMember(snapshot.Handoff.Value, "verifiers"); ok {
		for _, item := range verifiers.([]canonicaljson.Value) {
			object := item.(canonicaljson.Object)
			exitValue, _ := objectMember(object, "expected_exit")
			exit, _ := exitValue.(int64)
			evidenceValue, _ := objectMember(object, "evidence")
			evidenceObject := evidenceValue.(canonicaljson.Object)
			expected[objectString(object, "id")] = expectedVerifier{exit: exit, captureStdout: objectBool(evidenceObject, "capture_stdout"), captureStderr: objectBool(evidenceObject, "capture_stderr")}
		}
	}
	artifactIDs := map[string]bool{}
	if artifacts, ok := objectMember(snapshot.Terminal.Value, "artifacts"); ok {
		for _, item := range artifacts.([]canonicaljson.Value) {
			artifact := item.(canonicaljson.Object)
			artifactIDs[objectString(artifact, "artifact_id")] = true
			if objectString(artifact, "kind") == "withheld" {
				evidence = append(evidence, "terminal result contains withheld evidence")
			}
		}
	}
	seen := map[string]bool{}
	if results, ok := objectMember(snapshot.Terminal.Value, "verifier_results"); ok {
		for _, item := range results.([]canonicaljson.Value) {
			object := item.(canonicaljson.Object)
			id := objectString(object, "verifier_id")
			exitValue, _ := objectMember(object, "exit")
			exit, _ := exitValue.(int64)
			want, found := expected[id]
			if !found || exit != want.exit {
				policy = append(policy, "verifier result does not match handoff: "+id)
			}
			if found {
				for _, capture := range []struct {
					required bool
					field    string
				}{{want.captureStdout, "stdout_artifact_id"}, {want.captureStderr, "stderr_artifact_id"}} {
					value, _ := objectMember(object, capture.field)
					artifactID, hasArtifact := value.(string)
					if capture.required && (!hasArtifact || !artifactIDs[artifactID]) {
						evidence = append(evidence, "verifier result is missing captured evidence: "+id+"/"+capture.field)
					}
					if !capture.required && value != nil {
						evidence = append(evidence, "verifier result reports undeclared captured evidence: "+id+"/"+capture.field)
					}
				}
			}
			seen[id] = true
		}
	}
	for id := range expected {
		if !seen[id] {
			policy = append(policy, "missing verifier result: "+id)
		}
	}
	return uniqueSorted(policy), uniqueSorted(evidence)
}
func uniqueSorted(values []string) []string {
	sort.Strings(values)
	result := values[:0]
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
func redactHandoff(value canonicaljson.Object) canonicaljson.Object {
	replyValue, _ := objectMember(value, "reply_capability")
	reply, _ := replyValue.(canonicaljson.Object)
	reply = replaceObjectMember(reply, "secret", "[REDACTED]")
	return replaceObjectMember(value, "reply_capability", reply)
}
func redactCapabilityProof(value canonicaljson.Object) canonicaljson.Object {
	proofValue, ok := objectMember(value, "capability_proof")
	if !ok {
		return value
	}
	proof, ok := proofValue.(canonicaljson.Object)
	if !ok {
		return value
	}
	proof = replaceObjectMember(proof, "value", "[REDACTED]")
	return replaceObjectMember(value, "capability_proof", proof)
}
func documentView(locator, digest string, bytes []byte, content canonicaljson.Value) canonicaljson.Object {
	return canonicaljson.Object{{Name: "locator", Value: locator}, {Name: "sha256", Value: digest}, {Name: "size_bytes", Value: int64(len(bytes))}, {Name: "format_version", Value: int64(1)}, {Name: "schema_version", Value: int64(1)}, {Name: "content", Value: content}}
}
func identityValue(value identity) canonicaljson.Object {
	return canonicaljson.Object{{Name: "activity_id", Value: value.ActivityID}, {Name: "run_id", Value: value.RunID}, {Name: "handoff_id", Value: value.HandoffID}, {Name: "start_receipt_id", Value: value.StartReceiptID}, {Name: "terminal_result_id", Value: value.TerminalResultID}, {Name: "publication_key", Value: value.PublicationKey}, {Name: "activity_key", Value: value.ActivityKey}, {Name: "created_at_utc", Value: value.CreatedAtUTC}}
}
func snapshotArtifacts(snapshot Snapshot) []canonicaljson.Value {
	if snapshot.Terminal == nil {
		return []canonicaljson.Value{}
	}
	value, ok := objectMember(snapshot.Terminal.Value, "artifacts")
	if !ok {
		return []canonicaljson.Value{}
	}
	artifacts, ok := value.([]canonicaljson.Value)
	if !ok {
		return []canonicaljson.Value{}
	}
	return artifacts
}
func recommendation(state string) (string, string) {
	switch state {
	case "ready":
		return "human_start_recipient", "The immutable handoff is ready for the expected human-started recipient."
	case "started":
		return "wait_for_recipient", "Wait for the recipient terminal result."
	case "cancelled", "superseded":
		return "none", "The handoff is closed."
	default:
		return "separate_result_control", "A separate control must inspect the preserved evidence before any next transition."
	}
}
func stringValues(values []string) []canonicaljson.Value {
	result := make([]canonicaljson.Value, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
func nullable(value string) canonicaljson.Value {
	if value == "" {
		return nil
	}
	return value
}
