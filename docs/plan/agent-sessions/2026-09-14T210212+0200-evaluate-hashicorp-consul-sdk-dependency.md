# Agent Session: Evaluate Hashicorp Consul SDK Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T210212+0200-evaluate-hashicorp-consul-sdk-dependency`
Created: `2026-09-14T21:02:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ec5b93e6692d2eeaa8d72405ecf25e0f8ee5d5b84e72298dd12425e3bc3ce972`
Previous: [2026-09-14T193201+0200-decide-hashicorp-consul-api-product-direction.md](2026-09-14T193201+0200-decide-hashicorp-consul-api-product-direction.md)
Next: [2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction.md](2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction.md)
Outcome: no exact-path stable release qualifies under the Go 1.18 and behavior contracts; dependency metadata remains unchanged and P7 stops for a fresh bounded Consul SDK product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/consul/sdk v0.8.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, and all recorded retained-module
decisions through Consul API v1.18.0. All earlier outcomes and lifecycle
ancestry are final. Evaluate only Hashicorp Consul SDK in this session; do not
reopen or combine another dependency group. P8 remains queued.

The user's 2026-09-14 Consul API option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/consul/api v1.18.0` without
dependency metadata changes. It accepts only the documented non-standalone
release closure, Go 1.26 Unix native-test link incompatibility, broken
`consulent` test branch, expired TLS fixtures, silently discarded response-
metadata errors, inherited closure-only vulnerability findings, documented
nil/panic and mutation boundaries, and related completed qualification
findings. Its exception is target-specific and valid only while exact v1.18.0
and the sole selected-version Viper v1.15.0 edge remain unchanged, zero Consul
API packages load, the module remains runtime-unreachable, and no new advisory
or independent disqualifier appears. Direct import/loading, runtime
reachability, a version or incoming-edge change, or a new advisory or
independent defect requires a fresh Consul API dependency and product decision
before merge. Revalidate these guards before work and stop for that owning
decision if any fails. Do not transfer its exception to Consul SDK or change
Viper.

The Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC middleware v1.0.0, Gorilla
WebSocket v1.4.2, GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise
Certificate Proxy v0.2.1, and GAX v2.7.0 decisions remain final and separate
under their exact-selection, selected-version incoming-edge, zero-load,
runtime-unreachable, and no-new-finding guards. Revalidate those guards and
stop for the fresh owning decision if one expires. Do not change either
Gateway parent, mvn-pom-mutator, GoConvey, Viper, or transfer an earlier
exception.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Consul API and every retained
group since Google UUID have no dependency implementation or metadata commit.

