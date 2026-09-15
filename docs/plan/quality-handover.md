# Quality Upgrade Handover

Generated: 2026-09-15T23:43:09+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`. The go-hclog evaluation
  began from clean handoff HEAD
  `452fffbd865dcae811fba3b02821cf7fd48f774c`, parent
  `0d1d66480604da5b646c75593570b01ebec72dda`, tree
  `4858e0dd3c216b772476058324a265aef0cb005a`.
- The latest dependency implementation remains exact Google UUID v1.4.0
  commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. Consul API, Gateway,
  and every retained group since Google UUID have no dependency implementation
  or metadata commit.
- The go-hclog evaluation archive and every earlier archive are answered. The
  sole NEXT archive is `docs/plan/agent-sessions/2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction.md`. It authorizes only the bounded go-hclog
  product decision, requires one explicit user choice, and does not authorize
  implementation, another dependency group, or P8.
- No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not launch a successor, push,
  merge, publish, release, stash, revert, bypass cleanup, or remove the
  worktree.

## Roadmap And Guarded Decisions

P2A-P6 are complete. P7 is stopped at one bounded Hashicorp go-hclog product
decision because no exact-path stable release qualifies. Exact Go 1.26.7,
every accepted dependency move through Google UUID v1.4.0, qualified
go-cleanhttp, and all retained-module decisions through exact inherited,
unloaded Errwrap v1.0.0 are final under their target-specific guards. P8
remains queued. Do not infer the go-hclog choice, combine groups, or begin P8.

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
Viper; reopen an earlier choice; or transfer any exception between targets.

## Hashicorp go-hclog Evaluation And Decision Boundary

No exact-path stable `github.com/hashicorp/go-hclog` release qualifies under
the current behavior and compatibility contracts. The project remains
byte-for-byte unchanged with inherited, unloaded v1.2.0 selected through the
sole selected-version Viper v1.15.0 edge, but that physical state is not a
qualification or accepted exception. P7 now stops for the bounded decision in
the sole NEXT archive.

Fresh identity evidence resolves canonical public active unarchived non-fork
MIT repository `https://github.com/hashicorp/go-hclog.git`. The proxy lists 31
stable releases from v0.7.0 through proxy `@latest` v1.6.3, with no retraction,
module deprecation, redirect, or `/v2` line. Non-semver alias tags, nested
`hclogvet` tags, branch `f-v2`, and unreleased Go-1.25 main were not promoted.
Selected v1.2.0 is lightweight unsigned commit
`b6b55671f4e5b82443139ee3e9f4417603c4cd72`, tree
`45f0da0a3b7eb001522fc9891504cee03d951ef0`; latest v1.6.3 is commit
`d12136aa2e51933c460084f5083b6d5bb9d41960`, tree
`26c5c9247ad5b3a6f930ea2b1fc0a1ee532a7804`. Strict Git and proxy/Git byte
verification pass.

Every stable declares no more than Go 1.13. Imported selected/latest complete
production/test closures declare no more than Go 1.17 and pass contained Go
1.18.10. Serious releases contain one package, 12 production files, seven
tests, Unix/Windows color branches, and one benchmark; no command, example,
fuzz target, testdata, generated file, cgo, embed, generator, or symlink exists.
Selected v1.2.0 through last-compatible v1.3.1, first-breaking v1.4.0, and
latest v1.6.3 pass verification, native repeated/race/vet tests, and applicable
cross-builds under both SDKs.

Pinned API comparison finds v1.2.1/v1.2.2 identical to selected and
v1.3.0/v1.3.1 compatible with only `LoggerOptions.ColorHeaderAndFields` added.
V1.4.0 adds `Logger.GetLevel` to the public interface, an incompatible change
retained through latest. V1.3.1 is the last API-compatible candidate.

