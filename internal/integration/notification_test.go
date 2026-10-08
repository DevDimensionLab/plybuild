package integration

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
)

func integrationNotificationFixture(t *testing.T) (NotificationRoute, NotificationEvent, notification.Dependencies) {
	t.Helper()
	root := adapterPhysicalTemp(t)
	route := notification.Route{Kind: "ply.workflow.notification-route", SchemaVersion: 1, Name: "explicit-test", Transport: "slack_incoming_webhook", ChannelLabel: "Integration test", WebhookEnv: "PLY_TEST_INTEGRATION_WEBHOOK", StateRoot: filepath.Join(root, "native-notifications")}
	raw, _ := json.Marshal(route)
	path := filepath.Join(root, "route.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	binding, err := LoadNotificationRoute(path)
	if err != nil {
		t.Fatal(err)
	}
	event := NotificationEvent{Kind: "ply.integration.observed", SchemaVersion: 1, ID: "integration-event-exact-effect", IntegrationID: "int_test", TaskID: "task_test", DeliveryID: "dlv_test", CandidateOID: strings.Repeat("a", 40), TargetRef: "refs/heads/epic", IntegratedOID: strings.Repeat("c", 40), ObservedAtUTC: "2026-10-08T12:00:00Z", Remaining: []string{"install_pending", "cleanup_pending"}}
	d := notification.Dependencies{Now: func() time.Time { return time.Date(2026, 10, 8, 12, 1, 0, 0, time.UTC) }, LookupEnv: func(key string) (string, bool) {
		if key != route.WebhookEnv {
			t.Fatalf("unexpected credential lookup %s", key)
		}
		return "https://hooks.slack.com/services/T_TEST/B_TEST/local-fake-credential", true
	}, Transport: func(context.Context, string, []byte) notification.Observation {
		t.Fatal("test transport must be explicitly replaced")
		return notification.Observation{}
	}}
	return binding, event, d
}

func notificationTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if info.IsDir() {
			out[rel] = "directory"
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[rel] = string(raw)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestIntegrationNotificationPlanAndReadbackAreReadOnlyWithoutCredentials(t *testing.T) {
	route, event, _ := integrationNotificationFixture(t)
	root := filepath.Dir(route.Path)
	before := notificationTree(t, root)
	a := NativeNotificationAdapter{Dependencies: notification.Dependencies{LookupEnv: func(string) (string, bool) { t.Fatal("readback read credentials"); return "", false }, Transport: func(context.Context, string, []byte) notification.Observation {
		t.Fatal("readback used network")
		return notification.Observation{}
	}}}
	if _, err := LoadNotificationRoute(route.Path); err != nil {
		t.Fatal(err)
	}
	o, err := a.Observe(route, event)
	if err != nil || o.State != "not_sent" || len(o.Attempts) != 0 {
		t.Fatalf("readback: %+v %v", o, err)
	}
	if after := notificationTree(t, root); !reflect.DeepEqual(before, after) {
		t.Fatalf("planning/readback wrote files: before=%v after=%v", before, after)
	}
}

func TestIntegrationNotificationUsesStableEventAndNativeTransport(t *testing.T) {
	route, event, d := integrationNotificationFixture(t)
	posts := 0
	d.Transport = func(_ context.Context, _ string, raw []byte) notification.Observation {
		posts++
		var p notification.Payload
		if err := json.Unmarshal(raw, &p); err != nil {
			t.Fatal(err)
		}
		if p.Mrkdwn || p.Parse != "none" || p.UnfurlLinks || p.UnfurlMedia || !strings.Contains(p.Text, event.TaskID) || !strings.Contains(p.Text, event.CandidateOID) || !strings.Contains(p.Text, event.IntegratedOID) || !strings.Contains(p.Text, "install_pending, cleanup_pending") || strings.Contains(p.Text, "agent_finished") {
			t.Fatalf("misleading integration notification: %+v", p)
		}
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	a := NativeNotificationAdapter{Dependencies: d}
	first, err := a.Send(route, event)
	if err != nil || first.State != "transport_acknowledged" || first.EventID != event.ID || len(first.Attempts) != 1 {
		t.Fatalf("notification: %+v %v", first, err)
	}
	// Deduplicated Send and readback do not need a credential after an attempt.
	a.Dependencies.LookupEnv = func(string) (string, bool) { t.Fatal("duplicate send consulted credentials"); return "", false }
	for i := 0; i < 2; i++ {
		again, err := a.Send(route, event)
		if err != nil || again.ID != first.ID || again.State != first.State {
			t.Fatalf("duplicate: %+v %v", again, err)
		}
	}
	if posts != 1 {
		t.Fatalf("stable event sent %d times", posts)
	}
	if _, err = os.Stat(route.Route.StateRoot); !os.IsNotExist(err) {
		t.Fatal("integration transport touched the existing agent/gate namespace")
	}
}

func TestIntegrationNotificationCrashAndTimeoutNeverBlindRetry(t *testing.T) {
	for _, point := range []string{"integration_after_reservation", "integration_after_send", "transport-unknown"} {
		t.Run(point, func(t *testing.T) {
			route, event, d := integrationNotificationFixture(t)
			posts := 0
			d.Transport = func(context.Context, string, []byte) notification.Observation {
				posts++
				if point == "transport-unknown" {
					return notification.Observation{Dispatch: "unknown", Failed: true}
				}
				return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
			}
			d.Fault = func(name string) error {
				if point == name {
					return errors.New("injected process stop")
				}
				return nil
			}
			a := NativeNotificationAdapter{Dependencies: d}
			first, _ := a.Send(route, event)
			if first.State != "unknown" {
				t.Fatalf("interrupted send has invented outcome: %+v", first)
			}
			a.Dependencies.Fault = nil
			for i := 0; i < 2; i++ {
				o, err := a.Send(route, event)
				if err != nil || o.State != "unknown" {
					t.Fatalf("unknown readback: %+v %v", o, err)
				}
			}
			if _, err := a.RetryConfirmation(route, event); err == nil {
				t.Fatal("unknown send allowed retry")
			}
			want := 1
			if point == "integration_after_reservation" {
				want = 0
			}
			if posts != want {
				t.Fatalf("blind retry after %s: %d posts", point, posts)
			}
		})
	}
}

func TestIntegrationNotificationExplicitRetryOnlyForProvenNonDelivery(t *testing.T) {
	route, event, d := integrationNotificationFixture(t)
	posts := 0
	d.Transport = func(context.Context, string, []byte) notification.Observation {
		posts++
		if posts == 1 {
			return notification.Observation{Status: 400, Body: []byte("invalid_payload"), Dispatch: "started"}
		}
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	a := NativeNotificationAdapter{Dependencies: d}
	first, err := a.Send(route, event)
	if err != nil || first.State != "rejected" {
		t.Fatalf("rejection not retained: %+v %v", first, err)
	}
	if _, err = a.Send(route, event); err != nil || posts != 1 {
		t.Fatal("send repeated a rejection implicitly")
	}
	confirmation, err := a.RetryConfirmation(route, event)
	if err != nil {
		t.Fatal(err)
	}
	second, err := a.Retry(route, event, confirmation)
	if err != nil || second.State != "transport_acknowledged" || second.ID != first.ID || len(second.Attempts) != 2 || posts != 2 {
		t.Fatalf("explicit retry: %+v %v posts=%d", second, err, posts)
	}
	if _, err = a.Retry(route, event, confirmation); err != nil || posts != 2 {
		t.Fatal("replayed retry confirmation sent again")
	}
	if _, err = a.RetryConfirmation(route, event); err == nil {
		t.Fatal("acknowledged send can be retried")
	}
}

func TestIntegrationNotificationBindingChangesAndMissingReceiptBlockResend(t *testing.T) {
	for _, mode := range []string{"changed-event", "changed-route", "missing-state"} {
		t.Run(mode, func(t *testing.T) {
			route, event, d := integrationNotificationFixture(t)
			posts := 0
			d.Transport = func(context.Context, string, []byte) notification.Observation {
				posts++
				return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
			}
			a := NativeNotificationAdapter{Dependencies: d}
			if _, err := a.Send(route, event); err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "changed-event":
				event.Remaining = []string{"different-status"}
			case "changed-route":
				raw, err := os.ReadFile(route.Path)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(route.Path, append(raw, '\n'), 0600); err != nil {
					t.Fatal(err)
				}
			case "missing-state":
				if err := os.Remove(notification.IntegrationStatePath(route)); err != nil {
					t.Fatal(err)
				}
			}
			o, err := a.Send(route, event)
			if err == nil || o.State != "unknown" || posts != 1 {
				t.Fatalf("changed binding permitted a send: %+v %v posts=%d", o, err, posts)
			}
		})
	}
}

func TestIntegrationNotificationRejectsSyntheticPreEffectEvents(t *testing.T) {
	for _, kind := range []string{"agent_finished", "ply.integration.planned", "ply.integration.qa_pass", "ply.integration.remote_pending"} {
		route, event, d := integrationNotificationFixture(t)
		event.Kind = kind
		a := NativeNotificationAdapter{Dependencies: d}
		if _, err := a.Send(route, event); err == nil {
			t.Fatalf("non-integration event %q sent", kind)
		}
	}
}

func TestIntegrationNotificationConcurrentSendReservesOnlyOneAttempt(t *testing.T) {
	route, event, d := integrationNotificationFixture(t)
	var posts int32
	d.Transport = func(context.Context, string, []byte) notification.Observation {
		atomic.AddInt32(&posts, 1)
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	a := NativeNotificationAdapter{Dependencies: d}
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, _ = a.Send(route, event) }()
	}
	wg.Wait()
	o, err := a.Observe(route, event)
	if err != nil || o.State != "transport_acknowledged" || atomic.LoadInt32(&posts) != 1 || len(o.Attempts) != 1 {
		t.Fatalf("concurrent send duplicated or lost result: %+v %v posts=%d", o, err, posts)
	}
}
