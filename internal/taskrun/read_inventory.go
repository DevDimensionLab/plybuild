package taskrun

// This adapter reads durable run evidence for workspace views. It deliberately
// does not use Show/WorkflowShow: those also qualify the present checkout and
// runtime, whereas a historical list must survive a moved checkout or upgrade.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type InventoryHerdr struct {
	WorkspaceID    *string `json:"workspace_id"`
	TabID          *string `json:"tab_id"`
	PaneID         *string `json:"pane_id"`
	TerminalID     *string `json:"terminal_id"`
	AgentSessionID *string `json:"agent_session_id"`
}

type InventoryReturn struct {
	State           string  `json:"state"`
	ReportSHA256    *string `json:"report_sha256"`
	TerminalSHA256  *string `json:"terminal_sha256"`
	ReportedOutcome *string `json:"reported_outcome"`
}

type InventoryRun struct {
	RunID            string                         `json:"run_id"`
	TaskID           *string                        `json:"task_id"`
	Provider         *string                        `json:"provider"`
	Transport        *string                        `json:"transport"`
	LogicalSessionID *string                        `json:"logical_session_id"`
	Herdr            *InventoryHerdr                `json:"herdr"`
	State            string                         `json:"state"`
	StateFreshness   string                         `json:"state_freshness"`
	HistoricalState  *string                        `json:"historical_state"`
	Unresolved       bool                           `json:"unresolved"`
	StartedAtUTC     *string                        `json:"started_at_utc"`
	LastSeenUTC      *string                        `json:"last_seen_utc"`
	Return           *InventoryReturn               `json:"return"`
	Freshness        string                         `json:"freshness"`
	Sources          []FileBinding                  `json:"sources"`
	Reasons          []Reason                       `json:"reasons"`
	Delivery         *InventoryDelivery             `json:"delivery,omitempty"`
	StartupRecovery  *InventoryStartupRecovery      `json:"startup_recovery,omitempty"`
	Closeout         *workspace.TaskCloseoutReceipt `json:"closeout,omitempty"`
}

type Inventory struct {
	Runs    []InventoryRun `json:"runs"`
	Reasons []Reason       `json:"reasons"`
}

var workflowRunIDPattern = regexp.MustCompile(`^wfr_[0-9a-f]{64}$`)

// ReadInventory returns all preserved attempts, including incomplete and
// damaged ones. An absent store is empty. A damaged store/run is explicit and
// cannot hide valid siblings. No lock, process probe, provider call or repair is
// performed. "Fresh" describes the files read, never provider liveness.
func ReadInventory(d Dependencies, root string) (Inventory, error) {
	out := Inventory{Runs: []InventoryRun{}, Reasons: []Reason{}}
	observed, err := workspace.ObserveRoot(d.Workspace, root)
	if err != nil {
		return out, err
	}
	if observed.Root != root {
		return out, conflict("inventory root is not the physical workspace root")
	}
	workflowIDs := map[string]bool{}
	for _, store := range []struct{ path, kind string }{
		{filepath.Join(storeRoot(root), "runs"), "native"},
		{filepath.Join(workflowRoot(root), "requests"), "herdr"},
	} {
		entries, err := inventoryEntries(store.path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			out.Reasons = append(out.Reasons, Reason{"run_store_unavailable", store.kind + " run store could not be read: " + err.Error()})
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".publish-") {
				continue
			}
			id := name
			if store.kind == "herdr" {
				id = strings.TrimSuffix(name, ".json")
			}
			valid := store.kind == "native" && runPattern.MatchString(id) || store.kind == "herdr" && strings.HasSuffix(name, ".json") && workflowRunIDPattern.MatchString(id)
			if !valid {
				out.Reasons = append(out.Reasons, Reason{"run_entry_unrecognized", fmt.Sprintf("Unrecognized entry in %s run store: %q", store.kind, name)})
				continue
			}
			var row InventoryRun
			if store.kind == "herdr" {
				workflowIDs[id] = true
				row = readHerdrInventoryRun(d, root, id)
			} else {
				row = readNativeInventoryRun(d, root, id)
			}
			out.Runs = append(out.Runs, row)
		}
		inventoryCheckEntries(&out, store.path, entries)
	}
	// A run directory without its request cannot safely acquire a Task identity.
	// Keep it as an identifiable orphan instead of silently dropping evidence.
	path := filepath.Join(workflowRoot(root), "runs")
	entries, err := inventoryEntries(path)
	if err != nil && !os.IsNotExist(err) {
		out.Reasons = append(out.Reasons, Reason{"run_store_unavailable", "herdr run directories could not be read: " + err.Error()})
	}
	for _, entry := range entries {
		id := entry.Name()
		if strings.HasPrefix(id, ".publish-") || workflowIDs[id] {
			continue
		}
		if !workflowRunIDPattern.MatchString(id) {
			out.Reasons = append(out.Reasons, Reason{"run_entry_unrecognized", fmt.Sprintf("Unrecognized entry in herdr run directories: %q", id)})
			continue
		}
		row := inventoryEmpty(id, "herdr")
		row.Reasons = append(row.Reasons, Reason{"run_source_orphaned", "Preserved run directory has no captured request: " + filepath.Join(path, id)})
		out.Runs = append(out.Runs, row)
	}
	if err == nil {
		inventoryCheckEntries(&out, path, entries)
	}
	sort.Slice(out.Runs, func(i, j int) bool { return out.Runs[i].RunID < out.Runs[j].RunID })
	return out, nil
}

