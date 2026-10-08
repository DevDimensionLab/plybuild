package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func mergeTargetFixture() PRMergeTarget {
	return PRMergeTarget{Worktree: "/unused/integration-test", Repository: "Example/Project", Number: 31, HeadRef: "refs/heads/task", HeadOID: strings.Repeat("a", 40), BaseRef: "refs/heads/main", Method: "squash"}
}

func mergePRFixture(t PRMergeTarget) map[string]any {
	return map[string]any{"number": t.Number, "html_url": "https://github.com/" + t.Repository + "/pull/31", "state": "open", "draft": false, "head": map[string]any{"ref": "task", "sha": t.HeadOID, "repo": map[string]any{"full_name": t.Repository}}, "base": map[string]any{"ref": "main", "sha": strings.Repeat("b", 40), "repo": map[string]any{"full_name": t.Repository}}, "merged": false, "merged_at": nil, "merge_commit_sha": nil, "mergeable": true, "mergeable_state": "clean", "stack": nil}
}

func mergeSettingsFixture(t PRMergeTarget) map[string]any {
	return map[string]any{"full_name": t.Repository, "allow_merge_commit": true, "allow_squash_merge": true, "permissions": map[string]any{"push": true}}
}

func mergePendingFixture(t PRMergeTarget) map[string]any {
	return map[string]any{"status": "pending", "details": map[string]any{"uuid": "630b9d5e-3f2a-4f7e-8b0c-2d5f9a8c1e42", "merge_method": t.Method, "merge_action": "direct_merge", "expected_head_sha": t.HeadOID, "bypass_rules": false}}
}

func mergeHTTPFixture(status int, body any) []byte {
	raw, _ := json.Marshal(body)
	return []byte(fmt.Sprintf("HTTP/2.0 %d %s\r\nContent-Type: application/json\r\n\r\n%s", status, http.StatusText(status), raw))
}

func githubCall(t *testing.T, cwd, program string, args []string) (method, endpoint string) {
	t.Helper()
	if cwd != mergeTargetFixture().Worktree || program != "gh" || len(args) < 11 || args[0] != "api" || args[1] != "--hostname" || args[2] != "github.com" || args[3] != "--method" || args[5] != "--include" {
		t.Fatalf("unexpected or unpinned command %s %v at %s", program, args, cwd)
	}
	return args[4], args[10]
}

func TestGitHubMergeUsesExactHeadExplicitMethodAndNoBypass(t *testing.T) {
	for _, method := range []string{"merge", "squash"} {
		t.Run(method, func(t *testing.T) {
			target := mergeTargetFixture()
			target.Method = method
			p := mergePRFixture(target)
			puts, prReads := 0, 0
			actualMerge := strings.Repeat("c", 40)
			a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				m, path := githubCall(t, cwd, program, args)
				switch {
				case m == "GET" && path == prEndpoint(target):
					prReads++
					if prReads > 1 {
						p["merged"], p["state"], p["merged_at"], p["merge_commit_sha"] = true, "closed", "2026-10-08T10:11:12Z", actualMerge
					}
					return mergeHTTPFixture(200, p), nil
				case m == "GET" && path == "repos/"+target.Repository:
					return mergeHTTPFixture(200, mergeSettingsFixture(target)), nil
				case m == "PUT" && path == prEndpoint(target)+"/merge-async":
					puts++
					want := []string{"--raw-field", "sha=" + target.HeadOID, "--raw-field", "merge_method=" + method, "--raw-field", "merge_action=direct_merge", "--field", "bypass_rules=false"}
					if !reflect.DeepEqual(args[11:], want) {
						t.Fatalf("merge request escaped plan: %v", args)
					}
					return mergeHTTPFixture(202, mergePendingFixture(target)), nil
				case m == "GET" && strings.Contains(path, "/merge-async/"):
					return mergeHTTPFixture(200, map[string]any{"status": "merged", "details": map[string]any{"sha": actualMerge}}), nil
				default:
					t.Fatalf("unexpected command %v", args)
					return nil, nil
				}
			}}
			o, err := a.Merge(target)
			if err != nil || o.State != "merged" || o.HeadOID != target.HeadOID || o.MergeOID != actualMerge || o.MergeOID == o.HeadOID || o.BaseRef != target.BaseRef || o.MergedAt != "2026-10-08T10:11:12Z" || o.OperationID == "" || puts != 1 {
				t.Fatalf("merge receipt is incomplete: %+v %v puts=%d", o, err, puts)
			}
		})
	}
}

