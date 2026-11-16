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

// defineChangeProject defines one Project, whose registry lives under its
// wrapper; resolving it observes no Git.
func defineChangeProject(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "project"), 0o755); err != nil {
		t.Fatal(err)
	}
	definitions := `{"format_version": 1, "projects": [{"id": "project", "name": "Project", "wrapper": "project", "repos": [{"id": "repo", "path": "project/repo"}]}]}`
	if err := os.WriteFile(filepath.Join(root, workspace.DefinitionsFile), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
}

const emptyRegistry = "format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n"

func writeChangeFile(t *testing.T, root, relative, content string) string {
	t.Helper()
	p := filepath.Join(root, ".ply", relative)
	// Work items, lifecycle and journals live in the only Project's work.
	if strings.HasPrefix(relative, "work-item") || strings.HasPrefix(relative, "task-process/") {
		works, err := workspace.DefinedWork(root)
		if err != nil {
			t.Fatal(err)
		}
		p = filepath.Join(works[0].Dir, relative)
	}
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
	defineChangeProject(t, root)
	first, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	if first.RegistrySHA256 == nil || first.Freshness != "fresh" || len(first.Registers) != 3 {
		t.Fatalf("empty snapshot: %+v", first)
	}
	// init leaves projects.lock beside the marker; reads add nothing.
	entries, _ := os.ReadDir(filepath.Join(root, ".ply"))
	if len(entries) != 2 {
		t.Fatalf("read created files: %v", entries)
	}
	p := writeChangeFile(t, root, "work-items.yaml", emptyRegistry+"# content-one\n")
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
	writeChangeFile(t, root, "work-items.yaml", emptyRegistry+"# content-two\n")
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
	defineChangeProject(t, root)
	previous, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	// The definitions change only in formatting, so reads still accept them.
	definitions := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(definitions)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(definitions, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if current, err := ReadChanges(d); err != nil || current.RegistrySHA256 == nil || *previous.RegistrySHA256 == *current.RegistrySHA256 {
		t.Fatalf("not covered: ply-workspace.json %v", err)
	} else {
		previous = current
	}
	for _, relative := range []string{
		"project-metadata.json", "work-items.yaml", "work-item-lifecycle.json",
		"task-process/v1/task-a/000000001.json",
	} {
		content := "one"
		if relative == "work-items.yaml" {
			// Reads validate the registry; only its bytes matter here.
			content = emptyRegistry + "# one\n"
		}
		writeChangeFile(t, root, relative, content)
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
	// Old Run stores are no longer a read-model owner.
	writeChangeFile(t, root, "task-runs/v1/runs/trn_a/request.json", "not a read-model owner")
	writeChangeFile(t, root, "workflow-runs/runs/wfr_a/state.json", "not a read-model owner")
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
	if err := os.Symlink(external, filepath.Join(root, ".ply", "project-metadata.json")); err != nil {
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
	writeChangeFile(t, root, "project-metadata.json", "before")
	result, err := readChanges(d, func() { writeChangeFile(t, root, "project-metadata.json", "after") })
	if err != nil {
		t.Fatal(err)
	}
	if result.Freshness != "stale" || result.RegistrySHA256 != nil {
		t.Fatalf("racing snapshot was cacheable: %+v", result)
	}
}

func TestChangesTrackUnexpectedJournalEntriesAndModes(t *testing.T) {
	d, root := changesFixture(t)
	defineChangeProject(t, root)
	current, err := ReadChanges(d)
	if err != nil {
		t.Fatal(err)
	}
	var before ChangesResult
	for _, relative := range []string{"task-process/v1/task-a/unexpected.txt", "task-process/v1/task-a/.unexpected"} {
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
