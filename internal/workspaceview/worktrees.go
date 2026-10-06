package workspaceview

import (
	"sort"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const WorktreeListKind = "WorkspaceWorktreeListReadback@1"

type WorktreeOptions struct{ ProjectID, RepoID string }
type WorktreeBase struct {
	Kind     string `json:"kind"`
	Locator  string `json:"locator"`
	Ref      string `json:"ref"`
	OID      string `json:"oid"`
	Revision int    `json:"revision"`
}
type WorktreeItem struct {
	ProjectID      string                           `json:"project_id"`
	RepoID         string                           `json:"repo_id"`
	Locator        string                           `json:"locator"`
	Ref            *string                          `json:"ref"`
	OID            *string                          `json:"oid"`
	Clean          *bool                            `json:"clean"`
	Owner          string                           `json:"owner"`
	OwnerID        *string                          `json:"owner_id"`
	OwnerFreshness string                           `json:"owner_freshness"`
	Base           *WorktreeBase                    `json:"base"`
	Ahead          *int                             `json:"ahead"`
	Behind         *int                             `json:"behind"`
	MergedIntoBase *bool                            `json:"merged_into_base"`
	LastCommitUTC  *string                          `json:"last_commit_utc"`
	Locked         bool                             `json:"locked"`
	Prunable       bool                             `json:"prunable"`
	Bare           bool                             `json:"bare"`
	Detached       bool                             `json:"detached"`
	Freshness      workspace.InventoryFactFreshness `json:"freshness"`
	Reasons        []string                         `json:"reasons"`
}
type WorktreeRepository struct {
	ProjectID string   `json:"project_id"`
	RepoID    string   `json:"repo_id"`
	Freshness string   `json:"freshness"`
	Reasons   []string `json:"reasons"`
}
type WorktreeListResult struct {
	Kind          string               `json:"kind"`
	SchemaVersion int                  `json:"schema_version"`
	Workspace     WorkspaceRef         `json:"workspace"`
	AsOf          string               `json:"as_of"`
	Worktrees     []WorktreeItem       `json:"worktrees"`
	Repositories  []WorktreeRepository `json:"repositories"`
	Freshness     string               `json:"freshness"`
	Reasons       []string             `json:"reasons"`
}
type inventoryOwner struct {
	kind, id, ref, common string
	base                  *WorktreeBase
	freshness             string
}

func ListWorktrees(d workspace.Dependencies, s *Snapshot, o WorktreeOptions) (WorktreeListResult, error) {
	out := WorktreeListResult{Kind: WorktreeListKind, SchemaVersion: 1, Workspace: WorkspaceRef{s.Workspace.Root}, AsOf: s.ObservedAtUTC, Worktrees: []WorktreeItem{}, Repositories: []WorktreeRepository{}, Freshness: s.Freshness, Reasons: []string{}}
	if err := validateProjectFilter(s, o.ProjectID); err != nil {
		return out, err
	}
	if o.RepoID != "" {
		id := workspace.RepoID(o.RepoID)
		if err := ValidateTaskFilters(s, workspace.TaskListFilters{RepoID: &id}); err != nil {
			return out, err
		}
	}
	repos := map[workspace.RepoID]workspace.RepoRecord{}
	for _, r := range s.Projects.Repos {
		repos[r.ID] = r
	}
	for _, project := range s.Projects.Projects {
		if o.ProjectID != "" && string(project.ID) != o.ProjectID {
			continue
		}
		for _, id := range project.RepoIDs {
			if o.RepoID != "" && string(id) != o.RepoID {
				continue
			}
			repo, ok := repos[id]
			if !ok {
				out.Repositories = append(out.Repositories, WorktreeRepository{string(project.ID), string(id), "unknown", []string{"registered_repository_missing"}})
				out.Freshness = "unknown"
				continue
			}
			owners := inventoryOwners(s, project.ID, repo)
			bases := map[string]string{}
			for path, owner := range owners {
				if owner.base != nil {
					bases[path] = owner.base.OID
				}
			}
			facts, err := workspace.ReadRegisteredWorktreeInventory(d, repo, bases)
			if err != nil {
				out.Repositories = append(out.Repositories, WorktreeRepository{string(project.ID), string(id), "unknown", []string{"worktree_inventory_unavailable: " + err.Error()}})
				out.Freshness = "unknown"
				continue
			}
			out.Repositories = append(out.Repositories, WorktreeRepository{string(project.ID), string(id), facts.Freshness, append([]string{}, facts.Reasons...)})
			if facts.Freshness != "fresh" {
				out.Freshness = "unknown"
			}
			for _, f := range facts.Worktrees {
				entry := f.Entry
				row := WorktreeItem{ProjectID: string(project.ID), RepoID: string(id), Locator: entry.Locator, Ref: stringValue(entry.Ref), OID: stringValue(entry.OID), Clean: f.Clean, Owner: "none", OwnerFreshness: s.Freshness, Ahead: f.Ahead, Behind: f.Behind, MergedIntoBase: f.MergedIntoBase, LastCommitUTC: f.LastCommitUTC, Locked: entry.Locked, Prunable: entry.Prunable, Bare: entry.Bare, Detached: entry.Detached, Freshness: f.Freshness, Reasons: append([]string{}, f.Reasons...)}
				if owner, ok := owners[entry.Locator]; ok {
					row.Owner, row.OwnerID, row.OwnerFreshness, row.Base = owner.kind, stringValue(owner.id), owner.freshness, owner.base
					if owner.ref != "" && entry.Ref != owner.ref || owner.common != repo.GitCommonDir {
						row.OwnerFreshness = "stale"
						row.Reasons = append(row.Reasons, "registered_worktree_binding_differs")
					}
					if f.Freshness.Identity != "fresh" {
						row.OwnerFreshness = "unknown"
					}
				}
				row.Reasons = sortedUnique(row.Reasons)
				out.Worktrees = append(out.Worktrees, row)
			}
		}
	}
	sort.Slice(out.Worktrees, func(i, j int) bool {
		a, b := out.Worktrees[i], out.Worktrees[j]
		if a.ProjectID != b.ProjectID {
			return a.ProjectID < b.ProjectID
		}
		if a.RepoID != b.RepoID {
			return a.RepoID < b.RepoID
		}
		return a.Locator < b.Locator
	})
	sort.Slice(out.Repositories, func(i, j int) bool {
		a, b := out.Repositories[i], out.Repositories[j]
		if a.ProjectID != b.ProjectID {
			return a.ProjectID < b.ProjectID
		}
		return a.RepoID < b.RepoID
	})
	s.RefreshFreshness(d)
	if s.Freshness != "fresh" {
		out.Freshness = s.Freshness
		for i := range out.Worktrees {
			out.Worktrees[i].OwnerFreshness = "unknown"
		}
	}
	out.Reasons = append(out.Reasons, s.Reasons...)
	return out, nil
}

func inventoryOwners(s *Snapshot, project workspace.ProjectID, repo workspace.RepoRecord) map[string]inventoryOwner {
	owners := map[string]inventoryOwner{repo.Locator: {kind: "project_main", id: string(project), common: repo.GitCommonDir, freshness: s.Freshness}}
	bases := map[workspace.EpicID]*WorktreeBase{}
	for _, epic := range s.Registry.Epics {
		if epic.ProjectID != project {
			continue
		}
		for _, base := range workspace.RegisteredEpicBases(s.Registry, epic) {
			if base.RepoID == repo.ID {
				b := &WorktreeBase{"registered_epic_base", base.Locator, base.Ref, base.OID, base.Revision}
				bases[epic.ID] = b
				owners[base.Locator] = inventoryOwner{"epic", string(epic.ID), base.Ref, base.GitCommonDir, b, s.Freshness}
			}
		}
	}
	for _, task := range s.Registry.Tasks {
		if task.ProjectID != project || task.RepoID != repo.ID {
			continue
		}
		if task.Worktree != nil {
			w := task.Worktree
			owners[w.Locator] = inventoryOwner{"task", string(task.ID), w.Ref, w.GitCommonDir, bases[task.ParentEpicID], s.Freshness}
			continue
		}
		for _, op := range s.Registry.WorktreeOperations {
			if op.TaskID == task.ID {
				owners[op.TargetLocator] = inventoryOwner{"task", string(task.ID), op.SourceRef, op.GitCommonDir, bases[task.ParentEpicID], "unknown"}
			}
		}
	}
	return owners
}
