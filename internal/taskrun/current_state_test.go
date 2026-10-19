//go:build !windows

package taskrun

import (
	"testing"
	"time"
)

// The controller has consumed resume, but the persisted human guidance can
// still describe the earlier parking boundary. Its wording is presentation.
func TestAR01ConsumedResumeIgnoresDisplayWording(t *testing.T) {
	f := newAW01Fixture(t)
	f.result(&Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, nil)
	f.park("pause")
	f.own()
	if err := workerConsumeResume(f.d, f.s, f.control("resume")); err != nil {
		t.Fatal(err)
	}
	f.reload()
	f.s.Result.NextAction.Message = "The preserved work may continue."
	f.save()
	f.reads(aw01AssertObserving)
}

func TestAR01ProviderWaitIgnoresReasonWording(t *testing.T) {
	f := newAW01Fixture(t)
	deadline := f.now.Add(time.Hour).Format(time.RFC3339Nano)
	r := f.result(&Process{State: "exited", ExitCode: ptr(1), Quiescence: ptr(true)}, &WorkerProviderFailure{Kind: "usage_limit", ResetAt: deadline})
	f.own()
	if err := workerWaitForProvider(f.d, f.s, r); err != nil {
		t.Fatal(err)
	}
	f.reload()
	f.s.Result.Worker.Reason = "The selected provider is waiting for its reset."
	f.s.Result.NextAction.Message = "The visible wait wording has changed."
	f.save()
	f.control("pause")
	f.reads(func(t *testing.T, o WorkflowRun) {
		if o.NextAction.Actor != "worker" || o.Worker.ControlAction != "pause" {
			t.Fatalf("wording prevented observation of fresh control: %+v", o)
		}
	})
}
