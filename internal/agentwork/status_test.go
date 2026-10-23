package agentwork

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestAgentWorkPrivacyBindingAndReadOnly(t *testing.T) {
	s, _, _ := fixture(t)
	second := s.Runtime.(*fakeRuntime).agents[0]
	second.Session = "private-unbound"
	second.Pane = "private-pane"
	s.Runtime.(*fakeRuntime).agents = append(s.Runtime.(*fakeRuntime).agents, second)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, e := s.Overview()
			if e != nil {
				t.Error(e)
				return
			}
			b, _ := json.Marshal(v)
			for _, secret := range []string{"session-fixture", "terminal-fixture", "private-unbound"} {
				if strings.Contains(string(b), secret) {
					t.Error("private identity leaked")
				}
			}
			if len(v.Agents) != 2 || v.Agents[0].Work == nil && v.Agents[1].Work == nil {
				t.Error(v)
			}
			for _, a := range v.Agents {
				if a.Status == "done" && a.Work != nil && a.Work.Phase != "working" {
					t.Error("turn marked native delivered")
				}
			}
		}()
	}
	wg.Wait()
	db, e := s.read()
	if e != nil || db.Sequence != 0 {
		t.Fatal("read created state", e)
	}
}
func TestAgentWorkSourceOutcomesAndReusedPane(t *testing.T) {
	s, _, _ := fixture(t)
	rt := s.Runtime.(*fakeRuntime)
	rt.agents = nil
	v, e := s.Overview()
	if e != nil || v.Sources[0].Status != "empty" {
		t.Fatal(v, e)
	}
	rt.err = errors.New("secret source output")
	v, e = s.Overview()
	b, _ := json.Marshal(v)
	if e != nil || v.Sources[0].Status != "unavailable" || strings.Contains(string(b), "secret") {
		t.Fatal(v, e)
	}
	rt.err = nil
	rt.agents = []Observation{{Identity: Identity{Provider: "codex", Pane: "reuse", Session: "replacement", Workspace: "w", Tab: "t", Terminal: "term"}, Status: "done"}}
	v, e = s.Overview()
	if e != nil || v.Agents[0].Work != nil || v.Agents[0].BindingStatus != "unbound" {
		t.Fatal(v, e)
	}
}
