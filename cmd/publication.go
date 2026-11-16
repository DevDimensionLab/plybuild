package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/publication"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newPublicationCommand(d workspace.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:   "publication",
		Short: "Publish, list, show and withdraw runs that other applications deliver through Ply",
		Long: "A run is a finished result of an application, such as a code map or a workflow view, stored unchanged with publication.json written by Ply.\n" +
			"Runs live in <project wrapper>/data/publications/,\n" +
			"as <repo-id or _project>/<category>/[<agent>/]<run>/. A published run is never changed; Ply keeps every run with --plan and the newest five without plan per path.\n" +
			"Withdraw hides the runs of a path from list until the next publish; it removes nothing.\n" +
			"Without --project the single Project whose wrapper contains the current directory is used.",
		Example: "  ply publication publish --project ply --repo ply-code-map --category code-map --dir out --entry index.html --title 'Kodekart'\n" +
			"  ply publication publish --project ply --category workflow --agent ply-epic --plan PLAN-49 --dir out --entry index.html --title 'PLAN-49'\n" +
			"  ply publication list --project ply --plan PLAN-49 --format json\n" +
			"  ply publication show --project ply --repo ply-code-map --category code-map --format json\n" +
			"  ply publication withdraw --project ply --repo ply-planning --category code-health",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return publicationInvalidArguments("unknown publication command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	command.SetFlagErrorFunc(func(*cobra.Command, error) error {
		return publicationInvalidArguments("invalid flag")
	})

	var project, format string
	var request publication.Request
	var filter publication.Filter
	var run string
	addFlags := func(c *cobra.Command) *cobra.Command {
		c.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		c.Args = func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return publicationInvalidArguments("expected no positional arguments")
			}
			if format != "text" && format != "json" {
				return publicationInvalidArguments("--format must be text or json")
			}
			return nil
		}
		command.AddCommand(c)
		return c
	}
	var resolved workspace.ProjectRecord
	open := func() (publication.Store, error) {
		var err error
		if resolved, err = resolvePublicationProject(d, project); err != nil {
			return publication.Store{}, err
		}
		directory, err := resolved.Data().Publications()
		if err != nil {
			return publication.Store{}, err
		}
		return publication.Open(directory), nil
	}
	readback := func(cmd *cobra.Command, v publication.View) error {
		if format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), v)
		}
		_, err := io.WriteString(cmd.OutOrStdout(), publicationText(v))
		return err
	}

	publishCommand := addFlags(&cobra.Command{
		Use:   "publish --category <category> --dir <dir> --entry <file> --title <title>",
		Short: "Store the files of a directory as a new run",
		Long: "Copy every regular file under --dir as a new run of the repo (--repo) or the whole Project (no --repo), the category and optionally the agent.\n" +
			"--entry names the file to open; --plan links the run to a plan of the Project, which must exist; --source is a short producer ID such as a commit.\n" +
			"Symlinks, publication.json at the top, more than 1000 files and more than 20 MiB are refused before anything is stored. The actor is --actor or PLY_ACTOR.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			request.Project = string(resolved.ID)
			if request.Repo != "" && !slices.Contains(resolved.RepoIDs, workspace.RepoID(request.Repo)) {
				return publicationInvalidArguments("--repo must be a repo ID of project " + request.Project)
			}
			if request.Plan != "" {
				if err := checkPublicationPlan(resolved, request.Plan); err != nil {
					return err
				}
			}
			if request.Actor == "" {
				request.Actor = os.Getenv("PLY_ACTOR")
			}
			v, err := s.Publish(request)
			if v.Run == "" {
				return err
			}
			if writeErr := readback(cmd, v); writeErr != nil && err == nil {
				err = &publication.Error{Class: publication.ClassIOError, Detail: fmt.Sprintf("run %s was stored; writing its readback failed", v.Run)}
			}
			return err
		},
	})
	publishCommand.Flags().StringVar(&request.Repo, "repo", "", "repo ID of the Project (default: the whole Project)")
	publishCommand.Flags().StringVar(&request.Category, "category", "", "category such as code-map or workflow (required)")
	publishCommand.Flags().StringVar(&request.Agent, "agent", "", "agent that performed the work")
	publishCommand.Flags().StringVar(&request.Plan, "plan", "", "PLAN-<n> of the Project the run belongs to")
	publishCommand.Flags().StringVar(&request.Dir, "dir", "", "directory with the files to publish (required)")
	publishCommand.Flags().StringVar(&request.Entry, "entry", "", "file to open, relative to --dir (required)")
	publishCommand.Flags().StringVar(&request.Title, "title", "", "one line (required)")
	publishCommand.Flags().StringVar(&request.Source, "source", "", "short producer ID such as a commit")
	publishCommand.Flags().StringVar(&request.Actor, "actor", "", "who publishes (default: PLY_ACTOR)")

	listCommand := addFlags(&cobra.Command{
		Use:   "list",
		Short: "List the newest run of every path, newest first",
		Long: "Filters select runs before the newest of each path is taken; --all lists every kept run.\n" +
			"Runs whose publication.json cannot be read are listed under invalid.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			result, err := s.List(filter)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					publication.ListResult
				}{publication.ListKind, publication.SchemaVersion, resolved.ID, result})
			}
			_, err = io.WriteString(cmd.OutOrStdout(), publicationListText(result))
			return err
		},
	})
	listCommand.Flags().StringVar(&filter.Repo, "repo", "", "only runs of this repo ID")
	listCommand.Flags().StringVar(&filter.Category, "category", "", "only runs of this category")
	listCommand.Flags().StringVar(&filter.Agent, "agent", "", "only runs of this agent")
	listCommand.Flags().StringVar(&filter.Plan, "plan", "", "only runs of this plan")
	listCommand.Flags().BoolVar(&filter.All, "all", false, "every kept run, not only the newest per path")

	var showRepo, showCategory, showAgent string
	showCommand := addFlags(&cobra.Command{
		Use:   "show --category <category> [--run <run>]",
		Short: "Show one run with the absolute path of its directory",
		Long:  "Without --repo the run of the whole Project is shown; without --run the newest run of the path.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			v, err := s.Show(showRepo, showCategory, showAgent, run)
			if err != nil {
				return err
			}
			return readback(cmd, v)
		},
	})
	showCommand.Flags().StringVar(&showRepo, "repo", "", "repo ID (default: the whole Project)")
	showCommand.Flags().StringVar(&showCategory, "category", "", "category (required)")
	showCommand.Flags().StringVar(&showAgent, "agent", "", "agent")
	showCommand.Flags().StringVar(&run, "run", "", "run ID (default: the newest)")

	var withdrawRepo, withdrawCategory, withdrawAgent, withdrawActor string
	withdrawCommand := addFlags(&cobra.Command{
		Use:   "withdraw --category <category>",
		Short: "Hide the runs of a path from list and show until the next publish",
		Long: "Store withdrawn.json in the path with its newest run: list and show without --run leave out that run and every older one, and a later publish is shown again.\n" +
			"No run is removed and show --run still reads a withdrawn run. Without --repo the path of the whole Project is withdrawn; without --agent only the path without agent.\n" +
			"Withdrawing a path whose newest run is already withdrawn returns the stored withdrawal unchanged. The actor is --actor or PLY_ACTOR.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			if withdrawRepo != "" && !slices.Contains(resolved.RepoIDs, workspace.RepoID(withdrawRepo)) {
				return publicationInvalidArguments("--repo must be a repo ID of project " + string(resolved.ID))
			}
			if withdrawActor == "" {
				withdrawActor = os.Getenv("PLY_ACTOR")
			}
			w, err := s.Withdraw(string(resolved.ID), withdrawRepo, withdrawCategory, withdrawAgent, withdrawActor)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), w)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), withdrawalText(w))
			return err
		},
	})
	withdrawCommand.Flags().StringVar(&withdrawRepo, "repo", "", "repo ID of the Project (default: the whole Project)")
	withdrawCommand.Flags().StringVar(&withdrawCategory, "category", "", "category (required)")
	withdrawCommand.Flags().StringVar(&withdrawAgent, "agent", "", "agent")
	withdrawCommand.Flags().StringVar(&withdrawActor, "actor", "", "who withdraws (default: PLY_ACTOR)")
	return command
}

