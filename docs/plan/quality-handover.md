# Quality Upgrade Handover

Generated: 2026-08-26T19:57:29+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Product commit: `c382caa38be167fe17f847370ad8a12270644de3`.
- Evidence-review commit before this continuity handoff:
  `18beee0a806bcc962cb2f7fbdd5952e43a2cb633`.
- Exact parent of `18beee0`: `c382caa38be167fe17f847370ad8a12270644de3`.
- After launch, obtain the new docs-only continuity HEAD with
  `git rev-parse --short=12 HEAD`.
- No product or test file changed in the evidence review. No manual-evidence
  file was created in or outside the repository.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued.
Its active archive is
`docs/plan/agent-sessions/2026-08-26T195729+0200-repair-q19-empty-populations.md`.
The evidence-review archive is answered history and links reciprocally to the
new tail. The graph has exactly one NEXT archive. The launcher stable skeleton
remains SHA-256
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`,
and all 62 launcher controls remain unchanged.

## Evidence Decision

P4 did not complete. Q1.7 and Q1.9 have real current-project blockers, so the
three schema-2 receipts were not created. Q1.6, Q1.7, and Q1.9 consequently
remain UNMEASURABLE in the authoritative scorecard. There is no evidence-file
SHA-256 and there are no canonical receipt digests to record. No
`--manual-evidence` audit was run; a covered subset was never promoted to PASS.

The explicit scope rules were:

- Q1.6 includes every production struct in `cmd`, `internal`, or `pkg` whose
  type name ends in `Dependencies` or `dependencies`, and every field in those
  structs. It excludes interfaces, selector functions, ordinary configuration
  structs, and test-only doubles. Each field must have safe default behavior,
  an argument-observing double, and whole-struct delivery.
- Q1.7 includes a production collection-item failure path when it reports the
  failure and continues, and an operation when it can return a populated
  partial aggregate with an error. It excludes fail-fast errors before any
  partial state, ignored errors with no reported partial result, non-error
  business warnings, sequential non-collection phases, and fixture/support
  code.
- Q1.9 begins with every syntactic Go `range` statement in a `_test.go` file.
  Assertion-driving tables, expectations, recorded collections, and structural
  searches require executable empty-population failure. Fixture construction
  or delivery loops are excluded, as are result loops whose emptiness already
  fails an independent executable assertion. Classification never treats a
  comment, a non-empty literal, or a passing sample as a guard.

## Q1.6 Review

The complete production population is 57 structs and 74 fields: five adapter
structs with seven fields and 52 private composition structs with 67 fields.
The exact enumeration command was a read-only parser over all production Go
files; its report is `/private/tmp/ply-p4-q16-population.txt`, SHA-256
`0fe1db8870845b245adeb32a8a5f40cb537a9d5efc6d2b4efd1238c50a2deda9`.
Structure and contract review found every field safely defaulted and observed,
and every composition struct delivered whole rather than reconstructed from an
incoming value. This independently supports Q1.6, but no receipt was emitted
because Q1.7 and Q1.9 failed the all-three precondition.

The focused command was:

```sh
go test -json ./... \
  -run 'Test.*(Dependencies|Dependency|Recorded|PassesComplete|PreservesComplete|SelectsCompleteSystem)' \
  -count=1
