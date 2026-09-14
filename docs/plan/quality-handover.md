# Quality Upgrade Handover

Generated: 2026-09-14T21:55:34+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`. The Consul SDK evaluation
  began from clean HEAD
  `19e0f4a2e36c654c9e1810b8f660cec2a29f22fb`, parent
  `474909583bfe6efbc5ae49b9ff1d68ec7402ee41`, tree
  `e1d077544884d461318861fb1d6aad72e8caa202`.
- The latest dependency implementation remains exact Google UUID v1.4.0
  commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. Consul API, Gateway,
  and every retained group since Google UUID have no dependency implementation
  or metadata commit.
- The Consul SDK evaluation archive is answered. The sole NEXT archive is
  `docs/plan/agent-sessions/2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction.md`.
  The user has supplied option 1 with the recommended bounds. It does not
  authorize a dependency implementation, reopening Consul API, another
  dependency group, or P8 in that turn.
- No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not launch a successor, push,
  merge, publish, release, stash, revert, bypass cleanup, or remove the
  worktree.

## Roadmap And Guarded Decisions

P2A-P6 are complete. P7 is recording the user's bounded Consul SDK option 1
decision after exact Go 1.26.7, every accepted dependency move through Google
UUID v1.4.0, and all recorded retained-module decisions. P8 remains queued.
Earlier outcomes and lifecycle ancestry are final. Do not combine another
dependency group or begin P8 before this decision is recorded in a committed
handoff.

The user's 2026-09-14 Gateway option 1 decision retains exact inherited,
unloaded `github.com/grpc-ecosystem/grpc-gateway v1.16.0` under only its
documented discarded-cancel resource defect, repeated-suite global/test-state
failures, eight inherited called vulnerability IDs, malformed-generator
diagnostic nondeterminism, and related completed qualification findings. It is
valid only while exact v1.16.0 and both etcd/api/v3 v3.5.1 and OTLP v0.7.0
incoming edges remain unchanged, zero packages load, the module remains
runtime-unreachable, and no new advisory or independent defect appears.

The gRPC Prometheus v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy
v0.2.1, and GAX v2.7.0 decisions remain separate under their exact-selection,
selected-version incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. At SDK evaluation start Consul API retained its sole
Viper v1.15.0 edge; Gateway retained exactly its two recorded incoming edges;
gRPC Prometheus, gRPC middleware, and Gorilla retained their sole mvn-pom-
mutator v0.2.3 edges; GopherJS retained its sole GoConvey v1.6.4 edge; and
Enterprise Certificate Proxy and GAX retained their sole Viper v1.15.0 edges.
All nine why results, including Consul SDK, are negative; source imports and
target package loads are zero, and runtime unreachability remains intact.

Fresh primary data has no new exact record for those targets. Gorilla retains
only its recorded entries, including unwithdrawn GO-2026-6278. Direct import or
loading, runtime reachability, a version or incoming-edge change, or a new
advisory/independent defect expires the owning exception and requires its fresh
decision. Do not change either Gateway parent, mvn-pom-mutator, GoConvey, or
Viper; reopen an earlier choice; or transfer any exception to Consul SDK.

## Hashicorp Consul API Product Decision

No exact-path stable release qualifies. On 2026-09-14 the user explicitly
selected option 1 with the recommended bounds: retain exact selected,
inherited, unloaded `github.com/hashicorp/consul/api v1.18.0` without changing
dependency metadata.

The decision accepts only the documented non-standalone release closure, Go
1.26 Unix native-test link incompatibility, broken `consulent` test branch,
expired TLS fixtures, silently discarded response-metadata errors, inherited
closure-only vulnerability findings, documented nil/panic and mutation
boundaries, and related completed qualification findings. It does not accept a
new or independently discovered defect. The exception is target-specific and
non-transferable.

The exception remains valid only while exact v1.18.0 and its sole Viper v1.15.0
selected-version incoming edge remain unchanged, zero Consul API packages load,
the module remains runtime-unreachable, and no new advisory or independent
disqualifier appears. Direct import or loading, runtime reachability, a target
version or incoming-edge change, or a new advisory or independent defect
expires the exception and requires a fresh Consul API dependency and product
decision before merge. Do not add a direct edge, select another v1 release,
move to `/api/v2`, change Viper, raise the Go floor, authorize parent
modernization or removal, replace the architecture, move unrelated selections,
or manufacture a dependency commit. Do not ask for this same decision again
while all guards hold. All earlier exceptions remain separate.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
public active unarchived non-fork repository
`https://github.com/hashicorp/consul.git`, module subdirectory `api`, and
MPL-2.0 at selected source. The v1 proxy lists 96 entries: 84 stable exact-
semver releases and 12 prerelease/suffix entries. Selected v1.18.0 is
lightweight tag `api/v1.18.0`, commit
`13836d5ca84c71b35f201da06dd75f5ad6699c36`, parent
`18dffc51de5587f8fba5f188a636d1864bf2b98a`, repository tree
`1960701737d9ae23b0dc7d24a5087a4aa0541f83`, API subtree
`b088d6696716b2cb8b61c824299dca34cab203f4`, SDK subtree
`7c45f454d71c98d13e82a36c1eefb9d8762b9442`, dated
2022-11-30T18:59:52Z with a GitHub-verified commit signature. The tag is
unsigned and has no GitHub Release object. Proxy zip SHA-256 is
`0dc6cfca8c71b05b3ba859726378d2ee611c15304fc85c2c030e3366ee068062`;
proxy and Git source are byte-identical, strict Git verification passes, and
sumdb source/mod sums are
`h1:R7PPNzTCeN6VuQNDwwhZWJvzCtGSrNpJqfb22h3yH9g=` /
`h1:owRRGJ9M5xReDC5nfT8FTJrNAPbT4NM6p/k+d03q2v4=`.

