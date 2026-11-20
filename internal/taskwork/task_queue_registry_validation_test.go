package taskwork

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These tests pin validateQueueRegistry's own rules: which error wins when
// several rules are broken at once, and its text and class. The registries are
// built in memory on top of validationRegistry; no store or Git is involved.

const queueValidationTime = "2026-11-20T10:00:00Z"

// Base versions of queueValidationRegistry, in registry order.
const (
	epicBaseV1 = iota
	epicBaseV2
	zetaBaseV1
)

// Queue events of queueValidationRegistry, in registry order.
const (
	firstSetEvent = iota
	preparationEvent
	goalSetEvent
	otherQueueEvent
)

func queueValidationDigest(digit string) string { return "sha256:" + strings.Repeat(digit, 64) }

func queueValidationOID(digit string) string { return strings.Repeat(digit, 40) }

func queueValidationID(digit string) string { return "que_" + strings.Repeat(digit, 64) }

func intPointer(value int) *int { return &value }

func queueValidationTarget(epic EpicRecord, binding EpicRepoBinding) QueueTarget {
	wt := binding.Worktree
	return QueueTarget{ProjectID: epic.ProjectID, RepoID: binding.RepoID, EpicID: epic.ID, GitCommonDir: binding.GitCommonDir, ParentWorktreeID: wt.ID, ParentLocator: wt.Locator, ParentRef: wt.Ref}
}

// queueValidationPlan returns a valid plan of the current schema that moves
// the base of the Epic's first binding from previous to a new OID and tree.
func queueValidationPlan(epic EpicRecord, previous EpicBaseVersion, oidDigit, treeDigit string) WorkspaceEpicBaseUpdatePlan {
	plan := WorkspaceEpicBaseUpdatePlan{Kind: "WorkspaceEpicBaseUpdatePlan@3", SchemaVersion: currentEpicBaseUpdatePlanSchema, Workspace: IntegrationPlanWorkspace{Root: "/validation", MarkerSHA256: queueValidationDigest("e")}, RegistryFormatVersion: RegistryFormat, ExpectedRevision: previous.Revision, PreviousOID: previous.OID, PreviousTree: previous.Tree, NextOID: queueValidationOID(oidDigit), NextTree: queueValidationOID(treeDigit), AffectedSelections: []EpicAffectedSelection{}, UnresolvedOperations: []string{}}
	plan.ObservedParent = PlanWorktreeObservation{OID: plan.NextOID, Tree: plan.NextTree, ObjectFormat: "sha1", RefFormat: "files", Symbolic: true, Clean: true, StatusEntries: []IntegrationStatusEntry{}, InProgress: []string{}}
	retargetQueuePlan(&plan, epic, epic.RepoBindings[0])
	return plan
}

// retargetQueuePlan points the plan and its observation at binding of epic.
func retargetQueuePlan(plan *WorkspaceEpicBaseUpdatePlan, epic EpicRecord, binding EpicRepoBinding) {
	plan.Target = queueValidationTarget(epic, binding)
	plan.ObservedParent.WorktreeLocator, plan.ObservedParent.Ref, plan.ObservedParent.GitCommonDir = plan.Target.ParentLocator, plan.Target.ParentRef, plan.Target.GitCommonDir
}

// resealBaseUpdate gives the update the digest and ID of its plan.
func resealBaseUpdate(update *EpicBaseUpdate) {
	update.PlanSHA256 = queueDigest(update.Plan)
	update.ID = "ebu_" + strings.TrimPrefix(update.PlanSHA256, "sha256:")
}

func sealedBaseUpdate(plan WorkspaceEpicBaseUpdatePlan, phase string, outcome *EpicBaseUpdateOutcome) EpicBaseUpdate {
	update := EpicBaseUpdate{Plan: plan, Phase: phase, RecordedAtUTC: queueValidationTime, Outcome: outcome}
	resealBaseUpdate(&update)
	return update
}

func committedOutcome(revision int) *EpicBaseUpdateOutcome {
	return &EpicBaseUpdateOutcome{Revision: intPointer(revision), Reasons: []QueueReason{}}
}

// nextBaseVersion returns the base version that the committed update produced
// after previous.
func nextBaseVersion(previous EpicBaseVersion, update EpicBaseUpdate) EpicBaseVersion {
	next := previous
	next.Revision, next.PreviousRevision, next.UpdateID = previous.Revision+1, intPointer(previous.Revision), stringPointer(update.ID)
	next.OID, next.Tree = update.Plan.NextOID, update.Plan.NextTree
	return next
}

func sealedQueueEvent(queue string, revision int, kind string, request map[string]any) TaskQueueEvent {
	return TaskQueueEvent{QueueID: queue, Revision: revision, Kind: kind, RequestSHA256: queueDigest(request), Request: request, RecordedAtUTC: queueValidationTime}
}

func queueDraftFields(kind string, schema int, key string, expected int, entries []any) map[string]any {
	return map[string]any{"kind": kind, "schema_version": schema, "publication_key": key, "project_id": "project", "repo_id": "repo", "epic_id": "epic", "expected_revision": expected, "entries": append([]any{}, entries...), "registry_upgrade": nil}
}

// legacyQueueDraft returns a valid set request of schema 1.
func legacyQueueDraft(key string, expected int, entries ...any) map[string]any {
	draft := queueDraftFields("WorkspaceTaskQueueDraft@1", 1, key, expected, entries)
	draft["human_decision"] = map[string]any{"actor_claim": "human", "decided_at_utc": queueValidationTime, "source": "human_cli", "statement": "Queue it"}
	return draft
}

// goalQueueDraft returns a valid set request of schema 2.
func goalQueueDraft(key string, expected int, entries ...any) map[string]any {
	draft := queueDraftFields("WorkspaceTaskQueueDraft@2", 2, key, expected, entries)
	draft["recorder"] = map[string]any{"actor_claim": "agent", "control_surface": "cli", "recorded_at_utc": queueValidationTime}
	return draft
}

func selectionQueueEntry(task string, selection *TaskSelectionReference) map[string]any {
	if selection == nil {
		return map[string]any{"task_id": task, "selection": nil}
	}
	return map[string]any{"task_id": task, "selection": map[string]any{"id": selection.ID, "manifest_sha256": selection.ManifestSHA256}}
}

func goalQueueEntry(task string, spec TaskSpecReference) map[string]any {
	return map[string]any{"task_id": task, "goal": map[string]any{"spec_id": spec.SpecID, "spec": map[string]any{"revision": spec.Revision, "manifest_sha256": spec.ManifestSHA256}}}
}

// queueValidationRegistry extends validationRegistry with every kind of queue
// record, all valid. Epic "epic" has a committed base update that produced its
// base v2 and names two affected Tasks, a historical conflict against v1 and a
// later intent against v2; Epic "zeta" has only its adopted base. The first
// queue is set twice around a removed preparation event, once per draft
// schema; the second queue holds a single preparation event.
func queueValidationRegistry() WorkItemRegistry {
	r := validationRegistry()
	addValidationEpic(&r)
	rebaseValidationEpics(&r)
	epic := r.Epics[0]
	nativeSelection := TaskSelectionReference{TaskID: "native", Ordinal: 1, ID: "sel_" + strings.Repeat("1", 32), ManifestSHA256: queueValidationDigest("1")}
	unboundSelection := TaskSelectionReference{TaskID: "unbound", Ordinal: 1, ID: "sel_" + strings.Repeat("2", 32), ManifestSHA256: queueValidationDigest("2")}
	spec := TaskSpecReference{TaskID: "unbound", SpecID: "goal", Revision: 1, ManifestSHA256: queueValidationDigest("3")}
	r.TaskSolutionSelections = []TaskSelectionReference{nativeSelection, unboundSelection}
	r.TaskSpecRevisions = []TaskSpecReference{spec}

	v1 := r.EpicBaseVersions[epicBaseV1]
	committedPlan := queueValidationPlan(epic, v1, "c", "d")
	committedPlan.AffectedSelections = []EpicAffectedSelection{{TaskID: "creating"}, {TaskID: "native", Selection: &TaskDecisionRef{ID: nativeSelection.ID, ManifestSHA256: nativeSelection.ManifestSHA256}}}
	committed := sealedBaseUpdate(committedPlan, "committed", committedOutcome(2))
	v2 := nextBaseVersion(v1, committed)
	conflictPlan := queueValidationPlan(epic, v1, "e", "f")
	conflictPlan.Kind, conflictPlan.SchemaVersion, conflictPlan.CurrentPreparationID = "WorkspaceEpicBaseUpdatePlan@1", 1, stringPointer("prep_historical")
	conflict := sealedBaseUpdate(conflictPlan, "conflict", &EpicBaseUpdateOutcome{Reasons: []QueueReason{{Code: "epic_base_changed", Message: "The base changed"}}})
	intent := sealedBaseUpdate(queueValidationPlan(epic, v2, "1", "2"), "intent", nil)
	r.EpicBaseVersions = append(r.EpicBaseVersions, v2)
	r.EpicBaseUpdates = []EpicBaseUpdate{committed, conflict, intent}

	first, second := queueValidationID("a"), queueValidationID("b")
	r.TaskQueueEvents = []TaskQueueEvent{
		sealedQueueEvent(first, 1, "set", legacyQueueDraft("queue/first", 0, selectionQueueEntry("unbound", &unboundSelection))),
		sealedQueueEvent(first, 2, "reserve", map[string]any{"expected_revision": 1}),
		sealedQueueEvent(first, 3, "set", goalQueueDraft("queue/second", 2, selectionQueueEntry("reconciling", nil), goalQueueEntry("unbound", spec))),
		sealedQueueEvent(second, 1, "release", map[string]any{"expected_revision": 0}),
	}
	sortQueueRegistry(&r)
	return r
}

