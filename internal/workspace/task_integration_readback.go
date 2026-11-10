package workspace

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type WorkspaceTaskIntegrationReadback struct {
	HistoricalCompletion bool `json:"historical_completion,omitempty"`
	Value                canonicaljson.Object
	Classification       string
	GitChanged           *bool
	RecoveryStatus       string
	NextAction           IntegrationNextAction
	TaskID               TaskID
	ResultOID            string
	EpicID               EpicID
	ParentRef            string
	ParentOID            string
}

func MarshalTaskIntegrationReadback(result WorkspaceTaskIntegrationReadback) ([]byte, error) {
	return canonicaljson.Marshal(result.Value)
}

func newIntegrationReadback(root string, registry WorkItemRegistry, ctx integrationContext, plan WorkspaceTaskIntegrationPlan, digest string, authority *IntegrationAuthority, intent *IntegrationIntent, attempt *IntegrationAttempt, result *IntegrationResult) WorkspaceTaskIntegrationReadback {
	classification := plan.Readiness
	changed := boolPointer(false)
	recovery := "blocked"
	next := integrationPlanNextAction(ctx, plan, digest)
	if classification == "ready" {
		recovery = "safe-no-effect"
	}
	if classification == "already_integrated" {
		recovery = "complete"
		next = IntegrationNextAction{Kind: "none", Reason: "The Task is integrated.", Argv: []string{}}
	}
	if result != nil {
		classification = result.Outcome
		changed = result.GitChanged
		recovery = result.RecoveryStatus
		next = result.NextAction
	}
	next = integrationCurrentGuidance(registry, ctx, plan, classification, recovery, next)
	authorityValue := nullableStruct(authority)
	intentValue := nullableStruct(intent)
	attemptValue := nullableStruct(attempt)
	resultValue := nullableStruct(result)
	planValue := canonicaljson.Object{{Name: "sha256", Value: digest}, {Name: "content", Value: integrationPlanCanonical(plan)}}
	persistedAfter := integrationPersistedView(registry, authority, intent, attempt, result)
	persistedBefore := ctx.PersistedBefore
	if persistedBefore == nil {
		persistedBefore = persistedAfter
	}
	observedSource := observedWorktreeReadback(ctx.Source, ctx.SourceObserved)
	observedParent := observedWorktreeReadback(ctx.ParentObservation, ctx.ParentObserved)
	observedInventory := []IntegrationInventoryEntry{}
	if ctx.RepositoryObserved {
		observedInventory = ctx.Repository.Inventory
	}
	observedReflog := []IntegrationReflogEntry{}
	if ctx.ReflogObserved {
		observedReflog = ctx.Reflog
	}
	kind, schema := integrationReadbackVersion(plan)
	readback := canonicaljson.Object{
		{Name: "kind", Value: kind}, {Name: "schema_version", Value: int64(schema)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "workspace", Value: structCanonical(plan.Workspace)}, {Name: "project", Value: structCanonical(plan.Project)}, {Name: "repository", Value: structCanonical(plan.Repository)}, {Name: "epic", Value: structCanonical(plan.Epic)}, {Name: "task", Value: structCanonical(plan.Task)}, {Name: "source", Value: structCanonical(plan.ObservedSource)},
		{Name: "plan", Value: planValue}, {Name: "authority", Value: authorityValue}, {Name: "intent", Value: intentValue}, {Name: "attempt", Value: attemptValue}, {Name: "integration_result", Value: resultValue},
		{Name: "persisted_before", Value: persistedBefore}, {Name: "persisted_after", Value: persistedAfter}, {Name: "observed_source", Value: observedSource}, {Name: "observed_parent", Value: observedParent}, {Name: "observed_inventory", Value: canonicaljson.Object{{Name: "entries", Value: structCanonical(observedInventory)}, {Name: "freshness", Value: collectionFreshness(ctx.RepositoryObserved)}}}, {Name: "observed_reflog", Value: canonicaljson.Object{{Name: "entries", Value: structCanonical(observedReflog)}, {Name: "freshness", Value: collectionFreshness(ctx.ReflogObserved)}}},
		{Name: "classification", Value: classification}, {Name: "git_changed", Value: pointerValue(changed)}, {Name: "recovery_status", Value: recovery}, {Name: "next_action", Value: structCanonical(next)},
	}
	parentOID := plan.Epic.ExpectedParentOID
	if ctx.ParentObserved && ctx.ParentObservation.OID != "" {
		parentOID = ctx.ParentObservation.OID
	} else if result != nil && result.BeforeObservation.OID != "" {
		parentOID = result.BeforeObservation.OID
	} else if plan.ObservedParent.OID != "" {
		parentOID = plan.ObservedParent.OID
	}
	return WorkspaceTaskIntegrationReadback{Value: readback, Classification: classification, GitChanged: changed, RecoveryStatus: recovery, NextAction: next, TaskID: ctx.Task.ID, ResultOID: ctx.ResultOID, EpicID: ctx.Epic.ID, ParentRef: ctx.Parent.Worktree.Ref, ParentOID: parentOID}
}

