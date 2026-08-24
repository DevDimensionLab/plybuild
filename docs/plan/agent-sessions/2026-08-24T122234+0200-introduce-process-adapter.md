# Agent Session: Introduce Process Adapter

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T122234+0200-introduce-process-adapter`
Created: `2026-08-24T12:22:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `48374146595f2d319e21f5f460cac71798a1b4a384311d83438b74725044bd33`
Previous: [2026-08-24T113208+0200-remove-process-exits.md](2026-08-24T113208+0200-remove-process-exits.md)
Next: [2026-08-24T130235+0200-migrate-maven-process-flow.md](2026-08-24T130235+0200-migrate-maven-process-flow.md)
Outcome: P3 move 2 completed at 03d6242: the process adapter owns Git clone/pull/init/add/commit execution, Q1.3 fell to 77 of 79, both Git seam swaps raised Q1.4 to 2 of 8, and the complete clean gate passed with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the second P3 move in one measured change: introduce
`internal/adapter/process`, migrate the coherent git process flow in
`pkg/shell/git.go`, and kill the declared `git-process` and `git-commit`
argument swaps. Reduce Q1.3 by the migrated git sites while keeping Q1.2 at
zero and preserving all P2A contracts with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/shell/git.go`, `pkg/shell/command.go`, their callers,
the Q1.3/Q1.4 scanners, existing process tests, and the API/CLI/subprocess and
host acceptance contracts. Regenerate ignored reports. Implementation commit
`5c6f2fa` has 39 tests across 9 of 22 packages, Q1.2 at 0, Q1.3 at 80 of 80
direct sites outside five missing adapters, and Q1.4 at 0 of 8 covered seams.
Its clean gate passed and its full audit exited 1, never 2. The restart handoff
rephrased one inherited Q3.4 host-state sentence; require four improved, three
held, and zero regressed ratchets before accepting this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the git process-adapter move. Do not push, merge,
publish, remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands. Do not invoke a publisher or distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Default production behavior must
still execute the same git argv in the same directories. Recording doubles must
have safe defaults, capture complete arguments and working directories, and
fail when the asserted call population is empty.

Migrate only the coherent git flow in `pkg/shell/git.go`. Do not migrate Maven
or cloud process calls in the same move. Do not introduce HTTP, filesystem,
clock, or server adapters, and do not enter Docker, dependencies, cloud,
Spring, distribution, or mutation-harness scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/shell/git.go`, `pkg/shell/command.go`, every Git helper caller, the
Q1.3/Q1.4 implementation in `.quality/tools`, and the P2A compatibility and
host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording argument-order contracts. Prove clone keeps URL
   before target directory, add/commit keeps the target directory before the
   message, the complete dependency value reaches the production path, safe
   defaults perform no real process in tests, and an empty recorded population
   fails.

2. Add the thinnest `internal/adapter/process` boundary needed by the Git flow
   and migrate the direct Git clone, pull, init, add, and commit execution paths
   without changing public signatures or `shell.Output` behavior. Keep exact
   argv and working-directory semantics. Do not migrate the generic shell or
   Maven/cloud process paths in this move.

3. Run focused adapter and argument-swap tests, Q1.2/Q1.3/Q1.4 measurements,
   API/CLI compatibility, all four host acceptance flows, the full checkpoint
   gate, and the empty-HOME count-2 test from a clean commit. Expect the three
   direct `pkg/shell/git.go` process sites to leave Q1.3 and both declared git
   swaps to raise Q1.4; accept only the scanner's regenerated exact numbers.
   The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent git adapter move
or record an exact resumable state. Rewrite the rolling handover, update the P3
measurements, answer this archive, create one linked NEXT archive, replace the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active for the ordered Maven/cloud and later adapter moves. Do not
launch the next session. COMPLETE is valid only after every authorized
checkpoint through P8 is complete.
<!-- CODEX_SESSION_PROMPT_END -->
