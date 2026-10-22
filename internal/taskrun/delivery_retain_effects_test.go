package taskrun

import (
	"bytes"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
)

func TestAL40RetainsUncertainEffects(t *testing.T) {
	for _, effect := range []string{DeliveryEffectPartial, DeliveryEffectUnknown} {
		t.Run(effect, func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, effect)
			owner := f.D
			var warnings bytes.Buffer
			oldOutput := logrus.StandardLogger().Out
			logrus.SetOutput(&warnings)
			defer logrus.SetOutput(oldOutput)
			sourceFile := f.Prepared.Preparation.Preparation.Plan.WorktreePath + "/unique-source"
			if err := os.WriteFile(sourceFile, []byte("unique retained bytes"), 0600); err != nil {
				t.Fatal(err)
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
			f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
			retireAsCoordinator(&f)
			out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 modeled coordinator", Reason: "Resource ownership only", RetainEffects: true})
			if err != nil {
				t.Fatal(err)
			}
			if out.Delivery.Phase != "stopped" || out.Delivery.ActiveStop == nil {
				t.Fatal("stop lost")
			}
			if retireQueueCurrent(t, f) != nil {
				t.Fatal("owned queue not released")
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
				t.Fatal("closed pane")
			}
			if !strings.Contains(warnings.String(), "original_effect="+effect) || !strings.Contains(warnings.String(), "retained_record=") || !strings.Contains(warnings.String(), "source=") {
				t.Fatal("uncertain effects warning missing", warnings.String())
			}
			status := out.DeliveryStatus.ResourceRetirement
			if status.State != "retired" || status.Source.Clean || status.RetainedEffects[0].OriginalEffect != effect || out.DeliveryStatus.FinalDelivery.State != "not_delivered" {
				t.Fatalf("retention readback: %+v", status)
			}
			if bytes, err := os.ReadFile(sourceFile); err != nil || string(bytes) != "unique retained bytes" {
				t.Fatal("source changed", err)
			}
			for _, callback := range []func() (WorkflowRun, error){
				func() (WorkflowRun, error) {
					return WorkflowDeliveryVerify(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "/never-read-review")
				},
				func() (WorkflowRun, error) {
					return WorkflowDeliveryIntegrate(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
				},
				func() (WorkflowRun, error) {
					return WorkflowDeliveryReceipt(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: f.R.Runtime.Provider, NativeSession: o.Transport.AgentSessionID}, true)
				},
			} {
				if _, err := callback(); err == nil || !strings.Contains(err.Error(), "owner_retired") {
					t.Fatalf("late callback not rejected: %v", err)
				}
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = ownerResumePreview(f.D, f.R.WorkspaceRoot, DeliveryOwnerResumeInput{}, s); err == nil || !strings.Contains(err.Error(), "owner_retired") {
				t.Fatal("resume not rejected", err)
			}
		})
	}
}

func TestAL40BlockedBeforeResourceEffects(t *testing.T) {
	for _, tc := range []string{"live_pane", "local_writer", "process_unknown", "native_attempt", "resume_attempt", "released_owner", "completed", "evidence_corrupt", "source_drift", "foreign_stop", "source_ref", "empty_process_table"} {
		t.Run(tc, func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectPartial)
			workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": tc != "live_pane"})
			calls := 0
			f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
				calls++
				if tc == "empty_process_table" {
					return nil, nil
				}
				if tc == "process_unknown" {
					return nil, os.ErrPermission
				}
				if tc == "local_writer" {
					return []LocalProcess{{PID: os.Getpid()}, {PID: 424242, CWD: f.Prepared.Preparation.Preparation.Plan.WorktreePath}}, nil
				}
				if tc == "source_drift" {
					if err := os.WriteFile(f.Prepared.Preparation.Preparation.Plan.WorktreePath+"/drift", []byte("changed"), 0600); err != nil {
						t.Fatal(err)
					}
				}
				return []LocalProcess{{PID: os.Getpid()}}, nil
			}
			if tc == "native_attempt" || tc == "resume_attempt" || tc == "released_owner" || tc == "completed" {
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					switch tc {
					case "native_attempt":
						s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "integration", State: "attempted"}
					case "resume_attempt":
						s.OwnerResume = &FileBinding{}
						s.OwnerResumeStep = ownerResumeTabAttempted
					case "released_owner":
						s.Result.Delivery.OwnershipRelease = &FileBinding{}
					case "completed":
						s.Result.Delivery.Phase = "completed"
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			if tc == "evidence_corrupt" {
				s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				raw, err := workflowBound(*s.Result.Delivery.ActiveStop, 2<<20)
				if err != nil {
					t.Fatal(err)
				}
				var stop DeliveryStopRecord
				if err = decode(raw, 2<<20, &stop); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(stop.Evidence.Locator, []byte("damaged"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if tc == "source_ref" {
				runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "branch", "-m", "foreign-source-ref")
			}
			if tc == "foreign_stop" {
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					proof := writeAny(t, f.R.WorkspaceRoot, "foreign-stop-evidence.json", map[string]string{"modeled": "foreign owner"})
					b := FileBinding{proof, hashFileTest(t, proof)}
					record := DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: o.RunID, Input: DeliveryStopInput{EventID: "foreign-stop", Provider: f.R.Runtime.Provider, NativeSession: "foreign-session", EffectStatus: DeliveryEffectUnknown}, Source: b, Evidence: b}
					_, err := deliveryAppendEvent(f.D, s, "foreign-stop", "coordinator_stop", record)
					return err
				}); err != nil {
					t.Fatal(err)
				}
			}
			retireAsCoordinator(&f)
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 fixture", Reason: "Retain only", RetainEffects: true}); err == nil {
				t.Fatal("blocker accepted")
			}
			if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
				t.Fatal("closed owner")
			}
			if tc == "live_pane" || tc == "native_attempt" || tc == "resume_attempt" || tc == "released_owner" || tc == "completed" || tc == "evidence_corrupt" || tc == "foreign_stop" || tc == "source_ref" {
				if calls != 0 {
					t.Fatal("inspected local processes past known blocker")
				}
			}
			if tc == "process_unknown" && calls != 1 {
				t.Fatal("replayed inspection")
			}
			if retireQueueCurrent(t, f) == nil {
				t.Fatal("released blocked queue")
			}
		})
	}
}

