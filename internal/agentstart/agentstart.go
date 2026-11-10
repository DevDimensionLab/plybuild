package agentstart

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

const (
	ReachedReserved        = "reserved"
	ReachedTabCreated      = "tab_created"
	ReachedAgentStarted    = "agent_started"
	ReachedSessionBound    = "session_bound"
	ReachedPromptSubmitted = "prompt_submitted"

	StageTabCreate  = "tab_create"
	StageAgentStart = "agent_start"
	StageAwaitReady = "await_ready"
	StagePrompt     = "prompt"

	ClassPrecondition   = "precondition"
	ClassBinaryChanged  = "binary_changed"
	ClassTabResponse    = "tab_response"
	ClassAgentResponse  = "agent_response"
	ClassNotReady       = "not_ready"
	ClassPromptResponse = "prompt_response"

	awaitReadyTimeout = 90 * time.Second
	awaitReadyPause   = 250 * time.Millisecond

	CodeInvalidInput          = "invalid_input"
	CodeInputMismatch         = "input_mismatch"
	CodeLegacyUnbound         = "legacy_unbound"
	CodeIncompleteReservation = "incomplete_reservation"
	CodeAlreadyReserved       = "already_reserved"
	CodeReceiptUnavailable    = "receipt_unavailable"
	CodeNotResumable          = "not_resumable"

	ProviderClaude = "claude"
	ProviderCodex  = "codex"
)

type HerdrBinary struct {
	Path   string
	SHA256 string
}

type StartInput struct {
	Key          string
	ReceiptPath  string
	Provider     string
	Model        *string
	Permission   *string
	CodexProfile *string
	Cwd          string
	Workspace    string
	Label        string
	Env          []string
	Prompt       string
	Herdr        HerdrBinary
	Precondition func() error
}

type StartFailure struct {
	Stage         string `json:"stage"`
	Class         string `json:"class"`
	EffectUnknown bool   `json:"effect_unknown"`
}

type Receipt struct {
	Key         string        `json:"key"`
	Provider    string        `json:"provider"`
	InputSHA256 string        `json:"input_sha256"`
	Reached     string        `json:"reached"`
	Pending     string        `json:"pending,omitempty"`
	Tab         string        `json:"tab,omitempty"`
	Pane        string        `json:"pane,omitempty"`
	Terminal    string        `json:"terminal,omitempty"`
	Session     string        `json:"session,omitempty"`
	Failure     *StartFailure `json:"failure,omitempty"`
	Legacy      bool          `json:"legacy,omitempty"`
}

type Error struct {
	Code  string
	Stage string
	Err   error
}

func (e *Error) Error() string {
	msg := e.Code
	if e.Stage != "" {
		msg += " at " + e.Stage
	}
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}
func (e *Error) Unwrap() error { return e.Err }

type ops struct {
	call         func(ctx context.Context, h HerdrBinary, args []string) (json.RawMessage, error)
	fault        func(point string) error
	awaitTimeout time.Duration
	pause        time.Duration
	lockFail     func(op string) error
}

var sha256Hex = regexp.MustCompile(`^[0-9a-f]{64}$`)

func StartAgent(ctx context.Context, in StartInput) (Receipt, error) {
	return start(ctx, in, ops{call: callHerdr, awaitTimeout: awaitReadyTimeout, pause: awaitReadyPause})
}

func ResumeAgent(ctx context.Context, in StartInput) (Receipt, error) {
	return continueStart(ctx, in, ops{call: callHerdr, awaitTimeout: awaitReadyTimeout, pause: awaitReadyPause})
}

