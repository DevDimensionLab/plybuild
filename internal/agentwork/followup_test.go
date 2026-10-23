package agentwork

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func fixture(t *testing.T) (*Service, Draft, *atomic.Int32) {
	t.Helper()
	now := time.Date(2026, 10, 23, 13, 0, 0, 0, time.UTC)
	target := Identity{Provider: "codex", Session: "session-fixture", Workspace: "w-fixture", Tab: "tab-fixture", Pane: "pane-fixture", Terminal: "terminal-fixture"}
	calls := &atomic.Int32{}
	rt := &fakeRuntime{agents: []Observation{{Identity: target, Status: "idle", Ready: boolp(true)}}}
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	s := &Service{Root: root, Now: func() time.Time { return now }, Runtime: rt, Native: func() NativeSnapshot {
		return NativeSnapshot{Available: true, Bindings: []Binding{{Identity: target, Work: Work{ProjectID: "project-fixture", RepoID: "repo-fixture", EpicID: "epic-fixture", TaskID: "task-fixture", RunID: "run-fixture", AttemptRef: "run-fixture", Phase: "working", NextActor: "agent", NextAction: "continue_delivery", Reason: "Read the registered native action before continuing.", DetailRef: "workflow.trace:task-fixture", Source: "registered", Freshness: "fresh"}, Veto: "clear"}}}
	}, Send: func(Identity, string) error { calls.Add(1); return nil }}
	d := Draft{PublicationKey: "review/1", OwnerRef: "owner-fixture", MandateRef: "mandate-fixture", Registrar: Registrar{Kind: HumanRequested, ActorRef: "human-fixture"}, TargetRef: target.Ref(), WaitKind: "review", WaitRef: "review-fixture", NextCheckAtUTC: now.Add(time.Minute).Format(time.RFC3339)}
	return s, d, calls
}

type fakeRuntime struct {
	agents []Observation
	err    error
}