The 183-line independent fixture SHA-256 is
`b62af8023aa4e9d2b84a20db69627992998bfa29e7c268832025b345dcf5bb26`.
It directly proves five blockers in selected and every API-compatible patch:
NaN silently drops a JSON record; caller `@message`/`@level` fields overwrite
core metadata; absent-sink deregistration underflows the count and disables a
later real sink; a self-deregistering sink deadlocks under the held registry
mutex; and `SetDefault(nil)` violates the documented non-nil `FromContext`
result. Source history preserves these implementations through latest. The
latest-only 201-line extension SHA-256
`15e4bc71cde744a69ee0cecc18cea1ec6730acb595a87f11ea3de6d92d2ebd02`
also proves a `SyncParentLevel` concurrent level/epoch race under both SDKs.
The fixtures separately characterize deterministic formatting/routing, exact
reset errors, nil/panic and aliasing boundaries, allocations, immutable
concurrency, and caller-owned output cleanup.

Selecting v1.3.1 retains all five blockers. Selecting v1.4.0+ also breaks the
public interface. Selecting v1.1.0 or lower forces an unauthorized Viper
change; exact v1.1.0 downgrades Viper to v1.10.1. No candidate was applied.

Selected v1.2.0 has one selected-version incoming edge from Viper v1.15.0.
Historical Viper v1.10.1 and `sagikazarmark/crypt v0.4.0` request v1.0.0;
historical Consul API v1.12.0 and Consul SDK v0.8.0 request v0.12.0. The why
result is negative, source imports are zero, and production/complete-test loads
contain zero target packages, so the module is runtime-unreachable. Disposable
v1.2.0/v1.3.1/v1.4.0/v1.6.3 gets preserve 234 modules, every unrelated
selection, 429 complete-test entries, and zero target load; they manufacture a
direct root and checksums. Tidy restores selected v1.2.0 and the exact base
projection. Nothing was applied.

Fresh primary data remains 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z with no go-hclog record. Exact serious
OSV responses are empty. Selected has zero Go-1.26 source findings. V1.3.1+
has only module-level GO-2026-5024 through old `x/sys`; the affected Windows
package/symbol is not imported or called, and no target trace exists. Old-SDK
source traces call vulnerable Go 1.18 standard library but assign no advisory
to go-hclog. Project scans retain 30 IDs and zero target SBOM package, symbol,
test-symbol, or reachable trace.

The project remains 234 modules, 3,599 edges, 355 production entries, 429
complete-test entries, 197 module-backed entries across 41 loaded modules,
1,067 sum lines, and the 432-line tidy projection. `go.mod`/`go.sum` SHA-256
remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Applicable project, API/CLI, lint, repeated/race/vet, empty-HOME, cross-build,
host/snapshot acceptance, all 17 script meta-tests, 80/80 mutation kills, all
15 audit controls, and an independent 62-control lifecycle run pass. Contained
Go 1.18 retains only the two accepted Darwin shell wording failures. Full
preflight alone reproduced the known nested signal-retention timing race.
Accepted quality remains 27/27 Q0-Q2 PASS at L2. No dependency implementation
or metadata commit was created. The 1,328-entry disposable evidence manifest
SHA-256 is
`e6647301b469e2c915ca07bd1f2b0f7856e6b55df75c76a350943cc780783a9f`.

The next session must present exactly three choices: recommended guarded
retention of exact inherited unloaded v1.2.0 with the completed findings;
a separately scoped remediation study without implementation authority; or a
P7 block. It must not infer a choice, qualify v1.3.1, promote v1.4.0+, downgrade
Viper, add a direct edge, or implement a patch/fork/parent change.

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
Direct import or loading, runtime reachability, a target version or incoming-
edge change, or a new advisory or independent defect expires the exception and
requires a fresh SDK dependency and product decision before merge. Do not add
a direct edge, select v0.13.0 or a later release, change or remove the
historical parent, raise the Go floor, authorize dependency modernization,
alter Viper, move unrelated modules, or manufacture a dependency commit. Do
not ask for this same decision again while all guards hold. All earlier
exceptions remain separate.

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

