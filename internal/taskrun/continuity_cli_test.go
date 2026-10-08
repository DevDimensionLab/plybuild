package taskrun

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// This is the last source revision before delivery agreements were added to
// Task Goal/Spec@2. Building its real CLI, rather than changing a current decoder
// in a test, proves the preserved old-reader contract at the command boundary.
const continuityOldSource = "f2de0c638f1d2b1c38417aff70b79ca7bb4be31a"

func continuityOldBinary(t *testing.T) string {
	t.Helper()
	return continuitySourceBinary(t, continuityOldSource, "PLY_CONTINUITY_OLD_BINARY")
}

func continuitySourceBinary(t *testing.T, revision, override string) string {
	t.Helper()
	if path := os.Getenv(override); path != "" {
		physical, err := filepath.EvalSymlinks(path)
		if err != nil || !filepath.IsAbs(physical) {
			t.Fatalf("invalid historical CLI path: %s: %v", path, err)
		}
		return physical
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	archive := exec.Command("git", "archive", "--format=tar", revision)
	archive.Dir = repo
	raw, err := archive.Output()
	if err != nil {
		t.Fatalf("historical source %s must be available locally: %v", revision, err)
	}
	root := t.TempDir()
	r := tar.NewReader(bytes.NewReader(raw))
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		name := filepath.Clean(header.Name)
		if filepath.IsAbs(name) || name == ".." || strings.HasPrefix(name, ".."+string(os.PathSeparator)) {
			t.Fatalf("unsafe historical archive entry %q", header.Name)
		}
		path := filepath.Join(root, name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0700); err != nil {
				t.Fatal(err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, os.FileMode(header.Mode)&0777); err != nil {
				t.Fatal(err)
			}
		}
	}
	binary := filepath.Join(root, "ply-old")
	build := exec.Command("go", "build", "-o", binary, "./cmd/ply")
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real historical CLI: %v\n%s", err, output)
	}
	return binary
}

