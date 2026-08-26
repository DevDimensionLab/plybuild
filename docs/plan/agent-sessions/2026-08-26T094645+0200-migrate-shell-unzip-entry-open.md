# Agent Session: Migrate Shell Unzip Entry Open

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T094645+0200-migrate-shell-unzip-entry-open`
Created: `2026-08-26T09:46:45+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `01d203669b16b811b43fa302783b41bb66c7866f86ed30749884fa3f1a4d3876`
Previous: [2026-08-26T091301+0200-migrate-maven-debug-stdout.md](2026-08-26T091301+0200-migrate-maven-debug-stdout.md)
Next: [2026-08-26T101755+0200-narrow-shell-run-system-capability.md](2026-08-26T101755+0200-narrow-shell-run-system-capability.md)
Outcome: Product commit `9f56714` routes the exact archive entry through the
existing filesystem dependency. All focused and complete gates pass; the clean
audit exits 1 for the same 15 findings with zero comparable regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only shell Unzip's direct
archive-entry `zip.File.Open()` selection in `pkg/shell/command.go` through the
existing `internal/adapter/filesystem` dependency while preserving exact entry
identity, traversal, partial filenames, operation order, errors, ignored
results, and all completed behavior. Preserve P3.58, every completed move, and
zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.58 product commit `4376e05` passed its
truthful complete checkpoint.

This mission authorizes only focused filesystem-adapter and shell Unzip
recording contracts, one narrow archive-entry open operation on the existing
filesystem dependency, the production implementation that invokes the exact
injected `*zip.File`, mechanical preservation of the complete filesystem
dependency in its doubles/callers, and replacement of this one direct
selection:

```go
rc, err := f.Open()
```

It does not authorize another Unzip operation, traversal or path rules, file or
directory modes, logging, exported `shell.Unzip` signature, process/Maven,
plugin diagrams, another adapter, inventory, scanner, audit apparatus,
clock/server, mutation harnesses, or P4-P8 implementation.

# Measurements At Start

Clean product commit `4376e05` has 380 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 15/37 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 76 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`09cb1f364cbf55243f557d77c400b4befa7672d5b0a9378ad9e7d7cd6447ce32`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single shell Unzip entry-open move and its
recording contracts, then make the normal separate continuity-only commit. Do
not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the exact archive
source and destination, one archive-reader open, exact deferred archive-close
identity and ignored error, archive entry order, log text/order, zip-slip check,
partial-filename append point, directory handling, parent creation, destination
open flags and entry mode, and every early return.

For each non-directory entry, preserve destination open before entry open, one
entry-open attempt with the exact `*zip.File`, one copy with the exact opened
reader, ignored copy count/error, output close before entry-reader close, exact
close errors, and suppression of later operations after each existing failure.
The zero dependency must remain safe: no archive or entry open, no filesystem
mutation, and the exact existing `filesystem.ErrNoFilesystem` result. Pass the
complete dependency without fallback, retry, logging changes, extra calls,
error wrapping, cleanup changes, or global-state leakage.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Temporary archives outside the repository remain allowed
only through the guarded test helpers already used by the Unzip contracts.

Do not change path construction or validation, filenames, modes, copy/close
semantics, public API/CLI, process/Maven, `pkg/logger`, plugin diagrams, another
filesystem caller, inventory, scanner, baseline, audit repair, or any completed
product behavior.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `4376e05`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete shell Unzip code/tests/callers, the filesystem
adapter and every complete filesystem double/caller, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Start red with focused filesystem-adapter and shell Unzip recording
   contracts. Prove exact archive-entry identity, exact returned reader and
   error, safe zero dependency, the complete successful file-entry sequence,
   entry-open failure placement after destination open, absence of copy/close
   and later traversal after that failure, complete dependency preservation,
   rejection of an empty entry-open population, and absence of another
   filesystem request.
2. Add only the narrow archive-entry open operation to the existing filesystem
   dependency and production implementation, then replace only the direct
   `f.Open()` selection in shell Unzip. Preserve every existing operation and
   return path. Regenerate exact Q1.3 without changing the scanner or broadening
   the move to force a number.
3. Run focused filesystem/shell and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests,
   race, and vet; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The full audit
   may exit 1 for documented findings but never 2, and comparable ratchets must
   not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent shell Unzip
entry-open move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create exactly
one reciprocally linked NEXT archive for the next coherent authorized roadmap
move, replace only the launcher's mutable regions, run launcher and handoff
contracts, and make the separate `docs: prepare next agent session` commit. Do
not launch a real successor, push, merge, publish, distribute, stash, revert,
or remove the worktree. COMPLETE remains invalid while P3 or P4-P8 is
unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
