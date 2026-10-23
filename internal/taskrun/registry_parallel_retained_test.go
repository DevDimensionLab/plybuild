package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type activationNoProjectLock struct{ calls *int }

func (l activationNoProjectLock) WithSnapshotLock(string, func(workspace.ProjectSnapshot) error) error {
	*l.calls++
	return fmt.Errorf("forbidden activation check write lock")
}

type activationNoStoreWrites struct {
	workspace.WorkItemStore
	calls *int
}

func (s activationNoStoreWrites) WithLock(string, func(workspace.WorkItemStoreSession) error) error {
	*s.calls++
	return fmt.Errorf("forbidden activation check registry write lock")
}

func activationNativeState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := branchFixtureState(t, root)
	err := filepath.WalkDir(filepath.Join(root, ".ply"), func(path string, d os.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			out[path] = "directory"
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TD7671441: native qualification, integration and Git are real isolated facts;
// all provider/session/process/absence observations are simulated. No live owner
// or installed CLI is used, and the source remains retained throughout.
func retainedActivationFixture(t *testing.T) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires exact private make-built candidate via PLY_DELIVERY_TEST_BINARY")
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("candidate must be absolute")
	}
	raw, err := exec.Command(binary, "version", "--format", "json").Output()
	if err != nil {
		t.Fatal(err)
	}
	var identity buildinfo.Identity
	if err = json.Unmarshal(raw, &identity); err != nil || identity.Dirty || identity.SourceCommit == "" {
		t.Fatal("requires clean embedded source identity", err)
	}
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "simulated native coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, SeedSource: source, SeedOID: identity.SourceCommit})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, f.Prepared.Preparation.Plan.WorktreePath, "merge", "--ff-only", identity.SourceCommit)
	if err = os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); err != nil {
		t.Fatal(err)
	}
	o = parallelQualify(t, f, o, binary)
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if !parallelOwnerHandoff(err) || o.Delivery.Phase != "closing" {
		t.Fatalf("retained automatic delivery: %+v %v", o.Delivery, err)
	}
	return f, o, binary
}

