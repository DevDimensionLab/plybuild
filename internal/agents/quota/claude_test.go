package quota

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

const statusLineInput = `{"session_id":"CANARY-session","cwd":"/CANARY/cwd","transcript_path":"/CANARY/t.jsonl",
"model":{"id":"m"},"rate_limits":{"five_hour":{"used_percentage":23.5,"resets_at":1794080400,"CANARY_sub":1},
"seven_day":{"used_percentage":90,"resets_at":null},"spend_limit":{"CANARY":true}}}`

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestRecordClaudeStatusLineWritesOnlyWhitelistedWindows(t *testing.T) {
	target := filepath.Join(t.TempDir(), "new", "claude-rate-limits.json")
	recorded, err := RecordClaudeStatusLine([]byte(statusLineInput), target, testNow)
	if err != nil || !recorded {
		t.Fatalf("recorded %v, %v", recorded, err)
	}
	written := readFile(t, target)
	want := `{"version":1,"observedAt":1794070000,"rateLimits":{"five_hour":{"used_percentage":23.5,"resets_at":1794080400},"seven_day":{"used_percentage":90,"resets_at":null}}}`
	if strings.TrimSpace(written) != want || strings.Contains(written, "CANARY") {
		t.Fatalf("snapshot = %s", written)
	}
	for path, mode := range map[string]os.FileMode{target: 0o600, filepath.Dir(target): 0o700} {
		if info, err := os.Stat(path); err != nil || info.Mode().Perm() != mode {
			t.Fatalf("%s mode %v, %v; want %v", path, info.Mode().Perm(), err, mode)
		}
	}
}

func TestRecordClaudeStatusLineKeepsPreviousSnapshotOnBadInput(t *testing.T) {
	inputs := map[string]string{
		"invalid json":      `{"rate_limits":`,
		"no rate limits":    `{"session_id":"x"}`,
		"empty rate limits": `{"rate_limits":{}}`,
		"only unknown":      `{"rate_limits":{"spend_limit":{"used_percentage":1}}}`,
		"string used":       `{"rate_limits":{"five_hour":{"used_percentage":"5","resets_at":1}}}`,
		"missing used":      `{"rate_limits":{"five_hour":{"resets_at":1}}}`,
		"one window bad":    `{"rate_limits":{"five_hour":{"used_percentage":"x"},"seven_day":{"used_percentage":1}}}`,
		"empty":             ``,
	}
	for name, input := range inputs {
		t.Run(name, func(t *testing.T) {
			target := filepath.Join(t.TempDir(), "claude-rate-limits.json")
			previous := `{"version":1,"observedAt":1,"rateLimits":{"seven_day":{"used_percentage":1,"resets_at":null}}}`
			if err := os.WriteFile(target, []byte(previous), 0o600); err != nil {
				t.Fatal(err)
			}
			recorded, err := RecordClaudeStatusLine([]byte(input), target, testNow)
			if recorded || readFile(t, target) != previous {
				t.Fatalf("recorded %v, %v; file %s", recorded, err, readFile(t, target))
			}
			if err != nil && strings.Contains(err.Error(), "CANARY") {
				t.Fatalf("error leaks input: %v", err)
			}
			entries, _ := os.ReadDir(filepath.Dir(target))
			if len(entries) != 1 {
				t.Fatalf("left files %v", entries)
			}
		})
	}
}

