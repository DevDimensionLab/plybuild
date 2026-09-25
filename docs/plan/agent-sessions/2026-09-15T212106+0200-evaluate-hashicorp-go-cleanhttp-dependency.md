# Agent Session: Evaluate Hashicorp Go Cleanhttp Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-15T212106+0200-evaluate-hashicorp-go-cleanhttp-dependency`
Created: `2026-09-15T21:21:06+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `bacb220e23cce70fa12552bb23368530ba4de43661939eab6e54a2560377b37e`
Previous: [2026-09-14T234716+0200-decide-hashicorp-errwrap-product-direction.md](2026-09-14T234716+0200-decide-hashicorp-errwrap-product-direction.md)
Next: [2026-09-15T222945+0200-evaluate-hashicorp-go-hclog-dependency.md](2026-09-15T222945+0200-evaluate-hashicorp-go-hclog-dependency.md)
Outcome: retained qualified exact v0.5.2; no dependency metadata changed

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-cleanhttp v0.5.2` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, and all recorded retained-module
decisions through Errwrap v1.0.0. All earlier outcomes and lifecycle ancestry
are final. Evaluate only Hashicorp go-cleanhttp in this session; do not reopen
or combine another dependency group. P8 remains queued.

The user's 2026-09-15 Errwrap option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/errwrap v1.0.0` without dependency
metadata changes. It accepts only the completed concrete-type collision,
absent standard single/multi-error traversal, nil/panic, aliasing, allocation,
recursion, API, and related qualification findings. Its exception is Errwrap-
specific and valid only while exact v1.0.0 and the sole selected-version
go-multierror v1.1.0 incoming edge remain unchanged, zero Errwrap packages
load, the module remains runtime-unreachable, and no new advisory or
independent defect appears. Direct import/loading, runtime reachability, a
target version or incoming-edge change, or a new advisory or independent
defect requires a fresh Errwrap dependency and product decision before merge.
Revalidate these guards before work and stop for that owning decision if any
fails. Do not transfer its exception to go-cleanhttp or change its parent
chain.

The Consul SDK v0.8.0, Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus
v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain final and separate under their exact-
selection, recorded incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards and stop for the fresh owning decision
if one expires. Do not change either Gateway parent, mvn-pom-mutator,
GoConvey, Viper, the historical Consul API parent, go-multierror, Serf, or
transfer an earlier exception.

# Measurements At Start

