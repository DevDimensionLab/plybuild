package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestIntegrationCLIControllerReplacementAndInstallCrashRecovery(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationCLIFixture(t, binary, "main", true, false)
	installed := filepath.Join(f.R.WorkspaceRoot, "installed-ply")
	bytes, e := os.ReadFile(binary)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(installed, bytes, 0700); e != nil {
		t.Fatal(e)
	}
	artifact := []byte("#!/bin/sh\n# Synthetic installed artifact, not implementation QA.\nexit 0\n")
	staged := filepath.Join(f.R.WorkspaceRoot, "staged-artifact")
	if e = os.WriteFile(staged, artifact, 0700); e != nil {
		t.Fatal(e)
	}
	python, e := exec.LookPath("python3")
	if e != nil {
		t.Fatal(e)
	}
	installer := filepath.Join(f.R.WorkspaceRoot, "install-and-interrupt.py")
	script := `import os, shutil, signal, sys
from pathlib import Path
root=Path(sys.argv[1])
count=root/"install-count"
if count.exists(): raise SystemExit("installation repeated")
shutil.copyfile(root/"staged-artifact", root/"installed-ply")
os.chmod(root/"installed-ply", 0o700)
count.write_text("1")
os.kill(os.getppid(), signal.SIGKILL)
`
	if e = os.WriteFile(installer, []byte(script), 0600); e != nil {
		t.Fatal(e)
	}
	profile := writeAny(t, f.R.WorkspaceRoot, "controller-install.json", map[string]any{
		"kind": "ply.integration.install-profile", "schema_version": 1, "name": "controller-fixture", "candidate_oid": o.Delivery.Candidates[0].OID,
		"artifact_path": installed, "artifact_sha256": hash(artifact),
		"install": map[string]any{"executable": python, "arguments": []string{installer, f.R.WorkspaceRoot}, "cwd": f.Parent},
		"verify":  map[string]any{"executable": installed, "arguments": []string{}, "cwd": f.Parent},
	})
	before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	r, output, e := integrationCLIRun(installed, f.Parent, "pass\n", "--delivery", id, "--install-profile", profile)
	if e == nil || r == nil || deliveryCLIString(t, r, "state") != "controller_required" {
		t.Fatalf("self replacement did not preserve controller: %v %s", e, output)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
		t.Fatal("effect started before stable controller")
	}
	var controller FileBinding
	if e = json.Unmarshal(r["controller"], &controller); e != nil {
		t.Fatal(e)
	}
	operation := deliveryCLIString(t, r, "id")
	if _, _, e = integrationCLIRun(controller.Locator, f.Parent, "", "resume", operation); e == nil {
		t.Fatal("crashed install reported success")
	}
	if b, e := os.ReadFile(installed); e != nil || hash(b) != hash(artifact) {
		t.Fatalf("fixture did not replace ordinary executable: %v", e)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != o.Delivery.Candidates[0].OID {
		t.Fatal("fixture did not integrate before install")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); e != nil {
		t.Fatal("interrupted installation removed source")
	}
	done := integrationCLIOK(t, controller.Locator, f.Parent, "", "resume", operation)
	if deliveryCLIString(t, done, "state") != "completed" {
		t.Fatal("verified installed artifact did not recover")
	}
	// Replay the durable root receipt prefix at its after_closeout interruption
	// boundary: native cleanup has finished but the root has not saved that fact.
	path := filepath.Join(f.R.WorkspaceRoot, ".ply", "integrations", "v1", "records", operation+".json")
	interrupted := workflowProviderDocument(t, path)
	interrupted["closeout"], interrupted["state"] = nil, "closing"
	var prior []any
	for _, event := range interrupted["events"].([]any) {
		if event.(map[string]any)["kind"] != "integration.closed" {
			prior = append(prior, event)
		}
	}
	interrupted["events"] = prior
	writeAny(t, filepath.Dir(path), filepath.Base(path), interrupted)
	integrationCLIOK(t, controller.Locator, f.Parent, "", "resume", operation)
	if b, e := os.ReadFile(filepath.Join(f.R.WorkspaceRoot, "install-count")); e != nil || string(b) != "1" {
		t.Fatal("installation repeated")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); !os.IsNotExist(e) {
		t.Fatal("verified install did not close source")
	}
	registry, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || len(registry.IntegrationAttempts) != 1 {
		t.Fatal("installation recovery repeated integration")
	}
}

func TestIntegrationCLILegacyReconcileWithoutDeliveryOrNewQA(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}, ParentRef: "master"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryCLIQualified(t, f)
	var e error
	o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if e != nil {
		t.Fatal(e)
	}
	o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic legacy writer releases already integrated candidate."); e != nil {
		t.Fatal(e)
	}
	before, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	refBefore := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "HEAD")
	integrationCLIOK(t, binary, f.Parent, "released\n", "release", "--task", "task")
	preview := integrationCLIOK(t, binary, f.Parent, "", "--task", "task", "--reconcile", "--check")
	if deliveryCLIString(t, preview, "state") != "ready" {
		t.Fatalf("legacy blocked: %s", preview["reasons"])
	}
	r := integrationCLIOK(t, binary, f.Parent, "pass\n", "--task", "task", "--reconcile")
	if deliveryCLIString(t, r, "state") != "completed" {
		t.Fatal("legacy not complete")
	}
	integrationCLIOK(t, binary, f.Parent, "", "resume", deliveryCLIString(t, r, "id"))
	after, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if len(after.HumanQARecords) != len(before.HumanQARecords) || len(after.TaskResults) != len(before.TaskResults) || len(after.IntegrationAttempts) != len(before.IntegrationAttempts) {
		t.Fatal("reconcile invented native candidate/QA/integration")
	}
	if refBefore != gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "HEAD") {
		t.Fatal("legacy reconcile merged again")
	}
	list := integrationCLIOK(t, binary, f.Parent, "", "list")
	if string(list["deliveries"]) != "[]" {
		t.Fatal("legacy reconcile invented a Delivery")
	}
}

