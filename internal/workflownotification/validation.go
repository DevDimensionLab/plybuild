package workflownotification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var routePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
var envPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]{0,127}$`)
var webhookPath = regexp.MustCompile(`^/services/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+$`)
var webhookText = regexp.MustCompile(`(?i)https?://hooks\.(slack\.com|slack-gov\.com)(/|\b)`)
var idPattern = regexp.MustCompile(`^ntf_[0-9a-f]{64}$`)

func invalid() error {
	return fail(2, "invalid_input", "Invalid notification input or source binding.")
}
func strict(raw []byte, out any) (map[string]any, error) {
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return nil, invalid()
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var obj map[string]any
	if err := dec.Decode(&obj); err != nil || obj == nil {
		return nil, invalid()
	}
	if out != nil {
		d := json.NewDecoder(bytes.NewReader(raw))
		d.DisallowUnknownFields()
		if err := d.Decode(out); err != nil {
			return nil, invalid()
		}
	}
	return obj, nil
}
func keys(v any, names ...string) bool {
	m, ok := v.(map[string]any)
	if !ok || len(m) != len(names) {
		return false
	}
	for _, n := range names {
		if _, ok = m[n]; !ok {
			return false
		}
	}
	return true
}
func cleanText(s string, max int) bool {
	if s == "" || utf8.RuneCountInString(s) > max || strings.TrimSpace(s) != s {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
func publicText(s string) bool {
	l := strings.ToLower(s)
	return !strings.ContainsAny(s, "<>") && !strings.Contains(l, "@channel") && !strings.Contains(l, "@here") && !strings.Contains(l, "@everyone")
}
func privateSafe(v any, secret string) bool {
	switch t := v.(type) {
	case string:
		return !webhookText.MatchString(t) && (secret == "" || !strings.Contains(t, secret))
	case map[string]any:
		for k, x := range t {
			if !privateSafe(k, secret) || !privateSafe(x, secret) {
				return false
			}
		}
	case []any:
		for _, x := range t {
			if !privateSafe(x, secret) {
				return false
			}
		}
	}
	return true
}
func validPath(s string) bool {
	return filepath.IsAbs(s) && filepath.Clean(s) == s && cleanText(s, 16384)
}
func validLocator(l Locator) bool {
	return validPath(l.Path) && digestPattern.MatchString(l.SHA256) && (l.Kind == nil || cleanText(*l.Kind, 200))
}
func parseRequest(raw []byte, secret string) (Request, error) {
	return parseRequestMode(raw, secret, true)
}

// Historical requests freeze their local offset. Revalidating with today's
// tzdata would invalidate an otherwise intact database after a zone rule change.
func parseStoredRequest(raw []byte, secret string) (Request, error) {
	return parseRequestMode(raw, secret, false)
}
func parseRequestMode(raw []byte, secret string, fresh bool) (Request, error) {
	var r Request
	m, err := strict(raw, &r)
	if err != nil {
		return r, err
	}
	if r.SchemaVersion == 2 || r.SchemaVersion == 3 {
		return parseAgentRequest(r, m, secret, fresh)
	}
	if r.Gate == nil || r.Source.Start == nil || r.Source.Report == nil || !keys(m, "kind", "schema_version", "route", "source", "gate", "sender", "public") || !keys(m["gate"], "id", "revision", "opened_at", "reason", "state") || !keys(m["sender"], "actor_claim") || !keys(m["public"], "task_title", "next_action") {
		return r, invalid()
	}
	source, ok := m["source"].(map[string]any)
	if !ok {
		return r, invalid()
	}
	fields := []string{"kind", "activity", "run", "worktree", "handoff", "start_receipt", "report"}
	if r.Source.Kind == "native" {
		fields = append(fields, "handoff_id", "task_id")
	} else if r.Source.Kind != "external" {
		return r, invalid()
	}
	if !keys(source, fields...) || !keys(source["handoff"], "path", "sha256") {
		return r, invalid()
	}
	for _, name := range []string{"start_receipt", "report"} {
		f := []string{"path", "sha256"}
		if r.Source.Kind == "external" {
			f = append(f, "kind")
		}
		if !keys(source[name], f...) {
			return r, invalid()
		}
	}
	if r.Kind != "ply.workflow.notification-request" || r.SchemaVersion != 1 || !routePattern.MatchString(r.Route) || !cleanText(r.Source.Activity, 200) || !cleanText(r.Source.Run, 200) || !validPath(r.Source.Worktree) || !validLocator(r.Source.Handoff) || !validLocator(*r.Source.Start) || !validLocator(*r.Source.Report) || !cleanText(r.Gate.ID, 80) || r.Gate.Revision < 1 || r.Gate.Revision > 2147483647 || r.Gate.Reason != "result_control" || r.Gate.State != "waiting_for_human" || !cleanText(r.Sender.ActorClaim, 100) || !cleanText(r.Public.TaskTitle, 100) || !cleanText(r.Public.NextAction, 300) {
		return r, invalid()
	}
	t, e := time.Parse("2006-01-02T15:04:05Z", r.Gate.OpenedAt)
	if e != nil || t.Format("2006-01-02T15:04:05Z") != r.Gate.OpenedAt {
		return r, invalid()
	}
	if r.Source.Kind == "external" {
		if r.Source.Start.Kind == nil || r.Source.Report.Kind == nil {
			return r, invalid()
		}
	} else {
		if r.Source.HandoffID == nil || !cleanText(*r.Source.HandoffID, 200) {
			return r, invalid()
		}
		if source["task_id"] != nil {
			s, ok := source["task_id"].(string)
			if !ok || !cleanText(s, 200) {
				return r, invalid()
			}
		}
	}
	for _, s := range []string{r.Public.TaskTitle, r.Public.NextAction, r.Source.Run, r.Gate.ID, r.Gate.OpenedAt} {
		if !publicText(s) {
			return r, invalid()
		}
	}
	if !privateSafe(m, secret) {
		return r, invalid()
	}
	return r, nil
}
func parseRoute(raw []byte, secret string) (Route, error) {
	var r Route
	m, err := strict(raw, &r)
	if err != nil {
		return r, err
	}
	if !keys(m, "kind", "schema_version", "name", "transport", "channel_label", "webhook_env", "state_root") || r.Kind != "ply.workflow.notification-route" || r.SchemaVersion != 1 || !routePattern.MatchString(r.Name) || r.Transport != "slack_incoming_webhook" || !cleanText(r.ChannelLabel, 100) || !envPattern.MatchString(r.WebhookEnv) || !validPath(r.StateRoot) || !privateSafe(m, secret) {
		return r, invalid()
	}
	return r, nil
}
func credential(s string) bool {
	u, err := url.Parse(s)
	return err == nil && !strings.ContainsAny(s, "?#") && u.Scheme == "https" && u.Host == "hooks.slack.com" && u.User == nil && u.RawQuery == "" && !u.ForceQuery && u.Fragment == "" && u.RawFragment == "" && u.RawPath == "" && webhookPath.MatchString(u.Path)
}
func identity(r Request) string {
	agent := r.SchemaVersion == 2 || r.SchemaVersion == 3
	if agent && r.Event == nil || !agent && r.Gate == nil {
		return ""
	}
	if agent {
		return "ntf_" + strings.TrimPrefix(gateFamily(r), "sha256:")
	}
	return "ntf_" + strings.TrimPrefix(digestValue([]any{r.Source.Kind, r.Source.Activity, r.Source.Run, r.Gate.ID, r.Gate.Revision, r.Route}), "sha256:")
}
func gateFamily(r Request) string {
	if r.SchemaVersion == 2 || r.SchemaVersion == 3 {
		// This historical namespace is intentionally independent of new schema versions.
		return digestValue([]any{"agent-event-v2", r.Source.Activity, r.Source.Run, r.Event.ID, r.Route})
	}
	return digestValue([]any{r.Source.Kind, r.Source.Activity, r.Source.Run, r.Gate.ID, r.Route})
}
func message(r Request) Payload {
	if r.SchemaVersion == 3 {
		return compactMessage(r)
	}
	if r.SchemaVersion == 2 {
		return agentMessage(r)
	}
	return Payload{Text: fmt.Sprintf("Ply needs your attention\n%s — report ready\nReported only; not yet controlled. Human QA is not attested.\nNext: %s\nRun: %s · Gate: %s/%d\nOpened: %s", r.Public.TaskTitle, r.Public.NextAction, r.Source.Run, r.Gate.ID, r.Gate.Revision, r.Gate.OpenedAt), Parse: "none"}
}
func checkSources(r Request) error {
	d, err := openDirectory(r.Source.Worktree)
	if err != nil {
		return invalid()
	}
	d.Close()
	for _, item := range sourceLocators(r) {
		raw, err := readFile(item.loc.Path, item.limit, false)
		if err != nil {
			return err
		}
		if Digest(raw) != item.loc.SHA256 {
			return fail(2, "source_changed", "Source bytes do not match the bound digest.")
		}
		if (r.SchemaVersion == 2 || r.SchemaVersion == 3) && item.loc == r.Event.Record {
			if e := checkEventRecord(raw, r); e != nil {
				return e
			}
		} else if r.Source.Kind == "external" && item.loc.Kind != nil {
			m, e := strictSource(raw)
			if e != nil {
				return e
			}
			if m["kind"] != *item.loc.Kind || m["activity"] != r.Source.Activity || m["run"] != r.Source.Run {
				return fail(2, "source_changed", "Source identities do not match the request.")
			}
		}
	}
	if r.Source.Kind == "native" {
		return checkNative(r.Source)
	}
	return nil
}

// External metadata is ordinary JSON, not the native signed-integer canonical
// document domain. Validate its grammar first, then mask only numeric tokens to
// reuse the existing duplicate-key and Unicode checks without rejecting decimals
// or large JSON integers in fields the notification service deliberately ignores.
func strictSource(raw []byte) (map[string]any, error) {
	if !json.Valid(raw) {
		return nil, invalid()
	}
	normalized := make([]byte, 0, len(raw))
	inString := false
	escaped := false
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if inString {
			normalized = append(normalized, c)
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				inString = false
			}
			continue
		}
		if c == '"' {
			inString = true
			normalized = append(normalized, c)
			continue
		}
		if c == '-' || c >= '0' && c <= '9' {
			normalized = append(normalized, '0')
			for i+1 < len(raw) && strings.ContainsRune("0123456789.eE+-", rune(raw[i+1])) {
				i++
			}
			continue
		}
		normalized = append(normalized, c)
	}
	if _, e := canonicaljson.DecodeStrict(normalized); e != nil {
		return nil, invalid()
	}
	var result map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if e := decoder.Decode(&result); e != nil || result == nil {
		return nil, invalid()
	}
	return result, nil
}