All 83 retrievable stable v1 go.mod files were inspected. V1.0.0 through
v1.18.0 declare Go 1.12; v1.18.1 through v1.31.0 declare Go 1.19 and later
releases rise further. Latest v1.34.5 declares Go 1.26.7. V1.18.0 is therefore
the highest declaration-eligible stable candidate; v1.18.1 is the first floor-
ineligible one. The distinct `/api/v2 v2.0.0` path declares Go 1.26 and is out
of scope. Latest v1 records mutated/invalid-tag retractions, including proxy-
listed but unavailable v1.21.2. No redirect, alternate path, prerelease,
branch head, or different major line was promoted.

The selected release has 78 module files and 74 Go files: 40 production and 34
tests in `api` and `watch`, with 217 tests, no examples, benchmarks, fuzz
targets, commands, testdata, generated files, cgo, embeds, or go:generate
directives. Its published module is not standalone: go.mod replaces Consul SDK
with `../sdk`, and tests read 14 certificate fixtures from
`../test/client_certs`. Qualification used exact tagged API, SDK, and test
subtrees.

Both SDKs resolve 56 modules. Exact Go 1.26.7 has 154 graph edges, 206
production entries, and 270 complete-test entries; Go 1.18.10 has 153 edges,
144 production entries, and 207 complete-test entries. Fourteen production
and 33 test modules import, and none declares above Go 1.17. Verification,
production build, vet, and darwin/amd64, linux/amd64, linux/arm64, and
windows/amd64 production cross-builds pass under both SDKs. The watch package
also passes test compilation, count-one, two count-ten repeats, and race under
both SDKs.

The selected release fails qualification independently:

- its root API test binary cannot link on Darwin or Linux under Go 1.26.7
  because old x/net references `syscall.recvmsg`; adjacent v1.17.0 reproduces
  this failure;
- its `consulent` root test branch does not compile under either SDK because
  `defaultNamespace` and `defaultPartition` are absent;
- its module archive is not independently testable because of the relative SDK
  replacement and repository-external certificate fixtures;
- four native TLS cases fail with official Consul v1.14.2 because the tagged
  certificates expired on 2023-11-01; with a scratch-only sandbox freeport
  shim the other 213 tests pass both repeats and race with no data race; and
