package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// T0F15276E: every runtime observation is injected before starting a fixture.
// Native Git, journal, QA, integration, release and closeout records are real.
func nf01Runtime(t *testing.T, f *deliveryFixture) {
	t.Helper()
	parallelRuntime(t, f)
	f.D.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		return workspace.SourceWriterObservation{State: "absent", Method: "modeled_nf01", Limitation: "Disposable fixture; no OS observation"}, nil
	}
	f.D.RunVerificationCommand = func(c *exec.Cmd) error {
		if c.Path == "/usr/bin/sandbox-exec" {
			if len(c.Args) != 5 || c.Args[1] != "-f" || c.Args[3] != "/bin/sh" {
				return fmt.Errorf("unexpected modeled sandbox argv: %v", c.Args)
			}
			// Simulate OS sandbox acceptance, not real sandbox enforcement.
			// Still execute and wait for the exact native acceptance snapshot.
			c.Path = "/bin/sh"
			c.Args = append([]string(nil), c.Args[3:]...)
		}
		return c.Run()
	}
	f.D.Workflow.TaskWorkspace = &f.D.Workspace
}
func nf01State(t *testing.T, f workflowFixture, o WorkflowRun) workflowState {
	t.Helper()
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func nf01Terminal(t *testing.T, f workflowFixture, o WorkflowRun) {
	t.Helper()
	s := nf01State(t, f, o)
	fact, e := parallelConsumerClass(f.D, s, "")
	if e != nil || fact.Class != "PAR-A02" || !fact.Compatible || len(fact.RemainingEffects) != 0 {
		t.Fatalf("proven terminal must be A02 without live reader: %+v %v", fact, e)
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	d := f.D
	forbidden := 0
	d.Workspace.ProjectLocks = activationNoProjectLock{&forbidden}
	d.Workspace.WorkItems = activationNoStoreWrites{d.Workspace.WorkItems, &forbidden}
	d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		forbidden++
		return nil, fmt.Errorf("forbidden process observation")
	}
	d.StartupProcessTree = func(int) ([]StartupProcess, error) { forbidden++; return nil, fmt.Errorf("forbidden process tree") }
	d.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		forbidden++
		return workspace.SourceWriterObservation{}, fmt.Errorf("forbidden source writers")
	}
	before := nf01NativeState(t, f.R.WorkspaceRoot)
	first, sha, e := parallelActivationConsumers(d, f.R.WorkspaceRoot, r, "")
	if e != nil {
		t.Fatal(e)
	}
	second, sha2, e := parallelActivationConsumers(d, f.R.WorkspaceRoot, r, "")
	if e != nil || sha != sha2 || !equal(first, second) {
		t.Fatal("unstable collection", e)
	}
	found := false
	for _, c := range first {
		if c.ID == fact.ID {
			found = true
			if c.Class != "PAR-A02" || !c.Compatible {
				t.Fatalf("terminal collection: %+v", c)
			}
		}
	}
	if forbidden != 0 || !found || !equal(before, nf01NativeState(t, f.R.WorkspaceRoot)) {
		t.Fatal("missing terminal or collection wrote native state")
	}
}
func nf01StartedReport(t *testing.T, postStart bool, seeds ...string) (workflowFixture, WorkflowRun) {
	t.Helper()
	var f workflowFixture
	if len(seeds) == 0 {
		f = workflowTestFixture(t)
	} else {
		root, e := filepath.EvalSymlinks(t.TempDir())
		if e != nil {
			t.Fatal(e)
		}
		f = newWorkflowFixtureSeed(t, root, "", seeds[0], seeds[1])
		bin := filepath.Join(root, "bin")
		if e = os.Mkdir(bin, 0700); e != nil {
			t.Fatal(e)
		}
		if e = os.Symlink(f.R.Runtime.Executable.Path, filepath.Join(bin, "codex")); e != nil {
			t.Fatal(e)
		}
		t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
		t.Setenv("HERDR_ENV", "1")
	}
	b := mustReadAL01(t, f.R.Herdr.Executable.Path)
	b = bytes.Replace(b, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if e := os.WriteFile(f.R.Herdr.Executable.Path, b, 0700); e != nil {
		t.Fatal(e)
	}
	f.R.Herdr.Executable.SHA256 = hash(b)
	f.File = writeAny(t, f.R.WorkspaceRoot, "nf01-report.json", f.R)
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}}, nil
	}
	f.D.StartupProcessTree = func(pid int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	f.D.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		return workspace.SourceWriterObservation{State: "absent", Method: "modeled_nf01"}, nil
	}
	f.D.Workflow.TaskWorkspace = &f.D.Workspace
	o := workflowTestStart(t, f)
	if postStart {
		f.Acceptance.RunID, f.Acceptance.RequestSHA256, f.Acceptance.SessionID = o.RunID, o.RequestSHA256, o.SessionID
		o = workflowTestAccept(t, f, o)
		f.Report.RunID, f.Report.RequestSHA256, f.Report.SessionID = o.RunID, o.RequestSHA256, o.SessionID
		o = workflowTestReport(t, f, o)
	}
	return f, o
}
func nf01Report(t *testing.T, postStart bool, seeds ...string) (workflowFixture, WorkflowRun) {
	t.Helper()
	f, o := nf01StartedReport(t, postStart, seeds...)
	workflowTestModel(t, f, map[string]any{"closed_transport": true})
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	if postStart {
		proof := writeAny(t, f.R.WorkspaceRoot, "modeled-stop.json", map[string]any{"fixture": "no effect started"})
		var e error
		o, e = WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryStopInput{EventID: "nf01-stop", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "modeled coordinator", Reason: "Modeled stopped report", Action: "report transport", EffectStatus: DeliveryEffectNotStarted, EvidencePath: proof})
		if e != nil {
			t.Fatal(e)
		}
	}
	var e error
	o, e = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "modeled coordinator", Reason: "Retain history in disposable fixture"})
	if e != nil {
		t.Fatal(e)
	}
	return f, o
}
func TestNF01RetiredReportCompatibility(t *testing.T) {
	for _, post := range []bool{false, true} {
		t.Run(fmt.Sprint(post), func(t *testing.T) { f, o := nf01Report(t, post); nf01Terminal(t, f, o) })
	}
}
func nf01Headless(t *testing.T, complete bool) (deliveryFixture, WorkflowRun) {
	t.Helper()
	f, o := newWorkerDeliveryFixtureOptions(t, deliveryFixtureOptions{StartFixture: nf01ReserveOnly}, func(f *deliveryFixture) { nf01Runtime(t, f) })
	if !complete {
		return f, o
	}
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ngrep -q 'worker native candidate' README.md\n"), 0600); e != nil {
		t.Fatal(e)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	var e error
	o, e = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary)
	if e != nil {
		st := nf01State(t, f.workflowFixture, o)
		if a := st.Result.Delivery.Attempt; a != nil {
			filepath.WalkDir(a.Path, func(p string, de os.DirEntry, err error) error {
				if err == nil && !de.IsDir() && (strings.Contains(p, "stdout") || strings.Contains(p, "stderr") || strings.Contains(p, "completed-result")) {
					b, _ := os.ReadFile(p)
					t.Logf("%s: %s", p, b)
				}
				return nil
			})
		}
		t.Fatal(e)
	}
	workerTestReserveNative(t, f, o.RunID, "integrate")
	o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if e != nil {
		t.Fatal(e)
	}
	return f, o
}
func TestNF01HeadlessTerminalCompatibility(t *testing.T) {
	f, o := nf01Headless(t, true)
	nf01Terminal(t, f.workflowFixture, o)
}
func nf01Human(t *testing.T, human bool, seeds ...string) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}
	if human {
		a.SchemaVersion = 2
		a.IntegrationOwner = workspace.IntegrationOwnerHuman
	}
	option := deliveryFixtureOptions{Agreement: &a, ParentRef: "main"}
	if len(seeds) > 0 {
		option.SeedSource, option.SeedOID = seeds[0], seeds[1]
		option.ParentRef = "epic"
	}
	f := newDeliveryFixture(t, "codex", option)
	nf01Runtime(t, &f)
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "nf01-human.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	if len(seeds) == 0 {
		deliveryTestCandidate(t, f, "human terminal candidate")
	} else {
		runGit(t, f.Prepared.Preparation.Plan.WorktreePath, "merge", "--ff-only", seeds[1])
		if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest -f go.mod\n"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	var e error
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "nf01-human"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Modeled owner release"); e != nil {
		t.Fatal(e)
	}
	in := humanIntegrationTestInput(t, f, o, "pass", "nf01-human")
	if _, e = AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); e != nil {
		t.Fatal(e)
	}
	c := o.Delivery.Candidates[0]
	if _, e = CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in.PlanSHA256); e != nil {
		t.Fatal(e)
	}
	s := nf01State(t, f.workflowFixture, o)
	own := *s.Result.Delivery.OwnershipRelease
	closeIn := workspace.TaskCloseoutInput{OperationID: in.PlanID, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ExpectedResultOID: c.OID, ExpectedResultTree: c.Tree, TargetRef: f.R.Delivery.Agreement.TargetRef, Keep: true, Ownership: workspace.TaskCloseoutOwnership{Released: true, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, EvidenceLocator: own.Locator, EvidenceSHA256: own.SHA256}}
	p, e := workspace.PreviewTaskCloseout(f.D.Workspace, closeIn)
	if e != nil || !p.Ready {
		t.Fatalf("native closeout preview: %+v %v", p, e)
	}
	if _, e = workspace.ApplyTaskCloseout(f.D.Workspace, closeIn, p.PlanSHA256); e != nil {
		t.Fatal(e)
	}
	return f, o
}
func TestNF01HumanTerminalCompatibility(t *testing.T) {
	for _, human := range []bool{false, true} {
		t.Run(fmt.Sprint(human), func(t *testing.T) { f, o := nf01Human(t, human); nf01Terminal(t, f.workflowFixture, o) })
	}
}
func TestNF01PreparationBindingDiagnostics(t *testing.T) {
	f, o := nf01Headless(t, false)
	week := filepath.Join(f.R.WorkspaceRoot, ".ply", "week-runs", "modeled-week")
	if e := os.MkdirAll(week, 0700); e != nil {
		t.Fatal(e)
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	facts, _, e := parallelActivationConsumers(f.D, f.R.WorkspaceRoot, r, "")
	if e != nil {
		t.Fatal(e)
	}
	run, weekFound := false, false
	for _, c := range facts {
		if c.ID == "preparation/"+f.R.PreparationID {
			t.Fatalf("captured headless binding has false duplicate preparation blocker: %+v", c)
		}
		if c.ID == "run/"+o.RunID {
			run = c.Class == "unknown" && !c.Compatible && strings.Contains(strings.Join(c.Reasons, " "), f.R.PreparationID)
		}
		if c.ID == "week-runs/modeled-week" {
			weekFound = c.Class == "unknown"
		}
	}
	if !run || !weekFound {
		t.Fatalf("Run/preparation identity or week authority lost: %+v", facts)
	}
}

func nf01Clone(t *testing.T, s workflowState) workflowState {
	t.Helper()
	b, e := json.Marshal(s)
	if e != nil {
		t.Fatal(e)
	}
	var out workflowState
	if e = json.Unmarshal(b, &out); e != nil {
		t.Fatal(e)
	}
	return out
}
func nf01Blocked(t *testing.T, d Dependencies, s workflowState) {
	t.Helper()
	c, e := parallelConsumerClass(d, s, "")
	if e == nil && (c.Class == "PAR-A02" || c.Compatible) {
		t.Fatalf("incomplete/contradictory proof cleared authority: %+v", c)
	}
}
func nf01Envelope(t *testing.T, s *workflowState, change func(map[string]any)) {
	t.Helper()
	c := &s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	old := *c.Integration
	b, e := workflowBound(old, 4<<20)
	if e != nil {
		t.Fatal(e)
	}
	var v map[string]any
	if e = json.Unmarshal(b, &v); e != nil {
		t.Fatal(e)
	}
	change(v)
	path := writeAny(t, s.Request.WorkspaceRoot, "negative-"+strings.ReplaceAll(t.Name(), "/", "-")+".json", v)
	binding := FileBinding{path, hashFileTest(t, path)}
	c.Integration = &binding
	for i := range s.Result.Delivery.Events {
		ev := &s.Result.Delivery.Events[i]
		if equal(ev.Binding, old) {
			ev.Binding = binding
		}
	}
}

// These mutate captured inputs or negative copies of bound evidence only.
// No synthetic accepted/completed state is used to establish a positive case.
func TestNF01StrictCompatibilityNegatives(t *testing.T) {
	t.Run("report", func(t *testing.T) {
		f, o := nf01Report(t, true)
		s := nf01State(t, f, o)
		nf01Terminal(t, f, o)
		for _, which := range []string{"missing-retirement", "wrong-retirement-hash", "wrong-session", "missing-preparation", "missing-release", "other-release-preparation", "other-release-queue"} {
			t.Run(which, func(t *testing.T) {
				next := nf01Clone(t, s)
				d := f.D
				switch which {
				case "missing-retirement":
					next.ReportRetirement = &FileBinding{filepath.Join(f.R.WorkspaceRoot, "absent-retirement"), hash(nil)}
				case "wrong-retirement-hash":
					next.ReportRetirement.SHA256 = hash([]byte("wrong"))
				case "wrong-session":
					next.Result.Transport.AgentSessionID = "another-session"
				default:
					r, e := d.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
					if e != nil {
						t.Fatal(e)
					}
					if which == "missing-preparation" {
						r.TaskPreparations = nil
					} else {
						for i := range r.TaskQueueEvents {
							ev := &r.TaskQueueEvents[i]
							if ev.Kind == "release" || ev.Kind == "advance" {
								switch which {
								case "missing-release":
									ev.Kind = "reserve"
								case "other-release-preparation":
									ev.Request = map[string]any{"preparation_id": "other"}
								case "other-release-queue":
									ev.QueueID = "other"
								}
							}
						}
					}
					d.Workspace.WorkItems = nf01Snapshot{WorkItemStore: d.Workspace.WorkItems, registry: r}
				}
				nf01Blocked(t, d, next)
			})
		}
	})
	for _, route := range []string{"human", "headless"} {
		t.Run(route, func(t *testing.T) {
			var f deliveryFixture
			var o WorkflowRun
			if route == "human" {
				f, o = nf01Human(t, true)
			} else {
				f, o = nf01Headless(t, true)
			}
			s := nf01State(t, f.workflowFixture, o)
			nf01Terminal(t, f.workflowFixture, o)
			cases := []string{"missing-release", "wrong-release-hash", "unresolved-attempt", "phase-only", "wrong-result", "wrong-oid", "wrong-tree", "wrong-source", "missing-integration", "changed-integration", "wrong-event", "wrong-envelope-run", "wrong-envelope-request", "wrong-envelope-kind", "wrong-envelope-candidate", "envelope-error", "incomplete-recovery", "incomplete-result"}
			if route == "human" {
				cases = append(cases, "missing-human", "accepted-only", "missing-decision", "wrong-plan", "wrong-human-qa", "unregistered-integration", "wrong-authority-qa", "later-veto")
			}
			for _, which := range cases {
				t.Run(which, func(t *testing.T) {
					next := nf01Clone(t, s)
					d := f.D
					ds := next.Result.Delivery
					c := &ds.Candidates[len(ds.Candidates)-1]
					switch which {
					case "missing-release":
						ds.OwnershipRelease = nil
					case "wrong-release-hash":
						ds.OwnershipRelease.SHA256 = hash([]byte("bad-release"))
					case "unresolved-attempt":
						ds.Attempt = &DeliveryAttempt{ID: "unresolved", State: "attempted"}
					case "phase-only":
						ds.Phase = "working"
					case "wrong-result":
						c.TaskResult.ID = "other"
					case "wrong-oid":
						c.OID = strings.Repeat("a", 40)
					case "wrong-tree":
						c.Tree = strings.Repeat("a", 40)
					case "wrong-source":
						c.TaskResult.SourceLocator += "-other"
					case "missing-integration":
						c.Integration = nil
					case "changed-integration":
						c.Integration.SHA256 = hash([]byte("bad"))
					case "wrong-event":
						for i := range ds.Events {
							if ds.Events[i].Kind == "integration" {
								ds.Events[i].Kind = "report"
							}
						}
					case "missing-human":
						ds.HumanIntegration = nil
					case "accepted-only":
						ds.HumanIntegration.Accepted = false
					case "missing-decision":
						ds.HumanIntegration.Decision.Locator += "-absent"
					case "wrong-plan":
						ds.HumanIntegration.PlanSHA256 = hash([]byte("other-plan"))
					case "wrong-human-qa":
						c.HumanQA.ID = "other"
					case "unregistered-integration", "wrong-authority-qa", "later-veto":
						r, e := d.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
						if e != nil {
							t.Fatal(e)
						}
						switch which {
						case "unregistered-integration":
							r.IntegrationResults = nil
						case "wrong-authority-qa":
							r.IntegrationAuthorities[0].HumanQARecordID = "other"
						case "later-veto":
							q := *c.HumanQA
							q.ID = "qa_negative"
							q.Outcome = "fail"
							q.Actor.CompletedAtUTC = "2026-10-06T00:00:00Z"
							r.HumanQARecords = append(r.HumanQARecords, q)
						}
						d.Workspace.WorkItems = nf01Snapshot{WorkItemStore: d.Workspace.WorkItems, registry: r}
					default:
						nf01Envelope(t, &next, func(v map[string]any) {
							switch which {
							case "wrong-envelope-run":
								v["run_id"] = "other"
							case "wrong-envelope-request":
								v["request_sha256"] = hash([]byte("other"))
							case "wrong-envelope-kind":
								v["kind"] = "Other@1"
							case "wrong-envelope-candidate":
								v["candidate"] = "other"
							case "envelope-error":
								v["error"] = "failed"
							case "incomplete-result":
								v["result"].(map[string]any)["Completed"] = false
							case "incomplete-recovery":
								v["result"].(map[string]any)["Integration"].(map[string]any)["Readback"].(map[string]any)["RecoveryStatus"] = "unknown"
							}
						})
					}
					nf01Blocked(t, d, next)
				})
			}
			t.Run("missing-closeout", func(t *testing.T) {
				path := filepath.Join(f.R.WorkspaceRoot, ".ply", "task-closeouts.json")
				b, e := os.ReadFile(path)
				if e != nil {
					t.Fatal(e)
				}
				if e = os.Remove(path); e != nil {
					t.Fatal(e)
				}
				defer os.WriteFile(path, b, 0600)
				nf01Blocked(t, f.D, s)
			})
		})
	}
}

type nf01Snapshot struct {
	workspace.WorkItemStore
	registry workspace.WorkItemRegistry
}

func (s nf01Snapshot) Snapshot(string) (workspace.WorkItemRegistry, error) { return s.registry, nil }
func TestNF01VerifierExecutionSeam(t *testing.T) {
	for _, modeled := range []bool{false, true} {
		t.Run(fmt.Sprint(modeled), func(t *testing.T) {
			d := Dependencies{Now: time.Now}
			calls := 0
			if modeled {
				d.RunVerificationCommand = func(c *exec.Cmd) error {
					calls++
					if !equal(c.Args, []string{"/bin/sh", "-c", "exit 7"}) {
						t.Fatal(c.Args)
					}
					return c.Run()
				}
			}
			r := deliveryVerificationReceipt{Argv: []string{"/bin/sh", "-c", "exit 7"}, CWD: t.TempDir()}
			_, _, finish := runDeliveryVerification(d, &r)
			finish()
			if r.Exit == nil || *r.Exit != 7 || r.Error == "" || modeled && calls != 1 {
				t.Fatalf("lost actual waited command outcome: %+v calls=%d", r, calls)
			}
		})
	}
}

func TestNF01RemainingAuthority(t *testing.T) {
	a := workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalEpic}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a})
	nf01Runtime(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "remaining human authority")
	var e error
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "remaining"))
	if e != nil {
		t.Fatal(e)
	}
	for _, released := range []bool{false, true} {
		if released {
			if _, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Modeled release without integration"); e != nil {
				t.Fatal(e)
			}
		}
		s := nf01State(t, f.workflowFixture, o)
		want := "PAR-A01"
		if released {
			want = "PAR-A04"
		}
		for _, sha := range []string{"", s.Request.Runtime.PlyExecutable.SHA256} {
			c, e := parallelConsumerClass(f.D, s, sha)
			if e != nil || c.Class != want || c.Compatible != (sha != "") || len(c.RemainingEffects) == 0 {
				t.Fatalf("remaining native authority: %+v %v", c, e)
			}
		}
	}
}

