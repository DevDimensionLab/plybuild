package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type deliveryRejectStopState struct {
	systemFileSystem
	statePath string
}

func (f deliveryRejectStopState) Replace(path string, b []byte) error {
	if path == f.statePath {
		return errors.New("review injection: state publication failed after durable stop event")
	}
	return f.systemFileSystem.Replace(path, b)
}

func TestDeliveryObservedStopInterruptedPublicationRetainsBlock(t *testing.T) {
	for _, accepted := range []bool{false, true} {
		name := "before_acceptance"
		if accepted {
			name = "after_acceptance"
		}
		t.Run(name, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := workflowTestStart(t, f.workflowFixture)
			in := beginFixtureInput(f, o)
			var err error
			if accepted {
				o, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
				if err != nil {
					t.Fatal(err)
				}
			}
			proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic provider rejected callback transport"})
			stage := "acceptance"
			if accepted {
				stage = "report"
			}
			stop := DeliveryStopInput{EventID: "review-denial", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: stage, Observer: "review fixture", Reason: "Observed callback denial", EvidencePath: proof}
			fault := f.D
			instant := time.Now().UTC()
			fault.Now = func() time.Time { return instant }
			fault.Files = deliveryRejectStopState{statePath: filepath.Join(o.Paths.RunRoot, "state.json")}
			if _, err = WorkflowObserveDeliveryStop(fault, f.R.WorkspaceRoot, o.RunID, stop); err == nil {
				t.Fatal("injected state write failure did not occur")
			}
			if _, err = os.Stat(filepath.Join(o.Paths.RunRoot, "delivery", "events", "00000001.json")); err != nil {
				t.Fatalf("expected durable stop event: %v", err)
			}
			eventPath := filepath.Join(o.Paths.RunRoot, "delivery", "events", "00000001.json")
			original, err := os.ReadFile(eventPath)
			if err != nil {
				t.Fatal(err)
			}
			statePath := filepath.Join(o.Paths.RunRoot, "state.json")
			originalState, err := os.ReadFile(statePath)
			if err != nil {
				t.Fatal(err)
			}
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("after interrupted stop: phase=%s active_stop=%v", shown.Delivery.Phase, shown.Delivery.ActiveStop)
			if shown.Delivery.ActiveStop == nil {
				t.Error("persisted denial is invisible after state-write failure")
			}
			begun, beginErr := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
			if beginErr == nil {
				t.Errorf("begin ignored durable unresolved denial: phase=%s permission=%s", begun.Delivery.Phase, begun.Delivery.PermissionState)
			}
			unchangedState, err := os.ReadFile(statePath)
			if err != nil || string(unchangedState) != string(originalState) {
				t.Fatal("readback or blocked begin wrote recovered state")
			}
			f.D.Now = func() time.Time { return instant.Add(time.Second) }
			recovered, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
			if err != nil {
				t.Fatalf("identical observation cannot recover its durable event: %v", err)
			}
			after, err := os.ReadFile(eventPath)
			if err != nil || string(after) != string(original) {
				t.Fatal("recovery changed original denial or timestamp")
			}
			if recovered.Delivery.ActiveStop == nil || len(recovered.Delivery.Events) != 1 {
				t.Fatal("recovery lost or duplicated the denial")
			}
			changed := stop
			changed.Reason = "different reason"
			if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, changed); err == nil {
				t.Fatal("changed replay replaced denial")
			}
			stop.EventID = "later-denial"
			later, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
			if err != nil || len(later.Delivery.Events) != 2 {
				t.Fatalf("recovered index cannot preserve a later denial: %v", err)
			}
		})
	}
}

func TestDeliveryObservedStopResolutionReplayWithPendingDenials(t *testing.T) {
	f := newDeliveryFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic separate callback denials"})
	stop := DeliveryStopInput{EventID: "first-denial", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "acceptance", Observer: "review fixture", Reason: "Observed acceptance denial", EvidencePath: proof}
	first, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
	if err != nil {
		t.Fatal(err)
	}
	stop.EventID = "second-denial"
	stop.Stage = "report"
	stop.Reason = "Observed report denial"
	second, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
	if err != nil {
		t.Fatal(err)
	}
	sha := second.Delivery.ActiveStop.SHA256
	resolution := "Synthetic report transport cause repaired; actual permission confirmed"
	resolved, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, sha, resolution, "Continue the clarified changed fixture action")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Delivery.ActiveStop == nil || *resolved.Delivery.ActiveStop != *first.Delivery.ActiveStop {
		t.Fatal("earlier denial disappeared")
	}
	replay, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, sha, resolution, "Continue the clarified changed fixture action")
	if err != nil {
		t.Fatalf("lost resolution response cannot replay while earlier denial remains: %v", err)
	}
	if replay.Delivery.ActiveStop == nil || *replay.Delivery.ActiveStop != *first.Delivery.ActiveStop || len(replay.Delivery.Events) != 3 {
		t.Fatal("resolution replay changed earlier denial or history")
	}
	stop.EventID = "third-denial"
	third, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
	if err != nil {
		t.Fatal(err)
	}
	replay, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, sha, resolution, "Continue the clarified changed fixture action")
	if err != nil {
		t.Fatalf("recorded resolution cannot replay with a later pending denial: %v", err)
	}
	if replay.Delivery.ActiveStop == nil || *replay.Delivery.ActiveStop != *third.Delivery.ActiveStop || len(replay.Delivery.Events) != 4 {
		t.Fatal("resolution replay changed later denial or history")
	}
	if _, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, sha, "different repair claim", "Continue the clarified changed fixture action"); err == nil {
		t.Fatal("conflicting resolution replay accepted")
	}
}
