package workspace

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestParallelSourceReadOnlyPreflightPreservesUnstarted(t *testing.T) {
	for _, refusal := range []string{"driver", "ignored-path-read"} {
		for _, action := range []string{"retry", "cancel-unstarted"} {
			t.Run(refusal+"/"+action, func(t *testing.T) {
				f, p, in := sourceReconcileFixture(t, false)
				source := p.Preparation.Plan.WorktreePath
				before := runLocalGit(t, source, "rev-parse", "HEAD")
				preview, e := ReconcileTaskSource(f.dependencies, in)
				if e != nil {
					t.Fatal(e)
				}
				if refusal == "driver" {
					runLocalGit(t, source, "config", "--local", "merge.private-fixture.driver", "false")
				}
				git := f.dependencies.IntegrationGit.(*systemTaskIntegrationGit)
				runner, merges, readFailure := git.run, 0, refusal == "ignored-path-read"
				git.run = func(args, env []string) GitCommandOutcome {
					for _, a := range args {
						if a == "merge" {
							merges++
						}
						if a == "--ignored" && readFailure {
							return GitCommandOutcome{Exit: -1, Err: errInjected, Stderr: []byte("simulated read-only ignored-path failure")}
						}
					}
					return runner(args, env)
				}
				in.Apply, in.Confirmation = true, preview.Confirmation
				refused, e := ReconcileTaskSource(f.dependencies, in)
				if e == nil || refused.Operation == nil || refused.State != "intent" {
					t.Fatalf("read-only refusal: %+v %v", refused, e)
				}
				t.Logf("refusal=%s state=%s unstarted=%t merges=%d error=%v", refusal, refused.State, refused.Operation.Unstarted, merges, e)
				if !refused.Operation.Unstarted {
					t.Error("read-only preflight consumed original positive unstarted proof")
				}
				marker := runLocalGit(t, source, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+refused.Operation.ID)
				if _, e = os.Stat(marker); !os.IsNotExist(e) || merges != 0 || runLocalGit(t, source, "rev-parse", "HEAD") != before {
					t.Fatal("read-only refusal attempted a source effect")
				}
				state := queueStateBytes(t, f)
				in.OperationID = refused.Operation.ID
				if _, e = ReconcileTaskSource(f.dependencies, in); e == nil || merges != 0 || queueStateBytes(t, f) != state {
					t.Fatal("unchanged preflight retry changed native intent or attempted a merge")
				}
				if refusal == "driver" {
					runLocalGit(t, source, "config", "--local", "--unset", "merge.private-fixture.driver")
				}
				readFailure = false
				if action == "cancel-unstarted" {
					in.Action = action
				}
				progressed, e := ReconcileTaskSource(f.dependencies, in)
				want, wantMerges := "committed", 1
				if action == "cancel-unstarted" {
					want, wantMerges = "aborted", 0
				}
				if e != nil || progressed.State != want || merges != wantMerges || progressed.Operation.ID != refused.Operation.ID {
					t.Fatalf("same owner %s after preflight: %+v %v merges=%d", action, progressed, e, merges)
				}
			})
		}
	}
}

// Real isolated Git and the unchanged production runner/environment. The
// private driver only writes its own execution flag; no live runtime is used.
func TestParallelSourceRejectsEffectiveCustomDriversBeforeEffect(t *testing.T) {
	for _, scope := range []string{"local-include", "worktree", "worktree-include"} {
		t.Run(scope, func(t *testing.T) {
			f, p, in := sourceReconcileFixture(t, true)
			source := p.Preparation.Plan.WorktreePath
			if e := os.WriteFile(filepath.Join(source, ".gitattributes"), []byte("README.md merge=private-fixture\n"), 0600); e != nil {
				t.Fatal(e)
			}
			runLocalGit(t, source, "add", ".gitattributes")
			runLocalGit(t, source, "commit", "-m", "Private merge attribute")
			script := filepath.Join(t.TempDir(), "private-driver")
			if e := os.WriteFile(script, []byte("#!/bin/sh\n: > \"$0.executed\"\nexit 1\n"), 0700); e != nil {
				t.Fatal(e)
			}
			config := filepath.Join(t.TempDir(), "private.gitconfig")
			runLocalGit(t, source, "config", "--file", config, "merge.private-fixture.driver", strconv.Quote(script))
			if scope == "local-include" {
				runLocalGit(t, source, "config", "--local", "include.path", config)
			} else {
				runLocalGit(t, source, "config", "--local", "extensions.worktreeConfig", "true")
				if scope == "worktree" {
					runLocalGit(t, source, "config", "--worktree", "merge.private-fixture.driver", strconv.Quote(script))
				} else {
					runLocalGit(t, source, "config", "--worktree", "include.path", config)
				}
			}
			before := runLocalGit(t, source, "rev-parse", "HEAD")
			preview, e := ReconcileTaskSource(f.dependencies, in)
			if e != nil {
				t.Fatal(e)
			}
			in.Apply, in.Confirmation = true, preview.Confirmation
			result, e := ReconcileTaskSource(f.dependencies, in)
			_, driverErr := os.Stat(script + ".executed")
			t.Logf("scope=%s state=%s error=%v privateDriverExecuted=%t", scope, result.State, e, driverErr == nil)
			if !os.IsNotExist(driverErr) {
				t.Error("controlled source merge executed a custom driver")
			}
			if e == nil || !strings.Contains(e.Error(), "source_merge_driver_unsupported") || result.State != "intent" || result.Operation == nil {
				t.Fatalf("effective custom driver did not refuse before source effect: %+v %v", result, e)
			}
			marker := runLocalGit(t, source, "rev-parse", "--path-format=absolute", "--git-path", "ply-source-"+result.Operation.ID)
			if _, e = os.Stat(marker); !os.IsNotExist(e) {
				t.Fatal("driver refusal published an operation marker")
			}
			if runLocalGit(t, source, "rev-parse", "HEAD") != before || runLocalGit(t, source, "status", "--porcelain") != "" {
				t.Fatal("driver refusal changed the source or started a merge")
			}
		})
	}
}
