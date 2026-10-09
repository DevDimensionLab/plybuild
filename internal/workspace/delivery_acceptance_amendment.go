package workspace

import "fmt"

// DeliveryAcceptanceAmendment binds an explicit acceptance choice to one
// preserved Run and its original agreement. The decision bytes carry the
// user's provenance; this binding neither replaces them nor grants effects.
type DeliveryAcceptanceAmendment struct {
	SchemaVersion           int                      `yaml:"schema_version" json:"schema_version"`
	RunID                   string                   `yaml:"run_id" json:"run_id"`
	RequestSHA256           string                   `yaml:"request_sha256" json:"request_sha256"`
	MandateSHA256           string                   `yaml:"mandate_sha256" json:"mandate_sha256"`
	PreviousAgreementSHA256 string                   `yaml:"previous_agreement_sha256" json:"previous_agreement_sha256"`
	Policy                  DeliveryAcceptancePolicy `yaml:"policy" json:"policy"`
	DecisionLocator         string                   `yaml:"decision_locator" json:"decision_locator"`
	DecisionSHA256          string                   `yaml:"decision_sha256" json:"decision_sha256"`
}

// ApplyDeliveryAcceptanceAmendment derives an effective agreement without
// rewriting the original. Only a bound v1 local Epic agreement can opt in;
// every source, target, owner and allowed effect remains exactly the same.
func ApplyDeliveryAcceptanceAmendment(original DeliveryAgreement, a DeliveryAcceptanceAmendment) (DeliveryAgreement, error) {
	if err := ValidateDeliveryAgreement(original); err != nil {
		return DeliveryAgreement{}, err
	}
	if original.SchemaVersion != 1 || original.Mode != DeliveryLocalEpic || original.SourceRef == "" || original.IntegrationOwner != "" || original.Acceptance != nil {
		return DeliveryAgreement{}, fmt.Errorf("acceptance amendment requires an existing bound v1 local Epic agreement without a separate human integration owner")
	}
	if a.SchemaVersion != 1 || !validTaskText(a.RunID, 1, 256) || !digestPattern.MatchString(a.RequestSHA256) || !digestPattern.MatchString(a.MandateSHA256) || a.PreviousAgreementSHA256 != DeliveryAgreementDigest(original) || contentPath(a.DecisionLocator) != nil || !digestPattern.MatchString(a.DecisionSHA256) {
		return DeliveryAgreement{}, fmt.Errorf("acceptance amendment lacks its exact previous agreement, native Run/request/mandate or preserved decision")
	}
	if err := validateDeliveryAcceptancePolicy(a.Policy); err != nil {
		return DeliveryAgreement{}, err
	}
	effective := original
	effective.SchemaVersion = 3
	policy := a.Policy
	effective.Acceptance = &policy
	return effective, ValidateDeliveryAgreement(effective)
}

func deliveryAuthorizedAgreement(original DeliveryAgreement, auth *DeliveryAuthorization) (DeliveryAgreement, error) {
	if auth == nil {
		return DeliveryAgreement{}, fmt.Errorf("native delivery authorization is missing")
	}
	effective := original
	if auth.AcceptanceAmendment != nil {
		var err error
		effective, err = ApplyDeliveryAcceptanceAmendment(original, *auth.AcceptanceAmendment)
		if err != nil {
			return DeliveryAgreement{}, err
		}
	}
	if !contentTypedEqual(effective, auth.Agreement) {
		return DeliveryAgreement{}, fmt.Errorf("delivery authority differs from the immutable agreement and explicit acceptance amendment")
	}
	return effective, nil
}
