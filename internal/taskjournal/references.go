package taskjournal

func validateReferences(c Common, ev *EventInput, ob *ObservationInput, s Snapshot, runs map[string]RunBinding) error {
	if r := c.RunBinding; r != nil {
		actual, ok := runs[r.RunID]
		if !ok || !equal(*r, actual) {
			return conflict("run/request/preparation is not bound to this workspace and Task")
		}
	}
	events := map[string]Event{}
	steps := map[string]Step{}
	actors := map[string]bool{}
	knownRuns := map[string]bool{}
	for _, e := range s.Events {
		events[e.EventID] = e
		if e.Actor != nil {
			actors[e.Actor.ID] = true
		}
		if e.RunID != nil {
			knownRuns[*e.RunID] = true
		}
	}
	for _, p := range s.Steps {
		steps[p.StepID] = p
	}
	if ev != nil {
		if p := ev.Candidate; p != nil {
			if s.Task.Worktree == nil || p.RepoID != s.Task.RepoID || p.WorktreeID != string(s.Task.Worktree.ID) {
				return conflict("candidate repository/worktree differs from Task")
			}
		}
		for _, r := range ev.Relations {
			target, ok := events[r.EventID]
			if !ok {
				return conflict("relation is not an existing event in this Task")
			}
			if r.Type == "supersedes" && target.Origin != "contribution" {
				return conflict("contributions cannot supersede native facts")
			}
		}
		if p := ev.Step; p != nil {
			if p.ParentStepID != nil {
				if *p.ParentStepID == p.ID {
					return conflict("step cannot parent itself")
				}
				if _, ok := steps[*p.ParentStepID]; !ok {
					return conflict("parent step does not exist")
				}
			}
			old, ok := steps[p.ID]
			if ev.Type == "step_started" {
				if ok {
					return conflict("step already has an original start")
				}
			}
			if ev.Type == "step_finished" || ev.Type == "note" {
				if !ok {
					return conflict("step requires an existing start")
				}
				start := events[old.StartEventID]
				if old.Kind != p.Kind || !equal(old.ParentStepID, p.ParentStepID) {
					return conflict("step kind or parent differs")
				}
				if ev.Type == "step_finished" {
					if old.EndEventID != nil {
						return conflict("step already has an original finish")
					}
					var run *string
					if c.RunBinding != nil {
						run = &c.RunBinding.RunID
					}
					if !equal(start.Actor, &c.Actor) || !equal(start.RunID, run) {
						return conflict("finish actor or run differs from start")
					}
				}
			}
		}
		return nil
	}
	for _, id := range ob.Targets.EventIDs {
		if _, ok := events[id]; !ok {
			return conflict("observation event target does not exist")
		}
	}
	for _, id := range ob.Targets.StepIDs {
		if _, ok := steps[id]; !ok {
			return conflict("observation step target does not exist")
		}
	}
	if i := ob.Targets.Interval; i != nil {
		for _, id := range i.ActorIDs {
			if !actors[id] {
				return conflict("unknown interval actor")
			}
		}
		for _, id := range i.RunIDs {
			if !knownRuns[id] {
				return conflict("unknown interval run")
			}
		}
	}
	if ob.ProposalEventID != nil {
		proposal := false
		accepted, rejected, applied := false, false, false
		for _, o := range s.Observations {
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
	}
	return nil
}
