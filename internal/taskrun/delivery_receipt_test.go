package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Setting the additive protocol through JSON lets this regression reproduce the
// old behavior: the old reader ignores it and demands blanket permission.
func receiptFixture(t *testing.T, provider string, options ...deliveryFixtureOptions) deliveryFixture {
	t.Helper()
	f := newDeliveryFixture(t, provider, options...)
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
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

func TestDeliveryReceiptCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic receipt journey", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	f.File = writeAny(t, f.R.WorkspaceRoot, "receipt-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	args := []string{"execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID}
	for _, bad := range []struct {
		index int
		value string
	}{{4, o.Paths.Context + "wrong"}, {6, "claude"}, {8, "wrong-session"}} {
		invalid := append([]string{}, args...)
		invalid[bad.index] = bad.value
		if _, _, err := automaticCLI(binary, f, invalid...); err == nil {
			t.Fatal("invalid receipt accepted")
		}
	}
	o = automaticCLIOK(t, binary, f, args...)
	o = automaticCLIOK(t, binary, f, args...)
	if o.DeliveryStatus.Startup.SandboxCoverage != "unknown" || o.DeliveryStatus.Startup.RuntimeDecision != "per_action_not_observed" {
		t.Fatalf("receipt fabricated runtime facts: %+v", o.DeliveryStatus.Startup)
	}
	if _, _, err := automaticCLI(binary, f, "execute", "begin", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--policy", "unknown", "--read-root", f.R.WorkspaceRoot, "--write-root", f.R.WorkspaceRoot, "--permission-confirmed"); err == nil {
		t.Fatal("legacy entrypoint accepted new protocol")
	}
	deliveryTestCandidate(t, f, "receipt candidate")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'receipt candidate'\n\"$PLY_CANDIDATE_BINARY\" workflow execute check-in --help >/dev/null\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "receipt-journey")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
	if o.Delivery.Phase != "automatic_acceptance_passed" {
		t.Fatalf("receipt did not reach candidate qualification: %+v", o.Delivery)
	}
	assertAutomaticCLIClosure(t, binary, f, o)
}

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

func TestDeliveryReceiptGuideAndVisibleMandate(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	prompt := deliveryInstructions(s)
	if !strings.Contains(prompt, "Your authorized assignment is") || !strings.Contains(prompt, "workflow execute check-in") || strings.Contains(prompt, "permission-confirmed") {
		t.Fatalf("visible task omitted the neutral mandate: %s", prompt)
	}
	b, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "delivery", "callback-guide.md"))
	if err != nil {
		t.Fatal(err)
	}
	guide := string(b)
	if !strings.Contains(guide, "mandate_receipt_v1") || !strings.Contains(guide, "matches_contract=null") || strings.Contains(guide, "## Advanced acceptance file contract") || strings.Contains(guide, "--permission-confirmed") {
		t.Fatalf("guide routes new contract to legacy self-attestation: %s", guide)
	}
	if _, err = os.Stat(filepath.Join(o.Paths.RunRoot, "delivery", "templates", "acceptance.json")); !os.IsNotExist(err) {
		t.Fatal("new flow still supplies handwritten legacy acceptance template")
	}
}

func TestDeliveryReceiptAuthorityRemainsBoundToMandateAndEffects(t *testing.T) {
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic authority test", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	o := workflowTestStart(t, f.workflowFixture)
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	auth, err := deliveryAuthorityState(f.D, s)
	if err != nil {
		t.Fatal(err)
	}
	if auth.AuthorityKind != workspace.DeliveryMandateReceiptV1 || auth.PermissionConfirmed {
		t.Fatalf("native human mandate became permission self-attestation: %+v", auth)
	}
	var accepted DeliveryAcceptance
	b, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &accepted); err != nil {
		t.Fatal(err)
	}
	for _, change := range []string{"effects", "target", "protocol", "permission"} {
		t.Run(change, func(t *testing.T) {
			raw, _ := json.Marshal(accepted)
			var bad DeliveryAcceptance
			_ = json.Unmarshal(raw, &bad)
			switch change {
			case "effects":
				bad.DeliveryPermission.AllowedEffects = append(bad.DeliveryPermission.AllowedEffects, "push")
			case "target":
				bad.DeliveryPermission.DeliveryAgreementSHA256 = hash([]byte("foreign target"))
			case "protocol":
				bad.DeliveryPermission.StartupProtocol = ""
			case "permission":
				bad.DeliveryPermission.PermissionConfirmed = true
			}
			if err := positiveDeliveryClaim(bad.Acceptance, bad.DeliveryPermission, s.Request); err == nil {
				t.Fatal("receipt accepted expanded or changed authority")
			}
		})
	}
	altered := s
	altered.Result.Handoff.SHA256 = hash([]byte("different mandate"))
	if _, err = deliveryAuthorityState(f.D, altered); err == nil {
		t.Fatal("receipt accepted a changed mandate")
	}
	priorCWD := f.D.CWD
	f.D.CWD = func() (string, error) { return f.Parent, nil }
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}); err == nil {
		t.Fatal("receipt accepted another cwd")
	}
	f.D.CWD = priorCWD
}
