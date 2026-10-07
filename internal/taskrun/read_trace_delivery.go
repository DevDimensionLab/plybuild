package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func (r *traceReader) delivery(s workflowState) {
	if s.Result.Delivery == nil {
		r.problem("trace_delivery_unavailable", "Preserved run has no delivery state")
		return
	}
	candidates := r.deliveryCandidates(s)
	eventRoot := filepath.Join(s.Result.Paths.RunRoot, "delivery", "events")
	entries, entriesErr := inventoryEntries(eventRoot)
	boundPaths, seen := map[string]bool{}, map[string]bool{}
	var previous *string
	for i, event := range s.Result.Delivery.Events {
		predecessor := previous
		previous = ptr(event.Binding.SHA256)
		boundPaths[event.Binding.Locator] = true
		if !key(event.ID) || seen[event.ID] || event.Binding.Locator != filepath.Join(eventRoot, fmt.Sprintf("%08d.json", i+1)) {
			r.problem("trace_delivery_event_unbound", "Duplicate or invalid event slot: "+event.ID)
			continue
		}
		seen[event.ID] = true
		raw, err := r.bound(event.Binding, 8<<20)
		if err != nil {
			r.problem("trace_delivery_event_unavailable", event.ID+": "+err.Error())
			continue
		}
		e := traceEntry(event.ID, event.Kind, s.Result.RunID+":delivery", i+1, event.Binding)
		// Native event records have ordered slots, not a universal predecessor
		// assertion. Only delivery reports explicitly bind their predecessor.
		if event.Kind == "report" {
			var report DeliveryReport
			if err = decode(raw, 8<<20, &report); err == nil {
				err = validateDeliveryReportContent(report)
			}
			if err != nil || report.RunID != s.Result.RunID || report.RequestSHA256 != s.Result.RequestSHA256 || report.SessionID != s.Result.SessionID || report.EventID != event.ID || !equal(report.PreviousEventSHA256, predecessor) {
				r.problem("trace_delivery_report_unbound", "Report identity/predecessor differs: "+event.Binding.Locator)
				continue
			}
			e.Kind, e.Role, e.ActorClaim, e.EvidenceClass = "delivery_report", "agent", s.Request.Delivery.OwnerClaim, "agent_report"
			e.PreviousSHA256, e.Outcome, e.Data = report.PreviousEventSHA256, ptr(report.Phase), traceJSON(report)
			r.run.Entries = append(r.run.Entries, e)
			continue
		}
		if err = inventoryDeliveryEventIdentity(raw, s, event.Kind); err != nil {
			r.problem("trace_delivery_event_unbound", event.ID+": "+err.Error())
			continue
		}
		switch event.Kind {
		case "verification":
			r.verification(s, e, raw, candidates)
		case "human_qa":
			r.deliveryQA(s, e, raw, candidates)
		case "integration":
			r.deliveryIntegration(s, e, raw, candidates)
		case "pull_request":
			var value struct {
				Candidate     string                         `json:"candidate"`
				Receipt       FileBinding                    `json:"receipt"`
				Observed      PullRequestDeliveryObservation `json:"observed"`
				RecordedAtUTC string                         `json:"recorded_at_utc"`
			}
			if err := json.Unmarshal(raw, &value); err != nil {
				r.problem("trace_pr_unbound", e.ID+": invalid PR observation")
				continue
			}
			c, found := candidates[value.Candidate]
			if !found || c.PullRequest == nil || c.PullRequest.SHA256 != e.Source.SHA256 || value.Observed.HeadOID != c.OID {
				r.problem("trace_pr_unbound", e.ID+": PR observation does not bind the qualified candidate")
				continue
			}
			e.Kind, e.Role, e.ActorClaim, e.EvidenceClass = "delivery_pull_request", "ply", "ply PR delivery", "remote_effect_observation"
			e.Candidate, e.Outcome, e.Data = traceCandidate(c), ptr("delivered"), traceJSON(value)
			e.RegisteredAtUTC = r.time(value.RecordedAtUTC, e.ID)
			r.run.Entries = append(r.run.Entries, e)
		}
	}
	if !equal(previous, s.Result.Delivery.LastEventSHA256) {
		r.problem("trace_delivery_chain_unbound", "Last event digest differs from captured delivery state")
	}
	if entriesErr != nil && !os.IsNotExist(entriesErr) {
		r.problem("trace_delivery_store_unavailable", eventRoot+": "+entriesErr.Error())
	}
	for _, item := range entries {
		path := filepath.Join(eventRoot, item.Name())
		if !strings.HasPrefix(item.Name(), ".publish-") && !boundPaths[path] {
			r.problem("trace_delivery_event_orphaned", "Preserved event has no captured state binding: "+path)
		}
	}
	r.checkEntries(eventRoot, entries, entriesErr)
	if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
		r.problem("trace_delivery_attempt_unresolved", "Preserved delivery attempt has no settled effect: "+a.ID)
		if !seen[a.ID] && a.Kind == "verification" {
			path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", a.ID, "attempt.json")
			if a.Path != filepath.Dir(path) {
				r.problem("trace_delivery_attempt_unbound", "Reserved attempt path differs: "+a.ID)
				return
			}
			raw, err := r.bound(FileBinding{path, a.InputSHA256}, 256<<10)
			if err != nil {
				r.problem("trace_delivery_attempt_unavailable", path+": "+err.Error())
				return
			}
			e := traceEntry(a.ID, "verification_reserved", s.Result.RunID+":attempts", 1, FileBinding{path, a.InputSHA256})
			var receipt deliveryVerificationReceipt
			if decode(raw, 256<<10, &receipt) != nil || receipt.RunID != s.Result.RunID || receipt.RequestSHA256 != s.Result.RequestSHA256 || receipt.AttemptID != a.ID || receipt.CandidateOID != a.CandidateOID || receipt.CandidateTree != a.CandidateTree {
				r.problem("trace_delivery_attempt_unbound", "Reserved verification identity differs: "+a.ID)
				return
			}
			e.Role, e.ActorClaim = "ply", "ply"
			e.Candidate = &TraceCandidate{Key: receipt.AttemptID, OID: receipt.CandidateOID, Tree: receipt.CandidateTree}
			e.Data = traceJSON(map[string]any{"attempt_id": a.ID, "state": "unresolved", "meaning": "Reservation does not prove a verifier process executed or finished."})
			r.run.Entries = append(r.run.Entries, e)
		}
	}
}

