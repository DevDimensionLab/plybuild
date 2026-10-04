package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	notification "github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/spf13/cobra"
)

// NewWorkflowNotificationCommand is the production command graph with an internal
// effect boundary for native tests and the separate, network-free fixture driver.
func NewWorkflowNotificationCommand(d notification.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "notification", Short: "Preview and deliver human gates or agent events", Long: "Notify one fixed Slack channel about a v1 report-ready human gate or a v2 agent_finished, agent_stopped, or feedback_required event. V2 external sources bind a preserved event; agent_finished requires a report, and before_start is only valid for agent_stopped without a start receipt. Source, sender, gate, event and channel are claims; transport acknowledgement does not attest product correctness or human QA. No automatic detection or retry occurs.", Example: "  ply workflow notification send --file ./report-ready.json --route ./slack-route.json --check\n  ply workflow notification show NOTIFICATION_ID --route ./slack-route.json", PersistentPreRunE: func(*cobra.Command, []string) error { return nil }}
	parent.Args = func(c *cobra.Command, a []string) error {
		if len(a) != 0 {
			return notificationUsage(c, "Expected a notification subcommand.")
		}
		return nil
	}
	parent.RunE = func(c *cobra.Command, a []string) error {
		if len(a) != 0 {
			return notificationUsage(c, "Expected a notification subcommand.")
		}
		return c.Help()
	}
	parent.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return notificationUsage(c, "Invalid notification flags.") })
	for _, operation := range []string{"send", "show", "retry"} {
		operation := operation
		var file, route, format, confirm string
		var check, apply bool
		use := operation
		if operation != "send" {
			use += " ID"
		}
		short := map[string]string{"send": "Preview or send one human-gate or agent-event notification", "show": "Read a preserved notification without credentials or writes", "retry": "Preview or explicitly retry proven non-delivery"}[operation]
		long := short + ". "
		switch operation {
		case "send":
			long += "Read-only preview is the default. --apply requires the exact --confirm digest from preview and performs at most one Slack POST after durable reservation. Exact duplicate sends read the preserved result without another POST. Recheck that the gate or event is current and the route is authorized before apply. Feedback requires a necessary answer; optional questions and internal corrections do not trigger notifications."
		case "show":
			long += "Show separates transport state from source freshness. Unknown is a successful readback, never permission to resend. No locks, credentials, recovery writes or Slack calls are used."
		case "retry":
			long += "A fresh preview and its new --confirm digest are required for a new attempt. Only rejected, not_sent, or rate_limited after its deadline can retry. Acknowledged or unknown delivery cannot retry. No automatic retries, sleeping, force or gate-revision recovery exist. Recheck the gate or event before apply; never change event identity or route to bypass unknown."
		}
		example := "  ply workflow notification " + operation + " ID --route ./slack-route.json --format json"
		if operation == "send" {
			example = "  ply workflow notification send --file ./report-ready.json --route ./slack-route.json --check --format json\n  ply workflow notification send --file ./report-ready.json --route ./slack-route.json --apply --confirm 'sha256:<confirmation-from-preview>' --format json"
		}
		if operation == "retry" {
			example = "  ply workflow notification retry ID --route ./slack-route.json --check --format json\n  ply workflow notification retry ID --route ./slack-route.json --apply --confirm 'sha256:<new-retry-confirmation>' --format json"
		}
		c := &cobra.Command{Use: use, Short: short, Long: long, Example: example}
		c.Args = func(c *cobra.Command, args []string) error {
			count := 1
			if operation == "send" {
				count = 0
			}
			if len(args) != count || route == "" || operation == "send" && file == "" {
				return notificationUsage(c, "Provide the required route, request file or notification ID.")
			}
			if format != "text" && format != "json" {
				return notificationUsage(c, "--format must be text or json.")
			}
			if operation != "show" && (c.Flags().Changed("check") && c.Flags().Changed("apply") || apply && confirm == "" || !apply && c.Flags().Changed("confirm")) {
				return notificationUsage(c, "--check and --apply are exclusive; --confirm is required only with --apply.")
			}
			if c.Flags().Changed("force") {
				return notificationUsage(c, "Force is not supported for notifications.")
			}
			return nil
		}
		c.RunE = func(c *cobra.Command, args []string) error {
			in := notification.Input{Operation: operation, File: file, Route: route, Confirm: confirm, Apply: apply}
			if len(args) > 0 {
				in.ID = args[0]
			}
			out, err := notification.Execute(d, in)
			if renderErr := renderNotification(c, format, out); renderErr != nil {
				return renderErr
			}
			return err
		}
		c.Flags().StringVar(&route, "route", "", "Explicit non-secret route JSON file")
		c.Flags().StringVar(&format, "format", "text", "Output format: text or json")
		if operation == "send" {
			c.Flags().StringVar(&file, "file", "", "Strict v1 human-gate or v2 agent-event request JSON file")
		}
		if operation != "show" {
			c.Flags().BoolVar(&check, "check", false, "Read-only preview (the default)")
			c.Flags().BoolVar(&apply, "apply", false, "Apply one confirmed attempt or read an exact duplicate")
			c.Flags().StringVar(&confirm, "confirm", "", "Exact confirmation digest from the corresponding preview")
		}
		c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return notificationUsage(c, "Invalid notification flags.") })
		parent.AddCommand(c)
	}
	return parent
}
func notificationUsage(c *cobra.Command, message string) error {
	e := &notification.Error{Exit: 2, Code: "invalid_arguments", Message: message}
	format, _ := c.Flags().GetString("format")
	if format == "json" || notificationJSONRequested(os.Args[1:]) {
		_ = renderNotification(c, "json", notification.ErrorOutput(e))
	}
	return e
}
func renderNotification(c *cobra.Command, format string, value any) error {
	if format == "json" {
		enc := json.NewEncoder(c.OutOrStdout())
		enc.SetEscapeHTML(false)
		return enc.Encode(value)
	}
	switch v := value.(type) {
	case notification.Preview:
		fmt.Fprintf(c.OutOrStdout(), "Route: %s (%s; channel identity claimed)\nSource: %s\nAllowed: %t\nState root: %s\nPayload: mrkdwn=false parse=none unfurl_links=false unfurl_media=false\n%s\n", v.Route.Name, v.Route.ChannelLabel, v.SourceKind, v.Allowed, v.StateRoot, v.Payload.Text)
		for _, effect := range v.Effects {
			fmt.Fprintln(c.OutOrStdout(), "Effect:", effect)
		}
		if len(v.Effects) == 0 {
			fmt.Fprintln(c.OutOrStdout(), "Effects: none (read-only)")
		}
		if v.ExistingState != nil {
			fmt.Fprintln(c.OutOrStdout(), "Existing transport state:", *v.ExistingState)
		}
		if v.Confirmation != nil {
			fmt.Fprintln(c.OutOrStdout(), "Confirmation:", *v.Confirmation)
		}
		for _, reason := range v.Reasons {
			fmt.Fprintln(c.OutOrStdout(), reason.Message)
		}
	case notification.Result:
		fmt.Fprintf(c.OutOrStdout(), "Notification: %s\nTransport: %s\nSource freshness: %s\nPersistence: %s\nKnowledge: reported; not controlled; human QA is not attested.\nNext: %s\n", v.ID, v.State, v.Freshness, v.Persistence, v.NextAction)
		for _, reason := range v.Reasons {
			fmt.Fprintln(c.OutOrStdout(), reason.Message)
		}
	case notification.ErrorResult:
		for _, reason := range v.Reasons {
			fmt.Fprintln(c.OutOrStdout(), reason.Message)
		}
	}
	return nil
}

// Cobra stops flag parsing at the first invalid flag. At the executable boundary
// retain the caller's JSON error choice even when --format follows that flag.
func notificationJSONRequested(args []string) bool {
	format := ""
	for i, a := range args {
		if a == "--" {
			break
		}
		if a == "--format" && i+1 < len(args) {
			format = args[i+1]
		}
		if strings.HasPrefix(a, "--format=") {
			format = strings.TrimPrefix(a, "--format=")
		}
	}
	return format == "json"
}
