package workspace

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWorkItemRegistryEncodingIsDeterministicAndStrict(t *testing.T) {
	registry := emptyWorkItemRegistry()
	got, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	want := "format_version: 2\nepics: []\ntasks: []\nworktree_operations: []\ntask_results: []\nhuman_qa_records: []\nintegration_authorities: []\nintegration_intents: []\nintegration_attempts: []\nintegration_results: []\n"
	if string(got) != want {
		t.Fatalf("registry bytes:\n%s\nwant:\n%s", got, want)
	}
	decoded, err := decodeWorkItemRegistry(got)
	if err != nil || decoded.FormatVersion != 2 {
		t.Fatalf("decode = %#v, %v", decoded, err)
	}
	for name, invalid := range map[string]string{
		"unknown field":   strings.Replace(want, "format_version: 2", "format_version: 2\nunknown: true", 1),
		"wrong version":   strings.Replace(want, "format_version: 2", "format_version: 3", 1),
		"missing list":    strings.Replace(want, "epics: []\n", "", 1),
		"missing v2 list": strings.Replace(want, "task_results: []\n", "", 1),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeWorkItemRegistry([]byte(invalid)); err == nil {
				t.Fatalf("invalid registry accepted:\n%s", invalid)
			}
		})
	}
}

func TestWorkItemRegistryDecodesAndPreservesFormatOne(t *testing.T) {
	legacy := []byte("format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n")
	registry, err := decodeWorkItemRegistry(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if registry.FormatVersion != 1 || registry.TaskResults != nil || registry.HumanQARecords != nil || registry.IntegrationAuthorities != nil || registry.IntegrationIntents != nil || registry.IntegrationAttempts != nil || registry.IntegrationResults != nil {
		t.Fatalf("legacy registry changed shape: %#v", registry)
	}
	encoded, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != string(legacy) {
		t.Fatalf("format-1 round trip:\n%s\nwant:\n%s", encoded, legacy)
	}

	withLifecycleCollection := append(append([]byte(nil), legacy...), []byte("task_results: []\n")...)
	if _, err := decodeWorkItemRegistry(withLifecycleCollection); err == nil {
		t.Fatal("format-1 registry with a lifecycle collection was accepted")
	}
}

func TestWorkItemRegistryRejectsInvalidObservationShapes(t *testing.T) {
	fixture := newWorkItemJourneyFixture(t)
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err != nil {
		t.Fatal(err)
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*WorkItemRegistry){
		"unknown target kind": func(value *WorkItemRegistry) { value.WorktreeOperations[0].LastObservation.TargetKind = "directory" },
		"malformed observed OID": func(value *WorkItemRegistry) {
			value.WorktreeOperations[0].LastObservation.TargetOID = stringPointer("abc")
		},
		"readback-only reason": func(value *WorkItemRegistry) {
			value.WorktreeOperations[0].LastObservation.Reasons = []string{"project_binding_unknown"}
		},
		"mismatched exact effect": func(value *WorkItemRegistry) {
			value.WorktreeOperations[0].LastObservation.TargetRef = stringPointer("refs/heads/other")
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate, err := decodeWorkItemRegistry(baseline)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&candidate)
			if _, err := encodeWorkItemRegistry(candidate); err == nil {
				t.Fatalf("invalid observation accepted: %#v", candidate.WorktreeOperations[0].LastObservation)
			}
		})
	}
}

func TestWorkItemRegistryRejectsIntegrationPlanWithDifferentGitArgv(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	registry.IntegrationAuthorities[0].Plan.Effect.Argv = []string{"git", "reset", "--hard"}
	registry.IntegrationAuthorities[0].PlanSHA256 = integrationPlanDigest(registry.IntegrationAuthorities[0].Plan)
	registry.IntegrationIntents[0].PlanSHA256 = registry.IntegrationAuthorities[0].PlanSHA256
	registry.IntegrationIntents[0].IntentSHA256 = integrationIntentDigest(registry.IntegrationIntents[0])
	if _, err := encodeWorkItemRegistry(registry); err == nil {
		t.Fatal("integration plan with non-contract Git argv was accepted")
	}
}

