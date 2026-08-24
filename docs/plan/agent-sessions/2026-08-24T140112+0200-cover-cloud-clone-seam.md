# Agent Session: Cover Cloud Clone Seam

Status: NEXT
Session ID: `2026-08-24T140112+0200-cover-cloud-clone-seam`
Created: `2026-08-24T14:01:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `850562d0c3cd6251155b0afb50f0c3ddf071a2372efc37b45ee515150cf2ac16`
Previous: [2026-08-24T130235+0200-migrate-maven-process-flow.md](2026-08-24T130235+0200-migrate-maven-process-flow.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: expose the thinnest injectable
Git boundary for `GitCloudConfig.Refresh`, preserve its cache-first clone/pull
behavior, and kill the declared `cloud-clone` URL/target argument swap. Keep
Q1.2 at zero and preserve every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/config/cloud.go` and its tests, local-config parsing,
all `CloudConfig.Refresh` callers, the existing shell/process adapter and
recording tests, the Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A compatibility,
subprocess, and host acceptance contracts. Regenerate ignored reports outside
the measured tree or remove them before the clean audit. Implementation commit
`f5ee37d` has 50 tests across 11 of 23 packages, Q1.2 at 0, Q1.3 at 76
violations of 78 production effect sites with four adapter paths absent, and
Q1.4 at 3 of 8. The exact Q2.1 validator remains 0 of 8 executable harnesses;
the upstream filename-only denominator sees two non-executable `mutate-*`
paths. The clean gate passed, the full audit exited 1 for 16 documented
findings and never 2, and comparable ratchets were five improved, two held, and
zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the cloud-clone seam. Do not push, merge, publish,
remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands. Do not invoke a publisher or distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve the exported
`CloudConfig` interface, `GitCloudConfig`, `OpenGitCloudConfig`, every Refresh
caller, logging, and formatted errors. Production must still pull the target
when `<target>/.git` exists and otherwise clone the configured URL to that
target through the existing Git/process path.

Migrate only the cloud Refresh clone/pull dependency needed to cover the
declared swap. Do not duplicate Git argv or change the generic shell, Maven,
diagrams, profile/editor, browser-opening, HTTP, filesystem, clock, or server
effects. Do not enter Docker, dependencies, Spring, distribution, or formal
mutation-harness scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/config/cloud.go`, its tests and local-config types, every Refresh caller,
`pkg/shell/git.go` and its recording tests, `internal/adapter/process`, the
Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`, and the P2A
compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording clone/pull contracts. Prove clone keeps the complete
   configured URL before the complete target directory, an existing `.git`
   path selects pull without clone, dependency errors retain their formatted
   behavior, safe defaults perform no Git mutation, the complete dependency
   reaches Refresh, and an empty recorded population fails.

2. Introduce the thinnest resolvable interface boundary for the Refresh Git
   dependency and route production through the existing `shell.GitClone` and
   `shell.GitPull` behavior. Do not use stored function dependencies or
   reconstruct Git argv in config. Preserve the exact filesystem existence
   probe, branch selection, logging, URL/target order, pull target, and error
   semantics. Bind the immutable `cloud-clone` label through the declared
   non-executable `scripts/mutate-config-cloud` seam driver and a Q0.8
   meta-test, without claiming the later P5 harness.

3. Run focused config/shell/process/argument-swap tests, Q1.2/Q1.3/Q1.4/Q2.1
   measurements, API/CLI compatibility, all four host acceptance flows, the
   full checkpoint gate, and the empty-HOME count-2 test from a clean commit.
   Expect Q1.2 to stay zero, `cloud-clone` to raise Q1.4 from 3 of 8 to 4 of 8,
   and exact Q2.1 to remain 0 of 8. Accept only regenerated scanner values. The
   full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent cloud-clone seam
or record an exact resumable state. Rewrite the rolling handover, update the P3
measurements, answer this archive, create one linked NEXT archive, replace the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active for later adapter moves. Do not launch the next session.
COMPLETE is valid only after every authorized checkpoint through P8 is
complete.
<!-- CODEX_SESSION_PROMPT_END -->
