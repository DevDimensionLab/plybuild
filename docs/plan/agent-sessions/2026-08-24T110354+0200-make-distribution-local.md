# Agent Session: Make Distribution Local-Only

Status: NEXT
Session ID: `2026-08-24T110354+0200-make-distribution-local`
Created: `2026-08-24T11:03:54+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8a321c84978c98e5e199362a5cfcb2bd71017fb817497127cb194db75d51acba`
Previous: [2026-08-24T103558+0200-characterize-compatibility.md](2026-08-24T103558+0200-characterize-compatibility.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete checkpoint P2B in one measured move: make the supported default
distribution workflow local-only, remove inactive Homebrew/Snap publishers from
that path, and prove the snapshot command cannot publish. Leave zero comparable
ratchet regressions and do not execute any publisher.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P2A is complete, P2B is
active, and P3-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md. A queued checkpoint is approved work, so the
launcher must remain NEXT until all authorized checkpoints are complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P2B
plan, `.goreleaser.yml`, `.goreleaser.brews.yml`, the release targets, and their
current contract tests. Regenerate quality reports rather than relying on
ignored output. The clean P2A audit at 14764fd exited 1 for 16 documented
findings with two improved, five held, and zero regressed ratchets. Exit 2
invalidates a checkpoint.

# Role And Boundaries

Work autonomously in this worktree on codex/upgrade-quality. Make one focused
implementation commit for P2B. Do not push, merge, publish, remove the worktree,
stash inherited changes, revert user work, or run destructive Git commands.
Do not invoke `make release`, `make release-brew`, or an unguarded GoReleaser
release. Use static configuration checks and a recording executable for
negative publication tests.

Keep Go 1.18 and preserve the P2A API, Cobra, subprocess, and host acceptance
contracts. Do not broaden this move into Docker, dependencies, cloud, Spring,
process-exit refactoring, adapters, or mutation harnesses. GoReleaser syntax is
version-sensitive: inspect the repository context and current official,
version-matched help or documentation before binding production flags.

# Required Reading

Read docs/plan/quality-handover.md, the P2B section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.goreleaser.yml`,
`.goreleaser.brews.yml`, Makefile release targets, and
test/makefile_preflight_test.sh before editing.

# Three Moves

These are three ordered steps inside the single measured P2B move.

1. Start with a failing distribution contract. It must reject an active Homebrew
or Snap publisher in the default release configuration, reject a callable
inactive-package-manager target, and record the exact local snapshot invocation
without executing GoReleaser or using credentials.

2. Remove or isolate the `brews` publisher in `.goreleaser.yml`, disable or
explicitly fail closed the standalone brew configuration and
`make release-brew`, and expose a supported snapshot command that cannot create
a GitHub release or package-manager metadata. Preserve an explicit production
release path only if its opt-in boundary is unmistakable and contract-tested;
the ordinary/snapshot path must remain non-publishing.

3. Run focused negative tests, compatibility, host acceptance, the full checkpoint
gate, and the empty-HOME count-2 test from a clean commit. The full audit may
exit 1 for known findings, but must report zero ratchet regressions.

# Automatic Handoff

Before this agent session ends, finish and commit the one coherent P2B move or
record an exact resumable state. Rewrite the rolling handover, update checkpoint
statuses, answer this archive, create one linked NEXT archive, and replace the
launcher's mutable session regions. No separate agent-restart message is
required.

If P2B is complete, activate P3. If unfinished, keep P2B active with a concrete
resume mission. Stage only handoff files in the handoff commit named
`docs: prepare next agent session`. Do not launch the next session yourself.
COMPLETE is valid only after every authorized checkpoint through P8 is complete.
<!-- CODEX_SESSION_PROMPT_END -->
