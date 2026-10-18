package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteReviewCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryReviewRequestInput
	var format string
	c := &cobra.Command{Use: "review RUN_ID", Short: "Start an independent review round for the owner's current candidate", Long: "Coordinator operation from the containing workspace. Binds the exact clean committed candidate of this delivery Run as a new review round and writes a reviewer brief, a return template and a reviewer-only return context. No review Task, Spec, selection, queue move, preparation or owner callback is needed, and the owner keeps working. Repeating it for the same candidate returns the existing round; a corrected candidate gets a new round with earlier rounds preserved. It then starts the separate reviewer session through Ply in a new visible Herdr tab of the owner's workspace (or --herdr-workspace), with the Run's provider, model and permission mode in the Task worktree, binds the reviewer's observed native session (distinct from the owner's) and sends it the brief. Each start effect is reserved before it is sent and never replayed; an unknown start stays visible, and after inspection --new-round starts a fresh reviewer for the same candidate.", Example: "  ply workflow execute review wfr_<digest> --coordinator 'Actual coordinator'", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowDeliveryRequestReview(d, ws.Observation.Root, args[0], in)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&in.Coordinator, "coordinator", "", "actual review coordinator claim (required)")
	c.Flags().StringVar(&in.HerdrWorkspace, "herdr-workspace", "", "Herdr workspace for the reviewer tab (default: the owner's workspace)")
	c.Flags().BoolVar(&in.NewRound, "new-round", false, "replace an unknown reviewer start for the same candidate with a fresh round")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func newWorkflowExecuteReviewReturnCommand(d taskrun.Dependencies) *cobra.Command {
	var contextPath, file, format string
	c := &cobra.Command{Use: "review-return RUN_ID", Short: "Seal the reviewer's completed judgment for its round's candidate", Long: "Reviewer operation from inside the workspace. --context is the round context named in the brief; Run it from the reviewer pane Ply started: Ply checks that pane and its observed session. --file is an ordinary JSON file with decision (passed, changes_requested or blocked), summary, findings (severity, summary, optional evidence), evidence_gaps and optional reviewer_claim. A complete negative review is a completed review; truthful evidence gaps are kept and never block the return. Ply assigns finding IDs and ordering. The bound reviewer session is distinct from the owner's. Repeating the same file is safe; a different judgment for a sealed round is refused.", Example: "  ply workflow execute review-return wfr_<digest> --context /absolute/review/rounds/001/context.json --file /absolute/my-review.json", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowDeliveryReviewReturn(d, ws.Observation.Root, args[0], contextPath, file)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&contextPath, "context", "", "review round context from the brief (required)")
	c.Flags().StringVar(&file, "file", "", "absolute review return JSON file (required)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
