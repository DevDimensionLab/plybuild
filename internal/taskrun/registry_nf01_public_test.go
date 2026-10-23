package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func nf01Image(t *testing.T) (string, string, string) {
	t.Helper()
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires exact clean make candidate via PLY_DELIVERY_TEST_BINARY")
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("candidate must be absolute")
	}
	b, e := exec.Command(binary, "version", "--format", "json").Output()
	if e != nil {
		t.Fatal(e)
	}
	var id buildinfo.Identity
	if e = json.Unmarshal(b, &id); e != nil || id.Dirty || id.SourceCommit == "" {
		t.Fatal("clean embedded identity required", e)
	}
	source, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	return binary, source, id.SourceCommit
}

// Supply only a new goal draft; all qualification, review, release and native
// delivery facts below are created through native fixture paths.
func nf01GoalTemplate(t *testing.T, f workflowFixture) {
	t.Helper()
	root := f.R.WorkspaceRoot
	doc := filepath.Join(root, "nf01-reader-goal.md")
	b := []byte("# Disposable exact reader\nVerify the exact candidate without live provider observations.\n")
	if e := os.WriteFile(doc, b, 0600); e != nil {
		t.Fatal(e)
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, ProjectID: "ply", RepoID: "ply", EpicID: "epic", Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "modeled NF01 coordinator", TimeoutSeconds: 30}}
	writeAny(t, root, "goal-draft.json", map[string]any{"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "nf01/template", "task_id": "task", "spec_id": "solution", "registry_upgrade": nil, "expected_previous": nil, "problem": nil, "recorder": map[string]any{"actor_claim": "fixture planner", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"}, "contract_kind": "goal", "title": "Exact fixture reader", "objective": "Qualify an exact private candidate within a disposable native fixture.", "change_reason": "NF01 public read-only compatibility demo", "documents": []any{map[string]any{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(b), "size_bytes": len(b), "media_type": "text/markdown", "git_provenance": nil}}}, "design": []any{map[string]any{"document_id": "design", "section": nil}}, "constraints": []string{"Only disposable native fixture."}, "requirements": []any{map[string]any{"id": "behavior", "acceptance": "Exact private reader remains read-only."}}, "executor": map[string]any{"provider": "codex", "model": "fixture-model", "effort": "medium"}, "delivery": a})
}
func nf01PublicReader(t *testing.T, f workflowFixture, binary, oid string) (deliveryFixture, WorkflowRun) {
	t.Helper()
	s, e := workflowRead(f.R.WorkspaceRoot, workflowID(f.R))
	if e != nil {
		t.Fatal(e)
	}
	nf01GoalTemplate(t, f)
	base := deliveryFixture{workflowFixture: f, Parent: s.Observed.Epic.WorktreeLocator}
	b := parallelAdditionalOwner(t, base, false)
	nf01Runtime(t, &b)
	reviewRoundPatchHerdr(t, &b)
	// The additional owner's stand-in uses its own pane; the native reviewer
	// fixture patch is pane-agnostic and retains that actual modeled identity.
	o := workflowTestStart(t, b.workflowFixture)
	o, e = WorkflowDeliveryReceipt(b.D, b.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	runGit(t, b.Prepared.Preparation.Plan.WorktreePath, "merge", "--ff-only", oid)
	if e = os.WriteFile(b.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest -f go.mod\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, b, o, binary)
	o, e = WorkflowDeliveryIntegrate(b.D, b.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if !parallelOwnerHandoff(e) || o.Delivery.Phase != "closing" {
		t.Fatalf("native reader retained return: %+v %v", o.Delivery, e)
	}
	return b, o
}
func nf01PublicCheck(t *testing.T, f deliveryFixture, qualified WorkflowRun, binary, want string, classes map[string]string) {
	t.Helper()
	d := f.D
	d.CWD = func() (string, error) { return f.Parent, nil }
	forbidden := 0
	d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		forbidden++
		t.Error("unexpected process observation")
		return nil, nil
	}
	d.StartupProcessTree = func(int) ([]StartupProcess, error) { forbidden++; t.Error("unexpected process tree"); return nil, nil }
	d.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		forbidden++
		t.Error("unexpected writer observation")
		return workspace.SourceWriterObservation{}, nil
	}
	d.Workflow.TaskWorkspace = &d.Workspace
	d.Workspace.ProjectLocks = activationNoProjectLock{&forbidden}
	d.Workspace.WorkItems = activationNoStoreWrites{d.Workspace.WorkItems, &forbidden}
	before := activationNativeState(t, f.R.WorkspaceRoot)
	direct, e := ActivateParallelWorkspaceWithOptions(d, f.R.WorkspaceRoot, qualified.RunID, "synthetic NF01 coordinator", "", false, ParallelActivationOptions{AllowRetainedSource: true})
	if e != nil || direct.State != want {
		t.Fatalf("isolated native preview: %+v %v", direct, e)
	}
	var previous []byte
	for i := 0; i < 2; i++ {
		cmd := exec.Command(binary, "workspace", "registry", "activate-parallel", "--qualified-run", qualified.RunID, "--coordinator", "synthetic NF01 coordinator", "--allow-retained-source", "--check", "--format", "json")
		cmd.Dir = f.Parent
		raw, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("isolated public --check: %v %s", e, raw)
		}
		var out ParallelActivationReadback
		if e = json.Unmarshal(raw, &out); e != nil || !equal(out, direct) || out.State != want || out.Plan.FromVersion != 4 {
			t.Fatalf("public state/classes: %v %s", e, raw)
		}
		for id, class := range classes {
			found := false
			for _, c := range out.Plan.Evidence.Consumers {
				if c.ID == id {
					found = true
					if c.Class != class || class == "PAR-A02" && (!c.Compatible || len(c.RemainingEffects) != 0) {
						t.Fatalf("wrong public consumer: %+v", c)
					}
				}
			}
			if !found {
				t.Fatalf("missing public consumer %s: %s", id, raw)
			}
		}
		for _, c := range out.Plan.Evidence.Consumers {
			if strings.HasPrefix(c.ID, "preparation/") {
				t.Fatalf("false duplicate preparation blocker: %+v", c)
			}
		}
		if i > 0 && string(previous) != string(raw) {
			t.Fatal("nondeterministic public JSON")
		}
		previous = raw
	}
	if forbidden != 0 || !equal(before, activationNativeState(t, f.R.WorkspaceRoot)) {
		t.Fatal("check wrote native state or invoked a forbidden observation/lock")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || r.FormatVersion != 4 {
		t.Fatal("check changed registry format", e)
	}
	t.Logf("T0F15276E exact public --check %s; modeled OS/provider only; registryformat4 unchanged", want)
}
func TestNF01PublicCompatibilityCheck(t *testing.T) {
	binary, source, oid := nf01Image(t)
	t.Run("retired-report", func(t *testing.T) {
		f, o := nf01Report(t, false, source, oid)
		reader, q := nf01PublicReader(t, f, binary, oid)
		nf01PublicCheck(t, reader, q, binary, "preview", map[string]string{"run/" + o.RunID: "PAR-A02", "run/" + q.RunID: "PAR-A03"})
	})
	t.Run("human-terminal", func(t *testing.T) {
		f, o := nf01Human(t, true, source, oid)
		reader, q := nf01PublicReader(t, f.workflowFixture, binary, oid)
		nf01PublicCheck(t, reader, q, binary, "preview", map[string]string{"run/" + o.RunID: "PAR-A02", "run/" + q.RunID: "PAR-A03"})
	})
	t.Run("headless-and-unknown", func(t *testing.T) {
		f, o := newWorkerDeliveryFixtureOptions(t, deliveryFixtureOptions{SeedSource: source, SeedOID: oid}, func(f *deliveryFixture) {
			nf01Runtime(t, f)
			f.R.Runtime.PlyExecutable = Executable{binary, hashFileTest(t, binary)}
		})
		workerTestReserveNative(t, f, o.RunID, "verify")
		var e error
		o, e = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, binary)
		if e != nil {
			t.Fatal(e)
		}
		workerTestReserveNative(t, f, o.RunID, "integrate")
		o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
		if e != nil {
			t.Fatal(e)
		}
		nf01PublicCheck(t, f, o, binary, "preview", map[string]string{"run/" + o.RunID: "PAR-A02"})
		// A separate captured headless preparation remains active and incompatible.
		base := f
		base.R.Herdr.Executable = Executable{filepath.Join(f.R.WorkspaceRoot, "herdr-standin"), hashFileTest(t, filepath.Join(f.R.WorkspaceRoot, "herdr-standin"))}
		b := parallelAdditionalOwner(t, base, false)
		nf01Runtime(t, &b)
		b.R.Envelope = deliveryEnv("headless-run-request")
		b.R.Worker = &WorkerContract{SchemaVersion: 1, MaxAttempts: 3, AgentTimeoutSeconds: 60}
		b.R.Herdr = (WorkflowRequest{}).Herdr
		b.R.Runtime.Mode = "headless"
		b.R.Runtime.PermissionBinding.ProfileID = "ply-worker"
		b.R.HumanAuthority.StartSurface = "human_authorized_worker"
		b.R.Delivery.StartupProtocol = ""
		b.R.HandoffDraft, e = workflowhandoff.BuildDeliveryHandoffDraft(b.D.Workflow, b.R.PreparationID, "synthetic human", b.R.Delivery.OwnerClaim, "nf01/active-headless", b.AcceptancePath, workflowhandoff.DeliveryDraftAuthority{Agreement: b.Prepared.Delivery})
		if e != nil {
			t.Fatal(e)
		}
		b.File = writeAny(t, b.R.WorkspaceRoot, "nf01-active-headless.json", b.R)
		active := workflowTestStart(t, b.workflowFixture)
		if e = os.MkdirAll(filepath.Join(f.R.WorkspaceRoot, ".ply", "week-runs", "modeled-week"), 0700); e != nil {
			t.Fatal(e)
		}
		nf01PublicCheck(t, f, o, binary, "blocked", map[string]string{"run/" + o.RunID: "PAR-A02", "run/" + active.RunID: "unknown", "week-runs/modeled-week": "unknown"})
	})
}