const integrationReadbackKind = "WorkspaceTaskIntegrationReadback@6"
const integrationReadbackSchema = 6
const refreshIntegrationReadbackKind = "WorkspaceTaskIntegrationReadback@7"
const refreshIntegrationReadbackSchema = 7

func integrationReadbackVersion(plan WorkspaceTaskIntegrationPlan) (string, int) {
	if refreshIntegrationPlan(plan) {
		return refreshIntegrationReadbackKind, refreshIntegrationReadbackSchema
	}
	return integrationReadbackKind, integrationReadbackSchema
}

func integrationPersistedView(registry WorkItemRegistry, authority *IntegrationAuthority, intent *IntegrationIntent, attempt *IntegrationAttempt, result *IntegrationResult) canonicaljson.Value {
	var observedHash canonicaljson.Value
	freshness := "unknown"
	if digestPattern.MatchString(registry.RawSHA256) {
		observedHash, freshness = registry.RawSHA256, "fresh"
	}
	return canonicaljson.Object{{Name: "format_version", Value: int64(registry.FormatVersion)}, {Name: "sha256", Value: observedHash}, {Name: "authority_id", Value: nullableID(authority)}, {Name: "intent_id", Value: nullableID(intent)}, {Name: "attempt_id", Value: nullableID(attempt)}, {Name: "integration_result_id", Value: nullableID(result)}, {Name: "freshness", Value: freshness}}
}

// Task show reads the Task's latest integration authority, or otherwise the
// current Task worktree, through the same check as integrate --check.
func buildTaskShowIntegrationReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding) WorkspaceTaskIntegrationReadback {
	authority, ambiguous := selectTaskIntegrationLeaf(registry, task.ID)
	if ambiguous {
		return newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "conflict", IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: "Several integration attempts compete for this Task; inspect them before a new check.", Argv: []string{}})
	}
	var intent *IntegrationIntent
	var attempt *IntegrationAttempt
	var stored *IntegrationResult
	var input TaskIntegrationInput
	if authority != nil {
		intent, attempt, stored = findIntentForAuthority(registry, authority.ID), findAttemptForAuthority(registry, authority.ID), findIntegrationResultForAuthority(registry, authority.ID)
		input = TaskIntegrationInput{TaskID: task.ID, ExpectedResultOID: authority.Plan.Task.ResultOID, ExpectedParentOID: authority.Plan.Epic.ExpectedParentOID, RetryAfterResultID: authority.RetryAfterResultID, RefreshParent: refreshIntegrationPlan(authority.Plan)}
	} else if d.IntegrationGit != nil && task.Worktree != nil {
		if source, err := d.IntegrationGit.ObserveIntegrationWorktree(task.Worktree.Locator, task.Worktree.Ref); err == nil {
			input = TaskIntegrationInput{TaskID: task.ID, ExpectedResultOID: source.OID, ExpectedParentOID: task.Worktree.ParentOID}
		}
	}
	if input.ExpectedResultOID != "" && d.IntegrationGit != nil && d.Files != nil {
		if plan, digest, ctx, err := buildIntegrationPlan(d, root, projects, registry, input); err == nil {
			return newIntegrationReadback(root, registry, ctx, plan, digest, authority, intent, attempt, stored)
		}
	}
	if authority != nil {
		return newPersistedIntegrationFallbackReadback(d, root, projects, registry, task, epic, parent, authority, intent, attempt, stored)
	}
	return newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "blocked", IntegrationNextAction{Kind: "resolve_blocker", Reason: "Refresh the local repository observations before integration.", Argv: []string{}})
}

func newPersistedIntegrationFallbackReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding, authority *IntegrationAuthority, intent *IntegrationIntent, attempt *IntegrationAttempt, result *IntegrationResult) WorkspaceTaskIntegrationReadback {
	classification := "blocked"
	recovery := "blocked"
	next := IntegrationNextAction{Kind: "resolve_blocker", Reason: "Refresh the local repository observations before integration.", Argv: []string{}}
	var changed *bool
	if result != nil {
		classification = result.Outcome
		recovery = result.RecoveryStatus
		next = result.NextAction
		changed = result.GitChanged
	}
	readback := newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, classification, next)
	readback.Classification = classification
	readback.GitChanged = changed
	readback.RecoveryStatus = recovery
	readback.NextAction = next
	readback.ResultOID = authority.Plan.Task.ResultOID
	if result != nil && result.BeforeObservation.OID != "" {
		readback.ParentOID = result.BeforeObservation.OID
	}
	if refreshIntegrationPlan(authority.Plan) {
		kind, schema := integrationReadbackVersion(authority.Plan)
		replaceReadbackMember(&readback.Value, "kind", kind)
		replaceReadbackMember(&readback.Value, "schema_version", int64(schema))
		replaceReadbackMember(&readback.Value, "epic", structCanonical(authority.Plan.Epic))
		replaceReadbackMember(&readback.Value, "task", structCanonical(authority.Plan.Task))
		replaceReadbackMember(&readback.Value, "source", structCanonical(authority.Plan.ObservedSource))
		if result == nil || result.BeforeObservation.OID == "" {
			readback.ParentOID = authority.Plan.Epic.ExpectedParentOID
		}
	}
	replaceReadbackMember(&readback.Value, "authority", nullableStruct(authority))
	replaceReadbackMember(&readback.Value, "intent", nullableStruct(intent))
	replaceReadbackMember(&readback.Value, "attempt", nullableStruct(attempt))
	replaceReadbackMember(&readback.Value, "integration_result", nullableStruct(result))
	persisted := integrationPersistedView(registry, authority, intent, attempt, result)
	replaceReadbackMember(&readback.Value, "persisted_before", persisted)
	replaceReadbackMember(&readback.Value, "persisted_after", persisted)
	replaceReadbackMember(&readback.Value, "classification", classification)
	replaceReadbackMember(&readback.Value, "git_changed", pointerValue(changed))
	replaceReadbackMember(&readback.Value, "recovery_status", recovery)
	replaceReadbackMember(&readback.Value, "next_action", structCanonical(next))
	return readback
}

func selectTaskIntegrationLeaf(registry WorkItemRegistry, taskID TaskID) (*IntegrationAuthority, bool) {
	authorities := []*IntegrationAuthority{}
	predecessors := map[IntegrationAuthorityID]bool{}
	for i := range registry.IntegrationAuthorities {
		a := &registry.IntegrationAuthorities[i]
		if a.TaskID != taskID || authorityRetired(registry, a) {
			continue
		}
		authorities = append(authorities, a)
		if a.RetryAfterResultID != nil {
			if prior := findIntegrationResult(registry, *a.RetryAfterResultID); prior != nil {
				predecessors[prior.AuthorityID] = true
			}
		}
	}
	leaves := []*IntegrationAuthority{}
	for _, a := range authorities {
		if !predecessors[a.ID] {
			leaves = append(leaves, a)
		}
	}
	if len(leaves) > 1 {
		return nil, true
	}
	if len(leaves) == 0 {
		return nil, false
	}
	return leaves[0], false
}

func newNullableIntegrationReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding, classification string, next IntegrationNextAction) WorkspaceTaskIntegrationReadback {
	persisted := integrationPersistedView(registry, nil, nil, nil, nil)
	project := IntegrationPlanProject{ProjectID: task.ProjectID}
	repository := IntegrationPlanRepository{RepoID: task.RepoID, GitCommonDir: task.GitCommonDir}
	if _, repo, err := projectAndRepo(projects, task.ProjectID, task.RepoID); err == nil {
		repository.RegisteredLocator = repo.Locator
	}
	var markerDigest canonicaljson.Value
	if d.Files != nil {
		if marker, err := d.Files.ReadFile(filepath.Join(root, MarkerDirectory, MarkerFile)); err == nil {
			markerDigest = digestTaskBytes(marker)
		}
	}
	epicValue := IntegrationPlanEpic{EpicID: epic.ID, ParentWorktreeID: parent.Worktree.ID, ParentLocator: parent.Worktree.Locator, ParentRef: parent.Worktree.Ref, ExpectedParentOID: parent.Worktree.OID, ExpectedParentTree: parent.Worktree.Tree}
	taskValue := IntegrationPlanTask{TaskID: task.ID}
	var source canonicaljson.Value
	if task.Worktree != nil {
		taskValue = IntegrationPlanTask{TaskID: task.ID, TaskWorktreeID: task.Worktree.ID, SourceLocator: task.Worktree.Locator, SourceRef: task.Worktree.Ref, ResultOID: task.Worktree.OID, ResultTree: task.Worktree.Tree}
		source = structCanonical(PlanWorktreeObservation{WorktreeLocator: task.Worktree.Locator, Ref: task.Worktree.Ref, OID: task.Worktree.OID, Tree: task.Worktree.Tree, GitCommonDir: task.Worktree.GitCommonDir})
	}
	value := canonicaljson.Object{
		{Name: "kind", Value: integrationReadbackKind}, {Name: "schema_version", Value: int64(integrationReadbackSchema)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "workspace", Value: canonicaljson.Object{{Name: "root", Value: root}, {Name: "marker_sha256", Value: markerDigest}}}, {Name: "project", Value: structCanonical(project)}, {Name: "repository", Value: structCanonical(repository)}, {Name: "epic", Value: structCanonical(epicValue)}, {Name: "task", Value: structCanonical(taskValue)}, {Name: "source", Value: source},
		{Name: "plan", Value: nil}, {Name: "authority", Value: nil}, {Name: "intent", Value: nil}, {Name: "attempt", Value: nil}, {Name: "integration_result", Value: nil},
		{Name: "persisted_before", Value: persisted}, {Name: "persisted_after", Value: persisted}, {Name: "observed_source", Value: observedWorktreeReadback(IntegrationWorktreeObservation{}, false)}, {Name: "observed_parent", Value: observedWorktreeReadback(IntegrationWorktreeObservation{}, false)}, {Name: "observed_inventory", Value: canonicaljson.Object{{Name: "entries", Value: []canonicaljson.Value{}}, {Name: "freshness", Value: "unknown"}}}, {Name: "observed_reflog", Value: canonicaljson.Object{{Name: "entries", Value: []canonicaljson.Value{}}, {Name: "freshness", Value: "unknown"}}},
		{Name: "classification", Value: classification}, {Name: "git_changed", Value: nil}, {Name: "recovery_status", Value: "blocked"}, {Name: "next_action", Value: structCanonical(next)},
	}
	return WorkspaceTaskIntegrationReadback{Value: value, Classification: classification, RecoveryStatus: "blocked", NextAction: next, TaskID: task.ID, EpicID: epic.ID, ParentRef: parent.Worktree.Ref, ParentOID: parent.Worktree.OID}
}