func writeSettings(t *testing.T, dir, command string) string {
	t.Helper()
	path := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(path, []byte(command), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadClaudeRelayDetection(t *testing.T) {
	cases := map[string]error{
		``:                                  ErrNotConfigured,
		`{"model":"x"}`:                     ErrNotConfigured,
		`{"statusLine":{"type":"command"}}`: ErrNotConfigured,
		`{"statusLine":{"command":"~/.claude/statusline.sh"}}`:                                                ErrNotConfigured,
		`{"statusLine":{"command":"echo claude-statusline-relay"}}`:                                           ErrNotConfigured,
		`{"statusLine":{"command":"/bin/claude-statusline-relay-old"}}`:                                       ErrNotConfigured,
		`{"statusLine":{"command":"'/opt/my tools/claude-statusline-relay' --then 'printf '\\''%s'\\'' x'"}}`: nil,
		`{"statusLine":{"command":"\"/a \\\"b\\\"/claude-statusline-relay\""}}`:                               nil,
		`{"statusLine":{"command":"  /a\\ b/claude-statusline-relay"}}`:                                       nil,
		`{"statusLine":{"command":"claude-statusline-relay"}}`:                                                nil,
		`{"statusLine":{"command":"'/unterminated/claude-statusline-relay"}}`:                                 ErrInvalidResponse,
		`{"statusLine":{"command":7}}`:                                                                        ErrInvalidResponse,
		`{"statusLine":"x"}`:                                                                                  ErrInvalidResponse,
		`[1]`:                                                                                                 ErrInvalidResponse,
		`{"statusLine":`:                                                                                      ErrInvalidResponse,
	}
	for settings, want := range cases {
		dir := t.TempDir()
		got, err := ReadClaude(ClaudeOptions{SettingsPath: writeSettings(t, dir, settings), SnapshotPath: filepath.Join(dir, "missing.json")})
		if !errors.Is(err, want) || got != nil {
			t.Errorf("settings %s: %v, %v; want %v", settings, got, err, want)
		}
	}
	_, err := ReadClaude(ClaudeOptions{SettingsPath: filepath.Join(t.TempDir(), "absent.json"), SnapshotPath: "unused"})
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("missing settings: %v", err)
	}
}

func TestReadClaudeSnapshot(t *testing.T) {
	reset := time.Unix(1794080400, 0)
	cases := map[string]struct {
		file    string
		windows []Window
		err     error
	}{
		"v1 shared format": {`{"version":1,"observedAt":1794074413,"rateLimits":{"five_hour":{"used_percentage":14,"resets_at":1794080400},"seven_day":{"used_percentage":61,"resets_at":1794080400}},"extra":1}`,
			[]Window{{ID: FiveHour, UsedPercent: 14, ResetsAt: &reset}, {ID: Weekly, UsedPercent: 61, ResetsAt: &reset}}, nil},
		"weekly only, null reset": {`{"version":1,"observedAt":1794074413,"rateLimits":{"seven_day":{"used_percentage":0,"resets_at":null}}}`,
			[]Window{{ID: Weekly, UsedPercent: 0}}, nil},
		"unknown version": {`{"version":2,"observedAt":1,"rateLimits":{"seven_day":{"used_percentage":1}}}`, nil, ErrInvalidResponse},
		"no windows":      {`{"version":1,"observedAt":1,"rateLimits":{}}`, nil, ErrInvalidResponse},
		"missing used":    {`{"version":1,"observedAt":1,"rateLimits":{"five_hour":{"resets_at":1}}}`, nil, ErrInvalidResponse},
		"corrupt":         {`{"version":1,`, nil, ErrInvalidResponse},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			snapshotPath := filepath.Join(dir, "claude-rate-limits.json")
			if err := os.WriteFile(snapshotPath, []byte(c.file), 0o600); err != nil {
				t.Fatal(err)
			}
			settings := writeSettings(t, dir, `{"statusLine":{"command":"/x/claude-statusline-relay"}}`)
			got, err := ReadClaude(ClaudeOptions{SettingsPath: settings, SnapshotPath: snapshotPath})
			if !errors.Is(err, c.err) {
				t.Fatalf("err = %v, want %v", err, c.err)
			}
			if c.err != nil {
				return
			}
			if got.Source != SourceClaude || !got.ObservedAt.Equal(time.Unix(1794074413, 0)) || !reflect.DeepEqual(got.Windows, c.windows) {
				t.Fatalf("snapshot = %+v", got)
			}
		})
	}
}

func TestReadClaudeConfiguredWithoutSnapshotIsWaiting(t *testing.T) {
	dir := t.TempDir()
	settings := writeSettings(t, dir, `{"statusLine":{"command":"/x/claude-statusline-relay"}}`)
	got, err := ReadClaude(ClaudeOptions{SettingsPath: settings, SnapshotPath: filepath.Join(dir, "none.json")})
	if got != nil || err != nil {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestClaudeDefaultPaths(t *testing.T) {
	cases := []struct{ goos, xdg, want string }{
		{"darwin", "/xdg", "/h/Library/Application Support/AI Usage/claude-rate-limits.json"},
		{"linux", "", "/h/.local/state/ai-usage/claude-rate-limits.json"},
		{"linux", "/xdg", "/xdg/ai-usage/claude-rate-limits.json"},
	}
	for _, c := range cases {
		if got := claudeSnapshotPath(c.goos, "/h", c.xdg); got != c.want {
			t.Errorf("%s %q: %s", c.goos, c.xdg, got)
		}
	}
}
