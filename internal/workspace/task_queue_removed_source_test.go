package workspace

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func removedSourceFixture(t *testing.T) (workItemJourneyFixture, TaskPreparationReadback, string) {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	p := queuePrepareFixture(t, f, true)
	runLocalGit(t, p.Plan.WorktreePath, "commit", "--allow-empty", "-m", "legitimate later work")
	oid := runLocalGit(t, p.Plan.WorktreePath, "rev-parse", "HEAD")
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "remove", p.Plan.WorktreePath)
	return f, p, oid
}

func TestTaskQueueReleaseRemovedSource(t *testing.T) {
	f, p, oid := removedSourceFixture(t)
	before := mustQueueRegistry(t, f)
	gitBefore := queueStateBytes(t, f)
	for _, kind := range []string{"release", "advance"} {
		if _, e := CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", kind); e == nil {
			t.Fatalf("default %s accepted missing source", kind)
		}
	}
	if queueStateBytes(t, f) != gitBefore {
		t.Fatal("default rejection changed state")
	}
	q, e := ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", oid)
	if e != nil || q.Current != nil || q.Revision != 3 {
		t.Fatalf("removed source should release with preserved work: %#v %v", q, e)
	}
	after := mustQueueRegistry(t, f)
	if len(after.TaskQueueEvents) != len(before.TaskQueueEvents)+1 {
		t.Fatal("release did not append exactly one event")
	}
	event := after.TaskQueueEvents[len(after.TaskQueueEvents)-1]
	if event.Kind != "release" || event.Request["source_state"] != queueSourceRemoved || event.Request["expected_source_oid"] != oid {
		t.Fatalf("release lost removed-source evidence: %#v", event)
	}
	after.TaskQueueEvents = after.TaskQueueEvents[:len(after.TaskQueueEvents)-1]
	after.RawSHA256 = before.RawSHA256
	if !reflect.DeepEqual(before, after) {
		t.Fatal("release rewrote registry history beyond the new event")
	}
	stable := queueStateBytes(t, f)
	if _, e = ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", oid); e != nil || stable != queueStateBytes(t, f) {
		t.Fatalf("identical retry changed state: %v", e)
	}
	for _, change := range []string{"oid", "reason", "mode"} {
		switch change {
		case "oid":
			_, e = ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", f.oid)
		case "reason":
			_, e = ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "different reason", oid)
		case "mode":
			_, e = CloseTaskQueue(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", "release")
		}
		if e == nil || stable != queueStateBytes(t, f) {
			t.Fatalf("changed %s reused a prior request: %v", change, e)
		}
	}
	// An acknowledged retry reads the prior decision even if its old evidence has
	// since changed; it must not silently turn into a fresh release.
	runLocalGit(t, filepath.Join(f.wrapper, "main"), "update-ref", "refs/heads/"+p.Plan.Branch, f.oid)
	if e = os.Mkdir(p.Plan.WorktreePath, 0700); e != nil {
		t.Fatal(e)
	}
	stable = queueStateBytes(t, f)
	if _, e = ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "source removed, preserve work", oid); e != nil || stable != queueStateBytes(t, f) {
		t.Fatalf("historical retry revalidated old evidence: %v", e)
	}
}

type removedSourceFiles struct {
	FileSystem
	path string
}

func (f removedSourceFiles) Lstat(path string) (fs.FileInfo, error) {
	if path == f.path {
		return nil, fs.ErrPermission
	}
	return f.FileSystem.Lstat(path)
}

type removedSourceGit struct {
	WorkItemGit
	mode string
}

func (g removedSourceGit) ObserveRepo(repo RepoRecord) (GitRepoObservation, error) {
	o, e := g.WorkItemGit.ObserveRepo(repo)
	if g.mode == "repo" {
		o.Locator += "-different"
	} else if g.mode == "common-dir" {
		o.GitCommonDir += "-different"
	}
	return o, e
}
func (g removedSourceGit) ObserveRef(repo RepoRecord, ref string) (GitRefObservation, error) {
	o, e := g.WorkItemGit.ObserveRef(repo, ref)
	if g.mode == "ref" {
		o.Ref += "-different"
	} else if g.mode == "ref-unknown" {
		e = errors.New("injected unreadable ref")
	}
	return o, e
}
func (g removedSourceGit) ListWorktrees(repo RepoRecord) (GitWorktreeInventory, error) {
	if g.mode == "inventory-unknown" {
		return GitWorktreeInventory{}, errors.New("injected unreadable inventory")
	}
	return g.WorkItemGit.ListWorktrees(repo)
}

