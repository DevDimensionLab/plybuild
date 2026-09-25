# Agent Session: Evaluate Jtolds GLS Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-21T162538+0200-evaluate-jtolds-gls-dependency`
Created: `2026-09-21T16:25:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ffcba3dbc87f662197403e5e4f0990e7a2f621f9e38537da615e2504e82a9570`
Previous: [2026-09-21T160501+0200-decide-jstemmer-go-junit-report-product-direction.md](2026-09-21T160501+0200-decide-jstemmer-go-junit-report-product-direction.md)
Next: [2026-09-21T171848+0200-decide-jtolds-gls-product-direction.md](2026-09-21T171848+0200-decide-jtolds-gls-product-direction.md)
Outcome: No exact-path stable release qualifies: v4.2.0 and v4.2.1 fail their upstream behavior suite under both required SDKs, while v4.20.0 deterministically violates documented point-in-time `Go` propagation and races on the aliased values map; product and dependency metadata remain unchanged and P7 stops for one bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/jtolds/gls v4.20.0+incompatible` as one bounded dependency group.
Resolve its complete repository and stable-release identity, legacy
`+incompatible` module semantics, Go-floor closure, exported API and ordinary
behavior, actual project loading, exact MVS effects, vulnerability evidence,
and every applicable quality contract. Retain or select only a qualified
exact-path stable release whose complete minimal source/test closure preserves
Go 1.18 and whose relevant behavior passes every contract; otherwise leave
metadata unchanged and stop for one fresh bounded product decision. Do not
combine or alter the owning goconvey parent, evaluate another dependency
group, or begin P8.

# Defensive Scope

This is an ordinary dependency-quality review. Use public metadata, static
source/API inspection, admissible upstream tests, bounded deterministic
fixtures with ordinary documented inputs, project graph/build commands, and
public advisory evidence. Do not fuzz, stress, probe resource exhaustion,
generate oversized or deeply nested input, or perform security or
exploitability analysis. Any malformed-input check needed for documented error
behavior must be small, deterministic, and non-adversarial.

Every disposable archive, clone, cache, tool, binary, report, fixture, and
project copy must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write
to `/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2, and
all target-specific retained-module decisions through exact selected,
inherited, unloaded go-junit-report v0.9.1. Every earlier outcome is final
under its own guards. P8 remains queued.

The go-junit-report option-1 decision explicitly retains exact v0.9.1 without
product-source or dependency-metadata changes. It is not qualified. Its target-
specific, non-transferable exception accepts only the completed ordinary
writer-error, repository/archive, API/command, behavior, ownership/lifecycle,
Go-floor, MVS, loading, repeatability, vulnerability, and related findings. It
is bounded by exact v0.9.1, all 25 historical requests and genuine parent
identities, negative target why, no direct root/import/load/runtime
reachability, exact graph/module/tidy and earlier guards, and no new advisory,
independent defect, qualified stable release, genuine supported tidy-stable
owner, or compatible qualified route. Any change expires that decision and
requires its fresh owning decision. No go-junit-report, json-iterator,
clockwork, demangle, strcase, memberlist, or other exception transfers.

Selected jtolds/gls v4.20.0+incompatible is only a queue identity. Current
graph evidence shows one historical
`github.com/smartystreets/goconvey@v1.6.4 -> github.com/jtolds/gls@v4.20.0+incompatible`
request, negative target why, zero repository imports, zero production or
complete-test loads, and no runtime reachability. Independently verify each
fact. Physical MVS selection and zero loading are not qualification or
authorization to retain it. Do not add a direct edge merely to alter MVS.

# Measurements At Start

The go-junit-report decision began from clean ordinary and ignored state on
branch `codex/upgrade-quality` at handoff HEAD
`efbe47116c0c56a5a0f7bbf67133d49312286846`, parent
`5946d0958a1538bdea281bb66c59647d0e1a2824`, tree
`ed26c4945de0f1be4ba20d46fbf94bb538a2a985`. That handoff changed exactly the
launcher, answered go-junit-report evaluation archive, then-NEXT go-junit-
report decision archive, rolling handover, and roadmap. Verify the new
decision handoff, reciprocal archive chain, latest Google UUID implementation
ancestry, exact Go identity, and launcher check rather than assuming these
facts.

