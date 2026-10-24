package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const safeStartupContract = "safe_startup_v1"
const startupSkip = "skip_once"
const startupBootstrap = "bootstrap"
const startupTask = "task"

type startupBinding struct {
	WorkspaceID string `json:"workspace_id"`
	TabID       string `json:"tab_id"`
	TerminalID  string `json:"terminal_id"`
	Name        string `json:"name"`
	CWD         string `json:"cwd"`
	Provider    string `json:"provider"`
	RunID       string `json:"workflow_run_id"`
	Generation  string `json:"generation"`
}
type startupExpected struct {
	Epoch  string `json:"server_epoch"`
	Launch string `json:"managed_launch_id"`
	Token  string `json:"evidence_token"`
}
type startupSubmission struct {
	AttemptID     string         `json:"attempt_id"`
	Kind          string         `json:"kind"`
	Binding       startupBinding `json:"binding"`
	PayloadDigest string         `json:"payload_digest"`
	Outcome       string         `json:"outcome"`
	BytesWritten  string         `json:"bytes_written"`
	Reason        string         `json:"reason"`
}
type startupInspection struct {
	Contract string `json:"contract"`
	startupExpected
	Binding       startupBinding `json:"binding"`
	ObservedAt    int64          `json:"observed_at_unix_ms"`
	Surface       string         `json:"surface"`
	PromptReady   bool           `json:"prompt_ready"`
	LaunchPending bool           `json:"launch_pending"`
	Session       *string        `json:"native_session"`
	Menu          *struct {
		Layout   string `json:"layout_id"`
		Selected int    `json:"selected_index"`
		Action   string `json:"skip_action"`
	} `json:"menu"`
	Attempt *startupSubmission `json:"attempt"`
}
type startupSubmit struct {
	Target    string          `json:"target"`
	Binding   startupBinding  `json:"binding"`
	Expected  startupExpected `json:"expected"`
	AttemptID string          `json:"attempt_id"`
	Kind      string          `json:"kind"`
	Text      *string         `json:"text"`
	Deadline  int64           `json:"deadline_unix_ms"`
}
type startupGuard struct {
	Binding  startupBinding `json:"binding"`
	Epoch    string         `json:"server_epoch"`
	Launch   string         `json:"managed_launch_id"`
	Deadline int64          `json:"deadline_unix_ms"`
}

// This adapter is shared by owner generations and independent review rounds.
// Its private files extend startup evidence without migrating frozen state.
type startupTransport struct {
	State   workflowState
	Binding startupBinding
	Target  string
	Dir     string
	Review  *deliveryReviewRecord
}

func ownerStartup(s workflowState) startupTransport {
	t := s.Result.Transport
	generation := workflowStartupGeneration(s)
	dir := filepath.Dir(workflowStartupPath(s, "startup-guard.json"))
	if s.OwnerResume != nil {
		generation = "owner-resume/" + s.OwnerResume.SHA256
		dir = filepath.Dir(s.OwnerResume.Locator)
	}
	return startupTransport{s, startupBinding{t.WorkspaceID, t.TabID, t.TerminalID, workflowAgentName(s.Result.RunID), s.Observed.Target.WorktreeLocator, s.Request.Runtime.Provider, s.Result.RunID, generation}, t.PaneID, dir, nil}
}
func reviewerStartup(s workflowState, r deliveryReviewRecord, l deliveryReviewLaunch, p deliveryReviewPane) startupTransport {
	return startupTransport{s, startupBinding{p.WorkspaceID, p.TabID, p.TerminalID, l.AgentName, l.CWD, l.Provider, s.Result.RunID, fmt.Sprintf("review/%d/%s", r.Request.Round, workflowStartupGeneration(s))}, p.PaneID, r.Dir, &r}
}
func (t startupTransport) fence(d Dependencies) error {
	current, err := workflowRead(t.State.Request.WorkspaceRoot, t.State.Result.RunID)
	if err != nil {
		return err
	}
	if workflowStartupGeneration(current) != workflowStartupGeneration(t.State) || (current.OwnerResume == nil) != (t.State.OwnerResume == nil) || current.OwnerResume != nil && current.OwnerResume.SHA256 != t.State.OwnerResume.SHA256 {
		return workflowError(4, "startup generation changed; no guarded input sent")
	}
	if workflowReportStopped(current) || deliveryObservedStop(current) {
		return workflowError(4, "startup has a preserved runtime stop; no guarded input sent")
	}
	var refusal struct {
		Reason string `json:"reason"`
	}
	raw, readErr := readFile(filepath.Join(t.Dir, "startup-refusal.json"), 8192, true)
	if readErr == nil {
		if json.Unmarshal(raw, &refusal) != nil {
			return workflowError(4, "guarded startup refusal evidence is unreadable; no input sent")
		}
		return workflowError(4, "guarded startup refusal is preserved; no dependent input sent: "+refusal.Reason)
	}
	if !os.IsNotExist(readErr) {
		return readErr
	}
	if t.Review != nil {
		x, err := deliveryTarget(d, current)
		if err != nil {
			return err
		}
		if x.OID != t.Review.Request.CandidateOID || x.Tree != t.Review.Request.CandidateTree {
			return workflowError(4, "review candidate changed before guarded input")
		}
		rounds, err := deliveryReviewRounds(current)
		if err != nil {
			return err
		}
		if len(rounds) == 0 || digest(rounds[len(rounds)-1].Request) != digest(t.Review.Request) || rounds[len(rounds)-1].Result != nil {
			return workflowError(4, "review round is no longer the current pending candidate")
		}
	}
	return nil
}

