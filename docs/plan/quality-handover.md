# Quality Upgrade Handover

Generated: 2026-08-26T21:33:24+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Unchanged production commit: `c382caa38be167fe17f847370ad8a12270644de3`.
- Q1.7 test-only implementation commit:
  `7be93e7e822f05bdb007e4e5a1a00402e620a6cb`.
- Implementation tree: `0f0955f8bbe3e49f5c8d1b4f98454dcb04ea62f1`.
- The implementation commit's exact parent is continuity commit
  `609043491d1a7f7239dc7fd1c2b59c75248ed78c`; that commit's exact parent is
  Q1.9 test-only commit `3a8f0bcd030baf787a29440ee8e4e4a087edb33b`.
- After launch, obtain the new continuity HEAD with
  `git rev-parse --short=12 HEAD`; its exact parent must be `7be93e7`.
- No production, audit, scanner, parser, inventory, baseline, acceptance,
  mutation-harness, or manual-evidence file changed in the Q1.7 move.
- No push, merge, release, publication, distribution, stash, revert,
  successor launch, or worktree removal was performed.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active for one Q1.7
controllability repair and later combined manual evidence, while P5-P8 remain
queued. Its active archive is
`docs/plan/agent-sessions/2026-08-26T213324+0200-repair-q17-delete-demo-discovery.md`.
The Q1.7 characterization archive is answered history and links reciprocally
to the new tail. The graph has exactly one NEXT archive. The launcher's stable
execution region is unchanged; all 62 launcher controls must remain green.

## Q1.7 Result

Q1.7 includes a production collection-item failure path when it reports the
failure and continues, and an operation when it can return a populated partial
aggregate with an error. It excludes fail-fast errors before partial state,
intentionally ignored errors with no reported partial result, non-error
business warnings, sequential non-collection phases, and test support. Every
classification was regenerated from executable control flow; neither names,
search hits, the 23 pre-existing focused contracts, nor exit status was treated
as coverage.

The complete population remains exactly 17 operations. Ten new top-level
contracts in seven `_test.go` files bring the suite from 453 to 463 tests and
repair every branch exercisable through the seams present at the start:

- `initCmd` has exact missing-type, initialization, configuration-write, and
  POM-write failure content plus later-project continuation.
- `Bitbucket.SynchronizeAllRepos` has exact repository-query, clone, and pull
  warnings plus later-project or later-repository continuation.
- `spring.DeleteDemoFiles` has exact selected-test and fixed-demo deletion
  warnings plus later-file continuation.
- dependency upgrade has exact missing-version, invalid-maximum, invalid-current,
  and metadata failure content plus later dependency actions.
- plugin upgrade has exact missing-version, invalid-current, metadata, and
  release failure content plus later model updates.
- `maven.RemoveDeprecated` has exact populated partial result with error and
  exact replacement warning plus continuation.
- `template.MergeTemplates` has exact merge warning plus later-template
  continuation.
- `GitCloudConfig.ValidTemplatesFrom` has its exact ordered populated partial
  result with error.

The other nine operations remain independently covered: `GitCloudConfig.Templates`,
`Context.FindAndPopulate`, `Context.OnEach`, `Context.OnRoot`, `findFirst`,
`findAll`, `grepRecursive`, `filteredFilesFromTemplate`, and `unzip`.
Executable inspection also preserved exact exclusions: Bitbucket project-list
failure and `RemoveDeprecated` discovery failure are fail-fast; dependency
latest-release errors are converted to nil and action errors are intentionally
ignored without a reported partial result; the plugin setter error cannot
occur for an item obtained from the same unremoved model slice.

Sixteen of the 17 operations now have exact content plus continuation or exact
ordered-partial-result coverage. One included branch remains blocked:
`DeleteDemoFiles`' discovery warning. `DeleteDemoFiles` hard-codes public
`file.FindFirst`. That function hard-codes the system filesystem, its
`filepath.Walk` callback ignores every incoming walk error and returns only
`nil` or `io.EOF`, and final `io.EOF` is normalized to nil. Missing roots,
permission errors, and entry-read failures therefore cannot make the error
checked by `DeleteDemoFiles` non-nil. The current mission prohibited the one
production injection seam needed to exercise that branch. Q1.7 remains blocked;
no PASS claim or receipt was made.

External records:

- complete 17-operation branch/classification/coverage record:
  `/private/tmp/ply-p4-q17-classification.tsv`, SHA-256
  `4b247f3343a651374b96245d6c15b563f573a152799f3329713c578ace863e57`;
- exact 33-contract focused population manifest:
  `/private/tmp/ply-p4-q17-focused-population.txt`, SHA-256
  `43391020bb74aa696a750c68cbdefad40ceee4de5928c4f547356f581ddcae22`;
- focused population execution log:
  `/private/tmp/ply-p4-q17-focused-population.log`, SHA-256
  `4c537c3ca61053669d19948db5268cda87b1b364dee67e15277804216e397149`.

