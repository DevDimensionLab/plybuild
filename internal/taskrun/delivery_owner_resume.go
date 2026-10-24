package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DeliveryOwnerResumeInput resumes the same provider conversation of a begun
// interactive owner after its Herdr transport is positively absent. Only the
// transport and the callback control are replaced; the Run, Task worktree,
// frozen mandate, context, native session and all delivery evidence remain.
type DeliveryOwnerResumeInput struct {
	RunID                  string
	ProviderExecutable     Executable
	ControlExecutable      Executable
	DestinationWorkspaceID string
	Timeout                time.Duration
	InvokerPaneID          string
}

type DeliveryOwnerResumePreview struct {
	Kind                   string                         `json:"kind"`
	SchemaVersion          int                            `json:"schema_version"`
	RunID                  string                         `json:"run_id"`
	Provider               string                         `json:"provider"`
	State                  string                         `json:"state"`
	Phase                  string                         `json:"phase"`
	Reason                 string                         `json:"reason,omitempty"`
	Detail                 string                         `json:"detail,omitempty"`
	Generation             int                            `json:"generation,omitempty"`
	Directory              string                         `json:"directory,omitempty"`
	NativeSessionID        string                         `json:"native_session_id"`
	ProviderSession        string                         `json:"provider_session,omitempty"`
	ProviderExecutable     Executable                     `json:"provider_executable"`
	ControlExecutable      Executable                     `json:"control_executable"`
	DestinationWorkspaceID string                         `json:"destination_workspace_id,omitempty"`
	InvokerPaneID          string                         `json:"invoker_pane_id,omitempty"`
	PreviousTransport      WorkflowTransport              `json:"previous_transport"`
	Absence                *HumanOwnerRecoveryObservation `json:"absence,omitempty"`
	Liveness               *TaskWriterLiveness            `json:"writer_liveness,omitempty"`
	ResumeArgv             []string                       `json:"resume_argv,omitempty"`
	Confirmation           string                         `json:"confirmation"`
	NextAction             WorkflowAction                 `json:"next_action"`
}

// DeliveryOwnerResumeReadback separates the reserved start, the observed
// resumed session, the sent continuation and the first functional callback.
// TaskWriterLiveness is the reconciled local basis for starting a new writer:
// no process outside the invoking command chain has its current directory in
// the Task worktree. Excluded identities are diagnostic, not authority.
type TaskWriterLiveness struct {
	Method        string `json:"method"`
	Scope         string `json:"scope"`
	RemainingPIDs []int  `json:"remaining_pids"`
	ExcludedPIDs  []int  `json:"excluded_pids"`
	Limitation    string `json:"limitation"`
}

type DeliveryOwnerResumeReadback struct {
	Generation              int               `json:"generation"`
	Step                    string            `json:"step"`
	NativeSessionID         string            `json:"native_session_id"`
	PreviousTransport       WorkflowTransport `json:"previous_transport"`
	ProviderExecutable      Executable        `json:"provider_executable"`
	ControlExecutable       Executable        `json:"control_executable"`
	SessionObserved         bool              `json:"session_observed"`
	LivenessBasis           string            `json:"liveness_basis"`
	OwnerCallback           string            `json:"owner_callback"`
	CallbackEventID         string            `json:"callback_event_id,omitempty"`
	PreviousCallbackEventID string            `json:"previous_callback_event_id,omitempty"`
	Binding                 FileBinding       `json:"binding"`
}

type deliveryOwnerResumeRecord struct {
	Kind               string                        `json:"kind"`
	SchemaVersion      int                           `json:"schema_version"`
	Generation         int                           `json:"generation"`
	RunID              string                        `json:"run_id"`
	RequestSHA256      string                        `json:"request_sha256"`
	NativeSessionID    string                        `json:"native_session_id"`
	BeforeState        FileBinding                   `json:"before_state"`
	PreviousTransport  WorkflowTransport             `json:"previous_transport"`
	PreviousNextAction WorkflowAction                `json:"previous_next_action"`
	Absence            HumanOwnerRecoveryObservation `json:"absence"`
	// Liveness is absent in generations recorded before reconciliation
	// existed. Such a record is readable history only; every new start
	// observes liveness itself and never inherits a basis from history.
	Liveness               *TaskWriterLiveness `json:"writer_liveness,omitempty"`
	ProviderSession        string              `json:"provider_session"`
	ProviderExecutable     Executable          `json:"provider_executable"`
	SourceControl          Executable          `json:"source_control"`
	ControlExecutable      Executable          `json:"control_executable"`
	ContinuationGeneration int                 `json:"continuation_generation"`
	EventsAtResume         int                 `json:"events_at_resume"`
	DestinationWorkspaceID string              `json:"destination_workspace_id"`
	Confirmation           string              `json:"confirmation"`
	RecordedAtUTC          string              `json:"recorded_at_utc"`
}

const (
	ownerResumeTabAttempted   = "tab_create_attempted"
	ownerResumeStartAttempted = "agent_start_attempted"
	ownerResumePromptAttempt  = "prompt_attempted"
	ownerResumePromptSent     = "prompt_sent"
)

var ownerResumeSessionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

func ownerResumeGeneration(s workflowState) int {
	if s.OwnerResume == nil {
		return 0
	}
	n, _ := strconv.Atoi(filepath.Base(filepath.Dir(s.OwnerResume.Locator)))
	return n
}

func ownerResumeDirectory(s workflowState, generation int) string {
	return filepath.Join(s.Result.Paths.RunRoot, "owner-resume", fmt.Sprintf("%03d", generation))
}

