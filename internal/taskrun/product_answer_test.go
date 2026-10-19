//go:build !windows

package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const sl01Thread = "01a145a7-9221-7ac0-a56d-d64f4c7bab11"

// Isolated source-turn fixture uses physical native store/locks and hash-bound
// result/stream/session files. Its qualification callback labels a synthetic
// stand-in; it does not claim another real Codex probe.
func newSL01Fixture(t *testing.T) *aw01Fixture {
	return sl01Fixture(t, newAW01Fixture(t))
}
func sl01Fixture(t *testing.T, f *aw01Fixture) *aw01Fixture {
	f.job("developer")
	f.s.StartSHA256 = ptr("sha256:" + strings.Repeat("c", 64))
	f.s.Request.Worker = &WorkerContract{SchemaVersion: 2, MaxAttempts: 3, AgentTimeoutSeconds: 120, ProductQuestions: &ProductQuestionContract{1, BoundAnswerMode, 1}}
	f.s.Request.ControlContract = NewControlContract("week", "week_"+strings.Repeat("a", 32), "sha256:"+strings.Repeat("b", 64))
	task := filepath.Join(f.s.Request.WorkspaceRoot, "task")
	if e := privateDir(task); e != nil {
		t.Fatal(e)
	}
	f.file(filepath.Join(task, "guide.md"), []byte("# Guide\n"))
	f.s.Observed.Target.WorktreeLocator = task
	f.d.Workspace.IntegrationGit = aw01Git{x: workspace.IntegrationWorktreeObservation{Locator: task, Ref: f.s.Observed.Target.Ref, GitCommonDir: f.s.Observed.Target.GitCommonDir, OID: f.s.Observed.Target.OID, Tree: f.s.Observed.Target.Tree, Clean: false, InProgress: []string{}}}
	f.s.Request.Delivery.Agreement = &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: "local_epic_integration", ProjectID: "fixture", RepoID: "fixture", EpicID: "fixture"}
	goal := []byte(`{"kind":"WorkspaceTaskSpecRevision@2","schema_version":2,"contract_kind":"goal","task_id":"fixture-task","spec_id":"fixture-goal","revision":1}`)
	f.s.Request.Delivery.Goal = f.file(filepath.Join(f.s.Request.WorkspaceRoot, "sl01-goal.json"), goal)
	execution := []byte(`{"contract_kind":"execution","task_id":"fixture-task","spec_id":"fixture-goal","revision":2,"goal_origin":{"spec_id":"fixture-goal","spec":{"revision":1,"manifest_sha256":"` + hash(goal) + `"}}}`)
	executionBinding := f.file(filepath.Join(f.s.Request.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(hash(execution), "sha256:")+".json"), execution)
	f.s.Request.HandoffDraft = []byte(`{"binding_request":{"project_id":"fixture","repo_id":"fixture"},"task_spec_binding":{"task_id":"fixture-task","spec_id":"fixture-goal","spec":{"revision":2,"manifest_sha256":"` + executionBinding.SHA256 + `"}}}`)
	f.d.ControlQueue = func(string, string) (string, error) { return "sha256:" + strings.Repeat("d", 64), nil }
	f.d.ProductAnswerQualification = func(Runtime, string) (string, error) { return "sha256:" + strings.Repeat("f", 64), nil }
	gate := ParentControlGate{Contract: f.s.Request.ControlContract, Deadline: f.now.Add(time.Hour).Format(time.RFC3339Nano), DesiredAction: "continue", ActiveRun: f.s.Result.RunID, Dispatch: "fixture-dispatch"}
	if e := (Dependencies{}).writeValue(ParentControlPath(f.s.Request.WorkspaceRoot, *f.s.Request.ControlContract.WeekID), gate); e != nil {
		t.Fatal(e)
	}
	home := filepath.Join(f.s.Request.WorkspaceRoot, "codex-home")
	if e := privateDir(home); e != nil {
		t.Fatal(e)
	}
	in, e := workerIntent(f.s, *f.s.Result.Worker.Job)
	if e != nil {
		t.Fatal(e)
	}
	in.Environment = []string{"CODEX_HOME=" + home}
	in.TimeoutSeconds = 120
	in.Argv = []string{f.s.Request.Runtime.Executable.Path, "--no-daemon", "-a", "never", "exec", "--json", "-"}
	raw, _ := Canonical(in)
	if e = publishFile(filepath.Join(f.s.Result.Worker.Job.Directory, "intent.json"), raw, true); e != nil {
		t.Fatal(e)
	}
	f.s.Result.Worker.Job.IntentHash = hash(raw)
	entries := []any{}
	for path, access := range map[string]string{":root": "read", ":minimal": "read", task: "write", filepath.Join(task, ".git"): "read", f.s.Observed.Target.GitCommonDir: "read", filepath.Join(f.s.Result.Worker.Job.Directory, "output"): "write", filepath.Join(f.s.Result.Worker.Job.Directory, "tmp"): "write"} {
		var p any
		if strings.HasPrefix(path, ":") {
			p = map[string]any{"type": "special", "value": map[string]string{"kind": strings.TrimPrefix(path, ":")}}
		} else {
			p = map[string]string{"type": "path", "path": path}
		}
		entries = append(entries, map[string]any{"path": p, "access": access})
	}
	context := map[string]any{"type": "turn_context", "payload": map[string]any{"cwd": task, "approval_policy": "never", "model": f.s.Request.Runtime.Model, "effort": "max", "permission_profile": map[string]any{"type": "managed", "network": "restricted", "file_system": map[string]any{"type": "restricted", "entries": entries}}}}
	metadata, _ := Canonical(map[string]any{"type": "session_meta", "payload": map[string]string{"id": sl01Thread, "cwd": task}})
	contextRaw, _ := Canonical(context)
	f.file(filepath.Join(home, "sessions", "2026", "10", "10", "rollout-"+sl01Thread+".jsonl"), append(append(metadata, '\n'), append(contextRaw, '\n')...))
	f.s.Result.RequestSHA256 = digest(f.s.Request)
	f.save()
	f.own()
	r := f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
	r.ThreadID = sl01Thread
	r.StartedAt = f.now.Add(-time.Second).Format(time.RFC3339Nano)
	stream := []byte(`{"type":"thread.started","thread_id":"` + sl01Thread + `"}` + "\n" + `{"type":"turn.started"}` + "\n" + `{"type":"turn.completed"}` + "\n")
	if e = publishFile(r.Stdout.Locator, stream, true); e != nil {
		t.Fatal(e)
	}
	r.Stdout.SHA256 = hash(stream)
	r.Stderr = r.Stdout
	answer := f.file(filepath.Join(f.s.Result.Worker.Job.Directory, "answer.json"), []byte(`{"status":"needs_product_answer","summary":"A product decision is required","block_reason":"","product_question":{"purpose":"product_choice","text":"Which reading order?","reason":"Both orders satisfy the Goal","options":[{"option_id":"task_first","label":"Task first","consequence":"Tasks precede commands"},{"option_id":"command_first","label":"Commands first","consequence":"Commands precede tasks"}],"continuation":{"type":"resume_same_developer_thread","source_iteration":1,"source_job_id":"developer-0001"}}}`))
	r.Answer = &answer
	if e = (Dependencies{}).replaceValue(filepath.Join(f.s.Result.Worker.Job.Directory, "result.json"), r); e != nil {
		t.Fatal(e)
	}
	return f
}
func sl01Park(t *testing.T, f *aw01Fixture) {
	t.Helper()
	if e := workerStep(f.d, f.s); e != nil {
		t.Fatal(e)
	}
	f.reload()
}
func sl01Preview(t *testing.T, f *aw01Fixture, id string) BoundAnswerPreview {
	t.Helper()
	q, e := ReadProductQuestion(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if q.Eligibility != "ready" {
		t.Fatalf("question not ready: %+v", q)
	}
	p, e := PreviewBoundAnswer(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID, q.Question.QuestionID, 1, "task_first", id)
	if e != nil || p.Request == nil {
		t.Fatalf("preview: %+v %v", p, e)
	}
	return p
}

func TestSL01SourceTurnPublishesDurableQuestion(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	w := f.s.Result.Worker
	if w.Phase != "awaiting_product_answer" || w.Stage != "developer" || w.Iteration != 1 || w.CandidateBinary != "" || f.s.ProductQuestion == nil {
		t.Fatalf("question did not retain its developer turn: %+v", w)
	}
}

// PH01: the default qualification path, without an injected callback, must not
// require the provider executable to match one hardcoded approved Codex digest.
func TestPH01DefaultQualificationAcceptsAnotherProviderDigest(t *testing.T) {
	f := newAW01Fixture(t)
	f.s.Request.Runtime.Model = "gpt-6-astra"
	f = sl01Fixture(t, f)
	f.d.ProductAnswerQualification = nil
	if f.s.Request.Runtime.Executable.SHA256 == "sha256:74c6a6d263f40d65f99c27a815275f5b4ed0bcc2baecc396b4cf3255498bdfc0" {
		t.Fatal("fixture provider must differ from the formerly approved Codex digest")
	}
	sl01Park(t, f)
	if f.s.Result.Worker.Phase != "awaiting_product_answer" || f.s.ProductQuestion == nil {
		t.Fatalf("default qualification did not publish: %+v", f.s.Result.Worker)
	}
	p := sl01Preview(t, f, "fixture-answer-001")
	if p.Basis.Worker == nil || p.Basis.Worker.ProviderBinarySHA256 != f.s.Request.Runtime.Executable.SHA256 || p.Basis.Worker.QualificationSHA256 != productAnswerQualification {
		t.Fatalf("source validation lost the observed provider binding: %+v", p.Basis.Worker)
	}
}
func TestPH01DefaultQualificationKeepsModelAndEffortScope(t *testing.T) {
	f := newAW01Fixture(t)
	f.s.Request.Runtime.Model = "gpt-6-astra"
	if _, e := productRuntimeQualification(Dependencies{}, f.s); e != nil {
		t.Fatalf("supported Codex runtime rejected: %v", e)
	}
	f.s.Request.Runtime.Model = "fixture-model"
	if _, e := productRuntimeQualification(Dependencies{}, f.s); e == nil {
		t.Fatal("unsupported model qualified")
	}
	f.s.Request.Runtime.Model = "gpt-6-astra"
	f.s.Request.Delivery.ReasoningEffort = "high"
	if _, e := productRuntimeQualification(Dependencies{}, f.s); e == nil {
		t.Fatal("unsupported effort qualified")
	}
}
func TestSL01ReadAndPreviewHaveNoEffects(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	p := sl01Preview(t, f, "fixture-answer-001")
	if p.Request.DispatchPolicy != "when_parent_allows" || p.Basis.Question.Continuation.GlobalWriterQuiescence != "not_attested" {
		t.Fatal(p)
	}
	if _, e := ParseBoundAnswerRequest(mustSL01JSON(t, *p.Request)); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
		t.Fatal("read/preview mutated files")
	}
}