// queueUpdate returns the first base update in the given phase.
func queueUpdate(r *WorkItemRegistry, phase string) *EpicBaseUpdate {
	for i := range r.EpicBaseUpdates {
		if r.EpicBaseUpdates[i].Phase == phase {
			return &r.EpicBaseUpdates[i]
		}
	}
	panic("no base update in phase " + phase)
}

// changeBasePlan changes the plan of the update in the given phase and reseals
// it: the update gets the ID of its new plan, the base version it produced
// follows, and the updates are sorted again. Only the changed rule is broken.
func changeBasePlan(r *WorkItemRegistry, phase string, change func(*WorkspaceEpicBaseUpdatePlan)) {
	update := queueUpdate(r, phase)
	old := update.ID
	change(&update.Plan)
	resealBaseUpdate(update)
	for i := range r.EpicBaseVersions {
		if version := &r.EpicBaseVersions[i]; version.UpdateID != nil && *version.UpdateID == old {
			version.UpdateID = stringPointer(update.ID)
		}
	}
	sortQueueRegistry(r)
}

func basePlanViolation(phase string, change func(*WorkspaceEpicBaseUpdatePlan)) func(*WorkItemRegistry) {
	return func(r *WorkItemRegistry) { changeBasePlan(r, phase, change) }
}

// changeQueueRequest changes the request of an event and reseals its digest.
func changeQueueRequest(r *WorkItemRegistry, event int, change func(request map[string]any)) {
	e := &r.TaskQueueEvents[event]
	change(e.Request)
	e.RequestSHA256 = queueDigest(e.Request)
}

func queueRequestViolation(event int, change func(request map[string]any)) func(*WorkItemRegistry) {
	return func(r *WorkItemRegistry) { changeQueueRequest(r, event, change) }
}

func queueRequestEntry(request map[string]any, index int) map[string]any {
	return request["entries"].([]any)[index].(map[string]any)
}

// addQueueValidationBinding gives Epic "epic" a second repository "second"
// with its adopted base, which sorts between the bases of "epic" and "zeta".
func addQueueValidationBinding(r *WorkItemRegistry) EpicRepoBinding {
	binding := *addValidationBinding(r)
	r.EpicBaseVersions = append(r.EpicBaseVersions, adoptedEpicBase(r.Epics[0], binding))
	sortQueueRegistry(r)
	return binding
}

// addSecondBaseUpdate appends a committed update of Epic "epic" from base v2,
// with the base v3 that it produced.
func addSecondBaseUpdate(r *WorkItemRegistry, change func(*WorkspaceEpicBaseUpdatePlan)) {
	v2 := r.EpicBaseVersions[epicBaseV2]
	plan := queueValidationPlan(r.Epics[0], v2, "3", "4")
	change(&plan)
	update := sealedBaseUpdate(plan, "committed", committedOutcome(3))
	r.EpicBaseUpdates = append(r.EpicBaseUpdates, update)
	r.EpicBaseVersions = append(r.EpicBaseVersions, nextBaseVersion(v2, update))
	sortQueueRegistry(r)
}

// addUnappliedUpdate appends a valid committed update that no base version
// records. Its workspace root is chosen so that its ID sorts before or after
// the intent's.
func addUnappliedUpdate(r *WorkItemRegistry, beforeIntent bool) {
	update := *queueUpdate(r, "committed")
	intent := queueUpdate(r, "intent").ID
	for i := 0; ; i++ {
		update.Plan.Workspace.Root = fmt.Sprintf("/validation/%d", i)
		resealBaseUpdate(&update)
		if (update.ID < intent) == beforeIntent {
			break
		}
	}
	r.EpicBaseUpdates = append(r.EpicBaseUpdates, update)
	sortQueueRegistry(r)
}

// queueErrorClass names the class that a caller can read from the error; the
// plain errors of the queue rules have none.
func queueErrorClass(err error) string {
	var content *TaskContentError
	var work *workspace.WorkItemError
	var unknown *contentUnknownFields
	switch {
	case errors.As(err, &content):
		return "content " + content.Code
	case errors.As(err, &work):
		return "work " + string(work.Class)
	case errors.As(err, &unknown):
		return "unknown fields " + strings.Join(unknown.names, ",")
	}
	return ""
}

// assertQueueError requires the exact error text and class; without a class
// the error must be unclassified.
func assertQueueError(t *testing.T, err error, want string, class ...string) {
	t.Helper()
	if err == nil || err.Error() != want {
		t.Fatalf("validation = %v, want %q", err, want)
	}
	if got := queueErrorClass(err); got != strings.Join(class, "") {
		t.Fatalf("validation error %#v has class %q, want %q", err, got, strings.Join(class, ""))
	}
}

func assertQueueViolations(t *testing.T, violations []registryViolation) {
	t.Helper()
	for _, violation := range violations {
		t.Run(violation.name, func(t *testing.T) {
			r := queueValidationRegistry()
			violation.mutate(&r)
			assertQueueError(t, validateQueueRegistry(r), violation.want)
		})
	}
}

// assertQueuePrecedence checks that each violation is reported when it is
// alone, and that it wins over every violation listed after it when they all
// hold at once. Later violations are applied first, so an earlier one may
// build on the records they left.
func assertQueuePrecedence(t *testing.T, violations []registryViolation) {
	t.Helper()
	assertQueueViolations(t, violations)
	for i, violation := range violations {
		t.Run(violation.name+" wins over later rules", func(t *testing.T) {
			r := queueValidationRegistry()
			for j := len(violations) - 1; j >= i; j-- {
				violations[j].mutate(&r)
			}
			assertQueueError(t, validateQueueRegistry(r), violation.want)
		})
	}
}

