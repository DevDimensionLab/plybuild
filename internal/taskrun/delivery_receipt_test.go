package taskrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Setting the additive protocol through JSON lets this regression reproduce the
// old behavior: the old reader ignores it and demands blanket permission.
func receiptFixture(t *testing.T, provider string, options ...deliveryFixtureOptions) deliveryFixture {
	t.Helper()
	f := newDeliveryFixture(t, provider, options...)
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	raw, err := json.Marshal(f.R)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err = json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	fields["delivery"].(map[string]any)["startup_protocol"] = MandateReceiptV1
	draft := fields["handoff_draft"].(map[string]any)
	draft["delivery_binding"].(map[string]any)["startup_protocol"] = MandateReceiptV1
	raw, err = json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &f.R); err != nil {
		t.Fatal(err)
	}
	f.File = writeAny(t, f.R.WorkspaceRoot, "receipt-request.json", f.R)
	return f
}

func TestDeliveryReceiptDoesNotRequireBlanketPermission(t *testing.T) {
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := receiptFixture(t, provider)
			o := workflowTestStart(t, f.workflowFixture)
			in := DeliveryBeginInput{Provider: provider, NativeSession: o.Transport.AgentSessionID}
			first, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in)
			if err != nil {
				t.Fatalf("neutral mandate receipt must reach authorized work without blanket permission: %v", err)
			}
			if first.Delivery.PermissionState != RuntimePerAction {
				t.Fatalf("receipt fabricated permission: %+v", first.Delivery)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			raw, err := os.ReadFile(s.Acceptance.Locator)
			if err != nil {
				t.Fatal(err)
			}
			var a DeliveryAcceptance
			if err = json.Unmarshal(raw, &a); err != nil {
				t.Fatal(err)
			}
			if a.DeliveryPermission.PermissionConfirmed {
				t.Fatal("receipt asserted blanket permission")
			}
			var sandbox map[string]any
			if err = json.Unmarshal(a.Sandbox, &sandbox); err != nil {
				t.Fatal(err)
			}
			if sandbox["matches_contract"] != nil {
				t.Fatalf("unobserved sandbox coverage was attested: %s", a.Sandbox)
			}
			before := *s.Acceptance
			if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
				t.Fatal(err)
			}
			s, _ = workflowRead(f.R.WorkspaceRoot, o.RunID)
			if *s.Acceptance != before {
				t.Fatal("lost receipt response changed immutable evidence")
			}
		})
	}
}

var errReceiptSubprocessCloseout = errors.New("receipt closeout cannot inherit observers across a subprocess")

// Execute the same real closeout body in the journey and regression. The
// fixture's cleanupCLIAt guard observes any subprocess restoration inside it.
func receiptJourneyClosure(t *testing.T, f deliveryFixture, o WorkflowRun, observers *receiptObservers, failure string) {
	assertReceiptSimulatedObserverClosure(t, f, o, observers, failure)
}

