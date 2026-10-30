package workflowhandoff

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Version 3 separates a delivery owner's lifetime from candidate verification.
// Null limits are preserved in the contract. They are not a claim of unlimited
// runtime permission, nor a fabricated count of implementation iterations.
func validateDeliveryBinding(o canonicaljson.Object) error {
	if protocol, found := objectMember(o, "report_startup_protocol"); found && (taskSpecVersion(o) != 2 || protocol != workspace.DeliveryMandateReceiptV1) {
		return fmt.Errorf("report receipt requires Handoff@2 and mandate_receipt_v1")
	}

	v, found := objectMember(o, "delivery_binding")
	if err := validateParallelHandoff(o); err != nil {
		return err
	}
	if !deliveryHandoff(o) {
		if found {
			return fmt.Errorf("delivery_binding requires Handoff@3")
		}
		return nil
	}
	if !found {
		return fmt.Errorf("Handoff@3 requires delivery_binding")
	}
	fields := []string{"mode", "human_actor", "owner_claim", "parent_handoff_locator", "parent_handoff_sha256", "candidate_key"}
	binding, _ := v.(canonicaljson.Object)
	if _, ok := objectMember(binding, "startup_protocol"); ok {
		fields = append(fields, "startup_protocol")
		if objectString(binding, "startup_protocol") != workspace.DeliveryMandateReceiptV1 {
			return fmt.Errorf("unsupported startup protocol")
		}
	}
	if _, ok := objectMember(binding, "agreement"); ok {
		fields = append(fields, "agreement")
		if objectString(binding, "mode") == "candidate" {
			fields = append(fields, "workflow_run_id", "request_sha256")
		}
	}
	if _, ok := objectMember(binding, "acceptance_amendment"); ok {
		fields = append(fields, "acceptance_amendment")
	}
	m, e := exactObject(v, "delivery_binding", fields...)
	if e != nil {
		return e
	}
	for _, key := range []string{"human_actor", "owner_claim"} {
		if e = validatePlainText(key, objectMapString(m, key), 1, 256); e != nil {
			return e
		}
	}
	switch objectMapString(m, "mode") {
	case "owner":
		if m["parent_handoff_locator"] != nil || m["parent_handoff_sha256"] != nil || m["candidate_key"] != nil {
			return fmt.Errorf("owner cannot claim a candidate parent")
		}
	case "candidate":
		if _, e = absoluteCleanString(m, "parent_handoff_locator", "delivery_binding"); e != nil {
			return e
		}
		if !validateDigest(objectMapString(m, "parent_handoff_sha256")) || validateKey("candidate_key", objectMapString(m, "candidate_key")) != nil {
			return fmt.Errorf("invalid candidate parent or key")
		}
	default:
		return fmt.Errorf("invalid delivery mode")
	}
	a, e := deliveryAgreement(o)
	if e != nil {
		return e
	}
	if a != nil {
		if a.SourceRef == "" {
			return fmt.Errorf("delivery agreement requires the exact source ref")
		}
		if objectMapString(m, "mode") == "candidate" && (validatePlainText("workflow_run_id", objectMapString(m, "workflow_run_id"), 1, 256) != nil || !validateDigest(objectMapString(m, "request_sha256"))) {
			return fmt.Errorf("candidate requires its exact workflow authority binding")
		}
	}
	if _, e = effectiveDeliveryAgreement(o); e != nil {
		return e
	}
	_, e = taskSpecBasis(o)
	return e
}

func unlimitedDeliveryBudget(o canonicaljson.Object) bool {
	b, _ := objectMember(o, "budget")
	m, _ := b.(canonicaljson.Object)
	v, found := objectMember(m, "max_rounds")
	return deliveryHandoff(o) && found && v == nil
}

func deliveryMode(o canonicaljson.Object) string {
	v, _ := objectMember(o, "delivery_binding")
	m, _ := v.(canonicaljson.Object)
	return objectString(m, "mode")
}

