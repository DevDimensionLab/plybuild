//go:build !unix

package taskrun

import (
	"fmt"
	"os"
)

func workflowFileIdentity(st os.FileInfo) (uint64, uint64, error) {
	return 0, 0, fmt.Errorf("Codex process-local trust requires supported Unix filesystem identity")
}
