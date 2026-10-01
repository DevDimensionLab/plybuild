package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
)

func collectPreview(d Dependencies, root, id string) (CollectPreview, journal, error) {
	j, e := readJournal(root, id)
	p := CollectPreview{Envelope: env("collect-preview"), RunID: id, Reasons: []Reason{}, NextArgv: []string{}}
	if e != nil {
		return p, j, e
	}
	if e = validateRunEvidence(d, j); e != nil {
		return p, j, e
	}
	if record, err := recordedResult(d, j); err != nil {
		return p, j, err
	} else if record != nil {
		j.Result.Collection = historicalCollection(record)
		j.Result.ObservedTarget = historicalTarget(record)
		if j.Result.TaskExecution.StatusSHA256 == nil {
			j.Result.TaskExecution.HistoricalQualification = true
		}
	}
	proposed, e := proposedStatus(d, j)
	if e != nil {
		return p, j, e
	}
	r := j.Result
	p.BasisEventSHA256 = r.LastEventSHA256
	p.TaskExecution = r.TaskExecution
	inactive := taskInactive(j)
	if proposed != nil {
		h := digest(proposed)
		p.ProposedTaskStatusSHA256 = &h
		if !statusAlreadyRecorded(j, h) {
			p.TaskExecution = executionFromStatus(*proposed, h)
			inactive = proposed.State == "inactive"
		}
	}
	p.RequestSHA256 = r.RequestSHA256
	p.Process = r.Process
	p.ReportSHA256 = r.Delivery.ReportSHA256
	p.TerminalSHA256 = r.Delivery.TerminalSHA256
	if r.Collection.State == "qualified" {
		finishCollectPreview(&p, j.Request, d.TaskStatusPath)
		return p, j, nil
	}
	add := func(code, detail string) { p.Reasons = append(p.Reasons, Reason{code, detail}) }
	if !inactive {
		add("task_run_task_status_"+p.TaskExecution.State, "Task execution needs a fresh, unambiguous human observation of inactive after the received report. Preview collect with --task-status, then apply its confirmation.")
	}
	if !quiescent(r.Process) || r.Process.State != "exited" {
		add("task_run_process_unknown", "Process and group quiescence are not proven.")
	}
	if r.Process.ExitCode == nil || *r.Process.ExitCode != 0 || r.Process.Signal != nil {
		add("task_run_process_not_successful", "Child exit or signal does not qualify the delivery.")
	}
	if r.Acceptance.State != "started" {
		add("task_run_start_missing", "A bound started receipt is required.")
	}
	if r.Delivery.State != "received" || r.Delivery.TerminalSHA256 == nil {
		add("task_run_terminal_missing", "A complete representable terminal result is required.")
	}
	if r.Delivery.ReportedOutcome == nil || *r.Delivery.ReportedOutcome != "complete" {
		add("task_run_delivery_incomplete", "The semantic report is not complete.")
	}
	if r.Budget.WithinAgreement == nil || !*r.Budget.WithinAgreement {
		add("task_run_budget_unknown_or_exceeded", "Reported active usage does not establish completion within agreement A.")
	}
	if j.Binding != nil {
		facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, j.Binding.Handoff.Locator)
		if err != nil {
			return p, j, err
		}
		if facts.Start == nil || facts.SessionID != "ply:"+id || facts.Acceptance != "started" || r.Acceptance.ReceiptSHA256 == nil || facts.Start.SHA256 != *r.Acceptance.ReceiptSHA256 {
			add("task_run_receipt_binding", "Started receipt or session differs.")
		}
		if facts.Terminal == nil || r.Delivery.TerminalSHA256 == nil || facts.Terminal.SHA256 != *r.Delivery.TerminalSHA256 {
			add("task_run_terminal_binding", "Terminal digest differs.")
		}
		plan := j.Binding.Preparation.Plan
		x, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(plan.WorktreePath, "refs/heads/"+plan.Branch)
		if err != nil {
			add("task_run_target_unknown", "Current target is unobservable.")
		} else {
			p.Target = &workspace.PlanWorktreeObservation{WorktreeLocator: x.Locator, Ref: x.Ref, OID: x.OID, Tree: x.Tree, GitCommonDir: x.GitCommonDir, ObjectFormat: x.ObjectFormat, RefFormat: x.RefFormat, Symbolic: x.Symbolic, Clean: x.Clean, StatusEntries: x.StatusEntries, InProgress: x.InProgress}
			f := facts.FinalTarget
			if f == nil {
				add("task_run_final_target_missing", "No reported final target evidence is available.")
			} else if !x.Clean || !x.Symbolic || len(x.InProgress) != 0 || x.Locator != f.WorktreeLocator || x.Ref != f.Ref || x.OID != f.OID || x.Tree != f.Tree || x.GitCommonDir != f.GitCommonDir {
				add("task_run_candidate_drift", "Current target differs from the reported clean candidate.")
			}
		}
		inspection, err := workflowhandoff.Inspect(d.Workflow, workflowhandoff.InspectInput{HandoffLocator: j.Binding.Handoff.Locator, Format: "json"})
		if err != nil {
			return p, j, err
		}
		p.InspectionSHA256 = ptr(hash(inspection.Bytes))
	}
	p.Reasons = sortedReasons(p.Reasons)
	finishCollectPreview(&p, j.Request, d.TaskStatusPath)
	return p, j, nil
}
func finishCollectPreview(p *CollectPreview, r Request, statusPath string) {
	p.Confirmation = ptr(confirmation(*p))
	p.NextArgv = []string{r.Runtime.PlyExecutable.Path, "workspace", "task", "run", "collect", p.RunID, "--apply", "--confirm", *p.Confirmation}
	if statusPath != "" {
		p.NextArgv = append(p.NextArgv, "--task-status", statusPath)
	}
}
func PreviewCollect(d Dependencies, root, id string) (CollectPreview, error) {
	if e := containing(d, root); e != nil {
		return CollectPreview{}, e
	}
	p, _, e := collectPreview(d, root, id)
	return p, e
}
func Collect(d Dependencies, root, id, confirmation string) (Result, error) {
	return collectApply(d, root, id, confirmation, false)
}
func collectApply(d Dependencies, root, id, confirm string, parent bool) (Result, error) {
	if e := containing(d, root); e != nil {
		return Result{}, e
	}
	p, j, e := collectPreview(d, root, id)
	if e != nil {
		return Result{}, e
	}
	if j.Result.Collection.State == "qualified" && d.TaskStatusPath == "" {
		if !parent && (p.Confirmation == nil || *p.Confirmation != confirm) {
			return Result{}, conflict("collect confirmation differs from preview")
		}
		return recoverCollection(d, root, id)
	}
	if !parent && (p.Confirmation == nil || *p.Confirmation != confirm) {
		return Result{}, conflict("collect confirmation differs from preview")
	}
	e = withStore(root, func() error {
		fresh, j, e := collectPreview(d, root, id)
		if e != nil {
			return e
		}
		if !parent && confirmation(fresh) != confirm {
			return conflict("collect bindings changed")
		}
		if e = publishTaskStatus(d, j, fresh.ProposedTaskStatusSHA256); e != nil {
			return e
		}
		// The status is now durable. Re-read without proposing it against a later tip.
		d.TaskStatusPath = ""
		fresh, j, e = collectPreview(d, root, id)
		if e != nil {
			return e
		}
		if j.Result.Collection.State == "qualified" {
			persisted, err := readJournal(root, id)
			if err != nil {
				return err
			}
			if persisted.Result.Collection.State == "qualified" {
				return nil
			}
			return appendEvent(d, j.Request, "collection", map[string]any{"collection": j.Result.Collection, "observed_target": j.Result.ObservedTarget, "reasons": []Reason{}})
		}
		if e = recoverAcceptance(d, j); e != nil {
			return e
		}
		fresh, j, e = collectPreview(d, root, id)
		if e != nil {
			return e
		}
		// Only complete bookkeeping of preserved bytes, never a launch attempt.
		if j.Result.Delivery.State == "missing" {
			raw, err := readFile(filepath.Join(runPaths(j.Request).RunRoot, "reports", "accepted.json"), 4<<20, true)
			if err == nil {
				report, err := validateReport(raw)
				if err != nil {
					return err
				}
				if err = publishReport(d, j, report, raw); err != nil {
					return err
				}
				fresh, j, e = collectPreview(d, root, id)
				if e != nil {
					return e
				}
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		c := Collection{State: "blocked"}
		if len(fresh.Reasons) == 0 && j.Binding != nil {
			raw, e := readFile(filepath.Join(runPaths(j.Request).RunRoot, "reports", "accepted.json"), 4<<20, true)
			if e != nil {
				return e
			}
			report, e := validateReport(raw)
			if e != nil {
				return e
			}
			if report.TechnicalAssessment.Gate != "passed" && report.TechnicalAssessment.Gate != "good_enough_with_known_debt" {
				fresh.Reasons = append(fresh.Reasons, Reason{"task_run_technical_gate", "Technical assessment does not qualify."})
			} else {
				draftPath := filepath.Join(runPaths(j.Request).TempRoot, "task-result-draft.json")
				draft, e := readFile(draftPath, 1<<20, true)
				if os.IsNotExist(e) {
					technical, _ := Canonical(report.TechnicalAssessment)
					draft, _, e = workflowhandoff.BuildTaskRunResult(d.Workflow, j.Binding.Handoff.Locator, id, technical)
					if e == nil {
						e = d.writeOnce(draftPath, draft)
					}
				}
				if e != nil {
					return e
				}
				c.TaskResultDraftSHA256 = ptr(hash(draft))
				if e = d.fault("before_task_result"); e != nil {
					return e
				}
				w := d.Workspace
				w.HandoffEvidence = workflowhandoff.NewTaskHandoffEvidenceReader(d.Workflow)
				record, err := workspace.RecordTaskResult(w, workspace.TaskResultRecordInput{TaskID: j.Binding.Preparation.Plan.TaskID, File: draftPath})
				if err == nil {
					c.State = "qualified"
					c.TaskResultID = ptr(string(record.Record.ID))
					if e = d.fault("after_task_result"); e != nil {
						return e
					}
				} else {
					fresh.Reasons = append(fresh.Reasons, Reason{"task_run_result_rejected", "Existing TaskResult validation rejected the return: " + err.Error()})
				}
			}
		}
		if !quiescent(j.Result.Process) {
			c.State = "unknown"
		}
		fresh.Reasons = sortedReasons(fresh.Reasons)
		// Avoid a new journal event on an identical, already preserved collection.
		if equal(j.Result.Collection, c) && equal(j.Result.ObservedTarget, fresh.Target) && equal(j.Result.Reasons, fresh.Reasons) {
			return nil
		}
		return appendEvent(d, j.Request, "collection", map[string]any{"collection": c, "observed_target": fresh.Target, "reasons": fresh.Reasons})
	})
	out, readErr := Show(d, root, id)
	if e != nil {
		return out, e
	}
	return out, readErr
}

// recordedResult finds a publication from its frozen bytes before checking a new
// candidate. A crash after the registry write must not mint another result, nor
// invalidate the historical result when later work has changed the target.
func recordedResult(d Dependencies, j journal) (*workspace.TaskResultRecord, error) {
	path := filepath.Join(runPaths(j.Request).TempRoot, "task-result-draft.json")
	raw, e := readFile(path, 1<<20, true)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	if !quiescent(j.Result.Process) || j.Result.Process.State != "exited" || j.Result.Process.ExitCode == nil || *j.Result.Process.ExitCode != 0 || j.Result.Process.Signal != nil || j.Result.Acceptance.State != "started" || j.Result.Delivery.State != "received" {
		return nil, integrity("frozen TaskResult lacks qualifying process and return facts")
	}
	var found *workspace.TaskResultRecord
	e = workspace.WithTaskSpecSnapshot(d.Workspace, j.Request.WorkspaceRoot, func(s *workspace.TaskSpecSession) error {
		for _, x := range s.Registry.TaskResults {
			if x.PublicationKey == "task-run/"+j.Result.RunID {
				if x.DraftSHA256 != hash(raw) || j.Binding == nil || x.HandoffSHA256 != j.Binding.Handoff.SHA256 || j.Result.Acceptance.ReceiptSHA256 == nil || x.StartReceiptSHA256 != *j.Result.Acceptance.ReceiptSHA256 || j.Result.Delivery.TerminalSHA256 == nil || x.TerminalResultSHA256 != *j.Result.Delivery.TerminalSHA256 {
					return integrity("recorded TaskResult differs from frozen run evidence")
				}
				copy := x
				found = &copy
			}
		}
		return nil
	})
	return found, e
}
func historicalCollection(record *workspace.TaskResultRecord) Collection {
	return Collection{State: "qualified", TaskResultID: ptr(string(record.ID)), TaskResultDraftSHA256: &record.DraftSHA256}
}
func historicalTarget(record *workspace.TaskResultRecord) *workspace.PlanWorktreeObservation {
	return &workspace.PlanWorktreeObservation{WorktreeLocator: record.SourceLocator, Ref: record.SourceRef, OID: record.ResultOID, Tree: record.ResultTree, GitCommonDir: record.GitCommonDir, Symbolic: true, Clean: true, StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
}

func recoverCollection(d Dependencies, root, id string) (Result, error) {
	e := withStore(root, func() error {
		j, e := readJournal(root, id)
		if e != nil {
			return e
		}
		if j.Result.Collection.State == "qualified" {
			return nil
		}
		record, e := recordedResult(d, j)
		if e != nil {
			return e
		}
		if record == nil {
			return integrity("previously observed result disappeared")
		}
		return appendEvent(d, j.Request, "collection", map[string]any{"collection": historicalCollection(record), "observed_target": historicalTarget(record), "reasons": []Reason{}})
	})
	out, readErr := Show(d, root, id)
	if e != nil {
		return out, e
	}
	return out, readErr
}
