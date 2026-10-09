package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryAcceptanceChoice records an explicit policy choice, never a QA
// judgment or a test result. The active owner can preserve a human's choice on
// their behalf. Actor and answer are local provenance, not authentication.
type DeliveryAcceptanceChoice struct {
	Kind                    string                             `json:"kind"`
	SchemaVersion           int                                `json:"schema_version"`
	RunID                   string                             `json:"run_id"`
	RequestSHA256           string                             `json:"request_sha256"`
	PreviousAgreementSHA256 string                             `json:"previous_agreement_sha256"`
	CandidateOID            string                             `json:"candidate_oid"`
	CandidateTree           string                             `json:"candidate_tree"`
	ActorClaim              string                             `json:"actor_claim"`
	Answer                  string                             `json:"answer"`
	RequestedAtUTC          string                             `json:"requested_at_utc"`
	Policy                  workspace.DeliveryAcceptancePolicy `json:"policy"`
}

type deliveryAcceptanceSelection struct {
	Kind                string      `json:"kind"`
	SchemaVersion       int         `json:"schema_version"`
	RunID               string      `json:"run_id"`
	RequestSHA256       string      `json:"request_sha256"`
	SessionID           string      `json:"session_id"`
	NativeSessionID     string      `json:"native_session_id"`
	OwnerClaim          string      `json:"owner_claim"`
	Context             FileBinding `json:"context"`
	Mandate             FileBinding `json:"mandate"`
	RuntimeAcceptance   FileBinding `json:"runtime_acceptance"`
	NativeStartSHA256   string      `json:"native_start_sha256"`
	ActualPolicySHA256  string      `json:"actual_policy_sha256"`
	Choice              FileBinding `json:"choice"`
	ObservedEventSHA256 *string     `json:"observed_event_sha256"`
	RecordedAtUTC       string      `json:"recorded_at_utc"`
}

func deliveryEffectiveAgreement(s workflowState) *workspace.DeliveryAgreement {
	if s.EffectiveDeliveryAgreement != nil {
		return s.EffectiveDeliveryAgreement
	}
	if s.Request.Delivery == nil {
		return nil
	}
	return s.Request.Delivery.Agreement
}

func deliverySelectionDirectory(s workflowState) string {
	return filepath.Join(s.Result.Paths.RunRoot, "delivery", "acceptance-selection")
}

func validateAcceptanceChoice(s workflowState, c DeliveryAcceptanceChoice) error {
	if c.Kind != "DeliveryAcceptanceChoice@1" || c.SchemaVersion != 1 || c.RunID != s.Result.RunID || c.RequestSHA256 != s.Result.RequestSHA256 || !plain(c.ActorClaim, 1, 256) || !utf8.ValidString(c.Answer) || len(c.Answer) > 16384 || strings.TrimSpace(c.Answer) == "" || strings.ContainsRune(c.Answer, 0) {
		return workflowError(4, "acceptance choice requires the exact run/request, actor and actual policy-choice answer")
	}
	if _, err := time.Parse(time.RFC3339Nano, c.RequestedAtUTC); err != nil {
		return workflowError(2, "acceptance choice requires an actual UTC timestamp")
	}
	if c.CandidateOID == "" || c.CandidateTree == "" || c.Policy.RequireHumanQA {
		return workflowError(2, "this transition selects automatic acceptance without a new human QA requirement and names the exact current source")
	}
	if s.Request.Delivery == nil || s.Request.Delivery.Agreement == nil || c.PreviousAgreementSHA256 != workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement) {
		return workflowError(4, "acceptance choice differs from the original frozen agreement")
	}
	return nil
}

