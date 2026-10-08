package integration

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These tests exercise the effect journal below the native candidate gate.
// The actual CLI scenarios separately require real qualified native evidence.
func serviceEffectReceipt(t *testing.T, root string) Receipt {
	t.Helper()
	p := Plan{Kind: "ply.integration.plan", SchemaVersion: 1, Workspace: root, TaskResult: workspace.TaskResultRecord{TaskID: "task", ResultOID: strings.Repeat("a", 40)}, Agreement: workspace.DeliveryAgreement{TargetRef: "refs/heads/epic"}}
	h := digest(p)
	return Receipt{Kind: "ply.integration.receipt", SchemaVersion: 1, ID: operationID(h), Plan: p, PlanSHA256: h, State: "completed", Integrated: true, IntegratedOID: p.TaskResult.ResultOID, Events: []Event{}}
}

type uncertainInstaller struct{ install, observe int }

func (a *uncertainInstaller) Install(InstallProfile) (InstallObservation, error) {
	a.install++
	return InstallObservation{State: "effect_unknown", Reason: "lost response"}, errors.New("unknown process result")
}
func (a *uncertainInstaller) Observe(InstallProfile) (InstallObservation, error) {
	a.observe++
	return InstallObservation{State: "effect_unknown", Reason: "cannot verify artifact"}, nil
}

func TestServicePreservesUnknownInstallationAndNeverOffersBlindRetry(t *testing.T) {
	a := &uncertainInstaller{}
	s := NewService(taskrun.Dependencies{}, Options{Installer: a})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.Install = &InstallProfile{}
	if err := s.install(&r, false); err != nil {
		t.Fatal(err)
	}
	if r.State != "install_unknown" || !r.Integrated || strings.Contains(r.NextAction, "--retry-install") {
		t.Fatalf("uncertain installation was reported as safely retryable: state=%s next=%s", r.State, r.NextAction)
	}
	if err := s.install(&r, true); err != nil {
		t.Fatal(err)
	}
	if a.install != 1 || a.observe != 1 {
		t.Fatalf("unknown installation repeated: %+v", a)
	}
}

type failedInstaller struct{ calls int }

func (a *failedInstaller) Install(InstallProfile) (InstallObservation, error) {
	a.calls++
	if a.calls == 1 {
		return InstallObservation{State: "install_failed", Reason: "command exited nonzero"}, nil
	}
	return InstallObservation{State: "verified"}, nil
}
func (*failedInstaller) Observe(InstallProfile) (InstallObservation, error) {
	return InstallObservation{State: "effect_unknown", Reason: "expected artifact absent"}, errors.New("artifact does not exist")
}

func TestKnownFailedInstallSurvivesOrdinaryResumeBeforeExplicitRetry(t *testing.T) {
	a := &failedInstaller{}
	s := NewService(taskrun.Dependencies{}, Options{Installer: a})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.Install = &InstallProfile{}
	if e := s.install(&r, false); e != nil {
		t.Fatal(e)
	}
	if e := s.install(&r, false); e != nil {
		t.Fatal(e)
	}
	if r.State != "install_failed" || r.Installation.State != "install_failed" || a.calls != 1 {
		t.Fatal("observation discarded the known failed attempt and its explicit retry path")
	}
	if e := s.install(&r, true); e != nil {
		t.Fatal(e)
	}
	if r.Installation.State != "verified" || a.calls != 2 {
		t.Fatal("explicit retry did not use preserved known failure")
	}
}

func TestInterruptedInstallRetryCannotReusePreviousFailureAsNoEffectProof(t *testing.T) {
	a := &failedInstaller{}
	s := NewService(taskrun.Dependencies{}, Options{Installer: a, Fault: func(stage string) error {
		if stage == "after_install_effect" && a.calls == 2 {
			return errors.New("interrupted after retry dispatch")
		}
		return nil
	}})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.Install = &InstallProfile{}
	if e := s.install(&r, false); e != nil {
		t.Fatal(e)
	}
	if e := s.install(&r, true); e == nil {
		t.Fatal("retry interruption not exercised")
	}
	if e := s.install(&r, false); e != nil {
		t.Fatal(e)
	}
	if e := s.install(&r, true); e != nil {
		t.Fatal(e)
	}
	if a.calls != 2 || r.State != "install_unknown" {
		t.Fatalf("old failure allowed blind third attempt: calls=%d state=%s", a.calls, r.State)
	}
}

