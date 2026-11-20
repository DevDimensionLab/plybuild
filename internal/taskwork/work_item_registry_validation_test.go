package taskwork

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These tests pin validateWorkItemRegistry's own rules: which error wins when
// several rules are broken at once, and its text and class. The registries are
// built in memory; no store or Git is involved.

const validationCommonDir = "/validation/main/.git"

const (
	creatingTask = iota
	nativeTask
	readyTask
	reconcilingTask
	unboundTask
)

const (
	creatingOperation = iota
	readyOperation
	reconcilingOperation
)

func validationWorktreeID(digit string) WorktreeID {
	return WorktreeID("wt_" + strings.Repeat(digit, 32))
}

func validationOperationID(digit string) WorktreeOperationID {
	return WorktreeOperationID("wop_" + strings.Repeat(digit, 32))
}

// validationEpic returns a valid Epic bound to repository "repo". The digit
// keeps its worktree ID apart from other records; path and ref follow the ID.
func validationEpic(id EpicID, digit string) EpicRecord {
	return EpicRecord{ID: id, Title: "Epic " + string(id), ProjectID: "project", RepoBindings: []EpicRepoBinding{{RepoID: "repo", GitCommonDir: validationCommonDir, Worktree: EpicWorktreeBinding{ID: validationWorktreeID(digit), OwnerKind: "epic", OwnerID: id, Origin: "adopted", Locator: "/validation/" + string(id), Ref: "refs/heads/" + string(id), OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), GitCommonDir: validationCommonDir}}}}
}

func validationTask(id TaskID, epic EpicRecord) TaskRecord {
	return TaskRecord{ID: id, Title: "Task " + string(id), Description: "Description", ParentEpicID: epic.ID, ProjectID: epic.ProjectID, RepoID: "repo", GitCommonDir: validationCommonDir, WorktreeState: WorkItemUnbound}
}

// validationOperation returns a valid operation for task in the given state,
// with the observation that state requires.
func validationOperation(task TaskRecord, epic EpicRecord, digit, state string) WorktreeOperationRecord {
	parent := epic.RepoBindings[0].Worktree
	operation := WorktreeOperationRecord{ID: validationOperationID(digit), TaskID: task.ID, WorktreeID: validationWorktreeID(digit), State: state, ProjectID: task.ProjectID, RepoID: task.RepoID, GitCommonDir: task.GitCommonDir, ParentEpicID: epic.ID, ParentWorktreeID: parent.ID, ParentRef: parent.Ref, ParentOID: parent.OID, ParentTree: parent.Tree, SourceRef: "refs/heads/" + string(task.ID), TargetLocator: "/validation/" + string(task.ID)}
	switch state {
	case "creating":
		operation.LastObservation = emptyCreateObservation("no_effect")
	case "reconciliation_required":
		operation.LastObservation = emptyCreateObservation("partial_or_unknown")
	case "ready":
		operation.LastObservation = WorktreeObservationRecord{Classification: "exact_effect", ParentRefOID: stringPointer(operation.ParentOID), ParentRefTree: stringPointer(operation.ParentTree), ParentWorktreeLocator: stringPointer(parent.Locator), ParentWorktreeRef: stringPointer(operation.ParentRef), ParentWorktreeOID: stringPointer(operation.ParentOID), ParentWorktreeTree: stringPointer(operation.ParentTree), ParentWorktreeGitCommonDir: stringPointer(operation.GitCommonDir), ParentWorktreeClean: boolPointer(true), SourceRefOID: stringPointer(operation.ParentOID), SourceRefTree: stringPointer(operation.ParentTree), SourceRefCheckedOutAt: stringPointer(operation.TargetLocator), TargetKind: "worktree", TargetLocator: stringPointer(operation.TargetLocator), TargetRef: stringPointer(operation.SourceRef), TargetOID: stringPointer(operation.ParentOID), TargetTree: stringPointer(operation.ParentTree), TargetGitCommonDir: stringPointer(operation.GitCommonDir), TargetClean: boolPointer(true), InventoryMatch: boolPointer(true), Reasons: []string{}}
	}
	resealOperation(&operation)
	return operation
}

// resealOperation recomputes the intent digest after a test changed an
// identity field, so that only the intended rule is broken.
func resealOperation(operation *WorktreeOperationRecord) {
	operation.IntentDigest, _ = intentDigest(*operation)
}

func validationReadyBinding(task TaskRecord, operation WorktreeOperationRecord) *TaskWorktreeBinding {
	return &TaskWorktreeBinding{ID: operation.WorktreeID, OwnerKind: "task", OwnerID: task.ID, Origin: "ply_created", Locator: operation.TargetLocator, Ref: operation.SourceRef, OID: operation.ParentOID, Tree: operation.ParentTree, GitCommonDir: operation.GitCommonDir, ParentEpicID: operation.ParentEpicID, ParentWorktreeID: operation.ParentWorktreeID, ParentRef: operation.ParentRef, ParentOID: operation.ParentOID, ParentTree: operation.ParentTree, OperationID: operation.ID, IntentDigest: operation.IntentDigest}
}

func validationNativeReceipt(task TaskRecord, epic EpicRecord, digit string) *TaskWorktreeBinding {
	parent := epic.RepoBindings[0].Worktree
	task.Worktree = &TaskWorktreeBinding{NativeReceiptVersion: nativeReceiptVersion, ID: validationWorktreeID(digit), OwnerKind: "task", OwnerID: task.ID, Origin: "ply_created", Locator: "/validation/" + string(task.ID), Ref: "refs/heads/" + string(task.ID), OID: parent.OID, Tree: parent.Tree, GitCommonDir: task.GitCommonDir, ParentEpicID: epic.ID, ParentWorktreeID: parent.ID, ParentRef: parent.Ref, ParentOID: parent.OID, ParentTree: parent.Tree, OperationID: validationOperationID(digit)}
	resealNativeReceipt(&task)
	return task.Worktree
}

func resealNativeReceipt(task *TaskRecord) {
	w := task.Worktree
	w.IntentDigest, _ = nativeIntentDigest(task.GitCommonDir, task.ParentEpicID, w.ParentOID, w.ParentRef, w.ParentTree, w.ParentWorktreeID, task.ProjectID, task.RepoID, w.Ref, w.Locator, task.ID)
}

