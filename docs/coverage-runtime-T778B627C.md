# T778B627C: fresh coverage inputs and runtime budget

This Task follows the existing focused acceptance pattern with
`test/coverage-preflight.py`; it does not run a full suite or qualify the failed
original profile. `test/coverage-runtime-acceptance.sh` is the focused native
acceptance entrypoint. A's independently reviewed and natively passed
verify-00000003 at base `880b508369c61c40823eded3e378a0629e9e54a3` is the
inherited baseline; B does not repeat A's unrelated command suite. B runs both
consumers of its changed installation helper, the changed ordinary-mode journey,
writer/symlink negative guards, sampled runtime groups, Python regressions, vet
and endpoint integrity checks. Root wE:p3 owns independent review and Slack; ply-plan
w9:pE owns the later single fresh Epic-wide coverage attempt.

## Explicit prerequisites

Run from the physical repository/module root, with clean committed source.
The preflight rejects a subdirectory before hashing a partial tree or checking
inputs. It never builds a candidate, installs dependencies, or launches tests.
Provide existing absolute executables, exact source OID/tree, and one consistent
set of candidate/Python aliases. Python must import jsonschema **4.26.0** and
validate Draft 2020-12 under `-I -B`; PYTHONPATH alone is insufficient.

The existing validated local tools are:

- Go `/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go`,
  SHA256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- Python `/Users/perottochristensen/.local/share/ply/gs01-20261009/schema-venv/bin/python`,
  executable SHA256 `01328ad249b1c29f731960419d4145c427b4e3456378d0a53645c4586ae081f9`.
  Preflight also records Python version, package versions and hashes of the
  installed jsonschema, specifications, referencing, rpds and attrs files.

The private final readiness JSON supplies the actual committed candidate path,
OID/tree, SHA256, focused logs and delivery state. A's base-only binary is not a
substitute for the final integrated candidate. Use `make build` with an external
output and the same explicit environment as the test command, then check
`version --format json` and `go version -m`. There is no installation step.

Preflight requires explicit GOOS/GOARCH/CGO_ENABLED, offline module behavior,
GOENV/GOWORK off, GOTOOLCHAIN local and GOFLAGS containing `-mod=readonly`.
Optional `-tags=value` must match the candidate metadata; unsupported GOFLAGS
fail clearly. It checks embedded commit/dirty, toolchain version, target, cgo,
tags and recorded architecture/experiment settings against effective Go inputs.
It captures all effective `go env -json` fields, and compares them at finish,
except `GOGCCFLAGS`, which contains Go's freshly generated temporary prefix-map
path. Explicit cgo flags, compiler choice, architecture features and experiments
remain compared. Preserve the full recorded environment for the actual run.

Before/after manifests hash every tracked repository input, plus the source,
test and embed files selected by offline `go list -deps -test -json ./...`,
including dependency and standard-library source. The Go/Python executables,
candidate and schema package files have separate identities. This records
selected source inputs; it does not snapshot the entire OS SDK, compiler
installation, build cache, dynamic libraries or historical Git archives that
some compatibility tests build. Keep those unchanged and retain their ordinary
Go build provenance; do not claim hermetic reproduction of the whole machine.
Hash equality is endpoint evidence, not continuous monitoring. Declared argv
is explicitly a declaration; it is not proof that the consumer executed it.

## Bounded timeout diagnosis

The retained original invocation used Go 1.26.7 darwin/arm64, cgo 1, no tags,
readonly modules and offline dependency lookup, but omitted candidate/schema
inputs and used Go's default ten-minute **per-package** timeout. The original
`/private/tmp/ply-agentic-crap-ormuj5rm/test.stdout` shows:

- taskrun: 601.068s at `TestBackupGoalRetryAfterNativePreparation`, active 8s;
  the active test is publishing Task content and syncing a directory.
- workspace: 601.417s at
  `TestRecordedProgressDoesNotApplyHistoricalIntegrationToNewProblem`, active 0s;
  its stack is in fixture setup observing Git, with subprocess output readers.
  An old signal receiver waiting nine minutes is not evidence of a stuck test.

