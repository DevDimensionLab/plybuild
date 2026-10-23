package taskrun

import (
	"encoding/json"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// parallelHumanTerminal observes the historical human route; it never invokes
// integration checks, native owner authority, provider observations or effects.
// Empty proof means incomplete history, while contradictory bindings are errors.
func parallelHumanTerminal(d Dependencies, s workflowState, c DeliveryCandidate, closeout workspace.TaskCloseoutReceipt) (string, error) {
	human, err := preservedHumanIntegration(s)
	if err != nil || human == nil || c.Integration == nil {
		return "", err
	}
	raw, err := workflowBound(*c.Integration, 4<<20)
	if err != nil {
		return "", err
	}
	var v struct {
		Kind          string                                    `json:"kind"`
		RunID         string                                    `json:"run_id"`
		RequestSHA256 string                                    `json:"request_sha256"`
		Candidate     string                                    `json:"candidate"`
		Result        workflowhandoff.DeliveryIntegrationResult `json:"result"`
		Error         string                                    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	invalid := func() (string, error) {
		return "", workflowError(4, "activation_consumer_compatibility_unknown: human terminal integration differs from its exact native candidate/authority")
	}
	found := false
	for _, event := range s.Result.Delivery.Events {
		if event.Kind == "integration" && equal(event.Binding, *c.Integration) {
			found = true
		}
	}
	rb := v.Result.Integration.Readback
	if !found || v.Kind != "PlyDeliveryIntegration@1" || v.RunID != s.Result.RunID || v.RequestSHA256 != s.Result.RequestSHA256 || v.Candidate != c.Key || v.Error != "" || !v.Result.Completed || rb.RecoveryStatus != "complete" || rb.ParentOID != c.OID && !rb.HistoricalCompletion || rb.TaskResult == nil || !equal(*rb.TaskResult, c.TaskResult) {
		return invalid()
	}
	wire, err := json.Marshal(rb.Value)
	if err != nil {
		return "", err
	}
	var ir workspace.IntegrationResult
	if err = json.Unmarshal(traceWireJSON(traceWireObject(wire)["integration_result"]), &ir); err != nil {
		return invalid()
	}
	registry, err := d.Workspace.WorkItems.Snapshot(s.Request.WorkspaceRoot)
	if err != nil {
		return "", err
	}
	registered := false
	for _, actual := range registry.IntegrationResults {
		if equal(actual, ir) {
			registered = true
		}
	}
	if !registered || ir.ID == "" || ir.RecoveryStatus != "complete" || (ir.Outcome != "exact_effect" && ir.Outcome != "already_integrated") || closeout.NativeIntegrationID != ir.ID {
		return invalid()
	}
	var authority *workspace.IntegrationAuthority
	for i := range registry.IntegrationAuthorities {
		a := &registry.IntegrationAuthorities[i]
		if a.ID == ir.AuthorityID {
			authority = a
		}
	}
	if authority == nil {
		return invalid()
	}
	a := authority
	p := a.Plan
	authorization := p.DeliveryAuthorization
	agreement := deliveryEffectiveAgreement(s)
	if agreement == nil || authorization == nil || workspace.ValidateDeliveryAuthorization(*authorization) != nil || !equal(authorization.HumanIntegration, human) || !equal(&authorization.Agreement, agreement) || authorization.RunID != s.Result.RunID || authorization.CandidateRunID != c.TaskResult.RunID || authorization.RequestSHA256 != s.Result.RequestSHA256 || authorization.MandateSHA256 != s.Result.Handoff.SHA256 {
		return invalid()
	}
	owner := a.DeliveryOwner
	if owner == nil || owner.RunID != s.Result.RunID || owner.RequestSHA256 != s.Result.RequestSHA256 || owner.PreparationID != s.Request.PreparationID || owner.ActorClaim != s.Request.Delivery.OwnerClaim {
		return invalid()
	}
	if a.TaskID != c.TaskResult.TaskID || a.TaskResultID != c.TaskResult.ID || a.HumanQARecordID != human.HumanQARecordID || p.Task.TaskID != c.TaskResult.TaskID || p.Task.ResultOID != c.OID || p.Task.ResultTree != c.Tree || p.Task.SourceLocator != c.TaskResult.SourceLocator || p.Task.SourceRef != c.TaskResult.SourceRef || p.TaskResult.ID != c.TaskResult.ID || p.HumanQA.ID != human.HumanQARecordID || p.HumanQA.Outcome != "pass" || p.Epic.ParentLocator != agreement.TargetWorktree || p.Epic.ParentRef != agreement.TargetRef || p.Project.ProjectID != agreement.ProjectID || p.Repository.RepoID != agreement.RepoID || p.Epic.EpicID != agreement.EpicID {
		return invalid()
	}
	h := s.Result.Delivery.HumanIntegration
	if c.HumanQA == nil || !equal(*c.HumanQA, h.HumanQA) || rb.HumanQA == nil || !equal(*rb.HumanQA, h.HumanQA) {
		return invalid()
	}
	latest, err := LatestDeliveryAcceptanceHumanQA(registry.HumanQARecords, c.TaskResult, agreement, s.Request.SchemaVersion == 3)
	if err != nil || latest == nil || !equal(*latest, h.HumanQA) {
		return invalid()
	}
	own := closeout.Plan.Ownership
	if !own.Released || own.EvidenceLocator != human.ReleaseLocator || own.EvidenceSHA256 != human.ReleaseSHA256 || closeout.Plan.TargetRef != agreement.TargetRef || closeout.Plan.ReturnLocator != agreement.TargetWorktree {
		return invalid()
	}
	return digest([]any{c.Integration, v.Result, ir, a, h, registry.RawSHA256, closeout}), nil
}
