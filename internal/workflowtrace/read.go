package workflowtrace

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Read observes only registered stores. All optional histories degrade
// independently; discovery, unknown Task and authoritative registry errors fail.
func Read(d workspace.Dependencies, taskID string) (Result, error) {
	return read(d, taskID, nil)
}

func read(d workspace.Dependencies, taskID string, afterHistory func()) (Result, error) {
	id, err := workspace.ParseTaskID(taskID)
	if err != nil {
		return Result{}, err
	}
	basis, err := workspace.ReadViewBasis(d)
	if err != nil {
		return Result{}, err
	}
	var task *workspace.TaskRecord
	for i := range basis.Registry.Tasks {
		if basis.Registry.Tasks[i].ID == id {
			task = &basis.Registry.Tasks[i]
			break
		}
	}
	if task == nil {
		return Result{}, fmt.Errorf("workflow_trace_task_not_found: Task %s is not registered", taskID)
	}
	now := time.Now()
	if d.WorkClock != nil {
		now = d.WorkClock.Now()
	}
	b := newBuilder(basis.Workspace.Root, *task, now.UTC().Format(time.RFC3339Nano))
	journalPath := filepath.Join(basis.Workspace.Root, ".ply", "task-process", "v1", taskID)
	journalBefore, journalBeforeErr := directoryDigest(journalPath)
	if title := basis.Titles[id]; title.Title != nil {
		b.out.Task.Title = *title.Title
	}
	for path, digest := range basis.Sources {
		b.source(Source{Kind: "registration", Locator: path, SHA256: nonempty(digest), Status: "valid"})
	}
	lifecycle, err := workspace.ReadWorkItemLifecycle(d, basis.Workspace.Root)
	if err != nil {
		b.diagnostic("lifecycle_unavailable", err.Error(), nil, nil)
	}
	selected := basis.Registry
	selected.Tasks = []workspace.TaskRecord{*task}
	// Journal history needs this Task's exact preparation/content closure, not
	// the readiness of historical queues and Epic base updates for other work.
	selected.TaskQueueEvents, selected.EpicBaseUpdates = nil, nil
	selected.TaskPreparations = []workspace.TaskPreparation{}
	for _, p := range basis.Registry.TaskPreparations {
		if p.Plan.TaskID == id {
			selected.TaskPreparations = append(selected.TaskPreparations, p)
		}
	}
	batch, err := taskjournal.ReadEventBatch(d, taskjournal.EventBatchInput{Workspace: basis.Workspace, Projects: basis.Projects, Registry: selected, Lifecycle: lifecycle})
	if err != nil {
		b.diagnostic("journal_unavailable", err.Error(), nil, nil)
	} else {
		b.journal(batch.Tasks[id], basis.Registry)
		for _, r := range batch.Reasons {
			b.diagnostic(r.Code, r.Detail, r.SourceIDs, nil)
		}
	}
	b.goals(d, basis, *task)
	history, err := taskrun.ReadTraceHistory(taskrun.SystemDependencies(d), basis.Workspace.Root, taskID)
	if err != nil {
		b.diagnostic("run_history_unavailable", err.Error(), nil, nil)
	} else {
		b.history(history)
	}
	b.linkNative(basis.Registry)
	if afterHistory != nil {
		afterHistory()
	}
	b.recheckSources()
	journalAfter, journalAfterErr := directoryDigest(journalPath)
	if journalBeforeErr != nil || journalAfterErr != nil || journalBefore != journalAfter {
		b.out.Freshness = "unknown"
		b.diagnostic("journal_store_changed_or_unavailable", "Journal entries changed or could not be safely compared during reading: "+journalPath, nil, nil)
	}
	fresh, reasons := workspace.CheckReadViewBasis(d, basis)
	if fresh != "fresh" {
		b.out.Freshness = "unknown"
		for _, reason := range reasons {
			b.diagnostic(reason, "Registered sources changed during this read.", nil, nil)
		}
	}
	return b.finish(), nil
}

