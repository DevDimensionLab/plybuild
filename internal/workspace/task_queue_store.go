package workspace

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"gopkg.in/yaml.v3"
	"reflect"
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
		"registry_upgrade": contentNullable(contentExact(map[string]contentRule{"from_version": contentInteger(1, 3), "registry_sha256": contentDigest})),
	})
	return func(v canonicaljson.Value) error {
		if contentInt(contentFields(v), "schema_version") == 2 {
			return goalQueueDraftRule()(v)
		}
		return legacy(v)
	}
}
func queueCloseRule(close bool) contentRule {
	m := map[string]contentRule{"preparation_id": queuePrefixedID("pre_"), "expected_revision": contentInteger(0, 2147483647)}
	if close {
		m["reason"] = contentText(2000)
	}
	return contentExact(m)
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
	sort.Slice(r.TaskPreparations, func(i, j int) bool { return r.TaskPreparations[i].ID < r.TaskPreparations[j].ID })
	sort.Slice(r.TaskQueueEvents, func(i, j int) bool {
		a, b := r.TaskQueueEvents[i], r.TaskQueueEvents[j]
		if a.QueueID != b.QueueID {
			return a.QueueID < b.QueueID
		}
		return a.Revision < b.Revision
	})
}

type queueFold struct {
	Revision int
	Pending  []TaskGoalQueueEntry
	Current  *string
	Terminal map[string]string
}

