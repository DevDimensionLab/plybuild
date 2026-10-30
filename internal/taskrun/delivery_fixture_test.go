package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type deliveryFixture struct {
	workflowFixture
	Prepared       deliveryPreparedFixture
	Parent         string
	AcceptancePath string
}

type deliveryFixtureOptions struct {
	StartFixture      func(*testing.T, workflowFixture) WorkflowRun
	SeedSource        string
	SeedOID           string
	Parallel          bool
	AdoptExisting     bool
	Agreement         *workspace.DeliveryAgreement
	AllowLocalInstall bool
	ParentRef         string
	Root              string
	BeforePrepare     func(workspace.Dependencies, string, string, workspace.TaskGoalRef)
	SelectFreshGoal   func(workspace.Dependencies, string, string, workspace.TaskGoalRef) workspace.TaskGoalRef
	InformativeReport func(string) []byte
}

type deliveryPreparedFixture struct {
	Preparation workspace.TaskPreparationReadback
	Goal        workspace.TaskGoalRef
	Delivery    *workspace.DeliveryAgreement
}

func newDeliveryFixture(t *testing.T, provider string, options ...deliveryFixtureOptions) deliveryFixture {
	t.Helper()
	var option deliveryFixtureOptions
	if len(options) > 0 {
		option = options[0]
	}
	var roots []string
	if option.Root != "" {
		roots = []string{option.Root}
	}
	d, native, _ := fixtureSeed(t, option.SeedSource, option.SeedOID, roots...)
	root := native.WorkspaceRoot
	w := d.Workspace
	registry, e := w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	var preparation workspace.TaskPreparation
	for _, item := range registry.TaskPreparations {
		if item.ID == native.PreparationID {
			preparation = item
		}
	}
	if preparation.ID == "" {
		t.Fatal("manual fixture preparation missing")
	}
	epic := preparation.Plan.ObservedParent.WorktreeLocator
	acceptance := filepath.Join(root, "acceptance.sh")
	doc := filepath.Join(root, "goal.md")
	if e = os.WriteFile(doc, []byte("# Explicit fixture instruction\n\nChoose subsequent work separately.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	goal := workspace.TaskGoalRef{SpecID: preparation.Plan.SpecID, Spec: preparation.Plan.Spec}
	var fields map[string]any
	if e = json.Unmarshal(native.HandoffDraft, &fields); e != nil {
		t.Fatal(e)
	}
	fields["schema_version"] = 3
	fields["delivery_binding"] = map[string]any{"mode": "owner", "human_actor": "fixture human", "owner_claim": "fixture owner", "parent_handoff_locator": nil, "parent_handoff_sha256": nil, "candidate_key": nil}
	fields["budget"] = map[string]any{"max_rounds": nil, "round_definition": map[string]any{"unit": "candidate_verification_attempt", "command_retry_consumes_round": false, "retry_condition": "only_if_no_effect_started"}}
	fields["verifiers"] = []any{map[string]any{"id": "test", "argv": []string{"/bin/sh", acceptance}, "cwd": preparation.Plan.WorktreePath, "env": []string{}, "expected_exit": 0, "stop_on_failure": true, "evidence": map[string]any{"capture_stdout": true, "capture_stderr": true, "classification": "workspace_internal", "binding": "target_oid"}}}
	authority := fields["authority"].(map[string]any)
	authority["human_gates"] = []string{}
	effects := []any{}
	add := func(id, kind string, scope map[string]any) {
		var limit any
		if scope["kind"] == "delivery" {
			limit = 1
		}
		effects = append(effects, map[string]any{"id": id, "type": kind, "scope": scope, "max_occurrences": limit, "sequence": len(effects) + 1})
	}
	add("task-write", "filesystem_write", map[string]any{"kind": "task_worktree"})
	for _, pair := range [][2]string{{"git_index_write", "update_index"}, {"git_commit", "create_commit"}, {"git_ref_write", "update_ref"}} {
		add(pair[0], pair[0], map[string]any{"kind": "git", "ref": "refs/heads/" + preparation.Plan.Branch, "expected_before_oid": preparation.Plan.SourceStartOID(), "operation": pair[1]})
	}
	add("local-commands", "command_execute", map[string]any{"kind": "command", "procedure_ids": []string{"implement"}, "verifier_ids": []string{"test"}})
	add("process-start", "process_start", map[string]any{"kind": "boolean", "allowed": true})
	if option.AllowLocalInstall {
		add("install", "install", map[string]any{"kind": "boolean", "allowed": true})
	}
	agreement := option.Agreement
	if agreement != nil {
		copy := *agreement
		agreement = &copy
		agreement.ProjectID = "ply"
		agreement.RepoID = "ply"
		agreement.EpicID = "epic"
		agreement.SourceRef = "refs/heads/" + preparation.Plan.Branch
		if agreement.TargetRef == "" {
			agreement.TargetRef = preparation.Plan.ObservedParent.Ref
		}
		if agreement.Mode != workspace.DeliveryPullRequest {
			agreement.TargetWorktree = epic
		}
		fields["delivery_binding"].(map[string]any)["agreement"] = agreement
		scope := map[string]any{"kind": "delivery", "agreement": agreement}
		if agreement.Mode == workspace.DeliveryPullRequest {
			add("delivery-push", "push", scope)
			add("delivery-pr", "pull_request", scope)
		} else if !agreement.HumanOwnedIntegration() {
			add("delivery-integration", "integration", scope)
		}
	}
	authority["allowed_effects"] = effects
	authority["forbidden_effects"] = []any{}
	draft, e := json.Marshal(fields)
	if e != nil {
		t.Fatal(e)
	}
	p := deliveryPreparedFixture{Preparation: workspace.TaskPreparationReadback{Preparation: &preparation}, Goal: goal, Delivery: agreement}
	providerPath, ply := filepath.Join(root, "synthetic-provider"), filepath.Join(root, "synthetic-ply")
	for _, path := range []string{providerPath, ply} {
		if e = os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
			t.Fatal(e)
		}
	}
	python, e := exec.LookPath("python3")
	if e != nil {
		t.Fatal(e)
	}
	python, e = filepath.EvalSymlinks(python)
	if e != nil {
		t.Fatal(e)
	}
	herdr := filepath.Join(root, "herdr-standin")
	if e = os.WriteFile(herdr, []byte("#!"+python+"\n"+workflowStandin), 0700); e != nil {
		t.Fatal(e)
	}
	policy := filepath.Join(root, "launch-policy.json")
	if e = os.WriteFile(policy, []byte("Requested fixture launch policy; no OS attestation.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	profile := "synthetic-policy"
	if provider == "claude" {
		profile = "manual"
	}
	runtime := Runtime{Provider: provider, Mode: "interactive", Model: "fixture-model", Executable: Executable{providerPath, hashFileTest(t, providerPath)}, PlyExecutable: Executable{ply, hashFileTest(t, ply)}, PermissionBinding: Permission{AuthorityKind: "launch_contract_pending_runtime_acceptance", ProfileID: profile, EffectivePolicySHA256: hashFileTest(t, policy), Evidence: []Evidence{{Locator: policy, SHA256: hashFileTest(t, policy), Role: "runtime_contract"}}}}
	boundary := DeliveryLocalIntegrationBoundary(p.Delivery)
	goalBinding := FileBinding{doc, hashFileTest(t, doc)}
	if p.Delivery != nil && p.Delivery.AutomaticAcceptance() {
		goalBinding = FileBinding{filepath.Join(root, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(goal.Spec.ManifestSHA256, "sha256:")+".json"), goal.Spec.ManifestSHA256}
	}
	req, e := BuildDeliveryWorkflowRequest(DeliveryRequestInput{RequestKey: "fixture/delivery", WorkspaceRoot: root, PreparationID: preparation.ID, PreparationSHA256: digest(preparation), HandoffDraft: draft, Runtime: runtime, HumanAuthority: HumanAuthority{"synthetic human", "human_authorized_herdr", true}, HerdrExecutable: Executable{herdr, hashFileTest(t, herdr)}, HerdrWorkspaceID: "w-fixture", TabLabel: "Delivery fixture", Delivery: DeliveryContract{AdmissionMode: preparation.Plan.AdmissionMode, SourceReconciliation: preparation.Plan.SourceReconciliation, OwnerClaim: "fixture owner", Goal: goalBinding, AcceptancePath: acceptance, AllowSubagents: true, AllowLocalInstall: option.AllowLocalInstall, LocalIntegration: boundary, ReasoningEffort: "medium", Agreement: p.Delivery}})
	if e != nil {
		t.Fatal(e)
	}
	model := writeAny(t, root, "herdr-model.json", map[string]any{"agent_status": "idle", "agent_session_id": "fixture-session", "prompt_mode": "normal", "provider": provider, "foreground_cwd": preparation.Plan.WorktreePath})
	bin := filepath.Join(root, "bin")
	if e = os.Mkdir(bin, 0700); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(providerPath, filepath.Join(bin, provider)); e != nil {
		t.Fatal(e)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HERDR_ENV", "1")
	d.CWD = func() (string, error) { return preparation.Plan.WorktreePath, nil }
	d.Executable = func() (string, error) { return ply, nil }
	return deliveryFixture{workflowFixture: workflowFixture{D: d, R: req, File: writeAny(t, root, "delivery-request.json", req), Model: model, Calls: filepath.Join(root, "herdr-calls.jsonl")}, Prepared: p, Parent: epic, AcceptancePath: acceptance}
}

func deliveryTestAcceptance(t *testing.T, f deliveryFixture, o WorkflowRun) DeliveryAcceptance {
	t.Helper()
	policy := filepath.Join(f.R.WorkspaceRoot, "actual-policy.json")
	if e := os.WriteFile(policy, []byte("Actual synthetic fixture authority observation, distinct from launch contract.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	view, e := workflowhandoff.TaskRunHandoffView(f.D.Workflow, o.Handoff.Locator)
	if e != nil {
		t.Fatal(e)
	}
	var mandate struct {
		ReplyCapability struct {
			ReplyRoot string `json:"reply_root"`
		} `json:"reply_capability"`
	}
	if e = json.Unmarshal(view, &mandate); e != nil {
		t.Fatal(e)
	}
	// This isolated fixture's actual writable area is the whole private
	// workspace, including future candidate-control reply roots.
	roots := []string{f.R.WorkspaceRoot}
	sort.Strings(roots)
	tmp := filepath.Join(f.R.WorkspaceRoot, "recipient-tmp")
	if e = privateDir(tmp); e != nil {
		t.Fatal(e)
	}
	sandbox, _ := Canonical(map[string]any{"read_roots": []string{f.R.WorkspaceRoot}, "write_roots": roots, "temp_root": tmp, "matches_contract": true})
	a := Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, RuntimeClaim: RuntimeClaim{RuntimeID: ptr(f.R.Runtime.Provider), ProfileID: ptr(f.R.Runtime.PermissionBinding.ProfileID), EffectivePolicySHA256: ptr(hashFileTest(t, policy)), NativeSessionID: ptr(o.Transport.AgentSessionID)}, Sandbox: sandbox, Acceptance: "started", Issues: json.RawMessage("[]")}
	permission := DeliveryPermissionAcceptance{LaunchContractSHA256: f.R.Runtime.PermissionBinding.EffectivePolicySHA256, PermissionConfirmed: true, ActualPolicyEvidence: []Evidence{{Locator: policy, SHA256: hashFileTest(t, policy), Role: "effective_policy"}}}
	if f.R.Delivery.Agreement != nil {
		permission.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(*f.R.Delivery.Agreement)
		permission.AllowedEffects = workflowhandoff.DeliveryAllowedEffects(*f.R.Delivery.Agreement)
	}
	return DeliveryAcceptance{a, permission}
}

func deliveryTestAccept(t *testing.T, f deliveryFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	a := deliveryTestAcceptance(t, f, o)
	o, e := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "recipient-acceptance.json", a))
	if e != nil {
		t.Fatal(e)
	}
	return o
}
