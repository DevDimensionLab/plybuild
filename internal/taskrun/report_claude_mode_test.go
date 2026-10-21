package taskrun

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// A new ordinary Claude report Run with the selected native permission mode.
func al06ClaudeReportFixture(t *testing.T, mode, binary string) workflowFixture {
	t.Helper()
	f := workflowProviderFixture(t, "claude")
	f.R.Runtime.PermissionBinding.ProfileID = mode
	f.R.RequestKey += "-" + mode
	return al01ReceiptFixtureFrom(t, f, binary)
}

func al06AssertClaudeStart(t *testing.T, f workflowFixture, o WorkflowRun, mode string) {
	t.Helper()
	starts := workflowProviderCalls(t, f, "agent start")
	want := []string{"agent", "start", "ply-" + o.RunID[4:32], "--kind", "claude", "--pane", "w-fixture:p1", "--timeout", "30000", "--", "--model", f.R.Runtime.Model, "--permission-mode", mode}
	if len(starts) != 1 || !reflect.DeepEqual(starts[0], want) {
		t.Fatalf("launch argv = %q; want one %q", starts, want)
	}
}

func TestAL06ClaudeReportModeReceiptCorrectionAndReturn(t *testing.T) {
	for _, mode := range []string{ClaudePermissionManual, ClaudePermissionAuto} {
		t.Run(mode, func(t *testing.T) {
			f := al06ClaudeReportFixture(t, mode, "")
			if f.R.StartupProtocol != MandateReceiptV1 || f.R.Runtime.PermissionBinding.ProfileID != mode {
				t.Fatalf("preserved request lost the selected mode or receipt: %+v", f.R.Runtime.PermissionBinding)
			}
			v, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			p := v.(WorkflowPreview)
			if p.Effects[1] != "Create one background Herdr tab and one Claude session in "+mode+" permission mode" || p.RequestSHA256 != digest(f.R) {
				t.Fatalf("preview does not show the selected mode: %q", p.Effects[1])
			}
			o, err := WorkflowStart(f.D, f.File, *p.Confirmation)
			if err != nil {
				t.Fatal(err)
			}
			al06AssertClaudeStart(t, f, o, mode)

			in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID, Profile: mode}
			for i := 0; i < 2; i++ {
				if o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err != nil {
					t.Fatal(err)
				}
			}
			// The first receipt is already preserved: a later different label
			// cannot replace that observation. This is immutability, not a
			// requested==actual attestation gate.
			other := map[string]string{"manual": "auto", "auto": "manual"}[mode]
			if _, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "claude", NativeSession: in.NativeSession, Profile: other}, true); err == nil {
				t.Fatal("a later receipt with another label replaced the preserved first observation")
			}
			s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
			raw, _ := workflowBound(*s.Acceptance, 1<<20)
			var a Acceptance
			json.Unmarshal(raw, &a)
			if a.RuntimeClaim.ProfileID == nil || *a.RuntimeClaim.ProfileID != mode || s.Request.Runtime.PermissionBinding.ProfileID != mode || o.Delivery != nil {
				t.Fatalf("readback lost the preserved observation or frozen selected mode: %s", raw)
			}

			o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
			o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
			if err != nil || o.Round.Number != 1 || o.Transport.AgentSessionID != in.NativeSession {
				t.Fatalf("correction lost the selected session: %+v %v", o, err)
			}
			o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
			o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
			if err != nil || o.FinalReturn.State != "accepted" || o.FinalReturn.TerminalSHA256 == nil || o.TaskResultState != "not_published" {
				t.Fatalf("report return did not seal: %+v %v", o, err)
			}
			if workflowTestCalls(t, f, "tab create") != 1 || workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 2 {
				t.Fatal("report journey started another tab/session or sent the wrong prompts")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestAL06ClaudeAutoLostStartAndPreservedRefusal(t *testing.T) {
	f := al06ClaudeReportFixture(t, ClaudePermissionAuto, "")
	f.D.HerdrTimeout = 2 * time.Second
	workflowTestModel(t, f, map[string]any{"lost_command": "agent start"})
	o, err := workflowReadinessStart(t, f)
	if err != nil || o.Transport.AgentSessionID != "fixture-session" {
		t.Fatalf("lost start reply did not follow the same session: %+v %v", o, err)
	}
	al06AssertClaudeStart(t, f, o, ClaudePermissionAuto)
	again, err := WorkflowPreviewStart(f.D, f.File)
	if existing, ok := again.(WorkflowRun); err != nil || !ok || existing.Transport.AgentSessionID != o.Transport.AgentSessionID {
		t.Fatalf("repeated start did not return the same reserved session: %+v %v", again, err)
	}
	if workflowTestCalls(t, f, "tab create") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
		t.Fatal("lost reply or repeated start opened a duplicate tab or session")
	}

	evidence := writeAny(t, f.R.WorkspaceRoot, "denied-auto-receipt.json", map[string]any{"source": "synthetic Claude auto runtime", "action": "receipt callback", "result": "denied"})
	stop := DeliveryStopInput{EventID: "denied-auto-receipt", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "receipt", Observer: "fixture coordinator", Reason: "Claude auto refused the receipt callback", Action: "receipt callback", EffectStatus: DeliveryEffectNotStarted, EvidencePath: evidence}
	if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop); err != nil {
		t.Fatal(err)
	}
	for _, profile := range []string{"", "auto", "manual"} {
		in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID, Profile: profile}
		if _, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err == nil {
			t.Fatalf("receipt with profile %q erased the preserved refusal", profile)
		}
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if s.Acceptance != nil || s.ReportProviderStop == nil || workflowTestCalls(t, f, "agent prompt") != 1 {
		t.Fatal("preserved refusal gained acceptance or new effects")
	}
}

