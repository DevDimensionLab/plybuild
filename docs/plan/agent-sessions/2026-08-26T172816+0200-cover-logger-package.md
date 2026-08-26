# Agent Session: Cover Logger Package

Status: NEXT
Session ID: `2026-08-26T172816+0200-cover-logger-package`
Created: `2026-08-26T17:28:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `551a599d5cb7993bd7e215a9b0bdf985598d29e9ecac7305f49ede79ff4dabc6`
Previous: [2026-08-26T165809+0200-cover-webservice-api-package.md](2026-08-26T165809+0200-cover-webservice-api-package.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the untested package
`pkg/logger`. Add deterministic characterization contracts for the complete
existing Collector and public logger-helper behavior without changing
production code, globals, hook installation, formatting, levels, output
selection, callers, logging bytes, or observable semantics. Preserve every
completed checkpoint and produce zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.7 product commit `4da64d2` added complete
webservice-API handler characterization, improved Q1.1 from 3/27 to 2/27, and
held exact Q1.3 at zero of 27.

This mission authorizes only focused tests in `pkg/logger` that characterize
`Collector.Levels`, `Collector.Fire`, package initialization of the collector
hook, `DebugLogger`, `Context`, `ExternalError`, `SetJsonLogging`,
`SetFieldLogger`, `IsFieldLogger`, `StdOut`, and `LogEntries`, plus the exact
package and logrus globals they already use. It does not authorize changing any
production file, extracting a seam, changing hook or entry ownership, copying
entries, adding public reset behavior, changing log levels, fields, formatters,
output selection, error text, config, callers, scanners, audit apparatus,
inventory, mutation harnesses, acceptance scripts, public APIs, another
untested package, manual evidence, or P5-P8 implementation. If truthful tests
expose a preferred logger ownership, isolation, formatting, or output design,
preserve the existing behavior and record that product decision for later
scope.

# Measurements At Start

Clean product commit `4da64d2` has 434 tests across 25 of 27 packages. Q0.6 has
27 guarded safe-writer sites, 22 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 2/27, Q1.2 is zero, exact Q1.3 is 0/27 with
all five declared adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is
zero phrases across 91 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`0682df76893e0eac05302870ea986b1f18cab32e7513a7e0e91ab85da5f16c38`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make the
normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the exact exported
logger signatures, package globals and identities, private logger and collector
hook arrangement, and standard-logrus state split. Characterize the exact
ordered `Info` and `Warn` hook levels, append-only entry-pointer capture, and
the collector hook installed on the package logger. Characterize
`DebugLogger`'s caller-derived fields, empty data, and private debug-level side
effect; `Context`'s empty data; representative exact `ExternalError` text;
exact JSON formatter selection; field-logger state transitions; `StdOut`'s
existing dependence on the standard logrus logger's debug enablement and exact
`os.Stdout` or nil result; and the exact recorded entries exposed by
`LogEntries`. Do not add copying, validation, recovery, synchronization,
formatting changes, output redirection, reset APIs, or error propagation.

Focused tests may use only in-memory logrus entries, loggers, hooks, and
buffers. Every test that mutates package or standard-logrus globals must restore
the original `fieldLogger`, `collector`, `log`, level, formatter, output,
hooks, and any other touched state, including exact pointer identities. They
must open no socket, launch no external program, perform no sleep or timed wait,
write no fixture, and leave no hook or entry behind. Do not use parallel tests
around globals. Reject empty table-driven and recorded-entry populations before
iterating. Leave context, API handlers, templates, resources, sorting, server,
clock, Kibana, Spring, Maven, config, filesystem, process, HTTP adapters,
browser, plugin diagrams, profile, shell, Unzip, and every completed contract
unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `4da64d2`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, both complete production files under
`pkg/logger`, package initialization, every logger caller, relevant logrus
types and package-level state behavior, representative in-memory log capture,
exact-error, global-restoration, stdout-selection, pointer-identity, and
explicit non-empty contracts, the import-aware scanner, API/CLI contracts, and
the T15 repair and baseline reproduction README.

# Three Moves

1. Add one focused `pkg/logger` test file proving exact Collector levels and
   entry capture, installed hook behavior, logger-helper fields and effects,
   exact error text, formatter and field state, stdout selection, exposed log
   entries, restored globals, and non-empty populations. Do not change
   production logger behavior to make preferred ownership or isolation pass.
2. Run the focused package and relevant command and package caller tests,
   inspect package coverage, and confirm that only the new logger test file
   changed. Regenerate focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without
   changing scanner or inventory. Expect Q1.1 to improve from 2/27 to 1/27 and
   exact Q1.3 to hold at 0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent logger coverage
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P4 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
