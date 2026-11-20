package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// InitDiscoverResult reports the workspace, the Project defined for its
// directory and what discovery left out.
type InitDiscoverResult struct {
	Root      string
	ProjectID ProjectID
	Repos     []workspacedef.Repo
	Omissions []DiscoveryOmission
	Changed   bool
}

type DiscoveryOmission struct {
	Path   string
	Reason string
}

// InitDiscover initializes the current directory as a workspace and defines
// one Project in ply-workspace.json with wrapper "." and the Git repositories
// found in its children: <name> or <name>/main, with repository ID <name> and
// the URL of origin as remote. The Project ID defaults to the directory name.
// A rerun with the same repositories changes nothing.
func InitDiscover(d Dependencies, project string) (InitDiscoverResult, error) {
	if d.Files == nil || d.Repos == nil {
		return InitDiscoverResult{}, ProjectInvalidArguments("filesystem and repository observer are required")
	}
	remotes, ok := d.Repos.(RemoteObserver)
	if !ok {
		return InitDiscoverResult{}, ProjectInvalidArguments("the repository observer cannot read remotes")
	}
	root, err := physicalWorkingDirectory(d.Files)
	if err != nil {
		return InitDiscoverResult{}, err
	}
	if project == "" {
		project = filepath.Base(root)
	}
	id, err := ParseProjectID(project)
	if err != nil {
		return InitDiscoverResult{}, fmt.Errorf("%w; choose a Project ID with --project", err)
	}
	result := InitDiscoverResult{Root: root, ProjectID: id}
	found, omissions, err := discoverRepos(d, root)
	result.Omissions = omissions
	if err != nil {
		return result, err
	}
	for _, repo := range found {
		remote, err := remotes.Remote(filepath.Join(root, filepath.FromSlash(repo.Path)))
		if err != nil {
			return result, err
		}
		repo.Remote = remote
		result.Repos = append(result.Repos, repo)
	}
	if _, err = Init(d); err != nil {
		return result, err
	}
	defined := workspacedef.Project{ID: string(id), Name: string(id), Wrapper: ".", Repos: result.Repos}
	result.Changed, err = updateDefinitions(root, defineProject(defined))
	return result, err
}

// discoverRepos finds the Git worktree root in each child of root, <name> or
// else <name>/main. Several roots of one repository keep the ordinary root,
// then the lexically first.
func discoverRepos(d Dependencies, root string) ([]workspacedef.Repo, []DiscoveryOmission, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, projectPathError("read directory", root, err)
	}
	omissions := []DiscoveryOmission{}
	var blockers []string
	type candidate struct {
		name string
		RepoObservation
	}
	groups := map[string][]candidate{}
	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		omit := func(reason string) { omissions = append(omissions, DiscoveryOmission{path, reason}) }
		info, err := d.Files.Lstat(path)
		if err != nil {
			blockers = append(blockers, fmt.Sprintf("inspect %s: %v", path, err))
			continue
		}
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			omit("symlink; target not inspected")
			continue
		case strings.HasPrefix(entry.Name(), "."):
			omit("hidden directory or file")
			continue
		case !info.IsDir():
			omit("not a directory")
			continue
		}
		worktree := path
		marker, err := d.Files.Lstat(filepath.Join(worktree, ".git"))
		if errors.Is(err, fs.ErrNotExist) {
			worktree = filepath.Join(path, "main")
			marker, err = d.Files.Lstat(filepath.Join(worktree, ".git"))
		}
		if errors.Is(err, fs.ErrNotExist) {
			omit("no Git repository in " + entry.Name() + " or " + entry.Name() + "/main")
			continue
		}
		if err != nil {
			blockers = append(blockers, fmt.Sprintf("inspect Git marker in %s: %v", worktree, err))
			continue
		}
		if marker.Mode()&fs.ModeSymlink != 0 {
			blockers = append(blockers, fmt.Sprintf("Git marker in %s is a symlink", worktree))
			continue
		}
		observation, err := d.Repos.Observe(worktree)
		if err != nil {
			blockers = append(blockers, fmt.Sprintf("Git repository in %s: %v", worktree, err))
			continue
		}
		if observation.Locator != worktree || !ValidStoredPath(observation.GitCommonDir) {
			blockers = append(blockers, fmt.Sprintf("invalid physical Git observation for %s", worktree))
			continue
		}
		groups[observation.GitCommonDir] = append(groups[observation.GitCommonDir], candidate{entry.Name(), observation})
	}
	repos := []workspacedef.Repo{}
	for common, group := range groups {
		sort.Slice(group, func(i, j int) bool {
			ordinaryI := filepath.Join(group[i].Locator, ".git") == common
			ordinaryJ := filepath.Join(group[j].Locator, ".git") == common
			if ordinaryI != ordinaryJ {
				return ordinaryI
			}
			return group[i].Locator < group[j].Locator
		})
		chosen := group[0]
		if err := ValidateIdentifier("repository", chosen.name); err != nil {
			blockers = append(blockers, err.Error())
		}
		relative, _ := filepath.Rel(root, chosen.Locator)
		repos = append(repos, workspacedef.Repo{ID: chosen.name, Path: filepath.ToSlash(relative)})
		for _, other := range group[1:] {
			omissions = append(omissions, DiscoveryOmission{other.Locator, "same Git repository as " + chosen.Locator})
		}
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].ID < repos[j].ID })
	sort.Slice(omissions, func(i, j int) bool { return omissions[i].Path < omissions[j].Path })
	if len(repos) == 0 && len(blockers) == 0 {
		blockers = append(blockers, fmt.Sprintf("no Git repositories in the children of %s or their main directories", root))
	}
	if len(blockers) > 0 {
		sort.Strings(blockers)
		return nil, omissions, projectError(ErrorProjectConflict, strings.Join(blockers, "; "), nil)
	}
	return repos, omissions, nil
}
