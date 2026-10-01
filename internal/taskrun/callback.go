package taskrun

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"os"
	"path/filepath"
	"strings"
)

func callbackContext(d Dependencies, root, id string) (journal, Context, error) {
	var ctx Context
	j, e := readJournal(root, id)
	if e != nil {
		return j, ctx, e
	}
	if j.Binding == nil || j.Result.Launch.Attempts != 1 {
		return j, ctx, conflict("run has no bound launch intent")
	}
	cp := d.ContextPath()
	expected := filepath.Join(runPaths(j.Request).TempRoot, "context.json")
	if cp != expected {
		return j, ctx, conflict("callback requires its bound private context")
	}
	if e = readValue(cp, 1<<20, &ctx); e != nil {
		return j, ctx, e
	}
	want := Context{env("context"), id, digest(j.Request), "ply:" + id, root, digest(j.Binding), j.Request.Runtime.PlyExecutable}
	if !equal(ctx, want) {
		return j, ctx, conflict("callback context differs from binding")
	}
	cwd, e := d.CWD()
	if e != nil {
		return j, ctx, e
	}
	if cwd != j.Binding.Preparation.Plan.WorktreePath {
		return j, ctx, conflict("callback cwd differs from the bound Task worktree")
	}
	if e = physical(cwd, false); e != nil {
		return j, ctx, e
	}
	if d.Executable == nil {
		return j, ctx, conflict("callback executable observer is unavailable")
	}
	actual, e := d.Executable()
	if e != nil {
		return j, ctx, e
	}
	actual, e = filepath.EvalSymlinks(actual)
	if e != nil || actual != ctx.PlyExecutable.Path {
		return j, ctx, conflict("callback is not the bound Ply executable")
	}
	if e = verifyExecutable(ctx.PlyExecutable); e != nil {
		return j, ctx, e
	}
	return j, ctx, nil
}
func claimBinding(id, request, session string, j journal) error {
	if id != j.Result.RunID || request != j.Result.RequestSHA256 || session != "ply:"+id {
		return conflict("callback run, request or session differs")
	}
	return nil
}
func Accept(d Dependencies, root, id, file string) (Result, error) {
	if e := containing(d, root); e != nil {
		return Result{}, e
	}
	raw, e := readFile(file, 256<<10, true)
	if e != nil {
		return Result{}, e
	}
	var claim Acceptance
	if e = decode(raw, 256<<10, &claim); e != nil {
		return Result{}, e
	}
	if e = checkEnvelope(claim.Envelope, "acceptance"); e != nil {
		return Result{}, e
	}
	canonical, _ := Canonical(claim)
	ch := hash(canonical)
	e = withStore(root, func() error {
		j, _, e := callbackContext(d, root, id)
		if e != nil {
			return e
		}
		if e = claimBinding(claim.RunID, claim.RequestSHA256, claim.SessionID, j); e != nil {
			return e
		}
		r := j.Request
		c := claim.RuntimeClaim
		if c.RuntimeID != r.Runtime.Provider || c.ModelID != r.Runtime.Model || c.ProfileID != r.Runtime.PermissionBinding.ProfileID || c.EffectivePolicySHA256 != r.Runtime.PermissionBinding.EffectivePolicySHA256 || c.NativeSessionID != nil && !plain(*c.NativeSessionID, 1, 256) {
			return conflict("runtime, model or effective policy claim differs from request")
		}
		if _, e = runtimeBindings(r.Runtime); e != nil {
			return e
		}
		for _, ev := range j.Events {
			if ev.Type == "acceptance_received" {
				var p struct {
					ClaimSHA256 string `json:"claim_sha256"`
				}
				_ = json.Unmarshal(ev.Payload, &p)
				if p.ClaimSHA256 == ch {
					return nil
				}
				return conflict("accepted claim is immutable")
			}
		}
		path := filepath.Join(runPaths(r).RunRoot, "claims", strings.TrimPrefix(ch, "sha256:")+".json")
		draftPath := filepath.Join(runPaths(r).TempRoot, "start-"+strings.TrimPrefix(ch, "sha256:")+".json")
		draft, e := readFile(draftPath, 256<<10, true)
		if os.IsNotExist(e) {
			draft, e = workflowhandoff.BuildTaskRunStart(d.Workflow, j.Binding.Handoff.Locator, r.HumanAuthority.ActorClaim, r.HumanAuthority.StartSurface, claim.SessionID, c.RuntimeID, c.ModelID, claim.Sandbox, claim.Issues, claim.Acceptance)
		}
		if e != nil {
			return invalid(e.Error())
		}
		if e = d.writeOnce(path, canonical); e != nil {
			return e
		}
		if e = d.writeOnce(draftPath, draft); e != nil {
			return e
		}
		start, e := workflowhandoff.SubmitStart(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: j.Binding.Handoff.Locator, DraftPath: draftPath})
		if e != nil {
			return conflict(e.Error())
		}
		if e = d.fault("after_submit_start"); e != nil {
			return e
		}
		state := "rejected"
		if claim.Acceptance == "started" {
			state = "started"
		}
		return appendEvent(d, r, "acceptance_received", map[string]any{"claim_sha256": ch, "receipt_sha256": start.SHA256, "state": state})
	})
	out, _ := Show(d, root, id)
	return out, e
}
func SubmitReport(d Dependencies, root, id, file string) (Result, error) {
	if e := containing(d, root); e != nil {
		return Result{}, e
	}
	raw, e := readFile(file, 4<<20, true)
	if e != nil {
		return Result{}, e
	}
	report, e := validateReport(raw)
	if e != nil {
		return Result{}, e
	}
	canonical, _ := Canonical(report)
	rh := hash(canonical)
	e = withStore(root, func() error {
		j, _, e := callbackContext(d, root, id)
		if e != nil {
			return e
		}
		if e = claimBinding(report.RunID, report.RequestSHA256, report.SessionID, j); e != nil {
			return e
		}
		if j.Result.Delivery.ReportSHA256 != nil {
			if *j.Result.Delivery.ReportSHA256 == rh {
				return nil
			}
			return conflict("accepted semantic report is immutable")
		}
		r := j.Request
		run := runPaths(r).RunRoot
		// Freeze the observed final target before accepting the semantic slot.
		// Recovery may use these bytes; it must never observe a newer candidate
		// and attach it to an already received report.
		acceptedPath := filepath.Join(run, "reports", "accepted.json")
		if _, err := os.Lstat(acceptedPath); os.IsNotExist(err) && report.BudgetUsage != nil {
			rounds := report.BudgetUsage.CorrectionRounds
			if report.BudgetUsage.InitialExecutionStarted {
				rounds++
			}
			draftPath := terminalDraftPath(r, rh)
			if _, err = os.Lstat(draftPath); os.IsNotExist(err) {
				draft, buildErr := workflowhandoff.BuildTaskRunTerminal(d.Workflow, j.Binding.Handoff.Locator, canonical, rounds)
				if buildErr == nil {
					if e = d.writeOnce(draftPath, draft); e != nil {
						return e
					}
				}
			}
		}
		// This write-once slot reserves the first schema-valid bound semantic report,
		// including blocked/unknown reports that cannot become a WF terminal result.
		if e = d.writeOnce(filepath.Join(run, "reports", "accepted.json"), canonical); e != nil {
			return e
		}
		if e = d.writeOnce(filepath.Join(run, "reports", strings.TrimPrefix(rh, "sha256:")+".json"), canonical); e != nil {
			return e
		}
		if e = d.fault("after_semantic_report"); e != nil {
			return e
		}
		return publishReport(d, j, report, canonical)
	})
	out, _ := Show(d, root, id)
	return out, e
}
func publishReport(d Dependencies, j journal, report Report, raw []byte) error {
	var terminal *string
	state := "not_representable"
	r := j.Request
	if j.Binding != nil && report.BudgetUsage != nil {
		draftPath := terminalDraftPath(r, hash(raw))
		draft, e := readFile(draftPath, 2<<20, true)

		if e == nil {
			if e = d.writeOnce(draftPath, draft); e != nil {
				return e
			}
			var result workflowhandoff.SubmitResult
			recovered, submitErr := workflowhandoff.RecoverTaskRunTerminal(d.Workflow, j.Binding.Handoff.Locator, draft)
			if submitErr == nil {
				if recovered != nil {
					result = *recovered
				} else {
					result, submitErr = workflowhandoff.SubmitResultDocument(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: j.Binding.Handoff.Locator, DraftPath: draftPath})
				}
			}
			if submitErr == nil {
				terminal = &result.SHA256
				state = "received"
			} else {
				if e = d.writeValue(filepath.Join(runPaths(r).RunRoot, "reports", "terminal-diagnostic.json"), map[string]any{"code": "task_run_terminal_not_representable", "detail": "The existing workflow validator did not accept the semantic report; no qualifying TaskResult was created."}); e != nil {
					return e
				}
			}
		}
	}
	if e := d.fault("after_submit_terminal"); e != nil {
		return e
	}
	return appendEvent(d, r, "report_received", map[string]any{"report_sha256": hash(raw), "terminal_sha256": terminal, "state": state})
}