func TestParallelActivationRetainedSourceReachesCompatibility(t *testing.T) {
	f, o, binary := retainedActivationFixture(t)
	d := reviewRoundAsCoordinator(f)
	root := f.R.WorkspaceRoot
	forbidden := 0
	d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		forbidden++
		return nil, fmt.Errorf("forbidden process inventory")
	}
	d.StartupProcessTree = func(int) ([]StartupProcess, error) { forbidden++; return nil, fmt.Errorf("forbidden process tree") }
	d.Workspace.SourceWriters = func(string) (workspace.SourceWriterObservation, error) {
		forbidden++
		return workspace.SourceWriterObservation{}, fmt.Errorf("forbidden source writers")
	}
	d.Fault = func(stage string) error {
		forbidden++
		return fmt.Errorf("forbidden lifecycle/installation stage %s", stage)
	}
	checkOnly := d
	checkOnly.Workspace.ProjectLocks = activationNoProjectLock{&forbidden}
	checkOnly.Workspace.WorkItems = activationNoStoreWrites{d.Workspace.WorkItems, &forbidden}
	before := activationNativeState(t, root)
	if _, err := ActivateParallelWorkspace(d, root, o.RunID, "simulated coordinator", "", false); err == nil || !strings.Contains(err.Error(), "completed") {
		t.Fatalf("default must remain strict: %v", err)
	}
	options := ParallelActivationOptions{AllowRetainedSource: true}
	preview, err := ActivateParallelWorkspaceWithOptions(checkOnly, root, o.RunID, "simulated coordinator", "", false, options)
	if err != nil || preview.State != "preview" || preview.Plan.FromVersion != 4 || len(preview.Plan.Evidence.Consumers) != 1 || preview.Plan.Evidence.Consumers[0].Class != "PAR-A03" {
		t.Fatalf("TD7671441 retained reader must reach finite compatibility: %+v %v", preview, err)
	}
	if !preview.QualifiedReaderEligible || !preview.SourceCleanupPending || !preview.AllowRetainedSource || !equal(before, activationNativeState(t, root)) {
		t.Fatal("activation check lost qualification/cleanup readback or wrote native files/directories")
	}
	// Exercise the exact private image through the public CLI in this fixture.
	cmd := exec.Command(binary, "workspace", "registry", "activate-parallel", "--qualified-run", o.RunID, "--coordinator", "simulated coordinator", "--allow-retained-source", "--check", "--format", "json")
	cmd.Dir = f.Parent
	raw, err := cmd.CombinedOutput()
	var public ParallelActivationReadback
	if err != nil || json.Unmarshal(raw, &public) != nil || !equal(public, preview) {
		t.Fatalf("exact public activation preview: %s %v", raw, err)
	}
	if !equal(before, activationNativeState(t, root)) {
		t.Fatal("public --check wrote native files/directories")
	}
	// Unknown native remaining effects still reach the inventory and veto apply.
	for _, directory := range []string{"week-runs", "bound-goal-starts"} {
		t.Run(directory, func(t *testing.T) {
			path := filepath.Join(root, ".ply", directory, "simulated-unknown")
			if err := os.MkdirAll(path, 0700); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(filepath.Dir(path))
			current := activationNativeState(t, root)
			blocked, err := ActivateParallelWorkspaceWithOptions(checkOnly, root, o.RunID, "simulated coordinator", "", false, options)
			if err != nil || blocked.State != "blocked" || !blocked.QualifiedReaderEligible || len(blocked.Blockers) == 0 {
				t.Fatalf("unknown reader was ignored: %+v %v", blocked, err)
			}
			if _, err := ActivateParallelWorkspaceWithOptions(d, root, o.RunID, "simulated coordinator", blocked.Confirmation, true, options); err == nil {
				t.Fatal("unknown reader allowed a confirmed activation")
			}
			if _, err := ActivateParallelWorkspaceWithOptions(d, root, o.RunID, "simulated coordinator", preview.Confirmation, true, options); err == nil {
				t.Fatal("changed inventory accepted stale confirmation")
			}
			if !equal(current, activationNativeState(t, root)) {
				t.Fatal("blocked reader changed native state")
			}
		})
	}
	if _, err := ActivateParallelWorkspaceWithOptions(d, root, o.RunID, "simulated coordinator", "sha256:"+strings.Repeat("0", 64), true, options); err == nil {
		t.Fatal("wrong confirmation activated format8")
	}
	if !equal(before, activationNativeState(t, root)) {
		t.Fatal("stale confirmation wrote state")
	}
	predecessor, err := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := ActivateParallelWorkspaceWithOptions(d, root, o.RunID, "simulated coordinator", preview.Confirmation, true, options)
	if err != nil || result.State != "activated" || !result.SourceCleanupPending {
		t.Fatalf("confirmed boundary: %+v %v", result, err)
	}
	backup, err := os.ReadFile(result.Backup)
	if err != nil || string(backup) != string(predecessor) {
		t.Fatal("lost exact predecessor backup", err)
	}
	after := activationNativeState(t, root)
	again, err := ActivateParallelWorkspaceWithOptions(d, root, o.RunID, "simulated coordinator", preview.Confirmation, true, options)
	if err != nil || !equal(result, again) || !equal(after, activationNativeState(t, root)) {
		t.Fatalf("retry changed activation: %+v %v", again, err)
	}
	r, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil || r.FormatVersion != 8 || r.ParallelBoundary.FromVersion != 4 || r.ParallelActivation == nil || len(r.IntegrationResults) != 1 {
		t.Fatal("wrong boundary or repeated integration", err)
	}
	for path, bytes := range before {
		if path != filepath.Join(root, ".ply", "work-items.yaml") && after[path] != bytes {
			t.Fatalf("activation changed old control/history %s", path)
		}
	}
	if _, err := os.Stat(f.Prepared.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("activation removed retained source", err)
	}
	status, err := exec.Command(testGit(t), "-C", f.Parent, "status", "--porcelain").Output()
	if err != nil || len(status) != 0 || gitOutput(t, f.Parent, "rev-parse", "HEAD") != o.Delivery.Candidates[0].OID {
		t.Fatal("activation changed target")
	}
	for _, path := range []string{filepath.Join(root, "ply-parallel-reader"), filepath.Join(root, ".ply", "workflow-runs", "cli-installations")} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("activation installed a CLI: %s %v", path, err)
		}
	}
	if forbidden != 0 {
		t.Fatalf("forbidden process/cleanup/installation/write-lock calls: %d", forbidden)
	}
	t.Logf("TD7671441 exact reader=%s: isolated format4 -> 8 once; retained cleanup; all runtime observations simulated", preview.Plan.Evidence.ReaderSHA256)
}

