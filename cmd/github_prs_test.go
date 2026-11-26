package cmd

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/ghapi"
	"github.com/devdimensionlab/plybuild/internal/githubprs"
)

func githubPRsTestEnv(result githubprs.Result, err error) githubActionsEnv {
	env := githubActionsTestEnv("", nil)
	env.prs = func(context.Context) (githubprs.Result, error) { return result, err }
	return env
}

func githubPRsTestResult() githubprs.Result {
	state := "CLEAN"
	return githubprs.Result{Viewer: "me", IncompleteReasons: []string{}, PullRequests: []githubprs.PullRequest{
		{Repository: "org/app", Number: 1, Title: "Team review", URL: "https://github.com/org/app/pull/1", Author: "ann", CreatedAt: "c", UpdatedAt: "u",
			Reasons: []githubprs.Reason{{Kind: "review_requested", Via: "team"}}, Detail: "unavailable", NeedsYou: []string{"review_requested"}},
		{Repository: "me/tool", Number: 3, Title: "Own work", URL: "https://github.com/me/tool/pull/3", Author: "me", Draft: true, CreatedAt: "c", UpdatedAt: "u",
			Reasons: []githubprs.Reason{{Kind: "author"}}, Detail: "available", MergeState: &state, NeedsYou: []string{}},
	}}
}

