package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// All identities, failures and process observations here are synthetic.
func al33FailedReview(t *testing.T, model ...map[string]any) (deliveryFixture, WorkflowRun, Dependencies) {
	t.Helper()
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	reviewRoundPatchHerdr(t, &f)
	b, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	injection := `
if cmd == ["agent", "start"] and "--pane" in args and args[args.index("--pane")+1].startswith("w-fixture:r"):
    sys.stderr.write(model.get("al33_start_stderr", "SECRET-child-message"))
    raise SystemExit(1)
if cmd == ["agent", "startup"] and args[2] == "submit" and json.loads(Path(args[4]).read_text())["kind"] == "bootstrap" and model.get("al33_bootstrap"):
    time.sleep(3) # only this synthetic guarded bootstrap exhausts the shared fixture budget
    raise SystemExit(0)
if cmd == ["agent", "get"] and args[2].startswith("w-fixture:r"):
    if model.get("al33_bootstrap"):
        print(json.dumps({"result":{"agent":{"workspace_id":"w-fixture","tab_id":"tab-review-1","pane_id":"w-fixture:r1","terminal_id":"terminal-review-1","name":model["al33_name"],"agent":"codex","agent_status":"idle","interactive_ready":True,"agent_session":None}}}))
        raise SystemExit(0)
    if model.get("al33_live"):
        print(json.dumps({"result":{"agent":{"workspace_id":"w-fixture","tab_id":"tab-review-1","pane_id":"w-fixture:r1","terminal_id":"terminal-review-1","name":model["al33_name"],"agent":"codex","agent_status":"working","agent_session":{"agent":"codex","kind":"id","value":"synthetic-live-session"}}}}))
        raise SystemExit(0)
    print(json.dumps({"error":{"code":model.get("al33_observation_code","agent_not_found"),"message":"SECRET-observation"}}))
    raise SystemExit(1)
`
	b = []byte(strings.Replace(string(b), "cmd = args[:2]\n", "cmd = args[:2]\n"+injection, 1))
	if err = os.WriteFile(f.R.Herdr.Executable.Path, b, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(b)
	f.File = writeAny(t, f.R.WorkspaceRoot, "diagnostic-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	reviewRoundCandidate(t, f, "fixed diagnostic fixture")
	workflowTestModel(t, f.workflowFixture, map[string]any{"caller_shell": 12345})
	if len(model) > 0 {
		workflowTestModel(t, f.workflowFixture, model[0])
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"al33_name": deliveryReviewAgentName(o.RunID, 1)})
	d := reviewRoundAsCoordinator(f)
	d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), ParentPID: 12345}}, nil
	}
	d.HerdrTimeout = 2 * time.Second
	reviewRoundIn(t, "w-fixture:root")
	if _, err = WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic coordinator"}); err == nil {
		t.Fatal("expected synthetic startup failure")
	}
	return f, o, d
}

func TestAL33ReviewerStartupDiagnostics(t *testing.T) {
	f, o, d := al33FailedReview(t)
	shown := mustShow(t, f, o.RunID)
	round := reviewRoundLatest(t, shown)
	raw, _ := json.Marshal(round)
	if !strings.Contains(string(raw), `"first_failure"`) || !strings.Contains(string(raw), `"start_attempt":"attempted"`) || !strings.Contains(string(raw), `"task_prompt":"not_attempted"`) {
		t.Fatalf("missing durable phase diagnostics: %s", raw)
	}
	diag := round.StartupDiagnostics
	if diag.FirstFailure == nil || diag.FirstFailure.Phase != "agent start" || diag.FirstFailure.Class != "nonzero_exit" || diag.FirstFailure.Exit != "1" || diag.FirstFailure.ProviderCode != "unknown" || diag.FirstFailure.StderrBytes != 20 || diag.FirstFailureEvidence == nil || diag.LastObservation == nil || diag.ObservationEvidence == nil {
		t.Fatalf("first failure or observation lost: %s", raw)
	}
	if diag.Identity.PaneID != "w-fixture:r1" || diag.Identity.TabID != "tab-review-1" || diag.Identity.TerminalID != "terminal-review-1" || len(diag.InspectionArgv) != 3 || diag.InspectionArgv[2][3] != "--pane" {
		t.Fatalf("inspection identity: %s", raw)
	}
	if strings.Contains(string(raw), "SECRET") {
		t.Fatalf("child text leaked: %s", raw)
	}
	before := len(reviewRoundArgv(t, f, "agent start"))
	for i := 0; i < 2; i++ {
		mustShow(t, f, o.RunID)
		if _, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic coordinator"}); err != nil {
			t.Fatal(err)
		}
	}
	if len(reviewRoundArgv(t, f, "agent start")) != before || len(reviewRoundArgv(t, f, "agent prompt")) != 1 {
		t.Fatal("read/idempotence sent another effect")
	}
}

