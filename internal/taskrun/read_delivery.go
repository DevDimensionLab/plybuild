package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// InventoryDelivery is a historical, read-only delivery-owner projection. Its
// reported progress never supplies native technical, human QA or integration
// outcomes; those remain owned by the registered Task records.
type InventoryDelivery struct {
	Phase           string                    `json:"phase"`
	Basis           *workspace.TaskSpecBasis  `json:"basis"`
	Goal            *workspace.TaskGoalRef    `json:"goal"`
	GoalSource      FileBinding               `json:"goal_source"`
	Report          *InventoryDeliveryReport  `json:"report"`
	Reports         []InventoryDeliveryReport `json:"reports"`
	NextAction      InventoryDeliveryAction   `json:"next_action"`
	LastActivityUTC *string                   `json:"last_activity_utc"`
	Freshness       string                    `json:"freshness"`
	Reasons         []Reason                  `json:"reasons"`
}

type InventoryDeliveryReport struct {
	EventID            string      `json:"event_id"`
	Binding            FileBinding `json:"binding"`
	Phase              string      `json:"phase"`
	Summary            string      `json:"summary"`
	Meaning            string      `json:"meaning"`
	Question           *string     `json:"question"`
	ProcessObservation *string     `json:"process_observation"`
	PreventiveFollowup *string     `json:"preventive_followup"`
	Current            bool        `json:"current"`
	RecordedAtUTC      *string     `json:"recorded_at_utc"`
}

// RecordedActor preserves the run owner's vocabulary while Actor gives the
// status projection an explicit human/agent/unknown classification.
type InventoryDeliveryAction struct {
	Actor         string   `json:"actor"`
	RecordedActor string   `json:"recorded_actor"`
	Kind          string   `json:"kind"`
	Reason        string   `json:"reason"`
	EventID       *string  `json:"event_id"`
	SinceUTC      *string  `json:"since_utc"`
	EvidenceIDs   []string `json:"evidence_ids"`
}

func inventoryDeliveryProblem(row *InventoryRun, code, detail string) {
	reason := Reason{code, detail}
	row.Reasons = append(row.Reasons, reason)
	if row.Freshness != "stale" {
		row.Freshness = "unknown"
	}
	row.State, row.StateFreshness, row.Unresolved = "unknown", "unknown", true
	if row.Delivery != nil {
		row.Delivery.Reasons = append(row.Delivery.Reasons, reason)
		row.Delivery.Freshness = row.Freshness
	}
}

func inventoryDeliveryActor(actor string) string {
	switch actor {
	case "user":
		return "human"
	case "recipient", "coordinator":
		return "agent"
	default:
		return "unknown"
	}
}

func inventoryDeliveryUnknownAction(s workflowState, reason string) InventoryDeliveryAction {
	return InventoryDeliveryAction{Actor: "unknown", RecordedActor: s.Result.NextAction.Actor, Kind: "inspect_delivery", Reason: reason, EvidenceIDs: []string{s.Result.RunID}}
}

func inventoryDeliveryAction(s workflowState) InventoryDeliveryAction {
	a := s.Result.NextAction
	out := InventoryDeliveryAction{Actor: inventoryDeliveryActor(a.Actor), RecordedActor: a.Actor, Kind: "continue_delivery", Reason: a.Message, EvidenceIDs: []string{s.Result.RunID}}
	switch s.Result.Delivery.Phase {
	case "awaiting_acceptance":
		out.Kind = "delivery_startup"
	case "working", "verifying":
	case "needs_input":
		out.Kind = "answer_question"
	case "stopped":
		out.Kind = "delivery_stopped"
	case "awaiting_human_qa":
		out.Kind = "prepare_human_qa"
	case "human_qa_passed", "integrating":
		out.Kind = "continue_integration"
		if s.Request.Delivery.Agreement != nil && s.Request.Delivery.Agreement.Mode == workspace.DeliveryPullRequest {
			out.Kind = "continue_pull_request"
		}
	case "completed":
		out.Kind = "delivery_completed"
	default:
		return inventoryDeliveryUnknownAction(s, "The preserved delivery phase is unknown; inspect its source.")
	}
	return out
}

