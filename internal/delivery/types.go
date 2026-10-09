// Package delivery owns durable delivery of a native qualified Task candidate.
// Registration and observation never publish a branch or manufacture human QA.
package delivery

import (
	"fmt"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type Registration struct {
	Kind           string                 `json:"kind"`
	SchemaVersion  int                    `json:"schema_version"`
	PublicationKey string                 `json:"publication_key"`
	TaskID         workspace.TaskID       `json:"task_id"`
	TaskResultID   workspace.TaskResultID `json:"task_result_id"`
	WorkflowRunID  string                 `json:"workflow_run_id"`
	PredecessorID  string                 `json:"predecessor_id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Body           string                 `json:"body,omitempty"`
	Metadata       *MetadataSelection     `json:"metadata,omitempty"`
}

// Pointers distinguish an absent choice from an explicitly empty list. Empty
// defaults override inheritance, but only Remove fields remove existing people.
type MetadataSelection struct {
	Assignees       *[]string `json:"assignees,omitempty" yaml:"assignees,omitempty"`
	Reviewers       *[]string `json:"reviewers,omitempty" yaml:"reviewers,omitempty"`
	RemoveAssignees []string  `json:"remove_assignees,omitempty" yaml:"remove_assignees,omitempty"`
	RemoveReviewers []string  `json:"remove_reviewers,omitempty" yaml:"remove_reviewers,omitempty"`
}

type MetadataChoice struct {
	Assignees         []string `json:"assignees"`
	Reviewers         []string `json:"reviewers"`
	RemoveAssignees   []string `json:"remove_assignees"`
	RemoveReviewers   []string `json:"remove_reviewers"`
	AssigneesSource   string   `json:"assignees_source"`
	ReviewersSource   string   `json:"reviewers_source"`
	PreferencesSHA256 string   `json:"preferences_sha256,omitempty"`
	FrozenAtUTC       string   `json:"frozen_at_utc"`
	Revision          int      `json:"revision"`
	State             string   `json:"state"`
	Error             string   `json:"error,omitempty"`
}

type Manifest struct {
	Workspace         string                      `json:"workspace"`
	TaskID            workspace.TaskID            `json:"task_id"`
	EpicID            workspace.EpicID            `json:"epic_id"`
	ProjectID         workspace.ProjectID         `json:"project_id"`
	RepoID            workspace.RepoID            `json:"repo_id"`
	Spec              workspace.TaskSpecBasis     `json:"spec"`
	TaskResult        workspace.TaskResultRecord  `json:"task_result"`
	TaskResultSHA256  string                      `json:"task_result_sha256"`
	Agreement         workspace.DeliveryAgreement `json:"agreement"`
	WorkflowRunID     string                      `json:"workflow_run_id"`
	RequestSHA256     string                      `json:"request_sha256"`
	MandateSHA256     string                      `json:"mandate_sha256"`
	PreparationID     string                      `json:"preparation_id"`
	OwnerClaim        string                      `json:"owner_claim"`
	ExpectedParentOID string                      `json:"expected_parent_oid"`
	RequiredGates     []string                    `json:"required_gates"`
	EffectKey         string                      `json:"effect_key"`
}

type Event struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	AtUTC          string `json:"at_utc"`
	Actor          string `json:"actor"`
	Origin         string `json:"origin"`
	DeliveryID     string `json:"delivery_id"`
	CandidateOID   string `json:"candidate_oid"`
	AttemptID      string `json:"attempt_id,omitempty"`
	Detail         string `json:"detail"`
	URL            string `json:"url,omitempty"`
	PreviousSHA256 string `json:"previous_sha256,omitempty"`
	SHA256         string `json:"sha256"`
}

type PRReceipt struct {
	Repository      string `json:"repository"`
	Number          int    `json:"number"`
	URL             string `json:"url"`
	HeadRef         string `json:"head_ref"`
	HeadOID         string `json:"head_oid"`
	BaseRef         string `json:"base_ref"`
	Disposition     string `json:"disposition"`
	MetadataApplied bool   `json:"metadata_applied"`
	ObservedAtUTC   string `json:"observed_at_utc"`
	CreationEventID string `json:"creation_event_id,omitempty"`
}

type LocalReceipt struct {
	TargetRef      string                                     `json:"target_ref"`
	TargetWorktree string                                     `json:"target_worktree"`
	BeforeOID      string                                     `json:"before_oid"`
	AfterOID       string                                     `json:"after_oid"`
	Integration    workspace.WorkspaceTaskIntegrationReadback `json:"integration"`
	Base           workspace.EpicBaseResult                   `json:"base"`
	Queue          workspace.WorkspaceTaskQueueReadback       `json:"queue"`
	ObservedAtUTC  string                                     `json:"observed_at_utc"`
	Closeout       *workspace.TaskCloseoutReceipt             `json:"closeout,omitempty"`
}

type Receipt struct {
	Kind               string                                `json:"kind"`
	SchemaVersion      int                                   `json:"schema_version"`
	ID                 string                                `json:"id"`
	Registration       Registration                          `json:"registration"`
	RegistrationSHA256 string                                `json:"registration_sha256"`
	Manifest           Manifest                              `json:"manifest"`
	ManifestSHA256     string                                `json:"manifest_sha256"`
	State              string                                `json:"state"`
	Reasons            []string                              `json:"reasons"`
	NextAction         string                                `json:"next_action"`
	CreatedAtUTC       string                                `json:"created_at_utc"`
	UpdatedAtUTC       string                                `json:"updated_at_utc"`
	AttemptID          string                                `json:"attempt_id,omitempty"`
	Attempts           int                                   `json:"attempts"`
	HumanQA            *workspace.TaskHumanQARecord          `json:"human_qa"`
	Acceptance         *workspace.DeliveryAcceptanceDecision `json:"acceptance,omitempty"`
	Metadata           *MetadataChoice                       `json:"metadata"`
	PR                 *PRReceipt                            `json:"pull_request"`
	Local              *LocalReceipt                         `json:"local_integration"`
	NativeClosed       bool                                  `json:"native_closed"`
	ReusedDeliveryID   string                                `json:"reused_delivery_id,omitempty"`
	Events             []Event                               `json:"events"`
}

type ListResult struct {
	Kind          string    `json:"kind"`
	SchemaVersion int       `json:"schema_version"`
	Deliveries    []Receipt `json:"deliveries"`
}

type Options struct {
	PR              PRAdapter
	PreferencesPath string
	Now             func() time.Time
	// Fault is an in-process test boundary. There is no environment bypass.
	Fault func(string) error
}

type Service struct {
	dependencies taskrun.Dependencies
	options      Options
}

func NewService(d taskrun.Dependencies, options Options) *Service {
	if options.Now == nil {
		options.Now = time.Now
	}
	if options.PR == nil {
		options.PR = &GitHubAdapter{}
	}
	return &Service{dependencies: d, options: options}
}

func (s *Service) now() string { return s.options.Now().UTC().Format(time.RFC3339Nano) }
func (s *Service) fault(point string) error {
	if s.options.Fault != nil {
		return s.options.Fault(point)
	}
	return nil
}

func Text(r Receipt) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Delivery %s: %s\nTask: %s / Epic: %s\nCandidate: %s\nMode: %s\nTarget: %s", r.ID, r.State, r.Manifest.TaskID, r.Manifest.EpicID, r.Manifest.TaskResult.ResultOID, r.Manifest.Agreement.Mode, r.Manifest.Agreement.TargetRef)
	if r.Manifest.Agreement.TargetWorktree != "" {
		fmt.Fprintf(&b, " at %s", r.Manifest.Agreement.TargetWorktree)
	}
	if r.Manifest.Agreement.GitHubRepository != "" {
		fmt.Fprintf(&b, " in %s", r.Manifest.Agreement.GitHubRepository)
	}
	b.WriteByte('\n')
	if r.Acceptance != nil {
		fmt.Fprintf(&b, "Acceptance: %s (%s)\n", r.Acceptance.Outcome, r.Acceptance.Mode)
	} else if r.Manifest.Agreement.AutomaticAcceptance() {
		fmt.Fprintf(&b, "Acceptance: automatic; responsible actor: %s\n", r.Manifest.Agreement.Acceptance.ResponsibleActor)
	}
	if r.PR != nil {
		fmt.Fprintf(&b, "PR: %s (%s; head %s; metadata applied: %t)\n", r.PR.URL, r.PR.Disposition, r.PR.HeadOID, r.PR.MetadataApplied)
	}
	if r.ReusedDeliveryID != "" {
		fmt.Fprintf(&b, "Preserved effect owner: %s\n", r.ReusedDeliveryID)
	}
	if r.Local != nil {
		fmt.Fprintf(&b, "Observed local integration: %s -> %s\n", r.Local.BeforeOID, r.Local.AfterOID)
	}
	for _, reason := range r.Reasons {
		fmt.Fprintf(&b, "Reason: %s\n", reason)
	}
	fmt.Fprintf(&b, "Native task closure: %t\nNext action: %s\n", r.NativeClosed, r.NextAction)
	return b.String()
}

func ListText(r ListResult) string {
	var b strings.Builder
	for _, v := range r.Deliveries {
		fmt.Fprintf(&b, "%s  %s  %s  %s  %s\n", v.ID, v.State, v.Manifest.TaskID, v.Manifest.Agreement.Mode, v.Manifest.Agreement.TargetRef)
	}
	if len(r.Deliveries) == 0 {
		return "No registered deliveries.\n"
	}
	return b.String()
}
