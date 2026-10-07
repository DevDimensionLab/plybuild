// Package workflowtrace projects recorded Task history. It never grants a
// transition, probes a process, or estimates effort from elapsed time.
package workflowtrace

import (
	"encoding/json"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const Kind = "WorkflowTraceReadback@1"

type Result struct {
	Kind                     string       `json:"kind"`
	SchemaVersion            int          `json:"schema_version"`
	Workspace                string       `json:"workspace"`
	Task                     Task         `json:"task"`
	ObservedAtUTC            string       `json:"observed_at_utc"`
	Freshness                string       `json:"freshness"`
	LiveState                string       `json:"live_state"`
	HistoryState             string       `json:"history_state"`
	Ordering                 string       `json:"ordering"`
	CurrentGoals             []Goal       `json:"current_goals"`
	Runs                     []Run        `json:"runs"`
	Sources                  []Source     `json:"sources"`
	Events                   []Event      `json:"events"`
	Chains                   []Chain      `json:"chains"`
	Relations                []Relation   `json:"relations"`
	Candidates               []Candidate  `json:"candidates"`
	Analysis                 []Analysis   `json:"analysis"`
	Coverage                 Coverage     `json:"coverage"`
	Diagnostics              []Diagnostic `json:"diagnostics"`
	NextTransitionAuthorized bool         `json:"next_transition_authorized"`
}

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ProjectID string `json:"project_id"`
	RepoID    string `json:"repo_id"`
	EpicID    string `json:"epic_id"`
}

type Goal struct {
	SpecID         string          `json:"spec_id"`
	Revision       int             `json:"revision"`
	ManifestSHA256 string          `json:"manifest_sha256"`
	Status         string          `json:"status"`
	SourceIDs      []string        `json:"source_ids"`
	Declaration    json.RawMessage `json:"declaration"`
}

type Run struct {
	ID              string                   `json:"id"`
	Family          string                   `json:"family"`
	Provider        *string                  `json:"provider"`
	SessionID       *string                  `json:"session_id"`
	RequestSHA256   string                   `json:"request_sha256"`
	FrozenBasis     *workspace.TaskSpecBasis `json:"frozen_basis"`
	FrozenGoal      json.RawMessage          `json:"frozen_goal"`
	DeclaredProcess json.RawMessage          `json:"declared_process"`
	Coverage        string                   `json:"coverage"`
	SourceIDs       []string                 `json:"source_ids"`
	EventIDs        []string                 `json:"event_ids"`
}

type Source struct {
	ID      string  `json:"id"`
	Kind    string  `json:"kind"`
	Locator string  `json:"locator"`
	SHA256  *string `json:"sha256"`
	Status  string  `json:"status"`
}

type Event struct {
	ID              string                `json:"id"`
	NativeID        *string               `json:"native_id"`
	Type            string                `json:"type"`
	Title           string                `json:"title"`
	Role            string                `json:"role"`
	ActorClaim      *string               `json:"actor_claim"`
	Recorder        string                `json:"recorder"`
	EvidenceClass   string                `json:"evidence_class"`
	RunIDs          []string              `json:"run_ids"`
	CandidateID     *string               `json:"candidate_id"`
	ResultID        *string               `json:"result_id"`
	Outcome         *string               `json:"outcome"`
	OccurredAtUTC   *string               `json:"occurred_at_utc"`
	ReportedAtUTC   *string               `json:"reported_at_utc"`
	RegisteredAtUTC *string               `json:"registered_at_utc"`
	TimeBasis       taskjournal.TimeBasis `json:"time_basis"`
	SourceIDs       []string              `json:"source_ids"`
	Positions       []Position            `json:"positions"`
	Data            json.RawMessage       `json:"data"`
}

type Position struct {
	ChainID  string `json:"chain_id"`
	Sequence int    `json:"sequence"`
}

type Chain struct {
	ID       string   `json:"id"`
	Ordering string   `json:"ordering"`
	EventIDs []string `json:"event_ids"`
}

type Relation struct {
	Type      string   `json:"type"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Basis     string   `json:"basis"`
	SourceIDs []string `json:"source_ids"`
}

type Candidate struct {
	ID        string          `json:"id"`
	OID       string          `json:"oid"`
	Tree      string          `json:"tree"`
	RunIDs    []string        `json:"run_ids"`
	ResultIDs []string        `json:"result_ids"`
	EventIDs  []string        `json:"event_ids"`
	Axes      map[string]Axis `json:"axes"`
}

type Axis struct {
	State       string   `json:"state"`
	Outcomes    []string `json:"outcomes"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// EvidenceIDs always refer to Events; SourceIDs to Sources. Value is null for
// unmeasurable intervals and questions, never a manufactured zero.
type Analysis struct {
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	Definition  string   `json:"definition"`
	Value       *float64 `json:"value"`
	Unit        string   `json:"unit"`
	Coverage    string   `json:"coverage"`
	EvidenceIDs []string `json:"evidence_ids"`
	SourceIDs   []string `json:"source_ids"`
	Question    *string  `json:"question"`
	Unknowns    []string `json:"unknowns"`
}

type Coverage struct {
	State    string   `json:"state"`
	Read     []string `json:"read"`
	Unknowns []string `json:"unknowns"`
}

type Diagnostic struct {
	Code      string   `json:"code"`
	Detail    string   `json:"detail"`
	SourceIDs []string `json:"source_ids"`
	EventIDs  []string `json:"event_ids"`
}
