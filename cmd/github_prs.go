package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/githubprs"
	"github.com/spf13/cobra"
)

var errGitHubPRsFailed = errors.New("reading GitHub pull requests failed")

type githubPRs struct {
	Kind              string                   `json:"kind"`
	SchemaVersion     int                      `json:"schema_version"`
	Status            string                   `json:"status"`
	Error             *githubActionsError      `json:"error,omitempty"`
	ObservedAt        string                   `json:"observed_at,omitempty"`
	Viewer            string                   `json:"viewer,omitempty"`
	IncompleteReasons *[]string                `json:"incomplete_reasons,omitempty"`
	PullRequests      *[]githubprs.PullRequest `json:"pull_requests,omitempty"`
}

func newGitHubPRsCommand(env githubActionsEnv) *cobra.Command {
	var format string
	c := &cobra.Command{
		Use:   "prs",
		Short: "List the open pull requests that are on you, and what waits for you",
		Long: "List the open pull requests where review is requested from you or a team you are on, that are assigned to you,\n" +
			"or that you made, as gh is logged in. Ply gets what gh gets: a pull request GitHub gives no details for, for example\n" +
			"in an organization with SAML the gh token is not authorized for, is listed with detail unavailable.\n" +
			"needs_you: review_requested, assigned, and on your own pull requests changes_requested, checks_failed, conflict,\n" +
			"unresolved_threads and ready_to_merge (GitHub's merge state CLEAN). The read keeps nothing and takes at most 40 seconds.\n" +
			"The exit status is 1 when GitHub could not be read.",
		Example: "  ply github prs\n  ply github prs --format json",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("expected no positional arguments")
			}
			if format != "text" && format != "json" {
				return errors.New("--format must be text or json")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			doc := githubPRs{Kind: "PlyGitHubPullRequests@1", SchemaVersion: 1, Status: "ok"}
			result, err := env.prs(cmd.Context())
			var failed error
			if err != nil {
				doc.Status, doc.Error, failed = "error", githubError(err), errGitHubPRsFailed
			} else {
				doc.ObservedAt, doc.Viewer = githubActionsTime(env.now()), result.Viewer
				doc.IncompleteReasons, doc.PullRequests = &result.IncompleteReasons, &result.PullRequests
			}
			if err := githubActionsWrite(cmd.OutOrStdout(), format, doc, githubPRsText(doc)); err != nil {
				return err
			}
			return failed
		},
	}
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	return c
}

func githubPRsText(doc githubPRs) string {
	if doc.Error != nil {
		return fmt.Sprintf("GitHub pull requests: %s (%s)\n", doc.Error.Description, doc.Error.Code)
	}
	var waiting, own strings.Builder
	for _, pr := range *doc.PullRequests {
		var why []string
		for _, r := range pr.Reasons {
			if r.Via != "" {
				why = append(why, r.Kind+" ("+r.Via+")")
			} else {
				why = append(why, r.Kind)
			}
		}
		line := fmt.Sprintf("  %s#%d %s\n    %s", pr.Repository, pr.Number, pr.Title, strings.Join(why, ", "))
		if pr.Draft {
			line += "; draft"
		}
		if pr.Detail != "available" {
			line += "; no details from GitHub"
		}
		if len(pr.NeedsYou) > 0 {
			line += "; needs you: " + strings.Join(pr.NeedsYou, ", ")
		}
		line += "\n    " + pr.URL + "\n"
		if len(pr.NeedsYou) > 0 {
			waiting.WriteString(line)
		} else {
			own.WriteString(line)
		}
	}
	text := "Waiting for you\n" + orNone(waiting.String()) + "Your pull requests\n" + orNone(own.String())
	if len(*doc.IncompleteReasons) > 0 {
		text += "Incomplete: " + strings.Join(*doc.IncompleteReasons, ", ") + "\n"
	}
	return text
}

func orNone(lines string) string {
	if lines == "" {
		return "  none\n"
	}
	return lines
}
