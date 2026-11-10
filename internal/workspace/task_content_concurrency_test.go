package workspace

import (
	"path/filepath"
	"testing"
)

func TestTaskContentConcurrentCASPublishesExactlyOneRevision(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	a := contentFixtureProblem(t, f, "fixture/a")
	b := contentFixtureProblem(t, f, "fixture/b")
	start := make(chan struct{})
	done := make(chan error, 2)
	for _, in := range []TaskContentInput{a, b} {
		go func(in TaskContentInput) { <-start; _, err := RecordTaskProblem(f.dependencies, in); done <- err }(in)
	}
	close(start)
	wins := 0
	for i := 0; i < 2; i++ {
		if <-done == nil {
			wins++
		}
	}
	r, err := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if err != nil || wins != 1 || len(r.TaskProblemRevisions) != 2 {
		t.Fatalf("CAS wins=%d registry=%#v error=%v", wins, r, err)
	}
}
