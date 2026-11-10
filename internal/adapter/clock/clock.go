// Package clock isolates current-time selection from callers.
package clock

import "time"

// Clock returns the current time selected by a caller.
type Clock interface {
	Now() time.Time
}

// Dependencies contains the clock effect used by a caller. Its zero value
// returns the deterministic zero time; production callers must select System
// explicitly.
type Dependencies struct {
	Clock Clock
}

// Now returns the exact time from the configured dependency.
func Now(dependencies Dependencies) time.Time {
	if dependencies.Clock == nil {
		return time.Time{}
	}
	return dependencies.Clock.Now()
}

// System returns the production dependency that reads the system clock.
func System() Dependencies {
	return Dependencies{Clock: systemClock{}}
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}
