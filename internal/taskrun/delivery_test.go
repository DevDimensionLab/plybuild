package taskrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func TestDeliveryModelAliasMatchesOnlySelectedClaudeFamily(t *testing.T) {
	for _, tc := range []struct {
		provider, requested, actual string
		match                       bool
	}{
		{"claude", "opus", "opus", true}, {"claude", "opus", "claude-opus-4-6", true},
		{"claude", "sonnet", "claude-sonnet-4-5-20250929", true}, {"claude", "haiku", "claude-haiku-4-5", true},
		{"claude", "opus", "claude-sonnet-4-5", false}, {"claude", "opus", "claude-opus-", false},
		{"claude", "opus", "claude-opus-latest", false}, {"claude", "opus", "claude-opus-4-sonnet", false},
		{"claude", "claude-opus-4-5", "claude-opus-4-6", false}, {"codex", "opus", "claude-opus-4-6", false},
		{"codex", "gpt-6", "gpt-6", true}, {"codex", "gpt-6", "gpt-6.1", false},
	} {
		t.Run(tc.provider+"/"+tc.requested+"/"+tc.actual, func(t *testing.T) {
			if got := deliveryModelMatches(tc.provider, tc.requested, tc.actual); got != tc.match {
				t.Fatalf("model match=%v want %v", got, tc.match)
			}
		})
	}
}

func deliveryTestReview(t *testing.T, f deliveryFixture, key string) string {
	t.Helper()
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	return writeAny(t, f.R.WorkspaceRoot, "review-"+key+".json", map[string]any{"kind": "DeliveryCandidateReview@1", "schema_version": 1, "candidate_oid": gitOutput(t, target, "rev-parse", "HEAD"), "candidate_tree": gitOutput(t, target, "rev-parse", "HEAD^{tree}"), "reviewer_claim": "synthetic fixture reviewer", "reviewer_session_id": "fixture-review-session", "decision": "passed", "findings": []any{}, "fixes": []any{}, "open_actionable_findings": []any{}})
}

func deliveryTestCandidate(t *testing.T, f deliveryFixture, text string) string {
	t.Helper()
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if e := os.WriteFile(filepath.Join(target, "README.md"), []byte(text+"\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, target, "add", "README.md")
	runGit(t, target, "commit", "-m", "synthetic candidate "+text)
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = "+ShellQuote(text)+"\nprintf 'observed fixture behavior\\n'\n"), 0600); e != nil {
		t.Fatal(e)
	}
	return gitOutput(t, target, "rev-parse", "HEAD")
}

func deliveryTestHuman(t *testing.T, f deliveryFixture, o WorkflowRun, outcome string) string {
	t.Helper()
	c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	return writeAny(t, f.R.WorkspaceRoot, fmt.Sprintf("human-%d-%s.json", c.Number, outcome), workflowhandoff.DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: "task", TaskResultID: string(c.TaskResult.ID), ResultOID: c.OID, ResultTree: c.Tree, Outcome: outcome, ActorClaim: "Synthetic fixture human, not real product approval", StartSurface: "isolated automated test", StartedAtUTC: "2026-10-05T00:00:00Z", CompletedAtUTC: "2026-10-05T00:00:01Z", Answer: outcome, Observation: "Synthetic protocol fixture only; no actual human QA."})
}

func deliveryTestSingleReceipt(t *testing.T, o WorkflowRun) {
	t.Helper()
	for i, event := range o.Delivery.Events {
		if event.Kind != "verification" {
			continue
		}
		path := filepath.Join(o.Paths.RunRoot, "delivery", "attempts", event.ID, "verification.json")
		if event.Binding.Locator != path {
			t.Errorf("verification does not bind its one receipt: %s", event.Binding.Locator)
		}
		if _, err := workflowBound(event.Binding, 256<<10); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(o.Paths.RunRoot, "delivery", "events", fmt.Sprintf("%08d.json", i+1))); !os.IsNotExist(err) {
			t.Errorf("verification retained a duplicate numbered event: %v", err)
		}
		for _, candidate := range o.Delivery.Candidates {
			if candidate.Key == event.ID && candidate.Verification != event.Binding {
				t.Error("candidate and event bind different receipts")
			}
		}
	}
}

