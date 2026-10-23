package taskrun

import (
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelActivationExactQualifiedCandidateCheckIsReadOnly(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires the exact private make-built candidate via PLY_DELIVERY_TEST_BINARY")
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("candidate must be absolute; no installed CLI fallback")
	}
	source, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	raw, e := exec.Command(binary, "version", "--format", "json").Output()
	if e != nil {
		t.Fatal(e)
	}
	var identity struct {
		SourceCommit string `json:"source_commit"`
		Dirty        bool   `json:"dirty"`
	}
	if e = json.Unmarshal(raw, &identity); e != nil || identity.Dirty || identity.SourceCommit == "" {
		t.Fatal("candidate lacks exact clean embedded commit")
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "simulated native coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, SeedSource: source, SeedOID: identity.SourceCommit})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o, e = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	runGit(t, f.Prepared.Preparation.Plan.WorktreePath, "merge", "--ff-only", identity.SourceCommit)
	if e = os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if e != nil && !parallelOwnerHandoff(e) {
		t.Fatal(e)
	}
	coordinator := reviewRoundAsCoordinator(f)
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "inventory": []any{}})
	if _, handled, e := ResumeAutomaticLocalDelivery(coordinator, f.R.WorkspaceRoot, o.RunID, o.Delivery.Candidates[0].TaskResult.ID); e != nil || !handled {
		t.Fatalf("complete exact native closeout: %v", e)
	}
	before := branchFixtureState(t, f.R.WorkspaceRoot)
	preview, e := ActivateParallelWorkspace(coordinator, f.R.WorkspaceRoot, o.RunID, "simulated exact qualified coordinator", "", false)
	if e != nil || preview.State != "preview" || preview.Plan.FromVersion != 4 {
		t.Fatalf("qualified activation check: %+v %v", preview, e)
	}
	if !equal(before, branchFixtureState(t, f.R.WorkspaceRoot)) {
		t.Fatal("qualified activation preview wrote native records")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || r.FormatVersion != 4 {
		t.Fatal("qualification alone activated format8")
	}
	t.Logf("exact private image sourceOID=%s readerSHA=%s activation remains preview; simulated runtime only", identity.SourceCommit, preview.Plan.Evidence.ReaderSHA256)
}

