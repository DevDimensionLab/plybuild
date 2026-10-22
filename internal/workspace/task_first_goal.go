package workspace

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// GoalDocumentInput names supplied product bytes. Ply computes native descriptors.
type GoalDocumentInput struct {
	ID      string `json:"id"`
	Locator string `json:"locator"`
}
type GoalDeliveryInput struct {
	Mode             string `json:"mode"`
	Acceptance       string `json:"acceptance"`
	ResponsibleActor string `json:"responsible_actor,omitempty"`
	RequireHumanQA   bool   `json:"require_human_qa,omitempty"`
	TimeoutSeconds   int    `json:"timeout_seconds,omitempty"`
	GitHubRepository string `json:"github_repository,omitempty"`
	Remote           string `json:"remote,omitempty"`
	TargetRef        string `json:"target_ref,omitempty"`
}
type FirstGoalInput struct {
	TaskID      TaskID                 `json:"task_id"`
	SpecID      string                 `json:"spec_id"`
	Title       string                 `json:"title"`
	Objective   string                 `json:"objective"`
	Criteria    []string               `json:"criteria"`
	Constraints []string               `json:"constraints"`
	Documents   []GoalDocumentInput    `json:"documents"`
	Target      QueueTargetInput       `json:"target"`
	Delivery    GoalDeliveryInput      `json:"delivery"`
	Executor    TaskExecutorAssignment `json:"executor"`
}
type GoalSourceBinding struct {
	ID        string `json:"id"`
	Locator   string `json:"locator"`
	SHA256    string `json:"sha256"`
	SizeBytes int    `json:"size_bytes"`
}
type FirstGoalResult struct {
	Title        string                     `json:"title"`
	Objective    string                     `json:"objective"`
	Criteria     []string                   `json:"criteria"`
	Constraints  []string                   `json:"constraints"`
	Executor     TaskExecutorAssignment     `json:"executor"`
	Kind         string                     `json:"kind"`
	TaskID       TaskID                     `json:"task_id"`
	Problem      *TaskRevisionRef           `json:"problem"`
	Goal         TaskGoalRef                `json:"goal"`
	Target       QueueTarget                `json:"target"`
	Delivery     DeliveryAgreement          `json:"delivery"`
	Sources      []GoalSourceBinding        `json:"sources"`
	CheckCwd     string                     `json:"check_cwd"`
	Readiness    string                     `json:"readiness"`
	Effects      []string                   `json:"effects"`
	Confirmation string                     `json:"confirmation"`
	Applied      bool                       `json:"applied"`
	Phase        string                     `json:"phase"`
	Contract     *TaskGoalContract          `json:"contract"`
	Queue        WorkspaceTaskQueueReadback `json:"queue"`
	NextArgv     []string                   `json:"next_argv"`
	Run          *string                    `json:"run"`
}

