package workflowhandoff

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// This fixture executes local Git and the declared shell test in a disposable
// workspace. All human claims below are explicitly synthetic test inputs.
type deliveryFixture struct {
	d                                      Dependencies
	w                                      workspace.Dependencies
	root, epic, target, parentOID, locator string
	prepared                               workspace.TaskGoalPreparedExecution
}

func deliveryFixtureWrite(t *testing.T, path string, value any) {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, b, 0600); err != nil {
		t.Fatal(err)
	}
}

func prepareDeliveryFixture(t *testing.T) deliveryFixture {
	t.Helper()
	root, epic, _, oid := prepareServiceWorkspace(t)
	runGit(t, epic, "branch", "-m", "epic")
	w := workspace.SystemDependencies()
	if _, e := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Delivery fixture", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: oid}); e != nil {
		t.Fatal(e)
	}
	if _, e := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "task", Title: "Delivery fixture", Description: "Verify the delivery control bridge.", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	r, e := w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	problem := workspace.TaskRevisionRef{Revision: r.TaskProblemRevisions[0].Revision, ManifestSHA256: r.TaskProblemRevisions[0].ManifestSHA256}
	doc, body := filepath.Join(root, "design.md"), []byte("# Goal\nAdd the tested delivery behavior.\n")
	if e = os.WriteFile(doc, body, 0600); e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(root, "goal.json")
	deliveryFixtureWrite(t, file, map[string]any{
		"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "fixture/goal", "task_id": "task", "registry_upgrade": nil,
		"recorder": map[string]any{"actor_claim": "fixture planner", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"},
		"spec_id":  "goal", "expected_previous": nil, "problem": problem, "title": "Delivery fixture", "contract_kind": "goal", "objective": "Add one observable behavior with a real shell acceptance check.",
		"documents": []any{map[string]any{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": digestBytes(body), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}},
		"design":    []any{map[string]any{"document_id": "design", "section": nil}}, "requirements": []any{map[string]any{"id": "behavior", "acceptance": "The README contains the exact delivery behavior."}}, "constraints": []string{"No external effects."},
		"executor": map[string]any{"provider": "codex", "model": "fixture-model", "effort": "high"}, "change_reason": "Isolated test; not actual human approval.",
	})
	spec, e := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "task", File: file})
	if e != nil {
		t.Fatal(e)
	}
	goal := workspace.TaskGoalRef{SpecID: "goal", Spec: workspace.TaskRevisionRef{Revision: *spec.OutcomeRef.Revision, ManifestSHA256: spec.OutcomeRef.ManifestSHA256}}
	r, e = w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	q := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/queue", ProjectID: "ply", RepoID: "ply", EpicID: "epic", Entries: []workspace.TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}}
	if r.FormatVersion < 4 {
		q.RegistryUpgrade = &workspace.TaskRegistryUpgrade{FromVersion: r.FormatVersion, RegistrySHA256: r.RawSHA256}
	}
	file = filepath.Join(root, "queue.json")
	deliveryFixtureWrite(t, file, q)
	if _, e = workspace.SetTaskQueue(w, file); e != nil {
		t.Fatal(e)
	}
	in := workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, Next: true}
	preview, e := workspace.PreviewTaskGoalExecution(w, in)
	if e != nil {
		t.Fatal(e)
	}
	prepared, e := workspace.PrepareTaskGoalExecution(w, in, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "fixture human, not actual approval", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "human_cli", Statement: "Synthetic execute choice only."})
	if e != nil {
		t.Fatal(e)
	}
	d := SystemDependencies()
	raw, e := BuildDeliveryHandoffDraft(d, prepared.Preparation.Preparation.ID, "fixture human, not actual approval", "fixture-owner", "fixture/delivery", prepared.AcceptancePath)
	if e != nil {
		t.Fatal(e)
	}
	file = filepath.Join(root, "owner-draft.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	created, e := Create(d, CreateInput{DraftPath: file})
	if e != nil {
		t.Fatal(e)
	}
	temp := filepath.Join(root, "recipient-temp")
	if e = os.MkdirAll(temp, 0700); e != nil {
		t.Fatal(e)
	}
	sandbox, _ := json.Marshal(map[string]any{"read_roots": []string{root}, "write_roots": []string{root}, "temp_root": temp, "matches_contract": true})
	start, e := BuildDeliveryTaskRunStart(d, created.Locator, "fixture human, not actual approval", "synthetic test", "fixture-owner-session", "codex", "fixture-model", sandbox, []byte("[]"), "started")
	if e != nil {
		t.Fatal(e)
	}
	file = filepath.Join(root, "owner-start.json")
	if e = os.WriteFile(file, start, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = SubmitStart(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: file}); e != nil {
		t.Fatal(e)
	}
	return deliveryFixture{d: d, w: w, root: root, epic: epic, target: prepared.Preparation.Plan.WorktreePath, parentOID: oid, locator: created.Locator, prepared: prepared}
}

