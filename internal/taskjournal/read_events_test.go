package taskjournal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type noBatchGit struct{ workspace.WorkItemGit }

func (noBatchGit) ObserveWorktree(string) (workspace.GitWorktreeObservation, error) {
	panic("workspace activity must not observe Git")
}

func eventBatchInput(t *testing.T, f testFixture) EventBatchInput {
	t.Helper()
	b, err := workspace.ReadTaskJournalBasis(f.s.Workspace, "task")
	if err != nil {
		t.Fatal(err)
	}
	p, r, err := f.s.Workspace.Projects.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	l, err := workspace.ReadWorkItemLifecycle(f.s.Workspace, f.root)
	if err != nil {
		t.Fatal(err)
	}
	return EventBatchInput{b.Workspace, workspace.ProjectSnapshot{Projects: p, Repos: r}, b.Registry, l}
}

func TestEventBatchMatchesJournalEventsWithoutGitOrWrites(t *testing.T) {
	f := fixture(t)
	f.append(t, f.input)
	single := f.show(t)
	in := eventBatchInput(t, f)
	f.s.Workspace.WorkGit = noBatchGit{}
	before := nativeFingerprint(t, f.root)
	b, err := ReadEventBatch(f.s.Workspace, in)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := b.Tasks["task"]
	if !ok {
		t.Fatal("Task absent from batch")
	}
	if !reflect.DeepEqual(single.Events, got.Events) {
		t.Fatalf("event identities/content differ\nsingle: %+v\nbatch: %+v", single.Events, got.Events)
	}
	if nativeFingerprint(t, f.root) != before {
		t.Fatal("batch read wrote native sources")
	}
}

func TestEventBatchKeepsTaskAfterMissingContributionArtifact(t *testing.T) {
	f := fixture(t)
	f.append(t, f.input)
	if err := os.Remove(f.input.Sources[0].Locator); err != nil {
		t.Fatal(err)
	}
	in := eventBatchInput(t, f)
	b, err := ReadEventBatch(f.s.Workspace, in)
	if err != nil || len(b.Tasks["task"].Events) == 0 || b.Tasks["task"].Coverage.State != "partial" {
		t.Fatalf("batch=%+v err=%v", b, err)
	}
}

func TestEventBatchInvalidRunForOtherTaskPreservesValidPreparation(t *testing.T) {
	f := fixture(t)
	if _, err := workspace.CreateTask(f.s.Workspace, workspace.TaskCreateInput{TaskID: "a-task", Title: "Earlier Task", Description: "Has no preparation", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
		t.Fatal(err)
	}
	p := preparedBatchFixture(t, f)
	before := f.show(t)
	var expected *Event
	for i := range before.Events {
		if before.Events[i].Type == "preparation" {
			expected = &before.Events[i]
		}
	}
	if expected == nil {
		t.Fatal("fixture did not produce a valid native preparation event")
	}
	// An undecodable handoff cannot identify its Task. It references the valid
	// sibling's preparation, so projecting a-task must reject that binding only.
	bad := taskrun.Request{RequestKey: "invalid-sibling", WorkspaceRoot: f.root, PreparationID: p.ID}
	dir := filepath.Join(f.root, ".ply", "task-runs", "v1", "runs", taskrun.RunID(bad.RequestKey))
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	b, err := Canonical(bad)
	if err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(dir, "request.json"), b)
	in := eventBatchInput(t, f)
	if len(in.Registry.Tasks) != 2 || in.Registry.Tasks[0].ID != "a-task" {
		t.Fatal("fixture does not project the unprepared Task first")
	}
	out, err := ReadEventBatch(f.s.Workspace, in)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range out.Tasks["task"].Events {
		if event.EventID == expected.EventID {
			if !reflect.DeepEqual(event, *expected) {
				t.Fatalf("valid sibling preparation changed: %+v", event)
			}
			return
		}
	}
	t.Fatalf("invalid sibling request hid a valid preparation: %+v", out.Tasks["task"].Coverage.Reasons)
}

func preparedBatchFixture(t *testing.T, f testFixture) workspace.TaskPreparation {
	t.Helper()
	d := f.s.Workspace
	r, err := d.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	var problem workspace.TaskRevisionRef
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == "task" {
			problem = workspace.TaskRevisionRef{Revision: p.Revision, ManifestSHA256: p.ManifestSHA256}
		}
	}
	doc, body := filepath.Join(f.root, "goal.md"), []byte("# Goal\n\nPreserve the journal event for this prepared Task.\n")
	write(t, doc, body)
	draft := map[string]any{
		"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "publication_key": "batch/goal", "task_id": "task", "registry_upgrade": nil,
		"format": "json", "format_version": 1, "canonicalization": "RFC8785",
		"recorder": map[string]any{"actor_claim": "synthetic fixture", "control_surface": "test", "recorded_at_utc": "2026-10-05T00:00:00Z"},
		"spec_id":  "batch-goal", "expected_previous": nil, "problem": problem, "title": "Batch projection fixture", "contract_kind": "goal",
		"objective":    "Preserve independent historical Task events when another run is damaged.",
		"documents":    []any{map[string]any{"id": "goal", "source": map[string]any{"kind": "file", "locator": doc, "sha256": "sha256:" + hash(body), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}},
		"requirements": []any{map[string]any{"id": "independence", "acceptance": "A valid preparation remains visible."}},
		"design":       []any{map[string]any{"document_id": "goal", "section": nil}}, "constraints": []string{"No provider process is started."},
		"executor": map[string]any{"provider": "codex", "model": "synthetic-model", "effort": "medium"}, "change_reason": "Create a preparation fixture.",
	}
	spec, err := workspace.RecordTaskSpec(d, workspace.TaskContentInput{TaskID: "task", File: f.file(t, draft)})
	if err != nil {
		t.Fatal(err)
	}
	goal := workspace.TaskGoalRef{SpecID: *spec.OutcomeRef.SpecID, Spec: workspace.TaskRevisionRef{Revision: *spec.OutcomeRef.Revision, ManifestSHA256: spec.OutcomeRef.ManifestSHA256}}
	r, err = d.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	queue := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "batch/queue", ProjectID: "project", RepoID: "repo", EpicID: "epic", Entries: []workspace.TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "synthetic fixture", ControlSurface: "test", RecordedAtUTC: "2026-10-05T00:00:00Z"}, RegistryUpgrade: &workspace.TaskRegistryUpgrade{RegistrySHA256: r.RawSHA256, FromVersion: r.FormatVersion}}
	if _, err := workspace.SetTaskQueue(d, f.file(t, queue)); err != nil {
		t.Fatal(err)
	}
	in := workspace.TaskGoalExecuteInput{Target: workspace.QueueTargetInput{ProjectID: "project", RepoID: "repo", EpicID: "epic"}, Next: true}
	preview, err := workspace.PreviewTaskGoalExecution(d, in)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := workspace.PrepareTaskGoalExecution(d, in, preview.Confirmation, workspace.QueueHumanDecision{ActorClaim: "synthetic fixture only", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "human_cli", Statement: "Synthetic fixture, not human product approval."})
	if err != nil || prepared.Preparation.Preparation == nil {
		t.Fatalf("preparation fixture: %+v, %v", prepared, err)
	}
	return *prepared.Preparation.Preparation
}
