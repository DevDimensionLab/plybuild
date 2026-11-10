package workflownotification

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These acceptance inputs deliberately use public JSON. They exercise
// validation and actual effects.
func compactFixture(t *testing.T, event, status, provider string) (*fixture, map[string]any) {
	t.Helper()
	f := newFixture(t)
	f.request.Event.Type = event
	f.request.Public.Status = status
	f.request.Presentation.Provider = provider
	if event == "feedback_required" {
		f.request.Public.Summary = "Hvilket mål?"
		f.request.Public.NextAction = "Reply in the active agent conversation."
		f.request.Public.NextActor = "user"
	}
	f.save(t)
	b, _ := os.ReadFile(f.input.File)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return f, m
}
func saveCompact(t *testing.T, f *fixture, m map[string]any) {
	t.Helper()
	writeJSON(t, f.input.File, m)
}
func TestCompactGoldenStatusesAndProviders(t *testing.T) {
	for _, c := range []struct{ event, status, label string }{
		{"agent_finished", "ready_for_review", "Ready for review"},
		{"agent_finished", "ready_for_your_check", "Ready for your check"},
		{"agent_finished", "done", "Done"},
		{"agent_stopped", "stopped", "Stopped"},
		{"feedback_required", "needs_answer", "Needs answer"},
	} {
		for _, provider := range []string{"codex", "claude"} {
			t.Run(c.status+"/"+provider, func(t *testing.T) {
				f, m := compactFixture(t, c.event, c.status, provider)
				pub := m["public"].(map[string]any)
				want := "16:00:00 · " + provider + "\n" + c.label + ": Varsler — " + pub["summary"].(string) + "\nNext: " + pub["next_action"].(string)
				p := f.preview(t)
				if p.Payload.Text != want || p.Payload.Mrkdwn || p.Payload.Parse != "none" || p.Payload.UnfurlLinks || p.Payload.UnfurlMedia {
					t.Fatalf("payload: %+v", p.Payload)
				}
				var sent Payload
				f.deps.Transport = func(_ context.Context, _ string, b []byte) Observation {
					f.posts++
					if err := json.Unmarshal(b, &sent); err != nil {
						t.Fatal(err)
					}
					return Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
				}
				out, err := f.apply(t, p)
				if err != nil || sent.Text != want || out.(Result).Knowledge != "reported" {
					t.Fatal(out, err, sent)
				}
				state := stateBytes(t, f)
				t.Setenv("TZ", "Pacific/Honolulu")
				out, err = f.apply(t, p)
				if err != nil || f.posts != 1 || !bytes.Equal(state, stateBytes(t, f)) {
					t.Fatal("replay changed history", out, err)
				}
			})
		}
	}
}
func TestCompactLocalTimeAndRepeatedSummary(t *testing.T) {
	for _, c := range []struct{ utc, local, clock string }{
		{"2026-01-04T14:00:00Z", "2026-01-04T15:00:00+01:00", "15:00:00"},
		{"2026-07-04T14:00:00Z", "2026-07-04T16:00:00+02:00", "16:00:00"},
		{"2026-07-04T23:30:00Z", "2026-07-05T01:30:00+02:00", "01:30:00"},
	} {
		t.Run(c.utc, func(t *testing.T) {
			f, m := compactFixture(t, "agent_finished", "done", "codex")
			m["event"].(map[string]any)["occurred_at"] = c.utc
			m["presentation"].(map[string]any)["local_occurred_at"] = c.local
			m["public"].(map[string]any)["summary"] = "Varsler"
			saveCompact(t, f, m)
			want := c.clock + " · codex\nDone: Varsler\nNext: No action required."
			if p := f.preview(t); p.Payload.Text != want {
				t.Fatal(p.Payload)
			}
		})
	}
}
func TestCompactInvalidInputHasNoReservation(t *testing.T) {
	cases := map[string]func(map[string]any){
		"provider":     func(m map[string]any) { m["presentation"].(map[string]any)["provider"] = "model-name" },
		"zone":         func(m map[string]any) { m["presentation"].(map[string]any)["timezone"] = "No/Such_Zone" },
		"process-zone": func(m map[string]any) { m["presentation"].(map[string]any)["timezone"] = "Local" },
		"offset": func(m map[string]any) {
			m["presentation"].(map[string]any)["local_occurred_at"] = "2026-10-04T15:00:00+01:00"
		},
		"instant": func(m map[string]any) {
			m["presentation"].(map[string]any)["local_occurred_at"] = "2026-10-04T16:00:01+02:00"
		},
		"context":        func(m map[string]any) { m["presentation"].(map[string]any)["context"] = "wrong/path" },
		"traversal":      func(m map[string]any) { m["presentation"].(map[string]any)["context"] = "../planning" },
		"origin-outside": func(m map[string]any) { m["presentation"].(map[string]any)["origin_cwd"] = "/" },
		"origin-missing": func(m map[string]any) {
			p := m["presentation"].(map[string]any)
			p["origin_cwd"] = p["context_root"].(string) + "/missing"
			p["context"] = "missing"
		},
		"status":         func(m map[string]any) { m["public"].(map[string]any)["status"] = "approved" },
		"status-event":   func(m map[string]any) { m["public"].(map[string]any)["status"] = "stopped" },
		"title-limit":    func(m map[string]any) { m["public"].(map[string]any)["task_title"] = strings.Repeat("ø", 61) },
		"summary-limit":  func(m map[string]any) { m["public"].(map[string]any)["summary"] = strings.Repeat("ø", 101) },
		"next-limit":     func(m map[string]any) { m["public"].(map[string]any)["next_action"] = strings.Repeat("ø", 141) },
		"newline":        func(m map[string]any) { m["public"].(map[string]any)["summary"] = "first\nsecond" },
		"unicode-line":   func(m map[string]any) { m["public"].(map[string]any)["summary"] = "first\u2028second" },
		"bidi":           func(m map[string]any) { m["public"].(map[string]any)["summary"] = "abc\u202edef" },
		"mention":        func(m map[string]any) { m["public"].(map[string]any)["summary"] = "Hi @channel" },
		"injection":      func(m map[string]any) { m["public"].(map[string]any)["summary"] = "<@someone>" },
		"secret":         func(m map[string]any) { m["public"].(map[string]any)["summary"] = syntheticSecret },
		"unknown-field":  func(m map[string]any) { m["presentation"].(map[string]any)["extra"] = true },
		"missing-status": func(m map[string]any) { delete(m["public"].(map[string]any), "status") },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			f, m := compactFixture(t, "agent_finished", "done", "codex")
			change(m)
			saveCompact(t, f, m)
			if _, err := Execute(f.deps, f.input); err == nil || f.posts != 0 {
				t.Fatal("invalid input accepted", err)
			}
			if _, err := os.Stat(f.route.StateRoot); !os.IsNotExist(err) {
				t.Fatal("invalid input reserved state", err)
			}
		})
	}
}
func TestCompactStoredRequestIgnoresTodaysZoneRules(t *testing.T) {
	f, m := compactFixture(t, "agent_finished", "done", "codex")
	m["presentation"].(map[string]any)["timezone"] = "Future/Retired_Zone"
	saveCompact(t, f, m)
	raw, _ := os.ReadFile(f.input.File)
	if _, err := parseRequest(raw, ""); err == nil {
		t.Fatal("fresh unknown zone accepted")
	}
	if _, err := parseStoredRequest(raw, ""); err != nil {
		t.Fatal("history consulted tzdata", err)
	}
	m["presentation"].(map[string]any)["local_occurred_at"] = "2026-10-04T16:00:01+02:00"
	saveCompact(t, f, m)
	raw, _ = os.ReadFile(f.input.File)
	if _, err := parseStoredRequest(raw, ""); err == nil {
		t.Fatal("historical instant mismatch accepted")
	}
}
func TestCompactContextAndTextBoundaries(t *testing.T) {
	f, m := compactFixture(t, "agent_finished", "ready_for_your_check", "claude")
	presentation := m["presentation"].(map[string]any)
	presentation["origin_cwd"] = f.root
	presentation["context"] = "."
	saveCompact(t, f, m)
	if p := f.preview(t); !strings.HasPrefix(p.Payload.Text, "16:00:00 · claude\n") {
		t.Fatal(p.Payload)
	}
	pub := m["public"].(map[string]any)
	pub["task_title"] = strings.Repeat("ø", 60)
	pub["summary"] = strings.Repeat("ø", 100)
	pub["next_action"] = strings.Repeat("ø", 140)
	context := strings.Repeat("x", 96)
	presentation["origin_cwd"] = filepath.Join(f.root, context)
	presentation["context"] = context
	if err := os.Mkdir(presentation["origin_cwd"].(string), 0700); err != nil {
		t.Fatal(err)
	}
	saveCompact(t, f, m)
	f.preview(t)
}
func TestCompactRejectsRemovedVersionAndSymlinkOrigin(t *testing.T) {
	for _, variant := range []string{"removed-version", "symlink", "invalid-unicode"} {
		t.Run(variant, func(t *testing.T) {
			f, m := compactFixture(t, "agent_finished", "done", "codex")
			switch variant {
			case "removed-version":
				m["schema_version"] = 3
			case "symlink":
				p := m["presentation"].(map[string]any)
				path := p["origin_cwd"].(string)
				if err := os.Rename(path, path+"-real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(path+"-real", path); err != nil {
					t.Fatal(err)
				}
			}
			writeJSON(t, f.input.File, m)
			if variant == "invalid-unicode" {
				b, _ := os.ReadFile(f.input.File)
				b = bytes.Replace(b, []byte("Ferdig."), []byte(`\ud800`), 1)
				os.WriteFile(f.input.File, b, 0600)
			}
			if _, err := Execute(f.deps, f.input); err == nil || f.posts != 0 {
				t.Fatal("bad request accepted", err)
			}
			if _, err := os.Stat(f.route.StateRoot); !os.IsNotExist(err) {
				t.Fatal("bad request reserved state")
			}
		})
	}
}
