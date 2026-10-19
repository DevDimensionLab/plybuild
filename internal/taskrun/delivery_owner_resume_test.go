package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Modeled scenarios: the Herdr and provider stand-ins below are synthetic.
// They prove Ply's own reservations and bindings, not an actual provider resume.
const ownerResumeStandin = `
if model.get("closed_transport") and cmd == ["pane", "get"] and model.get("inspection_denied"):
    print(json.dumps({"id":"fixture", "error":{"code":"permission_denied", "message":"inspection denied"}}), file=sys.stderr)
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["agent", "get"] and args[2] == model.get("closed_pane", "w-fixture:p1"):
    print(json.dumps({"id":"fixture", "error":{"code":"agent_not_found", "message":"closed fixture agent"}}), file=sys.stderr)
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["agent", "list"]:
    print(json.dumps({"id":"fixture", "result":{"agents":model.get("recovery_agents", [])}}))
    raise SystemExit(0)
`

type ownerResumeFixture struct {
	deliveryFixture
	Owner WorkflowRun
	In    DeliveryOwnerResumeInput
}

func newOwnerResumeFixture(t *testing.T, provider string, qualified bool) ownerResumeFixture {
	t.Helper()
	var f deliveryFixture
	if qualified {
		f = newDeliveryFixture(t, provider)
	} else {
		f = receiptFixture(t, provider, deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}})
	}
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+ownerResumeStandin+closedTransportInspectStandin), 1)
	herdr = bytes.Replace(herdr, []byte("if \"foreground_cwd\" in model:\n"), []byte(closedTransportAgentStandin+"if \"foreground_cwd\" in model:\n"), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "owner-resume-request.json", f.R)
	var o WorkflowRun
	if qualified {
		o = deliveryCLIQualified(t, f)
		o = reportCandidateQuestion(t, f, o, "question-before-closed-tab")
	} else {
		o = workflowTestStart(t, f.workflowFixture)
		o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: provider, NativeSession: o.Transport.AgentSessionID}, true)
		if err != nil {
			t.Fatal(err)
		}
		o = ownerResumeReport(t, f, o, f.D, "begun-work")
	}
	// The provider's own transcript for this session survives its closed tab.
	t.Setenv("HOME", filepath.Join(f.R.WorkspaceRoot, "provider-home"))
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	t.Setenv("CODEX_HOME", "")
	ownerResumeTranscript(t, f, o)
	bin := filepath.Join(f.R.WorkspaceRoot, "provider-bin")
	if err = os.MkdirAll(bin, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(f.R.Runtime.Executable.Path, filepath.Join(bin, provider)); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HERDR_ENV", "1")
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "agent_status": "idle"})
	f.D.HerdrTimeout = 3 * time.Second
	return ownerResumeFixture{f, o, DeliveryOwnerResumeInput{RunID: o.RunID, ProviderExecutable: f.R.Runtime.Executable, ControlExecutable: f.R.Runtime.PlyExecutable, Timeout: 3 * time.Second}}
}