func directoryDigest(path string) (string, error) {
	physical, err := filepath.EvalSymlinks(path)
	if os.IsNotExist(err) {
		return "absent", nil
	}
	if err != nil {
		return "", err
	}
	if physical != path {
		return "", fmt.Errorf("journal directory is not physical")
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	names := []string{}
	for _, entry := range entries {
		if entry.Name() == "append.lock" || strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		names = append(names, entry.Name()+"/"+entry.Type().String())
	}
	return hash(names), nil
}

func (b *builder) goals(d workspace.Dependencies, basis workspace.ReadViewBasisSnapshot, task workspace.TaskRecord) {
	latest := map[string]Goal{}
	unclassified := map[string]int{}
	j := workspace.TaskJournalBasis{Workspace: basis.Workspace, Task: task, Registry: basis.Registry, RegistryLocator: filepath.Join(basis.Workspace.Root, workspace.MarkerDirectory, workspace.WorkItemsFile)}
	for _, ref := range basis.Registry.TaskSpecRevisions {
		if ref.TaskID != task.ID {
			continue
		}
		var sources []workspace.TaskJournalSource
		var err error
		found := false
		for _, pub := range basis.Registry.TaskContentPublications {
			if pub.TaskID == task.ID && pub.OutcomeRef.ManifestSHA256 == ref.ManifestSHA256 {
				found = true
				sources, err = workspace.ReadTaskJournalContent(d, j, pub)
				break
			}
		}
		if !found {
			err = fmt.Errorf("Spec %s@%d has no publication binding", ref.SpecID, ref.Revision)
		}
		ids := []string{}
		for _, src := range sources {
			ids = append(ids, b.source(Source{Kind: src.Kind, Locator: src.Locator, SHA256: ptr(src.SHA256), Status: "valid"}))
		}
		var fields map[string]json.RawMessage
		if err == nil {
			var data []byte
			data, err = d.TaskContent.Read(basis.Workspace.Root, "manifests", ref.ManifestSHA256)
			if err == nil {
				err = json.Unmarshal(data, &fields)
			}
		}
		if err != nil {
			b.diagnostic("goal_source_unavailable", err.Error(), ids, nil)
			if ref.Revision > unclassified[ref.SpecID] {
				unclassified[ref.SpecID] = ref.Revision
			}
			continue
		}
		var kind string
		_ = json.Unmarshal(fields["contract_kind"], &kind)
		if kind != "goal" {
			continue
		}
		if old, ok := latest[ref.SpecID]; !ok || ref.Revision > old.Revision {
			latest[ref.SpecID] = Goal{ref.SpecID, ref.Revision, ref.ManifestSHA256, "valid", unique(ids), raw(fields)}
		}
	}
	for id, g := range latest {
		if unclassified[id] > g.Revision {
			g.Status = "unknown"
			b.diagnostic("goal_head_unknown", "A newer Spec cannot be classified; the last readable goal may not be current.", g.SourceIDs, nil)
		}
		b.out.CurrentGoals = append(b.out.CurrentGoals, g)
	}
	sort.Slice(b.out.CurrentGoals, func(i, j int) bool { return b.out.CurrentGoals[i].SpecID < b.out.CurrentGoals[j].SpecID })
}

func (b *builder) history(h taskrun.TraceHistory) {
	for _, reason := range h.Reasons {
		b.diagnostic(reason.Code, reason.Detail, nil, nil)
	}
	for _, in := range h.Runs {
		for _, alias := range in.NativeRunIDs {
			b.runAliases[alias] = append(b.runAliases[alias], in.RunID)
		}
		r := Run{ID: in.RunID, Family: in.Family, Provider: nonempty(in.Provider), SessionID: nonempty(in.SessionID), RequestSHA256: in.RequestSHA256, FrozenBasis: in.FrozenBasis, FrozenGoal: raw(in.FrozenGoal), DeclaredProcess: raw(map[string]any{"contract": in.Contract, "obligations": in.DeclaredProcess}), Coverage: in.Coverage, SourceIDs: []string{}, EventIDs: []string{}}
		if r.Coverage != "complete" {
			r.Coverage = "partial"
		}
		for _, s := range in.Sources {
			r.SourceIDs = append(r.SourceIDs, b.source(Source{Kind: "run_history", Locator: s.Locator, SHA256: nonempty(s.SHA256), Status: "valid"}))
		}
		for _, reason := range in.Reasons {
			b.diagnostic(reason.Code, in.RunID+": "+reason.Detail, r.SourceIDs, nil)
		}
		for _, reason := range in.HistoricalReasons {
			b.out.Diagnostics = append(b.out.Diagnostics, Diagnostic{reason.Code, in.RunID + " (historical observation): " + reason.Detail, unique(r.SourceIDs), []string{}})
		}
		for _, x := range in.Entries {
			id := in.RunID + "/" + x.ID
			if x.NativeEventID != nil {
				id = *x.NativeEventID
			}
			e := Event{ID: id, NativeID: x.NativeID, Type: x.Kind, Title: x.Kind, Role: x.Role, ActorClaim: nonempty(x.ActorClaim), Recorder: x.Recorder, EvidenceClass: x.EvidenceClass, RunIDs: []string{in.RunID}, Outcome: x.Outcome, OccurredAtUTC: x.OccurredAtUTC, ReportedAtUTC: x.ReportedAtUTC, RegisteredAtUTC: x.RegisteredAtUTC, TimeBasis: taskjournal.TimeBasis{Kind: "unknown", Precision: "unknown"}, SourceIDs: []string{}, Positions: []Position{}, Data: x.Data}
			if e.NativeID == nil {
				e.NativeID = ptr(x.ID)
			}
			if e.Role == "" {
				e.Role = "unknown"
			}
			switch e.EvidenceClass {
			case "controlled_verification", "local_effect_observation":
				e.EvidenceClass = "controlled"
			case "agent_report", "review_claim":
				e.EvidenceClass = "reported"
			case "reported", "human_attestation", "controlled", "registered", "unknown":
			case "native_record", "declared_contract":
				e.EvidenceClass = "registered"
			default:
				e.EvidenceClass = "unknown"
			}
			if e.ReportedAtUTC != nil {
				e.TimeBasis = taskjournal.TimeBasis{Kind: "reported", Clock: ptr("native attestation"), Precision: "nanosecond"}
			}
			if e.OccurredAtUTC != nil {
				e.TimeBasis = taskjournal.TimeBasis{Kind: "observed", Clock: ptr("native verifier"), Precision: "nanosecond"}
			}
			if x.Candidate != nil {
				e.CandidateID = b.candidate(x.Candidate.OID, x.Candidate.Tree)
				e.ResultID = x.Candidate.TaskResultID
			}
			for _, s := range append(append([]taskrun.FileBinding{}, x.Sources...), x.Source) {
				if s.Locator != "" {
					e.SourceIDs = append(e.SourceIDs, b.source(Source{Kind: "run_history", Locator: s.Locator, SHA256: nonempty(s.SHA256), Status: "valid"}))
				}
			}
			if x.ChainID != "" {
				e.Positions = append(e.Positions, Position{x.ChainID, x.Sequence})
			}
			if x.Verification != nil {
				b.verifiers[id] = *x.Verification
				e.Data = raw(map[string]any{"record": x.Data, "verification": x.Verification})
			}
			b.event(e)
			r.EventIDs = append(r.EventIDs, id)
		}
		r.SourceIDs = unique(r.SourceIDs)
		r.EventIDs = unique(r.EventIDs)
		b.out.Runs = append(b.out.Runs, r)
	}
}

func (b *builder) linkNative(registry workspace.WorkItemRegistry) {
	results := map[string]string{}
	qas := map[string]string{}
	runAliases := map[string][]string{}
	for _, e := range b.out.Events {
		if e.Type == "task_result_technical" && e.ResultID != nil {
			results[*e.ResultID] = e.ID
			runAliases[*e.ResultID] = append(runAliases[*e.ResultID], e.RunIDs...)
		}
		if e.Type == "human_qa" && e.NativeID != nil {
			qas[*e.NativeID] = e.ID
		}
	}
	for i := range b.out.Events {
		e := &b.out.Events[i]
		if e.ResultID == nil {
			continue
		}
		e.RunIDs = unique(append(e.RunIDs, runAliases[*e.ResultID]...))
		if id := results[*e.ResultID]; id != "" && e.ID != id && (e.Type == "human_qa" || e.Type == "task_result_reported") {
			b.out.Relations = append(b.out.Relations, Relation{"applies_to", e.ID, id, "exact TaskResult identity", e.SourceIDs})
		}
	}
	for _, ir := range registry.IntegrationResults {
		for _, a := range registry.IntegrationAuthorities {
			if a.ID != ir.AuthorityID || string(a.TaskID) != b.out.Task.ID {
				continue
			}
			from := nativeID("integration", string(ir.ID), ir)
			if _, ok := b.events[from]; !ok {
				continue
			}
			if to := qas[string(a.HumanQARecordID)]; to != "" {
				b.out.Relations = append(b.out.Relations, Relation{"integrates_after_qa", from, to, "exact integration authority QA binding", []string{}})
			}
		}
	}
	known := map[string]bool{}
	for _, r := range b.out.Runs {
		known[r.ID] = true
	}
	for i := range b.out.Events {
		e := &b.out.Events[i]
		resolved, aliases := []string{}, []string{}
		for _, id := range e.RunIDs {
			if known[id] {
				resolved = append(resolved, id)
			} else {
				aliases = append(aliases, id)
				resolved = append(resolved, b.runAliases[id]...)
			}
		}
		e.RunIDs = unique(resolved)
		if len(aliases) > 0 {
			var data map[string]json.RawMessage
			_ = json.Unmarshal(e.Data, &data)
			if data == nil {
				data = map[string]json.RawMessage{}
			}
			data["native_run_ids"] = raw(unique(aliases))
			e.Data = raw(data)
		}
		// Exact alias binding also makes independent registry facts inspectable
		// from their outer run without manufacturing a second execution.
		for j := range b.out.Runs {
			for _, id := range e.RunIDs {
				if b.out.Runs[j].ID == id {
					b.out.Runs[j].EventIDs = unique(append(b.out.Runs[j].EventIDs, e.ID))
				}
			}
		}
	}
}

// Recheck every captured source once, comparing bytes without deriving clocks
// from filesystem metadata. Foreign symlinks and oversized files remain unknown.
func (b *builder) recheckSources() {
	cache := map[string]string{}
	fail := map[string]bool{}
	for i := range b.out.Sources {
		s := &b.out.Sources[i]
		if s.Status != "valid" || s.SHA256 == nil {
			continue
		}
		if _, ok := cache[s.Locator]; !ok && !fail[s.Locator] {
			h, err := sourceDigest(s.Locator)
			if err != nil {
				fail[s.Locator] = true
			} else {
				cache[s.Locator] = h
			}
		}
		if fail[s.Locator] {
			s.Status = "unknown"
			b.diagnostic("source_unavailable_at_recheck", "Source could not be safely reread: "+s.Locator, []string{s.ID}, nil)
		} else if cache[s.Locator] != strings.TrimPrefix(*s.SHA256, "sha256:") {
			s.Status = "changed"
			b.diagnostic("source_changed", "Source bytes differ from their captured binding: "+s.Locator, []string{s.ID}, nil)
		}
	}
}

func sourceDigest(path string) (string, error) {
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	if physical != path {
		return "", fmt.Errorf("source path is not physical")
	}
	before, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !before.Mode().IsRegular() || before.Size() > 64<<20 {
		return "", fmt.Errorf("source is not a bounded regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(before, opened) {
		return "", fmt.Errorf("source changed before read")
	}
	h := sha256.New()
	n, err := io.Copy(h, io.LimitReader(f, 64<<20+1))
	if err != nil {
		return "", err
	}
	after, err := os.Lstat(path)
	if err != nil || !os.SameFile(before, after) || n != before.Size() || after.Size() != before.Size() || !after.ModTime().Equal(before.ModTime()) {
		return "", fmt.Errorf("source changed during read")
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