// workflowOwnerResumeRecord validates the latest generation against its
// preserved before-state: one Run, request, session, Task and frozen context.
func workflowOwnerResumeRecord(s workflowState) (*deliveryOwnerResumeRecord, *workflowState, error) {
	if s.OwnerResume == nil {
		if s.OwnerResumeStep != "" || s.OwnerResumeCallback != "" {
			return nil, nil, workflowError(4, "owner resume step has no bound record")
		}
		return nil, nil, nil
	}
	n := ownerResumeGeneration(s)
	dir := ownerResumeDirectory(s, n)
	if n < 1 || n > 128 || s.OwnerResume.Locator != filepath.Join(dir, "resume.json") {
		return nil, nil, workflowError(4, "owner resume generation or path differs")
	}
	switch s.OwnerResumeStep {
	case ownerResumeTabAttempted, ownerResumeStartAttempted, ownerResumePromptAttempt, ownerResumePromptSent:
	default:
		return nil, nil, workflowError(4, "unrecognized owner resume step")
	}
	raw, err := workflowBound(*s.OwnerResume, 1<<20)
	if err != nil {
		return nil, nil, err
	}
	var r deliveryOwnerResumeRecord
	if err = decode(raw, 1<<20, &r); err != nil {
		return nil, nil, err
	}
	if r.Kind != "PlyDeliveryOwnerResume@1" || r.SchemaVersion != 1 || r.Generation != n || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.NativeSessionID != s.Result.Transport.AgentSessionID || r.BeforeState.Locator != filepath.Join(dir, "before-state.json") || r.ControlExecutable.Path != filepath.Join(dir, "ply-control") || r.ControlExecutable.SHA256 != r.SourceControl.SHA256 || !digestPattern.MatchString(r.Confirmation) {
		return nil, nil, workflowError(4, "owner resume record differs from the same Run, request, session or control")
	}
	before, err := workflowBound(r.BeforeState, 8<<20)
	if err != nil {
		return nil, nil, err
	}
	var old workflowState
	if err = decode(before, 8<<20, &old); err != nil {
		return nil, nil, err
	}
	if ownerResumeGeneration(old) != n-1 || !equal(old.Request, s.Request) || !equal(old.Observed, s.Observed) || old.Result.RunID != s.Result.RunID || old.Result.Handoff != s.Result.Handoff || old.Result.Paths != s.Result.Paths || old.ContextSHA256 != s.ContextSHA256 || old.Result.Transport.AgentSessionID != r.NativeSessionID || !equal(old.Acceptance, s.Acceptance) || !equal(old.StartSHA256, s.StartSHA256) || deliveryTransportIdentity(old.Result.Transport) != deliveryTransportIdentity(r.PreviousTransport) || deliveryContinuationGeneration(old) != r.ContinuationGeneration {
		return nil, nil, workflowError(4, "owner resume history changed Task, mandate, context, authority or session")
	}
	// A callback marker equal to the immutable before-state's is inherited
	// history from an earlier generation (as e3a25d3 left it); it is never this
	// generation's callback. Any other marker needs this generation's prompt.
	if s.OwnerResumeCallback != "" && s.OwnerResumeCallback != old.OwnerResumeCallback && s.OwnerResumeStep != ownerResumePromptSent {
		return nil, nil, workflowError(4, "owner resume callback precedes its continuation prompt")
	}
	return &r, &old, nil
}

// ownerResumeBase returns the state before every owner resume recorded after
// continuation generation n (n < 0: before all resumes). Earlier transitions
// validate their own transport against this base; the context is unchanged.
func ownerResumeBase(s workflowState, n int) (workflowState, error) {
	for s.OwnerResume != nil {
		r, old, err := workflowOwnerResumeRecord(s)
		if err != nil {
			return s, err
		}
		if n >= 0 && r.ContinuationGeneration < n {
			return s, nil
		}
		s = *old
	}
	return s, nil
}

func ownerResumeStartArgv(r WorkflowRequest, cwd, pane, session string) []string {
	argv := workflowStartArgv(r, cwd, pane)
	for i, a := range argv {
		if a != "--" {
			continue
		}
		resume := []string{"resume", session}
		if r.Runtime.Provider == "claude" {
			resume = []string{"--resume", session}
		}
		return append(append(append([]string{}, argv[:i+1]...), resume...), argv[i+1:]...)
	}
	return argv
}

// deliveryProviderSession finds the provider's own preserved transcript for
// the bound session. Its absence means the conversation cannot be resumed.
func deliveryProviderSession(s workflowState) (string, error) {
	id := s.Result.Transport.AgentSessionID
	if !ownerResumeSessionPattern.MatchString(id) {
		return "", workflowError(4, "bound native session identity is not a resumable provider session")
	}
	home := os.Getenv("HOME")
	switch s.Request.Runtime.Provider {
	case "claude":
		dir := os.Getenv("CLAUDE_CONFIG_DIR")
		if t := s.ClaudeTrust; t != nil && t.ConfigDir != "" {
			dir = t.ConfigDir
		} else if t != nil && t.UserHome != "" {
			dir = filepath.Join(t.UserHome, ".claude")
		} else if dir == "" {
			dir = filepath.Join(home, ".claude")
		}
		matches, err := filepath.Glob(filepath.Join(dir, "projects", "*", id+".jsonl"))
		if err == nil && len(matches) > 0 {
			return matches[0], nil
		}
		return "", workflowError(4, "no Claude transcript for session "+id+" under "+dir)
	case "codex":
		dir := ownerResumeCodexHome(s)
		found := ""
		_ = filepath.WalkDir(filepath.Join(dir, "sessions"), func(path string, e fs.DirEntry, err error) error {
			if err == nil && !e.IsDir() && strings.HasSuffix(e.Name(), "-"+id+".jsonl") {
				found = path
				return fs.SkipAll
			}
			return nil
		})
		if found != "" {
			return found, nil
		}
		return "", workflowError(4, "no Codex session for "+id+" under "+dir)
	}
	return "", workflowError(4, "provider has no supported resume adapter")
}