// Never let a resume replace the original send window.
func (t startupTransport) window(d Dependencies, deadline time.Time) (time.Time, error) {
	if deadline.IsZero() {
		deadline = time.Now().Add(90 * time.Second)
	}
	path := filepath.Join(t.Dir, "startup-deadline.json")
	var saved struct {
		Deadline int64 `json:"deadline_unix_ms"`
	}
	err := withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		raw, err := readFile(path, 4096, true)
		if os.IsNotExist(err) {
			saved.Deadline = deadline.UnixMilli()
			_, err = workflowKeep(d, path, saved)
			return err
		}
		if err != nil {
			return err
		}
		return json.Unmarshal(raw, &saved)
	})
	if err != nil {
		return deadline, err
	}
	if saved.Deadline <= 0 {
		return deadline, workflowError(4, "invalid preserved startup deadline")
	}
	original := time.UnixMilli(saved.Deadline)
	if original.Before(deadline) {
		deadline = original
	}
	if !time.Now().Before(deadline) {
		return deadline, workflowError(5, "Herdr startup deadline expired; consumed reservations cannot be revived")
	}
	return deadline, nil
}
func startupCapability(d Dependencies, s workflowState, deadline time.Time) error {
	raw, err := workflowCallUntil(d, s.Request, deadline, "status", "server", "--json")
	if err != nil {
		return workflowError(4, "selected Herdr server safe_startup_v1 capability unavailable; coordinator must inspect the qualified endpoint: "+err.Error())
	}
	var status struct {
		Running      bool `json:"running"`
		Compatible   bool `json:"compatible"`
		Capabilities struct {
			Safe bool `json:"safe_startup_v1"`
		} `json:"capabilities"`
	}
	if json.Unmarshal(raw, &status) != nil || !status.Running || !status.Compatible || !status.Capabilities.Safe {
		return workflowError(4, "selected Herdr server requires safe_startup_v1 (qualified macOS/Unix producer; Windows unsupported); no startup input sent")
	}
	return nil
}
func (t startupTransport) inspect(d Dependencies, deadline time.Time) (startupInspection, error) {
	var out startupInspection
	deadline, err := t.window(d, deadline)
	if err != nil {
		return out, err
	}
	if err = startupCapability(d, t.State, deadline); err != nil {
		return out, err
	}
	// Inspection uses one immutable request per exact binding. Reinspection
	// refreshes producer evidence without replacing any consumed input slot.
	path := filepath.Join(t.Dir, "startup-inspect.json")
	err = withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		_, err := workflowKeep(d, path, map[string]any{"target": t.Target, "binding": t.Binding, "attempt_id": nil})
		return err
	})
	if err != nil {
		return out, err
	}
	var raw json.RawMessage
	raw, err = workflowCallUntil(d, t.State.Request, deadline, "agent", "startup", "inspect", "--file", path)
	if err != nil {
		return out, err
	}
	var envelope struct {
		Startup json.RawMessage `json:"startup"`
	}
	if json.Unmarshal(raw, &envelope) != nil {
		return out, workflowError(4, "invalid safe startup inspection")
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(envelope.Startup, &fields) != nil {
		return out, workflowError(4, "missing safe startup inspection")
	}
	for _, key := range []string{"contract", "server_epoch", "managed_launch_id", "evidence_token", "binding", "observed_at_unix_ms", "surface", "prompt_ready", "launch_pending", "native_session", "menu", "attempt"} {
		value, ok := fields[key]
		if !ok || string(value) == "null" && key != "native_session" && key != "menu" && key != "attempt" {
			return out, workflowError(4, "safe startup inspection missing required field "+key)
		}
	}
	if json.Unmarshal(envelope.Startup, &out) != nil || out.Contract != safeStartupContract || out.Binding != t.Binding || !plain(out.Epoch, 1, 256) || !plain(out.Launch, 1, 256) || !plain(out.Token, 1, 4096) || out.ObservedAt <= 0 || time.Now().UnixMilli()-out.ObservedAt > 2000 || out.ObservedAt > time.Now().UnixMilli()+1000 {
		return out, workflowError(4, "safe startup contract, binding or current evidence differs; coordinator must inspect")
	}
	if out.Surface == "update_skip_once" {
		var menu map[string]json.RawMessage
		if json.Unmarshal(fields["menu"], &menu) != nil {
			return out, workflowError(4, "guarded skip menu is missing required fields")
		}
		for _, key := range []string{"layout_id", "selected_index", "skip_action"} {
			value, ok := menu[key]
			if !ok || string(value) == "null" {
				return out, workflowError(4, "guarded skip menu missing required field "+key)
			}
		}
	}
	if out.Surface == "provider_prompt" && out.Menu != nil {
		return out, workflowError(4, "conflicting startup menu and prompt evidence; no input sent")
	}
	if out.Surface != "provider_prompt" && out.Surface != "update_skip_once" && out.Surface != "blocked" && out.Surface != "unknown" {
		return out, workflowError(4, "unrecognized safe startup surface; no input sent")
	}
	if err = withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		return d.replaceValue(filepath.Join(t.Dir, "startup-last-inspection.json"), out)
	}); err != nil {
		return out, err
	}
	if out.Surface != "provider_prompt" && out.Surface != "update_skip_once" {
		return out, workflowError(4, "safe startup surface="+out.Surface+"; coordinator must inspect the same terminal; no input sent")
	}
	err = withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		guard := startupGuard{t.Binding, out.Epoch, out.Launch, deadline.UnixMilli()}
		path := filepath.Join(t.Dir, "startup-guard.json")
		raw, err := readFile(path, 8192, true)
		if os.IsNotExist(err) {
			_, err = workflowKeep(d, path, guard)
			return err
		}
		if err != nil {
			return err
		}
		var prior startupGuard
		if json.Unmarshal(raw, &prior) != nil || prior.Binding != guard.Binding || prior.Epoch != guard.Epoch || prior.Launch != guard.Launch {
			return workflowError(4, "safe startup epoch or managed launch changed; reservations remain consumed")
		}
		return nil
	})
	return out, err
}
func (t startupTransport) prepare(d Dependencies, deadline time.Time, kind string, text *string, inspection startupInspection) (startupSubmit, error) {
	deadline, err := t.window(d, deadline)
	p := startupSubmit{t.Target, t.Binding, inspection.startupExpected, "startup-" + digest([]string{t.Binding.RunID, t.Binding.Generation, kind})[7:39], kind, text, deadline.UnixMilli()}
	if err != nil {
		return p, err
	}
	if kind == startupSkip {
		if inspection.Surface != "update_skip_once" || inspection.Menu == nil || inspection.Menu.Layout != "codex-0.164-picker-v1" || inspection.Menu.Action != startupSkip || inspection.Menu.Selected < 0 || inspection.Menu.Selected > 2 || text != nil {
			return p, workflowError(4, "unrecognized guarded skip menu; no guessed key sent")
		}
	} else if inspection.Surface != "provider_prompt" || !inspection.PromptReady || inspection.LaunchPending || kind == startupBootstrap && inspection.Session != nil || kind == startupTask && (inspection.Session == nil || *inspection.Session == "") {
		return p, workflowError(4, "safe startup requires positive current provider prompt evidence before input")
	}
	return p, nil
}