func (r *traceReader) deliveryCandidates(s workflowState) map[string]DeliveryCandidate {
	out := map[string]DeliveryCandidate{}
	for i, c := range s.Result.Delivery.Candidates {
		_, duplicate := out[c.Key]
		registered := false
		for _, actual := range r.basis.Registry.TaskResults {
			registered = registered || actual.ID == c.TaskResult.ID && equal(actual, c.TaskResult)
		}
		if duplicate || c.Number != i+1 || !key(c.Key) || !registered || c.TaskResult.TaskID != r.basis.Task.ID || c.TaskResult.ResultOID != c.OID || c.TaskResult.ResultTree != c.Tree || c.TaskResult.TechnicalGate != "passed" || c.Handoff.ID != c.TaskResult.HandoffID || c.Handoff.Locator != c.TaskResult.HandoffLocator || c.Handoff.SHA256 != c.TaskResult.HandoffSHA256 || c.Verification.Locator != filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", c.Key, "verification.json") {
			r.problem("trace_candidate_unbound", "Qualified generation differs from its native TaskResult or verification: "+c.Key)
			continue
		}
		if !r.deliveryCandidateEvidence(s, c) {
			continue
		}
		out[c.Key] = c
	}
	return out
}

func (r *traceReader) deliveryCandidateEvidence(s workflowState, c DeliveryCandidate) bool {
	// Registry equality and equal commit bytes alone cannot bind a generation:
	// the same commit may be qualified repeatedly, including in another run.
	receipts := 0
	for _, artifact := range c.TaskResult.Artifacts {
		if artifact.ArtifactID == "verification-receipt" && artifact.SHA256 == c.Verification.SHA256 {
			receipts++
		}
	}
	if receipts != 1 {
		r.problem("trace_candidate_unbound", "Native TaskResult does not bind this exact verifier receipt/generation: "+c.Key)
		return false
	}
	raw, err := r.bound(FileBinding{c.Handoff.Locator, c.Handoff.SHA256}, 4<<20)
	var handoff struct {
		SchemaVersion int                      `json:"schema_version"`
		Basis         *workspace.TaskSpecBasis `json:"task_spec_binding"`
		Identity      struct {
			RunID     string `json:"run_id"`
			HandoffID string `json:"handoff_id"`
		} `json:"identity"`
		Delivery struct {
			Mode          string `json:"mode"`
			ParentLocator string `json:"parent_handoff_locator"`
			ParentSHA256  string `json:"parent_handoff_sha256"`
			CandidateKey  string `json:"candidate_key"`
			OwnerClaim    string `json:"owner_claim"`
		} `json:"delivery_binding"`
	}
	if err == nil {
		err = json.Unmarshal(raw, &handoff)
	}
	if err != nil || handoff.SchemaVersion != 3 || !equal(handoff.Basis, r.run.FrozenBasis) || handoff.Identity.RunID != c.TaskResult.RunID || handoff.Identity.HandoffID != c.Handoff.ID || handoff.Delivery.Mode != "candidate" || handoff.Delivery.ParentLocator != s.Result.Handoff.Locator || handoff.Delivery.ParentSHA256 != s.Result.Handoff.SHA256 || handoff.Delivery.CandidateKey != c.Key || handoff.Delivery.OwnerClaim != s.Request.Delivery.OwnerClaim {
		r.problem("trace_candidate_unbound", "Native candidate handoff does not bind this exact owner and generation: "+c.Key)
		return false
	}
	return true
}

