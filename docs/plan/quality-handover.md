# Quality Upgrade Handover

Generated: 2026-09-14T11:27:08+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`. The gRPC middleware
  option 1 decision was revalidated from clean HEAD
  `1444e88be20aa8b145c75d4006e30d2d49ac1c08`, tree
  `71824e9b30e68396630e6bf1aca15d70d2a0fb64`.
- The latest dependency implementation remains exact Google UUID v1.4.0
  commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. The gRPC middleware
  evaluation and decision have no dependency implementation or metadata
  commit.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. All earlier
  outcomes and product decisions are final.
- The gRPC middleware evaluation and option 1 decision archives are answered.
  The sole NEXT archive is
  `docs/plan/agent-sessions/2026-09-14T112708+0200-evaluate-grpc-ecosystem-go-grpc-prometheus-dependency.md`.
  It evaluates only selected exact-path gRPC Prometheus v1.2.0 as the next
  bounded P7 dependency group. Do not execute another dependency group or P8
  in that turn.
- No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish, release,
  stash, revert, bypass cleanup, or remove the worktree.

## Lifecycle And Retained Decisions

P2A-P6 are complete. P7 remains active after the user's bounded gRPC middleware
option 1 decision. Exact Go 1.26.7 and the accepted
Speakeasy, XXHash, Fatih Color, Go Logfmt, Go Stack, Godbus D-Bus, Golang
Protobuf, Golang Snappy, Google Martian, Google Renameio, and Google UUID moves
remain final. Google pprof, Gogo Protobuf, Crypt, OpenCensus Proto, Logex,
Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, historical root
GLFW, Googleapis GAX Go v2 v2.7.0, Google Cloud Go Testing, Enterprise
Certificate Proxy v0.2.1, exact selected GopherJS, exact selected Gorilla
WebSocket, and exact selected gRPC middleware remain retained. P8 is queued.
Do not combine another dependency group or begin P8.

The user's 2026-09-14 Gorilla WebSocket option 1 decision remains final. Exact
`github.com/gorilla/websocket v1.4.2` retains only GO-2026-6278's documented
weak `math/rand` client-mask behavior and the documented upstream full-source
Go 1.26 cross-test lifecycle race while:

- exact v1.4.2 and its sole mvn-pom-mutator v0.2.3 edge are unchanged;
- `go mod why -m` stays negative, source imports stay zero, and zero target
  packages load;
- the target remains runtime-unreachable; and
- no new advisory or independent disqualifier appears.

Direct import/loading, runtime reachability, a version or incoming-edge change,
or a new advisory/defect expires the exception and requires its own decision.
Do not change mvn-pom-mutator or GoConvey or transfer the exception.

The GopherJS option 1 exception remains valid only for exact
`v0.0.0-20181017120253-0766667cb4d1` and its sole GoConvey v1.6.4 edge while
zero packages load, it remains runtime-unreachable, and no new advisory or
independent disqualifier appears. Enterprise Certificate Proxy v0.2.1 retains
only its separate exception for the sole Viper v1.15.0 edge under the same
zero-load/unreachable/no-new-finding guards. The two GAX v2.7.0 exceptions
remain valid only while zero GAX packages load. Their guards were unchanged in
this evaluation. None of these exceptions transfers to gRPC middleware.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or bypass launcher cleanup.

## gRPC Middleware Product Decision

No exact-path stable release passes every existing contract. On 2026-09-14 the
user explicitly selected option 1 with the recommended bounds: retain exact
selected, inherited, unloaded
`github.com/grpc-ecosystem/go-grpc-middleware v1.0.0` without a dependency
edit.

The decision accepts only the already documented missing release module
metadata and deterministic complete closure, TLS test-certificate failures in
later candidates, universal retry cancellation and discarded-timer-cancel
resource defects, and related recorded source/test qualification findings. It
does not accept a new or independently discovered defect. The exception is
target-specific and non-transferable.

The exception remains valid only while exact v1.0.0 and its sole incoming
mvn-pom-mutator v0.2.3 edge remain unchanged, zero middleware packages load,
the module remains runtime-unreachable, and no new advisory or independent
disqualifier appears. Direct import or loading, runtime reachability, a target
version or incoming-edge change, or a new advisory or independent defect
expires the exception and requires a fresh dependency and product decision
before merge.
Do not add a direct edge, select v1.4.0 or another release, change
mvn-pom-mutator or GoConvey, raise the Go floor, authorize a patch/replacement
or parent/removal design, move unrelated selections, or manufacture a
dependency commit. Do not ask for this same decision again while all guards
hold. Existing Gorilla WebSocket, GopherJS, Enterprise Certificate Proxy, and
GAX exceptions remain separate.

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves
canonical `https://github.com/grpc-ecosystem/go-grpc-middleware.git`, a public,
active, unarchived, non-fork Apache-2.0 repository. The exact path has seven
stable releases: v1.0.0, v1.1.0, v1.2.0, v1.2.1, v1.2.2, v1.3.0, and
proxy-latest v1.4.0. There is no retraction or formal module deprecation.
Current main is the separate `/v2` module requiring Go 1.24, and the five
post-v1.4.0 v1-branch commits are unreleased.

