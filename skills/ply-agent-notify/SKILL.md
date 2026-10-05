---
name: ply-agent-notify
description: Bind completion notifications at the start of human-requested agent work and use Ply to notify the agreed Slack channel when the owning agent finishes, stops, or needs the user's answer. Honor standing notification preferences; keep internal delegated work quiet and preserve delivery identity across continuations.
---

Use this skill when a human starts a concrete task, then at its completion or necessary-
feedback point. An explicit request or standing user/project preference for the named
route is sufficient authorization; do not require the user to repeat it in each
handoff. Merely loading this skill or starting an agent does not grant outgoing
authority when no such preference exists. A narrower current instruction takes priority.

At task start, bind one notification owner (normally root), the authority source,
activity/run, physical worktree, stable mandate, report location, route and event root.
Use the existing task context template; keep owner/authority notes in the mandate,
not extra fields in its strict JSON. A short task note can preserve a chat request;
no separate developer handoff or invented start receipt is required.

Internal subagents and Herdr tasks return to the owner without their own Slack
notification. Pass that ownership rule when delegating. A separately agreed delivery
boundary may have its own notifier, but only one owner sends for that boundary.
Continuation, a human's QA answer, and context restart reuse the existing binding;
they are not new tasks just because another agent turn starts.

The owner preserves the local result, then invokes the helper before its final
response or genuine wait. Do not notify for each technical step while the coordinator
can continue. A completed delivery is an event even when there is no remaining user
action; describe the actual next actor rather than inventing a new start/approval gate.

Preserve the actual local report before `agent_finished`. For a real stop use
`agent_stopped`; before-start stops omit a nonexistent receipt and report. Include an
available agreed start receipt. For `feedback_required`, preserve the actual necessary
question first, address the user, then notify before waiting. Optional clarification
and a test failure that can be corrected within the task send no notification.

Read [usage and invocation](references/usage.md) for the helper and event shape. Use
[the task context template](references/task-template.json) to bind activity/run,
handoff, route and one fixed event directory in the task mandate. Preserve the same
event file and ID for the same logical question across retries, process restarts and
port revisions, and likewise preserve an existing completion event. An ending that still waits for that answer uses
`--trigger waiting_for_feedback`; it does not create another stop event. A genuinely
new event gets a new ID. Never change an ID, route, or state directory to bypass unknown.

Call `scripts/notify.py` with structured file arguments. It preserves the request,
uses actual `ply workflow notification send --check`, applies at most once, and reads
`show`. It never evaluates question text or performs automatic retries. Keep the
agreed local return and user question intact even if notification fails. Report its
JSON `state` separately from the product outcome: an acknowledged transport does not
mean read, approved, controlled, or product QA passed. A coordinator continuation must
say “Coordinator continues”; ask the user to act only when the user is the next actor.

Runtime restrictions still apply. If the fixed notification state or transport cannot
be used, keep the local return and report the actual blocked/not-attempted/unknown
outcome. Do not change permissions, route, event identity or state root to bypass it.

Credential setup and personal skill installation belong to the separately authorized
installation activity. Normal calls use the existing environment value or the one
named private file; never put a webhook in prompts, argv, reports or ordinary logs.
Do not set up an app, search secret stores, use another transport or ask for another
webhook on every task. The helper cannot report hard process death or a closed tab.
