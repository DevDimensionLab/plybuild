package integration

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
)

func TestServiceExplicitNotificationRetryPreservesIntegrationAndInstallFailure(t *testing.T) {
	route, _, d := integrationNotificationFixture(t)
	posts := 0
	d.Transport = func(context.Context, string, []byte) notification.Observation {
		posts++
		if posts == 1 {
			return notification.Observation{Status: 400, Body: []byte("invalid_payload"), Dispatch: "started"}
		}
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	s := NewService(taskrun.Dependencies{}, Options{Notifier: &NativeNotificationAdapter{Dependencies: d}})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
	r.Plan.NotificationChoice = &route
	r.State, r.InstallationStarted = "install_failed", true
	r.Installation = &InstallObservation{State: "install_failed", Reason: "The selected installer exited unsuccessfully."}
	s.event(&r, "integration.observed")
	if err := s.notify(&r); err != nil {
		t.Fatal(err)
	}
	before := notificationTree(t, filepath.Dir(route.Path))
	p, err := s.notificationRetryPreview(r)
	if err != nil || !p.Allowed || p.Confirmation == "" || p.ChannelLabel != route.Route.ChannelLabel || p.Observation.State != "rejected" {
		t.Fatalf("known rejection has no precise public retry: %+v %v", p, err)
	}
	if after := notificationTree(t, filepath.Dir(route.Path)); !reflect.DeepEqual(before, after) {
		t.Fatal("retry preview wrote transport state")
	}
	if err = s.retryNotification(&r, "sha256:"+strings.Repeat("0", 64)); err == nil {
		t.Fatal("unconfirmed retry succeeded")
	}
	if err = s.retryNotification(&r, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if r.NotificationState != "transport_acknowledged" || posts != 2 || len(r.Notification.Attempts) != 2 || r.State != "install_failed" || !r.Integrated || r.Installation.State != "install_failed" || r.Closeout != nil {
		t.Fatalf("notification retry altered integration or repeated transport: state=%s notification=%+v posts=%d", r.State, r.Notification, posts)
	}
	if err = s.retryNotification(&r, p.Confirmation); err != nil || posts != 2 {
		t.Fatalf("same explicit retry was replayed: %v posts=%d", err, posts)
	}
	if next, err := s.notificationRetryPreview(r); err != nil || next.Allowed {
		t.Fatalf("acknowledged delivery offered another retry: %+v %v", next, err)
	}
}

func TestServiceNotificationRetryNeverRepeatsUnknownTransport(t *testing.T) {
	for _, point := range []string{"initial-unknown", "interrupted-explicit-retry"} {
		t.Run(point, func(t *testing.T) {
			route, _, d := integrationNotificationFixture(t)
			posts := 0
			d.Transport = func(context.Context, string, []byte) notification.Observation {
				posts++
				if point == "initial-unknown" {
					return notification.Observation{Dispatch: "unknown", Failed: true}
				}
				if posts == 1 {
					return notification.Observation{Status: 400, Body: []byte("invalid_payload"), Dispatch: "started"}
				}
				return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
			}
			a := &NativeNotificationAdapter{Dependencies: d}
			s := NewService(taskrun.Dependencies{}, Options{Notifier: a})
			r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
			r.Plan.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
			s.event(&r, "integration.observed")
			if err := s.notify(&r); err != nil {
				t.Fatal(err)
			}
			if point == "interrupted-explicit-retry" {
				p, err := s.notificationRetryPreview(r)
				if err != nil || !p.Allowed {
					t.Fatalf("initial known rejection: %+v %v", p, err)
				}
				a.Dependencies.Fault = func(phase string) error {
					if phase == "integration_after_send" {
						return errors.New("injected lost retry response")
					}
					return nil
				}
				if err = s.retryNotification(&r, p.Confirmation); err == nil {
					t.Fatal("lost retry response reported acknowledged")
				}
				a.Dependencies.Fault = nil
				if err = s.retryNotification(&r, p.Confirmation); err != nil {
					t.Fatal(err)
				}
			}
			before := posts
			p, err := s.notificationRetryPreview(r)
			if err != nil || p.Allowed || p.Observation.State != "unknown" {
				t.Fatalf("unknown send offers another retry: %+v %v", p, err)
			}
			if err = s.retryNotification(&r, "sha256:"+strings.Repeat("e", 64)); err == nil || posts != before {
				t.Fatalf("unknown send retried: %v posts=%d", err, posts)
			}
			if r.State != "completed" || !r.Integrated {
				t.Fatal("unknown notification changed the integration result")
			}
		})
	}
}
