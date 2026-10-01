//go:build !darwin

package taskrun

func systemRunner() ProcessRunner { return unsupportedRunner{} }

type unsupportedRunner struct{}

func (unsupportedRunner) Check() error {
	return failure("task_run_platform_unsupported", 3, "Interactive provider start is supported only in a local macOS terminal.")
}
func (r unsupportedRunner) Run(LaunchSpec, func(Process) error) (Process, error) {
	return Process{State: "not_started"}, r.Check()
}
