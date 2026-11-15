package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/spf13/cobra"
)

// newPlanJournalCommand builds plan journal step start|finish and note: the
// plan owner's own steps, which belong to the plan rather than to one Task.
func newPlanJournalCommand(project *string, open func() (plan.Store, error)) *cobra.Command {
	journal := &cobra.Command{Use: "journal", Short: "Record the plan owner's own steps and notes, which ply plan trace shows"}
	step := &cobra.Command{Use: "step", Short: "Record the start or finish of a plan step"}
	for _, verb := range []string{"start", "finish", "note"} {
		verb := verb
		var id, kind, parent, title, outcome, actor, role string
		c := &cobra.Command{
			Use:   map[string]string{"start": "start <PLAN-n> --id <step-id> --kind <kind> --title <title>", "finish": "finish <PLAN-n> --id <step-id> --kind <kind> --title <title> [--outcome <outcome>]", "note": "note <PLAN-n> --title <title> [--outcome <outcome>]"}[verb],
			Short: map[string]string{"start": "Record that a plan step started now", "finish": "Record that a plan step finished now", "note": "Record a note on the plan now"}[verb],
			Long: "For work of the plan owner outside any Task, such as mapping, a user's choice, checking a return or clean-up. The record gets Ply's current time.\n" +
				"Step kinds: " + strings.Join(taskjournal.StepKinds, ", ") + ". Outcomes: " + strings.Join(taskjournal.Outcomes, ", ") + ".",
			Example: "  ply plan journal step start PLAN-1 --project ply --id map --kind planning --title 'Map the backlog' --actor ply-plan\n  ply plan journal note PLAN-1 --project ply --title 'User chose option 1' --actor ply-plan",
			Args:    cobra.ExactArgs(1),
			RunE: func(c *cobra.Command, args []string) error {
				e := plan.JournalEntry{Type: "note", Title: title, ActorID: plyActor(actor), Role: role}
				if verb != "note" {
					if id == "" || !slices.Contains(taskjournal.StepKinds, kind) {
						return planInvalidArguments("--id and a --kind of " + strings.Join(taskjournal.StepKinds, ", ") + " are required")
					}
					e.Type = "step_" + map[string]string{"start": "started", "finish": "finished"}[verb]
					e.StepID, e.StepKind = &id, &kind
					if parent != "" {
						e.ParentStepID = &parent
					}
				}
				if outcome != "" {
					if !slices.Contains(taskjournal.Outcomes, outcome) {
						return planInvalidArguments("--outcome must be one of " + strings.Join(taskjournal.Outcomes, ", "))
					}
					e.Outcome = &outcome
				}
				if !slices.Contains(taskjournal.Roles, role) {
					return planInvalidArguments("--role must be one of " + strings.Join(taskjournal.Roles, ", "))
				}
				s, err := open()
				if err != nil {
					return err
				}
				written, err := s.AppendJournal(args[0], e)
				if err != nil {
					return err
				}
				_, err = fmt.Fprintf(c.OutOrStdout(), "Recorded: %s\n", written.At)
				return err
			},
		}
		if verb != "note" {
			c.Flags().StringVar(&id, "id", "", "step ID, the same for start and finish (required)")
			c.Flags().StringVar(&kind, "kind", "", "step kind (required)")
			c.Flags().StringVar(&parent, "parent", "", "ID of the step this one is inside")
		}
		if verb != "start" {
			c.Flags().StringVar(&outcome, "outcome", "", "reported outcome")
		}
		c.Flags().StringVar(project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&title, "title", "", "one line (required)")
		c.Flags().StringVar(&actor, "actor", "", "who reports (default: PLY_ACTOR)")
		c.Flags().StringVar(&role, "role", "planner", "role: "+strings.Join(taskjournal.Roles, ", "))
		_ = c.MarkFlagRequired("title")
		if verb == "note" {
			journal.AddCommand(c)
		} else {
			step.AddCommand(c)
		}
	}
	journal.AddCommand(step)
	return journal
}
