//go:build unix

package taskrun

import (
	"fmt"
	"os"
	"syscall"
)

func workflowFileIdentity(st os.FileInfo) (uint64, uint64, error) {
	s, ok := st.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, fmt.Errorf("unsupported filesystem identity")
	}
	return uint64(s.Dev), s.Ino, nil
}
