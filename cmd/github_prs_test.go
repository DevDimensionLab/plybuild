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