The unchanged project has 234 selected modules, 3,599 graph edges, 355
production and 429 complete-test entries, 197 module-backed complete-test
entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
tidy projection. `go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Including go-junit-report, all 29 guarded selections and their 213 incoming
graph edges retain sorted snapshot SHA-256
`5d2a35a2961c04eb07dd927afc072e77f555d87c7317579ddda486bf128d2c93`;
all guarded why results are negative and guarded imports/loads are zero.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh exact go-junit-report stable/v1-branch OSV and GitHub advisory results
are empty. The Go vulnerability index remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
guard OSV retains only the recorded Gorilla/go-retryablehttp pairs, and the
2,807-byte PUBLISHED memberlist CNA response remains exact at
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

# Role And Boundaries

From fresh public proxy, sumdb, `go-import`, Git, forge, release, and advisory
evidence, resolve the exact jtolds/gls module path, every stable and serious
candidate, repository state, tags/releases, commits, ancestry, signatures,
license, retractions, deprecation, redirects, forks, alternate paths, and
major lines. Resolve the meaning and support status of the selected
`+incompatible` version without silently promoting a fork, branch, prerelease,
redirect, alternate path, version-masquerading replacement, or floor-
ineligible candidate.

Prove each serious candidate's complete minimal production and test closure
under exact Go 1.26.7 and contained Go 1.18.10. Inspect every package, exported
API, ordinary goroutine-local storage lifecycle and cleanup behavior, caller
ownership, aliasing, determinism, concurrency and global state, platform/build-
tag branches, examples, benchmarks, testdata, generated files, cgo, and
external boundaries. Use only bounded ordinary fixtures permitted by the
defensive scope and classify upstream-test or environment failures precisely.

Measure exact project module, graph, package, checksum, tidy, compatibility,
acceptance, and vulnerability effects in disposable copies. Explain why
jtolds/gls exists in MVS and whether any package loads. Preserve go-junit-
report, json-iterator, clockwork, mvn-pom-mutator, demangle, pprof, strcase,
memberlist, and every earlier guarded decision. A goconvey-parent, Go-floor,
unrelated-selection, product-source, or non-exact-path change requires its own
fresh bounded decision; do not manufacture a direct dependency owner.

# Required Reading

Read this archive, the answered go-junit-report decision/evaluation, the
answered json-iterator, clockwork, and demangle decisions/evaluations, pprof
and strcase records, the memberlist ownership chain, relevant retained-module
records, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify branch,
clean ordinary and ignored state, handoff HEAD/parent/tree and changed set,
reciprocal archive chain, latest Google UUID implementation ancestry, exact Go
identity, module hashes, all guarded selection/edge/why/import/load/advisory
conditions, scratch containment, and `./codex-dev-start.sh --check`. Earlier
outcomes are final.

# Three Moves

First, revalidate the exact starting guards and independently identify the
highest qualified exact-path stable jtolds/gls release. Second, if and only if
one candidate preserves Go 1.18 and passes every applicable contract,
implement that exact dependency-only selection and run the normal changed-
selection gate; otherwise leave source and metadata unchanged and stop for one
bounded product decision. Third, update the roadmap and rolling handover,
answer this archive, prepare exactly one reciprocal successor matching the
result, verify scratch containment, and commit the handoff without executing
the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, reopen go-junit-report/json-iterator/
clockwork/mvn-pom-mutator/demangle/pprof/strcase/memberlist work, write outside
the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

# Answer

No exact-path stable `github.com/jtolds/gls` release qualifies. Product source,
`go.mod`, and `go.sum` remain unchanged. The next and only NEXT archive is the
reciprocal bounded jtolds/gls product decision linked above; it was prepared
but not executed.

## Continuity and baseline

- Evaluation began on `codex/upgrade-quality` at clean HEAD
  `5dd19af0436917968bba1c25baca8a35a1b2e721`, parent
  `7f0b2c1e3506f3a7f3185e4730ba8140d22eaf99`, tree
  `7e0aaa8e15096f58518863026acad52a053ca4c0`. HEAD changes only the launcher
  mutable-boundary blank line. The preceding go-junit-report decision handoff
  has the required five-file changed set from `efbe471...`; the reciprocal
  archive chain, launcher check, and latest Google UUID implementation
  ancestry at `cf53bc64eeb69471d35c7536d196bf1da15f3973` passed.
- Fresh official darwin/arm64 SDKs identify exact Go 1.26.7 archive/binary
  SHA-256 as
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and contained Go 1.18.10 as
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
- Baseline remains 234 selected modules, 3,599 graph edges, 355 production and
  429 complete-test entries, 197 module-backed complete-test entries across 41
  modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  The unchanged 432-line tidy projection is SHA-256
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`
  and produces the recorded 52/948-line projected files.
