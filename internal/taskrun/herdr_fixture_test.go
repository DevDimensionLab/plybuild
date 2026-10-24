package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const workflowStandin = `import json, sys, time, subprocess, os
from pathlib import Path
root = Path(__file__).parent
path = root / "herdr-model.json"
model = json.loads(path.read_text())
args = sys.argv[1:]
if model.get("before_call_delay_command") == " ".join(args[:2]):
    time.sleep(model["before_call_delay_millis"] / 1000)
call = {"argv": args, "at_ns": time.time_ns()}
if model.get("claude_trust_config") and args[:2] == ["agent", "start"]:
    cfg = json.loads(Path(model["claude_trust_config"]).read_text())
    call["claude_trusted_before_start"] = cfg.get("projects", {}).get(model["claude_trust_key"], {}).get("hasTrustDialogAccepted") is True
if model.get("bootstrap_mode"):
    run_root = Path(model["bootstrap_run_root"])
    state_path = run_root / "state.json"
    state = json.loads(state_path.read_text()) if state_path.exists() else {}
    call.update(bootstrap_reserved=(run_root / "startup-bootstrap-attempt.json").is_file(),
                phase_before_call=state.get("phase"),
                session_before_call=state.get("result", {}).get("transport", {}).get("agent_session_id"),
                bootstrap_settled_before_call=model.get("bootstrap_settled_observed", False))
with (root / "herdr-calls.jsonl").open("a") as f:
    f.write(json.dumps(call) + "\n")
if args[:3] == ["status", "server", "--json"]:
    print(json.dumps({"running":True,"compatible":True,"capabilities":{"safe_startup_v1":model.get("safe_capability",True)}}))
    raise SystemExit(0)
# Model the producer's initialized-thread requirement from the actual child argv.
# This is invocation evidence only, never native-session authority.
if args[:2] == ["agent", "start"] and model.get("safe_require_thread_profile"):
    import tomllib
    child = args[args.index("--") + 1:]
    config = {}
    for i, arg in enumerate(child[:-1]):
        if arg == "-c":
            config.update(tomllib.loads(child[i + 1]))
    profile = config.get("tui", {}).get("status_line")
    model.setdefault("safe_thread_profiles", {})[args[args.index("--pane") + 1]] = (
        profile == ["thread-id", "context-remaining"] and not model.get("safe_omit_thread_marker", False))
def startup_fixture(agent):
    if args[:2] != ["agent", "startup"]:
        if args[:2] == ["agent", "get"]:
            model.setdefault("safe_last_agents",{})[args[2]] = agent
        return
    request = json.loads(Path(args[4]).read_text())
    binding = request["binding"]
    agent = model.get("safe_last_agents",{}).get(request["target"],agent)
    review = request["target"].startswith("w-fixture:r")
    prefix = "review_" if review else ""
    if args[2] == "inspect":
        index = model.get(prefix + "safe_inspect_count", 0)
        sequence = model.get(prefix + "safe_inspections", [{}])
        step = sequence[min(index, len(sequence) - 1)]
        model[prefix + "safe_inspect_count"] = index + 1
        if index == model.get(prefix + "safe_inspect_denied_at", -1):
            path.write_text(json.dumps(model))
            print(json.dumps({"error":{"code":"permission_denied","message":"synthetic inspection refusal"}}))
            raise SystemExit(1)
        surface = step.get("surface", model.get(prefix + "safe_surface", "provider_prompt"))
        if model.get(prefix + "safe_unknown_while_working") and agent.get("agent_status") == "working":
            surface = "unknown"
        if model.get(prefix + "safe_skip") and not model.get(prefix + "safe_skipped"):
            surface = "update_skip_once"
        if model.get("safe_require_thread_profile") and not model.get("safe_thread_profiles", {}).get(request["target"], False):
            surface = "unknown"
        session = agent.get("agent_session")
        ready = agent.get("agent_status") in ("idle", "done") and not agent.get("launch_pending", False)
        if not session:
            ready = ready and agent.get("interactive_ready") is True
        result = {"startup":{"contract":"safe_startup_v1", "server_epoch":model.get(prefix + "safe_epoch","epoch-fixture"),
                  "managed_launch_id":"launch-" + request["target"], "evidence_token":"token-fixture", "binding":binding,
                  "observed_at_unix_ms":int(time.time()*1000), "surface":surface,
                  "prompt_ready":model.get(prefix + "safe_prompt_ready",ready) if surface == "provider_prompt" else False,
                  "launch_pending":agent.get("launch_pending",False), "native_session":session["value"] if session else None,
                  "menu":{"layout_id":"codex-0.164-picker-v1","selected_index":0,"skip_action":"skip_once"} if surface == "update_skip_once" else None,
                  "attempt":model.get(request.get("attempt_id"))}}
        result["startup"].update(step)
        result["startup"].update(model.get(prefix + "safe_override",{}))
        for field in model.get(prefix + "safe_drop",[]):
            result["startup"].pop(field,None)
    else:
        kind = request["kind"]
        if kind == "skip_once":
            model[prefix + "safe_skipped"] = not model.get(prefix + "safe_menu_remains",False)
            if model.get(prefix + "safe_epoch_after_skip"):
                model[prefix + "safe_epoch"] = model[prefix + "safe_epoch_after_skip"]
        elif review:
            model.setdefault("review_prompts", {})[request["target"]] = request["text"]
            if kind == "bootstrap":
                model["review_bootstrapped"] = True
        elif kind == "bootstrap":
            model["bootstrap_prompt_count"] = model.get("bootstrap_prompt_count",0) + 1
            model["bootstrap_started"] = True
        else:
            model["task_prompt_count"] = model.get("task_prompt_count",0) + 1
            model["agent_status"] = "working"
        submission = {"attempt_id":request["attempt_id"], "kind":kind, "binding":binding,
                      "payload_digest":"sha256:" + "a"*64, "outcome":model.get(prefix + "safe_outcome","submitted"),
                      "bytes_written":model.get(prefix + "safe_bytes","some"), "reason":"synthetic writer result"}
        model[request["attempt_id"]] = submission
        result = {"submission":submission}
    model.setdefault("safe_history", []).append({"operation": args[2], "target": request["target"],
        "surface": result.get("startup", {}).get("surface"), "kind": request.get("kind"),
        "deadline": request.get("deadline_unix_ms"), "at_ms": int(time.time()*1000)})
    path.write_text(json.dumps(model))
    if args[2] == "submit" and model.get(prefix + "safe_denied") == request["kind"]:
        print(json.dumps({"error":{"code":"permission_denied","message":"synthetic refusal"}}))
        raise SystemExit(1)
    if args[2] == "submit" and request["kind"] != "skip_once":
        if model.get("delay_command") == "agent prompt":
            time.sleep(model["delay_millis"] / 1000 if "delay_millis" in model else model.get("delay_seconds",2))
        if model.get("prompt_mode") == "lost_reply" or model.get("lost_command") == "agent prompt" or model.get("bootstrap_mode") == "lost_reply" and request["kind"] == "bootstrap":
            raise SystemExit(7)
    if args[2] == "submit" and (model.get(prefix + "safe_lost") == request["kind"] or review and model.get("review_prompt_lost")):
        raise SystemExit(7)
    print(json.dumps({"id":"fixture", "result":result}))
    raise SystemExit(0)
cmd = args[:2]
observation_delay = 0
if cmd == ["agent", "start"]:
    model["name"] = args[2]
agent = {"workspace_id":"w-fixture", "tab_id":"tab-fixture", "pane_id":"w-fixture:p1",
         "terminal_id":"terminal-fixture", "name":model.get("name", "unstarted"),
         "agent":model.get("provider", "codex"), "agent_status":model["agent_status"],
         "agent_session":{"agent":model.get("session_provider", model.get("provider", "codex")),
                          "kind":"id", "value":model["agent_session_id"]}}
if "foreground_cwd" in model:
    agent["foreground_cwd"] = model["foreground_cwd"]
bootstrap_mode = model.get("bootstrap_mode")
if bootstrap_mode:
    if "interactive_ready" in model:
        agent["interactive_ready"] = model["interactive_ready"]
    agent["launch_pending"] = model.get("launch_pending", False)
    if not model.get("bootstrap_started") or bootstrap_mode == "no_session":
        agent["agent_session"] = None
    elif cmd == ["agent", "get"] and not model.get("task_prompt_count"):
        count = model.get("bootstrap_gets", 0) + 1
        model["bootstrap_gets"] = count
        agent["agent_status"] = "working" if count < 2 or bootstrap_mode == "working_forever" else "idle"
        if bootstrap_mode == "session_change" and count >= 2:
            agent["agent_session"]["value"] = "replacement-session"
        if bootstrap_mode == "session_change_before_task" and count >= 3:
            agent["agent_session"]["value"] = "replacement-session"
        if agent["agent_status"] == "idle":
            model["bootstrap_settled_observed"] = True
if cmd == ["agent", "get"] and model.get("observations"):
    index = model.get("observation_index", 0)
    observation = model["observations"][min(index, len(model["observations"]) - 1)]
    observation_delay = observation.get("delay_millis", 0) / 1000
    agent.update({k: v for k, v in observation.items() if k not in ("drop_fields", "delay_millis")})
    for field in observation.get("drop_fields", []):
        agent.pop(field, None)
    model["observation_index"] = index + 1
startup_fixture(agent)
result = {}
if cmd == ["tab", "create"]:
    result = {"root_pane": agent}
elif cmd == ["pane", "run"]:
    if not model.get("claude_environment_no_ack"):
        subprocess.run(["/bin/bash", "--noprofile", "--norc", "-c", args[3]], check=True)
    result = {"ok": True}
elif cmd in (["agent", "start"], ["agent", "get"]):
    result = {"agent": agent}
elif cmd == ["agent", "prompt"]:
    if bootstrap_mode and args[3] == model["bootstrap_expected_prompt"]:
        model["bootstrap_prompt_count"] = model.get("bootstrap_prompt_count", 0) + 1
        model["bootstrap_started"] = True
        result = {"acknowledged": True}
    else:
        if bootstrap_mode:
            model["task_prompt_count"] = model.get("task_prompt_count", 0) + 1
            if not model.get("bootstrap_settled_observed") or not call.get("session_before_call"):
                model["premature_task_prompt"] = True
        model["agent_status"] = "working"
        result = {"acknowledged": True} if model["prompt_mode"] == "ack_without_agent" else {"agent":{**agent,"agent_status":"working"}}
elif cmd == ["tab", "rename"]:
    model["label"] = args[3]
    result = {"ok": True}
else:
    raise SystemExit("Unexpected stand-in command: " + repr(args))
# Polling deadlines may stop this process during publication. Keep the
# previous complete model until its replacement has been fully written.
pending = path.with_name(path.name + "." + str(os.getpid()) + ".tmp")
pending.write_text(json.dumps(model))
pending.replace(path)
time.sleep(observation_delay)
if bootstrap_mode == "lost_reply" and cmd == ["agent", "prompt"] and args[3] == model["bootstrap_expected_prompt"]:
    raise SystemExit(7)
failure = model.get("failure", {})
if failure.get("command") == " ".join(cmd):
    sys.stderr.write(failure.get("stderr", ""))
    sys.stdout.write(failure.get("stdout", ""))
    sys.stdout.flush()
    sys.stderr.flush()
    time.sleep(failure.get("delay_seconds", 0))
    raise SystemExit(failure.get("exit", 0))
if model.get("delay_command") == " ".join(cmd):
    time.sleep(model["delay_millis"] / 1000 if "delay_millis" in model else model.get("delay_seconds", 2))
if cmd == ["agent", "prompt"] and model["prompt_mode"] == "lost_reply":
    raise SystemExit(7)
if model.get("lost_command") == " ".join(cmd):
    raise SystemExit(7)
print(json.dumps({"id":"fixture", "result":result}))
`

