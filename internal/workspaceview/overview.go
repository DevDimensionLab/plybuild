package workspaceview

import "github.com/devdimensionlab/plybuild/internal/workspace"

const OverviewKind = "WorkspaceOverviewReadback@1"

type Overview struct {
	Kind          string          `json:"kind"`
	SchemaVersion int             `json:"schema_version"`
	Workspace     WorkspaceRef    `json:"workspace"`
	Scope         EpicFilters     `json:"scope"`
	Tasks         []TaskRow       `json:"tasks"`
	Epics         []EpicRow       `json:"epics"`
	Attention     []AttentionItem `json:"attention"`
	SourceBasis   string          `json:"source_basis"`
	Freshness     string          `json:"freshness"`
	Reasons       []string        `json:"reasons"`
}

// ReadOverview reads sources once and derives all three Home collections from
// that captured basis. It does not claim an atomic transaction across stores.
func ReadOverview(d workspace.Dependencies, filters EpicFilters) (Overview, error) {
	taskFilters := workspace.TaskListFilters{ProjectID: filters.ProjectID}
	if err := taskFilters.Validate(); err != nil {
		return Overview{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return Overview{}, err
	}
	if err := ValidateTaskFilters(s, taskFilters); err != nil {
		return Overview{}, err
	}
	process := readProgressProcess(d, s)
	s.RefreshFreshness(d)
	tasks, err := BuildTaskList(s, taskFilters, process)
	if err != nil {
		return Overview{}, err
	}
	epics := epicsFromTasks(s, filters, tasks)
	attention := attentionFromTasks(tasks)
	return Overview{OverviewKind, 1, tasks.Workspace, filters, tasks.Tasks, epics.Epics, attention.Items, "registered", tasks.Freshness, tasks.Reasons}, nil
}