Q1.6 remains independently supported by its complete 57-struct, 74-field
production review. Q1.9 remains independently supported by its complete
169-range classification. No manual evidence was created and no schema-2
manual-evidence audit was run.

## Authoritative Measurement

The clean no-evidence audit ran from exact implementation commit `7be93e7`
with all generated output outside the worktree. Its identity is:

- commit: `7be93e7e822f05bdb007e4e5a1a00402e620a6cb`;
- commit tree: `0f0955f8bbe3e49f5c8d1b4f98454dcb04ea62f1`;
- status SHA-256:
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`;
- inventory SHA-256:
  `4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`;
- scorecard SHA-256:
  `b2ccc93dbb045f06404f955ac4b122a614a654f2c0d983d48941fe7026b352ec`;
- raw output SHA-256:
  `4cb2a6e8c05520f4fdb3d3a13568e60c3623a659e0130ceb88a41c15ba506c04`;
- output root: `/private/tmp/ply-q17-audit.Me3kf9`.

The instrument remains wrapper version 1 with wrapper SHA-256
`b896bb8fac4c81fd5c4aab40732b3064e8355fa0d320dbc4b797541192beda75`,
parser SHA-256
`f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a`,
call-scanner SHA-256
`73dba9d563fb0b3e85af032487d024c7f42b80fd76de9f928d2888fed7347879`,
and Go `go1.26.2 darwin/arm64`; the 75-test-file population SHA-256 is
`9d4d2b0e94d56ee0c14c4be685c34f5428570784eeee65e49bb17b4540a665dd`.

The scorecard records 463 test functions, all 27 packages tested, 169 Go files,
95 table-driven tests, 13 scripts, five adapters, eight seams, eight subjects,
six non-executable mutation drivers, and four acceptance scripts. It exits 1
for 13 documented non-passing criteria, never 2. L0 is 8/8; L1 is six PASS and
Q1.6/Q1.7/Q1.9 UNMEASURABLE; ratchets are six improved, two held, zero
regressed, and zero non-comparable; dirty paths are empty.

## Gate Result

The exact affected six-package population and complete 33-contract Q1.7
population pass. Complete uncached tests across all 27 packages, race, and vet
pass. API/CLI and CLI-surface compatibility, fresh subprocess meta-contracts,
pinned golangci-lint 2.12.2 with zero issues, `make preflight`, `make test`,
Make and production-script contracts, all 62 launcher controls, all four host
acceptance flows, the 15-control audit meta-suite including T15, and the exact
empty-HOME count-2 command pass. API, CLI, audit, Go, temporary, and lint
outputs and caches were kept outside the worktree. The authoritative ordinary
and ignored statuses were empty.

An initial broader run encountered the known nondeterministic partial-raw-log
signal fixture in the unchanged launcher suite. Unchanged standalone launcher,
complete `make test`, and complete preflight reruns all passed 62 of 62. The
authoritative launcher rerun SHA-256 is
`bb0054078e20aca6a546ad399b0b55120d045c414523fa50cf36f1038d55fd48`;
the complete preflight rerun SHA-256 is
`0be59f1344b51dd104f958ef088e9a0e42b645b9e0d604032769391faa90d019`.
API and CLI report SHA-256 values are respectively
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

## Next Objective

Add exactly one private function-parameter seam around `file.FindFirst` in
`pkg/spring/io.go`. Keep the exported `DeleteDemoFiles` signature unchanged and
delegate to an unexported helper with `file.FindFirst`; do not add a dependency
struct or inject deletion. Add one test that injects a sentinel discovery error,
asserts exact warning `Unable to find testfile, fileSuffix=.kt`, verifies exact
lookup arguments, and proves later deletion of `HELP.md`, `mvnw`, and
`mvnw.cmd`. This avoids expanding the completed Q1.6 struct/field population.

Reclassify all 17 Q1.7 operations after the repair. If all included paths have
exact content and continuation or ordered-partial-result proof, Q1.7 is
supported, but do not create manual evidence in the seam session. The next
coherent P4 checkpoint would then record combined schema-2 manual evidence for
Q1.6, Q1.7, and Q1.9. P4 stays active and P5-P8 stay queued.

## Start And Stop

Read this handover, the linked NEXT archive, complete P4 entry and gate, both
design documents, Q1.7 criterion, all 17 operations and their focused tests,
`pkg/spring/io.go`, and `pkg/file/file.go` before editing. Confirm branch, HEAD,
clean and ignored status, reciprocal archive links, launcher `--check`, and
exact ancestry.

Make one focused implementation commit containing only the authorized private
seam and exact test, then one separate continuity commit. Stop before another
production change or seam, exported API changes, manual evidence, Q1.6/Q1.9
changes, audit/inventory/baseline changes, P5 mutation work, T1-T10 changes,
P6-P8, publication, or distribution. Do not push, merge, stash, revert, launch
a successor, or remove the worktree.
