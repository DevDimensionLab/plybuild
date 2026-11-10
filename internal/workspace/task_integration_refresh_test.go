package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type refreshFixture struct {
	integrationJourneyFixture
	epicPath, secondPath string
	a, b, c, cTree       string
}

func (f refreshFixture) input(refresh bool) TaskIntegrationInput {
	return TaskIntegrationInput{TaskID: "second", ExpectedResultOID: f.c, ExpectedParentOID: f.b, RefreshParent: refresh}
}

func newRefreshFixture(t *testing.T) refreshFixture {
	t.Helper()
	f := newIntegrationJourneyFixture(t)
	second, _ := ParseTaskCreateInput("second", "Second", "Second Task from the same base", "epic", "ply", "ply")
	if err := fixtureCreateTaskWithProblem(f.dependencies, f.root, second); err != nil {
		t.Fatal(err)
	}
	secondPath := filepath.Join(f.wrapper, "second")
	worktree, _ := ParseTaskWorktreeCreateInput("second", "second", secondPath, f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, worktree); err != nil {
		t.Fatal(err)
	}
	first := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
	if out := refreshApply(t, f.dependencies, first); out.Classification != "exact_effect" {
		t.Fatalf("first Task = %s", out.Classification)
	}
	refreshBaseUpdate(t, f)
	epicPath := filepath.Join(f.wrapper, "epic")
	b := runLocalGit(t, epicPath, "rev-parse", "HEAD")
	runLocalGit(t, secondPath, "reset", "--hard", b)
	if err := os.WriteFile(filepath.Join(secondPath, "second.txt"), []byte("second\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, secondPath, "add", "second.txt")
	runLocalGit(t, secondPath, "commit", "-m", "adapted second Task")
	c := runLocalGit(t, secondPath, "rev-parse", "HEAD")
	return refreshFixture{integrationJourneyFixture: f, epicPath: epicPath, secondPath: secondPath, a: f.oid, b: b, c: c, cTree: runLocalGit(t, secondPath, "rev-parse", "HEAD^{tree}")}
}

func refreshBaseUpdate(t *testing.T, f integrationJourneyFixture) {
	t.Helper()
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	preview, err := UpdateEpicBase(f.dependencies, EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || preview.Confirmation == nil {
		t.Fatalf("base preview = %#v err=%v", preview, err)
	}
	in := EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *preview.Confirmation}
	if r.FormatVersion < 4 {
		in.Upgrade = r.RawSHA256
	}
	if _, err = UpdateEpicBase(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
}

func refreshCheck(t *testing.T, d Dependencies, in TaskIntegrationInput) WorkspaceTaskIntegrationReadback {
	t.Helper()
	out, err := CheckTaskIntegration(d, in)
	if err != nil {
		t.Fatal(err)
	}
	return out.Readback
}

func refreshApply(t *testing.T, d Dependencies, in TaskIntegrationInput) WorkspaceTaskIntegrationReadback {
	t.Helper()
	checked := refreshCheck(t, d, in)
	if checked.Classification != "ready" {
		t.Fatalf("check = %s %#v", checked.Classification, checked.NextAction)
	}
	in.Confirmation = integrationDigest(t, checked)
	out, err := ApplyTaskIntegration(d, in)
	if err != nil {
		t.Fatal(err)
	}
	return out.Readback
}

func readbackMember(r WorkspaceTaskIntegrationReadback, name string) any {
	for _, m := range r.Value {
		if m.Name == name {
			return m.Value
		}
	}
	return nil
}

func planReasons(t *testing.T, r WorkspaceTaskIntegrationReadback) string {
	t.Helper()
	return integrationPlanFromReadback(t, r)
}

func TestRefreshParentIntegratesAdaptedTaskFromRegisteredBase(t *testing.T) {
	f := newRefreshFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	task, _ := findTask(r, "second")
	originalTask := *task
	originalPublications := slices.Clone(r.TaskContentPublications)
	before, _ := os.ReadFile(workItemsPath(f.root))

	plain := refreshCheck(t, f.dependencies, f.input(false))
	if plain.Classification == "ready" || readbackMember(plain, "kind") != "WorkspaceTaskIntegrationReadback@6" || !strings.Contains(planReasons(t, plain), "expected_oid_mismatch") || !strings.Contains(plain.NextAction.Reason, "--refresh-parent") || !strings.Contains(plain.NextAction.Reason, f.b) {
		t.Fatalf("plain check = %s %v", plain.Classification, readbackMember(plain, "kind"))
	}
	checked := refreshCheck(t, f.dependencies, f.input(true))
	plan := integrationPlanFromReadback(t, checked)
	if checked.Classification != "ready" || readbackMember(checked, "kind") != "WorkspaceTaskIntegrationReadback@7" || readbackMember(checked, "schema_version") != int64(7) || !strings.Contains(plan, `"kind":"WorkspaceTaskIntegrationPlan@7"`) || !strings.Contains(plan, `"expected_parent_oid":"`+f.b+`"`) || strings.Contains(plan, `"expected_parent_oid":"`+f.a+`"`) {
		t.Fatalf("refresh check = %s %#v\n%s", checked.Classification, checked.NextAction, plan)
	}
	digest := integrationDigest(t, checked)
	want := []string{"ply", "workspace", "task", "integrate", "second", "--refresh-parent", "--expected-result-oid", f.c, "--expected-parent-oid", f.b, "--apply", "--confirm", digest}
	if !slices.Equal(checked.NextAction.Argv, want) {
		t.Fatalf("argv = %v", checked.NextAction.Argv)
	}
	after, _ := os.ReadFile(workItemsPath(f.root))
	if string(before) != string(after) || git.count() != 0 {
		t.Fatal("check had an effect")
	}

	in := f.input(true)
	in.Confirmation = digest
	applied, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || applied.Readback.Classification != "exact_effect" || git.count() != 1 || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.c {
		t.Fatalf("apply = %s merges=%d err=%v", applied.Readback.Classification, git.count(), err)
	}
	r = mustQueueRegistry(t, f.workItemJourneyFixture)
	task, _ = findTask(r, "second")
	if !reflect.DeepEqual(*task, originalTask) || !reflect.DeepEqual(r.TaskContentPublications[:len(originalPublications)], originalPublications) {
		t.Fatal("refresh changed the Task receipt or its Goal/start history")
	}
	authority := r.IntegrationAuthorities[len(r.IntegrationAuthorities)-1]
	if authority.historical != nil || authority.Plan.Kind != "WorkspaceTaskIntegrationPlan@7" || authority.AllowedEffect.ExpectedParentOID != f.b {
		t.Fatalf("authority = %#v", authority)
	}

	replay := refreshCheck(t, f.dependencies, f.input(true))
	if replay.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("replay = %s", replay.Classification)
	}
	again, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || again.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("replay apply = %s err=%v", again.Readback.Classification, err)
	}
	shown, err := ShowTask(f.dependencies, "second")
	if err != nil || shown.Integration == nil || shown.Integration.Classification != "exact_effect" || shown.Integration.ResultOID != f.c {
		t.Fatalf("show = %#v err=%v", shown.Integration, err)
	}
	if fact := RecordedTaskProgressFacts(r)["second"]; fact.IntegrationResult == nil || fact.IntegrationResult.Outcome != "exact_effect" || fact.IntegrationPending {
		t.Fatalf("progress = %#v", fact)
	}
	runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "later member")
	later := refreshCheck(t, f.dependencies, f.input(true))
	if later.Classification != "exact_effect" || !later.HistoricalCompletion || git.count() != 1 {
		t.Fatalf("later replay = %s historical=%t", later.Classification, later.HistoricalCompletion)
	}
	if _, err := f.dependencies.WorkItems.Snapshot(f.root); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshParentWithoutNewBasisIsRejected(t *testing.T) {
	cases := map[string]struct {
		prepare func(*testing.T, *refreshFixture) TaskIntegrationInput
		reason  string
	}{
		"same original parent": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			in := f.input(true)
			in.ExpectedParentOID = f.a
			return in
		}, "refresh_parent_unchanged"},
		"stale registered base": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "unregistered Epic work")
			d := runLocalGit(t, f.epicPath, "rev-parse", "HEAD")
			runLocalGit(t, f.secondPath, "rebase", d)
			f.c = runLocalGit(t, f.secondPath, "rev-parse", "HEAD")
			in := f.input(true)
			in.ExpectedParentOID = d
			return in
		}, "refresh_parent_not_registered_base"},
		"dirty source": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			if err := os.WriteFile(filepath.Join(f.secondPath, "dirty.txt"), []byte("dirty\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			return f.input(true)
		}, "source_not_clean"},
		"moved source": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			runLocalGit(t, f.secondPath, "commit", "--allow-empty", "-m", "moved")
			return f.input(true)
		}, "source_binding_mismatch"},
		"result not descended from new parent": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			runLocalGit(t, f.secondPath, "reset", "--hard", f.a)
			runLocalGit(t, f.secondPath, "commit", "--allow-empty", "-m", "not adapted")
			f.c = runLocalGit(t, f.secondPath, "rev-parse", "HEAD")
			return f.input(true)
		}, "result_not_descendant"},
		"new parent not descended from original": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			f.dependencies.IntegrationGit = forcedResetIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit, original: f.a, parent: f.b}
			return f.input(true)
		}, "refresh_parent_not_descendant"},
		"ambiguous inventory": {func(t *testing.T, f *refreshFixture) TaskIntegrationInput {
			f.dependencies.IntegrationGit = inventoryMismatchIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			return f.input(true)
		}, "inventory_mismatch"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := newRefreshFixture(t)
			in := tc.prepare(t, &f)
			before, _ := os.ReadFile(workItemsPath(f.root))
			out := refreshCheck(t, f.dependencies, in)
			if out.Classification == "ready" || !strings.Contains(planReasons(t, out), tc.reason) {
				t.Fatalf("check = %s\n%s", out.Classification, planReasons(t, out))
			}
			after, _ := os.ReadFile(workItemsPath(f.root))
			if string(before) != string(after) {
				t.Fatal("rejected refresh wrote state")
			}
		})
	}
}