func TestDeliveryReceiptClosureDispatch(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, observers := receiptCLIQualified(t, binary)
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	calls := map[string]int{}
	for _, command := range []string{"pane get", "pane list", "agent get", "agent list", "agent prompt", "agent start", "tab create", "pane close"} {
		calls[command] = workflowTestCalls(t, f.workflowFixture, command)
	}
	receiptJourneyClosure(t, f, o, observers, "herdr_denied")
	if observers.subprocess != 0 {
		// The actual callable ran up to cleanupCLIAt's pre-exec tripwire.
		// Compare with phase entry: startup/CLI verify already happened, but
		// no closeout launch, inspection or integration may have happened.
		after, readErr := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
		shown := mustShow(t, f, o.RunID)
		if readErr != nil || !equal(registry, after) || !equal(shown.Delivery.Events, o.Delivery.Events) || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != before || observers.local != 0 {
			t.Fatal("blocked closeout callable changed native state or inspected local processes")
		}
		c := o.Delivery.Candidates[0]
		if _, e := os.Stat(c.TaskResult.SourceLocator); e != nil || gitOutput(t, f.Parent, "rev-parse", c.TaskResult.SourceRef) != c.OID {
			t.Fatal("blocked closeout callable changed source directory/ref")
		}
		for command, count := range calls {
			if workflowTestCalls(t, f.workflowFixture, command) != count {
				t.Fatalf("blocked closeout callable reached downstream %s", command)
			}
		}
		observers.assertUnused(t)
		t.Fatalf("actual receipt closeout callable attempted subprocess before exec.Command: attempts=%d; closeout launch/inspection/integration effects=0; native state unchanged", observers.subprocess)
	}
	assertCleanupRemoved(t, f, mustShow(t, f, o.RunID), gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"))
}

// All adapters are installed before startup. Ordinary non-adopted startup does
// not consume SourceWriters or StartupProcessTree; their zero counts are not
// behavioral coverage of adoption or recovery. LocalProcesses is exercised only
// after modeled owner exit, through production WorkflowFollow/closeout.
type receiptObservers struct {
	startup, source, local int
	subprocess             int
	session, transcript    string
	table                  []LocalProcess
	err                    error
}

func receiptControlledRuntime(t *testing.T, f *deliveryFixture) *receiptObservers {
	t.Helper()
	home := filepath.Join(f.R.WorkspaceRoot, "receipt-home")
	transcript := filepath.Join(home, ".codex", "sessions", "2026", "10", "24", "rollout-fixture-session.jsonl")
	if err := os.MkdirAll(filepath.Dir(transcript), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(transcript, []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", filepath.Join(home, ".codex"))
	s := &receiptObservers{session: "fixture-session", transcript: transcript, table: []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}}
	for _, cwd := range []string{f.Parent, f.Prepared.Preparation.Preparation.Plan.WorktreePath} {
		guardCleanupCLIAt(t, cwd, func() error {
			s.subprocess++
			return errReceiptSubprocessCloseout
		})
	}
	f.D.LocalProcesses = func(session, transcript string) ([]LocalProcess, error) {
		s.local++
		if session != s.session || transcript != s.transcript {
			t.Fatalf("wrong local inspection identity: session=%q transcript=%q", session, transcript)
		}
		return s.table, s.err
	}
	f.D.StartupProcessTree = func(pid int) ([]StartupProcess, error) {
		s.startup++
		return []StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	f.D.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		s.source++
		return workspace.SourceWriterObservation{State: "absent", Method: "receipt_fixture", Limitation: "Controlled disposable fixture; no OS observation"}, nil
	}
	// SystemDependencies callbacks captured older values. Rebind explicitly.
	f.D.Workflow.TaskWorkspace = &f.D.Workspace
	ownerPaneStandin(t, f)
	return s
}

func (s *receiptObservers) assertUnused(t *testing.T) {
	t.Helper()
	if s.startup != 0 || s.source != 0 {
		t.Fatalf("ordinary journey consumed recovery/adoption observers: startup=%d source=%d", s.startup, s.source)
	}
	t.Logf("defensive observers unexercised: StartupProcessTree=%d SourceWriters=%d", s.startup, s.source)
}

func receiptCLICall(t *testing.T, binary string, f deliveryFixture, args ...string) (WorkflowRun, error) {
	t.Helper()
	if len(args) < 2 || args[0] != "execute" || args[1] != "check-in" && args[1] != "begin" && args[1] != "verify" {
		t.Fatalf("unmodeled receipt subprocess boundary: %q", args)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir, c.WaitDelay = f.Prepared.Preparation.Preparation.Plan.WorktreePath, time.Second
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run() // Exact owned child, synchronously waited; never a detached supervisor.
	t.Logf("receipt CLI argv=%q cwd=%q stdout=%s stderr=%s exit=%v", c.Args, c.Dir, stdout.Bytes(), stderr.Bytes(), err)
	var o WorkflowRun
	if decodeErr := json.Unmarshal(stdout.Bytes(), &o); decodeErr != nil && err == nil {
		err = decodeErr
	}
	return o, err
}

func receiptCLIQualified(t *testing.T, binary string) (deliveryFixture, WorkflowRun, *receiptObservers) {
	t.Helper()
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic receipt journey", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	observers := receiptControlledRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	if o.Transport.AgentSessionID != observers.session || f.R.Runtime.Mode != "interactive" {
		t.Fatal("receipt fixture changed its controlled session or interactive verification boundary")
	}
	ok := func(args ...string) WorkflowRun {
		t.Helper()
		o, err := receiptCLICall(t, binary, f, args...)
		if err != nil {
			t.Fatal(err)
		}
		return o
	}
	args := []string{"execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID}
	for _, bad := range []struct {
		index int
		value string
	}{{4, o.Paths.Context + "wrong"}, {6, "claude"}, {8, "wrong-session"}} {
		invalid := append([]string{}, args...)
		invalid[bad.index] = bad.value
		if _, err := receiptCLICall(t, binary, f, invalid...); err == nil {
			t.Fatal("invalid receipt accepted")
		}
	}
	o = ok(args...)
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Acceptance == nil {
		t.Fatalf("receipt binding missing: %v", err)
	}
	firstBinding := *state.Acceptance
	first, err := os.ReadFile(firstBinding.Locator)
	if err != nil {
		t.Fatal(err)
	}
	o = ok(args...)
	state, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Acceptance == nil || *state.Acceptance != firstBinding {
		t.Fatalf("repeated CLI receipt changed binding: %v", err)
	}
	repeated, err := os.ReadFile(state.Acceptance.Locator)
	if err != nil || string(first) != string(repeated) {
		t.Fatal("repeated CLI receipt changed immutable identity")
	}
	if o.DeliveryStatus.Startup.SandboxCoverage != "unknown" || o.DeliveryStatus.Startup.RuntimeDecision != "per_action_not_observed" {
		t.Fatalf("receipt fabricated runtime facts: %+v", o.DeliveryStatus.Startup)
	}
	if _, err := receiptCLICall(t, binary, f, "execute", "begin", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID, "--policy", "unknown", "--read-root", f.R.WorkspaceRoot, "--write-root", f.R.WorkspaceRoot, "--permission-confirmed"); err == nil {
		t.Fatal("legacy entrypoint accepted new protocol")
	}
	deliveryTestCandidate(t, f, "receipt candidate")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'receipt candidate'\n\"$PLY_CANDIDATE_BINARY\" workflow execute check-in --help >/dev/null\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "receipt-journey")
	o = ok(automaticCLIVerifyArgs(o, review, binary)...)
	if o.Delivery.Phase != "automatic_acceptance_passed" {
		t.Fatalf("receipt did not reach candidate qualification: %+v", o.Delivery)
	}
	observers.assertUnused(t)
	if observers.local != 0 {
		t.Fatal("startup/receipt/verify unexpectedly inspected local writers")
	}
	return f, o, observers
}

func TestDeliveryReceiptCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, observers := receiptCLIQualified(t, binary)
	receiptJourneyClosure(t, f, o, observers, "herdr_denied")
	if observers.subprocess != 0 {
		t.Fatalf("actual receipt closeout reached subprocess boundary: %d", observers.subprocess)
	}
}

// Public Delivery register/check/execute and subprocess follow remain covered by
// their own tests. This receipt-local closure uses production functions with
// simulated inventories and real native Git, integration, queue/base records.
func assertReceiptSimulatedObserverClosure(t *testing.T, f deliveryFixture, o WorkflowRun, observers *receiptObservers, failure string) {
	t.Helper()
	var err error
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	cleanupOwnerHandoff(t, f, o, err)
	c := o.Delivery.Candidates[0]
	if observers.local != 0 {
		t.Fatal("owner-present handoff reached local inspection")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.IntegrationResults) != 1 || len(registry.HumanQARecords) != 0 {
		t.Fatalf("wrong native integration or fabricated QA: %v", err)
	}
	q, err := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if err != nil || q.Current != nil || len(q.Active) != 0 || len(q.Pending) != 0 {
		t.Fatalf("integration did not close the queue: %+v %v", q, err)
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatalf("Epic base did not advance: %+v %v", base, err)
	}
	assertStable := func(current WorkflowRun, pending bool) {
		t.Helper()
		r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
		b, baseErr := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
		queue, queueErr := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
		if e != nil || baseErr != nil || queueErr != nil || !equal(r.IntegrationResults, registry.IntegrationResults) || !equal(r.TaskQueueEvents, registry.TaskQueueEvents) || !equal(b, base) || queue.Current != nil || len(queue.Active) != 0 || len(queue.Pending) != 0 || len(r.HumanQARecords) != 0 || current.DeliveryStatus.HumanJudgment != nil {
			t.Fatal("closeout changed integration/closed queue/advanced base or fabricated human QA")
		}
		cc := current.Delivery.Candidates[0]
		if cc.TaskResult.ID != c.TaskResult.ID || !equal(cc.Integration, c.Integration) || current.DeliveryStatus.FinalDelivery.State != "delivered" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
			t.Fatal("closeout changed exact delivery or repeated integration")
		}
		if pending {
			if current.Delivery.Phase != "closing" || !equal(current.Delivery.Events, o.Delivery.Events) || gitOutput(t, f.Parent, "rev-parse", c.TaskResult.SourceRef) != c.OID {
				t.Fatal("negative observation changed closing phase, events or source ref")
			}
			if _, e := os.Stat(c.TaskResult.SourceLocator); e != nil {
				t.Fatal("negative observation removed source")
			}
			closeout, e := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
			if e != nil || closeout != nil {
				t.Fatalf("negative observation manufactured a cleanup/absence receipt: %+v %v", closeout, e)
			}
		}
	}
	assertStable(o, true)
	cleanupOwnerLeaves(t, f)
	coordinator := f.D
	coordinator.CWD = func() (string, error) { return f.Parent, nil }
	reviewRoundIn(t, "w-fixture:root")
	follow := func() (WorkflowRun, error) {
		return WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 5*time.Second)
	}
	if failure == "herdr_denied" {
		assertCleanupUnknownInspectionPreservesDeliveredIntegration(t, f, o, reflog, func() (WorkflowRun, error) {
			blocked, err := follow()
			if err == nil || !strings.Contains(err.Error(), "owner_unknown") || observers.local != 0 {
				t.Fatalf("denied managed inventory reached local observer: %v calls=%d", err, observers.local)
			}
			assertStable(blocked, true)
			return blocked, err
		})
	} else {
		want := "owner_unknown"
		switch failure {
		case "local_denied":
			observers.table, observers.err = nil, os.ErrPermission
		case "empty":
			observers.table = nil
		case "incomplete":
			observers.table = []LocalProcess{{PID: os.Getpid() + 10000, CWD: f.Parent}}
		case "external_task", "external_session":
			want = "owner_active"
			writer := LocalProcess{PID: os.Getpid() + 10000, CWD: c.TaskResult.SourceLocator}
			if failure == "external_session" {
				writer.CWD, writer.Session = f.Parent, "transcript_open"
			}
			observers.table = append(observers.table, writer)
		default:
			t.Fatalf("unknown modeled failure %q", failure)
		}
		before := observers.local
		blocked, err := follow()
		if err == nil || !strings.Contains(err.Error(), want) || observers.local-before != 1 {
			t.Fatalf("modeled %s was retried or did not preserve uncertainty/writer: %v calls=%d", failure, err, observers.local-before)
		}
		assertStable(blocked, true)
	}
	// Explicitly replace only the modeled failure, then resume the same closeout.
	observers.table, observers.err = []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}, nil
	before := observers.local
	done, err := follow()
	if err != nil {
		t.Fatal(err)
	}
	if observers.local-before != 2 {
		t.Fatalf("WorkflowFollow must inspect locally before closeout and under its ownership guard: calls=%d", observers.local-before)
	}
	t.Logf("%s: LocalProcesses successful follow calls=%d; session=%q transcript=%q observing_pid=%d", failure, observers.local-before, observers.session, observers.transcript, os.Getpid())
	assertCleanupRemoved(t, f, done, reflog)
	assertStable(done, false)
	again, err := follow()
	if err != nil || !equal(again.Delivery.Events, done.Delivery.Events) {
		t.Fatalf("repeat completion changed delivery events: %v", err)
	}
	assertCleanupRemoved(t, f, again, reflog)
	assertStable(again, false)
	shown, err := WorkflowShow(coordinator, f.R.WorkspaceRoot, o.RunID)
	if err != nil || !equal(shown.Delivery.Events, done.Delivery.Events) {
		t.Fatalf("readback changed completed history: %v", err)
	}
	assertStable(shown, false)
	if observers.local-before != 2 {
		t.Fatal("completed follow/readback repeated local inspection")
	}
	observers.assertUnused(t)
}

