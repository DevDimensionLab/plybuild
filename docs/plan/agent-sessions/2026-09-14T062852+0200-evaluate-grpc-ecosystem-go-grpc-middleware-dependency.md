# Agent Session: Evaluate gRPC Middleware Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency`
Created: `2026-09-14T06:28:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ec63e86f50a7e00c9164c489b270a1ef28d8bc319dab6edb3014a10c453762eb`
Previous: [2026-09-14T001519+0200-decide-gorilla-websocket-product-direction.md](2026-09-14T001519+0200-decide-gorilla-websocket-product-direction.md)
Next: [2026-09-14T075103+0200-decide-grpc-ecosystem-go-grpc-middleware-product-direction.md](2026-09-14T075103+0200-decide-grpc-ecosystem-go-grpc-middleware-product-direction.md)
Outcome: Stopped without dependency metadata changes because no exact-path stable release passes every existing contract; preserved the full evaluation and prepared a fresh bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/grpc-ecosystem/go-grpc-middleware v1.0.0` as one bounded
dependency group. Resolve its complete repository and release identity, full
Go-floor closure, package and interceptor behavior, exported API, actual
project loading, exact MVS effects, vulnerability evidence, and every
applicable quality contract. Retain or select only a qualified exact-path
stable release whose complete minimal source/test closure preserves Go 1.18
and whose relevant behavior passes every contract; otherwise stop for a fresh
bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof, Gogo Protobuf, Crypt, OpenCensus Proto, Logex,
Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, historical root
GLFW, Googleapis GAX Go v2 v2.7.0, Google Cloud Go Testing, Enterprise
Certificate Proxy v0.2.1, exact selected GopherJS, and exact selected Gorilla
WebSocket remain retained. All earlier decisions and lifecycle ancestry are
final. Do not revisit them or combine another dependency group. P8 remains
queued.

The user's 2026-09-14 Gorilla WebSocket option 1 decision is final. Exact
`github.com/gorilla/websocket v1.4.2` retains only the documented
GO-2026-6278 weak `math/rand` client-mask behavior and upstream full-source
Go 1.26 cross-test lifecycle-race exceptions while zero target packages load,
it remains runtime-unreachable, its exact version and sole incoming
mvn-pom-mutator v0.2.3 edge remain unchanged, and no new advisory or
independent disqualifier appears. Direct import/loading, runtime reachability,
a version or incoming-edge change, or a new advisory or independent
disqualifier expires the exception. Revalidate these guards before work and
stop for a fresh dependency and product decision if any fails. Do not reopen
Gorilla WebSocket, change mvn-pom-mutator or GoConvey, or transfer its
exceptions to gRPC middleware.

The GopherJS option 1 exception remains valid only for exact
`v0.0.0-20181017120253-0766667cb4d1` and its sole GoConvey v1.6.4 edge
while zero target packages load, it remains runtime-unreachable, and no new
advisory or independent disqualifier appears. Enterprise Certificate Proxy
v0.2.1 retains only its separate exception for its sole Viper v1.15.0 edge
under the same zero-load and unreachable guards. The two GAX v2.7.0
exceptions remain valid only while zero GAX packages load. Revalidate those
guards and stop for a fresh owning decision if one expires. Do not reopen or
transfer any exception.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gorilla WebSocket, GopherJS,
Enterprise Certificate Proxy, Google Cloud Go Testing, GAX, Google pprof, and
Gogo Protobuf were retained without dependency implementation commits or
metadata edits.

The clean Gorilla decision revalidated exact selected v1.4.2 and its sole
mvn-pom-mutator v0.2.3 edge, a negative `go mod why -m` result, and zero
source imports. The unchanged project has 234 selected modules, 3,599 graph
edges, 429 complete-test entries, 197 module-backed packages, 41 loaded
modules, 1,067 `go.sum` lines, and a 432-line unapplied tidy projection.
Exactly zero Gorilla WebSocket, GopherJS, Enterprise Certificate Proxy, or GAX
packages load. `go.mod` and `go.sum` SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

Fresh primary vulnerability data still contains 1,398 module records and was
last modified 2026-09-10T16:28:28Z. Its exact Gorilla WebSocket entry retains
GO-2020-0019, fixed in v1.4.1, and selected-affecting GO-2026-6278; no new
exact advisory appeared. Selected gRPC middleware v1.0.0 is only a starting
MVS selection, not proof of canonical repository, release qualification,
ancestry, floor, loading, behavior, vulnerability state, or suitability.
Independently resolve its incoming graph paths, `go mod why -m` result,
loaded-package population, and every serious exact-path stable candidate. Do
not add a direct edge merely to alter MVS.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and the recorded apidiff source receipt. Preserve the known apidiff archive
reproducibility discrepancy and Python 3.14 Docker timestamp control; neither
is gRPC middleware evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, major
module lines, and every serious exact-path stable candidate. Do not silently
promote a redirect, fork, alternate path, different major module path,
prerelease, non-versioning tag, unreleased branch head, or floor-ineligible
release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph.

Inspect every package, exported API, example, benchmark, fuzz target, testdata,
generated file, and platform or build-tag branch. Characterize unary and
stream interceptor composition and ordering, context and metadata propagation,
logging/auth/recovery/retry/timeout behavior, cancellation and deadlines,
status/error identity, panic and nil boundaries, stream lifecycle, concurrent
use, and network/resource cleanup. Distinguish optional subpackages and example
effects from behavior actually loaded by this project.

Add independent fixtures where useful for interceptor order, short-circuiting,
context and metadata, unary and stream success/failure, cancellation,
deadlines, retries, panic recovery, supported concurrency, deterministic
outputs, cleanup, and selected-project compatibility. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every network, TLS,
clock, timeout, filesystem, subprocess, platform, resource, or test-design
failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, or unrelated-module change needs a
fresh bounded decision rather than silent implementation. Compare primary
vulnerability results at module, package, symbol, test-symbol, and
reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, unchanged exception guards, and
`./codex-dev-start.sh --check`. Read this archive, the answered Gorilla
decision/evaluation, answered GopherJS decision/evaluation, answered Enterprise
Certificate Proxy decision/evaluation, answered GAX archive, rolling handover,
roadmap, `go.mod`, `go.sum`, and every referenced quality, compatibility,
release, runner, evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose complete
minimal source/test closure preserves Go 1.18. Do not promote an unqualified
or floor-ineligible identity. If no candidate satisfies the existing
contracts, preserve the evidence and stop for a fresh bounded product decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an
authorized changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete P7
dependency gate. For a retained or blocked selection, prove the no-change
effect and run all applicable gates without manufacturing activity. Full
changed-selection quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held,
regressed, non-comparable, or dirty counts.

Third, record the exact decision and evidence in the roadmap and rolling
handover, answer this archive, and follow the lifecycle contract for one next
bounded authorized mission. Do not execute the successor in this turn.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the gRPC middleware decision or bounded product stop, rewrite the rolling
handover and roadmap, answer this archive, and follow the repository lifecycle
contract. Do not implement another dependency group, launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Stopped exact-path `github.com/grpc-ecosystem/go-grpc-middleware` for a fresh
bounded product decision. No stable release passes every existing contract,
and the exact `go get` required for a changed selection also moves unrelated
MVS selections. I did not edit `go.mod` or `go.sum`, did not create a
dependency implementation commit, and did not audit another dependency group.

## Repository And Release Identity

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves
canonical `https://github.com/grpc-ecosystem/go-grpc-middleware.git`, a public,
active, unarchived, non-fork Apache-2.0 repository. Strict fsck passes. The
exact module path has seven proxy/GitHub stable releases and no prerelease,
retraction, formal module deprecation, redirect, or qualified alternate path:
v1.0.0, v1.1.0, v1.2.0, v1.2.1, v1.2.2, v1.3.0, and proxy-latest v1.4.0.
Current main is the different
`github.com/grpc-ecosystem/go-grpc-middleware/v2` module and requires Go 1.24;
the v2 line, five unreleased v1-branch commits after v1.4.0, and current main
are not exact-path stable candidates.