func ownerResumeCodexHome(s workflowState) string {
	if t := s.CodexTrust; t != nil && t.CodexHome != "" {
		return t.CodexHome
	}
	if home := os.Getenv("CODEX_HOME"); home != "" {
		return home
	}
	return filepath.Join(os.Getenv("HOME"), ".codex")
}

// deliveryResumePhase names the actual current gate; it never infers progress.
func deliveryResumePhase(s workflowState) string {
	d := s.Result.Delivery
	switch {
	case d == nil:
		return "unknown"
	case ownerRetired(s):
		return "retired"
	case d.Phase == "completed":
		return "completed"
	case d.OwnershipRelease != nil || d.HumanIntegration != nil:
		return "owner_released"
	case !deliveryStarted(s):
		return "startup_interrupted"
	case d.Attempt != nil && d.Attempt.State == "attempted":
		return "delivery_interrupted"
	case s.Result.NextAction.Actor == "user":
		return "awaiting_human"
	case deliveryCurrentCandidateQualified(s):
		return "qualified_candidate"
	}
	return "begun_work"
}

func ownerRetired(s workflowState) bool {
	r, err := deliveryRetirement(s)
	return err == nil && r != nil
}

func ownerResumeBlock(p *DeliveryOwnerResumePreview, reason, detail string, next WorkflowAction) error {
	p.State, p.Reason, p.Detail, p.NextAction, p.Confirmation = "blocked", reason, detail, next, ""
	return workflowError(4, reason+": "+detail)
}

