package taskwork

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
	"reflect"
	"slices"
	"sort"
	"strings"
)

func validateQueueYAMLMap(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return fmt.Errorf("queue request must be a mapping")
	}
	var walk func(*yaml.Node) error
	walk = func(n *yaml.Node) error {
		if n.Kind == yaml.AliasNode {
			return fmt.Errorf("aliases are unsupported")
		}
		if n.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for i := 0; i < len(n.Content); i += 2 {
				k := n.Content[i]
				if k.Tag != "!!str" || seen[k.Value] {
					return fmt.Errorf("duplicate or non-string request key")
				}
				seen[k.Value] = true
				if e := walk(n.Content[i+1]); e != nil {
					return e
				}
			}
		} else if n.Kind == yaml.SequenceNode {
			for _, v := range n.Content {
				if e := walk(v); e != nil {
					return e
				}
			}
		}
		return nil
	}
	return walk(n)
}
func queueRequest(v any) map[string]any {
	var b []byte
	if object, ok := v.(canonicaljson.Object); ok {
		b, _ = canonicaljson.Marshal(object)
	} else {
		b, _ = contentCanonical(v)
	}
	var out map[string]any
	_ = yaml.Unmarshal(b, &out)
	return out
}
func queueRequestValue(m map[string]any) canonicaljson.Value { v, _ := contentValue(m); return v }
func queueDraftRule() contentRule {
	legacy := contentExact(map[string]contentRule{
		"kind": contentEnum("WorkspaceTaskQueueDraft@1"), "schema_version": contentInteger(1, 1), "publication_key": contentKey, "project_id": contentSlug, "repo_id": contentSlug, "epic_id": contentSlug, "expected_revision": contentInteger(0, 2147483647),
		"entries":          contentList(256, contentExact(map[string]contentRule{"task_id": contentSlug, "selection": contentNullable(contentDecisionRule("sel_"))}), func(v canonicaljson.Value) string { return contentString(contentFields(v), "task_id") }, true),
		"human_decision":   contentExact(map[string]contentRule{"actor_claim": contentText(256), "decided_at_utc": contentUTC, "source": contentEnum("human_cli", "explicit_human_instruction"), "statement": contentText(2000)}),
		"registry_upgrade": contentNoRegistryUpgrade,
	})
	return func(v canonicaljson.Value) error {
		if contentInt(contentFields(v), "schema_version") == 2 {
			return goalQueueDraftRule()(v)
		}
		return legacy(v)
	}
}
func queuePrefixedID(prefix string) contentRule {
	return func(v canonicaljson.Value) error {
		s, ok := v.(string)
		if !ok || !strings.HasPrefix(s, prefix) || !digestPattern.MatchString("sha256:"+strings.TrimPrefix(s, prefix)) {
			return fmt.Errorf("invalid %s ID", prefix)
		}
		return nil
	}
}
func sortQueueRegistry(r *WorkItemRegistry) {
	sort.Slice(r.EpicBaseVersions, func(i, j int) bool {
		a, b := r.EpicBaseVersions[i], r.EpicBaseVersions[j]
		if a.EpicID != b.EpicID {
			return a.EpicID < b.EpicID
		}
		if a.RepoID != b.RepoID {
			return a.RepoID < b.RepoID
		}
		return a.Revision < b.Revision
	})
	sort.Slice(r.EpicBaseUpdates, func(i, j int) bool { return r.EpicBaseUpdates[i].ID < r.EpicBaseUpdates[j].ID })
	sort.Slice(r.TaskQueueEvents, func(i, j int) bool {
		a, b := r.TaskQueueEvents[i], r.TaskQueueEvents[j]
		if a.QueueID != b.QueueID {
			return a.QueueID < b.QueueID
		}
		return a.Revision < b.Revision
	})
}

// Only "set" orders the queue. Events of the removed preparation keep their
// place in the append-only chain without meaning.
type queueFold struct {
	Revision int
	Pending  []TaskGoalQueueEntry
}

