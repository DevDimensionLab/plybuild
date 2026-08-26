// Package clock isolates current-time selection and sleeping from callers.
package clock

import "time"

// Clock returns the current time selected by a caller.
type Clock interface {
	Now() time.Time
}

// Sleeper waits for the exact duration selected by a caller.
type Sleeper interface {
	Sleep(time.Duration)
}

// Dependencies contains the clock effects used by a caller. Its zero value
// returns the deterministic zero time and sleeps as a no-op; production callers
// must select System explicitly.
type Dependencies struct {
	Clock   Clock
	Sleeper Sleeper
}

// Now returns the exact time from the configured dependency.
func Now(dependencies Dependencies) time.Time {
	if dependencies.Clock == nil {
		return time.Time{}
	}
	return dependencies.Clock.Now()
}

// Sleep passes the exact duration to the configured dependency.
func Sleep(dependencies Dependencies, duration time.Duration) {
	if dependencies.Sleeper == nil {
		return
	}
	dependencies.Sleeper.Sleep(duration)
}

// System returns the production dependency that reads and sleeps on the system clock.
func System() Dependencies {
	system := systemClock{}
	return Dependencies{Clock: system, Sleeper: system}
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func (systemClock) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