// Native legacy Run/control/release/closeout records and Git are real fixture
// facts. Runtime absence and the available qualified coordinator are simulated;
// this test does not install a CLI or claim live provider compatibility.
func TestParallelActivationLegacyCloseoutClasses(t *testing.T) {
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic compatibility coordinator", TimeoutSeconds: 30}}
	f := receiptFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	reviewRoundPatchHerdr(t, &f)
	parallelRuntime(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	o, e := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if e != nil {
		t.Fatal(e)
	}
	originalRequest, e := os.ReadFile(f.File)
	if e != nil {
		t.Fatal(e)
	}
	state, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	class, e := parallelConsumerClass(f.D, state, "sha256:"+strings.Repeat("f", 64))
	if e != nil || class.Class != "PAR-A01" || class.Compatible {
		t.Fatalf("incompatible resumable owner was not blocked: %+v %v", class, e)
	}
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		binary = filepath.Join(f.R.WorkspaceRoot, "synthetic-candidate")
		if e = os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
			t.Fatal(e)
		}
	}
	reviewRoundCandidate(t, f, "fixed legacy source")
	if e = os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ngrep -q fixed README.md\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if e != nil && !parallelOwnerHandoff(e) {
		t.Fatal(e)
	}
	if o.Delivery.Phase != "closing" || o.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("legacy automatic return did not preserve closing source: %+v", o.DeliveryStatus)
	}
	state, e = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	readerSHA := "sha256:" + strings.Repeat("e", 64)
	class, e = parallelConsumerClass(f.D, state, readerSHA)
	if e != nil || class.Class != "PAR-A03" || !class.Compatible {
		t.Fatalf("cleanup-only compatible coordinator was not accepted: %+v %v", class, e)
	}
	missing, e := parallelConsumerClass(f.D, state, "")
	if e != nil || missing.Compatible || !strings.Contains(strings.Join(missing.Reasons, " "), "coordinator") {
		t.Fatalf("missing closeout coordinator not diagnosed: %+v %v", missing, e)
	}
	coordinator := reviewRoundAsCoordinator(f)
	if _, _, e = qualifiedCLIInstallation(coordinator, f.R.WorkspaceRoot, o.RunID, filepath.Join(f.R.WorkspaceRoot, "not-installed-ply"), "synthetic coordinator"); e == nil || !strings.Contains(e.Error(), "completed") {
		t.Fatalf("closing feature bypassed native installer prerequisites: %v", e)
	}
	// Simulated process observation is unknown; it never falls through to ps/lsof.
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return nil, errors.New("simulated writer absence unknown")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "inventory": []any{}})
	c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	_, handled, e := ResumeAutomaticLocalDelivery(coordinator, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if !handled || e == nil {
		t.Fatal("unknown writer absence authorized cleanup")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Plan.WorktreePath); e != nil {
		t.Fatal("unknown absence removed source")
	}
	// This is the workspace effect seam with an explicitly simulated qualified
	// coordinator proof. The public command obtains that proof from native CLI
	// qualification and exact executable identity, never caller JSON.
	observe := func(r workspace.WorkItemRegistry) (workspace.ParallelActivationEvidence, error) {
		s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
		if e != nil {
			return workspace.ParallelActivationEvidence{}, e
		}
		fact, e := parallelConsumerClass(coordinator, s, readerSHA)
		if e != nil {
			return workspace.ParallelActivationEvidence{}, e
		}
		return workspace.ParallelActivationEvidence{Target: f.Prepared.Preparation.Plan.Target, ReaderVersion: "v0.0.0-simulated-qualified", ReaderCommit: c.OID, ReaderSHA256: readerSHA, QualifiedManifestSHA256: "sha256:" + strings.Repeat("d", 64), InventorySHA256: digest(fact), Consumers: []workspace.ParallelConsumerCompatibility{fact}}, nil
	}
	before := branchFixtureState(t, f.R.WorkspaceRoot)
	preview, e := workspace.ActivateParallelRegistry(f.D.Workspace, "", false, observe)
	if e != nil {
		t.Fatal(e)
	}
	after := branchFixtureState(t, f.R.WorkspaceRoot)
	if !equal(before, after) {
		t.Fatal("activation check changed legacy records")
	}
	activated, e := workspace.ActivateParallelRegistry(f.D.Workspace, preview.Confirmation, true, observe)
	if e != nil || activated.State != "activated" {
		t.Fatalf("supported isolated activation: %+v %v", activated, e)
	}
	_, handled, e = ResumeAutomaticLocalDelivery(coordinator, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if !handled || e == nil {
		t.Fatal("format8 erased unknown writer absence")
	}
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: 70001, CWD: f.Parent}, {PID: 70001, ParentPID: 1, CWD: f.Parent}}, nil
	}
	closed, handled, e := ResumeAutomaticLocalDelivery(coordinator, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if e != nil || !handled || closed.Closeout == nil || closed.Closeout.State != "complete" {
		t.Fatalf("new coordinator could not finish same old closeout: %v", e)
	}
	state, e = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	class, e = parallelConsumerClass(coordinator, state, readerSHA)
	if e != nil || class.Class != "PAR-A02" || !class.Compatible {
		t.Fatalf("terminal history blocked activation: %+v %v", class, e)
	}
	raw, e := os.ReadFile(f.File)
	if e != nil || string(raw) != string(originalRequest) {
		t.Fatal("activation/closeout changed old control contract")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || len(r.IntegrationResults) != 1 {
		t.Fatalf("legacy cleanup repeated integration: %v", e)
	}
	if os.Getenv("PLY_DELIVERY_TEST_BINARY") != "" {
		installation, _, e := qualifiedCLIInstallation(coordinator, f.R.WorkspaceRoot, o.RunID, filepath.Join(f.R.WorkspaceRoot, "not-installed-ply"), "synthetic coordinator")
		if e != nil || installation.Source.BinarySHA256 != hashFileTest(t, binary) {
			t.Fatalf("completed native CLI did not qualify after exact closeout: %v", e)
		}
	}
}

func TestParallelActivationHumanReleaseRequiresEveryRemainingControl(t *testing.T) {
	f, o, b, second := parallelOwnerPair(t)
	binary := filepath.Join(f.R.WorkspaceRoot, "simulated-candidate")
	if e := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nexit 0\n"), 0600); e != nil {
		t.Fatal(e)
	}
	o = parallelQualify(t, f, o, binary)
	before, e := workflowRead(b.R.WorkspaceRoot, second.RunID)
	if e != nil {
		t.Fatal(e)
	}
	reviewRoundIn(t, o.Transport.PaneID)
	if _, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Simulated handover to existing human integration/correction route"); e != nil {
		t.Fatal(e)
	}
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	wrong, e := parallelConsumerClass(f.D, s, "sha256:"+strings.Repeat("c", 64))
	if e != nil || wrong.Class != "PAR-A04" || wrong.Compatible || len(wrong.RemainingEffects) != 3 {
		t.Fatalf("human release removed remaining effects: %+v %v", wrong, e)
	}
	runtime, e := workflowEffectiveRuntime(s)
	if e != nil {
		t.Fatal(e)
	}
	compatible, e := parallelConsumerClass(f.D, s, runtime.PlyExecutable.SHA256)
	if e != nil || !compatible.Compatible || compatible.Class != "PAR-A04" {
		t.Fatalf("exact modeled compatible controls refused: %+v %v", compatible, e)
	}
	missing, e := parallelConsumerClass(f.D, s, "")
	if e != nil || missing.Compatible {
		t.Fatal("missing coordinator permitted activation")
	}
	after, e := workflowRead(b.R.WorkspaceRoot, second.RunID)
	if e != nil || !equal(before, after) {
		t.Fatal("A release/classification changed B")
	}
}
