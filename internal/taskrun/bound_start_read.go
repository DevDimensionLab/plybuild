package taskrun

import (
	"os"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ReadBoundStartRun is a historical read. It neither qualifies the installed
// runtime nor probes/restarts the provider. Missing receipts remain unknown to
// its caller, which owns the earlier immutable bound-start intent.
func ReadBoundStartRun(d Dependencies, root, id string) (*InventoryRun, error) {
	observed, err := workspace.ObserveRoot(d.Workspace, root)
	if err != nil {
		return nil, err
	}
	if observed.Root != root || !workflowRunIDPattern.MatchString(id) {
		return nil, conflict("invalid bound start workspace or Run identity")
	}
	if _, err = os.Lstat(workflowIndex(root, id)); err != nil {
		return nil, err
	}
	row := readHerdrInventoryRun(d, root, id)
	return &row, nil
}
