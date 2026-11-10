package agentwork

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"
)

func plain(v string) bool {
	return v != "" && len(utf16.Encode([]rune(v))) <= 256 && strings.TrimSpace(v) == v && !strings.ContainsAny(v, "/\\\n\r\t") && !strings.Contains(v, "://") && !strings.Contains(v, "@") && !strings.ContainsFunc(v, unicode.IsControl)
}
func reference(v string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.:-]{0,127}$`).MatchString(v)
}
func timeValid(v string) bool {
	t, e := time.Parse(time.RFC3339Nano, v)
	return e == nil && t.Location() == time.UTC && strings.HasSuffix(v, "Z")
}
func (s *Service) isDue(f Followup) bool {
	t, e := time.Parse(time.RFC3339Nano, f.NextCheckAtUTC)
	return e == nil && !s.Now().Before(t)
}
func (s *Service) native() NativeSnapshot {
	if s.Native == nil {
		return NativeSnapshot{}
	}
	return s.Native()
}
func binding(n NativeSnapshot, i Identity) (*Work, string) {
	var found *Binding
	for j := range n.Bindings {
		b := &n.Bindings[j]
		if b.Identity == i {
			if found != nil {
				return nil, "ambiguous_binding"
			}
			found = b
		}
	}
	if !n.Available {
		return nil, "native_coverage_unknown"
	}
	if found == nil {
		if !n.CoverageUnknown && n.VetoCoverageComplete {
			if veto := n.Vetoes[i.sessionKey()]; veto != "" {
				return nil, veto
			}
			return nil, "clear"
		}
		return nil, "native_veto_coverage_unknown"
	}
	w := found.Work
	if n.CoverageUnknown {
		return &w, "native_veto_coverage_unknown"
	}
	if veto := n.Vetoes[i.sessionKey()]; veto != "" && veto != "registered_session_requires_exact_binding" {
		return &w, veto
	}
	return &w, found.Veto
}
func (s *Service) find(ref string) (Identity, *Work, error) {
	if s.Runtime == nil {
		return Identity{}, nil, fmt.Errorf("agent_source_unavailable")
	}
	agents, e := s.Runtime.List()
	if e != nil {
		return Identity{}, nil, fmt.Errorf("agent_source_unavailable")
	}
	var target *Identity
	for _, a := range agents {
		if a.Ref() == ref && a.valid() {
			if target != nil {
				return Identity{}, nil, fmt.Errorf("agent_binding_ambiguous")
			}
			i := a.Identity
			target = &i
		}
	}
	if target == nil {
		return Identity{}, nil, fmt.Errorf("agent_target_unavailable")
	}
	w, _ := binding(s.native(), *target)
	return *target, w, nil
}
func (s *Service) Register(d Draft) (Followup, error) {
	var out Followup
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_./:-]{0,127}$`).MatchString(d.PublicationKey) || !reference(d.OwnerRef) || !reference(d.MandateRef) || !reference(d.Registrar.ActorRef) || (d.Registrar.Kind != HumanRequested && d.Registrar.Kind != AgentSelfRegistered) || !reference(d.WaitKind) || !reference(d.WaitRef) || !regexp.MustCompile(`^agr_[a-f0-9]{64}$`).MatchString(d.TargetRef) || !timeValid(d.NextCheckAtUTC) {
		return out, fmt.Errorf("followup_invalid_draft")
	}
	if d.Registrar.Kind == AgentSelfRegistered && d.Registrar.ActorRef != d.TargetRef {
		return out, fmt.Errorf("followup_self_registrar_must_match_target")
	}
	if d.WaitKind == "status" && d.WaitRef != d.TargetRef {
		return out, fmt.Errorf("followup_status_requires_canonical_target_ref")
	}
	id := "fup_" + hash([]string{d.OwnerRef, d.MandateRef, d.WaitKind, d.WaitRef})
	e := s.locked(func(db *database) error {
		if d.WaitKind == "status" {
			for existingID, r := range db.Records {
				if existingID == id || r.Value.OwnerRef != d.OwnerRef || r.Value.MandateRef != d.MandateRef || r.Value.WaitKind != "status" {
					continue
				}
				for _, ref := range r.TargetRefs {
					if ref == d.TargetRef {
						return fmt.Errorf("followup_status_target_already_registered")
					}
				}
			}
		}
		if old, ok := db.Keys[d.PublicationKey]; ok && old != id {
			return fmt.Errorf("followup_publication_conflict")
		}
		if r, ok := db.Records[id]; ok {
			a, b := r.Draft, d
			a.PublicationKey = ""
			b.PublicationKey = ""
			if hash(a) != hash(b) {
				return fmt.Errorf("followup_registration_conflict")
			}
			out = r.Value
			if _, ok := db.Keys[d.PublicationKey]; !ok {
				db.Keys[d.PublicationKey] = id
				return s.save(db)
			}
			return nil
		}
		if len(db.Records) >= 512 {
			return fmt.Errorf("followup_record_limit")
		}
		target, _, e := s.find(d.TargetRef)
		if e != nil {
			return e
		}
		out = Followup{Kind: FollowupKind, SchemaVersion: 1, ID: id, Revision: 1, Generation: 1, OwnerRef: d.OwnerRef, MandateRef: d.MandateRef, Registrar: d.Registrar, TargetRef: d.TargetRef, WaitKind: d.WaitKind, WaitRef: d.WaitRef, Work: nil, NextCheckAtUTC: d.NextCheckAtUTC, State: "waiting", Cycle: 1, Transport: Transport{State: "not_sent", Reason: "Registered; no input attempted.", CycleRef: cycleRef(id, 1)}, RegisteredAtUTC: s.now()}
		db.Records[id] = record{Value: out, Target: target, TargetRefs: []string{d.TargetRef}, Draft: d, Requests: map[string]operation{}}
		db.Keys[d.PublicationKey] = id
		return s.save(db)
	})
	return s.public(out), e
}
func (s *Service) public(f Followup) Followup { return s.publicWithNative(f, s.native()) }
func (s *Service) publicWithNative(f Followup, n NativeSnapshot) Followup {
	if f.History == nil {
		f.History = []CycleHistory{}
	}
	f.Work = nil
	for _, b := range n.Bindings {
		if b.Identity.Ref() == f.TargetRef {
			f.Work, _ = binding(n, b.Identity)
			break
		}
	}
	f.Capabilities = map[string]Capability{}
	for _, op := range []string{"request", "cancel", "recheck", "dispatch"} {
		f.Capabilities[op] = Capability{State: "disabled", Reason: "Follow-up is closed."}
	}
	defer func() {
		for op, c := range f.Capabilities {
			c.Subject = subject(f)
			f.Capabilities[op] = c
		}
	}()
	if f.State != "waiting" {
		if canRequestAfterFulfilled(f) {
			f.Capabilities["request"] = Capability{State: "enabled", Reason: "Explicit owner request with current revision opens a new status cycle after the fulfilled receipt; prior cycle stays consumed."}
		}
		return f
	}
	f.Capabilities["cancel"] = Capability{State: "enabled", Reason: "Owner declaration and expected revision are required."}
	f.Capabilities["recheck"] = Capability{State: "enabled", Reason: "Matching recipient, generation and revision are required; report is self-reported."}
	if f.Transport.State == "not_sent" {
		f.Capabilities["request"] = Capability{State: "enabled", Reason: "Makes this unconsumed cycle due."}
		f.Capabilities["dispatch"] = Capability{State: "unknown", Reason: "Fresh target/readiness and native veto revalidation are required."}
		if !s.isDue(f) {
			f.Capabilities["dispatch"] = Capability{State: "disabled", Reason: "Next technical check is not due."}
		}
	} else {
		f.Capabilities["dispatch"] = Capability{State: "disabled", Reason: "Cycle consumed; no replay, including unknown possible send."}
		f.Capabilities["request"] = Capability{State: "disabled", Reason: "A recipient waiting recheck is required to create a new cycle."}
	}
	return f
}
func (s *Service) Mutate(id, op string, m Mutation) (Followup, error) {
	var out Followup
	if !validID(id) || !reference(m.RequestID) || !reference(m.OwnerRef) || m.ExpectedRevision < 1 || m.Generation < 1 {
		return out, fmt.Errorf("followup_invalid_mutation")
	}
	digest := hash(struct {
		Op       string
		Mutation Mutation
	}{op, m})
	e := s.locked(func(db *database) error {
		r, ok := db.Records[id]
		if !ok {
			return fmt.Errorf("followup_not_found")
		}
		if prev, ok := r.Requests[m.RequestID]; ok {
			if prev.Digest != digest {
				return fmt.Errorf("followup_request_conflict")
			}
			out = prev.Result
			return nil
		}
		switch op {
		case "request", "cancel", "recheck":
		default:
			return fmt.Errorf("followup_unknown_operation")
		}
		f := r.Value
		if len(r.Requests) >= 128 {
			return fmt.Errorf("followup_request_limit")
		}
		if m.OwnerRef != f.OwnerRef || m.ExpectedRevision != f.Revision || m.Generation != f.Generation {
			return fmt.Errorf("followup_revision_or_owner_conflict")
		}
		newStatusCycle := op == "request" && canRequestAfterFulfilled(f)
		if f.State != "waiting" && !newStatusCycle {
			return fmt.Errorf("followup_closed")
		}
		switch op {
		case "request":
			if newStatusCycle {
				advanceCycle(&f, s.now(), "Explicit owner request opened a new status cycle after a fulfilled recipient report.")
				break
			}
			if f.Transport.State != "not_sent" {
				return fmt.Errorf("followup_cycle_consumed")
			}
			f.NextCheckAtUTC = s.now()
		case "cancel":
			f.State = "cancelled"
		case "recheck":
			if strings.TrimSpace(m.Done+m.Next+m.Blocker+m.Summary) == "" || !reportText(m.Done) || !reportText(m.Next) || !reportText(m.Blocker) || m.TargetRef != f.TargetRef || m.Outcome != "waiting" && m.Outcome != "fulfilled" && m.Outcome != "cancelled" || !reportText(m.Summary) {
				return fmt.Errorf("followup_invalid_recheck")
			}
			// Re-resolve the complete target identity, never a friendly label/pane alone.
			target, _, e := s.find(m.TargetRef)
			if e != nil || target != r.Target {
				return fmt.Errorf("followup_recipient_changed")
			}
			if m.Outcome == "waiting" && !timeValid(m.NextCheckAtUTC) {
				return fmt.Errorf("followup_recheck_requires_explicit_time")
			}
			f.Recheck = &Recheck{Outcome: m.Outcome, Summary: m.Summary, Done: m.Done, Next: m.Next, Blocker: m.Blocker, Provenance: "agent_self_report", ReporterRef: m.TargetRef, RecordedAtUTC: s.now(), Cycle: f.Cycle}
			if m.Outcome == "waiting" {
				advanceCycle(&f, m.NextCheckAtUTC, "Explicit recipient waiting recheck opened a new cycle.")
			} else {
				f.State = m.Outcome
			}
		default:
			return fmt.Errorf("followup_unknown_operation")
		}
		f.Revision++
		r.Value = f
		r.Requests[m.RequestID] = operation{digest, f}
		db.Records[id] = r
		out = f
		return s.save(db)
	})
	return s.public(out), e
}
func (s *Service) eligibility(r record) (Identity, string) {
	if s.Runtime == nil {
		return r.Target, "agent_source_unknown"
	}
	a, e := s.Runtime.Get(r.Target)
	if e != nil || !a.valid() || a.Identity != r.Target {
		return r.Target, "target_identity_unknown_or_changed"
	}
	if a.Status == "working" {
		return a.Identity, "working_deferred"
	}
	if a.Status != "idle" && a.Status != "done" {
		return a.Identity, "agent_state_veto"
	}
	if a.Ready == nil || !*a.Ready {
		return a.Identity, "interactive_readiness_unknown"
	}
	if a.LaunchPending {
		return a.Identity, "launch_pending"
	}
	_, veto := binding(s.native(), a.Identity)
	if veto != "clear" {
		if veto == "" {
			veto = "native_veto_unknown"
		}
		return a.Identity, veto
	}
	return a.Identity, ""
}