func continuityCLI(binary, cwd string, args ...string) ([]byte, error) {
	c := exec.Command(binary, append(args, "--format", "json")...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	if err != nil {
		return stdout.Bytes(), fmt.Errorf("%v: %w: %s", args, err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func continuityCLIOK(t *testing.T, binary, cwd string, args ...string) []byte {
	t.Helper()
	raw, err := continuityCLI(binary, cwd, args...)
	if err != nil {
		t.Fatalf("native command failed: %v\n%s", err, raw)
	}
	return raw
}

func continuityRun(t *testing.T, raw []byte) WorkflowRun {
	t.Helper()
	var run WorkflowRun
	if err := json.Unmarshal(raw, &run); err != nil || run.RunID == "" {
		t.Fatalf("invalid run readback: %v\n%s", err, raw)
	}
	return run
}

type continuityFixture struct {
	deliveryFixture
	Old, Current, CWD, Review, Counter, FirstOID, CorrectionOID string
	Run                                                         WorkflowRun
	Unrelated                                                   map[string]string
	Upgrade                                                     *continuityProviderUpgrade
}

func continuityCounter(t *testing.T, path string) int {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(raw), "acceptance executed\n")
}

func continuityCandidate(t *testing.T, f deliveryFixture, text, counter string) string {
	t.Helper()
	oid := deliveryTestCandidate(t, f, text)
	// A side-effect counter outside the worktree proves receipt reuse does not
	// execute the successful acceptance command again.
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = " + ShellQuote(text) + "\nprintf 'acceptance executed\\n' >> " + ShellQuote(counter) + "\nprintf 'observed fixture behavior\\n'\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	return oid
}

func continuityPublish(t *testing.T, f deliveryFixture, current, task string) map[string]string {
	t.Helper()
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	continuityCLIOK(t, current, cwd, "workspace", "task", "create", task, "--title", "Unrelated newer publication", "--description", "Valid delivery metadata published while a legacy Task is active.", "--epic", "epic", "--project", "ply", "--repo", "ply")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	var problem workspace.TaskRevisionRef
	for _, revision := range registry.TaskProblemRevisions {
		if string(revision.TaskID) == task {
			problem = workspace.TaskRevisionRef{Revision: revision.Revision, ManifestSHA256: revision.ManifestSHA256}
		}
	}
	if problem.Revision == 0 {
		t.Fatal("unrelated Task did not publish its actual problem")
	}
	draft := workflowProviderDocument(t, filepath.Join(f.R.WorkspaceRoot, "goal-draft.json"))
	draft["task_id"], draft["spec_id"], draft["publication_key"], draft["problem"] = task, task+"-solution", "continuity/"+task, problem
	draft["delivery"] = workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic, ProjectID: "ply", RepoID: "ply", EpicID: "epic", TargetRef: "refs/heads/epic", TargetWorktree: f.Parent}
	path := writeAny(t, f.R.WorkspaceRoot, task+"-goal.json", draft)
	continuityCLIOK(t, current, cwd, "workspace", "task", "spec", "record", task, "--file", path)
	registry, err = f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	bindings := map[string]string{}
	for _, publication := range registry.TaskContentPublications {
		if string(publication.TaskID) != task {
			continue
		}
		for kind, digest := range map[string]string{"manifests": publication.OutcomeRef.ManifestSHA256, "requests": publication.RequestSHA256} {
			file := filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", kind, "sha256", strings.TrimPrefix(digest, "sha256:")+".json")
			bindings[file] = hashFileTest(t, file)
		}
	}
	if len(bindings) < 4 {
		t.Fatal("unrelated publication bytes were not preserved")
	}
	return bindings
}

func newContinuityFixture(t *testing.T, current, old, root string) continuityFixture {
	t.Helper()
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Root: root})
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	control := filepath.Join(f.R.WorkspaceRoot, "historical-ply-control")
	raw, err := os.ReadFile(old)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(control, raw, 0500); err != nil {
		t.Fatal(err)
	}
	f.R.Runtime.PlyExecutable = Executable{control, hashFileTest(t, control)}
	f.D.Executable = func() (string, error) { return control, nil }
	f.File = writeAny(t, f.R.WorkspaceRoot, "legacy-request.json", f.R)
	var preview WorkflowPreview
	if err = json.Unmarshal(continuityCLIOK(t, control, cwd, "workflow", "run", "start", "--file", f.File, "--check"), &preview); err != nil || preview.Confirmation == nil {
		t.Fatalf("old-reader start did not preview: %+v %v", preview, err)
	}
	o := continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "run", "start", "--file", f.File, "--apply", "--confirm", *preview.Confirmation))
	claim := writeAny(t, f.R.WorkspaceRoot, "legacy-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", claim))
	counter := filepath.Join(f.R.WorkspaceRoot, "acceptance-executions.txt")
	first := continuityCandidate(t, f, "first behavior", counter)
	review := deliveryTestReview(t, f, "first")
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review))
	if len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].OID != first {
		t.Fatalf("old reader never qualified its earlier candidate: %+v", o.Delivery)
	}
	o = continuityRun(t, continuityCLIOK(t, control, cwd, "workflow", "execute", "qa", o.RunID, "fail", "--context", o.Paths.Context, "--evidence", deliveryTestHuman(t, f, o, "fail")))
	unrelated := continuityPublish(t, f, current, "newer")
	correction := continuityCandidate(t, f, "corrected behavior", counter)
	review = deliveryTestReview(t, f, "correction")
	failed, err := continuityCLI(control, cwd, "workflow", "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review)
	if err == nil || !strings.Contains(err.Error(), "task_content_invalid_input: object has missing or unknown fields") {
		t.Fatalf("historical reader did not reproduce the real qualification failure: %v\n%s", err, failed)
	}
	o = continuityRun(t, failed)
	if o.Delivery.Attempt == nil || o.Delivery.Attempt.State != "recorded" || o.Delivery.Attempt.CandidateOID != correction || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("recorded failure lost the correction or invented qualification: %+v", o.Delivery)
	}
	var receipt deliveryVerificationReceipt
	if err = readValue(filepath.Join(o.Delivery.Attempt.Path, "verification.json"), 1<<20, &receipt); err != nil || receipt.Exit == nil || *receipt.Exit != 0 || receipt.CandidateOID != correction || receipt.Review.SHA256 != hashFileTest(t, review) {
		t.Fatalf("real successful command evidence missing: %+v %v", receipt, err)
	}
	if continuityCounter(t, counter) != 2 {
		t.Fatal("native old reader did not execute each acceptance exactly once")
	}
	if _, err = continuityCLI(control, cwd, "workspace", "task", "show", "task"); err == nil || !strings.Contains(err.Error(), "task_content_invalid_input: object has missing or unknown fields") {
		t.Fatalf("old Task readback did not reproduce the compatibility boundary: %v", err)
	}
	t.Logf("Historical reader %s: first candidate %s qualified; correction %s acceptance exit 0, qualification failed; actual attempt %s", continuityOldSource, first, correction, o.Delivery.Attempt.ID)
	return continuityFixture{deliveryFixture: f, Old: control, Current: current, CWD: cwd, Review: review, Counter: counter, FirstOID: first, CorrectionOID: correction, Run: o, Unrelated: unrelated}
}