The Errwrap decision was recorded after guard-only revalidation from clean
decision HEAD `f77cb8a75c6813b28656c85b654d17d7df464731`, parent
`a2f1bc98dc6b978d802d42f7c36bae96627bda33`, tree
`9dd5c04d0cc3cdebf482b831faef8809ea25643e`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`. No retained group since
Google UUID has a dependency implementation or metadata commit.

Guard-only revalidation under exact Go 1.26.7 preserved all ten guarded
versions and recorded incoming edges. All ten `go mod why -m` results remain
negative, repository Go source contains zero guarded-path occurrences, and
production and complete-test loads contain zero guarded packages. The project
remains 234 selected modules, 3,599 graph edges, 429 complete-test entries,
197 module-backed entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line unapplied tidy projection. Its `go.mod`/`go.sum` SHA-256
values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. It has no exact Errwrap or other new
guarded-target record; exact Errwrap v1.0.0/v1.1.0 OSV responses remain empty.
Gorilla retains only GO-2020-0019 and unwithdrawn GO-2026-6278 at record
SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.

A bounded queue survey identifies selected Hashicorp go-cleanhttp v0.5.2 in
the current build list. That selection is not proof of repository identity,
release qualification, ancestry, floor, package loading, behavior,
vulnerability state, or suitability. Resolve those facts independently and do
not add a direct edge merely to alter MVS.

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
none is Hashicorp go-cleanhttp evidence.

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
source-time resolution from the project's selected MVS graph and from every
already-final guarded dependency decision.

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize deterministic behavior, errors and identity,
nil/panic behavior, mutation and aliasing, allocation, concurrency, global
state, resource cleanup, and malformed inputs. Distinguish tools, examples,
optional packages, and test-only helpers from behavior actually loaded by this
project.

Add independent fixtures where useful for deterministic behavior, error
identity, nil/panic boundaries, malformed inputs, mutation, aliasing, supported
concurrency, resource cleanup, and selected-project consumers. Run source
verification, package listing, native complete tests, two independent repeats,
race, vet, and meaningful cross-builds under both SDKs. Classify every failure
precisely.

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
`./codex-dev-start.sh --check`. Read this archive, the answered Errwrap
decision and evaluation, the Consul SDK and earlier guarded-decision archives,
rolling handover, roadmap, `go.mod`, `go.sum`, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose
complete minimal source/test closure preserves Go 1.18. Do not promote an
unqualified or floor-ineligible identity. If no candidate satisfies the
existing contracts, preserve the evidence and stop for a fresh bounded product
decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, update the roadmap and rolling handover with exact evidence, outcome,
commit identity, limitations, and next boundary. Answer this archive and
prepare one reciprocal NEXT mission only after the bounded outcome is coherent
and committed. Do not execute the successor.

# Automatic Handoff

After a completed coherent result, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Retain qualified exact selected `github.com/hashicorp/go-cleanhttp v0.5.2`
without changing `go.mod` or `go.sum`. It is the highest exact-path stable
release, its complete minimal source/test closure preserves Go 1.18, and its
source, API, behavior, tests, vulnerability evidence, and exact project effects
pass every applicable contract. No dependency implementation commit and no
product-risk exception were needed.

The evaluation began from clean feature-branch handoff HEAD
`0d1d66480604da5b646c75593570b01ebec72dda`, parent
`f77cb8a75c6813b28656c85b654d17d7df464731`, and tree
`46e8716418edbafc7095bb22b194ed79f67ffb83`. The latest dependency
implementation remains Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`. Startup branch, clean ordinary
and ignored status, archive ancestry, launcher check, and all ten earlier
decision guards were valid. All guarded versions and incoming edges remain
exact, all ten why results remain negative, source occurrences and loaded
packages remain zero, and no new guarded advisory or independent defect
appeared.

Fresh proxy, checksum-database, `go-import`, Git, and GitHub evidence resolves
the exact path to public active unarchived non-fork
`https://github.com/hashicorp/go-cleanhttp.git` under MPL-2.0. The proxy lists
exactly v0.5.0, v0.5.1, and v0.5.2; all are stable. V0.5.2 is proxy `@latest`
and the sole GitHub Release, published 2021-02-03T18:51:13Z. The repository
has one protected `master` branch and no repository security advisory. There
is no prerelease, `/v2` line, alternate module path, redirect, retraction,
deprecation, fork, or later stable identity. Current unreleased master
`2901fbf3e0ecb2512cd7d278977a6b4ae0342ac0` declares Go 1.24 and was not
promoted.

All three tags are lightweight unsigned commit objects and form linear
ancestry v0.5.0 -> v0.5.1 -> v0.5.2 -> master. V0.5.0 is commit
`e8ab9daed8d1ddd2d3c4efba338fe2eeae2e4f18`, parents
`d5fe4b57a186c716b0e00b8c301cbd9b4182694d` and
`003691524b78cc05492d4d4853bc6281ccd10d7c`, tree
`4e73aae6cb561e640cd9619ec84785d956638d1b`, dated
2018-08-30T03:37:06Z. V0.5.1 is commit
`eda1e5db218aad1db63ca4642c8906b26bcf2744`, parent v0.5.0, tree
`d7c6740e1c415649cbf335bed109358e84b7e552`, dated
2019-03-20T05:41:31Z. Selected v0.5.2 is commit
`6d9e2ac5d828e5f8594b97f88c4bde14a67bb6d2`, parent
`d3fcbee8e1810ecee4bdbf415f42f84cfd0e3361`, tree
`c1160f09cedce00dc3ef7b06169ac45cad1c12e8`, dated
2021-02-03T18:51:13Z.