type forcedResetIntegrationGit struct {
	TaskIntegrationGit
	original, parent string
}

func (g forcedResetIntegrationGit) CheckAncestor(repo RepoRecord, ancestor, descendant string) (bool, error) {
	if ancestor == g.original && descendant == g.parent {
		return false, nil
	}
	return g.TaskIntegrationGit.CheckAncestor(repo, ancestor, descendant)
}

func TestRefreshParentRequiresSettledRegisteredBase(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	second, _ := ParseTaskCreateInput("second", "Second", "Second Task", "epic", "ply", "ply")
	if err := fixtureCreateTaskWithProblem(f.dependencies, f.root, second); err != nil {
		t.Fatal(err)
	}
	secondPath := filepath.Join(f.wrapper, "second")
	worktree, _ := ParseTaskWorktreeCreateInput("second", "second", secondPath, f.oid)
	if _, err := CreateTaskWorktree(f.dependencies, worktree); err != nil {
		t.Fatal(err)
	}
	refreshApply(t, f.dependencies, TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid})
	b := runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD")
	runLocalGit(t, secondPath, "reset", "--hard", b)
	runLocalGit(t, secondPath, "commit", "--allow-empty", "-m", "adapted")
	c := runLocalGit(t, secondPath, "rev-parse", "HEAD")
	out := refreshCheck(t, f.dependencies, TaskIntegrationInput{TaskID: "second", ExpectedResultOID: c, ExpectedParentOID: b, RefreshParent: true})
	if out.Classification == "ready" || !strings.Contains(planReasons(t, out), "refresh_parent_not_registered_base") {
		t.Fatalf("unregistered base = %s\n%s", out.Classification, planReasons(t, out))
	}
	settled := newRefreshFixture(t)
	r := mustQueueRegistry(t, settled.workItemJourneyFixture)
	r.EpicBaseUpdates = append(r.EpicBaseUpdates, EpicBaseUpdate{ID: "pending", Phase: "intent", Plan: WorkspaceEpicBaseUpdatePlan{Target: QueueTarget{EpicID: "epic", RepoID: "ply"}}})
	projects, repos, err := settled.dependencies.Projects.Snapshot(settled.root)
	if err != nil {
		t.Fatal(err)
	}
	plan, _, _, err := buildIntegrationPlan(settled.dependencies, settled.root, ProjectSnapshot{Projects: projects, Repos: repos}, r, settled.input(true))
	if err != nil || plan.Readiness == "ready" || !slices.Contains(plan.Reasons, "epic_base_update_pending") {
		t.Fatalf("pending base update did not block refresh: %v %v", plan.Reasons, err)
	}
}