// validationRegistry returns a valid registry with one Epic and one Task per
// worktree situation, in ID order: creating, native (ready by native receipt),
// ready (by operation), reconciling and unbound.
func validationRegistry() WorkItemRegistry {
	r := emptyWorkItemRegistry()
	epic := validationEpic("epic", "1")
	r.Epics = []EpicRecord{epic}
	r.EpicBaseVersions = []EpicBaseVersion{adoptedEpicBase(epic, epic.RepoBindings[0])}
	creating := validationTask("creating", epic)
	creating.WorktreeState = WorkItemCreating
	native := validationTask("native", epic)
	native.WorktreeState = WorkItemReady
	native.Worktree = validationNativeReceipt(native, epic, "c")
	ready := validationTask("ready", epic)
	ready.WorktreeState = WorkItemReady
	readyOperation := validationOperation(ready, epic, "4", "ready")
	ready.Worktree = validationReadyBinding(ready, readyOperation)
	reconciling := validationTask("reconciling", epic)
	reconciling.WorktreeState = WorkItemReconciliationRequired
	r.Tasks = []TaskRecord{creating, native, ready, reconciling, validationTask("unbound", epic)}
	r.WorktreeOperations = []WorktreeOperationRecord{validationOperation(creating, epic, "3", "creating"), readyOperation, validationOperation(reconciling, epic, "5", "reconciliation_required")}
	return r
}

// rebaseValidationEpics rebuilds the base versions after a test changed the
// Epics. The queue phase runs first and compares every base with its Epic
// binding; with matching bases the Epic rules themselves answer.
func rebaseValidationEpics(r *WorkItemRegistry) {
	r.EpicBaseVersions = []EpicBaseVersion{}
	seen := map[string]bool{}
	for _, epic := range r.Epics {
		for _, binding := range epic.RepoBindings {
			if key := string(epic.ID) + "/" + string(binding.RepoID); !seen[key] {
				seen[key] = true
				r.EpicBaseVersions = append(r.EpicBaseVersions, adoptedEpicBase(epic, binding))
			}
		}
	}
	sortQueueRegistry(r)
}

// addValidationEpic appends a second valid Epic "zeta".
func addValidationEpic(r *WorkItemRegistry) *EpicRecord {
	r.Epics = append(r.Epics, validationEpic("zeta", "2"))
	return &r.Epics[len(r.Epics)-1]
}

// addValidationBinding gives Epic "epic" a second valid repository binding.
func addValidationBinding(r *WorkItemRegistry) *EpicRepoBinding {
	epic := &r.Epics[0]
	binding := EpicRepoBinding{RepoID: "second", GitCommonDir: "/validation/second/.git", Worktree: EpicWorktreeBinding{ID: validationWorktreeID("2"), OwnerKind: "epic", OwnerID: epic.ID, Origin: "adopted", Locator: "/validation/second-epic", Ref: "refs/heads/epic", OID: strings.Repeat("a", 40), Tree: strings.Repeat("b", 40), GitCommonDir: "/validation/second/.git"}}
	epic.RepoBindings = append(epic.RepoBindings, binding)
	return &epic.RepoBindings[1]
}

// addOrphanOperation appends a valid operation whose Task is not registered.
func addOrphanOperation(r *WorkItemRegistry) *WorktreeOperationRecord {
	r.WorktreeOperations = append(r.WorktreeOperations, validationOperation(validationTask("ghost", r.Epics[0]), r.Epics[0], "9", "creating"))
	return &r.WorktreeOperations[len(r.WorktreeOperations)-1]
}

func invalidIdentifierError(kind, value string) string {
	return fmt.Sprintf("%s: %s ID %q must match %s", workspace.ErrorWorkInvalidArguments, kind, value, workspace.IdentifierRule)
}

func invalidTextError(kind string, maximum int) string {
	return fmt.Sprintf("%s: %s must be valid UTF-8, contain 1..%d characters, and have no surrounding whitespace", workspace.ErrorWorkInvalidArguments, kind, maximum)
}

// assertRegistryError requires the exact error text. Identifier and text rules
// answer with the invalid-arguments class; every other rule answers with an
// unclassified error.
func assertRegistryError(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil || err.Error() != want {
		t.Fatalf("validation = %v, want %q", err, want)
	}
	var classified *workspace.WorkItemError
	wantClassified := strings.HasPrefix(want, string(workspace.ErrorWorkInvalidArguments)+": ")
	if errors.As(err, &classified) != wantClassified || wantClassified && classified.Class != workspace.ErrorWorkInvalidArguments {
		t.Fatalf("validation error %#v has the wrong class", err)
	}
}

type registryViolation struct {
	name   string
	mutate func(*WorkItemRegistry)
	want   string
}

func assertRegistryViolations(t *testing.T, violations []registryViolation) {
	t.Helper()
	for _, violation := range violations {
		t.Run(violation.name, func(t *testing.T) {
			r := validationRegistry()
			violation.mutate(&r)
			rebaseValidationEpics(&r)
			assertRegistryError(t, validateWorkItemRegistry(r), violation.want)
		})
	}
}

// assertRegistryPrecedence checks that each violation is reported when it is
// alone, and that it wins over every violation listed after it when they all
// hold at once. Later violations are applied first, so an earlier one may
// build on the records they left.
func assertRegistryPrecedence(t *testing.T, violations []registryViolation) {
	t.Helper()
	assertRegistryViolations(t, violations)
	for i, violation := range violations {
		t.Run(violation.name+" wins over later rules", func(t *testing.T) {
			r := validationRegistry()
			for j := len(violations) - 1; j >= i; j-- {
				violations[j].mutate(&r)
			}
			rebaseValidationEpics(&r)
			assertRegistryError(t, validateWorkItemRegistry(r), violation.want)
		})
	}
}

func TestValidateWorkItemRegistryAcceptsEveryWorktreeSituation(t *testing.T) {
	t.Parallel()
	r := validationRegistry()
	if err := validateWorkItemRegistry(r); err != nil {
		t.Fatal(err)
	}
	// A retired Task keeps the binding of a ready one, by operation or receipt.
	r.Tasks[readyTask].WorktreeState, r.Tasks[nativeTask].WorktreeState = WorkItemRetired, WorkItemRetired
	if err := validateWorkItemRegistry(r); err != nil {
		t.Fatalf("retired: %v", err)
	}
}

func TestValidateWorkItemRegistryEmptyAndNilLists(t *testing.T) {
	t.Parallel()
	if err := validateWorkItemRegistry(emptyWorkItemRegistry()); err != nil {
		t.Fatalf("empty lists: %v", err)
	}
	// Epics, Tasks and operations need not be explicit arrays here; the other
	// phases own their own lists.
	r := emptyWorkItemRegistry()
	r.Epics, r.Tasks, r.WorktreeOperations = nil, nil, nil
	if err := validateWorkItemRegistry(r); err != nil {
		t.Fatalf("nil lists: %v", err)
	}
	assertRegistryError(t, validateWorkItemRegistry(WorkItemRegistry{}), "epic_create_operations must be an explicit array")
}

