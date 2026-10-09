package workflowhandoff

import (
	"crypto/hmac"
	"errors"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func taskSpecVersion(value canonicaljson.Object) int64 { return objectInt(value, "schema_version") }
func taskSpecBasis(value canonicaljson.Object) (*workspace.TaskSpecBasis, error) {
	v, found := objectMember(value, "task_spec_binding")
	if taskSpecVersion(value) == 1 {
		if found {
			return nil, fmt.Errorf("version 1 cannot contain task_spec_binding")
		}
		return nil, nil
	}
	if (taskSpecVersion(value) != 2 && taskSpecVersion(value) != 3) || !found {
		return nil, fmt.Errorf("version 2 or 3 requires explicit task_spec_binding")
	}
	return workspace.ValidateTaskSpecBasisValue(v)
}
func taskSpecLegacyProjection(value canonicaljson.Object) canonicaljson.Object {
	return replaceObjectMember(removeObjectMember(removeObjectMember(value, "task_spec_binding"), "delivery_binding"), "schema_version", int64(1))
}
func decodeHandoffDraft(input []byte) (handoffDraft, error) {
	v, e := canonicaljson.DecodeStrict(input)
	if e != nil {
		return handoffDraft{}, schemaError("handoff draft", e)
	}
	o, ok := v.(canonicaljson.Object)
	if !ok {
		return handoffDraft{}, schemaError("handoff draft", fmt.Errorf("expected object"))
	}
	if taskSpecVersion(o) == 1 {
		return decodeHandoffDraftV1(input)
	}
	if len(input) > maxHandoffBytes {
		return handoffDraft{}, classified(ErrorPayloadTooLarge, "handoff draft exceeds 256 KiB", nil)
	}
	basis, e := taskSpecBasis(o)
	if e != nil {
		return handoffDraft{}, schemaError("handoff draft", e)
	}
	if e = validateDeliveryBinding(o); e != nil {
		return handoffDraft{}, schemaError("handoff draft", e)
	}
	legacy, e := canonicaljson.Marshal(taskSpecLegacyProjection(o))
	if e != nil {
		return handoffDraft{}, e
	}
	d, e := decodeHandoffDraftWithBudget(legacy, taskSpecVersion(o) == 3)
	if e != nil {
		return d, e
	}
	d.Value = o
	d.Canonical, e = canonicaljson.Marshal(o)
	if e != nil {
		return d, e
	}
	if len(d.Canonical) > maxHandoffBytes {
		return handoffDraft{}, classified(ErrorPayloadTooLarge, "canonical handoff exceeds 256 KiB", nil)
	}
	d.Digest = digestBytes(d.Canonical)
	if basis != nil && d.Binding.Status.Mode != "clean" {
		return handoffDraft{}, schemaError("handoff draft", fmt.Errorf("Task Spec delivery requires clean status policy"))
	}
	return d, nil
}
func validateStoredHandoff(o canonicaljson.Object) error {
	if taskSpecVersion(o) == 1 {
		return validateStoredHandoffV1(o)
	}
	if _, e := taskSpecBasis(o); e != nil {
		return e
	}
	if e := validateDeliveryBinding(o); e != nil {
		return e
	}
	if err := validateStoredHandoffWithBudget(taskSpecLegacyProjection(o), taskSpecVersion(o) == 3); err != nil {
		return err
	}
	a, err := deliveryAgreement(o)
	if err != nil {
		return err
	}
	return validateDeliveryAgreementEffects(o, a)
}
func decodeStartDraft(input []byte) (startDraft, error) {
	v, e := canonicaljson.DecodeStrict(input)
	if e != nil {
		return startDraft{}, schemaError("start draft", e)
	}
	o, ok := v.(canonicaljson.Object)
	if !ok {
		return startDraft{}, schemaError("start draft", fmt.Errorf("expected object"))
	}
	if taskSpecVersion(o) == 1 {
		return decodeStartDraftV1(input)
	}
	if len(input) > maxStartBytes {
		return startDraft{}, classified(ErrorPayloadTooLarge, "start draft exceeds 128 KiB", nil)
	}
	if _, e = taskSpecBasis(o); e != nil {
		return startDraft{}, schemaError("start draft", e)
	}
	if _, found := objectMember(o, "delivery_binding"); found {
		return startDraft{}, schemaError("start draft", fmt.Errorf("delivery_binding belongs to the handoff, not a start receipt"))
	}
	projected := taskSpecLegacyProjection(o)
	// Validate the added issue type with the same strict shape and ordering before
	// passing the remaining version-1 fields through the historical validator.
	if issues, found := objectMember(o, "issues"); found {
		a, ok := issues.([]canonicaljson.Value)
		if !ok {
			return startDraft{}, schemaError("start issues", fmt.Errorf("expected array"))
		}
		filtered := []canonicaljson.Value{}
		last := ""
		for _, v := range a {
			m, e := exactObject(v, "issue", "type", "detail")
			if e != nil {
				return startDraft{}, e
			}
			kind := objectMapString(m, "type")
			detail := objectMapString(m, "detail")
			key := kind + "\x00" + detail
			if last != "" && key <= last {
				return startDraft{}, fmt.Errorf("issues must be sorted unique")
			}
			last = key
			if kind == "task_spec" {
				if e = validatePlainText("issue detail", detail, 1, 2000); e != nil {
					return startDraft{}, e
				}
			} else {
				filtered = append(filtered, v)
			}
		}
		// Keep the nonempty-issues invariant in the historical validator even
		// when the only issue is the version-2 extension validated above.
		if len(a) > 0 && len(filtered) == 0 {
			filtered = []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: "target"}, {Name: "detail", Value: "Task Spec issue validated by version 2."}}}
		}
		projected = replaceObjectMember(projected, "issues", filtered)
	}
	b, e := canonicaljson.Marshal(projected)
	if e != nil {
		return startDraft{}, e
	}
	d, e := decodeStartDraftV1(b)
	if e != nil {
		return d, e
	}
	d.Value = o
	d.Canonical, e = canonicaljson.Marshal(o)
	return d, e
}
func validateStoredAccepted(value canonicaljson.Object, phase, capabilityID, secret string) error {
	if phase != "start" || taskSpecVersion(value) == 1 {
		return validateStoredAcceptedV1(value, phase, capabilityID, secret)
	}
	if _, e := taskSpecBasis(value); e != nil {
		return e
	}
	observation, found := objectMember(value, "task_spec_observation")
	if !found {
		return fmt.Errorf("version 2 start requires service observation")
	}
	basis, _ := taskSpecBasis(value)
	if e := validateTaskSpecObservation(observation, basis, objectString(value, "acceptance")); e != nil {
		return e
	}
	projected := removeObjectMember(taskSpecLegacyProjection(value), "task_spec_observation")
	if _, e := validateFinalDocument(projected, "ply.workflow.start-receipt", startDocumentFields); e != nil {
		return e
	}
	proofValue, _ := objectMember(value, "capability_proof")
	proof, e := exactObject(proofValue, "proof", "capability_id", "algorithm", "value")
	if e != nil {
		return e
	}
	if objectMapString(proof, "capability_id") != capabilityID {
		return fmt.Errorf("capability ID differs")
	}
	payload := removeObjectMember(value, "capability_proof")
	b, e := canonicaljson.Marshal(payload)
	if e != nil {
		return e
	}
	expected, e := capabilityProof(secret, b)
	if e != nil {
		return e
	}
	if !hmac.Equal([]byte(expected), []byte(objectMapString(proof, "value"))) {
		return fmt.Errorf("capability proof does not match version 2 payload")
	}
	draft := replaceObjectMember(removeObjectMember(payload, "task_spec_observation"), "kind", "ply.workflow.start-receipt-draft")
	b, e = canonicaljson.Marshal(draft)
	if e != nil {
		return e
	}
	_, e = decodeStartDraft(b)
	return e
}
func validateTaskSpecObservation(v canonicaljson.Value, b *workspace.TaskSpecBasis, acceptance string) error {
	if b == nil {
		if v != nil {
			return fmt.Errorf("general handoff requires null Task Spec observation")
		}
		return nil
	}
	m, e := exactObject(v, "task_spec_observation", "expected_basis_sha256", "observed_basis_sha256", "registry_sha256", "current_problem", "current_selection", "selection_freshness", "content_integrity", "target_freshness", "validated_at_utc", "reasons")
	if e != nil {
		return e
	}
	bv, e := taskSpecJSONValue(b)
	if e != nil {
		return e
	}
	bytes, e := canonicaljson.Marshal(bv)
	if e != nil {
		return e
	}
	if objectMapString(m, "expected_basis_sha256") != digestBytes(bytes) {
		return fmt.Errorf("expected Task basis digest differs")
	}
	for _, k := range []string{"observed_basis_sha256", "registry_sha256"} {
		if m[k] != nil {
			d, ok := m[k].(string)
			if !ok || !validateDigest(d) {
				return fmt.Errorf("invalid observation digest")
			}
		}
	}
	for _, k := range []string{"current_problem", "current_selection"} {
		if m[k] != nil {
			field := "problem"
			if k == "current_selection" {
				field = "selection"
			}
			if _, err := workspace.ValidateTaskSpecBasisValue(replaceObjectMember(bv.(canonicaljson.Object), field, m[k])); err != nil {
				return fmt.Errorf("invalid observation %s: %w", k, err)
			}
			expected := canonicaljson.Value(nil)
			if k == "current_problem" {
				expected, _ = taskSpecJSONValue(b.Problem)
			} else {
				expected, _ = taskSpecJSONValue(b.Selection)
			}
			if acceptance == "started" && !canonicalEqual(expected, m[k]) {
				return fmt.Errorf("started observation differs from basis")
			}
		}
	}
	if !setOf("current", "stale", "withdrawn", "none", "unknown")[objectMapString(m, "selection_freshness")] || !setOf("valid", "missing", "conflict", "unknown")[objectMapString(m, "content_integrity")] || !setOf("fresh", "stale", "unknown")[objectMapString(m, "target_freshness")] || !validateUTC(objectMapString(m, "validated_at_utc")) {
		return fmt.Errorf("invalid Task Spec observation")
	}
	reasons, e := sortedStringArray(m["reasons"], "task_spec reasons", true)
	if e != nil {
		return e
	}
	if acceptance == "started" && (objectMapString(m, "observed_basis_sha256") != objectMapString(m, "expected_basis_sha256") || objectMapString(m, "selection_freshness") != "current" || objectMapString(m, "content_integrity") != "valid" || objectMapString(m, "target_freshness") != "fresh" || m["registry_sha256"] == nil || m["current_problem"] == nil || m["current_selection"] == nil || len(reasons) != 0) {
		return fmt.Errorf("started requires an exact current Task Spec observation")
	}
	return nil
}
func taskSpecJSONValue(value any) (canonicaljson.Value, error) {
	b, e := workspace.MarshalTaskSpecValue(value)
	if e != nil {
		return nil, e
	}
	return canonicaljson.DecodeStrict(b)
}
func validateTaskSpecInputs(value canonicaljson.Object, required []workspace.TaskSpecRequiredInput, spec canonicaljson.Object) error {
	v, _ := objectMember(value, "inputs")
	inputs, e := validateInputs(v)
	if e != nil {
		return e
	}
	expected := map[string]workspace.TaskSpecRequiredInput{}
	for _, r := range required {
		expected[r.ID] = r
	}
	seen := map[string]bool{}
	for _, in := range inputs {
		if r, ok := expected[in.ID]; ok {
			if in.Locator != r.Locator || in.SHA256 != r.SHA256 || in.SizeBytes != r.SizeBytes || in.Role != r.Role || in.MediaType != r.MediaType || in.GitBinding != nil {
				return classified(ErrorTaskSpec, "required Task input differs: "+in.ID, nil)
			}
			seen[in.ID] = true
		} else if strings.HasPrefix(in.ID, "task-basis/") || strings.HasPrefix(in.ID, "task-doc/") || in.Role == "spec" || in.Role == "design" {
			return classified(ErrorTaskSpec, "unselected or reserved Task input: "+in.ID, nil)
		}
	}
	if len(seen) != len(expected) {
		return classified(ErrorTaskSpec, "required Task input coverage is incomplete", nil)
	}
	verifiers := map[string]bool{}
	vs, _ := objectMember(value, "verifiers")
	for _, v := range vs.([]canonicaljson.Value) {
		o := v.(canonicaljson.Object)
		verifiers[objectString(o, "id")] = true
	}
	reqs, _ := objectMember(spec, "requirements")
	if reqs != nil {
		for _, v := range reqs.([]canonicaljson.Value) {
			ids, _ := objectMember(v.(canonicaljson.Object), "verification_ids")
			for _, id := range ids.([]canonicaljson.Value) {
				if !verifiers[id.(string)] {
					return classified(ErrorTaskSpec, "requirement verification ID is not declared by handoff: "+id.(string), nil)
				}
			}
		}
	}
	// Reporting retains its exact version-1 field set. Its required artifacts field
	// is specialized by the explicit reporting instruction in the procedure.
	procedure, _ := objectMember(value, "procedure")
	p, _ := canonicaljson.Marshal(procedure)
	if !strings.Contains(string(p), "task-requirements") {
		return classified(ErrorTaskSpec, "procedure must explicitly require reporting the task-requirements artifact", nil)
	}
	return nil
}
func validateCreateTaskSpec(d Dependencies, draft handoffDraft, target TargetObservation) error {
	if d.taskSpecSession == nil {
		return nil
	}
	basis, e := taskSpecBasis(draft.Value)
	if e != nil {
		return e
	}
	evaluation, e := validateDeliveryOrTaskTarget(d, draft.Value, target, basis)
	if e != nil {
		return classified(ErrorTaskSpec, e.Error(), e)
	}
	if basis != nil {
		return validateTaskSpecInputs(draft.Value, evaluation.RequiredInputs, evaluation.Spec)
	}
	return nil
}
func withTaskSpecWorkspace(d Dependencies, root string, fn func(Dependencies) error) error {
	if d.taskSpecSession != nil {
		return fn(d)
	}
	if d.TaskWorkspace == nil {
		return classified(ErrorTaskSpec, "workspace Task guard dependency is unavailable", nil)
	}
	return workspace.WithTaskSpecRegistrationSnapshot(*d.TaskWorkspace, root, func(s *workspace.TaskSpecSession) error { d.taskSpecSession = s; return fn(d) })
}
func validateTaskSpecStart(d Dependencies, snapshot Snapshot, draft startDraft) (canonicaljson.Value, error) {
	if taskSpecVersion(snapshot.Handoff.Value) != taskSpecVersion(draft.Value) {
		return nil, classified(ErrorTaskSpec, "handoff and start schema versions differ", nil)
	}
	basis, e := taskSpecBasis(snapshot.Handoff.Value)
	if e != nil {
		return nil, e
	}
	reported, e := taskSpecBasis(draft.Value)
	if e != nil {
		return nil, e
	}
	bv, _ := taskSpecJSONValue(basis)
	rv, _ := taskSpecJSONValue(reported)
	if !canonicalEqual(bv, rv) {
		return nil, classified(ErrorTaskSpec, "start Task basis differs from handoff", nil)
	}
	if d.taskSpecSession == nil {
		return nil, classified(ErrorTaskSpec, "locked workspace Task snapshot is required", nil)
	}
	eval, e := validateDeliveryOrTaskTarget(d, snapshot.Handoff.Value, snapshot.Handoff.Target, basis)
	if e == nil && basis != nil {
		e = validateTaskSpecInputs(snapshot.Handoff.Value, eval.RequiredInputs, eval.Spec)
	}
	if e != nil {
		if draft.Acceptance == "started" || basis == nil {
			return nil, classified(ErrorTaskSpec, e.Error(), e)
		}
		code := "task_content_observation_unknown"
		var contentErr *workspace.TaskContentError
		if errors.As(e, &contentErr) {
			code = contentErr.Code
		}
		eval.Reasons = append(eval.Reasons, code)
		if eval.SelectionFreshness == "" {
			eval.SelectionFreshness = "unknown"
		}
		if eval.ContentIntegrity == "" || eval.ContentIntegrity == "not_recorded" {
			eval.ContentIntegrity = "unknown"
		}
		if eval.TargetFreshness == "" {
			eval.TargetFreshness = "unknown"
		}
	}
	return d.taskSpecSession.StartObservation(basis, eval), nil
}