func emptyQueueFold() queueFold {
	return queueFold{Pending: []TaskGoalQueueEntry{}}
}
func applyQueueEvent(q *queueFold, e TaskQueueEvent) {
	q.Revision = e.Revision
	if e.Kind == "set" {
		q.Pending = decodeQueueDraft(queueRequestValue(e.Request)).Entries
	}
}
func foldQueue(r WorkItemRegistry, id string) queueFold {
	q := emptyQueueFold()
	for _, e := range r.TaskQueueEvents {
		if e.QueueID == id {
			applyQueueEvent(&q, e)
		}
	}
	return q
}
func validateQueueTarget(r WorkItemRegistry, t QueueTarget) error {
	e, _ := findEpic(r, t.EpicID)
	if e == nil || e.ProjectID != t.ProjectID {
		return fmt.Errorf("foreign queue Epic")
	}
	b, _ := findEpicRepo(*e, t.RepoID)
	if b == nil || b.GitCommonDir != t.GitCommonDir || b.Worktree.ID != t.ParentWorktreeID || b.Worktree.Ref != t.ParentRef || b.Worktree.Locator != t.ParentLocator {
		return fmt.Errorf("foreign queue target")
	}
	return nil
}
func validateQueueSelection(r WorkItemRegistry, id TaskID, ref *TaskDecisionRef) error {
	if ref == nil {
		return nil
	}
	for _, v := range r.TaskSolutionSelections {
		if v.TaskID == id && v.ID == ref.ID && v.ManifestSHA256 == ref.ManifestSHA256 {
			return nil
		}
	}
	return fmt.Errorf("selection does not belong to Task %s", id)
}
func validateQueuePlanEnvelope(r WorkItemRegistry, w IntegrationPlanWorkspace, t QueueTarget, version int, u *TaskRegistryUpgrade) error {
	if !workspace.ValidStoredPath(w.Root) || !digestPattern.MatchString(w.MarkerSHA256) {
		return fmt.Errorf("invalid workspace binding")
	}
	if e := validateQueueTarget(r, t); e != nil {
		return e
	}
	if version != RegistryFormat || u != nil {
		return fmt.Errorf("plan must bind registry format %d without an upgrade", RegistryFormat)
	}
	return nil
}

// validateQueueRegistry runs the queue's phases in a fixed order; the first
// failing rule answers. Base updates, base versions and queue events are each
// validated in registry order, one record at a time.
func validateQueueRegistry(r WorkItemRegistry) error {
	if r.EpicBaseVersions == nil || r.EpicBaseUpdates == nil || r.TaskQueueEvents == nil || r.TaskPreparations == nil {
		return fmt.Errorf("queue arrays must be explicit")
	}
	if !queueArraysSorted(r) {
		return fmt.Errorf("queue arrays must be sorted")
	}
	validation := &queueRegistryValidation{registry: r, updates: map[string]EpicBaseUpdate{}, baseChains: map[string][]EpicBaseVersion{}, queueRevisions: map[string]int{}, publicationKeys: map[string]bool{}}
	if err := validation.validateBaseUpdates(); err != nil {
		return err
	}
	if err := validation.validateBaseVersions(); err != nil {
		return err
	}
	if err := validation.validateEpicsHaveBase(); err != nil {
		return err
	}
	if err := validation.validateBaseUpdatesAgainstBases(); err != nil {
		return err
	}
	return validation.validateQueueEvents()
}

// queueRegistryValidation carries what the records validated so far
// established for the ones after them: the base updates by ID, the base
// versions of each Epic binding in revision order, the last revision of each
// queue, and the publication keys in use.
type queueRegistryValidation struct {
	registry        WorkItemRegistry
	updates         map[string]EpicBaseUpdate
	baseChains      map[string][]EpicBaseVersion
	queueRevisions  map[string]int
	publicationKeys map[string]bool
}

