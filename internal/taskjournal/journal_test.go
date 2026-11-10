package taskjournal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type cwdFiles struct {
	workspace.FileSystem
	root string
}

func (f cwdFiles) Getwd() (string, error) { return f.root, nil }

type testFixture struct {
	s     Service
	root  string
	input EventInput
}

func fixture(t *testing.T) testFixture {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	repo := filepath.Join(root, "repo")
	if e = os.Mkdir(repo, 0700); e != nil {
		t.Fatal(e)
	}
	git := func(args ...string) string {
		t.Helper()
		prefix := []string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "gc.auto=0", "-c", "maintenance.auto=false", "-c", "submodule.recurse=false", "-c", "user.name=Synthetic Journal", "-c", "user.email=journal@example.invalid", "-C", repo}
		b, e := exec.Command("git", append(prefix, args...)...).CombinedOutput()
		if e != nil {
			t.Fatalf("fixture git: %v: %s", e, b)
		}
		return strings.TrimSpace(string(b))
	}
	git("init", "-b", "epic")
	write(t, filepath.Join(repo, "file"), []byte("fixture\n"))
	git("add", "file")
	git("commit", "-m", "synthetic fixture")
	d := workspace.SystemDependencies()
	d.Files = cwdFiles{d.Files, root}
	if _, e = workspace.Init(d); e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.AddProject(d, workspace.ProjectAddInput{ProjectID: "project", Name: "Fixture", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.AdoptEpic(d, workspace.EpicAdoptInput{EpicID: "epic", Title: "Fixture Epic", ProjectID: "project", RepoID: "repo", Worktree: repo, Ref: "refs/heads/epic", ExpectedOID: git("rev-parse", "HEAD")}); e != nil {
		t.Fatal(e)
	}
	if _, e = workspace.CreateTask(d, workspace.TaskCreateInput{TaskID: "task", Title: "Fixture Task", Description: "Process journal fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); e != nil {
		t.Fatal(e)
	}
	source := filepath.Join(root, "evidence.txt")
	write(t, source, []byte("synthetic evidence\n"))
	in := EventInput{Common: Common{Kind: "ply.workspace.task-journal-event-input", SchemaVersion: 1, PublicationKey: "test:1", TaskID: "task", Actor: Actor{ID: "developer", Role: "developer"}, OccurredAt: ptr("2026-10-03T10:00:00Z"), TimeBasis: TimeBasis{"reported", ptr("fixture"), "second", nil}, Sources: []SourceRef{{source, hash([]byte("synthetic evidence\n"))}}}, Type: "note", Title: "Synthetic statement", Relations: []Relation{}}
	return testFixture{New(d), root, in}
}
func write(t *testing.T, path string, b []byte) {
	t.Helper()
	if e := os.WriteFile(path, b, 0600); e != nil {
		t.Fatal(e)
	}
}
func (f testFixture) file(t *testing.T, v any) string {
	t.Helper()
	b, e := Canonical(v)
	if e != nil {
		t.Fatal(e)
	}
	p := filepath.Join(f.root, "input-"+hash(b)+".json")
	write(t, p, b)
	return p
}
func (f testFixture) append(t *testing.T, in EventInput) AppendResult {
	t.Helper()
	r, e := f.s.Append("task", f.file(t, in), "event")
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func (f testFixture) show(t *testing.T) Snapshot {
	t.Helper()
	s, e := f.s.Show("task", Options{})
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func nativeFingerprint(t *testing.T, root string) string {
	t.Helper()
	out := map[string]any{}
	e := filepath.Walk(filepath.Join(root, ".ply"), func(p string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if strings.Contains(p, "/task-process") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode().IsRegular() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			out[p] = []any{hash(b), info.Mode(), info.ModTime().UnixNano()}
		}
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return digest(out)
}

func TestNativeReadAndOfflineAreReadOnly(t *testing.T) {
	f := fixture(t)
	before := nativeFingerprint(t, f.root)
	s := f.show(t)
	if len(s.Events) == 0 || s.Events[0].Origin != "native" {
		t.Fatal("missing native facts")
	}
	for _, e := range s.Events {
		if e.OccurredAt != nil {
			t.Fatal("native state invented occurrence time")
		}
	}
	if _, e := os.Lstat(storePath(f.root, "task")); !os.IsNotExist(e) {
		t.Fatal("show created storage")
	}
	p := f.file(t, s)
	offline, e := ReadSnapshot(p, "task", Options{ActorID: "absent", Order: "recorded"})
	if e != nil {
		t.Fatal(e)
	}
	if len(offline.Selection.EventIDs) != 0 || offline.SnapshotID != s.SnapshotID || offline.AsOf != s.AsOf {
		t.Fatal("offline selection changed evidence")
	}
	if before != nativeFingerprint(t, f.root) {
		t.Fatal("read changed native bytes or mtime")
	}
}

func TestStrictInputRejectedBeforeJournalCreation(t *testing.T) {
	f := fixture(t)
	b := raw(f.input)
	var m map[string]any
	json.Unmarshal(b, &m)
	cases := map[string][]byte{"duplicate": []byte(strings.Replace(string(b), `"title":`, `"title":"duplicate","title":`, 1)), "extra": append(append([]byte{}, b...), []byte(" {}")...), "utf8": []byte("{\"x\":\"\xff\"}"), "oversize": []byte(strings.Repeat(" ", inputLimit+1) + string(b))}
	for label, change := range map[string]func(map[string]any){"unknown": func(x map[string]any) { x["payload"] = "secret-like opaque content is not allowed" }, "missing": func(x map[string]any) { delete(x, "activity_id") }, "null-list": func(x map[string]any) { x["relations"] = nil }, "version": func(x map[string]any) { x["schema_version"] = 2 }, "unknown-clock": func(x map[string]any) {
		x["time_basis"] = map[string]any{"kind": "unknown", "clock": "invented", "precision": "unknown", "uncertainty": nil}
	}, "control": func(x map[string]any) { x["title"] = "\x1b[2J" }, "long-id": func(x map[string]any) { x["activity_id"] = strings.Repeat("x", 257) }} {
		var x map[string]any
		json.Unmarshal(b, &x)
		change(x)
		cases[label] = raw(x)
	}
	for label, data := range cases {
		t.Run(label, func(t *testing.T) {
			p := filepath.Join(f.root, label+".json")
			write(t, p, data)
			if _, e := f.s.Append("task", p, "event"); e == nil {
				t.Fatal("invalid input accepted")
			}
		})
	}
	if _, e := os.Lstat(filepath.Join(f.root, ".ply", "task-process")); !os.IsNotExist(e) {
		t.Fatal("invalid input wrote journal")
	}
}

func TestConcurrentKeysAndLostResponse(t *testing.T) {
	f := fixture(t)
	inputs := []string{}
	for i := 0; i < 16; i++ {
		in := f.input
		in.PublicationKey = fmt.Sprintf("key:%d", i/2)
		inputs = append(inputs, f.file(t, in))
	}
	before := nativeFingerprint(t, f.root)
	var wg sync.WaitGroup
	results := make(chan AppendResult, 16)
	errs := make(chan error, 16)
	for _, p := range inputs {
		wg.Add(1)
		go func(p string) { defer wg.Done(); r, e := f.s.Append("task", p, "event"); results <- r; errs <- e }(p)
	}
	wg.Wait()
	close(results)
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	ids := map[string]int{}
	created := 0
	for r := range results {
		ids[r.EventID]++
		if r.Created {
			created++
		}
	}
	if len(ids) != 8 || created != 8 {
		t.Fatalf("lost or duplicated appends: %d/%d", len(ids), created)
	}
	if before != nativeFingerprint(t, f.root) {
		t.Fatal("append changed native state")
	}
	// A lost successful response does not repeat the described work.
	r, e := f.s.Append("task", inputs[0], "event")
	if e != nil || r.Created {
		t.Fatal("lost response retry not idempotent", e)
	}
}

func TestPublicationFaults(t *testing.T) {
	for _, stage := range []string{"lock", "write", "before_publish", "after_publish", "directory_sync", "unlock"} {
		t.Run(stage, func(t *testing.T) {
			f := fixture(t)
			before := nativeFingerprint(t, f.root)
			f.s.fault = func(point string) error {
				if point == stage {
					return errors.New("deterministic fault")
				}
				return nil
			}
			p := f.file(t, f.input)
			r, e := f.s.Append("task", p, "event")
			if e == nil || r.Created {
				t.Fatal("fault claimed success")
			}
			published := one(stage, "after_publish", "directory_sync", "unlock")
			if published && !strings.Contains(e.Error(), "publication_unknown") {
				t.Fatalf("post-publish error lost unknown state: %v", e)
			}
			f.s.fault = nil
			s := f.show(t)
			count := 0
			for _, e := range s.Events {
				if e.Origin == "contribution" {
					count++
				}
			}
			if (count == 1) != published {
				t.Fatalf("partial publication: %d", count)
			}
			r, e = f.s.Append("task", p, "event")
			if e != nil || r.Created == published {
				t.Fatalf("exact retry: %v %+v", e, r)
			}
			if before != nativeFingerprint(t, f.root) {
				t.Fatal("fault changed native state")
			}
		})
	}
}

func TestCorruptJournalPreservesPrefixAndBlocksAppend(t *testing.T) {
	f := fixture(t)
	first := f.append(t, f.input)
	in := f.input
	in.PublicationKey = "second"
	second := f.append(t, in)
	write(t, second.RecordLocator, []byte(`{"partial":`))
	before := nativeFingerprint(t, f.root)
	s := f.show(t)
	found := false
	for _, e := range s.Events {
		found = found || e.EventID == first.EventID
	}
	if !found || s.Coverage.State != "partial" {
		t.Fatal("valid prefix hidden")
	}
	in.PublicationKey = "third"
	if _, e := f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("damaged log appended")
	}
	if before != nativeFingerprint(t, f.root) {
		t.Fatal("damaged show altered native facts")
	}
}

func TestBindingAndRelationRejections(t *testing.T) {
	f := fixture(t)
	in := f.input
	in.RunBinding = &RunBinding{"trn_" + strings.Repeat("a", 64), "sha256:" + strings.Repeat("b", 64), "pre_" + strings.Repeat("c", 64), "sha256:" + strings.Repeat("d", 64)}
	if _, e := f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("unbound run accepted")
	}
	in = f.input
	in.Relations = []Relation{{"responds_to", "foreign-event"}}
	if _, e := f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("foreign relation accepted")
	}
	in = f.input
	in.Step = &StepRef{"child", "implementation", ptr("future-parent")}
	in.Type = "step_started"
	if _, e := f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("future parent accepted")
	}
	if _, e := os.Lstat(storePath(f.root, "task")); !os.IsNotExist(e) {
		t.Fatal("binding rejection created journal")
	}
}

func TestUncertaintyParentsAndClockConflict(t *testing.T) {
	f := fixture(t)
	start := f.input
	start.Type = "step_started"
	start.Step = &StepRef{"parent", "planning", nil}
	a := f.append(t, start)
	child := start
	child.PublicationKey = "child:start"
	child.Step = &StepRef{"child", "implementation", ptr("parent")}
	child.TimeBasis.Uncertainty = &Uncertainty{"2026-10-03T09:59:00Z", "2026-10-03T10:01:00Z"}
	f.append(t, child)
	child.PublicationKey = "child:finish"
	child.Type = "step_finished"
	child.OccurredAt = ptr("2026-10-03T10:10:00Z")
	child.TimeBasis.Uncertainty = &Uncertainty{"2026-10-03T10:09:00Z", "2026-10-03T10:11:00Z"}
	f.append(t, child)
	start.PublicationKey = "parent:finish"
	start.Type = "step_finished"
	start.OccurredAt = ptr("2026-10-03T10:10:00Z")
	f.append(t, start)
	s := f.show(t)
	for _, p := range s.Steps {
		if p.StepID == "parent" && (p.DurationSeconds == nil || *p.DurationSeconds != 600) {
			t.Fatal("parent duration changed by child")
		}
		if p.StepID == "child" && (p.DurationSeconds != nil || p.State != "time_unknown") {
			t.Fatal("uncertainty became precise duration")
		}
	}
	in := f.input
	in.PublicationKey = "correction"
	in.Relations = []Relation{{"supersedes", a.EventID}}
	f.append(t, in)
	s = f.show(t)
	for _, e := range s.Events {
		if e.EventID == a.EventID && len(e.SupersededBy) != 1 {
			t.Fatal("original not retained as superseded")
		}
	}
}

func TestRetryAfterSourceLossAndCrossCommandConflict(t *testing.T) {
	f := fixture(t)
	r := f.append(t, f.input)
	before, _ := os.ReadFile(r.RecordLocator)
	if e := os.Remove(f.input.Sources[0].Locator); e != nil {
		t.Fatal(e)
	}
	again, e := f.s.Append("task", f.file(t, f.input), "event")
	if e != nil || again.Created || again.RecordedAt != r.RecordedAt {
		t.Fatal("retry depends on mutable source", e)
	}
	if len(again.Warnings) != 1 {
		t.Fatal("retry omitted source warning")
	}
	after, _ := os.ReadFile(r.RecordLocator)
	if string(before) != string(after) {
		t.Fatal("retry rewrote record")
	}
	ob := ObservationInput{Common: f.input.Common, Type: "proposal", Targets: Targets{EventIDs: []string{r.EventID}, StepIDs: []string{}}, Finding: "Evidence missing", Action: ptr("Inspect"), Owner: ptr("planner"), NextSignal: ptr("New evidence")}
	ob.Kind = "ply.workspace.task-journal-observation-input"
	if _, e = f.s.Append("task", f.file(t, ob), "observation"); e == nil {
		t.Fatal("record/observe keys have separate namespaces")
	}
	s := f.show(t)
	if s.Coverage.State != "partial" {
		t.Fatal("missing source not reported")
	}
}

func TestSnapshotDigestAndReferences(t *testing.T) {
	f := fixture(t)
	f.append(t, f.input)
	s := f.show(t)
	s.Events[0].SourceIDs = []string{"missing"}
	s.SnapshotID = snapshotDigest(s)
	if _, e := ReadSnapshot(f.file(t, s), "task", Options{}); e == nil {
		t.Fatal("digest-valid but dangling snapshot accepted")
	}
	s = f.show(t)
	s.Task.Title = "tampered"
	if _, e := ReadSnapshot(f.file(t, s), "task", Options{}); e == nil {
		t.Fatal("tampered snapshot accepted")
	}
}

func TestSourceAndStorageSymlinks(t *testing.T) {
	f := fixture(t)
	link := filepath.Join(f.root, "link")
	if e := os.Symlink(f.input.Sources[0].Locator, link); e != nil {
		t.Fatal(e)
	}
	in := f.input
	in.Sources = []SourceRef{{link, in.Sources[0].SHA256}}
	if _, e := f.s.Append("task", f.file(t, in), "event"); e == nil {
		t.Fatal("symlink source accepted")
	}
	outside := filepath.Join(f.root, "elsewhere")
	if e := os.Mkdir(outside, 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(outside, filepath.Join(f.root, ".ply", "task-process")); e != nil {
		t.Fatal(e)
	}
	if _, e := f.s.Append("task", f.file(t, f.input), "event"); e == nil {
		t.Fatal("symlink storage accepted")
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("escaped journal writes")
	}
}

func TestIntervalsWithoutPointAreNotInvented(t *testing.T) {
	f := fixture(t)
	in := f.input
	in.OccurredAt = nil
	in.TimeBasis.Uncertainty = &Uncertainty{"2026-10-03T10:00:00Z", "2026-10-03T10:01:00Z"}
	r := f.append(t, in)
	s := f.show(t)
	for _, e := range s.Events {
		if e.EventID == r.EventID && (e.OccurredAt != nil || e.TimeBasis.Uncertainty == nil) {
			t.Fatal("interval acquired point time")
		}
	}
	if !strings.Contains(Text(s, "timeline"), "Unknown time") {
		t.Fatal("unknown group missing")
	}
	_ = time.UTC
}