func TestRefreshParentRecoversExactEffectAfterCrash(t *testing.T) {
	f := newRefreshFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	refreshApply(t, f.dependencies, f.input(true))
	err := f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		last := len(r.IntegrationResults) - 1
		for i, a := range r.IntegrationAttempts {
			if a.ResultID != nil && *a.ResultID == r.IntegrationResults[last].ID {
				r.IntegrationAttempts[i].State, r.IntegrationAttempts[i].ResultID = "prepared", nil
			}
		}
		r.IntegrationResults = r.IntegrationResults[:last]
		return session.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	pending := refreshCheck(t, f.dependencies, f.input(true))
	in := f.input(true)
	in.Confirmation = integrationDigest(t, pending)
	out, err := ApplyTaskIntegration(f.dependencies, in)
	if err != nil || out.Readback.Classification != "exact_effect" || !strings.Contains(integrationPlanFromReadback(t, out.Readback), `"recovered":true`) || git.count() != 1 {
		t.Fatalf("resume = %s merges=%d err=%v", out.Readback.Classification, git.count(), err)
	}
}

func TestRefreshParentNoEffectRetryKeepsItsBasis(t *testing.T) {
	f := newRefreshFixture(t)
	real := f.dependencies.IntegrationGit
	f.dependencies.IntegrationGit = &noEffectIntegrationGit{TaskIntegrationGit: real}
	out := refreshApply(t, f.dependencies, f.input(true))
	if out.Classification != "no_effect" {
		t.Fatalf("classification = %s", out.Classification)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	prior := r.IntegrationResults[len(r.IntegrationResults)-1]
	want := []string{"ply", "workspace", "task", "integrate", "second", "--refresh-parent", "--expected-result-oid", f.c, "--expected-parent-oid", f.b, "--retry-after", string(prior.ID), "--check"}
	if !slices.Equal(out.NextAction.Argv, want) {
		t.Fatalf("no-effect argv = %v", out.NextAction.Argv)
	}
	f.dependencies.IntegrationGit = real
	retry := f.input(true)
	retry.RetryAfterResultID = &prior.ID
	plain := retry
	plain.RefreshParent = false
	if out := refreshCheck(t, f.dependencies, plain); out.Classification == "ready" {
		t.Fatal("refresh retry without the flag became ready")
	}
	checked := refreshCheck(t, f.dependencies, retry)
	if checked.Classification != "ready" || !slices.Contains(checked.NextAction.Argv, "--refresh-parent") || !slices.Contains(checked.NextAction.Argv, string(prior.ID)) || !strings.Contains(integrationPlanFromReadback(t, checked), `"expected_parent_oid":"`+f.b+`"`) {
		t.Fatalf("retry check = %s %v", checked.Classification, checked.NextAction.Argv)
	}

	runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "later Epic work")
	refreshBaseUpdate(t, f.integrationJourneyFixture)
	moved := refreshCheck(t, f.dependencies, retry)
	if moved.Classification == "ready" || !strings.Contains(planReasons(t, moved), "refresh_parent_not_registered_base") || !strings.Contains(moved.NextAction.Reason, "never selects another basis") || !strings.Contains(moved.NextAction.Reason, f.b) {
		t.Fatalf("retry after Epic moved = %s\n%s", moved.Classification, planReasons(t, moved))
	}
	d := runLocalGit(t, f.epicPath, "rev-parse", "HEAD")
	other := retry
	other.ExpectedParentOID = d
	if out := refreshCheck(t, f.dependencies, other); out.Classification == "ready" || !strings.Contains(planReasons(t, out), "retry_binding_mismatch") {
		t.Fatalf("retry selected another basis = %s\n%s", out.Classification, planReasons(t, out))
	}
}

