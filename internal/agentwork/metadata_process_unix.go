//go:build !windows

package agentwork

import (
	"os"
	"os/exec"
	"syscall"
)

// When a process adapter made Ply its group leader, children must remain in
// that group: even SIGKILL of Ply must not let Herdr escape the adapter's final
// group cleanup. In this mode CommandContext owns the direct child and bounds
// pipe draining; the external adapter owns descendant cleanup.
// For ordinary invocations without a dedicated Ply group, isolate each read
// so cancellation can clean up its descendants without signaling the shell.
func configureMetadataProcess(c *exec.Cmd) func() {
	if syscall.Getpgrp() == os.Getpid() {
		return func() {}
	}
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	kill := func() error {
		if c.Process == nil {
			return os.ErrProcessDone
		}
		err := syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
		if err == syscall.ESRCH {
			return os.ErrProcessDone
		}
		return err
	}
	c.Cancel = kill
	return func() { _ = kill() }
}
