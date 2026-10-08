package taskrun

// The trace adapter reads preserved contracts and events. It deliberately does
// not qualify the present checkout, ask a transport about liveness, or invoke
// any callback/recovery path. Its output is evidence, never execution authority.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type TraceHistory struct {
	Runs    []TraceRun `json:"runs"`
	Reasons []Reason   `json:"reasons"`
}

type TraceRun struct {
	RunID             string                   `json:"run_id"`
	Family            string                   `json:"family"`
	TaskID            string                   `json:"task_id"`
	Provider          string                   `json:"provider"`
	SessionID         string                   `json:"session_id"`
	RequestSHA256     string                   `json:"request_sha256"`
	NativeRunIDs      []string                 `json:"native_run_ids"`
	FrozenBasis       *workspace.TaskSpecBasis `json:"frozen_basis"`
	FrozenGoal        *workspace.TaskGoalRef   `json:"frozen_goal"`
	FrozenGoalSource  *FileBinding             `json:"frozen_goal_source"`
	Contract          json.RawMessage          `json:"contract"`
	DeclaredProcess   []string                 `json:"declared_process"`
	Sources           []FileBinding            `json:"sources"`
	Entries           []TraceEntry             `json:"entries"`
	Reasons           []Reason                 `json:"reasons"`
	HistoricalReasons []Reason                 `json:"historical_reasons"`
	Coverage          string                   `json:"coverage"`
}

type TraceCandidate struct {
	Key          string  `json:"key"`
	Number       int     `json:"number"`
	OID          string  `json:"oid"`
	Tree         string  `json:"tree"`
	TaskResultID *string `json:"task_result_id"`
}

// NativeEventID is the exact taskjournal native event identity, when present.
// ID remains the original source event ID. Equal wording is never an identity.
// Sequence orders only ChainID; registration/reporting/occurrence clocks are
// intentionally separate and untimed events remain null.
type TraceEntry struct {
	ID              string             `json:"id"`
	Kind            string             `json:"kind"`
	ChainID         string             `json:"chain_id"`
	Sequence        int                `json:"sequence"`
	PreviousSHA256  *string            `json:"previous_sha256"`
	Source          FileBinding        `json:"source"`
	Sources         []FileBinding      `json:"sources"`
	Role            string             `json:"role"`
	ActorClaim      string             `json:"actor_claim"`
	Recorder        string             `json:"recorder"`
	EvidenceClass   string             `json:"evidence_class"`
	OccurredAtUTC   *string            `json:"occurred_at_utc"`
	ReportedAtUTC   *string            `json:"reported_at_utc"`
	RegisteredAtUTC *string            `json:"registered_at_utc"`
	Candidate       *TraceCandidate    `json:"candidate"`
	NativeKind      *string            `json:"native_kind"`
	NativeID        *string            `json:"native_id"`
	NativeEventID   *string            `json:"native_event_id"`
	Outcome         *string            `json:"outcome"`
	Data            json.RawMessage    `json:"data"`
	Verification    *TraceVerification `json:"verification"`
}

type TraceVerification struct {
	AttemptID          string      `json:"attempt_id"`
	CandidateOID       string      `json:"candidate_oid"`
	CandidateTree      string      `json:"candidate_tree"`
	Argv               []string    `json:"argv"`
	CWD                string      `json:"cwd"`
	Acceptance         FileBinding `json:"acceptance"`
	AcceptanceSnapshot FileBinding `json:"acceptance_snapshot"`
	Review             FileBinding `json:"review"`
	Exit               *int        `json:"exit"`
	Stdout             FileBinding `json:"stdout"`
	Stderr             FileBinding `json:"stderr"`
	StartedAtUTC       *string     `json:"started_at_utc"`
	FinishedAtUTC      *string     `json:"finished_at_utc"`
	Error              string      `json:"error"`
	// InputsBound is true only when preserved input bytes match the receipt.
	// A present candidate alone cannot establish repeated verification basis.
	InputsBound bool `json:"inputs_bound"`
}

type traceCWD struct {
	workspace.FileSystem
	root string
}

func (f traceCWD) Getwd() (string, error) { return f.root, nil }

