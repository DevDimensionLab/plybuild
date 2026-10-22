package taskrun

// DeliveryResourceRetirementStatus reads recorded observations only. It is not
// a current source digest, process observation or effect-resolution authority.
type DeliveryResourceRetirementStatus struct {
	Mode                 *string                  `json:"mode"`
	State                string                   `json:"state"`
	NativeSession        string                   `json:"native_session"`
	ActiveStop           *FileBinding             `json:"active_stop"`
	Source               *ReportRetainedSource    `json:"source"`
	RetainedEffects      []DeliveryRetainedEffect `json:"retained_effects"`
	Retirement           *FileBinding             `json:"retirement"`
	RetainedObservations []FileBinding            `json:"retained_observations"`
	QueueDisposition     string                   `json:"queue_disposition"`
	NextAction           string                   `json:"next_action"`
}

func deliveryResourceStatus(d Dependencies, s workflowState) (*DeliveryResourceRetirementStatus, error) {
	out := &DeliveryResourceRetirementStatus{State: "not_retired", NativeSession: s.Result.Transport.AgentSessionID, ActiveStop: s.Result.Delivery.ActiveStop, RetainedEffects: []DeliveryRetainedEffect{}, RetainedObservations: []FileBinding{}, QueueDisposition: "unknown", NextAction: "Resource retirement requires an explicit supported coordinator operation."}
	r, b, observations, bindings, err := deliveryRetentionHistory(s)
	if err != nil {
		return out, err
	}
	if r == nil || r.Mode != deliveryRetainEffects {
		return out, nil
	}
	out.Mode = &r.Mode
	out.State = "recorded_pending_release"
	out.Source = r.DeliverySource
	out.RetainedEffects = r.DeliveryEffects
	out.Retirement = &b
	out.RetainedObservations = bindings
	if len(observations) > 0 {
		o := observations[len(observations)-1]
		out.Source = &o.Source
		out.RetainedEffects = o.Effects
		out.ActiveStop = o.ActiveStop
	}
	p, err := workflowReportPreparation(d, s)
	out.NextAction = retainedResourcesCaution + " Exact preparation release is pending; repeat the same retire command to revalidate current facts."
	if err != nil {
		return out, err
	}
	out.QueueDisposition = p.Disposition
	if p.Disposition == "released" || p.Disposition == "advanced" {
		out.State = "retired"
		out.NextAction = retainedResourcesCaution
	}
	return out, nil
}