func TestPlainNoEffectRetryKeepsCurrentArgv(t *testing.T) {
	f := newIntegrationJourneyFixture(t)
	f.dependencies.IntegrationGit = &noEffectIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	out := refreshApply(t, f.dependencies, TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid})
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	want := []string{"ply", "workspace", "task", "integrate", "task", "--expected-result-oid", f.resultOID, "--expected-parent-oid", f.oid, "--retry-after", string(r.IntegrationResults[0].ID), "--check"}
	if out.Classification != "no_effect" || !slices.Equal(out.NextAction.Argv, want) || readbackMember(out, "kind") != "WorkspaceTaskIntegrationReadback@6" {
		t.Fatalf("plain no-effect = %s %v", out.Classification, out.NextAction.Argv)
	}
}

func TestRefreshParentConcurrentApplyMergesOnce(t *testing.T) {
	f := newRefreshFixture(t)
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	in := f.input(true)
	in.Confirmation = integrationDigest(t, refreshCheck(t, f.dependencies, in))
	var wg sync.WaitGroup
	results := make([]string, 2)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			out, err := ApplyTaskIntegration(f.dependencies, in)
			if err == nil {
				results[i] = out.Readback.Classification
			}
		}(i)
	}
	wg.Wait()
	if git.count() != 1 || results[0] != "exact_effect" || results[1] != "exact_effect" {
		t.Fatalf("concurrent apply = %v merges=%d", results, git.count())
	}
}

