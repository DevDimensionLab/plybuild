package workflownotification

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// legacyRequest writes the request and source documents a version before v4
// bound, as the removed delivery process produced them.
func legacyRequest(t *testing.T, f *fixture, m map[string]any, version int) []byte {
	t.Helper()
	source := m["source"].(map[string]any)
	event := m["event"].(map[string]any)
	public := m["public"].(map[string]any)
	document := func(name string, v any) map[string]any {
		path := filepath.Join(f.root, name)
		var b []byte
		if s, ok := v.(string); ok {
			b = []byte(s)
		} else {
			b, _ = json.Marshal(v)
		}
		if err := os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
		return map[string]any{"path": path, "sha256": Digest(b)}
	}
	identity := func(kind string) map[string]any {
		return map[string]any{"kind": kind, "activity": source["activity"], "run": source["run"]}
	}
	withKind := func(l map[string]any, kind string) map[string]any { l["kind"] = kind; return l }
	legacy := map[string]any{"kind": m["kind"], "schema_version": version, "route": m["route"], "sender": m["sender"]}
	src := map[string]any{"kind": "external", "activity": source["activity"], "run": source["run"], "worktree": source["worktree"], "handoff": document("handoff.md", "Fixture handoff\n"),
		"report": withKind(document("report.json", identity("FixtureReport@1")), "FixtureReport@1")}
	legacy["source"] = src
	if version == 1 {
		src["start_receipt"] = withKind(document("start.json", identity("FixtureStart@1")), "FixtureStart@1")
		legacy["gate"] = map[string]any{"id": "result-control", "revision": 1, "opened_at": "2026-10-03T19:00:00Z", "reason": "result_control", "state": "waiting_for_human"}
		legacy["public"] = map[string]any{"task_title": public["task_title"], "next_action": public["next_action"]}
		b, _ := json.Marshal(legacy)
		return b
	}
	recordKind, pub := "PlyAgentNotificationEvent@2", public
	if version == 2 {
		recordKind, pub = "PlyAgentNotificationEvent@1", map[string]any{}
		for k, v := range public {
			if k != "status" {
				pub[k] = v
			}
		}
	} else {
		legacy["presentation"] = m["presentation"]
	}
	legacy["public"] = pub
	record := document("event.json", map[string]any{"kind": recordKind, "activity": source["activity"], "run": source["run"], "event_id": event["id"], "event_type": event["type"], "phase": event["phase"], "occurred_at": event["occurred_at"], "public": pub})
	legacy["event"] = map[string]any{"id": event["id"], "type": event["type"], "occurred_at": event["occurred_at"], "phase": event["phase"], "record": withKind(record, recordKind)}
	b, _ := json.Marshal(legacy)
	return b
}

// Versions before v4 bound documents of the removed delivery process. A new send
// of one is refused before any state or post, even when its documents verify.
func TestRemovedRequestVersionsAreRefusedWithoutEffect(t *testing.T) {
	for version := 1; version <= 3; version++ {
		f, m := documentlessFixture(t)
		if err := os.WriteFile(f.input.File, legacyRequest(t, f, m, version), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := Execute(f.deps, f.input); exitOf(err) != 2 || f.posts != 0 {
			t.Errorf("v%d request accepted: %v", version, err)
		}
		if _, err := os.Stat(f.route.StateRoot); !os.IsNotExist(err) {
			t.Errorf("v%d request reserved state", version)
		}
	}
}

// A stored record of a removed version, including a v1 gate without an event, is history: it reads back from its raw
// bytes and stored outcome, and neither retry nor a new format of the same
// event ever posts again.
func TestStoredRecordOfRemovedVersionIsReadOnlyHistory(t *testing.T) {
	for version := 1; version <= 3; version++ {
		t.Run(fmt.Sprintf("v%d", version), func(t *testing.T) { readOnlyHistory(t, version) })
	}
}
func readOnlyHistory(t *testing.T, version int) {
	f, m := documentlessFixture(t)
	v4, err := os.ReadFile(f.input.File)
	if err != nil {
		t.Fatal(err)
	}
	f.deps.Transport = func(context.Context, string, []byte) Observation {
		f.posts++
		return Observation{Status: 404, Body: []byte("channel_not_found"), Dispatch: "started"}
	}
	p := f.preview(t)
	if out, _ := f.apply(t, p); out.(Result).State != "rejected" || f.posts != 1 {
		t.Fatal(out)
	}
	legacy := legacyRequest(t, f, m, version)
	if err = os.WriteFile(f.input.File, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	dir, err := stateDirectory(f.route.StateRoot, false)
	if err != nil {
		t.Fatal(err)
	}
	db, err := readDatabase(dir)
	if err != nil {
		t.Fatal(err)
	}
	r := db.Records[p.ID]
	r.RequestBytes = legacy
	r.Basis = digestValue([]any{Digest(legacy), r.RequestPath, r.RouteSHA})
	r.Payload = Payload{Text: "16:00:00 · claude · ply/planning\nDone: Varsler — Ferdig.\nNext: No action required.", Parse: "none"}
	if err = writeDatabase(dir, db); err != nil {
		t.Fatal(err)
	}
	dir.Close()
	stored := stateBytes(t, f)

	show := Input{Operation: "show", Route: f.input.Route, ID: p.ID}
	out, err := Execute(f.deps, show)
	if res, ok := out.(Result); err != nil || !ok || res.ID != p.ID || res.State != "rejected" || res.RequestSHA != Digest(legacy) || (version > 1) != (res.Event != nil) {
		t.Fatal("historical record does not read back", out, err)
	}
	if res := out.(Result); len(res.Reasons) != 1 || res.Reasons[0].Code != "historical_request" {
		t.Fatal("historical readback does not say it is never sent again", res.Reasons)
	}
	if _, err = Execute(f.deps, Input{Operation: "retry", Route: f.input.Route, ID: p.ID}); exitOf(err) != 4 {
		t.Fatal("historical record retried", err)
	}
	if err = os.WriteFile(f.input.File, v4, 0600); err != nil {
		t.Fatal(err)
	}
	if out, err = Execute(f.deps, f.input); exitOf(err) != 4 || out.(Result).ID != p.ID {
		t.Fatal("new format of a historical event opened a send", out, err)
	}
	if f.posts != 1 || string(stored) != string(stateBytes(t, f)) {
		t.Fatal("history changed or posted again", f.posts)
	}
}
