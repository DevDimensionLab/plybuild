// Package integration owns a human-started, plan-bound integration and Task
// closeout. Technical qualification, human judgment and observed effects remain
// separate native facts.
package integration

import (
	"fmt"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/delivery"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type Selection struct {
	DeliveryID            string `json:"delivery_id,omitempty"`
	TaskID                string `json:"task_id,omitempty"`
	Reconcile             bool   `json:"reconcile"`
	Method                string `json:"method,omitempty"`
	Keep                  bool   `json:"keep"`
	InstallProfilePath    string `json:"install_profile_path,omitempty"`
	NotificationRoutePath string `json:"notification_route_path,omitempty"`
}

// Plan contains only confirmed identities and choices. Transient observations
// and timestamps live in Preview/Receipt and cannot silently alter the plan.
type Plan struct {
	Kind               string                            `json:"kind"`
	SchemaVersion      int                               `json:"schema_version"`
	Workspace          string                            `json:"workspace"`
	Selection          Selection                         `json:"selection"`
	DeliveryID         string                            `json:"delivery_id,omitempty"`
	DeliverySHA256     string                            `json:"delivery_sha256,omitempty"`
	WorkflowRunID      string                            `json:"workflow_run_id,omitempty"`
	TaskResult         workspace.TaskResultRecord        `json:"task_result"`
	Agreement          workspace.DeliveryAgreement       `json:"agreement"`
	ExpectedParentOID  string                            `json:"expected_parent_oid"`
	Flow               string                            `json:"flow"`
	EffectKey          string                            `json:"effect_key"`
	PreviousQA         *workspace.TaskHumanQARecord      `json:"previous_qa"`
	Ownership          taskrun.HumanIntegrationOwnership `json:"ownership"`
	Closeout           workspace.TaskCloseoutPlan        `json:"closeout"`
	CloseoutSHA256     string                            `json:"closeout_sha256"`
	PR                 *PRMergeTarget                    `json:"pull_request"`
	Install            *InstallProfile                   `json:"installation"`
	InstallFile        *taskrun.FileBinding              `json:"installation_file"`
	NotificationRoute  *taskrun.FileBinding              `json:"notification_route"`
	NotificationChoice *NotificationRoute                `json:"notification_choice"`
}

type Preview struct {
	Kind          string              `json:"kind"`
	SchemaVersion int                 `json:"schema_version"`
	ID            string              `json:"id"`
	Plan          Plan                `json:"plan"`
	PlanSHA256    string              `json:"plan_sha256"`
	State         string              `json:"state"`
	Reasons       []string            `json:"reasons"`
	NextAction    string              `json:"next_action"`
	PRObservation *PRMergeObservation `json:"pull_request_observation"`
}

type Decision struct {
	Kind          string                                   `json:"kind"`
	SchemaVersion int                                      `json:"schema_version"`
	IntegrationID string                                   `json:"integration_id"`
	PlanSHA256    string                                   `json:"plan_sha256"`
	Human         workflowhandoff.DeliveryHumanAttestation `json:"human"`
}

type Event struct {
	ID             string           `json:"id"`
	Kind           string           `json:"kind"`
	IntegrationID  string           `json:"integration_id"`
	DeliveryID     string           `json:"delivery_id,omitempty"`
	TaskID         workspace.TaskID `json:"task_id"`
	CandidateOID   string           `json:"candidate_oid"`
	TargetRef      string           `json:"target_ref"`
	IntegratedOID  string           `json:"integrated_oid,omitempty"`
	PRURL          string           `json:"pr_url,omitempty"`
	AtUTC          string           `json:"at_utc"`
	PreviousSHA256 string           `json:"previous_sha256,omitempty"`
	SHA256         string           `json:"sha256"`
}

type Receipt struct {
	Kind                     string                                     `json:"kind"`
	SchemaVersion            int                                        `json:"schema_version"`
	ID                       string                                     `json:"id"`
	Plan                     Plan                                       `json:"plan"`
	PlanSHA256               string                                     `json:"plan_sha256"`
	Decision                 Decision                                   `json:"decision"`
	DecisionFile             taskrun.FileBinding                        `json:"decision_file"`
	Controller               *taskrun.FileBinding                       `json:"controller"`
	HumanAcceptance          *taskrun.HumanIntegrationAcceptance        `json:"human_acceptance"`
	State                    string                                     `json:"state"`
	Reasons                  []string                                   `json:"reasons"`
	NextAction               string                                     `json:"next_action"`
	CreatedAtUTC             string                                     `json:"created_at_utc"`
	UpdatedAtUTC             string                                     `json:"updated_at_utc"`
	IntegrationStarted       bool                                       `json:"integration_started"`
	Integrated               bool                                       `json:"integrated"`
	IntegratedOID            string                                     `json:"integrated_oid,omitempty"`
	Local                    *workflowhandoff.DeliveryIntegrationResult `json:"local_integration"`
	PR                       *PRMergeObservation                        `json:"pull_request"`
	MergeAttempts            []PRMergeObservation                       `json:"merge_attempts,omitempty"`
	MergeObservations        []PRMergeObservation                       `json:"merge_observations,omitempty"`
	MergeRetryAtUTC          string                                     `json:"merge_retry_at_utc,omitempty"`
	Reconsideration          *ReconsiderationReceipt                    `json:"reconsideration,omitempty"`
	InstallationStarted      bool                                       `json:"installation_started"`
	Installation             *InstallObservation                        `json:"installation"`
	InstallationAttempts     []InstallObservation                       `json:"installation_attempts,omitempty"`
	InstallationObservations []InstallObservation                       `json:"installation_observations,omitempty"`
	Closeout                 *workspace.TaskCloseoutReceipt             `json:"closeout"`
	NotificationStarted      bool                                       `json:"notification_started"`
	NotificationState        string                                     `json:"notification_state"`
	Notification             *NotificationObservation                   `json:"notification"`
	NotificationEvent        *NotificationEvent                         `json:"notification_event"`
	Events                   []Event                                    `json:"events"`
}