func TestAL40PendingReleaseSettledSourceAndDurableRecovery(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectUnknown)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	observations := 0
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		observations++
		return []LocalProcess{{PID: os.Getpid()}}, nil
	}
	retireAsCoordinator(&f)
	in := DeliveryRetireInput{Observer: "AL40 modeled root", Reason: "Resource closure only", RetainEffects: true}
	fault := f.D
	fault.Fault = func(stage string) error {
		if stage == "retain_effects_after_retirement" {
			return os.ErrInvalid
		}
		return nil
	}
	if _, err := WorkflowRetireDelivery(fault, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("fault absent")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	r, original, _, _, err := deliveryRetentionHistory(s)
	if err != nil || r == nil {
		t.Fatal(err)
	}
	raw, err := workflowBound(original, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	work := f.Prepared.Preparation.Preparation.Plan.WorktreePath + "/settled-source"
	if err = os.WriteFile(work, []byte("settled unique work"), 0600); err != nil {
		t.Fatal(err)
	}
	runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "add", "settled-source")
	runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "commit", "-m", "Disposable settled source")
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	// A fresh bounded observation can reject a writer without changing history.
	blocked := f.D
	blocked.LocalProcesses = func(string, string) ([]LocalProcess, error) { return nil, os.ErrPermission }
	if _, err = WorkflowRetireDelivery(blocked, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("unknown writer accepted")
	}
	if retireQueueCurrent(t, f) == nil {
		t.Fatal("blocked retry released")
	}
	drifting := f.D
	drifting.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		if err := os.WriteFile(work, []byte("later settled bytes"), 0600); err != nil {
			t.Fatal(err)
		}
		return []LocalProcess{{PID: os.Getpid()}}, nil
	}
	if _, err = WorkflowRetireDelivery(drifting, f.R.WorkspaceRoot, o.RunID, in); err == nil || !strings.Contains(err.Error(), "source changed") {
		t.Fatal("pending within-call drift did not block", err)
	}
	if retireQueueCurrent(t, f) == nil {
		t.Fatal("drifting retry released")
	}
	// Fault after append but before indexing must recover one immutable observation.
	fault = f.D
	fault.Fault = func(stage string) error {
		if stage == "retain_effects_after_observation" {
			return os.ErrInvalid
		}
		return nil
	}
	if _, err = WorkflowRetireDelivery(fault, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("observation fault absent")
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, retained, bindings, err := deliveryRetentionHistory(s)
	if err != nil || len(retained) != 1 {
		t.Fatalf("observation recovery: %v %d", err, len(retained))
	}
	if retained[0].Source.SourceSHA256 == r.DeliverySource.SourceSHA256 || retained[0].Effects[0].OriginalEffect != DeliveryEffectUnknown || len(retained[0].Effects) != 2 || retained[0].Source.OID == r.DeliverySource.OID {
		t.Fatal("settled source or unknown lost")
	}
	b := bindings[0]
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	if retireQueueCurrent(t, f) != nil {
		t.Fatal("queue still held")
	}
	originalRaw, err := workflowBound(original, 1<<20)
	if err != nil || !equal(raw, originalRaw) {
		t.Fatal("original retirement changed", err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, retained, bindings, err = deliveryRetentionHistory(s)
	if err != nil || len(retained) != 1 || bindings[0] != b {
		t.Fatal("observation rewritten", err)
	}
	// Completed retries need no live source or writer grant and keep history fixed.
	count := observations
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		t.Fatal("historical retry inspected processes")
		return nil, nil
	}
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatal(err)
	}
	if observations != count || len(out.Delivery.Events) != 4 {
		t.Fatal("retry changed history")
	}
	in.Reason = "different"
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err == nil {
		t.Fatal("changed reason accepted")
	}
}

