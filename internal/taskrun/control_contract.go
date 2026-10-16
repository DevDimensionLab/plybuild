package taskrun

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const BoundControlMode = "bound_control_v1"

type ControlContract struct {
	SchemaVersion int      `json:"schema_version"`
	Mode          string   `json:"mode"`
	Digest        string   `json:"contract_digest"`
	ParentKind    string   `json:"parent_kind"`
	WeekID        *string  `json:"week_id"`
	WeekBinding   *string  `json:"week_binding"`
	Schemas       []string `json:"schemas"`
}

func controlString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func newControlNonce() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b[:])
}

func NewControlContract(parent, week, binding string) *ControlContract {
	c := &ControlContract{SchemaVersion: 1, Mode: BoundControlMode, ParentKind: parent, WeekID: controlString(week), WeekBinding: controlString(binding), Schemas: []string{"ply.workflow.control-request@1", "ply.workflow.control-status@1"}}
	c.Digest = digest(c)
	return c
}

func ValidateControlContract(c *ControlContract) error {
	if c == nil {
		return nil
	}
	if c.SchemaVersion != 1 || c.Mode != BoundControlMode || (c.ParentKind != "standalone" && c.ParentKind != "week") {
		return ControlError("control_unsupported", "The frozen control contract is unsupported; use local recovery.")
	}
	week, binding := "", ""
	if c.ParentKind == "week" {
		if c.WeekID == nil || c.WeekBinding == nil || !regexp.MustCompile(`^week_[a-f0-9]{32}$`).MatchString(*c.WeekID) || !digestPattern.MatchString(*c.WeekBinding) {
			return ControlError("control_owner_unknown", "The parent week identity is missing or invalid.")
		}
		week, binding = *c.WeekID, *c.WeekBinding
	}
	if !equal(c, NewControlContract(c.ParentKind, week, binding)) {
		return ControlError("control_unsupported", "The frozen control contract differs; use local recovery.")
	}
	return nil
}

type ControlTarget struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}
type ControlInput struct {
	Target ControlTarget `json:"target"`
	Action string        `json:"action"`
}
type ControlValidity struct {
	IssuedAt  string `json:"issued_at"`
	ExpiresAt string `json:"expires_at"`
}
type ControlRequest struct {
	Kind          string          `json:"kind"`
	SchemaVersion int             `json:"schema_version"`
	RequestID     string          `json:"request_id"`
	Workspace     string          `json:"workspace"`
	Input         ControlInput    `json:"input"`
	Validity      ControlValidity `json:"validity"`
	Binding       string          `json:"binding"`
}
type ControlFailure struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (e *ControlFailure) Error() string      { return e.Code + ": " + e.Detail }
func ControlError(code, detail string) error { return &ControlFailure{code, detail} }
func ControlFailureOf(err error) *ControlFailure {
	if e, ok := err.(*ControlFailure); ok {
		return e
	}
	return &ControlFailure{"control_outcome_unknown", "Control evidence is unavailable or inconsistent. Inspect this same identity locally; do not retry with a new ID."}
}

var controlID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)

func ValidateControlID(id string) bool { return controlID.MatchString(id) }
func ValidateControlTarget(t ControlTarget) bool {
	return t.Kind == "run" && workflowRunIDPattern.MatchString(t.ID) || t.Kind == "week" && regexp.MustCompile(`^week_[a-f0-9]{32}$`).MatchString(t.ID)
}
func validControlAction(a string) bool { return a == "pause" || a == "stop" || a == "resume" }
func ValidateControlRequest(r ControlRequest) error {
	bad := func() error {
		return ControlError("control_invalid_request", "Expected a complete canonical control-request v1 with an exact target, action, identity, validity and binding.")
	}
	if r.Kind != "ply.workflow.control-request" || r.SchemaVersion != 1 || !ValidateControlID(r.RequestID) || !ValidateControlTarget(r.Input.Target) || !validControlAction(r.Input.Action) || !digestPattern.MatchString(r.Binding) || !filepath.IsAbs(r.Workspace) || filepath.Clean(r.Workspace) != r.Workspace {
		return bad()
	}
	if err := physical(r.Workspace, false); err != nil {
		return bad()
	}
	i, e1 := time.Parse(time.RFC3339Nano, r.Validity.IssuedAt)
	e, e2 := time.Parse(time.RFC3339Nano, r.Validity.ExpiresAt)
	if e1 != nil || e2 != nil || !e.After(i) || e.Sub(i) > 5*time.Minute {
		return bad()
	}
	return nil
}
func ParseControlRequest(raw []byte) (ControlRequest, error) {
	var r ControlRequest
	bad := ControlError("control_invalid_request", "Expected strict control-request v1 JSON; unknown, missing, duplicate or incorrectly typed fields are not accepted.")
	if len(raw) > 64<<10 {
		return r, bad
	}
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return r, bad
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(&r); err != nil {
		return r, bad
	}
	// Required strings cannot be null or omitted; structural validation covers all leaves.
	return r, ValidateControlRequest(r)
}