func TestDeliveryReceiptSimulatedObserverCloseout(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, failure := range []string{"local_denied", "empty", "incomplete", "external_task", "external_session"} {
		t.Run(failure, func(t *testing.T) {
			f, o, observers := receiptCLIQualified(t, binary)
			receiptJourneyClosure(t, f, o, observers, failure)
			if observers.subprocess != 0 {
				t.Fatalf("actual receipt closeout reached subprocess boundary: %d", observers.subprocess)
			}
		})
	}
}

func TestDeliveryReceiptReportedRootsNeverAttestSandbox(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID, ReadRoots: []string{f.R.WorkspaceRoot}, WriteRoots: []string{f.Prepared.Preparation.Preparation.Plan.WorktreePath}, Profile: "auto", Policy: "Source: synthetic recipient observation; future decisions unknown"}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	b, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	var a DeliveryAcceptance
	if err = json.Unmarshal(b, &a); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(a.Sandbox), `"matches_contract":null`) || a.DeliveryPermission.PermissionConfirmed {
		t.Fatalf("reported roots became authority: %s", b)
	}
}

func TestDeliveryReceiptLaterRuntimeRefusalPreservesStop(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
		t.Fatal(err)
	}
	proof := writeAny(t, f.R.WorkspaceRoot, "actual-denial.json", map[string]string{"source": "synthetic provider classifier", "reason": "action denied", "effect": "none observed"})
	var stop DeliveryStopInput
	raw, _ := json.Marshal(map[string]any{"event_id": "later-refusal", "provider": "claude", "native_session": o.Transport.AgentSessionID, "stage": "work", "observer": "fixture coordinator", "reason": "Provider refused the relevant test action", "evidence_path": proof, "action": "run relevant test", "effect_status": "none_observed", "rule": "synthetic-classifier-rule"})
	if err := json.Unmarshal(raw, &stop); err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop)
	if err != nil {
		t.Fatalf("later action refusal must preserve phase, source and effects: %v", err)
	}
	if _, err = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, stop); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
		t.Fatal("receipt cleared a known refusal")
	}
	b, err := workflowBound(*shown.Delivery.ActiveStop, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"action":"run relevant test"`) || !strings.Contains(string(b), `"effect_status":"none_observed"`) || !strings.Contains(string(b), `"rule":"synthetic-classifier-rule"`) {
		t.Fatalf("refusal lost action or effects: %s", b)
	}
	if shown.DeliveryStatus.Startup.RuntimeDecision != "preserved_stop" {
		t.Fatal("readback hid runtime refusal")
	}
}

func TestDeliveryReceiptDoesNotDowngradeLegacyContract(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	in := DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err == nil {
		t.Fatal("check-in downgraded frozen legacy startup")
	}
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
		t.Fatal("legacy permission requirement changed")
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if s.Acceptance != nil || s.StartSHA256 != nil || s.Result.Delivery.PermissionState != "pending_runtime_acceptance" {
		t.Fatal("legacy authority/history changed")
	}
}

func TestDeliveryReceiptGuideAndVisibleMandate(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	prompt := deliveryInstructions(s)
	if !strings.Contains(prompt, "Your authorized assignment is") || !strings.Contains(prompt, "workflow execute check-in") || strings.Contains(prompt, "permission-confirmed") {
		t.Fatalf("visible task omitted the neutral mandate: %s", prompt)
	}
	b, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "delivery", "callback-guide.md"))
	if err != nil {
		t.Fatal(err)
	}
	guide := string(b)
	for _, expected := range []string{"Progress is optional; verifier_results may be empty", "never a prerequisite\nfor starting separate review", "IDs are local stable declaration keys", "^TestDeliveryReportVerifierFieldDiagnostics$", "exit=null and a truthful reason"} {
		if !strings.Contains(guide, expected) {
			t.Fatalf("callback guide lacks optional progress guidance %q", expected)
		}
	}
	if !strings.Contains(guide, "mandate_receipt_v1") || !strings.Contains(guide, "matches_contract=null") || strings.Contains(guide, "## Advanced acceptance file contract") || strings.Contains(guide, "--permission-confirmed") {
		t.Fatalf("guide routes new contract to legacy self-attestation: %s", guide)
	}
	if _, err = os.Stat(filepath.Join(o.Paths.RunRoot, "delivery", "templates", "acceptance.json")); !os.IsNotExist(err) {
		t.Fatal("new flow still supplies handwritten legacy acceptance template")
	}
}

func TestDeliveryReceiptAuthorityRemainsBoundToMandateAndEffects(t *testing.T) {
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic authority test", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	o := workflowTestStart(t, f.workflowFixture)
	if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	auth, err := deliveryAuthorityState(f.D, s)
	if err != nil {
		t.Fatal(err)
	}
	if auth.AuthorityKind != workspace.DeliveryMandateReceiptV1 || auth.PermissionConfirmed {
		t.Fatalf("native human mandate became permission self-attestation: %+v", auth)
	}
	var accepted DeliveryAcceptance
	b, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &accepted); err != nil {
		t.Fatal(err)
	}
	for _, change := range []string{"effects", "target", "protocol", "permission"} {
		t.Run(change, func(t *testing.T) {
			raw, _ := json.Marshal(accepted)
			var bad DeliveryAcceptance
			_ = json.Unmarshal(raw, &bad)
			switch change {
			case "effects":
				bad.DeliveryPermission.AllowedEffects = append(bad.DeliveryPermission.AllowedEffects, "push")
			case "target":
				bad.DeliveryPermission.DeliveryAgreementSHA256 = hash([]byte("foreign target"))
			case "protocol":
				bad.DeliveryPermission.StartupProtocol = ""
			case "permission":
				bad.DeliveryPermission.PermissionConfirmed = true
			}
			if err := positiveDeliveryClaim(bad.Acceptance, bad.DeliveryPermission, s.Request); err == nil {
				t.Fatal("receipt accepted expanded or changed authority")
			}
		})
	}
	altered := s
	altered.Result.Handoff.SHA256 = hash([]byte("different mandate"))
	if _, err = deliveryAuthorityState(f.D, altered); err == nil {
		t.Fatal("receipt accepted a changed mandate")
	}
	priorCWD := f.D.CWD
	f.D.CWD = func() (string, error) { return f.Parent, nil }
	if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}); err == nil {
		t.Fatal("receipt accepted another cwd")
	}
	f.D.CWD = priorCWD
}

// A runtime refusal may follow ordinary work outside native Delivery.Attempt.
// Permission to continue cannot establish the outcome of that earlier action.
func TestDeliveryReceiptResolutionPreservesUncertainExternalEffects(t *testing.T) {
	for _, effect := range []string{"unknown", "partial"} {
		t.Run(effect, func(t *testing.T) {
			f := receiptFixture(t, "claude")
			o := workflowTestStart(t, f.workflowFixture)
			in := DeliveryBeginInput{Provider: "claude", NativeSession: o.Transport.AgentSessionID}
			if _, err := WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err != nil {
				t.Fatal(err)
			}
			proof := writeAny(t, f.R.WorkspaceRoot, "external-work-refusal.json", map[string]string{"source": "synthetic provider refusal", "effect_status": effect})
			stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "external-refusal", Provider: "claude", NativeSession: o.Transport.AgentSessionID, Stage: "work", Observer: "synthetic coordinator", Reason: "Tool response lost after ordinary work; earlier effects unresolved", EvidencePath: proof, Action: "ordinary work outside the native delivery attempt", EffectStatus: effect})
			if err != nil {
				t.Fatal(err)
			}
			stop := *stopped.Delivery.ActiveStop
			if stopped.Delivery.Attempt != nil {
				t.Fatal("fixture unexpectedly has a native delivery attempt")
			}
			for attempt := 0; attempt < 2; attempt++ {
				resumed, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stop.SHA256, "Current refusal cause repaired; earlier external effects remain unresolved", "Continue the clarified changed fixture action")
				if err == nil {
					t.Fatalf("permission and text cleared %s external effects: phase=%s stop=%v", effect, resumed.Delivery.Phase, resumed.Delivery.ActiveStop)
				}
				observed, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				if observed.Delivery.ActiveStop == nil || *observed.Delivery.ActiveStop != stop || observed.Delivery.Phase != "stopped" || len(observed.Delivery.Events) != 1 {
					t.Fatal("failed resolution changed the stop or its history")
				}
				if _, err = WorkflowDeliveryBegin(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in); err == nil {
					t.Fatal("receipt reopened unresolved external work")
				}
			}
		})
	}
}

func TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint(t *testing.T) {
	for _, receipt := range []bool{false, true} {
		name := "legacy"
		if receipt {
			name = "mandate_receipt"
		}
		t.Run(name, func(t *testing.T) {
			var f deliveryFixture
			if receipt {
				f = receiptFixture(t, "codex")
			} else {
				f = newDeliveryFixture(t, "codex")
			}
			o := workflowTestStart(t, f.workflowFixture)
			proof := writeAny(t, f.R.WorkspaceRoot, "pre-receipt-denial.json", map[string]string{"source": "synthetic provider refusal", "effect_status": "not_started"})
			stopped, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "before-receipt", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "acceptance", Observer: "fixture coordinator", Reason: "Initial receipt action refused before starting", EvidencePath: proof, Action: "record initial receipt", EffectStatus: DeliveryEffectNotStarted})
			if err != nil {
				t.Fatal(err)
			}
			want := "execute begin"
			if receipt {
				want = "execute check-in"
			}
			if !strings.Contains(stopped.NextAction.Message, strings.TrimPrefix(want, "execute ")+" is still required") {
				t.Fatalf("stopped %s contract recommends wrong startup entrypoint: %s", name, stopped.NextAction.Message)
			}
			resumed, err := WorkflowResolveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, stopped.Delivery.ActiveStop.SHA256, "Actual refusal cause repaired; initial receipt never started", "Continue the clarified changed fixture action")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(resumed.NextAction.Message, want) {
				t.Fatalf("resolved %s contract recommends wrong startup entrypoint: %s", name, resumed.NextAction.Message)
			}
			in := beginFixtureInput(f, o)
			if receipt {
				in = DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}
			}
			if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, !receipt); err == nil {
				t.Fatal("recovery allowed the other protocol's startup entrypoint")
			}
			received, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, receipt)
			if err != nil {
				t.Fatalf("recommended frozen entrypoint did not reach receipt: %v", err)
			}
			if received.Delivery.Phase != "working" || received.Delivery.ActiveStop != nil {
				t.Fatal("correct startup failed to resume the same mandate")
			}
		})
	}
}
