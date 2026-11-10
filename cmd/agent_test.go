package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agents/quota"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var agentNow = time.Date(2026, 11, 7, 18, 0, 0, 0, time.UTC)

func agentAt(offset time.Duration) *time.Time {
	t := agentNow.Add(offset)
	return &t
}

func codexOK() *quota.Snapshot {
	return &quota.Snapshot{Source: quota.SourceCodex, Plan: "pro", ObservedAt: agentNow.Add(-5 * time.Second), Windows: []quota.Window{
		{ID: quota.FiveHour, UsedPercent: 38, ResetsAt: agentAt(2 * time.Hour)},
		{ID: quota.Weekly, UsedPercent: 82},
	}}
}

func claudeOK() *quota.Snapshot {
	return &quota.Snapshot{Source: quota.SourceClaude, ObservedAt: agentNow.Add(-10 * time.Minute), Windows: []quota.Window{
		{ID: quota.FiveHour, UsedPercent: 14}, {ID: quota.Weekly, UsedPercent: 61},
	}}
}

type agentFakes struct {
	codex       func(context.Context) (*quota.Snapshot, error)
	claude      func() (*quota.Snapshot, error)
	codexCalls  int
	claudeCalls int
}

func (f *agentFakes) readers() agentQuotaReaders {
	return agentQuotaReaders{
		codex: func(ctx context.Context) (*quota.Snapshot, error) {
			f.codexCalls++
			return f.codex(ctx)
		},
		claude: func() (*quota.Snapshot, error) {
			f.claudeCalls++
			return f.claude()
		},
		now: func() time.Time { return agentNow },
	}
}

func agentTestRoot(t *testing.T, agent *cobra.Command) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	root := &cobra.Command{Use: "ply", PersistentPreRunE: func(*cobra.Command, []string) error {
		t.Error("agent invoked profile initialization")
		return errors.New("profile initialization must not run")
	}}
	root.PersistentFlags().Bool("json", false, "logging")
	root.PersistentFlags().Bool("debug", false, "debug logging")
	root.AddCommand(agent)
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SilenceErrors, root.SilenceUsage = true, true
	return root, &stdout, &stderr
}

func runAgent(t *testing.T, fakes *agentFakes, args ...string) (string, string, error) {
	t.Helper()
	root, out, stderr := agentTestRoot(t, newAgentCommand(fakes.readers()))
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), stderr.String(), err
}

func TestAgentTextShowsBothProvidersAndWindowsSeparately(t *testing.T) {
	fakes := &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return codexOK(), nil }, claude: func() (*quota.Snapshot, error) { return claudeOK(), nil }}
	out, _, err := runAgent(t, fakes, "agent")
	if err != nil {
		t.Fatal(err)
	}
	want := "Codex   5-hour: 62.0% remaining (resets 2026-11-07T20:00:00Z)   Weekly: 18.0% remaining\n" +
		"        plan pro · observed 2026-11-07T17:59:55Z (5s ago)\n" +
		"Claude  5-hour: 86.0% remaining   Weekly: 39.0% remaining\n" +
		"        observed 2026-11-07T17:50:00Z (10m0s ago)\n"
	if out != want {
		t.Fatalf("text =\n%s\nwant\n%s", out, want)
	}
}

func TestAgentTextStates(t *testing.T) {
	pendingPool := codexOK()
	pendingPool.LimitReached = true
	pendingPool.ObservedAt = agentNow.Add(-121 * time.Second)
	pendingPool.Windows[0].ResetsAt = agentAt(-time.Minute)
	pendingPool.Windows = append(pendingPool.Windows, quota.Window{ID: quota.Weekly, Pool: "Spark", UsedPercent: 99})
	fakes := &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return pendingPool, nil }, claude: func() (*quota.Snapshot, error) { return nil, quota.ErrNotConfigured }}
	out, _, err := runAgent(t, fakes, "agent")
	for _, want := range []string{"5-hour: pending refresh", "Spark Weekly: 1.0% remaining", "LIMIT REACHED", "STALE", "Claude  not configured — configure claude-statusline-relay (see ply agent --help)"} {
		if err != nil || !strings.Contains(out, want) {
			t.Fatalf("missing %q (err %v):\n%s", want, err, out)
		}
	}
	fakes = &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return nil, quota.ErrCLINotFound }, claude: func() (*quota.Snapshot, error) { return nil, nil }}
	out, _, err = runAgent(t, fakes, "agent")
	if err != nil || !strings.Contains(out, "Codex   not configured — codex CLI not found") || !strings.Contains(out, "Claude  waiting for the first status line snapshot") {
		t.Fatalf("err %v:\n%s", err, out)
	}
}