// assertQueueAccepts checks leniencies: each change leaves the registry valid.
func assertQueueAccepts(t *testing.T, changes map[string]func(*WorkItemRegistry)) {
	t.Helper()
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			r := queueValidationRegistry()
			change(&r)
			if err := validateQueueRegistry(r); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestValidateQueueRegistryAcceptsEveryRecordKind(t *testing.T) {
	t.Parallel()
	r := queueValidationRegistry()
	if err := validateQueueRegistry(r); err != nil {
		t.Fatal(err)
	}
	if len(r.EpicBaseVersions) != 3 || len(r.EpicBaseUpdates) != 3 || len(r.TaskQueueEvents) != 4 {
		t.Fatalf("fixture has %d bases, %d updates and %d events", len(r.EpicBaseVersions), len(r.EpicBaseUpdates), len(r.TaskQueueEvents))
	}
	assertQueueAccepts(t, map[string]func(*WorkItemRegistry){
		"legacy plan without preparation": basePlanViolation("conflict", func(p *WorkspaceEpicBaseUpdatePlan) { p.CurrentPreparationID = nil }),
		"sha256 object format":            basePlanViolation("committed", func(p *WorkspaceEpicBaseUpdatePlan) { p.ObservedParent.ObjectFormat = "sha256" }),
		"largest expected revision": basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) {
			p.ExpectedRevision = 2147483646
		}),
		"intent against an older base": basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) {
			p.PreviousOID, p.PreviousTree = queueValidationOID("a"), queueValidationOID("b")
		}),
		"second base update": func(r *WorkItemRegistry) { addSecondBaseUpdate(r, func(*WorkspaceEpicBaseUpdatePlan) {}) },
		"second repository":  func(r *WorkItemRegistry) { addQueueValidationBinding(r) },
		"Epic without bindings needs no base": func(r *WorkItemRegistry) {
			r.Epics = append(r.Epics, EpicRecord{ID: "zz", ProjectID: "project"})
		},
		"empty set": queueRequestViolation(goalSetEvent, func(request map[string]any) { request["entries"] = []any{} }),
		// Events of the removed preparation keep their place in the chain; only
		// their revision and expected revision are read.
		"reserve_parallel": func(r *WorkItemRegistry) { r.TaskQueueEvents[preparationEvent].Kind = "reserve_parallel" },
		"advance":          func(r *WorkItemRegistry) { r.TaskQueueEvents[preparationEvent].Kind = "advance" },
		"release":          func(r *WorkItemRegistry) { r.TaskQueueEvents[preparationEvent].Kind = "release" },
		"preparation request is not a draft": queueRequestViolation(preparationEvent, func(request map[string]any) {
			request["entries"], request["kind"] = "junk", 7
		}),
		"first event without expected revision": queueRequestViolation(otherQueueEvent, func(request map[string]any) {
			delete(request, "expected_revision")
		}),
		"set repeats in another queue": func(r *WorkItemRegistry) {
			r.TaskQueueEvents[otherQueueEvent] = sealedQueueEvent(queueValidationID("b"), 1, "set", legacyQueueDraft("queue/third", 0, selectionQueueEntry("unbound", nil)))
		},
	})
}

func TestValidateQueueRegistryEmptyAndNilLists(t *testing.T) {
	t.Parallel()
	if err := validateQueueRegistry(emptyWorkItemRegistry()); err != nil {
		t.Fatalf("empty lists: %v", err)
	}
	// Only the four queue arrays must be explicit here.
	if err := validateQueueRegistry(WorkItemRegistry{EpicBaseVersions: []EpicBaseVersion{}, EpicBaseUpdates: []EpicBaseUpdate{}, TaskQueueEvents: []TaskQueueEvent{}, TaskPreparations: []HistoricalRecord{}}); err != nil {
		t.Fatalf("other lists nil: %v", err)
	}
	assertQueueError(t, validateQueueRegistry(WorkItemRegistry{}), "queue arrays must be explicit")
	const explicit = "queue arrays must be explicit"
	assertQueueViolations(t, []registryViolation{
		{"base versions", func(r *WorkItemRegistry) { r.EpicBaseVersions = nil }, explicit},
		{"base updates", func(r *WorkItemRegistry) { r.EpicBaseUpdates = nil }, explicit},
		{"queue events", func(r *WorkItemRegistry) { r.TaskQueueEvents = nil }, explicit},
		{"preparations", func(r *WorkItemRegistry) { r.TaskPreparations = nil }, explicit},
		// With an Epic, an empty list of bases is explicit but incomplete.
		{"no bases", func(r *WorkItemRegistry) {
			r.EpicBaseVersions, r.EpicBaseUpdates = []EpicBaseVersion{}, []EpicBaseUpdate{}
		}, "Epic lacks base v1"},
	})
}

func TestValidateQueueRegistryZeroRecords(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		mutate func(*WorkItemRegistry)
		want   string
	}{
		"base version": {func(r *WorkItemRegistry) { r.EpicBaseVersions = []EpicBaseVersion{{}} }, "base has foreign Epic"},
		"base update":  {func(r *WorkItemRegistry) { r.EpicBaseUpdates = []EpicBaseUpdate{{}} }, "invalid base update identity"},
		"queue event":  {func(r *WorkItemRegistry) { r.TaskQueueEvents = []TaskQueueEvent{{}} }, "invalid queue event"},
	} {
		t.Run(name, func(t *testing.T) {
			r := emptyWorkItemRegistry()
			tc.mutate(&r)
			assertQueueError(t, validateQueueRegistry(r), tc.want)
		})
	}
	assertQueueViolations(t, []registryViolation{
		{"affected selection", basePlanViolation("committed", func(p *WorkspaceEpicBaseUpdatePlan) {
			p.AffectedSelections = []EpicAffectedSelection{{}}
		}), "invalid affected selection target or order"},
		{"outcome", func(r *WorkItemRegistry) { queueUpdate(r, "committed").Outcome = &EpicBaseUpdateOutcome{} }, "invalid committed base outcome"},
		{"queue request", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[firstSetEvent] = sealedQueueEvent(queueValidationID("a"), 1, "set", nil)
		}, "object has missing or unknown fields"},
	})
}

// The phases run in a fixed order, and the first failing phase answers: the
// arrays, every base update, every base version, a base for every Epic
// binding, every base update against the bases, and every queue event.
func TestValidateQueueRegistryPhasePrecedence(t *testing.T) {
	t.Parallel()
	assertQueuePrecedence(t, []registryViolation{
		{"explicit arrays", func(r *WorkItemRegistry) { r.TaskPreparations = nil }, "queue arrays must be explicit"},
		{"array order", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[firstSetEvent], r.TaskQueueEvents[preparationEvent] = r.TaskQueueEvents[preparationEvent], r.TaskQueueEvents[firstSetEvent]
		}, "queue arrays must be sorted"},
		{"base update", func(r *WorkItemRegistry) { queueUpdate(r, "conflict").Outcome.Reasons = []QueueReason{} }, "invalid conflict outcome"},
		{"base version", func(r *WorkItemRegistry) { r.EpicBaseVersions[epicBaseV2].OID = "abc" }, "invalid base chain"},
		{"base for every Epic", func(r *WorkItemRegistry) { r.EpicBaseVersions = r.EpicBaseVersions[:zetaBaseV1] }, "Epic lacks base v1"},
		{"base update against the bases", basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) {
			p.PreviousOID = queueValidationOID("9")
		}), "base request parent missing"},
		{"queue event", func(r *WorkItemRegistry) { r.TaskQueueEvents[otherQueueEvent].Kind = "bogus" }, "unknown queue event kind"},
	})
}

func TestValidateQueueRegistryRequiresSortedArrays(t *testing.T) {
	t.Parallel()
	const want = "queue arrays must be sorted"
	swap := func(i, j int) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) {
			r.TaskQueueEvents[i], r.TaskQueueEvents[j] = r.TaskQueueEvents[j], r.TaskQueueEvents[i]
		}
	}
	assertQueueViolations(t, []registryViolation{
		{"bases by Epic", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV2], r.EpicBaseVersions[zetaBaseV1] = r.EpicBaseVersions[zetaBaseV1], r.EpicBaseVersions[epicBaseV2]
		}, want},
		{"bases by repository", func(r *WorkItemRegistry) {
			addQueueValidationBinding(r)
			r.EpicBaseVersions[epicBaseV2], r.EpicBaseVersions[epicBaseV2+1] = r.EpicBaseVersions[epicBaseV2+1], r.EpicBaseVersions[epicBaseV2]
		}, want},
		{"bases by revision", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV1], r.EpicBaseVersions[epicBaseV2] = r.EpicBaseVersions[epicBaseV2], r.EpicBaseVersions[epicBaseV1]
		}, want},
		{"updates by ID", func(r *WorkItemRegistry) { slices.Reverse(r.EpicBaseUpdates) }, want},
		{"events by queue", swap(goalSetEvent, otherQueueEvent), want},
		{"events by revision", swap(preparationEvent, goalSetEvent), want},
	})
}

