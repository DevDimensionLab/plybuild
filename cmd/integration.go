package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/integration"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newIntegrationCommand(d taskrun.Dependencies) *cobra.Command {
	service := integration.NewService(d, integration.Options{})
	var selection integration.Selection
	var format, actor, observation string
	var check bool
	c := &cobra.Command{Use: "integration", Short: "Review and integrate one exact Delivery, then safely close its Task", Long: "Start a human-owned integration from an existing checkout. Inspect the exact candidate, target, merge method, optional installation and cleanup plan. An actual pass confirms this plan; fail, blocked or EOF starts none of those effects. Local return uses the registered target, never the current branch. PR merge requires an already published Delivery and an explicit merge or squash method. Read commands and --check do not mutate files, refs or external services.", Example: "  ply integration list --project example\n  ply integration --delivery dlv_<digest> --check\n  ply integration --delivery dlv_<digest> --keep\n  ply integration --delivery dlv_<digest> --method squash\n  ply integration resume int_<digest>\n  ply integration --task legacy-task --reconcile --check", Args: cobra.NoArgs, PersistentPreRunE: func(c *cobra.Command, _ []string) error {
		if force, e := c.Flags().GetBool("force"); e == nil && force {
			return workspace.WorkInvalidArguments("--force cannot replace a human integration decision")
		}
		return nil
	}, RunE: func(c *cobra.Command, _ []string) error {
		if e := executeFormat(format); e != nil {
			return e
		}
		cwd, e := os.Getwd()
		if e != nil {
			return e
		}
		p, e := service.Preview(cwd, selection)
		if e != nil {
			return e
		}
		if check {
			return writeDeliveryOutput(c, format, p, integration.PreviewText(p))
		}
		if p.State != "ready" {
			if e = writeDeliveryOutput(c, format, p, integration.PreviewText(p)); e != nil {
				return e
			}
			return &taskrun.Error{Code: "integration_blocked", Detail: p.NextAction, Exit: 4}
		}
		fmt.Fprint(c.ErrOrStderr(), integration.PreviewText(p))
		fmt.Fprint(c.ErrOrStderr(), "This is your candidate judgment and confirmation of the displayed effects. Answer exactly pass, fail or blocked: ")
		started := time.Now().UTC().Format(time.RFC3339Nano)
		answer, e := readExactHumanLine(c.InOrStdin())
		if e != nil {
			return e
		}
		if answer != "pass" && answer != "fail" && answer != "blocked" {
			return workspace.WorkInvalidArguments("human_qa_required: answer must be exactly pass, fail or blocked; nothing was changed")
		}
		if observation == "" {
			observation = "The local caller answered " + answer + " after viewing this exact candidate and integration effect plan."
		}
		tr := p.Plan.TaskResult
		h := workflowhandoff.DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: string(tr.TaskID), TaskResultID: string(tr.ID), ResultOID: tr.ResultOID, ResultTree: tr.ResultTree, Outcome: answer, Answer: answer, ActorClaim: actor, StartSurface: "ply integration interactive", StartedAtUTC: started, CompletedAtUTC: time.Now().UTC().Format(time.RFC3339Nano), Observation: observation}
		r, e := service.Decide(cwd, p, h)
		if r.ID != "" {
			if outputErr := writeDeliveryOutput(c, format, r, integration.Text(r)); outputErr != nil {
				return outputErr
			}
		}
		if e != nil {
			return e
		}
		return integrationExit(r)
	}}
	c.Flags().StringVar(&selection.DeliveryID, "delivery", "", "exact registered Delivery ID; required when several exist")
	c.Flags().StringVar(&selection.TaskID, "task", "", "one exact legacy Task ID, only with --reconcile")
	c.Flags().BoolVar(&selection.Reconcile, "reconcile", false, "close an already integrated legacy Task using existing native evidence; never merge again")
	c.Flags().StringVar(&selection.Method, "method", "", "explicit PR merge method: merge or squash")
	c.Flags().BoolVar(&selection.Keep, "keep", false, "complete the Task while preserving its source worktree and local branch")
	c.Flags().StringVar(&selection.InstallProfilePath, "install-profile", "", "absolute path to an explicit candidate-bound installation profile; default none")
	c.Flags().StringVar(&selection.NotificationRoutePath, "notification-route", "", "absolute path to an explicit local notification route; default off")
	c.Flags().BoolVar(&check, "check", false, "show the exact plan and blockers without recording a decision or performing effects")
	c.Flags().StringVar(&actor, "actor", "Local integration caller", "actor claim preserved with the actual answer; not identity authentication")
	c.Flags().StringVar(&observation, "observation", "", "human observation to preserve with the answer")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	for _, verb := range []string{"list", "show", "scan", "resume"} {
		verb := verb
		var output, project, task string
		var preview, retryInstall, retryMerge bool
		use := verb
		if verb == "show" || verb == "resume" {
			use += " INTEGRATION_ID"
		}
		child := &cobra.Command{Use: use, Short: map[string]string{"list": "List registered Deliveries and preserved integrations", "show": "Read one integration's plan, decision, receipts and history", "scan": "Inspect active, integrated, retained, retired and unknown Task resources without writes", "resume": "Observe the same operation and complete only its remaining steps"}[verb], Args: func(c *cobra.Command, args []string) error {
			want := 0
			if verb == "show" || verb == "resume" {
				want = 1
			}
			if len(args) != want {
				return workspace.WorkInvalidArguments(fmt.Sprintf("expected %d arguments", want))
			}
			if preview && (retryInstall || retryMerge) || retryInstall && retryMerge {
				return workspace.WorkInvalidArguments("choose read-only --check or exactly one explicit retry effect")
			}
			return executeFormat(output)
		}, RunE: func(c *cobra.Command, args []string) error {
			cwd, e := os.Getwd()
			if e != nil {
				return e
			}
			switch verb {
			case "list":
				r, e := service.List(cwd, project, task)
				if e != nil {
					return e
				}
				var b strings.Builder
				for _, d := range r.Deliveries {
					fmt.Fprintf(&b, "%s  %s  %s  %s\n", d.ID, d.Manifest.TaskID, d.State, d.Manifest.Agreement.TargetRef)
				}
				for _, v := range r.Integrations {
					fmt.Fprintf(&b, "%s  %s  %s\n", v.ID, v.Plan.TaskResult.TaskID, v.State)
				}
				if b.Len() == 0 {
					b.WriteString("No Deliveries or integrations.\n")
				}
				return writeDeliveryOutput(c, output, r, b.String())
			case "scan":
				r, e := service.Scan(cwd, project)
				if e != nil {
					return e
				}
				var b strings.Builder
				for _, v := range r.Entries {
					fmt.Fprintf(&b, "%s  %s  %s\n  Next: %s\n", v.TaskID, v.State, v.Worktree, v.NextAction)
				}
				return writeDeliveryOutput(c, output, r, b.String())
			case "show":
				r, e := service.Read(cwd, args[0])
				if e != nil {
					return e
				}
				return writeDeliveryOutput(c, output, r, integration.Text(r))
			case "resume":
				r, e := service.Resume(cwd, args[0], preview, retryInstall, retryMerge)
				if r.ID != "" {
					if oe := writeDeliveryOutput(c, output, r, integration.Text(r)); oe != nil {
						return oe
					}
				}
				if e != nil {
					return e
				}
				if preview {
					return nil
				}
				return integrationExit(r)
			}
			return nil
		}}
		child.Flags().StringVar(&output, "format", "text", "output format (text or json)")
		if verb == "list" || verb == "scan" {
			child.Flags().StringVar(&project, "project", "", "filter by exact registered project ID")
		}
		if verb == "list" {
			child.Flags().StringVar(&task, "task", "", "filter by exact registered Task ID")
		}
		if verb == "resume" {
			child.Flags().BoolVar(&preview, "check", false, "read the preserved operation without writes or retries")
			child.Flags().BoolVar(&retryInstall, "retry-install", false, "explicitly retry only a known failed installation; unknown effects are observation only")
			child.Flags().BoolVar(&retryMerge, "retry-merge", false, "explicitly retry the unchanged PR plan only after exact proven rejection, fresh preflight and current human pass")
		}
		setWorkFlagErrors(child)
		c.AddCommand(child)
	}
	c.AddCommand(newIntegrationReleaseCommand(service))
	c.AddCommand(newIntegrationReconsiderCommand(service), newIntegrationNotificationRetryCommand(service))
	setWorkFlagErrors(c)
	return c
}

