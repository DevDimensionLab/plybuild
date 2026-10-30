package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAL01OrdinaryReportRequestSelectsReceipt(t *testing.T) {
	f := workflowTestFixture(t)
	var fields map[string]any
	b, _ := os.ReadFile(f.File)
	json.Unmarshal(b, &fields)
	delete(fields, "startup_protocol")
	r, err := ReadWorkflowRequest(writeAny(t, f.R.WorkspaceRoot, "ordinary.json", fields))
	if err != nil {
		t.Fatal(err)
	}
	b, _ = Canonical(r)
	var got map[string]any
	json.Unmarshal(b, &got)
	if got["startup_protocol"] != MandateReceiptV1 {
		t.Fatal("ordinary report request still selects legacy attestation")
	}
}

func al01StopReport(f workflowFixture) WorkflowReport {
	r := f.Report
	r.Outcome = "blocked"
	r.BudgetUsage = &BudgetUsage{}
	r.StopReasons = json.RawMessage(`[{"type":"scope_or_authority_expansion","primary":true,"detail":"Profile and policy are unobservable before work."}]`)
	r.ObservedEffects = json.RawMessage(`[]`)
	r.VerifierResults = json.RawMessage(`[]`)
	r.Artifacts = json.RawMessage(`[]`)
	r.TechnicalAssessment.Gate = "unknown"
	r.TechnicalAssessment.RequiredVerifierIDs = []string{}
	return r
}

func TestAL01PreAcceptanceStop(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestStart(t, f)
	r := al01StopReport(f)
	_, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "pre-stop.json", r))
	if err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.StartSHA256 != nil || s.Acceptance != nil || len(s.Records) != 0 {
		t.Fatal("stop manufactured work acceptance or terminal")
	}
}

func TestAL01LegacyRetirementSurvivesCLIUpgrade(t *testing.T) {
	f := workflowTestFixture(t)
	ply := filepath.Join(f.R.WorkspaceRoot, "historical-ply")
	os.WriteFile(ply, []byte("old ply"), 0700)
	f.R.Runtime.PlyExecutable = Executable{ply, hashFileTest(t, ply)}
	f.D.Executable = func() (string, error) { return ply, nil }
	herdr, _ := os.ReadFile(f.R.Herdr.Executable.Path)
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700)
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "legacy-retire-request.json", f.R)
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}}, nil
	}
	o := workflowTestStart(t, f)
	a := f.Acceptance
	a.RequestSHA256 = o.RequestSHA256
	a.Acceptance = "unknown"
	a.RuntimeClaim = RuntimeClaim{NativeSessionID: a.RuntimeClaim.NativeSessionID}
	a.Sandbox = json.RawMessage(`null`)
	a.Issues = json.RawMessage(`[{"type":"unknown","detail":"Policy unknown"}]`)
	if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "negative.json", a)); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(ply, []byte("upgraded ply"), 0700)
	workflowTestModel(t, f, map[string]any{"closed_transport": true})
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Replace clean negatively accepted owner after upgrade"}); err != nil {
		t.Fatal(err)
	}
}

func al01ReceiptFixture(t *testing.T, binary string) workflowFixture {
	t.Helper()
	return al01ReceiptFixtureFrom(t, workflowTestFixture(t), binary)
}

