//go:build windows

package backlog

import (
	"errors"
	"os"
)

// Writes fail closed on Windows; the package depends only on golang.org/x/sys/unix.
func lockFile(*os.File) error   { return errors.New("locking is unsupported on Windows") }
func unlockFile(*os.File) error { return nil }
