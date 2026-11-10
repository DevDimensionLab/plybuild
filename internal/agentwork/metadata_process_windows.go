//go:build windows

package agentwork

import "os/exec"

// CommandContext and WaitDelay bound the direct child and pipe lifetime.
func configureMetadataProcess(c *exec.Cmd) func() { return func() {} }
