package workspace

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

type projectCwdFileSystem struct {
	FileSystem
	cwd string
}

func (filesystem projectCwdFileSystem) Getwd() (string, error) { return filesystem.cwd, nil }

type fakeRepoObserver struct {
	observations map[string]RepoObservation
	errors       map[string]error
	calls        []string
}

func (observer *fakeRepoObserver) Remote(string) (string, error) { return "", nil }

func (observer *fakeRepoObserver) Observe(path string) (RepoObservation, error) {
	observer.calls = append(observer.calls, path)
	if err := observer.errors[path]; err != nil {
		return RepoObservation{}, err
	}
	observation, found := observer.observations[path]
	if !found {
		return RepoObservation{}, fmt.Errorf("unexpected repository %s", path)
	}
	return observation, nil
}

func TestParseProjectAddInputValidatesIdentifiersNamesAndRepositories(t *testing.T) {
	input, err := ParseProjectAddInput("trip", "Trip", "../trip", []string{
		"service=../trip-service",
		"api=../trip-openapi=mirror",
	})
	if err != nil {
		t.Fatal(err)
	}
	if input.ProjectID != ProjectID("trip") || input.Name != "Trip" || len(input.Repos) != 2 {
		t.Fatalf("input = %#v", input)
	}
	if input.Repos[1] != (RepoInput{RepoID: RepoID("api"), Path: "../trip-openapi=mirror"}) {
		t.Fatalf("second repository = %#v", input.Repos[1])
	}

	valid63 := "a" + strings.Repeat("1", 62)
	for _, value := range []string{"a", "repo-1", valid63} {
		if _, err := ParseProjectID(value); err != nil {
			t.Fatalf("valid ID %q: %v", value, err)
		}
	}
	for _, value := range []string{"", "A", "1repo", "repo-", "repo--one", "repo_one", "répo", "a" + strings.Repeat("1", 63)} {
		if _, err := ParseProjectID(value); err == nil || !strings.HasPrefix(err.Error(), string(ErrorProjectInvalidArguments)+":") {
			t.Fatalf("invalid ID %q error = %v", value, err)
		}
	}

	for _, test := range []struct {
		name      string
		id        string
		label     string
		noWrapper bool
		repos     []string
		class     ErrorClass
	}{
		{name: "invalid project ID", id: "Trip", label: "Trip", repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "empty name", id: "trip", repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "trimmed name", id: "trip", label: " Trip", repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "multiline name", id: "trip", label: "Trip\nApp", repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "control name", id: "trip", label: "Trip\u007f", repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "long name", id: "trip", label: strings.Repeat("ø", 129), repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "missing wrapper", id: "trip", label: "Trip", noWrapper: true, repos: []string{"repo=path"}, class: ErrorProjectInvalidArguments},
		{name: "missing separator", id: "trip", label: "Trip", repos: []string{"repo"}, class: ErrorProjectInvalidArguments},
		{name: "empty path", id: "trip", label: "Trip", repos: []string{"repo="}, class: ErrorProjectInvalidArguments},
		{name: "duplicate repository ID", id: "trip", label: "Trip", repos: []string{"repo=one", "repo=two"}, class: ErrorProjectDuplicateRepo},
		{name: "empty repositories", id: "trip", label: "Trip", class: ErrorProjectInvalidArguments},
	} {
		t.Run(test.name, func(t *testing.T) {
			wrapper := "wrapper"
			if test.noWrapper {
				wrapper = ""
			}
			_, err := ParseProjectAddInput(test.id, test.label, wrapper, test.repos)
			assertProjectErrorClass(t, err, test.class)
		})
	}
}

func TestLocateContainingWorkspaceFromRootChildAndSymlink(t *testing.T) {
	root := createProjectWorkspace(t)
	child := filepath.Join(root, "one", "two")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(t.TempDir(), "workspace-link")
	if err := os.Symlink(root, alias); err != nil {
		t.Fatal(err)
	}
	files := SystemDependencies().Files
	for _, cwd := range []string{root, child, filepath.Join(alias, "one", "two")} {
		physical, err := canonicalDirectory(files, "", cwd, "working directory")
		if err != nil {
			t.Fatal(err)
		}
		got, err := locateContainingWorkspace(files, physical)
		if err != nil || got != root {
			t.Fatalf("locate(%s) = %q, %v; want %q", cwd, got, err, root)
		}
	}
}