// One base update is checked rule by rule: identity, workspace, target,
// registry format, plan, observation, affected selections, then phase and
// outcome.
func TestValidateQueueRegistryBaseUpdateRulePrecedence(t *testing.T) {
	t.Parallel()
	committed := func(change func(*WorkspaceEpicBaseUpdatePlan)) func(*WorkItemRegistry) {
		return basePlanViolation("committed", change)
	}
	assertQueuePrecedence(t, []registryViolation{
		{"identity", func(r *WorkItemRegistry) { queueUpdate(r, "committed").RecordedAtUTC = "yesterday" }, "invalid base update identity"},
		{"workspace", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.Workspace.Root = "relative" }), "invalid workspace binding"},
		{"target Epic", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ProjectID = "other" }), "foreign queue Epic"},
		{"target binding", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ParentLocator = "/elsewhere" }), "foreign queue target"},
		{"registry format", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.RegistryFormatVersion-- }), fmt.Sprintf("plan must bind registry format %d without an upgrade", RegistryFormat)},
		{"plan", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.UnresolvedOperations = []string{"wop"} }), "invalid base update plan"},
		{"observation", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.ObservedParent.Symbolic = false }), "plan does not bind an exact clean Epic observation"},
		{"affected Task", committed(func(p *WorkspaceEpicBaseUpdatePlan) { p.AffectedSelections[0].TaskID = "absent" }), "invalid affected selection target or order"},
		{"affected selection", committed(func(p *WorkspaceEpicBaseUpdatePlan) {
			p.AffectedSelections[1].Selection.ManifestSHA256 = queueValidationDigest("9")
		}), "selection does not belong to Task native"},
		{"outcome", func(r *WorkItemRegistry) { queueUpdate(r, "committed").Outcome.Revision = intPointer(7) }, "invalid committed base outcome"},
		{"reasons", func(r *WorkItemRegistry) { queueUpdate(r, "committed").Outcome.Reasons = nil }, "null reasons"},
	})
}

func TestValidateQueueRegistryBaseUpdateIdentity(t *testing.T) {
	t.Parallel()
	const want = "invalid base update identity"
	update := func(change func(*EpicBaseUpdate)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) {
			change(queueUpdate(r, "committed"))
			sortQueueRegistry(r)
		}
	}
	assertQueueViolations(t, []registryViolation{
		{"ID of another plan", update(func(u *EpicBaseUpdate) { u.ID = "ebu_" + strings.Repeat("9", 64) }), want},
		{"ID without prefix", update(func(u *EpicBaseUpdate) { u.ID = strings.TrimPrefix(u.ID, "ebu_") }), want},
		{"digest of another plan", update(func(u *EpicBaseUpdate) {
			u.PlanSHA256, u.ID = queueValidationDigest("9"), "ebu_"+strings.Repeat("9", 64)
		}), want},
		{"plan changed after sealing", update(func(u *EpicBaseUpdate) { u.Plan.NextTree = queueValidationOID("9") }), want},
		{"duplicate", func(r *WorkItemRegistry) {
			r.EpicBaseUpdates = append(r.EpicBaseUpdates, *queueUpdate(r, "intent"))
			sortQueueRegistry(r)
		}, want},
		{"time without Z", update(func(u *EpicBaseUpdate) { u.RecordedAtUTC = "2026-11-20T10:00:00+01:00" }), want},
		{"time not a timestamp", update(func(u *EpicBaseUpdate) { u.RecordedAtUTC = "noonZ" }), want},
		{"time empty", update(func(u *EpicBaseUpdate) { u.RecordedAtUTC = "" }), want},
	})
}

func TestValidateQueueRegistryBaseUpdateEnvelope(t *testing.T) {
	t.Parallel()
	format := fmt.Sprintf("plan must bind registry format %d without an upgrade", RegistryFormat)
	plan := func(change func(*WorkspaceEpicBaseUpdatePlan)) func(*WorkItemRegistry) {
		return basePlanViolation("intent", change)
	}
	assertQueueViolations(t, []registryViolation{
		{"workspace root", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Workspace.Root = "/validation/../x" }), "invalid workspace binding"},
		{"workspace marker", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Workspace.MarkerSHA256 = "sha256:abc" }), "invalid workspace binding"},
		{"unknown Epic", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.EpicID = "ghost" }), "foreign queue Epic"},
		{"other Project", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ProjectID = "other" }), "foreign queue Epic"},
		{"unknown repository", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.RepoID = "ghost" }), "foreign queue target"},
		{"common directory", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.GitCommonDir = "/elsewhere/.git" }), "foreign queue target"},
		{"worktree ID", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ParentWorktreeID = validationWorktreeID("9") }), "foreign queue target"},
		{"ref", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ParentRef = "refs/heads/elsewhere" }), "foreign queue target"},
		{"locator", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.Target.ParentLocator = "/elsewhere" }), "foreign queue target"},
		{"older format", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.RegistryFormatVersion = RegistryFormat - 1 }), format},
		{"newer format", plan(func(p *WorkspaceEpicBaseUpdatePlan) { p.RegistryFormatVersion = RegistryFormat + 1 }), format},
		{"upgrade", plan(func(p *WorkspaceEpicBaseUpdatePlan) {
			p.RegistryUpgrade = &TaskRegistryUpgrade{FromVersion: RegistryFormat}
		}), format},
	})
}

func TestValidateQueueRegistryBaseUpdatePlan(t *testing.T) {
	t.Parallel()
	const want = "invalid base update plan"
	plan := func(name string, change func(*WorkspaceEpicBaseUpdatePlan)) registryViolation {
		return registryViolation{name, basePlanViolation("intent", change), want}
	}
	// The observation follows a changed next OID or tree, so that only the rule
	// on the value itself is broken.
	nextOID := func(p *WorkspaceEpicBaseUpdatePlan, oid string) { p.NextOID, p.ObservedParent.OID = oid, oid }
	nextTree := func(p *WorkspaceEpicBaseUpdatePlan, tree string) { p.NextTree, p.ObservedParent.Tree = tree, tree }
	assertQueueViolations(t, []registryViolation{
		plan("unknown kind", func(p *WorkspaceEpicBaseUpdatePlan) { p.Kind, p.SchemaVersion = "WorkspaceEpicBaseUpdatePlan@2", 2 }),
		plan("current kind with legacy schema", func(p *WorkspaceEpicBaseUpdatePlan) { p.SchemaVersion = 1 }),
		plan("legacy kind with current schema", func(p *WorkspaceEpicBaseUpdatePlan) { p.Kind = "WorkspaceEpicBaseUpdatePlan@1" }),
		plan("current plan with a preparation", func(p *WorkspaceEpicBaseUpdatePlan) { p.CurrentPreparationID = stringPointer("prep") }),
		plan("expected revision zero", func(p *WorkspaceEpicBaseUpdatePlan) { p.ExpectedRevision = 0 }),
		plan("expected revision at the limit", func(p *WorkspaceEpicBaseUpdatePlan) { p.ExpectedRevision = 2147483647 }),
		plan("next OID", func(p *WorkspaceEpicBaseUpdatePlan) { nextOID(p, "abc") }),
		plan("next tree", func(p *WorkspaceEpicBaseUpdatePlan) { nextTree(p, strings.Repeat("G", 40)) }),
		plan("next OID is the previous one", func(p *WorkspaceEpicBaseUpdatePlan) { nextOID(p, p.PreviousOID) }),
		plan("affected selections null", func(p *WorkspaceEpicBaseUpdatePlan) { p.AffectedSelections = nil }),
		plan("unresolved operations null", func(p *WorkspaceEpicBaseUpdatePlan) { p.UnresolvedOperations = nil }),
		plan("unresolved operation", func(p *WorkspaceEpicBaseUpdatePlan) { p.UnresolvedOperations = []string{"wop"} }),
		plan("observed parent dirty", func(p *WorkspaceEpicBaseUpdatePlan) { p.ObservedParent.Clean = false }),
		plan("observed OID", func(p *WorkspaceEpicBaseUpdatePlan) { p.ObservedParent.OID = queueValidationOID("9") }),
		plan("observed tree", func(p *WorkspaceEpicBaseUpdatePlan) { p.ObservedParent.Tree = queueValidationOID("9") }),
	})
}

