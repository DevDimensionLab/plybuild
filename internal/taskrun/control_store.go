package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ControlBackend participates in parent -> TaskRun -> Project -> WorkItem order.
// Guard must not hold a process owner lock or wait for a provider/supervisor.
type ControlBackend interface {
	ControlBasis(string, string) (ControlBasis, error)
	ControlGuard(string, string, func() error) error
	ControlWake(string, string) error
}
type ControlService struct {
	Dependencies Dependencies
	Backend      ControlBackend
}
type ControlPreview struct {
	Kind          string          `json:"kind"`
	SchemaVersion int             `json:"schema_version"`
	State         string          `json:"state"`
	Basis         *ControlBasis   `json:"basis"`
	Request       *ControlRequest `json:"request"`
	Error         *ControlFailure `json:"error"`
}
type ControlRecord struct {
	Request    ControlRequest `json:"request"`
	Digest     string         `json:"request_digest"`
	Basis      ControlBasis   `json:"basis"`
	Sequence   string         `json:"sequence"`
	AcceptedAt string         `json:"accepted_at"`
}
type ControlReceipt struct {
	RequestDigest         string  `json:"request_digest"`
	Sequence              string  `json:"sequence"`
	Effect                string  `json:"control_effect"`
	ObservedAt            string  `json:"observed_at"`
	OwnerGeneration       *int    `json:"owner_generation"`
	PredecessorGeneration *int    `json:"predecessor_generation,omitempty"`
	JobID                 *string `json:"job_id"`
	JobHash               *string `json:"job_hash"`
	Reason                *string `json:"reason"`
	NativeStopSequence    *string `json:"native_stop_sequence"`
}
type ControlStatus struct {
	Kind            string             `json:"kind"`
	SchemaVersion   int                `json:"schema_version"`
	Workspace       string             `json:"workspace"`
	RequestID       string             `json:"request_id"`
	RequestDigest   *string            `json:"request_digest"`
	ApprovalBinding *string            `json:"approval_binding"`
	Action          *string            `json:"action"`
	Scope           *ControlTarget     `json:"scope"`
	Basis           *ControlBasis      `json:"basis"`
	Sequence        *string            `json:"control_sequence"`
	RequestState    string             `json:"request_state"`
	ControlEffect   string             `json:"control_effect"`
	Receipt         *ControlReceipt    `json:"receipt"`
	Observation     ControlObservation `json:"runtime_observation"`
	Error           *ControlFailure    `json:"error"`
	NextAction      string             `json:"next_action"`
}

