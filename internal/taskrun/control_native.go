package taskrun

import "github.com/devdimensionlab/plybuild/internal/workspace"

// These existing native snapshot locks are acquired after TaskRun. They remain
// held through publication, so independent native policy/queue mutations cannot
// pass between revalidation and the durable control reservation.
func WithNativeControlSnapshot(d Dependencies, root string, fn func(string) error) error {
	return workspace.WithTaskSpecRegistrationSnapshot(d.Workspace, root, func(s *workspace.TaskSpecSession) error {
		return fn(digest(struct {
			Registry workspace.WorkItemRegistry
			Projects workspace.ProjectSnapshot
		}{s.Registry, s.Projects}))
	})
}

type controlPublicationGuard interface {
	ControlPublicationGuard(string, string, func() error) error
}