func TestValidateWorkItemRegistryZeroRecords(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		mutate func(*WorkItemRegistry)
		want   string
	}{
		"Epic":      {func(r *WorkItemRegistry) { r.Epics = []EpicRecord{{}} }, invalidIdentifierError("Epic", "")},
		"Task":      {func(r *WorkItemRegistry) { r.Tasks = []TaskRecord{{}} }, invalidIdentifierError("Task", "")},
		"operation": {func(r *WorkItemRegistry) { r.WorktreeOperations = []WorktreeOperationRecord{{}} }, "invalid or duplicate operation ID "},
	} {
		t.Run(name, func(t *testing.T) {
			r := emptyWorkItemRegistry()
			tc.mutate(&r)
			assertRegistryError(t, validateWorkItemRegistry(r), tc.want)
		})
	}
}

// The phases run in a fixed order, and the first failing phase answers.
func TestValidateWorkItemRegistryPhasePrecedence(t *testing.T) {
	t.Parallel()
	assertRegistryPrecedence(t, []registryViolation{
		{"Epic create operations", func(r *WorkItemRegistry) { r.EpicCreateOperations = nil }, "epic_create_operations must be an explicit array"},
		{"queue", func(r *WorkItemRegistry) { r.TaskQueueEvents = nil }, "queue arrays must be explicit"},
		{"Task content", func(r *WorkItemRegistry) { r.TaskProblemRevisions = nil }, "Task content arrays must be explicit"},
		{"array order", func(r *WorkItemRegistry) { slices.Reverse(r.Tasks) }, "work-item arrays must be sorted by unique IDs"},
		{"Epic", func(r *WorkItemRegistry) { r.Epics[0].Title = "" }, invalidTextError("Epic title", 128)},
		{"operation", func(r *WorkItemRegistry) { r.WorktreeOperations[creatingOperation].State = "bogus" }, "operation wop_33333333333333333333333333333333 has invalid state"},
		{"Task", func(r *WorkItemRegistry) { r.Tasks[unboundTask].Description = "" }, invalidTextError("Task description", 2048)},
		{"operation without Task", func(r *WorkItemRegistry) { addOrphanOperation(r) }, "operation refers to missing Task ghost"},
		{"Task lifecycle", func(r *WorkItemRegistry) {
			r.IntegrationAuthorities = []IntegrationAuthority{{ID: "same"}, {ID: "same"}}
		}, "Task lifecycle arrays must be sorted by unique IDs"},
	})
}

func TestValidateWorkItemRegistryRequiresSortedUniqueIDs(t *testing.T) {
	t.Parallel()
	const want = "work-item arrays must be sorted by unique IDs"
	assertRegistryViolations(t, []registryViolation{
		{"Epics unsorted", func(r *WorkItemRegistry) { addValidationEpic(r); slices.Reverse(r.Epics) }, want},
		{"Epic duplicate", func(r *WorkItemRegistry) { r.Epics = append(r.Epics, r.Epics[0]) }, want},
		{"Tasks unsorted", func(r *WorkItemRegistry) { slices.Reverse(r.Tasks) }, want},
		{"Task duplicate", func(r *WorkItemRegistry) { r.Tasks = append(r.Tasks, r.Tasks[unboundTask]) }, want},
		{"operations unsorted", func(r *WorkItemRegistry) { slices.Reverse(r.WorktreeOperations) }, want},
		{"operation duplicate", func(r *WorkItemRegistry) {
			r.WorktreeOperations = append(r.WorktreeOperations, r.WorktreeOperations[reconcilingOperation])
		}, want},
	})
}

func TestValidateWorkItemRegistryEpicRulePrecedence(t *testing.T) {
	t.Parallel()
	assertRegistryPrecedence(t, []registryViolation{
		{"ID", func(r *WorkItemRegistry) { r.Epics[0].ID = "Epic" }, invalidIdentifierError("Epic", "Epic")},
		{"title", func(r *WorkItemRegistry) { r.Epics[0].Title = " padded" }, invalidTextError("Epic title", 128)},
		{"Project", func(r *WorkItemRegistry) { r.Epics[0].ProjectID = "" }, invalidIdentifierError("Project", "")},
		{"binding order", func(r *WorkItemRegistry) { addValidationBinding(r); slices.Reverse(r.Epics[0].RepoBindings) }, "Epic epic repository bindings are not sorted"},
	})
	for name, bindings := range map[string][]EpicRepoBinding{"nil": nil, "empty": {}} {
		t.Run("no binding/"+name, func(t *testing.T) {
			r := validationRegistry()
			r.Epics[0].RepoBindings = bindings
			rebaseValidationEpics(&r)
			assertRegistryError(t, validateWorkItemRegistry(r), "Epic epic has no repository binding")
			r.Epics[0].ProjectID = "Project"
			assertRegistryError(t, validateWorkItemRegistry(r), invalidIdentifierError("Project", "Project"))
		})
	}
}

func TestValidateWorkItemRegistryEpicBindingRulePrecedence(t *testing.T) {
	t.Parallel()
	const sharedWorktree = "worktree ID wt_11111111111111111111111111111111 is shared by Epic epic and Epic epic"
	later := []registryViolation{
		{"common directory", func(r *WorkItemRegistry) { r.Epics[0].RepoBindings[1].GitCommonDir = "relative/.git" }, "Epic epic has invalid common directory"},
		{"worktree binding", func(r *WorkItemRegistry) { r.Epics[0].RepoBindings[1].Worktree.OwnerKind = "task" }, "Epic epic has invalid worktree binding"},
		{"worktree ID", func(r *WorkItemRegistry) {
			r.Epics[0].RepoBindings[1].Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID
		}, sharedWorktree},
	}
	// A duplicate repository has the ID of a binding that already passed the
	// identifier rule, so the two rules are compared with the later ones
	// separately.
	for _, first := range []registryViolation{
		{"duplicate repository", func(r *WorkItemRegistry) { r.Epics[0].RepoBindings[1].RepoID = "repo" }, "Epic epic has duplicate repository repo"},
		{"repository ID", func(r *WorkItemRegistry) { r.Epics[0].RepoBindings[1].RepoID = "second!" }, invalidIdentifierError("repository", "second!")},
	} {
		t.Run(first.name, func(t *testing.T) {
			violations := []registryViolation{}
			for _, violation := range append([]registryViolation{first}, later...) {
				mutate := violation.mutate
				violation.mutate = func(r *WorkItemRegistry) {
					if len(r.Epics[0].RepoBindings) == 1 {
						addValidationBinding(r)
					}
					mutate(r)
				}
				violations = append(violations, violation)
			}
			assertRegistryPrecedence(t, violations)
		})
	}
}