func TestValidateQueueRegistryBaseUpdateObservation(t *testing.T) {
	t.Parallel()
	const want = "plan does not bind an exact clean Epic observation"
	observation := func(name string, change func(*PlanWorktreeObservation)) registryViolation {
		return registryViolation{name, basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) { change(&p.ObservedParent) }), want}
	}
	assertQueueViolations(t, []registryViolation{
		observation("locator", func(o *PlanWorktreeObservation) { o.WorktreeLocator = "/elsewhere" }),
		observation("ref", func(o *PlanWorktreeObservation) { o.Ref = "refs/heads/elsewhere" }),
		observation("common directory", func(o *PlanWorktreeObservation) { o.GitCommonDir = "/elsewhere/.git" }),
		observation("detached", func(o *PlanWorktreeObservation) { o.Symbolic = false }),
		observation("status entries null", func(o *PlanWorktreeObservation) { o.StatusEntries = nil }),
		observation("status entry", func(o *PlanWorktreeObservation) { o.StatusEntries = []IntegrationStatusEntry{{RecordKind: "ordinary"}} }),
		observation("in progress null", func(o *PlanWorktreeObservation) { o.InProgress = nil }),
		observation("in progress", func(o *PlanWorktreeObservation) { o.InProgress = []string{"merge"} }),
		observation("object format", func(o *PlanWorktreeObservation) { o.ObjectFormat = "md5" }),
		observation("ref format", func(o *PlanWorktreeObservation) { o.RefFormat = "reftable" }),
	})
}

func TestValidateQueueRegistryAffectedSelections(t *testing.T) {
	t.Parallel()
	const target = "invalid affected selection target or order"
	affected := func(change func(*[]EpicAffectedSelection)) func(*WorkItemRegistry) {
		return basePlanViolation("committed", func(p *WorkspaceEpicBaseUpdatePlan) { change(&p.AffectedSelections) })
	}
	assertQueueViolations(t, []registryViolation{
		{"descending", affected(func(a *[]EpicAffectedSelection) { slices.Reverse(*a) }), target},
		{"duplicate Task", affected(func(a *[]EpicAffectedSelection) { (*a)[1] = (*a)[0] }), target},
		{"unknown Task", affected(func(a *[]EpicAffectedSelection) { (*a)[0].TaskID = "absent" }), target},
		{"Task of another Project", func(r *WorkItemRegistry) { r.Tasks[creatingTask].ProjectID = "other" }, target},
		{"Task of another repository", func(r *WorkItemRegistry) { r.Tasks[creatingTask].RepoID = "other" }, target},
		{"Task of another Epic", func(r *WorkItemRegistry) { r.Tasks[creatingTask].ParentEpicID = "zeta" }, target},
		// Within one entry, the Task answers before the selection it does not own.
		{"unknown Task with a selection", affected(func(a *[]EpicAffectedSelection) { (*a)[1].TaskID = "absent" }), target},
		{"unknown selection", affected(func(a *[]EpicAffectedSelection) {
			(*a)[1].Selection.ID = "sel_" + strings.Repeat("9", 32)
		}), "selection does not belong to Task native"},
		{"selection manifest", affected(func(a *[]EpicAffectedSelection) {
			(*a)[1].Selection.ManifestSHA256 = queueValidationDigest("9")
		}), "selection does not belong to Task native"},
		{"selection of another Task", affected(func(a *[]EpicAffectedSelection) {
			(*a)[0].Selection = (*a)[1].Selection
		}), "selection does not belong to Task creating"},
		// The entries are checked one at a time, in plan order.
		{"first entry's selection before second entry's Task", affected(func(a *[]EpicAffectedSelection) {
			(*a)[0].Selection, (*a)[1].TaskID = (*a)[1].Selection, "absent"
		}), "selection does not belong to Task creating"},
		{"first entry's Task before second entry's selection", affected(func(a *[]EpicAffectedSelection) {
			(*a)[0].TaskID, (*a)[1].Selection.ID = "absent", "sel_"+strings.Repeat("9", 32)
		}), target},
	})
}

func TestValidateQueueRegistryBaseUpdatePhaseAndOutcome(t *testing.T) {
	t.Parallel()
	update := func(phase string, change func(*EpicBaseUpdate)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { change(queueUpdate(r, phase)) }
	}
	reasons := []QueueReason{{Code: "code", Message: "Message"}}
	assertQueueViolations(t, []registryViolation{
		{"intent with outcome", update("intent", func(u *EpicBaseUpdate) { u.Outcome = &EpicBaseUpdateOutcome{Reasons: []QueueReason{}} }), "intent has outcome"},
		{"committed without outcome", update("committed", func(u *EpicBaseUpdate) { u.Outcome = nil }), "invalid committed base outcome"},
		{"committed without revision", update("committed", func(u *EpicBaseUpdate) { u.Outcome.Revision = nil }), "invalid committed base outcome"},
		{"committed at the expected revision", update("committed", func(u *EpicBaseUpdate) { u.Outcome.Revision = intPointer(1) }), "invalid committed base outcome"},
		{"committed with reasons", update("committed", func(u *EpicBaseUpdate) { u.Outcome.Reasons = reasons }), "invalid committed base outcome"},
		{"committed with null reasons", update("committed", func(u *EpicBaseUpdate) { u.Outcome.Reasons = nil }), "null reasons"},
		{"conflict without outcome", update("conflict", func(u *EpicBaseUpdate) { u.Outcome = nil }), "invalid conflict outcome"},
		{"conflict with revision", update("conflict", func(u *EpicBaseUpdate) { u.Outcome.Revision = intPointer(2) }), "invalid conflict outcome"},
		{"conflict without reasons", update("conflict", func(u *EpicBaseUpdate) { u.Outcome.Reasons = []QueueReason{} }), "invalid conflict outcome"},
		{"conflict with null reasons", update("conflict", func(u *EpicBaseUpdate) { u.Outcome.Reasons = nil }), "invalid conflict outcome"},
		{"unknown phase", update("intent", func(u *EpicBaseUpdate) { u.Phase = "bogus" }), "invalid base phase"},
		{"empty phase", update("conflict", func(u *EpicBaseUpdate) { u.Phase = "" }), "invalid base phase"},
		// The phase is read before the outcome's reasons.
		{"unknown phase with null reasons", update("committed", func(u *EpicBaseUpdate) { u.Phase, u.Outcome.Reasons = "bogus", nil }), "invalid base phase"},
	})
}

// Base updates are checked completely one at a time, in ID order: the last rule
// of an earlier update wins over the first rule of a later one.
func TestValidateQueueRegistryBaseUpdatesInIDOrder(t *testing.T) {
	t.Parallel()
	identity := func(u *EpicBaseUpdate) { u.RecordedAtUTC = "yesterday" }
	phase := func(u *EpicBaseUpdate) { u.Phase = "bogus" }
	assertQueueViolations(t, []registryViolation{
		{"phase first", func(r *WorkItemRegistry) { phase(&r.EpicBaseUpdates[0]); identity(&r.EpicBaseUpdates[1]) }, "invalid base phase"},
		{"identity first", func(r *WorkItemRegistry) { identity(&r.EpicBaseUpdates[1]); phase(&r.EpicBaseUpdates[2]) }, "invalid base update identity"},
	})
}

// One base version is checked rule by rule: its Epic binding, its place and
// values in the chain, then either the adoption (v1) or the link to the
// previous version, the reuse of an OID and the committed update.
func TestValidateQueueRegistryBaseVersionRulePrecedence(t *testing.T) {
	t.Parallel()
	v2 := func(change func(*EpicBaseVersion)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { change(&r.EpicBaseVersions[epicBaseV2]) }
	}
	assertQueuePrecedence(t, []registryViolation{
		{"chain", v2(func(v *EpicBaseVersion) { v.Locator = "/elsewhere" }), "invalid base chain"},
		{"link", v2(func(v *EpicBaseVersion) { v.PreviousRevision = nil }), "broken base chain"},
		{"reused OID", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV2].OID = r.EpicBaseVersions[epicBaseV1].OID
		}, "reused base OID"},
		{"update", v2(func(v *EpicBaseVersion) { v.Tree = queueValidationOID("9") }), "base differs from update"},
	})
	zeta := func(change func(*EpicBaseVersion)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { change(&r.EpicBaseVersions[zetaBaseV1]) }
	}
	assertQueuePrecedence(t, []registryViolation{
		{"first chain", zeta(func(v *EpicBaseVersion) { v.Tree = "abc" }), "invalid base chain"},
		{"adoption", zeta(func(v *EpicBaseVersion) { v.OID = queueValidationOID("9") }), "base v1 differs from adoption"},
	})
}

