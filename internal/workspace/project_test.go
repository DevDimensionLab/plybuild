package workspace

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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

func TestAddProjectCanonicalizesSortsAndIsIdempotentByGitIdentity(t *testing.T) {
	root := createProjectWorkspace(t)
	child := filepath.Join(root, "child")
	wrapper := filepath.Join(t.TempDir(), "wrapper")
	repoA := filepath.Join(t.TempDir(), "repo-a")
	repoB := filepath.Join(t.TempDir(), "repo-b")
	repoAlias := filepath.Join(t.TempDir(), "repo-alias")
	for _, path := range []string{child, wrapper, repoA, repoB, repoAlias} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	child = physicalPath(t, child)
	wrapper = physicalPath(t, wrapper)
	repoA = physicalPath(t, repoA)
	repoB = physicalPath(t, repoB)
	repoAlias = physicalPath(t, repoAlias)
	commonA := physicalPath(t, t.TempDir())
	commonB := physicalPath(t, t.TempDir())
	observer := &fakeRepoObserver{observations: map[string]RepoObservation{
		repoA:     {Locator: repoA, GitCommonDir: commonA},
		repoAlias: {Locator: repoAlias, GitCommonDir: commonA},
		repoB:     {Locator: repoB, GitCommonDir: commonB},
	}, errors: map[string]error{}}
	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: child}
	dependencies.Repos = observer

	input := ProjectAddInput{ProjectID: "trip", Name: "Trip", Wrapper: wrapper, Repos: []RepoInput{{RepoID: "service", Path: repoB}, {RepoID: "api", Path: repoA}}}
	result, err := AddProject(dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || result.Workspace != root || !reflect.DeepEqual(result.Project.RepoIDs, []RepoID{"api", "service"}) {
		t.Fatalf("created result = %#v", result)
	}
	if got := []RepoID{result.Repos[0].ID, result.Repos[1].ID}; !reflect.DeepEqual(got, []RepoID{"api", "service"}) {
		t.Fatalf("sorted repositories = %#v", got)
	}

	input.Repos = []RepoInput{{RepoID: "api", Path: repoAlias}, {RepoID: "service", Path: repoB}}
	retry, err := AddProject(dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if retry.Created || retry.Repos[0].Locator != repoA {
		t.Fatalf("idempotent result = %#v", retry)
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
	if _, err := os.Lstat(filepath.Join(root, MarkerDirectory, ProjectsFile)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("registry exists after duplicate input: %v", err)
	}
}

func TestRelativeAddPathsUsePhysicalCWDAndReadbackDoesNotReobserve(t *testing.T) {
	root := createProjectWorkspace(t)
	cwd := filepath.Join(root, "nested", "cwd")
	wrapper := filepath.Join(root, "wrapper")
	repository := physicalPath(t, t.TempDir())
	for _, path := range []string{cwd, wrapper} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cwd = physicalPath(t, cwd)
	wrapper = physicalPath(t, wrapper)
	common := filepath.Join(root, "observed-common")
	observer := &fakeRepoObserver{observations: map[string]RepoObservation{
		repository: {Locator: repository, GitCommonDir: common},
	}, errors: map[string]error{}}
	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: cwd}
	dependencies.Repos = observer
	relativeWrapper, err := filepath.Rel(cwd, wrapper)
	if err != nil {
		t.Fatal(err)
	}
	relativeRepository, err := filepath.Rel(cwd, repository)
	if err != nil {
		t.Fatal(err)
	}
	result, err := AddProject(dependencies, ProjectAddInput{ProjectID: "relative", Name: "Relative", Wrapper: relativeWrapper, Repos: []RepoInput{{RepoID: "repo", Path: relativeRepository}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Project.Wrapper != wrapper || result.Repos[0].Locator != repository || !reflect.DeepEqual(observer.calls, []string{repository}) {
		t.Fatalf("relative result = %#v, observer calls = %#v", result, observer.calls)
	}
	if err := os.RemoveAll(wrapper); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(repository); err != nil {
		t.Fatal(err)
	}
	observer.errors[repository] = errors.New("readback must not observe Git")
	shown, err := ShowProject(dependencies, "relative")
	if err != nil || shown.Project.Wrapper != wrapper || shown.Repos[0].Locator != repository {
		t.Fatalf("show after deletion = %#v, %v", shown, err)
	}
	listed, err := ListProjects(dependencies)
	if err != nil || len(listed.Projects) != 1 {
		t.Fatalf("list after deletion = %#v, %v", listed, err)
	}
	if !reflect.DeepEqual(observer.calls, []string{repository}) {
		t.Fatalf("readback reobserved repositories: %#v", observer.calls)
	}
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

func TestRepoObserverRejectsMalformedOutputAndLaunchFailure(t *testing.T) {
	path := t.TempDir()
	files := SystemDependencies().Files
	malformed := systemRepoObserver{files: files, run: func(string, ...string) gitCommandResult {
		return gitCommandResult{stdout: []byte("true\nextra\n")}
	}}
	_, err := malformed.Observe(path)
	assertProjectErrorClass(t, err, ErrorProjectRepoObservation)

	launch := systemRepoObserver{files: files, run: func(string, ...string) gitCommandResult {
		return gitCommandResult{exit: -1, err: errors.New("git unavailable")}
	}}
	_, err = launch.Observe(path)
	assertProjectErrorClass(t, err, ErrorProjectRepoObservation)
}

func TestRepoObserverRunsOnlyThreeLocalIdentityCommands(t *testing.T) {
	repository := physicalPath(t, t.TempDir())
	common := physicalPath(t, t.TempDir())
	type invocation struct {
		path string
		args []string
	}
	var calls []invocation
	observer := systemRepoObserver{files: SystemDependencies().Files}
	observer.run = func(path string, arguments ...string) gitCommandResult {
		calls = append(calls, invocation{path: path, args: append([]string(nil), arguments...)})
		switch {
		case reflect.DeepEqual(arguments, []string{"rev-parse", "--is-inside-work-tree"}):
			return gitCommandResult{stdout: []byte("true\n")}
		case reflect.DeepEqual(arguments, []string{"rev-parse", "--path-format=absolute", "--show-toplevel"}):
			return gitCommandResult{stdout: []byte(repository + "\n")}
		case reflect.DeepEqual(arguments, []string{"rev-parse", "--path-format=absolute", "--git-common-dir"}):
			return gitCommandResult{stdout: []byte(common + "\n")}
		default:
			return gitCommandResult{exit: -1, err: fmt.Errorf("unexpected Git arguments %v", arguments)}
		}
	}
	observation, err := observer.Observe(repository)
	if err != nil {
		t.Fatal(err)
	}
	if observation != (RepoObservation{Locator: repository, GitCommonDir: common}) || len(calls) != 3 {
		t.Fatalf("observation = %#v, calls = %#v", observation, calls)
	}
	for _, call := range calls {
		if call.path != repository {
			t.Fatalf("Git path = %q, want %q", call.path, repository)
		}
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