func ownerResumeTranscript(t *testing.T, f deliveryFixture, o WorkflowRun) string {
	t.Helper()
	home := os.Getenv("HOME")
	path := filepath.Join(home, ".claude", "projects", "fixture-project", o.Transport.AgentSessionID+".jsonl")
	if f.R.Runtime.Provider == "codex" {
		path = filepath.Join(home, ".codex", "sessions", "2026", "10", "19", "rollout-2026-10-19T00-00-00-"+o.Transport.AgentSessionID+".jsonl")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func ownerResumeReport(t *testing.T, f deliveryFixture, o WorkflowRun, d Dependencies, id string) WorkflowRun {
	t.Helper()
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: id, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "Synthetic owner progress", Meaning: "Modeled owner callback in a disposable fixture.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	out, err := WorkflowDeliveryReport(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, id+".json", r))
	if err != nil {
		t.Fatalf("owner report %s: %v", id, err)
	}
	return out
}

func ownerResumePreserve(t *testing.T, f ownerResumeFixture, p DeliveryOwnerResumePreview) {
	t.Helper()
	deliveryRecoveryPreserveControl(t, DeliveryStartRecoveryInput{ControlExecutable: f.In.ControlExecutable}, DeliveryStartRecoveryPreview{ControlExecutable: p.ControlExecutable})
}

func ownerResumeArgv(t *testing.T, f deliveryFixture, command string) [][]string {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if err != nil {
		t.Fatal(err)
	}
	var out [][]string
	for _, line := range bytes.Split(b, []byte("\n")) {
		var x struct {
			Argv []string `json:"argv"`
		}
		if len(line) == 0 || json.Unmarshal(line, &x) != nil {
			continue
		}
		if strings.Join(x.Argv[:2], " ") == command {
			out = append(out, x.Argv)
		}
	}
	return out
}

// RESUME01 scenario 1: begun work loses its tab. The same conversation is
// resumed in one new visible tab; the resumed owner can call back and the
// replaced binding cannot compete.
func TestOwnerResumeContinuesBegunWorkInSameSession(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			f := newOwnerResumeFixture(t, provider, false)
			before := deliveryRecoveryFiles(t, f.Owner.Paths.RunRoot)
			p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err != nil || p.State != "ready" || p.Phase != "begun_work" || p.NativeSessionID != f.Owner.Transport.AgentSessionID || p.Absence == nil || p.Absence.ProcessLiveness != "unknown" || p.Confirmation == "" {
				t.Fatalf("begun work was not previewed for resume: %+v %v", p, err)
			}
			if !slices.Contains(p.ResumeArgv, f.Owner.Transport.AgentSessionID) {
				t.Fatalf("resume does not name the same provider session: %v", p.ResumeArgv)
			}
			deliveryRecoveryUnchanged(t, before, "")
			if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
				t.Fatal("preview created a tab")
			}
			ownerResumePreserve(t, f, p)
			o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
			if err != nil {
				t.Fatalf("resume failed: %v", err)
			}
			if o.RunID != f.Owner.RunID || o.Handoff != f.Owner.Handoff || o.Paths.Context != f.Owner.Paths.Context || o.Transport.AgentSessionID != f.Owner.Transport.AgentSessionID || o.Transport.PaneID == f.Owner.Transport.PaneID {
				t.Fatalf("resume changed Task ownership or kept the closed pane: %+v", o.Transport)
			}
			if o.OwnerResume == nil || o.OwnerResume.Generation != 1 || o.OwnerResume.Step != "prompt_sent" || o.OwnerResume.OwnerCallback != "awaiting" || o.OwnerResume.ControlExecutable != p.ControlExecutable {
				t.Fatalf("resume readback missing: %+v", o.OwnerResume)
			}
			if len(o.Delivery.Events) != len(f.Owner.Delivery.Events) || o.Delivery.Phase != f.Owner.Delivery.Phase {
				t.Fatal("resume changed preserved delivery history")
			}
			starts := ownerResumeArgv(t, f.deliveryFixture, "agent start")
			if len(starts) != 2 || ownerResumeArgv(t, f.deliveryFixture, "tab create") == nil || workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
				t.Fatal("resume did not create exactly one tab, provider and prompt")
			}
			resumed := starts[1]
			if !slices.Contains(resumed, f.Owner.Transport.AgentSessionID) || !slices.Contains(resumed, f.R.Runtime.Model) || resumed[slices.Index(resumed, "--pane")+1] != o.Transport.PaneID {
				t.Fatalf("provider was not resumed with its session and model in the new pane: %v", resumed)
			}
			// The replaced binding's control can no longer report.
			r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "old-binding", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "Late", Meaning: "Late callback from the replaced binding.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
			if _, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "old-binding.json", r)); err == nil {
				t.Fatal("the replaced control still delivered an owner callback")
			}
			owner := f.D
			owner.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
			o = ownerResumeReport(t, f.deliveryFixture, o, owner, "resumed-work")
			if o.OwnerResume == nil || o.OwnerResume.OwnerCallback != "observed" {
				t.Fatalf("functional owner callback is not visible: %+v", o.OwnerResume)
			}
			// Repeating the same confirmation never repeats a native effect.
			if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err != nil {
				t.Fatal(err)
			}
			if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
				t.Fatal("repeat replayed a native effect")
			}
			again, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err != nil || again.State != "live_owner" {
				t.Fatalf("resumed live owner was not recognized: %+v %v", again, err)
			}
		})
	}
}

