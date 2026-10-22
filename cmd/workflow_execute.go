package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskexecute.Input
	var format string
	var target queueTargetFlags
	var headless bool
	var boundControl bool
	var resume bool
	var resumeTimeout int
	var workerAttempts, workerTimeout, workerRetries int
	c := &cobra.Command{
		Use: "execute", Short: "Deliver the next goal from this Epic with its assigned agent",
		Long:    "From the registered return worktree, use --headless for a durable one-Spec Codex worker on macOS. It requires an explicitly automatic local Epic agreement, uses the assigned model and effort with existing account login, and performs developer, sandboxed tests, independent read-only review, bounded corrections and native local delivery. It has no terminal dependency and performs no remote push. Control it with workflow worker show/pause/stop/resume. Pause and stop finish the active attempt before preserving a boundary. Unknown attempts remain blocked. Without --headless, from a Herdr terminal, select the next eligible queued goal (default) or --spec, create its feature worktree, and start its assigned interactive Claude or Codex owner in Herdr. The owner defines the detailed solution and tests, handles review and fixes, and follows the frozen delivery agreement after candidate-bound acceptance under its explicit policy. New local Epic agreements may select automatic acceptance; an existing schema 1 local Epic run can record an explicit human policy choice with execute acceptance. Historical agreements retain their human gate without that supported choice. A structured agreement selects a pull request, local Epic return, or explicit local branch return; historical goals keep their original local contract. Repeat the same command to inspect the preserved execution. After begun work loses its agent tab, run --resume from the Task worktree (or --spec SPEC_ID --resume from the return worktree): Ply finds the same Run without its ID, and --check previews without writes or start. When the old pane, agent and every Task writer are positively absent in Herdr and the provider still has the session transcript, it starts the same Claude (--resume) or Codex (resume) conversation in one new visible tab, verifies the same native session, preserves a new callback control (closing the old one) and sends one continuation prompt; the Run, context, mandate, answers, candidates and evidence are unchanged. A live owner, a remaining writer, an unknown earlier effect, a missing session or an uncertain resume step stops with a reason and next step; repeating never duplicates a tab, provider start or prompt. Plain execute and --restart keep their meaning. Add --restart with --spec to recover an interrupted Codex startup before the Task prompt; Ply finds the attempt and preserves the Task, worktree and history. If no attempt exists, it performs the ordinary first start. A missing Task terminal can be replaced in the current Herdr workspace after fresh checks. Possible Task input or an uncertain new-tab creation prevents another start. --check previews all effects. New Claude launches request persistent folder trust for the Task worktree in Claude's configuration. Tool permissions and native prompts remain under provider control; preserved launches retain their original trust and permission choices.",
		Example: "  ply workflow execute --check\n  ply workflow execute\n  ply workflow execute --spec explain-errors\n  ply workflow execute --spec explain-errors --restart\n  ply workflow execute --spec explain-errors --restart --check\n  ply workflow execute --resume --check\n  ply workflow execute --spec explain-errors --resume",
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if resume {
				if in.Restart || headless || boundControl || c.Flags().Changed("next") || in.NotificationContext != "" || in.Runtime.PermissionProfile != "" {
					return workspace.WorkInvalidArguments("--resume keeps the Run's frozen choices; combine it only with --spec, --check, --herdr-workspace, --timeout and --format")
				}
				if resumeTimeout < 1 || resumeTimeout > 3600 {
					return workspace.WorkInvalidArguments("--timeout must be between 1 and 3600 seconds")
				}
				if err := executeFormat(format); err != nil {
					return err
				}
				result, err := taskexecute.ResumeOwner(d, taskexecute.OwnerResumeInput{SpecID: in.SpecID, Check: in.Check, Timeout: time.Duration(resumeTimeout) * time.Second, DestinationWorkspaceID: in.Runtime.HerdrWorkspace})
				if outErr := writeExecuteResult(c, format, result); outErr != nil {
					return outErr
				}
				return workflowRunError(err)
			} else if c.Flags().Changed("timeout") {
				return workspace.WorkInvalidArguments("--timeout requires --resume")
			}
			if boundControl {
				if !headless || in.SpecID == "" || in.Next {
					return workspace.WorkInvalidArguments("--bound-control requires --headless and an explicit --spec")
				}
				in.ControlContract = taskrun.NewControlContract("standalone", "", "")
			}
			if headless {
				if workerRetries < 1 || workerRetries > 1000 {
					return workspace.WorkInvalidArguments("--worker-retries must be between 1 and 1000")
				}
				in.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: workerAttempts, AgentTimeoutSeconds: workerTimeout, MaxProviderRetries: workerRetries}
			} else if c.Flags().Changed("worker-attempts") || c.Flags().Changed("worker-timeout") || c.Flags().Changed("worker-retries") {
				return workspace.WorkInvalidArguments("worker limits require --headless")
			}
			if err := executeFormat(format); err != nil {
				return err
			}
			if c.Flags().Changed("next") && (!in.Next || in.SpecID != "") {
				return workspace.WorkInvalidArguments("--next selects the next goal and cannot be combined with --spec")
			}
			if in.Restart && in.SpecID == "" {
				return workspace.WorkInvalidArguments("--restart requires an explicit --spec")
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
	c.Flags().BoolVar(&in.Restart, "restart", false, "recover this --spec's interrupted Codex startup before Task input, preserving its worktree and history")
	c.Flags().BoolVar(&resume, "resume", false, "resume the same begun owner after its tab closed: same Task, session and evidence; run from its Task worktree or with --spec from the return worktree")
	c.Flags().IntVar(&resumeTimeout, "timeout", 120, "seconds to observe the resumed session with --resume; never stops the agent")
	c.Flags().BoolVar(&headless, "headless", false, "start one durable Codex worker for an automatic local Epic goal without a terminal")
	c.Flags().BoolVar(&boundControl, "bound-control", false, "opt a new --headless --spec Run into bound runtime control v1; --check previews the frozen contract")
	c.Flags().IntVar(&workerAttempts, "worker-attempts", 3, "maximum development attempts for --headless (1-20)")
	c.Flags().IntVar(&workerTimeout, "worker-timeout", 1800, "seconds per headless provider attempt (10-86400)")
	c.Flags().IntVar(&workerRetries, "worker-retries", 20, "maximum persisted provider usage-limit retries (1-1000); no model or account switch")
	target.bind(c)
	c.Flags().StringVar(&in.Runtime.HerdrWorkspace, "herdr-workspace", "", "Herdr workspace ID (default: current Herdr workspace)")
	c.Flags().StringVar(&in.Runtime.PermissionProfile, "permission-profile", "", "existing Codex permission profile; Claude uses native auto mode")
	c.Flags().StringVar(&in.NotificationContext, "notification-context", "", "optional existing owner notification context to freeze into the execution")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return workspace.WorkInvalidArguments(e.Error()) })
	addExecuteOwnerCommands(c, d)
	addBoundStartCommands(c, d)
	addBoundAnswerCommands(c, d)
	c.AddCommand(newWorkflowExecuteReleaseCommand(d))
	c.AddCommand(newWorkflowExecuteInstallCLICommand(d))
	c.AddCommand(newWorkflowExecuteRecoverStartCommand(d))
	c.AddCommand(newWorkflowExecuteContinueCommand(d))
	c.AddCommand(newWorkflowExecuteAcceptanceCommand(d))
	c.AddCommand(newWorkflowExecuteBeginCommand(d))
	c.AddCommand(newWorkflowExecuteCheckInCommand(d))
	c.AddCommand(newWorkflowExecuteObserveStopCommand(d))
	c.AddCommand(newWorkflowExecuteResolveStopCommand(d))
	c.AddCommand(newWorkflowExecuteReconcileQuotaStopCommand(d))
	c.AddCommand(newWorkflowExecuteReconcileCapacityStopCommand(d))
	c.AddCommand(newWorkflowExecuteRetireCommand(d))
	c.AddCommand(newWorkflowExecuteBackupGoalCommand(d))
	c.AddCommand(newWorkflowExecuteReviewCommand(d))
	c.AddCommand(newWorkflowExecuteReviewReturnCommand(d))
	c.AddCommand(newWorkflowExecuteQARecoveryCommand(d))
	executeLogsToStderr(c)
	return c
}