func ownerResumePreview(d Dependencies, root string, in DeliveryOwnerResumeInput, s workflowState) (DeliveryOwnerResumePreview, error) {
	p := DeliveryOwnerResumePreview{Kind: "PlyDeliveryOwnerResumePreview@1", SchemaVersion: 1, RunID: s.Result.RunID, Provider: s.Request.Runtime.Provider, State: "blocked", Phase: deliveryResumePhase(s), NativeSessionID: s.Result.Transport.AgentSessionID, PreviousTransport: s.Result.Transport, ProviderExecutable: in.ProviderExecutable}
	coordinator := func(m string) WorkflowAction { return WorkflowAction{"coordinator", m} }
	if !deliveryRun(s.Request) || headlessRun(s.Request) || s.Result.Delivery == nil {
		return p, ownerResumeBlock(&p, "unsupported_run", "only interactive delivery-owner Runs can resume a provider conversation", coordinator("Use workflow worker show/resume for headless Runs; historical Runs keep their own commands."))
	}
	dl := s.Result.Delivery
	if dl.Phase == "completed" {
		p.State, p.NextAction = "completed", WorkflowAction{"user", "This delivery is complete; no owner is resumed and its worktree is not recreated."}
		return p, nil
	}
	if r, err := deliveryRetirement(s); err != nil {
		return p, err
	} else if r != nil {
		return p, ownerResumeBlock(&p, "owner_retired", "the coordinator retired this owner", coordinator("A fresh Task copy receives its own new mandate; this Run stays history."))
	}
	if dl.OwnershipRelease != nil || dl.HumanIntegration != nil {
		return p, ownerResumeBlock(&p, "owner_released", "candidate ownership was released to human integration", coordinator("Continue the recorded human integration plan; developer callbacks stay closed."))
	}
	if !deliveryStarted(s) {
		return p, ownerResumeBlock(&p, "startup_interrupted", "the owner never accepted its mandate, so no begun work exists", coordinator("Use ply workflow execute --spec SPEC_ID --restart for startup recovery, or execute resume RUN_ID after native onboarding."))
	}
	record, _, err := workflowOwnerResumeRecord(s)
	if err != nil {
		return p, err
	}
	pending := record != nil && s.OwnerResumeStep != ownerResumePromptSent
	if dl.Attempt != nil && dl.Attempt.State == "attempted" {
		return p, ownerResumeBlock(&p, "unknown_effect", "a reserved "+dl.Attempt.Kind+" attempt has no recorded result", coordinator("Inspect that reserved effect and its receipt before any resume; never replay it."))
	}
	stops, err := deliveryUnresolvedStops(s)
	if err != nil {
		return p, err
	}
	for _, stop := range stops {
		if e := stop.Input.EffectStatus; e == DeliveryEffectUnknown || e == DeliveryEffectPartial {
			return p, ownerResumeBlock(&p, "unknown_effect", "preserved stop "+stop.Input.EventID+" has "+e+" effects for "+stop.Input.Action, coordinator("Establish the actual effect from observations, then resolve the stop before resuming."))
		}
	}
	if !plain(p.NativeSessionID, 1, 256) {
		return p, ownerResumeBlock(&p, "session_unbound", "no native provider session is bound to this owner", coordinator("Inspect the Run; without a bound session only retirement and a fresh Task copy are supported."))
	}
	pane := s.Result.Transport.PaneID
	a, inspectErr := workflowAgentRead(d, s, time.Time{})
	if inspectErr != nil {
		var call *HerdrCallError
		if !errors.As(inspectErr, &call) || !call.ReadOnly || call.Phase != "agent get" || call.ProviderCode != "agent_not_found" {
			return p, ownerResumeBlock(&p, "inspection_unavailable", inspectErr.Error(), coordinator("Restore Herdr inspection and repeat the preview; no new writer starts while identity is unknown."))
		}
	}
	if inspectErr == nil && workflowLocation(s, a) == nil {
		switch {
		case a.Session == nil && a.CodexResumeIdentity == nil && pending && s.OwnerResumeStep == ownerResumeStartAttempted:
			detail := "the resumed provider in pane " + pane + " has no verified session; no input was sent"
			if a.SessionObservationError != nil {
				detail += ": " + a.SessionObservationError.Error()
			}
			return p, ownerResumeBlock(&p, "session_pending", detail, WorkflowAction{"coordinator", "Repeat ply workflow execute --resume to observe only this same pane. If identity remains unavailable, close no further generations: use workflow execute retire " + s.Result.RunID + " to close its exact pane, prove writer absence and preserve the source and history. Unknown effects must be reconciled first."})
		case workflowSession(s, a, false) == nil && !pending:
			p.State, p.NextAction = "live_owner", WorkflowAction{"recipient", "The same owner is live in Herdr pane " + pane + "; continue there. No tab, provider or prompt is started."}
			return p, nil
		case workflowSession(s, a, false) == nil && s.OwnerResumeStep == ownerResumeStartAttempted:
			// The same reserved attempt continues: its tab and provider exist and
			// the same session is observed. Only the reserved prompt remains.
			state, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
			if err != nil {
				return p, err
			}
			p.State, p.Generation, p.Directory, p.ControlExecutable = "continue_ready", record.Generation, filepath.Dir(s.OwnerResume.Locator), record.ControlExecutable
			p.ProviderExecutable, p.ProviderSession, p.DestinationWorkspaceID = record.ProviderExecutable, record.ProviderSession, record.DestinationWorkspaceID
			p.Confirmation = digest(struct{ Mode, State, Pane, Session string }{"continue", hash(state), pane, p.NativeSessionID})
			p.NextAction = WorkflowAction{"user", "Apply to continue the same reserved resume: the same session is live in pane " + pane + "; only its one continuation prompt is sent."}
			return p, nil
		case workflowSession(s, a, false) == nil && s.OwnerResumeStep == ownerResumePromptAttempt:
			runtime, err := workflowEffectiveRuntime(s)
			if err != nil {
				return p, err
			}
			return p, ownerResumeBlock(&p, "prompt_unconfirmed", "the continuation prompt to pane "+pane+" was attempted but its arrival is unknown; Ply never resends it", WorkflowAction{"user", "Look at Herdr pane " + pane + ": if the resumed owner received the continuation, its next callback confirms it. Otherwise ask it to run " + ShellQuote(runtime.PlyExecutable.Path) + " workflow execute show " + s.Result.RunID + " and continue."})
		case a.Session != nil && pending:
			return p, ownerResumeBlock(&p, "session_mismatch", "pane "+pane+" runs a different provider session than the bound one", WorkflowAction{"user", "Close Herdr pane " + pane + " (it is not this owner), then repeat ply workflow execute --resume."})
		}
	}
	writers, transcript := (taskWriters{d, s}).resume(in.InvokerPaneID)
	if writers.absence.PaneInventorySHA != "" {
		p.Absence = &writers.absence
	}
	p.Liveness, p.ProviderSession, p.InvokerPaneID = writers.liveness, transcript, in.InvokerPaneID
	if writers.state != taskWritersAbsent {
		return p, ownerResumeBlock(&p, writers.reason, writers.err.Error(), coordinator("Stop the named Task/session writer, or restore Herdr and local process inspection, then repeat the preview from the return worktree. No new writer starts while evidence is unknown."))
	}
	resolved, err := exec.LookPath(s.Request.Runtime.Provider)
	if err == nil {
		resolved, err = filepath.EvalSymlinks(resolved)
	}
	if err != nil || resolved != in.ProviderExecutable.Path || verifyExecutable(in.ProviderExecutable) != nil {
		return p, ownerResumeBlock(&p, "provider_unavailable", "the installed "+s.Request.Runtime.Provider+" executable differs from the proposed resume binding", coordinator("Install the same provider on PATH in this Herdr terminal and repeat the preview."))
	}
	if err = verifyExecutable(in.ControlExecutable); err != nil {
		return p, err
	}
	cwd, err := d.CWD()
	if err == nil {
		cwd, err = filepath.EvalSymlinks(cwd)
	}
	if err != nil {
		return p, err
	}
	if cwd != s.Observed.Target.WorktreeLocator && cwd != s.Observed.Epic.WorktreeLocator {
		return p, ownerResumeBlock(&p, "wrong_cwd", "resume runs from the exact Task or return worktree", WorkflowAction{"user", "Run the command from " + s.Observed.Target.WorktreeLocator + " or " + s.Observed.Epic.WorktreeLocator + "."})
	}
	p.DestinationWorkspaceID = in.DestinationWorkspaceID
	if p.DestinationWorkspaceID == "" {
		p.DestinationWorkspaceID = s.Request.Herdr.WorkspaceID
	}
	if !plain(p.DestinationWorkspaceID, 1, 128) {
		return p, workflowError(2, "invalid resume destination workspace")
	}
	raw, err := workflowCall(d, s.Request, "workspace", "get", p.DestinationWorkspaceID)
	var w struct {
		Workspace struct {
			ID string `json:"workspace_id"`
		} `json:"workspace"`
	}
	if err != nil || json.Unmarshal(raw, &w) != nil || w.Workspace.ID != p.DestinationWorkspaceID {
		return p, ownerResumeBlock(&p, "inspection_unavailable", "resume destination Herdr workspace is unavailable", coordinator("Run from an open Herdr workspace or pass --herdr-workspace with an existing workspace."))
	}
	p.Generation = ownerResumeGeneration(s) + 1
	p.Directory = ownerResumeDirectory(s, p.Generation)
	p.ControlExecutable = Executable{filepath.Join(p.Directory, "ply-control"), in.ControlExecutable.SHA256}
	if err = physical(p.Directory, true); err != nil {
		return p, err
	}
	if _, err = os.Lstat(filepath.Join(p.Directory, "resume.json")); !os.IsNotExist(err) {
		return p, workflowError(4, "owner resume generation already has preserved artifacts")
	}
	p.ResumeArgv = ownerResumeStartArgv(workflowEffectiveRequest(s), s.Observed.Target.WorktreeLocator, "<new-pane>", p.NativeSessionID)
	state, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
	if err != nil {
		return p, err
	}
	authority := writers.absence
	authority.PaneInventorySHA, authority.AgentInventorySHA = "", ""
	p.Confirmation = digest(struct {
		State                            string
		Provider, SourceControl, Control Executable
		Destination, Session, Invoker    string
		Absence                          HumanOwnerRecoveryObservation
		LivenessMethod                   string
		Remaining                        []int
	}{hash(state), in.ProviderExecutable, in.ControlExecutable, p.ControlExecutable, p.DestinationWorkspaceID, p.ProviderSession, in.InvokerPaneID, authority, writers.liveness.Method, writers.liveness.RemainingPIDs})
	p.State = "ready"
	p.NextAction = WorkflowAction{"user", "Apply this checked resume to start the same " + p.Provider + " conversation in one new Herdr tab. This preview checks startup prerequisites; the new live runtime session must still be independently observed before any continuation."}
	return p, nil
}