func TestSL01PreviewObservationIncludesItsBasisAndIssuedRequest(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	now := f.now
	f.d.Now = func() time.Time { now = now.Add(time.Microsecond); return now }
	p := sl01Preview(t, f, "preview-time-0001")
	observed, e := boundAnswerTime(p.ObservedAtUTC)
	if e != nil {
		t.Fatal(e)
	}
	basis, _ := boundAnswerTime(p.Basis.ObservedAtUTC)
	issued, _ := boundAnswerTime(p.Request.Validity.IssuedAtUTC)
	if observed.Before(basis) || observed.Before(issued) {
		t.Fatalf("preview predates its observed basis/request: preview=%s basis=%s issued=%s", p.ObservedAtUTC, p.Basis.ObservedAtUTC, p.Request.Validity.IssuedAtUTC)
	}
}
func TestSL01AwaitingStepDoesNotRedispatch(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	if e := workerStep(f.d, f.s); e != nil {
		t.Fatalf("awaiting must remain idle: %v", e)
	}
	if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
		t.Fatal("awaiting step changed preserved source")
	}
}
func mustSL01JSON(t *testing.T, v any) []byte {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	return b
}

func TestSL01ApplyOneWinnerAndReplayIsReadOnly(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	p := sl01Preview(t, f, "fixture-answer-001")
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e != nil || o.State != "answer_recorded" || o.AnswerRecorded == nil || !*o.AnswerRecorded {
		t.Fatalf("apply: %+v %v", o, e)
	}
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	f.now = f.now.Add(2 * time.Hour)
	o, e = ApplyBoundAnswer(f.d, *p.Request)
	if e != nil || o.AnswerRecorded == nil || !*o.AnswerRecorded {
		t.Fatalf("expired replay: %+v %v", o, e)
	}
	if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
		t.Fatal("replay wrote/woke")
	}
	other := *p.Request
	other.OptionID = "command_first"
	if _, e = ApplyBoundAnswer(f.d, other); BoundAnswerFailureOf(e).Code != "bound_answer_request_conflict" {
		t.Fatal(e)
	}
}
func TestSL01StatusNotRecordedHasNullUnknownFacts(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	o, e := ShowBoundAnswer(Dependencies{}, root, "unknown-answer-001")
	if e != nil || o.State != "not_recorded" || o.AnswerRecorded != nil || o.ContinuationRequested != nil || o.RunID != nil || o.ExecutionEffect != "unknown" {
		t.Fatalf("%+v %v", o, e)
	}
	if _, e = os.Stat(filepath.Join(root, ".ply")); !os.IsNotExist(e) {
		t.Fatal("status created store")
	}
}