The v0.5.0/v0.5.1/v0.5.2 proxy zip SHA-256 values are respectively
`f5ca0da7f2432a3e961064529f0af8de1113b854bf239a880d16f89424538b44`,
`e3cc9964b0bc80c6156d6fb064abcb62ff8c00df8be8009b6f6d3aefc2776a23`,
and `e9f3dcfcb33172ba499b4f8e888169252d7f1e072082182124a6e2053523f7df`.
Proxy and Git regular-file bytes agree for every release. Their sumdb source
sums are `h1:wvCrVc9TjDls6+YGAF2hAifE1E5U1+b4tH6KdvN3Gig=`,
`h1:dH3aiDG9Jvb5r5+bYHsikaOUIpcM0xvgMXVoDkXMzJM=`, and
`h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=`; v0.5.2's mod sum is
`h1:kO/YDlP8L1346E6Sodw+PrpBSV4/SoxCXGY6BqNFT48=`.

Each release contains exactly LICENSE, README, `cleanhttp.go`, `doc.go`,
`go.mod`, `handlers.go`, and `handlers_test.go`: one package, three production
Go files, and one test file. There are no commands, examples, benchmarks, fuzz
targets, testdata, generated files, build tags, platform branches, cgo, embeds,
go:generate directives, symlinks, or non-standard-library dependencies.
V0.5.0/v0.5.1 have no Go directive; v0.5.2 declares Go 1.13. The complete
minimal imported closure therefore preserves Go 1.18 without an indirect
floor risk. Exact Go 1.26.7 resolves 184 production and 209 complete-test
entries; contained Go 1.18.10 resolves 123 and 147. The target is the sole
external selected module in each isolated closure.

Selected v0.5.2 passes module verification, build, native count-one, two
independent count-ten repeats, race, vet, and production/test cross-builds for
Darwin amd64/arm64, Linux amd64/arm64, Windows amd64, and js/wasm under both
SDKs. V0.5.0 and v0.5.1 pass build, vet, and all cross-builds but reproducibly
fail count-one, both repeats, and race under both SDKs because their release
tests place raw LF, CR, and NUL bytes in request URLs; current Go rejects those
URLs as invalid. V0.5.2 percent-encodes those test inputs, preserves v0.5.1's
nil-request/nil-next fix, and adds `ForceAttemptHTTP2`. Thus selected v0.5.2 is
the only and highest release satisfying the native-test contract.

Pinned apidiff finds identical public declarations in both directions across
all releases: `DefaultTransport`, `DefaultPooledTransport`, `DefaultClient`,
`DefaultPooledClient`, `PrintablePathCheckHandler`, and `HandlerInput` with
exported `ErrStatus`. All three export snapshots have SHA-256
`e96e3d2d100e788eac2ba58bf3e1b07237bedfd3d2ab5969fca30da56b946c64`.

Independent fixtures verify exact transport fields, fresh transport/client
identity without caller aliasing, no package-owned global pointer mutation,
deterministic printable and non-printable routing, Unicode and decoded-control
handling, transient HTTP/1 connection replacement, pooled reuse and
`CloseIdleConnections`, HTTP/2 negotiation, supported immutable handler
concurrency, and allocation counts under both SDKs. A zero caller
`HandlerInput.ErrStatus` is mutated to 400 at construction; the returned
handler retains that pointer, so later mutation changes behavior and concurrent
mutation is unsupported. Nil request and nil next are safe no-ops. A non-nil
request with nil URL and an invalid nonzero response status panic through
standard-library preconditions. Invalid UTF-8 path bytes become printable
replacement runes. Pooled clients own idle resources and must be reused or
closed. Allocation counts are three for a transport, four for a pooled client,
and two for a handler under both SDKs. The fixture SHA-256 is
`1b3aa948358fdbfd136041fc6c930ba904ff8bc6260bd915ea9af50f40dd73a2`.
The open performance discussion about the configurable `GOMAXPROCS+1` idle
limit is not a correctness disqualifier.