func TestWorkItemRegistryRejectsIntegrationPlanWithDifferentRetryBinding(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	other := IntegrationResultID("ires_ffffffffffffffffffffffffffffffff")
	registry.IntegrationAuthorities[0].Plan.RetryAfterResultID = &other
	registry.IntegrationAuthorities[0].PlanSHA256 = integrationPlanDigest(registry.IntegrationAuthorities[0].Plan)
	registry.IntegrationIntents[0].PlanSHA256 = registry.IntegrationAuthorities[0].PlanSHA256
	registry.IntegrationIntents[0].IntentSHA256 = integrationIntentDigest(registry.IntegrationIntents[0])
	if _, err := encodeWorkItemRegistry(registry); err == nil {
		t.Fatal("integration plan with a different retry predecessor was accepted")
	}
}

func TestWorkItemRegistryRejectsMalformedIntegrationCommandEvidence(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	baseline := registry.IntegrationResults[0]
	for name, mutate := range map[string]func(*IntegrationResult){
		"invalid base64":          func(result *IntegrationResult) { result.Command.Stdout.Base64 = "!" },
		"captured size mismatch":  func(result *IntegrationResult) { result.Command.Stdout.CapturedBytes++ },
		"invalid launch nullform": func(result *IntegrationResult) { result.Command.ExitCode = nil },
		"invalid truncation":      func(result *IntegrationResult) { result.Command.Stdout.Truncated = !result.Command.Stdout.Truncated },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := baseline
			mutate(&candidate)
			if err := validateIntegrationResultShape(candidate); err == nil {
				t.Fatalf("malformed command evidence accepted: %#v", candidate.Command)
			}
		})
	}
}

func TestWorkItemRegistryRejectsContradictoryIntegrationRecoveryActions(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*IntegrationResult){
		"wrong recovery status": func(result *IntegrationResult) { result.RecoveryStatus = "unknown" },
		"unsafe next action": func(result *IntegrationResult) {
			result.NextAction = IntegrationNextAction{Kind: "apply_confirmed_plan", Reason: "Apply again.", Argv: []string{"git", "merge"}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate, err := decodeWorkItemRegistry(baseline)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&candidate.IntegrationResults[0])
			if _, err := encodeWorkItemRegistry(candidate); err == nil {
				t.Fatalf("contradictory result accepted: %#v", candidate.IntegrationResults[0])
			}
		})
	}
}

func TestWorkItemRegistryRejectsBrokenIntegrationObservationChain(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	otherOID := strings.Repeat("f", len(fixture.oid))
	for name, mutate := range map[string]func(*WorkItemRegistry){
		"result before differs from Attempt": func(candidate *WorkItemRegistry) {
			candidate.IntegrationResults[0].BeforeObservation.OID = otherOID
		},
		"Attempt prestate differs from plan": func(candidate *WorkItemRegistry) {
			candidate.IntegrationAttempts[0].PreObservation.OID = otherOID
			candidate.IntegrationResults[0].BeforeObservation.OID = otherOID
		},
		"exact effect inventory differs": func(candidate *WorkItemRegistry) {
			candidate.IntegrationResults[0].AfterObservation.InventoryEntries[0].OID = fixture.oid
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate, err := decodeWorkItemRegistry(baseline)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&candidate)
			if _, err := encodeWorkItemRegistry(candidate); err == nil {
				t.Fatalf("broken observation chain accepted: %#v", candidate.IntegrationResults[0])
			}
		})
	}
}

func TestWorkItemRegistryRequiresExactNestedFieldsExplicitArraysAndNulls(t *testing.T) {
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
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := encodeWorkItemRegistry(registry)
	if err != nil {
		t.Fatal(err)
	}
	baseline := string(encoded)
	for name, candidate := range map[string]string{
		"missing nullable": strings.Replace(baseline, "    retry_after_result_id: null\n", "", 1),
		"null array":       strings.Replace(baseline, "    verifier_results: []\n", "    verifier_results: null\n", 1),
		"wrong root order": strings.Replace(baseline, "task_results:", "human_qa_records:", 1),
	} {
		t.Run(name, func(t *testing.T) {
			if candidate == baseline {
				t.Fatalf("test mutation %q did not change encoded bytes", name)
			}
			if _, err := decodeWorkItemRegistry([]byte(candidate)); err == nil {
				t.Fatalf("invalid exact shape accepted:\n%s", candidate)
			}
		})
	}
}