The decision was recorded after guard-only revalidation from clean decision
HEAD `2346fcecd1f9f464697d51c3fa223b6709e84972`. Consul SDK remains exact
v0.8.0 through only the historical Consul API v1.12.0 graph edge; all nine
guarded why results remain negative, repository source imports remain zero,
and the 429-entry complete test load contains 197 module-backed packages across
41 modules with zero guarded-target packages. The production load also has
zero guarded-target packages. The project remains 234 selected modules and
3,599 graph edges, so SDK and every earlier guarded target remain runtime-
unreachable. No dependency implementation or metadata change was created.

Fresh primary data remains byte-identical at 1,398 records and the recorded
index SHA-256/Last-Modified values. SDK and every non-Gorilla guarded target
retain no exact module record; the exact SDK v0.8.0 OSV result remains empty.
Gorilla retains only its recorded module entries, including unwithdrawn
GO-2026-6278; that record's SHA-256 remains
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
No new advisory or independent disqualifier appeared.

## Hashicorp Errwrap Product Decision

On 2026-09-15 the user explicitly selected option 1 with the recommended
bounds: retain exact selected, inherited, unloaded
`github.com/hashicorp/errwrap v1.0.0` without changing `go.mod` or `go.sum`.
Accept only the completed concrete-type collision, absent standard single/
multi-error traversal, nil/panic, aliasing, allocation, recursion, API, and
related qualification findings. This does not accept a new or independently
discovered defect. The exception is Errwrap-specific and non-transferable.

It remains valid only while exact v1.0.0 and its sole selected-version incoming
edge from `github.com/hashicorp/go-multierror v1.1.0` remain unchanged, zero
Errwrap packages load, the module remains runtime-unreachable, and no new
advisory or independent defect appears. Direct import/loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent defect
expires the exception and requires a fresh Errwrap dependency and product
decision before merge. Do not add a direct edge, select v1.1.0 or another
version, remove or change the historical parent chain, change go-multierror,
Serf, Viper or mvn-pom-mutator, raise the Go floor, authorize replacement or
modernization, move unrelated selections, or manufacture a dependency commit.
Do not ask for this same choice again while all guards hold. All earlier
exceptions remain separate.

No exact-path stable release qualifies. The proxy exposes only v1.0.0 and
v1.1.0. Both have stdlib-only complete source/test closures with no Go
directive and pass under contained Go 1.18.10, so v1.1.0 is the highest floor-
eligible stable candidate. Canonical metadata resolves public active
unarchived non-fork `https://github.com/hashicorp/errwrap`, MPL-2.0, no GitHub
Release objects, no retractions or module deprecation, and no `/v2` line.
Unreleased master now declares Go 1.24 and was not promoted.

Selected v1.0.0 is lightweight tag/commit
`8a6fb523712970c966eefc6b39ed2c5e74880354`, parent
`d6c0cd88035724dd42e0f335ae30161c20575ecc`, tree
`9863613ad8fe960290d1631d84ede660af5f0746`, dated
2018-08-24T00:39:10Z, with proxy zip SHA-256
`ccdf4c90f894d8a5fde4e79d5828c5d27a13e9f7ce3006dd72ce76e6e17cdeb2`.
Its commit contains a PGP payload but GitHub reports `unknown_key`. V1.1.0 is
lightweight tag/verified merge commit
`7b00e5db719c64d14dd0caaacbd13e76254d02c0`, parents v1.0.0 and
`96a78ad11c51762df122b738e6f4f30f58e03d8d`, tree
`aefd62cf8e9549e154a65afb8b4524b7b6ed5f2e`, dated
2020-07-14T15:51:01Z, with proxy zip SHA-256
`209ae99bc039443e28e4d6bb66517d1756d9468b7578d31f1b63a28103d8e18c`.
V1.0.0 is its ancestor. Proxy and Git source are byte-identical and strict Git
verification passes. Sumdb source sums are
`h1:hLrqtXhLNNiApEnnlnZBrrPHi3JYFw4JYHsiTfAzrGc=` and
`h1:OxrO0bFjS2NK6ZaMyfRMIaByCoOUqmfzLLXGmPAQZ7w=`; both use mod sum
`h1:YH+1FKiLX6mgECZlMWQKBywILFU/P0P8Iy1URd86jBQ=`.

