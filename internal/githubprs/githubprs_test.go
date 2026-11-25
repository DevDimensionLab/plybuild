package githubprs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/ghapi"
	"github.com/devdimensionlab/plybuild/internal/testutil"
)

func TestMain(m *testing.M) { os.Exit(testutil.Isolate(m)) }

type reply struct {
	stdout string
	code   int
	err    error
}

// fakeGH answers each search by its qualifier and the GraphQL call by "graphql".
type fakeGH struct {
	replies   map[string]reply
	deadlines []time.Duration
	graphql   string
}

func (f *fakeGH) Run(c process.Command) error {
	if c.Name != "gh" || c.Context == nil {
		return fmt.Errorf("unexpected command %s %q", c.Name, c.Args)
	}
	deadline, _ := c.Context.Deadline()
	f.deadlines = append(f.deadlines, time.Until(deadline).Round(time.Second))
	key := ""
	for i, a := range c.Args {
		if a == "graphql" {
			key = "graphql"
			f.graphql = strings.Join(c.Args[i+1:], " ")
		}
		if strings.HasPrefix(a, "q=") {
			for _, q := range []string{"user-review-requested:@me", "review-requested:@me", "assignee:@me", "author:@me"} {
				if strings.Contains(a, q) {
					key = q
					break
				}
			}
		}
	}
	r, ok := f.replies[key]
	if !ok {
		return fmt.Errorf("no reply for %q", c.Args)
	}
	fmt.Fprint(c.Stdout, r.stdout)
	if r.err != nil {
		return r.err
	}
	if r.code != 0 {
		return exec.Command("/bin/sh", "-c", fmt.Sprintf("exit %d", r.code)).Run()
	}
	return nil
}

func item(id, repo string, number int, author string) string {
	return fmt.Sprintf(`{"node_id":%q,"number":%d,"title":"T%d","html_url":"https://github.com/%s/pull/%d","user":{"login":%q},"draft":false,`+
		`"created_at":"2026-11-0%dT00:00:00Z","updated_at":"2026-11-2%dT00:00:00Z","repository_url":"https://api.github.com/repos/%s"}`,
		id, number, number, repo, number, author, number, number, repo)
}

func search(total int, incomplete bool, items ...string) reply {
	return reply{stdout: fmt.Sprintf(`{"total_count":%d,"incomplete_results":%t,"items":[%s]}`, total, incomplete, strings.Join(items, ","))}
}

func node(id, decision, mergeable, state, checks string, unresolved, resolved int) string {
	threads := []string{}
	for i := 0; i < unresolved; i++ {
		threads = append(threads, `{"isResolved":false}`)
	}
	for i := 0; i < resolved; i++ {
		threads = append(threads, `{"isResolved":true}`)
	}
	rollup := "null"
	if checks != "" {
		rollup = fmt.Sprintf(`{"state":%q}`, checks)
	}
	dec := "null"
	if decision != "" {
		dec = fmt.Sprintf("%q", decision)
	}
	return fmt.Sprintf(`{"id":%q,"reviewDecision":%s,"mergeable":%q,"mergeStateStatus":%q,"commits":{"nodes":[{"commit":{"statusCheckRollup":%s}}]},"reviewThreads":{"nodes":[%s]}}`,
		id, dec, mergeable, state, rollup, strings.Join(threads, ","))
}

func str(s string) *string { return &s }
func num(n int) *int       { return &n }