func TestAgentJSONStatusesAndExit(t *testing.T) {
	stale := claudeOK()
	stale.ObservedAt = agentNow.Add(-2 * time.Hour)
	poolOnly := &quota.Snapshot{Source: quota.SourceCodex, ObservedAt: agentNow, Windows: []quota.Window{{ID: quota.Weekly, Pool: "Spark", UsedPercent: 10}}}
	cases := []struct {
		name               string
		codexErr           error
		codexSnap          *quota.Snapshot
		claudeErr          error
		claudeSnap         *quota.Snapshot
		codex, claude      string
		code               string
		stale, fail, cSnap bool
	}{
		{name: "both ok", codexSnap: codexOK(), claudeSnap: claudeOK(), codex: "ok", claude: "ok", cSnap: true},
		{name: "claude stale", codexSnap: codexOK(), claudeSnap: stale, codex: "ok", claude: "ok", stale: true, cSnap: true},
		{name: "not configured", codexErr: quota.ErrCLINotFound, claudeErr: quota.ErrNotConfigured, codex: "not_configured", claude: "not_configured"},
		{name: "waiting", codexErr: quota.ErrNoQuotaReported, codex: "not_reported", claude: "waiting"},
		{name: "pool only", codexSnap: poolOnly, claudeSnap: claudeOK(), codex: "not_reported", claude: "ok", cSnap: true},
		{name: "unsupported", codexErr: quota.ErrUnsupportedPlatform, claudeErr: quota.ErrUnsupportedPlatform, codex: "unsupported", claude: "unsupported"},
		{name: "codex timeout keeps claude", codexErr: quota.ErrTimeout, claudeSnap: claudeOK(), codex: "error", claude: "ok", code: "timeout", fail: true, cSnap: true},
		{name: "upstream", codexErr: &quota.UpstreamError{Code: -32603}, claudeErr: quota.ErrInvalidResponse, codex: "error", claude: "error", code: "upstream_error", fail: true},
		{name: "signed out", codexErr: quota.ErrNotSignedIn, claudeErr: quota.ErrNotConfigured, codex: "error", claude: "not_configured", code: "not_signed_in", fail: true},
		{name: "io", codexErr: fmt.Errorf("open /CANARY: %w", os.ErrPermission), claudeErr: quota.ErrNotConfigured, codex: "error", claude: "not_configured", code: "io_error", fail: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakes := &agentFakes{
				codex:  func(context.Context) (*quota.Snapshot, error) { return c.codexSnap, c.codexErr },
				claude: func() (*quota.Snapshot, error) { return c.claudeSnap, c.claudeErr },
			}
			out, stderr, err := runAgent(t, fakes, "agent", "--json")
			if (err != nil) != c.fail || (c.fail && ExitCode(err) != 1) {
				t.Fatalf("err = %v, want failure %v", err, c.fail)
			}
			var report agentReport
			if jsonErr := json.Unmarshal([]byte(out), &report); jsonErr != nil || report.SchemaVersion != 1 || len(report.Providers) != 2 {
				t.Fatalf("report %s: %v", out, jsonErr)
			}
			codex, claude := report.Providers[0], report.Providers[1]
			if codex.Source != "codex" || claude.Source != "claude" || codex.Status != c.codex || claude.Status != c.claude || claude.Stale != c.stale || (claude.Snapshot != nil) != c.cSnap {
				t.Fatalf("providers %s", out)
			}
			if c.code != "" && (codex.Error == nil || codex.Error.Code != c.code || codex.Error.Description == "") {
				t.Fatalf("codex error %+v", codex.Error)
			}
			if c.code == "" && codex.Error != nil {
				t.Fatalf("unexpected error %+v", codex.Error)
			}
			if strings.Contains(out+stderr, "CANARY") || fakes.codexCalls != 1 || fakes.claudeCalls != 1 {
				t.Fatalf("leak or calls %d/%d: %s %s", fakes.codexCalls, fakes.claudeCalls, out, stderr)
			}
		})
	}
}

