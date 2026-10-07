package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (r *traceReader) workflow(root, id string) {
	reservationPath := workflowIndex(root, id)
	raw, err := r.read(reservationPath, 4<<20)
	var reserved workflowState
	if err == nil {
		err = decode(raw, 4<<20, &reserved)
	}
	expected := filepath.Join(workflowRoot(root), "runs", id)
	if err != nil || workflowID(reserved.Request) != id || reserved.Request.WorkspaceRoot != root || reserved.Result.RunID != id || reserved.Result.RequestSHA256 != digest(reserved.Request) || reserved.Result.SessionID != "ply:"+id || reserved.Result.Paths.RunRoot != expected {
		r.problem("trace_run_unbound", "Workflow reservation identity is unavailable: "+reservationPath)
		return
	}
	delivery := deliveryRun(reserved.Request)
	r.run.Family = "workflow_run_v1"
	if delivery {
		r.run.Family = "goal_execute_v2"
	}
	r.run.Provider, r.run.SessionID, r.run.RequestSHA256 = reserved.Request.Runtime.Provider, reserved.Result.SessionID, reserved.Result.RequestSHA256
	if !r.bindDraft(reserved.Request.HandoffDraft, reserved.Request.PreparationID, reserved.Request.PreparationSHA256, delivery) {
		return
	}
	if delivery {
		r.contractField("delivery", reserved.Request.Delivery)
	} else {
		r.contractField("agreement", reserved.Request.Agreement)
		r.contractField("return_mode", reserved.Request.ReturnMode)
	}
	r.contractField("human_authority", reserved.Request.HumanAuthority)
	binding := FileBinding{reservationPath, hash(raw)}
	entry := traceEntry(id+":request", "run_request", id+":request", 1, binding)
	entry.EvidenceClass = "declared_contract"
	entry.ActorClaim = reserved.Request.HumanAuthority.ActorClaim
	entry.Data = traceJSON(map[string]any{"family": r.run.Family, "request_sha256": r.run.RequestSHA256, "human_authority": reserved.Request.HumanAuthority, "meaning": "The preserved request declares authority; it is not a present runtime observation or an authenticated human action."})
	r.run.Entries = append(r.run.Entries, entry)
	s := reserved
	statePath := filepath.Join(expected, "state.json")
	stateBytes, err := r.read(statePath, 8<<20)
	if err == nil {
		var captured inventoryWorkflowState
		if err = decode(stateBytes, 8<<20, &captured); err == nil {
			saved := captured.state()
			if !equal(saved.Request, reserved.Request) || !equal(saved.Observed, reserved.Observed) || !equal(saved.CodexTrust, reserved.CodexTrust) || !equal(saved.ClaudeTrust, reserved.ClaudeTrust) || saved.Result.RunID != id || saved.Result.RequestSHA256 != r.run.RequestSHA256 || saved.Result.SessionID != r.run.SessionID || saved.Result.Paths.RunRoot != expected {
				err = integrity("preserved workflow state differs from reservation")
			} else {
				s = saved
				sources, recoveryErr := readInventoryRecovery(captured)
				for _, source := range sources {
					if _, readErr := r.bound(source.binding, source.max); readErr != nil {
						r.problem("trace_recovery_source_unavailable", readErr.Error())
					}
				}
				if recoveryErr != nil {
					r.problem("trace_recovery_unavailable", recoveryErr.Error())
				}
			}
		}
	}
	if err != nil {
		r.problem("trace_state_unavailable", statePath+": "+err.Error())
	}
	r.nativeHandoff(s.Result.Handoff, s.Request.HandoffDraft)
	if s.Acceptance != nil {
		r.workflowAcceptance(s)
	}
	if delivery {
		r.deliveryGoal(s)
		r.delivery(s)
	} else {
		r.rounds(s)
	}
}

