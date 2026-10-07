package delivery

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func testTarget() PRTarget {
	return PRTarget{Worktree: "/test/worktree", Remote: "origin", Repository: "Example/Project", SourceRef: "refs/heads/task-branch", BaseRef: "refs/heads/main", OID: strings.Repeat("a", 40), EffectKey: "sha256:" + strings.Repeat("b", 64)}
}
func testPR(t PRTarget) PullRequest {
	return PullRequest{Number: 7, URL: "https://github.com/Example/Project/pull/7", HeadRef: t.SourceRef, HeadOID: t.OID, BaseRef: t.BaseRef, Repository: t.Repository, Assignees: []string{"existing-owner"}, Reviewers: []string{"existing-reviewer"}}
}

func TestGitHubCommandsRestrictSourceAndPreservePeopleRoles(t *testing.T) {
	target := testTarget()
	var calls [][]string
	var body string
	a := GitHubAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
		if cwd != target.Worktree {
			t.Fatalf("wrong cwd %q", cwd)
		}
		calls = append(calls, append([]string{program}, args...))
		if program == "git" && args[0] == "remote" {
			return []byte("git@github.com:Example/Project.git\n"), nil
		}
		if program == "gh" && args[1] == "create" {
			for i, v := range args {
				if v == "--body-file" {
					raw, e := os.ReadFile(args[i+1])
					if e != nil {
						return nil, e
					}
					body = string(raw)
				}
			}
		}
		return []byte{}, nil
	}}
	if e := a.Push(target); e != nil {
		t.Fatal(e)
	}
	if e := a.Create(target, "A title with $(literal)", "A body with `literal`."); e != nil {
		t.Fatal(e)
	}
	p := testPR(target)
	m := MetadataChoice{Assignees: []string{"existing-owner", "new-owner"}, Reviewers: []string{"new-reviewer"}}
	if e := a.ApplyMetadata(target, p, m); e != nil {
		t.Fatal(e)
	}
	wantPush := []string{"git", "-c", "remote.origin.mirror=false", "-c", "core.hooksPath=/dev/null", "push", "--no-follow-tags", "--recurse-submodules=no", "--no-verify", "--", "origin", target.OID + ":" + target.SourceRef}
	if !reflect.DeepEqual(calls[2], wantPush) {
		t.Fatalf("push differs: %v", calls)
	}
	wantMetadata := []string{"gh", "pr", "edit", "7", "--repo", "github.com/Example/Project", "--add-assignee", "new-owner", "--add-reviewer", "new-reviewer"}
	if !reflect.DeepEqual(calls[len(calls)-1], wantMetadata) {
		t.Fatalf("roles changed or existing people removed: %v", calls)
	}
	if !strings.Contains(body, marker(target)) {
		t.Fatalf("create lacks recoverable effect marker: %q", body)
	}
	for _, call := range calls {
		for _, arg := range call {
			if arg == "--force" || arg == "merge" || arg == "--auto" || arg == target.BaseRef {
				t.Fatalf("forbidden effect: %v", call)
			}
		}
	}
	before := len(calls)
	if e := a.ApplyMetadata(target, p, MetadataChoice{}); e != nil || len(calls) != before {
		t.Fatalf("empty defaults removed existing people: %v", calls)
	}
}

func TestRemoteBindingRejectsWrongRepositoryAndMultiplePushTargets(t *testing.T) {
	for _, remote := range []string{"https://github.com/Other/Repo.git\n", "git@github.com:Example/Project.git\ngit@github.com:Other/Repo.git\n", "https://github.com.evil/Example/Project.git\n"} {
		calls := 0
		a := GitHubAdapter{Run: func(_, _ string, _ ...string) ([]byte, error) { calls++; return []byte(remote), nil }}
		if _, e := a.Observe(testTarget()); e == nil {
			t.Fatalf("accepted remote %q", remote)
		}
		if calls > 2 {
			t.Fatal("GitHub queried after conflicting remote")
		}
	}
}

func TestPRObservationRequiresExactHeadBaseAndUnambiguousIdentity(t *testing.T) {
	target := testTarget()
	p := testPR(target)
	for name, change := range map[string]func(*PRObservation){
		"ambiguous":    func(o *PRObservation) { o.PullRequests = append(o.PullRequests, p) },
		"head":         func(o *PRObservation) { o.PullRequests[0].HeadOID = strings.Repeat("c", 40) },
		"base":         func(o *PRObservation) { o.PullRequests[0].BaseRef = "refs/heads/other" },
		"repository":   func(o *PRObservation) { o.PullRequests[0].Repository = "Other/Repo" },
		"url":          func(o *PRObservation) { o.PullRequests[0].URL = "https://github.com/Other/Repo/pull/7" },
		"missing-base": func(o *PRObservation) { o.BaseExists = false },
	} {
		t.Run(name, func(t *testing.T) {
			o := PRObservation{RemoteOID: target.OID, BaseExists: true, PullRequests: []PullRequest{p}}
			change(&o)
			if _, e := validateObservation(target, o); e == nil {
				t.Fatalf("accepted conflicting observation: %+v", o)
			}
		})
	}
	if _, e := validateObservation(target, PRObservation{RemoteOID: target.OID, BaseExists: true, PullRequests: []PullRequest{p}}); e != nil {
		t.Fatal(e)
	}
}

