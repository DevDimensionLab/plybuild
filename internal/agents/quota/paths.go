package quota

import (
	"os"
	"path/filepath"
	"runtime"
)

const claudeSnapshotFile = "claude-rate-limits.json"

// homeDir prefers $HOME so callers and tests can isolate it.
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	home, _ := os.UserHomeDir()
	return home
}

// DefaultClaudeSettingsPath is ~/.claude/settings.json. Project settings and
// CLAUDE_CONFIG_DIR are not discovered in v1.
func DefaultClaudeSettingsPath() string {
	return filepath.Join(homeDir(), ".claude", "settings.json")
}

// DefaultClaudeSnapshotPath is shared with AI Usage on macOS and uses
// $XDG_STATE_HOME/ai-usage on Linux.
func DefaultClaudeSnapshotPath() string {
	return claudeSnapshotPath(runtime.GOOS, homeDir(), os.Getenv("XDG_STATE_HOME"))
}

func claudeSnapshotPath(goos, home, xdgState string) string {
	if goos == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "AI Usage", claudeSnapshotFile)
	}
	if xdgState == "" {
		xdgState = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(xdgState, "ai-usage", claudeSnapshotFile)
}