func queueArraysSorted(r WorkItemRegistry) bool {
	sorted := r
	sorted.EpicBaseVersions = append([]EpicBaseVersion{}, r.EpicBaseVersions...)
	sorted.EpicBaseUpdates = append([]EpicBaseUpdate{}, r.EpicBaseUpdates...)
	sorted.TaskQueueEvents = append([]TaskQueueEvent{}, r.TaskQueueEvents...)
	sortQueueRegistry(&sorted)
	return reflect.DeepEqual(sorted.EpicBaseVersions, r.EpicBaseVersions) &&
		reflect.DeepEqual(sorted.EpicBaseUpdates, r.EpicBaseUpdates) &&
		reflect.DeepEqual(sorted.TaskQueueEvents, r.TaskQueueEvents)
}

// epicBaseKey names the chain of base versions of one Epic binding.
func epicBaseKey(epic EpicID, repo workspace.RepoID) string {
	return string(epic) + "/" + string(repo)
}

// taskOfTarget reports whether a registered Task belongs to the Project,
// repository and Epic that a base update or a queue targets.
func taskOfTarget(task *TaskRecord, project workspace.ProjectID, repo workspace.RepoID, epic EpicID) bool {
	return task != nil && task.ProjectID == project && task.RepoID == repo && task.ParentEpicID == epic
}

func (validation *queueRegistryValidation) validateBaseUpdates() error {
	for _, update := range validation.registry.EpicBaseUpdates {
		if err := validation.validateBaseUpdate(update); err != nil {
			return err
		}
		validation.updates[update.ID] = update
	}
	return nil
}

// validateBaseUpdate checks one update on its own; whether its bases exist is
// asked in validateBaseUpdatesAgainstBases, once the base versions are known.
func (validation *queueRegistryValidation) validateBaseUpdate(update EpicBaseUpdate) error {
	plan := update.Plan
	_, duplicate := validation.updates[update.ID]
	if update.ID != "ebu_"+strings.TrimPrefix(update.PlanSHA256, "sha256:") || update.PlanSHA256 != queueDigest(plan) || duplicate || contentUTC(update.RecordedAtUTC) != nil {
		return fmt.Errorf("invalid base update identity")
	}
	if err := validateQueuePlanEnvelope(validation.registry, plan.Workspace, plan.Target, plan.RegistryFormatVersion, plan.RegistryUpgrade); err != nil {
		return err
	}
	if !wellFormedBaseUpdatePlan(plan) {
		return fmt.Errorf("invalid base update plan")
	}
	if err := validateQueuePlanObservation(plan.ObservedParent, plan.Target, plan.NextOID, plan.NextTree); err != nil {
		return err
	}
	if err := validateAffectedSelections(validation.registry, plan); err != nil {
		return err
	}
	return validateBaseUpdateOutcome(update)
}

// A plan moves the base to a new commit that was observed clean, with nothing
// left unresolved. Only a historical plan may name a preparation.
func wellFormedBaseUpdatePlan(plan WorkspaceEpicBaseUpdatePlan) bool {
	historical := plan.Kind == "WorkspaceEpicBaseUpdatePlan@1" && plan.SchemaVersion == 1
	current := plan.Kind == "WorkspaceEpicBaseUpdatePlan@3" && plan.SchemaVersion == currentEpicBaseUpdatePlanSchema && plan.CurrentPreparationID == nil
	return (historical || current) &&
		plan.ExpectedRevision >= 1 && plan.ExpectedRevision < 2147483647 &&
		validOIDText(plan.NextOID) &&
		validOIDText(plan.NextTree) &&
		plan.NextOID != plan.PreviousOID &&
		plan.AffectedSelections != nil &&
		plan.UnresolvedOperations != nil && len(plan.UnresolvedOperations) == 0 &&
		plan.ObservedParent.Clean &&
		plan.ObservedParent.OID == plan.NextOID &&
		plan.ObservedParent.Tree == plan.NextTree
}

