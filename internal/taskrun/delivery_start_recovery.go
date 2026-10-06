package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// DeliveryStartRecoveryInput replaces only executable identities, before any
// Task prompt. ControlExecutable is the source binary, not its preserved copy.
type DeliveryStartRecoveryInput struct {
	RunID              string
	ProviderExecutable Executable
	ControlExecutable  Executable
	Timeout            time.Duration
}

type DeliveryStartRecoveryReadback struct {
	Generation                 int         `json:"generation"`
	Binding                    FileBinding `json:"binding"`
	OriginalProviderExecutable Executable  `json:"original_provider_executable"`
	OriginalControlExecutable  Executable  `json:"original_control_executable"`
	ProviderExecutable         Executable  `json:"provider_executable"`
	ControlExecutable          Executable  `json:"control_executable"`
}

type DeliveryStartRecoveryPreview struct {
	Kind               string                       `json:"kind"`
	SchemaVersion      int                          `json:"schema_version"`
	RunID              string                       `json:"run_id"`
	State              string                       `json:"state"`
	Confirmation       string                       `json:"confirmation"`
	ProviderExecutable Executable                   `json:"provider_executable"`
	ControlExecutable  Executable                   `json:"control_executable"`
	RecoveryDirectory  string                       `json:"recovery_directory"`
	Run                *WorkflowRun                 `json:"run"`
	ExitEvidence       *DeliveryStartupExitEvidence `json:"exit_evidence"`
}

// Only typed process and pane facts are kept. Terminal text, argv and process
// command lines are neither a release proof nor safe diagnostic output.
type DeliveryStartupExitEvidence struct {
	WorkspaceID              string           `json:"workspace_id"`
	TabID                    string           `json:"tab_id"`
	PaneID                   string           `json:"pane_id"`
	TerminalID               string           `json:"terminal_id"`
	CWD                      string           `json:"cwd"`
	ShellPID                 int              `json:"shell_pid"`
	ForegroundProcessGroupID int              `json:"foreground_process_group_id"`
	ShellName                string           `json:"shell_name"`
	ProcessTree              []StartupProcess `json:"process_tree"`
}

type deliveryStartRecoveryRecord struct {
	Kind               string                      `json:"kind"`
	SchemaVersion      int                         `json:"schema_version"`
	RunID              string                      `json:"run_id"`
	RequestSHA256      string                      `json:"request_sha256"`
	BeforeState        FileBinding                 `json:"before_state"`
	SourceControl      Executable                  `json:"source_control"`
	ProviderExecutable Executable                  `json:"provider_executable"`
	ControlExecutable  Executable                  `json:"control_executable"`
	Context            FileBinding                 `json:"context"`
	ExitEvidence       DeliveryStartupExitEvidence `json:"exit_evidence"`
	Confirmation       string                      `json:"confirmation"`
	RecordedAtUTC      string                      `json:"recorded_at_utc"`
}

func deliveryRecoveryDirectory(s workflowState) string {
	return filepath.Join(s.Result.Paths.RunRoot, "startup-recovery", "001")
}

func workflowStartupPath(s workflowState, name string) string {
	if s.Recovery != nil {
		return filepath.Join(deliveryRecoveryDirectory(s), name)
	}
	return filepath.Join(s.Result.Paths.RunRoot, name)
}

func workflowStartupGeneration(s workflowState) string {
	if s.Recovery != nil {
		return s.Recovery.SHA256
	}
	return "original"
}

func workflowDeliveryGuideDirectory(s workflowState) string {
	if s.Recovery != nil {
		return filepath.Join(deliveryRecoveryDirectory(s), "delivery")
	}
	return filepath.Join(s.Result.Paths.RunRoot, "delivery")
}

