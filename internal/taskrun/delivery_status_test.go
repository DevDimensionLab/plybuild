package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

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
	o, err := WorkflowDeliveryVerify(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
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
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
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

func TestDeliveryStatusKeepsCommandCompletionWhenInputChanges(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "successful acceptance before failed qualification")
	review := f.AcceptancePath
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(review, []byte("changed command input"), 0600); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	status := observedDeliveryStatus(t, o)
	v := status["verification"].(map[string]any)
	if v["outcome"] != "passed" || v["exit"] != float64(0) {
		t.Fatalf("changed live input erased the preserved successful command outcome: %#v", v)
	}
	if v["inputs_match"] != false || len(status["reasons"].([]any)) == 0 {
		t.Fatalf("changed input did not prevent reuse: %#v", status)
	}
}