func TestIntegrationCLIExportHumanJourney(t *testing.T) {
	root := os.Getenv("PLY_INTEGRATION_JOURNEY_ROOT")
	if root == "" {
		t.Skip("export only for an explicitly requested installed disposable human journey")
	}
	if !filepath.IsAbs(root) {
		t.Fatal("journey root must be absolute")
	}
	binary := deliveryCLIBinary(t)
	f, o, id := integrationCLIFixture(t, binary, "main", true, false, root)
	path := writeAny(t, f.R.WorkspaceRoot, "human-journey.json", map[string]any{
		"kind": "ply.integration.human-journey", "schema_version": 1,
		"binary": binary, "workspace": f.R.WorkspaceRoot, "task_worktree": f.Prepared.Preparation.Preparation.Plan.WorktreePath,
		"target_worktree": f.Parent, "target_ref": "refs/heads/main", "delivery_id": id, "candidate_oid": o.Delivery.Candidates[0].OID,
		"fixture_notice":   "Disposable product exercise with a synthetic provider and technical candidate. No human QA has been supplied for this integration. Only the person reviewing the installed implementation can judge that implementation.",
		"preview_argv":     []string{binary, "integration", "--delivery", id, "--check"},
		"interactive_argv": []string{binary, "integration", "--delivery", id},
	})
	t.Log("Prepared installed human journey: " + path)
}

func TestIntegrationCLIInstalledCloseoutCrashRecoversFromRemovedSource(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, id := integrationCLIFixture(t, binary, "main", true, false)
	python, e := exec.LookPath("python3")
	if e != nil {
		t.Fatal(e)
	}
	git, e := exec.LookPath("git")
	if e != nil {
		t.Fatal(e)
	}
	artifact := []byte("#!/bin/sh\nexit 0\n")
	staged := filepath.Join(f.R.WorkspaceRoot, "artifact-source")
	installed := filepath.Join(f.R.WorkspaceRoot, "artifact-installed")
	if e = os.WriteFile(staged, artifact, 0700); e != nil {
		t.Fatal(e)
	}
	install, e := exec.LookPath("install")
	if e != nil {
		t.Fatal(e)
	}
	profile := writeAny(t, f.R.WorkspaceRoot, "install-closeout.json", map[string]any{
		"kind": "ply.integration.install-profile", "schema_version": 1, "name": "cleanup-recovery", "candidate_oid": o.Delivery.Candidates[0].OID,
		"artifact_path": installed, "artifact_sha256": hash(artifact),
		"install": map[string]any{"executable": install, "arguments": []string{"-m", "0700", staged, installed}, "cwd": f.Parent},
		"verify":  map[string]any{"executable": installed, "arguments": []string{}, "cwd": f.Parent},
	})
	model := writeAny(t, f.R.WorkspaceRoot, "cleanup-crash.json", map[string]any{"git": git, "crash": true})
	bin := filepath.Join(f.R.WorkspaceRoot, "cleanup-adapters")
	if e = os.Mkdir(bin, 0700); e != nil {
		t.Fatal(e)
	}
	script := `import json,os,signal,subprocess,sys
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"cleanup-crash.json"
m=json.loads(p.read_text())
a=sys.argv[1:]
if "worktree" in a and "remove" in a:
    result=subprocess.run([m["git"],*a])
    if result.returncode==0:
        m["remove_count"]=m.get("remove_count",0)+1
        crash=m.pop("crash",False)
        p.write_text(json.dumps(m))
        if crash:
            os.kill(os.getppid(),signal.SIGKILL)
            raise SystemExit(9)
    raise SystemExit(result.returncode)
os.execv(m["git"],[m["git"],*a])
`
	if e = os.WriteFile(filepath.Join(bin, "git"), []byte("#!"+python+"\n"+script), 0700); e != nil {
		t.Fatal(e)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	preview := integrationCLIOK(t, binary, f.Parent, "", "--delivery", id, "--install-profile", profile, "--check")
	operation := deliveryCLIString(t, preview, "id")
	if _, _, e = integrationCLIRun(binary, f.Parent, "pass\n", "--delivery", id, "--install-profile", profile); e == nil {
		t.Fatal("cleanup crash reported success")
	}
	if _, e = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); !os.IsNotExist(e) {
		t.Fatal("fixture did not interrupt after source removal")
	}
	done := integrationCLIOK(t, binary, f.Parent, "", "resume", operation)
	if deliveryCLIString(t, done, "state") != "completed" {
		t.Fatal("installed cleanup recovery did not complete")
	}
	if state := workflowProviderDocument(t, model); state["remove_count"] != float64(1) {
		t.Fatal("source removal repeated")
	}
}
