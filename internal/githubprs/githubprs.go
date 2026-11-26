// Package githubprs reads the open pull requests that are on the GitHub user gh
// is logged in as: review requested from the user or a team, assigned to the
// user, or made by the user. Ply gets what gh gets; a pull request GitHub
// gives search data but no details for keeps its search data.
package githubprs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/ghapi"
)

// Deadline bounds a read with details: four searches and the details calls.
const Deadline = 40 * time.Second

// ListDeadline bounds the overview: the four searches only.
const ListDeadline = 20 * time.Second

// Reasons a pull request has no details in Details.
const (
	DetailRefused = "refused"   // GitHub answered without it, for example SAML; lasting
	DetailFailed  = "gh_failed" // its batch got no readable answer, for example 502
	DetailTimeout = "timeout"   // the deadline passed before its batch was read
	detailOK      = ""
)

// Reason says why a pull request is on the user.
type Reason struct {
	Kind string `json:"kind"`
	Via  string `json:"via,omitempty"`
}

// PullRequest is one open pull request; detail fields are GitHub's values
// unchanged and null when GitHub gave no details.
type PullRequest struct {
	Repository        string   `json:"repository"`
	Number            int      `json:"number"`
	Title             string   `json:"title"`
	URL               string   `json:"url"`
	Author            string   `json:"author"`
	Draft             bool     `json:"draft"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	Reasons           []Reason `json:"reasons"`
	Detail            string   `json:"detail"`
	ReviewDecision    *string  `json:"review_decision"`
	Checks            *string  `json:"checks"`
	Mergeable         *string  `json:"mergeable"`
	MergeState        *string  `json:"merge_state"`
	UnresolvedThreads *int     `json:"unresolved_threads"`
	NeedsYou          []string `json:"needs_you"`
}

// Listed is one pull request of the overview: what the searches give.
type Listed struct {
	ID         string   `json:"id"`
	Repository string   `json:"repository"`
	Number     int      `json:"number"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Author     string   `json:"author"`
	Draft      bool     `json:"draft"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Reasons    []Reason `json:"reasons"`
	NeedsYou   []string `json:"needs_you"`
}

// ListResult is the overview.
type ListResult struct {
	IncompleteReasons []string
	PullRequests      []Listed
}

// Detailed is the details of one pull request; NeedsYou has only the reasons
// the details give.
type Detailed struct {
	ID                string   `json:"id"`
	Repository        string   `json:"repository"`
	Number            int      `json:"number"`
	Detail            string   `json:"detail"`
	DetailReason      *string  `json:"detail_reason"`
	ReviewDecision    *string  `json:"review_decision"`
	Checks            *string  `json:"checks"`
	Mergeable         *string  `json:"mergeable"`
	MergeState        *string  `json:"merge_state"`
	UnresolvedThreads *int     `json:"unresolved_threads"`
	NeedsYou          []string `json:"needs_you"`
}

// Result is the user and the pull requests on the user.
type Result struct {
	Viewer            string
	IncompleteReasons []string
	PullRequests      []PullRequest
}

// searches are the qualifiers in the order their reasons are listed.
var searches = []string{"review-requested:@me", "user-review-requested:@me", "assignee:@me", "author:@me"}

var nodeID = regexp.MustCompile(`^[A-Za-z0-9_=-]+$`)

type searchItem struct {
	NodeID        string `json:"node_id"`
	Number        int    `json:"number"`
	Title         string `json:"title"`
	HTMLURL       string `json:"html_url"`
	Draft         bool   `json:"draft"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	RepositoryURL string `json:"repository_url"`
	User          struct {
		Login string `json:"login"`
	} `json:"user"`
}

// found is the pull requests of the four searches, and their node ids in the same order.
type found struct {
	incomplete []string
	prs        []PullRequest
	ids        []string
}

// searchAll runs the four searches and merges their hits with their reasons.
func searchAll(ctx context.Context, runner process.Runner) (found, error) {
	f := found{incomplete: []string{}, prs: []PullRequest{}}
	hits := map[string][]searchItem{}
	for _, q := range searches {
		items, incomplete, err := searchPRs(ctx, runner, q)
		if err != nil {
			return found{}, err
		}
		hits[q] = items
		for _, reason := range incomplete {
			f.incomplete = appendOnce(f.incomplete, reason)
		}
	}
	direct := map[string]bool{}
	for _, it := range hits["user-review-requested:@me"] {
		direct[it.NodeID] = true
	}
	byID := map[string]int{}
	for _, q := range searches {
		for _, it := range hits[q] {
			i, seen := byID[it.NodeID]
			if !seen {
				i = len(f.prs)
				byID[it.NodeID] = i
				f.ids = append(f.ids, it.NodeID)
				f.prs = append(f.prs, it.pullRequest())
			}
			if r, ok := reasonOf(q, direct[it.NodeID]); ok && !hasReason(f.prs[i], r.Kind) {
				f.prs[i].Reasons = append(f.prs[i].Reasons, r)
			}
		}
	}
	return f, nil
}

