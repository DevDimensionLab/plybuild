# Agent Session: Cover Resources Package

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T160625+0200-cover-resources-package`
Created: `2026-08-26T16:06:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8d08e706ab7b69c3da7f7766dacd013556a6a1e576af143479a1f6a703128230`
Previous: [2026-08-26T153936+0200-cover-sorting-package.md](2026-08-26T153936+0200-cover-sorting-package.md)
Next: [2026-08-26T163426+0200-cover-webservice-templates-package.md](2026-08-26T163426+0200-cover-webservice-templates-package.md)
Outcome: Product `64c6189` added three deterministic resource path/read characterization contracts at 100% package coverage; Q1.1 improved to 4/27 and the clean full audit exited 1 for 13 documented findings with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the untested logic
package `pkg/resources`. Add deterministic characterization contracts for its
complete existing resource-directory and file-content behavior without changing
production code, exported API, callers, filesystem boundaries, or observable
semantics. Preserve every completed checkpoint and produce zero comparable
ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.4 product commit `cbb620a` added the complete
sorting characterization and improved Q1.1 from 6/27 to 5/27 while exact Q1.3
held at zero of 27.

This mission authorizes only focused tests in `pkg/resources` that characterize
the existing `LocalDir` and `ResourceAsString` implementations. It does not
authorize changing `pkg/resources/resources.go`, adding a dependency seam or
inventory label, changing config, file, template, Maven, or webservice callers,
modifying scanners, audit apparatus, mutation harnesses, acceptance scripts,
public APIs, templates, logging, filesystem or process behavior, another
untested package, manual evidence, or P5-P8 implementation. If truthful tests
reveal a product decision rather than existing behavior, preserve the behavior
and record the decision for a later authorized move.

# Measurements At Start

Clean product commit `cbb620a` has 423 tests across 22 of 27 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 5/27, Q1.2 is zero, exact Q1.3 is 0/27 with
all five declared adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is
zero phrases across 88 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`4f7550ad8fd184f09237c5f8ae339a1d746ab8c2b366706d4cd9c5feb79ae0af`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make
the normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Tests must preserve the exact
`config.CloudConfig` parameter and exported `LocalDir(config.CloudConfig)
string` and `ResourceAsString(config.CloudConfig, string) (string, error)`
function signatures.
Characterize the exact single `Implementation` and `Dir` selection, existing
`file.Path` resource-directory construction, literal resource/filename slash
composition, arbitrary root and filename bytes supported by the host,
successful byte-to-string conversion, empty and non-UTF-8 content, nested
resource names, and direct missing-read error behavior without wrapping or
replacement. Do not clean or normalize legacy path composition in production.

Focused tests must make no request, open no socket, launch no external program,
perform no wait, mutate no global logger or production state, and write no
repository fixture. Any required resource file must live under `t.TempDir()`
and be created through the central safe writer after its parent directory is
prepared. Reject empty table-driven populations before iterating. Leave
sorting, server, clock, Kibana, Spring, Maven, template production, filesystem,
process, HTTP, browser, plugin-diagrams, profile, shell, Unzip, and every
completed contract unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `cbb620a`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, complete `pkg/resources/resources.go`, every
resources caller and its relevant tests, the complete `config.CloudConfig`,
`config.Directory`, `GitCloudConfig`, and `DirConfig` shapes, `file.Path` and
`file.Open`, representative safe temporary-fixture and explicit non-empty
population contracts, the import-aware scanner, API/CLI contracts, and the T15
repair and baseline reproduction README.

# Three Moves

1. Add focused `pkg/resources` contracts that prove exact implementation and
   directory selection counts, `LocalDir` path bytes for representative, empty,
   and trailing-separator roots, exact successful resource reads for empty,
   arbitrary, non-UTF-8, and nested content, direct missing-file result and
   error shape, and non-empty test populations. Do not change production code
   to make a preferred normalized path or injected read pass.
2. Run the focused package and relevant template caller tests, inspect coverage,
   and confirm that only the new resources test file changed. Regenerate focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without changing the scanner or inventory.
   Expect Q1.1 to improve from 5/27 to 4/27 and Q1.3 to hold at 0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent resources
coverage move or record an exact resumable state. Rewrite the rolling handover,
record the measured P4 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch
a real successor, push, merge, publish, distribute, stash, revert, or remove
the worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