```

It passed 245 top-level tests across 27 packages with zero failures. The JSON
report SHA-256 is
`ab5a59ebbb65de5564799792ea83cc1775b514a258535f0cfc9de809c26d464e`.

## Q1.7 Review

The complete scope contains 17 operations. These nine have exact content
assertions for their partial failures:

- `Context.FindAndPopulateMavenProjects`
- `Context.OnEachMavenProject`
- `Context.OnRootProject`
- `GitCloudConfig.templates`
- `findFirst`
- `findAll`
- `grepRecursive`
- `filteredFilesFromTemplateWithDependencies`
- `unzipWithDependencies`

These eight do not have complete failure-content coverage and are the exact
blockers:

- `cmd/build_init.go`: the `initCmd` project loop reports missing type,
  initialization, configuration-write, and POM-write failures without a
  command contract.
- `pkg/bitbucket/bitbucket.go`: repository-query warning content is asserted,
  but the clone/pull warning branch in `SynchronizeAllRepos` is not.
- `pkg/spring/io.go`: `DeleteDemoFiles` reports discovery and deletion
  failures while continuing, with no contract reference.
- `pkg/maven/dependency.go`: `upgradeDependencies` continues after version,
  maximum-version, and upgrade failures without log-content assertions.
- `pkg/maven/plugin.go`: `upgradePluginsOnModel` continues after per-plugin
  upgrade errors without a content assertion.
- `pkg/maven/deprecated.go`: `RemoveDeprecated` can return a partial template
  result or warn and continue, with no direct contract.
- `pkg/template/template.go`: `MergeTemplates` warns and continues after a
  template merge failure, with no direct contract.
- `pkg/config/cloud.go`: `ValidTemplatesFrom` returns an ordered partial
  template slice with an error, with no direct contract.

The focused covered-subset run passed 23 top-level tests and 46 total test
events across six packages, with zero failures. Its JSON SHA-256 is
`f7b9bf9827b47c82db7a089ecf317dd3b5d2ca2c57beee511a438f1522f530ee`.
Those 23 tests are evidence of the nine covered operations only; they do not
support a Q1.7 receipt.

## Q1.9 Review

The complete syntactic population is 169 Go test `range` statements. The exact
enumeration is `/private/tmp/ply-p4-q19-all-range-population.txt`, SHA-256
`c4978c74184f94a8369f9f3e4b8c1a48e90db15ab9b6db24749dce45f7e72c25`.
The complete local table class contains 64 loops over the identifier `tests`:
13 enclosing test functions have an executable `len(tests)` assertion and 51
do not. The exact classification report is
`/private/tmp/ply-p4-q19-named-table-audit.txt`, SHA-256
`d832e212662a36c472b1f8cac2d1d8f0e57f87e98dbb4551ab20352c1e521c1c`.

The 51 confirmed table blockers are at:

- `cmd/profile_editor_contract_test.go:82`,
  `cmd/tips_read_contract_test.go:156`.
- `internal/adapter/filesystem/filesystem_test.go:245,305,418,562,776` and
  `internal/adapter/httpclient/client_test.go:185`.
- `pkg/bitbucket/bitbucket_contract_test.go:257,300,454`.
- `pkg/config/cloud_examples_test.go:233`,
  `cloud_git_hooks_test.go:236`, `cloud_refresh_test.go:147`,
  `local_directory_create_contract_test.go:151`,
  `local_directory_stat_contract_test.go:174`,
  `local_touch_create_contract_test.go:168`,
  `local_touch_write_contract_test.go:163`,
  `local_update_create_contract_test.go:168`,
  `local_update_write_contract_test.go:163`, and
  `project_write_contract_test.go:240`.
- `pkg/file/copy_test.go:204,300` and
  `pkg/file/file_test.go:1780,1847,1915,2009,2122,2206,2275,2350,2424,2467`.
- `pkg/http/client_test.go:91`, `download_test.go:260,429`, and
  `pkg/kibana/post_contract_test.go:433`.
- `pkg/maven/command_test.go:86`, `graph_write_contract_test.go:162`, and
  `pkg/shell/git_test.go:168,268,364`.
- `pkg/spring/archive_path_contract_test.go:230`,
  `discovery_contract_test.go:125,175`, and
  `download_contract_test.go:149`.
- `pkg/structurizr/output_write_contract_test.go:160`,
  `pkg/template/markdown_write_contract_test.go:178`,
  `pkg/tips/tips_test.go:307`,
  `pkg/webservice/browser_launcher_contract_test.go:96`, and
  `pkg/webservice/server_contract_test.go:292`.

The remaining 105 ranges include fixture/support loops, structurally guarded
AST searches, output loops with independent cardinality assertions, and
additional unguarded expectation/helper populations. Confirmed additional
gaps include `cmd/entry_boundary_test.go:23,37`,
`cmd/plugin_diagrams_export_contract_test.go:166`,
`internal/adapter/filesystem/archive_test.go:410`,
`internal/adapter/httpclient/client_test.go:330`,
`pkg/config/cloud_templates_test.go:286,347`,
`pkg/file/file_test.go:1552`, `pkg/file/grep_test.go:229`,
`pkg/maven/version_test.go:215`,
`pkg/spring/download_contract_test.go:225`,
`pkg/template/filtered_walk_contract_test.go:283,348`, and
`pkg/tips/tips_test.go:278`. The next move must finish the 105-loop
classification rather than assuming these are the only additional gaps.

The existing empty-population regex run passed 69 named top-level contracts
across 27 packages, with zero failures. Its JSON SHA-256 is
`103ef6620ce810f663e1783cafbc6f6debb227f911dba973133aa60fd3122f7a`.
That covered subset cannot support Q1.9 while the 51-table class fails.

## Authoritative Measurement

Both valid clean audits ran without manual evidence:

```sh
bash .quality/tools/quality-audit.sh . \
  --baseline .quality/baseline/scorecard.json --out "$external_out"
```

They are byte-identical. The authoritative identity is:

- commit: `18beee0a806bcc962cb2f7fbdd5952e43a2cb633`
- commit tree: `85a6c67f38b8f764b849dcfac4da9eec071d22d0`
- status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`
- inventory SHA-256:
  `4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`
