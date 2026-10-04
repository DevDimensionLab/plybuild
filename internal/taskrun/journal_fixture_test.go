package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
)

// This named synthetic fixture uses the native public validators. No real provider runs.
func prepareJournalWorkspace(t *testing.T, roots ...string) (string, string, string, string) {
	t.Helper()
	var requested string
	if len(roots) == 0 {
		requested = t.TempDir()
	} else {
		requested = roots[0]
		if err := os.MkdirAll(requested, 0700); err != nil {
			t.Fatal(err)
		}
	}
	root, err := filepath.EvalSymlinks(requested)
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	runGit(t, target, "init", "-b", "main")
	runGit(t, target, "config", "user.name", "WF01 Test")
	runGit(t, target, "config", "user.email", "wf01@example.invalid")
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, target, "add", "README.md")
	runGit(t, target, "commit", "-m", "base")
	ref := gitOutput(t, target, "symbolic-ref", "HEAD")
	oid := gitOutput(t, target, "rev-parse", "HEAD")

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	workspaceDependencies := workspace.SystemDependencies()
	if _, err := workspace.Init(workspaceDependencies); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AddProject(workspaceDependencies, workspace.ProjectAddInput{ProjectID: "ply-journal-fixture", Name: "Ply", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "fixture-repo", Path: target}}}); err != nil {
		t.Fatal(err)
	}
	return root, target, ref, oid
}

