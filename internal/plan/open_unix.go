//go:build !windows

package plan

import (
	"errors"
	"os"
	"syscall"
)

// openPlanFile opens path for reading without following a symlink and
// without blocking on a FIFO; the caller checks the opened file is regular.
func openPlanFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if errors.Is(err, syscall.ELOOP) {
		return nil, errNotRegular
	}
	return file, err
}
