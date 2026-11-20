package workspaceview

import (
	"sort"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const AttentionKind = "WorkspaceAttentionReadback@1"

type AttentionItem struct {
	Subject     string              `json:"subject"`
	SubjectID   string              `json:"subject_id"`
	ProjectID   workspace.ProjectID `json:"project_id"`
	EpicID      taskwork.EpicID     `json:"epic_id"`
	Title       *string             `json:"title"`
	Kind        string              `json:"kind"`
	Actor       string              `json:"actor"`
	Reason      string              `json:"reason"`
	Severity    string              `json:"severity"`
	SinceUTC    *string             `json:"since_utc"`
	Source      string              `json:"source"`
	Freshness   string              `json:"freshness"`
	EvidenceIDs []string            `json:"evidence_ids"`
}

type AttentionList struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Workspace     WorkspaceRef             `json:"workspace"`
	Scope         taskwork.TaskListFilters `json:"scope"`
	Items         []AttentionItem          `json:"items"`
	Freshness     string                   `json:"freshness"`
	SourceBasis   string                   `json:"source_basis"`
	Reasons       []string                 `json:"reasons"`
}

func ReadAttention(d taskwork.Dependencies, filters taskwork.TaskListFilters) (AttentionList, error) {
	if err := filters.Validate(); err != nil {
		return AttentionList{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return AttentionList{}, err
	}
	if err := ValidateTaskFilters(s, filters); err != nil {
		return AttentionList{}, err
	}
	process := readProgressProcess(d, s)
	s.RefreshFreshness(d)
	return BuildAttention(s, filters, process)
}

func BuildAttention(s *Snapshot, filters taskwork.TaskListFilters, process map[taskwork.TaskID]ProcessFacts) (AttentionList, error) {
	tasks, err := BuildTaskList(s, filters, process)
	if err != nil {
		return AttentionList{}, err
	}
	return attentionFromTasks(tasks), nil
}

func attentionFromTasks(tasks TaskList) AttentionList {
	out := AttentionList{AttentionKind, 1, tasks.Workspace, tasks.Scope, []AttentionItem{}, tasks.Freshness, "registered", tasks.Reasons}
	for _, task := range tasks.Tasks {
		if (task.Lifecycle != taskwork.LifecycleActive || task.EpicLifecycle != taskwork.LifecycleActive) && task.Progress.State != "cleanup_pending" || !task.Progress.HasProgress {
			continue
		}
		for _, action := range task.Progress.NextActions {
			if action.Actor != "human" {
				continue
			}
			out.Items = append(out.Items, AttentionItem{"task", string(task.TaskID), task.ProjectID, task.ParentEpicID, task.Title, action.Kind, action.Actor, action.Reason, action.Severity, action.SinceUTC, action.Source, task.Progress.Freshness, action.EvidenceIDs})
		}
	}
	sort.Slice(out.Items, func(i, j int) bool {
		a, b := out.Items[i], out.Items[j]
		if a.Severity != b.Severity {
			return a.Severity == "attention"
		}
		if a.SinceUTC == nil && b.SinceUTC != nil {
			return false
		}
		if a.SinceUTC != nil && b.SinceUTC == nil {
			return true
		}
		if a.SinceUTC != nil && b.SinceUTC != nil {
			at, ae := time.Parse(time.RFC3339Nano, *a.SinceUTC)
			bt, be := time.Parse(time.RFC3339Nano, *b.SinceUTC)
			if ae == nil && be == nil && !at.Equal(bt) {
				return at.Before(bt)
			}
		}
		if a.SubjectID != b.SubjectID {
			return a.SubjectID < b.SubjectID
		}
		if a.Source != b.Source {
			return a.Source < b.Source
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		return a.Reason < b.Reason
	})
	return out
}
