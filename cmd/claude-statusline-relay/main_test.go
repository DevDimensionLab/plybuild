//go:build !windows

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/agents/quota"
)

const input = `{"session_id":"CANARY","cwd":"/tmp/'q' \u00e9 ø","rate_limits":{"five_hour":{"used_percentage":14,"resets_at":1794080400}}}` + "  \n"

func buildRelay(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "claude-statusline-relay")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return binary
}

func runRelay(t *testing.T, binary, home string, args ...string) (stdout, stderr string, exit int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = []string{"HOME=" + home, "PATH=" + os.Getenv("PATH")}
	cmd.Stdin = strings.NewReader(input)
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return out.String(), errOut.String(), exitErr.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return out.String(), errOut.String(), 0
}

func snapshotPath(home string) string {
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "AI Usage", "claude-rate-limits.json")
	}
	return filepath.Join(home, ".local", "state", "ai-usage", "claude-rate-limits.json")
}

func TestRelayPassesThroughOriginalCommand(t *testing.T) {
	binary := buildRelay(t)
	home := t.TempDir()
	copied := filepath.Join(home, "stdin.copy")
	stdout, stderr, exit := runRelay(t, binary, home, "--then", "cat > '"+copied+"'; printf X; printf E >&2")
	if stdout != "X" || stderr != "E" || exit != 0 {
		t.Fatalf("stdout %q stderr %q exit %d", stdout, stderr, exit)
	}
	if got, _ := os.ReadFile(copied); string(got) != input {
		t.Fatalf("original command stdin = %q", got)
	}
	snapshot, err := os.ReadFile(snapshotPath(home))
	if err != nil || !strings.Contains(string(snapshot), `"used_percentage":14`) || strings.Contains(string(snapshot), "CANARY") {
		t.Fatalf("snapshot %s, %v", snapshot, err)
	}
	for command, want := range map[string]int{"exit 3": 3, "definitely-not-a-command-aq01": 127, "kill -TERM $$": 128 + 15} {
		if _, _, exit := runRelay(t, binary, home, "--then", command); exit != want {
			t.Fatalf("%s: exit %d, want %d", command, exit, want)
		}
	}
}

func TestRelayRecordingFailureKeepsStatusLine(t *testing.T) {
	binary := buildRelay(t)
	home := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(home, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if stdout, stderr, exit := runRelay(t, binary, home, "--then", "cat >/dev/null; printf X"); stdout != "X" || stderr != "" || exit != 0 {
		t.Fatalf("stdout %q stderr %q exit %d", stdout, stderr, exit)
	}
	if stdout, stderr, exit := runRelay(t, binary, home); stdout != "" || stderr != "" || exit != 0 {
		t.Fatalf("without --then: stdout %q stderr %q exit %d", stdout, stderr, exit)
	}
}

func TestRelayRejectsInvalidArguments(t *testing.T) {
	binary := buildRelay(t)
	for _, args := range [][]string{{"--then"}, {"--bogus"}, {"--then", "x", "extra"}} {
		if _, stderr, exit := runRelay(t, binary, t.TempDir(), args...); exit != 2 || !strings.Contains(stderr, "usage") {
			t.Fatalf("%v: exit %d stderr %q", args, exit, stderr)
		}
	}
}

func TestReadmeRelaySetupBuildsAndIsRecognized(t *testing.T) {
	readme, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	const stablePath = "/absolute/stable/path/claude-statusline-relay"
	build := "go build -o " + stablePath + " ./cmd/claude-statusline-relay"
	var settings string
	for _, line := range strings.Split(string(readme), "\n") {
		if strings.HasPrefix(line, `{"statusLine"`) && strings.Contains(line, stablePath) {
			settings = line
		}
	}
	if !strings.Contains(string(readme), "\n"+build+"\n") || settings == "" {
		t.Fatalf("README relay setup is missing the build command or settings example")
	}
	binary := filepath.Join(t.TempDir(), "claude-statusline-relay")
	args := strings.Fields(strings.Replace(build, stablePath, binary, 1))
	command := exec.Command(args[0], args[1:]...)
	command.Dir = filepath.Join("..", "..")
	if out, err := command.CombinedOutput(); err != nil {
		t.Fatalf("README build: %v\n%s", err, out)
	}
	dir := t.TempDir()
	settingsPath := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(settings), 0o600); err != nil {
		t.Fatal(err)
	}
	snapshot, err := quota.ReadClaude(quota.ClaudeOptions{SettingsPath: settingsPath, SnapshotPath: filepath.Join(dir, "none.json")})
	if snapshot != nil || err != nil {
		t.Fatalf("README settings not recognized as configured relay: %v, %v", snapshot, err)
	}
}
