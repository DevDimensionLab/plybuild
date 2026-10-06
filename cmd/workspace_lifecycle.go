package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceLifecycleServices struct {
	show func(string, string) (workspace.WorkItemLifecycleReadback, error)
	set  func(workspace.WorkItemLifecycleInput) (workspace.WorkItemLifecycleMutation, error)
}

func newWorkspaceLifecycleCommand(d workspace.Dependencies, subjectKind string) *cobra.Command {
	return newWorkspaceLifecycleCommandWithServices(subjectKind, workspaceLifecycleServices{
		show: func(kind, id string) (workspace.WorkItemLifecycleReadback, error) {
			return workspace.ShowWorkItemLifecycle(d, kind, id)
		},
		set: func(input workspace.WorkItemLifecycleInput) (workspace.WorkItemLifecycleMutation, error) {
			return workspace.SetWorkItemLifecycle(d, input)
		},
	})
}

func newWorkspaceLifecycleCommandWithServices(kind string, services workspaceLifecycleServices) *cobra.Command {
	label := strings.ToUpper(kind[:1]) + kind[1:]
	command := &cobra.Command{Use: "lifecycle", Short: "Organize " + label + " lifecycle", Long: "Record active, parked, frozen, or archived lifecycle with an actor and durable history. Lifecycle organizes lists and attention; it does not stop agents, change Git, or authorize workflow transitions.", Example: "  ply workspace " + kind + " lifecycle show an-item\n  ply workspace " + kind + " lifecycle set an-item archived --actor codex"}
	var showFormat string
	show := &cobra.Command{Use: "show <" + kind + "-id>", Short: "Show lifecycle and its recorded changes", Long: "Read the current lifecycle and complete change history for one registered " + label + ". Older items default to active without writing metadata. Revision is the shared workspace lifecycle revision.", Example: "  ply workspace " + kind + " lifecycle show an-item --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments("expected exactly one " + label + " ID")
		}
		if err := validateLifecycleCommandID(kind, args[0]); err != nil {
			return err
		}
		return validateWorkFormat(showFormat)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		out, err := services.show(kind, args[0])
		if err != nil {
			return err
		}
		if showFormat == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s lifecycle: %s (workspace revision %d).\n", label, out.SubjectID, out.Lifecycle, out.Revision); err != nil {
			return err
		}
		for _, event := range out.Events {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %d: %s -> %s by %s at %s\n", event.Sequence, event.From, event.To, event.ActorClaim, event.RecordedAtUTC); err != nil {
				return err
			}
		}
		return nil
	}}
	show.Flags().StringVar(&showFormat, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	var actor, setFormat string
	var expectedRevision int
	input := func(cmd *cobra.Command, args []string) workspace.WorkItemLifecycleInput {
		out := workspace.WorkItemLifecycleInput{SubjectKind: kind, SubjectID: args[0], State: workspace.LifecycleState(args[1]), ActorClaim: actor}
		if cmd.Flags().Changed("expected-revision") {
			out.ExpectedRevision = &expectedRevision
		}
		return out
	}
	set := &cobra.Command{Use: "set <" + kind + "-id> <active|parked|frozen|archived>", Short: "Record an explicit lifecycle change", Long: "Record one lifecycle change and the supplied actor claim in an atomic history update. Repeating the current state writes nothing. --expected-revision checks the shared workspace lifecycle revision, including changes to other items. Actor claims describe the recorder and do not attest human QA or grant authority.", Example: "  ply workspace " + kind + " lifecycle set an-item archived --actor codex\n  ply workspace " + kind + " lifecycle set an-item active --actor human --expected-revision 3 --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return workspace.WorkInvalidArguments("expected one " + label + " ID and one lifecycle value")
		}
		if err := validateWorkFormat(setFormat); err != nil {
			return err
		}
		return input(cmd, args).Validate()
	}, RunE: func(cmd *cobra.Command, args []string) error {
		out, err := services.set(input(cmd, args))
		if err != nil {
			return err
		}
		if setFormat == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
		}
		verb := "Set"
		if !out.Changed {
			verb = "Kept"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s lifecycle to %s (workspace revision %d).\n", verb, label, out.SubjectID, out.Lifecycle, out.Revision)
		return err
	}}
	set.Flags().StringVar(&actor, "actor", "", "recorder identity or actor claim")
	set.Flags().IntVar(&expectedRevision, "expected-revision", 0, "require this workspace lifecycle revision before writing")
	set.Flags().StringVar(&setFormat, "format", "text", "output format (text or json)")
	_ = set.MarkFlagRequired("actor")
	setWorkFlagErrors(set)
	command.AddCommand(show, set)
	return command
}

func validateLifecycleCommandID(kind, id string) error {
	if kind == "task" {
		_, err := workspace.ParseTaskID(id)
		return err
	}
	if kind == "epic" {
		_, err := workspace.ParseEpicID(id)
		return err
	}
	return workspace.WorkInvalidArguments("lifecycle subject must be task or epic")
}
