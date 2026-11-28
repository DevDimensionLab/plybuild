package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agents/quota"
	"github.com/devdimensionlab/plybuild/internal/quotahistory"
)

type recorded struct {
	series string
	point  quotahistory.Point
}

func TestAgentRecordsTheSharedWindowsInTheHistory(t *testing.T) {
	var got []recorded
	fakes := &agentFakes{
		codex: func(context.Context) (*quota.Snapshot, error) {
			s := codexOK()
			s.Windows = append(s.Windows, quota.Window{ID: quota.Weekly, Pool: "spark", UsedPercent: 3}, quota.Window{ID: quota.WindowID{Minutes: 60}, UsedPercent: 9})
			return s, nil
		},
		claude: func() (*quota.Snapshot, error) { return claudeOK(), nil },
	}
	readers := fakes.readers()
	readers.record = func(series string, p quotahistory.Point) error {
		got = append(got, recorded{series, p})
		return nil
	}
	root, _, _ := agentTestRoot(t, newAgentCommand(readers))
	root.SetArgs([]string{"agent", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	codexAt, claudeAt := agentNow.Add(-5*time.Second), agentNow.Add(-10*time.Minute)
	want := []recorded{
		{"codex.five_hour", quotahistory.Point{At: codexAt, Used: 38, ResetsAt: agentAt(2 * time.Hour)}},
		{"codex.weekly", quotahistory.Point{At: codexAt, Used: 82}},
		{"claude.five_hour", quotahistory.Point{At: claudeAt, Used: 14}},
		{"claude.weekly", quotahistory.Point{At: claudeAt, Used: 61}},
	}
	if len(got) != len(want) {
		t.Fatalf("recorded %+v", got)
	}
	for i := range want {
		g, w := got[i], want[i]
		if g.series != w.series || !g.point.At.Equal(w.point.At) || g.point.Used != w.point.Used || (g.point.ResetsAt == nil) != (w.point.ResetsAt == nil) ||
			g.point.ResetsAt != nil && !g.point.ResetsAt.Equal(*w.point.ResetsAt) || g.point.Limit != nil {
			t.Fatalf("%d: %+v, want %+v", i, g, w)
		}
	}
}

func TestAgentOutputIgnoresAFailedHistoryWrite(t *testing.T) {
	run := func(record func(string, quotahistory.Point) error) (string, error) {
		fakes := &agentFakes{codex: func(context.Context) (*quota.Snapshot, error) { return codexOK(), nil }, claude: func() (*quota.Snapshot, error) { return claudeOK(), nil }}
		readers := fakes.readers()
		readers.record = record
		root, out, _ := agentTestRoot(t, newAgentCommand(readers))
		root.SetArgs([]string{"agent", "--json"})
		err := root.Execute()
		return out.String(), err
	}
	plain, err := run(nil)
	if err != nil {
		t.Fatal(err)
	}
	failing, err := run(func(string, quotahistory.Point) error { return errors.New("read-only file system") })
	if err != nil || failing != plain {
		t.Fatalf("output changed: %q vs %q, %v", failing, plain, err)
	}
}

func TestGitHubActionsQuotaRecordsTheMonthInTheHistory(t *testing.T) {
	var got []recorded
	env := githubActionsTestEnv(githubActionsHome(t, 3000), nil)
	env.record = func(series string, p quotahistory.Point) error {
		got = append(got, recorded{series, p})
		return nil
	}
	if _, err := executeGitHub(t, env, "actions", "quota", "--format", "json"); err != nil {
		t.Fatal(err)
	}
	reset := time.Date(2026, time.December, 1, 0, 0, 0, 0, time.UTC)
	if len(got) != 1 || got[0].series != "github.actions" || !got[0].point.At.Equal(githubActionsNow) || got[0].point.Used != 1308 ||
		got[0].point.Limit == nil || *got[0].point.Limit != 3000 || got[0].point.ResetsAt == nil || !got[0].point.ResetsAt.Equal(reset) {
		t.Fatalf("recorded %+v", got)
	}
	got = nil
	if _, err := executeGitHub(t, githubActionsTestEnvRecording(githubActionsHome(t, 0), errors.New("gh failed"), &got), "actions", "quota", "--format", "json"); err == nil || len(got) != 0 {
		t.Fatalf("a failed read recorded %+v (%v)", got, err)
	}
}

func githubActionsTestEnvRecording(home string, err error, got *[]recorded) githubActionsEnv {
	env := githubActionsTestEnv(home, err)
	env.record = func(series string, p quotahistory.Point) error {
		*got = append(*got, recorded{series, p})
		return nil
	}
	return env
}

func runQuota(t *testing.T, now time.Time, args ...string) (string, error) {
	t.Helper()
	command := newQuotaCommand(quotaEnv{now: func() time.Time { return now }})
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func TestQuotaHistoryReadsTheStateDirectory(t *testing.T) {
	state := t.TempDir()
	t.Setenv("XDG_STATE_HOME", state)
	store := quotahistory.Open(state + "/ply/quota-history")
	now := time.Date(2026, 11, 28, 12, 0, 0, 0, time.UTC)
	for _, p := range []struct {
		series string
		at     time.Time
		used   float64
	}{{"claude.weekly", now.Add(-time.Hour), 40}, {"claude.weekly", now.Add(-9 * 24 * time.Hour), 10}, {"github.actions", now.Add(-2 * time.Hour), 100}} {
		if _, err := store.Record(p.series, quotahistory.Point{At: p.at, Used: p.used}); err != nil {
			t.Fatal(err)
		}
	}
	out, err := runQuota(t, now, "history", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var h struct {
		Kind          string                    `json:"kind"`
		SchemaVersion int                       `json:"schema_version"`
		Since         string                    `json:"since"`
		Until         string                    `json:"until"`
		Series        []quotahistory.SeriesView `json:"series"`
		Invalid       []quotahistory.Invalid    `json:"invalid"`
	}
	if err := json.Unmarshal([]byte(out), &h); err != nil || h.Kind != "PlyQuotaHistory@1" || h.SchemaVersion != 1 || h.Since != "2026-11-20T12:00:00Z" ||
		len(h.Series) != 5 || len(h.Series[1].Points) != 1 || h.Series[1].Points[0].Used != 40 || len(h.Series[4].Points) != 1 || h.Invalid == nil {
		t.Fatalf("%s %v", out, err)
	}
	out, err = runQuota(t, now, "history", "--series", "github.actions", "--since", "2026-11-28T11:00:00Z", "--format", "json")
	if err != nil || strings.Count(out, `"id"`) != 1 || !strings.Contains(out, `"since": "2026-11-28T11:00:00Z"`) || strings.Contains(out, `"used": 100`) {
		t.Fatalf("filtered %s %v", out, err)
	}
	if _, err := runQuota(t, now, "history", "--series", "claude.daily"); err == nil || !strings.HasPrefix(err.Error(), "quota_invalid_arguments: ") {
		t.Fatalf("unknown series: %v", err)
	}
	if _, err := runQuota(t, now, "history", "--since", "yesterday"); err == nil || !strings.HasPrefix(err.Error(), "quota_invalid_arguments: ") {
		t.Fatalf("bad since: %v", err)
	}
	if text, err := runQuota(t, now, "history"); err != nil || !strings.Contains(text, "claude.weekly") {
		t.Fatalf("text %q %v", text, err)
	}
}
