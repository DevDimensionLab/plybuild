package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agents/quota"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	agentStatusOK            = "ok"
	agentStatusNotConfigured = "not_configured"
	agentStatusWaiting       = "waiting"
	agentStatusNotReported   = "not_reported"
	agentStatusUnsupported   = "unsupported"
	agentStatusError         = "error"
)

var errAgentQuotaFailed = errors.New("one or more quota sources failed")

type agentQuotaReaders struct {
	codex  func(context.Context) (*quota.Snapshot, error)
	claude func() (*quota.Snapshot, error)
	now    func() time.Time
}

func systemAgentQuotaReaders() agentQuotaReaders {
	return agentQuotaReaders{
		codex:  func(ctx context.Context) (*quota.Snapshot, error) { return quota.ReadCodex(ctx, quota.CodexOptions{}) },
		claude: func() (*quota.Snapshot, error) { return quota.ReadClaude(quota.ClaudeOptions{}) },
		now:    time.Now,
	}
}

type agentReport struct {
	SchemaVersion int             `json:"schema_version"`
	Providers     []agentProvider `json:"providers"`
}

type agentProvider struct {
	Source   string         `json:"source"`
	Status   string         `json:"status"`
	Snapshot *agentSnapshot `json:"snapshot,omitempty"`
	Stale    bool           `json:"stale"`
	Error    *agentError    `json:"error,omitempty"`
}

type agentSnapshot struct {
	ObservedAt   string        `json:"observed_at"`
	Plan         string        `json:"plan,omitempty"`
	LimitReached bool          `json:"limit_reached"`
	Windows      []agentWindow `json:"windows"`
}

type agentWindow struct {
	IDMinutes        int     `json:"id_minutes"`
	Label            string  `json:"label"`
	Pool             string  `json:"pool"`
	UsedPercent      float64 `json:"used_percent"`
	RemainingPercent float64 `json:"remaining_percent"`
	ResetsAt         string  `json:"resets_at,omitempty"`
	ResetPending     bool    `json:"reset_pending"`
}

type agentError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

const agentLong = `Show the remaining subscription quota reported for Codex and Claude on this
machine. This is the provider-reported account quota, not Task usage, tokens
left in a conversation, money or a guarantee that the next model call succeeds.

Codex is read directly: a short-lived 'codex app-server' answers
account/rateLimits/read with the existing Codex login (15 s deadline; no
threads, turns or inference). Values are stale after 120 s.

Claude has no quota command. Claude Code reports rate_limits only to its status
line, so the separate claude-statusline-relay records a private snapshot that
this command reads. It updates only while a Claude Code session is running and
is stale after 3600 s. Only ~/.claude/settings.json is inspected; project
status line overrides and CLAUDE_CONFIG_DIR are not detected.

One-time relay setup (not done automatically; the relay is not part of Ply
releases or Homebrew): build it from a Ply checkout to a stable absolute path,

  go build -o /stable/path/claude-statusline-relay ./cmd/claude-statusline-relay

then set statusLine.command in ~/.claude/settings.json to
'/stable/path/claude-statusline-relay' --then '<previous command>' (POSIX
single quotes) or just the quoted relay path. Undo by restoring the previous
command. See README "Agent quota" for details.

Status: ok, not_configured, waiting, not_reported, unsupported or error. The
exit status is 1 when any source has status error. --json (before or after
agent) prints one schema_version 1 document; raw provider data is never shown.`

func newAgentCommand(readers agentQuotaReaders) *cobra.Command {
	command := &cobra.Command{
		Use:               "agent",
		Short:             "Show remaining Codex and Claude subscription quota",
		Long:              agentLong,
		Example:           "  ply agent\n  ply agent --json",
		Args:              cobra.NoArgs,
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		RunE: func(cmd *cobra.Command, args []string) error {
			previousOutput := logrus.StandardLogger().Out
			logrus.SetOutput(cmd.ErrOrStderr())
			defer logrus.SetOutput(previousOutput)
			codex, codexErr := readers.codex(cmd.Context())
			claude, claudeErr := readers.claude()
			now := readers.now()
			report := agentReport{SchemaVersion: 1, Providers: []agentProvider{
				agentProviderResult(quota.SourceCodex, codex, codexErr, now),
				agentProviderResult(quota.SourceClaude, claude, claudeErr, now),
			}}
			var err error
			if workflowStatusJSON(cmd) {
				err = json.NewEncoder(cmd.OutOrStdout()).Encode(report)
			} else {
				err = writeAgentText(cmd.OutOrStdout(), report, now)
			}
			if err != nil {
				return err
			}
			for _, provider := range report.Providers {
				if provider.Status == agentStatusError {
					return errAgentQuotaFailed
				}
			}
			return nil
		},
	}
	return command
}

