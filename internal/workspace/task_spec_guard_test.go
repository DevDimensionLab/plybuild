package workspace

import (
	"bytes"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"testing"
)

type unavailableTaskBasisGit struct {
	WorkItemGit
	TaskContentGit
}

func (g unavailableTaskBasisGit) ObserveWorktree(string) (GitWorktreeObservation, error) {
	return GitWorktreeObservation{}, errors.New("fixture Git observation unavailable")
}

func TestTaskSpecUnknownBasisDoesNotReportStale(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: filepath.Join(f.wrapper, "task"), ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	basis := fixtureSelectTaskSpec(t, f)
	f.dependencies.WorkGit = unavailableTaskBasisGit{f.dependencies.WorkGit, f.dependencies.WorkGit.(TaskContentGit)}
	got, err := ShowTaskSpec(f.dependencies, TaskContentQuery{TaskID: "task", SpecID: "solution", Revision: 1})
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(got.Value)
	if contentString(m, "selection_freshness") != "unknown" {
		t.Errorf("Spec readback reports %v selection freshness for an unknown basis", m["selection_freshness"])
	}
	if contentString(m, "target_freshness") != "unknown" || m["task_spec_binding"] != nil {
		t.Fatalf("unknown observation became a usable basis: %#v", m)
	}
	unknown := false
	for _, reason := range contentArray(m, "reasons") {
		if reason == "task_spec_basis_stale" {
			t.Fatal("unavailable Git observation was reported as known drift")
		}
		unknown = unknown || reason == "task_content_observation_unknown"
	}
	if !unknown {
		t.Fatal("missing normative unknown reason")
	}
	show, err := ShowTask(f.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	readback := contentFields(show.Content)
	if contentString(contentFields(readback["next_action"]), "kind") != "inspect_task" || contentString(contentFields(readback["selection"]), "freshness") != "unknown" {
		t.Errorf("unknown basis did not prioritize inspection: %#v", readback)
	}
	err = WithTaskSpecSnapshot(f.dependencies, filepath.Dir(f.wrapper), func(s *TaskSpecSession) error {
		_, err := s.ValidateTarget(filepath.Join(f.wrapper, "task"), "refs/heads/task", s.Registry.Tasks[0].GitCommonDir, &basis)
		return err
	})
	var contentErr *TaskContentError
	if !errors.As(err, &contentErr) || contentErr.Code != "task_content_observation_unknown" {
		t.Fatalf("start guard misclassified unknown basis: %v", err)
	}
}

func TestTaskSpecGuardRejectsUnboundAndForgedBasis(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	path := filepath.Join(f.wrapper, "task")
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: path, ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	basis := fixtureSelectTaskSpec(t, f)
	for _, name := range []string{"null", "worktree", "digest", "valid"} {
		t.Run(name, func(t *testing.T) {
			b := basis
			ptr := &b
			switch name {
			case "null":
				ptr = nil
			case "worktree":
				b.TaskWorktreeID = "wt_00000000000000000000000000000000"
			case "digest":
				b.DeliveryBindingSHA256 = digestTaskBytes([]byte("forged"))
			}
			err := WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
				paths, err := s.TaskSpecReadPaths(&basis)
				if err != nil || len(paths) != 1 || paths[0] != s.Registry.Epics[0].RepoBindings[0].Worktree.Locator {
					t.Fatalf("missing Epic read coverage: %v %v", paths, err)
				}
				_, err = s.ValidateTarget(path, "refs/heads/task", s.Registry.Tasks[0].GitCommonDir, ptr)
				return err
			})
			if (err == nil) != (name == "valid") {
				t.Fatalf("basis %s: %v", name, err)
			}
		})
	}
}

