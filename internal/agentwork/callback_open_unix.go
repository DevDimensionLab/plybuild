//go:build !windows

package agentwork

import (
	"os"

	"golang.org/x/sys/unix"
)

func openCallbackExecutable(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|unix.O_NONBLOCK|unix.O_NOFOLLOW, 0)
}
