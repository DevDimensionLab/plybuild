# Agent Session: Narrow Unzip Output File Provenance

Status: NEXT
Session ID: `2026-08-26T124026+0200-narrow-unzip-output-file-provenance`
Created: `2026-08-26T12:40:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ce9587d3ed287014fd7012c23b620583b619ab5eb375c69d72fcfae79fddfc88`
Previous: [2026-08-26T121214+0200-narrow-browser-launcher-system-capability.md](2026-08-26T121214+0200-narrow-browser-launcher-system-capability.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 provenance move: return shell Unzip's existing private
destination-file open through the existing `filesystem.File` interface so the
following adapter `Copy` and `Close` calls no longer inherit concrete
`*os.File` provenance. Preserve exact open, archive, path, directory, entry,
copy, close, error, traversal, partial-result, and public behavior; P3.64; every
completed move; and zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.64 product commit `0b96f10` passed its
truthful complete checkpoint.

This mission authorizes only focused filesystem adapter and shell Unzip
recording contracts; one narrow helper that delegates to the complete existing
`FileSystem.OpenFile` operation while returning the existing `filesystem.File`
interface; and replacement of the one destination `filesystem.OpenFile` call
inside `unzipWithDependencies()` in `pkg/shell/command.go` with that helper.

It does not authorize changing the complete `FileSystem` interface, its system
implementation, `filesystem.OpenFile`, public `file.OpenFile`, or another
filesystem operation or caller; changing the Unzip command flow; Maven;
browser launcher; plugin diagrams; profile editor; shell Run or Git; process,
HTTP, clock, or server work; inventory; scanner; audit apparatus; mutation
harnesses; or P4-P8 implementation.

# Measurements At Start

Clean product commit `0b96f10` has 388 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 7/29 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 81 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`f817974080ab4b1641155e24e7d2b7f4d0fddf1ae7c06d7c4922320a70f5a842`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Unzip output-file provenance move and
its recording contracts, then make the normal separate continuity-only commit.
Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the existing complete
caller-owned Unzip filesystem dependency value. The new helper must make
exactly one call to that dependency's existing `OpenFile` operation with the
complete path, flags, and mode; return the same file as the existing
`filesystem.File` interface and the exact error; remain safe at zero; and add
no selection, fallback, retry, wrapping, logging, cleanup, or global state.

Preserve the exact source and destination bytes, archive open and deferred
close, zip-slip boundary and error, archive directory metadata detection,
joined output paths, recursive directories and modes, ordered filenames,
destination flags `os.O_WRONLY|os.O_CREATE|os.O_TRUNC`, each entry mode, entry
open, copy, output close, entry-reader close, all attempts and operation order,
ignored copy results, direct open and close errors, exact partial filenames,
error precedence, traversal stopping or continuation, safe zero behavior,
system selection, exported `Unzip` signature and return, callers, and every
public behavior.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Preserve complete caller-owned dependency fields in every
affected recording double. Leave completed filesystem, process, browser,
plugin-diagrams, profile, shell Run and Git, Maven, server, and clock behavior
unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `0b96f10`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete filesystem adapter code/tests and every complete
filesystem double/caller, complete shell Unzip code and all its contracts and
callers, the P3.64 browser-launcher contracts, P3.63 plugin-diagrams contracts,
P3.62 profile-editor contracts, P3.61 shell Git contracts, P3.60 shell Run
contracts, P3.58 Maven stdout contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Start red by strengthening the focused filesystem-open and shell Unzip
   recording contracts. Prove exact complete dependency delivery, path, flags,
   mode, returned `filesystem.File` identity, exact error, safe zero behavior,
   non-empty population, complete caller-owned dependency preservation, every
   existing Unzip operation and ordering guarantee, partial result, error
   precedence, traversal rule, and absence of another operation.
2. Add only the authorized interface-returning helper and replace only Unzip's
   authorized destination-file open call. Preserve the complete adapter
   interface and system operation, exported `filesystem.OpenFile`, public
   `file.OpenFile`, `unzipWithDependencies`, exported `Unzip`, every operation
   field and sequence, caller, and return path. Regenerate exact Q1.3 without
   changing the scanner or broadening the move to force a number.
3. Run focused filesystem/Unzip and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests,
   race, and vet; all four host acceptance flows; the 15-control audit meta-
   suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean
   audit; and empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The
   full audit may exit 1 for documented findings but never 2, and comparable
   ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Unzip output-
file provenance move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
