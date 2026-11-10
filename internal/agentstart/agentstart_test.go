package agentstart

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

var errCrash = errors.New("simulated process death")

type fakeHerdr struct {
	t       *testing.T
	receipt string
	mu      sync.Mutex
	calls   [][]string
	log     *[]string
	drop    map[string]bool
	name    string
	kind    string
	never   bool
	working bool
	agent   func(a map[string]any)
	tab     string
	// startFails makes agent start fail after Herdr named the pane, as when
	// Claude waits on its folder trust dialog past the start timeout.
	startFails bool
}

func newFake(t *testing.T, in StartInput) *fakeHerdr {
	return &fakeHerdr{t: t, receipt: in.ReceiptPath, drop: map[string]bool{}}
}

func stageOf(args []string) string {
	switch {
	case args[0] == "tab" && args[1] == "create":
		return StageTabCreate
	case args[0] == "agent" && args[1] == "start":
		return StageAgentStart
	case args[0] == "agent" && args[1] == "prompt":
		return StagePrompt
	}
	return ""
}

func (f *fakeHerdr) call(_ context.Context, _ HerdrBinary, args []string) (json.RawMessage, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, args)
	stage := stageOf(args)
	if f.log != nil {
		if stage != "" {
			*f.log = append(*f.log, stage)
		} else {
			*f.log = append(*f.log, args[0]+" "+args[1])
		}
	}
	if stage != "" {
		got, err := ReadReceipt(f.receipt)
		if err != nil || got.Pending != stage {
			f.t.Errorf("effect %s before durable pending marker: %+v %v", stage, got, err)
		}
	}
	if stage != "" && f.drop[stage] {
		return nil, errors.New("reply lost")
	}
	switch stage {
	case StageTabCreate:
		if f.tab != "" {
			return json.RawMessage(f.tab), nil
		}
		return json.RawMessage(`{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","terminal_id":"term-2","workspace_id":"w1"}}`), nil
	case StageAgentStart:
		f.name = args[2]
		f.kind = args[4]
		if f.startFails {
			return nil, errors.New("herdr agent start: exit status 1")
		}
	}
	if args[0] == "agent" && args[1] == "get" {
		cwd := ""
		if len(f.calls) > 0 {
			cwd = f.cwd()
		}
		a := map[string]any{"workspace_id": "w1", "tab_id": "w1:t2", "pane_id": args[2], "terminal_id": "term-2", "foreground_cwd": cwd, "name": f.name, "agent": f.kind, "agent_status": "idle", "interactive_ready": true, "launch_pending": false}
		if !f.never {
			a["agent_session"] = map[string]any{"agent": f.kind, "kind": "id", "value": "sess-1"}
		}
		if f.working {
			a["agent_status"] = "working"
		}
		if f.agent != nil {
			f.agent(a)
		}
		raw, _ := json.Marshal(map[string]any{"agent": a})
		return raw, nil
	}
	return json.RawMessage(`{"ok":true}`), nil
}

func (f *fakeHerdr) cwd() string {
	for _, c := range f.calls {
		if c[0] == "tab" {
			for i := range c {
				if c[i] == "--cwd" {
					return c[i+1]
				}
			}
		}
	}
	return ""
}

func (f *fakeHerdr) count(prefix ...string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		if strings.HasPrefix(strings.Join(c, " "), strings.Join(prefix, " ")) {
			n++
		}
	}
	return n
}

func (f *fakeHerdr) total() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func testInput(t *testing.T) StartInput {
	dir := t.TempDir()
	return StartInput{
		Key:         "fixture-key",
		ReceiptPath: filepath.Join(dir, "nested", "start.json"),
		Provider:    ProviderCodex,
		Cwd:         filepath.Join(dir, "worktree"),
		Workspace:   "w1",
		Label:       "fixture",
		Prompt:      "do the bounded work",
		Herdr:       HerdrBinary{Path: "/usr/local/bin/herdr", SHA256: strings.Repeat("a", 64)},
	}
}

func str(s string) *string { return &s }

func run(in StartInput, f *fakeHerdr, fault func(string) error) (Receipt, error) {
	return start(context.Background(), in, ops{call: f.call, fault: fault, awaitTimeout: 300 * time.Millisecond, pause: time.Millisecond})
}

func resume(in StartInput, f *fakeHerdr) (Receipt, error) {
	return resumeAgent(context.Background(), in, f)
}

func resumeAgent(ctx context.Context, in StartInput, f *fakeHerdr) (Receipt, error) {
	return continueStart(ctx, in, ops{call: f.call, awaitTimeout: 300 * time.Millisecond, pause: time.Millisecond})
}

func wantCode(t *testing.T, err error, codes ...string) *Error {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("want *Error %v, got %v", codes, err)
	}
	for _, c := range codes {
		if e.Code == c {
			return e
		}
	}
	t.Fatalf("want code %v, got %q (%v)", codes, e.Code, e.Err)
	return e
}

func wantName(key string) string {
	sum := sha256.Sum256([]byte(key))
	return "ply-" + hex.EncodeToString(sum[:])[:28]
}

func TestAgentStartArgvPerProviderKeepsDefaults(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*StartInput)
		start  []string
	}{
		{"codex inherited defaults", func(in *StartInput) {}, []string{"agent", "start", wantName("fixture-key"), "--kind", "codex", "--pane", "w1:p2", "--timeout", "30000"}},
		{"claude inherited defaults", func(in *StartInput) { in.Provider = ProviderClaude }, []string{"agent", "start", wantName("fixture-key"), "--kind", "claude", "--pane", "w1:p2", "--timeout", "30000"}},
		{"claude model and permission mode", func(in *StartInput) {
			in.Provider, in.Model, in.Permission = ProviderClaude, str("opus"), str("auto")
		}, []string{"agent", "start", wantName("fixture-key"), "--kind", "claude", "--pane", "w1:p2", "--timeout", "30000", "--", "--model", "opus", "--permission-mode", "auto"}},
		{"codex model permission profile", func(in *StartInput) {
			in.Model, in.Permission, in.CodexProfile = str("gpt-x"), str("workspace"), str("fast")
		}, []string{"agent", "start", wantName("fixture-key"), "--kind", "codex", "--pane", "w1:p2", "--timeout", "30000", "--", "--model", "gpt-x", "-c", `default_permissions="workspace"`, "--profile", "fast"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := testInput(t)
			c.mutate(&in)
			f := newFake(t, in)
			r, err := run(in, f, nil)
			if err != nil || r.Reached != ReachedPromptSubmitted || r.Pending != "" || r.Failure != nil {
				t.Fatalf("%+v %v", r, err)
			}
			want := [][]string{
				{"tab", "create", "--workspace", "w1", "--cwd", in.Cwd, "--label", "fixture", "--no-focus"},
				c.start,
				{"agent", "get", "w1:p2"},
				{"agent", "prompt", "w1:p2", "do the bounded work", "--wait", "--until", "working", "--timeout", "10000"},
			}
			if fmt.Sprint(f.calls) != fmt.Sprint(want) {
				t.Fatalf("argv\n got %q\nwant %q", f.calls, want)
			}
			if len(wantName("fixture-key")) > 32 {
				t.Fatal("Herdr name too long")
			}
			if r.Pane != "w1:p2" || r.Tab != "w1:t2" || r.Terminal != "term-2" || r.Session != "sess-1" || r.Key != "fixture-key" || r.Provider != in.Provider || len(r.InputSHA256) != 64 {
				t.Fatalf("identity not recorded: %+v", r)
			}
		})
	}
}

