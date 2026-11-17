package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/githubbilling"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

var githubActionsNow = time.Date(2026, time.November, 17, 20, 10, 0, 0, time.UTC)

func githubActionsItems() []githubbilling.Item {
	minutes := func(repo, sku string, q float64) githubbilling.Item {
		return githubbilling.Item{Product: "actions", SKU: sku, UnitType: "Minutes", Repository: repo, Quantity: q}
	}
	return []githubbilling.Item{
		minutes("App", "Actions Linux", 1000),
		minutes("App", "Actions Linux", 281),
		minutes("tui", "Actions macOS 3-core", 2),
		minutes("tui", "Actions Linux", 1),
		minutes("win", "Actions Windows", 3),
		minutes("gpu", "Actions GPU", 50),
		{Product: "actions", SKU: "Actions storage", UnitType: "GigabyteHours", Repository: "App", Quantity: 9},
	}
}

// githubActionsHome writes the included minutes, or no file when minutes is 0.
func githubActionsHome(t *testing.T, minutes int) string {
	t.Helper()
	home := t.TempDir()
	if minutes > 0 {
		dir := filepath.Join(home, ".config", "ply")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := fmt.Sprintf(`{"actions": {"included_minutes": %d}}`, minutes)
		if err := os.WriteFile(filepath.Join(dir, "github.json"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func executeGitHub(t *testing.T, env githubActionsEnv, args ...string) (string, error) {
	t.Helper()
	command := newGitHubCommand(env)
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	err := command.ExecuteContext(context.Background())
	return stdout.String(), err
}

func githubActionsTestEnv(home string, err error) githubActionsEnv {
	return githubActionsEnv{
		read: func(_ context.Context, year int, month time.Month) (githubbilling.Usage, error) {
			if year != 2026 || month != time.November {
				return githubbilling.Usage{}, fmt.Errorf("asked for %d-%d", year, month)
			}
			return githubbilling.Usage{Login: "octo", Items: githubActionsItems()}, err
		},
		now:  func() time.Time { return githubActionsNow },
		home: func() (string, error) { return home, nil },
	}
}

func TestGitHubActionsQuotaCountsIncludedMinutesAgainstTheLimit(t *testing.T) {
	stdout, err := executeGitHub(t, githubActionsTestEnv(githubActionsHome(t, 3000), nil), "actions", "quota", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	// 1281 + 1 Linux + 2×10 macOS + 3×2 Windows = 1308; GPU and storage are not counted.
	want := `{"kind":"PlyGitHubActionsQuota@1","schema_version":1,"id":"github.actions","source":"github","status":"ok",` +
		`"observed_at":"2026-11-17T20:10:00Z","unit":"minutes","used":1308,"limit":3000,"remaining":1692,` +
		`"remaining_percent":56,"resets_at":"2026-12-01T00:00:00Z","level":"ok"}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGitHubActionsQuotaWithoutConfigOmitsLimitAndLevel(t *testing.T) {
	stdout, err := executeGitHub(t, githubActionsTestEnv(githubActionsHome(t, 0), nil), "actions", "quota", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"PlyGitHubActionsQuota@1","schema_version":1,"id":"github.actions","source":"github","status":"ok",` +
		`"observed_at":"2026-11-17T20:10:00Z","unit":"minutes","used":1308,"resets_at":"2026-12-01T00:00:00Z"}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGitHubActionsLevel(t *testing.T) {
	for _, test := range []struct {
		remaining, percent int
		want               string
	}{
		{0, 0, "exhausted"}, {1, 0, "low"}, {570, 19, "low"}, {600, 20, "ok"}, {3000, 100, "ok"},
	} {
		if got := githubActionsLevel(test.remaining, test.percent); got != test.want {
			t.Fatalf("level(%d, %d%%) = %q, want %q", test.remaining, test.percent, got, test.want)
		}
	}
}

func TestGitHubActionsUsageWithoutMinutesIsAnEmptyList(t *testing.T) {
	t.Chdir(t.TempDir())
	env := githubActionsTestEnv(githubActionsHome(t, 3000), nil)
	env.read = func(context.Context, int, time.Month) (githubbilling.Usage, error) {
		return githubbilling.Usage{Login: "octo", Items: []githubbilling.Item{}}, nil
	}
	stdout, err := executeGitHub(t, env, "actions", "usage", "--format", "json")
	if err != nil || !strings.HasSuffix(stdout, `"unit":"minutes","repos":[]}`+"\n") {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
}

func TestGitHubActionsErrorIsReportedWithItsCodeOnly(t *testing.T) {
	failure := fmt.Errorf("%w: billing usage needs the gh scope user", githubbilling.ErrScope)
	for _, command := range []string{"quota", "usage"} {
		stdout, err := executeGitHub(t, githubActionsTestEnv(githubActionsHome(t, 3000), failure), "actions", command, "--format", "json")
		if err == nil {
			t.Fatalf("%s: error status exited 0", command)
		}
		if !strings.Contains(stdout, `"status":"error","error":{"code":"gh_scope","description":"billing usage needs the gh scope user"}`) || strings.Contains(stdout, "used") {
			t.Fatalf("%s: stdout %s", command, stdout)
		}
	}
}

func TestGitHubActionsUsageListsReposWithTheirProjects(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if _, err := workspace.Init(workspace.SystemDependencies()); err != nil {
		t.Fatal(err)
	}
	definitions := `{"format_version": 1, "projects": [
  {"id": "a", "name": "A", "wrapper": "a", "repos": [{"id": "app", "path": "a/app", "remote": "git@GitHub.com:Octo/App.git"}]},
  {"id": "b", "name": "B", "wrapper": "b", "repos": [{"id": "app", "path": "b/app", "remote": "https://github.com/octo/app"},
    {"id": "other", "path": "b/other", "remote": "git@github.com:someone/tui.git"}]}
]}
`
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, err := executeGitHub(t, githubActionsTestEnv(githubActionsHome(t, 3000), nil), "actions", "usage", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"PlyGitHubActionsUsage@1","schema_version":1,"source":"github","status":"ok",` +
		`"observed_at":"2026-11-17T20:10:00Z","period_start":"2026-11-01T00:00:00Z","unit":"minutes","repos":[` +
		`{"key":"octo/App","project_ids":["a","b"],"minutes":1281,"macos_minutes":0},` +
		`{"key":"octo/tui","project_ids":[],"minutes":21,"macos_minutes":20},` +
		`{"key":"octo/win","project_ids":[],"minutes":6,"macos_minutes":0}]}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGitHubActionsUsageOutsideAWorkspaceOmitsProjects(t *testing.T) {
	t.Chdir(t.TempDir())
	stdout, err := executeGitHub(t, githubActionsTestEnv(githubActionsHome(t, 3000), nil), "actions", "usage", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout, "project_ids") || !strings.Contains(stdout, `{"key":"octo/tui","minutes":21,"macos_minutes":20}`) {
		t.Fatalf("stdout %s", stdout)
	}
}

func TestGitHubActionsText(t *testing.T) {
	t.Chdir(t.TempDir())
	env := githubActionsTestEnv(githubActionsHome(t, 3000), nil)
	quota, err := executeGitHub(t, env, "actions", "quota")
	if err != nil || quota != "GitHub Actions: 1308 of 3000 minutes used, 1692 left (56%, ok), resets 2026-12-01\n" {
		t.Fatalf("quota %q, error %v", quota, err)
	}
	usage, err := executeGitHub(t, env, "actions", "usage")
	if err != nil || usage != "octo/App  1281\nocto/tui  21 (macOS 20)\nocto/win  6\n" {
		t.Fatalf("usage %q, error %v", usage, err)
	}
}
