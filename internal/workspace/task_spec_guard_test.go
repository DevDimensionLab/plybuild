package workspace

import (
	"bytes"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"strings"
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
	if err != nil {
		t.Fatalf("current observation must not revoke an intact historical input: %v", err)
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
			if (err == nil) != (name == "valid" || name == "null") {
				t.Fatalf("basis %s: %v", name, err)
			}
		})
	}
}

func TestTaskSpecTargetRejectsRelevantMissingPathsAndOwnershipConflicts(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root, target := filepath.Dir(f.wrapper), filepath.Join(f.wrapper, "task")
	if _, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: target, ExpectedParentOID: f.oid}); err != nil {
		t.Fatal(err)
	}
	basis := fixtureSelectTaskSpec(t, f)
	alias := filepath.Join(f.wrapper, "alias")
	if err := os.Symlink(target, alias); err != nil {
		t.Fatal(err)
	}
	loop := filepath.Join(f.wrapper, "loop")
	if err := os.Symlink(loop, loop); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"valid_alias", "missing_selected_target", "unresolvable_selected_target",
		"wrong_task", "wrong_spec", "unbound_selected_target", "unregistered_target",
		"selected_registration_missing", "conflicting_physical_path", "conflicting_alias",
		"conflicting_ref_common", "missing_same_ref_common", "unresolvable_other_registration",
	} {
		t.Run(name, func(t *testing.T) {
			err := WithTaskSpecSnapshot(f.dependencies, root, func(s *TaskSpecSession) error {
				b, locator := basis, target
				basisPtr := &b
				task, index := findTask(s.Registry, "task")
				selected := *task
				binding := *task.Worktree
				selected.Worktree = &binding
				other, otherBinding := selected, binding
				other.ID, other.Worktree = "conflict", &otherBinding
				other.Worktree.Ref = "refs/heads/other"
				other.GitCommonDir = filepath.Join(f.wrapper, "other.git")
				expect := "task_spec_binding_conflict"
				switch name {
				case "valid_alias":
					locator, expect = alias, ""
				case "missing_selected_target":
					locator, expect = filepath.Join(f.wrapper, "absent"), "no such file"
				case "unresolvable_selected_target":
					locator, expect = loop, "too many links"
				case "wrong_task":
					b.TaskID = "wrong"
				case "wrong_spec":
					b.Spec.ManifestSHA256 = digestTaskBytes([]byte("different selected Spec"))
					expect = "task_spec_binding_conflict"
				case "unbound_selected_target":
					basisPtr, expect = nil, ""
				case "unregistered_target":
					locator = filepath.Join(f.wrapper, "main")
				case "selected_registration_missing":
					selected.Worktree.Locator = filepath.Join(f.wrapper, "absent")
					s.Registry.Tasks[index] = selected
					expect = "task_content_observation_unknown"
				case "conflicting_physical_path":
					s.Registry.Tasks = append(s.Registry.Tasks, other)
				case "conflicting_alias":
					other.Worktree.Locator = alias
					s.Registry.Tasks = append(s.Registry.Tasks, other)
				case "conflicting_ref_common", "missing_same_ref_common":
					other.Worktree.Ref, other.GitCommonDir = selected.Worktree.Ref, selected.GitCommonDir
					other.Worktree.Locator = filepath.Join(f.wrapper, "main")
					if name == "missing_same_ref_common" {
						other.Worktree.Locator = filepath.Join(f.wrapper, "absent")
						expect = "task_content_observation_unknown"
					}
					s.Registry.Tasks = append(s.Registry.Tasks, other)
				case "unresolvable_other_registration":
					other.Worktree.Locator = loop
					s.Registry.Tasks = append(s.Registry.Tasks, other)
					expect = "task_content_observation_unknown"
				}
				_, err := s.ValidateTarget(locator, "refs/heads/task", task.GitCommonDir, basisPtr)
				if expect == "" {
					if err != nil {
						t.Fatalf("valid physical alias rejected: %v", err)
					}
				} else if err == nil || !strings.Contains(err.Error(), expect) {
					t.Fatalf("wanted %s before startup, got %v", expect, err)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("target validation changed the registered Task history")
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
	if err != nil || recovered.Readback.Classification != "exact_effect" || recovered.Readback.GitChanged == nil || !*recovered.Readback.GitChanged || git.count() != 1 {
		t.Fatalf("optional document blocked exact effect recovery: %#v %v merges=%d", recovered, err, git.count())
	}
	v := contentFields(recovered.Readback.Value)
	r, err = f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || v["integration_result"] == nil || len(r.IntegrationResults) != 1 {
		t.Fatalf("recovery did not persist exactly one observed effect: %v", err)
	}
}

func TestTaskResultGreenSpecRequiresFunctionalRequirementCoverage(t *testing.T) {
	draft := taskResultDraft{TechnicalGate: "passed", RequiredVerifierIDs: []string{}, EvidenceArtifacts: []TaskArtifactRecord{}, AcceptedDebt: []TaskAcceptedDebtRecord{}}
	evidence := TaskHandoffEvidence{TaskSpecBasis: &TaskSpecBasis{}, TaskSpecValid: true, TaskRequirementsValid: false, ReportedOutcome: "complete", SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, EvidenceCoverageValid: true, ExpectedVerifierIDs: []string{}, VerifierResults: []TaskVerifierResultRecord{}, Review: TaskReviewRecord{Findings: []TaskReviewEntry{}, Fixes: []TaskReviewEntry{}, OpenActionableFindings: []TaskReviewEntry{}}, Artifacts: []TaskArtifactRecord{}, EvidenceGaps: []string{}}
	if err := validateTaskEvidence(draft, evidence); err == nil {
		t.Fatal("green result accepted without selected requirement coverage")
	}
}

func TestTaskIntegrationResumeIgnoresUnrelatedProblemChange(t *testing.T) {
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
	if err != nil || git.count() != 1 || len(r.IntegrationAttempts) != 1 || resumed.Readback.Classification != "exact_effect" {
		t.Fatalf("document update incorrectly vetoed the unchanged explicit Git effect: %#v %v merges=%d", resumed, err, git.count())
	}
}

func TestTaskIntegrationLegacyPlanRetainsExactVersionOneAfterRegistration(t *testing.T) {
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
	if _, err = CreateTask(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	after, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	newBytes, _ := canonicaljson.Marshal(integrationPlanCanonical(after.IntegrationAuthorities[0].Plan))
	if !bytes.Equal(oldBytes, newBytes) || after.FormatVersion != before.FormatVersion || after.IntegrationAuthorities[0].PlanSHA256 != before.IntegrationAuthorities[0].PlanSHA256 || len(after.TaskResults) != 1 || len(after.HumanQARecords) != 1 || len(after.IntegrationResults) != 1 || len(after.TaskProblemRevisions) != 0 {
		t.Fatal("unrelated registration altered preserved result, QA, integration or registry format")
	}
}

func TestTaskContentPublicationPreservesUnresolvedIntegration(t *testing.T) {
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
			if err != nil || !out.Attempted.ObjectPublicationStarted {
				t.Fatalf("explicit document publication was blocked by unrelated integration state: %#v %v", out, err)
			}
			create, _ := ParseTaskCreateInput("other", "Other fixture", "Independent fixture Task", "epic", "ply", "ply")
			if _, err = CreateTask(f.dependencies, create); err != nil {
				t.Fatal(err)
			}
			after, err := f.dependencies.WorkItems.Snapshot(f.root)
			if err != nil || len(after.IntegrationResults) != 1 || after.IntegrationResults[0].Outcome != outcome || after.IntegrationAuthorities[0].PlanSHA256 != before.IntegrationAuthorities[0].PlanSHA256 {
				t.Fatalf("migration altered unresolved history: %v", err)
			}
		})
	}
}
