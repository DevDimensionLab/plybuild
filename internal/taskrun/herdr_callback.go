package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func filepathForRound(s workflowState, name string) string {
	return filepath.Join(s.Result.Paths.RunRoot, "rounds", workflowRoundName(s.Result.Round.Number), name)
}
func workflowCallback(d Dependencies, s workflowState, contextPath string) error {
	if contextPath == "" {
		return workflowError(2, "--context is required")
	}
	if contextPath != s.Result.Paths.Context {
		return workflowError(4, "callback requires the exact current private context")
	}
	if e := workflowFresh(d, s, false); e != nil {
		return e
	}
	cwd, e := d.CWD()
	if e != nil {
		return e
	}
	if cwd != s.Observed.Target.WorktreeLocator {
		return workflowError(4, "callback cwd differs from the exact Task worktree")
	}
	if e = physical(cwd, false); e != nil {
		return e
	}
	actual, e := d.Executable()
	if e == nil {
		actual, e = filepath.EvalSymlinks(actual)
	}
	if e != nil || actual != s.Request.Runtime.PlyExecutable.Path {
		return workflowError(4, "callback is not the bound Ply executable")
	}
	return verifyExecutable(s.Request.Runtime.PlyExecutable)
}
func workflowClaim(s workflowState, id, request, session string) error {
	if id != s.Result.RunID || request != s.Result.RequestSHA256 || session != s.Result.SessionID {
		return workflowError(4, "callback run, request or session differs")
	}
	return nil
}
func WorkflowAccept(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	b, e := readFile(file, 256<<10, true)
	if e != nil {
		return WorkflowRun{}, e
	}
	var a Acceptance
	if e = decode(b, 256<<10, &a); e != nil {
		return WorkflowRun{}, e
	}
	if a.Envelope != workflowEnv("run-acceptance") {
		return WorkflowRun{}, workflowError(2, "invalid workflow acceptance kind or version")
	}
	native := a
	native.Envelope = Envelope{"ply.workspace.task-run-acceptance", 2}
	raw, _ := Canonical(native)
	if _, e = parseAcceptance(raw); e != nil {
		return WorkflowRun{}, e
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if e := workflowCallback(d, *s, contextPath); e != nil {
			return e
		}
		if e := workflowClaim(*s, a.RunID, a.RequestSHA256, a.SessionID); e != nil {
			return e
		}
		canonical, _ := Canonical(a)
		if s.Acceptance != nil && s.Acceptance.SHA256 != hash(canonical) {
			return workflowError(4, "acceptance slot is immutable")
		}
		if s.StartSHA256 != nil {
			return nil
		}
		if s.Result.Round.Number != 0 || s.Result.Transport.AgentSessionID == "" {
			return workflowError(4, "run has no bound initial session")
		}
		if s.Phase != "prompt_attempted" && s.Phase != "prompt_unknown" && s.Phase != "following" {
			return workflowError(4, "initial prompt has not been attempted; coordinator inspection is required before recipient acceptance")
		}
		c := a.RuntimeClaim
		if c.NativeSessionID != nil && *c.NativeSessionID != s.Result.Transport.AgentSessionID {
			return workflowError(4, "native session claim differs from Herdr")
		}
		if a.Acceptance == "started" {
			if c.NativeSessionID == nil {
				return workflowError(4, "positive acceptance requires the bound native session")
			}
			if e := positiveClaim(native, workflowNativeRequest(s.Request)); e != nil {
				return e
			}
			// The recipient's claim cannot establish that Herdr still hosts the
			// selected provider/session. Observe it before granting Task writes.
			// A callback normally runs while the agent is working, not idle.
			fresh, e := workflowAgentGet(d, *s, false)
			if e != nil {
				return e
			}
			workflowObserve(d, s, fresh)
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "start-draft.json")
		draft, err := readFile(path, 1<<20, true)
		if os.IsNotExist(err) && c.RuntimeID != nil && string(a.Sandbox) != "null" {
			model := "unknown"
			if c.ModelID != nil {
				model = *c.ModelID
			}
			draft, err = workflowhandoff.BuildTaskRunStart(d.Workflow, s.Result.Handoff.Locator, s.Request.HumanAuthority.ActorClaim, s.Request.HumanAuthority.StartSurface, a.SessionID, *c.RuntimeID, model, a.Sandbox, a.Issues, a.Acceptance)
		}
		if err == nil {
			err = workflowhandoff.ValidateTaskRunStart(d.Workflow, s.Result.Handoff.Locator, draft)
		}
		if err != nil && a.Acceptance == "started" {
			return workflowError(4, "native runtime acceptance rejected: "+err.Error())
		}
		if err == nil {
			if e := d.writeOnce(path, draft); e != nil {
				return e
			}
			s.StartDraft = &FileBinding{path, hash(draft)}
		}
		binding, e := workflowKeep(d, filepath.Join(s.Result.Paths.RunRoot, "acceptance.json"), a)
		if e != nil {
			return e
		}
		s.Acceptance = &binding
		s.Result.RuntimeFacts = RuntimeFacts{s.Request.Runtime.Model, c.ModelID, "unknown", ptr("recipient_claim"), ptr(binding.SHA256)}
		if c.ModelID != nil {
			if *c.ModelID == s.Request.Runtime.Model {
				s.Result.RuntimeFacts.ModelState = "matches"
			} else {
				s.Result.RuntimeFacts.ModelState = "differs"
			}
		}
		if e = workflowSave(d, *s); e != nil {
			return e
		}
		if e = d.fault("workflow_after_acceptance"); e != nil {
			return e
		}
		if s.StartDraft != nil {
			recovered, e := workflowhandoff.RecoverTaskRunStart(d.Workflow, s.Result.Handoff.Locator, draft)
			if e != nil {
				return e
			}
			if recovered == nil {
				result, e := workflowhandoff.SubmitStart(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: s.Result.Handoff.Locator, DraftPath: path})
				if e != nil {
					return e
				}
				recovered = &result
			}
			s.StartSHA256 = &recovered.SHA256
		}
		if a.Acceptance == "started" {
			s.Result.Round.State = "working"
		} else {
			s.Result.Round.State = "unknown"
			s.Result.FinalReturn.State = "blocked"
			s.Result.NextAction = WorkflowAction{"coordinator", "Runtime acceptance stopped the run; inspect the recipient's issues."}
		}
		return nil
	})
	o, _ := WorkflowShow(d, root, id)
	return o, e
}