// The four searches give the reasons; GraphQL gives the details, except for a
// pull request GitHub refuses (SAML), which keeps its search data.
func TestReadMergesTheSearchesAndTheDetails(t *testing.T) {
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me":      search(2, false, item("PR_team", "org/app", 1, "ann"), item("PR_direct", "me/lib", 2, "bob")),
		"user-review-requested:@me": search(1, false, item("PR_direct", "me/lib", 2, "bob")),
		"assignee:@me":              search(1, false, item("PR_own", "me/tool", 3, "me")),
		"author:@me":                search(1, false, item("PR_own", "me/tool", 3, "me")),
		"graphql": {code: 1, stdout: `{"data":{"viewer":{"login":"me"},"nodes":[null,` +
			node("PR_direct", "REVIEW_REQUIRED", "MERGEABLE", "BLOCKED", "SUCCESS", 0, 0) + `,` +
			node("PR_own", "CHANGES_REQUESTED", "CONFLICTING", "DIRTY", "FAILURE", 2, 1) + `]},` +
			`"errors":[{"type":"FORBIDDEN","path":["nodes",0],"extensions":{"saml_failure":true}}]}`},
	}}
	got, err := Read(context.Background(), gh)
	if err != nil {
		t.Fatal(err)
	}
	if got.Viewer != "me" || len(got.IncompleteReasons) != 0 {
		t.Fatalf("viewer %q incomplete %v", got.Viewer, got.IncompleteReasons)
	}
	if !strings.Contains(gh.graphql, `"PR_team","PR_direct","PR_own"`) {
		t.Fatalf("graphql asked %s", gh.graphql)
	}
	want := []PullRequest{
		{Repository: "org/app", Number: 1, Title: "T1", URL: "https://github.com/org/app/pull/1", Author: "ann", CreatedAt: "2026-11-01T00:00:00Z", UpdatedAt: "2026-11-21T00:00:00Z",
			Reasons: []Reason{{Kind: "review_requested", Via: "team"}}, Detail: "unavailable", NeedsYou: []string{"review_requested"}},
		{Repository: "me/lib", Number: 2, Title: "T2", URL: "https://github.com/me/lib/pull/2", Author: "bob", CreatedAt: "2026-11-02T00:00:00Z", UpdatedAt: "2026-11-22T00:00:00Z",
			Reasons: []Reason{{Kind: "review_requested", Via: "direct"}}, Detail: "available",
			ReviewDecision: str("REVIEW_REQUIRED"), Checks: str("SUCCESS"), Mergeable: str("MERGEABLE"), MergeState: str("BLOCKED"), NeedsYou: []string{"review_requested"}},
		{Repository: "me/tool", Number: 3, Title: "T3", URL: "https://github.com/me/tool/pull/3", Author: "me", CreatedAt: "2026-11-03T00:00:00Z", UpdatedAt: "2026-11-23T00:00:00Z",
			Reasons: []Reason{{Kind: "assigned"}, {Kind: "author"}}, Detail: "available",
			ReviewDecision: str("CHANGES_REQUESTED"), Checks: str("FAILURE"), Mergeable: str("CONFLICTING"), MergeState: str("DIRTY"), UnresolvedThreads: num(2),
			NeedsYou: []string{"assigned", "changes_requested", "checks_failed", "conflict", "unresolved_threads"}},
	}
	if !reflect.DeepEqual(got.PullRequests, want) {
		t.Fatalf("got  %+v\nwant %+v", got.PullRequests, want)
	}
	for i, d := range gh.deadlines {
		if d > ghapi.CallTimeout {
			t.Fatalf("call %d deadline %v", i, d)
		}
	}
}

// ready_to_merge is GitHub's own CLEAN; checks ERROR also fails; null review
// decision and checks are kept as null.
func TestNeedsYouOfOwnPullRequests(t *testing.T) {
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(0, false), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false),
		"author:@me": search(2, false, item("PR_a", "me/a", 1, "me"), item("PR_b", "me/b", 2, "me")),
		"graphql": {stdout: `{"data":{"viewer":{"login":"me"},"nodes":[` +
			node("PR_a", "", "MERGEABLE", "CLEAN", "", 0, 3) + `,` + node("PR_b", "APPROVED", "MERGEABLE", "UNSTABLE", "ERROR", 0, 0) + `]}}`},
	}}
	got, err := Read(context.Background(), gh)
	if err != nil {
		t.Fatal(err)
	}
	a, b := got.PullRequests[0], got.PullRequests[1]
	if a.ReviewDecision != nil || a.Checks != nil || *a.UnresolvedThreads != 0 || !reflect.DeepEqual(a.NeedsYou, []string{"ready_to_merge"}) {
		t.Fatalf("a %+v", a)
	}
	if !reflect.DeepEqual(b.NeedsYou, []string{"checks_failed"}) {
		t.Fatalf("b %+v", b)
	}
}