func TestWorkItemRegistryRejectsTaskOperationStateMismatches(t *testing.T) {
	fixture := newWorkItemJourneyFixture(t)
	fixture.dependencies.WorkGit = noEffectGit{WorkItemGit: fixture.dependencies.WorkGit}
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkGitEffect) {
		t.Fatalf("creating error = %v", err)
	}
	root := filepath.Dir(fixture.wrapper)
	creating, err := fixture.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("reconciliation Task with creating operation", func(t *testing.T) {
		contents, err := encodeWorkItemRegistry(creating)
		if err != nil {
			t.Fatal(err)
		}
		invalid := strings.Replace(string(contents), "    worktree_state: creating\n", "    worktree_state: reconciliation_required\n", 1)
		if invalid == string(contents) {
			t.Fatal("Task worktree_state was not found in encoded registry")
		}
		if _, err := decodeWorkItemRegistry([]byte(invalid)); err == nil {
			t.Fatal("reconciliation_required Task with creating operation was accepted")
		}
	})

	t.Run("creating Task with reconciliation operation", func(t *testing.T) {
		reconciliation := creating
		reconciliation.Tasks = append([]TaskRecord(nil), creating.Tasks...)
		reconciliation.WorktreeOperations = append([]WorktreeOperationRecord(nil), creating.WorktreeOperations...)
		reconciliation.Tasks[0].WorktreeState = WorkItemReconciliationRequired
		reconciliation.WorktreeOperations[0].State = "reconciliation_required"
		reconciliation.WorktreeOperations[0].LastObservation = emptyCreateObservation("partial_or_unknown")
		reconciliation.WorktreeOperations[0].LastObservation.Reasons = []string{"task_worktree_reconciliation_required"}
		contents, err := encodeWorkItemRegistry(reconciliation)
		if err != nil {
			t.Fatal(err)
		}
		invalid := strings.Replace(string(contents), "    worktree_state: reconciliation_required\n", "    worktree_state: creating\n", 1)
		if invalid == string(contents) {
			t.Fatal("Task worktree_state was not found in encoded registry")
		}
		if _, err := decodeWorkItemRegistry([]byte(invalid)); err == nil {
			t.Fatal("creating Task with reconciliation_required operation was accepted")
		}
	})
}

func TestWorkItemStoreFaultsLeaveOldOrCompleteNewRegistry(t *testing.T) {
	tests := []struct {
		name, operation           string
		shortWrite, expectPublish bool
	}{
		{name: "lock open", operation: "lock-open"},
		{name: "lock", operation: "lock"},
		{name: "temp open", operation: "temp-open"},
		{name: "write", operation: "write"},
		{name: "short write", shortWrite: true},
		{name: "chmod", operation: "chmod"},
		{name: "file sync", operation: "file-sync"},
		{name: "close", operation: "close"},
		{name: "replace", operation: "replace"},
		{name: "directory sync", operation: "directory-sync", expectPublish: true},
		{name: "unlock", operation: "unlock", expectPublish: true},
		{name: "lock close", operation: "lock-close", expectPublish: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := createProjectWorkspace(t)
			healthy := newSystemWorkItemStore()
			if err := healthy.WithLock(root, func(session WorkItemStoreSession) error { return session.Publish(emptyWorkItemRegistry()) }); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			candidate := minimalEpicRegistry(root)
			faulty := &systemWorkItemStore{faults: &workItemStoreFaults{shortWrite: test.shortWrite, fail: func(operation string) error {
				if operation == test.operation {
					return errInjected
				}
				return nil
			}}}
			err = faulty.WithLock(root, func(session WorkItemStoreSession) error { return session.Publish(candidate) })
			if err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
				t.Fatalf("error = %v", err)
			}
			after, readErr := os.ReadFile(workItemsPath(root))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if test.expectPublish == reflect.DeepEqual(before, after) {
				t.Fatalf("expect publish=%t bytes equal=%t", test.expectPublish, reflect.DeepEqual(before, after))
			}
			if err := healthy.WithLock(root, func(session WorkItemStoreSession) error { return session.Publish(candidate) }); err != nil {
				t.Fatal(err)
			}
			final, err := healthy.Snapshot(root)
			if err != nil || len(final.Epics) != 1 {
				t.Fatalf("final=%#v error=%v", final, err)
			}
		})
	}
}

