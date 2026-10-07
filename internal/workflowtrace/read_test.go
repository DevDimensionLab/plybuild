package workflowtrace

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type cwdFiles struct {
	workspace.FileSystem
	root string
}

func (f cwdFiles) Getwd() (string, error) { return f.root, nil }

type readOnlyFiles struct{ workspace.FileSystem }

func (readOnlyFiles) MkdirTemp(string, string) (string, error) {
	panic("trace must not create directories")
}
func (readOnlyFiles) OpenFile(string, int, fs.FileMode) (workspace.File, error) {
	panic("trace must not open for mutation")
}
func (readOnlyFiles) Chmod(string, fs.FileMode) error { panic("trace must not chmod") }
func (readOnlyFiles) Rename(string, string) error     { panic("trace must not rename") }
func (readOnlyFiles) RemoveAll(string) error          { panic("trace must not delete") }

type noGit struct{ workspace.WorkItemGit }

func (noGit) ObserveWorktree(string) (workspace.GitWorktreeObservation, error) {
	panic("trace Git probe")
}
func (noGit) ObserveRepo(workspace.RepoRecord) (workspace.GitRepoObservation, error) {
	panic("trace Git probe")
}
func (noGit) ObserveRef(workspace.RepoRecord, string) (workspace.GitRefObservation, error) {
	panic("trace Git probe")
}
func (noGit) ListWorktrees(workspace.RepoRecord) (workspace.GitWorktreeInventory, error) {
	panic("trace Git probe")
}