// Affected selections name Tasks in ascending order.
func validateAffectedSelections(r WorkItemRegistry, plan WorkspaceEpicBaseUpdatePlan) error {
	previous := TaskID("")
	for _, affected := range plan.AffectedSelections {
		if err := validateAffectedSelection(r, plan.Target, affected, previous); err != nil {
			return err
		}
		previous = affected.TaskID
	}
	return nil
}

// An affected Task belongs to the plan's target and comes after the previous
// one, with a selection of its own or none.
func validateAffectedSelection(r WorkItemRegistry, target QueueTarget, affected EpicAffectedSelection, previous TaskID) error {
	task, _ := findTask(r, affected.TaskID)
	if affected.TaskID <= previous || !taskOfTarget(task, target.ProjectID, target.RepoID, target.EpicID) {
		return fmt.Errorf("invalid affected selection target or order")
	}
	return validateQueueSelection(r, affected.TaskID, affected.Selection)
}

// validateBaseUpdateOutcome answers with the first problem that applies: an
// intent has no outcome yet, a committed update has the revision it produced
// and no reasons, and a conflict has reasons and no revision.
func validateBaseUpdateOutcome(update EpicBaseUpdate) error {
	outcome, phase := update.Outcome, update.Phase
	switch {
	case phase == "intent" && outcome != nil:
		return fmt.Errorf("intent has outcome")
	case phase == "committed" && (outcome == nil || outcome.Revision == nil || *outcome.Revision != update.Plan.ExpectedRevision+1 || len(outcome.Reasons) != 0):
		return fmt.Errorf("invalid committed base outcome")
	case phase == "conflict" && (outcome == nil || outcome.Revision != nil || len(outcome.Reasons) == 0):
		return fmt.Errorf("invalid conflict outcome")
	case phase != "intent" && phase != "committed" && phase != "conflict":
		return fmt.Errorf("invalid base phase")
	case outcome != nil && outcome.Reasons == nil:
		return fmt.Errorf("null reasons")
	}
	return nil
}

func (validation *queueRegistryValidation) validateBaseVersions() error {
	for _, version := range validation.registry.EpicBaseVersions {
		key := epicBaseKey(version.EpicID, version.RepoID)
		chain := validation.baseChains[key]
		if err := validation.validateBaseVersion(version, chain); err != nil {
			return err
		}
		validation.baseChains[key] = append(chain, version)
	}
	return nil
}

// validateBaseVersion checks one base version against its Epic binding and
// against chain, the earlier versions of that binding.
func (validation *queueRegistryValidation) validateBaseVersion(version EpicBaseVersion, chain []EpicBaseVersion) error {
	epic, _ := findEpic(validation.registry, version.EpicID)
	if epic == nil {
		return fmt.Errorf("base has foreign Epic")
	}
	binding, _ := findEpicRepo(*epic, version.RepoID)
	if binding == nil {
		return fmt.Errorf("base has foreign repo")
	}
	adopted := adoptedEpicBase(*epic, *binding)
	if version.Revision != len(chain)+1 || version.Revision > 2147483647 || !baseOnAdoptedWorktree(version, adopted) || !validOIDText(version.OID) || !validOIDText(version.Tree) {
		return fmt.Errorf("invalid base chain")
	}
	if len(chain) > 0 {
		return validation.validateBaseSuccessor(version, chain)
	}
	if !contentTypedEqual(version, adopted) {
		return fmt.Errorf("base v1 differs from adoption")
	}
	return nil
}

// Every base version stays on the worktree that its Epic binding adopted.
func baseOnAdoptedWorktree(version, adopted EpicBaseVersion) bool {
	return version.ProjectID == adopted.ProjectID &&
		version.WorktreeID == adopted.WorktreeID &&
		version.Locator == adopted.Locator &&
		version.Ref == adopted.Ref &&
		version.GitCommonDir == adopted.GitCommonDir
}

