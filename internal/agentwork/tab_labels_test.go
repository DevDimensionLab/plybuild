package agentwork

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOverviewTabLabelRenameAndMissing(t *testing.T) {
	s, _, _ := fixture(t)
	i := s.Runtime.(*fakeRuntime).agents[0].Identity
	path := filepath.Join(t.TempDir(), "label")
	h := runtimeScript(t, `case "$1 $2" in
 "agent list") printf '%s' '{"result":{"agents":[{"agent":"codex","agent_status":"idle","workspace_id":"w-fixture","tab_id":"tab-fixture","pane_id":"pane-fixture","terminal_id":"terminal-fixture","agent_session":{"agent":"codex","kind":"id","value":"session-fixture"}}]}}';;
 "tab list") cat "`+path+`";;
 *) exit 9;; esac`)
	// The response file is replaced between requests: metadata is never cached.
	s.Runtime = h
	var first Overview
	for n, label := range []string{"Actual label", "Renamed 🌱", ""} {
		writeLabelResponse(t, path, label)
		got, err := s.Overview()
		if err != nil || len(got.Agents) != 1 {
			t.Fatal(got, err)
		}
		b, _ := json.Marshal(got.Agents[0])
		var wire map[string]any
		_ = json.Unmarshal(b, &wire)
		if label == "" {
			if _, ok := wire["tab_label"]; ok {
				t.Fatal("empty label emitted")
			}
		} else if wire["tab_label"] != label {
			t.Fatalf("label=%v want %q", wire["tab_label"], label)
		}
		if got.Agents[0].AgentRef != i.Ref() {
			t.Fatal("label changed identity")
		}
		if n == 0 {
			first = got
		} else {
			a, c := first.Agents[0], got.Agents[0]
			// Compare every existing projection field, including source/capabilities.
			x, _ := json.Marshal(a)
			y, _ := json.Marshal(c)
			var ax, ay map[string]any
			_ = json.Unmarshal(x, &ax)
			_ = json.Unmarshal(y, &ay)
			delete(ax, "tab_label")
			delete(ay, "tab_label")
			if !reflect.DeepEqual(ax, ay) || !reflect.DeepEqual(first.Sources, got.Sources) {
				t.Fatal("metadata changed base result")
			}
		}
	}
}

func writeLabelResponse(t *testing.T, path, label string) {
	t.Helper()
	b, _ := json.Marshal(map[string]any{"result": map[string]any{"tabs": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "tab-fixture", "label": label}}}})
	if e := os.WriteFile(path, b, 0600); e != nil {
		t.Fatal(e)
	}
}
