package workspace

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

var epicCreateWorktreeEffects = map[string]string{"no_effect": "no_effect", "branch_only": "partial", "exact_effect": "completed", "partial_or_unknown": "unknown"}

func epicDefinitionsChanged(root string) error {
	return workError(ErrorWorkIdentityConflict, fmt.Sprintf("%s changed during Epic create and was not written by Ply; the worktree is kept; repeat the identical create request", filepath.Join(root, DefinitionsFile)), nil)
}

func definitionLocator(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, path)
}

func epicCreateDefinition(root string, in EpicCreateInput) workspacedef.Epic {
	path := in.Path
	if rel, err := filepath.Rel(root, in.Path); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		path = rel
	}
	return workspacedef.Epic{ID: string(in.EpicID), Title: in.Title, Repo: string(in.RepoID), Path: path, Ref: "refs/heads/" + in.Branch}
}

func epicDefined(defs *workspacedef.Definitions, root string, in EpicCreateInput) (bool, error) {
	path := filepath.Join(root, DefinitionsFile)
	want := epicCreateDefinition(root, in)
	member, defined := false, false
	for _, p := range defs.Projects {
		for _, e := range p.Epics {
			if e.ID != string(in.EpicID) {
				continue
			}
			if p.ID != string(in.ProjectID) || e.Title != want.Title || e.Repo != want.Repo || e.Ref != want.Ref || definitionLocator(root, e.Path) != in.Path {
				return false, WorkInvalidArguments(fmt.Sprintf("Epic %s is already defined differently in %s; select a new Epic ID, or repeat the create request exactly as defined there", in.EpicID, path))
			}
			defined = true
		}
		for _, r := range p.Repos {
			member = member || p.ID == string(in.ProjectID) && r.ID == string(in.RepoID)
		}
	}
	if !member {
		return false, WorkInvalidArguments(fmt.Sprintf("repository %s of Project %s is not defined in %s; add it to that Project there, then repeat the create request", in.RepoID, in.ProjectID, path))
	}
	return defined, nil
}

// recordEpicDefinition appends the Epic to its Project through the single
// definitions writer, which owns projects.lock for the whole update.
func recordEpicDefinition(root, basis string, in EpicCreateInput) (string, error) {
	path := filepath.Join(root, DefinitionsFile)
	retry := "the worktree is kept; repeat the identical create request"
	mutate := func(defs *workspacedef.Definitions) (bool, error) {
		if defined, err := epicDefined(defs, root, in); err != nil || defined {
			return false, err
		}
		for i := range defs.Projects {
			if defs.Projects[i].ID == string(in.ProjectID) {
				defs.Projects[i].Epics = append(defs.Projects[i].Epics, epicCreateDefinition(root, in))
			}
		}
		return true, nil
	}
	result, err := workspacedef.Update(root, basis, mutate, validateDefinitions)
	if err == nil {
		return string(result.Effect), nil
	}
	var release *filelock.ReleaseError
	var defErr *workspacedef.Error
	if errors.As(err, &release) || !errors.As(err, &defErr) {
		return string(result.Effect), definitionUpdateError(err)
	}
	switch defErr.Kind {
	case workspacedef.KindMissing, workspacedef.KindConflict:
		return string(result.Effect), epicDefinitionsChanged(root)
	case workspacedef.KindPublish:
		return string(result.Effect), workError(ErrorWorkIO, fmt.Sprintf("could not update %s; %s", path, retry), defErr.Err)
	case workspacedef.KindReadback:
		return string(result.Effect), workError(ErrorWorkReconciliation, fmt.Sprintf("%s holds unexpected content after the Epic definition update; inspect it, then repeat the identical create request", path), defErr.Err)
	default:
		return string(result.Effect), definitionUpdateError(err)
	}
}