func TestSL01OpaqueBindingCoversTheApprovedRequest(t *testing.T) {
	for _, change := range []string{"request_id", "validity", "expires", "dispatch_policy"} {
		t.Run(change, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			p := sl01Preview(t, f, "fixture-answer-001")
			r := *p.Request
			switch change {
			case "request_id":
				r.RequestID = "fixture-answer-002"
			case "validity":
				r.Validity.IssuedAtUTC = f.now.Add(-time.Second).Format(time.RFC3339Nano)
				r.Validity.ExpiresAtUTC = f.now.Add(299 * time.Second).Format(time.RFC3339Nano)
			case "expires":
				r.Validity.ExpiresAtUTC = f.now.Add(299 * time.Second).Format(time.RFC3339Nano)
			case "dispatch_policy":
				r.DispatchPolicy = "deferred_until_explicit_parent_resume"
			}
			o, e := ApplyBoundAnswer(f.d, r)
			if e == nil || o.State != "rejected" || o.AnswerRecorded == nil || *o.AnswerRecorded {
				t.Fatalf("changed approved request accepted: %+v %v", o, e)
			}
			if _, e = os.Stat(productClaimPath(r.Workspace, r.RunID, r.QuestionID)); !os.IsNotExist(e) {
				t.Fatal("modified request reserved a winner")
			}
		})
	}
}

