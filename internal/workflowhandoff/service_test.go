package workflowhandoff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestCreatePublishesBoundHandoffAndRetriesIdentically(t *testing.T) {
	workspaceRoot, target, ref, oid := prepareServiceWorkspace(t)
	draft := minimalHandoffDraftValue(target, ref, oid)
	draftBytes, err := canonicaljson.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	draftPath := filepath.Join(workspaceRoot, "draft.json")
	if err := os.WriteFile(draftPath, draftBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	dependencies := SystemDependencies()
	first, err := Create(dependencies, CreateInput{DraftPath: draftPath})
	if err != nil {
		t.Fatal(err)
	}
	if !first.Created || first.Purpose != "Test handoff" || first.Worktree != target || first.Locator == "" {
		t.Fatalf("first result = %#v", first)
	}
	second, err := Create(dependencies, CreateInput{DraftPath: draftPath})
	if err != nil {
		t.Fatal(err)
	}
	if second.Created || second.HandoffID != first.HandoffID || second.Locator != first.Locator {
		t.Fatalf("retry result = %#v, first %#v", second, first)
	}
	if _, err := os.Stat(first.Locator); err != nil {
		t.Fatal(err)
	}
}

func TestSubmitStartAndTerminalResultRoundTrip(t *testing.T) {
	workspaceRoot, target, ref, oid := prepareServiceWorkspace(t)
	draftBytes, err := canonicaljson.Marshal(minimalHandoffDraftValue(target, ref, oid))
	if err != nil {
		t.Fatal(err)
	}
	handoffDraftPath := filepath.Join(workspaceRoot, "handoff-draft.json")
	if err := os.WriteFile(handoffDraftPath, draftBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	dependencies := SystemDependencies()
	created, err := Create(dependencies, CreateInput{DraftPath: handoffDraftPath})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principal := canonicaljson.Object{
		{Name: "expected_principal_id", Value: "codex-delivery-agent"},
		{Name: "human_start_principal", Value: "product-owner"},
		{Name: "start_surface", Value: "Codex chat"},
		{Name: "session_id", Value: "session-test"},
		{Name: "runtime_id", Value: "codex"},
		{Name: "model_id", Value: "test-model"},
		{Name: "started_at_utc", Value: "2026-09-27T00:00:00Z"},
	}
	startValue := minimalStartDraftValue(t, snapshot, principal)
	startBytes, _ := canonicaljson.Marshal(startValue)
	startPath := filepath.Join(workspaceRoot, "start.json")
	if err := os.WriteFile(startPath, startBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	start, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil || !start.Created || start.Phase != "start" {
		t.Fatalf("start = %#v, %v", start, err)
	}
	startRetry, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil || startRetry.Created || startRetry.SHA256 != start.SHA256 {
		t.Fatalf("start retry = %#v, %v", startRetry, err)
	}

	resultValue := minimalTerminalDraftValue(snapshot, principal, start)
	resultBytes, _ := canonicaljson.Marshal(resultValue)
	resultPath := filepath.Join(workspaceRoot, "result.json")
	if err := os.WriteFile(resultPath, resultBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath})
	if err != nil || !result.Created || result.Phase != "terminal" {
		t.Fatalf("result = %#v, %v", result, err)
	}
	resultRetry, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath})
	if err != nil || resultRetry.Created || resultRetry.SHA256 != result.SHA256 {
		t.Fatalf("result retry = %#v, %v", resultRetry, err)
	}
}