func plain(s string) bool {
	if s == "" || len(s) > 256 || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func invalid(format string, args ...any) error {
	return &Error{Code: CodeInvalidInput, Err: fmt.Errorf(format, args...)}
}

func validate(in StartInput) error {
	switch {
	case !plain(in.Key):
		return invalid("key must be a plain non-empty string")
	case in.Provider != ProviderClaude && in.Provider != ProviderCodex:
		return invalid("provider must be claude or codex")
	case !filepath.IsAbs(in.ReceiptPath) || !filepath.IsAbs(in.Cwd) || !filepath.IsAbs(in.Herdr.Path):
		return invalid("receipt path, cwd and herdr path must be absolute")
	case !plain(in.Workspace) || !plain(in.Label):
		return invalid("workspace and label must be plain non-empty strings")
	case in.Prompt == "" || !utf8.ValidString(in.Prompt):
		return invalid("prompt must be non-empty UTF-8")
	case !sha256Hex.MatchString(in.Herdr.SHA256):
		return invalid("herdr digest must be 64 lowercase hex characters")
	case in.CodexProfile != nil && in.Provider != ProviderCodex:
		return invalid("codex profile requires the codex provider")
	}
	for _, option := range []*string{in.Model, in.Permission, in.CodexProfile} {
		if option != nil && !plain(*option) {
			return invalid("model, permission and profile must be plain non-empty strings when set")
		}
	}
	for _, env := range in.Env {
		name, _, ok := strings.Cut(env, "=")
		if !ok || name == "" || !utf8.ValidString(env) {
			return invalid("environment entries must be KEY=VALUE")
		}
	}
	return nil
}

func optional(s *string) canonicaljson.Value {
	if s == nil {
		return nil
	}
	return *s
}

func inputDigest(in StartInput) (string, error) {
	env := []canonicaljson.Value{}
	for _, entry := range in.Env {
		env = append(env, entry)
	}
	raw, err := canonicaljson.Marshal(canonicaljson.Object{
		{Name: "provider", Value: in.Provider},
		{Name: "model", Value: optional(in.Model)},
		{Name: "permission", Value: optional(in.Permission)},
		{Name: "codex_profile", Value: optional(in.CodexProfile)},
		{Name: "cwd", Value: in.Cwd},
		{Name: "workspace", Value: in.Workspace},
		{Name: "label", Value: in.Label},
		{Name: "env", Value: env},
		{Name: "prompt", Value: in.Prompt},
		{Name: "herdr_path", Value: in.Herdr.Path},
		{Name: "herdr_sha256", Value: in.Herdr.SHA256},
	})
	if err != nil {
		return "", invalid("input is not canonical JSON: %v", err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func start(ctx context.Context, in StartInput, o ops) (Receipt, error) {
	if err := validate(in); err != nil {
		return Receipt{}, err
	}
	digest, err := inputDigest(in)
	if err != nil {
		return Receipt{}, err
	}
	if err = os.MkdirAll(filepath.Dir(in.ReceiptPath), 0700); err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	f, err := os.OpenFile(in.ReceiptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, receiptMode)
	if errors.Is(err, fs.ErrExist) {
		return existing(in, digest)
	}
	if err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	if err = f.Close(); err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	s := &session{ctx: ctx, in: in, o: o, r: Receipt{Key: in.Key, Provider: in.Provider, InputSHA256: digest, Reached: ReachedReserved}}
	if err = s.hit("after_excl"); err != nil {
		return Receipt{}, err
	}
	if err = writeReceipt(in.ReceiptPath, s.r); err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	if err = s.hit("after_reserve"); err != nil {
		return s.r, err
	}
	err = s.drive()
	return s.r, err
}

func continueStart(ctx context.Context, in StartInput, o ops) (_ Receipt, err error) {
	if err := validate(in); err != nil {
		return Receipt{}, err
	}
	digest, err := inputDigest(in)
	if err != nil {
		return Receipt{}, err
	}
	lock, err := filelock.Acquire(in.ReceiptPath+".lock", filelock.Options{Fail: o.lockFail})
	if err != nil {
		return Receipt{}, &Error{Code: CodeReceiptUnavailable, Err: err}
	}
	// A release failure is reported beside the outcome; the returned receipt
	// stays the recorded one, so a retry reads it back.
	defer func() {
		if released := lock.Release(); released != nil {
			if err == nil {
				err = &Error{Code: CodeReceiptUnavailable, Err: released}
			} else {
				err = errors.Join(err, released)
			}
		}
	}()
	stored, err := bound(in, digest)
	switch {
	case err != nil:
		return stored, err
	case stored.complete():
		return stored, nil
	case !stored.resumable():
		return stored, &Error{Code: CodeNotResumable, Err: errors.New("only a finished start stopped before its prompt can continue; nothing was sent")}
	}
	s := &session{ctx: ctx, in: in, o: o, r: stored, observing: true}
	if s.r.Reached == ReachedTabCreated {
		s.r.Reached = ReachedAgentStarted
	}
	if err = s.awaitReady(); err == nil {
		err = s.prompt()
	}
	if err != nil && s.observing {
		return stored, err
	}
	return s.r, err
}

// bound reads the stored receipt and requires it to bind the same key and input.
func bound(in StartInput, digest string) (Receipt, error) {
	stored, err := ReadReceipt(in.ReceiptPath)
	switch {
	case err != nil:
		return Receipt{}, err
	case stored.Legacy:
		return stored, &Error{Code: CodeLegacyUnbound, Err: errors.New("receipt predates input binding; input equality is not claimed")}
	case stored.Key != in.Key || stored.InputSHA256 != digest:
		return stored, &Error{Code: CodeInputMismatch, Err: fmt.Errorf("stored input %s, requested %s", stored.InputSHA256, digest)}
	}
	return stored, nil
}

func existing(in StartInput, digest string) (Receipt, error) {
	stored, err := bound(in, digest)
	switch {
	case err != nil:
		return stored, err
	case stored.complete():
		return stored, nil
	case stored.Reached == ReachedReserved && stored.Pending == "" && stored.Failure == nil:
		return stored, &Error{Code: CodeIncompleteReservation, Err: errors.New("reservation never recorded a first effect marker; effect unknown")}
	}
	return stored, &Error{Code: CodeAlreadyReserved, Err: errors.New("start already owned; inspect the receipt, nothing was sent")}
}

// observing marks a continuation that has not yet written its prompt marker:
// until then the stored receipt is left byte-identical.
type session struct {
	ctx       context.Context
	in        StartInput
	o         ops
	r         Receipt
	observing bool
}

func (s *session) hit(point string) error {
	if s.o.fault == nil {
		return nil
	}
	return s.o.fault(point)
}

func (s *session) fail(stage, class string, unknown bool, cause error) error {
	var coded *Error
	if errors.As(cause, &coded) && coded.Code == ClassBinaryChanged {
		class = ClassBinaryChanged
	}
	s.r.Failure = &StartFailure{Stage: stage, Class: class, EffectUnknown: unknown}
	if !s.observing {
		_ = writeReceipt(s.in.ReceiptPath, s.r)
	}
	return &Error{Code: class, Stage: stage, Err: cause}
}

func (s *session) commit(stage string) error {
	pending := s.r.Pending
	s.r.Pending = ""
	if s.observing {
		return nil
	}
	if err := writeReceipt(s.in.ReceiptPath, s.r); err != nil {
		s.r.Pending = pending
		return &Error{Code: CodeReceiptUnavailable, Stage: stage, Err: err}
	}
	return nil
}

func (s *session) begin(stage string) error {
	if pre := s.in.Precondition; pre != nil {
		if err := pre(); err != nil {
			return s.fail(stage, ClassPrecondition, false, err)
		}
	}
	if err := s.ctx.Err(); err != nil && s.observing {
		return s.fail(stage, ClassNotReady, false, err)
	}
	s.r.Pending = stage
	failure := s.r.Failure
	if s.observing {
		s.r.Failure = nil
	}
	if err := writeReceipt(s.in.ReceiptPath, s.r); err != nil {
		s.r.Pending, s.r.Failure = "", failure
		return &Error{Code: CodeReceiptUnavailable, Stage: stage, Err: err}
	}
	s.observing = false
	return s.hit("after_pending:" + stage)
}

func (s *session) effect(stage, class string, args []string) (json.RawMessage, error) {
	if err := s.begin(stage); err != nil {
		return nil, err
	}
	raw, err := s.o.call(s.ctx, s.in.Herdr, args)
	if err != nil {
		return nil, s.fail(stage, class, true, err)
	}
	if err = s.hit("after_call:" + stage); err != nil {
		return nil, err
	}
	return raw, nil
}

func (s *session) drive() error {
	raw, err := s.effect(StageTabCreate, ClassTabResponse, tabArgs(s.in))
	if err != nil {
		return err
	}
	var tab struct {
		RootPane agentView `json:"root_pane"`
	}
	err = json.Unmarshal(raw, &tab)
	root := tab.RootPane
	if err != nil || !plain(root.PaneID) || !plain(root.TabID) || !plain(root.TerminalID) || root.WorkspaceID != s.in.Workspace {
		return s.fail(StageTabCreate, ClassTabResponse, true, errors.New("created tab identity is unknown"))
	}
	s.r.Tab, s.r.Pane, s.r.Terminal, s.r.Reached = root.TabID, root.PaneID, root.TerminalID, ReachedTabCreated
	if err = s.commit(StageTabCreate); err != nil {
		return err
	}
	if _, err = s.effect(StageAgentStart, ClassAgentResponse, startArgs(s.in, s.r.Pane)); err != nil {
		return err
	}
	s.r.Reached = ReachedAgentStarted
	if err = s.commit(StageAgentStart); err != nil {
		return err
	}
	if err = s.awaitReady(); err != nil {
		return err
	}
	return s.prompt()
}

func (s *session) prompt() error {
	if _, err := s.effect(StagePrompt, ClassPromptResponse, promptArgs(s.r.Pane, s.in.Prompt)); err != nil {
		return err
	}
	s.r.Reached = ReachedPromptSubmitted
	return s.commit(StagePrompt)
}

type agentView struct {
	WorkspaceID      string `json:"workspace_id"`
	TabID            string `json:"tab_id"`
	PaneID           string `json:"pane_id"`
	TerminalID       string `json:"terminal_id"`
	ForegroundCWD    string `json:"foreground_cwd"`
	Name             string `json:"name"`
	Agent            string `json:"agent"`
	Status           string `json:"agent_status"`
	InteractiveReady *bool  `json:"interactive_ready"`
	LaunchPending    *bool  `json:"launch_pending"`
	Session          *struct {
		Agent string `json:"agent"`
		Kind  string `json:"kind"`
		Value string `json:"value"`
	} `json:"agent_session"`
}

func (a agentView) settled() bool {
	return (a.Status == "idle" || a.Status == "done") && (a.LaunchPending == nil || !*a.LaunchPending)
}

// A native session is not required to address the agent: the prompt goes to the
// pane, whose Herdr name is cleared when the agent is replaced. Without one, the
// agent must report that it is interactive and where it runs.
func (a agentView) sessionlessReady() bool {
	return a.InteractiveReady != nil && *a.InteractiveReady && a.ForegroundCWD != ""
}

func (s *session) identity(a agentView) error {
	r, in := s.r, s.in
	switch {
	case a.PaneID != r.Pane || a.Name != agentName(in.Key) || a.Agent != in.Provider:
		return errors.New("observed pane, name or agent differs from the reserved start")
	case a.TerminalID != r.Terminal || a.TabID != r.Tab:
		return errors.New("observed tab or terminal differs from the created tab")
	case a.WorkspaceID != in.Workspace:
		return errors.New("observed workspace differs from the requested workspace")
	case (in.Provider == ProviderClaude || a.ForegroundCWD != "") && a.ForegroundCWD != in.Cwd:
		return errors.New("observed foreground cwd differs from the requested cwd")
	case a.Session != nil && (a.Session.Agent != in.Provider || a.Session.Kind != "id" || !plain(a.Session.Value)):
		return errors.New("observed native session does not belong to the requested provider")
	}
	return nil
}

func (s *session) awaitReady() error {
	deadline := time.Now().Add(s.o.awaitTimeout)
	observed := s.r.Session
	for {
		raw, err := s.o.call(s.ctx, s.in.Herdr, []string{"agent", "get", s.r.Pane})
		if err != nil {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
		}
		var got struct {
			Agent agentView `json:"agent"`
		}
		if err = json.Unmarshal(raw, &got); err != nil {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, errors.New("agent get returned an unreadable result"))
		}
		a := got.Agent
		if err = s.identity(a); err != nil {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, err)
		}
		if observed != "" && (a.Session == nil || a.Session.Value != observed) {
			return s.fail(StageAwaitReady, ClassAgentResponse, false, errors.New("native session changed during startup"))
		}
		if a.Session != nil {
			observed = a.Session.Value
			if a.settled() {
				s.r.Session, s.r.Reached = observed, ReachedSessionBound
				return s.commit(StageAwaitReady)
			}
		} else if a.settled() && a.sessionlessReady() {
			return nil
		}
		if !time.Now().Before(deadline) {
			return s.fail(StageAwaitReady, ClassNotReady, false, errors.New("deadline expired before the agent was observed ready"))
		}
		select {
		case <-s.ctx.Done():
			return s.fail(StageAwaitReady, ClassNotReady, false, s.ctx.Err())
		case <-time.After(s.o.pause):
		}
	}
}
