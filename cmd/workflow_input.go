package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/devdimensionlab/plybuild/internal/humaninput"
	"github.com/spf13/cobra"
)

func inputInvalid(detail string) error {
	return fmt.Errorf("%s: %s", humaninput.CodeInvalidArguments, detail)
}

type inputFlags struct {
	workspace, format string
}

// configureInput keeps every consumer command explicit: --workspace is never
// inferred from cwd and the JSON contract is identical in text and json.
func configureInput(c *cobra.Command, f *inputFlags, workspace bool) {
	c.SilenceUsage = true
	c.PersistentPreRunE = func(_ *cobra.Command, _ []string) error { return nil }
	c.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return inputInvalid(err.Error()) })
	c.Flags().StringVar(&f.format, "format", "json", "output format (json or text; both print the complete JSON contract)")
	if workspace {
		c.Flags().StringVar(&f.workspace, "workspace", "", "physical absolute Ply workspace root (required; never inferred from cwd)")
	}
}

func (f inputFlags) check(workspace bool) error {
	if f.format != "json" && f.format != "text" {
		return inputInvalid("--format must be json or text")
	}
	if workspace && f.workspace == "" {
		return inputInvalid("--workspace is required")
	}
	return nil
}

func exactArgs(n int, what string) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != n {
			return inputInvalid(what)
		}
		return nil
	}
}

// writeInput prints the public result, or the semantic refusal with the
// current question when known. Argument errors print nothing on stdout.
func writeInput(c *cobra.Command, workspace string, value any, err error) error {
	var refusal *humaninput.Error
	if err != nil && errors.As(err, &refusal) && refusal.Code != humaninput.CodeInvalidArguments {
		out := humaninput.ErrorResult{Kind: humaninput.ErrorKind, SchemaVersion: humaninput.SchemaVersion, Mode: humaninput.Mode, Code: refusal.Code, Detail: refusal.Detail, Question: refusal.Question, ObservedAtUTC: time.Now().UTC().Format(time.RFC3339Nano)}
		if workspace != "" {
			out.Workspace = &workspace
		}
		if refusal.Question != nil {
			out.QuestionID = &refusal.Question.QuestionID
		}
		value = out
	} else if err != nil {
		return err
	}
	body, e := humaninput.Encode(value)
	if e != nil {
		return e
	}
	if _, e = c.OutOrStdout().Write(append(body, '\n')); e != nil {
		return e
	}
	return err
}

