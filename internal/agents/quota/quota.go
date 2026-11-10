package quota

import (
	"fmt"
	"math"
	"time"
)

// Source names the provider that reported a snapshot.
type Source string

const (
	SourceCodex  Source = "codex"
	SourceClaude Source = "claude"
)

const (
	MinutesFiveHour = 300
	MinutesWeekly   = 10080

	// StaleAfterCodex and StaleAfterClaude are the default staleness policies.
	StaleAfterCodex  = 120 * time.Second
	StaleAfterClaude = 3600 * time.Second
)

// WindowID identifies a window by its reported duration, never by position.
type WindowID struct{ Minutes int }

var (
	FiveHour = WindowID{Minutes: MinutesFiveHour}
	Weekly   = WindowID{Minutes: MinutesWeekly}
)

func (id WindowID) Label() string {
	switch {
	case id == FiveHour:
		return "5-hour"
	case id == Weekly:
		return "Weekly"
	case id.Minutes%1440 == 0:
		return fmt.Sprintf("%d-day", id.Minutes/1440)
	case id.Minutes%60 == 0:
		return fmt.Sprintf("%d-hour", id.Minutes/60)
	}
	return fmt.Sprintf("%d-minute", id.Minutes)
}

// Window is one reported quota window. Pool "" is the shared plan limit;
// any other value is a separate model-specific pool.
type Window struct {
	ID          WindowID
	Pool        string
	UsedPercent float64
	ResetsAt    *time.Time
}

// Remaining is clamp(100 - used, 0, 100).
func (w Window) Remaining() float64 {
	return math.Min(math.Max(100-w.UsedPercent, 0), 100)
}

// ResetPending reports a passed reset; it never means full quota.
func (w Window) ResetPending(now time.Time) bool {
	return w.ResetsAt != nil && !w.ResetsAt.After(now)
}

// Snapshot is one normalized provider observation.
type Snapshot struct {
	Source       Source
	Plan         string
	Windows      []Window
	ObservedAt   time.Time
	LimitReached bool
}

func (s Snapshot) Stale(now time.Time) bool {
	limit := StaleAfterCodex
	if s.Source == SourceClaude {
		limit = StaleAfterClaude
	}
	return now.Sub(s.ObservedAt) > limit
}

type Choice int

const (
	ChooseFiveHour Choice = iota
	ChooseWeekly
	ChooseLowest
)

type SelectionKind int

const (
	SelectedNotReported SelectionKind = iota
	SelectedPending
	SelectedValue
)

// Selection is the outcome of Select; Window is the controlling window.
type Selection struct {
	Kind      SelectionKind
	Remaining float64
	Window    *Window
}

// Select picks among shared windows only and never adds or averages them.
func Select(choice Choice, windows []Window, now time.Time) Selection {
	var candidates []Window
	for _, w := range windows {
		if w.Pool != "" || (choice == ChooseFiveHour && w.ID != FiveHour) || (choice == ChooseWeekly && w.ID != Weekly) {
			continue
		}
		if w.ResetPending(now) {
			return Selection{Kind: SelectedPending, Window: &w}
		}
		candidates = append(candidates, w)
	}
	if len(candidates) == 0 {
		return Selection{Kind: SelectedNotReported}
	}
	lowest := candidates[0]
	for _, w := range candidates[1:] {
		if w.Remaining() < lowest.Remaining() {
			lowest = w
		}
	}
	return Selection{Kind: SelectedValue, Remaining: lowest.Remaining(), Window: &lowest}
}