func inventoryCheckEntries(out *Inventory, path string, before []os.DirEntry) {
	after, err := inventoryEntries(path)
	same := err == nil && len(before) == len(after)
	if same {
		for i := range before {
			if before[i].Name() != after[i].Name() || before[i].Type() != after[i].Type() {
				same = false
				break
			}
		}
	}
	if !same {
		out.Reasons = append(out.Reasons, Reason{"run_store_changed", "Run store entries changed during reading: " + path})
	}
}

func inventoryEntries(path string) ([]os.DirEntry, error) {
	r, err := physicalRoot(path, false)
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	f, err := r.Open(".")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	entries, err := f.ReadDir(-1)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, err
}

func inventoryEmpty(id, transport string) InventoryRun {
	return InventoryRun{RunID: id, Transport: nonempty(transport), State: "unknown", StateFreshness: "unknown", Unresolved: true, Freshness: "unknown", Sources: []FileBinding{}, Reasons: []Reason{}}
}
func nonempty(v string) *string {
	if v == "" {
		return nil
	}
	return ptr(v)
}
func parseInventoryTime(v string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, v)
	if err != nil || t.Location() != time.UTC {
		return time.Time{}, integrity("run observation time is not UTC")
	}
	return t, nil
}
func inventoryTask(draft []byte) (*string, error) {
	x, err := workflowhandoff.ValidateTaskRunDraft(draft)
	if err != nil {
		return nil, err
	}
	if x.Basis == nil {
		return nil, integrity("run has no registered Task binding")
	}
	return ptr(string(x.Basis.TaskID)), nil
}
func inventoryProblem(row *InventoryRun, err error) {
	row.State, row.StateFreshness, row.Unresolved = "unknown", "unknown", true
	if row.Freshness != "stale" {
		row.Freshness = "unknown"
	}
	row.Reasons = append(row.Reasons, Reason{"run_evidence_unavailable", err.Error()})
}
func inventorySource(row *InventoryRun, path string, before []byte, max int) {
	row.Sources = append(row.Sources, FileBinding{path, hash(before)})
	after, err := readFile(path, max, true)
	if err != nil || !bytes.Equal(before, after) {
		inventoryProblem(row, integrity("run source changed during reading"))
		row.Freshness = "stale"
	}
}

