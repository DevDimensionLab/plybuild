package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func durableFixture(t *testing.T) (string, string, Dependencies) {
	t.Helper()
	root := physicalPath(t, t.TempDir())
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "init", "-b", "main")
	runLocalGit(t, repo, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, repo, "config", "user.name", "Ply tests")
	if err := os.WriteFile(filepath.Join(repo, "README"), []byte("fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "add", "README")
	runLocalGit(t, repo, "commit", "-m", "fixture")
	config := `{"format_version":1,"projects":[{"id":"sample","name":"Sample","wrapper":".","repos":[{"id":"code","path":"repo"}],"epics":[{"id":"feature","title":"Feature","repo":"code","path":"repo","ref":"refs/heads/main"}]}]}`
	if err := os.WriteFile(filepath.Join(root, "ply-workspace.json"), []byte(config), 0644); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: repo}
	return root, repo, d
}

func TestDurableConfigOrdinaryUseAfterEmptyRuntime(t *testing.T) {
	root, repo, d := durableFixture(t)
	var previous WorktreeID
	for i := 0; i < 2; i++ {
		projects, err := ListProjects(d)
		if err != nil || len(projects.Projects) != 1 {
			t.Fatalf("ordinary project list: %#v, %v", projects, err)
		}
		epics, err := ListEpics(d)
		if err != nil || len(epics.Epics) != 1 {
			t.Fatalf("ordinary epic list: %#v, %v", epics, err)
		}
		binding := epics.Epics[0].RepoBindings[0].Worktree
		if binding.OID != runLocalGit(t, repo, "rev-parse", "HEAD") {
			t.Fatal("bootstrap did not observe current Git")
		}
		if binding.ID == previous {
			t.Fatal("reused runtime identity after reset")
		}
		previous = binding.ID
		_, err = CreateTask(d, TaskCreateInput{TaskID: "task", Title: "Task", Description: "Fixture", ParentEpicID: "feature", ProjectID: "sample", RepoID: "code"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = os.Stat(filepath.Join(root, MarkerDirectory, ProjectsFile)); !os.IsNotExist(err) {
			t.Fatalf("duplicated project definitions: %v", err)
		}
		if err = os.RemoveAll(filepath.Join(root, MarkerDirectory)); err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			runLocalGit(t, repo, "commit", "--allow-empty", "-m", "advance after reset")
		}
	}
}

func TestObserveContainingRootDoesNotAdoptDefinedEpics(t *testing.T) {
	root, _, d := durableFixture(t)
	if _, err := Init(Dependencies{Files: projectCwdFileSystem{FileSystem: d.Files, cwd: root}}); err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, root)
	observed, err := ObserveContainingRoot(d)
	if err != nil || observed.Root != root {
		t.Fatalf("observed %#v, %v", observed, err)
	}
	if after := snapshotTree(t, root); !reflect.DeepEqual(after, before) {
		t.Fatal("root observation changed the workspace")
	}
	// The same workspace is adopted by ordinary discovery, so the check above is not vacuous.
	if _, err := ObserveContaining(d); err != nil {
		t.Fatal(err)
	}
	if after := snapshotTree(t, root); reflect.DeepEqual(after, before) {
		t.Fatal("ordinary discovery did not adopt the defined Epic")
	}
}

func TestDurableConfigRejectsUndeclaredEpic(t *testing.T) {
	_, repo, d := durableFixture(t)
	_, err := AdoptEpic(d, EpicAdoptInput{EpicID: "extra", Title: "Extra", ProjectID: "sample", RepoID: "code", Worktree: repo, Ref: "refs/heads/main", ExpectedOID: runLocalGit(t, repo, "rev-parse", "HEAD")})
	if err == nil {
		t.Fatal("native adoption bypassed authoritative definitions")
	}
}

func TestDurableConfigGuards(t *testing.T) {
	for _, test := range []string{"dirty", "ref", "repo", "duplicate-git", "duplicate-json", "unknown-field", "symlink"} {
		t.Run(test, func(t *testing.T) {
			root, repo, d := durableFixture(t)
			path := filepath.Join(root, "ply-workspace.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			switch test {
			case "dirty":
				err = os.WriteFile(filepath.Join(repo, "dirty"), []byte("dirty"), 0644)
			case "ref":
				raw = []byte(strings.Replace(string(raw), "refs/heads/main", "refs/heads/other", 1))
			case "repo":
				other := filepath.Join(root, "other")
				if err = os.Mkdir(other, 0755); err != nil {
					t.Fatal(err)
				}
				runLocalGit(t, other, "init", "-b", "main")
				runLocalGit(t, other, "config", "user.name", "Ply tests")
				runLocalGit(t, other, "config", "user.email", "tests@example.invalid")
				runLocalGit(t, other, "commit", "--allow-empty", "-m", "other")
				// A clean worktree on the right ref in a separate repository.
				raw = []byte(strings.Replace(string(raw), `"path":"repo","ref"`, `"path":"other","ref"`, 1))
			case "duplicate-git":
				raw = []byte(strings.Replace(string(raw), `{"id":"code","path":"repo"}`, `{"id":"code","path":"repo"},{"id":"alias","path":"repo"}`, 1))
			case "duplicate-json":
				raw = []byte(strings.Replace(string(raw), `"format_version":1`, `"format_version":1,"format_version":1`, 1))
			case "unknown-field":
				raw = []byte(strings.Replace(string(raw), `"format_version":1`, `"format_version":1,"old_oid":"stale"`, 1))
			case "symlink":
				if err = os.Rename(path, path+".source"); err != nil {
					t.Fatal(err)
				}
				err = os.Symlink(path+".source", path)
			}
			if err != nil {
				t.Fatal(err)
			}
			if test != "symlink" {
				if err = os.WriteFile(path, raw, 0644); err != nil {
					t.Fatal(err)
				}
			}
			if _, err = ListEpics(d); err == nil {
				t.Fatal("invalid definition/Git state accepted")
			}
			r, err := d.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Epics) != 0 {
				t.Fatal("failed bootstrap persisted Epic")
			}
		})
	}
}

