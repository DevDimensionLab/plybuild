package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/spf13/cobra"
)

type workflowHandoffServices struct {
	create       func(workflowhandoff.CreateInput) (workflowhandoff.CreateResult, error)
	show         func(workflowhandoff.HandoffID) (workflowhandoff.ShowResult, error)
	inspect      func(workflowhandoff.InspectInput) (workflowhandoff.InspectResult, error)
	submitStart  func(workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error)
	submitResult func(workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error)
	cancel       func(workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error)
	supersede    func(workflowhandoff.SupersedeInput) (workflowhandoff.CreateResult, error)
	abandon      func(workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error)
}

func newWorkflowHandoffCommand(dependencies workflowhandoff.Dependencies) *cobra.Command {
	return newWorkflowHandoffCommandWithServices(workflowHandoffServices{
		create: func(input workflowhandoff.CreateInput) (workflowhandoff.CreateResult, error) {
			return workflowhandoff.Create(dependencies, input)
		},
		show: func(id workflowhandoff.HandoffID) (workflowhandoff.ShowResult, error) {
			return workflowhandoff.Show(dependencies, id)
		},
		inspect: func(input workflowhandoff.InspectInput) (workflowhandoff.InspectResult, error) {
			return workflowhandoff.Inspect(dependencies, input)
		},
		submitStart: func(input workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error) {
			return workflowhandoff.SubmitStart(dependencies, input)
		},
		submitResult: func(input workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error) {
			return workflowhandoff.SubmitResultDocument(dependencies, input)
		},
		cancel: func(input workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
			return workflowhandoff.Cancel(dependencies, input)
		},
		supersede: func(input workflowhandoff.SupersedeInput) (workflowhandoff.CreateResult, error) {
			return workflowhandoff.Supersede(dependencies, input)
		},
		abandon: func(input workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
			return workflowhandoff.Abandon(dependencies, input)
		},
	})
}