// readInventoryDelivery uses the same captured state as the transport inventory.
// All effects and all native outcome qualification belong to their existing
// owner APIs. This adapter does not read runtime binaries, policy, Git, Herdr or
// the referenced working artifacts of a report.
func readInventoryDelivery(s workflowState, basis *workspace.TaskSpecBasis, row *InventoryRun) {
	out := &InventoryDelivery{Basis: basis, GoalSource: s.Request.Delivery.Goal, Reports: []InventoryDeliveryReport{}, Freshness: "fresh", Reasons: []Reason{}}
	row.Delivery = out
	if s.Result.Delivery == nil {
		out.Phase = "unknown"
		out.NextAction = inventoryDeliveryUnknownAction(s, "The preserved run has no delivery state.")
		inventoryDeliveryProblem(row, "delivery_state_unavailable", "The delivery request has no preserved delivery state")
		return
	}
	out.Phase, out.NextAction = s.Result.Delivery.Phase, inventoryDeliveryAction(s)
	if out.NextAction.Kind == "inspect_delivery" {
		inventoryDeliveryProblem(row, "delivery_phase_unknown", "Unrecognized preserved delivery phase: "+out.Phase)
	}
	bindingValid := s.Request.Envelope == deliveryEnv("herdr-run-request") && s.Result.Envelope == deliveryEnv("run") && s.Request.ReturnMode == "delivery_owner"
	if !bindingValid {
		inventoryDeliveryProblem(row, "delivery_binding_invalid", "Delivery request and state envelopes differ from the delivery-owner contract")
	}

	type capturedSource struct {
		binding FileBinding
		bytes   []byte
	}
	sources := []capturedSource{}
	readSource := func(binding FileBinding) ([]byte, error) {
		raw, err := workflowBound(binding, 8<<20)
		if err == nil {
			sources = append(sources, capturedSource{binding, raw})
		}
		return raw, err
	}
	goalBytes, err := readSource(out.GoalSource)
	if err == nil {
		out.Goal, err = inventoryDeliveryGoal(goalBytes, *basis, out.GoalSource.SHA256)
	}
	if err == nil {
		specBinding := FileBinding{filepath.Join(s.Request.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(basis.Spec.ManifestSHA256, "sha256:")+".json"), basis.Spec.ManifestSHA256}
		var specBytes []byte
		if specBytes, err = readSource(specBinding); err == nil {
			err = inventoryDeliveryGoalOrigin(specBytes, *basis, *out.Goal)
		}
	}
	if err != nil {
		out.Goal = nil
		inventoryDeliveryProblem(row, "delivery_goal_unavailable", "Preserved goal binding could not be read: "+err.Error())
	}

	events := s.Result.Delivery.Events
	eventRoot := filepath.Join(s.Result.Paths.RunRoot, "delivery", "events")
	eventEntries, eventErr := inventoryEntries(eventRoot)
	if eventErr != nil && !os.IsNotExist(eventErr) {
		inventoryDeliveryProblem(row, "delivery_event_store_unavailable", "Preserved delivery events could not be listed: "+eventErr.Error())
	}
	boundPaths := map[string]bool{}
	for _, event := range events {
		boundPaths[event.Binding.Locator] = true
	}
	for _, entry := range eventEntries {
		path := filepath.Join(eventRoot, entry.Name())
		if !strings.HasPrefix(entry.Name(), ".publish-") && !boundPaths[path] {
			inventoryDeliveryProblem(row, "delivery_event_orphaned", "Delivery event source has no binding in captured state: "+path)
		}
	}
	seen := map[string]bool{}
	var previous *string
	currentReport := false
	for i, event := range events {
		predecessor := previous
		previous = ptr(event.Binding.SHA256)
		if !key(event.ID) || seen[event.ID] || event.Binding.Locator != filepath.Join(s.Result.Paths.RunRoot, "delivery", "events", fmt.Sprintf("%08d.json", i+1)) {
			inventoryDeliveryProblem(row, "delivery_event_invalid", fmt.Sprintf("Delivery event %q has a duplicate or invalid slot binding", event.ID))
			continue
		}
		seen[event.ID] = true
		raw, err := readSource(event.Binding)
		if err != nil {
			inventoryDeliveryProblem(row, "delivery_event_unavailable", fmt.Sprintf("Delivery event %q: %s", event.ID, err))
			continue
		}
		if event.Kind != "report" {
			// Keep non-report events as immutable sources. They cannot supply
			// TaskResult, QA or integration facts through this adapter.
			if err = inventoryDeliveryEventIdentity(raw, s, event.Kind); err != nil {
				inventoryDeliveryProblem(row, "delivery_event_invalid", fmt.Sprintf("Delivery event %q: %s", event.ID, err))
			}
			continue
		}
		var report DeliveryReport
		if err = decode(raw, 8<<20, &report); err == nil {
			err = validateDeliveryReportContent(report)
		}
		if err == nil && (report.RunID != row.RunID || report.RequestSHA256 != s.Result.RequestSHA256 || report.SessionID != s.Result.SessionID || report.EventID != event.ID || !equal(report.PreviousEventSHA256, predecessor)) {
			err = integrity("report identity or predecessor differs from its preserved event")
		}
		if err != nil {
			inventoryDeliveryProblem(row, "delivery_report_invalid", fmt.Sprintf("Delivery report %q: %s", event.ID, err))
			continue
		}
		item := InventoryDeliveryReport{EventID: event.ID, Binding: event.Binding, Phase: report.Phase, Summary: report.Summary, Meaning: report.Meaning, Question: report.Question, ProcessObservation: report.ProcessObservation, PreventiveFollowup: report.PreventiveFollowup}
		item.Current = bindingValid && i == len(events)-1 && out.Phase == report.Phase && (s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.State == "recorded")
		if item.Current && report.Phase == "needs_input" && (s.Result.NextAction.Actor != "user" || s.Result.NextAction.Message != *report.Question) {
			item.Current = false
			inventoryDeliveryProblem(row, "delivery_question_binding_invalid", "The preserved needs_input action does not match the report's exact question and human actor")
		}
		out.Reports = append(out.Reports, item)
		if item.Current {
			currentReport = true
			out.NextAction.EventID = ptr(event.ID)
			out.NextAction.EvidenceIDs = append(out.NextAction.EvidenceIDs, event.ID, event.Binding.SHA256)
			if report.Phase == "needs_input" {
				out.NextAction.Reason = *report.Question
			}
		}
	}
	if !equal(previous, s.Result.Delivery.LastEventSHA256) {
		inventoryDeliveryProblem(row, "delivery_event_chain_invalid", "The last delivery event digest differs from the preserved event list")
		currentReport = false
		for i := range out.Reports {
			out.Reports[i].Current = false
		}
	}
	if len(out.Reports) > 0 {
		out.Report = &out.Reports[len(out.Reports)-1]
	}
	if (out.Phase == "needs_input" || out.Phase == "stopped") && !currentReport {
		inventoryDeliveryProblem(row, "delivery_current_report_unavailable", "The current delivery phase has no matching intact report")
		out.NextAction = inventoryDeliveryUnknownAction(s, "The preserved delivery need has no matching intact current report; inspect the recorded sources.")
	}
	if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.State == "attempted" {
		inventoryDeliveryProblem(row, "delivery_attempt_unresolved", "The preserved delivery effect has an unresolved attempt: "+attempt.ID)
		out.NextAction = inventoryDeliveryUnknownAction(s, "Inspect the unresolved delivery attempt "+attempt.ID+" before continuing; its effect is unknown.")
		out.NextAction.EvidenceIDs = append(out.NextAction.EvidenceIDs, attempt.ID)
	}
	for _, source := range sources {
		inventorySource(row, source.binding.Locator, source.bytes, 8<<20)
	}
	if eventErr == nil {
		captured := Inventory{}
		inventoryCheckEntries(&captured, eventRoot, eventEntries)
		if len(captured.Reasons) > 0 {
			row.Freshness = "stale"
		}
	}
	if row.Freshness == "stale" {
		inventoryDeliveryProblem(row, "delivery_source_changed", "A captured delivery source changed during reading")
		out.NextAction = inventoryDeliveryUnknownAction(s, "Delivery sources changed during reading; inspect a new read before acting.")
		for i := range out.Reports {
			out.Reports[i].Current = false
		}
	}
}