func TestServiceIntegrationEventReachesNativeNotificationOnce(t *testing.T) {
	route, _, d := integrationNotificationFixture(t)
	posts := 0
	d.Transport = func(context.Context, string, []byte) notification.Observation {
		posts++
		return notification.Observation{Status: 200, Body: []byte("ok"), Dispatch: "started"}
	}
	s := NewService(taskrun.Dependencies{}, Options{Notifier: &NativeNotificationAdapter{Dependencies: d}})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
	s.event(&r, "integration.observed")
	if e := s.notify(&r); e != nil {
		t.Fatal(e)
	}
	if r.NotificationState != "transport_acknowledged" || posts != 1 {
		t.Fatalf("integration event was rejected or never sent: state=%s next=%s posts=%d", r.NotificationState, r.NextAction, posts)
	}
	if e := s.notify(&r); e != nil {
		t.Fatal(e)
	}
	if posts != 1 {
		t.Fatalf("notification repeated %d times", posts)
	}
}

type interruptedNotifier struct{}

func (interruptedNotifier) Send(_ NotificationRoute, e NotificationEvent) (NotificationObservation, error) {
	return NotificationObservation{State: "unknown", EventID: e.ID, ID: "reserved-attempt"}, errors.New("interrupted after durable reservation")
}
func (interruptedNotifier) Observe(_ NotificationRoute, e NotificationEvent) (NotificationObservation, error) {
	return NotificationObservation{State: "unknown", EventID: e.ID, ID: "reserved-attempt"}, nil
}
func TestServiceKeepsNativeNotificationReservationWhenTransportReturnsError(t *testing.T) {
	route, _, _ := integrationNotificationFixture(t)
	s := NewService(taskrun.Dependencies{}, Options{Notifier: interruptedNotifier{}})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
	s.event(&r, "integration.observed")
	if e := s.notify(&r); e != nil {
		t.Fatal(e)
	}
	if r.Notification == nil || r.Notification.ID != "reserved-attempt" || !r.Integrated || r.State != "completed" {
		t.Fatalf("native error observation lost or integration status changed: %+v", r)
	}
}

func TestVerifiedInstallationResumesCloseoutWithoutReadingRemovedSource(t *testing.T) {
	s := NewService(taskrun.Dependencies{}, Options{})
	r := serviceEffectReceipt(t, adapterPhysicalTemp(t))
	r.Plan.Flow = "legacy_reconcile" // Journal-only fixture: no native QA fabrication.
	r.Plan.Install = &InstallProfile{}
	r.Installation = &InstallObservation{State: "verified"}
	r.Closeout = &workspace.TaskCloseoutReceipt{State: "complete"}
	d := taskrun.Dependencies{Executable: func() (string, error) {
		return "", errors.New("source executable no longer exists after successful cleanup")
	}}
	if err := s.execute(d, r.Plan.Workspace, &r, false); err != nil {
		t.Fatalf("completed installation prevented closeout recovery: %v", err)
	}
	if r.State != "completed" || !r.Integrated {
		t.Fatalf("lost completed effects: %s", r.State)
	}
}

func TestPlanTextShowsExactInstallCommandsAndNotificationDestination(t *testing.T) {
	route, _, _ := integrationNotificationFixture(t)
	p := Preview{Plan: Plan{Install: &InstallProfile{Name: "fixture", ArtifactPath: "/artifacts/result", ArtifactSHA256: "sha256:expected", Install: CommandSpec{Executable: "/tools/install", Arguments: []string{"--candidate", "approved"}, CWD: "/source/build"}, Verify: CommandSpec{Executable: "/tools/verify", Arguments: []string{"--offline"}, CWD: "/installed/check"}}, NotificationRoute: &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}, NotificationChoice: &route}}
	shown := PreviewText(p)
	for _, detail := range []string{"/tools/install", "--candidate", "approved", "/source/build", "/tools/verify", "--offline", "/installed/check", route.Route.ChannelLabel, route.Route.Name} {
		if !strings.Contains(shown, detail) {
			t.Fatalf("confirmation hides %q: %s", detail, shown)
		}
	}
}

func TestControlPlanRejectsNotificationStateOrInstalledArtifactInsideSource(t *testing.T) {
	for _, choice := range []string{"notification-state", "installed-artifact"} {
		t.Run(choice, func(t *testing.T) {
			p := Plan{TaskResult: workspace.TaskResultRecord{SourceLocator: "/workspace/task"}}
			if choice == "notification-state" {
				p.NotificationChoice = &NotificationRoute{StateRoot: "/workspace/task/notification-integration-events"}
			} else {
				p.Install = &InstallProfile{ArtifactPath: "/workspace/task/bin/product"}
			}
			out := Preview{State: "ready"}
			validateControlLocations(&out, p)
			if out.State != "blocked" {
				t.Fatal("cleanup may remove required durable integration output")
			}
		})
	}
}
