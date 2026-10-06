package taskrun

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"

	"golang.org/x/text/unicode/norm"
)

// ClaudeProjectTrust authorizes one persistent folder-trust entry, independently
// of the provider's tool permissions. Omission preserves historical launches.
type ClaudeProjectTrust struct {
	Mode         string `json:"mode"`
	WorktreeRoot string `json:"worktree_root"`
}

// ClaudeTrustPreview contains locations, never configuration contents. The
// configuration can change during a session without invalidating its callbacks.
type ClaudeTrustPreview struct {
	WorktreeRoot     string             `json:"worktree_root"`
	ProjectKey       string             `json:"project_key"`
	ConfigPath       string             `json:"config_path"`
	UserHome         string             `json:"user_home"`
	ConfigDir        string             `json:"config_dir"`
	Reason           string             `json:"reason,omitempty"`
	WorktreeIdentity *workflowTrustPath `json:"worktree_identity,omitempty"`
}

func workflowValidateClaudeTrust(r WorkflowRequest) error {
	g := r.ClaudeProjectTrust
	if g == nil {
		return nil
	}
	if r.Runtime.Provider != "claude" || g.Mode != "configuration" || !plain(g.WorktreeRoot, 1, 4096) || !filepath.IsAbs(g.WorktreeRoot) || filepath.Clean(g.WorktreeRoot) != g.WorktreeRoot {
		return workflowError(2, "claude_project_trust requires Claude, mode configuration and an absolute physical worktree_root")
	}
	return nil
}

// PreviewClaudeProjectTrust also works before prepare creates the Task worktree.
// Config resolution follows the inspected Claude Code 2.1.285 CLI: an existing
// config-directory/.config.json takes precedence over the global .claude.json.
// Unknown OAuth redirects are deliberately left to native onboarding.
func PreviewClaudeProjectTrust(g *ClaudeProjectTrust) (*ClaudeTrustPreview, error) {
	if g == nil {
		return nil, nil
	}
	if e := workflowValidateClaudeTrust(WorkflowRequest{Runtime: Runtime{Provider: "claude"}, ClaudeProjectTrust: g}); e != nil {
		return nil, e
	}
	p := &ClaudeTrustPreview{WorktreeRoot: g.WorktreeRoot, ProjectKey: norm.NFC.String(g.WorktreeRoot)}
	home, err := os.UserHomeDir()
	p.UserHome = home
	if err != nil || !filepath.IsAbs(home) || filepath.Clean(home) != home {
		p.Reason = "user home is not an absolute clean path; native onboarding is unchanged"
		return p, nil
	}
	dir, explicit := os.LookupEnv("CLAUDE_CONFIG_DIR")
	if explicit && (!filepath.IsAbs(dir) || filepath.Clean(dir) != dir) {
		p.Reason = "CLAUDE_CONFIG_DIR must be an absolute clean path; native onboarding is unchanged"
		return p, nil
	}
	p.ConfigDir = dir
	if os.Getenv("CLAUDE_CODE_CUSTOM_OAUTH_URL") != "" {
		p.Reason = "custom OAuth configuration is unsupported; native onboarding is unchanged"
		return p, nil
	}
	configDir := dir
	if !explicit {
		configDir = filepath.Join(home, ".claude")
	}
	legacy := filepath.Join(norm.NFC.String(configDir), ".config.json")
	_, err = os.Stat(legacy)
	if err == nil {
		p.ConfigPath = legacy
	} else if !os.IsNotExist(err) {
		p.Reason = "legacy configuration location cannot be inspected; native onboarding is unchanged"
	} else if explicit {
		p.ConfigPath = filepath.Join(dir, ".claude.json")
	} else {
		p.ConfigPath = filepath.Join(home, ".claude.json")
	}
	return p, nil
}

func workflowObserveClaudeTrust(r WorkflowRequest, cwd string) (*ClaudeTrustPreview, error) {
	if e := workflowValidateClaudeTrust(r); e != nil {
		return nil, e
	}
	if r.ClaudeProjectTrust == nil {
		return nil, nil
	}
	if r.ClaudeProjectTrust.WorktreeRoot != cwd {
		return nil, workflowError(4, "Claude trust grant differs from the prepared Task worktree; no trust was granted")
	}
	if e := physical(cwd, false); e != nil {
		return nil, workflowError(4, "Claude trust requires the physical prepared Task worktree; no trust was granted")
	}
	p, err := PreviewClaudeProjectTrust(r.ClaudeProjectTrust)
	if err != nil {
		return nil, err
	}
	identity, err := workflowTrustPathAt(cwd)
	if err != nil {
		return nil, err
	}
	p.WorktreeIdentity = &identity
	return p, nil
}