func TestAgentStartPassesExplicitEnvironmentToTab(t *testing.T) {
	in := testInput(t)
	in.Env = []string{"A=1", "B=2"}
	f := newFake(t, in)
	if _, err := run(in, f, nil); err != nil {
		t.Fatal(err)
	}
	want := []string{"tab", "create", "--workspace", "w1", "--cwd", in.Cwd, "--label", "fixture", "--env", "A=1", "--env", "B=2", "--no-focus"}
	if fmt.Sprint(f.calls[0]) != fmt.Sprint(want) {
		t.Fatalf("%q", f.calls[0])
	}
}

func TestAgentStartRejectsInvalidInputBeforeAnyEffect(t *testing.T) {
	cases := map[string]func(*StartInput){
		"empty key":              func(in *StartInput) { in.Key = "" },
		"unknown provider":       func(in *StartInput) { in.Provider = "gemini" },
		"empty provider":         func(in *StartInput) { in.Provider = "" },
		"relative receipt":       func(in *StartInput) { in.ReceiptPath = "start.json" },
		"relative cwd":           func(in *StartInput) { in.Cwd = "worktree" },
		"empty workspace":        func(in *StartInput) { in.Workspace = "" },
		"empty prompt":           func(in *StartInput) { in.Prompt = "" },
		"short herdr sha":        func(in *StartInput) { in.Herdr.SHA256 = "abc" },
		"relative herdr path":    func(in *StartInput) { in.Herdr.Path = "herdr" },
		"profile on claude":      func(in *StartInput) { in.Provider, in.CodexProfile = ProviderClaude, str("p") },
		"env without assignment": func(in *StartInput) { in.Env = []string{"NOPE"} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			mutate(&in)
			f := newFake(t, in)
			_, err := run(in, f, nil)
			wantCode(t, err, CodeInvalidInput)
			if f.total() != 0 {
				t.Fatal("effect before validation")
			}
			if _, statErr := os.Stat(in.ReceiptPath); statErr == nil && filepath.IsAbs(in.ReceiptPath) {
				t.Fatal("reservation written for invalid input")
			}
		})
	}
}

func TestAgentStartSameInputTwiceSendsNothingAgain(t *testing.T) {
	in := testInput(t)
	pre := 0
	in.Precondition = func() error { pre++; return nil }
	f := newFake(t, in)
	first, err := run(in, f, nil)
	if err != nil {
		t.Fatal(err)
	}
	calls, checks := f.total(), pre
	before, _ := os.ReadFile(in.ReceiptPath)
	second, err := run(in, f, nil)
	if err != nil || !reflect.DeepEqual(second, first) {
		t.Fatalf("%+v %v", second, err)
	}
	if f.total() != calls || pre != checks {
		t.Fatal("repeat touched Herdr or the precondition")
	}
	after, _ := os.ReadFile(in.ReceiptPath)
	if string(before) != string(after) {
		t.Fatal("receipt rewritten")
	}
}

func TestAgentStartChangedInputMismatchesBeforeEffect(t *testing.T) {
	fields := map[string]func(*StartInput){
		"provider":      func(in *StartInput) { in.Provider = ProviderClaude },
		"model":         func(in *StartInput) { in.Model = str("m") },
		"permission":    func(in *StartInput) { in.Permission = str("p") },
		"codex profile": func(in *StartInput) { in.CodexProfile = str("p") },
		"cwd":           func(in *StartInput) { in.Cwd = in.Cwd + "-other" },
		"workspace":     func(in *StartInput) { in.Workspace = "w2" },
		"label":         func(in *StartInput) { in.Label = "other" },
		"env":           func(in *StartInput) { in.Env = []string{"A=1"} },
		"prompt":        func(in *StartInput) { in.Prompt = "changed" },
		"herdr sha":     func(in *StartInput) { in.Herdr.SHA256 = strings.Repeat("b", 64) },
		"herdr path":    func(in *StartInput) { in.Herdr.Path = "/opt/herdr" },
	}
	for name, mutate := range fields {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			if _, err := run(in, f, nil); err != nil {
				t.Fatal(err)
			}
			calls := f.total()
			before, _ := os.ReadFile(in.ReceiptPath)
			changed := in
			pre := 0
			changed.Precondition = func() error { pre++; return nil }
			mutate(&changed)
			got, err := run(changed, f, nil)
			wantCode(t, err, CodeInputMismatch)
			if got.InputSHA256 == "" || got.Reached != ReachedPromptSubmitted {
				t.Fatalf("stored receipt not returned: %+v", got)
			}
			after, _ := os.ReadFile(in.ReceiptPath)
			if f.total() != calls || pre != 0 || string(before) != string(after) {
				t.Fatal("mismatch had an effect")
			}
		})
	}
}

func TestAgentStartDigestIgnoresKeyReceiptPathAndPrecondition(t *testing.T) {
	a, b := testInput(t), testInput(t)
	b.Cwd, b.Key, b.Precondition = a.Cwd, "other-key", func() error { return nil }
	fa, fb := newFake(t, a), newFake(t, b)
	ra, err := run(a, fa, nil)
	if err != nil {
		t.Fatal(err)
	}
	rb, err := run(b, fb, nil)
	if err != nil {
		t.Fatal(err)
	}
	if ra.InputSHA256 != rb.InputSHA256 {
		t.Fatal("digest depends on identity, location or function")
	}
}

func TestAgentStartReceiptStoresDigestNotPromptAndWritesOnlyReceipt(t *testing.T) {
	in := testInput(t)
	in.Prompt = "SECRET-GOAL-TEXT"
	f := newFake(t, in)
	if _, err := run(in, f, nil); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(in.ReceiptPath)
	if strings.Contains(string(raw), "SECRET-GOAL-TEXT") {
		t.Fatal("prompt duplicated into receipt")
	}
	entries, _ := os.ReadDir(filepath.Dir(in.ReceiptPath))
	if len(entries) != 1 || entries[0].Name() != "start.json" {
		t.Fatalf("unexpected files: %v", entries)
	}
}

func TestAgentStartRecordsPendingBeforeEveryMutatingCallAndPreconditionBeforeEach(t *testing.T) {
	in := testInput(t)
	var log []string
	in.Precondition = func() error { log = append(log, "pre"); return nil }
	f := newFake(t, in)
	f.log = &log
	if _, err := run(in, f, nil); err != nil {
		t.Fatal(err)
	}
	want := "pre tab_create pre agent_start agent get pre prompt"
	if got := strings.Join(log, " "); got != want {
		t.Fatalf("order\n got %s\nwant %s", got, want)
	}
}

