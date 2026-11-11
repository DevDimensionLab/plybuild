package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/inbox"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// inboxRequestLimit reads one byte past the 512 KiB request limit so the
// decoder, not the CLI, rejects oversized requests.
const inboxRequestLimit = 512*1024 + 1

func newInboxCommand(dependencies workspace.Dependencies) *cobra.Command {
	var root string
	command := newInboxCommandWithDirectory(func(project string) (string, error) {
		d := dependencies
		if root != "" {
			var err error
			if d, err = workspace.AtRoot(dependencies, root); err != nil {
				return "", inboxInvalidArguments("workspace must be its physical root")
			}
		}
		return inboxDirectory(func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ShowProject(d, id)
		})(project)
	})
	command.PersistentFlags().StringVar(&root, "workspace", "", "physical workspace root (default: the workspace containing the current directory)")
	return command
}

// newInboxCommandWithDirectory builds `ply inbox`; resolve maps a project ID to its Inbox data path.
func newInboxCommandWithDirectory(resolve func(string) (string, error)) *cobra.Command {
	var project, owner, format, file, conversation, message, value, cursor string
	var limit int
	// fail prints InboxError@1 on stdout for --format json and returns err for stderr and exit 1.
	fail := func(cmd *cobra.Command, err error) error {
		if err == nil || format != "json" {
			return err
		}
		if result, ok := inboxErrorResult(err); ok {
			if writeErr := inbox.Encode(cmd.OutOrStdout(), result); writeErr != nil {
				return writeErr
			}
		}
		return err
	}
	inboxCommand := &cobra.Command{
		Use:   "inbox",
		Short: "Store and read canonical agent replies and their read marks",
		Long: "Store and read the canonical agent replies of one owner in a project registered in the containing Ply workspace, with shared per-item read marks.\n" +
			"Replies live as one write-once JSON file each in the project's data_paths.inbox (default <project wrapper>/inbox/), under --owner. " +
			"put stores one InboxPutRequest@1 from --file (- for stdin): a replay of the same payload is unchanged, a different payload under the same conversation and message is a conflict. " +
			"Both created and unchanged confirm that the payload is stored. show and list read; list pages newest first with an opaque --cursor. mark sets read or unread and never changes a message or the order.\n" +
			"Nothing here starts agents, sends notifications or uses the network. The contract bundle is schemas/ply-inbox-v1.",
		Example: "  ply inbox put --project ply --owner owner-1 --file reply.json --format json\n" +
			"  ply inbox list --project ply --owner owner-1 --limit 20\n" +
			"  ply inbox show --project ply --owner owner-1 --conversation conv_… --message msg_…\n" +
			"  ply inbox mark --project ply --owner owner-1 --conversation conv_… --message msg_… --value read",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return inboxInvalidArguments("unknown inbox command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() },
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	inboxCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return fail(cmd, inboxInvalidArguments("invalid flag"))
	})
	addFlags := func(command *cobra.Command) *cobra.Command {
		command.Flags().StringVar(&project, "project", "", "project ID")
		command.Flags().StringVar(&owner, "owner", "", "owner ID ([A-Za-z0-9_.-]{1,64})")
		command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		return command
	}
	addIdentity := func(command *cobra.Command) *cobra.Command {
		command.Flags().StringVar(&conversation, "conversation", "", "conversation ID (conv_ and 32 hex)")
		command.Flags().StringVar(&message, "message", "", "message ID (msg_ and 32 hex)")
		return command
	}
	// verb validates the common flags, opens the store and runs; every failure goes through fail.
	verb := func(run func(*cobra.Command, *inbox.Store) error) func(*cobra.Command, []string) error {
		return func(cmd *cobra.Command, args []string) error {
			switch {
			case format != "text" && format != "json":
				return inboxInvalidArguments("--format must be text or json")
			case len(args) != 0:
				return fail(cmd, inboxInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args))))
			case project == "":
				return fail(cmd, inboxInvalidArguments("--project is required"))
			case owner == "":
				return fail(cmd, inboxInvalidArguments("--owner is required"))
			}
			if err := inbox.CheckOwner(owner); err != nil {
				return fail(cmd, err)
			}
			directory, err := resolve(project)
			if err != nil {
				return fail(cmd, err)
			}
			store, err := inbox.Open(directory, project, owner)
			if err != nil {
				return fail(cmd, err)
			}
			return fail(cmd, run(cmd, store))
		}
	}
	render := func(cmd *cobra.Command, value any, text string) error {
		if format == "json" {
			return inbox.Encode(cmd.OutOrStdout(), value)
		}
		_, err := io.WriteString(cmd.OutOrStdout(), text)
		return err
	}

	putCommand := addFlags(&cobra.Command{
		Use:   "put",
		Short: "Store one agent reply from an InboxPutRequest@1",
		RunE: verb(func(cmd *cobra.Command, store *inbox.Store) error {
			if file == "" {
				return inboxInvalidArguments("--file is required")
			}
			raw, err := readInboxRequest(cmd, file)
			if err != nil {
				return err
			}
			req, err := inbox.DecodePutRequest(raw)
			if err != nil {
				return err
			}
			result, err := store.Put(req)
			if err != nil {
				return err
			}
			return render(cmd, result, fmt.Sprintf("%s %s/%s %s\n", result.Effect, result.Item.ConversationID, result.Item.MessageID, result.Item.PayloadSHA256))
		}),
	})
	putCommand.Flags().StringVar(&file, "file", "", "InboxPutRequest@1 JSON file, or - for stdin")
	showCommand := addIdentity(addFlags(&cobra.Command{
		Use:   "show",
		Short: "Show one reply with its read state",
		RunE: verb(func(cmd *cobra.Command, store *inbox.Store) error {
			item, err := store.Show(conversation, message)
			if err != nil {
				return err
			}
			return render(cmd, item, inboxItemText(item))
		}),
	}))
	listCommand := addFlags(&cobra.Command{
		Use:   "list",
		Short: "List replies newest first, one page at a time",
		RunE: verb(func(cmd *cobra.Command, store *inbox.Store) error {
			list, err := store.List(limit, cursor)
			if err != nil {
				return err
			}
			return render(cmd, list, inboxListText(list))
		}),
	})
	listCommand.Flags().IntVar(&limit, "limit", 20, "page size (1-50)")
	listCommand.Flags().StringVar(&cursor, "cursor", "", "next_cursor of an earlier list of the same project and owner")
	markCommand := addIdentity(addFlags(&cobra.Command{
		Use:   "mark",
		Short: "Mark one reply read or unread",
		RunE: verb(func(cmd *cobra.Command, store *inbox.Store) error {
			result, err := store.Mark(conversation, message, inbox.MarkValue(value))
			if err != nil {
				return err
			}
			return render(cmd, result, fmt.Sprintf("%s %s/%s %s\n", result.Effect, result.ConversationID, result.MessageID, readWord(&result.Read)))
		}),
	}))
	markCommand.Flags().StringVar(&value, "value", "", "read or unread")
	inboxCommand.AddCommand(putCommand, showCommand, listCommand, markCommand)
	return inboxCommand
}