// RESUME01 scenario 2: a qualified candidate and its stored question survive
// the closed tab; the resumed owner continues at the same gate.
func TestOwnerResumePreservesQualifiedCandidateAndQuestion(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", true)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" || p.Phase != "awaiting_human" {
		t.Fatalf("qualified candidate with question was not previewed: %+v %v", p, err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(o.Delivery.Candidates, f.Owner.Delivery.Candidates) || !equal(o.Delivery.Events, f.Owner.Delivery.Events) || o.NextAction != f.Owner.NextAction {
		t.Fatalf("resume changed candidate, stored question or gate: %+v", o.NextAction)
	}
	if o.DeliveryStatus == nil || f.Owner.DeliveryStatus == nil || o.DeliveryStatus.QualifiedCandidate == nil || !equal(o.DeliveryStatus.QualifiedCandidate, f.Owner.DeliveryStatus.QualifiedCandidate) {
		t.Fatal("resume lost the qualified candidate")
	}
}

// RESUME01 scenario 3: every unsafe or unsupported case stops before effects
// with a concrete reason and next step.
func TestOwnerResumeBoundariesStopBeforeEffects(t *testing.T) {
	cases := []struct {
		name, code string
		setup      func(*testing.T, ownerResumeFixture)
	}{
		{"other_writer_in_task", "owner_active", func(t *testing.T, f ownerResumeFixture) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_agents": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "t2", "pane_id": "p2", "terminal_id": "term2", "agent": "codex", "agent_status": "working", "foreground_cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath}}})
		}},
		{"missing_provider_session", "provider_session_missing", func(t *testing.T, f ownerResumeFixture) {
			if err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".claude")); err != nil {
				t.Fatal(err)
			}
		}},
		{"inspection_denied", "inspection_unavailable", func(t *testing.T, f ownerResumeFixture) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"inspection_denied": true})
		}},
		{"unknown_effect", "unknown_effect", func(t *testing.T, f ownerResumeFixture) {
			retireStopFixture(t, f.deliveryFixture, f.Owner, DeliveryEffectUnknown)
		}},
		{"unresolved_attempt", "unknown_effect", func(t *testing.T, f ownerResumeFixture) {
			if err := workflowUpdate(f.D, f.R.WorkspaceRoot, f.Owner.RunID, func(s *workflowState) error {
				s.Result.Delivery.Attempt = &DeliveryAttempt{Kind: "integration", State: "attempted"}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		}},
		{"retired_owner", "owner_retired", func(t *testing.T, f ownerResumeFixture) {
			coordinator := f.deliveryFixture
			retireAsCoordinator(&coordinator)
			if _, err := WorkflowRetireDelivery(coordinator.D, f.R.WorkspaceRoot, f.Owner.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Selected replacement"}); err != nil {
				t.Fatal(err)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newOwnerResumeFixture(t, "claude", false)
			tc.setup(t, f)
			before := deliveryRecoveryFiles(t, f.Owner.Paths.RunRoot)
			p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err == nil || p.State != "blocked" || p.Reason != tc.code || p.NextAction.Message == "" || p.NextAction.Actor == "" {
				t.Fatalf("expected blocked %s with next step: %+v %v", tc.code, p, err)
			}
			if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, "sha256:"+strings.Repeat("0", 64)); err == nil {
				t.Fatal("apply passed a blocked resume")
			}
			deliveryRecoveryUnchanged(t, before, "")
			if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
				t.Fatal("blocked resume started a native effect")
			}
		})
	}
	t.Run("live_owner", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": false})
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil || p.State != "live_owner" || !strings.Contains(p.NextAction.Message, f.Owner.Transport.PaneID) {
			t.Fatalf("live owner was not found without a new tab: %+v %v", p, err)
		}
		if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
			t.Fatal("live owner got another tab")
		}
	})
	t.Run("completed", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		if err := workflowUpdate(f.D, f.R.WorkspaceRoot, f.Owner.RunID, func(s *workflowState) error {
			s.Result.Delivery.Phase = "completed"
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil || p.State != "completed" || p.Confirmation != "" {
			t.Fatalf("completed delivery was offered a resume: %+v %v", p, err)
		}
	})
}