func al01ReceiptFixtureFrom(t *testing.T, f workflowFixture, binary string) workflowFixture {
	t.Helper()
	f.R.StartupProtocol = ""
	f.R.RequestKey += "-receipt"
	var draft map[string]any
	json.Unmarshal(f.R.HandoffDraft, &draft)
	draft["publication_key"] = "al01/report-receipt"
	f.R.HandoffDraft, _ = Canonical(draft)
	if binary != "" {
		f.R.Runtime.PlyExecutable = Executable{binary, hashFileTest(t, binary)}
		f.D.Executable = func() (string, error) { return binary, nil }
	}
	f.File = writeAny(t, f.R.WorkspaceRoot, "ordinary-receipt.json", f.R)
	r, err := ReadWorkflowRequest(f.File)
	if err != nil {
		t.Fatal(err)
	}
	f.R = r
	f.Report.RunID, f.Report.RequestSHA256, f.Report.SessionID = workflowID(r), digest(r), "ply:"+workflowID(r)
	created, err := workflowhandoff.Create(f.D.Workflow, workflowhandoff.CreateInput{DraftPath: writeAny(t, f.R.WorkspaceRoot, "receipt-handoff.json", json.RawMessage(r.HandoffDraft))})
	if err != nil {
		t.Fatal(err)
	}
	f.Handoff = created.Locator
	view, err := workflowhandoff.TaskRunHandoffView(f.D.Workflow, f.Handoff)
	if err != nil {
		t.Fatal(err)
	}
	var h struct {
		Reply struct {
			Root string `json:"reply_root"`
		} `json:"reply_capability"`
	}
	json.Unmarshal(view, &h)
	staging := filepath.Join(h.Reply.Root, "staging")
	if err = privateDir(staging); err != nil {
		t.Fatal(err)
	}
	var artifacts []map[string]any
	json.Unmarshal(f.Report.Artifacts, &artifacts)
	for _, a := range artifacts {
		old := a["locator"].(string)
		b, _ := os.ReadFile(old)
		path := filepath.Join(staging, filepath.Base(old))
		if err = os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
		a["locator"] = path
	}
	f.Report.Artifacts, _ = Canonical(artifacts)
	return f
}

func al01CLIOK(t *testing.T, binary string, f workflowFixture, args ...string) WorkflowRun {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir, _ = f.D.CWD()
	out, stderr, err := al14Command(t, c)
	if err != nil {
		t.Fatalf("CLI %v: %v %s %s", args, err, out, stderr)
	}
	var o WorkflowRun
	if err = json.Unmarshal(out, &o); err != nil {
		t.Fatal(err)
	}
	return o
}

func TestAL01ReportReceiptUnknownFactsAndCompletedReview(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := workflowTestStart(t, f)
	in := DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}
	first, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true)
	if err != nil {
		t.Fatal(err)
	}
	if first.Delivery != nil || first.RuntimeFacts.ModelState != "unknown" {
		t.Fatal("receipt changed authority or invented model")
	}
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	raw, _ := workflowBound(*s.Acceptance, 1<<20)
	var a Acceptance
	json.Unmarshal(raw, &a)
	var sandbox map[string]any
	json.Unmarshal(a.Sandbox, &sandbox)
	if a.RuntimeClaim.ProfileID != nil || a.RuntimeClaim.EffectivePolicySHA256 != nil || sandbox["matches_contract"] != nil || s.StartSHA256 == nil {
		t.Fatalf("unknown facts became attested: %s", raw)
	}
	if _, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, in, true); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []DeliveryBeginInput{{Provider: "claude", NativeSession: in.NativeSession}, {Provider: "codex", NativeSession: "wrong"}, {Provider: "codex", NativeSession: in.NativeSession, PermissionConfirmed: true}, {Provider: "codex", NativeSession: in.NativeSession, Profile: "changed"}} {
		if _, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, bad, true); err == nil {
			t.Fatal("wrong/changed receipt accepted")
		}
	}
	o = workflowTestSettle(t, f, workflowTestReport(t, f, first))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || o.FinalReturn.TerminalSHA256 == nil || o.TaskResultState != "not_published" || o.Delivery != nil {
		t.Fatalf("report review lost meaning: %+v %v", o, err)
	}
}