func TestAgentStartPromptsSettledCodexWithoutNativeSession(t *testing.T) {
	in := testInput(t)
	var log []string
	f := newFake(t, in)
	f.log, f.never = &log, true
	// Shape of the real Codex 0.164.0 / Herdr 0.9.5 agent get: no agent_session, no launch_pending.
	f.agent = func(a map[string]any) {
		a["agent_status"] = "done"
		delete(a, "launch_pending")
	}
	r, err := run(in, f, nil)
	if err != nil || r.Reached != ReachedPromptSubmitted || r.Pending != "" || r.Failure != nil || r.Session != "" {
		t.Fatalf("%+v %v", r, err)
	}
	if got, want := strings.Join(log, " "), "tab_create agent_start agent get prompt"; got != want {
		t.Fatalf("order\n got %s\nwant %s", got, want)
	}
	if f.count("agent", "prompt") != 1 || f.count("agent", "prompt", "w1:p2", "do the bounded work") != 1 {
		t.Fatalf("want exactly the task prompt: %q", f.calls)
	}
	raw, err := os.ReadFile(in.ReceiptPath)
	if err != nil || strings.Contains(string(raw), "session") {
		t.Fatalf("receipt must leave an unobserved session unknown: %s %v", raw, err)
	}
	before := f.total()
	if again, err := run(in, f, nil); err != nil || again.Reached != ReachedPromptSubmitted || f.total() != before {
		t.Fatalf("same key repeated an effect: %+v %v", again, err)
	}
}

func TestAgentStartNeverPromptsBeforeSessionlessAgentIsReady(t *testing.T) {
	for name, mutate := range map[string]func(map[string]any){
		"still working":             func(a map[string]any) { a["agent_status"] = "working" },
		"launch pending":            func(a map[string]any) { a["launch_pending"] = true },
		"not interactive":           func(a map[string]any) { a["interactive_ready"] = false },
		"interactive unknown":       func(a map[string]any) { delete(a, "interactive_ready") },
		"foreground cwd unreported": func(a map[string]any) { delete(a, "foreground_cwd") },
	} {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			f.never, f.agent = true, mutate
			r, err := run(in, f, nil)
			wantCode(t, err, ClassNotReady)
			if r.Reached != ReachedAgentStarted || r.Failure == nil || r.Failure.Stage != StageAwaitReady || r.Failure.Class != ClassNotReady || r.Failure.EffectUnknown || r.Pending != "" {
				t.Fatalf("%+v", r)
			}
			if f.count("agent", "prompt") != 0 {
				t.Fatalf("prompt sent before ready: %q", f.calls)
			}
			before := f.total()
			if _, err = run(in, f, nil); err == nil || f.total() != before {
				t.Fatal("retry after failure repeated an effect")
			}
		})
	}
}

func TestAgentStartRefusesPromptWhenObservedNativeSessionDoesNotHold(t *testing.T) {
	for name, later := range map[string]map[string]any{
		"changed":  {"agent": "codex", "kind": "id", "value": "sess-2"},
		"vanished": nil,
	} {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			polls := 0
			f.agent = func(a map[string]any) {
				if polls++; polls == 1 {
					a["agent_status"] = "working"
					return
				}
				delete(a, "agent_session")
				if later != nil {
					a["agent_session"] = later
				}
			}
			r, err := run(in, f, nil)
			wantCode(t, err, ClassAgentResponse)
			if f.count("agent", "prompt") != 0 || r.Reached != ReachedAgentStarted || r.Failure == nil {
				t.Fatalf("%+v calls=%q", r, f.calls)
			}
		})
	}
}

func TestAgentStartRefusesPromptWhenObservedIdentityDiffers(t *testing.T) {
	cases := map[string]func(map[string]any){
		"name":         func(a map[string]any) { a["name"] = "someone-else" },
		"agent kind":   func(a map[string]any) { a["agent"] = "codex-other" },
		"pane":         func(a map[string]any) { a["pane_id"] = "w1:p9" },
		"terminal":     func(a map[string]any) { a["terminal_id"] = "term-9" },
		"tab":          func(a map[string]any) { a["tab_id"] = "w1:t9" },
		"workspace":    func(a map[string]any) { a["workspace_id"] = "w9" },
		"no workspace": func(a map[string]any) { delete(a, "workspace_id") },
		"no tab":       func(a map[string]any) { delete(a, "tab_id") },
		"no pane":      func(a map[string]any) { delete(a, "pane_id") },
		"no terminal":  func(a map[string]any) { delete(a, "terminal_id") },
		"empty tab":    func(a map[string]any) { a["tab_id"] = "" },
		"session kind": func(a map[string]any) {
			a["agent_session"] = map[string]any{"agent": "codex", "kind": "path", "value": "x"}
		},
		"session agent": func(a map[string]any) {
			a["agent_session"] = map[string]any{"agent": "claude", "kind": "id", "value": "x"}
		},
	}
	for name, mutate := range cases {
		for _, sessionless := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s sessionless=%v", name, sessionless), func(t *testing.T) {
				in := testInput(t)
				f := newFake(t, in)
				f.agent, f.never = mutate, sessionless
				r, err := run(in, f, nil)
				wantCode(t, err, ClassAgentResponse)
				if f.count("agent", "prompt") != 0 || r.Reached != ReachedAgentStarted || r.Failure == nil {
					t.Fatalf("%+v", r)
				}
			})
		}
	}
	t.Run("claude foreground cwd", func(t *testing.T) {
		in := testInput(t)
		in.Provider = ProviderClaude
		f := newFake(t, in)
		f.agent = func(a map[string]any) { a["foreground_cwd"] = "/elsewhere" }
		_, err := run(in, f, nil)
		wantCode(t, err, ClassAgentResponse)
		if f.count("agent", "prompt") != 0 {
			t.Fatal("prompt sent from wrong cwd")
		}
	})
}