func sl01Gate(t *testing.T, f *aw01Fixture, change func(*ParentControlGate)) {
	t.Helper()
	if e := WithWorkerParentGuard(f.s.Request.WorkspaceRoot, f.s.Request.ControlContract, func() error {
		g, e := ReadParentControl(f.s.Request.WorkspaceRoot, f.s.Request.ControlContract)
		if e != nil {
			return e
		}
		change(&g)
		g.Revision++
		return (Dependencies{}).replaceValue(ParentControlPath(f.s.Request.WorkspaceRoot, *f.s.Request.ControlContract.WeekID), g)
	}); e != nil {
		t.Fatal(e)
	}
}
func sl01Claim(t *testing.T, f *aw01Fixture) productAnswerClaim {
	t.Helper()
	p := sl01Preview(t, f, "fixture-answer-001")
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e != nil || o.AnswerRecorded == nil || !*o.AnswerRecorded {
		t.Fatalf("apply %+v %v", o, e)
	}
	c, e := productClaimForQuestion(f.s.Request.WorkspaceRoot, f.s.Result.RunID, p.Request.QuestionID)
	if e != nil {
		t.Fatal(e)
	}
	return c
}

func TestSL01FinalStartRechecksDeadlineAfterLaunchAuthorization(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	c := sl01Claim(t, f)
	if e := workerReserveProductContinuation(f.d, f.s, c); e != nil {
		t.Fatal(e)
	}
	f.reload()
	in, e := workerIntent(f.s, *f.s.ProductContinuation)
	if e != nil {
		t.Fatal(e)
	}
	f.d.Fault = func(point string) error {
		if point == "bound_answer_after_launch_authorization" {
			f.now = f.now.Add(2 * time.Hour)
		}
		return nil
	}
	starts := 0
	e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil })
	if e == nil || starts != 0 {
		t.Fatalf("expired absolute deadline still started provider: starts=%d err=%v", starts, e)
	}
	f.reload()
	if f.s.ProductLaunch == nil || f.s.ProductContinuation == nil {
		t.Fatal("single-launch authorization was erased")
	}
	f.d.Fault = nil
	if e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil }); e == nil || starts != 0 {
		t.Fatal("expired/possibly-started launch replayed")
	}
}

func TestSL01FirstApplyCannotCrossItsValidityAtPublication(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	p := sl01Preview(t, f, "expires-at-claim-01")
	f.d.Fault = func(point string) error {
		if point == "bound_answer_before_claim" {
			f.now = f.now.Add(301 * time.Second)
		}
		return nil
	}
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e == nil || o.State != "rejected" {
		t.Fatalf("expired first publication accepted: %+v %v", o, e)
	}
	if _, e = os.Stat(productClaimPath(f.s.Request.WorkspaceRoot, f.s.Result.RunID, p.Request.QuestionID)); !os.IsNotExist(e) {
		t.Fatal("expired apply reserved winner")
	}
}

func TestSL01PartialLaunchCannotBecomeAnUnstartedReservation(t *testing.T) {
	for _, missing := range []string{"marker", "state"} {
		t.Run(missing, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			c := sl01Claim(t, f)
			if e := workerReserveProductContinuation(f.d, f.s, c); e != nil {
				t.Fatal(e)
			}
			f.reload()
			in, e := workerIntent(f.s, *f.s.ProductContinuation)
			if e != nil {
				t.Fatal(e)
			}
			f.d.Fault = func(point string) error {
				if point == "bound_answer_after_launch_authorization" {
					return fmt.Errorf("crash at the single-start boundary")
				}
				return nil
			}
			starts := 0
			if e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil }); e == nil || starts != 0 {
				t.Fatal("fault did not preserve launch boundary")
			}
			f.reload()
			f.d.Fault = nil
			if missing == "marker" {
				if e = os.Remove(f.s.ProductLaunch.Locator); e != nil {
					t.Fatal(e)
				}
			} else {
				f.s.ProductLaunch = nil
				f.save()
			}
			o, e := ShowBoundAnswer(f.d, f.s.Request.WorkspaceRoot, c.Request.RequestID)
			if e != nil || o.State != "outcome_unknown" {
				t.Errorf("partial launch downgraded to unstarted: %+v %v", o, e)
			}
			if e = f.d.WorkerOwner.lock.Close(); e != nil {
				t.Fatal(e)
			}
			lock, e := workerLock(filepath.Join(workerDirectory(f.s), "owner.lock"), false)
			if e != nil {
				t.Fatal(e)
			}
			defer lock.Close()
			f.d.WorkerOwner = nil
			if e = workerClaimOwner(&f.d, f.s, lock); e == nil {
				t.Error("partial launch acquired fresh generation without a known attempt outcome")
			}
		})
	}
}