A current `go test -list '^Test'` records 668 taskrun tests and 345 workspace
tests. The alarm tests occur at positions 41 and 282 respectively. Therefore
both packages made substantial sequential progress before their alarms; neither
alarm identifies a single ten-minute test. Repeated fallback builds without a
supplied candidate also add cumulative cost (`deliveryCLIBinary` otherwise
builds per request). Candidate inputs remove that fallback where supported;
intentional historical/self-build contracts remain enabled.

Bounded instrumented measurements on A's exact delivered base, with explicit
candidate/Python aliases and `-count=1 -timeout=300s -coverprofile=...`:

| Selection | Package elapsed | Observation |
| --- | ---: | --- |
| taskrun `TestBackupGoal.*` (6 tests) | 44.793s | all pass; alarm test 8.65s |
| workspace `TestRecordedProgress.*` (7 tests) | 15.422s | all pass; alarm test 1.78s |
| taskrun `TestAL05.*`, `TestAL16.*`, `TestBackupCommitted.*` | 128.786s | completed; two fixture failures retained below |

Logs, profiles and list order are preserved in `/private/tmp/T778B627C` as
`runtime.jsonl`, `runtime.cover`, `runtime-prefix.jsonl`,
`runtime-prefix.cover`, and `test-order.txt`. The prefix profile is failed
measurement evidence, never a baseline. Neither measurement is a full package.

For the **one** full consumer attempt, use `-timeout=4h`, `-count=1`, `-json`.
This is a conservative per-package planning ceiling, not a predicted duration:
the original taskrun prefix needed roughly 600s for 40 completed tests; simple
prefix-rate projection across 668 tests is about 10,020s (167min). Four hours
adds about 44% headroom to that deliberately conservative projection. Workspace's
281-test prefix projects roughly 737s across 345 tests; its sampled active test
also completes promptly. A ten-minute limit is demonstrably insufficient for
this filesystem-heavy suite. The sample cannot guarantee the remaining tests
will pass or rule out a later hang. Retain per-test JSON to identify any newly
stalled test; do not automatically retry or further raise the budget on failure.
No sleeps, tests or durability operations were removed to meet this budget.

## Fixture corrections and remaining runtime constraint

The prefix measurement preserved a real AL05 cleanup denial:
`fork/exec /bin/ps: operation not permitted`; writer liveness remained unknown,
with integration delivered and source cleanup pending. That inspection is not
retried by this Task. AL05's ordinary-file-mode journey now stops at the normal
owner-side `owner_handoff:` exit 5 with a positively existing synthetic owner
pane. It asserts delivered exact OID/tree, ownership release, queue closure,
fresh Epic base, one integration, no human QA, and retained source. It does not
claim coordinator cleanup or weaken production liveness guards.

AL16's VCS-free binary used `t.TempDir()` through macOS's `/var` alias; independent
path observation confirms `/var` resolves to `/private/var`. The fixture now
resolves that valid directory before creating its binary. Production still
rejects symlinked inputs. The installation consumer's modeled lifecycle uses
the existing in-process `LocalProcesses` dependency and explicit modeled owner
exit for fixture closeout; public install/check/smoke/rollback remain actual CLI
calls. This avoids crossing into real OS inspection for a synthetic owner.
Unknown/denied/incomplete writer observations and negative symlink tests remain
meaningful and are included in focused acceptance.

A full suite still contains **intentional real process-inspection tests** such
as `TestSystemLocalProcessesObservesRealWriters` and public cleanup journeys.
The consumer needs an authorized runtime where those supported read operations
are available. B does not bypass the preserved denial, change global policy,
probe it again, or claim that its restricted environment can pass that whole
suite. Root must resolve this runtime prerequisite through a supported route
before asking ply-plan to spend the single full attempt.

A's product fixes remain part of the input: exactly one registry error JSON
without polluted stdout, plus the schema's existing backup-goal effect
`preview_or_fresh_goal_registration`. See `coverage-regressions-TB2546BF5.md`.

## Independent review correction

