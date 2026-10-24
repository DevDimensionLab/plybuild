package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func provisionalSequence() []any {
	return []any{map[string]any{"surface": "unknown", "launch_pending": true}, map[string]any{}}
}

// Assert against the producer's temporal history, rather than the final contents
// of the reused inspect request. Every input must follow positive evidence.
func assertProvisionalHistory(t *testing.T, f workflowFixture, target string, want string) {
	t.Helper()
	var model struct {
		History []struct {
			Operation, Target, Surface, Kind string
			Deadline                         int64
		} `json:"safe_history"`
	}
	raw, err := os.ReadFile(f.Model)
	if err != nil || json.Unmarshal(raw, &model) != nil {
		t.Fatalf("producer history: %v", err)
	}
	last, unknown, kinds, deadline := "", 0, []string{}, int64(0)
	for _, h := range model.History {
		if h.Target != target {
			continue
		}
		if h.Operation == "inspect" {
			last = h.Surface
			if last == "unknown" {
				unknown++
			}
		} else {
			if last != "provider_prompt" && !(last == "update_skip_once" && h.Kind == startupSkip) {
				t.Fatalf("%s input followed %s", h.Kind, last)
			}
			if deadline != 0 && deadline != h.Deadline {
				t.Fatal("startup input deadline renewed")
			}
			deadline = h.Deadline
			kinds = append(kinds, h.Kind)
		}
	}
	minimumUnknown := 1
	if want == "bootstrap,task" {
		minimumUnknown = 2 // initial loading and post-bootstrap working
	}
	if unknown < minimumUnknown || strings.Join(kinds, ",") != want {
		t.Fatalf("unknown observations=%d, inputs=%v; want %s", unknown, kinds, want)
	}
}