func TestSL01PartialIntentDoesNotRepairItselfIntoDispatch(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	c := sl01Claim(t, f)
	f.d.Fault = func(point string) error {
		if point == "bound_answer_after_continuation_intent" {
			return fmt.Errorf("interrupted after intent")
		}
		return nil
	}
	if e := workerReserveProductContinuation(f.d, f.s, c); e == nil {
		t.Fatal("fault not reached")
	}
	f.reload()
	f.d.Fault = nil
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	if e := workerReserveProductContinuation(f.d, f.s, c); e == nil {
		t.Fatal("partial reservation silently repaired")
	}
	if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
		t.Fatal("partial reservation changed")
	}
}

func TestSL01IncompleteAnswerCannotReauthorizeOwner(t *testing.T) {
	for _, boundary := range []string{"claim", "intent"} {
		t.Run(boundary, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			if boundary == "claim" {
				p := sl01Preview(t, f, "partial-owner-01")
				f.d.Fault = func(point string) error {
					if point == "bound_answer_after_claim" {
						return fmt.Errorf("partial claim")
					}
					return nil
				}
				if _, e := ApplyBoundAnswer(f.d, *p.Request); e == nil {
					t.Fatal("fault not reached")
				}
			} else {
				c := sl01Claim(t, f)
				f.d.Fault = func(point string) error {
					if point == "bound_answer_after_continuation_intent" {
						return fmt.Errorf("partial intent")
					}
					return nil
				}
				if e := workerReserveProductContinuation(f.d, f.s, c); e == nil {
					t.Fatal("fault not reached")
				}
			}
			f.d.Fault = nil
			f.reload()
			if e := f.d.WorkerOwner.lock.Close(); e != nil {
				t.Fatal(e)
			}
			lock, e := workerLock(filepath.Join(workerDirectory(f.s), "owner.lock"), false)
			if e != nil {
				t.Fatal(e)
			}
			defer lock.Close()
			f.d.WorkerOwner = nil
			if e = workerClaimOwner(&f.d, f.s, lock); e == nil {
				t.Fatal("incomplete answer granted a fresh owner generation")
			}
		})
	}
}

func TestSL01LaunchAuthorizationDoesNotAttestProviderDispatch(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	c := sl01Claim(t, f)
	if e := workerReserveProductContinuation(f.d, f.s, c); e != nil {
		t.Fatal(e)
	}
	f.reload()
	in, e := workerIntent(f.s, *f.s.ProductContinuation)
	if e != nil {
		t.Fatal(e)
	}
	f.d.Fault = func(point string) error {
		if point == "bound_answer_after_launch_authorization" {
			return fmt.Errorf("crash before Start")
		}
		return nil
	}
	if e = workerProductFinalStart(f.d, f.s, in, func() error { t.Fatal("Start reached"); return nil }); e == nil {
		t.Fatal("fault not reached")
	}
	q, e := ReadProductQuestion(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID)
	if e != nil || q.Worker == nil || q.Worker.ContinuationsDispatched != 0 {
		t.Fatalf("authorization falsely attested a physical dispatch: %+v %v", q.Worker, e)
	}
	// NA01: the recorded answer never hides the unknown launch outcome.
	sl01NextAction(t, f, "human", c.Request.RequestID, "not confirmed", "Do not answer again or resend")
}

func sl01NextAction(t *testing.T, f *aw01Fixture, actor string, want ...string) {
	t.Helper()
	for i := 0; i < 2; i++ {
		before, e := os.ReadFile(filepath.Join(f.s.Result.Paths.RunRoot, "state.json"))
		if e != nil {
			t.Fatal(e)
		}
		// This synthetic fixture reports binding drift through WorkflowShow,
		// which keeps precedence; project its persisted Run directly.
		f.reload()
		a := currentRunAction(f.d, f.s, f.s.Result)
		if a.Actor != actor {
			t.Fatalf("next action %+v (persisted %+v), want actor %s", a, f.s.Result.NextAction, actor)
		}
		for _, w := range want {
			if !strings.Contains(a.Message, w) {
				t.Fatalf("next action %q lacks %q", a.Message, w)
			}
		}
		after, e := os.ReadFile(filepath.Join(f.s.Result.Paths.RunRoot, "state.json"))
		if e != nil || string(after) != string(before) {
			t.Fatalf("read changed native state: %v", e)
		}
	}
}

