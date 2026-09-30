package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func contentFixtureDraft(t *testing.T, root, kind, key string, fields map[string]canonicaljson.Value) TaskContentInput {
	t.Helper()
	fields["publication_key"] = key
	fields["task_id"] = "task"
	fields["registry_upgrade"] = nil
	fields["recorder"] = contentObject(map[string]canonicaljson.Value{"actor_claim": "fixture agent", "control_surface": "isolated test fixture", "recorded_at_utc": "2026-09-29T12:00:00Z"})
	b, e := canonicaljson.Marshal(contentEnvelope(kind, fields))
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(root, filepath.Base(key)+"-draft.json")
	if e = os.WriteFile(path, b, 0o600); e != nil {
		t.Fatal(e)
	}
	return TaskContentInput{TaskID: "task", File: path}
}

func contentFixtureProblem(t *testing.T, f workItemJourneyFixture, key string) TaskContentInput {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	head := taskContentState(r, "task").ProblemHead
	doc := filepath.Join(root, "problem.md")
	b := []byte("# Revised fixture problem\n\nExplain the changed fixture need.\n")
	if err := os.WriteFile(doc, b, 0600); err != nil {
		t.Fatal(err)
	}
	return contentFixtureDraft(t, root, "WorkspaceTaskProblemDraft@1", key, map[string]canonicaljson.Value{
		"expected_previous": contentRefValue(head), "origin": contentObject(map[string]canonicaljson.Value{"kind": "authored"}),
		"title": "Revised fixture problem", "summary": "Explain the changed fixture need.", "problem_document_id": "problem",
		"documents": []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "problem", "source": contentObject(map[string]canonicaljson.Value{"kind": "file", "locator": doc, "sha256": digestTaskBytes(b), "size_bytes": int64(len(b)), "media_type": "text/markdown", "git_provenance": nil})})},
		"sources":   []canonicaljson.Value{}, "claims": []canonicaljson.Value{}, "deadline": nil, "change_reason": "Fixture problem changed.",
	})
}
func contentFixtureSpec(t *testing.T, f workItemJourneyFixture, key string, previous *TaskRevisionRef) TaskContentInput {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	task, _ := findTask(r, "task")
	epic, _ := findEpic(r, task.ParentEpicID)
	parent, _ := findEpicRepo(*epic, task.RepoID)
	problem := taskContentState(r, "task").ProblemHead
	doc := filepath.Join(root, "solution.md")
	b := []byte("# Fixture solution\n\nPreserve behavior, scope, verification and recovery.\n")
	if e = os.WriteFile(doc, b, 0o600); e != nil {
		t.Fatal(e)
	}
	ref := contentObject(map[string]canonicaljson.Value{"document_id": "solution", "section": nil})
	refs := []canonicaljson.Value{ref}
	part := contentObject(map[string]canonicaljson.Value{"state": "present", "reason": nil, "documents": refs})
	return contentFixtureDraft(t, root, "WorkspaceTaskSpecDraft@1", key, map[string]canonicaljson.Value{"spec_id": "solution", "expected_previous": contentRefValue(previous), "problem": contentRefValue(problem), "title": "Fixture solution", "parts": contentObject(map[string]canonicaljson.Value{"abstract": part, "functional": part, "technical": part}), "documents": []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "solution", "source": contentObject(map[string]canonicaljson.Value{"kind": "file", "locator": doc, "sha256": digestTaskBytes(b), "size_bytes": int64(len(b)), "media_type": "text/markdown", "git_provenance": nil})})}, "supporting": []canonicaljson.Value{}, "requirements": []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "f-01", "functional_refs": refs, "acceptance": "Fixture behavior is preserved.", "verification_ids": []canonicaljson.Value{"unit"}, "technical_refs": refs})}, "removed_requirement_ids": []canonicaljson.Value{}, "phases": []canonicaljson.Value{contentObject(map[string]canonicaljson.Value{"id": "implement", "purpose": "Implement the fixture solution.", "requirement_ids": []canonicaljson.Value{"f-01"}, "entry_criteria": []canonicaljson.Value{}, "exit_criteria": []canonicaljson.Value{"Fixture verified."}, "verification_ids": []canonicaljson.Value{"unit"}})}, "implementation_basis": contentObject(map[string]canonicaljson.Value{"project_id": string(task.ProjectID), "repo_id": string(task.RepoID), "git_common_dir": task.GitCommonDir, "epic_id": string(task.ParentEpicID), "parent_worktree_id": string(parent.Worktree.ID), "parent_ref": parent.Worktree.Ref, "parent_oid": parent.Worktree.OID, "parent_tree": parent.Worktree.Tree, "start_oid": parent.Worktree.OID, "start_tree": parent.Worktree.Tree}), "dependencies": []canonicaljson.Value{}, "change_reason": "Register a fixture solution."})
}
func fixtureSelectTaskSpec(t *testing.T, f workItemJourneyFixture) TaskSpecBasis {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	spec, e := RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "fixture/spec", nil))
	if e != nil {
		t.Fatal(e)
	}
	ref := &TaskRevisionRef{*spec.OutcomeRef.Revision, spec.OutcomeRef.ManifestSHA256}
	checks := []canonicaljson.Value{}
	for _, id := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, contentObject(map[string]canonicaljson.Value{"id": id, "outcome": "pass", "reason": "Fixture assessment of the preserved solution.", "evidence_document_ids": []canonicaljson.Value{}}))
	}
	assessment, e := AssessTaskSpec(f.dependencies, contentFixtureDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", "fixture/assessment", map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(ref), "expected_previous_assessment": nil, "outcome": "ready", "reason": "Fixture readiness claim.", "open_questions": []canonicaljson.Value{}, "checks": checks, "documents": []canonicaljson.Value{}}))
	if e != nil {
		t.Fatal(e)
	}
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	p := taskContentState(r, "task").ProblemHead
	_, e = SelectTaskSolution(f.dependencies, contentFixtureDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", "fixture/selection", map[string]canonicaljson.Value{"expected_previous_selection": nil, "action": "select", "solution": contentObject(map[string]canonicaljson.Value{"spec_id": "solution", "spec": contentRefValue(ref), "problem": contentRefValue(p), "assessment": contentRefValue(&TaskDecisionRef{*assessment.OutcomeRef.ID, assessment.OutcomeRef.ManifestSHA256})}), "reason": "Fixture selection.", "human_decision": contentObject(map[string]canonicaljson.Value{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-09-29T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: select this solution for the isolated test."})}))
	if e != nil {
		t.Fatal(e)
	}
	r, e = f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	task, _ := findTask(r, "task")
	ps, rs, e := f.dependencies.Projects.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	eval, e := currentTaskSpec(f.dependencies, root, r, ProjectSnapshot{Projects: ps, Repos: rs}, *task, true)
	if e != nil || eval.Basis == nil {
		t.Fatalf("basis: %#v %v", eval, e)
	}
	return *eval.Basis
}
func TestTaskContentCreatePreservesAtomicProblemAndIdenticalRetry(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	if r.FormatVersion != 3 || len(r.TaskProblemRevisions) != 1 || !taskRequiresSpec(r, "task") || len(r.TaskContentPublications) != 1 {
		t.Fatalf("incomplete initial problem: %#v", r)
	}
	before, e := os.ReadFile(workItemsPath(root))
	if e != nil {
		t.Fatal(e)
	}
	task, _ := findTask(r, "task")
	input, _ := ParseTaskCreateInput("task", task.Title, task.Description, string(task.ParentEpicID), string(task.ProjectID), string(task.RepoID))
	got, e := CreateTask(f.dependencies, input)
	if e != nil || got.Created {
		t.Fatalf("retry: %#v %v", got, e)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if string(before) != string(after) {
		t.Fatal("identical create rewrote registry")
	}
	p, e := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"})
	if e != nil || contentString(contentFields(p.Value), "integrity") != "valid" {
		t.Fatalf("problem: %#v %v", p, e)
	}
}
func TestTaskContentSourceDeletionRetryCASAndSelection(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(f.wrapper, "task"), f.oid)
	if _, e := CreateTaskWorktree(f.dependencies, input); e != nil {
		t.Fatal(e)
	}
	basis := fixtureSelectTaskSpec(t, f)
	draft := contentFixtureSpec(t, f, "fixture/spec", nil)
	if e := os.Remove(filepath.Join(root, "solution.md")); e != nil {
		t.Fatal(e)
	}
	retry, e := RecordTaskSpec(f.dependencies, draft)
	if e != nil || retry.Classification != "existing" {
		t.Fatalf("retry requires deleted source: %#v %v", retry, e)
	}
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := historicalTaskSpec(f.dependencies, root, r, basis); e != nil {
		t.Fatal(e)
	}
	conflict, e := RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "fixture/conflicting-spec", nil))
	if e == nil || conflict.Attempted.ObjectPublicationStarted {
		t.Fatalf("CAS conflict: %#v %v", conflict, e)
	}
	r2, e := RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "fixture/spec-2", &basis.Spec))
	if e != nil {
		t.Fatal(e)
	}
	if *r2.OutcomeRef.Revision != 2 {
		t.Fatal("revision sequence")
	}
	r, _ = f.dependencies.WorkItems.Snapshot(root)
	ps, rs, _ := f.dependencies.Projects.Snapshot(root)
	task, _ := findTask(r, "task")
	eval, e := currentTaskSpec(f.dependencies, root, r, ProjectSnapshot{Projects: ps, Repos: rs}, *task, true)
	if e != nil || !contentTypedEqual(eval.Basis, &basis) || eval.SelectionFreshness != "current" {
		t.Fatalf("new draft rebound selection: %#v %v", eval, e)
	}
}