func TestGitHubPRsWritesTheContract(t *testing.T) {
	stdout, err := executeGitHub(t, githubPRsTestEnv(githubPRsTestResult(), nil), "prs", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"PlyGitHubPullRequests@1","schema_version":1,"status":"ok","observed_at":"2026-11-17T20:10:00Z","viewer":"me","incomplete_reasons":[],"pull_requests":[` +
		`{"repository":"org/app","number":1,"title":"Team review","url":"https://github.com/org/app/pull/1","author":"ann","draft":false,"created_at":"c","updated_at":"u",` +
		`"reasons":[{"kind":"review_requested","via":"team"}],"detail":"unavailable","review_decision":null,"checks":null,"mergeable":null,"merge_state":null,"unresolved_threads":null,"needs_you":["review_requested"]},` +
		`{"repository":"me/tool","number":3,"title":"Own work","url":"https://github.com/me/tool/pull/3","author":"me","draft":true,"created_at":"c","updated_at":"u",` +
		`"reasons":[{"kind":"author"}],"detail":"available","review_decision":null,"checks":null,"mergeable":null,"merge_state":"CLEAN","unresolved_threads":null,"needs_you":[]}]}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGitHubPRsTextGroupsWhatWaitsForYou(t *testing.T) {
	stdout, err := executeGitHub(t, githubPRsTestEnv(githubPRsTestResult(), nil), "prs")
	if err != nil {
		t.Fatal(err)
	}
	waiting, own := strings.Index(stdout, "Waiting for you"), strings.Index(stdout, "Your pull requests")
	team, work := strings.Index(stdout, "org/app#1 Team review"), strings.Index(stdout, "me/tool#3 Own work")
	if waiting < 0 || own < waiting || team < waiting || team > own || work < own ||
		!strings.Contains(stdout, "review_requested (team); no details from GitHub") || !strings.Contains(stdout, "draft") {
		t.Fatalf("stdout\n%s", stdout)
	}
}

// A failure is the same document with status error and exit status 1, as in ply github actions.
func TestGitHubPRsReportsAFailureOnStdout(t *testing.T) {
	failure := fmt.Errorf("%w: gh is not logged in", ghapi.ErrAuth)
	stdout, err := executeGitHub(t, githubPRsTestEnv(githubprs.Result{}, failure), "prs", "--format", "json")
	if err == nil {
		t.Fatal("no error")
	}
	want := `{"kind":"PlyGitHubPullRequests@1","schema_version":1,"status":"error","error":{"code":"gh_auth","description":"gh is not logged in"}}` + "\n"
	// The test command has no root that silences cobra's usage text after the document.
	if !strings.HasPrefix(stdout, want) {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGitHubPRsListWritesTheOverview(t *testing.T) {
	env := githubActionsTestEnv("", nil)
	env.list = func(context.Context) (githubprs.ListResult, error) {
		return githubprs.ListResult{IncompleteReasons: []string{}, PullRequests: []githubprs.Listed{
			{ID: "PR_1", Repository: "org/app", Number: 1, Title: "Team review", URL: "https://github.com/org/app/pull/1", Author: "ann", CreatedAt: "c", UpdatedAt: "u",
				Reasons: []githubprs.Reason{{Kind: "review_requested", Via: "team"}}, NeedsYou: []string{"review_requested"}},
		}}, nil
	}
	stdout, err := executeGitHub(t, env, "prs", "list", "--format", "json")
	want := `{"kind":"PlyGitHubPullRequestList@1","schema_version":1,"status":"ok","observed_at":"2026-11-17T20:10:00Z","incomplete_reasons":[],"pull_requests":[` +
		`{"id":"PR_1","repository":"org/app","number":1,"title":"Team review","url":"https://github.com/org/app/pull/1","author":"ann","draft":false,"created_at":"c","updated_at":"u",` +
		`"reasons":[{"kind":"review_requested","via":"team"}],"needs_you":["review_requested"]}]}` + "\n"
	if err != nil || stdout != want {
		t.Fatalf("%v\n%s\nwant\n%s", err, stdout, want)
	}
	if stdout, err = executeGitHub(t, env, "prs", "list"); err != nil || !strings.Contains(stdout, "org/app#1 Team review\n    review_requested (team); needs you: review_requested\n    https://github.com/org/app/pull/1") {
		t.Fatalf("text %v\n%s", err, stdout)
	}
	env.list = func(context.Context) (githubprs.ListResult, error) {
		return githubprs.ListResult{}, fmt.Errorf("%w: gh is not logged in", ghapi.ErrAuth)
	}
	stdout, err = executeGitHub(t, env, "prs", "list", "--format", "json")
	if err == nil || !strings.HasPrefix(stdout, `{"kind":"PlyGitHubPullRequestList@1","schema_version":1,"status":"error","error":{"code":"gh_auth","description":"gh is not logged in"}}`+"\n") {
		t.Fatalf("failure %v\n%s", err, stdout)
	}
}

func TestGitHubPRsDetailsWritesEachPullRequest(t *testing.T) {
	env := githubActionsTestEnv("", nil)
	failed, state := githubprs.DetailFailed, "CLEAN"
	env.details = func(context.Context) ([]githubprs.Detailed, error) {
		return []githubprs.Detailed{
			{ID: "PR_1", Repository: "org/app", Number: 1, Detail: "unavailable", DetailReason: &failed, NeedsYou: []string{}},
			{ID: "PR_3", Repository: "me/tool", Number: 3, Detail: "available", MergeState: &state, NeedsYou: []string{"ready_to_merge"}},
		}, nil
	}
	stdout, err := executeGitHub(t, env, "prs", "details", "--format", "json")
	want := `{"kind":"PlyGitHubPullRequestDetails@1","schema_version":1,"status":"ok","observed_at":"2026-11-17T20:10:00Z","pull_requests":[` +
		`{"id":"PR_1","repository":"org/app","number":1,"detail":"unavailable","detail_reason":"gh_failed","review_decision":null,"checks":null,"mergeable":null,"merge_state":null,"unresolved_threads":null,"needs_you":[]},` +
		`{"id":"PR_3","repository":"me/tool","number":3,"detail":"available","detail_reason":null,"review_decision":null,"checks":null,"mergeable":null,"merge_state":"CLEAN","unresolved_threads":null,"needs_you":["ready_to_merge"]}]}` + "\n"
	if err != nil || stdout != want {
		t.Fatalf("%v\n%s\nwant\n%s", err, stdout, want)
	}
	if stdout, err = executeGitHub(t, env, "prs", "details"); err != nil ||
		!strings.Contains(stdout, "org/app#1 no details (gh_failed)") || !strings.Contains(stdout, "me/tool#3 merge state CLEAN; needs you: ready_to_merge") {
		t.Fatalf("text %v\n%s", err, stdout)
	}
}
