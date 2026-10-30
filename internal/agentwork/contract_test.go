package agentwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func artifact(t *testing.T, name string, v any) {
	t.Helper()
	p := filepath.Join("..", "..", "schemas", "agent-work-v1", "examples", "v1", name+".json")
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		t.Fatal(e)
	}
	b = append(b, '\n')
	if os.Getenv("PLY_UPDATE_AGENT_WORK_FIXTURES") == "1" {
		if e = os.MkdirAll(filepath.Dir(p), 0755); e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(p, b, 0644); e != nil {
			t.Fatal(e)
		}
	}
	old, e := os.ReadFile(p)
	if e != nil || !bytes.Equal(old, b) {
		t.Fatalf("fixture %s does not match actual serialization: %v", name, e)
	}
}
func TestAgentWorkContractFixtures(t *testing.T) {
	s, d, _ := fixture(t)
	d.PublicationKey = "review-1"
	artifact(t, "register-request", d)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "registered", f)
	m := mutation(f, "request-1")
	artifact(t, "request-mutation", m)
	f, e = s.Mutate(f.ID, "request", m)
	if e != nil {
		t.Fatal(e)
	}
	s.Runtime.(*fakeRuntime).agents[0].Status = "working"
	f, e = s.Dispatch(f.ID)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "deferred", f)
	s.Runtime.(*fakeRuntime).agents[0].Status = "done"
	f, e = s.Dispatch(f.ID)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "sent-without-recheck", f)
	m = mutation(f, "recheck-waiting")
	m.Outcome = "waiting"
	m.NextCheckAtUTC = s.Now().Add(time.Hour).Format(time.RFC3339)
	m.Done = "Checked the review source"
	m.Next = "Recheck when due"
	m.Blocker = "Review return is pending"
	artifact(t, "recheck-waiting-mutation", m)
	f, e = s.Mutate(f.ID, "recheck", m)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "recheck-waiting", f)
	m = mutation(f, "fulfilled")
	m.Outcome = "fulfilled"
	m.Done = "Review returned"
	m.Next = "Owner checks the native review"
	f, e = s.Mutate(f.ID, "recheck", m)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "recheck-fulfilled", f)
	// r001-f01: same canonical status obligation, new explicit control only after
	// a real recipient fulfilled report. The old dispatch subject stays archived.
	status, sd, _ := fixture(t)
	sd.WaitKind = "status"
	sd.WaitRef = sd.TargetRef
	sd.NextCheckAtUTC = status.now()
	sf, err := status.Register(sd)
	if err != nil {
		t.Fatal(err)
	}
	sf, err = status.Dispatch(sf.ID)
	if err != nil {
		t.Fatal(err)
	}
	sm := mutation(sf, "summary-fulfilled")
	sm.Outcome = "fulfilled"
	sm.Summary = "Checked actual fixture sources"
	sf, err = status.Mutate(sf.ID, "recheck", sm)
	if err != nil {
		t.Fatal(err)
	}
	artifact(t, "status-summary-fulfilled", sf)
	sm = mutation(sf, "later-explicit-status")
	artifact(t, "status-new-cycle-request", sm)
	sf, err = status.Mutate(sf.ID, "request", sm)
	if err != nil {
		t.Fatal(err)
	}
	artifact(t, "status-explicit-new-cycle", sf)

	u, ud, _ := fixture(t)
	ud.NextCheckAtUTC = u.now()
	uf, e := u.Register(ud)
	if e != nil {
		t.Fatal(e)
	}
	u.Send = func(Identity, string) error { return errors.New("lost ACK fixture") }
	uf, e = u.Dispatch(uf.ID)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "unknown-no-replay", uf)
	n := u.Runtime.(*fakeRuntime).agents[0]
	n.Session = "replacement-fixture"
	n.Pane = "replacement-pane-fixture"
	n.Terminal = "replacement-terminal-fixture"
	u.Runtime.(*fakeRuntime).agents = append(u.Runtime.(*fakeRuntime).agents, n)
	m = mutation(uf, "retarget-1")
	m.TargetRef = n.Ref()
	artifact(t, "retarget-mutation", m)
	uf, e = u.Mutate(uf.ID, "retarget", m)
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "retarget-consumed", uf)
	uf, e = u.Mutate(uf.ID, "cancel", mutation(uf, "cancel-1"))
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "cancelled", uf)
	v, e := u.Overview()
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "overview-bound-unbound", v)
	for _, state := range []string{"clear", "refusal", "partial"} {
		self, sd, _ := fixture(t)
		sd.NextCheckAtUTC = self.now()
		sd.Registrar = Registrar{Kind: AgentSelfRegistered, ActorRef: sd.TargetRef}
		identity := self.Runtime.(*fakeRuntime).agents[0].Identity
		coverage := NativeSnapshot{Available: true, VetoCoverageComplete: true, Vetoes: map[string]string{}}
		if state == "refusal" {
			coverage.Vetoes[identity.sessionKey()] = "native_refusal"
		}
		if state == "partial" {
			coverage.CoverageUnknown = true
			coverage.VetoCoverageComplete = false
			coverage.Partial = true
		}
		self.Native = func() NativeSnapshot { return coverage }
		sf, err := self.Register(sd)
		if err != nil {
			t.Fatal(err)
		}
		sf, err = self.Dispatch(sf.ID)
		if err != nil {
			t.Fatal(err)
		}
		artifact(t, "unbound-self-registered-"+state, sf)
	}

	partial := u.native()
	partial.Partial = true
	u.Native = func() NativeSnapshot { return partial }
	v, e = u.Overview()
	if e != nil {
		t.Fatal(e)
	}
	artifact(t, "overview-partial", v)
}
