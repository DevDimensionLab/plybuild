package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
	"strings"
)

// ExitCode maps command errors at the executable boundary. Child process status
// remains a separate TaskRun fact and is never interpreted here.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var typed *taskrun.Error
	if errors.As(err, &typed) && typed != nil && typed.Exit >= 2 && typed.Exit <= 5 {
		return typed.Exit
	}
	return 1
}
func runUsage(c *cobra.Command, detail string) error {
	return &taskrun.Error{Code: "task_run_invalid_arguments", Detail: detail, Exit: 2}
}
func addWorkspaceTaskRunCommands(task *cobra.Command, w workspace.Dependencies) {
	task.AddCommand(newTaskRunCommand(taskrun.SystemDependencies(w)))
}
func newTaskRunCommand(d taskrun.Dependencies) *cobra.Command {
	run := &cobra.Command{Use: "run", Short: "Start one bounded Codex run and preserve its return", Long: "Preview an exact prepared Task, confirm one interactive or factory exec Codex start, and collect its typed return. Callbacks use a bound --context path; reported model may be unknown. Client exit and human-observed task inactivity are separate facts. Plan result control, human QA and integration are separate gates.", Example: "  ply workspace task run start --file /absolute/request.json\n  ply workspace task run show trn_<digest>", Args: func(c *cobra.Command, a []string) error {
		if len(a) != 0 {
			return runUsage(c, "expected a run subcommand")
		}
		return nil
	}, RunE: func(c *cobra.Command, a []string) error { return c.Help() }}
	run.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return runUsage(c, e.Error()) })
	for _, name := range []string{"start", "show", "collect", "accept", "report"} {
		name := name
		var file, format, confirm, contextPath, statusPath string
		var check, apply bool
		use := name
		if name != "start" {
			use += " RUN_ID"
		}
		short := map[string]string{"start": "Preview or confirm one bounded Codex start", "show": "Read preserved run facts without writing", "collect": "Check or recover publication of a received result", "accept": "Accept the exact run from its bound recipient session", "report": "Submit the final semantic report from its bound session"}[name]
		long := short + ". "
		switch name {
		case "start":
			long += "--check is the default and writes nothing. --apply requires --confirm with the exact preview digest. Request schema 1 requires text output and a foreground macOS terminal. Factory exec request schema 2 permits --format json without a terminal; stdout is one Ply result and provider streams stay in bound private files. The child inherits the terminal; result collection runs after it exits. Without a human task observation, exit 5 with task_run_task_status_unknown means collection is waiting for that observation. A reserved run or factory slot is never relaunched. Exec requires a hash-bound factory authorization and narrow managed policy; one native terminal turn, reaped exit zero and empty owned process group establish inactivity. Unknown completion exits 5 and preserves the environment."
		case "show":
			long += "Shows requested and reported model provenance, client process, task execution source, acceptance, delivery and collection separately. Unknown is a successful readback, not completion. No locks or recovery files are created."
		case "collect":
			long += "--check is the default. --apply --confirm publishes only an already received return with proven client quiescence and an inactive task observation after the report. Use --task-status with a private canonical task-run-task-status@1 JSON document built by a return helper from the bound run and an explicit human observation. This source is an attestation, not provider authentication or human QA. Check is read-only; apply confirms the status digest and journal tip. It never starts an agent. Exec uses runner-owned ProviderCompletion@1 and rejects --task-status, including manual overrides. Exec result and collect-preview use schema 3. Historical TaskResults remain immutable."
		case "accept":
			long += "Use --file with a private ply.workspace.task-run-acceptance document. Use --context for the exact private context path; PLY_TASK_RUN_CONTEXT remains supported and must agree when both are set. Acceptance@2 permits a null model; known model mismatch or unknown necessary authority cannot start. Bound cwd, executable, policy and run session are checked. Negative claims are preserved even when missing facts prevent a WF receipt. This is the recipient's first action before target writes."
		case "report":
			long += "Use --file with a private ply.workspace.task-run-report document after the last correction. Use --context for the bound private context path, or the matching PLY_TASK_RUN_CONTEXT environment value. Identical submissions are idempotent. Include the mandatory task-requirements artifact from the frozen Task Spec. Report received does not prove provider or process completion."
		}
		c := &cobra.Command{Use: use, Short: short, Long: long, Example: "  ply workspace task run " + name + map[string]string{"start": " --file /absolute/request.json --check", "show": " trn_<digest> --format json", "collect": " trn_<digest> --check", "accept": " trn_<digest> --context /absolute/run/tmp/context.json --file /absolute/acceptance.json", "report": " trn_<digest> --context /absolute/run/tmp/context.json --file /absolute/report.json"}[name]}
		c.Args = func(c *cobra.Command, a []string) error {
			want := 1
			if name == "start" {
				want = 0
			}
			if len(a) != want {
				return runUsage(c, fmt.Sprintf("expected %d positional arguments", want))
			}
			if format != "text" && format != "json" {
				return runUsage(c, "--format must be text or json")
			}
			if (name == "start" || name == "accept" || name == "report") && file == "" {
				return runUsage(c, "--file is required")
			}
			if c.Flags().Changed("check") && apply {
				return runUsage(c, "--check and --apply are mutually exclusive")
			}
			if confirm != "" && !apply || apply && confirm == "" {
				return runUsage(c, "--apply requires --confirm, and --confirm requires --apply")
			}
			if name == "start" && apply && format == "json" {
				request, err := taskrun.ReadRequest(file)
				if err != nil {
					return runUsage(c, "JSON apply requires a valid exec request: "+err.Error())
				}
				if request.SchemaVersion != 2 {
					return runUsage(c, "interactive apply requires text output")
				}
			}
			if c.Flags().Changed("context") && contextPath == "" {
				return runUsage(c, "--context must not be empty")
			}
			if c.Flags().Changed("task-status") && statusPath == "" {
				return runUsage(c, "--task-status must not be empty")
			}
			if f := c.Flag("force"); f != nil && f.Changed {
				return runUsage(c, "--force cannot authorize a Task run")
			}
			return nil
		}
		c.RunE = func(c *cobra.Command, a []string) error {
			d := d
			if c.Flags().Changed("context") {
				d.CallbackContext = &contextPath
			}
			if c.Flags().Changed("task-status") {
				d.TaskStatusPath = statusPath
			}
			var result any
			var e error
			if name == "start" {
				if apply {
					preview, err := taskrun.PreviewStart(d, file)
					if err != nil {
						return err
					}
					if format == "text" {
						if err = renderTaskRun(c, "text", preview); err != nil {
							return err
						}
					}
					result, e = taskrun.Start(d, file, confirm)
				} else {
					result, e = taskrun.PreviewStart(d, file)
				}
			} else {
				ws, err := d.Workflow.Workspace.ObserveContaining()
				if err != nil {
					return err
				}
				root := ws.Observation.Root
				switch name {
				case "show":
					result, e = taskrun.Show(d, root, a[0])
				case "collect":
					if apply {
						result, e = taskrun.Collect(d, root, a[0], confirm)
					} else {
						result, e = taskrun.PreviewCollect(d, root, a[0])
					}
				case "accept":
					result, e = taskrun.Accept(d, root, a[0], file)
				case "report":
					result, e = taskrun.SubmitReport(d, root, a[0], file)
				}
			}
			if result != nil {
				if writeErr := renderTaskRun(c, format, result); writeErr != nil {
					return writeErr
				}
			}
			return e
		}
		if name == "accept" || name == "report" {
			c.Flags().StringVar(&contextPath, "context", "", "exact absolute private callback context path (must agree with environment)")
		}
		if name == "collect" {
			c.Flags().StringVar(&statusPath, "task-status", "", "absolute private canonical JSON with an explicit human task observation")
		}
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if name == "start" || name == "collect" {
			c.Flags().BoolVar(&check, "check", true, "preview without writing (default)")
			c.Flags().BoolVar(&apply, "apply", false, "apply the exact confirmed preview")
			c.Flags().StringVar(&confirm, "confirm", "", "exact sha256 preview confirmation")
		}
		if name == "start" || name == "accept" || name == "report" {
			c.Flags().StringVar(&file, "file", "", "absolute physical JSON document path")
		}
		c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return runUsage(c, e.Error()) })
		run.AddCommand(c)
	}
	return run
}
func renderTaskRun(c *cobra.Command, format string, v any) error {
	if format == "json" {
		b, e := taskrun.Canonical(v)
		if e != nil {
			return e
		}
		_, e = c.OutOrStdout().Write(append(b, '\n'))
		return e
	}
	switch x := v.(type) {
	case taskrun.Preview:
		var draft struct {
			Goal struct {
				Title string `json:"title"`
			} `json:"goal"`
		}
		_ = json.Unmarshal(x.Request.HandoffDraft, &draft)
		fmt.Fprintf(c.OutOrStdout(), "Task: %s — %s\nSpec: %s revision %d\nBase: %s\nWorktree: %s\nProvider: %s; Requested model: %s\nBudget: initial execution + 3 correction rounds, 5400 active seconds, 2 environment measures\nEffects: bounded Task work; one run reservation and WF handoff; at most one qualified TaskResult.\nReturn: %s\n", x.Preparation.Plan.TaskID, draft.Goal.Title, x.Preparation.Plan.SpecID, x.Preparation.Plan.Spec.Revision, x.Preparation.Plan.ParentOID, x.Preparation.Plan.WorktreePath, x.Request.Runtime.Provider, x.Request.Runtime.Model, x.Paths.ReportPath)
		return renderRunNext(c, x.NextArgv, x.Reasons)
	case taskrun.Result:
		if x.Kind == "" {
			return nil
		}
		_, e := fmt.Fprintf(c.OutOrStdout(), "Run: %s\nLaunch: %s; acceptance: %s\nClient process: %s; delivery: %s; collection: %s\nNext action: %s\n", x.RunID, x.Launch.State, x.Acceptance.State, x.Process.State, x.Delivery.State, x.Collection.State, x.NextAction)
		if e != nil {
			return e
		}
		reported := "missing"
		if x.RuntimeFacts.Source != nil {
			reported = "unknown (recipient claim)"
			if x.RuntimeFacts.ReportedModel != nil {
				reported = *x.RuntimeFacts.ReportedModel + " (recipient claim; " + x.RuntimeFacts.ModelState + ")"
			}
		}
		fmt.Fprintf(c.OutOrStdout(), "Requested model: %s\nReported model: %s\nTask execution: %s (source: %s)\n", x.RuntimeFacts.RequestedModel, reported, x.TaskExecution.State, x.TaskExecution.Source)
		for _, r := range x.Reasons {
			fmt.Fprintf(c.OutOrStdout(), "Needs attention: %s\n", r.Detail)
		}
	case *taskrun.Result:
		if x != nil {
			return renderTaskRun(c, format, *x)
		}
	case taskrun.CollectPreview:
		fmt.Fprintf(c.OutOrStdout(), "Run: %s\nClient process: %s\nTask execution: %s (source: %s)\nCollection checks: %d unresolved reasons\n", x.RunID, x.Process.State, x.TaskExecution.State, x.TaskExecution.Source, len(x.Reasons))
		return renderRunNext(c, x.NextArgv, x.Reasons)
	}
	return nil
}
func renderRunNext(c *cobra.Command, argv []string, reasons []taskrun.Reason) error {
	for _, r := range reasons {
		fmt.Fprintf(c.OutOrStdout(), "Needs attention: %s\n", r.Detail)
	}
	if len(argv) > 0 {
		quoted := make([]string, len(argv))
		for i, s := range argv {
			quoted[i] = taskrun.ShellQuote(s)
		}
		_, e := fmt.Fprintf(c.OutOrStdout(), "Next command: %s\n", strings.Join(quoted, " "))
		return e
	}
	return nil
}
