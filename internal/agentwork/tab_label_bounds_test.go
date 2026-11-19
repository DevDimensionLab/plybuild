package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

// tabProcessFixture answers Herdr's public commands in memory from files in
// root, without Herdr or provider access. The calls file and the max-active
// file make command, concurrency and workspace boundaries observable. A
// "block" file holds tab reads until the request budget ends; an "overlap"
// file holds the first tab read until a second one has started.
func tabProcessFixture(t *testing.T) (HerdrRuntime, string) {
	t.Helper()
	root := t.TempDir()
	var mu sync.Mutex
	active, maxActive := 0, 0
	overlapped := make(chan struct{})
	var once sync.Once
	answer := func(ctx context.Context, stdout, stderr io.Writer, args ...string) error {
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
			return fmt.Errorf("unsupported command %v", args)
		}
		mu.Lock()
		f, e := os.OpenFile(filepath.Join(root, "calls"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e == nil {
			_, _ = fmt.Fprintln(f, name)
			e = f.Close()
		}
		mu.Unlock()
		if e != nil {
			return e
		}
		if args[0] == "tab" {
			mu.Lock()
			active++
			if active > maxActive {
				maxActive = active
				_ = os.WriteFile(filepath.Join(root, "max-active"), []byte(strconv.Itoa(maxActive)), 0600)
			}
			if active > 1 {
				once.Do(func() { close(overlapped) })
			}
			mu.Unlock()
			defer func() { mu.Lock(); active--; mu.Unlock() }()
			if _, e := os.Stat(filepath.Join(root, "overlap")); e == nil {
				select {
				case <-overlapped:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			if _, e := os.Stat(filepath.Join(root, "block")); e == nil {
				<-ctx.Done()
				return ctx.Err()
			}
			if b, _ := os.ReadFile(filepath.Join(root, "stderr")); len(b) > 0 {
				if _, e := stderr.Write(b); e != nil {
					return e
				}
			}
		}
		b, e := os.ReadFile(filepath.Join(root, name))
		if e != nil {
			return e
		}
		_, e = stdout.Write(b)
		return e
	}
	return HerdrRuntime{Executable: "herdr", Timeout: time.Second, run: answer}, root
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
	putFixture(t, root, "overlap", "")
	got := h.tabLabels(agents)
	calls := fixtureCalls(t, root)
	sort.Strings(calls)
	if len(calls) != 32 || len(got) != 32 || calls[0] != "w00" || calls[31] != "w31" {
		t.Fatal("workspace bound/sharing", calls, len(got))
	}
	b, e := os.ReadFile(filepath.Join(root, "max-active"))
	if e != nil {
		t.Fatal(e)
	}
	max, _ := strconv.Atoi(string(b))
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
	for _, tc := range []string{"deadline", "stdout", "stderr", "invalid-utf8"} {
		t.Run(tc, func(t *testing.T) {
			h, root := tabProcessFixture(t)
			for w := 0; w < 32; w++ {
				id := fmt.Sprintf("w%02d", w)
				putFixture(t, root, id, labelEnvelope(labelRow(id, "t1", "Known")))
			}
			agents := []Agent{{Workspace: "w00", Tab: "t1"}}
			switch tc {
			case "deadline":
				putFixture(t, root, "block", "")
				for w := 1; w < 32; w++ {
					agents = append(agents, Agent{Workspace: fmt.Sprintf("w%02d", w), Tab: "t1"})
				}
			case "stdout":
				putFixture(t, root, "w00", labelEnvelope(labelRow("w00", "t1", "Known"))+strings.Repeat(" ", 256<<10))
			case "stderr":
				putFixture(t, root, "stderr", strings.Repeat("e", (256<<10)+1))
			case "invalid-utf8":
				putFixture(t, root, "w00", labelEnvelope(labelRow("w00", "t1", "Bad"))+string([]byte{0xff}))
			}
			// A blocked read returns only because the deadline ends it.
			if got := h.tabLabels(agents); len(got) != 0 {
				t.Fatal("bad source produced labels", got)
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
	// A real process: a descendant keeps both pipes open after the command
	// exits, until the test releases it. The read must end without waiting for
	// that release; otherwise it blocks and the test times out.
	t.Run("inherited-pipes", func(t *testing.T) {
		release := filepath.Join(t.TempDir(), "release")
		h := runtimeScript(t, `(while [ ! -e '`+release+`' ]; do sleep 0.05; done) &
printf '%s' '`+labelEnvelope(labelRow("w00", "t1", "Known"))+`'`)
		defer func() { _ = os.WriteFile(release, nil, 0600) }()
		if got := h.tabLabels([]Agent{{Workspace: "w00", Tab: "t1"}}); len(got) != 0 {
			t.Fatal("read waited for the descendant", got)
		}
		if _, e := os.Stat(release); e == nil {
			t.Fatal("descendant was released before the read ended")
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
			f, e := s.Register(s.onlyProject(), d)
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
