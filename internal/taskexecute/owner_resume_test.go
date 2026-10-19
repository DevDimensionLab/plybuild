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
	selected, choices, next := selectOwnerRun([]taskrun.OwnerRunCandidate{run("old", "/w/t1", "retired"), run("new", "/w/t2", "begun_work")}, nil, epic, "s")
	if selected == nil || selected.RunID != "new" || len(choices) != 0 {
		t.Fatalf("retired history made the active Run ambiguous: %+v %+v %s", selected, choices, next)
	}
	selected, _, _ = selectOwnerRun([]taskrun.OwnerRunCandidate{run("done", "/w/t1", "completed")}, nil, epic, "s")
	if selected == nil || selected.RunID != "done" {
		t.Fatal("an explicitly selected completed Run must still show its phase")
	}
	selected, choices, next = selectOwnerRun([]taskrun.OwnerRunCandidate{run("a", "/w/t1", "begun_work"), run("b", "/w/t2", "qualified_candidate")}, nil, epic, "s")
	if selected != nil || len(choices) != 2 || !strings.Contains(next, "/w/t1") || !strings.Contains(next, "/w/t2") || strings.Contains(next, "--spec s") {
		t.Fatalf("ambiguity does not offer distinct executable choices: %s", next)
	}
	selected, _, _ = selectOwnerRun([]taskrun.OwnerRunCandidate{run("a", "/w/t1", "begun_work"), run("b", "/w/t2", "begun_work")}, nil, "/w/t2", "")
	if selected == nil || selected.RunID != "b" {
		t.Fatal("the Task worktree did not select its own Run")
	}
}

// r004-f02: when nothing readable matches, unreadable Runs are named with
// their exact errors and a native next step instead of "no Run found";
// unrelated unreadable history does not block a readable match.
func TestSelectOwnerRunNamesUnreadableRuns(t *testing.T) {
	bad := []taskrun.OwnerRunUnreadable{{RunID: "wfr_bad", Error: "task_run_schema_invalid: missing writer_liveness"}}
	selected, _, next := selectOwnerRun(nil, bad, "/w/t1", "")
	if selected != nil || !strings.Contains(next, "wfr_bad") || !strings.Contains(next, "missing writer_liveness") || !strings.Contains(next, "workflow execute show wfr_bad") || strings.Contains(next, "no Run was found") {
		t.Fatalf("unreadable Run hidden: %s", next)
	}
	good := []taskrun.OwnerRunCandidate{{RunID: "wfr_good", SpecID: "s", TaskWorktree: "/w/t1", ReturnWorktree: "/w/epic", Phase: "begun_work"}}
	if selected, _, _ = selectOwnerRun(good, bad, "/w/t1", ""); selected == nil || selected.RunID != "wfr_good" {
		t.Fatal("unrelated unreadable history blocked a readable Task Run")
	}
}
