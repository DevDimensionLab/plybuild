package agentwork

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
	"unicode/utf16"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const MetadataKind = "AgentStatusMetadataReadback@1"
const MetadataTimeout = 2 * time.Second
const metadataScope = "Recognized codex/claude agents in the current Herdr user agent list; maximum 4096 source rows."
const metadataMaxAgents = 4096

type StatusAgent struct {
	Pane      string `json:"pane"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
	Tab       string `json:"tab"`
	Workspace string `json:"workspace"`
}
type MetadataSource struct {
	Source              string  `json:"source"`
	Status              string  `json:"status"`
	Freshness           string  `json:"freshness"`
	SourceObservedAtUTC *string `json:"source_observed_at_utc"`
	Coverage            string  `json:"coverage"`
	Scope               string  `json:"scope"`
}
type MetadataReadback struct {
	Kind          string           `json:"kind"`
	SchemaVersion int              `json:"schema_version"`
	ObservedAtUTC string           `json:"observed_at_utc"`
	Sources       []MetadataSource `json:"sources"`
	Agents        []StatusAgent    `json:"agents"`
}

// StatusMetadata performs exactly one public List call, independently of
// execution identity, the full Service, workspace discovery and native stores.
// Transport takes at most 2s plus the existing 1s pipe cleanup allowance.
func StatusMetadata(h HerdrRuntime) MetadataReadback {
	if h.Timeout <= 0 || h.Timeout > MetadataTimeout {
		h.Timeout = MetadataTimeout
	}
	raw, err := h.callBytes("agent", "list")
	observed := time.Now().UTC().Format(time.RFC3339Nano)
	out := MetadataReadback{Kind: MetadataKind, SchemaVersion: 1, ObservedAtUTC: observed, Agents: []StatusAgent{}, Sources: []MetadataSource{{Source: "herdr.agent.list", Status: "unavailable", Freshness: "unknown", Coverage: "unknown", Scope: metadataScope}}}
	if err != nil {
		return out
	}
	return decodeMetadata(raw, out)
}

// Decode only selected object members. Execution fields are opaque raw JSON;
// their types and duplicate members never qualify metadata completeness.
func metadataObject(raw []byte, selected map[string]bool) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	tok, e := d.Token()
	if e != nil || tok != json.Delim('{') {
		return nil, fmt.Errorf("metadata_object_invalid")
	}
	out := map[string]json.RawMessage{}
	for d.More() {
		tok, e = d.Token()
		if e != nil {
			return nil, e
		}
		key, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("metadata_key_invalid")
		}
		var v json.RawMessage
		if e = d.Decode(&v); e != nil {
			return nil, e
		}
		if selected == nil || selected[key] {
			if _, exists := out[key]; exists {
				return nil, fmt.Errorf("metadata_duplicate")
			}
			out[key] = v
		}
	}
	if _, e = d.Token(); e != nil {
		return nil, e
	}
	if _, e = d.Token(); e != io.EOF {
		return nil, fmt.Errorf("metadata_trailing")
	}
	return out, nil
}
func metadataText(raw []byte, nonempty bool) (string, bool) {
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok || len(utf16.Encode([]rune(s))) > 256 || nonempty && s == "" {
		return "", false
	}
	// Reuse the public text safety policy, allowing empty tab/workspace.
	return s, s == "" || plain(s)
}
func decodeMetadata(raw []byte, out MetadataReadback) MetadataReadback {
	env, e := metadataObject(raw, map[string]bool{"result": true, "error": true})
	if e != nil {
		return out
	}
	if v, exists := env["error"]; exists && string(bytes.TrimSpace(v)) != "null" {
		return out
	}
	result, e := metadataObject(env["result"], nil)
	if e != nil {
		return out
	}
	typ, ok := metadataText(result["type"], true)
	if !ok || typ != "agent_list" {
		return out
	}
	var rows []json.RawMessage
	if len(result["agents"]) == 0 || bytes.Equal(bytes.TrimSpace(result["agents"]), []byte("null")) || json.Unmarshal(result["agents"], &rows) != nil {
		return out
	}
	src := &out.Sources[0]
	src.Status = "available"
	src.Freshness = "fresh"
	src.Coverage = "complete"
	src.SourceObservedAtUTC = &out.ObservedAtUTC
	partial := func() { src.Status = "partial"; src.Coverage = "partial" }
	// agent_list is the supported nonpaginated public enumeration. Any additional
	// result member could narrow that enumeration: preserve uncertainty.
	if len(result) != 2 {
		partial()
	}
	if len(rows) > metadataMaxAgents {
		partial()
		rows = rows[:metadataMaxAgents]
	}
	selected := map[string]bool{"agent": true, "agent_status": true, "pane_id": true, "tab_id": true, "workspace_id": true}
	seen := map[string]int{}
	for _, rawRow := range rows {
		row, e := metadataObject(rawRow, selected)
		if e != nil {
			partial()
			continue
		}
		provider, ok := metadataText(row["agent"], true)
		if !ok {
			partial()
			continue
		}
		if provider != "codex" && provider != "claude" {
			continue
		}
		pane, pok := metadataText(row["pane_id"], true)
		tab, tok := metadataText(row["tab_id"], false)
		workspace, wok := metadataText(row["workspace_id"], false)
		status, sok := metadataText(row["agent_status"], false)
		if !pok || !tok || !wok || !sok {
			partial()
			continue
		}
		switch status {
		case "working", "idle", "done", "blocked", "unknown":
		default:
			status = "unknown"
		}
		seen[pane]++
		out.Agents = append(out.Agents, StatusAgent{pane, provider, status, tab, workspace})
	}
	// Pane is the consumer's unique key; discard every occurrence of ambiguity.
	safe := out.Agents[:0]
	for _, a := range out.Agents {
		if seen[a.Pane] != 1 {
			partial()
			continue
		}
		safe = append(safe, a)
	}
	out.Agents = safe
	sort.Slice(out.Agents, func(i, j int) bool { return out.Agents[i].Pane < out.Agents[j].Pane })
	if src.Status == "available" && len(out.Agents) == 0 {
		src.Status = "empty"
	}
	return out
}