// Capture complete disposable native stores, directories, source/control bytes,
// Git refs and history. Symlinks are compared as links, never traversed.
func nf01NativeState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, de os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := de.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if de.Type()&os.ModeSymlink != 0 {
			target, e := os.Readlink(path)
			if e != nil {
				return e
			}
			value += " " + target
		} else if !de.IsDir() {
			b, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			value += " " + hash(b)
		}
		out[path] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func TestNF01DamagedPreparationBinding(t *testing.T) {
	f, o := nf01Headless(t, false)
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	r.TaskPreparations[0].Plan.WorktreePath += "-different"
	before := nf01NativeState(t, f.R.WorkspaceRoot)
	facts, _, e := parallelActivationConsumers(f.D, f.R.WorkspaceRoot, r, "")
	if e != nil {
		t.Fatal(e)
	}
	run, prep := false, false
	for _, c := range facts {
		if c.ID == "run/"+o.RunID {
			run = c.Class == "unknown"
		}
		if c.ID == "preparation/"+f.R.PreparationID {
			prep = c.Class == "unknown"
		}
	}
	if !run || !prep || !equal(before, nf01NativeState(t, f.R.WorkspaceRoot)) {
		t.Fatalf("damaged binding suppressed unresolved preparation: %+v", facts)
	}
}

