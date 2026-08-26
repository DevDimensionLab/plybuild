# Agent Session: Narrow Shell Run System Capability

Status: NEXT
Session ID: `2026-08-26T101755+0200-narrow-shell-run-system-capability`
Created: `2026-08-26T10:17:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3e150735be39600469a295d93fe7c9167c717b27270b4c741dc48e4314e1fba7`
Previous: [2026-08-26T094645+0200-migrate-shell-unzip-entry-open.md](2026-08-26T094645+0200-migrate-shell-unzip-entry-open.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 provenance move: give public shell `Run` a runner-only
system process dependency so its production composition no longer inherits the
unused standard-output capability added in P3.58. Preserve the exact system
runner, command and buffer identities, logs, errors, public behavior, P3.59,
every completed move, and zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.59 product commit `9f56714` passed its
truthful complete checkpoint.

This mission authorizes only focused process-adapter and shell Run recording
contracts, one narrow runner-only production selector on the existing process
adapter, mechanical preservation of complete process dependencies in affected
doubles, and replacement of only the `process.System()` selection inside
`systemRunDependencies()` in `pkg/shell/command.go`.

It does not authorize a process execution change, command or buffer change,
Git wrapper, Maven, profile editor, plugin diagrams, browser launching, Unzip,
another process caller, another adapter, inventory, scanner, audit apparatus,
clock/server, mutation harnesses, or P4-P8 implementation. Do not change
`process.System()` or Maven's completed standard-output behavior.

# Measurements At Start

Clean product commit `9f56714` has 384 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 15/37 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 76 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`41b51871f6436851613af04c86af93b6fb7f7f0e7f74833c47d8e3d45ddfb10e`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single shell Run selection move and its
recording contracts, then make the normal separate continuity-only commit. Do
not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve public
`Run(name string, args ...string) Output`, exact command name and ordered args,
debug log text and placement, distinct stdout and stderr buffer identities,
nil stdin and directory, synchronous execution, one exact runner request,
returned `Output`, exact dependency error, and safe zero dependency. Production
must use the exact existing system runner while exposing no standard-output
capability on this one dependency value. Do not add fallback, retry, wrapping,
logging, execution, cleanup, or global-state changes.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Preserve `process.System()` as the complete runner plus
exact `os.Stdout` dependency used by Maven. Pass complete dependencies without
silently discarding caller-owned fields in doubles.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `9f56714`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete process adapter code/tests and every complete
process double/caller, complete shell Run code/tests/callers, the P3.58 Maven
stdout contracts, the import-aware scanner, API/CLI contracts, and the T15
repair and baseline reproduction README.

# Three Moves

1. Start red with focused process-adapter and shell Run recording contracts.
   Prove the runner-only selector carries the exact system runner and nil
   stdout, `process.System()` still carries exact `os.Stdout`, the complete
   arbitrary command and buffer identities, exact error, safe zero behavior,
   complete dependency preservation, a non-empty request population, and
   absence of another process request.
2. Add only the narrow runner-only system selector to the existing process
   adapter and use it only in `systemRunDependencies()`. Preserve every command
   field and return path. Regenerate exact Q1.3 without changing the scanner or
   broadening the move to force a number.
3. Run focused process/shell and relevant caller tests; API/CLI and subprocess
   compatibility; launcher and Make contracts; complete tests, race, and vet;
   the 15-control audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4
   measurement; full clean audit; and empty-HOME count-2. Expect Q0.6 to hold at
   26 guarded sites. The full audit may exit 1 for documented findings but
   never 2, and comparable ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent shell Run system
capability move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create exactly
one reciprocally linked NEXT archive for the next coherent authorized roadmap
move, replace only the launcher's mutable regions, run launcher and handoff
contracts, and make the separate `docs: prepare next agent session` commit. Do
not launch a real successor, push, merge, publish, distribute, stash, revert,
or remove the worktree. COMPLETE remains invalid while P3 or P4-P8 is
unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
