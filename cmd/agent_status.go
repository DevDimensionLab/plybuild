package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agentstatus"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// Fixed JSON contract names of the agent-status readbacks.
const (
	agentStatusReadbackKind  = "PlyAgentStatusReadback@1"
	agentStatusListKind      = "PlyAgentStatusList@1"
	agentStatusClearKind     = "PlyAgentStatusClear@1"
	agentStatusSchemaVersion = 1
)

func newAgentStatusCommand(d workspace.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:   "agent-status",
		Short: "Set, clear and list the status agents report about themselves in a Project",
		Long: "An agent reports what a terminal cannot show: that it works, waits for its own background job, delegate or plan (waiting_background), waits for a human (needs_answer) or is done. Ply stores one post per actor and never contacts Herdr or the agent.\n" +
			"Every post expires: --for is required, at most 4h, or 24h for needs_answer. list shows expired posts with expired: true; set removes other actors' posts expired for over 24 hours.\n" +
			"Posts are files in <project wrapper>/data/agent-status/.\n" +
			"Without --project the single Project whose wrapper contains the current directory is used.",
		Example: "  ply agent-status set --project ply --actor claude:<session-id> --pane w7:pT --state waiting_background --ref PLAN-40 --for 30m\n" +
			"  ply agent-status list --project ply --format json\n" +
			"  ply agent-status clear --project ply --actor claude:<session-id>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return agentStatusInvalidArguments("unknown agent-status command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	command.SetFlagErrorFunc(func(*cobra.Command, error) error {
		return agentStatusInvalidArguments("invalid flag")
	})

	var project, format, actor, pane, state, ref, duration string
	addFlags := func(c *cobra.Command) *cobra.Command {
		c.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		command.AddCommand(c)
		return c
	}
	validate := func(required ...string) cobra.PositionalArgs {
		return func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return agentStatusInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			if format != "text" && format != "json" {
				return agentStatusInvalidArguments("--format must be text or json")
			}
			for _, name := range required {
				if !cmd.Flags().Changed(name) {
					return agentStatusInvalidArguments("--" + name + " is required")
				}
			}
			return nil
		}
	}
	var resolved workspace.ProjectRecord
	open := func() (agentstatus.Store, error) {
		var err error
		if resolved, err = resolveAgentStatusProject(d, project); err != nil {
			return agentstatus.Store{}, err
		}
		directory, err := resolved.Data().AgentStatus()
		if err != nil {
			return agentstatus.Store{}, err
		}
		return agentstatus.Open(directory), nil
	}

	setCommand := addFlags(&cobra.Command{
		Use:   "set --actor <provider:id> --pane <wX:pY> --state <state> [--ref <id>] --for <duration>",
		Short: "Replace your own post and print it",
		Long: "States: " + strings.Join(agentstatus.States, ", ") + ". --ref is an ID such as ASSIGNMENT-61, PLAN-36/arch or a backlog ID, never free text. --for is a Go duration such as 30m or 2h.\n" +
			"Ply sets reported_at and expires_at. Set again to renew; clear when you continue, finish or restart.",
		Args: validate("actor", "pane", "state", "for"),
		RunE: func(cmd *cobra.Command, args []string) error {
			parsed, err := time.ParseDuration(duration)
			if err != nil {
				return agentStatusInvalidArguments("--for must be a duration such as 30m or 2h")
			}
			s, err := open()
			if err != nil {
				return err
			}
			item, err := s.Set(actor, pane, state, ref, parsed)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					Item          agentstatus.Item    `json:"item"`
				}{agentStatusReadbackKind, agentStatusSchemaVersion, resolved.ID, item})
			}
			_, err = io.WriteString(cmd.OutOrStdout(), agentStatusLine(item))
			return err
		},
	})
	setCommand.Flags().StringVar(&actor, "actor", "", "your agent label, <provider>:<id> (required)")
	setCommand.Flags().StringVar(&pane, "pane", "", "your Herdr pane, such as w7:pT (required)")
	setCommand.Flags().StringVar(&state, "state", "", "one of "+strings.Join(agentstatus.States, ", ")+" (required)")
	setCommand.Flags().StringVar(&ref, "ref", "", "optional ID of what you wait for or work on")
	setCommand.Flags().StringVar(&duration, "for", "", "how long the post holds, at most 4h or 24h for needs_answer (required)")

	clearCommand := addFlags(&cobra.Command{
		Use:   "clear --actor <provider:id>",
		Short: "Remove your own post",
		Long:  "Clearing a missing post is not an error; removed says whether a post was removed.",
		Args:  validate("actor"),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			removed, err := s.Clear(actor)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					Actor         string              `json:"actor"`
					Removed       bool                `json:"removed"`
				}{agentStatusClearKind, agentStatusSchemaVersion, resolved.ID, actor, removed})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s removed: %t\n", actor, removed)
			return err
		},
	})
	clearCommand.Flags().StringVar(&actor, "actor", "", "your agent label, <provider>:<id> (required)")

	addFlags(&cobra.Command{
		Use:   "list",
		Short: "List every post of the Project by actor, with expiry computed now",
		Long:  "Files that cannot be read are listed under invalid. Text lines are <actor> <state> <pane> [<ref>] until <expires_at>, with (expired) when it has passed.",
		Args:  validate(),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			result, err := s.List()
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					agentstatus.ListResult
				}{agentStatusListKind, agentStatusSchemaVersion, resolved.ID, result})
			}
			var out strings.Builder
			for _, item := range result.Items {
				out.WriteString(agentStatusLine(item))
			}
			for _, invalid := range result.Invalid {
				fmt.Fprintf(&out, "invalid %s: %s\n", invalid.File, invalid.Error)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	})
	return command
}

func agentStatusLine(item agentstatus.Item) string {
	fields := []string{item.Actor, item.State, item.Pane}
	if item.Ref != "" {
		fields = append(fields, item.Ref)
	}
	line := strings.Join(fields, " ") + " until " + item.ExpiresAt
	if item.Expired {
		line += " (expired)"
	}
	return line + "\n"
}

// resolveAgentStatusProject maps Project lookup failures to fixed details
// without paths, as resolvePlanProject does.
func resolveAgentStatusProject(d workspace.Dependencies, value string) (workspace.ProjectRecord, error) {
	project, err := resolveAssignmentProject(d, value)
	var projectError *workspace.ProjectError
	if err != nil && !errors.As(err, &projectError) {
		return workspace.ProjectRecord{}, &agentstatus.Error{Class: agentstatus.ClassIOError, Detail: "could not resolve the project"}
	}
	return project, err
}

func agentStatusInvalidArguments(detail string) error {
	return &agentstatus.Error{Class: agentstatus.ClassInvalidArguments, Detail: detail}
}

func init() {
	RootCmd.AddCommand(newAgentStatusCommand(workspace.SystemDependencies()))
}
