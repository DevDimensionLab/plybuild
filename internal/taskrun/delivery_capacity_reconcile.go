package taskrun

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func deliveryReconciliationScope(provider string) string {
	if provider == "codex" {
		return capacityRetirementScope
	}
	return quotaRetirementScope
}

// Codex's native item_completed carries both execution timestamps, even when
// the tool returned a running process earlier. A capacity error is a model-turn
// boundary, not a command completion. Examine the exact selected launch and the full tail, plus all previously
// reported running processes. Historical synchronous effects are not assessed.
// Terminal commands completed before capacity remain historical work, with
// their actual failures preserved. This proves completion, never that shell
// text was read-only or effects absent.
func deliveryValidateCapacityTranscript(s workflowState, in DeliveryQuotaReconcileInput, raw []byte) error {
	fail := func(reason string) error { return workflowError(4, "native Codex capacity proof: "+reason) }
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return fail("incomplete JSONL source")
	}
	if err := capacityValidateExecutorInventory(in, raw); err != nil {
		return err
	}
	if !plain(in.ProcessID, 1, 128) || !digestPattern.MatchString(in.CommandSHA256) {
		return fail("exact process and command digest required")
	}
	type item struct {
		Type    string          `json:"type"`
		ID      string          `json:"id"`
		Process string          `json:"process_id"`
		CWD     string          `json:"cwd"`
		Command json.RawMessage `json:"command"`
		Source  string          `json:"source"`
		Status  string          `json:"status"`
		Exit    *int            `json:"exit_code"`
		Kind    string          `json:"kind"`
		Output  string          `json:"aggregated_output"`
	}
	type payload struct {
		Arguments string          `json:"arguments"`
		Output    json.RawMessage `json:"output"`
		Type      string          `json:"type"`
		ID        string          `json:"id"`
		CWD       string          `json:"cwd"`
		Thread    string          `json:"thread_id"`
		Session   string          `json:"session_id"`
		Turn      string          `json:"turn_id"`
		RootTurn  string          `json:"root_turn_id"`
		Name      string          `json:"name"`
		Namespace string          `json:"namespace"`
		Call      string          `json:"call_id"`
		Error     json.RawMessage `json:"error"`
		Item      item            `json:"item"`
		Start     int64           `json:"started_at_ms"`
		End       int64           `json:"completed_at_ms"`
	}
	type row struct {
		Timestamp string  `json:"timestamp"`
		Type      string  `json:"type"`
		Payload   payload `json:"payload"`
	}
	var rows []row
	capacity, begin, meta := -1, -1, 0
	var prior time.Time
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var r row
		if err := json.Unmarshal(line, &r); err != nil {
			return fail("malformed native row")
		}
		tm, err := time.Parse(time.RFC3339Nano, r.Timestamp)
		if err != nil || tm.Before(prior) {
			return fail("missing or unordered source timestamp")
		}
		prior = tm
		p := r.Payload
		if p.Thread != "" && p.Thread != in.NativeSession || p.Session != "" && p.Session != in.NativeSession {
			return fail("wrong native session")
		}
		if r.Type == "session_meta" {
			meta++
			if meta != 1 || len(rows) != 0 || p.ID != in.NativeSession || p.CWD != s.Observed.Target.WorktreeLocator {
				return fail("wrong session metadata or Task cwd")
			}
		}
		if r.Timestamp == in.QuotaEvent {
			var e struct {
				Code string `json:"codex_error_info"`
			}
			if err := json.Unmarshal(p.Error, &e); err != nil || capacity >= 0 || r.Type != "event_msg" || p.Type != "task_complete" || e.Code != "server_overloaded" || p.Turn == "" || p.RootTurn != p.Turn {
				return fail("selected event is not a unique native server_overloaded turn")
			}
			capacity = len(rows)
			for i := len(rows) - 1; i >= 0; i-- {
				if rows[i].Type == "event_msg" && rows[i].Payload.Type == "task_started" && rows[i].Payload.Turn == p.Turn {
					begin = i
					break
				}
			}
		}
		rows = append(rows, r)
	}
	if capacity < 0 || begin < 0 || meta != 1 {
		return fail("missing capacity turn or session metadata")
	}
	// The audit window starts at the selected command's positively recorded
	// executor launch, not at a convenient final return. Older synchronous
	// results remain history; every reported running process in the full source
	// is still an obligation below, including launches before the selected turn.
	launch := -1
	for _, r := range rows {
		if r.Type != "response_item" || r.Payload.Type != "custom_tool_call_output" {
			continue
		}
		results, err := capacityExecutorResults(r.Payload.Output)
		if err != nil {
			continue
		}
		for _, result := range results {
			if result.Session == nil || strconv.FormatInt(*result.Session, 10) != in.ProcessID {
				continue
			}
			for j := begin; j < capacity; j++ {
				c := rows[j]
				if c.Type == "response_item" && c.Payload.Type == "custom_tool_call" && c.Payload.Name == "exec" && c.Payload.Call == r.Payload.Call {
					if launch < 0 || j < launch {
						launch = j
					}
				}
			}
		}
	}
	if launch < 0 {
		return fail("exact native running process launch is missing")
	}
	begin = launch
	// Yielded cells from earlier history must also have an exact completed wait.
	historicalCells := map[string]bool{}
	historicalWaits := map[string]string{}
	for _, r := range rows[:begin] {
		p := r.Payload
		if p.Type == "function_call" && p.Name == "wait" {
			var args struct {
				Cell      string `json:"cell_id"`
				Terminate bool   `json:"terminate"`
			}
			if json.Unmarshal([]byte(p.Arguments), &args) != nil || args.Terminate {
				return fail("unknown historical wait")
			}
			historicalWaits[p.Call] = args.Cell
		}
		if p.Type == "custom_tool_call_output" || p.Type == "function_call_output" {
			if cell := capacityYieldedCell(p.Output); cell != "" {
				historicalCells[cell] = true
			}
			if cell := historicalWaits[p.Call]; cell != "" && capacityYieldedCell(p.Output) == "" {
				if _, err := capacityExecutorResults(p.Output); err != nil {
					return fail("unresolved historical cell")
				}
				delete(historicalCells, cell)
			}
		}
	}
	if len(historicalCells) > 0 {
		return fail("unresolved executor cell before selected launch")
	}
	capTime, _ := time.Parse(time.RFC3339Nano, in.QuotaEvent)
	type call struct {
		parent     *call
		name       string
		start, end int64
		completed  bool
		commands   int
	}
	calls := map[string]*call{}
	cells := map[string]*call{}
	commands, processes, items := map[string]bool{}, map[string]bool{}, map[string]bool{}
	selected := false
	running := map[string]bool{}
	for index, r := range rows {
		p := r.Payload
		// No incomplete native execution, including one launched before this turn,
		// may be silently treated as historical completion.
		if r.Type == "event_msg" && p.Type == "item_completed" && p.Item.Type == "CommandExecution" {
			it := p.Item
			if it.ID == "" || it.Process == "" || commands[it.ID] || processes[it.Process] || p.Thread != in.NativeSession || p.Turn == "" || p.Start <= 0 || p.End < p.Start || (it.Status != "completed" && it.Status != "failed") || it.Exit == nil || it.Status == "completed" && *it.Exit != 0 || it.Status == "failed" && *it.Exit == 0 || it.Source != "unified_exec_startup" {
				return fail("unresolved, duplicate or conflicting command completion")
			}
			commands[it.ID], processes[it.Process] = true, true
			if index >= begin {
				cwd, err := url.Parse(it.CWD)
				if err != nil || cwd.Scheme != "file" || cwd.Host != "" || cwd.Path != s.Observed.Target.WorktreeLocator {
					return fail("command cwd differs from Task")
				}
				var argv []string
				if err := json.Unmarshal(it.Command, &argv); err != nil || len(argv) == 0 {
					return fail("missing native command identity")
				}
				if (index >= capacity || p.End >= capTime.UnixMilli()) && it.Status == "failed" && !capacityRejectedStopResolution(s, in, argv, it.Output, *it.Exit) {
					return fail("failed command is not the exact preserved stop-resolution rejection")
				}
				if p.End > mustCapacityRowTime(r.Timestamp) {
					return fail("completion timestamp precedes native command return")
				}
				if it.ID == in.LastToolResult {
					compact := new(bytes.Buffer)
					if err := json.Compact(compact, it.Command); err != nil {
						return fail("malformed command identity")
					}
					if selected || it.Process != in.ProcessID || hash(compact.Bytes()) != in.CommandSHA256 || it.Status != "completed" || p.Turn != rows[capacity].Payload.Turn || p.Start >= capTime.UnixMilli() || p.End <= capTime.UnixMilli() || index <= capacity {
						return fail("selected command must span capacity and complete successfully with exact identity")
					}
					selected = true
				}
			}
		}
		if r.Type == "response_item" && (p.Type == "custom_tool_call_output" || p.Type == "function_call_output" && calls[p.Call] != nil && calls[p.Call].name == "wait") {
			cell := capacityYieldedCell(p.Output)
			if cell != "" {
				if index < begin {
					continue
				}
				c := calls[p.Call]
				if c == nil {
					return fail("unbound executor cell")
				}
				if c.parent != nil {
					c = c.parent
				}
				if old := cells[cell]; old != nil && old != c {
					return fail("duplicate executor cell")
				}
				cells[cell] = c
			} else {
				results, err := capacityExecutorResults(p.Output)
				if err != nil {
					return fail("executor return " + p.Call + ": " + err.Error())
				}
				for _, result := range results {
					if result.Session != nil {
						running[strconv.FormatInt(*result.Session, 10)] = true
					}
				}
			}
		}
		if index < begin {
			continue
		}
		if len(p.Error) > 0 && string(p.Error) != "null" && index != capacity {
			return fail("later native error or permission denial remains unresolved")
		}
		switch r.Type {
		case "turn_context", "token_usage_record":
			if p.CWD != "" && p.CWD != s.Observed.Target.WorktreeLocator {
				return fail("later Task cwd changed")
			}
		case "response_item":
			switch p.Type {
			case "reasoning", "message":
			case "custom_tool_call", "function_call":
				if p.Call == "" || calls[p.Call] != nil {
					return fail("duplicate tool call")
				}
				name := ""
				if p.Type == "custom_tool_call" && p.Name == "exec" {
					name = "exec"
				}
				if p.Type == "function_call" && p.Namespace == "" && p.Name == "wait" {
					name = "wait"
				}
				if p.Type == "function_call" && p.Namespace == "clock" && p.Name == "sleep" {
					name = "clock.sleep"
				}
				if name == "" {
					return fail("unsupported later tool activity")
				}
				c := &call{name: name, start: mustCapacityRowTime(r.Timestamp)}
				if name == "wait" {
					var args struct {
						Cell      string `json:"cell_id"`
						Terminate bool   `json:"terminate"`
					}
					if json.Unmarshal([]byte(p.Arguments), &args) != nil || args.Terminate || cells[args.Cell] == nil || cells[args.Cell].completed {
						return fail("unbound or cancelling executor wait")
					}
					c.parent = cells[args.Cell]
				}
				calls[p.Call] = c
			case "custom_tool_call_output", "function_call_output":
				c := calls[p.Call]
				if c == nil || c.completed {
					return fail("unbound or duplicate tool return")
				}
				c.completed = true
				c.end = mustCapacityRowTime(r.Timestamp)
				if c.name == "exec" && capacityYieldedCell(p.Output) != "" {
					c.completed = false
				}
				if c.parent != nil {
					c.parent.end = c.end
					c.parent.completed = capacityYieldedCell(p.Output) == ""
				}
			default:
				return fail("unknown response activity")
			}
		case "event_msg":
			switch p.Type {
			case "task_started", "task_complete", "token_count", "thread_settings_applied":
			case "item_completed":
				it := p.Item
				if it.ID == "" || items[it.ID] || p.Thread != in.NativeSession {
					return fail("duplicate item or wrong session")
				}
				items[it.ID] = true
				switch it.Type {
				case "Reasoning", "AgentMessage", "UserMessage", "CommandExecution":
				case "Extension":
					c := calls[it.ID]
					if it.Kind != "clock.sleep" || c == nil || c.name != "clock.sleep" || p.Start <= 0 || p.End < p.Start {
						return fail("unresolved extension")
					}
					c.commands++
				default:
					return fail("unknown native effect")
				}
			default:
				return fail("unknown native event")
			}
		default:
			return fail("unknown source activity")
		}
	}
	// Bind each supported executor call to native completion evidence. Tool text
	// such as "process running" cannot substitute for a native command return.
	for _, r := range rows[begin:] {
		p := r.Payload
		if r.Type != "event_msg" || p.Type != "item_completed" || p.Item.Type != "CommandExecution" {
			continue
		}
		n := 0
		for _, c := range calls {
			if c.name == "exec" && p.Start >= c.start && p.Start <= c.end {
				c.commands++
				n++
			}
		}
		if n != 1 {
			// A process launched before the selected command can still be
			// active at capacity. Its earlier structured running result and
			// exact native process completion are both required.
			if n != 0 || p.Start >= mustCapacityRowTime(rows[begin].Timestamp) || !running[p.Item.Process] {
				return fail("command start does not bind exactly one executor call or earlier running process")
			}
		}
	}
	for id, c := range calls {
		if !c.completed || c.name != "wait" && c.commands == 0 {
			return fail("outstanding tool or executor without native completion: " + id + " (" + c.name + ")")
		}
	}
	for process := range running {
		if !processes[process] {
			return fail("started process lacks native completion: " + process)
		}
	}
	if !selected {
		return fail("exact eventual command completion missing")
	}
	return nil
}