- all response-metadata call sites discard parsing errors, so malformed Consul
  index/cache/hash/query-backend headers are silently accepted with defaults.

The independent consumer fixture SHA-256 is
`40e09f7b5f6003da6baf17cb9aa45d60cf8319d53732559186f2aeb6e0c7a90d`.
It passes verify, count-one, two count-ten repeats, race, and vet under both
SDKs while covering configuration, request/options encoding, authentication,
token/TLS inputs, cancellation/deadline and error identity, no-retry behavior,
cleanup, malformed inputs, nil/panic boundaries, and concurrency. Pinned
apidiff finds zero incompatible changes from v1.17.0 to v1.18.0 and v1.18.0 to
v1.18.1, but five from v1.18.0 to v1.34.5.

Selected v1.18.0 has exactly one selected-version incoming edge, from Viper
v1.15.0. Its why result is negative, repository imports are zero, zero target
packages occur in the 429-entry complete project load, and it is runtime-
unreachable. Disposable exact `go get` effects retain zero load: v1.17.0
downgrades Viper; v1.18.0 manufactures eight indirect roots and 26 sum lines
before tidy removes them byte-identically to base; v1.18.1 exceeds the floor
and upgrades x/net/x/text; v1.34.5 raises the main Go declaration and moves
many unrelated selections. No projection was applied.

Fresh primary vulnerability data remains 1,398 records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z, with no Consul API record. Exact OSV
queries for all serious candidates are empty. Base and selected project scans
are identical at 30 module, 22 package, and 20 called IDs with zero target
occurrence. The fixture has zero findings. The source production closure has
inherited module-only GO-2026-5024 in old x/sys with no vulnerable package or
called symbol; tests add inherited GO-2022-0603 in yaml.v3 and load that test
package but reach no vulnerable symbol.

The decision was recorded after guard-only revalidation from clean decision
HEAD `474909583bfe6efbc5ae49b9ff1d68ec7402ee41`. Consul API remains exact
v1.18.0 through only Viper v1.15.0; all eight guarded why results remain
negative, repository source imports remain zero, and the 429-entry complete
test load contains 197 module-backed packages across 41 modules with zero
guarded-target packages. The project remains 234 selected modules and 3,599
graph edges, so Consul API and every earlier guarded target remain runtime-
unreachable. No dependency implementation was created.

Fresh primary data remains byte-identical at 1,398 records and the recorded
index SHA-256/Last-Modified values. Consul API and the six non-Gorilla earlier
targets retain empty exact OSV responses. Gorilla retains only GO-2020-0019 and
unwithdrawn GO-2026-6278; the latter's record SHA-256 remains
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
No new advisory or independent disqualifier appeared.

## Hashicorp Consul SDK Product Decision

No exact-path stable `github.com/hashicorp/consul/sdk` release qualifies under
the current Go 1.18 and behavior contracts. On 2026-09-14 the user explicitly
selected option 1 with the recommended bounds: retain exact selected,
inherited, unloaded v0.8.0 without changing dependency metadata.

The decision accepts only the completed target findings: retry nontermination
for negative count, post-deadline retry, nil-Stop panic and absent cancellation;
freeport duplicate, zero-block, and process-global-state behavior; malformed
and non-atomic iptables handling with lost wrapped error identity; writable
descriptors after TempFile cleanup; the documented test-server timeout,
context, deadline, response-body, path, log-fatal, and double-Wait defects; the
recorded API incompatibilities; inherited closure-only x/sys findings; and
related completed qualification findings. It does not accept a new or
independently discovered defect. The exception is SDK-specific and non-
transferable.

The exception remains valid only while exact v0.8.0 and its sole selected-
version incoming edge from the historical Consul API v1.12.0 graph vertex
remain unchanged, zero SDK packages load, the module remains runtime-
unreachable, and no new advisory or independent disqualifier appears.
Revalidate and record those guards in the committed decision. Direct import or
loading, runtime reachability, a target version or incoming-edge change, or a
new advisory or independent defect expires the exception and requires a fresh
SDK dependency and product decision before merge. Do not add a direct edge,
select v0.13.0 or a later release, change or remove the historical parent,
raise the Go floor, authorize dependency modernization, alter Viper, move
unrelated modules, or manufacture a dependency commit. Do not ask for this
same decision again while all guards hold. All earlier exceptions remain
separate.

