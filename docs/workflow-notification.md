# Report-ready notifications

## Agent completion notifications

The installed `ply-agent-notify` skill also supports `agent_finished`, `agent_stopped`
and `feedback_required` events through the same notification CLI. At the start of
human-requested work, bind the standing or explicit authorization for the named route,
one notification owner, and a stable task context. A standing preference does not
need to be repeated in every handoff. Without that authorization, the skill sends nothing.

Root normally owns the notification. Internal delegates, including Herdr tasks, return
to root without sending their own message. At the agreed delivery boundary root
preserves the report and invokes the helper before returning to the user. Necessary
human questions are preserved and notified before waiting; optional questions and
internal corrections stay local. Continuation and restart retain the task/event
identity so an existing attempt is read back rather than sent again.

See [the skill](../skills/ply-agent-notify/SKILL.md) and
[its invocation contract](../skills/ply-agent-notify/references/usage.md). The helper
preserves the local outcome when Slack is unavailable. Transport acknowledgement
is separate from human acceptance; a skill cannot detect a dead agent or closed tab.

## Report-ready human gates

`ply workflow notification` sends an explicitly authorized report-ready signal to one
fixed Slack channel. The procedure owner first preserves the report and local return,
then checks that the gate still needs a human to start result control. An active
coordinator continues its authorized work without sending a human-gate notification.

The source, actor, gate, worktree and channel are claims. Hashes identify exact files;
they do not authenticate the sender or authorize a send. A transport acknowledgement
means Slack accepted the POST. It does not attest product correctness, report control,
reading, human QA, or authority for another activity. Local return instructions remain
usable after every notification outcome. The CLI never executes `next_action`.

## Fixed route setup

When separately authorized to set up the route, create a Slack Incoming Webhook for
the chosen fixed channel. Put its URL in `PLY_SLACK_WEBHOOK` using your usual secret
practice. Keep the URL out of route files, command arguments, source control and logs.
Only `send` and `retry` read that one named environment variable. `show` needs none.
The channel label is a claim; this CLI does not verify membership or channel identity.

Save the route below as `slack-route.json`. Use the same dedicated physical state root
for this route on every invocation. Its parent must already exist. A new root is
created only by apply with mode 0700; state files use 0600. Existing state must be
private and owned by the current user. Symlinks are refused. Unix filesystems with
file and directory sync and advisory locks are supported; Windows fails closed.

```json
{
  "kind": "ply.workflow.notification-route",
  "schema_version": 1,
  "name": "report-ready",
  "transport": "slack_incoming_webhook",
  "channel_label": "#workflow-reports",
  "webhook_env": "PLY_SLACK_WEBHOOK",
  "state_root": "/private/tmp/ply-notification-example/notification-state"
}
```

## Complete synthetic external request

The following JSON is the validator-tested example in
`test/fixtures/notification-docs/request.json`. Its matching handoff, start receipt and
report are in the same fixture directory. Copy those exact fixture files into
`/private/tmp/ply-notification-example` and create its `worktree` directory to reproduce
this synthetic example. This example does not authorize a real send. For a real report,
use its absolute physical paths, exact SHA-256 file digests, activity/run and expected
start/report kinds. Report contents are never copied into the Slack message.

Save the request as `report-ready.json`:

```json
{
  "kind": "ply.workflow.notification-request",
  "schema_version": 1,
  "route": "report-ready",
  "source": {
    "kind": "external",
    "activity": "example/report-ready",
    "run": "example-delivery-001",
    "worktree": "/private/tmp/ply-notification-example/worktree",
    "handoff": {
      "path": "/private/tmp/ply-notification-example/handoff.md",
      "sha256": "sha256:c0647fa48ec41ea53dcf95c11628d7471241d4d12b9449c55ec17b8f028b206a"
    },
    "start_receipt": {
      "path": "/private/tmp/ply-notification-example/start-receipt.json",
      "sha256": "sha256:5a00a3af794163b5b2093954b7ffa4a602d50b952dade47fa2ecaed1d0439ea4",
      "kind": "ExampleDeliveryStart@1"
    },
    "report": {
      "path": "/private/tmp/ply-notification-example/report.json",
      "sha256": "sha256:733ad44e4e463dd6d61510f1e0fc429254f7749bf1a6cbdd6bcd5ccafcd5ba61",
      "kind": "ExampleDeliveryReport@1"
    }
  },
  "gate": {
    "id": "result-control",
    "revision": 1,
    "opened_at": "2026-10-03T19:00:00Z",
    "reason": "result_control",
    "state": "waiting_for_human"
  },
  "sender": {
    "actor_claim": "delivery agent"
  },
  "public": {
    "task_title": "Example delivery",
    "next_action": "Start the planning launcher on your development machine for result control."
  }
}
```

