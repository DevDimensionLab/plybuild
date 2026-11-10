package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

// Registry observations do not consume optional document bytes.
func TestWorkItemSnapshotDoesNotReadUnselectedDocuments(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	root := filepath.Dir(f.wrapper)
	store := newSystemWorkItemStore()
	reads := 0
	store.faults = &workItemStoreFaults{fail: func(stage string) error {
		if stage == "managed-read" {
			reads++
		}
		return nil
	}}
	if _, err := store.Snapshot(root); err != nil {
		t.Fatal(err)
	}
	if reads != 0 {
		t.Fatalf("registry observation read %d unselected document files", reads)
	}

}

func TestContentObservationOwnsBytesAndRetriesMissingSources(t *testing.T) {
	root := createProjectWorkspace(t)
	source := &TaskContentStorage{}
	raw := []byte("verified immutable content")
	digest := digestTaskBytes(raw)
	view := observeTaskContent(source)
	if _, err := view.Read(root, "objects", digest); err == nil {
		t.Fatal("accepted missing object")
	}
	if _, err := source.Publish(root, "objects", raw); err != nil {
		t.Fatal(err)
	}
	first, err := view.Read(root, "objects", digest)
	if err != nil {
		t.Fatal(err)
	}
	first[0] = '!'
	again, err := view.Read(root, "objects", digest)
	if err != nil || string(again) != string(raw) {
		t.Fatalf("aliased content: %q %v", again, err)
	}
	if _, err := view.Publish(root, "objects", raw); err == nil {
		t.Fatal("observation published")
	}
	// Another root with the same digest does not inherit verified existence.
	other := physicalPath(t, t.TempDir())
	if _, err := view.Read(other, "objects", digest); err == nil {
		t.Fatal("cross-root reuse")
	}
	// The next phase must detect corruption even if the preceding one succeeded.
	path := taskContentPath(root, "objects", digest)
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("corrupt"), 0400); err != nil {
		t.Fatal(err)
	}
	if _, err := observeTaskContent(source).Read(root, "objects", digest); err == nil {
		t.Fatal("next phase hid corruption")
	}
}

func TestWorkItemObservationKeepsUnavailableDocumentsAsSeparateFacts(t *testing.T) {
	for _, operation := range []string{"snapshot", "publish", "relock"} {
		t.Run(operation, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			goal := goalFixture(t, f, "task", "goal")
			goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
			root := filepath.Dir(f.wrapper)
			store := f.dependencies.WorkItems
			r, err := store.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			corrupt := func() {
				path := taskContentPath(root, "manifests", goal.Spec.ManifestSHA256)
				if err := os.Chmod(path, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("corrupt"), 0400); err != nil {
					t.Fatal(err)
				}
			}
			switch operation {
			case "snapshot":
				corrupt()
				_, err = store.Snapshot(root)
			case "relock":
				corrupt()
				err = store.WithLock(root, func(s WorkItemStoreSession) error { _, e := s.Snapshot(); return e })
			case "publish":
				err = store.WithLock(root, func(s WorkItemStoreSession) error {
					if _, e := s.Snapshot(); e != nil {
						t.Fatal(e)
					}
					corrupt()
					return s.Publish(r)
				})
			}
			if err != nil {
				t.Fatalf("registry operation consumed unrelated document bytes: %v", err)
			}
			if _, err := f.dependencies.TaskContent.Read(root, "manifests", goal.Spec.ManifestSHA256); err == nil {
				t.Fatal("explicit document read hid corruption")
			}
		})
	}
}
