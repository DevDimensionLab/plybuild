package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// newWorkspaceTaskCleanBuildCommand empties .build/ in a Task's registered
// worktree, so an agent cleans up its temporary files with a Ply command.
func newWorkspaceTaskCleanBuildCommand(d taskwork.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:     "clean-build <task-id>",
		Short:   "Empty .build/ in the Task's worktree",
		Long:    "Remove everything inside .build/ of the Task's registered worktree and keep the directory. Nothing outside .build/ is touched and symbolic links are removed, not followed. A worktree without .build/ is left as it is.",
		Example: "  ply workspace task clean-build code-map-card-name",
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := taskwork.ParseTaskID(args[0])
			if err != nil {
				return err
			}
			task, err := taskwork.ShowTask(d, id)
			if err != nil {
				return err
			}
			if task.Task.Worktree == nil || task.Task.WorktreeState != taskwork.WorkItemReady {
				return workspace.WorkInvalidArguments(fmt.Sprintf("Task %s has no worktree", id))
			}
			n, err := cleanBuild(task.Task.Worktree.Locator)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(c.OutOrStdout(), "Removed %d entries from %s\n", n, filepath.Join(task.Task.Worktree.Locator, ".build"))
			return err
		},
	}
}

// cleanBuild removes the entries of worktree/.build and returns how many.
func cleanBuild(worktree string) (int, error) {
	dir := filepath.Join(worktree, ".build")
	info, err := os.Lstat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("%s is not a directory", dir)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for i, entry := range entries {
		if err = os.RemoveAll(filepath.Join(dir, entry.Name())); err != nil {
			return i, err
		}
	}
	return len(entries), nil
}
