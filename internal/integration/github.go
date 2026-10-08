package integration

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PRMergeTarget is copied from the published Delivery receipt. OperationID is
// recovery evidence, never a mutable part of a human-confirmed plan.
type PRMergeTarget struct {
	Worktree    string `json:"worktree"`
	Repository  string `json:"repository"`
	Number      int    `json:"number"`
	HeadRef     string `json:"head_ref"`
	HeadOID     string `json:"head_oid"`
	BaseRef     string `json:"base_ref"`
	Method      string `json:"method"`
	OperationID string `json:"operation_id,omitempty"`
}

type PRMergeObservation struct {
	State string `json:"state"`
	// NoEffect is positive evidence about this attempt, not an inference from
	// an open PR. Only a separate explicit retry may use it.
	NoEffect        bool     `json:"no_effect,omitempty"`
	Reason          string   `json:"reason,omitempty"`
	NextAction      string   `json:"next_action"`
	Repository      string   `json:"repository"`
	Number          int      `json:"number"`
	HeadRef         string   `json:"head_ref"`
	HeadOID         string   `json:"head_oid"`
	BaseRef         string   `json:"base_ref"`
	ObservedBaseOID string   `json:"observed_base_oid,omitempty"`
	MergeOID        string   `json:"merge_oid,omitempty"`
	URL             string   `json:"url,omitempty"`
	MergedAt        string   `json:"merged_at,omitempty"`
	OperationID     string   `json:"operation_id,omitempty"`
	AllowedMethods  []string `json:"allowed_methods,omitempty"`
}

// Callers durably reserve the effect before Merge. An incomplete reservation
// permits Observe only: observing an open PR does not prove that a timed-out
// request failed to reach GitHub.
type PRMergeAdapter interface {
	Observe(PRMergeTarget) (PRMergeObservation, error)
	Merge(PRMergeTarget) (PRMergeObservation, error)
}

type GitHubMergeAdapter struct {
	Run func(cwd, program string, args ...string) ([]byte, error)
}

var mergeRepositoryPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*/[A-Za-z0-9][A-Za-z0-9_.-]*$`)
var mergeOIDPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var mergeOperationPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func validMergeRef(s string) bool {
	return strings.HasPrefix(s, "refs/heads/") && len(s) > len("refs/heads/") &&
		!strings.ContainsAny(s, " \t\r\n\\~^:?*[") && !strings.Contains(s, "..") &&
		!strings.Contains(s, "@{") && !strings.Contains(s, "//") && !strings.HasSuffix(s, "/") &&
		!strings.HasSuffix(s, ".") && !strings.HasSuffix(s, ".lock")
}

func ValidatePRMergeTarget(t PRMergeTarget) error {
	if !mergeRepositoryPattern.MatchString(t.Repository) || t.Number < 1 ||
		!validMergeRef(t.HeadRef) || !validMergeRef(t.BaseRef) || t.HeadRef == t.BaseRef ||
		!mergeOIDPattern.MatchString(t.HeadOID) || (t.Method != "merge" && t.Method != "squash") ||
		(t.OperationID != "" && !mergeOperationPattern.MatchString(t.OperationID)) {
		return fmt.Errorf("merge requires a published PR, exact repository/head/base, and explicit merge or squash method")
	}
	return nil
}

type adapterBuffer struct {
	bytes.Buffer
	limit int
	over  bool
}

func (b *adapterBuffer) Write(p []byte) (int, error) {
	n := len(p)
	left := b.limit - b.Len()
	if len(p) > left {
		b.over = true
		p = p[:left]
	}
	_, _ = b.Buffer.Write(p)
	return n, nil
}

func runAdapterCommand(cwd, program string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(), "GH_PROMPT_DISABLED=1", "GH_PAGER=cat", "GIT_TERMINAL_PROMPT=0")
	stdout, stderr := adapterBuffer{limit: 2 << 20}, adapterBuffer{limit: 8 << 10}
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if stdout.over {
		return nil, fmt.Errorf("adapter response exceeded its observation limit")
	}
	if err != nil {
		// The HTTP response, when present, is the only authoritative source for
		// remote classification. Do not expose CLI stderr or credentials.
		return stdout.Bytes(), fmt.Errorf("adapter command failed: %w", err)
	}
	return stdout.Bytes(), nil
}

type mergeHTTP struct {
	status int
	body   []byte
}

func (a *GitHubMergeAdapter) request(t PRMergeTarget, method, endpoint string, fields ...string) (mergeHTTP, error) {
	run := a.Run
	if run == nil {
		run = runAdapterCommand
	}
	args := []string{"api", "--hostname", "github.com", "--method", method, "--include", "--header", "Accept: application/vnd.github+json", "--header", "X-GitHub-Api-Version: 2026-03-10", endpoint}
	args = append(args, fields...)
	raw, runErr := run(t.Worktree, "gh", args...)
	// --include retains the status even when gh exits nonzero for HTTP errors.
	response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(raw)), nil)
	if err != nil {
		if runErr != nil {
			return mergeHTTP{}, runErr
		}
		return mergeHTTP{}, fmt.Errorf("GitHub response lacks a valid HTTP status")
	}
	defer response.Body.Close()
	var body bytes.Buffer
	if _, err = body.ReadFrom(response.Body); err != nil || body.Len() > 2<<20 {
		return mergeHTTP{}, fmt.Errorf("GitHub response body is incomplete or too large")
	}
	if runErr != nil && response.StatusCode >= 200 && response.StatusCode < 300 {
		return mergeHTTP{}, runErr
	}
	return mergeHTTP{status: response.StatusCode, body: body.Bytes()}, nil
}

func mergeFailure(o PRMergeObservation, state, reason, next string) PRMergeObservation {
	o.State, o.Reason, o.NextAction = state, reason, next
	return o
}

func mergeUnknown(o PRMergeObservation) PRMergeObservation {
	o.NoEffect = false
	return mergeFailure(o, "effect_unknown", "effect_unknown", "Resume this integration to observe the same PR and preserved merge operation; do not submit another merge.")
}

func prEndpoint(t PRMergeTarget) string {
	return "repos/" + t.Repository + "/pulls/" + strconv.Itoa(t.Number)
}

type mergePRWire struct {
	Number         int             `json:"number"`
	URL            string          `json:"html_url"`
	State          string          `json:"state"`
	Draft          bool            `json:"draft"`
	Head           mergeRefWire    `json:"head"`
	Base           mergeRefWire    `json:"base"`
	Merged         *bool           `json:"merged"`
	MergedAt       string          `json:"merged_at"`
	MergeOID       string          `json:"merge_commit_sha"`
	Mergeable      *bool           `json:"mergeable"`
	MergeableState string          `json:"mergeable_state"`
	Stack          json.RawMessage `json:"stack"`
}

type mergeRefWire struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	Repo struct {
		FullName string `json:"full_name"`
	} `json:"repo"`
}

func exactMergeURL(raw string, t PRMergeTarget) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host == "github.com" && u.User == nil && u.RawQuery == "" && u.Fragment == "" &&
		strings.EqualFold(u.Path, "/"+t.Repository+"/pull/"+strconv.Itoa(t.Number))
}

func (a *GitHubMergeAdapter) observePR(t PRMergeTarget) (PRMergeObservation, *mergePRWire, error) {
	o := PRMergeObservation{Repository: t.Repository, Number: t.Number, OperationID: t.OperationID}
	if err := ValidatePRMergeTarget(t); err != nil {
		return mergeFailure(o, "blocked", "invalid_pr_target", "Select the published Delivery PR and an explicitly allowed merge method."), nil, err
	}
	r, err := a.request(t, "GET", prEndpoint(t))
	if err != nil {
		return mergeUnknown(o), nil, err
	}
	if r.status != http.StatusOK {
		return mergeFailure(o, "blocked", "merge_blocked", "Restore read access to the selected PR and resume; no merge was requested."), nil, nil
	}
	var p mergePRWire
	if err = json.Unmarshal(r.body, &p); err != nil {
		return mergeUnknown(o), nil, fmt.Errorf("invalid PR response")
	}
	o.HeadRef, o.HeadOID, o.BaseRef, o.ObservedBaseOID, o.URL = "refs/heads/"+p.Head.Ref, p.Head.SHA, "refs/heads/"+p.Base.Ref, p.Base.SHA, p.URL
	if p.Number != t.Number || !strings.EqualFold(p.Head.Repo.FullName, t.Repository) || !strings.EqualFold(p.Base.Repo.FullName, t.Repository) || o.HeadRef != t.HeadRef || o.BaseRef != t.BaseRef || !exactMergeURL(p.URL, t) {
		return mergeFailure(o, "blocked", "pr_target_changed", "The PR repository, branches or URL changed; prepare a new explicit plan."), &p, nil
	}
	if p.Head.SHA != t.HeadOID {
		return mergeFailure(o, "blocked", "pr_head_changed", "Qualify the changed head and prepare a new plan; the old candidate cannot be merged."), &p, nil
	}
	if p.Merged == nil {
		return mergeUnknown(o), &p, fmt.Errorf("PR response omits merge observation")
	}
	if *p.Merged {
		if p.State != "closed" || !mergeOIDPattern.MatchString(p.MergeOID) {
			return mergeUnknown(o), &p, fmt.Errorf("PR merge receipt is incomplete")
		}
		if _, err = time.Parse(time.RFC3339, p.MergedAt); err != nil {
			return mergeUnknown(o), &p, fmt.Errorf("PR merge receipt omits its timestamp")
		}
		o.State, o.MergeOID, o.MergedAt, o.NextAction = "merged", p.MergeOID, p.MergedAt, "Continue the remaining local closeout steps."
		return o, &p, nil
	}
	if p.State != "open" {
		return mergeFailure(o, "blocked", "pr_closed_unmerged", "The selected PR is closed without merge; resolve it before integration."), &p, nil
	}
	o.State, o.NextAction = "ready", "Confirm the exact PR integration plan."
	return o, &p, nil
}

type mergeAsyncWire struct {
	Status  string `json:"status"`
	Details struct {
		UUID            string `json:"uuid"`
		MergeMethod     string `json:"merge_method"`
		MergeAction     string `json:"merge_action"`
		ExpectedHeadSHA string `json:"expected_head_sha"`
		BypassRules     *bool  `json:"bypass_rules"`
		SHA             string `json:"sha"`
	} `json:"details"`
}

func pendingBinding(w mergeAsyncWire, t PRMergeTarget) bool {
	return mergeOperationPattern.MatchString(w.Details.UUID) && (t.OperationID == "" || t.OperationID == w.Details.UUID) &&
		w.Details.MergeMethod == t.Method && w.Details.MergeAction == "direct_merge" &&
		w.Details.ExpectedHeadSHA == t.HeadOID && w.Details.BypassRules != nil && !*w.Details.BypassRules
}

func (a *GitHubMergeAdapter) observeOperation(t PRMergeTarget, o PRMergeObservation) (PRMergeObservation, error) {
	r, err := a.request(t, "GET", prEndpoint(t)+"/merge-async/"+t.OperationID)
	if err != nil || r.status != http.StatusOK {
		return mergeUnknown(o), err
	}
	var w mergeAsyncWire
	if json.Unmarshal(r.body, &w) != nil {
		return mergeUnknown(o), fmt.Errorf("invalid asynchronous merge observation")
	}
	switch w.Status {
	case "pending":
		if !pendingBinding(w, t) {
			return mergeUnknown(o), fmt.Errorf("asynchronous merge binding differs from the confirmed plan")
		}
		return mergeFailure(o, "remote_pending", "remote_pending", "Resume this integration ID to observe the preserved remote operation; keep the Task worktree."), nil
	case "enqueued":
		return mergeFailure(o, "remote_pending", "remote_pending", "The PR is queued but no merge is observed; resume this integration without submitting another request."), nil
	case "failed":
		// Reobserve the exact unmerged PR after the terminal remote result. An
		// earlier open-PR read does not establish the current candidate identity.
		fresh, _, readErr := a.observePR(t)
		if readErr != nil {
			return mergeUnknown(o), readErr
		}
		if fresh.State != "ready" {
			return fresh, nil
		}
		fresh.NoEffect = true
		return mergeFailure(fresh, "blocked", "merge_blocked", "GitHub reports the preserved merge failed; resolve branch rules, checks or reviews before an explicit retry."), nil
	case "merged":
		// The operation result alone lacks the exact current PR identity and
		// timestamp. Read the PR again before reporting an observed merge.
		fresh, _, err := a.observePR(t)
		if err != nil || fresh.State != "merged" || fresh.MergeOID != w.Details.SHA {
			return mergeUnknown(o), err
		}
		return fresh, nil
	default:
		return mergeUnknown(o), fmt.Errorf("unrecognized asynchronous merge state")
	}
}

func (a *GitHubMergeAdapter) Observe(t PRMergeTarget) (PRMergeObservation, error) {
	o, p, err := a.observePR(t)
	if err != nil || o.State != "ready" {
		return o, err
	}
	if t.OperationID != "" {
		return a.observeOperation(t, o)
	}
	if len(p.Stack) > 0 && string(p.Stack) != "null" {
		return mergeFailure(o, "blocked", "stacked_pr_unsupported", "Stacked PR integration is outside v1; select a standalone PR."), nil
	}
	if p.Draft || p.Mergeable == nil || !*p.Mergeable || p.MergeableState != "clean" {
		return mergeFailure(o, "blocked", "merge_blocked", "Wait for GitHub to report a clean, mergeable PR with satisfied checks and reviews; merge queues are not requested."), nil
	}
	r, err := a.request(t, "GET", "repos/"+t.Repository)
	if err != nil {
		return mergeUnknown(o), err
	}
	var settings struct {
		FullName    string `json:"full_name"`
		AllowMerge  *bool  `json:"allow_merge_commit"`
		AllowSquash *bool  `json:"allow_squash_merge"`
		Permissions struct {
			Push *bool `json:"push"`
		} `json:"permissions"`
	}
	if r.status != http.StatusOK || json.Unmarshal(r.body, &settings) != nil || !strings.EqualFold(settings.FullName, t.Repository) || settings.Permissions.Push == nil || !*settings.Permissions.Push {
		return mergeFailure(o, "blocked", "merge_permission_missing", "Obtain normal write permission to the selected repository; administrator bypass is not supported."), nil
	}
	if settings.AllowMerge != nil && *settings.AllowMerge {
		o.AllowedMethods = append(o.AllowedMethods, "merge")
	}
	if settings.AllowSquash != nil && *settings.AllowSquash {
		o.AllowedMethods = append(o.AllowedMethods, "squash")
	}
	for _, method := range o.AllowedMethods {
		if method == t.Method {
			return o, nil
		}
	}
	return mergeFailure(o, "blocked", "merge_method_unavailable", "The repository disallows the selected method; choose a permitted method in a new plan."), nil
}

func (a *GitHubMergeAdapter) Merge(t PRMergeTarget) (PRMergeObservation, error) {
	// A known operation is observation-only, including an already failed one.
	// A caller cannot accidentally repeat it by reusing this entrypoint.
	o, err := a.Observe(t)
	if err != nil || o.State != "ready" || t.OperationID != "" {
		if err == nil && o.State == "blocked" && t.OperationID == "" {
			// This invocation stopped in preflight before any merge request. Do
			// not add this claim to Observe, which may follow an unknown attempt.
			o.NoEffect = true
		}
		return o, err
	}
	r, err := a.request(t, "PUT", prEndpoint(t)+"/merge-async", "--raw-field", "sha="+t.HeadOID, "--raw-field", "merge_method="+t.Method, "--raw-field", "merge_action=direct_merge", "--field", "bypass_rules=false")
	if err != nil {
		return mergeUnknown(o), err
	}
	var w mergeAsyncWire
	parsed := json.Unmarshal(r.body, &w) == nil
	if parsed && (r.status == http.StatusAccepted || r.status == http.StatusConflict) && w.Status == "pending" {
		if !pendingBinding(w, t) {
			return mergeUnknown(o), fmt.Errorf("merge acknowledgement does not match the confirmed request")
		}
		t.OperationID, o.OperationID = w.Details.UUID, w.Details.UUID
		// Bound the wait to one read. Further progress is an explicit resume.
		return a.observeOperation(t, o)
	}
	if parsed && r.status == http.StatusOK && w.Status == "enqueued" {
		return mergeFailure(o, "remote_pending", "remote_pending", "GitHub reports a queue position, not integration; resume this ID to observe the PR."), nil
	}
	if parsed && r.status == http.StatusOK && w.Status == "merged" {
		fresh, _, readErr := a.observePR(t)
		if readErr != nil || fresh.State != "merged" || fresh.MergeOID != w.Details.SHA {
			return mergeUnknown(o), readErr
		}
		return fresh, nil
	}
	if r.status == http.StatusBadRequest || r.status == http.StatusForbidden || r.status == http.StatusNotFound || r.status == http.StatusMethodNotAllowed || r.status == http.StatusUnprocessableEntity {
		o.NoEffect = true
		return mergeFailure(o, "blocked", "merge_blocked", "GitHub rejected this merge; resolve head, branch rules, checks, reviews or permissions before preparing a new attempt."), nil
	}
	// A conflict can identify an existing request. Without its exact binding,
	// neither rejection nor permission to retry has been established.
	return mergeUnknown(o), nil
}
