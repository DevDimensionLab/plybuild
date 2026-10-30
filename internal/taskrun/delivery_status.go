package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryStatus is a read-only projection, never a receipt or an authority
// input. Each outcome remains bound to its own candidate and evidence.
type DeliveryStatus struct {
	ScopedContinuation *DeliveryDispositionStatus        `json:"scoped_continuation,omitempty"`
	ResourceRetirement *DeliveryResourceRetirementStatus `json:"resource_retirement"`
	Startup            *DeliveryStartupStatus            `json:"startup,omitempty"`
	Source             DeliverySourceStatus              `json:"source"`
	Verification       *DeliveryVerificationStatus       `json:"verification"`
	QualifiedCandidate *DeliveryQualifiedStatus          `json:"qualified_candidate"`
	HumanJudgment      *DeliveryHumanStatus              `json:"human_judgment"`
	Acceptance         *DeliveryAcceptanceStatus         `json:"acceptance,omitempty"`
	FinalDelivery      DeliveryFinalStatus               `json:"final_delivery"`
	Reasons            []Reason                          `json:"reasons"`
}

// These are separate observations, not a claim that a launched terminal has
// accepted permission or begun meaningful product work.
type DeliveryStartupStatus struct {
	Terminal          string       `json:"terminal"`
	NativeSession     string       `json:"native_session"`
	TaskPrompt        string       `json:"task_prompt"`
	TaskReceived      string       `json:"task_received"`
	RuntimePermission string       `json:"runtime_permission"`
	TargetWork        string       `json:"target_work"`
	Stop              *FileBinding `json:"stop"`
	StartupProtocol   string       `json:"startup_protocol,omitempty"`
	ObservationSource string       `json:"observation_source,omitempty"`
	SandboxCoverage   string       `json:"sandbox_coverage,omitempty"`
	HumanScope        *FileBinding `json:"human_scope,omitempty"`
	RuntimeDecision   string       `json:"runtime_decision,omitempty"`
}

type DeliverySourceStatus struct {
	OID   string `json:"oid"`
	Tree  string `json:"tree"`
	Clean bool   `json:"clean"`
	State string `json:"state"`
}

type DeliveryVerificationStatus struct {
	AttemptID     string       `json:"attempt_id"`
	CandidateOID  string       `json:"candidate_oid"`
	CandidateTree string       `json:"candidate_tree"`
	Outcome       string       `json:"outcome"`
	Exit          *int         `json:"exit"`
	Qualification string       `json:"qualification"`
	Current       bool         `json:"current"`
	InputsMatch   bool         `json:"inputs_match"`
	Receipt       *FileBinding `json:"receipt"`
}

type DeliveryQualifiedStatus struct {
	Key          string                 `json:"key"`
	OID          string                 `json:"oid"`
	Tree         string                 `json:"tree"`
	TaskResultID workspace.TaskResultID `json:"task_result_id"`
	Current      bool                   `json:"current"`
}

type DeliveryHumanStatus struct {
	CandidateOID string                    `json:"candidate_oid"`
	TaskResultID workspace.TaskResultID    `json:"task_result_id"`
	HumanQAID    workspace.HumanQARecordID `json:"human_qa_id"`
	Outcome      string                    `json:"outcome"`
	Current      bool                      `json:"current"`
}

type DeliveryAcceptanceStatus struct {
	workspace.DeliveryAcceptanceDecision
	Current          bool   `json:"current"`
	ResponsibleActor string `json:"responsible_actor,omitempty"`
}

type DeliveryFinalStatus struct {
	State        string       `json:"state"`
	CandidateOID string       `json:"candidate_oid,omitempty"`
	Evidence     *FileBinding `json:"evidence"`
}

