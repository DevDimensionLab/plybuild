# Agent Session: Evaluate Hashicorp Golang LRU Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T180252+0200-evaluate-hashicorp-golang-lru-dependency`
Created: `2026-09-20T18:02:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f45f5b50dea9a7c8315e5da6d054518310a22d0d98272801d0d07be7393cfff5`
Previous: [2026-09-20T174526+0200-decide-hashicorp-go-net-product-direction.md](2026-09-20T174526+0200-decide-hashicorp-go-net-product-direction.md)
Next: [2026-09-20T190911+0200-decide-hashicorp-golang-lru-product-direction.md](2026-09-20T190911+0200-decide-hashicorp-golang-lru-product-direction.md)
Outcome: No exact-path stable release qualified: v0.5.4 deadlocks on callback reentrancy and poisons its mutex after callback panic, while fixed v0.6.0/v1.0.2 break exported Cache comparability; metadata stayed unchanged and one bounded product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/golang-lru v0.5.4` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go.net v0.0.1. Every earlier outcome and lifecycle ancestor is final.
Evaluate only Hashicorp golang-lru in this session; do not reopen or combine
another dependency group. P8 remains queued.

The authorized 2026-09-20 go.net option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/go.net v0.0.1` without product
source or dependency metadata changes. Its target-specific, non-transferable
exception accepts only the completed unresolved production/test closure,
obsolete imports, Linux-arm64/wasm gaps, panic, blocking, error-identity,
aliasing, global-state, timeout, resource, platform, API, MVS, vulnerability,
and related recorded findings. It remains valid only while exact v0.0.1, its
sole exact request from `github.com/hashicorp/mdns v1.0.0`, no direct root or
repository import, zero target loads, runtime unreachability, every earlier
guard, and the no-new-advisory-or-defect condition remain exact. Any change
requires the owning fresh decision.

The go-uuid v1.0.1, go-syslog v1.0.0, go-sockaddr v1.0.0, go-rootcerts
v1.0.2, go-retryablehttp v0.5.3, go-multierror v1.1.0, go-msgpack v0.5.3,
go-immutable-radix v1.3.1, go-hclog v1.2.0, Errwrap v1.0.0, qualified
go-cleanhttp v0.5.2, and every other recorded exception or qualification
remain separate under their exact selection, incoming-request, zero-load,
runtime-unreachable, and no-new-finding guards. Revalidate those guards and
stop for the owning decision if any expires. No earlier exception transfers
to golang-lru. Do not change a guarded parent, the Go floor, or an unrelated
module.

Selected `github.com/hashicorp/golang-lru v0.5.4` is inherited through exact
v0.5.4 requests from Viper v1.15.0, historical Viper v1.10.1, and
`sagikazarmark/crypt v0.4.0`; lower requests include the recorded v0.5.0 and
v0.5.1 graph edges. The current repository import search and production and
complete-test loads contain zero target packages, its why result is negative,
and there is no direct main-module root. These queue observations and the
physical MVS selection are not proof of repository identity, release
qualification, ancestry, floor, behavior, vulnerability state, or suitability.
Resolve them independently and do not add a direct edge merely to alter MVS.

# Measurements At Start

The go.net decision recording began from clean HEAD
`fa9af5e9b168520f024da6f2049065ad186735b8`, parent
`70867a4a7701f69cbc4551de3d4db53fd7d86fe2`, tree
`53de66fc8ea39ffdfa68611677564af70f487f8e`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive chain, and `./codex-dev-start.sh --check` rather than
assuming them. The latest dependency implementation remains exact Google UUID
v1.4.0 commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 20
earlier guarded selections and recorded requests plus exact go.net v0.0.1.
All 21 why results remain negative, repository imports are zero, and production
and complete-test loads contain zero guarded packages. The project remains 234
modules, 3,599 graph edges, 355 production entries, 429 complete-test entries,
197 module-backed entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line unapplied tidy projection. Base `go.mod` and `go.sum`
SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
Exact-version OSV remains empty for go.net and every guarded module except the
recorded Gorilla GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp
GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Go.net's exact GitHub global and
repository advisory queries remain empty. No new guarded advisory or
independent defect appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, Make-variable inheritance boundary, and nested
launcher signal-retention timing race; none is golang-lru evidence.

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
testdata, generated file, platform/build-tag branch, and applicable API and
runtime boundary. Characterize capacity, eviction and recency ordering,
resize/purge behavior, ARC and 2Q adaptation, callback and reentrancy behavior,
key/value ownership and aliasing, nil/panic and malformed-input behavior,
allocation, concurrency, global state, resources, environment interaction, and
actual project consumers. Add independent fixtures where useful and run source
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

