package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func addBoundStartCommands(parent *cobra.Command, d taskrun.Dependencies) {
	var in taskexecute.BoundStartInput
	var id, format string
	var target queueTargetFlags
	preview := &cobra.Command{Use: "preview", Short: "Read an externally approvable initial Goal start without effects", Long: "Read-only version 1 bound-start preview. Shows the exact explicit Goal, queue/base, worktree, requested Codex runtime and policy, executables, Herdr target, trust and delivery agreement. No intent, reservation, trust grant, worktree, Run, agent or tab is created. The consumer supplies a stable request ID before approval and preserves the returned request unchanged. This binding does not authenticate the user or grant additional delivery authority. Only interactive Codex is supported; no --next, headless or startup recovery fallback.", Example: "  ply workflow execute preview --spec selected-goal --request-id 550e8400-e29b-41d4-a716-446655440000 --herdr-workspace w9 --format json", Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if err := executeFormat(format); err != nil {
				return err
			}
			var err error
			in.Target, err = target.input()
			if err != nil {
				return err
			}
			out, err := taskexecute.PreviewBoundStart(d, id, in)
			if e := writeBoundStart(c, format, out); e != nil {
				return e
			}
			return err
		}}
	preview.Flags().StringVar(&in.SpecID, "spec", "", "exact pending Goal Spec ID (required)")
	preview.Flags().StringVar(&id, "request-id", "", "consumer identity: 8-128 letters, digits, dots, underscores or hyphens (required)")
	preview.Flags().StringVar(&in.StartMode, "start-mode", "interactive", "bound mode; v1 supports only interactive Codex")
	preview.Flags().StringVar(&in.HerdrWorkspace, "herdr-workspace", "", "exact Herdr workspace (default: current Herdr workspace)")
	preview.Flags().StringVar(&in.PermissionProfile, "permission-profile", "", "existing Codex permission profile to bind")
	preview.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	target.bind(preview)
	parent.AddCommand(preview)
	var file, applyFormat string
	apply := &cobra.Command{Use: "apply", Short: "Start exactly one approved immutable Goal request", Long: "Apply only the unchanged request object from execute preview. Native locks recheck the full approved basis before new effects. The same identity and same request only read its retained attempt; a different request under that identity is rejected. Unknown or interrupted outcomes never automatically resume, retry, reopen, or create another agent. No fallback to ordinary execute. Explicit bound Herdr starts work from the registered return worktree without requiring HERDR_ENV; user authentication remains the consumer's responsibility.", Example: "  ply workflow execute apply --file /absolute/approved-request.json --format json", Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if e := executeFormat(applyFormat); e != nil {
				return e
			}
			r, err := taskexecute.ReadBoundStartRequest(file)
			if err != nil {
				return err
			}
			out, err := taskexecute.ApplyBoundStart(d, r)
			if e := writeBoundStart(c, applyFormat, out); e != nil {
				return e
			}
			return err
		}}
	apply.Flags().StringVar(&file, "file", "", "absolute immutable bound-start-request@1 JSON (required)")
	apply.Flags().StringVar(&applyFormat, "format", "text", "output format (text or json)")
	parent.AddCommand(apply)
	var root, statusFormat string
	status := &cobra.Command{Use: "status REQUEST_ID", Short: "Read a bound request after a lost reply without starting anything", Long: "Read-only lookup by the consumer's stable request ID and physical workspace root. Returns not_recorded, retained_attempt, existing_run or outcome_unknown. No intent alone does not prove no effect. Preserved Run state is historical evidence, not live provider observation. This command never starts, resumes, retries, repairs or reopens an attempt and does not need the original provider or installed control bytes.", Example: "  ply workflow execute status 550e8400-e29b-41d4-a716-446655440000 --workspace /physical/workspace --format json", Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if e := executeFormat(statusFormat); e != nil {
				return e
			}
			out, err := taskexecute.StatusBoundStart(d, root, args[0])
			if e := writeBoundStart(c, statusFormat, out); e != nil {
				return e
			}
			return err
		}}
	status.Flags().StringVar(&root, "workspace", "", "physical workspace root (required)")
	status.Flags().StringVar(&statusFormat, "format", "text", "output format (text or json)")
	parent.AddCommand(status)
	for _, c := range []*cobra.Command{preview, apply, status} {
		c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return fmt.Errorf("bound_start_invalid_arguments: %w", err) })
	}
}

func writeBoundStart(c *cobra.Command, format string, value any) error {
	if format == "json" {
		return writeExecuteResult(c, format, value)
	}
	switch v := value.(type) {
	case taskexecute.BoundStartPreview:
		fmt.Fprintf(c.OutOrStdout(), "Bound start preview: %s\nRequest: %s\n", v.State, v.Request.RequestID)
		if v.State == "ready" {
			fmt.Fprintf(c.OutOrStdout(), "Goal: %s · Task %s\nReturn: %s\nPlanned Task: %s\nRuntime: %s · %s · %s\nHerdr: %s\nBinding: %s\nUse --format json and preserve the complete request object for apply. Runtime acceptance is still pending.\n", v.Basis.Goal.Goal.Goal.SpecID, v.Basis.Goal.Goal.TaskID, v.Basis.Goal.Target.ParentLocator, v.Basis.Goal.WorktreePath, v.Basis.Runtime.Runtime.Provider, v.Basis.Runtime.Runtime.Model, v.Basis.Runtime.ReasoningEffort, v.Basis.Runtime.HerdrWorkspace, v.Request.Binding)
		}
	case taskexecute.BoundStartStatus:
		fmt.Fprintf(c.OutOrStdout(), "Bound start: %s\nRequest: %s\nExecution effect: %s\n", v.State, v.RequestID, v.ExecutionEffect)
		if v.Run != nil {
			fmt.Fprintf(c.OutOrStdout(), "Run: %s\n", v.Run.RunID)
		}
		fmt.Fprintln(c.OutOrStdout(), v.NextAction)
	}
	return nil
}
