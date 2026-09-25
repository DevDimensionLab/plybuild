# Agent Session: Migrate Maven Process Flow

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T130235+0200-migrate-maven-process-flow`
Created: `2026-08-24T13:02:35+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `accbeba9c63b46dfe7d50a59e6707c53dfacdcdb8b5fa102831ba337f6d2571c`
Previous: [2026-08-24T122234+0200-introduce-process-adapter.md](2026-08-24T122234+0200-introduce-process-adapter.md)
Next: [2026-08-24T140112+0200-cover-cloud-clone-seam.md](2026-08-24T140112+0200-cover-cloud-clone-seam.md)
Outcome: P3 move 3 completed at f5ee37d: Maven subprocess execution now uses the process adapter, Q1.2 remains zero, Q1.3 is 76 of 78, the Maven swap raises Q1.4 to 3 of 8, and the clean gate passed with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the third P3 move in one measured change: reuse
`internal/adapter/process`, migrate the coherent Maven subprocess flow in
`pkg/maven/command.go`, and kill the declared `maven-process` argument swap.
Reduce Q1.3 by the migrated Maven sites while keeping Q1.2 at zero and
preserving all P2A contracts with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, the process adapter and tests, `pkg/maven/command.go`, all
`maven.RunOn` callers, relevant Maven tests, the Q1.2/Q1.3/Q1.4 scanners, and
the API/CLI/subprocess and host acceptance contracts. Regenerate ignored
reports. Implementation commit `03d6242` has 47 tests across 11 of 23 packages,
Q1.2 at 0, Q1.3 at 77 violations of 79 production effect sites with four
adapter paths absent, and Q1.4 at 2 of 8. The exact Q2.1 validator remains 0 of
8 executable harnesses; the upstream filename-only denominator sees the
non-executable Git seam driver as one `mutate-*` path. The clean gate passed,
the full audit exited 1 for 16 documented findings and never 2, and comparable
ratchets were five improved, two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Maven process-adapter move. Do not push, merge,
publish, remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands. Do not invoke a publisher or distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`maven.RunOn`, its returned callback signature, error behavior, executable and
argument order, project working directory, and logger stdout wiring. Production
must still execute the same Maven/dot/ktlint argv in the same directories.
Recording doubles must have safe defaults, capture complete commands, and fail
when the asserted call population is empty.

Migrate only the two direct process effects in `pkg/maven/command.go`. Do not
migrate cloud, generic shell, Git, diagrams, profile/editor, browser-opening,
HTTP, filesystem, clock, or server calls in the same move. Do not enter Docker,
dependencies, cloud, Spring, distribution, or formal mutation-harness scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`internal/adapter/process/process.go`, its tests, `pkg/maven/command.go`, every
`maven.RunOn` caller, relevant Maven tests, the Q1.2/Q1.3/Q1.4 implementation
in `.quality/tools`, and the P2A compatibility and host acceptance contracts
before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with a recording argument-order contract. Prove the executable
   stays before every argument, the full argument values are preserved, the
   project path becomes the process working directory, the logger stdout writer
   reaches the dependency unchanged, the complete dependency value reaches the
   production callback, safe defaults execute no process, and an empty recorded
   population fails.

2. Reuse the existing `process.Runner` boundary and migrate `exec.Command` plus
   `cmd.Run` from `pkg/maven/command.go` without changing public signatures or
   behavior. Keep exact executable, argv, `Dir`, inherited environment/stdin,
   nil stderr, logger stdout, and returned error semantics. Do not use a stored
   function dependency: the type-aware scanner requires the resolvable Runner
   interface. Bind the immutable `maven-process` label through the declared
   non-executable `scripts/mutate-maven-sorting` seam driver and a Q0.8
   meta-test, without claiming the later P5 harness.

3. Run focused adapter/Maven/argument-swap tests, Q1.2/Q1.3/Q1.4 measurements,
   API/CLI compatibility, all four host acceptance flows, the full checkpoint
   gate, and the empty-HOME count-2 test from a clean commit. Expect the two
   direct Maven process effects to leave Q1.3 and `maven-process` to raise Q1.4
   from 2 of 8 to 3 of 8; accept only the scanner's regenerated exact numbers.
   The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Maven adapter
move or record an exact resumable state. Rewrite the rolling handover, update
the P3 measurements, answer this archive, create one linked NEXT archive,
replace the launcher's mutable regions, run the launcher contract, and make the
separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active for the ordered cloud and later adapter moves. Do not launch the
next session. COMPLETE is valid only after every authorized checkpoint through
P8 is complete.
<!-- CODEX_SESSION_PROMPT_END -->
