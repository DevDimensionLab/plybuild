package file

import (
	"fmt"
	"path/filepath"
)

// Path safe path function for windows compatibility
func Path(format string, a ...interface{}) string {
	return filepath.FromSlash(fmt.Sprintf(format, a...))
}
