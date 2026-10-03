package taskjournal

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"sort"
	"strings"
	"time"
)

func (s Service) build(b workspace.TaskJournalBasis, records []stored) (Snapshot, map[string]RunBinding, error) {
	out, runs, e := s.native(b)
	if e != nil {
		return out, runs, e
	}
	var previousRegistration *time.Time
	for _, r := range records {
		c := r.Common
		src := out.source("contribution", r.Locator, &r.SHA256, &r.Record.Sequence, r.Record.PreviousSHA256)
		ev := Event{EventID: r.Record.EventID, Origin: "contribution", Actor: ptr(c.Actor), ActivityID: c.ActivityID, OccurredAt: normalize(c.OccurredAt), OriginalOccurredAt: c.OccurredAt, RecordedAt: ptr(r.Record.RecordedAt), TimeBasis: c.TimeBasis, SourceIDs: []string{src}, SourceSequence: ptr(r.Record.Sequence), Relations: []Relation{}, SupersededBy: []string{}, Data: r.Record.Input}
		if c.RunBinding != nil {
			ev.RunID = &c.RunBinding.RunID
			if bound, ok := runs[c.RunBinding.RunID]; !ok || !equal(bound, *c.RunBinding) {
				out.reason("run_unbound", "A historical contribution's native run can no longer be fully validated", src)
			}
		}
		for _, source := range c.Sources {
			ev.SourceIDs = append(ev.SourceIDs, out.source("artifact", source.Locator, &source.SHA256, nil, nil))
		}
		if r.Event != nil {
			x := r.Event
			ev.Type = x.Type
			ev.Title = x.Title
			ev.Outcome = x.Outcome
			ev.Candidate = x.Candidate
			ev.Relations = x.Relations
			if x.Step != nil {
				ev.StepID = &x.Step.ID
			}
		} else {
			x := r.Observation
			ev.Type = "observation_" + x.Type
			ev.Title = x.Finding
			if len([]rune(ev.Title)) > 256 {
				ev.Title = string([]rune(ev.Title)[:253]) + "..."
			}
			out.Observations = append(out.Observations, Observation{*x, ev.EventID})
		}
		out.Events = append(out.Events, ev)
		registration, _ := parseTime(r.Record.RecordedAt)
		if previousRegistration != nil && registration.Before(*previousRegistration) {
			out.Coverage.TimeConflictEventIDs = append(out.Coverage.TimeConflictEventIDs, ev.EventID)
			out.reason("time_conflict", "Contribution sequence contradicts its registration clock", src)
		}
		previousRegistration = &registration
	}
	out.closeReferences()
	out.derive()
	return out, runs, nil
}