// ReadTraceHistory returns all safely Task-bound run families, in stable run-ID
// order. Unknown Tasks and unreadable registries are fatal. Damaged individual
// sources add diagnostics while independent runs/events remain inspectable.
// No timestamp is derived from mtime, and no current executable is read.
func ReadTraceHistory(d Dependencies, root, taskID string) (TraceHistory, error) {
	out := TraceHistory{Runs: []TraceRun{}, Reasons: []Reason{}}
	wd := d.Workspace
	if wd.Files == nil {
		wd.Files = workspace.SystemDependencies().Files
	}
	wd.Files = traceCWD{wd.Files, root}
	basis, err := workspace.ReadTaskJournalBasis(wd, workspace.TaskID(taskID))
	if err != nil {
		return out, err
	}
	if basis.Workspace.Root != root {
		return out, conflict("trace root is not the physical workspace root")
	}
	inventory, err := ReadInventory(d, root)
	if err != nil {
		return out, err
	}
	out.Reasons = append(out.Reasons, inventory.Reasons...)
	for _, item := range inventory.Runs {
		if item.TaskID == nil {
			for _, reason := range item.Reasons {
				out.Reasons = append(out.Reasons, Reason{reason.Code, item.RunID + ": " + reason.Detail})
			}
			continue
		}
		if *item.TaskID != taskID {
			continue
		}
		r := traceReader{d: d, basis: basis, captured: map[string]inventoryCapturedSource{}}
		r.run = TraceRun{RunID: item.RunID, TaskID: taskID, NativeRunIDs: []string{}, Sources: []FileBinding{}, Entries: []TraceEntry{}, Reasons: []Reason{}, HistoricalReasons: []Reason{}, DeclaredProcess: []string{}, Contract: json.RawMessage(`null`), Coverage: "complete"}
		for _, reason := range item.Reasons {
			if traceInventoryCoverageReason(reason.Code) {
				r.problem(reason.Code, reason.Detail)
			} else {
				r.run.HistoricalReasons = append(r.run.HistoricalReasons, reason)
			}
		}
		if item.Freshness == "stale" {
			r.run.Coverage = "stale"
		}
		if strings.HasPrefix(item.RunID, "trn_") {
			r.native(root, item.RunID)
		} else {
			r.workflow(root, item.RunID)
		}
		r.recheck()
		appendTraceRun(&out, r.run)
	}
	registry, err := readFile(basis.RegistryLocator, 32<<20, false)
	if err != nil || hash(registry) != "sha256:"+strings.TrimPrefix(basis.Registry.RawSHA256, "sha256:") {
		out.Reasons = append(out.Reasons, Reason{"trace_registry_changed", "Task registry changed during trace reading: " + basis.RegistryLocator})
		for i := range out.Runs {
			out.Runs[i].Coverage = "stale"
		}
	}
	return out, nil
}

func appendTraceRun(out *TraceHistory, run TraceRun) {
	if run.Family == "" || !digestPattern.MatchString(run.RequestSHA256) || run.FrozenBasis == nil {
		// Inventory and detail reads are separate observations, not a locked
		// snapshot. A disappearing or newly unbound request cannot supply a
		// fabricated typed run identity in an otherwise successful projection.
		out.Reasons = append(out.Reasons, Reason{"trace_run_identity_unavailable", run.RunID + ": run identity no longer has a safely bound request and preparation"})
		for _, reason := range run.Reasons {
			out.Reasons = append(out.Reasons, Reason{reason.Code, run.RunID + ": " + reason.Detail})
		}
		return
	}
	out.Runs = append(out.Runs, run)
}

func traceInventoryCoverageReason(code string) bool {
	switch code {
	case "run_evidence_unavailable", "run_source_orphaned", "delivery_state_unavailable", "delivery_phase_unknown", "delivery_binding_invalid", "delivery_goal_unavailable", "delivery_event_store_unavailable", "delivery_event_orphaned", "delivery_event_invalid", "delivery_event_unavailable", "delivery_report_invalid", "delivery_question_binding_invalid", "delivery_event_chain_invalid", "delivery_current_report_unavailable", "delivery_attempt_unresolved", "delivery_source_changed":
		return true
	default:
		return false
	}
}

