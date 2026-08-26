# Quality Upgrade Handover

Generated: 2026-08-26T20:41:47+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Unchanged production commit: `c382caa38be167fe17f847370ad8a12270644de3`.
- Q1.9 test-only implementation commit:
  `3a8f0bcd030baf787a29440ee8e4e4a087edb33b`.
- Implementation tree: `75dee3dcc1cb2fd91f7eb359e96b31cc4be3d986`.
- Exact ancestry is `3a8f0bc` -> continuity `fc454225` -> evidence review
  `18beee0a` -> production `c382caa`.
- After launch, obtain the new continuity HEAD with
  `git rev-parse --short=12 HEAD`; its exact parent must be `3a8f0bc`.
- No production, audit, scanner, parser, inventory, baseline, acceptance,
  mutation-harness, or manual-evidence file changed in the Q1.9 move.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active for Q1.7 and P5-P8 are
queued. Its active archive is
`docs/plan/agent-sessions/2026-08-26T204147+0200-repair-q17-partial-failures.md`.
The Q1.9 archive is answered history and links reciprocally to the new tail.
The graph has exactly one NEXT archive. The launcher's stable execution region
is unchanged; all 62 launcher controls must remain green.

## Q1.9 Result

The exhaustive scope is every syntactic `go/ast.RangeStmt` in every current
repository `_test.go` file, including helpers and nested function literals.
Assertion-driving tables, expectations, recorded collections, and structural
searches require an executable empty-population failure. Fixture construction
or delivery is excluded only after inspecting executable control flow. Result
loops are excluded only when a separate executable cardinality assertion
rejects unintended emptiness. A comment, a non-empty literal, a passing sample,
or one of the pre-existing 69 empty-population contracts is never treated as a
guard for another population.

The final AST inventory still contains exactly 169 range statements. Complete
item-by-item classification records:

- 72 range sites that required and now have a new executable assertion;
- 83 independently guarded range sites; and
- 14 fixture/support range sites with an explicit executable-flow exclusion.

All 64 ranges over an identifier named `tests` now have a direct executable
non-empty guard: the 13 existing guarded tables remain and all 51 confirmed
blockers were repaired. The additional 21 repaired range sites cover the
entrypoint caller and helper populations, plugin-diagrams call expectations,
archive entries, filesystem names, HTTP headers, cloud template names and
metadata, context paths, cloud/file/grep/template walk errors, Kibana results,
Maven version order, the Spring ordered-log helper, and tips identities and
metadata. Helper boundaries reject empty input where non-empty input is part
of the helper contract.

The 14 exclusions are archive and zip fixture construction, six recording
callback-delivery loops, context log-field copying, environment restoration,
render byte joining after `assertedWrites`, Kibana timestamp fixture and
serialization support, and an optional closed-reader negative identity loop.
Each row's exact independent assertion or support rationale is recorded
outside the worktree.

External records:

- final 169-range inventory:
  `/private/tmp/ply-p4-q19-final-inventory.txt`, SHA-256
  `05b2fc916c1c9d94a44322fcbb997d3e142a9ee70fe11b40de0b43e84a88b343`;
- complete path/function/collection/classification/guard record:
  `/private/tmp/ply-p4-q19-classification.tsv`, SHA-256
  `5f8cb9f00dec919769fde376ae240ae9d5a281063f2a48508d7439f9786eb005`;
- rerun of all 69 named empty-population contracts:
  `/private/tmp/ply-p4-q19-named-empty-tests.json`, SHA-256
  `bc8b171c0de3e4b2aa356ba8bf12033835b25b6c881fa30548ca1342547fdb09`.

This complete population supports Q1.9. The clean scorecard still records it
as UNMEASURABLE because the mission correctly did not create a schema-2 manual
receipt while Q1.7 remains blocked. Q1.6 remains independently supported by
its complete 57-struct, 74-field production review. No manual evidence was
created and no manual-evidence audit was run.

## Remaining Q1.7 Blockers

Q1.7 includes a production collection-item failure path when it reports the
failure and continues, and an operation when it can return a populated partial
aggregate with an error. It excludes fail-fast errors before partial state,
ignored errors without a reported partial result, non-error business warnings,
sequential non-collection phases, and test support.

The complete population remains 17 operations. Nine already have exact
failure-content coverage. These eight are the remaining blockers:

- `cmd/build_init.go`: `initCmd` project-loop missing-type, initialization,
  configuration-write, and POM-write failures;
- `pkg/bitbucket/bitbucket.go`: `Bitbucket.SynchronizeAllRepos` clone/pull
  warning content;