// RESUME01 scenario 1 negative: a different conversation in the new tab is not
// the same owner; no Task prompt is sent to it.
func TestOwnerResumeRejectsDifferentResumedSession(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_session_id": "another-session"})
	if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
		t.Fatal("a different provider session was accepted as the resumed owner")
	}
	if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("resume prompted an unverified session")
	}
	o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, f.Owner.RunID)
	if err != nil || o.OwnerResume == nil || o.OwnerResume.Step != "agent_start_attempted" {
		t.Fatalf("unverified resume is not visible as unresolved: %+v %v", o.OwnerResume, err)
	}
}

// RESUME01 scenario 4: lost replies and concurrent calls keep one attempt.
func TestOwnerResumeLostRepliesAndConcurrencyKeepOneAttempt(t *testing.T) {
	t.Run("lost_tab_reply", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil {
			t.Fatal(err)
		}
		ownerResumePreserve(t, f, p)
		workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": "tab create"})
		if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
			t.Fatal("lost tab reply was reported as confirmed")
		}
		for n := 0; n < 2; n++ {
			again, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err == nil || again.State != "blocked" || again.Reason != "resume_unknown" {
				t.Fatalf("unknown tab creation could be replaced: %+v %v", again, err)
			}
			_, _ = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
		}
		if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
			t.Fatal("uncertain tab creation was replayed")
		}
	})
	t.Run("lost_prompt_reply", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil {
			t.Fatal(err)
		}
		ownerResumePreserve(t, f, p)
		workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": "agent prompt"})
		if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
			t.Fatal("lost prompt reply was reported as delivered")
		}
		workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": ""})
		_, _ = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
		o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, f.Owner.RunID)
		if err != nil || o.OwnerResume == nil || o.OwnerResume.Step != "prompt_attempted" {
			t.Fatalf("uncertain prompt is not visible: %+v %v", o.OwnerResume, err)
		}
		if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatal("uncertain resume prompt was resent")
		}
	})
	t.Run("concurrent_calls", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil {
			t.Fatal(err)
		}
		ownerResumePreserve(t, f, p)
		var wg sync.WaitGroup
		for n := 0; n < 3; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
			}()
		}
		wg.Wait()
		if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatal("concurrent resume calls duplicated a native effect")
		}
	})
}

// Coordinator review focus: only an actual callback through the resumed
// control counts. A coordinator stop or an installed-reader incomplete report
// increases events without proving the resumed owner works.
func TestOwnerResumeCallbackIsSourcedToResumedControl(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	coordinator := f.deliveryFixture
	retireAsCoordinator(&coordinator)
	o = retireStopFixture(t, coordinator, o, DeliveryEffectNotStarted)
	if o.OwnerResume == nil || o.OwnerResume.OwnerCallback != "awaiting" {
		t.Fatalf("coordinator stop was counted as an owner callback: %+v", o.OwnerResume)
	}
	installed := f.D
	installed.Executable = func() (string, error) { return f.R.Runtime.PlyExecutable.Path, nil }
	question := "Is the resumed owner still needed?"
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "incomplete-installed", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "needs_input", Question: &question, Summary: "Incomplete", Meaning: "Installed-reader report, not the resumed control.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	o, err = WorkflowDeliveryIncompleteReport(installed, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "incomplete.json", r))
	if err != nil {
		t.Fatal(err)
	}
	if o.OwnerResume == nil || o.OwnerResume.OwnerCallback != "awaiting" {
		t.Fatalf("incomplete installed-reader report was counted as the resumed owner's callback: %+v", o.OwnerResume)
	}
}

