// Package humaninput lets an ordinary cooperative agent publish one necessary
// human question and read back one immutable answer. It owns no process,
// provider, Task or transport authority; it only records local data.
package humaninput

import (
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Schema 1 of cooperative_agent_v1. Nullable fields are always present as
// explicit JSON null; public objects deliberately have no omitempty.
const (
	Mode          = "cooperative_agent_v1"
	SchemaVersion = 1

	DraftKind    = "ply.workflow.human-input-question-draft"
	QuestionKind = "ply.workflow.human-input-question"
	ListKind     = "ply.workflow.human-input-list"
	PreviewKind  = "ply.workflow.human-input-answer-preview"
	RequestKind  = "ply.workflow.human-input-answer-request"
	StatusKind   = "ply.workflow.human-input-answer-status"
	WaitKind     = "ply.workflow.human-input-wait"
	ErrorKind    = "ply.workflow.human-input-error"

	PurposeAnswer   = "human_input_answer_v1"
	ProvenanceAgent = "reported_by_agent"

	AnswerChoice = "choice"
	AnswerText   = "text"

	StateWaiting   = "waiting_for_answer"
	StateAnswered  = "answered"
	StateReceived  = "agent_received"
	StateContinued = "agent_reported_continued"
	StateWithdrawn = "withdrawn"

	ActorHuman = "human"
	ActorAgent = "agent"
	ActorNone  = "none"

	StatusRecorded    = "recorded"
	StatusOtherWinner = "other_winner"
	StatusNotRecorded = "not_recorded"
	StatusUnavailable = "unavailable"
	StatusUnknown     = "unknown"

	WaitAnswered     = "answered"
	WaitWithdrawn    = "withdrawn"
	WaitStillWaiting = "still_waiting"

	DefaultListLimit = 20
	MaxListLimit     = 20
	MaxOutputBytes   = 4 << 20
	MaxWaitSeconds   = 600
)

type Option struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Consequence string `json:"consequence"`
}

// QuestionDraft is what the agent writes for ask. The workspace comes from the
// explicit --workspace flag; provider and session are agent-reported facts.
type QuestionDraft struct {
	Kind           string   `json:"kind"`
	SchemaVersion  int      `json:"schema_version"`
	PublicationKey string   `json:"publication_key"`
	ProjectID      string   `json:"project_id"`
	RepoID         *string  `json:"repo_id"`
	TaskID         *string  `json:"task_id"`
	ActivityID     string   `json:"activity_id"`
	ActivityTitle  string   `json:"activity_title"`
	Provider       string   `json:"provider"`
	SessionID      *string  `json:"session_id"`
	OriginCWD      string   `json:"origin_cwd"`
	Text           string   `json:"text"`
	Reason         string   `json:"reason"`
	AnswerKind     string   `json:"answer_kind"`
	Options        []Option `json:"options"`
}

type Scope struct {
	ProjectID     string  `json:"project_id"`
	RepoID        *string `json:"repo_id"`
	TaskID        *string `json:"task_id"`
	ActivityID    string  `json:"activity_id"`
	ActivityTitle string  `json:"activity_title"`
}

type Source struct {
	Provider   string  `json:"provider"`
	SessionID  *string `json:"session_id"`
	OriginCWD  string  `json:"origin_cwd"`
	Provenance string  `json:"provenance"`
}

type Answer struct {
	RequestID     string  `json:"request_id"`
	RequestSHA256 string  `json:"request_sha256"`
	AnswerKind    string  `json:"answer_kind"`
	OptionID      *string `json:"option_id"`
	Text          *string `json:"text"`
	RecordedAtUTC string  `json:"recorded_at_utc"`
}

type Receipt struct {
	RequestID     string `json:"request_id"`
	RequestSHA256 string `json:"request_sha256"`
	ReportedAtUTC string `json:"reported_at_utc"`
	Provenance    string `json:"provenance"`
}

type Continuation struct {
	RequestID     string `json:"request_id"`
	RequestSHA256 string `json:"request_sha256"`
	Summary       string `json:"summary"`
	ReportedAtUTC string `json:"reported_at_utc"`
	Provenance    string `json:"provenance"`
}

type Withdrawal struct {
	Reason         string `json:"reason"`
	WithdrawnAtUTC string `json:"withdrawn_at_utc"`
	Provenance     string `json:"provenance"`
}

type NextAction struct {
	Actor   string `json:"actor"`
	Message string `json:"message"`
}

// Question is the one public row used by ask, list, show, wait, acknowledge,
// continue and withdraw. question_sha256 covers only the immutable content.
type Question struct {
	Kind           string        `json:"kind"`
	SchemaVersion  int           `json:"schema_version"`
	Mode           string        `json:"mode"`
	Workspace      string        `json:"workspace"`
	QuestionID     string        `json:"question_id"`
	QuestionSHA256 string        `json:"question_sha256"`
	PublicationKey string        `json:"publication_key"`
	Scope          Scope         `json:"scope"`
	Source         Source        `json:"source"`
	Text           string        `json:"text"`
	Reason         string        `json:"reason"`
	AnswerKind     string        `json:"answer_kind"`
	Options        []Option      `json:"options"`
	CreatedAtUTC   string        `json:"created_at_utc"`
	State          string        `json:"state"`
	Answer         *Answer       `json:"answer"`
	Receipt        *Receipt      `json:"receipt"`
	Continuation   *Continuation `json:"continuation"`
	Withdrawal     *Withdrawal   `json:"withdrawal"`
	NextAction     NextAction    `json:"next_action"`
}

