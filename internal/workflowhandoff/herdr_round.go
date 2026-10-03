package workflowhandoff

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ValidateWorkflowRound validates the frozen Spec's required evidence before a
// transport reserves a report. It reuses native semantics without publication.
func ValidateWorkflowRound(d Dependencies, locator string, raw []byte) error {
	if e := ValidateTaskRunTerminal(d, locator, raw); e != nil {
		return e
	}
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return e
	}
	draft, e := decodeTerminalDraft(raw)
	if e != nil {
		return e
	}
	s.Terminal = &acceptedDocument{Value: draft.Value, Bytes: raw, Outcome: draft.Outcome}
	b, eval, e := historicalHandoffTaskSpec(d, s)
	if e != nil {
		return e
	}
	return taskRequirementsCoverage(d, s, b, eval)
}

// ValidateWorkflowPassedAssessment checks the native verifier projection before
// a reviewed report claims the passed gate, without creating a TaskResult.
func ValidateWorkflowPassedAssessment(d Dependencies, locator string, raw, technical []byte) error {
	if e := workspace.ValidateTaskRunTechnicalAssessment(technical); e != nil {
		return e
	}
	var assessment struct {
		Gate     string            `json:"gate"`
		Required []string          `json:"required_verifier_ids"`
		Debt     []json.RawMessage `json:"accepted_debt"`
	}
	if e := json.Unmarshal(technical, &assessment); e != nil {
		return e
	}
	if assessment.Gate != "passed" || len(assessment.Debt) != 0 {
		return fmt.Errorf("passed gate requires no accepted evidence debt")
	}
	s, e := d.Store.ReadByLocator(locator)
	if e != nil {
		return e
	}
	draft, e := decodeTerminalDraft(raw)
	if e != nil {
		return e
	}
	s.Terminal = &acceptedDocument{Value: draft.Value}
	policy, evidence := validateReportedPolicy(s)
	if len(policy) != 0 || len(evidence) != 0 {
		return fmt.Errorf("native report policy or evidence coverage failed: %v %v", policy, evidence)
	}
	expected, _ := objectMember(s.Handoff.Value, "verifiers")
	verifiers := expected.([]canonicaljson.Value)
	results := workspaceVerifierProjection(s)
	if len(assessment.Required) != len(verifiers) || len(results) != len(verifiers) {
		return fmt.Errorf("required verifier IDs differ from the native handoff")
	}
	for i, v := range verifiers {
		id := objectString(v.(canonicaljson.Object), "id")
		if assessment.Required[i] != id || results[i].VerifierID != id || results[i].Outcome != "passed" {
			return fmt.Errorf("required verifier is missing or failed: %s", id)
		}
	}
	return nil
}