func replaceReadbackMember(object *canonicaljson.Object, name string, value canonicaljson.Value) {
	for i := range *object {
		if (*object)[i].Name == name {
			(*object)[i].Value = value
			return
		}
	}
}

func integrationPlanNextAction(ctx integrationContext, plan WorkspaceTaskIntegrationPlan, digest string) IntegrationNextAction {
	switch plan.Readiness {
	case "ready":
		argv := append(integrationArgv(ctx.Task.ID, plan), "--expected-result-oid", ctx.ResultOID, "--expected-parent-oid", plan.Epic.ExpectedParentOID)
		if plan.RetryAfterResultID != nil {
			argv = append(argv, "--retry-after", string(*plan.RetryAfterResultID))
		}
		return IntegrationNextAction{Kind: "apply_confirmed_plan", Reason: "Apply the exact confirmed local fast-forward.", Argv: append(argv, "--apply", "--confirm", digest)}
	case "already_integrated":
		return IntegrationNextAction{Kind: "none", Reason: "The Task is already integrated.", Argv: []string{}}
	case "conflict":
		return IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: "Resolve the conflicting exact binding in a separate activity.", Argv: []string{}}
	default:
		return IntegrationNextAction{Kind: "resolve_blocker", Reason: "Resolve the named prerequisite before a new check.", Argv: []string{}}
	}
}

func integrationCurrentGuidance(registry WorkItemRegistry, ctx integrationContext, plan WorkspaceTaskIntegrationPlan, classification, recovery string, next IntegrationNextAction) IntegrationNextAction {
	if classification == "partial" || classification == "unknown" || recovery == "reconciliation-required" || next.Kind == "read_only_recovery_control" {
		return next
	}
	if (classification == "exact_effect" || classification == "already_integrated") && recovery == "complete" {
		return integrationBaseGuidance(registry, ctx, classification, recovery, next)
	}
	if classification == "conflict" && plan.Readiness != "conflict" && plan.Readiness != "blocked" {
		plan.Readiness = "conflict"
		plan.Reasons = nil
		if ctx.SourceObserved && ctx.Source.OID != plan.Task.ResultOID {
			plan.Reasons = append(plan.Reasons, "source_binding_mismatch")
		}
		if ctx.ParentObserved && (ctx.ParentObservation.OID != plan.Epic.ExpectedParentOID || ctx.ParentObservation.Tree != plan.Epic.ExpectedParentTree) {
			plan.Reasons = append(plan.Reasons, "parent_moved")
		}
	}
	if plan.Readiness != "conflict" && plan.Readiness != "blocked" {
		return integrationBaseGuidance(registry, ctx, classification, recovery, next)
	}
	return integrationConflictGuidance(registry, ctx, plan)
}

type integrationReasonSet map[string]bool

func (reasons integrationReasonSet) any(names ...string) bool {
	for _, name := range names {
		if reasons[name] {
			return true
		}
	}
	return false
}

func integrationStop(reason string) IntegrationNextAction {
	return IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: reason, Argv: []string{}}
}

