package taskjournal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Text(s Snapshot, view string) string {
	var b strings.Builder
	place := "Unbound Task"
	if s.Task.Worktree != nil {
		place = s.Task.Worktree.Locator
	}
	fmt.Fprintf(&b, "Task: %s — %s\nDelivery: %s\n", s.Task.TaskID, s.Task.Title, place)
	for _, p := range [][2]string{{"transport", "Transport"}, {"reported", "Reported"}, {"technical", "Technical"}, {"human", "Human QA"}, {"integration", "Integration"}} {
		a := s.Current.Axes[p[0]]
		fmt.Fprintf(&b, "%s: %s (%s basis)\n", p[1], a.Value, a.BasisStatus)
	}
	fmt.Fprintf(&b, "State: %s coverage; snapshot as of %s\nNext action: %s\nNext actor: %s\n", s.Coverage.State, s.AsOf, s.Current.NextAction.Text, value(s.Current.NextAction.Actor, "Task owner"))
	fmt.Fprintln(&b, "Coverage: available recorded sources only; no complete work/conversation history or active-time total is asserted.")
	for _, r := range s.Coverage.Reasons {
		fmt.Fprintf(&b, "Coverage warning [%s]: %s\n", r.Code, r.Detail)
	}
	fmt.Fprintln(&b, "Next transition authorized: false")
	events := map[string]Event{}
	lanes := map[string]Lane{}
	selected := map[string]bool{}
	for _, e := range s.Events {
		events[e.EventID] = e
	}
	for _, l := range s.Lanes {
		lanes[l.LaneID] = l
	}
	for _, id := range s.Selection.EventIDs {
		selected[id] = true
	}
	ids := s.Selection.EventIDs
	if view == "summary" || view == "" {
		if len(ids) > 8 {
			ids = ids[len(ids)-8:]
		}
		fmt.Fprintln(&b, "\nLatest selected points (up to 8; use --view timeline or --view details):")
	} else {
		fmt.Fprintf(&b, "\nShared UTC timeline (%s order; ordering does not imply causality):\n", s.Selection.Order)
	}
	unknown := []Event{}
	line := func(e Event, stamp *string) {
		l := lanes[e.LaneID]
		fmt.Fprintf(&b, "%s | %s / run %s | %s [%s]", value(stamp, "Unknown time"), value(l.ActorID, "native state"), value(l.RunID, "unknown"), e.Title, e.Type)
		if e.Outcome != nil {
			label := "Reported"
			if e.Origin == "native" {
				label = "Native"
			}
			fmt.Fprintf(&b, "; %s: %s", label, *e.Outcome)
		}
		if len(e.SupersededBy) > 0 {
			fmt.Fprintf(&b, "; superseded by %s", strings.Join(e.SupersededBy, ", "))
		}
		fmt.Fprintf(&b, "\n  occurred: %s; recorded: %s; sequence: %v\n", value(e.OccurredAt, "Unknown time"), value(e.RecordedAt, "Unknown time"), seqText(e.SourceSequence))
		if e.TimeBasis.Uncertainty != nil {
			fmt.Fprintf(&b, "  uncertainty: %s .. %s (%s, %s)\n", e.TimeBasis.Uncertainty.Earliest, e.TimeBasis.Uncertainty.Latest, value(e.TimeBasis.Clock, "unknown clock"), e.TimeBasis.Precision)
		}
	}
	for _, id := range ids {
		e := events[id]
		stamp := e.OccurredAt
		if s.Selection.Order == "recorded" {
			stamp = e.RecordedAt
		}
		if stamp == nil {
			unknown = append(unknown, e)
		} else {
			line(e, stamp)
		}
	}
	if len(unknown) > 0 {
		fmt.Fprintln(&b, "\nUnknown time — outside the selected time axis:")
		for _, e := range unknown {
			line(e, nil)
		}
	}
	if view == "timeline" || view == "details" {
		var earliest, latest time.Time
		for _, p := range s.Steps {
			if !selected[p.StartEventID] && (p.EndEventID == nil || !selected[*p.EndEventID]) {
				continue
			}
			a := events[p.StartEventID]
			if a.OccurredAt != nil {
				t, _ := parseTime(*a.OccurredAt)
				if earliest.IsZero() || t.Before(earliest) {
					earliest = t
				}
				if latest.IsZero() || t.After(latest) {
					latest = t
				}
			}
			if p.EndEventID != nil && events[*p.EndEventID].OccurredAt != nil {
				t, _ := parseTime(*events[*p.EndEventID].OccurredAt)
				if latest.IsZero() || t.After(latest) {
					latest = t
				}
			}
		}
		fmt.Fprintf(&b, "\nIntervals on one UTC axis: %s .. %s\n", axisTime(earliest), axisTime(latest))
		for _, p := range s.Steps {
			if !selected[p.StartEventID] && (p.EndEventID == nil || !selected[*p.EndEventID]) {
				continue
			}
			l := lanes[p.LaneID]
			a := events[p.StartEventID]
			var end *string
			if p.EndEventID != nil {
				end = events[*p.EndEventID].OccurredAt
			}
			bar := strings.Repeat(" ", 48)
			label := map[string]string{"closed": "Closed", "no_end_recorded": "No end recorded", "time_unknown": "Unknown time", "time_conflict": "Time conflict"}[p.State]
			if p.DurationSeconds != nil && !earliest.IsZero() && latest.After(earliest) {
				start, _ := parseTime(*a.OccurredAt)
				finish, _ := parseTime(*end)
				left := int(start.Sub(earliest).Seconds() / latest.Sub(earliest).Seconds() * 47)
				right := int(finish.Sub(earliest).Seconds() / latest.Sub(earliest).Seconds() * 47)
				if left >= 0 && right >= left && right < 48 {
					cells := []byte(bar)
					for i := left; i <= right; i++ {
						cells[i] = '='
					}
					bar = string(cells)
				}
			}
			fmt.Fprintf(&b, "|%s| %s / run %s — %s (%s): %s\n  %s .. %s", bar, value(l.ActorID, "native"), value(l.RunID, "unknown"), p.StepID, p.Kind, label, value(a.OccurredAt, "Unknown time"), value(end, "No end recorded"))
			if p.DurationSeconds != nil {
				fmt.Fprintf(&b, "; %s seconds", strconv.FormatFloat(*p.DurationSeconds, 'f', -1, 64))
			}
			if p.ParentStepID != nil {
				fmt.Fprintf(&b, "; parent %s", *p.ParentStepID)
			}
			if p.Waiting != nil {
				fmt.Fprintf(&b, "; waiting: %s; dependency: %s; next actor: %s", p.Waiting.Reason, p.Waiting.Dependency, p.Waiting.NextActor)
			}
			fmt.Fprintln(&b)
		}
	}
	if view == "details" {
		fmt.Fprintln(&b, "\nSources:")
		for _, src := range s.Sources {
			fmt.Fprintf(&b, "%s | %s | %s | sha256 %s | sequence %v | previous %s\n", src.SourceID, src.Status, src.Locator, value(src.SHA256, "unknown"), seqText(src.Sequence), value(src.PreviousSHA256, "none"))
		}
		fmt.Fprintln(&b, "\nRelations:")
		for _, e := range s.Events {
			for _, r := range e.Relations {
				fmt.Fprintf(&b, "%s %s %s\n", e.EventID, r.Type, r.EventID)
			}
		}
		fmt.Fprintln(&b, "\nObservations (self-reported; proposals are not performed actions):")
		for _, o := range s.Observations {
			fmt.Fprintf(&b, "%s | %s | %s | proposal %s\nFinding: %s\nHypothesis: %s\nProposed action: %s\nOwner: %s; next signal: %s; decision: %s; assessment: %s\n", o.EventID, o.Actor.ID, o.Type, value(o.ProposalEventID, "none"), o.Finding, value(o.Hypothesis, "none"), value(o.Action, "unchanged"), value(o.Owner, "unchanged"), value(o.NextSignal, "unchanged"), value(o.Decision, "none"), value(o.Assessment, "none"))
		}
	}
	return b.String()
}
func value(p *string, fallback string) string {
	if p == nil {
		return fallback
	}
	return *p
}
func seqText(p *int) any {
	if p == nil {
		return "unknown"
	}
	return *p
}
func axisTime(t time.Time) string {
	if t.IsZero() {
		return "Unknown time"
	}
	return t.UTC().Format(time.RFC3339Nano)
}
