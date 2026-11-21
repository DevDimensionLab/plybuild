package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/intelligence"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// Fixed JSON contract names of the intelligence readbacks.
const (
	intelligenceReceiptKind   = "PlyIntelligenceReceipt@1"
	intelligenceListKind      = "PlyIntelligenceList@1"
	intelligencePostKind      = "PlyIntelligencePost@1"
	intelligenceSchemaVersion = 1
)

// maxIntelligenceInput bounds a draft read from a file or stdin.
const maxIntelligenceInput = 8 << 20

func newIntelligenceCommand(d workspace.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:   "intelligence",
		Short: "Record and read what agents learn in a Project",
		Long: "Agents record experience as immutable posts and read it back only through these commands.\n" +
			"A post's schema is its kind: " + strings.Join(intelligence.Kinds, ", ") + ". draft gives a post with ID, time and origin to fill in; observe validates and stores it.\n" +
			"The same post again is already_recorded; the same ID with other content is a conflict. Posts live in <project wrapper>/data/intelligence/.\n" +
			"Without --project the single Project whose wrapper contains the current directory is used.",
		Example: "  ply intelligence draft --kind observation --repo ply --activity PLAN-77/spec --role author --provider claude > draft.json\n" +
			"  ply intelligence observe --file draft.json\n" +
			"  ply intelligence observations --after obs-<uuid> --format json\n" +
			"  ply intelligence show obs-<uuid>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return intelligenceInvalidArguments("unknown intelligence command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	command.SetFlagErrorFunc(func(*cobra.Command, error) error {
		return intelligenceInvalidArguments("invalid flag")
	})

	var project, format, file string
	var kind, repo, activity, role, provider, session, since, after, linksTo string
	addFlags := func(c *cobra.Command) *cobra.Command {
		c.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		command.AddCommand(c)
		return c
	}
	validate := func(positional int, required ...string) cobra.PositionalArgs {
		return func(cmd *cobra.Command, args []string) error {
			if len(args) != positional {
				return intelligenceInvalidArguments(fmt.Sprintf("expected %d positional arguments, got %d", positional, len(args)))
			}
			if format != "text" && format != "json" {
				return intelligenceInvalidArguments("--format must be text or json")
			}
			for _, name := range required {
				if !cmd.Flags().Changed(name) {
					return intelligenceInvalidArguments("--" + name + " is required")
				}
			}
			return nil
		}
	}
	var resolved workspace.ProjectRecord
	open := func() (intelligence.Store, error) {
		var err error
		if resolved, err = resolveIntelligenceProject(d, project); err != nil {
			return intelligence.Store{}, err
		}
		directory, err := resolved.Data().Intelligence()
		if err != nil {
			return intelligence.Store{}, err
		}
		repos := make([]string, 0, len(resolved.RepoIDs))
		for _, id := range resolved.RepoIDs {
			repos = append(repos, string(id))
		}
		return intelligence.Open(directory, string(resolved.ID), repos), nil
	}

	draftCommand := addFlags(&cobra.Command{
		Use:   "draft --kind <kind> --repo <repo-id> --activity <case> --role <role> --provider <codex|claude|other> [--session <id>]",
		Short: "Print a new post of a kind to fill in",
		Long:  "The ID and recorded_at stay on every retry of the same post. --repo must be a repository of the Project. --activity is the stable logical case, such as PLAN-77/spec or <repo>:<base7>..<head7>.",
		Args:  validate(0, "kind", "repo", "activity", "role", "provider"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := open(); err != nil {
				return err
			}
			known := false
			for _, id := range resolved.RepoIDs {
				known = known || string(id) == repo
			}
			if !known {
				return intelligenceInvalidArguments("--repo must be a repository of Project " + string(resolved.ID))
			}
			origin := intelligence.Origin{Project: string(resolved.ID), Repo: repo, Activity: activity, Role: role, Provider: provider}
			if session != "" {
				origin.Session = &session
			}
			draft, err := intelligence.Draft(kind, origin, time.Now())
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(draft)
			return err
		},
	})
	draftCommand.Flags().StringVar(&kind, "kind", "", "one of "+strings.Join(intelligence.Kinds, ", ")+" (required)")
	draftCommand.Flags().StringVar(&repo, "repo", "", "repository ID the post concerns (required)")
	draftCommand.Flags().StringVar(&activity, "activity", "", "stable logical case (required)")
	draftCommand.Flags().StringVar(&role, "role", "", "your actual role, such as author or reviewer (required)")
	draftCommand.Flags().StringVar(&provider, "provider", "", "codex, claude or other (required)")
	draftCommand.Flags().StringVar(&session, "session", "", "known session ID")

	observeCommand := addFlags(&cobra.Command{
		Use:   "observe --file <draft.json|->",
		Short: "Validate and store a filled post",
		Args:  validate(0, "file"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readIntelligenceInput(cmd.InOrStdin(), file)
			if err != nil {
				return err
			}
			s, err := open()
			if err != nil {
				return err
			}
			receipt, err := s.Observe(raw)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string               `json:"kind"`
					SchemaVersion int                  `json:"schema_version"`
					Project       workspace.ProjectID  `json:"project"`
					Receipt       intelligence.Receipt `json:"receipt"`
				}{intelligenceReceiptKind, intelligenceSchemaVersion, resolved.ID, receipt})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n", receipt.State, receipt.ID, receipt.SHA256)
			return err
		},
	})
	observeCommand.Flags().StringVar(&file, "file", "", "the filled draft, or - for stdin (required)")

	listCommand := addFlags(&cobra.Command{
		Use:   "observations",
		Short: "List posts in recorded order, optionally filtered",
		Long:  "--after gives the posts ordered after a known post, so a new sweep sees only what is new. --links-to takes an obs- or k- ID. JSON gives each whole post with the sha256 of its bytes. Files that cannot be read are listed under invalid.",
		Args:  validate(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			result, err := s.List(intelligence.Filter{Kind: kind, Repo: repo, Role: role, Activity: activity, Since: since, After: after, LinksTo: linksTo})
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					intelligence.ListResult
				}{intelligenceListKind, intelligenceSchemaVersion, resolved.ID, result})
			}
			var out strings.Builder
			for _, item := range result.Items {
				var head struct {
					Origin intelligence.Origin `json:"origin"`
				}
				_ = json.Unmarshal(item.Post, &head)
				fmt.Fprintf(&out, "%s %s %s %s %s %s\n", item.ID, item.RecordedAt, item.Kind, head.Origin.Repo, head.Origin.Role, head.Origin.Activity)
			}
			for _, invalid := range result.Invalid {
				fmt.Fprintf(&out, "invalid %s: %s\n", invalid.File, invalid.Error)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	})
	listCommand.Flags().StringVar(&kind, "kind", "", "only posts of this kind")
	listCommand.Flags().StringVar(&repo, "repo", "", "only posts about this repository")
	listCommand.Flags().StringVar(&role, "role", "", "only posts recorded in this role")
	listCommand.Flags().StringVar(&activity, "activity", "", "only posts of this logical case")
	listCommand.Flags().StringVar(&since, "since", "", "only posts recorded at or after this UTC time")
	listCommand.Flags().StringVar(&after, "after", "", "only posts ordered after this obs- ID")
	listCommand.Flags().StringVar(&linksTo, "links-to", "", "only posts linking to this obs- or k- ID")

	addFlags(&cobra.Command{
		Use:   "show <obs-id>",
		Short: "Show one post with the sha256 of its bytes",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			item, err := s.Show(args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					Item          intelligence.Item   `json:"item"`
				}{intelligencePostKind, intelligenceSchemaVersion, resolved.ID, item})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", item.Post, item.SHA256)
			return err
		},
	})
	command.AddCommand(newIntelligenceKnowledgeCommand(d))
	return command
}

