package taskjournal

func validateReferences(c Common, ev *EventInput, ob *ObservationInput, s Snapshot) error {
	if c.RunBinding != nil {
		return conflict("run/request/preparation is not bound to this workspace and Task")
	}
	x := indexReferences(s)
	if ev != nil {
		return validateEventReferences(c, ev, s.Task, x)
	}
	return validateObservationReferences(ob, s.Observations, x)
}

// referenceIndex is the part of a built snapshot that a new contribution may
// refer to: existing events, steps, and the actors and runs seen on events.
type referenceIndex struct {
	events map[string]Event
	steps  map[string]Step
	actors map[string]bool
	runs   map[string]bool
}

func indexReferences(s Snapshot) referenceIndex {
	x := referenceIndex{map[string]Event{}, map[string]Step{}, map[string]bool{}, map[string]bool{}}
	for _, e := range s.Events {
		x.events[e.EventID] = e
		if e.Actor != nil {
			x.actors[e.Actor.ID] = true
		}
		if e.RunID != nil {
			x.runs[*e.RunID] = true
		}
	}
	for _, p := range s.Steps {
		x.steps[p.StepID] = p
	}
	return x
}
func validateEventReferences(c Common, ev *EventInput, task Task, x referenceIndex) error {
	if p := ev.Candidate; p != nil {
		if task.Worktree == nil || p.RepoID != task.RepoID || p.WorktreeID != string(task.Worktree.ID) {
			return conflict("candidate repository/worktree differs from Task")
		}
	}
	for _, r := range ev.Relations {
		target, ok := x.events[r.EventID]
		if !ok {
			return conflict("relation is not an existing event in this Task")
		}
		if r.Type == "supersedes" && target.Origin != "contribution" {
			return conflict("contributions cannot supersede native facts")
		}
	}
	if ev.Step != nil {
		return validateStepReference(c, ev, x)
	}
	return nil
}

// A start must be new; a finish or note must match the original start, and a
// finish must come from the same actor and run as that start.
func validateStepReference(c Common, ev *EventInput, x referenceIndex) error {
	p := ev.Step
	if p.ParentStepID != nil {
		if *p.ParentStepID == p.ID {
			return conflict("step cannot parent itself")
		}
		if _, ok := x.steps[*p.ParentStepID]; !ok {
			return conflict("parent step does not exist")
		}
	}
	old, ok := x.steps[p.ID]
	if ev.Type == "step_started" && ok {
		return conflict("step already has an original start")
	}
	if ev.Type != "step_finished" && ev.Type != "note" {
		return nil
	}
	if !ok {
		return conflict("step requires an existing start")
	}
	if old.Kind != p.Kind || !equal(old.ParentStepID, p.ParentStepID) {
		return conflict("step kind or parent differs")
	}
	if ev.Type != "step_finished" {
		return nil
	}
	if old.EndEventID != nil {
		return conflict("step already has an original finish")
	}
	var run *string
	if c.RunBinding != nil {
		run = &c.RunBinding.RunID
	}
	start := x.events[old.StartEventID]
	if !equal(start.Actor, &c.Actor) || !equal(start.RunID, run) {
		return conflict("finish actor or run differs from start")
	}
	return nil
}
func validateObservationReferences(ob *ObservationInput, observations []Observation, x referenceIndex) error {
	for _, id := range ob.Targets.EventIDs {
		if _, ok := x.events[id]; !ok {
			return conflict("observation event target does not exist")
		}
	}
	for _, id := range ob.Targets.StepIDs {
		if _, ok := x.steps[id]; !ok {
			return conflict("observation step target does not exist")
		}
	}
	if i := ob.Targets.Interval; i != nil {
		for _, id := range i.ActorIDs {
			if !x.actors[id] {
				return conflict("unknown interval actor")
			}
		}
		for _, id := range i.RunIDs {
			if !x.runs[id] {
				return conflict("unknown interval run")
			}
		}
	}
	if ob.ProposalEventID != nil {
		return validateFollowUp(ob, observations)
	}
	return nil
}

// A follow-up names an existing proposal; applied needs an unconflicted
// acceptance and assessment needs a prior applied record for that proposal.
func validateFollowUp(ob *ObservationInput, observations []Observation) error {
	proposal := false
	accepted, rejected, applied := false, false, false
	for _, o := range observations {
		if o.EventID == *ob.ProposalEventID && o.Type == "proposal" {
			proposal = true
		}
		if o.ProposalEventID != nil && *o.ProposalEventID == *ob.ProposalEventID {
			if o.Decision != nil {
				accepted = accepted || *o.Decision == "accepted"
				rejected = rejected || *o.Decision == "rejected"
			}
			applied = applied || o.Type == "applied"
		}
	}
	if !proposal {
		return conflict("follow-up requires an existing proposal in this Task")
	}
	if ob.Type == "applied" && (!accepted || rejected) {
		return conflict("applied requires an unconflicted accepted decision")
	}
	if ob.Type == "assessment" && !applied {
		return conflict("assessment requires a prior applied record")
	}
	return nil
}
