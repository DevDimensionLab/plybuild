package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func observedDeliveryStatus(t *testing.T, o WorkflowRun) map[string]any {
	t.Helper()
	raw, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err = json.Unmarshal(raw, &value); err != nil {
		t.Fatal(err)
	}
	status, ok := value["delivery_status"].(map[string]any)
	if !ok {
		t.Fatal("run readback omits the distinct source, verification, qualification, human judgment and delivery outcomes")
	}
	return status
}

func TestDeliveryStatusSeparatesCorrectionFromEarlierCandidate(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	first := deliveryTestCandidate(t, f, "first qualified candidate")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "first"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "fail", deliveryTestHuman(t, f, o, "fail"))
	if err != nil {
		t.Fatal(err)
	}
	failed := observedDeliveryStatus(t, o)
	if failed["human_judgment"].(map[string]any)["current"] != true {
		t.Fatal("the actual negative human answer for the unchanged source is not historical")
	}
	second := deliveryTestCandidate(t, f, "correction acceptance succeeds")
	review := workflowProviderDocument(t, deliveryTestReview(t, f, "second"))
	review["decision"] = "unobserved"
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "pending-review.json", review))
	if err == nil {
		t.Fatal("unobserved review unexpectedly qualified")
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	before, _ := os.ReadFile(statePath)
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	status := observedDeliveryStatus(t, o)
	source := status["source"].(map[string]any)
	verification := status["verification"].(map[string]any)
	qualified := status["qualified_candidate"].(map[string]any)
	judgment := status["human_judgment"].(map[string]any)
	delivery := status["final_delivery"].(map[string]any)
	if source["oid"] != second || verification["candidate_oid"] != second || verification["outcome"] != "passed" || verification["exit"] != float64(0) || verification["qualification"] != "not_qualified" {
		t.Fatalf("successful correction verification was lost or promoted: %#v", status)
	}
	if !strings.Contains(o.NextAction.Message, "--reuse "+o.Delivery.Attempt.ID) || !strings.Contains(o.NextAction.Message, "new verification") {
		t.Fatalf("readback still advises repeating a successful check without the reuse boundary: %s", o.NextAction.Message)
	}
	if qualified["oid"] != first || qualified["current"] != false || judgment["outcome"] != "fail" || judgment["current"] != false || delivery["state"] != "not_delivered" {
		t.Fatalf("old candidate/QA promoted to correction: %#v", status)
	}
	after, _ := os.ReadFile(statePath)
	if !bytes.Equal(before, after) {
		t.Fatal("status projection mutated preserved state")
	}
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "corrected"))
	if err != nil {
		t.Fatal(err)
	}
	status = observedDeliveryStatus(t, o)
	if status["qualified_candidate"].(map[string]any)["current"] != true || status["human_judgment"] != nil {
		t.Fatalf("new qualification inherited prior human answer: %#v", status)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	status = observedDeliveryStatus(t, o)
	if status["final_delivery"].(map[string]any)["state"] != "not_delivered" {
		t.Fatal("human pass claimed integration")
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	status = observedDeliveryStatus(t, o)
	if status["final_delivery"].(map[string]any)["state"] != "delivered" {
		t.Fatalf("observed integration absent: %#v", status)
	}
}

func TestDeliveryStatusPreservesCompletedReceiptBeforeQualification(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "receipt before bookkeeping")
	d := f.D
	d.Fault = func(stage string) error {
		if stage == "delivery_after_verification_receipt" {
			return errors.New("synthetic interruption after actual command completion")
		}
		return nil
	}
	o, err := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "interrupted"))
	if err == nil {
		t.Fatal("interruption did not occur")
	}
	status := observedDeliveryStatus(t, o)
	v := status["verification"].(map[string]any)
	if v["outcome"] != "passed" || v["qualification"] != "pending" || status["qualified_candidate"] != nil {
		t.Fatalf("receipt was hidden or treated as qualification: %#v", status)
	}
	path := filepath.Join(o.Delivery.Attempt.Path, "stdout.txt")
	if err = os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	status = observedDeliveryStatus(t, o)
	if status["verification"].(map[string]any)["outcome"] != "unknown" {
		t.Fatal("tampered output still represented as verified evidence")
	}
}

func TestDeliveryStatusDoesNotPresentTamperedQualificationAsCurrent(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "qualified candidate with preserved evidence")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "tamper"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(o.Delivery.Attempt.Path, "stdout.txt")
	if err = os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("changed after qualification"), 0600); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	status := observedDeliveryStatus(t, o)
	if status["qualified_candidate"].(map[string]any)["current"] != false {
		t.Fatal("tampered qualification evidence still labels the source as the current qualified candidate")
	}
	if status["verification"].(map[string]any)["outcome"] != "unknown" {
		t.Fatal("tampered command evidence retained a successful outcome")
	}
}

func TestDeliveryStatusKeepsCommandCompletionWhenReviewChanges(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "successful acceptance before failed qualification")
	review := deliveryTestReview(t, f, "original-review")
	broken := f.D
	broken.Workflow.TaskWorkspace = nil
	o, err := WorkflowDeliveryVerify(broken, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err == nil {
		t.Fatal("expected qualification failure")
	}
	if err = os.WriteFile(review, []byte("changed review input"), 0600); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	status := observedDeliveryStatus(t, o)
	v := status["verification"].(map[string]any)
	if v["outcome"] != "passed" || v["exit"] != float64(0) {
		t.Fatalf("changed live review erased the preserved successful command outcome: %#v", v)
	}
	if v["inputs_match"] != false || len(status["reasons"].([]any)) == 0 {
		t.Fatalf("changed review did not prevent reuse: %#v", status)
	}
}
