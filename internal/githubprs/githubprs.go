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

// Deadline bounds the whole read: four searches and one GraphQL call.
const Deadline = 40 * time.Second

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

// Read runs the searches and the details call within Deadline. Errors carry
// one of the ghapi codes.
func Read(ctx context.Context, runner process.Runner) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, Deadline)
	defer cancel()
	result := Result{IncompleteReasons: []string{}, PullRequests: []PullRequest{}}
	hits := map[string][]searchItem{}
	for _, q := range searches {
		items, incomplete, err := searchPRs(ctx, runner, q)
		if err != nil {
			return Result{}, err
		}
		hits[q] = items
		for _, reason := range incomplete {
			result.IncompleteReasons = appendOnce(result.IncompleteReasons, reason)
		}
	}
	direct := map[string]bool{}
	for _, it := range hits["user-review-requested:@me"] {
		direct[it.NodeID] = true
	}
	byID := map[string]int{}
	var ids []string
	for _, q := range searches {
		for _, it := range hits[q] {
			i, seen := byID[it.NodeID]
			if !seen {
				i = len(result.PullRequests)
				byID[it.NodeID] = i
				ids = append(ids, it.NodeID)
				result.PullRequests = append(result.PullRequests, it.pullRequest())
			}
			if r, ok := reasonOf(q, direct[it.NodeID]); ok && !hasReason(result.PullRequests[i], r.Kind) {
				result.PullRequests[i].Reasons = append(result.PullRequests[i].Reasons, r)
			}
		}
	}
	viewer, err := details(ctx, runner, ids, result.PullRequests)
	if err != nil {
		return Result{}, err
	}
	result.Viewer = viewer
	for i := range result.PullRequests {
		result.PullRequests[i].NeedsYou = needsYou(result.PullRequests[i])
	}
	return result, nil
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

const detailsQuery = `query{viewer{login} nodes(ids:[%s]){... on PullRequest{id reviewDecision mergeable mergeStateStatus ` +
	`commits(last:1){nodes{commit{statusCheckRollup{state}}}} reviewThreads(first:100){nodes{isResolved}}}}}`

// details fills the detail fields of prs, in the order of ids, and returns the
// viewer. A node GitHub refuses stays unavailable; a call without data fails.
func details(ctx context.Context, runner process.Runner, ids []string, prs []PullRequest) (string, error) {
	quoted := make([]string, len(ids))
	for i, id := range ids {
		quoted[i] = `"` + id + `"`
	}
	stdout, err := ghapi.Run(ctx, runner, []string{"api", "graphql", "-f", "query=" + fmt.Sprintf(detailsQuery, strings.Join(quoted, ","))})
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
	if json.Unmarshal(stdout, &body) != nil || body.Data == nil || body.Data.Viewer.Login == "" || len(body.Data.Nodes) != len(ids) {
		return "", fmt.Errorf("%w: unreadable GitHub details", ghapi.ErrFailed)
	}
	for i, n := range body.Data.Nodes {
		if n != nil && n.ID == ids[i] {
			n.fill(&prs[i])
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

// needsYou follows the table in the contract; each reason comes from one field.
func needsYou(pr PullRequest) []string {
	is := func(v *string, values ...string) bool {
		for _, want := range values {
			if v != nil && *v == want {
				return true
			}
		}
		return false
	}
	own := hasReason(pr, "author")
	needs := []string{}
	for _, rule := range []struct {
		name string
		when bool
	}{
		{"review_requested", hasReason(pr, "review_requested")},
		{"assigned", hasReason(pr, "assigned")},
		{"changes_requested", own && is(pr.ReviewDecision, "CHANGES_REQUESTED")},
		{"checks_failed", own && is(pr.Checks, "FAILURE", "ERROR")},
		{"conflict", own && is(pr.Mergeable, "CONFLICTING")},
		{"unresolved_threads", own && pr.UnresolvedThreads != nil && *pr.UnresolvedThreads > 0},
		{"ready_to_merge", own && is(pr.MergeState, "CLEAN")},
	} {
		if rule.when {
			needs = append(needs, rule.name)
		}
	}
	return needs
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
