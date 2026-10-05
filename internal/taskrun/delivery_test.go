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
	if e != nil || !strings.Contains(string(guide), "not cryptographic human authentication") {
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

func TestDeliveryOwnerCandidateFailureCorrectionAndIntegration(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := newDeliveryFixture(t, provider)
			o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
			calls := workflowProviderCalls(t, f.workflowFixture, "agent start")
			if len(calls) != 1 {
				t.Fatalf("provider starts=%d", len(calls))
			}
			argv := strings.Join(calls[0], " ")
			if provider == "claude" && (!strings.Contains(argv, "--effort medium") || strings.Contains(argv, "default_permissions")) || provider == "codex" && !strings.Contains(argv, `model_reasoning_effort="medium"`) {
				t.Fatalf("wrong provider/effort argv: %s", argv)
			}
			first := deliveryTestCandidate(t, f, "first behavior")
			var e error
			o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "first"))
			if e != nil {
				t.Fatal(e)
			}
			if o.Delivery == nil {
				_, readErr := workflowRead(f.R.WorkspaceRoot, workflowID(f.R))
				t.Fatalf("verification lost readback: %+v; %v", o, readErr)
			}
			if o.Delivery.Phase != "awaiting_human_qa" || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != first {
				t.Fatalf("candidate not qualified: %+v", o.Delivery)
			}
			s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if e != nil {
				t.Fatal(e)
			}
			if e = workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); e == nil {
				t.Fatal("TaskResult released delivery owner before QA/integration")
			}
			if _, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
				t.Fatal("technical candidate counted as human pass")
			}
			unanswered := filepath.Join(o.Paths.RunRoot, "delivery", "templates", "human-qa.json")
			if _, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", unanswered); e == nil {
				t.Fatal("empty human template became a verdict")
			}
			o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "fail", deliveryTestHuman(t, f, o, "fail"))
			if e != nil {
				t.Fatal(e)
			}
			if o.Delivery.Phase != "working" || o.Delivery.Candidates[0].HumanQA.Outcome != "fail" {
				t.Fatal("fail did not retain owner and exact old verdict")
			}
			second := deliveryTestCandidate(t, f, "corrected behavior")
			o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "corrected"))
			if e != nil {
				t.Fatal(e)
			}
			if len(o.Delivery.Candidates) != 2 || o.Transport.AgentSessionID != "fixture-session" || o.Delivery.Candidates[0].OID != first {
				t.Fatal("correction lost same session or candidate history")
			}
			qa := deliveryTestHuman(t, f, o, "pass")
			o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", qa)
			if e != nil {
				t.Fatal(e)
			}
			o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", qa)
			if e != nil {
				t.Fatalf("same QA replay: %v", e)
			}
			o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
			if e != nil {
				t.Fatal(e)
			}
			if o.Delivery.Phase != "completed" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != second {
				t.Fatalf("local delivery incomplete: %+v", o)
			}
			registry, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
			if e != nil {
				t.Fatal(e)
			}
			if len(registry.TaskResults) != 2 || len(registry.HumanQARecords) != 2 {
				t.Fatalf("history incomplete: results=%d human=%d", len(registry.TaskResults), len(registry.HumanQARecords))
			}
			if e = workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); e != nil {
				t.Fatalf("completed delivery retained reservation: %v", e)
			}
			if n := workflowTestCalls(t, f.workflowFixture, "agent start"); n != 1 {
				t.Fatalf("delivery restarted %d times", n)
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
			if _, e = WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second); e != nil {
				t.Fatal(e)
			}
			if label := workflowProviderDocument(t, f.Model)["label"]; label != "finished Delivery fixture" {
				t.Fatalf("completed tab label=%v", label)
			}
		})
	}
}

func TestDeliveryOwnerUnknownVerificationNeverReplays(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "verify once")
	review := deliveryTestReview(t, f, "once")
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
	if _, e := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review); e == nil {
		t.Fatal("fault did not stop")
	}
	if _, e := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review); e == nil || !strings.Contains(e.Error(), "unknown") {
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

func TestDeliveryOwnerFailedVerifierReviewAndCandidateDrift(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "candidate")
	review := deliveryTestReview(t, f, "candidate")
	if e := os.WriteFile(f.AcceptancePath, []byte("printf 'actual failed check\\n'\nexit 7\n"), 0600); e != nil {
		t.Fatal(e)
	}
	var e error
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if e == nil || o.Delivery.Phase != "working" || len(o.Delivery.Candidates) != 0 {
		t.Fatalf("failed command qualified: %+v %v", o, e)
	}
	var failed deliveryVerificationReceipt
	if e = readValue(o.Delivery.Events[0].Binding.Locator, 1<<20, &failed); e != nil {
		t.Fatal(e)
	}
	if failed.Exit == nil || *failed.Exit != 7 || failed.CWD != f.Prepared.Preparation.Preparation.Plan.WorktreePath {
		t.Fatalf("actual failure evidence missing: %+v", failed)
	}
	if e = os.WriteFile(f.AcceptancePath, []byte("test \"$(cat README.md)\" = candidate\n"), 0600); e != nil {
		t.Fatal(e)
	}
	unobserved := filepath.Join(o.Paths.RunRoot, "delivery", "templates", "review.json")
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, unobserved)
	if e == nil || len(o.Delivery.Candidates) != 0 {
		t.Fatalf("unobserved review qualified: %+v %v", o, e)
	}
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if e != nil {
		t.Fatal(e)
	}
	human := deliveryTestHuman(t, f, o, "pass")
	deliveryTestCandidate(t, f, "changed after QA invitation")
	if _, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", human); e == nil {
		t.Fatal("human verdict transferred to changed candidate")
	}
}

func TestDeliveryOwnerCompletedVerificationRecoversWithoutRerun(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "one execution")
	review := deliveryTestReview(t, f, "once")
	counter := filepath.Join(f.R.WorkspaceRoot, "verifier-started")
	if e := os.WriteFile(f.AcceptancePath, []byte("printf 'run\\n' >> "+ShellQuote(counter)+"\ntest \"$(cat README.md)\" = 'one execution'\n"), 0600); e != nil {
		t.Fatal(e)
	}
	d := f.D
	d.Fault = func(point string) error {
		if point == "delivery_after_verification_receipt" {
			return errors.New("synthetic lost response")
		}
		return nil
	}
	if _, e := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review); e == nil {
		t.Fatal("fault did not stop")
	}
	var e error
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if e != nil || o.Delivery.Phase != "awaiting_human_qa" {
		t.Fatalf("receipt recovery failed: %+v %v", o, e)
	}
	b, e := os.ReadFile(counter)
	if e != nil || string(b) != "run\n" {
		t.Fatalf("verifier was replayed: %q %v", b, e)
	}
}
