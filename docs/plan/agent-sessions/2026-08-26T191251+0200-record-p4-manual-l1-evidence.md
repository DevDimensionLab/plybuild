# Agent Session: Record P4 Manual L1 Evidence

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T191251+0200-record-p4-manual-l1-evidence`
Created: `2026-08-26T19:12:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a37605edc2c426e6ea0984c0f2436977c0249158d020523fd56c262a699dade0`
Previous: [2026-08-26T184054+0200-close-loopback-seam.md](2026-08-26T184054+0200-close-loopback-seam.md)
Next: [2026-08-26T195729+0200-repair-q19-empty-populations.md](2026-08-26T195729+0200-repair-q19-empty-populations.md)
Outcome: The exhaustive evidence review created no schema-2 document because
Q1.7 has eight uncovered operations among 17 in scope and Q1.9 has 51
unguarded local-table loops among 64, plus confirmed non-table gaps among 169
total test ranges. Q1.6's complete 57-struct, 74-field population is supported
independently, but all three rows remain UNMEASURABLE. Initial and final clean
no-evidence audits are byte-identical at scorecard SHA-256
`37ef9bd20cfddf450f381bbc1afe5bd8ea5a037f5a8b0f50713572d22e613f2b`,
exit 1 for 13 documented findings, and contain zero dirty paths. P4 remains
active; the reciprocal NEXT move repairs Q1.9 only.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Finish P4 only if an exhaustive current-project review truthfully supports the
three remaining manual L1 rows. Record exact non-empty schema-2 receipts for
Q1.6 default dependency doubles, Q1.7 partial-failure content assertions, and
Q1.9 executable empty-population assertions, then prove them with the clean
authoritative audit. Do not infer PASS from a sample or manufacture evidence.

# Authorized Roadmap

P3 is complete, P4 is active only because Q1.6, Q1.7, and Q1.9 have no
current-project receipts, and P5-P8 remain queued in the machine-readable block
in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until every
authorized checkpoint finishes. P4.10 product commit `c382caa` changed only the
private interactive server address to `127.0.0.1:7999`, added four focused
contracts and the eighth regular non-executable seam driver, and improved Q1.4
from 7/8 to 8/8. Every automated P4 exit now passes.

This evidence-only mission authorizes exhaustive read-only source and contract
review; external temporary scripts, reports, and one commit-bound schema-2
manual-evidence document; exact Q1.6/Q1.7/Q1.9 receipt validation; clean focused
and full audits using that explicit external document; and tracked continuity
records of the commands, populations, exclusions, canonical receipt digests,
evidence-document hash, scorecard identity, and result. It does not authorize
production or test changes, a checked-in or worktree-local manual-evidence file,
validator or scanner changes, `.quality/inventory` or baseline changes, an
executable P5 mutation harness, T1-T10 work, another seam, acceptance changes,
P6-P8 work, or a PASS claim that the complete reviewed population does not
support.

# Measurements At Start

Clean product commit `c382caa` has 453 tests across all 27 packages. Q0.6 has
29 guarded safe-writer sites, 23 write and 6 copy, zero skipped tests, and zero
unsafe direct test writes. Q0.8 is reciprocal for all 13 production scripts.
Q1.1 is 0/27, Q1.2 is zero, exact Q1.3 is 0/27 with all five declared adapters
valid, Q1.4 is 8/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases across 94
Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, L1 six PASS and three UNMEASURABLE, six improved
ratchets, two held, zero regressed, zero non-comparable, and zero dirty paths.
Its scorecard SHA-256 is
`b4c7ddc7d163a833663f8b7cd876e94aadf4c1a875b945252a46441e354c6b72`;
the measured product tree is `6e44c014c73fb8c99c6c90f20a64a41f314fce43`.
The repaired T15 proof still passes all 15 controls and reproduces the exact
stored debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. This checkpoint
must not create a product implementation commit: the receipts bind the exact
clean current continuity commit and are supplied from outside the repository.
After the measurement, make only the normal continuity commit that records the
result and prepares the next authorized session. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

The measured-tree implementation counts tracked, untracked, and ignored
repository bytes. Keep the evidence document and every generated report or
helper outside the worktree; never create `.quality/manual-evidence.json`.
Derive the exact current commit, commit tree, clean status digest, inventory
object, and complete instrument object from a clean audit of the current HEAD.
Use schema version 2, include only Q1.6/Q1.7/Q1.9 receipts, and compute each
`evidence_sha256` from canonical UTF-8 JSON with sorted keys, no whitespace, and
`ensure_ascii=False`. Pass the document explicitly with `--manual-evidence`.

For Q1.6, enumerate the complete in-scope dependency-struct population and
prove every dependency has a default and argument recorder and that each struct
is passed whole rather than reconstructed field by field. For Q1.7, enumerate
all commands or operations that can report per-item or partial failures and
prove every such failure has a content assertion and none is exit-only. For
Q1.9, enumerate the complete in-scope test collection-iteration population and
prove every iterated collection has an executable assertion that fails when its
population is empty. State exact inclusion and exclusion rules before counting;
run the named contracts behind every claim. If any denominator is empty or any
item fails its required truth, record the exact blocker, keep that row
UNMEASURABLE, and leave P4 active.

# Required Reading

Before measuring, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product parent `c382caa`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/README.md`, `.quality/inventory`, complete audit wrapper,
complete schema-2 parser and manual-evidence negative meta-test, baseline manual
example and reproduction README, T15 repair, all dependency structs and their
default/recording contracts, all command and operation partial-failure paths and
their assertions, and every contract that iterates or validates a non-empty
recorded collection. Read the API/CLI and launcher contracts needed by the final
gate. Do not rely only on names or text search where Go structure or runtime
execution is needed to establish the population.

# Three Moves

1. Run a clean audit without manual evidence to capture the exact current
   repository, tree, status, inventory, instrument, denominator, and automated
   criterion identities. Define explicit exhaustive scope rules for Q1.6,
   Q1.7, and Q1.9; enumerate every subject; inspect its production and test
   structure; and run a non-empty exact contract population. Record exact
   counts, commands, exclusions, and any failing item.
2. Only if all three complete populations satisfy their required truth, create
   an external schema-2 document bound to the clean current HEAD, validate its
   canonical receipt digests, run focused Q1.6/Q1.7/Q1.9 audit, and then run the
   full authoritative audit with the same explicit document. Require all three
   rows to become PASS, every automated result to hold, audit exit never 2, and
   zero dirty paths. Record the document and scorecard SHA-256 identities. If a
   claim cannot be proved, do not change code in this mission; hand off the
   exact smallest truthful remediation instead.
3. Run API/CLI compatibility, launcher and Make contracts, complete tests,
   race, vet, pinned lint, all four host acceptance flows, the 15-control audit
   meta-suite, and empty-HOME count-2 in proportion to the evidence checkpoint.
   Confirm the worktree stayed byte-clean throughout measurement. Mark P4
   complete and advance the roadmap to the first P5 subject only if every P4
   exit, including all nine L1 rows, is valid and non-empty.

# Automatic Handoff

Before this agent session ends, finish the coherent evidence checkpoint or
record an exact resumable blocker. Rewrite the rolling handover, record the
measured P4 result, answer this archive, create exactly one reciprocally linked
NEXT archive for the next coherent authorized roadmap move, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
single continuity-only `docs: prepare next agent session` commit. Do not launch
a real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
