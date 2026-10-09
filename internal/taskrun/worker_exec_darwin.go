//go:build darwin

package taskrun

import "time"

func workerOwnedExec(s LaunchSpec, started func(Process) error) (Process, error) {
	return runOwnedExec(s, started, 86400*time.Second)
}