func integrationConflictGuidance(registry WorkItemRegistry, ctx integrationContext, plan WorkspaceTaskIntegrationPlan) IntegrationNextAction {
	reasons := integrationReasonSet{}
	for _, reason := range plan.Reasons {
		reasons[reason] = true
	}
	if reasons.any("prior_integration_requires_reconciliation", "parent_integration_already_active", "prior_no_effect_requires_retry", "retry_predecessor_already_used", "invalid_retry_predecessor", "retry_binding_mismatch") {
		return integrationPriorGuidance(registry, ctx, reasons)
	}
	if reasons.any("epic_base_update_pending") {
		return integrationStop(fmt.Sprintf("Epic %s repository %s already has a pending base update. Complete or reconcile that operation before planning Task %s integration.", ctx.Epic.ID, ctx.Parent.RepoID, ctx.Task.ID))
	}
	if reasons.any("source_inventory_mismatch", "parent_inventory_mismatch", "repository_binding_mismatch", "object_id_length_mismatch", "unsupported_object_format", "unsupported_ref_format", "unsupported_repository_state") {
		return integrationStop(fmt.Sprintf("Task %s and Epic %s repository/worktree identity or inventory does not match the observed Git state. Resolve the named plan reasons (%s), then run check again.", ctx.Task.ID, ctx.Epic.ID, strings.Join(plan.Reasons, ", ")))
	}
	if reasons.any("source_not_clean", "parent_not_clean", "result_not_descendant") {
		return integrationStop(fmt.Sprintf("Task %s source or Epic %s parent is dirty, in progress, or the result is not descended from the expected parent (%s). Resolve the Git condition and run check again.", ctx.Task.ID, ctx.Epic.ID, strings.Join(plan.Reasons, ", ")))
	}
	return integrationObservedConflictGuidance(registry, ctx, plan, reasons)
}

func integrationPriorGuidance(registry WorkItemRegistry, ctx integrationContext, reasons integrationReasonSet) IntegrationNextAction {
	if reasons.any("parent_integration_already_active") {
		blockers := []IntegrationAuthority{}
		for _, authority := range registry.IntegrationAuthorities {
			if authority.Plan.Repository.GitCommonDir == ctx.Task.GitCommonDir && authority.Plan.Epic.ParentRef == ctx.Parent.Worktree.Ref && findIntegrationResultForAuthority(registry, authority.ID) == nil {
				blockers = append(blockers, authority)
			}
		}
		if len(blockers) == 1 {
			blocker := blockers[0]
			return IntegrationNextAction{Kind: "read_only_recovery_control", Reason: fmt.Sprintf("Task %s is blocked by Task %s integration authority %s on parent %s. Inspect that Task before another check or apply.", ctx.Task.ID, blocker.TaskID, blocker.ID, ctx.Parent.Worktree.Ref), Argv: []string{"ply", "workspace", "task", "show", string(blocker.TaskID), "--format", "json"}}
		}
		if len(blockers) > 1 {
			ids := []string{}
			for _, blocker := range blockers {
				ids = append(ids, fmt.Sprintf("%s/%s", blocker.TaskID, blocker.ID))
			}
			sort.Strings(ids)
			return integrationStop(fmt.Sprintf("Task %s has several active authorities on parent %s: %s. Resolve their identities before a new check.", ctx.Task.ID, ctx.Parent.Worktree.Ref, strings.Join(ids, ", ")))
		}
	}
	ids := []string{}
	for _, authority := range registry.IntegrationAuthorities {
		if authority.TaskID == ctx.Task.ID && !authorityRetired(registry, &authority) {
			ids = append(ids, string(authority.ID))
		}
	}
	for _, result := range registry.IntegrationResults {
		if authority := findAuthority(registry, result.AuthorityID); authority != nil && authority.TaskID == ctx.Task.ID {
			ids = append(ids, string(result.ID))
		}
	}
	sort.Strings(ids)
	return IntegrationNextAction{Kind: "read_only_recovery_control", Reason: fmt.Sprintf("Task %s has an unresolved or competing integration authority/result (%s). Inspect it before another check or apply.", ctx.Task.ID, strings.Join(ids, ", ")), Argv: []string{"ply", "workspace", "task", "show", string(ctx.Task.ID), "--format", "json"}}
}