func (r *traceReader) workflowAcceptance(s workflowState) {
	raw, err := r.bound(*s.Acceptance, 1<<20)
	var a Acceptance
	if err == nil {
		if deliveryRun(s.Request) {
			var delivery DeliveryAcceptance
			err = decode(raw, 1<<20, &delivery)
			a = delivery.Acceptance
			if err == nil && a.Envelope != deliveryEnv("run-acceptance") {
				err = integrity("delivery acceptance envelope differs")
			}
		} else {
			err = decode(raw, 1<<20, &a)
			if err == nil && a.Envelope != workflowEnv("run-acceptance") {
				err = integrity("legacy acceptance envelope differs")
			}
		}
		if err == nil {
			native := a
			native.Envelope = Envelope{"ply.workspace.task-run-acceptance", 2}
			_, err = parseAcceptance(traceJSON(native))
		}
	}
	if err != nil || a.RunID != s.Result.RunID || a.RequestSHA256 != s.Result.RequestSHA256 || a.SessionID != s.Result.SessionID {
		r.problem("trace_acceptance_unavailable", "Preserved acceptance differs from its run: "+s.Acceptance.Locator)
		return
	}
	e := traceEntry(s.Result.RunID+":acceptance", "run_acceptance", s.Result.RunID+":acceptance", 1, *s.Acceptance)
	e.Role, e.ActorClaim, e.EvidenceClass = "agent", "recipient", "agent_report"
	e.Outcome = nonempty(a.Acceptance)
	e.Data = traceJSON(map[string]any{"acceptance": a.Acceptance, "runtime_claim": a.RuntimeClaim, "meaning": "Recipient runtime and permission claims are preserved; provider identity is not human authentication."})
	r.run.Entries = append(r.run.Entries, e)
}

func (r *traceReader) deliveryGoal(s workflowState) {
	binding := s.Request.Delivery.Goal
	r.run.FrozenGoalSource = &binding
	raw, err := r.bound(binding, 8<<20)
	if err == nil {
		r.run.FrozenGoal, err = inventoryDeliveryGoal(raw, *r.run.FrozenBasis, binding.SHA256)
	}
	if err != nil {
		r.problem("trace_frozen_goal_unavailable", binding.Locator+": "+err.Error())
	} else {
		r.contractField("frozen_goal_declaration", traceManifestDeclaration(raw))
	}
	basis := r.run.FrozenBasis
	specBinding := FileBinding{filepath.Join(s.Request.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(basis.Spec.ManifestSHA256, "sha256:")+".json"), basis.Spec.ManifestSHA256}
	spec, specErr := r.bound(specBinding, 8<<20)
	if specErr == nil && r.run.FrozenGoal != nil {
		specErr = inventoryDeliveryGoalOrigin(spec, *basis, *r.run.FrozenGoal)
	}
	if specErr != nil {
		r.problem("trace_frozen_execution_unavailable", specBinding.Locator+": "+specErr.Error())
	} else {
		r.contractField("frozen_execution_declaration", traceManifestDeclaration(spec))
	}
}

func traceManifestDeclaration(raw []byte) map[string]json.RawMessage {
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(raw, &fields)
	out := map[string]json.RawMessage{}
	for _, name := range []string{"kind", "schema_version", "contract_kind", "task_id", "spec_id", "revision", "title", "objective", "requirements", "constraints", "goal_origin", "phases", "parts"} {
		if value, ok := fields[name]; ok {
			out[name] = value
		}
	}
	return out
}