Release identities are:

- v1.0.0 lightweight tag commit
  `c250d6563d4d4c20252cd865923440e829844f4e`, parent
  `52645229b435086da88a951d01df177cf3c7cff5`, tree
  `40f5f04f2dab7666784d1596a1a8a0894493c4fa`, at
  2018-05-02T10:16:42+01:00; GitHub reports a verified commit signature.
- v1.1.0 unsigned annotated tag commit
  `dd15ed025b6054e5253963e355991f3070d4e593`, parent
  `6e1e746222237c7c04fd6e0c7ca13ca44d94b5ed`, tree
  `0fade509e9fc0685baceb7e79bad6256fe3ed125`, with author time
  2019-09-12T18:59:40+09:00.
- v1.2.0 unsigned lightweight tag commit
  `3c51f7f332123e8be5a157c0802a228ac85bf9db`, parent
  `6f8030a0b4ee588a3f33556266b552a90a5574e2`, tree
  `0af732c656ad99f8dd100cea801ea9cf40728120`, with committer time
  2020-01-28T17:06:57+01:00.
- v1.2.1 commit `7a5efaad6c58f5b6f96d76e658608730125ac82a`, parent
  `a4522057abfe5f1be6662f16b2e9b44b9fa6d8da`, tree
  `3230f06e51d522cc91adfcf121a1b2d31521a2d4`, dated 2020-07-24;
  GitHub reports a verified commit signature.
