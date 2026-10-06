package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These fixtures record historical files directly. No synthetic provider
// process or reported verifier result is treated as actual execution or QA.
func deliveryInventoryState(t *testing.T, f deliveryFixture, key string) workflowState {
	t.Helper()
	r := f.R
	r.RequestKey = "inventory/" + key
	contract := *r.Delivery
	r.Delivery = &contract
	r.Delivery.Goal = FileBinding{filepath.Join(r.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	x := workspace.PlanWorktreeObservation{StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
	s := workflowInitial(r, Observed{Target: x, Epic: x, RuntimeBindings: []FileBinding{}})
	s.Result.Transport = WorkflowTransport{WorkspaceID: "workspace", TabID: key, PaneID: key, TerminalID: key, AgentSessionID: "native-" + key, State: "working", ObservedAt: "2026-10-05T10:00:00Z"}
	for _, path := range []string{filepath.Dir(workflowIndex(r.WorkspaceRoot, s.Result.RunID)), s.Result.Paths.RunRoot} {
		if err := os.MkdirAll(path, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := writeValue(workflowIndex(r.WorkspaceRoot, s.Result.RunID), s); err != nil {
		t.Fatal(err)
	}
	return s
}

func deliveryInventoryReport(t *testing.T, f deliveryFixture, s *workflowState, id, phase string) DeliveryReport {
	t.Helper()
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, EventID: id, PreviousEventSHA256: s.Result.Delivery.LastEventSHA256, Phase: phase, Summary: "Implementing the selected behavior.", Meaning: "The owner preserves its actual reported state.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	action := WorkflowAction{"recipient", "Continue the selected delivery within the bound authority."}
	if phase == "needs_input" {
		report.Question = ptr("Which default should this command use?")
		action = WorkflowAction{"user", *report.Question}
	} else if phase == "stopped" {
		action = WorkflowAction{"user", report.Meaning}
	}
	if _, err := deliveryAppendEvent(f.D, s, id, "report", report); err != nil {
		t.Fatal(err)
	}
	s.Result.Delivery.Phase, s.Result.NextAction = phase, action
	if err := workflowSave(f.D, *s); err != nil {
		t.Fatal(err)
	}
	return report
}

func deliveryInventoryRows(t *testing.T, f deliveryFixture) map[string]InventoryRun {
	t.Helper()
	// These dependencies must never be reached by the historical reader.
	d := f.D
	d.Executable = func() (string, error) { t.Fatal("inventory observed the running executable"); return "", nil }
	d.Runner, d.ExecRunner = nil, nil
	d.Workspace.WorkGit, d.Workspace.IntegrationGit, d.Workflow.Git = nil, nil, nil
	before := treeState(t, f.R.WorkspaceRoot)
	result, err := ReadInventory(d, f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if treeState(t, f.R.WorkspaceRoot) != before {
		t.Fatal("inventory wrote to its captured workspace")
	}
	out := map[string]InventoryRun{}
	for _, row := range result.Runs {
		out[row.RunID] = row
	}
	return out
}

func hasInventoryReason(row InventoryRun, code string) bool {
	for _, reason := range row.Reasons {
		if reason.Code == code {
			return true
		}
	}
	return false
}

func TestReadInventoryDeliveryReportsKeepExactQuestionAndSeparateNativeOutcomes(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	runs := map[string]string{}
	for _, phase := range []string{"working", "needs_input", "stopped"} {
		s := deliveryInventoryState(t, f, phase)
		deliveryInventoryReport(t, f, &s, phase+"-event", phase)
		runs[phase] = s.Result.RunID
	}
	// Reading a preserved delivery must not requalify removed runtime binaries.
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	rows := deliveryInventoryRows(t, f)
	if len(rows) != 3 {
		t.Fatalf("lost sibling attempts: %+v", rows)
	}
	for phase, id := range runs {
		row := rows[id]
		if row.TaskID == nil || *row.TaskID != "task" || row.Provider == nil || *row.Provider != "codex" || row.Herdr == nil || row.Herdr.AgentSessionID == nil || row.LogicalSessionID == nil {
			t.Fatalf("lost run provenance: %+v", row)
		}
		if row.State != "unknown" || row.StateFreshness != "unknown" || row.Freshness != "fresh" || !row.Unresolved || row.Return != nil {
			t.Fatalf("reported phase became native result or live state: %+v", row)
		}
		x := row.Delivery
		if x == nil || x.Phase != phase || x.Basis == nil || !equal(x.Goal, &f.Prepared.Goal) || x.GoalSource.SHA256 != f.Prepared.Goal.Spec.ManifestSHA256 || x.Report == nil || !x.Report.Current || x.Report.EventID != phase+"-event" || len(x.Reports) != 1 {
			t.Fatalf("lost delivery/goal provenance: %+v", x)
		}
		if x.LastActivityUTC != nil || x.Report.RecordedAtUTC != nil || x.NextAction.SinceUTC != nil || row.StartedAtUTC != nil {
			t.Fatalf("invented report time from read time or cached transport: %+v", x)
		}
		switch phase {
		case "working":
			if x.NextAction.Actor != "agent" || x.NextAction.RecordedActor != "recipient" || x.NextAction.Kind != "continue_delivery" {
				t.Fatal(x.NextAction)
			}
		case "needs_input":
			if x.NextAction.Actor != "human" || x.NextAction.RecordedActor != "user" || x.NextAction.Kind != "answer_question" || x.Report.Question == nil || x.NextAction.Reason != *x.Report.Question {
				t.Fatalf("question or explicit human actor lost: %+v", x)
			}
		case "stopped":
			if x.NextAction.Actor != "human" || x.NextAction.RecordedActor != "user" || x.NextAction.Kind != "delivery_stopped" {
				t.Fatal(x.NextAction)
			}
		}
	}
}

func TestReadInventoryDeliveryCorruptionAndOrphansPreserveValidSiblings(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	valid := deliveryInventoryState(t, f, "valid")
	deliveryInventoryReport(t, f, &valid, "valid-question", "needs_input")
	brokenReport := deliveryInventoryState(t, f, "broken-report")
	deliveryInventoryReport(t, f, &brokenReport, "broken-question", "needs_input")
	if err := os.WriteFile(brokenReport.Result.Delivery.Events[0].Binding.Locator, []byte("corrupt report"), 0600); err != nil {
		t.Fatal(err)
	}
	brokenState := deliveryInventoryState(t, f, "broken-state")
	if err := os.WriteFile(filepath.Join(brokenState.Result.Paths.RunRoot, "state.json"), []byte("corrupt state"), 0600); err != nil {
		t.Fatal(err)
	}
	orphanID := "wfr_" + strings.Repeat("a", 64)
	if err := os.Mkdir(filepath.Join(workflowRoot(f.R.WorkspaceRoot), "runs", orphanID), 0700); err != nil {
		t.Fatal(err)
	}
	rows := deliveryInventoryRows(t, f)
	if len(rows) != 4 || rows[valid.Result.RunID].Delivery.Report == nil || !rows[valid.Result.RunID].Delivery.Report.Current || rows[valid.Result.RunID].Freshness != "fresh" {
		t.Fatalf("corruption hid a valid sibling: %+v", rows)
	}
	bad := rows[brokenReport.Result.RunID]
	if bad.TaskID == nil || bad.Delivery.Report != nil || bad.Delivery.NextAction.Kind != "inspect_delivery" || bad.Delivery.NextAction.Actor != "unknown" || bad.Freshness != "unknown" || !hasInventoryReason(bad, "delivery_event_unavailable") {
		t.Fatalf("corrupt report became a current question: %+v", bad)
	}
	bad = rows[brokenState.Result.RunID]
	if bad.TaskID == nil || *bad.TaskID != "task" || bad.Freshness != "unknown" || !hasInventoryReason(bad, "delivery_state_unavailable") {
		t.Fatalf("valid immutable Task binding lost with damaged mutable state: %+v", bad)
	}
	orphan := rows[orphanID]
	if orphan.TaskID != nil || orphan.Delivery != nil || orphan.Freshness != "unknown" || !hasInventoryReason(orphan, "run_source_orphaned") {
		t.Fatalf("orphan source acquired an invented Task: %+v", orphan)
	}
}

func TestReadInventoryDeliveryRejectsForgedReportBindingsAndOutcomes(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	for _, condition := range []string{"run", "request", "session", "event", "predecessor", "envelope", "native_outcome", "fabricated_exit", "actor", "question"} {
		t.Run(condition, func(t *testing.T) {
			s := deliveryInventoryState(t, f, condition)
			report := deliveryInventoryReport(t, f, &s, condition+"-event", "needs_input")
			switch condition {
			case "run":
				report.RunID = "wfr_" + strings.Repeat("0", 64)
			case "request":
				report.RequestSHA256 = hash([]byte("another request"))
			case "session":
				report.SessionID = "another session"
			case "event":
				report.EventID = "another-event"
			case "predecessor":
				report.PreviousEventSHA256 = ptr(hash([]byte("another predecessor")))
			case "envelope":
				report.SchemaVersion = 1
			case "native_outcome":
				report.Phase, report.Question = "completed", nil
			case "fabricated_exit":
				report.VerifierResults = []DeliveryVerifierResult{{ID: "reported", Outcome: "not_run", Exit: ptr(0), Reason: "Not executed.", Argv: []string{}, Evidence: []FileBinding{}}}
			case "actor":
				s.Result.NextAction.Actor = "recipient"
			case "question":
				s.Result.NextAction.Message = "A different question?"
			}
			raw, err := Canonical(report)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(s.Result.Delivery.Events[0].Binding.Locator, raw, 0600); err != nil {
				t.Fatal(err)
			}
			s.Result.Delivery.Events[0].Binding.SHA256 = hash(raw)
			s.Result.Delivery.LastEventSHA256 = ptr(hash(raw))
			if err = workflowSave(f.D, s); err != nil {
				t.Fatal(err)
			}
			row := readHerdrInventoryRun(f.D, f.R.WorkspaceRoot, s.Result.RunID)
			if row.TaskID == nil || row.Delivery == nil || row.Delivery.NextAction.Kind != "inspect_delivery" || row.Delivery.NextAction.Actor != "unknown" || row.Delivery.Report != nil && row.Delivery.Report.Current || row.State != "unknown" || row.Freshness != "unknown" {
				t.Fatalf("%s forged report was accepted: %+v delivery=%+v", condition, row, row.Delivery)
			}
		})
	}
}

func TestReadInventoryDeliveryOlderQuestionsAndReportedVerifiersStayHistorical(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s := deliveryInventoryState(t, f, "history")
	deliveryInventoryReport(t, f, &s, "old-question", "needs_input")
	report := deliveryInventoryReport(t, f, &s, "new-working", "working")
	report.VerifierResults = []DeliveryVerifierResult{{ID: "reported-test", Outcome: "passed", Exit: ptr(0), CWD: f.R.WorkspaceRoot, Argv: []string{"true"}, Evidence: []FileBinding{}, Reason: "A reported claim, never qualification."}}
	raw, err := Canonical(report)
	if err != nil {
		t.Fatal(err)
	}
	last := &s.Result.Delivery.Events[1]
	if err = os.WriteFile(last.Binding.Locator, raw, 0600); err != nil {
		t.Fatal(err)
	}
	last.Binding.SHA256, s.Result.Delivery.LastEventSHA256 = hash(raw), ptr(hash(raw))
	if err = workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	row := deliveryInventoryRows(t, f)[s.Result.RunID]
	if row.State != "unknown" || row.Return != nil || row.Delivery.NextAction.Actor != "agent" || len(row.Delivery.Reports) != 2 || row.Delivery.Reports[0].Current || !row.Delivery.Reports[1].Current {
		t.Fatalf("old question or reported verifier became native status: %+v", row)
	}
	encoded, err := json.Marshal(row.Delivery)
	if err != nil || strings.Contains(string(encoded), "verifier_results") || strings.Contains(string(encoded), "technical_assessment") || strings.Contains(string(encoded), "human_qa") {
		t.Fatalf("report projection exposes native qualification claims: %s %v", encoded, err)
	}
	// A later domain event supersedes the reported question without promoting
	// the cached delivery phase to a native technical/QA result.
	_, err = deliveryAppendEvent(f.D, &s, "preserved-verification", "verification", deliveryVerificationReceipt{Kind: "PlyDeliveryVerification@1", SchemaVersion: 1, RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256})
	if err != nil {
		t.Fatal(err)
	}
	s.Result.Delivery.Phase = "awaiting_human_qa"
	s.Result.NextAction = WorkflowAction{"recipient", "Prepare the exact installed candidate journey."}
	if err = workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	row = deliveryInventoryRows(t, f)[s.Result.RunID]
	if row.State != "unknown" || row.Return != nil || row.Delivery.Report.Current || row.Delivery.NextAction.Kind != "prepare_human_qa" || row.Delivery.NextAction.Actor != "agent" {
		t.Fatalf("cached delivery phase became native QA or revived a question: %+v", row)
	}
}

func TestReadInventoryDeliveryOrphanEventAndUnknownAttemptStayInspectable(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s := deliveryInventoryState(t, f, "partial")
	deliveryInventoryReport(t, f, &s, "working", "working")
	orphan := filepath.Join(s.Result.Paths.RunRoot, "delivery", "events", "00000002.json")
	if err := os.WriteFile(orphan, []byte("an unbound event written before a state-save failure"), 0600); err != nil {
		t.Fatal(err)
	}
	s.Result.Delivery.Attempt = &DeliveryAttempt{ID: "verify-unknown", State: "attempted", Kind: "verification"}
	if err := workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	row := deliveryInventoryRows(t, f)[s.Result.RunID]
	if !hasInventoryReason(row, "delivery_event_orphaned") || !hasInventoryReason(row, "delivery_attempt_unresolved") || row.TaskID == nil || row.Delivery.Report == nil || row.Delivery.Report.Current || row.Delivery.NextAction.Actor != "unknown" || row.Delivery.NextAction.Kind != "inspect_delivery" || row.Freshness != "unknown" {
		t.Fatalf("partial effects were hidden or became live progress: %+v delivery=%+v", row, row.Delivery)
	}
}

func TestInventoryCapturedSourceAndDirectoryChangesAreExplicit(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "state.json")
	before := []byte("captured source")
	if err := os.WriteFile(path, before, 0600); err != nil {
		t.Fatal(err)
	}
	entries, err := inventoryEntries(root)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("changed source"), 0600); err != nil {
		t.Fatal(err)
	}
	row := inventoryEmpty("run", "herdr")
	row.Freshness = "fresh"
	inventorySource(&row, path, before, 1<<20)
	if row.Freshness != "stale" || row.StateFreshness != "unknown" || !hasInventoryReason(row, "run_evidence_unavailable") || len(row.Sources) != 1 || row.Sources[0].SHA256 != hash(before) {
		t.Fatalf("source change lost captured digest or claimed freshness: %+v", row)
	}
	if err = os.WriteFile(filepath.Join(root, "another.json"), []byte("new source"), 0600); err != nil {
		t.Fatal(err)
	}
	out := Inventory{Runs: []InventoryRun{}, Reasons: []Reason{}}
	inventoryCheckEntries(&out, root, entries)
	if len(out.Reasons) != 1 || out.Reasons[0].Code != "run_store_changed" {
		t.Fatalf("directory change was hidden: %+v", out)
	}
}
