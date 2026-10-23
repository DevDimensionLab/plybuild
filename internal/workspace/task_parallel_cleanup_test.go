package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelCleanupDefersSharedGitMetadata(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	if _, e := CreateTask(f.dependencies, TaskCreateInput{TaskID: "task-b", Title: "B", Description: "Retain B while A closes", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	queueSetFixture(t, f.workItemJourneyFixture, "upgrade-before-b", []QueueEntry{})
	base, e := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *base.Confirmation}); e != nil {
		t.Fatal(e)
	}
	sel := queueSelect(t, f.workItemJourneyFixture, "task-b", "cleanup-b", nil)
	queueSetFixture(t, f.workItemJourneyFixture, "queue-b", []QueueEntry{{TaskID: "task-b", Selection: &sel}})
	b := queuePrepareFixture(t, f.workItemJourneyFixture, true)
	config := filepath.Join(b.Plan.Target.GitCommonDir, "config")
	before, e := os.Stat(config)
	if e != nil {
		t.Fatal(e)
	}
	preview, e := PreviewTaskCloseout(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	if preview.Ready || !strings.Contains(strings.Join(preview.Reasons, " "), "parallel_participant_active") {
		t.Fatalf("cleanup did not defer while B depends on shared config: %+v", preview)
	}
	if _, e = ApplyTaskCloseout(f.dependencies, in, preview.PlanSHA256); e == nil {
		t.Fatal("cleanup changed shared Git metadata while B is active")
	}
	after, e := os.Stat(config)
	if e != nil || !os.SameFile(before, after) {
		t.Fatal("deferred cleanup replaced shared config identity")
	}
	if _, e = os.Stat(f.taskPath); e != nil {
		t.Fatal("deferred cleanup removed A")
	}
}

func TestParallelControlledGitConfigIdentity(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	config := filepath.Join(f.wrapper, "main", ".git", "config")
	before, err := os.Stat(config)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !preview.Ready {
		t.Fatalf("isolated cleanup preview: %+v %v", preview, err)
	}
	if _, err = ApplyTaskCloseout(f.dependencies, in, preview.PlanSHA256); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(config)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("isolated native cleanup config observation: same_inode=%t same_bytes=%t; no live provider trust behavior is inferred", os.SameFile(before, after), string(raw) == string(updated))
}
