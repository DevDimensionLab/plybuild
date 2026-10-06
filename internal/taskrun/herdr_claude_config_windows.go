//go:build windows

package taskrun

import (
	"errors"
	"os"
)

func workflowOpenClaudeConfig(root *os.Root, name string) (*os.File, error) {
	return nil, errors.New("Claude configuration writes are unsupported on Windows")
}
