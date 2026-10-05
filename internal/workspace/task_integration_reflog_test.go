package workspace

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTaskIntegrationReflogStorageCompatibility(t *testing.T) {
	for _, scenario := range []struct {
		name          string
		legacy        bool
		priorEpicWork bool
		format        int
		entries       int
	}{
		{"current initial", false, false, 3, 1},
		{"current prior work", false, true, 3, 2},
		{"legacy initial", true, false, 2, 1},
		{"legacy prior work", true, true, 2, 2},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			fixture := newIntegrationJourneyFixtureHistory(t, scenario.legacy, scenario.priorEpicWork)
			reflog, err := fixture.dependencies.IntegrationGit.ObserveParentReflog(filepath.Join(fixture.wrapper, "epic"), "refs/heads/epic", 2)
			if err != nil || len(reflog) != scenario.entries || reflog[0].OID != fixture.oid {
				t.Fatalf("initial history = %#v, %v", reflog, err)
			}
			input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
			checked, err := CheckTaskIntegration(fixture.dependencies, input)
			if err != nil {
				t.Fatal(err)
			}
			input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
			if applied, err := ApplyTaskIntegration(fixture.dependencies, input); err != nil || applied.Readback.Classification != "exact_effect" {
				t.Fatalf("apply = %#v, %v", applied, err)
			}
			stored, err := os.ReadFile(workItemsPath(fixture.root))
			if err != nil {
				t.Fatal(err)
			}
			registry, err := decodeWorkItemRegistry(stored)
			if err != nil {
				t.Fatal(err)
			}
			plan := registry.IntegrationAuthorities[0].Plan
			if registry.FormatVersion != scenario.format || plan.FormatVersion != 1 || plan.Format != "json" || plan.Canonicalization != "RFC8785" || !reflect.DeepEqual(plan.ObservedReflog, reflog) {
				t.Fatalf("stored plan changed format or history: %#v", plan)
			}
			encoded, err := encodeWorkItemRegistry(registry)
			if err != nil || !bytes.Equal(encoded, stored) {
				t.Fatalf("stored registry did not round-trip unchanged: %v", err)
			}
			fixture.dependencies.WorkItems = newSystemWorkItemStore()
			shown, err := ShowTask(fixture.dependencies, "task")
			if err != nil || shown.Integration == nil || shown.Integration.Classification != "exact_effect" {
				t.Fatalf("fresh reader = %#v, %v", shown, err)
			}
		})
	}
}

type oneEntryAfterMergeGit struct {
	TaskIntegrationGit
	merges int
}

func (g *oneEntryAfterMergeGit) MergeFastForward(input IntegrationMergeInput) GitCommandOutcome {
	g.merges++
	return g.TaskIntegrationGit.MergeFastForward(input)
}

func (g *oneEntryAfterMergeGit) ObserveParentReflog(path, ref string, limit int) ([]IntegrationReflogEntry, error) {
	entries, err := g.TaskIntegrationGit.ObserveParentReflog(path, ref, limit)
	if g.merges > 0 && err == nil {
		return entries[:1], nil
	}
	return entries, err
}

func TestTaskIntegrationOnePostEffectEntryRequiresRecovery(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	git := &oneEntryAfterMergeGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	for i := 0; i < 2; i++ {
		applied, err := ApplyTaskIntegration(fixture.dependencies, input)
		if err != nil || applied.Readback.Classification != "partial" || applied.Readback.GitChanged != nil || applied.Readback.NextAction.Kind != "read_only_recovery_control" || git.merges != 1 {
			t.Fatalf("apply/replay = %#v, %v, merges = %d", applied, err, git.merges)
		}
	}
	registry, err := newSystemWorkItemStore().Snapshot(fixture.root)
	if err != nil || len(registry.IntegrationResults) != 1 || registry.IntegrationResults[0].Outcome != "partial" || len(registry.IntegrationResults[0].Reflog) != 1 {
		t.Fatalf("persisted recovery = %#v, %v", registry, err)
	}
	if got := runLocalGit(t, filepath.Join(fixture.wrapper, "epic"), "rev-parse", "HEAD"); got != fixture.resultOID {
		t.Fatal("fixture did not actually perform the fast-forward")
	}
}

func TestTaskIntegrationStoredExactEffectRequiresReflogProof(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err != nil {
		t.Fatal(err)
	}
	stored, err := os.ReadFile(workItemsPath(fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*IntegrationResult){
		"one entry": func(r *IntegrationResult) { r.Reflog = r.Reflog[:1] },
		"wrong attempt": func(r *IntegrationResult) {
			r.Reflog[0].Action = "ply-workspace-task-integrate:iat_ffffffffffffffffffffffffffffffff: Fast-forward"
		},
		"wrong before": func(r *IntegrationResult) { r.Reflog[1].OID = r.Reflog[0].OID },
		"wrong result": func(r *IntegrationResult) { r.Reflog[0].OID = r.Reflog[1].OID },
	} {
		t.Run(name, func(t *testing.T) {
			registry, err := decodeWorkItemRegistry(stored)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&registry.IntegrationResults[0])
			if _, err := encodeWorkItemRegistry(registry); err == nil || !strings.Contains(err.Error(), "does not prove its exact effect") {
				t.Fatalf("writer accepted false exact-effect proof: %v", err)
			}
			// Simulate untrusted disk bytes without invoking the validating writer.
			malformed, err := yaml.Marshal(registry)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := decodeWorkItemRegistry(malformed); err == nil || !strings.Contains(err.Error(), "does not prove its exact effect") {
				t.Fatalf("reader accepted false exact-effect proof: %v", err)
			}
		})
	}
}

func TestTaskIntegrationReflogObservationRejectsMissingOrMalformedHistory(t *testing.T) {
	for name, output := range map[string]GitCommandOutcome{
		"empty":          {},
		"unavailable":    {Err: errors.New("missing reflog")},
		"invalid oid":    {Stdout: []byte("bad\x00refs/heads/epic@{0}\x00branch: Created from HEAD\x00")},
		"missing action": {Stdout: []byte(strings.Repeat("a", 40) + "\x00refs/heads/epic@{0}\x00")},
	} {
		t.Run(name, func(t *testing.T) {
			git := &systemTaskIntegrationGit{run: func(argv, env []string) GitCommandOutcome {
				if containsExact(argv, "exists") && name != "unavailable" {
					return GitCommandOutcome{}
				}
				return output
			}}
			if entries, err := git.ObserveParentReflog("/repo", "refs/heads/epic", 2); err == nil || entries != nil {
				t.Fatalf("invalid history produced evidence: %#v, %v", entries, err)
			}
		})
	}
}