func TestLifecycleControlsAndRequiredSequence(t *testing.T) {
	t.Run("terminal before start", func(t *testing.T) {
		dependencies, created, snapshot, principal := createTestHandoff(t)
		resultValue := minimalTerminalDraftValue(snapshot, principal, SubmitResult{DocumentID: snapshot.Handoff.Identity.StartReceiptID, SHA256: digestBytes([]byte("missing"))})
		resultPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "early-result.json", resultValue)
		if _, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: resultPath}); !IsClass(err, ErrorStartRequired) {
			t.Fatalf("terminal-before-start error = %v", err)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		dependencies, created, snapshot, _ := createTestHandoff(t)
		if _, err := Cancel(dependencies, ControlInput{HandoffID: created.HandoffID, Reason: "The work is no longer needed."}); err != nil {
			t.Fatal(err)
		}
		read, err := dependencies.Store.ReadByLocator(snapshot.Handoff.Locator)
		if err != nil || read.HeadState != "cancelled" {
			t.Fatalf("cancelled snapshot = %#v, %v", read, err)
		}
	})

	t.Run("supersede", func(t *testing.T) {
		dependencies, created, snapshot, _ := createTestHandoff(t)
		replacement := replaceObjectMember(minimalHandoffDraftValue(snapshot.Handoff.Target.Worktree, snapshot.Handoff.Target.Ref, snapshot.Handoff.Target.OID), "publication_key", "wf01/service-test-replacement")
		replacementPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "replacement.json", replacement)
		createdReplacement, err := Supersede(dependencies, SupersedeInput{HandoffID: created.HandoffID, DraftPath: replacementPath, Reason: "Correct the immutable input."})
		if err != nil {
			t.Fatal(err)
		}
		old, err := dependencies.Store.ReadByLocator(created.Locator)
		if err != nil || old.HeadState != "superseded" || old.ReplacementID != string(createdReplacement.HandoffID) {
			t.Fatalf("old snapshot = %#v, %v", old, err)
		}
		current, err := dependencies.Store.ReadByLocator(createdReplacement.Locator)
		if err != nil || current.HeadState != "ready" || current.Handoff.Identity.ActivityID != old.Handoff.Identity.ActivityID {
			t.Fatalf("replacement snapshot = %#v, %v", current, err)
		}
	})

	t.Run("abandon after start", func(t *testing.T) {
		dependencies, created, snapshot, principal := createTestHandoff(t)
		startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "start.json", minimalStartDraftValue(t, snapshot, principal))
		if _, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath}); err != nil {
			t.Fatal(err)
		}
		if _, err := Abandon(dependencies, ControlInput{HandoffID: created.HandoffID, Reason: "The recipient session disappeared.", AcknowledgeEffectsUnknown: true}); err != nil {
			t.Fatal(err)
		}
		read, err := dependencies.Store.ReadByLocator(created.Locator)
		if err != nil || read.HeadState != "abandoned_unknown" {
			t.Fatalf("abandoned snapshot = %#v, %v", read, err)
		}
	})
}

func TestStartRejectsTargetDrift(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "start.json", minimalStartDraftValue(t, snapshot, principal))
	if err := os.WriteFile(filepath.Join(snapshot.Handoff.Target.Worktree, "README.md"), []byte("drift\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath}); !IsClass(err, ErrorConflict) {
		t.Fatalf("drift error = %v", err)
	}
}

func createTestHandoff(t *testing.T) (Dependencies, CreateResult, Snapshot, canonicaljson.Object) {
	t.Helper()
	workspaceRoot, target, ref, oid := prepareServiceWorkspace(t)
	draftPath := writeCanonicalTestFile(t, workspaceRoot, "handoff-draft.json", minimalHandoffDraftValue(target, ref, oid))
	dependencies := SystemDependencies()
	created, err := Create(dependencies, CreateInput{DraftPath: draftPath})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principal := canonicaljson.Object{{Name: "expected_principal_id", Value: "codex-delivery-agent"}, {Name: "human_start_principal", Value: "product-owner"}, {Name: "start_surface", Value: "Codex chat"}, {Name: "session_id", Value: "session-test"}, {Name: "runtime_id", Value: "codex"}, {Name: "model_id", Value: "test-model"}, {Name: "started_at_utc", Value: "2026-09-27T00:00:00Z"}}
	return dependencies, created, snapshot, principal
}