Selected v0.5.2 exists through three selected-version graph edges: Viper
v1.15.0, historical Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`.
Historical Consul API/SDK vertices request v0.5.1 and go-retryablehttp v0.5.3
requests v0.5.0. The shortest selected path is main -> Viper v1.15.0 ->
go-cleanhttp v0.5.2. `go mod why -m` is negative, repository Go imports are
zero, and production and complete-test target loads are zero, so the target is
runtime-unreachable.

A disposable exact v0.5.2 get retains 234 selected modules, 429 complete-test
entries, every unrelated version, and zero target load. It manufactures only a
direct indirect root edge and one source checksum line, producing 3,600 graph
edges and 1,068 sum lines. Its `go.mod`/`go.sum` SHA-256 pair is
`55698b6242da86f1e456b0dcf9cadc52f72f2b419500a6196c04b0117281b063` /
`e1b800f07e83f6d7b9e566da23d0de7841f65ccb9866fdd547b90214f0d88cd5`.
Tidy returns exactly to the base projection hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No direct edge was added.

Exact v0.5.0 and v0.5.1 gets force broad unrelated downgrades, remove the
required mvn-pom-mutator root, and make complete project loading fail. They
leave only 182/180 total selected modules and 2,311/2,953 graph edges. Their
tidy copies restore the missing parent but converge to a different historical
Viper projection with hashes
`62e718ee391bcd7dead1a2e1a20bf86b6a12c7af417858d4de7125dac6b7715c` /
`21efc2ced251ca99083e3f0ca82e8c982523a9fd35cee3908bf5bc557ff29d9b`.
Those prohibited unrelated changes were not applied.

Fresh primary vulnerability data contains 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z, with no go-cleanhttp record. Exact OSV
responses for all three releases are empty. Exact Go 1.26.7 isolated module,
package, symbol, and test-symbol scans have zero findings. Go 1.18.10 reports
only that old SDK's standard-library population: 89 unique IDs and
90/154/186/335 module/package/symbol/test-symbol finding paths. Go-cleanhttp
frames occur in seven production IDs/32 paths and 45 test IDs/178 paths that
call the vulnerable old standard library; no advisory is assigned to the
module. Base and disposable v0.5.2 project scans have identical normalized
populations at 30 IDs and 30/52/74 Darwin plus 75 Windows finding paths, with
zero target SBOM package, symbol, test-symbol, or reachable-trace occurrence.

The project remains byte-for-byte unchanged at 234 modules, 3,599 graph edges,
429 complete-test entries, 197 module-backed entries across 41 loaded modules,
1,067 sum lines, and the recorded 432-line tidy projection. `go.mod` and
`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 module verification, build, count-one, two count-ten repeats,
race, vet, pinned golangci-lint 2.12.2, API/CLI compatibility, host and snapshot
acceptance, empty-HOME count-two, four production cross-builds, all 17 script
meta-tests, and all 15 quality-audit controls pass. Contained Go 1.18.10 loads
366 complete-test entries; 26 unaffected packages pass count-one, both
count-ten repeats, race, and vet, and all four production cross-builds pass.
The full Go 1.18 suite retains only the two accepted Darwin `pkg/shell` closed-
file wording failures. Full preflight repeatedly reached passing substantive
stages before reproducing the known outer or nested launcher signal-log timing
race; an independent run passed all 62 launcher controls. A scratch-only
`mktemp` wrapper was required for one unchanged legacy meta-test because the
managed sandbox denied its default macOS temp root. These are target-
independent harness boundaries. The final post-handoff launcher run passed all
62 controls on its second timing attempt. Because no source or dependency
metadata changed, no changed-selection scorecard applies; accepted quality
remains 27/27 Q0-Q2 PASS at L2.

Exact Go 1.26.7 archive/binary SHA-256 remains
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
Go 1.18.10 remains
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Portable golangci-lint and GoReleaser receipts match their required hashes. The
412-entry disposable source/evaluation evidence manifest SHA-256 is
`40b0a6fc11f4378aed62586bfc8bdd94bef33b0a0ba88051dff51912807137c1`.

This qualified result remains final while exact v0.5.2 and its three selected-
version incoming edges remain unchanged, zero target packages load, it remains
runtime-unreachable, and no new advisory or independently disqualifying
behavior appears. Direct import/loading, runtime reachability, a target version
or incoming-edge change, or a new advisory/defect requires a fresh bounded
go-cleanhttp decision. This guard is target-specific and transfers no earlier
exception.