func inventoryDeliveryGoal(raw []byte, basis workspace.TaskSpecBasis, sha string) (*workspace.TaskGoalRef, error) {
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return nil, err
	}
	var manifest struct {
		Kind          string                    `json:"kind"`
		SchemaVersion int                       `json:"schema_version"`
		ContractKind  string                    `json:"contract_kind"`
		TaskID        workspace.TaskID          `json:"task_id"`
		SpecID        string                    `json:"spec_id"`
		Revision      int                       `json:"revision"`
		Problem       workspace.TaskRevisionRef `json:"problem"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, err
	}
	_, specErr := workspace.ParseTaskID(manifest.SpecID)
	if manifest.Kind != "WorkspaceTaskSpecRevision@2" || manifest.SchemaVersion != 2 || manifest.ContractKind != "goal" || manifest.TaskID != basis.TaskID || specErr != nil || manifest.Revision < 1 || manifest.Revision > 2147483647 || manifest.Problem != basis.Problem {
		return nil, integrity("goal manifest identity differs from the delivery Task and Problem binding")
	}
	return &workspace.TaskGoalRef{SpecID: manifest.SpecID, Spec: workspace.TaskRevisionRef{Revision: manifest.Revision, ManifestSHA256: sha}}, nil
}

func inventoryDeliveryGoalOrigin(raw []byte, basis workspace.TaskSpecBasis, goal workspace.TaskGoalRef) error {
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return err
	}
	var manifest struct {
		Kind          string                    `json:"kind"`
		SchemaVersion int                       `json:"schema_version"`
		ContractKind  string                    `json:"contract_kind"`
		TaskID        workspace.TaskID          `json:"task_id"`
		SpecID        string                    `json:"spec_id"`
		Revision      int                       `json:"revision"`
		Problem       workspace.TaskRevisionRef `json:"problem"`
		GoalOrigin    workspace.TaskGoalRef     `json:"goal_origin"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return err
	}
	if manifest.Kind != "WorkspaceTaskSpecRevision@2" || manifest.SchemaVersion != 2 || manifest.ContractKind != "execution" || manifest.TaskID != basis.TaskID || manifest.SpecID != basis.SpecID || manifest.Revision != basis.Spec.Revision || manifest.Problem != basis.Problem || manifest.GoalOrigin != goal {
		return integrity("the execution Spec's goal origin differs from the delivery goal")
	}
	return nil
}

func inventoryDeliveryEventIdentity(raw []byte, s workflowState, kind string) error {
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return err
	}
	var event struct {
		Kind          string `json:"kind"`
		RunID         string `json:"run_id"`
		RequestSHA256 string `json:"request_sha256"`
	}
	if err := json.Unmarshal(raw, &event); err != nil {
		return err
	}
	expected := map[string]string{"verification": "PlyDeliveryVerification@1", "human_qa": "PlyDeliveryHumanQA@1", "integration": "PlyDeliveryIntegration@1", "pull_request": "PlyDeliveryPullRequest@1"}[kind]
	if expected == "" || event.Kind != expected || event.RunID != s.Result.RunID || event.RequestSHA256 != s.Result.RequestSHA256 {
		return integrity("event kind, run or request binding differs")
	}
	return nil
}