// WorkflowPreviewOwnerResume observes only. It never copies controls, writes
// state, creates a tab, starts a provider or sends input.
func WorkflowPreviewOwnerResume(d Dependencies, root string, in DeliveryOwnerResumeInput) (DeliveryOwnerResumePreview, error) {
	if err := containing(d, root); err != nil {
		return DeliveryOwnerResumePreview{}, err
	}
	s, err := workflowRead(root, in.RunID)
	if err != nil {
		return DeliveryOwnerResumePreview{}, err
	}
	return ownerResumePreview(d, root, in, s)
}

// WorkflowResumeOwner reserves each external effect before sending it. An
// unknown reply stays attached to this generation and is never replayed.
func WorkflowResumeOwner(d Dependencies, root string, in DeliveryOwnerResumeInput, confirmation string) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if in.Timeout <= 0 {
		return WorkflowRun{}, workflowError(2, "timeout must be positive")
	}
	var owned *workflowState
	continuing := 0
	err := withStore(root, func() error {
		s, err := workflowRead(root, in.RunID)
		if err != nil {
			return err
		}
		prior, _, err := workflowOwnerResumeRecord(s)
		if err != nil {
			return err
		} else if prior != nil && prior.Confirmation == confirmation {
			return nil
		}
		if os.Getenv("HERDR_ENV") != "1" {
			return workflowError(4, "owner resume requires a local Herdr terminal (HERDR_ENV=1)")
		}
		p, err := ownerResumePreview(d, root, in, s)
		if err == nil && p.State == "live_owner" {
			return nil
		}
		if err == nil && p.State == "continue_ready" && confirmation != "" && p.Confirmation == confirmation {
			continuing = p.Generation
			return nil
		}
		if err != nil {
			return err
		}
		if p.State != "ready" || confirmation == "" || p.Confirmation != confirmation {
			return workflowError(4, "owner resume confirmation differs from current observations; repeat the preview")
		}
		if err = verifyExecutable(p.ControlExecutable); err != nil {
			return err
		}
		before, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
		if err != nil {
			return err
		}
		// An abandoned earlier generation keeps history; the gate it displaced
		// is the one restored after this resume.
		previousNext := s.Result.NextAction
		if prior != nil && s.OwnerResumeStep != ownerResumePromptSent {
			previousNext = prior.PreviousNextAction
		}
		beforeBinding := FileBinding{filepath.Join(p.Directory, "before-state.json"), hash(before)}
		if err = d.writeOnce(beforeBinding.Locator, before); err != nil {
			return err
		}
		r := deliveryOwnerResumeRecord{Kind: "PlyDeliveryOwnerResume@1", SchemaVersion: 1, Generation: p.Generation, RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, NativeSessionID: p.NativeSessionID, BeforeState: beforeBinding, PreviousTransport: s.Result.Transport, PreviousNextAction: previousNext, Absence: *p.Absence, Liveness: p.Liveness, ProviderSession: p.ProviderSession, ProviderExecutable: in.ProviderExecutable, SourceControl: in.ControlExecutable, ControlExecutable: p.ControlExecutable, ContinuationGeneration: deliveryContinuationGeneration(s), EventsAtResume: len(s.Result.Delivery.Events), DestinationWorkspaceID: p.DestinationWorkspaceID, Confirmation: p.Confirmation, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		binding, err := workflowKeep(d, filepath.Join(p.Directory, "resume.json"), r)
		if err != nil {
			return err
		}
		// The earlier generation's callback stays in the bound before-state.
		s.OwnerResume, s.OwnerResumeStep, s.OwnerResumeCallback = &binding, ownerResumeTabAttempted, ""
		if err = d.writeValue(filepath.Join(p.Directory, "tab-create-attempt.json"), map[string]any{"resume_sha256": binding.SHA256, "argv": workflowTabCreateArgv(s, p.DestinationWorkspaceID)}); err != nil {
			return err
		}
		s.Result.Transport.State, s.Result.Transport.Observation = "unknown", "cached"
		s.Result.NextAction = WorkflowAction{"coordinator", "Follow this one owner resume; never repeat an uncertain tab, provider start or prompt."}
		if err = workflowSave(d, s); err != nil {
			return err
		}
		owned = &s
		return nil
	})
	deadline := time.Now().Add(in.Timeout)
	if err == nil && continuing > 0 {
		return deliveryReadback(d, root, in.RunID, ownerResumeFinish(d, root, in.RunID, continuing, deadline))
	}
	if err != nil || owned == nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	s := *owned
	generation := ownerResumeGeneration(s)
	record, _, err := workflowOwnerResumeRecord(s)
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if err = d.fault("owner_resume_before_tab_send"); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	b, err := workflowCallUntil(d, s.Request, deadline, workflowTabCreateArgv(s, record.DestinationWorkspaceID)...)
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	var tab struct {
		RootPane workflowAgent `json:"root_pane"`
	}
	if err = json.Unmarshal(b, &tab); err != nil {
		return deliveryReadback(d, root, in.RunID, workflowError(5, "resume tab response has no confirmed identity; do not repeat creation"))
	}
	a := tab.RootPane
	cwd := a.ForegroundCWD
	if cwd == "" {
		cwd = a.CWD
	}
	previous := record.PreviousTransport
	if a.WorkspaceID != record.DestinationWorkspaceID || !plain(a.TabID, 1, 128) || !plain(a.PaneID, 1, 128) || !plain(a.TerminalID, 1, 128) || a.PaneID == previous.PaneID || a.TerminalID == previous.TerminalID || cwd != s.Observed.Target.WorktreeLocator {
		return deliveryReadback(d, root, in.RunID, workflowError(5, "resume tab identity or Task directory is unknown; do not repeat creation"))
	}
	startArgv := ownerResumeStartArgv(workflowEffectiveRequest(s), s.Observed.Target.WorktreeLocator, a.PaneID, record.NativeSessionID)
	err = ownerResumeAdvance(d, root, in.RunID, generation, ownerResumeTabAttempted, ownerResumeStartAttempted, func(s *workflowState) error {
		transport, err := workflowKeep(d, filepath.Join(filepath.Dir(s.OwnerResume.Locator), "transport.json"), deliveryRecoveryTransport{WorkspaceID: a.WorkspaceID, TabID: a.TabID, PaneID: a.PaneID, TerminalID: a.TerminalID, CWD: cwd})
		if err != nil {
			return err
		}
		s.Result.Transport = WorkflowTransport{WorkspaceID: a.WorkspaceID, TabID: a.TabID, PaneID: a.PaneID, TerminalID: a.TerminalID, AgentSessionID: record.NativeSessionID, State: "unknown", Observation: "fresh", ObservedAt: d.Now().UTC().Format(time.RFC3339Nano)}
		if s.Request.Runtime.Provider == "codex" {
			if _, err := workflowKeep(d, filepath.Join(filepath.Dir(s.OwnerResume.Locator), "startup-deadline.json"), map[string]int64{"deadline_unix_ms": deadline.UnixMilli()}); err != nil {
				return err
			}
		}
		return d.writeValue(filepath.Join(filepath.Dir(s.OwnerResume.Locator), "agent-start-attempt.json"), map[string]any{"resume_sha256": s.OwnerResume.SHA256, "transport_sha256": transport.SHA256, "argv": startArgv})
	})
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if err = d.fault("owner_resume_before_agent_send"); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if _, startErr := workflowCallUntil(d, s.Request, deadline, startArgv...); startErr != nil {
		_ = workflowUpdate(d, root, in.RunID, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"delivery_owner_resume_start_response", startErr.Error()})
			return nil
		})
	}
	return deliveryReadback(d, root, in.RunID, ownerResumeFinish(d, root, in.RunID, generation, deadline))
}

