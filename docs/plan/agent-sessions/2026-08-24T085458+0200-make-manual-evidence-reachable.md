# Agent Session: Make Manual Evidence Reachable

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T085458+0200-make-manual-evidence-reachable`
Created: `2026-08-24T08:54:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5eb956f83bd5c5f398f196a0d1cd068eb26917c50ae673e07a6a811c95619023`
Previous: [2026-08-24T061532+0200-close-absolute-l0.md](2026-08-24T061532+0200-close-absolute-l0.md)
Next: [2026-08-24T103558+0200-characterize-compatibility.md](2026-08-24T103558+0200-characterize-compatibility.md)
Outcome: P1B completed in exactly three measured moves at eb987fd: all six manual rows are reachable without overriding automated verdicts or empty populations, the baseline migration preserved all 228 numeric debt leaves and Q3.9, and the checkpoint ended with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete checkpoint P1B: make the six manual L1/L2 rows reachable through
commit-, tree-, inventory-, and instrument-bound structured evidence, then
migrate the exact-toolchain baseline without changing any numeric debt. Use no
more than three measured moves and leave zero comparable ratchet regressions.

# Changes Since The Previous Prompt

1. P1 closed absolute L0 at 8 of 8 and added the non-publishing daily
   `make preflight` gate in exactly three quality moves.
2. golangci-lint is pinned at `v2.12.2`; its five-linter gate is clean and
   remains separate from source-rewriting `make format`.
3. The sole production script has an independent negative meta-test, so Q0.8
   is 0 of 1 scripts without a counterpart.
4. The clean P1 audit reports two improved, five held, and zero regressed
   ratchets; its exit 1 is caused only by 20 documented findings above L0.
5. The next ordered checkpoint is P1B. P2 compatibility, publishing, cloud,
   Spring, packaging, and dependency work remain out of scope.

# Measurements At Start

Run these commands before editing. Treat command output as evidence. Inspect
status and diffs first whenever the launcher warned that the tree is dirty.

```sh
git status --short --branch
git diff --stat
git rev-parse --show-toplevel
git rev-parse --short=12 HEAD
git branch --show-current
sed -n '1,300p' docs/plan/quality-handover.md
sed -n '120,240p' docs/plan/quality-upgrade.md
sed -n '1,260p' docs/design/quality-lift.md
sed -n '1,260p' .quality/baseline/README.md
bash .quality/tools/test-quality-audit.sh
bash .quality/tools/quality-audit.sh . \
  --baseline .quality/baseline/scorecard.json