// inboxDirectory resolves a project ID to its Inbox data path. Workspace
// error details carry absolute paths, so only their class is kept; data path
// errors carry none and are kept whole.
func inboxDirectory(show func(workspace.ProjectID) (workspace.ProjectResult, error)) func(string) (string, error) {
	return func(value string) (string, error) {
		id, err := workspace.ParseProjectID(value)
		if err != nil {
			return "", inboxInvalidArguments("invalid project ID")
		}
		result, err := show(id)
		if err == nil {
			return workspace.ProjectDataPath(result.Project, workspace.DataPathInbox)
		}
		var projectError *workspace.ProjectError
		switch {
		case errors.As(err, &projectError) && projectError.Class == workspace.ErrorProjectNotFound:
			return "", &workspace.ProjectError{Class: workspace.ErrorProjectNotFound, Detail: "project " + value + " is not registered"}
		case errors.As(err, &projectError):
			return "", &workspace.ProjectError{Class: projectError.Class, Detail: "could not resolve project " + value}
		}
		return "", &inbox.Error{Class: inbox.ClassIOError, Detail: "could not resolve project " + value}
	}
}

func readInboxRequest(cmd *cobra.Command, file string) ([]byte, error) {
	if file == "-" {
		raw, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), inboxRequestLimit))
		if err != nil {
			return nil, &inbox.Error{Class: inbox.ClassIOError, Detail: "read stdin"}
		}
		return raw, nil
	}
	opened, err := os.Open(file)
	if err != nil {
		return nil, inboxInvalidArguments("cannot open --file")
	}
	defer opened.Close()
	raw, err := io.ReadAll(io.LimitReader(opened, inboxRequestLimit))
	if err != nil {
		return nil, inboxInvalidArguments("cannot read --file")
	}
	return raw, nil
}