func TestTaskQueueReleaseRemovedSourceRejectsChangedEvidence(t *testing.T) {
	f, p, oid := removedSourceFixture(t)
	main := filepath.Join(f.wrapper, "main")
	ref := "refs/heads/" + p.Plan.Branch
	for _, mode := range []string{"wrong-oid", "malformed-oid", "wrong-preparation", "wrong-revision", "empty-reason", "wrong-target", "directory", "file", "symlink", "dangling-symlink", "unreadable", "missing-ref", "other-checkout", "stale-inventory", "repo", "common-dir", "ref", "ref-unknown", "inventory-unknown"} {
		t.Run(mode, func(t *testing.T) {
			d := f.dependencies
			id, expected, reason, expectedOID, target := p.Preparation.ID, 2, "preserve history", oid, queueTargetFixture()
			cleanup := func() {}
			switch mode {
			case "wrong-oid":
				expectedOID = f.oid
			case "malformed-oid":
				expectedOID = oid[:12]
			case "wrong-preparation":
				id = "pre_" + strings.Repeat("f", 64)
			case "wrong-revision":
				expected = 1
			case "empty-reason":
				reason = ""
			case "wrong-target":
				target.RepoID = "other"
			case "directory", "file", "symlink", "dangling-symlink":
				var e error
				if mode == "directory" {
					e = os.Mkdir(p.Plan.WorktreePath, 0700)
				} else if mode == "file" {
					e = os.WriteFile(p.Plan.WorktreePath, []byte("reused"), 0600)
				} else {
					target := main
					if mode == "dangling-symlink" {
						target += "-absent"
					}
					e = os.Symlink(target, p.Plan.WorktreePath)
				}
				if e != nil {
					t.Fatal(e)
				}
				cleanup = func() {
					if e := os.Remove(p.Plan.WorktreePath); e != nil {
						t.Fatal(e)
					}
				}
			case "unreadable":
				d.Files = removedSourceFiles{d.Files, p.Plan.WorktreePath}
			case "missing-ref":
				runLocalGit(t, main, "update-ref", "-d", ref)
				cleanup = func() { runLocalGit(t, main, "update-ref", ref, oid) }
			case "other-checkout":
				other := filepath.Join(f.wrapper, "other-checkout")
				runLocalGit(t, main, "worktree", "add", other, p.Plan.Branch)
				cleanup = func() { runLocalGit(t, main, "worktree", "remove", other) }
			case "stale-inventory":
				runLocalGit(t, main, "worktree", "add", "--detach", p.Plan.WorktreePath, oid)
				if e := os.RemoveAll(p.Plan.WorktreePath); e != nil {
					t.Fatal(e)
				}
				cleanup = func() { runLocalGit(t, main, "worktree", "prune") }
			default:
				d.WorkGit = removedSourceGit{d.WorkGit, mode}
			}
			defer cleanup()
			before := queueStateBytes(t, f)
			if _, e := ReleaseTaskQueueRemovedSource(d, target, id, expected, reason, expectedOID); e == nil {
				t.Fatal("changed evidence was accepted")
			}
			if before != queueStateBytes(t, f) {
				t.Fatal("rejected release changed state")
			}
		})
	}
}

func TestTaskQueueReleaseRemovedSourceRequiresKnownCreation(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	sel := queueSelect(t, f, "task", "a", nil)
	queueSetFixture(t, f, "initial", []QueueEntry{{"task", &sel}})
	preview := queuePrepareFixture(t, f, false)
	original := f.dependencies.WorkGit
	f.dependencies.WorkGit = queueFaultGit{original, original.(TaskContentGit), "branch"}
	if _, e := PrepareTask(f.dependencies, TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "next"}, Apply: true, Confirmation: *preview.Confirmation}); e == nil {
		t.Fatal("injected partial creation did not fail")
	}
	f.dependencies.WorkGit = original
	r := mustQueueRegistry(t, f)
	if r.TaskPreparations[0].Outcome != nil {
		t.Fatal("fixture unexpectedly has known creation")
	}
	before := queueStateBytes(t, f)
	if _, e := ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), r.TaskPreparations[0].ID, 2, "cannot abandon unknown creation", f.oid); e == nil || before != queueStateBytes(t, f) {
		t.Fatalf("unknown creation was released or changed: %v", e)
	}
}

