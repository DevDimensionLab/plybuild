//go:build !windows

package agentwork

import (
	"os"
	"os/exec"
	"syscall"
)

// Keep a group leader's children in its group so a caller's group cleanup can
// still reach them after Ply is hard-killed. Group leadership proves topology,
// not that such a cleanup owner exists. Without caller cleanup, this branch
// retains the F02 standalone descendant leak on normal timeout (r001-f01).
// CommandContext/WaitDelay bound only the direct child and pipe draining here.
// For nonleaders, isolate each read so its own cancellation can clean up the
// child group without signaling the invoking shell. See the documented scope
// and open Ply-plan backlog; this helper does not establish caller authority.
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