func TestObserveUsesOnlyRemoteReadsAndExactGitHubQuery(t *testing.T) {
	target := testTarget()
	var calls [][]string
	row := map[string]any{"number": 7, "url": "https://github.com/Example/Project/pull/7", "headRefName": "task-branch", "headRefOid": target.OID, "baseRefName": "main", "headRepository": map[string]string{"name": "Project"}, "headRepositoryOwner": map[string]string{"login": "Example"}, "isCrossRepository": false, "body": "old", "assignees": []map[string]string{{"login": "owner"}}, "reviewRequests": []map[string]string{{"login": "reviewer"}, {"slug": "team"}}}
	a := GitHubAdapter{Run: func(_ string, program string, args ...string) ([]byte, error) {
		calls = append(calls, append([]string{program}, args...))
		switch {
		case program == "git" && args[0] == "remote":
			return []byte("https://github.com/Example/Project.git\n"), nil
		case program == "git" && args[0] == "ls-remote":
			return []byte(target.OID + "\t" + target.SourceRef + "\n" + strings.Repeat("f", 40) + "\t" + target.BaseRef + "\n"), nil
		case program == "gh" && args[0] == "pr" && args[1] == "list":
			return json.Marshal([]any{row})
		default:
			return nil, fmt.Errorf("unexpected command")
		}
	}}
	o, e := a.Observe(target)
	if e != nil {
		t.Fatal(e)
	}
	if !o.BaseExists || o.RemoteOID != target.OID || len(o.PullRequests) != 1 || !reflect.DeepEqual(o.PullRequests[0].Reviewers, []string{"reviewer", "Example/team"}) {
		t.Fatalf("observation: %+v", o)
	}
	for _, call := range calls {
		if call[1] == "push" || call[1] == "fetch" || call[1] == "merge" || call[1] == "rebase" {
			t.Fatalf("read had effect %v", call)
		}
	}
}

func TestPushPublishesOnlySelectedRefDespiteUserGitConfiguration(t *testing.T) {
	for _, setting := range []string{"follow-tags", "mirror", "pre-push-hook"} {
		t.Run(setting, func(t *testing.T) {
			root := physicalTemp(t)
			source := filepath.Join(root, "source")
			remote := filepath.Join(root, "remote.git")
			git := func(cwd string, args ...string) []byte {
				t.Helper()
				c := exec.Command("git", args...)
				c.Dir = cwd
				out, e := c.CombinedOutput()
				if e != nil {
					t.Fatalf("git %v: %v %s", args, e, out)
				}
				return out
			}
			git(root, "init", "--bare", remote)
			git(root, "init", "-b", "main", source)
			git(source, "config", "user.name", "Delivery test")
			git(source, "config", "user.email", "delivery@example.invalid")
			if e := os.WriteFile(filepath.Join(source, "file"), []byte("candidate\n"), 0600); e != nil {
				t.Fatal(e)
			}
			git(source, "add", "file")
			git(source, "commit", "-m", "candidate")
			git(source, "checkout", "-b", "task-branch")
			git(source, "remote", "add", "origin", remote)
			git(source, "tag", "-a", "unrelated-tag", "-m", "must remain local")
			hookObserved := filepath.Join(root, "hook-ran")
			switch setting {
			case "follow-tags":
				git(source, "config", "push.followTags", "true")
			case "mirror":
				git(source, "config", "remote.origin.mirror", "true")
			case "pre-push-hook":
				hook := "#!/bin/sh\nprintf touched > '" + hookObserved + "'\n"
				if e := os.WriteFile(filepath.Join(source, ".git", "hooks", "pre-push"), []byte(hook), 0700); e != nil {
					t.Fatal(e)
				}
			}
			target := testTarget()
			target.Worktree = source
			target.OID = strings.TrimSpace(string(git(source, "rev-parse", "HEAD")))
			a := GitHubAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				if len(args) > 1 && args[0] == "remote" && args[1] == "get-url" {
					return []byte("https://github.com/Example/Project.git\n"), nil
				}
				c := exec.Command(program, args...)
				c.Dir = cwd
				return c.CombinedOutput()
			}}
			if e := a.Push(target); e != nil {
				t.Fatalf("source-only push failed under %s: %v", setting, e)
			}
			refs := strings.TrimSpace(string(git(root, "--git-dir", remote, "for-each-ref", "--format=%(refname) %(objectname)")))
			if refs != target.SourceRef+" "+target.OID {
				t.Fatalf("push escaped source-only agreement: %s", refs)
			}
			if _, e := os.Stat(hookObserved); !os.IsNotExist(e) {
				t.Fatal("push executed an unrelated configured hook")
			}
		})
	}
}

func TestGitHubHostIsPinnedForEveryReadAndMutation(t *testing.T) {
	target := testTarget()
	hostChecks := 0
	a := GitHubAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
		if program == "git" {
			if args[0] == "remote" {
				return []byte("git@github.com:Example/Project.git\n"), nil
			}
			return []byte(target.OID + "\t" + target.SourceRef + "\n" + target.OID + "\t" + target.BaseRef), nil
		}
		for i, arg := range args {
			if arg == "--repo" {
				hostChecks++
				if args[i+1] != "github.com/"+target.Repository {
					return nil, fmt.Errorf("GH_HOST could redirect unqualified repository %s", args[i+1])
				}
			}
		}
		if args[1] == "list" {
			return []byte("[]"), nil
		}
		return []byte{}, nil
	}}
	if _, e := a.Observe(target); e != nil {
		t.Fatal(e)
	}
	if e := a.Create(target, "test", "body"); e != nil {
		t.Fatal(e)
	}
	if e := a.ApplyMetadata(target, testPR(target), MetadataChoice{Assignees: []string{"new-owner"}}); e != nil {
		t.Fatal(e)
	}
	if hostChecks != 3 {
		t.Fatalf("expected explicit host for every GitHub operation: %d", hostChecks)
	}
}
