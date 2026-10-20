package taskrun

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// The bound pane exists until the coordinator closes it. Closing it makes the
// pane and its managed agent positively absent in the synthetic Herdr.
const retireOpenPaneStandin = `
if cmd == ["pane", "list"] and model.get("retire_inventory_denied"):
    print(json.dumps({"error":{"code":"permission_denied", "message":"inspection denied"}}), file=sys.stderr)
    raise SystemExit(1)
if cmd == ["pane", "close"]:
    if model.get("quota_append_on_close"):
        with open(model["quota_source_path"], "a") as transcript:
            transcript.write(model["quota_append_on_close"])
    if not model.get("close_ignored"):
        model["closed_transport"] = True
    path.write_text(json.dumps(model))
    print(json.dumps({"id":"fixture", "result":{"ok":True}}))
    raise SystemExit(0)
if not model.get("closed_transport") and cmd == ["pane", "get"]:
    pane = {"workspace_id":"w-fixture", "tab_id":"tab-fixture", "pane_id":"w-fixture:p1", "terminal_id":"terminal-fixture",
            "foreground_cwd":model.get("foreground_cwd"), "agent":model.get("provider", "codex"), "agent_status":model["agent_status"],
            "agent_session":{"agent":model.get("provider", "codex"), "kind":"id", "value":model["agent_session_id"]}}
    print(json.dumps({"id":"fixture", "result":{"pane":pane}}))
    raise SystemExit(0)
`

func TestAR03RetirementDoesNotReplayDeniedHerdrInspection(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	workflowTestModel(t, f.workflowFixture, map[string]any{"retire_inventory_denied": true})
	retireAsCoordinator(&f)
	f.D.HerdrTimeout = 350 * time.Millisecond
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AR03 fixture root", Reason: "Preserve denied inspection"})
	if err == nil {
		t.Fatal("denied inspection permitted retirement")
	}
	after, readErr := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !equal(before.Result.Delivery.Events, after.Result.Delivery.Events) {
		t.Fatal("blocked retirement changed history")
	}
	if c := retireQueueCurrent(t, f); c == nil || c.PreparationID != f.R.PreparationID {
		t.Fatal("blocked retirement released queue")
	}
	if _, e := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); e != nil {
		t.Fatal("blocked retirement removed source")
	}
	if n := workflowTestCalls(t, f.workflowFixture, "pane list"); n != 1 {
		t.Fatalf("retirement replayed denied Herdr inspection %d times; want one; error=%v", n, err)
	}
}

func retireFixture(t *testing.T) (deliveryFixture, WorkflowRun) {
	t.Helper()
	f := receiptFixture(t, "claude", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}})
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	f.D.HerdrTimeout = 3 * time.Second
	return f, o
}

func retireQueueCurrent(t *testing.T, f deliveryFixture) *workspace.QueueCurrent {
	t.Helper()
	a := f.R.Delivery.Agreement
	q, err := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: a.ProjectID, RepoID: a.RepoID, EpicID: a.EpicID}, false)
	if err != nil {
		t.Fatal(err)
	}
	return q.Current
}

func retireAsCoordinator(f *deliveryFixture) {
	f.D.CWD = func() (string, error) { return f.Parent, nil }
	f.D.Executable = func() (string, error) { return "/not-the-owner-control", nil }
}

func retireStopFixture(t *testing.T, f deliveryFixture, o WorkflowRun, effect string) WorkflowRun {
	t.Helper()
	proof := writeAny(t, f.R.WorkspaceRoot, "refusal-"+effect+".json", map[string]string{"source": "synthetic provider refusal", "effect_status": effect})
	stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "refusal-" + effect, Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "fixture coordinator", Reason: "Synthetic refused action", EvidencePath: proof, Action: "synthetic refused action", EffectStatus: effect})
	if err != nil {
		t.Fatal(err)
	}
	return stopped
}