func deliveryStatus(d Dependencies, s workflowState) *DeliveryStatus {
	if !deliveryRun(s.Request) || s.Result.Delivery == nil {
		return nil
	}
	out := &DeliveryStatus{Source: DeliverySourceStatus{State: "unknown"}, FinalDelivery: DeliveryFinalStatus{State: "not_delivered"}, Reasons: []Reason{}}
	problem := func(code string, err error) { out.Reasons = append(out.Reasons, Reason{code, err.Error()}) }
	out.Startup = &DeliveryStartupStatus{Terminal: "not_observed", NativeSession: s.Result.Transport.AgentSessionID, TaskPrompt: s.Phase, TaskReceived: "unknown", RuntimePermission: s.Result.Delivery.PermissionState, TargetWork: "unknown", Stop: s.Result.Delivery.ActiveStop}
	if s.Result.Delivery.ActiveStop != nil {
		out.Startup.RuntimeDecision = "preserved_stop"
	}
	if s.Result.Transport.PaneID != "" {
		out.Startup.Terminal = "created"
	}
	if s.Acceptance != nil {
		out.Startup.TaskReceived = "owner_callback_observed"
	}
	source, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if err != nil {
		problem("source_unavailable", err)
	} else if source.Ref == s.Observed.Target.Ref && source.GitCommonDir == s.Observed.Target.GitCommonDir {
		out.Source = DeliverySourceStatus{OID: source.OID, Tree: source.Tree, Clean: source.Clean && len(source.InProgress) == 0, State: "observed"}
		out.Startup.TargetWork = "no_changes_observed"
		if !out.Source.Clean || source.OID != s.Observed.Target.OID || source.Tree != s.Observed.Target.Tree {
			out.Startup.TargetWork = "changes_observed"
		}
	}
	observe := func(receipt deliveryVerificationReceipt, binding FileBinding, state string) *DeliveryVerificationStatus {
		v := &DeliveryVerificationStatus{AttemptID: receipt.AttemptID, CandidateOID: receipt.CandidateOID, CandidateTree: receipt.CandidateTree, Exit: receipt.Exit, Outcome: "unknown", Qualification: state, Receipt: &binding}
		if err := deliveryPreservedVerificationInputs(receipt); err != nil {
			problem("verification_evidence_changed", err)
		} else if receipt.Exit != nil && receipt.Error == "" {
			v.Outcome = "failed"
			if *receipt.Exit == 0 {
				v.Outcome = "passed"
			}
		}
		if err := deliveryVerificationInputs(receipt); err != nil {
			problem("verification_input_changed", err)
		} else {
			v.InputsMatch = true
		}
		v.Current = v.Outcome != "unknown" && out.Source.State == "observed" && out.Source.Clean && out.Source.OID == receipt.CandidateOID && out.Source.Tree == receipt.CandidateTree && v.InputsMatch
		return v
	}
	for i := len(s.Result.Delivery.Events) - 1; i >= 0; i-- {
		event := s.Result.Delivery.Events[i]
		if event.Kind != "verification" {
			continue
		}
		raw, err := workflowBound(event.Binding, 256<<10)
		var receipt deliveryVerificationReceipt
		if err == nil {
			err = decode(raw, 256<<10, &receipt)
		}
		if err != nil {
			problem("verification_unavailable", err)
			break
		}
		out.Verification = observe(receipt, event.Binding, "recorded")
		break
	}
	if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.Kind == "verification" && attempt.State == "attempted" {
		if out.Verification == nil || out.Verification.AttemptID != attempt.ID {
			out.Verification = &DeliveryVerificationStatus{AttemptID: attempt.ID, CandidateOID: attempt.CandidateOID, CandidateTree: attempt.CandidateTree, Outcome: "unknown", Qualification: "unrecorded"}
			if receipt, binding, err := deliveryRecoverVerification(s); err == nil {
				out.Verification = observe(receipt, binding, "pending")
			}
		}
	}
	if n := len(s.Result.Delivery.Candidates); n > 0 {
		c := s.Result.Delivery.Candidates[n-1]
		current := out.Verification != nil && out.Verification.AttemptID == c.Key && out.Verification.Current
		if err := workspace.ValidateTaskVerificationReceipt(d.Workspace, c.TaskResult); err != nil {
			problem("verification_record_changed", err)
			current = false
		}
		out.QualifiedCandidate = &DeliveryQualifiedStatus{Key: c.Key, OID: c.OID, Tree: c.Tree, TaskResultID: c.TaskResult.ID, Current: current}
	}
	if resource, err := deliveryResourceStatus(d, s); err == nil {
		out.ResourceRetirement = resource
	} else {
		problem("resource_observation_unknown", err)
	}
	return out
}
