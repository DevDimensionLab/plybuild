package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	"testing"
	"time"
)

// Modeled scenarios: the Herdr and provider stand-ins below are synthetic.
// They prove Ply's own reservations and bindings, not an actual provider resume.
const ownerResumeStandin = `
if model.get("closed_transport") and cmd == ["agent", "get"] and model.get("managed_inspection_denied"):
    print(json.dumps({"error":{"code":"permission_denied", "message":"inspection denied"}}), file=sys.stderr)
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["pane", "get"] and model.get("inspection_denied"):
    print(json.dumps({"id":"fixture", "error":{"code":"permission_denied", "message":"inspection denied"}}), file=sys.stderr)
    raise SystemExit(1)
if model.get("closed_transport") and cmd == ["agent", "get"] and args[2] == model.get("closed_pane", "w-fixture:p1"):
    print(json.dumps({"id":"fixture", "error":{"code":"agent_not_found", "message":"closed fixture agent"}}), file=sys.stderr)
    raise SystemExit(1)
if cmd == ["pane", "close"]:
    model["closed_transport"] = True
    model["closed_pane"] = args[2]
    path.write_text(json.dumps(model))
    print(json.dumps({"result":{"ok":True}}))
    raise SystemExit(0)
if cmd == ["pane", "get"] and model.get("replacement_count") and args[2] != model.get("closed_pane", "w-fixture:p1"):
    suffix = str(model["replacement_count"])
    pane = {"workspace_id":model.get("destination", "w-fixture"),"tab_id":"recovery-tab-"+suffix,"pane_id":"recovery-pane-"+suffix,"terminal_id":"recovery-terminal-"+suffix}
    print(json.dumps({"result":{"pane":pane}}))
    raise SystemExit(0)
if cmd == ["pane", "process-info"] and model.get("resume_inspection_denied"):
    print(json.dumps({"error":{"code":"permission_denied"}}),file=sys.stderr)
    raise SystemExit(1)
if cmd == ["pane", "process-info"]:
    print(json.dumps({"id":"fixture", "result":{"process_info":{"pane_id":args[3], "shell_pid":model.get("invoker_shell_pid"), "foreground_process_group_id":model.get("invoker_shell_pid"), "foreground_processes":model.get("resume_processes", [])}}}))
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

// RESUME01 scenario 2: a qualified candidate and its stored question survive
// the closed tab; the resumed owner continues at the same gate.

// RESUME01 scenario 3: every unsafe or unsupported case stops before effects
// with a concrete reason and next step.

// RESUME01 scenario 1 negative: a different conversation in the new tab is not
// the same owner; no Task prompt is sent to it.

// RESUME01 scenario 4: lost replies and concurrent calls keep one attempt.

// Coordinator review focus: only an actual callback through the resumed
// control counts. A coordinator stop or an installed-reader incomplete report
// increases events without proving the resumed owner works.

// Coordinator review focus: explicit lookup must not depend on unrelated
// historical Runs, while the selected Run keeps its full integrity checks.

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
	again, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, r)
	if err != nil || again.Confirmation != rp.Confirmation {
		t.Fatalf("fresh observation changed confirmation without changed resume facts: first=%s second=%s error=%v", rp.Confirmation, again.Confirmation, err)
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

// r001-f02/f03: a resumed session that Herdr discovers late is awaited, not
// rejected; a deadline leaves the same reserved attempt continuable once.

// r001-f06: a logical cwd through a symlink is the same physical worktree.

// r001-f04: without the provider on PATH the public CLI still reports the
// actual phase and next action, and names the missing provider precisely.

// r002-f01: Herdr may never report a resumed Codex session. The public result
// states that coverage honestly, and session_pending has a bounded exit:
// closing that pane lets a new generation start once Herdr proves it absent.

// r002-f02: scenario 2 through the existing QA contract. The stored blocked
// answer for the unchanged candidate survives the resume; the resumed control verifies a
// changed candidate, the old answer cannot qualify it, and the resumed
// control's own QA callback on the new candidate is accepted.
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

// r003-f01: an empty Herdr inventory does not prove the old writer stopped.
// A local process still in the Task worktree or holding the same provider
// session, or an unavailable observation, blocks before any tab. Only the
// invoking command chain up to the Herdr-identified invoker shell is excluded,
// and the reconciled basis is recorded on the positive path.

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

// ownerResumeCloseResumedTab models the human closing the resumed owner's tab
// after it worked; Herdr then reports neither its pane nor any Task writer.
func ownerResumeCloseResumedTab(t *testing.T, f ownerResumeFixture, o WorkflowRun) {
	t.Helper()
	workflowTestModel(t, f.workflowFixture, map[string]any{"closed_pane": o.Transport.PaneID, "inventory": []any{}, "agent_status": "idle"})
}

func ownerResumeApplyNext(t *testing.T, f ownerResumeFixture, generation int) (WorkflowRun, DeliveryOwnerResumePreview) {
	t.Helper()
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil || p.State != "ready" || p.Generation != generation {
		t.Fatalf("generation %d not offered: %+v %v", generation, p, err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil || o.OwnerResume == nil || o.OwnerResume.Generation != generation || o.OwnerResume.Step != "prompt_sent" {
		t.Fatalf("generation %d did not resume: %+v %v", generation, o.OwnerResume, err)
	}
	return o, p
}

// r005-f01: the same Task resumes twice in sequence. Each generation gets its
// own actual owner callback through its own control; the earlier callback stays
// history, the earlier control is closed, and show and lookup stay readable.

// r005-f01: e3a25d3 left a reserved, unsent generation 2 carrying generation
// 1's callback marker. Exactly that inherited marker, proven by generation 2's
// immutable before-state, reads as history without start authority; the next
// generation is offered only after absence and liveness are reconciled again.