type traceReader struct {
	d        Dependencies
	basis    workspace.TaskJournalBasis
	run      TraceRun
	captured map[string]inventoryCapturedSource
}

func (r *traceReader) problem(code, detail string) {
	r.run.Reasons = append(r.run.Reasons, Reason{code, detail})
	if r.run.Coverage != "stale" {
		r.run.Coverage = "partial"
	}
}

func (r *traceReader) read(path string, max int) ([]byte, error) {
	raw, err := readFile(path, max, true)
	if err != nil {
		return nil, err
	}
	binding := FileBinding{path, hash(raw)}
	if prior, ok := r.captured[path]; ok {
		if !bytes.Equal(prior.bytes, raw) {
			r.problem("trace_source_changed", "Source changed between reads: "+path)
			r.run.Coverage = "stale"
		}
	} else {
		r.captured[path] = inventoryCapturedSource{binding, raw, max}
		r.run.Sources = append(r.run.Sources, binding)
	}
	return raw, nil
}

func (r *traceReader) bound(b FileBinding, max int) ([]byte, error) {
	raw, err := r.read(b.Locator, max)
	if os.IsNotExist(err) {
		if retained, e := workspace.ResolveTaskRetainedEvidence(r.basis.Workspace.Root, r.basis.Task.ID, b.Locator, b.SHA256); e == nil {
			raw, err = r.read(retained, max)
		}
	}
	if err == nil && hash(raw) != b.SHA256 {
		err = integrity("bound source digest differs: " + b.Locator)
	}
	return raw, err
}

func (r *traceReader) recheck() {
	paths := make([]string, 0, len(r.captured))
	for path := range r.captured {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		before := r.captured[path]
		after, err := readFile(path, before.max, true)
		if err != nil || !bytes.Equal(before.bytes, after) {
			r.problem("trace_source_changed", "Source changed during trace reading: "+path)
			r.run.Coverage = "stale"
		}
	}
	sort.Slice(r.run.Sources, func(i, j int) bool { return r.run.Sources[i].Locator < r.run.Sources[j].Locator })
}

func (r *traceReader) bindDraft(raw []byte, preparationID, preparationSHA string, delivery bool) bool {
	draft, err := workflowhandoff.ValidateTaskRunDraft(raw)
	if delivery {
		draft, err = workflowhandoff.ValidateDeliveryTaskRunDraft(raw)
	}
	if err != nil || draft.Basis == nil {
		r.problem("trace_run_unbound", "Run draft has no valid Task binding")
		return false
	}
	// Validate this preparation's content closure without making unrelated
	// historical queues/preparations a prerequisite for observing this run.
	// The full registry was already structurally decoded by its native reader.
	selected := r.basis
	selected.Registry.TaskPreparations = nil
	selected.Registry.TaskQueueEvents = nil
	selected.Registry.EpicBaseUpdates = nil
	for _, p := range r.basis.Registry.TaskPreparations {
		if p.ID == preparationID {
			selected.Registry.TaskPreparations = append(selected.Registry.TaskPreparations, p)
		}
	}
	if len(selected.Registry.TaskPreparations) != 1 {
		r.problem("trace_run_unbound", "Run has no unique registered preparation: "+preparationID)
		return false
	}
	p := selected.Registry.TaskPreparations[0]
	if p.Outcome == nil || preparationSHA != digest(p) || p.Plan.TaskID != r.basis.Task.ID || p.Plan.Workspace.Root != r.basis.Workspace.Root || p.Plan.Workspace.MarkerSHA256 != r.basis.Workspace.MarkerSHA256 || draft.Basis.TaskID != r.basis.Task.ID || draft.Basis.TaskWorktreeID != p.Outcome.WorktreeID || draft.Worktree != p.Plan.WorktreePath || draft.Ref != "refs/heads/"+p.Plan.Branch || draft.OID != p.Plan.ParentOID || draft.Basis.SpecID != p.Plan.SpecID || draft.Basis.Spec != p.Plan.Spec || draft.Basis.Selection != p.Plan.Selection || draft.Basis.Problem != p.Plan.Problem || draft.Basis.Assessment != p.Plan.Assessment {
		r.problem("trace_run_unbound", fmt.Sprintf("Run request differs from exact historical preparation %s", preparationID))
		return false
	}
	if _, err := workspace.ValidateTaskJournalPreparation(r.d.Workspace, selected, preparationID); err != nil {
		// The exact registered Task/preparation remains identifiable when an
		// optional historical contract document is gone. Preserve independent
		// run events and mark the missing declaration instead of hiding them.
		r.problem("trace_preparation_content_unavailable", preparationID+": "+err.Error())
	}
	r.run.FrozenBasis = draft.Basis
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(raw, &fields)
	// Whitelist declared process fields; never expose recipient credentials,
	// reply capabilities, the whole mandate, or an arbitrary request object.
	contract := map[string]json.RawMessage{}
	for _, key := range []string{"goal", "authority", "budget", "procedure", "verifiers", "stop_conditions", "reporting"} {
		if value, ok := fields[key]; ok {
			contract[key] = value
		}
	}
	r.run.Contract = traceJSON(contract)
	var steps []struct {
		Instruction string `json:"instruction"`
	}
	if json.Unmarshal(fields["procedure"], &steps) == nil {
		for _, step := range steps {
			r.run.DeclaredProcess = append(r.run.DeclaredProcess, step.Instruction)
		}
	}
	return true
}

