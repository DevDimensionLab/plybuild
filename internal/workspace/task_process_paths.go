package workspace

import "path/filepath"

// GoalStartReceiptPath is where workflow goal-start keeps the start receipt of a Task.
func GoalStartReceiptPath(root string, task TaskID) string {
	return filepath.Join(root, MarkerDirectory, "task-process", "goal-start", string(task)+".json")
}