// RegisterFirstGoal publishes metadata under the existing registry CAS lock. The
// immutable Task description fences the whole supplied intent; native content and
// queue publications supply recovery receipts, including after a lost response.
func RegisterFirstGoal(d Dependencies, in FirstGoalInput, apply bool, confirmation string) (out FirstGoalResult, err error) {
	if err = validateFirstGoalInput(in); err != nil {
		return out, err
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	docs := []canonicaljson.Value{}
	for _, doc := range in.Documents {
		b, e := d.TaskContent.ReadSource(doc.Locator, taskDocumentLimit)
		if e != nil {
			return out, e
		}
		if e = taskDocumentBytesValid(b); e != nil {
			return out, e
		}
		media := "text/plain"
		if filepath.Ext(doc.Locator) == ".md" {
			media = "text/markdown"
		}
		docs = append(docs, contentRefValue(map[string]any{"id": doc.ID, "source": map[string]any{"kind": "file", "locator": doc.Locator, "sha256": digestTaskBytes(b), "size_bytes": len(b), "media_type": media, "git_provenance": nil}}))
	}
	sort.Slice(docs, func(i, j int) bool { return contentIDKey(docs[i]) < contentIDKey(docs[j]) })
	run := func(projects ProjectSnapshot, session WorkItemStoreSession, r WorkItemRegistry) (runErr error) {
		target, _, e := resolveQueueTarget(d, root, r, projects, in.Target)
		if e != nil {
			return e
		}
		observed, e := d.WorkGit.ObserveWorktree(target.ParentLocator)
		if e != nil {
			return e
		}
		if observed.Ref != target.ParentRef {
			return queueError("first_goal_target_conflict", "registered return ref is not checked out")
		}
		delivery, e := firstGoalDelivery(in.Delivery, target)
		if e != nil {
			return e
		}
		// Only real input and the stable registered return identity enter this digest.
		token := queueDigest(struct {
			Input     FirstGoalInput
			Documents []canonicaljson.Value
			Target    QueueTarget
			Delivery  DeliveryAgreement
		}{in, docs, target, delivery})
		key := "first-goal/" + string(in.TaskID)
		if contentKey(key+"/queue") != nil {
			return WorkInvalidArguments("Task ID is too long for the native registration key")
		}
		task := TaskRecord{ID: in.TaskID, Title: in.Title, Description: "First Goal registration " + token, ProjectID: target.ProjectID, RepoID: target.RepoID, ParentEpicID: target.EpicID, GitCommonDir: target.GitCommonDir, WorktreeState: WorkItemUnbound}
		existing, _ := findTask(r, in.TaskID)
		identity := task
		if existing != nil {
			identity.Worktree, identity.WorktreeState = existing.Worktree, existing.WorktreeState
		}
		if existing != nil && (!contentTypedEqual(*existing, identity) || taskPublication(r, "task-create/"+string(in.TaskID)) == nil) {
			return queueError("first_goal_identity_conflict", "Task ID already belongs to different first-registration input")
		}
		for _, s := range r.TaskSpecRevisions {
			if s.SpecID == in.SpecID && s.TaskID != in.TaskID {
				return queueError("first_goal_identity_conflict", "Spec ID belongs to another Task")
			}
		}
		out = FirstGoalResult{Title: in.Title, Objective: in.Objective, Criteria: in.Criteria, Constraints: in.Constraints, Executor: in.Executor, Kind: "ply.workspace.first-goal", TaskID: in.TaskID, Goal: TaskGoalRef{SpecID: in.SpecID, Spec: TaskRevisionRef{Revision: 1}}, Target: target, Delivery: delivery, Sources: []GoalSourceBinding{}, CheckCwd: target.ParentLocator, Readiness: "not_checked", Effects: []string{"publish Task and initial Problem", "preserve supplied documents and publish first Goal", "append exact Goal to queue"}, Confirmation: token, Phase: "preview", NextArgv: []string{"ply", "workflow", "execute", "--project", string(target.ProjectID), "--repo", string(target.RepoID), "--epic", string(target.EpicID), "--spec", in.SpecID, "--check"}}
		for _, v := range docs {
			m := contentFields(v)
			src := contentFields(m["source"])
			out.Sources = append(out.Sources, GoalSourceBinding{ID: contentString(m, "id"), Locator: contentString(src, "locator"), SHA256: contentString(src, "sha256"), SizeBytes: contentInt(src, "size_bytes")})
		}
		readback := func() error {
			if pub := taskPublication(r, "task-create/"+string(in.TaskID)); pub != nil {
				out.Problem = &TaskRevisionRef{Revision: *pub.OutcomeRef.Revision, ManifestSHA256: pub.OutcomeRef.ManifestSHA256}
			}
			if pub := taskPublication(r, key+"/goal"); pub != nil {
				out.Goal.Spec = TaskRevisionRef{Revision: *pub.OutcomeRef.Revision, ManifestSHA256: pub.OutcomeRef.ManifestSHA256}
				c, e := loadTaskGoal(d, root, r, in.TaskID, out.Goal)
				if e != nil {
					return e
				}
				out.Contract = &c
			}
			var e error
			out.Queue, e = buildQueueReadback(d, root, r, projects, in.Target, false)
			return e
		}
		if e = readback(); e != nil {
			return e
		}
		effectsStarted := false
		defer func() {
			if runErr != nil && effectsStarted && session != nil {
				out.Phase = "partial"
				observed, readErr := session.Snapshot()
				if readErr == nil {
					r = observed
					readErr = readback()
				}
				if readErr != nil {
					runErr = fmt.Errorf("%w; after-state readback failed: %v", runErr, readErr)
				}
			}
		}()
		// A completed append is historical success even if a later execution has
		// prepared, removed or completed this Goal. Never reconstruct its old queue.
		completed := false
		for _, ev := range r.TaskQueueEvents {
			if ev.Kind == "set" && contentString(contentFields(queueRequestValue(ev.Request)), "publication_key") == key+"/queue" {
				entries := decodeQueueDraft(queueRequestValue(ev.Request)).Entries
				if out.Contract == nil || ev.QueueID != queueID(root, target) || len(entries) == 0 || !contentTypedEqual(entries[len(entries)-1], TaskGoalQueueEntry{TaskID: in.TaskID, Goal: &out.Goal}) {
					return queueError("first_goal_publication_conflict", "append receipt differs from the exact registered Goal")
				}
				completed = true
			}
		}
		if apply && confirmation != token {
			return queueError("first_goal_confirmation_conflict", "preview confirmation differs from supplied product/source/target/policy/executor input")
		}
		if completed {
			out.Applied = true
			out.Phase = "registered"
			return nil
		}
		lifecycle, e := readWorkItemLifecycle(root)
		if e != nil {
			return e
		}
		if lifecycle.Task(in.TaskID) != LifecycleActive || lifecycle.Epic(target.EpicID) != LifecycleActive {
			return queueError("first_goal_inactive", "Task and Epic must be active")
		}
		if existing != nil && (existing.Worktree != nil || existing.WorktreeState != WorkItemUnbound || preparationForTask(r, in.TaskID) != nil) {
			return queueError("first_goal_task_bound", "Task already has started or bound work without a completed registration")
		}
		if op, _ := findOperation(r, in.TaskID); op != nil {
			return queueError("first_goal_task_bound", "Task already has a worktree operation without a completed registration")
		}
		state := taskContentState(r, in.TaskID)
		create := taskPublication(r, "task-create/"+string(in.TaskID))
		if existing != nil && (state.ProblemHead == nil || !contentTypedEqual(state.ProblemHead, &TaskRevisionRef{Revision: *create.OutcomeRef.Revision, ManifestSHA256: create.OutcomeRef.ManifestSHA256}) || !firstGoalHeadMatches(r, state.SpecHead, key+"/goal")) {
			return queueError("first_goal_identity_conflict", "Task history is not this first registration")
		}
		q := foldQueue(r, queueID(root, target))
		if len(q.Pending) >= 256 {
			return queueError("first_goal_queue_full", "queue append exceeds 256 entries")
		}
		for _, entry := range q.Pending {
			if entry.TaskID == in.TaskID {
				return queueError("first_goal_identity_conflict", "Task is already queued without this registration receipt")
			}
		}
		requirements := []TaskGoalRequirement{}
		for i, c := range in.Criteria {
			requirements = append(requirements, TaskGoalRequirement{ID: fmt.Sprintf("ac-%03d", i+1), Acceptance: c})
		}
		constraints := append([]string{}, in.Constraints...)
		sort.Strings(constraints)
		design := []canonicaljson.Value{}
		for _, doc := range docs {
			design = append(design, contentRefValue(map[string]any{"document_id": contentIDKey(doc), "section": nil}))
		}
		// Validate all caller content before creating the Task. The real Problem
		// head and original publication time replace these schema-valid placeholders.
		when := queueNow(d)
		problem := TaskRevisionRef{Revision: 1, ManifestSHA256: "sha256:" + fmt.Sprintf("%064d", 0)}
		if create != nil {
			when = create.RecordedAtUTC
			problem = TaskRevisionRef{Revision: *create.OutcomeRef.Revision, ManifestSHA256: create.OutcomeRef.ManifestSHA256}
		}
		draft := func() canonicaljson.Object {
			m := contentFields(contentEnvelope("WorkspaceTaskSpecDraft@2", map[string]canonicaljson.Value{"publication_key": key + "/goal", "task_id": string(in.TaskID), "registry_upgrade": nil, "expected_previous": nil, "title": in.Title, "change_reason": "Register the supplied first Goal.", "contract_kind": "goal", "objective": in.Objective, "problem": contentRefValue(problem), "requirements": contentRefValue(requirements), "design": design, "constraints": contentRefValue(constraints), "executor": contentRefValue(in.Executor), "documents": docs, "delivery": contentRefValue(delivery), "recorder": contentRefValue(TaskGoalQueueRecorder{ActorClaim: "local_cli_user", ControlSurface: "ply workspace task goal register", RecordedAtUTC: when})}))
			m["schema_version"] = int64(2)
			m["spec_id"] = in.SpecID
			return contentObject(m)
		}
		raw, e := canonicaljson.Marshal(draft())
		if e != nil {
			return e
		}
		if _, e = decodeTaskContent(raw, "spec_record", false, false); e != nil {
			return e
		}
		if !apply {
			return nil
		}
		// Return actual durable refs after a failed acknowledgement. The caller must
		// inspect these and repeat identical input; no rollback or dependent retry.
		effectsStarted = true
		if existing == nil {
			var upgrade *TaskRegistryUpgrade
			if r.FormatVersion < 3 {
				upgrade = &TaskRegistryUpgrade{FromVersion: r.FormatVersion, RegistrySHA256: r.RawSHA256}
			}
			_, e = createTaskProblem(d, root, session, r, projects, task, upgrade)
			if e != nil {
				return fmt.Errorf("Task/Problem may be published; inspect Task %s, then repeat identical registration: %w", in.TaskID, e)
			}
			r, e = session.Snapshot()
			if e != nil {
				return e
			}
			create = taskPublication(r, "task-create/"+string(in.TaskID))
			when = create.RecordedAtUTC
			problem = TaskRevisionRef{Revision: *create.OutcomeRef.Revision, ManifestSHA256: create.OutcomeRef.ManifestSHA256}
		}
		out.Phase = "task_published"
		raw, e = canonicaljson.Marshal(draft())
		if e != nil {
			return e
		}
		mutation := newTaskContentMutation(root, in.TaskID, key+"/goal", digestTaskBytes(raw))
		e = publishTaskContent(d, root, session, r, task, projects, draft(), "spec_record", raw, &mutation, false)
		if e != nil {
			return fmt.Errorf("Goal may be published; inspect Task %s key %s/goal, then repeat identical registration: %w", in.TaskID, key, e)
		}
		r, e = session.Snapshot()
		if e != nil {
			return e
		}
		if e = readback(); e != nil {
			return e
		}
		out.Phase = "goal_published"
		// The same registry lock encloses read + append. CAS in Publish protects
		// against stale replacement; unrelated current/pending entries survive.
		if e = upgradeQueueRegistry(d, root, &r, queueUpgrade(r)); e != nil {
			return e
		}
		q = foldQueue(r, queueID(root, target))
		entries := append(append([]TaskGoalQueueEntry{}, q.Pending...), TaskGoalQueueEntry{TaskID: in.TaskID, Goal: &out.Goal})
		v := contentRefValue(WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: key + "/queue", ProjectID: target.ProjectID, RepoID: target.RepoID, EpicID: target.EpicID, ExpectedRevision: q.Revision, Entries: entries, Recorder: TaskGoalQueueRecorder{ActorClaim: "local_cli_user", ControlSurface: "ply workspace task goal register", RecordedAtUTC: when}})
		if e = goalQueueDraftRule()(v); e != nil {
			return e
		}
		request := queueRequest(v)
		r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{queueID(root, target), q.Revision + 1, "set", queueDigest(request), request, queueNow(d)})
		sortQueueRegistry(&r)
		if e = session.Publish(r); e != nil {
			return fmt.Errorf("queue append may be published; inspect Task %s and queue, then repeat identical registration: %w", in.TaskID, e)
		}
		r, e = session.Snapshot()
		if e != nil {
			return e
		}
		out.Applied = true
		out.Phase = "registered"
		return readback()
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

func validateFirstGoalInput(in FirstGoalInput) error {
	if contentSlug(string(in.TaskID)) != nil || contentSlug(in.SpecID) != nil || contentText(128)(in.Title) != nil || contentText(8000)(in.Objective) != nil {
		return WorkInvalidArguments("explicit valid Task/Spec IDs, title and objective are required")
	}
	if in.Target.ProjectID == "" || in.Target.RepoID == "" || in.Target.EpicID == "" {
		return WorkInvalidArguments("explicit --project, --repo and --epic are required")
	}
	if len(in.Criteria) < 1 || len(in.Criteria) > 128 || len(in.Documents) < 1 || len(in.Documents) > 31 {
		return WorkInvalidArguments("provide 1..128 observable criteria and 1..31 supplied documents")
	}
	for _, c := range in.Criteria {
		if e := contentText(2000)(c); e != nil {
			return WorkInvalidArguments("invalid observable criterion: " + e.Error())
		}
	}
	seen := map[string]bool{}
	for _, doc := range in.Documents {
		if contentKey(doc.ID) != nil || seen[doc.ID] || contentPath(doc.Locator) != nil {
			return WorkInvalidArguments("documents require unique IDs and absolute physical file paths")
		}
		seen[doc.ID] = true
	}
	return nil
}

func firstGoalDelivery(in GoalDeliveryInput, t QueueTarget) (DeliveryAgreement, error) {
	a := DeliveryAgreement{SchemaVersion: 1, Mode: in.Mode, ProjectID: t.ProjectID, RepoID: t.RepoID, EpicID: t.EpicID, TargetRef: t.ParentRef, TargetWorktree: t.ParentLocator}
	if in.Mode == DeliveryPullRequest {
		a.TargetWorktree = ""
		a.TargetRef = in.TargetRef
		a.Remote = in.Remote
		a.GitHubRepository = in.GitHubRepository
	} else if in.TargetRef != "" || in.Remote != "" || in.GitHubRepository != "" {
		return a, WorkInvalidArguments("local delivery uses the explicit registered return; remote/PR target options are inapplicable")
	}
	switch in.Acceptance {
	case "human":
		if in.ResponsibleActor != "" || in.TimeoutSeconds != 0 || in.RequireHumanQA {
			return a, WorkInvalidArguments("automatic policy options require automatic acceptance")
		}
	case "automatic":
		a.SchemaVersion = 3
		a.Acceptance = &DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: in.ResponsibleActor, RequireHumanQA: in.RequireHumanQA, TimeoutSeconds: in.TimeoutSeconds}
	default:
		return a, WorkInvalidArguments("explicit acceptance policy human or automatic is required")
	}
	return a, ValidateDeliveryAgreement(a)
}

// Before the append receipt, a changed Goal head is a conflict, never an
// invitation to queue a stale first revision. Completed historical retries are
// handled earlier and may legitimately observe a derived execution head.
func firstGoalHeadMatches(r WorkItemRegistry, head *TaskRevisionRef, key string) bool {
	pub := taskPublication(r, key)
	if pub == nil {
		return head == nil
	}
	return contentTypedEqual(head, &TaskRevisionRef{Revision: *pub.OutcomeRef.Revision, ManifestSHA256: pub.OutcomeRef.ManifestSHA256})
}
