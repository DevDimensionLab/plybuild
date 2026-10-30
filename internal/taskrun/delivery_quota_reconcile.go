package taskrun

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const quotaSourceLimit = 32 << 20
const capacityRetirementScope = "native_completion_of_commands_active_at_capacity_and_full_tail;precapacity_terminal_results_preserved_as_history;retirement_only;no_effect_replay"
const quotaRetirementScope = "only_the_unstarted_continuation_after_last_tool_result;historical_effects_not_assessed"

// This is a retirement-only annotation of one mistaken pre-execution quota
// observation. It never resolves the refusal for an owner callback or resume.
type DeliveryQuotaReconcileInput struct {
	EventID        string `json:"event_id"`
	StopSHA256     string `json:"stop_sha256"`
	Provider       string `json:"provider"`
	NativeSession  string `json:"native_session"`
	Action         string `json:"action"`
	QuotaEvent     string `json:"quota_event"`
	LastToolResult string `json:"last_tool_result"`
	ProcessID      string `json:"process_id,omitempty"`
	CommandSHA256  string `json:"command_sha256,omitempty"`
	SourceSHA256   string `json:"source_sha256"`
	Observer       string `json:"observer"`
	Reason         string `json:"reason"`
}

type deliveryQuotaReconciliation struct {
	Kind          string                      `json:"kind"`
	SchemaVersion int                         `json:"schema_version"`
	RunID         string                      `json:"run_id"`
	RequestSHA256 string                      `json:"request_sha256"`
	Scope         string                      `json:"scope"`
	Input         DeliveryQuotaReconcileInput `json:"input"`
	Stop          FileBinding                 `json:"stop"`
	Source        FileBinding                 `json:"source"`
	Evidence      FileBinding                 `json:"evidence"`
	History       []DeliveryEventRecord       `json:"history,omitempty"`
	RecordedAtUTC string                      `json:"recorded_at_utc"`
}

func WorkflowReconcileQuotaStop(d Dependencies, root, id string, in DeliveryQuotaReconcileInput) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if !plain(in.EventID, 1, 128) || !plain(in.Observer, 1, 256) || !plain(in.Reason, 1, 4096) ||
		!plain(in.Action, 1, 2000) || !plain(in.QuotaEvent, 1, 128) || !plain(in.LastToolResult, 1, 128) ||
		!digestPattern.MatchString(in.StopSHA256) || !digestPattern.MatchString(in.SourceSHA256) {
		return WorkflowRun{}, workflowError(2, "exact stop, action, quota event, last tool result, source digest, observer and reason are required")
	}
	err := workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "quota reconciliation requires an interactive delivery")
		}
		for _, event := range s.Result.Delivery.Events {
			if event.ID != in.EventID {
				continue
			}
			if event.Kind != "quota_stop_reconciliation" {
				return workflowError(4, "event ID is already used")
			}
			old, err := deliveryReadQuotaReconciliation(*s, event.Binding)
			if err != nil {
				return err
			}
			if !equal(old.Input, in) {
				return workflowError(4, "quota correction already has different immutable input")
			}
			return nil // Lost replies reuse the preserved observation, not a new transcript snapshot.
		}
		dl := s.Result.Delivery
		if dl.Phase == "completed" || dl.OwnershipRelease != nil || dl.HumanIntegration != nil || ownerRetired(*s) {
			return workflowError(4, "owner is already closed or released")
		}
		if err := deliveryQuotaAttemptGuard(*s); err != nil {
			return err
		}
		if dl.ActiveStop == nil || dl.ActiveStop.SHA256 != in.StopSHA256 {
			return workflowError(4, "correction must bind the current unresolved stop")
		}
		stop := *dl.ActiveStop
		if err := deliveryQuotaStopBinding(*s, stop, in); err != nil {
			return err
		}
		for _, e := range dl.Events {
			if e.Kind != "quota_stop_reconciliation" {
				continue
			}
			r, err := deliveryReadQuotaReconciliation(*s, e.Binding)
			if err != nil {
				return err
			}
			if r.Stop == stop && r.Source.SHA256 == in.SourceSHA256 {
				return workflowError(4, "this exact stop and source already have an immutable quota correction; reuse its event ID")
			}
		}
		path, err := deliveryReconciliationSession(*s)
		if err != nil {
			return err
		}
		source, raw, err := deliveryReadBinding(path, quotaSourceLimit)
		if err != nil {
			return err
		}
		if source.SHA256 != in.SourceSHA256 {
			return workflowError(4, "native transcript differs from the supplied source digest")
		}
		if err := deliveryValidateQuotaTranscript(*s, in, raw); err != nil {
			return err
		}
		copyPath := filepath.Join(s.Result.Paths.RunRoot, "delivery", "quota-evidence", strings.TrimPrefix(source.SHA256, "sha256:"))
		if err := d.writeOnce(copyPath, raw); err != nil {
			return err
		}
		r := deliveryQuotaReconciliation{Kind: "ply.workflow.quota-stop-reconciliation", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, Scope: deliveryReconciliationScope(in.Provider), Input: in, Stop: stop, Source: source, Evidence: FileBinding{copyPath, source.SHA256}, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		if in.Provider == "codex" {
			r.History = append([]DeliveryEventRecord(nil), dl.Events...)
		}
		_, err = deliveryAppendEvent(d, s, in.EventID, "quota_stop_reconciliation", r)
		return err
	})
	return deliveryReadback(d, root, id, err)
}

