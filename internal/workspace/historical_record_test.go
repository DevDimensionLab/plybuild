package workspace

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Registers written before DP01 keep records of the removed delivery process.
// They must still open, and a write must keep those records byte for byte.
func TestRegistryWithHistoricalDeliveryRecordsOpensAndKeepsThem(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	queueSetFixture(t, f, "historical", []TaskGoalQueueEntry{})
	raw, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	historical := "task_results:\n    - id: trs_11111111111111111111111111111111\n      task_id: task\n      technical_gate: passed\n      verifier_results: []\n"
	text := string(raw)
	if !strings.Contains(text, "task_results: []\n") {
		t.Fatalf("fixture has no empty task_results section:\n%s", text)
	}
	text = strings.Replace(text, "task_results: []\n", historical, 1)
	if err := os.WriteFile(workItemsPath(root), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	r := mustQueueRegistry(t, f)
	if len(r.TaskResults) != 1 {
		t.Fatalf("historical records = %d", len(r.TaskResults))
	}
	if _, err := CreateTask(f.dependencies, TaskCreateInput{TaskID: "after", Title: "After", Description: "A later write", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(workItemsPath(root))
	if !strings.Contains(string(after), historical) {
		t.Fatalf("a later write changed historical records:\n%s", after)
	}
}

// A current authority is validated on its own facts. An irrelevant old member
// must not turn it into an unvalidated historical record.
func TestCurrentAuthorityWithLegacyKeyIsStillValidated(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	checked, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	encoded, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	authority := r.IntegrationAuthorities[0]
	forged := "sha256:" + strings.Repeat("e", 64)
	for name, mutate := range map[string]func(string) string{
		"authority": func(s string) string {
			s = withMember(s, string(authority.ID), "mode: goal_delivery")
			return strings.Replace(s, "plan_sha256: "+authority.PlanSHA256+"\n", "plan_sha256: "+forged+"\n", 1)
		},
		"authority and intent": func(s string) string {
			s = withMember(s, string(authority.ID), "mode: goal_delivery")
			s = withMember(s, string(r.IntegrationIntents[0].ID), `task_result_id: ""`)
			return strings.ReplaceAll(s, "plan_sha256: "+authority.PlanSHA256+"\n", "plan_sha256: "+forged+"\n")
		},
		"intent": func(s string) string {
			return withMember(s, string(r.IntegrationIntents[0].ID), `task_result_id: ""`)
		},
	} {
		t.Run(name, func(t *testing.T) {
			changed := mutate(string(encoded))
			if changed == string(encoded) {
				t.Fatal("mutation did not apply")
			}
			if _, err := decodeWorkItemRegistry([]byte(changed)); err == nil {
				t.Fatal("current record with a legacy member skipped validation")
			}
		})
	}
}

// withMember adds one member to the list item whose id is id.
func withMember(s, id, member string) string {
	at := strings.Index(s, "- id: "+id+"\n")
	if at < 0 {
		return s
	}
	indent := strings.Repeat(" ", at-strings.LastIndex(s[:at], "\n")-1+2)
	end := at + len("- id: "+id+"\n")
	return s[:end] + indent + member + "\n" + s[end:]
}

// A new Epic base update carries no member of the removed preparation, in its
// stored plan or in the plan its confirmation digest binds.
func TestNewEpicBaseUpdatePlanHasNoPreparationMembers(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "base", []TaskGoalQueueEntry{})
	runLocalGit(t, filepath.Join(f.wrapper, "epic"), "commit", "--allow-empty", "-m", "advance base")
	preview, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || preview.Plan == nil || preview.Confirmation == nil {
		t.Fatalf("preview = %#v %v", preview, err)
	}
	digested, err := contentCanonical(*preview.Plan)
	if err != nil || strings.Contains(string(digested), "preparation") {
		t.Fatalf("digested plan names a preparation: %s %v", digested, err)
	}
	if _, err = UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *preview.Confirmation}); err != nil {
		t.Fatal(err)
	}
	stored, _ := os.ReadFile(workItemsPath(filepath.Dir(f.wrapper)))
	if strings.Contains(string(stored), "preparation_id") {
		t.Fatalf("stored base update names a preparation:\n%s", stored)
	}
}

// pendingAuthority leaves an authority and intent without attempt or result, as
// a crash right before the merge does. With forged it then rewrites the plan
// kind to an old one.
func pendingAuthority(t *testing.T, forged bool) (goalDeliveryFixture, TaskIntegrationInput) {
	t.Helper()
	f := taskWithoutGoalFixture(t)
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil {
		t.Fatal(err)
	}
	in := f.input
	in.Confirmation = integrationDigest(t, checked.Readback)
	crashing := f.dependencies
	crashing.TaskLifecycleIDs = &failAttemptIDs{fail: true}
	if _, err := ApplyTaskIntegration(crashing, in); err == nil {
		t.Fatal("injected crash before attempt did not stop apply")
	}
	if !forged {
		return f, in
	}
	err = f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		if len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationAttempts) != 0 {
			t.Fatalf("want one pending authority: %+v", registry.IntegrationAuthorities)
		}
		inline, err := encodeWorkItemRegistry(registry)
		if err != nil {
			return err
		}
		text := strings.Replace(string(inline), "kind: WorkspaceTaskIntegrationPlan@6", "kind: WorkspaceTaskIntegrationPlan@5", 1)
		if text == string(inline) {
			t.Fatal("no pending authority plan to forge")
		}
		if registry, err = decodeWorkItemRegistry([]byte(text)); err != nil {
			return err
		}
		return session.Publish(registry)
	})
	if err != nil {
		t.Fatal(err)
	}
	return f, in
}

