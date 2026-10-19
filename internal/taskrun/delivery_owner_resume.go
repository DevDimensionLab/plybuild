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
	PreviousTransport      WorkflowTransport              `json:"previous_transport"`
	Absence                *HumanOwnerRecoveryObservation `json:"absence,omitempty"`
	ResumeArgv             []string                       `json:"resume_argv,omitempty"`
	Confirmation           string                         `json:"confirmation"`
	NextAction             WorkflowAction                 `json:"next_action"`
}

// DeliveryOwnerResumeReadback separates the reserved start, the observed
// resumed session, the sent continuation and the first functional callback.
type DeliveryOwnerResumeReadback struct {
	Generation         int               `json:"generation"`
	Step               string            `json:"step"`
	NativeSessionID    string            `json:"native_session_id"`
	PreviousTransport  WorkflowTransport `json:"previous_transport"`
	ProviderExecutable Executable        `json:"provider_executable"`
	ControlExecutable  Executable        `json:"control_executable"`
	SessionObserved    bool              `json:"session_observed"`
	OwnerCallback      string            `json:"owner_callback"`
	CallbackEventID    string            `json:"callback_event_id,omitempty"`
	Binding            FileBinding       `json:"binding"`
}

type deliveryOwnerResumeRecord struct {
	Kind                   string                        `json:"kind"`
	SchemaVersion          int                           `json:"schema_version"`
	Generation             int                           `json:"generation"`
	RunID                  string                        `json:"run_id"`
	RequestSHA256          string                        `json:"request_sha256"`
	NativeSessionID        string                        `json:"native_session_id"`
	BeforeState            FileBinding                   `json:"before_state"`
	PreviousTransport      WorkflowTransport             `json:"previous_transport"`
	PreviousNextAction     WorkflowAction                `json:"previous_next_action"`
	Absence                HumanOwnerRecoveryObservation `json:"absence"`
	ProviderSession        string                        `json:"provider_session"`
	ProviderExecutable     Executable                    `json:"provider_executable"`
	SourceControl          Executable                    `json:"source_control"`
	ControlExecutable      Executable                    `json:"control_executable"`
	ContinuationGeneration int                           `json:"continuation_generation"`
	EventsAtResume         int                           `json:"events_at_resume"`
	DestinationWorkspaceID string                        `json:"destination_workspace_id"`
	Confirmation           string                        `json:"confirmation"`
	RecordedAtUTC          string                        `json:"recorded_at_utc"`
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
	if s.OwnerResumeCallback != "" && s.OwnerResumeStep != ownerResumePromptSent {
		return nil, nil, workflowError(4, "owner resume callback precedes its continuation prompt")
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
		dir := os.Getenv("CODEX_HOME")
		if t := s.CodexTrust; t != nil && t.CodexHome != "" {
			dir = t.CodexHome
		} else if dir == "" {
			dir = filepath.Join(home, ".codex")
		}
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

// deliveryResumePhase names the actual current gate; it never infers progress.
func deliveryResumePhase(s workflowState) string {
	d := s.Result.Delivery
	switch {
	case d == nil:
		return "unknown"
	case d.Phase == "completed":
		return "completed"
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
	if record != nil && s.OwnerResumeStep != ownerResumePromptSent {
		return p, ownerResumeBlock(&p, "resume_unknown", "owner resume generation "+strconv.Itoa(record.Generation)+" stopped at "+s.OwnerResumeStep+" with an unconfirmed effect", coordinator("Inspect Herdr and the reserved attempt in "+filepath.Dir(s.OwnerResume.Locator)+"; never repeat an uncertain tab, provider start or prompt."))
	}
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
	if _, err = workflowAgentGet(d, s, false); err == nil {
		p.State, p.NextAction = "live_owner", WorkflowAction{"recipient", "The same owner is live in Herdr pane " + s.Result.Transport.PaneID + "; continue there. No tab, provider or prompt is started."}
		return p, nil
	}
	_, err = workflowCall(d, s.Request, "pane", "get", s.Result.Transport.PaneID)
	var call *HerdrCallError
	if err == nil {
		return p, ownerResumeBlock(&p, "terminal_present", "the bound pane still exists without the observed owner session", coordinator("Inspect Herdr pane "+s.Result.Transport.PaneID+"; close it only when no writer remains, then repeat the preview."))
	}
	if !errors.As(err, &call) || call.ProviderCode != "pane_not_found" {
		return p, ownerResumeBlock(&p, "inspection_unavailable", "Herdr could not prove the bound pane absent: "+err.Error(), coordinator("Restore read access to Herdr from the Task or return worktree, then repeat the preview; absence is never assumed."))
	}
	absence, err := observeAbsentDeliveryOwner(d, s)
	if err != nil {
		reason := "inspection_unavailable"
		if strings.Contains(err.Error(), "owner_active") {
			reason = "owner_active"
		}
		return p, ownerResumeBlock(&p, reason, err.Error(), coordinator("Stop or close the remaining Task/session writer in Herdr (or restore inventory access), then repeat the preview."))
	}
	p.Absence = &absence
	if p.ProviderSession, err = deliveryProviderSession(s); err != nil {
		return p, ownerResumeBlock(&p, "provider_session_missing", err.Error(), coordinator("The conversation cannot be resumed; retire this owner with workflow execute retire RUN_ID and start a fresh Task copy."))
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
	p.ResumeArgv = ownerResumeStartArgv(s.Request, s.Observed.Target.WorktreeLocator, "<new-pane>", p.NativeSessionID)
	state, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
	if err != nil {
		return p, err
	}
	authority := absence
	authority.PaneInventorySHA, authority.AgentInventorySHA = "", ""
	p.Confirmation = digest(struct {
		State                            string
		Provider, SourceControl, Control Executable
		Destination, Session             string
		Absence                          HumanOwnerRecoveryObservation
	}{hash(state), in.ProviderExecutable, in.ControlExecutable, p.ControlExecutable, p.DestinationWorkspaceID, p.ProviderSession, authority})
	p.State = "ready"
	p.NextAction = WorkflowAction{"user", "Apply this checked resume to start the same " + p.Provider + " conversation in one new Herdr tab."}
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
	err := withStore(root, func() error {
		s, err := workflowRead(root, in.RunID)
		if err != nil {
			return err
		}
		if r, _, err := workflowOwnerResumeRecord(s); err != nil {
			return err
		} else if r != nil && r.Confirmation == confirmation {
			return nil
		}
		if os.Getenv("HERDR_ENV") != "1" {
			return workflowError(4, "owner resume requires a local Herdr terminal (HERDR_ENV=1)")
		}
		p, err := ownerResumePreview(d, root, in, s)
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
		beforeBinding := FileBinding{filepath.Join(p.Directory, "before-state.json"), hash(before)}
		if err = d.writeOnce(beforeBinding.Locator, before); err != nil {
			return err
		}
		r := deliveryOwnerResumeRecord{Kind: "PlyDeliveryOwnerResume@1", SchemaVersion: 1, Generation: p.Generation, RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, NativeSessionID: p.NativeSessionID, BeforeState: beforeBinding, PreviousTransport: s.Result.Transport, PreviousNextAction: s.Result.NextAction, Absence: *p.Absence, ProviderSession: p.ProviderSession, ProviderExecutable: in.ProviderExecutable, SourceControl: in.ControlExecutable, ControlExecutable: p.ControlExecutable, ContinuationGeneration: deliveryContinuationGeneration(s), EventsAtResume: len(s.Result.Delivery.Events), DestinationWorkspaceID: p.DestinationWorkspaceID, Confirmation: p.Confirmation, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		binding, err := workflowKeep(d, filepath.Join(p.Directory, "resume.json"), r)
		if err != nil {
			return err
		}
		s.OwnerResume, s.OwnerResumeStep = &binding, ownerResumeTabAttempted
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
	if err != nil || owned == nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	deadline := time.Now().Add(in.Timeout)
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
	startArgv := ownerResumeStartArgv(s.Request, s.Observed.Target.WorktreeLocator, a.PaneID, record.NativeSessionID)
	err = ownerResumeAdvance(d, root, in.RunID, generation, ownerResumeTabAttempted, ownerResumeStartAttempted, func(s *workflowState) error {
		transport, err := workflowKeep(d, filepath.Join(filepath.Dir(s.OwnerResume.Locator), "transport.json"), deliveryRecoveryTransport{WorkspaceID: a.WorkspaceID, TabID: a.TabID, PaneID: a.PaneID, TerminalID: a.TerminalID, CWD: cwd})
		if err != nil {
			return err
		}
		s.Result.Transport = WorkflowTransport{WorkspaceID: a.WorkspaceID, TabID: a.TabID, PaneID: a.PaneID, TerminalID: a.TerminalID, AgentSessionID: record.NativeSessionID, State: "unknown", Observation: "fresh", ObservedAt: d.Now().UTC().Format(time.RFC3339Nano)}
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
	if err = ownerResumeAwaitSession(d, root, in.RunID, deadline); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	current, err := workflowRead(root, in.RunID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if err = workflowDeliveryGuide(d, current); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	prompt, err := ownerResumePrompt(current)
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	promptArgv := []string{"agent", "prompt", current.Result.Transport.PaneID, prompt, "--wait", "--until", "working", "--timeout", "10000"}
	err = ownerResumeAdvance(d, root, in.RunID, generation, ownerResumeStartAttempted, ownerResumePromptAttempt, func(s *workflowState) error {
		return d.writeValue(filepath.Join(filepath.Dir(s.OwnerResume.Locator), "prompt-attempt.json"), map[string]any{"resume_sha256": s.OwnerResume.SHA256, "argv": promptArgv})
	})
	if err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if err = d.fault("owner_resume_before_prompt_send"); err != nil {
		return deliveryReadback(d, root, in.RunID, err)
	}
	if _, callErr := workflowCallUntil(d, s.Request, deadline, promptArgv...); callErr != nil {
		_ = workflowUpdate(d, root, in.RunID, func(s *workflowState) error {
			s.Result.Reasons = append(s.Result.Reasons, Reason{"delivery_owner_resume_prompt_response", callErr.Error()})
			s.Result.NextAction = WorkflowAction{"coordinator", "Inspect the resumed session in Herdr pane " + s.Result.Transport.PaneID + "; the continuation prompt may have arrived and is never resent."}
			return nil
		})
		return deliveryReadback(d, root, in.RunID, callErr)
	}
	err = ownerResumeAdvance(d, root, in.RunID, generation, ownerResumePromptAttempt, ownerResumePromptSent, func(s *workflowState) error {
		s.Result.Transport.State = "working"
		s.Result.NextAction = record.PreviousNextAction
		return nil
	})
	return deliveryReadback(d, root, in.RunID, err)
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
		if err == nil {
			if err = workflowIdentity(s, a, false); err != nil {
				return workflowError(4, "resumed provider is not the same bound session; no input was sent: "+err.Error())
			}
			if workflowSettled(a) {
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
			return err
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
	if s.OwnerResume != nil && s.OwnerResumeStep == ownerResumePromptSent && s.OwnerResumeCallback == "" {
		s.OwnerResumeCallback = eventID
	}
}

func ownerResumeReadback(s workflowState) (*DeliveryOwnerResumeReadback, error) {
	r, _, err := workflowOwnerResumeRecord(s)
	if err != nil || r == nil {
		return nil, err
	}
	out := &DeliveryOwnerResumeReadback{Generation: r.Generation, Step: s.OwnerResumeStep, NativeSessionID: r.NativeSessionID, PreviousTransport: r.PreviousTransport, ProviderExecutable: r.ProviderExecutable, ControlExecutable: r.ControlExecutable, Binding: *s.OwnerResume, OwnerCallback: "not_ready"}
	out.SessionObserved = s.OwnerResumeStep == ownerResumePromptAttempt || s.OwnerResumeStep == ownerResumePromptSent
	if s.OwnerResumeStep == ownerResumePromptSent {
		out.OwnerCallback = "awaiting"
		if s.OwnerResumeCallback != "" {
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
	Runs       []OwnerRunCandidate `json:"runs"`
	Unreadable []string            `json:"unreadable,omitempty"`
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
			out.Unreadable = append(out.Unreadable, id)
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
