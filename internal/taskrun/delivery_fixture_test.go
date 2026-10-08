package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type deliveryFixture struct {
	workflowFixture
	Prepared       workspace.TaskGoalPreparedExecution
	Parent         string
	AcceptancePath string
}

type deliveryFixtureOptions struct {
	Agreement *workspace.DeliveryAgreement
	ParentRef string
	Root      string
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
	root, epic, _, oid := prepareServiceWorkspace(t, roots...)
	parent := option.ParentRef
	if parent == "" {
		parent = "epic"
	}
	runGit(t, epic, "branch", "-m", parent)
	w := workspace.SystemDependencies()
	if _, e := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Delivery fixture", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/" + parent, ExpectedOID: oid}); e != nil {
		t.Fatal(e)
	}
	if _, e := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "task", Title: "Delivery fixture", Description: "A synthetic goal tests delivery semantics; no native provider or human approval.", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	r, e := w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	problem := workspace.TaskRevisionRef{Revision: r.TaskProblemRevisions[0].Revision, ManifestSHA256: r.TaskProblemRevisions[0].ManifestSHA256}
	doc := filepath.Join(root, "goal.md")
	contents := []byte("# Synthetic goal\nAdd the fixture behavior; own meaningful acceptance. Root owns any notification.\n")
	if e = os.WriteFile(doc, contents, 0600); e != nil {
		t.Fatal(e)
	}
	input := map[string]any{
		"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "format": "json", "format_version": 1, "canonicalization": "RFC8785",
		"publication_key": "fixture/goal", "task_id": "task", "spec_id": "solution", "registry_upgrade": nil, "expected_previous": nil, "problem": problem,
		"recorder":      map[string]any{"actor_claim": "fixture planner", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"},
		"contract_kind": "goal", "title": "Delivery fixture", "objective": "Observe a complete delivery without attributing synthetic evidence to a real human.", "change_reason": "Synthetic fixture goal",
		"documents": []any{map[string]any{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(contents), "size_bytes": len(contents), "media_type": "text/markdown", "git_provenance": nil}}},
		"design":    []any{map[string]any{"document_id": "design", "section": nil}}, "constraints": []string{"Local isolated fixture only."},
		"requirements": []any{map[string]any{"id": "behavior", "acceptance": "The actual acceptance entrypoint observes the fixture behavior."}},
		"executor":     map[string]any{"provider": provider, "model": "fixture-model", "effort": "medium"},
	}
	if option.Agreement != nil {
		a := *option.Agreement
		if a.SchemaVersion == 0 {
			a.SchemaVersion = 1
		}
		a.ProjectID, a.RepoID, a.EpicID = "ply", "ply", "epic"
		if a.TargetRef == "" {
			a.TargetRef = "refs/heads/" + parent
		}
		if a.Mode != workspace.DeliveryPullRequest {
			a.TargetWorktree = epic
		}
		input["delivery"] = a
	}
	publication, e := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "task", File: writeAny(t, root, "goal-draft.json", input)})
	if e != nil {
		t.Fatal(e)
	}
	goal := workspace.TaskGoalRef{SpecID: "solution", Spec: workspace.TaskRevisionRef{Revision: *publication.OutcomeRef.Revision, ManifestSHA256: publication.OutcomeRef.ManifestSHA256}}
	r, e = w.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	queue := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/goals", ProjectID: "ply", RepoID: "ply", EpicID: "epic", Entries: []workspace.TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}, RegistryUpgrade: &workspace.TaskRegistryUpgrade{RegistrySHA256: r.RawSHA256, FromVersion: r.FormatVersion}}
	if _, e = workspace.SetTaskQueue(w, writeAny(t, root, "goal-queue.json", queue)); e != nil {
		t.Fatal(e)
	}
	acceptance := filepath.Join(root, "acceptance.sh")
	preview, e := workspace.PreviewTaskGoalExecution(w, workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, Next: true, AcceptancePath: acceptance})
	if e != nil {
		t.Fatal(e)
	}
	p, e := workspace.PrepareTaskGoalExecution(w, preview.Input, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "synthetic human", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Synthetic execution choice; not actual product QA."})
	if e != nil {
		t.Fatal(e)
	}
	d := SystemDependencies(w)
	preparation := *p.Preparation.Preparation
	draft, e := workflowhandoff.BuildDeliveryHandoffDraft(d.Workflow, preparation.ID, "synthetic human", "fixture owner", "fixture/delivery", acceptance, workflowhandoff.DeliveryDraftAuthority{AllowLocalInstall: true, Agreement: p.Delivery})
	if e != nil {
		t.Fatal(e)
	}
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
	boundary := "after_human_pass"
	if p.Delivery != nil && (p.Delivery.Mode == workspace.DeliveryPullRequest || p.Delivery.HumanOwnedIntegration()) {
		boundary = "none"
	}
	req, e := BuildDeliveryWorkflowRequest(DeliveryRequestInput{RequestKey: "fixture/delivery", WorkspaceRoot: root, PreparationID: preparation.ID, PreparationSHA256: digest(preparation), HandoffDraft: draft, Runtime: runtime, HumanAuthority: HumanAuthority{"synthetic human", "human_authorized_herdr", true}, HerdrExecutable: Executable{herdr, hashFileTest(t, herdr)}, HerdrWorkspaceID: "w-fixture", TabLabel: "Delivery fixture", Delivery: DeliveryContract{OwnerClaim: "fixture owner", Goal: FileBinding{doc, hashFileTest(t, doc)}, AcceptancePath: acceptance, AllowSubagents: true, AllowLocalInstall: true, LocalIntegration: boundary, ReasoningEffort: "medium", Agreement: p.Delivery}})
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
