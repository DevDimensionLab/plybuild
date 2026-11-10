package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/backlog"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// backlogStdinLimit reads one byte past the store's 64 KiB text limit so the
// store, not the CLI, rejects oversized text.
const backlogStdinLimit = 64*1024 + 1

func newBacklogCommand(dependencies workspace.Dependencies) *cobra.Command {
	var root string
	command := newBacklogCommandWithDirectory(func(project string) (backlogPaths, error) {
		d := dependencies
		if root != "" {
			var err error
			if d, err = workspace.AtRoot(dependencies, root); err != nil {
				return backlogPaths{}, backlogInvalidArguments("workspace must be its physical root")
			}
		}
		return backlogDirectory(func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ShowProject(d, id)
		})(project)
	})
	command.PersistentFlags().StringVar(&root, "workspace", "", "physical workspace root (default: the workspace containing the current directory)")
	return command
}

// backlogPaths are the directories and repo IDs of one project's backlog.
type backlogPaths struct {
	directory, archive, wrapper string
	members                     []workspace.RepoID
}

// newBacklogCommandWithDirectory builds `ply backlog`; resolve maps a project ID
// to its backlog and archive directories, its wrapper and its repo IDs.
func newBacklogCommandWithDirectory(resolve func(string) (backlogPaths, error)) *cobra.Command {
	backlogCommand := &cobra.Command{
		Use:     "backlog",
		Short:   "Read, append to and set status in a project backlog",
		Long:    "Read, append to and set the status (open, ready or done) of items in the backlog of a project registered in the containing Ply workspace.\nItems live as one JSON file each in the project's data_paths.backlog in ply-workspace.json, relative to the project wrapper (default <project wrapper>/backlog/). Item and note text is read from stdin; note --author records who wrote it, and tags replaces the item's repo and topic lists. archive moves a done item into a named group in data_paths.backlog-archive (default <project wrapper>/backlog-archive/); groups and group read them. add --file instead stores a copy of that file and leaves the text empty. attach adds copies of one or more files to an existing item in argument order without changing it; it stops at the first failure, keeps the files attached before it and reports how many. attachments lists them and show includes them. archive refuses items with attachments.\n--workspace selects a workspace by its physical root instead of the current directory.",
		Example: "  ply backlog list --project ply\n  echo 'Text' | ply backlog add --project ply --title 'Title'\n  ply backlog --workspace ~/ws add --project ply --title 'Spec' --file 'my spec.md'\n  ply backlog attach --project ply 20261104-47530e79 trace.png 'my log.txt'\n  ply backlog attachments --project ply --format json 20261104-47530e79\n  echo 'Note' | ply backlog note --project ply 20261104-47530e79\n  ply backlog status --project ply --value ready 20261104-47530e79\n  ply backlog tags --project ply --repo ply --topic 'Web app' 20261104-47530e79\n  ply backlog list --project ply --repo ply --topic 'web app'\n  ply backlog archive --project ply --group backlog-kjerne --new-group 20261104-47530e79\n  ply backlog groups --project ply\n  ply backlog group --project ply backlog-kjerne",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return backlogInvalidArguments("unknown backlog command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	backlogCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return backlogInvalidArguments("invalid flag")
	})

	var project, format, title, value, filter, file, author, repoFilter, topicFilter string
	var repos, topics []string
	var clearTags bool
	var noteNumber int
	addFlags := func(command *cobra.Command) *cobra.Command {
		command.Flags().StringVar(&project, "project", "", "project ID")
		command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		return command
	}
	validate := func(positional int) cobra.PositionalArgs {
		return func(cmd *cobra.Command, args []string) error {
			if len(args) != positional {
				return backlogInvalidArguments(fmt.Sprintf("expected %d positional arguments, got %d", positional, len(args)))
			}
			if project == "" {
				return backlogInvalidArguments("--project is required")
			}
			if format != "text" && format != "json" {
				return backlogInvalidArguments("--format must be text or json")
			}
			return nil
		}
	}
	var members []workspace.RepoID
	var paths backlogPaths
	open := func() (backlog.Store, error) {
		var err error
		if paths, err = resolve(project); err != nil {
			return backlog.Store{}, err
		}
		members = paths.members
		return backlog.Open(paths.directory), nil
	}
	checkMembers := func() error {
		for _, repo := range repos {
			if !slices.Contains(members, workspace.RepoID(strings.TrimSpace(repo))) {
				return backlogInvalidArguments("--repo must be a repo ID of project " + project)
			}
		}
		return nil
	}
	readText := func(cmd *cobra.Command) (string, error) {
		text, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), backlogStdinLimit))
		if err != nil {
			return "", &backlog.Error{Class: backlog.ClassIOError, Detail: "read stdin"}
		}
		return string(text), nil
	}
	renderItem := func(cmd *cobra.Command, store backlog.Store, item backlog.Item) error {
		if format == "json" {
			return writeBacklogJSON(cmd.OutOrStdout(), item)
		}
		return writeBacklogItemText(cmd.OutOrStdout(), store, item, nil)
	}
	renderAttached := func(cmd *cobra.Command, store backlog.Store, item backlog.Item, attachments []backlog.Attachment) error {
		if format != "json" {
			return writeBacklogItemText(cmd.OutOrStdout(), store, item, attachments)
		}
		return writeBacklogJSON(cmd.OutOrStdout(), struct {
			backlog.Item
			Attachments []backlog.Attachment `json:"attachments,omitempty"`
		}{item, attachments})
	}

	listCommand := addFlags(&cobra.Command{
		Use:   "list",
		Short: "List backlog items, optionally only those with one status",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(0)(cmd, args); err != nil {
				return err
			}
			if cmd.Flags().Changed("status") && !validBacklogStatus(filter) {
				return backlogInvalidArguments("--status must be open, ready or done")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			result, err := store.List()
			if err != nil {
				return err
			}
			result.Items = filterBacklogItems(result.Items, backlog.Status(filter), repoFilter, topicFilter)
			if format == "json" {
				return writeBacklogJSON(cmd.OutOrStdout(), struct {
					Project string `json:"project"`
					backlog.ListResult
				}{project, result})
			}
			return writeBacklogListText(cmd.OutOrStdout(), project, filter, result)
		},
	})
	listCommand.Flags().StringVar(&filter, "status", "", "only items with this status (open, ready or done)")
	listCommand.Flags().StringVar(&repoFilter, "repo", "", "only items tagged with this repo ID")
	listCommand.Flags().StringVar(&topicFilter, "topic", "", "only items tagged with this topic, ignoring case")
	showCommand := addFlags(&cobra.Command{
		Use:   "show <item-id>",
		Short: "Show one backlog item",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			item, err := store.Read(args[0])
			if err != nil {
				return err
			}
			attachments, err := store.Attachments(item.ID)
			if err != nil {
				return err
			}
			return renderAttached(cmd, store, item, attachments)
		},
	})
	attachCommand := addFlags(&cobra.Command{
		Use:   "attach <item-id> <file>...",
		Short: "Add copies of files to an existing backlog item, in argument order",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return backlogInvalidArguments("expected an item ID and at least one file")
			}
			return validate(len(args))(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			var attachments []backlog.Attachment
			for i, file := range args[1:] {
				source, err := filepath.Abs(file)
				if err != nil {
					err = backlogInvalidArguments("could not resolve a file")
				} else {
					attachments, err = store.Attach(args[0], source)
				}
				var failure *backlog.Error
				if errors.As(err, &failure) {
					return &backlog.Error{Class: failure.Class, Detail: fmt.Sprintf("attached %d of %d files before: %s", i, len(args)-1, failure.Detail)}
				} else if err != nil {
					return err
				}
			}
			if format == "json" {
				return writeBacklogJSON(cmd.OutOrStdout(), backlogAttachments{args[0], attachments})
			}
			item, err := store.Read(args[0])
			if err != nil {
				return err
			}
			return renderAttached(cmd, store, item, attachments)
		},
	})
	attachmentsCommand := addFlags(&cobra.Command{
		Use:   "attachments <item-id>",
		Short: "List the attachments of a backlog item",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			attachments, err := store.Attachments(args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeBacklogJSON(cmd.OutOrStdout(), backlogAttachments{args[0], attachments})
			}
			return writeBacklogAttachmentsText(cmd.OutOrStdout(), store, args[0], attachments)
		},
	})
	addCommand := addFlags(&cobra.Command{
		Use:   "add",
		Short: "Add a backlog item with text from stdin or a copy of --file",
		Args:  validate(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("file") && file == "" {
				return backlogInvalidArguments("--file must not be empty")
			}
			store, err := open()
			if err != nil {
				return err
			}
			if file != "" {
				source, err := filepath.Abs(file)
				if err != nil {
					return backlogInvalidArguments("could not resolve --file")
				}
				item, err := store.AddFile(title, source)
				if err != nil {
					return err
				}
				return renderItem(cmd, store, item)
			}
			text, err := readText(cmd)
			if err != nil {
				return err
			}
			item, err := store.Add(title, text)
			if err != nil {
				return err
			}
			return renderItem(cmd, store, item)
		},
	})
	addCommand.Flags().StringVar(&title, "title", "", "one-line item title")
	addCommand.Flags().StringVar(&file, "file", "", "store a copy of this file instead of reading stdin")
	noteCommand := addFlags(&cobra.Command{
		Use:   "note <item-id>",
		Short: "Append a note from stdin to a backlog item",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			if err := checkMembers(); err != nil {
				return err
			}
			text, err := readText(cmd)
			if err != nil {
				return err
			}
			item, err := store.AddNoteTagged(args[0], text, author, repos, topics)
			if err != nil {
				return err
			}
			return renderItem(cmd, store, item)
		},
	})
	noteCommand.Flags().StringVar(&author, "author", "", "who wrote the note (shown as given, not verified)")
	noteCommand.Flags().StringArrayVar(&repos, "repo", nil, "repo ID of the project (repeatable)")
	noteCommand.Flags().StringArrayVar(&topics, "topic", nil, "topic (repeatable)")
	tagsCommand := addFlags(&cobra.Command{
		Use:   "tags <item-id>",
		Short: "Replace the repo IDs and topics of a backlog item or, with --note N, of its Nth note; --clear removes both",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			if clearTags == (len(repos)+len(topics) > 0) {
				return backlogInvalidArguments("give --repo or --topic, or --clear alone")
			}
			if cmd.Flags().Changed("note") && noteNumber < 1 {
				return backlogInvalidArguments("--note must be a note number from 1")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			if err := checkMembers(); err != nil {
				return err
			}
			var item backlog.Item
			if cmd.Flags().Changed("note") {
				item, err = store.SetNoteTags(args[0], noteNumber-1, repos, topics)
			} else {
				item, err = store.SetTags(args[0], repos, topics)
			}
			if err != nil {
				return err
			}
			return renderItem(cmd, store, item)
		},
	})
	tagsCommand.Flags().StringArrayVar(&repos, "repo", nil, "repo ID of the project (repeatable)")
	tagsCommand.Flags().StringArrayVar(&topics, "topic", nil, "topic (repeatable)")
	tagsCommand.Flags().BoolVar(&clearTags, "clear", false, "remove all repo IDs and topics")
	tagsCommand.Flags().IntVar(&noteNumber, "note", 0, "tag the note with this number as shown (1 is the first) instead of the item")
	statusCommand := addFlags(&cobra.Command{
		Use:     "status <item-id>",
		Short:   "Set the status of a backlog item to open, ready or done",
		Example: "  ply backlog status --project ply --value ready 20261104-47530e79",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			if !validBacklogStatus(value) {
				return backlogInvalidArguments("--value must be open, ready or done")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := open()
			if err != nil {
				return err
			}
			item, err := store.SetStatus(args[0], backlog.Status(value))
			if err != nil {
				return err
			}
			return renderItem(cmd, store, item)
		},
	})
	statusCommand.Flags().StringVar(&value, "value", "", "status (open, ready or done)")
	var group string
	var newGroup bool
	openArchive := func() (backlog.Store, backlog.Archive, error) {
		store, err := open()
		return store, backlog.OpenArchive(paths.archive, paths.wrapper), err
	}
	archiveCommand := addFlags(&cobra.Command{
		Use:   "archive <item-id>",
		Short: "Move a done backlog item into an archive group; --new-group creates the group",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := validate(1)(cmd, args); err != nil {
				return err
			}
			if group == "" {
				return backlogInvalidArguments("--group is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, archive, err := openArchive()
			if err != nil {
				return err
			}
			item, err := store.Archive(args[0], archive, group, newGroup)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeBacklogJSON(cmd.OutOrStdout(), struct {
					backlog.ArchivedItem
					Group string `json:"group"`
				}{item, group})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Archived %s to %s.\n", item.ID, group)
			return err
		},
	})
	archiveCommand.Flags().StringVar(&group, "group", "", "archive group (lowercase letters, digits and -)")
	archiveCommand.Flags().BoolVar(&newGroup, "new-group", false, "create the group; it must not exist")
	groupsCommand := addFlags(&cobra.Command{
		Use:   "groups",
		Short: "List archive groups",
		Args:  validate(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, archive, err := openArchive()
			if err != nil {
				return err
			}
			result, err := archive.ListGroups()
			if err != nil {
				return err
			}
			if format == "json" {
				return writeBacklogJSON(cmd.OutOrStdout(), struct {
					Project string `json:"project"`
					backlog.GroupListResult
				}{project, result})
			}
			return writeBacklogGroupsText(cmd.OutOrStdout(), project, result)
		},
	})
	groupCommand := addFlags(&cobra.Command{
		Use:   "group <group>",
		Short: "Show the archived items of one group",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, archive, err := openArchive()
			if err != nil {
				return err
			}
			read, stored, err := archive.ReadGroupFile(args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				_, err := cmd.OutOrStdout().Write(stored)
				return err
			}
			return writeBacklogGroupText(cmd.OutOrStdout(), archive, read)
		},
	})
	backlogCommand.AddCommand(listCommand, showCommand, addCommand, attachCommand, attachmentsCommand, noteCommand, statusCommand, tagsCommand, archiveCommand, groupsCommand, groupCommand)
	return backlogCommand
}

