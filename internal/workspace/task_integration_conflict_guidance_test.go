package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestIntegrationConflictGuidanceUsesFreshQualifiedFacts(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	projects, repos, err := fixture.dependencies.Projects.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	plan, digest, ctx, err := buildIntegrationPlan(fixture.dependencies, fixture.root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, input)
	if err != nil {
		t.Fatal(err)
	}
	plan.Readiness = "conflict"
	plan.Reasons = []string{"expected_oid_mismatch", "source_binding_mismatch"}
	ctx.ResultOID = strings.Repeat("f", len(fixture.resultOID))
	got := newIntegrationReadback(fixture.root, registry, ctx, plan, digest, nil, nil, nil, nil)
	want := []string{"ply", "workspace", "task", "integrate", "task", "--expected-result-oid", fixture.resultOID, "--expected-parent-oid", fixture.oid, "--check", "--format", "json"}
	if !reflect.DeepEqual(got.NextAction.Argv, want) {
		t.Fatalf("next action = %#v, want %v", got.NextAction, want)
	}
	if !strings.Contains(RenderTaskIntegrationText(got), strings.Join(want, " ")) {
		t.Fatalf("text = %q", RenderTaskIntegrationText(got))
	}
	top := integrationReadbackObject(t, integrationReadbackMember(t, got.Value, "next_action"), "next_action")
	if integrationReadbackMember(t, top, "reason") != got.NextAction.Reason {
		t.Fatalf("JSON next action = %#v", top)
	}

	cases := []struct {
		name       string
		reasons    []string
		change     func(*integrationContext)
		wantArgv   []string
		wantReason string
	}{
		{"stale parent argument", []string{"expected_oid_mismatch"}, nil, want, "stale"},
		{"pending before stale", []string{"epic_base_update_pending", "expected_oid_mismatch"}, nil, nil, "pending base update"},
		{"unknown before stale", []string{"prior_integration_requires_reconciliation", "expected_oid_mismatch"}, nil, []string{"ply", "workspace", "task", "show", "task", "--format", "json"}, "unresolved"},
		{"dirty before stale", []string{"source_not_clean", "expected_oid_mismatch"}, nil, nil, "dirty"},
		{"non descendant before stale", []string{"result_not_descendant", "expected_oid_mismatch"}, nil, nil, "not descended"},
		{"identity before stale", []string{"source_binding_mismatch"}, func(c *integrationContext) { c.Source.Ref = "refs/heads/other" }, nil, "identity mismatch"},
		{"foreign parent in inventory", []string{"expected_oid_mismatch"}, func(c *integrationContext) {
			c.ParentObservation.Ref = "refs/heads/foreign"
			for i := range c.Repository.Inventory {
				if c.Repository.Inventory[i].Locator == c.ParentObservation.Locator {
					c.Repository.Inventory[i].Ref = c.ParentObservation.Ref
				}
			}
		}, nil, "parent binding"},
		{"foreign parent locator in inventory", []string{"expected_oid_mismatch"}, func(c *integrationContext) {
			old := c.ParentObservation.Locator
			c.ParentObservation.Locator += "-foreign"
			for i := range c.Repository.Inventory {
				if c.Repository.Inventory[i].Locator == old {
					c.Repository.Inventory[i].Locator = c.ParentObservation.Locator
				}
			}
		}, nil, "parent binding"},
		{"foreign parent common directory", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.ParentObservation.GitCommonDir += "-foreign" }, nil, "parent binding"},
		{"detached source", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.Source.Symbolic = false }, nil, "symbolic"},
		{"detached parent", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.ParentObservation.Symbolic = false }, nil, "symbolic"},
		{"unsupported object format", []string{"expected_oid_mismatch"}, func(c *integrationContext) {
			c.Source.ObjectFormat, c.ParentObservation.ObjectFormat, c.Repository.ObjectFormat = "other", "other", "other"
		}, nil, "repository format"},
		{"unsupported ref format", []string{"expected_oid_mismatch"}, func(c *integrationContext) {
			c.Source.RefFormat, c.ParentObservation.RefFormat, c.Repository.RefFormat = "reftable", "reftable", "reftable"
		}, nil, "repository format"},
		{"shallow repository", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.Repository.Shallow = true }, nil, "repository state"},
		{"partial clone", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.Repository.PartialClone = true }, nil, "repository state"},
		{"sparse checkout", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.Repository.SparseCheckout = true }, nil, "repository state"},
		{"unknown observations", []string{"expected_oid_mismatch"}, func(c *integrationContext) { c.SourceObserved = false }, nil, "Fresh source"},
		{"Task basis older than recorded base", []string{"epic_base_changed", "expected_oid_mismatch"}, func(c *integrationContext) { c.Task.Worktree.ParentOID = strings.Repeat("e", len(fixture.oid)) }, nil, "--refresh-parent --expected-result-oid <adapted commit>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, c := plan, ctx
			p.Reasons = slices.Clone(tc.reasons)
			c.Repository.Inventory = slices.Clone(ctx.Repository.Inventory)
			if tc.change != nil {
				tc.change(&c)
			}
			action := newIntegrationReadback(fixture.root, registry, c, p, digest, nil, nil, nil, nil).NextAction
			if tc.wantArgv == nil {
				if len(action.Argv) != 0 {
					t.Fatalf("unexpected argv = %#v", action)
				}
			} else if !slices.Equal(action.Argv, tc.wantArgv) {
				t.Fatalf("argv = %#v", action)
			}
			if !strings.Contains(action.Reason, tc.wantReason) {
				t.Fatalf("reason = %q", action.Reason)
			}
		})
	}
}

