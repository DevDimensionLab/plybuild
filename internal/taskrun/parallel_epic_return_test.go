package taskrun

import (
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// All pane/session/writer observations below are simulated. Native registry,
// separate review sealing, executable acceptance and Git return are real.
func parallelRuntime(t *testing.T, f *deliveryFixture) {
	t.Helper()
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: 70001, CWD: f.Parent}, {PID: 70001, ParentPID: 1, CWD: f.Parent}}, nil
	}
	f.D.StartupProcessTree = func(pid int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
	}
	raw, err := os.ReadFile(f.Model)
	if err != nil {
		t.Fatal(err)
	}
	var model map[string]any
	if err = json.Unmarshal(raw, &model); err != nil {
		t.Fatal(err)
	}
	model["caller_shell"] = 70001
	model["review_session_prefix"] = f.R.PreparationID + "-"
	writeAny(t, filepath.Dir(f.Model), filepath.Base(f.Model), model)
}
func parallelSecondOwner(t *testing.T, a deliveryFixture) deliveryFixture {
	return parallelAdditionalOwner(t, a, true)
}

func parallelAdditionalOwner(t *testing.T, a deliveryFixture, parallel bool) deliveryFixture {
	t.Helper()
	root := a.R.WorkspaceRoot
	w := a.D.Workspace
	if _, e := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "task-b", Title: "B", Description: "Simulated independent owner", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	r, e := w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	raw, e := os.ReadFile(filepath.Join(root, "goal-draft.json"))
	if e != nil {
		t.Fatal(e)
	}
	var goal map[string]any
	if e = json.Unmarshal(raw, &goal); e != nil {
		t.Fatal(e)
	}
	goal["publication_key"] = "fixture/goal-b"
	goal["task_id"] = "task-b"
	goal["spec_id"] = "solution-b"
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == "task-b" {
			goal["problem"] = workspace.TaskRevisionRef{Revision: p.Revision, ManifestSHA256: p.ManifestSHA256}
		}
	}
	pub, e := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "task-b", File: writeAny(t, root, "goal-b-draft.json", goal)})
	if e != nil {
		t.Fatal(e)
	}
	g := workspace.TaskGoalRef{SpecID: "solution-b", Spec: workspace.TaskRevisionRef{Revision: *pub.OutcomeRef.Revision, ManifestSHA256: pub.OutcomeRef.ManifestSHA256}}
	q, e := workspace.ListTaskQueueVersion(w, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false, 2)
	if e != nil {
		t.Fatal(e)
	}
	queue := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/queue-b", ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: q.Revision, Entries: []workspace.TaskGoalQueueEntry{{TaskID: "task-b", Goal: &g}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}}
	if _, e = workspace.SetTaskQueue(w, writeAny(t, root, "queue-b.json", queue)); e != nil {
		t.Fatal(e)
	}
	acceptance := filepath.Join(root, "acceptance-b.sh")
	in := workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, SpecID: g.SpecID, Parallel: parallel, AcceptancePath: acceptance}
	preview, e := workspace.PreviewTaskGoalExecution(w, in)
	if e != nil {
		t.Fatal(e)
	}
	prepared, e := workspace.PrepareTaskGoalExecution(w, preview.Input, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "fixture only", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Simulated start authority, not human product QA"})
	if e != nil {
		t.Fatal(e)
	}
	d := SystemDependencies(w)
	prep := *prepared.Preparation.Preparation
	draft, e := workflowhandoff.BuildDeliveryHandoffDraft(d.Workflow, prep.ID, "fixture human", "fixture owner B", "fixture/delivery-b", acceptance, workflowhandoff.DeliveryDraftAuthority{StartupProtocol: MandateReceiptV1, Agreement: prepared.Delivery})
	if e != nil {
		t.Fatal(e)
	}
	directory := filepath.Join(root, "b-runtime")
	if e = os.Mkdir(directory, 0700); e != nil {
		t.Fatal(e)
	}
	herdr, e := os.ReadFile(a.R.Herdr.Executable.Path)
	if e != nil {
		t.Fatal(e)
	}
	herdr = []byte(strings.ReplaceAll(string(herdr), "w-fixture:p1", "w-fixture:p2"))
	path := filepath.Join(directory, "herdr")
	if e = os.WriteFile(path, herdr, 0700); e != nil {
		t.Fatal(e)
	}
	admission, reconciliation := "", ""
	if parallel {
		admission, reconciliation = workspace.AdmissionParallel, workspace.SourceReconciliationMergeParent
	}
	req, e := BuildDeliveryWorkflowRequest(DeliveryRequestInput{RequestKey: "fixture/delivery-b", WorkspaceRoot: root, PreparationID: prep.ID, PreparationSHA256: digest(prep), HandoffDraft: draft, Runtime: a.R.Runtime, HumanAuthority: a.R.HumanAuthority, HerdrExecutable: Executable{path, hash(herdr)}, HerdrWorkspaceID: "w-fixture", TabLabel: "B", Delivery: DeliveryContract{StartupProtocol: MandateReceiptV1, AdmissionMode: admission, SourceReconciliation: reconciliation, OwnerClaim: "fixture owner B", Goal: FileBinding{filepath.Join(root, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(g.Spec.ManifestSHA256, "sha256:")+".json"), g.Spec.ManifestSHA256}, AcceptancePath: acceptance, LocalIntegration: DeliveryLocalIntegrationBoundary(prepared.Delivery), Agreement: prepared.Delivery}})
	if e != nil {
		t.Fatal(e)
	}
	model := writeAny(t, directory, "herdr-model.json", map[string]any{"agent_status": "idle", "agent_session_id": "fixture-session-b", "prompt_mode": "normal", "provider": "codex", "foreground_cwd": prep.Plan.WorktreePath})
	d.CWD = func() (string, error) { return prep.Plan.WorktreePath, nil }
	d.Executable = a.D.Executable
	f := deliveryFixture{workflowFixture: workflowFixture{D: d, R: req, File: writeAny(t, root, "delivery-b-request.json", req), Model: model, Calls: filepath.Join(directory, "herdr-calls.jsonl")}, Prepared: prepared, Parent: a.Parent, AcceptancePath: acceptance}
	parallelRuntime(t, &f)
	return f
}
func parallelQualify(t *testing.T, f deliveryFixture, o WorkflowRun, binary string) WorkflowRun {
	t.Helper()
	defer func() {
		if t.Failed() {
			raw, _ := os.ReadFile(f.Model)
			t.Logf("simulated model: %s", raw)
			raw, _ = os.ReadFile(f.Calls)
			t.Logf("simulated calls: %s", raw)
		}
	}()
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, e := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic parallel coordinator"})
	if e != nil {
		t.Fatalf("separate review/verify step: %v", e)
	}
	round := reviewRoundLatest(t, requested)
	reviewRoundIn(t, round.ReviewerPane)
	_, e = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, reviewRoundReturnFile(t, f, "parallel-review-"+o.RunID+".json", `{"decision":"passed","summary":"Simulated independent exact candidate review","findings":[],"evidence_gaps":[]}`))
	if e != nil {
		t.Fatalf("separate review/verify step: %v", e)
	}
	reviewRoundIn(t, o.Transport.PaneID)
	result, e := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary)
	if e != nil {
		t.Fatalf("separate review/verify step: %v", e)
	}
	return result
}
func TestParallelExactTwoQualifiedSerialReturns(t *testing.T) {
	parallelExactSerialReturns(t, false)
}

