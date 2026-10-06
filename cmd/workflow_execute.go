package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskexecute.Input
	var format string
	var target queueTargetFlags
	c := &cobra.Command{
		Use: "execute", Short: "Deliver the next goal from this Epic with its assigned agent",
		Long:    "From a Herdr terminal inside the registered return worktree, select the next eligible queued goal (default) or --spec, create its feature worktree, and start its assigned interactive Claude or Codex owner in Herdr. The owner defines the detailed solution and tests, handles review and fixes, then completes local integration after an actual candidate-bound human pass. New Claude launches request persistent folder trust for the Task worktree in Claude's configuration; --check previews the config path and project key without writing. Tool permissions and other native prompts remain under provider control. Repeated starts inspect the preserved execution; an unknown start never launches another agent. Older preserved launches retain their original trust and permission choices.",
		Example: "  ply workflow execute --check\n  ply workflow execute\n  ply workflow execute --spec explain-errors\n  ply workflow execute show wfr_<digest>",
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if err := executeFormat(format); err != nil {
				return err
			}
			if c.Flags().Changed("next") && (!in.Next || in.SpecID != "") {
				return workspace.WorkInvalidArguments("--next selects the next goal and cannot be combined with --spec")
			}
			if in.SpecID == "" {
				in.Next = true
			}
			var err error
			in.Target, err = target.input()
			if err != nil {
				return err
			}
			result, err := taskexecute.Execute(d, in)
			if outErr := writeExecuteResult(c, format, result); outErr != nil {
				return outErr
			}
			return workflowRunError(err)
		},
	}
	c.Flags().BoolVar(&in.Next, "next", false, "select the first eligible goal in this Epic's queue (default)")
	c.Flags().StringVar(&in.SpecID, "spec", "", "exact pending goal Spec ID in this Epic; ambiguous IDs are rejected")
	c.Flags().BoolVar(&in.Check, "check", false, "preview goal, base, worktree, runtime and Claude folder trust effect without writes or startup")
	target.bind(c)
	c.Flags().StringVar(&in.Runtime.HerdrWorkspace, "herdr-workspace", "", "Herdr workspace ID (default: current Herdr workspace)")
	c.Flags().StringVar(&in.Runtime.PermissionProfile, "permission-profile", "", "existing Codex permission profile; Claude uses native auto mode")
	c.Flags().StringVar(&in.NotificationContext, "notification-context", "", "optional existing owner notification context to freeze into the execution")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return workspace.WorkInvalidArguments(e.Error()) })
	addExecuteOwnerCommands(c, d)
	return c
}

func executeFormat(format string) error {
	if format != "text" && format != "json" {
		return workspace.WorkInvalidArguments("--format must be text or json")
	}
	return nil
}