func TestIntegrationConflictGuidanceNamesOtherTaskBlockingParent(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	projects, repos, err := fixture.dependencies.Projects.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	plan, digest, ctx, err := buildIntegrationPlan(fixture.dependencies, fixture.root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, input)
	if err != nil {
		t.Fatal(err)
	}
	other := TaskID("other-task")
	registry.Tasks = append(registry.Tasks, TaskRecord{ID: other, ParentEpicID: ctx.Epic.ID, RepoID: ctx.Task.RepoID})
	otherPlan := plan
	otherPlan.Task.TaskID = other
	registry.IntegrationAuthorities = append(registry.IntegrationAuthorities, IntegrationAuthority{ID: "int_auth_other", TaskID: other, Plan: otherPlan})
	plan.Readiness, plan.Reasons = "conflict", []string{"parent_integration_already_active"}
	got := newIntegrationReadback(fixture.root, registry, ctx, plan, digest, nil, nil, nil, nil).NextAction
	want := []string{"ply", "workspace", "task", "show", string(other), "--format", "json"}
	if !slices.Equal(got.Argv, want) || !strings.Contains(got.Reason, "int_auth_other") || !strings.Contains(got.Reason, string(other)) {
		t.Fatalf("other Task guidance = %#v", got)
	}
}

func TestIntegrationConflictGuidanceChangesAfterBasePreview(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	epicPath := filepath.Join(fixture.wrapper, "epic")
	if err := os.WriteFile(filepath.Join(epicPath, "later.txt"), []byte("later\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, epicPath, "add", "later.txt")
	runLocalGit(t, epicPath, "commit", "-m", "later")
	parentOID := runLocalGit(t, epicPath, "rev-parse", "HEAD")
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if checked.Readback.Classification != "conflict" || !slices.Equal(checked.Readback.NextAction.Argv, wantBaseGuidanceArgv()) || !strings.Contains(checked.Readback.NextAction.Reason, parentOID) {
		t.Fatalf("base preview guidance = %#v", checked.Readback.NextAction)
	}
	applyGuidedBaseUpdate(t, fixture, checked.Readback.NextAction)
	after, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Readback.NextAction.Argv) != 0 || !strings.Contains(after.Readback.NextAction.Reason, "--refresh-parent") || !strings.Contains(after.Readback.NextAction.Reason, fixture.oid) || !strings.Contains(after.Readback.NextAction.Reason, parentOID) {
		t.Fatalf("after base update = %#v", after.Readback.NextAction)
	}
}

func TestIntegrationConflictGuidanceRechecksObservedOIDsWithoutOldDigest(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.oid, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"ply", "workspace", "task", "integrate", "task", "--expected-result-oid", fixture.resultOID, "--expected-parent-oid", fixture.oid, "--check", "--format", "json"}
	if !slices.Equal(checked.Readback.NextAction.Argv, want) {
		t.Fatalf("next = %#v", checked.Readback.NextAction)
	}
	if !strings.Contains(checked.Readback.NextAction.Reason, "verify ancestry and other preconditions") {
		t.Fatalf("ancestry guidance = %q", checked.Readback.NextAction.Reason)
	}
	if strings.Contains(strings.Join(checked.Readback.NextAction.Argv, " "), "--apply") || strings.Contains(strings.Join(checked.Readback.NextAction.Argv, " "), "--confirm") {
		t.Fatalf("mutating guidance = %#v", checked.Readback.NextAction)
	}
}

func TestIntegrationConflictGuidanceAfterConfirmationPreservesStoredResult(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	registry, err := fixture.dependencies.WorkItems.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	projects, repos, err := fixture.dependencies.Projects.Snapshot(fixture.root)
	if err != nil {
		t.Fatal(err)
	}
	plan, digest, ctx, err := buildIntegrationPlan(fixture.dependencies, fixture.root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, input)
	if err != nil {
		t.Fatal(err)
	}
	ctx.Source.OID = fixture.oid
	for i := range ctx.Repository.Inventory {
		if ctx.Repository.Inventory[i].Locator == ctx.Source.Locator {
			ctx.Repository.Inventory[i].OID = fixture.oid
		}
	}
	stored := IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: "Source changed after confirmation.", Argv: []string{}}
	unchanged := false
	result := IntegrationResult{Outcome: "conflict", GitChanged: &unchanged, RecoveryStatus: "conflict", NextAction: stored}
	readback := newIntegrationReadback(fixture.root, registry, ctx, plan, digest, nil, nil, nil, &result)
	if readback.Classification != "conflict" || readback.GitChanged == nil || *readback.GitChanged || readback.NextAction.Reason == stored.Reason || len(readback.NextAction.Argv) == 0 {
		t.Fatalf("fresh guidance = %#v", readback)
	}
	if !reflect.DeepEqual(result.NextAction, stored) {
		t.Fatalf("stored result changed = %#v", result.NextAction)
	}
	storedValue := integrationReadbackObject(t, integrationReadbackMember(t, readback.Value, "integration_result"), "integration_result")
	storedAction := integrationReadbackObject(t, integrationReadbackMember(t, storedValue, "next_action"), "integration_result.next_action")
	if integrationReadbackMember(t, storedAction, "reason") != stored.Reason {
		t.Fatalf("stored JSON = %#v", storedAction)
	}
}
