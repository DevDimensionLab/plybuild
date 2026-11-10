package taskjournal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These tests pin the exact error class and message of each validation phase so
// that restructuring the validators cannot silently reorder or reclassify them.

func pathCommon(kind string) Common {
	return Common{Kind: "ply.workspace.task-journal-" + kind + "-input", SchemaVersion: 1, PublicationKey: "key", TaskID: "task", Actor: Actor{ID: "developer", Role: "developer"}, OccurredAt: ptr("2026-10-03T10:00:00Z"), TimeBasis: TimeBasis{"reported", ptr("fixture"), "second", nil}, Sources: []SourceRef{{"/evidence", strings.Repeat("a", 64)}}}
}
func pathEvent() EventInput {
	return EventInput{Common: pathCommon("event"), Type: "note", Title: "Statement", Relations: []Relation{}}
}
func pathObservation() ObservationInput {
	return ObservationInput{Common: pathCommon("observation"), Type: "proposal", Targets: Targets{EventIDs: []string{"e1"}, StepIDs: []string{}}, Finding: "finding", Action: ptr("act"), Owner: ptr("owner"), NextSignal: ptr("signal")}
}
func assertError(t *testing.T, e error, want string) {
	t.Helper()
	if e == nil || e.Error() != want {
		t.Fatalf("error = %v, want %q", e, want)
	}
}

func TestParseInputErrorPaths(t *testing.T) {
	oid := strings.Repeat("1", 40)
	events := map[string]struct {
		edit func(*EventInput)
		want string
	}{
		"type":               {func(e *EventInput) { e.Type = "other" }, "invalid event type, title or detail (detail must be empty in v1)"},
		"detail":             {func(e *EventInput) { e.Detail = "x" }, "invalid event type, title or detail (detail must be empty in v1)"},
		"outcome":            {func(e *EventInput) { e.Outcome = ptr("maybe") }, "invalid reported outcome"},
		"start without step": {func(e *EventInput) { e.Type = "step_started" }, "start and finish require a step"},
		"decision step":      {func(e *EventInput) { e.Type = "decision"; e.Step = &StepRef{"s", "planning", nil} }, "decision is a point, not a step"},
		"step kind":          {func(e *EventInput) { e.Step = &StepRef{"s", "dreaming", nil} }, "invalid step"},
		"waiting reason":     {func(e *EventInput) { e.Step = &StepRef{"s", "waiting", nil} }, "waiting step needs explicit reason, dependency and next actor"},
		"waiting step":       {func(e *EventInput) { e.Waiting = &Waiting{"r", "d", "n"} }, "invalid waiting metadata"},
		"candidate":          {func(e *EventInput) { e.Candidate = &Candidate{"repo", "wt", "bad", oid} }, "invalid candidate"},
		"relation":           {func(e *EventInput) { e.Relations = []Relation{{"likes", "e1"}} }, "invalid relation"},
		"supersede step": {func(e *EventInput) {
			e.Type = "step_started"
			e.Step = &StepRef{"s", "planning", nil}
			e.Relations = []Relation{{"supersedes", "e1"}}
		}, "only notes and decisions supersede contributions"},
		"common": {func(e *EventInput) { e.TaskID = "other" }, "input kind, version, Task or publication key differs"},
	}
	for name, c := range events {
		t.Run("event/"+name, func(t *testing.T) {
			in := pathEvent()
			c.edit(&in)
			b, _ := Canonical(in)
			_, _, _, _, e := parseInput(b, "event", "task")
			assertError(t, e, "task_journal_invalid_input: "+c.want)
		})
	}
	observations := map[string]struct {
		edit func(*ObservationInput)
		want string
	}{
		"type":             {func(o *ObservationInput) { o.Type = "other" }, "invalid observation"},
		"proposal fields":  {func(o *ObservationInput) { o.Owner = nil }, "proposal requires action, owner and next signal"},
		"follow-up":        {func(o *ObservationInput) { o.Type = "applied" }, "follow-up requires a proposal"},
		"decision value":   {func(o *ObservationInput) { o.Type = "decision"; o.ProposalEventID = ptr("p") }, "invalid decision"},
		"stray decision":   {func(o *ObservationInput) { o.Decision = ptr("accepted") }, "decision only belongs to decision records"},
		"assessment value": {func(o *ObservationInput) { o.Type = "assessment"; o.ProposalEventID = ptr("p") }, "invalid assessment"},
		"stray assessment": {func(o *ObservationInput) { o.Assessment = ptr("better") }, "assessment only belongs to assessment records"},
		"no target":        {func(o *ObservationInput) { o.Targets.EventIDs = []string{} }, "observation needs a target"},
		"unscoped interval": {func(o *ObservationInput) {
			o.Targets.Interval = &Interval{"2026-10-03T10:00:00Z", "2026-10-03T11:00:00Z", []string{}, []string{}}
		}, "invalid or unscoped target interval"},
	}
	for name, c := range observations {
		t.Run("observation/"+name, func(t *testing.T) {
			in := pathObservation()
			c.edit(&in)
			b, _ := Canonical(in)
			_, _, _, _, e := parseInput(b, "observation", "task")
			assertError(t, e, "task_journal_invalid_input: "+c.want)
		})
	}
	_, _, _, _, e := parseInput([]byte("{}"), "other", "task")
	assertError(t, e, "task_journal_invalid_input: unknown input type")
	b, _ := Canonical(pathEvent())
	c, ev, ob, canonical, e := parseInput(b, "event", "task")
	if e != nil || ev == nil || ob != nil || c.PublicationKey != "key" || string(canonical) != string(b) {
		t.Fatalf("valid event parsed as %v %v %v %v", c, ev, ob, e)
	}
}

