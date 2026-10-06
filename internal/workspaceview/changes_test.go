package workspaceview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type changesWorkingDirectory struct {
	workspace.FileSystem
	root string
}

func (f changesWorkingDirectory) Getwd() (string, error) { return f.root, nil }

func changesFixture(t *testing.T) (workspace.Dependencies, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = changesWorkingDirectory{d.Files, root}
	if _, err := workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	return d, root
}

func writeChangeFile(t *testing.T, root, relative, content string) string {
	t.Helper()
	p := filepath.Join(root, ".ply", relative)
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestChangesAreContentBoundAndReadOnly(t *testing.T) {
	d, root := changesFixture(t)
	first, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if first.RegistrySHA256 == nil || first.Freshness != "fresh" || len(first.Registers) != 4 {
		t.Fatalf("empty snapshot: %+v", first)
	}
	entries, _ := os.ReadDir(filepath.Join(root, ".ply"))
	if len(entries) != 1 {
		t.Fatalf("read created files: %v", entries)
	}
	p := writeChangeFile(t, root, "work-items.yaml", "content-one")
	second, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *first.RegistrySHA256 == *second.RegistrySHA256 {
		t.Fatal("new register did not change digest")
	}
	if err := os.Chtimes(p, time.Now().Add(time.Hour), time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	writeChangeFile(t, root, "work-items.lock", "lock bytes are not facts")
	writeChangeFile(t, root, "unrelated.json", "not a register")
	third, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *second.RegistrySHA256 != *third.RegistrySHA256 {
		t.Fatal("mtime or unrelated file invalidated content digest")
	}
	writeChangeFile(t, root, "work-items.yaml", "content-two")
	fourth, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *third.RegistrySHA256 == *fourth.RegistrySHA256 {
		t.Fatal("effective content change missing")
	}
}

func TestChangesCoverAllReadModelOwners(t *testing.T) {
	d, root := changesFixture(t)
	previous, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{
		"projects.yaml", "project-metadata.json", "work-items.yaml", "work-item-lifecycle.json",
		"task-process/v1/task-a/000000001.json",
		"task-runs/v1/runs/trn_a/request.json", "task-runs/v1/runs/trn_a/events/000001.json",
		"workflow-runs/requests/wfr_a.json", "workflow-runs/runs/wfr_a/state.json",
	} {
		writeChangeFile(t, root, relative, "one")
		current, err := ReadChanges(d)
		if err != nil {
			t.Fatal(err)
		}
		if current.RegistrySHA256 == nil || *previous.RegistrySHA256 == *current.RegistrySHA256 {
			t.Fatalf("not covered: %s", relative)
		}
		previous = current
	}
	writeChangeFile(t, root, "task-process/v1/task-a/.publish-staging.json", "temporary")
	writeChangeFile(t, root, "workflow-runs/runs/wfr_a/terminal.txt", "not a read-model owner")
	current, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *previous.RegistrySHA256 != *current.RegistrySHA256 {
		t.Fatal("temporary/output files changed registered-facts tag")
	}
}

func TestChangesExposeUnknownWithoutFollowingSymlinks(t *testing.T) {
	d, root := changesFixture(t)
	external := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(external, []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, ".ply", "projects.yaml")); err != nil {
		t.Fatal(err)
	}
	result, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if result.RegistrySHA256 != nil || result.Freshness != "unknown" {
		t.Fatalf("symlink treated as reliable: %+v", result)
	}
	for _, r := range result.Registers {
		if r.ID == "projects" && (r.RegistrySHA256 != nil || r.Freshness != "unknown" || len(r.Reasons) == 0) {
			t.Fatalf("missing source diagnosis: %+v", r)
		}
	}
}

func TestChangesDetectMovementBetweenPasses(t *testing.T) {
	d, root := changesFixture(t)
	writeChangeFile(t, root, "projects.yaml", "before")
	result, err := readChanges(d, func() { writeChangeFile(t, root, "projects.yaml", "after") })
	if err != nil {
		t.Fatal(err)
	}
	if result.Freshness != "stale" || result.RegistrySHA256 != nil {
		t.Fatalf("racing snapshot was cacheable: %+v", result)
	}
}

func TestChangesTrackIncompleteRunsUnexpectedJournalEntriesAndModes(t *testing.T) {
	d, root := changesFixture(t)
	before, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	run := filepath.Join(root, ".ply", "task-runs", "v1", "runs", "trn_"+strings.Repeat("a", 64))
	if err := os.MkdirAll(run, 0700); err != nil {
		t.Fatal(err)
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		t.Fatal(err)
	}
	runs, err := ListRuns(d, s, RunOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(runs.Runs) != 1 || runs.Runs[0].State != "unknown" {
		t.Fatalf("partial run not exposed: %+v", runs)
	}
	current, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *before.RegistrySHA256 == *current.RegistrySHA256 {
		t.Fatal("visible incomplete run did not invalidate cache")
	}
	for _, relative := range []string{"task-process/v1/task-a/unexpected.txt", "task-process/v1/task-a/.unexpected", "task-runs/v1/runs/trn_a/events/extra.lock"} {
		before = current
		writeChangeFile(t, root, relative, "changes read validity")
		current, err = ReadChanges(d)
		if err != nil {
			t.Fatal(err)
		}
		if *before.RegistrySHA256 == *current.RegistrySHA256 {
			t.Fatalf("source entry not fingerprinted: %s", relative)
		}
	}
	p := writeChangeFile(t, root, "task-process/v1/task-a/000000001.json", "{}")
	before, err = ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(p, 0644); err != nil {
		t.Fatal(err)
	}
	current, err = ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if *before.RegistrySHA256 == *current.RegistrySHA256 {
		t.Fatal("loss of journal privacy did not invalidate source validity")
	}
}