- `pkg/spring/io.go`: `DeleteDemoFiles` discovery and deletion failures;
- `pkg/maven/dependency.go`: `Repository.upgradeDependencies` version,
  maximum-version, and upgrade failures;
- `pkg/maven/plugin.go`: `Repository.upgradePluginsOnModel` per-plugin
  failures;
- `pkg/maven/deprecated.go`: `maven.RemoveDeprecated` partial result and
  warn-and-continue paths;
- `pkg/template/template.go`: `template.MergeTemplates` merge warning path;
- `pkg/config/cloud.go`: `GitCloudConfig.ValidTemplatesFrom` ordered partial
  result plus error.

The existing 23 focused top-level partial-failure contracts cover only the
other nine operations and cannot stand in for these eight. The next move adds
test-only, exact content and continuation/partial-result assertions for all
eight where current seams permit them. If a blocker cannot be exercised
without a production change, record its exact controllability gap and keep
Q1.7 blocked; do not add a seam or change production under that mission.

## Authoritative Measurement

The clean no-evidence audit ran from exact implementation commit `3a8f0bc`
with all generated output outside the worktree. Its identity is:

- commit: `3a8f0bcd030baf787a29440ee8e4e4a087edb33b`;
- commit tree: `75dee3dcc1cb2fd91f7eb359e96b31cc4be3d986`;
- status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`;
- inventory SHA-256:
  `4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`;
- scorecard SHA-256:
  `08e5ff9d450a6fa3b8824d82c720d339cf0d44003a490d01d3a412681081a615`;
- output root: `/private/tmp/ply-q19-audit.Y2Yoyl`.

The instrument remains wrapper version 1 with wrapper SHA-256
`b896bb8fac4c81fd5c4aab40732b3064e8355fa0d320dbc4b797541192beda75`,
parser SHA-256
`f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a`,
call-scanner SHA-256
`73dba9d563fb0b3e85af032487d024c7f42b80fd76de9f928d2888fed7347879`,
and Go `go1.26.2 darwin/arm64`; the test-file population SHA-256 is
`d4ac119c13efbf2945703781e582ad4523287956337901c3394b42eaa4ff62b3`.

The scorecard records 453 test functions, all 27 packages tested, 164 Go files,
92 table-driven tests, 13 scripts, five adapters, eight seams, eight subjects,
six non-executable mutation drivers, and four acceptance scripts. It exits 1
for 13 documented non-passing criteria, never 2. L0 is 8/8; L1 is six PASS and
Q1.6/Q1.7/Q1.9 UNMEASURABLE; ratchets are six improved, two held, zero
regressed, and zero non-comparable; dirty paths are empty.

## Gate Result

The exact affected 16-package population and all 69 named empty-population
contracts pass. Complete uncached tests across all 27 packages, race, and vet
pass. API/CLI and CLI-surface compatibility, fresh subprocess meta-contracts,
pinned golangci-lint 2.12.2 with zero issues, `make preflight`, `make test`,
Make and production-script contracts, all 62 launcher controls, all four host
acceptance flows, the 15-control audit meta-suite including T15, and the exact
empty-HOME count-2 command pass. API, CLI, audit, Go, temporary, and lint
outputs and caches were kept outside the worktree. The authoritative ordinary
and ignored statuses are empty.

Initial compatibility/lint attempts that lacked writable temporary/cache
directories were runner-invalid; unchanged reruns with isolated external paths
passed. An initial combined CLI shell command masked the compatibility status;
the exact CLI compatibility command was rerun independently and passed.

## Next Objective

Repair Q1.7 only with focused test contracts for the eight named gaps. Assert
the exact failure content and the observable continuation or ordered partial
result for every exercisable per-item failure path. Preserve production code,
public behavior, process-state restoration, current order, and Go 1.18
compatibility. Re-run the complete 17-operation classification and exact
focused population before concluding whether Q1.7 is supported.

Do not create manual evidence in that move. Q1.6 and Q1.9 support can be
combined with Q1.7 only in a later, separately authorized evidence checkpoint.
P4 stays active and P5-P8 stay queued.

## Start And Stop

Read this handover, the linked NEXT archive, complete P4 entry and gate, both
design documents, Q1.7 criterion, all 17 operations and their tests, and the
relevant logger/double helpers before editing. Confirm branch, HEAD, clean and
ignored status, reciprocal archive links, launcher `--check`, and exact
ancestry.

Make one focused test-only implementation commit and one separate continuity
commit. Stop before production changes, new seams, manual evidence, Q1.9
changes, audit/inventory/baseline changes, P5 mutation work, T1-T10 changes,
P6-P8, publication, or distribution. Do not push, merge, stash, revert, launch
a successor, or remove the worktree.