func mustCapacityRowTime(timestamp string) int64 {
	t, _ := time.Parse(time.RFC3339Nano, timestamp)
	return t.UnixMilli()
}

// Only native executor envelopes are accepted here. Never inspect shell stdout
// for a claim of read-only behavior, completion, or permission. A running process
// ID is an obligation even when a sibling command already returned successfully.
type capacityExecutorResult struct {
	Chunk   string `json:"chunk_id"`
	Session *int64 `json:"session_id"`
	Exit    *int   `json:"exit_code"`
}

func capacityExecutorResults(raw json.RawMessage) ([]capacityExecutorResult, error) {
	fail := func() error { return workflowError(4, "unsupported or incomplete executor result") }
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &blocks); err != nil || len(blocks) < 2 {
		return nil, fail()
	}
	var results []capacityExecutorResult
	for i, block := range blocks {
		if block.Type != "input_text" {
			return nil, fail()
		}
		if i == 0 {
			if !strings.HasPrefix(block.Text, "Script completed\n") || !strings.HasSuffix(block.Text, "\nOutput:\n") {
				return nil, fail()
			}
			continue
		}
		var envelope struct {
			Status string          `json:"status"`
			Value  json.RawMessage `json:"value"`
		}
		data := []byte(block.Text)
		if err := json.Unmarshal(data, &envelope); err != nil {
			return nil, fail()
		}
		if envelope.Status != "" {
			if envelope.Status != "fulfilled" {
				return nil, fail()
			}
			data = envelope.Value
		}
		var result capacityExecutorResult
		if err := json.Unmarshal(data, &result); err != nil || result.Chunk == "" || (result.Session == nil && result.Exit == nil) || result.Session != nil && *result.Session <= 0 {
			return nil, fail()
		}
		results = append(results, result)
	}
	return results, nil
}