- v1.2.2 commit `46f2eb369b917e60df91057c4d37847d17e5a9a4`, parent
  `a0dd2b94302baf7917c42e80a14d50b94f78c7a0`, tree
  `8125f94bfc70636e2dc8c1a29a65bd8601f06acb`, dated 2020-09-09;
  GitHub reports a verified commit signature.
- v1.3.0 commit `df0f91b29bbbdfc3a686a7a8edbe2b9de2072fdd`, parent
  `165f605a7bd1a023d57002bc4e0b64c03c862cd6`, tree
  `f302342837354c61fffa3d88a138efb2a00676a9`, dated 2021-04-22/23;
  GitHub reports a verified commit signature.
- v1.4.0 commit `d42ae9d517069c2bd7f9339147a0eafa86b3d4a3`, parent
  `da1b13ec28bbdd492bdc876045791b69c4be5b81`, tree
  `649f748f6177e19ca3c9922674cf919d05a80096`, with committer time
  2023-03-15T10:41:00Z; GitHub and fresh Go origin metadata report the same
  verified tag commit.

All seven tag commits are ancestors of both the v1 branch and current main.
The v1.1.0-v1.4.0 proxy archives have exactly the Git file population and
byte-identical content. V1.0.0 differs only because the module zip correctly
omits nine Git symlinked README files; all 124 retained files are byte-identical.
The Apache-2.0 license digest is identical across the line. Sumdb, proxy, and
local source/module sums agree:

| Version | Source sum | Module-file sum |
| --- | --- | --- |
| v1.0.0 | `h1:Iju5GlWwrvL6UBg4zJJt3btmonfrMlCDdsejg4CZE7c=` | `h1:FiyG127CGDf3tlThmgyCl78X/SZQqEOJBCDaAfeWzPs=` |
| v1.1.0 | `h1:THDBEeQ9xZ8JEaCLyLQqXMMdRqNr0QAUJTIkQAUtFjg=` | `h1:f5nM7jw/oeRSadq3xCzHAvxcr8HZnzsqU6ILg/0NiiE=` |
| v1.2.0 | `h1:0IKlLyQ3Hs9nDaiK5cSHAGmcQEIC8l2Ts1u6x5Dfrqg=` | `h1:mJzapYve32yjrKlk9GbyCZHuPgZsrbyIbyKhSzOpg6s=` |
| v1.2.1 | `h1:V59tBiPuMkySHwJkuq/OYkK0WnOLwCwD3UkTbEMr12U=` | `h1:EaizFBKfUKtMIF5iaDEhniwNedqGo9FuLFzppDr3uwI=` |
| v1.2.2 | `h1:FlFbCRLd5Jr4iYXZufAvgWN6Ao0JrI5chLINnUXDDr0=` | `h1:EaizFBKfUKtMIF5iaDEhniwNedqGo9FuLFzppDr3uwI=` |
| v1.3.0 | `h1:+9834+KizmvFV7pXQGSXQTsaWhq2GjuNUt0aUU0YBYw=` | `h1:z0ButlSOZa5vEBq9m2m2hlwIgKw+rp3sdCBRoJY+30Y=` |
| v1.4.0 | `h1:UH//fgunKIs4JdUbpDl1VZCDaL56wXCB/5+wF6uHfaI=` | `h1:g5qyo/la0ALbONm6Vbp88Yd8NsDy6rZz+RcrMPxvld8=` |