func traceCandidate(c DeliveryCandidate) *TraceCandidate {
	return &TraceCandidate{Key: c.Key, Number: c.Number, OID: c.OID, Tree: c.Tree, TaskResultID: ptr(string(c.TaskResult.ID))}
}

func (r *traceReader) registrySource() FileBinding {
	return FileBinding{r.basis.RegistryLocator, "sha256:" + strings.TrimPrefix(r.basis.Registry.RawSHA256, "sha256:")}
}

func (r *traceReader) verification(s workflowState, e TraceEntry, raw []byte, candidates map[string]DeliveryCandidate) {
	var v deliveryVerificationReceipt
	if err := decode(raw, 256<<10, &v); err != nil || v.SchemaVersion != 1 || v.AttemptID != e.ID || v.RunID != s.Result.RunID || v.RequestSHA256 != s.Result.RequestSHA256 || !traceOIDPattern.MatchString(v.CandidateOID) || !traceOIDPattern.MatchString(v.CandidateTree) || len(v.Argv) == 0 || !filepath.IsAbs(v.CWD) {
		r.problem("trace_verification_unbound", "Verification receipt identity/basis differs: "+e.Source.Locator)
		return
	}
	vp := &TraceVerification{AttemptID: v.AttemptID, CandidateOID: v.CandidateOID, CandidateTree: v.CandidateTree, Argv: v.Argv, CWD: v.CWD, Acceptance: v.Acceptance, AcceptanceSnapshot: v.AcceptanceSnapshot, Review: v.Review, Exit: v.Exit, Stdout: v.Stdout, Stderr: v.Stderr, Error: v.Error}
	for _, tm := range []struct {
		value  string
		target **string
	}{{v.StartedAt, &vp.StartedAtUTC}, {v.FinishedAt, &vp.FinishedAtUTC}} {
		if _, err := parseInventoryTime(tm.value); err == nil {
			*tm.target = ptr(tm.value)
		} else {
			r.problem("trace_verification_time_unknown", e.ID+": missing or invalid verifier endpoint")
		}
	}
	if vp.StartedAtUTC != nil && vp.FinishedAtUTC != nil {
		start, _ := time.Parse(time.RFC3339Nano, *vp.StartedAtUTC)
		finish, _ := time.Parse(time.RFC3339Nano, *vp.FinishedAtUTC)
		if finish.Before(start) {
			r.problem("trace_verification_time_conflict", e.ID+": finish precedes start; elapsed time is unknown")
		}
	}
	path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", v.AttemptID)
	// acceptance.sh is a preserved input snapshot; the original entrypoint is
	// mutable across corrections and is deliberately not requalified here.
	vp.InputsBound = v.AcceptanceSnapshot.Locator == filepath.Join(path, "acceptance.sh") && v.AcceptanceSnapshot.SHA256 == v.Acceptance.SHA256
	for _, input := range []FileBinding{v.AcceptanceSnapshot, v.Review} {
		if _, err := r.bound(input, 4<<20); err != nil {
			vp.InputsBound = false
			r.problem("trace_verification_input_unavailable", e.ID+": "+err.Error())
		} else {
			e.Sources = append(e.Sources, input)
		}
	}
	for _, output := range []FileBinding{v.Stdout, v.Stderr} {
		if _, err := r.bound(output, 4<<20); err != nil {
			r.problem("trace_verification_output_unavailable", e.ID+": "+err.Error())
		} else {
			e.Sources = append(e.Sources, output)
		}
	}
	e.Role, e.ActorClaim, e.EvidenceClass, e.Verification = "ply", "ply verifier", "controlled_verification", vp
	e.Candidate = &TraceCandidate{Key: v.AttemptID, OID: v.CandidateOID, Tree: v.CandidateTree}
	e.Outcome = ptr("unknown")
	if v.Exit != nil {
		e.Outcome = ptr("failed")
		if *v.Exit == 0 && v.Error == "" {
			e.Outcome = ptr("passed")
		}
	}
	e.Data = traceJSON(map[string]any{"attempt_id": v.AttemptID, "error": v.Error, "started_at": v.StartedAt, "finished_at": v.FinishedAt})
	c, qualified := candidates[v.AttemptID]
	if qualified {
		bound, err := r.bound(c.Verification, 256<<10)
		if err != nil || hash(bound) != hash(raw) || c.OID != v.CandidateOID || c.Tree != v.CandidateTree || v.Exit == nil || *v.Exit != 0 || v.Error != "" {
			r.problem("trace_candidate_verification_unbound", "Qualified generation does not match its successful verifier event: "+c.Key)
			qualified = false
		} else {
			e.Candidate = traceCandidate(c)
			e.Sources = append(e.Sources, c.Verification, r.registrySource())
		}
	}
	r.run.Entries = append(r.run.Entries, e)
	if !qualified {
		return
	}
	r.nativeAlias(c.TaskResult.RunID)
	for _, axis := range []string{"reported", "technical"} {
		kind := "task_result_" + axis
		q := traceEntry(string(c.TaskResult.ID)+":"+axis, kind, s.Result.RunID+":candidates", c.Number, r.registrySource())
		q.Sources = append(q.Sources, c.Verification, e.Source)
		q.Candidate, q.Role, q.ActorClaim, q.Recorder = traceCandidate(c), "ply", "ply", c.TaskResult.Recorder.ActorClaim
		q.RegisteredAtUTC = r.time(c.TaskResult.Recorder.RecordedAtUTC, q.ID)
		q.Outcome = ptr(c.TaskResult.TechnicalGate)
		if axis == "reported" {
			q.Role, q.ActorClaim, q.EvidenceClass, q.Outcome = "agent", s.Request.Delivery.OwnerClaim, "agent_report", ptr(c.TaskResult.ReportedOutcome)
		}
		q.Data = traceJSON(map[string]any{"task_result_id": c.TaskResult.ID, "generation": c.Number, "verification_event_id": e.ID, "candidate_key": c.Key, "recorder": c.TaskResult.Recorder})
		traceNative(&q, kind, string(c.TaskResult.ID), c.TaskResult)
		r.run.Entries = append(r.run.Entries, q)
	}
}