- All 29 prior guarded selections and their 213 sorted incoming edges remain
  exact; the edge snapshot independently reproduces SHA-256
  `5d2a35a2961c04eb07dd927afc072e77f555d87c7317579ddda486bf128d2c93`.
  All 30 target-plus-prior `go mod why -m` results are negative and repository
  imports plus production/complete-test loads are zero.

## Repository, releases, and legacy module identity

- The old GitHub URL redirects to the public, enabled, unarchived, non-fork
  MIT repository `github.com/jtolio/gls`; old-path `go-import` metadata also
  names `github.com/jtolio/gls` and that repository. This rename is an
  alternate module/repository identity and was not silently promoted.
- The exact old-path proxy exposes only stable v4.2.0+incompatible,
  v4.2.1+incompatible, and v4.20.0+incompatible. There is no `/v4` module,
  prerelease, GitHub release object, retraction, deprecation, replacement, or
  later master commit. The canonical new unsuffixed path exposes the same
  releases but is not the selected exact path.
- Tags are lightweight and unsigned. V4.2.0 is commit
  `8ddce2a84170772b95dd5d576c48d517b22cac63`, v4.2.1 is
  `77f18212c9c7edc9bd6a33d383a7b545ce62f064`, and both `v4.20` and
  `v4.20.0` name current master commit
  `b4936e06046bbecbb94cae9c18127ebe510a2cb9`, tree
  `21801a3091b3d1c1822c0621039366472633fafd`. Their ancestry is ordered and
  GitHub reports each commit unsigned.
- `+incompatible` means these v4 tags belong to a legacy unsuffixed path with
  no source `go.mod`. The proxy admits them only with that suffix and
  synthesizes the one-line `module github.com/jtolds/gls` file. The selected
  source/synthetic-mod sumdb records are
  `h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=` /
  `h1:QJZ7F/aHp+rZTRtaJ1ow/lLfFfVYBRgL+9YlvaHOwJU=`.
- Proxy archive SHA-256 values are v4.2.0
  `13a0e92d7a7ad03ec0e97bc407dd2f69ae31b042b3936f199055829c559c4f85`,
  v4.2.1
  `f75147a45f26352f2ca66499a698a503d2eb9c52e1d67e6004f82db28a750a20`,
  and v4.20.0
  `2f51f8cb610e846dc4bd9b3c0fbf6bebab24bb06d866db7804e123a61b0bd9ec`;
  their 9/10/10 regular-file manifests are byte-identical to Git exports.
  The common MIT license SHA-256 is
  `c7729e27f504fa6fa6c6e67a3309b89b151010f8a8a0d83206be875f20f7a430`.
  No symlink or submodule is present.

## Closure, API, and ordinary behavior

- There is one library package and no command, cgo, embed, generated marker,
  testdata, filesystem/network/subprocess boundary, or external production
  resource. The selected tree has seven production Go files, one test file,
  one upstream test, two examples, and two benchmarks. Benchmarks were not run
  under the defensive scope. Legacy `js`/`!js` build tags select GopherJS or
  `runtime.Callers` stack-tag implementations.
- Normalized `go doc -all` for v4.2.0 is 62 lines at SHA-256
  `a5c6f0d73ba45825f1a746fa9de227980562e9ad25d59465c729ae04d0043c01`.
  V4.2.1/v4.20.0 expose identical 73-line APIs at
  `7a4e14cb42d693700a6833d4bd26fd94e93a6810aa01b2d4b37101cf32c9cfb5`:
  `Values`, `ContextKey`, `ContextManager`, construction/registration/value
  methods, `GenSym`, `GetGoroutineId`, `EnsureGoroutineId`, and `Go`.
  V4.20.0 differs from v4.2.1 only by adding `//go:noinline` to the stack-tag
  marker functions.
- Under both exact Go 1.26.7 and Go 1.18.10, v4.2.0 and v4.2.1 compile and
  vet but deterministically fail the upstream test/examples because modern
  compilers erase the expected stack tags. V4.20.0 passes upstream count-one,
  count-ten, race count-ten, and vet under both SDKs. All three test binaries
  cross-compile under both SDKs for darwin/amd64, linux/amd64, linux/arm64,
  windows/amd64, and freebsd/amd64.
