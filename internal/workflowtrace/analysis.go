package workflowtrace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
)

func (b *builder) measure(id, kind, definition, unit string, v *float64, ids []string, unknowns ...string) {
	ids = unique(ids)
	sources := []string{}
	coverage := "exact"
	for _, id := range ids {
		if i, ok := b.events[id]; ok {
			e := b.out.Events[i]
			sources = append(sources, e.SourceIDs...)
			if !b.validEvent(e) {
				coverage = "unknown"
			}
		} else {
			coverage = "unknown"
		}
	}
	if v == nil {
		coverage = "unknown"
	} else if coverage == "exact" && kind == "count" && b.out.Coverage.State != "complete" {
		coverage = "lower_bound"
	}
	if coverage == "unknown" && v != nil {
		v = nil
		unknowns = append(unknowns, "At least one evidence endpoint is unavailable or changed; a numeric value cannot be supported.")
	}
	b.out.Analysis = append(b.out.Analysis, Analysis{ID: id, Kind: kind, Definition: definition, Value: v, Unit: unit, Coverage: coverage, EvidenceIDs: ids, SourceIDs: unique(sources), Unknowns: unique(unknowns)})
}

func (b *builder) count(id, definition string, n int, ids []string) {
	b.measure(id, "count", definition, "records", ptr(float64(n)), ids, "Counts cover preserved records, not all real-world actions or work.")
}

func (b *builder) pattern(id, definition, question string, ids []string, missing ...string) {
	b.measure(id, "pattern", definition, "question", nil, ids, missing...)
	p := &b.out.Analysis[len(b.out.Analysis)-1]
	p.Question = ptr(question)
	// A question can have exact evidence while its answer remains unknown.
	p.Coverage = "exact"
	for _, id := range p.EvidenceIDs {
		if i, ok := b.events[id]; !ok || !b.validEvent(b.out.Events[i]) {
			p.Coverage = "unknown"
		}
	}
}

func (b *builder) analyze() {
	questions, handoffs := []string{}, []string{}
	candidateEvents := map[string][]string{}
	questionGroups := map[string][]string{}
	for _, e := range b.out.Events {
		if !b.validEvent(e) {
			continue
		}
		if e.CandidateID != nil {
			candidateEvents[*e.CandidateID] = append(candidateEvents[*e.CandidateID], e.ID)
		}
		if question := eventQuestion(e); question != "" {
			questions = append(questions, e.ID)
			// This groups equal recorded text, never authenticates or merges an action.
			questionGroups[question] = append(questionGroups[question], e.ID)
			b.pattern("question/"+e.ID, "A recorded request for an answer; later working reports do not establish a response.", "What evidence records the answer to this request?", []string{e.ID}, "An explicit answer link and compatible event clocks are required to measure a response interval.")
		}
	}
	for _, r := range b.out.Relations {
		if r.Type == "handoff" {
			handoffs = append(handoffs, r.From, r.To)
		}
	}
	ce := []string{}
	for _, ids := range candidateEvents {
		ce = append(ce, ids...)
	}
	b.count("distinct_observed_candidates", "Distinct exact repository/commit/tree pairs in validated recorded events; no live Git verification and no inferred correction.", len(candidateEvents), ce)
	b.count("question_records", "Distinct records containing explicit questions; equal wording remains separate records and need not be distinct real actions.", len(questions), questions)
	handoffCount := 0
	for _, r := range b.out.Relations {
		if r.Type == "handoff" {
			if i, ok := b.events[r.From]; ok && b.validEvent(b.out.Events[i]) {
				handoffCount++
			}
		}
	}
	b.count("explicit_handoff_links", "Explicit source handoff relations only; adjacent roles, report prose and timestamp proximity create no handoff.", handoffCount, handoffs)
	for question, ids := range questionGroups {
		if len(ids) > 1 {
			b.pattern("repeated_question/"+hash(question), "Several preserved records contain the same exact question text; they are not silently merged or counted as authenticated duplicate actions.", "Could this request remain available across handoffs without repeating the human-facing question?", ids, "The records do not establish how often a person saw the question, why it was repeated, or time savings.")
		}
	}
	b.journalIntervals()
}

func eventQuestion(e Event) string {
	var m map[string]json.RawMessage
	if json.Unmarshal(e.Data, &m) != nil {
		return ""
	}
	var question string
	_ = json.Unmarshal(m["question"], &question)
	if question != "" {
		return question
	}
	if e.Type == "decision_requested" || e.Type == "question" {
		return e.Title
	}
	return ""
}

