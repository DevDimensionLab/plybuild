package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PRTarget struct {
	Worktree, Remote, Repository, SourceRef, BaseRef, OID, EffectKey string
}
type PullRequest struct {
	Number     int      `json:"number"`
	URL        string   `json:"url"`
	HeadRef    string   `json:"head_ref"`
	HeadOID    string   `json:"head_oid"`
	BaseRef    string   `json:"base_ref"`
	Repository string   `json:"repository"`
	Body       string   `json:"body"`
	Assignees  []string `json:"assignees"`
	Reviewers  []string `json:"reviewers"`
}
type PRObservation struct {
	RemoteOID    string
	BaseExists   bool
	PullRequests []PullRequest
}

// PRAdapter contains only the three explicitly authorized publication effects.
// It has no merge, rebase, force-push, target-branch update or notification API.
type PRAdapter interface {
	Observe(PRTarget) (PRObservation, error)
	CanFastForward(PRTarget, string) (bool, error)
	Push(PRTarget) error
	Create(PRTarget, string, string) error
	ApplyMetadata(PRTarget, PullRequest, MetadataChoice) error
}

// Run accepts argv directly, never shell text. Tests can return controlled
// command responses without any real network, GitHub account or credential.
type GitHubAdapter struct {
	Run func(cwd, program string, args ...string) ([]byte, error)
}

