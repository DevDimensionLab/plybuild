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
// It is absent from v1/v2 JSON to preserve their persisted integrity digests.
type Presentation struct {
	Provider        string `json:"provider"`
	OriginCWD       string `json:"origin_cwd"`
	ContextRoot     string `json:"context_root"`
	Context         string `json:"context"`
	Timezone        string `json:"timezone"`
	LocalOccurredAt string `json:"local_occurred_at"`
}

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
		!validPath(p.OriginCWD) || !validPath(p.ContextRoot) || !compactText(p.Context, 96) ||
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
	return Payload{Text: fmt.Sprintf("%s · %s · %s\n%s: %s%s\nNext: %s", p.LocalOccurredAt[11:19], p.Provider, p.Context, statusLabels[r.Public.Status], r.Public.TaskTitle, summary, r.Public.NextAction), Parse: "none"}
}
