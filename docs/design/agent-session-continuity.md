# Agent Session Continuity Design

Status: accepted for the `codex/upgrade-quality` worktree.

## Problem

Long agent sessions accumulate tool output and stale assumptions. The previous
startup path split the next task between a tracked handover and an ignored
`.agent-task/current.md`; those copies had already drifted to different heads.
The upgrade needs short sessions whose task and learned state survive a fresh
Codex process.

## Authority

| Information | Authoritative location | Update rule |
| --- | --- | --- |
| Stable decisions and premises | `docs/design/*.md` | Change when evidence or an approved decision changes. |
| Roadmap, checkpoint status, and measured debt | `docs/plan/quality-upgrade.md` | Update after a measured quality move or checkpoint. |
| Resume state, learned facts, and next boundary | `docs/plan/quality-handover.md` | Rewrite at every session handoff; do not append a session diary. |
| One next-session mission | `codex-dev-start.sh` mutable prompt | Replace automatically before an agent session ends. |
| Active prompt mirror and previous missions | `docs/plan/agent-sessions/*.md` | Require the active byte-exact mirror at startup; preserve answered entries as linked history, never as independent task authority. |

The complete tracked continuity set - launcher, design notes, upgrade plan,
rolling handover, and archive graph - is sufficient to resume. The launcher's
mutable prompt remains the sole next-task authority; the active archive is its
required integrity mirror. The ignored `.agent-task/current.md` is not an input
to this workflow. The machine-readable authorized-checkpoint block in the
upgrade plan constrains lifecycle state: queued work cannot coexist with a
`COMPLETE` launcher.

## Launcher Contract

`codex-dev-start.sh` has a stable executable section followed by
`CODEX_STABLE_EXECUTION_END` and two mutable data regions. Stable execution
ends with `exit 70` before either mutable region. Every mutable data line starts
with `#|`, so the tail is comment data rather than executable shell.

The header has exactly four ordered lifecycle records:

```text
SESSION_STATUS
SESSION_ID
SESSION_ARCHIVE_REL
PREVIOUS_SESSION_ARCHIVE_REL
```

The prompt region stores one decoded prompt line per `#|` record. The launcher
snapshots its own file once and reads that snapshot with `awk`; it does not
`source`, `eval`, or execute mutable bytes. The contract test removes both
regions, hashes the remaining skeleton, and rejects raw executable lines in
either region.

The launcher:

- derives the repository root from its own physical path;
- requires its own source to be a regular, non-symlink file;
- requires the attached `codex/upgrade-quality` branch and expected plan files;
- accepts `--check`, `--print-prompt`, and a normal supervised start;
- allows dirty inspection through `--check` and `--print-prompt`, emits one
  static warning, and never inserts Git filenames or status output into the
  prompt; normal supervised start requires a clean worktree;
- fails startup when Git status cannot be measured;
- resolves `CODEX_BIN` to an external executable file, requires an external
  `python3` for structured-event validation, and inherits the user's Codex
  profile except for explicit normal-service and workspace-write overrides;
- never resumes an old session or invokes publication commands.

The archive block and decoded launcher prompt are byte-equal, including one
terminal LF. `--print-prompt` and the Codex argument use those same bytes. Dirty
state does not rewrite the prompt.

## Non-Interactive Supervisor Protocol

Normal invocation snapshots the validated HEAD, session ID, active archive,
archive count, and normalized stable-skeleton digest, then starts exactly one
fresh process with these arguments:

```text
codex exec --ephemeral -c service_tier="default" --sandbox workspace-write -C <repo> --json --output-last-message <external-path> <exact-prompt>
```

It never uses `resume` or hands task authority to agent output. Before the first
turn it creates one unique physical log directory beneath `TMPDIR` (or `/tmp`)
and rejects a log root inside the worktree. Each fresh turn has its own session-
named directory containing the byte-exact `events.jsonl`, Codex stderr, the
requested final-message path, and one isolated `scratch` directory. It exports
that path as both `TMPDIR` and `CODEX_SESSION_SCRATCH_ROOT`; tasks put all
disposable caches, archives, build contexts, reports, and evidence beneath it
instead of creating independent `/private/tmp/ply-*` roots. The launcher prints
the log root once and emits
short progress records for thread start, commands, file changes, agent messages,
and the terminal turn event. Event strings are printed only as `%s` data by the
parser; no event field is sourced, evaluated, or executed.