func TestTaskQueueRemovedSourceEventStrictness(t *testing.T) {
	f, p, oid := removedSourceFixture(t)
	if _, e := ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, 2, "preserve history", oid); e != nil {
		t.Fatal(e)
	}
	r := mustQueueRegistry(t, f)
	for _, mode := range []string{"unknown-field", "missing-oid", "bad-oid", "missing-state", "wrong-state", "advance", "no-outcome"} {
		t.Run(mode, func(t *testing.T) {
			copy := r
			copy.TaskQueueEvents = append([]TaskQueueEvent{}, r.TaskQueueEvents...)
			event := &copy.TaskQueueEvents[len(copy.TaskQueueEvents)-1]
			event.Request = queueRequest(event.Request)
			switch mode {
			case "unknown-field":
				event.Request["completed"] = true
			case "missing-oid":
				delete(event.Request, "expected_source_oid")
			case "bad-oid":
				event.Request["expected_source_oid"] = "HEAD"
			case "missing-state":
				delete(event.Request, "source_state")
			case "wrong-state":
				event.Request["source_state"] = "unknown"
			case "advance":
				event.Kind = "advance"
			case "no-outcome":
				copy.TaskPreparations = append([]TaskPreparation{}, r.TaskPreparations...)
				copy.TaskPreparations[0].Outcome = nil
			}
			event.RequestSHA256 = queueDigest(event.Request)
			if e := validateQueueRegistry(copy); e == nil {
				t.Fatal("invalid removed-source event accepted")
			}
		})
	}
}

func TestTaskQueueRemovedSourceConcurrentReleaseAndStart(t *testing.T) {
	f, p, oid := removedSourceFixture(t)
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "next", Title: "Next", Description: "Next work", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	next := queueSelect(t, f, "next", "b", nil)
	q := queueSetFixture(t, f, "pending", []QueueEntry{{"next", &next}})
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	start := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, e := ReleaseTaskQueueRemovedSource(f.dependencies, queueTargetFixture(), p.Preparation.ID, q.Revision, "preserve history", oid)
			errs <- e
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		id := TaskID("next")
		in := TaskPrepareInput{Target: queueTargetFixture(), Selector: QueueSelector{Kind: "task", TaskID: &id}}
		preview, e := PrepareTask(f.dependencies, in)
		if e != nil {
			if !strings.Contains(e.Error(), "task_queue_current_conflict") {
				t.Errorf("unexpected start rejection: %v", e)
			}
			return
		}
		if preview.Plan == nil || preview.Plan.QueueRevision != q.Revision+1 {
			t.Errorf("start reused release revision: %#v", preview.Plan)
			return
		}
		in.Apply, in.Confirmation = true, *preview.Confirmation
		if _, e = PrepareTask(f.dependencies, in); e != nil {
			t.Errorf("start after release: %v", e)
		}
	}()
	close(start)
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	r := mustQueueRegistry(t, f)
	releases := 0
	for _, ev := range r.TaskQueueEvents {
		if ev.Kind == "release" {
			releases++
			if ev.Revision != q.Revision+1 {
				t.Fatal("release used another revision")
			}
		}
	}
	if releases != 1 {
		t.Fatalf("concurrent retries published %d releases", releases)
	}
	if len(r.TaskPreparations) == 1 {
		queuePrepareFixture(t, f, true)
		r = mustQueueRegistry(t, f)
	}
	last := r.TaskQueueEvents[len(r.TaskQueueEvents)-1]
	if len(r.TaskPreparations) != 2 || last.Kind != "reserve" || last.Revision != q.Revision+2 {
		t.Fatalf("new start did not use the next revision: %#v", last)
	}
}
