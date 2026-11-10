package cmd

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/assignment"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// assignmentStdinLimit reads one byte past the store's 64 KiB text limit so
// the store, not the CLI, rejects oversized text.
const assignmentStdinLimit = 64*1024 + 1

func newAssignmentCommand(d workspace.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:   "assignment",
		Short: "Add, show, list, take, release and finish the Assignments of a Project",
		Long: "An Assignment is a short work description stored in a Project of the containing Ply workspace: ASSIGNMENT-<n>, a title, text read from stdin, one or more repo IDs of the Project and optional topics. It is never changed after add.\n" +
			"Assignments live as one JSON file each in the Project's data_paths.assignments in ply-workspace.json, relative to the Project wrapper (default <project wrapper>/assignments/).\n" +
			"Without --project the single Project whose wrapper contains the current directory is used.\n" +
			"An agent takes an Assignment before working on it and marks it done when finished; done puts it in the archive view of list. The original file is never changed: the lifecycle is kept in data_paths.assignments/.lifecycle/.\n" +
			"An agent session starts on one with the ply-assignment skill: /ply-assignment ASSIGNMENT-<n>.",
		Example: "  printf '%s\\n' 'Vis repo-tags i backlog list.' | ply assignment add --title 'Vis repo-tags' --repo ply --topic backlog\n" +
			"  ply assignment show ASSIGNMENT-7 --format json\n" +
			"  ply assignment take ASSIGNMENT-7 --revision 0 --as claude:<session-id>\n" +
			"  ply assignment done ASSIGNMENT-7 --claim <claim>\n" +
			"  ply assignment list --project ply --view archive",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return assignmentInvalidArguments("unknown assignment command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	command.SetFlagErrorFunc(func(*cobra.Command, error) error {
		return assignmentInvalidArguments("invalid flag")
	})

	var project, format, title, actor, claim, reason, view string
	var repos, topics []string
	var revision int
	var revoke bool
	addFlags := func(c *cobra.Command) *cobra.Command {
		c.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		command.AddCommand(c)
		return c
	}
	validate := func(positional int) cobra.PositionalArgs {
		return func(cmd *cobra.Command, args []string) error {
			if len(args) != positional {
				return assignmentInvalidArguments(fmt.Sprintf("expected %d positional arguments, got %d", positional, len(args)))
			}
			if format != "text" && format != "json" {
				return assignmentInvalidArguments("--format must be text or json")
			}
			return nil
		}
	}
	var resolved workspace.ProjectRecord
	open := func() (assignment.Store, error) {
		var err error
		if resolved, err = resolveAssignmentProject(d, project); err != nil {
			return assignment.Store{}, err
		}
		directory, err := workspace.ProjectDataPath(resolved, assignment.DataPathCategory)
		if err != nil {
			return assignment.Store{}, err
		}
		return assignment.Open(directory), nil
	}
	render := func(cmd *cobra.Command, a assignment.Assignment, text func(io.Writer, assignment.Assignment) error) error {
		if format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), struct {
				Project    workspace.ProjectID   `json:"project"`
				Assignment assignment.Assignment `json:"assignment"`
			}{resolved.ID, a})
		}
		return text(cmd.OutOrStdout(), a)
	}

	addCommand := addFlags(&cobra.Command{
		Use:   "add --title <title> --repo <repo-id>... [--topic <topic>]...",
		Short: "Add an Assignment with text from stdin",
		Long: "Add an Assignment with the text read byte for byte from stdin (non-blank UTF-8, at most 64 KiB). Each --repo must be a repo ID of the Project. It prints the new ID.\n" +
			"If the answer is lost, or add fails after the file may have been written, run list and look up the title before trying again.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(0)(cmd, args); err != nil {
				return err
			}
			if strings.TrimSpace(title) == "" {
				return assignmentInvalidArguments("--title is required")
			}
			if len(repos) == 0 {
				return assignmentInvalidArguments("--repo is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			for _, repo := range repos {
				if !slices.Contains(resolved.RepoIDs, workspace.RepoID(strings.TrimSpace(repo))) {
					return assignmentInvalidArguments("--repo must be a repo ID of project " + string(resolved.ID))
				}
			}
			text, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), assignmentStdinLimit))
			if err != nil {
				return &assignment.Error{Class: assignment.ClassIOError, Detail: "read stdin"}
			}
			added, err := store.Add(title, string(text), repos, topics)
			if err != nil {
				return err
			}
			return render(cmd, added, func(out io.Writer, a assignment.Assignment) error {
				_, err := fmt.Fprintln(out, a.ID)
				return err
			})
		},
	})
	addCommand.Flags().StringVar(&title, "title", "", "one-line title")
	addCommand.Flags().StringArrayVar(&repos, "repo", nil, "repo ID of the Project the Assignment concerns (repeatable, at least one)")
	addCommand.Flags().StringArrayVar(&topics, "topic", nil, "topic (repeatable)")

	addFlags(&cobra.Command{
		Use:   "show <ASSIGNMENT-n>",
		Short: "Show one Assignment with its text",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			a, err := store.Read(args[0])
			if err != nil {
				return err
			}
			lifecycle, err := store.Lifecycle(args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Project    workspace.ProjectID   `json:"project"`
					Assignment assignment.Assignment `json:"assignment"`
					Lifecycle  assignment.Lifecycle  `json:"lifecycle"`
				}{resolved.ID, a, lifecycle})
			}
			return writeAssignmentText(cmd.OutOrStdout(), a, lifecycle)
		},
	})

	// transition runs one lifecycle change and prints "<id> <event> revision <n>".
	transition := func(cmd *cobra.Command, id string, change func(assignment.Store) (assignment.Lifecycle, error)) error {
		store, err := open()
		if err != nil {
			return err
		}
		lifecycle, err := change(store)
		if err != nil {
			return err
		}
		if format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), struct {
				Project   workspace.ProjectID  `json:"project"`
				ID        string               `json:"id"`
				Lifecycle assignment.Lifecycle `json:"lifecycle"`
			}{resolved.ID, id, lifecycle})
		}
		last := lifecycle.Events[len(lifecycle.Events)-1]
		var out string
		if last.Kind == assignment.EventTaken {
			out = fmt.Sprintf("%s taken claim %s revision %d\n", id, last.Claim, lifecycle.Revision)
		} else {
			out = fmt.Sprintf("%s %s revision %d\n", id, last.Kind, lifecycle.Revision)
		}
		_, err = io.WriteString(cmd.OutOrStdout(), out)
		return err
	}
	required := func(cmd *cobra.Command, names ...string) error {
		for _, name := range names {
			if !cmd.Flags().Changed(name) {
				return assignmentInvalidArguments("--" + name + " is required")
			}
		}
		return nil
	}

	takeCommand := addFlags(&cobra.Command{
		Use:   "take <ASSIGNMENT-n> --revision <r> --as <actor>",
		Short: "Take an available Assignment before working on it",
		Long: "Take an available Assignment for one agent. --revision is the lifecycle revision just read with show; a take with an older revision is refused, so of two simultaneous takes exactly one succeeds. --as is a one-line display label of at most 100 characters, such as claude:<session-id>.\n" +
			"It prints the new claim, which release and done need. If the answer is lost, run show: the holder label says whether the take succeeded.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			return required(cmd, "revision", "as")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return transition(cmd, args[0], func(store assignment.Store) (assignment.Lifecycle, error) {
				return store.Take(args[0], revision, actor)
			})
		},
	})
	takeCommand.Flags().IntVar(&revision, "revision", 0, "lifecycle revision shown by show (required)")
	takeCommand.Flags().StringVar(&actor, "as", "", "display label of the agent taking it (required)")

	releaseCommand := addFlags(&cobra.Command{
		Use:   "release <ASSIGNMENT-n> (--claim <claim> | --revoke --revision <r> --as <actor> --reason <text>)",
		Short: "Give a taken Assignment back unfinished",
		Long: "With --claim the holder gives the Assignment back; it becomes available again. Repeating it right after it succeeded changes nothing.\n" +
			"With --revoke a human recovers an Assignment whose holder is gone: --revision is the revision shown by show, and --reason is one line of at most 200 characters. Agents do not revoke.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			if cmd.Flags().Changed("claim") == revoke {
				return assignmentInvalidArguments("exactly one of --claim and --revoke is required")
			}
			if !revoke {
				for _, name := range []string{"revision", "as", "reason"} {
					if cmd.Flags().Changed(name) {
						return assignmentInvalidArguments("--revision, --as and --reason require --revoke")
					}
				}
				return nil
			}
			return required(cmd, "revision", "as", "reason")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return transition(cmd, args[0], func(store assignment.Store) (assignment.Lifecycle, error) {
				if revoke {
					return store.Revoke(args[0], revision, actor, reason)
				}
				return store.Release(args[0], claim)
			})
		},
	})
	releaseCommand.Flags().StringVar(&claim, "claim", "", "claim printed by take")
	releaseCommand.Flags().BoolVar(&revoke, "revoke", false, "revoke the current claim (human recovery)")
	releaseCommand.Flags().IntVar(&revision, "revision", 0, "lifecycle revision shown by show (with --revoke)")
	releaseCommand.Flags().StringVar(&actor, "as", "", "display label of who revokes (with --revoke)")
	releaseCommand.Flags().StringVar(&reason, "reason", "", "why the claim is revoked (with --revoke)")

	doneCommand := addFlags(&cobra.Command{
		Use:   "done <ASSIGNMENT-n> --claim <claim>",
		Short: "Mark a taken Assignment done; it moves to the archive view",
		Long: "The holder marks the work done. In the same atomic write the Assignment leaves the available and taken views and appears in list --view archive. Done is final; a follow-up is a new Assignment.\n" +
			"Repeating it with the same claim changes nothing. On an uncertain result (result unknown, run show) run show before anything else.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			return required(cmd, "claim")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return transition(cmd, args[0], func(store assignment.Store) (assignment.Lifecycle, error) {
				return store.Done(args[0], claim)
			})
		},
	})
	doneCommand.Flags().StringVar(&claim, "claim", "", "claim printed by take (required)")

	listCommand := addFlags(&cobra.Command{
		Use:   "list [--view available|taken|archive|all]",
		Short: "List Assignments by number, and files that cannot be read",
		Long: "List the Assignments of one view: available (default) are not taken and not done, taken are being worked on, archive are done, all is every view. Files that cannot be read are always listed under invalid.\n" +
			"Text lines are <id> <created_at> <repos> <title>; other views than available add <state> <taken_at|done_at|-> before the title.",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(0)(cmd, args); err != nil {
				return err
			}
			if !slices.Contains(assignment.Views, assignment.View(view)) {
				return assignmentInvalidArguments("--view must be available, taken, archive or all")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			result, err := store.ListView(assignment.View(view))
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Project workspace.ProjectID `json:"project"`
					assignment.ViewResult
				}{resolved.ID, result})
			}
			var out strings.Builder
			for _, item := range result.Items {
				if result.View == assignment.ViewAvailable {
					fmt.Fprintf(&out, "%s %s %s %s\n", item.ID, item.CreatedAt, strings.Join(item.Repos, ","), item.Title)
					continue
				}
				at := cmp.Or(item.TakenAt, item.DoneAt, "-")
				fmt.Fprintf(&out, "%s %s %s %s %s %s\n", item.ID, item.CreatedAt, strings.Join(item.Repos, ","), item.State, at, item.Title)
			}
			for _, invalid := range result.Invalid {
				fmt.Fprintf(&out, "invalid: %s: %s\n", invalid.File, invalid.Error)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	})
	listCommand.Flags().StringVar(&view, "view", string(assignment.ViewAvailable), "available, taken, archive or all")
	return command
}

