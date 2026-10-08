package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestReadInventoryAbsentStoresStayAbsent(t *testing.T) {
	d, r, _ := fixture(t)
	out, err := ReadInventory(d, r.WorkspaceRoot)
	if err != nil || len(out.Runs) != 0 || len(out.Reasons) != 0 {
		t.Fatalf("empty inventory: %+v, %v", out, err)
	}
	for _, p := range []string{storeRoot(r.WorkspaceRoot), workflowRoot(r.WorkspaceRoot)} {
		if _, err := os.Lstat(p); !os.IsNotExist(err) {
			t.Fatalf("read created run storage: %s (%v)", p, err)
		}
	}
}

func TestReadInventoryPreservesProviderAndDoesNotClaimCachedLiveness(t *testing.T) {
	d, r, _ := fixture(t)
	for _, provider := range []string{"codex", "claude"} {
		w := WorkflowRequest{Envelope: workflowEnv("herdr-run-request"), RequestKey: "inventory/" + provider, WorkspaceRoot: r.WorkspaceRoot, PreparationID: r.PreparationID, PreparationSHA256: r.PreparationSHA256, HandoffDraft: r.HandoffDraft, Runtime: r.Runtime}
		w.Runtime.Provider = provider
		x := workspace.PlanWorktreeObservation{StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
		s := workflowInitial(w, Observed{Target: x, Epic: x, RuntimeBindings: []FileBinding{}})
		s.Result.Transport = WorkflowTransport{WorkspaceID: "w1", TabID: "t7", PaneID: "p7", AgentSessionID: provider + "-native", State: "working", ObservedAt: "2026-10-05T10:00:00Z"}
		for _, path := range []string{filepath.Dir(workflowIndex(r.WorkspaceRoot, s.Result.RunID)), s.Result.Paths.RunRoot} {
			if err := os.MkdirAll(path, 0700); err != nil {
				t.Fatal(err)
			}
		}
		if err := writeValue(workflowIndex(r.WorkspaceRoot, s.Result.RunID), s); err != nil {
			t.Fatal(err)
		}
		if err := workflowSave(d, s); err != nil {
			t.Fatal(err)
		}
	}
	// A historical reader must not execute or hash current provider binaries.
	if err := os.Remove(r.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	out, err := ReadInventory(d, r.WorkspaceRoot)
	if err != nil || len(out.Runs) != 2 {
		t.Fatalf("inventory: %+v, %v", out, err)
	}
	providers := map[string]bool{}
	for _, row := range out.Runs {
		if row.Provider == nil {
			t.Fatalf("missing provider: %+v", row)
		}
		providers[*row.Provider] = true
		if row.Transport == nil || *row.Transport != "herdr" || row.State != "unknown" || row.HistoricalState == nil || *row.HistoricalState != "working" || row.StateFreshness != "unknown" || !row.Unresolved || row.TaskID == nil || *row.TaskID != "task" || row.Herdr == nil || row.Herdr.AgentSessionID == nil {
			t.Fatalf("invented liveness or lost binding: %+v", row)
		}
		if row.StartedAtUTC != nil {
			t.Fatalf("invented historical start time: %+v", row)
		}
	}
	if !providers["codex"] || !providers["claude"] {
		t.Fatal(providers)
	}
}

func TestReadInventoryKeepsValidSiblingsWhenOneRunIsCorrupt(t *testing.T) {
	d, r, _ := fixture(t)
	p := runPaths(r).RunRoot
	if err := os.MkdirAll(p, 0700); err != nil {
		t.Fatal(err)
	}
	if err := writeValue(filepath.Join(p, "request.json"), r); err != nil {
		t.Fatal(err)
	}
	badID := RunID("broken")
	bad := filepath.Join(storeRoot(r.WorkspaceRoot), "runs", badID)
	if err := os.MkdirAll(bad, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bad, "request.json"), []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	out, err := ReadInventory(d, r.WorkspaceRoot)
	if err != nil || len(out.Runs) != 2 {
		t.Fatalf("inventory: %+v, %v", out, err)
	}
	var valid, invalid *InventoryRun
	for i := range out.Runs {
		if out.Runs[i].RunID == badID {
			invalid = &out.Runs[i]
		} else {
			valid = &out.Runs[i]
		}
	}
	if valid == nil || valid.Transport == nil || *valid.Transport != "terminal" || valid.TaskID == nil || invalid == nil || invalid.Transport != nil || invalid.State != "unknown" || invalid.Freshness != "unknown" || len(invalid.Reasons) == 0 {
		t.Fatalf("valid=%+v invalid=%+v", valid, invalid)
	}
}

func TestReadInventoryFactoryTransportIsHeadless(t *testing.T) {
	d, r, _, _ := factoryFixture(t)
	if err := os.MkdirAll(runPaths(r).RunRoot, 0700); err != nil {
		t.Fatal(err)
	}
	if err := writeValue(filepath.Join(runPaths(r).RunRoot, "request.json"), r); err != nil {
		t.Fatal(err)
	}
	out, err := ReadInventory(d, r.WorkspaceRoot)
	if err != nil || len(out.Runs) != 1 || out.Runs[0].Transport == nil || *out.Runs[0].Transport != "headless" {
		t.Fatalf("inventory: %+v, %v", out, err)
	}
}

func TestReadInventoryQualifiedReturnSurvivesCheckoutDrift(t *testing.T) {
	d, r, file := fixture(t)
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(LaunchSpec) error {
		if _, err := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "claim.json", acceptance(t, d, r))); err != nil {
			return err
		}
		_, err := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), writeAny(t, r.WorkspaceRoot, "semantic-report.json", semanticReport(t, d, r)))
		return err
	}
	p, err := PreviewStart(d, file)
	if err != nil {
		t.Fatal(err)
	}
	run, err := startAndObserve(t, d, r, file, *p.(Preview).Confirmation)
	if err != nil || run.Collection.State != "qualified" {
		t.Fatalf("fixture qualification: %+v %v", run, err)
	}
	if err := os.WriteFile(filepath.Join(run.ObservedTarget.WorktreeLocator, "later.txt"), []byte("independent later change"), 0600); err != nil {
		t.Fatal(err)
	}
	before := treeState(t, r.WorkspaceRoot)
	rows, err := ReadInventory(d, r.WorkspaceRoot)
	if err != nil || len(rows.Runs) != 1 {
		t.Fatalf("inventory: %+v %v", rows, err)
	}
	row := rows.Runs[0]
	if row.State != "completed" || row.Unresolved || row.Return == nil || row.Return.ReportSHA256 == nil || row.Freshness != "fresh" {
		t.Fatalf("lost historical completion: %+v", row)
	}
	if treeState(t, r.WorkspaceRoot) != before {
		t.Fatal("inventory wrote historical state")
	}
}

