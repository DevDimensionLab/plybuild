package agentwork

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"testing"
)

func ptr(s string) *string { return &s }
func TestNativeVetoScopePreservesAWhenKnownBSourceFails(t *testing.T) {
	s, _, _ := fixture(t)
	a := s.Runtime.(*fakeRuntime).agents[0].Identity
	b := taskrun.InventoryRun{Provider: ptr("claude"), Herdr: &taskrun.InventoryHerdr{AgentSessionID: ptr("session-B")}, Freshness: "unknown", StatusIdentityKnown: true}
	n := NativeSnapshot{Available: true, Vetoes: map[string]string{}}
	scopeInventory(&n, taskrun.Inventory{Runs: []taskrun.InventoryRun{b}}, nil)
	if !n.VetoCoverageComplete || n.CoverageUnknown {
		t.Fatal("known B failure became global")
	}
	if _, v := binding(n, a); v != "clear" {
		t.Fatal(v)
	}
	n.Vetoes[a.sessionKey()] = "native_question"
	if _, v := binding(n, a); v != "native_question" {
		t.Fatal("own question lost", v)
	}
	n = NativeSnapshot{Available: true, Vetoes: map[string]string{}}
	b.StatusIdentityKnown = false
	scopeInventory(&n, taskrun.Inventory{Runs: []taskrun.InventoryRun{b}}, nil)
	if _, v := binding(n, a); v == "clear" {
		t.Fatal("unidentifiable source supplied absence")
	}
	n = NativeSnapshot{Available: true, Vetoes: map[string]string{}}
	b.StatusIdentityKnown = true
	b.Freshness = "stale"
	scopeInventory(&n, taskrun.Inventory{Runs: []taskrun.InventoryRun{b}}, nil)
	if _, v := binding(n, a); v == "clear" {
		t.Fatal("changed identity source supplied absence")
	}
}

func TestFollowupOwnQuestionAndCoverageUnknownVetoBoundA(t *testing.T) {
	for _, reason := range []string{"native_question", "native_sources_unknown", "coverage_unknown"} {
		t.Run(reason, func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.now()
			a := s.Runtime.(*fakeRuntime).agents[0].Identity
			n := s.native()
			n.Partial = true
			n.VetoCoverageComplete = true
			n.Vetoes = map[string]string{}
			if reason == "coverage_unknown" {
				n.CoverageUnknown = true
			} else {
				n.Vetoes[a.sessionKey()] = reason
			}
			s.Native = func() NativeSnapshot { return n }
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			f, e = s.Dispatch(f.ID)
			if e != nil || c.Load() != 0 {
				t.Fatal(f, e)
			}
		})
	}
}
func TestNativeSameSessionVetoCannotBeHiddenByEarlierClearAttempt(t *testing.T) {
	s, _, _ := fixture(t)
	a := s.Runtime.(*fakeRuntime).agents[0].Identity
	good := taskrun.InventoryRun{Provider: ptr(a.Provider), Herdr: &taskrun.InventoryHerdr{AgentSessionID: ptr(a.Session)}, Freshness: "fresh", StatusIdentityKnown: true, StatusFollowupVeto: "clear"}
	bad := good
	bad.Freshness = "unknown"
	n := NativeSnapshot{Available: true, Vetoes: map[string]string{}}
	scopeInventory(&n, taskrun.Inventory{Runs: []taskrun.InventoryRun{good, bad}}, nil)
	if n.Vetoes[a.sessionKey()] != "native_sources_unknown" {
		t.Fatal("earlier clear record hid matching uncertain source", n.Vetoes)
	}
}