func TestValidateQueueRegistryBaseVersionBinding(t *testing.T) {
	t.Parallel()
	assertQueueViolations(t, []registryViolation{
		{"unknown Epic", func(r *WorkItemRegistry) { r.Epics = r.Epics[:1] }, "base has foreign Epic"},
		{"unknown repository", func(r *WorkItemRegistry) { r.Epics[1].RepoBindings[0].RepoID = "other" }, "base has foreign repo"},
		{"no base for the Epic", func(r *WorkItemRegistry) { r.EpicBaseVersions = r.EpicBaseVersions[:zetaBaseV1] }, "Epic lacks base v1"},
		{"no base for one binding", func(r *WorkItemRegistry) { addValidationBinding(r) }, "Epic lacks base v1"},
	})
}

func TestValidateQueueRegistryBaseChain(t *testing.T) {
	t.Parallel()
	const want = "invalid base chain"
	v2 := func(name string, change func(*EpicBaseVersion)) registryViolation {
		return registryViolation{name, func(r *WorkItemRegistry) { change(&r.EpicBaseVersions[epicBaseV2]) }, want}
	}
	assertQueueViolations(t, []registryViolation{
		v2("revision gap", func(v *EpicBaseVersion) { v.Revision = 3 }),
		{"first base is not v1", func(r *WorkItemRegistry) { r.EpicBaseVersions = r.EpicBaseVersions[epicBaseV2:] }, want},
		{"repeated revision", func(r *WorkItemRegistry) {
			r.EpicBaseVersions = slices.Insert(r.EpicBaseVersions, epicBaseV2, r.EpicBaseVersions[epicBaseV1])
		}, want},
		v2("Project", func(v *EpicBaseVersion) { v.ProjectID = "other" }),
		v2("worktree ID", func(v *EpicBaseVersion) { v.WorktreeID = validationWorktreeID("9") }),
		v2("locator", func(v *EpicBaseVersion) { v.Locator = "/elsewhere" }),
		v2("ref", func(v *EpicBaseVersion) { v.Ref = "refs/heads/elsewhere" }),
		v2("common directory", func(v *EpicBaseVersion) { v.GitCommonDir = "/elsewhere/.git" }),
		v2("OID", func(v *EpicBaseVersion) { v.OID = "abc" }),
		v2("tree", func(v *EpicBaseVersion) { v.Tree = "abc" }),
		// An Epic binding with an invalid OID is rejected here, with its base v1.
		{"adopted OID", func(r *WorkItemRegistry) {
			r.Epics[1].RepoBindings[0].Worktree.OID = "abc"
			r.EpicBaseVersions[zetaBaseV1].OID = "abc"
		}, want},
	})
}

func TestValidateQueueRegistryBaseAdoption(t *testing.T) {
	t.Parallel()
	const want = "base v1 differs from adoption"
	v1 := func(name string, change func(*EpicBaseVersion)) registryViolation {
		return registryViolation{name, func(r *WorkItemRegistry) { change(&r.EpicBaseVersions[zetaBaseV1]) }, want}
	}
	assertQueueViolations(t, []registryViolation{
		v1("OID", func(v *EpicBaseVersion) { v.OID = queueValidationOID("9") }),
		v1("tree", func(v *EpicBaseVersion) { v.Tree = queueValidationOID("9") }),
		v1("previous revision", func(v *EpicBaseVersion) { v.PreviousRevision = intPointer(0) }),
		v1("update", func(v *EpicBaseVersion) { v.UpdateID = stringPointer("ebu_x") }),
	})
}

func TestValidateQueueRegistryBaseSuccessor(t *testing.T) {
	t.Parallel()
	v2 := func(change func(*EpicBaseVersion)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { change(&r.EpicBaseVersions[epicBaseV2]) }
	}
	committed := func(change func(*WorkItemRegistry, *WorkspaceEpicBaseUpdatePlan)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) {
			changeBasePlan(r, "committed", func(p *WorkspaceEpicBaseUpdatePlan) { change(r, p) })
		}
	}
	const update = "base differs from update"
	assertQueueViolations(t, []registryViolation{
		{"no previous revision", v2(func(v *EpicBaseVersion) { v.PreviousRevision = nil }), "broken base chain"},
		{"wrong previous revision", v2(func(v *EpicBaseVersion) { v.PreviousRevision = intPointer(2) }), "broken base chain"},
		{"no update", v2(func(v *EpicBaseVersion) { v.UpdateID = nil }), "broken base chain"},
		{"OID of the previous base", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV2].OID = r.EpicBaseVersions[epicBaseV1].OID
		}, "reused base OID"},
		// The update of v3 is valid and matches v3; only the OID has been a base before.
		{"OID of an older base", func(r *WorkItemRegistry) {
			addSecondBaseUpdate(r, func(p *WorkspaceEpicBaseUpdatePlan) {
				p.NextOID, p.ObservedParent.OID = r.EpicBaseVersions[epicBaseV1].OID, r.EpicBaseVersions[epicBaseV1].OID
			})
		}, "reused base OID"},
		// Every earlier base counts, not only the first: v3 repeats v2.
		{"OID of the previous of three bases", func(r *WorkItemRegistry) {
			addSecondBaseUpdate(r, func(*WorkspaceEpicBaseUpdatePlan) {})
			r.EpicBaseVersions[epicBaseV2+1].OID = r.EpicBaseVersions[epicBaseV2].OID
		}, "reused base OID"},
		{"unknown update", v2(func(v *EpicBaseVersion) { v.UpdateID = stringPointer("ebu_" + strings.Repeat("9", 64)) }), update},
		// Only the phase is wrong: the plan moves the base exactly from v1 to v2.
		{"update is an intent", func(r *WorkItemRegistry) {
			applied := queueUpdate(r, "committed")
			applied.Phase, applied.Outcome = "intent", nil
		}, update},
		{"update is a conflict", func(r *WorkItemRegistry) {
			applied := queueUpdate(r, "committed")
			applied.Phase, applied.Outcome = "conflict", queueUpdate(r, "conflict").Outcome
		}, update},
		{"update of another Epic", committed(func(r *WorkItemRegistry, p *WorkspaceEpicBaseUpdatePlan) {
			retargetQueuePlan(p, r.Epics[1], r.Epics[1].RepoBindings[0])
			p.AffectedSelections = []EpicAffectedSelection{}
		}), update},
		{"update of another repository", func(r *WorkItemRegistry) {
			binding := addQueueValidationBinding(r)
			changeBasePlan(r, "committed", func(p *WorkspaceEpicBaseUpdatePlan) {
				retargetQueuePlan(p, r.Epics[0], binding)
				p.AffectedSelections = []EpicAffectedSelection{}
			})
		}, update},
		{"update expected another revision", func(r *WorkItemRegistry) {
			changeBasePlan(r, "committed", func(p *WorkspaceEpicBaseUpdatePlan) { p.ExpectedRevision = 2 })
			queueUpdate(r, "committed").Outcome.Revision = intPointer(3)
		}, update},
		{"update from another OID", committed(func(_ *WorkItemRegistry, p *WorkspaceEpicBaseUpdatePlan) { p.PreviousOID = queueValidationOID("9") }), update},
		{"update from another tree", committed(func(_ *WorkItemRegistry, p *WorkspaceEpicBaseUpdatePlan) { p.PreviousTree = queueValidationOID("9") }), update},
		{"update to another OID", committed(func(_ *WorkItemRegistry, p *WorkspaceEpicBaseUpdatePlan) {
			p.NextOID, p.ObservedParent.OID = queueValidationOID("9"), queueValidationOID("9")
		}), update},
		{"update to another tree", committed(func(_ *WorkItemRegistry, p *WorkspaceEpicBaseUpdatePlan) {
			p.NextTree, p.ObservedParent.Tree = queueValidationOID("9"), queueValidationOID("9")
		}), update},
	})
}

