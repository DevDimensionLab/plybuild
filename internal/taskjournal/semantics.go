package taskjournal

import (
	"encoding/json"
	"sort"
	"strings"
)

// Validate immutable relationships independently of mutable artifacts. A source
// may later disappear, but malformed historical steps must still block append.
func validateHistory(prior []stored, ev *EventInput, ob *ObservationInput) error {
	starts := map[string]EventInput{}
	finished := map[string]bool{}
	events := map[string]bool{}
	proposals := map[string]bool{}
	decisions := map[string]map[string]bool{}
	applied := map[string]bool{}
	for _, r := range prior {
		events[r.Record.EventID] = true
		if e := r.Event; e != nil && e.Step != nil {
			if e.Type == "step_started" {
				starts[e.Step.ID] = *e
			}
			if e.Type == "step_finished" {
				finished[e.Step.ID] = true
			}
		}
		if o := r.Observation; o != nil {
			if o.Type == "proposal" {
				proposals[r.Record.EventID] = true
			}
			if o.ProposalEventID != nil {
				id := *o.ProposalEventID
				if o.Decision != nil {
					if decisions[id] == nil {
						decisions[id] = map[string]bool{}
					}
					decisions[id][*o.Decision] = true
				}
				if o.Type == "applied" {
					applied[id] = true
				}
			}
		}
	}
	exists := func(id string) bool {
		return events[id] || strings.HasPrefix(id, "native_") && hexHash.MatchString(strings.TrimPrefix(id, "native_"))
	}
	if ev != nil {
		for _, r := range ev.Relations {
			if !exists(r.EventID) || r.Type == "supersedes" && !events[r.EventID] {
				return conflict("invalid historical relation")
			}
		}
		if p := ev.Step; p != nil {
			if p.ParentStepID != nil {
				if _, ok := starts[*p.ParentStepID]; !ok || *p.ParentStepID == p.ID {
					return conflict("invalid historical parent")
				}
			}
			start, ok := starts[p.ID]
			if ev.Type == "step_started" && ok {
				return conflict("duplicate historical step start")
			}
			if ev.Type == "step_finished" || ev.Type == "note" {
				if !ok || !equal(start.Step, p) {
					return conflict("historical step has no matching start")
				}
				if ev.Type == "step_finished" && (finished[p.ID] || !equal(start.Actor, ev.Actor) || !equal(start.RunBinding, ev.RunBinding)) {
					return conflict("invalid historical step finish")
				}
			}
		}
	} else {
		for _, id := range ob.Targets.EventIDs {
			if !exists(id) {
				return conflict("unknown historical observation target")
			}
		}
		for _, id := range ob.Targets.StepIDs {
			if _, ok := starts[id]; !ok {
				return conflict("unknown historical observation step")
			}
		}
		if ob.ProposalEventID != nil {
			id := *ob.ProposalEventID
			if !proposals[id] {
				return conflict("unknown historical proposal")
			}
			if ob.Type == "applied" && (!decisions[id]["accepted"] || decisions[id]["rejected"]) {
				return conflict("historical applied without unconflicted acceptance")
			}
			if ob.Type == "assessment" && !applied[id] {
				return conflict("historical assessment without applied")
			}
		}
	}
	return nil
}

