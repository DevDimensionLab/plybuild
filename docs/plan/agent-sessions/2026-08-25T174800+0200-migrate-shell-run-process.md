# Agent Session: Migrate Shell Run Process

Status: NEXT
Session ID: `2026-08-25T174800+0200-migrate-shell-run-process`
Created: `2026-08-25T17:48:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `cb3efd125fdf1a2cbbb6375cc4adfc44b26587be684e540425c91d8f3340a6ff`
Previous: [2026-08-25T171500+0200-migrate-shell-unzip-file-open.md](2026-08-25T171500+0200-migrate-shell-unzip-file-open.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`exec.Command` process construction and synchronous execution flow behind
public `pkg/shell.Run` through the established `process.Execute` adapter
operation. Preserve the public signature, `Output` type and methods, exact
command name and argument order, stdout and stderr capture, synchronous
completion, established returned-output and process-error behavior, callers,
and every completed unzip, filesystem, Spring, process, browser, profile,
HTTP, tips, config, Maven, structurizr, Bitbucket, Wpost, local-config, and
supervisor move, with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes one private
complete shell Run composition containing `process.Dependencies`, one public
production selection of `process.System()`, focused recording of the existing
`process.Command` request, and replacement of the direct public-Run process
construction/execution flow with `process.Execute`. It authorizes removal or
reshaping of the now-private `run(*exec.Cmd)` helper only if it has no caller.

It does not authorize a process interface or adapter implementation change,
another process operation or caller, a `Git*` migration, structurizr or
plugin-diagram work, another shell function, any remaining unzip operation,
filesystem or HTTP work, Spring, clock/server adapters, public API, inventory,
mutation harnesses, or later-roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/shell` source/tests/callers,
complete `internal/adapter/process` source/tests and every complete process
double, Maven analyze and every public `shell.Run` caller, relevant command,
context, config, HTTP, filesystem, Spring, Maven, structurizr, Bitbucket,
Wpost, local-config, browser, profile, tips, file, template, and Kibana
code/tests, `.quality/inventory`, both design documents, and the focused audit
implementation. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `2a684a0` has 332 tests across 19 of 25 packages. Q0.6
has 25 guarded safe-writer sites, 20 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 6 of 25, Q1.2 is 0, Q1.3 is 12 violations of 33
production effect sites with clock and server absent, Q1.4 is 7 of 8, and
exact Q2.1 is 0 of 8 executable harnesses. The clean full audit exited 1 for
15 documented findings and never 2; comparable ratchets were five improved,
two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing or finalizing.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the shell Run process move and focused private
recording contracts, with measured planning notes. Then make the normal
separate handoff-only commit. Do not push, merge, publish, distribute, remove
the worktree, stash inherited changes, revert user work, or run destructive
Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve public
`Run(string, ...string) Output`, `Output`, `String`, and `FormatError`, every
caller, exact requested command name and variadic argument order, separate
stdout/stderr buffers and delivered bytes, one synchronous execution attempt,
the established debug-log position and command description, and the exact
legacy returned `Output.Err` and early-return behavior on dependency error.
Characterize that legacy error behavior before choosing assignments; do not
silently normalize or improve it.

Reuse the existing zero-safe process adapter and its complete
`process.Command`. Select `process.System()` only in the public production
entry. The private complete dependency composition must be passed whole. Do
not add fallback, retry, wrapping, extra logging, environment, working
directory, stdin, Start mode, shell parsing, cleanup, or error normalization.
The focused test double may record only the existing process request and may
write guarded in-memory stdout/stderr bytes; it must not launch a process,
touch the network, change the working directory, or write a file.

Do not change Git process composition, `shell.Unzip`, structurizr command
construction, Maven process composition, plugin diagrams, the process adapter,
another caller, public API, or `.quality/inventory`.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `pkg/shell` source
and tests, every `shell.Run` caller, complete `internal/adapter/process`
source/tests and every complete double, completed filesystem/unzip/Spring and
process contracts, and relevant caller packages before editing. Before the
full gate, read the complete launcher contract, Make meta-tests, P2A API/CLI
and subprocess contracts, and all four host acceptance flows.

# Three Moves

1. Start red with only focused private shell Run recording contracts. Prove
   complete production process selection, exact command name and ordered
   argument bytes, a single synchronous request with Start false and empty
   Dir/Stdin, exact stdout/stderr writer identity and delivered arbitrary
   bytes, exact dependency-error identity as observed through the legacy
   Output contract, safe-zero behavior, rejection of an empty recording
   population, and absence of unrelated process requests. Invoke no real
   process.

2. Add one private complete Run dependency composition and select
   `process.System()` only in public `Run`. Build the exact complete
   `process.Command`, preserve logging and buffer wiring order, and call the
   existing `process.Execute` helper once. Remove the direct `exec.Command`
   construction from this flow and change no adapter, caller, or other effect.

3. Run focused shell/process/Maven and relevant command, context, config,
   filesystem, HTTP, Spring, structurizr, profile, browser, tips, file,
   template, Bitbucket, Wpost, local-config, Kibana, and caller package tests;
   the launcher contract from `/bin/bash`; Make preflight meta-contracts;
   API/CLI and subprocess compatibility; full Go tests and race/vet; all four
   host acceptance flows; the audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/
   Q1.4/Q2.1/Q3.4 measurements; full clean checkpoint audit; and empty-HOME
   count-2. Expect nominal Q1.3 to improve from 12 of 33 to 11 of 32 when the
   one caller construction site leaves and the existing adapter implementation
   stays singular. Regenerate exact values; the full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent shell Run
process move or record an exact resumable state. Rewrite the rolling handover,
record the measured P3 result, answer this archive, create one linked NEXT
archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