func readNativeInventoryRun(d Dependencies, root, id string) InventoryRun {
	row := inventoryEmpty(id, "")
	path := filepath.Join(storeRoot(root), "runs", id, "request.json")
	b, err := readFile(path, 1<<20, true)
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
	j, err := readJournal(root, id)
	if j.Request.RequestKey == "" {
		inventoryProblem(&row, err)
		return row
	}
	row.Freshness = "fresh"
	row.Provider = nonempty(j.Request.Runtime.Provider)
	row.Transport = ptr("terminal")
	if j.Request.Runtime.Mode == "exec" {
		row.Transport = ptr("headless")
	}
	row.TaskID, _ = inventoryTask(j.Request.HandoffDraft)
	if row.TaskID == nil {
		inventoryProblem(&row, integrity("run Task identity is unavailable"))
	}
	if j.Binding != nil {
		row.LogicalSessionID = nonempty(j.Binding.SessionID)
	}
	if len(j.Events) > 0 {
		row.LastSeenUTC = ptr(j.Events[len(j.Events)-1].RecordedAtUTC)
		for _, e := range j.Events {
			if e.Type == "process_started" && row.StartedAtUTC == nil {
				row.StartedAtUTC = ptr(e.RecordedAtUTC)
			}
		}
	}
	row.HistoricalState = nonempty(j.Result.Process.State)
	for _, event := range j.Events {
		if event.Type == "process_started" || event.Type == "process_exited" {
			var value struct {
				Process Process `json:"process"`
			}
			if json.Unmarshal(event.Payload, &value) == nil {
				row.HistoricalState = nonempty(value.Process.State)
			}
		}
	}
	if j.Result.Delivery.ReportSHA256 != nil || j.Result.Delivery.TerminalSHA256 != nil {
		row.Return = &InventoryReturn{j.Result.Delivery.State, j.Result.Delivery.ReportSHA256, j.Result.Delivery.TerminalSHA256, j.Result.Delivery.ReportedOutcome}
		row.State, row.StateFreshness = "returned", "fresh"
	}
	if j.Result.Collection.State == "qualified" {
		row.State, row.StateFreshness, row.Unresolved = "completed", "fresh", false
	} else if row.Return == nil && quiescent(j.Result.Process) && j.Result.Process.State != "not_started" {
		row.State, row.StateFreshness, row.Unresolved = "stopped", "fresh", false
	}
	if err != nil {
		inventoryProblem(&row, err)
	} else if err = validateRunEvidence(d, j); err != nil {
		inventoryProblem(&row, err)
	}
	row.Reasons = append(row.Reasons, j.Result.Reasons...)
	inventorySource(&row, path, b, 1<<20)
	return row
}