Round 1 found that absent candidate GOEXPERIMENT metadata and empty/unset test
cgo flags could bypass the initial mismatch comparison, even though unchanged
endpoints later passed. `red-r001-f01-isolated.log` preserves nine reproduced
mismatches before the fix. The correction compares both sides after normalizing
missing/empty settings to their default representation. Go 1.26.7's
`cmd/go/internal/load/pkg.go` stamps raw cgo environment values and omits empty
GOEXPERIMENT; GOENV is off here, so absent/empty raw cgo inputs normalize to the
same empty stamped value, rather than the synthesized `go env` compiler defaults.
Architecture/experiment fields compare effective values even when metadata is
absent. Eleven Python tests cover the added experiment mismatch, all four
stamped cgo flags against both empty and unset inputs, and valid default/default
and matching explicit flag cases. Fresh review and final exact-candidate
qualification remain native evidence, with the runtime capability gap separate.

## Exact consumer command

After B's native integration and root's handoff, run the following once from
the clean physical Epic root. Set the two placeholders from root's actual
integrated candidate receipt, and retain the resulting directory. This command
is documentation for ply-plan; B does not execute it.

```sh
set -eu
export PLY_REGRESSION_GO=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go
export PLY_MR05_SCHEMA_PYTHON=/Users/perottochristensen/.local/share/ply/gs01-20261009/schema-venv/bin/python
export PLY_CANDIDATE_BINARY=/absolute/private/final-integrated-candidate/ply
export PLY_CANDIDATE_OID=EXACT_INTEGRATED_OID
export PLY_CANDIDATE_TREE=$(git rev-parse 'HEAD^{tree}')
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_TASK_RUN_TEST_BINARY="$PLY_CANDIDATE_BINARY"
export PLY_SL01_SCHEMA_PYTHON="$PLY_MR05_SCHEMA_PYTHON"
export PATH="$(dirname "$PLY_REGRESSION_GO"):/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 GOARM64=v8.0 GOEXPERIMENT=
export GOENV=off GOWORK=off GOFLAGS=-mod=readonly GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
coverage_run=$(mktemp -d /private/tmp/ply-fresh-coverage.XXXXXX)
export TMPDIR="$coverage_run/" GOTMPDIR="$coverage_run"
"$PLY_CANDIDATE_BINARY" version --format json > "$coverage_run/candidate-version.json"
"$PLY_REGRESSION_GO" version -m "$PLY_CANDIDATE_BINARY" > "$coverage_run/candidate-build.txt"
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/coverage-preflight.py --output "$coverage_run/before.json" -- \
  "$PLY_REGRESSION_GO" test -json -count=1 -timeout=4h "-coverprofile=$coverage_run/coverage.out" ./...
# Stop on preflight failure. Only after its successful exit:
set +e
"$PLY_REGRESSION_GO" test -json -count=1 -timeout=4h "-coverprofile=$coverage_run/coverage.out" ./... \
  > "$coverage_run/test.jsonl" 2> "$coverage_run/test.stderr"
coverage_exit=$?
set -e
printf '%s\n' "$coverage_exit" > "$coverage_run/test.exit"
"$PLY_MR05_SCHEMA_PYTHON" -I -B test/coverage-preflight.py --after "$coverage_run/before.json" \
  --test-exit "$coverage_exit" --output "$coverage_run/after.json"
```

Preflight records cwd, declared argv, effective context, source and executable
hashes. Finish reads sibling `test.jsonl`, records test/package elapsed and
terminal actions, actual supplied exit, profile SHA256 and endpoint differences;
it preserves a failing report and exits nonzero on gaps. Outputs are exclusive
creation, so retries cannot overwrite retained evidence. Keep outputs outside
source. `baseline_qualified` stays false: the consumer must independently
confirm the actual full-suite argv/exit 0, complete log/profile and matching
immutable source/build context before running any analyzer. Never analyze a
profile still being written or reuse the failed original profile.

Go's profile covers instrumented in-process Go execution. Calls to the separate,
uninstrumented CLI candidate do **not** contribute hits to that profile. Do not
merge subprocess profiles silently, infer missing tests from zero statements,
or report a CRAP baseline before the full successful attempt is observed.