type Options struct {
	PR        PRMergeAdapter
	Installer InstallAdapter
	Notifier  NotificationAdapter
	Now       func() time.Time
	// Fault is an in-process test seam, never a command-line permission bypass.
	Fault func(string) error
}

type Service struct {
	dependencies taskrun.Dependencies
	options      Options
}

func NewService(d taskrun.Dependencies, o Options) *Service {
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.PR == nil {
		o.PR = &GitHubMergeAdapter{}
	}
	if o.Installer == nil {
		o.Installer = &LocalInstallAdapter{}
	}
	if o.Notifier == nil {
		o.Notifier = &NativeNotificationAdapter{}
	}
	return &Service{d, o}
}
func (s *Service) now() string { return s.options.Now().UTC().Format(time.RFC3339Nano) }
func (s *Service) fault(point string) error {
	if s.options.Fault != nil {
		return s.options.Fault(point)
	}
	return nil
}

type ListResult struct {
	Kind          string             `json:"kind"`
	SchemaVersion int                `json:"schema_version"`
	Deliveries    []delivery.Receipt `json:"deliveries"`
	Integrations  []Receipt          `json:"integrations"`
}

func PreviewText(p Preview) string {
	var b strings.Builder
	r := p.Plan.TaskResult
	fmt.Fprintf(&b, "Integration plan %s: %s\nTask: %s\nDelivery: %s\nCandidate: %s (tree %s)\nSource: %s at %s\nTarget: %s", p.ID, p.State, r.TaskID, p.Plan.DeliveryID, r.ResultOID, r.ResultTree, r.SourceRef, r.SourceLocator, p.Plan.Agreement.TargetRef)
	if p.Plan.Agreement.TargetWorktree != "" {
		fmt.Fprintf(&b, " at %s", p.Plan.Agreement.TargetWorktree)
	}
	fmt.Fprintf(&b, "\nFlow: %s\n", p.Plan.Flow)
	if p.Plan.PR != nil {
		fmt.Fprintf(&b, "Pull request: %s #%d; method: %s\n", p.Plan.PR.Repository, p.Plan.PR.Number, p.Plan.PR.Method)
	}
	if p.Plan.Install == nil {
		b.WriteString("Installation: none\n")
	} else {
		fmt.Fprintf(&b, "Installation: %s; artifact %s; expected %s\n", p.Plan.Install.Name, p.Plan.Install.ArtifactPath, p.Plan.Install.ArtifactSHA256)
		fmt.Fprintf(&b, "Install command: %q %q\nInstall directory: %s\nVerify command: %q %q\nVerify directory: %s\n", p.Plan.Install.Install.Executable, p.Plan.Install.Install.Arguments, p.Plan.Install.Install.CWD, p.Plan.Install.Verify.Executable, p.Plan.Install.Verify.Arguments, p.Plan.Install.Verify.CWD)
	}
	if p.Plan.Selection.Keep {
		b.WriteString("Closeout: keep the Task worktree and branch\n")
	} else {
		fmt.Fprintf(&b, "Closeout: remove only %s and local branch %s\n", r.SourceLocator, r.SourceRef)
	}
	if p.Plan.NotificationRoute == nil {
		b.WriteString("Notification: off\n")
	} else {
		fmt.Fprintf(&b, "Notification route: %s\n", p.Plan.NotificationRoute.Locator)
		if p.Plan.NotificationChoice != nil {
			fmt.Fprintf(&b, "Route name: %s; channel label: %s\n", p.Plan.NotificationChoice.Route.Name, p.Plan.NotificationChoice.Route.ChannelLabel)
		}
	}
	if p.Plan.PreviousQA != nil {
		fmt.Fprintf(&b, "Recorded candidate QA: %s (%s)\n", p.Plan.PreviousQA.Outcome, p.Plan.PreviousQA.ID)
	}
	fmt.Fprintf(&b, "Plan SHA256: %s\n", p.PlanSHA256)
	for _, reason := range p.Reasons {
		fmt.Fprintf(&b, "Reason: %s\n", reason)
	}
	fmt.Fprintf(&b, "Next: %s\n", p.NextAction)
	return b.String()
}

func Text(r Receipt) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Integration %s: %s\nTask: %s\nCandidate: %s\nObserved integration: %t", r.ID, r.State, r.Plan.TaskResult.TaskID, r.Plan.TaskResult.ResultOID, r.Integrated)
	if r.Integrated {
		fmt.Fprintf(&b, " (%s -> %s)", r.Plan.Agreement.TargetRef, r.IntegratedOID)
	}
	b.WriteByte('\n')
	for _, reason := range r.Reasons {
		fmt.Fprintf(&b, "Reason: %s\n", reason)
	}
	fmt.Fprintf(&b, "Notification: %s\nNext: %s\n", r.NotificationState, r.NextAction)
	return b.String()
}
