package taskrun

import (
	"strings"

	"github.com/devdimensionlab/plybuild/internal/currentstate"
)

// WorkerOutcome records interpretation of an attempt, independently of its
// diagnostic. Empty is reserved for historical records; unknown is explicit so
// changing display text can never turn a new uncertain effect into a failure.
type WorkerOutcome string

const (
	WorkerOutcomeUnknown           WorkerOutcome = "unknown"
	WorkerOutcomeAttemptsExhausted WorkerOutcome = "attempts_exhausted"
	WorkerOutcomeDeveloperBlocked  WorkerOutcome = "developer_blocked"
	WorkerOutcomeReviewBlocked     WorkerOutcome = "review_blocked"
)

type workerOutcomeError struct {
	outcome WorkerOutcome
	detail  string
}

func (e *workerOutcomeError) Error() string { return e.detail }

// legacyWorkerOutcome is the only compatibility interpretation of the old
// diagnostic format. It does not rewrite records or broaden native policy.
func legacyWorkerOutcome(reason string) WorkerOutcome {
	switch {
	case strings.HasPrefix(reason, "worker reached its explicit limit"):
		return WorkerOutcomeAttemptsExhausted
	case strings.HasPrefix(reason, "developer blocked:"):
		return WorkerOutcomeDeveloperBlocked
	case strings.HasPrefix(reason, "review did not give an unambiguous pass: blocked;"):
		return WorkerOutcomeReviewBlocked
	default:
		return WorkerOutcomeUnknown
	}
}

// QueueReleaseAdvice describes a parked concrete failure. It grants no effect
// authority: WithParkedWorker reacquires ownership and checks the preparation,
// control, bound attempt, process completion and provider evidence at release.
// Stopped failures remain eligible for native recovery; Week only automatically
// releases blocked Runs under its existing schema-1 policy.
// Completed releases retain advice so a caller that lost the reply can reconcile
// the same disposition; the native preparation prevents another release effect.
func (w *WorkerStatus) QueueReleaseAdvice() bool {
	if w == nil || (w.Phase != "blocked" && w.Phase != "stopped") || w.RetryNotBefore != "" || w.ReasonCode == currentstate.ProviderWait || w.ControlAction == "resume" || w.Job == nil || (w.Stage != "developer" && w.Stage != "reviewer" && w.Stage != "pretest") {
		return false
	}
	outcome := w.Outcome
	if outcome == "" {
		outcome = legacyWorkerOutcome(w.Reason)
	}
	switch outcome {
	case WorkerOutcomeAttemptsExhausted, WorkerOutcomeDeveloperBlocked, WorkerOutcomeReviewBlocked:
		return true
	default:
		return false
	}
}