func workflowClaudeTrustFresh(d Dependencies, s workflowState) error {
	if s.Request.ClaudeProjectTrust == nil {
		return nil
	}
	fresh, err := workflowObserveClaudeTrust(s.Request, s.Observed.Target.WorktreeLocator)
	if err != nil {
		return err
	}
	if !equal(fresh, s.ClaudeTrust) {
		return workflowError(4, "Claude trust locations changed after preview; inspect the preserved start before any new attempt")
	}
	observe := workflowhandoff.ObserveTaskRun
	if deliveryRun(s.Request) {
		observe = workflowhandoff.ObserveDeliveryTaskRun
	}
	actual, err := observe(d.Workflow, s.Request.WorkspaceRoot, s.Request.PreparationID, s.Request.HandoffDraft)
	if err != nil || !equal(actual.Target, s.Observed.Target) || digest(actual.Preparation) != s.Request.PreparationSHA256 {
		return workflowError(4, "prepared Task identity changed before Claude trust; no trust was granted")
	}
	return nil
}

// Bind the actual interactive shell, including absence of a config override.
// Tab environment alone cannot clear an override inherited from the Herdr server
// or a shell startup file. The one reserved command creates a private nonce ack
// only after all exports succeed; terminal echo is never treated as completion.
func workflowBindClaudeEnvironment(d Dependencies, s workflowState, pane string) error {
	p := s.ClaudeTrust
	if p == nil || p.Reason != "" {
		return nil
	}
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	token := hex.EncodeToString(nonce)
	ack := filepath.Join(s.Result.Paths.RunRoot, "claude-environment-"+token+".ack")
	command := workflowClaudeEnvironmentCommand(p, token, ack)
	attempt := struct{ PaneID, Command, AckPath, Token string }{pane, command, ack, token}
	if err := d.writeValue(filepath.Join(s.Result.Paths.RunRoot, "claude-environment-attempt.json"), attempt); err != nil {
		return err
	}
	if err := workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		current.Phase = "claude_environment_attempted"
		return nil
	}); err != nil {
		return err
	}
	limit := 5 * time.Second
	if d.HerdrTimeout > 0 && d.HerdrTimeout < limit {
		limit = d.HerdrTimeout
	}
	deadline := time.Now().Add(limit)
	_, sendErr := workflowCallUntil(d, s.Request, deadline, "pane", "run", pane, command)
	for {
		b, err := readFile(ack, 128, true)
		if err == nil && string(b) == token {
			return workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
				current.Phase = "agent_start_attempted"
				if sendErr != nil {
					current.Result.Reasons = append(current.Result.Reasons, Reason{"claude_environment_response", "Shell setup response was lost; the exact private completion acknowledgement was observed."})
				}
				return nil
			})
		}
		if !time.Now().Before(deadline) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return workflowError(5, "Claude shell environment setup is unconfirmed; no trust was written or provider started. Inspect this preserved tab; do not replay its shell input.")
}

func workflowClaudeEnvironmentCommand(p *ClaudeTrustPreview, token, ack string) string {
	// Bash/zsh builtins bypass ordinary unset/export/printf aliases/functions.
	// Unsupported shells produce no ack and leave the dependent start stopped.
	command := `\builtin unset CLAUDE_CONFIG_DIR CLAUDE_CODE_CUSTOM_OAUTH_URL && \builtin export HOME=` + ShellQuote(p.UserHome)
	if p.ConfigDir != "" {
		command += ` && \builtin export CLAUDE_CONFIG_DIR=` + ShellQuote(p.ConfigDir)
	}
	wantSet := ""
	if p.ConfigDir != "" {
		wantSet = "x"
	}
	command += ` && \builtin test "$HOME" = ` + ShellQuote(p.UserHome) + ` && \builtin test "${CLAUDE_CONFIG_DIR-}" = ` + ShellQuote(p.ConfigDir) + ` && \builtin test "${CLAUDE_CONFIG_DIR+x}" = ` + ShellQuote(wantSet) + ` && \builtin test "${CLAUDE_CODE_CUSTOM_OAUTH_URL+x}" = ''`
	return command + ` && (\builtin umask 077 && \builtin set -C && \builtin printf '%s' ` + ShellQuote(token) + " > " + ShellQuote(ack) + ")"
}

// Record an uncertain attempt before the first possible config write. Repeated
// starts only show the saved run; they never call this function a second time.
func workflowLaunchClaudeTrust(d Dependencies, s workflowState) error {
	if s.ClaudeTrust == nil {
		return nil
	}
	p := s.ClaudeTrust
	effect := ClaudeTrustEffect{State: "unknown", ConfigPath: p.ConfigPath, ProjectKey: p.ProjectKey, Reason: "trust write was reserved; its outcome is not yet recorded"}
	if err := workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		if current.Result.ClaudeTrust != nil {
			return workflowError(4, "Claude trust already has an attempt; inspect the saved outcome without replay")
		}
		current.Result.ClaudeTrust = &effect
		return nil
	}); err != nil {
		return err
	}
	if p.Reason != "" {
		effect.State, effect.Reason = "skipped", p.Reason
	} else {
		effect = workflowApplyClaudeTrust(d, p.ConfigPath, p.ProjectKey)
	}
	return workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		current.Result.ClaudeTrust = &effect
		if effect.State != "written" && effect.State != "already" {
			current.Result.Reasons = append(current.Result.Reasons, Reason{"claude_trust_" + effect.State, effect.Reason + "; continuing the same Claude start. Inspect any native prompt in this tab; do not restart."})
		}
		return nil
	})
}