Each release contains exactly LICENSE, README, go.mod, one production Go file,
and one test file: one package, no commands, examples, benchmarks, fuzz
targets, testdata, generated files, build tags, platform branches, cgo, embeds,
go:generate directives, symlinks, or external dependencies. Under exact Go
1.26.7 and contained Go 1.18.10 both pass module verification, build, native
count-one, two count-ten repeats, race, vet, and Darwin amd64, Linux
amd64/arm64, Windows amd64, and js/wasm production/test cross-compilation.
Pinned API diff finds the same exported `Contains`, `ContainsType`, `Get`,
`GetAll`, `GetAllType`, `GetType`, `Walk`, `Wrap`, and `Wrapf` functions plus
`WalkFunc` and `Wrapper` types in both releases. V1.1.0 adds the unlisted
standard `Unwrap() error` method and deprecates `Wrapf`; public declarations
have zero forward or reverse incompatibilities.

Both stable candidates independently violate behavior contracts. Their
exported type matching compares `reflect.Type.String()`, so distinct concrete
types in different import paths with the same package/type spelling falsely
match. Both also fail to walk the standard `Unwrap() []error` graph, including
`errors.Join`, under Go 1.26.7; v1.0.0 additionally cannot interoperate through
standard single-error unwrapping. Independent fixtures reproduce the concrete-
type defect under both Go lines and the multi-error defect where supported.
They also characterize deterministic formatting/order and identity, nil
lookup, explicit nil outer-error and nil callback panics, error aliasing,
isolated `WrappedErrors` result slices, custom child order, supported immutable
concurrent reads, and allocation. No global state or resource ownership exists;
recursive walking is unguarded against cycles.

Selected v1.0.0 has exactly one selected-version incoming edge from
`github.com/hashicorp/go-multierror v1.1.0`; a historical multierror v1.0.0
edge also requests it. The shortest project graph path is main -> mvn-pom-
mutator v0.2.3 -> historical Viper v1.10.1 -> Serf v0.9.6 -> go-multierror
v1.1.0 -> Errwrap v1.0.0. Its why result is negative, repository imports are
zero, and zero target packages occur in production or complete-test loads, so
it is runtime-unreachable. Disposable exact v1.0.0 and v1.1.0 gets preserve
234 modules, 429 complete-test entries, zero load, and every unrelated
selection; they add only a direct root/one graph edge and respectively one or
two checksum lines. All tidy projections converge to selected v1.0.0 and the
base projection. Nothing was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no Errwrap record. Exact candidate
OSV responses are empty. Go 1.26.7 direct source scans report no module,
package, symbol, or test-symbol finding. Go 1.18.10 source scans report only
that legacy SDK's standard-library findings, not Errwrap. Base and v1.1.0
project reports are byte-identical in all scan modes with zero target package,
symbol, test-symbol, or reachable trace.

The project remains unchanged. The fixture file-list receipt SHA-256 is
`6eb0df6dedb99aaf6d00c9e3b9b9771543ed09f394aa8f509c0184eafc884765`;
the 322-entry disposable evidence manifest SHA-256 is
`181ef08a5ce2f8132d9bab9b080bb5cdfb53aedef14f0e87f0ad78e152853cc8`.
The decision was recorded from clean decision HEAD
`f77cb8a75c6813b28656c85b654d17d7df464731`. Exact v1.0.0 retains its sole
selected-version go-multierror v1.1.0 incoming edge; the historical
go-multierror v1.0.0 request also remains. All ten guarded versions and
recorded incoming edges remain unchanged. All ten why results remain
negative, repository Go source has zero guarded-path occurrences, and
production and 429-entry complete-test loads contain zero guarded packages.
The complete-test load retains 197 module-backed entries across 41 modules,
so all guarded modules remain runtime-unreachable.

