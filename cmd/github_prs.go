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
	c.AddCommand(newGitHubPRsListCommand(env), newGitHubPRsDetailsCommand(env))
	return c
}

type githubPRsList struct {
	Kind              string              `json:"kind"`
	SchemaVersion     int                 `json:"schema_version"`
	Status            string              `json:"status"`
	Error             *githubActionsError `json:"error,omitempty"`
	ObservedAt        string              `json:"observed_at,omitempty"`
	IncompleteReasons *[]string           `json:"incomplete_reasons,omitempty"`
	PullRequests      *[]githubprs.Listed `json:"pull_requests,omitempty"`
}

type githubPRsDetails struct {
	Kind          string                `json:"kind"`
	SchemaVersion int                   `json:"schema_version"`
	Status        string                `json:"status"`
	Error         *githubActionsError   `json:"error,omitempty"`
	ObservedAt    string                `json:"observed_at,omitempty"`
	PullRequests  *[]githubprs.Detailed `json:"pull_requests,omitempty"`
}

func newGitHubPRsListCommand(env githubActionsEnv) *cobra.Command {
	return githubPRsSubcommand("list", "List the open pull requests that are on you, from the searches only",
		"List the open pull requests on you from GitHub's search alone, quickly also with many: id (GitHub's node id),\n"+
			"repository, number, title, url, author, draft, times, reasons and needs_you (review_requested and assigned only;\n"+
			"the reasons on your own pull requests come with ply github prs details). At most 20 seconds; nothing is kept.",
		func(cmd *cobra.Command) (any, string, error) {
			doc := githubPRsList{Kind: "PlyGitHubPullRequestList@1", SchemaVersion: 1, Status: "ok"}
			result, err := env.list(cmd.Context())
			if err != nil {
				doc.Status, doc.Error = "error", githubError(err)
				return doc, githubPRsFailureText(doc.Error), errGitHubPRsFailed
			}
			doc.ObservedAt, doc.IncompleteReasons, doc.PullRequests = githubActionsTime(env.now()), &result.IncompleteReasons, &result.PullRequests
			var b strings.Builder
			for _, pr := range result.PullRequests {
				fmt.Fprintf(&b, "  %s#%d %s\n    %s", pr.Repository, pr.Number, pr.Title, reasonsText(pr.Reasons))
				if pr.Draft {
					b.WriteString("; draft")
				}
				if len(pr.NeedsYou) > 0 {
					b.WriteString("; needs you: " + strings.Join(pr.NeedsYou, ", "))
				}
				b.WriteString("\n    " + pr.URL + "\n")
			}
			return doc, "Pull requests on you\n" + orNone(b.String()) + incompleteText(result.IncompleteReasons), nil
		})
}

func newGitHubPRsDetailsCommand(env githubActionsEnv) *cobra.Command {
	return githubPRsSubcommand("details", "Read the details of the open pull requests that are on you, in batches",
		"Run the searches of ply github prs list and read the details of every hit in batches of at most 20, review threads\n"+
			"only for your own: review_decision, checks, mergeable, merge_state, unresolved_threads and needs_you (changes_requested,\n"+
			"checks_failed, conflict, unresolved_threads and ready_to_merge on your own pull requests). A pull request without\n"+
			"details is unavailable with detail_reason refused (GitHub left it out, for example SAML), gh_failed or timeout (its\n"+
			"batch failed or the deadline passed); the status is still ok. At most 40 seconds; nothing is kept.",
		func(cmd *cobra.Command) (any, string, error) {
			doc := githubPRsDetails{Kind: "PlyGitHubPullRequestDetails@1", SchemaVersion: 1, Status: "ok"}
			rows, err := env.details(cmd.Context())
			if err != nil {
				doc.Status, doc.Error = "error", githubError(err)
				return doc, githubPRsFailureText(doc.Error), errGitHubPRsFailed
			}
			doc.ObservedAt, doc.PullRequests = githubActionsTime(env.now()), &rows
			var b strings.Builder
			for _, pr := range rows {
				fmt.Fprintf(&b, "  %s#%d %s\n", pr.Repository, pr.Number, detailText(pr))
			}
			return doc, "Pull request details\n" + orNone(b.String()), nil
		})
}

// githubPRsSubcommand is one ply github prs subcommand that writes the
// document run gives; a failure document exits with status 1.
func githubPRsSubcommand(use, short, long string, run func(*cobra.Command) (any, string, error)) *cobra.Command {
	var format string
	c := &cobra.Command{Use: use, Short: short, Long: long, Args: cobra.NoArgs,
		Example: "  ply github prs " + use + "\n  ply github prs " + use + " --format json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if format != "text" && format != "json" {
				return errors.New("--format must be text or json")
			}
			doc, text, failed := run(cmd)
			if err := githubActionsWrite(cmd.OutOrStdout(), format, doc, text); err != nil {
				return err
			}
			return failed
		}}
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	return c
}

func detailText(pr githubprs.Detailed) string {
	if pr.DetailReason != nil {
		return "no details (" + *pr.DetailReason + ")"
	}
	var parts []string
	for _, f := range []struct {
		name  string
		value *string
	}{{"review", pr.ReviewDecision}, {"checks", pr.Checks}, {"mergeable", pr.Mergeable}, {"merge state", pr.MergeState}} {
		if f.value != nil {
			parts = append(parts, f.name+" "+*f.value)
		}
	}
	if pr.UnresolvedThreads != nil {
		parts = append(parts, fmt.Sprintf("%d unresolved threads", *pr.UnresolvedThreads))
	}
	if len(pr.NeedsYou) > 0 {
		parts = append(parts, "needs you: "+strings.Join(pr.NeedsYou, ", "))
	}
	return strings.Join(parts, "; ")
}

func githubPRsFailureText(e *githubActionsError) string {
	return fmt.Sprintf("GitHub pull requests: %s (%s)\n", e.Description, e.Code)
}

func reasonsText(reasons []githubprs.Reason) string {
	var why []string
	for _, r := range reasons {
		if r.Via != "" {
			why = append(why, r.Kind+" ("+r.Via+")")
		} else {
			why = append(why, r.Kind)
		}
	}
	return strings.Join(why, ", ")
}

func incompleteText(reasons []string) string {
	if len(reasons) == 0 {
		return ""
	}
	return "Incomplete: " + strings.Join(reasons, ", ") + "\n"
}

func githubPRsText(doc githubPRs) string {
	if doc.Error != nil {
		return githubPRsFailureText(doc.Error)
	}
	var waiting, own strings.Builder
	for _, pr := range *doc.PullRequests {
		line := fmt.Sprintf("  %s#%d %s\n    %s", pr.Repository, pr.Number, pr.Title, reasonsText(pr.Reasons))
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
	return "Waiting for you\n" + orNone(waiting.String()) + "Your pull requests\n" + orNone(own.String()) + incompleteText(*doc.IncompleteReasons)
}

func orNone(lines string) string {
	if lines == "" {
		return "  none\n"
	}
	return lines
}
