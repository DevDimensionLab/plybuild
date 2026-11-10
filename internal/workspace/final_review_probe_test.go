package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func finalReviewInstallDefinitions(t *testing.T, f workItemJourneyFixture) {
	t.Helper()
	b, err := ReadViewBasis(f.dependencies)
	if err != nil {
		t.Fatal(err)
	}
	defs := workspacedef.Definitions{FormatVersion: 1, Projects: []workspacedef.Project{}}
	for _, p := range b.Projects.Projects {
		pd := workspacedef.Project{ID: string(p.ID), Name: p.Name, Wrapper: p.Wrapper, Repos: []workspacedef.Repo{}, Epics: []workspacedef.Epic{}}
		for _, id := range p.RepoIDs {
			for _, r := range b.Projects.Repos {
				if r.ID == id {
					pd.Repos = append(pd.Repos, workspacedef.Repo{ID: string(r.ID), Path: r.Locator})
				}
			}
		}
		for _, e := range b.Registry.Epics {
			if e.ProjectID == p.ID {
				for _, binding := range e.RepoBindings {
					pd.Epics = append(pd.Epics, workspacedef.Epic{ID: string(e.ID), Title: e.Title, Repo: string(binding.RepoID), Path: binding.Worktree.Locator, Ref: binding.Worktree.Ref})
				}
			}
		}
		defs.Projects = append(defs.Projects, pd)
	}
	raw, err := json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(b.Workspace.Root, DefinitionsFile), raw, 0600); err != nil {
		t.Fatal(err)
	}
}
func TestFinalReviewConfigPreservesRegisteredQueueIdentity(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	before := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if len(before.Queues) != 1 || len(before.Queues[0].Pending) != 1 || before.Queues[0].Pending[0].State != "ready" {
		t.Fatalf("control fixture not ready: %+v", before)
	}
	finalReviewInstallDefinitions(t, f)
	after := readRegisteredQueueFixture(f, registeredQueueBasis(t, f))
	if len(after.Queues) != 1 || len(after.Queues[0].Pending) != 1 || after.Queues[0].Pending[0].State != "ready" {
		t.Fatalf("same definitions lost ready registered queue: %+v", after)
	}
}
func TestFinalReviewFilterUnknownAndANDWithoutCheckout(t *testing.T) {
	root, repo, d := durableFixture(t)
	if _, err := CreateTask(d, TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "feature", ProjectID: "sample", RepoID: "code"}); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, DefinitionsFile)
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var defs workspacedef.Definitions
	if err = json.Unmarshal(raw, &defs); err != nil {
		t.Fatal(err)
	}
	defs.Projects = append(defs.Projects, workspacedef.Project{ID: "second", Name: "Second", Wrapper: ".", Repos: []workspacedef.Repo{{ID: "other", Path: "offline-other"}}, Epics: []workspacedef.Epic{}})
	raw, err = json.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(p, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Rename(repo, repo+"-offline"); err != nil {
		t.Fatal(err)
	}
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	sample, second, missing := ProjectID("sample"), ProjectID("second"), ProjectID("missing")
	code, other, missingRepo := RepoID("code"), RepoID("other"), RepoID("missing")
	missingEpic := EpicID("missing")
	for _, f := range []TaskListFilters{{ProjectID: &missing}, {RepoID: &missingRepo}, {EpicID: &missingEpic}} {
		if _, err := ListTasksWithFilters(d, f); err == nil {
			t.Fatalf("unknown filter accepted: %+v", f)
		}
	}
	for _, c := range []struct {
		f TaskListFilters
		n int
	}{{TaskListFilters{ProjectID: &sample, RepoID: &code}, 1}, {TaskListFilters{ProjectID: &sample, RepoID: &other}, 0}, {TaskListFilters{ProjectID: &second, RepoID: &code}, 0}} {
		out, err := ListTasksWithFilters(d, c.f)
		if err != nil || len(out.Tasks) != c.n {
			t.Fatalf("AND filter %+v count=%d want=%d error=%v", c.f, len(out.Tasks), c.n, err)
		}
	}
}