func integrationObservedConflictGuidance(registry WorkItemRegistry, ctx integrationContext, plan WorkspaceTaskIntegrationPlan, reasons integrationReasonSet) IntegrationNextAction {
	if !ctx.SourceObserved || !ctx.ParentObserved || !ctx.RepositoryObserved || !digestPattern.MatchString(registry.RawSHA256) {
		return integrationStop(fmt.Sprintf("Fresh source, parent, repository and registry observations are required for Task %s guidance; run check after restoring those observations.", ctx.Task.ID))
	}
	if ctx.PreconditionReason != "" {
		return integrationStop(fmt.Sprintf("Task %s confirmation precondition changed: %s Reconcile that condition and run check again.", ctx.Task.ID, ctx.PreconditionReason))
	}
	if ctx.Source.Ref != ctx.Task.Worktree.Ref || ctx.Source.GitCommonDir != ctx.Task.GitCommonDir || ctx.Source.Locator != ctx.Task.Worktree.Locator {
		return integrationStop(fmt.Sprintf("Task %s source locator, ref or common directory differs from its recorded worktree. Resolve that identity mismatch and run check again.", ctx.Task.ID))
	}
	if ctx.ParentObservation.Locator != ctx.Parent.Worktree.Locator || ctx.ParentObservation.Ref != ctx.Parent.Worktree.Ref || ctx.ParentObservation.GitCommonDir != ctx.Parent.GitCommonDir || ctx.Repo.GitCommonDir != ctx.Parent.GitCommonDir {
		return integrationStop(fmt.Sprintf("Epic %s parent binding (locator %s, ref %s, common directory %s) differs from the observed parent (locator %s, ref %s, common directory %s). Resolve parent binding before a new check.", ctx.Epic.ID, ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref, ctx.Parent.GitCommonDir, ctx.ParentObservation.Locator, ctx.ParentObservation.Ref, ctx.ParentObservation.GitCommonDir))
	}
	if !ctx.Source.Symbolic || !ctx.ParentObservation.Symbolic {
		return integrationStop(fmt.Sprintf("Task %s source and Epic %s parent must both have symbolic refs before a new check.", ctx.Task.ID, ctx.Epic.ID))
	}
	if (ctx.Repository.ObjectFormat != "sha1" && ctx.Repository.ObjectFormat != "sha256") || ctx.Repository.RefFormat != "files" {
		return integrationStop(fmt.Sprintf("Task %s repository format %s/%s is unsupported for integration guidance.", ctx.Task.ID, ctx.Repository.ObjectFormat, ctx.Repository.RefFormat))
	}
	if !objectIDsMatchFormat(ctx.Repository.ObjectFormat, ctx.Source.OID, ctx.Source.Tree, ctx.ParentObservation.OID, ctx.ParentObservation.Tree) {
		return integrationStop(fmt.Sprintf("Task %s observed source or Epic %s parent OID does not match the repository object format. Resolve object identity before a new check.", ctx.Task.ID, ctx.Epic.ID))
	}
	if ctx.Repository.Shallow || ctx.Repository.PartialClone || ctx.Repository.SparseCheckout {
		return integrationStop(fmt.Sprintf("Task %s repository state is shallow, partial clone, or sparse checkout. Resolve that state before a new check.", ctx.Task.ID))
	}
	if !ctx.Source.Clean || !ctx.ParentObservation.Clean || len(ctx.Source.InProgress) != 0 || len(ctx.ParentObservation.InProgress) != 0 {
		return integrationStop(fmt.Sprintf("Task %s source %s or Epic %s parent %s is dirty or has Git work in progress. Resolve it and run check again.", ctx.Task.ID, ctx.Source.OID, ctx.Epic.ID, ctx.ParentObservation.OID))
	}
	if !integrationInventoryContains(ctx.Repository.Inventory, ctx.Source) || !integrationInventoryContains(ctx.Repository.Inventory, ctx.ParentObservation) || ctx.Source.ObjectFormat != ctx.Repository.ObjectFormat || ctx.ParentObservation.ObjectFormat != ctx.Repository.ObjectFormat || ctx.Source.RefFormat != ctx.Repository.RefFormat || ctx.ParentObservation.RefFormat != ctx.Repository.RefFormat {
		return integrationStop(fmt.Sprintf("Task %s source or Epic %s parent does not match the fresh repository inventory and format. Resolve worktree identity before a new check.", ctx.Task.ID, ctx.Epic.ID))
	}
	return integrationConflictNextStep(registry, ctx, plan, reasons)
}