func TestRefreshParentDriftBetweenCheckAndApplyChangesDigest(t *testing.T) {
	f := newRefreshFixture(t)
	in := f.input(true)
	in.Confirmation = integrationDigest(t, refreshCheck(t, f.dependencies, in))
	runLocalGit(t, f.epicPath, "commit", "--allow-empty", "-m", "drift")
	if _, err := ApplyTaskIntegration(f.dependencies, in); err == nil || !strings.Contains(err.Error(), "fresh plan digest") {
		t.Fatalf("drift apply err = %v", err)
	}
	if runLocalGit(t, f.secondPath, "rev-parse", "HEAD") != f.c {
		t.Fatal("source moved")
	}
}

func TestStoredRefreshPlanNeedsPreservedEpicBasis(t *testing.T) {
	f := newRefreshFixture(t)
	refreshApply(t, f.dependencies, f.input(true))
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	if err := validateTaskLifecycleRegistry(r); err != nil {
		t.Fatal(err)
	}
	last := len(r.IntegrationAuthorities) - 1
	for name, oid := range map[string]string{"original basis": f.a, "unrecorded parent": f.c} {
		t.Run(name, func(t *testing.T) {
			tampered := cloneWorkItemRegistry(r)
			plan := &tampered.IntegrationAuthorities[last].Plan
			plan.Epic.ExpectedParentOID, plan.Effect.ExpectedParentOID = oid, oid
			if err := validateStoredIntegrationPlan(*plan, tampered); err == nil {
				t.Fatal("refresh plan without a preserved new Epic basis was accepted")
			}
		})
	}
	plain := cloneWorkItemRegistry(r)
	plain.IntegrationAuthorities[last].Plan.Kind = currentIntegrationPlanKind
	if err := validateStoredIntegrationPlan(plain.IntegrationAuthorities[last].Plan, plain); err == nil {
		t.Fatal("a plan@6 kind accepted a refreshed basis")
	}
}