func workflowRecoveryRecord(s workflowState) (*deliveryStartRecoveryRecord, error) {
	if s.Recovery == nil {
		if s.Result.StartupRecovery != nil {
			return nil, workflowError(4, "startup recovery readback has no bound record")
		}
		return nil, nil
	}
	dir := deliveryRecoveryDirectory(s)
	if s.Recovery.Locator != filepath.Join(dir, "recovery.json") {
		return nil, workflowError(4, "startup recovery path differs")
	}
	b, err := workflowBound(*s.Recovery, 1<<20)
	if err != nil {
		return nil, err
	}
	var r deliveryStartRecoveryRecord
	if err = decode(b, 1<<20, &r); err != nil {
		return nil, err
	}
	if r.Kind != "PlyDeliveryStartupRecovery@1" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.BeforeState.Locator != filepath.Join(dir, "before-state.json") || r.ControlExecutable.Path != filepath.Join(dir, "ply-control") || r.Context.Locator != filepath.Join(dir, "context.json") || s.Result.Paths.Context != r.Context.Locator || s.ContextSHA256 != r.Context.SHA256 {
		return nil, workflowError(4, "startup recovery differs from its original run")
	}
	before, err := workflowBound(r.BeforeState, 8<<20)
	if err != nil {
		return nil, err
	}
	if _, err = workflowBound(r.Context, 1<<20); err != nil {
		return nil, err
	}
	var old workflowState
	if err = decode(before, 8<<20, &old); err != nil {
		return nil, err
	}
	if old.Recovery != nil || !equal(old.Request, s.Request) || !equal(old.Observed, s.Observed) || old.Result.RunID != s.Result.RunID || old.Result.Handoff != s.Result.Handoff || old.Result.Transport.AgentSessionID != "" || old.Acceptance != nil || old.StartSHA256 != nil || old.StartDraft != nil || old.Result.Delivery == nil || len(old.Result.Delivery.Events) != 0 || len(old.Result.Delivery.Candidates) != 0 {
		return nil, workflowError(4, "startup recovery before-state is not the original unstarted delivery")
	}
	want := DeliveryStartRecoveryReadback{1, *s.Recovery, s.Request.Runtime.Executable, s.Request.Runtime.PlyExecutable, r.ProviderExecutable, r.ControlExecutable}
	if s.Result.StartupRecovery == nil || !equal(*s.Result.StartupRecovery, want) {
		return nil, workflowError(4, "startup recovery readback binding differs")
	}
	return &r, nil
}

// The immutable request remains the original mandate. Only this resolver may
// select executable replacements; callers must never persist an effective copy.
func workflowEffectiveRuntime(s workflowState) (Runtime, error) {
	r := s.Request.Runtime
	recovery, err := workflowRecoveryRecord(s)
	if err != nil {
		return r, err
	}
	if recovery != nil {
		r.Executable, r.PlyExecutable = recovery.ProviderExecutable, recovery.ControlExecutable
	}
	return r, nil
}

func deliveryRecoveryExisting(d Dependencies, root string, s workflowState) (DeliveryStartRecoveryPreview, error) {
	p := DeliveryStartRecoveryPreview{Kind: "PlyDeliveryStartupRecoveryPreview@1", SchemaVersion: 1, RunID: s.Result.RunID, State: "existing", RecoveryDirectory: deliveryRecoveryDirectory(s)}
	r, err := workflowRecoveryRecord(s)
	if err != nil {
		return p, err
	}
	if r == nil {
		return p, workflowError(4, "startup recovery is not bound")
	}
	p.Confirmation, p.ProviderExecutable, p.ControlExecutable = r.Confirmation, r.ProviderExecutable, r.ControlExecutable
	o, err := WorkflowShow(d, root, s.Result.RunID)
	p.Run = &o
	return p, err
}

func deliveryRecoveryPreTask(d Dependencies, s workflowState) error {
	if !deliveryRun(s.Request) || s.Request.Runtime.Provider != "codex" || s.Recovery != nil || s.Result.Delivery == nil || s.Acceptance != nil || s.StartDraft != nil || s.StartSHA256 != nil || s.Result.Transport.AgentSessionID != "" || s.Result.Round.Number != 0 || len(s.Records) != 0 || len(s.Result.Delivery.Events) != 0 || len(s.Result.Delivery.Candidates) != 0 || s.Result.Delivery.Attempt != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" || s.Result.Delivery.Phase != "awaiting_acceptance" {
		return workflowError(4, "startup recovery requires the original Codex delivery before any Task prompt, session or acceptance")
	}
	if s.Phase != "agent_start_attempted" && s.Phase != "bootstrap_attempted" {
		return workflowError(4, "startup recovery is not available at this phase")
	}
	for _, path := range []string{filepathForRound(s, "prompt-attempt.json"), filepath.Join(s.Result.Paths.RunRoot, "acceptance.json"), filepath.Join(s.Result.Paths.RunRoot, "start-draft.json"), filepath.Join(deliveryRecoveryDirectory(s), "recovery.json")} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			return workflowError(4, "startup recovery cannot replace existing or unreadable Task/start artifacts: "+path)
		}
	}
	// A crash may publish an immutable artifact before updating mutable state.
	// Inspect every preserved round, not only the state's claimed round zero.
	rounds := filepath.Join(s.Result.Paths.RunRoot, "rounds")
	if err := physical(rounds, false); err != nil {
		return err
	}
	entries, err := os.ReadDir(rounds)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(rounds, entry.Name())
		if err := physical(path, false); err != nil {
			return err
		}
		if !entry.IsDir() {
			return workflowError(4, "unexpected startup round artifact")
		}
		if _, err := os.Lstat(filepath.Join(path, "prompt-attempt.json")); !os.IsNotExist(err) {
			return workflowError(4, "a preserved round may already have attempted Task input")
		}
	}
	for _, name := range []string{"events", "attempts", "candidates"} {
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", name)
		if err := physical(path, true); err != nil {
			return err
		}
		entries, err := os.ReadDir(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if len(entries) > 0 {
			return workflowError(4, "preserved delivery effects prevent startup recovery")
		}
	}
	facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
	if err != nil {
		return err
	}
	if facts.Start != nil || facts.Terminal != nil {
		return workflowError(4, "native Task start or terminal already exists")
	}
	return nil
}

