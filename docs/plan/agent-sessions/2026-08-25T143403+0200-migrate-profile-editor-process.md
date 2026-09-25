# Agent Session: Migrate Profile Editor Process

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T143403+0200-migrate-profile-editor-process`
Created: `2026-08-25T14:34:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `44e57e5b74b6f13dab3164312d6d4f618beb0a878fd614a079275940b40057f8`
Previous: [2026-08-25T135854+0200-migrate-tips-list-read-dir.md](2026-08-25T135854+0200-migrate-tips-list-read-dir.md)
Next: [2026-08-25T151736+0200-migrate-open-browser-process-start.md](2026-08-25T151736+0200-migrate-open-browser-process-start.md)
Outcome: completed by `0e10282`; exact process stdin and the profile editor boundary are green at Q1.3 19/38, and the next focused P3 process-start move is linked.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`exec.Command(editor, ctx.LocalConfig.FilePath())` construction and `cmd.Run()`
in the `profileCmd` edit branch through the existing process adapter. Preserve
exact EDITOR lookup and fallback, local-config path selection, executable and
argument values, stdin/stdout/stderr and working-directory behavior, one run
attempt, exact error behavior, every later profile branch and follow-up, every
command object and registration, every completed process, filesystem, HTTP,
tips, config, Maven, structurizr, Bitbucket, Wpost, local-config, and supervisor
move, and every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`exec.Command` and `cmd.Run` operations inside `if configOpts.Edit` in
`cmd/profile.go`; one `Stdin io.Reader` addition to the existing internal
`process.Command` value and its exact direct system mapping to `exec.Cmd.Stdin`;
focused process-stdin contracts; and focused private profile-editor recording
contracts. It does not authorize another process operation or caller, an
environment adapter, a function-valued dependency, another profile/config
branch or effect, a public API, inventory, mutation harnesses, or later roadmap
implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `cmd/profile.go`, its registration, every
profile/config command caller and relevant command, context, local-config, and
terminal-width test, complete `internal/adapter/process` source and tests, every
`process.Command`, `process.Dependencies`, `process.Runner`, and `process.Execute`
caller and recording double, the completed Maven process flow, complete tips
and filesystem changes from move 43, relevant config/file/template/Maven/
structurizr/Bitbucket/Wpost code and tests, `.quality/inventory`, the continuity
and quality-lift designs, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation. Regenerate ignored reports outside the measured tree or remove
them before a clean audit.

Implementation commit `85f4c2b` has 304 tests across 18 of 25 packages. Q0.6
has 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 7 of 25, Q1.2 is 0, Q1.3 is 21 violations of 40 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for exact process stdin representation and the profile
editor boundary and contracts, with measured planning notes. Then perform the
normal separate handoff-only commit. Do not push, merge, publish, distribute,
remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, `profileCmd`, `configCmd`, all flags and registrations, public
API/CLI behavior, all four host acceptance flows, `ConfigOpts`, global context,
local terminal width behavior, config types, Maven process behavior, tips
behavior, and every caller's observable behavior.

Add only `Stdin io.Reader` to the existing complete `process.Command` request
and map it directly to `cmd.Stdin` before the existing `cmd.Run`. Keep name,
args, directory, stdout, stderr, safe zero behavior, `Runner`, `Dependencies`,
`Execute`, `System`, and every current caller unchanged. Extend the existing
adapter recording and system contracts to prove exact stdin identity and
delivered input. Do not add a new interface method, reconstruct an incomplete
dependency, invoke a shell, split the process operation, or change error
selection.

Preserve the current edit order: read `os.Getenv("EDITOR")` once, replace only
an empty value with exact `vim`, evaluate `ctx.LocalConfig.FilePath()` once,
then select `process.System()` only for a private complete profile-editor
composition and make one adapter execution attempt. The request must contain
the exact editor as `Name`, one exact path argument, exact `os.Stdin`, exact
`os.Stdout`, empty `Dir`, and nil `Stderr`. Return the exact process error.
On error, still stop before sync, reset, and print; on success, preserve every
existing later branch, condition, order, and result.

Do not parse or split EDITOR, use a shell, add stderr forwarding, change the
fallback, clean the path, preflight the executable or file, retry, wrap an
error, reorder environment/path selection, change the historically broad final
print predicate, or broaden into profile switching, sync, reset, print,
terminal config, filesystem, HTTP, tips, Maven, structurizr, Wpost, clock,
server, P4, P5, dependency, Docker, distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete `cmd/profile.go`
and relevant command, root, context, config, and local-config source/tests,
every profile command registration/caller, complete `internal/adapter/process`
source/tests and every complete process caller and double, the completed Maven
process contracts, complete move-43 tips/filesystem source and focused tests,
relevant file/template/Maven/structurizr/Bitbucket/Wpost code and tests, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing. Before
the full gate, read the complete launcher contract, relevant Make meta-tests,
P2A API/CLI/subprocess contracts, and all four host acceptance flows.

# Three Moves

1. Start red by extending focused process adapter contracts to prove the exact
   `io.Reader` identity reaches a recording runner and direct system execution
   receives exact input while preserving name, args, directory, stdout, stderr,
   one attempt, exact errors, and safe zero behavior. Add focused profile-editor
   contracts that prove complete production dependency selection, exact
   non-empty and empty EDITOR selection with `vim` fallback, exact arbitrary
   local-config path as the only argument, exact stdin/stdout identities, empty
   directory, nil stderr, one attempt, exact error identity, non-empty recorded
   populations, safe zero behavior, and no unrelated command or real process.

2. Add only `Stdin io.Reader` to `process.Command` and its direct system mapping.
   Keep `profileCmd` as the production entry, preserve EDITOR-before-path
   evaluation, select `process.System()` only for a private complete editor
   composition, and replace only the direct command construction/run. Do not
   change another process call, profile branch, command effect, adapter
   operation, caller, inventory, or public API.

3. Run focused command/process/config/context and relevant Maven, tips,
   filesystem, file, template, structurizr, Bitbucket, HTTP, local-config, and
   caller package tests; the launcher contract from `/bin/bash`; Make preflight
   meta-contracts; API/CLI and subprocess compatibility; full Go tests and
   race/vet; all four host acceptance flows; the audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements; full clean checkpoint
   audit; and empty-HOME count-2. Expect Q1.1 to hold at 7 of 25 and nominal
   Q1.3 to improve from 21 of 40 to 19 of 38 when the two direct process sites
   leave. Q0.6, Q1.2, Q1.4, and exact Q2.1 must hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent process-stdin
and profile-editor move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