Fresh primary vulnerability data now has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. It has no Errwrap or other new
guarded-target record. Exact Errwrap v1.0.0/v1.1.0 OSV results remain empty.
Gorilla retains only GO-2020-0019 and unwithdrawn GO-2026-6278; that record's
SHA-256 remains
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
No new advisory or independently observed defect appeared.

No dependency implementation or metadata change exists. Parent-edge removal
and broader replacement remain unauthorized, and no earlier exception
transfers.

## Hashicorp Go Cleanhttp Decision

Retain qualified exact selected `github.com/hashicorp/go-cleanhttp v0.5.2`
without changing `go.mod` or `go.sum`. It is the highest exact-path stable
release, its complete production/test closure preserves Go 1.18, and its
source, API, behavior, tests, vulnerability evidence, and project effects pass
the existing contracts. No exception or product-risk acceptance was required.

Fresh canonical evidence resolves public active unarchived non-fork
`https://github.com/hashicorp/go-cleanhttp.git`, MPL-2.0, and exactly stable
v0.5.0, v0.5.1, and v0.5.2. V0.5.2 is both proxy `@latest` and the sole GitHub
Release. There is no prerelease, `/v2` line, alternate path, redirect,
retraction, or deprecation. Unreleased master declares Go 1.24 and was not
promoted. The lightweight unsigned tags form linear ancestry through master.
Selected v0.5.2 is commit
`6d9e2ac5d828e5f8594b97f88c4bde14a67bb6d2`, parent
`d3fcbee8e1810ecee4bdbf415f42f84cfd0e3361`, tree
`c1160f09cedce00dc3ef7b06169ac45cad1c12e8`, dated
2021-02-03T18:51:13Z. Its proxy zip SHA-256 is
`e9f3dcfcb33172ba499b4f8e888169252d7f1e072082182124a6e2053523f7df`,
and proxy and Git bytes agree.

Every release is one standard-library-only package with three production Go
files and one test file; no commands, examples, benchmarks, fuzz targets,
testdata, generated files, build tags, platform branches, cgo, embeds,
go:generate directives, or symlinks exist. V0.5.2 declares Go 1.13. Under
exact Go 1.26.7 it resolves 184 production and 209 complete-test entries;
contained Go 1.18.10 resolves 123 and 147. It passes verify, build, count-one,
two count-ten repeats, race, vet, and six production/test cross-builds under
both SDKs. V0.5.0/v0.5.1 pass build/vet/cross-builds but their own native
tests consistently fail because they construct request URLs with raw control
bytes rejected by current Go; v0.5.2 percent-encodes those test inputs.

Pinned API exports are identical across all three releases. Independent
fixtures cover exact transport defaults and fresh identity, client/transport
aliasing, global state, printable and control paths, Unicode and invalid UTF-8,
HTTP/1 transient versus pooled reuse, idle cleanup, HTTP/2, allocations, and
supported immutable concurrency under both SDKs. A zero `HandlerInput.ErrStatus`
is mutated to 400; the handler retains the caller pointer; later mutation
changes behavior and concurrent mutation is unsupported. Nil request and nil
next are no-ops. A non-nil request with nil URL and invalid nonzero status
panic through standard-library preconditions. Pooled users own idle-resource
cleanup. The fixture SHA-256 is
`1b3aa948358fdbfd136041fc6c930ba904ff8bc6260bd915ea9af50f40dd73a2`.

