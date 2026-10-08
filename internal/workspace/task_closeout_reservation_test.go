package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTaskCloseoutReservationProtectsPendingInstallationAndStartup(t *testing.T) {
	f, in := closeoutFixture(t, true, true)
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview: %+v %v", p, err)
	}
	// Native startup takes the same WorkItems lock to validate and reserve its
	// source. Queue it while the closeout owner guard owns that lock: it must
	// observe the published reservation before any new writer can be started.
	startup := make(chan error, 1)
	guardCalls := 0
	guard := func() error {
		guardCalls++
		if _, err := f.dependencies.WorkItems.Snapshot(f.root); err != nil {
			return err
		}
		go func() {
			startup <- WithTaskSpecSnapshot(f.dependencies, f.root, func(s *TaskSpecSession) error {
				_, err := s.ValidateTarget(f.taskPath, p.Plan.Source.Ref, p.Plan.Source.GitCommonDir, nil)
				return err
			})
		}()
		return nil
	}
	r, err := ReserveTaskCloseout(f.dependencies, in, p.PlanSHA256, "installation_pending", guard)
	if err != nil || guardCalls != 1 || r.Retention != nil || r.WorktreeRemoved || r.LifecycleCompleted || r.NativeIntegrationID == "" || r.State != "cleanup_pending" {
		t.Fatalf("reservation invented effects: %+v %v", r, err)
	}
	select {
	case err := <-startup:
		if err == nil || !strings.Contains(err.Error(), "task_closeout_reserved") {
			t.Fatalf("new startup acquired reserved source: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("startup lock did not release")
	}
	for _, phase := range []string{"installation_pending", "install_failed", "install_unknown"} {
		if _, err = SetTaskCloseoutPhase(f.dependencies, in.TaskID, in.OperationID, p.PlanSHA256, phase); err != nil {
			t.Fatal(err)
		}
		show, err := ShowTask(f.dependencies, "task")
		if err != nil || show.Closeout == nil || show.Closeout.Phase != phase || show.ReadyForHandoff || show.Integration == nil || show.Integration.NextAction.Kind == "none" || !strings.Contains(strings.Join(show.Reasons, ","), phase) {
			t.Fatalf("installation followup disappeared: %+v %v", show, err)
		}
		if _, err := ApplyTaskCloseoutWithOwnershipGuard(f.dependencies, in, p.PlanSHA256, func() error { return nil }); err == nil || !strings.Contains(err.Error(), phase) {
			t.Fatalf("cleanup bypassed required installation %s: %v", phase, err)
		}
		if _, err := os.Stat(f.taskPath); err != nil {
			t.Fatal("pending installation source removed", err)
		}
	}
	if _, err := SetTaskCloseoutPhase(f.dependencies, in.TaskID, in.OperationID, p.PlanSHA256, "awaiting_closeout"); err != nil {
		t.Fatal(err)
	}
	if _, err := ApplyTaskCloseoutWithOwnershipGuard(f.dependencies, in, p.PlanSHA256, func() error { return errors.New("owner_active: a new writer appeared") }); err == nil || !strings.Contains(err.Error(), "owner_active") {
		t.Fatalf("cleanup bypassed native owner observation: %v", err)
	}
	if _, err := os.Stat(f.taskPath); err != nil {
		t.Fatal("failed owner check removed source", err)
	}
	r, err = ApplyTaskCloseoutWithOwnershipGuard(f.dependencies, in, p.PlanSHA256, func() error { return nil })
	if err != nil || r.State != "complete" {
		t.Fatalf("resume: %+v %v", r, err)
	}
	// An unrelated target remains available even when this Task's historical
	// path is intentionally absent. A retired resource is not an unknown path.
	if err := WithTaskSpecSnapshot(f.dependencies, f.root, func(s *TaskSpecSession) error {
		_, err := s.ValidateTarget(filepath.Join(f.wrapper, "main"), "refs/heads/main", p.Plan.Source.GitCommonDir, nil)
		return err
	}); err != nil {
		t.Fatalf("retired Task blocked unrelated source: %v", err)
	}
}

func TestTaskCloseoutReservationCannotClaimUnobservedIntegrationOrUnknownOwner(t *testing.T) {
	f, in := closeoutFixture(t, false, true)
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ReserveTaskCloseout(f.dependencies, in, p.PlanSHA256, "installation_pending", func() error { return nil }); err == nil || !strings.Contains(err.Error(), "integration_not_observed") {
		t.Fatalf("reservation invented integration: %v", err)
	}
	closeoutIntegrateFixture(t, f)
	if _, err := ReserveTaskCloseout(f.dependencies, in, p.PlanSHA256, "installation_pending", func() error { return errors.New("owner_active") }); err == nil || !strings.Contains(err.Error(), "owner_active") {
		t.Fatalf("reservation ignored owner: %v", err)
	}
	if r, err := ReadTaskCloseoutAt(f.root, in.TaskID); err != nil || r != nil {
		t.Fatalf("blocked reservation persisted: %+v %v", r, err)
	}
}