func (f *fakeRuntime) List() ([]Observation, error) { return f.agents, f.err }
func (f *fakeRuntime) Get(i Identity) (Observation, error) {
	for _, a := range f.agents {
		if a.Identity == i {
			return a, nil
		}
	}
	return Observation{}, errors.New("missing")
}
func boolp(b bool) *bool { return &b }
func mutation(f Followup, id string) Mutation {
	return Mutation{RequestID: id, OwnerRef: f.OwnerRef, ExpectedRevision: f.Revision, Generation: f.Generation, TargetRef: f.TargetRef}
}
func TestFollowupCycleNoReplayAndExplicitRecheck(t *testing.T) {
	s, d, calls := fixture(t)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.Register(d)
	if e != nil || again.ID != f.ID {
		t.Fatal(again, e)
	}
	alias := d
	alias.PublicationKey = "alias"
	again, e = s.Register(alias)
	if e != nil || again.ID != f.ID {
		t.Fatal("alias escaped canonical obligation", e)
	}
	changed := d
	changed.OwnerRef = "other"
	changed.PublicationKey = d.PublicationKey
	if _, e = s.Register(changed); e == nil {
		t.Fatal("key takeover")
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || f.Transport.State != "not_sent" || calls.Load() != 0 {
		t.Fatal(f, e)
	}
	f, e = s.Mutate(f.ID, "request", mutation(f, "explicit"))
	if e != nil {
		t.Fatal(e)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || f.Transport.State != "sent" || calls.Load() != 1 || f.Recheck != nil {
		t.Fatal(f, e)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || calls.Load() != 1 {
		t.Fatal("replayed", e)
	}
	m := mutation(f, "waiting")
	m.Outcome = "waiting"
	m.Done = "Checked actual sources"
	m.Next = "Wait for review"
	m.Blocker = "Review is pending"
	m.NextCheckAtUTC = s.Now().Add(time.Minute).Format(time.RFC3339)
	f, e = s.Mutate(f.ID, "recheck", m)
	if e != nil || f.Cycle != 2 || f.Recheck.Provenance != "agent_self_report" {
		t.Fatal(f, e)
	}
	if _, e = s.Mutate(f.ID, "cancel", m); e == nil {
		t.Fatal("request bytes reused for another operation")
	}
}
func TestFollowupConcurrentDispatchAndLostACK(t *testing.T) {
	s, d, calls := fixture(t)
	d.NextCheckAtUTC = s.Now().Format(time.RFC3339)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	s.Send = func(Identity, string) error { calls.Add(1); return errors.New("lost ACK") }
	var wg sync.WaitGroup
	for n := 0; n < 12; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, e := s.Dispatch(f.ID); e != nil {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	got, e := s.Show(f.ID)
	if e != nil || got.Transport.State != "unknown" || calls.Load() != 1 {
		t.Fatal(got, e, calls.Load())
	}
	alias := d
	alias.PublicationKey = "new-ui-id"
	got, e = s.Register(alias)
	if e != nil {
		t.Fatal(e)
	}
	got, e = s.Dispatch(got.ID)
	if e != nil || calls.Load() != 1 {
		t.Fatal("unknown replayed", e)
	}
	restart := *s
	got, e = restart.Dispatch(got.ID)
	if e != nil || calls.Load() != 1 {
		t.Fatal("restart replayed", e)
	}
}
func TestFollowupVetoAndTargetRevalidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*Service)
	}{
		{"working", func(s *Service) { s.Runtime.(*fakeRuntime).agents[0].Status = "working" }},
		{"readiness_unknown", func(s *Service) { s.Runtime.(*fakeRuntime).agents[0].Ready = nil }},
		{"launch_pending", func(s *Service) { s.Runtime.(*fakeRuntime).agents[0].LaunchPending = true }},
		{"new_session", func(s *Service) { s.Runtime.(*fakeRuntime).agents[0].Session = "reused-pane" }},
		{"native_unknown", func(s *Service) { s.Native = func() NativeSnapshot { return NativeSnapshot{} } }},
		{"question", func(s *Service) {
			n := s.Native()
			n.Bindings[0].Veto = "question"
			s.Native = func() NativeSnapshot { return n }
		}},
		{"refusal", func(s *Service) {
			n := s.Native()
			n.Bindings[0].Veto = "refusal"
			s.Native = func() NativeSnapshot { return n }
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.Now().Format(time.RFC3339)
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			tc.change(s)
			f, e = s.Dispatch(f.ID)
			if e != nil || c.Load() != 0 || f.Transport.State == "sent" {
				t.Fatal(f, e)
			}
		})
	}
}
func TestFollowupCrashClaimSurvivesAndCancelWins(t *testing.T) {
	s, d, c := fixture(t)
	d.NextCheckAtUTC = s.Now().Format(time.RFC3339)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	s.Fault = func(stage string) error {
		if stage == "after_claim" {
			return errors.New("crash")
		}
		return nil
	}
	if _, e = s.Dispatch(f.ID); e == nil {
		t.Fatal("fault missed")
	}
	s.Fault = nil
	f, e = s.Dispatch(f.ID)
	if e != nil || f.Transport.State != "unknown" || c.Load() != 0 {
		t.Fatal(f, e)
	}
	f, e = s.Mutate(f.ID, "cancel", mutation(f, "cancel"))
	if e != nil {
		t.Fatal(e)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || f.State != "cancelled" || c.Load() != 0 {
		t.Fatal(f, e)
	}
}

func TestFollowupFaultBoundariesAndRecipientRevision(t *testing.T) {
	for _, stage := range []string{"before_claim", "after_claim", "after_input", "after_ack"} {
		t.Run(stage, func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.now()
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			s.Fault = func(at string) error {
				if at == stage {
					return errors.New("fixture crash")
				}
				return nil
			}
			if _, e = s.Dispatch(f.ID); e == nil {
				t.Fatal("fault missing")
			}
			s.Fault = nil
			f, e = s.Dispatch(f.ID)
			if e != nil {
				t.Fatal(e)
			}
			if stage == "before_claim" {
				if c.Load() != 1 || f.Transport.State != "sent" {
					t.Fatal(f, c.Load())
				}
			} else {
				if f.Transport.State != "unknown" || c.Load() > 1 {
					t.Fatal(f, c.Load())
				}
			}
		})
	}
	s, d, c := fixture(t)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	m := mutation(f, "cancel")
	m.OwnerRef = "other"
	if _, e = s.Mutate(f.ID, "cancel", m); e == nil {
		t.Fatal("owner bypass")
	}
	m = mutation(f, "cancel")
	f, e = s.Mutate(f.ID, "cancel", m)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Mutate(f.ID, "recheck", m); e == nil {
		t.Fatal("operation takeover")
	}
	if _, e = s.Dispatch(f.ID); e != nil || c.Load() != 0 {
		t.Fatal(e)
	}
}
func TestFollowupLostMutationReplyAndNoTaskStatusCopy(t *testing.T) {
	s, d, _ := fixture(t)
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	m := mutation(f, "make-due")
	f, e = s.Mutate(f.ID, "request", m)
	if e != nil {
		t.Fatal(e)
	}
	same, e := s.RequestStatus(f.ID, m.RequestID)
	if e != nil || same.Revision != f.Revision {
		t.Fatal(same, e)
	}
	replayed, e := s.Mutate(f.ID, "request", m)
	if e != nil || replayed.Revision != f.Revision {
		t.Fatal(replayed, e)
	}
	db, e := s.read()
	if e != nil || db.Records[f.ID].Value.Work != nil {
		t.Fatal("native status copied to followup store")
	}
}
func TestFollowupUnboundSelfRegisteredRuntime(t *testing.T) {
	s, d, c := fixture(t)
	d.NextCheckAtUTC = s.now()
	d.Registrar = Registrar{Kind: AgentSelfRegistered, ActorRef: d.TargetRef}
	s.Native = func() NativeSnapshot {
		return NativeSnapshot{Available: true, VetoCoverageComplete: true, Vetoes: map[string]string{}}
	}
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	if f.Work != nil {
		t.Fatal("invented Task")
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || f.Transport.State != "sent" || c.Load() != 1 {
		t.Fatal("unbound exact runtime cannot receive reminder", f, e)
	}
	for _, veto := range []string{"question", "refusal", "unknown"} {
		t.Run(veto, func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.now()
			d.Registrar = Registrar{Kind: AgentSelfRegistered, ActorRef: d.TargetRef}
			identity := s.Runtime.(*fakeRuntime).agents[0].Identity
			n := NativeSnapshot{Available: true, VetoCoverageComplete: true, Vetoes: map[string]string{identity.sessionKey(): veto}}
			if veto == "unknown" {
				n.VetoCoverageComplete = false
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

// R1 regression: unrelated historical diagnostics do not become this target's
// veto. Own question or unidentifiable relevant coverage must still prevent input.
func TestFollowupTargetSpecificVetoCoverage(t *testing.T) {
	for _, bound := range []bool{true, false} {
		name := "unbound"
		if bound {
			name = "bound"
		}
		t.Run(name, func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.now()
			d.Registrar = Registrar{Kind: AgentSelfRegistered, ActorRef: d.TargetRef}
			a := s.Runtime.(*fakeRuntime).agents[0].Identity
			n := s.native()
			n.Partial = true
			n.VetoCoverageComplete = true
			n.Vetoes = map[string]string{"claude:old-target-b": "native_sources_unknown"}
			if !bound {
				n.Bindings = nil
			}
			s.Native = func() NativeSnapshot { return n }
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			f, e = s.Dispatch(f.ID)
			if e != nil || c.Load() != 1 || f.Transport.State != "sent" {
				t.Fatal("unrelated B diagnostic vetoed A", f, e)
			}
			_ = a
		})
	}
}
func TestAgentWorkRosterFixturesHaveDistinctCurrentPanes(t *testing.T) {
	b, e := os.ReadFile("../../schemas/agent-work-v1/examples/v1/overview-bound-unbound.json")
	if e != nil {
		t.Fatal(e)
	}
	var v Overview
	if e = json.Unmarshal(b, &v); e != nil {
		t.Fatal(e)
	}
	seen := map[string]bool{}
	for _, a := range v.Agents {
		key := a.Workspace + ":" + a.Pane
		if seen[key] {
			t.Fatal("R2: two current agents share one pane")
		}
		seen[key] = true
	}
}
func TestFollowupRetargetPreservesOriginalDispatchSubject(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		t.Run(fmt.Sprint(unknown), func(t *testing.T) {
			s, d, c := fixture(t)
			d.NextCheckAtUTC = s.now()
			if unknown {
				s.Send = func(Identity, string) error { c.Add(1); return errors.New("lost ACK") }
			}
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			f, e = s.Dispatch(f.ID)
			if e != nil {
				t.Fatal(e)
			}
			oldTarget, oldGeneration := f.TargetRef, f.Generation
			next := s.Runtime.(*fakeRuntime).agents[0]
			next.Session = "retarget-session"
			next.Pane = "retarget-pane"
			next.Terminal = "retarget-terminal"
			s.Runtime.(*fakeRuntime).agents = append(s.Runtime.(*fakeRuntime).agents, next)
			m := mutation(f, "retarget-subject")
			m.TargetRef = next.Ref()
			f, e = s.Mutate(f.ID, "retarget", m)
			if e != nil {
				t.Fatal(e)
			}
			shown, e := s.Show(f.ID)
			if e != nil {
				t.Fatal(e)
			}
			b, _ := json.Marshal(shown)
			var wire map[string]any
			_ = json.Unmarshal(b, &wire)
			transport := wire["transport"].(map[string]any)
			subject, ok := transport["dispatch_subject"].(map[string]any)
			if !ok || subject["agent_ref"] != oldTarget || subject["generation"] != float64(oldGeneration) {
				t.Fatal("old attempt reassigned to new target", transport)
			}
			if shown.TargetRef == oldTarget || shown.Generation != oldGeneration+1 {
				t.Fatal(shown)
			}
			if _, e = s.Dispatch(f.ID); e != nil || c.Load() != 1 {
				t.Fatal("retarget replayed", e, c.Load())
			}
		})
	}
}
func TestFollowupStatusAliasAfterRetargetCannotOpenAnotherSend(t *testing.T) {
	s, d, c := fixture(t)
	d.WaitKind = "status"
	d.WaitRef = d.TargetRef
	d.NextCheckAtUTC = s.now()
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil {
		t.Fatal(e)
	}
	next := s.Runtime.(*fakeRuntime).agents[0]
	next.Session = "new-target"
	next.Pane = "new-pane"
	s.Runtime.(*fakeRuntime).agents = append(s.Runtime.(*fakeRuntime).agents, next)
	m := mutation(f, "new-generation")
	m.TargetRef = next.Ref()
	if _, e = s.Mutate(f.ID, "retarget", m); e != nil {
		t.Fatal(e)
	}
	d.PublicationKey = "new-ui-action"
	d.TargetRef = next.Ref()
	d.WaitRef = d.TargetRef
	if _, e = s.Register(d); e == nil {
		t.Fatal("retargeted consumed status obligation acquired a new alias cycle")
	}
	if c.Load() != 1 {
		t.Fatal(c.Load())
	}
}

// r001-f01: a fulfilled status answer closes only that cycle; a later explicit
// owner request can open a new control, with all prior send identity preserved.
func TestFollowupExplicitStatusRequestAfterFulfilled(t *testing.T) {
	s, d, c := fixture(t)
	d.WaitKind = "status"
	d.WaitRef = d.TargetRef
	d.NextCheckAtUTC = s.now()
	f, e := s.Register(d)
	if e != nil {
		t.Fatal(e)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil {
		t.Fatal(e)
	}
	firstTransport := f.Transport
	done := mutation(f, "recipient-fulfilled")
	done.Outcome = "fulfilled"
	done.Summary = "Checked the actual sources"
	f, e = s.Mutate(f.ID, "recheck", done)
	if e != nil {
		t.Fatal(e)
	}
	// Register/read/dispatch are not an automatic rearm after fulfillment.
	again, e := s.Register(d)
	if e != nil || again.State != "fulfilled" || again.Cycle != 1 {
		t.Fatal(again, e)
	}
	if _, e = s.Dispatch(f.ID); e != nil || c.Load() != 1 {
		t.Fatal("closed cycle auto-rearmed", e)
	}
	req := mutation(f, "explicit-second-status")
	f, e = s.Mutate(f.ID, "request", req)
	if e != nil {
		t.Fatal("legitimate new status control rejected", e)
	}
	if f.State != "waiting" || f.Cycle != 2 || f.Transport.State != "not_sent" || f.Transport.DispatchSubject != nil || f.Transport.CycleRef == firstTransport.CycleRef {
		t.Fatal(f)
	}
	raw, _ := json.Marshal(f)
	var wire map[string]any
	_ = json.Unmarshal(raw, &wire)
	history, ok := wire["history"].([]any)
	if !ok || len(history) != 1 {
		t.Fatal("previous cycle receipt missing", string(raw))
	}
	prior := history[0].(map[string]any)
	if prior["cycle"] != float64(1) || prior["transport"].(map[string]any)["cycle_ref"] != firstTransport.CycleRef || prior["recheck"].(map[string]any)["outcome"] != "fulfilled" {
		t.Fatal(prior)
	}
	replay, e := s.Mutate(f.ID, "request", req)
	if e != nil || replay.Cycle != 2 || replay.Revision != f.Revision {
		t.Fatal("lost new-cycle reply opened another cycle", replay, e)
	}
	if f.Work == nil || f.Work.Phase != "working" {
		t.Fatal("native work changed", f.Work)
	}
	f, e = s.Dispatch(f.ID)
	if e != nil || f.Transport.State != "sent" || c.Load() != 2 {
		t.Fatal(f, e, c.Load())
	}
	if _, e = s.Mutate(f.ID, "request", mutation(f, "new-ui-alias")); e == nil {
		t.Fatal("consumed second cycle rearmed without receipt")
	}
}
func TestFollowupExplicitNewStatusCycleStillRejectsUnknownAndCancellation(t *testing.T) {
	for _, state := range []string{"unknown", "cancelled"} {
		t.Run(state, func(t *testing.T) {
			s, d, _ := fixture(t)
			d.WaitKind = "status"
			d.WaitRef = d.TargetRef
			d.NextCheckAtUTC = s.now()
			if state == "unknown" {
				s.Send = func(Identity, string) error { return errors.New("lost ACK") }
			}
			f, e := s.Register(d)
			if e != nil {
				t.Fatal(e)
			}
			f, e = s.Dispatch(f.ID)
			if e != nil {
				t.Fatal(e)
			}
			if state == "cancelled" {
				f, e = s.Mutate(f.ID, "cancel", mutation(f, "cancel"))
				if e != nil {
					t.Fatal(e)
				}
			}
			if _, e = s.Mutate(f.ID, "request", mutation(f, "new-explicit-id")); e == nil {
				t.Fatal("unfulfilled cycle rearmed")
			}
			got, e := s.Show(f.ID)
			if e != nil || got.Cycle != 1 || len(got.History) != 0 {
				t.Fatal(got, e)
			}
		})
	}
}
