package workflownotification

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

const agentEventKind = "PlyAgentNotificationEvent@1"
const compactEventKind = "PlyAgentNotificationEvent@2"

type Event struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	OccurredAt string  `json:"occurred_at"`
	Phase      string  `json:"phase"`
	Record     Locator `json:"record"`
}

type eventRecord struct {
	Kind       string `json:"kind"`
	Activity   string `json:"activity"`
	Run        string `json:"run"`
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	Phase      string `json:"phase"`
	OccurredAt string `json:"occurred_at"`
	Public     Public `json:"public"`
}

func parseAgentRequest(r Request, m map[string]any, secret string, fresh bool) (Request, error) {
	fields := []string{"kind", "schema_version", "route", "source", "event", "sender", "public"}
	publicFields := []string{"task_title", "summary", "next_action", "next_actor"}
	eventKind := agentEventKind
	if r.SchemaVersion == 3 {
		fields = append(fields, "presentation")
		publicFields = append(publicFields, "status")
		eventKind = compactEventKind
	}
	if !keys(m, fields...) ||
		!keys(m["event"], "id", "type", "occurred_at", "phase", "record") ||
		!keys(m["sender"], "actor_claim") || !keys(m["public"], publicFields...) || r.Event == nil {
		return r, invalid()
	}
	source, ok := m["source"].(map[string]any)
	if !ok {
		return r, invalid()
	}
	fields = []string{"kind", "activity", "run", "worktree", "handoff"}
	for _, name := range []string{"start_receipt", "report"} {
		if value, exists := source[name]; exists {
			fields = append(fields, name)
			if !keys(value, "path", "sha256", "kind") {
				return r, invalid()
			}
		}
	}
	event := m["event"].(map[string]any)
	if !keys(source, fields...) || !keys(source["handoff"], "path", "sha256") || !keys(event["record"], "path", "sha256", "kind") {
		return r, invalid()
	}
	e := r.Event
	if r.Kind != "ply.workflow.notification-request" || r.Source.Kind != "external" || !routePattern.MatchString(r.Route) ||
		!cleanText(r.Source.Activity, 200) || !cleanText(r.Source.Run, 200) || !validPath(r.Source.Worktree) ||
		!validLocator(r.Source.Handoff) || !validLocator(e.Record) || e.Record.Kind == nil || *e.Record.Kind != eventKind ||
		!cleanText(e.ID, 80) || !cleanText(r.Sender.ActorClaim, 100) || !cleanText(r.Public.TaskTitle, 100) ||
		!cleanText(r.Public.Summary, 300) || !cleanText(r.Public.NextAction, 300) ||
		(r.Public.NextActor != "user" && r.Public.NextActor != "coordinator") {
		return r, invalid()
	}
	for _, loc := range []*Locator{r.Source.Start, r.Source.Report} {
		if loc != nil && (!validLocator(*loc) || loc.Kind == nil) {
			return r, invalid()
		}
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", e.OccurredAt)
	if err != nil || t.Format("2006-01-02T15:04:05Z") != e.OccurredAt {
		return r, invalid()
	}
	if e.Phase != "before_start" && e.Phase != "after_start" {
		return r, invalid()
	}
	if e.Phase == "before_start" && (e.Type != "agent_stopped" || r.Source.Start != nil) {
		return r, invalid()
	}
	switch e.Type {
	case "agent_finished":
		if r.Source.Report == nil {
			return r, invalid()
		}
	case "agent_stopped":
	case "feedback_required":
		// Whether the answer is necessary is a caller claim. Require an explicit
		// question as well as the user continuation, without inferring product state.
		if e.Phase != "after_start" || r.Public.NextActor != "user" || !strings.Contains(r.Public.Summary, "?") {
			return r, invalid()
		}
	default:
		return r, invalid()
	}
	for _, s := range []string{r.Public.TaskTitle, r.Public.Summary, r.Public.NextAction, r.Source.Run, e.ID, e.OccurredAt} {
		if !publicText(s) {
			return r, invalid()
		}
	}
	if !privateSafe(m, secret) {
		return r, invalid()
	}
	if r.SchemaVersion == 3 && (!keys(m["presentation"], "provider", "origin_cwd", "context_root", "context", "timezone", "local_occurred_at") || !validCompact(r, fresh)) {
		return r, invalid()
	}
	return r, nil
}

func checkEventRecord(raw []byte, r Request) error {
	var actual eventRecord
	m, err := strict(raw, &actual)
	publicFields := []string{"task_title", "summary", "next_action", "next_actor"}
	eventKind := agentEventKind
	if r.SchemaVersion == 3 {
		publicFields = append(publicFields, "status")
		eventKind = compactEventKind
	}
	if err != nil || !keys(m, "kind", "activity", "run", "event_id", "event_type", "phase", "occurred_at", "public") ||
		!keys(m["public"], publicFields...) {
		return invalid()
	}
	want := eventRecord{eventKind, r.Source.Activity, r.Source.Run, r.Event.ID, r.Event.Type, r.Event.Phase, r.Event.OccurredAt, r.Public}
	if !reflect.DeepEqual(actual, want) {
		return fail(2, "source_changed", "Event record does not match the request.")
	}
	return nil
}

type sourceFile struct {
	loc   Locator
	limit int64
}

func sourceLocators(r Request) []sourceFile {
	files := []sourceFile{{r.Source.Handoff, 1 << 20}}
	if r.Source.Start != nil {
		files = append(files, sourceFile{*r.Source.Start, 16 << 20})
	}
	if r.Source.Report != nil {
		files = append(files, sourceFile{*r.Source.Report, 16 << 20})
	}
	if r.Event != nil {
		files = append(files, sourceFile{r.Event.Record, 1 << 20})
	}
	return files
}

func agentMessage(r Request) Payload {
	heading := map[string]string{"agent_finished": "Agent finished — reported, not yet controlled", "agent_stopped": "Agent stopped", "feedback_required": "Your answer is needed"}[r.Event.Type]
	next := "User continues in the active task"
	if r.Public.NextActor == "coordinator" {
		next = "Coordinator continues"
	}
	return Payload{Text: fmt.Sprintf("%s\n%s\n%s\n%s: %s\nRun: %s · Event: %s\nOccurred: %s", heading, r.Public.TaskTitle, r.Public.Summary, next, r.Public.NextAction, r.Source.Run, r.Event.ID, r.Event.OccurredAt), Parse: "none"}
}