At start read this archive, the answered go.net decision and evaluation, the
answered go-uuid decision and evaluation, rolling handover, roadmap, `go.mod`,
and `go.sum`. Verify the recorded handoff identity, changed-file set, clean
ordinary and ignored status, exact toolchain identity, reciprocal archive
history, every guarded selection/request/load result, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, resolve golang-lru identity, release line, source/test closure, API,
behavior, vulnerability, load, and exact MVS facts without changing the
worktree. Second, if and only if one exact-path stable release preserves Go
1.18 and passes every applicable contract, implement that exact dependency-
only selection and run the complete changed-selection gate; otherwise leave
metadata unchanged and stop for a bounded product decision. Third, update the
roadmap and rolling handover, answer this archive, prepare exactly one
reciprocal NEXT mission for the authorized result, and commit the handoff. Do
not execute the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable release qualifies under the authorized contracts. The
selected v0.5.4 release has a callback-under-lock defect that deadlocks on
reentrancy and permanently poisons the cache mutex when a callback panics.
V0.6.0 and v1.0.2 repair that behavior but make the exported `Cache` type
non-comparable, breaking code that uses a `Cache` value as a map key or in any
other comparable-only context. Earlier releases remove selected methods or
signatures, and v1.0.1 contains no Go package. No dependency implementation or
metadata commit was created. P7 therefore stops on the sole prepared bounded
golang-lru product decision.