// ownerResumeFinish sends the one reserved continuation after the same session
// is observed in the new pane. It is shared by a fresh resume and by a later
// continuation of the same reserved attempt; the step guard admits one sender.
func ownerResumeFinish(d Dependencies, root, id string, generation int, deadline time.Time) error {
	if err := ownerResumeAwaitSession(d, root, id, deadline); err != nil {
		return err
	}
	current, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	record, _, err := workflowOwnerResumeRecord(current)
	if err != nil {
		return err
	}
	if err = workflowDeliveryGuide(d, current); err != nil {
		return err
	}
	prompt, err := ownerResumePrompt(current)
	if err != nil {
		return err
	}
	promptArgv := []string{"agent", "prompt", current.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000"}
	// Reobserve at the prompt boundary; readiness evidence alone is not a
	// durable grant to a process that has since exited or switched threads.
	agent, err := workflowAgentGetUntil(d, current, false, deadline)
	if err != nil {
		return err
	}
	if !workflowSettled(agent) {
		return workflowError(4, "resumed owner is no longer idle; no input sent")
	}
	var safe startupSubmit
	safePath := filepath.Join(filepath.Dir(current.OwnerResume.Locator), "prompt-safe-submit.json")
	if current.Request.Runtime.Provider == "codex" {
		transport := ownerStartup(current)
		inspection, err := transport.inspect(d, deadline)
		if err != nil {
			return err
		}
		if inspection.Session == nil || *inspection.Session != current.Result.Transport.AgentSessionID {
			return workflowError(4, "safe startup cannot prove resumed native session; no continuation sent")
		}
		safe, err = transport.prepare(d, deadline, startupTask, &prompt, inspection)
		if err != nil {
			return err
		}
	}
	err = ownerResumeAdvance(d, root, id, generation, ownerResumeStartAttempted, ownerResumePromptAttempt, func(s *workflowState) error {
		if s.Request.Runtime.Provider == "codex" {
			if err := ownerStartup(*s).reserve(d, safePath, safe); err != nil {
				return err
			}
		}
		return d.writeValue(filepath.Join(filepath.Dir(s.OwnerResume.Locator), "prompt-attempt.json"), map[string]any{"resume_sha256": s.OwnerResume.SHA256, "argv": promptArgv, "codex_runtime_identity": agent.CodexResumeIdentity})
	})
	if err != nil {
		return err
	}
	if err = d.fault("owner_resume_before_prompt_send"); err != nil {
		return err
	}
	var callErr error
	if current.Request.Runtime.Provider == "codex" {
		callErr = ownerStartup(current).submit(d, safePath, safe)
	} else {
		_, callErr = workflowCallUntil(d, current.Request, deadline, promptArgv...)
	}
	if callErr != nil {
		_ = workflowUpdate(d, root, id, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"delivery_owner_resume_prompt_response", callErr.Error()})
			s.Result.NextAction = WorkflowAction{"user", "Look at Herdr pane " + s.Result.Transport.PaneID + ": the continuation prompt may have arrived and is never resent; repeat ply workflow execute --resume --check for the next step."}
			return nil
		})
		return callErr
	}
	return ownerResumeAdvance(d, root, id, generation, ownerResumePromptAttempt, ownerResumePromptSent, func(s *workflowState) error {
		s.Result.Transport.State = "working"
		s.Result.NextAction = record.PreviousNextAction
		return nil
	})
}