// A preserved reference is evidence that a target existed, not evidence of the
// missing native fact. Keep an explicit unavailable point so partial exports
// retain reference closure without inventing its time, outcome or candidate.
func (s *Snapshot) closeReferences() {
	known := map[string]bool{}
	for _, e := range s.Events {
		known[e.EventID] = true
	}
	add := func(id string, source string) {
		if known[id] {
			return
		}
		known[id] = true
		var ref Source
		for _, src := range s.Sources {
			if src.SourceID == source {
				ref = src
				break
			}
		}
		ref.SourceID = "src_" + digest([]string{"unavailable_reference", id, source})
		ref.Kind = "unavailable_reference"
		ref.Status = "unbound"
		s.Sources = append(s.Sources, ref)
		s.Events = append(s.Events, Event{EventID: id, Origin: "native", Type: "unavailable_reference", Title: "Unavailable historical reference", TimeBasis: unknownTime(), SourceIDs: []string{ref.SourceID}, Relations: []Relation{}, SupersededBy: []string{}, Data: raw(map[string]any{"basis": "Historical reference only; native source unavailable."})})
		s.reason("source_missing", "Referenced native fact is unavailable; its content and outcome remain unknown", ref.SourceID)
	}
	for _, e := range append([]Event{}, s.Events...) {
		for _, r := range e.Relations {
			if len(e.SourceIDs) > 0 {
				add(r.EventID, e.SourceIDs[0])
			}
		}
	}
	for _, o := range s.Observations {
		for _, e := range s.Events {
			if e.EventID == o.EventID && len(e.SourceIDs) > 0 {
				for _, id := range o.Targets.EventIDs {
					add(id, e.SourceIDs[0])
				}
				break
			}
		}
	}
}
func (s *Snapshot) derive() {
	lanes := map[string]Lane{}
	events := map[string]*Event{}
	steps := map[string]*Step{}
	for i := range s.Events {
		e := &s.Events[i]
		var actor *string
		if e.Actor != nil {
			actor = &e.Actor.ID
		}
		e.LaneID = "lane_" + digest([]any{actor, e.RunID})
		lanes[e.LaneID] = Lane{e.LaneID, actor, e.RunID}
		events[e.EventID] = e
		if e.OccurredAt == nil {
			s.Coverage.UnknownTimeEventIDs = append(s.Coverage.UnknownTimeEventIDs, e.EventID)
		}
	}
	for i := range s.Events {
		e := &s.Events[i]
		for _, r := range e.Relations {
			target, ok := events[r.EventID]
			if !ok {
				s.reason("source_missing", "Historical relation target is no longer available", e.SourceIDs...)
				continue
			}
			if r.Type == "supersedes" && target.Origin == "contribution" {
				target.SupersededBy = append(target.SupersededBy, e.EventID)
			}
		}
	}
	for i := range s.Events {
		e := &s.Events[i]
		if len(e.SupersededBy) > 1 {
			s.reason("contribution_conflict", "Several contributions supersede the same original; no timestamp wins", e.SourceIDs...)
		}
		if e.Origin != "contribution" || e.Type != "step_started" {
			continue
		}
		var in EventInput
		if json.Unmarshal(e.Data, &in) != nil || in.Step == nil {
			continue
		}
		p := in.Step
		steps[p.ID] = &Step{p.ID, p.Kind, p.ParentStepID, e.LaneID, e.EventID, nil, "no_end_recorded", nil, in.Waiting}
	}
	for _, e := range s.Events {
		if e.Origin != "contribution" || e.Type != "step_finished" || e.StepID == nil {
			continue
		}
		p := steps[*e.StepID]
		if p == nil {
			s.reason("source_invalid", "Step finish has no preserved start", e.SourceIDs...)
			continue
		}
		p.EndEventID = ptr(e.EventID)
		start := events[p.StartEventID]
		p.State = "time_unknown"
		if start.OccurredAt == nil || e.OccurredAt == nil {
			continue
		}
		a, _ := parseTime(*start.OccurredAt)
		z, _ := parseTime(*e.OccurredAt)
		if z.Before(a) {
			p.State = "time_conflict"
			s.Coverage.TimeConflictEventIDs = append(s.Coverage.TimeConflictEventIDs, start.EventID, e.EventID)
			s.reason("time_conflict", "Step finish precedes its start", e.SourceIDs...)
			continue
		}
		if start.TimeBasis.Uncertainty != nil || e.TimeBasis.Uncertainty != nil || !equal(start.TimeBasis.Clock, e.TimeBasis.Clock) {
			continue
		}
		p.State = "closed"
		p.DurationSeconds = ptr(z.Sub(a).Seconds())
	}
	for _, p := range steps {
		s.Steps = append(s.Steps, *p)
	}
	for _, l := range lanes {
		s.Lanes = append(s.Lanes, l)
	}
	decisions := map[string]map[string]bool{}
	for _, o := range s.Observations {
		if o.ProposalEventID != nil && o.Decision != nil {
			m := decisions[*o.ProposalEventID]
			if m == nil {
				m = map[string]bool{}
				decisions[*o.ProposalEventID] = m
			}
			m[*o.Decision] = true
		}
	}
	for id, m := range decisions {
		if len(m) > 1 {
			s.reason("observation_conflict", "Conflicting observer decisions for proposal "+id)
		}
	}
	s.deriveAxes()
	if len(s.Steps) == 0 {
		s.Coverage.Reasons = append(s.Coverage.Reasons, Reason{"no_step_reporting", []string{}, "No meaningful agent steps have been recorded; native facts are not a complete work or conversation history."})
	}
	sort.Slice(s.Events, func(i, j int) bool { return s.Events[i].EventID < s.Events[j].EventID })
	sort.Slice(s.Sources, func(i, j int) bool { return s.Sources[i].SourceID < s.Sources[j].SourceID })
	sort.Slice(s.Steps, func(i, j int) bool { return s.Steps[i].StepID < s.Steps[j].StepID })
	sort.Slice(s.Lanes, func(i, j int) bool { return s.Lanes[i].LaneID < s.Lanes[j].LaneID })
	sort.Slice(s.Observations, func(i, j int) bool { return s.Observations[i].EventID < s.Observations[j].EventID })
	sort.Strings(s.Coverage.UnknownTimeEventIDs)
	sort.Strings(s.Coverage.TimeConflictEventIDs)
}
func (s *Snapshot) deriveAxes() {
	statuses := map[string]string{}
	for _, src := range s.Sources {
		statuses[src.SourceID] = src.Status
	}
	for _, axis := range []string{"transport", "reported", "technical", "human", "integration"} {
		current := []string{}
		historical := []string{}
		values := map[string]bool{}
		for _, e := range s.Events {
			var data struct {
				Axis  string `json:"axis"`
				Value string `json:"value"`
			}
			_ = json.Unmarshal(e.Data, &data)
			if e.Origin == "contribution" {
				if e.Candidate == nil || e.Outcome == nil || len(e.SupersededBy) > 0 {
					continue
				}
				data.Axis = "reported"
				data.Value = *e.Outcome
			}
			if data.Axis != axis {
				continue
			}
			historical = append(historical, e.EventID)
			valid := true
			for _, id := range e.SourceIDs {
				valid = valid && statuses[id] == "valid"
			}
			if s.Current.Candidate == nil || e.Candidate == nil || !equal(e.Candidate, s.Current.Candidate) || !valid {
				continue
			}
			current = append(current, e.EventID)
			values[data.Value] = true
		}
		a := Axis{"unknown", []string{}, "missing"}
		if len(historical) > 0 {
			a.EventIDs = historical
			a.BasisStatus = "historical"
		}
		if len(values) == 1 {
			for v := range values {
				a.Value = v
			}
			a.EventIDs = current
			a.BasisStatus = "current"
		} else if len(values) > 1 {
			a.Value = "conflict"
			a.EventIDs = current
			a.BasisStatus = "conflict"
		}
		sort.Strings(a.EventIDs)
		s.Current.Axes[axis] = a
	}
}