func emptyQueueFold() queueFold {
	return queueFold{Pending: []TaskGoalQueueEntry{}, Terminal: map[string]string{}}
}
func applyQueueEvent(q *queueFold, e TaskQueueEvent) {
	m := contentFields(queueRequestValue(e.Request))
	q.Revision = e.Revision
	switch e.Kind {
	case "set":
		d := decodeQueueDraft(queueRequestValue(e.Request))
		q.Pending = d.Entries
	case "reserve":
		id := contentString(m, "preparation_id")
		q.Current = &id
	case "advance", "release":
		id := contentString(m, "preparation_id")
		q.Current = nil
		q.Terminal[id] = map[string]string{"advance": "advanced", "release": "released"}[e.Kind]
	}
}
func foldQueue(r WorkItemRegistry, id string) queueFold {
	q := emptyQueueFold()
	for _, e := range r.TaskQueueEvents {
		if e.QueueID == id {
			applyQueueEvent(&q, e)
			if e.Kind == "reserve" {
				p := findPreparation(r, *q.Current)
				if p != nil {
					pending := []TaskGoalQueueEntry{}
					for _, v := range q.Pending {
						if v.TaskID != p.Plan.TaskID {
							pending = append(pending, v)
						}
					}
					q.Pending = pending
				}
			}
		}
	}
	return q
}
func findPreparation(r WorkItemRegistry, id string) *TaskPreparation {
	for _, p := range r.TaskPreparations {
		if p.ID == id {
			return &p
		}
	}
	return nil
}
func preparationForTask(r WorkItemRegistry, id TaskID) *TaskPreparation {
	for _, p := range r.TaskPreparations {
		if p.Plan.TaskID == id {
			return &p
		}
	}
	return nil
}
func preparationDisposition(r WorkItemRegistry, id string) string {
	p := findPreparation(r, id)
	if p == nil {
		return "unknown"
	}
	q := foldQueue(r, p.Plan.QueueID)
	if v := q.Terminal[id]; v != "" {
		return v
	}
	if q.Current != nil && *q.Current == id {
		return "current"
	}
	return "unknown"
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
	if !validStoredPath(w.Root) || !digestPattern.MatchString(w.MarkerSHA256) {
		return fmt.Errorf("invalid workspace binding")
	}
	if e := validateQueueTarget(r, t); e != nil {
		return e
	}
	if version < 1 || version > 4 {
		return fmt.Errorf("invalid plan registry version")
	}
	if version == 4 {
		if u != nil {
			return fmt.Errorf("format 4 upgrade must be null")
		}
	} else if u == nil || u.FromVersion != version || !digestPattern.MatchString(u.RegistrySHA256) {
		return fmt.Errorf("invalid upgrade binding")
	}
	return nil
}
func validateQueueRegistry(r WorkItemRegistry) error {
	if r.FormatVersion != 4 {
		if r.EpicBaseVersions != nil || r.EpicBaseUpdates != nil || r.TaskQueueEvents != nil || r.TaskPreparations != nil {
			return fmt.Errorf("legacy registry contains format 4 records")
		}
		return nil
	}
	if r.EpicBaseVersions == nil || r.EpicBaseUpdates == nil || r.TaskQueueEvents == nil || r.TaskPreparations == nil {
		return fmt.Errorf("format 4 requires explicit arrays")
	}
	copy := r
	copy.EpicBaseVersions = append([]EpicBaseVersion{}, r.EpicBaseVersions...)
	copy.EpicBaseUpdates = append([]EpicBaseUpdate{}, r.EpicBaseUpdates...)
	copy.TaskQueueEvents = append([]TaskQueueEvent{}, r.TaskQueueEvents...)
	copy.TaskPreparations = append([]TaskPreparation{}, r.TaskPreparations...)
	sortQueueRegistry(&copy)
	if !reflect.DeepEqual(copy.EpicBaseVersions, r.EpicBaseVersions) || !reflect.DeepEqual(copy.EpicBaseUpdates, r.EpicBaseUpdates) || !reflect.DeepEqual(copy.TaskQueueEvents, r.TaskQueueEvents) || !reflect.DeepEqual(copy.TaskPreparations, r.TaskPreparations) {
		return fmt.Errorf("format 4 arrays must be sorted")
	}
	updates := map[string]EpicBaseUpdate{}
	for _, u := range r.EpicBaseUpdates {
		p := u.Plan
		if u.ID != "ebu_"+strings.TrimPrefix(u.PlanSHA256, "sha256:") || u.PlanSHA256 != queueDigest(p) || updates[u.ID].ID != "" || contentUTC(u.RecordedAtUTC) != nil {
			return fmt.Errorf("invalid base update identity")
		}
		if e := validateQueuePlanEnvelope(r, p.Workspace, p.Target, p.RegistryFormatVersion, p.RegistryUpgrade); e != nil {
			return e
		}
		if p.Kind != "WorkspaceEpicBaseUpdatePlan@1" || p.SchemaVersion != 1 || p.ExpectedRevision < 1 || p.ExpectedRevision >= 2147483647 || !validOIDText(p.NextOID) || !validOIDText(p.NextTree) || p.NextOID == p.PreviousOID || p.AffectedSelections == nil || p.UnresolvedOperations == nil || len(p.UnresolvedOperations) != 0 || !p.ObservedParent.Clean || p.ObservedParent.OID != p.NextOID || p.ObservedParent.Tree != p.NextTree {
			return fmt.Errorf("invalid base update plan")
		}
		if err := validateQueuePlanObservation(p.ObservedParent, p.Target, p.NextOID, p.NextTree); err != nil {
			return err
		}
		lastAffected := ""
		for _, v := range p.AffectedSelections {
			task, _ := findTask(r, v.TaskID)
			if string(v.TaskID) <= lastAffected || task == nil || task.ProjectID != p.Target.ProjectID || task.RepoID != p.Target.RepoID || task.ParentEpicID != p.Target.EpicID {
				return fmt.Errorf("invalid affected selection target or order")
			}
			if err := validateQueueSelection(r, v.TaskID, v.Selection); err != nil {
				return err
			}
			lastAffected = string(v.TaskID)
		}
		if p.CurrentPreparationID != nil {
			prep := findPreparation(r, *p.CurrentPreparationID)
			if prep == nil || prep.Plan.Target != p.Target {
				return fmt.Errorf("invalid current preparation binding")
			}
		}
		switch u.Phase {
		case "intent":
			if u.Outcome != nil {
				return fmt.Errorf("intent has outcome")
			}
		case "committed":
			if u.Outcome == nil || u.Outcome.Revision == nil || *u.Outcome.Revision != p.ExpectedRevision+1 || len(u.Outcome.Reasons) != 0 {
				return fmt.Errorf("invalid committed base outcome")
			}
		case "conflict":
			if u.Outcome == nil || u.Outcome.Revision != nil || len(u.Outcome.Reasons) == 0 {
				return fmt.Errorf("invalid conflict outcome")
			}
		default:
			return fmt.Errorf("invalid base phase")
		}
		if u.Outcome != nil && u.Outcome.Reasons == nil {
			return fmt.Errorf("null reasons")
		}
		updates[u.ID] = u
	}
	versions := map[string][]EpicBaseVersion{}
	for _, v := range r.EpicBaseVersions {
		key := string(v.EpicID) + "/" + string(v.RepoID)
		chain := versions[key]
		ep, _ := findEpic(r, v.EpicID)
		if ep == nil {
			return fmt.Errorf("base has foreign Epic")
		}
		b, _ := findEpicRepo(*ep, v.RepoID)
		if b == nil {
			return fmt.Errorf("base has foreign repo")
		}
		adopted := adoptedEpicBase(*ep, *b)
		if v.Revision != len(chain)+1 || v.Revision > 2147483647 || v.ProjectID != adopted.ProjectID || v.WorktreeID != adopted.WorktreeID || v.Locator != adopted.Locator || v.Ref != adopted.Ref || v.GitCommonDir != adopted.GitCommonDir || !validOIDText(v.OID) || !validOIDText(v.Tree) {
			return fmt.Errorf("invalid base chain")
		}
		if len(chain) == 0 {
			if !contentTypedEqual(v, adopted) {
				return fmt.Errorf("base v1 differs from adoption")
			}
		} else {
			prev := chain[len(chain)-1]
			if v.PreviousRevision == nil || *v.PreviousRevision != prev.Revision || v.UpdateID == nil {
				return fmt.Errorf("broken base chain")
			}
			for _, old := range chain {
				if old.OID == v.OID {
					return fmt.Errorf("reused base OID")
				}
			}
			u := updates[*v.UpdateID]
			if u.Phase != "committed" || u.Plan.Target.EpicID != v.EpicID || u.Plan.Target.RepoID != v.RepoID || u.Plan.ExpectedRevision != prev.Revision || u.Plan.PreviousOID != prev.OID || u.Plan.PreviousTree != prev.Tree || u.Plan.NextOID != v.OID || u.Plan.NextTree != v.Tree {
				return fmt.Errorf("base differs from update")
			}
		}
		versions[key] = append(chain, v)
	}
	for _, ep := range r.Epics {
		for _, b := range ep.RepoBindings {
			if len(versions[string(ep.ID)+"/"+string(b.RepoID)]) == 0 {
				return fmt.Errorf("Epic lacks base v1")
			}
		}
	}
	for _, u := range r.EpicBaseUpdates {
		p := u.Plan
		if !historicalEpicBaseMatches(r, p.Target.EpicID, p.Target.RepoID, p.Target.ParentWorktreeID, p.Target.ParentRef, p.PreviousOID, p.PreviousTree) {
			return fmt.Errorf("base request parent missing")
		}
		if u.Phase == "committed" {
			found := false
			for _, v := range r.EpicBaseVersions {
				if v.UpdateID != nil && *v.UpdateID == u.ID {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("committed update lacks version")
			}
		}
	}
	seen := map[string]bool{}
	tasks := map[TaskID]bool{}
	for _, prep := range r.TaskPreparations {
		p := prep.Plan
		if prep.ID != "pre_"+strings.TrimPrefix(prep.PlanSHA256, "sha256:") || prep.PlanSHA256 != queueDigest(p) || seen[prep.ID] || tasks[p.TaskID] || contentUTC(prep.CreatedAtUTC) != nil {
			return fmt.Errorf("invalid preparation identity")
		}
		seen[prep.ID] = true
		tasks[p.TaskID] = true
		if e := validateQueuePlanEnvelope(r, p.Workspace, p.Target, p.RegistryFormatVersion, p.RegistryUpgrade); e != nil {
			return e
		}
		t, _ := findTask(r, p.TaskID)
		if t == nil || t.ProjectID != p.Target.ProjectID || t.RepoID != p.Target.RepoID || t.ParentEpicID != p.Target.EpicID {
			return fmt.Errorf("foreign preparation Task")
		}
		if p.Kind != "WorkspaceTaskPreparePlan@1" || p.SchemaVersion != 1 || p.QueueID != queueID(p.Workspace.Root, p.Target) || p.QueueRevision < 0 || p.QueueRevision >= 2147483647 || p.RequiredInputs == nil || !validShortBranchText(p.Branch) || !validStoredPath(p.WorktreePath) || !historicalEpicBaseMatches(r, p.Target.EpicID, p.Target.RepoID, p.Target.ParentWorktreeID, p.Target.ParentRef, p.ParentOID, p.ParentTree) {
			return fmt.Errorf("invalid prepare plan")
		}
		if err := validateQueuePlanObservation(p.ObservedParent, p.Target, p.ParentOID, p.ParentTree); err != nil {
			return err
		}
		baseFound := false
		for _, v := range r.EpicBaseVersions {
			if v.EpicID == p.Target.EpicID && v.RepoID == p.Target.RepoID && v.OID == p.ParentOID && v.Tree == p.ParentTree && v.Revision == p.BaseRevision {
				baseFound = true
			}
		}
		if !baseFound {
			return fmt.Errorf("preparation base revision differs")
		}
		if (p.Selector.Kind == "next" && p.Selector.TaskID != nil) || (p.Selector.Kind == "task" && (p.Selector.TaskID == nil || *p.Selector.TaskID != p.TaskID)) || (p.Selector.Kind != "next" && p.Selector.Kind != "task") {
			return fmt.Errorf("invalid selector")
		}
		if e := validateQueueSelection(r, p.TaskID, &p.Selection); e != nil {
			return e
		}
		op, _ := findOperation(r, p.TaskID)
		if op == nil || op.ID != prep.OperationID || op.ParentOID != p.ParentOID || op.ParentTree != p.ParentTree || op.SourceRef != "refs/heads/"+p.Branch || op.TargetLocator != p.WorktreePath {
			return fmt.Errorf("preparation operation differs")
		}
		if prep.Outcome != nil {
			v := prep.Outcome
			if v.Kind != "prepared" || v.WorktreeID != op.WorktreeID || v.ObservedOID != p.ParentOID || v.ObservedTree != p.ParentTree || op.State != "ready" || contentUTC(v.RecordedAtUTC) != nil {
				return fmt.Errorf("invalid prepared outcome")
			}
		}
	}
	folds := map[string]queueFold{}
	keys := map[string]bool{}
	reserved := map[string]bool{}
	for _, e := range r.TaskQueueEvents {
		if queuePrefixedID("que_")(e.QueueID) != nil || contentUTC(e.RecordedAtUTC) != nil || e.RequestSHA256 != queueDigest(e.Request) {
			return fmt.Errorf("invalid queue event")
		}
		q, ok := folds[e.QueueID]
		if !ok {
			q = emptyQueueFold()
		}
		v := queueRequestValue(e.Request)
		m := contentFields(v)
		if q.Revision == 2147483647 || e.Revision != q.Revision+1 || contentInt(m, "expected_revision") != q.Revision {
			return fmt.Errorf("broken queue CAS chain")
		}
		switch e.Kind {
		case "set":
			if err := queueDraftRule()(v); err != nil {
				return err
			}
			key := contentString(m, "publication_key")
			if keys[key] {
				return fmt.Errorf("duplicate queue publication key")
			}
			keys[key] = true
			draft := decodeQueueDraft(v)
			ep, _ := findEpic(r, draft.EpicID)
			if ep == nil || ep.ProjectID != draft.ProjectID {
				return fmt.Errorf("foreign queue target")
			}
			b, _ := findEpicRepo(*ep, draft.RepoID)
			if b == nil {
				return fmt.Errorf("foreign queue repo")
			}
			for _, entry := range draft.Entries {
				t, _ := findTask(r, entry.TaskID)
				if t == nil || t.ProjectID != draft.ProjectID || t.RepoID != draft.RepoID || t.ParentEpicID != draft.EpicID {
					return fmt.Errorf("foreign queue Task")
				}
				if err := validateQueueSelection(r, t.ID, entry.Selection); err != nil {
					return err
				}
				if err := validateRegisteredGoalRef(r, t.ID, entry.Goal); err != nil {
					return err
				}
				if q.Current != nil && findPreparation(r, *q.Current).Plan.TaskID == t.ID {
					return fmt.Errorf("set includes current Task")
				}
				for id := range q.Terminal {
					if findPreparation(r, id).Plan.TaskID == t.ID {
						return fmt.Errorf("terminal Task requeued")
					}
				}
			}
		case "reserve":
			if err := queueCloseRule(false)(v); err != nil {
				return err
			}
			id := contentString(m, "preparation_id")
			p := findPreparation(r, id)
			if p == nil || reserved[id] || q.Current != nil || p.Plan.QueueID != e.QueueID || p.Plan.QueueRevision != q.Revision {
				return fmt.Errorf("invalid reservation")
			}
			found := false
			pending := []TaskGoalQueueEntry{}
			for _, entry := range q.Pending {
				if entry.TaskID == p.Plan.TaskID && contentTypedEqual(entry.Selection, &p.Plan.Selection) {
					found = true
				} else {
					pending = append(pending, entry)
				}
			}
			if !found {
				return fmt.Errorf("reservation is not a pending entry")
			}
			q.Pending = pending
			reserved[id] = true
		case "advance", "release":
			if err := queueCloseRule(true)(v); err != nil {
				return err
			}
			id := contentString(m, "preparation_id")
			if q.Current == nil || *q.Current != id {
				return fmt.Errorf("closing a noncurrent reservation")
			}
			if e.Kind == "advance" && findPreparation(r, id).Outcome == nil {
				return fmt.Errorf("advance needs prepared outcome")
			}
		default:
			return fmt.Errorf("unknown queue event kind")
		}
		applyQueueEvent(&q, e)
		folds[e.QueueID] = q
	}
	for _, p := range r.TaskPreparations {
		if !reserved[p.ID] {
			return fmt.Errorf("preparation lacks reservation")
		}
	}
	return nil
}

// Registry publication enforces append-only history, independently of public commands.
func validateQueueTransition(old, next WorkItemRegistry) error {
	if old.FormatVersion != 4 {
		return nil
	}
	if next.FormatVersion != 4 {
		return fmt.Errorf("format 4 cannot be downgraded")
	}
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
	for _, v := range old.TaskPreparations {
		n := findPreparation(next, v.ID)
		if n == nil || v.PlanSHA256 != n.PlanSHA256 || !contentTypedEqual(v.Plan, n.Plan) || v.OperationID != n.OperationID || v.CreatedAtUTC != n.CreatedAtUTC || v.Outcome != nil && !contentTypedEqual(v.Outcome, n.Outcome) {
			return fmt.Errorf("immutable preparation was removed or changed")
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
