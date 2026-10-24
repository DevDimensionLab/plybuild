package agentwork

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func metadataRow(pane string) map[string]any {
	return map[string]any{"agent": "codex", "agent_status": "future", "pane_id": pane, "tab_id": "t", "workspace_id": "w", "terminal_id": false, "agent_session": map[string]any{"invalid": true}}
}
func metadataRaw(v any) string { b, _ := json.Marshal(v); return string(b) }
func metadataFixture(t *testing.T, raw string) MetadataReadback {
	t.Helper()
	p := filepath.Join(t.TempDir(), "response")
	if e := os.WriteFile(p, []byte(raw), 0600); e != nil {
		t.Fatal(e)
	}
	h := runtimeScript(t, "cat '"+p+"'")
	return StatusMetadata(h)
}
func TestStatusMetadataQualification(t *testing.T) {
	row := metadataRow("p")
	envelope := func(rows any) string {
		return metadataRaw(map[string]any{"result": map[string]any{"type": "agent_list", "agents": rows}})
	}
	tests := []struct {
		name, raw, status, coverage string
		count                       int
	}{
		{"sessionless", envelope([]any{row}), "available", "complete", 1},
		{"empty", envelope([]any{}), "empty", "complete", 0},
		{"null", envelope(nil), "unavailable", "unknown", 0},
		{"malformed-envelope", "{", "unavailable", "unknown", 0},
		{"array-envelope", "[]", "unavailable", "unknown", 0},
		{"unknown-type", `{"result":{"type":"future","agents":[]}}`, "unavailable", "unknown", 0},
		{"missing-type", `{"result":{"agents":[]}}`, "unavailable", "unknown", 0},
		{"source-error", `{"result":{"type":"agent_list","agents":[]},"error":{}}`, "unavailable", "unknown", 0},
		{"duplicate-envelope", `{"result":{},"result":{"type":"agent_list","agents":[]}}`, "unavailable", "unknown", 0},
		{"truncated", `{"result":{"type":"agent_list","agents":[]`, "unavailable", "unknown", 0},
		{"envelope-truncation-signal", `{"id":"read","result":{"type":"agent_list","agents":[]},"truncated":true}`, "partial", "partial", 0},
		{"truncation-signal", `{"result":{"type":"agent_list","agents":[],"truncated":true}}`, "partial", "partial", 0},
		{"unknown-scope", `{"result":{"type":"agent_list","agents":[],"scope":"subset"}}`, "partial", "partial", 0},
		{"bad-row", envelope([]any{row, 42}), "partial", "partial", 1},
		{"duplicate-pane", envelope([]any{row, row}), "partial", "partial", 0},
		{"duplicate-metadata", `{"result":{"type":"agent_list","agents":[{"agent":"codex","pane_id":"a","pane_id":"b","tab_id":"t","workspace_id":"w","agent_status":"idle"}]}}`, "partial", "partial", 0},
		{"execution-opaque", `{"result":{"type":"agent_list","agents":[{"agent":"codex","pane_id":"a","tab_id":"t","workspace_id":"w","agent_status":"idle","agent_session":false,"terminal_id":{},"agent_session":42}]}}`, "available", "complete", 1},
	}
	for _, field := range []string{"pane_id", "tab_id", "workspace_id", "agent_status"} {
		for _, value := range []any{nil, 42, "bad\ntext", "/private/path", strings.Repeat("😀", 128) + "x", "\u007f"} {
			r := metadataRow("p")
			r[field] = value
			tests = append(tests, struct {
				name, raw, status, coverage string
				count                       int
			}{fmt.Sprintf("invalid-%s-%v", field, value), envelope([]any{r}), "partial", "partial", 0})
		}
	}
	for _, text := range []string{strings.Repeat("a", 256), strings.Repeat("😀", 128)} {
		r := metadataRow(text)
		tests = append(tests, struct {
			name, raw, status, coverage string
			count                       int
		}{"utf16-256", envelope([]any{r}), "available", "complete", 1})
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := metadataFixture(t, tc.raw)
			s := v.Sources[0]
			if s.Status != tc.status || s.Coverage != tc.coverage || len(v.Agents) != tc.count {
				t.Fatalf("%+v", v)
			}
			if s.Status != "unavailable" && (s.SourceObservedAtUTC == nil || *s.SourceObservedAtUTC != v.ObservedAtUTC) {
				t.Fatal("not response-observed", v)
			}
			for _, a := range v.Agents {
				if a.Status == "future" {
					t.Fatal(a)
				}
			}
		})
	}
}
func TestStatusMetadataRowAndByteBounds(t *testing.T) {
	for _, count := range []int{4096, 4097} {
		rows := []any{}
		for n := 0; n < count; n++ {
			rows = append(rows, metadataRow(fmt.Sprintf("p%04d", n)))
		}
		v := metadataFixture(t, metadataRaw(map[string]any{"result": map[string]any{"type": "agent_list", "agents": rows}}))
		want := "complete"
		if count > 4096 {
			want = "partial"
		}
		if len(v.Agents) != 4096 || v.Sources[0].Coverage != want {
			t.Fatal(len(v.Agents), v.Sources)
		}
	}
	v := metadataFixture(t, strings.Repeat(" ", 8<<20)+`{"result":{"type":"agent_list","agents":[]}}`)
	if v.Sources[0].Status != "unavailable" {
		t.Fatal(v)
	}
}
func TestStatusMetadataDeadlineAndResponseClock(t *testing.T) {
	for _, body := range []string{"exit 7", "sleep 20", "sleep 20 & wait"} {
		h := runtimeScript(t, body)
		h.Timeout = 50 * time.Millisecond
		start := time.Now()
		v := StatusMetadata(h)
		elapsed := time.Since(start)
		observed, e := time.Parse(time.RFC3339Nano, v.ObservedAtUTC)
		if e != nil || observed.Before(start) || v.Sources[0].Status != "unavailable" || elapsed > 1500*time.Millisecond {
			t.Fatal(v, elapsed, e)
		}
	}
}
func TestStatusMetadataDoesNotUpgradeFullExecutionPartial(t *testing.T) {
	s, _, _ := fixture(t)
	s.Runtime.(*fakeRuntime).agents[0].Session = ""
	v, e := s.Overview()
	if e != nil || v.Sources[0].Status != "partial" || v.Agents[0].BindingStatus != "unknown" {
		t.Fatal(v, e)
	}
}