func historicalHandoffTaskSpec(d Dependencies, snapshot Snapshot) (*workspace.TaskSpecBasis, *workspace.TaskSpecEvaluation, error) {
	b, e := taskSpecBasis(snapshot.Handoff.Value)
	if e != nil {
		return nil, nil, e
	}
	if b == nil {
		return nil, nil, nil
	}
	if d.TaskWorkspace == nil {
		return b, nil, fmt.Errorf("Task workspace reader is unavailable")
	}
	eval, e := workspace.ReadHistoricalTaskSpec(*d.TaskWorkspace, snapshot.Handoff.Workspace.Root, *b)
	if e != nil {
		return b, &eval, e
	}
	if e = validateTaskSpecInputs(snapshot.Handoff.Value, eval.RequiredInputs, eval.Spec); e != nil {
		return b, &eval, e
	}
	if taskSpecVersion(snapshot.Handoff.Value) == 3 {
		amendment, err := deliveryAcceptanceAmendment(snapshot.Handoff.Value)
		if err != nil {
			return b, &eval, err
		}
		if err = validateDeliveryAmendmentEvidence(d, amendment); err != nil {
			return b, &eval, err
		}
		if amendment != nil {
			if err = validateDeliverySpecAgreement(snapshot.Handoff.Value, eval.Spec); err != nil {
				return b, &eval, err
			}
		}
	}
	if snapshot.Start != nil {
		if taskSpecVersion(snapshot.Start.Value) < 2 {
			return b, &eval, fmt.Errorf("Task basis requires version 2 start")
		}
		v, _ := objectMember(snapshot.Start.Value, "task_spec_binding")
		bv, _ := taskSpecJSONValue(b)
		if !canonicalEqual(v, bv) {
			return b, &eval, fmt.Errorf("historical start basis differs")
		}
		obs, _ := objectMember(snapshot.Start.Value, "task_spec_observation")
		if e = validateTaskSpecObservation(obs, b, snapshot.Start.Outcome); e != nil {
			return b, &eval, e
		}
	}
	return b, &eval, nil
}
func taskRequirementsCoverage(d Dependencies, snapshot Snapshot, b *workspace.TaskSpecBasis, eval *workspace.TaskSpecEvaluation) error {
	if b == nil {
		return nil
	}
	if snapshot.Terminal == nil || eval == nil {
		return fmt.Errorf("task-requirements evidence is not recorded")
	}
	values, _ := objectMember(snapshot.Terminal.Value, "artifacts")
	declared := false
	for _, raw := range values.([]canonicaljson.Value) {
		a := raw.(canonicaljson.Object)
		if objectString(a, "artifact_id") == "task-requirements" {
			declared = objectString(a, "kind") == "managed" && objectString(a, "media_type") == "application/json"
		}
	}
	if !declared {
		return fmt.Errorf("task-requirements must be a managed application/json artifact")
	}
	verifiers := workspaceVerifierProjection(snapshot)
	artifacts := workspaceArtifactProjection(snapshot, verifiers, workspaceReviewProjection(snapshot))
	target, _ := objectMember(snapshot.Terminal.Value, "final_target")
	t, _ := target.(canonicaljson.Object)
	for _, a := range artifacts {
		if a.ArtifactID == "task-requirements" {
			raw, e := d.Files.ReadFile(a.Locator)
			if e != nil {
				return e
			}
			if int64(len(raw)) != a.SizeBytes || digestBytes(raw) != a.SHA256 {
				return fmt.Errorf("task-requirements artifact differs")
			}
			return workspace.ValidateTaskRequirementEvidence(raw, *b, eval.Spec, objectString(t, "oid"), objectString(t, "tree"), verifiers, artifacts)
		}
	}
	return fmt.Errorf("required task-requirements artifact is missing")
}
