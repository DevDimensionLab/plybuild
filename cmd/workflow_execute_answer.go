package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

var answerCLIRunID = regexp.MustCompile(`^wfr_[a-f0-9]{64}$`)
var answerCLIQuestionID = regexp.MustCompile(`^pqw_[a-f0-9]{32}$`)
var answerCLIRequestID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)
var answerCLIOptionID = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

func answerCLIInvalid(detail string) error {
	return fmt.Errorf("bound_answer_invalid_arguments: %s", detail)
}

func answerCLIFormat(format string) error {
	if format != "text" && format != "json" {
		return answerCLIInvalid("--format must be text or json")
	}
	return nil
}

// CLI shape checks do not discover a workspace from cwd or create missing
// directories. The native operations independently revalidate under their own
// descriptor-anchored store/parent boundaries before an effect.
func answerCLIWorkspace(root string) error {
	if len(root) == 0 || len(root) > 4096 || !filepath.IsAbs(root) || filepath.Clean(root) != root || strings.ContainsRune(root, 0) {
		return answerCLIInvalid("--workspace must name the physical absolute normalized workspace")
	}
	physical, err := filepath.EvalSymlinks(root)
	if err != nil || physical != root {
		return answerCLIInvalid("--workspace must exist without symlink components")
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return answerCLIInvalid("--workspace must be a physical directory")
	}
	return nil
}

func configureAnswerCommand(c *cobra.Command, format *string) {
	c.SilenceUsage = true
	c.PersistentPreRunE = func(_ *cobra.Command, _ []string) error { return nil }
	c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return answerCLIInvalid(err.Error()) })
	c.Flags().StringVar(format, "format", "json", "output format (json or text; both preserve the complete structured contract)")
}