func validateSnapshotSemantics(s Snapshot) error {
	events := map[string]Event{}
	steps := map[string]Step{}
	observations := map[string]Observation{}
	lanes := map[string]Lane{}
	for _, l := range s.Lanes {
		lanes[l.LaneID] = l
	}
	for _, p := range s.Steps {
		steps[p.StepID] = p
	}
	for _, o := range s.Observations {
		if _, ok := observations[o.EventID]; ok {
			return invalid("duplicate observation")
		}
		observations[o.EventID] = o
	}
	edges := map[string][]string{}
	parents := map[string][]string{}
	superseders := map[string][]string{}
	for _, e := range s.Events {
		events[e.EventID] = e
		var actor *string
		if e.Actor != nil {
			actor = &e.Actor.ID
		}
		l := lanes[e.LaneID]
		if !equal(actor, l.ActorID) || !equal(e.RunID, l.RunID) {
			return invalid("event lane differs from actor/run")
		}
		for _, r := range e.Relations {
			if !one(r.Type, "responds_to", "corrects", "verifies", "supersedes") {
				return invalid("unknown relation type")
			}
			edges[e.EventID] = append(edges[e.EventID], r.EventID)
			if r.Type == "supersedes" {
				superseders[r.EventID] = append(superseders[r.EventID], e.EventID)
			}
		}
		if e.Origin == "contribution" {
			kind := "event"
			if strings.HasPrefix(e.Type, "observation_") {
				kind = "observation"
			}
			c, input, ob, _, err := parseInput(e.Data, kind, s.Task.TaskID)
			if err != nil {
				return err
			}
			var run *string
			if c.RunBinding != nil {
				run = &c.RunBinding.RunID
			}
			if e.EventID != eventID(s.Workspace.ID, s.Task.TaskID, c.PublicationKey) || !equal(e.Actor, &c.Actor) || !equal(e.ActivityID, c.ActivityID) || !equal(e.RunID, run) || !equal(e.OriginalOccurredAt, c.OccurredAt) || !equal(e.TimeBasis, c.TimeBasis) || e.RecordedAt == nil || e.SourceSequence == nil || *e.SourceSequence < 1 {
				return invalid("contribution projection differs from input")
			}
			if input != nil {
				var step *string
				if input.Step != nil {
					step = &input.Step.ID
				}
				if e.Type != input.Type || e.Title != input.Title || !equal(e.StepID, step) || !equal(e.Outcome, input.Outcome) || !equal(e.Candidate, input.Candidate) || !equal(e.Relations, input.Relations) {
					return invalid("event projection differs from input")
				}
			} else {
				o, ok := observations[e.EventID]
				if !ok || !equal(o.ObservationInput, *ob) || e.Type != "observation_"+ob.Type {
					return invalid("observation projection differs")
				}
				if err = validateObservationSnapshot(o, s); err != nil {
					return err
				}
				delete(observations, e.EventID)
			}
		} else if err := validateNativeData(e); err != nil {
			return err
		}
	}
	if len(observations) != 0 {
		return invalid("observation does not refer to an observation event")
	}
	for _, e := range s.Events {
		actual := append([]string{}, e.SupersededBy...)
		want := append([]string{}, superseders[e.EventID]...)
		sort.Strings(actual)
		sort.Strings(want)
		if !equal(actual, want) {
			return invalid("supersession projection differs")
		}
	}
	for _, p := range s.Steps {
		if p.ParentStepID != nil {
			parents[p.StepID] = []string{*p.ParentStepID}
		}
		start := events[p.StartEventID]
		if start.Type != "step_started" || start.LaneID != p.LaneID {
			return invalid("step start differs")
		}
		var in EventInput
		if json.Unmarshal(start.Data, &in) != nil || in.Step == nil || p.Kind != in.Step.Kind || !equal(p.ParentStepID, in.Step.ParentStepID) || !equal(p.Waiting, in.Waiting) {
			return invalid("step projection differs")
		}
		if p.EndEventID != nil {
			end := events[*p.EndEventID]
			if end.Type != "step_finished" || !equal(start.Actor, end.Actor) || !equal(start.RunID, end.RunID) {
				return invalid("step finish differs")
			}
		}
		state := "no_end_recorded"
		var duration *float64
		if p.EndEventID != nil {
			end := events[*p.EndEventID]
			state = "time_unknown"
			if start.OccurredAt != nil && end.OccurredAt != nil {
				a, _ := parseTime(*start.OccurredAt)
				z, _ := parseTime(*end.OccurredAt)
				if z.Before(a) {
					state = "time_conflict"
				} else if start.TimeBasis.Uncertainty == nil && end.TimeBasis.Uncertainty == nil && equal(start.TimeBasis.Clock, end.TimeBasis.Clock) {
					state = "closed"
					duration = ptr(z.Sub(a).Seconds())
				}
			}
		}
		sameDuration := p.DurationSeconds == nil && duration == nil || p.DurationSeconds != nil && duration != nil && *p.DurationSeconds == *duration
		if p.State != state || !sameDuration {
			return invalid("step duration or state differs from its evidence")
		}
	}
	if cyclic(edges) || cyclic(parents) {
		return invalid("cyclic snapshot relationships")
	}
	copy := s
	copy.Current.Axes = map[string]Axis{}
	copy.deriveAxes()
	if !equal(copy.Current.Axes, s.Current.Axes) {
		return invalid("current axes differ from typed evidence")
	}
	return nil
}
func validateNativeData(e Event) error {
	var m map[string]any
	if err := json.Unmarshal(e.Data, &m); err != nil || m == nil {
		return invalid("native data must be a typed object")
	}
	str := func(k string) string { v, _ := m[k].(string); return v }
	switch {
	case m["axis"] != nil:
		if len(m) != 2 || !one(str("axis"), "transport", "reported", "technical", "human", "integration") || !textOK(str("value"), 256, false) || e.Outcome == nil || *e.Outcome != str("value") {
			return invalid("invalid native axis data")
		}
	case m["basis"] != nil:
		if len(m) != 1 || !textOK(str("basis"), 256, false) {
			return invalid("invalid native state data")
		}
	case m["operation"] != nil:
		if len(m) != 2 || str("operation") != e.Type || !hexHash.MatchString(strings.TrimPrefix(str("manifest_sha256"), "sha256:")) {
			return invalid("invalid content data")
		}
	case m["preparation_id"] != nil:
		if len(m) != 2 || e.Type != "preparation" || !textOK(str("preparation_id"), 256, false) || !hexHash.MatchString(strings.TrimPrefix(str("preparation_sha256"), "sha256:")) {
			return invalid("invalid preparation data")
		}
	case m["run_event_type"] != nil:
		if len(m) != 3 || e.Type != "run_"+str("run_event_type") || !hexHash.MatchString(strings.TrimPrefix(str("event_sha256"), "sha256:")) {
			return invalid("invalid run data")
		}
	default:
		return invalid("unknown native data shape")
	}
	return nil
}
func validateObservationSnapshot(o Observation, s Snapshot) error {
	actors := map[string]bool{}
	runs := map[string]bool{}
	for _, e := range s.Events {
		if e.Actor != nil {
			actors[e.Actor.ID] = true
		}
		if e.RunID != nil {
			runs[*e.RunID] = true
		}
	}
	if i := o.Targets.Interval; i != nil {
		a, e := parseTime(i.From)
		z, f := parseTime(i.To)
		if e != nil || f != nil || z.Before(a) || len(i.ActorIDs)+len(i.RunIDs) == 0 {
			return invalid("invalid snapshot target interval")
		}
		for _, id := range i.ActorIDs {
			if !actors[id] {
				return invalid("unknown target actor")
			}
		}
		for _, id := range i.RunIDs {
			if !runs[id] {
				return invalid("unknown target run")
			}
		}
	}
	return nil
}
func cyclic(edges map[string][]string) bool {
	state := map[string]int{}
	var visit func(string) bool
	visit = func(id string) bool {
		if state[id] == 1 {
			return true
		}
		if state[id] == 2 {
			return false
		}
		state[id] = 1
		for _, next := range edges[id] {
			if visit(next) {
				return true
			}
		}
		state[id] = 2
		return false
	}
	for id := range edges {
		if visit(id) {
			return true
		}
	}
	return false
}