func TestAL40ClosedPreparationAndLostReleaseReply(t *testing.T) {
	for _, closedFirst := range []bool{false, true} {
		t.Run(map[bool]string{true: "already_released", false: "lost_release_reply"}[closedFirst], func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectPartial)
			workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
			f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
			retireAsCoordinator(&f)
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			p, err := workflowReportPreparation(f.D, s)
			if err != nil {
				t.Fatal(err)
			}
			target := p.Plan.Target
			qt := workspace.QueueTargetInput{ProjectID: target.ProjectID, RepoID: target.RepoID, EpicID: target.EpicID}
			if closedFirst {
				q, err := workspace.ListTaskQueue(f.D.Workspace, qt, false)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = workspace.CloseTaskQueue(f.D.Workspace, qt, f.R.PreparationID, q.Revision, "fixture preclosed exact owner", "release"); err != nil {
					t.Fatal(err)
				}
			}
			if closedFirst {
				integrationCLIPrepareSuccessor(t, f)
			}
			before, err := workspace.ListTaskQueue(f.D.Workspace, qt, false)
			if err != nil {
				t.Fatal(err)
			}
			in := DeliveryRetireInput{Observer: "AL40 root", Reason: "Retain all effects", RetainEffects: true}
			d := f.D
			if !closedFirst {
				d.Fault = func(stage string) error {
					if stage == "retain_effects_after_release" {
						return os.ErrInvalid
					}
					return nil
				}
			}
			_, err = WorkflowRetireDelivery(d, f.R.WorkspaceRoot, o.RunID, in)
			if closedFirst && err != nil {
				t.Fatal(err)
			}
			if !closedFirst && err == nil {
				t.Fatal("lost-reply fault absent")
			}
			after, err := workspace.ListTaskQueue(f.D.Workspace, qt, false)
			if err != nil {
				t.Fatal(err)
			}
			f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
				t.Fatal("closed preparation triggered new observation")
				return nil, nil
			}
			if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
				t.Fatal(err)
			}
			again, err := workspace.ListTaskQueue(f.D.Workspace, qt, false)
			if err != nil {
				t.Fatal(err)
			}
			if !equal(after, again) || closedFirst && !equal(before, after) {
				t.Fatal("historical queue changed")
			}
		})
	}
}

func TestAL40ReportAliasRefusesRetainEffects(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestStart(t, f)
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		t.Fatal("unsupported report inspected processes")
		return nil, nil
	}
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 modeled root", Reason: "Unsupported explicit mode", RetainEffects: true}); err == nil || !strings.Contains(err.Error(), "local Epic") {
		t.Fatal("report silently accepted mode", err)
	}
}

func TestAL40LegacyStartAndEarlierUncertaintyRemainHistory(t *testing.T) {
	f := newDeliveryFixture(t, "claude", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic}})
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "legacy-retain-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	o = retireStopFixture(t, f, o, DeliveryEffectUnknown)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	acceptance := hashFileTest(t, before.Acceptance.Locator)
	retireAsCoordinator(&f)
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 modeled root", Reason: "Legacy resource closure", RetainEffects: true})
	if err != nil {
		t.Fatal(err)
	}
	effects := out.DeliveryStatus.ResourceRetirement.RetainedEffects
	if len(effects) != 2 || effects[0].OriginalEffect != DeliveryEffectUnknown || effects[1].OriginalEffect != DeliveryEffectNotStarted {
		t.Fatal("earlier uncertainty lost")
	}
	if hashFileTest(t, before.Acceptance.Locator) != acceptance || out.Delivery.PermissionState != before.Result.Delivery.PermissionState {
		t.Fatal("legacy permission history changed")
	}
}

func TestAL40NativeStartSessionCannotBeRebound(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectPartial)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	retireAsCoordinator(&f)
	s.Result.Transport.AgentSessionID = "foreign-session"
	if err = deliveryRetentionGuard(f.D, s); err == nil {
		t.Fatal("native start session was rebound")
	}
}
