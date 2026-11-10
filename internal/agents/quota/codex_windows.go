//go:build windows

package quota

import "context"

// ReadCodex is not supported on Windows in v1.
func ReadCodex(context.Context, CodexOptions) (*Snapshot, error) {
	return nil, ErrUnsupportedPlatform
}
