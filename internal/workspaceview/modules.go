package workspaceview

import (
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const ModuleListKind = "WorkspaceModuleListReadback@1"

// ModuleRepoUnavailable is the reason of a repo whose checkout cannot be read.
const ModuleRepoUnavailable = "repo_unavailable"

type ModuleOptions struct {
	ProjectID, RepoID, Commit string
}

// Module is a directory with pom.xml and its nearest enclosing module.
type Module struct {
	Module string  `json:"module"`
	Parent *string `json:"parent"`
}

type ModuleRepo struct {
	Repo    string   `json:"repo"`
	Commit  *string  `json:"commit"`
	Modules []Module `json:"modules"`
	Reason  *string  `json:"reason"`
}

type ModuleListResult struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	Project       string       `json:"project"`
	Repos         []ModuleRepo `json:"repos"`
}

// ListModules reads the module tree of each selected repo from a commit:
// HEAD of the registered checkout, or o.Commit for one repo. A repo that
// cannot be read gets a reason on its row and does not stop the others.
func ListModules(d taskwork.Dependencies, o ModuleOptions) (ModuleListResult, error) {
	out := ModuleListResult{Kind: ModuleListKind, SchemaVersion: 1, Project: o.ProjectID, Repos: []ModuleRepo{}}
	if o.ProjectID == "" {
		return out, workspace.WorkInvalidArguments("--project is required")
	}
	if o.Commit != "" && o.RepoID == "" {
		return out, workspace.WorkInvalidArguments("--commit requires --repo")
	}
	id, err := moduleProject(d, o)
	if err != nil {
		return out, err
	}
	shown, err := workspace.ShowProjectMember(d.Dependencies, id, workspace.RepoID(o.RepoID))
	if err != nil {
		return out, err
	}
	repos := map[workspace.RepoID]workspace.RepoRecord{}
	for _, r := range shown.Repos {
		repos[r.ID] = r
	}
	for _, repoID := range shown.Project.RepoIDs {
		if o.RepoID != "" && string(repoID) != o.RepoID {
			continue
		}
		out.Repos = append(out.Repos, readModuleRepo(d, repos, repoID, o.Commit))
	}
	return out, nil
}

// moduleProject is the Project of o, checked against the registered Projects
// without Git: an unknown Project or repository is not found, and a repository
// of another Project is an invalid argument.
func moduleProject(d taskwork.Dependencies, o ModuleOptions) (workspace.ProjectID, error) {
	id, err := workspace.ParseProjectID(o.ProjectID)
	if err != nil {
		return "", err
	}
	if o.RepoID != "" {
		if err := workspace.ValidateIdentifier("repository", o.RepoID); err != nil {
			return "", err
		}
	}
	listed, err := workspace.ListProjects(d.Dependencies)
	if err != nil {
		return "", err
	}
	repo := workspace.RepoID(o.RepoID)
	var members []workspace.RepoID
	found, registered := false, repo == ""
	for _, p := range listed.Projects {
		if p.ID == id {
			found, members = true, p.RepoIDs
		}
		registered = registered || slices.Contains(p.RepoIDs, repo)
	}
	switch {
	case !found:
		return "", notFound("project", o.ProjectID)
	case !registered:
		return "", notFound("repository", o.RepoID)
	case repo != "" && !slices.Contains(members, repo):
		return "", workspace.WorkInvalidArguments("--repo must be a repo ID of project " + o.ProjectID)
	}
	return id, nil
}

// readModuleRepo reads the module tree of member id, which project show has
// already observed.
func readModuleRepo(d taskwork.Dependencies, repos map[workspace.RepoID]workspace.RepoRecord, id workspace.RepoID, commit string) ModuleRepo {
	row := ModuleRepo{Repo: string(id)}
	repo, ok := repos[id]
	if !ok || repo.Unavailable != "" {
		return withReason(row, ModuleRepoUnavailable)
	}
	tree := taskwork.ReadRegisteredCommitTree(d, repo, commit)
	if tree.Reason != "" {
		return withReason(row, tree.Reason)
	}
	row.Commit, row.Modules = &tree.OID, ModuleTree(tree.Entries)
	return row
}

func withReason(row ModuleRepo, reason string) ModuleRepo {
	row.Reason = &reason
	return row
}

// ModuleTree is the module rule: every directory holding a regular pom.xml,
// outside directories plycq ignores, with its nearest enclosing module as
// parent. It is sorted by module path and never nil.
func ModuleTree(entries []taskwork.CommitTreeEntry) []Module {
	found := map[string]bool{}
	for _, e := range entries {
		if (e.Mode == "100644" || e.Mode == "100755") && path.Base(e.Path) == "pom.xml" && !moduleExcluded(e.Path) {
			found[path.Dir(e.Path)] = true
		}
	}
	modules := make([]Module, 0, len(found))
	for dir := range found {
		m := Module{Module: dir}
		for parent := dir; parent != "."; {
			parent = path.Dir(parent)
			if found[parent] {
				m.Parent = &parent
				break
			}
		}
		modules = append(modules, m)
	}
	sort.Slice(modules, func(i, j int) bool { return modules[i].Module < modules[j].Module })
	return modules
}

// moduleExcluded mirrors plycq: a directory that starts with a dot or is
// vendor, node_modules or testdata anywhere in the path, or target or build
// before the first src.
func moduleExcluded(file string) bool {
	dirs := strings.Split(file, "/")
	dirs = dirs[:len(dirs)-1]
	beforeSrc := true
	for _, name := range dirs {
		if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "testdata" {
			return true
		}
		if name == "src" {
			beforeSrc = false
		}
		if beforeSrc && (name == "target" || name == "build") {
			return true
		}
	}
	return false
}
