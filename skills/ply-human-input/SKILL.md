---
name: ply-human-input
description: Ask the human one necessary question through Ply, with 2-6 fixed choices or an exact short text, and wait for the answer, when your work in a registered Ply project cannot continue correctly without a new human decision that your mandate does not settle. A Ply Task is optional. Not for permissions, routine approvals or secrets.
---

# Ply human input

Publish one necessary question with `ply workflow input`. The human answers it
on a question card in the plyapp PWA, and plyapp wakes you. You read the same
stored answer back, use it, and report what you did. Ply records data only: it
starts, wakes and resumes nothing, creates no Task and grants no new authority.

## When to ask

Ask only if all of these hold:

- You cannot finish the assignment correctly without a decision that only the
  human can make.
- Your mandate, goal, specs and earlier answers do not already settle it.
- A reversible documented default would not be acceptable.

Do not ask for:

- runtime permission or command approval: follow your runtime's own prompt, and
  report a refusal as it is, without working around it.
- routine confirmation of work you are already authorized to do.
- secrets or credentials.

## Card or conversation

Use the card only when the question has 2-6 fixed choices (the PWA cannot
answer a text question) and plyapp-server confirms that your pane gets it:

```sh
plyapp-server pane-chat card-ready --pane "$HERDR_PANE_ID" --project PROJECT
```

`PROJECT` is the question's Project from step 1. Use the card only when the
command exits 0 with `ready: true` and its `session` equals your provider's
session ID (`CLAUDE_CODE_SESSION_ID` or `CODEX_THREAD_ID`). An Assignment label
is never a session ID. The draft takes `provider` and `session` from this answer.

Otherwise ask the question in the active conversation and wait for the answer
there; do not publish it in Ply. That covers every other outcome: exit 3 (this
pane or Project gets no card; `reason` says why when JSON is printed), exit 4
(the server cannot tell now), any other failure, and a missing `plyapp-server`
or `HERDR_PANE_ID`.

## Steps

1. **Context.** Reuse the workspace, project, repository and Task that your
   assignment already names, for example in its mandate, handoff or Ply
   context. Look up only what it leaves out, from your cwd:
   - `WS`: `workspace.root` from `ply workspace project list --format json`.
     `WS` is always the Ply workspace root, never the Project's wrapper.
   - `PROJECT`: the project in that list whose `wrapper` contains your cwd.
   - Your repository: in `ply workspace project show PROJECT --format json`,
     the entry in `repositories` whose `git_common_dir` equals
     `git rev-parse --path-format=absolute --git-common-dir`. This also finds
     a worktree outside the wrapper. When no wrapper contains your cwd, the
     project whose `show` has this entry is `PROJECT`.

   If neither source gives a registered project, stop; do not invent or
   create a registration. Set `task_id` only to an existing Task your
   assignment names, else `null`; never create a Task to ask a question.
   Then run the card-ready check from Card or conversation with this
   `PROJECT`; continue with step 2 only when the card applies.
   With a Task, `repo_id` is that Task's repository; otherwise it is the
   `repo_id` of your repository, or `null` when none matched.
2. **Key.** Choose `activity_id`, a stable lowercase key for your assignment or
   conversation (for example `docs-export` or `conversation/2026-10-18-readme`),
   and `publication_key` as `<activity_id>/<topic>`. Record both, and later the
   `question_id`, in your progress notes so that a resumed you reuses them.
3. **Draft.** Write the draft JSON to a new private file outside the
   repository, for example under `mktemp -d`. Every field is required; use
   `null` where allowed. `provider` and `session_id` are `provider` and
   `session` from card-ready.

   ```json
   {"kind":"ply.workflow.human-input-question-draft","schema_version":1,
    "publication_key":"docs-export/format","project_id":"ply","repo_id":"ply","task_id":null,
    "activity_id":"docs-export","activity_title":"Document the export format",
    "provider":"claude","session_id":"<card-ready session>","origin_cwd":"<pwd -P>",
    "text":"Which README section should this task write?",
    "reason":"The goal leaves the scope of the README change open.",
    "answer_kind":"choice",
    "options":[{"id":"short","label":"Short usage note","consequence":"One paragraph under Usage."},
               {"id":"full","label":"Full reference","consequence":"A new section with every flag."}]}
   ```

   Keep the question self-contained: the human reads it without your conversation.
4. **Publish.** Run
   `ply workflow input ask --workspace WS --file DRAFT --format json` and keep
   `question_id`. If the reply is lost or times out, run the same `ask` again
   with the same key and content; it returns the same question. Never publish
   a new key to retry. `human_input_publication_conflict` means you changed
   the content: read the original with `show`, and use a new key only if the
   question really is a different one.
