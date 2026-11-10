package quota

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// DefaultCodexTimeout bounds process start and the whole exchange.
	DefaultCodexTimeout = 15 * time.Second
	// maxCodexLine bounds one newline-delimited message from app-server.
	maxCodexLine = 1 << 20
)

// CodexOptions configures ReadCodex. An explicit Executable is never replaced
// by another installation. Home "" uses $HOME for the default search dirs.
type CodexOptions struct {
	Executable string
	Timeout    time.Duration
	Home       string
}

func locateCodex(explicit, pathEnv string, dirs []string) (string, error) {
	candidates := []string{explicit}
	if explicit == "" {
		candidates = nil
		for _, dir := range append(filepath.SplitList(pathEnv), dirs...) {
			if dir != "" {
				candidates = append(candidates, filepath.Join(dir, "codex"))
			}
		}
	}
	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
			return path, nil
		}
	}
	return "", ErrCLINotFound
}

func codexSearchDirs(home string) []string {
	return []string{"/opt/homebrew/bin", "/usr/local/bin", filepath.Join(home, ".local/bin"),
		filepath.Join(home, ".npm-global/bin"), filepath.Join(home, "bin"), "/usr/bin"}
}

type codexWindow struct {
	UsedPercent        *float64 `json:"usedPercent"`
	WindowDurationMins *int     `json:"windowDurationMins"`
	ResetsAt           *int64   `json:"resetsAt"`
}

type codexLimit struct {
	LimitID              *string      `json:"limitId"`
	LimitName            *string      `json:"limitName"`
	Primary              *codexWindow `json:"primary"`
	Secondary            *codexWindow `json:"secondary"`
	PlanType             *string      `json:"planType"`
	RateLimitReachedType interface{}  `json:"rateLimitReachedType"`
}

type codexRateLimits struct {
	RateLimits          *codexLimit            `json:"rateLimits"`
	RateLimitsByLimitID map[string]*codexLimit `json:"rateLimitsByLimitId"`
}

func (l *codexLimit) windows(pool string) ([]Window, error) {
	var out []Window
	for _, w := range []*codexWindow{l.Primary, l.Secondary} {
		if w == nil {
			continue
		}
		if w.UsedPercent == nil {
			return nil, ErrInvalidResponse
		}
		if w.WindowDurationMins == nil || *w.WindowDurationMins <= 0 {
			continue
		}
		window := Window{ID: WindowID{Minutes: *w.WindowDurationMins}, Pool: pool, UsedPercent: *w.UsedPercent}
		if w.ResetsAt != nil {
			reset := time.Unix(*w.ResetsAt, 0)
			window.ResetsAt = &reset
		}
		out = append(out, window)
	}
	return out, nil
}

// decodeCodexRateLimits normalizes an account/rateLimits/read result.
func decodeCodexRateLimits(raw json.RawMessage, observedAt time.Time) (*Snapshot, error) {
	var result codexRateLimits
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("%w: codex rate limits", ErrInvalidResponse)
	}
	main := result.RateLimits
	if main == nil {
		return nil, ErrNoQuotaReported
	}
	snapshot := &Snapshot{Source: SourceCodex, ObservedAt: observedAt, LimitReached: main.RateLimitReachedType != nil}
	if main.PlanType != nil {
		snapshot.Plan = *main.PlanType
	}
	windows, err := main.windows("")
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(result.RateLimitsByLimitID))
	for key := range result.RateLimitsByLimitID {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		limit := result.RateLimitsByLimitID[key]
		if limit == nil || (main.LimitID != nil && key == *main.LimitID) {
			continue
		}
		pool := key
		if limit.LimitName != nil && *limit.LimitName != "" {
			pool = *limit.LimitName
		}
		poolWindows, err := limit.windows(pool)
		if err != nil {
			return nil, err
		}
		windows = append(windows, poolWindows...)
	}
	snapshot.Windows = windows
	return snapshot, nil
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// classifyRPCError uses the private upstream message only to detect auth.
func classifyRPCError(e *rpcError) error {
	message := strings.ToLower(e.Message)
	for _, word := range []string{"auth", "login", "log in", "sign in", "logged"} {
		if strings.Contains(message, word) {
			return ErrNotSignedIn
		}
	}
	return &UpstreamError{Code: e.Code}
}