func TestTaskIntegrationRecoveryPreservesEffectAfterEvidenceDrift(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply = true
	in.Confirmation = integrationDigest(t, check.Readback)
	publishes := 0
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
		if point == "temp-open" {
			publishes++
			if publishes == 3 {
				return errInjected
			}
		}
		return nil
	}}}
	if _, err = ApplyTaskIntegration(f.dependencies, in); err == nil {
		t.Fatal("result publication fault was hidden")
	}
	if git.count() != 1 {
		t.Fatal("expected one completed merge")
	}
	f.dependencies.WorkItems = newSystemWorkItemStore()
	if _, err = RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f.workItemJourneyFixture, "fixture/p2")); err != nil {
		t.Fatal(err)
	}
	f.evidence.evidence.SchemaValid = false
	recovered, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || recovered.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("lost actual effect after later drift: %#v %v merges=%d", recovered, err, git.count())
	}
}

func TestTaskIntegrationRecoveryReportsEffectWhenContentIsMissing(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	publishes := 0
	f.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(point string) error {
		if point == "temp-open" {
			publishes++
			if publishes == 3 {
				return errInjected
			}
		}
		return nil
	}}}
	if _, err = ApplyTaskIntegration(f.dependencies, in); err == nil {
		t.Fatal("result publication fault was hidden")
	}
	f.dependencies.WorkItems = newSystemWorkItemStore()
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(taskContentPath(f.root, "manifests", r.TaskSpecRevisions[0].ManifestSHA256)); err != nil {
		t.Fatal(err)
	}
	recovered, err := ApplyTaskIntegration(f.dependencies, in)
	if err == nil || recovered.Readback.Classification != "exact_effect" || recovered.Readback.GitChanged == nil || !*recovered.Readback.GitChanged || git.count() != 1 {
		t.Fatalf("lost observed effect or hid failed persistence: %#v %v merges=%d", recovered, err, git.count())
	}
	v := contentFields(recovered.Readback.Value)
	if v["integration_result"] != nil || recovered.Readback.NextAction.Kind != "read_only_recovery_control" {
		t.Fatal("claimed unpersisted result or lost recovery action")
	}
}

