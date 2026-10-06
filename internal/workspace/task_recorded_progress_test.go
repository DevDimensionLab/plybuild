package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func recordProgressFixtureResult(t *testing.T, f taskDeliveryFixture, key, gate, recorded string) TaskResultRecord {
	t.Helper()
	raw, err := os.ReadFile(f.draftPath)
	if err != nil {
		t.Fatal(err)
	}
	var draft map[string]any
	if err = json.Unmarshal(raw, &draft); err != nil {
		t.Fatal(err)
	}
	draft["publication_key"] = key
	terminalID := "res_" + digestTaskBytes([]byte(key))[7:39]
	draft["handoff"].(map[string]any)["terminal_result_id"] = terminalID
	f.reader.evidence.TerminalResultID = terminalID
	draft["technical_assessment"].(map[string]any)["gate"] = gate
	draft["recorder"].(map[string]any)["recorded_at_utc"] = recorded
	raw, err = json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.root, filepath.Base(key)+"-result.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	result, err := RecordTaskResult(f.dependencies, TaskResultRecordInput{TaskID: "task", File: path})
	if err != nil {
		t.Fatal(err)
	}
	return result.Record
}

func registeredProgressFixtureFacts(t *testing.T, d Dependencies, root string) TaskRecordedProgressFacts {
	t.Helper()
	r, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		t.Fatal(err)
	}
	return RecordedTaskProgressFacts(r)["task"]
}

func TestRecordedProgressDoesNotChooseAnOlderGreenResultOverCompetingResults(t *testing.T) {
	for _, gate := range []string{"failed", "unknown", "passed"} {
		t.Run(gate, func(t *testing.T) {
			f := newTaskDeliveryFixture(t, false)
			recordProgressFixtureResult(t, f, "progress/a", "passed", "2026-09-28T12:00:00Z")
			recordProgressFixtureResult(t, f, "progress/b", gate, "2026-09-28T13:00:00Z")
			fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
			if !fact.Ambiguous || fact.TaskResult != nil || fact.HumanQA != nil || fact.IntegrationResult != nil || fact.ResultCount != 2 {
				t.Fatalf("competing %s result was hidden behind green history: %+v", gate, fact)
			}
		})
	}
}

func applyProgressFixtureIntegration(t *testing.T, f integrationJourneyFixture, retry *IntegrationResultID) TaskIntegrationResult {
	t.Helper()
	in := TaskIntegrationInput{TaskID: "task", TaskResultID: f.resultID, HumanQARecordID: f.qaID, ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid, RetryAfterResultID: retry}
	checked, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	return applied
}

func TestRecordedProgressDoesNotApplyHistoricalIntegrationToAnotherResult(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	if applied := applyProgressFixtureIntegration(t, f, nil); applied.Readback.Classification != "exact_effect" {
		t.Fatalf("fixture did not integrate: %+v", applied.Readback)
	}
	err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		other := r.TaskResults[0]
		other.ID, other.PublicationKey = "trs_22222222222222222222222222222222", "progress/new-result"
		other.TerminalResultID = "res_22222222222222222222222222222222"
		other.TechnicalGate, other.Recorder.RecordedAtUTC = "unknown", "2026-09-29T13:00:00Z"
		r.TaskResults = append(r.TaskResults, other)
		link := r.TaskResultSpecBindings[0]
		link.TaskResultID = other.ID
		r.TaskResultSpecBindings = append(r.TaskResultSpecBindings, link)
		sortWorkRegistry(&r)
		return s.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	if !fact.Ambiguous || fact.TaskResult != nil || fact.HumanQA != nil || fact.IntegrationResult != nil || fact.IntegrationPending {
		t.Fatalf("historical exact effect claimed for competing current result: %+v", fact)
	}
}

func TestRecordedProgressPreservesOneExplicitIntegrationRetryChain(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	realGit := f.dependencies.IntegrationGit
	f.dependencies.IntegrationGit = &noEffectIntegrationGit{TaskIntegrationGit: realGit}
	if first := applyProgressFixtureIntegration(t, f, nil); first.Readback.Classification != "no_effect" {
		t.Fatalf("fixture did not record no effect: %+v", first.Readback)
	}
	before := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	prior := before.IntegrationResult.ID
	f.dependencies.IntegrationGit = realGit
	if retry := applyProgressFixtureIntegration(t, f, &prior); retry.Readback.Classification != "exact_effect" {
		t.Fatalf("fixture did not integrate retry: %+v", retry.Readback)
	}
	fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	if fact.Ambiguous || fact.TaskResult == nil || fact.TaskResult.ID != f.resultID || fact.HumanQA == nil || fact.HumanQA.ID != f.qaID || fact.IntegrationResult == nil || fact.IntegrationResult.Outcome != "exact_effect" {
		t.Fatalf("explicit retry chain lost its exact binding: %+v", fact)
	}
}

