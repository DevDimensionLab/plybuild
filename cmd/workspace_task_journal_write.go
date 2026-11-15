package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// addWorkspaceTaskJournalWriteCommands adds the commands that write the
// common records without a JSON file: link, step start|finish and note.
func addWorkspaceTaskJournalWriteCommands(journal *cobra.Command, d workspace.Dependencies) {
	var planID, planStep, linkActor string
	link := &cobra.Command{
		Use:   "link <task-id> --plan <PLAN-n> [--plan-step <id>]",
		Short: "Link an existing Task to a plan step",
		Long: "Record a plan link for a Task that is already registered, for example one registered before --plan existed or with a link from an earlier Ply.\n" +
			"The plan and step are validated against the plans of the Task's Project. The latest link of a Task counts; ply plan trace reads it.",
		Example: "  ply workspace task journal link workflow-viz-d10-impl --plan PLAN-39 --plan-step impl-d10 --actor ply-plan",
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			id, err := workspace.ParseTaskID(args[0])
			if err != nil {
				return err
			}
			if planID == "" {
				return workspace.WorkInvalidArguments("--plan is required")
			}
			task, err := workspace.ShowTask(d, id)
			if err != nil {
				return err
			}
			l, err := taskPlanLink(d, task.Task.ProjectID, planID, planStep, task.Task.Title)
			if err != nil {
				return err
			}
			if err = recordPly(d, args[0], plyActor(linkActor), "nothing changed", *l); err != nil {
				return err
			}
			_, err = fmt.Fprintf(c.OutOrStdout(), "%s\n", plyTitle(*l))
			return err
		},
	}
	link.Flags().StringVar(&planID, "plan", "", "PLAN-<n> of the Task's Project (required)")
	link.Flags().StringVar(&planStep, "plan-step", "", "step ID of --plan the Task carries")
	link.Flags().StringVar(&linkActor, "actor", "", "who links (default: PLY_ACTOR)")
	journal.AddCommand(link)

	step := &cobra.Command{Use: "step", Short: "Record the start or finish of a work step without a JSON file"}
	for _, verb := range []string{"start", "finish"} {
		step.AddCommand(newJournalWriteCommand(d, verb))
	}
	journal.AddCommand(step)
	journal.AddCommand(newJournalWriteCommand(d, "note"))
}

// newJournalWriteCommand builds step start, step finish or note. The record
// gets Ply's current time; its outcome and text stay the reporter's claims.
func newJournalWriteCommand(d workspace.Dependencies, verb string) *cobra.Command {
	var id, kind, parent, title, outcome, actor, role string
	var sources []string
	use := map[string]string{
		"start":  "start <task-id> --id <step-id> --kind <kind> --title <title>",
		"finish": "finish <task-id> --id <step-id> --kind <kind> --title <title> [--outcome <outcome>]",
		"note":   "note <task-id> --title <title> [--outcome <outcome>]",
	}[verb]
	c := &cobra.Command{
		Use:   use,
		Short: map[string]string{"start": "Record that a work step started now", "finish": "Record that a work step finished now", "note": "Record a note now, for example at a return"}[verb],
		Long: "The record gets Ply's current time and needs no source; --source adds an existing file that Ply keeps a copy of.\n" +
			"A finish repeats the --id and --kind of its start. A step inside another, such as a correction inside an implementation, names the outer step with --parent.\n" +
			"Step kinds: clarification, planning, implementation, verification, review, correction, coordination. Outcomes: pass, fail, complete, blocked, unknown, deferred, declined.",
		Example: map[string]string{
			"start":  "  ply workspace task journal step start my-task --id impl --kind implementation --title 'Implement' --actor dev-1\n  ply workspace task journal step start my-task --id fix-1 --kind correction --parent impl --title 'Fix review finding'",
			"finish": "  ply workspace task journal step finish my-task --id impl --kind implementation --title 'Implement' --outcome complete",
			"note":   "  ply workspace task journal note my-task --title 'Returned to coordinator' --outcome complete",
		}[verb],
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			now := time.Now().UTC().Format(time.RFC3339Nano)
			nonce := make([]byte, 8)
			if _, err := rand.Read(nonce); err != nil {
				return err
			}
			in := taskjournal.EventInput{Common: taskjournal.Common{Kind: "ply.workspace.task-journal-event-input", SchemaVersion: 1, PublicationKey: fmt.Sprintf("cli:%x", nonce), TaskID: args[0],
				Actor: taskjournal.Actor{ID: plyActor(actor), Role: role}, OccurredAt: &now, TimeBasis: taskjournal.TimeBasis{Kind: "observed", Clock: &[]string{"ply"}[0], Precision: "nanosecond"}, Sources: []taskjournal.SourceRef{}},
				Type: "note", Title: title, Relations: []taskjournal.Relation{}}
			if verb != "note" {
				if id == "" || kind == "" {
					return workspace.WorkInvalidArguments("--id and --kind are required")
				}
				in.Type = "step_" + map[string]string{"start": "started", "finish": "finished"}[verb]
				in.Step = &taskjournal.StepRef{ID: id, Kind: kind}
				if parent != "" {
					in.Step.ParentStepID = &parent
				}
			}
			if outcome != "" {
				in.Outcome = &outcome
			}
			for _, path := range sources {
				ref, err := journalSource(path)
				if err != nil {
					return err
				}
				in.Sources = append(in.Sources, ref)
			}
			r, err := taskjournal.New(d).AppendEvent(args[0], in)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(c.OutOrStdout(), "Event: %s\nRecorded: %s\n", r.EventID, r.RecordedAt)
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
	c.Flags().StringVar(&title, "title", "", "one line (required)")
	c.Flags().StringVar(&actor, "actor", "", "who reports (default: PLY_ACTOR)")
	c.Flags().StringVar(&role, "role", "developer", "role: human, planner, developer, reviewer, observer or tool")
	c.Flags().StringArrayVar(&sources, "source", nil, "existing file as evidence; Ply keeps a copy (repeatable)")
	_ = c.MarkFlagRequired("title")
	return c
}

func journalSource(path string) (taskjournal.SourceRef, error) {
	abs, err := filepath.Abs(path)
	if err == nil {
		abs, err = filepath.EvalSymlinks(abs)
	}
	if err != nil {
		return taskjournal.SourceRef{}, workspace.WorkInvalidArguments("--source " + path + " is not an existing file")
	}
	f, err := os.Open(abs)
	if err != nil {
		return taskjournal.SourceRef{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return taskjournal.SourceRef{}, err
	}
	return taskjournal.SourceRef{Locator: abs, SHA256: fmt.Sprintf("%x", h.Sum(nil))}, nil
}