- scorecard SHA-256:
  `37ef9bd20cfddf450f381bbc1afe5bd8ea5a037f5a8b0f50713572d22e613f2b`
- output root: `/private/tmp/ply-p4-evidence-review.lUcTJA`

The complete instrument object is:

```json
{"call_scanner_sha256":"73dba9d563fb0b3e85af032487d024c7f42b80fd76de9f928d2888fed7347879","go_build":{"cgo_enabled":"0","go386":"","goamd64":"","goarch":"arm64","goarm":"","goarm64":"v8.0","goenv":"off","goexperiment":"","goflags":"","gomips":"","gomips64":"","goos":"darwin","goppc64":"","goriscv64":"","gowasm":"","gowork":"off","production_files":{"count":93,"sha256":"c9d17a29325acec93f4db4ca11bfcf41b65d02d0ef3e460290fe5b8986743ec5"},"test_files":{"count":70,"sha256":"8b1c804a39428d39cf9609176a24a8c6e261812ab5a31ddffbaa915afa125840"},"version":"go version go1.26.2 darwin/arm64"},"mode":"full mode","parser_sha256":"f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a","partial":false,"q06_contract_sha256":"7881988a831b02e367b21aa0d2ad3d08e25afb7524a5c30f8bfa2a2533c64772","report_template_sha256":"3138b45cd2ad3f3e396ad3c5ea1325130619d6d2c61c2d2c64db262035c0f07d","upstream_exit":1,"upstream_sha256":"45787c231a0f07f27093824b19685aee4fb4cab5159cdb7864834b79fbc9cee9","upstream_version":"1","wrapper_sha256":"b896bb8fac4c81fd5c4aab40732b3064e8355fa0d320dbc4b797541192beda75","wrapper_version":1}
```

Denominators are 453 test functions, 27 packages with tests of 27 packages,
164 Go files, 95 Markdown files, 92 table-driven tests, 13 scripts, five
adapters, eight seams, eight subjects, six non-executable mutation drivers,
and four acceptance scripts. The audit exits 1 for 13 documented non-passing
criteria, never 2. L0 is 8/8; L1 is six PASS and three UNMEASURABLE; ratchets
are six improved, two held, zero regressed, and zero non-comparable; dirty
paths are empty.

A deliberately rejected post-gate audit detected four generated ignored files
under `target/compatibility` and `target/quality-audit` and set
`measurement_clean=false`. Its scorecard SHA-256 was
`a5eb6dd0cbf5db8ffbd46a58d71523e2939166bd62d1d03edaf82ff7ca73c251`;
it is not evidence. Those exact generated files and their now-empty directories
were removed, and the final clean audit reproduced the initial scorecard byte
for byte. Future compatibility gates should set `API_COMPAT_REPORT_OUT` and
`CLI_COMPAT_REPORT_OUT` outside the worktree; the authoritative audit output
must also stay external.

## Gate Result

The corrected complete gate passes API/CLI and CLI-surface compatibility,
fresh subprocess meta-contracts, build, all 453 tests, race, vet, pinned
golangci-lint 2.12.2 with zero issues, `make preflight`, `make test`,
`make test-install`, `make test-agent-start`, all 62 launcher controls, all
Make and production-script contracts, all four host acceptance flows, the
15-control audit meta-suite including the repaired T15 reproduction, and the
empty-HOME count-2 run. The first compatibility attempt was runner-invalid
because a fresh module cache was incompatible with the check's required
offline phase; the unchanged rerun against the populated module cache passed.

## Next Objective

Repair Q1.9 only. Add executable fail-on-empty assertions to every in-scope
assertion-driving collection among the 169 enumerated test ranges. Start with
all 51 confirmed `tests`-table blockers, classify each of the remaining 105
ranges under the rules above, and repair every additional expectation,
recording, or helper population whose emptiness can skip its assertions. Do
not add guards to fixture-construction loops or loops already protected by an
independent executable cardinality assertion. Do not change production code,
the audit apparatus, inventory, baseline, or acceptance behavior.

After the test-only move, rerun exact Q1.9 enumeration and guard validation,
the full gate, and the clean no-evidence audit. Q1.7 will remain blocked and no
manual-evidence document may be created in that session. P4 stays active and
P5-P8 stay queued.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and gate,
both design documents, the Q1.9 criterion source, every one of the 169 range
sites, and the existing empty-population contracts before editing. Confirm
branch, HEAD, clean and ignored status, reciprocal archive links, launcher
`--check`, and exact product ancestry.

Make one focused test-only implementation commit and one separate continuity
commit. Stop before Q1.7 remediation, schema-2 evidence, production changes,
P5 mutation work, T1-T10 changes, P6-P8, publication, or distribution. Do not
push, merge, stash, revert, launch a successor, or remove the worktree.
