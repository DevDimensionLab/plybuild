# Quality Upgrade Handover

Generated: 2026-08-26T22:49:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P4 evidence commit: `88a95ad6effe7c2198d9505963dc19022bcabc10`.
- Its exact parent is focused Q1.7 implementation commit
  `4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`.
- Measured commit tree: `120d8db0b83a7724b0026f7c18a77a706e9f1cda`.
- After launch, obtain the new continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `88a95ad`.
- The P4 session made no production, test, audit, scanner, parser, inventory,
  baseline, acceptance, mutation-harness, or checked-in evidence change.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P5 is active while P6-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-26T224925+0200-build-p5-cli-context-harness.md`.
The P4 evidence archive is answered history and links reciprocally to this new
tail. The graph has exactly one NEXT archive. The launcher's stable execution
region is unchanged and all 62 controls pass in isolation and in complete
preflight.

The tracked launcher/archive apparatus is the task source. Do not create
`.agent-task/current.md`; the audit counts ignored and untracked bytes, and the
continuity contract requires both ordinary and ignored status to be empty at a
measured checkpoint.

## P4 Result

P4 is complete. Fresh current-project review supports every remaining manual
L1 row without sampling:

- Q1.6 includes every production unexported type ending `Dependencies` plus
  every exact exported `internal/adapter/*` `Dependencies` type. It excludes
  aliases, test-only structs, and production structs that do not compose
  injectable dependencies. The population is 57 structs and 74 fields: all 74
  have safe defaults and argument recorders, every struct is delivered whole,
  and 245 of 245 exact named contracts run and pass.
- Q1.7 includes production collection-item failures that report and continue,
  and operations returning a populated partial aggregate with an error. It
  excludes fail-fast pre-population errors, intentionally ignored errors with
  no reported partial result, business warnings, sequential non-collection
  phases, and test support. All 17 operations and 42 distinct included failure
  branches have exact content plus continuation or exact ordered partial-result
  assertions; 34 of 34 named contracts run and pass; none is exit-only.
- Q1.9 includes every assertion-driving syntactic `go/ast.RangeStmt` in every
  selected `_test.go`, including helpers and nested literals. It excludes only
  fixture construction/delivery and state-copy support, with result ranges
  excluded only when another executable assertion makes empty unsafe or empty
  is the topology under test. Fresh AST enumeration corrects the inherited
  claim: current HEAD has 183, not 169, total ranges. Exactly 168 require
  executable empty-population assertions and 15 are classified support
  exclusions. All 168 satisfy the criterion. The 69 named guard contracts and
  the 34 Q1.7 contracts covering the 14 later sites all run and pass.

External population identities:

- Q1.6 AST inventory:
  `38ff8cb2ad78f9ebf41fcfc5aeba2f9c5ddcf6a7dcf467f14f5fc1886a6345da`;
  review:
  `2dff55671c9b545a28de2e7ff7bcb614eb3064d01b2436ce58a20ce9a875f45b`;
  245-contract manifest:
  `6e045b462101ae8feab21ee6d6e117cb308fdaaec7a6615e70de215588914955`.
- Q1.7 classification:
  `01cb1abc6293f780fd0caffd6a37da7398ea6d49e9b6b7589ae15ff23d83b9eb`;
  34-contract manifest:
  `0772d019160fe45a8eeefce5fcadba7821fc18da896374d19cb2b48a1dbb05e3`.
- Q1.9 AST inventory:
  `417d9add0e20ae9791a5d5a106ccb5dc1accdccacc9d99355092b25a2cf10787`;
  classification:
  `8326a5ef3037f9cdd72e16a470ead8c2d4e4de96e960215af8002dce5ccc5c26`;
  168-range required-population record:
  `107d5953139a2959a0777047183893e894ff5743480979c04e2cb6e9bd5d08b7`;
  69-contract manifest:
  `bfd7210785bc672192575e372b098ebb1cbf226b5793d0ac8c15065be0e6295b`.

## Authoritative Measurement

All reports and the evidence document are under external root
`/private/tmp/ply-p4-evidence-audit.D7UNNg`. The no-evidence audit binds:

- commit: `88a95ad6effe7c2198d9505963dc19022bcabc10`;
- commit tree: `120d8db0b83a7724b0026f7c18a77a706e9f1cda`;
- clean status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`;
- inventory SHA-256:
  `4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`;
- no-evidence scorecard SHA-256:
  `0d038c41c160c8851b405e32f69bb9dcf03729749fca5f423ef18435b2cc58e2`.

The instrument is wrapper version 1 with wrapper SHA-256
`b896bb8fac4c81fd5c4aab40732b3064e8355fa0d320dbc4b797541192beda75`,
parser SHA-256
`f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a`,
call-scanner SHA-256
`73dba9d563fb0b3e85af032487d024c7f42b80fd76de9f928d2888fed7347879`,
upstream SHA-256
`45787c231a0f07f27093824b19685aee4fb4cab5159cdb7864834b79fbc9cee9`,
Q0.6 contract SHA-256
`7881988a831b02e367b21aa0d2ad3d08e25afb7524a5c30f8bfa2a2533c64772`,
and report-template SHA-256
`3138b45cd2ad3f3e396ad3c5ea1325130619d6d2c61c2d2c64db262035c0f07d`.
Go is `go1.26.2 darwin/arm64`; the selected production and test populations
remain 93 and 75 files.

The external schema-2 document includes only Q1.6, Q1.7, and Q1.9 criterion
receipts. Its SHA-256 is
`c97c0336998508fe410fc3f56b65cc41762281850ba5822d53c623dc74be34c0`.
Canonical evidence-object SHA-256 values are:

- Q1.6: `85db15d942b850f99baa518428043533ceee7b94d911b2de94b6e45a5f3f0ad7`;
- Q1.7: `e94221777b1f8ea0300cdc286fbb7086440d1038b5b9e9da0f7d845d18845e2d`;
- Q1.9: `6d4bf0c4429fa7a731078b84bc850d401a1c636f4cf72bb8d110dfe4700d9616`.

The focused audit exits 0 and its scorecard SHA-256 is
`63ede0de00d52d950f10845ee8d6e16c1e948e73c8dd183c368ef0e87cc4c729`.
The full audit exits 1, never 2, and its scorecard SHA-256 is
`52490a43d26c5b6c1e6031460352a9fbc6de31bba52128301d3d3e1f31d86bc6`.
It records 464 test functions across all 27 packages, L0 8/8, all nine L1
rows PASS, six improved ratchets, two held, zero regressed, zero
non-comparable, and zero dirty paths. Its nine remaining non-passing criteria
belong to P5-P8. Every non-manual criterion object is identical to the
no-evidence audit.

## Gate Result

API/CLI and CLI-surface compatibility, entry/subprocess contracts, pinned
golangci-lint 2.12.2 with zero issues, complete uncached tests across all 27
packages, race, vet, launcher and Make contracts, all four host acceptance
flows, the standalone 15-control audit meta-suite including T15, exact
empty-HOME count-2, and complete preflight all pass. API and CLI report
SHA-256 values are
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and
`955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

Two launcher invocations missed the known signal-fixture partial-log timing
assertion, once under concurrent heavy Go gates and once in the final
source-archive nested generation. Unchanged isolated reruns and complete
preflight pass all 62 controls. Ordinary and ignored status were empty before
both authoritative audits and after the full gate.

## Next Objective

Begin P5 with only the first `.quality/inventory` subject, `cli-context`, whose
production roots are `cmd` and `pkg/context` and whose declared harness is
`scripts/mutate-cli-context`. Create that one executable mutation harness and
its `scripts/test-mutate-cli-context` meta-test. Define at least eight
meaningful, deterministic source mutations limited to the subject roots; run
each mutation in a disposable external copy; prove the harness detects its own
falsification controls; implement the methodology T1-T10 controls; and report
exact declared, killed, survived, and unusable counts.

Do not implement another inventory subject in the same move. If a mutation
survives, classify the exact reachability, observability, or controllability
gap and make only the smallest in-scope test or seam improvement needed to kill
it. Require `declared == killed`, `survived == 0`, and `unusable == 0` before
claiming this subject complete. Run the focused harness/meta-test, automated
Q2.1-Q2.4 audit view, and the complete repository gate. P5 remains active after
this first subject; hand off the second bounded inventory subject only after a
clean measured checkpoint.

## Start And Stop

Read this handover, the linked NEXT archive, P5 and the checkpoint gate, both
design documents, `.quality/README.md`, `.quality/inventory`, the complete
mutation-related audit/parser/meta-test logic, the production roots `cmd` and
`pkg/context`, and the tests that can kill every proposed mutation before
editing. Confirm branch, HEAD, clean and ignored status, reciprocal archive
links, launcher `--check`, exact ancestry, and the authorized checkpoint block.

Stop before a second mutation subject, P6-P8, acceptance expansion, Go or
dependency upgrades, domain modernization, exported API or CLI changes,
audit/parser/scanner/baseline changes, publication, or distribution. Keep
generated mutant copies, reports, caches, and artifacts outside the worktree.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
