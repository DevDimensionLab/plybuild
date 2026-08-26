# Quality Upgrade Handover

Generated: 2026-08-26T22:04:03+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Q1.7 discovery-seam implementation commit:
  `4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`.
- Implementation tree: `7326302c171604dc315e2f79cd6871c55935aa60`.
- Its exact parent is continuity commit
  `58f83a9661ba6b9f1d72766db23982cbaf0e2ee0`; that commit's exact
  parent is Q1.7 characterization commit
  `7be93e7e822f05bdb007e4e5a1a00402e620a6cb`.
- After launch, obtain the new continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `4f1f45f`.
- The implementation commit changes only `pkg/spring/io.go` and
  `pkg/spring/delete_demo_files_partial_failure_test.go`.
- No exported API, dependency struct, deletion injection, audit, scanner,
  parser, inventory, baseline, acceptance, mutation harness, or manual evidence
  changed.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active for combined
schema-2 Q1.6/Q1.7/Q1.9 manual evidence, while P5-P8 remain queued. Its active
archive is
`docs/plan/agent-sessions/2026-08-26T220403+0200-record-p4-combined-manual-evidence.md`.
The discovery-repair archive is answered history and links reciprocally to the
new tail. The graph has exactly one NEXT archive. The launcher's stable
execution region is unchanged; all 62 launcher controls must remain green.

The tracked launcher/archive apparatus is the task source. Do not create
`.agent-task/current.md`: the audit counts ignored and untracked bytes,
and the established continuity contract requires both ordinary and ignored
status to remain empty.

## Q1.7 Result

The inclusion rule remains: include a production collection-item failure path
when it reports the failure and continues, and include an operation when it can
return a populated partial aggregate with an error. Exclude fail-fast errors
before partial state, intentionally ignored errors with no reported partial
result, non-error business warnings, sequential non-collection phases, and
test support. Every classification was regenerated from executable control
flow, not inferred from names, search hits, or exit status.

The complete population remains exactly 17 operations. All included paths now
have exact failure-content plus continuation coverage or exact ordered
partial-result coverage. The last repaired branch is
`spring.DeleteDemoFiles` discovery failure:

- public `DeleteDemoFiles(string, config.ProjectConfiguration)` keeps its
  exact signature and delegates with production `file.FindFirst`;
- private `deleteDemoFiles` accepts only
  `func(string, string) (string, error)`;
- the focused test injects a sentinel discovery error and proves one exact
  request with suffix `.kt` and path `src/test/kotlin` below the
  target;
- it asserts exact warning
  `Unable to find testfile, fileSuffix=.kt`;
- it proves continuation deletes the later fixed items `HELP.md`,
  `mvnw`, and `mvnw.cmd`.

Suffix selection, lookup arguments, existing warning spelling, deletion
behavior, fixed-item ordering, side effects, and Go 1.18 compatibility are
preserved. Deletion behavior was not injected.

External Q1.7 records:

- starting 17-row record:
  `/private/tmp/ply-p4-q17-classification-start.tsv`, SHA-256
  `4b247f3343a651374b96245d6c15b563f573a152799f3329713c578ace863e57`;
- final 17-row record:
  `/private/tmp/ply-p4-q17-classification-final.tsv`, SHA-256
  `01cb1abc6293f780fd0caffd6a37da7398ea6d49e9b6b7589ae15ff23d83b9eb`;
- exact 34-contract population manifest:
  `/private/tmp/ply-p4-q17-focused-population-final.txt`, SHA-256
  `0772d019160fe45a8eeefce5fcadba7821fc18da896374d19cb2b48a1dbb05e3`;
- complete focused execution log:
  `/private/tmp/ply-p4-q17-focused-population-final.log`, SHA-256
  `ee0a0abc8d06aa7014c96e238e1b1a019615c1a7a2e5ba3dd631e6b6f7834ec5`.

Q1.6 remains independently supported by its complete 57-struct, 74-field
production review. Q1.9 remains independently supported by its complete
169-range classification. No schema-2 evidence document or receipt was
created, so all three rows remain formally UNMEASURABLE until the next
evidence-only checkpoint.

## Authoritative Measurement

The clean no-evidence audit ran from exact implementation commit `4f1f45f`
with all output outside the worktree. Its identity is:

- commit: `4f1f45f3f758f2f09dd7d20967efe8ef74a0c613`;
- commit tree: `7326302c171604dc315e2f79cd6871c55935aa60`;
- status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`;
- inventory SHA-256:
  `4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`;
- scorecard SHA-256:
  `ce4190dc03b7bc37aa285569f737b2a59a38c15c5e99db125e340c64e6301b5e`;
- raw upstream scorecard SHA-256:
  `41d31c3df35b23c8950a14963681fec89c5b0d53ef72a066f6d621342ba800ba`;
- output root: `/private/tmp/ply-q17-clean-audit.dPEiK2`.

The instrument is wrapper version 1 with wrapper SHA-256
`b896bb8fac4c81fd5c4aab40732b3064e8355fa0d320dbc4b797541192beda75`,
parser SHA-256
`f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a`,
call-scanner SHA-256
`73dba9d563fb0b3e85af032487d024c7f42b80fd76de9f928d2888fed7347879`,
upstream SHA-256
`45787c231a0f07f27093824b19685aee4fb4cab5159cdb7864834b79fbc9cee9`,
and Go `go1.26.2 darwin/arm64`. The production population is 93 files,
SHA-256
`d2eebd6b3356e30d14feb6ab43e0e2cc276b14e1354cecbca481aa9a48ffa392`;
the test population is 75 files, SHA-256
`1856f4bdbead4fc1c2c17a65ec8b7f3d084a0506c0ca178dace03c121fc17d7f`.

The scorecard records 464 test functions, all 27 packages tested, 169 Go files,
95 table-driven tests, 13 scripts, five adapters, eight seams, eight subjects,
six non-executable mutation drivers, and four acceptance scripts. It exits 1
for 13 documented non-passing criteria, never 2. L0 is 8/8; final wrapper L1
is six PASS with Q1.6/Q1.7/Q1.9 UNMEASURABLE; ratchets are six improved, two
held, zero regressed, and zero non-comparable; dirty paths are empty. The
upstream prose still prints its pre-wrapper Q1.3 failure, but the authoritative
structured wrapper result correctly records exact local Q1.3 as PASS at 0 of
27.

## Gate Result

Focused Spring contracts and the complete 34-contract Q1.7 population pass.
Complete uncached tests across all 27 packages, race, vet, API/CLI and
CLI-surface compatibility, fresh subprocess contracts, pinned golangci-lint
2.12.2 with zero issues, `make preflight`, `make test`, Make and
production-script contracts, all four host acceptance flows, the 15-control
audit meta-suite including T15, and exact empty-HOME count-2 all pass. API and
CLI reports were external and have SHA-256 values
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and
`955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

The default Go and lint cache locations were not writable in this sandbox, so
unchanged commands were rerun with external caches. Two initial standalone
launcher invocations hit the known partial-raw-log signal-fixture timing miss;
the unchanged third standalone run, complete preflight, and complete
`make test` each passed all 62 controls. Ordinary and ignored status were
empty immediately before and after the authoritative audit.

## Next Objective

Run one evidence-only combined P4 checkpoint for Q1.6, Q1.7, and Q1.9. Start
from the clean continuity HEAD whose exact parent is `4f1f45f`.
Regenerate and inspect the complete 57-struct/74-field Q1.6 population, all 17
Q1.7 operations, and all 169 Q1.9 ranges from executable structure. Run the
exact non-empty focused populations behind every claim.

Only if all three complete populations still support their criteria, create
one external schema-2 evidence document bound to the exact clean current HEAD,
inventory, status, and complete instrument identity. Include only Q1.6, Q1.7,
and Q1.9 receipts; compute canonical receipt digests; pass the document
explicitly with `--manual-evidence`; and require focused and full audits
to make all three rows PASS with exit never 2 and zero dirty paths. Do not
create `.quality/manual-evidence.json` or any other worktree-local
evidence file.

If the combined evidence is valid, complete P4 and hand off the first bounded
P5 mutation-evidence move. If any claim fails, create no receipt for that
claim, leave P4 active, and hand off the exact smallest truthful blocker. P5-P8
remain queued until P4 is complete.

## Start And Stop

Read this handover, the linked NEXT archive, complete P4 entry and gate, both
design documents, `.quality/README.md`, schema-2 parser and negative
meta-tests, baseline manual example, T15, and every production/test population
behind the three receipts before measuring. Confirm branch, HEAD, clean and
ignored status, reciprocal archive links, launcher `--check`, exact
ancestry, and the authorized checkpoint block.

This next checkpoint is read-only except for an external temporary evidence
document, external reports, and the final tracked continuity records. Stop
before production or test changes, audit/parser/scanner/inventory/baseline
changes, executable P5 harnesses, T1-T10 changes, acceptance changes,
publication, or distribution. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