// inboxErrorResult is the InboxError@1 of an inbox or workspace project error.
func inboxErrorResult(err error) (inbox.ErrorResult, bool) {
	var inboxError *inbox.Error
	if errors.As(err, &inboxError) {
		return inbox.NewErrorResult(inboxError), true
	}
	var projectError *workspace.ProjectError
	if errors.As(err, &projectError) {
		return inbox.ErrorResult{Kind: inbox.ErrorKind, SchemaVersion: inbox.SchemaVersion, Class: string(projectError.Class), Detail: projectError.Detail}, true
	}
	return inbox.ErrorResult{}, false
}

func inboxInvalidArguments(detail string) error {
	return &inbox.Error{Class: inbox.ClassInvalidArguments, Detail: detail}
}

func readWord(read *bool) string {
	switch {
	case read == nil:
		return "unknown"
	case *read:
		return "read"
	}
	return "unread"
}

func inboxItemText(item inbox.Item) string {
	var b strings.Builder
	answer := map[string]string{"read": "yes", "unread": "no", "unknown": "unknown"}[readWord(item.Read)]
	fmt.Fprintf(&b, "%s\nItem: %s/%s\nReplied: %s\nTitle source: %s\nRead: %s\nAgent: %s (%s)\n", item.DisplayTitle, item.ConversationID, item.MessageID, item.RepliedAtUTC, item.TitleSource, answer, item.AgentLabel, item.Provider)
	if item.Truncated && item.FullSHA256 != nil {
		fmt.Fprintf(&b, "Truncated: full text sha256 %s\n", *item.FullSHA256)
	}
	fmt.Fprintf(&b, "\n%s", withNewline(item.Body))
	return b.String()
}

func inboxListText(list inbox.List) string {
	var b strings.Builder
	unread := "unknown"
	if list.UnreadCount != nil {
		unread = fmt.Sprint(*list.UnreadCount)
	}
	fmt.Fprintf(&b, "Inbox %s/%s: %d items, %s unread (%s)\n", list.ProjectID, list.Owner, list.TotalCount, unread, list.UnreadCountBasis)
	for _, item := range list.Items {
		fmt.Fprintf(&b, "%s  %s  %s/%s  %s\n", item.RepliedAtUTC, readWord(item.Read), item.ConversationID, item.MessageID, item.DisplayTitle)
	}
	for _, invalid := range list.Invalid {
		fmt.Fprintf(&b, "invalid: %s: %s\n", invalid.Entry, invalid.Error)
	}
	if list.NextCursor != nil {
		fmt.Fprintf(&b, "Next cursor: %s\n", *list.NextCursor)
	}
	return b.String()
}

func init() {
	RootCmd.AddCommand(newInboxCommand(workspace.SystemDependencies()))
}