// resolvePublicationProject maps Project lookup failures to fixed details
// without paths, as resolvePlanProject does.
func resolvePublicationProject(d workspace.Dependencies, value string) (workspace.ProjectRecord, error) {
	project, err := resolveAssignmentProject(d, value)
	var projectError *workspace.ProjectError
	if err != nil && !errors.As(err, &projectError) {
		return workspace.ProjectRecord{}, &publication.Error{Class: publication.ClassIOError, Detail: "could not resolve the project"}
	}
	return project, err
}

func checkPublicationPlan(project workspace.ProjectRecord, id string) error {
	directory, err := project.Data().Plans()
	if err != nil {
		return err
	}
	_, err = plan.Open(directory).Show(id, 0)
	var planError *plan.Error
	switch {
	case errors.As(err, &planError) && (planError.Class == plan.ClassNotFound || planError.Class == plan.ClassInvalidArguments):
		return publicationInvalidArguments("--plan " + id + " is not a plan of project " + string(project.ID))
	case err != nil:
		return &publication.Error{Class: publication.ClassIOError, Detail: "could not read plan " + id}
	}
	return nil
}

func publicationSubject(v publication.View) string {
	subject := v.Repo
	if subject == "" {
		subject = publication.ProjectSegment
	}
	subject += "/" + v.Category
	if v.Agent != "" {
		subject += "/" + v.Agent
	}
	return subject
}

func publicationText(v publication.View) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s %s %s\n", publicationSubject(v), v.Run, v.Title)
	fmt.Fprintf(&out, "published: %s", v.PublishedAt)
	if v.Actor != "" {
		fmt.Fprintf(&out, " by %s", v.Actor)
	}
	out.WriteString("\n")
	writeField(&out, "", "plan", v.Plan)
	writeField(&out, "", "source", v.Source)
	fmt.Fprintf(&out, "entry: %s (%d files)\npath: %s\n", v.Entry, len(v.Files), v.Path)
	return out.String()
}

func withdrawalText(w publication.Withdrawal) string {
	subject := publicationSubject(publication.View{Publication: publication.Publication{Repo: w.Repo, Category: w.Category, Agent: w.Agent}})
	text := fmt.Sprintf("%s withdrawn through %s at %s", subject, w.Through, w.WithdrawnAt)
	if w.Actor != "" {
		text += " by " + w.Actor
	}
	return text + "\n"
}

func publicationListText(result publication.ListResult) string {
	var out strings.Builder
	if len(result.Items) == 0 {
		out.WriteString("no publications\n")
	}
	for _, v := range result.Items {
		fmt.Fprintf(&out, "%s %s %s\n", publicationSubject(v), v.Run, v.Title)
	}
	for _, invalid := range result.Invalid {
		fmt.Fprintf(&out, "warning: invalid run %s: %s\n", invalid.Run, invalid.Error)
	}
	return out.String()
}

func publicationInvalidArguments(detail string) error {
	return &publication.Error{Class: publication.ClassInvalidArguments, Detail: detail}
}

func init() {
	RootCmd.AddCommand(newPublicationCommand(workspace.SystemDependencies()))
}