Fresh proxy, sumdb, Git, and GitHub evidence resolves the module to canonical
public active unarchived non-fork `https://github.com/hashicorp/consul.git`,
subdirectory `sdk`, with MPL-2.0 candidate source. The proxy lists 30 stable
releases and five RCs; each has its exact `sdk/` Git tag. Serious tags are
lightweight and have no GitHub Release object. Strict Git verification passes,
and no redirect, fork, alternate path, major module line, prerelease, non-
versioning tag, or branch head was promoted.

Selected v0.8.0 is commit
`b1ee900870b9f04377d64936d4af452a5799b72c`, parent
`a4a43460e51d66bc562fbea476844b6c13418543`, repository tree
`8e328b47681316d07b7abbd89dff743dfbf1293d`, and SDK subtree
`dcbc915786cb1a668b9ca0f10bfd129b4526b4a8`; its proxy zip SHA-256 is
`cf29fff6c000ee67eda1b8cacec9648d06944e3cdbb80e2e22dc0165708974c6`.
Highest floor candidate v0.13.0 is commit
`e297e3e75e357c532ced05123a4e9b48f5419393`, parent
`5ce0132e9b13e9afa668ef7af91005ce2883f11f`, repository tree
`7cad367c383f5e93a20f641ad0e342d6dac9d405`, and SDK subtree
`7c45f4545bc9530db00fc92a15aaeecf84697216`.

All 30 stable declarations were inspected. V0.1.0-v0.13.0 declare Go 1.12;
v0.13.1 is the first Go 1.19 release. V0.13.0's complete imported source/test
closure declares no higher than Go 1.17 and passes under contained Go 1.18.10,
so it is the highest real floor-eligible candidate. Latest v0.18.2 declares Go
1.26.7. Selected and highest candidates each resolve 19 modules under both Go
lines. No floor-ineligible release was promoted.

Every floor-eligible stable release fails: its exported retry `Counter` stops
only when `count == Count`, so negative `Count` never terminates. Selected
v0.8.0 and v0.13.0 also reproduce retry after a timer deadline, nil `Stop`
panic, duplicate freeport returns, a zero-block allocator panic, malformed
iptables input acceptance, and a writable descriptor after `TempFile` cleanup.
Retry has no cancellation; freeport owns process-lifetime global state; Linux
iptables mutation is non-atomic and loses wrapped error identity.

Selected test-server helpers additionally ignore `ReadyTimeout`, use HTTP and
subprocesses without contexts or deadlines, leak successful service/check
response bodies, concatenate unescaped paths, and can terminate the caller
through `log.Fatal`. V0.8.0 contains a double-`Wait`/unbuffered-goroutine
cleanup defect. These behaviors were reproduced independently under both Go
lines and are exact target defects.

The source exposes `freeport`, `iptables`, `testutil`, and `testutil/retry` and
no commands, examples, benchmarks, fuzz targets, testdata, generated files,
cgo, embeds, or go:generate directives. Exact v0.8.0/v0.13.0 source passes
verify, build, vet, and Darwin/Linux/Windows production/test cross-builds under
Go 1.26.7 and Go 1.18.10. Native Darwin tests encounter only the managed
sandbox's `/usr/sbin/sysctl` denial; a scratch overlay for that environmental
probe makes count-one, two count-ten repeats, and race pass for both candidates
and both Go lines. The overlay was not applied.

Pinned API diff finds incompatible testing-interface and token changes from
v0.8.0 to v0.13.0, only a compatible `R.Logf` addition to v0.13.1, and many
incompatible retry/testing/iptables changes through v0.18.2.

