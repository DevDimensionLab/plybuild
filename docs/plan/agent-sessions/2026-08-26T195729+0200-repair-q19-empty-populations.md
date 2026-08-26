# Agent Session: Repair Q1.9 Empty Populations

Status: NEXT
Session ID: `2026-08-26T195729+0200-repair-q19-empty-populations`
Created: `2026-08-26T19:57:29+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `971da199b2dc19466bfb32f7b48a913dbe6f33fb778149ac18d138102b2da3af`
Previous: [2026-08-26T191251+0200-record-p4-manual-l1-evidence.md](2026-08-26T191251+0200-record-p4-manual-l1-evidence.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one Q1.9-only test move. Exhaustively classify all 169 current
Go test `range` statements, then add executable assertions that fail when every
in-scope assertion-driving collection is empty. Repair all 51 confirmed local
`tests`-table blockers and every additional unguarded expectation, recording,
or helper population found in the remaining 105 ranges. Do not infer coverage
from the 69 existing empty-population contracts or from non-empty literals.

# Authorized Roadmap

P3 is complete. P4 remains active because Q1.7 and Q1.9 failed exhaustive
manual review; P5-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`. Q1.6's complete 57-struct, 74-field production
population is independently supported, but no schema-2 receipt was created
because the evidence mission required all three manual rows to be valid. The
launcher must remain NEXT until P4 and P5-P8 finish.

This mission authorizes only read-only reclassification of the 169 `_test.go`
range sites; test-only fail-on-empty assertions for every assertion-driving
table, expected collection, recorded collection, or validation helper that can
otherwise false-green; focused tests needed to keep their exact behavior; one
focused test-only implementation commit; clean automated measurement; and the
normal separate continuity commit. It does not authorize production changes,
Q1.7 remediation, a manual-evidence document or receipt, audit/parser/scanner
changes, `.quality/inventory` or baseline changes, executable P5 harnesses,
T1-T10 changes, another seam, acceptance changes, P6-P8 work, or a Q1.9 PASS
claim unless the complete classified population supports it.

# Measurements At Start

The evidence review measured clean commit
`18beee0a806bcc962cb2f7fbdd5952e43a2cb633`, exact parent product commit
`c382caa38be167fe17f847370ad8a12270644de3`, and commit tree
`85a6c67f38b8f764b849dcfac4da9eec071d22d0`. The clean no-evidence scorecard
SHA-256 is
`37ef9bd20cfddf450f381bbc1afe5bd8ea5a037f5a8b0f50713572d22e613f2b`;
the initial and final copies are byte-identical. It records 453 tests across all
27 packages, L0 8/8, L1 six PASS and Q1.6/Q1.7/Q1.9 UNMEASURABLE, six improved
ratchets, two held, zero regressed, zero non-comparable, and zero dirty paths.
The audit exits 1 for 13 documented non-passing criteria, never 2.

All 169 test ranges were enumerated. The complete local `tests` table class is
64 loops: 13 guarded and 51 unguarded. Confirmed non-table gaps include the
entrypoint collection, plugin-diagrams call expectations, archive-entry
validation, HTTP header expectations, cloud/template walk-error tables,
find/grep walk-error tables, Maven version order, Spring ordered-log helper,
filtered-template metadata and walk errors, and tips metadata. The remaining
105 ranges still require explicit item-by-item classification; the confirmed
list is not permission to stop early.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Confirm the new
continuity HEAD and exact ancestry before editing. Make one focused test-only
implementation commit, then the separate `docs: prepare next agent session`
continuity commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

An in-scope assertion must be executable and must fail if its intended
population is empty. Prefer a direct pre-loop `len(collection) == 0` fatal
assertion with a subject-specific message. For helpers, reject an empty caller
population at the helper boundary when non-empty input is part of that helper's
contract. Preserve subtest names, table order, exact values, expected side
effects, process-state restoration, public behavior, and Go 1.18 compatibility.
Do not add meaningless guards to fixture construction/delivery loops or to
result loops whose empty state already fails an independent executable
cardinality assertion. Record every exclusion and the assertion that makes it
safe.

Keep every generated report outside the worktree. Set
`API_COMPAT_REPORT_OUT` and `CLI_COMPAT_REPORT_OUT` to external paths whenever
running compatibility. The measured-tree identity counts ignored, untracked,
and tracked bytes, so confirm both ordinary and ignored status before every
authoritative audit.

# Required Reading

Before editing, confirm branch, HEAD, clean and ignored status, reciprocal
archive links, launcher `--check`, exact product ancestry, and the authorized
checkpoint block. Read the rolling handover, this archive, the complete P4
entry and checkpoint gate, both design documents, `.quality/README.md`, the
Q1.9 criterion and schema-2 validation path, the complete range reports or
regenerate them from source, all 169 enclosing tests/helpers, all existing
empty-population assertions, the complete audit wrapper, T15 repair, API/CLI
and launcher contracts, and the Make and acceptance gates. Do not classify a
range by name or search hit alone; inspect its executable control flow.

# Three Moves

1. Regenerate the exact 169-range inventory from the clean current tree. State
   the inclusion and exclusion rules before editing. Classify every range as
   fixture/support, independently guarded, or requiring a new executable
   assertion. Add guards to all 51 confirmed `tests` tables and every other
   in-scope unguarded population. Keep a complete path, function, collection,
   classification, and guard/exclusion record outside the worktree.
2. Run the exact affected test population, all named empty-population tests,
   complete uncached tests, race, and vet. Re-enumerate every range and inspect
   each new guard in executable context. If any intended population can still
   be empty without failure, keep Q1.9 blocked. Do not create manual evidence
   because Q1.7 remains blocked even if Q1.9 is repaired.
3. Run API/CLI and subprocess compatibility with external report paths; pinned
   lint; launcher and Make contracts; complete preflight; all four host
   acceptance flows; the 15-control audit meta-suite; empty-HOME count-2; and a
   clean no-evidence audit. Require audit exit never 2 and zero dirty paths.
   Record the exact new population, product commit, scorecard, and next P4
   blocker in continuity.

# Automatic Handoff

Before this session ends, finish and commit the coherent Q1.9 test-only move or
record an exact resumable blocker. Rewrite the rolling handover, record the P4
result in the roadmap, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent P4 move, replace only the launcher's
mutable regions, run launcher and handoff contracts, and make the separate
continuity-only commit. Do not launch a real successor, push, merge, publish,
distribute, stash, revert, or remove the worktree. P4 remains active for Q1.7,
and COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
