package workflowhandoff

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryAllowedEffects names the complete delivery boundary, independently
// of implementation commands. It deliberately grants no merge or target push.
func DeliveryAllowedEffects(a workspace.DeliveryAgreement) []string {
	if a.Mode == workspace.DeliveryPullRequest {
		return []string{"pull_request", "push"}
	}
	return []string{a.Mode}
}

func deliveryAgreement(o canonicaljson.Object) (*workspace.DeliveryAgreement, error) {
	v, _ := objectMember(o, "delivery_binding")
	b, _ := v.(canonicaljson.Object)
	v, found := objectMember(b, "agreement")
	if !found {
		return nil, nil
	}
	return workspace.DeliveryAgreementFromSpec(canonicaljson.Object{{Name: "delivery", Value: v}})
}

// ValidateDeliveryMandate rejects mode/target substitution before startup or a
// delivery effect. Historical contracts must retain their absent agreement.
func ValidateDeliveryMandate(raw []byte, agreement *workspace.DeliveryAgreement) error {
	d, err := decodeHandoffDraft(raw)
	if err != nil {
		return err
	}
	a, err := deliveryAgreement(d.Value)
	if err != nil {
		return err
	}
	if !canonicalEqual(bridgeValue(a), bridgeValue(agreement)) {
		return fmt.Errorf("delivery agreement differs from the frozen native mandate; revise and authorize a new execution")
	}
	if agreement == nil {
		return nil
	}
	if agreement.SourceRef != d.Binding.TargetRef || agreement.ProjectID != d.Binding.ProjectID || agreement.RepoID != d.Binding.RepoID {
		return fmt.Errorf("delivery agreement source or repository differs from the native target")
	}
	return validateDeliveryAgreementEffects(d.Value, agreement)
}

func validateDeliveryAgreementEffects(o canonicaljson.Object, a *workspace.DeliveryAgreement) error {
	if a == nil {
		return nil
	}
	if a.SourceRef == "" {
		return fmt.Errorf("delivery mandate requires the bound Task source ref")
	}
	if deliveryMode(o) != "owner" {
		return nil
	}
	v, _ := objectMember(o, "authority")
	authority, _ := v.(canonicaljson.Object)
	allowed, _ := objectMember(authority, "allowed_effects")
	seen := map[string]int{}
	for _, item := range allowed.([]canonicaljson.Value) {
		effect := item.(canonicaljson.Object)
		kind := objectString(effect, "type")
		scope, _ := objectMember(effect, "scope")
		so, _ := scope.(canonicaljson.Object)
		switch kind {
		case "git_ref_write", "git_index_write", "git_commit":
			if objectString(so, "ref") != a.SourceRef {
				return fmt.Errorf("delivery implementation writes may only target the Task source ref")
			}
		case "filesystem_write":
			if objectString(so, "kind") != "task_worktree" {
				return fmt.Errorf("delivery implementation filesystem writes must bind the Task worktree")
			}
		case "command_execute", "process_start", "install", "push", "pull_request", "integration":
		default:
			return fmt.Errorf("effect %s is outside the single Task delivery mandate", kind)
		}
		if kind != "push" && kind != "pull_request" && kind != "integration" {
			continue
		}
		seen[kind]++
		want := deliveryObject(map[string]any{"kind": "delivery", "agreement": a})
		if !canonicalEqual(scope, want) || objectInt(effect, "max_occurrences") != 1 {
			return fmt.Errorf("delivery effect %s differs from the exact frozen agreement", kind)
		}
	}
	if a.Mode == workspace.DeliveryPullRequest {
		if seen["push"] != 1 || seen["pull_request"] != 1 || seen["integration"] != 0 {
			return fmt.Errorf("PR mandate requires only its exact push and pull_request effects")
		}
	} else if seen["integration"] != 1 || seen["push"] != 0 || seen["pull_request"] != 0 {
		return fmt.Errorf("local mandate requires only its exact integration effect")
	}
	return nil
}

func validateDeliveryEffectScope(kind string, v canonicaljson.Value) error {
	if kind != "push" && kind != "pull_request" && kind != "integration" {
		return fmt.Errorf("unsupported scoped delivery effect")
	}
	m, err := exactObject(v, "delivery effect scope", "kind", "agreement")
	if err != nil || objectMapString(m, "kind") != "delivery" {
		return fmt.Errorf("invalid delivery effect scope")
	}
	a, err := workspace.DeliveryAgreementFromSpec(canonicaljson.Object{{Name: "delivery", Value: m["agreement"]}})
	if err != nil {
		return err
	}
	if a == nil || a.SourceRef == "" || (kind == "integration") == (a.Mode == workspace.DeliveryPullRequest) {
		return fmt.Errorf("effect is outside the selected delivery mode")
	}
	return nil
}

func validateDeliverySpecAgreement(handoff canonicaljson.Object, spec canonicaljson.Value) error {
	a, err := deliveryAgreement(handoff)
	if err != nil {
		return err
	}
	want, err := workspace.DeliveryAgreementFromSpec(spec)
	if err != nil {
		return err
	}
	if !canonicalEqual(bridgeValue(a), bridgeValue(want)) {
		return fmt.Errorf("delivery mandate differs from its frozen execution Spec")
	}
	return validateDeliveryAgreementEffects(handoff, a)
}

// DeliveryCandidateRuntimeBinding is read only; it binds native candidate
// evidence to the owning workflow without using a mutable registry lookup.
func DeliveryCandidateRuntimeBinding(raw []byte) (runID, requestSHA string, err error) {
	var value struct {
		Delivery struct {
			RunID      string `json:"workflow_run_id"`
			RequestSHA string `json:"request_sha256"`
		} `json:"delivery_binding"`
	}
	err = json.Unmarshal(raw, &value)
	return value.Delivery.RunID, value.Delivery.RequestSHA, err
}