func controlDirectory(root string) string {
	return filepath.Join(root, ".ply", "runtime-controls", "v1")
}
func controlPath(root, id, name string) string {
	return filepath.Join(controlDirectory(root), id, name+".json")
}
func readControlRecord(root, id string) (ControlRecord, error) {
	var r ControlRecord
	if !ValidateControlID(id) {
		return r, controlSafeError("control_invalid_arguments")
	}
	err := readValue(controlPath(root, id, "request"), 1<<20, &r)
	if err != nil {
		return r, err
	}
	if ValidateControlRequest(r.Request) != nil || r.Request.RequestID != id || r.Request.Workspace != root || r.Digest != digest(r.Request) || r.Sequence != r.Digest || r.Request.Binding != controlBinding(r.Request, r.Basis) {
		return r, controlSafeError("control_outcome_unknown")
	}
	return r, nil
}
func readControlReceipt(root string, r ControlRecord) (*ControlReceipt, error) {
	var p ControlReceipt
	err := readValue(controlPath(root, r.Request.RequestID, "receipt"), 64<<10, &p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if p.RequestDigest != r.Digest || p.Sequence != r.Sequence || p.ObservedAt == "" {
		return nil, controlSafeError("control_outcome_unknown")
	}
	switch p.Effect {
	case "paused", "stopped", "resumed", "completed_before_boundary", "superseded":
	default:
		return nil, controlSafeError("control_outcome_unknown")
	}
	return &p, nil
}

// ControlHistory is read-only. Every durable reservation changes the basis even
// if its publisher died before waking the consumer. No pointer-publication gap.
func ControlHistory(root string, target ControlTarget) (string, *ControlRecord, error) {
	entries, err := os.ReadDir(controlDirectory(root))
	if os.IsNotExist(err) {
		return digest([]string{}), nil, nil
	}
	if err != nil {
		return "", nil, err
	}
	history := []string{}
	var pending *ControlRecord
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".publish-") {
			continue
		}
		if !e.IsDir() || !ValidateControlID(e.Name()) {
			return "", nil, controlSafeError("control_outcome_unknown")
		}
		r, er := readControlRecord(root, e.Name())
		if er != nil {
			return "", nil, er
		}
		if r.Request.Input.Target != target {
			continue
		}
		p, er := readControlReceipt(root, r)
		if er != nil {
			return "", nil, er
		}
		history = append(history, digest(struct {
			Record  ControlRecord
			Receipt *ControlReceipt
		}{r, p}))
		if p == nil {
			if pending != nil {
				return "", nil, controlSafeError("control_outcome_unknown")
			}
			copy := r
			pending = &copy
		}
	}
	return digest(history), pending, nil
}
func WithControlStore(root string, fn func() error) error { return withStore(root, fn) }
func (c ControlService) basis(root, id string) (ControlBasis, error) {
	b, e := c.Backend.ControlBasis(root, id)
	if e != nil {
		return b, e
	}
	h, p, e := ControlHistory(root, b.Scope)
	if e != nil {
		b.Unavailable("control_outcome_unknown")
		return b, nil
	}
	b.StateRevision = digest([]string{b.StateRevision, h})
	if p != nil {
		b.Unavailable("control_in_progress")
	}
	return b, nil
}
func (c ControlService) Show(root, id string) (ControlBasis, error) { return c.basis(root, id) }
func (c ControlService) Preview(root, id, action, requestID string) (ControlPreview, error) {
	p := ControlPreview{Kind: "ply.workflow.control-preview", SchemaVersion: 1, State: "rejected"}
	fail := func(err error) (ControlPreview, error) { p.Error = ControlFailureOf(err); return p, err }
	if !validControlAction(action) || !ValidateControlID(requestID) {
		return fail(controlSafeError("control_invalid_arguments"))
	}
	b, e := c.basis(root, id)
	if e != nil {
		return fail(e)
	}
	p.Basis = &b
	available := false
	for _, a := range b.Actions {
		available = available || a.Action == action
	}
	if !available {
		return fail(controlSafeError("control_unavailable"))
	}
	now := controlNow(c.Dependencies)
	r := ControlRequest{Kind: "ply.workflow.control-request", SchemaVersion: 1, RequestID: requestID, Workspace: root, Input: ControlInput{b.Scope, action}, Validity: ControlValidity{ControlTimestamp(now), ControlTimestamp(now.Add(5 * time.Minute))}}
	r.Binding = controlBinding(r, b)
	p.Request = &r
	p.State = "ready"
	return p, nil
}
func emptyControlStatus(d Dependencies, root, id string) ControlStatus {
	return ControlStatus{Kind: "ply.workflow.control-status", SchemaVersion: 1, Workspace: root, RequestID: id, RequestState: "not_recorded", ControlEffect: "not_confirmed", Observation: ControlObservation{"unknown", "native_retained_state", "historical", ControlTimestamp(controlNow(d))}, NextAction: "Read the same request identity; absence of a record is not proof that an earlier invocation had no effect."}
}
func (c ControlService) Status(root, id string) (ControlStatus, error) {
	s := emptyControlStatus(c.Dependencies, root, id)
	r, e := readControlRecord(root, id)
	if os.IsNotExist(e) {
		if _, dirErr := os.Lstat(filepath.Dir(controlPath(root, id, "request"))); os.IsNotExist(dirErr) {
			return s, nil
		}
		e = controlSafeError("control_outcome_unknown")
	}
	if e != nil {
		s.RequestState = "outcome_unknown"
		s.ControlEffect = "unknown"
		s.Error = ControlFailureOf(e)
		return s, s.Error
	}
	s.RequestDigest = &r.Digest
	s.ApprovalBinding = &r.Request.Binding
	s.Action = &r.Request.Input.Action
	s.Scope = &r.Request.Input.Target
	s.Basis = &r.Basis
	s.Sequence = &r.Sequence
	s.RequestState = "accepted"
	s.ControlEffect = "pending_boundary"
	s.NextAction = "Observe this same request until the native boundary is confirmed. Status never wakes or recovers a controller."
	p, e := readControlReceipt(root, r)
	if e == nil && p != nil {
		s.Receipt = p
		s.ControlEffect = p.Effect
	}
	if p == nil && e == nil {
		for _, name := range []string{"consuming", "wake-unknown"} {
			if _, er := os.Lstat(controlPath(root, id, name)); er == nil {
				e = controlSafeError("control_outcome_unknown")
				break
			} else if !os.IsNotExist(er) {
				e = er
				break
			}
		}
	}
	if c.Backend != nil {
		if b, er := c.Backend.ControlBasis(root, r.Request.Input.Target.ID); er == nil {
			s.Observation = b.Observation
		}
	}
	if e != nil {
		s.RequestState = "outcome_unknown"
		s.ControlEffect = "unknown"
		s.Error = ControlFailureOf(e)
		return s, s.Error
	}
	return s, nil
}
func (c ControlService) Apply(r ControlRequest) (ControlStatus, error) {
	s := emptyControlStatus(c.Dependencies, r.Workspace, r.RequestID)
	if e := ValidateControlRequest(r); e != nil {
		s.RequestState = "rejected"
		s.ControlEffect = "none"
		s.Error = ControlFailureOf(e)
		return s, e
	}
	// Identity lookup precedes freshness, expiry and current ownership checks.
	if old, e := readControlRecord(r.Workspace, r.RequestID); e == nil {
		if !equal(old.Request, r) {
			s, _ = c.Status(r.Workspace, r.RequestID)
			s.Error = ControlFailureOf(controlSafeError("control_request_conflict"))
			return s, s.Error
		}
		return c.Status(r.Workspace, r.RequestID)
	} else if !os.IsNotExist(e) {
		s.RequestState = "outcome_unknown"
		s.ControlEffect = "unknown"
		s.Error = ControlFailureOf(e)
		return s, s.Error
	}
	reserved := false
	var reservation *ControlRecord
	e := c.Backend.ControlGuard(r.Workspace, r.Input.Target.ID, func() error {
		return withStore(r.Workspace, func() error {
			if old, er := readControlRecord(r.Workspace, r.RequestID); er == nil {
				if !equal(old.Request, r) {
					return controlSafeError("control_request_conflict")
				}
				return nil
			} else if !os.IsNotExist(er) {
				return er
			}
			// A partial identity is not a fresh request. Check under the store
			// lock so a concurrent first publisher can finish its atomic write.
			if _, er := os.Lstat(filepath.Dir(controlPath(r.Workspace, r.RequestID, "request"))); !os.IsNotExist(er) {
				return controlSafeError("control_outcome_unknown")
			}
			now := controlNow(c.Dependencies)
			issued, _ := time.Parse(time.RFC3339Nano, r.Validity.IssuedAt)
			expires, _ := time.Parse(time.RFC3339Nano, r.Validity.ExpiresAt)
			if !now.Before(expires) {
				return controlSafeError("control_expired")
			}
			if issued.After(now) {
				return controlSafeError("control_invalid_request")
			}
			b, er := c.basis(r.Workspace, r.Input.Target.ID)
			if er != nil {
				return er
			}
			if b.Scope != r.Input.Target || r.Binding != controlBinding(r, b) {
				return controlSafeError("control_stale")
			}
			allowed := false
			for _, a := range b.Actions {
				allowed = allowed || a.Action == r.Input.Action
			}
			if !allowed {
				return controlSafeError("control_unavailable")
			}
			if er = c.Dependencies.fault("control_before_reserve"); er != nil {
				return er
			}
			record := ControlRecord{r, digest(r), b, digest(r), ControlTimestamp(now)}
			publish := func() error {
				// Publication can return a directory-sync error after the link
				// committed. Preserve the attempted identity before entering I/O.
				reservation = &record
				if er := c.Dependencies.writeValue(controlPath(r.Workspace, r.RequestID, "request"), record); er != nil {
					return er
				}
				reserved = true
				return c.Dependencies.fault("control_after_reserve")
			}
			if b.NativeRevision != nil {
				guard, ok := c.Backend.(controlPublicationGuard)
				if !ok {
					return controlSafeError("control_unavailable")
				}
				return guard.ControlPublicationGuard(r.Workspace, *b.NativeRevision, publish)
			}
			return publish()
		})
	})
	if e == nil && reserved {
		e = c.Dependencies.fault("control_before_wake")
		if e == nil {
			e = c.Backend.ControlWake(r.Workspace, r.Input.Target.ID)
		}
		if e == nil {
			e = c.Dependencies.fault("control_after_wake")
		}
	}
	if e != nil {
		// A write may have committed before reporting its fsync failure. Re-read
		// the immutable identity; never retry its wake on the public replay path.
		if old, re := readControlRecord(r.Workspace, r.RequestID); re == nil {
			if !equal(old.Request, r) {
				s, _ = c.Status(r.Workspace, r.RequestID)
				s.Error = ControlFailureOf(controlSafeError("control_request_conflict"))
				return s, s.Error
			}
			_ = c.Dependencies.writeValue(controlPath(r.Workspace, r.RequestID, "wake-unknown"), map[string]string{"sequence": digest(r)})
			s, _ = c.Status(r.Workspace, r.RequestID)
			s.RequestState = "outcome_unknown"
			s.ControlEffect = "unknown"
			s.Error = ControlFailureOf(controlSafeError("control_outcome_unknown"))
			return s, s.Error
		}
		if reservation != nil {
			// A completed or ambiguous publication is stronger evidence than a
			// later failed read. Missing/corrupt storage cannot imply no effect.
			s.RequestDigest = &reservation.Digest
			s.ApprovalBinding = &reservation.Request.Binding
			s.Action = &reservation.Request.Input.Action
			s.Scope = &reservation.Request.Input.Target
			s.Basis = &reservation.Basis
			s.Sequence = &reservation.Sequence
			s.RequestState = "outcome_unknown"
			s.ControlEffect = "unknown"
			s.Error = ControlFailureOf(controlSafeError("control_outcome_unknown"))
			return s, s.Error
		}
		if _, er := os.Lstat(filepath.Dir(controlPath(r.Workspace, r.RequestID, "request"))); !os.IsNotExist(er) {
			s.RequestState = "outcome_unknown"
			s.ControlEffect = "unknown"
			s.Error = ControlFailureOf(controlSafeError("control_outcome_unknown"))
			return s, s.Error
		}
		s.RequestState = "rejected"
		s.ControlEffect = "none"
		s.Error = ControlFailureOf(e)
		return s, s.Error
	}
	return c.Status(r.Workspace, r.RequestID)
}
func BeginControlConsumption(d Dependencies, r ControlRecord) error {
	if e := d.fault("control_before_consume"); e != nil {
		return e
	}
	if e := d.writeValue(controlPath(r.Request.Workspace, r.Request.RequestID, "consuming"), map[string]string{"sequence": r.Sequence}); e != nil {
		return e
	}
	return d.fault("control_after_consume")
}
func FinishControl(d Dependencies, r ControlRecord, p ControlReceipt) error {
	if e := d.fault("control_before_receipt"); e != nil {
		return e
	}
	if e := d.writeValue(controlPath(r.Request.Workspace, r.Request.RequestID, "receipt"), p); e != nil {
		return e
	}
	return d.fault("control_after_receipt")
}
func ReceiptForControl(d Dependencies, r ControlRecord, effect string, b ControlBasis) ControlReceipt {
	return ControlReceipt{RequestDigest: r.Digest, Sequence: r.Sequence, Effect: effect, ObservedAt: ControlTimestamp(controlNow(d)), OwnerGeneration: b.OwnerGeneration, JobID: b.JobID, JobHash: b.JobHash}
}
