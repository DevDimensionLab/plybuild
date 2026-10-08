package integration

import "fmt"

type controlledNotificationRetry interface {
	RetryConfirmation(NotificationRoute, NotificationEvent) (string, error)
	Retry(NotificationRoute, NotificationEvent, string) (NotificationObservation, error)
}

type NotificationRetryPreview struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	ID            string                   `json:"id"`
	EventID       string                   `json:"event_id"`
	RouteName     string                   `json:"route_name"`
	ChannelLabel  string                   `json:"channel_label"`
	Allowed       bool                     `json:"allowed"`
	Confirmation  string                   `json:"confirmation,omitempty"`
	Observation   *NotificationObservation `json:"observation,omitempty"`
	Reason        string                   `json:"reason,omitempty"`
	NextAction    string                   `json:"next_action"`
}

func notificationRetryBinding(r Receipt) (NotificationRoute, NotificationEvent, error) {
	if !r.Integrated || !r.NotificationStarted || r.Plan.NotificationRoute == nil || r.NotificationEvent == nil {
		return NotificationRoute{}, NotificationEvent{}, fmt.Errorf("retry_forbidden: this integration has no preserved notification attempt")
	}
	e := *r.NotificationEvent
	if e.IntegrationID != r.ID || e.TaskID != string(r.Plan.TaskResult.TaskID) || e.DeliveryID != r.Plan.DeliveryID || e.CandidateOID != r.Plan.TaskResult.ResultOID || e.TargetRef != r.Plan.Agreement.TargetRef || e.IntegratedOID != r.IntegratedOID {
		return NotificationRoute{}, e, fmt.Errorf("event_changed: the notification differs from the preserved integration effect")
	}
	route, err := LoadNotificationRoute(r.Plan.NotificationRoute.Locator)
	if err != nil {
		return route, e, err
	}
	if route.SHA256 != r.Plan.NotificationRoute.SHA256 || r.Plan.NotificationChoice != nil && digest(route) != digest(*r.Plan.NotificationChoice) {
		return route, e, fmt.Errorf("plan_changed: the originally selected notification route changed")
	}
	return route, e, nil
}

func (s *Service) notificationRetryPreview(r Receipt) (NotificationRetryPreview, error) {
	p := NotificationRetryPreview{Kind: "ply.integration.notification-retry-preview", SchemaVersion: 1, ID: r.ID, NextAction: "Inspect the notification outcome; unknown or acknowledged sends cannot be retried."}
	route, event, err := notificationRetryBinding(r)
	if err != nil {
		p.Reason = err.Error()
		return p, nil
	}
	p.EventID, p.RouteName, p.ChannelLabel = event.ID, route.Route.Name, route.Route.ChannelLabel
	o, err := s.options.Notifier.Observe(route, event)
	p.Observation = &o
	if err != nil {
		p.Reason = err.Error()
		return p, nil
	}
	a, ok := s.options.Notifier.(controlledNotificationRetry)
	if !ok {
		p.Reason = "retry_forbidden: selected notification transport has no controlled retry"
		return p, nil
	}
	p.Confirmation, err = a.RetryConfirmation(route, event)
	if err != nil {
		p.Reason = err.Error()
		return p, nil
	}
	p.Allowed = true
	p.NextAction = "Explicitly retry only this proven undelivered notification with --apply --confirm " + p.Confirmation + "; integration, installation and cleanup will not run."
	return p, nil
}

func (s *Service) PreviewNotificationRetry(cwd, id string) (NotificationRetryPreview, error) {
	r, err := s.Read(cwd, id)
	if err != nil {
		return NotificationRetryPreview{}, err
	}
	return s.notificationRetryPreview(r)
}

func (s *Service) RetryNotification(cwd, id, confirmation string) (Receipt, error) {
	_, root, err := s.context(cwd)
	if err != nil {
		return Receipt{}, err
	}
	var r Receipt
	err = withLock(root, func() error {
		var err error
		r, err = readReceipt(root, id)
		if err != nil {
			return err
		}
		return s.retryNotification(&r, confirmation)
	})
	return r, err
}

// The native transport journals every explicit retry before dispatch. This
// controller never invokes the integration/install/cleanup execution path.
func (s *Service) retryNotification(r *Receipt, confirmation string) error {
	if !sha256Pattern.MatchString(confirmation) {
		return fmt.Errorf("confirmation_required: inspect the exact current notification retry first")
	}
	route, event, err := notificationRetryBinding(*r)
	if err != nil {
		return err
	}
	o, err := s.options.Notifier.Observe(route, event)
	if err != nil {
		return err
	}
	for _, prior := range o.Attempts {
		if prior.Confirmation == confirmation {
			r.Notification, r.NotificationState = &o, o.State
			return s.save(r)
		}
	}
	a, ok := s.options.Notifier.(controlledNotificationRetry)
	if !ok {
		return fmt.Errorf("retry_forbidden: selected notification transport has no controlled retry")
	}
	expected, err := a.RetryConfirmation(route, event)
	if err != nil {
		return err
	}
	if confirmation != expected {
		return fmt.Errorf("confirmation_changed: inspect the current event, route and exact non-delivery before retry")
	}
	o, err = a.Retry(route, event, confirmation)
	r.Notification, r.NotificationState = &o, o.State
	if err != nil {
		r.NextAction = "Integration status is preserved; inspect the separately recorded notification retry: " + err.Error()
	} else if o.State != "transport_acknowledged" {
		r.NextAction = "Integration status is preserved. Inspect the notification outcome; only positively proven non-delivery allows another explicit retry."
	}
	if saveErr := s.save(r); saveErr != nil {
		return saveErr
	}
	return err
}
