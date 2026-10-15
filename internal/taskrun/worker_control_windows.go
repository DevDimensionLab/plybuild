//go:build windows

package taskrun

type workerControlBackend struct{ d Dependencies }

func WorkerControls(d Dependencies) ControlService { return ControlService{d, workerControlBackend{d}} }
func ShowControl(d Dependencies, root, id string) (ControlBasis, error) {
	return WorkerControls(d).Show(root, id)
}
func ControlTargetForRun(root, id string) (ControlTarget, error) {
	return ControlTarget{"run", id}, nil
}
func (b workerControlBackend) ControlBasis(root, id string) (ControlBasis, error) {
	o := NewControlBasis(root, ControlTarget{"run", id}, controlNow(b.d))
	o.Unavailable("runtime_unsupported")
	o.StateRevision = digest(o.Scope)
	return o, nil
}
func (b workerControlBackend) ControlGuard(root, id string, fn func() error) error {
	return controlSafeError("control_unsupported")
}
func (b workerControlBackend) ControlWake(root, id string) error {
	return controlSafeError("control_unsupported")
}