func readHerdrInventoryRun(d Dependencies, root, id string) (row InventoryRun) {
	row = inventoryEmpty(id, "herdr")
	defer func() {
		if row.Delivery != nil && row.Freshness != "fresh" {
			row.Delivery.Freshness = row.Freshness
			if row.Freshness == "stale" {
				row.Delivery.NextAction = InventoryDeliveryAction{Actor: "unknown", RecordedActor: row.Delivery.NextAction.RecordedActor, Kind: "inspect_delivery", Reason: "Run sources changed during reading; inspect a new read before acting.", EvidenceIDs: []string{row.RunID}}
				for i := range row.Delivery.Reports {
					row.Delivery.Reports[i].Current = false
				}
			}
		}
	}()
	reservation := workflowIndex(root, id)
	b, err := readFile(reservation, 4<<20, true)
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
	defer func() { inventorySource(&row, reservation, b, 4<<20) }()
	// Do not follow a stored run path until it matches its request-derived path.
	var reserved workflowState
	if err = decode(b, 4<<20, &reserved); err != nil {
		inventoryProblem(&row, err)
		return row
	}
	if workflowID(reserved.Request) != id || reserved.Request.WorkspaceRoot != root || reserved.Result.RunID != id || reserved.Result.RequestSHA256 != digest(reserved.Request) || reserved.Result.SessionID != "ply:"+id {
		inventoryProblem(&row, integrity("reservation binding changed"))
		return row
	}
	expected := filepath.Join(workflowRoot(root), "runs", id)
	if reserved.Result.Paths.RunRoot != expected {
		inventoryProblem(&row, integrity("run root differs from its registered store"))
		return row
	}
	statePath := filepath.Join(expected, "state.json")
	stateBytes, stateErr := readFile(statePath, 8<<20, true)
	// Decode the bytes captured here instead of rereading mutable state through
	// WorkflowShow/workflowRead. The immutable request supplies identity even when
	// its later state is unavailable; a bad state must not hide valid siblings.
	s := reserved
	if stateErr == nil {
		defer func() { inventorySource(&row, statePath, stateBytes, 8<<20) }()
		var captured inventoryWorkflowState
		if stateErr = decode(stateBytes, 8<<20, &captured); stateErr == nil {
			saved := captured.state()
			if !equal(reserved.Request, saved.Request) || !equal(reserved.CodexTrust, saved.CodexTrust) || !equal(reserved.ClaudeTrust, saved.ClaudeTrust) || saved.Result.RunID != id || !equal(reserved.Observed, saved.Observed) || saved.Result.Paths.RunRoot != expected || saved.Result.RequestSHA256 != digest(reserved.Request) || saved.Result.SessionID != "ply:"+id {
				stateErr = integrity("preserved state differs from reservation")
			} else {
				var sources []inventoryCapturedSource
				sources, stateErr = readInventoryRecovery(captured)
				defer func() {
					for _, source := range sources {
						inventorySource(&row, source.binding.Locator, source.bytes, source.max)
					}
				}()
				if stateErr == nil {
					s, row.StartupRecovery = saved, captured.Result.StartupRecovery
				}
			}
		}
	}
	var draft workflowhandoff.TaskRunDraft
	if deliveryRun(s.Request) {
		draft, err = workflowhandoff.ValidateDeliveryTaskRunDraft(s.Request.HandoffDraft)
	} else {
		draft, err = workflowhandoff.ValidateTaskRunDraft(s.Request.HandoffDraft)
	}
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
	if draft.Basis == nil {
		inventoryProblem(&row, integrity("run has no registered Task binding"))
		return row
	}
	row.TaskID = ptr(string(draft.Basis.TaskID))
	if s.Request.Runtime.Provider != "codex" && s.Request.Runtime.Provider != "claude" {
		inventoryProblem(&row, integrity("unknown bound provider"))
		return row
	}
	row.Freshness = "fresh"
	row.Provider = ptr(s.Request.Runtime.Provider)
	row.LogicalSessionID = nonempty(s.Result.SessionID)
	t := s.Result.Transport
	row.Herdr = &InventoryHerdr{nonempty(t.WorkspaceID), nonempty(t.TabID), nonempty(t.PaneID), nonempty(t.TerminalID), nonempty(t.AgentSessionID)}
	row.HistoricalState, row.LastSeenUTC = nonempty(t.State), nonempty(t.ObservedAt)
	if row.LastSeenUTC != nil {
		if _, err := parseInventoryTime(*row.LastSeenUTC); err != nil {
			row.LastSeenUTC = nil
			inventoryProblem(&row, err)
		}
	}
	// Herdr has no durable creation clock in older reservations. Never use the
	// latest observation or filesystem mtime as an invented start timestamp.
	if s.Result.Round.ReportSHA256 != nil || s.Result.FinalReturn.ReportSHA256 != nil || s.Result.FinalReturn.TerminalSHA256 != nil {
		report := s.Result.FinalReturn.ReportSHA256
		if report == nil {
			report = s.Result.Round.ReportSHA256
		}
		row.Return = &InventoryReturn{s.Result.FinalReturn.State, report, s.Result.FinalReturn.TerminalSHA256, nil}
		row.State, row.StateFreshness = "returned", "fresh"
	}
	if s.Result.FinalReturn.State == "accepted" || s.Result.Delivery != nil && s.Result.Delivery.Phase == "completed" {
		row.State, row.StateFreshness, row.Unresolved = "completed", "fresh", false
	}
	// Report hashes are verified against immutable slots. Current executable
	// identity, permission policy and checkout freshness are deliberately not
	// requalified by this historical view.
	if err := workflowInventoryReturn(d, s, &row); err != nil {
		inventoryProblem(&row, err)
	}
	if deliveryRun(s.Request) {
		readInventoryDelivery(s, draft.Basis, &row)
	}
	if s.Result.Delivery != nil {
		for _, candidate := range s.Result.Delivery.Candidates {
			if _, err := workflowBound(candidate.Verification, 8<<20); err != nil {
				inventoryProblem(&row, err)
			}
			if candidate.Integration != nil {
				if _, err := workflowBound(*candidate.Integration, 8<<20); err != nil {
					inventoryProblem(&row, err)
				}
			}
		}
	}
	row.Reasons = append(row.Reasons, s.Result.Reasons...)
	if stateErr != nil {
		inventoryProblem(&row, stateErr)
		if row.Delivery != nil {
			inventoryDeliveryProblem(&row, "delivery_state_unavailable", stateErr.Error())
			row.Delivery.NextAction = inventoryDeliveryUnknownAction(s, "Preserved delivery state could not be read or bound; inspect the recorded source.")
		}
	}
	if closeout, e := workflowTaskCloseout(root, s); e != nil {
		inventoryProblem(&row, e)
	} else if closeout != nil {
		row.Closeout = closeout
		if row.Freshness == "fresh" && closeout.State == "complete" {
			row.State, row.StateFreshness, row.Unresolved = "completed", "fresh", false
		}
		if row.Delivery != nil && row.Freshness == "fresh" {
			if closeout.State == "complete" {
				row.Delivery.NextAction = InventoryDeliveryAction{Actor: "human", RecordedActor: "user", Kind: "task_closed", Reason: "The Task is closed with its source resource " + closeout.ResourceState + "; inspect retained history.", EvidenceIDs: []string{closeout.OperationID}}
			} else {
				reason, kind, detail := workspace.TaskCloseoutPendingAction(*closeout)
				row.Reasons = append(row.Reasons, Reason{reason, detail})
				row.Delivery.NextAction = InventoryDeliveryAction{Actor: "human", RecordedActor: "user", Kind: kind, Reason: detail, EvidenceIDs: []string{closeout.OperationID}}
			}
		}
	}
	return row
}