func workflowParseReport(b []byte) (WorkflowReport, error) {
	var r WorkflowReport
	if e := decode(b, 1<<20, &r); e != nil {
		return r, e
	}
	if r.Envelope != workflowEnv("round-report") || r.Round < 0 || r.Round > 3 {
		return r, workflowError(2, "invalid report kind, version or round")
	}
	n := r.Report
	n.Envelope = env("report")
	native, _ := Canonical(n)
	if _, e := validateReport(native); e != nil {
		return r, e
	}
	if r.ControlID != nil && !workflowReviewID.MatchString(*r.ControlID) || r.PreviousReportSHA256 != nil && !digestPattern.MatchString(*r.PreviousReportSHA256) {
		return r, workflowError(2, "invalid report predecessor binding")
	}
	return r, nil
}
func workflowUsageAtLeast(used *BudgetUsage, floor BudgetUsage) bool {
	return used != nil && (!floor.InitialExecutionStarted || used.InitialExecutionStarted) && used.CorrectionRounds >= floor.CorrectionRounds && used.ActiveSeconds >= floor.ActiveSeconds && used.EnvironmentMeasures >= floor.EnvironmentMeasures
}

type workflowReportSlot struct {
	Report       WorkflowReport  `json:"report"`
	Terminal     json.RawMessage `json:"terminal"`
	TargetSHA256 string          `json:"target_sha256"`
}