func newWorkflowHandoffCommandWithServices(services workflowHandoffServices) *cobra.Command {
	parent := &cobra.Command{Use: "handoff", Short: "Manage file-based agent handoffs", Long: "Create, control, receive, and inspect immutable local agent handoffs.", Example: "  ply workflow handoff create --file /absolute/handoff-draft.json\n  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef"}
	setHandoffFlagErrors(parent)

	var createFile string
	create := &cobra.Command{Use: "create", Short: "Create an immutable agent handoff", Long: "Validate a handoff draft and publish one immutable handoff in the containing Ply workspace.", Example: "  ply workflow handoff create --file /absolute/handoff-draft.json", Args: noArguments(func() error {
		if createFile == "" {
			return fmt.Errorf("--file is required")
		}
		return nil
	}), RunE: func(cmd *cobra.Command, args []string) error {
		result, err := services.create(workflowhandoff.CreateInput{DraftPath: createFile})
		if err != nil {
			return err
		}
		first := "Created agent handoff %s.\n"
		if !result.Created {
			first = "Agent handoff %s already exists with identical content.\n"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), first+"Purpose: %s\nWorking directory: %s\nHandoff: %s\nNext action: Open a fresh recipient agent in the working directory and tell it: \"Read and execute the handoff at %s.\"\n", result.HandoffID, result.Purpose, result.Worktree, result.Locator, result.Locator)
		return err
	}}
	create.Flags().StringVar(&createFile, "file", "", "handoff draft JSON file")
	_ = create.MarkFlagRequired("file")
	setHandoffFlagErrors(create)

	show := &cobra.Command{Use: "show <handoff-id>", Short: "Show an agent handoff for a human", Long: "Show the lifecycle result, practical meaning, and one next action for an agent handoff.", Example: "  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef", Args: oneHandoffID, RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := workflowhandoff.ParseHandoffID(args[0])
		result, err := services.show(id)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Status: %s\nResult: %s\nMeaning: %s\nNext action: %s\n", result.Status, result.Result, result.Meaning, result.NextAction)
		return err
	}}
	setHandoffFlagErrors(show)

	var inspectHandoff, inspectFormat, inspectRaw string
	var acknowledgeSecret bool
	inspect := &cobra.Command{Use: "inspect", Short: "Inspect bound agent handoff evidence", Long: "Read the canonical machine projection or exact stored bytes for an agent handoff.", Example: "  ply workflow handoff inspect --handoff /absolute/handoff.json --format json\n  ply workflow handoff inspect --handoff /absolute/handoff.json --raw result", Args: noArguments(func() error {
		if err := validateHandoffLocator(inspectHandoff); err != nil {
			return err
		}
		if (inspectFormat == "") == (inspectRaw == "") {
			return fmt.Errorf("exactly one of --format or --raw is required")
		}
		if inspectFormat != "" && inspectFormat != "json" {
			return fmt.Errorf("--format must be json")
		}
		if inspectRaw != "" && inspectRaw != "handoff" && inspectRaw != "start" && inspectRaw != "result" {
			return fmt.Errorf("--raw must be handoff, start, or result")
		}
		if acknowledgeSecret && inspectRaw != "handoff" {
			return fmt.Errorf("--acknowledge-secret-exposure is only valid with --raw handoff")
		}
		return nil
	}), RunE: func(cmd *cobra.Command, args []string) error {
		result, err := services.inspect(workflowhandoff.InspectInput{HandoffLocator: inspectHandoff, Format: inspectFormat, Raw: inspectRaw, AcknowledgeSecret: acknowledgeSecret})
		if err != nil {
			return err
		}
		if _, err = cmd.OutOrStdout().Write(result.Bytes); err != nil {
			return err
		}
		if result.AppendLF {
			_, err = cmd.OutOrStdout().Write([]byte{'\n'})
		}
		return err
	}}
	inspect.Flags().StringVar(&inspectHandoff, "handoff", "", "absolute immutable handoff JSON path")
	inspect.Flags().StringVar(&inspectFormat, "format", "", "inspection output format (json)")
	inspect.Flags().StringVar(&inspectRaw, "raw", "", "exact stored document (handoff, start, or result)")
	inspect.Flags().BoolVar(&acknowledgeSecret, "acknowledge-secret-exposure", false, "acknowledge that raw handoff output exposes the reply secret")
	_ = inspect.MarkFlagRequired("handoff")
	setHandoffFlagErrors(inspect)

	submitStart := newSubmitCommand("submit-start", "Submit a bound agent start receipt", "Validate and atomically publish the first start receipt for an exact agent handoff.", "  ply workflow handoff submit-start --handoff /absolute/handoff.json --file /absolute/start.json", "start receipt draft JSON file", services.submitStart)
	submitResult := newSubmitCommand("submit-result", "Submit a bound agent terminal result", "Validate and atomically publish the first terminal result for an exact started handoff.", "  ply workflow handoff submit-result --handoff /absolute/handoff.json --file /absolute/result.json", "terminal result draft JSON file", services.submitResult)

	cancel := newControlCommand("cancel <handoff-id>", "Cancel an unstarted agent handoff", "Close an unstarted handoff without deleting its immutable contract.", "  ply workflow handoff cancel hnd_0123456789abcdef0123456789abcdef --reason \"No longer needed\"", "human-readable cancellation reason", false, func(input workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
		return services.cancel(input)
	}, func(cmd *cobra.Command, result workflowhandoff.ControlResult) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "Cancelled agent handoff %s.\nReason: %s\nNext action: No action now; the handoff is closed.\n", result.HandoffID, result.Reason)
		return err
	})

	var supersedeFile, supersedeReason string
	supersede := &cobra.Command{Use: "supersede <handoff-id>", Short: "Supersede an unstarted agent handoff", Long: "Close an unstarted handoff and publish a replacement run under the same activity.", Example: "  ply workflow handoff supersede hnd_0123456789abcdef0123456789abcdef --file /absolute/replacement.json --reason \"Corrected input binding\"", Args: func(cmd *cobra.Command, args []string) error {
		if err := oneHandoffID(cmd, args); err != nil {
			return err
		}
		if supersedeFile == "" {
			return workflowhandoff.InvalidArguments("--file is required")
		}
		if supersedeReason == "" {
			return workflowhandoff.InvalidArguments("--reason is required")
		}
		return nil
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := workflowhandoff.ParseHandoffID(args[0])
		result, err := services.supersede(workflowhandoff.SupersedeInput{HandoffID: id, DraftPath: supersedeFile, Reason: supersedeReason})
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Superseded agent handoff %s with %s.\nReason: %s\nPurpose: %s\nWorking directory: %s\nHandoff: %s\nNext action: Open a fresh recipient agent in the working directory and tell it: \"Read and execute the handoff at %s.\"\n", id, result.HandoffID, supersedeReason, result.Purpose, result.Worktree, result.Locator, result.Locator)
		return err
	}}
	supersede.Flags().StringVar(&supersedeFile, "file", "", "replacement handoff draft JSON file")
	supersede.Flags().StringVar(&supersedeReason, "reason", "", "human-readable superseding reason")
	_ = supersede.MarkFlagRequired("file")
	_ = supersede.MarkFlagRequired("reason")
	setHandoffFlagErrors(supersede)

	abandon := newControlCommand("abandon <handoff-id>", "Abandon a started handoff with unknown effects", "Close a started handoff after explicitly acknowledging that target effects are unknown.", "  ply workflow handoff abandon hnd_0123456789abcdef0123456789abcdef --reason \"Recipient disappeared after start\" --acknowledge-effects-unknown", "human-readable abandonment reason", true, func(input workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error) {
		return services.abandon(input)
	}, func(cmd *cobra.Command, result workflowhandoff.ControlResult) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "Abandoned agent handoff %s with target effects marked unknown.\nReason: %s\nNext action: Start the separate result control with `ply workflow handoff inspect --handoff %s --format json`.\n", result.HandoffID, result.Reason, result.Locator)
		return err
	})

	parent.AddCommand(create, show, inspect, submitStart, submitResult, cancel, supersede, abandon)
	return parent
}

