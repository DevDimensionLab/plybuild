//go:build darwin

package taskrun_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/weekrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Only initial registration transport is substituted: the test creates its
// actual Goal/preparation/WorkflowStart while the parent retains its dispatch.
// All subsequent queue, child/control observations and controls use NativeAdapter.
type sl01NativeParent struct {
	weekrun.NativeAdapter
	id     string
	starts int
}

func (*sl01NativeParent) Preview(weekrun.Contract, weekrun.Selection, taskrun.Executable) error {
	return nil
}
func (n *sl01NativeParent) Start(c weekrun.Contract, _ weekrun.Selection, _ taskrun.Executable, check bool) (weekrun.Run, error) {
	if !check || n.id == "" {
		n.starts++
		return weekrun.Run{}, errors.New("fixture forbids another native launch")
	}
	return n.NativeAdapter.Observe(c, n.id)
}
func sl01ReadWeek(t *testing.T, root, id string) (weekrun.Binding, string, weekrun.State) {
	t.Helper()
	dir := filepath.Join(root, ".ply", "week-runs", id)
	var b weekrun.Binding
	var s weekrun.State
	for name, v := range map[string]any{"binding.json": &b, "state.json": &s} {
		raw, e := os.ReadFile(filepath.Join(dir, name))
		if e != nil {
			t.Fatal(e)
		}
		if e = json.Unmarshal(raw, v); e != nil {
			t.Fatal(e)
		}
	}
	return b, dir, s
}
func TestSL01NativeParentQuestionAnswerAndCompletedPauseResume(t *testing.T) {
	sl01NativeAnswerJourney(t, 3, 2, false)
}

func TestSL01NativeSingleAttemptZeroRetryAnswer(t *testing.T) {
	sl01NativeAnswerJourney(t, 1, 0, true)
}

func TestAR01ConfirmedAnswerKeepsReviewerProviderWait(t *testing.T) {
	sl01NativeAnswerJourney(t, 3, 2, false, true)
}