func TestSL01CompletedEvidenceRequiresOrderedTimesAndConsumedIdentity(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	c := sl01Claim(t, f)
	q, e := productClaimQuestion(c)
	if e != nil {
		t.Fatal(e)
	}
	r, e := workerRunReserved(f.d, f.s)
	if e != nil {
		t.Fatal(e)
	}
	session, e := os.ReadFile(q.Session.Path)
	if e != nil {
		t.Fatal(e)
	}
	lines := bytes.Split(bytes.TrimSpace(session), []byte{'\n'})
	if e = os.WriteFile(q.Session.Path, append(session, append(lines[len(lines)-1], '\n')...), 0600); e != nil {
		t.Fatal(e)
	}
	answer := productCompletedAnswer{Status: "completed", Summary: "synthetic consumed proof", BlockReason: "", ConsumedAnswerRequestID: c.Request.RequestID, ConsumedAnswerRequestSHA256: c.RequestSHA256}
	binding, e := workflowKeep(f.d, filepath.Join(f.s.Request.WorkspaceRoot, "completed-answer.json"), answer)
	if e != nil {
		t.Fatal(e)
	}
	r.Answer = &binding
	r.StartedAt = f.now.Format(time.RFC3339Nano)
	r.FinishedAt = r.StartedAt
	for _, bad := range []string{"before_answer", "finish_before_start", "future_result", "consumed_id"} {
		t.Run(bad, func(t *testing.T) {
			x := r
			switch bad {
			case "before_answer":
				x.StartedAt = f.now.Add(-time.Second).Format(time.RFC3339Nano)
			case "finish_before_start":
				x.FinishedAt = f.now.Add(-time.Second).Format(time.RFC3339Nano)
			case "future_result":
				x.FinishedAt = f.now.Add(time.Second).Format(time.RFC3339Nano)
			case "consumed_id":
				a := answer
				a.ConsumedAnswerRequestID = "another-request-01"
				b, e := workflowKeep(f.d, filepath.Join(f.s.Request.WorkspaceRoot, "wrong-consumed.json"), a)
				if e != nil {
					t.Fatal(e)
				}
				x.Answer = &b
			}
			if _, e := productValidateCompleted(f.d, f.s, c, q, x); e == nil {
				t.Fatalf("invalid terminal evidence accepted: %s", bad)
			}
		})
	}
}

func TestSL01RequestWorkspaceMustBeDirectory(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	p := sl01Preview(t, f, "workspace-file-01")
	r := *p.Request
	r.Workspace = filepath.Join(f.s.Observed.Target.WorktreeLocator, "guide.md")
	if _, e := ParseBoundAnswerRequest(mustSL01JSON(t, r)); e == nil || BoundAnswerFailureOf(e).Code != "bound_answer_invalid_request" {
		t.Fatalf("regular file accepted as request workspace: %v", e)
	}
}

func TestSL01DriftBlocksEveryNewEffectBoundary(t *testing.T) {
	for _, boundary := range []string{"apply", "reserve", "start"} {
		t.Run(boundary, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			p := sl01Preview(t, f, "file-drift-0001")
			var c productAnswerClaim
			if boundary != "apply" {
				if _, e := ApplyBoundAnswer(f.d, *p.Request); e != nil {
					t.Fatal(e)
				}
				var e error
				c, e = productClaimForQuestion(f.s.Request.WorkspaceRoot, f.s.Result.RunID, p.Request.QuestionID)
				if e != nil {
					t.Fatal(e)
				}
			}
			if boundary == "start" {
				if e := workerReserveProductContinuation(f.d, f.s, c); e != nil {
					t.Fatal(e)
				}
				f.reload()
			}
			if e := os.WriteFile(filepath.Join(f.s.Observed.Target.WorktreeLocator, "guide.md"), []byte("changed while waiting\n"), 0600); e != nil {
				t.Fatal(e)
			}
			var e error
			switch boundary {
			case "apply":
				_, e = ApplyBoundAnswer(f.d, *p.Request)
			case "reserve":
				e = workerReserveProductContinuation(f.d, f.s, c)
			case "start":
				var in workerAttemptIntent
				in, e = workerIntent(f.s, *f.s.ProductContinuation)
				if e != nil {
					t.Fatal(e)
				}
				e = workerProductFinalStart(f.d, f.s, in, func() error { t.Error("drift started provider"); return nil })
			}
			if e == nil {
				t.Fatal("Task file drift authorized new effect")
			}
		})
	}
}

