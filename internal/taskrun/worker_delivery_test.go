//go:build darwin

package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkerActualFailedPretestStartsBoundedCorrection(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		w := s.Result.Worker
		w.Stage = "pretest"
		dir := filepath.Join(workerDirectory(*s), "attempts", "pretest-0001")
		if err := privateDir(dir); err != nil {
			return err
		}
		in, err := workflowKeep(f.D, filepath.Join(dir, "intent.json"), workerAttemptIntent{RunID: o.RunID, JobID: "pretest-0001", Kind: "pretest", CandidateOID: w.CandidateOID, CandidateTree: w.CandidateTree, Argv: []string{"/bin/sh", "-c", "printf 'real failure'; exit 1"}, Environment: []string{}, TimeoutSeconds: 10})
		if err != nil {
			return err
		}
		w.Job = &WorkerJob{ID: "pretest-0001", Kind: "pretest", Directory: dir, IntentHash: in.SHA256}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	sh, err := filepath.EvalSymlinks("/bin/sh")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(sh)
	if err != nil {
		t.Fatal(err)
	}
	completion := ProviderCompletion{}
	p, runErr := workerOwnedExec(LaunchSpec{Executable: Executable{Path: sh, SHA256: hash(raw)}, Argv: []string{sh, "-c", "printf 'real failure'; exit 1"}, CWD: s.Observed.Target.WorktreeLocator, Timeout: 10 * time.Second, StreamsRoot: filepath.Join(s.Result.Worker.Job.Directory, "streams"), Completion: &completion, Environment: []string{}}, func(Process) error { return nil })
	if p.State != "exited" || p.ExitCode == nil || *p.ExitCode != 1 || !quiescent(p) || !completion.StreamComplete {
		t.Fatalf("not a positive failed-test exit: %+v %+v %v", p, completion, runErr)
	}
	r := workerAttemptResult{RunID: o.RunID, JobID: s.Result.Worker.Job.ID, IntentSHA256: s.Result.Worker.Job.IntentHash, FinishedAt: "2026-10-09T00:00:01Z", Exit: p.ExitCode, Process: &p, Error: workerAttemptError("pretest", runErr), Stdout: completion.Stdout, Stderr: completion.Stderr}
	if _, err = workflowKeep(f.D, filepath.Join(s.Result.Worker.Job.Directory, "result.json"), r); err != nil {
		t.Fatal(err)
	}
	if err = workerStep(f.D, s); err != nil {
		t.Fatalf("normal test failure blocked instead of correcting: %v", err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Worker.Stage != "developer" || s.Result.Worker.Iteration != 2 || !strings.Contains(s.Result.Worker.Feedback, "Acceptance failed") || s.Result.Worker.Job != nil {
		t.Fatalf("wrong correction: %+v", s.Result.Worker)
	}
	s.Result.Worker.Iteration = s.Request.Worker.MaxAttempts
	if err = workerCorrection(f.D, s, "same actual test failure"); err == nil {
		t.Fatal("ignored correction limit")
	}
}

// The provider and its JSONL below are deliberately synthetic unit-test input.
// Native registry, test subprocess, candidate, Git return and closeout are real.
// The separately invoked public-CLI journey exercises actual Codex sessions.
func newWorkerDeliveryFixture(t *testing.T, customize ...func(*deliveryFixture)) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic worker fixture", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a})
	f.R.Envelope = deliveryEnv("headless-run-request")
	f.R.Worker = &WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 60}
	f.R.Herdr = (WorkflowRequest{}).Herdr
	f.R.Runtime.Mode = "headless"
	f.R.Runtime.PermissionBinding.ProfileID = "ply-worker"
	f.R.HumanAuthority.StartSurface = "human_authorized_worker"
	f.R.Delivery.AllowLocalInstall = false
	for _, apply := range customize {
		apply(&f)
	}
	var err error
	f.R.HandoffDraft, err = workflowhandoff.BuildDeliveryHandoffDraft(f.D.Workflow, f.R.PreparationID, "synthetic human", f.R.Delivery.OwnerClaim, "fixture/delivery", f.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{AllowLocalInstall: false, Agreement: f.Prepared.Delivery})
	if err != nil {
		t.Fatal(err)
	}
	f.File = writeAny(t, f.R.WorkspaceRoot, "headless-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	lock, err := workerLock(filepath.Join(o.Paths.RunRoot, "worker", "owner.lock"), true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { lock.Close() })
	f.D.WorkerOwner = &workerOwner{root: f.R.WorkspaceRoot, runID: o.RunID, generation: 1, lock: lock}
	if err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Worker.Generation = 1
		s.Result.Worker.Phase = "running"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	acceptance := deliveryTestAcceptance(t, f, o)
	acceptance.RuntimeClaim.RuntimeID = ptr("ply-worker")
	acceptance.RuntimeClaim.NativeSessionID = ptr(o.Worker.SessionID)
	o, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "worker-test-acceptance.json", acceptance))
	if err != nil {
		t.Fatal(err)
	}
	deliveryTestCandidate(t, f, "worker native candidate")
	if err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		w := s.Result.Worker
		w.Stage = "reviewer"
		w.Iteration = 1
		w.DeveloperSessions = []string{"synthetic-developer-thread"}
		w.CandidateOID = gitOutput(t, s.Observed.Target.WorktreeLocator, "rev-parse", "HEAD")
		w.CandidateTree = gitOutput(t, s.Observed.Target.WorktreeLocator, "rev-parse", "HEAD^{tree}")
		w.CandidateBinary = f.R.Runtime.PlyExecutable.Path
		dir := filepath.Join(workerDirectory(*s), "attempts", "reviewer-0001")
		if err := privateDir(dir); err != nil {
			return err
		}
		intent, err := workflowKeep(f.D, filepath.Join(dir, "intent.json"), workerAttemptIntent{RunID: o.RunID, JobID: "reviewer-0001", Kind: "reviewer", Generation: 1, CandidateOID: w.CandidateOID, CandidateTree: w.CandidateTree, Argv: []string{"synthetic-provider-fixture"}, Environment: []string{}, TimeoutSeconds: 60})
		if err != nil {
			return err
		}
		w.Job = &WorkerJob{ID: "reviewer-0001", Kind: "reviewer", Directory: dir, IntentHash: intent.SHA256}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	dir := s.Result.Worker.Job.Directory
	stdout := []byte("{\"type\":\"thread.started\",\"thread_id\":\"synthetic-reviewer-thread\"}\n{\"type\":\"turn.started\"}\n{\"type\":\"turn.completed\",\"usage\":{\"input_tokens\":1,\"cached_input_tokens\":0,\"output_tokens\":1}}\n")
	if err = f.D.writeOnce(filepath.Join(dir, "stdout.txt"), stdout); err != nil {
		t.Fatal(err)
	}
	if err = f.D.writeOnce(filepath.Join(dir, "stderr.txt"), nil); err != nil {
		t.Fatal(err)
	}
	answer, err := workflowKeep(f.D, filepath.Join(dir, "answer.json"), workerReviewAnswer{Decision: "passed", Summary: "Synthetic reviewer fixture only", Findings: []string{}})
	if err != nil {
		t.Fatal(err)
	}
	result := workerAttemptResult{RunID: o.RunID, JobID: s.Result.Worker.Job.ID, IntentSHA256: s.Result.Worker.Job.IntentHash, StartedAt: "2026-10-09T00:00:00Z", FinishedAt: "2026-10-09T00:00:01Z", Exit: ptr(0), ThreadID: "synthetic-reviewer-thread", Process: &Process{State: "exited", ExitCode: ptr(0), Quiescence: ptr(true)}, Stdout: FileBinding{filepath.Join(dir, "stdout.txt"), hash(stdout)}, Stderr: FileBinding{filepath.Join(dir, "stderr.txt"), hash(nil)}, Answer: &answer}
	if _, err = workflowKeep(f.D, filepath.Join(dir, "result.json"), result); err != nil {
		t.Fatal(err)
	}
	if err = workerRecordReview(f.D, s, result); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	return f, o
}

