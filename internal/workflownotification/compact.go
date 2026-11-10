package workflownotification

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Presentation is a declared origin, separate from the target worktree and model.
type Presentation struct {
	Provider        string `json:"provider"`
	OriginCWD       string `json:"origin_cwd"`
	ContextRoot     string `json:"context_root"`
	Context         string `json:"context"`
	Timezone        string `json:"timezone"`
	LocalOccurredAt string `json:"local_occurred_at"`
	// Omitted below v5, so earlier request bytes and payloads are unchanged.
	Herdr *HerdrOrigin `json:"herdr,omitempty"`
	// Omitted below v6. The owner-selected link is frozen with the first request.
	AppLink string `json:"app_link,omitempty"`
}

// HerdrOrigin holds display labels only; Herdr IDs are never stored.
type HerdrOrigin struct {
	Workspace string `json:"workspace"`
	Tab       string `json:"tab"`
}

// One https host with an optional plain path: no port, user, query, fragment or Slack markup.
var appLink = regexp.MustCompile(`^https://[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+(/[A-Za-z0-9._~-]+)*/?$`)
var zoneName = regexp.MustCompile(`^[A-Za-z0-9_+-]+(/[A-Za-z0-9_+-]+)*$`)
var statusLabels = map[string]string{
	"ready_for_review":     "Ready for review",
	"ready_for_your_check": "Ready for your check",
	"done":                 "Done", "stopped": "Stopped", "needs_answer": "Needs answer",
}

func compactText(s string, max int) bool {
	if !utf8.ValidString(s) || !cleanText(s, max) || !publicText(s) {
		return false
	}
	for _, r := range s {
		if unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Zl, r) || unicode.Is(unicode.Zp, r) {
			return false
		}
	}
	return true
}

func validCompact(r Request, fresh bool) bool {
	p := r.Presentation
	if p == nil || (p.Provider != "codex" && p.Provider != "claude") ||
		!validPath(p.OriginCWD) || !validPath(p.ContextRoot) ||
		r.SchemaVersion != appLinkSchema && (r.SchemaVersion == herdrSchema) != (p.Herdr != nil) ||
		(r.SchemaVersion == appLinkSchema) != (p.AppLink != "") || p.AppLink != "" && (len(p.AppLink) > 200 || !appLink.MatchString(p.AppLink)) || p.Herdr != nil && (!compactText(p.Herdr.Workspace, 32) || !compactText(p.Herdr.Tab, 32)) ||
		!zoneName.MatchString(p.Timezone) || p.Timezone == "Local" || len(p.Timezone) > 128 ||
		!compactText(r.Public.TaskTitle, 60) || !compactText(r.Public.Summary, 100) || !compactText(r.Public.NextAction, 140) {
		return false
	}
	rel, err := filepath.Rel(p.ContextRoot, p.OriginCWD)
	if err != nil || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.ToSlash(rel) != p.Context {
		return false
	}
	switch r.Event.Type {
	case "agent_finished":
		if r.Public.Status != "ready_for_review" && r.Public.Status != "ready_for_your_check" && r.Public.Status != "done" {
			return false
		}
	case "agent_stopped":
		if r.Public.Status != "stopped" {
			return false
		}
	case "feedback_required":
		if r.Public.Status != "needs_answer" {
			return false
		}
	default:
		return false
	}
	utc, err := time.Parse(time.RFC3339, r.Event.OccurredAt)
	if err != nil {
		return false
	}
	local, err := time.Parse(time.RFC3339, p.LocalOccurredAt)
	if err != nil || local.Format(time.RFC3339) != p.LocalOccurredAt || !local.Equal(utc) {
		return false
	}
	if fresh {
		zone, err := time.LoadLocation(p.Timezone)
		if err != nil || utc.In(zone).Format(time.RFC3339) != p.LocalOccurredAt {
			return false
		}
		for _, path := range []string{p.ContextRoot, p.OriginCWD} {
			dir, err := openDirectory(path)
			if err != nil {
				return false
			}
			dir.Close()
		}
	}
	return utf8.RuneCountInString(compactMessage(r).Text) <= 480
}

func compactMessage(r Request) Payload {
	p := r.Presentation
	// Database integrity validation can reach rendering before request validation.
	if p == nil || len(p.LocalOccurredAt) < 19 {
		return Payload{Parse: "none"}
	}
	summary := ""
	if r.Public.Summary != r.Public.TaskTitle {
		summary = " — " + r.Public.Summary
	}
	// The directory name stays local.
	origin := p.LocalOccurredAt[11:19] + " · " + p.Provider
	if r.SchemaVersion >= herdrSchema && p.Herdr != nil {
		origin += " · " + p.Herdr.Workspace + " › " + p.Herdr.Tab
	}
	text := fmt.Sprintf("%s\n%s: %s%s\nNext: %s", origin, statusLabels[r.Public.Status], r.Public.TaskTitle, summary, r.Public.NextAction)
	payload := Payload{Text: text, Parse: "none"}
	if r.SchemaVersion == appLinkSchema {
		// Agent text stays plain_text; only Ply builds the mrkdwn link from the validated URL.
		payload.Blocks = []Block{{Type: "section", Text: &PlainText{Type: "plain_text", Text: text}},
			{Type: "context", Elements: []MrkdwnText{{Type: "mrkdwn", Text: "<" + p.AppLink + "|Åpne plyapp>", Verbatim: true}}}}
	}
	return payload
}