func TestGitHubMergeRejectsChangedOrBlockedPRWithoutMutation(t *testing.T) {
	for name, change := range map[string]func(map[string]any, map[string]any){
		"head":     func(p, _ map[string]any) { p["head"].(map[string]any)["sha"] = strings.Repeat("f", 40) },
		"head-ref": func(p, _ map[string]any) { p["head"].(map[string]any)["ref"] = "another-task" },
		"head-repository": func(p, _ map[string]any) {
			p["head"].(map[string]any)["repo"] = map[string]any{"full_name": "Other/Project"}
		},
		"base-ref": func(p, _ map[string]any) { p["base"].(map[string]any)["ref"] = "master" },
		"base-repository": func(p, _ map[string]any) {
			p["base"].(map[string]any)["repo"] = map[string]any{"full_name": "Other/Project"}
		},
		"url":                      func(p, _ map[string]any) { p["html_url"] = "https://github.com.evil/Example/Project/pull/31" },
		"number":                   func(p, _ map[string]any) { p["number"] = 32 },
		"closed":                   func(p, _ map[string]any) { p["state"] = "closed" },
		"draft":                    func(p, _ map[string]any) { p["draft"] = true },
		"stack":                    func(p, _ map[string]any) { p["stack"] = map[string]any{"number": 1, "size": 2} },
		"checks-fail":              func(p, _ map[string]any) { p["mergeable_state"] = "unstable" },
		"review-or-queue-block":    func(p, _ map[string]any) { p["mergeable_state"] = "blocked" },
		"mergeability-unknown":     func(p, _ map[string]any) { p["mergeable"] = nil },
		"missing-write-permission": func(_, s map[string]any) { s["permissions"] = map[string]any{"push": false, "admin": true} },
		"method-disallowed":        func(_, s map[string]any) { s["allow_squash_merge"] = false },
	} {
		t.Run(name, func(t *testing.T) {
			target := mergeTargetFixture()
			p, settings := mergePRFixture(target), mergeSettingsFixture(target)
			change(p, settings)
			puts := 0
			a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				method, path := githubCall(t, cwd, program, args)
				if method != "GET" {
					puts++
					return nil, errors.New("forbidden mutation")
				}
				if path == prEndpoint(target) {
					return mergeHTTPFixture(200, p), nil
				}
				return mergeHTTPFixture(200, settings), nil
			}}
			o, _ := a.Merge(target)
			if o.State != "blocked" || puts != 0 {
				t.Fatalf("accepted %s: %+v puts=%d", name, o, puts)
			}
		})
	}
}

func TestGitHubAlreadyMergedRequiresExactHeadAndActualReceipt(t *testing.T) {
	for _, mode := range []string{"valid", "different-head", "missing-merge-oid", "missing-time", "open-state"} {
		t.Run(mode, func(t *testing.T) {
			target := mergeTargetFixture()
			p := mergePRFixture(target)
			p["state"], p["merged"], p["merge_commit_sha"], p["merged_at"] = "closed", true, strings.Repeat("c", 40), "2026-10-08T12:00:00Z"
			switch mode {
			case "different-head":
				p["head"].(map[string]any)["sha"] = strings.Repeat("d", 40)
			case "missing-merge-oid":
				p["merge_commit_sha"] = nil
			case "missing-time":
				p["merged_at"] = nil
			case "open-state":
				p["state"] = "open"
			}
			calls := 0
			a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				m, path := githubCall(t, cwd, program, args)
				calls++
				if m != "GET" || path != prEndpoint(target) {
					t.Fatalf("externally merged PR triggered effect: %v", args)
				}
				return mergeHTTPFixture(200, p), nil
			}}
			o, err := a.Merge(target)
			if mode == "valid" {
				if err != nil || o.State != "merged" || o.MergeOID == target.HeadOID {
					t.Fatalf("missing external merge evidence: %+v %v", o, err)
				}
			} else if o.State == "merged" {
				t.Fatalf("invalid merged receipt accepted: %+v", o)
			}
			if calls != 1 {
				t.Fatalf("unnecessary GitHub call after final observation: %d", calls)
			}
		})
	}
}

func TestGitHubPendingOperationResumesWithoutAnotherMerge(t *testing.T) {
	target := mergeTargetFixture()
	puts := 0
	a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
		m, path := githubCall(t, cwd, program, args)
		if m == "PUT" {
			puts++
			return mergeHTTPFixture(202, mergePendingFixture(target)), nil
		}
		if path == prEndpoint(target) {
			return mergeHTTPFixture(200, mergePRFixture(target)), nil
		}
		if path == "repos/"+target.Repository {
			return mergeHTTPFixture(200, mergeSettingsFixture(target)), nil
		}
		return mergeHTTPFixture(200, mergePendingFixture(target)), nil
	}}
	o, err := a.Merge(target)
	if err != nil || o.State != "remote_pending" || o.OperationID == "" || o.MergeOID != "" {
		t.Fatalf("asynchronous acceptance treated as integration: %+v %v", o, err)
	}
	target.OperationID = o.OperationID
	for i := 0; i < 2; i++ {
		o, err = a.Merge(target)
		if err != nil || o.State != "remote_pending" {
			t.Fatalf("resume: %+v %v", o, err)
		}
	}
	if puts != 1 {
		t.Fatalf("pending merge repeated %d requests", puts)
	}
}