func TestParallelActivationRetainedSourceRejectsMismatchedReaderIdentity(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires private make-built candidate")
	}
	// This fixture deliberately delivers another source commit with the exact
	// binary. Installation evidence alone cannot establish reader eligibility.
	f, o := retainedInstallFixture(t, binary)
	d := reviewRoundAsCoordinator(f)
	before := activationNativeState(t, f.R.WorkspaceRoot)
	result, err := ActivateParallelWorkspaceWithOptions(d, f.R.WorkspaceRoot, o.RunID, "simulated coordinator", "", false, ParallelActivationOptions{AllowRetainedSource: true})
	if err == nil || !strings.Contains(err.Error(), "delivered commit and format8") || result.QualifiedReaderEligible || result.State != "blocked" {
		t.Fatalf("mismatched reader eligible: %+v %v", result, err)
	}
	if !equal(before, activationNativeState(t, f.R.WorkspaceRoot)) {
		t.Fatal("mismatched identity wrote state")
	}
}

func TestParallelActivationRetainedSourceIncompatibleOwnerVeto(t *testing.T) {
	f, o, binary := retainedActivationFixture(t)
	b := parallelAdditionalOwner(t, f, false)
	second := workflowTestStart(t, b.workflowFixture)
	second, err := WorkflowDeliveryReceipt(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: second.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	d := reviewRoundAsCoordinator(f)
	assertVeto := func(class string) {
		t.Helper()
		before := activationNativeState(t, f.R.WorkspaceRoot)
		blocked, err := ActivateParallelWorkspaceWithOptions(d, f.R.WorkspaceRoot, o.RunID, "simulated coordinator", "", false, ParallelActivationOptions{AllowRetainedSource: true})
		if err != nil || !blocked.QualifiedReaderEligible || blocked.State != "blocked" {
			t.Fatalf("%s check: %+v %v", class, blocked, err)
		}
		found := false
		for _, row := range blocked.Plan.Evidence.Consumers {
			if row.ID == "run/"+second.RunID && row.Class == class && !row.Compatible {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing incompatible %s consumer: %+v", class, blocked)
		}
		if _, err = ActivateParallelWorkspaceWithOptions(d, f.R.WorkspaceRoot, o.RunID, "simulated coordinator", blocked.Confirmation, true, ParallelActivationOptions{AllowRetainedSource: true}); err == nil {
			t.Fatalf("%s apply ignored reader veto", class)
		}
		if !equal(before, activationNativeState(t, f.R.WorkspaceRoot)) {
			t.Fatalf("%s check/apply changed old owner/control/history", class)
		}
	}
	assertVeto("PAR-A01")
	if err = os.WriteFile(b.AcceptancePath, []byte("#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); err != nil {
		t.Fatal(err)
	}
	second = parallelQualify(t, b, second, binary)
	if _, err = ReleaseDeliveryOwnership(b.D, b.R.WorkspaceRoot, second.RunID, second.Paths.Context, "Simulated human integration/correction handover; no product QA"); err != nil {
		t.Fatal(err)
	}
	assertVeto("PAR-A04")
}