func workflowInventoryReturn(d Dependencies, s workflowState, row *InventoryRun) error {
	if row.Return == nil {
		if row.State == "completed" {
			return integrity("completed run has no preserved return")
		}
		return nil
	}
	if s.Result.Delivery != nil {
		// Delivery completion is owned by its bound integration observation,
		// not a manually changed phase field or a provider's reported outcome.
		if row.State != "completed" {
			return nil
		}
		if row.Return.ReportSHA256 == nil || len(s.Result.Delivery.Candidates) == 0 {
			return integrity("delivery completion is unbound")
		}
		candidate := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if a := s.Request.Delivery.Agreement; a != nil && a.Mode == workspace.DeliveryPullRequest {
			if candidate.PullRequest == nil || candidate.PullRequest.SHA256 != *row.Return.ReportSHA256 || candidate.Integration != nil {
				return integrity("delivery PR completion differs from its candidate")
			}
			bound := false
			for _, event := range s.Result.Delivery.Events {
				if event.Kind == "pull_request" && event.Binding == *candidate.PullRequest {
					bound = true
				}
			}
			if !bound {
				return integrity("delivery PR completion has no native publication event")
			}
			raw, err := workflowBound(*candidate.PullRequest, 8<<20)
			if err != nil {
				return err
			}
			var publication struct {
				Kind          string                         `json:"kind"`
				RunID         string                         `json:"run_id"`
				RequestSHA256 string                         `json:"request_sha256"`
				Candidate     string                         `json:"candidate"`
				Receipt       FileBinding                    `json:"receipt"`
				Observed      PullRequestDeliveryObservation `json:"observed"`
			}
			if json.Unmarshal(raw, &publication) != nil || publication.Kind != "PlyDeliveryPullRequest@1" || publication.RunID != row.RunID || publication.RequestSHA256 != s.Result.RequestSHA256 || publication.Candidate != candidate.Key || !strings.EqualFold(publication.Observed.Repository, a.GitHubRepository) || publication.Observed.HeadRef != a.SourceRef || publication.Observed.HeadOID != candidate.OID || publication.Observed.BaseRef != a.TargetRef || publication.Observed.URL == "" || publication.Observed.Number < 1 || !publication.Observed.MetadataApplied {
				return integrity("delivery PR completion does not bind the exact published candidate")
			}
			_, err = workflowBound(publication.Receipt, 8<<20)
			return err
		}
		if candidate.Integration == nil || candidate.Integration.SHA256 != *row.Return.ReportSHA256 {
			return integrity("delivery completion differs from its candidate")
		}
		bound := false
		for _, event := range s.Result.Delivery.Events {
			if event.Kind == "integration" && event.Binding == *candidate.Integration {
				bound = true
			}
		}
		if !bound {
			return integrity("delivery completion has no integration event")
		}
		b, err := workflowBound(*candidate.Integration, 8<<20)
		if err != nil {
			return err
		}
		var observation struct {
			Kind          string `json:"kind"`
			RunID         string `json:"run_id"`
			RequestSHA256 string `json:"request_sha256"`
			Candidate     string `json:"candidate"`
			Result        struct {
				Completed bool `json:"completed"`
			} `json:"result"`
		}
		if json.Unmarshal(b, &observation) != nil || observation.Kind != "PlyDeliveryIntegration@1" || observation.RunID != row.RunID || observation.RequestSHA256 != s.Result.RequestSHA256 || observation.Candidate != candidate.Key || !observation.Result.Completed {
			return integrity("delivery integration does not record completion")
		}
		return nil
	}
	if row.Return.ReportSHA256 == nil {
		return integrity("workflow return report is missing")
	}
	var matching *workflowRecord
	for i := range s.Records {
		r := &s.Records[i]
		report, err := workflowReadReport(*r)
		if err != nil {
			return err
		}
		if report.RunID != s.Result.RunID || report.RequestSHA256 != s.Result.RequestSHA256 || report.SessionID != s.Result.SessionID || report.Round != r.Number {
			return integrity("workflow report identity differs")
		}
		terminal, err := workflowBound(r.Terminal, 2<<20)
		if err != nil {
			return err
		}
		slot := FileBinding{filepath.Join(filepath.Dir(r.Report.Locator), "report-slot.json"), digest(workflowReportSlot{report, terminal, r.TargetSHA256})}
		if _, err = workflowBound(slot, 4<<20); err != nil {
			return err
		}
		if r.Report.SHA256 == *row.Return.ReportSHA256 {
			matching = r
			row.Return.ReportedOutcome = ptr(report.Outcome)
		}
	}
	if matching == nil {
		return integrity("workflow return has no matching immutable report")
	}
	if row.State == "completed" {
		if matching.Review == nil {
			return integrity("accepted workflow return has no review")
		}
		b, err := workflowBound(*matching.Review, 256<<10)
		if err != nil {
			return err
		}
		review, err := workflowParseReview(b)
		if err != nil {
			return err
		}
		if review.RunID != row.RunID || review.RequestSHA256 != s.Result.RequestSHA256 || review.ReportSHA256 != matching.Report.SHA256 || review.Decision != "accepted" {
			return integrity("workflow completion review differs")
		}
		facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
		if err != nil {
			return err
		}
		if row.Return.TerminalSHA256 == nil || facts.Terminal == nil || facts.Terminal.SHA256 != *row.Return.TerminalSHA256 {
			return integrity("workflow completion terminal differs")
		}
	}
	return nil
}