func validateRunEvidence(d Dependencies, j journal) error {
	if j.Binding == nil {
		return nil
	}
	link, e := workflowhandoff.ReadTaskRunLink(d.Workflow, j.Binding.Handoff.Locator)
	if e != nil {
		return e
	}
	want := Handoff{link.ActivityID, link.RunID, link.HandoffID, link.Locator, link.SHA256}
	if !equal(j.Binding.Handoff, want) {
		return integrity("WF handoff differs from run binding")
	}
	facts, e := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, j.Binding.Handoff.Locator)
	if e != nil {
		return e
	}
	if j.Result.Acceptance.ReceiptSHA256 != nil && (facts.Start == nil || facts.Start.SHA256 != *j.Result.Acceptance.ReceiptSHA256 || facts.SessionID != "ply:"+j.Result.RunID) {
		return integrity("WF receipt differs from acceptance event")
	}
	if j.Result.Delivery.TerminalSHA256 != nil && (facts.Terminal == nil || facts.Terminal.SHA256 != *j.Result.Delivery.TerminalSHA256) {
		return integrity("WF terminal differs from report event")
	}
	for _, ev := range j.Events {
		if ev.Type == "acceptance_received" {
			var p struct {
				ClaimSHA256 string `json:"claim_sha256"`
			}
			_ = json.Unmarshal(ev.Payload, &p)
			raw, e := readFile(filepath.Join(runPaths(j.Request).RunRoot, "claims", strings.TrimPrefix(p.ClaimSHA256, "sha256:")+".json"), 256<<10, true)
			if e != nil || hash(raw) != p.ClaimSHA256 {
				return integrity("accepted claim bytes differ")
			}
			var a Acceptance
			if decode(raw, 256<<10, &a) != nil || claimBinding(a.RunID, a.RequestSHA256, a.SessionID, j) != nil {
				return integrity("accepted claim binding differs")
			}
		}
	}
	if j.Result.Delivery.ReportSHA256 != nil {
		raw, e := readFile(filepath.Join(runPaths(j.Request).RunRoot, "reports", "accepted.json"), 4<<20, true)
		if e != nil || hash(raw) != *j.Result.Delivery.ReportSHA256 {
			return integrity("accepted semantic report slot differs")
		}
	}
	return nil
}