The clean Consul API decision retained exact v1.18.0 through exactly the sole
selected-version incoming edge from Viper v1.15.0. Its `go mod why -m` result
remains negative, repository source imports are zero, and zero Consul API
packages occur in the complete project load. The unchanged project has 234
selected modules, 3,599 graph edges, 429 complete-test entries, 197 module-
backed packages, 41 loaded modules, 1,067 `go.sum` lines, and a 432-line
unapplied tidy projection. Exactly zero Consul API, Gateway, gRPC Prometheus,
gRPC middleware, Gorilla WebSocket, GopherJS, Enterprise Certificate Proxy, or
GAX packages load. `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains byte-identical at 1,398 module
records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact Consul API, Gateway,
gRPC Prometheus, gRPC middleware, GopherJS, Enterprise Certificate Proxy, or
GAX record. Their exact target/version OSV responses remain empty. Gorilla
retains only its recorded entries, including unwithdrawn GO-2026-6278 at
record SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.

A bounded starting survey shows selected Hashicorp Consul SDK v0.8.0 in the
current build list. The module graph has one exact selected-version incoming
edge, from the historical Consul API v1.12.0 vertex. Its `go mod why -m`
result is negative, repository imports are zero, and zero SDK packages occur
in the current complete project load. These are starting facts, not proof of
repository identity, release qualification, ancestry, floor, behavior,
vulnerability state, or suitability. Resolve them independently and do not
add a direct edge merely to alter MVS.

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
Docker timestamp control, and nested launcher signal-retention timing race;
none is Hashicorp Consul SDK evidence.

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
source-time resolution from the project's selected MVS graph and from the
distinct Consul API module's already-final decision.

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize configuration, helpers, test utilities,
network/listener behavior, timeouts, retries, concurrency, global state,
resource cleanup, error identity, nil/panic boundaries, filesystem and process
effects, and mutation/aliasing behavior. Distinguish tools, examples, optional
packages, and test-only helpers from behavior actually loaded by this project.

Add independent fixtures where useful for deterministic behavior, malformed
inputs, cancellation/deadlines, error identity, cleanup, supported
concurrency, platform branches, and selected-project consumers. Run source
verification, package listing, native complete tests, two independent repeats,
race, vet, and meaningful cross-builds under both SDKs. Classify every network,
clock, timeout, filesystem, subprocess, platform, resource, global-state, or
test-design failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. Compare primary vulnerability results at module, package,
symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, every unchanged exception guard, and
`./codex-dev-start.sh --check`. Read this archive, the answered Consul API
decision and evaluation, the earlier guarded-decision archives, rolling
handover, roadmap, `go.mod`, `go.sum`, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose
complete minimal source/test closure preserves Go 1.18. Do not promote an
unqualified or floor-ineligible identity. If no candidate satisfies the
existing contracts, preserve the evidence and stop for a fresh bounded product
decision.

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

After the Hashicorp Consul SDK decision or bounded product stop, rewrite the
rolling handover and roadmap, answer this archive, and follow the repository
lifecycle contract. Do not implement another dependency group, launch a
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable `github.com/hashicorp/consul/sdk` release qualifies under
the existing Go 1.18 floor and behavior contracts. The highest declaration-
and closure-eligible stable release is v0.13.0; v0.13.1 is the first release
declaring Go 1.19. Every stable floor-eligible release from v0.1.0 through
v0.13.0 retains an exported retry counter whose negative `Count` never reaches
its equality-only stop condition. Independent fixtures reproduce this and
additional material defects in selected v0.8.0 and v0.13.0. The evaluation
therefore stops for a fresh bounded product decision. `go.mod` and `go.sum`
were not changed, no direct target edge was added, and no dependency
implementation commit was manufactured.

The evaluation began from clean feature-branch HEAD
`19e0f4a2e36c654c9e1810b8f660cec2a29f22fb`, parent
`474909583bfe6efbc5ae49b9ff1d68ec7402ee41`, and tree
`e1d077544884d461318861fb1d6aad72e8caa202`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. The launcher check, reciprocal
archive history, roadmap state, ordinary and ignored status, and all earlier
exception guards were valid before qualification.

Fresh Go proxy, checksum database, Git, and GitHub evidence resolves the exact
module to the public, active, unarchived, non-fork Consul repository, `sdk`
subdirectory, and MPL-2.0 candidate source. The proxy lists 35 versions: 30
stable releases and five release candidates. Every proxy version has its exact
`sdk/` Git tag. Serious tags are lightweight, have no GitHub Release object,
and are not signed tag objects. Strict Git object verification passes. No
retraction, deprecation, redirect, fork, alternate module path, different major
line, prerelease, non-versioning tag, or branch head was promoted.

Selected v0.8.0 is tag `sdk/v0.8.0`, commit
`b1ee900870b9f04377d64936d4af452a5799b72c`, parent
`a4a43460e51d66bc562fbea476844b6c13418543`, repository tree
`8e328b47681316d07b7abbd89dff743dfbf1293d`, and SDK subtree
`dcbc915786cb1a668b9ca0f10bfd129b4526b4a8`, dated 2021-06-22. Its
proxy zip SHA-256 is
`cf29fff6c000ee67eda1b8cacec9648d06944e3cdbb80e2e22dc0165708974c6`;
sumdb source/mod sums are
`h1:OJtKBtEjboEZvG6AOUdh4Z1Zbyu0WcxQ0qatRrZHTVU=` /
`h1:GBvyrGALthsZObzUGsfgHZQDXjg4lOjagTIwIR1vPms=`.

Highest floor candidate v0.13.0 is tag `sdk/v0.13.0`, commit
`e297e3e75e357c532ced05123a4e9b48f5419393`, parent
`5ce0132e9b13e9afa668ef7af91005ce2883f11f`, repository tree
`7cad367c383f5e93a20f641ad0e342d6dac9d405`, and SDK subtree
`7c45f4545bc9530db00fc92a15aaeecf84697216`, dated 2022-11-21. Its
proxy zip SHA-256 is
`681fd9081e6b3266eda9f5c6ed3476f6fbf5e0f2a9edc472be96bdcb4ca363ab`.
First floor-ineligible v0.13.1 is commit
`983a1b8ddb5ce02192aa58bb2d451a5e4140f588`, parent
`a6180659d0f3d3c3fd4a74de4b98af037dda1afd`, tree
`6b7903c82e59813a134057b5b9348f453180b075`, SDK subtree
`d1349eaa7867eb292b62285cfc7da5c73737eb30`, and proxy zip SHA-256
`a92f901c29837740399220948e4e86dc5802c67c1568543c59a24f2a7941f690`.
Latest v0.18.2 is commit
`c115f1d08666ccd319ffda35f28adea5e162b7ab`, parent
`78ab247c8551f080c3d73ad5b7b6c381ee871465`, tree
`39d5a9d84f3ca63b3eb6bf31bcd52f10282ddb8a`, SDK subtree
`549857b8db6b39127b522a985bef9fa7f2801ee7`, and proxy zip SHA-256
`29c0f0952ed29f082438b5326a4f41208790a87c706b2d8b76e7211e004c5976`.
V0.8.0/v0.13.0 commits are unverified; v0.13.1/v0.18.2 commits have GitHub-
verified signatures. Release-pair merge bases show separate backport histories
rather than linear tag ancestry. Proxy and Git regular SDK source agree for
every serious candidate; Go module packaging omits only a repository symlink
and provides the module license in the normal proxy form. The GitHub exact-path
`?go-get=1` request returns an ordinary 404 with no `go-import` metadata and no
redirect; the proxy VCS origin and exact subdirectory tags establish identity.
GitHub's repository license field is `NOASSERTION`; the current repository root
uses BSL while the tagged SDK module source carries MPL-2.0.

All 30 stable go.mod declarations were inspected. V0.1.0 through v0.13.0
declare Go 1.12, v0.13.1 through v0.16.1 declare Go 1.19, and later releases
rise through Go 1.22 and Go 1.25 to Go 1.26.7. Complete imported production
and test closure was resolved under exact Go 1.26.7 and contained Go 1.18.10.
Selected v0.8.0 and v0.13.0 each resolve 19 modules; their imported closures
declare no higher than Go 1.13 and Go 1.17 respectively. Go 1.26.7 resolves
207/231 production/test entries for v0.8.0 and 208/232 for v0.13.0; Go 1.18.10
resolves 145/168 and 146/169. Thus v0.13.0 really preserves Go 1.18 across
source and tests, while v0.13.1 does not.

Selected and highest-candidate source each expose four packages: `freeport`,
`iptables`, `testutil`, and `testutil/retry`. There are no commands, examples,
benchmarks, fuzz targets, testdata trees, generated files, cgo, embeds, or
go:generate directives. Darwin, Linux, Windows, fallback Unix, and unsupported
iptables branches were inspected. The module README describes the SDK as
internal to Consul, public only for sharing, unsuitable for new consumers, and
free to change its 0.x APIs.

Material reproduced behavior includes:

- negative retry `Counter` limits continue indefinitely in every stable
  floor-eligible release; a timer may sleep beyond its deadline and then permit
  another attempt, retry has no cancellation boundary, `R.Stop(nil)` panics,
  and retry instances are mutable, non-reusable, and not concurrency-safe;
- freeport reserves only a block-selection listener, not returned ports; uses
  process-global random/state and a lifetime listener/goroutine; may wait
  forever; accepts duplicate or never-taken returns; and panics when a
  zero-block derived range reaches `rand.Int31n(0)`;
- iptables accepts malformed UIDs, negative ports, and invalid CIDRs; Linux
  mutation is non-atomic, subprocesses have no context or deadline, partial
  rules survive a later failure, and `%v` wrapping discards error identity;
- `TempFile` leaves its returned descriptor open after cleanup, selected
  testutil creates `/tmp/consul-test` during package initialization, cleanup
  errors are ignored, and global environment/logging state is captured once;
- test-server helpers use real subprocesses, network, and clocks. Selected
  v0.8.0 ignores `ReadyTimeout`, waits a fixed two seconds, has no HTTP client
  timeout, contains a double-`Wait`/unbuffered-goroutine cleanup defect, and can
  terminate the caller through `log.Fatal`. Successful service/check responses
  are discarded without closing their bodies, paths are concatenated without
  escaping, and the exposed request methods have no context or timeout.

Independent external and same-package fixtures reproduce 1,024 continuing
negative-counter iterations, a retry approved after its deadline, the nil
`Stop` panic, three unclosed successful response bodies, process exit from a
nil server receiver, malformed iptables rule acceptance, a writable file
descriptor after pathname cleanup, duplicate freeport returns, and the
zero-block allocator panic under both Go lines. Native upstream Darwin tests
hit only the managed sandbox's denial of `/usr/sbin/sysctl`; a scratch-only
overlay replaced that environmental probe. With the overlay, both v0.8.0 and
v0.13.0 pass count-one, two independent count-ten repeats, and race under both
Go lines. Exact source passes verify, build, and vet, and production plus test
cross-compilation passes darwin/amd64, linux/amd64, linux/arm64, and
windows/amd64. The overlay is not dependency evidence and was not applied.

Pinned API comparison finds an incompatible testing-interface type change and
removed token fields from v0.8.0 to v0.13.0, only the compatible `R.Logf`
addition from v0.13.0 to v0.13.1, and numerous incompatible retry, testing, and
iptables changes through v0.18.2. Export snapshots contain 273, 197, 209,
and 271 lines respectively.

The current MVS selection exists only because the historical
`github.com/hashicorp/consul/api v1.12.0` graph vertex requires SDK v0.8.0.
That is the sole exact selected-version incoming edge. `go mod why -m` is
negative, repository Go imports are zero, zero SDK packages occur in the
429-entry complete load, and the target is runtime-unreachable. Disposable
exact gets retain zero target load but manufacture a direct indirect root.
V0.8.0 produces 234 modules/3,600 edges/1,068 sum lines and changes only that
root/checksum. V0.13.0 produces 235/3,612/1,071 and additionally adds
go-version. V0.13.1 produces 235/3,615/1,069 but exceeds the floor. Latest
v0.18.2 raises the main Go declaration to 1.26.7, removes the toolchain
directive, moves numerous unrelated selections, and produces
235/3,619/1,080. Their go.mod/go.sum SHA-256 pairs are respectively
`f7df814b24a9549d4c38ba3f055ddf764b31b73cc209789e8cce1fbcfeec762b` /
`5d6feffbc6dd011550de36b8e6853d19a5d2e9d2905163dd84118eb0e7c921f2`,
`62f78358107904330bad12e3807eb19a7625e99dab9c22386dd6e03b42a44a41` /
`f7d8ec7262950f9ee40689dcd06410d2cee7d5c040fc765f468259972c852fd3`,
`4702be37b647684d4594d1a48bf8192631db33fbde9f44aeb7b559f4859f0a68` /
`d96eee1dc9def239bb737145bad859d9d4131b4a358eafadbb77414014852633`,
and `414dc5b6b7972f136fddadd7d791d49c6f8f196eb657d4a877406e32f2c96e94` /
`28a10a594318af661b7729eab6e6d3dc88da409c8447a95d9e62829b04b686a5`.
Base and v0.8.0/v0.13.0/v0.13.1 tidy projections all converge to selected
v0.8.0, 234 modules, 3,557 pruned graph edges, 948 sum lines, and hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Latest tidy restores SDK v0.8.0 but retains its Go-floor and unrelated changes.
None was applied because a direct edge would manufacture activity without a
loaded consumer.

Fresh primary vulnerability data remains 1,398 records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact Consul SDK record, and
exact OSV queries for v0.8.0, v0.13.0, v0.13.1, and v0.18.2 are empty. The
isolated v0.8.0 closure has module-level GO-2026-5024 and GO-2022-0493 in old
x/sys; its Darwin package scan retains GO-2022-0493, while symbol and test-
symbol scans are empty. V0.13.0/v0.13.1 retain only module-level GO-2026-5024
and have empty Darwin package, symbol, and test-symbol scans. Latest is empty.
Baseline and every candidate project scan are identical at 30 module IDs, 22
loaded-package IDs, and 20 called IDs/22 traces, with zero target package,
symbol, test-symbol, or reachable-trace occurrence. These inherited source-
closure findings are distinct from the project's selected x/sys v0.30.0 and
are not exact SDK advisories.

All earlier guarded selections, incoming edges, negative why results, zero
loads, runtime-unreachable status, and advisory state remain unchanged. The
project remains 234 modules, 3,599 graph edges, 429 complete-test entries, 197
module-backed packages across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 verify, build, count-one, two count-ten repeats, race, vet,
golangci-lint 2.12.2, API/CLI compatibility, host and snapshot acceptance,
empty-HOME count-two, four production cross-builds, audit meta-controls, and
the 62-control launcher suite pass. No changed-selection scorecard applies;
accepted quality remains 27/27 Q0-Q2 PASS at L2. Portable golangci-lint and
GoReleaser receipts match the required hashes. The known apidiff archive,
Python 3.14 Docker timestamp, and nested launcher signal-timing findings remain
target-independent.