func historyRecord(id string, ev *EventInput, ob *ObservationInput) stored {
	return stored{Record: Record{EventID: id}, Event: ev, Observation: ob}
}

func TestValidateHistoryErrorPaths(t *testing.T) {
	start := pathEvent()
	start.Type = "step_started"
	start.Step = &StepRef{"s", "planning", nil}
	finish := start
	finish.Type = "step_finished"
	proposal := pathObservation()
	accepted := pathObservation()
	accepted.Type = "decision"
	accepted.ProposalEventID = ptr("p")
	accepted.Decision = ptr("accepted")
	prior := []stored{historyRecord("e1", &start, nil), historyRecord("p", nil, &proposal)}
	events := map[string]struct {
		prior []stored
		edit  func(*EventInput)
		want  string
	}{
		"relation":       {prior, func(e *EventInput) { e.Relations = []Relation{{"responds_to", "missing"}} }, "invalid historical relation"},
		"supersede":      {prior, func(e *EventInput) { e.Relations = []Relation{{"supersedes", "native_" + strings.Repeat("a", 64)}} }, "invalid historical relation"},
		"parent":         {prior, func(e *EventInput) { e.Step = &StepRef{"c", "planning", ptr("missing")} }, "invalid historical parent"},
		"self parent":    {prior, func(e *EventInput) { e.Step = &StepRef{"s", "planning", ptr("s")} }, "invalid historical parent"},
		"duplicate":      {prior, func(e *EventInput) { e.Type = "step_started"; e.Step = &StepRef{"s", "planning", nil} }, "duplicate historical step start"},
		"no start":       {prior, func(e *EventInput) { e.Step = &StepRef{"other", "planning", nil} }, "historical step has no matching start"},
		"double finish":  {append(prior, historyRecord("e2", &finish, nil)), func(e *EventInput) { e.Type = "step_finished"; e.Step = start.Step }, "invalid historical step finish"},
		"finish actor":   {prior, func(e *EventInput) { e.Type = "step_finished"; e.Step = start.Step; e.Actor.ID = "other" }, "invalid historical step finish"},
		"native related": {prior, func(e *EventInput) { e.Relations = []Relation{{"responds_to", "native_" + strings.Repeat("a", 64)}} }, ""},
	}
	for name, c := range events {
		t.Run("event/"+name, func(t *testing.T) {
			in := pathEvent()
			c.edit(&in)
			e := validateHistory(c.prior, &in, nil)
			if c.want == "" {
				if e != nil {
					t.Fatal(e)
				}
				return
			}
			assertError(t, e, "task_journal_conflict: "+c.want)
		})
	}
	observations := map[string]struct {
		prior []stored
		edit  func(*ObservationInput)
		want  string
	}{
		"target":        {prior, func(o *ObservationInput) { o.Targets.EventIDs = []string{"missing"} }, "unknown historical observation target"},
		"step":          {prior, func(o *ObservationInput) { o.Targets.StepIDs = []string{"missing"} }, "unknown historical observation step"},
		"proposal":      {prior, func(o *ObservationInput) { o.ProposalEventID = ptr("e1") }, "unknown historical proposal"},
		"applied":       {prior, func(o *ObservationInput) { o.Type = "applied"; o.ProposalEventID = ptr("p") }, "historical applied without unconflicted acceptance"},
		"assessment":    {append(prior, historyRecord("d", nil, &accepted)), func(o *ObservationInput) { o.Type = "assessment"; o.ProposalEventID = ptr("p") }, "historical assessment without applied"},
		"accepted only": {append(prior, historyRecord("d", nil, &accepted)), func(o *ObservationInput) { o.Type = "applied"; o.ProposalEventID = ptr("p") }, ""},
	}
	for name, c := range observations {
		t.Run("observation/"+name, func(t *testing.T) {
			in := pathObservation()
			c.edit(&in)
			e := validateHistory(c.prior, nil, &in)
			if c.want == "" {
				if e != nil {
					t.Fatal(e)
				}
				return
			}
			assertError(t, e, "task_journal_conflict: "+c.want)
		})
	}
}