func (b *builder) interval(id, definition, clock string, start, end *string, ids []string, compatible bool) {
	var elapsed *float64
	missing := []string{}
	if start == nil || end == nil {
		missing = append(missing, "A start or end is missing; no interval is inferred from observation time or file metadata.")
	} else if !compatible || clock == "" {
		missing = append(missing, "Time kinds, clocks or uncertainty do not permit an elapsed interval.")
	} else {
		a, ae := time.Parse(time.RFC3339Nano, *start)
		z, ze := time.Parse(time.RFC3339Nano, *end)
		if ae != nil || ze != nil || z.Before(a) {
			missing = append(missing, "Endpoints are invalid or reversed.")
			b.diagnostic("time_interval_invalid", "Interval endpoints are invalid or reversed: "+id, nil, ids)
		} else {
			elapsed = ptr(z.Sub(a).Seconds())
		}
	}
	b.measure(id, "interval", definition+" Clock: "+clock+"; endpoints: "+timeLabel(start)+" -> "+timeLabel(end)+".", "seconds", elapsed, ids, append(missing, "Elapsed time is not active work, attention, cost or measured improvement.")...)
}
func timeLabel(t *string) string {
	if t == nil {
		return "unknown"
	}
	return *t
}
func eventTime(e Event) *string {
	if e.TimeBasis.Kind == "reported" {
		return e.ReportedAtUTC
	}
	if e.TimeBasis.Kind == "observed" {
		return e.OccurredAtUTC
	}
	return nil
}
func compatible(a, z Event) bool {
	return a.TimeBasis.Kind == z.TimeBasis.Kind && (a.TimeBasis.Kind == "reported" || a.TimeBasis.Kind == "observed") && a.TimeBasis.Clock != nil && value(a.TimeBasis.Clock) != "" && value(a.TimeBasis.Clock) == value(z.TimeBasis.Clock) && a.TimeBasis.Uncertainty == nil && z.TimeBasis.Uncertainty == nil && a.TimeBasis.Precision != "unknown" && z.TimeBasis.Precision != "unknown"
}

func (b *builder) journalIntervals() {
	ends := map[string][]Event{}
	for _, e := range b.out.Events {
		if e.Type == "step_finished" {
			var in taskjournal.EventInput
			if json.Unmarshal(e.Data, &in) == nil && in.Step != nil {
				ends[in.Step.ID] = append(ends[in.Step.ID], e)
			}
		}
	}
	for _, e := range b.out.Events {
		if e.Type != "step_started" {
			continue
		}
		var in taskjournal.EventInput
		if json.Unmarshal(e.Data, &in) != nil || in.Step == nil || in.Waiting == nil {
			continue
		}
		ids := []string{e.ID}
		var end *string
		valid := false
		clock := value(e.TimeBasis.Clock)
		for _, z := range ends[in.Step.ID] {
			ids = append(ids, z.ID)
		}
		if len(ends[in.Step.ID]) == 1 {
			z := ends[in.Step.ID][0]
			end = eventTime(z)
			valid = compatible(e, z)
		}
		b.interval("reported_wait/"+e.ID, "Reported waiting step: "+in.Waiting.Reason+"; dependency: "+in.Waiting.Dependency+"; next actor: "+in.Waiting.NextActor, clock, eventTime(e), end, ids, valid)
	}
	for _, r := range b.out.Relations {
		if r.Type != "responds_to" {
			continue
		}
		i, ok := b.events[r.To]
		j, jok := b.events[r.From]
		if !ok || !jok {
			continue
		}
		start, end := b.out.Events[i], b.out.Events[j]
		b.interval("reported_response/"+r.From+"/"+r.To, "Recorded elapsed interval for an explicit responds_to link; it does not authenticate the answer or prove waiting.", value(start.TimeBasis.Clock), eventTime(start), eventTime(end), []string{start.ID, end.ID}, compatible(start, end))
	}
}

// ValidateReferences is useful to schema consumers and acceptance fixtures.
// It checks projection closure without interpreting text or granting authority.
func ValidateReferences(r Result) error {
	events, sources, candidates := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, s := range r.Sources {
		if sources[s.ID] {
			return fmt.Errorf("duplicate source %s", s.ID)
		}
		sources[s.ID] = true
	}
	for _, c := range r.Candidates {
		if candidates[c.ID] {
			return fmt.Errorf("duplicate candidate %s", c.ID)
		}
		candidates[c.ID] = true
	}
	for _, e := range r.Events {
		if events[e.ID] {
			return fmt.Errorf("duplicate event %s", e.ID)
		}
		events[e.ID] = true
		for _, s := range e.SourceIDs {
			if !sources[s] {
				return fmt.Errorf("event %s missing source %s", e.ID, s)
			}
		}
		if e.CandidateID != nil && !candidates[*e.CandidateID] {
			return fmt.Errorf("missing candidate %s", *e.CandidateID)
		}
	}
	for _, a := range r.Analysis {
		for _, id := range a.EvidenceIDs {
			if !events[id] {
				return fmt.Errorf("analysis %s missing event %s", a.ID, id)
			}
		}
		for _, id := range a.SourceIDs {
			if !sources[id] {
				return fmt.Errorf("analysis missing source %s", id)
			}
		}
	}
	for _, c := range r.Chains {
		for _, id := range c.EventIDs {
			if !events[id] {
				return fmt.Errorf("chain missing event %s", id)
			}
		}
	}
	for _, rel := range r.Relations {
		if !events[rel.From] || !events[rel.To] {
			return fmt.Errorf("relation %s lacks endpoint", rel.Type)
		}
	}
	return nil
}