func (r *traceReader) time(value, id string) *string {
	if _, err := parseInventoryTime(value); err != nil {
		r.problem("trace_time_unknown", id+": missing or invalid preserved UTC time")
		return nil
	}
	return ptr(value)
}

func (r *traceReader) deliveryQA(s workflowState, e TraceEntry, raw []byte, candidates map[string]DeliveryCandidate) {
	var v struct {
		Kind          string                      `json:"kind"`
		RunID         string                      `json:"run_id"`
		RequestSHA256 string                      `json:"request_sha256"`
		Candidate     string                      `json:"candidate"`
		Record        workspace.TaskHumanQARecord `json:"record"`
		Source        FileBinding                 `json:"source"`
	}
	if decode(raw, 8<<20, &v) != nil {
		r.problem("trace_qa_unbound", "Invalid QA event: "+e.Source.Locator)
		return
	}
	c, ok := candidates[v.Candidate]
	registered := false
	for _, q := range r.basis.Registry.HumanQARecords {
		registered = registered || q.ID == v.Record.ID && equal(q, v.Record)
	}
	q := v.Record
	if !ok || !registered || q.TaskID != r.basis.Task.ID || q.TaskResultID != c.TaskResult.ID || q.ResultOID != c.OID || q.ResultTree != c.Tree {
		r.problem("trace_qa_unbound", "QA event differs from exact native candidate/result binding: "+e.ID)
		return
	}
	e.Candidate, e.Role, e.ActorClaim, e.EvidenceClass, e.Outcome = traceCandidate(c), "human", q.Actor.ActorClaim, "human_attestation", ptr(q.Outcome)
	e.ReportedAtUTC = r.time(q.Actor.CompletedAtUTC, e.ID)
	e.Sources = append(e.Sources, r.registrySource())
	data := map[string]any{"human_qa_id": q.ID, "task_result_id": q.TaskResultID, "actor": q.Actor, "observation": q.Observation, "answer": nil, "meaning": "Human identity and time are explicit attestation claims, not provider authentication or measured effort."}
	// Prefer the immutable native attestation copy; a caller's original private
	// evidence path may legitimately disappear after successful registration.
	for _, source := range q.Evidence {
		binding := FileBinding{source.Locator, source.SHA256}
		bytes, err := r.bound(binding, 1<<20)
		if err != nil {
			r.problem("trace_qa_evidence_unavailable", e.ID+": "+err.Error())
			continue
		}
		e.Sources = append(e.Sources, binding)
		if attestation, err := workflowhandoff.ValidateDeliveryHumanAttestation(bytes, r.run.TaskID, c.TaskResult, q.Outcome); err == nil {
			data["answer"] = attestation.Answer
		}
	}
	e.Data = traceJSON(data)
	traceNative(&e, "human_qa", string(q.ID), q)
	r.run.Entries = append(r.run.Entries, e)
}

