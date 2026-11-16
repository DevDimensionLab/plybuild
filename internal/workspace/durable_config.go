package workspace

import (
	"errors"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// DefinitionsFile lives beside .ply, independently of disposable runtime data.
const DefinitionsFile = workspacedef.File

// readDefinitions reads the definitions through workspacedef and applies the
// workspace Project and repository rules.
func readDefinitions(root string) (*workspacedef.Definitions, string, error) {
	defs, digest, err := workspacedef.Read(root)
	var defErr *workspacedef.Error
	if errors.As(err, &defErr) {
		return nil, "", definitionError(defErr)
	}
	if err != nil || defs == nil {
		return nil, digest, err
	}
	if err = validateDefinitions(*defs); err != nil {
		return nil, "", err
	}
	return defs, digest, nil
}

type definitionRootFiles struct {
	FileSystem
	root string
}

func (f definitionRootFiles) Getwd() (string, error) { return f.root, nil }

// Config discovery can recreate only the workspace marker. Init retains its
// existing atomic publication, ancestor validation and concurrent-init guards.
func initDefinitionWorkspace(files FileSystem, root string) (bool, error) {
	defs, _, err := readDefinitions(root)
	if err != nil || defs == nil {
		return false, err
	}
	_, err = Init(Dependencies{Files: definitionRootFiles{files, root}})
	return err == nil, err
}

func registrationProjects(d Dependencies, root string) ([]ProjectRecord, []RepoRecord, error) {
	registries, err := SnapshotWorkspace(d, root)
	if err != nil {
		return nil, nil, err
	}
	return registryProjects(root, MergeRegistries(registries))
}

// registryProjects resolves the defined Projects statically, with the Git
// common directories the registry recorded for their repositories.
func registryProjects(root string, registry WorkItemRegistry) ([]ProjectRecord, []RepoRecord, error) {
	defs, _, err := readDefinitions(root)
	if err != nil {
		return nil, nil, err
	}
	if defs == nil {
		return nil, nil, missingDefinitions(root)
	}
	var bindings []repohistory.Binding
	for _, epic := range registry.Epics {
		for _, binding := range epic.RepoBindings {
			bindings = append(bindings, repohistory.Binding{RepoID: string(binding.RepoID), GitCommonDir: binding.GitCommonDir})
		}
	}
	for _, task := range registry.Tasks {
		bindings = append(bindings, repohistory.Binding{RepoID: string(task.RepoID), GitCommonDir: task.GitCommonDir})
	}
	identities, err := repohistory.CommonDirs(bindings)
	if err != nil {
		return nil, nil, WorkInvalidArguments("conflicting historical repository identity")
	}
	projects := []ProjectRecord{}
	repos := []RepoRecord{}
	seen := map[RepoID]bool{}
	path := func(value string) string {
		if filepath.IsAbs(value) {
			return filepath.Clean(value)
		}
		return filepath.Join(root, value)
	}
	for _, p := range defs.Projects {
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: path(p.Wrapper)}
		for _, r := range p.Repos {
			project.RepoIDs = append(project.RepoIDs, RepoID(r.ID))
			if !seen[RepoID(r.ID)] {
				repos = append(repos, RepoRecord{ID: RepoID(r.ID), Locator: path(r.Path), GitCommonDir: identities[r.ID]})
				seen[RepoID(r.ID)] = true
			}
		}
		projects = append(projects, project)
	}
	return projects, repos, nil
}