func snapshotDigest(s Snapshot) string {
	var m map[string]any
	_ = json.Unmarshal(raw(s), &m)
	delete(m, "snapshot_id")
	delete(m, "selection")
	return "sha256:" + digest(m)
}
func Select(s *Snapshot, o Options) error {
	if o.Order == "" {
		o.Order = "occurred"
	}
	if !one(o.Order, "occurred", "recorded") {
		return invalid("--order must be occurred or recorded")
	}
	if o.View != "" && !one(o.View, "summary", "timeline", "details") {
		return invalid("--view must be summary, timeline or details")
	}
	selection := Selection{Order: o.Order, EventIDs: []string{}}
	if o.RunID != "" {
		selection.RunID = &o.RunID
	}
	if o.ActorID != "" {
		selection.ActorID = &o.ActorID
	}
	chosen := []Event{}
	for _, e := range s.Events {
		if o.RunID != "" && (e.RunID == nil || *e.RunID != o.RunID) {
			continue
		}
		if o.ActorID != "" && (e.Actor == nil || e.Actor.ID != o.ActorID) {
			continue
		}
		chosen = append(chosen, e)
	}
	get := func(e Event) *string {
		if o.Order == "recorded" {
			return e.RecordedAt
		}
		return e.OccurredAt
	}
	sort.Slice(chosen, func(i, j int) bool {
		a, b := get(chosen[i]), get(chosen[j])
		if a == nil && b != nil {
			return false
		}
		if a != nil && b == nil {
			return true
		}
		if a != nil && b != nil {
			at, _ := parseTime(*a)
			bt, _ := parseTime(*b)
			if !at.Equal(bt) {
				return at.Before(bt)
			}
		}
		return chosen[i].EventID < chosen[j].EventID
	})
	for _, e := range chosen {
		selection.EventIDs = append(selection.EventIDs, e.EventID)
	}
	s.Selection = selection
	return nil
}
func (s Service) Show(task string, o Options) (Snapshot, error) {
	b, e := workspace.ReadTaskJournalBasis(s.Workspace, workspace.TaskID(task))
	if e != nil {
		return Snapshot{}, e
	}
	records, readErr := readStore(b.Workspace.Root, task, workspaceID(b.Workspace.Root, b.Workspace.MarkerSHA256))
	out, _, e := s.build(b, records)
	if e != nil {
		return out, e
	}
	if readErr != nil {
		out.reason("source_invalid", "Journal is damaged; showing the validated prefix and independent native facts: "+readErr.Error())
	}
	// Observe freshness, without a global-transaction claim or any repairs.
	if e = s.check("before_freshness"); e != nil {
		return out, e
	}
	for _, src := range append([]Source{}, out.Sources...) {
		if src.Status != "valid" || src.SHA256 == nil {
			continue
		}
		h, e := sourceHash(src.Locator)
		if e != nil || h != *src.SHA256 {
			out.reason("source_changed", "Source changed during snapshot reading", src.SourceID)
			for i := range out.Sources {
				if out.Sources[i].SourceID == src.SourceID {
					out.Sources[i].Status = "changed"
				}
			}
		}
	}
	if b.Task.Worktree != nil && out.Current.Candidate != nil {
		g, e := workspace.ObserveTaskJournalWorktree(s.Workspace, b.Task.Worktree.Locator)
		if e != nil || !g.Clean || g.OID != out.Current.Candidate.OID || g.Tree != out.Current.Candidate.Tree {
			out.Current.Candidate = nil
			for k, a := range out.Current.Axes {
				a.Value = "unknown"
				if len(a.EventIDs) > 0 {
					a.BasisStatus = "historical"
				} else {
					a.BasisStatus = "missing"
				}
				out.Current.Axes[k] = a
			}
			out.reason("candidate_changed", "Candidate changed during snapshot reading")
		}
	}
	out.deriveAxes()
	sort.Slice(out.Coverage.Reasons, func(i, j int) bool { return digest(out.Coverage.Reasons[i]) < digest(out.Coverage.Reasons[j]) })
	out.SnapshotID = snapshotDigest(out)
	if e = Select(&out, o); e != nil {
		return out, e
	}
	return out, nil
}