func TestSL01ApplyAndParentPauseSerializeInBothOrders(t *testing.T) {
	for _, first := range []string{"apply", "pause"} {
		t.Run(first, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			p := sl01Preview(t, f, "parent-race-0001")
			if first == "pause" {
				sl01Gate(t, f, func(g *ParentControlGate) { g.DesiredAction = "pause" })
				if _, e := ApplyBoundAnswer(f.d, *p.Request); e == nil {
					t.Fatal("preview from before pause won")
				}
				p = sl01Preview(t, f, "parent-race-0002")
				if o, e := ApplyBoundAnswer(f.d, *p.Request); e != nil || o.State != "continuation_deferred" {
					t.Fatalf("paused apply %+v %v", o, e)
				}
				sl01NextAction(t, f, "parent_control", p.Request.RequestID, "explicit parent resume", "no new human answer")
				return
			}
			entered, release, done := make(chan struct{}), make(chan struct{}), make(chan error, 1)
			f.d.Fault = func(point string) error {
				if point == "bound_answer_before_claim" {
					close(entered)
					<-release
				}
				return nil
			}
			go func() { _, e := ApplyBoundAnswer(f.d, *p.Request); done <- e }()
			<-entered
			paused := make(chan error, 1)
			go func() {
				paused <- WithWorkerParentGuard(f.s.Request.WorkspaceRoot, f.s.Request.ControlContract, func() error {
					g, e := ReadParentControl(f.s.Request.WorkspaceRoot, f.s.Request.ControlContract)
					if e != nil {
						return e
					}
					g.DesiredAction = "pause"
					g.Revision++
					return (Dependencies{}).replaceValue(ParentControlPath(f.s.Request.WorkspaceRoot, *f.s.Request.ControlContract.WeekID), g)
				})
			}()
			close(release)
			if e := <-done; e != nil {
				t.Fatal(e)
			}
			if e := <-paused; e != nil {
				t.Fatal(e)
			}
			f.d.Fault = nil
			o, e := ShowBoundAnswer(f.d, f.s.Request.WorkspaceRoot, p.Request.RequestID)
			if e != nil || o.State != "continuation_deferred" || o.AnswerRecorded == nil || !*o.AnswerRecorded {
				t.Fatalf("winning answer lost across pause %+v %v", o, e)
			}
			if _, e = workerProductWaiting(f.d, f.s); e != nil {
				t.Fatal(e)
			}
			f.reload()
			if f.s.ProductContinuation != nil {
				t.Fatal("parent pause dispatched continuation")
			}
		})
	}
}