func changeProgressFixtureContent(t *testing.T, f workItemJourneyFixture, basis TaskSpecBasis, change string) {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	if change == "selected-spec" {
		selectNewProgressFixtureSpec(t, f, basis)
		return
	}
	if change == "problem" {
		if _, err := RecordTaskProblem(f.dependencies, contentFixtureProblem(t, f, "progress/problem")); err != nil {
			t.Fatal(err)
		}
		return
	}
	if change == "unselected-spec" {
		if _, err := RecordTaskSpec(f.dependencies, contentFixtureSpec(t, f, "progress/unselected-spec", &basis.Spec)); err != nil {
			t.Fatal(err)
		}
		return
	}
	name := "selection"
	if change == "assessment" {
		name = "assessment"
	}
	raw, err := os.ReadFile(filepath.Join(root, name+"-draft.json"))
	if err != nil {
		t.Fatal(err)
	}
	value, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		t.Fatal(err)
	}
	m := contentFields(value)
	m["publication_key"] = "progress/" + change
	if change == "assessment" {
		m["expected_previous_assessment"], m["outcome"] = contentRefValue(basis.Assessment), "needs_work"
	} else {
		m["expected_previous_selection"] = contentRefValue(basis.Selection)
		if change == "withdraw" {
			m["action"], m["solution"] = "withdraw", nil
		}
	}
	raw, err = canonicaljson.Marshal(contentObject(m))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "progress-content-change.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	input := TaskContentInput{TaskID: "task", File: path}
	switch change {
	case "assessment":
		_, err = AssessTaskSpec(f.dependencies, input)
	case "withdraw":
		_, err = WithdrawTaskSolution(f.dependencies, input)
	default:
		_, err = SelectTaskSolution(f.dependencies, input)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecordedProgressDoesNotReuseResultsAfterRegisteredBasisChanges(t *testing.T) {
	for _, change := range []string{"problem", "selected-spec", "withdraw", "assessment", "unselected-spec"} {
		t.Run(change, func(t *testing.T) {
			f := newTaskDeliveryFixture(t, false)
			result := recordProgressFixtureResult(t, f, "progress/result", "passed", "2026-09-28T12:00:00Z")
			changeProgressFixtureContent(t, f.workItemJourneyFixture, *f.reader.evidence.TaskSpecBasis, change)
			fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
			if change == "unselected-spec" {
				if fact.Ambiguous || fact.TaskResult == nil || fact.TaskResult.ID != result.ID {
					t.Fatalf("unselected draft replaced the explicit selection: %+v", fact)
				}
			} else if fact.TaskResult != nil || fact.HumanQA != nil || fact.IntegrationResult != nil || !fact.HasProgress || fact.ResultCount != 1 {
				t.Fatalf("historical result remained current after %s: %+v", change, fact)
			}
		})
	}
}

func selectNewProgressFixtureSpec(t *testing.T, f workItemJourneyFixture, basis TaskSpecBasis) TaskSpecBasis {
	t.Helper()
	root := filepath.Dir(f.wrapper)
	input := contentFixtureSpec(t, f, "progress/new-spec", &basis.Spec)
	mutate := func(input TaskContentInput, fn func(map[string]canonicaljson.Value)) TaskContentInput {
		t.Helper()
		raw, err := os.ReadFile(input.File)
		if err != nil {
			t.Fatal(err)
		}
		v, err := canonicaljson.DecodeStrict(raw)
		if err != nil {
			t.Fatal(err)
		}
		m := contentFields(v)
		fn(m)
		raw, err = canonicaljson.Marshal(contentObject(m))
		if err != nil {
			t.Fatal(err)
		}
		input.File = filepath.Join(root, filepath.Base(contentString(m, "publication_key"))+".json")
		if err = os.WriteFile(input.File, raw, 0600); err != nil {
			t.Fatal(err)
		}
		return input
	}
	input = mutate(input, func(m map[string]canonicaljson.Value) {
		binding := contentFields(m["implementation_basis"])
		binding["start_oid"] = runLocalGit(t, filepath.Join(f.wrapper, "task"), "rev-parse", "HEAD")
		binding["start_tree"] = runLocalGit(t, filepath.Join(f.wrapper, "task"), "rev-parse", "HEAD^{tree}")
		m["implementation_basis"] = contentObject(binding)
	})
	spec, err := RecordTaskSpec(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	ref := TaskRevisionRef{*spec.OutcomeRef.Revision, spec.OutcomeRef.ManifestSHA256}
	input = mutate(TaskContentInput{TaskID: "task", File: filepath.Join(root, "assessment-draft.json")}, func(m map[string]canonicaljson.Value) {
		m["publication_key"], m["spec"] = "progress/new-assessment", contentRefValue(ref)
	})
	assessment, err := AssessTaskSpec(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input = mutate(TaskContentInput{TaskID: "task", File: filepath.Join(root, "selection-draft.json")}, func(m map[string]canonicaljson.Value) {
		m["publication_key"], m["expected_previous_selection"] = "progress/new-selection", contentRefValue(basis.Selection)
		solution := contentFields(m["solution"])
		solution["spec"] = contentRefValue(ref)
		solution["assessment"] = contentRefValue(TaskDecisionRef{*assessment.OutcomeRef.ID, assessment.OutcomeRef.ManifestSHA256})
		m["solution"] = contentObject(solution)
	})
	selected, err := SelectTaskSolution(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	basis.Spec = ref
	basis.Assessment = TaskDecisionRef{*assessment.OutcomeRef.ID, assessment.OutcomeRef.ManifestSHA256}
	basis.Selection = TaskDecisionRef{*selected.OutcomeRef.ID, selected.OutcomeRef.ManifestSHA256}
	return basis
}

func TestRecordedProgressSelectsOnlyResultForCurrentRegisteredBasis(t *testing.T) {
	f := newTaskDeliveryFixture(t, false)
	recordProgressFixtureResult(t, f, "progress/historical", "passed", "2026-09-29T12:00:00Z")
	basis := selectNewProgressFixtureSpec(t, f.workItemJourneyFixture, *f.reader.evidence.TaskSpecBasis)
	f.reader.evidence.TaskSpecBasis = &basis
	// An explicitly bound result wins over stale history even when its recorder
	// clock is earlier; the choice must follow references, not wall-clock order.
	current := recordProgressFixtureResult(t, f, "progress/current", "unknown", "2026-09-28T12:00:00Z")
	fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	if fact.Ambiguous || fact.ResultBasisStale || fact.ResultCount != 2 || fact.TaskResult == nil || fact.TaskResult.ID != current.ID || fact.TaskResult.TechnicalGate != "unknown" || fact.HumanQA != nil || fact.IntegrationResult != nil {
		t.Fatalf("current explicit basis did not select the relevant result: %+v", fact)
	}
}

func TestRecordedProgressKeepsLegacyResultWithoutSpecBinding(t *testing.T) {
	f := newTaskDeliveryFixture(t, true)
	result := recordProgressFixtureResult(t, f, "progress/legacy", "passed", "2026-09-28T12:00:00Z")
	fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	if fact.Ambiguous || fact.ResultBasisStale || fact.TaskResult == nil || fact.TaskResult.ID != result.ID || fact.HumanQA != nil {
		t.Fatalf("legacy result lost its compatible registered semantics: %+v", fact)
	}
}

func TestRecordedProgressDoesNotApplyHistoricalIntegrationToNewProblem(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	if applied := applyProgressFixtureIntegration(t, f, nil); applied.Readback.Classification != "exact_effect" {
		t.Fatalf("fixture did not integrate: %+v", applied.Readback)
	}
	changeProgressFixtureContent(t, f.workItemJourneyFixture, *f.evidence.evidence.TaskSpecBasis, "problem")
	fact := registeredProgressFixtureFacts(t, f.dependencies, f.root)
	if fact.TaskResult != nil || fact.HumanQA != nil || fact.IntegrationResult != nil || fact.IntegrationPending || !fact.HasProgress || fact.ResultCount != 1 {
		t.Fatalf("historical integration remained current for changed problem: %+v", fact)
	}
}
