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
	RunID            string           `json:"run_id"`
	TaskID           *string          `json:"task_id"`
	Provider         *string          `json:"provider"`
	Transport        *string          `json:"transport"`
	LogicalSessionID *string          `json:"logical_session_id"`
	Herdr            *InventoryHerdr  `json:"herdr"`
	State            string           `json:"state"`
	StateFreshness   string           `json:"state_freshness"`
	HistoricalState  *string          `json:"historical_state"`
	Unresolved       bool             `json:"unresolved"`
	StartedAtUTC     *string          `json:"started_at_utc"`
	LastSeenUTC      *string          `json:"last_seen_utc"`
	Return           *InventoryReturn `json:"return"`
	Freshness        string           `json:"freshness"`
	Sources          []FileBinding    `json:"sources"`
	Reasons          []Reason         `json:"reasons"`
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
				row = readHerdrInventoryRun(d, root, id)
			} else {
				row = readNativeInventoryRun(d, root, id)
			}
			out.Runs = append(out.Runs, row)
		}
	}
	sort.Slice(out.Runs, func(i, j int) bool { return out.Runs[i].RunID < out.Runs[j].RunID })
	return out, nil
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
	row.State, row.StateFreshness, row.Freshness, row.Unresolved = "unknown", "unknown", "unknown", true
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

func readHerdrInventoryRun(d Dependencies, root, id string) InventoryRun {
	row := inventoryEmpty(id, "herdr")
	reservation := workflowIndex(root, id)
	b, err := readFile(reservation, 4<<20, true)
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
	// Do not follow a stored run path until it matches its request-derived path.
	var reserved workflowState
	if err = decode(b, 4<<20, &reserved); err != nil {
		inventoryProblem(&row, err)
		return row
	}
	expected := filepath.Join(workflowRoot(root), "runs", id)
	if reserved.Result.Paths.RunRoot != expected {
		inventoryProblem(&row, integrity("run root differs from its registered store"))
		return row
	}
	statePath := filepath.Join(expected, "state.json")
	stateBytes, stateErr := readFile(statePath, 8<<20, true)
	s, err := workflowRead(root, id)
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
	row.TaskID, err = inventoryTask(s.Request.HandoffDraft)
	if err != nil {
		inventoryProblem(&row, err)
		return row
	}
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
	if s.Result.Delivery != nil {
		for _, event := range s.Result.Delivery.Events {
			if _, err := workflowBound(event.Binding, 8<<20); err != nil {
				inventoryProblem(&row, err)
			}
		}
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
	} else {
		inventorySource(&row, statePath, stateBytes, 8<<20)
	}
	inventorySource(&row, reservation, b, 4<<20)
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