func TestReadInventoryLegacyReviewedReturnAndDamagedSlot(t *testing.T) {
	f := workflowTestFixture(t)
	run := workflowTestReady(t, f)
	run, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, run.RunID, workflowTestReview(t, f, run, "accepted"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	before := treeState(t, f.R.WorkspaceRoot)
	rows, err := ReadInventory(f.D, f.R.WorkspaceRoot)
	if err != nil || len(rows.Runs) != 1 {
		t.Fatalf("inventory: %+v %v", rows, err)
	}
	row := rows.Runs[0]
	if row.State != "completed" || row.Unresolved || row.Return == nil || row.Return.ReportedOutcome == nil || *row.Return.ReportedOutcome != "complete" {
		t.Fatalf("lost reviewed return: %+v", row)
	}
	if treeState(t, f.R.WorkspaceRoot) != before {
		t.Fatal("inventory contacted provider or wrote state")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, run.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Records[0].Report.Locator, []byte("damaged"), 0600); err != nil {
		t.Fatal(err)
	}
	rows, err = ReadInventory(f.D, f.R.WorkspaceRoot)
	if err != nil || rows.Runs[0].State != "unknown" || !rows.Runs[0].Unresolved || rows.Runs[0].Freshness != "unknown" {
		t.Fatalf("damaged evidence was completed: %+v %v", rows, err)
	}
}

func TestInventoryDeliveryCompletionRequiresBoundSuccessfulIntegration(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	id, request := "wfr_fixture", hash([]byte("request"))
	makeBinding := func(completed bool) FileBinding {
		p := writeAny(t, root, "integration.json", map[string]any{"kind": "PlyDeliveryIntegration@1", "run_id": id, "request_sha256": request, "candidate": "candidate", "result": map[string]any{"Completed": completed}})
		return FileBinding{p, hashFileTest(t, p)}
	}
	check := func(contract *DeliveryContract, completed bool) error {
		binding := makeBinding(completed)
		s := workflowState{Request: WorkflowRequest{Delivery: contract}, Result: WorkflowRun{RunID: id, RequestSHA256: request, Delivery: &DeliveryState{Phase: "completed", Candidates: []DeliveryCandidate{{Key: "candidate", Integration: &binding}}, Events: []DeliveryEventRecord{{Kind: "integration", Binding: binding}}}}}
		row := InventoryRun{RunID: id, State: "completed", Return: &InventoryReturn{ReportSHA256: &binding.SHA256}}
		return workflowInventoryReturn(Dependencies{}, s, &row)
	}
	for _, tc := range []struct {
		name     string
		contract *DeliveryContract
	}{
		{"missing_contract", nil},
		{"missing_agreement", &DeliveryContract{}},
		{"explicit_local_agreement", &DeliveryContract{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalBranch}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := check(tc.contract, true); err != nil {
				t.Fatal(err)
			}
			if err := check(tc.contract, false); err == nil {
				t.Fatal("unsuccessful integration became completed")
			}
		})
	}
	row := InventoryRun{State: "completed"}
	if err := workflowInventoryReturn(Dependencies{}, workflowState{}, &row); err == nil {
		t.Fatal("phase alone became completed")
	}
}