func deliveryRecoveryExit(d Dependencies, s workflowState) (DeliveryStartupExitEvidence, error) {
	var out DeliveryStartupExitEvidence
	var pane struct {
		Pane struct {
			WorkspaceID   string          `json:"workspace_id"`
			TabID         string          `json:"tab_id"`
			PaneID        string          `json:"pane_id"`
			TerminalID    string          `json:"terminal_id"`
			ForegroundCWD string          `json:"foreground_cwd"`
			Agent         *string         `json:"agent"`
			Session       json.RawMessage `json:"agent_session"`
			Status        string          `json:"agent_status"`
		} `json:"pane"`
	}
	b, err := workflowCall(d, s.Request, "pane", "get", s.Result.Transport.PaneID)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(b, &pane); err != nil {
		return out, workflowError(4, "invalid recovery pane observation")
	}
	p, t := pane.Pane, s.Result.Transport
	if p.WorkspaceID != t.WorkspaceID || p.TabID != t.TabID || p.PaneID != t.PaneID || p.TerminalID != t.TerminalID || p.Agent != nil || len(p.Session) > 0 && string(p.Session) != "null" || p.Status != "unknown" || p.ForegroundCWD != s.Observed.Target.WorktreeLocator {
		return out, workflowError(4, "recovery requires the exact preserved terminal without a provider or native session")
	}
	var process struct {
		Info struct {
			PaneID    string `json:"pane_id"`
			ShellPID  *int   `json:"shell_pid"`
			Group     *int   `json:"foreground_process_group_id"`
			Processes []struct {
				PID  int     `json:"pid"`
				Name string  `json:"name"`
				CWD  *string `json:"cwd"`
			} `json:"foreground_processes"`
		} `json:"process_info"`
	}
	b, err = workflowCall(d, s.Request, "pane", "process-info", "--pane", t.PaneID)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(b, &process); err != nil {
		return out, workflowError(4, "invalid recovery process observation")
	}
	i := process.Info
	if i.PaneID != t.PaneID || i.ShellPID == nil || *i.ShellPID <= 0 || i.Group == nil || *i.Group != *i.ShellPID || len(i.Processes) != 1 || i.Processes[0].PID != *i.ShellPID || i.Processes[0].CWD == nil || *i.Processes[0].CWD != s.Observed.Target.WorktreeLocator {
		return out, workflowError(4, "provider exit is not proven by a sole foreground shell in the exact Task worktree")
	}
	name := i.Processes[0].Name
	if name != "zsh" && name != "bash" && name != "fish" && name != "sh" {
		return out, workflowError(4, "foreground process is not a supported shell")
	}
	tree, err := deliveryRecoveryProcessTree(d, *i.ShellPID, *i.Group)
	if err != nil {
		return out, err
	}
	return DeliveryStartupExitEvidence{p.WorkspaceID, p.TabID, p.PaneID, p.TerminalID, p.ForegroundCWD, *i.ShellPID, *i.Group, name, tree}, nil
}