func TestValidateWorkItemRegistryEpicWorktreeBindingFields(t *testing.T) {
	t.Parallel()
	violations := []registryViolation{}
	for name, mutate := range map[string]func(*EpicRepoBinding){
		"ID":               func(b *EpicRepoBinding) { b.Worktree.ID = "wt_1" },
		"owner kind":       func(b *EpicRepoBinding) { b.Worktree.OwnerKind = "task" },
		"owner":            func(b *EpicRepoBinding) { b.Worktree.OwnerID = "other" },
		"origin":           func(b *EpicRepoBinding) { b.Worktree.Origin = "copied" },
		"common directory": func(b *EpicRepoBinding) { b.Worktree.GitCommonDir = "/validation/other/.git" },
		"locator":          func(b *EpicRepoBinding) { b.Worktree.Locator = "relative" },
		"ref":              func(b *EpicRepoBinding) { b.Worktree.Ref = "epic" },
	} {
		violations = append(violations, registryViolation{name, func(r *WorkItemRegistry) { mutate(&r.Epics[0].RepoBindings[0]) }, "Epic epic has invalid worktree binding"})
	}
	// An invalid OID or tree never reaches this rule: the queue phase rejects
	// the Epic's base version first.
	assertRegistryViolations(t, violations)
}

// Records are validated one at a time: the last rule of an earlier record (an
// Epic's worktree ID, an operation's ref, a ready Task's observation) answers
// before the first rule of the record after it, its ID.
func TestValidateWorkItemRegistryEarlierRecordWins(t *testing.T) {
	t.Parallel()
	assertRegistryViolations(t, []registryViolation{
		{"Epics", func(r *WorkItemRegistry) {
			addValidationEpic(r).ID = "zeta!"
			addValidationBinding(r).Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID
		}, "worktree ID wt_11111111111111111111111111111111 is shared by Epic epic and Epic epic"},
		{"operations", func(r *WorkItemRegistry) {
			r.WorktreeOperations[readyOperation].ID = "wop_4"
			r.WorktreeOperations[creatingOperation].SourceRef = r.Epics[0].RepoBindings[0].Worktree.Ref
			resealOperation(&r.WorktreeOperations[creatingOperation])
		}, "Git resource shared by Epic epic and operation wop_33333333333333333333333333333333"},
		{"Tasks", func(r *WorkItemRegistry) {
			r.Tasks[reconcilingTask].ID = "reconciling!"
			r.WorktreeOperations[readyOperation].LastObservation.InventoryMatch = boolPointer(false)
		}, "ready Task ready exact observation differs from operation"},
	})
}

func TestValidateWorkItemRegistryEpicsShareWorktreeIDsNeverButGitResourcesFreely(t *testing.T) {
	t.Parallel()
	r := validationRegistry()
	zeta := addValidationEpic(&r)
	zeta.RepoBindings[0].Worktree.Locator = r.Epics[0].RepoBindings[0].Worktree.Locator
	zeta.RepoBindings[0].Worktree.Ref = r.Epics[0].RepoBindings[0].Worktree.Ref
	rebaseValidationEpics(&r)
	if err := validateWorkItemRegistry(r); err != nil {
		t.Fatalf("Epics on one integration base: %v", err)
	}
	// The first Epic keeps the claim that later conflicts are reported against.
	operation := &r.WorktreeOperations[creatingOperation]
	operation.SourceRef = zeta.RepoBindings[0].Worktree.Ref
	resealOperation(operation)
	assertRegistryError(t, validateWorkItemRegistry(r), "Git resource shared by Epic epic and operation wop_33333333333333333333333333333333")

	r = validationRegistry()
	addValidationEpic(&r).RepoBindings[0].Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID
	rebaseValidationEpics(&r)
	assertRegistryError(t, validateWorkItemRegistry(r), "worktree ID wt_11111111111111111111111111111111 is shared by Epic epic and Epic zeta")
}

func TestValidateWorkItemRegistryOperationRulePrecedence(t *testing.T) {
	t.Parallel()
	operation := func(r *WorkItemRegistry) *WorktreeOperationRecord {
		return &r.WorktreeOperations[reconcilingOperation]
	}
	assertRegistryPrecedence(t, []registryViolation{
		{"ID", func(r *WorkItemRegistry) { operation(r).ID = "wop_x" }, "invalid or duplicate operation ID wop_x"},
		{"second operation of a Task", func(r *WorkItemRegistry) { operation(r).TaskID = "creating" }, "Task creating has multiple operations"},
		{"intent", func(r *WorkItemRegistry) { operation(r).ParentRef = "epic" }, "operation wop_55555555555555555555555555555555 has invalid intent"},
		{"worktree ID", func(r *WorkItemRegistry) { operation(r).WorktreeID = r.Epics[0].RepoBindings[0].Worktree.ID }, "worktree ID wt_11111111111111111111111111111111 is shared by Epic epic and operation wop_55555555555555555555555555555555"},
		{"intent digest", func(r *WorkItemRegistry) { operation(r).IntentDigest = "sha256:" + strings.Repeat("0", 64) }, "operation wop_55555555555555555555555555555555 intent digest does not match"},
		{"state", func(r *WorkItemRegistry) { operation(r).State = "bogus" }, "operation wop_55555555555555555555555555555555 has invalid state"},
		{"observation", func(r *WorkItemRegistry) { operation(r).LastObservation.Reasons = nil }, "operation wop_55555555555555555555555555555555: observation reasons must be an explicit sorted array"},
		{"Git resource", func(r *WorkItemRegistry) {
			operation(r).TargetLocator = r.Epics[0].RepoBindings[0].Worktree.Locator
			resealOperation(operation(r))
		}, "Git resource shared by Epic epic and operation wop_55555555555555555555555555555555"},
	})
}