func referenceSnapshot() Snapshot {
	start := Event{EventID: "start", Origin: "contribution", Actor: &Actor{ID: "developer", Role: "developer"}, StepID: ptr("s")}
	native := Event{EventID: "native", Origin: "native", RunID: ptr("run")}
	p := pathObservation()
	return Snapshot{Task: Task{RepoID: "repo"}, Events: []Event{start, native, {EventID: "p", Origin: "contribution"}}, Steps: []Step{{StepID: "s", Kind: "planning", StartEventID: "start"}, {StepID: "done", Kind: "planning", StartEventID: "start", EndEventID: ptr("start")}}, Observations: []Observation{{p, "p"}}}
}

func TestValidateReferencesErrorPaths(t *testing.T) {
	events := map[string]struct {
		edit func(*EventInput, *Snapshot)
		want string
	}{
		"run":           {func(e *EventInput, _ *Snapshot) { e.RunBinding = &RunBinding{RunID: "run"} }, "run/request/preparation is not bound to this workspace and Task"},
		"candidate":     {func(e *EventInput, _ *Snapshot) { e.Candidate = &Candidate{RepoID: "repo"} }, "candidate repository/worktree differs from Task"},
		"relation":      {func(e *EventInput, _ *Snapshot) { e.Relations = []Relation{{"responds_to", "missing"}} }, "relation is not an existing event in this Task"},
		"supersede":     {func(e *EventInput, _ *Snapshot) { e.Relations = []Relation{{"supersedes", "native"}} }, "contributions cannot supersede native facts"},
		"self parent":   {func(e *EventInput, _ *Snapshot) { e.Step = &StepRef{"c", "planning", ptr("c")} }, "step cannot parent itself"},
		"parent":        {func(e *EventInput, _ *Snapshot) { e.Step = &StepRef{"c", "planning", ptr("missing")} }, "parent step does not exist"},
		"restart":       {func(e *EventInput, _ *Snapshot) { e.Type = "step_started"; e.Step = &StepRef{"s", "planning", nil} }, "step already has an original start"},
		"no start":      {func(e *EventInput, _ *Snapshot) { e.Step = &StepRef{"c", "planning", nil} }, "step requires an existing start"},
		"kind":          {func(e *EventInput, _ *Snapshot) { e.Step = &StepRef{"s", "review", nil} }, "step kind or parent differs"},
		"double finish": {func(e *EventInput, _ *Snapshot) { e.Type = "step_finished"; e.Step = &StepRef{"done", "planning", nil} }, "step already has an original finish"},
		"finish actor": {func(e *EventInput, _ *Snapshot) {
			e.Type = "step_finished"
			e.Step = &StepRef{"s", "planning", nil}
			e.Actor.ID = "other"
		}, "finish actor or run differs from start"},
		"valid finish": {func(e *EventInput, _ *Snapshot) { e.Type = "step_finished"; e.Step = &StepRef{"s", "planning", nil} }, ""},
	}
	for name, c := range events {
		t.Run("event/"+name, func(t *testing.T) {
			in := pathEvent()
			s := referenceSnapshot()
			c.edit(&in, &s)
			e := validateReferences(in.Common, &in, nil, s)
			if c.want == "" {
				if e != nil {
					t.Fatal(e)
				}
				return
			}
			assertError(t, e, "task_journal_conflict: "+c.want)
		})
	}
	decision := pathObservation()
	decision.Type = "decision"
	decision.ProposalEventID = ptr("p")
	observations := map[string]struct {
		edit func(*ObservationInput, *Snapshot)
		want string
	}{
		"target": {func(o *ObservationInput, _ *Snapshot) { o.Targets.EventIDs = []string{"missing"} }, "observation event target does not exist"},
		"step":   {func(o *ObservationInput, _ *Snapshot) { o.Targets.StepIDs = []string{"missing"} }, "observation step target does not exist"},
		"actor":  {func(o *ObservationInput, _ *Snapshot) { o.Targets.Interval = &Interval{ActorIDs: []string{"missing"}} }, "unknown interval actor"},
		"run": {func(o *ObservationInput, _ *Snapshot) {
			o.Targets.Interval = &Interval{ActorIDs: []string{"developer"}, RunIDs: []string{"missing"}}
		}, "unknown interval run"},
		"proposal":   {func(o *ObservationInput, _ *Snapshot) { o.ProposalEventID = ptr("start") }, "follow-up requires an existing proposal in this Task"},
		"applied":    {func(o *ObservationInput, _ *Snapshot) { o.Type = "applied"; o.ProposalEventID = ptr("p") }, "applied requires an unconflicted accepted decision"},
		"assessment": {func(o *ObservationInput, _ *Snapshot) { o.Type = "assessment"; o.ProposalEventID = ptr("p") }, "assessment requires a prior applied record"},
		"rejected too": {func(o *ObservationInput, s *Snapshot) {
			o.Type = "applied"
			o.ProposalEventID = ptr("p")
			s.Observations = append(s.Observations, decisionObservation(decision, "accepted"), decisionObservation(decision, "rejected"))
		}, "applied requires an unconflicted accepted decision"},
		"accepted": {func(o *ObservationInput, s *Snapshot) {
			o.Type = "applied"
			o.ProposalEventID = ptr("p")
			s.Observations = append(s.Observations, decisionObservation(decision, "accepted"))
		}, ""},
		"valid run": {func(o *ObservationInput, _ *Snapshot) { o.Targets.Interval = &Interval{RunIDs: []string{"run"}} }, ""},
	}
	for name, c := range observations {
		t.Run("observation/"+name, func(t *testing.T) {
			in := pathObservation()
			in.Targets.EventIDs = []string{"start"}
			s := referenceSnapshot()
			c.edit(&in, &s)
			e := validateReferences(in.Common, nil, &in, s)
			if c.want == "" {
				if e != nil {
					t.Fatal(e)
				}
				return
			}
			assertError(t, e, "task_journal_conflict: "+c.want)
		})
	}
}
func decisionObservation(o ObservationInput, decision string) Observation {
	o.Decision = ptr(decision)
	return Observation{o, "decision-" + decision}
}