Evaluation began from clean ordinary and ignored worktree state on branch
`codex/upgrade-quality` at HEAD
`f867a05fa1334c7e7fc0b774edc060e78a9be88d`, parent
`fa9af5e9b168520f024da6f2049065ad186735b8`, tree
`cf389273ef952c0107541facf8671eee1c9b56ef`. That handoff changed exactly the
launcher, answered go.net decision archive, this then-NEXT golang-lru
evaluation archive, rolling handover, and roadmap. Its reciprocal 221-archive
chain, latest Google UUID v1.4.0 implementation ancestry, launcher check, and
exact five-file delta passed. Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`.

Fresh proxy enumeration returns v0.5.0-v0.5.4, v0.6.0, v1.0.1, and v1.0.2
as retrievable stable exact-path releases. The proxy list also contains stale
v1.0.0, but every proxy endpoint rejects it as unknown revision and Git has no
such tag, so it is not a candidate. Proxy latest is v1.0.2. `go-import`
resolves `github.com/hashicorp/golang-lru` to Hashicorp's public, active,
unarchived, non-fork MPL-2.0 repository. The exact-path v1 branch ends at
v1.0.2. Current main instead declares the distinct Go-1.19 module
`github.com/hashicorp/golang-lru/v2`, with ARC separately declared as
`github.com/hashicorp/golang-lru/arc/v2`; neither alternate path was promoted.
No retraction, module deprecation, redirect, prerelease, or qualifying
unreleased branch exists.

Selected v0.5.4 is a lightweight unsigned tag at commit
`14eae340515388ca95aa8e7b86f0de668e981f54`, parent
`7f827b33c0f158ec5dfbba01bb0b14a4541fd81d` (v0.5.3), tree
`635407bf580448ac91cc1265a72cac484c273f21`, time
2020-01-16T13:29:26-05:00. Proxy ZIP SHA-256 is
`7b2a8b1739c858727fca497a6415323edb801dc97b8aca04f7bac4ab9fb5c66b`;
sumdb verifies module sum
`h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=` and go.mod sum
`h1:iADmTwqILo4mZ8BN3D2Q6+9jd8WM5uGBxy+E8yxSoD4=`. Proxy and Git contents
match byte-for-byte. V0.5.4 is an ancestor of v0.6.0 and v1.0.2. V1.0.2 is
GitHub-verified commit `a032ef5a154020ffc7a74ba73702f9c5d3ea11f2`, parent
v0.6.0, tree `ddabbfd3c054c45dccdc7c3a04378ea85eff781d`, time
2023-08-08T11:23:39Z; its lightweight tag remains unsigned. The GitHub release
calls it the last v0/v1 release. V0.5.4 has no GitHub release object.

Selected v0.5.4 contains 14 regular files, ten Go files, six production Go
files, four tests, two packages, 27 test functions, and six benchmarks. It has
no command, example, fuzz target, testdata, generated file, build constraint,
platform branch, cgo, embed, generate directive, symlink, special file,
resource, environment, network, process, or goroutine interaction. V0.6.0 and
v1.0.2 contain 17 regular files and seven production Go files; their added
ordinary `testing.go` uses only standard-library crypto/rand, math/big, and
testing helpers. V1.0.1 has five regular files and no Go package or test.

Every sourceful release has a standard-library-only minimal production/test
closure. V0.5.2 and later declare Go 1.12; v0.5.0/v0.5.1 have no `go`
directive. Under exact Go 1.26.7 and contained official Go 1.18.10, every
sourceful release passes source verification, package listing, build,
count-one, two independent count-ten repeats, race, and vet. Selected v0.5.4,
v0.6.0, and v1.0.2 also pass production and test cross-compilation for Darwin
amd64, Linux amd64/arm64, Windows amd64, FreeBSD amd64, and `js/wasm` under
both SDKs. V1.0.1 cannot qualify because package listing is empty and vet
correctly reports no packages. The contained Go 1.18.10 archive/binary SHA-256
values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.

The root package exports constructors and synchronized operations for
`Cache`, `TwoQueueCache`, and `ARCCache`; `simplelru` exports `LRUCache`,
`LRU`, `EvictCallback`, its constructor, and matching non-synchronized LRU
operations. Pinned API archives show v0.5.0/v0.5.1 remove Resize, GetOldest,
PeekOrAdd, and selected return signatures; v0.5.2 still removes GetOldest,
PeekOrAdd, and signatures; and v0.5.3 removes PeekOrAdd. V0.6.0/v1.0.2 add
compatible `DefaultEvictedBufferSize`, but pinned apidiff reports the
incompatible change from comparable `Cache` to non-comparable `Cache`.
Independent `map[lru.Cache]bool` compilation passes at v0.5.4 and fails at
v0.6.0/v1.0.2 under both SDKs. V0.5.4 API archive SHA-256 is
`29046184b0c17247cfa4cb5a9af623f80a8ac3cd719978d886a380f128eaf40d`;
v0.6.0/v1.0.2 share
`a97b8fd67934c651a2b188f166417a0c3d13ca3bcd85c6668c15d8f976f7068c`.

Independent behavior fixture SHA-256
`891d2144f6d6adb7acba5788859eb992b7c0238bbff2806f0e9444dfdd93f73e`
covers capacity, recency, eviction ordering, Peek/Get/update, Resize, Purge,
2Q/ARC scan resistance and adaptation, callback reentrancy and panic safety,
key/value identity and aliasing, nil and malformed inputs, concurrent access,
and allocation. V0.5.4 invokes the user eviction callback from simplelru while
holding `Cache.lock`; a callback calling `Len` deadlocks, and a callback panic
propagates before unlock and leaves the mutex permanently locked under both
SDKs. V0.6.0/v1.0.2 buffer evictions and call users after unlocking, passing
those fixtures, repeats, race, and vet under both SDKs, but their API
comparability regression independently disqualifies them.

The characterized common behavior is exact. Keys enumerate oldest to newest;
Get and update refresh recency while Peek and Contains do not. 2Q and ARC keep
frequently reused entries through a scan and respect configured capacity;
Purge clears them. Caches retain caller-owned key/value identity, so referenced
value mutations remain visible. Nil comparable keys and nil values work;
non-comparable keys and nil/zero receivers panic. Constructors reject
nonpositive sizes. Default `New2Q(1)` truncates its ghost partition to zero and
fails; `New2QParams` accepts ghost ratio zero but its internal zero-size LRU
then fails. `simplelru.Resize(-1)` evicts both entries, returns three, retains
negative capacity, and immediately evicts later additions. Cache/2Q/ARC are
thread-safe; simplelru is explicitly not. Get and Peek measure zero
allocations. There is no mutable global state or external resource ownership.

Project MVS selects v0.5.4 through eighteen incoming graph edges: exact v0.5.4
from Viper v1.15.0, Viper v1.10.1, and sagikazarmark/crypt v0.4.0; v0.5.0
from go-immutable-radix v1.0.0/v1.3.1 and OpenCensus v0.21.0; and v0.5.1
from OpenCensus v0.22.0 plus Google API v0.7.0-v0.9.0, v0.13.0-v0.15.0,
v0.17.0-v0.20.0, and v0.22.0. There is no direct root, the why result is
negative, repository Go imports are zero, and production/complete-test loads
contain zero target packages.

A disposable direct v0.5.4 root retains all 234 selected modules, adds only
one main graph edge and the content checksum, yielding 3,600 edges and 1,068
sum lines. Direct v0.6.0/v1.0.1/v1.0.2 roots move only the target, retain 234
modules and 3,600 edges, and yield 1,069 sum lines. Every projection preserves
355 production entries and 429 complete-test entries with zero target loads.
Tidy removes each manufactured root and returns exactly to common projection
hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection was applied.

All four disposable project projections pass exact-Go module verification,
build, count-one, two count-ten repeats, race, vet, API/CLI compatibility,
host acceptance, and Darwin amd64, Linux amd64/arm64, Windows amd64, and
FreeBSD amd64 cross-builds. The project's pre-existing `js/wasm` failure is in
`github.com/chzyer/readline`, not the unloaded target. The unchanged base also
passes pinned golangci-lint 2.12.2, empty-HOME count-two, canonical full
preflight, every production-script/meta-test pair, all eight mutation meta-
stages with 80/80 kills, host acceptance, pinned GoReleaser snapshot
acceptance, daemon-backed Docker acceptance under Python 3.14.6, and every
quality-audit meta-control. Initial default-cache, cold offline API-base,
bare-`mktemp`, Make-variable inheritance, and unsupported project-wasm attempts
were superseded by the contained canonical runs and are not target evidence.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no target entry. Exact OSV queries
for every retrievable stable version and both exact GitHub advisory feeds are
empty. Exact-Go isolated module/package/symbol/test-symbol scans are empty;
Go-1.18.10 findings belong only to that standard library. Base and all four
direct-root project scans are byte-identical at 30 module, 22 vulnerable-
package, and 20 called-symbol/test-symbol IDs, with zero golang-lru assignment
or trace.

All 21 earlier guarded selections and recorded requests remain exact; all 21
why results remain negative, repository imports are zero, and production and
complete-test loads contain zero guarded packages. The base remains 234
modules, 3,599 edges, 355 production entries, 429 complete-test entries, 197
module-backed entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line tidy projection. `go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Earlier advisory states remain exact: only the recorded Gorilla and
go-retryablehttp pairs are non-empty.

Because no source or dependency metadata changed, no changed-selection
scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS at L2. The
sole reciprocal successor is a decision-only golang-lru session offering
guarded v0.5.4 retention, a separately scoped Go-1.18-compatible and
API-preserving remediation design, or a separately scoped parent/graph-removal
decision. It was prepared but not executed.