func integrationConflictNextStep(registry WorkItemRegistry, ctx integrationContext, plan WorkspaceTaskIntegrationPlan, reasons integrationReasonSet) IntegrationNextAction {
	if basePreviewQualified(registry, ctx) {
		base := currentEpicBase(registry, ctx.Epic, ctx.Parent)
		if base.OID != ctx.ParentObservation.OID || base.Tree != ctx.ParentObservation.Tree {
			return IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: fmt.Sprintf("The recorded Epic base %s differs from the observed parent %s. Preview the base update; this does not refresh Task %s's original basis.", base.OID, ctx.ParentObservation.OID, ctx.Task.ID), Argv: basePreviewArgv(ctx)}
		}
	}
	if reasons.any("refresh_parent_not_registered_base") && plan.RetryAfterResultID != nil {
		return integrationStop(fmt.Sprintf("Task %s no-effect retry %s keeps its predecessor basis %s, but the recorded Epic base is %s. A retry never selects another basis; restore that parent or choose a supported continuation.", ctx.Task.ID, *plan.RetryAfterResultID, plan.Epic.ExpectedParentOID, currentEpicBase(registry, ctx.Epic, ctx.Parent).OID))
	}
	if reasons.any("refresh_parent_unchanged") {
		return integrationStop(fmt.Sprintf("Task %s already has basis %s; check integration without --refresh-parent.", ctx.Task.ID, ctx.Task.Worktree.ParentOID))
	}
	if reasons.any("epic_base_changed") {
		base := currentEpicBase(registry, ctx.Epic, ctx.Parent).OID
		return integrationStop(fmt.Sprintf("Task %s was created from %s; the current recorded Epic base is %s. Adapt the Task so %s is an ancestor of its commit, then run `ply workspace task integrate %s --refresh-parent --expected-result-oid <adapted commit> --expected-parent-oid %s --check`. The Task's original basis stays recorded.", ctx.Task.ID, ctx.Task.Worktree.ParentOID, base, base, ctx.Task.ID, base))
	}
	if reasons.any("parent_moved") && (ctx.ParentObservation.OID != ctx.Task.Worktree.ParentOID || ctx.ParentObservation.Tree != ctx.Task.Worktree.ParentTree) {
		return integrationStop(fmt.Sprintf("Epic %s parent is now %s, while Task %s was created from %s. Settle the recorded Epic base first; --refresh-parent integrates only from the registered base.", ctx.Epic.ID, ctx.ParentObservation.OID, ctx.Task.ID, ctx.Task.Worktree.ParentOID))
	}
	if reasons.any("source_binding_mismatch", "expected_oid_mismatch") && ctx.ParentObservation.OID == ctx.Task.Worktree.ParentOID && ctx.ParentObservation.Tree == ctx.Task.Worktree.ParentTree && validOIDText(ctx.Source.OID) && validOIDText(ctx.ParentObservation.OID) {
		return IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: fmt.Sprintf("The supplied OID for Task %s is stale. Recheck the observed source %s and parent %s to verify ancestry and other preconditions for a new plan; discard the old digest.", ctx.Task.ID, ctx.Source.OID, ctx.ParentObservation.OID), Argv: []string{"ply", "workspace", "task", "integrate", string(ctx.Task.ID), "--expected-result-oid", ctx.Source.OID, "--expected-parent-oid", ctx.ParentObservation.OID, "--check", "--format", "json"}}
	}
	return integrationStop(fmt.Sprintf("Task %s integration has unresolved plan reasons (%s). Resolve them and run check again; no safe command can be inferred from these observations.", ctx.Task.ID, strings.Join(plan.Reasons, ", ")))
}

func nullableStruct(v any) canonicaljson.Value {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return nil
	}
	return structCanonical(v)
}
func nullableID(v any) canonicaljson.Value {
	switch x := v.(type) {
	case *IntegrationAuthority:
		if x != nil {
			return string(x.ID)
		}
	case *IntegrationIntent:
		if x != nil {
			return string(x.ID)
		}
	case *IntegrationAttempt:
		if x != nil {
			return string(x.ID)
		}
	case *IntegrationResult:
		if x != nil {
			return string(x.ID)
		}
	}
	return nil
}
func observedWorktreeReadback(observation IntegrationWorktreeObservation, observed bool) canonicaljson.Object {
	values := canonicaljson.Value(unknownWorktreeObservation())
	if observed {
		values = structCanonical(planObservation(observation))
	}
	return canonicaljson.Object{{Name: "values", Value: values}, {Name: "freshness", Value: observationFreshness(observed)}}
}

func unknownWorktreeObservation() canonicaljson.Object {
	return canonicaljson.Object{
		{Name: "worktree_locator", Value: nil}, {Name: "ref", Value: nil}, {Name: "oid", Value: nil}, {Name: "tree", Value: nil},
		{Name: "git_common_dir", Value: nil}, {Name: "object_format", Value: nil}, {Name: "ref_format", Value: nil}, {Name: "symbolic", Value: nil},
		{Name: "clean", Value: nil}, {Name: "status_entries", Value: []canonicaljson.Value{}}, {Name: "in_progress", Value: []canonicaljson.Value{}},
	}
}

func observationFreshness(observed bool) canonicaljson.Value {
	value := "unknown"
	if observed {
		value = "fresh"
	}
	return canonicaljson.Object{{Name: "worktree_locator", Value: value}, {Name: "ref", Value: value}, {Name: "oid", Value: value}, {Name: "tree", Value: value}, {Name: "git_common_dir", Value: value}, {Name: "object_format", Value: value}, {Name: "ref_format", Value: value}, {Name: "symbolic", Value: value}, {Name: "clean", Value: value}, {Name: "status_entries", Value: value}, {Name: "in_progress", Value: value}}
}

func collectionFreshness(known bool) string {
	if known {
		return "fresh"
	}
	return "unknown"
}