func ReadSnapshot(path, task string, o Options) (Snapshot, error) {
	b, e := readFile(path, 64<<20)
	if e != nil {
		return Snapshot{}, e
	}
	var s Snapshot
	if e = decode(b, 64<<20, &s); e != nil {
		return s, e
	}
	if s.Kind != "ply.workspace.task-journal-snapshot" || s.SchemaVersion != 1 || s.Task.TaskID != task || s.NextTransitionAuthorized || s.SnapshotID != snapshotDigest(s) {
		return s, invalid("snapshot identity, version or digest differs")
	}
	if _, e = workspace.ParseTaskID(task); e != nil {
		return s, e
	}
	if _, e = parseTime(s.AsOf); e != nil || !strings.HasSuffix(s.AsOf, "Z") {
		return s, invalid("invalid snapshot as_of")
	}
	if e = validateSnapshot(s); e != nil {
		return s, e
	}
	e = Select(&s, o)
	return s, e
}
func validateSnapshot(s Snapshot) error {
	if s.Workspace.ID != workspaceID(s.Workspace.Root, s.Workspace.MarkerSHA256) || !one(s.Coverage.State, "complete", "partial") {
		return invalid("invalid workspace or coverage")
	}
	sources := map[string]bool{}
	events := map[string]Event{}
	lanes := map[string]bool{}
	steps := map[string]bool{}
	for _, src := range s.Sources {
		if sources[src.SourceID] || src.SourceID == "" || !one(src.Status, "valid", "missing", "changed", "invalid", "unbound") {
			return invalid("invalid source")
		}
		sources[src.SourceID] = true
	}
	for _, l := range s.Lanes {
		if lanes[l.LaneID] || l.LaneID != "lane_"+digest([]any{l.ActorID, l.RunID}) {
			return invalid("invalid lane")
		}
		lanes[l.LaneID] = true
	}
	for _, e := range s.Events {
		if _, ok := events[e.EventID]; ok || e.EventID == "" || !one(e.Origin, "native", "contribution") || !lanes[e.LaneID] || !textOK(e.Title, 256, false) {
			return invalid("invalid event")
		}
		if e.OccurredAt != nil {
			if _, err := parseTime(*e.OccurredAt); err != nil {
				return invalid("invalid event time")
			}
		}
		if err := validateTime(e.OriginalOccurredAt, e.TimeBasis); err != nil {
			return err
		}
		if !equal(e.OccurredAt, normalize(e.OriginalOccurredAt)) {
			return invalid("normalized event time differs")
		}
		if e.RecordedAt != nil {
			if _, err := parseTime(*e.RecordedAt); err != nil {
				return invalid("invalid recorded time")
			}
		}
		if len(e.SourceIDs) == 0 {
			return invalid("event has no source")
		}
		for _, id := range e.SourceIDs {
			if !sources[id] {
				return invalid("unknown event source")
			}
		}
		events[e.EventID] = e
	}
	for _, p := range s.Steps {
		if steps[p.StepID] || !lanes[p.LaneID] || !one(p.State, "closed", "no_end_recorded", "time_unknown", "time_conflict") {
			return invalid("invalid step")
		}
		start, ok := events[p.StartEventID]
		if !ok || start.StepID == nil || *start.StepID != p.StepID {
			return invalid("unknown step start")
		}
		if p.EndEventID != nil {
			end, ok := events[*p.EndEventID]
			if !ok || end.StepID == nil || *end.StepID != p.StepID {
				return invalid("unknown step finish")
			}
		}
		if p.DurationSeconds != nil && (*p.DurationSeconds < 0 || p.State != "closed") {
			return invalid("invalid duration")
		}
		steps[p.StepID] = true
	}
	for _, p := range s.Steps {
		if p.ParentStepID != nil && !steps[*p.ParentStepID] {
			return invalid("unknown parent step")
		}
	}
	for _, e := range s.Events {
		if e.StepID != nil && !steps[*e.StepID] {
			return invalid("unknown event step")
		}
		for _, r := range e.Relations {
			target, ok := events[r.EventID]
			if !ok || r.Type == "supersedes" && target.Origin == "native" {
				return invalid("invalid event relation")
			}
		}
		for _, id := range e.SupersededBy {
			if _, ok := events[id]; !ok {
				return invalid("unknown superseding event")
			}
		}
	}
	for _, o := range s.Observations {
		if _, ok := events[o.EventID]; !ok {
			return invalid("unknown observation event")
		}
		if o.ProposalEventID != nil {
			p, ok := events[*o.ProposalEventID]
			if !ok || p.Type != "observation_proposal" {
				return invalid("invalid observation proposal")
			}
		}
		for _, id := range o.Targets.EventIDs {
			if _, ok := events[id]; !ok {
				return invalid("unknown observation target")
			}
		}
		for _, id := range o.Targets.StepIDs {
			if !steps[id] {
				return invalid("unknown observation step")
			}
		}
	}
	if len(s.Current.Axes) != 5 {
		return invalid("invalid current axes")
	}
	for _, key := range []string{"transport", "reported", "technical", "human", "integration"} {
		a, ok := s.Current.Axes[key]
		if !ok || !one(a.BasisStatus, "current", "historical", "missing", "conflict") {
			return invalid("invalid axis")
		}
		for _, id := range a.EventIDs {
			if _, ok := events[id]; !ok {
				return invalid("unknown axis event")
			}
		}
	}
	for _, ids := range [][]string{s.Selection.EventIDs, s.Coverage.UnknownTimeEventIDs, s.Coverage.TimeConflictEventIDs, s.Current.NextAction.SourceEventIDs} {
		for _, id := range ids {
			if _, ok := events[id]; !ok {
				return invalid("unknown snapshot event reference")
			}
		}
	}
	for _, r := range s.Coverage.Reasons {
		for _, id := range r.SourceIDs {
			if !sources[id] {
				return invalid("unknown coverage source")
			}
		}
	}
	return validateSnapshotSemantics(s)
}
