package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

var jsonUnmarshal = json.Unmarshal

func ReadWorkflowRequest(file string) (WorkflowRequest, error) {
	var r WorkflowRequest
	b, e := readFile(file, 1<<20, false)
	if e != nil {
		return r, e
	}
	if e = decode(b, 1<<20, &r); e != nil {
		return r, e
	}
	if r.Envelope != workflowEnv("herdr-run-request") || r.ReturnMode != "reviewed_report_only" || r.HumanAuthority.StartSurface != "human_authorized_herdr" || !r.Coordinator.MayRequestChanges || !plain(r.Coordinator.ActorClaim, 1, 256) || !plain(r.Herdr.TabLabel, 1, 80) || !plain(r.Herdr.WorkspaceID, 1, 128) {
		return r, workflowError(2, "invalid workflow request, authority or Herdr binding")
	}
	// Reuse the unchanged native schema on an in-memory projection only.
	b, e = Canonical(workflowNativeRequest(r))
	if e != nil {
		return r, e
	}
	if _, e = parseRequest(b); e != nil {
		return r, e
	}
	if !digestPattern.MatchString(r.Herdr.Executable.SHA256) {
		return r, workflowError(2, "invalid Herdr executable digest")
	}
	if e = verifyExecutable(r.Herdr.Executable); e != nil {
		return r, e
	}
	return r, nil
}
func workflowInitial(r WorkflowRequest, observed Observed) workflowState {
	id := workflowID(r)
	o := WorkflowRun{Envelope: workflowEnv("run"), RunID: id, RequestSHA256: digest(r), SessionID: "ply:" + id, Paths: workflowPaths(r, 0), Transport: WorkflowTransport{WorkspaceID: r.Herdr.WorkspaceID, State: "unknown", Observation: "cached"}, Round: WorkflowRound{State: "awaiting_acceptance"}, Budget: WorkflowBudget{MeasurementSource: "unknown"}, FinalReturn: WorkflowFinal{State: "pending"}, TaskResultState: "not_published", Reasons: []Reason{}, NextAction: WorkflowAction{"recipient", "Accept the bound runtime before target writes, then report this round."}, RuntimeFacts: RuntimeFacts{RequestedModel: r.Runtime.Model, ModelState: "unknown"}}
	return workflowState{Request: r, Observed: observed, Result: o, Phase: "reserved", Records: []workflowRecord{}}
}
func workflowExisting(r WorkflowRequest) (*WorkflowRun, error) {
	s, e := workflowRead(r.WorkspaceRoot, workflowID(r))
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	if !equal(s.Request, r) {
		return nil, workflowError(4, "request key already has different canonical bytes")
	}
	return &s.Result, nil
}
func workflowPreview(d Dependencies, r WorkflowRequest, file string) (WorkflowPreview, error) {
	p := WorkflowPreview{Envelope: workflowEnv("run-preview"), RunID: workflowID(r), RequestSHA256: digest(r), Reasons: []Reason{}, Effects: []string{"Reserve this Task target and one native handoff", "Create one background Herdr tab and one Codex session", "Preserve acceptance, immutable round reports and coordinator review", "Return a reviewed report; do not publish TaskResult or attest provider inactivity"}, Paths: workflowPaths(r, 0)}
	n, e := previewLocked(d, workflowNativeRequest(r), file)
	p.Observed = n.Observed
	if e == nil {
		e = verifyExecutable(r.Herdr.Executable)
	}
	if e != nil {
		p.Reasons = append(p.Reasons, Reason{"workflow_start_blocked", e.Error()})
		return p, e
	}
	p.Confirmation = ptr(digest(struct {
		Request  WorkflowRequest
		Observed Observed
		Paths    WorkflowPaths
		Effects  []string
	}{r, p.Observed, p.Paths, p.Effects}))
	return p, nil
}
func WorkflowPreviewStart(d Dependencies, file string) (any, error) {
	r, e := ReadWorkflowRequest(file)
	if e != nil {
		return nil, e
	}
	if e = containing(d, r.WorkspaceRoot); e != nil {
		return nil, e
	}
	if o, e := workflowExisting(r); o != nil || e != nil {
		if o != nil {
			return WorkflowShow(d, r.WorkspaceRoot, o.RunID)
		}
		return nil, e
	}
	return workflowPreview(d, r, file)
}
func WorkflowStart(d Dependencies, file, confirm string) (WorkflowRun, error) {
	r, e := ReadWorkflowRequest(file)
	if e != nil {
		return WorkflowRun{}, e
	}
	if e = containing(d, r.WorkspaceRoot); e != nil {
		return WorkflowRun{}, e
	}
	if o, e := workflowExisting(r); o != nil || e != nil {
		if o != nil {
			return WorkflowShow(d, r.WorkspaceRoot, o.RunID)
		}
		return WorkflowRun{}, e
	}
	p, e := workflowPreview(d, r, file)
	if e != nil {
		return WorkflowRun{}, e
	}
	if p.Confirmation == nil || confirm != *p.Confirmation {
		return WorkflowRun{}, workflowError(4, "confirmation differs from current preview")
	}
	if os.Getenv("HERDR_ENV") != "1" {
		return WorkflowRun{}, workflowError(4, "apply requires the local Herdr context (HERDR_ENV=1)")
	}
	resolved, e := exec.LookPath("codex")
	if e == nil {
		resolved, e = filepath.EvalSymlinks(resolved)
	}
	if e != nil || resolved != r.Runtime.Executable.Path {
		return WorkflowRun{}, workflowError(4, "local codex executable differs from the bound runtime")
	}
	owned := false
	e = withStore(r.WorkspaceRoot, func() error {
		if o, e := workflowExisting(r); o != nil || e != nil {
			return e
		}
		return workflowhandoff.WithTaskRunSnapshot(d.Workflow, r.WorkspaceRoot, func(wf workflowhandoff.Dependencies) error {
			locked := d
			locked.Workflow = wf
			fresh, e := workflowPreview(locked, r, file)
			if e != nil {
				return e
			}
			if fresh.Confirmation == nil || *fresh.Confirmation != confirm {
				return workflowError(4, "start bindings changed before reservation")
			}
			s := workflowInitial(r, fresh.Observed)
			if e = d.writeValue(workflowIndex(r.WorkspaceRoot, s.Result.RunID), s); e != nil {
				return e
			}
			if e = d.fault("workflow_after_reservation"); e != nil {
				return e
			}
			if e = workflowSave(d, s); e != nil {
				return e
			}
			draft := filepath.Join(s.Result.Paths.RunRoot, "handoff-draft.json")
			if e = d.writeOnce(draft, r.HandoffDraft); e != nil {
				return e
			}
			created, e := workflowhandoff.Create(wf, workflowhandoff.CreateInput{DraftPath: draft})
			if e != nil {
				return e
			}
			link, e := workflowhandoff.ReadTaskRunLink(wf, created.Locator)
			if e != nil {
				return e
			}
			s.Result.Handoff = WorkflowHandoff{link.HandoffID, link.Locator, link.SHA256}
			view, e := workflowhandoff.TaskRunHandoffView(wf, link.Locator)
			if e != nil {
				return e
			}
			if e = d.writeOnce(filepath.Join(s.Result.Paths.RunRoot, "mandate.json"), view); e != nil {
				return e
			}
			if e = workflowNewContext(d, &s); e != nil {
				return e
			}
			s.Phase = "tab_create_attempted"
			if e = workflowSave(d, s); e != nil {
				return e
			}
			owned = true
			return nil
		})
	})
	if e == nil && owned {
		e = workflowLaunch(d, r)
	}
	o, se := WorkflowShow(d, r.WorkspaceRoot, workflowID(r))
	if e == nil {
		e = se
	}
	if e != nil && o.RunID != "" {
		return o, workflowError(5, "Start attempt preserved; inspect the same run without restarting: "+e.Error())
	}
	return o, e
}
func workflowNewContext(d Dependencies, s *workflowState) error {
	o := &s.Result
	o.Paths = workflowPaths(s.Request, o.Round.Number)
	c := workflowContext{workflowEnv("run-context"), o.RunID, o.RequestSHA256, o.SessionID, o.Handoff.SHA256, o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256, s.Request.Runtime.PlyExecutable}
	b, e := workflowKeep(d, o.Paths.Context, c)
	if e != nil {
		return e
	}
	s.ContextSHA256 = b.SHA256
	return nil
}
func workflowInstructions(s workflowState, findings []WorkflowFinding) string {
	o := s.Result
	return fmt.Sprintf("Execute only the frozen prepared Task mandate at %s and its exact Spec inputs. The private context is %s. Read the immutable request reservation at %s. Keep its scope, human gates and agreement unchanged. Your first action before target writes is: %s workflow run accept %s --context %s --file <private-acceptance.json>. Acceptance kind ply.workflow.run-acceptance schema_version 1 uses native Acceptance@2 fields. Report actual runtime, policy and native session %s; unknown model is null, never infer it. Negative authority stops work. After your last write in each round submit kind ply.workflow.round-report schema_version 1 with all native TaskRun Report fields, the task-requirements artifact and this binding: run_id=%s request_sha256=%s session_id=%s round=%d control_id=%s previous_report_sha256=%s. Run: %s workflow run report %s --context %s --file <private-round-report.json>. Agreement A: initial execution, at most 3 corrections, 5400 active seconds, 2 environment measures; cumulative usage includes internal corrections. Budget floor: %s. Remaining budget: %s. Stop at the first reached limit. No new agents, scope expansion, TaskResult publication, installation, human QA or integration. After reporting, wait for coordinator review in this same session. Findings: %s", filepath.Join(o.Paths.RunRoot, "mandate.json"), o.Paths.Context, workflowIndex(s.Request.WorkspaceRoot, o.RunID), ShellQuote(s.Request.Runtime.PlyExecutable.Path), o.RunID, ShellQuote(o.Paths.Context), o.Transport.AgentSessionID, o.RunID, o.RequestSHA256, o.SessionID, o.Round.Number, workflowJSON(o.Round.ControlID), workflowJSON(o.Round.PreviousReportSHA256), ShellQuote(s.Request.Runtime.PlyExecutable.Path), o.RunID, ShellQuote(o.Paths.Context), workflowJSON(o.Budget.Floor), workflowJSON(map[string]int{"correction_rounds": 3 - o.Budget.Floor.CorrectionRounds, "active_seconds": 5400 - o.Budget.Floor.ActiveSeconds, "environment_measures": 2 - o.Budget.Floor.EnvironmentMeasures}), workflowJSON(findings))
}
func workflowJSON(v any) string { b, _ := Canonical(v); return string(b) }

