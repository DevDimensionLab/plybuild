package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
// and returns the link to record, with the project and Task title so that a
// plan trace needs no registry read, or nil without --plan.
func taskPlanLink(d workspace.Dependencies, project workspace.ProjectID, planID, step, title string) (*taskjournal.PlyAction, error) {
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

// recordEpicIntegration records epic_integrate in the journal of each Task
// whose integrated result reached the target with this fast-forward: it is
// in the new target commit and was not in the old one.
func recordEpicIntegration(d workspace.Dependencies, input workspace.EpicIntegrateInput, result workspace.EpicIntegrateResult, actor string) error {
	done := "the target is fast-forwarded"
	epic, err := workspace.ShowEpic(d, input.EpicID)
	if err != nil {
		return fmt.Errorf("%s; reading the Epic for its Task journals failed: %w", done, err)
	}
	gitDir := ""
	for _, b := range epic.Epic.RepoBindings {
		if b.RepoID == input.RepoID {
			gitDir = b.GitCommonDir
		}
	}
	w, err := workspace.ObserveContainingRoot(d)
	if err != nil {
		return fmt.Errorf("%s; %w", done, err)
	}
	tasks, err := taskjournal.JournalTasks(w.Root)
	if err != nil {
		return fmt.Errorf("%s; %w", done, err)
	}
	ancestor := func(commit, of string) bool {
		return exec.Command("git", "--git-dir", gitDir, "merge-base", "--is-ancestor", commit, of).Run() == nil
	}
	refs := map[string]string{"epic": string(input.EpicID), "repo": string(input.RepoID), "target_ref": result.TargetRef, "before_oid": input.ExpectedTargetOID, "after_oid": result.TargetOID, "effect": result.Effect}
	for _, task := range tasks {
		records, _ := taskjournal.ReadContributions(w, task)
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
