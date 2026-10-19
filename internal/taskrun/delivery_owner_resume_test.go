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
if cmd == ["pane", "process-info"]:
    print(json.dumps({"id":"fixture", "result":{"process_info":{"pane_id":args[3], "shell_pid":model.get("invoker_shell_pid"), "foreground_process_group_id":model.get("invoker_shell_pid"), "foreground_processes":[]}}}))
    raise SystemExit(0)
if model.get("closed_transport") and cmd == ["agent", "list"]:
    print(json.dumps({"id":"fixture", "result":{"agents":model.get("recovery_agents", [])}}))
    raise SystemExit(0)
`

// Real providers can be idle in Herdr before their resumed session is
// discovered. The first N observations of the new pane report no session.
const ownerResumeSessionStandin = `
if cmd == ["agent", "get"] and model.get("replacement_count") and model.get("session_absent_gets", 0) > 0:
    model["session_absent_gets"] -= 1
    agent["agent_session"] = None
`

type ownerResumeFixture struct {
	deliveryFixture
	Owner WorkflowRun
	In    DeliveryOwnerResumeInput
}

func newOwnerResumeFixture(t *testing.T, provider string, qualified bool, beforeClose ...func(deliveryFixture, WorkflowRun) WorkflowRun) ownerResumeFixture {
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
	herdr = bytes.Replace(herdr, []byte("if \"foreground_cwd\" in model:\n"), []byte(closedTransportAgentStandin+ownerResumeSessionStandin+"if \"foreground_cwd\" in model:\n"), 1)
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
	for _, fn := range beforeClose {
		o = fn(f, o)
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
			if err != nil || p.State != "ready" || p.Phase != "begun_work" || p.NativeSessionID != f.Owner.Transport.AgentSessionID || p.Absence == nil || p.Absence.ProcessLiveness != "no_task_or_session_process_observed" || p.Liveness == nil || p.Confirmation == "" {
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
	// r001-f03: an unknown tab creation is never replayed while a pane may
	// exist; once Herdr proves it absent, a new generation is the next step.
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
		task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
		workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": "", "inventory": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "stray-tab", "pane_id": "stray-pane", "terminal_id": "stray-terminal", "cwd": task}}})
		for n := 0; n < 2; n++ {
			again, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if err == nil || again.State != "blocked" || again.Reason != "owner_active" || again.NextAction.Message == "" {
				t.Fatalf("possibly created tab was not reported with a next step: %+v %v", again, err)
			}
			_, _ = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
		}
		if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
			t.Fatal("uncertain tab creation was replayed")
		}
		workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{}})
		next, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil || next.State != "ready" || next.Generation != 2 {
			t.Fatalf("proven-absent unknown tab has no executable next step: %+v %v", next, err)
		}
		ownerResumePreserve(t, f, next)
		o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, next.Confirmation)
		if err != nil || o.OwnerResume == nil || o.OwnerResume.Generation != 2 || o.OwnerResume.Step != "prompt_sent" {
			t.Fatalf("second generation did not resume: %+v %v", o.OwnerResume, err)
		}
		if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatal("resume prompt count differs from one")
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
		// r001-f03: the human gets an executable step, and an actual callback
		// through the resumed control proves the prompt arrived.
		again, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err == nil || again.Reason != "prompt_unconfirmed" || !strings.Contains(again.NextAction.Message, o.Transport.PaneID) || again.NextAction.Actor != "user" {
			t.Fatalf("uncertain prompt has no executable next step: %+v %v", again, err)
		}
		owner := f.D
		owner.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
		o = ownerResumeReport(t, f.deliveryFixture, o, owner, "resumed-after-lost-prompt")
		if o.OwnerResume.Step != "prompt_sent" || o.OwnerResume.OwnerCallback != "observed" {
			t.Fatalf("resumed callback did not settle the uncertain prompt: %+v", o.OwnerResume)
		}
		if live, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In); err != nil || live.State != "live_owner" {
			t.Fatalf("resumed owner not live: %+v %v", live, err)
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
	if !found || len(all.Unreadable) != 1 || all.Unreadable[0].RunID != "wfr_"+strings.Repeat("a", 64) || all.Unreadable[0].Error == "" {
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

// r001-f01: the human runs --resume from a plain shell in the Task worktree
// inside Herdr. That identified invoker pane runs no provider and is not a
// remaining writer; any other Task pane or an agent in the invoker still is.
func TestOwnerResumeAllowsIdentifiedInvokerShellInTaskWorktree(t *testing.T) {
	task := func(f ownerResumeFixture) string { return f.Prepared.Preparation.Preparation.Plan.WorktreePath }
	shell := func(f ownerResumeFixture, id string) map[string]any {
		return map[string]any{"workspace_id": "w-fixture", "tab_id": "tab-" + id, "pane_id": id, "terminal_id": "term-" + id, "cwd": task(f), "foreground_cwd": task(f)}
	}
	t.Run("invoker_shell", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		f.In.InvokerPaneID = "invoker-pane"
		workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{shell(f, "invoker-pane")}})
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil || p.State != "ready" || p.InvokerPaneID != "invoker-pane" {
			t.Fatalf("invoking shell blocked its own resume: %+v %v", p, err)
		}
		ownerResumePreserve(t, f, p)
		if o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err != nil || o.OwnerResume == nil || o.OwnerResume.Step != "prompt_sent" {
			t.Fatalf("resume from the Task worktree failed: %v", err)
		}
	})
	t.Run("unidentified_task_shell", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		f.In.InvokerPaneID = "invoker-pane"
		workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{shell(f, "invoker-pane"), shell(f, "other-pane")}})
		if p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In); err == nil || p.Reason != "owner_active" {
			t.Fatalf("another Task pane was excluded: %+v %v", p, err)
		}
	})
	t.Run("agent_in_invoker", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		f.In.InvokerPaneID = "invoker-pane"
		pane := shell(f, "invoker-pane")
		pane["agent"] = "claude"
		workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{pane}, "recovery_agents": []any{pane}})
		if p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In); err == nil || p.Reason != "owner_active" {
			t.Fatalf("an agent in the invoker pane was excluded: %+v %v", p, err)
		}
	})
}

// r001-f02/f03: a resumed session that Herdr discovers late is awaited, not
// rejected; a deadline leaves the same reserved attempt continuable once.
func TestOwnerResumeAwaitsLateSessionAndContinuesReservedAttempt(t *testing.T) {
	t.Run("late_session_within_deadline", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 2})
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil {
			t.Fatal(err)
		}
		ownerResumePreserve(t, f, p)
		o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
		if err != nil || o.OwnerResume == nil || o.OwnerResume.Step != "prompt_sent" {
			t.Fatalf("late session discovery stranded the resume: %+v %v", o.OwnerResume, err)
		}
	})
	t.Run("deadline_then_continue", func(t *testing.T) {
		f := newOwnerResumeFixture(t, "claude", false)
		workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 1000})
		f.In.Timeout = 1500 * time.Millisecond
		p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil {
			t.Fatal(err)
		}
		ownerResumePreserve(t, f, p)
		if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
			t.Fatal("undiscovered session was reported as resumed")
		}
		if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 1 {
			t.Fatalf("input was sent before the session was observed: %d", n)
		}
		f.In.Timeout = 3 * time.Second
		starting, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err == nil || starting.Reason != "session_pending" || starting.NextAction.Message == "" {
			t.Fatalf("pending session has no next step: %+v %v", starting, err)
		}
		workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 0})
		c, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		if err != nil || c.State != "continue_ready" || c.Generation != 1 || c.Confirmation == "" {
			t.Fatalf("observed same session cannot continue the reserved attempt: %+v %v", c, err)
		}
		o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, c.Confirmation)
		if err != nil || o.OwnerResume == nil || o.OwnerResume.Generation != 1 || o.OwnerResume.Step != "prompt_sent" {
			t.Fatalf("reserved attempt did not continue: %+v %v", o.OwnerResume, err)
		}
		if workflowTestCalls(t, f.workflowFixture, "tab create") != 2 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatal("continuation repeated a tab, provider start or prompt")
		}
		if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, c.Confirmation); err != nil || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatalf("repeated continuation replayed the prompt: %v", err)
		}
	})
}

// r001-f06: a logical cwd through a symlink is the same physical worktree.
func TestOwnerResumeAcceptsSymlinkedTaskCWD(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	link := filepath.Join(t.TempDir(), "task-link")
	if err := os.Symlink(f.Prepared.Preparation.Preparation.Plan.WorktreePath, link); err != nil {
		t.Fatal(err)
	}
	f.D.CWD = func() (string, error) { return link, nil }
	if p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In); err != nil || p.State != "ready" {
		t.Fatalf("symlinked Task cwd was rejected: %+v %v", p, err)
	}
}

// r001-f04: without the provider on PATH the public CLI still reports the
// actual phase and next action, and names the missing provider precisely.
func TestOwnerResumePublicCLIWithoutProviderKeepsPhaseAndReason(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newOwnerResumeFixture(t, "claude", false)
	t.Setenv("HERDR_WORKSPACE_ID", "")
	task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	var path []string
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir != filepath.Join(f.R.WorkspaceRoot, "provider-bin") {
			if _, err := os.Stat(filepath.Join(dir, "claude")); err != nil {
				path = append(path, dir)
			}
		}
	}
	t.Setenv("PATH", strings.Join(path, string(os.PathListSeparator)))
	out, output, err := ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err == nil || !strings.Contains(string(out["preview"]), `"reason":"provider_unavailable"`) || deliveryCLIString(t, out, "next_action") == "" {
		t.Fatalf("missing provider hid the precise reason: %v\n%s", err, output)
	}
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, f.Owner.RunID, func(s *workflowState) error {
		s.Result.Delivery.Phase = "completed"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	out, output, err = ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err != nil || deliveryCLIString(t, out, "state") != "completed" || deliveryCLIString(t, out, "next_action") == "" {
		t.Fatalf("completed phase was hidden by provider lookup: %v\n%s", err, output)
	}
}

// r002-f01: Herdr may never report a resumed Codex session. The public result
// states that coverage honestly, and session_pending has a bounded exit:
// closing that pane lets a new generation start once Herdr proves it absent.
func TestOwnerResumeSessionPendingHasBoundedExit(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 100000})
	f.In.Timeout = 1500 * time.Millisecond
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err == nil {
		t.Fatal("unreported Codex session was treated as resumed")
	}
	pending, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err == nil || pending.Reason != "session_pending" || !strings.Contains(pending.NextAction.Message, "close") || !strings.Contains(pending.NextAction.Message, "retire") {
		t.Fatalf("session_pending has no bounded exit: %+v %v", pending.NextAction, err)
	}
	// The human closes the unproven pane; Herdr then proves it absent.
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_pane": o.Transport.PaneID, "session_absent_gets": 0, "inventory": []any{}})
	next, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || next.State != "ready" || next.Generation != 2 {
		t.Fatalf("closing the unproven pane did not allow a new generation: %+v %v", next, err)
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 1 {
		t.Fatalf("a prompt was sent to an unproven session: %d", n)
	}
	if !strings.Contains(OwnerResumeProviderSupport("codex"), "unverified") || strings.Contains(OwnerResumeProviderSupport("claude"), "unverified") {
		t.Fatalf("public provider coverage is not stated honestly: %q / %q", OwnerResumeProviderSupport("codex"), OwnerResumeProviderSupport("claude"))
	}
}

// r002-f02: scenario 2 through the existing QA contract. The stored blocked
// answer for the unchanged candidate survives the resume; the resumed control verifies a
// changed candidate, the old answer cannot qualify it, and the resumed
// control's own QA callback on the new candidate is accepted.
func TestOwnerResumeKeepsStoredQAAndRequalifiesChangedCandidate(t *testing.T) {
	var first DeliveryCandidate
	f := newOwnerResumeFixture(t, "codex", true, func(f deliveryFixture, o WorkflowRun) WorkflowRun {
		// The stored actual answer for this candidate: blocked, so the owner
		// returns to correction work after the resume.
		o, err := WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "blocked", deliveryTestHuman(t, f, o, "blocked"))
		if err != nil {
			t.Fatal(err)
		}
		first = o.Delivery.Candidates[0]
		return o
	})
	if first.HumanQA == nil || f.Owner.DeliveryStatus.HumanJudgment == nil || f.Owner.DeliveryStatus.HumanJudgment.Outcome != "blocked" {
		t.Fatalf("fixture has no stored QA answer: %+v", f.Owner.DeliveryStatus.HumanJudgment)
	}
	oldAnswer := deliveryTestHuman(t, f.deliveryFixture, f.Owner, "pass")
	// The pass attestation above names candidate 1 and is never submitted for it.
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(o.Delivery.Candidates[0], first) || !equal(o.DeliveryStatus.HumanJudgment, f.Owner.DeliveryStatus.HumanJudgment) {
		t.Fatal("resume changed the stored QA answer for the unchanged candidate")
	}
	owner := f.D
	owner.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	deliveryTestCandidate(t, f.deliveryFixture, "changed after resume")
	o, err = WorkflowDeliveryVerify(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f.deliveryFixture, "after-resume"))
	if err != nil {
		t.Fatalf("resumed control could not verify a changed candidate: %v", err)
	}
	if o.OwnerResume == nil || o.OwnerResume.OwnerCallback != "observed" {
		t.Fatalf("resumed verify callback not observed: %+v", o.OwnerResume)
	}
	latest := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	if latest.OID == first.OID || latest.HumanQA != nil || o.DeliveryStatus.HumanJudgment != nil && o.DeliveryStatus.HumanJudgment.Current {
		t.Fatalf("old QA was reused for the changed candidate: %+v", o.DeliveryStatus.HumanJudgment)
	}
	if _, err = WorkflowDeliveryQA(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", oldAnswer); err == nil {
		t.Fatal("the old candidate's stored answer qualified the changed candidate")
	}
	o, err = WorkflowDeliveryQA(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f.deliveryFixture, o, "pass"))
	if err != nil {
		t.Fatalf("resumed control QA callback on the changed candidate failed: %v", err)
	}
	if latest = o.Delivery.Candidates[len(o.Delivery.Candidates)-1]; latest.HumanQA == nil || o.DeliveryStatus.HumanJudgment == nil || !o.DeliveryStatus.HumanJudgment.Current || o.DeliveryStatus.HumanJudgment.CandidateOID != latest.OID {
		t.Fatalf("resumed QA is not bound to the changed candidate: %+v", o.DeliveryStatus.HumanJudgment)
	}
	if !equal(o.Delivery.Candidates[0], first) {
		t.Fatal("history of the first candidate changed")
	}
}

func ownerResumeInput(t *testing.T, binary, cwd string, args ...string) map[string]json.RawMessage {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow", "input"}, args...)...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	if err := c.Run(); err != nil {
		t.Fatalf("workflow input %v: %v\n%s%s", args, err, stdout.String(), stderr.String())
	}
	var out map[string]json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		t.Fatalf("workflow input %v output: %v\n%s", args, err, stdout.String())
	}
	return out
}

// r002-f02 / RESUME01 scenario 2 (modeled provider, real human-input store):
// the owner asked a cooperative-agent question and the human answered while
// the tab was closed. After the resume, the same question, request and bytes
// are read and used through the existing contract; nothing is republished.
func TestOwnerResumeKeepsStoredHumanInputAnswer(t *testing.T) {
	binary := deliveryCLIBinary(t)
	var question, requestFile string
	var request map[string]json.RawMessage
	f := newOwnerResumeFixture(t, "claude", false, func(f deliveryFixture, o WorkflowRun) WorkflowRun {
		task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
		draft := writeAny(t, f.R.WorkspaceRoot, "question.json", map[string]any{"kind": "ply.workflow.human-input-question-draft", "schema_version": 1, "publication_key": "resume01/choice", "project_id": "ply", "repo_id": "ply", "task_id": "task", "activity_id": "resume01/owner", "activity_title": "Resumed owner fixture", "provider": "claude", "session_id": o.Transport.AgentSessionID, "origin_cwd": task, "text": "Which synthetic variant should the fixture use?", "reason": "Synthetic isolated question; no real human product decision.", "answer_kind": "choice", "options": []any{map[string]any{"id": "short", "label": "Short", "consequence": "Short fixture."}, map[string]any{"id": "full", "label": "Full", "consequence": "Full fixture."}}})
		asked := ownerResumeInput(t, binary, task, "ask", "--workspace", f.R.WorkspaceRoot, "--file", draft)
		question = deliveryCLIString(t, asked, "question_id")
		// The synthetic human answers while the owner's tab is closing.
		preview := ownerResumeInput(t, binary, f.R.WorkspaceRoot, "answer", "preview", question, "--workspace", f.R.WorkspaceRoot, "--request-id", "resume01-answer-0001", "--option-id", "full")
		if err := json.Unmarshal(preview["request"], &request); err != nil {
			t.Fatal(err)
		}
		requestFile = writeAny(t, f.R.WorkspaceRoot, "answer-request.json", request)
		ownerResumeInput(t, binary, f.R.WorkspaceRoot, "answer", "apply", "--workspace", f.R.WorkspaceRoot, "--file", requestFile)
		return o
	})
	store := filepath.Join(f.R.WorkspaceRoot, ".ply", "human-input")
	before := deliveryRecoveryFiles(t, store)
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err != nil {
		t.Fatal(err)
	}
	after := deliveryRecoveryFiles(t, store)
	if len(after) != len(before) {
		t.Fatal("resume published or removed human-input records")
	}
	deliveryRecoveryUnchanged(t, before, "")
	// The resumed owner, same session and cwd, finds the same question again
	// without a new publication and reads the exact stored answer.
	task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	again := ownerResumeInput(t, binary, task, "ask", "--workspace", f.R.WorkspaceRoot, "--file", filepath.Join(f.R.WorkspaceRoot, "question.json"))
	if deliveryCLIString(t, again, "question_id") != question {
		t.Fatal("resumed owner published a duplicate question")
	}
	if len(deliveryRecoveryFiles(t, store)) != len(before) {
		t.Fatal("asking again after resume wrote a new publication")
	}
	waited := ownerResumeInput(t, binary, task, "wait", question, "--workspace", f.R.WorkspaceRoot, "--timeout-seconds", "1")
	var answered struct {
		Question struct {
			State  string `json:"state"`
			Answer struct {
				RequestID     string  `json:"request_id"`
				RequestSHA256 string  `json:"request_sha256"`
				OptionID      *string `json:"option_id"`
			} `json:"answer"`
		} `json:"question"`
	}
	raw, _ := json.Marshal(waited)
	if err = json.Unmarshal(raw, &answered); err != nil {
		t.Fatal(err)
	}
	a := answered.Question.Answer
	if answered.Question.State != "answered" || a.RequestID != "resume01-answer-0001" || a.OptionID == nil || *a.OptionID != "full" {
		t.Fatalf("resumed owner did not read the stored answer: %s", raw)
	}
	ack := ownerResumeInput(t, binary, task, "acknowledge", question, "--workspace", f.R.WorkspaceRoot, "--request-id", a.RequestID, "--request-sha256", a.RequestSHA256)
	if deliveryCLIString(t, ack, "state") != "agent_received" || deliveryCLIString(t, ack, "question_id") != question {
		t.Fatalf("resumed owner did not receive the stored answer: %v", ack)
	}
	// Receipt is distinct from use: the resumed owner applies the selected
	// option to a controlled fixture result, then reports the continuation.
	result := filepath.Join(f.R.WorkspaceRoot, "resumed-variant.txt")
	if err = os.WriteFile(result, []byte("variant="+*a.OptionID+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	summary := "Resumed owner applied variant " + *a.OptionID + " to the fixture result."
	continued := ownerResumeInput(t, binary, task, "continue", question, "--workspace", f.R.WorkspaceRoot, "--request-id", a.RequestID, "--request-sha256", a.RequestSHA256, "--summary", summary)
	var used struct {
		QuestionID   string `json:"question_id"`
		State        string `json:"state"`
		Continuation struct {
			RequestID     string `json:"request_id"`
			RequestSHA256 string `json:"request_sha256"`
			Summary       string `json:"summary"`
		} `json:"continuation"`
	}
	raw, _ = json.Marshal(continued)
	if err = json.Unmarshal(raw, &used); err != nil {
		t.Fatal(err)
	}
	if used.QuestionID != question || used.State != "agent_reported_continued" || used.Continuation.RequestID != a.RequestID || used.Continuation.RequestSHA256 != a.RequestSHA256 || used.Continuation.Summary != summary {
		t.Fatalf("continuation is not bound to the same question and answer: %s", raw)
	}
	if got, _ := os.ReadFile(result); string(got) != "variant=full\n" {
		t.Fatalf("fixture result does not use the stored answer: %q", got)
	}
}

// r003-f01: an empty Herdr inventory does not prove the old writer stopped.
// A local process still in the Task worktree or holding the same provider
// session, or an unavailable observation, blocks before any tab. Only the
// invoking command chain up to the Herdr-identified invoker shell is excluded,
// and the reconciled basis is recorded on the positive path.
func TestOwnerResumeReconcilesWriterLivenessBeforeNewWriter(t *testing.T) {
	self := os.Getpid()
	cases := []struct {
		name, state, reason, pid string
		invoker                  bool
		table                    func(task string) ([]LocalProcess, error)
	}{
		{"writer_in_task_worktree", "blocked", "writer_process_alive", "4242", false, func(task string) ([]LocalProcess, error) {
			return []LocalProcess{{PID: 1, CWD: "/"}, {PID: 4242, ParentPID: 1, CWD: filepath.Join(task, "src")}, {PID: self, ParentPID: 1, CWD: "/elsewhere"}}, nil
		}},
		{"session_holder_elsewhere", "blocked", "writer_process_alive", "4343", false, func(string) ([]LocalProcess, error) {
			return []LocalProcess{{PID: 1, CWD: "/"}, {PID: 4343, ParentPID: 1, CWD: "/somewhere/else", Session: "transcript_open"}, {PID: self, ParentPID: 1, CWD: "/elsewhere"}}, nil
		}},
		{"agent_ancestor_not_invoker", "blocked", "writer_process_alive", "800", true, func(task string) ([]LocalProcess, error) {
			// An agent between the invoker shell and this process is a writer.
			return []LocalProcess{{PID: 1, CWD: "/"}, {PID: 900, ParentPID: 1, CWD: task}, {PID: 800, ParentPID: 900, CWD: task}, {PID: self, ParentPID: 800, CWD: task}}, nil
		}},
		{"observation_unavailable", "blocked", "liveness_unknown", "", false, func(string) ([]LocalProcess, error) {
			return nil, os.ErrPermission
		}},
		{"only_invoking_chain", "ready", "", "", true, func(task string) ([]LocalProcess, error) {
			return []LocalProcess{{PID: 1, CWD: "/"}, {PID: 900, ParentPID: 1, CWD: task}, {PID: self, ParentPID: 900, CWD: task}, {PID: 901, ParentPID: self, CWD: task}}, nil
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newOwnerResumeFixture(t, "claude", false)
			task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
			if tc.invoker {
				f.In.InvokerPaneID = "invoker-pane"
				shellPID := 900
				if tc.name == "agent_ancestor_not_invoker" {
					shellPID = 900
				}
				workflowTestModel(t, f.workflowFixture, map[string]any{"invoker_shell_pid": shellPID, "inventory": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "tab-invoker", "pane_id": "invoker-pane", "terminal_id": "term-invoker", "cwd": task, "foreground_cwd": task}}})
			}
			f.D.LocalProcesses = func(session, transcript string) ([]LocalProcess, error) {
				if session != f.Owner.Transport.AgentSessionID || !strings.HasSuffix(transcript, session+".jsonl") {
					t.Fatalf("observer did not receive the bound session and transcript: %s %s", session, transcript)
				}
				return tc.table(task)
			}
			before := deliveryRecoveryFiles(t, f.Owner.Paths.RunRoot)
			p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
			if p.State != tc.state || p.Reason != tc.reason || (tc.state == "ready") != (err == nil) {
				t.Fatalf("liveness was not reconciled: %+v %v", p, err)
			}
			if tc.state == "blocked" {
				if p.NextAction.Message == "" || tc.pid != "" && !strings.Contains(p.Detail, tc.pid) {
					t.Fatalf("blocked liveness has no concrete next step: %+v", p)
				}
				if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, "sha256:"+strings.Repeat("1", 64)); err == nil {
					t.Fatal("apply started a writer without reconciled liveness")
				}
				deliveryRecoveryUnchanged(t, before, "")
				if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 {
					t.Fatal("a tab was created while writer liveness was unreconciled")
				}
				return
			}
			if p.Liveness == nil || p.Liveness.Method != "local_process_cwd_and_session_v1" || len(p.Liveness.RemainingPIDs) != 0 || p.Absence.ProcessLiveness != "no_task_or_session_process_observed" || !strings.Contains(p.Liveness.Limitation, "cwd") {
				t.Fatalf("positive basis is not stated: %+v %+v", p.Liveness, p.Absence)
			}
			ownerResumePreserve(t, f, p)
			o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
			if err != nil || o.OwnerResume == nil || o.OwnerResume.Step != "prompt_sent" {
				t.Fatalf("reconciled resume failed: %v", err)
			}
			raw, err := os.ReadFile(o.OwnerResume.Binding.Locator)
			if err != nil || !strings.Contains(string(raw), `"method":"local_process_cwd_and_session_v1"`) {
				t.Fatalf("resume record lacks its liveness basis: %s %v", raw, err)
			}
		})
	}
}

// ownerResumeLegacyRecord rewrites, inside this disposable fixture only, the
// latest generation record into the exact shape written by 5f1908d (no
// writer_liveness) and rebinds its digest, as that candidate would have.
func ownerResumeLegacyRecord(t *testing.T, f ownerResumeFixture, o WorkflowRun) {
	t.Helper()
	path := o.OwnerResume.Binding.Locator
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err = json.Unmarshal(raw, &record); err != nil {
		t.Fatal(err)
	}
	delete(record, "writer_liveness")
	legacy, err := Canonical(record)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	state, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	state = bytes.Replace(state, []byte(o.OwnerResume.Binding.SHA256), []byte(hash(legacy)), 1)
	if err = os.WriteFile(statePath, state, 0600); err != nil {
		t.Fatal(err)
	}
}

// r004-f01: a generation record written before writer_liveness existed stays
// readable as history (show, lookup, owner callback), states that its basis
// was not recorded, and never stands in for the gate on a new start.
func TestOwnerResumeReadsPreLivenessGenerationAsHistoryOnly(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f := newOwnerResumeFixture(t, "claude", false)
	t.Setenv("HERDR_WORKSPACE_ID", "")
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumeLegacyRecord(t, f, o)
	o, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || o.OwnerResume == nil || o.OwnerResume.LivenessBasis != "not_recorded_before_reconciliation" {
		t.Fatalf("pre-liveness generation is unreadable or claims a basis: %+v %v", o.OwnerResume, err)
	}
	owner := f.D
	owner.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	o = ownerResumeReport(t, f.deliveryFixture, o, owner, "after-legacy-record")
	if o.OwnerResume.OwnerCallback != "observed" {
		t.Fatal("owner callback lost on the historical generation")
	}
	task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	check, output, err := ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err != nil || deliveryCLIString(t, check, "state") != "live_owner" {
		t.Fatalf("public lookup cannot read the historical generation: %v\n%s", err, output)
	}
	// The resumed tab closes; the next generation must pass the live gate.
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_pane": o.Transport.PaneID, "inventory": []any{}})
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: 1, CWD: "/"}, {PID: 5151, ParentPID: 1, CWD: task}, {PID: os.Getpid(), ParentPID: 1, CWD: "/elsewhere"}}, nil
	}
	next, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err == nil || next.Reason != "writer_process_alive" || next.Generation != 0 && next.State == "ready" {
		t.Fatalf("historical generation stood in for the liveness gate: %+v %v", next, err)
	}
}
