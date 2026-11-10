package workflownotification

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

// A request-v5 is v4 with the proven Herdr workspace and tab display labels.
func herdrFixture(t *testing.T, workspace, tab string) (*fixture, map[string]any) {
	t.Helper()
	f, m := documentlessFixture(t)
	m["schema_version"] = 5
	m["presentation"].(map[string]any)["herdr"] = map[string]any{"workspace": workspace, "tab": tab}
	writeJSON(t, f.input.File, m)
	return f, m
}
func TestHerdrOriginExtendsOnlyTheFirstLine(t *testing.T) {
	f, _ := herdrFixture(t, "ply-master", "pn01-herdr-notification-spec")
	p := f.preview(t)
	if p.Payload.Text != "16:00:00 · claude · ply-master › pn01-herdr-notification-spec\nDone: Varsler — Ferdig.\nNext: No action required." {
		t.Fatal(p.Payload.Text)
	}
	if out, err := f.apply(t, p); err != nil || out.(Result).State != "transport_acknowledged" || f.posts != 1 {
		t.Fatal(out, err)
	}
	f.input.Operation, f.input.ID = "show", p.ID
	if out, err := Execute(f.deps, f.input); err != nil || out.(Result).Freshness != "current" {
		t.Fatal("stored v5 record does not read back", out, err)
	}
}
func TestHerdrLabelsFitTheMessageLimit(t *testing.T) {
	label := strings.Repeat("w", 31) + "…"
	f, m := herdrFixture(t, label, label)
	public := m["public"].(map[string]any)
	public["task_title"], public["summary"], public["next_action"] = strings.Repeat("t", 60), strings.Repeat("s", 100), strings.Repeat("n", 140)
	writeJSON(t, f.input.File, m)
	if p := f.preview(t); utf8.RuneCountInString(p.Payload.Text) > 480 || !strings.Contains(p.Payload.Text, label+" › "+label+"\n") {
		t.Fatal(p.Payload.Text)
	}
}
func TestHerdrOriginIsRejectedWhenIncompleteOrUnsafe(t *testing.T) {
	for name, change := range map[string]func(map[string]any){
		"missing":        func(h map[string]any) { delete(h, "tab") },
		"extra":          func(h map[string]any) { h["pane_id"] = "w9:pE" },
		"empty":          func(h map[string]any) { h["workspace"] = "" },
		"33 characters":  func(h map[string]any) { h["tab"] = strings.Repeat("x", 33) },
		"mention":        func(h map[string]any) { h["tab"] = "<!here>" },
		"control":        func(h map[string]any) { h["tab"] = "a\u0007b" },
		"format":         func(h map[string]any) { h["workspace"] = "a\u202eb" },
		"untrimmed":      func(h map[string]any) { h["workspace"] = " ply" },
		"not an object":  nil,
		"in v4 requests": nil,
	} {
		t.Run(name, func(t *testing.T) {
			f, m := herdrFixture(t, "ply-master", "spec")
			presentation := m["presentation"].(map[string]any)
			switch name {
			case "not an object":
				presentation["herdr"] = "ply-master › spec"
			case "in v4 requests":
				m["schema_version"] = 4
			default:
				change(presentation["herdr"].(map[string]any))
			}
			writeJSON(t, f.input.File, m)
			if _, err := Execute(f.deps, f.input); exitOf(err) != 2 {
				t.Fatal(err)
			}
		})
	}
}

// v4 and v5 for one event share the identity, so adding an origin never sends twice.
func TestHerdrOriginSharesTheV4EventIdentity(t *testing.T) {
	for _, firstV5 := range []bool{false, true} {
		f, v5 := herdrFixture(t, "ply-master", "spec")
		raw, _ := os.ReadFile(f.input.File)
		var v4 map[string]any
		if err := json.Unmarshal(raw, &v4); err != nil {
			t.Fatal(err)
		}
		v4["schema_version"] = 4
		delete(v4["presentation"].(map[string]any), "herdr")
		first, second := v4, v5
		if firstV5 {
			first, second = v5, v4
		}
		writeJSON(t, f.input.File, first)
		p := f.preview(t)
		if out, _ := f.apply(t, p); out.(Result).State != "transport_acknowledged" {
			t.Fatal(out)
		}
		writeJSON(t, f.input.File, second)
		out, err := Execute(f.deps, f.input)
		if r, ok := out.(Result); exitOf(err) != 4 || !ok || r.ID != p.ID || f.posts != 1 {
			t.Fatal(out, err, f.posts)
		}
		in := f.input
		in.Apply, in.Confirm = true, *p.Confirmation
		if _, err = Execute(f.deps, in); exitOf(err) != 4 || f.posts != 1 {
			t.Fatal("second apply sent", err, f.posts)
		}
	}
}
