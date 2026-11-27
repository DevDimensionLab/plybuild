// Package workspaceview combines read-only domain projections for workspace
// clients. Domain stores remain the owners of validation and recorded facts.
package workspaceview

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workitem"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type WorkspaceRef struct {
	Root string `json:"root"`
}

type Snapshot struct {
	Workspace     workspace.WorkspaceObservation
	Projects      workspace.ProjectSnapshot
	Registry      taskwork.WorkItemRegistry
	Titles        map[taskwork.TaskID]taskwork.TaskListTitle
	Lifecycle     taskwork.WorkItemLifecycles
	Freshness     string
	Reasons       []string
	ObservedAtUTC string
	basis         taskwork.ReadViewBasisSnapshot
	process       *processHistory
	liveProcess   map[taskwork.TaskID]ProcessFacts
}

// ProcessFacts is a batch adapter's contribution to the registered progress
// view. Reported work never supplies native technical or human QA outcomes.
type ProcessFacts struct {
	HasProgress     bool
	LastActivityUTC *string
	Waiting         []ProcessWaiting
	Freshness       string
	Reasons         []string
}

type ProcessWaiting struct {
	Kind, Reason, Actor, EventID string
	SinceUTC                     *string
}

func LoadSnapshot(d taskwork.Dependencies) (*Snapshot, error) {
	basis, err := taskwork.ReadViewBasis(d)
	if err != nil {
		return nil, err
	}
	lifecycle, err := taskwork.ReadWorkItemLifecycle(d, basis.Workspace.Root)
	if err != nil {
		return nil, err
	}
	// The lifecycle reader intentionally does not reread the work-item registry.
	// Validate its subject references against this request's existing snapshot.
	known := map[string]bool{}
	for _, task := range basis.Registry.Tasks {
		known["task/"+string(task.ID)] = true
	}
	for _, epic := range basis.Registry.Epics {
		known["epic/"+string(epic.ID)] = true
	}
	for _, event := range lifecycle.Events() {
		if !known[event.SubjectKind+"/"+event.SubjectID] {
			return nil, &workspace.WorkItemError{Class: workspace.ErrorWorkStoreConflict, Detail: "lifecycle history refers to an unregistered subject"}
		}
	}
	now := time.Now()
	if d.WorkClock != nil {
		now = d.WorkClock.Now()
	}
	s := &Snapshot{Workspace: basis.Workspace, Projects: basis.Projects, Registry: basis.Registry, Titles: basis.Titles, Lifecycle: lifecycle, Freshness: "fresh", Reasons: []string{}, ObservedAtUTC: now.UTC().Format(time.RFC3339Nano), basis: basis}
	s.RefreshFreshness(d)
	return s, nil
}

// RefreshFreshness checks captured register bytes after optional batch reads.
// It does not reload typed registries or assert a cross-register transaction.
func (s *Snapshot) RefreshFreshness(d taskwork.Dependencies) {
	fresh, reasons := taskwork.CheckReadViewBasis(d, s.basis)
	if fresh != "fresh" {
		s.Freshness = "unknown"
		s.Reasons = append(s.Reasons, reasons...)
	}
	for _, lifecycle := range s.Lifecycle {
		if !lifecycleFresh(d, lifecycle) {
			s.Freshness = "unknown"
			s.Reasons = append(s.Reasons, "lifecycle_source_changed_or_unavailable")
			break
		}
	}
	s.Reasons = sortedUnique(s.Reasons)
}

// lifecycleFresh reports whether a Project's lifecycle file still has the
// bytes, or the absence, that the snapshot read.
func lifecycleFresh(d taskwork.Dependencies, l taskwork.WorkItemLifecycleSnapshot) bool {
	info, inspectErr := d.Files.Lstat(l.Locator)
	if inspectErr != nil && !os.IsNotExist(inspectErr) || inspectErr == nil && !info.Mode().IsRegular() {
		return false
	}
	b, err := d.Files.ReadFile(l.Locator)
	if !l.Exists && os.IsNotExist(err) {
		// Absence is the valid, read-only default.
		return true
	}
	return err == nil && l.Exists && fmt.Sprintf("sha256:%x", sha256.Sum256(b)) == l.SHA256
}

func ValidateTaskFilters(s *Snapshot, filters taskwork.TaskListFilters) error {
	if err := filters.Validate(); err != nil {
		return err
	}
	if filters.ProjectID != nil {
		found := false
		for _, p := range s.Projects.Projects {
			found = found || p.ID == *filters.ProjectID
		}
		if !found {
			return notFound("project", string(*filters.ProjectID))
		}
	}
	if filters.RepoID != nil {
		found := false
		for _, r := range s.Projects.Repos {
			found = found || r.ID == *filters.RepoID
		}
		if !found {
			return notFound("repository", string(*filters.RepoID))
		}
	}
	if filters.EpicID != nil {
		if epic, _ := workitem.FindEpic(s.Registry.Epics, *filters.EpicID); epic == nil {
			return notFound("Epic", string(*filters.EpicID))
		}
	}
	return nil
}

func MatchesTaskFilters(task taskwork.TaskRecord, filters taskwork.TaskListFilters) bool {
	return (filters.ProjectID == nil || task.ProjectID == *filters.ProjectID) && (filters.RepoID == nil || task.RepoID == *filters.RepoID) && (filters.EpicID == nil || task.ParentEpicID == *filters.EpicID)
}

func notFound(kind, id string) error {
	return &workspace.WorkItemError{Class: workspace.ErrorWorkNotFound, Detail: fmt.Sprintf("%s %s is not registered", kind, id)}
}

func sortedUnique(values []string) []string {
	result := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

func latestTime(previous *string, candidate *string) *string {
	if candidate == nil {
		return previous
	}
	t, err := time.Parse(time.RFC3339Nano, *candidate)
	if err != nil {
		return previous
	}
	if previous != nil {
		old, err := time.Parse(time.RFC3339Nano, *previous)
		if err == nil && !t.After(old) {
			return previous
		}
	}
	normalized := t.UTC().Format(time.RFC3339Nano)
	return &normalized
}

func stringValue(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