func TestAL01PreStartStopIsBoundImmutableAndNonQualifying(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := workflowTestStart(t, f)
	r := al01StopReport(f)
	for _, mutate := range []func(*WorkflowReport){func(r *WorkflowReport) { r.SessionID = "wrong" }, func(r *WorkflowReport) { r.RequestSHA256 = "sha256:" + strings.Repeat("0", 64) }, func(r *WorkflowReport) { r.Outcome = "complete" }, func(r *WorkflowReport) { r.ObservedEffects = f.Report.ObservedEffects }, func(r *WorkflowReport) { r.VerifierResults = f.Report.VerifierResults }, func(r *WorkflowReport) { r.TechnicalAssessment.Gate = "passed" }} {
		bad := r
		mutate(&bad)
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "bad-stop.json", bad)); err == nil {
			t.Fatal("pre-start stop manufactured authority")
		}
	}
	p := writeAny(t, f.R.WorkspaceRoot, "stop.json", r)
	for i := 0; i < 2; i++ {
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, p); err != nil {
			t.Fatal(err)
		}
	}
	r.Summary = "Changed observation"
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "changed-stop.json", r)); err == nil {
		t.Fatal("stop changed")
	}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true); err == nil {
		t.Fatal("stop cleared by receipt")
	}
}

func TestAL01ReportCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	f := al01ReceiptFixture(t, binary)
	o := workflowTestStart(t, f)
	args := []string{"run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID}
	for i := 0; i < 2; i++ {
		o = al01CLIOK(t, binary, f, args...)
	}
	if o.Delivery != nil || o.RuntimeFacts.ModelState != "unknown" {
		t.Fatal("CLI report receipt fabricated delivery or model")
	}
	stop := al01StopReport(f)
	// A separate synthetic pre-start report Run exercises public stop transport.
	stopFixture := workflowTestFixture(t)
	stopFixture.R.Runtime.PlyExecutable = Executable{binary, hashFileTest(t, binary)}
	stopFixture.File = writeAny(t, stopFixture.R.WorkspaceRoot, "cli-stop-request.json", stopFixture.R)
	started := workflowTestStart(t, stopFixture)
	stop.RunID, stop.RequestSHA256, stop.SessionID = started.RunID, started.RequestSHA256, started.SessionID
	file := writeAny(t, stopFixture.R.WorkspaceRoot, "cli-stop.json", stop)
	stopped := al01CLIOK(t, binary, stopFixture, "run", "report", started.RunID, "--context", started.Paths.Context, "--file", file)
	if stopped.Round.State != "stopped" || stopped.FinalReturn.TerminalSHA256 != nil {
		t.Fatal("CLI pre-start stop claimed a native result")
	}
}

func al01PreviewCLI(t *testing.T, binary, cwd, file string) WorkflowPreview {
	t.Helper()
	c := exec.Command(binary, "workflow", "run", "start", "--file", file, "--check", "--format", "json")
	c.Dir = cwd
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("preview: %v %s", err, out)
	}
	var p WorkflowPreview
	if err = json.Unmarshal(out, &p); err != nil {
		t.Fatal(err)
	}
	return p
}

