//go:build windows

package weekrun

import (
	"context"
	"fmt"
)

func launchDetached(Binding, ...string) error { return fmt.Errorf("week controller requires macOS") }
func (m Manager) Supervise(context.Context, string, string) error {
	return fmt.Errorf("week controller requires macOS")
}