The embedded parser writes each raw input line before decoding it as JSON. A
successful stream has a non-empty population, exactly one `thread.started`, one
`turn.started`, and one final `turn.completed`, with no event after that
terminal. Invalid UTF-8/JSON, an unterminated final line, an empty stream,
duplicate or contradictory terminals, `turn.failed`, or `error` makes the turn
fail. The Codex process must independently exit zero. Scratch is deleted after
every turn, including failed or interrupted turns, after restoring owner
permissions on read-only directories. Successful completion also deletes the
supervisor log root; failed and partial logs remain for diagnosis. `--ephemeral`
prevents Codex from persisting an additional session rollout outside this
lifecycle.

Only after both process and stream succeed does the original parent re-snapshot
the on-disk launcher and validate the full repository and archive contract. The
normalized stable skeleton must equal the pre-turn digest, the worktree must be
clean and attached, and HEAD must have changed. Continuation additionally
requires a changed session ID, the former archive answered with a reciprocal
link, and exactly one newly committed archive that is the sole valid `NEXT`
tail. A `COMPLETE` handoff instead retains and answers the terminal session,
adds no archive, and succeeds only when the authorized queue has no active or
queued checkpoint. Any other state stops before another process starts.

`HUP`, `INT`, or `TERM` is forwarded to the active Codex process and stops the
event parser and parent with a non-zero status. The parent waits for the child,
does not start another turn, retains the partial logs, and does not alter task
authority.

## Archive Graph

Archive names and session IDs use
`YYYY-MM-DDTHHMMSS+ZZZZ-<lowercase-checkpoint>`. A relative path is valid only
when it is exactly
`docs/plan/agent-sessions/<session-id>.md`; nested paths, `..`, whitespace,
absolute paths, symlinked directories, and symlinked files are rejected.

Every archive has one value for `Status`, `Session ID`, `Created`, `Source`,
`Prompt SHA-256`, `Previous`, `Next`, and `Outcome` before its prompt block.
Runtime validation checks:

1. IDs match filenames, timestamp fields have valid calendar/timezone ranges,
   and `Created` equals the timestamp encoded in the ID.
2. `Source` is `codex-dev-start.sh`.
3. A `NEXT` archive has `Outcome: pending` and `Next: none`.
4. An answered archive has a non-pending outcome.
5. Previous and next links are local basename links and are reciprocal.
6. Every historical prompt block matches its stored SHA-256 digest.
7. Walking backward from the active archive visits every archive exactly once.
8. The active tail has `Next: none`; `NEXT` state has exactly one `NEXT`
   archive and `COMPLETE` has none.

The active archive's prompt block must match the launcher's decoded prompt
before Codex can start.

## Automatic Handoff Protocol

The user has authorized the ordered quality roadmap through P8. At a natural
boundary, after three quality moves, or before context quality declines, the
agent must prepare the next session before stopping. No separate trigger is
required.

Session finalization requires the agent to:

1. Re-measure repository and quality state.
2. Preserve corrected claims, expensive findings, failed approaches, unrun
   checks, and unfinished work in their owning plan documents.
3. Rewrite the rolling handover and select one next mission.
4. Answer the current archive, create one linked `NEXT` archive, and replace
   only the launcher's two mutable tail regions.
5. Run the launcher contract and applicable quality gates.
6. Create one local commit named `docs: prepare next agent session`.
7. Report the commit and next command, then stop without starting Codex.

The handoff commit may include only the launcher, affected design/plan files,
the answered archive, and the new archive. Implementation, test, inherited, or
unrelated dirty files are not auto-staged. Automatic handoff does not authorize
push, merge, release, stash, revert, or worktree removal.

Normal task execution may correct plans and the rolling handover. Launcher
session data and archives change only during session bootstrap or finalization,
after a coherent implementation move is committed or an exact resumable state
is recorded.

If no authorized mission remains after P8, finalization answers the tail
archive, leaves its `Next` value as `none`, and changes the launcher header to
`COMPLETE`. The launcher rejects `COMPLETE` while the authorized queue contains
an active or queued checkpoint. `--check` still validates the graph; normal
start and `--print-prompt` fail instead of replaying the answered task.