// A later base version continues its chain from the previous one, on a commit
// that was never the base before, through the committed update that moved it.
func (validation *queueRegistryValidation) validateBaseSuccessor(version EpicBaseVersion, chain []EpicBaseVersion) error {
	previous := chain[len(chain)-1]
	if version.PreviousRevision == nil || *version.PreviousRevision != previous.Revision || version.UpdateID == nil {
		return fmt.Errorf("broken base chain")
	}
	if slices.ContainsFunc(chain, func(earlier EpicBaseVersion) bool { return earlier.OID == version.OID }) {
		return fmt.Errorf("reused base OID")
	}
	// An unknown update ID reads as the zero update, which is not committed.
	update := validation.updates[*version.UpdateID]
	if update.Phase != "committed" || !planMovesBase(update.Plan, previous, version) {
		return fmt.Errorf("base differs from update")
	}
	return nil
}

// planMovesBase reports whether the plan moved the base of next's Epic
// binding exactly from previous to next.
func planMovesBase(plan WorkspaceEpicBaseUpdatePlan, previous, next EpicBaseVersion) bool {
	return plan.Target.EpicID == next.EpicID &&
		plan.Target.RepoID == next.RepoID &&
		plan.ExpectedRevision == previous.Revision &&
		plan.PreviousOID == previous.OID &&
		plan.PreviousTree == previous.Tree &&
		plan.NextOID == next.OID &&
		plan.NextTree == next.Tree
}

func (validation *queueRegistryValidation) validateEpicsHaveBase() error {
	for _, epic := range validation.registry.Epics {
		for _, binding := range epic.RepoBindings {
			if len(validation.baseChains[epicBaseKey(epic.ID, binding.RepoID)]) == 0 {
				return fmt.Errorf("Epic lacks base v1")
			}
		}
	}
	return nil
}

// validateBaseUpdatesAgainstBases runs once every base version is known.
func (validation *queueRegistryValidation) validateBaseUpdatesAgainstBases() error {
	for _, update := range validation.registry.EpicBaseUpdates {
		if err := validateBaseUpdateAgainstBases(validation.registry, update); err != nil {
			return err
		}
	}
	return nil
}

// An update started from a registered base of its target, and a committed
// update produced a base version.
func validateBaseUpdateAgainstBases(r WorkItemRegistry, update EpicBaseUpdate) error {
	plan := update.Plan
	if !historicalEpicBaseMatches(r, plan.Target.EpicID, plan.Target.RepoID, plan.Target.ParentWorktreeID, plan.Target.ParentRef, plan.PreviousOID, plan.PreviousTree) {
		return fmt.Errorf("base request parent missing")
	}
	recordsUpdate := func(version EpicBaseVersion) bool { return version.UpdateID != nil && *version.UpdateID == update.ID }
	if update.Phase == "committed" && !slices.ContainsFunc(r.EpicBaseVersions, recordsUpdate) {
		return fmt.Errorf("committed update lacks version")
	}
	return nil
}

func (validation *queueRegistryValidation) validateQueueEvents() error {
	for _, event := range validation.registry.TaskQueueEvents {
		if err := validation.validateQueueEvent(event); err != nil {
			return err
		}
		validation.queueRevisions[event.QueueID] = event.Revision
	}
	return nil
}

// validateQueueEvent checks one event and its place in its queue's chain: each
// event has the next revision and expects the one before it, from 0.
func (validation *queueRegistryValidation) validateQueueEvent(event TaskQueueEvent) error {
	if queuePrefixedID("que_")(event.QueueID) != nil || contentUTC(event.RecordedAtUTC) != nil || event.RequestSHA256 != queueDigest(event.Request) {
		return fmt.Errorf("invalid queue event")
	}
	request := queueRequestValue(event.Request)
	last := validation.queueRevisions[event.QueueID]
	if last == 2147483647 || event.Revision != last+1 || contentInt(contentFields(request), "expected_revision") != last {
		return fmt.Errorf("broken queue CAS chain")
	}
	switch event.Kind {
	case "set":
		return validation.validateQueueSet(request)
	case "reserve", "reserve_parallel", "advance", "release":
		// Preparation history is kept in the chain without meaning.
		return nil
	default:
		return fmt.Errorf("unknown queue event kind")
	}
}

