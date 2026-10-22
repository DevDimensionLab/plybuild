package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceExternalIssuesCommand(d workspace.Dependencies, kind string) *cobra.Command {
	group := &cobra.Command{Use: "external-issue", Short: "Link offline Jira or GitHub issues", Long: "Manage optional offline supplied issue metadata. No provider request or Git change occurs. Identities are unique across all Epics and Tasks in this workspace. No reference is implicitly primary. Primary changes and unlink require the current shared metadata revision.", Example: "  ply workspace " + kind + " external-issue add an-item --provider jira --key PLY-123 --url https://jira.example.invalid/browse/PLY-123 --actor human\n  ply workspace " + kind + " external-issue show an-item --format json\n  ply workspace " + kind + " external-issue primary an-item ext_ID --expected-revision 1 --actor human\n  ply workspace " + kind + " external-issue unlink an-item ext_ID --expected-revision 2 --actor human"}
	var showFormat string
	show := &cobra.Command{Use: "show <" + kind + "-id>", Aliases: []string{"list"}, Short: "Read refs, primary, effective origin and history offline", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateWorkFormat(showFormat); err != nil {
			return err
		}
		out, err := workspace.ShowExternalIssues(d, kind, args[0])
		if err != nil {
			return err
		}
		if showFormat == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
		}
		return renderExternalIssues(cmd.OutOrStdout(), out)
	}}
	show.Flags().StringVar(&showFormat, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	group.AddCommand(show)
	for _, operation := range []string{workspace.ExternalIssueAdd, workspace.ExternalIssuePrimary, workspace.ExternalIssueUnlink} {
		operation := operation
		var provider, host, key, owner, repo, url, actor, format string
		var number int64
		var primary, clear bool
		var revision int
		use := operation + " <" + kind + "-id>"
		if operation != workspace.ExternalIssueAdd {
			use += " [ref-id]"
		}
		command := &cobra.Command{Use: use, Short: map[string]string{workspace.ExternalIssueAdd: "Add a typed offline reference (identical retry writes nothing)", workspace.ExternalIssuePrimary: "Select or explicitly clear the primary commit reference", workspace.ExternalIssueUnlink: "Unlink exactly one reference and clear its primary selection"}[operation], Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 || operation == workspace.ExternalIssueAdd && len(args) != 1 {
				return workspace.WorkInvalidArguments("expected work-item ID and, for primary/unlink, the selected ref ID")
			}
			if operation == workspace.ExternalIssuePrimary {
				if (clear && len(args) != 1) || (!clear && len(args) != 2) {
					return workspace.WorkInvalidArguments("primary requires a ref ID or --clear")
				}
			}
			if operation == workspace.ExternalIssueUnlink && len(args) != 2 {
				return workspace.WorkInvalidArguments("unlink requires the exact ref ID")
			}
			return validateWorkFormat(format)
		}, RunE: func(cmd *cobra.Command, args []string) error {
			in := workspace.ExternalIssueInput{SubjectKind: kind, SubjectID: args[0], Operation: operation, ActorClaim: actor, Primary: primary}
			if len(args) == 2 {
				in.RefID = args[1]
			}
			if cmd.Flags().Changed("expected-revision") {
				in.ExpectedRevision = &revision
			}
			if operation == workspace.ExternalIssueAdd {
				in.Issue = &workspace.ExternalIssue{Provider: workspace.ExternalProvider(provider), Host: host, Key: key, Owner: owner, Repository: repo, Number: number, URL: url}
			}
			out, err := workspace.UpdateExternalIssues(d, in)
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Changed: %t\n", out.Changed); err != nil {
				return err
			}
			return renderExternalIssues(cmd.OutOrStdout(), out.ExternalIssuesReadback)
		}}
		command.Flags().StringVar(&actor, "actor", "", "recorder identity or actor claim")
		_ = command.MarkFlagRequired("actor")
		command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		command.Flags().IntVar(&revision, "expected-revision", 0, "require this shared workspace external-issue revision")
		if operation == workspace.ExternalIssueAdd {
			command.Flags().StringVar(&provider, "provider", "", "jira or github")
			command.Flags().StringVar(&host, "host", "", "normalized URL authority (required for GitHub)")
			command.Flags().StringVar(&key, "key", "", "Jira issue key, for example PLY-123")
			command.Flags().StringVar(&owner, "owner", "", "GitHub owner")
			command.Flags().StringVar(&repo, "repository", "", "GitHub repository")
			command.Flags().Int64Var(&number, "number", 0, "positive GitHub issue number")
			command.Flags().StringVar(&url, "url", "", "absolute http or https issue URL without credentials, query or fragment")
			command.Flags().BoolVar(&primary, "primary", false, "explicitly select this ref as primary (requires expected revision)")
			_ = command.MarkFlagRequired("provider")
			_ = command.MarkFlagRequired("url")
		}
		if operation == workspace.ExternalIssuePrimary {
			command.Flags().BoolVar(&clear, "clear", false, "clear the explicit primary, exposing any Epic fallback")
		}
		setWorkFlagErrors(command)
		group.AddCommand(command)
	}
	return group
}
func renderExternalIssues(w io.Writer, out workspace.ExternalIssuesReadback) error {
	if _, err := fmt.Fprintf(w, "External issues for %s %s (workspace revision %d; %s):\n", out.SubjectKind, out.SubjectID, out.Revision, out.ProviderStatus); err != nil {
		return err
	}
	for _, ref := range out.Refs {
		if _, err := fmt.Fprintf(w, "  %s %s %s %s\n", ref.ID, ref.Provider, ref.CommitKey(), ref.URL); err != nil {
			return err
		}
	}
	primary := "none"
	if out.Primary != nil {
		primary = out.Primary.CommitKey
	}
	effective := "none"
	if out.Effective != nil {
		effective = fmt.Sprintf("%s from %s %s (%s)", out.Effective.CommitKey, out.Effective.SubjectKind, out.Effective.SubjectID, out.Effective.RefID)
	}
	_, err := fmt.Fprintf(w, "Primary: %s\nEffective commit reference: %s\n", primary, effective)
	return err
}