func addBoundAnswerCommands(parent *cobra.Command, d taskrun.Dependencies) {
	var workspace, format string
	question := &cobra.Command{
		Use: "question RUN_ID", Short: "Read a structured same-owner product question without effects",
		Long:    "Read the schema 2 week_codex_same_owner_turn_v1 product-question contract for one exact Run and physical workspace. Scope, parent, worker and fixed choices are separate from eligibility. Only a ready opted-in same-owner question can be previewed. Old Runs remain unsupported. Reading creates no intent, lock file, reservation, worker wake or provider turn. Retained evidence is not provider liveness or global writer quiescence.",
		Example: "  ply workflow execute question RUN_ID --workspace /physical/workspace --format json",
		Args: func(_ *cobra.Command, args []string) error {
			if err := answerCLIFormat(format); err != nil {
				return err
			}
			if len(args) != 1 || !answerCLIRunID.MatchString(args[0]) {
				return answerCLIInvalid("question requires one exact Run ID")
			}
			return answerCLIWorkspace(workspace)
		},
		RunE: func(c *cobra.Command, args []string) error {
			out, err := taskrun.ReadProductQuestion(d, workspace, args[0])
			return writeBoundAnswerResult(c, out.Kind, out, err)
		},
	}
	configureAnswerCommand(question, &format)
	question.Flags().StringVar(&workspace, "workspace", "", "physical absolute workspace root (required; never inferred from cwd)")
	parent.AddCommand(question)

	answer := &cobra.Command{
		Use: "answer", Short: "Preview, submit and read one bound same-owner product answer",
		Long: "The schema 2 same-owner answer contract preserves the existing Task, native reservation and exact Codex session. Apply accepts only the unchanged request from preview; it does not accept a prompt or general resume. Read and preview are effect-free. Keep the same request identity after an uncertain result.",
		Args: func(_ *cobra.Command, _ []string) error {
			return answerCLIInvalid("choose answer preview, apply or status")
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return answerCLIInvalid("choose answer preview, apply or status")
		},
		SilenceUsage: true,
	}
	answer.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return answerCLIInvalid(err.Error()) })

	var previewWorkspace, previewFormat, questionID, revision, optionID, requestID string
	preview := &cobra.Command{
		Use: "preview RUN_ID", Short: "Bind one fixed answer to the same-owner question without effects",
		Long:    "Read-only schema 2 same-owner answer preview. Supply the exact question ID and revision 1, one listed option and your stable request ID. A ready result carries the full unchanged request, opaque binding, explicit original deadline and at most five minutes of validity. Other results have request null. Parent pause may allow a deferred answer but answering never resumes the parent. No intent, file, reservation, wake or provider turn is created.",
		Example: "  ply workflow execute answer preview RUN_ID --workspace /physical/workspace --question-id QUESTION_ID --question-revision 1 --option-id task_first --request-id answer-choice-0001 --format json",
		Args: func(_ *cobra.Command, args []string) error {
			if err := answerCLIFormat(previewFormat); err != nil {
				return err
			}
			if len(args) != 1 || !answerCLIRunID.MatchString(args[0]) {
				return answerCLIInvalid("preview requires one exact Run ID")
			}
			if !answerCLIQuestionID.MatchString(questionID) || revision != "1" || !answerCLIOptionID.MatchString(optionID) || !answerCLIRequestID.MatchString(requestID) {
				return answerCLIInvalid("preview requires an exact question ID, revision 1, fixed option ID and stable request ID")
			}
			return answerCLIWorkspace(previewWorkspace)
		},
		RunE: func(c *cobra.Command, args []string) error {
			out, err := taskrun.PreviewBoundAnswer(d, previewWorkspace, args[0], questionID, 1, optionID, requestID)
			return writeBoundAnswerResult(c, out.Kind, out, err)
		},
	}
	configureAnswerCommand(preview, &previewFormat)
	preview.Flags().StringVar(&previewWorkspace, "workspace", "", "physical absolute workspace root (required)")
	preview.Flags().StringVar(&questionID, "question-id", "", "exact published pqw_ question ID (required)")
	preview.Flags().StringVar(&revision, "question-revision", "", "exact question revision; schema 2 requires 1 (required)")
	preview.Flags().StringVar(&optionID, "option-id", "", "one fixed option ID from the active question (required)")
	preview.Flags().StringVar(&requestID, "request-id", "", "stable consumer ID: 8-128 ASCII letters, digits, dots, underscores or hyphens (required)")
	answer.AddCommand(preview)

	var file, applyFormat string
	apply := &cobra.Command{
		Use: "apply", Short: "Record one approved same-owner answer with at most one continuation",
		Long:    "Apply only the complete unchanged schema 2 same-owner request returned by answer preview. The native parent/store transaction revalidates the approved basis before publishing one durable winner. Same ID and canonical request only read the same attempt; changed content conflicts. No blind retry, new session, new Task or parent resume is allowed. A nonzero exit can still carry valid JSON status: preserve the same ID and read status after a lost reply. Argument/file/parser failures have stable stderr codes and no invented stdout status.",
		Example: "  ply workflow execute answer apply --file /absolute/approved-answer.json --format json",
		Args: func(_ *cobra.Command, args []string) error {
			if err := answerCLIFormat(applyFormat); err != nil {
				return err
			}
			if len(args) != 0 || file == "" {
				return answerCLIInvalid("apply requires only --file with an approved answer request")
			}
			return nil
		},
		RunE: func(c *cobra.Command, _ []string) error {
			request, err := taskrun.ReadBoundAnswerRequest(file)
			if err != nil {
				return err
			}
			out, err := taskrun.ApplyBoundAnswer(d, request)
			return writeBoundAnswerResult(c, out.Kind, out, err)
		},
	}
	configureAnswerCommand(apply, &applyFormat)
	apply.Flags().StringVar(&file, "file", "", "absolute physical regular file containing the unchanged approved bound-answer-request schema 2 (maximum 64 KiB)")
	answer.AddCommand(apply)

	var statusWorkspace, statusFormat string
	status := &cobra.Command{
		Use: "status REQUEST_ID", Short: "Read the same-owner answer receipt after a lost reply without effects",
		Long:    "Read the schema 2 same-owner answer identity from its exact physical workspace. Distinguishes durable answer, deferred or reserved continuation, observed same-session start, positively consumed answer and unknown outcome. not_recorded means no readable request, not proof of no effect; unknown fields remain null. The approved basis is historical evidence. Status never repairs, wakes, consumes, retries or resumes anything and cannot release a Task for a new start.",
		Example: "  ply workflow execute answer status answer-choice-0001 --workspace /physical/workspace --format json",
		Args: func(_ *cobra.Command, args []string) error {
			if err := answerCLIFormat(statusFormat); err != nil {
				return err
			}
			if len(args) != 1 || !answerCLIRequestID.MatchString(args[0]) {
				return answerCLIInvalid("status requires one stable request ID")
			}
			return answerCLIWorkspace(statusWorkspace)
		},
		RunE: func(c *cobra.Command, args []string) error {
			out, err := taskrun.ShowBoundAnswer(d, statusWorkspace, args[0])
			return writeBoundAnswerResult(c, out.Kind, out, err)
		},
	}
	configureAnswerCommand(status, &statusFormat)
	status.Flags().StringVar(&statusWorkspace, "workspace", "", "physical absolute workspace root (required)")
	answer.AddCommand(status)
	parent.AddCommand(answer)
}

// A semantic error must not discard a known public result. Invalid arguments,
// files and syntax return before reaching this boundary and print no status.
func writeBoundAnswerResult(c *cobra.Command, kind string, value any, operationErr error) error {
	if kind == "" {
		if operationErr != nil {
			return operationErr
		}
		return fmt.Errorf("bound_answer_outcome_unknown: no trustworthy public result is available")
	}
	body, err := taskrun.Canonical(value)
	if err != nil || len(body)+1 > 256<<10 {
		return fmt.Errorf("bound_answer_outcome_unknown: public result cannot be encoded within 256 KiB; retain the request identity")
	}
	if _, err := c.OutOrStdout().Write(append(body, '\n')); err != nil {
		return err
	}
	return operationErr
}
