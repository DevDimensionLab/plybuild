package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHumanChosenTaskTargetNeedsNoSelectedSpec(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	before := mustQueueRegistry(t, f)
	if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "direct", Title: "Direct", Description: "Use the explicit source", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	after := mustQueueRegistry(t, f)
	if len(after.TaskProblemRevisions) != len(before.TaskProblemRevisions) || len(after.TaskContentPublications) != len(before.TaskContentPublications) {
		t.Fatal("Task creation generated process documents")
	}
	path := filepath.Join(f.wrapper, "direct")
	ready, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "direct", Branch: "direct", Path: path, ExpectedParentOID: f.oid})
	if err != nil {
		t.Fatal(err)
	}
	err = WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
		_, err := s.ValidateTarget(path, "refs/heads/direct", ready.Task.GitCommonDir, nil)
		return err
	})
	if err != nil {
		t.Fatalf("explicit Task target requires process paperwork: %v", err)
	}
	err = WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
		_, err := s.ValidateTarget(path, "refs/heads/wrong", ready.Task.GitCommonDir, nil)
		return err
	})
	if err == nil {
		t.Fatal("explicit target accepted a mismatched Git ref")
	}

}

func TestHumanChosenIntegrationNeedsNoQAOrTechnicalPass(t *testing.T) {
	f := newIntegrationJourneyFixtureVersion(t, true)
	err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		r.HumanQARecords = nil
		for i := range r.TaskResults {
			r.TaskResults[i].TechnicalGate = "failed"
		}
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	checked, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatalf("human-selected integration requires QA: %v", err)
	}
	if checked.Readback.Classification != "ready" {
		t.Fatalf("integration blocked by process: %#v", checked.Readback)
	}
	input.Apply = true
	input.Confirmation = checked.Readback.NextAction.Argv[len(checked.Readback.NextAction.Argv)-1]
	applied, err := ApplyTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if applied.Readback.Classification != "exact_effect" {
		t.Fatalf("integration: %#v", applied.Readback)
	}
	if got := runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD"); got != f.resultOID {
		t.Fatalf("parent = %s", got)
	}
}

func TestExplicitWorktreeDoesNotRequireUnselectedDocumentBytes(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "optional-empty-queue", []QueueEntry{})
	r := mustQueueRegistry(t, f)
	root := filepath.Dir(f.wrapper)
	if err := os.Remove(taskContentPath(root, "manifests", r.TaskProblemRevisions[0].ManifestSHA256)); err != nil {
		t.Fatal(err)
	}
	ready, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.wrapper, "task"), ExpectedParentOID: f.oid})
	if err != nil {
		t.Fatalf("unselected document blocked explicit Git target: %v", err)
	}
	if ready.Task.Worktree == nil || ready.Task.Worktree.OID != f.oid {
		t.Fatal("explicit worktree source binding was lost")
	}
	if read, err := ShowTaskProblem(f.dependencies, TaskContentQuery{TaskID: "task"}); err == nil && contentString(contentFields(read.Value), "integrity") == "valid" {
		t.Fatal("missing document was presented as readable")
	}
}
