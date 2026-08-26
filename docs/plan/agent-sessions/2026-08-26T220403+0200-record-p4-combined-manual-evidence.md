# Agent Session: Record P4 Combined Manual Evidence

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T220403+0200-record-p4-combined-manual-evidence`
Created: `2026-08-26T22:04:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `dd2d52e4513e02a362c1b51ebd895913cbcb8b5d953499015b83761bae8e245f`
Previous: [2026-08-26T213324+0200-repair-q17-delete-demo-discovery.md](2026-08-26T213324+0200-repair-q17-delete-demo-discovery.md)
Next: [2026-08-26T224925+0200-build-p5-cli-context-harness.md](2026-08-26T224925+0200-build-p5-cli-context-harness.md)
Outcome: P4 complete: external schema-2 Q1.6/Q1.7/Q1.9 evidence is valid, focused audit exits 0, full audit records all nine L1 rows PASS with zero dirty paths, and every P4 gate passes.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Finish P4 only if exhaustive current-project review still truthfully supports
the three remaining manual L1 rows. Record exact non-empty schema-2 receipts
for Q1.6 default dependency doubles, Q1.7 partial-failure content assertions,
and Q1.9 executable empty-population assertions, then prove them with clean
focused and full authoritative audits. Do not infer PASS from a sample or
manufacture evidence.

# Authorized Roadmap

