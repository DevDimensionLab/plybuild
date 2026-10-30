package taskrun

import (
	"encoding/json"
	"errors"

	"os"
	"path/filepath"

	"testing"
)

// All pane/session/writer observations below are simulated. Native registry,
// separate review sealing, executable acceptance and Git return are real.
func parallelRuntime(t *testing.T, f *deliveryFixture) {
	t.Helper()
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: 70001, CWD: f.Parent}, {PID: 70001, ParentPID: 1, CWD: f.Parent}}, nil
	}
	f.D.StartupProcessTree = func(pid int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	raw, err := os.ReadFile(f.Model)
	if err != nil {
		t.Fatal(err)
	}
	var model map[string]any
	if err = json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	model["caller_shell"] = 70001
	model["review_session_prefix"] = f.R.PreparationID + "-"
	writeAny(t, filepath.Dir(f.Model), filepath.Base(f.Model), model)
}

func parallelOwnerHandoff(err error) bool { var e *Error; return errors.As(err, &e) && e.Exit == 5 }