func mustReadAL01(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func al01QueueReplacementTask(t *testing.T, f workflowFixture) {
	t.Helper()
	root := f.R.WorkspaceRoot
	w := f.D.Workspace
	replacementDraft := func(t *testing.T, root, kind, key string, fields map[string]any) workspace.TaskContentInput {
		in := taskSpecFixtureDraft(t, root, kind, key, fields)
		var v map[string]any
		json.Unmarshal(mustReadAL01(t, in.File), &v)
		v["task_id"] = "replacement"
		in.TaskID = "replacement"
		in.File = writeAny(t, root, "replacement-"+filepath.Base(in.File), v)
		return in
	}
	if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "replacement", Title: "Fixture Task", Description: "Fixture need", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	recordFixtureProblem(t, w, root, "replacement")
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	parent := r.Epics[0].RepoBindings[0].Worktree
	problem := workspace.TaskRevisionRef{}
	for _, revision := range r.TaskProblemRevisions {
		if revision.TaskID == "replacement" {
			problem = workspace.TaskRevisionRef{Revision: revision.Revision, ManifestSHA256: revision.ManifestSHA256}
		}
	}
	doc := filepath.Join(root, "solution.md")
	b := []byte("# Fixture solution\n\nScope README.md; verify test; preserve snapshots on failure.\n")
	if err = os.WriteFile(doc, b, 0600); err != nil {
		t.Fatal(err)
	}
	refs := []any{map[string]any{"document_id": "solution", "section": nil}}
	part := map[string]any{"state": "present", "reason": nil, "documents": refs}
	spec, err := workspace.RecordTaskSpec(w, replacementDraft(t, root, "WorkspaceTaskSpecDraft@1", "replacement/spec", map[string]any{
		"spec_id": "solution", "expected_previous": nil, "problem": problem, "title": "Fixture solution",
		"parts":      map[string]any{"abstract": part, "functional": part, "technical": part},
		"documents":  []any{map[string]any{"id": "solution", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(b), "size_bytes": len(b), "media_type": "text/markdown", "git_provenance": nil}}},
		"supporting": []any{}, "requirements": []any{map[string]any{"id": "f-01", "functional_refs": refs, "technical_refs": refs, "acceptance": "Fixture behavior preserved.", "verification_ids": []string{"test"}}},
		"removed_requirement_ids": []any{}, "phases": []any{map[string]any{"id": "implement", "purpose": "Fixture implementation", "requirement_ids": []string{"f-01"}, "entry_criteria": []any{}, "exit_criteria": []string{"Verified"}, "verification_ids": []string{"test"}}},
		"implementation_basis": map[string]any{"project_id": "ply", "repo_id": "ply", "git_common_dir": r.Tasks[0].GitCommonDir, "epic_id": "epic", "parent_worktree_id": parent.ID, "parent_ref": parent.Ref, "parent_oid": parent.OID, "parent_tree": parent.Tree, "start_oid": parent.OID, "start_tree": parent.Tree}, "dependencies": []any{}, "change_reason": "Fixture solution"}))
	if err != nil {
		t.Fatal(err)
	}
	sr := workspace.TaskRevisionRef{Revision: *spec.OutcomeRef.Revision, ManifestSHA256: spec.OutcomeRef.ManifestSHA256}
	checks := []any{}
	for _, id := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, map[string]any{"id": id, "outcome": "pass", "reason": "Fixture assessment only", "evidence_document_ids": []any{}})
	}
	asm, err := workspace.AssessTaskSpec(w, replacementDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", "replacement/assess", map[string]any{"spec_id": "solution", "spec": sr, "expected_previous_assessment": nil, "outcome": "ready", "reason": "Fixture readiness only", "open_questions": []any{}, "checks": checks, "documents": []any{}}))
	if err != nil {
		t.Fatal(err)
	}
	fields := map[string]any{"expected_previous_selection": nil, "action": "select", "solution": map[string]any{"spec_id": "solution", "spec": sr, "problem": problem, "assessment": workspace.TaskDecisionRef{ID: *asm.OutcomeRef.ID, ManifestSHA256: asm.OutcomeRef.ManifestSHA256}}, "reason": "Fixture choice", "human_decision": map[string]any{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-09-29T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: choose this revision."}}
	choice, err := workspace.SelectTaskSolution(w, replacementDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", "replacement/select", fields))
	if err != nil {
		t.Fatal(err)
	}
	targetInput := workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
	q, err := workspace.ListTaskQueue(w, targetInput, false)
	if err != nil {
		t.Fatal(err)
	}
	queue := workspace.WorkspaceTaskQueueDraft{Kind: "WorkspaceTaskQueueDraft@1", SchemaVersion: 1, PublicationKey: "replacement/queue", ProjectID: "ply", RepoID: "ply", EpicID: "epic", ExpectedRevision: q.Revision, Entries: []workspace.QueueEntry{{TaskID: "replacement", Selection: &workspace.TaskDecisionRef{ID: *choice.OutcomeRef.ID, ManifestSHA256: choice.OutcomeRef.ManifestSHA256}}}, HumanDecision: workspace.QueueHumanDecision{ActorClaim: "synthetic human", DecidedAtUTC: "2026-10-01T00:00:00Z", Source: "explicit_human_instruction", Statement: "Synthetic fixture only."}, RegistryUpgrade: nil}
	qp := writeAny(t, root, "queue.json", queue)
	if _, err = workspace.SetTaskQueue(w, qp); err != nil {
		t.Fatal(err)
	}
}

func al01AbsentReportFixture(t *testing.T) (workflowFixture, WorkflowRun) {
	t.Helper()
	f := workflowTestFixture(t)
	herdr := mustReadAL01(t, f.R.Herdr.Executable.Path)
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+retireOpenPaneStandin+closedTransportInspectStandin+integrationOwnerRecoveryStandin), 1)
	if err := os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "absence-request.json", f.R)
	o := workflowTestStart(t, f)
	workflowTestModel(t, f, map[string]any{"closed_transport": true})
	f.D.CWD = func() (string, error) { return f.R.WorkspaceRoot, nil }
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}}, nil
	}
	return f, o
}