func TestAgentJSONDocument(t *testing.T) {
	codex := codexOK()
	codex.Windows = append(codex.Windows, quota.Window{ID: quota.WindowID{Minutes: 1440}, Pool: "Spark", UsedPercent: 130, ResetsAt: agentAt(-time.Second)})
	fakes := &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return codex, nil }, claude: func() (*quota.Snapshot, error) { return nil, nil }}
	for _, args := range [][]string{{"agent", "--json"}, {"--json", "agent"}} {
		out, _, err := runAgent(t, fakes, args...)
		want := `{"schema_version":1,"providers":[` +
			`{"source":"codex","status":"ok","snapshot":{"observed_at":"2026-11-07T17:59:55Z","plan":"pro","limit_reached":false,"windows":[` +
			`{"id_minutes":300,"label":"5-hour","pool":"","used_percent":38,"remaining_percent":62,"resets_at":"2026-11-07T20:00:00Z","reset_pending":false},` +
			`{"id_minutes":10080,"label":"Weekly","pool":"","used_percent":82,"remaining_percent":18,"reset_pending":false},` +
			`{"id_minutes":1440,"label":"1-day","pool":"Spark","used_percent":130,"remaining_percent":0,"resets_at":"2026-11-07T17:59:59Z","reset_pending":true}]},"stale":false},` +
			`{"source":"claude","status":"waiting","stale":false}]}` + "\n"
		if err != nil || out != want {
			t.Fatalf("%v: err %v\n%s\nwant\n%s", args, err, out, want)
		}
	}
}

func TestAgentRejectsArgumentsBeforeReading(t *testing.T) {
	fakes := &agentFakes{}
	for _, args := range [][]string{{"agent", "extra"}, {"agent", "--bogus"}} {
		if _, _, err := runAgent(t, fakes, args...); err == nil {
			t.Fatalf("%v accepted", args)
		}
	}
	if fakes.codexCalls+fakes.claudeCalls != 0 {
		t.Fatal("readers called on argument error")
	}
}

func TestAgentPassesCommandContextAndLogsToStderr(t *testing.T) {
	var legacyLog bytes.Buffer
	previous := logrus.StandardLogger().Out
	logrus.SetOutput(&legacyLog)
	defer logrus.SetOutput(previous)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fakes := &agentFakes{
		codex: func(ctx context.Context) (*quota.Snapshot, error) {
			logrus.Warn("agent diagnostic")
			return nil, fmt.Errorf("codex quota read: %w", ctx.Err())
		},
		claude: func() (*quota.Snapshot, error) { return claudeOK(), nil },
	}
	root, out, stderr := agentTestRoot(t, newAgentCommand(fakes.readers()))
	root.SetArgs([]string{"--debug", "agent", "--json"})
	err := root.ExecuteContext(ctx)
	if ExitCode(err) != 1 || logrus.StandardLogger().Out != &legacyLog || legacyLog.Len() != 0 || !strings.Contains(stderr.String(), "agent diagnostic") {
		t.Fatalf("err %v restored %v log %q stderr %q", err, logrus.StandardLogger().Out == &legacyLog, &legacyLog, stderr)
	}
	if !json.Valid(out.Bytes()) || !strings.Contains(out.String(), `"code":"canceled"`) {
		t.Fatalf("stdout %s", out)
	}
}

func TestAgentSkipsProfileInitializationWithEmptyHome(t *testing.T) {
	home, workdir := t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(workdir)
	root := &cobra.Command{Use: "ply", PersistentPreRunE: func(cmd *cobra.Command, _ []string) error { return InitGlobals(cmd) }}
	root.PersistentFlags().Bool("json", false, "")
	root.PersistentFlags().Bool("debug", false, "")
	fakes := &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return nil, quota.ErrCLINotFound }, claude: func() (*quota.Snapshot, error) { return nil, quota.ErrNotConfigured }}
	root.AddCommand(newAgentCommand(fakes.readers()))
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"agent"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{home, workdir} {
		if entries, _ := os.ReadDir(dir); len(entries) != 0 {
			t.Fatalf("%s written: %v", dir, entries)
		}
	}
}

func TestAgentRegisteredOnRoot(t *testing.T) {
	agent, _, err := RootCmd.Find([]string{"agent"})
	if err != nil || agent.Name() != "agent" || agent.PersistentPreRunE == nil || !strings.Contains(agent.Long, "go build -o") {
		t.Fatalf("agent command %v, %v", agent, err)
	}
}