func TestInventoryPullRequestCompletionRequiresFrozenContractAndExactPublication(t *testing.T) {
	for _, scenario := range []string{"valid", "missing_contract", "missing_agreement", "changed_head", "changed_base", "metadata_incomplete", "unbound_event", "changed_receipt"} {
		t.Run(scenario, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			id, request, oid := "wfr_fixture", hash([]byte("request")), strings.Repeat("a", 40)
			contract := &DeliveryContract{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryPullRequest, SourceRef: "refs/heads/task", TargetRef: "refs/heads/main", GitHubRepository: "fixture/product", Remote: "origin"}}
			observed := PullRequestDeliveryObservation{Repository: "fixture/product", HeadRef: "refs/heads/task", HeadOID: oid, BaseRef: "refs/heads/main", URL: "https://example.invalid/fixture/pull/1", Number: 1, MetadataApplied: true}
			receiptPath := writeAny(t, root, "receipt.json", observed)
			receipt := FileBinding{receiptPath, hashFileTest(t, receiptPath)}
			switch scenario {
			case "missing_contract":
				contract = nil
			case "missing_agreement":
				contract.Agreement = nil
			case "changed_head":
				observed.HeadOID = strings.Repeat("b", 40)
			case "changed_base":
				observed.BaseRef = "refs/heads/other"
			case "metadata_incomplete":
				observed.MetadataApplied = false
			case "changed_receipt":
				if err := os.WriteFile(receiptPath, []byte("changed evidence"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			path := writeAny(t, root, "publication.json", map[string]any{"kind": "PlyDeliveryPullRequest@1", "run_id": id, "request_sha256": request, "candidate": "candidate", "receipt": receipt, "observed": observed})
			binding := FileBinding{path, hashFileTest(t, path)}
			s := workflowState{Request: WorkflowRequest{Delivery: contract}, Result: WorkflowRun{RunID: id, RequestSHA256: request, Delivery: &DeliveryState{Phase: "completed", Candidates: []DeliveryCandidate{{Key: "candidate", OID: oid, PullRequest: &binding}}, Events: []DeliveryEventRecord{{Kind: "pull_request", Binding: binding}}}}}
			if scenario == "unbound_event" {
				s.Result.Delivery.Events = nil
			}
			row := InventoryRun{RunID: id, State: "completed", Return: &InventoryReturn{ReportSHA256: &binding.SHA256}}
			err = workflowInventoryReturn(Dependencies{}, s, &row)
			if (err == nil) != (scenario == "valid") {
				t.Fatalf("PR completion %s: %v", scenario, err)
			}
		})
	}
}
