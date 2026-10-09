//go:build windows

package taskrun

import "fmt"

func workerUnsupported() error {
	return fmt.Errorf("headless worker requires the verified macOS sandbox adapter")
}
func ShowWorker(Dependencies, string, string) (WorkflowRun, error) {
	return WorkflowRun{}, workerUnsupported()
}
func StartWorker(Dependencies, string, string, bool) (WorkflowRun, error) {
	return WorkflowRun{}, workerUnsupported()
}
func ControlWorker(Dependencies, string, string, string) (WorkflowRun, error) {
	return WorkflowRun{}, workerUnsupported()
}
func ServeWorker(Dependencies, string, string) error                    { return workerUnsupported() }
func SuperviseWorker(Dependencies, string, string) error                { return workerUnsupported() }
func RunWorkerAttempt(Dependencies, string, string, string) error       { return workerUnsupported() }
func WithParkedWorker(Dependencies, string, string, func() error) error { return workerUnsupported() }
func workerVerifierPolicy(Dependencies, workflowState, string) (string, error) {
	return "", workerUnsupported()
}
func workerValidateReview(s workflowState, _ FileBinding, _, _ string) error {
	if headlessRun(s.Request) {
		return workerUnsupported()
	}
	return nil
}
