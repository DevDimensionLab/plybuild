package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProjectInitInput selects a read-only discovery or an exact confirmed apply.
type ProjectInitInput struct {
	ProjectID    ProjectID
	Name         string
	Wrapper      string
	Apply        bool
	Confirmation string
}
type ProjectInitRepository struct {
	ID           RepoID `json:"repo_id"`
	Locator      string `json:"locator"`
	GitCommonDir string `json:"git_common_dir"`
}
type ProjectInitOmission struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}
type ProjectInitPlan struct {
	Kind         string                  `json:"kind"`
	Workspace    string                  `json:"workspace"`
	ProjectID    ProjectID               `json:"project_id"`
	Name         string                  `json:"name"`
	Wrapper      string                  `json:"wrapper"`
	Repositories []ProjectInitRepository `json:"repositories"`
	Omissions    []ProjectInitOmission   `json:"omissions"`
	Blockers     []string                `json:"blockers"`
	Confirmation string                  `json:"confirmation"`
	Applied      bool                    `json:"applied"`
	NextAction   string                  `json:"next_action"`
}

// InitProject discovers only immediate children. All registration effects go
// through AddProject; applyProject is reused for read-only conflict validation.
func InitProject(d Dependencies, in ProjectInitInput) (ProjectInitPlan, error) {
	p := ProjectInitPlan{Kind: "WorkspaceProjectInitPlan@1", ProjectID: in.ProjectID,
		Repositories: []ProjectInitRepository{}, Omissions: []ProjectInitOmission{}, Blockers: []string{},
		NextAction: "Use ply workspace project add <project-id> --name <name> --wrapper <wrapper> --repo <repo-id>=<worktree-root>; for an existing Project use ply workspace project add-member <project-id> --repo <repo-id>=<worktree-root>."}
	if _, err := ParseProjectID(string(in.ProjectID)); err != nil {
		return p, err
	}
	if in.Apply && in.Confirmation == "" || !in.Apply && in.Confirmation != "" {
		return p, ProjectInvalidArguments("--apply requires --confirmation; --confirmation is only for apply")
	}
	if d.Files == nil || d.Repos == nil || d.Projects == nil {
		return p, ProjectInvalidArguments("filesystem, repository observer and project store are required")
	}
	cwd, err := projectPhysicalWorkingDirectory(d.Files)
	if err != nil {
		return p, err
	}
	p.Workspace, err = locateContainingWorkspace(d.Files, cwd)
	if err != nil {
		return p, err
	}
	wrapper := in.Wrapper
	if wrapper == "" {
		wrapper = cwd
	}
	p.Wrapper, err = canonicalDirectory(d.Files, cwd, wrapper, "wrapper")
	if err != nil {
		return p, err
	}
	projects, repos, err := d.Projects.Snapshot(p.Workspace)
	if err != nil {
		return p, err
	}
	p.Name = in.Name
	if p.Name == "" {
		p.Name = string(in.ProjectID)
		for _, project := range projects {
			if project.ID == in.ProjectID {
				p.Name = project.Name
			}
		}
	}
	if err = validateName(p.Name); err != nil {
		return p, err
	}
	metadata, err := readProjectMetadata(p.Workspace)
	if err != nil {
		return p, err
	}
	companions := map[string]bool{}
	for _, project := range projects {
		for _, c := range metadata.Companions(project.ID) {
			companions[c.Locator] = true
		}
	}
	// A wrapper that is itself a worktree is outside this direct-child layout.
	if _, err = d.Files.Lstat(filepath.Join(p.Wrapper, ".git")); err == nil {
		p.Blockers = append(p.Blockers, fmt.Sprintf("wrapper %s itself has a Git marker; choose explicit repository roots", p.Wrapper))
	} else if !errors.Is(err, fs.ErrNotExist) {
		p.Blockers = append(p.Blockers, fmt.Sprintf("inspect wrapper Git marker: %v", err))
	}
	entries, err := os.ReadDir(p.Wrapper)
	if err != nil {
		return p, projectPathError("read wrapper", p.Wrapper, err)
	}
	groups := map[string][]RepoObservation{}
	for _, entry := range entries {
		path := filepath.Join(p.Wrapper, entry.Name())
		omit := func(reason string) { p.Omissions = append(p.Omissions, ProjectInitOmission{path, reason}) }
		info, err := d.Files.Lstat(path)
		if err != nil {
			p.Blockers = append(p.Blockers, fmt.Sprintf("inspect %s: %v", path, err))
			continue
		}
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			omit("warning: symlink omitted; target not inspected")
			continue
		case entry.Name() == ".git" || entry.Name() == MarkerDirectory:
			omit("system metadata omitted")
			continue
		case !info.IsDir():
			omit("not a directory")
			continue
		case companions[path]:
			omit("explicitly registered companion")
			continue
		}
		marker, err := d.Files.Lstat(filepath.Join(path, ".git"))
		if errors.Is(err, fs.ErrNotExist) {
			omit("no direct Git marker; nested repository-wrapper/main layouts are not searched; use explicit project add/add-member")
			continue
		}
		if err != nil {
			p.Blockers = append(p.Blockers, fmt.Sprintf("inspect Git marker in %s: %v", path, err))
			continue
		}
		if marker.Mode()&fs.ModeSymlink != 0 {
			p.Blockers = append(p.Blockers, fmt.Sprintf("Git marker in %s is a symlink; use explicit roots", path))
			continue
		}
		observation, err := d.Repos.Observe(path)
		if err != nil {
			p.Blockers = append(p.Blockers, fmt.Sprintf("Git marker in %s: %v", path, err))
			continue
		}
		if observation.Locator != path || !validStoredPath(observation.GitCommonDir) {
			p.Blockers = append(p.Blockers, fmt.Sprintf("invalid physical Git observation for %s", path))
			continue
		}
		groups[observation.GitCommonDir] = append(groups[observation.GitCommonDir], observation)
	}
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
		id := RepoID(filepath.Base(chosen.Locator))
		if len(groups) == 1 {
			id = RepoID(in.ProjectID)
		}
		if err := validateIdentifier("repository", string(id)); err != nil {
			p.Blockers = append(p.Blockers, err.Error())
		}
		p.Repositories = append(p.Repositories, ProjectInitRepository{id, chosen.Locator, common})
		for _, other := range group[1:] {
			p.Omissions = append(p.Omissions, ProjectInitOmission{other.Locator, "same physical Git common directory as " + chosen.Locator})
		}
	}
	sort.Slice(p.Repositories, func(i, j int) bool { return p.Repositories[i].ID < p.Repositories[j].ID })
	sort.Slice(p.Omissions, func(i, j int) bool { return p.Omissions[i].Path < p.Omissions[j].Path })
	if len(p.Repositories) == 0 {
		p.Blockers = append(p.Blockers, "no direct-child Git worktree roots; nested repository-wrapper/main layouts are not searched")
	}
	input := ProjectAddInput{ProjectID: p.ProjectID, Name: p.Name, Wrapper: p.Wrapper}
	project := ProjectRecord{ID: p.ProjectID, Name: p.Name, Wrapper: p.Wrapper}
	records := []RepoRecord{}
	for _, repo := range p.Repositories {
		input.Repos = append(input.Repos, RepoInput{RepoID: repo.ID, Path: repo.Locator})
		project.RepoIDs = append(project.RepoIDs, repo.ID)
		records = append(records, RepoRecord{ID: repo.ID, Locator: repo.Locator, GitCommonDir: repo.GitCommonDir})
	}
	if len(p.Blockers) == 0 {
		if err := validateProjectAddInput(input); err != nil {
			p.Blockers = append(p.Blockers, err.Error())
		} else if _, _, _, _, err := applyProject(projectRegistry{Projects: projects, Repos: repos}, project, records); err != nil {
			p.Blockers = append(p.Blockers, err.Error())
		}
	}
	sort.Strings(p.Blockers)
	if len(p.Blockers) > 0 {
		return p, projectError(ErrorProjectConflict, strings.Join(p.Blockers, "; ")+". "+p.NextAction, nil)
	}
	// Bind only stable registration intent. Diagnostics and omitted files can
	// change without changing membership, physical identity or explicit choices.
	p.Confirmation = queueDigest(struct {
		Workspace    string                  `json:"workspace"`
		ProjectID    ProjectID               `json:"project_id"`
		Name         string                  `json:"name"`
		Wrapper      string                  `json:"wrapper"`
		Repositories []ProjectInitRepository `json:"repositories"`
	}{p.Workspace, p.ProjectID, p.Name, p.Wrapper, p.Repositories})
	if !in.Apply {
		return p, nil
	}
	if in.Confirmation != p.Confirmation {
		return p, projectError(ErrorProjectConflict, "preview confirmation changed; read a fresh preview before --apply. "+p.NextAction, nil)
	}
	// AddProject reobserves every selected root before its sole store mutation.
	// Bind those observations to this confirmed plan to reject identity drift.
	d.Repos = projectInitConfirmedObserver{RepoObserver: d.Repos, expected: records}
	if _, err := AddProject(d, input); err != nil {
		return p, fmt.Errorf("%w. %s", err, p.NextAction)
	}
	p.Applied = true
	return p, nil
}

type projectInitConfirmedObserver struct {
	RepoObserver
	expected []RepoRecord
}

func (o projectInitConfirmedObserver) Observe(path string) (RepoObservation, error) {
	observed, err := o.RepoObserver.Observe(path)
	if err != nil {
		return observed, err
	}
	for _, repo := range o.expected {
		if repo.Locator == path && observed.Locator == path && repo.GitCommonDir == observed.GitCommonDir {
			return observed, nil
		}
	}
	return RepoObservation{}, projectError(ErrorProjectConflict, "Git identity changed since confirmed discovery at "+path, nil)
}