// Base versions are checked completely one at a time, in registry order, and
// all of them before any Epic is asked for its base.
func TestValidateQueueRegistryBaseVersionsInRegistryOrder(t *testing.T) {
	t.Parallel()
	link := func(r *WorkItemRegistry) { r.EpicBaseVersions[epicBaseV2].UpdateID = nil }
	adoption := func(r *WorkItemRegistry) { r.EpicBaseVersions[zetaBaseV1].Tree = queueValidationOID("9") }
	assertQueueViolations(t, []registryViolation{
		{"update of v2 before adoption of the next Epic", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV2].Tree = queueValidationOID("9")
			adoption(r)
		}, "base differs from update"},
		{"link of v2 before adoption of the next Epic", func(r *WorkItemRegistry) { link(r); adoption(r) }, "broken base chain"},
		{"adoption before chain of the next Epic", func(r *WorkItemRegistry) {
			r.EpicBaseVersions[epicBaseV1].Tree = queueValidationOID("9")
			r.EpicBaseVersions[zetaBaseV1].Locator = "/elsewhere"
		}, "base v1 differs from adoption"},
		{"foreign repository before a missing base", func(r *WorkItemRegistry) {
			addValidationBinding(r)
			r.Epics[1].RepoBindings[0].RepoID = "other"
		}, "base has foreign repo"},
	})
}

func TestValidateQueueRegistryBaseUpdatesAgainstBases(t *testing.T) {
	t.Parallel()
	const parent, version = "base request parent missing", "committed update lacks version"
	lostParent := basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) { p.PreviousOID = queueValidationOID("9") })
	assertQueueViolations(t, []registryViolation{
		{"previous OID was never a base", lostParent, parent},
		{"previous tree was never a base", basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) { p.PreviousTree = queueValidationOID("9") }), parent},
		{"previous OID and tree of different bases", basePlanViolation("intent", func(p *WorkspaceEpicBaseUpdatePlan) { p.PreviousTree = queueValidationOID("b") }), parent},
		{"committed update without a version", func(r *WorkItemRegistry) { addUnappliedUpdate(r, true) }, version},
		{"version of the update was removed", func(r *WorkItemRegistry) {
			r.EpicBaseVersions = slices.Delete(r.EpicBaseVersions, epicBaseV2, epicBaseV2+1)
			changeBasePlan(r, "intent", func(p *WorkspaceEpicBaseUpdatePlan) {
				p.PreviousOID, p.PreviousTree = queueValidationOID("a"), queueValidationOID("b")
			})
		}, version},
		// One update is checked completely before the next, in ID order.
		{"missing version of an earlier update", func(r *WorkItemRegistry) { lostParent(r); addUnappliedUpdate(r, true) }, version},
		{"missing parent of an earlier update", func(r *WorkItemRegistry) { lostParent(r); addUnappliedUpdate(r, false) }, parent},
		{"missing parent before missing version of one update", func(r *WorkItemRegistry) {
			addUnappliedUpdate(r, true)
			update := &r.EpicBaseUpdates[slices.IndexFunc(r.EpicBaseUpdates, func(u EpicBaseUpdate) bool {
				return strings.HasPrefix(u.Plan.Workspace.Root, "/validation/")
			})]
			update.Plan.PreviousOID = queueValidationOID("9")
			resealBaseUpdate(update)
			sortQueueRegistry(r)
		}, parent},
	})
}

// One queue event is checked rule by rule: the event itself, its place in the
// queue's chain, then for a set the draft, its publication key, its target and
// its entries.
func TestValidateQueueRegistryEventRulePrecedence(t *testing.T) {
	t.Parallel()
	request := func(change func(request map[string]any)) func(*WorkItemRegistry) {
		return queueRequestViolation(goalSetEvent, change)
	}
	assertQueuePrecedence(t, []registryViolation{
		{"event", func(r *WorkItemRegistry) { r.TaskQueueEvents[goalSetEvent].RecordedAtUTC = "yesterday" }, "invalid queue event"},
		{"chain", request(func(request map[string]any) { request["expected_revision"] = 7 }), "broken queue CAS chain"},
		{"draft", request(func(request map[string]any) { request["kind"] = "bogus" }), "kind: must be one of WorkspaceTaskQueueDraft@2"},
		{"publication key", request(func(request map[string]any) { request["publication_key"] = "queue/first" }), "duplicate queue publication key"},
		{"target Epic", request(func(request map[string]any) { request["project_id"] = "other" }), "foreign queue target"},
		{"target repository", request(func(request map[string]any) { request["repo_id"] = "ghost" }), "foreign queue repo"},
		{"Task", request(func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "ghost" }), "foreign queue Task"},
	})
	assertQueuePrecedence(t, []registryViolation{
		{"chain", func(r *WorkItemRegistry) { r.TaskQueueEvents[otherQueueEvent].Revision = 2 }, "broken queue CAS chain"},
		{"kind", func(r *WorkItemRegistry) { r.TaskQueueEvents[otherQueueEvent].Kind = "" }, "unknown queue event kind"},
	})
}

func TestValidateQueueRegistryEventIdentity(t *testing.T) {
	t.Parallel()
	const want = "invalid queue event"
	event := func(change func(*TaskQueueEvent)) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) {
			change(&r.TaskQueueEvents[otherQueueEvent])
			sortQueueRegistry(r)
		}
	}
	assertQueueViolations(t, []registryViolation{
		{"queue ID without prefix", event(func(e *TaskQueueEvent) { e.QueueID = strings.Repeat("b", 64) }), want},
		{"queue ID too short", event(func(e *TaskQueueEvent) { e.QueueID = "que_" + strings.Repeat("b", 63) }), want},
		{"queue ID not hexadecimal", event(func(e *TaskQueueEvent) { e.QueueID = "que_" + strings.Repeat("z", 64) }), want},
		{"time without Z", event(func(e *TaskQueueEvent) { e.RecordedAtUTC = "2026-11-20T10:00:00+01:00" }), want},
		{"time empty", event(func(e *TaskQueueEvent) { e.RecordedAtUTC = "" }), want},
		{"request changed after sealing", event(func(e *TaskQueueEvent) { e.Request["expected_revision"] = 1 }), want},
		{"digest of another request", event(func(e *TaskQueueEvent) { e.RequestSHA256 = queueValidationDigest("9") }), want},
	})
}

func TestValidateQueueRegistryEventChain(t *testing.T) {
	t.Parallel()
	const want = "broken queue CAS chain"
	assertQueueViolations(t, []registryViolation{
		{"first revision is not 1", func(r *WorkItemRegistry) { r.TaskQueueEvents[otherQueueEvent].Revision = 2 }, want},
		{"first revision is 0", func(r *WorkItemRegistry) { r.TaskQueueEvents[otherQueueEvent].Revision = 0 }, want},
		{"first expected revision is not 0", queueRequestViolation(otherQueueEvent, func(request map[string]any) { request["expected_revision"] = 1 }), want},
		{"revision gap", func(r *WorkItemRegistry) {
			r.TaskQueueEvents = slices.Delete(r.TaskQueueEvents, preparationEvent, preparationEvent+1)
		}, want},
		{"repeated revision", func(r *WorkItemRegistry) {
			r.TaskQueueEvents = slices.Insert(r.TaskQueueEvents, preparationEvent, r.TaskQueueEvents[firstSetEvent])
		}, want},
		{"expected revision behind", queueRequestViolation(goalSetEvent, func(request map[string]any) { request["expected_revision"] = 1 }), want},
		{"expected revision ahead", queueRequestViolation(goalSetEvent, func(request map[string]any) { request["expected_revision"] = 3 }), want},
		// A missing or non-integer expected revision reads as 0.
		{"expected revision missing", queueRequestViolation(preparationEvent, func(request map[string]any) { delete(request, "expected_revision") }), want},
		{"expected revision as text", queueRequestViolation(preparationEvent, func(request map[string]any) { request["expected_revision"] = "1" }), want},
		// Each queue has its own chain: the second queue does not continue the first.
		{"chain continued from another queue", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[otherQueueEvent] = sealedQueueEvent(queueValidationID("b"), 4, "release", map[string]any{"expected_revision": 3})
		}, want},
	})
}

func TestValidateQueueRegistryEventKind(t *testing.T) {
	t.Parallel()
	const want = "unknown queue event kind"
	kind := func(kind string) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { r.TaskQueueEvents[preparationEvent].Kind = kind }
	}
	assertQueueViolations(t, []registryViolation{
		{"unknown", kind("bogus"), want},
		{"empty", kind(""), want},
		{"other case", kind("Set"), want},
		// A set is validated as a draft whatever its request was written for.
		{"preparation request as set", kind("set"), "object has missing or unknown fields"},
	})
}