type ListResult struct {
	Kind          string     `json:"kind"`
	SchemaVersion int        `json:"schema_version"`
	Mode          string     `json:"mode"`
	Workspace     string     `json:"workspace"`
	ProjectID     string     `json:"project_id"`
	State         *string    `json:"state"`
	Limit         int        `json:"limit"`
	Cursor        *string    `json:"cursor"`
	Questions     []Question `json:"questions"`
	NextCursor    *string    `json:"next_cursor"`
	HasMore       bool       `json:"has_more"`
	Complete      bool       `json:"complete"`
	ObservedAtUTC string     `json:"observed_at_utc"`
}

type ListQuery struct {
	ProjectID string
	State     *string
	Limit     int // zero selects DefaultListLimit
	Cursor    *string
}

// AnswerValue carries exactly one of option_id (choice) or text (text).
type AnswerValue struct {
	AnswerKind string  `json:"answer_kind"`
	OptionID   *string `json:"option_id"`
	Text       *string `json:"text"`
}

// AnswerRequest is immutable after preview. binding_sha256 is the RFC8785
// digest of the request without binding_sha256; request_sha256 is the digest
// of the complete canonical request.
type AnswerRequest struct {
	Kind           string      `json:"kind"`
	SchemaVersion  int         `json:"schema_version"`
	Mode           string      `json:"mode"`
	Purpose        string      `json:"purpose"`
	Workspace      string      `json:"workspace"`
	QuestionID     string      `json:"question_id"`
	QuestionSHA256 string      `json:"question_sha256"`
	Scope          Scope       `json:"scope"`
	RequestID      string      `json:"request_id"`
	Answer         AnswerValue `json:"answer"`
	BindingSHA256  string      `json:"binding_sha256"`
}

type Preview struct {
	Kind          string        `json:"kind"`
	SchemaVersion int           `json:"schema_version"`
	Mode          string        `json:"mode"`
	Workspace     string        `json:"workspace"`
	Question      Question      `json:"question"`
	Request       AnswerRequest `json:"request"`
	RequestSHA256 string        `json:"request_sha256"`
	ObservedAtUTC string        `json:"observed_at_utc"`
}

type Reason struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type WinnerRef struct {
	RequestID     string `json:"request_id"`
	RequestSHA256 string `json:"request_sha256"`
}

// AnswerStatus is the authoritative outcome of one (question_id, request_id)
// pair, returned by apply and status. unknown and unavailable open no write.
type AnswerStatus struct {
	Kind          string     `json:"kind"`
	SchemaVersion int        `json:"schema_version"`
	Mode          string     `json:"mode"`
	Workspace     string     `json:"workspace"`
	QuestionID    string     `json:"question_id"`
	RequestID     string     `json:"request_id"`
	State         string     `json:"state"`
	RequestSHA256 *string    `json:"request_sha256"`
	Winner        *WinnerRef `json:"winner"`
	Question      *Question  `json:"question"`
	Reasons       []Reason   `json:"reasons"`
	ObservedAtUTC string     `json:"observed_at_utc"`
}

type WaitResult struct {
	Kind           string   `json:"kind"`
	SchemaVersion  int      `json:"schema_version"`
	Mode           string   `json:"mode"`
	Workspace      string   `json:"workspace"`
	QuestionID     string   `json:"question_id"`
	Outcome        string   `json:"outcome"`
	TimeoutSeconds int      `json:"timeout_seconds"`
	Question       Question `json:"question"`
	ObservedAtUTC  string   `json:"observed_at_utc"`
}

// ErrorResult is printed for a semantic refusal. question is the current
// question when it is known, so a refusal never invents a next actor.
type ErrorResult struct {
	Kind          string    `json:"kind"`
	SchemaVersion int       `json:"schema_version"`
	Mode          string    `json:"mode"`
	Workspace     *string   `json:"workspace"`
	QuestionID    *string   `json:"question_id"`
	Code          string    `json:"code"`
	Detail        string    `json:"detail"`
	Question      *Question `json:"question"`
	ObservedAtUTC string    `json:"observed_at_utc"`
}

type Dependencies struct {
	Workspace workspace.Dependencies
	Now       func() time.Time
	// Fault is an in-process test seam around answer publication.
	Fault func(string) error
	// Poll is the read-only wait interval; zero selects 500ms.
	Poll time.Duration
}

func SystemDependencies() Dependencies {
	return Dependencies{Workspace: workspace.SystemDependencies(), Now: time.Now}
}

func (d Dependencies) now() string {
	now := time.Now
	if d.Now != nil {
		now = d.Now
	}
	return now().UTC().Format(time.RFC3339Nano)
}
func (d Dependencies) fault(point string) error {
	if d.Fault != nil {
		return d.Fault(point)
	}
	return nil
}
func ptr[T any](v T) *T { return &v }
