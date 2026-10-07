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
	Value          canonicaljson.Object
	TaskResult     *TaskResultRecord
	HumanQA        *TaskHumanQARecord
	Classification string
	GitChanged     *bool
	RecoveryStatus string
	NextAction     IntegrationNextAction
	TaskID         TaskID
	ResultOID      string
	TechnicalGate  string
	HumanQAOutcome string
	EpicID         EpicID
	ParentRef      string
	ParentOID      string
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
	readback := canonicaljson.Object{
		{Name: "kind", Value: "WorkspaceTaskIntegrationReadback@1"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "workspace", Value: structCanonical(plan.Workspace)}, {Name: "project", Value: structCanonical(plan.Project)}, {Name: "repository", Value: structCanonical(plan.Repository)}, {Name: "epic", Value: structCanonical(plan.Epic)}, {Name: "task", Value: structCanonical(plan.Task)}, {Name: "source", Value: structCanonical(plan.ObservedSource)},
		{Name: "task_result", Value: structCanonical(plan.TaskResult)}, {Name: "human_qa", Value: structCanonical(plan.HumanQA)}, {Name: "plan", Value: planValue}, {Name: "authority", Value: authorityValue}, {Name: "intent", Value: intentValue}, {Name: "attempt", Value: attemptValue}, {Name: "integration_result", Value: resultValue},
		{Name: "persisted_before", Value: persistedBefore}, {Name: "persisted_after", Value: persistedAfter}, {Name: "observed_source", Value: observedSource}, {Name: "observed_parent", Value: observedParent}, {Name: "observed_inventory", Value: canonicaljson.Object{{Name: "entries", Value: structCanonical(observedInventory)}, {Name: "freshness", Value: collectionFreshness(ctx.RepositoryObserved)}}}, {Name: "observed_reflog", Value: canonicaljson.Object{{Name: "entries", Value: structCanonical(observedReflog)}, {Name: "freshness", Value: collectionFreshness(ctx.ReflogObserved)}}},
		{Name: "classification", Value: classification}, {Name: "git_changed", Value: pointerValue(changed)}, {Name: "recovery_status", Value: recovery}, {Name: "next_action", Value: structCanonical(next)},
	}
	if plan.SchemaVersion >= 2 || taskRequiresSpec(registry, ctx.Task.ID) {
		guard := ctx.TaskSpecRelevance
		if guard == nil {
			guard = plan.TaskSpecGuard
		}
		if guard == nil {
			guard = &TaskSpecRelevance{Relevance: "unknown", Reasons: []string{"task_spec_result_basis_missing"}}
		}
		replaceReadbackMember(&readback, "kind", "WorkspaceTaskIntegrationReadback@2")
		replaceReadbackMember(&readback, "schema_version", int64(2))
		readback = append(readback, canonicaljson.Member{Name: "task_spec_relevance", Value: contentRefValue(guard)})
	}
	parentOID := plan.Epic.ExpectedParentOID
	if ctx.ParentObserved && ctx.ParentObservation.OID != "" {
		parentOID = ctx.ParentObservation.OID
	} else if result != nil && result.BeforeObservation.OID != "" {
		parentOID = result.BeforeObservation.OID
	} else if plan.ObservedParent.OID != "" {
		parentOID = plan.ObservedParent.OID
	}
	return WorkspaceTaskIntegrationReadback{Value: readback, TaskResult: &ctx.Result, HumanQA: &ctx.QA, Classification: classification, GitChanged: changed, RecoveryStatus: recovery, NextAction: next, TaskID: ctx.Task.ID, ResultOID: ctx.Result.ResultOID, TechnicalGate: ctx.Result.TechnicalGate, HumanQAOutcome: ctx.QA.Outcome, EpicID: ctx.Epic.ID, ParentRef: ctx.Parent.Worktree.Ref, ParentOID: parentOID}
}