func readAcceptanceSelection(s workflowState, binding FileBinding) (deliveryAcceptanceSelection, DeliveryAcceptanceChoice, workspace.DeliveryAcceptanceAmendment, error) {
	var r deliveryAcceptanceSelection
	var c DeliveryAcceptanceChoice
	var amendment workspace.DeliveryAcceptanceAmendment
	dir := deliverySelectionDirectory(s)
	if binding.Locator != filepath.Join(dir, "selection.json") {
		return r, c, amendment, workflowError(4, "acceptance selection is not its managed immutable record")
	}
	raw, err := workflowBound(binding, 1<<20)
	if err == nil {
		err = decode(raw, 1<<20, &r)
	}
	if err != nil {
		return r, c, amendment, err
	}
	if r.Kind != "PlyDeliveryAcceptanceSelection@1" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.SessionID != s.Result.SessionID || r.NativeSessionID != s.Result.Transport.AgentSessionID || s.Request.Delivery == nil || r.OwnerClaim != s.Request.Delivery.OwnerClaim || r.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || r.Mandate != (FileBinding{s.Result.Handoff.Locator, s.Result.Handoff.SHA256}) || s.Acceptance == nil || r.RuntimeAcceptance != *s.Acceptance || s.StartSHA256 == nil || r.NativeStartSHA256 != *s.StartSHA256 || s.Result.Delivery == nil || s.Result.Delivery.ActualPolicySHA256 == nil || r.ActualPolicySHA256 != *s.Result.Delivery.ActualPolicySHA256 || r.Choice.Locator != filepath.Join(dir, "choice.json") {
		return r, c, amendment, workflowError(4, "acceptance selection differs from the exact original owner, context, mandate or accepted runtime")
	}
	if _, err = time.Parse(time.RFC3339Nano, r.RecordedAtUTC); err != nil {
		return r, c, amendment, err
	}
	raw, err = workflowBound(r.Choice, 1<<20)
	if err == nil {
		err = decode(raw, 1<<20, &c)
	}
	if err == nil {
		err = validateAcceptanceChoice(s, c)
	}
	if err != nil {
		return r, c, amendment, err
	}
	amendment = workspace.DeliveryAcceptanceAmendment{SchemaVersion: 1, RunID: r.RunID, RequestSHA256: r.RequestSHA256, MandateSHA256: r.Mandate.SHA256, PreviousAgreementSHA256: c.PreviousAgreementSHA256, Policy: c.Policy, DecisionLocator: binding.Locator, DecisionSHA256: binding.SHA256}
	_, err = workspace.ApplyDeliveryAcceptanceAmendment(*s.Request.Delivery.Agreement, amendment)
	return r, c, amendment, err
}

func workflowLoadAcceptanceSelection(s *workflowState) error {
	s.EffectiveDeliveryAgreement, s.EffectiveAcceptanceAmendment = nil, nil
	if s.Result.Delivery == nil {
		return nil
	}
	if s.Result.Delivery.AcceptanceSelection == nil {
		for _, event := range s.Result.Delivery.Events {
			if event.Kind == "acceptance_selection" {
				return workflowError(4, "native acceptance selection event is orphaned from its current binding")
			}
		}
		return nil
	}
	b := *s.Result.Delivery.AcceptanceSelection
	r, _, a, err := readAcceptanceSelection(*s, b)
	if err != nil {
		return err
	}
	found := false
	basisSeen := r.ObservedEventSHA256 == nil
	for _, e := range s.Result.Delivery.Events {
		if e.Kind == "acceptance_selection" {
			if found || e.ID != "acceptance-selection" || e.Binding.SHA256 != b.SHA256 || !basisSeen {
				return workflowError(4, "acceptance selection differs from its preserved event history")
			}
			if _, err = workflowBound(e.Binding, 1<<20); err != nil {
				return err
			}
			found = true
		}
		if r.ObservedEventSHA256 != nil && e.Binding.SHA256 == *r.ObservedEventSHA256 {
			basisSeen = true
		}
	}
	if !found {
		return workflowError(4, "acceptance selection has no preserved native event")
	}
	effective, err := workspace.ApplyDeliveryAcceptanceAmendment(*s.Request.Delivery.Agreement, a)
	if err != nil {
		return err
	}
	s.EffectiveDeliveryAgreement, s.EffectiveAcceptanceAmendment = &effective, &a
	return nil
}