// backlogDirectory resolves a project ID to its backlog data path. Workspace
// error details carry absolute paths, so only their class is kept; data path
// errors carry none and are kept whole.
func backlogDirectory(show func(workspace.ProjectID) (workspace.ProjectResult, error)) func(string) (backlogPaths, error) {
	return func(value string) (backlogPaths, error) {
		id, err := workspace.ParseProjectID(value)
		if err != nil {
			return backlogPaths{}, backlogInvalidArguments("invalid project ID")
		}
		result, err := show(id)
		if err == nil {
			paths := backlogPaths{wrapper: result.Project.Wrapper, members: result.Project.RepoIDs}
			if paths.directory, err = workspace.ProjectDataPath(result.Project, workspace.DataPathBacklog); err != nil {
				return backlogPaths{}, err
			}
			paths.archive, err = workspace.ProjectDataPath(result.Project, workspace.DataPathBacklogArchive)
			return paths, err
		}
		var projectError *workspace.ProjectError
		switch {
		case errors.As(err, &projectError) && projectError.Class == workspace.ErrorProjectNotFound:
			return backlogPaths{}, &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: "project " + value + " is not registered"}
		case errors.As(err, &projectError):
			return backlogPaths{}, &workspace.ProjectError{Class: projectError.Class, Detail: "could not resolve project " + value}
		}
		return backlogPaths{}, &backlog.Error{Class: backlog.ClassIOError, Detail: "could not resolve project " + value}
	}
}