func (f deliveryFixture) candidate(t *testing.T, key string) DeliveryCandidateInput {
	t.Helper()
	if e := os.WriteFile(filepath.Join(f.target, "README.md"), []byte("delivered\n"), 0644); e != nil {
		t.Fatal(e)
	}
	runGit(t, f.target, "add", "README.md")
	runGit(t, f.target, "commit", "--allow-empty", "-m", "delivery candidate "+key)
	in := DeliveryCandidateInput{ParentHandoffLocator: f.locator, CandidateKey: key, Summary: "Observed fixture acceptance and preserved technical review.", CandidateOID: gitOutput(t, f.target, "rev-parse", "HEAD"), CandidateTree: gitOutput(t, f.target, "rev-parse", "HEAD^{tree}"), VerifierID: "goal-acceptance", CWD: f.target, Argv: []string{"/bin/sh", f.prepared.AcceptancePath}}
	in.RunID, in.RequestSHA256 = "fixture-run", digestBytes([]byte("fixture request"))
	area := filepath.Join(f.root, "evidence-"+key)
	if e := os.MkdirAll(area, 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.MkdirAll(filepath.Dir(f.prepared.AcceptancePath), 0700); e != nil {
		t.Fatal(e)
	}
	script := []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = delivered\nprintf 'behavior passed\\n'\n")
	if e := os.WriteFile(f.prepared.AcceptancePath, script, 0600); e != nil {
		t.Fatal(e)
	}
	snapshot := filepath.Join(area, "acceptance.sh")
	if e := os.WriteFile(snapshot, script, 0600); e != nil {
		t.Fatal(e)
	}
	in.ReviewPath = filepath.Join(area, "review.json")
	template, e := DeliveryCandidateReviewTemplate(in.CandidateOID, in.CandidateTree)
	if e != nil {
		t.Fatal(e)
	}
	var review map[string]any
	if e = json.Unmarshal(template, &review); e != nil {
		t.Fatal(e)
	}
	review["reviewer_claim"], review["reviewer_session_id"], review["decision"] = "synthetic fixture review", "fixture-review-session", "passed"
	deliveryFixtureWrite(t, in.ReviewPath, review)
	reviewRaw, e := os.ReadFile(in.ReviewPath)
	if e != nil {
		t.Fatal(e)
	}
	var stdout, stderr bytes.Buffer
	command := exec.Command(in.Argv[0], in.Argv[1:]...)
	command.Dir, command.Stdout, command.Stderr = in.CWD, &stdout, &stderr
	started := time.Now().UTC().Format(time.RFC3339Nano)
	if e = command.Run(); e != nil {
		t.Fatal(e)
	}
	finished := time.Now().UTC().Format(time.RFC3339Nano)
	in.StdoutPath, in.StderrPath = filepath.Join(area, "stdout.txt"), filepath.Join(area, "stderr.txt")
	if e = os.WriteFile(in.StdoutPath, stdout.Bytes(), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(in.StderrPath, stderr.Bytes(), 0600); e != nil {
		t.Fatal(e)
	}
	in.VerificationPath = filepath.Join(area, "verification.json")
	binding := func(p string, b []byte) any { return map[string]any{"locator": p, "sha256": digestBytes(b)} }
	deliveryFixtureWrite(t, in.VerificationPath, map[string]any{"kind": "PlyDeliveryVerification@1", "schema_version": 1, "run_id": "fixture-run", "request_sha256": digestBytes([]byte("fixture request")), "attempt_id": key, "candidate_oid": in.CandidateOID, "candidate_tree": in.CandidateTree, "argv": in.Argv, "cwd": in.CWD, "acceptance": binding(f.prepared.AcceptancePath, script), "acceptance_snapshot": binding(snapshot, script), "review": binding(in.ReviewPath, reviewRaw), "exit": 0, "stdout": binding(in.StdoutPath, stdout.Bytes()), "stderr": binding(in.StderrPath, stderr.Bytes()), "started_at": started, "finished_at": finished})
	return in
}

func (f deliveryFixture) human(t *testing.T, result workspace.TaskResultRecord, outcome string) (workspace.TaskHumanQARecord, error) {
	t.Helper()
	path := filepath.Join(f.root, "human-"+outcome+"-"+string(result.ID)+".json")
	deliveryFixtureWrite(t, path, DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: "task", TaskResultID: string(result.ID), ResultOID: result.ResultOID, ResultTree: result.ResultTree, Outcome: outcome, ActorClaim: "synthetic human fixture; not actual product approval", StartSurface: "isolated test", StartedAtUTC: "2026-10-05T00:00:00Z", CompletedAtUTC: "2026-10-05T00:01:00Z", Answer: outcome, Observation: "Synthetic fixture answer exercises evidence binding only."})
	return RecordDeliveryHumanQA(f.d, "task", result, outcome, path)
}

func TestDeliveryCandidateHumanQAAndSingleLocalIntegration(t *testing.T) {
	f := prepareDeliveryFixture(t)
	input := f.candidate(t, "attempt-1")
	qualified, e := QualifyDeliveryCandidate(f.d, input)
	if e != nil {
		t.Fatal(e)
	}
	if qualified.TaskResult.ResultOID != input.CandidateOID || qualified.TaskResult.Recorder.ActorClaim != "fixture-owner" {
		t.Fatal("candidate or recorder binding lost")
	}
	control, e := f.d.Store.ReadByLocator(qualified.Handoff.Locator)
	if e != nil {
		t.Fatal(e)
	}
	effects, _ := objectMember(control.Terminal.Value, "observed_effects")
	if len(effects.([]canonicaljson.Value)) != 0 {
		t.Fatal("candidate control falsely reported the owner's prior command as a new effect")
	}
	repeated, e := QualifyDeliveryCandidate(f.d, input)
	if e != nil || repeated.TaskResult.ID != qualified.TaskResult.ID {
		t.Fatalf("same candidate did not recover: %v", e)
	}
	owner, e := f.d.Store.ReadByLocator(f.locator)
	if e != nil || owner.Terminal != nil {
		t.Fatal("technical candidate terminated delivery owner")
	}
	qa, e := f.human(t, qualified.TaskResult, "pass")
	if e != nil {
		t.Fatal(e)
	}
	authority := DeliveryOwnerAuthority{RunID: "fixture-run", RequestSHA256: digestBytes([]byte("fixture request")), ActorClaim: "fixture-owner", PreparationID: f.prepared.Preparation.Preparation.ID}
	returned, e := IntegrateDeliveryCandidate(f.d, "task", qualified.TaskResult.ID, qa.ID, f.parentOID, input.CandidateOID, authority)
	if e != nil {
		t.Fatal(e)
	}
	if !returned.Completed || gitOutput(t, f.epic, "rev-parse", "HEAD") != input.CandidateOID || returned.Queue.Current != nil {
		t.Fatal("delivery did not complete exact local return")
	}
	r, e := f.w.WorkItems.Snapshot(f.root)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.IntegrationAttempts) != 1 || len(r.IntegrationAuthorities) != 1 || r.IntegrationAuthorities[0].Mode != "delivery_owner_after_human_pass" || r.IntegrationAuthorities[0].DeliveryOwner == nil {
		t.Fatal("missing exact delivery-owner integration provenance")
	}
	returned, e = IntegrateDeliveryCandidate(f.d, "task", qualified.TaskResult.ID, qa.ID, f.parentOID, input.CandidateOID, authority)
	if e != nil || !returned.Completed {
		t.Fatalf("completed return did not recover: %v", e)
	}
	r, e = f.w.WorkItems.Snapshot(f.root)
	if e != nil || len(r.IntegrationAttempts) != 1 {
		t.Fatal("return replay attempted another merge")
	}
}

