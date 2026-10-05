package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"regexp"
)

var workflowReviewID = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,128}$`)

func workflowParseReview(b []byte) (WorkflowReview, error) {
	var r WorkflowReview
	if e := decode(b, 256<<10, &r); e != nil {
		return r, e
	}
	if r.Envelope != workflowEnv("run-review") || !workflowReviewID.MatchString(r.ReviewID) || r.Round < 0 || r.Round > 3 || !digestPattern.MatchString(r.RequestSHA256) || !digestPattern.MatchString(r.HandoffSHA256) || !digestPattern.MatchString(r.ReportSHA256) || !plain(r.Reviewer, 1, 256) {
		return r, workflowError(2, "invalid review binding")
	}
	if r.Decision != "accepted" && r.Decision != "blocked" && r.Decision != "changes_requested" {
		return r, workflowError(2, "invalid review decision")
	}
	if r.Decision != "accepted" && len(r.Findings) == 0 {
		return r, workflowError(2, "findings are required")
	}
	for _, f := range r.Findings {
		if !plain(f.Observed, 1, 2000) || !plain(f.Expected, 1, 2000) || !plain(f.Acceptance, 1, 2000) {
			return r, workflowError(2, "invalid review finding")
		}
	}
	return r, nil
}
func workflowSeal(d Dependencies, s *workflowState, r WorkflowReview) error {
	if r.Decision == "accepted" {
		record := s.Records[len(s.Records)-1]
		terminal, e := workflowBound(record.Terminal, 2<<20)
		if e != nil {
			return e
		}
		if e = workflowhandoff.ValidateWorkflowRound(d.Workflow, s.Result.Handoff.Locator, terminal); e != nil {
			return e
		}
		recovered, e := workflowhandoff.RecoverTaskRunTerminal(d.Workflow, s.Result.Handoff.Locator, terminal)
		if e != nil {
			return e
		}
		if recovered == nil {
			v, e := workflowhandoff.SubmitResultDocument(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: s.Result.Handoff.Locator, DraftPath: record.Terminal.Locator})
			if e != nil {
				return e
			}
			recovered = &v
		}
		if e = d.fault("workflow_after_terminal_publish"); e != nil {
			return e
		}
		s.Result.FinalReturn.TerminalSHA256 = &recovered.SHA256
	}
	s.Result.FinalReturn.State = r.Decision
	s.Result.FinalReturn.ReportSHA256 = &r.ReportSHA256
	s.Result.Round.State = "sealed"
	s.Phase = "sealed"
	s.Result.NextAction = WorkflowAction{"coordinator", "Report review is complete. TaskResult is not published; provider inactivity is not attested. Prepare separate human QA and integration gates."}
	return nil
}
func WorkflowReviewRun(d Dependencies, root, id, file string) (WorkflowRun, error) {
	if s, e := workflowRead(root, id); e == nil && deliveryRun(s.Request) {
		return WorkflowRun{}, workflowError(4, "delivery owners qualify candidates and record actual human QA; legacy coordinator review cannot control this run")
	}
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	b, e := readFile(file, 256<<10, true)
	if e != nil {
		return WorkflowRun{}, e
	}
	r, e := workflowParseReview(b)
	if e != nil {
		return WorkflowRun{}, e
	}
	canonical, _ := Canonical(r)
	send := false
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if e := workflowFresh(d, *s, true); e != nil {
			return e
		}
		if r.RunID != id || r.RequestSHA256 != s.Result.RequestSHA256 || r.HandoffSHA256 != s.Result.Handoff.SHA256 || r.Reviewer != s.Request.Coordinator.ActorClaim {
			return workflowError(4, "review does not match the authorized run and coordinator")
		}
		for _, record := range s.Records {
			if record.Review == nil {
				continue
			}
			old, e := workflowBound(*record.Review, 256<<10)
			if e != nil {
				return e
			}
			var prior WorkflowReview
			if e = decode(old, 256<<10, &prior); e != nil {
				return e
			}
			if prior.ReviewID == r.ReviewID || record.Number == r.Round {
				if record.Review.SHA256 != hash(canonical) {
					return workflowError(4, "review ID or round already has a different immutable decision")
				}
				if r.Decision == "changes_requested" {
					if s.Phase == "prompt_unknown" || s.Phase == "prompt_attempted" || s.Phase == "correction_reserved" {
						return workflowError(5, "feedback attempt is unknown or was not sent; it will not be replayed")
					}
					return nil
				}
				if s.Phase == "sealed" {
					return nil
				}
				return workflowSeal(d, s, r)
			}
		}
		if s.Result.Round.State != "ready_for_review" || r.Round != s.Result.Round.Number || s.Result.Round.ReportSHA256 == nil || r.ReportSHA256 != *s.Result.Round.ReportSHA256 {
			return workflowError(4, "review requires the current immutable report and a settled round")
		}
		if len(s.Records) == 0 {
			return workflowError(4, "round report is missing")
		}
		record := &s.Records[len(s.Records)-1]
		report, e := workflowReadReport(*record)
		if e != nil {
			return e
		}
		a, e := workflowAgentGet(d, *s, false)
		if e != nil {
			return e
		}
		if !workflowSettled(a) {
			return workflowError(4, "same session must be freshly idle for review")
		}
		workflowObserve(d, s, a)
		if r.Decision == "accepted" {
			terminal, e := workflowBound(record.Terminal, 2<<20)
			if e != nil {
				return e
			}
			technical, _ := Canonical(report.TechnicalAssessment)
			if e = workflowhandoff.ValidateWorkflowPassedAssessment(d.Workflow, s.Result.Handoff.Locator, terminal, technical); e != nil {
				return workflowError(4, e.Error())
			}
			x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
			if e != nil {
				return e
			}
			valid := budgetValid(report.BudgetUsage)
			if report.Outcome != "complete" || report.TechnicalAssessment.Gate != "passed" || valid == nil || !*valid || !workflowNoBlockers(report) || !x.Clean || x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir {
				return workflowError(4, "acceptance requires complete native technical evidence, known usage within A and a clean bound candidate without blockers")
			}
		}
		if r.Decision == "changes_requested" {
			u := report.BudgetUsage
			if !workflowUsageAtLeast(u, s.Result.Budget.Floor) || u == nil || !u.InitialExecutionStarted || u.CorrectionRounds >= 3 || u.ActiveSeconds >= 5400 || u.EnvironmentMeasures >= 2 {
				return workflowError(4, "agreement A is unknown or exhausted; no correction is authorized")
			}
		}
		binding, e := workflowKeep(d, filepathForRound(*s, "review.json"), r)
		if e != nil {
			return e
		}
		record.Review = &binding
		s.Result.Review = WorkflowReviewState{&r.ReviewID, &r.Decision}
		if r.Decision == "changes_requested" {
			floor := *report.BudgetUsage
			floor.CorrectionRounds++
			s.Result.Budget.Floor = floor
			s.Result.Round = WorkflowRound{Number: r.Round + 1, State: "working", ControlID: &r.ReviewID, PreviousReportSHA256: &r.ReportSHA256}
			if e = workflowNewContext(d, s); e != nil {
				return e
			}
			s.Phase = "correction_reserved"
			s.Result.NextAction = WorkflowAction{"recipient", "Address only the coordinator's bounded findings and report the new round with cumulative usage."}
			if e = workflowSave(d, *s); e != nil {
				return e
			}
			if e = d.fault("workflow_after_correction_reservation"); e != nil {
				return e
			}
			send = true
			return nil
		}
		s.Phase = "sealing"
		if e = workflowSave(d, *s); e != nil {
			return e
		}
		if e = d.fault("workflow_before_terminal_publish"); e != nil {
			return e
		}
		return workflowSeal(d, s, r)
	})
	if e == nil && send {
		s, err := workflowRead(root, id)
		if err == nil {
			_, err = workflowCall(d, s.Request, "tab", "rename", s.Result.Transport.TabID, "run "+s.Request.Herdr.TabLabel)
		}
		if err == nil {
			err = workflowPrompt(d, root, id, r.Findings)
		}
		if err != nil {
			e = workflowError(5, "Correction reservation preserved without retry: "+err.Error())
		}
	}
	o, _ := WorkflowShow(d, root, id)
	return o, e
}