func TestValidateWorkItemRegistryOperationIntentFields(t *testing.T) {
	t.Parallel()
	violations := []registryViolation{}
	for name, mutate := range map[string]func(*WorktreeOperationRecord){
		"worktree ID":      func(o *WorktreeOperationRecord) { o.WorktreeID = "wt_3" },
		"intent digest":    func(o *WorktreeOperationRecord) { o.IntentDigest = "sha256:3" },
		"common directory": func(o *WorktreeOperationRecord) { o.GitCommonDir = "relative/.git" },
		"target locator":   func(o *WorktreeOperationRecord) { o.TargetLocator = "relative" },
		"parent ref":       func(o *WorktreeOperationRecord) { o.ParentRef = "epic" },
		"source ref":       func(o *WorktreeOperationRecord) { o.SourceRef = "creating" },
		"parent OID":       func(o *WorktreeOperationRecord) { o.ParentOID = "abc" },
		"parent tree":      func(o *WorktreeOperationRecord) { o.ParentTree = "abc" },
	} {
		violations = append(violations, registryViolation{name, func(r *WorkItemRegistry) { mutate(&r.WorktreeOperations[creatingOperation]) }, "operation wop_33333333333333333333333333333333 has invalid intent"})
	}
	assertRegistryViolations(t, violations)
}

func TestValidateWorkItemRegistryOperationObservationErrorIsWrapped(t *testing.T) {
	t.Parallel()
	r := validationRegistry()
	r.WorktreeOperations[creatingOperation].LastObservation.Classification = "exact_effect"
	err := validateWorkItemRegistry(r)
	assertRegistryError(t, err, "operation wop_33333333333333333333333333333333: creating operation has invalid classification")
	if cause := errors.Unwrap(err); cause == nil || cause.Error() != "creating operation has invalid classification" {
		t.Fatalf("wrapped observation error = %v", cause)
	}
}

func TestValidateWorkItemRegistryOperationsShareNoWorktreeIDOrGitResource(t *testing.T) {
	t.Parallel()
	epicWorktree := func(r *WorkItemRegistry) EpicWorktreeBinding { return r.Epics[0].RepoBindings[0].Worktree }
	second := func(r *WorkItemRegistry) *WorktreeOperationRecord { return &r.WorktreeOperations[readyOperation] }
	const sharedWithEpic = "Git resource shared by Epic epic and operation wop_44444444444444444444444444444444"
	const sharedWithOperation = "Git resource shared by operation wop_33333333333333333333333333333333 and operation wop_44444444444444444444444444444444"
	assertRegistryViolations(t, []registryViolation{
		{"worktree ID of an operation", func(r *WorkItemRegistry) { second(r).WorktreeID = r.WorktreeOperations[creatingOperation].WorktreeID }, "worktree ID wt_33333333333333333333333333333333 is shared by operation wop_33333333333333333333333333333333 and operation wop_44444444444444444444444444444444"},
		{"Epic path", func(r *WorkItemRegistry) {
			second(r).TargetLocator = epicWorktree(r).Locator
			resealOperation(second(r))
		}, sharedWithEpic},
		{"Epic ref", func(r *WorkItemRegistry) { second(r).SourceRef = epicWorktree(r).Ref; resealOperation(second(r)) }, sharedWithEpic},
		{"operation path", func(r *WorkItemRegistry) {
			second(r).TargetLocator = r.WorktreeOperations[creatingOperation].TargetLocator
			resealOperation(second(r))
		}, sharedWithOperation},
		{"operation ref", func(r *WorkItemRegistry) {
			second(r).SourceRef = r.WorktreeOperations[creatingOperation].SourceRef
			resealOperation(second(r))
		}, sharedWithOperation},
		// With two conflicts, the holder of the path answers before the holder of the ref.
		{"Epic path and operation ref", func(r *WorkItemRegistry) {
			second(r).TargetLocator, second(r).SourceRef = epicWorktree(r).Locator, r.WorktreeOperations[creatingOperation].SourceRef
			resealOperation(second(r))
		}, sharedWithEpic},
		{"operation path and Epic ref", func(r *WorkItemRegistry) {
			second(r).TargetLocator, second(r).SourceRef = r.WorktreeOperations[creatingOperation].TargetLocator, epicWorktree(r).Ref
			resealOperation(second(r))
		}, sharedWithOperation},
	})
	// The same path and ref in another repository are other Git resources. The
	// operation then answers for its Task, whose repository it left.
	r := validationRegistry()
	operation := &r.WorktreeOperations[creatingOperation]
	operation.GitCommonDir, operation.TargetLocator, operation.SourceRef = "/validation/other/.git", epicWorktree(&r).Locator, epicWorktree(&r).Ref
	resealOperation(operation)
	assertRegistryError(t, validateWorkItemRegistry(r), "Task creating operation differs from its parent repository binding")
}

func TestValidateWorkItemRegistryTaskRulePrecedence(t *testing.T) {
	t.Parallel()
	assertRegistryPrecedence(t, []registryViolation{
		{"ID", func(r *WorkItemRegistry) { r.Tasks[creatingTask].ID = "Creating" }, invalidIdentifierError("Task", "Creating")},
		{"title", func(r *WorkItemRegistry) { r.Tasks[creatingTask].Title = "" }, invalidTextError("Task title", 128)},
		{"description", func(r *WorkItemRegistry) { r.Tasks[creatingTask].Description = "two\nlines" }, string(workspace.ErrorWorkInvalidArguments) + ": Task description must be one line without control characters"},
		{"unknown Epic", func(r *WorkItemRegistry) { r.Tasks[creatingTask].ParentEpicID = "missing" }, "Task creating has invalid Epic/Project reference"},
		{"unknown repository", func(r *WorkItemRegistry) { r.Tasks[creatingTask].RepoID = "missing" }, "Task creating has invalid repository binding"},
		{"operation and parent binding", func(r *WorkItemRegistry) {
			r.WorktreeOperations[creatingOperation].ParentRef = "refs/heads/elsewhere"
			resealOperation(&r.WorktreeOperations[creatingOperation])
		}, "Task creating operation differs from its parent repository binding"},
		{"worktree state", func(r *WorkItemRegistry) { r.Tasks[creatingTask].WorktreeState = "bogus" }, "Task creating has invalid worktree_state"},
	})
}

func TestValidateWorkItemRegistryOperationMatchesParentBindingFields(t *testing.T) {
	t.Parallel()
	violations := []registryViolation{}
	for name, mutate := range map[string]func(*WorktreeOperationRecord){
		"Project":          func(o *WorktreeOperationRecord) { o.ProjectID = "other" },
		"repository":       func(o *WorktreeOperationRecord) { o.RepoID = "other" },
		"common directory": func(o *WorktreeOperationRecord) { o.GitCommonDir = "/validation/other/.git" },
		"Epic":             func(o *WorktreeOperationRecord) { o.ParentEpicID = "other" },
		"parent worktree":  func(o *WorktreeOperationRecord) { o.ParentWorktreeID = validationWorktreeID("9") },
		"parent ref":       func(o *WorktreeOperationRecord) { o.ParentRef = "refs/heads/elsewhere" },
		"parent OID":       func(o *WorktreeOperationRecord) { o.ParentOID = strings.Repeat("e", 40) },
		"parent tree":      func(o *WorktreeOperationRecord) { o.ParentTree = strings.Repeat("e", 40) },
	} {
		violations = append(violations, registryViolation{name, func(r *WorkItemRegistry) {
			mutate(&r.WorktreeOperations[creatingOperation])
			resealOperation(&r.WorktreeOperations[creatingOperation])
		}, "Task creating operation differs from its parent repository binding"})
	}
	assertRegistryViolations(t, violations)
}

