package quota

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// ClaudeSnapshotVersion is the shared v1 snapshot format (also read by AI Usage).
const ClaudeSnapshotVersion = 1

// ClaudeOptions configures ReadClaude; empty paths use the OS defaults.
type ClaudeOptions struct {
	SettingsPath string
	SnapshotPath string
}

// claudeWindow is the whitelisted shape of five_hour and seven_day.
type claudeWindow struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *float64 `json:"resets_at"`
}

type claudeRateLimits struct {
	FiveHour *claudeWindow `json:"five_hour,omitempty"`
	SevenDay *claudeWindow `json:"seven_day,omitempty"`
}

type claudeSnapshot struct {
	Version    int              `json:"version"`
	ObservedAt int64            `json:"observedAt"`
	RateLimits claudeRateLimits `json:"rateLimits"`
}

// windows requires at least one known window and a number in each present one.
func (l claudeRateLimits) windows() ([]Window, error) {
	var out []Window
	for _, known := range []struct {
		id     WindowID
		window *claudeWindow
	}{{FiveHour, l.FiveHour}, {Weekly, l.SevenDay}} {
		if known.window == nil {
			continue
		}
		if known.window.UsedPercentage == nil {
			return nil, ErrInvalidResponse
		}
		window := Window{ID: known.id, UsedPercent: *known.window.UsedPercentage}
		if known.window.ResetsAt != nil {
			reset := time.Unix(int64(*known.window.ResetsAt), 0)
			window.ResetsAt = &reset
		}
		out = append(out, window)
	}
	if len(out) == 0 {
		return nil, ErrInvalidResponse
	}
	return out, nil
}

// RecordClaudeStatusLine stores only five_hour/seven_day from Claude Code's
// status line input. Missing or invalid input writes nothing.
func RecordClaudeStatusLine(input []byte, target string, now time.Time) (bool, error) {
	var status struct {
		RateLimits *claudeRateLimits `json:"rate_limits"`
	}
	if err := json.Unmarshal(input, &status); err != nil {
		return false, fmt.Errorf("%w: status line input", ErrInvalidResponse)
	}
	if status.RateLimits == nil || (status.RateLimits.FiveHour == nil && status.RateLimits.SevenDay == nil) {
		return false, nil
	}
	if _, err := status.RateLimits.windows(); err != nil {
		return false, err
	}
	encoded, err := json.Marshal(claudeSnapshot{Version: ClaudeSnapshotVersion, ObservedAt: now.Unix(), RateLimits: *status.RateLimits})
	if err != nil {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return false, err
	}
	if err := atomicfile.Replace(target, append(encoded, '\n'), 0o600); err != nil {
		return false, err
	}
	return true, nil
}

// ReadClaude returns ErrNotConfigured without a relay, (nil, nil) while the
// relay has not written a snapshot, and ErrInvalidResponse for invalid data.
func ReadClaude(opts ClaudeOptions) (*Snapshot, error) {
	if runtime.GOOS == "windows" {
		return nil, ErrUnsupportedPlatform
	}
	settingsPath, snapshotPath := opts.SettingsPath, opts.SnapshotPath
	if settingsPath == "" {
		settingsPath = DefaultClaudeSettingsPath()
	}
	if snapshotPath == "" {
		snapshotPath = DefaultClaudeSnapshotPath()
	}
	configured, err := relayConfigured(settingsPath)
	if err != nil {
		return nil, err
	}
	if !configured {
		return nil, ErrNotConfigured
	}
	raw, err := os.ReadFile(snapshotPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snapshot claudeSnapshot
	if json.Unmarshal(raw, &snapshot) != nil || snapshot.Version != ClaudeSnapshotVersion {
		return nil, fmt.Errorf("%w: claude snapshot", ErrInvalidResponse)
	}
	windows, err := snapshot.RateLimits.windows()
	if err != nil {
		return nil, fmt.Errorf("%w: claude snapshot", err)
	}
	return &Snapshot{Source: SourceClaude, Windows: windows, ObservedAt: time.Unix(snapshot.ObservedAt, 0)}, nil
}

func relayConfigured(settingsPath string) (bool, error) {
	raw, err := os.ReadFile(settingsPath)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return false, nil
	}
	var settings struct {
		StatusLine *struct {
			Command *string `json:"command"`
		} `json:"statusLine"`
	}
	if err := json.Unmarshal(raw, &settings); err != nil {
		return false, fmt.Errorf("%w: claude settings", ErrInvalidResponse)
	}
	if settings.StatusLine == nil || settings.StatusLine.Command == nil {
		return false, nil
	}
	return isRelayCommand(*settings.StatusLine.Command)
}
