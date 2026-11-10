package agentwork

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

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
	TabLabel  string `json:"tab_label,omitempty"`
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

// StatusMetadata reads public membership and optional tab names, without native
// stores or execution identity. All transport shares one request budget.
func StatusMetadata(h HerdrRuntime) MetadataReadback {
	return StatusMetadataContext(context.Background(), h)
}

func StatusMetadataContext(parent context.Context, h HerdrRuntime) MetadataReadback {
	if h.Timeout <= 0 || h.Timeout > MetadataTimeout {
		h.Timeout = MetadataTimeout
	}
	// Reserve pipe cleanup inside the shared budget, not after a second timeout.
	transport := h.Timeout - tabLabelCleanup
	if transport < 0 {
		transport = 0
	}
	ctx, cancel := context.WithTimeout(parent, transport)
	defer cancel()
	raw, err := h.callBytesContext(ctx, tabLabelCleanup, "agent", "list")
	observed := time.Now().UTC().Format(time.RFC3339Nano)
	out := MetadataReadback{Kind: MetadataKind, SchemaVersion: 1, ObservedAtUTC: observed, Agents: []StatusAgent{}, Sources: []MetadataSource{{Source: "herdr.agent.list", Status: "unavailable", Freshness: "unknown", Coverage: "unknown", Scope: metadataScope}}}
	if err != nil {
		return out
	}
	out = decodeMetadata(raw, out)
	agents := make([]Agent, 0, len(out.Agents))
	for _, a := range out.Agents {
		if a.Workspace != "" && a.Tab != "" {
			agents = append(agents, Agent{Workspace: a.Workspace, Tab: a.Tab})
		}
	}
	labels := h.tabLabelsContext(ctx, agents)
	for i := range out.Agents {
		a := &out.Agents[i]
		if label := labels[tabKey{a.Workspace, a.Tab}]; validPublicTabLabel(label) {
			a.TabLabel = label
		}
	}
	return out
}

func validPublicTabLabel(label string) bool {
	if privateText(label) || !utf8.ValidString(label) || !plain(label) || strings.ContainsFunc(label, func(r rune) bool {
		return unicode.Is(unicode.Cf, r) || r == '\u2028' || r == '\u2029' || r == '\ufffd'
	}) {
		return false
	}
	lower := strings.ToLower(label)
	for _, marker := range []string{"session_id", "terminal_id", "cwd=", "token=", "secret=", "password="} {
		if strings.Contains(lower, marker) {
			return false
		}
	}
	return true
}

// Decode only selected object members. Execution fields are opaque raw JSON;
// their types and duplicate members never qualify metadata completeness.
func metadataObject(raw []byte, selected map[string]bool, observe func(string, json.RawMessage)) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	tok, e := d.Token()
	if e != nil || tok != json.Delim('{') {
		return nil, fmt.Errorf("metadata_object_invalid")
	}
	out := map[string]json.RawMessage{}
	duplicate := false
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
		if observe != nil {
			observe(key, v)
		}
		if selected == nil || selected[key] {
			if _, exists := out[key]; exists {
				// Keep scanning: a malformed member cannot hide later pane claims.
				duplicate = true
			} else {
				out[key] = v
			}
		}
	}
	if _, e = d.Token(); e != nil {
		return nil, e
	}
	if _, e = d.Token(); e != io.EOF {
		return nil, fmt.Errorf("metadata_trailing")
	}
	if duplicate {
		return nil, fmt.Errorf("metadata_duplicate")
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
	env, e := metadataObject(raw, nil, nil)
	if e != nil {
		return out
	}
	if v, exists := env["error"]; exists && string(bytes.TrimSpace(v)) != "null" {
		return out
	}
	result, e := metadataObject(env["result"], nil, nil)
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
	for key := range env {
		if key != "id" && key != "result" && key != "error" {
			partial()
		}
	}
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
		// Every valid pane claim participates in uniqueness, even when another
		// field or duplicate key makes this row ineligible for publication.
		row, e := metadataObject(rawRow, selected, func(key string, value json.RawMessage) {
			if key == "pane_id" {
				if pane, ok := metadataText(value, true); ok {
					seen[pane]++
				}
			}
		})
		if e != nil {
			partial()
			continue
		}
		if raw, exists := row["agent"]; !exists || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
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
		out.Agents = append(out.Agents, StatusAgent{Pane: pane, Provider: provider, Status: status, Tab: tab, Workspace: workspace})
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
