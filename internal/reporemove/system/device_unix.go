//go:build unix

package system

import (
	"os"
	"syscall"
)

// deviceOf returns the st_dev of an Lstat result.
func deviceOf(info os.FileInfo) (uint64, bool) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Dev), true
	}
	return 0, false
}
