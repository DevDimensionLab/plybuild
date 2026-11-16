package agentwork

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// defineProjects writes the workspace definitions of root with one Project
// per id, each with an existing wrapper <root>/<id> holding one repository
// directory, which is never observed through Git.
func defineProjects(t *testing.T, root string, ids ...string) {
	t.Helper()
	projects := ""
	for i, id := range ids {
		if err := os.MkdirAll(filepath.Join(root, id, "main"), 0o755); err != nil {
			t.Fatal(err)
		}
		if i > 0 {
			projects += ","
		}
		projects += `{"id":"` + id + `","name":"` + id + `","wrapper":"` + id + `","repos":[{"id":"` + id + `","path":"` + id + `/main"}]}`
	}
	config := `{"format_version":1,"projects":[` + projects + `]}`
	if err := os.WriteFile(filepath.Join(root, workspace.DefinitionsFile), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A follow-up lives in the work of the Project it was registered for, never
// in .ply. Lookups without a Project find it in any Project, and the same
// draft in two Projects is two follow-ups with distinct IDs.
func TestFollowupsLiveInTheirProject(t *testing.T) {
	s, d, _ := fixture(t)
	defineProjects(t, s.Root, "project-fixture", "other")
	a, err := s.Register("project-fixture", d)
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.Register("other", d)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID {
		t.Fatalf("same ID %s in two Projects", a.ID)
	}
	works, err := workspace.DefinedWork(s.Root)
	if err != nil || len(works) != 2 {
		t.Fatalf("works %v %v", works, err)
	}
	for _, work := range works {
		if entries, err := os.ReadDir(filepath.Join(work.AgentWork(), "revisions")); err != nil || len(entries) != 1 {
			t.Fatalf("store of %s: %v %v", work.Project, entries, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(s.Root, ".ply", "agent-work")); !os.IsNotExist(err) {
		t.Fatalf(".ply/agent-work exists: %v", err)
	}
	if got, err := s.Show(b.ID); err != nil || got.ID != b.ID {
		t.Fatalf("show %v %v", got.ID, err)
	}
	if list, err := s.List(false); err != nil || len(list.Items) != 2 {
		t.Fatalf("list %v %v", list.Items, err)
	}
	if _, err := s.Mutate(b.ID, "cancel", mutation(b, "cancel-other")); err != nil {
		t.Fatal(err)
	}
	if got, err := s.Show(a.ID); err != nil || got.State != "waiting" {
		t.Fatalf("project-fixture follow-up changed: %v %v", got.State, err)
	}
	if _, err := s.Register("missing", d); err == nil || err.Error() != "followup_project_not_found" {
		t.Fatalf("unknown Project: %v", err)
	}
}

// only is the only Project of the workspace at s.Root. Tests of one Project
// reach its ID and store only through it.
func (s *Service) only() workspace.ProjectWork {
	works, err := workspace.DefinedWork(s.Root)
	if err != nil || len(works) != 1 {
		panic(fmt.Sprintf("workspace has no single Project: %v %v", works, err))
	}
	return works[0]
}
func (s *Service) onlyProject() string                  { return string(s.only().Project) }
func (s *Service) dir() string                          { return s.only().AgentWork() }
func (s *Service) read() (database, error)              { return read(s.dir()) }
func (s *Service) save(db *database) error              { return save(s.dir(), db) }
func (s *Service) locked(f func(*database) error) error { return locked(s.dir(), f) }
