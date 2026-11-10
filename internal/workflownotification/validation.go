package workflownotification

import (
	"bytes"
	"encoding/json"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
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
	if !agentSchema(r.SchemaVersion) {
		return r, invalid()
	}
	return parseAgentRequest(r, m, secret, fresh)
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
	if r.Event == nil {
		return ""
	}
	// This historical namespace is intentionally independent of request versions.
	return "ntf_" + strings.TrimPrefix(digestValue([]any{"agent-event-v2", r.Source.Activity, r.Source.Run, r.Event.ID, r.Route}), "sha256:")
}
func message(r Request) Payload { return compactMessage(r) }
func checkSources(r Request) error {
	d, err := openDirectory(r.Source.Worktree)
	if err != nil {
		return invalid()
	}
	d.Close()
	return nil
}
