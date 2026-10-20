package taskrun

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type capacityInvocation struct {
	Command string
	Session int64
}

// This recognizes two closed syntax forms, not JavaScript in general. Every
// tool result must be printed directly, with literal arguments and no other
// executable statements. Even an unprinted executor call makes the source
// unsupported. Shell contents are opaque command identity, never effect claims.
func capacityExecutorInput(source string) ([]capacityInvocation, error) {
	p := capacityLiteralParser{rest: strings.TrimSpace(source)}
	var out []capacityInvocation
	if p.take("const results") {
		if !p.take("=") || !p.take("await Promise.allSettled(") || !p.take("[") {
			return nil, p.fail()
		}
		for {
			call, err := p.tool()
			if err != nil {
				return nil, err
			}
			out = append(out, call)
			if p.take("]") {
				break
			}
			if !p.take(",") {
				return nil, p.fail()
			}
		}
		if !p.take(");") || !p.take("results.forEach(text);") || p.rest != "" {
			return nil, p.fail()
		}
	} else {
		for p.rest != "" {
			if !p.take("text(await ") {
				return nil, p.fail()
			}
			call, err := p.tool()
			if err != nil {
				return nil, err
			}
			out = append(out, call)
			if !p.take(");") {
				return nil, p.fail()
			}
		}
	}
	if len(out) == 0 {
		return nil, p.fail()
	}
	return out, nil
}

type capacityLiteralParser struct{ rest string }

func (p *capacityLiteralParser) fail() error {
	return workflowError(4, "unsupported executor source; complete literal invocation/result inventory required")
}
func (p *capacityLiteralParser) take(token string) bool {
	p.rest = strings.TrimSpace(p.rest)
	if !strings.HasPrefix(p.rest, token) {
		return false
	}
	p.rest = strings.TrimSpace(strings.TrimPrefix(p.rest, token))
	return true
}
func (p *capacityLiteralParser) literal() (json.RawMessage, error) {
	d := json.NewDecoder(strings.NewReader(p.rest))
	var v json.RawMessage
	if err := d.Decode(&v); err != nil {
		return nil, p.fail()
	}
	p.rest = strings.TrimSpace(p.rest[d.InputOffset():])
	return v, nil
}
func (p *capacityLiteralParser) tool() (capacityInvocation, error) {
	var out capacityInvocation
	execute := p.take("tools.exec_command(")
	if !execute && !p.take("tools.write_stdin(") {
		return out, p.fail()
	}
	if !p.take("{") {
		return out, p.fail()
	}
	seen := map[string]bool{}
	for {
		key := ""
		if strings.HasPrefix(p.rest, "\"") {
			b, err := p.literal()
			if err != nil || json.Unmarshal(b, &key) != nil {
				return out, p.fail()
			}
		} else {
			for i, c := range p.rest {
				if !(c >= 'a' && c <= 'z' || c == '_') {
					key = p.rest[:i]
					p.rest = strings.TrimSpace(p.rest[i:])
					break
				}
			}
		}
		if key == "" || seen[key] || !p.take(":") {
			return out, p.fail()
		}
		seen[key] = true
		b, err := p.literal()
		if err != nil {
			return out, err
		}
		switch key {
		case "cmd":
			if !execute || json.Unmarshal(b, &out.Command) != nil || out.Command == "" {
				return out, p.fail()
			}
		case "session_id":
			if execute || json.Unmarshal(b, &out.Session) != nil || out.Session <= 0 {
				return out, p.fail()
			}
		case "chars":
			var chars string
			if execute || json.Unmarshal(b, &chars) != nil || chars != "" {
				return out, p.fail()
			}
		case "max_output_tokens", "yield_time_ms":
			var n int
			if json.Unmarshal(b, &n) != nil || n < 0 {
				return out, p.fail()
			}
		default:
			return out, p.fail()
		}
		if p.take("}") {
			break
		}
		if !p.take(",") {
			return out, p.fail()
		}
	}
	if !p.take(")") || execute && out.Command == "" || !execute && out.Session == 0 {
		return out, p.fail()
	}
	return out, nil
}

