package workspace

import (
	"errors"
	"testing"
)

type unavailableTaskBasisGit struct {
	WorkItemGit
	TaskContentGit
}

func (g unavailableTaskBasisGit) ObserveWorktree(string) (GitWorktreeObservation, error) {
	return GitWorktreeObservation{}, errors.New("fixture Git observation unavailable")
}

func TestTaskIntegrationResumeIgnoresUnrelatedProblemChange(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	f.dependencies.TaskLifecycleIDs = &failFirstAttemptIDs{TaskLifecycleIDSource: f.dependencies.TaskLifecycleIDs}
	in := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply = true
	in.Confirmation = integrationDigest(t, checked.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err == nil {
		t.Fatal("attempt-ID fault was hidden")
	}
	if _, err = RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f.workItemJourneyFixture, "fixture/p2")); err != nil {
		t.Fatal(err)
	}
	resumed, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || git.count() != 1 || len(r.IntegrationAttempts) != 1 || resumed.Readback.Classification != "exact_effect" {
		t.Fatalf("document update incorrectly vetoed the unchanged explicit Git effect: %#v %v merges=%d", resumed, err, git.count())
	}
}
