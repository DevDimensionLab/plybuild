package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
)

func wantBaseGuidanceArgv() []string {
	return []string{"ply", "workspace", "epic", "base", "update", "epic", "--repo", "ply", "--check", "--format", "json"}
}

func assertBaseGuidance(t *testing.T, readback WorkspaceTaskIntegrationReadback) {
	t.Helper()
	if readback.NextAction.Kind != "new_check" || !slices.Equal(readback.NextAction.Argv, wantBaseGuidanceArgv()) {
		t.Fatalf("next action = %#v", readback.NextAction)
	}
	top := integrationReadbackObject(t, integrationReadbackMember(t, readback.Value, "next_action"), "next_action")
	if got := integrationReadbackMember(t, top, "kind"); got != "new_check" {
		t.Fatalf("top-level next_action.kind = %#v", got)
	}
	if !strings.Contains(RenderTaskIntegrationText(readback), "Run `"+strings.Join(wantBaseGuidanceArgv(), " ")+"`.") {
		t.Fatalf("text = %q", RenderTaskIntegrationText(readback))
	}
}

func assertNoBaseGuidance(t *testing.T, readback WorkspaceTaskIntegrationReadback) {
	t.Helper()
	if readback.NextAction.Kind == "new_check" {
		t.Fatalf("unexpected base guidance: %#v", readback.NextAction)
	}
	top := integrationReadbackObject(t, integrationReadbackMember(t, readback.Value, "next_action"), "next_action")
	if got := integrationReadbackMember(t, top, "kind"); got != readback.NextAction.Kind {
		t.Fatalf("top-level next_action.kind = %#v, want %s", got, readback.NextAction.Kind)
	}
}

func assertStoredSuccessNextAction(t *testing.T, readback WorkspaceTaskIntegrationReadback) {
	t.Helper()
	stored := integrationReadbackObject(t, integrationReadbackMember(t, readback.Value, "integration_result"), "integration_result")
	next := integrationReadbackObject(t, integrationReadbackMember(t, stored, "next_action"), "integration_result.next_action")
	if got := integrationReadbackMember(t, next, "kind"); got != "none" {
		t.Fatalf("stored next_action.kind = %#v", got)
	}
	if argv, ok := integrationReadbackMember(t, next, "argv").([]any); ok && len(argv) != 0 {
		t.Fatalf("stored next_action.argv = %#v", argv)
	}
}

func applyIntegrationFixture(t *testing.T, fixture integrationJourneyFixture) (TaskIntegrationInput, WorkspaceTaskIntegrationReadback) {
	t.Helper()
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	apply := input
	apply.Apply, apply.Confirmation = true, integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(fixture.dependencies, apply)
	if err != nil {
		t.Fatal(err)
	}
	return input, applied.Readback
}

func applyGuidedBaseUpdate(t *testing.T, fixture integrationJourneyFixture, guidance IntegrationNextAction) EpicBaseResult {
	t.Helper()
	if !slices.Equal(guidance.Argv, wantBaseGuidanceArgv()) {
		t.Fatalf("guidance = %#v", guidance)
	}
	in := EpicBaseInput{EpicID: EpicID(guidance.Argv[5]), RepoID: RepoID(guidance.Argv[7])}
	preview, err := UpdateEpicBase(fixture.dependencies, in)
	if err != nil || preview.Confirmation == nil {
		t.Fatalf("preview = %#v, error = %v", preview, err)
	}
	if preview.Plan.RegistryUpgrade != nil {
		in.Upgrade = preview.Plan.RegistryUpgrade.RegistrySHA256
	}
	in.Apply, in.Confirmation = true, *preview.Confirmation
	applied, err := UpdateEpicBase(fixture.dependencies, in)
	if err != nil {
		t.Fatalf("apply = %#v, error = %v", applied, err)
	}
	return applied
}

