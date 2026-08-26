# Agent Session: Cover Sorting Package

Status: NEXT
Session ID: `2026-08-26T153936+0200-cover-sorting-package`
Created: `2026-08-26T15:39:36+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b99cfa23f141551c169e130d4b6438aae212a419a914b07d7d44590f10447ee8`
Previous: [2026-08-26T145252+0200-migrate-web-server-effects.md](2026-08-26T145252+0200-migrate-web-server-effects.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the untested logic
package `pkg/sorting`. Add deterministic characterization contracts for its
complete existing dependency ordering behavior without changing production
code, exported API, callers, or observable sort semantics. Preserve every
completed checkpoint and produce zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.3 product commit `9601d29` completed the
declared server boundary and reduced exact Q1.3 to zero of 27.

This mission authorizes only focused tests in `pkg/sorting` that characterize
the existing `DependencySort` implementation and its private comparison
helpers. It does not authorize changing `pkg/sorting/sort.go`, fixing or
normalizing surprising legacy ordering, changing Maven/config callers, adding
an adapter or seam label, modifying inventory, scanners, audit apparatus,
mutation harnesses, acceptance scripts, public APIs, templates, logging,
filesystem or process behavior, the loopback-address policy, another untested
package, manual evidence, or P5-P8 implementation. If truthful tests reveal a
product decision rather than existing behavior, preserve the behavior and
record the decision for a later authorized move.

# Measurements At Start

Clean product commit `9601d29` has 414 tests across 21 of 27 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/27, Q1.2 is zero, exact Q1.3 is 0/27 with
all five declared adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is
zero phrases across 87 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`010923f3718d99723c99c76a6b541669cb7f99d7a4a8eb119044705be005cdfa`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make
the normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Tests must preserve the exact
`DependencySort` slice and `SortKey` API, value-receiver behavior, length,
in-place swap semantics, and strict less-than comparison. Characterize the
existing concatenated key order exactly: decimal scope weight, colon, group
weight, colon, and artifact ID. Preserve exact scope weights for empty,
`compile`, `provided`, `runtime`, `system`, `import`, `test`, and unknown
values; exact case-sensitive `strings.Contains` handling of the sort key;
matching `1-` and nonmatching `100-` group prefixes; group and artifact
tie-breaks; empty sort-key behavior; and any legacy lexical consequence of
formatting weight 100 as text. Do not silently assert a preferred numeric or
stable order that the implementation does not provide.

Focused tests must make no request, open no socket, launch no external program,
perform no wait, mutate no global logger or production state, and write no
repository fixture. Reject empty table-driven populations before iterating.
Leave server, clock, Kibana, Spring, Maven production, filesystem, process,
HTTP, browser, plugin-diagrams, profile, shell, Unzip, and every completed
contract unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `9601d29`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, complete `pkg/sorting/sort.go`, every
`sorting.DependencySort` caller and its relevant tests, the `pom.Dependency`
shape used by the repository, representative package-level characterization
tests with explicit non-empty table guards, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Add focused `pkg/sorting` contracts that prove `Len`, `Swap`, `Less`, helper
   key construction, every scope class, sort-key containment and case
   sensitivity, empty sort key, group/artifact tie-breaking, exact in-place
   full-sort results, legacy weight-100 lexical behavior, and non-empty test
   populations. Do not change production code to make a preferred order pass.
2. Run the focused package and relevant config/Maven caller tests, inspect
   coverage, and confirm that only the new sorting test file changed. Regenerate
   focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without changing the scanner or
   inventory. Expect Q1.1 to improve from 6/27 to 5/27 and Q1.3 to hold at
   0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent sorting coverage
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P4 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace
only the launcher's mutable regions, run launcher and handoff contracts, and
make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