func validBacklogStatus(value string) bool {
	switch backlog.Status(value) {
	case backlog.StatusOpen, backlog.StatusReady, backlog.StatusDone:
		return true
	}
	return false
}

// filterBacklogItems keeps items matching every non-empty filter: the status,
// the exact repo ID and the topic ignoring case.
func filterBacklogItems(items []backlog.Summary, status backlog.Status, repo, topic string) []backlog.Summary {
	filtered := []backlog.Summary{}
	for _, item := range items {
		if (status == "" || item.Status == status) && (repo == "" || slices.Contains(item.Repos, repo)) &&
			(topic == "" || slices.ContainsFunc(item.Topics, func(t string) bool { return strings.EqualFold(t, topic) })) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func backlogInvalidArguments(detail string) error {
	return &backlog.Error{Class: backlog.ClassInvalidArguments, Detail: detail}
}

// writeBacklogJSON uses the same encoding as the item files.
func writeBacklogJSON(out io.Writer, value any) error {
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

func writeBacklogListText(out io.Writer, project, filter string, result backlog.ListResult) error {
	var b strings.Builder
	if len(result.Items) == 0 && len(result.Invalid) == 0 {
		if filter != "" {
			fmt.Fprintf(&b, "No %s backlog items in project %s.\n", filter, project)
		} else {
			fmt.Fprintf(&b, "No backlog items in project %s.\n", project)
		}
	}
	for _, item := range result.Items {
		fmt.Fprintf(&b, "%s  %s  %d notes  %s  %s%s\n", item.ID, item.UpdatedAt, item.NoteCount, item.Status, item.Title, backlogTagsText("  repos: ", "  topics: ", "", item.Repos, item.Topics))
	}
	for _, invalid := range result.Invalid {
		fmt.Fprintf(&b, "invalid: %s: %s\n", invalid.File, invalid.Error)
	}
	_, err := io.WriteString(out, b.String())
	return err
}

type backlogAttachments struct {
	ID          string               `json:"id"`
	Attachments []backlog.Attachment `json:"attachments"`
}

func writeBacklogAttachmentsText(out io.Writer, store backlog.Store, id string, attachments []backlog.Attachment) error {
	var b strings.Builder
	if len(attachments) == 0 {
		fmt.Fprintf(&b, "No attachments on item %s.\n", id)
	}
	for _, attachment := range attachments {
		path, state, err := store.AttachmentState(id, attachment)
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "%s  %d bytes  %s  %s  %s\n", attachment.Name, attachment.Size, attachment.AddedAt, state, path)
	}
	_, err := io.WriteString(out, b.String())
	return err
}

func writeBacklogItemText(out io.Writer, store backlog.Store, item backlog.Item, attachments []backlog.Attachment) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%s  %s\nCreated: %s\nStatus: %s\n%s", item.ID, item.Title, item.CreatedAt, item.Status, backlogTagsText("Repos: ", "Topics: ", "\n", item.Repos, item.Topics))
	if item.File != nil {
		path, state, err := store.FileState(item)
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "File: %s (%d bytes, %s)\nPath: %s\n", item.File.Name, item.File.Size, state, path)
	}
	for _, attachment := range attachments {
		path, state, err := store.AttachmentState(item.ID, attachment)
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "Attachment: %s (%d bytes, %s)\nPath: %s\n", attachment.Name, attachment.Size, state, path)
	}
	fmt.Fprintf(&b, "\n%s", withNewline(item.Text))
	for i, note := range item.Notes {
		by := ""
		if note.Author != "" {
			by = " by " + note.Author
		}
		fmt.Fprintf(&b, "\nNote %d %s%s:\n%s%s", i+1, note.At, by, backlogTagsText("Repos: ", "Topics: ", "\n", note.Repos, note.Topics), withNewline(note.Text))
	}
	_, err := io.WriteString(out, b.String())
	return err
}

