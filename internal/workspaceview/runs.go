package workspaceview

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const RunListKind = "WorkspaceRunListReadback@1"

type RunOptions struct {
	ActiveOnly bool
	ProjectID  string
}
type RunItem struct {
	taskrun.InventoryRun
	ProjectID *string `json:"project_id"`
	EpicID    *string `json:"epic_id"`
}
type RunListResult struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	Workspace     WorkspaceRef `json:"workspace"`
	AsOf          string       `json:"as_of"`
	ActiveOnly    bool         `json:"active_only"`
	ActiveMeaning string       `json:"active_meaning"`
	Runs          []RunItem    `json:"runs"`
	Freshness     string       `json:"freshness"`
	Reasons       []string     `json:"reasons"`
}

func ListRuns(d workspace.Dependencies, s *Snapshot, o RunOptions) (RunListResult, error) {
	out := RunListResult{Kind: RunListKind, SchemaVersion: 1, Workspace: WorkspaceRef{s.Workspace.Root}, AsOf: s.ObservedAtUTC, ActiveOnly: o.ActiveOnly, ActiveMeaning: "unresolved_including_unknown", Runs: []RunItem{}, Freshness: s.Freshness, Reasons: []string{}}
	if err := validateProjectFilter(s, o.ProjectID); err != nil {
		return out, err
	}
	// Run-list does not need to render every Task journal. Reuse an existing
	// process batch if present; otherwise read only the run domains.
	var runs taskrun.Inventory
	var err error
	if s.process != nil {
		runs = s.process.runs
	} else {
		runs, err = taskrun.ReadInventory(taskrun.SystemDependencies(d), s.Workspace.Root)
	}
	if err != nil {
		return out, err
	}
	tasks := map[string]workspace.TaskRecord{}
	for _, t := range s.Registry.Tasks {
		tasks[string(t.ID)] = t
	}
	for _, run := range runs.Runs {
		if o.ActiveOnly && !run.Unresolved {
			continue
		}
		row := RunItem{InventoryRun: run}
		if run.TaskID != nil {
			if t, ok := tasks[*run.TaskID]; ok {
				row.ProjectID = stringValue(string(t.ProjectID))
				row.EpicID = stringValue(string(t.ParentEpicID))
			} else {
				row.Freshness = "unknown"
				row.Reasons = append(row.Reasons, taskrun.Reason{Code: "task_unregistered", Detail: "The preserved run Task is not in the captured registry"})
			}
		}
		if o.ProjectID != "" && (row.ProjectID == nil || *row.ProjectID != o.ProjectID) {
			continue
		}
		out.Runs = append(out.Runs, row)
		if row.Freshness != "fresh" {
			out.Freshness = "unknown"
		}
	}
	for _, r := range runs.Reasons {
		out.Reasons = append(out.Reasons, r.Code)
		out.Freshness = "unknown"
	}
	s.RefreshFreshness(d)
	if s.Freshness != "fresh" {
		out.Freshness = s.Freshness
	}
	out.Reasons = sortedUnique(append(out.Reasons, s.Reasons...))
	return out, nil
}