## Floor, Source, API, And Tests

Selected v1.0.0 has no repository `go.mod` or `go.sum`; its proxy module file
is synthesized as only the module declaration. Its 124-file/90-Go-file source
contains 30 test files, 41 tests, and 18 examples, but cannot provide a
complete deterministic minimal source/test closure: project-graph testing
still lacks declared optional OpenTracing, Zap, Gogo Protobuf, and OAuth
requirements. Its remaining network suites also fail their TLS startup.

V1.1.0 has no `go` directive; v1.2.0 declares Go 1.13; v1.2.1-v1.4.0
declare Go 1.14. Their complete isolated graphs contain respectively
36/40/53/53/49/54 modules and 59/63/179/179/178/208 edges under Go 1.18.10.
The actual imported production/test closures remain within the Go 1.18 floor;
v1.4.0 peaks at Go 1.17, uses 19 production and 21 test modules, and loads
245/309 production/test entries under Go 1.18.10. Thus v1.4.0 is the highest
stable floor-preserving release, but floor eligibility does not qualify it.

The complete release inventories are:

| Version | Files | Go | Test files | Tests | Examples | Benchmarks | Fuzz | Generated |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| v1.0.0 | 124 | 90 | 30 | 41 | 18 | 0 | 0 | 0 |
| v1.1.0 | 112 | 94 | 32 | 46 | 20 | 1 | 0 | 0 |
| v1.2.0 | 125 | 107 | 38 | 53 | 24 | 2 | 0 | 0 |
| v1.2.1 | 126 | 108 | 38 | 57 | 25 | 2 | 0 | 0 |
| v1.2.2 | 122 | 107 | 38 | 57 | 25 | 2 | 0 | 0 |
| v1.3.0 | 130 | 114 | 44 | 67 | 26 | 2 | 0 | 0 |
| v1.4.0 | 129 | 114 | 44 | 68 | 26 | 2 | 0 | 2 |

No release has testdata, cgo, embed directives, or platform/build-tag branches.
Only v1.4.0 has generated files, both test protobuf sources.

V1.4.0 contains 129 files, 114 Go files, 44 test files, 68 tests, 26 examples,
two generated protobuf files, and 23 packages. There are no benchmarks, fuzz
targets, testdata, symlinks, cgo files, embed directives, or platform/build-tag
branches. The packages cover root chaining, auth, rate limiting, recovery,
retry, request tags, validation, logging adapters/context helpers/payloads,
OpenTracing, metadata/backoff utilities, generated test protos, and a reusable
network/TLS interceptor test suite.

Pinned apidiff exported every release. Selected to v1.4.0 has incompatible
generated-protobuf comparability and Logrus/Zap `JsonPbMarshaller` type
changes. V1.2.1 to v1.2.2 removes the exported `testing/certs` package.
V1.3.0 to v1.4.0 incompatibly changes every logging/kit public logger type
from `github.com/go-kit/kit/log` to `github.com/go-kit/log`. The project itself
loads no target package, and the project API/CLI reports remain unchanged and
passing.