func TestProvisionalStartupOwner(t *testing.T) {
	for _, bootstrap := range []bool{false, true} {
		t.Run(map[bool]string{false: "configured", true: "bootstrap_working"}[bootstrap], func(t *testing.T) {
			var f workflowFixture
			if bootstrap {
				f = workflowBootstrapFixture(t, "normal")
			} else {
				f = workflowTestFixture(t)
			}
			f.D.HerdrTimeout = 5 * time.Second
			workflowTestModel(t, f, map[string]any{"safe_inspections": provisionalSequence(), "safe_unknown_while_working": true})
			o, err := workflowReadinessStart(t, f)
			if err != nil || o.Transport.State != "working" {
				t.Fatalf("provisional owner failed: %v %+v", err, o.Transport)
			}
			want := "task"
			if bootstrap {
				want = "bootstrap,task"
			}
			assertProvisionalHistory(t, f, o.Transport.PaneID, want)
			if workflowTestCalls(t, f, "agent start") != 1 || workflowUnguardedCalls(t, f) != 0 {
				t.Fatal("startup restarted or used unguarded input")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestProvisionalStartupReviewer(t *testing.T) {
	for _, bootstrap := range []bool{false, true} {
		t.Run(map[bool]string{false: "configured", true: "bootstrap_working"}[bootstrap], func(t *testing.T) {
			f, o, _ := reviewRoundFixture(t)
			f.D.HerdrTimeout = 5 * time.Second
			reviewRoundCandidate(t, f, "fixed")
			reviewRoundIn(t, "w-fixture:root")
			workflowTestModel(t, f.workflowFixture, map[string]any{"review_safe_inspections": provisionalSequence(), "review_sessionless": bootstrap, "review_safe_unknown_while_working": true})
			r, err := WorkflowDeliveryRequestReview(reviewRoundAsCoordinator(f), f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
			if err != nil || reviewRoundLatest(t, r).ReviewerStart != "started" {
				t.Fatalf("provisional reviewer failed: %v", err)
			}
			want := "task"
			if bootstrap {
				want = "bootstrap,task"
			}
			starts := workflowProviderCalls(t, f.workflowFixture, "agent start")
			if len(starts) != 2 || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
				t.Fatal("review restarted or used unguarded input")
			}
			for i, arg := range starts[1] {
				if arg == "--pane" {
					assertProvisionalHistory(t, f.workflowFixture, starts[1][i+1], want)
					return
				}
			}
			t.Fatal("missing reviewer pane")
		})
	}
}

func TestProvisionalStartupMenuThenUnknown(t *testing.T) {
	f := workflowTestFixture(t)
	f.D.HerdrTimeout = 5 * time.Second
	workflowTestModel(t, f, map[string]any{"safe_skip": true, "safe_inspections": []any{map[string]any{}, map[string]any{"surface": "unknown"}, map[string]any{}}})
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatal(err)
	}
	assertProvisionalHistory(t, f, o.Transport.PaneID, "skip_once,task")
	workflowAssertNoReplay(t, f, o)
}

func TestProvisionalStartupOwnerResume(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_inspections": provisionalSequence(), "safe_inspect_count": 0})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" {
		t.Fatalf("resume preview: %v %+v", err, p)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	assertProvisionalHistory(t, f.workflowFixture, o.Transport.PaneID, "task")
	if o.Transport.AgentSessionID != f.Owner.Transport.AgentSessionID || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatal("resume lost session or guarded input")
	}
}

func TestProvisionalStartupUnknownExpiresWithoutInput(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	f.D.HerdrTimeout = 1200 * time.Millisecond
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_surface": "unknown"})
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil || !(strings.Contains(err.Error(), "deadline") || strings.Contains(err.Error(), "timeout")) {
		t.Fatalf("unknown did not expire in original window: %v", err)
	}
	var window struct {
		Deadline int64 `json:"deadline_unix_ms"`
	}
	windowRaw, windowErr := os.ReadFile(filepath.Join(o.Paths.RunRoot, "startup-deadline.json"))
	if windowErr != nil || json.Unmarshal(windowRaw, &window) != nil || window.Deadline > time.Now().UnixMilli() {
		t.Fatalf("stop preceded original deadline: %v", windowErr)
	}
	assertProvisionalHistory(t, f.workflowFixture, o.Transport.PaneID, "")
	var last startupInspection
	raw, readErr := os.ReadFile(filepath.Join(o.Paths.RunRoot, "startup-last-inspection.json"))
	if readErr != nil || json.Unmarshal(raw, &last) != nil {
		t.Fatalf("last inspection unavailable: %v", readErr)
	}
	if last.Surface != "unknown" || len(safeStartupKinds(t, f.workflowFixture)) != 0 {
		t.Fatal("unknown diagnosis lost or input sent")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_surface": "provider_prompt"})
	_, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, 5*time.Second)
	if err == nil || !strings.Contains(err.Error(), "deadline") || len(safeStartupKinds(t, f.workflowFixture)) != 0 {
		t.Fatalf("resume renewed expired window: %v", err)
	}
}

func TestProvisionalStartupInvalidUnknownStops(t *testing.T) {
	for name, steps := range map[string][]any{
		"prompt_conflict": {map[string]any{"surface": "unknown", "prompt_ready": true}, map[string]any{}},
		"menu_conflict":   {map[string]any{"surface": "unknown", "menu": map[string]any{"layout_id": "foreign"}}, map[string]any{}},
		"session_drift":   {map[string]any{"surface": "unknown", "native_session": "foreign"}, map[string]any{}},
		"binding_drift":   {map[string]any{"surface": "unknown", "binding": map[string]any{"name": "foreign"}}, map[string]any{}},
		"epoch_drift":     {map[string]any{"surface": "unknown"}, map[string]any{"server_epoch": "replacement"}},
		"launch_drift":    {map[string]any{"surface": "unknown"}, map[string]any{"managed_launch_id": "replacement"}},
		"schema":          {map[string]any{"surface": "unknown", "prompt_ready": nil}, map[string]any{}},
	} {
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			workflowTestModel(t, f, map[string]any{"safe_inspections": steps})
			o, err := workflowReadinessStart(t, f)
			if err == nil || strings.Contains(err.Error(), "deadline") || len(safeStartupKinds(t, f)) != 0 {
				t.Fatalf("invalid unknown was awaited or sent input: %v", err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestProvisionalStartupAdmissionStaysStrict(t *testing.T) {
	// Ready at the wait boundary, provisional again at the actual input boundary.
	// The strict inspection must refuse rather than sending or starting a new wait.
	for _, bootstrap := range []bool{false, true} {
		t.Run(map[bool]string{false: "task", true: "bootstrap"}[bootstrap], func(t *testing.T) {
			var f workflowFixture
			if bootstrap {
				f = workflowBootstrapFixture(t, "normal")
			} else {
				f = workflowTestFixture(t)
			}
			workflowTestModel(t, f, map[string]any{"safe_inspections": []any{map[string]any{}, map[string]any{"surface": "unknown"}, map[string]any{}}})
			o, err := workflowReadinessStart(t, f)
			if err == nil || !strings.Contains(err.Error(), "surface=unknown") || len(safeStartupKinds(t, f)) != 0 {
				t.Fatalf("provisional input admission: %v", err)
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestProvisionalStartupUncertainBootstrapCannotWait(t *testing.T) {
	for _, outcome := range []string{"unknown", "partial", "not_sent", "lost_reply"} {
		t.Run(outcome, func(t *testing.T) {
			f := workflowBootstrapFixture(t, "normal")
			model := map[string]any{"safe_unknown_while_working": true, "safe_outcome": "unknown", "safe_bytes": "unknown"}
			if outcome == "partial" {
				model["safe_outcome"] = "submitted"
			} else if outcome == "not_sent" {
				model["safe_outcome"], model["safe_bytes"] = "not_sent", "none"
			}
			if outcome == "lost_reply" {
				model = map[string]any{"safe_unknown_while_working": true, "safe_lost": "bootstrap"}
			}
			workflowTestModel(t, f, model)
			o, err := workflowReadinessStart(t, f)
			if err == nil || strings.Join(safeStartupKinds(t, f), ",") != "bootstrap" {
				t.Fatalf("uncertain bootstrap became provisional readiness: %v, %v", err, safeStartupKinds(t, f))
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestProvisionalStartupUncertainSkipCannotWait(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"safe_skip": true, "safe_lost": "skip_once", "safe_inspections": []any{map[string]any{}, map[string]any{"surface": "unknown"}, map[string]any{}}})
	o, err := workflowReadinessStart(t, f)
	if err == nil || strings.Join(safeStartupKinds(t, f), ",") != "skip_once" {
		t.Fatalf("uncertain skip became provisional readiness: %v", err)
	}
	workflowAssertNoReplay(t, f, o)
}

func TestProvisionalStartupReadRefusalIsFatal(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"safe_inspections": provisionalSequence(), "safe_inspect_denied_at": 1})
	o, err := workflowReadinessStart(t, f)
	var model struct {
		Count int `json:"safe_inspect_count"`
	}
	raw, readErr := os.ReadFile(f.Model)
	if readErr != nil || json.Unmarshal(raw, &model) != nil {
		t.Fatalf("inspect count: %v", readErr)
	}
	if err == nil || !strings.Contains(err.Error(), "permission_denied") || model.Count != 2 || len(safeStartupKinds(t, f)) != 0 {
		t.Fatalf("read refusal retried or gained input: %v, calls=%d", err, model.Count)
	}
	workflowAssertNoReplay(t, f, o)
}

func TestProvisionalStartupCLI(t *testing.T) {
	binary := automaticCLIBinary(t)
	f := newDeliveryFixture(t, "codex")
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.File = writeAny(t, f.R.WorkspaceRoot, "provisional-request.json", f.R)
	workflowTestModel(t, f.workflowFixture, map[string]any{"safe_inspections": provisionalSequence(), "safe_require_thread_profile": true})
	_, output, err := automaticCLI(binary, f, "run", "start", "--file", f.File, "--check")
	var p WorkflowPreview
	if err != nil || json.Unmarshal(output, &p) != nil || p.Confirmation == nil {
		t.Fatalf("CLI preview: %v %s", err, output)
	}
	o := automaticCLIOK(t, binary, f, "run", "start", "--file", f.File, "--apply", "--confirm", *p.Confirmation)
	if o.Transport.State != "working" || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowUnguardedCalls(t, f.workflowFixture) != 0 {
		t.Fatalf("CLI did not finish in one guarded start: %+v", o.Transport)
	}
	assertProvisionalHistory(t, f.workflowFixture, o.Transport.PaneID, "task")
}
