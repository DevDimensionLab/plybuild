package workflowhandoff

import (
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Version 3 separates a delivery owner's lifetime from candidate verification.
// Null limits are preserved in the contract. They are not a claim of unlimited
// runtime permission, nor a fabricated count of implementation iterations.
func validateDeliveryBinding(o canonicaljson.Object) error {
	v, found := objectMember(o, "delivery_binding")
	if taskSpecVersion(o) != 3 {
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
	if _, ok := objectMember(binding, "agreement"); ok {
		fields = append(fields, "agreement")
		if objectString(binding, "mode") == "candidate" {
			fields = append(fields, "workflow_run_id", "request_sha256")
		}
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
	b, e := taskSpecBasis(o)
	if e != nil || b == nil {
		return fmt.Errorf("delivery requires an exact execution Spec basis")
	}
	return nil
}

func unlimitedDeliveryBudget(o canonicaljson.Object) bool {
	b, _ := objectMember(o, "budget")
	m, _ := b.(canonicaljson.Object)
	v, found := objectMember(m, "max_rounds")
	return taskSpecVersion(o) == 3 && found && v == nil
}

func deliveryMode(o canonicaljson.Object) string {
	v, _ := objectMember(o, "delivery_binding")
	m, _ := v.(canonicaljson.Object)
	return objectString(m, "mode")
}

func validateDeliveryOrTaskTarget(d Dependencies, contract canonicaljson.Object, target TargetObservation, basis *workspace.TaskSpecBasis) (workspace.TaskSpecEvaluation, error) {
	if taskSpecVersion(contract) != 3 || deliveryMode(contract) != "candidate" {
		eval, err := d.taskSpecSession.ValidateTarget(target.Worktree, target.Ref, target.GitCommonDir, basis)
		if err == nil && taskSpecVersion(contract) == 3 {
			err = validateDeliverySpecAgreement(contract, eval.Spec)
		}
		return eval, err
	}
	binding, _ := objectMember(contract, "delivery_binding")
	b := binding.(canonicaljson.Object)
	parent, e := d.Store.ReadByLocator(objectString(b, "parent_handoff_locator"))
	if e != nil {
		return workspace.TaskSpecEvaluation{}, e
	}
	if taskSpecVersion(parent.Handoff.Value) != 3 || deliveryMode(parent.Handoff.Value) != "owner" || parent.HeadState != "ready" || parent.Start == nil || parent.Start.Outcome != "started" || parent.Terminal != nil || parent.Handoff.SHA256 != objectString(b, "parent_handoff_sha256") {
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
	return d.taskSpecSession.ValidateDeliveryCandidateTarget(target.Worktree, target.Ref, target.GitCommonDir, basis)
}

func ValidateDeliveryTaskRunDraft(raw []byte) (TaskRunDraft, error) {
	d, e := decodeHandoffDraft(raw)
	if e != nil {
		return TaskRunDraft{}, e
	}
	b, e := taskSpecBasis(d.Value)
	if e != nil || taskSpecVersion(d.Value) != 3 || deliveryMode(d.Value) != "owner" || b == nil || d.Binding.Status.Mode != "clean" {
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
	if taskSpecVersion(s.Handoff.Value) != 3 {
		return nil, fmt.Errorf("delivery start requires Handoff@3")
	}
	return BuildTaskRunStart(d, locator, actor, surface, session, runtime, model, sandbox, issues, acceptance)
}

func BuildDeliveryTaskRunResult(d Dependencies, locator, publicationKey string, technical []byte, recorder workspace.TaskRecorderRecord) ([]byte, string, error) {
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return nil, "", e
	}
	if taskSpecVersion(s.Handoff.Value) != 3 || deliveryMode(s.Handoff.Value) != "candidate" {
		return nil, "", fmt.Errorf("delivery result requires candidate-control evidence")
	}
	return buildTaskRunResult(d, locator, publicationKey, technical, recorder)
}

// BuildDeliveryHandoffDraft binds prepared execution, not future test success.
type DeliveryDraftAuthority struct {
	AllowLocalInstall bool
	Agreement         *workspace.DeliveryAgreement
}

func BuildDeliveryHandoffDraft(d Dependencies, preparationID, humanActor, ownerClaim, publicationKey, acceptancePath string, options ...DeliveryDraftAuthority) ([]byte, error) {
	if len(options) > 1 {
		return nil, fmt.Errorf("at most one delivery authority scope is supported")
	}
	allowInstall := len(options) == 1 && options[0].AllowLocalInstall
	ws, e := d.Workspace.ObserveContaining()
	if e != nil {
		return nil, e
	}
	var raw []byte
	e = WithTaskRunSnapshot(d, ws.Observation.Root, func(locked Dependencies) error {
		s := locked.taskSpecSession
		p, e := s.TaskRunPreparation(preparationID)
		if e != nil {
			return e
		}
		if p.Preparation == nil || p.State != "prepared" || p.Disposition != "current" || p.Freshness != "fresh" || len(p.Reasons) != 0 {
			return fmt.Errorf("delivery preparation is not current and fresh")
		}
		plan := p.Preparation.Plan
		eval, e := s.CurrentDeliveryTaskSpec(plan.TaskID)
		if e != nil {
			return e
		}
		if eval.Basis == nil {
			return fmt.Errorf("execution Spec basis missing")
		}
		agreement, e := workspace.DeliveryAgreementFromSpec(eval.Spec)
		if e != nil {
			return e
		}
		if len(options) == 1 && options[0].Agreement != nil && !canonicalEqual(bridgeValue(agreement), bridgeValue(options[0].Agreement)) {
			return fmt.Errorf("selected delivery agreement differs from the prepared execution Spec")
		}
		if !filepath.IsAbs(acceptancePath) || filepath.Clean(acceptancePath) != acceptancePath {
			return fmt.Errorf("acceptance path must be absolute and clean")
		}
		if objectString(eval.Spec, "acceptance_path") != acceptancePath {
			return fmt.Errorf("acceptance path differs from execution Spec")
		}
		inputs, e := taskSpecJSONValue(plan.RequiredInputs)
		if e != nil {
			return e
		}
		basis, _ := taskSpecJSONValue(eval.Basis)
		goal := deliveryObject(map[string]any{"title": string(plan.TaskID), "recipient_role": "Responsible feature implementor", "objective": "Deliver the selected goal and its exact execution Spec within the authorized local Task and Epic.", "done_when": "The actual human has passed the exact verified candidate and its authorized local integration is observed."})
		if agreement != nil && agreement.Mode == workspace.DeliveryPullRequest {
			goal = replaceObjectMember(goal, "done_when", "The exact verified candidate has actual human pass and a verified PR to the agreed repository/base; stop before merge.")
		} else if agreement != nil && agreement.HumanOwnedIntegration() {
			goal = replaceObjectMember(goal, "done_when", "Qualify the exact candidate, register Delivery and explicitly release source ownership to the human-started integration command. Stop before integration.")
		}
		if agreement != nil && agreement.AutomaticAcceptance() && !agreement.Acceptance.RequireHumanQA && !agreement.HumanOwnedIntegration() {
			goal = replaceObjectMember(goal, "done_when", "The exact candidate passed its automatic acceptance and review, and authorized native local Epic integration and closure are observed.")
		}
		procedure := []canonicaljson.Value{deliveryObject(map[string]any{"id": "implement", "instruction": "Own design, tests, permitted delegation, review and local delivery. Preserve actual results and report task-requirements for each verified candidate. Human QA and integration are separate facts; never fabricate them.", "required_before": []string{}})}
		if agreement != nil {
			goal = replaceObjectMember(goal, "objective", "Deliver the selected goal and exact execution Spec within the frozen single Task delivery agreement.")
			procedure[0] = replaceObjectMember(procedure[0].(canonicaljson.Object), "instruction", "Own design, tests, permitted delegation, review and the selected delivery. Preserve actual results and report the task-requirements artifact for each verified candidate. Technical qualification, actual human QA and observed delivery are separate facts.")
		}
		verifier := deliveryVerifier(acceptancePath, plan.WorktreePath)
		effects := []canonicaljson.Value{}
		add := func(id, kind string, scope canonicaljson.Object, max any) {
			effects = append(effects, deliveryObject(map[string]any{"id": id, "type": kind, "scope": scope, "max_occurrences": max, "sequence": len(effects) + 1}))
		}
		add("task-write", "filesystem_write", deliveryObject(map[string]any{"kind": "task_worktree"}), nil)
		for _, v := range []struct{ kind, op string }{{"git_index_write", "update_index"}, {"git_commit", "create_commit"}, {"git_ref_write", "update_ref"}} {
			add(v.kind, v.kind, deliveryObject(map[string]any{"kind": "git", "ref": "refs/heads/" + plan.Branch, "expected_before_oid": plan.ParentOID, "operation": v.op}), nil)
		}
		add("local-commands", "command_execute", deliveryObject(map[string]any{"kind": "command", "procedure_ids": []string{"implement"}, "verifier_ids": []string{"goal-acceptance"}}), nil)
		for _, kind := range []string{"process_start"} {
			add(kind, kind, deliveryObject(map[string]any{"kind": "boolean", "allowed": true}), nil)
		}
		if allowInstall {
			add("install", "install", deliveryObject(map[string]any{"kind": "boolean", "allowed": true}), nil)
		}
		if agreement != nil {
			scope := deliveryObject(map[string]any{"kind": "delivery", "agreement": agreement})
			if agreement.Mode == workspace.DeliveryPullRequest {
				add("delivery-push", "push", scope, 1)
				add("delivery-pr", "pull_request", scope, 1)
			} else if !agreement.HumanOwnedIntegration() {
				add("delivery-integration", "integration", scope, 1)
			}
		}
		forbidden := []canonicaljson.Value{}
		for _, kind := range []string{"cleanup", "deploy", "merge", "pull_request", "push", "release"} {
			if agreement != nil && agreement.Mode == workspace.DeliveryPullRequest && (kind == "pull_request" || kind == "push") {
				continue
			}
			forbidden = append(forbidden, deliveryObject(map[string]any{"type": kind, "reason": "Outside the selected local delivery contract."}))
		}
		if agreement != nil && (agreement.Mode == workspace.DeliveryPullRequest || agreement.HumanOwnedIntegration()) {
			forbidden = append(forbidden, deliveryObject(map[string]any{"type": "integration", "reason": "This developer delivery stops before human-started integration."}))
		}
		if !allowInstall {
			forbidden = append(forbidden, deliveryObject(map[string]any{"type": "install", "reason": "Local installation was not selected."}))
		}
		value := bridgeEnvelope("ply.workflow.handoff-draft")
		value = replaceObjectMember(value, "schema_version", int64(3))
		fields := map[string]any{"publication_key": publicationKey, "activity_key": publicationKey, "goal": goal, "recipient": map[string]any{"principal_id": "delivery-owner", "principal_kind": "human_started_agent", "runtime_constraints": []string{"actual-runtime-authority-required"}}, "binding_request": map[string]any{"project_id": plan.Target.ProjectID, "repo_id": plan.Target.RepoID, "target_worktree": plan.WorktreePath, "target_ref": "refs/heads/" + plan.Branch, "expected_oid": plan.ParentOID, "status_policy": map[string]any{"mode": "clean"}}, "inputs": inputs, "authority": map[string]any{"allowed_effects": effects, "forbidden_effects": forbidden, "human_gates": []string{"human_task_qa"}}, "budget": deliveryBudget(nil), "procedure": procedure, "verifiers": []canonicaljson.Value{verifier}, "stop_conditions": deliveryStops(), "reporting": deliveryReporting(), "task_spec_binding": basis, "delivery_binding": map[string]any{"mode": "owner", "human_actor": humanActor, "owner_claim": ownerClaim, "parent_handoff_locator": nil, "parent_handoff_sha256": nil, "candidate_key": nil}}
		if agreement != nil {
			fields["delivery_binding"].(map[string]any)["agreement"] = agreement
			if agreement.AutomaticAcceptance() && !agreement.Acceptance.RequireHumanQA {
				fields["authority"].(map[string]any)["human_gates"] = []string{}
			}
		}
		for k, v := range fields {
			value = append(value, canonicaljson.Member{Name: k, Value: bridgeValue(v)})
		}
		raw, e = canonicaljson.Marshal(value)
		if e != nil {
			return e
		}
		_, e = ValidateDeliveryTaskRunDraft(raw)
		if e != nil {
			return fmt.Errorf("validate delivery owner draft: %w", e)
		}
		return e
	})
	return raw, e
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
