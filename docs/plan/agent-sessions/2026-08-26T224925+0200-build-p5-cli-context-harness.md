# Agent Session: Build P5 CLI Context Harness

Status: NEXT
Session ID: `2026-08-26T224925+0200-build-p5-cli-context-harness`
Created: `2026-08-26T22:49:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `531f519e3df4d31d9d7b8683fc3ba26eace358c11d07e27f0ddc885fac0a3125`
Previous: [2026-08-26T220403+0200-record-p4-combined-manual-evidence.md](2026-08-26T220403+0200-record-p4-combined-manual-evidence.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P5 with exactly one mutation-evidence subject: `cli-context`. Create its
executable mutation harness and falsifiability meta-test, run at least eight
meaningful mutations against disposable external copies, and finish the move
only if every declared mutation is killed with zero survived and zero unusable.
Do not treat compilation failures, empty test selections, or a sampled subset
as mutation evidence.

# Authorized Roadmap

P2A, P2B, P3, and P4 are complete. P5 is active; P6-P8 remain queued in the
machine-readable block in `docs/plan/quality-upgrade.md`. The first
`.quality/inventory` subject is `cli-context`, with production roots `cmd` and
`pkg/context` and declared harness `scripts/mutate-cli-context`. This session
may create only that harness and `scripts/test-mutate-cli-context`, plus the
smallest focused in-subject test or private seam repair if an actually executed
mutation exposes a classified gap. P5 remains active after this first subject.

This mission does not authorize another mutation subject, `.quality/inventory`
changes, audit/parser/scanner/baseline changes, acceptance expansion, P6-P8,
Go or dependency upgrades, exported API or CLI changes, packaging,
publication, or distribution. Keep every mutant checkout, cache, report, and
generated artifact outside the worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Measurements At Start

The P4 evidence checkpoint is clean continuity commit
`88a95ad6effe7c2198d9505963dc19022bcabc10`, exact parent
`4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`, and tree
`120d8db0b83a7724b0026f7c18a77a706e9f1cda`. The clean status SHA-256 is
`6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`.
After launch, the new continuity HEAD must have exact parent `88a95ad`.

The external schema-2 document has SHA-256
`c97c0336998508fe410fc3f56b65cc41762281850ba5822d53c623dc74be34c0`.
Its Q1.6, Q1.7, and Q1.9 receipt SHA-256 values are respectively
`85db15d942b850f99baa518428043533ceee7b94d911b2de94b6e45a5f3f0ad7`,
`e94221777b1f8ea0300cdc286fbb7086440d1038b5b9e9da0f7d845d18845e2d`,
and `6d4bf0c4429fa7a731078b84bc850d401a1c636f4cf72bb8d110dfe4700d9616`.
The focused audit exits 0 with scorecard SHA-256
`63ede0de00d52d950f10845ee8d6e16c1e948e73c8dd183c368ef0e87cc4c729`.
The full audit exits 1, never 2, with scorecard SHA-256
`52490a43d26c5b6c1e6031460352a9fbc6de31bba52128301d3d3e1f31d86bc6`:
L0 is 8/8, all nine L1 rows PASS, six ratchets improve, two hold, none
regress, and dirty paths are empty. Its nine remaining non-passing rows are
P5-P8 debt.

Fresh P4 gate results pass API/CLI and subprocess compatibility, pinned lint,
all 27 package tests, race, vet, launcher and Make contracts, all four host
acceptance flows, empty-HOME count-2, complete preflight, and the 15-control
audit meta-suite. Current Q1.9 structure contains 183 total test ranges: 168
assertion-required and 15 fixture/support exclusions.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Confirm the new
continuity HEAD and exact ancestry before editing. Preserve current behavior,
public Go API, Cobra surface, output, ordering, error text, exit status, and
compatibility. Inspect all relevant `cmd` and `pkg/context` production paths
and their complete test populations before selecting mutations.

The production harness must be a regular executable script and the meta-test
must prove its contract can fail. Each declared mutation must be deterministic,
unique, meaningful, applied exactly once inside the authorized production
roots, and run in a fresh disposable copy. Require a clean unmodified control,
a non-empty exact selected test population for every mutant, and evidence that
the selected tests actually run. Do not count a mutation as killed merely
because source selection, setup, compilation, tooling, or the harness failed.

Implement and name the roadmap's T1-T10 methodology controls. At minimum they
must fail closed on an empty or duplicate mutation manifest, out-of-scope or
zero/multiple source replacements, empty test selection, a broken clean
control, a mutant that was not compiled and exercised, falsified result
accounting, a surviving mutation without classification, repository-local
artifacts, and non-deterministic declared/killed totals. The final report must
state exact `declared`, `killed`, `survived`, and `unusable` counts and identify
every mutation and its killing test population.

If a genuine survivor appears, classify it as a reachability, observability,
or controllability gap before changing tests or adding a private seam. Make
only the smallest `cmd` or `pkg/context` repair needed to kill that mutation,
then rerun the complete harness. Do not manufacture easy mutants around known
assertions or weaken a mutation until it passes.

# Required Reading

Before editing, confirm branch, HEAD, clean and ignored status, reciprocal
archive links, launcher `--check`, exact ancestry, and the authorized
checkpoint block. Read the rolling handover, this archive, the complete P5
entry and checkpoint gate, both design documents, `.quality/README.md`,
`.quality/inventory`, the complete mutation-harness discovery and Q2.1-Q2.4
logic in the vendored audit and structured parser, the six existing
non-executable P3/P4 seam drivers and their meta-tests, Make/preflight script
population contracts, and all production/tests under `cmd` and `pkg/context`
that establish the proposed mutant population. Do not infer reachability or a
kill from names, grep results, compilation failure, or exit status alone.

# Three Moves

1. Define an explicit `cli-context` mutation scope and non-empty manifest of at
   least eight behaviorally meaningful mutations. Map each mutation to exact
   production syntax and the non-empty named tests expected to kill it. Record
   exclusions and reject any mutation that cannot be applied exactly once.
2. Implement `scripts/mutate-cli-context` and
   `scripts/test-mutate-cli-context`. Run the unmodified control and every
   mutant in isolated external copies; prove all T1-T10 negative controls; fix
   only a classified in-subject gap if required. Require exact
   `declared == killed`, `survived == 0`, and `unusable == 0`.
3. Make one focused implementation commit, then measure it cleanly with the
   focused harness/meta-test, automated Q2.1-Q2.4 audit view, API/CLI and
   subprocess compatibility, pinned lint, complete tests/race/vet, launcher
   and Make contracts, complete preflight, host acceptance, audit meta-suite,
   and empty-HOME count-2. Record the result and prepare the next bounded P5
   subject with one continuity-only commit.

# Automatic Handoff

Before this session ends, finish the coherent `cli-context` mutation-harness
move or record an exact resumable blocker. Rewrite the rolling handover, record
the result in the roadmap, answer this archive, create exactly one reciprocally
linked NEXT archive, replace only the launcher's mutable regions, run launcher
and handoff contracts, and make the normal continuity-only
`docs: prepare next agent session` commit after the focused implementation
commit. Do not launch a real successor, push, merge, publish, distribute,
stash, revert, or remove the worktree. COMPLETE remains invalid while P5-P8
are unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
