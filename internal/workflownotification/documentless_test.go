package workflownotification

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A request-v4 carries a public event without any source document.
func documentlessFixture(t *testing.T) (*fixture, map[string]any) {
	t.Helper()
	return compactFixture(t, "agent_finished", "done", "claude")
}
func unknownTransport(f *fixture) {
	f.deps.Transport = func(context.Context, string, []byte) Observation {
		f.posts++
		return Observation{Failed: true, Dispatch: "started"}
	}
}
func TestDocumentlessFinishedSendsOnceWithoutDocuments(t *testing.T) {
	f, _ := documentlessFixture(t)
	p := f.preview(t)
	if p.Payload.Text != "16:00:00 · claude\nDone: Varsler — Ferdig.\nNext: No action required." {
		t.Fatal(p.Payload.Text)
	}
	out, err := f.apply(t, p)
	if err != nil || out.(Result).State != "transport_acknowledged" || f.posts != 1 {
		t.Fatal(out, err)
	}
	b, _ := json.Marshal(out)
	for _, field := range []string{`"handoff"`, `"record"`, `"report"`, `"start_receipt"`} {
		if strings.Contains(string(b), field) {
			t.Fatal("fabricated document", field, string(b))
		}
	}
	if out, err = f.apply(t, p); err != nil || out.(Result).State != "transport_acknowledged" || f.posts != 1 {
		t.Fatal("replayed POST", out, err)
	}
	f.input.Operation, f.input.ID = "show", p.ID
	if out, err = Execute(f.deps, f.input); err != nil || out.(Result).Freshness != "current" {
		t.Fatal("readback", out, err)
	}
}
func TestDirectoryNameIsPrivate(t *testing.T) {
	for _, name := range []string{"a<b>", strings.Repeat("x", 97)} {
		f, m := documentlessFixture(t)
		origin := filepath.Join(f.root, name)
		if err := os.Mkdir(origin, 0700); err != nil {
			t.Fatal(err)
		}
		presentation := m["presentation"].(map[string]any)
		presentation["origin_cwd"], presentation["context"] = origin, name
		writeJSON(t, f.input.File, m)
		if p := f.preview(t); strings.Contains(p.Payload.Text, name) {
			t.Fatal("directory in payload", p.Payload.Text)
		}
		presentation["context"] = "other"
		writeJSON(t, f.input.File, m)
		if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
			t.Fatal("context unrelated to the physical origin", err)
		}
	}
}
func TestDocumentlessRejectsDocumentsAndOtherVersions(t *testing.T) {
	locator := map[string]any{"path": "/private/tmp/x", "sha256": "sha256:" + strings.Repeat("a", 64), "kind": "K@1"}
	for name, change := range map[string]func(map[string]any){
		"handoff":         func(m map[string]any) { m["source"].(map[string]any)["handoff"] = locator },
		"report":          func(m map[string]any) { m["source"].(map[string]any)["report"] = locator },
		"start":           func(m map[string]any) { m["source"].(map[string]any)["start_receipt"] = locator },
		"record":          func(m map[string]any) { m["event"].(map[string]any)["record"] = locator },
		"v5":              func(m map[string]any) { m["schema_version"] = 5 },
		"no-presentation": func(m map[string]any) { delete(m, "presentation") },
		"bad-status":      func(m map[string]any) { m["public"].(map[string]any)["status"] = "stopped" },
	} {
		t.Run(name, func(t *testing.T) {
			f, m := documentlessFixture(t)
			change(m)
			writeJSON(t, f.input.File, m)
			if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
				t.Fatal(err)
			}
		})
	}
}
func TestDocumentlessUsedIDConflictsAndUnknownOnlyReads(t *testing.T) {
	f, m := documentlessFixture(t)
	unknownTransport(f)
	p := f.preview(t)
	if out, _ := f.apply(t, p); out.(Result).State != "unknown" || f.posts != 1 {
		t.Fatal(out)
	}
	if out, err := f.apply(t, p); out.(Result).State != "unknown" || f.posts != 1 {
		t.Fatal("exact replay must only read", out, err)
	}
	retry := Input{Operation: "retry", ID: p.ID, Route: f.input.Route}
	if _, err := Execute(f.deps, retry); exitOf(err) != 4 {
		t.Fatal("unknown permitted retry", err)
	}
	m["public"].(map[string]any)["summary"] = "Endret tekst."
	writeJSON(t, f.input.File, m)
	if out, err := Execute(f.deps, f.input); exitOf(err) != 4 || out.(Result).State != "unknown" || f.posts != 1 {
		t.Fatal("changed text on a used ID", out, err)
	}
}