func (r *traceReader) deliveryIntegration(s workflowState, e TraceEntry, raw []byte, candidates map[string]DeliveryCandidate) {
	var v struct {
		Candidate string `json:"candidate"`
		Result    struct {
			Completed   bool
			Integration struct {
				Readback struct {
					Value          json.RawMessage
					TaskResult     *workspace.TaskResultRecord
					HumanQA        *workspace.TaskHumanQARecord
					Classification string
				}
			}
		} `json:"result"`
		Error *string `json:"error"`
	}
	if json.Unmarshal(raw, &v) != nil {
		r.problem("trace_integration_unbound", "Invalid integration event: "+e.ID)
		return
	}
	c, ok := candidates[v.Candidate]
	if !ok {
		r.problem("trace_integration_unbound", "Integration has no exact qualified candidate: "+e.ID)
		return
	}
	e.Kind, e.Candidate, e.Role, e.EvidenceClass = "integration_observation", traceCandidate(c), "unknown", "unbound_observation"
	data := map[string]any{"candidate_key": c.Key, "task_result_id": c.TaskResult.ID, "reported_completed": v.Result.Completed, "reported_classification": v.Result.Integration.Readback.Classification, "error": v.Error, "human_qa_id": nil, "integration_result_id": nil}
	// The existing delivery envelope embeds canonicaljson.Object as Go members.
	// Decode that preserved wire form without exporting command streams or the
	// entire integration readback, which are not needed by this projection.
	fields := traceWireObject(v.Result.Integration.Readback.Value)
	var ir workspace.IntegrationResult
	if rawResult := fields["integration_result"]; len(rawResult) > 0 {
		_ = json.Unmarshal(traceWireJSON(rawResult), &ir)
	}
	var registered *workspace.IntegrationResult
	for i := range r.basis.Registry.IntegrationResults {
		actual := &r.basis.Registry.IntegrationResults[i]
		if actual.ID == ir.ID && equal(*actual, ir) {
			registered = actual
		}
	}
	if registered != nil {
		var authority *workspace.IntegrationAuthority
		for i := range r.basis.Registry.IntegrationAuthorities {
			a := &r.basis.Registry.IntegrationAuthorities[i]
			if a.ID == registered.AuthorityID && a.TaskID == r.basis.Task.ID && a.TaskResultID == c.TaskResult.ID && a.Plan.Task.ResultOID == c.OID && a.Plan.Task.ResultTree == c.Tree {
				authority = a
			}
		}
		qaBound := false
		if authority != nil {
			for _, q := range r.basis.Registry.HumanQARecords {
				qaBound = qaBound || q.ID == authority.HumanQARecordID && q.TaskResultID == c.TaskResult.ID && q.ResultOID == c.OID && q.ResultTree == c.Tree && q.Outcome == "pass"
			}
		}
		if authority == nil || !qaBound {
			r.problem("trace_integration_unbound", "Native integration authority differs from exact candidate: "+e.ID)
		} else {
			e.Kind, e.Role, e.ActorClaim, e.EvidenceClass = "integration", "ply", "ply local integration", "local_effect_observation"
			e.Sources = append(e.Sources, r.registrySource())
			e.RegisteredAtUTC, e.Outcome = r.time(ir.RecordedAtUTC, e.ID), ptr(ir.Outcome)
			data["human_qa_id"], data["integration_result_id"] = authority.HumanQARecordID, ir.ID
			data["completed"] = v.Result.Completed
			traceNative(&e, "integration", string(ir.ID), ir)
		}
	} else if ir.ID != "" || v.Result.Completed {
		r.problem("trace_integration_unbound", "Integration observation has no matching preserved native result: "+e.ID)
	}
	e.Data = traceJSON(data)
	r.run.Entries = append(r.run.Entries, e)
}