func TestValidateWorkItemRegistryWorktreeStateRequiresMatchingOperation(t *testing.T) {
	t.Parallel()
	dropOperation := func(r *WorkItemRegistry, index int) {
		r.WorktreeOperations = slices.Delete(r.WorktreeOperations, index, index+1)
	}
	// A binding without the native receipt marker follows the worktree state.
	bind := func(r *WorkItemRegistry, task int) {
		binding := *r.Tasks[readyTask].Worktree
		r.Tasks[task].Worktree = &binding
	}
	setState := func(r *WorkItemRegistry, index int, state string) {
		operation := &r.WorktreeOperations[index]
		*operation = validationOperation(r.Tasks[map[int]int{creatingOperation: creatingTask, readyOperation: readyTask, reconcilingOperation: reconcilingTask}[index]], r.Epics[0], string(operation.ID[len(operation.ID)-1]), state)
	}
	const unbound = "unbound Task unbound has worktree state"
	const creating = "Task creating has incomplete operation state"
	const reconciling = "Task reconciling has incomplete operation state"
	const ready = "ready Task ready has invalid binding"
	assertRegistryViolations(t, []registryViolation{
		{"unbound with binding", func(r *WorkItemRegistry) { bind(r, unboundTask) }, unbound},
		{"unbound with operation", func(r *WorkItemRegistry) { r.Tasks[creatingTask].WorktreeState = WorkItemUnbound }, "unbound Task creating has worktree state"},
		{"creating with binding", func(r *WorkItemRegistry) { bind(r, creatingTask) }, creating},
		{"creating without operation", func(r *WorkItemRegistry) { dropOperation(r, creatingOperation) }, creating},
		{"creating with reconciling operation", func(r *WorkItemRegistry) { setState(r, creatingOperation, "reconciliation_required") }, creating},
		{"reconciling with binding", func(r *WorkItemRegistry) { bind(r, reconcilingTask) }, reconciling},
		{"reconciling without operation", func(r *WorkItemRegistry) { dropOperation(r, reconcilingOperation) }, reconciling},
		{"reconciling with creating operation", func(r *WorkItemRegistry) { setState(r, reconcilingOperation, "creating") }, reconciling},
		{"ready without binding", func(r *WorkItemRegistry) { r.Tasks[readyTask].Worktree = nil }, ready},
		{"ready without operation", func(r *WorkItemRegistry) { dropOperation(r, readyOperation) }, ready},
		{"ready with creating operation", func(r *WorkItemRegistry) { setState(r, readyOperation, "creating") }, ready},
		{"retired without binding", func(r *WorkItemRegistry) {
			r.Tasks[readyTask].WorktreeState, r.Tasks[readyTask].Worktree = WorkItemRetired, nil
		}, ready},
		{"retired without operation", func(r *WorkItemRegistry) {
			r.Tasks[readyTask].WorktreeState = WorkItemRetired
			dropOperation(r, readyOperation)
		}, ready},
		{"empty state", func(r *WorkItemRegistry) { r.Tasks[unboundTask].WorktreeState = "" }, "Task unbound has invalid worktree_state"},
		{"unknown state", func(r *WorkItemRegistry) { r.Tasks[readyTask].WorktreeState = "bogus" }, "Task ready has invalid worktree_state"},
	})
}

func TestValidateWorkItemRegistryReadyTaskRulePrecedence(t *testing.T) {
	t.Parallel()
	assertRegistryPrecedence(t, []registryViolation{
		{"operation state", func(r *WorkItemRegistry) {
			operation := &r.WorktreeOperations[readyOperation]
			operation.State, operation.LastObservation = "reconciliation_required", emptyCreateObservation("partial_or_unknown")
		}, "ready Task ready has invalid binding"},
		{"binding", func(r *WorkItemRegistry) { r.Tasks[readyTask].Worktree.Ref = "refs/heads/elsewhere" }, "ready Task ready binding differs from operation"},
		{"observation", func(r *WorkItemRegistry) {
			r.WorktreeOperations[readyOperation].LastObservation.TargetClean = boolPointer(false)
		}, "ready Task ready exact observation differs from operation"},
	})
}

func TestValidateWorkItemRegistryReadyBindingFields(t *testing.T) {
	t.Parallel()
	violations := []registryViolation{}
	for name, mutate := range map[string]func(*TaskWorktreeBinding){
		"ID":               func(w *TaskWorktreeBinding) { w.ID = validationWorktreeID("9") },
		"operation":        func(w *TaskWorktreeBinding) { w.OperationID = validationOperationID("9") },
		"intent digest":    func(w *TaskWorktreeBinding) { w.IntentDigest = "sha256:" + strings.Repeat("0", 64) },
		"owner kind":       func(w *TaskWorktreeBinding) { w.OwnerKind = "epic" },
		"owner":            func(w *TaskWorktreeBinding) { w.OwnerID = "other" },
		"origin":           func(w *TaskWorktreeBinding) { w.Origin = "adopted" },
		"locator":          func(w *TaskWorktreeBinding) { w.Locator = "/validation/elsewhere" },
		"ref":              func(w *TaskWorktreeBinding) { w.Ref = "refs/heads/elsewhere" },
		"OID":              func(w *TaskWorktreeBinding) { w.OID = strings.Repeat("e", 40) },
		"tree":             func(w *TaskWorktreeBinding) { w.Tree = strings.Repeat("e", 40) },
		"common directory": func(w *TaskWorktreeBinding) { w.GitCommonDir = "/validation/other/.git" },
		"Epic":             func(w *TaskWorktreeBinding) { w.ParentEpicID = "other" },
		"parent worktree":  func(w *TaskWorktreeBinding) { w.ParentWorktreeID = validationWorktreeID("9") },
		"parent ref":       func(w *TaskWorktreeBinding) { w.ParentRef = "refs/heads/elsewhere" },
		"parent OID":       func(w *TaskWorktreeBinding) { w.ParentOID = strings.Repeat("e", 40) },
		"parent tree":      func(w *TaskWorktreeBinding) { w.ParentTree = strings.Repeat("e", 40) },
	} {
		violations = append(violations, registryViolation{name, func(r *WorkItemRegistry) { mutate(r.Tasks[readyTask].Worktree) }, "ready Task ready binding differs from operation"})
	}
	assertRegistryViolations(t, violations)
}