func deliveryQuotaAttemptGuard(s workflowState) error {
	if a := s.Result.Delivery.Attempt; a != nil && a.State != "recorded" {
		return workflowError(4, "unresolved native attempt blocks quota reconciliation and retirement")
	}
	if s.OwnerResume != nil && s.OwnerResumeStep != ownerResumePromptSent {
		return workflowError(4, "unresolved owner resume attempt blocks quota reconciliation and retirement")
	}
	return nil
}

func deliveryQuotaStopBinding(s workflowState, binding FileBinding, in DeliveryQuotaReconcileInput) error {
	raw, err := workflowBound(binding, 2<<20)
	if err != nil {
		return err
	}
	var stop DeliveryStopRecord
	if err := decode(raw, 2<<20, &stop); err != nil {
		return err
	}
	if stop.Kind != "ply.workflow.coordinator-stop" || stop.RunID != s.Result.RunID ||
		(in.Provider != "claude" && in.Provider != "codex") || in.Provider != s.Request.Runtime.Provider || in.Provider != stop.Input.Provider ||
		in.NativeSession != s.Result.Transport.AgentSessionID || in.NativeSession != stop.Input.NativeSession ||
		in.StopSHA256 != binding.SHA256 || in.Action != stop.Input.Action || stop.Input.Stage != "work" || stop.Input.Rule != "" ||
		(stop.Input.EffectStatus != DeliveryEffectPartial && stop.Input.EffectStatus != DeliveryEffectUnknown) {
		return workflowError(4, "quota correction must bind the exact provider work stop, session and action; rule refusals cannot be corrected")
	}
	return nil
}

func deliveryReadQuotaReconciliation(s workflowState, binding FileBinding) (deliveryQuotaReconciliation, error) {
	var r deliveryQuotaReconciliation
	raw, err := workflowBound(binding, 1<<20)
	if err != nil {
		return r, err
	}
	if err := decode(raw, 1<<20, &r); err != nil {
		return r, err
	}
	if r.Kind != "ply.workflow.quota-stop-reconciliation" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.Scope != deliveryReconciliationScope(r.Input.Provider) || r.Source.SHA256 != r.Input.SourceSHA256 || r.Evidence.SHA256 != r.Source.SHA256 {
		return r, workflowError(4, "quota correction Run, request or source binding changed")
	}
	if r.Input.Provider == "codex" {
		index := -1
		for i, event := range s.Result.Delivery.Events {
			if event.ID == r.Input.EventID && event.Kind == "quota_stop_reconciliation" {
				index = i
				break
			}
		}
		if index != len(r.History) {
			return r, workflowError(4, "capacity audit does not bind its exact preceding native history")
		}
		if len(r.History) == 0 || len(r.History) > len(s.Result.Delivery.Events) || !equal(r.History, s.Result.Delivery.Events[:len(r.History)]) {
			return r, workflowError(4, "capacity audit native event history changed")
		}
		for _, event := range r.History {
			if _, err := workflowBound(event.Binding, 2<<20); err != nil {
				return r, err
			}
		}
	}
	if err := deliveryQuotaStopBinding(s, r.Stop, r.Input); err != nil {
		return r, err
	}
	proof, err := workflowBound(r.Evidence, quotaSourceLimit)
	if err != nil {
		return r, err
	}
	return r, deliveryValidateQuotaTranscript(s, r.Input, proof)
}