// Check every orchestration call, including history before the capacity turn.
// The grammar and one-to-one result count close the script-selected-output gap;
// all reported running process IDs still need actual native completion.
func capacityValidateExecutorInventory(in DeliveryQuotaReconcileInput, raw []byte) error {
	type command struct {
		Type    string   `json:"type"`
		ID      string   `json:"id"`
		Kind    string   `json:"kind"`
		Process string   `json:"process_id"`
		Argv    []string `json:"command"`
		Exit    *int     `json:"exit_code"`
	}
	type payload struct {
		Namespace string          `json:"namespace"`
		Type      string          `json:"type"`
		Call      string          `json:"call_id"`
		Name      string          `json:"name"`
		Input     string          `json:"input"`
		Arguments string          `json:"arguments"`
		Output    json.RawMessage `json:"output"`
		Item      command         `json:"item"`
		Start     int64           `json:"started_at_ms"`
		End       int64           `json:"completed_at_ms"`
	}
	type row struct {
		Timestamp string  `json:"timestamp"`
		Type      string  `json:"type"`
		Payload   payload `json:"payload"`
	}
	type script struct {
		specs      []capacityInvocation
		results    []capacityExecutorResult
		start, end int64
		done       bool
	}
	scripts := map[string]*script{}
	cells := map[string]*script{}
	waits := map[string]*script{}
	var commands []row
	chunks := map[string]bool{}
	callIDs := map[string]bool{}
	functionReturns := map[string]bool{}
	sleeps := map[string]bool{}
	extensions := map[string]bool{}
	fail := func() error {
		return workflowError(4, "executor source/result/native command inventory is incomplete or conflicting")
	}
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var r row
		if json.Unmarshal(line, &r) != nil {
			return fail()
		}
		p := r.Payload
		if r.Type == "event_msg" && p.Type == "item_completed" && p.Item.Type == "Extension" {
			if p.Item.Kind != "clock.sleep" || extensions[p.Item.ID] || p.Start <= 0 || p.End < p.Start {
				return fail()
			}
			extensions[p.Item.ID] = true
		}
		if r.Type == "event_msg" && p.Type == "item_completed" && p.Item.Process != "" {
			commands = append(commands, r)
		}
		if r.Type != "response_item" {
			continue
		}
		if p.Type == "custom_tool_call" || p.Type == "function_call" {
			if p.Call == "" || callIDs[p.Call] {
				return fail()
			}
			callIDs[p.Call] = true
		}
		if p.Type == "custom_tool_call" {
			if p.Name != "exec" || scripts[p.Call] != nil {
				return fail()
			}
			specs, err := capacityExecutorInput(p.Input)
			if err != nil {
				return err
			}
			scripts[p.Call] = &script{specs: specs, start: mustCapacityRowTime(r.Timestamp)}
		}
		if p.Type == "function_call" {
			if !(p.Name == "wait" && p.Namespace == "" || p.Name == "sleep" && p.Namespace == "clock") {
				return fail()
			}
			functionReturns[p.Call] = false
			if p.Name == "sleep" {
				sleeps[p.Call] = true
			}
		}
		if p.Type == "function_call_output" {
			returned, ok := functionReturns[p.Call]
			if !ok || returned {
				return fail()
			}
			functionReturns[p.Call] = true
		}
		if p.Type == "function_call" && p.Name == "wait" {
			var args struct {
				Cell      string `json:"cell_id"`
				Terminate bool   `json:"terminate"`
			}
			if json.Unmarshal([]byte(p.Arguments), &args) != nil || args.Terminate || cells[args.Cell] == nil || cells[args.Cell].done || waits[p.Call] != nil {
				return fail()
			}
			waits[p.Call] = cells[args.Cell]
		}
		if p.Type != "custom_tool_call_output" && !(p.Type == "function_call_output" && waits[p.Call] != nil) {
			continue
		}
		c := scripts[p.Call]
		if c == nil {
			c = waits[p.Call]
		}
		if c == nil || c.done {
			return fail()
		}
		if cell := capacityYieldedCell(p.Output); cell != "" {
			if old := cells[cell]; old != nil && old != c {
				return fail()
			}
			cells[cell] = c
			continue
		}
		results, err := capacityExecutorResults(p.Output)
		if err != nil {
			return err
		}
		for _, result := range results {
			if chunks[result.Chunk] {
				return fail()
			}
			chunks[result.Chunk] = true
		}
		c.results = append(c.results, results...)
		c.end = mustCapacityRowTime(r.Timestamp)
		c.done = true
		if len(c.results) != len(c.specs) {
			return fail()
		}
	}
	for id, returned := range functionReturns {
		if !returned || sleeps[id] && !extensions[id] {
			return fail()
		}
	}
	used := map[int]bool{}
	capacity := mustCapacityRowTime(in.QuotaEvent)
	// Map iteration order cannot decide which native command belongs to a script;
	// a command's start is within that script's exact sourced interval.
	for _, c := range scripts {
		if !c.done {
			return fail()
		}
		for i, spec := range c.specs {
			result := c.results[i]
			if spec.Session != 0 {
				if result.Session != nil && *result.Session != spec.Session {
					return fail()
				}
				found := false
				for _, r := range commands {
					if r.Payload.Item.Process == strconv.FormatInt(spec.Session, 10) && r.Payload.Item.Exit != nil {
						if result.Exit != nil && (*result.Exit != *r.Payload.Item.Exit || r.Payload.End > c.end) {
							return fail()
						}
						found = true
					}
				}
				if !found {
					return fail()
				}
				continue
			}
			match := -1
			for j, r := range commands {
				p := r.Payload
				if used[j] || len(p.Item.Argv) == 0 || p.Item.Argv[len(p.Item.Argv)-1] != spec.Command || p.Start < c.start || p.Start > c.end {
					continue
				}
				if result.Session != nil && p.Item.Process != strconv.FormatInt(*result.Session, 10) {
					continue
				}
				if result.Exit != nil && (p.Item.Exit == nil || *p.Item.Exit != *result.Exit || p.End > c.end) {
					continue
				}
				if match >= 0 {
					return fail()
				}
				match = j
			}
			if match < 0 {
				// A terminal nonzero executor result can lack a CommandExecution
				// item. Before capacity it remains a historical failure, never proof
				// of no effects. Running or successful missing returns fail.
				if result.Session != nil || result.Exit == nil || *result.Exit == 0 || c.end >= capacity {
					return fail()
				}
			} else {
				used[match] = true
			}
		}
	}
	return nil
}
