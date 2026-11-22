package workspace

import "fmt"

// RepoUses tells the register which work items of a Project use its
// repositories. The work-item store implements it.
type RepoUses interface {
	// RepoUses reads the registrations of work: every Epic binding, Task,
	// worktree operation and Epic create operation, in that order. A registry
	// the store rejects is an error.
	RepoUses(ProjectWork) ([]RepoUse, error)
}

// RepoUseKind is the kind of work item that uses a repository, as messages
// name it.
type RepoUseKind string

const (
	RepoUseEpic              RepoUseKind = "Epic"
	RepoUseTask              RepoUseKind = "Task"
	RepoUseWorktreeOperation RepoUseKind = "worktree operation"
	RepoUseEpicCreate        RepoUseKind = "Epic create operation"
)

// RepoUse is one use of a repository by a work item.
type RepoUse struct {
	Kind         RepoUseKind
	ID           string
	RepoID       RepoID
	GitCommonDir string // recorded by an Epic binding or a Task; "" when none
	Worktree     string // the worktree locator of an Epic binding
}

// binds reports whether the use binds the repository; otherwise it only names it.
func (u RepoUse) binds() bool { return u.Kind == RepoUseEpic || u.Kind == RepoUseTask }

// repoRegistration names the first use of repo in uses, read from path; ""
// when none uses it.
func repoRegistration(uses []RepoUse, repo, path string) string {
	for _, use := range uses {
		if string(use.RepoID) != repo {
			continue
		}
		if use.binds() {
			return fmt.Sprintf("%s %s in %s binds it", use.Kind, use.ID, path)
		}
		return fmt.Sprintf("%s %s in %s names it", use.Kind, use.ID, path)
	}
	return ""
}
