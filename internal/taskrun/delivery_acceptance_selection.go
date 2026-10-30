package taskrun

import "github.com/devdimensionlab/plybuild/internal/workspace"

func deliveryEffectiveAgreement(s workflowState) *workspace.DeliveryAgreement {
	if s.Request.Delivery == nil {
		return nil
	}
	return s.Request.Delivery.Agreement
}
