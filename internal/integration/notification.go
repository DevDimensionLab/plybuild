package integration

import notification "github.com/devdimensionlab/plybuild/internal/workflownotification"

type NotificationRoute = notification.IntegrationRouteBinding
type NotificationEvent = notification.IntegrationEvent
type NotificationObservation = notification.IntegrationObservation

type NotificationAdapter interface {
	Send(NotificationRoute, NotificationEvent) (NotificationObservation, error)
	Observe(NotificationRoute, NotificationEvent) (NotificationObservation, error)
}

type NativeNotificationAdapter struct {
	Dependencies notification.Dependencies
}

func LoadNotificationRoute(path string) (NotificationRoute, error) {
	return notification.LoadIntegrationRoute(path)
}

func (a *NativeNotificationAdapter) dependencies() notification.Dependencies {
	d := a.Dependencies
	system := notification.SystemDependencies()
	if d.Now == nil {
		d.Now = system.Now
	}
	if d.LookupEnv == nil {
		d.LookupEnv = system.LookupEnv
	}
	if d.Transport == nil {
		d.Transport = system.Transport
	}
	return d
}

func (a *NativeNotificationAdapter) Observe(route NotificationRoute, event NotificationEvent) (NotificationObservation, error) {
	return notification.ObserveIntegration(route, event)
}

func (a *NativeNotificationAdapter) Send(route NotificationRoute, event NotificationEvent) (NotificationObservation, error) {
	return notification.SendIntegration(a.dependencies(), route, event)
}

func (a *NativeNotificationAdapter) RetryConfirmation(route NotificationRoute, event NotificationEvent) (string, error) {
	return notification.IntegrationRetryConfirmation(a.dependencies(), route, event)
}

func (a *NativeNotificationAdapter) Retry(route NotificationRoute, event NotificationEvent, confirmation string) (NotificationObservation, error) {
	return notification.RetryIntegration(a.dependencies(), route, event, confirmation)
}