// Read runs the searches and the details calls within Deadline, for
// PlyGitHubPullRequests@1: a batch without details fails the whole read.
// Errors carry one of the ghapi codes.
func Read(ctx context.Context, runner process.Runner) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, Deadline)
	defer cancel()
	f, err := searchAll(ctx, runner)
	if err != nil {
		return Result{}, err
	}
	viewer, _, err := details(ctx, runner, f.ids, f.prs, true)
	if err != nil {
		return Result{}, err
	}
	for i := range f.prs {
		f.prs[i].NeedsYou = append(searchNeeds(f.prs[i]), detailNeeds(f.prs[i])...)
	}
	return Result{Viewer: viewer, IncompleteReasons: f.incomplete, PullRequests: f.prs}, nil
}

// List runs the four searches within ListDeadline.
func List(ctx context.Context, runner process.Runner) (ListResult, error) {
	ctx, cancel := context.WithTimeout(ctx, ListDeadline)
	defer cancel()
	f, err := searchAll(ctx, runner)
	if err != nil {
		return ListResult{}, err
	}
	list := make([]Listed, len(f.prs))
	for i, pr := range f.prs {
		list[i] = Listed{ID: f.ids[i], Repository: pr.Repository, Number: pr.Number, Title: pr.Title, URL: pr.URL, Author: pr.Author,
			Draft: pr.Draft, CreatedAt: pr.CreatedAt, UpdatedAt: pr.UpdatedAt, Reasons: pr.Reasons, NeedsYou: searchNeeds(pr)}
	}
	return ListResult{IncompleteReasons: f.incomplete, PullRequests: list}, nil
}

// Details runs the searches and reads the details of every hit within
// Deadline. A batch without details leaves its pull requests unavailable
// with a reason; only gh missing, gh not logged in or a failed search fail it.
func Details(ctx context.Context, runner process.Runner) ([]Detailed, error) {
	ctx, cancel := context.WithTimeout(ctx, Deadline)
	defer cancel()
	f, err := searchAll(ctx, runner)
	if err != nil {
		return nil, err
	}
	_, reasons, err := details(ctx, runner, f.ids, f.prs, false)
	if err != nil {
		return nil, err
	}
	rows := make([]Detailed, len(f.prs))
	for i, pr := range f.prs {
		rows[i] = Detailed{ID: f.ids[i], Repository: pr.Repository, Number: pr.Number, Detail: pr.Detail,
			ReviewDecision: pr.ReviewDecision, Checks: pr.Checks, Mergeable: pr.Mergeable, MergeState: pr.MergeState,
			UnresolvedThreads: pr.UnresolvedThreads, NeedsYou: detailNeeds(pr)}
		if reasons[i] != detailOK {
			rows[i].DetailReason = &reasons[i]
		}
	}
	return rows, nil
}

// searchPRs runs one search and returns its pull requests and what it left out.
func searchPRs(ctx context.Context, runner process.Runner, q string) ([]searchItem, []string, error) {
	var body struct {
		TotalCount        int          `json:"total_count"`
		IncompleteResults bool         `json:"incomplete_results"`
		Items             []searchItem `json:"items"`
	}
	err := ghapi.Get(ctx, runner, "search/issues", []string{"q=is:pr is:open archived:false " + q, "per_page=100"}, &body)
	if errors.Is(err, ghapi.ErrNotPermitted) {
		return nil, nil, fmt.Errorf("%w: GitHub refused the search %s", ghapi.ErrFailed, q)
	}
	if err != nil {
		return nil, nil, err
	}
	for _, it := range body.Items {
		if !nodeID.MatchString(it.NodeID) {
			return nil, nil, fmt.Errorf("%w: unreadable search result", ghapi.ErrFailed)
		}
	}
	incomplete := []string{}
	if body.TotalCount > len(body.Items) {
		incomplete = append(incomplete, "search_truncated")
	}
	if body.IncompleteResults {
		incomplete = append(incomplete, "search_incomplete")
	}
	return body.Items, incomplete, nil
}