func deliveryRecoveryPreview(d Dependencies, root string, in DeliveryStartRecoveryInput, s workflowState) (DeliveryStartRecoveryPreview, error) {
	p := DeliveryStartRecoveryPreview{Kind: "PlyDeliveryStartupRecoveryPreview@1", SchemaVersion: 1, RunID: in.RunID, State: "blocked", RecoveryDirectory: deliveryRecoveryDirectory(s), ProviderExecutable: in.ProviderExecutable, ControlExecutable: Executable{filepath.Join(deliveryRecoveryDirectory(s), "ply-control"), in.ControlExecutable.SHA256}}
	if err := deliveryRecoveryPreTask(d, s); err != nil {
		return p, err
	}
	if in.Timeout <= 0 {
		return p, workflowError(2, "timeout must be positive")
	}
	if err := verifyExecutable(in.ProviderExecutable); err != nil {
		return p, err
	}
	resolved, err := exec.LookPath(s.Request.Runtime.Provider)
	if err == nil {
		resolved, err = filepath.EvalSymlinks(resolved)
	}
	if err != nil || resolved != in.ProviderExecutable.Path {
		return p, workflowError(4, "local provider executable differs from the proposed recovery binding")
	}
	if err := verifyExecutable(in.ControlExecutable); err != nil {
		return p, err
	}
	if in.ProviderExecutable == s.Request.Runtime.Executable {
		return p, workflowError(4, "recovery requires an explicit replacement provider executable")
	}
	if err := physical(p.RecoveryDirectory, true); err != nil {
		return p, err
	}
	if b, err := os.Lstat(p.ControlExecutable.Path); err == nil {
		if !b.Mode().IsRegular() {
			return p, workflowError(4, "preserved recovery control is not a regular file")
		}
		if err = verifyExecutable(p.ControlExecutable); err != nil {
			return p, err
		}
	} else if !os.IsNotExist(err) {
		return p, err
	}
	// The old provider may no longer exist. Every other frozen input retains its
	// original checks, using only the explicitly proposed executable replacements.
	if err := workflowFreshRuntime(d, s, false, func() Runtime {
		r := s.Request.Runtime
		r.Executable = in.ProviderExecutable
		r.PlyExecutable = in.ControlExecutable
		return r
	}()); err != nil {
		return p, err
	}
	obs, err := workflowhandoff.ObserveDeliveryTaskRun(d.Workflow, root, s.Request.PreparationID, s.Request.HandoffDraft)
	if err != nil {
		return p, err
	}
	if digest(obs.Preparation) != s.Request.PreparationSHA256 || !equal(obs.Target, s.Observed.Target) || !equal(obs.Epic, s.Observed.Epic) || obs.MarkerSHA256 != s.Observed.WorkspaceMarkerSHA256 {
		return p, workflowError(4, "prepared Task, return worktree or workspace changed before recovery")
	}
	cwd, err := d.CWD()
	if err != nil {
		return p, err
	}
	if cwd != s.Observed.Target.WorktreeLocator && cwd != s.Observed.Epic.WorktreeLocator {
		return p, workflowError(4, "recovery must run from the exact Task or return worktree")
	}
	x, err := deliveryTarget(d, s)
	if err != nil {
		return p, err
	}
	if x.OID != s.Observed.Target.OID || x.Tree != s.Observed.Target.Tree {
		return p, workflowError(4, "Task target changed before recovery")
	}
	exit, err := deliveryRecoveryExit(d, s)
	if err != nil {
		return p, err
	}
	p.ExitEvidence = &exit
	state, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
	if err != nil {
		return p, err
	}
	var captured workflowState
	if err = decode(state, 8<<20, &captured); err != nil {
		return p, err
	}
	if !equal(captured, s) {
		return p, workflowError(4, "startup state changed during recovery preview")
	}
	p.Confirmation = digest(struct {
		State                            string
		Provider, SourceControl, Control Executable
		Exit                             DeliveryStartupExitEvidence
	}{hash(state), in.ProviderExecutable, in.ControlExecutable, p.ControlExecutable, exit})
	p.State = "ready"
	return p, nil
}

// Preview performs only observations. It neither copies a control binary nor
// locks, reserves, starts or prompts an agent.
func WorkflowPreviewDeliveryStartRecovery(d Dependencies, root string, in DeliveryStartRecoveryInput) (DeliveryStartRecoveryPreview, error) {
	if err := containing(d, root); err != nil {
		return DeliveryStartRecoveryPreview{}, err
	}
	s, err := workflowRead(root, in.RunID)
	if err != nil {
		return DeliveryStartRecoveryPreview{}, err
	}
	if s.Recovery != nil {
		return deliveryRecoveryExisting(d, root, s)
	}
	return deliveryRecoveryPreview(d, root, in, s)
}

