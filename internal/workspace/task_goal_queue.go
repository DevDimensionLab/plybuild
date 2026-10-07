package workspace

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// TaskGoalQueueEntry also permits an existing solution selection, so a queue can
// retain legacy work while planners publish base-free goals. Exactly one field
// is present in a version-2 wire entry.
type TaskGoalQueueEntry struct {
	TaskID    TaskID           `json:"task_id"`
	Selection *TaskDecisionRef `json:"selection,omitempty"`
	Goal      *TaskGoalRef     `json:"goal,omitempty"`
}

func (e TaskGoalQueueEntry) MarshalJSON() ([]byte, error) {
	if e.Goal != nil {
		if e.Selection != nil {
			return nil, fmt.Errorf("a queue entry cannot bind both a goal and a solution selection")
		}
		return json.Marshal(struct {
			TaskID TaskID       `json:"task_id"`
			Goal   *TaskGoalRef `json:"goal"`
		}{e.TaskID, e.Goal})
	}
	return json.Marshal(struct {
		TaskID    TaskID           `json:"task_id"`
		Selection *TaskDecisionRef `json:"selection"`
	}{e.TaskID, e.Selection})
}

type TaskGoalQueueRecorder struct {
	ActorClaim     string `json:"actor_claim"`
	ControlSurface string `json:"control_surface"`
	RecordedAtUTC  string `json:"recorded_at_utc"`
}

type WorkspaceTaskGoalQueueDraft struct {
	Kind             string                `json:"kind"`
	SchemaVersion    int                   `json:"schema_version"`
	PublicationKey   string                `json:"publication_key"`
	ProjectID        ProjectID             `json:"project_id"`
	RepoID           RepoID                `json:"repo_id"`
	EpicID           EpicID                `json:"epic_id"`
	ExpectedRevision int                   `json:"expected_revision"`
	Entries          []TaskGoalQueueEntry  `json:"entries"`
	Recorder         TaskGoalQueueRecorder `json:"recorder"`
	RegistryUpgrade  *TaskRegistryUpgrade  `json:"registry_upgrade"`
}

func goalQueueDraftRule() contentRule {
	entry := func(v canonicaljson.Value) error {
		if _, goal := contentFields(v)["goal"]; goal {
			return contentExact(map[string]contentRule{"task_id": contentSlug, "goal": taskGoalRefRule()})(v)
		}
		return contentExact(map[string]contentRule{"task_id": contentSlug, "selection": contentNullable(contentDecisionRule("sel_"))})(v)
	}
	return contentExact(map[string]contentRule{
		"kind": contentEnum("WorkspaceTaskQueueDraft@2"), "schema_version": contentInteger(2, 2),
		"publication_key": contentKey, "project_id": contentSlug, "repo_id": contentSlug, "epic_id": contentSlug,
		"expected_revision": contentInteger(0, 2147483647),
		"entries":           contentList(256, entry, func(v canonicaljson.Value) string { return contentString(contentFields(v), "task_id") }, true),
		"recorder":          contentRecorderRule(),
		"registry_upgrade":  contentNullable(contentExact(map[string]contentRule{"from_version": contentInteger(1, 3), "registry_sha256": contentDigest})),
	})
}

type decodedTaskQueueDraft struct {
	PublicationKey   string               `json:"publication_key"`
	ProjectID        ProjectID            `json:"project_id"`
	RepoID           RepoID               `json:"repo_id"`
	EpicID           EpicID               `json:"epic_id"`
	ExpectedRevision int                  `json:"expected_revision"`
	Entries          []TaskGoalQueueEntry `json:"entries"`
	RegistryUpgrade  *TaskRegistryUpgrade `json:"registry_upgrade"`
}

func decodeQueueDraft(v canonicaljson.Value) decodedTaskQueueDraft {
	var out decodedTaskQueueDraft
	_ = contentDecode(v, &out)
	return out
}

func queueGoalEvaluation(d Dependencies, root string, r WorkItemRegistry, entry TaskGoalQueueEntry, rank int) QueuePending {
	row := QueuePending{Rank: rank, TaskID: entry.TaskID, Goal: entry.Goal, State: "blocked", Reasons: []QueueReason{}}
	task, _ := findTask(r, entry.TaskID)
	if task == nil {
		row.Reasons = queueReasons("task_queue_task_missing", "queued Task is missing")
		return row
	}
	goal, e := loadTaskGoal(d, root, r, entry.TaskID, *entry.Goal)
	if e != nil {
		row.Reasons = queueErrorReasons(e)
		return row
	}
	row.Title = goal.Title
	row.SpecID, row.SpecRevision = &goal.Goal.SpecID, &goal.Goal.Spec.Revision
	row.Executor = &goal.Executor
	row.Delivery = goal.Delivery
	if !contentTypedEqual(taskContentState(r, task.ID).ProblemHead, &goal.Problem) {
		row.Reasons = queueReasons("task_goal_problem_changed", "queued goal no longer refers to the current problem")
		return row
	}
	if task.Worktree != nil || task.WorktreeState != WorkItemUnbound || preparationForTask(r, task.ID) != nil {
		row.Reasons = queueReasons("task_queue_task_bound", "Task already has a worktree or preparation")
		return row
	}
	if op, _ := findOperation(r, task.ID); op != nil {
		row.Reasons = queueReasons("task_queue_task_bound", "Task already has an intent")
		return row
	}
	row.State = "ready"
	return row
}
