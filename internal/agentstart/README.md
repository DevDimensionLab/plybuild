# Agent start

`StartAgent(ctx, StartInput) (Receipt, error)` starts one visible Herdr agent
(`claude` or `codex`) and sends it one prompt. `ResumeAgent(ctx, StartInput)`
continues a start that stopped before its prompt (see Same key). `ReadReceipt(path)` reads the
receipt without effect. `ResolveHerdr()` binds the `herdr` found on `PATH` by
resolved path and SHA-256. The package imports only the standard library,
`canonicaljson`, `atomicfile` and `filelock` (with `golang.org/x/sys/unix`); a
test checks this with `go list -deps`. It
knows no caller concepts. The only way for a caller to veto an effect is
`Precondition`, and `Key` is an opaque string.

## Sequence

reserve receipt → tab create → agent start → observe a ready agent → prompt.
Ready means settled with a native session, or, when Herdr reports no native
session, settled, `interactive_ready` and a reported `foreground_cwd`. The prompt
is addressed to the pane and needs no native ID: Herdr clears the agent name when
the agent is replaced, and the full identity below is checked on every
observation. An unreported session stays unknown; none is invented. Nothing is
ever sent again.

The tab response and every `agent get` must carry plain workspace, tab, pane and
terminal identities equal to the request and the created tab; a missing or
different value stops the start before further mutation. The agent name derived
from `Key` and provider kind must match. An observed native session must have
session kind `id` and the requested provider, and may neither change nor vanish
during startup. `foreground_cwd` must equal `Cwd` for Claude (always) and for
Codex when Herdr reports it; a missing Codex value is accepted only together
with a native session.

Each Herdr call checks the bound binary digest first, has a 40 s limit and a
1 MiB output limit, and must return a successful `result` envelope. Errors name
the subcommand but never echo arguments (they carry the prompt).

## Provider arguments

`agent start NAME --kind KIND --pane PANE --timeout 30000`, followed by
`-- FLAGS` only when a flag is set. `nil` omits a flag and inherits the
provider configuration.

| Field | Claude | Codex |
|---|---|---|
| `Model` | `--model` | `--model` |
| `Permission` | `--permission-mode` | `-c default_permissions="…"` |
| `CodexProfile` | rejected | `--profile` |

No trust or approval flag is ever sent. Both providers inherit the tab cwd;
`-C` is not sent. `Env` entries become `--env KEY=VALUE` on tab create; nothing
is added implicitly. The receipt states what was requested, not an observed
model, profile or sandbox.

## Input binding

`InputSHA256` is the SHA-256 of canonical JSON of `Provider`, `Model`,
`Permission`, `CodexProfile`, `Cwd`, `Workspace`, `Label`, `Env`, `Prompt`,
`Herdr.Path` and `Herdr.SHA256`. `Key`, `ReceiptPath` and `Precondition` are
not part of it. The receipt stores the digest, not the prompt.

## Receipt

`Reached` is the highest confirmed effect: `reserved`, `tab_created`,
`agent_started`, `session_bound`, `prompt_submitted`. It never decreases, and
`session_bound` and `Session` appear only when a native session was observed.
`Pending` names a mutating call (`tab_create`, `agent_start`, `prompt`) and is
written durably before that call and
cleared when its result is recorded. `Failure{Stage, Class, EffectUnknown}`
describes the step that stopped. Classes: `precondition`, `binary_changed`,
`tab_response`, `agent_response`, `not_ready`, `prompt_response`.
`EffectUnknown=false` is claimed only for the failed step when no mutating
call was made for it (precondition or read-only observation).

`Receipt.EffectUnknown()` is true when `Pending` is set, when the failure says
so, or when the start has neither failed nor reached `prompt_submitted`.
`reserved` alone is not "no effect". `EffectUnknown=false` concerns only the
step that stopped: it never means that no earlier effect happened (see
`Reached`, `Tab`, `Pane`) and never means that a restart or cleanup is safe.

## Same key

`StartAgent` runs the sequence at most once per key, in the process that wins
the exclusive create of the receipt. Every later call reads the receipt and
makes no Herdr call and no `Precondition` call:

| Stored receipt | Result |
|---|---|
| complete, same digest | receipt, no error |
| different digest or key | `input_mismatch`, stored receipt returned |
| no digest, `plan` and `state` (older receipts) | `legacy_unbound`, mapped receipt returned with `Legacy=true`; equality is not claimed |
| empty, unparseable or unrecognized | `incomplete_reservation` |
| `reserved` without marker, failure or pending | `incomplete_reservation` |
| anything else | `already_reserved`, receipt shows `Pending`/`Failure` |

`StartAgent` never retries or sends a second prompt. A reservation is never
deleted by this package. Continuing after a partial start is a human decision
that starts from inspecting the Herdr agent named after `Key`; when the start
stopped before its prompt (for example Claude waiting on folder trust), the
caller may then call `ResumeAgent` with the same input.

`ResumeAgent` never calls `tab create` or `agent start`. It holds an exclusive
lock on `ReceiptPath + ".lock"` for the whole call and decides on the receipt
read under that lock, so a waiting call sees the previous call's result:

| Stored receipt | Result |
|---|---|
| different digest or key, legacy, unreadable | as for `StartAgent`, no call |
| complete, same digest | receipt, no error, no call |
| `Failure` set, `Pending` empty or `agent_start`, `Reached` `tab_created`, `agent_started` or `session_bound`, and `Tab`, `Pane`, `Terminal` set | eligible for observation (below) |
| anything else, notably `Pending=prompt`, `Pending=tab_create`, `reserved` or no `Failure` (crash or owner still running) | `not_resumable`, no call |

Eligible only means the pane is observed again; it does not promise that an
agent exists. A start stopped by a precondition before `agent start` (for
example `tab_created` with `Pending` empty) holds no agent, and the
continuation is refused by the identity check. A continuation observes the
stored pane read-only with the same identity check as a start; a stored `Session` must be observed unchanged. Until the prompt
marker is written the receipt is not written at all: identity differences, a
missing agent, the deadline, cancellation and `Precondition` failures return
the typed error (`agent_response`, `not_ready`, `precondition`) with the stored
receipt byte-identical, so the original failure stays and a later
`ResumeAgent` may continue once the cause is settled. Cancellation is checked
again after `Precondition`, before the marker. Then `Precondition`, one
write with `Pending=prompt`, no `Failure`, and `Reached` at least
`agent_started` (`session_bound` and `Session` when observed), then the same
prompt call. From that write on the start rules apply: a lost reply leaves
`Pending=prompt` and `prompt_response`, which is never resumable. A failed lock
release is returned as well (`receipt_unavailable` wrapping
`*filelock.ReleaseError` when the call otherwise succeeded, joined to the
original error otherwise); the returned receipt is still the recorded one.

Errors are `*Error{Code, Stage, Err}`. Codes: `invalid_input` (nothing is
written), `input_mismatch`, `legacy_unbound`, `incomplete_reservation`,
`already_reserved`, `not_resumable`, `receipt_unavailable`, and the failure
classes.

## Writes

Only `ReceiptPath` is written (parent directory created with mode 0700; file
mode 0600): created exclusively, then replaced with `atomicfile.Replace`.
`ResumeAgent` also creates or opens `ReceiptPath + ".lock"` (regular file,
mode 0600). The lock file carries no state and may remain.
Tests inject a Herdr double and crash points; no provider is started.
