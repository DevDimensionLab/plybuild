package workflowhandoff

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type stagingPublicationStore struct {
	workspace.WorkItemStore
	fail bool
}

type stagingPublicationSession struct {
	workspace.WorkItemStoreSession
	fail *bool
}

func (s *stagingPublicationStore) WithLock(root string, f func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		return f(stagingPublicationSession{session, &s.fail})
	})
}

func (s stagingPublicationSession) Publish(r workspace.WorkItemRegistry) error {
	if *s.fail && len(r.TaskResults) != 0 {
		return errors.New("fixture TaskResult publication interrupted")
	}
	return s.WorkItemStoreSession.Publish(r)
}

func TestCandidateStagingReleasedOnlyAfterDurableTaskResult(t *testing.T) {
	f := prepareDeliveryFixture(t)
	input := f.candidate(t, "staging")
	store := &stagingPublicationStore{WorkItemStore: f.d.TaskWorkspace.WorkItems, fail: true}
	f.d.TaskWorkspace.WorkItems = store
	interrupted, err := QualifyDeliveryCandidate(f.d, input)
	if err == nil || !strings.Contains(err.Error(), "publication interrupted") {
		t.Fatalf("fixture did not reach the TaskResult publication boundary: %v", err)
	}
	control, err := f.d.Store.ReadByLocator(interrupted.Handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	staging := filepath.Join(control.Handoff.ReplyRoot, "staging")
	before, err := os.ReadDir(staging)
	if err != nil || len(before) != 6 {
		t.Fatalf("unpublished candidate lost recoverable staging: %d %v", len(before), err)
	}
	sentinel := filepath.Join(staging, "human-qa-unrelated")
	if err := os.WriteFile(sentinel, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	store.fail = false
	qualified, err := QualifyDeliveryCandidate(f.d, input)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range qualified.TaskResult.Artifacts {
		if _, err := os.Stat(filepath.Join(staging, a.ArtifactID)); !os.IsNotExist(err) {
			t.Fatalf("durably published artifact left another staging copy: %s", a.ArtifactID)
		}
		b, err := os.ReadFile(a.Locator)
		if err != nil || digestBytes(b) != a.SHA256 || int64(len(b)) != a.SizeBytes {
			t.Fatalf("canonical artifact disappeared after cleanup: %s %v", a.ArtifactID, err)
		}
	}
	repeated, err := QualifyDeliveryCandidate(f.d, input)
	if err != nil || repeated.TaskResult.ID != qualified.TaskResult.ID {
		t.Fatalf("completed qualification could not resume: %v", err)
	}
	if b, err := os.ReadFile(sentinel); err != nil || string(b) != "keep" {
		t.Fatal("cleanup touched unrelated staging")
	}
	r, err := f.w.WorkItems.Snapshot(f.root)
	if err != nil || len(r.TaskResults) != 1 {
		t.Fatalf("qualification replay published another TaskResult: %v", err)
	}
}

type stagingRemoveFailure struct{ FileSystem }

func (s stagingRemoveFailure) Remove(path string) error {
	if filepath.Base(filepath.Dir(path)) == "staging" {
		return errors.New("fixture cleanup unavailable")
	}
	return s.FileSystem.Remove(path)
}

func TestCandidateStagingCleanupFailureDoesNotUndoPublication(t *testing.T) {
	f := prepareDeliveryFixture(t)
	input := f.candidate(t, "cleanup-failure")
	files := f.d.Files
	f.d.Files = stagingRemoveFailure{files}
	qualified, err := QualifyDeliveryCandidate(f.d, input)
	if err != nil {
		t.Fatalf("cleanup failure invalidated a published TaskResult: %v", err)
	}
	f.d.Files = files
	repeated, err := QualifyDeliveryCandidate(f.d, input)
	if err != nil || repeated.TaskResult.ID != qualified.TaskResult.ID {
		t.Fatalf("cleanup retry lost published candidate: %v", err)
	}
	control, err := f.d.Store.ReadByLocator(qualified.Handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	staging := filepath.Join(control.Handoff.ReplyRoot, "staging")
	remaining, err := os.ReadDir(staging)
	if err != nil || len(remaining) != 0 {
		t.Fatalf("retry did not clear only the published artifact copies: %d %v", len(remaining), err)
	}
}