// validateQueueSet checks the request of the one event kind that orders the
// queue: the draft, its publication key, its target and its entries.
func (validation *queueRegistryValidation) validateQueueSet(request canonicaljson.Value) error {
	if err := queueDraftRule()(request); err != nil {
		return err
	}
	key := contentString(contentFields(request), "publication_key")
	if validation.publicationKeys[key] {
		return fmt.Errorf("duplicate queue publication key")
	}
	validation.publicationKeys[key] = true
	r := validation.registry
	draft := decodeQueueDraft(request)
	epic, _ := findEpic(r, draft.EpicID)
	if epic == nil || epic.ProjectID != draft.ProjectID {
		return fmt.Errorf("foreign queue target")
	}
	if binding, _ := findEpicRepo(*epic, draft.RepoID); binding == nil {
		return fmt.Errorf("foreign queue repo")
	}
	for _, entry := range draft.Entries {
		if err := validateQueueEntry(r, draft, entry); err != nil {
			return err
		}
	}
	return nil
}

// A queued Task belongs to the draft's target, with a selection or a goal that
// is registered to that Task.
func validateQueueEntry(r WorkItemRegistry, draft decodedTaskQueueDraft, entry TaskGoalQueueEntry) error {
	task, _ := findTask(r, entry.TaskID)
	if !taskOfTarget(task, draft.ProjectID, draft.RepoID, draft.EpicID) {
		return fmt.Errorf("foreign queue Task")
	}
	if err := validateQueueSelection(r, task.ID, entry.Selection); err != nil {
		return err
	}
	return validateRegisteredGoalRef(r, task.ID, entry.Goal)
}

// Registry publication enforces append-only history, independently of public commands.
func validateQueueTransition(old, next WorkItemRegistry) error {
	for _, v := range old.EpicBaseVersions {
		found := false
		for _, n := range next.EpicBaseVersions {
			if v.EpicID == n.EpicID && v.RepoID == n.RepoID && v.Revision == n.Revision {
				found = contentTypedEqual(v, n)
			}
		}
		if !found {
			return fmt.Errorf("immutable Epic base was removed or changed")
		}
	}
	for _, v := range old.TaskQueueEvents {
		found := false
		for _, n := range next.TaskQueueEvents {
			if v.QueueID == n.QueueID && v.Revision == n.Revision {
				found = contentTypedEqual(v, n)
			}
		}
		if !found {
			return fmt.Errorf("immutable queue event was removed or changed")
		}
	}
	for _, v := range old.EpicBaseUpdates {
		found := false
		for _, n := range next.EpicBaseUpdates {
			if v.ID != n.ID {
				continue
			}
			found = v.PlanSHA256 == n.PlanSHA256 && contentTypedEqual(v.Plan, n.Plan) && v.RecordedAtUTC == n.RecordedAtUTC
			if v.Phase != "intent" {
				found = found && contentTypedEqual(v, n)
			}
		}
		if !found {
			return fmt.Errorf("immutable base request was removed or changed")
		}
	}
	return nil
}
func validateQueuePlanObservation(p PlanWorktreeObservation, t QueueTarget, oid, tree string) error {
	if p.WorktreeLocator != t.ParentLocator || p.Ref != t.ParentRef || p.GitCommonDir != t.GitCommonDir || p.OID != oid || p.Tree != tree || !p.Clean || !p.Symbolic || p.StatusEntries == nil || len(p.StatusEntries) != 0 || p.InProgress == nil || len(p.InProgress) != 0 || p.ObjectFormat != "sha1" && p.ObjectFormat != "sha256" || p.RefFormat != "files" {
		return fmt.Errorf("plan does not bind an exact clean Epic observation")
	}
	return nil
}