func newIntegrationReconsiderCommand(s *integration.Service) *cobra.Command {
	var check bool
	var format, actor, observation string
	c := &cobra.Command{Use: "reconsider INTEGRATION_ID", Short: "Revoke an unexecuted plan with an actual later fail or blocked answer", Long: "Preserve a later human negative answer only when this plan has positively produced no integration effect. Keep the original decision and effect reservation, record native revocation, and return correction ownership. Unknown, pending, merged, installed or closed operations cannot be cancelled.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		p, err := s.CheckReconsideration(cwd, args[0])
		if err != nil {
			return err
		}
		shown := fmt.Sprintf("Integration: %s\nCandidate: %s\nReconsideration allowed: %t\nReason: %s\nNext: %s\n", p.ID, p.CandidateOID, p.Allowed, p.Reason, p.NextAction)
		if check {
			return writeDeliveryOutput(c, format, p, shown)
		}
		if !p.Allowed {
			return &taskrun.Error{Code: "reconsider_forbidden", Detail: p.Reason, Exit: 4}
		}
		if strings.TrimSpace(observation) == "" {
			return workspace.WorkInvalidArguments("reconsideration requires --observation with the human's reason for fail or blocked")
		}
		r, err := s.Read(cwd, args[0])
		if err != nil {
			return err
		}
		fmt.Fprint(c.ErrOrStderr(), shown)
		fmt.Fprint(c.ErrOrStderr(), "Record your later judgment and revoke this unexecuted plan. Answer exactly fail or blocked: ")
		started := time.Now().UTC().Format(time.RFC3339Nano)
		answer, err := readExactHumanLine(c.InOrStdin())
		if err != nil {
			return err
		}
		if answer != "fail" && answer != "blocked" {
			return workspace.WorkInvalidArguments("negative_answer_required: answer must be exactly fail or blocked; nothing was changed")
		}
		tr := r.Plan.TaskResult
		h := workflowhandoff.DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: string(tr.TaskID), TaskResultID: string(tr.ID), ResultOID: tr.ResultOID, ResultTree: tr.ResultTree, Outcome: answer, Answer: answer, ActorClaim: actor, StartSurface: "ply integration reconsider interactive", StartedAtUTC: started, CompletedAtUTC: time.Now().UTC().Format(time.RFC3339Nano), Observation: observation}
		r, err = s.Reconsider(cwd, args[0], p.BasisSHA256, h)
		if r.ID != "" {
			if outputErr := writeDeliveryOutput(c, format, r, integration.Text(r)); outputErr != nil {
				return outputErr
			}
		}
		return err
	}}
	c.Flags().BoolVar(&check, "check", false, "inspect exact no-effect eligibility without writes or a human decision")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.Flags().StringVar(&actor, "actor", "Local integration caller", "actor claim preserved with the actual later answer")
	c.Flags().StringVar(&observation, "observation", "", "required human reason for rejecting or blocking this plan")
	setWorkFlagErrors(c)
	return c
}