func TestTaskResultRetainsStartBasisWhenProblemChangesBeforeRecord(t *testing.T) {
	f := newTaskDeliveryFixture(t, false)
	if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f.workItemJourneyFixture, "fixture/p2")); err != nil {
		t.Fatal(err)
	}
	result, err := RecordTaskResult(f.dependencies, TaskResultRecordInput{TaskID: "task", File: f.draftPath})
	if err != nil {
		t.Fatal(err)
	}
	if result.SpecBinding == nil || result.SpecBinding.Basis.Problem.Revision != 1 || result.SpecBinding.TaskResultID != result.Record.ID {
		t.Fatalf("result rebound to current problem: %#v", result)
	}
}

func TestTaskContentReselectAndWithdrawPreserveHistory(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.wrapper, "task"), ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	basis := fixtureSelectTaskSpec(t, f)
	raw, err := os.ReadFile(filepath.Join(root, "selection-draft.json"))
	if err != nil {
		t.Fatal(err)
	}
	v, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(v)
	m["publication_key"] = "fixture/reselect"
	m["expected_previous_selection"] = contentRefValue(basis.Selection)
	raw, _ = canonicaljson.Marshal(contentObject(m))
	path := filepath.Join(root, "reselect.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	reselected, err := SelectTaskSolution(f.dependencies, TaskContentInput{TaskID: "task", File: path})
	if err != nil || *reselected.OutcomeRef.ID == basis.Selection.ID {
		t.Fatalf("reselect did not create a new decision: %#v %v", reselected, err)
	}
	m["publication_key"] = "fixture/withdraw"
	m["expected_previous_selection"] = contentRefValue(TaskDecisionRef{*reselected.OutcomeRef.ID, reselected.OutcomeRef.ManifestSHA256})
	m["action"] = "withdraw"
	m["solution"] = nil
	raw, _ = canonicaljson.Marshal(contentObject(m))
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	withdrawn, err := WithdrawTaskSolution(f.dependencies, TaskContentInput{TaskID: "task", File: path})
	if err != nil {
		t.Fatal(err)
	}
	retry, err := WithdrawTaskSolution(f.dependencies, TaskContentInput{TaskID: "task", File: path})
	if err != nil || *retry.OutcomeRef.ID != *withdrawn.OutcomeRef.ID {
		t.Fatalf("withdraw retry: %#v %v", retry, err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil || len(r.TaskSolutionSelections) != 3 {
		t.Fatalf("lost decision history: %#v %v", r, err)
	}
	if _, err = historicalTaskSpec(f.dependencies, root, r, basis); err != nil {
		t.Fatalf("withdraw erased old basis: %v", err)
	}
	if err = WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
		_, e := s.ValidateTarget(r.Tasks[0].Worktree.Locator, r.Tasks[0].Worktree.Ref, r.Tasks[0].GitCommonDir, &basis)
		return e
	}); err == nil {
		t.Fatal("withdrawn selection still authorizes start")
	}
}

