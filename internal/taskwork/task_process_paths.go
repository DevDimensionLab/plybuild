package taskwork

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// GoalStartKey names the start of one Goal revision of a Task.
func GoalStartKey(task TaskID, spec string, revision int) string {
	return fmt.Sprintf("%s@%s@%d", task, spec, revision)
}

// GoalStartReceiptPath is where workflow goal-start keeps the start receipt of
// one Goal revision of a Task, in the work of the Task's Project.
func GoalStartReceiptPath(work workspace.ProjectWork, task TaskID, spec string, revision int) string {
	return filepath.Join(work.TaskProcess(), "goal-start", string(task), fmt.Sprintf("%s@%d.json", spec, revision))
}

// TaskGoalStartReceiptPath is the one receipt per Task that goal-start kept
// before it kept one per Goal revision. It does not say which revision it
// started, so where it exists it still owns every start of the Task.
func TaskGoalStartReceiptPath(work workspace.ProjectWork, task TaskID) string {
	return filepath.Join(work.TaskProcess(), "goal-start", string(task)+".json")
}

// GoalStartReceiptPaths lists the goal-start receipts of a Task: the one kept
// per Task first, then those per Goal revision with the highest revision last.
func GoalStartReceiptPaths(work workspace.ProjectWork, task TaskID) ([]string, error) {
	var paths []string
	single := TaskGoalStartReceiptPath(work, task)
	if _, err := os.Lstat(single); err == nil {
		paths = append(paths, single)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	dir := strings.TrimSuffix(single, ".json")
	entries, err := os.ReadDir(dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	revisions := map[string]int{}
	for _, entry := range entries {
		name, isReceipt := strings.CutSuffix(entry.Name(), ".json")
		revision, err := strconv.Atoi(name[strings.LastIndex(name, "@")+1:])
		if !isReceipt || err != nil {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		revisions[path] = revision
		paths = append(paths, path)
	}
	sort.SliceStable(paths, func(i, j int) bool { return revisions[paths[i]] < revisions[paths[j]] })
	return paths, nil
}