// FS02 AC04/AC05: one coordinator operation ends the selected owner, proves no
// writer remains, closes the old owner's callbacks and frees the queue slot for
// a fresh Task copy. No final owner callback or manual state edit is needed.
func TestDeliveryRetireEndsSelectedOwnerAndFreesQueueWithoutOwnerCallback(t *testing.T) {
	for _, paneOpen := range []bool{true, false} {
		name := map[bool]string{true: "open_pane", false: "already_closed_with_stopped_history"}[paneOpen]
		t.Run(name, func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
			if !paneOpen {
				workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
			}
			if c := retireQueueCurrent(t, f); c == nil || c.PreparationID != f.R.PreparationID {
				t.Fatalf("fixture queue does not hold this Task: %+v", c)
			}
			owner := f.D
			retireAsCoordinator(&f)
			in := DeliveryRetireInput{Observer: "fixture root", Reason: "Selected cattle replacement after synthetic refusal"}
			out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if err != nil {
				t.Fatal(err)
			}
			if want := map[bool]int{true: 1, false: 0}[paneOpen]; workflowTestCalls(t, f.workflowFixture, "pane close") != want {
				t.Fatalf("pane close calls differ from %d", want)
			}
			if out.Delivery.Phase != "stopped" || out.NextAction.Actor != "coordinator" || !strings.Contains(out.NextAction.Message, "fresh Task") {
				t.Fatalf("retirement not visible: %+v %+v", out.Delivery.Phase, out.NextAction)
			}
			kinds := []string{}
			for _, e := range out.Delivery.Events {
				kinds = append(kinds, e.Kind)
			}
			if strings.Join(kinds, ",") != "coordinator_stop,owner_retirement" {
				t.Fatalf("history not preserved: %v", kinds)
			}
			if c := retireQueueCurrent(t, f); c != nil {
				t.Fatalf("queue slot still held: %+v", c)
			}
			if _, err := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
				t.Fatal("retirement removed the preserved Task worktree")
			}
			report := writeAny(t, f.R.WorkspaceRoot, "late.json", DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "late-owner", PreviousEventSHA256: out.Delivery.LastEventSHA256, Phase: "working", Summary: "Late", Meaning: "Late owner callback", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}})
			if _, err = WorkflowDeliveryReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, report); err == nil || !strings.Contains(err.Error(), "owner_retired") {
				t.Fatalf("retired owner callback was not closed: %v", err)
			}
			again, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if err != nil || len(again.Delivery.Events) != 2 {
				t.Fatalf("identical replay was not idempotent: %v", err)
			}
			in.Reason = "different reason"
			if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err == nil {
				t.Fatal("changed retirement overwrote the recorded one")
			}
		})
	}
}

func TestDeliveryRetireStopsOnActiveWriterOrUnknownEffect(t *testing.T) {
	cases := []struct {
		name, want string
		setup      func(*testing.T, deliveryFixture, WorkflowRun)
		closes     int
	}{
		{"pane_close_ignored", "owner_active", func(t *testing.T, f deliveryFixture, _ WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"close_ignored": true})
		}, 1},
		{"other_pane_writes_task", "owner_active", func(t *testing.T, f deliveryFixture, _ WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "t2", "pane_id": "p2", "terminal_id": "term2", "agent": "codex", "agent_status": "working", "foreground_cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath}}})
		}, 1},
		{"unknown_effect", "unknown", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			retireStopFixture(t, f, o, DeliveryEffectUnknown)
		}, 0},
		{"partial_effect", "partial", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			retireStopFixture(t, f, o, DeliveryEffectPartial)
		}, 0},
		{"attempt_in_progress", "active_writer", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
				s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "integration", State: "attempted"}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, o := retireFixture(t)
			tc.setup(t, f, o)
			before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			retireAsCoordinator(&f)
			if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Selected replacement"}); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("retirement did not stop with concrete %q diagnosis: %v", tc.want, err)
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != tc.closes {
				t.Fatal("dependent process effect started despite a known blocker")
			}
			after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if len(after.Result.Delivery.Events) != len(before.Result.Delivery.Events) {
				t.Fatal("blocked retirement recorded an event")
			}
			if c := retireQueueCurrent(t, f); c == nil || c.PreparationID != f.R.PreparationID {
				t.Fatal("blocked retirement released the queue slot")
			}
		})
	}
}

// Review fs02-r1/fs02-r2 regressions, adopted from the independent reviewer.
func TestFS02ReviewRetirePreservesEarlierUnknownEffect(t *testing.T) {
	f, o := retireFixture(t)
	retireStopFixture(t, f, o, DeliveryEffectUnknown)
	retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	retireAsCoordinator(&f)
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "independent reviewer fixture", Reason: "Selected isolated replacement"})
	if err == nil {
		t.Fatalf("retirement accepted unresolved earlier unknown effect: phase=%s queue_current=%v pane_close_calls=%d", out.Delivery.Phase, retireQueueCurrent(t, f), workflowTestCalls(t, f.workflowFixture, "pane close"))
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("closed pane despite earlier unresolved unknown effect")
	}
}

func TestFS02ReviewRetireRecoversDurableEventAfterStateFailure(t *testing.T) {
	f, o := retireFixture(t)
	retireAsCoordinator(&f)
	in := DeliveryRetireInput{Observer: "independent reviewer fixture", Reason: "Selected isolated replacement"}
	fault := f.D
	fault.Files = deliveryRejectStopState{statePath: o.Paths.RunRoot + "/state.json"}
	if _, err := WorkflowRetireDelivery(fault, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("fault not reached")
	}
	time.Sleep(10 * time.Millisecond)
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatalf("identical retirement cannot recover durable event after state publication failure: %v", err)
	}
}
