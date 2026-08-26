# Agent Session: Cover Context Package

Status: NEXT
Session ID: `2026-08-26T180111+0200-cover-context-package`
Created: `2026-08-26T18:01:11+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d15a436b770b306c5003da989c238153a0822b4b2a4b632e56265bca8ecc676a`
Previous: [2026-08-26T172816+0200-cover-logger-package.md](2026-08-26T172816+0200-cover-logger-package.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the sole remaining
untested package, `pkg/context`. Add deterministic characterization contracts
for the complete existing public Context operations and logger selection
without changing production code, filesystem/config/Maven behavior, project
ownership, logging, callers, or observable semantics. Preserve every completed
checkpoint and produce zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.8 product commit `67344a8` added complete
logger-package characterization, improved Q1.1 from 2/27 to 1/27, and held
exact Q1.3 at zero of 27.

This mission authorizes only focused tests in `pkg/context` that characterize
the existing `Context` state, `FindAndPopulateMavenProjects`,
`OnEachMavenProject`, `OnRootProject`, `LoadProfile`, `GetMavenRepository`,
package logger initialization, and `SetLogger`, plus the exact existing file,
config, Maven, project, and logger behavior those operations already invoke. It
does not authorize changing any production file, extracting a seam, changing
project or config ownership, changing walk or append order, validation,
logging, errors, writes, Maven repository selection, globals, callers,
scanners, audit apparatus, inventory, mutation harnesses, acceptance scripts,
public APIs, another package, manual evidence, or P5-P8 implementation. If
truthful tests expose a preferred isolation, ownership, validation, logging,
or repository design, preserve the existing behavior and record that product
decision for later scope.

# Measurements At Start

Clean product commit `67344a8` has 440 tests across 26 of 27 packages and
`pkg/logger` has 100% statement coverage. Q0.6 has 27 guarded safe-writer
sites, 22 write and 5 copy, zero skipped tests, and zero unsafe direct test
writes. Q1.1 is 1/27, Q1.2 is zero, exact Q1.3 is 0/27 with all five declared
adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases across
92 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`38e938f060312840cb244a575536aa613266b172de98542cb2bb2e0998da576b`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make the
normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the exact exported
Context fields and method signatures, private package logger and its
initialization, `SetLogger` assignment, value-versus-pointer behavior, project
append and traversal order, root-project selection, nil-job skipping, error
continuation and early-return choices, cloud-default merge, stealth and dry-run
effects, profile creation, Maven repository selection, exact log levels and
messages, and all existing file/config/Maven calls. Characterize representative
success, empty, error, skip, and partial-result paths that can be exercised
deterministically through the current interfaces. Do not add copying,
validation, recovery, synchronization, sorting, filtering, normalization,
fallbacks, output redirection, reset APIs, or error propagation.

Focused tests may use only `t.TempDir()`, guarded test writers, in-memory
models, logrus loggers/hooks/buffers, and isolated test-owned HOME/config/Maven
state. Every test that changes the private package logger or process state must
restore the exact original logger pointer/interface identity, environment,
working directory, and any other touched value. Do not use parallel tests
around globals or process state. Tests must open no socket, access no network,
launch no external program, perform no sleep or timed wait, write no repository
fixture, and leave no log hook, entry, file, environment value, or process-state
change behind. Reject empty table-driven, callback, project, job, and recorded-
entry populations before iterating. Leave logger, API handlers, templates,
resources, sorting, server, clock, Kibana, Spring, Maven production code,
config production code, filesystem production code, adapters, browser, plugin
diagrams, profile, shell, Unzip, and every completed contract unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `67344a8`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, both complete production files under
`pkg/context`, every context caller, and the exact config, project, Maven, file,
logger, profile, and guarded-writer code reached by each method. Read
representative exact-error, log-capture, global-restoration, temporary-home,
working-directory, pointer-identity, partial-result, callback-order,
nil-callback, dry-run, and explicit non-empty contracts, the import-aware
scanner, API/CLI contracts, and the T15 repair and baseline reproduction
README.

# Three Moves

1. Add one focused `pkg/context` test file proving the complete deterministic
   existing Context method branches and state transitions, exact logger
   selection and restoration, representative logs/errors, traversal and job
   order, value and pointer behavior, dry-run/write behavior, profile behavior,
   repository selection, partial results, and non-empty populations. Do not
   change production context behavior to make preferred isolation pass.
2. Run the focused package and command caller tests, inspect package coverage,
   and confirm that only the new context test file changed. Regenerate focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without changing scanner or inventory.
   Expect Q1.1 to improve from 1/27 to 0/27 and exact Q1.3 to hold at 0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent context coverage
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P4 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