// Only retirement consumes this annotation. Re-read the native source before
// any pane close; an appended tool or replaced transcript invalidates the proof.
func deliveryQuotaAllowsRetirement(s workflowState, stop DeliveryStopRecord) (bool, error) {
	// A later source-bound audit may cover benign final metadata. Older audit
	// snapshots stay immutable history; only the newest for this exact stop is
	// eligible, and it still must match the complete current native source.
	for i := len(s.Result.Delivery.Events) - 1; i >= 0; i-- {
		event := s.Result.Delivery.Events[i]
		if event.Kind != "quota_stop_reconciliation" {
			continue
		}
		r, err := deliveryReadQuotaReconciliation(s, event.Binding)
		if err != nil {
			return false, err
		}
		if r.Input.Action != stop.Input.Action {
			continue
		}
		raw, err := workflowBound(r.Stop, 2<<20)
		if err != nil {
			return false, err
		}
		var bound DeliveryStopRecord
		if err := decode(raw, 2<<20, &bound); err != nil {
			return false, err
		}
		if !equal(bound, stop) {
			continue
		}
		if err := deliveryQuotaAttemptGuard(s); err != nil {
			return false, err
		}
		path, err := deliveryReconciliationSession(s)
		if err != nil {
			return false, err
		}
		source, _, err := deliveryReadBinding(path, quotaSourceLimit)
		if err != nil {
			return false, err
		}
		if source != r.Source {
			return false, workflowError(4, "native quota source changed from "+r.Source.SHA256+" to "+source.SHA256+" at "+path+"; inspect actual activity before a new source-bound audit, never replay a close or clear new tool effects")
		}
		return true, nil
	}
	return false, nil
}

type quotaTranscriptRow struct {
	Type         string          `json:"type"`
	UUID         string          `json:"uuid"`
	Session      string          `json:"sessionId"`
	SessionAlias string          `json:"session_id"`
	CWD          string          `json:"cwd"`
	Sidechain    bool            `json:"isSidechain"`
	APIError     string          `json:"apiError"`
	IsAPIError   bool            `json:"isApiErrorMessage"`
	APIStatus    int             `json:"apiErrorStatus"`
	SourceTool   string          `json:"sourceToolAssistantUUID"`
	Message      json.RawMessage `json:"message"`
	ToolResult   struct {
		Interrupted      *bool  `json:"interrupted"`
		Async            bool   `json:"isAsync"`
		Status           string `json:"status"`
		AgentID          string `json:"agentId"`
		TaskID           string `json:"task_id"`
		BackgroundTaskID string `json:"backgroundTaskId"`
	} `json:"-"`
	Attachment struct {
		Type        string `json:"type"`
		CommandMode string `json:"commandMode"`
		Prompt      string `json:"prompt"`
		Origin      struct {
			Kind     string `json:"kind"`
			Producer string `json:"producer"`
		} `json:"origin"`
	} `json:"attachment"`
}