Selected v0.5.2 has three selected-version incoming edges: Viper v1.15.0,
historical Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`. The shortest path
is main -> Viper v1.15.0 -> go-cleanhttp. Its why result is negative, source
imports are zero, production and complete-test target loads are zero, and it
is runtime-unreachable. A disposable exact v0.5.2 get adds only a direct
indirect root edge and one checksum line while preserving all 234 modules,
429 complete-test entries, every unrelated version, and zero target load; tidy
returns to the base projection. V0.5.0/v0.5.1 exact gets instead force broad
unrelated downgrades, remove mvn-pom-mutator, and make the project unloadable.
No projection was applied.

Fresh 1,399-record primary vulnerability data has no go-cleanhttp record, and
all three exact OSV responses are empty. Go 1.26.7 isolated scans have no
findings. Go 1.18.10 reports only its old standard-library findings; target
frames call into that library but no advisory is assigned to go-cleanhttp.
Base and disposable v0.5.2 project scans have identical normalized finding
populations and zero target SBOM package, symbol, test-symbol, or reachable
trace. The 412-entry evaluation manifest SHA-256 is
`40b0a6fc11f4378aed62586bfc8bdd94bef33b0a0ba88051dff51912807137c1`.

This qualified result remains final while exact v0.5.2 and all three incoming
edges remain unchanged, zero target packages load, it stays runtime-
unreachable, and no new advisory or independently disqualifying behavior
appears. Direct import/loading, runtime reachability, a target version or
incoming-edge change, or a new advisory/defect requires a fresh bounded
go-cleanhttp decision. This guard transfers no earlier exception.

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

The Errwrap decision-recording session's exact Go 1.26.7 module verification,
build, count-one tests, race tests, vet, and launcher lifecycle check pass. No
changed-selection scorecard applies, so accepted quality remains 27/27 Q0-Q2
PASS at L2.

The Errwrap evaluation's exact Go 1.26.7 module verification, build,
count-one, two independent count-ten repeats, race, vet, pinned lint, API/CLI
compatibility, host and snapshot acceptance, empty-HOME count-two, and four
production cross-builds pass. Contained Go 1.18.10 loads 366 complete-test
entries; 26 unaffected packages pass both count-ten repeats, race, and vet,
while the full suite retains only the two accepted Darwin shell wording
failures. All 15 audit meta-controls pass. The raw audit without external
manual evidence is deliberately non-comparable; accepted quality remains
27/27 Q0-Q2 PASS at L2. One preflight launcher invocation reproduced the known
signal/log-retention timing race after all earlier stages passed; an immediate
independent rerun passed all 62 controls. No changed-selection scorecard
applies.

The go-cleanhttp evaluation's exact Go 1.26.7 module verification, build,
count-one, two independent count-ten repeats, race, vet, pinned lint, API/CLI
compatibility, host and snapshot acceptance, empty-HOME count-two, and four
production cross-builds pass. All 17 production-script meta-tests and all 15
quality-audit meta-controls pass. Contained Go 1.18.10 loads 366 complete-test
entries; its 26 unaffected packages pass count-one, both count-ten repeats,
race, and vet, while the full suite retains only the two accepted Darwin shell
wording failures. Four production cross-builds pass under Go 1.18.10. No
changed-selection scorecard applies, so accepted quality remains 27/27 Q0-Q2
PASS at L2.

Full preflight repeatedly reached passing API/CLI compatibility, build,
count-one, vet, and pinned lint before reproducing the known outer or nested
launcher signal/log-retention timing race. An independent launcher run passed
all 62 controls on its sixth timing attempt, and the final post-handoff run
passed all 62 on its second attempt. One legacy meta-test's bare
`mktemp -d` selected an unwritable managed macOS temp root; a scratch-only
command wrapper let the unchanged contract pass. These are known harness and
managed-environment boundaries, not go-cleanhttp findings.

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

Evaluate only selected exact-path `github.com/hashicorp/go-hclog v1.2.0` as
one bounded P7 dependency group. Resolve its repository and release identity,
complete Go-floor closure, behavior and API, actual project loading, MVS
effects, vulnerability state, and applicable quality contracts before
retaining or changing it. Preserve qualified go-cleanhttp v0.5.2, the Errwrap
decision, and every earlier separate exception. Do not combine another
dependency group or begin P8.
