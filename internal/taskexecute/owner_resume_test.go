package taskexecute

import (
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// r001-f05: from the return worktree, a retired or completed Run of the same
// Spec does not make the active one ambiguous; a real ambiguity names a
// distinguishing, executable choice instead of repeating the same --spec.
func TestSelectOwnerRunPrefersActiveRunAndOffersDistinctChoices(t *testing.T) {
	epic := "/w/epic"
	run := func(id, task, phase string) taskrun.OwnerRunCandidate {
		return taskrun.OwnerRunCandidate{RunID: id, TaskID: "goal", SpecID: "s", Provider: "claude", TaskWorktree: task, ReturnWorktree: epic, Phase: phase}
	}
	selected, choices, next := selectOwnerRun([]taskrun.OwnerRunCandidate{run("old", "/w/t1", "retired"), run("new", "/w/t2", "begun_work")}, epic, "s")
	if selected == nil || selected.RunID != "new" || len(choices) != 0 {
		t.Fatalf("retired history made the active Run ambiguous: %+v %+v %s", selected, choices, next)
	}
	selected, _, _ = selectOwnerRun([]taskrun.OwnerRunCandidate{run("done", "/w/t1", "completed")}, epic, "s")
	if selected == nil || selected.RunID != "done" {
		t.Fatal("an explicitly selected completed Run must still show its phase")
	}
	selected, choices, next = selectOwnerRun([]taskrun.OwnerRunCandidate{run("a", "/w/t1", "begun_work"), run("b", "/w/t2", "qualified_candidate")}, epic, "s")
	if selected != nil || len(choices) != 2 || !strings.Contains(next, "/w/t1") || !strings.Contains(next, "/w/t2") || strings.Contains(next, "--spec s") {
		t.Fatalf("ambiguity does not offer distinct executable choices: %s", next)
	}
	selected, _, _ = selectOwnerRun([]taskrun.OwnerRunCandidate{run("a", "/w/t1", "begun_work"), run("b", "/w/t2", "begun_work")}, "/w/t2", "")
	if selected == nil || selected.RunID != "b" {
		t.Fatal("the Task worktree did not select its own Run")
	}
}