func deliveryValidateQuotaTranscript(s workflowState, in DeliveryQuotaReconcileInput, raw []byte) error {
	if in.Provider == "codex" {
		return deliveryValidateCapacityTranscript(s, in, raw)
	}
	fail := func(reason string) error { return workflowError(4, "native pre-execution quota proof: "+reason) }
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return fail("incomplete JSONL source")
	}
	type tool struct {
		row, name, cwd string
		background     bool
	}
	pending := map[string]tool{}
	async := map[string]string{} // tool use -> native task id, until native completion
	seenTools, seenRows := map[string]bool{}, map[string]bool{}
	quota, window := false, false
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var row quotaTranscriptRow
		if err := json.Unmarshal(line, &row); err != nil {
			return fail("malformed native row")
		}
		if row.Type == "assistant" || row.Type == "user" || row.Type == "system" || row.Type == "attachment" {
			if row.Session != in.NativeSession || row.SessionAlias != "" && row.SessionAlias != row.Session || row.Sidechain || row.UUID == "" || seenRows[row.UUID] {
				return fail("wrong session, sidechain or duplicate row")
			}
			if (window || row.UUID == in.LastToolResult || row.UUID == in.QuotaEvent) && row.CWD != s.Observed.Target.WorktreeLocator {
				return fail("wrong cwd in the bound stop window")
			}
			seenRows[row.UUID] = true
		}
		// Tool failures can be strings; inspect only structured result metadata.
		var object map[string]json.RawMessage
		if err := json.Unmarshal(line, &object); err != nil {
			return err
		}
		if tr := bytes.TrimSpace(object["toolUseResult"]); len(tr) > 0 && tr[0] == '{' {
			if err := json.Unmarshal(tr, &row.ToolResult); err != nil {
				return fail("malformed tool result metadata")
			}
		}
		var content []struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			Name   string `json:"name"`
			ToolID string `json:"tool_use_id"`
			Error  bool   `json:"is_error"`
			Input  struct {
				Background json.RawMessage `json:"run_in_background"`
			} `json:"input"`
		}
		var message struct {
			Content json.RawMessage `json:"content"`
		}
		if m := bytes.TrimSpace(row.Message); len(m) > 0 && m[0] == '{' {
			if err := json.Unmarshal(m, &message); err != nil {
				return fail("malformed native message")
			}
		}
		if c := bytes.TrimSpace(message.Content); len(c) > 0 && c[0] == '[' {
			if err := json.Unmarshal(c, &content); err != nil {
				return fail("malformed message content")
			}
		}
		for _, block := range content {
			switch block.Type {
			case "tool_use":
				if row.Type != "assistant" || window || quota || block.ID == "" || seenTools[block.ID] {
					return fail("tool use after the selected boundary or duplicate/unbound tool")
				}
				seenTools[block.ID] = true
				// Claude also records native wire arguments as strings. Recognizing
				// this form only marks async work pending; it never proves completion.
				background := false
				switch string(bytes.TrimSpace(block.Input.Background)) {
				case "", "false", `"false"`:
				case "true", `"true"`:
					background = true
				default:
					return fail("unsupported background argument in native tool call")
				}
				pending[block.ID] = tool{row.UUID, block.Name, row.CWD, background}
			case "tool_result":
				call, ok := pending[block.ToolID]
				if window || quota || row.Type != "user" || !ok || row.SourceTool != call.row {
					return fail("unbound or later tool result")
				}
				if row.ToolResult.Interrupted != nil && *row.ToolResult.Interrupted {
					return fail("interrupted tool effect remains unresolved")
				}
				if row.ToolResult.TaskID != "" || row.ToolResult.BackgroundTaskID != "" {
					return fail("background tool effect is unsupported and unresolved")
				}
				if call.background || row.ToolResult.Async {
					if call.name != "Agent" || !row.ToolResult.Async || row.ToolResult.Status != "async_launched" || row.ToolResult.AgentID == "" {
						return fail("unresolved asynchronous tool effect")
					}
					async[block.ToolID] = row.ToolResult.AgentID
				}
				delete(pending, block.ToolID)
				if row.UUID == in.LastToolResult {
					if call.name != "Bash" || call.cwd != s.Observed.Target.WorktreeLocator || call.background || row.ToolResult.Async || row.ToolResult.Interrupted == nil || *row.ToolResult.Interrupted || block.Error || len(pending) > 0 || len(async) > 0 {
						return fail("selected boundary requires a successful synchronous Bash result with native interrupted=false and no pending tools or async effects")
					}
					// Earlier matched returns remain history. This boundary proves only
					// that the next named continuation has no started tool activity.
					window = true
				}
			}
		}
		a := row.Attachment
		if a.Type == "queued_command" && a.CommandMode == "task-notification" && a.Origin.Kind == "task-notification" && a.Origin.Producer == "session-task" {
			var notice struct {
				XMLName xml.Name `xml:"task-notification"`
				Task    string   `xml:"task-id"`
				Tool    string   `xml:"tool-use-id"`
				Status  string   `xml:"status"`
			}
			if err := xml.Unmarshal([]byte(a.Prompt), &notice); err != nil {
				return fail("malformed native task notification")
			}
			if window || quota || notice.Task == "" || async[notice.Tool] != notice.Task || notice.Status != "completed" {
				return fail("unresolved or later asynchronous effect")
			}
			delete(async, notice.Tool)
		}
		if row.UUID == in.QuotaEvent {
			if quota || row.Type != "assistant" || row.APIError != "usage_limit_reached" || !row.IsAPIError || row.APIStatus != 429 || len(pending) > 0 || len(async) > 0 || !window {
				return fail("quota event lacks a completed tool boundary; async launch is not completion")
			}
			quota = true
		} else if row.Type == "assistant" && window {
			return fail("assistant activity after the selected boundary")
		}
	}
	if !quota {
		return fail(fmt.Sprintf("exact quota event %s is missing", in.QuotaEvent))
	}
	return nil
}