func TestAL33ReviewerObservationDoesNotProveExit(t *testing.T) {
	for _, tc := range []struct {
		name  string
		model map[string]any
		want  string
	}{
		{"permission_denied", map[string]any{"al33_observation_code": "permission_denied", "al33_start_stderr": `{"error":{"code":"permission_denied","message":"SECRET"}}`}, "permission_denied"},
		{"unknown", map[string]any{"al33_observation_code": "arbitrary-private-code"}, "unknown"},
		{"live", map[string]any{"al33_live": true}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, o, _ := al33FailedReview(t, tc.model)
			r := reviewRoundLatest(t, mustShow(t, f, o.RunID))
			diag := r.StartupDiagnostics
			if diag.ProviderLiveness != "unknown" || diag.SessionBinding != "not_bound" || diag.TaskPrompt != "not_attempted" {
				t.Fatalf("overstated observation: %+v", diag)
			}
			obs := diag.LastObservation
			if tc.name == "live" {
				if obs.SessionID != "synthetic-live-session" || obs.AgentStatus != "working" {
					t.Fatalf("lost living agent observation: %+v", obs)
				}
			} else if obs.FirstFailure == nil || obs.FirstFailure.ProviderCode != tc.want {
				t.Fatalf("lost failure: %+v", obs)
			}
			if tc.name == "permission_denied" && diag.FirstFailure.ProviderCode != "permission_denied" {
				t.Fatal("first refusal lost")
			}
		})
	}
}

// A real candidate CLI reads an isolated synthetic Run; it launches nothing.
func TestAL33ReviewerDiagnosticsCandidateCLI(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("requires exact make-built PLY_CANDIDATE_BINARY; no CLI built by this test")
	}
	f, o, _ := al33FailedReview(t)
	before := len(reviewRoundArgv(t, f, "agent start"))
	for _, format := range []string{"json", "text"} {
		output, err := reviewRoundCLI(t, binary, f.Parent, "w-fixture:root", "show", o.RunID, "--format", format)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"agent start", "nonzero_exit", "agent_not_found", "process-info", "--pane", "w-fixture:r1"} {
			if !strings.Contains(output, want) {
				t.Fatalf("candidate %s lacks %s: %s", format, want, output)
			}
		}
		if strings.Contains(output, "SECRET") {
			t.Fatal("untrusted child text leaked")
		}
	}
	// Execute only the synthetic stand-in's advertised read; a shell-only result
	// does not strengthen the preserved startup evidence or establish provider exit.
	diag := reviewRoundLatest(t, mustShow(t, f, o.RunID)).StartupDiagnostics
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": map[string]any{"w-fixture:r1": 12345}})
	argv := diag.InspectionArgv[2]
	c := exec.Command(argv[0], argv[1:]...)
	raw, err := c.Output()
	if err != nil || !strings.Contains(string(raw), "shell_pid") {
		t.Fatalf("synthetic inspection: %s %v", raw, err)
	}
	reread := reviewRoundLatest(t, mustShow(t, f, o.RunID)).StartupDiagnostics
	if reread.ProviderLiveness != "unknown" || reread.SessionBinding != "not_bound" || len(reviewRoundArgv(t, f, "agent start")) != before {
		t.Fatal("shell-only read changed startup or launched another agent")
	}
}

func TestAL33HistoricalReviewerEvidenceRemainsUnknown(t *testing.T) {
	f, o, _ := al33FailedReview(t)
	r := reviewRoundLatest(t, mustShow(t, f, o.RunID))
	dir := filepath.Dir(r.Request.Locator)
	for _, name := range []string{"start-failure.json", "session-observation.json"} {
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	diag := reviewRoundLatest(t, mustShow(t, f, o.RunID)).StartupDiagnostics
	if diag.FirstFailure != nil || diag.LastObservation != nil || diag.TaskPrompt != "not_recorded" || diag.BootstrapPrompt != "not_recorded" || diag.ProviderLiveness != "unknown" {
		t.Fatalf("invented historical evidence: %+v", diag)
	}
	// An old observation without explicit phase facts must not imply no prompt.
	writeAny(t, dir, "session-observation.json", deliveryReviewObservation{Outcome: "unavailable"})
	if d := reviewRoundLatest(t, mustShow(t, f, o.RunID)).StartupDiagnostics; d.TaskPrompt != "not_recorded" {
		t.Fatalf("failed read inferred no attempt: %+v", d)
	}
}

func TestAL33BootstrapDeadlineRetainsReservedPromptStage(t *testing.T) {
	f, o, d := al33FailedReview(t, map[string]any{"al33_bootstrap": true})
	round := reviewRoundLatest(t, mustShow(t, f, o.RunID))
	diag := round.StartupDiagnostics
	if diag.LastObservation == nil || diag.LastObservation.BootstrapPromptStage != "attempted" || diag.BootstrapPrompt != "attempted" || diag.LastObservation.TaskPromptStage != "not_attempted" || diag.TaskPrompt != "not_attempted" {
		t.Fatalf("bootstrap reservation contradicts retained observation: %+v / %+v", diag, diag.LastObservation)
	}
	prompts := reviewRoundArgv(t, f, "agent prompt")
	if len(prompts) != 2 || prompts[1][3] != workflowBootstrapPrompt {
		t.Fatalf("expected owner prompt and one tool-free bootstrap only: %q", prompts)
	}
	if _, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "synthetic coordinator"}); err != nil {
		t.Fatal(err)
	}
	mustShow(t, f, o.RunID)
	if len(reviewRoundArgv(t, f, "agent prompt")) != len(prompts) {
		t.Fatal("bootstrap or task prompt replayed")
	}
}