func newWorkflowInputCommand(d humaninput.Dependencies) *cobra.Command {
	input := &cobra.Command{
		Use:   "input",
		Short: "Publish and answer necessary human input for an ordinary cooperative agent",
		Long:  "Mode cooperative_agent_v1, schema 1. An ordinary Claude or Codex agent publishes one necessary question (2-6 fixed choices or one exact short text) for a named activity in a registered project, with an optional existing Task. A configured app discovers questions with list/show, previews and applies one immutable answer, and reads status after a lost reply. The agent waits, acknowledges the exact answer and reports its use. Nothing here starts, wakes or resumes an agent, creates a Task, or grants Git, process or tool authority. Reads create no files or locks. The contract bundle is schemas/workflow-human-input-v1.",
		Example: "  ply workflow input ask --workspace /physical/workspace --file /absolute/question.json\n" +
			"  ply workflow input list --workspace /physical/workspace --project ply\n" +
			"  ply workflow input wait QUESTION_ID --workspace /physical/workspace --timeout-seconds 300",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return nil },
		SilenceUsage:      true,
		Args: func(_ *cobra.Command, _ []string) error {
			return inputInvalid("choose ask, list, show, answer, wait, acknowledge, continue or withdraw")
		},
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	}
	input.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return inputInvalid(err.Error()) })

	var askFlags inputFlags
	var askFile string
	ask := &cobra.Command{
		Use:     "ask",
		Short:   "Publish, or find again, one question from a stable publication key",
		Long:    "Read one closed ply.workflow.human-input-question-draft (schema 1) from an absolute physical file and publish it in the explicit workspace. The same publication_key with the same content returns the original question, also after a lost reply or from a resumed agent with another session; changed content under the same key is refused (publish a changed question under a new key). The draft's project, optional repository and optional existing Task are checked against the registers and origin_cwd must belong to the project. No transport, provider start, Task or registration change occurs.",
		Example: "  ply workflow input ask --workspace /physical/workspace --file /absolute/question.json --format json",
		Args:    exactArgs(0, "ask takes no positional arguments"),
		RunE: func(c *cobra.Command, _ []string) error {
			if e := askFlags.check(true); e != nil {
				return e
			}
			if askFile == "" {
				return inputInvalid("--file is required")
			}
			draft, e := humaninput.ReadDraft(askFile)
			if e != nil {
				return writeInput(c, askFlags.workspace, nil, e)
			}
			q, e := humaninput.AskDraft(d, askFlags.workspace, draft)
			return writeInput(c, askFlags.workspace, q, e)
		},
	}
	configureInput(ask, &askFlags, true)
	ask.Flags().StringVar(&askFile, "file", "", "absolute physical regular file with the question draft (maximum 256 KiB)")
	input.AddCommand(ask)

	var listFlags inputFlags
	var project, state, cursor string
	var limit int
	list := &cobra.Command{
		Use:     "list",
		Short:   "List a project's questions in every state, one page at a time",
		Long:    "Read-only discovery in immutable question-ID order. Every state is included: waiting_for_answer, answered, agent_received, agent_reported_continued and withdrawn; nothing is aged out. --limit is 1-20 (default 20). Pass next_cursor as --cursor until complete is true. complete means the page scan finished successfully, not a transactional snapshot; a new poll starts from the beginning. A failed scan returns an error and no page, never a partial list.",
		Example: "  ply workflow input list --workspace /physical/workspace --project ply --format json\n  ply workflow input list --workspace /physical/workspace --project ply --state waiting_for_answer --limit 20 --cursor NEXT_CURSOR",
		Args:    exactArgs(0, "list takes no positional arguments"),
		RunE: func(c *cobra.Command, _ []string) error {
			if e := listFlags.check(true); e != nil {
				return e
			}
			if project == "" {
				return inputInvalid("--project is required")
			}
			q := humaninput.ListQuery{ProjectID: project, Limit: limit}
			if c.Flags().Changed("limit") && limit == 0 {
				return inputInvalid("--limit must be between 1 and 20")
			}
			if state != "" {
				q.State = &state
			}
			if cursor != "" {
				q.Cursor = &cursor
			}
			out, e := humaninput.List(d, listFlags.workspace, q)
			return writeInputOrArgs(c, listFlags.workspace, out, e)
		},
	}
	configureInput(list, &listFlags, true)
	list.Flags().StringVar(&project, "project", "", "registered project ID (required)")
	list.Flags().StringVar(&state, "state", "", "optional state filter")
	list.Flags().IntVar(&limit, "limit", humaninput.DefaultListLimit, "rows per page, 1-20")
	list.Flags().StringVar(&cursor, "cursor", "", "next_cursor from the previous page")
	input.AddCommand(list)

	var showFlags inputFlags
	show := &cobra.Command{
		Use:     "show QUESTION_ID",
		Short:   "Show one question, its answer and its next actor without effects",
		Long:    "Read-only. Shows the immutable question and digest, the current state, the recorded answer and agent-reported receipt and continuation when present, and next_action for this question only.",
		Example: "  ply workflow input show QUESTION_ID --workspace /physical/workspace --format json",
		Args:    exactArgs(1, "show requires one question ID"),
		RunE: func(c *cobra.Command, args []string) error {
			if e := showFlags.check(true); e != nil {
				return e
			}
			q, e := humaninput.Show(d, showFlags.workspace, args[0])
			return writeInputOrArgs(c, showFlags.workspace, q, e)
		},
	}
	configureInput(show, &showFlags, true)
	input.AddCommand(show)

	answer := &cobra.Command{
		Use:   "answer",
		Short: "Preview, apply and read one immutable human answer",
		Long:  "preview returns the exact request (purpose human_input_answer_v1) without effects; apply records it as the question's one winner; status reads the outcome of the same (question ID, request ID) after a lost reply. CLI data alone does not attest a physical passkey.",
		Args: func(_ *cobra.Command, _ []string) error {
			return inputInvalid("choose answer preview, apply or status")
		},
		RunE:         func(_ *cobra.Command, _ []string) error { return nil },
		SilenceUsage: true,
	}
	answer.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return inputInvalid(err.Error()) })

	var previewFlags inputFlags
	var requestID, optionID, textFile string
	preview := &cobra.Command{
		Use:     "preview QUESTION_ID",
		Short:   "Prepare the exact immutable answer request without effects",
		Long:    "Read-only. Give your stable --request-id and exactly one of --option-id (choice) or --text-file (text: an absolute physical file whose exact UTF-8 bytes are the answer, at most 4000 codepoints and 16000 bytes, never trimmed or normalized). The result carries the complete request and request_sha256 (RFC8785). Keep those exact bytes and apply them unchanged. The request has no expiry; after an authoritative not_recorded status the same request may be applied again.",
		Example: "  ply workflow input answer preview QUESTION_ID --workspace /physical/workspace --request-id app-0001 --option-id full --format json\n  ply workflow input answer preview QUESTION_ID --workspace /physical/workspace --request-id app-0002 --text-file /absolute/answer.txt",
		Args:    exactArgs(1, "preview requires one question ID"),
		RunE: func(c *cobra.Command, args []string) error {
			if e := previewFlags.check(true); e != nil {
				return e
			}
			if (optionID == "") == (textFile == "") {
				return inputInvalid("give exactly one of --option-id or --text-file")
			}
			if requestID == "" {
				return inputInvalid("--request-id is required")
			}
			v := humaninput.AnswerValue{AnswerKind: humaninput.AnswerChoice, OptionID: &optionID}
			if textFile != "" {
				text, e := humaninput.ReadTextFile(textFile)
				if e != nil {
					return writeInput(c, previewFlags.workspace, nil, e)
				}
				v = humaninput.AnswerValue{AnswerKind: humaninput.AnswerText, Text: &text}
			}
			p, e := humaninput.PreviewAnswer(d, previewFlags.workspace, args[0], requestID, v)
			return writeInputOrArgs(c, previewFlags.workspace, p, e)
		},
	}
	configureInput(preview, &previewFlags, true)
	preview.Flags().StringVar(&requestID, "request-id", "", "stable consumer request ID: 8-128 ASCII letters, digits, dots, underscores or hyphens")
	preview.Flags().StringVar(&optionID, "option-id", "", "one option ID of a choice question")
	preview.Flags().StringVar(&textFile, "text-file", "", "absolute physical file with the exact text answer")
	answer.AddCommand(preview)

	var applyFlags inputFlags
	var applyFile string
	apply := &cobra.Command{
		Use:     "apply",
		Short:   "Record one previewed answer request as the question's one winner",
		Long:    "Apply only the complete unchanged request from answer preview. The same request again is an idempotent read of the same winner (state recorded, exit 0). Any other request, or the same request ID with other bytes, is refused without overwriting (state other_winner with winner). Refusals before any change are not_recorded with a reason. After a lost or interrupted reply, read answer status for the same request; never create a new request to retry. A nonzero exit still prints the status JSON when the request identity is known.",
		Example: "  ply workflow input answer apply --workspace /physical/workspace --file /absolute/answer-request.json --format json",
		Args:    exactArgs(0, "apply takes no positional arguments"),
		RunE: func(c *cobra.Command, _ []string) error {
			if e := applyFlags.check(true); e != nil {
				return e
			}
			if applyFile == "" {
				return inputInvalid("--file is required")
			}
			r, e := humaninput.ReadAnswerRequest(applyFile)
			if e != nil {
				return writeInput(c, applyFlags.workspace, nil, e)
			}
			s, e := humaninput.Apply(d, applyFlags.workspace, r)
			return writeStatus(c, s, e)
		},
	}
	configureInput(apply, &applyFlags, true)
	apply.Flags().StringVar(&applyFile, "file", "", "absolute physical regular file with the unchanged answer request (maximum 256 KiB)")
	answer.AddCommand(apply)

	var statusFlags inputFlags
	var statusQuestion string
	status := &cobra.Command{
		Use:     "status REQUEST_ID",
		Short:   "Read the authoritative outcome of one answer request without effects",
		Long:    "Read-only outcome of (--question-id, REQUEST_ID): recorded (this request is the answer), other_winner (another request is, see winner), not_recorded (the question is readable and has no answer from this request), unavailable (the question or workspace cannot be found) or unknown (stored state is unreadable or inconsistent). unavailable and unknown are not not_recorded and allow no new write.",
		Example: "  ply workflow input answer status app-0001 --question-id QUESTION_ID --workspace /physical/workspace --format json",
		Args:    exactArgs(1, "status requires one request ID"),
		RunE: func(c *cobra.Command, args []string) error {
			if e := statusFlags.check(true); e != nil {
				return e
			}
			if statusQuestion == "" {
				return inputInvalid("--question-id is required")
			}
			s, e := humaninput.Status(d, statusFlags.workspace, statusQuestion, args[0])
			if e != nil {
				return e
			}
			return writeStatus(c, s, nil)
		},
	}
	configureInput(status, &statusFlags, true)
	status.Flags().StringVar(&statusQuestion, "question-id", "", "the question the request answers (required)")
	answer.AddCommand(status)
	input.AddCommand(answer)

	var waitFlags inputFlags
	var timeout string
	wait := &cobra.Command{
		Use:     "wait QUESTION_ID",
		Short:   "Wait read-only for the question's answer",
		Long:    "Poll the question without a lock until it is answered or withdrawn, or until the technical --timeout-seconds (1-600) ends. outcome still_waiting after a timeout means keep waiting; it is not a human deadline, no new question is published and the assignment does not stop. Exit 0 for every outcome.",
		Example: "  ply workflow input wait QUESTION_ID --workspace /physical/workspace --timeout-seconds 300 --format json",
		Args:    exactArgs(1, "wait requires one question ID"),
		RunE: func(c *cobra.Command, args []string) error {
			if e := waitFlags.check(true); e != nil {
				return e
			}
			seconds, e := strconv.Atoi(timeout)
			if e != nil || seconds < 1 || seconds > humaninput.MaxWaitSeconds || strconv.Itoa(seconds) != timeout {
				return inputInvalid("--timeout-seconds must be an integer from 1 to 600")
			}
			w, e := humaninput.Wait(d, waitFlags.workspace, args[0], time.Duration(seconds)*time.Second)
			return writeInputOrArgs(c, waitFlags.workspace, w, e)
		},
	}
	configureInput(wait, &waitFlags, true)
	wait.Flags().StringVar(&timeout, "timeout-seconds", "", "technical wait bound in seconds, 1-600 (required)")
	input.AddCommand(wait)

	report := func(use, short, long, example string, summary bool) *cobra.Command {
		var f inputFlags
		var rid, sha, text string
		c := &cobra.Command{
			Use: use + " QUESTION_ID", Short: short, Long: long, Example: example,
			Args: exactArgs(1, use+" requires one question ID"),
			RunE: func(c *cobra.Command, args []string) error {
				if e := f.check(true); e != nil {
					return e
				}
				if rid == "" || sha == "" {
					return inputInvalid("--request-id and --request-sha256 of the recorded answer are required")
				}
				var q humaninput.Question
				var e error
				if summary {
					if !c.Flags().Changed("summary") {
						return inputInvalid("--summary is required")
					}
					q, e = humaninput.Continue(d, f.workspace, args[0], rid, sha, text)
				} else {
					q, e = humaninput.Acknowledge(d, f.workspace, args[0], rid, sha)
				}
				return writeInputOrArgs(c, f.workspace, q, e)
			},
		}
		configureInput(c, &f, true)
		c.Flags().StringVar(&rid, "request-id", "", "request_id of the recorded answer")
		c.Flags().StringVar(&sha, "request-sha256", "", "request_sha256 of the recorded answer")
		if summary {
			c.Flags().StringVar(&text, "summary", "", "what the agent did with the answer, at most 1000 codepoints")
		}
		return c
	}
	input.AddCommand(report("acknowledge", "Report that the agent received exactly the recorded answer",
		"The agent reports receipt of the recorded answer (state agent_received, provenance reported_by_agent). The request ID and digest must be the recorded winner's. This attests no new model turn, process or delivery. Repeating the same acknowledgement is idempotent.",
		"  ply workflow input acknowledge QUESTION_ID --workspace /physical/workspace --request-id ID --request-sha256 DIGEST", false))
	input.AddCommand(report("continue", "Report how the agent used the answer after doing the work",
		"After actually continuing its work with the answer, the agent reports what it did (state agent_reported_continued, provenance reported_by_agent). Requires a prior acknowledge of the same answer. This is not native process observation and not Task, Epic or main delivery. The same summary again is idempotent; another summary is refused.",
		"  ply workflow input continue QUESTION_ID --workspace /physical/workspace --request-id ID --request-sha256 DIGEST --summary 'Wrote the full version into README.md.'", true))

	var withdrawFlags inputFlags
	var reason string
	withdraw := &cobra.Command{
		Use:     "withdraw QUESTION_ID",
		Short:   "Withdraw an unanswered question without deleting it",
		Long:    "The agent ends a question it no longer needs (state withdrawn, terminal). The reason is shown as agent-reported text. An answered question cannot be withdrawn. The same reason again is idempotent.",
		Example: "  ply workflow input withdraw QUESTION_ID --workspace /physical/workspace --reason 'The goal already answers this.'",
		Args:    exactArgs(1, "withdraw requires one question ID"),
		RunE: func(c *cobra.Command, args []string) error {
			if e := withdrawFlags.check(true); e != nil {
				return e
			}
			if !c.Flags().Changed("reason") {
				return inputInvalid("--reason is required")
			}
			q, e := humaninput.Withdraw(d, withdrawFlags.workspace, args[0], reason)
			return writeInputOrArgs(c, withdrawFlags.workspace, q, e)
		},
	}
	configureInput(withdraw, &withdrawFlags, true)
	withdraw.Flags().StringVar(&reason, "reason", "", "why the question is no longer needed, at most 1000 codepoints")
	input.AddCommand(withdraw)
	return input
}

func isArgumentError(e error) bool {
	var x *humaninput.Error
	return errors.As(e, &x) && x.Code == humaninput.CodeInvalidArguments
}

func writeInputOrArgs(c *cobra.Command, workspace string, value any, e error) error {
	if isArgumentError(e) {
		return e
	}
	return writeInput(c, workspace, value, e)
}

func writeStatus(c *cobra.Command, s humaninput.AnswerStatus, e error) error {
	body, err := humaninput.Encode(s)
	if err != nil {
		return err
	}
	if _, err = c.OutOrStdout().Write(append(body, '\n')); err != nil {
		return err
	}
	return e
}