func TestValidateWorkItemRegistryReadyObservationFields(t *testing.T) {
	t.Parallel()
	otherOID, otherPath, otherRef := stringPointer(strings.Repeat("e", 40)), stringPointer("/validation/elsewhere"), stringPointer("refs/heads/elsewhere")
	violations := []registryViolation{}
	for name, mutate := range map[string]func(*WorktreeObservationRecord){
		"parent ref OID":                   func(o *WorktreeObservationRecord) { o.ParentRefOID = otherOID },
		"parent ref tree":                  func(o *WorktreeObservationRecord) { o.ParentRefTree = otherOID },
		"parent worktree locator":          func(o *WorktreeObservationRecord) { o.ParentWorktreeLocator = otherPath },
		"parent worktree ref":              func(o *WorktreeObservationRecord) { o.ParentWorktreeRef = otherRef },
		"parent worktree OID":              func(o *WorktreeObservationRecord) { o.ParentWorktreeOID = otherOID },
		"parent worktree tree":             func(o *WorktreeObservationRecord) { o.ParentWorktreeTree = otherOID },
		"parent worktree common directory": func(o *WorktreeObservationRecord) { o.ParentWorktreeGitCommonDir = otherPath },
		"parent worktree dirty":            func(o *WorktreeObservationRecord) { o.ParentWorktreeClean = boolPointer(false) },
		"source ref OID":                   func(o *WorktreeObservationRecord) { o.SourceRefOID = otherOID },
		"source ref tree":                  func(o *WorktreeObservationRecord) { o.SourceRefTree = otherOID },
		"source ref checked out elsewhere": func(o *WorktreeObservationRecord) { o.SourceRefCheckedOutAt = otherPath },
		"target locator":                   func(o *WorktreeObservationRecord) { o.TargetLocator = otherPath },
		"target ref":                       func(o *WorktreeObservationRecord) { o.TargetRef = otherRef },
		"target OID":                       func(o *WorktreeObservationRecord) { o.TargetOID = otherOID },
		"target tree":                      func(o *WorktreeObservationRecord) { o.TargetTree = otherOID },
		"target common directory":          func(o *WorktreeObservationRecord) { o.TargetGitCommonDir = otherPath },
		"target dirty":                     func(o *WorktreeObservationRecord) { o.TargetClean = boolPointer(false) },
		"inventory":                        func(o *WorktreeObservationRecord) { o.InventoryMatch = boolPointer(false) },
	} {
		violations = append(violations, registryViolation{name, func(r *WorkItemRegistry) { mutate(&r.WorktreeOperations[readyOperation].LastObservation) }, "ready Task ready exact observation differs from operation"})
	}
	assertRegistryViolations(t, violations)
}

func TestValidateWorkItemRegistryNativeReceiptRulePrecedence(t *testing.T) {
	t.Parallel()
	native := func(r *WorkItemRegistry) *TaskRecord { return &r.Tasks[nativeTask] }
	// An operation for the native Task: wop_2… sorts first.
	addOperation := func(r *WorkItemRegistry) *WorktreeOperationRecord {
		operation := validationOperation(*native(r), r.Epics[0], "2", "ready")
		operation.SourceRef, operation.TargetLocator = "refs/heads/second", "/validation/second"
		resealOperation(&operation)
		r.WorktreeOperations = append([]WorktreeOperationRecord{operation}, r.WorktreeOperations...)
		return &r.WorktreeOperations[0]
	}
	assertRegistryPrecedence(t, []registryViolation{
		{"operation and parent binding", func(r *WorkItemRegistry) {
			operation := addOperation(r)
			operation.ParentRef = "refs/heads/elsewhere"
			resealOperation(operation)
		}, "Task native operation differs from its parent repository binding"},
		{"receipt", func(r *WorkItemRegistry) { native(r).Worktree.Origin = "adopted" }, "Task native has invalid native receipt"},
		{"receipt digest", func(r *WorkItemRegistry) { native(r).Worktree.IntentDigest = "sha256:" + strings.Repeat("0", 64) }, "Task native native receipt digest differs"},
		{"worktree ID", func(r *WorkItemRegistry) { native(r).Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID }, "native receipt reuses worktree or operation ID"},
		{"Git resource", func(r *WorkItemRegistry) {
			native(r).Worktree.Ref = r.Epics[0].RepoBindings[0].Worktree.Ref
			resealNativeReceipt(native(r))
		}, "native receipt shares Git resource with Epic epic"},
	})
	assertRegistryViolations(t, []registryViolation{
		{"matching operation", func(r *WorkItemRegistry) { addOperation(r) }, "Task native has invalid native receipt"},
		// The receipt marker decides the rules, before the worktree state does.
		{"unbound state", func(r *WorkItemRegistry) { native(r).WorktreeState = WorkItemUnbound }, "Task native has invalid native receipt"},
		{"unknown state", func(r *WorkItemRegistry) { native(r).WorktreeState = "bogus" }, "Task native has invalid native receipt"},
		{"unknown receipt version", func(r *WorkItemRegistry) { native(r).Worktree.NativeReceiptVersion = -1 }, "Task native has invalid native receipt"},
		{"no receipt marker", func(r *WorkItemRegistry) { native(r).Worktree.NativeReceiptVersion = 0 }, "ready Task native has invalid binding"},
	})
}