func TestTaskContentNegativeAssessmentBlocksPreviousChoice(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.wrapper, "task"), ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	basis := fixtureSelectTaskSpec(t, f)
	raw, _ := os.ReadFile(filepath.Join(root, "assessment-draft.json"))
	v, _ := canonicaljson.DecodeStrict(raw)
	m := contentFields(v)
	m["publication_key"], m["expected_previous_assessment"], m["outcome"] = "fixture/negative", contentRefValue(basis.Assessment), "needs_work"
	checks := contentArray(m, "checks")
	check := contentFields(checks[0])
	check["outcome"] = "fail"
	checks[0] = contentObject(check)
	m["checks"] = checks
	raw, _ = canonicaljson.Marshal(contentObject(m))
	path := filepath.Join(root, "negative.json")
	if err := os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	negative, err := AssessTaskSpec(f.dependencies, TaskContentInput{TaskID: "task", File: path})
	if err != nil {
		t.Fatal(err)
	}
	r, _ := f.dependencies.WorkItems.Snapshot(root)
	task := r.Tasks[0]
	if err = WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
		_, e := s.ValidateTarget(task.Worktree.Locator, task.Worktree.Ref, task.GitCommonDir, &basis)
		return e
	}); err == nil {
		t.Fatal("later negative assessment did not block start")
	}
	raw, _ = os.ReadFile(filepath.Join(root, "selection-draft.json"))
	v, _ = canonicaljson.DecodeStrict(raw)
	m = contentFields(v)
	m["publication_key"] = "fixture/reselect-negative"
	m["expected_previous_selection"] = contentRefValue(basis.Selection)
	solution := contentFields(m["solution"])
	solution["assessment"] = contentRefValue(TaskDecisionRef{*negative.OutcomeRef.ID, negative.OutcomeRef.ManifestSHA256})
	m["solution"] = contentObject(solution)
	raw, _ = canonicaljson.Marshal(contentObject(m))
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = SelectTaskSolution(f.dependencies, TaskContentInput{TaskID: "task", File: path}); err == nil {
		t.Fatal("negative assessment was selected")
	}
}
