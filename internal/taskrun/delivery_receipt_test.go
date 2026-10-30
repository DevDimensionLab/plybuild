package taskrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Setting the additive protocol through JSON lets this regression reproduce the
// old behavior: the old reader ignores it and demands blanket permission.
func receiptFixture(t *testing.T, provider string, options ...deliveryFixtureOptions) deliveryFixture {
	t.Helper()
	f := newDeliveryFixture(t, provider, options...)
	goal := writeAny(t, f.R.WorkspaceRoot, "receipt-goal.json", deliveryReceiptGoal{TaskID: "task", Title: "Fixture task", Objective: "Observe the selected operation."})
	f.R.Delivery.Goal = FileBinding{goal, hashFileTest(t, goal)}
	raw, err := json.Marshal(f.R)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err = json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	fields["delivery"].(map[string]any)["startup_protocol"] = MandateReceiptV1
	draft := fields["handoff_draft"].(map[string]any)
	draft["delivery_binding"].(map[string]any)["startup_protocol"] = MandateReceiptV1
	raw, err = json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &f.R); err != nil {
		t.Fatal(err)
	}
	f.File = writeAny(t, f.R.WorkspaceRoot, "receipt-request.json", f.R)
	return f
}

func TestDeliveryReceiptDoesNotRequireBlanketPermission(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := receiptFixture(t, provider)
			o := workflowTestStart(t, f.workflowFixture)
			in := DeliveryBeginInput{Provider: provider, NativeSession: o.Transport.AgentSessionID}
			first, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
			if err != nil {
				t.Fatalf("neutral mandate receipt must reach authorized work without blanket permission: %v", err)
			}
			if first.Delivery.PermissionState != RuntimePerAction {
				t.Fatalf("receipt fabricated permission: %+v", first.Delivery)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(s.Acceptance.Locator)
			if err != nil {
				t.Fatal(err)
			}
			var a DeliveryAcceptance
			if err = json.Unmarshal(raw, &a); err != nil {
				t.Fatal(err)
			}
			if a.DeliveryPermission.PermissionConfirmed {
				t.Fatal("receipt asserted blanket permission")
			}
			var sandbox map[string]any
			if err = json.Unmarshal(a.Sandbox, &sandbox); err != nil {
				t.Fatal(err)
			}
			if sandbox["matches_contract"] != nil {
				t.Fatalf("unobserved sandbox coverage was attested: %s", a.Sandbox)
			}
			before := *s.Acceptance
			if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
				t.Fatal(err)
			}
			s, _ = workflowRead(f.R.WorkspaceRoot, o.RunID)
			if *s.Acceptance != before {
				t.Fatal("lost receipt response changed immutable evidence")
			}
		})
	}
}

var errReceiptSubprocessCloseout = errors.New("receipt closeout cannot inherit observers across a subprocess")

// Execute the same real closeout body in the journey and regression. The
// fixture's cleanupCLIAt guard observes any subprocess restoration inside it.

// All adapters are installed before startup. Ordinary non-adopted startup does
// not consume SourceWriters or StartupProcessTree; their zero counts are not
// behavioral coverage of adoption or recovery. LocalProcesses is exercised only
// after modeled owner exit, through production WorkflowFollow/closeout.
type receiptObservers struct {
	startup, source, local int
	subprocess             int
	session, transcript    string
	table                  []LocalProcess
	err                    error
}

func (s *receiptObservers) assertUnused(t *testing.T) {
	t.Helper()
	if s.startup != 0 || s.source != 0 {
		t.Fatalf("ordinary journey consumed recovery/adoption observers: startup=%d source=%d", s.startup, s.source)
	}
	t.Logf("defensive observers unexercised: StartupProcessTree=%d SourceWriters=%d", s.startup, s.source)
}

func receiptCLICall(t *testing.T, binary string, f deliveryFixture, args ...string) (WorkflowRun, error) {
	t.Helper()
	if len(args) < 2 || args[0] != "execute" || args[1] != "check-in" && args[1] != "begin" && args[1] != "verify" {
		t.Fatalf("unmodeled receipt subprocess boundary: %q", args)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir, c.WaitDelay = f.Prepared.Preparation.Preparation.Plan.WorktreePath, time.Second
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run() // Exact owned child, synchronously waited; never a detached supervisor.
	t.Logf("receipt CLI argv=%q cwd=%q stdout=%s stderr=%s exit=%v", c.Args, c.Dir, stdout.Bytes(), stderr.Bytes(), err)
	var o WorkflowRun
	if decodeErr := json.Unmarshal(stdout.Bytes(), &o); decodeErr != nil && err == nil {
		err = decodeErr
	}
	return o, err
}

// Public Delivery register/check/execute and subprocess follow remain covered by
// their own tests. This receipt-local closure uses production functions with
// simulated inventories and real native Git, integration, queue/base records.

func TestDeliveryReceiptReportedRootsNeverAttestSandbox(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID, ReadRoots: []string{f.R.WorkspaceRoot}, WriteRoots: []string{f.Prepared.Preparation.Preparation.Plan.WorktreePath}, Profile: "auto", Policy: "Source: synthetic recipient observation; future decisions unknown"}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	b, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	var a DeliveryAcceptance
	if err = json.Unmarshal(b, &a); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(a.Sandbox), `"matches_contract":null`) || a.DeliveryPermission.PermissionConfirmed {
		t.Fatalf("reported roots became authority: %s", b)
	}
}

