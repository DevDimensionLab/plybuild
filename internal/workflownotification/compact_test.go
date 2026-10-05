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

// These acceptance inputs deliberately use public JSON, including before the
// implementation knew the v3 fields. They exercise validation and actual effects.
func compactFixture(t *testing.T, event, status, provider string) (*fixture, map[string]any) {
	t.Helper()
	f := agentFixture(t, event)
	b, _ := os.ReadFile(f.input.File)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	m["schema_version"] = 3
	pub := m["public"].(map[string]any)
	pub["status"] = status
	pub["task_title"] = "Varsler"
	pub["summary"] = "Ferdig."
	pub["next_action"] = "No action required."
	if event == "feedback_required" {
		pub["summary"] = "Hvilket mål?"
		pub["next_action"] = "Reply in the active agent conversation."
	}
	origin := filepath.Join(f.root, "ply", "planning")
	if err := os.MkdirAll(origin, 0700); err != nil {
		t.Fatal(err)
	}
	m["presentation"] = map[string]any{"provider": provider, "origin_cwd": origin, "context_root": f.root, "context": "ply/planning", "timezone": "Europe/Oslo", "local_occurred_at": "2026-10-04T16:00:00+02:00"}
	saveCompact(t, f, m)
	return f, m
}
func saveCompact(t *testing.T, f *fixture, m map[string]any) {
	t.Helper()
	ev := m["event"].(map[string]any)
	rec := ev["record"].(map[string]any)
	rec["kind"] = "PlyAgentNotificationEvent@2"
	event := map[string]any{"kind": rec["kind"], "activity": f.request.Source.Activity, "run": f.request.Source.Run, "event_id": ev["id"], "event_type": ev["type"], "phase": ev["phase"], "occurred_at": ev["occurred_at"], "public": m["public"]}
	writeJSON(t, rec["path"].(string), event)
	b, err := os.ReadFile(rec["path"].(string))
	if err != nil {
		t.Fatal(err)
	}
	rec["sha256"] = Digest(b)
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
				want := "16:00:00 · " + provider + " · ply/planning\n" + c.label + ": Varsler — " + pub["summary"].(string) + "\nNext: " + pub["next_action"].(string)
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
			want := c.clock + " · codex · ply/planning\nDone: Varsler\nNext: No action required."
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
func TestCompactVersionCollisionNeverDispatches(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		f, m := compactFixture(t, "agent_finished", "done", "codex")
		// First preserve v2 for exactly the same identity and source paths.
		if unknown {
			f.deps.Transport = func(context.Context, string, []byte) Observation {
				f.posts++
				return Observation{Failed: true, Dispatch: "started"}
			}
		}
		f.save(t)
		saveAgent(t, f)
		p := f.preview(t)
		f.apply(t, p)
		saveCompact(t, f, m)
		if _, err := Execute(f.deps, f.input); exitOf(err) != 4 || f.posts != 1 {
			t.Fatal("version collision was not blocked", err)
		}
	}
}

func TestCompactStructuralReadbackAndSourceDrift(t *testing.T) {
	f, m := compactFixture(t, "agent_finished", "done", "codex")
	p := f.preview(t)
	record := m["event"].(map[string]any)["record"].(map[string]any)["path"].(string)
	raw, _ := os.ReadFile(record)
	if err := os.WriteFile(record, append(raw, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := f.apply(t, p); exitOf(err) != 4 || f.posts != 0 {
		t.Fatal("source drift dispatched", err)
	}
	saveCompact(t, f, m)
	m["presentation"].(map[string]any)["timezone"] = "Future/Retired_Zone"
	saveCompact(t, f, m)
	raw, _ = os.ReadFile(f.input.File)
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
	if p := f.preview(t); !strings.HasPrefix(p.Payload.Text, "16:00:00 · claude · .\n") {
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
	context += "x"
	presentation["origin_cwd"] = filepath.Join(f.root, context)
	presentation["context"] = context
	if err := os.Mkdir(presentation["origin_cwd"].(string), 0700); err != nil {
		t.Fatal(err)
	}
	saveCompact(t, f, m)
	if _, err := Execute(f.deps, f.input); err == nil {
		t.Fatal("long display context accepted")
	}
}
func TestCompactRejectsMixedEventAndSymlinkOrigin(t *testing.T) {
	for _, variant := range []string{"event-kind", "legacy-with-presentation", "symlink", "event-status", "invalid-unicode"} {
		t.Run(variant, func(t *testing.T) {
			f, m := compactFixture(t, "agent_finished", "done", "codex")
			switch variant {
			case "event-kind":
				m["event"].(map[string]any)["record"].(map[string]any)["kind"] = agentEventKind
			case "legacy-with-presentation":
				m["schema_version"] = 2
			case "symlink":
				p := m["presentation"].(map[string]any)
				path := p["origin_cwd"].(string)
				if err := os.Rename(path, path+"-real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(path+"-real", path); err != nil {
					t.Fatal(err)
				}
			case "event-status":
				path := m["event"].(map[string]any)["record"].(map[string]any)["path"].(string)
				b, _ := os.ReadFile(path)
				b = bytes.Replace(b, []byte(`"status":"done"`), []byte(`"status":"ready_for_review"`), 1)
				os.WriteFile(path, b, 0600)
				m["event"].(map[string]any)["record"].(map[string]any)["sha256"] = Digest(b)
			}
			writeJSON(t, f.input.File, m)
			if variant == "invalid-unicode" {
				b, _ := os.ReadFile(f.input.File)
				b = bytes.Replace(b, []byte("Ferdig."), []byte(`\ud800`), 1)
				os.WriteFile(f.input.File, b, 0600)
			}
			if _, err := Execute(f.deps, f.input); err == nil || f.posts != 0 {
				t.Fatal("bad v3 accepted", err)
			}
			if _, err := os.Stat(f.route.StateRoot); !os.IsNotExist(err) {
				t.Fatal("bad v3 reserved state")
			}
		})
	}
}