// Coordinator review focus: explicit lookup must not depend on unrelated
// historical Runs, while the selected Run keeps its full integrity checks.
func TestOwnerRunLookupIgnoresUnrelatedInvalidHistory(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	requests := filepath.Join(workflowRoot(f.R.WorkspaceRoot), "requests")
	unrelated := filepath.Join(requests, "wfr_"+strings.Repeat("a", 64)+".json")
	if err := os.WriteFile(unrelated, []byte("{\"not\":\"a reservation\"}"), 0600); err != nil {
		t.Fatal(err)
	}
	all, err := WorkflowOwnerRunCandidates(f.D, f.R.WorkspaceRoot)
	if err != nil {
		t.Fatalf("unrelated invalid history blocked lookup: %v", err)
	}
	found := false
	for _, c := range all.Runs {
		found = found || c.RunID == f.Owner.RunID && c.SpecID != "" && c.TaskWorktree == f.Prepared.Preparation.Preparation.Plan.WorktreePath
	}
	if !found || len(all.Unreadable) != 1 || all.Unreadable[0] != "wfr_"+strings.Repeat("a", 64) {
		t.Fatalf("valid Run or unreadable history not reported: %+v", all)
	}
	// The selected Run itself still fails closed on changed bytes.
	state := filepath.Join(f.Owner.Paths.RunRoot, "state.json")
	raw, err := os.ReadFile(state)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(state, bytes.Replace(raw, []byte(`"phase"`), []byte(`"phase_changed"`), 1), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In); err == nil {
		t.Fatal("selected Run integrity was weakened")
	}
}

func ownerResumeCLI(binary, cwd string, args ...string) (map[string]json.RawMessage, string, error) {
	c := exec.Command(binary, append([]string{"workflow", "execute"}, args...)...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var out map[string]json.RawMessage
	_ = json.Unmarshal(stdout.Bytes(), &out)
	return out, stdout.String() + stderr.String(), err
}

// RESUME01 public CLI (modeled Herdr/provider): no Run ID, private control or
// pane is supplied. The resumed control is an actual copy of the built CLI and
// performs a real owner callback.
func TestOwnerResumePublicCLIFromTaskWorktree(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newOwnerResumeFixture(t, "claude", false)
	t.Setenv("HERDR_WORKSPACE_ID", "")
	task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	if _, output, err := ownerResumeCLI(binary, f.R.WorkspaceRoot, "--resume", "--check", "--format", "json"); err == nil {
		t.Fatalf("wrong cwd selected a Run: %s", output)
	}
	if _, output, err := ownerResumeCLI(binary, task, "--resume", "--restart", "--spec", "x"); err == nil || !strings.Contains(output, "--resume keeps") {
		t.Fatalf("--resume accepted --restart: %s", output)
	}
	before := deliveryRecoveryFiles(t, f.Owner.Paths.RunRoot)
	check, output, err := ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err != nil || deliveryCLIString(t, check, "state") != "ready" || !strings.Contains(string(check["run"]), f.Owner.RunID) || !strings.Contains(string(check["provider_support"]), "--resume") {
		t.Fatalf("public preview failed: %v\n%s", err, output)
	}
	if after := deliveryRecoveryFiles(t, f.Owner.Paths.RunRoot); len(after) != len(before) {
		t.Fatal("public preview wrote Run artifacts")
	}
	deliveryRecoveryUnchanged(t, before, "")
	text, output, err := ownerResumeCLI(binary, task, "--resume", "--timeout", "5")
	_ = text
	if err != nil || !strings.Contains(output, "Owner resume: generation 1 · step prompt_sent · same session observed true · owner callback awaiting") {
		t.Fatalf("public resume failed: %v\n%s", err, output)
	}
	o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, f.Owner.RunID)
	if err != nil || o.OwnerResume == nil {
		t.Fatal(err)
	}
	control := o.OwnerResume.ControlExecutable.Path
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "resumed-cli-callback", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "Resumed owner callback", Meaning: "Modeled provider; the resumed control binary is real.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	c := exec.Command(control, "workflow", "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, f.R.WorkspaceRoot, "resumed-cli.json", r))
	c.Dir = task
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("resumed control could not call back: %v\n%s", err, out)
	}
	if o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, f.Owner.RunID); err != nil || o.OwnerResume.OwnerCallback != "observed" || o.OwnerResume.CallbackEventID != "resumed-cli-callback" {
		t.Fatalf("functional callback not recorded: %+v %v", o.OwnerResume, err)
	}
	again, output, err := ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err != nil || deliveryCLIString(t, again, "state") != "live_owner" {
		t.Fatalf("resumed owner not seen as live: %v\n%s", err, output)
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
		t.Fatal("public CLI duplicated a native effect")
	}
}