func TestValidateWorkItemRegistryNativeReceiptSharesNoIDOrGitResource(t *testing.T) {
	t.Parallel()
	native := func(r *WorkItemRegistry) *TaskRecord { return &r.Tasks[nativeTask] }
	// A second native Task that sorts last and so answers for the conflict.
	addNative := func(r *WorkItemRegistry) *TaskRecord {
		task := validationTask("zulu", r.Epics[0])
		task.WorktreeState = WorkItemRetired
		task.Worktree = validationNativeReceipt(task, r.Epics[0], "d")
		r.Tasks = append(r.Tasks, task)
		return &r.Tasks[len(r.Tasks)-1]
	}
	const reused = "native receipt reuses worktree or operation ID"
	assertRegistryViolations(t, []registryViolation{
		{"Epic worktree ID", func(r *WorkItemRegistry) { native(r).Worktree.ID = r.Epics[0].RepoBindings[0].Worktree.ID }, reused},
		{"operation worktree ID", func(r *WorkItemRegistry) { native(r).Worktree.ID = r.WorktreeOperations[creatingOperation].WorktreeID }, reused},
		{"operation ID", func(r *WorkItemRegistry) { native(r).Worktree.OperationID = r.WorktreeOperations[creatingOperation].ID }, reused},
		{"native worktree ID", func(r *WorkItemRegistry) { addNative(r).Worktree.ID = native(r).Worktree.ID }, reused},
		{"native operation ID", func(r *WorkItemRegistry) { addNative(r).Worktree.OperationID = native(r).Worktree.OperationID }, reused},
		{"Epic path", func(r *WorkItemRegistry) {
			native(r).Worktree.Locator = r.Epics[0].RepoBindings[0].Worktree.Locator
			resealNativeReceipt(native(r))
		}, "native receipt shares Git resource with Epic epic"},
		{"operation ref", func(r *WorkItemRegistry) {
			native(r).Worktree.Ref = r.WorktreeOperations[reconcilingOperation].SourceRef
			resealNativeReceipt(native(r))
		}, "native receipt shares Git resource with operation wop_55555555555555555555555555555555"},
		{"native path", func(r *WorkItemRegistry) {
			second := addNative(r)
			second.Worktree.Locator = native(r).Worktree.Locator
			resealNativeReceipt(second)
		}, "native receipt shares Git resource with Task native"},
		{"native ref", func(r *WorkItemRegistry) {
			second := addNative(r)
			second.Worktree.Ref = native(r).Worktree.Ref
			resealNativeReceipt(second)
		}, "native receipt shares Git resource with Task native"},
		// With two conflicts, the holder of the path answers before the holder of the ref.
		{"Epic path and operation ref", func(r *WorkItemRegistry) {
			native(r).Worktree.Locator, native(r).Worktree.Ref = r.Epics[0].RepoBindings[0].Worktree.Locator, r.WorktreeOperations[reconcilingOperation].SourceRef
			resealNativeReceipt(native(r))
		}, "native receipt shares Git resource with Epic epic"},
		{"operation path and Epic ref", func(r *WorkItemRegistry) {
			native(r).Worktree.Locator, native(r).Worktree.Ref = r.WorktreeOperations[reconcilingOperation].TargetLocator, r.Epics[0].RepoBindings[0].Worktree.Ref
			resealNativeReceipt(native(r))
		}, "native receipt shares Git resource with operation wop_55555555555555555555555555555555"},
	})
	r := validationRegistry()
	addNative(&r)
	if err := validateWorkItemRegistry(r); err != nil {
		t.Fatalf("two native receipts: %v", err)
	}
}

func TestValidateWorkItemRegistryOperationNeedsRegisteredTask(t *testing.T) {
	t.Parallel()
	assertRegistryViolations(t, []registryViolation{
		{"unknown Task", func(r *WorkItemRegistry) { addOrphanOperation(r) }, "operation refers to missing Task ghost"},
		{"removed Task", func(r *WorkItemRegistry) { r.Tasks = slices.Delete(r.Tasks, creatingTask, creatingTask+1) }, "operation refers to missing Task creating"},
		{"only operations", func(r *WorkItemRegistry) {
			r.Tasks = nil
			r.WorktreeOperations = r.WorktreeOperations[:1]
		}, "operation refers to missing Task creating"},
	})
}

// The three tests below pin rules that an earlier phase or an earlier record
// currently keeps from answering through validateWorkItemRegistry. They call
// the helpers, so that the rules are still there when the earlier phase changes.

// The queue phase rejects an Epic's invalid OID or tree with its base version,
// and the Epic create phase decides whether a created Epic is allowed.
func TestValidateWorkItemRegistryShadowedEpicWorktreeRules(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		mutate func(*EpicWorktreeBinding)
		valid  bool
	}{
		"adopted":      {func(w *EpicWorktreeBinding) {}, true},
		"created":      {func(w *EpicWorktreeBinding) { w.Origin = "ply_created" }, true},
		"invalid OID":  {func(w *EpicWorktreeBinding) { w.OID = "abc" }, false},
		"invalid tree": {func(w *EpicWorktreeBinding) { w.Tree = "abc" }, false},
	} {
		epic := validationEpic("epic", "1")
		tc.mutate(&epic.RepoBindings[0].Worktree)
		if validEpicWorktree(epic, epic.RepoBindings[0]) != tc.valid {
			t.Fatalf("%s Epic worktree: valid = %v", name, !tc.valid)
		}
	}
}

// The queue phase keeps every base version on its Epic binding's worktree ID
// and ref, so an operation on a registered base also matches the binding.
func TestValidateWorkItemRegistryShadowedOperationParentRule(t *testing.T) {
	t.Parallel()
	r := validationRegistry()
	validation := &workItemRegistryValidation{registry: r}
	parent := r.Epics[0].RepoBindings[0].Worktree
	for name, tc := range map[string]struct {
		mutate  func(*EpicWorktreeBinding)
		matches bool
	}{
		"same binding":      {func(w *EpicWorktreeBinding) {}, true},
		"other worktree ID": {func(w *EpicWorktreeBinding) { w.ID = validationWorktreeID("9") }, false},
		"other ref":         {func(w *EpicWorktreeBinding) { w.Ref = "refs/heads/elsewhere" }, false},
	} {
		binding := parent
		tc.mutate(&binding)
		if validation.operationMatchesParent(r.Tasks[creatingTask], r.WorktreeOperations[creatingOperation], binding) != tc.matches {
			t.Fatalf("%s: matches = %v", name, !tc.matches)
		}
	}
}

// No other record can hold the worktree ID of a ready Task's operation: the
// operation claimed it, and a native receipt may not reuse a claimed ID.
func TestValidateWorkItemRegistryShadowedReadyWorktreeIDRule(t *testing.T) {
	t.Parallel()
	r := validationRegistry()
	ready, operation, parent := r.Tasks[readyTask], r.WorktreeOperations[readyOperation], r.Epics[0].RepoBindings[0].Worktree
	validation := &workItemRegistryValidation{registry: r, worktreeOwners: map[WorktreeID]string{ready.Worktree.ID: "operation " + string(operation.ID)}}
	if err := validation.validateReadyTask(ready, operation, true, parent); err != nil {
		t.Fatalf("worktree ID held by the Task's own operation: %v", err)
	}
	validation.worktreeOwners[ready.Worktree.ID] = "Task other"
	assertRegistryError(t, validation.validateReadyTask(ready, operation, true, parent), "worktree ID wt_44444444444444444444444444444444 reused")
}
