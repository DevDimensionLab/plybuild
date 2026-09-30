package workflowhandoff

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTaskSpecGeneralVersionTwoProofAndVersionMatrix(t *testing.T) {
	root, target, ref, oid := prepareServiceWorkspace(t)
	d := SystemDependencies()
	draft := append(replaceObjectMember(minimalHandoffDraftValue(target, ref, oid), "schema_version", int64(2)), canonicaljson.Member{Name: "task_spec_binding", Value: nil})
	path := writeCanonicalTestFile(t, root, "v2-draft.json", draft)
	created, err := Create(d, CreateInput{DraftPath: path})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := d.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principal := canonicaljson.Object{{Name: "expected_principal_id", Value: "codex-delivery-agent"}, {Name: "human_start_principal", Value: "fixture human, not actual approval"}, {Name: "start_surface", Value: "isolated fixture"}, {Name: "session_id", Value: "fixture-session"}, {Name: "runtime_id", Value: "codex"}, {Name: "model_id", Value: "test-model"}, {Name: "started_at_utc", Value: "2026-09-29T12:00:00Z"}}
	old := minimalStartDraftValue(t, snapshot, principal)
	oldPath := writeCanonicalTestFile(t, root, "v1-start.json", old)
	if _, err = SubmitStart(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: oldPath}); err == nil {
		t.Fatal("version 1 start accepted for version 2 null handoff")
	}
	start := append(replaceObjectMember(old, "schema_version", int64(2)), canonicaljson.Member{Name: "task_spec_binding", Value: nil})
	negative := replaceObjectMember(replaceObjectMember(start, "acceptance", "conflict"), "issues", []canonicaljson.Value{canonicaljson.Object{{Name: "type", Value: "task_spec"}, {Name: "detail", Value: "Fixture selection changed."}}})
	negativeBytes, _ := canonicaljson.Marshal(negative)
	if _, err := decodeStartDraft(negativeBytes); err != nil {
		t.Fatalf("version 2 task_spec issue rejected: %v", err)
	}
	sandbox, _ := objectMember(start, "sandbox")
	start = replaceObjectMember(start, "sandbox", replaceObjectMember(sandbox.(canonicaljson.Object), "read_roots", []canonicaljson.Value{root}))
	sandbox, _ = objectMember(start, "sandbox")
	externalEpic := t.TempDir()
	if err := validateSandboxContract(d.Files, snapshot, sandbox.(canonicaljson.Object), externalEpic); err == nil {
		t.Fatal("sandbox omitted required external Epic Git source")
	}
	startPath := writeCanonicalTestFile(t, root, "v2-start.json", start)
	accepted, err := SubmitStart(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err = d.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	observation, ok := objectMember(snapshot.Start.Value, "task_spec_observation")
	if !ok || observation != nil {
		t.Fatal("general start must carry explicit null observation")
	}
	if err = validateStoredAccepted(snapshot.Start.Value, "start", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret); err != nil {
		t.Fatal(err)
	}
	proof, _ := objectMember(snapshot.Start.Value, "capability_proof")
	forged := replaceObjectMember(snapshot.Start.Value, "capability_proof", replaceObjectMember(proof.(canonicaljson.Object), "algorithm", "none"))
	if err = validateStoredAccepted(forged, "start", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret); err == nil {
		t.Fatal("version 2 accepted unsupported proof algorithm")
	}
	terminal := minimalTerminalDraftValue(snapshot, principal, accepted)
	terminalPath := writeCanonicalTestFile(t, root, "v1-terminal.json", terminal)
	if _, err = SubmitResultDocument(d, SubmitInput{HandoffLocator: created.Locator, DraftPath: terminalPath}); err != nil {
		t.Fatal(err)
	}
	if _, err = d.Store.ReadByLocator(created.Locator); err != nil {
		t.Fatal(err)
	}
}

