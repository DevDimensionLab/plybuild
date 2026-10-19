package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Real frozen Goal/Spec/Handoff and public CLI verification/integration/retention.
// Provider, runtime claim and review are synthetic; this is no Claude acceptance.
func informativeCLIFixture(t *testing.T, binary string) (deliveryFixture, WorkflowRun, []byte) {
	t.Helper()
	var report []byte
	root := ""
	if parent := os.Getenv("PLY_INFORMATIVE_JOURNEY_ROOT"); parent != "" {
		var err error
		root, err = os.MkdirTemp(parent, "consumer-")
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Preserved synthetic consumer workspace:", root)
	}
	agreement := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic informative consumer", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: agreement, Root: root, InformativeReport: func(root string) []byte {
		origin := filepath.Join(root, "historical-acceptance.sh")
		first := []byte("historical first acceptance bytes\n")
		if err := os.WriteFile(origin, first, 0600); err != nil {
			t.Fatal(err)
		}
		reportBytes, err := json.MarshalIndent(map[string]any{"kind": "PlyDeliveryVerification@2", "events": []any{
			map[string]any{"acceptance": FileBinding{origin, hash(first)}},
			map[string]any{"acceptance": FileBinding{origin, hash([]byte("historical second acceptance bytes\n"))}},
		}}, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		report = append(reportBytes, '\n')
		return report
	}})
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	acceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance)
	deliveryTestCandidate(t, f, "informative historical consumer")
	review := deliveryTestReview(t, f, "informative")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
	return f, o, report
}

func assertInformativeRetained(t *testing.T, f deliveryFixture, o WorkflowRun, report []byte) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	resource := map[bool]string{true: "kept", false: "retired"}[closeout != nil && closeout.Plan.Keep]
	if err != nil || closeout == nil || closeout.State != "complete" || !closeout.LifecycleCompleted || closeout.ResourceState != resource || closeout.WorktreeRemoved == closeout.Plan.Keep || closeout.BranchRemoved == closeout.Plan.Keep || closeout.Retention == nil {
		t.Fatalf("native closeout with retained evidence incomplete: %+v %v", closeout, err)
	}
	found := false
	for _, e := range closeout.Retention.Entries {
		if e.OriginLocator == filepath.Join(f.R.WorkspaceRoot, "historical-acceptance.sh") {
			t.Fatal("historical citation became a new required dependency")
		}
		if e.SHA256 != hash(report) {
			continue
		}
		found = true
		if e.SizeBytes != int64(len(report)) {
			t.Fatal("report size changed")
		}
		got, err := workspace.ReadTaskRetainedEvidence(f.R.WorkspaceRoot, c.TaskResult.TaskID, e.OriginLocator, e.SHA256)
		if err != nil || !bytes.Equal(got, report) {
			t.Fatal("exact original binding did not resolve retained report", err)
		}
		original, err := os.ReadFile(e.OriginLocator)
		if err != nil || !bytes.Equal(original, report) {
			t.Fatal("frozen report rewritten", err)
		}
	}
	if !found {
		t.Fatal("report missing from retained evidence")
	}
	if _, err := os.Lstat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); closeout.Plan.Keep && err != nil || !closeout.Plan.Keep && !os.IsNotExist(err) {
		t.Fatal("source state differs from the frozen cleanup grant", err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("wrong integrated target")
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatal("Epic base did not advance to exact candidate", err)
	}
	t.Logf("Native closeout: run=%s result=%s candidate=%s manifest=%s; report bytes/hash/size retained, source %s", o.RunID, c.TaskResult.ID, c.OID, closeout.Retention.ManifestSHA256, closeout.ResourceState)
}

func TestInformativeHistoricalReportCLIRetainedCloseout(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, report := informativeCLIFixture(t, binary)
	handoff, output, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if err == nil || !bytes.Contains(output, []byte("owner_handoff:")) || handoff.Delivery == nil || handoff.Delivery.Phase != "closing" {
		t.Fatalf("owner integrate did not end in the pending-cleanup handoff: %v\n%s", err, output)
	}
	cleanupOwnerLeaves(t, f)
	if o, output, err = cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30"); err != nil || o.Delivery.Phase != "completed" {
		t.Fatalf("native closeout did not finish: %v\n%s", err, output)
	}
	assertInformativeRetained(t, f, o, report)
}

// Set PLY_INFORMATIVE_OLD_READER to the unmodified base CLI. It creates the
// actual already-integrated retention failure, then the candidate uses native
// same-Run continuation. No private state is edited to construct this boundary.
func TestInformativeHistoricalReportCLIResumesOldReaderCloseout(t *testing.T) {
	old := os.Getenv("PLY_INFORMATIVE_OLD_READER")
	if old == "" {
		t.Skip("requires a preserved pre-fix CLI")
	}
	binary := automaticCLIBinary(t)
	f, o, report := informativeCLIFixture(t, old)
	failed, output, err := automaticCLI(old, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if err == nil || !strings.Contains(string(output), "conflicting evidence references") || failed.Delivery == nil || failed.Delivery.Phase != "closing" || failed.Delivery.OwnershipRelease == nil || failed.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("old reader did not reproduce integrated/released retention failure: %v %s", err, output)
	}
	c := failed.Delivery.Candidates[0]
	t.Logf("Old reader reached actual retention conflict after integration/release: run=%s result=%s candidate=%s", failed.RunID, c.TaskResult.ID, c.OID)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("old reader had not actually integrated")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	// All immutable inputs and existing delivery event bytes must survive resume.
	bindings := []FileBinding{f.R.Delivery.Goal, {o.Paths.Context, hashFileTest(t, o.Paths.Context)}, {o.Handoff.Locator, o.Handoff.SHA256}, *failed.Delivery.OwnershipRelease, *c.Integration}
	for _, e := range failed.Delivery.Events {
		bindings = append(bindings, e.Binding)
	}
	done := automaticCLIOK(t, binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if done.RunID != o.RunID || done.Delivery.Phase != "completed" || len(done.Delivery.Candidates) != 1 || done.Delivery.Candidates[0].TaskResult.ID != c.TaskResult.ID || !equal(done.Delivery.OwnershipRelease, failed.Delivery.OwnershipRelease) {
		t.Fatal("resume changed result/owner or omitted completion")
	}
	assertInformativeRetained(t, f, done, report)
	if actual := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); actual != reflog {
		t.Fatal("closeout resume repeated Git integration")
	}
	after, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(registry.IntegrationResults, after.IntegrationResults) || !equal(registry.IntegrationAttempts, after.IntegrationAttempts) || len(after.HumanQARecords) != 0 {
		t.Fatal("resume changed integration or manufactured human QA")
	}
	for _, b := range bindings {
		if hashFileTest(t, b.Locator) != b.SHA256 {
			t.Fatalf("resume rewrote immutable input %s", b.Locator)
		}
	}
	automaticCLIOK(t, binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("completed continuation repeated Git")
	}
	t.Log("Same-Run/result closeout resumed; prior Goal/context/handoff/events/release unchanged; integration attempts/results and target reflog unchanged")
}
