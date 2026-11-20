package taskwork

import "github.com/devdimensionlab/plybuild/internal/workspace"

// Dependencies are the effects of Task work, beside the register's.
type Dependencies struct {
	workspace.Dependencies
	SourceReconciliationFault func(string, string) error
	TaskContent               *TaskContentStorage
	TaskContentIDs            TaskContentIDSource
	WorkItems                 WorkItemStore
	WorkGit                   WorkItemGit
	WorkIDs                   WorkItemIDSource
	TaskLifecycleIDs          TaskLifecycleIDSource
	WorkClock                 WorkClock
	IntegrationGit            TaskIntegrationGit
}

// SystemDependencies are the system effects, with the work-item store
// as the register's RepoUses.
func SystemDependencies() Dependencies {
	d := workspace.SystemDependencies()
	workItems := newSystemWorkItemStore()
	d.RepoUses = workItems
	return Dependencies{
		Dependencies:     d,
		TaskContent:      &TaskContentStorage{},
		TaskContentIDs:   cryptoTaskContentIDs{},
		WorkItems:        workItems,
		WorkGit:          newSystemWorkItemGit(d.Files),
		WorkIDs:          cryptoWorkItemIDSource{},
		TaskLifecycleIDs: cryptoTaskLifecycleIDSource{},
		WorkClock:        systemWorkClock{},
		IntegrationGit:   newSystemTaskIntegrationGit(d.Files),
	}
}