func TestValidateQueueRegistrySetDraft(t *testing.T) {
	t.Parallel()
	legacy := func(change func(request map[string]any)) func(*WorkItemRegistry) {
		return queueRequestViolation(firstSetEvent, change)
	}
	goal := func(change func(request map[string]any)) func(*WorkItemRegistry) {
		return queueRequestViolation(goalSetEvent, change)
	}
	assertQueueViolations(t, []registryViolation{
		{"legacy kind", legacy(func(request map[string]any) { request["kind"] = "WorkspaceTaskQueueDraft@2" }), "kind: must be one of WorkspaceTaskQueueDraft@1"},
		{"goal kind", goal(func(request map[string]any) { request["kind"] = "WorkspaceTaskQueueDraft@1" }), "kind: must be one of WorkspaceTaskQueueDraft@2"},
		{"unknown schema", legacy(func(request map[string]any) { request["schema_version"] = 3 }), "schema_version: must be an integer in 1..1"},
		{"missing field", legacy(func(request map[string]any) { delete(request, "human_decision") }), "object has missing or unknown fields"},
		{"entries null", legacy(func(request map[string]any) { request["entries"] = nil }), "entries: must be an explicit array with at most 256 entries"},
		{"repeated Task", goal(func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "unbound" }), "entries: entries must be unique"},
		{"goal revision", goal(func(request map[string]any) {
			queueRequestEntry(request, 1)["goal"].(map[string]any)["spec"].(map[string]any)["revision"] = 0
		}), "entries: entry 1: goal: spec: revision: must be an integer in 1..2147483647"},
		{"registry upgrade", legacy(func(request map[string]any) { request["registry_upgrade"] = map[string]any{} }), fmt.Sprintf("registry_upgrade: must be null; registry format %d has no upgrade", RegistryFormat)},
	})
}

// A draft's own errors keep the class that the content rules gave them.
func TestValidateQueueRegistryClassifiedErrors(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		mutate      func(*WorkItemRegistry)
		want, class string
	}{
		"unknown field": {queueRequestViolation(firstSetEvent, func(request map[string]any) { request["extra"] = true }),
			"object has missing or unknown fields", "unknown fields extra"},
		"field of the other schema": {queueRequestViolation(firstSetEvent, func(request map[string]any) { request["schema_version"] = 2 }),
			"unknown or duplicate field human_decision", "unknown fields human_decision"},
		// Only schema 2 selects the goal draft; every other value is a legacy
		// draft, which has no goal entries.
		"goal draft of an unknown schema": {queueRequestViolation(goalSetEvent, func(request map[string]any) { request["schema_version"] = 3 }),
			"entries: entry 1: unknown or duplicate field goal", "unknown fields goal"},
		"Task ID": {queueRequestViolation(firstSetEvent, func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "Bad ID" }),
			"entries: entry 0: task_id: " + invalidIdentifierError("content", "Bad ID"), "work " + string(workspace.ErrorWorkInvalidArguments)},
		"unregistered goal": {func(r *WorkItemRegistry) { r.TaskSpecRevisions[0].Revision = 2 },
			"task_goal_invalid_reference: goal reference is not registered to this Task", "content task_goal_invalid_reference"},
	} {
		t.Run(name, func(t *testing.T) {
			r := queueValidationRegistry()
			tc.mutate(&r)
			assertQueueError(t, validateQueueRegistry(r), tc.want, tc.class)
		})
	}
}

func TestValidateQueueRegistrySetTargetAndEntries(t *testing.T) {
	t.Parallel()
	const task, goal = "foreign queue Task", "task_goal_invalid_reference: goal reference is not registered to this Task"
	request := func(event int, change func(request map[string]any)) func(*WorkItemRegistry) {
		return queueRequestViolation(event, change)
	}
	otherSelection := map[string]any{"id": "sel_" + strings.Repeat("1", 32), "manifest_sha256": queueValidationDigest("1")}
	for _, violation := range []registryViolation{
		{"duplicate publication key", request(goalSetEvent, func(request map[string]any) { request["publication_key"] = "queue/first" }), "duplicate queue publication key"},
		// Publication keys are unique across all queues.
		{"publication key of another queue", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[otherQueueEvent] = sealedQueueEvent(queueValidationID("b"), 1, "set", legacyQueueDraft("queue/second", 0))
		}, "duplicate queue publication key"},
		{"unknown Epic", request(firstSetEvent, func(request map[string]any) { request["epic_id"] = "ghost" }), "foreign queue target"},
		{"Epic of another Project", request(firstSetEvent, func(request map[string]any) { request["project_id"] = "other" }), "foreign queue target"},
		{"unknown repository", request(firstSetEvent, func(request map[string]any) { request["repo_id"] = "ghost" }), "foreign queue repo"},
		{"unknown Task", request(firstSetEvent, func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "ghost" }), task},
		{"Task of another Project", func(r *WorkItemRegistry) { r.Tasks[unboundTask].ProjectID = "other" }, task},
		{"Task of another repository", func(r *WorkItemRegistry) { r.Tasks[unboundTask].RepoID = "other" }, task},
		{"Task of another Epic", func(r *WorkItemRegistry) { r.Tasks[unboundTask].ParentEpicID = "zeta" }, task},
		{"Task of the other Epic's queue", request(firstSetEvent, func(request map[string]any) { request["epic_id"] = "zeta" }), task},
		{"unknown selection", request(firstSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 0)["selection"].(map[string]any)["id"] = "sel_" + strings.Repeat("9", 32)
		}), "selection does not belong to Task unbound"},
		{"selection manifest", request(firstSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 0)["selection"].(map[string]any)["manifest_sha256"] = queueValidationDigest("9")
		}), "selection does not belong to Task unbound"},
		{"selection of another Task", request(goalSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 0)["selection"] = otherSelection
		}), "selection does not belong to Task reconciling"},
		{"goal of another Spec", request(goalSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 1)["goal"].(map[string]any)["spec_id"] = "other"
		}), goal},
		{"goal manifest", request(goalSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 1)["goal"].(map[string]any)["spec"].(map[string]any)["manifest_sha256"] = queueValidationDigest("9")
		}), goal},
		{"goal of another Task", func(r *WorkItemRegistry) { r.TaskSpecRevisions[0].TaskID = "ready" }, goal},
		// The entries are checked one at a time, in draft order.
		{"first entry's selection before second entry's Task", request(goalSetEvent, func(request map[string]any) {
			queueRequestEntry(request, 0)["selection"] = otherSelection
			queueRequestEntry(request, 1)["task_id"] = "ghost"
		}), "selection does not belong to Task reconciling"},
		{"first entry's Task before second entry's goal", func(r *WorkItemRegistry) {
			r.TaskSpecRevisions = []TaskSpecReference{}
			changeQueueRequest(r, goalSetEvent, func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "ghost" })
		}, task},
	} {
		t.Run(violation.name, func(t *testing.T) {
			r := queueValidationRegistry()
			violation.mutate(&r)
			class := ""
			if violation.want == goal {
				class = "content task_goal_invalid_reference"
			}
			assertQueueError(t, validateQueueRegistry(r), violation.want, class)
		})
	}
}

// Queue events are checked completely one at a time, in registry order: the
// last rule of an earlier event wins over the first rule of a later one.
func TestValidateQueueRegistryEventsInRegistryOrder(t *testing.T) {
	t.Parallel()
	entry := queueRequestViolation(firstSetEvent, func(request map[string]any) { queueRequestEntry(request, 0)["task_id"] = "ghost" })
	identity := func(event int) func(*WorkItemRegistry) {
		return func(r *WorkItemRegistry) { r.TaskQueueEvents[event].RecordedAtUTC = "yesterday" }
	}
	assertQueueViolations(t, []registryViolation{
		{"entry of the first event", func(r *WorkItemRegistry) { entry(r); identity(preparationEvent)(r) }, "foreign queue Task"},
		{"kind of the second event", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[preparationEvent].Kind = "bogus"
			identity(goalSetEvent)(r)
		}, "unknown queue event kind"},
		{"identity of the first event", func(r *WorkItemRegistry) {
			identity(firstSetEvent)(r)
			r.TaskQueueEvents[preparationEvent].Kind = "bogus"
		}, "invalid queue event"},
		{"event of the first queue before the second", func(r *WorkItemRegistry) {
			r.TaskQueueEvents[goalSetEvent].Kind = "bogus"
			identity(otherQueueEvent)(r)
		}, "unknown queue event kind"},
	})
}
