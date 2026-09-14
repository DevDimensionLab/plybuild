# Quality Upgrade Handover

Generated: 2026-09-14T15:19:03+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`. The gRPC Gateway
  evaluation began from clean HEAD
  `5ad98f2c1eb987b153cfcb9078103c64233ea6c1`, parent
  `547460d63ae6eecce2212ce70a00f4b3a56e05c5`, tree
  `5ec9f74d91097b42cb16022e26488678932469c0`.
- The latest dependency implementation remains exact Google UUID v1.4.0
  commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. Gateway and all
  retained modules since Google UUID have no dependency implementation or
  metadata commit.
- The Gateway evaluation archive is answered. The sole NEXT archive is
  `docs/plan/agent-sessions/2026-09-14T151903+0200-decide-grpc-ecosystem-grpc-gateway-product-direction.md`.
  The user has supplied option 1 with the recommended bounds. It does not
  authorize a dependency edit, parent change, `/v2` migration, Go-floor
  change, or another dependency group in that turn.
- No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not launch a successor, push,
  merge, publish, release, stash, revert, bypass cleanup, or remove the
  worktree.

## Roadmap And Guarded Decisions

P2A-P6 are complete. P7 is recording the user's bounded Gateway option 1
decision; P8 remains queued. Exact Go 1.26.7, every accepted dependency move
through Google UUID v1.4.0, and all earlier outcomes are final. Do not combine
another dependency group or begin P8 before this decision is recorded in a
committed handoff.

The user's 2026-09-14 gRPC Prometheus option 1 decision retains exact
`github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0` only under its documented
release-closure, native stream-test, metric counting/classification, archived,
global-state, lifecycle, and qualification exceptions. The gRPC middleware
option 1 decision separately retains exact
`github.com/grpc-ecosystem/go-grpc-middleware v1.0.0` only under its recorded
missing-release-metadata/deterministic-closure, TLS-test, retry cancellation/
timer-resource, and qualification exceptions. Gorilla WebSocket v1.4.2,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy
v0.2.1, and GAX v2.7.0 retain only their separate documented exceptions.

All guards revalidated unchanged during the Gateway evaluation. gRPC
Prometheus v1.2.0, gRPC middleware v1.0.0, and Gorilla WebSocket v1.4.2 each
retain their sole selected-version incoming edge from mvn-pom-mutator v0.2.3;
GopherJS retains its sole GoConvey v1.6.4 edge; Enterprise Certificate Proxy
v0.2.1 and GAX v2.7.0 retain their sole Viper v1.15.0 edges. All six `go mod
why -m` results are negative, source imports are zero, target package loads are
zero, and all remain runtime-unreachable. Fresh primary data has no new exact
record for those targets; Gorilla retains only its recorded entries, including
unwithdrawn GO-2026-6278. A guarded target's direct import/loading, runtime
reachability, version or incoming-edge change, or new advisory/independent
defect expires its exception and requires the owning decision. Do not change
mvn-pom-mutator or GoConvey, reopen an earlier choice, or transfer any
exception to Gateway.

## gRPC Gateway Product Decision

No exact-path stable release qualifies. On 2026-09-14 the user explicitly
selected option 1 with its recommended bounds: retain exact selected,
inherited, unloaded `github.com/grpc-ecosystem/grpc-gateway v1.16.0` without
changing dependency metadata.

The decision accepts only the documented universal discarded-cancel resource
defect, native repeated-suite global and test-state failures, eight inherited
called/reachable vulnerability IDs in the release-native closure, malformed-
generator diagnostic nondeterminism, and related recorded qualification
findings. It does not accept a new or independently discovered defect. The
exception is target-specific and non-transferable.

The exception remains valid only while exact v1.16.0 and both selected-version
incoming edges from etcd/api/v3 v3.5.1 and OTLP v0.7.0 remain unchanged, zero
Gateway packages load, the module remains runtime-unreachable, and no new
advisory or independent disqualifier appears. Revalidate and record those
guards in the committed decision. Direct import or loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent defect expires the exception and requires a fresh dependency and
product decision before merge. Do not add a direct edge, select another v1
release, change either parent, move to `/v2`, raise the Go floor, authorize
parent-removal or architecture work, move unrelated selections, or manufacture
a dependency commit. Do not ask for this same decision again while all guards
hold. All earlier exceptions remain separate.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
public active unarchived non-fork BSD-3-Clause repository
`https://github.com/grpc-ecosystem/grpc-gateway.git`. The proxy exposes 55
versions: 53 stable and two prereleases. V1.16.0 is the latest exact-path
stable release and is lightweight tag/commit
`094a6fe78b3ca888297d090185cdf30f0e42e157`, parent
`e0a026aeb20c3eb0f32cbfdc537c9966634895b9`, tree
`2283b306f85e2b97858d8575954fa0db23db4b5d`, dated
2020-10-28T10:29:51Z with a GitHub-verified commit signature. Its proxy zip
SHA-256 is
`377b03aef288b34ed894449d3ddba40d525dd7fb55de6e79045cdf499e7fe565`;
sumdb source/mod sums are
`h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=` /
`h1:BDjrQk3hbvj6Nolgz8mAMFbcEtjT1g+wF4CSlocrBnw=`. Proxy and Git source are
byte-identical, strict Git verification passes, and license SHA-256 is
`a15b1d1b168954c92ff7fb1620382418f7c72f4f4d251ee791d1098ad68ab0c4`.
The later v1 branch is unreleased. Current main is the distinct
`github.com/grpc-ecosystem/grpc-gateway/v2` module requiring Go 1.26 and is
outside the exact-path decision.

