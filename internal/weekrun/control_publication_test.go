package weekrun

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func mr05Await(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("deterministic race participant did not reach its barrier")
		}
		time.Sleep(time.Millisecond)
	}
}

// Apply is held after the real parent's basis validation, with control.lock and
// the TaskRun store lock still owned. Serve must not publish a generation through
// this window, even when its context was cancelled before the first tick.
func mr05OwnerPublicationRace(t *testing.T) {
	f := newMR05Week(t)
	r := weekPreview(t, f, "pause", "owner-publication-0001")
	validated, release := make(chan struct{}), make(chan struct{})
	f.m.Fault = func(point string) error {
		if point == "control_before_reserve" {
			close(validated)
			<-release
		}
		return nil
	}
	applyDone := make(chan error, 1)
	go func() { _, err := f.m.Controls().Apply(r); applyDone <- err }()
	<-validated
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	serveDone := make(chan error, 1)
	go func() { serveDone <- f.m.Serve(ctx, f.c.Workspace, f.b.ID) }()
	mr05Await(t, func() bool {
		_, err := os.Stat(filepath.Join(f.dir, "owner.lock"))
		return err == nil
	})
	var early error
	returned := false
	select {
	case early = <-serveDone:
		returned = true
	case <-time.After(200 * time.Millisecond):
	}
	whileHeld, readErr := readState(f.b, f.dir)
	close(release)
	if err := <-applyDone; err != nil && taskrun.ControlFailureOf(err).Code != "control_outcome_unknown" {
		t.Fatal("apply lost its serialized basis", err)
	}
	if !returned {
		early = <-serveDone
	}
	if !errors.Is(early, context.Canceled) {
		t.Fatal("owner did not stop before dispatch", early)
	}
	if readErr != nil || whileHeld.Generation != f.s.Generation {
		t.Errorf("owner generation changed while approved basis was guarded: before=%d during=%d err=%v", f.s.Generation, whileHeld.Generation, readErr)
	}
	if len(f.n.starts) != 0 {
		t.Fatal("cancelled owner dispatched")
	}
	// This later owner was not caused by the accepted wake and must not silently
	// consume the old generation's approval.
	f.m.Fault = nil
	f.reload(t)
	_ = f.m.Tick(f.b, f.dir, &f.s)
	status, _ := f.m.Controls().Status(f.c.Workspace, r.RequestID)
	if status.ControlEffect == "paused" || status.ControlEffect == "stopped" || status.ControlEffect == "resumed" {
		t.Error("unlinked successor consumed the prior generation's approval")
	}
}

func mr05BoundWakeSuccessor(t *testing.T) {
	for _, restarts := range []int{1, 2} {
		t.Run(map[int]string{1: "exact_successor", 2: "later_generation_rejected"}[restarts], func(t *testing.T) {
			f := newMR05Week(t)
			r := weekPreview(t, f, "pause", "wake-successor-0001")
			weekApply(t, f, r)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			for i := 0; i < restarts; i++ {
				if err := f.m.Serve(ctx, f.c.Workspace, f.b.ID); !errors.Is(err, context.Canceled) {
					t.Fatal(err)
				}
			}
			f.reload(t)
			err := f.m.Tick(f.b, f.dir, &f.s)
			status, _ := f.m.Controls().Status(f.c.Workspace, r.RequestID)
			if restarts == 1 {
				if err != nil || status.ControlEffect != "paused" || status.Receipt == nil || status.Receipt.OwnerGeneration == nil || *status.Receipt.OwnerGeneration != 1 || status.Receipt.PredecessorGeneration == nil || *status.Receipt.PredecessorGeneration != 0 {
					t.Fatal("bound successor did not preserve predecessor proof", status, err)
				}
			} else if err == nil || status.ControlEffect == "paused" {
				t.Fatal("a second unlinked generation reused a prior wake", status, err)
			}
			if len(f.n.starts) != 0 {
				t.Fatal("wake test started a Goal")
			}
		})
	}
}

func mr05FailurePublicationRace(t *testing.T) {
	f := newMR05Week(t)
	dispatched, failDispatch := make(chan struct{}), make(chan struct{})
	validated, release := make(chan struct{}), make(chan struct{})
	f.m.Fault = func(point string) error {
		switch point {
		case "week_after_dispatch_reservation":
			close(dispatched)
			<-failDispatch
			return errors.New("injected dispatch failure before spawn")
		case "control_before_reserve":
			close(validated)
			<-release
		}
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serveDone := make(chan error, 1)
	go func() { serveDone <- f.m.Serve(ctx, f.c.Workspace, f.b.ID) }()
	<-dispatched
	r := weekPreview(t, f, "pause", "failure-publication-0001")
	before, err := readState(f.b, f.dir)
	if err != nil {
		t.Fatal(err)
	}
	applyDone := make(chan error, 1)
	go func() { _, err := f.m.Controls().Apply(r); applyDone <- err }()
	<-validated
	close(failDispatch)
	// The failure path runs while apply retains its publication barrier.
	deadline := time.Now().Add(200 * time.Millisecond)
	var held State
	for time.Now().Before(deadline) {
		held, err = readState(f.b, f.dir)
		if err != nil || held.Phase != before.Phase {
			break
		}
		time.Sleep(time.Millisecond)
	}
	close(release)
	applyErr := <-applyDone
	cancel()
	select {
	case <-serveDone:
	case <-time.After(3 * time.Second):
		t.Fatal("failure owner did not exit after release")
	}
	if applyErr != nil {
		t.Fatal("serialized apply unexpectedly failed", applyErr)
	}
	if err != nil || held.Phase != before.Phase || held.ControlRevision != before.ControlRevision {
		t.Errorf("failure publication changed parent under validated approval: before=%s/%d during=%s/%d err=%v", before.Phase, before.ControlRevision, held.Phase, held.ControlRevision, err)
	}
	if len(f.n.starts) != 0 {
		t.Fatal("fault barrier spawned native work")
	}
	status, err := f.m.Controls().Status(f.c.Workspace, r.RequestID)
	if status.RequestDigest == nil || *status.RequestDigest != taskrun.ControlRevision(r) || status.Scope == nil || *status.Scope != r.Input.Target || status.RequestState != "accepted" && status.RequestState != "outcome_unknown" {
		t.Fatal("race lost the accepted identity", status, err)
	}
}
