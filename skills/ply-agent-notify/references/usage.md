# Installation and diagnosis

Read only the section needed for an authorized installation/upgrade or diagnosis of the
helper's transport result. Normal operation is in [the skill](../SKILL.md).

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

An installation that still has the retired `references/task-template.json` fails closed in this
installer: its owner applies a controlled patch instead. An explicit upgrade uses the reviewed previous skill source and a separate private backup path.
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

### Owner-selected app link

Only when separately authorized, and only after the installed `ply` supports request-v6, the owner
writes the link file beside the route. Nothing in the helper source changes:

```sh
umask 077
printf '{"app_link":"https://plyapp-pwa.web.app"}\n' > /Users/perottochristensen/.config/ply/notifications/ply-log/app-link.json
```

New short-form events then send v6 with an `Åpne plyapp` link. To change the link, rewrite the
file; to stop, remove it. Neither affects events that already have a request or attempt. An
invalid file makes new events fail as `invalid_app_link_owner_file` without a request. Once a
v6 event exists, keep the v6-capable `ply`; an older one reports that route's state as unknown.

## Event types and statuses

The short form's `--status` selects the event type:

| `event_type` | `public.status` |
| --- | --- |
| `agent_finished` | `done` |
| `agent_finished` | `ready_for_review` |
| `agent_finished` | `ready_for_your_check` |
| `agent_stopped` | `stopped` |
| `feedback_required` | `needs_answer` |

The earlier bound-Task entrance (Task and event files with handoff, start and report documents,
request v2/v3) is removed. Ply keeps its stored events as history and never sends them again. Do
not re-create such an event through the short form: its activity, run and event ID were already used.

Physical origin must be inside the context root and match the helper's actual working directory.
Never relabel a Codex task as Claude. Requires Python 3.9+ with IANA timezone data; unknown zones
and `Local` fail. A saved request retains its local event timestamp across changes to TZ or
timezone data.

## Diagnose a helper result

The helper returns `product_outcome=unchanged`, transport `state` and command exits. Exit 0 means acknowledged or suppressed; 2 means not attempted; 5 means rejected,
not sent, rate-limited, or unknown. The JSON state is authoritative for the report.

A crash after the durable apply marker permits only readback. `--port-revision` is a caller
annotation and never changes event identity or grants another attempt. Explicit Ply retries are
a separate authorized operation under the CLI's guards.

The helper writes no receipt: its JSON output and `ply workflow notification show` are the record.
Its `origin_lookup` is `resolved` when the request carries proven Herdr labels (v5), otherwise `none` (v4).

`--trigger optional_feedback`, `internal_correction`, and `waiting_for_feedback` return
`suppressed` without invoking Ply or writing notification state. Use the original question when
still waiting for feedback; do not construct a replacement ending event.