func integrationPersistedView(registry WorkItemRegistry, authority *IntegrationAuthority, intent *IntegrationIntent, attempt *IntegrationAttempt, result *IntegrationResult) canonicaljson.Value {
	storeBytes, _ := encodeWorkItemRegistry(registry)
	return canonicaljson.Object{{Name: "format_version", Value: int64(registry.FormatVersion)}, {Name: "sha256", Value: digestTaskBytes(storeBytes)}, {Name: "authority_id", Value: nullableID(authority)}, {Name: "intent_id", Value: nullableID(intent)}, {Name: "attempt_id", Value: nullableID(attempt)}, {Name: "integration_result_id", Value: nullableID(result)}, {Name: "freshness", Value: "fresh"}}
}

func buildTaskShowIntegrationReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding) (out WorkspaceTaskIntegrationReadback) {
	defer func() {
		if !taskRequiresSpec(registry, task.ID) || out.Value == nil {
			return
		}
		var id TaskResultID
		if out.TaskResult != nil {
			id = out.TaskResult.ID
		}
		guard := taskResultSpecRelevance(d, root, projects, registry, task, id)
		replaceReadbackMember(&out.Value, "kind", "WorkspaceTaskIntegrationReadback@2")
		replaceReadbackMember(&out.Value, "schema_version", int64(2))
		if contentFields(out.Value)["task_spec_relevance"] == nil {
			out.Value = append(out.Value, canonicaljson.Member{Name: "task_spec_relevance", Value: contentRefValue(guard)})
		} else {
			replaceReadbackMember(&out.Value, "task_spec_relevance", contentRefValue(guard))
		}
	}()

	tr, qa, authority, ambiguous := selectTaskIntegrationLeaf(registry, task.ID)
	if ambiguous {
		return newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "conflict", IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: "Select exact Task result and human QA records in a separate check.", Argv: []string{}})
	}
	if tr == nil {
		return newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "blocked", IntegrationNextAction{Kind: "record_task_result", Reason: "Record one controlled Task result.", Argv: []string{"ply", "workspace", "task", "result", "record", string(task.ID), "--file", "<absolute-task-result.json>"}})
	}
	if qa == nil {
		r := newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "blocked", IntegrationNextAction{Kind: "record_human_qa", Reason: "Record human QA for the selected Task result.", Argv: []string{"ply", "workspace", "task", "qa", "record", string(task.ID), "--file", "<absolute-human-qa.json>"}})
		r.TaskResult = tr
		r.ResultOID = tr.ResultOID
		r.TechnicalGate = tr.TechnicalGate
		replaceReadbackMember(&r.Value, "task_result", taskResultCanonical(*tr))
		return r
	}
	input := TaskIntegrationInput{TaskID: task.ID, TaskResultID: tr.ID, HumanQARecordID: qa.ID, ExpectedResultOID: tr.ResultOID, ExpectedParentOID: task.Worktree.ParentOID}
	if authority != nil {
		input.RetryAfterResultID = authority.RetryAfterResultID
	}
	var intent *IntegrationIntent
	var attempt *IntegrationAttempt
	var integrationResult *IntegrationResult
	if authority != nil {
		intent = findIntentForAuthority(registry, authority.ID)
		attempt = findAttemptForAuthority(registry, authority.ID)
		integrationResult = findIntegrationResultForAuthority(registry, authority.ID)
	}
	if d.IntegrationGit != nil && d.Files != nil && task.Worktree != nil {
		if plan, digest, ctx, err := buildIntegrationPlan(d, root, projects, registry, input); err == nil {
			return newIntegrationReadback(root, registry, ctx, plan, digest, authority, intent, attempt, integrationResult)
		}
	}
	if authority != nil {
		return newPersistedIntegrationFallbackReadback(d, root, projects, registry, task, epic, parent, tr, qa, authority, intent, attempt, integrationResult)
	}
	r := newNullableIntegrationReadback(d, root, projects, registry, task, epic, parent, "blocked", IntegrationNextAction{Kind: "resolve_blocker", Reason: "Refresh the local repository observations before integration.", Argv: []string{}})
	r.TaskResult, r.HumanQA = tr, qa
	r.ResultOID, r.TechnicalGate, r.HumanQAOutcome = tr.ResultOID, tr.TechnicalGate, qa.Outcome
	replaceReadbackMember(&r.Value, "task_result", taskResultCanonical(*tr))
	replaceReadbackMember(&r.Value, "human_qa", taskHumanQACanonical(*qa))
	return r
}

func newPersistedIntegrationFallbackReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding, taskResult *TaskResultRecord, humanQA *TaskHumanQARecord, authority *IntegrationAuthority, intent *IntegrationIntent, attempt *IntegrationAttempt, result *IntegrationResult) WorkspaceTaskIntegrationReadback {
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
	readback.TaskResult = taskResult
	readback.HumanQA = humanQA
	readback.Classification = classification
	readback.GitChanged = changed
	readback.RecoveryStatus = recovery
	readback.NextAction = next
	readback.ResultOID = taskResult.ResultOID
	readback.TechnicalGate = taskResult.TechnicalGate
	readback.HumanQAOutcome = humanQA.Outcome
	if result != nil && result.BeforeObservation.OID != "" {
		readback.ParentOID = result.BeforeObservation.OID
	}
	replaceReadbackMember(&readback.Value, "task_result", taskResultCanonical(*taskResult))
	replaceReadbackMember(&readback.Value, "human_qa", taskHumanQACanonical(*humanQA))
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

func selectTaskIntegrationLeaf(registry WorkItemRegistry, taskID TaskID) (*TaskResultRecord, *TaskHumanQARecord, *IntegrationAuthority, bool) {
	authorities := []*IntegrationAuthority{}
	predecessors := map[IntegrationAuthorityID]bool{}
	for i := range registry.IntegrationAuthorities {
		a := &registry.IntegrationAuthorities[i]
		if a.TaskID != taskID {
			continue
		}
		authorities = append(authorities, a)
		if a.RetryAfterResultID != nil {
			if prior := findIntegrationResult(registry, *a.RetryAfterResultID); prior != nil {
				predecessors[prior.AuthorityID] = true
			}
		}
	}
	if len(authorities) > 0 {
		leaves := []*IntegrationAuthority{}
		for _, a := range authorities {
			if !predecessors[a.ID] {
				leaves = append(leaves, a)
			}
		}
		if len(leaves) != 1 {
			return nil, nil, nil, true
		}
		a := leaves[0]
		return findTaskResult(registry, a.TaskResultID), findHumanQA(registry, a.HumanQARecordID), a, false
	}
	green := []*TaskResultRecord{}
	for i := range registry.TaskResults {
		r := &registry.TaskResults[i]
		if r.TaskID == taskID && (r.TechnicalGate == "passed" || r.TechnicalGate == "good_enough_with_known_debt") {
			green = append(green, r)
		}
	}
	if len(green) == 0 {
		return nil, nil, nil, false
	}
	if len(green) != 1 {
		return nil, nil, nil, true
	}
	passes := []*TaskHumanQARecord{}
	for i := range registry.HumanQARecords {
		q := &registry.HumanQARecords[i]
		if q.TaskID == taskID && q.TaskResultID == green[0].ID && q.ResultOID == green[0].ResultOID && q.ResultTree == green[0].ResultTree && q.Outcome == "pass" {
			passes = append(passes, q)
		}
	}
	if len(passes) > 1 {
		return nil, nil, nil, true
	}
	if len(passes) == 0 {
		return green[0], nil, nil, false
	}
	return green[0], passes[0], nil, false
}

func newNullableIntegrationReadback(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, task TaskRecord, epic EpicRecord, parent EpicRepoBinding, classification string, next IntegrationNextAction) WorkspaceTaskIntegrationReadback {
	storeBytes, _ := encodeWorkItemRegistry(registry)
	persisted := canonicaljson.Object{{Name: "format_version", Value: int64(registry.FormatVersion)}, {Name: "sha256", Value: digestTaskBytes(storeBytes)}, {Name: "authority_id", Value: nil}, {Name: "intent_id", Value: nil}, {Name: "attempt_id", Value: nil}, {Name: "integration_result_id", Value: nil}, {Name: "freshness", Value: "fresh"}}
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
		{Name: "kind", Value: "WorkspaceTaskIntegrationReadback@1"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"},
		{Name: "workspace", Value: canonicaljson.Object{{Name: "root", Value: root}, {Name: "marker_sha256", Value: markerDigest}}}, {Name: "project", Value: structCanonical(project)}, {Name: "repository", Value: structCanonical(repository)}, {Name: "epic", Value: structCanonical(epicValue)}, {Name: "task", Value: structCanonical(taskValue)}, {Name: "source", Value: source},
		{Name: "task_result", Value: nil}, {Name: "human_qa", Value: nil}, {Name: "plan", Value: nil}, {Name: "authority", Value: nil}, {Name: "intent", Value: nil}, {Name: "attempt", Value: nil}, {Name: "integration_result", Value: nil},
		{Name: "persisted_before", Value: persisted}, {Name: "persisted_after", Value: persisted}, {Name: "observed_source", Value: observedWorktreeReadback(IntegrationWorktreeObservation{}, false)}, {Name: "observed_parent", Value: observedWorktreeReadback(IntegrationWorktreeObservation{}, false)}, {Name: "observed_inventory", Value: canonicaljson.Object{{Name: "entries", Value: []canonicaljson.Value{}}, {Name: "freshness", Value: "unknown"}}}, {Name: "observed_reflog", Value: canonicaljson.Object{{Name: "entries", Value: []canonicaljson.Value{}}, {Name: "freshness", Value: "unknown"}}},
		{Name: "classification", Value: classification}, {Name: "git_changed", Value: nil}, {Name: "recovery_status", Value: "blocked"}, {Name: "next_action", Value: structCanonical(next)},
	}
	return WorkspaceTaskIntegrationReadback{Value: value, Classification: classification, RecoveryStatus: "blocked", NextAction: next, TaskID: task.ID, EpicID: epic.ID, ParentRef: parent.Worktree.Ref, ParentOID: parent.Worktree.OID, TechnicalGate: "unknown", HumanQAOutcome: "not_recorded"}
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
		return IntegrationNextAction{Kind: "apply_confirmed_plan", Reason: "Apply the exact confirmed local fast-forward.", Argv: []string{"ply", "workspace", "task", "integrate", string(ctx.Task.ID), "--result", string(ctx.Result.ID), "--qa", string(ctx.QA.ID), "--expected-result-oid", ctx.Result.ResultOID, "--expected-parent-oid", ctx.Task.Worktree.ParentOID, "--apply", "--confirm", digest}}
	case "already_integrated":
		return IntegrationNextAction{Kind: "none", Reason: "The Task is already integrated.", Argv: []string{}}
	case "conflict":
		return IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: "Resolve the conflicting exact binding in a separate activity.", Argv: []string{}}
	default:
		return IntegrationNextAction{Kind: "resolve_blocker", Reason: "Resolve the named prerequisite before a new check.", Argv: []string{}}
	}
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
	return fmt.Sprintf("Status: %s\nTask result: %s / %s (%s; human QA %s)\nParent: %s / %s@%s\nGit changed: %s\nRecovery: %s\nNext action: %s\n", status, r.TaskID, r.ResultOID, r.TechnicalGate, r.HumanQAOutcome, r.EpicID, r.ParentRef, r.ParentOID, changed, r.RecoveryStatus, next)
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
			if t == reflect.TypeOf(WorkspaceTaskIntegrationPlan{}) && field.Name == "DeliveryAuthorization" && v.Field(i).IsNil() {
				continue
			}
			if t == reflect.TypeOf(DeliveryAgreement{}) && strings.Contains(field.Tag.Get("json"), ",omitempty") && v.Field(i).IsZero() {
				continue
			}
			if t == reflect.TypeOf(IntegrationAuthority{}) && field.Name == "DeliveryOwner" && v.Field(i).IsNil() {
				continue
			}
			if t == reflect.TypeOf(WorkspaceTaskIntegrationPlan{}) && field.Name == "TaskSpecGuard" && v.FieldByName("SchemaVersion").Int() == 1 {
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