func continuityContinue(t *testing.T, f continuityFixture) string {
	t.Helper()
	args := []string{"workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context}
	if f.Upgrade != nil {
		args = append(args, "--runtime-evidence", f.Upgrade.RuntimeObservation.Locator)
	}
	raw := continuityCLIOK(t, f.Current, f.CWD, args...)
	return continuityControl(t, raw)
}

func continuityControl(t *testing.T, raw []byte) string {
	t.Helper()
	var result struct {
		State   string `json:"state"`
		Preview struct {
			Control Executable `json:"control_executable"`
		} `json:"preview"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || result.Preview.Control.Path == "" || result.State != "continued" && result.State != "existing" {
		t.Fatalf("continuation did not preserve a compatible control: %v\n%s", err, raw)
	}
	return result.Preview.Control.Path
}

func TestContinuityNativeOldReaderCorrection(t *testing.T) {
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	f := newContinuityFixture(t, current, old, "")
	continuityUpgradeProvider(t, &f, before, "removed")
	protected := map[string]string{f.Old: hashFileTest(t, f.Old), f.Run.Paths.Context: hashFileTest(t, f.Run.Paths.Context), f.Run.Handoff.Locator: f.Run.Handoff.SHA256, filepath.Join(f.Run.Delivery.Attempt.Path, "verification.json"): hashFileTest(t, filepath.Join(f.Run.Delivery.Attempt.Path, "verification.json"))}
	shown := continuityRun(t, continuityCLIOK(t, current, f.CWD, "workflow", "execute", "show", f.Run.RunID))
	status := shown.DeliveryStatus
	if status == nil || status.Source.OID != f.CorrectionOID || status.Verification == nil || status.Verification.AttemptID != f.Run.Delivery.Attempt.ID || status.Verification.Outcome != "passed" || status.Verification.Qualification != "not_qualified" || status.QualifiedCandidate == nil || status.QualifiedCandidate.OID != f.FirstOID || status.QualifiedCandidate.Current || status.HumanJudgment == nil || status.HumanJudgment.Outcome != "fail" || status.HumanJudgment.Current || status.FinalDelivery.State != "not_delivered" {
		t.Fatalf("native status mixed successful correction verification with the older candidate: %+v", status)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "different-live-session"})
	blocked, err := continuityCLI(current, f.CWD, "workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context, "--runtime-evidence", f.Upgrade.RuntimeObservation.Locator, "--check")
	var diagnosis struct {
		State      string `json:"state"`
		Diagnostic struct {
			Code, Operation, Cause string
			ReaderCapability       string `json:"reader_capability"`
			NextAction             string `json:"next_action"`
		} `json:"diagnostic"`
	}
	if decodeErr := json.Unmarshal(blocked, &diagnosis); err == nil || decodeErr != nil || diagnosis.State != "blocked" || diagnosis.Diagnostic.Code != "delivery_continuation_session" || diagnosis.Diagnostic.Operation == "" || diagnosis.Diagnostic.ReaderCapability == "" || diagnosis.Diagnostic.Cause == "" || diagnosis.Diagnostic.NextAction == "" {
		t.Fatalf("stale live session lacks a precise machine-readable diagnosis: %v\n%s", err, blocked)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "fixture-session"})
	control := continuityContinue(t, f)
	o := continuityRun(t, continuityCLIOK(t, control, f.CWD, "workflow", "execute", "verify", f.Run.RunID, "--context", f.Run.Paths.Context, "--review", f.Review, "--reuse", f.Run.Delivery.Attempt.ID))
	if len(o.Delivery.Candidates) != 2 || o.Delivery.Candidates[1].OID != f.CorrectionOID || o.Delivery.Candidates[0].OID != f.FirstOID || o.Delivery.Candidates[0].HumanQA == nil || o.Delivery.Candidates[0].HumanQA.Outcome != "fail" || o.Delivery.Candidates[1].HumanQA != nil || o.Delivery.Phase != "awaiting_human_qa" {
		t.Fatalf("continuation confused corrected qualification with human approval: %+v", o.Delivery)
	}
	if o.DeliveryStatus == nil || o.DeliveryStatus.Verification.Qualification != "qualified" || !o.DeliveryStatus.QualifiedCandidate.Current || o.DeliveryStatus.HumanJudgment != nil || o.DeliveryStatus.FinalDelivery.State != "not_delivered" {
		t.Fatalf("native status confused requalification with human judgment or delivery: %+v", o.DeliveryStatus)
	}
	continuityCLIOK(t, control, f.CWD, "workflow", "execute", "verify", f.Run.RunID, "--context", f.Run.Paths.Context, "--review", f.Review, "--reuse", f.Run.Delivery.Attempt.ID)
	if continuityCounter(t, f.Counter) != 2 {
		t.Fatal("receipt reuse reran successful acceptance")
	}
	if _, err := continuityCLI(control, f.CWD, "workflow", "execute", "integrate", f.Run.RunID, "--context", f.Run.Paths.Context); err == nil {
		t.Fatal("correction integrated without its own actual QA record")
	}
	for path, digest := range f.Unrelated {
		protected[path] = digest
	}
	for path, digest := range protected {
		if hashFileTest(t, path) != digest {
			t.Fatalf("continuation rewrote preserved evidence: %s", path)
		}
	}
	if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
		t.Fatal("continuation restarted or replaced the delivery owner")
	}
	continuityCLIOK(t, current, f.CWD, "workspace", "task", "show", "task")
	continuityCLIOK(t, current, f.CWD, "workspace", "task", "show", "newer")
	// Progress remains a callback from the same accepted session after the
	// technical candidate. It cannot substitute for that candidate's human QA.
	report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "continued-progress", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "The correction is ready for an actual product check.", Meaning: "Synthetic fixture only; this progress report is not a human pass.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	reportPath := writeAny(t, f.R.WorkspaceRoot, "continued-report.json", report)
	continuityCLIOK(t, control, f.CWD, "workflow", "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", reportPath)
	continuityCLIOK(t, control, f.CWD, "workflow", "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", reportPath)
	o = continuityRun(t, continuityCLIOK(t, control, f.CWD, "workflow", "execute", "qa", o.RunID, "pass", "--context", o.Paths.Context, "--evidence", deliveryTestHuman(t, f.deliveryFixture, o, "pass")))
	o = continuityRun(t, continuityCLIOK(t, control, f.CWD, "workflow", "execute", "integrate", o.RunID, "--context", o.Paths.Context))
	if o.Delivery.Phase != "completed" || gitOutput(t, f.Parent, "rev-parse", "HEAD") != f.CorrectionOID {
		t.Fatalf("authorized legacy local return did not complete: %+v", o.Delivery)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 2 || len(registry.HumanQARecords) != 2 {
		t.Fatalf("continuation duplicated or lost qualification/QA records: %v", err)
	}
}

func TestContinuityNativeConcurrentContinuationAndPublication(t *testing.T) {
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	f := newContinuityFixture(t, current, old, "")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	draft := workflowProviderDocument(t, filepath.Join(f.R.WorkspaceRoot, "newer-goal.json"))
	for _, revision := range registry.TaskSpecRevisions {
		if revision.TaskID == "newer" {
			draft["expected_previous"] = workspace.TaskRevisionRef{Revision: revision.Revision, ManifestSHA256: revision.ManifestSHA256}
		}
	}
	draft["publication_key"], draft["change_reason"] = "continuity/newer-revision", "Concurrent unrelated publication retains both writers."
	publication := writeAny(t, f.R.WorkspaceRoot, "concurrent-publication.json", draft)
	var wg sync.WaitGroup
	var outputs [3][]byte
	var failures [3]error
	for i := range outputs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := []string{"workflow", "execute", "continue", f.Run.RunID, "--context", f.Run.Paths.Context}
			if i == 2 {
				args = []string{"workspace", "task", "spec", "record", "newer", "--file", publication}
			}
			outputs[i], failures[i] = continuityCLI(current, f.CWD, args...)
		}(i)
	}
	wg.Wait()
	for i, err := range failures {
		if err != nil {
			t.Fatalf("concurrent native writer %d failed: %v\n%s", i, err, outputs[i])
		}
	}
	control := continuityControl(t, outputs[0])
	if continuityControl(t, outputs[1]) != control {
		t.Fatal("concurrent continuation acquired more than one control")
	}
	continuityCLIOK(t, control, f.CWD, "workflow", "execute", "verify", f.Run.RunID, "--context", f.Run.Paths.Context, "--review", f.Review, "--reuse", f.Run.Delivery.Attempt.ID)
	registry, err = f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	revisions := 0
	for _, revision := range registry.TaskSpecRevisions {
		if revision.TaskID == "newer" {
			revisions++
		}
	}
	if revisions != 2 || len(registry.TaskResults) != 2 || continuityCounter(t, f.Counter) != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
		t.Fatal("concurrent continuation lost publication or duplicated an effect")
	}
	for path, digest := range f.Unrelated {
		if hashFileTest(t, path) != digest {
			t.Fatalf("concurrent publication rewrote history: %s", path)
		}
	}
}

func continuityJourneyManifest(t *testing.T, f continuityFixture) map[string]any {
	t.Helper()
	return map[string]any{
		"kind": "PlyContinuityJourney@1", "schema_version": 1,
		"binary": f.Current, "workspace": f.R.WorkspaceRoot,
		"task_worktree": f.CWD, "target_worktree": f.Parent,
		"old_control": f.Old, "old_control_sha256": hashFileTest(t, f.Old),
		"old_source_revision": continuityOldSource, "run_id": f.Run.RunID,
		"context": f.Run.Paths.Context, "context_sha256": hashFileTest(t, f.Run.Paths.Context),
		"first_candidate": f.FirstOID, "corrected_candidate": f.CorrectionOID,
		"attempt": f.Run.Delivery.Attempt.ID, "review": f.Review,
		"counter": f.Counter, "provider_calls": f.Calls,
		"native_session_id":      f.Run.Transport.AgentSessionID,
		"provider_upgrade":       f.Upgrade,
		"unrelated_publications": f.Unrelated,
		"fixture_notice":         "Disposable product exercise. The provider, authority and earlier QA are synthetic fixtures. This journey does not approve the implementation candidate in the real delivery.",
	}
}

func TestContinuityJourneyPreservesExplicitOutcomeAndExactFeedback(t *testing.T) {
	driver, err := filepath.Abs("../../test/continuity_journey.py")
	if err != nil {
		t.Fatal(err)
	}
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	f := newContinuityFixture(t, current, old, "")
	continuityUpgradeProvider(t, &f, before, "removed")
	manifest := continuityJourneyManifest(t, f)
	manifest["qa_actor_claim"] = "Synthetic automated journey regression, not actual human QA"
	path := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", manifest)
	for _, action := range []string{"inspect", "report-incomplete", "report-incomplete", "continue", "continue", "verify", "verify", "inspect"} {
		command := exec.Command("python3", driver, "--fixture", path, action)
		command.Dir = f.CWD
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("upgraded installed journey %s failed: %v\n%s", action, err, output)
		}
	}
	qualified := continuityRun(t, continuityCLIOK(t, current, f.CWD, "workflow", "execute", "show", f.Run.RunID))
	if len(qualified.Delivery.Candidates) != 2 || qualified.Delivery.Candidates[1].HumanQA != nil || continuityCounter(t, f.Counter) != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
		t.Fatalf("short journey commands duplicated effects or supplied human QA: %+v", qualified.Delivery)
	}
	feedback := "Synthetic regression only: the preserved correction is ready. Ikke ekte menneskelig QA."
	for i := 0; i < 2; i++ {
		command := exec.Command("python3", driver, "--fixture", path, "qa", "pass", "--answer", feedback)
		command.Dir = f.CWD
		var stderr bytes.Buffer
		command.Stderr = &stderr
		if _, err := command.Output(); err != nil {
			t.Fatalf("native journey must preserve explicit pass plus exact feedback: %v\n%s", err, stderr.String())
		}
	}
	o := continuityRun(t, continuityCLIOK(t, current, f.CWD, "workflow", "execute", "show", f.Run.RunID))
	qa := o.Delivery.Candidates[len(o.Delivery.Candidates)-1].HumanQA
	if qa == nil || qa.Outcome != "pass" || !strings.Contains(qa.Observation, feedback) || qa.Actor.ActorClaim != manifest["qa_actor_claim"] {
		t.Fatalf("native QA lost the explicit synthetic outcome, actor or exact feedback: %+v", qa)
	}
	var attestation map[string]any
	raw, err := os.ReadFile(qa.Evidence[0].Locator)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &attestation); err != nil || attestation["answer"] != "pass" || !strings.Contains(fmt.Sprint(attestation["observation"]), feedback) {
		t.Fatalf("native attestation did not preserve outcome and feedback: %v\n%s", err, raw)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.HumanQARecords) != 2 {
		t.Fatalf("identical journey QA repeated a human record: %v", err)
	}
}

func TestContinuityCLIExportHumanJourney(t *testing.T) {
	root := os.Getenv("PLY_CONTINUITY_JOURNEY_ROOT")
	if root == "" {
		t.Skip("export only for an explicitly requested installed disposable journey")
	}
	if !filepath.IsAbs(root) {
		t.Fatal("journey root must be absolute")
	}
	current, old := deliveryCLIBinary(t), continuityOldBinary(t)
	before := continuitySourceBinary(t, continuityBeforeUpgradeSource, "PLY_CONTINUITY_BEFORE_UPGRADE_BINARY")
	f := newContinuityFixture(t, current, old, root)
	continuityUpgradeProvider(t, &f, before, "removed")
	path := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", continuityJourneyManifest(t, f))
	t.Log("Prepared installed continuity journey: " + path)
}
