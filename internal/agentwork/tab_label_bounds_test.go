package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func labelEnvelope(rows ...any) string {
	b, _ := json.Marshal(map[string]any{"result": map[string]any{"tabs": rows}})
	return string(b)
}
func labelRow(workspace, tab, label string) map[string]any {
	return map[string]any{"workspace_id": workspace, "tab_id": tab, "label": label}
}

func TestTabLabelsValidationAndConflict(t *testing.T) {
	good := labelRow("w1", "t1", "Actual 🌱")
	for _, tc := range []struct{ name, raw, want string }{
		{"valid", labelEnvelope(good), "Actual 🌱"},
		{"256-codepoints", labelEnvelope(labelRow("w1", "t1", strings.Repeat("🌱", 256))), strings.Repeat("🌱", 256)},
		{"257-codepoints", labelEnvelope(labelRow("w1", "t1", strings.Repeat("🌱", 257))), ""},
		{"empty", labelEnvelope(labelRow("w1", "t1", "")), ""},
		{"whitespace", labelEnvelope(labelRow("w1", "t1", " \u2003 ")), ""},
		{"newline", labelEnvelope(labelRow("w1", "t1", "bad\nlabel")), ""},
		{"ansi", labelEnvelope(labelRow("w1", "t1", "\x1b[31mred")), ""},
		{"c1", labelEnvelope(labelRow("w1", "t1", "bad\u009blabel")), ""},
		{"emoji-joiner", labelEnvelope(labelRow("w1", "t1", "👩‍💻")), "👩‍💻"},
		{"duplicate", labelEnvelope(good, good), ""},
		{"conflicting", labelEnvelope(good, labelRow("w1", "t1", "Renamed")), ""},
		{"invalid-duplicate", labelEnvelope(good, map[string]any{"workspace_id": "w1", "tab_id": "t1", "label": 42}), ""},
		{"invalid-first", labelEnvelope(map[string]any{"workspace_id": "w1", "tab_id": "t1", "label": 42}, good), ""},
		{"triplicate", labelEnvelope(good, good, good), ""},
		{"workspace-mismatch", labelEnvelope(labelRow("w2", "t1", "Other")), ""},
		{"missing-workspace", labelEnvelope(map[string]any{"tab_id": "t1", "label": "Other"}), ""},
		{"tab-mismatch", labelEnvelope(labelRow("w1", "t2", "Other")), ""},
		{"terminal-title", labelEnvelope(map[string]any{"workspace_id": "w1", "tab_id": "t1", "terminal_title": "Wrong"}), ""},
		{"malformed", "{", ""}, {"null-tabs", `{"result":{"tabs":null}}`, ""},
		{"missing-tabs", `{"result":{}}`, ""}, {"error", `{"result":{"tabs":[]},"error":{"code":"fail"}}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := decodeTabLabels(context.Background(), []byte(tc.raw), "w1")[tabKey{"w1", "t1"}]
			if got != tc.want {
				t.Fatalf("label %q want %q", got, tc.want)
			}
		})
	}
	rows := []any{good, labelRow("w1", "t2", "Kept"), good}
	if got := decodeTabLabels(context.Background(), []byte(labelEnvelope(rows...)), "w1"); len(got) != 1 || got[tabKey{"w1", "t2"}] != "Kept" {
		t.Fatal("conflict removed unrelated tab", got)
	}
	for _, count := range []int{4096, 4097} {
		rows := make([]any, count)
		for n := range rows {
			rows[n] = labelRow("w1", fmt.Sprintf("t%d", n), "Known")
		}
		got := decodeTabLabels(context.Background(), []byte(labelEnvelope(rows...)), "w1")
		if count == 4096 && len(got) != 4096 || count == 4097 && len(got) != 0 {
			t.Fatalf("%d tabs accepted %d", count, len(got))
		}
	}
}

// The public-command fixture is a real subprocess, with no Herdr/provider
// access. Separate log files make concurrency, byte and command boundaries
// observable without production injection hooks.
func TestTabLabelProcessChild(t *testing.T) {
	root := os.Getenv("PLY_TAB_LABEL_TEST_ROOT")
	if root == "" {
		t.Skip("test-child entrypoint")
	}
	args := os.Args
	for n, a := range args {
		if a == "--" {
			args = args[n+1:]
			break
		}
	}
	if len(args) >= 3 && args[0] == "agent" && (args[1] == "get" || args[1] == "prompt") {
		_ = os.WriteFile(filepath.Join(root, args[1]+"-target"), []byte(args[2]), 0600)
	}
	var name string
	switch {
	case len(args) == 2 && args[0] == "agent" && args[1] == "list":
		name = "agents"
	case len(args) == 3 && args[0] == "agent" && args[1] == "get":
		name = "agent"
	case len(args) == 4 && args[0] == "agent" && args[1] == "prompt":
		name = "prompt"
	case len(args) == 4 && args[0] == "tab" && args[1] == "list" && args[2] == "--workspace":
		name = args[3]
	default:
		os.Exit(5)
	}
	f, e := os.OpenFile(filepath.Join(root, "calls"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		os.Exit(6)
	}
	_, _ = fmt.Fprintln(f, name)
	_ = f.Close()
	if name != "agents" && name != "agent" && name != "prompt" {
		_ = os.WriteFile(filepath.Join(root, name+".start"), []byte(strconv.FormatInt(time.Now().UnixNano(), 10)), 0600)
		delay, _ := os.ReadFile(filepath.Join(root, "delay"))
		ms, _ := strconv.Atoi(string(delay))
		time.Sleep(time.Duration(ms) * time.Millisecond)
		if b, _ := os.ReadFile(filepath.Join(root, "stderr")); len(b) > 0 {
			_, _ = os.Stderr.Write(b)
		}
		if _, e := os.Stat(filepath.Join(root, "hold-pipes")); e == nil {
			// A descendant holds both pipes after the command exits. WaitDelay must
			// bound the parent's drain even when the metadata deadline has not fired.
			c := exec.Command("/bin/sleep", "2")
			c.Stdout, c.Stderr = os.Stdout, os.Stderr
			_ = c.Start()
		}
	}
	b, e := os.ReadFile(filepath.Join(root, name))
	if e != nil {
		os.Exit(7)
	}
	_, _ = os.Stdout.Write(b)
	_ = os.WriteFile(filepath.Join(root, name+".end"), []byte(strconv.FormatInt(time.Now().UnixNano(), 10)), 0600)
	os.Exit(0)
}
func tabProcessFixture(t *testing.T) (HerdrRuntime, string) {
	t.Helper()
	root := t.TempDir()
	t.Setenv("PLY_TAB_LABEL_TEST_ROOT", root)
	t.Setenv("GORACE", "atexit_sleep_ms=0")
	binary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	h := runtimeScript(t, "exec '"+strings.ReplaceAll(binary, "'", "'\\''")+"' -test.run='^TestTabLabelProcessChild$' -- \"$@\"\n")
	return h, root
}
func putFixture(t *testing.T, root, name, body string) {
	t.Helper()
	if e := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); e != nil {
		t.Fatal(e)
	}
}
func fixtureCalls(t *testing.T, root string) []string {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(root, "calls"))
	if e != nil {
		t.Fatal(e)
	}
	return strings.Fields(string(b))
}

func TestTabLabelsWorkspaceSharingFreshnessAndConcurrency(t *testing.T) {
	h, root := tabProcessFixture(t)
	var agents []Agent
	// Reverse the input to prove deterministic workspace selection.
	for w := 33; w >= 0; w-- {
		for a := 0; a < 2; a++ {
			agents = append(agents, Agent{Workspace: fmt.Sprintf("w%02d", w), Tab: "t1"})
		}
	}
	for w := 0; w < 34; w++ {
		id := fmt.Sprintf("w%02d", w)
		putFixture(t, root, id, labelEnvelope(labelRow(id, "t1", "First")))
	}
	putFixture(t, root, "delay", "20")
	got := h.tabLabels(agents)
	calls := fixtureCalls(t, root)
	sort.Strings(calls)
	if len(calls) != 32 || len(got) != 32 || calls[0] != "w00" || calls[31] != "w31" {
		t.Fatal("workspace bound/sharing", calls, len(got))
	}
	type event struct {
		at     int64
		change int
	}
	var events []event
	for _, id := range calls {
		for _, v := range []struct {
			suffix string
			change int
		}{{"start", 1}, {"end", -1}} {
			b, e := os.ReadFile(filepath.Join(root, id+"."+v.suffix))
			if e != nil {
				t.Fatal(e)
			}
			at, _ := strconv.ParseInt(string(b), 10, 64)
			events = append(events, event{at, v.change})
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].at < events[j].at })
	active, max := 0, 0
	for _, e := range events {
		active += e.change
		if active > max {
			max = active
		}
	}
	if max < 2 || max > 4 {
		t.Fatalf("parallel calls %d", max)
	}
	putFixture(t, root, "w00", labelEnvelope(labelRow("w00", "t1", "Renamed")))
	got = h.tabLabels([]Agent{{Workspace: "w00", Tab: "t1"}, {Workspace: "w00", Tab: "t1"}})
	if got[tabKey{"w00", "t1"}] != "Renamed" || len(fixtureCalls(t, root)) != 33 {
		t.Fatal("request cache or repeated workspace read", got)
	}
}

func TestTabLabelsDeadlineBytesAndInheritedPipes(t *testing.T) {
	for _, tc := range []string{"deadline", "stdout", "stderr", "inherited-pipes", "invalid-utf8"} {
		t.Run(tc, func(t *testing.T) {
			h, root := tabProcessFixture(t)
			for w := 0; w < 32; w++ {
				id := fmt.Sprintf("w%02d", w)
				putFixture(t, root, id, labelEnvelope(labelRow(id, "t1", "Known")))
			}
			agents := []Agent{{Workspace: "w00", Tab: "t1"}}
			switch tc {
			case "deadline":
				putFixture(t, root, "delay", "3000")
				for w := 1; w < 32; w++ {
					agents = append(agents, Agent{Workspace: fmt.Sprintf("w%02d", w), Tab: "t1"})
				}
			case "stdout":
				putFixture(t, root, "w00", labelEnvelope(labelRow("w00", "t1", "Known"))+strings.Repeat(" ", 256<<10))
			case "stderr":
				putFixture(t, root, "stderr", strings.Repeat("e", (256<<10)+1))
			case "inherited-pipes":
				putFixture(t, root, "hold-pipes", "")
			case "invalid-utf8":
				putFixture(t, root, "w00", labelEnvelope(labelRow("w00", "t1", "Bad"))+string([]byte{0xff}))
			}
			start := time.Now()
			got := h.tabLabels(agents)
			elapsed := time.Since(start)
			if len(got) != 0 {
				t.Fatal("bad source produced labels", got)
			}
			if elapsed > 1500*time.Millisecond {
				t.Fatalf("unbounded metadata read %s", elapsed)
			}
			calls := fixtureCalls(t, root)
			if tc == "deadline" && len(calls) > 4 {
				t.Fatal("queued calls started after deadline", calls)
			}
			if tc != "deadline" && len(calls) != 1 {
				t.Fatal("retried metadata", calls)
			}
		})
	}
	// Exactly the byte boundary remains usable (including harmless stderr).
	t.Run("exact-byte-boundary", func(t *testing.T) {
		h, root := tabProcessFixture(t)
		b := labelEnvelope(labelRow("w00", "t1", "Known"))
		putFixture(t, root, "w00", b+strings.Repeat(" ", (256<<10)-len(b)))
		putFixture(t, root, "stderr", strings.Repeat("e", 256<<10))
		if got := h.tabLabels([]Agent{{Workspace: "w00", Tab: "t1"}}); got[tabKey{"w00", "t1"}] != "Known" {
			t.Fatal(got)
		}
	})
}

func TestOverviewTabMetadataFailuresPreserveBase(t *testing.T) {
	s, _, _ := fixture(t)
	i := s.Runtime.(*fakeRuntime).agents[0].Identity
	h, root := tabProcessFixture(t)
	a := map[string]any{"agent": i.Provider, "agent_status": "idle", "workspace_id": i.Workspace, "tab_id": i.Tab, "pane_id": i.Pane, "terminal_id": i.Terminal, "agent_session": map[string]string{"agent": i.Provider, "kind": "id", "value": i.Session}}
	b, _ := json.Marshal(map[string]any{"result": map[string]any{"agents": []any{a}}})
	putFixture(t, root, "agents", string(b))
	s.Runtime = h
	baseline, e := s.Overview()
	if e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{"{", `{"result":{"tabs":null}}`, labelEnvelope(labelRow(i.Workspace, i.Tab, "bad\nlabel")), labelEnvelope(labelRow(i.Workspace, i.Tab, "One"), labelRow(i.Workspace, i.Tab, "Two"))} {
		putFixture(t, root, i.Workspace, raw)
		got, e := s.Overview()
		if e != nil || !reflect.DeepEqual(baseline, got) {
			t.Fatal("metadata failure changed base/source/freshness/capabilities", got, e)
		}
	}
}

func TestTabLabelsRejectInvalidUnicodeEscapes(t *testing.T) {
	for _, tc := range []struct{ encoded, want string }{
		{`"\ud83c\udf31"`, "🌱"}, {`"\ufffd"`, "�"},
		{`"\ud83c"`, ""}, {`"\udf31"`, ""}, {`"bad\ud83c text"`, ""},
		{`"\\ud83c"`, `\ud83c`},
	} {
		raw := `{"result":{"tabs":[{"workspace_id":"w1","tab_id":"t1","label":` + tc.encoded + `}]}}`
		if got := decodeTabLabels(context.Background(), []byte(raw), "w1")[tabKey{"w1", "t1"}]; got != tc.want {
			t.Errorf("%s gave %q want %q", tc.encoded, got, tc.want)
		}
	}
}

func TestTabLabelRenamePreservesDispatchAndFreshVeto(t *testing.T) {
	for _, state := range []string{"clear", "same-session-veto", "unknown-coverage", "new-session"} {
		t.Run(state, func(t *testing.T) {
			s, d, _ := fixture(t)
			i := s.Runtime.(*fakeRuntime).agents[0].Identity
			h, root := tabProcessFixture(t)
			wire := map[string]any{"agent": i.Provider, "agent_status": "idle", "interactive_ready": true, "workspace_id": i.Workspace, "tab_id": i.Tab, "pane_id": i.Pane, "terminal_id": i.Terminal, "agent_session": map[string]string{"agent": i.Provider, "kind": "id", "value": i.Session}}
			b, _ := json.Marshal(map[string]any{"result": map[string]any{"agents": []any{wire}}})
			putFixture(t, root, "agents", string(b))
			b, _ = json.Marshal(map[string]any{"result": map[string]any{"agent": wire}})
			putFixture(t, root, "agent", string(b))
			putFixture(t, root, "prompt", `{"result":{"type":"agent_prompted"}}`)
			putFixture(t, root, i.Workspace, labelEnvelope(labelRow(i.Workspace, i.Tab, "Before")))
			s.Runtime = h
			s.Send = h.Prompt
			d.NextCheckAtUTC = s.now()
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			before, e := s.Overview()
			if e != nil {
				t.Fatal(e)
			}
			putFixture(t, root, i.Workspace, labelEnvelope(labelRow(i.Workspace, i.Tab, "After")))
			after, e := s.Overview()
			if e != nil {
				t.Fatal(e)
			}
			if after.Agents[0].TabLabel != "After" || before.Agents[0].AgentRef != after.Agents[0].AgentRef || !reflect.DeepEqual(before.Followups, after.Followups) {
				t.Fatal("rename changed followup identity/generation/claim")
			}
			native := s.native()
			var reads atomic.Int32
			s.Native = func() NativeSnapshot { reads.Add(1); return native }
			switch state {
			case "same-session-veto":
				native.Vetoes = map[string]string{i.sessionKey(): "native_refusal"}
			case "unknown-coverage":
				native.CoverageUnknown = true
				native.Partial = true
			case "new-session":
				wire["agent_session"] = map[string]string{"agent": i.Provider, "kind": "id", "value": "replacement"}
				b, _ = json.Marshal(map[string]any{"result": map[string]any{"agent": wire}})
				putFixture(t, root, "agent", string(b))
			}
			callsBefore := fixtureCalls(t, root)
			if _, e := h.Get(i); e != nil {
				t.Fatal(e)
			}
			var workers sync.WaitGroup
			for n := 0; n < 6; n++ {
				workers.Add(1)
				go func() {
					defer workers.Done()
					if _, e := s.Dispatch(f.ID); e != nil {
						t.Error(e)
					}
				}()
			}
			workers.Wait()
			got, e := s.Show(f.ID)
			if e != nil {
				t.Fatal(e)
			}
			calls := fixtureCalls(t, root)
			prompts := 0
			for _, call := range calls[len(callsBefore):] {
				if call == i.Workspace {
					t.Fatal("Get/dispatch fetched metadata")
				}
				if call == "prompt" {
					prompts++
				}
			}
			if state == "clear" {
				if got.Transport.State != "sent" || prompts != 1 || got.TargetRef != f.TargetRef || got.Generation != f.Generation || reads.Load() == 0 {
					t.Fatal("rename changed transport/claim/no-replay", got, prompts, reads.Load())
				}
				for _, op := range []string{"get", "prompt"} {
					b, e := os.ReadFile(filepath.Join(root, op+"-target"))
					if e != nil || string(b) != i.Pane {
						t.Fatal("transport target changed", op, string(b), e)
					}
				}
			} else if got.Transport.State != "not_sent" || prompts != 0 {
				t.Fatal("stale authority dispatched", state, got, prompts)
			}
		})
	}
}
