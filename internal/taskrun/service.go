package taskrun

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SystemDependencies(w workspace.Dependencies) Dependencies {
	f := workflowhandoff.SystemDependencies()
	f.TaskWorkspace = &w
	w.HandoffEvidence = workflowhandoff.NewTaskHandoffEvidenceReader(f)
	return Dependencies{Executable: os.Executable, Workspace: w, Workflow: f, Runner: systemRunner(), ExecRunner: systemExecRunner(), Now: time.Now, ContextPath: func() string { return os.Getenv("PLY_TASK_RUN_CONTEXT") }, CWD: os.Getwd}
}
func containing(d Dependencies, root string) error {
	o, e := d.Workflow.Workspace.ObserveContaining()
	if e != nil {
		return e
	}
	if o.Observation.Root != root {
		return conflict("request does not name the containing workspace")
	}
	return nil
}
func existing(r Request) (*Result, error) {
	index := filepath.Join(storeRoot(r.WorkspaceRoot), "requests", strings.TrimPrefix(hash([]byte(r.RequestKey)), "sha256:")+".json")
	var reservation targetRef
	ie := readValue(index, 4096, &reservation)
	if ie == nil && reservation.RequestSHA256 != digest(r) {
		return nil, conflict("request key was already reserved with different bytes")
	}
	if ie != nil && !os.IsNotExist(ie) {
		return nil, ie
	}
	path := runPaths(r).RunRoot
	if e := physical(path, true); e != nil {
		return nil, e
	}
	_, e := os.Lstat(path)
	if os.IsNotExist(e) {
		if r.SchemaVersion == 2 {
			out, err := existingFactoryReservation(r)
			if out != nil || err != nil {
				return out, err
			}
		}
		if ie == nil {
			out := initial(r)
			out.Launch.State = "unknown"
			out.Process.State = "unknown"
			out.Collection.State = "unknown"
			return &out, nil
		}
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if os.IsNotExist(e) {
		out := initial(r)
		out.Launch.State = "unknown"
		out.Process.State = "unknown"
		out.Collection.State = "unknown"
		out.Reasons = []Reason{{"task_run_reservation_unknown", "An incomplete reservation exists; it must never be relaunched."}}
		return &out, nil
	}
	if e != nil {
		return nil, e
	}
	if !equal(j.Request, r) {
		return nil, conflict("request key already reserved with different canonical bytes")
	}
	return &j.Result, nil
}
func previewLocked(d Dependencies, r Request, file string) (p Preview, err error) {
	defer func() {
		if err != nil {
			code := "task_run_preflight_blocked"
			if typed, ok := err.(*Error); ok {
				code = typed.Code
			}
			p.Reasons = []Reason{{code, err.Error()}}
			p.Confirmation = nil
			p.NextArgv = []string{}
		}
	}()
	p = Preview{Envelope: env("preview"), RequestSHA256: digest(r), RunID: RunID(r.RequestKey), Request: r, Paths: runPaths(r), Reasons: []Reason{}, NextArgv: []string{}}
	obs, e := workflowhandoff.ObserveTaskRun(d.Workflow, r.WorkspaceRoot, r.PreparationID, r.HandoffDraft)
	p.Preparation = obs.Preparation
	p.Observed = Observed{obs.MarkerSHA256, obs.RegistrySHA256, obs.Target, obs.Epic, []FileBinding{}}
	if e != nil {
		return p, conflict(e.Error())
	}
	if digest(obs.Preparation) != r.PreparationSHA256 {
		return p, conflict("whole preparation digest differs")
	}
	p.Preparation = obs.Preparation
	if r.SchemaVersion == 2 {
		if _, _, e = validateFactory(r, p.Preparation, file, d.Now()); e != nil {
			return p, e
		}
	}
	bindings, e := runtimeBindings(r.Runtime)
	if e != nil {
		return p, e
	}
	p.Observed = Observed{obs.MarkerSHA256, obs.RegistrySHA256, obs.Target, obs.Epic, bindings}
	if e = workflowReservations(d, r.WorkspaceRoot, obs.Target); e != nil {
		return p, e
	}
	if e = physical(p.Paths.RunRoot, true); e != nil {
		return p, e
	}
	// Request reservations are the durable source; target indexes are only caches.
	reservations, err := os.ReadDir(filepath.Join(storeRoot(r.WorkspaceRoot), "requests"))
	if err != nil && !os.IsNotExist(err) {
		return p, err
	}
	for _, entry := range reservations {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		var rr targetRef
		if err = readValue(filepath.Join(storeRoot(r.WorkspaceRoot), "requests", entry.Name()), 4096, &rr); err != nil {
			return p, err
		}
		if rr.Target.WorktreeLocator == obs.Target.WorktreeLocator || rr.Target.GitCommonDir == obs.Target.GitCommonDir && rr.Target.Ref == obs.Target.Ref {
			old, err := readJournal(r.WorkspaceRoot, rr.RunID)
			if err != nil || !reservationReleased(old) {
				return p, conflict("an active or unknown reservation owns this target")
			}
		}
	}
	index := targetPath(r, obs.Target)
	if b, e := readFile(index, 4096, true); e == nil {
		var ref targetRef
		if e = decode(b, 4096, &ref); e != nil {
			return p, integrity(e.Error())
		}
		j, e := readJournal(r.WorkspaceRoot, ref.RunID)
		if e != nil {
			return p, e
		}
		if !reservationReleased(j) {
			return p, conflict("another active or unknown run reserves this worktree")
		}
	} else if !os.IsNotExist(e) {
		return p, e
	}
	p.Confirmation = ptr(confirmation(p))
	p.NextArgv = []string{r.Runtime.PlyExecutable.Path, "workspace", "task", "run", "start", "--file", file, "--apply", "--confirm", *p.Confirmation}
	encoded, e := Canonical(p)
	if e != nil {
		return p, e
	}
	if len(encoded) > 1<<20 {
		return p, invalid("preview exceeds 1 MiB")
	}
	return p, nil
}

type targetRef struct {
	RunID         string                            `json:"run_id"`
	RequestSHA256 string                            `json:"request_sha256"`
	Target        workspace.PlanWorktreeObservation `json:"target"`
}

func targetPath(r Request, x workspace.PlanWorktreeObservation) string {
	return filepath.Join(storeRoot(r.WorkspaceRoot), "targets", strings.TrimPrefix(digest([]string{x.GitCommonDir, x.Ref, x.WorktreeLocator}), "sha256:")+".json")
}
func PreviewStart(d Dependencies, file string) (any, error) {
	r, e := ReadRequest(file)
	if e != nil {
		return nil, e
	}
	if e = containing(d, r.WorkspaceRoot); e != nil {
		return nil, e
	}
	if out, e := existing(r); out != nil || e != nil {
		return out, e
	}
	return previewLocked(d, r, file)
}
func Start(d Dependencies, file, confirm string) (Result, error) {
	r, e := ReadRequest(file)
	if e != nil {
		return Result{}, e
	}
	if e = containing(d, r.WorkspaceRoot); e != nil {
		return Result{}, e
	}
	if out, e := existing(r); out != nil || e != nil {
		if out != nil {
			return *out, e
		}
		return Result{}, e
	}
	p, e := previewLocked(d, r, file)
	if e != nil {
		return Result{}, e
	}
	if p.Confirmation == nil || confirm != *p.Confirmation {
		return Result{}, conflict("start confirmation differs from preview")
	}
	runner := d.Runner
	if r.SchemaVersion == 2 {
		runner = d.ExecRunner
	}
	if runner == nil {
		return Result{}, conflict("exec runner is unavailable")
	}
	if e = runner.Check(); e != nil {
		return Result{}, e
	}
	launched := false
	var spec LaunchSpec
	e = withStore(r.WorkspaceRoot, func() error {
		if out, e := existing(r); out != nil || e != nil {
			return e
		}
		return workflowhandoff.WithTaskRunSnapshot(d.Workflow, r.WorkspaceRoot, func(wf workflowhandoff.Dependencies) error {
			locked := d
			locked.Workflow = wf
			fresh, e := previewLocked(locked, r, file)
			if e != nil {
				return e
			}
			if fresh.Confirmation == nil || confirm != *fresh.Confirmation {
				return conflict("start bindings changed before reservation")
			}
			if e = runner.Check(); e != nil {
				return e
			}
			paths := runPaths(r)
			if r.SchemaVersion == 2 {
				if e = reserveFactory(d, r, fresh.Preparation, file); e != nil {
					return e
				}
			}
			requestIndex := filepath.Join(storeRoot(r.WorkspaceRoot), "requests", strings.TrimPrefix(hash([]byte(r.RequestKey)), "sha256:")+".json")
			ref := targetRef{RunID(r.RequestKey), digest(r), fresh.Observed.Target}
			if e = d.writeValue(requestIndex, ref); e != nil {
				return e
			}
			if e = d.fault("after_request_reservation"); e != nil {
				return e
			}
			if e = privateDir(filepath.Dir(paths.RunRoot)); e != nil {
				return e
			}
			if e = os.Mkdir(paths.RunRoot, 0700); e != nil {
				return e
			}
			if e = syncDir(filepath.Dir(paths.RunRoot)); e != nil {
				return e
			}
			if e = d.writeValue(filepath.Join(paths.RunRoot, "request.json"), r); e != nil {
				return e
			}
			if e = d.fault("after_request"); e != nil {
				return e
			}
			if e = appendEvent(d, r, "reserved", map[string]any{"target": fresh.Observed.Target, "request_key": r.RequestKey}); e != nil {
				return e
			}
			tp := targetPath(r, fresh.Observed.Target)
			if e = privateDir(filepath.Dir(tp)); e != nil {
				return e
			}
			if e = d.replaceValue(tp, ref); e != nil {
				return e
			}
			if e = d.fault("after_reservation"); e != nil {
				return e
			}
			if e = d.writeOnce(filepath.Join(paths.RunRoot, "tmp", "handoff-draft.json"), r.HandoffDraft); e != nil {
				return e
			}
			created, e := workflowhandoff.Create(wf, workflowhandoff.CreateInput{DraftPath: filepath.Join(paths.TempRoot, "handoff-draft.json")})
			if e != nil {
				return e
			}
			link, e := workflowhandoff.ReadTaskRunLink(wf, created.Locator)
			if e != nil {
				return e
			}
			b := Binding{env("binding"), RunID(r.RequestKey), digest(r), fresh.Preparation, Handoff{link.ActivityID, link.RunID, link.HandoffID, link.Locator, link.SHA256}, "ply:" + RunID(r.RequestKey), paths.RunRoot, paths.TempRoot, paths.ReportPath}
			if e = d.writeValue(filepath.Join(paths.RunRoot, "binding.json"), b); e != nil {
				return e
			}
			if e = appendEvent(d, r, "handoff_bound", map[string]any{"binding_sha256": digest(b)}); e != nil {
				return e
			}
			ctx := Context{env("context"), b.RunID, b.RequestSHA256, b.SessionID, r.WorkspaceRoot, digest(b), r.Runtime.PlyExecutable}
			cp := filepath.Join(paths.TempRoot, "context.json")
			if e = d.writeValue(cp, ctx); e != nil {
				return e
			}
			view, e := workflowhandoff.TaskRunHandoffView(wf, link.Locator)
			if e != nil {
				return e
			}
			if e = d.writeOnce(filepath.Join(paths.TempRoot, "mandate.json"), view); e != nil {
				return e
			}
			instructions := recipientInstructions(r, b)
			ip := filepath.Join(paths.TempRoot, "recipient.md")
			if e = d.writeOnce(ip, []byte(instructions)); e != nil {
				return e
			}
			argv := []string{r.Runtime.Executable.Path, "--cd", fresh.Observed.Target.WorktreeLocator, "--model", r.Runtime.Model}
			if r.Runtime.ConfigProfile != nil {
				argv = append(argv, "--profile", *r.Runtime.ConfigProfile)
			}
			argv = append(argv, "--no-alt-screen", "Read and execute the private Task run instructions at "+ip+". Your first action is the typed accept callback before any target write.")
			spec = LaunchSpec{Executable: r.Runtime.Executable, Argv: argv, CWD: fresh.Observed.Target.WorktreeLocator, ContextPath: cp}
			if r.SchemaVersion == 2 {
				a, policy, err := validateFactory(r, fresh.Preparation, file, d.Now())
				if err != nil {
					return err
				}
				argv, err = execArgv(r, policy, spec.CWD, ip)
				if err != nil {
					return err
				}
				spec.Argv = argv
				spec.Timeout = time.Duration(r.FactoryTest.TimeoutSeconds) * time.Second
				deadline, _ := time.Parse(time.RFC3339Nano, a.DeadlineUTC)
				if remaining := deadline.Sub(d.Now()); remaining < spec.Timeout {
					spec.Timeout = remaining
				}
				spec.StreamsRoot = factorySlotRoot(r, a)
				spec.Completion = &ProviderCompletion{}
				if e = d.writeValue(filepath.Join(spec.StreamsRoot, "launch.json"), ExecLaunch{r.Runtime.Executable, argv, spec.CWD, digest(r), spec.Timeout.Milliseconds()}); e != nil {
					return e
				}
			}
			if _, e = runtimeBindings(r.Runtime); e != nil {
				return e
			}
			if e = appendEvent(d, r, "launch_intent", map[string]any{"executable_sha256": r.Runtime.Executable.SHA256, "argv_sha256": digest(argv), "cwd": spec.CWD, "session_id": b.SessionID, "effective_policy_sha256": r.Runtime.PermissionBinding.EffectivePolicySHA256}); e != nil {
				return e
			}
			if e = d.fault("after_launch_intent"); e != nil {
				return e
			}
			launched = true
			return nil
		})
	})
	if e != nil {
		out, _ := Show(d, r.WorkspaceRoot, RunID(r.RequestKey))
		return out, e
	}
	if !launched {
		return Show(d, r.WorkspaceRoot, RunID(r.RequestKey))
	}
	process, runErr := runner.Run(spec, func(p Process) error {
		return withStore(r.WorkspaceRoot, func() error { return appendEvent(d, r, "process_started", map[string]any{"process": p}) })
	})
	e = withStore(r.WorkspaceRoot, func() error {
		if runErr != nil && process.State == "not_started" {
			return appendEvent(d, r, "launch_failed", map[string]any{"code": "task_run_launch_failed", "detail": "The bound provider could not be started."})
		}
		if e := appendEvent(d, r, "process_exited", map[string]any{"process": process}); e != nil {
			return e
		}
		if r.SchemaVersion == 2 && spec.Completion != nil && spec.Completion.Kind != "" {
			return publishProviderCompletion(d, r, spec, process)
		}
		return nil
	})
	if e != nil {
		return Result{}, failure("task_run_process_unknown", 5, "Unable to preserve process outcome: "+e.Error())
	}
	out, ce := collectApply(d, r.WorkspaceRoot, RunID(r.RequestKey), "", true)
	if ce != nil {
		return out, ce
	}
	if runErr != nil && process.State == "not_started" {
		return out, failure("task_run_launch_failed", 4, "Provider start failed; attempt preserved.")
	}
	if out.Collection.State != "qualified" && out.TaskExecution.State == "unknown" {
		if r.SchemaVersion == 2 {
			return out, failure("task_run_provider_completion_unknown", 5, "Exec completion is unknown; preserve the run and work environment without restart or teardown.")
		}
		return out, failure("task_run_task_status_unknown", 5, "Client outcome and return are preserved. Obtain an explicit human task observation, then preview collect --task-status and apply its confirmation; plan result control remains the next gate.")
	}
	if out.Collection.State != "qualified" {
		return out, failure("task_run_outcome_unknown", 5, "Delivery is preserved and requires plan result control.")
	}
	return out, nil
}
func recipientInstructions(r Request, b Binding) string {
	text := fmt.Sprintf("# Task run delivery\n\nRead %s for the exact preserved Task Spec inputs, procedure, verifiers, authority and human gates. Read %s for the immutable request and agreement. Do not select a newer Spec.\n\nYour first action, before target writes, is to report your actual runtime facts and scope contract in ply.workspace.task-run-acceptance schema version 2. Use run_id %s, request_sha256 %s, session_id %s. Model may be null when unknown; never infer it from the requested model. Known runtime, profile and effective policy are required for started. Negative claims may use null for unknown runtime, model, profile, policy and sandbox and must explain issues. Native provider session may be null; never invent one. runtime_claim has runtime_id, model_id, profile_id, effective_policy_sha256, native_session_id. sandbox and issues use the existing WF types; sandbox is the reported contract, not an OS policy claim. Acceptance must be negative when authority is unknown. Do not read or print raw WF reply secrets.\n\nRun: %s workspace task run accept %s --context %s --file <absolute-private-claim.json>\n\nAgreement A: initial execution plus at most three correction rounds, 5400 active seconds, and two environment measures. Stop at the first limit. No other agent, nested provider, human QA, integration, install, deployment or remote effects. The provider contact is owned by the parent transport only.\n\nAfter your last correction, write the semantic ply.workspace.task-run-report schema version 1 to %s and run: %s workspace task run report %s --context %s --file %s\nRequired fields: run_id, request_sha256, session_id, outcome, summary, meaning, budget_usage, stop_reasons, observed_effects, verifier_results, review, artifacts, evidence_gaps, forbidden_effects_observed, technical_assessment, process_observation, preventive_followup. WF terminal fields retain existing types. budget_usage is initial_execution_started, correction_rounds, active_seconds, environment_measures, or null when unknown. technical_assessment is gate, required_verifier_ids, accepted_debt. Last two fields are null for none. Never report an unexecuted verifier as exit 0. No TaskResult or human QA claim is yours to manufacture.\n\nThe explicit context works in a clean tool shell without inherited environment. The parent collects after the client process exits. Client exit alone does not prove the managed task inactive; a separate explicit human observation is required for qualification. Pending task status is expected. After reporting, make no more target writes. Tell the human the report was received and they can exit Codex normally. Plan result control is the next gate.\n", filepath.Join(b.TempRoot, "mandate.json"), filepath.Join(b.RunRoot, "request.json"), b.RunID, b.RequestSHA256, b.SessionID, ShellQuote(r.Runtime.PlyExecutable.Path), b.RunID, ShellQuote(filepath.Join(b.TempRoot, "context.json")), b.ReportPath, ShellQuote(r.Runtime.PlyExecutable.Path), b.RunID, ShellQuote(filepath.Join(b.TempRoot, "context.json")), ShellQuote(b.ReportPath))
	if r.SchemaVersion == 2 {
		text = "Follow the frozen mandate's delivery procedure first. When it supplies a return helper, use that helper for acceptance and reporting; it constructs the bound callback fields. Check that the report has delivery.state=received and terminal_sha256 before ending the turn. Correct rejected input within the same call and budget.\n\n" + text
		text = strings.Replace(text, "Client exit alone does not prove the managed task inactive; a separate explicit human observation is required for qualification. Pending task status is expected. After reporting, make no more target writes. Tell the human the report was received and they can exit Codex normally.", "The parent validates native provider completion and process group quiescence. After reporting, make no more target writes and finish this single turn. Never submit manual task status.", 1)
		text += fmt.Sprintf("\nFactory limit: %d seconds for this one call, reasoning low; no retries or other agent starts. Report the required task-requirements artifact with all requirement IDs, phases and verifier IDs from the frozen Spec. Preserve unknown model as null.\n", r.FactoryTest.TimeoutSeconds)
	}
	return text
}
func ShellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'" }
func quiescent(p Process) bool {
	return p.Quiescence != nil && *p.Quiescence && (p.State == "exited" || p.State == "not_started")
}
func Show(d Dependencies, root, id string) (Result, error) {
	if e := containing(d, root); e != nil {
		return Result{}, e
	}
	j, e := readJournal(root, id)
	if e != nil {
		return Result{}, e
	}
	if e = validateRunEvidence(d, j); e != nil {
		return Result{}, e
	}
	out := j.Result
	if record, err := recordedResult(d, j); err != nil {
		return Result{}, err
	} else if record != nil {
		out.Collection = historicalCollection(record)
		if out.TaskExecution.StatusSHA256 == nil && j.Request.SchemaVersion == 1 {
			out.TaskExecution.HistoricalQualification = true
		}
		if out.ObservedTarget == nil {
			out.ObservedTarget = historicalTarget(record)
		}
	}
	if j.Binding != nil {
		p := j.Binding.Preparation.Plan
		x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(p.WorktreePath, "refs/heads/"+p.Branch)
		if e == nil {
			actual := workspace.PlanWorktreeObservation{WorktreeLocator: x.Locator, Ref: x.Ref, OID: x.OID, Tree: x.Tree, GitCommonDir: x.GitCommonDir, ObjectFormat: x.ObjectFormat, RefFormat: x.RefFormat, Symbolic: x.Symbolic, Clean: x.Clean, StatusEntries: x.StatusEntries, InProgress: x.InProgress}
			if out.ObservedTarget != nil && (out.ObservedTarget.OID != actual.OID || out.ObservedTarget.Tree != actual.Tree || out.ObservedTarget.Ref != actual.Ref || out.ObservedTarget.GitCommonDir != actual.GitCommonDir || out.ObservedTarget.WorktreeLocator != actual.WorktreeLocator || !actual.Clean) {
				out.Reasons = append(out.Reasons, Reason{"task_run_current_target_changed", "The current target differs from the historical collection; historical receipts and result are retained."})
			}
			out.ObservedTarget = &actual
		} else {
			out.Reasons = append(out.Reasons, Reason{"task_run_current_target_unknown", "The current target could not be observed."})
		}
	}
	if len(j.Events) > 0 {
		start, _ := time.Parse(time.RFC3339Nano, j.Events[0].RecordedAtUTC)
		end := d.Now()
		for _, ev := range j.Events {
			if ev.Type == "process_exited" {
				end, _ = time.Parse(time.RFC3339Nano, ev.RecordedAtUTC)
			}
		}
		secs := int64(end.Sub(start).Seconds())
		if secs < 0 {
			secs = 0
		}
		out.Budget.ElapsedSeconds = &secs
	}
	if out.Acceptance.State != "missing" && out.Acceptance.ReceiptSHA256 == nil {
		out.Reasons = append(out.Reasons, Reason{"task_run_acceptance_not_representable", "The recipient claim was received; missing facts prevent a WF StartReceipt. This is a representation gap, not permission to work."})
	}
	if out.Delivery.State == "not_representable" {
		detail := "The preserved report has no workflow terminal; missing budget, start or terminal facts prevent representation."
		var diagnostic struct {
			Code   string `json:"code"`
			Detail string `json:"detail"`
		}
		if readValue(filepath.Join(runPaths(j.Request).RunRoot, "reports", "terminal-diagnostic.json"), 64<<10, &diagnostic) == nil && diagnostic.Detail != "" {
			detail = diagnostic.Detail
		}
		out.Reasons = append(out.Reasons, Reason{"task_run_terminal_not_representable", detail})
	}
	if j.Request.SchemaVersion == 1 && (out.TaskExecution.State == "unknown" || out.TaskExecution.State == "active") {
		out.NextAction = "Obtain an explicit human task observation; use collect --task-status with check, then apply its confirmation. Plan result control precedes separate human QA or integration."
	}
	out.Reasons = sortedReasons(out.Reasons)
	return out, nil
}