func writeCanonicalTestFile(t *testing.T, directory, name string, value canonicaljson.Value) string {
	t.Helper()
	contents, err := canonicaljson.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func minimalStartDraftValue(t *testing.T, snapshot Snapshot, principal canonicaljson.Object) canonicaljson.Object {
	t.Helper()
	tempRoot := filepath.Join(snapshot.Handoff.Workspace.Root, "recipient-temp")
	if err := os.MkdirAll(tempRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(tempRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	contractDigest := func(name string) string {
		value, _ := objectMember(snapshot.Handoff.Value, name)
		bytes, err := canonicaljson.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return digestBytes(bytes)
	}
	projectValue, _ := objectMember(snapshot.Handoff.Value, "project_binding")
	project := projectValue.(canonicaljson.Object)
	return canonicaljson.Object{
		{Name: "kind", Value: "ply.workflow.start-receipt-draft"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "receipt_id", Value: snapshot.Handoff.Identity.StartReceiptID},
		{Name: "binding", Value: canonicaljson.Object{{Name: "activity_id", Value: snapshot.Handoff.Identity.ActivityID}, {Name: "run_id", Value: snapshot.Handoff.Identity.RunID}, {Name: "handoff_id", Value: snapshot.Handoff.Identity.HandoffID}, {Name: "handoff_sha256", Value: snapshot.Handoff.SHA256}}},
		{Name: "principal", Value: principal},
		{Name: "observed_workspace", Value: canonicaljson.Object{{Name: "root", Value: snapshot.Handoff.Workspace.Root}, {Name: "marker_format_version", Value: int64(snapshot.Handoff.Workspace.MarkerFormatVersion)}, {Name: "marker_sha256", Value: snapshot.Handoff.Workspace.MarkerSHA256}, {Name: "matches_expected", Value: true}}},
		{Name: "observed_project", Value: canonicaljson.Object{{Name: "project_id", Value: objectString(project, "project_id")}, {Name: "repo_id", Value: objectString(project, "repo_id")}, {Name: "registered_locator", Value: objectString(project, "registered_locator")}, {Name: "registered_git_common_dir", Value: objectString(project, "registered_git_common_dir")}, {Name: "matches_expected", Value: true}}},
		{Name: "observed_target", Value: append(targetValue(snapshot.Handoff.Target), canonicaljson.Member{Name: "matches_expected", Value: true})},
		{Name: "observed_inputs", Value: []canonicaljson.Value{}},
		{Name: "contract_digests", Value: canonicaljson.Object{{Name: "authority_sha256", Value: contractDigest("authority")}, {Name: "budget_sha256", Value: contractDigest("budget")}, {Name: "verifiers_sha256", Value: contractDigest("verifiers")}, {Name: "stop_conditions_sha256", Value: contractDigest("stop_conditions")}}},
		{Name: "sandbox", Value: canonicaljson.Object{{Name: "read_roots", Value: []canonicaljson.Value{snapshot.Handoff.Locator, snapshot.Handoff.Target.Worktree}}, {Name: "write_roots", Value: []canonicaljson.Value{snapshot.Handoff.ReplyRoot, snapshot.Handoff.Target.Worktree}}, {Name: "temp_root", Value: tempRoot}, {Name: "matches_contract", Value: true}}},
		{Name: "acceptance", Value: "started"}, {Name: "issues", Value: []canonicaljson.Value{}},
	}
}

func minimalTerminalDraftValue(snapshot Snapshot, principal canonicaljson.Object, start SubmitResult) canonicaljson.Object {
	authorityValue, _ := objectMember(snapshot.Handoff.Value, "authority")
	authority := authorityValue.(canonicaljson.Object)
	allowedValue, _ := objectMember(authority, "allowed_effects")
	allowed := allowedValue.([]canonicaljson.Value)
	observedEffects := make([]canonicaljson.Value, 0, len(allowed))
	for _, effectValue := range allowed {
		effect := effectValue.(canonicaljson.Object)
		observedEffects = append(observedEffects, canonicaljson.Object{
			{Name: "effect_id", Value: objectString(effect, "id")},
			{Name: "type", Value: objectString(effect, "type")},
			{Name: "scope", Value: func() canonicaljson.Value { value, _ := objectMember(effect, "scope"); return value }()},
			{Name: "occurrences", Value: int64(1)},
			{Name: "within_authority", Value: true},
		})
	}
	verifierResults := []canonicaljson.Value{canonicaljson.Object{
		{Name: "verifier_id", Value: "test"},
		{Name: "argv", Value: []canonicaljson.Value{"go", "test", "./..."}},
		{Name: "cwd", Value: snapshot.Handoff.Target.Worktree},
		{Name: "exit", Value: int64(0)},
		{Name: "bound_oid_or_sha256", Value: snapshot.Handoff.Target.OID},
		{Name: "stdout_artifact_id", Value: nil},
		{Name: "stderr_artifact_id", Value: nil},
	}}
	return canonicaljson.Object{
		{Name: "kind", Value: "ply.workflow.terminal-result-draft"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "result_id", Value: snapshot.Handoff.Identity.TerminalResultID},
		{Name: "binding", Value: canonicaljson.Object{{Name: "activity_id", Value: snapshot.Handoff.Identity.ActivityID}, {Name: "run_id", Value: snapshot.Handoff.Identity.RunID}, {Name: "handoff_id", Value: snapshot.Handoff.Identity.HandoffID}, {Name: "handoff_sha256", Value: snapshot.Handoff.SHA256}}},
		{Name: "start_binding", Value: canonicaljson.Object{{Name: "receipt_id", Value: start.DocumentID}, {Name: "start_receipt_sha256", Value: start.SHA256}}},
		{Name: "principal", Value: principal}, {Name: "reported_outcome", Value: "complete"}, {Name: "rounds_used", Value: int64(1)}, {Name: "stop_reasons", Value: []canonicaljson.Value{}},
		{Name: "summary", Value: "The bounded change is complete."}, {Name: "meaning", Value: "Technical checks were reported as passing."},
		{Name: "final_target", Value: append(targetValue(snapshot.Handoff.Target), canonicaljson.Member{Name: "matches_expected", Value: true})},
		{Name: "observed_effects", Value: observedEffects}, {Name: "verifier_results", Value: verifierResults},
		{Name: "review", Value: canonicaljson.Object{{Name: "findings", Value: []canonicaljson.Value{}}, {Name: "fixes", Value: []canonicaljson.Value{}}, {Name: "open_actionable_findings", Value: []canonicaljson.Value{}}}},
		{Name: "artifacts", Value: []canonicaljson.Value{}}, {Name: "evidence_gaps", Value: []canonicaljson.Value{}}, {Name: "forbidden_effects_observed", Value: []canonicaljson.Value{}},
	}
}

func prepareServiceWorkspace(t *testing.T) (string, string, string, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
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
	if _, err := workspace.AddProject(workspaceDependencies, workspace.ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "ply", Path: target}}}); err != nil {
		t.Fatal(err)
	}
	return root, target, ref, oid
}

func minimalHandoffDraftValue(target, ref, oid string) canonicaljson.Object {
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
		{Name: "binding_request", Value: canonicaljson.Object{{Name: "project_id", Value: "ply"}, {Name: "repo_id", Value: "ply"}, {Name: "target_worktree", Value: target}, {Name: "target_ref", Value: ref}, {Name: "expected_oid", Value: oid}, {Name: "status_policy", Value: canonicaljson.Object{{Name: "mode", Value: "clean"}}}}},
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
		{Name: "procedure", Value: []canonicaljson.Value{canonicaljson.Object{{Name: "id", Value: "implement"}, {Name: "instruction", Value: "Implement the bounded change."}, {Name: "required_before", Value: []canonicaljson.Value{}}}}},
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