```

The full audit may exit 1 for documented findings above L0. Exit 2 invalidates
the measurement. Regenerate reports; do not use an old ignored report as
evidence. The launcher never inserts Git paths or other worktree content into
this prompt.

# Role, Permissions, And Stop Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Inspect before
editing, protect inherited changes, use focused tests, and make one focused
commit per quality move. You may update implementation, tests, quality
apparatus, design decisions, the plan, and the rolling handover when the active
objective requires it.

Do not mutate this launcher's session regions or any session archive until the
user sends the exact restart trigger. Queue prompt corrections in the rolling
handover meanwhile. Do not push, merge, release, publish, remove the worktree,
stash inherited changes, revert user work, or run destructive Git commands.
Never run `make release` or `make release-brew`. Stop and report when the audit
exits 2, a comparable ratchet regresses, compatibility cannot be established,
later cloud/Spring/packaging scope becomes necessary, or three moves are done.

# Required Reading

Read in this order:

1. `docs/plan/quality-handover.md` for the exact P1 exit and resumption state.
2. The P1B section and checkpoint gate in `docs/plan/quality-upgrade.md`.
3. `docs/design/quality-lift.md` for evidence and ratchet invariants.
4. `.quality/baseline/README.md`, `manual-evidence.json`, and `scorecard.json`
   for the current schema, receipts, instrument identity, and reproduction.
5. Manual-evidence validation and baseline comparison in
   `.quality/tools/scorecard.py`, plus T6b-T8b and T15 in
   `.quality/tools/test-quality-audit.sh`.
6. `.quality/inventory` and `docs/design/agent-session-continuity.md` before
   changing instrument-bound evidence or preparing another restart.

Do not use `.agent-task/current.md` as task authority. If a statement here is
wrong, record the correction in its owning plan/handover document and carry it
into the next prompt only during an authorized restart.

# Environment Constraints

- Use `LC_ALL=C LANG=C` for deterministic shell tooling.
- Use fresh `GOCACHE` and `GOLANGCI_LINT_CACHE` directories under
  `/private/tmp`; both shared caches denied access during restart measurement.
- Probe Docker and shellcheck availability before claiming evidence from them.
- Keep tests independent of developer HOME, public services, and repository
  writes.
- The module stays at Go 1.18. The exact baseline context is Go 1.26.2 on
  Darwin arm64 with the build selectors recorded in the baseline README.
- Parser or schema changes alter instrument identity. Reproduce commit
  `5635d50` with both old and new instruments, preserve every numeric debt
  value, and record both identities; never merely bless a regenerated file.
- A receipt may resolve only an upstream `UNMEASURABLE` row. An explicit
  upstream `FAIL`, an empty declared population, or invalid evidence must not
  become PASS.
- Homebrew remains active in release configuration despite being outside the
  accepted distribution matrix; do not execute publication paths.

# First Task

Implement P1B as at most three focused moves:

1. Start with failing audit meta-tests, then evolve the structured evidence
   schema so valid, non-empty receipts can resolve Q1.6, Q1.7, Q1.9, Q2.4,
   Q2.8, and Q2.9. Bind receipts to module, commit, measured tree, inventory,
   instrument identity, criterion, and evidence digest. Reject stale,
   duplicate, empty, wrong-kind, dirty-tree, and false-PASS evidence.
2. Prove precedence and populations independently: an upstream FAIL always
   wins, and rows with no underlying harness/script/test population remain
   non-passing. Synthetic non-empty fixtures may demonstrate reachability;
   do not claim the current project passes rows whose population is zero.
3. Perform an explicit debt-preserving baseline migration. Reproduce commit
   `5635d50` under the old and new instruments, compare every numeric metric,
   record both hashes and the reproducible recipe, and retain Q3.9 evidence.

After each move, run its focused tests and measure from a clean commit. At the
checkpoint, run the complete gate and empty-HOME count-2 command from the plan.

# Restart Protocol

Suggest an agent restart at a natural checkpoint, after three moves, or when
context quality starts falling. Preparation starts only when the user's
trimmed, case-sensitive message is exactly `agent-restart`. That message
authorizes restart preparation and one local allowlisted commit; it does not
authorize push, merge, release, stash, revert, or staging unrelated changes.

On `agent-restart`:

1. Read
   `/Users/perottochristensen/github/spk-.agents/skills/spk-agent-restart/SKILL.md`
   when available, then follow this repository protocol as controlling.
2. Re-measure branch, HEAD, status, tests, and quality. Finish a coherent slice
   when feasible. Preserve wrong claims, expensive findings, blind alleys,
   tests not run, and unfinished work when the next agent needs them.
3. Update design notes only for durable decisions or premises, the quality
   plan for roadmap or measured-state changes, and rewrite the rolling
   handover with the exact resumption state.
4. Replace only the two mutable regions below
   `CODEX_STABLE_EXECUTION_END`. Keep every data line prefixed with `#|`, keep
   the normalized stable skeleton unchanged, and set one concrete next
   mission, new session ID/archive path, numbered delta, measurements,
   boundaries, reading order, environment constraints, task, and protocol.
5. Create
   `docs/plan/agent-sessions/YYYY-MM-DDTHHMMSS+ZZZZ-<checkpoint>.md` with
   status `NEXT`, timestamp matching its ID, reciprocal previous link,
   SHA-256 of the prompt block, and the byte-exact decoded launcher prompt,
   including its terminal LF, between standalone
   `<!-- CODEX_SESSION_PROMPT_BEGIN -->` and
   `<!-- CODEX_SESSION_PROMPT_END -->` markers.
6. Mark the prior archive `ANSWERED - HISTORY`, set a non-pending outcome and
   reciprocal next link, and preserve its prompt block and digest unchanged.
   There must be exactly one connected archive chain and one `NEXT` tail.
7. Run `make test-agent-start`, focused tests for changed behavior, and the
   relevant checkpoint gates before staging. A failed required check blocks
   restart preparation unless its external blocker and missing evidence are
   explicit in the handover.
8. Stage only the launcher, plan/design files changed for this checkpoint,
   the prior archive, and the new archive. Commit once with
   `docs: prepare next agent session`. Never stage inherited dirty files.
9. Run `./codex-dev-start.sh --check` after the commit, report the commit and
   next mission, tell the user to run `./codex-dev-start.sh`, then stop. Do not
   launch the next session yourself.

When the total objective is complete and no follow-up objective has explicit
approval, answer the final archive with a non-pending outcome, leave `Next:
none`, and set `SESSION_STATUS=COMPLETE`. `--check` must pass, while start and
`--print-prompt` fail closed instead of replaying the answered task.
<!-- CODEX_SESSION_PROMPT_END -->