// Interrupt at the existing filesystem seam after native reservation bytes are
// durably saved, before WorkflowStart can call StartWorker. Bytes pass through
// unchanged; the injected failure is observed, never reported as successful
// supervisor startup. No detached process/provider exists in this fixture.
type nf01ReservationFiles struct {
	FileSystem
	captured bool
}

const nf01ReservationInterrupted = "modeled NF01 stop after native worker reservation, before launch"

func (f *nf01ReservationFiles) Replace(path string, raw []byte) error {
	if err := f.FileSystem.Replace(path, raw); err != nil {
		return err
	}
	var s workflowState
	if filepath.Base(path) == "state.json" && json.Unmarshal(raw, &s) == nil && headlessRun(s.Request) && s.Phase == "worker_reserved" && s.Result.Handoff.SHA256 != "" {
		f.captured = true
		return fmt.Errorf("%s", nf01ReservationInterrupted)
	}
	return nil
}
func nf01ReserveOnly(t *testing.T, f workflowFixture) WorkflowRun {
	t.Helper()
	files := &nf01ReservationFiles{FileSystem: f.D.files()}
	d := f.D
	d.Files = files
	p, e := WorkflowPreviewStart(d, f.File)
	if e != nil {
		t.Fatal(e)
	}
	preview := p.(WorkflowPreview)
	o, e := WorkflowStart(d, f.File, *preview.Confirmation)
	if e == nil || !strings.Contains(e.Error(), nf01ReservationInterrupted) || !files.captured || o.RunID == "" {
		t.Fatalf("native reserve-only boundary not reached: %+v %v", o, e)
	}
	s := nf01State(t, f, o)
	if s.Phase != "worker_reserved" || s.Result.Worker == nil || s.Result.Worker.OwnerPID != 0 {
		t.Fatal("fixture started a worker instead of capturing native reservation")
	}
	if e = privateDir(workerDirectory(s)); e != nil {
		t.Fatal(e)
	}
	return o
}

