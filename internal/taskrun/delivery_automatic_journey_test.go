package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// This exporter leaves an installed, disposable automatic Delivery ready for
// observation and execution. Its provider, runtime permission claim and review
// are synthetic fixture data; no HumanQA record or implementation approval is
// created. The pinned public CLI performs actual acceptance and later delivery.
func TestAutomaticDeliveryCLIExportJourney(t *testing.T) {
	root := os.Getenv("PLY_AUTOMATIC_JOURNEY_ROOT")
	if root == "" {
		t.Skip("export only when explicitly requested for an installed automatic journey")
	}
	if !filepath.IsAbs(root) || filepath.Clean(root) != root || root == string(filepath.Separator) {
		t.Fatal("PLY_AUTOMATIC_JOURNEY_ROOT must be a clean absolute disposable directory")
	}
	if entries, err := os.ReadDir(root); err == nil && len(entries) != 0 {
		t.Fatal("automatic journey root must be empty; preserve any existing journey")
	} else if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if pinned := os.Getenv("PLY_DELIVERY_TEST_BINARY"); pinned == "" || !filepath.IsAbs(pinned) {
		t.Fatal("export requires PLY_DELIVERY_TEST_BINARY pinned to an absolute installed CLI")
	}
	binary := automaticCLIBinary(t)
	if info, err := os.Stat(binary); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("pinned CLI must be a regular file: %v", err)
	}
	a := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic journey orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: a, Root: root})
	// Bind the installed CLI before startup, so the immutable request and actual
	// public acceptance callback agree on the executable being exercised.
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := workflowTestStart(t, f.workflowFixture)
	runtimeAcceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", runtimeAcceptance)
	parentBefore := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	deliveryTestCandidate(t, f, "automatic local Epic journey")
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'automatic local Epic journey'\n\"$PLY_CANDIDATE_BINARY\" capabilities --format json >/dev/null\nprintf 'observed automatic fixture candidate\\n'\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "automatic-journey")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
	if o.Delivery == nil || o.Delivery.Phase != "automatic_acceptance_passed" || len(o.Delivery.Candidates) != 1 || o.Delivery.Candidates[0].TaskResult.ID == "" {
		t.Fatal("public CLI did not preserve one automatically accepted candidate")
	}
	candidate := o.Delivery.Candidates[0]
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registration := deliveryCLIRegisterInput(t, f, o, "automatic-journey")
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", registration)
	id := deliveryCLIString(t, registered, "id")
	ready := deliveryCLIOK(t, binary, cwd, "check", id)
	if deliveryCLIString(t, ready, "state") != "ready" {
		t.Fatalf("automatic journey is not ready: %s", ready["reasons"])
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.HumanQARecords) != 0 || len(r.IntegrationAuthorities) != 0 || len(r.IntegrationAttempts) != 0 || len(r.IntegrationResults) != 0 || gitOutput(t, f.Parent, "rev-parse", "HEAD") != parentBefore {
		t.Fatalf("journey export created human QA or advanced integration: %v", err)
	}
	show := [][]string{
		{binary, "workflow", "execute", "show", o.RunID, "--format", "json"},
		{binary, "workflow", "delivery", "show", id, "--format", "json"},
		{binary, "workspace", "task", "show", "task", "--format", "json"},
	}
	deliver := []string{binary, "workflow", "delivery", "execute", id, "--format", "json"}
	showPath := writeAutomaticJourneyScript(t, f.R.WorkspaceRoot, "show.sh", cwd, show)
	deliverPath := writeAutomaticJourneyScript(t, f.R.WorkspaceRoot, "deliver.sh", cwd, [][]string{deliver})
	path := writeAny(t, f.R.WorkspaceRoot, "automatic-journey.json", map[string]any{
		"kind": "ply.delivery.automatic-journey", "schema_version": 1,
		"binary": binary, "binary_sha256": f.R.Runtime.PlyExecutable.SHA256,
		"workspace": f.R.WorkspaceRoot, "task_worktree": cwd, "task_id": "task",
		"target_worktree": f.Parent, "target_ref": "refs/heads/epic", "target_oid_at_export": parentBefore,
		"workflow_run_id": o.RunID, "context": o.Paths.Context, "delivery_id": id,
		"candidate_oid": candidate.OID, "candidate_tree": candidate.Tree, "task_result_id": candidate.TaskResult.ID,
		"agreement": f.R.Delivery.Agreement, "registration_file": registration,
		"runtime_acceptance_file": runtimeAcceptance, "review_file": review, "acceptance_path": f.AcceptancePath,
		"exported_state": "ready", "human_qa_records_at_export": 0, "integration_performed_at_export": false,
		"show_script": showPath, "deliver_script": deliverPath, "show_argv": show, "deliver_argv": deliver,
		"fixture_provenance": map[string]any{
			"provider": "synthetic fixture provider", "runtime_permission": "synthetic fixture claim, not an OS or active agent session attestation",
			"review":      "synthetic fixture reviewer; not independent review of this implementation",
			"acceptance":  "actual pinned public CLI executed acceptance.sh against the disposable candidate",
			"integration": "not executed by exporter; deliver.sh uses the pinned public CLI for real local Git integration and native queue, base and lifecycle completion",
		},
		"fixture_notice": "Disposable product exercise. Provider, runtime and review records are synthetic. No human QA or approval of the implementation is recorded. Repeating deliver.sh resumes the same Delivery and must not repeat an observed Git effect. Automatic closeout retains the Task worktree and branch.",
	})
	t.Log("Prepared installed automatic journey: " + path)
}

func writeAutomaticJourneyScript(t *testing.T, root, name, cwd string, commands [][]string) string {
	t.Helper()
	var script strings.Builder
	script.WriteString("#!/bin/sh\nset -eu\n# Disposable fixture; no human QA or implementation approval.\ncd " + ShellQuote(cwd) + "\n")
	for i, argv := range commands {
		if i == len(commands)-1 {
			script.WriteString("exec ")
		}
		quoted := make([]string, len(argv))
		for j, arg := range argv {
			quoted[j] = ShellQuote(arg)
		}
		script.WriteString(strings.Join(quoted, " ") + "\n")
	}
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(script.String()), 0700); err != nil {
		t.Fatal(err)
	}
	return path
}
