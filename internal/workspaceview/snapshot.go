// Package workspaceview combines read-only domain projections for workspace
// clients. Domain stores remain the owners of validation and recorded facts.
package workspaceview

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type WorkspaceRef struct {
	Root string `json:"root"`
}

type Snapshot struct {
	Workspace     workspace.WorkspaceObservation
	Projects      workspace.ProjectSnapshot
	Registry      workspace.WorkItemRegistry
	Titles        map[workspace.TaskID]workspace.TaskListTitle
	Closeouts     map[workspace.TaskID]workspace.TaskCloseoutReceipt
	Lifecycle     workspace.WorkItemLifecycleSnapshot
	Freshness     string
	Reasons       []string
	ObservedAtUTC string
	basis         workspace.ReadViewBasisSnapshot
	process       *processHistory
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

func LoadSnapshot(d workspace.Dependencies) (*Snapshot, error) {
	basis, err := workspace.ReadViewBasis(d)
	if err != nil {
		return nil, err
	}
	lifecycle, err := workspace.ReadWorkItemLifecycle(d, basis.Workspace.Root)
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
	for _, event := range lifecycle.Events {
		if !known[event.SubjectKind+"/"+event.SubjectID] {
			return nil, &workspace.WorkItemError{Class: workspace.ErrorWorkStoreConflict, Detail: "lifecycle history refers to an unregistered subject"}
		}
	}
	now := time.Now()
	if d.WorkClock != nil {
		now = d.WorkClock.Now()
	}
	s := &Snapshot{Workspace: basis.Workspace, Projects: basis.Projects, Registry: basis.Registry, Titles: basis.Titles, Lifecycle: lifecycle, Freshness: "fresh", Reasons: []string{}, ObservedAtUTC: now.UTC().Format(time.RFC3339Nano), basis: basis}
	s.Closeouts = map[workspace.TaskID]workspace.TaskCloseoutReceipt{}
	closeouts, err := workspace.ReadTaskCloseoutsAt(basis.Workspace.Root)
	if err != nil {
		return nil, err
	}
	for _, receipt := range closeouts {
		s.Closeouts[receipt.Plan.TaskID] = receipt
	}
	s.RefreshFreshness(d)
	return s, nil
}

// RefreshFreshness checks captured register bytes after optional batch reads.
// It does not reload typed registries or assert a cross-register transaction.
func (s *Snapshot) RefreshFreshness(d workspace.Dependencies) {
	fresh, reasons := workspace.CheckReadViewBasis(d, s.basis)
	if fresh != "fresh" {
		s.Freshness = "unknown"
		s.Reasons = append(s.Reasons, reasons...)
	}
	if s.Lifecycle.Locator != "" {
		info, inspectErr := d.Files.Lstat(s.Lifecycle.Locator)
		if inspectErr != nil && !os.IsNotExist(inspectErr) || inspectErr == nil && !info.Mode().IsRegular() {
			s.Freshness = "unknown"
			s.Reasons = sortedUnique(append(s.Reasons, "lifecycle_source_changed_or_unavailable"))
			return
		}
		b, err := d.Files.ReadFile(s.Lifecycle.Locator)
		if !s.Lifecycle.Exists && os.IsNotExist(err) {
			// Absence is the valid, read-only default.
		} else if err != nil || !s.Lifecycle.Exists || fmt.Sprintf("sha256:%x", sha256.Sum256(b)) != s.Lifecycle.SHA256 {
			s.Freshness = "unknown"
			s.Reasons = append(s.Reasons, "lifecycle_source_changed_or_unavailable")
		}
	}
	s.Reasons = sortedUnique(s.Reasons)
}

func ValidateTaskFilters(s *Snapshot, filters workspace.TaskListFilters) error {
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
		found := false
		for _, e := range s.Registry.Epics {
			found = found || e.ID == *filters.EpicID
		}
		if !found {
			return notFound("Epic", string(*filters.EpicID))
		}
	}
	return nil
}

func MatchesTaskFilters(task workspace.TaskRecord, filters workspace.TaskListFilters) bool {
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