func TestAgentStartRejectsIncompleteTabIdentityBeforeAgentStart(t *testing.T) {
	for name, tab := range map[string]string{
		"only pane":    `{"root_pane":{"pane_id":"w1:p2"}}`,
		"no workspace": `{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","terminal_id":"term-2"}}`,
		"no tab":       `{"root_pane":{"pane_id":"w1:p2","terminal_id":"term-2","workspace_id":"w1"}}`,
		"no terminal":  `{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","workspace_id":"w1"}}`,
		"no pane":      `{"root_pane":{"tab_id":"w1:t2","terminal_id":"term-2","workspace_id":"w1"}}`,
		"other space":  `{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","terminal_id":"term-2","workspace_id":"w9"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			f.tab = tab
			r, err := run(in, f, nil)
			wantCode(t, err, ClassTabResponse)
			if f.count("agent") != 0 || r.Reached != ReachedReserved || r.Pending != StageTabCreate || r.Failure == nil || !r.Failure.EffectUnknown || !r.EffectUnknown() {
				t.Fatalf("%+v calls=%q", r, f.calls)
			}
		})
	}
}

func TestAgentStartChecksCodexForegroundCwdOnlyWhenReported(t *testing.T) {
	in := testInput(t)
	f := newFake(t, in)
	f.agent = func(a map[string]any) { a["foreground_cwd"] = "/somewhere/else" }
	r, err := run(in, f, nil)
	wantCode(t, err, ClassAgentResponse)
	if f.count("agent", "prompt") != 0 || r.Reached != ReachedAgentStarted {
		t.Fatalf("%+v", r)
	}
	in = testInput(t)
	f = newFake(t, in)
	f.never, f.agent = true, func(a map[string]any) { a["foreground_cwd"] = "/somewhere/else" }
	r, err = run(in, f, nil)
	wantCode(t, err, ClassAgentResponse)
	if f.count("agent", "prompt") != 0 || r.Reached != ReachedAgentStarted {
		t.Fatalf("sessionless wrong cwd: %+v", r)
	}
	in = testInput(t)
	f = newFake(t, in)
	f.agent = func(a map[string]any) { delete(a, "foreground_cwd") }
	if r, err = run(in, f, nil); err != nil || r.Reached != ReachedPromptSubmitted {
		t.Fatalf("missing codex cwd must stay accepted: %+v %v", r, err)
	}
}

func TestAgentStartLostReplyKeepsPendingAndNeverResends(t *testing.T) {
	cases := []struct {
		stage, reached, class string
		sessionless           bool
	}{
		{StageTabCreate, ReachedReserved, ClassTabResponse, false},
		{StageAgentStart, ReachedTabCreated, ClassAgentResponse, false},
		{StagePrompt, ReachedSessionBound, ClassPromptResponse, false},
		{StagePrompt, ReachedAgentStarted, ClassPromptResponse, true},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s sessionless=%v", c.stage, c.sessionless), func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			f.drop[c.stage], f.never = true, c.sessionless
			r, err := run(in, f, nil)
			wantCode(t, err, c.class)
			if r.Pending != c.stage || r.Reached != c.reached || r.Failure == nil || !r.Failure.EffectUnknown || r.Failure.Stage != c.stage || r.Failure.Class != c.class || !r.EffectUnknown() {
				t.Fatalf("%+v", r)
			}
			saved, err := ReadReceipt(in.ReceiptPath)
			if err != nil || !reflect.DeepEqual(saved, r) {
				t.Fatalf("not durable: %+v %v", saved, err)
			}
			before := f.total()
			again, err := run(in, f, nil)
			wantCode(t, err, CodeAlreadyReserved)
			if f.total() != before || again.Pending != c.stage || !again.EffectUnknown() {
				t.Fatal("retry after lost reply sent again or hid unknown effect")
			}
		})
	}
}

func TestAgentStartPreconditionFailureAfterPartialEffectIsKnownPerStep(t *testing.T) {
	cases := []struct {
		failAt  int
		stage   string
		reached string
		calls   int
	}{
		{1, StageTabCreate, ReachedReserved, 0},
		{2, StageAgentStart, ReachedTabCreated, 1},
		{3, StagePrompt, ReachedSessionBound, 3},
	}
	for _, c := range cases {
		t.Run(c.stage, func(t *testing.T) {
			in := testInput(t)
			n := 0
			in.Precondition = func() error {
				n++
				if n == c.failAt {
					return errors.New("task moved")
				}
				return nil
			}
			f := newFake(t, in)
			r, err := run(in, f, nil)
			e := wantCode(t, err, ClassPrecondition)
			if e.Stage != c.stage || r.Reached != c.reached || r.Pending != "" || r.Failure == nil || r.Failure.Class != ClassPrecondition || r.Failure.Stage != c.stage || r.Failure.EffectUnknown {
				t.Fatalf("%+v %+v", r, e)
			}
			if f.total() != c.calls {
				t.Fatalf("calls %d, want %d", f.total(), c.calls)
			}
			before := f.total()
			if _, err = run(in, f, nil); err == nil || f.total() != before {
				t.Fatal("retry after precondition failure repeated an effect")
			}
		})
	}
}

func TestAgentStartCrashPointsNeverReadAsFreeOrLegacyAndNeverRepeat(t *testing.T) {
	points := []struct {
		point       string
		sessionless bool
		empty       bool
	}{
		{"after_excl", false, true},
		{"after_reserve", false, false},
		{"after_pending:" + StageTabCreate, false, false},
		{"after_call:" + StageTabCreate, false, false},
		{"after_pending:" + StageAgentStart, false, false},
		{"after_call:" + StageAgentStart, false, false},
		{"after_pending:" + StagePrompt, false, false},
		{"after_call:" + StagePrompt, false, false},
		{"after_pending:" + StagePrompt, true, false},
		{"after_call:" + StagePrompt, true, false},
	}
	for _, p := range points {
		t.Run(fmt.Sprintf("%s sessionless=%v", p.point, p.sessionless), func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			f.never = p.sessionless
			_, err := run(in, f, func(point string) error {
				if point == p.point {
					return errCrash
				}
				return nil
			})
			if !errors.Is(err, errCrash) {
				t.Fatalf("crash point %s not reached: %v", p.point, err)
			}
			before := f.total()
			read, readErr := ReadReceipt(in.ReceiptPath)
			if p.empty {
				wantCode(t, readErr, CodeIncompleteReservation)
			} else if readErr != nil || read.Legacy || !read.EffectUnknown() {
				t.Fatalf("unknown effect read as safe: %+v %v", read, readErr)
			}
			again, err := run(in, f, nil)
			wantCode(t, err, CodeIncompleteReservation, CodeAlreadyReserved)
			if again.Legacy || f.total() != before {
				t.Fatal("restart repeated an effect or read crash as legacy")
			}
		})
	}
}

func TestAgentStartConcurrentOwnersStartExactlyOnce(t *testing.T) {
	in := testInput(t)
	f := newFake(t, in)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var okReceipts []Receipt
	var errs []error
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := run(in, f, nil)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errs = append(errs, err)
			} else {
				okReceipts = append(okReceipts, r)
			}
		}()
	}
	wg.Wait()
	if f.count("tab", "create") != 1 || f.count("agent", "start") != 1 || f.count("agent", "prompt") != 1 {
		t.Fatalf("duplicate effects: %q", f.calls)
	}
	if len(okReceipts) == 0 {
		t.Fatal("no owner succeeded")
	}
	for _, r := range okReceipts {
		if r.Reached != ReachedPromptSubmitted {
			t.Fatalf("success reported for unfinished start: %+v", r)
		}
	}
	for _, err := range errs {
		wantCode(t, err, CodeAlreadyReserved, CodeIncompleteReservation)
	}
}

func TestAgentStartConcurrentDifferentInputNeverStartsSecond(t *testing.T) {
	in := testInput(t)
	f := newFake(t, in)
	var wg sync.WaitGroup
	var mu sync.Mutex
	codes := map[string]int{}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mine := in
			if i%2 == 1 {
				mine.Prompt = "different prompt"
			}
			_, err := run(mine, f, nil)
			mu.Lock()
			defer mu.Unlock()
			var e *Error
			if errors.As(err, &e) {
				codes[e.Code]++
			} else {
				codes["ok"]++
			}
		}(i)
	}
	wg.Wait()
	if f.count("tab", "create") != 1 || f.count("agent", "prompt") != 1 {
		t.Fatalf("second owner started: %q", f.calls)
	}
	for code := range codes {
		switch code {
		case "ok", CodeInputMismatch, CodeAlreadyReserved, CodeIncompleteReservation:
		default:
			t.Fatalf("unexpected outcome %q", code)
		}
	}
}

const legacySubmitted = `{"plan":{"task":"t","worktree":"/w","ref":"refs/heads/t","oid":"abc","herdr_workspace":"wD","goal":{"constraints":["x"],"executor":{"effort":null,"model":null,"provider":"codex"}}},"state":"prompt_submitted","pane":"wD:p3"}`
const legacyTab = `{"plan":{"task":"t","worktree":"/w","ref":"refs/heads/t","oid":"abc","herdr_workspace":"wD","goal":{}},"state":"tab_created","pane":"wD:p4","error":"herdr agent: exit status 1"}`

func TestAgentStartReadsLegacyReceiptsWithoutClaimingInputBinding(t *testing.T) {
	cases := []struct {
		name, raw, reached, pending string
		unknown                     bool
	}{
		{"submitted", legacySubmitted, ReachedPromptSubmitted, "", false},
		{"tab created after failure", legacyTab, ReachedTabCreated, "", true},
		{"reserved", `{"plan":{"task":"t"},"state":"reserved"}`, ReachedReserved, "", true},
		{"agent ready", `{"plan":{"task":"t"},"state":"agent_ready","pane":"p"}`, ReachedAgentStarted, "", true},
		{"prompt attempting", `{"plan":{"task":"t"},"state":"prompt_attempting","pane":"p"}`, ReachedAgentStarted, StagePrompt, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := testInput(t)
			if err := os.MkdirAll(filepath.Dir(in.ReceiptPath), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(in.ReceiptPath, []byte(c.raw), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := ReadReceipt(in.ReceiptPath)
			if err != nil || !got.Legacy || got.Reached != c.reached || got.Pending != c.pending || got.InputSHA256 != "" || got.EffectUnknown() != c.unknown {
				t.Fatalf("%+v %v", got, err)
			}
			f := newFake(t, in)
			pre := 0
			in.Precondition = func() error { pre++; return nil }
			again, err := run(in, f, nil)
			wantCode(t, err, CodeLegacyUnbound)
			if !reflect.DeepEqual(again, got) || f.total() != 0 || pre != 0 {
				t.Fatal("legacy receipt re-run or treated as bound")
			}
			after, _ := os.ReadFile(in.ReceiptPath)
			if string(after) != c.raw {
				t.Fatal("legacy receipt rewritten")
			}
		})
	}
}

func TestAgentStartReceiptFromBootstrapCandidateStaysUnknownAndRefused(t *testing.T) {
	in := testInput(t)
	digest, err := inputDigest(in)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.MkdirAll(filepath.Dir(in.ReceiptPath), 0700); err != nil {
		t.Fatal(err)
	}
	raw := fmt.Sprintf(`{"key":"fixture-key","provider":"codex","input_sha256":%q,"reached":"agent_started","pane":"w1:p2","failure":{"stage":"await_session","class":"session_unbound","effect_unknown":true},"bootstrap":"acknowledged"}`, digest)
	if err = os.WriteFile(in.ReceiptPath, []byte(raw), 0600); err != nil {
		t.Fatal(err)
	}
	f := newFake(t, in)
	r, err := run(in, f, nil)
	wantCode(t, err, CodeAlreadyReserved)
	if f.total() != 0 || r.Legacy || !r.EffectUnknown() || r.Failure.Class != "session_unbound" {
		t.Fatalf("%+v calls=%q", r, f.calls)
	}
}

func TestAgentStartUnparseableReceiptIsIncompleteNotFreeOrLegacy(t *testing.T) {
	for name, raw := range map[string]string{"empty": "", "garbage": "{not json", "foreign json": `{"other":1}`, "truncated legacy": `{"plan":{"task":"t"`} {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			if err := os.MkdirAll(filepath.Dir(in.ReceiptPath), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(in.ReceiptPath, []byte(raw), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := ReadReceipt(in.ReceiptPath)
			wantCode(t, err, CodeIncompleteReservation)
			f := newFake(t, in)
			_, err = run(in, f, nil)
			wantCode(t, err, CodeIncompleteReservation)
			after, _ := os.ReadFile(in.ReceiptPath)
			if f.total() != 0 || string(after) != raw {
				t.Fatal("effect or rewrite on unreadable reservation")
			}
		})
	}
}

func TestAgentStartEffectUnknownClassification(t *testing.T) {
	for name, c := range map[string]struct {
		r    Receipt
		want bool
	}{
		"complete":             {Receipt{Reached: ReachedPromptSubmitted}, false},
		"pending":              {Receipt{Reached: ReachedTabCreated, Pending: StageAgentStart}, true},
		"reserved only":        {Receipt{Reached: ReachedReserved}, true},
		"session bound only":   {Receipt{Reached: ReachedSessionBound}, true},
		"known precondition":   {Receipt{Reached: ReachedTabCreated, Failure: &StartFailure{Class: ClassPrecondition}}, false},
		"unknown failure":      {Receipt{Reached: ReachedTabCreated, Failure: &StartFailure{EffectUnknown: true}}, true},
		"complete with marker": {Receipt{Reached: ReachedPromptSubmitted, Pending: StagePrompt}, true},
	} {
		if c.r.EffectUnknown() != c.want {
			t.Errorf("%s: got %v", name, !c.want)
		}
	}
}

func TestAgentStartPackageHasNoProcessDependencies(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}", ".").Output()
	if err != nil {
		t.Fatal(err)
	}
	const module = "github.com/devdimensionlab/plybuild"
	allowed := map[string]bool{module + "/internal/agentstart": true, module + "/internal/canonicaljson": true, module + "/internal/atomicfile": true, module + "/internal/filelock": true, "golang.org/x/sys/unix": true}
	for _, path := range strings.Fields(string(out)) {
		first := strings.SplitN(path, "/", 2)[0]
		if !strings.Contains(first, ".") || allowed[path] {
			continue
		}
		t.Errorf("unexpected dependency %s", path)
	}
}

func TestAgentStartExportedSurfaceNamesNoProcessConcepts(t *testing.T) {
	files, _ := filepath.Glob("*.go")
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, _ := os.ReadFile(name)
		for _, word := range []string{"Task", "Goal", "workspace.", "handoff", "QA", "review"} {
			if strings.Contains(string(raw), word) {
				t.Errorf("%s mentions %q", name, word)
			}
		}
	}
}

func executable(t *testing.T, body string) HerdrBinary {
	t.Helper()
	path := filepath.Join(t.TempDir(), "herdr")
	if err := os.WriteFile(path, []byte(body), 0700); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(body))
	return HerdrBinary{Path: path, SHA256: hex.EncodeToString(sum[:])}
}

func TestAgentStartHerdrClientVerifiesBinaryBeforeExecuting(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "ran")
	h := executable(t, "#!/bin/sh\necho ran > "+marker+"\necho '{\"result\":{\"ok\":true}}'\n")
	got, err := callHerdr(context.Background(), h, []string{"agent", "get", "p"})
	if err != nil || string(got) != `{"ok":true}` {
		t.Fatalf("%s %v", got, err)
	}
	_ = os.Remove(marker)
	h.SHA256 = strings.Repeat("0", 64)
	_, err = callHerdr(context.Background(), h, nil)
	wantCode(t, err, ClassBinaryChanged)
	if _, statErr := os.Stat(marker); statErr == nil {
		t.Fatal("changed binary was executed")
	}
}

func TestAgentStartHerdrClientRejectsUnsuccessfulOrOversizedOutput(t *testing.T) {
	for name, body := range map[string]string{
		"error envelope": "#!/bin/sh\necho '{\"error\":{\"message\":\"no\"}}'\n",
		"null result":    "#!/bin/sh\necho '{\"result\":null}'\n",
		"not json":       "#!/bin/sh\necho nope\n",
		"exit status":    "#!/bin/sh\necho '{\"result\":{}}'\nexit 3\n",
		"oversized":      "#!/bin/sh\nhead -c 2000000 /dev/zero | tr '\\0' 'a'\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := callHerdr(context.Background(), executable(t, body), []string{"x"}); err == nil {
				t.Fatal("unsuccessful Herdr reply accepted")
			}
		})
	}
}

func TestAgentStartResolveHerdrBindsResolvedBinary(t *testing.T) {
	dir := t.TempDir()
	body := "#!/bin/sh\nexit 0\n"
	real := filepath.Join(dir, "real-herdr")
	if err := os.WriteFile(real, []byte(body), 0700); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(bin, "herdr")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	got, err := ResolveHerdr()
	wantPath, _ := filepath.EvalSymlinks(real)
	sum := sha256.Sum256([]byte(body))
	if err != nil || got.Path != wantPath || got.SHA256 != hex.EncodeToString(sum[:]) {
		t.Fatalf("%+v %v", got, err)
	}
	t.Setenv("PATH", t.TempDir())
	if _, err = ResolveHerdr(); err == nil {
		t.Fatal("missing Herdr accepted")
	}
}

// stoppedAtTrust leaves a start whose agent start reply was lost while the
// pane already carries the reserved agent, as observed when Claude waits on
// folder trust.
func stoppedAtTrust(t *testing.T, in StartInput, f *fakeHerdr) []byte {
	t.Helper()
	f.startFails = true
	r, err := run(in, f, nil)
	wantCode(t, err, ClassAgentResponse)
	if r.Reached != ReachedTabCreated || r.Pending != StageAgentStart || r.Failure == nil || r.Failure.Stage != StageAgentStart || !r.Failure.EffectUnknown {
		t.Fatalf("not stopped at agent start: %+v", r)
	}
	f.startFails = false
	raw, err := os.ReadFile(in.ReceiptPath)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func wantResumed(t *testing.T, in StartInput, f *fakeHerdr, before int) Receipt {
	t.Helper()
	r, err := resume(in, f)
	if err != nil || r.Reached != ReachedPromptSubmitted || r.Pending != "" || r.Failure != nil {
		t.Fatalf("not resumed: %+v %v", r, err)
	}
	saved, err := ReadReceipt(in.ReceiptPath)
	if err != nil || !reflect.DeepEqual(saved, r) {
		t.Fatalf("resumed receipt not durable: %+v %v", saved, err)
	}
	f.mu.Lock()
	after := f.calls[before:]
	f.mu.Unlock()
	if len(after) < 2 {
		t.Fatalf("calls after stop: %q", after)
	}
	for _, c := range after[:len(after)-1] {
		if fmt.Sprint(c) != fmt.Sprint([]string{"agent", "get", r.Pane}) {
			t.Fatalf("resume made a call other than agent get before the prompt: %q", after)
		}
	}
	if fmt.Sprint(after[len(after)-1]) != fmt.Sprint(promptArgs(r.Pane, in.Prompt)) || f.count("agent", "prompt") != 1 {
		t.Fatalf("want exactly one prompt to the reserved pane: %q", after)
	}
	return r
}

func wantUnchanged(t *testing.T, in StartInput, raw []byte) {
	t.Helper()
	after, err := os.ReadFile(in.ReceiptPath)
	if err != nil || string(after) != string(raw) {
		t.Fatalf("receipt bytes changed:\n%s\n%s %v", raw, after, err)
	}
}

func TestAgentResumeContinuesStartStoppedBeforePromptInTheSameAgent(t *testing.T) {
	for _, provider := range []string{ProviderClaude, ProviderCodex} {
		t.Run(provider+" trust", func(t *testing.T) {
			in := testInput(t)
			in.Provider = provider
			f := newFake(t, in)
			stoppedAtTrust(t, in, f)
			before := f.total()
			again, err := run(in, f, nil)
			wantCode(t, err, CodeAlreadyReserved)
			if f.total() != before || again.Pending != StageAgentStart {
				t.Fatal("plain start continued a stopped start")
			}
			r := wantResumed(t, in, f, before)
			if r.Session != "sess-1" || r.Tab != "w1:t2" || r.Pane != "w1:p2" || r.Terminal != "term-2" {
				t.Fatalf("identity not kept: %+v", r)
			}
			if f.count("tab", "create") != 1 || f.count("agent", "start") != 1 {
				t.Fatalf("resume repeated a start effect: %q", f.calls)
			}
		})
		t.Run(provider+" await ready deadline", func(t *testing.T) {
			in := testInput(t)
			in.Provider = provider
			f := newFake(t, in)
			f.working = true
			r, err := run(in, f, nil)
			wantCode(t, err, ClassNotReady)
			if r.Reached != ReachedAgentStarted || r.Failure == nil || r.Failure.Stage != StageAwaitReady {
				t.Fatalf("%+v", r)
			}
			f.working = false
			wantResumed(t, in, f, f.total())
		})
	}
	t.Run("sessionless codex", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		f.never = true
		stoppedAtTrust(t, in, f)
		f.agent = func(a map[string]any) { a["agent_status"] = "done"; delete(a, "launch_pending") }
		r := wantResumed(t, in, f, f.total())
		if r.Session != "" || r.Reached != ReachedPromptSubmitted {
			t.Fatalf("unobserved session invented: %+v", r)
		}
	})
}

func TestAgentResumeRefusesDifferentObservedAgentKeepingReceiptForLaterResume(t *testing.T) {
	cases := map[string]func(map[string]any){
		"pane":            func(a map[string]any) { a["pane_id"] = "w1:p9" },
		"terminal":        func(a map[string]any) { a["terminal_id"] = "term-9" },
		"tab":             func(a map[string]any) { a["tab_id"] = "w1:t9" },
		"workspace":       func(a map[string]any) { a["workspace_id"] = "w9" },
		"name":            func(a map[string]any) { a["name"] = "someone-else" },
		"agent not named": func(a map[string]any) { a["name"] = ""; a["agent"] = "" },
		"provider":        func(a map[string]any) { a["agent"] = ProviderCodex },
		"cwd":             func(a map[string]any) { a["foreground_cwd"] = "/elsewhere" },
		"session agent": func(a map[string]any) {
			a["agent_session"] = map[string]any{"agent": ProviderCodex, "kind": "id", "value": "sess-1"}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			in.Provider = ProviderClaude
			f := newFake(t, in)
			raw := stoppedAtTrust(t, in, f)
			f.agent = mutate
			got, err := resume(in, f)
			e := wantCode(t, err, ClassAgentResponse)
			if e.Stage != StageAwaitReady || f.count("agent", "prompt") != 0 || f.count("agent", "start") != 1 || f.count("tab", "create") != 1 {
				t.Fatalf("%+v calls=%q", e, f.calls)
			}
			if stored, _ := ReadReceipt(in.ReceiptPath); !reflect.DeepEqual(got, stored) {
				t.Fatalf("stored receipt not returned: %+v", got)
			}
			wantUnchanged(t, in, raw)
			f.agent = nil
			wantResumed(t, in, f, f.total())
		})
	}
	t.Run("agent start never named the pane", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		f.drop[StageAgentStart] = true
		if _, err := run(in, f, nil); err == nil {
			t.Fatal("lost agent start reply not reported")
		}
		raw, _ := os.ReadFile(in.ReceiptPath)
		_, err := resume(in, f)
		wantCode(t, err, ClassAgentResponse)
		if f.count("agent", "prompt") != 0 {
			t.Fatal("prompt sent to an unnamed pane")
		}
		wantUnchanged(t, in, raw)
		f.name, f.kind = wantName(in.Key), in.Provider
		wantResumed(t, in, f, f.total())
	})
	t.Run("different native session than stored", func(t *testing.T) {
		in := testInput(t)
		n := 0
		in.Precondition = func() error {
			if n++; n == 3 {
				return errors.New("task moved")
			}
			return nil
		}
		f := newFake(t, in)
		r, err := run(in, f, nil)
		wantCode(t, err, ClassPrecondition)
		if r.Reached != ReachedSessionBound || r.Session != "sess-1" {
			t.Fatalf("%+v", r)
		}
		raw, _ := os.ReadFile(in.ReceiptPath)
		in.Precondition = nil
		f.agent = func(a map[string]any) {
			a["agent_session"] = map[string]any{"agent": ProviderCodex, "kind": "id", "value": "sess-2"}
		}
		_, err = resume(in, f)
		wantCode(t, err, ClassAgentResponse)
		if f.count("agent", "prompt") != 0 {
			t.Fatal("prompt sent to a replaced session")
		}
		wantUnchanged(t, in, raw)
		f.agent = nil
		if r = wantResumed(t, in, f, f.total()); r.Session != "sess-1" {
			t.Fatalf("%+v", r)
		}
	})
}

func TestAgentResumeNotReadyOrPreconditionKeepsReceiptForLaterResume(t *testing.T) {
	t.Run("deadline", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		raw := stoppedAtTrust(t, in, f)
		f.working = true
		_, err := resume(in, f)
		if e := wantCode(t, err, ClassNotReady); e.Stage != StageAwaitReady || f.count("agent", "prompt") != 0 {
			t.Fatalf("%+v", e)
		}
		wantUnchanged(t, in, raw)
		f.working = false
		wantResumed(t, in, f, f.total())
	})
	t.Run("cancelled", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		raw := stoppedAtTrust(t, in, f)
		f.working = true
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := resumeAgent(ctx, in, f)
		wantCode(t, err, ClassNotReady)
		wantUnchanged(t, in, raw)
		f.working = false
		wantResumed(t, in, f, f.total())
	})
	t.Run("precondition", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		raw := stoppedAtTrust(t, in, f)
		in.Precondition = func() error { return errors.New("task moved") }
		_, err := resume(in, f)
		if e := wantCode(t, err, ClassPrecondition); e.Stage != StagePrompt || f.count("agent", "prompt") != 0 {
			t.Fatalf("%+v", e)
		}
		wantUnchanged(t, in, raw)
		in.Precondition = func() error { return nil }
		wantResumed(t, in, f, f.total())
	})
}

func TestAgentResumeRefusesChangedOrLegacyInputWithoutEffect(t *testing.T) {
	for name, mutate := range map[string]func(*StartInput){
		"prompt":    func(in *StartInput) { in.Prompt = "changed" },
		"herdr sha": func(in *StartInput) { in.Herdr.SHA256 = strings.Repeat("b", 64) },
		"key":       func(in *StartInput) { in.Key = "other-key" },
	} {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			raw := stoppedAtTrust(t, in, f)
			changed, pre := in, 0
			changed.Precondition = func() error { pre++; return nil }
			mutate(&changed)
			before := f.total()
			got, err := resume(changed, f)
			wantCode(t, err, CodeInputMismatch)
			if f.total() != before || pre != 0 || got.Pending != StageAgentStart {
				t.Fatalf("mismatch had an effect or hid the stored receipt: %+v", got)
			}
			wantUnchanged(t, in, raw)
		})
	}
	t.Run("legacy", func(t *testing.T) {
		in := testInput(t)
		if err := os.MkdirAll(filepath.Dir(in.ReceiptPath), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(in.ReceiptPath, []byte(legacyTab), 0600); err != nil {
			t.Fatal(err)
		}
		f := newFake(t, in)
		got, err := resume(in, f)
		wantCode(t, err, CodeLegacyUnbound)
		if f.total() != 0 || !got.Legacy {
			t.Fatalf("%+v", got)
		}
		wantUnchanged(t, in, []byte(legacyTab))
	})
}

func TestAgentResumeNeverResendsUnknownPromptOrContinuesUnfinishedOwner(t *testing.T) {
	stops := map[string]struct {
		drop  string
		fault string
		pre   int
	}{
		"prompt reply lost":         {drop: StagePrompt},
		"tab create reply lost":     {drop: StageTabCreate},
		"precondition before tab":   {pre: 1},
		"reserved crash":            {fault: "after_reserve"},
		"tab create crash":          {fault: "after_call:" + StageTabCreate},
		"agent start crash":         {fault: "after_call:" + StageAgentStart},
		"pending agent start crash": {fault: "after_pending:" + StageAgentStart},
		"pending prompt crash":      {fault: "after_pending:" + StagePrompt},
		"prompt crash after call":   {fault: "after_call:" + StagePrompt},
	}
	for name, stop := range stops {
		t.Run(name, func(t *testing.T) {
			in := testInput(t)
			n := 0
			in.Precondition = func() error {
				if n++; n == stop.pre {
					return errors.New("task moved")
				}
				return nil
			}
			f := newFake(t, in)
			if stop.drop != "" {
				f.drop[stop.drop] = true
			}
			_, err := run(in, f, func(point string) error {
				if point == stop.fault {
					return errCrash
				}
				return nil
			})
			if err == nil {
				t.Fatal("start did not stop")
			}
			delete(f.drop, stop.drop)
			raw, _ := os.ReadFile(in.ReceiptPath)
			before := f.total()
			got, err := resume(in, f)
			wantCode(t, err, CodeNotResumable)
			if f.total() != before || got.InputSHA256 == "" {
				t.Fatalf("not resumable start touched Herdr or hid the receipt: %+v calls=%q", got, f.calls[before:])
			}
			wantUnchanged(t, in, raw)
		})
	}
	t.Run("prompt reply lost during resume", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		stoppedAtTrust(t, in, f)
		f.drop[StagePrompt] = true
		r, err := resume(in, f)
		wantCode(t, err, ClassPromptResponse)
		if r.Pending != StagePrompt || r.Failure == nil || !r.Failure.EffectUnknown || r.Reached != ReachedSessionBound {
			t.Fatalf("%+v", r)
		}
		if saved, _ := ReadReceipt(in.ReceiptPath); !reflect.DeepEqual(saved, r) {
			t.Fatalf("lost prompt not durable: %+v", saved)
		}
		delete(f.drop, StagePrompt)
		before := f.total()
		_, err = resume(in, f)
		wantCode(t, err, CodeNotResumable)
		if f.total() != before || f.count("agent", "prompt") != 1 {
			t.Fatal("unknown prompt sent again")
		}
	})
}

func TestAgentResumeReadsBackCompleteStartWithoutEffect(t *testing.T) {
	in := testInput(t)
	pre := 0
	in.Precondition = func() error { pre++; return nil }
	f := newFake(t, in)
	first, err := run(in, f, nil)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(in.ReceiptPath)
	calls, checks := f.total(), pre
	got, err := resume(in, f)
	if err != nil || !reflect.DeepEqual(got, first) || f.total() != calls || pre != checks {
		t.Fatalf("readback had an effect: %+v %v", got, err)
	}
	wantUnchanged(t, in, raw)
}

func TestAgentResumeConcurrentCallsPromptTheSameAgentOnce(t *testing.T) {
	in := testInput(t)
	f := newFake(t, in)
	stoppedAtTrust(t, in, f)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var receipts []Receipt
	var errs []error
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := resume(in, f)
			mu.Lock()
			defer mu.Unlock()
			receipts, errs = append(receipts, r), append(errs, err)
		}()
	}
	wg.Wait()
	if f.count("agent", "prompt") != 1 || f.count("agent", "start") != 1 || f.count("tab", "create") != 1 {
		t.Fatalf("duplicate effects: %q", f.calls)
	}
	for i, err := range errs {
		if err != nil || receipts[i].Reached != ReachedPromptSubmitted || receipts[i].Failure != nil {
			t.Fatalf("resume %d: %+v %v", i, receipts[i], err)
		}
	}
}

func resumeWithLockFault(in StartInput, f *fakeHerdr, op string) (Receipt, error) {
	fault := errors.New("injected " + op + " fault")
	return continueStart(context.Background(), in, ops{call: f.call, awaitTimeout: 300 * time.Millisecond, pause: time.Millisecond, lockFail: func(got string) error {
		if got == op {
			return fault
		}
		return nil
	}})
}

func wantReleaseError(t *testing.T, err error) {
	t.Helper()
	var released *filelock.ReleaseError
	if !errors.As(err, &released) {
		t.Fatalf("lock release failure hidden: %v", err)
	}
}

func TestAgentResumeReportsLockReleaseFailureWithoutHidingOutcome(t *testing.T) {
	for _, op := range []string{"unlock", "lock-close"} {
		t.Run(op+" after prompt", func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			stoppedAtTrust(t, in, f)
			r, err := resumeWithLockFault(in, f, op)
			wantReleaseError(t, err)
			if r.Reached != ReachedPromptSubmitted || r.Pending != "" || r.Failure != nil || f.count("agent", "prompt") != 1 {
				t.Fatalf("known prompt_submitted lost: %+v", r)
			}
			before := f.total()
			again, err := resume(in, f)
			if err != nil || !reflect.DeepEqual(again, r) || f.total() != before {
				t.Fatalf("retry after release failure was not idempotent readback: %+v %v", again, err)
			}
		})
		t.Run(op+" after precondition failure", func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			raw := stoppedAtTrust(t, in, f)
			in.Precondition = func() error { return errors.New("task moved") }
			_, err := resumeWithLockFault(in, f, op)
			if e := wantCode(t, err, ClassPrecondition); e.Stage != StagePrompt {
				t.Fatalf("original failure lost: %+v", e)
			}
			wantReleaseError(t, err)
			wantUnchanged(t, in, raw)
			in.Precondition = nil
			wantResumed(t, in, f, f.total())
		})
		t.Run(op+" on complete readback", func(t *testing.T) {
			in := testInput(t)
			f := newFake(t, in)
			first, err := run(in, f, nil)
			if err != nil {
				t.Fatal(err)
			}
			before := f.total()
			got, err := resumeWithLockFault(in, f, op)
			wantReleaseError(t, err)
			if !reflect.DeepEqual(got, first) || f.total() != before {
				t.Fatalf("readback changed or touched Herdr: %+v", got)
			}
		})
	}
}

func TestAgentResumeCancelledBeforePromptMarkerKeepsReceiptForLaterResume(t *testing.T) {
	t.Run("cancelled before observing a ready agent", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		raw := stoppedAtTrust(t, in, f)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := resumeAgent(ctx, in, f)
		wantCode(t, err, ClassNotReady)
		if f.count("agent", "prompt") != 0 {
			t.Fatalf("prompt attempted after cancellation: %q", f.calls)
		}
		wantUnchanged(t, in, raw)
		wantResumed(t, in, f, f.total())
	})
	t.Run("cancelled during a passing precondition", func(t *testing.T) {
		in := testInput(t)
		f := newFake(t, in)
		raw := stoppedAtTrust(t, in, f)
		ctx, cancel := context.WithCancel(context.Background())
		in.Precondition = func() error { cancel(); return nil }
		_, err := resumeAgent(ctx, in, f)
		wantCode(t, err, ClassNotReady)
		if f.count("agent", "prompt") != 0 {
			t.Fatalf("prompt attempted after cancellation: %q", f.calls)
		}
		wantUnchanged(t, in, raw)
		in.Precondition = nil
		wantResumed(t, in, f, f.total())
	})
}