func traceWireObject(raw json.RawMessage) map[string]json.RawMessage {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) == nil {
		return fields
	}
	var members []struct {
		Name  string
		Value json.RawMessage
	}
	if json.Unmarshal(raw, &members) != nil {
		return nil
	}
	fields = map[string]json.RawMessage{}
	for _, member := range members {
		if member.Name == "" || fields[member.Name] != nil {
			return nil
		}
		fields[member.Name] = member.Value
	}
	return fields
}

func traceWireJSON(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`null`)
	}
	if raw[0] == '[' {
		var elements []json.RawMessage
		if json.Unmarshal(raw, &elements) == nil && len(elements) > 0 {
			var first map[string]json.RawMessage
			if json.Unmarshal(elements[0], &first) == nil && first["Name"] != nil && first["Value"] != nil {
				fields := traceWireObject(raw)
				for name, value := range fields {
					fields[name] = traceWireJSON(value)
				}
				return traceJSON(fields)
			}
			for i, element := range elements {
				elements[i] = traceWireJSON(element)
			}
			return traceJSON(elements)
		}
	}
	if raw[0] == '{' {
		var fields map[string]json.RawMessage
		if json.Unmarshal(raw, &fields) == nil {
			for name, value := range fields {
				fields[name] = traceWireJSON(value)
			}
			return traceJSON(fields)
		}
	}
	// Strict parsing has already excluded duplicate keys and non-JSON values.
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return json.RawMessage(`null`)
	}
	return raw
}
