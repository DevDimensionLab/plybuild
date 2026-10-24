package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// These inspect actual published CLI calls; legacy agent.prompt is never counted
// as a guarded effect.
func safeStartupKinds(t *testing.T, f workflowFixture) []string {
	t.Helper()
	b, _ := os.ReadFile(f.Calls)
	var kinds []string
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var call struct{ Argv []string }
		if json.Unmarshal([]byte(line), &call) != nil || len(call.Argv) < 5 {
			continue
		}
		if strings.Join(call.Argv[:3], " ") == "agent startup submit" {
			raw, err := os.ReadFile(call.Argv[4])
			if err != nil {
				t.Fatal(err)
			}
			var p struct{ Kind string }
			if err = json.Unmarshal(raw, &p); err != nil {
				t.Fatal(err)
			}
			kinds = append(kinds, p.Kind)
		}
	}
	return kinds
}

func TestWorkflowStartupRequiresPositivePromptInBothBranches(t *testing.T) {
	for _, sessionless := range []bool{false, true} {
		for _, surface := range []string{"blocked", "unknown", "provider_prompt"} {
			t.Run(surface+map[bool]string{true: "/sessionless", false: "/session"}[sessionless], func(t *testing.T) {
				f := workflowTestFixture(t)
				f.D.HerdrTimeout = 2 * time.Second
				model := map[string]any{"safe_surface": surface, "safe_prompt_ready": false}
				if sessionless {
					model["bootstrap_mode"] = "normal"
					model["interactive_ready"] = true
					model["bootstrap_run_root"] = workflowPaths(f.R, 0).RunRoot
					model["bootstrap_expected_prompt"] = workflowExpectedBootstrapPrompt
				}
				workflowTestModel(t, f, model)
				o, err := workflowReadinessStart(t, f)
				if err == nil || len(safeStartupKinds(t, f)) != 0 || workflowUnguardedCalls(t, f) != 0 {
					t.Fatalf("unsafe surface gained input: %v", err)
				}
				workflowAssertNoReplay(t, f, o)
			})
		}
	}
}
func TestWorkflowStartupRequiresProducerCapability(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"safe_capability": false})
	o, err := workflowReadinessStart(t, f)
	if err == nil || !strings.Contains(err.Error(), "safe_startup_v1") || workflowUnguardedCalls(t, f) != 0 || len(safeStartupKinds(t, f)) != 0 {
		t.Fatalf("unsupported server received input: %v", err)
	}
	workflowAssertNoReplay(t, f, o)
}
func TestWorkflowStartupSkipsOnceBeforeReadiness(t *testing.T) {
	for _, lost := range []bool{false, true} {
		t.Run(map[bool]string{false: "reply", true: "lost"}[lost], func(t *testing.T) {
			f := workflowTestFixture(t)
			f.D.HerdrTimeout = 5 * time.Second
			m := map[string]any{"safe_skip": true}
			if lost {
				m["safe_lost"] = "skip_once"
			}
			workflowTestModel(t, f, m)
			o, err := workflowReadinessStart(t, f)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Join(safeStartupKinds(t, f), ","); got != "skip_once,task" {
				t.Fatalf("wrong guarded order %s", got)
			}
			if workflowUnguardedCalls(t, f) != 0 {
				t.Fatal("unguarded fallback")
			}
			if _, err = os.Stat(filepath.Join(o.Paths.RunRoot, "startup-skip-attempt.json")); err != nil {
				t.Fatal(err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}
func TestDeliveryReviewStartupUsesGuardedInput(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		t.Run(map[bool]string{false: "positive", true: "blocked"}[blocked], func(t *testing.T) {
			f, o, _ := reviewRoundFixture(t)
			reviewRoundCandidate(t, f, "fixed")
			reviewRoundIn(t, "w-fixture:root")
			m := map[string]any{"review_safe_skip": true}
			if blocked {
				m["review_safe_surface"] = "blocked"
				m["review_safe_skip"] = false
			}
			workflowTestModel(t, f.workflowFixture, m)
			before := len(safeStartupKinds(t, f.workflowFixture))
			r, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
			if err != nil && !blocked {
				t.Fatal(err)
			}
			kinds := safeStartupKinds(t, f.workflowFixture)[before:]
			if blocked {
				if len(kinds) != 0 || reviewRoundLatest(t, r).ReviewerStart == "started" {
					t.Fatal("blocked reviewer gained input")
				}
			} else if strings.Join(kinds, ",") != "skip_once,task" {
				t.Fatalf("review input not guarded: %v", kinds)
			}
		})
	}
}

func TestWorkflowStartupRejectsIncompleteOrDriftedEvidence(t *testing.T) {
	cases := map[string]map[string]any{
		"contract": {"safe_override": map[string]any{"contract": "future_contract"}},
		"binding":  {"safe_override": map[string]any{"binding": map[string]any{"name": "foreign"}}},
		"session":  {"safe_override": map[string]any{"native_session": "foreign-session"}},
		"stale":    {"safe_override": map[string]any{"observed_at_unix_ms": 1}},
		"menu":     {"safe_skip": true, "safe_override": map[string]any{"menu": map[string]any{"layout_id": "unqualified", "selected_index": 0, "skip_action": "skip_once"}}},
	}
	cases["menu_missing_selection"] = map[string]any{"safe_skip": true, "safe_override": map[string]any{"menu": map[string]any{"layout_id": "codex-0.164-picker-v1", "skip_action": "skip_once"}}}
	for _, field := range []string{"prompt_ready", "launch_pending", "binding", "native_session", "server_epoch", "evidence_token", "observed_at_unix_ms", "surface"} {
		cases["missing_"+field] = map[string]any{"safe_drop": []string{field}}
	}
	for name, model := range cases {
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			workflowTestModel(t, f, model)
			o, err := workflowReadinessStart(t, f)
			if err == nil || len(safeStartupKinds(t, f)) != 0 || workflowUnguardedCalls(t, f) != 0 {
				t.Fatalf("invalid inspection allowed startup input: %v", err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}
func TestWorkflowStartupUnknownEffectsNeverReplay(t *testing.T) {
	for _, outcome := range []string{"unknown", "not_sent", "lost_reply"} {
		t.Run(outcome, func(t *testing.T) {
			f := workflowTestFixture(t)
			workflowTestModel(t, f, map[string]any{"safe_outcome": outcome, "safe_bytes": "unknown"})
			if outcome == "lost_reply" {
				workflowTestModel(t, f, map[string]any{"safe_lost": "task"})
			}
			o, err := workflowReadinessStart(t, f)
			if err == nil || strings.Join(safeStartupKinds(t, f), ",") != "task" {
				t.Fatalf("uncertain Task input was not consumed: %v", err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}
func TestDeliveryReviewPendingStartupNeverReplays(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_safe_lost": "task"})
	d := reviewRoundAsCoordinator(f)
	_, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err == nil {
		t.Fatal("lost brief became a successful start")
	}
	before, _ := os.ReadFile(f.Calls)
	again, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err != nil {
		t.Fatal(err)
	}
	if reviewRoundLatest(t, again).ReviewerStart != "prompt_unknown" {
		t.Fatal("unknown brief lost")
	}
	// The real pending-round branch may inspect the coordinator pane, but must
	// neither launch nor submit any input into the reviewer again.
	after, _ := os.ReadFile(f.Calls)
	if len(after) <= len(before) {
		t.Fatal("pending review never reached native coordinator checks")
	}
	kinds := safeStartupKinds(t, f.workflowFixture)
	if strings.Join(kinds, ",") != "task,task" || workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 {
		t.Fatalf("pending review replayed startup: %v", kinds)
	}
}
func TestWorkflowStartupDeadlineCannotBeRevived(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	f.D.HerdrTimeout = 1200 * time.Millisecond
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_prompt_ready": false})
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil {
		t.Fatal("expected shared deadline")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_prompt_ready": true})
	if _, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, 5*time.Second); err == nil || !strings.Contains(err.Error(), "deadline") {
		t.Fatalf("resume revived expired input window: %v", err)
	}
	if len(safeStartupKinds(t, f.workflowFixture)) != 0 {
		t.Fatal("expired startup sent input")
	}
}

func TestDeliveryReviewSessionlessStartupUsesGuardedBootstrap(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	// This positive round includes skip, bootstrap, settlement and brief. Its
	// fixture budget must accommodate all of those real CLI exchanges.
	f.D.HerdrTimeout = 10 * time.Second
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_sessionless": true, "review_safe_skip": true, "review_safe_lost": "bootstrap"})
	r, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err != nil {
		t.Fatal(err)
	}
	if reviewRoundLatest(t, r).ReviewerStart != "started" || strings.Join(safeStartupKinds(t, f.workflowFixture), ",") != "task,skip_once,bootstrap,task" || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatal("sessionless reviewer did not guard bootstrap and brief")
	}
}

func TestOwnerResumeStartupRequiresPositivePrompt(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.Now = func() time.Time { return time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC) }
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_surface": "blocked"})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("resume fixture not ready: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	before := len(safeStartupKinds(t, f.workflowFixture))
	_, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil && !strings.Contains(err.Error(), "safe startup") {
		t.Fatalf("resume stopped before guarded branch: %v", err)
	}
	if err == nil || len(safeStartupKinds(t, f.workflowFixture)) != before || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatalf("closed-tab resume sent into blocked provider: %v", err)
	}
}

func TestWorkflowStartupSkipCannotReplayOrCrossEpoch(t *testing.T) {
	for _, change := range []string{"menu_remains", "epoch_after_skip"} {
		t.Run(change, func(t *testing.T) {
			f := workflowTestFixture(t)
			m := map[string]any{"safe_skip": true, "safe_lost": "skip_once"}
			if change == "menu_remains" {
				m["safe_menu_remains"] = true
			} else {
				m["safe_epoch_after_skip"] = "replacement-epoch"
			}
			workflowTestModel(t, f, m)
			o, err := workflowReadinessStart(t, f)
			if err == nil || strings.Join(safeStartupKinds(t, f), ",") != "skip_once" || workflowUnguardedCalls(t, f) != 0 {
				t.Fatalf("uncertain skip gained another input: %v", err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowStartupRefusalStopsDependentInput(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"safe_skip": true, "safe_denied": "skip_once"})
	o, err := workflowReadinessStart(t, f)
	if err == nil || strings.Join(safeStartupKinds(t, f), ",") != "skip_once" || !strings.Contains(err.Error(), "permission_denied") {
		t.Fatalf("known refusal became readiness: %v", err)
	}
	workflowAssertNoReplay(t, f, o)
}

func TestDeliveryReviewConcurrentStartupConsumesOneRound(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	d := reviewRoundAsCoordinator(f)
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			_, _ = WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
		}()
	}
	<-done
	<-done
	shown, err := WorkflowShow(d, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(shown.DeliveryStatus.Review.Rounds) != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") > 2 {
		t.Fatal("concurrent requests overlapped reviewer launches")
	}
	kinds := safeStartupKinds(t, f.workflowFixture)
	if len(kinds) > 2 || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatalf("concurrent requests replayed input: %v", kinds)
	}
}

func TestReviewR1RejectOwnerSessionBeforeSkip(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed")
	reviewRoundIn(t, "w-fixture:root")
	workflowTestModel(t, f.workflowFixture, map[string]any{"review_session_is_owner": true, "review_safe_skip": true})
	before := len(safeStartupKinds(t, f.workflowFixture))
	_, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	kinds := safeStartupKinds(t, f.workflowFixture)[before:]
	if err == nil {
		t.Fatal("owner session was accepted")
	}
	if len(kinds) != 0 {
		t.Fatalf("input reached owner's native session before identity rejection: kinds=%v err=%v", kinds, err)
	}
}

func TestReviewR1RejectInspectionSessionDriftBeforeSkip(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"safe_skip": true, "safe_override": map[string]any{"native_session": "foreign-session"}})
	_, err := workflowReadinessStart(t, f)
	kinds := safeStartupKinds(t, f)
	if err == nil {
		t.Fatal("inspection drift was accepted")
	}
	if len(kinds) != 0 {
		t.Fatalf("input sent before native session drift rejection: kinds=%v err=%v", kinds, err)
	}
}

func TestReviewR1OwnerResumeRejectsSessionBeforeSkip(t *testing.T) {
	for _, drift := range []string{"inspection", "agent"} {
		t.Run(drift, func(t *testing.T) {
			f := newOwnerResumeFixture(t, "codex", false)
			model := map[string]any{"safe_skip": true}
			if drift == "inspection" {
				model["safe_override"] = map[string]any{"native_session": "foreign-session"}
			} else {
				model["agent_session_id"] = "foreign-session"
				model["safe_override"] = map[string]any{"native_session": f.Owner.Transport.AgentSessionID}
			}
			workflowTestModel(t, f.workflowFixture, model)
			p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err != nil || p.State != "ready" {
				t.Fatalf("fixture failed before resume: %+v %v", p, err)
			}
			ownerResumePreserve(t, f, p)
			before := len(safeStartupKinds(t, f.workflowFixture))
			_, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
			if err == nil {
				t.Fatal("foreign session became ready")
			}
			if kinds := safeStartupKinds(t, f.workflowFixture)[before:]; len(kinds) != 0 {
				t.Fatalf("resume input preceded identity check: %v, %v", kinds, err)
			}
		})
	}
}

func TestReviewR1OwnerSessionlessMenuStillBootstraps(t *testing.T) {
	f := workflowBootstrapFixture(t, "normal")
	f.D.HerdrTimeout = 10 * time.Second
	workflowTestModel(t, f, map[string]any{"safe_skip": true})
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(safeStartupKinds(t, f), ",") != "skip_once,bootstrap,task" || workflowUnguardedCalls(t, f) != 0 || o.Transport.AgentSessionID != "fixture-session" {
		t.Fatal("legitimate sessionless menu path changed")
	}
}
func TestReviewR1ResumeSessionlessMenuStillBindsOwner(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.HerdrTimeout = 10 * time.Second
	f.In.Timeout = 10 * time.Second
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_skip": true, "session_absent_gets": 1})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("fixture failed before resume: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	before := len(safeStartupKinds(t, f.workflowFixture))
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(safeStartupKinds(t, f.workflowFixture)[before:], ",") != "skip_once,task" || workflowUnguardedCalls(t, f.workflowFixture) != 0 || o.Transport.AgentSessionID != f.Owner.Transport.AgentSessionID {
		t.Fatal("resume sessionless menu did not preserve bound owner")
	}
}