V1.16.0 declares Go 1.14 and its complete isolated production/test closure
preserves Go 1.18. Exact Go 1.26.7/Go 1.18.10 resolve 33 modules, 136/135 graph
edges, 29 packages, 310/247 production entries, and 342/279 complete-test
entries; 12 external modules load and no declared floor exceeds Go 1.14. The
release has 348 files, 179 Go files, 32 native test files, 164 tests, six
benchmarks, one old go-fuzz build-tag target, 30 generated files, four
commands, no Go examples/testdata, and no cgo/embed/symlinks. Pinned v1.15.2
-> v1.16.0 API comparison has zero incompatible changes and nine compatible
Swagger-option additions.

Native count-1, race, broad cross-builds, and targeted test compilation pass
under both SDKs. Both independent count-10 repeats fail under both SDKs:
codegenerator reuses a consumed buffer, runtime tests leak
`DefaultContextTimeout`, and Swagger tests leak global flag state. Vet fails
because production `runtime/context.go` discards the cancel from
`context.WithTimeout`, leaking timer/context resources; two more findings are
test-only `Fatalf` calls from goroutines. Every one of the 53 exact-path stable
releases contains the production discarded-cancel defect. Adjacent v1.15.2
reproduces the repeat and vet failures.

The independent runtime fixture SHA-256 is
`8f097645c3f31bde4b6e4a1b7addd7f92473d89bfe2647862a676d0df3ad3d17`.
It passes count-10 and race under both SDKs while preserving the source vet
finding. It covers HTTP routing, verbs and path parameters; query reflection,
filtering, and the package-global parser setter; JSON/proto/HTTP-body
marshaling; metadata, binary headers, and malformed values; unary/streaming
status and error translation; cancellation/deadlines; cleanup, concurrency,
and nil/panic boundaries. Generator protocol output is deterministic within
and across SDKs at gateway response SHA-256
`ff0829e0e0d3a76ee14a72d0e3e3941b85be7e6d70f5159c91e5b022615dc28c`
and Swagger response SHA-256
`d4e950e6163c8769b44d574fe8350e74b5f35435144048433f0272808a79d85b`.
Malformed generator input exits 255 but emits timestamped, stack-bearing
nondeterministic diagnostics.