func TestPendingPlainAuthorityKeepsItsOwnContract(t *testing.T) {
	f := newRefreshFixture(t)
	runLocalGit(t, f.secondPath, "reset", "--hard", f.a)
	if err := os.WriteFile(filepath.Join(f.secondPath, "plain.txt"), []byte("plain\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, f.secondPath, "add", "plain.txt")
	runLocalGit(t, f.secondPath, "commit", "-m", "unadapted")
	plain := TaskIntegrationInput{TaskID: "second", ExpectedResultOID: runLocalGit(t, f.secondPath, "rev-parse", "HEAD"), ExpectedParentOID: f.a}
	err := f.dependencies.ProjectLocks.WithSnapshotLock(f.root, func(projects ProjectSnapshot) error {
		return f.dependencies.WorkItems.WithLock(f.root, func(session WorkItemStoreSession) error {
			r, err := session.Snapshot()
			if err != nil {
				return err
			}
			plan, digest, _, err := buildIntegrationPlan(f.dependencies, f.root, projects, r, plain)
			if err != nil {
				return err
			}
			plan.Readiness, plan.Reasons = "ready", []string{}
			plan.ObservedParent.OID, plan.ObservedParent.Tree = f.a, runLocalGit(t, f.secondPath, "rev-parse", f.a+"^{tree}")
			for i := range plan.ObservedInventory {
				if plan.ObservedInventory[i].Locator == plan.Epic.ParentLocator {
					plan.ObservedInventory[i].OID = f.a
				}
			}
			digest = integrationPlanDigest(plan)
			a := IntegrationAuthority{ID: IntegrationAuthorityID("iauth_" + strings.Repeat("9", 32)), CreatedAtUTC: "2026-09-28T12:02:00Z", PlanSHA256: digest, Plan: plan, TaskID: "second", AllowedEffect: IntegrationAllowedEffect{Kind: "local_ff_only", ParentRef: plan.Epic.ParentRef, ExpectedParentOID: f.a, ResultOID: plan.Task.ResultOID, MaxOccurrences: 1}}
			i := intentFromPlan(IntegrationIntentID("iint_"+strings.Repeat("9", 32)), a.ID, plan, digest)
			i.IntentSHA256 = integrationIntentDigest(i)
			r.IntegrationAuthorities = append(r.IntegrationAuthorities, a)
			r.IntegrationIntents = append(r.IntegrationIntents, i)
			sortWorkRegistry(&r)
			return session.Publish(r)
		})
	})
	if err != nil {
		t.Fatalf("pending plan@6 is no longer a valid stored contract: %v", err)
	}
	refresh := f.input(true)
	refresh.ExpectedResultOID = plain.ExpectedResultOID
	if out := refreshCheck(t, f.dependencies, refresh); out.Classification == "ready" || !strings.Contains(planReasons(t, out), "parent_integration_already_active") {
		t.Fatalf("refresh replaced a pending plan@6: %s", out.Classification)
	}
	r := mustQueueRegistry(t, f.workItemJourneyFixture)
	pending := r.IntegrationAuthorities[len(r.IntegrationAuthorities)-1]
	plain.Confirmation = pending.PlanSHA256
	out, err := ApplyTaskIntegration(f.dependencies, plain)
	if err != nil || out.Readback.Classification != "conflict" || readbackMember(out.Readback, "kind") != "WorkspaceTaskIntegrationReadback@6" || runLocalGit(t, f.epicPath, "rev-parse", "HEAD") != f.b {
		t.Fatalf("pending plan@6 continued outside its exact contract: %s %v", out.Readback.Classification, err)
	}
}

func TestTaskShowFallbackKeepsTheAuthorityReadbackVersion(t *testing.T) {
	f := newRefreshFixture(t)
	refreshApply(t, f.dependencies, f.input(true))
	f.dependencies.IntegrationGit = unsupportedIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	for task, want := range map[TaskID]string{"second": "WorkspaceTaskIntegrationReadback@7", "task": "WorkspaceTaskIntegrationReadback@6"} {
		shown, err := ShowTask(f.dependencies, task)
		if err != nil || shown.Integration == nil || shown.Integration.Classification != "exact_effect" {
			t.Fatalf("%s fallback = %#v err=%v", task, shown.Integration, err)
		}
		if kind := readbackMember(*shown.Integration, "kind"); kind != want || readbackMember(*shown.Integration, "schema_version") != int64(want[len(want)-1]-'0') {
			t.Fatalf("%s fallback wrapper = %v", task, kind)
		}
		raw, err := MarshalTaskIntegrationReadback(*shown.Integration)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(raw), `"plan":null`) || !strings.Contains(string(raw), `"observed_source":{"freshness":{"clean":"unknown"`) {
			t.Fatalf("%s fallback invented fresh observations: %s", task, raw)
		}
	}
	shown, _ := ShowTask(f.dependencies, "second")
	epic, _ := readbackMember(*shown.Integration, "epic").(canonicaljson.Object)
	task, _ := readbackMember(*shown.Integration, "task").(canonicaljson.Object)
	source, _ := readbackMember(*shown.Integration, "source").(canonicaljson.Object)
	if member(epic, "expected_parent_oid") != f.b || member(epic, "expected_parent_tree") != runLocalGit(t, f.epicPath, "rev-parse", f.b+"^{tree}") || member(task, "result_oid") != f.c || member(task, "result_tree") != f.cTree || member(source, "oid") != f.c {
		t.Fatalf("fallback presents another integration basis: epic=%v task=%v source=%v", epic, task, source)
	}
}

func member(object canonicaljson.Object, name string) any {
	for _, m := range object {
		if m.Name == name {
			return m.Value
		}
	}
	return nil
}
