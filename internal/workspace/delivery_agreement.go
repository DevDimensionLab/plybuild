package workspace

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const (
	DeliveryPullRequest   = "pull_request"
	DeliveryLocalEpic     = "local_epic_integration"
	DeliveryLocalBranch   = "local_branch_integration"
	IntegrationOwnerHuman = "human"
)

// DeliveryAgreement v1 describes one Task in one repository. Every mode requires
// meaningful tests, review and actual human pass for the exact candidate. A PR
// stops before merge; local delivery includes the native base and queue close.
// Registering an agreement does not grant runtime authority.
type DeliveryAgreement struct {
	SchemaVersion    int       `yaml:"schema_version" json:"schema_version"`
	Mode             string    `yaml:"mode" json:"mode"`
	ProjectID        ProjectID `yaml:"project_id" json:"project_id"`
	RepoID           RepoID    `yaml:"repo_id" json:"repo_id"`
	EpicID           EpicID    `yaml:"epic_id" json:"epic_id"`
	SourceRef        string    `yaml:"source_ref,omitempty" json:"source_ref,omitempty"`
	TargetRef        string    `yaml:"target_ref" json:"target_ref"`
	TargetWorktree   string    `yaml:"target_worktree,omitempty" json:"target_worktree,omitempty"`
	GitHubRepository string    `yaml:"github_repository,omitempty" json:"github_repository,omitempty"`
	Remote           string    `yaml:"remote,omitempty" json:"remote,omitempty"`
	IntegrationOwner string    `yaml:"integration_owner,omitempty" json:"integration_owner,omitempty"`
}

// DeliveryAuthorization is supplied by the validated native run evidence. An
// input copy alone is insufficient: integration re-reads that evidence and the
// immutable result-bound Spec before allowing any new effect.
type DeliveryAuthorization struct {
	Agreement                DeliveryAgreement              `yaml:"agreement" json:"agreement"`
	AgreementSHA256          string                         `yaml:"agreement_sha256" json:"agreement_sha256"`
	RunID                    string                         `yaml:"run_id" json:"run_id"`
	CandidateRunID           string                         `yaml:"candidate_run_id" json:"candidate_run_id"`
	RequestSHA256            string                         `yaml:"request_sha256" json:"request_sha256"`
	MandateSHA256            string                         `yaml:"mandate_sha256" json:"mandate_sha256"`
	AllowedEffects           []string                       `yaml:"allowed_effects" json:"allowed_effects"`
	PermissionConfirmed      bool                           `yaml:"permission_confirmed" json:"permission_confirmed"`
	HumanIntegration         *HumanIntegrationAuthorization `yaml:"human_integration,omitempty" json:"human_integration,omitempty"`
	HumanIntegrationRequired bool                           `yaml:"human_integration_required,omitempty" json:"human_integration_required,omitempty"`
}

// HumanIntegrationAuthorization is separate from the frozen developer mandate.
// Native run evidence revalidates these preserved human decision and release bytes.
type HumanIntegrationAuthorization struct {
	SchemaVersion   int             `yaml:"schema_version" json:"schema_version"`
	PlanID          string          `yaml:"plan_id" json:"plan_id"`
	PlanSHA256      string          `yaml:"plan_sha256" json:"plan_sha256"`
	DecisionLocator string          `yaml:"decision_locator" json:"decision_locator"`
	DecisionSHA256  string          `yaml:"decision_sha256" json:"decision_sha256"`
	ReleaseLocator  string          `yaml:"release_locator" json:"release_locator"`
	ReleaseSHA256   string          `yaml:"release_sha256" json:"release_sha256"`
	HumanQARecordID HumanQARecordID `yaml:"human_qa_record_id" json:"human_qa_record_id"`
	TaskResultID    TaskResultID    `yaml:"task_result_id" json:"task_result_id"`
}

var deliveryRepositoryPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*/[A-Za-z0-9][A-Za-z0-9_.-]*$`)
var deliveryRemotePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

func ValidateDeliveryAgreement(a DeliveryAgreement) error {
	if a.SchemaVersion != 1 && a.SchemaVersion != 2 || contentSlug(string(a.ProjectID)) != nil || contentSlug(string(a.RepoID)) != nil || contentSlug(string(a.EpicID)) != nil {
		return fmt.Errorf("delivery agreement requires schema_version 1 or 2 and exact project, repository and Epic")
	}
	if a.SchemaVersion == 1 && a.IntegrationOwner != "" || a.SchemaVersion == 2 && a.IntegrationOwner != IntegrationOwnerHuman {
		return fmt.Errorf("delivery agreement v2 requires integration_owner human; v1 retains its historical owner")
	}
	if !validFullBranchRef(a.TargetRef) || a.SourceRef != "" && (!validFullBranchRef(a.SourceRef) || a.SourceRef == a.TargetRef) {
		return fmt.Errorf("delivery requires distinct full source and target branch refs")
	}
	switch a.Mode {
	case DeliveryPullRequest:
		if a.TargetWorktree != "" || !deliveryRepositoryPattern.MatchString(a.GitHubRepository) || !deliveryRemotePattern.MatchString(a.Remote) {
			return fmt.Errorf("pull_request requires an exact GitHub owner/repository and remote, without a local target worktree")
		}
	case DeliveryLocalEpic, DeliveryLocalBranch:
		if contentPath(a.TargetWorktree) != nil || a.GitHubRepository != "" || a.Remote != "" {
			return fmt.Errorf("local delivery requires an exact absolute target worktree and no remote or GitHub repository")
		}
		if a.Mode == DeliveryLocalEpic && ordinaryProductRef(a.TargetRef) {
			return fmt.Errorf("local_epic_integration cannot target main/master; explicitly choose local_branch_integration")
		}
	default:
		return fmt.Errorf("unknown delivery mode %q", a.Mode)
	}
	return nil
}

func ordinaryProductRef(ref string) bool {
	return ref == "refs/heads/main" || ref == "refs/heads/master"
}

// BindDeliveryAgreement only fills an unbound source. It never infers a mode or
// target from a branch name, current directory or repository visibility.
func BindDeliveryAgreement(a DeliveryAgreement, target QueueTarget, sourceRef string) (DeliveryAgreement, error) {
	if e := ValidateDeliveryAgreement(a); e != nil {
		return a, e
	}
	if a.ProjectID != target.ProjectID || a.RepoID != target.RepoID || a.EpicID != target.EpicID || a.SourceRef != "" && a.SourceRef != sourceRef {
		return a, fmt.Errorf("delivery agreement differs from the selected Task scope or source")
	}
	if a.Mode != DeliveryPullRequest && (a.TargetRef != target.ParentRef || a.TargetWorktree != target.ParentLocator) {
		return a, fmt.Errorf("local delivery target differs from the registered return ref/worktree")
	}
	a.SourceRef = sourceRef
	if sourceRef == "" {
		return a, fmt.Errorf("execution delivery requires the exact Task source ref")
	}
	return a, ValidateDeliveryAgreement(a)
}

func DeliveryAgreementDigest(a DeliveryAgreement) string { return queueDigest(a) }

func (a DeliveryAgreement) RequiredGates() []string {
	return []string{"meaningful_tests", "review", "exact_human_pass"}
}

func (a DeliveryAgreement) AllowedEffects() []string {
	if a.Mode == DeliveryPullRequest {
		return []string{"push", "pull_request"}
	}
	if a.HumanOwnedIntegration() {
		return nil
	}
	return []string{a.Mode}
}

func (a DeliveryAgreement) HumanOwnedIntegration() bool {
	return a.SchemaVersion == 2 && a.IntegrationOwner == IntegrationOwnerHuman
}

func (a DeliveryAgreement) StopAfter() string {
	if a.Mode == DeliveryPullRequest {
		return "verified_pull_request_before_merge"
	}
	if a.HumanOwnedIntegration() {
		return "qualified_candidate_before_human_integration"
	}
	return "observed_local_integration_and_native_closure"
}

func deliveryAgreementRule(v canonicaljson.Value) error {
	var a DeliveryAgreement
	if e := contentDecode(v, &a); e != nil {
		return e
	}
	if !contentEqual(v, contentRefValue(a)) {
		return fmt.Errorf("delivery agreement has missing or unknown fields")
	}
	return ValidateDeliveryAgreement(a)
}

// DeliveryAgreementFromSpec returns nil for historical contracts. Absence never
// authorizes PR publication or direct main/master integration.
func DeliveryAgreementFromSpec(spec canonicaljson.Value) (*DeliveryAgreement, error) {
	v, present := contentFields(spec)["delivery"]
	if !present {
		return nil, nil
	}
	if e := deliveryAgreementRule(v); e != nil {
		return nil, e
	}
	var a DeliveryAgreement
	if e := contentDecode(v, &a); e != nil {
		return nil, e
	}
	return &a, nil
}

func ValidateDeliveryAuthorization(a DeliveryAuthorization) error {
	if e := ValidateDeliveryAgreement(a.Agreement); e != nil {
		return e
	}
	if a.Agreement.SourceRef == "" || a.AgreementSHA256 != DeliveryAgreementDigest(a.Agreement) || !validTaskText(a.RunID, 1, 256) || !validTaskText(a.CandidateRunID, 1, 256) || !digestPattern.MatchString(a.RequestSHA256) || !digestPattern.MatchString(a.MandateSHA256) || !a.PermissionConfirmed {
		return fmt.Errorf("delivery authorization lacks an exact agreement, native run/request/mandate or confirmed permission")
	}
	want := a.Agreement.AllowedEffects()
	if len(a.AllowedEffects) != len(want) {
		return fmt.Errorf("delivery authorization effects differ from its mode")
	}
	for _, effect := range want {
		if !slices.Contains(a.AllowedEffects, effect) {
			return fmt.Errorf("delivery authorization lacks %s", effect)
		}
	}
	if h := a.HumanIntegration; h != nil {
		if h.SchemaVersion != 1 || !validTaskText(h.PlanID, 1, 256) || !digestPattern.MatchString(h.PlanSHA256) || contentPath(h.DecisionLocator) != nil || !digestPattern.MatchString(h.DecisionSHA256) || contentPath(h.ReleaseLocator) != nil || !digestPattern.MatchString(h.ReleaseSHA256) || !lifecycleIDPatterns["human QA"].MatchString(string(h.HumanQARecordID)) || !lifecycleIDPatterns["task result"].MatchString(string(h.TaskResultID)) {
			return fmt.Errorf("human integration authority lacks its exact plan, decision, release and candidate QA")
		}
	}
	return nil
}

func deliveryAgreementForResult(d Dependencies, root string, registry WorkItemRegistry, result TaskResultRecord) (*DeliveryAgreement, error) {
	link := taskResultSpecLink(registry, result.ID)
	if link == nil {
		return nil, nil
	}
	evaluation, e := historicalTaskSpec(d, root, registry, link.Basis)
	if e != nil {
		return nil, e
	}
	return DeliveryAgreementFromSpec(evaluation.Spec)
}

func validateIntegrationDelivery(d Dependencies, root string, r WorkItemRegistry, result TaskResultRecord, qaID HumanQARecordID, auth *DeliveryAuthorization, target QueueTarget) error {
	agreement, e := deliveryAgreementForResult(d, root, r, result)
	if e != nil {
		return e
	}
	if agreement == nil && auth == nil {
		return nil
	}
	if agreement == nil || auth == nil || ValidateDeliveryAuthorization(*auth) != nil || auth.CandidateRunID != result.RunID || !contentTypedEqual(agreement, &auth.Agreement) {
		return workError(ErrorTaskIntegrationBlocked, "delivery authorization differs from the result-bound immutable agreement and native run", nil)
	}
	if agreement.Mode == DeliveryPullRequest {
		return workError(ErrorTaskIntegrationBlocked, "pull_request authorizes no local integration", nil)
	}
	if (agreement.HumanOwnedIntegration() || auth.HumanIntegrationRequired) && auth.HumanIntegration == nil {
		return workError(ErrorTaskIntegrationBlocked, "human_integration_required: the developer mandate stops before integration; use ply integration", nil)
	}
	if h := auth.HumanIntegration; h != nil && (h.TaskResultID != result.ID || h.HumanQARecordID != qaID) {
		return workError(ErrorTaskIntegrationBlocked, "human integration decision differs from the exact candidate and QA", nil)
	}
	latestQA, e := LatestTaskHumanQA(r.HumanQARecords, result.ID, result.ResultOID, result.ResultTree)
	if e != nil || latestQA == nil || latestQA.Outcome != "pass" || latestQA.ID != qaID {
		return workError(ErrorTaskIntegrationBlocked, "delivery requires the latest unambiguous human pass for this exact candidate", e)
	}
	bound, e := BindDeliveryAgreement(*agreement, target, result.SourceRef)
	if e != nil || !contentTypedEqual(bound, *agreement) {
		return workError(ErrorTaskIntegrationBlocked, "delivery target differs from the exact registered return context", e)
	}
	evidence, e := d.HandoffEvidence.ReadTaskEvidence(TaskHandoffEvidenceRequest{ActivityID: result.ActivityID, RunID: result.RunID, HandoffID: result.HandoffID, HandoffLocator: result.HandoffLocator, HandoffSHA256: result.HandoffSHA256, StartReceiptID: result.StartReceiptID, StartReceiptLocator: result.StartReceiptLocator, StartReceiptSHA256: result.StartReceiptSHA256, TerminalResultID: result.TerminalResultID, TerminalResultLocator: result.TerminalResultLocator, TerminalResultSHA256: result.TerminalResultSHA256, InspectionSHA256: result.InspectionSHA256})
	if e != nil || !evidence.SchemaValid || !evidence.DigestValid || !evidence.LifecycleValid || !evidence.BindingValid || !evidence.CapabilityValid || !evidence.PrincipalSessionValid || !evidence.PolicyValid || !contentTypedEqual(evidence.DeliveryAuthorization, auth) {
		return workError(ErrorTaskIntegrationBlocked, "delivery authority is not confirmed by revalidated native run evidence", e)
	}
	return nil
}