func RenderTaskIntegrationText(r WorkspaceTaskIntegrationReadback) string {
	status := map[string]string{"ready": "The exact integration plan is ready; no write was made.", "exact_effect": "The local Epic parent was fast-forwarded to the exact Task result.", "already_integrated": "The local Epic parent already contains the exact Task result; no Git attempt was made.", "no_effect": "The Git attempt made no observable change.", "blocked": "Integration is blocked by a prerequisite.", "conflict": "Integration conflicts with persisted or observed state.", "partial": "A non-exact Git effect was observed.", "unknown": "The integration effect could not be classified safely."}[r.Classification]
	if status == "" {
		status = r.Classification
	}
	changed := "unknown"
	if r.GitChanged != nil {
		if *r.GitChanged {
			changed = "yes"
		} else {
			changed = "no"
		}
	}
	next := r.NextAction.Reason
	if r.NextAction.Kind == "none" {
		next = "No action now; the Task is integrated."
	} else if len(r.NextAction.Argv) > 0 {
		next += " Run `" + strings.Join(r.NextAction.Argv, " ") + "`."
	}
	return fmt.Sprintf("Status: %s\nTask: %s / %s\nParent: %s / %s@%s\nGit changed: %s\nRecovery: %s\nNext action: %s\n", status, r.TaskID, r.ResultOID, r.EpicID, r.ParentRef, r.ParentOID, changed, r.RecoveryStatus, next)
}

func canonicalFromReflect(value any) canonicaljson.Value {
	return canonicalReflectValue(reflect.ValueOf(value))
}
func canonicalReflectValue(v reflect.Value) canonicaljson.Value {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		return canonicalReflectValue(v.Elem())
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return canonicalReflectValue(v.Elem())
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Slice, reflect.Array:
		r := make([]canonicaljson.Value, v.Len())
		for i := 0; i < v.Len(); i++ {
			r[i] = canonicalReflectValue(v.Index(i))
		}
		return r
	case reflect.Struct:
		t := v.Type()
		members := make(canonicaljson.Object, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}
			tag := field.Tag.Get("json")
			name := strings.Split(tag, ",")[0]
			if name == "" {
				name = field.Name
			}
			if name == "-" {
				continue
			}
			members = append(members, canonicaljson.Member{Name: name, Value: canonicalReflectValue(v.Field(i))})
		}
		sort.Slice(members, func(i, j int) bool { return members[i].Name < members[j].Name })
		return members
	}
	return nil
}

func integrationBaseGuidance(registry WorkItemRegistry, ctx integrationContext, classification, recovery string, next IntegrationNextAction) IntegrationNextAction {
	if (classification != "exact_effect" && classification != "already_integrated") || recovery != "complete" || !basePreviewQualified(registry, ctx) {
		return next
	}
	bases := epicBases(registry, ctx.Epic, ctx.Parent)
	if len(bases) == 0 {
		return next
	}
	if base := bases[len(bases)-1]; base.OID == ctx.ParentObservation.OID && base.Tree == ctx.ParentObservation.Tree {
		return next
	}
	return IntegrationNextAction{Kind: "new_check", Reason: "The Task is integrated. The recorded Epic base differs from the observed parent. Preview the base update before creating another Task.", Argv: basePreviewArgv(ctx)}
}

func basePreviewQualified(registry WorkItemRegistry, ctx integrationContext) bool {
	if !digestPattern.MatchString(registry.RawSHA256) || !ctx.ParentObserved || !qualifiedBaseParentObservation(ctx) {
		return false
	}
	if _, ambiguous := selectTaskIntegrationLeaf(registry, ctx.Task.ID); ambiguous {
		return false
	}
	target := QueueTarget{ctx.Epic.ProjectID, ctx.Parent.RepoID, ctx.Epic.ID, ctx.Parent.GitCommonDir, ctx.Parent.Worktree.ID, ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref}
	return len(queueUnresolved(registry, target)) == 0 && !pendingBaseUpdate(registry, ctx.Epic.ID, ctx.Parent.RepoID)
}

func qualifiedBaseParentObservation(ctx integrationContext) bool {
	parent := ctx.ParentObservation
	if parent.Locator != ctx.Parent.Worktree.Locator || parent.Ref != ctx.Parent.Worktree.Ref || parent.GitCommonDir != ctx.Parent.GitCommonDir || ctx.Repo.GitCommonDir != ctx.Parent.GitCommonDir || parent.ObjectFormat != ctx.Repository.ObjectFormat || (parent.ObjectFormat != "sha1" && parent.ObjectFormat != "sha256") || parent.RefFormat != "files" || !parent.Clean || !parent.Symbolic || len(parent.InProgress) != 0 {
		return false
	}
	return true
}

func basePreviewArgv(ctx integrationContext) []string {
	return []string{"ply", "workspace", "epic", "base", "update", string(ctx.Epic.ID), "--repo", string(ctx.Parent.RepoID), "--check", "--format", "json"}
}