P3 is complete. P4 is active only for combined Q1.6/Q1.7/Q1.9 manual evidence;
P5-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`. Q1.6 is independently supported by its
complete 57-struct, 74-field review, Q1.7 by its complete 17-operation
executable-flow classification, and Q1.9 by its complete 169-range
classification. The launcher must remain NEXT until P4 and P5-P8 finish.

This evidence-only mission authorizes exhaustive read-only source and contract
review; external temporary scripts, reports, and one commit-bound schema-2
manual-evidence document; exact Q1.6/Q1.7/Q1.9 receipt validation; clean
focused and full audits using that explicit external document; and the normal
single continuity-only commit. It does not authorize production or test
changes, a checked-in or worktree-local evidence file, another seam, exported
API changes, validator or scanner changes, `.quality/inventory` or
baseline changes, executable P5 harnesses, T1-T10 changes, acceptance changes,
P6-P8 work, or a PASS claim that a complete reviewed population does not
support.

# Measurements At Start

Focused Q1.7 implementation commit
`4f1f45f3f758f2f09dd7d20967efe8ef74a0c613` has exact tree
`7326302c171604dc315e2f79cd6871c55935aa60`. It changes only
`pkg/spring/io.go` and
`pkg/spring/delete_demo_files_partial_failure_test.go`. Public
`DeleteDemoFiles(string, config.ProjectConfiguration)` delegates to a
private helper accepting only the discovery function and passes
`file.FindFirst` in production. One exact test injects the discovery
error, proves the `.kt` and `src/test/kotlin` request, asserts
`Unable to find testfile, fileSuffix=.kt`, and proves continuation to
delete `HELP.md`, `mvnw`, and `mvnw.cmd`.

The complete Q1.7 population remains 17 operations and every included path now
has exact failure content plus continuation or exact ordered partial-result
coverage. The final classification SHA-256 is
`01cb1abc6293f780fd0caffd6a37da7398ea6d49e9b6b7589ae15ff23d83b9eb`;
the exact 34-contract population manifest SHA-256 is
`0772d019160fe45a8eeefce5fcadba7821fc18da896374d19cb2b48a1dbb05e3`.
Q1.6's complete population remains 57 dependency structs and 74 fields.
Q1.9's complete population remains 169 syntactic test ranges.

All focused Q1.7 contracts, complete uncached tests, race, vet, API/CLI and
subprocess compatibility, pinned lint, launcher and Make contracts, complete
preflight, all four host acceptance flows, the 15-control audit meta-suite, and
empty-HOME count-2 pass. The no-evidence audit scorecard SHA-256 is
`ce4190dc03b7bc37aa285569f737b2a59a38c15c5e99db125e340c64e6301b5e`.
It records 464 tests across all 27 packages, L0 8/8, final wrapper L1 six PASS
and Q1.6/Q1.7/Q1.9 UNMEASURABLE, six improved ratchets, two held, zero
regressed, zero non-comparable, and zero dirty paths. It exits 1 for 13
documented non-passing criteria, never 2. No evidence document or receipt was
created in the repair session.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Confirm
the new continuity HEAD and exact ancestry before measuring; its exact parent
must be `4f1f45f`. This checkpoint must not create a product
implementation commit. The receipts bind the exact clean current continuity
HEAD and are supplied from outside the repository. After measurement, make
only the normal `docs: prepare next agent session` continuity commit.

Keep the evidence document and every generated report or helper outside the
worktree; never create `.quality/manual-evidence.json` or
`.agent-task/current.md`. The measured tree counts tracked, untracked,
and ignored bytes. Confirm ordinary and ignored status before every
authoritative audit. Set `API_COMPAT_REPORT_OUT` and
`CLI_COMPAT_REPORT_OUT` to external paths whenever running
compatibility.

Use schema version 2 and include only Q1.6, Q1.7, and Q1.9 receipts. Derive the
exact current commit, commit tree, clean status digest, inventory object, and
complete instrument object from a clean no-evidence audit of current HEAD.
Compute every `evidence_sha256` from canonical UTF-8 JSON with sorted
keys, no whitespace, and `ensure_ascii=False`. Pass the evidence
explicitly with `--manual-evidence`.

For Q1.6, enumerate the complete in-scope dependency-struct population and
prove every dependency has a safe default and argument recorder and each
struct is passed whole rather than reconstructed field by field. For Q1.7,
enumerate all production operations that report collection-item failures and
continue or return a populated partial aggregate with an error, and prove
every included failure has exact content and continuation or ordered partial
result. For Q1.9, enumerate the complete in-scope test collection-iteration
population and prove every required collection has an executable assertion
that fails when empty. State exact inclusion and exclusion rules before
counting and run the named contracts behind every claim.

If any denominator is empty or any item fails its criterion, do not emit a
false receipt. Record the exact blocker, keep the affected row UNMEASURABLE,
leave P4 active, and prepare the smallest truthful successor task. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, launch a successor, or run destructive Git commands.

# Required Reading

Before measuring, confirm branch, HEAD, clean and ignored status, reciprocal
archive links, launcher `--check`, exact ancestry, and the authorized
checkpoint block. Read the rolling handover, this archive, the complete P4
entry and gate, both design documents, `.quality/README.md`,
`.quality/inventory`, the complete audit wrapper, complete schema-2
parser and manual-evidence negative meta-test, baseline manual example and
reproduction README, T15 repair, all 57 dependency structs and 74 fields with
their default/recording contracts, all 17 Q1.7 production operations and their
exact 34-contract population, and all 169 Q1.9 ranges with their executable
guards and exclusions. Read API/CLI, launcher, Make, and acceptance contracts
needed by the final gate. Do not rely on names or text search alone where Go
structure or runtime execution establishes the population.

# Three Moves

1. Run a clean no-evidence audit to capture exact current repository, tree,
   status, inventory, instrument, denominator, and automated criterion
   identities. Regenerate explicit exhaustive scope and records for Q1.6,
   Q1.7, and Q1.9; inspect every subject in executable context; run exact
   non-empty focused contract populations; and record counts, exclusions, and
   hashes outside the worktree.
2. Only if all three complete populations satisfy their required truth, create
   one external schema-2 document bound to clean current HEAD. Validate its
   canonical receipt digests, run focused Q1.6/Q1.7/Q1.9 audit, and run the
   full authoritative audit with the same explicit document. Require all three
   rows PASS, all automated results held, audit exit never 2, and zero dirty
   paths. Record document, receipt, record, and scorecard SHA-256 identities.
3. Run API/CLI and subprocess compatibility with external report paths; pinned
   lint; launcher and Make contracts; complete tests, race, and vet; complete
   preflight; all four host acceptance flows; the 15-control audit meta-suite;
   and empty-HOME count-2 in proportion to the evidence checkpoint. Mark P4
   complete and hand off the first bounded P5 move only if all nine L1 rows and
   every P4 exit are valid and non-empty.

# Automatic Handoff

Before this session ends, finish the coherent combined evidence checkpoint or
record an exact resumable blocker. Rewrite the rolling handover, record the P4
result in the roadmap, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace
only the launcher's mutable regions, run launcher and handoff contracts, and
make the single continuity-only `docs: prepare next agent session`
commit. If P4 completes, the next archive must begin the first narrowly bounded
P5 mutation-evidence subject; otherwise it must name the exact P4 blocker. Do
not launch a real successor, push, merge, publish, distribute, stash, revert,
or remove the worktree. COMPLETE remains invalid while P5-P8 are unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