func readIntelligenceInput(stdin io.Reader, file string) ([]byte, error) {
	reader := stdin
	if file != "-" {
		opened, err := os.Open(file)
		if err != nil {
			return nil, intelligenceInvalidArguments("cannot read --file")
		}
		defer opened.Close()
		reader = opened
	}
	raw, err := io.ReadAll(io.LimitReader(reader, maxIntelligenceInput+1))
	if err != nil {
		return nil, intelligenceInvalidArguments("cannot read --file")
	}
	if len(raw) > maxIntelligenceInput {
		return nil, intelligenceInvalidArguments("the draft is larger than 8 MiB")
	}
	return raw, nil
}

// resolveIntelligenceProject maps Project lookup failures to fixed details
// without paths, as resolveAgentStatusProject does.
func resolveIntelligenceProject(d workspace.Dependencies, value string) (workspace.ProjectRecord, error) {
	project, err := resolveAssignmentProject(d, value)
	var projectError *workspace.ProjectError
	if err != nil && !errors.As(err, &projectError) {
		return workspace.ProjectRecord{}, &intelligence.Error{Class: intelligence.ClassIOError, Detail: "could not resolve the project"}
	}
	return project, err
}

func intelligenceInvalidArguments(detail string) error {
	return &intelligence.Error{Class: intelligence.ClassInvalidArguments, Detail: detail}
}

func init() {
	RootCmd.AddCommand(newIntelligenceCommand(workspace.SystemDependencies()))
}
