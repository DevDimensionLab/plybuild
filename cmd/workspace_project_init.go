package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceProjectInitCommand(d workspace.Dependencies) *cobra.Command {
	var in workspace.ProjectInitInput
	var format string
	c := &cobra.Command{
		Use: "init <project-id>", Short: "Preview direct-child repository registration, then explicitly apply",
		Long:    "Discover only immediate physical Git worktree children of the wrapper (default cwd). Group exact Git common directories, preferring the ordinary root, then lexical path. One repository uses the Project ID; multiple repositories use representative directory names without slug conversion. Name defaults to the existing Project name or Project ID. Symlinks and nested repository-wrapper/main layouts are omitted. Preview is read-only. Apply requires its exact --confirmation digest and uses project add; existing membership is never updated implicitly. Use project add/add-member for explicit layouts.",
		Example: "  ply workspace project init example --format json\n  ply workspace project init example --apply --confirmation sha256:<preview-digest>",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.ProjectInvalidArguments("expected exactly one project ID")
			}
			if _, err := workspace.ParseProjectID(args[0]); err != nil {
				return err
			}
			return validateProjectFormat(format)
		},
		RunE: func(c *cobra.Command, args []string) error {
			in.ProjectID = workspace.ProjectID(args[0])
			plan, err := workspace.InitProject(d, in)
			// Blocked plans retain concrete candidates and omissions in both formats.
			if err != nil && plan.Wrapper == "" {
				return err
			}
			if renderErr := renderProjectInit(c, format, plan); renderErr != nil {
				return renderErr
			}
			return err
		},
	}
	c.Flags().StringVarP(&in.Wrapper, "wrapper", "w", "", "existing wrapper directory (default cwd)")
	c.Flags().StringVarP(&in.Name, "name", "n", "", "display name (existing name or project ID)")
	c.Flags().BoolVar(&in.Apply, "apply", false, "register the exact confirmed ProjectAdd plan")
	c.Flags().StringVar(&in.Confirmation, "confirmation", "", "exact sha256: preview confirmation; requires --apply")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setProjectFlagErrors(c)
	return c
}

func renderProjectInit(c *cobra.Command, format string, p workspace.ProjectInitPlan) error {
	if format == "json" {
		return json.NewEncoder(c.OutOrStdout()).Encode(p)
	}
	out := c.OutOrStdout()
	action := "Preview"
	if p.Applied {
		action = "Registered"
	}
	if _, err := fmt.Fprintf(out, "%s project %s (%s).\nWrapper: %s\nRepositories:\n", action, p.ProjectID, p.Name, p.Wrapper); err != nil {
		return err
	}
	for _, repo := range p.Repositories {
		if _, err := fmt.Fprintf(out, "  %s:\n    locator: %s\n    git common directory: %s\n", repo.ID, repo.Locator, repo.GitCommonDir); err != nil {
			return err
		}
	}
	for _, omit := range p.Omissions {
		if _, err := fmt.Fprintf(out, "Omitted %s: %s\n", omit.Path, omit.Reason); err != nil {
			return err
		}
	}
	for _, blocker := range p.Blockers {
		if _, err := fmt.Fprintf(out, "Blocked: %s\n", blocker); err != nil {
			return err
		}
	}
	if p.Confirmation != "" {
		if _, err := fmt.Fprintf(out, "To apply, repeat these options with --apply --confirmation %s\n", p.Confirmation); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(out, p.NextAction)
	return err
}