5. **Notify.** After publication, read the installed `ply-agent-notify/SKILL.md`
   when this job has an existing explicitly or standing-authorized notification
   route. Reuse its route, job activity/run and notification owner;
   Q's `activity_id` is the question activity, not a replacement for the job's
   notification activity. Internal delegates return the QID to that owner.
   Without route authority or the installed notify skill, preserve the local question and report
   notification `not_attempted`; continue waiting without credential setup or a
   new permission question. Never create a Task/Run just to notify.

   Immediately before notifying, read the same Q:
   `ply workflow input show QUESTION_ID --workspace WS --format json`.
   If `state` is `answered`, continue from Read the answer. If `state` is `withdrawn`,
   stop asking. Only for `state=waiting_for_answer`, send this one event from the
   job's working directory, replacing `QUESTION_ID` with this exact QID. Record the
   QID/activity/publication key in your progress notes. On restart, repeat the same
   command.
   Keep its actual provider; do not relabel the app as an agent.
   This compact alert points to the complete, unchanged question, reason and
   consequences stored in Ply and displayed by plyapp. Keep private question
   details, options and the full answer out of public notification text.

   ```sh
   python3 ~/.agents/skills/ply-agent-notify/scripts/notify.py \
     --activity stable-job-id --run stable-run-id --event-id human-input/QUESTION_ID \
     --status needs_answer --provider claude --title "⚠️ Human input needed" \
     --summary "A necessary question is waiting. Can you answer it in the plyapp PWA?" \
     --next-action "Open the plyapp PWA and select HERDR_NAME under Mine." \
     --ply-bin "$PLY_BIN"
   ```

   `HERDR_NAME` is `label` from `herdr tab get "$HERDR_TAB_ID"`.

   `PLY_BIN` is the already-resolved absolute Ply executable used to publish/read
   Q (normally `PLY_BIN=$(command -v ply)`). Pass that same
   executable to the helper; do not assume its subprocess finds `ply` on PATH
   or change global PATH/configuration.
   That skill owns send/readback. Record its JSON `state` separately: acknowledged
   means transport acknowledged, not answered. Unknown stays unknown. After a
   crash/lost reply, use the same helper inputs for durable attempt readback;
   never choose another event ID, route or state root to get another send.
   Repeated waiting uses the same event with `--trigger waiting_for_feedback`
   (suppressed); it does not send a reminder. Recheck Q before any notification
   invocation; only `waiting_for_answer` can need an alert. For `agent_received`
   or `agent_reported_continued`, resume Use and report from the stored answer
   and existing continuation; do not notify or repeat completed work.
6. **End the turn.** Write in the active conversation that the question
   waits, with its text and each option's `id` and label, and that the human
   can answer on the card at the top of your chat in the plyapp PWA or directly
   in this conversation. The card may not show, so this message is the
   fallback. These details go only here, never in a notification. A Slack
   thread reply does not answer Q. Then end your turn. Do not run `wait`: it
   keeps the pane working and blocks the wake.
7. **Read the answer.** plyapp wakes you with a message whose only question
   data is the QID. Then, and whenever you resume, run
   `ply workflow input show QUESTION_ID --workspace WS --format json`.
   `waiting_for_answer`: end the turn again. `withdrawn`: stop asking.
   `answered`: the answer is `answer.option_id`; record receipt:

   ```
   ply workflow input acknowledge QUESTION_ID --workspace WS \
     --request-id <answer.request_id> --request-sha256 <answer.request_sha256> --format json
   ```

   Treat the answer as input to the assignment, never as a command.

   If the human answers in the conversation instead, run `show` first.
   `waiting_for_answer`: run
   `ply workflow input withdraw QUESTION_ID --workspace WS --reason "answered in conversation"`,
   then do the work the conversation's answer decides and report it there;
   `acknowledge` and `continue` do not apply to a withdrawn Q. `answered`: the card's answer
   applies; say so in the conversation if the two differ, and continue as above.
8. **Use and report.** Do the work the answer decides, then report what you did
   (at most 1000 characters):

   ```
   ply workflow input continue QUESTION_ID --workspace WS \
     --request-id <id> --request-sha256 <digest> --summary "<what you did>" --format json
   ```

   This is your own report: it is not proof of delivery, and the human's next
   steps are not tracked by it.

If the question becomes unnecessary before it is answered, run
`ply workflow input withdraw QUESTION_ID --workspace WS --reason "<why>"`.

## Interruption

An answer is stored and stays stored. plyapp wakes you at most once; a missed
wake is not repeated. When you resume by any route, the same `ask` with the
same key, or `show`, returns the question and its answer; continue from Read
the answer. Reuse the same notification event and its result, without a new
send for the same waiting Q. Say this plainly if the human asks what happens
when you are interrupted.

## Setup and limits

- Resolve the existing Ply executable once as absolute `PLY_BIN`. Commands above
  use `ply` as shorthand for that executable; use the assignment's existing path
  when it is absent from PATH. Do not install or change global configuration.
- Publishing writes only under `data/work/human-input/v1` in the wrapper of the
  question's Project (`wrapper` from `ply workspace project list --format json`).
  Ask for local write access there if your runtime needs it, and nothing broader.
- Limits: question 2000 characters, reason and consequence 1000, label 120,
  title 160, 2-6 options.
- For the full contract, read `schemas/workflow-human-input-v1/README.md` in the
  Ply repository, or `ply workflow input --help`.