func sl01NativeAnswerJourney(t *testing.T, attempts, retries int, capture bool, reviewerWait ...bool) {
	var m weekrun.Manager
	var n *sl01NativeParent
	var b weekrun.Binding
	var dir string
	var state weekrun.State
	worker := &taskrun.WorkerContract{SchemaVersion: 2, MaxAttempts: attempts, AgentTimeoutSeconds: 30, MaxProviderRetries: retries, ProductQuestions: &taskrun.ProductQuestionContract{SchemaVersion: 1, Mode: taskrun.BoundAnswerMode, MaxContinuationTurns: 1}}
	f := taskrun.NewSL01NativeFixture(t, func(d taskrun.Dependencies, root, epic string, goal workspace.TaskGoalRef) *taskrun.ControlContract {
		n = &sl01NativeParent{NativeAdapter: weekrun.NativeAdapter{Dependencies: d}}
		control := filepath.Join(root, "parent-control-standin")
		if e := os.WriteFile(control, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
			t.Fatal(e)
		}
		m = weekrun.Manager{Native: n, Now: time.Now, Executable: func() (string, error) { return control, nil }, Launch: func(weekrun.Binding) error { return nil }}
		c := weekrun.Contract{Kind: "ply.workflow.week-contract", SchemaVersion: 2, Name: "sl01-native", Workspace: root, ReturnWorktree: epic, ReturnRef: "refs/heads/epic", ProjectID: "ply", RepoID: "ply", EpicID: "epic", Selected: []weekrun.Selection{{TaskID: "task", Goal: goal, DependsOn: []workspace.TaskID{}}}, Model: "fixture-model", Effort: "medium", DurationSeconds: 3600, MaxAttempts: attempts, AgentTimeoutSeconds: 30, MaxProviderRetries: retries, RuntimeControl: &weekrun.RuntimeControl{SchemaVersion: 1, Mode: taskrun.BoundControlMode}, ProductQuestions: &taskrun.ProductQuestionContract{SchemaVersion: 1, Mode: taskrun.BoundAnswerMode, MaxContinuationTurns: 1}}
		var e error
		state, e = m.Start(c, false)
		if e != nil {
			t.Fatal(e)
		}
		b, dir, state = sl01ReadWeek(t, root, state.ID)
		barrier := errors.New("owned fixture constructs native child after dispatch reservation")
		m.Fault = func(point string) error {
			if point == "week_after_dispatch_reservation" {
				return barrier
			}
			return nil
		}
		if e = m.Tick(b, dir, &state); !errors.Is(e, barrier) {
			t.Fatalf("reservation barrier: %v", e)
		}
		m.Fault = nil
		if capture {
			// Exercise the actual NativeAdapter/Execute path from the frozen
			// contract. Stop at its existing intent-before-preparation barrier;
			// this fixture then performs its single native Task preparation.
			bin := filepath.Join(root, "preview-bin")
			if e = os.Mkdir(bin, 0700); e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(filepath.Join(bin, "codex"), []byte("#!/bin/sh\nexit 99\n"), 0700); e != nil {
				t.Fatal(e)
			}
			t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
			intentBarrier := errors.New("inspect the native preserved worker before preparation")
			d.Fault = func(point string) error {
				if point == "execute_before_preparation" {
					return intentBarrier
				}
				return nil
			}
			controlRaw, er := os.ReadFile(control)
			if er != nil {
				t.Fatal(er)
			}
			adapter := weekrun.NativeAdapter{Dependencies: d}
			_, e = adapter.Start(b.Contract, b.Contract.Selected[0], taskrun.Executable{Path: control, SHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(controlRaw))}, false)
			if !errors.Is(e, intentBarrier) {
				t.Fatalf("actual NativeAdapter did not reach preserved worker intent: %v", e)
			}
			found := 0
			e = filepath.WalkDir(filepath.Join(root, ".ply"), func(path string, entry fs.DirEntry, er error) error {
				if er != nil || entry.IsDir() || entry.Name() != "intent.json" {
					return er
				}
				raw, er := os.ReadFile(path)
				if er != nil {
					return er
				}
				var intent struct {
					Kind    string                   `json:"kind"`
					Worker  *taskrun.WorkerContract  `json:"worker"`
					Control *taskrun.ControlContract `json:"control_contract"`
				}
				if json.Unmarshal(raw, &intent) != nil || intent.Kind != "ply.workflow.execute-intent" {
					return nil
				}
				if !reflect.DeepEqual(intent.Worker, worker) || !reflect.DeepEqual(intent.Control, b.ControlContract) {
					return fmt.Errorf("native intent differs from frozen week: worker=%+v control=%+v", intent.Worker, intent.Control)
				}
				*worker = *intent.Worker
				found++
				t.Logf("actual NativeAdapter preserved attempts=%d/retries=%d/continuations=%d, intent sha256:%x", worker.MaxAttempts, worker.MaxProviderRetries, worker.ProductQuestions.MaxContinuationTurns, sha256.Sum256(raw))
				return nil
			})
			if e != nil || found != 1 {
				t.Fatalf("expected one actual NativeAdapter worker intent, found=%d err=%v", found, e)
			}
		}
		return b.ControlContract
	}, worker)
	export := func(name string, value any) {
		if capture {
			f.Export(t, name, value)
		}
	}
	n.Dependencies = f.D
	n.id = f.RunID
	if e := m.Tick(b, dir, &state); e != nil {
		t.Fatal(e)
	}
	f.Step(t)      // native developer stage
	f.Step(t)      // actual inherited-owner attempt -> owned provider -> result/streams
	o := f.Step(t) // publish the actual structured question
	if o.Worker.Phase != "awaiting_product_answer" {
		t.Fatal(o.Worker)
	}
	f.AssertNativeSource(t)
	if e := m.Tick(b, dir, &state); e != nil {
		t.Fatal(e)
	}
	// NA01: each public worker/week read agrees on phase, entry and next actor.
	// With PLY_BIN, the same reads also go through that public CLI executable.
	cli := func(t *testing.T, v any, args ...string) bool {
		t.Helper()
		bin := os.Getenv("PLY_BIN")
		if bin == "" {
			return false
		}
		cmd := exec.Command(bin, append(args, "--workspace", f.Root, "--format", "json")...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		raw, e := cmd.Output()
		if e != nil || json.Unmarshal(raw, v) != nil {
			t.Fatalf("public CLI %v: %v\n%s\n%s", args, e, raw, stderr.Bytes())
		}
		t.Logf("na01 public CLI %s %v", bin, args)
		return true
	}
	next := func(t *testing.T, actor string, want ...string) {
		t.Helper()
		o, e := taskrun.ShowWorker(f.D, f.Root, f.RunID)
		if e != nil || o.NextAction.Actor != actor {
			t.Fatalf("worker next action %+v (phase %s stage %s), want actor %s: %v", o.NextAction, o.Worker.Phase, o.Worker.Stage, actor, e)
		}
		for _, w := range want {
			if !strings.Contains(o.NextAction.Message, w) {
				t.Fatalf("worker next action %q lacks %q", o.NextAction.Message, w)
			}
		}
		var public taskrun.WorkflowRun
		if cli(t, &public, "workflow", "worker", "show", f.RunID) && (public.NextAction != o.NextAction || public.Worker.Phase != o.Worker.Phase || public.Worker.Stage != o.Worker.Stage) {
			t.Fatalf("public worker show %+v %s/%s differs from %+v %s/%s", public.NextAction, public.Worker.Phase, public.Worker.Stage, o.NextAction, o.Worker.Phase, o.Worker.Stage)
		}
	}
	week := func(t *testing.T, phase, entry string, asksAnswer bool) {
		t.Helper()
		if e := m.Tick(b, dir, &state); e != nil {
			t.Fatal(e)
		}
		shown, e := m.Show(f.Root, b.ID)
		views := []weekrun.State{state, shown}
		var public weekrun.State
		if cli(t, &public, "workflow", "week", "show", b.ID) {
			views = append(views, public)
		}
		for _, s := range views {
			if e != nil || s.Phase != phase || s.Entries[0].State != entry || strings.Contains(s.NextAction, "product answer") != asksAnswer || s.NextAction != state.NextAction {
				t.Fatalf("week %s/%s %q, want %s/%s asks=%v as %q: %v", s.Phase, s.Entries[0].State, s.NextAction, phase, entry, asksAnswer, state.NextAction, e)
			}
		}
	}
	// Question and answer reads keep their own public meaning at each stage.
	asked := func(t *testing.T, question, answer, request string) {
		t.Helper()
		q, e := taskrun.ReadProductQuestion(f.D, f.Root, f.RunID)
		if e != nil || q.QuestionState != question {
			t.Fatalf("question state %s, want %s: %v", q.QuestionState, question, e)
		}
		var publicQ taskrun.ProductQuestionReadback
		if cli(t, &publicQ, "workflow", "execute", "question", f.RunID) && publicQ.QuestionState != question {
			t.Fatalf("public question state %s, want %s", publicQ.QuestionState, question)
		}
		if request == "" {
			return
		}
		a, e := taskrun.ShowBoundAnswer(f.D, f.Root, request)
		if e != nil || a.State != answer {
			t.Fatalf("answer state %s, want %s: %v", a.State, answer, e)
		}
		var publicA taskrun.BoundAnswerStatus
		if cli(t, &publicA, "workflow", "execute", "answer", "status", request) && publicA.State != answer {
			t.Fatalf("public answer state %s, want %s", publicA.State, answer)
		}
	}
	next(t, "product_answer", "product choice")
	week(t, "awaiting_product_answer", "awaiting_product_answer", true)
	asked(t, "active", "", "")
	q, e := taskrun.ReadProductQuestion(f.D, f.Root, f.RunID)
	if e != nil || q.Eligibility != "ready" {
		t.Fatalf("read %+v %v", q, e)
	}
	if q.Scope.GoalRevision != b.Contract.Selected[0].Goal.Spec.Revision || q.Scope.GoalSHA256 != b.Contract.Selected[0].Goal.Spec.ManifestSHA256 || q.ParentWeek.WorksetSHA256 != b.SHA256 || q.ParentWeek.ContractSHA256 != b.ControlContract.Digest {
		t.Fatalf("public scope does not match the actual native Goal and parent contract: scope=%+v selected=%+v parent=%+v binding=%s control=%+v", q.Scope, b.Contract.Selected[0], q.ParentWeek, b.SHA256, b.ControlContract)
	}
	export("native-question-ready", q)
	before := *q.Worker
	// Two actual controller ownership lifetimes reload their immutable binding,
	// generation and gate. The held child retains its same reservation/session.
	for i := 0; i < 2; i++ {
		generation := state.Generation
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		e = m.Serve(ctx, f.Root, b.ID)
		cancel()
		if e != nil && !errors.Is(e, context.DeadlineExceeded) {
			t.Fatal(e)
		}
		b, dir, state = sl01ReadWeek(t, f.Root, b.ID)
		if state.Generation != generation+1 || state.Phase != "awaiting_product_answer" || state.ActiveTaskID != "task" || state.Entries[0].RunID != f.RunID || state.Deadline != b.Deadline {
			t.Fatal(state)
		}
	}
	q, e = taskrun.ReadProductQuestion(f.D, f.Root, f.RunID)
	if e != nil || q.Eligibility != "ready" || !reflect.DeepEqual(before, *q.Worker) || f.TurnCount(t) != 1 || n.starts != 0 {
		t.Fatalf("parent restart changed child %+v %v", q, e)
	}
	p, e := taskrun.PreviewBoundAnswer(f.D, f.Root, f.RunID, q.Question.QuestionID, 1, "task_first", "native-answer-0001")
	if e != nil || p.Request == nil {
		t.Fatalf("preview %+v %v", p, e)
	}
	export("native-preview-ready", p)
	export("native-request", *p.Request)
	t.Run("source_session_drift_blocks_apply", func(t *testing.T) { f.AssertSourceSessionDriftBlocks(t, *p.Request, "apply") })
	status, e := taskrun.ApplyBoundAnswer(f.D, *p.Request)
	if e != nil || status.State != "answer_recorded" {
		t.Fatalf("apply %+v %v", status, e)
	}
	export("native-answer-recorded", status)
	t.Run("na01_received_answer_is_not_requested_again", func(t *testing.T) {
		next(t, "worker", p.Request.RequestID, "received", "use is not confirmed")
		week(t, "running", "running", false)
		asked(t, "answered", "answer_recorded", p.Request.RequestID)
	})
	f.Step(t) // immutable continuation reservation
	t.Run("source_session_drift_blocks_start", func(t *testing.T) { f.AssertSourceSessionDriftBlocks(t, *p.Request, "start") })
	release := f.StartHeldContinuation(t) // owned process held after real turn.started
	for until := time.Now().Add(5 * time.Second); time.Now().Before(until); {
		status, e = taskrun.ShowBoundAnswer(f.D, f.Root, p.Request.RequestID)
		if e == nil && status.State == "continuation_observed" {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if e != nil || status.State != "continuation_observed" {
		release()
		t.Fatalf("live owned turn not observed %+v %v", status, e)
	}
	observedStart := *status.Receipt.ContinuationStartedAtUTC
	t.Run("na01_started_continuation_is_not_confirmed_use", func(t *testing.T) {
		next(t, "worker", p.Request.RequestID, "continuing", "use is not confirmed")
		asked(t, "answered", "continuation_observed", p.Request.RequestID)
	})
	export("native-continuation-observed", status)
	// Parent control completes while the provider is still held. Neither the
	// parent nor TaskRun store guard may be held across the provider wait.
	if _, e = m.Control(f.Root, b.ID, "pause"); e != nil {
		t.Fatal(e)
	}
	if e = m.Tick(b, dir, &state); e != nil {
		t.Fatal(e)
	}
	t.Run("observed_status_retains_current_parent_reason", func(t *testing.T) {
		current, e := taskrun.ShowBoundAnswer(f.D, f.Root, p.Request.RequestID)
		found := false
		for _, r := range current.Reasons {
			found = found || r.Code == "bound_answer_parent_paused"
		}
		if e != nil || current.State != "continuation_observed" || !found {
			t.Fatalf("observed effect erased current pause %+v %v", current, e)
		}
	})
	release()
	status, e = taskrun.ShowBoundAnswer(f.D, f.Root, p.Request.RequestID)
	if e != nil || status.State != "continuation_completed" || f.TurnCount(t) != 2 {
		t.Fatalf("completion %+v %v", status, e)
	}
	t.Run("receipt_start_is_stable_and_ordered", func(t *testing.T) {
		if *status.Receipt.ContinuationStartedAtUTC != observedStart {
			t.Errorf("receipt start moved between actual observed/completed: observed=%s completed=%s", observedStart, *status.Receipt.ContinuationStartedAtUTC)
		}
		previous := time.Time{}
		for _, value := range []string{status.Receipt.RecordedAtUTC, *status.Receipt.ContinuationStartedAtUTC, *status.Receipt.ContinuationCompletedAtUTC, status.ObservedAtUTC} {
			at, e := time.Parse(time.RFC3339Nano, value)
			if e != nil || at.Before(previous) {
				t.Fatalf("receipt timestamps are not recorded <= start <= finished <= observed: %+v %v", status, e)
			}
			previous = at
		}
	})
	export("native-continuation-completed", status)
	t.Run("resumed_context_drift_blocks_completion", func(t *testing.T) { f.AssertResumedContextDriftBlocksCompletion(t, p.Request.RequestID) })
	// Reviewed regression: completed continuation is paused before workerStep
	// consumes it. Lawful parent resume must consume that same result once.
	if o = f.Step(t); o.Worker.Phase != "paused" {
		t.Fatal(o.Worker)
	}
	if e = m.Tick(b, dir, &state); e != nil || state.Phase != "paused" {
		t.Fatalf("parent pause %+v %v", state, e)
	}
	t.Run("completed_status_retains_current_parent_reason", func(t *testing.T) {
		current, e := taskrun.ShowBoundAnswer(f.D, f.Root, p.Request.RequestID)
		found := false
		for _, r := range current.Reasons {
			found = found || r.Code == "bound_answer_parent_paused"
		}
		if e != nil || current.State != "continuation_completed" || !found {
			t.Fatalf("completed evidence erased current parent pause: %+v %v", current, e)
		}
	})
	if _, e = m.Control(f.Root, b.ID, "resume"); e != nil {
		t.Fatal(e)
	}
	if e = m.Tick(b, dir, &state); e != nil {
		t.Fatal(e)
	}
	f.Step(t)
	o = f.Step(t)
	if o.Worker.Stage != "commit" || o.Worker.Job != nil || len(o.Worker.DeveloperSessions) != 1 || f.TurnCount(t) != 2 || o.Delivery.Phase == "completed" {
		t.Fatalf("completed continuation was not consumed into ordinary commit exactly once: %+v", o.Worker)
	}
	t.Run("na01_resumed_commit_shows_confirmed_use", func(t *testing.T) {
		next(t, "worker", p.Request.RequestID, "use is confirmed", "commit")
		week(t, "running", "running", false)
		asked(t, "answered", "continuation_completed", p.Request.RequestID)
	})
	status, e = taskrun.ApplyBoundAnswer(f.D, *p.Request)
	if e != nil || status.State != "continuation_completed" || f.TurnCount(t) != 2 {
		t.Fatalf("replay %+v %v", status, e)
	}
	f.Step(t)     // native commit attempt
	f.Step(t)     // consume native commit
	f.Step(t)     // actual isolated acceptance subprocess
	o = f.Step(t) // accepted pretest -> ordinary independent review
	if o.Worker.Stage != "reviewer" {
		t.Fatal(o.Worker)
	}
	t.Run("na01_review_names_worker_not_human_answer", func(t *testing.T) {
		next(t, "worker", "use is confirmed", "reviewer")
		week(t, "running", "running", false)
		asked(t, "answered", "continuation_completed", p.Request.RequestID)
	})
	t.Run("missing_answer_blocks_reviewer", func(t *testing.T) { f.AssertMissingAnswerBlocksReviewer(t) })
	if len(reviewerWait) != 0 && reviewerWait[0] {
		t.Run("reviewer_usage_limit_after_confirmed_answer", func(t *testing.T) {
			// The owned stand-in exits with a real usage-limit event stream.
			// Its completed attempt reaches the normal controller retry path.
			if e := os.WriteFile(filepath.Join(filepath.Dir(f.Turns), "usage-limit-role"), []byte("reviewer"), 0600); e != nil {
				t.Fatal(e)
			}
			f.Step(t)      // reserve and execute the owned reviewer attempt
			o := f.Step(t) // consume its observed provider failure
			reset := o.Worker.RetryNotBefore
			if o.Worker.Phase != "waiting" || reset == "" || o.Worker.LastProviderFailure == "" {
				t.Fatalf("fixture did not reach the provider wait: %+v", o.Worker)
			}
			if o.NextAction.Actor != "worker_wait" || !strings.Contains(o.NextAction.Message, reset) {
				t.Errorf("ShowWorker hid the current provider wait behind answer history: %+v; reset %s", o.NextAction, reset)
			}
			var public taskrun.WorkflowRun
			if cli(t, &public, "workflow", "worker", "show", f.RunID) && public.NextAction != o.NextAction {
				t.Errorf("public worker show disagrees on the provider wait: %+v / %+v", public.NextAction, o.NextAction)
			}
			shown, e := m.Show(f.Root, b.ID) // no Tick has projected this wait yet
			if e != nil || shown.Phase != "waiting" || !strings.Contains(shown.NextAction, "worker_wait") || !strings.Contains(shown.NextAction, reset) {
				t.Errorf("Week Show hid the current provider wait: %s %q; reset %s: %v", shown.Phase, shown.NextAction, reset, e)
			}
			week(t, "waiting", "running", false)
			if !strings.Contains(state.NextAction, "worker_wait") || !strings.Contains(state.NextAction, reset) {
				t.Errorf("Week Tick hid the current provider wait: %q; reset %s", state.NextAction, reset)
			}
			if _, e := m.Control(f.Root, b.ID, "pause"); e != nil {
				t.Fatal(e)
			}
			if e := m.Tick(b, dir, &state); e != nil {
				t.Fatal(e)
			}
			next(t, "worker", "pause requested", "safe boundary")
		})
		return
	}
	f.Step(t) // actual reserved reviewer intent and owned stand-in
	t.Run("reviewer_intent_binds_selected_answer", func(t *testing.T) { f.AssertReviewerChoice(t) })
	if f.TurnCount(t) != 3 {
		t.Fatalf("expected exactly initial, one answer and independent reviewer; got %d actual launches", f.TurnCount(t))
	}
	if attempts > 1 {
		t.Run("correction_keeps_choice_after_restart", func(t *testing.T) { f.AssertCorrectionRetainsChoiceAcrossOwnerRestart(t) })
	} else {
		t.Run("single_attempt_forbids_correction", func(t *testing.T) { f.AssertSingleAttemptForbidsCorrection(t) })
	}
	if attempts == 1 {
		// NA01 review r001-f01: continue the same public journey through native
		// verify, integration and completion.
		t.Run("na01_delivery_and_completion", func(t *testing.T) {
			if o = f.Step(t); o.Worker.Stage != "verify" || o.Worker.Job != nil {
				t.Fatalf("reviewed candidate did not reach verify: %+v", o.Worker)
			}
			next(t, "worker", p.Request.RequestID, "use is confirmed", "verify")
			week(t, "running", "running", false)
			if o = f.Step(t); o.Delivery.Phase != "automatic_acceptance_passed" {
				t.Fatalf("native verify did not qualify: %+v", o.Delivery)
			}
			next(t, "recipient", "Automatic acceptance and review passed")
			week(t, "running", "running", false)
			f.Step(t) // integrate stage
			// CL02 releases source ownership during integration, then closes
			// from the controller after the native attempt process has exited.
			// Preserve both facts; an integration receipt alone is not closeout.
			if o = f.Step(t); o.Delivery.Phase != "closing" || o.Delivery.OwnershipRelease == nil || o.DeliveryStatus.FinalDelivery.State != "delivered" {
				t.Fatalf("native integration did not reach its released cleanup boundary: %+v", o)
			}
			if o = f.Step(t); o.Delivery.Phase != "completed" {
				t.Fatalf("native integration did not complete: %+v", o.Delivery)
			}
			next(t, "user", "delivery is closed")
			week(t, "completed", "completed", false)
			asked(t, "answered", "continuation_completed", p.Request.RequestID)
			// A completed worker without a closeout record (older or pending
			// closeout) never names the worker as owner of further work.
			for _, phase := range []string{"completed", "running"} {
				a := f.ProjectPersistedAnswerAction(t, phase, "completed")
				if a.Actor == "worker" || a.Actor == "product_answer" || !strings.Contains(a.Message, "use is confirmed") || !strings.Contains(a.Message, "delivery is completed") {
					t.Fatalf("completed worker (%s) projects %+v", phase, a)
				}
			}
		})
	}
}
