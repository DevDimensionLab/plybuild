//go:build windows

package weekrun

import (
	"errors"
	"os"
)

var errOwned = errors.New("week controller is already owned")

func lock(string) (*os.File, error) { return nil, errors.New("week service requires macOS") }
func ownerState(string) string      { return "unsupported" }