func TestDeliveryReceiptLaterRuntimeRefusalPreservesStop(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	proof := writeAny(t, f.R.WorkspaceRoot, "actual-denial.json", map[string]string{"source": "synthetic provider classifier", "reason": "action denied", "effect": "none observed"})
	var stop DeliveryStopInput
	raw, _ := json.Marshal(map[string]any{"event_id": "later-refusal", "provider": "claude", "native_session": o.Transport.AgentSessionID, "stage": "work", "observer": "fixture coordinator", "reason": "Provider refused the relevant test action", "evidence_path": proof, "action": "run relevant test", "effect_status": "none_observed", "rule": "synthetic-classifier-rule"})
	if err := json.Unmarshal(raw, &stop); err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
	if err != nil {
		t.Fatalf("later action refusal must preserve phase, source and effects: %v", err)
	}
	if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
		t.Fatal("receipt cleared a known refusal")
	}
	b, err := workflowBound(*shown.Delivery.ActiveStop, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"action":"run relevant test"`) || !strings.Contains(string(b), `"effect_status":"none_observed"`) || !strings.Contains(string(b), `"rule":"synthetic-classifier-rule"`) {
		t.Fatalf("refusal lost action or effects: %s", b)
	}
	if shown.DeliveryStatus.Startup.RuntimeDecision != "preserved_stop" {
		t.Fatal("readback hid runtime refusal")
	}
}

func TestDeliveryReceiptDoesNotDowngradeLegacyContract(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err == nil {
		t.Fatal("check-in downgraded frozen legacy startup")
	}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
		t.Fatal("legacy permission requirement changed")
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" {
		t.Fatal("legacy authority/history changed")
	}
}

// A runtime refusal may follow ordinary work outside native Delivery.Attempt.
// Permission to continue cannot establish the outcome of that earlier action.
func TestDeliveryReceiptResolutionPreservesUncertainExternalEffects(t *testing.T) {
	for _, effect := range []string{"unknown", "partial"} {
		t.Run(effect, func(t *testing.T) {
			f := receiptFixture(t, "claude")
			o := workflowTestStart(t, f.workflowFixture)
			in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}
			if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
				t.Fatal(err)
			}
			proof := writeAny(t, f.R.WorkspaceRoot, "external-work-refusal.json", map[string]string{"source": "synthetic provider refusal", "effect_status": effect})
			stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "external-refusal", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "synthetic coordinator", Reason: "Tool response lost after ordinary work; earlier effects unresolved", EvidencePath: proof, Action: "ordinary work outside the native delivery attempt", EffectStatus: effect})
			if err != nil {
				t.Fatal(err)
			}
			stop := *stopped.Delivery.ActiveStop
			if stopped.Delivery.Attempt != nil {
				t.Fatal("fixture unexpectedly has a native delivery attempt")
			}
			for attempt := 0; attempt < 2; attempt++ {
				resumed, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stop.SHA256, "Current refusal cause repaired; earlier external effects remain unresolved", "Continue the clarified changed fixture action")
				if err == nil {
					t.Fatalf("permission and text cleared %s external effects: phase=%s stop=%v", effect, resumed.Delivery.Phase, resumed.Delivery.ActiveStop)
				}
				observed, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				if observed.Delivery.ActiveStop == nil || *observed.Delivery.ActiveStop != stop || observed.Delivery.Phase != "stopped" || len(observed.Delivery.Events) != 1 {
					t.Fatal("failed resolution changed the stop or its history")
				}
				if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
					t.Fatal("receipt reopened unresolved external work")
				}
			}
		})
	}
}

func TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint(t *testing.T) {
	for _, receipt := range []bool{false, true} {
		name := "legacy"
		if receipt {
			name = "mandate_receipt"
		}
		t.Run(name, func(t *testing.T) {
			var f deliveryFixture
			if receipt {
				f = receiptFixture(t, "codex")
			} else {
				f = newDeliveryFixture(t, "codex")
			}
			o := workflowTestStart(t, f.workflowFixture)
			proof := writeAny(t, f.R.WorkspaceRoot, "pre-receipt-denial.json", map[string]string{"source": "synthetic provider refusal", "effect_status": "not_started"})
			stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "before-receipt", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "acceptance", Observer: "fixture coordinator", Reason: "Initial receipt action refused before starting", EvidencePath: proof, Action: "record initial receipt", EffectStatus: DeliveryEffectNotStarted})
			if err != nil {
				t.Fatal(err)
			}
			want := "execute begin"
			if receipt {
				want = "execute check-in"
			}
			if !strings.Contains(stopped.NextAction.Message, strings.TrimPrefix(want, "execute ")+" is still required") {
				t.Fatalf("stopped %s contract recommends wrong startup entrypoint: %s", name, stopped.NextAction.Message)
			}
			resumed, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stopped.Delivery.ActiveStop.SHA256, "Actual refusal cause repaired; initial receipt never started", "Continue the clarified changed fixture action")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(resumed.NextAction.Message, want) {
				t.Fatalf("resolved %s contract recommends wrong startup entrypoint: %s", name, resumed.NextAction.Message)
			}
			in := beginFixtureInput(f, o)
			if receipt {
				in = DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}
			}
			if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, !receipt); err == nil {
				t.Fatal("recovery allowed the other protocol's startup entrypoint")
			}
			received, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, receipt)
			if err != nil {
				t.Fatalf("recommended frozen entrypoint did not reach receipt: %v", err)
			}
			if received.Delivery.Phase != "working" || received.Delivery.ActiveStop != nil {
				t.Fatal("correct startup failed to resume the same mandate")
			}
		})
	}
}
