package workflownotification

import (
	"strings"
	"time"
)

// Request v4 is the oldest request this code sends: a public event without
// source documents; its local request is the only record.
const documentlessSchema = 4

// Request v5 is v4 with the proven Herdr workspace and tab display labels.
const herdrSchema = 5

// Request v6 is v4, with an optional v5 Herdr origin, plus one owner-selected app link.
const appLinkSchema = 6

type Event struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	OccurredAt string `json:"occurred_at"`
	Phase      string `json:"phase"`
}

func agentSchema(version int) bool { return version >= documentlessSchema && version <= appLinkSchema }

func parseAgentRequest(r Request, m map[string]any, secret string, fresh bool) (Request, error) {
	if !keys(m, "kind", "schema_version", "route", "source", "event", "sender", "public", "presentation") ||
		!keys(m["event"], "id", "type", "occurred_at", "phase") ||
		!keys(m["sender"], "actor_claim") || !keys(m["public"], "task_title", "summary", "next_action", "next_actor", "status") ||
		!keys(m["source"], "kind", "activity", "run", "worktree") || r.Event == nil {
		return r, invalid()
	}
	e := r.Event
	if r.Kind != "ply.workflow.notification-request" || r.Source.Kind != "external" || !routePattern.MatchString(r.Route) ||
		!cleanText(r.Source.Activity, 200) || !cleanText(r.Source.Run, 200) || !validPath(r.Source.Worktree) ||
		!cleanText(e.ID, 80) || !cleanText(r.Sender.ActorClaim, 100) || !cleanText(r.Public.TaskTitle, 100) ||
		!cleanText(r.Public.Summary, 300) || !cleanText(r.Public.NextAction, 300) ||
		(r.Public.NextActor != "user" && r.Public.NextActor != "coordinator") {
		return r, invalid()
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", e.OccurredAt)
	if err != nil || t.Format("2006-01-02T15:04:05Z") != e.OccurredAt {
		return r, invalid()
	}
	if e.Phase != "before_start" && e.Phase != "after_start" {
		return r, invalid()
	}
	if e.Phase == "before_start" && e.Type != "agent_stopped" {
		return r, invalid()
	}
	switch e.Type {
	case "agent_finished", "agent_stopped":
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
	presentationFields := []string{"provider", "origin_cwd", "context_root", "context", "timezone", "local_occurred_at"}
	presentation, _ := m["presentation"].(map[string]any)
	if _, herdr := presentation["herdr"]; r.SchemaVersion == herdrSchema || r.SchemaVersion == appLinkSchema && herdr {
		presentationFields = append(presentationFields, "herdr")
		if !keys(presentation["herdr"], "workspace", "tab") {
			return r, invalid()
		}
	}
	if r.SchemaVersion == appLinkSchema {
		presentationFields = append(presentationFields, "app_link")
	}
	if !keys(m["presentation"], presentationFields...) || !validCompact(r, fresh) {
		return r, invalid()
	}
	return r, nil
}