func TestWorkItemStoreRejectsRegistryAndLockSymlinks(t *testing.T) {
	for _, name := range []string{WorkItemsFile, workItemsLock} {
		t.Run(name, func(t *testing.T) {
			root := createProjectWorkspace(t)
			target := filepath.Join(root, "target")
			if err := os.WriteFile(target, []byte("target"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, filepath.Join(root, MarkerDirectory, name)); err != nil {
				t.Fatal(err)
			}
			store := newSystemWorkItemStore()
			if name == WorkItemsFile {
				_, err := store.Snapshot(root)
				if err == nil || !hasWorkErrorClass(err, ErrorWorkStoreConflict) {
					t.Fatalf("error=%v", err)
				}
			} else {
				err := store.WithLock(root, func(WorkItemStoreSession) error { return errors.New("must not run") })
				if err == nil || !hasWorkErrorClass(err, ErrorWorkStoreConflict) {
					t.Fatalf("error=%v", err)
				}
			}
			contents, err := os.ReadFile(target)
			if err != nil || string(contents) != "target" {
				t.Fatalf("target=%q error=%v", contents, err)
			}
		})
	}
}

func TestTaskWorktreeStoreFailuresRespectTheGitEffectBoundary(t *testing.T) {
	t.Run("intent failure starts no Git effect", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		fixture.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
			if operation == "temp-open" {
				return errInjected
			}
			return nil
		}}}
		target := filepath.Join(fixture.wrapper, "task")
		input, _ := ParseTaskWorktreeCreateInput("task", "task", target, fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
			t.Fatalf("error=%v", err)
		}
		if _, err := os.Lstat(target); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("target exists: %v", err)
		}
		if output, err := execGitOutput(filepath.Join(fixture.wrapper, "main"), "show-ref", "--verify", "refs/heads/task"); err == nil {
			t.Fatalf("source ref exists: %s", output)
		}
	})

	t.Run("ready directory sync failure recovers exact state", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		directorySyncs := 0
		faulty := &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
			if operation == "directory-sync" {
				directorySyncs++
				if directorySyncs == 2 {
					return errInjected
				}
			}
			return nil
		}}}
		fixture.dependencies.WorkItems = faulty
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
			t.Fatalf("error=%v", err)
		}
		fixture.dependencies.WorkItems = newSystemWorkItemStore()
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "already_ready" || result.Task.WorktreeState != WorkItemReady {
			t.Fatalf("retry=%#v error=%v", result, err)
		}
	})
}

func execGitOutput(directory string, arguments ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	output, err := command.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func minimalEpicRegistry(root string) WorkItemRegistry {
	oid := strings.Repeat("a", 40)
	tree := strings.Repeat("b", 40)
	registry := emptyWorkItemRegistry()
	registry.Epics = []EpicRecord{{ID: "epic", Title: "Epic", ProjectID: "ply", RepoBindings: []EpicRepoBinding{{RepoID: "ply", GitCommonDir: filepath.Join(root, "repo.git"), Worktree: EpicWorktreeBinding{ID: "wt_11111111111111111111111111111111", OwnerKind: "epic", OwnerID: "epic", Origin: "adopted", Locator: filepath.Join(root, "epic"), Ref: "refs/heads/epic", OID: oid, Tree: tree, GitCommonDir: filepath.Join(root, "repo.git")}}}}}
	return registry
}

func TestWorkItemStoreTreatsMissingRegistryAsEmpty(t *testing.T) {
	root := createProjectWorkspace(t)
	registry, err := newSystemWorkItemStore().Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if registry.FormatVersion != 2 || len(registry.Epics) != 0 || len(registry.Tasks) != 0 || len(registry.WorktreeOperations) != 0 || registry.TaskResults == nil || registry.IntegrationResults == nil {
		t.Fatalf("registry = %#v", registry)
	}
	if filepath.Base(workItemsPath(root)) != "work-items.yaml" {
		t.Fatalf("path = %s", workItemsPath(root))
	}
}
