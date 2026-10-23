// Package agentwork projects public agent metadata and owns durable status-only
// follow-ups. It never writes Task, Run, registry, QA or delivery state.
package agentwork

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

const (
	HumanRequested      = "human_requested"
	AgentSelfRegistered = "agent_self_registered"
	FollowupKind        = "AgentFollowup@1"
	OverviewKind        = "AgentWorkReadback@1"
	ListKind            = "AgentFollowupList@1"
)

type Identity struct {
	Provider  string `json:"provider"`
	Session   string `json:"session"`
	Workspace string `json:"workspace"`
	Tab       string `json:"tab"`
	Pane      string `json:"pane"`
	Terminal  string `json:"terminal"`
}

func hash(v any) string {
	b, _ := json.Marshal(v)
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
func (i Identity) sessionKey() string { return i.Provider + ":" + i.Session }
func (i Identity) Ref() string        { return "agr_" + hash(i) }
func (i Identity) valid() bool {
	return (i.Provider == "codex" || i.Provider == "claude") && plain(i.Session) && plain(i.Workspace) && plain(i.Tab) && plain(i.Pane) && plain(i.Terminal)
}

type Observation struct {
	Identity
	Status        string
	Ready         *bool
	LaunchPending bool
}
type Runtime interface {
	List() ([]Observation, error)
	Get(Identity) (Observation, error)
}
type Work struct {
	ProjectID           string  `json:"project_id"`
	RepoID              string  `json:"repo_id"`
	EpicID              string  `json:"epic_id"`
	TaskID              string  `json:"task_id"`
	RunID               string  `json:"run_id"`
	AttemptRef          string  `json:"attempt_ref"`
	Phase               string  `json:"phase"`
	NextActor           string  `json:"next_actor"`
	NextAction          string  `json:"next_action"`
	Reason              string  `json:"reason"`
	DetailRef           string  `json:"detail_ref"`
	Source              string  `json:"source"`
	Freshness           string  `json:"freshness"`
	SourceObservedAtUTC *string `json:"source_observed_at_utc"`
}
type Binding struct {
	Identity
	Work Work
	Veto string
}
type NativeSnapshot struct {
	CoverageUnknown      bool
	VetoCoverageComplete bool
	Vetoes               map[string]string
	Available            bool
	Partial              bool
	Bindings             []Binding
}
type Registrar struct {
	Kind     string `json:"kind"`
	ActorRef string `json:"actor_ref"`
}
type Draft struct {
	PublicationKey string    `json:"publication_key"`
	OwnerRef       string    `json:"owner_ref"`
	MandateRef     string    `json:"mandate_ref"`
	Registrar      Registrar `json:"registrar"`
	TargetRef      string    `json:"target_ref"`
	WaitKind       string    `json:"wait_kind"`
	WaitRef        string    `json:"wait_ref"`
	NextCheckAtUTC string    `json:"next_check_at_utc"`
}
type Mutation struct {
	RequestID        string `json:"request_id"`
	OwnerRef         string `json:"owner_ref"`
	ExpectedRevision int    `json:"expected_revision"`
	Generation       int    `json:"generation"`
	TargetRef        string `json:"target_ref,omitempty"`
	Outcome          string `json:"outcome,omitempty"`
	NextCheckAtUTC   string `json:"next_check_at_utc,omitempty"`
	Summary          string `json:"summary,omitempty"`
	Done             string `json:"done,omitempty"`
	Next             string `json:"next,omitempty"`
	Blocker          string `json:"blocker,omitempty"`
}
type Subject struct {
	AgentRef   string `json:"agent_ref"`
	FollowupID string `json:"followup_id"`
	Revision   int    `json:"revision"`
	Generation int    `json:"generation"`
	Cycle      int    `json:"cycle"`
}

func subject(f Followup) Subject {
	return Subject{f.TargetRef, f.ID, f.Revision, f.Generation, f.Cycle}
}

type Capability struct {
	State   string  `json:"state"`
	Reason  string  `json:"reason"`
	Subject Subject `json:"subject"`
}
type Transport struct {
	DispatchSubject   *Subject `json:"dispatch_subject"`
	State             string   `json:"state"`
	Reason            string   `json:"reason"`
	CycleRef          string   `json:"cycle_ref"`
	AttemptedAtUTC    *string  `json:"attempted_at_utc"`
	AcknowledgedAtUTC *string  `json:"acknowledged_at_utc"`
}
type Recheck struct {
	Outcome       string `json:"outcome"`
	Summary       string `json:"summary"`
	Done          string `json:"done"`
	Next          string `json:"next"`
	Blocker       string `json:"blocker"`
	Provenance    string `json:"provenance"`
	ReporterRef   string `json:"reporter_ref"`
	RecordedAtUTC string `json:"recorded_at_utc"`
	Cycle         int    `json:"cycle"`
}
type CycleHistory struct {
	Cycle     int       `json:"cycle"`
	Transport Transport `json:"transport"`
	Recheck   *Recheck  `json:"recheck"`
}
type Followup struct {
	History         []CycleHistory        `json:"history"`
	Kind            string                `json:"kind"`
	SchemaVersion   int                   `json:"schema_version"`
	ID              string                `json:"id"`
	Revision        int                   `json:"revision"`
	Generation      int                   `json:"generation"`
	OwnerRef        string                `json:"owner_ref"`
	MandateRef      string                `json:"mandate_ref"`
	Registrar       Registrar             `json:"registrar"`
	TargetRef       string                `json:"target_ref"`
	WaitKind        string                `json:"wait_kind"`
	WaitRef         string                `json:"wait_ref"`
	Work            *Work                 `json:"work"`
	NextCheckAtUTC  string                `json:"next_check_at_utc"`
	State           string                `json:"state"`
	Cycle           int                   `json:"cycle"`
	Transport       Transport             `json:"transport"`
	Recheck         *Recheck              `json:"recheck"`
	RegisteredAtUTC string                `json:"registered_at_utc"`
	Capabilities    map[string]Capability `json:"capabilities"`
}
type Source struct {
	Source              string  `json:"source"`
	Status              string  `json:"status"`
	Freshness           string  `json:"freshness"`
	SourceObservedAtUTC *string `json:"source_observed_at_utc"`
	Coverage            string  `json:"coverage"`
}
type Agent struct {
	AgentRef      string                `json:"agent_ref"`
	Provider      string                `json:"provider"`
	Status        string                `json:"status"`
	Pane          string                `json:"pane"`
	FollowupRefs  []Subject             `json:"followup_refs"`
	Workspace     string                `json:"workspace"`
	Tab           string                `json:"tab"`
	BindingStatus string                `json:"binding_status"`
	Work          *Work                 `json:"work"`
	Capabilities  map[string]Capability `json:"capabilities"`
}
type Overview struct {
	Kind          string     `json:"kind"`
	SchemaVersion int        `json:"schema_version"`
	ObservedAtUTC string     `json:"observed_at_utc"`
	Sources       []Source   `json:"sources"`
	Agents        []Agent    `json:"agents"`
	Followups     []Followup `json:"followups"`
}
type List struct {
	Kind          string     `json:"kind"`
	SchemaVersion int        `json:"schema_version"`
	ObservedAtUTC string     `json:"observed_at_utc"`
	Items         []Followup `json:"items"`
}
type Service struct {
	Root    string
	Now     func() time.Time
	Runtime Runtime
	Native  func() NativeSnapshot
	Send    func(Identity, string) error
	Fault   func(string) error
}

func (s *Service) now() string { return s.Now().UTC().Format(time.RFC3339Nano) }
