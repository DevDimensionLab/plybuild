package githubprs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
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

// fakeGH answers each search by its qualifier and the GraphQL call by
// "graphql", or, with nodes, each GraphQL call from the ids it asks for: an id
// without a node is refused as SAML refuses it.
type fakeGH struct {
	replies   map[string]reply
	nodes     map[string]string
	deadlines []time.Duration
	graphql   string
	graphqls  []string
}

var askedIDs = regexp.MustCompile(`ids:\[([^\]]*)\]`)

func (f *fakeGH) answer(query string) reply {
	m := askedIDs.FindStringSubmatch(query)
	nodes, errs := []string{}, []string{}
	for i, id := range strings.Split(m[1], ",") {
		if id == "" {
			continue
		}
		if n, ok := f.nodes[strings.Trim(id, `"`)]; ok {
			nodes = append(nodes, n)
		} else {
			nodes = append(nodes, "null")
			errs = append(errs, fmt.Sprintf(`{"type":"FORBIDDEN","path":["nodes",%d],"extensions":{"saml_failure":true}}`, i))
		}
	}
	r := reply{stdout: `{"data":{"viewer":{"login":"me"},"nodes":[` + strings.Join(nodes, ",") + `]}`}
	if len(errs) > 0 {
		r.code, r.stdout = 1, r.stdout+`,"errors":[`+strings.Join(errs, ",")+`]`
	}
	r.stdout += "}"
	return r
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
			f.graphqls = append(f.graphqls, f.graphql)
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
	if key == "graphql" && f.nodes != nil {
		r, ok = f.answer(f.graphql), true
	}
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
	}, nodes: map[string]string{
		"PR_direct": node("PR_direct", "REVIEW_REQUIRED", "MERGEABLE", "BLOCKED", "SUCCESS", 0, 0),
		"PR_own":    node("PR_own", "CHANGES_REQUESTED", "CONFLICTING", "DIRTY", "FAILURE", 2, 1),
	}}
	got, err := Read(context.Background(), gh)
	if err != nil {
		t.Fatal(err)
	}
	if got.Viewer != "me" || len(got.IncompleteReasons) != 0 {
		t.Fatalf("viewer %q incomplete %v", got.Viewer, got.IncompleteReasons)
	}
	if len(gh.graphqls) != 2 || !strings.Contains(gh.graphqls[0], `ids:["PR_team","PR_direct"]`) || !strings.Contains(gh.graphqls[1], `ids:["PR_own"]`) {
		t.Fatalf("graphql asked %q", gh.graphqls)
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

// Details come in batches of at most 20, and review threads, which only the
// user's own pull requests use, are asked for only in their batches: GitHub
// answered 502 when one query asked for the threads of dozens of pull requests.
func TestDetailsComeInBatches(t *testing.T) {
	var team, own []string
	nodes := map[string]string{}
	for n := 1; n <= 45; n++ {
		id := fmt.Sprintf("PR_t%d", n)
		team, nodes[id] = append(team, item(id, "org/app", n, "ann")), node(id, "", "MERGEABLE", "BLOCKED", "", 0, 0)
	}
	for n := 1; n <= 2; n++ {
		id := fmt.Sprintf("PR_o%d", n)
		own, nodes[id] = append(own, item(id, "me/tool", n, "me")), node(id, "", "MERGEABLE", "CLEAN", "", 1, 0)
	}
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(45, false, team...), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false),
		"author:@me": search(2, false, own...),
	}, nodes: nodes}
	got, err := Read(context.Background(), gh)
	if err != nil {
		t.Fatal(err)
	}
	var sizes []int
	for _, q := range gh.graphqls {
		ids := strings.Count(askedIDs.FindStringSubmatch(q)[1], `"`) / 2
		sizes = append(sizes, ids)
		if threads := strings.Contains(q, "reviewThreads"); threads != strings.Contains(q, "PR_o") {
			t.Fatalf("threads %v in %s", threads, q)
		}
	}
	if fmt.Sprint(sizes) != "[20 20 5 2]" {
		t.Fatalf("batches %v", sizes)
	}
	for _, pr := range got.PullRequests {
		if pr.Detail != "available" || (pr.Author == "me") != (pr.UnresolvedThreads != nil) {
			t.Fatalf("%+v", pr)
		}
	}
}

// A batch GitHub fails, as with 502 Bad Gateway, fails the read as before.
func TestAFailedBatchFailsTheRead(t *testing.T) {
	var team []string
	for n := 1; n <= 25; n++ {
		team = append(team, item(fmt.Sprintf("PR_t%d", n), "org/app", n, "ann"))
	}
	gh := &fakeGH{replies: map[string]reply{
		"review-requested:@me": search(25, false, team...), "user-review-requested:@me": search(0, false), "assignee:@me": search(0, false), "author:@me": search(0, false),
		"graphql": {code: 1, stdout: `{"message":"Bad Gateway","status":"502"}`},
	}}
	if _, err := Read(context.Background(), gh); !errors.Is(err, ghapi.ErrFailed) || len(gh.graphqls) != 1 {
		t.Fatalf("%v after %d calls", err, len(gh.graphqls))
	}
}
