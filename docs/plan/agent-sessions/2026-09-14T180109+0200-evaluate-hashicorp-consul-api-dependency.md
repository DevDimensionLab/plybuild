# Agent Session: Evaluate Hashicorp Consul API Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T180109+0200-evaluate-hashicorp-consul-api-dependency`
Created: `2026-09-14T18:01:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8eeffe8e08009da9f2654dbb59aa3a0d6c45ab7bea15d69b3091ce9058a5b75b`
Previous: [2026-09-14T151903+0200-decide-grpc-ecosystem-grpc-gateway-product-direction.md](2026-09-14T151903+0200-decide-grpc-ecosystem-grpc-gateway-product-direction.md)
Next: [2026-09-14T193201+0200-decide-hashicorp-consul-api-product-direction.md](2026-09-14T193201+0200-decide-hashicorp-consul-api-product-direction.md)
Outcome: No exact-path stable release qualified; preserved exact inherited and unloaded v1.18.0 without metadata changes, recorded the complete evaluation, and stopped for a fresh bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/consul/api v1.18.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, and all recorded retained-module
decisions through gRPC Gateway v1.16.0. All earlier outcomes and lifecycle
ancestry are final. Evaluate only Hashicorp Consul API in this session; do not
reopen or combine another dependency group. P8 remains queued.

The user's 2026-09-14 Gateway option 1 decision retains exact selected,
inherited, unloaded `github.com/grpc-ecosystem/grpc-gateway v1.16.0` without
dependency metadata changes. It accepts only the documented universal
discarded-cancel resource defect, native repeated-suite global and test-state
failures, eight inherited called/reachable vulnerability IDs in the release-
native closure, malformed-generator diagnostic nondeterminism, and related
recorded qualification findings. Its exception is target-specific and valid
only while exact v1.16.0 and the selected-version etcd/api/v3 v3.5.1 and OTLP
v0.7.0 edges remain unchanged, zero Gateway packages load, the module remains
runtime-unreachable, and no new advisory or independent disqualifier appears.
Direct import/loading, runtime reachability, a version or incoming-edge change,
or a new advisory or independent defect requires a fresh Gateway dependency
and product decision before merge. Revalidate these guards before work and
stop for that owning decision if any fails. Do not transfer its exception to
Hashicorp Consul API or change either Gateway parent.

The gRPC Prometheus v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy
v0.2.1, and GAX v2.7.0 decisions remain final and separate under their exact-
selection, selected-version incoming-edge, zero-load, runtime-unreachable, and
no-new-finding guards. Revalidate those guards and stop for the fresh owning
decision if one expires. Do not change mvn-pom-mutator, GoConvey, Viper, or
transfer an earlier exception.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gateway and all retained
modules since Google UUID have no dependency implementation or metadata
commit.

The clean Gateway decision retained exact v1.16.0 through exactly the two
selected-version incoming edges from etcd/api/v3 v3.5.1 and OTLP v0.7.0. Its
`go mod why -m` result remains negative, repository source imports are zero,
and zero Gateway packages occur in the complete project load. The unchanged
project has 234 selected modules, 3,599 graph edges, 429 complete-test entries,
197 module-backed packages, 41 loaded modules, 1,067 `go.sum` lines, and a
432-line unapplied tidy projection. Exactly zero Gateway, gRPC Prometheus,
gRPC middleware, Gorilla WebSocket, GopherJS, Enterprise Certificate Proxy,
or GAX packages load. `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains byte-identical at 1,398 module
records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact Gateway, gRPC
Prometheus, gRPC middleware, GopherJS, Enterprise Certificate Proxy, or GAX
record. Gorilla retains only its recorded entries, including unwithdrawn
GO-2026-6278 at record SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
The exact Gateway OSV response remains empty. Selected Hashicorp Consul API
v1.18.0 exists in the current build list, but that starting fact is not proof
of its repository, release qualification, ancestry, floor, incoming graph,
loading, behavior, vulnerability state, or suitability. Resolve those facts
independently and do not add a direct edge merely to alter MVS.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`
and GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
Preserve the known apidiff archive reproducibility discrepancy, Python 3.14
Docker timestamp control, Docker buildx absence, and nested launcher signal-
retention timing race; none is Hashicorp Consul API evidence.

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

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize client configuration, request construction,
transport behavior, authentication and TLS inputs, query and option encoding,
context/cancellation/deadline handling, retries, concurrency, global state,
resource cleanup, error identity, and nil/panic boundaries. Distinguish tools,
examples, and optional packages from behavior actually loaded by this project.