func readRegistryBytes(t *testing.T, fixture integrationJourneyFixture) []byte {
	t.Helper()
	b, err := os.ReadFile(workItemsPath(fixture.root))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestIntegrationBaseGuidanceAfterExactEffectKeepsStoredResult(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input, applied := applyIntegrationFixture(t, fixture)
	if applied.Classification != "exact_effect" || applied.RecoveryStatus != "complete" {
		t.Fatalf("applied = %#v", applied)
	}
	assertBaseGuidance(t, applied)
	assertStoredSuccessNextAction(t, applied)
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil || len(registry.IntegrationResults) != 1 {
		t.Fatalf("registry results = %#v, error = %v", registry.IntegrationResults, err)
	}
	if stored := registry.IntegrationResults[0]; stored.NextAction.Kind != "none" || len(stored.NextAction.Argv) != 0 {
		t.Fatalf("stored next action = %#v", stored.NextAction)
	}
	before := readRegistryBytes(t, fixture)
	replayed, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	assertBaseGuidance(t, replayed.Readback)
	assertStoredSuccessNextAction(t, replayed.Readback)
	if got := integrationReadbackMember(t, integrationReadbackObject(t, integrationReadbackMember(t, replayed.Readback.Value, "observed_inventory"), "observed_inventory"), "freshness"); got != "unknown" {
		t.Fatalf("persisted replay presents historical inventory as %#v", got)
	}
	shown, err := ShowTask(fixture.dependencies, "task")
	if err != nil || shown.Integration == nil {
		t.Fatalf("shown = %#v, error = %v", shown.Integration, err)
	}
	assertBaseGuidance(t, *shown.Integration)
	if after := readRegistryBytes(t, fixture); string(after) != string(before) {
		t.Fatal("readback changed the registry")
	}
}

func TestIntegrationBaseGuidanceAlreadyIntegratedFollowsBaseFreshness(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	epicPath := filepath.Join(fixture.wrapper, "epic")
	runLocalGit(t, epicPath, "merge", "--ff-only", fixture.resultOID)
	reflog := runLocalGit(t, epicPath, "reflog", "show", "--format=%H", "refs/heads/epic")
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || checked.Readback.Classification != "already_integrated" {
		t.Fatalf("checked = %#v, error = %v", checked.Readback, err)
	}
	assertBaseGuidance(t, checked.Readback)
	apply := input
	apply.Apply, apply.Confirmation = true, integrationDigest(t, checked.Readback)
	applied, err := ApplyTaskIntegration(fixture.dependencies, apply)
	if err != nil || applied.Readback.Classification != "already_integrated" || applied.Readback.RecoveryStatus != "complete" {
		t.Fatalf("applied = %#v, error = %v", applied.Readback, err)
	}
	assertBaseGuidance(t, applied.Readback)
	assertStoredSuccessNextAction(t, applied.Readback)
	applyGuidedBaseUpdate(t, fixture, applied.Readback.NextAction)
	fresh, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || fresh.Readback.Classification != "already_integrated" || fresh.Readback.NextAction.Kind != "none" {
		t.Fatalf("fresh-base readback = %#v, error = %v", fresh.Readback, err)
	}
	if got := runLocalGit(t, epicPath, "reflog", "show", "--format=%H", "refs/heads/epic"); got != reflog {
		t.Fatalf("readback or base update moved Git: %q, want %q", got, reflog)
	}
}

func TestIntegrationBaseGuidanceUsesCurrentParentNotHistoricalAfterObservation(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input, applied := applyIntegrationFixture(t, fixture)
	applyGuidedBaseUpdate(t, fixture, applied.NextAction)
	current, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if current.Readback.NextAction.Kind != "none" {
		t.Fatalf("fresh base readback = %#v", current.Readback.NextAction)
	}
	epicPath := filepath.Join(fixture.wrapper, "epic")
	if err := os.WriteFile(filepath.Join(epicPath, "later.txt"), []byte("later Epic work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, epicPath, "add", "later.txt")
	runLocalGit(t, epicPath, "commit", "-m", "later Epic work")
	historical, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || !historical.Readback.HistoricalCompletion {
		t.Fatalf("historical = %#v, error = %v", historical.Readback, err)
	}
	assertBaseGuidance(t, historical.Readback)
	assertStoredSuccessNextAction(t, historical.Readback)
}

func TestIntegrationBaseGuidanceConcurrentReadsWriteNothingAndFreshBaseClearsIt(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input, applied := applyIntegrationFixture(t, fixture)
	original, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	before := readRegistryBytes(t, fixture)
	var wg sync.WaitGroup
	readbacks := make([]WorkspaceTaskIntegrationReadback, 2)
	errs := make([]error, 2)
	for i := range readbacks {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var result TaskIntegrationResult
			result, errs[i] = CheckTaskIntegration(fixture.dependencies, input)
			readbacks[i] = result.Readback
		}(i)
	}
	wg.Wait()
	for i := range readbacks {
		if errs[i] != nil {
			t.Fatal(errs[i])
		}
		assertBaseGuidance(t, readbacks[i])
	}
	if after := readRegistryBytes(t, fixture); string(after) != string(before) {
		t.Fatal("concurrent readbacks changed the registry")
	}
	applyGuidedBaseUpdate(t, fixture, applied.NextAction)
	refreshed, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil || refreshed.Readback.NextAction.Kind != "none" {
		t.Fatalf("refreshed = %#v, error = %v", refreshed.Readback.NextAction, err)
	}
	updated, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updated.IntegrationAuthorities, original.IntegrationAuthorities) || !reflect.DeepEqual(updated.IntegrationIntents, original.IntegrationIntents) || !reflect.DeepEqual(updated.IntegrationAttempts, original.IntegrationAttempts) || !reflect.DeepEqual(updated.IntegrationResults, original.IntegrationResults) {
		t.Fatal("base update changed the original integration records")
	}
}

func TestIntegrationBaseGuidanceNewTaskUsesUpdatedBase(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	_, applied := applyIntegrationFixture(t, fixture)
	updated := applyGuidedBaseUpdate(t, fixture, applied.NextAction)
	latest := updated.Versions[len(updated.Versions)-1]
	if latest.OID != fixture.resultOID || latest.Tree != fixture.resultTree || latest.PreviousRevision == nil || latest.Revision != *latest.PreviousRevision+1 {
		t.Fatalf("latest base = %#v", latest)
	}
	if _, err := CreateTask(fixture.dependencies, TaskCreateInput{TaskID: "next-task", Title: "Next", Description: "Next task", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	in, err := ParseTaskWorktreeCreateInput("next-task", "next-task", filepath.Join(fixture.wrapper, "next-task"), latest.OID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateTaskWorktree(fixture.dependencies, in); err != nil {
		t.Fatal(err)
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	task, _ := findTask(registry, "next-task")
	epic, _ := findEpic(registry, "epic")
	binding, _ := findEpicRepo(*epic, "ply")
	current := currentEpicBase(registry, *epic, *binding)
	if task == nil || task.Worktree == nil || current.Revision != latest.Revision || current.UpdateID == nil || task.Worktree.ParentOID != current.OID || task.Worktree.ParentTree != current.Tree {
		t.Fatalf("new Task basis = %#v, current base = %#v", task, current)
	}
}

func baseGuidanceContext(t *testing.T) (WorkItemRegistry, integrationContext, IntegrationNextAction) {
	t.Helper()
	fixture := newIntegrationJourneyFixture(t)
	_, applied := applyIntegrationFixture(t, fixture)
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := contextFromPlan(registry, registry.IntegrationAuthorities[0].Plan)
	if err != nil {
		t.Fatal(err)
	}
	ctx.ParentObservation, err = fixture.dependencies.IntegrationGit.ObserveIntegrationWorktree(ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref)
	if err != nil {
		t.Fatal(err)
	}
	ctx.ParentObserved = true
	if applied.NextAction.Kind != "new_check" {
		t.Fatalf("fixture guidance = %#v", applied.NextAction)
	}
	return registry, ctx, registry.IntegrationResults[0].NextAction
}

func TestIntegrationBaseGuidanceRequiresFreshQualifiedFactsAndNoPendingWork(t *testing.T) {
	registry, ctx, stored := baseGuidanceContext(t)
	if ctx.RepositoryObserved {
		t.Fatal("persisted replay context must keep historical repository facts")
	}
	if got := integrationBaseGuidance(registry, ctx, "exact_effect", "complete", stored); got.Kind != "new_check" || !slices.Equal(got.Argv, wantBaseGuidanceArgv()) {
		t.Fatalf("persisted replay with fresh parent = %#v", got)
	}
	target := QueueTarget{ctx.Epic.ProjectID, ctx.Parent.RepoID, ctx.Epic.ID, ctx.Parent.GitCommonDir, ctx.Parent.Worktree.ID, ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref}
	other := registry.IntegrationAuthorities[0]
	other.ID = "int_auth_other"
	cases := map[string]struct {
		classification, recovery string
		registry                 func(*WorkItemRegistry)
		ctx                      func(*integrationContext)
	}{
		"partial":                {classification: "partial", recovery: "reconciliation-required"},
		"unknown":                {classification: "unknown", recovery: "unknown"},
		"no effect":              {classification: "no_effect", recovery: "safe-no-effect"},
		"incomplete recovery":    {recovery: "unknown"},
		"unobserved parent":      {ctx: func(c *integrationContext) { c.ParentObserved = false }},
		"dirty parent":           {ctx: func(c *integrationContext) { c.ParentObservation.Clean = false }},
		"detached parent":        {ctx: func(c *integrationContext) { c.ParentObservation.Symbolic = false }},
		"operation in progress":  {ctx: func(c *integrationContext) { c.ParentObservation.InProgress = []string{"merge"} }},
		"other parent locator":   {ctx: func(c *integrationContext) { c.ParentObservation.Locator += "-other" }},
		"other parent ref":       {ctx: func(c *integrationContext) { c.ParentObservation.Ref = "refs/heads/other" }},
		"other common dir":       {ctx: func(c *integrationContext) { c.ParentObservation.GitCommonDir += "-other" }},
		"other repository":       {ctx: func(c *integrationContext) { c.Repo.GitCommonDir += "-other" }},
		"unsupported object":     {ctx: func(c *integrationContext) { c.ParentObservation.ObjectFormat = "" }},
		"unsupported ref format": {ctx: func(c *integrationContext) { c.ParentObservation.RefFormat = "reftable" }},
		"other supported object format": {ctx: func(c *integrationContext) {
			if c.Repository.ObjectFormat != "sha1" || c.ParentObservation.ObjectFormat != "sha1" {
				t.Fatalf("fixture formats = %q/%q", c.Repository.ObjectFormat, c.ParentObservation.ObjectFormat)
			}
			c.ParentObservation.ObjectFormat = "sha256"
		}},
		"unqualified registry": {registry: func(r *WorkItemRegistry) { r.RawSHA256 = "" }},
		"pending base update": {registry: func(r *WorkItemRegistry) {
			r.EpicBaseUpdates = append(r.EpicBaseUpdates, EpicBaseUpdate{ID: "ebu_pending", Plan: WorkspaceEpicBaseUpdatePlan{Target: target}, Phase: "intent"})
		}},
		"pending worktree": {registry: func(r *WorkItemRegistry) {
			r.WorktreeOperations = append(r.WorktreeOperations, WorktreeOperationRecord{ID: "wop_pending", ParentEpicID: ctx.Epic.ID, RepoID: ctx.Parent.RepoID, State: "creating"})
		}},
		"shared ref authority": {registry: func(r *WorkItemRegistry) {
			shared := other
			shared.TaskID, shared.Plan.Epic.EpicID = "other-task", "other-epic"
			r.IntegrationAuthorities = append(r.IntegrationAuthorities, shared)
		}},
		"several Task leaves": {registry: func(r *WorkItemRegistry) {
			r.IntegrationAuthorities = append(r.IntegrationAuthorities, other)
			result := r.IntegrationResults[0]
			result.ID, result.AuthorityID = "int_res_other", other.ID
			r.IntegrationResults = append(r.IntegrationResults, result)
		}},
		"fresh format 3 base": {ctx: func(c *integrationContext) {
			c.Parent.Worktree.OID, c.Parent.Worktree.Tree = c.ParentObservation.OID, c.ParentObservation.Tree
		}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := registry
			r.IntegrationAuthorities = slices.Clone(registry.IntegrationAuthorities)
			r.IntegrationResults = slices.Clone(registry.IntegrationResults)
			r.EpicBaseVersions = slices.Clone(registry.EpicBaseVersions)
			r.EpicBaseUpdates = slices.Clone(registry.EpicBaseUpdates)
			r.WorktreeOperations = slices.Clone(registry.WorktreeOperations)
			x := ctx
			x.ParentObservation.InProgress = slices.Clone(ctx.ParentObservation.InProgress)
			if c.registry != nil {
				c.registry(&r)
			}
			if c.ctx != nil {
				c.ctx(&x)
			}
			classification, recovery := c.classification, c.recovery
			if classification == "" {
				classification = "exact_effect"
			}
			if recovery == "" {
				recovery = "complete"
			}
			if got := integrationBaseGuidance(r, x, classification, recovery, stored); !reflect.DeepEqual(got, stored) {
				t.Fatalf("guidance = %#v, want stored %#v", got, stored)
			}
		})
	}
}

func TestIntegrationBaseGuidanceNeverReplacesRecoveryAfterResultPublicationFailure(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	publishes := 0
	fixture.dependencies.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
		if operation == "temp-open" {
			publishes++
			if publishes == 3 {
				return errInjected
			}
		}
		return nil
	}}}
	failed, err := ApplyTaskIntegration(fixture.dependencies, input)
	if err == nil || failed.Readback.NextAction.Kind != "read_only_recovery_control" {
		t.Fatalf("failed = %#v, error = %v", failed.Readback.NextAction, err)
	}
	assertNoBaseGuidance(t, failed.Readback)
}

func TestIntegrationBaseGuidanceAbsentFromPartialAndHistoricalFallback(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	git := &partialAfterCommandGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}
	fixture.dependencies.IntegrationGit = git
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	partial, _ := ApplyTaskIntegration(fixture.dependencies, input)
	if partial.Readback.Classification == "exact_effect" || partial.Readback.Classification == "already_integrated" {
		t.Fatalf("partial = %#v", partial.Readback)
	}
	assertNoBaseGuidance(t, partial.Readback)

	complete := newIntegrationJourneyFixture(t)
	applyIntegrationFixture(t, complete)
	complete.dependencies.IntegrationGit = unsupportedIntegrationGit{TaskIntegrationGit: complete.dependencies.IntegrationGit}
	shown, err := ShowTask(complete.dependencies, "task")
	if err != nil || shown.Integration == nil || shown.Integration.Classification != "exact_effect" || shown.Integration.NextAction.Kind != "none" {
		t.Fatalf("historical fallback = %#v, error = %v", shown.Integration, err)
	}
	assertNoBaseGuidance(t, *shown.Integration)
}