The handoff is at most 1 MiB; start and report are JSON objects at most 16 MiB each.
Request JSON is at most 64 KiB and route JSON at most 16 KiB. Duplicate or unknown
request/route fields, invalid Unicode, non-integer numbers, control characters,
mentions and webhook values in retained fields are rejected. All sources are read-only.

Native requests instead use `source.kind="native"`, include `handoff_id` and `task_id`
(string, or null for a handoff without a Task basis), and omit `kind` from the start and
report locators. The locators must match `workflow handoff inspect` exactly. Native
reported failures or evidence gaps can still need result control; notification does
not convert them into success or create a Task, TaskResult, Journal entry or run.

## Preview, apply, and readback

```sh
ply workflow notification send --file ./report-ready.json --route ./slack-route.json --check --format json
ply workflow notification send --file ./report-ready.json --route ./slack-route.json --apply --confirm 'sha256:<confirmation-from-preview>' --format json
ply workflow notification show '<notification-id>' --route ./slack-route.json --format json
```

Read `confirmation` and `notification_id` from the first JSON object. Substitute those
exact values in the later commands. Check is the default and performs no writes,
locks or HTTP calls. It shows the complete payload and planned effects. Text mode is
the default; every leaf also accepts `--format json`. `--check` and `--apply` are
exclusive. Confirmation is invalid without apply. There is no force or endpoint flag.

Immediately before every apply or retry, the procedure owner checks that the port
still requires human action and has not been taken over or closed. Confirmation binds
the exact request/source/route/credential basis; it is not a human approval token.
Each new apply durably reserves one attempt before at most one POST. Repeating an
identical send with its original confirmation reads the existing receipt without
another POST. A source, path, route, public-text or credential change conflicts even
when the environment variable name stays the same.

A new request prints exactly this six-line message, without a final newline:

```text
Ply needs your attention
Example delivery — report ready
Reported only; not yet controlled. Human QA is not attested.
Next: Start the planning launcher on your development machine for result control.
Run: example-delivery-001 · Gate: result-control/1
Opened: 2026-10-03T19:00:00Z
```

The only other payload fields are `mrkdwn=false`, `parse="none"`,
`unfurl_links=false`, and `unfurl_media=false`. Actor claims, source paths and report
bodies stay local. Local file paths are not mobile return links.

## Explicit retry and uncertainty

Only `rejected`, proven `not_sent`, or `rate_limited` after its recorded deadline can
retry. There is no background loop. Preview the next attempt separately:

```sh
ply workflow notification retry '<notification-id>' --route ./slack-route.json --check --format json
ply workflow notification retry '<notification-id>' --route ./slack-route.json --apply --confirm 'sha256:<new-retry-confirmation>' --format json
```

The original send digest does not authorize retry. An already recorded retry digest
reads back its result without another attempt. Acknowledged delivery and `unknown`
cannot retry. A crash after reservation, a lost response, an ambiguous server error,
or a lost final receipt leaves uncertainty. Preserve the local return and check Slack
separately. Do not change revision, ID, route or state root to evade an unknown gate.
Corrupt state must be investigated without blind resend or automated repair.

`show` preserves transport knowledge while reporting source freshness as `current`,
`changed` or `unavailable`. A stale acknowledged message does not become unsent.
`persistence` distinguishes a stored result, a reservation without a completed receipt,
and failure to create a new receipt. A nonzero exit never means an effect was undone.

| Exit | Meaning |
| --- | --- |
| 0 | Valid preview/readback, acknowledged transport, or exact duplicate |
| 2 | Invalid CLI/input/format or missing/invalid credential before reservation |
| 4 | Conflict, unknown ID, busy state, forbidden transition or early retry |
| 5 | Preserved attempt with rejected, not-sent, rate-limited or unknown delivery |
| 1 | Local I/O failure or inability to durably preserve the required result |

Preview/show and automated fixture tests do not establish live Slack delivery or human
product acceptance. Installation, real setup, sending, QA and integration remain the
owning procedure's separately authorized activities.