func TestDurableConfigConcurrentBootstrap(t *testing.T) {
	root, _, d := durableFixture(t)
	const n = 6
	results := make(chan EpicListResult, n)
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() { r, e := ListEpics(d); results <- r; errs <- e }()
	}
	var id WorktreeID
	for i := 0; i < n; i++ {
		r := <-results
		err := <-errs
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Epics) != 1 {
			t.Fatalf("epics %#v", r)
		}
		got := r.Epics[0].RepoBindings[0].Worktree.ID
		if id != "" && got != id {
			t.Fatal("concurrent bootstrap changed identity")
		}
		id = got
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil || len(r.Epics) != 1 {
		t.Fatalf("registry %#v, %v", r, err)
	}
}

type configChangingGit struct {
	WorkItemGit
	path string
}

func (g configChangingGit) ObserveRef(repo RepoRecord, ref string) (GitRefObservation, error) {
	result, err := g.WorkItemGit.ObserveRef(repo, ref)
	raw, e := os.ReadFile(g.path)
	if e != nil {
		return result, e
	}
	if e = os.WriteFile(g.path, append(raw, ' '), 0644); e != nil {
		return result, e
	}
	return result, err
}
func TestDurableConfigChangeDuringGitObservationDoesNotPublish(t *testing.T) {
	root, _, d := durableFixture(t)
	d.WorkGit = configChangingGit{WorkItemGit: d.WorkGit, path: filepath.Join(root, "ply-workspace.json")}
	if _, err := ListEpics(d); err == nil {
		t.Fatal("concurrent config change accepted")
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Epics) != 0 {
		t.Fatal("published binding from stale config")
	}
}

func TestDurableConfigExistingRuntimeRejectsChangedGitIdentity(t *testing.T) {
	root, _, d := durableFixture(t)
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(root, "other")
	if err := os.Mkdir(other, 0755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, other, "init", "-b", "main")
	runLocalGit(t, other, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, other, "config", "user.name", "Ply tests")
	runLocalGit(t, other, "commit", "--allow-empty", "-m", "unrelated")
	// Same configured repo ID, but a different Git common directory.
	path := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.Replace(string(raw), `"path":"repo"`, `"path":"other"`, 1))
	if err = os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err = ListEpics(d); err == nil {
		t.Fatal("runtime accepted a different repository under the same repo ID")
	}
}

func TestDurableConfigAuthorityAndStableReadback(t *testing.T) {
	root, _, d := durableFixture(t)
	first, err := ListEpics(d)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	again, err := ListEpics(d)
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if first.Epics[0].RepoBindings[0].Worktree.ID != again.Epics[0].RepoBindings[0].Worktree.ID || string(before) != string(after) {
		t.Fatal("ordinary read rewrote existing runtime binding")
	}
	if err = os.WriteFile(filepath.Join(root, MarkerDirectory, ProjectsFile), []byte("invalid old cache"), 0644); err != nil {
		t.Fatal(err)
	}
	if got, err := ListProjects(d); err != nil || len(got.Projects) != 1 || got.Projects[0].ID != "sample" {
		t.Fatalf("config not authoritative: %#v, %v", got, err)
	}
	input := ProjectAddInput{ProjectID: "extra", Name: "Extra", Wrapper: root, Repos: []RepoInput{{RepoID: "code", Path: filepath.Join(root, "repo")}}}
	if _, err = AddProject(d, input); err == nil {
		t.Fatal("Project add overrode config")
	}
	basis, err := ReadViewBasis(d)
	if err != nil {
		t.Fatal(err)
	}
	if basis.Sources[filepath.Join(root, "ply-workspace.json")] == "" {
		t.Fatal("view omitted authoritative source")
	}
	path := filepath.Join(root, "ply-workspace.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, append(raw, ' '), 0644); err != nil {
		t.Fatal(err)
	}
	if state, _ := CheckReadViewBasis(d, basis); state == "fresh" {
		t.Fatal("config change did not stale view")
	}
}

type definitionRegistrationReader struct{ WorkItemStore }

func (s definitionRegistrationReader) Snapshot(string) (WorkItemRegistry, error) {
	return WorkItemRegistry{}, fmt.Errorf("full content read is unavailable")
}
func TestDurableConfigViewsKeepRegistrationOnlyRead(t *testing.T) {
	_, _, d := durableFixture(t)
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	d.WorkItems = definitionRegistrationReader{d.WorkItems}
	if _, err := ReadViewBasis(d); err != nil {
		t.Fatalf("view needed unrelated full content: %v", err)
	}
}
