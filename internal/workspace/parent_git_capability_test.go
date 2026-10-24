package workspace

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type parentFormatCase struct {
	name, output, code, state string
	err                       bool
}

var parentFormatCases = []parentFormatCase{
	{"echo", "OPTION\n", "task_queue_git_capability_unsupported", "blocked", false},
	{"newline", "\n", "task_queue_git_capability_unsupported", "blocked", false},
	{"unknown", "future-format\n", "task_queue_git_capability_unsupported", "blocked", false},
	{"reftable", "reftable\n", "task_queue_git_capability_unsupported", "blocked", false},
	{"quoted", "bad\t\"format\n", "task_queue_git_capability_unsupported", "blocked", false},
	{"empty", "", "task_queue_observation_unknown", "unknown", false},
	{"multiline", "files\nfiles\n", "task_queue_observation_unknown", "unknown", false},
	{"command-error", "", "task_queue_observation_unknown", "unknown", true},
}

// Install only after fixture creation. All Git calls made by the consuming path
// are counted and restricted to reads in this disposable repository. Each real
// child is pinned and bounded; no provider or unsupported Git is invoked.
func parentReadRunner(t *testing.T, f workItemJourneyFixture, counts map[string]int) workGitRunner {
	t.Helper()
	return func(args, env []string) GitCommandOutcome {
		original := append([]string(nil), args...)
		for len(args) >= 2 && args[0] == "-c" {
			args = args[2:]
		}
		if len(args) > 0 && args[0] == "check-ref-format" {
			// Pure syntax validation has no repository argument.
		} else {
			if len(args) < 3 || args[0] != "-C" || !strings.HasPrefix(args[1], f.wrapper+string(filepath.Separator)) {
				t.Fatalf("unexpected Git location/command: %q", original)
			}
			args = args[2:]
		}
		allowed := args[0] == "rev-parse" || args[0] == "symbolic-ref" && reflect.DeepEqual(args[1:], []string{"-q", "HEAD"}) ||
			args[0] == "merge-base" && len(args) == 4 && args[1] == "--is-ancestor" ||
			args[0] == "diff-tree" && len(args) == 7 && reflect.DeepEqual(args[1:5], []string{"-r", "-z", "--no-commit-id", "--name-only"}) ||
			args[0] == "status" || args[0] == "show-ref" || args[0] == "check-ref-format" ||
			args[0] == "worktree" && len(args) > 1 && args[1] == "list"
		if !allowed {
			t.Fatalf("unexpected Git effect/command: %q", original)
		}
		counts[args[0]]++
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "/opt/homebrew/bin/git", original...)
		cmd.Env = safeIntegrationEnvironment(env)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		exit := 0
		if err != nil {
			exit = -1
			var ee *exec.ExitError
			if errors.As(err, &ee) {
				exit = ee.ExitCode()
			}
		}
		return GitCommandOutcome{Exit: exit, Err: err, Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	}
}

func parentFormatSeam(t *testing.T, f *workItemJourneyFixture, option string, tc parentFormatCase) *int {
	t.Helper()
	counts := map[string]int{}
	runner := parentReadRunner(t, *f, counts)
	f.dependencies.WorkGit.(*systemWorkItemGit).run = runner
	hits := 0
	f.dependencies.IntegrationGit.(*systemTaskIntegrationGit).run = func(args, env []string) GitCommandOutcome {
		if reflect.DeepEqual(args, append(integrationGitPrefix(filepath.Join(f.wrapper, "epic")), "rev-parse", option)) {
			hits++
			if tc.err {
				return GitCommandOutcome{Exit: 128, Err: errors.New("controlled Git failure"), Stderr: []byte("controlled parent format read failure")}
			}
			return GitCommandOutcome{Stdout: []byte(strings.ReplaceAll(tc.output, "OPTION", option))}
		}
		return runner(args, env)
	}
	t.Cleanup(func() { t.Logf("T0523E9D1 option=%s intercepted=%d allowlisted reads=%v", option, hits, counts) })
	return &hits
}

func parentCapabilityFixture(t *testing.T, goal bool) workItemJourneyFixture {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	if goal {
		g := goalFixture(t, f, "task", "goal")
		goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &g}})
	} else {
		s := queueSelect(t, f, "task", "selected", nil)
		queueSetFixture(t, f, "queue", []QueueEntry{{TaskID: "task", Selection: &s}})
	}
	return f
}

func parentCapabilityState(t *testing.T, f workItemJourneyFixture) string {
	t.Helper()
	// Include empty directories as well as the existing registry/content/ref and
	// worktree-inventory digest: even an empty execution directory is an effect.
	dirs := []string{}
	root := filepath.Dir(f.wrapper)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			dirs = append(dirs, path)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return queueDigest([]any{queueStateBytes(t, f), dirs})
}

func TestParentGitCapabilityPreparationPreview(t *testing.T) {
	f := parentCapabilityFixture(t, false)
	before := parentCapabilityState(t, f)
	for _, option := range []string{"--show-object-format=storage", "--show-ref-format"} {
		for _, tc := range parentFormatCases {
			for _, selector := range []string{"task", "next"} {
				t.Run(option+"/"+tc.name+"/"+selector, func(t *testing.T) {
					hits := parentFormatSeam(t, &f, option, tc)
					id := TaskID("task")
					sel := QueueSelector{Kind: selector}
					if selector == "task" {
						sel.TaskID = &id
					}
					out, err := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: sel})
					if *hits == 0 {
						t.Fatal("fixture did not reach parent format command")
					}
					if err != nil {
						t.Fatalf("preview error: %v", err)
					}
					if out.Confirmation != nil || out.Plan != nil {
						t.Fatalf("unsupported parent issued confirmation: state=%s confirmation=%v", out.State, out.Confirmation)
					}
					if out.State != tc.state || out.Freshness != "unknown" || len(out.Reasons) != 1 || out.Reasons[0].Code != tc.code {
						t.Fatalf("lost specific readback: state=%s freshness=%s reasons=%+v", out.State, out.Freshness, out.Reasons)
					}
					assertParentFormatDetail(t, option, tc, out.Reasons[0].Message)
					if parentCapabilityState(t, f) != before {
						t.Fatal("preview changed registry/content/refs/inventory")
					}
				})
			}
		}
	}
}