func TestDeliveryOwnerActualAcceptanceAndTruthfulProgress(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	if o.Delivery == nil || o.Delivery.PermissionState != "pending_runtime_acceptance" {
		t.Fatalf("launch granted observed permission: %+v", o.Delivery)
	}
	template := workflowProviderDocument(t, filepath.Join(o.Paths.RunRoot, "delivery", "templates", "acceptance.json"))
	if template["runtime_claim"].(map[string]any)["effective_policy_sha256"] != nil || template["delivery_permission"].(map[string]any)["permission_confirmed"] != false {
		t.Fatalf("template fabricated acceptance: %#v", template)
	}
	guide, e := os.ReadFile(filepath.Join(o.Paths.RunRoot, "delivery", "callback-guide.md"))
	if e != nil || !strings.Contains(string(guide), "Choose reviews, testing, integration, installation and cleanup explicitly") {
		t.Fatalf("callback guide missing: %s %v", guide, e)
	}
	a := deliveryTestAcceptance(t, f, o)
	a.DeliveryPermission.PermissionConfirmed = false
	if _, e = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "unconfirmed.json", a)); e == nil {
		t.Fatal("requested policy was enough to grant acceptance")
	}
	a.DeliveryPermission.PermissionConfirmed = true
	o, e = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "confirmed.json", a))
	if e != nil {
		t.Fatal(e)
	}
	if o.Delivery.PermissionState != "recipient_confirmed_contract" || *o.Delivery.ActualPolicySHA256 == f.R.Runtime.PermissionBinding.EffectivePolicySHA256 {
		t.Fatalf("actual policy confused with requested: %+v", o.Delivery)
	}
	question := "Which of these two product behaviors do you want?"
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "needs-choice", Phase: "needs_input", Summary: "A concrete product choice needs the human.", Meaning: "No completed candidate is claimed.", Question: &question, Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{{ID: "future-human-qa", Outcome: "not_run", Reason: "The human journey happens after a candidate is verified.", Argv: []string{}, Evidence: []FileBinding{}}}}
	file := writeAny(t, f.R.WorkspaceRoot, "incomplete-report.json", r)
	o, e = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
	if e != nil {
		t.Fatal(e)
	}
	if o.Delivery.Phase != "needs_input" || len(o.Delivery.Events) != 1 || len(o.Delivery.Candidates) != 0 || o.NextAction.Actor != "user" {
		t.Fatalf("incomplete report lost or qualified: %+v", o)
	}
	if _, e = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); e != nil {
		t.Fatalf("identical report retry failed: %v", e)
	}
	r.VerifierResults[0].Exit = ptr(0)
	if _, e = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "fabricated-exit.json", r)); e == nil {
		t.Fatal("unrun verifier fabricated a zero exit")
	}
	if _, e = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, writeAny(t, f.R.WorkspaceRoot, "legacy-review.json", WorkflowReview{})); e == nil {
		t.Fatal("legacy coordinator review acquired delivery ownership")
	}
}

func TestDeliveryOwnerUnknownVerificationNeverReplays(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "verify once")

	counter := filepath.Join(f.R.WorkspaceRoot, "verifier-started")
	if e := os.WriteFile(f.AcceptancePath, []byte("printf 'run\\n' >> "+ShellQuote(counter)+"\n"), 0600); e != nil {
		t.Fatal(e)
	}
	d := f.D
	d.Fault = func(point string) error {
		if point == "delivery_after_verification_reservation" {
			return errors.New("synthetic interruption")
		}
		return nil
	}
	if _, e := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
		t.Fatal("fault did not stop")
	}
	if _, e := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil || !strings.Contains(e.Error(), "unknown") {
		t.Fatalf("unknown attempt replay allowed: %v", e)
	}
	if _, e := os.Stat(counter); !os.IsNotExist(e) {
		t.Fatal("verifier started after unknown reservation")
	}
	before := workflowTestCalls(t, f.workflowFixture, "agent start")
	if _, e := WorkflowStart(f.D, f.File, "unused"); e == nil { /* existing request is legitimate readback, never another start */
	}
	if workflowTestCalls(t, f.workflowFixture, "agent start") != before {
		t.Fatal("same request replay started another agent")
	}
	_, _ = WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 10*time.Millisecond)
}

type deliveryInterruptVerificationState struct {
	systemFileSystem
	after bool
}

func (f deliveryInterruptVerificationState) Replace(path string, raw []byte) error {
	var s workflowState
	if filepath.Base(path) == "state.json" && decode(raw, 8<<20, &s) == nil && s.Result.Delivery != nil && s.Result.Delivery.Attempt != nil && s.Result.Delivery.Attempt.State == "recorded" {
		if f.after {
			if err := f.systemFileSystem.Replace(path, raw); err != nil {
				return err
			}
		}
		return errors.New("interrupted verification state publication")
	}
	return f.systemFileSystem.Replace(path, raw)
}

func TestDeliveryOwnerCompletedVerificationRecoversWithoutRerun(t *testing.T) {
	for _, point := range []string{"delivery_after_verification_receipt", "state_before_replace", "state_after_replace"} {
		t.Run(point, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			deliveryTestCandidate(t, f, "one execution")

			counter := filepath.Join(f.R.WorkspaceRoot, "verifier-started")
			if e := os.WriteFile(f.AcceptancePath, []byte("printf 'run\\n' >> "+ShellQuote(counter)+"\ntest \"$(cat README.md)\" = 'one execution'\n"), 0600); e != nil {
				t.Fatal(e)
			}
			d := f.D
			d.Fault = func(stage string) error {
				if point == stage {
					return errors.New("synthetic lost response")
				}
				return nil
			}
			if strings.HasPrefix(point, "state_") {
				d.Files = deliveryInterruptVerificationState{after: point == "state_after_replace"}
			}
			if _, e := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
				t.Fatal("fault did not stop")
			}
			path := filepath.Join(o.Paths.RunRoot, "delivery", "attempts", "verify-00000001", "verification.json")
			before, e := os.ReadFile(path)
			if e != nil {
				t.Fatal(e)
			}
			o, e = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "verify-00000001")
			if e != nil || o.Delivery.Phase != "verified" {
				t.Fatalf("receipt recovery failed: %+v %v", o, e)
			}
			deliveryTestSingleReceipt(t, o)
			after, err := os.ReadFile(path)
			registry, registryErr := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if err != nil || registryErr != nil || string(before) != string(after) || len(o.Delivery.Events) != 1 || len(o.Delivery.Candidates) != 1 || len(registry.TaskResults) != 1 {
				t.Fatalf("recovery changed receipt or duplicated publication: %v %v", err, registryErr)
			}
			b, e := os.ReadFile(counter)
			if e != nil || string(b) != "run\n" {
				t.Fatalf("verifier was replayed: %q %v", b, e)
			}
		})
	}
}