func readFixture(t *testing.T) (workspace.Dependencies, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(root, "repo")
	if err = os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		prefix := []string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "core.fsmonitor=false", "-c", "gc.auto=0", "-c", "maintenance.auto=false", "-c", "user.name=Trace fixture", "-c", "user.email=trace@example.invalid", "-C", repo}
		out, err := exec.Command("git", append(prefix, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("fixture Git %v: %s", err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("init", "-b", "epic")
	if err = os.WriteFile(filepath.Join(repo, "README"), []byte("fixture\n"), 0600); err != nil {
		t.Fatal(err)
	}
	git("add", "README")
	git("commit", "-m", "fixture")
	d := workspace.SystemDependencies()
	d.Files = cwdFiles{d.Files, root}
	if _, err = workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.AddProject(d, workspace.ProjectAddInput{ProjectID: "project", Name: "Trace fixture", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.AdoptEpic(d, workspace.EpicAdoptInput{EpicID: "epic", Title: "Epic", ProjectID: "project", RepoID: "repo", Worktree: repo, Ref: "refs/heads/epic", ExpectedOID: git("rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.CreateTask(d, workspace.TaskCreateInput{TaskID: "task", Title: "Untouched Task", Description: "Read fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
		t.Fatal(err)
	}
	return d, root
}
func fingerprint(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := e.Info()
		if err != nil {
			return err
		}
		out[path] = info.Mode().String()
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out[path] += fmt.Sprintf("/%x", sha256.Sum256(data))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func appendJournal(t *testing.T, d workspace.Dependencies, root, key string) taskjournal.AppendResult {
	t.Helper()
	evidence := filepath.Join(root, "evidence")
	if err := os.WriteFile(evidence, []byte("recorded evidence"), 0600); err != nil {
		t.Fatal(err)
	}
	input := taskjournal.EventInput{Common: taskjournal.Common{Kind: "ply.workspace.task-journal-event-input", SchemaVersion: 1, PublicationKey: key, TaskID: "task", Actor: taskjournal.Actor{ID: "Developer claim", Role: "developer"}, TimeBasis: taskjournal.TimeBasis{Kind: "unknown", Precision: "unknown"}, Sources: []taskjournal.SourceRef{{Locator: evidence, SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte("recorded evidence")))}}}, Type: "note", Title: "Reported work", Relations: []taskjournal.Relation{}}
	file := filepath.Join(root, key+".json")
	if err := os.WriteFile(file, raw(input), 0600); err != nil {
		t.Fatal(err)
	}
	result, err := taskjournal.New(d).Append("task", file, "event")
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestTraceReadBoundaryEmptyUntouchedAndUnknownTask(t *testing.T) {
	d, root := readFixture(t)
	d.WorkGit = noGit{}
	d.Files = readOnlyFiles{d.Files}
	before := fingerprint(t, root)
	for i := 0; i < 2; i++ {
		r, err := Read(d, "task")
		if err != nil {
			t.Fatal(err)
		}
		if r.HistoryState != "empty" || len(r.Runs) != 0 || r.LiveState != "unknown" || r.NextTransitionAuthorized {
			t.Fatalf("empty trace %+v", r)
		}
		if err = ValidateReferences(r); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Read(d, "missing"); err == nil {
		t.Fatal("unknown Task succeeded")
	}
	if !reflect.DeepEqual(before, fingerprint(t, root)) {
		t.Fatal("read changed directories, modes or source bytes")
	}
	other, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	outside := d
	outside.Files = readOnlyFiles{cwdFiles{workspace.SystemDependencies().Files, other}}
	snapshot := fingerprint(t, other)
	if _, err = Read(outside, "task"); err == nil {
		t.Fatal("absent workspace succeeded")
	}
	if !reflect.DeepEqual(snapshot, fingerprint(t, other)) {
		t.Fatal("read initialized absent workspace")
	}
}

func TestTraceReadPreservesJournalPrefixMissingSourcesAndMtimeNulls(t *testing.T) {
	d, root := readFixture(t)
	first := appendJournal(t, d, root, "first")
	second := appendJournal(t, d, root, "second")
	d.WorkGit = noGit{}
	d.Files = readOnlyFiles{d.Files}
	a, err := Read(d, "task")
	if err != nil {
		t.Fatal(err)
	}
	e := eventByID(t, a, first.EventID)
	if e.Role != "agent" || e.OccurredAtUTC != nil || e.ReportedAtUTC != nil || e.RegisteredAtUTC == nil {
		t.Fatalf("journal clocks/actor %+v", e)
	}
	if err = os.Chtimes(first.RecordLocator, time.Now().Add(time.Hour), time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	after, err := Read(d, "task")
	if err != nil {
		t.Fatal(err)
	}
	again := eventByID(t, after, first.EventID)
	if !reflect.DeepEqual(e, again) {
		t.Fatal("mtime changed recorded event")
	}
	if err = os.WriteFile(second.RecordLocator, []byte("corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(filepath.Join(root, "evidence")); err != nil {
		t.Fatal(err)
	}
	before := fingerprint(t, root)
	partial, err := Read(d, "task")
	if err != nil {
		t.Fatal(err)
	}
	if partial.Coverage.State != "partial" || partial.Freshness != "unknown" {
		t.Fatalf("damaged history %+v", partial.Coverage)
	}
	eventByID(t, partial, first.EventID)
	if len(partial.Events) < 2 {
		t.Fatal("journal damage hid independent native Task facts")
	}
	if !reflect.DeepEqual(before, fingerprint(t, root)) {
		t.Fatal("read repaired damaged sources")
	}
}

func TestTraceSourceRecheckRejectsChangedAndSymlinkedEvidence(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "source")
	if err = os.WriteFile(p, []byte("before"), 0600); err != nil {
		t.Fatal(err)
	}
	digest, err := sourceDigest(p)
	if err != nil {
		t.Fatal(err)
	}
	b := testBuilder()
	b.source(Source{ID: "changing", Locator: p, SHA256: ptr(digest), Status: "valid"})
	if err = os.WriteFile(p, []byte("after"), 0600); err != nil {
		t.Fatal(err)
	}
	b.recheckSources()
	if b.out.Sources[b.sources["changing"]].Status != "changed" {
		t.Fatal("concurrent changed bytes passed")
	}
	alias := filepath.Join(root, "alias")
	if err = os.Symlink(p, alias); err != nil {
		t.Fatal(err)
	}
	if _, err = sourceDigest(alias); err == nil {
		t.Fatal("source symlink followed")
	}
}

func TestTraceCurrentGoalRemainsSeparateFromFrozenRun(t *testing.T) {
	d, root := readFixture(t)
	registry, err := d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	problem := registry.TaskProblemRevisions[0]
	doc := filepath.Join(root, "design.md")
	body := []byte("# Trace\nRecorded facts only.\n")
	if err = os.WriteFile(doc, body, 0600); err != nil {
		t.Fatal(err)
	}
	var previous *workspace.TaskRevisionRef
	var first workspace.TaskGoalRef
	for n := 1; n <= 2; n++ {
		draft := map[string]any{"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "publication_key": fmt.Sprintf("goal/%d", n), "task_id": "task", "registry_upgrade": nil, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "recorder": map[string]any{"actor_claim": "Synthetic planner", "control_surface": "fixture", "recorded_at_utc": "2026-10-07T10:00:00Z"}, "spec_id": "trace-goal", "expected_previous": previous, "problem": workspace.TaskRevisionRef{Revision: problem.Revision, ManifestSHA256: problem.ManifestSHA256}, "title": fmt.Sprintf("Goal %d", n), "contract_kind": "goal", "objective": fmt.Sprintf("Observe revision %d", n), "requirements": []map[string]string{{"id": "read", "acceptance": "Read recorded sources without writing."}}, "design": []map[string]any{{"document_id": "design", "section": nil}}, "constraints": []string{"Read only."}, "executor": map[string]any{"provider": "codex", "model": nil, "effort": nil}, "change_reason": "Fixture goal revision.", "documents": []map[string]any{{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": fmt.Sprintf("sha256:%x", sha256.Sum256(body)), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}}}
		file := filepath.Join(root, fmt.Sprintf("goal-%d.json", n))
		if err = os.WriteFile(file, raw(draft), 0600); err != nil {
			t.Fatal(err)
		}
		record, err := workspace.RecordTaskSpec(d, workspace.TaskContentInput{TaskID: "task", File: file})
		if err != nil {
			t.Fatal(err)
		}
		previous = &workspace.TaskRevisionRef{Revision: *record.OutcomeRef.Revision, ManifestSHA256: record.OutcomeRef.ManifestSHA256}
		if n == 1 {
			first = workspace.TaskGoalRef{SpecID: "trace-goal", Spec: *previous}
		}
	}
	d.WorkGit = noGit{}
	d.Files = readOnlyFiles{d.Files}
	r, err := Read(d, "task")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.CurrentGoals) != 1 || r.CurrentGoals[0].Revision != 2 || !strings.Contains(string(r.CurrentGoals[0].Declaration), "Observe revision 2") {
		t.Fatalf("current goal %+v", r.CurrentGoals)
	}
	b := testBuilder()
	b.out.CurrentGoals = r.CurrentGoals
	b.history(taskrun.TraceHistory{Runs: []taskrun.TraceRun{{RunID: "frozen", Family: "goal_execute_v2", FrozenGoal: &first, Contract: json.RawMessage(`{"procedure":["original procedure"]}`), Coverage: "complete"}}})
	combined := b.finish()
	if !strings.Contains(string(combined.Runs[0].FrozenGoal), first.Spec.ManifestSHA256) || combined.CurrentGoals[0].ManifestSHA256 == first.Spec.ManifestSHA256 {
		t.Fatal("current goal overwrote frozen goal")
	}
}

func TestTraceFatalRegistryDoesNotReturnPartialSuccess(t *testing.T) {
	d, root := readFixture(t)
	if err := os.WriteFile(filepath.Join(root, ".ply", "work-items.yaml"), []byte("not: [valid"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(d, "task"); err == nil {
		t.Fatal("fatal registry error hidden")
	}
}

func TestTraceConcurrentJournalAppendMakesCapturedPrefixPartial(t *testing.T) {
	d, root := readFixture(t)
	first := appendJournal(t, d, root, "first")
	observed, err := read(d, "task", func() { appendJournal(t, d, root, "second") })
	if err != nil {
		t.Fatal(err)
	}
	if observed.Coverage.State != "partial" || observed.Freshness != "unknown" {
		t.Fatalf("concurrent append hidden: %+v", observed.Coverage)
	}
	eventByID(t, observed, first.EventID)
	for _, diag := range observed.Diagnostics {
		if diag.Code == "journal_store_changed_or_unavailable" {
			return
		}
	}
	t.Fatal("missing directory-change diagnostic")
}
