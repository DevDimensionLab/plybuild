package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoalOperationUnboundSessionCannotPublish(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	goalQueueFixture(t, f, []TaskGoalQueueEntry{})
	root := filepath.Dir(f.wrapper)
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	d, _ := beginGoalOperationContent(f.dependencies)
	err = d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		return session.Publish(r)
	})
	if err == nil || !strings.Contains(err.Error(), "dependencies are not bound") {
		t.Fatalf("unbound publication accepted: %v", err)
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("unbound operation changed shared registry")
	}
}