Exact tag identities are recorded in the answered archive. All tags are
ancestors of v1 and main. V1.1.0-v1.4.0 proxy archives match Git exactly;
v1.0.0 differs only by nine correctly omitted symlinked README files. All
retained content is byte-identical and the Apache-2.0 license digest is stable.

V1.4.0 is the highest stable release whose complete minimal source/test closure
preserves Go 1.18. It declares Go 1.14; its actual closure peaks at Go 1.17,
uses 19 production and 21 test modules, and loads 245/309 production/test
entries under Go 1.18.10. Floor eligibility alone does not qualify it.

## Disqualifying Source And Test Evidence

Every exact v1 release has the same production cancellation defect under both
SDKs. With retry enabled, zero backoff, and an already-canceled parent, the
invoker still runs three times and returns Unavailable rather than Canceled.
The independent 29-line release fixture has SHA-256
`17bf692b527e101a0a3048e3d8da13b559a414db85541fb6cf98e409c8d1ab76`.

Every release also executes `context.WithTimeout` in retry per-call setup while
discarding the cancel function, retaining timer resources until expiry.
V1.1.0-v1.4.0 fail modern vet on this production source. Selected v1.0.0 has
the same statement but no release `go.mod` or deterministic complete test
closure; project-graph testing lacks its optional OpenTracing, Zap, Gogo
Protobuf, and OAuth dependencies.

Release-suite classification:

- v1.1.0-v1.2.1 fail both SDKs because their 2017 CN-only test certificate has
  no SAN; v1.1.0 and v1.2.0 additionally hang until the bounded timeout;
- v1.2.2 passes default/repeats under both SDKs and Go 1.18 race, but Go 1.26
  race finds unsynchronized retry-suite server listener/server restart state;
- v1.3.0 passes default tests, two count-10 repeats, race, and Linux/Windows
  cross-builds under both SDKs, but retains both universal retry defects; and
- v1.4.0 passes count-1/race/cross-builds, but both independent count-10 runs
  on both SDKs panic when unclosed client work reaches a subtest-bound global
  logger after `TestReplaceGrpcLoggerV2` or its subtest returns.

V1.4.0 `logging/settable` forwards `s.get().Info(args)` instead of
`s.get().Info(args...)`, nesting variadic arguments. The independent 254-line
v1.4.0 fixture proves this defect and otherwise passes count-1/count-10/race
positive behavior under both SDKs. Its SHA-256 is
`6c3d2ba53a405da762493a3a2db8480ff70a64cc1761a93f23bbc1116c3dd0c4`.

Pinned apidiff reports selected-to-later generated protobuf comparability and
JSON marshaller incompatibilities. V1.2.1 to v1.2.2 removes the exported
`testing/certs` package. V1.3.0 to v1.4.0 incompatibly changes all logging/kit
public logger types from `github.com/go-kit/kit/log` to `github.com/go-kit/log`.
The known apidiff archive reproducibility discrepancy is preserved separately;
the exact source receipt and executed binary produced valid comparisons.

## Interceptor Behavior

- Unary and stream chains execute left-to-right before handlers and unwind
  right-to-left; short-circuit errors preserve identity. Stream context
  replacement requires `WrapServerStream`.
- Auth is server-side, propagates callback/override context and errors, and
  parses the first authorization metadata value case-insensitively by scheme.
  Rate limiting runs once at RPC start and rejects with ResourceExhausted.
- Recovery catches ordinary and nil panics, defaults to Internal, and must be
  placed innermost to catch later handlers/interceptors. Validation prioritizes
  legacy `Validate()`, returns InvalidArgument, and validates successful stream
  receives.
- Retry is client-side and disabled by default. `WithMax` is total attempts,
  default codes are Unavailable/ResourceExhausted, the public attempt key is
  misspelled `x-retry-attempty`, and stream retry supports only server streams
  with shallow buffering/resends.
- Logging adapters expose global mutable logger/marshaller state and real-clock
  duration. Stream logging covers creation rather than terminal lifecycle;
  payload wrappers log successful sends/receives. Tags and logging context
  holders are not safe for concurrent mutation.
- OpenTracing clones metadata and finishes spans at terminal operations, but an
  abandoned stream can retain its span and `CloseSend` may finish before later
  receives. Testing helpers bind loopback TCP, generate TLS/RSA state, and use
  real deadlines/goroutines, requiring explicit cleanup.

