# Installation and legacy usage

Read only the section needed for an authorized installation/upgrade, a legacy task continuation,
or diagnosis of the helper's transport result. Normal Task@2 operation is in [the skill](../SKILL.md).

## One-time installation and credential bootstrap

Only during the separately authorized installation activity, from the owning source repository:

```sh
make install-agent-notify
python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/bootstrap.py
```

Bootstrap prefers `PLY_SLACK_WEBHOOK_URL`, otherwise reuses its named private value or requests
hidden terminal input. An invalid environment value is an error, with no file fallback. It creates
only the named route, private credential file, and state directory; it sends nothing. Conflicting
files, symlinks, unsafe ownership, and unsafe modes fail without repair.

- Route: `/Users/perottochristensen/.config/ply/notifications/ply-log/route.json`.
- Credential: `webhook` alongside it, mode `0600` in a `0700` directory.
- State: `/Users/perottochristensen/.local/state/ply/agent-notifications/ply-log`, mode `0700`.

The helper passes the credential only to send/check children, never to show. Normal calls reuse
this installation; they do not search secret stores, set up another app/route, or request another
webhook. Never put a webhook in prompts, argv, reports, or ordinary logs.

An explicit upgrade uses the reviewed previous skill source and a separate private backup path.
The installer checks the complete owned file set, exact previous bytes, ownership, and modes,
then preserves a complete backup before replacing owned files. Inspect the destination and backup
after an interrupted upgrade; do not blindly replace a mixed installation.

```sh
python3 skills/ply-agent-notify/scripts/install.py   --destination /absolute/installed/ply-agent-notify   --upgrade-from /absolute/reviewed-base/skills/ply-agent-notify   --backup /absolute/private/ply-agent-notify-before-upgrade
```

For fixtures, use `make install-agent-notify AGENT_NOTIFY_DEST=/absolute/private/fixture/ply-agent-notify`.
Bootstrap accepts explicit `--config-dir` and `--state-root`; notify accepts `--credential-file`
and a single `--ply-bin` executable. Keep fixture paths under the authorized test root.
`make test-agent-notify` builds the CLI and uses a network-free production-command driver. This
checks technical behavior; actual Slack delivery requires an authorized real notification.

## Continue an older task

New tasks use Task@2 and Event@2, producing request-v3. An existing Task@1 must keep Event@1 and
request-v2. Mixed pairs fail. Preserve existing task/event bytes and IDs; converting a record to
change its presentation does not create permission for another send.

Event `kind`, `activity`, and `run` may be omitted: the helper binds them from the existing Task.
Supply the actual report/start receipt with the existing flags. Complete sources retain their
original locators. When a source lacks notification metadata, the helper preserves its exact
bytes and original locator/digest beside a small `PlyAgentNotificationSource@1` binding in the
same private event directory. The native document's kind, fields, and meaning remain unchanged;
the binding only associates that source with the authorized Task. Explicit legacy or native
identity conflicts fail locally before credential lookup or any Ply call. An explicit worktree
must match the Task; Task@2 also checks an explicit provider or origin cwd. No separate handwritten
notification report is needed.

Task@1 contains `kind`, `activity`, `run`, `worktree`, `handoff`, `route`, `event_root`,
`actor_claim`, and `notification_authorized`. Task@2 adds `provider`, `origin_cwd`, `context_root`,
and `timezone`. Use the existing records and the schema validation in the helper rather than
rebuilding legacy records from a new template.

Physical origin must be inside the context root and match the helper's actual working directory.
A separately authorized native proxy preserves the provider's actual origin. Never relabel a
Codex task as Claude. Requires Python 3.9+ with IANA timezone data; unknown zones and `Local` fail.
A saved request retains its local event timestamp across changes to TZ or timezone data.

## Diagnose a helper result

The helper returns `product_outcome=unchanged`, transport `state`, command exits, and a receipt
when available. Exit 0 means acknowledged or suppressed; 2 means not attempted; 5 means rejected,
not sent, rate-limited, or unknown. The JSON state is authoritative for the report.

A crash after the durable apply marker permits only readback. `--port-revision` is a caller
annotation and never changes event identity or grants another attempt. Explicit Ply retries are
a separate authorized operation under the CLI's guards.

The helper's existing receipt includes `source_inputs` when it generated bindings, so the original
source path and digest remain reviewable. Reuse the original input paths and event ID for readback;
changing a source does not authorize a resend. A binding does not establish QA or transport success.

`--trigger optional_feedback`, `internal_correction`, and `waiting_for_feedback` return
`suppressed` without invoking Ply or writing notification state. Use the original question when
still waiting for feedback; do not construct a replacement ending event.