// Every execute command can read or call back a report Run whose launch facts
// changed. Its warnings go to stderr, as in workflow run, so a JSON result stays
// the only value on stdout.
func executeLogsToStderr(c *cobra.Command) {
	for _, child := range c.Commands() {
		executeLogsToStderr(child)
	}
	if run := c.RunE; run != nil {
		c.RunE = func(c *cobra.Command, args []string) error {
			previousOutput := logrus.StandardLogger().Out
			logrus.SetOutput(c.ErrOrStderr())
			defer logrus.SetOutput(previousOutput)
			return run(c, args)
		}
	}
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
	case taskexecute.OwnerResumeResult:
		out := c.OutOrStdout()
		fmt.Fprintf(out, "Owner resume: %s\n", v.State)
		for _, choice := range v.Choices {
			fmt.Fprintf(out, "Choice: --spec %s · Task %s · %s\n", choice.SpecID, choice.TaskID, choice.Phase)
		}
		if r := v.Run; r != nil {
			fmt.Fprintf(out, "Task: %s · Spec %s\nTask worktree: %s\nRun: %s\nPhase: %s\n", r.TaskID, r.SpecID, r.TaskWorktree, r.RunID, r.Phase)
		}
		if p := v.Preview; p != nil {
			fmt.Fprintf(out, "Provider: %s · session %s · %s\n", p.Provider, p.NativeSessionID, v.ProviderSupport)
			if p.Reason != "" {
				fmt.Fprintf(out, "Reason: %s — %s\n", p.Reason, p.Detail)
			}
			if l := p.Liveness; l != nil {
				fmt.Fprintf(out, "Writer liveness (%s): %d remaining process(es); %s\n", l.Method, len(l.RemainingPIDs), l.Limitation)
			}
			if p.State == "ready" && v.Workflow == nil {
				fmt.Fprintf(out, "Preview only: old pane %s absent; one new tab is planned in Herdr workspace %s\nPlanned control: %s\n", p.PreviousTransport.PaneID, p.DestinationWorkspaceID, p.ControlExecutable.Path)
			}
		}
		if v.Workflow != nil {
			return writeExecuteResult(c, format, *v.Workflow)
		}
		if v.NextAction != "" {
			fmt.Fprintln(out, "Next: "+v.NextAction)
		}
	case taskexecute.RecoveryResult:
		fmt.Fprintf(c.OutOrStdout(), "Startup recovery: %s\n", v.State)
		if v.Preview != nil && v.Preview.State != "existing" && v.Run == nil {
			fmt.Fprintf(c.OutOrStdout(), "Run: %s\nReplacement Codex: %s\n", v.Preview.RunID, v.Preview.ProviderExecutable.Path)
			if v.Preview.State == "ready" {
				if v.Preview.TransportMode == "replace_missing_terminal" {
					fmt.Fprintf(c.OutOrStdout(), "Task terminal: missing; a new tab is planned in Herdr workspace %s\n", v.Preview.DestinationWorkspaceID)
				} else {
					fmt.Fprintln(c.OutOrStdout(), "Task terminal: existing idle shell will be reused")
				}
				fmt.Fprintf(c.OutOrStdout(), "Planned control: %s\n", v.Preview.ControlExecutable.Path)
			}
		}
		if v.Run != nil {
			return writeExecuteResult(c, format, *v.Run)
		}
		if v.NextAction != "" {
			fmt.Fprintln(c.OutOrStdout(), "Next: "+v.NextAction)
		}
	case taskexecute.Result:
		if v.Goal != nil {
			fmt.Fprintf(c.OutOrStdout(), "Task: %s — %s\nReturn worktree: %s\nFeature worktree: %s\n", v.Goal.Goal.TaskID, v.Goal.Goal.Title, v.Goal.Target.ParentLocator, v.Goal.WorktreePath)
			a := v.Goal.Delivery
			if v.Run != nil && v.Run.DeliveryStatus != nil && v.Run.DeliveryStatus.AcceptanceSelection != nil {
				// Display the native validated current choice without changing the
				// frozen Goal that the JSON response preserves separately.
				a = &v.Run.DeliveryStatus.AcceptanceSelection.EffectiveAgreement
			}
			if a != nil {
				fmt.Fprintf(c.OutOrStdout(), "Delivery mode: %s\nSource: %s\nTarget: %s", a.Mode, a.SourceRef, a.TargetRef)
				if a.TargetWorktree != "" {
					fmt.Fprintf(c.OutOrStdout(), " at %s", a.TargetWorktree)
				}
				if a.GitHubRepository != "" {
					fmt.Fprintf(c.OutOrStdout(), " in %s via %s", a.GitHubRepository, a.Remote)
				}
				gates := "exact candidate human pass"
				if a.AutomaticAcceptance() {
					gates = "exact candidate automatic pass and preserved review"
					if a.Acceptance.RequireHumanQA || a.HumanOwnedIntegration() {
						gates += ", exact candidate human pass"
					}
				}
				fmt.Fprintf(c.OutOrStdout(), "\nDelivery gates: %s, actual runtime authority\n", gates)
				if a.Mode == workspace.DeliveryPullRequest {
					fmt.Fprintln(c.OutOrStdout(), "Delivery effects: publish source branch and open/update PR; stop before merge")
				} else {
					fmt.Fprintln(c.OutOrStdout(), "Delivery effects: exact local return, registered base update and current Task queue closure")
				}
			}
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
		if v.Recovery != nil {
			return writeExecuteResult(c, format, *v.Recovery)
		}
		if v.Run != nil {
			return writeExecuteResult(c, format, *v.Run)
		}
		if v.NextAction != "" {
			fmt.Fprintln(c.OutOrStdout(), "Next: "+v.NextAction)
		}
	case taskrun.WorkflowRun:
		if v.Worker != nil {
			fmt.Fprintf(c.OutOrStdout(), "Run: %s\nWorker: %s; stage %s; attempt %d; generation %d; owner %s\n", v.RunID, v.Worker.Phase, v.Worker.Stage, v.Worker.Iteration, v.Worker.Generation, v.Worker.OwnerState)
			if v.Worker.ControlAction != "" {
				fmt.Fprintf(c.OutOrStdout(), "Control: %s requested at %s\n", v.Worker.ControlAction, v.Worker.ControlRequestedAt)
			}
			if v.Worker.Reason != "" {
				fmt.Fprintln(c.OutOrStdout(), v.Worker.Reason)
			}
		} else {
			fmt.Fprintf(c.OutOrStdout(), "Run: %s\nAgent: %s\nHerdr: %s / %s / %s\nTransport: %s (%s)\n", v.RunID, v.Provider, v.Transport.WorkspaceID, v.Transport.TabID, v.Transport.PaneID, v.Transport.State, v.Transport.Observation)
		}
		if v.Delivery != nil {
			fmt.Fprintf(c.OutOrStdout(), "Delivery: %s\nRuntime permission: %s\n", v.Delivery.Phase, v.Delivery.PermissionState)
		}
		writeDeliveryStatus(c.OutOrStdout(), v.DeliveryStatus)
		if r := v.OwnerResume; r != nil {
			fmt.Fprintf(c.OutOrStdout(), "Owner resume: generation %d · step %s · same session observed %t · owner callback %s\nResumed control: %s\n", r.Generation, r.Step, r.SessionObserved, r.OwnerCallback, r.ControlExecutable.Path)
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
		var format, contextPath, file, review, evidence, reuse, candidateBinary string
		var timeout int
		var incomplete bool
		use := verb + " RUN_ID"
		if verb == "qa" {
			use += " pass|fail|blocked"
		}
		short := map[string]string{"resume": "Continue the same startup after native onboarding without restarting it", "show": "Read the preserved delivery and its next action", "follow": "Observe the same interactive delivery session", "report": "Preserve truthful progress, questions or incomplete results", "verify": "Run the owner's acceptance entrypoint and qualify this candidate", "qa": "Record an actual human answer for the exact candidate", "integrate": "Complete the authorized local return after matching acceptance"}[verb]
		detail := map[string]string{
			"resume":    "Continue an existing delivery startup after the human resolves native onboarding in its preserved tab. Verify the same provider and terminal, bind readiness, and send the first Task prompt only when it was never attempted. Never create another tab, start another agent, grant permissions or replay uncertain input.",
			"show":      "Read-only. Cached provider state is not a fresh observation or product approval.",
			"follow":    "The observation timeout does not stop the provider. Reuse the existing tab when startup or return is unknown.",
			"report":    "Owner callback: use the frozen control executable, exact Task cwd and --context. --file is a delivery-report@2; failed, not_run and unknown facts remain readable and do not qualify a candidate. If a runtime dependency blocks callbacks, the installed CLI's --incomplete path accepts only stopped/needs_input from the same live owner. It preserves the original context, artifacts and unresolved attempts without granting runtime authority or switching controls.",
			"verify":    "Owner callback: use the frozen control executable, exact Task cwd and --context. Without --review, Ply uses the completed passed review round for exactly this candidate (see execute review); --review supplies an actual candidate-bound DeliveryCandidateReview@1 record instead. Execute the declared acceptance script and preserve real evidence. Native verification assembles and preserves task-requirements from the frozen Spec requirements and this exact candidate's acceptance and review evidence. Do not hand-author a second requirements/status artifact for that managed ID; cite candidate-review for observations preserved in the review. Managed artifacts become available during successful candidate verification. Keep actual evidence and truthful requirement outcomes; a progress report or guide does not qualify a candidate. An explicit automatic agreement requires --candidate-binary and freezes the executable, test environment and policy timeout into the attempt. A successful test plus review supplies machine acceptance; required human QA and later negative human judgments remain separate gates. Missing, failed or unknown tests cannot produce a pass. Repeating the same reserved invocation recovers its completed receipt without starting another test; unknown execution is blocked.",
			"qa":        "Owner callback: use the frozen control executable, exact Task cwd and --context. --evidence preserves the actual user's answer and candidate binding. The answer is a local human attestation, not cryptographic identity proof. Never infer pass from tests, silence or a provider's claim.",
			"integrate": "Owner callback: use the frozen control executable, exact Task cwd and --context. Requires the exact qualified candidate, its frozen automatic or human acceptance gate, unchanged clean parent, and the original local delivery authority. Preserves integration evidence, advances this preparation and updates the registered base. Local main/master return requires an explicit local_branch_integration contract and matching runtime authority. PR mode cannot integrate locally. Does not push or delete worktrees.",
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
				if verb == "report" && file == "" || verb == "qa" && evidence == "" {
					return workspace.WorkInvalidArguments("provide the required callback evidence file")
				}
				if verb == "verify" && reuse != "" && candidateBinary != "" {
					return workspace.WorkInvalidArguments("--reuse preserves the original binary; omit --candidate-binary")
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
					if incomplete {
						r, err = taskrun.WorkflowDeliveryIncompleteReport(d, root, args[0], contextPath, file)
					} else {
						r, err = taskrun.WorkflowDeliveryReport(d, root, args[0], contextPath, file)
					}
				case "verify":
					if reuse != "" {
						r, err = taskrun.WorkflowDeliveryRequalify(d, root, args[0], contextPath, review, reuse)
					} else {
						r, err = taskrun.WorkflowDeliveryVerifyWithBinary(d, root, args[0], contextPath, review, candidateBinary)
					}
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
			child.Flags().BoolVar(&incomplete, "incomplete", false, "record only stopped/needs_input through the installed reader; no runtime authority or control transition")
		}
		if verb == "verify" {
			child.Flags().StringVar(&review, "review", "", "absolute candidate-bound review evidence (default: the passed native review round for this candidate)")
			child.Flags().StringVar(&candidateBinary, "candidate-binary", "", "absolute built candidate executable; required by explicit automatic acceptance")
			child.Flags().StringVar(&reuse, "reuse", "", "requalify this exact preserved verification attempt without executing acceptance again")
		}
		if verb == "qa" {
			child.Flags().StringVar(&evidence, "evidence", "", "absolute evidence of the actual human answer (required)")
		}
		child.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return workspace.WorkInvalidArguments(e.Error()) })
		parent.AddCommand(child)
	}
}
