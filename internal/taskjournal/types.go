// Package taskjournal records process statements without granting operational
// authority. Native Task, run, result, QA and integration stores remain owners
// of their facts; the journal never changes them.
package taskjournal

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"time"
)

type Actor struct {
	ID        string  `json:"id"`
	Role      string  `json:"role"`
	SessionID *string `json:"session_id"`
}
type RunBinding struct {
	RunID             string `json:"run_id"`
	RequestSHA256     string `json:"request_sha256"`
	PreparationID     string `json:"preparation_id"`
	PreparationSHA256 string `json:"preparation_sha256"`
}
type Uncertainty struct {
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
}
type TimeBasis struct {
	Kind        string       `json:"kind"`
	Clock       *string      `json:"clock"`
	Precision   string       `json:"precision"`
	Uncertainty *Uncertainty `json:"uncertainty"`
}
type SourceRef struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}
type Common struct {
	Kind           string      `json:"kind"`
	SchemaVersion  int         `json:"schema_version"`
	PublicationKey string      `json:"publication_key"`
	TaskID         string      `json:"task_id"`
	Actor          Actor       `json:"actor"`
	ActivityID     *string     `json:"activity_id"`
	RunBinding     *RunBinding `json:"run_binding"`
	OccurredAt     *string     `json:"occurred_at"`
	TimeBasis      TimeBasis   `json:"time_basis"`
	Sources        []SourceRef `json:"sources"`
}
type StepRef struct {
	ID           string  `json:"id"`
	Kind         string  `json:"kind"`
	ParentStepID *string `json:"parent_step_id"`
}
type Candidate struct {
	RepoID     string `json:"repo_id"`
	WorktreeID string `json:"worktree_id"`
	OID        string `json:"oid"`
	Tree       string `json:"tree"`
}
type Relation struct {
	Type    string `json:"type"`
	EventID string `json:"event_id"`
}
type Waiting struct {
	Reason     string `json:"reason"`
	Dependency string `json:"dependency"`
	NextActor  string `json:"next_actor"`
}
type EventInput struct {
	Common
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Detail    string     `json:"detail"`
	Step      *StepRef   `json:"step"`
	Outcome   *string    `json:"outcome"`
	Candidate *Candidate `json:"candidate"`
	Relations []Relation `json:"relations"`
	Waiting   *Waiting   `json:"waiting"`
}
type Interval struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	ActorIDs []string `json:"actor_ids"`
	RunIDs   []string `json:"run_ids"`
}
type Targets struct {
	EventIDs []string  `json:"event_ids"`
	StepIDs  []string  `json:"step_ids"`
	Interval *Interval `json:"interval"`
}
type ObservationInput struct {
	Common
	Type            string  `json:"type"`
	ProposalEventID *string `json:"proposal_event_id"`
	Targets         Targets `json:"targets"`
	Finding         string  `json:"finding"`
	Hypothesis      *string `json:"hypothesis"`
	Action          *string `json:"action"`
	Owner           *string `json:"owner"`
	NextSignal      *string `json:"next_signal"`
	Decision        *string `json:"decision"`
	Assessment      *string `json:"assessment"`
}
type Record struct {
	Kind           string          `json:"kind"`
	SchemaVersion  int             `json:"schema_version"`
	WorkspaceID    string          `json:"workspace_id"`
	TaskID         string          `json:"task_id"`
	EventID        string          `json:"event_id"`
	Sequence       int             `json:"sequence"`
	PreviousSHA256 *string         `json:"previous_sha256"`
	RecordedAt     string          `json:"recorded_at"`
	InputSHA256    string          `json:"input_sha256"`
	Input          json.RawMessage `json:"input"`
}
type AppendResult struct {
	Kind                     string   `json:"kind"`
	SchemaVersion            int      `json:"schema_version"`
	TaskID                   string   `json:"task_id"`
	PublicationKey           string   `json:"publication_key"`
	EventID                  string   `json:"event_id"`
	InputSHA256              string   `json:"input_sha256"`
	RecordedAt               string   `json:"recorded_at"`
	Created                  bool     `json:"created"`
	RecordLocator            string   `json:"record_locator"`
	RecordSHA256             string   `json:"record_sha256"`
	NextTransitionAuthorized bool     `json:"next_transition_authorized"`
	Warnings                 []string `json:"-"`
}
type Source struct {
	SourceID       string  `json:"source_id"`
	Kind           string  `json:"kind"`
	Locator        string  `json:"locator"`
	SHA256         *string `json:"sha256"`
	Sequence       *int    `json:"sequence"`
	PreviousSHA256 *string `json:"previous_sha256"`
	Status         string  `json:"status"`
}
type Event struct {
	EventID            string          `json:"event_id"`
	Origin             string          `json:"origin"`
	Type               string          `json:"type"`
	Title              string          `json:"title"`
	Actor              *Actor          `json:"actor"`
	ActivityID         *string         `json:"activity_id"`
	RunID              *string         `json:"run_id"`
	StepID             *string         `json:"step_id"`
	LaneID             string          `json:"lane_id"`
	OccurredAt         *string         `json:"occurred_at"`
	OriginalOccurredAt *string         `json:"original_occurred_at"`
	RecordedAt         *string         `json:"recorded_at"`
	TimeBasis          TimeBasis       `json:"time_basis"`
	SourceIDs          []string        `json:"source_ids"`
	SourceSequence     *int            `json:"source_sequence"`
	Relations          []Relation      `json:"relations"`
	Outcome            *string         `json:"outcome"`
	Candidate          *Candidate      `json:"candidate"`
	SupersededBy       []string        `json:"superseded_by"`
	Data               json.RawMessage `json:"data"`
}
type Step struct {
	StepID          string   `json:"step_id"`
	Kind            string   `json:"kind"`
	ParentStepID    *string  `json:"parent_step_id"`
	LaneID          string   `json:"lane_id"`
	StartEventID    string   `json:"start_event_id"`
	EndEventID      *string  `json:"end_event_id"`
	State           string   `json:"state"`
	DurationSeconds *float64 `json:"duration_seconds"`
	Waiting         *Waiting `json:"waiting"`
}
type Lane struct {
	LaneID  string  `json:"lane_id"`
	ActorID *string `json:"actor_id"`
	RunID   *string `json:"run_id"`
}
type Observation struct {
	ObservationInput
	EventID string `json:"event_id"`
}
type Axis struct {
	Value       string   `json:"value"`
	EventIDs    []string `json:"event_ids"`
	BasisStatus string   `json:"basis_status"`
}
type NextAction struct {
	Text           string   `json:"text"`
	Actor          *string  `json:"actor"`
	SourceEventIDs []string `json:"source_event_ids"`
}
type Current struct {
	Candidate  *Candidate      `json:"candidate"`
	Axes       map[string]Axis `json:"axes"`
	NextAction NextAction      `json:"next_action"`
}
type Reason struct {
	Code      string   `json:"code"`
	SourceIDs []string `json:"source_ids"`
	Detail    string   `json:"detail"`
}
type Coverage struct {
	State                string   `json:"state"`
	Reasons              []Reason `json:"reasons"`
	UnknownTimeEventIDs  []string `json:"unknown_time_event_ids"`
	TimeConflictEventIDs []string `json:"time_conflict_event_ids"`
}
type Selection struct {
	Order    string   `json:"order"`
	RunID    *string  `json:"run_id"`
	ActorID  *string  `json:"actor_id"`
	EventIDs []string `json:"event_ids"`
}
type Workspace struct {
	ID           string `json:"id"`
	Root         string `json:"root"`
	MarkerSHA256 string `json:"marker_sha256"`
}
type Task struct {
	TaskID       string                         `json:"task_id"`
	Title        string                         `json:"title"`
	ProjectID    string                         `json:"project_id"`
	RepoID       string                         `json:"repo_id"`
	EpicID       string                         `json:"epic_id"`
	Worktree     *workspace.TaskWorktreeBinding `json:"worktree"`
	RecordSHA256 string                         `json:"record_sha256"`
}
type Snapshot struct {
	Kind                     string        `json:"kind"`
	SchemaVersion            int           `json:"schema_version"`
	SnapshotID               string        `json:"snapshot_id"`
	AsOf                     string        `json:"as_of"`
	Workspace                Workspace     `json:"workspace"`
	Task                     Task          `json:"task"`
	Sources                  []Source      `json:"sources"`
	Events                   []Event       `json:"events"`
	Steps                    []Step        `json:"steps"`
	Lanes                    []Lane        `json:"lanes"`
	Observations             []Observation `json:"observations"`
	Current                  Current       `json:"current"`
	Coverage                 Coverage      `json:"coverage"`
	Selection                Selection     `json:"selection"`
	NextTransitionAuthorized bool          `json:"next_transition_authorized"`
}
type Options struct{ View, Order, RunID, ActorID string }
type Service struct {
	Workspace workspace.Dependencies
	Now       func() time.Time
	fault     func(string) error
}

func New(d workspace.Dependencies) Service { return Service{Workspace: d, Now: time.Now} }
func ptr[T any](v T) *T                    { return &v }
func unknownTime() TimeBasis               { return TimeBasis{Kind: "unknown", Precision: "unknown"} }
