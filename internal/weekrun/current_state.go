package weekrun

import "github.com/devdimensionlab/plybuild/internal/currentstate"

// Missing optional origins identify the older initial record, not activity.
// Recovery still compares every original byte of semantic state.
func isInitialState(b Binding, s State) bool {
	want := initialState(b)
	if s.GuidanceOrigin == "" {
		want.GuidanceOrigin = ""
	}
	for i := range want.Entries {
		if i < len(s.Entries) && s.Entries[i].BlockOrigin == "" {
			want.Entries[i].BlockOrigin = ""
		}
	}
	return hash(canonical(s)) == hash(canonical(want))
}

func currentWeek(s State) currentstate.Week {
	return currentstate.Week{Phase: s.Phase, DesiredAction: s.DesiredAction, NextAction: s.NextAction, GuidanceOrigin: s.GuidanceOrigin, UnresolvedDispatch: s.Dispatch != nil, ControllerBlocked: s.Phase == "blocked"}
}

// Only Tick, after successful guarded native reconciliation, can confirm the
// recovery of a Week-owned obstacle. A read cannot attest that controller step.
func projectCurrentRun(s *State, entry *Entry, r Run, reconciliationConfirmed bool) {
	w := currentWeek(*s)
	w.ReconciliationConfirmed = reconciliationConfirmed
	w.OtherEntriesComplete = true
	for i := range s.Entries {
		if &s.Entries[i] != entry && s.Entries[i].State != "completed" {
			w.OtherEntriesComplete = false
		}
	}
	e := currentstate.Entry{State: entry.State, BlockOrigin: entry.BlockOrigin}
	w.GuidanceOrigin, e.BlockOrigin = currentstate.LegacyWeekOrigins(w, e)
	native := currentstate.Native{Phase: r.Phase, Completed: r.Completed, Action: currentstate.Action(r.NextAction), Unavailable: r.ObservationError != nil}
	if r.ObservationError != nil {
		native.Failure = r.ObservationError.Error()
	}
	w, e = currentstate.WeekEntry(w, e, native)
	s.Phase, s.NextAction, s.GuidanceOrigin = w.Phase, w.NextAction, w.GuidanceOrigin
	entry.State, entry.BlockOrigin = e.State, e.BlockOrigin
}

func projectQueueGuidance(s *State) {
	w := currentWeek(*s)
	w.GuidanceOrigin, _ = currentstate.LegacyWeekOrigins(w, currentstate.Entry{})
	w = currentstate.WeekQueue(w)
	s.NextAction, s.GuidanceOrigin = w.NextAction, w.GuidanceOrigin
}