func validateDeliveryOrTaskTarget(d Dependencies, contract canonicaljson.Object, target TargetObservation, basis *workspace.TaskSpecBasis) (workspace.TaskSpecEvaluation, error) {
	if !deliveryHandoff(contract) || deliveryMode(contract) != "candidate" {
		return d.taskSpecSession.ValidateTarget(target.Worktree, target.Ref, target.GitCommonDir, basis)
	}
	binding, _ := objectMember(contract, "delivery_binding")
	b := binding.(canonicaljson.Object)
	parent, e := d.Store.ReadByLocator(objectString(b, "parent_handoff_locator"))
	if e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	if !deliveryHandoff(parent.Handoff.Value) || deliveryMode(parent.Handoff.Value) != "owner" || parent.HeadState != "ready" || parent.Start == nil || parent.Start.Outcome != "started" || parent.Terminal != nil || parent.Handoff.SHA256 != objectString(b, "parent_handoff_sha256") {
		return workspace.TaskSpecEvaluation{}, fmt.Errorf("candidate requires its active accepted delivery owner")
	}
	pb, _ := objectMember(parent.Handoff.Value, "task_spec_binding")
	cb, _ := objectMember(contract, "task_spec_binding")
	pd, _ := objectMember(parent.Handoff.Value, "delivery_binding")
	owner := pd.(canonicaljson.Object)
	if !canonicalEqual(pb, cb) || objectString(owner, "owner_claim") != objectString(b, "owner_claim") || objectString(owner, "human_actor") != objectString(b, "human_actor") || target.Worktree != parent.Handoff.Target.Worktree || target.Ref != parent.Handoff.Target.Ref || target.GitCommonDir != parent.Handoff.Target.GitCommonDir {
		return workspace.TaskSpecEvaluation{}, fmt.Errorf("candidate differs from the delivery owner binding")
	}
	pa, e := deliveryAgreement(parent.Handoff.Value)
	if e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	ca, e := deliveryAgreement(contract)
	if e != nil || !canonicalEqual(bridgeValue(pa), bridgeValue(ca)) {
		return workspace.TaskSpecEvaluation{}, fmt.Errorf("candidate delivery agreement differs from its frozen owner")
	}
	if _, e = effectiveDeliveryAgreement(contract); e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	amendment, e := deliveryAcceptanceAmendment(contract)
	if e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	if e = validateDeliveryAmendmentEvidence(d, amendment); e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	effective, err := parallelEffectiveBasis(contract)
	if err != nil {
		return workspace.TaskSpecEvaluation{}, err
	}
	return d.taskSpecSession.ValidateDeliveryCandidateTargetBasis(target.Worktree, target.Ref, target.GitCommonDir, basis, effective)
}

func ValidateDeliveryTaskRunDraft(raw []byte) (TaskRunDraft, error) {
	d, e := decodeHandoffDraft(raw)
	if e != nil {
		return TaskRunDraft{}, e
	}
	b, e := taskSpecBasis(d.Value)
	if e != nil || !deliveryHandoff(d.Value) || deliveryMode(d.Value) != "owner" || b == nil || d.Binding.Status.Mode != "clean" {
		return TaskRunDraft{}, fmt.Errorf("delivery start requires an owner Handoff@3, execution Spec and clean target")
	}
	a, e := deliveryAgreement(d.Value)
	if e != nil {
		return TaskRunDraft{}, e
	}
	if e = validateDeliveryAgreementEffects(d.Value, a); e != nil {
		return TaskRunDraft{}, e
	}
	return TaskRunDraft{b, d.Binding.TargetWorktree, d.Binding.TargetRef, d.Binding.ExpectedOID, d.GoalTitle, d.PublicationKey}, nil
}

func ObserveDeliveryTaskRun(d Dependencies, root, preparationID string, raw []byte) (TaskRunObservation, error) {
	return observeTaskRun(d, root, preparationID, raw, true)
}

func BuildDeliveryTaskRunStart(d Dependencies, locator, actor, surface, session, runtime, model string, sandbox, issues []byte, acceptance string) ([]byte, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, e
	}
	if !deliveryHandoff(s.Handoff.Value) {
		return nil, fmt.Errorf("delivery start requires Handoff@3")
	}
	return BuildTaskRunStart(d, locator, actor, surface, session, runtime, model, sandbox, issues, acceptance)
}

func BuildDeliveryTaskRunResult(d Dependencies, locator, publicationKey string, technical []byte, recorder workspace.TaskRecorderRecord) ([]byte, string, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, "", e
	}
	if !deliveryHandoff(s.Handoff.Value) || deliveryMode(s.Handoff.Value) != "candidate" {
		return nil, "", fmt.Errorf("delivery result requires candidate-control evidence")
	}
	return buildTaskRunResult(d, locator, publicationKey, technical, recorder)
}

// BuildDeliveryHandoffDraft binds prepared execution, not future test success.
type DeliveryDraftAuthority struct {
	AllowLocalInstall bool
	StartupProtocol   string
	Agreement         *workspace.DeliveryAgreement
}

