package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDeliveryResolvingLatestStopDoesNotForgetEarlierDenial(t *testing.T) {
	f := newDeliveryFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic denials"})
	in := DeliveryStopInput{EventID: "acceptance-denied", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "acceptance", Observer: "coordinator", Reason: "Rejected acceptance", EvidencePath: proof}
	first, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	in.EventID = "report-denied"
	in.Stage = "report"
	in.Reason = "Rejected stop report"
	second, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, second.Delivery.ActiveStop.SHA256, "Report cause repaired with actual permission", "Continue the clarified changed fixture action")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Delivery.ActiveStop == nil || *resolved.Delivery.ActiveStop != *first.Delivery.ActiveStop || resolved.Delivery.Phase != "stopped" {
		t.Fatal("resolving a later report denial hid unresolved acceptance denial")
	}
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, beginFixtureInput(f, o)); err == nil {
		t.Fatal("earlier denial no longer blocked acceptance")
	}
	resolved, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, first.Delivery.ActiveStop.SHA256, "Acceptance cause repaired with actual permission", "Continue the clarified changed fixture action")
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Delivery.ActiveStop != nil {
		t.Fatal("fully resolved stop history remained active")
	}
}

func TestDeliveryObservedStopKeepsLegacyStateReadable(t *testing.T) {
	f := newDeliveryFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic acceptance rejection"})
	_, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "denial", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "acceptance", Observer: "coordinator", Reason: "Rejected callback", EvidencePath: proof})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"active_stop"`) {
		t.Fatal("derived stop added a field that older frozen controls reject")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Delivery.ActiveStop == nil {
		t.Fatal("retained events did not reconstruct active stop")
	}
}

func TestDeliveryObservedStopRequiresOwnerResolutionBeforeProgress(t *testing.T) {
	for _, accepted := range []bool{false, true} {
		f := newDeliveryFixture(t, "codex")
		o := workflowTestStart(t, f.workflowFixture)
		in := beginFixtureInput(f, o)
		if accepted {
			var err error
			o, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
			if err != nil {
				t.Fatal(err)
			}
		}
		proof := writeAny(t, f.R.WorkspaceRoot, "denial.json", map[string]string{"observed": "Synthetic callback denial"})
		stop := DeliveryStopInput{EventID: "denial", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "coordinator", Reason: "Observed callback denial", EvidencePath: proof}
		shown, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Delivery.Phase != "stopped" {
			t.Fatal("stale working state hid stop")
		}
		digest := shown.Delivery.ActiveStop.SHA256
		if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
			t.Fatal("begin bypassed known stop")
		}
		if _, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, digest, "unknown current permission", ""); err == nil {
			t.Fatal("unknown permission resolved stop")
		}
		if _, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "sha256:wrong", "cause repaired", "Continue the clarified changed fixture action"); err == nil {
			t.Fatal("wrong stop resolved")
		}
		shown, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, digest, "Fixture cause repaired; current synthetic policy covers the same isolated scope", "Continue the clarified changed fixture action")
		if err != nil {
			t.Fatal(err)
		}
		if shown.Delivery.ActiveStop != nil {
			t.Fatal("resolved stop stayed active")
		}
		if !accepted && shown.Delivery.PermissionState != "pending_runtime_acceptance" {
			t.Fatal("resolution granted acceptance")
		}
		shown, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Delivery.Phase != "working" {
			t.Fatal("supported continuation made no progress")
		}
		shown, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
		if err != nil {
			t.Fatal(err)
		}
		if shown.Delivery.ActiveStop != nil || len(shown.Delivery.Events) != 2 {
			t.Fatal("replayed observation reactivated or lost the resolved history")
		}
	}
}

func TestDeliveryCoordinatorPreservesDeniedCallbacksWithoutAcceptance(t *testing.T) {
	for _, stage := range []string{"acceptance", "report"} {
		t.Run(stage, func(t *testing.T) {
			f := newDeliveryFixture(t, "claude")
			o := workflowTestStart(t, f.workflowFixture)
			source := writeAny(t, f.R.WorkspaceRoot, "actual-denial.json", map[string]string{"observation": "Synthetic separate provider denial of " + stage})
			in := DeliveryStopInput{EventID: "denied-" + stage, Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: stage, Observer: "fixture coordinator", Reason: "Provider denied callback; owner has not accepted", EvidencePath: source}
			// The coordinator is outside the Task and is not the bound control.
			f.D.CWD = func() (string, error) { return f.Parent, nil }
			f.D.Executable = func() (string, error) { return "/not-the-owner-control", nil }
			shown, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if err != nil {
				t.Fatal(err)
			}
			if shown.Round.State != "stopped" || shown.NextAction.Actor != "coordinator" {
				t.Fatalf("stop hidden: %+v", shown)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" {
				t.Fatal("observer impersonated owner acceptance")
			}
			before := s.Result.Delivery.LastEventSHA256
			if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
				t.Fatal(err)
			}
			s, _ = workflowRead(f.R.WorkspaceRoot, o.RunID)
			if !equal(before, s.Result.Delivery.LastEventSHA256) || len(s.Result.Delivery.Events) != 1 {
				t.Fatal("repeated observation duplicated the event")
			}
			if err = workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); err == nil {
				t.Fatal("stop falsely released writer")
			}
			in.Reason = "changed denial"
			if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err == nil {
				t.Fatal("denial was overwritten")
			}
		})
	}
}

// FS02: the owner names the clarified changed action instead of attesting
// future runtime permission. The refused action itself stays refused.
func TestDeliveryResolveStopContinuesNamedChangedActionWithoutPermissionClaim(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err != nil {
		t.Fatal(err)
	}
	denied := "ply workflow run review wfr_old --candidate abc"
	proof := writeAny(t, f.R.WorkspaceRoot, "work-refusal.json", map[string]string{"source": "synthetic provider refusal", "rule": "Auto-Mode Bypass"})
	stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "work-refusal", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "fixture coordinator", Reason: "Review command refused before execution", EvidencePath: proof, Action: denied, EffectStatus: DeliveryEffectNotStarted, Rule: "Auto-Mode Bypass"})
	if err != nil {
		t.Fatal(err)
	}
	stop := *stopped.Delivery.ActiveStop
	for _, action := range []string{"", "  ply workflow run review   wfr_old --candidate abc "} {
		if _, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stop.SHA256, "Root selected a separate review", action); err == nil {
			t.Fatalf("resolution with action %q cleared a refusal without a changed action", action)
		}
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if shown.Delivery.ActiveStop == nil || *shown.Delivery.ActiveStop != stop || len(shown.Delivery.Events) != 1 {
		t.Fatal("rejected resolution changed the preserved refusal")
	}
	changed := "Continue product correction in this Task; root owns the separate review"
	resolved, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stop.SHA256, "Root selected a separate review under its standing authority", changed)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Delivery.ActiveStop != nil || resolved.Delivery.Phase != "working" || len(resolved.Delivery.Events) != 2 {
		t.Fatalf("named continuation made no progress: %+v", resolved.Delivery)
	}
	var record DeliveryStopResolution
	raw, err := workflowBound(resolved.Delivery.Events[1].Binding, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if err = decode(raw, 1<<20, &record); err != nil {
		t.Fatal(err)
	}
	if record.PermissionConfirmed || record.ContinuationAction != changed || record.Stop != stop || record.RequestSHA256 != o.RequestSHA256 {
		t.Fatalf("resolution did not preserve the named action, mandate and refusal: %+v", record)
	}
	if strings.Contains(string(raw), "permission_confirmed\":true") {
		t.Fatal("resolution recorded a self-attested permission")
	}
	if _, err = WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stop.SHA256, "Root selected a separate review under its standing authority", changed); err != nil {
		t.Fatalf("identical replay after a lost reply failed: %v", err)
	}
	if s, err := workflowRead(f.R.WorkspaceRoot, o.RunID); err != nil || s.Result.Delivery.ActiveStop != nil || len(s.Result.Delivery.Events) != 2 {
		t.Fatalf("reload of resolved history failed: %v", err)
	}
}

// FS02 AC06: observing the same bound session is not a Task mutation. A
// transient read cwd outside the Task must not hide that it still works;
// owner callbacks keep their own fresh cwd check.
func TestDeliveryFollowObservesSameSessionDuringTransientReadCWD(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}, true); err != nil {
		t.Fatal(err)
	}
	before, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	outside := f.Parent
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "working", "foreground_cwd": outside})
	time.Sleep(time.Millisecond)
	shown, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 300*time.Millisecond)
	if err == nil || strings.Contains(err.Error(), "cwd") {
		t.Fatalf("follow treated a transient read cwd as a binding failure: %v", err)
	}
	if shown.Transport.ObservedAt == before.Transport.ObservedAt || shown.Transport.State != "working" || shown.Delivery.Phase != "working" {
		t.Fatalf("follow hid the still-working session: %+v", shown.Transport)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = deliveryLiveOwner(f.D, s); err == nil || !strings.Contains(err.Error(), "cwd") {
		t.Fatalf("owner callback check accepted a live session outside the Task: %v", err)
	}
}
