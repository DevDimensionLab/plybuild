package agentwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
	binary, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	qualifyTestCallback(t, root, binary)
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
		{"refusal", func(s *Service) {
			n := s.Native()
			n.Bindings[0].Veto = "refusal"
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
	for _, veto := range []string{"refusal", "unknown"} {
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
// veto. Actual refusal or unidentifiable relevant coverage must still prevent input.
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
func seedHistoricalRetarget(t *testing.T, s *Service, f Followup, next Identity) (Followup, Mutation) {
	t.Helper()
	m := mutation(f, "historical-retarget")
	m.TargetRef = next.Ref()
	err := s.locked(func(db *database) error {
		r := db.Records[f.ID]
		f.Generation++
		f.Revision++
		f.TargetRef = next.Ref()
		f.Work = nil
		r.Value = f
		r.Target = next
		r.TargetRefs = append(r.TargetRefs, next.Ref())
		r.Requests[m.RequestID] = operation{hash(struct {
			Op       string
			Mutation Mutation
		}{"retarget", m}), f}
		db.Records[f.ID] = r
		return s.save(db)
	})
	if err != nil {
		t.Fatal(err)
	}
	return f, m
}

func TestFollowupHistoricalRetarget(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		t.Run(fmt.Sprint(unknown), func(t *testing.T) {
			s, d, calls := fixture(t)
			d.WaitKind, d.WaitRef, d.NextCheckAtUTC = "status", d.TargetRef, s.now()
			if unknown {
				s.Send = func(Identity, string) error { calls.Add(1); return errors.New("lost ACK") }
			}
			f, err := s.Register(d)
			if err != nil {
				t.Fatal(err)
			}
			f, err = s.Dispatch(f.ID)
			if err != nil {
				t.Fatal(err)
			}
			originalTransport := f.Transport
			next := s.Runtime.(*fakeRuntime).agents[0]
			next.Session, next.Pane, next.Terminal = "historical-session", "historical-pane", "historical-terminal"
			s.Runtime.(*fakeRuntime).agents = append(s.Runtime.(*fakeRuntime).agents, next)
			historical, m := seedHistoricalRetarget(t, s, f, next.Identity)
			original, err := s.RequestStatus(f.ID, m.RequestID)
			if err != nil || original.Generation != 2 || original.TargetRef != next.Ref() {
				t.Fatal(original, err)
			}
			beforeDB, err := s.read()
			if err != nil {
				t.Fatal(err)
			}
			if _, err = s.Mutate(f.ID, "retarget", mutation(historical, "new-while-waiting")); err == nil || err.Error() != "followup_unknown_operation" {
				t.Fatal("new retarget while waiting", err)
			}
			if _, err = s.Dispatch(f.ID); err != nil || calls.Load() != 1 {
				t.Fatal("historical consumed cycle replay", err, calls.Load())
			}
			afterDB, err := s.read()
			if err != nil || !reflect.DeepEqual(beforeDB, afterDB) {
				t.Fatal("historical waiting record changed", err)
			}
			_, err = s.Mutate(f.ID, "cancel", mutation(historical, "later-cancel"))
			if err != nil {
				t.Fatal(err)
			}
			snapshot := func() []byte {
				t.Helper()
				db, err := s.read()
				if err != nil {
					t.Fatal(err)
				}
				b, err := os.ReadFile(filepath.Join(s.dir(), "revisions", fmt.Sprintf("%012d.json", db.Sequence)))
				if err != nil {
					t.Fatal(err)
				}
				entries, err := os.ReadDir(filepath.Join(s.dir(), "revisions"))
				if err != nil || len(entries) != db.Sequence {
					t.Fatal(entries, err)
				}
				return b
			}
			before := snapshot()
			for _, mutation := range []Mutation{m, mutation(historical, "unknown-retarget")} {
				got, err := s.Mutate(f.ID, "retarget", mutation)
				if mutation.RequestID == m.RequestID {
					if err != nil || !reflect.DeepEqual(got, original) {
						t.Fatal("original replay", got, err)
					}
				} else if err == nil || err.Error() != "followup_unknown_operation" {
					t.Fatal("new retarget", err)
				}
			}
			changed := m
			changed.TargetRef = d.TargetRef
			if _, err := s.Mutate(f.ID, "retarget", changed); err == nil || err.Error() != "followup_request_conflict" {
				t.Fatal("request conflict", err)
			}
			got, err := s.RequestStatus(f.ID, m.RequestID)
			if err != nil || !reflect.DeepEqual(got, original) {
				t.Fatal(got, err)
			}
			shown, err := s.Show(f.ID)
			if err != nil || shown.Generation != 2 || shown.TargetRef != next.Ref() || !reflect.DeepEqual(shown.Transport, originalTransport) {
				t.Fatal(shown, err)
			}
			listed, err := s.List(false)
			if err != nil || len(listed.Items) != 1 || !reflect.DeepEqual(listed.Items[0], shown) {
				t.Fatal(listed, err)
			}
			overview, err := s.Overview()
			if err != nil {
				t.Fatal(err)
			}
			for _, value := range []Followup{shown, original, listed.Items[0], overview.Followups[0]} {
				if _, ok := value.Capabilities["retarget"]; ok {
					t.Fatal("retarget capability", value)
				}
			}
			if _, err = s.Dispatch(f.ID); err != nil || calls.Load() != 1 {
				t.Fatal("replayed dispatch", err, calls.Load())
			}
			d.PublicationKey, d.TargetRef, d.WaitRef = "alias-new-target", next.Ref(), next.Ref()
			if _, err = s.Register(d); err == nil || err.Error() != "followup_status_target_already_registered" {
				t.Fatal("historical lineage", err)
			}
			if !bytes.Equal(before, snapshot()) {
				t.Fatal("historical read or rejected mutation wrote store")
			}
			db, err := s.read()
			if err != nil || db.Records[f.ID].Target != next.Identity || len(db.Records[f.ID].TargetRefs) != 2 {
				t.Fatal(db, err)
			}
		})
	}
}

func TestFollowupStatusCancelRegisterNewTarget(t *testing.T) {
	s, d, calls := fixture(t)
	d.WaitKind, d.WaitRef, d.NextCheckAtUTC = "status", d.TargetRef, s.now()
	f, err := s.Register(d)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Mutate(f.ID, "cancel", mutation(f, "cancel-old-target")); err != nil {
		t.Fatal(err)
	}
	next := s.Runtime.(*fakeRuntime).agents[0]
	next.Session, next.Pane, next.Terminal = "new-session", "new-pane", "new-terminal"
	s.Runtime.(*fakeRuntime).agents = append(s.Runtime.(*fakeRuntime).agents, next)
	s.Native = func() NativeSnapshot {
		return NativeSnapshot{Available: true, VetoCoverageComplete: true, Vetoes: map[string]string{}}
	}
	s.Send = func(target Identity, _ string) error {
		if target != next.Identity {
			t.Fatal("wrong target", target)
		}
		calls.Add(1)
		return nil
	}
	d.PublicationKey, d.TargetRef, d.WaitRef = "new-status-target", next.Ref(), next.Ref()
	replacement, err := s.Register(d)
	if err != nil || replacement.ID == f.ID {
		t.Fatal(replacement, err)
	}
	for _, id := range []string{f.ID, replacement.ID, replacement.ID, f.ID} {
		if _, err = s.Dispatch(id); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatal(calls.Load())
	}
	sent, err := s.Show(replacement.ID)
	if err != nil || sent.Transport.State != "sent" || sent.Transport.DispatchSubject.AgentRef != next.Ref() {
		t.Fatal(sent, err)
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
