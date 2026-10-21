package workspace

import (
	"fmt"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// BackupGoalRegistration is supplied only after the workflow layer validates
// the native retired owner. Content publication and append share the existing
// registry lock; each durable intermediate publication has an exact retry key.
type BackupGoalRegistration struct {
	SourcePreparation string           `json:"source_preparation"`
	SourceTask        TaskID           `json:"source_task"`
	SourceGoal        TaskGoalRef      `json:"source_goal"`
	Target            QueueTargetInput `json:"target"`
	TaskID            TaskID           `json:"task_id"`
	SpecID            string           `json:"spec_id"`
	Provider          string           `json:"provider"`
	Model             string           `json:"model"`
	Effort            string           `json:"effort"`
	ExecutorSource    string           `json:"executor_source"`
	Provenance        string           `json:"provenance"`
	RecordedAtUTC     string           `json:"recorded_at_utc"`
	PublicationKey    string           `json:"publication_key"`
}
type BackupGoalRegistrationResult struct {
	Kind           string                     `json:"kind"`
	SchemaVersion  int                        `json:"schema_version"`
	Provenance     string                     `json:"provenance"`
	Effects        []string                   `json:"effects"`
	TaskID         TaskID                     `json:"task_id"`
	Goal           TaskGoalRef                `json:"goal"`
	Requirements   []TaskGoalRequirement      `json:"requirements"`
	Constraints    []string                   `json:"constraints"`
	Target         QueueTarget                `json:"target"`
	Delivery       *DeliveryAgreement         `json:"delivery"`
	Queue          WorkspaceTaskQueueReadback `json:"queue"`
	Confirmation   string                     `json:"confirmation"`
	Applied        bool                       `json:"applied"`
	Executor       TaskExecutorAssignment     `json:"executor"`
	ExecutorSource string                     `json:"executor_source"`
	NextArgv       []string                   `json:"next_argv"`
}

func RegisterBackupGoal(d Dependencies, in BackupGoalRegistration, apply bool, confirmation string) (out BackupGoalRegistrationResult, err error) {
	if contentSlug(string(in.TaskID)) != nil || contentSlug(in.SpecID) != nil || in.TaskID == in.SourceTask || (in.Provider != "codex" && in.Provider != "claude") {
		return out, WorkInvalidArguments("explicit fresh Task/Spec IDs and codex or claude provider are required")
	}
	executor := TaskExecutorAssignment{Provider: &in.Provider, Model: &in.Model, Effort: &in.Effort}
	if err = ValidateTaskExecutorChoice(in.Provider, in.Model, in.Effort); err != nil {
		return out, err
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	run := func(projects ProjectSnapshot, session WorkItemStoreSession, r WorkItemRegistry) error {
		source, err := loadTaskGoal(d, root, r, in.SourceTask, in.SourceGoal)
		if err != nil {
			return err
		}
		target, _, err := resolveQueueTarget(d, root, r, projects, in.Target)
		if err != nil {
			return err
		}
		if foldQueue(r, queueID(root, target)).Terminal[in.SourcePreparation] != "released" {
			return queueError("backup_goal_source_unresolved", "retired source preparation has not been natively released")
		}
		a := source.Delivery
		if a == nil || a.Mode != DeliveryLocalEpic || !a.AutomaticAcceptance() || a.TargetWorktree != target.ParentLocator || a.TargetRef != target.ParentRef || a.ProjectID != target.ProjectID || a.RepoID != target.RepoID || a.EpicID != target.EpicID {
			return queueError("backup_goal_target_conflict", "source must bind the same automatic local Epic target")
		}
		// A single run-scoped key fences changed fresh IDs/provider as well as retries.
		token := queueDigest(in)
		description := "Fresh backup registration " + token + ". " + in.Provenance
		task := TaskRecord{ID: in.TaskID, Title: source.Title, Description: description, ParentEpicID: target.EpicID, ProjectID: target.ProjectID, RepoID: target.RepoID, GitCommonDir: target.GitCommonDir, WorktreeState: WorkItemUnbound}
		existing, _ := findTask(r, in.TaskID)
		createKey := "task-create/" + string(in.TaskID)
		if existing != nil && (!contentTypedEqual(*existing, task) || taskPublication(r, createKey) == nil) {
			return queueError("backup_goal_identity_conflict", "fresh Task ID already belongs to different input")
		}
		for _, pub := range r.TaskContentPublications {
			if pub.PublicationKey == in.PublicationKey+"/problem" && pub.TaskID != in.TaskID {
				return queueError("backup_goal_publication_conflict", "source registration key already identifies another fresh Task")
			}
		}
		for _, t := range r.Tasks {
			if t.ID != in.TaskID && strings.HasPrefix(t.Description, "Fresh backup registration ") && strings.HasSuffix(t.Description, ". "+in.Provenance) {
				return queueError("backup_goal_publication_conflict", "source already identifies a different fresh Task")
			}
		}
		for _, ref := range r.TaskSpecRevisions {
			if ref.SpecID == in.SpecID && ref.TaskID != in.TaskID {
				return queueError("backup_goal_identity_conflict", "fresh Spec ID is already registered")
			}
		}
		normative := fmt.Sprintf("Fresh runtime binding (normative for runtime identity over preserved source prose): Task %s, Spec %s, provider %s; model %s, effort %s. Defaults/overrides: %s. %s No old acceptance, permission, review, tests, QA, result, candidate, qualification or delivery is inherited. Product requirements and documents remain binding. No source patch transfer or provider start.", in.TaskID, in.SpecID, in.Provider, in.Model, in.Effort, in.ExecutorSource, in.Provenance)
		constraints := append(append([]string{}, source.Constraints...), normative)
		sort.Strings(constraints)
		out = BackupGoalRegistrationResult{Kind: "ply.workflow.backup-goal", SchemaVersion: 1, Provenance: in.Provenance, Effects: []string{"publish fresh Task/Problem/Goal", "append exact Goal to existing queue", "no branch/worktree or provider start", "no inherited runtime acceptance or result"}, ExecutorSource: in.ExecutorSource, Executor: executor, TaskID: in.TaskID, Requirements: source.Requirements, Constraints: constraints, Target: target, Delivery: a, Confirmation: token, NextArgv: []string{"ply", "workflow", "execute", "--spec", in.SpecID, "--check"}}
		out.Queue, err = buildQueueReadback(d, root, r, projects, in.Target, false)
		if err != nil {
			return err
		}
		// Generate all immutable drafts deterministically, including after interruption.
		sourceSpec, err := registeredTaskSpec(d, root, r, in.SourceTask, in.SourceGoal.SpecID, &in.SourceGoal.Spec)
		if err != nil {
			return err
		}
		problem, err := readRegisteredTaskManifest(d.TaskContent, root, r, source.Problem.ManifestSHA256)
		if err != nil {
			return err
		}
		recorder := contentRefValue(TaskGoalQueueRecorder{ActorClaim: "backup goal operator", ControlSurface: "ply workflow execute backup-goal", RecordedAtUTC: in.RecordedAtUTC})
		draft := func(v canonicaljson.Object, kind, key, hash string) canonicaljson.Object {
			m := contentFields(v)
			for _, k := range []string{"revision", "previous", "source_draft_sha256", "recorded_at_utc"} {
				delete(m, k)
			}
			m["kind"], m["publication_key"], m["task_id"], m["registry_upgrade"], m["expected_previous"], m["recorder"] = kind, key, string(in.TaskID), nil, nil, recorder
			m["change_reason"] = in.Provenance
			docs := []canonicaljson.Value{}
			for _, doc := range contentArray(m, "documents") {
				id := contentString(contentFields(doc), "id")
				docs = append(docs, contentRefValue(map[string]any{"id": id, "source": map[string]any{"kind": "snapshot", "task_id": in.SourceTask, "manifest_sha256": hash, "document_id": id}}))
			}
			m["documents"] = docs
			return contentObject(m)
		}
		pd := contentFields(draft(problem, "WorkspaceTaskProblemDraft@1", in.PublicationKey+"/problem", source.Problem.ManifestSHA256))
		pd["origin"] = contentRefValue(map[string]any{"kind": "authored"})
		// The initial generated problem is a native Task identity receipt; revision 2
		// preserves the original problem documents/claims without accepting history.
		if existing != nil {
			p := taskPublication(r, createKey)
			pd["expected_previous"] = contentRefValue(TaskRevisionRef{Revision: *p.OutcomeRef.Revision, ManifestSHA256: p.OutcomeRef.ManifestSHA256})
		}
		// Validate both planned payloads before the first Task write.
		sdCheck := contentFields(draft(sourceSpec, "WorkspaceTaskSpecDraft@2", in.PublicationKey+"/goal", in.SourceGoal.Spec.ManifestSHA256))
		sdCheck["spec_id"], sdCheck["problem"], sdCheck["constraints"], sdCheck["executor"] = in.SpecID, contentRefValue(source.Problem), contentRefValue(constraints), contentRefValue(executor)
		for _, planned := range []struct {
			value canonicaljson.Object
			op    string
		}{{contentObject(pd), "problem_record"}, {contentObject(sdCheck), "spec_record"}} {
			raw, e := canonicaljson.Marshal(planned.value)
			if e != nil {
				return e
			}
			if _, e = decodeTaskContent(raw, planned.op, false, false); e != nil {
				return e
			}
		}
		queueAlreadyPublished := false
		for _, event := range r.TaskQueueEvents {
			if event.Kind == "set" && contentString(contentFields(queueRequestValue(event.Request)), "publication_key") == in.PublicationKey+"/queue" {
				queueAlreadyPublished = true
			}
		}
		if !queueAlreadyPublished && len(foldQueue(r, queueID(root, target)).Pending) >= 256 {
			return queueError("backup_goal_queue_full", "queue append exceeds 256 entries")
		}
		if !apply {
			out.Goal = TaskGoalRef{SpecID: in.SpecID, Spec: TaskRevisionRef{Revision: 1}}
			if pub := taskPublication(r, in.PublicationKey+"/goal"); pub != nil {
				out.Goal.Spec = TaskRevisionRef{Revision: *pub.OutcomeRef.Revision, ManifestSHA256: pub.OutcomeRef.ManifestSHA256}
			}
			return nil
		}
		if confirmation != token {
			return queueError("backup_goal_confirmation_conflict", "preview confirmation differs from the fresh binding")
		}
		if existing == nil {
			if _, err = createTaskProblem(d, root, session, r, projects, task, nil); err != nil {
				return fmt.Errorf("backup Task creation requires publication readback; inspect Task %s key %s before identical retry: %w", in.TaskID, createKey, err)
			}
			r, err = session.Snapshot()
			if err != nil {
				return err
			}
			p := taskPublication(r, createKey)
			pd["expected_previous"] = contentRefValue(TaskRevisionRef{Revision: *p.OutcomeRef.Revision, ManifestSHA256: p.OutcomeRef.ManifestSHA256})
		}
		publish := func(v canonicaljson.Object, op string) (*TaskContentOutcomeRef, error) {
			raw, e := canonicaljson.Marshal(v)
			if e != nil {
				return nil, e
			}
			if _, e = decodeTaskContent(raw, op, false, false); e != nil {
				return nil, e
			}
			result := newTaskContentMutation(root, in.TaskID, contentString(contentFields(v), "publication_key"), digestTaskBytes(raw))
			e = publishTaskContent(d, root, session, r, task, projects, v, op, raw, &result, false)
			if e != nil {
				return nil, fmt.Errorf("backup registration may be partially published; inspect Task %s publications and queue, then repeat identical input: %w", in.TaskID, e)
			}
			r, e = session.Snapshot()
			return result.OutcomeRef, e
		}
		p, err := publish(contentObject(pd), "problem_record")
		if err != nil {
			return err
		}
		sd := contentFields(draft(sourceSpec, "WorkspaceTaskSpecDraft@2", in.PublicationKey+"/goal", in.SourceGoal.Spec.ManifestSHA256))
		sd["spec_id"], sd["problem"], sd["constraints"] = in.SpecID, contentRefValue(TaskRevisionRef{Revision: *p.Revision, ManifestSHA256: p.ManifestSHA256}), contentRefValue(constraints)
		sd["executor"] = contentRefValue(executor)
		sp, err := publish(contentObject(sd), "spec_record")
		if err != nil {
			return err
		}
		out.Goal = TaskGoalRef{SpecID: in.SpecID, Spec: TaskRevisionRef{Revision: *sp.Revision, ManifestSHA256: sp.ManifestSHA256}}
		qid := queueID(root, target)
		q := foldQueue(r, qid)
		key := in.PublicationKey + "/queue"
		for _, event := range r.TaskQueueEvents {
			if event.Kind == "set" && contentString(contentFields(queueRequestValue(event.Request)), "publication_key") == key {
				entries := decodeQueueDraft(queueRequestValue(event.Request)).Entries
				if event.QueueID != qid || len(entries) == 0 || !contentTypedEqual(entries[len(entries)-1], TaskGoalQueueEntry{TaskID: in.TaskID, Goal: &out.Goal}) {
					return queueError("backup_goal_publication_conflict", "queue publication differs")
				}
				out.Applied = true
				out.Queue, err = buildQueueReadback(d, root, r, projects, in.Target, false)
				return err
			}
		}
		if existing != nil && (existing.Worktree != nil || preparationForTask(r, in.TaskID) != nil) {
			return queueError("backup_goal_task_bound", "backup Task already prepared without the append receipt")
		}
		entries := append(append([]TaskGoalQueueEntry{}, q.Pending...), TaskGoalQueueEntry{TaskID: in.TaskID, Goal: &out.Goal})
		v := contentRefValue(WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: key, ProjectID: target.ProjectID, RepoID: target.RepoID, EpicID: target.EpicID, ExpectedRevision: q.Revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "backup goal operator", ControlSurface: "ply workflow execute backup-goal", RecordedAtUTC: in.RecordedAtUTC}})
		if err = goalQueueDraftRule()(v); err != nil {
			return err
		}
		request := queueRequest(v)
		r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{qid, q.Revision + 1, "set", queueDigest(request), request, queueNow(d)})
		sortQueueRegistry(&r)
		if err = session.Publish(r); err != nil {
			return fmt.Errorf("backup queue append outcome requires readback; repeat identical input after inspection: %w", err)
		}
		r, err = session.Snapshot()
		if err != nil {
			return err
		}
		out.Applied = true
		out.Queue, err = buildQueueReadback(d, root, r, projects, in.Target, false)
		return err
	}
	if apply {
		err = d.ProjectLocks.WithSnapshotLock(root, func(p ProjectSnapshot) error {
			return d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
				r, e := s.Snapshot()
				if e != nil {
					return e
				}
				return run(p, s, r)
			})
		})
	} else {
		err = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error { return run(s.Projects, nil, s.Registry) })
	}
	return out, err
}