func (r *traceReader) rounds(s workflowState) {
	seen := map[int]bool{}
	boundPaths := map[string]bool{}
	var previousReport *string
	var previousReview *WorkflowReview
	for _, record := range s.Records {
		prefix := filepath.Join(s.Result.Paths.RunRoot, "rounds", workflowRoundName(record.Number))
		if record.Number < 0 || record.Number > 3 || seen[record.Number] || record.Report.Locator != filepath.Join(prefix, "report.json") || record.Terminal.Locator != filepath.Join(prefix, "terminal-draft.json") {
			r.problem("trace_round_unbound", fmt.Sprintf("Round %d has an invalid or duplicate source binding", record.Number))
			continue
		}
		seen[record.Number] = true
		boundPaths[prefix] = true
		raw, err := r.bound(record.Report, 1<<20)
		var report WorkflowReport
		if err == nil {
			report, err = workflowParseReport(raw)
		}
		validIdentity := report.RunID == s.Result.RunID && report.RequestSHA256 == s.Result.RequestSHA256 && report.SessionID == s.Result.SessionID && report.Round == record.Number
		if err != nil || !validIdentity {
			r.problem("trace_round_report_unavailable", "Round report cannot be bound: "+record.Report.Locator)
		} else {
			e := traceEntry(s.Result.RunID+":round:"+fmt.Sprint(record.Number)+":report", "round_report", s.Result.RunID+":rounds", record.Number*2+1, record.Report)
			e.PreviousSHA256 = report.PreviousReportSHA256
			e.Role, e.ActorClaim, e.EvidenceClass = "agent", "recipient", "agent_report"
			e.Outcome = ptr(report.Outcome)
			e.Data = traceJSON(report)
			terminal, terminalErr := r.bound(record.Terminal, 2<<20)
			if terminalErr == nil {
				slot := FileBinding{filepath.Join(prefix, "report-slot.json"), digest(workflowReportSlot{report, terminal, record.TargetSHA256})}
				if _, terminalErr = r.bound(slot, 4<<20); terminalErr == nil {
					e.Sources = append(e.Sources, record.Terminal, slot)
				}
			}
			if terminalErr != nil {
				r.problem("trace_round_slot_unavailable", record.Report.Locator+": "+terminalErr.Error())
			}
			if record.Number == 0 && (report.ControlID != nil || report.PreviousReportSHA256 != nil) || record.Number > 0 && (!equal(report.PreviousReportSHA256, previousReport) || previousReview == nil || previousReview.Decision != "changes_requested" || report.ControlID == nil || *report.ControlID != previousReview.ReviewID || previousReview.Round+1 != record.Number) {
				r.problem("trace_round_chain_unavailable", "Round report predecessor/review binding is missing or different: "+record.Report.Locator)
			}
			r.run.Entries = append(r.run.Entries, e)
		}
		previousReport = ptr(record.Report.SHA256)
		previousReview = nil
		if record.Review == nil {
			continue
		}
		if record.Review.Locator != filepath.Join(prefix, "review.json") {
			r.problem("trace_round_review_unbound", "Review path differs from its round: "+record.Review.Locator)
			continue
		}
		raw, err = r.bound(*record.Review, 256<<10)
		var review WorkflowReview
		if err == nil {
			review, err = workflowParseReview(raw)
		}
		if err != nil || review.RunID != s.Result.RunID || review.RequestSHA256 != s.Result.RequestSHA256 || review.HandoffSHA256 != s.Result.Handoff.SHA256 || review.Round != record.Number || review.ReportSHA256 != record.Report.SHA256 || review.Reviewer != s.Request.Coordinator.ActorClaim {
			r.problem("trace_round_review_unavailable", "Review cannot be bound to its exact report: "+record.Review.Locator)
			continue
		}
		previousReview = &review
		e := traceEntry(review.ReviewID, "round_review", s.Result.RunID+":rounds", record.Number*2+2, *record.Review)
		e.PreviousSHA256 = ptr(review.ReportSHA256)
		// Legacy coordinator claims do not declare whether the reviewer is a
		// person or an agent. A reviewer name must not authenticate a human.
		e.ActorClaim, e.EvidenceClass, e.Outcome = review.Reviewer, "review_claim", ptr(review.Decision)
		e.Data = traceJSON(review)
		r.run.Entries = append(r.run.Entries, e)
	}
	path := filepath.Join(s.Result.Paths.RunRoot, "rounds")
	entries, err := inventoryEntries(path)
	if err != nil && !os.IsNotExist(err) {
		r.problem("trace_round_store_unavailable", path+": "+err.Error())
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") || boundPaths[filepath.Join(path, entry.Name())] {
			continue
		}
		// The current round may have only a context and no report yet.
		p := filepath.Join(path, entry.Name(), "report.json")
		if _, readErr := readFile(p, 1<<20, true); !os.IsNotExist(readErr) {
			r.problem("trace_round_orphaned", "Unbound preserved round report: "+p)
		}
	}
	r.checkEntries(path, entries, err)
}

func (r *traceReader) checkEntries(path string, entries []os.DirEntry, err error) {
	if err != nil {
		return
	}
	check := Inventory{Reasons: []Reason{}}
	inventoryCheckEntries(&check, path, entries)
	for _, reason := range check.Reasons {
		r.problem(reason.Code, reason.Detail)
		r.run.Coverage = "stale"
	}
}