// DeliveryLocalInstallAllowed reads the effective authority of a validated
// handoff. A request's descriptive flag must not widen or conceal this grant.
func DeliveryLocalInstallAllowed(raw []byte) (bool, error) {
	draft, err := decodeHandoffDraft(raw)
	if err != nil {
		return false, err
	}
	effects, _ := objectMember(draft.Authority, "allowed_effects")
	for _, value := range effects.([]canonicaljson.Value) {
		if objectString(value.(canonicaljson.Object), "type") == "install" {
			return true, nil
		}
	}
	return false, nil
}

// deliveryCleanupAgreement selects the agreements whose new mandates grant
// physical removal of the released Task source after automatic local delivery.
func deliveryCleanupAgreement(a *workspace.DeliveryAgreement) bool {
	return a != nil && a.Mode != workspace.DeliveryPullRequest && !a.HumanOwnedIntegration() && a.AutomaticAcceptance()
}

// DeliveryCleanupAllowed reads whether a validated frozen handoff grants
// source cleanup. Older mandates forbid it and keep their retained source.
func DeliveryCleanupAllowed(raw []byte) (bool, error) {
	draft, err := decodeHandoffDraft(raw)
	if err != nil {
		return false, err
	}
	effects, _ := objectMember(draft.Authority, "allowed_effects")
	for _, value := range effects.([]canonicaljson.Value) {
		if objectString(value.(canonicaljson.Object), "type") == "cleanup" {
			return true, nil
		}
	}
	return false, nil
}

func bridgeValue(v any) canonicaljson.Value {
	switch value := v.(type) {
	case canonicaljson.Object:
		return value
	case []canonicaljson.Value:
		return value
	case map[string]any:
		o := canonicaljson.Object{}
		for key, item := range value {
			o = append(o, canonicaljson.Member{Name: key, Value: bridgeValue(item)})
		}
		return o
	}
	b, _ := workspace.MarshalTaskSpecValue(v)
	out, _ := canonicaljson.DecodeStrict(b)
	return out
}

func deliveryObject(v map[string]any) canonicaljson.Object {
	return bridgeValue(v).(canonicaljson.Object)
}
func deliveryBudget(max any) canonicaljson.Object {
	return deliveryObject(map[string]any{"max_rounds": max, "round_definition": map[string]any{"unit": "candidate_verification_attempt", "command_retry_consumes_round": false, "retry_condition": "only_if_no_effect_started"}})
}
func deliveryVerifier(path, cwd string) canonicaljson.Object {
	return deliveryObject(map[string]any{"id": "goal-acceptance", "argv": []string{"/bin/sh", path}, "cwd": cwd, "env": []string{}, "expected_exit": 0, "stop_on_failure": true, "evidence": map[string]any{"capture_stdout": true, "capture_stderr": true, "classification": "workspace_internal", "binding": "target_oid"}})
}
func deliveryStops() []canonicaljson.Value {
	out := []canonicaljson.Value{}
	for _, kind := range []string{"product_decision_required", "scope_or_authority_expansion", "target_or_input_drift", "unknown_or_partial_effect", "unexpected_sensitive_data", "round_budget_exhausted"} {
		out = append(out, deliveryObject(map[string]any{"type": kind, "description": "Stop the dependent effect when this condition occurs; a null budget has no fixed round limit."}))
	}
	return out
}
func deliveryReporting() canonicaljson.Object {
	return deliveryObject(map[string]any{"summary_max_codepoints": 240, "meaning_max_codepoints": 600, "required_start_fields": []string{"acceptance", "binding", "contract_digests", "issues", "observed_inputs", "observed_project", "observed_target", "observed_workspace", "principal", "receipt_id", "sandbox"}, "required_terminal_fields": []string{"artifacts", "binding", "evidence_gaps", "final_target", "forbidden_effects_observed", "meaning", "observed_effects", "principal", "reported_outcome", "result_id", "review", "rounds_used", "start_binding", "stop_reasons", "summary", "verifier_results"}})
}

// DeliveryStartupProtocol reads the immutable mandate, never a recipient permission claim.
func DeliveryStartupProtocol(raw []byte) (string, error) {
	d, err := decodeHandoffDraft(raw)
	if err != nil {
		return "", err
	}
	return deliveryStartupProtocol(d.Value), nil
}
func deliveryStartupProtocol(o canonicaljson.Object) string {
	v, _ := objectMember(o, "delivery_binding")
	b, _ := v.(canonicaljson.Object)
	return objectString(b, "startup_protocol")
}
