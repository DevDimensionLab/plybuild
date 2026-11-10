//go:build !unix

package system

import "os"

// deviceOf reports no device where the platform has no st_dev; the wrapper
// check then fails closed.
func deviceOf(os.FileInfo) (uint64, bool) { return 0, false }