func TestAL01RetirementGuardsAndDurableRetry(t *testing.T) {
	for _, which := range []string{"writer", "unknown_inventory", "dirty", "unknown_effect", "durable_retry"} {
		t.Run(which, func(t *testing.T) {
			f, o := al01AbsentReportFixture(t)
			switch which {
			case "writer":
				f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
					return []LocalProcess{{PID: os.Getpid(), CWD: f.R.WorkspaceRoot}, {PID: 999999, CWD: mustReportWorktree(t, f, o)}}, nil
				}
			case "unknown_inventory":
				workflowTestModel(t, f, map[string]any{"retire_inventory_denied": true})
			case "dirty":
				if err := os.WriteFile(filepath.Join(mustReportWorktree(t, f, o), "unreported"), []byte("change"), 0600); err != nil {
					t.Fatal(err)
				}
			case "unknown_effect":
				evidence := writeAny(t, f.R.WorkspaceRoot, "actual-denial.json", map[string]any{"observation": "synthetic runtime refusal; prior effect unknown"})
				in := DeliveryStopInput{EventID: "synthetic-stop", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "fixture root", Reason: "Unknown prior effect", Action: "report transport", EffectStatus: DeliveryEffectUnknown, EvidencePath: evidence}
				if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
					t.Fatal(err)
				}
			case "durable_retry":
				f.D.Fault = func(stage string) error {
					if stage == "report_after_retirement" {
						return fmt.Errorf("synthetic state publication interruption")
					}
					return nil
				}
			}
			in := DeliveryRetireInput{Observer: "fixture coordinator", Reason: "Selected fresh replacement"}
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err == nil {
				t.Fatal("guard/interruption not reached")
			}
			if which != "durable_retry" {
				s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if s.ReportRetirement != nil {
					t.Fatal("blocked retirement was recorded")
				}
				return
			}
			f.D.Fault = nil
			if _, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
				t.Fatal(err)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || s.ReportRetirement == nil {
				t.Fatalf("durable retirement missing: %v", err)
			}
			if _, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "late.json", f.Acceptance)); err == nil {
				t.Fatal("retired callback revived")
			}
		})
	}
}
func mustReportWorktree(t *testing.T, f workflowFixture, o WorkflowRun) string {
	t.Helper()
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	return s.Observed.Target.WorktreeLocator
}

func TestAL01CoordinatorDenialBlocksReceiptWithoutImpersonation(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := workflowTestStart(t, f)
	evidence := writeAny(t, f.R.WorkspaceRoot, "denied-receipt.json", map[string]any{"source": "synthetic runtime", "action": "receipt callback", "result": "denied before start"})
	in := DeliveryStopInput{EventID: "denied-receipt", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "receipt", Observer: "fixture coordinator", Reason: "Synthetic provider refused callback transport", Action: "receipt callback", EffectStatus: DeliveryEffectNotStarted, EvidencePath: evidence}
	for i := 0; i < 2; i++ {
		if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true); err == nil {
		t.Fatal("receipt erased a preserved runtime denial")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Acceptance != nil || s.StartSHA256 != nil || s.ReportProviderStop == nil || s.Result.Delivery != nil {
		t.Fatal("coordinator manufactured owner acceptance")
	}
}

func TestAL01ReportReceiptCannotChangeNativeTerminalContract(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	native := workflowNativeRequest(f.R)
	raw, _ := Canonical(native)
	if _, err := parseRequest(raw); err == nil {
		t.Fatal("report receipt marker escaped into native terminal/factory contract")
	}
}