// Called under the generation/round lock. Identical bytes never release a slot.
func (t startupTransport) reserve(d Dependencies, path string, p startupSubmit) error {
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		return workflowError(4, "guarded startup input already consumed or unknown; it is not replayed")
	}
	_, err := workflowKeep(d, path, p)
	return err
}
func (t startupTransport) submit(d Dependencies, path string, p startupSubmit) error {
	return withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		rawFile, err := readFile(path, 2<<20, true)
		if err != nil {
			return err
		}
		if hash(rawFile) != digest(p) {
			return workflowError(4, "reserved guarded request changed; no input sent")
		}
		if t.Review == nil && p.Kind == startupTask {
			current, err := workflowRead(t.State.Request.WorkspaceRoot, t.State.Result.RunID)
			if err != nil {
				return err
			}
			if err = workflowFresh(d, current, false); err != nil {
				return err
			}
			if deliveryRun(current.Request) && current.OwnerResume == nil {
				x, err := deliveryTarget(d, current)
				if err != nil {
					return err
				}
				if x.OID != current.Observed.Target.OID || x.Tree != current.Observed.Target.Tree {
					return workflowError(4, "Task start target changed before guarded input")
				}
			}
		}
		deadline := time.UnixMilli(p.Deadline)
		if !time.Now().Before(deadline) {
			return workflowError(5, "Herdr startup deadline expired before guarded input")
		}
		raw, err := workflowCallUntil(d, t.State.Request, deadline, "agent", "startup", "submit", "--file", path)
		if err != nil {
			var call *HerdrCallError
			if errors.As(err, &call) && call.ProviderCode == "permission_denied" {
				if _, keepErr := workflowKeep(d, filepath.Join(t.Dir, "startup-refusal.json"), map[string]string{"reason": err.Error(), "attempt_id": p.AttemptID}); keepErr != nil {
					return keepErr
				}
			}
			return err
		}
		var response struct {
			Submission *startupSubmission `json:"submission"`
		}
		if json.Unmarshal(raw, &response) != nil || response.Submission == nil {
			return workflowError(4, "guarded startup submission outcome missing; slot remains consumed")
		}
		result := response.Submission
		if _, err = workflowKeep(d, path+".result.json", result); err != nil {
			return err
		}
		if result.AttemptID != p.AttemptID || result.Kind != p.Kind || result.Binding != p.Binding || !digestPattern.MatchString(result.PayloadDigest) || !plain(result.Reason, 1, 2000) || result.Outcome != "submitted" || result.BytesWritten != "some" {
			outcome, bytes := result.Outcome, result.BytesWritten
			if outcome != "submitted" && outcome != "not_sent" && outcome != "unknown" {
				outcome = "unrecognized"
			}
			if bytes != "some" && bytes != "none" && bytes != "unknown" {
				bytes = "unrecognized"
			}
			return workflowError(4, "guarded startup outcome="+outcome+" bytes_written="+bytes+"; slot remains consumed; inspect without replay")
		}
		return nil
	})
}
func (t startupTransport) skip(d Dependencies, deadline time.Time, inspection startupInspection) error {
	p, err := t.prepare(d, deadline, startupSkip, nil, inspection)
	if err != nil {
		return err
	}
	path := filepath.Join(t.Dir, "startup-skip-attempt.json")
	if err = withStore(t.State.Request.WorkspaceRoot, func() error {
		if err := t.fence(d); err != nil {
			return err
		}
		return t.reserve(d, path, p)
	}); err != nil {
		return err
	}
	// A lost skip reply may be resolved by a fresh same-launch provider prompt.
	// The result or transport failure is retained, never interpreted as no bytes.
	err = t.submit(d, path, p)
	if err != nil {
		_, keepErr := workflowKeep(d, path+".failure.json", map[string]string{"error": err.Error()})
		if keepErr != nil {
			return keepErr
		}
	}
	return nil
}