func TestParallelReturnRunPublicationRecovery(t *testing.T) {
	parallelExactSerialReturns(t, true)
}

func parallelExactSerialReturns(t *testing.T, interrupt bool) {
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic parallel coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, Parallel: true})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	first := workflowTestStart(t, f.workflowFixture)
	first, e := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, first.RunID, first.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: first.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	b := parallelSecondOwner(t, f)
	second := workflowTestStart(t, b.workflowFixture)
	second, e = WorkflowDeliveryReceipt(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: second.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	if first.RunID == second.RunID || first.Transport.AgentSessionID == second.Transport.AgentSessionID {
		t.Fatal("owners share Run/session")
	}
	binary := filepath.Join(f.R.WorkspaceRoot, "simulated-candidate")
	if e = os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if exact := os.Getenv("PLY_DELIVERY_TEST_BINARY"); exact != "" {
		if !filepath.IsAbs(exact) {
			t.Fatal("exact candidate must be absolute")
		}
		binary = exact
	}
	for _, entry := range []struct {
		f    deliveryFixture
		name string
	}{{f, "a.txt"}, {b, "b.txt"}} {
		source := entry.f.Prepared.Preparation.Plan.WorktreePath
		if e = os.WriteFile(filepath.Join(source, entry.name), []byte("fixed "+entry.name+"\n"), 0600); e != nil {
			t.Fatal(e)
		}
		runGit(t, source, "add", entry.name)
		runGit(t, source, "commit", "-m", "fixed "+entry.name)
		candidateCheck := "test \"$(\"$PLY_CANDIDATE_BINARY\")\" = 'fixture candidate'\n"
		if os.Getenv("PLY_DELIVERY_TEST_BINARY") != "" {
			candidateCheck = "\"$PLY_CANDIDATE_BINARY\" version --format json\n"
		}
		if e = os.WriteFile(entry.f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ngrep -q fixed "+entry.name+"\n"+candidateCheck), 0600); e != nil {
			t.Fatal(e)
		}
	}
	first = parallelQualify(t, f, first, binary)
	second = parallelQualify(t, b, second, binary)
	firstAuthority, e := nativeDeliveryCandidate(f.D, f.R.WorkspaceRoot, first.RunID, first.Delivery.Candidates[0].TaskResult.ID)
	if e != nil {
		t.Fatal(e)
	}
	reviewRoundIn(t, first.Transport.PaneID)
	if interrupt {
		f.D.Fault = func(at string) error {
			if at == "local_delivery_after_native_return" {
				return errors.New("simulated Run publication interruption")
			}
			return nil
		}
	}
	delivered, e := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, first.RunID, first.Paths.Context)
	if !interrupt && e != nil && !parallelOwnerHandoff(e) {
		t.Fatalf("A return: %v", e)
	}
	if !interrupt && delivered.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("A return not observed: %+v", delivered.DeliveryStatus)
	}
	if interrupt && (e == nil || !strings.Contains(e.Error(), "simulated Run publication")) {
		t.Fatalf("Run interruption not observed: %v", e)
	}
	f.D.Fault = nil
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	reviewRoundIn(t, second.Transport.PaneID)
	if _, e = WorkflowDeliveryIntegrate(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context); e == nil {
		t.Fatal("B reused old parent proof")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatal("stale B return changed target")
	}
	preview, e := WorkflowDeliveryReconcileSource(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	merged, e := WorkflowDeliveryReconcileSource(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: preview.Confirmation})
	if e != nil || merged.State != "committed" {
		t.Fatalf("native reconcile: %+v %v", merged, e)
	}
	second = parallelQualify(t, b, second, binary)
	delivered, e = WorkflowDeliveryIntegrate(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context)
	if e != nil && !parallelOwnerHandoff(e) {
		t.Fatalf("B exact return: %v", e)
	}
	if delivered.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("B return not observed: %+v", delivered.DeliveryStatus)
	}
	r, e := b.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.IntegrationResults) != 2 {
		t.Fatalf("wrong integration count: %d", len(r.IntegrationResults))
	}
	q, e := workspace.ListTaskQueueVersion(b.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false, 2)
	if e != nil || len(q.Active) != 0 {
		t.Fatalf("named closure: %+v %v", q, e)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if _, e = os.Stat(filepath.Join(f.Parent, name)); e != nil {
			t.Fatal(e)
		}
	}
	// Model a lost Run publication after A's native integration/queue/base
	// closure. B has now lawfully advanced the parent; observing A's completed
	// native return must not demand the old live parent or repeat any effect.
	auth := firstAuthority
	recovered, e := workflowhandoff.IntegrateDeliveryCandidate(f.D.Workflow, string(auth.Candidate.TaskResult.TaskID), auth.Candidate.TaskResult.ID, "", auth.ExpectedParentOID, auth.Candidate.OID, workflowhandoff.DeliveryOwnerAuthority{RunID: first.RunID, RequestSHA256: auth.RequestSHA256, ActorClaim: auth.OwnerClaim, PreparationID: auth.PreparationID, Authorization: auth.Authorization})
	if e != nil || !recovered.Completed {
		t.Fatalf("completed A return lost after B advanced: %v", e)
	}
	r, e = f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || len(r.IntegrationResults) != 2 || gitOutput(t, f.Parent, "rev-parse", "HEAD") != delivered.Delivery.Candidates[len(delivered.Delivery.Candidates)-1].OID {
		t.Fatal("historical return recovery repeated or rewound effect")
	}
	if interrupt {
		reviewRoundIn(t, first.Transport.PaneID)
		again, e := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, first.RunID, first.Paths.Context)
		if e != nil && !parallelOwnerHandoff(e) {
			t.Fatalf("historical Run publication recovery: %v", e)
		}
		if again.DeliveryStatus.FinalDelivery.State != "delivered" {
			t.Fatal("completed native A return not published after B")
		}
	}
}

func parallelOwnerHandoff(err error) bool { var e *Error; return errors.As(err, &e) && e.Exit == 5 }