func journalHandoffDraft(target, ref, oid string) canonicaljson.Object {
	requiredStart := []canonicaljson.Value{"acceptance", "binding", "contract_digests", "issues", "observed_inputs", "observed_project", "observed_target", "observed_workspace", "principal", "receipt_id", "sandbox"}
	requiredTerminal := []canonicaljson.Value{"artifacts", "binding", "evidence_gaps", "final_target", "forbidden_effects_observed", "meaning", "observed_effects", "principal", "reported_outcome", "result_id", "review", "rounds_used", "start_binding", "stop_reasons", "summary", "verifier_results"}
	stop := func(kind string) canonicaljson.Value {
		return canonicaljson.Object{{Name: "type", Value: kind}, {Name: "description", Value: "Stop safely when this condition occurs."}}
	}
	return canonicaljson.Object{
		{Name: "kind", Value: "ply.workflow.handoff-draft"},
		{Name: "schema_version", Value: int64(1)},
		{Name: "format", Value: "json"},
		{Name: "format_version", Value: int64(1)},
		{Name: "canonicalization", Value: "RFC8785"},
		{Name: "publication_key", Value: "wf01/service-test"},
		{Name: "activity_key", Value: "wf01/service-test"},
		{Name: "goal", Value: canonicaljson.Object{{Name: "title", Value: "Test handoff"}, {Name: "recipient_role", Value: "Delivery agent"}, {Name: "objective", Value: "Make the bounded test change."}, {Name: "done_when", Value: "The verifier is green."}}},
		{Name: "recipient", Value: canonicaljson.Object{{Name: "principal_id", Value: "codex-delivery-agent"}, {Name: "principal_kind", Value: "human_started_agent"}, {Name: "runtime_constraints", Value: []canonicaljson.Value{"local"}}}},
		{Name: "binding_request", Value: canonicaljson.Object{{Name: "project_id", Value: "ply-journal-fixture"}, {Name: "repo_id", Value: "fixture-repo"}, {Name: "target_worktree", Value: target}, {Name: "target_ref", Value: ref}, {Name: "expected_oid", Value: oid}, {Name: "status_policy", Value: canonicaljson.Object{{Name: "mode", Value: "clean"}}}}},
		{Name: "inputs", Value: []canonicaljson.Value{}},
		{Name: "authority", Value: canonicaljson.Object{
			{Name: "allowed_effects", Value: []canonicaljson.Value{
				canonicaljson.Object{{Name: "id", Value: "write-readme"}, {Name: "type", Value: "filesystem_write"}, {Name: "scope", Value: canonicaljson.Object{{Name: "kind", Value: "filesystem"}, {Name: "paths", Value: []canonicaljson.Value{"README.md"}}, {Name: "directory_prefixes", Value: []canonicaljson.Value{}}}}, {Name: "max_occurrences", Value: int64(1)}, {Name: "sequence", Value: int64(1)}},
				canonicaljson.Object{{Name: "id", Value: "run-test"}, {Name: "type", Value: "command_execute"}, {Name: "scope", Value: canonicaljson.Object{{Name: "kind", Value: "command"}, {Name: "procedure_ids", Value: []canonicaljson.Value{"implement"}}, {Name: "verifier_ids", Value: []canonicaljson.Value{"test"}}}}, {Name: "max_occurrences", Value: int64(1)}, {Name: "sequence", Value: int64(2)}},
			}},
			{Name: "forbidden_effects", Value: []canonicaljson.Value{
				canonicaljson.Object{{Name: "type", Value: "network"}, {Name: "reason", Value: "Network is outside this local task."}},
				canonicaljson.Object{{Name: "type", Value: "push"}, {Name: "reason", Value: "Push requires a later human gate."}},
			}},
			{Name: "human_gates", Value: []canonicaljson.Value{"human_task_qa", "local_integration"}},
		}},
		{Name: "budget", Value: canonicaljson.Object{{Name: "max_rounds", Value: int64(2)}, {Name: "round_definition", Value: canonicaljson.Object{{Name: "unit", Value: "implementation_or_review_fix_iteration"}, {Name: "command_retry_consumes_round", Value: false}, {Name: "retry_condition", Value: "only_if_no_effect_started"}}}}},
		{Name: "procedure", Value: []canonicaljson.Value{canonicaljson.Object{{Name: "id", Value: "implement"}, {Name: "instruction", Value: "Implement the bounded change and report the task-requirements artifact."}, {Name: "required_before", Value: []canonicaljson.Value{}}}}},
		{Name: "verifiers", Value: []canonicaljson.Value{
			canonicaljson.Object{
				{Name: "id", Value: "test"},
				{Name: "argv", Value: []canonicaljson.Value{"go", "test", "./..."}},
				{Name: "cwd", Value: target},
				{Name: "env", Value: []canonicaljson.Value{}},
				{Name: "expected_exit", Value: int64(0)},
				{Name: "stop_on_failure", Value: true},
				{Name: "evidence", Value: canonicaljson.Object{{Name: "capture_stdout", Value: false}, {Name: "capture_stderr", Value: false}, {Name: "classification", Value: "workspace_internal"}, {Name: "binding", Value: "target_oid"}}},
			},
		}},
		{Name: "stop_conditions", Value: []canonicaljson.Value{stop("product_decision_required"), stop("scope_or_authority_expansion"), stop("target_or_input_drift"), stop("unknown_or_partial_effect"), stop("unexpected_sensitive_data"), stop("round_budget_exhausted")}},
		{Name: "reporting", Value: canonicaljson.Object{{Name: "summary_max_codepoints", Value: int64(240)}, {Name: "meaning_max_codepoints", Value: int64(600)}, {Name: "required_start_fields", Value: requiredStart}, {Name: "required_terminal_fields", Value: requiredTerminal}}},
	}
}

func journalSpecDraft(t *testing.T, root, kind, key string, fields map[string]any) workspace.TaskContentInput {
	t.Helper()
	fields["kind"], fields["schema_version"], fields["format"], fields["format_version"], fields["canonicalization"] = kind, 1, "json", 1, "RFC8785"
	fields["publication_key"], fields["task_id"], fields["registry_upgrade"] = key, "journal-acceptance", nil
	fields["recorder"] = map[string]any{"actor_claim": "fixture agent", "control_surface": "isolated service test", "recorded_at_utc": "2026-09-29T12:00:00Z"}
	return workspace.TaskContentInput{TaskID: "journal-acceptance", File: writeCanonicalTestFile(t, root, filepath.Base(key)+".json", taskSpecFixtureValue(t, fields).(canonicaljson.Object))}
}