func TestLocateContainingWorkspaceStopsAtFirstInvalidCandidate(t *testing.T) {
	outer := createProjectWorkspace(t)
	child := filepath.Join(outer, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	files := SystemDependencies().Files

	t.Run("ordinary marker directory without marker is skipped", func(t *testing.T) {
		cwd := filepath.Join(child, "missing-marker", "leaf")
		if err := os.MkdirAll(filepath.Join(child, "missing-marker", MarkerDirectory), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(cwd, 0o755); err != nil {
			t.Fatal(err)
		}
		got, err := locateContainingWorkspace(files, cwd)
		if err != nil || got != outer {
			t.Fatalf("root = %q, error = %v", got, err)
		}
	})

	for _, test := range []struct {
		name  string
		setup func(string)
	}{
		{name: "reserved file", setup: func(cwd string) { mustWriteFile(t, filepath.Join(cwd, MarkerDirectory), "reserved") }},
		{name: "reserved symlink", setup: func(cwd string) {
			if err := os.Symlink(outer, filepath.Join(cwd, MarkerDirectory)); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "marker directory", setup: func(cwd string) {
			if err := os.MkdirAll(filepath.Join(cwd, MarkerDirectory, MarkerFile), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "invalid marker", setup: func(cwd string) {
			if err := os.Mkdir(filepath.Join(cwd, MarkerDirectory), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteFile(t, filepath.Join(cwd, MarkerDirectory, MarkerFile), "format_version: 2\nroot: /wrong\n")
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			cwd := filepath.Join(child, strings.ReplaceAll(test.name, " ", "-"))
			if err := os.Mkdir(cwd, 0o755); err != nil {
				t.Fatal(err)
			}
			test.setup(cwd)
			_, err := locateContainingWorkspace(files, cwd)
			assertProjectErrorClass(t, err, ErrorProjectWorkspace)
		})
	}

	missing := t.TempDir()
	_, err := locateContainingWorkspace(files, missing)
	assertProjectErrorClass(t, err, ErrorWorkspaceNotFound)
	if !strings.Contains(err.Error(), "no Ply workspace contains "+missing) {
		t.Fatalf("missing error = %v", err)
	}
}

func TestAddProjectRejectsDuplicateObservedCommonDirectoryBeforeStore(t *testing.T) {
	root := createProjectWorkspace(t)
	repoA := filepath.Join(root, "repo-a")
	repoB := filepath.Join(root, "repo-b")
	for _, path := range []string{repoA, repoB} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repoA = physicalPath(t, repoA)
	repoB = physicalPath(t, repoB)
	common := filepath.Join(root, "common")
	observer := &fakeRepoObserver{observations: map[string]RepoObservation{
		repoA: {Locator: repoA, GitCommonDir: common},
		repoB: {Locator: repoB, GitCommonDir: common},
	}, errors: map[string]error{}}
	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: root}
	dependencies.Repos = observer
	_, err := AddProject(dependencies, ProjectAddInput{ProjectID: "trip", Name: "Trip", Wrapper: root, Repos: []RepoInput{{RepoID: "a", Path: repoA}, {RepoID: "b", Path: repoB}}})
	assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo)
}

func TestSystemRepoObserverUsesOnlyLocalIdentityCommands(t *testing.T) {
	repository := filepath.Join(t.TempDir(), "repository")
	if err := os.Mkdir(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	repository = physicalPath(t, repository)
	runGitForTest(t, repository, "init", "-q")
	mustWriteFile(t, filepath.Join(repository, "dirty.txt"), "dirty\n")
	observer := newSystemRepoObserver(SystemDependencies().Files)
	observation, err := observer.Observe(repository)
	if err != nil {
		t.Fatal(err)
	}
	if observation.Locator != repository || observation.GitCommonDir != filepath.Join(repository, ".git") {
		t.Fatalf("observation = %#v", observation)
	}

	subdirectory := filepath.Join(repository, "subdirectory")
	if err := os.Mkdir(subdirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err = observer.Observe(subdirectory)
	assertProjectErrorClass(t, err, ErrorProjectRepoInvalid)

	nonGit := t.TempDir()
	_, err = observer.Observe(nonGit)
	assertProjectErrorClass(t, err, ErrorProjectRepoInvalid)

	bare := filepath.Join(t.TempDir(), "bare.git")
	runGitForTest(t, "", "init", "--bare", "-q", bare)
	_, err = observer.Observe(bare)
	assertProjectErrorClass(t, err, ErrorProjectRepoInvalid)
}

func TestSystemRepoObserverCanonicalizesLinkedWorktreeCommonDirectory(t *testing.T) {
	parent := physicalPath(t, t.TempDir())
	mainRepository := filepath.Join(parent, "main")
	linkedRepository := filepath.Join(parent, "linked")
	if err := os.Mkdir(mainRepository, 0o755); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, mainRepository, "init", "-q")
	runGitForTest(t, mainRepository, "config", "user.email", "test@example.invalid")
	runGitForTest(t, mainRepository, "config", "user.name", "Test")
	mustWriteFile(t, filepath.Join(mainRepository, "tracked.txt"), "tracked\n")
	runGitForTest(t, mainRepository, "add", "tracked.txt")
	runGitForTest(t, mainRepository, "commit", "-q", "-m", "initial")
	runGitForTest(t, mainRepository, "worktree", "add", "-q", "-b", "linked", linkedRepository)
	runGitForTest(t, mainRepository, "checkout", "-q", "--detach")

	observer := newSystemRepoObserver(SystemDependencies().Files)
	mainObservation, err := observer.Observe(mainRepository)
	if err != nil {
		t.Fatal(err)
	}
	linkedObservation, err := observer.Observe(linkedRepository)
	if err != nil {
		t.Fatal(err)
	}
	if mainObservation.Locator == linkedObservation.Locator || mainObservation.GitCommonDir != linkedObservation.GitCommonDir {
		t.Fatalf("main = %#v, linked = %#v", mainObservation, linkedObservation)
	}
}

func TestSystemRepoObserverIgnoresInheritedGitDirectory(t *testing.T) {
	parent := physicalPath(t, t.TempDir())
	repository, other := filepath.Join(parent, "repository"), filepath.Join(parent, "other")
	for _, directory := range []string{repository, other} {
		if err := os.Mkdir(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		runGitForTest(t, directory, "init", "-q")
	}
	t.Setenv("GIT_DIR", filepath.Join(other, ".git"))

	observation, err := newSystemRepoObserver(SystemDependencies().Files).Observe(repository)

	if err != nil || observation != (RepoObservation{Locator: repository, GitCommonDir: filepath.Join(repository, ".git")}) {
		t.Fatalf("observation = %#v, %v", observation, err)
	}
}

type fakeGitResponse struct {
	stdout, stderr string
	exit           int
	err            error
}

type fakeGitRunner struct {
	t         *testing.T
	path      string
	responses map[string]fakeGitResponse
	calls     int
}

func (runner *fakeGitRunner) Run(command process.Command) error {
	runner.t.Helper()
	runner.calls++
	environment := "\n" + strings.Join(command.Env, "\n") + "\n"
	if command.Name != "git" || len(command.Args) < 2 || command.Args[0] != "-C" || command.Args[1] != runner.path ||
		command.Context == nil || !strings.Contains(environment, "\nLC_ALL=C\n") || strings.Contains(environment, "\nGIT_DIR=") {
		runner.t.Fatalf("unexpected Git request %#v", command)
	}
	response, ok := runner.responses[strings.Join(command.Args[2:], " ")]
	if !ok {
		runner.t.Fatalf("unexpected Git arguments %v", command.Args)
	}
	_, _ = io.WriteString(command.Stdout, response.stdout)
	_, _ = io.WriteString(command.Stderr, response.stderr)
	if response.exit != 0 {
		return exec.Command("/bin/sh", "-c", fmt.Sprintf("exit %d", response.exit)).Run()
	}
	return response.err
}

func TestRepoObserverClassifiesIdentifyOutcomes(t *testing.T) {
	t.Setenv("GIT_DIR", "/outside/.git")
	repository := physicalPath(t, t.TempDir())
	common := physicalPath(t, t.TempDir())
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(repository, link); err != nil {
		t.Fatal(err)
	}
	const inside, top, commonDir = "rev-parse --is-inside-work-tree", "rev-parse --path-format=absolute --show-toplevel", "rev-parse --path-format=absolute --git-common-dir"
	worktree := fakeGitResponse{stdout: "true\n"}
	exit128 := fakeGitResponse{exit: 128}
	tests := []struct {
		name      string
		responses map[string]fakeGitResponse
		calls     int
		want      ErrorClass
	}{
		{"three local identity commands", map[string]fakeGitResponse{inside: worktree, top: {stdout: repository + "\n"}, commonDir: {stdout: common + "\n"}}, 3, ""},
		{"toplevel canonicalized before comparison", map[string]fakeGitResponse{inside: worktree, top: {stdout: link + "\n"}, commonDir: {stdout: common + "\n"}}, 3, ""},
		{"malformed inside", map[string]fakeGitResponse{inside: {stdout: "true\nextra\n"}}, 1, ErrorProjectRepoObservation},
		{"launch failure", map[string]fakeGitResponse{inside: {err: errors.New("git unavailable")}}, 1, ErrorProjectRepoObservation},
		{"dubious ownership", map[string]fakeGitResponse{inside: {stderr: "fatal: detected dubious ownership in repository\n", exit: 128}}, 1, ErrorProjectRepoInvalid},
		{"not worktree", map[string]fakeGitResponse{inside: {stdout: "false\n"}}, 1, ErrorProjectRepoInvalid},
		{"toplevel exit", map[string]fakeGitResponse{inside: worktree, top: exit128}, 2, ErrorProjectRepoObservation},
		{"other toplevel wins over common-dir failure", map[string]fakeGitResponse{inside: worktree, top: {stdout: common + "\n"}, commonDir: exit128}, 3, ErrorProjectRepoInvalid},
		{"common-dir exit", map[string]fakeGitResponse{inside: worktree, top: {stdout: repository + "\n"}, commonDir: exit128}, 3, ErrorProjectRepoObservation},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &fakeGitRunner{t: t, path: repository, responses: test.responses}

			observation, err := systemRepoObserver{files: SystemDependencies().Files, runner: runner}.Observe(repository)

			if runner.calls != test.calls {
				t.Fatalf("made %d Git calls, want %d", runner.calls, test.calls)
			}
			if test.want != "" {
				assertProjectErrorClass(t, err, test.want)
			} else if err != nil || observation != (RepoObservation{Locator: repository, GitCommonDir: common}) {
				t.Fatalf("observation = %#v, %v", observation, err)
			}
		})
	}
}

func createProjectWorkspace(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, MarkerDirectory), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(root, MarkerDirectory, MarkerFile), fmt.Sprintf("format_version: 1\nroot: %s\n", root))
	mustWriteFile(t, filepath.Join(root, DefinitionsFile), "{\"format_version\": 1, \"projects\": []}\n")
	return root
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGitForTest(t *testing.T, directory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	if directory != "" {
		command.Dir = directory
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
}

func assertProjectErrorClass(t *testing.T, err error, class ErrorClass) {
	t.Helper()
	var typed *ProjectError
	if !errors.As(err, &typed) || typed.Class != class {
		t.Fatalf("error = %v, want class %s", err, class)
	}
}

func TestParseRepoIDUsesTheIdentifierRule(t *testing.T) {
	if id, err := ParseRepoID("ply-code-map"); err != nil || id != RepoID("ply-code-map") {
		t.Fatalf("id %q, error %v", id, err)
	}
	for _, value := range []string{"", "Ply", "1repo", "repo-", "a/b", "repo_one"} {
		if _, err := ParseRepoID(value); err == nil || !strings.HasPrefix(err.Error(), string(ErrorProjectInvalidArguments)+": repository ID") {
			t.Fatalf("invalid ID %q error = %v", value, err)
		}
	}
}

// project add writes the Project to ply-workspace.json with paths relative to
// the workspace and the origin remote. A repeat is no change; a difference is
// a conflict that leaves the file untouched.
func TestAddProjectWritesDefinitions(t *testing.T) {
	d, root := discoverFixture(t)
	discoverRepo(t, filepath.Join(root, "api", "main"), "git@example.invalid:api.git")
	discoverRepo(t, filepath.Join(root, "web"), "")
	if _, err := Init(d); err != nil {
		t.Fatal(err)
	}
	input := ProjectAddInput{ProjectID: "trip", Name: "Trip", Wrapper: ".", Repos: []RepoInput{{RepoID: "web", Path: "web"}, {RepoID: "api", Path: "api/main"}}}
	result, err := AddProject(d, input)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || result.Project.Wrapper != root || !reflect.DeepEqual(result.Project.RepoIDs, []RepoID{"api", "web"}) || result.Repos[0].Locator != filepath.Join(root, "api", "main") || result.Repos[0].Remote != "git@example.invalid:api.git" {
		t.Fatalf("result = %#v", result)
	}
	defs, _, err := workspacedef.Read(root)
	if err != nil || defs == nil {
		t.Fatalf("definitions %v %v", defs, err)
	}
	want := []workspacedef.Project{{ID: "trip", Name: "Trip", Wrapper: ".", Repos: []workspacedef.Repo{
		{ID: "api", Path: "api/main", Remote: "git@example.invalid:api.git"},
		{ID: "web", Path: "web"},
	}}}
	if !reflect.DeepEqual(defs.Projects, want) {
		t.Fatalf("projects %#v", defs.Projects)
	}
	path := filepath.Join(root, DefinitionsFile)
	before, _ := os.ReadFile(path)
	retry, err := AddProject(d, input)
	if err != nil || retry.Created {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(filepath.Join(root, "api", "main"), alias); err != nil {
		t.Fatal(err)
	}
	input.Repos[1].Path = alias
	if retry, err := AddProject(d, input); err != nil || retry.Created {
		t.Fatalf("retry through alias = %#v, %v", retry, err)
	}
	input.Name = "Other"
	_, err = AddProject(d, input)
	assertProjectErrorClass(t, err, ErrorProjectConflict)
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Fatal("retry or conflict rewrote the definitions")
	}
}

// project add-member appends to the defined Project in ply-workspace.json.
func TestAddProjectMemberWritesDefinitions(t *testing.T) {
	d, root := discoverFixture(t)
	discoverRepo(t, filepath.Join(root, "api"), "")
	discoverRepo(t, filepath.Join(root, "web"), "git@example.invalid:web.git")
	if _, err := Init(d); err != nil {
		t.Fatal(err)
	}
	if _, err := AddProject(d, ProjectAddInput{ProjectID: "trip", Name: "Trip", Wrapper: ".", Repos: []RepoInput{{RepoID: "api", Path: "api"}}}); err != nil {
		t.Fatal(err)
	}
	result, err := AddProjectMember(d, "trip", RepoInput{RepoID: "web", Path: "web"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || !reflect.DeepEqual(result.Project.RepoIDs, []RepoID{"api", "web"}) {
		t.Fatalf("result = %#v", result)
	}
	defs, _, err := workspacedef.Read(root)
	if err != nil || defs == nil || !reflect.DeepEqual(defs.Projects[0].Repos, []workspacedef.Repo{{ID: "api", Path: "api"}, {ID: "web", Path: "web", Remote: "git@example.invalid:web.git"}}) {
		t.Fatalf("definitions %#v %v", defs, err)
	}
	path := filepath.Join(root, DefinitionsFile)
	before, _ := os.ReadFile(path)
	if retry, err := AddProjectMember(d, "trip", RepoInput{RepoID: "web", Path: "web"}); err != nil || retry.Created {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	if _, err := AddProjectMember(d, "trip", RepoInput{RepoID: "web", Path: "api"}); err == nil {
		t.Fatal("other path for a defined member was accepted")
	}
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Fatal("retry or conflict rewrote the definitions")
	}
}

// A workspace without ply-workspace.json has no Projects to read.
func TestProjectsWithoutDefinitionsFail(t *testing.T) {
	root := createProjectWorkspace(t)
	if err := os.Remove(filepath.Join(root, DefinitionsFile)); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	_, err := ListProjects(d)
	if err == nil || !strings.Contains(err.Error(), DefinitionsFile) || !strings.Contains(err.Error(), "ply workspace init") {
		t.Fatalf("list without definitions: %v", err)
	}
}

// defineTestProject writes project and its repository locators into
// ply-workspace.json, creating the file when the fixture has none.
func defineTestProject(t *testing.T, root string, project ProjectRecord, repos []RepoRecord) {
	t.Helper()
	if err := workspacedef.Create(root); err != nil {
		t.Fatal(err)
	}
	defined := workspacedef.Project{ID: string(project.ID), Name: project.Name, Wrapper: definitionPath(root, project.Wrapper), Repos: []workspacedef.Repo{}}
	for _, r := range repos {
		defined.Repos = append(defined.Repos, workspacedef.Repo{ID: string(r.ID), Path: definitionPath(root, r.Locator)})
	}
	if _, err := updateDefinitions(root, defineProject(defined)); err != nil {
		t.Fatal(err)
	}
}

// ply-workspace.json holds only workspace structure: an Epic is adopted
// without any definition of it, and the definitions stay unchanged.
func TestAdoptEpicNeedsNoDefinition(t *testing.T) {
	d, root := discoverFixture(t)
	repo := filepath.Join(root, "api")
	discoverRepo(t, repo, "")
	if _, err := Init(d); err != nil {
		t.Fatal(err)
	}
	if _, err := AddProject(d, ProjectAddInput{ProjectID: "trip", Name: "Trip", Wrapper: ".", Repos: []RepoInput{{RepoID: "api", Path: "api"}}}); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(root, DefinitionsFile))
	oid := strings.TrimSpace(runLocalGit(t, repo, "rev-parse", "HEAD"))
	adopted, err := AdoptEpic(d, EpicAdoptInput{EpicID: "main", Title: "Main", ProjectID: "trip", RepoID: "api", Worktree: repo, Ref: "refs/heads/main", ExpectedOID: oid})
	if err != nil || !adopted.Created {
		t.Fatalf("adopt = %#v, %v", adopted, err)
	}
	if after, _ := os.ReadFile(filepath.Join(root, DefinitionsFile)); string(after) != string(before) {
		t.Fatalf("adoption changed the definitions: %s", after)
	}
}