// An authority whose plan new code does not understand is history. Check says
// its outcome is unrecorded and must be reconciled against Git; it never offers
// an apply, and applying it anyway never runs Git, even when its plan kind was
// edited to look old.
func TestPendingAuthorityWithOldPlanKindNeverRunsGit(t *testing.T) {
	f, in := pendingAuthority(t, true)
	raw, err := os.ReadFile(workItemsPath(f.root))
	if err != nil {
		t.Fatal(err)
	}
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || checked.Readback.Classification != "blocked" || checked.Readback.RecoveryStatus != "blocked" || checked.Readback.NextAction.Kind != "resolve_blocker" || len(checked.Readback.NextAction.Argv) != 0 || !strings.Contains(checked.Readback.NextAction.Reason, "reconcile") {
		t.Fatalf("check of historical pending authority: %+v %v", checked.Readback.NextAction, err)
	}
	if b, _ := MarshalTaskIntegrationReadback(checked.Readback); !strings.Contains(string(b), `"classification":"blocked"`) || strings.Contains(string(b), "apply_confirmed_plan") {
		t.Fatalf("readback offers the old plan: %s", b)
	}
	_, err = ApplyTaskIntegration(f.dependencies, in)
	if git.count() != 0 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.input.ExpectedParentOID {
		t.Fatalf("historical authority ran Git: merges=%d err=%v", git.count(), err)
	}
	if err == nil || !strings.Contains(err.Error(), "reconcile") || strings.Contains(err.Error(), "new check") {
		t.Fatalf("historical authority apply reason: %v", err)
	}
	after, readErr := os.ReadFile(workItemsPath(f.root))
	if readErr != nil || string(after) != string(raw) {
		t.Fatalf("refused apply wrote the register: %v", readErr)
	}
}

// The same pending authority with its current plan resumes and merges once.
func TestPendingCurrentAuthorityResumesWithOneMerge(t *testing.T) {
	f, in := pendingAuthority(t, false)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	checked, err := CheckTaskIntegration(f.dependencies, f.input)
	if err != nil || checked.Readback.NextAction.Kind != "apply_confirmed_plan" {
		t.Fatalf("check of current pending authority: %+v %v", checked.Readback.NextAction, err)
	}
	applied, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || applied.Readback.Classification != "exact_effect" || git.count() != 1 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.resultOID {
		t.Fatalf("current pending authority: %v merges=%d", err, git.count())
	}
}

// A session's second publish starts from the snapshot confirmed by its first,
// as integration does when it records an attempt after its authority. History
// must survive that copy: old authorities stay historical and old sections keep
// their bytes.
func TestSecondPublishInOneSessionKeepsHistoricalRecords(t *testing.T) {
	f := taskWithoutGoalFixture(t)
	expectIntegratedOnce(t, f)
	historical := "task_results:\n    - id: trs_11111111111111111111111111111111\n      task_id: task\n      technical_gate: passed\n      verifier_results: []\n"
	err := f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		inline, err := encodeWorkItemRegistry(registry)
		if err != nil {
			return err
		}
		text := strings.Replace(string(inline), "kind: WorkspaceTaskIntegrationPlan@6", "kind: WorkspaceTaskIntegrationPlan@5", 1)
		text = strings.Replace(text, "task_results: []\n", historical, 1)
		text = withMember(text, string(registry.IntegrationIntents[0].ID), `task_result_id: ""`)
		if registry, err = decodeWorkItemRegistry([]byte(text)); err != nil {
			return err
		}
		return session.Publish(registry)
	})
	if err != nil {
		t.Fatal(err)
	}
	err = f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		for publish := 1; publish <= 2; publish++ {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			if len(registry.IntegrationAuthorities) != 1 || registry.IntegrationAuthorities[0].historical == nil || len(registry.IntegrationIntents) != 1 || registry.IntegrationIntents[0].historical == nil {
				t.Fatalf("publish %d: snapshot lost the historical authority or intent", publish)
			}
			if _, err := publishWorkItemRegistryRecover(session, registry); err != nil {
				t.Fatalf("publish %d: %v", publish, err)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	after := mustQueueRegistry(t, f.workItemJourneyFixture)
	encoded, err := encodeWorkItemRegistry(after)
	kept := regexp.MustCompile(`task_results:\n(\s+)- id: trs_11111111111111111111111111111111\n\s+task_id: task\n\s+technical_gate: passed\n\s+verifier_results: \[\]\n`)
	if err != nil || len(after.TaskResults) != 1 || !kept.Match(encoded) {
		t.Fatalf("historical section changed after two publishes: %v\n%s", err, encoded)
	}
}
