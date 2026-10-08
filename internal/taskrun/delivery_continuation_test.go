package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func continuationTestInput(t *testing.T, f deliveryFixture, o WorkflowRun) (Dependencies, DeliveryContinuationInput, DeliveryContinuationPreview) {
	t.Helper()
	source := filepath.Join(f.R.WorkspaceRoot, "compatible-installed-ply")
	if err := os.WriteFile(source, []byte("#!/bin/sh\n# compatible reader fixture\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	d := f.D
	d.Executable = func() (string, error) { return source, nil }
	in := DeliveryContinuationInput{RunID: o.RunID, ContextPath: o.Paths.Context, ControlExecutable: Executable{source, hashFileTest(t, source)}}
	p, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in)
	if err != nil {
		t.Fatal(err)
	}
	if err = privateDir(p.Directory); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(p.ControlExecutable.Path, raw, 0500); err != nil {
		t.Fatal(err)
	}
	return d, in, p
}

func TestDeliveryContinuationPreservesAcceptedSessionAndFrozenInputs(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	prompts := len(workflowProviderCalls(t, f.workflowFixture, "agent prompt"))
	frozen := map[string][]byte{}
	for _, path := range []string{workflowIndex(f.R.WorkspaceRoot, o.RunID), o.Paths.Context, f.R.Runtime.PlyExecutable.Path, filepath.Join(o.Paths.RunRoot, "acceptance.json"), filepath.Join(o.Paths.RunRoot, "mandate.json")} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		frozen[path] = raw
	}
	d, in, p := continuationTestInput(t, f, o)
	result, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || result.State != "continued" || result.Context.Locator != o.Paths.Context {
		t.Fatalf("same-session continuation: %+v %v", result, err)
	}
	for path, want := range frozen {
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("frozen artifact changed: %s %v", path, err)
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Continuation == nil || s.Result.SessionID != o.SessionID || s.Result.Transport.AgentSessionID != o.Transport.AgentSessionID {
		t.Fatalf("ownership changed: %+v %v", s.Result.Transport, err)
	}
	if err = workflowCallback(f.D, s, o.Paths.Context); err == nil {
		t.Fatal("old control retained callback mutation authority after explicit transition")
	}
	d.Executable = func() (string, error) { return result.ControlExecutable.Path, nil }
	if err = deliveryCallback(d, s, o.Paths.Context, true); err != nil {
		t.Fatalf("new preserved control cannot use unchanged context: %v", err)
	}
	if calls := workflowProviderCalls(t, f.workflowFixture, "agent start"); len(calls) != 1 {
		t.Fatalf("continuation restarted provider: %d", len(calls))
	}
	if calls := workflowProviderCalls(t, f.workflowFixture, "agent prompt"); len(calls) != prompts {
		t.Fatalf("continuation sent another prompt: %d", len(calls))
	}
	if _, err = WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("continuation supplied missing human QA authority")
	}
}

func TestDeliveryContinuationConcurrentAndInterruptedPublication(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	fault := d
	fault.Fault = func(point string) error {
		if point == "delivery_after_continuation_proof" {
			return errors.New("simulated lost publication reply")
		}
		return nil
	}
	if _, err := WorkflowContinueDelivery(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fault did not interrupt publication")
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation)
			if err == nil && result.State != "continued" && result.State != "existing" {
				err = errors.New("unexpected continuation state")
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || deliveryContinuationGeneration(s) != 1 {
		t.Fatalf("duplicate continuation generation: %d %v", deliveryContinuationGeneration(s), err)
	}
	entries, err := os.ReadDir(filepath.Dir(p.Directory))
	if err != nil || len(entries) != 1 {
		t.Fatalf("duplicate continuation slots: %v %v", entries, err)
	}
}

func TestDeliveryContinuationRejectsUnknownVerificationAndStaleSession(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, _ := continuationTestInput(t, f, o)
	model := workflowProviderDocument(t, f.Model)
	model["agent_session_id"] = "different-session"
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), model)
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil || !strings.Contains(err.Error(), "delivery_continuation_session") {
		t.Fatalf("stale session granted continuation: %v", err)
	}
	model["agent_session_id"] = o.Transport.AgentSessionID
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), model)
	deliveryTestCandidate(t, f, "unknown command")
	fault := f.D
	fault.Fault = func(point string) error {
		if point == "delivery_after_verification_reservation" {
			return errors.New("interrupted before known execution")
		}
		return nil
	}
	if _, err := WorkflowDeliveryVerify(fault, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "unknown")); err == nil {
		t.Fatal("fault did not reserve unknown command")
	}
	if _, err := WorkflowPreviewDeliveryContinuation(d, f.R.WorkspaceRoot, in); err == nil || !strings.Contains(err.Error(), "delivery_effect_unknown") {
		t.Fatalf("unknown effect was replayable: %v", err)
	}
}

func TestDeliveryContinuationRejectsTamperedPriorControlAndAuthority(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	d, in, p := continuationTestInput(t, f, o)
	if _, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.R.Runtime.PlyExecutable.Path, []byte("#!/bin/sh\nexit 9\n"), 0700); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = workflowEffectiveRuntime(s); err == nil || !strings.Contains(err.Error(), "delivery_continuation_integrity") {
		t.Fatalf("tampered prior executable was ignored: %v", err)
	}
	s.Request.SchemaVersion = 9
	if err = deliveryContinuationAuthority(s); err == nil || !strings.Contains(err.Error(), "delivery_continuation_unsupported") {
		t.Fatalf("unsupported contract granted continuation: %v", err)
	}
}