func TestGitHubUnknownAndEnqueuedAreNeverMergeReceipts(t *testing.T) {
	for _, mode := range []string{"timeout", "malformed", "missing-operation-id", "unexpected-head", "bypass", "enqueued", "expired-operation"} {
		t.Run(mode, func(t *testing.T) {
			target := mergeTargetFixture()
			puts := 0
			a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				m, path := githubCall(t, cwd, program, args)
				if m == "GET" && path == prEndpoint(target) {
					return mergeHTTPFixture(200, mergePRFixture(target)), nil
				}
				if m == "GET" && path == "repos/"+target.Repository {
					return mergeHTTPFixture(200, mergeSettingsFixture(target)), nil
				}
				if m == "GET" {
					return mergeHTTPFixture(404, map[string]any{"message": "expired"}), errors.New("gh exit 1")
				}
				puts++
				w := mergePendingFixture(target)
				switch mode {
				case "timeout":
					return nil, errors.New("timeout after dispatch")
				case "malformed":
					return []byte("unexpected response"), nil
				case "missing-operation-id":
					delete(w["details"].(map[string]any), "uuid")
				case "unexpected-head":
					w["details"].(map[string]any)["expected_head_sha"] = strings.Repeat("f", 40)
				case "bypass":
					w["details"].(map[string]any)["bypass_rules"] = true
				case "enqueued":
					return mergeHTTPFixture(200, map[string]any{"status": "enqueued", "details": map[string]any{"message": "queued"}}), nil
				}
				return mergeHTTPFixture(202, w), nil
			}}
			o, _ := a.Merge(target)
			want := "effect_unknown"
			if mode == "enqueued" {
				want = "remote_pending"
			}
			if o.State != want || o.MergeOID != "" || puts != 1 {
				t.Fatalf("ambiguous acceptance incorrectly qualified: %+v puts=%d", o, puts)
			}
			if mode == "expired-operation" && o.OperationID == "" {
				t.Fatal("lost operation ID after observation failure")
			}
		})
	}
}

func TestGitHubMethodMustBeExplicitAndBaseOIDIsOnlyAnObservation(t *testing.T) {
	for _, method := range []string{"", "rebase", "auto", "squash --admin"} {
		target := mergeTargetFixture()
		target.Method = method
		a := GitHubMergeAdapter{Run: func(string, string, ...string) ([]byte, error) {
			t.Fatal("invalid method performed network lookup")
			return nil, nil
		}}
		if _, err := a.Merge(target); err == nil {
			t.Fatalf("accepted method %q", method)
		}
	}
	target := mergeTargetFixture()
	a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
		_, path := githubCall(t, cwd, program, args)
		if path == prEndpoint(target) {
			p := mergePRFixture(target)
			p["base"].(map[string]any)["sha"] = strings.Repeat("f", 40)
			return mergeHTTPFixture(200, p), nil
		}
		return mergeHTTPFixture(200, mergeSettingsFixture(target)), nil
	}}
	o, err := a.Observe(target)
	if err != nil || o.State != "ready" || o.ObservedBaseOID != strings.Repeat("f", 40) {
		t.Fatalf("ordinary base progress rejected or stale tested-OID asserted: %+v %v", o, err)
	}
}

func TestGitHubProvenNoEffectIsDistinctFromUnknown(t *testing.T) {
	for _, mode := range []string{"server-rejected", "operation-failed", "operation-mismatch", "connection-lost", "preflight-blocked"} {
		t.Run(mode, func(t *testing.T) {
			target := mergeTargetFixture()
			puts := 0
			a := GitHubMergeAdapter{Run: func(cwd, program string, args ...string) ([]byte, error) {
				method, path := githubCall(t, cwd, program, args)
				if method == "GET" && path == prEndpoint(target) {
					p := mergePRFixture(target)
					if mode == "preflight-blocked" {
						p["mergeable_state"] = "blocked"
					}
					return mergeHTTPFixture(200, p), nil
				}
				if method == "GET" && path == "repos/"+target.Repository {
					return mergeHTTPFixture(200, mergeSettingsFixture(target)), nil
				}
				if method == "GET" {
					return mergeHTTPFixture(200, map[string]any{"status": "failed", "details": map[string]any{"message": "Required checks did not pass."}}), nil
				}
				puts++
				if mode == "server-rejected" {
					return mergeHTTPFixture(403, map[string]any{"message": "Forbidden"}), errors.New("gh exit 1")
				}
				if mode == "connection-lost" {
					return nil, errors.New("timeout")
				}
				pending := mergePendingFixture(target)
				if mode == "operation-mismatch" {
					pending["details"].(map[string]any)["expected_head_sha"] = strings.Repeat("f", 40)
				}
				return mergeHTTPFixture(202, pending), nil
			}}
			o, _ := a.Merge(target)
			want := mode == "server-rejected" || mode == "operation-failed" || mode == "preflight-blocked"
			if o.NoEffect != want {
				t.Fatalf("proven no-effect=%t want=%t: %+v", o.NoEffect, want, o)
			}
			if puts > 1 {
				t.Fatal("adapter automatically retried rejected merge")
			}
			if mode == "operation-failed" && o.OperationID == "" {
				t.Fatal("terminal remote rejection lost operation identity")
			}
		})
	}
}
