package taskjournal

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const inputLimit = 64 << 10

var keyPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,128}$`)
var hexHash = regexp.MustCompile(`^[0-9a-f]{64}$`)
var oidPattern = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)

func hash(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func Canonical(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	x, e := decodeJournalJSON(b)
	if e != nil {
		return nil, e
	}
	return canonicaljson.Marshal(x)
}
func digest(v any) string { b, _ := Canonical(v); return hash(b) }
func equal(a, b any) bool {
	x, e := Canonical(a)
	y, f := Canonical(b)
	return e == nil && f == nil && bytes.Equal(x, y)
}
func raw(v any) json.RawMessage { b, _ := Canonical(v); return b }
func invalid(s string) error    { return fmt.Errorf("task_journal_invalid_input: %s", s) }
func conflict(s string) error   { return fmt.Errorf("task_journal_conflict: %s", s) }

// Require every typed field, including explicit nulls, recursively. Canonical
// decoding rejects duplicate keys, invalid UTF-8, extra values and unsafe numbers.
func shape(v any, t reflect.Type) error {
	if t == reflect.TypeOf(json.RawMessage{}) {
		if v == nil {
			return invalid("null data")
		}
		return nil
	}
	if t.Kind() == reflect.Pointer {
		if v == nil {
			return nil
		}
		return shape(v, t.Elem())
	}
	if v == nil {
		return invalid("unexpected null")
	}
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok {
			return invalid("expected object")
		}
		// Fields tagged omitempty were added later and may be absent; every
		// other field is required, also as an explicit null or empty list.
		fields, optional := map[string]reflect.Type{}, map[string]bool{}
		var add func(reflect.Type)
		add = func(t reflect.Type) {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Anonymous {
					add(f.Type)
				} else {
					tag := strings.Split(f.Tag.Get("json"), ",")
					if tag[0] != "-" {
						fields[tag[0]] = f.Type
						optional[tag[0]] = len(tag) > 1 && tag[1] == "omitempty"
					}
				}
			}
		}
		add(t)
		for n := range m {
			if _, ok := fields[n]; !ok {
				return invalid("missing or unknown fields in " + t.Name())
			}
		}
		for n, ft := range fields {
			x, ok := m[n]
			if !ok && optional[n] {
				continue
			}
			if !ok {
				return invalid("missing or unknown fields in " + t.Name())
			}
			if e := shape(x, ft); e != nil {
				return fmt.Errorf("%s: %w", n, e)
			}
		}
	case reflect.Slice:
		a, ok := v.([]any)
		if !ok {
			return invalid("expected array")
		}
		for _, x := range a {
			if e := shape(x, t.Elem()); e != nil {
				return e
			}
		}
	case reflect.Map:
		m, ok := v.(map[string]any)
		if !ok {
			return invalid("expected object")
		}
		for _, x := range m {
			if e := shape(x, t.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}
func decode(b []byte, limit int, out any) error {
	if len(b) > limit {
		return invalid("document exceeds size limit")
	}
	x, e := decodeJournalJSON(b)
	if e != nil {
		return invalid(e.Error())
	}
	c, e := canonicaljson.Marshal(x)
	if e != nil {
		return e
	}
	var v any
	if e = json.Unmarshal(c, &v); e != nil {
		return e
	}
	if e = shape(v, reflect.TypeOf(out).Elem()); e != nil {
		return e
	}
	d := json.NewDecoder(bytes.NewReader(c))
	d.DisallowUnknownFields()
	return d.Decode(out)
}
func textOK(s string, max int, empty bool) bool {
	if !utf8.ValidString(s) || (!empty && strings.TrimSpace(s) == "") || utf8.RuneCountInString(s) > max {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func one(s string, values ...string) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}
func optionalText(s *string, max int) bool  { return s == nil || textOK(*s, max, false) }
func parseTime(s string) (time.Time, error) { return time.Parse(time.RFC3339Nano, s) }
func normalize(s *string) *string {
	if s == nil {
		return nil
	}
	t, e := parseTime(*s)
	if e != nil {
		return nil
	}
	return ptr(t.UTC().Format(time.RFC3339Nano))
}
func validateTime(at *string, b TimeBasis) error {
	if !one(b.Kind, "observed", "reported", "unknown") || !one(b.Precision, "second", "millisecond", "microsecond", "nanosecond", "unknown") {
		return invalid("invalid time basis")
	}
	if b.Kind == "unknown" {
		if at != nil || b.Clock != nil || b.Uncertainty != nil || b.Precision != "unknown" {
			return invalid("unknown time must remain explicitly unknown")
		}
		return nil
	}
	if b.Clock == nil || !textOK(*b.Clock, 256, false) || b.Precision == "unknown" || (at == nil && b.Uncertainty == nil) {
		return invalid("known time needs a clock, precision and point or interval")
	}
	var point time.Time
	var e error
	if at != nil {
		point, e = parseTime(*at)
		if e != nil {
			return invalid("invalid occurred_at")
		}
	}
	if b.Uncertainty != nil {
		a, e := parseTime(b.Uncertainty.Earliest)
		z, f := parseTime(b.Uncertainty.Latest)
		if e != nil || f != nil || z.Before(a) || (at != nil && (point.Before(a) || point.After(z))) {
			return invalid("invalid uncertainty interval")
		}
	}
	return nil
}

// Roles are the actor roles a journal input may claim; Ply's own records use
// plyRole.
var Roles = []string{"human", "planner", "developer", "reviewer", "observer", "tool"}

const plyRole = "ply"

// StepKinds are the kinds of a work step.
var StepKinds = []string{"clarification", "planning", "implementation", "verification", "review", "correction", "coordination", "waiting"}

// PlyOperations are the operations Ply records in a Task journal.
var PlyOperations = []string{"plan_link", "task_create", "goal_register", "worktree_create", "goal_start_reserved", "goal_start_started", "goal_start_prompt_submitted", "integration", "lifecycle", "epic_integrate"}

var planIDPattern = regexp.MustCompile(`^PLAN-[1-9][0-9]{0,8}$`)

func validateCommon(c Common, kind, task string) error {
	if c.Kind != "ply.workspace.task-journal-"+kind+"-input" || c.SchemaVersion != 1 || c.TaskID != task || !keyPattern.MatchString(c.PublicationKey) {
		return invalid("input kind, version, Task or publication key differs")
	}
	switch {
	case !textOK(c.Actor.ID, 256, false):
		return invalid("actor.id must be one line of 1 to 256 characters")
	case !one(c.Actor.Role, append(Roles, plyRole)...):
		return invalid("actor.role must be one of " + strings.Join(Roles, ", "))
	case !optionalText(c.Actor.SessionID, 256):
		return invalid("actor.session_id must be null or one line of 1 to 256 characters")
	case !optionalText(c.ActivityID, 256):
		return invalid("activity_id must be null or one line of 1 to 256 characters")
	}
	if e := validateTime(c.OccurredAt, c.TimeBasis); e != nil {
		return e
	}
	for _, s := range c.Sources {
		if !filepath.IsAbs(s.Locator) || filepath.Clean(s.Locator) != s.Locator || !hexHash.MatchString(s.SHA256) {
			return invalid("source requires an absolute physical locator and lowercase SHA-256")
		}
	}
	if r := c.RunBinding; r != nil {
		if !regexp.MustCompile(`^trn_[0-9a-f]{64}$`).MatchString(r.RunID) || !hexHash.MatchString(strings.TrimPrefix(r.RequestSHA256, "sha256:")) || !hexHash.MatchString(strings.TrimPrefix(r.PreparationSHA256, "sha256:")) || !textOK(r.PreparationID, 256, false) {
			return invalid("invalid native run binding")
		}
	}
	return nil
}
func parseInput(b []byte, kind, task string) (Common, *EventInput, *ObservationInput, []byte, error) {
	var c Common
	var ev *EventInput
	var ob *ObservationInput
	if kind == "event" {
		ev = &EventInput{}
		if e := decode(b, inputLimit, ev); e != nil {
			return c, nil, nil, nil, e
		}
		c = ev.Common
	} else if kind == "observation" {
		ob = &ObservationInput{}
		if e := decode(b, inputLimit, ob); e != nil {
			return c, nil, nil, nil, e
		}
		c = ob.Common
	} else {
		return c, nil, nil, nil, invalid("unknown input type")
	}
	if e := validateCommon(c, kind, task); e != nil {
		return c, nil, nil, nil, e
	}
	ply := c.Actor.Role == plyRole
	if ply != strings.HasPrefix(c.PublicationKey, "ply:") || ply != (ev != nil && ev.Type == "ply_action") || ev != nil && (ev.Ply != nil) != (ev.Type == "ply_action") {
		return c, nil, nil, nil, invalid("role ply, type ply_action, the ply field and a ply: publication key belong together; such records are written by Ply")
	}
	if len(c.Sources) == 0 && (ev == nil || !one(ev.Type, "note", "step_started", "step_finished", "ply_action")) {
		return c, nil, nil, nil, invalid("at least one source is required for a decision or an observation")
	}
	var value any
	var e error
	if ev != nil {
		value, e = ev, validateEventInput(ev)
	} else {
		value, e = ob, validateObservationInput(ob)
	}
	if e != nil {
		return c, nil, nil, nil, e
	}
	canonical, e := Canonical(value)
	return c, ev, ob, canonical, e
}
func validateEventInput(ev *EventInput) error {
	if !one(ev.Type, "step_started", "step_finished", "note", "decision", "ply_action") || !textOK(ev.Title, 256, false) || ev.Detail != "" {
		return invalid("invalid event type, title or detail (detail must be empty in v1)")
	}
	if p := ev.Ply; p != nil {
		if !one(p.Operation, PlyOperations...) || p.Plan != nil && !planIDPattern.MatchString(*p.Plan) || p.PlanStep != nil && (p.Plan == nil || !textOK(*p.PlanStep, 64, false)) || len(p.Refs) > 16 {
			return invalid("invalid ply action")
		}
		for k, v := range p.Refs {
			if !textOK(k, 64, false) || !textOK(v, 256, false) {
				return invalid("invalid ply action reference")
			}
		}
	}
	if ev.Outcome != nil && !one(*ev.Outcome, "pass", "fail", "complete", "blocked", "unknown", "deferred", "declined") {
		return invalid("invalid reported outcome")
	}
	if one(ev.Type, "step_started", "step_finished") && ev.Step == nil {
		return invalid("start and finish require a step")
	}
	if ev.Type == "decision" && ev.Step != nil {
		return invalid("decision is a point, not a step")
	}
	if s := ev.Step; s != nil {
		switch {
		case !textOK(s.ID, 256, false):
			return invalid("step.id must be one line of 1 to 256 characters")
		case !optionalText(s.ParentStepID, 256):
			return invalid("step.parent_step_id must be null or one line of 1 to 256 characters")
		case !one(s.Kind, StepKinds...):
			return invalid("step.kind must be one of " + strings.Join(StepKinds, ", "))
		}
		if s.Kind == "waiting" && ev.Waiting == nil {
			return invalid("waiting step needs explicit reason, dependency and next actor")
		}
	}
	if w := ev.Waiting; w != nil {
		if ev.Step == nil || ev.Step.Kind != "waiting" || !textOK(w.Reason, 2048, false) || !textOK(w.Dependency, 256, false) || !textOK(w.NextActor, 256, false) {
			return invalid("invalid waiting metadata")
		}
	}
	if p := ev.Candidate; p != nil {
		if !textOK(p.RepoID, 256, false) || !textOK(p.WorktreeID, 256, false) || !oidPattern.MatchString(p.OID) || !oidPattern.MatchString(p.Tree) {
			return invalid("invalid candidate")
		}
	}
	for _, r := range ev.Relations {
		if !one(r.Type, "responds_to", "corrects", "verifies", "supersedes") || !textOK(r.EventID, 256, false) {
			return invalid("invalid relation")
		}
		if r.Type == "supersedes" && !one(ev.Type, "note", "decision") {
			return invalid("only notes and decisions supersede contributions")
		}
	}
	return nil
}
func validateObservationInput(ob *ObservationInput) error {
	if !one(ob.Type, "proposal", "decision", "applied", "assessment") || !textOK(ob.Finding, 2048, false) || !optionalText(ob.Hypothesis, 2048) || !optionalText(ob.Action, 2048) || !optionalText(ob.Owner, 256) || !optionalText(ob.NextSignal, 2048) {
		return invalid("invalid observation")
	}
	if ob.Type == "proposal" {
		if ob.ProposalEventID != nil || ob.Action == nil || ob.Owner == nil || ob.NextSignal == nil {
			return invalid("proposal requires action, owner and next signal")
		}
	} else if ob.ProposalEventID == nil || !textOK(*ob.ProposalEventID, 256, false) {
		return invalid("follow-up requires a proposal")
	}
	if ob.Type == "decision" {
		if ob.Decision == nil || !one(*ob.Decision, "accepted", "rejected") {
			return invalid("invalid decision")
		}
	} else if ob.Decision != nil {
		return invalid("decision only belongs to decision records")
	}
	if ob.Type == "assessment" {
		if ob.Assessment == nil || !one(*ob.Assessment, "better", "unchanged", "worse", "inconclusive") {
			return invalid("invalid assessment")
		}
	} else if ob.Assessment != nil {
		return invalid("assessment only belongs to assessment records")
	}
	if len(ob.Targets.EventIDs)+len(ob.Targets.StepIDs) == 0 && ob.Targets.Interval == nil {
		return invalid("observation needs a target")
	}
	if i := ob.Targets.Interval; i != nil {
		a, e := parseTime(i.From)
		z, f := parseTime(i.To)
		if e != nil || f != nil || z.Before(a) || len(i.ActorIDs)+len(i.RunIDs) == 0 {
			return invalid("invalid or unscoped target interval")
		}
	}
	return nil
}