// Bind dirty candidates by identity, HEAD, index/status and all tracked/untracked
// file contents. Symlinks are recorded as links and never traversed.
func workflowTarget(d Dependencies, s workflowState) (string, error) {
	t := s.Observed.Target
	observe := func() (any, error) {
		return d.Workspace.IntegrationGit.ObserveIntegrationWorktree(t.WorktreeLocator, t.Ref)
	}
	before, e := observe()
	if e != nil {
		return "", e
	}
	cmd := exec.Command("git", "--no-optional-locks", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-C", t.WorktreeLocator, "ls-files", "--cached", "--others", "-z")
	b, e := cmd.Output()
	if e != nil {
		return "", e
	}
	paths := strings.Split(strings.TrimSuffix(string(b), "\x00"), "\x00")
	sort.Strings(paths)
	files := map[string]string{}
	for _, rel := range paths {
		if rel == "" {
			continue
		}
		if filepath.IsAbs(rel) || filepath.Clean(rel) != rel || rel == ".." || strings.HasPrefix(rel, "../") {
			return "", workflowError(4, "unsafe candidate path")
		}
		p := filepath.Join(t.WorktreeLocator, rel)
		if e = physical(filepath.Dir(p), false); e != nil {
			return "", e
		}
		info, err := os.Lstat(p)
		if os.IsNotExist(err) {
			files[rel] = "missing"
			continue
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(p)
			if err != nil {
				return "", err
			}
			files[rel] = "symlink:" + link
			continue
		}
		if !info.Mode().IsRegular() {
			return "", workflowError(4, "candidate contains an unsupported nonregular entry")
		}
		data, err := readFile(p, 64<<20, false)
		if err != nil {
			return "", err
		}
		files[rel] = fmt.Sprint(info.Mode().Perm()) + ":" + hash(data)
	}
	after, e := observe()
	if e != nil {
		return "", e
	}
	if !equal(before, after) {
		return "", workflowError(4, "candidate changed during observation")
	}
	return digest(struct {
		Target any
		Files  map[string]string
	}{before, files}), nil
}