func TestIncompleteSearchesAreReported(t *testing.T) {
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(150, false, item("PR_x", "o/r", 1, "a")), "user-review-requested:@me": search(0, true), "assignee:@me": search(0, false), "author:@me": search(0, false),
		"graphql": {stdout: `{"data":{"viewer":{"login":"me"},"nodes":[` + node("PR_x", "", "UNKNOWN", "UNKNOWN", "", 0, 0) + `]}}`},
	}}
	got, err := Read(context.Background(), gh)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.IncompleteReasons, []string{"search_truncated", "search_incomplete"}) {
		t.Fatalf("incomplete %v", got.IncompleteReasons)
	}
}

// No pull requests still reads the viewer; an empty list is a complete answer.
func TestNothingOnTheUser(t *testing.T) {
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(0, false), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false), "author:@me": search(0, false),
		"graphql": {stdout: `{"data":{"viewer":{"login":"me"},"nodes":[]}}`},
	}}
	got, err := Read(context.Background(), gh)
	if err != nil || got.Viewer != "me" || got.PullRequests == nil || len(got.PullRequests) != 0 || got.IncompleteReasons == nil {
		t.Fatalf("%+v %v", got, err)
	}
}

// A failed search or GraphQL call fails the whole read with a ply github code.
func TestFailuresFailTheRead(t *testing.T) {
	ok := map[string]reply{
		"review-requested:@me": search(0, false), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false), "author:@me": search(0, false),
		"graphql": {stdout: `{"data":{"viewer":{"login":"me"},"nodes":[]}}`},
	}
	for name, c := range map[string]struct {
		key   string
		reply reply
		want  error
	}{
		"gh missing":           {"review-requested:@me", reply{err: &exec.Error{Name: "gh", Err: exec.ErrNotFound}}, ghapi.ErrMissing},
		"not logged in":        {"author:@me", reply{code: 4}, ghapi.ErrAuth},
		"search refused":       {"assignee:@me", reply{code: 1, stdout: `{"message":"API rate limit exceeded","status":"403"}`}, ghapi.ErrFailed},
		"graphql without data": {"graphql", reply{code: 1, stdout: `{"errors":[{"message":"boom"}]}`}, ghapi.ErrFailed},
		"graphql unreadable":   {"graphql", reply{stdout: `not json`}, ghapi.ErrFailed},
		"graphql refused":      {"graphql", reply{code: 1, stdout: `{"message":"Forbidden","status":"403"}`}, ghapi.ErrFailed},
		"graphql timeout":      {"graphql", reply{err: context.DeadlineExceeded}, ghapi.ErrTimeout},
	} {
		replies := map[string]reply{}
		for k, v := range ok {
			replies[k] = v
		}
		replies[c.key] = c.reply
		if _, err := Read(context.Background(), &fakeGH{replies: replies}); !errors.Is(err, c.want) {
			t.Fatalf("%s: %v, want %v", name, err, c.want)
		}
	}
}

// All calls share one deadline of at most 40 seconds.
func TestReadHasOneDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(0, false), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false), "author:@me": search(0, false),
		"graphql": {stdout: `{"data":{"viewer":{"login":"me"},"nodes":[]}}`},
	}}
	if _, err := Read(ctx, gh); err != nil {
		t.Fatal(err)
	}
	if Deadline != 40*time.Second || len(gh.deadlines) != 5 {
		t.Fatalf("deadline %v calls %d", Deadline, len(gh.deadlines))
	}
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	gh.deadlines = nil
	if _, err := Read(ctx, gh); err != nil {
		t.Fatal(err)
	}
	for _, d := range gh.deadlines {
		if d > 5*time.Second {
			t.Fatalf("call deadline %v beyond the caller's", d)
		}
	}
}
