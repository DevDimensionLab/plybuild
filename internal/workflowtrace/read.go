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
	// Journal history needs this Task's content closure, not the readiness of
	// historical queues and Epic base updates for other work.
	selected.TaskQueueEvents, selected.EpicBaseUpdates = nil, nil
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
