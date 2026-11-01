package workflowhandoff

// ValidateWorkflowRound validates the frozen Spec's required evidence before a
// transport reserves a report. It reuses native semantics without publication.
// Only a complete report must show every requirement passed; a non-complete
// report keeps the exact binding with its actual failed/not_run/unknown outcomes.
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
	if draft.Outcome != "complete" {
		return incompleteRequirementsCoverage(d, s, b, eval)
	}
	return taskRequirementsCoverage(d, s, b, eval)
}