// reasonOf is the reason a hit of search q gives; the direct search only
// decides via, through the review-requested search that also finds it.
func reasonOf(q string, direct bool) (Reason, bool) {
	switch q {
	case "review-requested:@me":
		if direct {
			return Reason{Kind: "review_requested", Via: "direct"}, true
		}
		return Reason{Kind: "review_requested", Via: "team"}, true
	case "user-review-requested:@me":
		return Reason{Kind: "review_requested", Via: "direct"}, true
	case "assignee:@me":
		return Reason{Kind: "assigned"}, true
	case "author:@me":
		return Reason{Kind: "author"}, true
	}
	return Reason{}, false
}

func (it searchItem) pullRequest() PullRequest {
	return PullRequest{Repository: strings.TrimPrefix(it.RepositoryURL, "https://api.github.com/repos/"),
		Number: it.Number, Title: it.Title, URL: it.HTMLURL, Author: it.User.Login, Draft: it.Draft, CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt,
		Reasons: []Reason{}, Detail: "unavailable"}
}

// batchSize bounds the pull requests of one details call; GitHub answered 502
// when one call asked for dozens of pull requests with their review threads.
const batchSize = 20

const detailsQuery = `query{viewer{login} nodes(ids:[%s]){... on PullRequest{id reviewDecision mergeable mergeStateStatus ` +
	`commits(last:1){nodes{commit{statusCheckRollup{state}}}}%s}}}`

// reviewThreads is asked for only for the user's own pull requests, the only ones that use it.
const reviewThreads = ` reviewThreads(first:100){nodes{isResolved}}`

// details fills the detail fields of prs, whose node ids are ids, in the
// batches of batchesOf. It returns the viewer and, per pull request, why it
// has no details. gh missing or not
// logged in fails it; with stop, so does any batch without details.
func details(ctx context.Context, runner process.Runner, ids []string, prs []PullRequest, stop bool) (string, []string, error) {
	reasons := make([]string, len(ids))
	viewer := ""
	for _, b := range batchesOf(prs) {
		v, err := b.readInto(ctx, runner, ids, prs, reasons)
		if err != nil && (stop || errors.Is(err, ghapi.ErrMissing) || errors.Is(err, ghapi.ErrAuth)) {
			return "", nil, err
		}
		if err == nil {
			viewer = v
		}
	}
	return viewer, reasons, nil
}

// batchesOf is the batches for prs: first the others, then the user's own
// with review threads; one empty batch without pull requests, for the viewer.
func batchesOf(prs []PullRequest) []batch {
	var others, own []int
	for i := range prs {
		if hasReason(prs[i], "author") {
			own = append(own, i)
		} else {
			others = append(others, i)
		}
	}
	batches := append(batched(others, ""), batched(own, reviewThreads)...)
	if len(batches) == 0 {
		return []batch{{}}
	}
	return batches
}

// readInto reads b, unless the deadline has passed, and records in reasons
// why each of its pull requests has no details.
func (b batch) readInto(ctx context.Context, runner process.Runner, ids []string, prs []PullRequest, reasons []string) (string, error) {
	if ctx.Err() != nil {
		b.mark(reasons, DetailTimeout)
		return "", fmt.Errorf("%w: the details did not come within %s", ghapi.ErrTimeout, Deadline)
	}
	viewer, err := b.read(ctx, runner, ids, prs)
	if err != nil {
		b.mark(reasons, batchReason(err))
		return "", err
	}
	for _, at := range b.indexes {
		if prs[at].Detail != "available" {
			reasons[at] = DetailRefused
		}
	}
	return viewer, nil
}

func batchReason(err error) string {
	if errors.Is(err, ghapi.ErrTimeout) {
		return DetailTimeout
	}
	return DetailFailed
}

func (b batch) mark(reasons []string, reason string) {
	for _, at := range b.indexes {
		reasons[at] = reason
	}
}

// batch is the indexes of at most batchSize pull requests and the extra fields to ask for.
type batch struct {
	indexes []int
	extra   string
}

func batched(indexes []int, extra string) []batch {
	var batches []batch
	for len(indexes) > 0 {
		n := min(batchSize, len(indexes))
		batches = append(batches, batch{indexes[:n], extra})
		indexes = indexes[n:]
	}
	return batches
}