func WorkflowReportRound(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	b, e := readFile(file, 1<<20, true)
	if e != nil {
		return WorkflowRun{}, e
	}
	r, e := workflowParseReport(b)
	if e != nil {
		return WorkflowRun{}, e
	}
	canonical, _ := Canonical(r)
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if e := workflowCallback(d, *s, contextPath); e != nil {
			return e
		}
		if e := workflowClaim(*s, r.RunID, r.RequestSHA256, r.SessionID); e != nil {
			return e
		}
		if r.Round != s.Result.Round.Number || !equal(r.ControlID, s.Result.Round.ControlID) || !equal(r.PreviousReportSHA256, s.Result.Round.PreviousReportSHA256) {
			return workflowError(4, "report does not match the current reserved round")
		}
		if s.Result.Round.ReportSHA256 != nil {
			if *s.Result.Round.ReportSHA256 == hash(canonical) {
				return nil
			}
			return workflowError(4, "round report slot is immutable")
		}
		if s.Acceptance == nil || s.StartSHA256 == nil {
			return workflowError(4, "report requires positive native runtime acceptance")
		}
		var a Acceptance
		claim, e := workflowBound(*s.Acceptance, 256<<10)
		if e != nil {
			return e
		}
		if e = decode(claim, 256<<10, &a); e != nil {
			return e
		}
		if a.Acceptance != "started" {
			return workflowError(4, "negative acceptance cannot report a working round")
		}
		if !workflowUsageAtLeast(r.BudgetUsage, s.Result.Budget.Floor) {
			return workflowError(4, "report usage is unknown or below the cumulative budget floor")
		}
		if s.Result.FinalReturn.State != "pending" {
			return workflowError(4, "run is sealed")
		}
		target, e := workflowTarget(d, *s)
		if e != nil {
			return e
		}
		terminalPath := filepathForRound(*s, "terminal-draft.json")
		// Reserve the report, terminal semantics and candidate as one atomic slot.
		// Later files are deterministic projections, never independent decisions.
		var slot workflowReportSlot
		e = readValue(filepathForRound(*s, "report-slot.json"), 4<<20, &slot)
		var terminal []byte
		if os.IsNotExist(e) {
			terminal, e = workflowhandoff.BuildTaskRunTerminal(d.Workflow, s.Result.Handoff.Locator, canonical, 1+r.BudgetUsage.CorrectionRounds)
		} else if e == nil {
			if !equal(slot.Report, r) || slot.TargetSHA256 != target {
				return workflowError(4, "immutable report reservation or its candidate differs")
			}
			terminal = slot.Terminal
		}
		if e != nil {
			return workflowError(4, "round cannot be represented by the native terminal: "+e.Error())
		}
		if e = workflowhandoff.ValidateWorkflowRound(d.Workflow, s.Result.Handoff.Locator, terminal); e != nil {
			return workflowError(4, "native round validation failed: "+e.Error())
		}
		after, e := workflowTarget(d, *s)
		if e != nil {
			return e
		}
		if after != target {
			return workflowError(4, "candidate changed during report validation")
		}
		if e = d.writeValue(filepathForRound(*s, "report-slot.json"), workflowReportSlot{r, terminal, target}); e != nil {
			return e
		}
		if e = d.writeOnce(terminalPath, terminal); e != nil {
			return e
		}
		report, e := workflowKeep(d, filepathForRound(*s, "report.json"), r)
		if e != nil {
			return e
		}
		record := workflowRecord{r.Round, report, FileBinding{terminalPath, hash(terminal)}, target, nil}
		if e = d.writeValue(filepathForRound(*s, "record.json"), record); e != nil {
			return e
		}
		if e = d.fault("workflow_after_report_slot"); e != nil {
			return e
		}
		s.Records = append(s.Records, record)
		s.Result.Round.ReportPath = &report.Locator
		s.Result.Round.ReportSHA256 = &report.SHA256
		s.Result.Round.State = "report_received"
		s.Result.Budget.Used = r.BudgetUsage
		s.Result.Budget.Floor = *r.BudgetUsage
		s.Result.Budget.WithinAgreement = budgetValid(r.BudgetUsage)
		s.Result.Budget.MeasurementSource = "reported"
		s.Result.NextAction = WorkflowAction{"coordinator", "Follow the same session, then review this immutable report."}
		return nil
	})
	o, _ := WorkflowShow(d, root, id)
	return o, e
}

func workflowReadReport(record workflowRecord) (WorkflowReport, error) {
	b, e := workflowBound(record.Report, 1<<20)
	if e != nil {
		return WorkflowReport{}, e
	}
	return workflowParseReport(b)
}
func workflowNoBlockers(r WorkflowReport) bool {
	var review struct {
		Open []json.RawMessage `json:"open_actionable_findings"`
	}
	if json.Unmarshal(r.Review, &review) != nil {
		return false
	}
	var forbidden, gaps []json.RawMessage
	if json.Unmarshal(r.ForbiddenEffectsObserved, &forbidden) != nil || json.Unmarshal(r.EvidenceGaps, &gaps) != nil {
		return false
	}
	return len(review.Open) == 0 && len(forbidden) == 0 && len(gaps) == 0
}