- Native source/test closure is target plus standard library. A complete
  build-tag tidy today resolves omitted JS import
  `github.com/gopherjs/gopherjs/js` to gopherjs v1.21.0 across an 18-module
  graph. That dependency declares Go 1.21, so the release metadata does not
  preserve the Go 1.18 floor even though both SDKs compile the js/wasm test
  binary. The actual project separately selects the historical GopherJS
  pseudo-version, but no GLS package loads and that inherited graph fact does
  not repair the release's incomplete module metadata.
- `NewContextManager` registers a pointer in package-global state and the
  caller must eventually `Unregister`; registry, ID pool, and symbol counter
  operations use mutexes. `SetValues` makes a shallow map copy and restores it
  after normal or panic return; keys and values retain caller identity. The
  underlying returned/captured values maps are not independently synchronized,
  and goroutine IDs are explicitly callback-scoped and reused.
- A bounded ordinary fixture (SHA-256
  `bd3a0422253d0fa507c6db28ec0d16be3532a53354c34e18a11a3e6925c0b091`)
  passes nesting/restoration, panic cleanup, four ordinary propagated
  goroutines, isolation, and unregister behavior under both SDKs. It exposes
  one deterministic selected-release defect: with one ordinary scheduler
  slot, `gls.Go` called inside nested `SetValues` returns the outer value,
  rather than the documented point-in-time inner value, because `getValues`
  returns the live map and copying is deferred to the child goroutine. The
  same case reports a map read/write race under `-race -count=10` while the
  parent restores the map. Count-one and count-ten reproduce the wrong value
  under both SDKs. This is an ordinary behavior/aliasing/concurrency contract
  failure, not security analysis.

## MVS, project effects, and advisories

- MVS selects v4.20.0+incompatible solely through
  `github.com/smartystreets/goconvey@v1.6.4 -> github.com/jtolds/gls@v4.20.0+incompatible`.
  The genuine route begins at direct `mvn-pom-mutator v0.2.3`; goconvey and
  GLS both have negative why results and load no package. There is no direct
  target root, repository import, production/test load, or runtime reachability.
- A disposable direct selected root retains 234 modules and 355/429 loads,
  changes the graph only from 3,599 to 3,600 edges, and changes sums only from
  1,067 to 1,068 lines. Its hashes are
  `55a1e913f947f74f6b8397a05db2c087d477524def58a4d4187abdd57f7a0cbe` /
  `d5eeb0f3c5efd1f9e6c74eb328802b1c49a2ded6ec08d32d18b54129f30e3a90`.
  It remains unloaded and passes verify, build, count-one, count-ten, race,
  vet, API/CLI/surface compatibility, and all four host acceptance scripts.
  Tidy removes the manufactured root, retains the historical goconvey request,
  and restores the common 234-module/3,557-edge/948-sum projection and exact
  baseline projection hashes. No projection was retained.
- Fresh exact-version and package-wide OSV and GitHub global queries, the
  canonical repository advisory endpoint, and the 1,402-record Go index are
  empty for all three candidates. Pinned govulncheck v1.8.0 module/package/
  symbol/test-symbol scans are empty for each. Base and direct-selected project
  scans are byte-identical at 30/22/20/20 findings with no target trace; their
  normalized hashes are
  `26eeb337b34f2d252dc9afd47e7e0efc08aaa67e4233cf360005d900d4dd0dbf`,
  `ad0897a7284d45475df0176e7bb19dfa7bfc66dc70ed7e7142a183b3e676b2db`,
  and `111bb9668ebcd5f7d5aada04ba5b4637dd25fe90fbc03885673b41dfe3b59264`
  for module, package, and both symbol populations.
- The Go index remains 518,501 bytes/1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
  Guard OSV retains only Gorilla v1.4.2 GO-2026-6278/GHSA-w67g-5rqw-f597 and
  go-retryablehttp v0.5.3 GO-2024-2947/GHSA-v6v8-xj6m-xwqh. The PUBLISHED
  memberlist CNA response remains 2,807 bytes at SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Because the older releases fail their own ordinary suite and selected/latest
v4.20.0 fails documented propagation, race, and complete floor-metadata
contracts, none can be selected as qualified. The unchanged project passes
final exact-Go module verification, build, count-one, race, and vet under
`umask 022`. No changed-selection scorecard applies; accepted quality remains
27/27 Q0-Q2 PASS at L2. Task-owned scratch evidence was confined to
`${CODEX_SESSION_SCRATCH_ROOT:?}` and removed before handoff.