Add independent fixtures where useful for request/response semantics,
configuration and options, cancellation/deadlines, error identity, malformed
inputs, deterministic outputs, cleanup, supported concurrency, and selected-
project consumers. Run source verification, package listing, native complete
tests, two independent repeats, race, vet, and meaningful cross-builds under
both SDKs. Classify every network, TLS, clock, timeout, filesystem, subprocess,
platform, resource, global-state, or test-design failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent implementation.
Compare primary vulnerability results at module, package, symbol, test-symbol,
and reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, every unchanged exception guard, and
`./codex-dev-start.sh --check`. Read this archive, the answered Gateway
decision and evaluation, the earlier guarded-decision archives, rolling
handover, roadmap, `go.mod`, `go.sum`, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose complete
minimal source/test closure preserves Go 1.18. Do not promote an unqualified or
floor-ineligible identity. If no candidate satisfies the existing contracts,
preserve the evidence and stop for a fresh bounded product decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, record the exact decision and evidence in the roadmap and rolling
handover, answer this archive, and follow the lifecycle contract for one next
bounded authorized mission. Do not execute the successor in this turn.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the Hashicorp Consul API decision or bounded product stop, rewrite the
rolling handover and roadmap, answer this archive, and follow the repository
lifecycle contract. Do not implement another dependency group, launch a
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Stopped exact-path `github.com/hashicorp/consul/api` for a fresh bounded
product decision without changing dependency metadata. V1.18.0 is the highest
stable exact-path release whose complete imported closure preserves the Go
1.18 floor, but it does not satisfy every qualification contract. No exception
was granted and no dependency implementation commit was manufactured.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
public active unarchived non-fork repository
`https://github.com/hashicorp/consul.git`, with module subdirectory `api` and
MPL-2.0 at the selected source. The v1 proxy lists 96 entries: 84 stable exact-
semver releases and 12 prerelease/suffix entries. Selected v1.18.0 is
lightweight tag `api/v1.18.0`, commit
`13836d5ca84c71b35f201da06dd75f5ad6699c36`, parent
`18dffc51de5587f8fba5f188a636d1864bf2b98a`, repository tree
`1960701737d9ae23b0dc7d24a5087a4aa0541f83`, API subtree
`b088d6696716b2cb8b61c824299dca34cab203f4`, dated
2022-11-30T18:59:52Z with a GitHub-verified commit signature. Its proxy zip
SHA-256 is
`0dc6cfca8c71b05b3ba859726378d2ee611c15304fc85c2c030e3366ee068062`;
proxy and Git API sources are byte-identical, strict shallow Git verification
passes, and sumdb source/mod sums are
`h1:R7PPNzTCeN6VuQNDwwhZWJvzCtGSrNpJqfb22h3yH9g=` /
`h1:owRRGJ9M5xReDC5nfT8FTJrNAPbT4NM6p/k+d03q2v4=`.

All 83 retrievable stable v1 go.mod files were inspected. V1.0.0 through
v1.18.0 declare Go 1.12; v1.18.1 through v1.31.0 declare Go 1.19 and later
releases rise further. Latest v1.34.5 declares Go 1.26.7. Thus v1.18.0 is the
highest declaration-eligible stable release and v1.18.1 is the first floor-
ineligible release. The distinct `github.com/hashicorp/consul/api/v2 v2.0.0`
path declares Go 1.26 and is outside this exact-path decision. Latest v1 also
records retractions for mutated or invalid tags, including unavailable
v1.21.2; no redirected, alternate, prerelease, branch-tip, or v2 identity was
silently promoted.

The selected release has 78 module files and 74 Go files: 40 production and 34
tests across `api` and `watch`, with 217 tests, no examples, benchmarks, fuzz
targets, commands, testdata, generated files, cgo, embeds, or go:generate
directives. Its published go.mod declares Go 1.12 but replaces Consul SDK with
`../sdk`, and tests require 14 certificate fixtures from
`../test/client_certs`; the proxy module archive is therefore not standalone.
Qualification used exact tagged API, SDK, and test-fixture subtrees.

Both SDKs resolve 56 modules. Exact Go 1.26.7 has 154 graph edges, 206
production entries, and 270 complete-test entries; Go 1.18.10 has 153 edges,
144 production entries, and 207 complete-test entries. Fourteen production and
33 test modules import, and no imported module declares above Go 1.17. Module
verification, production build, vet, and darwin/amd64, linux/amd64,
linux/arm64, and windows/amd64 production cross-builds pass under both SDKs.
The watch package also passes test compilation, count-one, two independent
count-ten repeats, and race under both SDKs.

