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
func ListModules(d taskwork.Dependencies, s *Snapshot, o ModuleOptions) (ModuleListResult, error) {
	out := ModuleListResult{Kind: ModuleListKind, SchemaVersion: 1, Project: o.ProjectID, Repos: []ModuleRepo{}}
	if o.ProjectID == "" {
		return out, workspace.WorkInvalidArguments("--project is required")
	}
	if o.Commit != "" && o.RepoID == "" {
		return out, workspace.WorkInvalidArguments("--commit requires --repo")
	}
	if err := validateProjectFilter(s, o.ProjectID); err != nil {
		return out, err
	}
	if o.RepoID != "" {
		id := workspace.RepoID(o.RepoID)
		if err := ValidateTaskFilters(s, taskwork.TaskListFilters{RepoID: &id}); err != nil {
			return out, err
		}
	}
	repos := map[workspace.RepoID]workspace.RepoRecord{}
	for _, r := range s.Projects.Repos {
		repos[r.ID] = r
	}
	for _, project := range s.Projects.Projects {
		if string(project.ID) != o.ProjectID {
			continue
		}
		if o.RepoID != "" && !slices.Contains(project.RepoIDs, workspace.RepoID(o.RepoID)) {
			return out, workspace.WorkInvalidArguments("--repo must be a repo ID of project " + o.ProjectID)
		}
		for _, id := range project.RepoIDs {
			if o.RepoID != "" && string(id) != o.RepoID {
				continue
			}
			out.Repos = append(out.Repos, readModuleRepo(d, repos, id, o.Commit))
		}
	}
	return out, nil
}

func readModuleRepo(d taskwork.Dependencies, repos map[workspace.RepoID]workspace.RepoRecord, id workspace.RepoID, commit string) ModuleRepo {
	row := ModuleRepo{Repo: string(id)}
	repo, ok := repos[id]
	if !ok {
		return withReason(row, ModuleRepoUnavailable)
	}
	if availability, _ := workspace.ObserveMember(repo); availability != "" {
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