// A yielded orchestration cell is not a completed tool. Only its exact native
// wait/result can close it; cancellation or an unresolved cell remains blocked.
func capacityYieldedCell(raw json.RawMessage) string {
	var text string
	if json.Unmarshal(raw, &text) != nil {
		return ""
	}
	lines := strings.Split(text, "\n")
	if len(lines) != 4 || lines[2] != "Output:" || lines[3] != "" || !strings.HasPrefix(lines[1], "Wall time ") {
		return ""
	}
	id := strings.TrimPrefix(lines[0], "Script running with cell ID ")
	if id == lines[0] || !plain(id, 1, 128) {
		return ""
	}
	return id
}

// The sole failed command accepted in this window is the same owner's native
// resolve-stop rejection for this exact partial/unknown stop. This matches a
// single literal invocation and its exact native rejection, never a shell
// read-only classifier or a general allowance for nonzero exits.
func capacityRejectedStopResolution(s workflowState, in DeliveryQuotaReconcileInput, argv []string, output string, exit int) bool {
	if exit != 4 || len(argv) != 3 || argv[0] != "/bin/zsh" || argv[1] != "-lc" || s.Result.RunID == "" || s.Result.Delivery == nil || s.Result.Delivery.ActiveStop == nil {
		return false
	}
	if _, err := workflowEffectiveRuntime(s); err != nil {
		return false
	}
	// The transcript binds the actual invocation and native rejection. Ply's
	// historical executable path is not authority for this same owner's proof.
	control, command, ok := strings.Cut(argv[2], " ")
	if !ok || control == "" || strings.Contains(control, "=") {
		return false
	}
	prefix := "workflow execute resolve-stop " + s.Result.RunID + " --context " + s.Result.Paths.Context + " --stop " + in.StopSHA256 + " --resolution '"
	if !strings.HasPrefix(command, prefix) {
		return false
	}
	rest := strings.TrimPrefix(command, prefix)
	resolution, rest, ok := strings.Cut(rest, "'")
	if !ok || resolution == "" || !strings.HasPrefix(rest, " --action '") {
		return false
	}
	action := strings.TrimPrefix(rest, " --action '")
	if len(action) < 2 || action[len(action)-1] != '\'' || strings.Contains(action[:len(action)-1], "'") {
		return false
	}
	// Unquoted executable/identity fields must be literal shell words.
	for _, word := range []string{control, s.Result.RunID, s.Result.Paths.Context, in.StopSHA256} {
		if strings.ContainsAny(word, " \t\r\n;|&<>$`\\\"'(){}[]*?!#~") {
			return false
		}
	}
	return strings.HasPrefix(output, "Run: "+s.Result.RunID+"\nAgent: codex\n") && strings.Contains(output, "\nDelivery: stopped\n") && strings.Contains(output, "\nObserved stop: "+s.Result.Delivery.ActiveStop.Locator+"\n") && strings.HasSuffix(output, "Error: workflow_run: unknown or partial external effects remain unresolved; text and permission confirmation cannot clear this stop\n")
}

// Audit source identity must be unique, even though older resume callers only
// needed a resumable session locator. Never choose one of two matching journals.
func deliveryReconciliationSession(s workflowState) (string, error) {
	if s.Request.Runtime.Provider != "codex" {
		return deliveryProviderSession(s)
	}
	path, err := deliveryProviderSession(s)
	if err != nil {
		return "", err
	}
	count := 0
	err = filepath.WalkDir(filepath.Join(ownerResumeCodexHome(s), "sessions"), func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !e.IsDir() && strings.HasSuffix(e.Name(), "-"+s.Result.Transport.AgentSessionID+".jsonl") {
			count++
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if count != 1 {
		return "", workflowError(4, "ambiguous native Codex session source")
	}
	return path, nil
}

func WorkflowReconcileCapacityStop(d Dependencies, root, id string, in DeliveryQuotaReconcileInput) (WorkflowRun, error) {
	if in.Provider != "codex" {
		return WorkflowRun{}, workflowError(2, "capacity reconciliation requires the bound Codex provider")
	}
	return WorkflowReconcileQuotaStop(d, root, id, in)
}