The independent disqualifiers are:

- the root API test binary cannot link on Darwin or Linux under Go 1.26.7
  because old `x/net/internal/socket` references `syscall.recvmsg`; adjacent
  v1.17.0 reproduces the defect;
- the `consulent` root test branch does not compile under either SDK because
  `defaultNamespace` and `defaultPartition` exist only in the excluded
  `!consulent` test file;
- the release module is not independently testable because of its relative SDK
  replacement and repository-external certificate fixtures;
- with official Consul v1.14.2 and a scratch-only freeport sandbox shim, four
  `TestAPI_ClientTLSOptions` cases fail because the tagged certificates expired
  on 2023-11-01; the other 213 tests pass each repeated and race iteration and
  no data race appears; and
- every selected call site discards response-metadata parsing errors, so
  malformed Consul index, cache, hash, and query-backend headers are silently
  accepted with default metadata.

The independent consumer fixture SHA-256 is
`40e09f7b5f6003da6baf17cb9aa45d60cf8319d53732559186f2aeb6e0c7a90d`.
It passes verify, count-one, two count-ten repeats, race, and vet under both
SDKs. It covers configuration, paths and request/options encoding, auth/token/
TLS inputs, cancellation/deadlines, error identity, no-retry behavior,
cleanup, malformed JSON and metadata, nil/panic boundaries, invalid schemes/
keys, and concurrent requests/header use. Pinned apidiff reports zero
incompatible and one compatible change from v1.17.0 to v1.18.0, zero
incompatible and seven compatible additions from v1.18.0 to v1.18.1, and five
incompatible changes from v1.18.0 to v1.34.5.

Selected v1.18.0 exists solely through Viper v1.15.0. `go mod why -m` is
negative, repository imports are zero, zero target packages occur in the
429-entry complete project load, and it is runtime-unreachable. Disposable
exact `go get` projections all retain zero load: v1.17.0 downgrades Viper;
v1.18.0 manufactures eight indirect roots and 26 sum lines before tidy removes
them byte-identically to base tidy; v1.18.1 upgrades x/net and x/text while
exceeding the floor; and v1.34.5 raises the main Go declaration, moves many
unrelated selections, and has incompatible API. No projection was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no Consul API module record.
Exact OSV queries for v1.17.0, v1.18.0, v1.18.1, and v1.34.5 are empty. Base
and selected project scans are identical at 30 module IDs, 22 package IDs, and
20 called IDs with zero target occurrence. The consumer fixture has zero
findings. The selected production closure has inherited module-only
GO-2026-5024 in old x/sys with no loaded vulnerable package or called symbol;
tests add inherited GO-2022-0603 in yaml.v3 and load that test package but call
no vulnerable symbol.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 verify, build, count-one, two count-ten repeats, race, vet,
pinned lint, API/CLI compatibility, host and snapshot acceptance, empty-HOME
count-two, and four production cross-builds pass. The Go 1.18 projection loads
366 complete-test entries; its 26 unaffected packages pass two count-ten
repeats and race, and vet plus host/four cross-builds pass. Full Go 1.18 runs
retain only the two accepted Darwin shell closed-file wording assertions.

All 15 audit meta-controls pass. A raw audit without the required external
manual-evidence receipt exits 1 by contract and is not comparable to the
accepted manual-evidence-adjusted 27/27 Q0-Q2 L2 state. No implementation or
dependency metadata changed, so accepted quality remains unchanged. Full
preflight's substantive stages pass, while its nested launcher test retains
the known signal-retention timing race. One completed-handoff run passes all
62 controls; three later final-text runs pass outer controls 1-50 and reproduce
only that nested signal-log timing failure at control 51. Buildx is now present
at v0.33.0-
desktop.1; Docker acceptance builds and tests the image, then reproduces the
known Python 3.14 nanosecond timestamp-control incompatibility. Neither is
Consul API evidence.

All previously guarded dependency exceptions revalidated unchanged; none
transfers to Consul API. The 9,780-entry evidence manifest and decision-summary
SHA-256 values are
`966263870f7529dcdd4713eba204dc1a8a00e17901a5f702557467b3c79ac46e`
and `5957399af9cf4564b9864e4b6dea8c84e8912db3846e446c953345dadbddb08a`.
Preserve exact inherited, unloaded v1.18.0 and all metadata unchanged pending
the next bounded decision; do not silently retain it under an exception,
select another release, change Viper, move to `/v2`, patch or replace it, or
raise the Go floor.