// Compatibility: a Run that already replaced its startup transport keeps
// validating that recovery after a later owner resume replaces it again.
func TestOwnerResumeAfterStartupRecoveryKeepsRecoveryHistory(t *testing.T) {
	f := deliveryRecoveryNewFixture(t)
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+ownerResumeStandin+closedTransportInspectStandin), 1)
	herdr = bytes.Replace(herdr, []byte("if \"foreground_cwd\" in model:\n"), []byte(closedTransportAgentStandin+"if \"foreground_cwd\" in model:\n"), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil {
		t.Fatal("fixture did not stop before Task")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true, "recovery_armed": true})
	f.D.HerdrTimeout = 3 * time.Second
	in := DeliveryStartRecoveryInput{RunID: o.RunID, ProviderExecutable: f.R.Runtime.Executable, ControlExecutable: f.R.Runtime.PlyExecutable, Timeout: 3 * time.Second}
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	if o, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err != nil || o.StartupRecovery == nil {
		t.Fatalf("startup recovery fixture: %v", err)
	}
	f.D.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	o = deliveryTestAccept(t, f, o)
	o = ownerResumeReport(t, f, o, f.D, "after-recovery")
	// The recovered tab now closes as well.
	t.Setenv("HOME", filepath.Join(f.R.WorkspaceRoot, "provider-home"))
	t.Setenv("CODEX_HOME", "")
	ownerResumeTranscript(t, f, o)
	bin := filepath.Join(f.R.WorkspaceRoot, "provider-bin")
	if err = os.MkdirAll(bin, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(f.R.Runtime.Executable.Path, filepath.Join(bin, "codex")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	// After startup, the stand-in models an ordinary session-bearing provider.
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_pane": o.Transport.PaneID, "inventory": []any{}, "agent_status": "idle", "agent_session_id": o.Transport.AgentSessionID, "recovery_armed": false, "bootstrap_mode": "", "bootstrap_started": false})
	r := DeliveryOwnerResumeInput{RunID: o.RunID, ProviderExecutable: f.R.Runtime.Executable, ControlExecutable: f.R.Runtime.PlyExecutable, DestinationWorkspaceID: o.Transport.WorkspaceID, Timeout: 3 * time.Second}
	rp, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, r)
	if err != nil || rp.State != "ready" {
		t.Fatalf("recovered owner could not preview resume: %+v %v", rp, err)
	}
	deliveryRecoveryPreserveControl(t, DeliveryStartRecoveryInput{ControlExecutable: r.ControlExecutable}, DeliveryStartRecoveryPreview{ControlExecutable: rp.ControlExecutable})
	resumed, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, r, rp.Confirmation)
	if err != nil || resumed.StartupRecovery == nil || !equal(resumed.StartupRecovery, o.StartupRecovery) || resumed.Transport.PaneID == o.Transport.PaneID {
		t.Fatalf("resume lost or broke startup recovery history: %v", err)
	}
	f.D.Executable = func() (string, error) { return rp.ControlExecutable.Path, nil }
	ownerResumeReport(t, f, resumed, f.D, "after-resume")
}