func TestStatusMetadataClockIsResponseObservation(t *testing.T) {
	h := runtimeScript(t, `sleep 0.05
printf '%s' '{"result":{"type":"agent_list","agents":[]}}'`)
	start := time.Now()
	v := StatusMetadata(h)
	observed, e := time.Parse(time.RFC3339Nano, v.ObservedAtUTC)
	if e != nil || observed.Before(start.Add(40*time.Millisecond)) || v.Sources[0].SourceObservedAtUTC == nil || *v.Sources[0].SourceObservedAtUTC != v.ObservedAtUTC {
		t.Fatal(v, e)
	}
}

// SM01 r001-f01: malformed rows still claim panes before other fields qualify.
func TestStatusMetadataOmitsEveryMalformedDuplicatePane(t *testing.T) {
	valid := metadataRaw(metadataRow("w1:p1"))
	unaffected := metadataRaw(metadataRow("w1:p-safe"))
	cases := map[string]string{}
	for name, field := range map[string]string{"invalid-status": "agent_status", "null-tab": "tab_id", "invalid-provider": "agent"} {
		row := metadataRow("w1:p1")
		row[field] = 42
		if name == "null-tab" {
			row[field] = nil
		}
		cases[name] = metadataRaw(row)
	}
	cases["duplicate-status-before-pane"] = `{"agent":"codex","agent_status":"idle","agent_status":42,"pane_id":"w1:p1","tab_id":"t","workspace_id":"w"}`
	cases["duplicate-status-after-pane"] = `{"pane_id":"w1:p1","agent":"codex","agent_status":"idle","agent_status":42,"tab_id":"t","workspace_id":"w"}`
	cases["duplicate-provider-before-pane"] = `{"agent":"codex","agent":42,"pane_id":"w1:p1","tab_id":"t","workspace_id":"w","agent_status":"idle"}`
	cases["duplicate-pane-same-claim"] = `{"agent":"codex","pane_id":"w1:p1","pane_id":"w1:p1","tab_id":"t","workspace_id":"w","agent_status":"idle"}`
	cases["duplicate-pane-different-claims"] = `{"agent":"codex","pane_id":"w1:p1","pane_id":"w1:p2","tab_id":"t","workspace_id":"w","agent_status":"idle"}`
	for name, malformed := range cases {
		for _, reverse := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/reverse=%t", name, reverse), func(t *testing.T) {
				rows := []string{valid, malformed, unaffected}
				if reverse {
					rows[0], rows[1] = rows[1], rows[0]
				}
				raw := `{"result":{"type":"agent_list","agents":[` + strings.Join(rows, ",") + `]}}`
				v := metadataFixture(t, raw)
				if v.Sources[0].Status != "partial" || v.Sources[0].Coverage != "partial" || len(v.Agents) != 1 || v.Agents[0].Pane != "w1:p-safe" {
					t.Fatalf("ambiguous pane survives or unrelated pane lost: %+v", v)
				}
			})
		}
	}
}
