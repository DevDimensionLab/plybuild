package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var traceOIDPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

func (r *traceReader) native(root, id string) {
	r.run.Family = "legacy_task_run"
	path := filepath.Join(storeRoot(root), "runs", id, "request.json")
	raw, err := r.read(path, 1<<20)
	if err != nil {
		r.problem("trace_native_request_unavailable", path+": "+err.Error())
		return
	}
	request, err := parseRequest(raw)
	if err != nil || RunID(request.RequestKey) != id || request.WorkspaceRoot != root {
		r.problem("trace_run_unbound", "Native run request differs from its workspace/run: "+path)
		return
	}
	r.run.Provider, r.run.SessionID, r.run.RequestSHA256 = request.Runtime.Provider, "ply:"+id, digest(request)
	if !r.bindDraft(request.HandoffDraft, request.PreparationID, request.PreparationSHA256, false) {
		return
	}
	r.contractField("agreement", request.Agreement)
	r.contractField("return_policy", request.ReturnPolicy)
	r.contractField("human_authority", request.HumanAuthority)
	binding := FileBinding{path, hash(raw)}
	entry := traceEntry(id+":request", "run_request", id+":request", 1, binding)
	entry.EvidenceClass, entry.ActorClaim = "declared_contract", request.HumanAuthority.ActorClaim
	entry.Data = traceJSON(map[string]any{"request_sha256": r.run.RequestSHA256, "family": r.run.Family, "human_authority": request.HumanAuthority})
	r.run.Entries = append(r.run.Entries, entry)
	j, journalErr := readJournal(root, id)
	if !equal(j.Request, request) {
		r.problem("trace_native_request_changed", "Native request changed during journal reading: "+path)
		r.run.Coverage = "stale"
		return
	}
	if journalErr != nil {
		r.problem("trace_native_chain_unavailable", journalErr.Error())
	}
	if j.Binding != nil {
		r.nativeHandoff(WorkflowHandoff{j.Binding.Handoff.HandoffID, j.Binding.Handoff.Locator, j.Binding.Handoff.SHA256}, request.HandoffDraft)
		bindingPath := filepath.Join(filepath.Dir(path), "binding.json")
		if _, err := r.bound(FileBinding{bindingPath, digest(j.Binding)}, 2<<20); err != nil {
			r.problem("trace_native_binding_unavailable", err.Error())
		}
		if j.Binding.Preparation.ID != request.PreparationID || digest(j.Binding.Preparation) != request.PreparationSHA256 {
			r.problem("trace_native_binding_unbound", "Native binding differs from its exact historical preparation")
			return
		}
	}
	eventRoot := filepath.Join(filepath.Dir(path), "events")
	entries, entriesErr := inventoryEntries(eventRoot)
	for i, event := range j.Events {
		source := FileBinding{filepath.Join(eventRoot, fmt.Sprintf("%06d.json", event.Sequence)), j.Hashes[i]}
		if _, err := r.bound(source, 64<<10); err != nil {
			r.problem("trace_native_event_unavailable", err.Error())
			continue
		}
		kind := "run_" + event.Type
		e := traceEntry(id+":"+fmt.Sprint(event.Sequence), kind, id+":events", event.Sequence, source)
		e.PreviousSHA256, e.RegisteredAtUTC, e.Role, e.ActorClaim = event.PreviousSHA256, ptr(event.RecordedAtUTC), "ply", "ply"
		e.Sources = append(e.Sources, binding)
		e.Data = event.Payload
		traceNative(&e, kind, e.ID, event)
		switch event.Type {
		case "report_received":
			var payload struct {
				ReportSHA256 string `json:"report_sha256"`
			}
			_ = json.Unmarshal(event.Payload, &payload)
			reportBinding := FileBinding{filepath.Join(filepath.Dir(path), "reports", strings.TrimPrefix(payload.ReportSHA256, "sha256:")+".json"), payload.ReportSHA256}
			bytes, err := r.bound(reportBinding, 4<<20)
			var report Report
			if err == nil {
				report, err = validateReport(bytes)
			}
			if err != nil || report.RunID != id || report.RequestSHA256 != r.run.RequestSHA256 || report.SessionID != r.run.SessionID {
				r.problem("trace_native_report_unavailable", "Native semantic report cannot be bound: "+reportBinding.Locator)
			} else {
				e.Sources = append(e.Sources, reportBinding)
				e.Role, e.ActorClaim, e.EvidenceClass, e.Outcome = "agent", "recipient", "agent_report", ptr(report.Outcome)
				e.Data = traceJSON(report)
			}
		case "acceptance_received":
			var payload struct {
				ClaimSHA256 string `json:"claim_sha256"`
			}
			_ = json.Unmarshal(event.Payload, &payload)
			source := FileBinding{filepath.Join(filepath.Dir(path), "claims", strings.TrimPrefix(payload.ClaimSHA256, "sha256:")+".json"), payload.ClaimSHA256}
			bytes, err := r.bound(source, 256<<10)
			if err == nil {
				var claim Acceptance
				if claim, err = parseAcceptance(bytes); err == nil {
					e.Role, e.ActorClaim, e.EvidenceClass, e.Outcome = "agent", "recipient", "agent_report", ptr(claim.Acceptance)
					e.Sources = append(e.Sources, source)
					e.Data = traceJSON(map[string]any{"acceptance": claim.Acceptance, "runtime_claim": claim.RuntimeClaim})
				}
			}
			if err != nil {
				r.problem("trace_native_acceptance_unavailable", err.Error())
			}
		case "task_status_observed":
			var payload struct {
				StatusSHA256 string `json:"status_sha256"`
			}
			_ = json.Unmarshal(event.Payload, &payload)
			source := FileBinding{taskStatusPath(request, payload.StatusSHA256), payload.StatusSHA256}
			bytes, err := r.bound(source, 64<<10)
			if err == nil {
				var status TaskStatus
				if status, err = parseTaskStatus(bytes); err == nil {
					e.Role, e.ActorClaim, e.EvidenceClass, e.Outcome = "unknown", status.ActorClaim, "status_claim", ptr(status.State)
					e.ReportedAtUTC = r.time(status.ObservedAtUTC, e.ID)
					e.Sources = append(e.Sources, source)
					e.Data = traceJSON(status)
				}
			}
			if err != nil {
				r.problem("trace_native_status_unavailable", err.Error())
			}
		}
		r.run.Entries = append(r.run.Entries, e)
	}
	if entriesErr != nil && !os.IsNotExist(entriesErr) {
		r.problem("trace_native_store_unavailable", eventRoot+": "+entriesErr.Error())
	}
	r.checkEntries(eventRoot, entries, entriesErr)
}