func TestSL01ConcurrentAnswersHaveOneImmutableWinner(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	requests := []BoundAnswerRequest{}
	for i := 0; i < 12; i++ {
		requests = append(requests, *sl01Preview(t, f, fmt.Sprintf("race-answer-%04d", i)).Request)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	wins := 0
	for _, r := range requests {
		wg.Add(1)
		go func(r BoundAnswerRequest) {
			defer wg.Done()
			o, e := ApplyBoundAnswer(f.d, r)
			mu.Lock()
			defer mu.Unlock()
			if e == nil && o.AnswerRecorded != nil && *o.AnswerRecorded {
				wins++
			}
		}(r)
	}
	wg.Wait()
	if wins != 1 {
		t.Fatalf("winners=%d", wins)
	}
	files, e := os.ReadDir(filepath.Join(productAnswerRoot(f.s.Request.WorkspaceRoot), "claims"))
	if e != nil || len(files) != 1 {
		t.Fatalf("claim inventory %v %v", files, e)
	}
}
func TestSL01PartialIndexPreservesIdentityAndBlocksReplacement(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	p := sl01Preview(t, f, "partial-answer-001")
	other := sl01Preview(t, f, "partial-answer-002")
	f.d.Fault = func(point string) error {
		if point == "bound_answer_after_claim" {
			return fmt.Errorf("injected index interruption")
		}
		return nil
	}
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e == nil || o.State != "outcome_unknown" || o.RequestSHA256 == nil || *o.RequestSHA256 != *p.RequestSHA256 {
		t.Fatalf("lost partial identity %+v %v", o, e)
	}
	f.d.Fault = nil
	before := aw01Inventory(t, f.s.Request.WorkspaceRoot)
	o, e = ApplyBoundAnswer(f.d, *p.Request)
	if e != nil || o.State != "outcome_unknown" {
		t.Fatalf("replay %+v %v", o, e)
	}
	_, _ = ShowBoundAnswer(f.d, f.s.Request.WorkspaceRoot, p.Request.RequestID)
	if !reflect.DeepEqual(before, aw01Inventory(t, f.s.Request.WorkspaceRoot)) {
		t.Fatal("partial replay/status repaired or woke")
	}
	o, e = ApplyBoundAnswer(f.d, *other.Request)
	if e == nil || o.State != "rejected" {
		t.Fatalf("replacement %+v %v", o, e)
	}
}
func TestSL01PausedPreviewKeepsQuestionAndDefers(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	before := sl01Preview(t, f, "before-pause-001")
	sl01Gate(t, f, func(g *ParentControlGate) { g.DesiredAction = "pause" })
	if handled, e := workerProductWaiting(f.d, f.s); !handled || e != nil {
		t.Fatal(handled, e)
	}
	f.reload()
	p := sl01Preview(t, f, "after-pause-0001")
	if p.Request.DispatchPolicy != "deferred_until_explicit_parent_resume" || p.Basis.Question.QuestionSHA256 != before.Basis.Question.QuestionSHA256 {
		t.Fatal("pause changed question or omitted request policy")
	}
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e != nil || o.State != "continuation_deferred" {
		t.Fatalf("%+v %v", o, e)
	}
	if _, e = workerProductWaiting(f.d, f.s); e != nil {
		t.Fatal(e)
	}
	f.reload()
	if f.s.ProductContinuation != nil {
		t.Fatal("pause reserved continuation")
	}
}
func TestSL01FinalStartRechecksParentAndPreservesReservedIntent(t *testing.T) {
	for _, obstacle := range []string{"pause", "stop", "deadline", "active_run"} {
		t.Run(obstacle, func(t *testing.T) {
			f := newSL01Fixture(t)
			sl01Park(t, f)
			c := sl01Claim(t, f)
			if e := workerReserveProductContinuation(f.d, f.s, c); e != nil {
				t.Fatal(e)
			}
			f.reload()
			job := *f.s.ProductContinuation
			in, e := workerIntent(f.s, job)
			if e != nil {
				t.Fatal(e)
			}
			sl01Gate(t, f, func(g *ParentControlGate) {
				switch obstacle {
				case "pause", "stop":
					g.DesiredAction = obstacle
				case "deadline":
					f.now = f.now.Add(2 * time.Hour)
				case "active_run":
					g.ActiveRun = "wfr_" + strings.Repeat("0", 64)
				}
			})
			starts := 0
			e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil })
			if e == nil || starts != 0 {
				t.Fatalf("parent obstacle bypassed %d %v", starts, e)
			}
			o, e := ShowBoundAnswer(f.d, f.s.Request.WorkspaceRoot, c.Request.RequestID)
			if e != nil || o.State != "continuation_reserved" || o.Receipt.ContinuationIntentSHA256 == nil || *o.Receipt.ContinuationIntentSHA256 != job.IntentHash {
				t.Fatalf("reserved identity lost %+v %v", o, e)
			}
			if obstacle == "pause" {
				if _, e = workerProductWaiting(f.d, f.s); e != nil {
					t.Fatal(e)
				}
				f.reload()
				basis, e := WorkerControls(f.d).Show(f.s.Request.WorkspaceRoot, f.s.Result.RunID)
				if e != nil || basis.UnavailableReason == nil || *basis.UnavailableReason != "parent_control_unsupported" {
					t.Fatalf("unstarted reservation cannot acknowledge parent pause %+v %v", basis, e)
				}
				sl01Gate(t, f, func(g *ParentControlGate) { g.DesiredAction = "continue"; g.ResumeAllowed = true })
				if _, e = workerProductWaiting(f.d, f.s); e != nil {
					t.Fatal(e)
				}
				f.reload()
				if e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil }); e != nil || starts != 1 {
					t.Fatalf("same intent cannot start once after explicit parent resume %d %v", starts, e)
				}
				if e = workerProductFinalStart(f.d, f.s, in, func() error { starts++; return nil }); e == nil || starts != 1 {
					t.Fatalf("duplicate start %d %v", starts, e)
				}
			}
		})
	}
}
func TestSL01RestartRebindsGenerationWithoutAnotherTurn(t *testing.T) {
	f := newSL01Fixture(t)
	sl01Park(t, f)
	p := sl01Preview(t, f, "restart-answer-01")
	old := f.d
	oldGeneration := f.s.Result.Worker.Generation
	if e := f.d.WorkerOwner.lock.Close(); e != nil {
		t.Fatal(e)
	}
	lock, e := workerLock(filepath.Join(workerDirectory(f.s), "owner.lock"), false)
	if e != nil {
		t.Fatal(e)
	}
	defer lock.Close()
	f.d.WorkerOwner = nil
	if e = workerClaimOwner(&f.d, f.s, lock); e != nil {
		t.Fatal(e)
	}
	f.reload()
	if f.s.Result.Worker.Generation != oldGeneration+1 || f.s.Result.Worker.Phase != "awaiting_product_answer" || f.s.Result.Worker.Iteration != 1 {
		t.Fatal(f.s.Result.Worker)
	}
	if e = workerLiveOwner(old, f.s); e == nil {
		t.Fatal("stale owner gained authority")
	}
	o, e := ApplyBoundAnswer(f.d, *p.Request)
	if e == nil || o.State != "rejected" {
		t.Fatalf("old generation preview accepted %+v %v", o, e)
	}
	p2 := sl01Preview(t, f, "restart-answer-02")
	if p2.Basis.Question.QuestionSHA256 != p.Basis.Question.QuestionSHA256 {
		t.Fatal("restart changed immutable question")
	}
}