func TestTaskSpecDraftRejectsMissingAndUnknownVersionFields(t *testing.T) {
	_, target, ref, oid := prepareServiceWorkspace(t)
	base := minimalHandoffDraftValue(target, ref, oid)
	for _, v := range []canonicaljson.Object{replaceObjectMember(base, "schema_version", int64(2)), append(base, canonicaljson.Member{Name: "task_spec_binding", Value: nil}), append(replaceObjectMember(base, "schema_version", int64(3)), canonicaljson.Member{Name: "task_spec_binding", Value: nil})} {
		b, err := canonicaljson.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = decodeHandoffDraft(b); err == nil {
			t.Fatal("unsupported wire accepted")
		}
	}
}

// These are local receipt/human-claim fixtures; no agent or human QA runs.
func taskSpecFixtureValue(t *testing.T, v any) canonicaljson.Value {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	out, err := canonicaljson.DecodeStrict(b)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func taskSpecFixtureDraft(t *testing.T, root, kind, key string, fields map[string]any) workspace.TaskContentInput {
	t.Helper()
	fields["kind"], fields["schema_version"], fields["format"], fields["format_version"], fields["canonicalization"] = kind, 1, "json", 1, "RFC8785"
	fields["publication_key"], fields["task_id"], fields["registry_upgrade"] = key, "task", nil
	fields["recorder"] = map[string]any{"actor_claim": "fixture agent", "control_surface": "isolated service test", "recorded_at_utc": "2026-09-29T12:00:00Z"}
	return workspace.TaskContentInput{TaskID: "task", File: writeCanonicalTestFile(t, root, filepath.Base(key)+".json", taskSpecFixtureValue(t, fields).(canonicaljson.Object))}
}
func prepareBoundTaskStart(t *testing.T) (Dependencies, workspace.Dependencies, string, SubmitInput, workspace.TaskContentInput) {
	t.Helper()
	root, epic, ref, oid := prepareServiceWorkspace(t)
	w := workspace.SystemDependencies()
	if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Fixture Epic", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: ref, ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "task", Title: "Fixture Task", Description: "Fixture need", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
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
	spec, err := workspace.RecordTaskSpec(w, taskSpecFixtureDraft(t, root, "WorkspaceTaskSpecDraft@1", "fixture/spec", map[string]any{
		"spec_id": "solution", "expected_previous": nil, "problem": problem, "title": "Fixture solution",
		"parts":      map[string]any{"abstract": part, "functional": part, "technical": part},
		"documents":  []any{map[string]any{"id": "solution", "source": map[string]any{"kind": "file", "locator": doc, "sha256": digestBytes(b), "size_bytes": len(b), "media_type": "text/markdown", "git_provenance": nil}}},
		"supporting": []any{}, "requirements": []any{map[string]any{"id": "f-01", "functional_refs": refs, "technical_refs": refs, "acceptance": "Fixture behavior preserved.", "verification_ids": []string{"test"}}},
		"removed_requirement_ids": []any{}, "phases": []any{map[string]any{"id": "implement", "purpose": "Fixture implementation", "requirement_ids": []string{"f-01"}, "entry_criteria": []any{}, "exit_criteria": []string{"Verified"}, "verification_ids": []string{"test"}}},
		"implementation_basis": map[string]any{"project_id": "ply", "repo_id": "ply", "git_common_dir": r.Tasks[0].GitCommonDir, "epic_id": "epic", "parent_worktree_id": parent.ID, "parent_ref": parent.Ref, "parent_oid": parent.OID, "parent_tree": parent.Tree, "start_oid": parent.OID, "start_tree": parent.Tree}, "dependencies": []any{}, "change_reason": "Fixture solution"}))
	if err != nil {
		t.Fatal(err)
	}
	sr := workspace.TaskRevisionRef{Revision: *spec.OutcomeRef.Revision, ManifestSHA256: spec.OutcomeRef.ManifestSHA256}
	checks := []any{}
	for _, id := range []string{"acceptance_coverage", "implementation_basis", "problem_coverage", "recovery", "scope_and_phases", "three_parts"} {
		checks = append(checks, map[string]any{"id": id, "outcome": "pass", "reason": "Fixture assessment only", "evidence_document_ids": []any{}})
	}
	asm, err := workspace.AssessTaskSpec(w, taskSpecFixtureDraft(t, root, "WorkspaceTaskSpecAssessmentDraft@1", "fixture/assess", map[string]any{"spec_id": "solution", "spec": sr, "expected_previous_assessment": nil, "outcome": "ready", "reason": "Fixture readiness only", "open_questions": []any{}, "checks": checks, "documents": []any{}}))
	if err != nil {
		t.Fatal(err)
	}
	fields := map[string]any{"expected_previous_selection": nil, "action": "select", "solution": map[string]any{"spec_id": "solution", "spec": sr, "problem": problem, "assessment": workspace.TaskDecisionRef{ID: *asm.OutcomeRef.ID, ManifestSHA256: asm.OutcomeRef.ManifestSHA256}}, "reason": "Fixture choice", "human_decision": map[string]any{"actor_claim": "fixture human, not actual approval", "decided_at_utc": "2026-09-29T12:00:00Z", "source": "explicit_human_instruction", "statement": "Fixture only: choose this revision."}}
	choice, err := workspace.SelectTaskSolution(w, taskSpecFixtureDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", "fixture/select", fields))
	if err != nil {
		t.Fatal(err)
	}
	// A ready human choice can precede worktree creation, without start authority.
	show, err := workspace.ShowTask(w, "task")
	if err != nil {
		t.Fatal(err)
	}
	ready, _ := objectMember(show.Content, "ready_for_spec_handoff")
	if ready != false {
		t.Fatal("choice without a worktree opened handoff readiness")
	}
	target := filepath.Join(root, "task")
	if _, err = workspace.CreateTaskWorktree(w, workspace.TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: target, ExpectedParentOID: oid}); err != nil {
		t.Fatal(err)
	}
	selected, err := workspace.ShowTaskSpec(w, workspace.TaskContentQuery{TaskID: "task", SpecID: "solution", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	basis, _ := objectMember(selected.Value, "task_spec_binding")
	required, _ := objectMember(selected.Value, "required_inputs")
	inputs := []canonicaljson.Value{}
	for _, v := range required.([]canonicaljson.Value) {
		inputs = append(inputs, v)
	}
	draft := minimalHandoffDraftValue(target, "refs/heads/task", oid)
	draft = replaceObjectMember(draft, "schema_version", int64(2))
	draft = append(draft, canonicaljson.Member{Name: "task_spec_binding", Value: basis})
	draft = replaceObjectMember(draft, "inputs", inputs)
	procedure, _ := objectMember(draft, "procedure")
	steps := procedure.([]canonicaljson.Value)
	steps[0] = replaceObjectMember(steps[0].(canonicaljson.Object), "instruction", "Fixture only: preserve a task-requirements artifact.")
	draft = replaceObjectMember(draft, "procedure", steps)
	d := SystemDependencies()
	created, err := Create(d, CreateInput{DraftPath: writeCanonicalTestFile(t, root, "handoff.json", draft)})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := d.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principal := taskSpecFixtureValue(t, map[string]any{"expected_principal_id": "codex-delivery-agent", "human_start_principal": "fixture human, not actual approval", "start_surface": "isolated service fixture", "session_id": "fixture-session", "runtime_id": "fixture", "model_id": "fixture", "started_at_utc": "2026-09-29T12:00:00Z"}).(canonicaljson.Object)
	start := minimalStartDraftValue(t, snapshot, principal)
	start = replaceObjectMember(start, "schema_version", int64(2))
	start = append(start, canonicaljson.Member{Name: "task_spec_binding", Value: basis})
	observedInputs := []canonicaljson.Value{}
	for _, v := range inputs {
		in := v.(canonicaljson.Object)
		in = removeObjectMember(removeObjectMember(removeObjectMember(in, "role"), "media_type"), "git_binding")
		observedInputs = append(observedInputs, append(in, canonicaljson.Member{Name: "matches_expected", Value: true}))
	}
	start = replaceObjectMember(start, "observed_inputs", observedInputs)
	sandbox, _ := objectMember(start, "sandbox")
	start = replaceObjectMember(start, "sandbox", replaceObjectMember(sandbox.(canonicaljson.Object), "read_roots", []canonicaljson.Value{root}))
	input := SubmitInput{HandoffLocator: created.Locator, DraftPath: writeCanonicalTestFile(t, root, "start.json", start)}
	fields["expected_previous_selection"] = workspace.TaskDecisionRef{ID: *choice.OutcomeRef.ID, ManifestSHA256: choice.OutcomeRef.ManifestSHA256}
	reselect := taskSpecFixtureDraft(t, root, "WorkspaceTaskSolutionSelectionDraft@1", "fixture/reselect", fields)
	return d, w, root, input, reselect
}

type barrierTaskSelectionStore struct {
	workspace.WorkItemStore
	entered chan struct{}
	release chan struct{}
}

func (s barrierTaskSelectionStore) WithLock(root string, fn func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error { close(s.entered); <-s.release; return fn(session) })
}

type barrierTaskStartStore struct {
	Store
	entered chan struct{}
	release chan struct{}
}

func (s barrierTaskStartStore) SubmitStart(in SubmitStoreInput) (SubmitStoreResult, error) {
	close(s.entered)
	<-s.release
	return s.Store.SubmitStart(in)
}

func TestTaskSpecServiceSelectAndAcceptedStartUseExistingLocks(t *testing.T) {
	for _, selectFirst := range []bool{true, false} {
		t.Run(map[bool]string{true: "select-before-start", false: "accepted-start-before-select"}[selectFirst], func(t *testing.T) {
			d, w, _, input, reselect := prepareBoundTaskStart(t)
			entered, release := make(chan struct{}), make(chan struct{})
			first, second := make(chan error, 1), make(chan error, 1)
			originalStore := d.Store
			if selectFirst {
				w.WorkItems = barrierTaskSelectionStore{w.WorkItems, entered, release}
				go func() { _, err := workspace.SelectTaskSolution(w, reselect); first <- err }()
				waitTaskSpecBarrier(t, entered, first)
				go func() { _, err := SubmitStart(d, input); second <- err }()
			} else {
				d.Store = barrierTaskStartStore{d.Store, entered, release}
				go func() { _, err := SubmitStart(d, input); first <- err }()
				waitTaskSpecBarrier(t, entered, first)
				go func() { _, err := workspace.SelectTaskSolution(w, reselect); second <- err }()
			}
			select {
			case err := <-second:
				close(release)
				<-first
				t.Fatalf("operation bypassed cooperating lock: %v", err)
			case <-time.After(60 * time.Millisecond):
			}
			close(release)
			if err := <-first; err != nil {
				t.Fatal(err)
			}
			err := <-second
			if selectFirst && err == nil {
				t.Fatal("old tuple started after a newer choice")
			}
			if !selectFirst && err != nil {
				t.Fatal(err)
			}
			snapshot, err := originalStore.ReadByLocator(input.HandoffLocator)
			if err != nil {
				t.Fatal(err)
			}
			if selectFirst {
				if snapshot.Start != nil {
					t.Fatal("rejected start occupied accepted slot")
				}
				return
			}
			if snapshot.Start == nil || taskSpecVersion(snapshot.Start.Value) != 2 {
				t.Fatal("missing accepted @2 start")
			}
			if err = validateStoredAccepted(snapshot.Start.Value, "start", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret); err != nil {
				t.Fatal(err)
			}
			d.Store = originalStore
			retry, err := SubmitStart(d, input)
			if err != nil || retry.Created || retry.SHA256 != snapshot.Start.SHA256 {
				t.Fatalf("later choice rewrote accepted start: %#v %v", retry, err)
			}
			observation, _ := objectMember(snapshot.Start.Value, "task_spec_observation")
			forged := replaceObjectMember(snapshot.Start.Value, "task_spec_observation", replaceObjectMember(observation.(canonicaljson.Object), "validated_at_utc", "2026-09-29T12:00:01Z"))
			if err = validateStoredAccepted(forged, "start", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret); err == nil {
				t.Fatal("proof did not bind service observation")
			}
		})
	}
}

func waitTaskSpecBarrier(t *testing.T, entered <-chan struct{}, done <-chan error) {
	t.Helper()
	select {
	case <-entered:
	case err := <-done:
		t.Fatalf("operation failed before lock barrier: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("operation did not reach lock barrier")
	}
}