func writeExecuteResult(c *cobra.Command, format string, result any) error {
	if format == "json" {
		b, e := taskrun.Canonical(result)
		if e != nil {
			return e
		}
		_, e = c.OutOrStdout().Write(append(b, '\n'))
		return e
	}
	switch v := result.(type) {
	case taskexecute.Result:
		if v.Goal != nil {
			fmt.Fprintf(c.OutOrStdout(), "Task: %s — %s\nReturn worktree: %s\nFeature worktree: %s\n", v.Goal.Goal.TaskID, v.Goal.Goal.Title, v.Goal.Target.ParentLocator, v.Goal.WorktreePath)
		}
		if v.Runtime != nil {
			fmt.Fprintf(c.OutOrStdout(), "Implementor: %s · model %s · effort %s\n", v.Runtime.Runtime.Provider, v.Runtime.Runtime.Model, v.Runtime.ReasoningEffort)
		}
		if trust := v.ClaudeProjectTrust; trust != nil && (v.Run == nil || v.Run.ClaudeTrust == nil) {
			config := trust.ConfigPath
			if config == "" {
				config = "unresolved"
			}
			fmt.Fprintf(c.OutOrStdout(), "Claude folder trust: planned persistent trust for %s\nClaude configuration: %s\n", trust.ProjectKey, config)
			if trust.Reason != "" {
				fmt.Fprintln(c.OutOrStdout(), "Trust configuration warning: "+trust.Reason)
			}
		}
		fmt.Fprintf(c.OutOrStdout(), "Execution: %s\n", v.State)
		if v.Run != nil {
			return writeExecuteResult(c, format, *v.Run)
		}
		if v.NextAction != "" {
			fmt.Fprintln(c.OutOrStdout(), "Next: "+v.NextAction)
		}
	case taskrun.WorkflowRun:
		fmt.Fprintf(c.OutOrStdout(), "Run: %s\nAgent: %s\nHerdr: %s / %s / %s\nTransport: %s (%s)\n", v.RunID, v.Provider, v.Transport.WorkspaceID, v.Transport.TabID, v.Transport.PaneID, v.Transport.State, v.Transport.Observation)
		if v.Delivery != nil {
			fmt.Fprintf(c.OutOrStdout(), "Delivery: %s\nRuntime permission: %s\n", v.Delivery.Phase, v.Delivery.PermissionState)
		}
		writeWorkflowClaudeTrust(c.OutOrStdout(), v.ClaudeTrust)
		for _, r := range v.Reasons {
			fmt.Fprintln(c.OutOrStdout(), "Needs attention: "+r.Detail)
		}
		fmt.Fprintf(c.OutOrStdout(), "Next: %s — %s\n", v.NextAction.Actor, v.NextAction.Message)
	}
	return nil
}

// writeWorkflowClaudeTrust renders only the performed effect. Locations and
// backup references remain available in JSON for inspection and recovery.
func writeWorkflowClaudeTrust(out io.Writer, trust *taskrun.ClaudeTrustEffect) {
	if trust == nil {
		return
	}
	fmt.Fprintf(out, "Claude trust: %s", trust.State)
	if trust.State != "written" && trust.State != "already" && trust.Reason != "" {
		fmt.Fprint(out, " — "+trust.Reason)
	}
	fmt.Fprintln(out)
}

