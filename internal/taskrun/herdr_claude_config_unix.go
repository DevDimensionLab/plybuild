//go:build !windows

package taskrun

import (
	"os"

	"golang.org/x/sys/unix"
)

func workflowOpenClaudeConfig(root *os.Root, name string) (*os.File, error) {
	return root.OpenFile(name, os.O_RDONLY|unix.O_NOFOLLOW, 0)
}