// read runs one details call for the batch and fills its pull requests.
func (b batch) read(ctx context.Context, runner process.Runner, ids []string, prs []PullRequest) (string, error) {
	quoted := make([]string, len(b.indexes))
	for i, at := range b.indexes {
		quoted[i] = `"` + ids[at] + `"`
	}
	stdout, err := ghapi.Run(ctx, runner, []string{"api", "graphql", "-f", "query=" + fmt.Sprintf(detailsQuery, strings.Join(quoted, ","), b.extra)})
	if errors.Is(err, ghapi.ErrNotPermitted) {
		return "", fmt.Errorf("%w: GitHub refused the details", ghapi.ErrFailed)
	}
	if err != nil && !errors.Is(err, ghapi.ErrFailed) {
		return "", err
	}
	var body struct {
		Data *struct {
			Viewer struct {
				Login string `json:"login"`
			} `json:"viewer"`
			Nodes []*detailNode `json:"nodes"`
		} `json:"data"`
	}
	if json.Unmarshal(stdout, &body) != nil || body.Data == nil || body.Data.Viewer.Login == "" || len(body.Data.Nodes) != len(b.indexes) {
		return "", fmt.Errorf("%w: unreadable GitHub details", ghapi.ErrFailed)
	}
	for i, n := range body.Data.Nodes {
		if at := b.indexes[i]; n != nil && n.ID == ids[at] {
			n.fill(&prs[at])
		}
	}
	return body.Data.Viewer.Login, nil
}

// detailNode is one pull request in the GraphQL answer; null when GitHub refused it.
type detailNode struct {
	ID               string  `json:"id"`
	ReviewDecision   *string `json:"reviewDecision"`
	Mergeable        *string `json:"mergeable"`
	MergeStateStatus *string `json:"mergeStateStatus"`
	Commits          struct {
		Nodes []struct {
			Commit struct {
				StatusCheckRollup *struct {
					State string `json:"state"`
				} `json:"statusCheckRollup"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"commits"`
	ReviewThreads struct {
		Nodes []struct {
			IsResolved bool `json:"isResolved"`
		} `json:"nodes"`
	} `json:"reviewThreads"`
}

// fill copies the details into pr; unresolved threads only for the user's own.
func (n *detailNode) fill(pr *PullRequest) {
	pr.Detail, pr.ReviewDecision, pr.Mergeable, pr.MergeState = "available", n.ReviewDecision, n.Mergeable, n.MergeStateStatus
	if c := n.Commits.Nodes; len(c) == 1 && c[0].Commit.StatusCheckRollup != nil {
		pr.Checks = &c[0].Commit.StatusCheckRollup.State
	}
	if hasReason(*pr, "author") {
		unresolved := 0
		for _, t := range n.ReviewThreads.Nodes {
			if !t.IsResolved {
				unresolved++
			}
		}
		pr.UnresolvedThreads = &unresolved
	}
}

// searchNeeds and detailNeeds follow the table in the contract, in its order;
// each reason comes from one field.
func searchNeeds(pr PullRequest) []string {
	return needs([]rule{
		{"review_requested", hasReason(pr, "review_requested")},
		{"assigned", hasReason(pr, "assigned")},
	})
}

func detailNeeds(pr PullRequest) []string {
	is := func(v *string, values ...string) bool {
		for _, want := range values {
			if v != nil && *v == want {
				return true
			}
		}
		return false
	}
	own := hasReason(pr, "author")
	return needs([]rule{
		{"changes_requested", own && is(pr.ReviewDecision, "CHANGES_REQUESTED")},
		{"checks_failed", own && is(pr.Checks, "FAILURE", "ERROR")},
		{"conflict", own && is(pr.Mergeable, "CONFLICTING")},
		{"unresolved_threads", own && pr.UnresolvedThreads != nil && *pr.UnresolvedThreads > 0},
		{"ready_to_merge", own && is(pr.MergeState, "CLEAN")},
	})
}

type rule struct {
	name string
	when bool
}

func needs(rules []rule) []string {
	out := []string{}
	for _, r := range rules {
		if r.when {
			out = append(out, r.name)
		}
	}
	return out
}

func hasReason(pr PullRequest, kind string) bool {
	for _, r := range pr.Reasons {
		if r.Kind == kind {
			return true
		}
	}
	return false
}

func appendOnce(list []string, value string) []string {
	for _, v := range list {
		if v == value {
			return list
		}
	}
	return append(list, value)
}