func ownerResumeAdvance(d Dependencies, root, id string, generation int, from, to string, fn func(*workflowState) error) error {
	return workflowUpdate(d, root, id, func(s *workflowState) error {
		if ownerResumeGeneration(*s) != generation || s.OwnerResumeStep != from {
			return workflowError(4, "owner resume generation or step changed; another caller owns it")
		}
		if err := fn(s); err != nil {
			return err
		}
		s.OwnerResumeStep = to
		return nil
	})
}

// Only the same native session, idle in the new pane, may receive input.
func ownerResumeAwaitSession(d Dependencies, root, id string, deadline time.Time) error {
	for {
		s, err := workflowRead(root, id)
		if err != nil {
			return err
		}
		a, err := workflowAgentRead(d, s, deadline)
		positive := true
		if err == nil && s.Request.Runtime.Provider == "codex" {
			if e := workflowLocation(s, a); e != nil {
				return e
			}
			// A guarded skip is still provider input. Establish every available
			// current/bound session fact before admitting even that action.
			if a.Session != nil || a.CodexResumeIdentity != nil {
				if e := workflowSession(s, a, false); e != nil {
					return workflowError(4, "resumed provider is not the same bound session; no input was sent: "+e.Error())
				}
			} else if a.SessionObservationError != nil {
				var call *HerdrCallError
				var failure *Error
				if errors.As(a.SessionObservationError, &call) || errors.As(a.SessionObservationError, &failure) && failure.Exit == 4 {
					return a.SessionObservationError
				}
			}
			transport := ownerStartup(s)
			inspection, e := transport.inspect(d, deadline)
			if e != nil {
				return e
			}
			if inspection.Session != nil && *inspection.Session != s.Result.Transport.AgentSessionID {
				return workflowError(4, "safe startup resumed native session differs from bound owner")
			}
			if a.Session != nil && inspection.Session == nil || a.Session == nil && a.CodexResumeIdentity == nil && inspection.Session != nil {
				return workflowError(4, "safe startup resumed session differs from fresh identity evidence")
			}
			if inspection.Surface == "update_skip_once" {
				if e = transport.skip(d, deadline, inspection); e != nil {
					return e
				}
				continue
			}
			positive = inspection.PromptReady && !inspection.LaunchPending && inspection.Session != nil
		}

		if err == nil {
			// Herdr can report the provider before its cwd or session settles.
			// Absent or unsettled is not different; only another session is.
			if e := workflowLocation(s, a); e != nil {
				err = e
			} else if a.Session == nil && a.CodexResumeIdentity == nil {
				err = workflowError(5, "resumed session not yet verified; no input was sent")
				if a.SessionObservationError != nil {
					err = a.SessionObservationError
					var call *HerdrCallError
					var failure *Error
					// Failed inspection and known mismatches are not startup lag.
					// Preserve the same attempt without automatically repeating a
					// rejected read or sending input to an unproven process.
					if errors.As(err, &call) || errors.As(err, &failure) && failure.Exit == 4 {
						return err
					}
				}
			} else if err = workflowSession(s, a, false); err != nil {
				return workflowError(4, "resumed provider is not the same bound session; no input was sent: "+err.Error())
			} else if workflowSettled(a) && positive {
				return workflowUpdate(d, root, id, func(s *workflowState) error {
					workflowObserve(d, s, a)
					return nil
				})
			}
		}
		if !time.Now().Before(deadline) {
			if err == nil {
				err = workflowError(5, "resumed session did not become idle before the deadline; no input was sent")
			}
			return workflowError(5, err.Error()+"; repeat ply workflow execute --resume to continue this same reserved attempt")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func ownerResumePrompt(s workflowState) (string, error) {
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return "", err
	}
	control := ShellQuote(runtime.PlyExecutable.Path)
	return fmt.Sprintf("Ply resumed this same Task after its previous Herdr tab closed. Your conversation, Task worktree %s, frozen mandate, private context %s, answers, candidates and recorded evidence are unchanged. From now on use only this control executable for owner callbacks: %s; the earlier control is closed. The updated callback guide is %s. First read the current native state with %s workflow execute show %s and continue from its actual next action and gate (%s). Do not repeat answered questions, valid evidence or completed effects, and never resend uncertain input.", s.Observed.Target.WorktreeLocator, ShellQuote(s.Result.Paths.Context), control, filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md"), control, s.Result.RunID, deliveryResumePhase(s)), nil
}

// ownerResumeNoteCallback is called only at owner callbacks that passed
// deliveryCallback, which requires the resumed control, exact context and the
// same live session. Coordinator events and installed-reader reports never call it.
func ownerResumeNoteCallback(s *workflowState, eventID string) {
	if s.OwnerResume == nil {
		return
	}
	_, old, err := workflowOwnerResumeRecord(*s)
	if err != nil || s.OwnerResumeCallback != "" && s.OwnerResumeCallback != old.OwnerResumeCallback {
		return
	}
	switch s.OwnerResumeStep {
	case ownerResumePromptAttempt:
		// A callback through the resumed control proves the continuation
		// arrived even though its send reply was lost.
		s.OwnerResumeStep = ownerResumePromptSent
		fallthrough
	case ownerResumePromptSent:
		s.OwnerResumeCallback = eventID
	}
}

func ownerResumeReadback(s workflowState) (*DeliveryOwnerResumeReadback, error) {
	r, old, err := workflowOwnerResumeRecord(s)
	if err != nil || r == nil {
		return nil, err
	}
	out := &DeliveryOwnerResumeReadback{PreviousCallbackEventID: old.OwnerResumeCallback, Generation: r.Generation, Step: s.OwnerResumeStep, NativeSessionID: r.NativeSessionID, PreviousTransport: r.PreviousTransport, ProviderExecutable: r.ProviderExecutable, ControlExecutable: r.ControlExecutable, Binding: *s.OwnerResume, OwnerCallback: "not_ready"}
	out.SessionObserved = s.OwnerResumeStep == ownerResumePromptAttempt || s.OwnerResumeStep == ownerResumePromptSent
	out.LivenessBasis = "not_recorded_before_reconciliation"
	if r.Liveness != nil {
		out.LivenessBasis = r.Liveness.Method
	}
	if s.OwnerResumeStep == ownerResumePromptSent {
		out.OwnerCallback = "awaiting"
		if s.OwnerResumeCallback != "" && s.OwnerResumeCallback != old.OwnerResumeCallback {
			out.OwnerCallback = "observed"
			out.CallbackEventID = s.OwnerResumeCallback
		}
	}
	return out, nil
}

// OwnerRunCandidate identifies an interactive delivery Run by its public
// Task, Spec and worktrees so a human never needs the private Run ID.
type OwnerRunCandidate struct {
	RunID          string `json:"run_id"`
	TaskID         string `json:"task_id"`
	SpecID         string `json:"spec_id"`
	Provider       string `json:"provider"`
	TaskWorktree   string `json:"task_worktree"`
	ReturnWorktree string `json:"return_worktree"`
	Phase          string `json:"phase"`
}

type OwnerRunCandidates struct {
	Runs       []OwnerRunCandidate  `json:"runs"`
	Unreadable []OwnerRunUnreadable `json:"unreadable,omitempty"`
}

// OwnerRunUnreadable names a Run that failed its native integrity read.
type OwnerRunUnreadable struct {
	RunID string `json:"run_id"`
	Error string `json:"error"`
}

// WorkflowOwnerRunCandidates reads every interactive delivery reservation in
// this workspace without observing Herdr or writing anything.
func WorkflowOwnerRunCandidates(d Dependencies, root string) (OwnerRunCandidates, error) {
	var out OwnerRunCandidates
	if err := containing(d, root); err != nil {
		return out, err
	}
	dir := filepath.Join(workflowRoot(root), "requests")
	if err := physical(dir, true); err != nil {
		return out, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return out, nil
	}
	if err != nil {
		return out, err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		s, err := workflowRead(root, id)
		if err != nil {
			// Unrelated history cannot block lookup. It stays visible, and the
			// selected Run is reread with full integrity checks before any use.
			out.Unreadable = append(out.Unreadable, OwnerRunUnreadable{RunID: id, Error: err.Error()})
			continue
		}
		if !deliveryRun(s.Request) || headlessRun(s.Request) {
			continue
		}
		var draft struct {
			Binding struct {
				TaskID string `json:"task_id"`
				SpecID string `json:"spec_id"`
			} `json:"task_spec_binding"`
		}
		if err = json.Unmarshal(s.Request.HandoffDraft, &draft); err != nil {
			return out, err
		}
		out.Runs = append(out.Runs, OwnerRunCandidate{RunID: s.Result.RunID, TaskID: draft.Binding.TaskID, SpecID: draft.Binding.SpecID, Provider: s.Request.Runtime.Provider, TaskWorktree: s.Observed.Target.WorktreeLocator, ReturnWorktree: s.Observed.Epic.WorktreeLocator, Phase: deliveryResumePhase(s)})
	}
	return out, nil
}

// ownerResumeProviderSupport states each existing adapter's coverage. Both
// resume the same session id without forking. Herdr reports Claude sessions;
// it may not report a resumed Codex session before its first turn, so a Codex
// resume can stop at session_pending without any input being sent.
func OwnerResumeProviderSupport(provider string) string {
	switch provider {
	case "claude":
		return "claude --resume SESSION_ID in a new Herdr tab (same session id; no --fork-session)"
	case "codex":
		return "codex resume SESSION_ID in a new Herdr tab; verify the same live session through Herdr or the pane-bound provider executable holding its exclusive Codex thread lock before one continuation. Unavailable or ambiguous runtime identity stops at session_pending; the coordinator can retire that exact attempt without losing work. Installed provider journey remains unverified until exercised."
	}
	return "no resume adapter"
}
