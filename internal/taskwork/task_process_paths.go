package taskwork

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"path/filepath"
)

// GoalStartReceiptPath is where workflow goal-start keeps the start receipt of
// a Task, in the work of the Task's Project.
func GoalStartReceiptPath(work workspace.ProjectWork, task TaskID) string {
	return filepath.Join(work.TaskProcess(), "goal-start", string(task)+".json")
}