func TestDeliveryCandidateRejectsDriftAndNonGreenEvidence(t *testing.T) {
	f := prepareDeliveryFixture(t)
	in := f.candidate(t, "attempt-1")
	bad := in
	bad.Exit = 1
	if _, e := QualifyDeliveryCandidate(f.d, bad); e == nil {
		t.Fatal("failed verifier qualified candidate")
	}
	bad = in
	bad.CWD = f.epic
	if _, e := QualifyDeliveryCandidate(f.d, bad); e == nil {
		t.Fatal("wrong command directory qualified candidate")
	}
	bad = in
	bad.RunID = "another-run"
	if _, e := QualifyDeliveryCandidate(f.d, bad); e == nil {
		t.Fatal("verification receipt from another run qualified candidate")
	}
	if e := os.WriteFile(filepath.Join(f.target, "untracked"), []byte("changed"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := QualifyDeliveryCandidate(f.d, in); e == nil {
		t.Fatal("dirty candidate qualified")
	}
	r, e := f.w.WorkItems.Snapshot(f.root)
	if e != nil || len(r.TaskResults) != 0 {
		t.Fatal("rejected observations published a TaskResult")
	}
}

func TestDeliveryReviewRequiresExplicitCandidateBoundAssessment(t *testing.T) {
	in := DeliveryCandidateInput{CandidateOID: strings.Repeat("a", 40), CandidateTree: strings.Repeat("b", 40)}
	raw, e := DeliveryCandidateReviewTemplate(in.CandidateOID, in.CandidateTree)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e == nil {
		t.Fatal("unobserved template is not a review")
	}
	var review map[string]any
	if e = json.Unmarshal(raw, &review); e != nil {
		t.Fatal(e)
	}
	review["reviewer_claim"], review["decision"] = "fixture reviewer", "passed"
	raw, _ = json.Marshal(review)
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e != nil {
		t.Fatal(e)
	}
	in.RequireIndependentReview = true
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e == nil {
		t.Fatal("missing independent reviewer accepted")
	}
	review["reviewer_session_id"] = "owner-session"
	raw, _ = json.Marshal(review)
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e == nil {
		t.Fatal("owner accepted as independent reviewer")
	}
	review["reviewer_session_id"] = "review-session"
	raw, _ = json.Marshal(review)
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e != nil {
		t.Fatal(e)
	}
	review["candidate_oid"] = strings.Repeat("c", 40)
	raw, _ = json.Marshal(review)
	if _, e = deliveryCandidateReview(raw, in, "owner-session"); e == nil {
		t.Fatal("review of another candidate accepted")
	}
}

func TestDeliveryBudgetAndInstallScopeAreExplicit(t *testing.T) {
	f := prepareDeliveryFixture(t)
	s, e := f.d.Store.ReadByLocator(f.locator)
	if e != nil {
		t.Fatal(e)
	}
	if !unlimitedDeliveryBudget(s.Handoff.Value) {
		t.Fatal("unbounded owner acquired an invented round limit")
	}
	a, _ := objectMember(s.Handoff.Value, "authority")
	allowed, _ := objectMember(a.(canonicaljson.Object), "allowed_effects")
	for _, v := range allowed.([]canonicaljson.Value) {
		if objectString(v.(canonicaljson.Object), "type") == "install" {
			t.Fatal("installation enabled by default")
		}
	}
	legacy := taskSpecLegacyProjection(s.Handoff.Value)
	if e := validateStoredHandoffV1(legacy); e == nil {
		t.Fatal("legacy schema accepted new delivery budget")
	}
	if _, _, e = validateHandoffBudget(deliveryBudget(0), true); e == nil {
		t.Fatal("zero was silently treated as unlimited")
	}
	if _, _, e = validateHandoffBudget(deliveryBudget(nil), false); e == nil {
		t.Fatal("legacy accepted unlimited budget")
	}
	raw, e := BuildDeliveryHandoffDraft(f.d, f.prepared.Preparation.Preparation.ID, "fixture human, not actual approval", "fixture-owner", "fixture/install-scope", f.prepared.AcceptancePath, DeliveryDraftAuthority{AllowLocalInstall: true})
	if e != nil {
		t.Fatal(e)
	}
	draft, e := decodeHandoffDraft(raw)
	if e != nil {
		t.Fatal(e)
	}
	effects, _ := objectMember(draft.Authority, "allowed_effects")
	foundInstall := false
	for _, v := range effects.([]canonicaljson.Value) {
		if objectString(v.(canonicaljson.Object), "type") == "install" {
			foundInstall = true
		}
	}
	if !foundInstall {
		t.Fatal("explicit installation choice was not included")
	}
}

func TestDeliveryHumanFailureKeepsOwnerAndNewCandidateSeparate(t *testing.T) {
	f := prepareDeliveryFixture(t)
	firstInput := f.candidate(t, "attempt-1")
	first, e := QualifyDeliveryCandidate(f.d, firstInput)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.human(t, first.TaskResult, ""); e == nil {
		t.Fatal("missing human answer became QA")
	}
	failed, e := f.human(t, first.TaskResult, "fail")
	if e != nil {
		t.Fatal(e)
	}
	authority := DeliveryOwnerAuthority{RunID: firstInput.RunID, RequestSHA256: firstInput.RequestSHA256, ActorClaim: "fixture-owner", PreparationID: f.prepared.Preparation.Preparation.ID}
	if _, e = IntegrateDeliveryCandidate(f.d, "task", first.TaskResult.ID, failed.ID, f.parentOID, firstInput.CandidateOID, authority); e == nil {
		t.Fatal("human failure allowed integration")
	}
	secondInput := f.candidate(t, "attempt-2")
	second, e := QualifyDeliveryCandidate(f.d, secondInput)
	if e != nil {
		t.Fatal(e)
	}
	if second.Handoff.HandoffID == first.Handoff.HandoffID || second.TaskResult.ID == first.TaskResult.ID || second.TaskResult.ResultOID == first.TaskResult.ResultOID {
		t.Fatal("new candidate rewrote the prior immutable result")
	}
	if _, e = IntegrateDeliveryCandidate(f.d, "task", second.TaskResult.ID, failed.ID, f.parentOID, secondInput.CandidateOID, authority); e == nil {
		t.Fatal("old QA answer moved to a new candidate")
	}
	owner, e := f.d.Store.ReadByLocator(f.locator)
	if e != nil || owner.Terminal != nil {
		t.Fatal("human failure terminated the responsible owner")
	}
	r, e := f.w.WorkItems.Snapshot(f.root)
	if e != nil || len(r.TaskResults) != 2 || len(r.HumanQARecords) != 1 || len(r.IntegrationAttempts) != 0 {
		t.Fatal("candidate generations or failed QA were not preserved")
	}
}

func TestDeliveryChangedParentStopsBeforeMerge(t *testing.T) {
	f := prepareDeliveryFixture(t)
	in := f.candidate(t, "attempt-1")
	qualified, e := QualifyDeliveryCandidate(f.d, in)
	if e != nil {
		t.Fatal(e)
	}
	qa, e := f.human(t, qualified.TaskResult, "pass")
	if e != nil {
		t.Fatal(e)
	}
	runGit(t, f.epic, "commit", "--allow-empty", "-m", "concurrent parent movement")
	changed := gitOutput(t, f.epic, "rev-parse", "HEAD")
	authority := DeliveryOwnerAuthority{RunID: in.RunID, RequestSHA256: in.RequestSHA256, ActorClaim: "fixture-owner", PreparationID: f.prepared.Preparation.Preparation.ID}
	if _, e = IntegrateDeliveryCandidate(f.d, "task", qualified.TaskResult.ID, qa.ID, f.parentOID, in.CandidateOID, authority); e == nil {
		t.Fatal("changed parent accepted without refresh")
	}
	if gitOutput(t, f.epic, "rev-parse", "HEAD") != changed {
		t.Fatal("changed parent was modified")
	}
	r, e := f.w.WorkItems.Snapshot(f.root)
	if e != nil || len(r.IntegrationAttempts) != 0 {
		t.Fatal("changed parent started a merge attempt")
	}
}

func TestDeliveryStagingSymlinkIsRejectedBeforeDirectoryCreation(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	outside := filepath.Join(root, "outside")
	if e = os.Mkdir(outside, 0700); e != nil {
		t.Fatal(e)
	}
	link := filepath.Join(root, "staging")
	if e = os.Symlink(outside, link); e != nil {
		t.Fatal(e)
	}
	if e = deliveryEnsureDir(SystemDependencies().Files, filepath.Join(link, "candidate")); e == nil {
		t.Fatal("staging symlink accepted")
	}
	if _, e = os.Stat(filepath.Join(outside, "candidate")); !os.IsNotExist(e) {
		t.Fatal("directory created through rejected staging symlink")
	}
}