func TestAL06ClaudeReportUnsupportedModeFailsBeforeEffects(t *testing.T) {
	for _, mode := range []string{"other", "Auto", "default", "bypassPermissions"} {
		t.Run(mode, func(t *testing.T) {
			f := workflowProviderFixture(t, "claude")
			f.R.Runtime.PermissionBinding.ProfileID = mode
			workflowProviderRebind(t, &f)
			_, err := WorkflowPreviewStart(f.D, f.File)
			if err == nil || !strings.Contains(err.Error(), "unsupported Claude permission mode") || !strings.Contains(err.Error(), "manual or auto") {
				t.Fatalf("unsupported mode lacks a precise supported way forward: %v", err)
			}
			if _, err = WorkflowStart(f.D, f.File, "invalid"); err == nil {
				t.Fatal("unsupported Claude mode started")
			}
			workflowProviderNoReservation(t, f)
		})
	}
}

func TestAL06ClaudeAutoReportCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	f := al06ClaudeReportFixture(t, ClaudePermissionAuto, binary)
	cwd, _ := f.D.CWD()
	p := al01PreviewCLI(t, binary, cwd, f.File)
	if p.Provider != "claude" || p.Effects[1] != "Create one background Herdr tab and one Claude session in auto permission mode" {
		t.Fatalf("candidate CLI preview does not show auto: %+v", p.Effects)
	}
	o := workflowTestStart(t, f)
	al06AssertClaudeStart(t, f, o, ClaudePermissionAuto)
	args := []string{"run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "claude", "--native-session", o.Transport.AgentSessionID, "--profile", "auto"}
	for i := 0; i < 2; i++ {
		o = al01CLIOK(t, binary, f, args...)
	}
	if o.Delivery != nil || o.ReportStartup == nil || o.ReportStartup.Receipt == nil || o.ReportStartup.Protocol != MandateReceiptV1 {
		t.Fatalf("candidate CLI receipt lost the report contract: %+v", o.ReportStartup)
	}
	t.Log("Candidate CLI previewed, and received an auto Claude report mandate in the same session")
}

// A launch selection is not an attestation of the runtime's observed policy.
func TestAL06ClaudeReceiptPreservesUnknownOrDifferentObservation(t *testing.T) {
	for _, profile := range []string{"", ClaudePermissionManual} {
		t.Run("observed-"+profile, func(t *testing.T) {
			f := al06ClaudeReportFixture(t, ClaudePermissionAuto, "")
			o := workflowTestStart(t, f)
			in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID, Profile: profile}
			for i := 0; i < 2; i++ {
				if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err != nil {
					t.Fatal(err)
				}
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := workflowBound(*s.Acceptance, 1<<20)
			if err != nil {
				t.Fatal(err)
			}
			var a Acceptance
			if err = json.Unmarshal(raw, &a); err != nil {
				t.Fatal(err)
			}
			if profile == "" && a.RuntimeClaim.ProfileID != nil || profile != "" && (a.RuntimeClaim.ProfileID == nil || *a.RuntimeClaim.ProfileID != profile) {
				t.Fatalf("receipt invented or changed profile observation: %s", raw)
			}
			if s.Request.Runtime.PermissionBinding.ProfileID != ClaudePermissionAuto {
				t.Fatal("receipt changed frozen launch selection")
			}
			changed := in
			changed.Profile = ClaudePermissionAuto
			if _, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, changed, true); err == nil {
				t.Fatal("later receipt rewrote initial observation")
			}
			al06AssertClaudeStart(t, f, o, ClaudePermissionAuto)
		})
	}
}