func agentProviderResult(source quota.Source, snapshot *quota.Snapshot, err error, now time.Time) agentProvider {
	provider := agentProvider{Source: string(source)}
	switch {
	case err == nil && snapshot == nil:
		provider.Status = agentStatusWaiting
	case err == nil:
		provider.Status = agentStatusNotReported
		if quota.Select(quota.ChooseLowest, snapshot.Windows, now).Kind != quota.SelectedNotReported {
			provider.Status = agentStatusOK
		}
		provider.Stale = snapshot.Stale(now)
		provider.Snapshot = &agentSnapshot{ObservedAt: agentTime(snapshot.ObservedAt), Plan: snapshot.Plan, LimitReached: snapshot.LimitReached, Windows: []agentWindow{}}
		for _, w := range snapshot.Windows {
			window := agentWindow{IDMinutes: w.ID.Minutes, Label: w.ID.Label(), Pool: w.Pool, UsedPercent: w.UsedPercent, RemainingPercent: w.Remaining(), ResetPending: w.ResetPending(now)}
			if w.ResetsAt != nil {
				window.ResetsAt = agentTime(*w.ResetsAt)
			}
			provider.Snapshot.Windows = append(provider.Snapshot.Windows, window)
		}
	case errors.Is(err, quota.ErrNotConfigured), errors.Is(err, quota.ErrCLINotFound):
		provider.Status = agentStatusNotConfigured
	case errors.Is(err, quota.ErrNoQuotaReported):
		provider.Status = agentStatusNotReported
	case errors.Is(err, quota.ErrUnsupportedPlatform):
		provider.Status = agentStatusUnsupported
	default:
		provider.Status = agentStatusError
		provider.Error = agentSafeError(err)
	}
	return provider
}

// agentSafeError maps errors to fixed descriptions; error text with paths or
// upstream data is never shown.
func agentSafeError(err error) *agentError {
	var upstream *quota.UpstreamError
	switch {
	case errors.As(err, &upstream):
		return &agentError{Code: "upstream_error", Description: upstream.Error()}
	case errors.Is(err, quota.ErrNotSignedIn):
		return &agentError{Code: "not_signed_in", Description: "Codex is not signed in; sign in with the Codex CLI"}
	case errors.Is(err, quota.ErrTimeout):
		return &agentError{Code: "timeout", Description: "the quota read timed out"}
	case errors.Is(err, quota.ErrInvalidResponse):
		return &agentError{Code: "invalid_response", Description: "the quota data or settings could not be read"}
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return &agentError{Code: "canceled", Description: "the quota read was canceled"}
	}
	return &agentError{Code: "io_error", Description: "reading the quota source failed"}
}

func agentTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func writeAgentText(w io.Writer, report agentReport, now time.Time) error {
	var out strings.Builder
	for _, provider := range report.Providers {
		name := map[string]string{string(quota.SourceCodex): "Codex", string(quota.SourceClaude): "Claude"}[provider.Source]
		fmt.Fprintf(&out, "%-8s", name)
		switch provider.Status {
		case agentStatusNotConfigured:
			if provider.Source == string(quota.SourceCodex) {
				out.WriteString("not configured — codex CLI not found; install Codex and sign in")
			} else {
				out.WriteString("not configured — configure claude-statusline-relay (see ply agent --help)")
			}
		case agentStatusWaiting:
			out.WriteString("waiting for the first status line snapshot from Claude Code")
		case agentStatusUnsupported:
			out.WriteString("unsupported on this platform")
		case agentStatusError:
			fmt.Fprintf(&out, "error — %s", provider.Error.Description)
		case agentStatusNotReported:
			out.WriteString("no shared quota reported")
		}
		if snapshot := provider.Snapshot; snapshot != nil {
			var parts []string
			for _, window := range snapshot.Windows {
				label := strings.TrimSpace(window.Pool + " " + window.Label)
				switch {
				case window.ResetPending:
					parts = append(parts, label+": pending refresh")
				case window.ResetsAt != "":
					parts = append(parts, fmt.Sprintf("%s: %.1f%% remaining (resets %s)", label, window.RemainingPercent, window.ResetsAt))
				default:
					parts = append(parts, fmt.Sprintf("%s: %.1f%% remaining", label, window.RemainingPercent))
				}
			}
			if provider.Status == agentStatusNotReported && len(parts) > 0 {
				out.WriteString("   ")
			}
			out.WriteString(strings.Join(parts, "   "))
			out.WriteString("\n        ")
			if snapshot.Plan != "" {
				fmt.Fprintf(&out, "plan %s · ", snapshot.Plan)
			}
			observed, _ := time.Parse(time.RFC3339, snapshot.ObservedAt)
			fmt.Fprintf(&out, "observed %s (%s ago)", snapshot.ObservedAt, now.Sub(observed).Round(time.Second))
			if provider.Stale {
				out.WriteString(" · STALE, not current quota")
			}
			if snapshot.LimitReached {
				out.WriteString(" · LIMIT REACHED")
			}
		}
		out.WriteString("\n")
	}
	_, err := io.WriteString(w, out.String())
	return err
}

func init() {
	RootCmd.AddCommand(newAgentCommand(systemAgentQuotaReaders()))
}