These optional packages and examples are not loaded by this project.

## Loading And MVS

Selected v1.0.0 exists only through main -> direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` -> v1.0.0. This is the
sole incoming target edge. `go mod why -m` is negative, repository source has
zero target imports, and the 429-entry complete test load has zero target
packages. The target is runtime-unreachable.

Exact `go get` projections for every candidate retain zero loading. V1.4.0
changes the target plus four unrelated existing selections: BigQuery v1.8.0
to v1.44.0, Datastore v1.1.0 to v1.10.0, Pub/Sub v1.3.1 to v1.27.1, and Envoy
control-plane v0.10.1 to
`v0.10.2-0.20220325020618-49ff273808a1`. It adds 116 graph modules and moves
234 modules/3,599 edges/1,067 sum lines to 350/3,784/1,084. Tidy returns
byte-identically to the base tidy projection and restores inherited v1.0.0.
A manual target-only requirement avoids the four existing-version moves but
violates the mandated exact-`go get` workflow and is not authorized.

## Vulnerability And Project Gates

Fresh primary vulnerability data has 1,398 module records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
Last-Modified 2026-09-10T16:28:28Z, and no exact target record. Reviewed base
and v1.4.0 project scans are identical at 30 module findings, 22 package
findings, and 20 unique IDs/22 production and test-symbol traces, with zero
target occurrence. The isolated v1.4.0 module scan has 28 findings; production
package/symbol analysis has 32 module, 18 package, and 63 symbol traces across
five called IDs; test-symbol analysis has 33 module, 19 package, and 75 traces
across the same five IDs. These are inherited old gRPC/x/net findings reached
through testing helpers, not an exact middleware advisory. Guarded Gorilla
data remains only GO-2020-0019 and selected-affecting GO-2026-6278.

The unchanged project remains:

- 234 selected modules and 3,599 graph edges;
- 429 complete-test entries, 197 module-backed packages, and 41 loaded modules;
- 1,067 `go.sum` lines and a 432-line unapplied tidy projection; and
- `go.mod`/`go.sum` SHA-256
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 full preflight, verify, build, count-1, two count-10 repeats,
race, vet, pinned golangci-lint 2.12.2, empty-HOME count-2, Linux/Windows
builds, API/CLI, 62 lifecycle controls, all mutation suites, and 15 audit meta
controls pass. API/CLI report hashes remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

The contained Go 1.18.10 projection passes verify/load/build/vet/cross-build,
two count-10 repeats, and race for 26 compatible packages. Full count-1 retains
only the two accepted `pkg/shell` closed-file wording failures. Accepted
quality stays 27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

The selected 560-file evidence manifest and decision summary hashes are
`bab1aa974f1182ba6d953831dcfccc0257b6f32f5064c477c68df170b7113b8e`
and `a07237e54d3fac720966f6fc33767d6e98d61db6a2efc8a44de02f78316d3150`.

The decision-recording revalidation used a freshly unpacked exact Go 1.26.7
archive and binary at SHA-256
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Exact middleware v1.0.0 and its sole mvn-pom-mutator v0.2.3 edge remain
unchanged. Gorilla WebSocket retains its sole mvn-pom-mutator edge, GopherJS
its sole GoConvey v1.6.4 edge, and Enterprise Certificate Proxy and GAX their
Viper v1.15.0 edges. All five guarded `go mod why -m` results are negative,
repository Go source imports are zero, and the 429-entry load has 197
module-backed packages across 41 modules with zero packages from every guarded
target. All remain runtime-unreachable.

The project remains 234 modules, 3,599 graph edges, and 1,067 sum lines;
`go.mod` and `go.sum` retain their recorded hashes. Exact Go 1.26.7 module
verification, build, full count-1 tests, full race tests, and vet pass. The
62-control launcher lifecycle suite passes on an immediate unchanged rerun
after its first nested source-archive signal-retention probe hit the already
documented timing flake. Fresh primary vulnerability data is byte-identical at
1,398 records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. There is no exact middleware,
GopherJS, Enterprise Certificate Proxy, or GAX record. Gorilla retains only
GO-2020-0019 and selected-affecting unwithdrawn GO-2026-6278 at record
SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
No new advisory or independent disqualifier appeared.

## Next Bounded Objective

The next session must independently evaluate selected exact-path
`github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0` as one bounded P7
dependency group. Resolve its exact repository and release identity, complete
Go-1.18 floor closure, Prometheus collector/interceptor behavior, exported API,
actual loading, MVS effects, vulnerability evidence, and every applicable
quality contract before retaining or selecting a qualified exact-path stable
release. Preserve the separate gRPC middleware, Gorilla WebSocket, GopherJS,
Enterprise Certificate Proxy, and GAX guards and stop for the owning decision
if any expires. Do not execute another group or begin P8.
