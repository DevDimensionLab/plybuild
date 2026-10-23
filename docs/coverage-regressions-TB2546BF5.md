# TB2546BF5: deterministic regression acceptance

Task `ply-coverage-regressions-20261023` restores focused acceptance on execution
base `d58dbd4def24bee385307db7dc18fe1583958533`. It does not qualify a full
coverage profile or change provider permissions, retained history, retirement,
registry formats, or integration authority.

## Reproductions and contracts

The unchanged four-test reproduction used the cached Go 1.26.7 darwin/arm64
executable, `GOFLAGS=-mod=readonly GOPROXY=off GOSUMDB=off GOENV=off GOWORK=off
GOTOOLCHAIN=local`, and this command:

```sh
go test ./cmd ./internal/taskexecute ./internal/taskrun \
  -run '^(TestRegistryFailureProducesExactlyOneJSONDiagnostic|TestExecuteKeepsEarlierManualLaunchWhenAutoBecomesDefault|TestSharedEpicUniqueSpecPublicCLIPreparesAndPreservesRetry|TestAL09FollowCandidateCLI)$' \
  -count=1 -timeout=180s
```

It exited 1 with all four original failures. Raw evidence is
`/private/tmp/TB2546BF5/red-original.log`, SHA-256
`2b8839dd9864eb745a3e669f1ea6ed25dce100959ce0efd95955496df78f7a03`.

- `internal/workspace/registry_compatibility.go` supports formats 1 through 7;
  `TestRegistryDiagnosticsUseDeclaredFormatAndValidation` independently uses 99
  for a future version. The command fixture now separates format 99 from
  supported format 6 with missing lifecycle arrays. It preserves rejection,
  exactly one structured stderr diagnostic, and byte-for-byte registry identity.
  Tightening stdout assertions first failed because execute emitted an empty
  result on registry errors (`green-focused-1.log`). Execute now returns that
  error before writing its result, letting the existing root diagnostic own stderr.
- `checkPreservedRuntimeChoices` explicitly retains omitted choices and rejects
  changed explicit profiles, even without a Run. Runtime defaults and the nearby
  backup tests establish that a new Claude launch requests auto. The historical
  fixture now omits that choice when resuming manual, explicitly tests auto drift
  rejection both before and after reservation, and retains its one reservation,
  original request, policy, executable, and absence of new trust authority.
- `delivery_automatic_closeout.go` and existing cleanup tests establish that
  `owner_handoff:` means delivered integration with source cleanup pending.
  The shared Epic journey requires exit 5 plus that marker, closing phase,
  ownership release, delivered exact candidate, target OID/tree, native queue
  closure, and updated native Epic base. It checks one integration result and
  unchanged registry/reflog after readback. Unqualified integration still fails
  with exit 4 and leaves its target untouched. AL09's real timeout exit 5 must
  not claim delivered handoff. No integration command is replayed to prove this.
- `deliveryCLIBinary` already accepts `PLY_DELIVERY_TEST_BINARY` before building.
  That exact provided candidate remains preferred. Its fallback now derives the
  module from the compiled helper source filename using `runtime.Caller`, rather
  than the fixture cwd. The child build alone gets `Cmd.Dir`; it never changes
  process cwd. `TestDeliveryCLIBinaryBuildIgnoresFixtureCWD` clears the override,
  builds and executes the helper from an unrelated temporary cwd, and checks cwd
  preservation. AL09 retains its identity and no-replayed-provider assertions.
  Source builds with trimmed/relative caller paths must supply the exact binary;
  they fail clearly instead of building unrelated bytes.

## Exact candidate acceptance

Supplying the real candidate/Python exposed SL01's capabilities-schema failure
in `acceptance-9892a1d.log`: the catalog's existing backup-goal effect,
`preview_or_fresh_goal_registration`, was absent from the schema enum.
`workflow execute backup-goal` help and `docs/task-writer-evidence.md`
independently establish preview/fresh-registration semantics with no inherited
acceptance or authority. The schema now recognizes that one existing effect.
`red-schema-effect.log` preserves a failing schema regression before the fix;
its positive fixture and negative unknown/missing/wrong-type effect cases keep
the enum restrictive. No command effect, permission or locked answer/input
bundle changed. SL01 still validates the actual full catalog against the schema.

The versioned entrypoint is `test/coverage-regressions-acceptance.sh`; the native
private `acceptance.sh` supplies these existing, validated inputs:

- Go: `/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go`,
  SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- Python: `/Users/perottochristensen/.local/share/ply/gs01-20261009/schema-venv/bin/python`,
  executable SHA-256 `01328ad249b1c29f731960419d4145c427b4e3456378d0a53645c4586ae081f9`.
  `-I -B -c 'import jsonschema'` succeeds with the existing Python 3.14 environment.
- Candidate: native `PLY_CANDIDATE_BINARY`, built from the clean committed source
  with `make build PLY_BUILD_OUTPUT=/absolute/private/ply`, without installation.
  The script checks embedded commit/dirty identity, source OID/tree/cleanliness,
  executable hash stability, and `go version -m` provenance.

The script runs the command-level registry, MR05, SL01 and HI07 contract tests;
registry validation, retained automatic authorization and format upgrade guards;
historical manual launch and backup drift/recovery tests; shared Epic CLI and
MR05 projection; all AL09 observations plus the cwd-independent helper build;
and vet for the three affected packages. Its explicit commands and native
captured output are the verification evidence. No candidate/schema skip,
installation, network dependency download, or unrelated binary fallback is used
for candidate-dependent contracts. The forced helper build separately tests its
established source-build contract.

Final candidate identity, binary digest, focused results, independent review and
integration receipt are private native artifacts, avoiding a self-referential
commit ID in this document. Working logs are preserved under
`/private/tmp/TB2546BF5`; an intermediate log name is not a success claim.

## Remaining timeout scope and coordinator return

The supplied failed aggregate attempt under `/private/tmp/ply-agentic-crap-ormuj5rm`
was not a valid coverage baseline and ran no CRAP analysis. It recorded package
10-minute timeouts while `TestBackupGoalRetryAfterNativePreparation` (8 seconds)
and `TestRecordedProgressDoesNotApplyHistoricalIntegrationToNewProblem`
(0 seconds) were active. Those labels do not establish individual test duration
or causality. The next bounded timeout/preflight Task must inspect retained
`test.stdout`/`test.stderr` stacks and package elapsed budgets rather than simply
increase the timeout. This Task deliberately performs no full `go test ./...`
or full coverage attempt.

Root wE:p3 owns native independent review, integration coordination and Slack.
After exact-candidate review/acceptance and observed local Epic integration,
return the exact integrated OID, candidate/Python provenance and focused logs to
root for the later single Epic-wide attempt owned by ply-plan w9:pE. A delivered
owner handoff still requires root's supported source cleanup after the owner
ends its session; it must not trigger repeated integration or a claim that the
entire Epic is complete.
