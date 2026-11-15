package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// plyActor is the actor Ply records for its own actions: --actor, else
// PLY_ACTOR, else unknown.
func plyActor(flag string) string {
	if flag != "" {
		return flag
	}
	if actor := os.Getenv("PLY_ACTOR"); actor != "" {
		return actor
	}
	return "unknown"
}

// taskPlanLink validates --plan and --plan-step against the plans of project
// and returns the link to record, or nil without --plan.
func taskPlanLink(d workspace.Dependencies, project workspace.ProjectID, planID, step string) (*taskjournal.PlyAction, error) {
	if planID == "" {
		if step != "" {
			return nil, workspace.WorkInvalidArguments("--plan-step requires --plan")
		}
		return nil, nil
	}
	resolved, err := resolveAssignmentProject(d, string(project))
	if err != nil {
		return nil, err
	}
	directory, err := workspace.ProjectDataPath(resolved, plan.DataPathCategory)
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
	link := &taskjournal.PlyAction{Operation: "plan_link", Plan: &planID}
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
func recordPly(d workspace.Dependencies, task, actor, done string, records ...taskjournal.PlyAction) error {
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