// resolveAssignmentProject maps Project lookup failures to fixed details
// without paths, as the backlog command does.
func resolveAssignmentProject(d workspace.Dependencies, value string) (workspace.ProjectRecord, error) {
	project, err := workspace.ResolveDefinedProject(d, workspace.ProjectID(value))
	if err == nil {
		return project, nil
	}
	subject := "project " + value
	if value == "" {
		subject = "the project of the current directory"
	}
	var projectError *workspace.ProjectError
	switch {
	case errors.As(err, &projectError) && projectError.Class == workspace.ErrorProjectNotFound:
		return workspace.ProjectRecord{}, &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: subject + " is not registered"}
	case errors.As(err, &projectError) && projectError.Class == workspace.ErrorProjectInvalidArguments && value == "":
		// The selection detail only counts Projects; it carries no path.
		return workspace.ProjectRecord{}, &workspace.ProjectError{Class: projectError.Class, Detail: projectError.Detail}
	case errors.As(err, &projectError):
		return workspace.ProjectRecord{}, &workspace.ProjectError{Class: projectError.Class, Detail: "could not resolve " + subject}
	}
	return workspace.ProjectRecord{}, &assignment.Error{Class: assignment.ClassIOError, Detail: "could not resolve " + subject}
}

func writeAssignmentText(out io.Writer, a assignment.Assignment, l assignment.Lifecycle) error {
	_, err := fmt.Fprintf(out, "%s %s\nrepos: %s\ntopics: %s\ncreated: %s\nstate: %s\n\n%s", a.ID, a.Title, strings.Join(a.Repos, ", "), strings.Join(a.Topics, ", "), a.CreatedAt, lifecycleText(l), a.Text)
	return err
}

func lifecycleText(l assignment.Lifecycle) string {
	switch {
	case l.Holder != nil:
		return fmt.Sprintf("taken by %s since %s (revision %d, claim %s)", l.Holder.Actor, l.Holder.TakenAt, l.Revision, l.Holder.Claim)
	case l.Done != nil:
		return fmt.Sprintf("done by %s at %s (revision %d)", l.Done.Actor, l.Done.At, l.Revision)
	}
	return fmt.Sprintf("%s (revision %d)", l.State, l.Revision)
}

func assignmentInvalidArguments(detail string) error {
	return &assignment.Error{Class: assignment.ClassInvalidArguments, Detail: detail}
}

func writeAssignmentJSON(out io.Writer, value any) error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return err
	}
	_, err := out.Write(buffer.Bytes())
	return err
}

func init() {
	RootCmd.AddCommand(newAssignmentCommand(workspace.SystemDependencies()))
}