## Premises

| Premise | Evidence | Recheck | Consequence if false |
| --- | --- | --- | --- |
| A fresh non-interactive Codex turn accepts explicit ephemeral, working-directory, sandbox, JSONL, final-message, service-tier, and prompt values. | Measured 2026-09-06: the [official non-interactive-mode documentation](https://learn.chatgpt.com/docs/non-interactive-mode) defines fresh `codex exec`, `--ephemeral`, and JSONL turn events. | Recording `CODEX_BIN` asserts all twelve argument bytes, the isolated scratch environment, and characterized success/failure streams without invoking real Codex. | Stop before launch and repair the invocation or lifecycle contract. |
| Mutable task data can remain non-executable inside one script. | Contract test injects raw commands into each tail region, invokes the script, and observes that neither sentinel is created. | `/bin/bash test/codex_dev_start_test.sh`. | Stop startup and move data behind a verified execution boundary. |
| One linked archive graph detects duplicate or lost handoffs. | Contract test exercises a second generation, duplicate `NEXT`, disconnected history, a cycle, missing predecessor, broken backlink, historical prompt drift, and `COMPLETE`. | `/bin/bash test/codex_dev_start_test.sh`. | Reject startup until the graph is repaired from Git history. |
| Dirty recovery requires facts, not injected filenames. | Dirty fixture records the exact Codex argument and static stderr warning. | Compare clean, dirty, archived, and recorded prompt bytes. | Fail startup rather than insert raw status. |
| Restart commits can exclude unrelated work. | Git supports explicit path staging while other changes remain unstaged. | Inspect staged names before every restart commit. | Do not commit; leave a decision-ready handover. |
| One startup can avoid a globally selected Fast mode. | Measured 2026-08-24: local `codex-cli 0.149.0` accepts `-c key=value`, and `codex --strict-config -c 'service_tier="default"' --help` exits 0; the [official service-tier reference](https://developers.openai.com/api/reference/cli/resources/responses/methods/create) defines `default` as standard pricing/performance and `fast`/`priority` as Fast mode. | Recording `CODEX_BIN` asserts the exact normal-service argument. | Stop startup and update the pinned invocation contract. |

## Portability

`make test-agent-start` invokes `/bin/bash`. That is Bash 3.2 on the supported
macOS host and Bash supplied by the Linux image. The 62-control suite never
invokes real Codex: its recording executable captures argv and prompt bytes,
emits JSONL, performs synthetic committed handoffs, and waits for a forwarded
signal. It rejects empty event and turn populations, event-text execution,
malformed/truncated/failed streams, no progress, dirty or invalid handoffs,
contract drift, replay, and premature completion. The test chooses
`sha256sum` first, then macOS `shasum -a 256`, then `openssl`. Git fixtures
disable system/global config, templates, signing, and hooks. The checked-in
launcher and archive graph are copied into a clean synthetic repository before
they are executed, but only after the original graph's files and directory
components pass non-symlink checks. A nested second-generation run starts from
a source archive with no `.git`, covering Docker `COPY` and release-source
execution without a host worktree pointer. Nested mode is an internal argument;
ambient state cannot reduce the asserted control count. The test clears
`CDPATH`, system/global/command Git configuration, and Git's repository-local
environment variables. It enumerates visible and hidden Markdown archives
before copying, matching runtime `dotglob` behavior.

Docker daemon execution remains a separate environment measurement. The
launcher test itself has no network dependency.

## Invariants

- One worktree has one connected session archive graph.
- A `NEXT` launcher has exactly one active `NEXT` tail.
- The mission is the first decoded prompt line and names one primary outcome.
- The launcher never evaluates mutable task bytes.
- The supervisor never evaluates JSONL event bytes or trusts final prose as
  lifecycle evidence.
- One successful process may authorize at most one fresh successor, and only
  through a clean committed reciprocal handoff.
- Failed and partial turn logs remain outside the worktree; successful logs and
  every turn's scratch tree are removed automatically.
- Raw worktree data never becomes prompt input.
- Every session leaves one accurate `NEXT` tail while authorized work remains.
- `COMPLETE` is valid only when all authorized checkpoints are complete.
- Restart history supplements Git history; it does not replace measured state.
