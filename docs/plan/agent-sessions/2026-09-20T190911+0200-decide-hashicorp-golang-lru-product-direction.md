# Agent Session: Decide Hashicorp Golang LRU Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T190911+0200-decide-hashicorp-golang-lru-product-direction`
Created: `2026-09-20T19:09:11+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `bebc1caad27fdf195999470a5d045ffa6fa92be3f04cbbe10c6f6da300a8a290`
Previous: [2026-09-20T180252+0200-evaluate-hashicorp-golang-lru-dependency.md](2026-09-20T180252+0200-evaluate-hashicorp-golang-lru-dependency.md)
Next: [2026-09-20T192912+0200-evaluate-hashicorp-mdns-dependency.md](2026-09-20T192912+0200-evaluate-hashicorp-mdns-dependency.md)
Outcome: Authorized option 1 for exact inherited, unloaded golang-lru v0.5.4 after every guard passed; product source and dependency metadata stayed unchanged, and one bounded mdns v1.0.4 evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one fresh bounded product decision for selected
exact-path `github.com/hashicorp/golang-lru v0.5.4`. Its independent evaluation
is complete and found no qualifying exact-path stable release. Choose exactly
one option below, record the authorized direction and precise bounds, and
prepare one coherent follow-up without executing it. Do not repeat the audit,
silently accept a defect or API regression, implement before authorization,
combine another dependency group, change the Go floor, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go.net v0.0.1. Every earlier outcome and lifecycle ancestor is final.
No earlier exception transfers to golang-lru. Preserve every guarded selection
and stop for its owning decision if any guard differs. P8 remains queued.

Selected `github.com/hashicorp/golang-lru v0.5.4` is inherited through exact
v0.5.4 requests from Viper v1.15.0, historical Viper v1.10.1, and
`sagikazarmark/crypt v0.4.0`. Its graph also contains the completed fifteen
lower requests: v0.5.0 from go-immutable-radix v1.0.0/v1.3.1 and OpenCensus
v0.21.0, plus v0.5.1 from OpenCensus v0.22.0 and Google API v0.7.0-v0.9.0,
v0.13.0-v0.15.0, v0.17.0-v0.20.0, and v0.22.0. It has no direct main-module
root or repository Go import; `go mod why -m` is negative; production and
complete-test loads contain zero target packages; and it is runtime-
unreachable in this project. These facts bound current exposure but do not
qualify the release or authorize an exception.

# Measurements At Start

The evaluation handoff began from clean ordinary and ignored state at HEAD
`f867a05fa1334c7e7fc0b774edc060e78a9be88d`, parent
`fa9af5e9b168520f024da6f2049065ad186735b8`, tree
`cf389273ef952c0107541facf8671eee1c9b56ef`. Its exact five-file delta,
reciprocal archive chain, latest Google UUID implementation ancestry, and
launcher check passed. Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The project starts this decision unchanged at 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 1,067 sum lines, and the
recorded 432-line tidy projection. Verify the new handoff rather than assuming
these facts.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
public active unarchived non-fork MPL-2.0 repository. Retrievable exact-path
stable releases are v0.5.0-v0.5.4, v0.6.0, v1.0.1, and v1.0.2. A stale
v1.0.0 proxy-list entry has no retrievable proxy endpoint or Git tag and is not
a candidate. V1.0.2 is proxy latest and the final v0/v1 release; current main
declares the distinct Go-1.19 `github.com/hashicorp/golang-lru/v2` path, with
ARC also split to `github.com/hashicorp/golang-lru/arc/v2`. Neither alternate
major path, an unreleased branch, a fork, or a non-versioning tag is an exact-
path candidate. There is no retraction or module deprecation.

Selected v0.5.4 is a lightweight unsigned tag at commit
`14eae340515388ca95aa8e7b86f0de668e981f54`, parent v0.5.3, tree
`635407bf580448ac91cc1265a72cac484c273f21`, time
2020-01-16T13:29:26-05:00. It is an ancestor of v0.6.0 and v1.0.2. Latest
v1.0.2 is commit `a032ef5a154020ffc7a74ba73702f9c5d3ea11f2`, parent v0.6.0,
tree `ddabbfd3c054c45dccdc7c3a04378ea85eff781d`, time
2023-08-08T11:23:39Z, and is GitHub-verified although its lightweight tag is
unsigned. Every retrievable proxy archive matches its Git tag bytes and
sumdb verifies it.

V0.5.4 has two packages, six production Go files, four tests, 27 test
functions, and six benchmarks. It has no command, example, fuzz target,
testdata, generated file, build tag, platform branch, cgo, embed, environment,
network, file, process, or goroutine boundary. V0.6.0/v1.0.2 add one ordinary
production file using only standard-library testing/random helpers. Every
sourceful release has a standard-library-only minimal closure and passes
verification, listing, build, count-one, two independent count-ten repeats,
race, and vet under exact Go 1.26.7 and contained Go 1.18.10. Selected,
v0.6.0, and v1.0.2 also pass production and test cross-compilation for Darwin,
Linux amd64/arm64, Windows, FreeBSD, and `js/wasm` under both SDKs. V1.0.1 is
a package-empty release with no Go API or tests and is not a usable candidate.

The root package exports the synchronized `Cache`, `TwoQueueCache`, and
`ARCCache` constructors and cache operations; `simplelru` exports its
non-thread-safe `LRU`, interface, eviction callback, and matching LRU
operations. Pinned API comparison proves v0.5.0-v0.5.3 remove selected
methods or signatures. V0.6.0/v1.0.2 add only the compatible
`DefaultEvictedBufferSize` symbol, but changing `Cache.lru` from an interface
to a pointer and adding slices makes exported `Cache` non-comparable. Pinned
apidiff reports that incompatibility, and an independent `map[lru.Cache]bool`
fixture builds with v0.5.4 but fails with v0.6.0/v1.0.2 under both SDKs.

Independent behavior fixture SHA-256
`891d2144f6d6adb7acba5788859eb992b7c0238bbff2806f0e9444dfdd93f73e`
characterizes capacity, recency, resize, purge, LRU/2Q/ARC scan resistance,
callbacks, aliasing, nil/panic inputs, allocation, and concurrency. Selected
v0.5.4 invokes an eviction callback while holding `Cache.lock`; a callback
that calls `Len` deadlocks, and a callback panic escapes before unlock and
permanently poisons the cache mutex under both SDKs. V0.6.0/v1.0.2 buffer
evictions and invoke callbacks after unlock, passing reentrancy and panic lock-
safety fixtures, but fail the comparability contract. Earlier releases either
retain the callback defect or already fail selected API compatibility.

Common recorded boundaries are exact: keys are ordered oldest to newest;
Get and update refresh recency while Peek and Contains do not; keys and values
retain caller identity; nil comparable keys and nil values work; non-comparable
keys and nil/zero receivers panic; constructors reject nonpositive capacity;
default `New2Q(1)` and an accepted zero ghost ratio fail through a zero-sized
internal LRU; and `simplelru.Resize(-1)` evicts every item, returns an
over-capacity count, retains negative capacity, and immediately evicts later
adds. Cache/2Q/ARC are synchronized, simplelru is explicitly not, and Get/Peek
measure zero allocations. There is no global mutable state or external
resource ownership.

MVS selects v0.5.4 through the eighteen recorded incoming edges above. A
disposable direct v0.5.4 root keeps the selected version and 234 modules,
changes graph edges 3,599 -> 3,600 and sums 1,067 -> 1,068. Direct v0.6.0,
v1.0.1, and v1.0.2 each change only the target selection, keep 234 modules and
3,600 edges, and produce 1,069 sum lines. All preserve 355 production and 429
complete-test entries with zero target loads and no unrelated version move.
Tidy removes every manufactured root and returns byte-for-byte to the common
base projection. No projection was applied.

All four direct-root project projections pass exact-Go module verification,
build, count-one, two count-ten repeats, race, vet, API/CLI compatibility,
host acceptance, and supported Darwin/Linux/Windows/FreeBSD cross-builds.
The project's existing `js/wasm` failure remains in `chzyer/readline`, not
golang-lru. The canonical unchanged base also passes pinned lint, empty-HOME
count-two, full preflight, every script/meta-test, all eight mutation meta-
stages with 80/80 kills, host/snapshot/Docker acceptance, and all quality-
audit meta-controls. No changed-selection scorecard applies.

Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and has no target entry. Exact OSV results for every retrievable stable version
and both target GitHub advisory feeds are empty. Exact-Go isolated scans are
empty; Go-1.18 findings belong only to that standard library. Base and all four
project projections are byte-identical at 30 module, 22 vulnerable-package,
and 20 called-symbol/test-symbol IDs with no target assignment or trace. All
21 earlier guarded selections, requests, negative why results, zero imports,
zero loads, runtime-unreachability conditions, and advisory states remain
exact.

# Role And Boundaries

Choose exactly one:

1. Retain exact selected, inherited, unloaded golang-lru v0.5.4 without
   changing product source, `go.mod`, or `go.sum`, under a new target-specific,
   non-transferable exception. This is the recommended bounded option for the
   current zero-load graph because it avoids manufacturing a direct root or
   accepting an exported API regression. It must explicitly accept only the
   completed callback-under-lock reentrancy deadlock, callback-panic mutex
   poisoning, constructor/2Q errors, negative-resize behavior, nil/non-
   comparable-key panics,
   caller-identity aliasing, simplelru non-concurrent-use, and related recorded
   API, behavior, MVS, vulnerability, allocation, and resource findings. It
   must expire if v0.5.4 or any of the eighteen incoming requests changes, a
   direct root/import/load or runtime reachability appears, an earlier guard
   changes, or a new advisory or independent defect appears.
2. Authorize a separate Go-1.18-compatible exact-path patch/fork/replacement
   design. That follow-up must establish ownership and provenance, invoke
   callbacks outside locks while preserving exported `Cache` comparability and
   all required API/consumer semantics, repair or explicitly decide every
   recorded boundary, qualify the full source/test/platform/vulnerability
   closure, and measure every MVS/project effect before implementation may be
   retained. It may not silently select v0.6.0, v1.0.2, a v2 path, main, a
   fork, or an alternate module path.
3. Authorize a separate parent/graph-removal decision for the exact Viper,
   crypt, immutable-radix, OpenCensus, and historical Google API request
   population. That follow-up must independently qualify each required parent
   upgrade, replacement, or removal and every transitive selection/API/
   behavior effect. It may not change a guarded parent or remove golang-lru in
   this decision session.

Do not infer option 1 merely from physical selection or zero reachability. Do
not combine options, manufacture a direct root, accept v0.6.0/v1.0.2's API
regression, select a package-empty or alternate-major release, patch source,
change a parent, raise the Go floor, move an unrelated selection, transfer
another exception, or implement anything before the choice is explicitly
authorized. If none is acceptable, stop P7 explicitly without a dependency or
source change.

# Required Reading

This is a product-decision session, not a renewed audit. At start read this
archive, its answered golang-lru evaluation, the answered go.net decision and
evaluation, the rolling handover, roadmap, `go.mod`, and `go.sum`. Verify
branch, clean ordinary and ignored status, handoff ancestry and exact changed-
file set, reciprocal archive chain, exact Go 1.26.7 identity, all eighteen
target requests and zero-load state, every earlier guarded selection/request/
why/import/load condition, module hashes, vulnerability guards, and
`./codex-dev-start.sh --check`. Stop for the owning decision if any guard
differs.

# Three Moves

First, revalidate only the narrow decision guards. Second, obtain and record
one explicit product choice and its exact bounds; do not implement an option
unless the authorization explicitly includes implementation in this session.
Third, update the roadmap and rolling handover, answer this archive, prepare
exactly one reciprocal NEXT mission appropriate to the authorized result, and
commit the handoff. Do not execute the successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit unless
the explicit choice separately authorizes a later implementation commit. Do
not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is explicitly authorized. Retain exact selected
`github.com/hashicorp/golang-lru v0.5.4` as an inherited, unloaded module
without changing product source, `go.mod`, or `go.sum`. This bounded direction
uses the verified zero-load exposure while avoiding a manufactured direct root
and the exported `Cache` comparability regression in v0.6.0/v1.0.2; it is an
explicit product choice, not an inference from physical selection. The new
golang-lru-specific, non-transferable exception accepts only the completed callback-
under-lock reentrancy deadlock; callback-panic mutex poisoning; nonpositive-
capacity constructor errors; the default `New2Q(1)` and zero-ghost-ratio
failures; negative `simplelru.Resize` over-capacity counting, retained negative
capacity, and immediate later eviction; nil and zero receiver and non-
comparable-key panics; caller key/value identity and aliasing; explicitly non-
concurrent simplelru use; and the recorded API, behavior, MVS, vulnerability,
allocation, resource, and related completed findings. It accepts no new or
independently discovered defect.

The exception remains valid only while exact v0.5.4 and all eighteen recorded
incoming requests remain unchanged: v0.5.4 from Viper v1.15.0, historical
Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`; v0.5.0 from
go-immutable-radix v1.0.0/v1.3.1 and OpenCensus v0.21.0; and v0.5.1 from
OpenCensus v0.22.0 and Google API v0.7.0-v0.9.0, v0.13.0-v0.15.0,
v0.17.0-v0.20.0, and v0.22.0. It also requires no direct main-module root or
repository Go import, zero production and complete-test target loads, runtime
unreachability, every earlier guard remaining intact, and no new advisory or
independent defect. A target version or incoming-request change, direct root/
import/load, runtime reachability, earlier owning-guard change, or new advisory
or independent defect expires the exception and requires the owning fresh
dependency and product decision before merge. The exact-path remediation and
parent/graph-removal alternatives are not authorized. V0.6.0, v1.0.2, the
package-empty v1.0.1 release, v2 paths, main, forks, and alternate paths remain
unauthorized, and no exception transfers to another target.