Selected v1.16.0 exists through exactly two selected-version incoming edges:
etcd/api/v3 v3.5.1 and OTLP v0.7.0. `go mod why -m` is negative, project
source imports are zero, and zero Gateway packages occur in the normal
429-entry complete-test load, so it is runtime-unreachable. Exact disposable
`go get v1.16.0` only adds a direct root/source checksum, leaves selection and
234 modules unchanged, grows edges 3,599 -> 3,600 and sum lines 1,067 ->
1,068, and leaves target loading zero; tidy removes the root. Exact v1.15.2
instead cascades 234 -> 160 modules and 3,599 -> 2,247 edges, removes the
direct mvn-pom-mutator requirement, and makes the project unloadable. No
projection was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no exact Gateway record and an
empty exact OSV response. Govulncheck v1.8.0 reports 33 inherited module IDs,
17 imported-package IDs, and eight called/reachable IDs in both the v1.16.0
production and test closures: GO-2020-0036, GO-2022-0956, GO-2023-1571,
GO-2023-2153, GO-2024-2687, GO-2025-3372, GO-2026-4762, and GO-2026-6061.
V1.15.2 is identical. The normal project scan remains 30 module IDs, 22
package IDs, and 20 called IDs with zero Gateway occurrence because no target
package loads.

## Project And Quality State

No dependency metadata changed. The project remains 234 modules, 3,599 edges,
429 complete-test entries, 197 module-backed packages, 41 loaded modules,
1,067 `go.sum` lines, and a 432-line tidy projection. `go.mod`/`go.sum`
SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Freshly unpacked exact Go 1.26.7 archive/binary SHA-256 values are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Exact Go 1.18.10 archive/binary values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Both ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.

Exact Go 1.26.7 project verify, count-1, two count-10 repeats, race, vet,
golangci-lint 2.12.2, empty-HOME count-2, Linux/Windows builds, and API/CLI
compatibility pass. The Go 1.18.10 projection loads 366 complete-test entries;
verify, vet, host/Linux/Windows builds, two compatible-package count-10 repeats,
and compatible-package race pass. Full runs retain only the two accepted
Darwin `pkg/shell` closed-file wording assertions.

Focused quality remains 27/27 Q0-Q2 PASS at L2, scorecard SHA-256
`04039eb917cc8aa674c093e866b8fbfb303ab562a28dcbaf2414dba1c9d01c3a`,
with zero held/regressed/non-comparable and seven improved ratchet rows. All
eight mutation meta-suites pass, all 80/80 live mutations are killed,
host/snapshot acceptance pass, and all 15 audit meta-controls pass. Docker
acceptance is unavailable because this host's Docker CLI lacks `buildx`. Two
full quality-wrapper attempts and three nested lifecycle source-archive runs
hit the known signal-retention timing assertion; direct evaluation-time
50-control lifecycle runs and the final 62-check handoff contract pass, and
every applicable quality stage passes independently. This is target-independent
environment/test-design evidence. Preserve the known
apidiff archive reproducibility discrepancy and Python 3.14 Docker timestamp
control; neither is Gateway evidence.

The 7,701-entry external evidence manifest and decision-summary SHA-256 values
are `873610fe31b0da11ba8e91bfd74665be47428ef60c1b589244324ff9b7dc13d7`
and `b73b6153778b43978c2c81738c246d0823a6731a995068ed5c0f8409e0dee354`.
All disposable evidence belongs beneath `$CODEX_SESSION_SCRATCH_ROOT`; never
run `go mod download all` in a measured worktree.

## Next Bounded Objective

The next session must record the user's explicit option 1 decision without a
dependency implementation. Revalidate exact v1.16.0 and both incoming parent
edges, negative why result, zero packages loaded from Gateway and every earlier
guarded target, unchanged module hashes, runtime unreachability, and fresh
advisory state. Then answer the decision archive and prepare exactly one next
bounded P7 mission. Do not execute that successor in the decision-recording
turn.