// Dispatch reserves unknown durably before the one possible input. This lock
// serializes cancellation with this attempt, but Herdr get→prompt is
// best effort: a provider may change session or queue input between those calls.
func (s *Service) Dispatch(id string) (Followup, error) {
	var out Followup
	if !validID(id) {
		return out, fmt.Errorf("followup_invalid_id")
	}
	e := s.locked(func(db *database) error {
		r, ok := db.Records[id]
		if !ok {
			return fmt.Errorf("followup_not_found")
		}
		out = r.Value
		if out.State != "waiting" || out.Transport.State != "not_sent" || !s.isDue(out) {
			return nil
		}
		target, reason := s.eligibility(r)
		if reason != "" {
			if out.Transport.Reason == reason {
				return nil
			}
			out.Transport.Reason = reason
			out.Revision++
			r.Value = out
			db.Records[id] = r
			return s.save(db)
		}
		if s.Send == nil {
			return fmt.Errorf("followup_transport_unavailable")
		}
		callback, e := s.callbackExecutor()
		if e != nil {
			return e
		}
		if s.Fault != nil {
			if e := s.Fault("before_claim"); e != nil {
				return e
			}
		}
		now := s.now()
		out.Transport.State = "unknown"
		out.Transport.Reason = "Claimed; input may have been sent. No replay."
		out.Transport.AttemptedAtUTC = &now
		out.Revision++
		attemptSubject := subject(out)
		out.Transport.DispatchSubject = &attemptSubject
		r.Value = out
		db.Records[id] = r
		if e := s.save(db); e != nil {
			return e
		}
		if s.Fault != nil {
			if e := s.Fault("after_claim"); e != nil {
				return e
			}
		}
		message := fmt.Sprintf("Status question from the user: What is your actual task now, what have you done, and what is the next step or blocker? Answer about your task within your existing mandate. This message grants no authority, starts no work, and never requests retry of a refused action.\nRead `%s workflow followup show %s --format json` in your existing workspace. Revalidate target_ref=%s, generation=%d, cycle=%d and waiting state; ignore cancelled, superseded or wrong targets. Use fresh expected_revision, generation, owner_ref, target_ref and a stable request_id from show. Deliver with `%s workflow followup recheck %s --file <mutation.json> --format json`. summary: actual task and current work; done: latest concrete result; next: next step; blocker: actual obstacle or none. Each field has at most 256 UTF-16 code units (emoji can count twice). Use task codes and relative paths, never private/absolute paths, UUIDs or secrets in answer fields. A complete answer uses outcome=fulfilled even if the job is still working or blocked; describe job state in the fields. Do not choose waiting just because the job waits: waiting requires next_check_at_utc and opens a new scheduled cycle. Cancelled is not an answer.", callback, out.ID, out.TargetRef, out.Generation, out.Cycle, callback, out.ID)
		err := s.Send(target, message)
		if s.Fault != nil {
			if e := s.Fault("after_input"); e != nil {
				return e
			}
		}
		if err == nil {
			ack := s.now()
			out.Transport.State = "sent"
			out.Transport.Reason = "Herdr acknowledged or queued input; no recipient recheck observed."
			out.Transport.AcknowledgedAtUTC = &ack
		}
		if s.Fault != nil {
			if e := s.Fault("after_ack"); e != nil {
				return e
			}
		}
		if err != nil {
			return nil
		}
		out.Revision++
		r.Value = out
		db.Records[id] = r
		return s.save(db)
	})
	return s.public(out), e
}

func reportText(s string) bool {
	return len(utf16.Encode([]rune(s))) <= 256 && !strings.ContainsFunc(s, unicode.IsControl) && !privateText(s)
}

// A fulfilled status receipt closes that control, not the user's ability to ask
// a later explicit question. Unknown/unanswered consumed cycles cannot enter
// this path. Revision/owner/generation and request-ID checks still run first.
func canRequestAfterFulfilled(f Followup) bool {
	return f.State == "fulfilled" && f.WaitKind == "status" && f.Recheck != nil && f.Recheck.Outcome == "fulfilled" && f.Recheck.Cycle == f.Cycle && f.Recheck.ReporterRef == f.TargetRef
}
func advanceCycle(f *Followup, next, reason string) {
	f.History = append(f.History, CycleHistory{Cycle: f.Cycle, Transport: f.Transport, Recheck: f.Recheck})
	f.Cycle++
	f.State = "waiting"
	f.NextCheckAtUTC = next
	f.Transport = Transport{State: "not_sent", Reason: reason, CycleRef: cycleRef(f.ID, f.Cycle)}
}