func journalFixture(t *testing.T, roots ...string) (Dependencies, Request, string) {
	t.Helper()
	root, epic, ref, oid := prepareJournalWorkspace(t, roots...)
	w := workspace.SystemDependencies()
	if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "journal-fixture-epic", Title: "Fixture Epic", ProjectID: "ply-journal-fixture", RepoID: "fixture-repo", Worktree: epic, Ref: ref, ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "journal-acceptance", Title: "Fixture Task", Description: "Fixture need", ParentEpicID: "journal-fixture-epic", ProjectID: "ply-journal-fixture", RepoID: "fixture-repo"}); err != nil {
		t.Fatal(err)
	}
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	parent := r.Epics[0].RepoBindings[0].Worktree
	problem := workspace.TaskRevisionRef{Revision: r.TaskProblemRevisions[0].Revision, ManifestSHA256: r.TaskProblemRevisions[0].ManifestSHA256}
	doc := filepath.Join(root, "solution.md")
	b := []byte("# Fixture solution\n\nScope README.md; verify test; preserve snapshots on failure.\n")
	if err = os.WriteFile(doc, b, 0600); err != nil {
		t.Fatal(err)
	}
	refs := []any{map[string]any{"document_id": "solution", "section": nil}}
	part := map[string]any{"state": "present", "reason": nil, "documents": refs}
	spec, err := workspace.RecordTaskSpec(w, journalSpecDraft(t, root, "WorkspaceTaskSpecDraft@1", "fixture/spec", map[string]any{
		"spec_id": "solution", "expected_previous": nil, "problem": problem, "title": "Fixture solution",
		"parts":      map[string]any{"abstract": part, "functional": part, "technical": part},
		"documents":  []any{map[string]any{"id": "solution", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(b), "size_bytes": len(b), "media_type": "text/markdown", "git_provenance": nil}}},
		"supporting": []any{}, "requirements": []any{map[string]any{"id": "f-01", "functional_refs": refs, "technical_refs": refs, "acceptance": "Fixture behavior preserved.", "verification_ids": []string{"test"}}},
		"removed_requirement_ids": []any{}, "phases": []any{map[string]any{"id": "implement", "purpose": "Fixture implementation", "requirement_ids": []string{"f-01"}, "entry_criteria": []any{}, "exit_criteria": []string{"Verified"}, "verification_ids": []string{"test"}}},
		"implementation_basis": map[string]any{"project_id": "ply-journal-fixture", "repo_id": "fixture-repo", "git_common_dir": r.Tasks[0].GitCommonDir, "epic_id": "journal-fixture-epic", "parent_worktree_id": parent.ID, "parent_ref": parent.Ref, "parent_oid": parent.OID, "parent_tree": parent.Tree, "start_oid": parent.OID, "start_tree": parent.Tree}, "dependencies": []any{}, "change_reason": "Fixture solution"}))
	if err != nil {
		t.Fatal(err)
	}
	sr := workspace.TaskRevisionRef{Revision: *spec.OutcomeRef.Revision, ManifestSHA256: spec.OutcomeRef.ManifestSHA256}
	checks := []any{}
	for _, id := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, map[string]any{"id": id, "outcome": "pass", "reason": "Fixture assessment only", "evidence_document_ids": []any{}})
	}
	asm, err := workspace.AssessTaskSpec(w, journalSpecDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", "fixture/assess", map[string]any{"spec_id": "solution", "spec": sr, "expected_previous_assessment": nil, "outcome": "ready", "reason": "Fixture readiness only", "open_questions": []any{}, "checks": checks, "documents": []any{}}))
	if err != nil {
		t.Fatal(err)
	}
	fields := map[string]any{"expected_previous_selection": nil, "action": "select", "solution": map[string]any{"spec_id": "solution", "spec": sr, "problem": problem, "assessment": workspace.TaskDecisionRef{ID: *asm.OutcomeRef.ID, ManifestSHA256: asm.OutcomeRef.ManifestSHA256}}, "reason": "Fixture choice", "human_decision": map[string]any{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-09-29T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: choose this revision."}}
	choice, err := workspace.SelectTaskSolution(w, journalSpecDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", "fixture/select", fields))
	if err != nil {
		t.Fatal(err)
	}
	targetInput := workspace.QueueTargetInput{ProjectID: "ply-journal-fixture", RepoID: "fixture-repo", EpicID: "journal-fixture-epic"}
	reg, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	queue := workspace.WorkspaceTaskQueueDraft{Kind: "WorkspaceTaskQueueDraft@1", SchemaVersion: 1, PublicationKey: "fixture/queue", ProjectID: "ply-journal-fixture", RepoID: "fixture-repo", EpicID: "journal-fixture-epic", ExpectedRevision: 0, Entries: []workspace.QueueEntry{{TaskID: "journal-acceptance", Selection: &workspace.TaskDecisionRef{ID: *choice.OutcomeRef.ID, ManifestSHA256: choice.OutcomeRef.ManifestSHA256}}}, HumanDecision: workspace.QueueHumanDecision{ActorClaim: "synthetic human", DecidedAtUTC: "2026-10-01T00:00:00Z", Source: "explicit_human_instruction", Statement: "Synthetic fixture only."}, RegistryUpgrade: &workspace.TaskRegistryUpgrade{RegistrySHA256: reg.RawSHA256, FromVersion: reg.FormatVersion}}
	qp := writeAny(t, root, "queue.json", queue)
	if _, err = workspace.SetTaskQueue(w, qp); err != nil {
		t.Fatal(err)
	}
	in := workspace.TaskPrepareInput{Target: targetInput, Selector: workspace.QueueSelector{Kind: "next"}}
	plan, err := workspace.PrepareTask(w, in)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Confirmation == nil {
		t.Fatalf("no confirmation: %+v", plan)
	}
	in.Apply = true
	in.Confirmation = *plan.Confirmation
	preparation, err := workspace.PrepareTask(w, in)
	if err != nil {
		t.Fatal(err)
	}
	p := *preparation.Preparation
	selected, err := workspace.ShowTaskSpec(w, workspace.TaskContentQuery{TaskID: "journal-acceptance", SpecID: "solution", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	basis, _ := objectMember(selected.Value, "task_spec_binding")
	inputs, _ := objectMember(selected.Value, "required_inputs")
	draft := journalHandoffDraft(p.Plan.WorktreePath, "refs/heads/"+p.Plan.Branch, oid)
	draft = replaceObjectMember(draft, "schema_version", int64(2))
	draft = append(draft, canonicaljson.Member{Name: "task_spec_binding", Value: basis})
	draft = replaceObjectMember(draft, "inputs", inputs)
	budget, _ := objectMember(draft, "budget")
	draft = replaceObjectMember(draft, "budget", replaceObjectMember(budget.(canonicaljson.Object), "max_rounds", int64(4)))
	raw, err := canonicaljson.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "synthetic-provider")
	if err = os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	ply := filepath.Join(root, "synthetic-ply")
	if err = os.WriteFile(ply, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	policy := filepath.Join(root, "policy.json")
	if err = os.WriteFile(policy, []byte("synthetic policy evidence, not OS attestation"), 0600); err != nil {
		t.Fatal(err)
	}
	request := Request{Envelope: env("request"), RequestKey: "fixture/run", WorkspaceRoot: root, PreparationID: p.ID, PreparationSHA256: digest(p), HandoffDraft: raw, Runtime: Runtime{Provider: "codex", Mode: "interactive", Model: "synthetic-model", Executable: Executable{executable, hashFileTest(t, executable)}, PlyExecutable: Executable{ply, hashFileTest(t, ply)}, PermissionBinding: Permission{AuthorityKind: "reported_contract_with_effective_policy", ProfileID: "synthetic-policy", EffectivePolicySHA256: hash([]byte("policy")), Evidence: []Evidence{{policy, hashFileTest(t, policy), "effective_policy"}}}}, Agreement: Agreement{"A", 3, 5400, 2}, ReturnPolicy: ReturnPolicy{true, true, "separate", "separate"}, HumanAuthority: HumanAuthority{"synthetic human", "human_ordinary_terminal", true}}
	d := SystemDependencies(w)
	d.Executable = func() (string, error) { return request.Runtime.PlyExecutable.Path, nil }
	d.Runner = &fakeRunner{}
	d.CWD = func() (string, error) { return p.Plan.WorktreePath, nil }
	d.ContextPath = func() string { return filepath.Join(runPaths(request).TempRoot, "context.json") }
	return d, request, writeAny(t, root, "request.json", request)
}
