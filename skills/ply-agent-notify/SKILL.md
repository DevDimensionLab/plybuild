---
name: ply-agent-notify
description: Use Ply to notify the agreed Slack channel when an agent finishes, stops, or requires the user's answer. Preserve the local return or question and report delivery status separately. Excludes optional questions, internal corrections and crash monitoring.
---

Use this skill at the task's agreed completion or necessary-feedback point. The task
must already authorize outgoing notifications to the named route. This skill does
not grant that authority or require another approval for an already authorized call.

Preserve the actual local report before `agent_finished`. For a real stop use
`agent_stopped`; before-start stops omit a nonexistent receipt and report. Include an
available agreed start receipt. For `feedback_required`, preserve the actual necessary
question first, address the user, then notify before waiting. Optional clarification
and a test failure that can be corrected within the task send no notification.

Read [usage and invocation](references/usage.md) for the helper and event shape. Use
[the task context template](references/task-template.json) to bind activity/run,
handoff, route and one fixed event directory in the task mandate. Preserve the same
event file and ID for the same logical question across retries, process restarts and
port revisions. An ending that still waits for that answer uses
`--trigger waiting_for_feedback`; it does not create another stop event. A genuinely
new event gets a new ID. Never change an ID, route, or state directory to bypass unknown.

Call `scripts/notify.py` with structured file arguments. It preserves the request,
uses actual `ply workflow notification send --check`, applies at most once, and reads
`show`. It never evaluates question text or performs automatic retries. Keep the
agreed local return and user question intact even if notification fails. Report its
JSON `state` separately from the product outcome: an acknowledged transport does not
mean read, approved, controlled, or product QA passed. A coordinator continuation must
say “Coordinator continues”; ask the user to act only when the user is the next actor.

Credential setup and personal skill installation belong to the separately authorized
installation activity. Normal calls use the existing environment value or the one
named private file; never put a webhook in prompts, argv, reports or ordinary logs.
Do not set up an app, search secret stores, use another transport or ask for another
webhook on every task. The helper cannot report hard process death or a closed tab.