func assertParentFormatDetail(t *testing.T, option string, tc parentFormatCase, detail string) {
	t.Helper()
	if tc.code == "task_queue_git_capability_unsupported" {
		value := strings.TrimSuffix(strings.ReplaceAll(tc.output, "OPTION", option), "\n")
		allowed := `"files"`
		if option == "--show-object-format=storage" {
			allowed = `"sha1" or "sha256"`
		}
		if !strings.Contains(detail, option) || !strings.Contains(detail, strconv.Quote(value)) || !strings.Contains(detail, allowed) {
			t.Fatalf("missing safe observed/required format detail: %q", detail)
		}
	} else if tc.err && !strings.Contains(detail, "controlled parent format read failure") {
		t.Fatalf("lost command error: %q", detail)
	} else if !tc.err && !strings.Contains(detail, "malformed output") {
		t.Fatalf("lost parser error: %q", detail)
	}
}

func TestParentGitCapabilityGoalPreview(t *testing.T) {
	f := parentCapabilityFixture(t, true)
	before := parentCapabilityState(t, f)
	for _, option := range []string{"--show-object-format=storage", "--show-ref-format"} {
		for _, tc := range parentFormatCases {
			t.Run(option+"/"+tc.name, func(t *testing.T) {
				hits := parentFormatSeam(t, &f, option, tc)
				out, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
				if *hits == 0 {
					t.Fatal("fixture did not reach parent format command")
				}
				if err == nil || out.Confirmation != "" {
					t.Fatal("unsupported parent issued goal confirmation")
				}
				reason := queueErrorReasons(err)[0]
				if reason.Code != tc.code {
					t.Fatalf("lost goal reason: %+v", reason)
				}
				assertParentFormatDetail(t, option, tc, reason.Message)
				if parentCapabilityState(t, f) != before {
					t.Fatal("goal preview changed registry/content/refs/inventory")
				}
			})
		}
	}
}

func TestParentGitCapabilityChangedBeforeApply(t *testing.T) {
	for _, goal := range []bool{false, true} {
		for _, option := range []string{"--show-object-format=storage", "--show-ref-format"} {
			t.Run(fmt.Sprintf("goal=%t/%s", goal, option), func(t *testing.T) {
				f := parentCapabilityFixture(t, goal)
				before := parentCapabilityState(t, f)
				if goal {
					p, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
					if err != nil || p.Confirmation == "" {
						t.Fatalf("supported preview: %v", err)
					}
					hits := parentFormatSeam(t, &f, option, parentFormatCases[0])
					if _, err = PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
						t.Fatal("changed format applied")
					}
					if *hits == 0 {
						t.Fatal("apply did not reobserve format")
					}
				} else {
					in := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}}
					p, err := PrepareTask(f.dependencies, in)
					if err != nil || p.Confirmation == nil {
						t.Fatalf("supported preview: %v", err)
					}
					hits := parentFormatSeam(t, &f, option, parentFormatCases[0])
					in.Apply, in.Confirmation = true, *p.Confirmation
					out, err := PrepareTask(f.dependencies, in)
					if err == nil || out.Preparation != nil || out.Confirmation != nil || !strings.Contains(err.Error(), "task_prepare_conflict") {
						t.Fatalf("changed format apply: %+v %v", out, err)
					}
					if *hits == 0 {
						t.Fatal("apply did not reobserve format")
					}
				}
				if parentCapabilityState(t, f) != before {
					t.Fatal("refused apply changed registry/content/refs/inventory")
				}
			})
		}
	}
}

func TestParentGitCapabilitySupportedFormats(t *testing.T) {
	for _, object := range []string{"sha1", "sha256"} {
		t.Run(object, func(t *testing.T) {
			t.Setenv("GIT_DEFAULT_HASH", object)
			for _, goal := range []bool{false, true} {
				f := parentCapabilityFixture(t, goal)
				before := parentCapabilityState(t, f)
				hits := parentFormatSeam(t, &f, "--show-object-format=storage", parentFormatCase{output: object + "\n"})
				var oid, tree string
				if goal {
					out, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{Target: queueTargetFixture(), Next: true})
					if err != nil || out.Confirmation == "" {
						t.Fatalf("supported goal: %v", err)
					}
					oid, tree = out.ParentOID, out.ParentTree
				} else {
					out, err := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}})
					if err != nil || out.Confirmation == nil {
						t.Fatalf("supported preparation: %+v %v", out, err)
					}
					oid, tree = out.Plan.ParentOID, out.Plan.ParentTree
				}
				if *hits == 0 || !objectIDsMatchFormat(object, oid, tree) {
					t.Fatalf("format/OID mismatch: %s %s %s", object, oid, tree)
				}
				if parentCapabilityState(t, f) != before {
					t.Fatal("supported preview wrote state")
				}
			}
		})
	}
}
