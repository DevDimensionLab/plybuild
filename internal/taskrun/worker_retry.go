//go:build !windows

package taskrun

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/currentstate"
)

// WorkerProviderFailure is an observation from a provider error envelope, never
// a successful answer or tool output. The original stream stays bound to the job.
type WorkerProviderFailure struct {
	Kind    string `json:"kind"`
	ResetAt string `json:"reset_at,omitempty"`
}

func workerProviderFailure(raw []byte, now time.Time) *WorkerProviderFailure {
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return nil
	}
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 4096), providerLineLimit+1)
	var result *WorkerProviderFailure
	for scanner.Scan() {
		var event map[string]any
		if json.Unmarshal(scanner.Bytes(), &event) != nil {
			return nil
		}
		typ, _ := event["type"].(string)
		if typ == "turn.completed" {
			return nil
		}
		if typ != "error" && typ != "turn.failed" {
			continue
		}
		fields := event
		if inner, ok := event["error"].(map[string]any); ok {
			fields = inner
		}
		var parts []string
		for _, key := range []string{"code", "type", "codex_error_info", "message"} {
			if value, ok := fields[key].(string); ok {
				parts = append(parts, strings.ToLower(value))
			}
		}
		text := strings.Join(parts, " ")
		kind := ""
		for _, s := range []string{"usage_limit_reached", "usage_limit_exceeded", "rate_limit_exceeded", "rate_limit_reached", "you've hit your usage limit", "usage limit reached", "rate limit exceeded", "rate limit reached"} {
			if strings.Contains(text, s) {
				kind = "usage_limit"
			}
		}
		for _, s := range []string{"authentication_error", "authentication failed", "unauthorized", "invalid_api_key", "token_expired", "token expired", "refresh_token_reused", "not logged in"} {
			if strings.Contains(text, s) {
				kind = "authentication"
			}
		}
		if kind == "" {
			continue
		}
		if result != nil && result.Kind == "authentication" {
			continue
		}
		failure := &WorkerProviderFailure{Kind: kind}
		for _, key := range []string{"resets_at", "reset_at", "retry_at", "retry_after_seconds"} {
			var reset time.Time
			switch value := fields[key].(type) {
			case string:
				reset, _ = time.Parse(time.RFC3339Nano, value)
			case float64:
				if value > 0 && value < 1e12 {
					if key == "retry_after_seconds" {
						reset = now.Add(time.Duration(value) * time.Second)
					} else {
						reset = time.Unix(int64(value), 0)
					}
				}
			}
			if reset.After(now) && reset.Before(now.Add(30*24*time.Hour)) {
				failure.ResetAt = reset.UTC().Format(time.RFC3339Nano)
			}
		}
		// Codex may repeat the same error as a terminal turn.failed event with
		// only a message. That must not discard an already observed reset.
		if result != nil && result.Kind == failure.Kind && result.ResetAt != "" {
			previous, _ := time.Parse(time.RFC3339Nano, result.ResetAt)
			current, _ := time.Parse(time.RFC3339Nano, failure.ResetAt)
			if previous.After(current) {
				failure.ResetAt = result.ResetAt
			}
		}
		result = failure
	}
	if scanner.Err() != nil {
		return nil
	}
	return result
}

func workerWaitForProvider(d Dependencies, s workflowState, result workerAttemptResult) error {
	w := s.Result.Worker
	if w.Job == nil || (w.Stage != "developer" && w.Stage != "reviewer") || result.ProviderFailure == nil || result.Process == nil || !quiescent(*result.Process) || result.Process.State != "exited" || result.Exit == nil || *result.Exit == 0 {
		return fmt.Errorf("provider retry requires a positively completed failed provider attempt")
	}
	raw, err := workflowBound(result.Stdout, providerStreamLimit)
	if err != nil {
		return err
	}
	failure := workerProviderFailure(raw, d.Now())
	if failure == nil || failure.Kind != result.ProviderFailure.Kind {
		return fmt.Errorf("provider failure classification differs from its preserved stream")
	}
	if failure.Kind == "authentication" {
		return fmt.Errorf("authentication_required: existing Codex login failed; inspect %s; no account, model or API fallback was attempted", result.Stdout.Locator)
	}
	limit := s.Request.Worker.MaxProviderRetries
	// Schema 1 preserves its historical omitted/zero default. The explicitly
	// opted-in schema-2 product-question contract treats zero as no extra turn.
	if limit == 0 && s.Request.Worker.SchemaVersion == 1 {
		limit = 20
	}
	if w.ProviderRetries >= limit {
		return fmt.Errorf("provider usage-limit retry bound %d exhausted; preserved attempt %s", limit, w.Job.ID)
	}
	now := d.Now().UTC()
	// Without a provider reset, persist exponential backoff (5m to 1h).
	backoff := 5 * time.Minute
	for i := 0; i < w.ProviderRetries && backoff < time.Hour; i++ {
		backoff *= 2
	}
	if backoff > time.Hour {
		backoff = time.Hour
	}
	next := now.Add(backoff)
	if result.ProviderFailure.ResetAt != "" {
		reset, e := time.Parse(time.RFC3339Nano, result.ProviderFailure.ResetAt)
		if e != nil {
			return e
		}
		next = reset
		if next.Before(now.Add(30 * time.Second)) {
			next = now.Add(30 * time.Second)
		}
	}
	return workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		w := current.Result.Worker
		w.ProviderRetries++
		w.RetryNotBefore = next.Format(time.RFC3339Nano)
		w.LastProviderFailure = result.Stdout.Locator
		w.Phase = "waiting"
		w.Reason = "usage_limit: same Run waits until " + w.RetryNotBefore
		w.ReasonCode = currentstate.ProviderWait
		current.Result.NextAction = WorkflowAction{"worker_wait", w.Reason}
		return nil
	})
}

func workerRetryDue(d Dependencies, s workflowState) (bool, error) {
	if s.Result.Worker.RetryNotBefore == "" {
		return true, nil
	}
	at, err := time.Parse(time.RFC3339Nano, s.Result.Worker.RetryNotBefore)
	if err != nil {
		return false, err
	}
	return !d.Now().Before(at), nil
}

func workerContinueAfterWait(d Dependencies, s workflowState) error {
	due, err := workerRetryDue(d, s)
	if err != nil {
		return err
	}
	if !due {
		return fmt.Errorf("provider reset has not been reached")
	}
	// Re-read the exact positive completion before releasing this reservation.
	r, err := workerRunReserved(d, s)
	if err != nil {
		return err
	}
	if r.ProviderFailure == nil || r.ProviderFailure.Kind != "usage_limit" {
		return fmt.Errorf("waiting job is not a preserved usage-limit attempt")
	}
	return workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		w := current.Result.Worker
		w.Phase, w.RetryNotBefore, w.Reason = "running", "", ""
		w.Job = nil
		return nil
	})
}
