# Agent Session: Migrate Maven Debug Stdout

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T091301+0200-migrate-maven-debug-stdout`
Created: `2026-08-26T09:13:01+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `58e0daf53a3bef445f0f1fff667e8c2223f4ceb47c39f44526366e524e0fd13d`
Previous: [2026-08-26T084000+0200-migrate-plugin-diagrams-dot.md](2026-08-26T084000+0200-migrate-plugin-diagrams-dot.md)
Next: [2026-08-26T094645+0200-migrate-shell-unzip-entry-open.md](2026-08-26T094645+0200-migrate-shell-unzip-entry-open.md)
Outcome: product commit `4376e05` routes only Maven's conditional stdout selection through the complete process dependency; 380 tests, exact Q1.3 15/37, the clean gate, T15, full audit, and empty-HOME count-2 pass with zero comparable regressions

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only Maven command's
conditional `logger.StdOut()` capability selection in `pkg/maven/command.go`
through the existing `internal/adapter/process` dependency while preserving the
exact Maven command, conditional stdout identity, returned process error, and
all completed behavior. Preserve P3.57, every completed move, and zero
comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.57 product commit `9164e06` passed its
truthful complete checkpoint.

This mission authorizes only focused process-adapter and Maven recording
contracts, one narrow standard-output capability on the existing process
dependency, the production selection needed to supply `os.Stdout` from
`process.System()`, the conditional forwarding helper, mechanical preservation
of that complete dependency in relevant process doubles/callers, and replacement
of this one caller-side selection:

```go
Stdout: logger.StdOut(),
```

It does not authorize another Maven flow, command arguments, logging, the
exported `maven.RunOn` signature or returned callback, `pkg/logger` behavior or
public API, plugin diagrams, another process request, another adapter, shell
unzip provenance, inventory, scanner, audit apparatus, clock/server, mutation
harnesses, or P4-P8 implementation.

# Measurements At Start

Clean product commit `9164e06` has 375 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, Q1.3 is 7/29 with clock
and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases
across 75 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`9f584097422a17fd6eea005dda97182e042aad05a34837f0461ce0933248f700`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Maven stdout-selection move and its
recording contracts, then make the normal separate continuity-only commit. Do
not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve exact executable and
ordered argument bytes, `project.Path` as `Dir`, nil stdin and stderr,
`Start: false`, one synchronous process attempt, and the exact dependency error.
When the global logrus debug level is enabled, preserve exact `os.Stdout`
identity as the command stdout. When it is disabled, preserve nil stdout. Keep
the log call before process execution and preserve every level/state behavior.

Select the process runner and standard-output capability together only in
`process.System()`. The narrow forwarding helper must return the injected writer
only when enabled and nil otherwise. Pass the complete dependency value without
fallback, retry, logging changes, another process attempt, error wrapping,
environment changes, cleanup, or global-state leakage. The zero dependency must
remain safe: no tool launch and nil output. Focused tests must launch no real
Maven or another program, touch no network, and write no repository fixture.

Do not change the command's name, arguments, directory, error behavior,
`Repository` value behavior, exported callback shape, logging text/order,
`pkg/logger`, plugin diagrams, shell unzip, another caller, public API/CLI,
inventory, scanner, baseline, audit repair, or any completed product behavior.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `9164e06`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete Maven and logger code/tests/callers, the process
adapter and every complete process double/caller, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Start red with focused process-adapter and Maven recording contracts. Prove
   exact system `os.Stdout` selection, enabled injected-writer identity, disabled
   and zero-dependency nil output, the complete Maven command in both level
   states, one request, exact error, complete dependency preservation, rejection
   of an empty recorded population, and absence of another process request. Use
   recording boundaries so focused tests invoke no external program.
2. Add only the narrow process standard-output dependency capability and
   conditional forwarding helper, then route the one Maven `logger.StdOut()`
   selection through it. Preserve every existing command field and return path.
   Leave `pkg/logger` untouched. Regenerate exact Q1.3 without changing the
   scanner or broadening the move to force a number.
3. Run focused Maven/process/logger and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests, race,
   and vet; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The full audit
   may exit 1 for documented findings but never 2, and comparable ratchets must
   not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Maven stdout
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make the
separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