func writeBacklogGroupsText(out io.Writer, project string, result backlog.GroupListResult) error {
	var b strings.Builder
	if len(result.Groups) == 0 && len(result.Invalid) == 0 {
		fmt.Fprintf(&b, "No archive groups in project %s.\n", project)
	}
	for _, group := range result.Groups {
		fmt.Fprintf(&b, "%s  %d items  %s\n", group.Group, group.ItemCount, group.LastArchivedAt)
	}
	for _, invalid := range result.Invalid {
		fmt.Fprintf(&b, "invalid: %s: %s\n", invalid.File, invalid.Error)
	}
	_, err := io.WriteString(out, b.String())
	return err
}

// writeBacklogGroupText shows each archived item as show does, with Archived in place of Status.
func writeBacklogGroupText(out io.Writer, archive backlog.Archive, group backlog.ArchiveGroup) error {
	var b strings.Builder
	for i, item := range group.Items {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "%s  %s\nCreated: %s\nArchived: %s\n%s", item.ID, item.Title, item.CreatedAt, item.ArchivedAt, backlogTagsText("Repos: ", "Topics: ", "\n", item.Repos, item.Topics))
		if item.File != nil {
			path, state, err := archive.FileState(item)
			if err != nil {
				return err
			}
			fmt.Fprintf(&b, "File: %s (%d bytes, %s)\nPath: %s\n", item.File.Name, item.File.Size, state, path)
		}
		fmt.Fprintf(&b, "\n%s", withNewline(item.Text))
		for j, note := range item.Notes {
			by := ""
			if note.Author != "" {
				by = " by " + note.Author
			}
			fmt.Fprintf(&b, "\nNote %d %s%s:\n%s%s", j+1, note.At, by, backlogTagsText("Repos: ", "Topics: ", "\n", note.Repos, note.Topics), withNewline(note.Text))
		}
	}
	_, err := io.WriteString(out, b.String())
	return err
}

// backlogTagsText writes each non-empty list after its label and before end.
func backlogTagsText(reposLabel, topicsLabel, end string, repos, topics []string) string {
	var b strings.Builder
	if len(repos) > 0 {
		b.WriteString(reposLabel + strings.Join(repos, ", ") + end)
	}
	if len(topics) > 0 {
		b.WriteString(topicsLabel + strings.Join(topics, ", ") + end)
	}
	return b.String()
}

func withNewline(text string) string {
	if strings.HasSuffix(text, "\n") {
		return text
	}
	return text + "\n"
}

func init() {
	RootCmd.AddCommand(newBacklogCommand(workspace.SystemDependencies()))
}