Guard-only revalidation began from clean ordinary and ignored worktree state
on branch `codex/upgrade-quality` at HEAD
`961c48582ca568f0867d91e7781c6fa0b970dcea`, parent
`f867a05fa1334c7e7fc0b774edc060e78a9be88d`, tree
`716bf2f170c296836aa1903414c55ab4f5450109`. That handoff changed exactly the
launcher, answered golang-lru evaluation archive, this then-NEXT decision
archive, rolling handover, and roadmap. The evaluation handoff `f867a05`,
parent `fa9af5e`, tree `cf389273`, retains its exact five-file delta. The
reciprocal 222-archive chain, latest Google UUID v1.4.0 implementation
ancestry, and launcher check pass.

Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`. Golang-lru remains
exact v0.5.4 through all eighteen recorded requests. There is no direct root,
its `go mod why -m` result remains negative, repository Go imports remain
zero, and the 355-entry production and 429-entry complete-test loads contain
zero target packages.

All 21 earlier guarded modules retain their exact selections and recorded
requests. All 22 target-plus-earlier why results remain negative, repository
Go imports remain zero, and production and complete-test loads contain zero
guarded packages. The complete-test load retains 197 module-backed entries
across 41 loaded modules, so golang-lru and every earlier guarded target remain
runtime-unreachable. The project remains 234 modules, 3,599 graph edges, 1,067
sum lines, and the recorded 432-line tidy projection. `go.mod` and `go.sum`
retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact module verification passes.

Fresh primary vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no golang-lru record. Exact-version
OSV results remain empty for golang-lru and every guarded target except the
recorded Gorilla WebSocket GO-2026-6278/GHSA-w67g-5rqw-f597 and
go-retryablehttp GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Both exact target
GitHub advisory feeds remain empty. No new advisory or independently observed
defect appeared.

No dependency implementation or metadata commit was created. No changed-
selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS at
L2. The sole reciprocal successor is the bounded P7 evaluation of selected
exact-path `github.com/hashicorp/mdns v1.0.4`; it was prepared but not
executed.
