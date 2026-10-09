//go:build !darwin && !windows

package taskrun

import "fmt"

func workerOwnedExec(LaunchSpec, func(Process) error) (Process, error) {
	return Process{}, fmt.Errorf("headless worker requires the verified macOS process adapter")
}
