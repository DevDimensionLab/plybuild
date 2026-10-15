package taskrun

import (
	"path/filepath"
	"time"
)

// ParentControlGate is the sole dispatch authority for an opted-in week. Only
// the native parent writes it, under control.lock; a child cannot clear it.
type ParentControlGate struct {
	Contract      *ControlContract `json:"control_contract"`
	Deadline      string           `json:"deadline"`
	DesiredAction string           `json:"desired_action"`
	Revision      uint64           `json:"revision"`
	ActiveRun     string           `json:"active_run"`
	Dispatch      string           `json:"dispatch"`
	ResumeAllowed bool             `json:"resume_allowed"`
	StopSequence  string           `json:"stop_sequence"`
}

func ParentControlPath(root, week string) string {
	return filepath.Join(root, ".ply", "week-runs", week, "control-gate.json")
}
func WithParentControlGuard(root, week string, fn func() error) error {
	if !ValidateControlTarget(ControlTarget{"week", week}) || physical(root, false) != nil {
		return controlSafeError("control_invalid_arguments")
	}
	path := filepath.Join(root, ".ply", "week-runs", week, "control.lock")
	if err := privateDir(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := lockFile(path)
	if err != nil {
		return err
	}
	defer unlockFile(f)
	return fn()
}
func WithWorkerParentGuard(root string, c *ControlContract, fn func() error) error {
	if c == nil || c.ParentKind == "standalone" {
		return fn()
	}
	if ValidateControlContract(c) != nil {
		return controlSafeError("control_owner_unknown")
	}
	return WithParentControlGuard(root, *c.WeekID, fn)
}
func ReadParentControl(root string, c *ControlContract) (ParentControlGate, error) {
	var g ParentControlGate
	if c == nil || c.ParentKind != "week" || ValidateControlContract(c) != nil {
		return g, controlSafeError("control_owner_unknown")
	}
	if e := readValue(ParentControlPath(root, *c.WeekID), 1<<20, &g); e != nil {
		return g, controlSafeError("parent_control_unsupported")
	}
	if !equal(g.Contract, c) {
		return g, controlSafeError("parent_control_unsupported")
	}
	return g, nil
}
func CheckParentDispatch(root string, c *ControlContract, now time.Time) error {
	if c == nil || c.ParentKind == "standalone" {
		return nil
	}
	g, e := ReadParentControl(root, c)
	if e != nil {
		return e
	}
	deadline, e := time.Parse(time.RFC3339Nano, g.Deadline)
	if e != nil || !now.Before(deadline) || g.DesiredAction != "continue" {
		return controlSafeError("control_unavailable")
	}
	_, pending, e := ControlHistory(root, ControlTarget{"week", *c.WeekID})
	if e != nil {
		return e
	}
	if pending != nil {
		return controlSafeError("control_in_progress")
	}
	return nil
}
func CheckParentAction(root string, c *ControlContract, action string, now time.Time) error {
	g, e := ReadParentControl(root, c)
	if e != nil {
		return e
	}
	if action == "resume" {
		deadline, e := time.Parse(time.RFC3339Nano, g.Deadline)
		if e != nil || !now.Before(deadline) || g.DesiredAction != "continue" || !g.ResumeAllowed {
			return controlSafeError("control_unavailable")
		}
	} else if g.DesiredAction != action {
		return controlSafeError("control_unavailable")
	}
	return nil
}
