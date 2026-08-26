# Agent Session: Narrow Maven Standard Output Provenance

Status: NEXT
Session ID: `2026-08-26T130843+0200-narrow-maven-standard-output-provenance`
Created: `2026-08-26T13:08:43+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3c880bbf28373de5dd123a7cd01ae045fdef9518c6520002ea67b66d1e569ce0`
Previous: [2026-08-26T124026+0200-narrow-unzip-output-file-provenance.md](2026-08-26T124026+0200-narrow-unzip-output-file-provenance.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 provenance move: preserve Maven's exact debug-only
system standard output while preventing public `RunOn` from inheriting concrete
`os.Stdout` provenance through the complete `process.System()` dependency.
Preserve exact command, log, directory, arguments, streams, error, callback,
dependency, and public behavior; P3.65; every completed move; and zero
comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.65 product commit `40cece5` passed its
truthful complete checkpoint.

This mission authorizes only focused process-standard-output and Maven command
recording contracts; one narrow process adapter helper that returns the exact
existing system standard output by delegating to the existing `Stdout(System(),
true)` composition; one private Maven production dependency selector carrying
the exact existing `SystemRunner()` and that exact standard output; and
replacement of the one `process.System()` selection in public `maven.RunOn`
with that selector.

It does not authorize changing `process.Dependencies`, `Runner`, `System`,
`SystemRunner`, `Stdout`, `Execute`, or system execution; changing Maven command
flow, logging, debug selection, callbacks, callers, arguments, or public API;
Unzip; browser launcher; plugin diagrams; profile editor; shell Run or Git;
filesystem, HTTP, clock, or server work; inventory; scanner; audit apparatus;
mutation harnesses; or P4-P8 implementation.

# Measurements At Start

Clean product commit `40cece5` has 388 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 5/27 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 83 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`1bae6abec78d2306c2205bdf1f96b7bffdabcef7e82a232a2aedafe536f3bf56`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Maven standard-output provenance move
and its recording contracts, then make the normal separate continuity-only
commit. Do not push, merge, publish, distribute, remove the worktree, stash
inherited changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the existing complete
caller-owned Maven process dependency value. The new helper must select the
same exact `os.Stdout` identity through the existing process composition,
without fallback, wrapping, buffering, copying, logging, cleanup, retry, or
global state. The private Maven selector must retain the exact system runner
and exact standard output while adding no other capability.

Preserve every arbitrary command and argument byte, argument order, callback
shape and invocation timing, repository value, exact project directory, info
log before execution, conditional stdout selection for debug and trace only,
nil stdin and stderr, synchronous execution, one exact attempt, direct process
error, safe zero injected dependency behavior, non-empty population, exported
`RunOn` signature and return, callers, and every public behavior.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Preserve complete caller-owned dependency fields in every
affected recording double. Leave all completed filesystem, process execution,
browser, plugin-diagrams, profile, shell, Maven graph, Unzip, server, and clock
behavior unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `40cece5`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete process adapter code/tests and every complete
process double/caller, complete Maven command code/contracts and callers, the
P3.65 Unzip output-file contracts, P3.64 browser-launcher contracts, P3.63
plugin-diagrams contracts, P3.62 profile-editor contracts, P3.61 shell Git
contracts, P3.60 shell Run contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Start red by strengthening the focused process-standard-output and Maven
   recording contracts. Prove exact system stdout identity, exact runner,
   complete dependency delivery and preservation, every log and command field,
   the full log-level matrix, one attempt, exact error, safe zero behavior,
   non-empty population, public composition, and absence of another operation.
2. Add only the authorized narrow system-standard-output helper and private
   Maven selector, then replace only public `RunOn`'s authorized
   `process.System()` selection. Preserve the complete process adapter API and
   system implementation, `runOn`, every callback and caller, and the return
   path. Regenerate exact Q1.3 without changing the scanner or broadening the
   move to force a number.
3. Run focused process/Maven and relevant caller tests; API/CLI and subprocess
   compatibility; launcher and Make contracts; complete tests, race, and vet;
   all four host acceptance flows; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The full audit
   may exit 1 for documented findings but never 2, and comparable ratchets must
   not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Maven standard-
output provenance move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
