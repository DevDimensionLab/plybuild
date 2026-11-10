package quota

import (
	"errors"
	"fmt"
)

var (
	ErrNotConfigured       = errors.New("claude-statusline-relay is not configured")
	ErrCLINotFound         = errors.New("codex CLI not found")
	ErrNotSignedIn         = errors.New("codex is not signed in")
	ErrNoQuotaReported     = errors.New("no quota reported")
	ErrInvalidResponse     = errors.New("invalid quota data")
	ErrTimeout             = errors.New("quota read timed out")
	ErrUnsupportedPlatform = errors.New("quota sources are not supported on this platform")
)

// UpstreamError is a Codex RPC error. It never carries the upstream message.
type UpstreamError struct{ Code int }

func (e *UpstreamError) Error() string {
	return fmt.Sprintf("codex reported an error (rpc code %d)", e.Code)
}
