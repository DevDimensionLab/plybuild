# Agent Session: Repair Q1.7 Partial Failures

Status: NEXT
Session ID: `2026-08-26T204147+0200-repair-q17-partial-failures`
Created: `2026-08-26T20:41:47+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `71727d1e6f93f11f45bba986fefe84a0f48dee5f64e07479691a892111422801`
Previous: [2026-08-26T195729+0200-repair-q19-empty-populations.md](2026-08-26T195729+0200-repair-q19-empty-populations.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one Q1.7-only test move. Repair the eight confirmed operations
whose per-item failure content, continuation, or ordered partial result lacks
an exact executable contract. Reclassify all 17 in-scope operations from
executable control flow and do not infer coverage from the 23 existing focused
tests or from exit-code assertions alone.

# Authorized Roadmap

P3 is complete. P4 remains active only for Q1.7 remediation and later combined
manual evidence; P5-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`. Q1.6 is independently supported by its complete
57-struct, 74-field review, and Q1.9 is independently supported by its complete
169-range classification. Neither support permits a schema-2 receipt in this
test-only session. The launcher must remain NEXT until P4 and P5-P8 finish.

This mission authorizes read-only reclassification of the 17 Q1.7 operations;
test-only characterization of `initCmd`, `Bitbucket.SynchronizeAllRepos`,
`spring.DeleteDemoFiles`, `Repository.upgradeDependencies`,
`Repository.upgradePluginsOnModel`, `maven.RemoveDeprecated`,
`template.MergeTemplates`, and `GitCloudConfig.ValidTemplatesFrom`; focused
test helpers needed to exercise their existing behavior; one focused test-only
implementation commit; clean automated measurement; and the normal separate
continuity commit. It does not authorize production changes, another seam,
manual evidence or receipts, Q1.9 changes, audit/parser/scanner changes,
`.quality/inventory` or baseline changes, executable P5 harnesses, T1-T10
changes, acceptance changes, P6-P8 work, or a Q1.7 PASS claim unless the
complete classified population supports it.

# Measurements At Start

The Q1.9 move produced test-only commit
`3a8f0bcd030baf787a29440ee8e4e4a087edb33b`, exact tree
`75dee3dcc1cb2fd91f7eb359e96b31cc4be3d986`, on unchanged production commit
`c382caa38be167fe17f847370ad8a12270644de3`. Its exhaustive 169-range record
classifies 72 repaired, 83 independently guarded, and 14 fixture/support range
sites; all 64 local `tests` loops are now directly guarded. Q1.9 is supported,
but no manual-evidence document or receipt was created.

The clean no-evidence scorecard SHA-256 is
`08e5ff9d450a6fa3b8824d82c720d339cf0d44003a490d01d3a412681081a615`.
It records 453 tests across all 27 packages, L0 8/8, L1 six PASS and
Q1.6/Q1.7/Q1.9 UNMEASURABLE, six improved ratchets, two held, zero regressed,
zero non-comparable, and zero dirty paths. The audit exits 1 for 13 documented
non-passing criteria, never 2.

# Q1.7 Scope And Blockers

Include a production collection-item failure path when it reports the failure
and continues, and an operation when it can return a populated partial
aggregate with an error. Exclude fail-fast errors before partial state,
intentionally ignored errors with no reported partial result, non-error
business warnings, sequential non-collection phases, and test support. Inspect
executable control flow; names and search hits are not classification evidence.

The complete scope is 17 operations. Nine have exact content coverage. The
eight confirmed gaps are the `initCmd` project-loop failures;
`SynchronizeAllRepos` clone/pull warning; `DeleteDemoFiles` discovery/deletion
failures; dependency version, maximum-version, and upgrade failures; per-plugin
upgrade failures; `RemoveDeprecated` partial-result and warn/continue paths;
`MergeTemplates` merge warning; and `ValidTemplatesFrom` ordered partial result
with error. A new contract must assert the exact observable failure content and
also prove later-item continuation or the exact populated partial result where
that behavior exists. Prefer existing doubles and log-capture helpers. Preserve
all production behavior and process-state restoration.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Confirm the new
continuity HEAD and exact ancestry before editing. Make one focused test-only
implementation commit, then the separate `docs: prepare next agent session`
continuity commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

If any confirmed blocker cannot be exercised through current seams without a
production change, exhaust safe test-only approaches, record the exact
controllability gap, and keep Q1.7 blocked. Do not change production or add a
seam under this mission. Preserve test names and data already in place, exact
failure strings and ordering, public behavior, side effects, callback order,
and Go 1.18 compatibility.

Keep every generated report outside the worktree. Set
`API_COMPAT_REPORT_OUT` and `CLI_COMPAT_REPORT_OUT` to external paths whenever
running compatibility. The measured-tree identity counts ignored, untracked,
and tracked bytes, so confirm both ordinary and ignored status before every
authoritative audit.

# Required Reading

Before editing, confirm branch, HEAD, clean and ignored status, reciprocal
archive links, launcher `--check`, exact ancestry, and the authorized
checkpoint block. Read the rolling handover, this archive, the complete P4
entry and gate, both design documents, `.quality/README.md`, the Q1.7 criterion
and schema-2 validation path, all 17 production operations and enclosing call
flows, every existing partial-failure test and helper, the complete audit
wrapper, T15 repair, API/CLI and launcher contracts, and Make and acceptance
gates. Do not classify an operation by name or search hit alone.

# Three Moves

1. Regenerate the complete 17-operation Q1.7 record from the clean current
   tree. For each operation record path, function, every included failure
   branch, existing exact-content/continuation coverage, new contract, or exact
   exclusion outside the worktree. Add focused test-only contracts for every
   exercisable branch among the eight gaps. Require exact failure content and
   later-item continuation or exact ordered partial results as applicable.
2. Run the exact affected tests, the complete 17-operation focused population,
   complete uncached tests, race, and vet. Inspect every new assertion in its
   executable context. If any included failure path can still occur without an
   exact content assertion and continuation/partial-result proof, keep Q1.7
   blocked. Do not create manual evidence in this session.
3. Run API/CLI and subprocess compatibility with external report paths; pinned
   lint; launcher and Make contracts; complete preflight; all four host
   acceptance flows; the 15-control audit meta-suite; empty-HOME count-2; and a
   clean no-evidence audit. Require audit exit never 2 and zero dirty paths.
   Record the exact population, implementation commit, scorecard, and next P4
   blocker or evidence-only next step in continuity.

# Automatic Handoff

Before this session ends, finish and commit the coherent Q1.7 test-only move or
record an exact resumable blocker. Rewrite the rolling handover, record the P4
result in the roadmap, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent P4 move, replace only the launcher's
mutable regions, run launcher and handoff contracts, and make the separate
continuity-only commit. Do not launch a real successor, push, merge, publish,
distribute, stash, revert, or remove the worktree. COMPLETE remains invalid
while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