func newIntegrationNotificationRetryCommand(s *integration.Service) *cobra.Command {
	var check, apply bool
	var format, confirmation string
	c := &cobra.Command{Use: "retry-notification INTEGRATION_ID", Short: "Inspect or explicitly retry one proven undelivered integration notification", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		if check && apply || !apply && confirmation != "" || apply && confirmation == "" {
			return workspace.WorkInvalidArguments("inspect first, then choose --apply --confirm with the exact notification retry digest")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if !apply {
			p, err := s.PreviewNotificationRetry(cwd, args[0])
			if err != nil {
				return err
			}
			shown := fmt.Sprintf("Integration: %s\nNotification event: %s\nRoute: %s (%s)\nRetry allowed: %t\nReason: %s\nNext: %s\n", p.ID, p.EventID, p.RouteName, p.ChannelLabel, p.Allowed, p.Reason, p.NextAction)
			return writeDeliveryOutput(c, format, p, shown)
		}
		r, err := s.RetryNotification(cwd, args[0], confirmation)
		if r.ID != "" {
			if outputErr := writeDeliveryOutput(c, format, r, integration.Text(r)); outputErr != nil {
				return outputErr
			}
		}
		return err
	}}
	c.Flags().BoolVar(&check, "check", false, "read-only inspection; also the default without --apply")
	c.Flags().BoolVar(&apply, "apply", false, "perform only the explicitly confirmed notification retry")
	c.Flags().StringVar(&confirmation, "confirm", "", "exact current retry digest returned by the read-only preview")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func readExactHumanLine(reader io.Reader) (string, error) {
	line, e := bufio.NewReader(io.LimitReader(reader, 4096)).ReadString('\n')
	if e != nil {
		return "", fmt.Errorf("human_qa_required: input ended without an exact complete answer; no effects started")
	}
	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}
func integrationExit(r integration.Receipt) error {
	if r.State != "completed" {
		return &taskrun.Error{Code: "integration_" + r.State, Detail: r.NextAction, Exit: 4}
	}
	return nil
}

func newIntegrationReleaseCommand(s *integration.Service) *cobra.Command {
	var deliveryID, taskID, format, actor string
	var check bool
	c := &cobra.Command{Use: "release", Short: "Preserve ownership release after a positively observed owner exit", Long: "For an existing Delivery, positively observe its native provider exit before releasing the exact candidate. An active, idle or unknown provider is not released. For one integrated legacy Task, check known owners and preserve the caller's explicit ownership release; this records neither new QA nor integration.", Args: cobra.NoArgs, RunE: func(c *cobra.Command, _ []string) error {
		if e := executeFormat(format); e != nil {
			return e
		}
		if (deliveryID == "") == (taskID == "") {
			return workspace.WorkInvalidArguments("choose exactly one --delivery or --task")
		}
		cwd, e := os.Getwd()
		if e != nil {
			return e
		}
		if check {
			p, e := s.Preview(cwd, integration.Selection{DeliveryID: deliveryID, TaskID: taskID, Reconcile: taskID != "", Keep: true})
			if e != nil {
				return e
			}
			return writeDeliveryOutput(c, format, p, integration.PreviewText(p))
		}
		if deliveryID != "" {
			r, e := s.ReleaseExited(cwd, deliveryID)
			if e != nil {
				return e
			}
			return writeExecuteResult(c, format, r)
		}
		p, e := s.Preview(cwd, integration.Selection{TaskID: taskID, Reconcile: true, Keep: true})
		if e != nil {
			return e
		}
		fmt.Fprintf(c.ErrOrStderr(), "Legacy Task %s, candidate %s, source %s. Confirm you have stopped all writers and release ownership of this exact source. Type released: ", taskID, p.Plan.TaskResult.ResultOID, p.Plan.TaskResult.SourceLocator)
		started := time.Now().UTC().Format(time.RFC3339Nano)
		answer, e := readExactHumanLine(c.InOrStdin())
		if e != nil {
			return e
		}
		r, e := s.ReleaseLegacy(cwd, taskID, actor, answer, started)
		if e != nil {
			return e
		}
		return writeDeliveryOutput(c, format, r, "Ownership release preserved. Preview the legacy reconciliation before confirming closeout.\n")
	}}
	c.Flags().StringVar(&deliveryID, "delivery", "", "Delivery whose exact provider exit can be observed")
	c.Flags().StringVar(&taskID, "task", "", "one already integrated legacy Task")
	c.Flags().StringVar(&actor, "actor", "Local integration caller", "actual caller's actor claim")
	c.Flags().BoolVar(&check, "check", false, "inspect release readiness without writes")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	return c
}

func init() {
	RootCmd.AddCommand(newIntegrationCommand(taskrun.SystemDependencies(workspace.SystemDependencies())))
}
