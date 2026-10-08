package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func addWorkspaceTaskContentCommands(task *cobra.Command, d workspace.Dependencies) {
	parents := map[string]*cobra.Command{}
	for _, name := range []string{"problem", "spec", "publication"} {
		p := &cobra.Command{Use: name, Short: map[string]string{"problem": "Preserve Task problem revisions", "spec": "Record Task goals and assess or select detailed solutions", "publication": "Inspect Task content publication outcomes"}[name], Long: "Record or inspect Task metadata and preserved documents. Goal Spec@2 keeps the objective, design, acceptance criteria and assigned implementor without an early Git base. A solution selection records a human claim. Goal recording does not start an agent; use workflow execute from the registered Epic to select and deliver a queued goal.", Example: "  ply workspace task " + name + " --help", RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() }}
		task.AddCommand(p)
		parents[name] = p
	}
	type mutation struct {
		parent, name, short string
		run                 func(workspace.Dependencies, workspace.TaskContentInput) (workspace.TaskContentMutationResult, error)
	}
	for _, entry := range []mutation{{"problem", "record", "Record an immutable problem revision", workspace.RecordTaskProblem}, {"spec", "record", "Record an immutable goal or solution revision", workspace.RecordTaskSpec}, {"spec", "assess", "Record a readiness assessment without selecting a solution", workspace.AssessTaskSpec}, {"spec", "select", "Record the human's exact solution choice", workspace.SelectTaskSolution}, {"spec", "withdraw", "Withdraw the current choice without deleting its history", workspace.WithdrawTaskSolution}} {
		entry := entry
		var file, format string
		c := &cobra.Command{Use: entry.name + " <task-id>", Short: entry.short, Long: entry.short + ". Read a complete strict JSON draft from an absolute physical regular UTF-8 file. This records metadata; it does not start an agent or authorize worktree creation, product QA, or integration. Select and withdraw preserve an explicitly reported human claim.", Example: "  ply workspace task " + entry.parent + " " + entry.name + " explain-start-errors --file /absolute/" + entry.name + ".json --format json", Args: cobra.ExactArgs(1), PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, e := workspace.ParseTaskID(args[0]); e != nil {
				return e
			}
			return validateWorkFormat(format)
		}, RunE: func(cmd *cobra.Command, args []string) error {
			id, e := workspace.ParseTaskID(args[0])
			if e != nil {
				return e
			}
			r, mutationErr := entry.run(d, workspace.TaskContentInput{TaskID: id, File: file})
			if format == "json" && r.Workspace != "" {
				b, e := workspace.MarshalTaskContentMutation(r)
				if e != nil {
					return e
				}
				if _, e = cmd.OutOrStdout().Write(append(b, '\n')); e != nil {
					return fmt.Errorf("%w; inspect ply workspace task publication show %s --key %s", e, id, r.PublicationKey)
				}
			} else if mutationErr == nil {
				if _, e = fmt.Fprintf(cmd.OutOrStdout(), "Task %s: %s publication %s.\nMetadata recorded; no next transition authorized.\nNext action: ply workspace task show %s\n", id, r.Classification, r.PublicationKey, id); e != nil {
					return fmt.Errorf("%w; inspect ply workspace task publication show %s --key %s", e, id, r.PublicationKey)
				}
			}
			return mutationErr
		}}
		if entry.parent == "spec" && (entry.name == "record" || entry.name == "assess" || entry.name == "select") {
			c.Long += " Solution Spec@1 uses the current Epic base; after a base update record, assess and select a new solution revision. Goal Spec@2 (contract_kind goal) has no implementation_basis or detailed technical plan: queue its exact revision with QueueDraft@2, then workflow execute binds the current base and preserves a derived execution Spec. Planner goal publication does not claim human selection. Historical revisions remain immutable."
		}
		c.Flags().StringVar(&file, "file", "", "absolute physical path to a complete JSON draft (maximum 256 KiB)")
		_ = c.MarkFlagRequired("file")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		setWorkFlagErrors(c)
		parents[entry.parent].AddCommand(c)
	}
	type query struct {
		parent, name, short string
		run                 func(workspace.Dependencies, workspace.TaskContentQuery) (workspace.TaskContentReadbackResult, error)
	}
	for _, entry := range []query{{"problem", "show", "Show the recorded problem and preserved sources", workspace.ShowTaskProblem}, {"spec", "list", "List goal and solution revisions and their assessments", workspace.ListTaskSpecs}, {"spec", "show", "Show an exact goal or solution revision and its preserved inputs", workspace.ShowTaskSpec}, {"publication", "show", "Inspect the outcome of a Task content publication", workspace.ShowTaskPublication}} {
		entry := entry
		var format, spec, key string
		var revision int
		c := &cobra.Command{Use: entry.name + " <task-id>", Short: entry.short, Long: entry.short + ". This is read-only and does not migrate a registry, create locks, repair snapshots, or authorize another action. JSON includes preserved descriptors and exact input paths and hashes for handoff preparation.", Example: "  ply workspace task " + entry.parent + " " + entry.name + " explain-start-errors --format json", Args: cobra.ExactArgs(1), PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, e := workspace.ParseTaskID(args[0]); e != nil {
				return e
			}
			if cmd.Flags().Changed("revision") && revision < 1 {
				return workspace.WorkInvalidArguments("--revision must be a positive integer")
			}
			return validateWorkFormat(format)
		}, RunE: func(cmd *cobra.Command, args []string) error {
			id, e := workspace.ParseTaskID(args[0])
			if e != nil {
				return e
			}
			r, e := entry.run(d, workspace.TaskContentQuery{TaskID: id, SpecID: spec, Revision: revision, PublicationKey: key})
			if e != nil {
				return taskContentCommandError(cmd, format, e)
			}
			if format == "json" {
				b, e := workspace.MarshalTaskContentReadback(r)
				if e != nil {
					return e
				}
				_, e = cmd.OutOrStdout().Write(append(b, '\n'))
				return e
			}
			_, e = fmt.Fprint(cmd.OutOrStdout(), workspace.TaskContentText(r.Value))
			return e
		}}
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if entry.name == "show" && entry.parent != "publication" {
			c.Flags().IntVar(&revision, "revision", 0, "exact positive revision number")
		}
		if entry.parent == "spec" && entry.name == "show" {
			c.Flags().StringVar(&spec, "spec", "", "Task-owned Spec ID")
			_ = c.MarkFlagRequired("spec")
			_ = c.MarkFlagRequired("revision")
			c.Example = "  ply workspace task spec show explain-start-errors --spec explain-errors --revision 1 --format json"
		}
		if entry.parent == "publication" {
			c.Flags().StringVar(&key, "key", "", "exact publication key")
			_ = c.MarkFlagRequired("key")
			c.Example = "  ply workspace task publication show explain-start-errors --key solution/r1 --format json"
		}
		setWorkFlagErrors(c)
		parents[entry.parent].AddCommand(c)
	}
}