// A shown snapshot with a step start, finish and observation exercises every
// structural phase of validateSnapshot.
func structuralSnapshot(t *testing.T) Snapshot {
	t.Helper()
	f := fixture(t)
	start := f.input
	start.Type = "step_started"
	start.Step = &StepRef{"step", "implementation", nil}
	first := f.append(t, start)
	finish := start
	finish.PublicationKey = "finish"
	finish.Type = "step_finished"
	f.append(t, finish)
	o := ObservationInput{Common: f.input.Common, Type: "proposal", Targets: Targets{EventIDs: []string{first.EventID}, StepIDs: []string{"step"}}, Finding: "finding", Action: ptr("act"), Owner: ptr("owner"), NextSignal: ptr("signal")}
	o.Kind = "ply.workspace.task-journal-observation-input"
	o.PublicationKey = "observation"
	if _, e := f.s.Append("task", f.file(t, o), "observation"); e != nil {
		t.Fatal(e)
	}
	s := f.show(t)
	if e := validateSnapshot(s); e != nil {
		t.Fatal(e)
	}
	return s
}

func TestValidateSnapshotErrorPaths(t *testing.T) {
	base := structuralSnapshot(t)
	contribution := func(s *Snapshot, typ string) *Event {
		for i := range s.Events {
			if s.Events[i].Type == typ {
				return &s.Events[i]
			}
		}
		t.Fatalf("no %s event", typ)
		return nil
	}
	cases := map[string]struct {
		edit func(*Snapshot)
		want string
	}{
		"workspace":         {func(s *Snapshot) { s.Workspace.ID = "ws_other" }, "invalid workspace or coverage"},
		"coverage":          {func(s *Snapshot) { s.Coverage.State = "full" }, "invalid workspace or coverage"},
		"source status":     {func(s *Snapshot) { s.Sources[0].Status = "fine" }, "invalid source"},
		"duplicate source":  {func(s *Snapshot) { s.Sources = append(s.Sources, s.Sources[0]) }, "invalid source"},
		"lane":              {func(s *Snapshot) { s.Lanes[0].LaneID = "lane_x" }, "invalid lane"},
		"event origin":      {func(s *Snapshot) { s.Events[0].Origin = "other" }, "invalid event"},
		"event time":        {func(s *Snapshot) { contribution(s, "step_started").OccurredAt = ptr("bad") }, "invalid event time"},
		"normalized time":   {func(s *Snapshot) { contribution(s, "step_started").OccurredAt = ptr("2026-10-03T11:00:00Z") }, "normalized event time differs"},
		"recorded time":     {func(s *Snapshot) { contribution(s, "step_started").RecordedAt = ptr("bad") }, "invalid recorded time"},
		"no source":         {func(s *Snapshot) { s.Events[0].SourceIDs = []string{} }, "event has no source"},
		"unknown source":    {func(s *Snapshot) { s.Events[0].SourceIDs = []string{"src_missing"} }, "unknown event source"},
		"step state":        {func(s *Snapshot) { s.Steps[0].State = "open" }, "invalid step"},
		"step start":        {func(s *Snapshot) { s.Steps[0].StartEventID = "missing" }, "unknown step start"},
		"step finish":       {func(s *Snapshot) { s.Steps[0].EndEventID = ptr("missing") }, "unknown step finish"},
		"duration":          {func(s *Snapshot) { s.Steps[0].State = "time_unknown" }, "invalid duration"},
		"parent":            {func(s *Snapshot) { s.Steps[0].ParentStepID = ptr("missing") }, "unknown parent step"},
		"event step":        {func(s *Snapshot) { contribution(s, "task_state").StepID = ptr("missing") }, "unknown event step"},
		"relation":          {func(s *Snapshot) { s.Events[0].Relations = []Relation{{"responds_to", "missing"}} }, "invalid event relation"},
		"superseding":       {func(s *Snapshot) { s.Events[0].SupersededBy = []string{"missing"} }, "unknown superseding event"},
		"observation event": {func(s *Snapshot) { s.Observations[0].EventID = "missing" }, "unknown observation event"},
		"proposal":          {func(s *Snapshot) { s.Observations[0].ProposalEventID = ptr(s.Steps[0].StartEventID) }, "invalid observation proposal"},
		"target":            {func(s *Snapshot) { s.Observations[0].Targets.EventIDs = []string{"missing"} }, "unknown observation target"},
		"target step":       {func(s *Snapshot) { s.Observations[0].Targets.StepIDs = []string{"missing"} }, "unknown observation step"},
		"axes count":        {func(s *Snapshot) { delete(s.Current.Axes, "human") }, "invalid current axes"},
		"axis status":       {func(s *Snapshot) { a := s.Current.Axes["human"]; a.BasisStatus = "x"; s.Current.Axes["human"] = a }, "invalid axis"},
		"axis event": {func(s *Snapshot) {
			a := s.Current.Axes["human"]
			a.EventIDs = []string{"missing"}
			s.Current.Axes["human"] = a
		}, "unknown axis event"},
		"selection":             {func(s *Snapshot) { s.Selection.EventIDs = []string{"missing"} }, "unknown snapshot event reference"},
		"coverage source":       {func(s *Snapshot) { s.Coverage.Reasons = []Reason{{"x", []string{"src_missing"}, "d"}} }, "unknown coverage source"},
		"duplicate observation": {func(s *Snapshot) { s.Observations = append(s.Observations, s.Observations[0]) }, "duplicate observation"},
		"lane differs":          {func(s *Snapshot) { contribution(s, "step_started").Actor.ID = "other" }, "event lane differs from actor/run"},
		"relation type": {func(s *Snapshot) {
			contribution(s, "step_finished").Relations = []Relation{{"likes", contribution(s, "step_started").EventID}}
		}, "unknown relation type"},
		"contribution":     {func(s *Snapshot) { contribution(s, "step_started").SourceSequence = nil }, "contribution projection differs from input"},
		"event projection": {func(s *Snapshot) { contribution(s, "step_started").Title = "other" }, "event projection differs from input"},
		"observation proj": {func(s *Snapshot) { s.Observations[0].Finding = "other" }, "observation projection differs"},
		"native data":      {func(s *Snapshot) { contribution(s, "task_state").Data = raw(map[string]any{"x": 1}) }, "unknown native data shape"},
		"supersession": {func(s *Snapshot) {
			contribution(s, "step_started").SupersededBy = []string{contribution(s, "step_finished").EventID}
		}, "supersession projection differs"},
		"step projection": {func(s *Snapshot) { s.Steps[0].Kind = "review" }, "step projection differs"},
		"step evidence":   {func(s *Snapshot) { s.Steps[0].State = "time_unknown"; s.Steps[0].DurationSeconds = nil }, "step duration or state differs from its evidence"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var s Snapshot
			if e := decode(raw(base), 64<<20, &s); e != nil {
				t.Fatal(e)
			}
			c.edit(&s)
			assertError(t, validateSnapshot(s), "task_journal_invalid_input: "+c.want)
		})
	}
}

