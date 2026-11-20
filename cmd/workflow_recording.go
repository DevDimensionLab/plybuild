package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskactor"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// namedActor is the actor someone named: --actor, else PLY_ACTOR, else "".
func namedActor(flag string) string {
	if flag != "" {
		return flag
	}
	return os.Getenv("PLY_ACTOR")
}

// plyActor is the actor of an action outside a Task: the named actor, else
// unknown.
func plyActor(flag string) string {
	if actor := namedActor(flag); actor != "" {
		return actor
	}
	return taskactor.Unknown
}

// taskPlyActor is the actor of Ply's own action on task when nobody named
// one: the owner of the plan in a plan link among records, else of the plan
// the Task's journal links it to, else unknown.
func taskPlyActor(d taskwork.Dependencies, task string, records []taskjournal.PlyAction) string {
	link := (*taskjournal.PlyAction)(nil)
	for i := range records {
		if records[i].Operation == "plan_link" {
			link = &records[i]
		}
	}
	if link == nil {
		if _, journal, ok := taskRecords(d, task); ok {
			link = taskjournal.LatestPlanLink(journal)
		}
	}
	return taskactor.Of(true, "", planOwner(d, link))
}

// taskRecords reads the journal of task for the actor rule; a Task whose
// journal Ply cannot read gets no derived actor.
func taskRecords(d taskwork.Dependencies, task string) (workspace.ProjectWork, []taskjournal.Contribution, bool) {
	b, err := taskwork.ReadTaskJournalBasis(d, taskwork.TaskID(task))
	if err != nil {
		return workspace.ProjectWork{}, nil, false
	}
	records, _ := taskjournal.ReadContributions(b.Workspace, b.Work, task)
	return b.Work, records, true
}

// planOwner is the owner of the plan link names, or "" when it cannot be read.
func planOwner(d taskwork.Dependencies, link *taskjournal.PlyAction) string {
	if link == nil || link.Plan == nil {
		return ""
	}
	resolved, err := resolveAssignmentProject(d.Dependencies, link.Refs["project"])
	if err != nil {
		return ""
	}
	directory, err := resolved.Data().Plans()
	if err != nil {
		return ""
	}
	v, err := plan.Open(directory).Show(*link.Plan, 0)
	if err != nil {
		return ""
	}
	return v.Plan.Owner
}

// taskPlanLink validates --plan and --plan-step against the plans of project
// and returns the link to record, with the project and Task title so that a
// plan trace needs no registry read, or nil without --plan.
func taskPlanLink(d taskwork.Dependencies, project workspace.ProjectID, planID, step, title string) (*taskjournal.PlyAction, error) {
	if planID == "" {
		if step != "" {
			return nil, workspace.WorkInvalidArguments("--plan-step requires --plan")
		}
		return nil, nil
	}
	resolved, err := resolveAssignmentProject(d.Dependencies, string(project))
	if err != nil {
		return nil, err
	}
	directory, err := resolved.Data().Plans()
	if err != nil {
		return nil, err
	}
	v, err := plan.Open(directory).Show(planID, 0)
	var planError *plan.Error
	switch {
	case errors.As(err, &planError) && (planError.Class == plan.ClassNotFound || planError.Class == plan.ClassInvalidArguments):
		return nil, workspace.WorkInvalidArguments(fmt.Sprintf("--plan %s is not a plan of project %s", planID, project))
	case err != nil:
		return nil, err
	}
	link := &taskjournal.PlyAction{Operation: "plan_link", Plan: &planID, Refs: map[string]string{"project": string(project), "task_title": title}}
	if step != "" {
		if !slices.ContainsFunc(v.Plan.Steps, func(s plan.Step) bool { return s.ID == step }) {
			return nil, workspace.WorkInvalidArguments(fmt.Sprintf("--plan-step %s is not a step of %s", step, planID))
		}
		link.PlanStep = &step
	}
	return link, nil
}

// recordPly appends Ply's own records to the Task journal after an operation
// succeeded; a failure names the operation as done and the record as missing.
// An empty actor is the one taskPlyActor finds.
func recordPly(d taskwork.Dependencies, task, actor, done string, records ...taskjournal.PlyAction) error {
	if actor == "" {
		actor = taskPlyActor(d, task, records)
	}
	service := taskjournal.New(d)
	for _, action := range records {
		if _, err := service.AppendPly(task, actor, plyTitle(action), action); err != nil {
			return fmt.Errorf("%s; recording %s in the Task journal failed: %w", done, action.Operation, err)
		}
	}
	return nil
}

func plyTitle(action taskjournal.PlyAction) string {
	switch {
	case action.Operation == "plan_link" && action.PlanStep != nil:
		return fmt.Sprintf("Task belongs to %s step %s", *action.Plan, *action.PlanStep)
	case action.Operation == "plan_link":
		return "Task belongs to " + *action.Plan
	}
	return "Ply " + action.Operation
}

// recordEpicIntegration records epic_integrate in the journal of each Task
// whose integrated result reached the target with this fast-forward: it is
// in the new target commit and was not in the old one.
func recordEpicIntegration(d taskwork.Dependencies, input taskwork.EpicIntegrateInput, result taskwork.EpicIntegrateResult, actor string) error {
	done := "the target is fast-forwarded"
	epic, err := taskwork.ShowEpic(d, input.EpicID)
	if err != nil {
		return fmt.Errorf("%s; reading the Epic for its Task journals failed: %w", done, err)
	}
	gitDir := ""
	for _, b := range epic.Epic.RepoBindings {
		if b.RepoID == input.RepoID {
			gitDir = b.GitCommonDir
		}
	}
	w, err := workspace.ObserveContaining(d.Dependencies)
	if err != nil {
		return fmt.Errorf("%s; %w", done, err)
	}
	journals, err := taskjournal.Journals(w.Root)
	if err != nil {
		return fmt.Errorf("%s; %w", done, err)
	}
	ancestor := func(commit, of string) bool {
		return exec.Command("git", "--git-dir", gitDir, "merge-base", "--is-ancestor", commit, of).Run() == nil
	}
	refs := map[string]string{"epic": string(input.EpicID), "repo": string(input.RepoID), "target_ref": result.TargetRef, "before_oid": input.ExpectedTargetOID, "after_oid": result.TargetOID, "effect": result.Effect}
	for _, journal := range journals {
		task := journal.Task
		records, _ := taskjournal.ReadContributions(w, journal.Work, task)
		for _, r := range records {
			p := r.Event
			if p == nil || p.Ply == nil || p.Ply.Operation != "integration" || p.Ply.Refs["epic"] != string(input.EpicID) {
				continue
			}
			if oid := p.Ply.Refs["result_oid"]; oid != "" && ancestor(oid, result.TargetOID) && !ancestor(oid, input.ExpectedTargetOID) {
				if err := recordPly(d, task, actor, done, taskjournal.PlyAction{Operation: "epic_integrate", Refs: refs}); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}