func TestWorkerNativeReadbackAndReceiptRecovery(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if o.Transport.AgentSessionID != "" || strings.Contains(o.NextAction.Message, "same tab") {
		t.Fatalf("worker was projected as Herdr: %+v", o.NextAction)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	f.D.Fault = func(point string) error {
		if point == "delivery_before_verification_receipt" {
			return errors.New("synthetic interruption after actual test")
		}
		return nil
	}
	_, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary)
	if err == nil || !strings.Contains(err.Error(), "synthetic interruption") {
		t.Fatalf("did not reach actual-test interruption: %v", err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	spool := filepath.Join(s.Result.Delivery.Attempt.Path, "completed-result.json")
	// Recovery uses the existing native completion object, never a fresh test.
	before, _ := os.ReadDir(s.Result.Delivery.Attempt.Path)
	f.D.Fault = nil
	if err = workerRecoverNative(f.D, s); err != nil {
		t.Fatal(err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Worker.Stage != "integrate" || len(s.Result.Delivery.Candidates) != 1 || s.Result.Delivery.Candidates[0].HumanQA != nil {
		t.Fatalf("native verification recovery lost candidate: %+v", s.Result.Delivery)
	}
	t.Logf("preserved native attempt files=%d; spool family under %s", len(before), spool)
	workerTestReserveNative(t, f, o.RunID, "integrate")
	f.D.Fault = func(point string) error {
		if point == "automatic_closeout_after_release" {
			return errors.New("synthetic interruption after integration")
		}
		return nil
	}
	_, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err == nil {
		t.Fatal("closeout interruption missing")
	}
	candidate := s.Result.Delivery.Candidates[0]
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate.OID {
		t.Fatal("native integration did not occur before interruption")
	}
	f.D.Fault = nil
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	// A native guard can preserve its noncompleted result after the effect.
	// Resume must still use the existing native receipt recovery seam.
	stream := filepath.Join(s.Result.Worker.Job.Directory, "native-output.txt")
	if err = f.D.writeOnce(stream, nil); err != nil {
		t.Fatal(err)
	}
	result := workerAttemptResult{RunID: o.RunID, JobID: s.Result.Worker.Job.ID, IntentSHA256: s.Result.Worker.Job.IntentHash, FinishedAt: "2026-10-09T00:00:01Z", Error: "synthetic interruption after integration", NativeCompleted: false, Stdout: FileBinding{stream, hash(nil)}, Stderr: FileBinding{stream, hash(nil)}}
	if _, err = workflowKeep(f.D, filepath.Join(s.Result.Worker.Job.Directory, "result.json"), result); err != nil {
		t.Fatal(err)
	}
	if err = workerStep(f.D, s); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || o.Closeout == nil || o.Closeout.State != "complete" || o.Delivery.Phase != "completed" {
		t.Fatalf("worker native closeout: %+v %v", o.Closeout, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.HumanQARecords) != 0 || len(registry.IntegrationResults) != 1 {
		t.Fatal("recovery invented HumanQA or repeated integration")
	}
	inventory, err := ReadInventory(f.D, f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range inventory.Runs {
		if row.RunID == o.RunID {
			if row.Herdr != nil || row.Transport == nil || *row.Transport != "headless" || row.Freshness != "fresh" {
				t.Fatalf("invalid worker inventory: %+v", row)
			}
		}
	}
	trace, err := ReadTraceHistory(f.D, f.R.WorkspaceRoot, "task")
	if err != nil {
		t.Fatal(err)
	}
	for _, run := range trace.Runs {
		if run.RunID == o.RunID && len(run.Reasons) != 0 {
			t.Fatalf("worker trace rejected native evidence: %+v", run.Reasons)
		}
	}
}

func TestWorkerRecoversActualFailedTestSpoolIntoCorrection(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'actual native failure'\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	f.D.Fault = func(point string) error {
		if point == "delivery_before_verification_receipt" {
			return errors.New("synthetic interruption after failed native test")
		}
		return nil
	}
	_, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary)
	if err == nil || !strings.Contains(err.Error(), "synthetic interruption") {
		t.Fatalf("did not interrupt after native test: %v", err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	previousAttempt := s.Result.Delivery.Attempt.ID
	f.D.Fault = nil
	if err = workerRecoverNative(f.D, s); err != nil {
		t.Fatalf("actual failed-test spool blocked correction: %v", err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Delivery.Attempt.ID != previousAttempt || s.Result.Delivery.Attempt.State != "recorded" || s.Result.Worker.Iteration != 2 || s.Result.Worker.Stage != "developer" || !strings.Contains(s.Result.Worker.Feedback, "Native acceptance failed") {
		t.Fatalf("failure was replayed or did not request correction: %+v %+v", s.Result.Delivery.Attempt, s.Result.Worker)
	}
}

func TestWorkerResumeCannotBeLostDuringSupervisorShutdown(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Worker.Phase = "paused"; return nil }); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = workerControlWrite(f.D, s, "pause"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(workerDirectory(s), "supervisor.lock")
	lock, err := workerLock(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	decided, release, exited := make(chan struct{}), make(chan struct{}), make(chan struct{})
	released := false
	defer func() {
		if !released {
			close(release)
		}
		<-exited
	}()
	d := f.D
	d.Fault = func(point string) error {
		if point == "worker_supervisor_before_exit" {
			close(decided)
			<-release
		}
		return nil
	}
	go func() { defer close(exited); _, _ = workerSupervisorExit(d, f.R.WorkspaceRoot, o.RunID, lock) }()
	select {
	case <-decided:
	case <-time.After(5 * time.Second):
		t.Fatal("supervisor did not reach shutdown boundary")
	}
	requestDone := make(chan error, 1)
	go func() {
		if e := workerControlWrite(f.D, s, "resume"); e != nil {
			requestDone <- e
			return
		}
		next, e := workerLock(path, false)
		if next != nil {
			next.Close()
		}
		requestDone <- e
	}()
	select {
	case e := <-requestDone:
		t.Fatalf("control was published before old supervisor released ownership; replacement cannot own it: %v", e)
	case <-time.After(100 * time.Millisecond):
	}
	close(release)
	released = true
	select {
	case e := <-requestDone:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("resume was lost at supervisor exit")
	}
}

func TestWorkerStopWakesSupervisorWhenInitialPhaseBecameStale(t *testing.T) {
	var marker string
	f, o := newWorkerDeliveryFixture(t, func(f *deliveryFixture) {
		path := filepath.Join(f.R.WorkspaceRoot, "native-start-observer")
		script := []byte("#!/bin/sh\nprintf launched >\"$0.started\"\n")
		if err := os.WriteFile(path, script, 0700); err != nil {
			t.Fatal(err)
		}
		f.R.Runtime.PlyExecutable = Executable{Path: path, SHA256: hash(script)}
		f.D.Executable = func() (string, error) { return path, nil }
		marker = path + ".started"
	})
	// Ignore the fixture's initial start; only a new wake after the race counts.
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := workerLock(filepath.Join(workerDirectory(s), "supervisor.lock"), true)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	d := f.D
	d.Fault = func(point string) error {
		if point != "worker_control_after_read" {
			return nil
		}
		if e := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(current *workflowState) error { current.Result.Worker.Phase = "paused"; return nil }); e != nil {
			return e
		}
		if e := workerControlWrite(f.D, s, "pause"); e != nil {
			return e
		}
		parked, e := workerSupervisorExit(f.D, f.R.WorkspaceRoot, o.RunID, lock)
		if e != nil {
			return e
		}
		if !parked {
			return errors.New("fixture supervisor did not park")
		}
		return nil
	}
	if _, err = ControlWorker(d, f.R.WorkspaceRoot, o.RunID, "stop"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, err = os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("stop published after supervisor exit but stale initial phase prevented a replacement start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	var control workerControl
	if err = readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); err != nil || control.Action != "stop" {
		t.Fatalf("stop request lost: %+v %v", control, err)
	}
}

func TestWorkerNativeReviewRejectsDeveloperIdentityAndForgedReview(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = workerValidateReview(s, *o.Worker.Review, s.Result.Worker.CandidateOID, s.Result.Worker.CandidateTree); err != nil {
		t.Fatal(err)
	}
	s.Result.Worker.DeveloperSessions = append(s.Result.Worker.DeveloperSessions, "synthetic-reviewer-thread")
	if err = workerValidateReview(s, *o.Worker.Review, s.Result.Worker.CandidateOID, s.Result.Worker.CandidateTree); err == nil {
		t.Fatal("same developer/reviewer passed")
	}
	s.Result.Worker.DeveloperSessions = []string{"synthetic-developer-thread"}
	s.Result.Worker.ReviewProof = nil
	if err = workerValidateReview(s, *o.Worker.Review, s.Result.Worker.CandidateOID, s.Result.Worker.CandidateTree); err == nil {
		t.Fatal("self-declared review without provider result passed")
	}
	withoutOwner := f.D
	withoutOwner.WorkerOwner = nil
	if _, err = WorkflowDeliveryVerifyWithBinary(withoutOwner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary); err == nil {
		t.Fatal("external callback bypassed native owner capability")
	}
}

func TestWorkerRecoveryCannotConsumeAnotherCandidate(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary); err != nil {
		t.Fatal(err)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	s.Result.Delivery.Attempt.CandidateOID = strings.Repeat("e", 40)
	if err = workerRecoverNative(f.D, s); err == nil {
		t.Fatal("worker recovered a receipt belonging to the previous candidate")
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Result.Worker.Stage != "verify" || after.Result.Worker.Job == nil {
		t.Fatal("mismatched receipt advanced the pending worker")
	}
}

func workerTestReserveNative(t *testing.T, f deliveryFixture, id, kind string) {
	t.Helper()
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, id, func(s *workflowState) error {
		w := s.Result.Worker
		w.Stage = kind
		dir := filepath.Join(workerDirectory(*s), "attempts", kind+"-0001")
		if err := privateDir(dir); err != nil {
			return err
		}
		binary, _, err := deliveryReadBinding(w.CandidateBinary, 64<<20)
		if err != nil {
			return err
		}
		b, err := workflowKeep(f.D, filepath.Join(dir, "intent.json"), workerAttemptIntent{RunID: id, JobID: kind + "-0001", Kind: kind, Generation: w.Generation, CandidateOID: w.CandidateOID, CandidateTree: w.CandidateTree, Review: w.Review, CandidateBinary: &binary, Argv: []string{}, Environment: []string{}})
		if err != nil {
			return err
		}
		w.Job = &WorkerJob{ID: kind + "-0001", Kind: kind, Directory: dir, IntentHash: b.SHA256}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkerRecoveryRejectsNewBinaryAtSameCommit(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nexit 1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary); err == nil || !strings.Contains(err.Error(), "verification did not qualify") {
		t.Fatal(err)
	}
	replacement := filepath.Join(f.R.WorkspaceRoot, "new-private-candidate")
	if err := os.WriteFile(replacement, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Worker.CandidateBinary = replacement; return nil }); err != nil {
		t.Fatal(err)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = workerRecoverNative(f.D, state); err == nil || !strings.Contains(err.Error(), "private executable") {
		t.Fatalf("old receipt consumed for a new private executable: %v", err)
	}
	after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Result.Worker.Iteration != state.Result.Worker.Iteration || after.Result.Worker.Stage != "verify" {
		t.Fatal("old failure started a correction for the new candidate")
	}
}