Selected v0.8.0 exists only through the historical Consul API v1.12.0 graph
vertex, its sole exact selected-version incoming edge. Its why result is
negative, source imports are zero, zero target packages occur in the complete
load, and it is runtime-unreachable. Disposable exact gets keep zero load but
manufacture a direct root. V0.13.0 also adds go-version; v0.13.1 exceeds the
floor; latest raises the main Go declaration and moves unrelated selections.
The base and first three candidate tidy projections converge to selected
v0.8.0 and the same base tidy result. No projection was applied.

Fresh primary data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no SDK record. Exact OSV queries
for v0.8.0, v0.13.0, v0.13.1, and v0.18.2 are empty. Isolated older closures
contain only inherited x/sys findings: v0.8.0 has GO-2026-5024 and
GO-2022-0493; v0.13.0/v0.13.1 have GO-2026-5024. No Darwin called symbol is
reached. Project scans remain identical with zero SDK package, symbol, test-
symbol, or reachable-trace occurrence.

## Project And Quality State

No dependency metadata changed. The project remains 234 modules, 3,599 graph
edges, 429 complete-test entries, 197 module-backed packages, 41 loaded
modules, 1,067 sum lines, and the recorded 432-line tidy projection. `go.mod`
and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 archive/binary SHA-256 values are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Go 1.18.10 archive/binary values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Both ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.

Exact Go 1.26.7 verify, build, count-one, two count-ten repeats, race, vet,
golangci-lint 2.12.2, API/CLI compatibility, host and snapshot acceptance,
empty-HOME count-two, and four production cross-builds pass. Go 1.18.10 loads
366 complete-test entries; 26 unaffected packages pass two count-ten repeats
and race, and vet plus host/four cross-builds pass. The full Go 1.18 suite
retains only the two accepted Darwin `pkg/shell` closed-file wording failures.

All 15 audit meta-controls pass. A raw audit without the required external
manual-evidence receipt exits 1 by contract and is not comparable to the
accepted manual-evidence-adjusted scorecard. Because no source or dependency
metadata changed, accepted quality remains 27/27 Q0-Q2 PASS at L2, with its
recorded scorecard SHA-256
`04039eb917cc8aa674c093e866b8fbfb303ab562a28dcbaf2414dba1c9d01c3a`.

The SDK evaluation's exact-Go no-change module verification, build, count-one,
two count-ten repeats, race, vet/lint, API/CLI compatibility, host and snapshot
acceptance, empty-HOME count-two, four production cross-builds, audit meta-
controls, and launcher lifecycle gates pass. Because no source or dependency
metadata changed, no changed-selection scorecard run applies and accepted
27/27 Q0-Q2 L2 remains unchanged.

Full preflight's substantive stages pass; its nested launcher self-test hit
the known signal-retention timing race. Docker buildx is now available at
v0.33.0-desktop.1. Docker acceptance builds the image and completes the in-
image suite, then reaches the preserved Python 3.14 nanosecond timestamp parser
failure. These are target-independent harness/environment findings, not Consul
API evidence. One completed-handoff lifecycle run passes all 62 controls;
three later final-text runs pass outer controls 1-50 and reproduce only the
known nested signal-log timing failure at control 51. Preserve the known
apidiff archive reproducibility discrepancy.

The 9,780-entry external evidence manifest and decision-summary SHA-256 values
are `966263870f7529dcdd4713eba204dc1a8a00e17901a5f702557467b3c79ac46e`
and `5957399af9cf4564b9864e4b6dea8c84e8912db3846e446c953345dadbddb08a`.
All disposable evidence remains beneath `$CODEX_SESSION_SCRATCH_ROOT`; never
run `go mod download all` in a measured worktree.

## Next Bounded Objective

The next session must record the user's explicit option 1 decision without a
dependency implementation. Revalidate exact v0.8.0 and its sole historical
Consul API v1.12.0 incoming edge, negative why result, zero packages loaded
from SDK and every earlier guarded target, unchanged module hashes, runtime
unreachability, and fresh advisory state. Then answer the decision archive and
prepare exactly one next bounded P7 mission. Do not execute that successor in
the decision-recording turn.