func (r *traceReader) contractField(key string, value any) {
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(r.run.Contract, &fields)
	if fields == nil {
		fields = map[string]json.RawMessage{}
	}
	fields[key] = traceJSON(value)
	r.run.Contract = traceJSON(fields)
}

func (r *traceReader) nativeHandoff(h WorkflowHandoff, draft json.RawMessage) {
	if h.Locator == "" {
		return
	}
	binding := FileBinding{h.Locator, h.SHA256}
	raw, err := r.bound(binding, 4<<20)
	var source struct {
		SourceDraftSHA256 string                   `json:"source_draft_sha256"`
		Basis             *workspace.TaskSpecBasis `json:"task_spec_binding"`
	}
	if err == nil {
		err = json.Unmarshal(raw, &source)
	}
	if err != nil || source.SourceDraftSHA256 != digest(draft) || !equal(source.Basis, r.run.FrozenBasis) {
		r.problem("trace_handoff_unbound", "Native handoff differs from exact frozen request draft: "+h.Locator)
		return
	}
	link, err := workflowhandoff.ReadTaskRunLink(r.d.Workflow, h.Locator)
	if err != nil || link.HandoffID != h.ID || link.SHA256 != h.SHA256 || link.Locator != h.Locator {
		r.problem("trace_handoff_unavailable", "Native handoff identity could not be validated: "+h.Locator)
		return
	}
	r.nativeAlias(link.RunID)
}

func (r *traceReader) nativeAlias(id string) {
	if id == "" {
		return
	}
	for _, prior := range r.run.NativeRunIDs {
		if prior == id {
			return
		}
	}
	r.run.NativeRunIDs = append(r.run.NativeRunIDs, id)
	sort.Strings(r.run.NativeRunIDs)
}

func traceJSON(v any) json.RawMessage { b, _ := Canonical(v); return b }

func traceEntry(id, kind, chain string, sequence int, source FileBinding) TraceEntry {
	return TraceEntry{ID: id, Kind: kind, ChainID: chain, Sequence: sequence, Source: source, Sources: []FileBinding{source}, Role: "unknown", Recorder: "ply", EvidenceClass: "native_record", Data: json.RawMessage(`{}`)}
}

func traceNative(e *TraceEntry, kind, id string, value any) {
	e.NativeKind, e.NativeID = ptr(kind), ptr(id)
	// taskjournal uses unprefixed canonical digests in nativeID.
	recordHash := strings.TrimPrefix(digest(value), "sha256:")
	e.NativeEventID = ptr("native_" + strings.TrimPrefix(digest([]string{kind, id, recordHash}), "sha256:"))
}