func TestNF01HumanHistoricalParent(t *testing.T) {
	f, o := nf01Human(t, true)
	runGit(t, f.Parent, "commit", "--allow-empty", "-m", "later disposable parent history")
	nf01Terminal(t, f.workflowFixture, o)
}
func TestNF01DamagedReportSlot(t *testing.T) {
	f, o := nf01Report(t, true)
	s := nf01State(t, f, o)
	if len(s.Records) == 0 {
		t.Fatal("post-start fixture lacks native report slot")
	}
	if e := os.WriteFile(s.Records[0].Report.Locator, []byte("damaged bound report"), 0600); e != nil {
		t.Fatal(e)
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	facts, _, e := parallelActivationConsumers(f.D, f.R.WorkspaceRoot, r, "")
	if e != nil {
		t.Fatal(e)
	}
	for _, c := range facts {
		if c.ID == "run/"+o.RunID {
			if c.Class != "unknown" || c.Compatible {
				t.Fatalf("damaged report slot became terminal: %+v", c)
			}
			return
		}
	}
	t.Fatal("damaged Run disappeared")
}

func TestNF01LiveAndSealedReportsRemainBlocked(t *testing.T) {
	f, o := nf01StartedReport(t, true)
	nf01Blocked(t, f.D, nf01State(t, f, o))
	o = workflowTestSettle(t, f, o)
	nf01Blocked(t, f.D, nf01State(t, f, o))
	var e error
	o, e = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if e != nil {
		t.Fatal(e)
	}
	nf01Blocked(t, f.D, nf01State(t, f, o))
}
