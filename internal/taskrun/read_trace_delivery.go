package taskrun

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"path/filepath"
	"strings"
	"time"
)

func (r *traceReader) delivery(s workflowState) {
	if s.Result.Delivery == nil {
		r.problem("trace_delivery_unavailable", "Run has no delivery observations")
		return
	}
	candidates := r.deliveryCandidates(s)
	seen := map[string]bool{}
	var previous *string
	for i, event := range s.Result.Delivery.Events {
		if seen[event.ID] || event.Binding.Locator != deliveryEventPath(s.Result.Paths.RunRoot, event.ID, event.Kind, i+1) {
			r.problem("trace_delivery_event_unbound", "invalid event identity: "+event.ID)
			continue
		}
		seen[event.ID] = true
		raw, err := r.bound(event.Binding, 8<<20)
		if err != nil {
			r.problem("trace_delivery_event_unavailable", err.Error())
			continue
		}
		e := traceEntry(event.ID, event.Kind, s.Result.RunID+":delivery", i+1, event.Binding)
		switch event.Kind {
		case "verification":
			r.verification(s, e, raw, candidates)
		case "report":
			var report DeliveryReport
			if decode(raw, 8<<20, &report) != nil || report.RunID != s.Result.RunID || report.RequestSHA256 != s.Result.RequestSHA256 || report.SessionID != s.Result.SessionID || report.EventID != event.ID || !equal(report.PreviousEventSHA256, previous) {
				r.problem("trace_delivery_report_unbound", event.ID)
				break
			}
			e.Kind, e.Role, e.ActorClaim, e.EvidenceClass = "delivery_report", "agent", s.Request.Delivery.OwnerClaim, "agent_report"
			e.PreviousSHA256, e.Outcome, e.Data = report.PreviousEventSHA256, ptr(report.Phase), traceJSON(report)
			r.run.Entries = append(r.run.Entries, e)
		default:
			e.Data = append(json.RawMessage(nil), raw...)
			r.run.Entries = append(r.run.Entries, e)
		}
		previous = ptr(event.Binding.SHA256)
	}
}
func (r *traceReader) deliveryCandidates(s workflowState) map[string]DeliveryCandidate {
	out := map[string]DeliveryCandidate{}
	for _, c := range s.Result.Delivery.Candidates {
		if r.deliveryCandidateEvidence(s, c) {
			out[c.Key] = c
		}
	}
	return out
}
func (r *traceReader) deliveryCandidateEvidence(s workflowState, c DeliveryCandidate) bool {
	receipt := c.TaskResult.VerificationReceipt
	if receipt == nil || receipt.Locator != c.Verification.Locator || receipt.SHA256 != c.Verification.SHA256 || c.TaskResult.RunID != s.Result.RunID || c.OID != c.TaskResult.ResultOID || c.Tree != c.TaskResult.ResultTree {
		r.problem("trace_candidate_unbound", c.Key)
		return false
	}
	found := false
	for _, record := range r.basis.Registry.TaskResults {
		if record.ID == c.TaskResult.ID {
			found = equal(record, c.TaskResult)
		}
	}
	if !found {
		r.problem("trace_candidate_unbound", "record differs from registry: "+c.Key)
		return false
	}
	if err := workspace.ValidateTaskVerificationReceipt(r.d.Workspace, c.TaskResult); err != nil {
		r.problem("trace_candidate_unbound", err.Error())
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
	if decode(raw, 256<<10, &v) != nil || !deliveryVerificationVersion(v) || v.AttemptID != e.ID || v.RunID != s.Result.RunID || v.RequestSHA256 != s.Result.RequestSHA256 || !filepath.IsAbs(v.CWD) {
		r.problem("trace_verification_unbound", e.ID)
		return
	}
	for _, binding := range []FileBinding{v.AcceptanceSnapshot, v.Stdout, v.Stderr} {
		if _, err := r.bound(binding, 64<<20); err != nil {
			r.problem("trace_verification_unbound", err.Error())
			return
		}
	}
	e.Kind, e.Role, e.EvidenceClass = "verification", "executor", "observed_command"
	e.Data = append(json.RawMessage(nil), raw...)
	outcome := "unknown"
	if v.Exit != nil {
		outcome = "failed"
		if *v.Exit == 0 && v.Error == "" {
			outcome = "passed"
		}
	}
	e.Outcome = &outcome
	e.RegisteredAtUTC = r.time(v.FinishedAt, e.ID)
	if c, ok := candidates[v.AttemptID]; ok {
		e.Candidate = traceCandidate(c)
	}
	r.run.Entries = append(r.run.Entries, e)
}
func (r *traceReader) time(value, id string) *string {
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		r.problem("trace_timestamp_invalid", id)
		return nil
	}
	return ptr(value)
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
