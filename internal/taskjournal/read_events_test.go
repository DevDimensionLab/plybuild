package taskjournal

import (
	"os"
	"reflect"
	"testing"

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
