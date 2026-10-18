package taskexecute

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Genuine public CLI preparation and native Git/registry effects. Only the
// provider and Herdr are stand-ins, confined to the disposable fixture.
func TestSharedEpicUniqueSpecPublicCLIPreparesAndPreservesRetry(t *testing.T) {
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		binary = filepath.Join(physicalTemp(t), "ply")
		build := exec.Command("go", "build", "-o", binary, "./cmd/ply")
		build.Dir = source
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build public CLI: %v\n%s", err, output)
		}
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("public CLI binary must be absolute")
	}
	for _, moveTarget := range []bool{false, true} {
		t.Run(map[bool]string{false: "prepare-and-retry", true: "changed-target"}[moveTarget], func(t *testing.T) {
			agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic CLI fixture", TimeoutSeconds: 30}}
			d, in, root, epic := launcherFixtureWithCLI(t, "codex", binary, agreement)
			if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
				t.Fatal(err)
			}
			run := func(executable, cwd string, wantSuccess bool, args ...string) []byte {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
				defer cancel()
				command := exec.CommandContext(ctx, executable, args...)
				command.Dir = cwd
				var stderr strings.Builder
				command.Stderr = &stderr
				output, err := command.Output()
				if (err == nil) != wantSuccess {
					t.Fatalf("CLI %v: %v\nstdout=%s\nstderr=%s", args, err, output, stderr.String())
				}
				if !wantSuccess {
					return []byte(stderr.String())
				}
				return output
			}
			call := func(wantSuccess bool, args ...string) []byte { return run(binary, epic, wantSuccess, args...) }
			base := fixtureGit(t, epic, "rev-parse", "HEAD")
			call(true, "workspace", "epic", "adopt", "sibling", "--title", "Sibling", "--project", "fixture", "--repo", "fixture", "--worktree", epic, "--ref", "refs/heads/epic", "--expected-oid", base)
			sibling := filepath.Join(root, "sibling-task")
			call(true, "workspace", "task", "create", "sibling-task", "--title", "Sibling Task", "--description", "Separate shared-base Task", "--epic", "sibling", "--project", "fixture", "--repo", "fixture")
			call(true, "workspace", "task", "worktree", "create", "sibling-task", "--branch", "sibling-task", "--path", sibling, "--expected-parent-oid", base)
			if output := call(false, "workflow", "execute", "--check", "--format", "json"); !strings.Contains(string(output), "--epic sibling") || !strings.Contains(string(output), "--epic epic") {
				t.Fatalf("ambiguous cwd did not give explicit choices: %s", output)
			}
			var preview Result
			if err := json.Unmarshal(call(true, "workflow", "execute", "--spec", "observable-goal", "--check", "--format", "json"), &preview); err != nil || preview.Goal == nil || preview.Goal.Target.EpicID != "epic" {
				t.Fatalf("Spec-only CLI preview = %+v, %v", preview, err)
			}
			if moveTarget {
				fixtureGit(t, epic, "commit", "--allow-empty", "-m", "shared target advanced")
				before, err := d.Workspace.WorkItems.Snapshot(root)
				if err != nil {
					t.Fatal(err)
				}
				if output := call(false, "workflow", "execute", "--spec", "observable-goal", "--format", "json"); !strings.Contains(string(output), string(workspace.ErrorWorkParentStale)) {
					t.Fatalf("expected changed-parent refusal: %s", output)
				}
				after, err := d.Workspace.WorkItems.Snapshot(root)
				if err != nil || after.RawSHA256 != before.RawSHA256 {
					t.Fatalf("changed target refusal wrote registry: %v", err)
				}
				if _, err := os.Stat(preview.Goal.WorktreePath); !os.IsNotExist(err) {
					t.Fatalf("changed target created Task copy: %v", err)
				}
				return
			}
			var started Result
			if err := json.Unmarshal(call(true, "workflow", "execute", "--spec", "observable-goal", "--format", "json"), &started); err != nil || started.Run == nil {
				t.Fatalf("Spec-only execution did not start: %+v, %v", started, err)
			}
			if fixtureGit(t, preview.Goal.WorktreePath, "symbolic-ref", "HEAD") != "refs/heads/"+preview.Goal.Branch || fixtureGit(t, epic, "rev-parse", "HEAD") != base {
				t.Fatal("CLI preparation failed isolation or changed shared target")
			}
			before, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			var repeated Result
			if err := json.Unmarshal(call(true, "workflow", "execute", "--spec", "observable-goal", "--format", "json"), &repeated); err != nil || repeated.State != "existing" || repeated.Run == nil || repeated.Run.RunID != started.Run.RunID {
				t.Fatalf("CLI retry changed Run: %+v, %v", repeated, err)
			}
			after, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil || after.RawSHA256 != before.RawSHA256 || len(after.TaskPreparations) != 1 || after.TaskPreparations[0].Plan.Target != preview.Goal.Target || len(after.TaskSpecRevisions) != 2 || len(after.TaskSolutionSelections) != 1 {
				t.Fatalf("retry changed selected target/history: %v", err)
			}
			if preview.Goal.Input.Target != (workspace.QueueTargetInput{}) {
				t.Fatal("test did not exercise Spec-only selection")
			}
			var request taskrun.WorkflowRequest
			if _, err := readPreservedJSON(started.RequestPath, &request); err != nil {
				t.Fatal(err)
			}
			taskPath := preview.Goal.WorktreePath
			owner := func(args ...string) []byte { return run(request.Runtime.PlyExecutable.Path, taskPath, true, args...) }
			o := started.Run
			owner("workflow", "execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--format", "json")
			if err := os.WriteFile(filepath.Join(taskPath, "README.md"), []byte("shared Epic candidate\n"), 0644); err != nil {
				t.Fatal(err)
			}
			fixtureGit(t, taskPath, "add", "README.md")
			fixtureGit(t, taskPath, "commit", "-m", "synthetic shared-target implementation")
			candidate := fixtureGit(t, taskPath, "rev-parse", "HEAD")
			review := filepath.Join(root, "synthetic-review.json")
			fixtureJSON(t, review, map[string]any{"kind": "DeliveryCandidateReview@1", "schema_version": 1, "candidate_oid": candidate, "candidate_tree": fixtureGit(t, taskPath, "rev-parse", "HEAD^{tree}"), "reviewer_claim": "synthetic fixture reviewer; no actual provider review", "reviewer_session_id": nil, "decision": "passed", "findings": []any{}, "fixes": []any{}, "open_actionable_findings": []any{}})
			if err := os.WriteFile(preview.Goal.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'shared Epic candidate'\n"), 0600); err != nil {
				t.Fatal(err)
			}
			owner("workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--candidate-binary", binary, "--format", "json")
			if fixtureGit(t, epic, "rev-parse", "HEAD") != base {
				t.Fatal("shared target changed before authorized integration")
			}
			var delivered taskrun.WorkflowRun
			if err := json.Unmarshal(owner("workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context, "--format", "json"), &delivered); err != nil || delivered.DeliveryStatus == nil || delivered.DeliveryStatus.FinalDelivery.State != "delivered" {
				t.Fatalf("native return did not close: %+v, %v", delivered, err)
			}
			if fixtureGit(t, epic, "rev-parse", "HEAD") != candidate || fixtureGit(t, sibling, "rev-parse", "HEAD") != base {
				t.Fatal("integration returned another candidate or changed sibling Task")
			}
			closed, err := workspace.ListTaskQueue(d.Workspace, workspace.QueueTargetInput{ProjectID: "fixture", RepoID: "fixture", EpicID: "epic"}, false)
			if err != nil || closed.Current != nil || len(closed.Pending) != 0 {
				t.Fatalf("native queue not closed: %+v, %v", closed, err)
			}
		})
	}
}