// WorkflowDeliverySelectAcceptance preserves a single bounded transition for
// the active local Epic owner. It neither executes acceptance nor qualifies the
// current candidate. A duplicate resumes publication of the identical record.
func WorkflowDeliverySelectAcceptance(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	raw, err := readFile(file, 1<<20, true)
	var choice DeliveryAcceptanceChoice
	if err == nil {
		err = decode(raw, 1<<20, &choice)
	}
	if err != nil {
		return WorkflowRun{}, err
	}
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		if err := deliveryCallback(d, *s, contextPath, true); err != nil {
			return err
		}
		if err := deliveryAcceptedAuthority(*s); err != nil {
			return err
		}
		if err := validateAcceptanceChoice(*s, choice); err != nil {
			return err
		}
		if existing := s.Result.Delivery.AcceptanceSelection; existing != nil {
			r, _, _, err := readAcceptanceSelection(*s, *existing)
			if err != nil {
				return err
			}
			if r.Choice.SHA256 != hash(raw) {
				return workflowError(4, "this Run already preserves a different acceptance choice")
			}
			return nil
		}
		if s.Result.Delivery.Phase == "integrating" || s.Result.Delivery.Phase == "closing" {
			return workflowError(4, "acceptance cannot change during an existing delivery effect")
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
			return workflowError(4, "acceptance choice is blocked by an unresolved delivery attempt; observe the same attempt without replay")
		}
		x, err := deliveryTarget(d, *s)
		if err != nil {
			return err
		}

		// Revalidate frozen publications and runtime, but permit changed source
		// to require its own new verification rather than borrow old QA.
		if _, err = deliveryAuthorityState(d, *s); err != nil {
			return err
		}
		dir := deliverySelectionDirectory(*s)
		b := FileBinding{Locator: filepath.Join(dir, "selection.json")}
		var r deliveryAcceptanceSelection
		if saved, readErr := readFile(b.Locator, 1<<20, true); readErr == nil {
			b.SHA256 = hash(saved)
			r, _, _, err = readAcceptanceSelection(*s, b)
			if err != nil {
				return err
			}
			if r.Choice.SHA256 != hash(raw) {
				return workflowError(4, "a different selection is preserved; resume the same pending choice")
			}
		} else if !os.IsNotExist(readErr) {
			return readErr
		} else {
			if x.OID != choice.CandidateOID || x.Tree != choice.CandidateTree {
				return workflowError(4, "acceptance choice names stale candidate bytes; no candidate is approved")
			}
			r = deliveryAcceptanceSelection{Kind: "PlyDeliveryAcceptanceSelection@1", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, NativeSessionID: s.Result.Transport.AgentSessionID, OwnerClaim: s.Request.Delivery.OwnerClaim, Context: FileBinding{s.Result.Paths.Context, s.ContextSHA256}, Mandate: FileBinding{s.Result.Handoff.Locator, s.Result.Handoff.SHA256}, RuntimeAcceptance: *s.Acceptance, NativeStartSHA256: *s.StartSHA256, ActualPolicySHA256: *s.Result.Delivery.ActualPolicySHA256, Choice: FileBinding{filepath.Join(dir, "choice.json"), hash(raw)}, ObservedEventSHA256: s.Result.Delivery.LastEventSHA256, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
			a := workspace.DeliveryAcceptanceAmendment{SchemaVersion: 1, RunID: id, RequestSHA256: r.RequestSHA256, MandateSHA256: r.Mandate.SHA256, PreviousAgreementSHA256: choice.PreviousAgreementSHA256, Policy: choice.Policy, DecisionLocator: b.Locator, DecisionSHA256: digest(r)}
			if _, err = workspace.ApplyDeliveryAcceptanceAmendment(*s.Request.Delivery.Agreement, a); err != nil {
				return err
			}
			if err = d.writeOnce(r.Choice.Locator, raw); err != nil {
				return err
			}
			b, err = workflowKeep(d, b.Locator, r)
			if err != nil {
				return err
			}
		}
		if err = d.fault("delivery_after_acceptance_selection_record"); err != nil {
			return err
		}
		if _, err = deliveryAppendEvent(d, s, "acceptance-selection", "acceptance_selection", r); err != nil {
			return err
		}
		s.Result.Delivery.AcceptanceSelection = &b
		if err = workflowLoadAcceptanceSelection(s); err != nil {
			return err
		}
		s.Result.Delivery.Phase, s.Result.Round.State = "working", "working"
		s.Result.NextAction = WorkflowAction{"recipient", "Automatic acceptance explicitly selected for this same Run. Execute verify with the private candidate binary and candidate-bound review; this choice is neither a test pass nor HumanQA."}
		return nil
	})
	return deliveryReadback(d, root, id, err)
}