type workflowFixture struct {
	D          Dependencies
	R          WorkflowRequest
	File       string
	Acceptance Acceptance
	Report     WorkflowReport
	Model      string
	Calls      string
	Handoff    string
}

func newWorkflowFixture(t *testing.T, root, ply string) workflowFixture {
	return newWorkflowFixtureSeed(t, root, ply, "", "")
}
func newWorkflowFixtureSeed(t *testing.T, root, ply, seedSource, seedOID string) workflowFixture {
	t.Helper()
	d, native, _ := fixtureSeed(t, seedSource, seedOID, root)
	// This fixture also supports a real local candidate commit and a dirty round.
	var draftFields map[string]any
	if e := json.Unmarshal(native.HandoffDraft, &draftFields); e != nil {
		t.Fatal(e)
	}
	authority := draftFields["authority"].(map[string]any)
	effects := authority["allowed_effects"].([]any)
	bind := draftFields["binding_request"].(map[string]any)
	for _, item := range []struct{ kind, operation string }{{"git_index_write", "update_index"}, {"git_commit", "create_commit"}} {
		effects = append(effects, map[string]any{"id": item.kind, "type": item.kind, "scope": map[string]any{"kind": "git", "ref": bind["target_ref"], "expected_before_oid": bind["expected_oid"], "operation": item.operation}, "max_occurrences": 3, "sequence": len(effects) + 1})
	}
	authority["allowed_effects"] = effects
	var e error
	native.HandoffDraft, e = Canonical(draftFields)
	if e != nil {
		t.Fatal(e)
	}
	if ply != "" {
		native.Runtime.PlyExecutable = Executable{ply, hashFileTest(t, ply)}
	}
	python := os.Getenv("PLY_WORKFLOW_PYTHON")
	if python == "" {
		python, _ = exec.LookPath("python3")
	}
	python, _ = filepath.EvalSymlinks(python)
	herdr := filepath.Join(root, "herdr-standin")
	if e := os.WriteFile(herdr, []byte("#!"+python+"\n"+workflowStandin), 0700); e != nil {
		t.Fatal(e)
	}
	r := WorkflowRequest{StartupProtocol: ReportLegacyAcceptanceV1, Envelope: workflowEnv("herdr-run-request"), RequestKey: native.RequestKey, WorkspaceRoot: root, PreparationID: native.PreparationID, PreparationSHA256: native.PreparationSHA256, HandoffDraft: native.HandoffDraft, Runtime: native.Runtime, Agreement: native.Agreement, HumanAuthority: HumanAuthority{"fixture human", "human_authorized_herdr", true}, ReturnMode: "reviewed_report_only"}
	r.Herdr.Executable = Executable{herdr, hashFileTest(t, herdr)}
	r.Herdr.WorkspaceID = "w-fixture"
	r.Herdr.TabLabel = "Fixture Task"
	r.Coordinator.ActorClaim = "fixture coordinator"
	r.Coordinator.MayRequestChanges = true
	f := writeAny(t, root, "workflow-request.json", r)
	if _, e := ReadWorkflowRequest(f); e != nil {
		t.Fatal(e)
	}
	obs, e := workflowhandoff.ObserveTaskRun(d.Workflow, root, r.PreparationID, r.HandoffDraft)
	if e != nil {
		t.Fatal(e)
	}
	// Native idempotent creation establishes fixture-only callback templates. The
	// production start reuses this exact unused publication; no TaskRun is made.
	draft := writeAny(t, root, "fixture-handoff.json", json.RawMessage(r.HandoffDraft))
	h, e := workflowhandoff.Create(d.Workflow, workflowhandoff.CreateInput{DraftPath: draft})
	if e != nil {
		t.Fatal(e)
	}
	view, e := workflowhandoff.TaskRunHandoffView(d.Workflow, h.Locator)
	if e != nil {
		t.Fatal(e)
	}
	var mandate map[string]json.RawMessage
	if e = json.Unmarshal(view, &mandate); e != nil {
		t.Fatal(e)
	}
	var cap struct {
		ReplyRoot string `json:"reply_root"`
	}
	json.Unmarshal(mandate["reply_capability"], &cap)
	staging := filepath.Join(cap.ReplyRoot, "staging")
	if e = privateDir(staging); e != nil {
		t.Fatal(e)
	}
	tmp := filepath.Join(root, "recipient-tmp")
	if e = privateDir(tmp); e != nil {
		t.Fatal(e)
	}
	roots := []string{obs.Target.WorktreeLocator, cap.ReplyRoot}
	sort.Strings(roots)
	sandbox, _ := Canonical(map[string]any{"read_roots": []string{root}, "write_roots": roots, "temp_root": tmp, "matches_contract": true})
	a := Acceptance{workflowEnv("run-acceptance"), workflowID(r), digest(r), "ply:" + workflowID(r), RuntimeClaim{ptr("codex"), nil, ptr(r.Runtime.PermissionBinding.ProfileID), ptr(r.Runtime.PermissionBinding.EffectivePolicySHA256), ptr("fixture-session")}, sandbox, "started", json.RawMessage("[]")}
	p := obs.Preparation.Plan
	requirements := map[string]any{"kind": "WorkspaceTaskRequirementEvidence@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "task_id": p.TaskID, "spec_id": p.SpecID, "spec": p.Spec, "result_oid": p.ParentOID, "result_tree": p.ParentTree, "requirements": []any{map[string]any{"id": "f-01", "outcome": "passed", "verifier_ids": []string{"test"}, "artifact_ids": []string{}, "reason": "Synthetic verifier evidence only."}}}
	artifact := writeAny(t, staging, "requirements-0.json", requirements)
	info, _ := os.Stat(artifact)
	raw := func(v any) json.RawMessage { b, _ := Canonical(v); return b }
	report := WorkflowReport{Report: Report{Envelope: workflowEnv("round-report"), RunID: workflowID(r), RequestSHA256: digest(r), SessionID: "ply:" + workflowID(r), Outcome: "complete", Summary: "Synthetic Task report.", Meaning: "Stand-in evidence; no actual provider or human QA was run.", BudgetUsage: &BudgetUsage{true, 0, 10, 0}, StopReasons: raw([]any{}), ObservedEffects: syntheticEffects(t, mandate["authority"]), VerifierResults: raw([]any{map[string]any{"verifier_id": "test", "argv": []string{"go", "test", "./..."}, "cwd": p.WorktreePath, "exit": 0, "bound_oid_or_sha256": p.ParentOID, "stdout_artifact_id": nil, "stderr_artifact_id": nil}}), Review: raw(map[string]any{"findings": []any{}, "fixes": []any{}, "open_actionable_findings": []any{}}), Artifacts: raw([]any{map[string]any{"artifact_id": "task-requirements", "kind": "managed", "description": "Fixture requirements", "media_type": "application/json", "classification": "workspace_internal", "size_bytes": info.Size(), "sha256": hashFileTest(t, artifact), "locator": artifact}}), EvidenceGaps: raw([]any{}), ForbiddenEffectsObserved: raw([]any{}), TechnicalAssessment: TechnicalAssessment{"passed", []string{"test"}, []workspace.TaskAcceptedDebtRecord{}}}}
	model := writeAny(t, root, "herdr-model.json", map[string]any{"agent_status": "idle", "agent_session_id": "fixture-session", "prompt_mode": "normal", "foreground_cwd": obs.Target.WorktreeLocator})
	d.Executable = func() (string, error) { return r.Runtime.PlyExecutable.Path, nil }
	d.CWD = func() (string, error) { return obs.Target.WorktreeLocator, nil }
	return workflowFixture{d, r, f, a, report, model, filepath.Join(root, "herdr-calls.jsonl"), h.Locator}
}
func TestWorkflowFixtureExport(t *testing.T) {
	root := os.Getenv("PLY_WORKFLOW_FIXTURE_ROOT")
	if root == "" {
		t.Skip("fixture export is explicitly requested by the black-box harness")
	}
	fx := newWorkflowFixture(t, root, os.Getenv("PLY_WORKFLOW_FIXTURE_BINARY"))
	manifest := map[string]any{"request_file": fx.File, "worktree": fx.D.CWD, "model_file": fx.Model, "calls_file": fx.Calls, "acceptance_template": fx.Acceptance, "report_template": fx.Report, "protected_paths": []string{filepath.Join(root, ".ply", "workspace.yaml"), filepath.Join(root, ".ply", "projects.yaml"), filepath.Join(root, ".ply", "work-items.yaml")}}
	cwd, _ := fx.D.CWD()
	manifest["worktree"] = cwd
	b, e := json.Marshal(manifest)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println("WH01_FIXTURE:" + string(b))
}