var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*/[A-Za-z0-9][A-Za-z0-9_.-]*$`)
var oidPattern = regexp.MustCompile(`^(?:[a-f0-9]{40}|[a-f0-9]{64})$`)

type limitedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if b.Len() < b.limit {
		left := b.limit - b.Len()
		if len(p) > left {
			p = p[:left]
		}
		_, _ = b.Buffer.Write(p)
	}
	return n, nil
}

func (a *GitHubAdapter) run(cwd, program string, args ...string) ([]byte, error) {
	if a.Run != nil {
		return a.Run(cwd, program, args...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir = cwd
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GH_PROMPT_DISABLED=1", "GH_PAGER=cat")
	var stdout, stderr limitedBuffer
	stdout.limit = 2 << 20
	stderr.limit = 8 << 10
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("%s failed: %w: %s", program, err, strings.TrimSpace(stderr.String()))
	}
	if stdout.Len() >= stdout.limit {
		return nil, fmt.Errorf("%s observation exceeded output limit", program)
	}
	return stdout.Bytes(), nil
}

func githubRepository(remote string) (string, error) {
	remote = strings.TrimSpace(remote)
	path := ""
	if strings.HasPrefix(remote, "git@github.com:") {
		path = strings.TrimPrefix(remote, "git@github.com:")
	} else {
		u, e := url.Parse(remote)
		if e != nil || u.Host != "github.com" || u.RawQuery != "" || u.Fragment != "" || !(u.Scheme == "https" || u.Scheme == "ssh") {
			return "", fmt.Errorf("remote must name the agreed GitHub repository")
		}
		if u.User != nil {
			if u.Scheme != "ssh" || u.User.Username() != "git" {
				return "", fmt.Errorf("remote has an unsupported credential or user binding")
			}
			if _, ok := u.User.Password(); ok {
				return "", fmt.Errorf("remote has an unsupported credential binding")
			}
		}
		path = strings.TrimPrefix(u.Path, "/")
	}
	path = strings.TrimSuffix(path, ".git")
	if !repositoryPattern.MatchString(path) {
		return "", fmt.Errorf("remote must identify exactly one owner/repository")
	}
	return path, nil
}

func (a *GitHubAdapter) remote(t PRTarget) error {
	for _, push := range []bool{false, true} {
		args := []string{"remote", "get-url", "--all"}
		if push {
			args = append(args, "--push")
		}
		args = append(args, t.Remote)
		raw, e := a.run(t.Worktree, "git", args...)
		if e != nil {
			return e
		}
		lines := strings.Fields(string(raw))
		if len(lines) != 1 {
			return fmt.Errorf("delivery requires exactly one matching remote URL")
		}
		repo, e := githubRepository(lines[0])
		if e != nil {
			return e
		}
		if !strings.EqualFold(repo, t.Repository) {
			return fmt.Errorf("remote repository differs from the frozen delivery agreement")
		}
	}
	return nil
}

func (a *GitHubAdapter) Observe(t PRTarget) (PRObservation, error) {
	var out PRObservation
	if e := a.remote(t); e != nil {
		return out, e
	}
	raw, e := a.run(t.Worktree, "git", "ls-remote", "--heads", t.Remote, t.SourceRef, t.BaseRef)
	if e != nil {
		return out, e
	}
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		if line == "" {
			continue
		}
		v := strings.Fields(line)
		if len(v) != 2 || !oidPattern.MatchString(v[0]) || seen[v[1]] {
			return out, fmt.Errorf("remote ref observation is ambiguous")
		}
		seen[v[1]] = true
		switch v[1] {
		case t.SourceRef:
			out.RemoteOID = v[0]
		case t.BaseRef:
			out.BaseExists = true
		default:
			return out, fmt.Errorf("remote ref observation returned an unrequested ref")
		}
	}
	raw, e = a.run(t.Worktree, "gh", "pr", "list", "--repo", "github.com/"+t.Repository, "--state", "open", "--head", strings.TrimPrefix(t.SourceRef, "refs/heads/"), "--limit", "1000", "--json", "number,url,headRefName,headRefOid,baseRefName,headRepository,headRepositoryOwner,isCrossRepository,body,assignees,reviewRequests")
	if e != nil {
		return out, e
	}
	var wire []struct {
		Number         int    `json:"number"`
		URL            string `json:"url"`
		HeadRefName    string `json:"headRefName"`
		HeadRefOID     string `json:"headRefOid"`
		BaseRefName    string `json:"baseRefName"`
		Body           string `json:"body"`
		HeadRepository struct {
			Name string `json:"name"`
		} `json:"headRepository"`
		HeadRepositoryOwner struct {
			Login string `json:"login"`
		} `json:"headRepositoryOwner"`
		IsCrossRepository bool `json:"isCrossRepository"`
		Assignees         []struct {
			Login string `json:"login"`
		} `json:"assignees"`
		ReviewRequests []struct {
			Login string `json:"login"`
			Name  string `json:"name"`
			Slug  string `json:"slug"`
		} `json:"reviewRequests"`
	}
	if e = json.Unmarshal(raw, &wire); e != nil {
		return out, fmt.Errorf("invalid GitHub PR observation: %w", e)
	}
	if len(wire) >= 1000 {
		return out, fmt.Errorf("PR lookup limit reached; uniqueness cannot be established")
	}
	for _, w := range wire {
		repo := w.HeadRepositoryOwner.Login + "/" + w.HeadRepository.Name
		if w.IsCrossRepository || !strings.EqualFold(repo, t.Repository) {
			return out, fmt.Errorf("PR head repository differs from the selected repository")
		}
		p := PullRequest{Number: w.Number, URL: w.URL, HeadRef: "refs/heads/" + w.HeadRefName, HeadOID: w.HeadRefOID, BaseRef: "refs/heads/" + w.BaseRefName, Repository: t.Repository, Body: w.Body, Assignees: []string{}, Reviewers: []string{}}
		for _, v := range w.Assignees {
			p.Assignees = append(p.Assignees, v.Login)
		}
		for _, v := range w.ReviewRequests {
			if v.Login != "" {
				p.Reviewers = append(p.Reviewers, v.Login)
			} else if v.Slug != "" {
				p.Reviewers = append(p.Reviewers, strings.Split(t.Repository, "/")[0]+"/"+v.Slug)
			} else {
				return out, fmt.Errorf("unknown GitHub reviewer identity")
			}
		}
		out.PullRequests = append(out.PullRequests, p)
	}
	return out, nil
}

func validateObservation(t PRTarget, o PRObservation) (*PullRequest, error) {
	if !o.BaseExists {
		return nil, fmt.Errorf("agreed PR base branch is missing")
	}
	if len(o.PullRequests) > 1 {
		return nil, fmt.Errorf("ambiguous PR matches for the selected source branch")
	}
	if len(o.PullRequests) == 0 {
		return nil, nil
	}
	p := o.PullRequests[0]
	if p.Number < 1 || !strings.EqualFold(p.Repository, t.Repository) || p.HeadRef != t.SourceRef || p.BaseRef != t.BaseRef || p.HeadOID != t.OID {
		return nil, fmt.Errorf("existing PR head, candidate or base differs from the agreed target")
	}
	u, e := url.Parse(p.URL)
	if e != nil || u.Scheme != "https" || u.Host != "github.com" || !strings.EqualFold(strings.TrimSuffix(u.Path, "/"), "/"+t.Repository+"/pull/"+strconv.Itoa(p.Number)) || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return nil, fmt.Errorf("PR URL does not identify the observed repository and PR")
	}
	return &p, nil
}

func (a *GitHubAdapter) CanFastForward(t PRTarget, old string) (bool, error) {
	_, e := a.run(t.Worktree, "git", "merge-base", "--is-ancestor", old, t.OID)
	if e == nil {
		return true, nil
	}
	var exit *exec.ExitError
	if errors.As(e, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return false, e
}

func (a *GitHubAdapter) Push(t PRTarget) error {
	if e := a.remote(t); e != nil {
		return e
	}
	// Configured defaults must not expand this one-ref agreement into tags,
	// mirrored refs, submodule publication or arbitrary local pre-push hooks.
	_, e := a.run(t.Worktree, "git", "-c", "remote."+t.Remote+".mirror=false", "-c", "core.hooksPath=/dev/null", "push", "--no-follow-tags", "--recurse-submodules=no", "--no-verify", "--", t.Remote, t.OID+":"+t.SourceRef)
	return e
}
func marker(t PRTarget) string { return "<!-- ply-delivery-effect:" + t.EffectKey + " -->" }
func (a *GitHubAdapter) Create(t PRTarget, title, body string) error {
	f, e := os.CreateTemp("", "ply-delivery-pr-*.md")
	if e != nil {
		return e
	}
	path := f.Name()
	defer os.Remove(path)
	_, e = f.WriteString(body + "\n\n" + marker(t) + "\n")
	ce := f.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return e
	}
	_, e = a.run(t.Worktree, "gh", "pr", "create", "--repo", "github.com/"+t.Repository, "--head", strings.TrimPrefix(t.SourceRef, "refs/heads/"), "--base", strings.TrimPrefix(t.BaseRef, "refs/heads/"), "--title", title, "--body-file", path)
	return e
}

func metadataDifference(p PullRequest, m MetadataChoice) (addAssignees, addReviewers, removeAssignees, removeReviewers []string) {
	for _, v := range m.Assignees {
		if !includes(p.Assignees, v) {
			addAssignees = append(addAssignees, v)
		}
	}
	for _, v := range m.Reviewers {
		if !includes(p.Reviewers, v) {
			addReviewers = append(addReviewers, v)
		}
	}
	for _, v := range m.RemoveAssignees {
		if includes(p.Assignees, v) {
			removeAssignees = append(removeAssignees, v)
		}
	}
	for _, v := range m.RemoveReviewers {
		if includes(p.Reviewers, v) {
			removeReviewers = append(removeReviewers, v)
		}
	}
	return
}
func metadataSatisfied(p PullRequest, m MetadataChoice) bool {
	a, b, c, d := metadataDifference(p, m)
	return len(a)+len(b)+len(c)+len(d) == 0
}
func (a *GitHubAdapter) ApplyMetadata(t PRTarget, p PullRequest, m MetadataChoice) error {
	assignees, reviewers, removeA, removeR := metadataDifference(p, m)
	args := []string{"pr", "edit", strconv.Itoa(p.Number), "--repo", "github.com/" + t.Repository}
	for _, field := range []struct {
		name   string
		values []string
	}{{"--add-assignee", assignees}, {"--add-reviewer", reviewers}, {"--remove-assignee", removeA}, {"--remove-reviewer", removeR}} {
		if len(field.values) > 0 {
			args = append(args, field.name, strings.Join(field.values, ","))
		}
	}
	if len(args) == 5 {
		return nil
	}
	_, e := a.run(t.Worktree, "gh", args...)
	return e
}