func newSubmitCommand(use, short, long, example, fileHelp string, service func(workflowhandoff.SubmitInput) (workflowhandoff.SubmitResult, error)) *cobra.Command {
	var handoff, file string
	command := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Args: noArguments(func() error {
		if err := validateHandoffLocator(handoff); err != nil {
			return err
		}
		if file == "" {
			return fmt.Errorf("--file is required")
		}
		return nil
	}), RunE: func(cmd *cobra.Command, args []string) error {
		result, err := service(workflowhandoff.SubmitInput{HandoffLocator: handoff, DraftPath: file})
		if err != nil {
			return err
		}
		noun := "start receipt"
		label := "Start receipt"
		if result.Phase == "terminal" {
			noun = "terminal result"
			label = "Terminal result"
		}
		first := "Accepted " + noun + " %s.\n"
		if !result.Created {
			first = label + " %s already exists with identical content.\n"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), first+label+": %s\nSHA-256: %s\n", result.DocumentID, result.Locator, result.SHA256)
		return err
	}}
	command.Flags().StringVar(&handoff, "handoff", "", "absolute immutable handoff JSON path")
	command.Flags().StringVar(&file, "file", "", fileHelp)
	_ = command.MarkFlagRequired("handoff")
	_ = command.MarkFlagRequired("file")
	setHandoffFlagErrors(command)
	return command
}

func newControlCommand(use, short, long, example, reasonHelp string, ackRequired bool, service func(workflowhandoff.ControlInput) (workflowhandoff.ControlResult, error), render func(*cobra.Command, workflowhandoff.ControlResult) error) *cobra.Command {
	var reason string
	var acknowledge bool
	command := &cobra.Command{Use: use, Short: short, Long: long, Example: example, Args: func(cmd *cobra.Command, args []string) error {
		if err := oneHandoffID(cmd, args); err != nil {
			return err
		}
		if reason == "" {
			return workflowhandoff.InvalidArguments("--reason is required")
		}
		if ackRequired && !acknowledge {
			return workflowhandoff.InvalidArguments("--acknowledge-effects-unknown must be true")
		}
		return nil
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := workflowhandoff.ParseHandoffID(args[0])
		result, err := service(workflowhandoff.ControlInput{HandoffID: id, Reason: reason, AcknowledgeEffectsUnknown: acknowledge})
		if err != nil {
			return err
		}
		return render(cmd, result)
	}}
	command.Flags().StringVar(&reason, "reason", "", reasonHelp)
	_ = command.MarkFlagRequired("reason")
	if ackRequired {
		command.Flags().BoolVar(&acknowledge, "acknowledge-effects-unknown", false, "acknowledge that target effects may be unknown")
	}
	setHandoffFlagErrors(command)
	return command
}

func oneHandoffID(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return workflowhandoff.InvalidArguments(fmt.Sprintf("expected exactly one handoff ID, got %d arguments", len(args)))
	}
	_, err := workflowhandoff.ParseHandoffID(args[0])
	return err
}
func noArguments(extra func() error) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != 0 {
			return workflowhandoff.InvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
		}
		if err := extra(); err != nil {
			return workflowhandoff.InvalidArguments(err.Error())
		}
		return nil
	}
}
func validateHandoffLocator(value string) error {
	if value == "" {
		return fmt.Errorf("--handoff is required")
	}
	if !filepath.IsAbs(value) || filepath.Clean(value) != value {
		return fmt.Errorf("--handoff must be an absolute clean path")
	}
	return nil
}
func setHandoffFlagErrors(command *cobra.Command) {
	command.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error { return workflowhandoff.InvalidArguments(err.Error()) })
}