func addExecuteOwnerCommands(parent *cobra.Command, d taskrun.Dependencies) {
	for _, verb := range []string{"show", "follow", "resume", "report", "verify", "qa", "integrate"} {
		verb := verb
		ownerCallback := verb != "show" && verb != "follow" && verb != "resume"
		var format, contextPath, file, review, evidence string
		var timeout int
		use := verb + " RUN_ID"
		if verb == "qa" {
			use += " pass|fail|blocked"
		}
		short := map[string]string{"resume": "Continue the same startup after native onboarding without restarting it", "show": "Read the preserved delivery and its next action", "follow": "Observe the same interactive delivery session", "report": "Preserve truthful progress, questions or incomplete results", "verify": "Run the owner's acceptance entrypoint and qualify this candidate", "qa": "Record an actual human answer for the exact candidate", "integrate": "Complete the authorized local return after a matching human pass"}[verb]
		detail := map[string]string{
			"resume":    "Continue an existing delivery startup after the human resolves native onboarding in its preserved tab. Verify the same provider and terminal, bind readiness, and send the first Task prompt only when it was never attempted. Never create another tab, start another agent, grant permissions or replay uncertain input.",
			"show":      "Read-only. Cached provider state is not a fresh observation or product approval.",
			"follow":    "The observation timeout does not stop the provider. Reuse the existing tab when startup or return is unknown.",
			"report":    "Owner callback: use the frozen control executable, exact Task cwd and --context. --file is a delivery-report@2; failed, not_run and unknown facts remain readable and do not qualify a candidate.",
			"verify":    "Owner callback: use the frozen control executable, exact Task cwd and --context. --review supplies an actual candidate-bound DeliveryCandidateReview@1 record. Execute the declared acceptance script and preserve real evidence. Missing or failed tests cannot produce a technical pass. Human product judgment remains separate.",
			"qa":        "Owner callback: use the frozen control executable, exact Task cwd and --context. --evidence preserves the actual user's answer and candidate binding. The answer is a local human attestation, not cryptographic identity proof. Never infer pass from tests, silence or a provider's claim.",
			"integrate": "Owner callback: use the frozen control executable, exact Task cwd and --context. Requires the exact qualified candidate and recorded human pass, unchanged clean parent, and the original local delivery authority. Preserves integration evidence, advances this preparation and updates the Epic base. Does not push, merge master or delete worktrees.",
		}[verb]
		example := "  ply workflow execute " + verb + " wfr_<digest>"
		if ownerCallback {
			example = "  /absolute/preserved/ply-control workflow execute " + verb + " wfr_<digest>"
		}
		if verb == "follow" || verb == "resume" {
			example += " --timeout 60"
		}
		if verb == "report" {
			example += " --context /absolute/context.json --file /absolute/report.json"
		}
		if verb == "verify" {
			example += " --context /absolute/context.json --review /absolute/review.json"
		}
		if verb == "qa" {
			example += " pass --context /absolute/context.json --evidence /absolute/human-answer.json"
		}
		if verb == "integrate" {
			example += " --context /absolute/context.json"
		}
		child := &cobra.Command{Use: use, Short: short, Long: short + ". " + detail, Example: example,
			Args: func(c *cobra.Command, args []string) error {
				want := 1
				if verb == "qa" {
					want = 2
				}
				if len(args) != want {
					return workspace.WorkInvalidArguments(fmt.Sprintf("expected %d positional arguments", want))
				}
				if err := executeFormat(format); err != nil {
					return err
				}
				if ownerCallback && contextPath == "" {
					return workspace.WorkInvalidArguments("--context is required for an owner callback")
				}
				if verb == "report" && file == "" || verb == "verify" && review == "" || verb == "qa" && evidence == "" {
					return workspace.WorkInvalidArguments("provide the required callback evidence file")
				}
				if verb == "qa" && args[1] != "pass" && args[1] != "fail" && args[1] != "blocked" {
					return workspace.WorkInvalidArguments("human outcome must be pass, fail or blocked")
				}
				if (verb == "follow" || verb == "resume") && (timeout < 1 || timeout > 86400) {
					return workspace.WorkInvalidArguments("--timeout must be between 1 and 86400 seconds")
				}
				return nil
			},
			RunE: func(c *cobra.Command, args []string) error {
				ws, err := d.Workflow.Workspace.ObserveContaining()
				if err != nil {
					return err
				}
				root := ws.Observation.Root
				var r taskrun.WorkflowRun
				switch verb {
				case "show":
					r, err = taskrun.WorkflowShow(d, root, args[0])
				case "follow":
					r, err = taskrun.WorkflowFollow(d, root, args[0], time.Duration(timeout)*time.Second)
				case "resume":
					r, err = taskrun.WorkflowResumeDeliveryStart(d, root, args[0], time.Duration(timeout)*time.Second)
				case "report":
					r, err = taskrun.WorkflowDeliveryReport(d, root, args[0], contextPath, file)
				case "verify":
					r, err = taskrun.WorkflowDeliveryVerify(d, root, args[0], contextPath, review)
				case "qa":
					r, err = taskrun.WorkflowDeliveryQA(d, root, args[0], contextPath, args[1], evidence)
				case "integrate":
					r, err = taskrun.WorkflowDeliveryIntegrate(d, root, args[0], contextPath)
				}
				if outErr := writeExecuteResult(c, format, r); outErr != nil {
					return outErr
				}
				return workflowRunError(err)
			},
		}
		child.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if verb == "follow" || verb == "resume" {
			child.Flags().IntVar(&timeout, "timeout", 60, "wait timeout in seconds; never stops the agent")
		}
		if ownerCallback {
			child.Flags().StringVar(&contextPath, "context", "", "exact private execution context (required)")
		}
		if verb == "report" {
			child.Flags().StringVar(&file, "file", "", "absolute immutable delivery report file (required)")
		}
		if verb == "verify" {
			child.Flags().StringVar(&review, "review", "", "absolute candidate-bound review evidence (required)")
		}
		if verb == "qa" {
			child.Flags().StringVar(&evidence, "evidence", "", "absolute evidence of the actual human answer (required)")
		}
		child.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return workspace.WorkInvalidArguments(e.Error()) })
		parent.AddCommand(child)
	}
}