func TestReadStoreErrorPaths(t *testing.T) {
	for name, c := range map[string]struct {
		edit func(t *testing.T, dir string)
		want string
	}{
		"gap":     {func(t *testing.T, dir string) { write(t, filepath.Join(dir, "000000004.json"), []byte("{}")) }, "task_journal_conflict: journal sequence gap or unexpected file: 000000004.json"},
		"private": {func(t *testing.T, dir string) { _ = os.Chmod(filepath.Join(dir, "000000001.json"), 0644) }, "task_journal_conflict: record must be private"},
		"chain": {func(t *testing.T, dir string) {
			b, _ := os.ReadFile(filepath.Join(dir, "000000002.json"))
			var r Record
			if e := decode(b, 128<<10, &r); e != nil {
				t.Fatal(e)
			}
			r.PreviousSHA256 = ptr(strings.Repeat("0", 64))
			b, _ = Canonical(r)
			write(t, filepath.Join(dir, "000000002.json"), b)
		}, "task_journal_conflict: journal record binding or chain differs"},
	} {
		t.Run(name, func(t *testing.T) {
			f := fixture(t)
			f.append(t, f.input)
			in := f.input
			in.PublicationKey = "second"
			f.append(t, in)
			ws := f.show(t).Workspace.ID
			c.edit(t, storePath(f.root, "task"))
			records, e := readStore(f.root, "task", ws)
			assertError(t, e, c.want)
			if name == "chain" && len(records) != 1 {
				t.Fatalf("validated prefix = %d records, want 1", len(records))
			}
		})
	}
}