func TestTaskResultGreenSpecRequiresFunctionalRequirementCoverage(t *testing.T) {
	draft := taskResultDraft{TechnicalGate: "passed", RequiredVerifierIDs: []string{}, EvidenceArtifacts: []TaskArtifactRecord{}, AcceptedDebt: []TaskAcceptedDebtRecord{}}
	evidence := TaskHandoffEvidence{TaskSpecBasis: &TaskSpecBasis{}, TaskSpecValid: true, TaskRequirementsValid: false, ReportedOutcome: "complete", SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true, ExpectedVerifierIDs: []string{}, VerifierResults: []TaskVerifierResultRecord{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, Artifacts: []TaskArtifactRecord{}, EvidenceGaps: []string{}}
	if err := validateTaskEvidence(draft, evidence); err == nil {
		t.Fatal("green result accepted without selected requirement coverage")
	}
}

func TestTaskIntegrationResumeRejectsProblemChangeBeforeAttempt(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	f.dependencies.TaskLifecycleIDs = &failFirstAttemptIDs{TaskLifecycleIDSource: f.dependencies.TaskLifecycleIDs}
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
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
	if err != nil || git.count() != 0 || len(r.IntegrationAttempts) != 0 || resumed.Readback.Classification != "conflict" {
		t.Fatalf("stale intent launched Git: %#v %v merges=%d", resumed, err, git.count())
	}
}

func TestTaskIntegrationLegacyPlanRetainsExactVersionOneAfterMigration(t *testing.T) {
	f := newIntegrationJourneyFixtureVersion(t, true)
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := MarshalTaskIntegrationReadback(checked.Readback)
	if err != nil || !bytes.Contains(encoded, []byte(`"kind":"WorkspaceTaskIntegrationPlan@1"`)) || bytes.Contains(encoded, []byte("task_spec_guard")) {
		t.Fatalf("legacy plan changed: %s %v", encoded, err)
	}
	in.Apply = true
	in.Confirmation = integrationDigest(t, checked.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	before, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	plan := before.IntegrationAuthorities[0].Plan
	oldBytes, err := canonicaljson.Marshal(integrationPlanCanonical(plan))
	if err != nil {
		t.Fatal(err)
	}
	input, _ := ParseTaskCreateInput("other", "Other fixture Task", "Independent fixture problem", "epic", "ply", "ply")
	input.RegistryUpgrade = &TaskRegistryUpgrade{FromVersion: 2, RegistrySHA256: before.RawSHA256}
	rawBefore, _ := os.ReadFile(workItemsPath(f.root))
	if _, err = CreateTask(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	after, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	newBytes, _ := canonicaljson.Marshal(integrationPlanCanonical(after.IntegrationAuthorities[0].Plan))
	if !bytes.Equal(oldBytes, newBytes) || after.IntegrationAuthorities[0].PlanSHA256 != before.IntegrationAuthorities[0].PlanSHA256 || len(after.TaskResults) != 1 || len(after.HumanQARecords) != 1 || len(after.IntegrationResults) != 1 || taskRequiresSpec(after, "task") {
		t.Fatal("migration altered legacy result/QA/integration or policy")
	}
	backup, err := f.dependencies.TaskContent.Read(f.root, "backups", before.RawSHA256)
	if err != nil || !bytes.Equal(backup, rawBefore) {
		t.Fatalf("migration backup differs: %v", err)
	}
}

func TestTaskContentMigrationPreservesUnresolvedIntegrationAndRejectsActivation(t *testing.T) {
	for _, outcome := range []string{"partial", "unknown"} {
		t.Run(outcome, func(t *testing.T) {
			f := newIntegrationJourneyFixtureVersion(t, true)
			if outcome == "partial" {
				f.dependencies.IntegrationGit = &partialAfterCommandGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			} else {
				f.dependencies.IntegrationGit = &repositoryFailureAfterCommandGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			}
			in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
			checked, err := CheckTaskIntegration(f.dependencies, in)
			if err != nil {
				t.Fatal(err)
			}
			in.Apply = true
			in.Confirmation = integrationDigest(t, checked.Readback)
			applied, err := ApplyTaskIntegration(f.dependencies, in)
			if err != nil || applied.Readback.Classification != outcome {
				t.Fatalf("fixture outcome %s: %#v %v", outcome, applied, err)
			}
			before, err := f.dependencies.WorkItems.Snapshot(f.root)
			if err != nil {
				t.Fatal(err)
			}
			problem := contentFixtureProblem(t, f.workItemJourneyFixture, "fixture/activate")
			raw, _ := os.ReadFile(problem.File)
			v, _ := canonicaljson.DecodeStrict(raw)
			m := contentFields(v)
			m["registry_upgrade"] = contentRefValue(TaskRegistryUpgrade{FromVersion: 2, RegistrySHA256: before.RawSHA256})
			raw, _ = canonicaljson.Marshal(contentObject(m))
			if err = os.WriteFile(problem.File, raw, 0600); err != nil {
				t.Fatal(err)
			}
			out, err := RecordTaskProblem(f.dependencies, problem)
			if err == nil || out.Attempted.ObjectPublicationStarted {
				t.Fatalf("activated unresolved Task: %#v %v", out, err)
			}
			create, _ := ParseTaskCreateInput("other", "Other fixture", "Independent fixture Task", "epic", "ply", "ply")
			create.RegistryUpgrade = &TaskRegistryUpgrade{FromVersion: 2, RegistrySHA256: before.RawSHA256}
			if _, err = CreateTask(f.dependencies, create); err != nil {
				t.Fatal(err)
			}
			after, err := f.dependencies.WorkItems.Snapshot(f.root)
			if err != nil || len(after.IntegrationResults) != 1 || after.IntegrationResults[0].Outcome != outcome || after.IntegrationAuthorities[0].PlanSHA256 != before.IntegrationAuthorities[0].PlanSHA256 || taskRequiresSpec(after, "task") {
				t.Fatalf("migration altered unresolved history: %v", err)
			}
		})
	}
}
