package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func runtimeScript(t *testing.T, body string) HerdrRuntime {
	t.Helper()
	p := filepath.Join(t.TempDir(), "herdr-fixture")
	if e := os.WriteFile(p, []byte("#!/bin/sh\n"+body), 0700); e != nil {
		t.Fatal(e)
	}
	return HerdrRuntime{Executable: p, Timeout: time.Second}
}
func TestHerdrRuntimeMalformedTruncatedUnavailable(t *testing.T) {
	for _, tc := range []struct {
		name   string
		answer func(context.Context, []string) (string, error)
	}{
		{"malformed", func(context.Context, []string) (string, error) { return "synthetic-private-source", nil }},
		{"truncated", func(context.Context, []string) (string, error) { return `{"result":{"agents":[`, nil }},
		{"null", func(context.Context, []string) (string, error) { return `{"result":{"agents":null}}`, nil }},
		{"unavailable", func(context.Context, []string) (string, error) { return "", fmt.Errorf("synthetic-private-error") }},
		{"timeout", func(ctx context.Context, _ []string) (string, error) { <-ctx.Done(); return "", ctx.Err() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := fakeHerdr(tc.answer)
			h.Timeout = 50 * time.Millisecond
			if _, e := h.List(); e == nil || strings.Contains(e.Error(), "synthetic-private") {
				t.Fatal(e)
			}
		})
	}
}
func TestHerdrRuntimeSerializerDefaultsAndPrivacy(t *testing.T) {
	raw := `{"result":{"agents":[{"agent":"codex","agent_status":"done","workspace_id":"w1","tab_id":"w1:t1","pane_id":"w1:p1","terminal_id":"secret-terminal","interactive_ready":true,"agent_session":{"agent":"codex","kind":"id","value":"secret-session"},"title":"secret-title","tokens":{"token":"secret-token"}}]}}`
	h := fakeHerdr(func(context.Context, []string) (string, error) { return raw, nil })
	agents, e := h.List()
	if e != nil || len(agents) != 1 || agents[0].LaunchPending || agents[0].Ready == nil || !*agents[0].Ready {
		t.Fatal(agents, e)
	}
	s, d, _ := fixture(t)
	s.Runtime = h
	v, e := s.Overview()
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(v)
	if strings.Contains(string(b), "secret-") {
		t.Fatal("private source metadata escaped")
	}
	d.TargetRef = agents[0].Ref()
	if _, e = s.Register(s.onlyProject(), d); e != nil {
		t.Fatal(e)
	}
}
func TestAgentWorkLimitsAndPartialPreserveKnown(t *testing.T) {
	s, _, _ := fixture(t)
	rt := s.Runtime.(*fakeRuntime)
	one := rt.agents[0]
	for n := 0; n < 4097; n++ {
		rt.agents = append(rt.agents, one)
	}
	v, e := s.Overview()
	if e != nil || v.Sources[0].Status != "partial" {
		t.Fatal(v, e)
	}
	n := s.native()
	n.Partial = true
	s.Native = func() NativeSnapshot { return n }
	rt.agents = []Observation{one}
	v, e = s.Overview()
	if e != nil || v.Agents[0].Work == nil || v.Sources[1].Status != "partial" {
		t.Fatal(v, e)
	}
}