V1.1.0-v1.2.1 default tests fail both SDKs because their 2017 CN-only
localhost certificate has no subject alternative name. V1.2.2 and v1.3.0
pass native count-1 and two count-10 repeats under both SDKs. V1.2.2 passes
Go 1.18 race but its Go 1.26 race exposes unsynchronized server listener/
server state during retry-suite restart. V1.3.0 passes race under both SDKs.
V1.4.0 count-1 and race pass, but both independent count-10 runs on both SDKs
panic because an unclosed gRPC client connection continues logging after
`TestReplaceGrpcLoggerV2` or its subtest completes. V1.2.2-v1.4.0 production
cross-builds pass for Linux/amd64 and Windows/amd64 under both SDKs.

Every v1.1.0-v1.4.0 vet run fails on production `retry/perCallContext`
discarding the cancel function returned by `context.WithTimeout`, retaining
timer resources until expiry; v1.0.0 contains the same source statement but
lacks the declared closure needed for an authoritative whole-module vet.

## Interceptor Behavior And Independent Findings

Unary and stream chains execute interceptors left-to-right before the handler
and unwind right-to-left; short-circuit errors preserve identity. Stream
context propagation requires `WrapServerStream`. Auth propagates its returned
context and error verbatim, including service overrides; metadata parsing uses
the first authorization value with case-insensitive scheme matching. Rate
limiting occurs once at RPC start and rejects with ResourceExhausted. Recovery
catches ordinary and nil panics and defaults to Internal; it must be innermost
to catch later handlers/interceptors. Validation prefers legacy `Validate()`
over `Validate(bool)`, returns InvalidArgument, and validates successful stream
receives. Nil functions/receivers/options generally retain panic boundaries.

Retry is client-only and disabled by default. `WithMax` is implemented as the
total attempt count despite its maximum-retries wording. Default codes are
Unavailable and ResourceExhausted. Server-stream retries shallow-buffer sent
message references, resend them, and do not support client/bidirectional
streams. Retry-attempt metadata uses the public misspelled key
`x-retry-attempty`. Backoff observes cancellation only while actually waiting;
per-call timeout contexts are never canceled by the interceptor.

The 29-line release-line fixture at SHA-256
`17bf692b527e101a0a3048e3d8da13b559a414db85541fb6cf98e409c8d1ab76`
proves a production defect in every release under both SDKs: with retry
enabled, zero backoff, and an already-canceled parent, the invoker still runs
three times and the interceptor returns Unavailable rather than Canceled.

The 254-line v1.4.0 fixture at SHA-256
`6c3d2ba53a405da762493a3a2db8480ff70a64cc1761a93f23bbc1116c3dd0c4`
passes count-1/count-10/race positive cases under both SDKs for unary/stream
order, short-circuiting, context/metadata propagation, auth, rate limiting,
panic recovery, validation, retry attempts/headers, unsupported client-stream
retry, and tag aliasing. Its negative cases reproduce the cancellation defect
and independently prove that `logging/settable` forwards variadic arguments
as one nested slice. The settable package claims thread-safe replacement, but
its forwarding changes logger semantics. Request tag maps and logging context
field holders are explicitly not safe for concurrent mutation; metadata clone
is shallow only at API-specified boundaries.

Logrus, Zap, and go-kit adapters attach method/code/duration fields and expose
payload logging, global protobuf marshaller/logger configuration, and real
clock dependence. Stream client logging covers creation rather than terminal
lifecycle; payload wrappers log only successful sends/receives. OpenTracing
propagates cloned metadata and finishes on terminal stream operations, but a
caller that abandons a stream can retain its span, while `CloseSend` may finish
before later receives. The reusable testing package binds loopback TCP,
generates RSA/TLS state in newer releases, uses real deadlines and goroutines,
and needs explicit cleanup. These optional packages and examples are not
loaded by this project.

## Loading, MVS, And Vulnerability Evidence