func recoverAcceptance(d Dependencies, j journal) error {
	if j.Binding == nil || j.Result.Acceptance.State != "missing" {
		return nil
	}
	facts, e := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, j.Binding.Handoff.Locator)
	if e != nil {
		return e
	}
	if facts.Start == nil {
		return nil
	}
	entries, e := os.ReadDir(filepath.Join(runPaths(j.Request).RunRoot, "claims"))
	if e != nil {
		return e
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		raw, e := readFile(filepath.Join(runPaths(j.Request).RunRoot, "claims", entry.Name()), 256<<10, true)
		if e != nil {
			return e
		}
		var a Acceptance
		if e = decode(raw, 256<<10, &a); e != nil {
			return e
		}
		if e = claimBinding(a.RunID, a.RequestSHA256, a.SessionID, j); e != nil {
			return e
		}
		h := hash(raw)
		draftPath := filepath.Join(runPaths(j.Request).TempRoot, "start-"+strings.TrimPrefix(h, "sha256:")+".json")
		// SubmitStart's existing immutable equality path executes before freshness.
		start, e := workflowhandoff.SubmitStart(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: j.Binding.Handoff.Locator, DraftPath: draftPath})
		if e == nil && start.SHA256 == facts.Start.SHA256 {
			state := "rejected"
			if a.Acceptance == "started" {
				state = "started"
			}
			return appendEvent(d, j.Request, "acceptance_received", map[string]any{"claim_sha256": h, "receipt_sha256": start.SHA256, "state": state})
		}
	}
	return integrity("WF start exists without its matching preserved claim")
}

func terminalDraftPath(r Request, h string) string {
	return filepath.Join(runPaths(r).TempRoot, "terminal-"+strings.TrimPrefix(h, "sha256:")+".json")
}
