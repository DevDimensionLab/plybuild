# Task binding and normal use

At human task start, check the current request and standing user/project preference
for the named Slack route. A standing choice is reusable authorization, not a reason
to ask again. Without one, do not send. Bind the single notification owner and that
source in the stable task mandate. For a chat task this may be a short local task
note; a separate agent handoff is not required.

Copy `task-template.json` to the task's agreed return directory and replace the
example activity/run and physical paths. Reuse an existing context for a continuation.
`notification_authorized` records the existing choice for that exact route; it is a
claim, not authentication. The task JSON has a strict schema: ownership and authority
notes belong in the mandate, not additional JSON fields.
The event directory and route state root must remain fixed across callers/restarts.
The event directory's parent must exist; the helper creates it privately if absent.
The task's completion procedure preserves its report, then calls the helper before
the human-readable final response. Its necessary-feedback procedure preserves the
question, calls the helper, then waits for the answer in the active task.

The root/coordinator owns the parent task's notification. Tell internal delegates to
return locally without sending Slack; opening a Herdr tab is not a separate human
notification request. If the user separately chooses a delegated notification boundary,
bind one owner for it and avoid sending the same completion again from root.

Reuse the same task context and event file when resuming, restarting, or receiving a
QA answer. Do not recreate a sent completion or still-pending question. A later real
completion after an answered question is a new event within the same activity/run.
Keep the original event/request bytes and let the helper read back an existing attempt.
Record transport status separately; a final local result is still usable if sending fails.

For an agent that the user asked to deliver work, `agent_finished` means that the
agreed delivery boundary is reached, not that each internal coding/review step ended.
Use `feedback_required` for an actual required human decision, and `agent_stopped` for
a real stop. Do not turn an authorized coordinator continuation into a human start gate.
When no human action is required, say so plainly in `next_action`; `next_actor` still
identifies the actor responsible for any remaining follow-up. The only supported
values are `user` and `coordinator`. For completed work with no required follow-up,
use `user` with, for example, `The requested work is complete. No action is required.`
This identifies the recipient of the return without inventing a new task or approval.

Preserve a strict event JSON file at the agreed return location, for example:

```json
{
  "kind": "PlyAgentNotificationEvent@1",
  "activity": "example/delivery/1",
  "run": "example-run-1",
  "event_id": "target-question-1",
  "event_type": "feedback_required",
  "phase": "after_start",
  "occurred_at": "2026-10-04T14:00:00Z",
  "public": {
    "task_title": "Choose the deployment target",
    "summary": "Which of the two agreed targets should I use?",
    "next_action": "Answer in the active agent task so work can continue.",
    "next_actor": "user"
  }
}
```

The activity/run must match the task context. Use the actual UTC event time without
fractional seconds; preserve it on replay. Public text is the only exported free
text. No mentions, webhook URLs, secrets, raw report text or transcript belong there.
`agent_finished` requires `--report` pointing to an actual JSON object with its
kind/activity/run. `--start-receipt` uses the same metadata binding when available.
`agent_stopped` permits missing report/receipt; only this event permits `before_start`,
which forbids a start receipt. Feedback requires `after_start`, `next_actor=user`
and an actual question containing `?`. Necessity remains the caller's responsibility.

```sh
python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/notify.py \
  --task /absolute/task/notify-task.json --event /absolute/task/question.json \
  --start-receipt /absolute/task/start-receipt.json

python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/notify.py \
  --task /absolute/task/notify-task.json --event /absolute/task/finished.json \
  --start-receipt /absolute/task/start-receipt.json --report /absolute/task/report.json
```

For an optional question, internal correction, or ending still waiting for the same
feedback, use `--trigger optional_feedback`, `--trigger internal_correction`, or
`--trigger waiting_for_feedback`. These produce `suppressed` without invoking Ply or
writing notification state. `--port-revision` is only a caller annotation: it never
changes the event identity or grants another attempt.

The helper returns JSON with `product_outcome=unchanged`, transport `state`, command
exits and a local delivery receipt when available. Exit 0 means acknowledged or
suppressed; 2 means not attempted; 5 means delivery was rejected, not sent,
rate-limited or is unknown. Always parse `state`, including for Ply exit 0. A helper
crash after its apply marker permits only readback; uncertainty is preserved. Do not
change the frozen request or bypass a conflict. Keep the question visible in the
terminal, and the local report usable, even when transport fails. Explicit Ply retry
remains a separate operation under its existing guards; this helper never retries.

# One-time installation and credential bootstrap

Only during the separately authorized installation activity:

```sh
make install-agent-notify
python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/bootstrap.py
```

Bootstrap prefers `PLY_SLACK_WEBHOOK_URL` if present, otherwise reuses its already
stored value or requests the existing webhook through hidden terminal input. An
invalid environment value is an error, with no file fallback. It creates only the
named route, private credential file and state directory; it sends nothing. It
preserves conflicting files and rejects symlinks, unsafe ownership and modes.

Normal route: `/Users/perottochristensen/.config/ply/notifications/ply-log/route.json`.
Credential: the `webhook` file alongside it (0600 in a 0700 directory).
State: `/Users/perottochristensen/.local/state/ply/agent-notifications/ply-log` (0700).
The helper passes the credential only to send/check children, never to show. Ply's
CLI still stores no secret. Historical QA fixtures are not operational state.

Fixture installation uses `make install-agent-notify AGENT_NOTIFY_DEST=/absolute/private/fixture/ply-agent-notify`.
Bootstrap accepts explicit `--config-dir` and `--state-root`; notify accepts the exact
`--credential-file` and a single `--ply-bin` executable. Keep all fixture paths under
the authorized test root. `make test-agent-notify` builds the actual CLI and its
network-free production-command driver, exercises the installed helper, and retains
its isolated fixture path. Fake transport proves technical behavior; actual skill
activation, useful text and Slack delivery require the later human agent journey.