Selected v1.0.0 exists only through main -> direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` -> target. That edge is
the target's sole incoming edge. `go mod why -m` is negative, repository Go
source has zero target imports, and the 429-entry complete test load contains
zero target packages. It is runtime-unreachable in the current project.

Disposable exact `go get` projections for every release retain zero target
loading but add a direct indirect target requirement and explicit current
Genproto/gRPC/Protobuf roots. Every projection changes four unrelated existing
selections: BigQuery v1.8.0 -> v1.44.0, Datastore v1.1.0 -> v1.10.0, Pub/Sub
v1.3.1 -> v1.27.1, and Envoy control-plane v0.10.1 ->
`v0.10.2-0.20220325020618-49ff273808a1`. V1.2.1 and v1.2.2 additionally move
`golang.org/x/exp` from `v0.0.0-20200224162631-6cc2880d07d6` to
`v0.0.0-20200331195152-e8c3332aa8e5`.

| Exact `go get` | Modules | Graph edges | `go.sum` lines | Delta from 234/3,599/1,067 |
| --- | ---: | ---: | ---: | --- |
| v1.0.0 | 346 | 3,749 | 1,074 | +112 / +150 / +7 |
| v1.1.0 | 347 | 3,765 | 1,082 | +113 / +166 / +15 |
| v1.2.0 | 347 | 3,768 | 1,082 | +113 / +169 / +15 |
| v1.2.1 | 347 | 3,789 | 1,086 | +113 / +190 / +19 |
| v1.2.2 | 347 | 3,789 | 1,086 | +113 / +190 / +19 |
| v1.3.0 | 347 | 3,770 | 1,080 | +113 / +171 / +13 |
| v1.4.0 | 350 | 3,784 | 1,084 | +116 / +185 / +17 |

Each `go.mod` projection adds four requirements. Every tidy projection removes
the target root and returns byte-identically to the base tidy projection,
restoring inherited v1.0.0. A manual v1.4.0 target-only requirement changes
only the target and adds four closure module paths, producing 238 modules and
3,635 edges, but it violates the mandated exact-`go get` workflow and was not
authorized or applied.

Fresh primary vulnerability data contains 1,398 module records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact target record and
preserves every Gorilla/GopherJS/ECP/GAX guard. Reviewed govulncheck v1.8.0
reports base and v1.4.0 project populations identically: 30 module findings,
22 package findings, and 20 IDs/22 production and test-symbol traces, with zero
target occurrence. The isolated v1.4.0 module scan has 28 findings. Its
production package scan has 32 module and 18 package findings; production
symbol analysis has 32 module and 18 package findings plus 63 traces across
five called IDs. Test-symbol analysis has 33 module and 19 package findings
plus 75 traces across the same five called IDs. The first vulnerable frames are
inherited old gRPC/x/net symbols reached through exported testing helpers
(GO-2023-1571, GO-2023-2153, GO-2024-2687, GO-2026-4762, and GO-2026-6061),
not an exact middleware advisory.

## Project Gates And Decision Boundary

The unchanged project remains 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed packages, 41 loaded modules, 1,067
`go.sum` lines, and a 432-line unapplied tidy projection. `go.mod` and `go.sum`
remain byte-identical at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 full preflight, module verification, build, count-1, two
count-10 repeats, race, vet, pinned golangci-lint 2.12.2, empty-HOME count-2,
Linux/Windows builds, API, CLI, 62-control lifecycle, all mutation suites, and
15-control quality-audit meta-suite pass. API/CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The Go 1.18.10 projection removes only the unsupported toolchain line and
passes verify/load/build/vet/Linux/Windows, two count-10 repeats, and race for
the 26 compatible packages; full count-1 retains only the two accepted
`pkg/shell` closed-file wording failures. Accepted quality remains 27/27
Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

The selected 560-file evidence manifest and decision summary SHA-256 values
are `bab1aa974f1182ba6d953831dcfccc0257b6f32f5064c477c68df170b7113b8e`
and `a07237e54d3fac720966f6fc33767d6e98d61db6a2efc8a44de02f78316d3150`.

No stable release qualifies because the canceled-context retry defect is
universal and v1.0.0 additionally lacks a complete declared closure. V1.4.0
is not a resolution: it preserves the retry/timer defects and adds logger,
test-lifecycle, API, and MVS problems. The next bounded session therefore asks
the user to choose among retaining exact unloaded v1.0.0 under a new narrowly
scoped exception, authorizing a separately bounded maintained patch/
replacement, or authorizing a separately bounded parent/removal redesign.