type ControlAction struct {
	Action       string   `json:"action"`
	Boundary     string   `json:"boundary"`
	Effect       string   `json:"effect"`
	Continuation string   `json:"continuation_policy"`
	Caveats      []string `json:"caveats"`
}
type ControlObservation struct {
	State      string `json:"state"`
	Source     string `json:"source"`
	Freshness  string `json:"freshness"`
	ObservedAt string `json:"observed_at"`
}
type ControlWeek struct {
	ID            string          `json:"id"`
	Binding       string          `json:"binding"`
	Selected      json.RawMessage `json:"selected"`
	Deadline      string          `json:"deadline"`
	DesiredAction string          `json:"desired_action"`
	Revision      uint64          `json:"control_revision"`
	Generation    int             `json:"generation"`
	FrozenControl string          `json:"frozen_control"`
	Dispatch      *string         `json:"dispatch"`
}
type ControlBasis struct {
	Kind              string             `json:"kind"`
	SchemaVersion     int                `json:"schema_version"`
	Workspace         string             `json:"workspace"`
	Scope             ControlTarget      `json:"scope"`
	ProjectID         *string            `json:"project_id"`
	RepoID            *string            `json:"repo_id"`
	EpicID            *string            `json:"epic_id"`
	TaskID            *string            `json:"task_id"`
	Goal              *string            `json:"goal"`
	RunID             *string            `json:"run_id"`
	RunRequestHash    *string            `json:"run_request_hash"`
	WorkerSession     *string            `json:"worker_session"`
	RuntimeMode       *string            `json:"runtime_mode"`
	Provider          *string            `json:"provider"`
	Contract          *ControlContract   `json:"control_contract"`
	FrozenControl     *string            `json:"frozen_control"`
	OwnerGeneration   *int               `json:"owner_generation"`
	Phase             string             `json:"phase"`
	Stage             *string            `json:"stage"`
	JobID             *string            `json:"job_id"`
	JobHash           *string            `json:"job_hash"`
	StateRevision     string             `json:"state_revision"`
	QueueRevision     *string            `json:"queue_revision"`
	NativeRevision    *string            `json:"native_revision"`
	RetryNotBefore    *string            `json:"retry_not_before"`
	OwnerLock         string             `json:"owner_lock"`
	SupervisorLock    string             `json:"supervisor_lock"`
	Week              *ControlWeek       `json:"week"`
	Observation       ControlObservation `json:"runtime_observation"`
	Actions           []ControlAction    `json:"actions"`
	UnavailableReason *string            `json:"unavailable_reason"`
	NextAction        string             `json:"next_action"`
	Error             *ControlFailure    `json:"error,omitempty"`
}

// A missing revision means no state was observed, not a digest of an invented
// state. Keep the internal string for stable comparison and emit explicit null.
func (b ControlBasis) MarshalJSON() ([]byte, error) {
	type basis ControlBasis
	return json.Marshal(struct {
		basis
		Revision *string `json:"state_revision"`
	}{basis(b), controlString(b.StateRevision)})
}

// PublicControlState projects internal phases without asserting quiescence
// while a pause/stop is still approaching its native boundary.
func PublicControlState(phase string) string {
	switch phase {
	case "running", "waiting", "awaiting_product_answer", "paused", "stopped", "completed":
		return phase
	case "pausing", "stopping":
		return "running"
	case "expired":
		return "stopped"
	case "completed_with_blocks":
		return "completed"
	default:
		return "unknown"
	}
}

func (b *ControlBasis) Unavailable(reason string) {
	b.Actions = []ControlAction{}
	b.UnavailableReason = &reason
	b.NextAction = "Inspect the preserved Run or week locally and resolve " + reason + " before requesting a new preview."
}
func (b *ControlBasis) Offer(actions ...string) {
	b.Actions = []ControlAction{}
	b.UnavailableReason = nil
	for _, a := range actions {
		effect := map[string]string{"pause": "paused", "stop": "stopped", "resume": "resumed"}[a]
		boundary := "after the active or reserved native attempt has positively completed"
		if a == "resume" {
			boundary = "when the native controller consumes this exact reservation"
		}
		b.Actions = append(b.Actions, ControlAction{a, boundary, effect, "Preserve the same Run, worker session, agreement, workset and original deadlines.", []string{"An already reserved integration may complete before parking.", "Receipt does not attest provider liveness; no signal or terminal input is sent."}})
	}
	b.NextAction = "Preview the exact action and approve the complete request before apply."
}
func NewControlBasis(root string, target ControlTarget, now time.Time) ControlBasis {
	return ControlBasis{Kind: "ply.workflow.control-capabilities", SchemaVersion: 1, Workspace: root, Scope: target, Phase: "unknown", OwnerLock: "unknown", SupervisorLock: "unknown", Actions: []ControlAction{}, Observation: ControlObservation{"unknown", "native_retained_state", "historical", now.UTC().Format(time.RFC3339Nano)}}
}
func controlBinding(r ControlRequest, b ControlBasis) string {
	r.Binding = ""
	b.Observation.ObservedAt = ""
	b.OwnerLock, b.SupervisorLock = "", ""
	return digest(struct {
		Request ControlRequest
		Basis   ControlBasis
	}{r, b})
}
func ControlRevision(v any) string { return digest(v) }
func controlSemanticState(s workflowState) string {
	var worker *WorkerStatus
	if s.Result.Worker != nil {
		copy := *s.Result.Worker
		copy.UpdatedAt = ""
		copy.OwnerState = ""
		copy.ControlRequestedAt = ""
		worker = &copy
	}
	return digest(struct {
		Request         WorkflowRequest
		Worker          *WorkerStatus
		Delivery        *DeliveryState
		Phase           string
		Wake, Successor *workerControlWake
	}{s.Request, worker, s.Result.Delivery, s.Phase, s.ControlWake, s.ControlSuccessor})
}
func ControlTimestamp(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) }
func controlNow(d Dependencies) time.Time {
	if d.Now != nil {
		return d.Now().UTC()
	}
	return time.Now().UTC()
}
func ControlSequence(r ControlRequest) string { return digest(r) }
func controlSafeError(code string) error {
	return ControlError(code, fmt.Sprintf("%s: obtain a fresh preview for this exact target or inspect its preserved identity locally.", code))
}
