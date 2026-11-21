package taskwork

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestGoalStartReceiptPathsListTheTaskReceiptThenRevisionsInOrder(t *testing.T) {
	work := workspace.ProjectWork{Dir: t.TempDir()}
	if paths, err := GoalStartReceiptPaths(work, "task"); err != nil || len(paths) != 0 {
		t.Fatalf("Task without receipts: %v %v", paths, err)
	}
	single := TaskGoalStartReceiptPath(work, "task")
	want := []string{single, GoalStartReceiptPath(work, "task", "spec", 2), GoalStartReceiptPath(work, "task", "spec", 10)}
	other := GoalStartReceiptPath(work, "task-2", "spec", 1)
	dir := filepath.Dir(want[1])
	for _, path := range []string{want[2], want[1], single, other, want[1] + ".lock", filepath.Join(dir, "notes.json")} {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("{}"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if paths, err := GoalStartReceiptPaths(work, "task"); err != nil || !reflect.DeepEqual(paths, want) {
		t.Fatalf("receipts %v %v, want %v", paths, err, want)
	}
	if err := os.Remove(single); err != nil {
		t.Fatal(err)
	}
	if paths, err := GoalStartReceiptPaths(work, "task"); err != nil || !reflect.DeepEqual(paths, want[1:]) {
		t.Fatalf("receipts %v %v, want %v", paths, err, want[1:])
	}
}