func WorkflowRecoverDeliveryStart(d Dependencies, root string, in DeliveryStartRecoveryInput, confirmation string) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	owned := false
	err := withStore(root, func() error {
		s, err := workflowRead(root, in.RunID)
		if err != nil {
			return err
		}
		if s.Recovery != nil {
			return nil
		}
		if os.Getenv("HERDR_ENV") != "1" {
			return workflowError(4, "recovery start requires the local Herdr context")
		}
		p, err := deliveryRecoveryPreview(d, root, in, s)
		if err != nil {
			return err
		}
		if confirmation == "" || p.Confirmation != confirmation {
			return workflowError(4, "startup recovery confirmation differs from current observations")
		}
		if err = verifyExecutable(p.ControlExecutable); err != nil {
			return err
		}
		before, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
		if err != nil {
			return err
		}
		beforeBinding := FileBinding{filepath.Join(p.RecoveryDirectory, "before-state.json"), hash(before)}
		if err = d.writeOnce(beforeBinding.Locator, before); err != nil {
			return err
		}
		contextPath := filepath.Join(p.RecoveryDirectory, "context.json")
		contextValue := workflowContext{workflowEnv("run-context"), s.Result.RunID, s.Result.RequestSHA256, s.Result.SessionID, s.Result.Handoff.SHA256, s.Result.Round.Number, s.Result.Round.ControlID, s.Result.Round.PreviousReportSHA256, p.ControlExecutable}
		contextBinding, err := workflowKeep(d, contextPath, contextValue)
		if err != nil {
			return err
		}
		r := deliveryStartRecoveryRecord{"PlyDeliveryStartupRecovery@1", 1, in.RunID, s.Result.RequestSHA256, beforeBinding, in.ControlExecutable, in.ProviderExecutable, p.ControlExecutable, contextBinding, *p.ExitEvidence, p.Confirmation, d.Now().UTC().Format(time.RFC3339Nano)}
		binding, err := workflowKeep(d, filepath.Join(p.RecoveryDirectory, "recovery.json"), r)
		if err != nil {
			return err
		}
		s.Recovery = &binding
		s.Result.StartupRecovery = &DeliveryStartRecoveryReadback{1, binding, s.Request.Runtime.Executable, s.Request.Runtime.PlyExecutable, in.ProviderExecutable, p.ControlExecutable}
		s.Result.Paths.Context, s.ContextSHA256 = contextPath, contextBinding.SHA256
		// Reserve the one new start before its external effect. A lost reply or
		// crash never permits another start on repetition of this operation.
		if err = d.writeValue(workflowStartupPath(s, "agent-start-attempt.json"), map[string]any{"recovery_sha256": binding.SHA256, "argv": workflowStartArgv(s.Request, s.Observed.Target.WorktreeLocator, s.Result.Transport.PaneID)}); err != nil {
			return err
		}
		s.Phase = "agent_start_attempted"
		s.Result.Transport.State = "unknown"
		s.Result.Transport.Observation = "cached"
		s.Result.NextAction = WorkflowAction{"coordinator", "Follow this one recovery startup; do not restart or resend uncertain input."}
		if err = workflowSave(d, s); err != nil {
			return err
		}
		owned = true
		return nil
	})
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if !owned {
		return WorkflowShow(d, root, in.RunID)
	}
	if err = d.fault("delivery_recovery_before_agent_send"); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	s, err := workflowRead(root, in.RunID)
	if err != nil {
		return WorkflowRun{}, err
	}
	deadline := time.Now().Add(in.Timeout)
	_, startErr := workflowStartupCall(d, s, deadline, workflowStartArgv(s.Request, s.Observed.Target.WorktreeLocator, s.Result.Transport.PaneID)...)
	if startErr != nil {
		if err = workflowUpdate(d, root, in.RunID, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"delivery_recovery_start_response", startErr.Error()})
			return nil
		}); err != nil {
			return deliveryReadback(d, root, in.RunID, err)
		}
	}
	err = workflowAwaitReadinessGeneration(d, root, in.RunID, deadline, workflowStartupGeneration(s))
	if err == nil {
		err = workflowPromptUntil(d, root, in.RunID, nil, deadline)
	}
	return deliveryReadback(d, root, in.RunID, err)
}
