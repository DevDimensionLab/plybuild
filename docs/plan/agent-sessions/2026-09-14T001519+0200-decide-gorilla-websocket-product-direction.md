# Agent Session: Decide Gorilla WebSocket Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T001519+0200-decide-gorilla-websocket-product-direction`
Created: `2026-09-14T00:15:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9477119a4153363ab50f93e8cbf8f15cda0e05754b5933d60e142cee85f7f4a5`
Previous: [2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md](2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md)
Next: [2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md](2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md)
Outcome: Recorded the user's bounded option 1 decision, retained exact inherited and unloaded Gorilla WebSocket v1.4.2 without metadata changes, revalidated every exception guard, and prepared the next bounded P7 group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-14 explicit selection of
option 1 for exact-path `github.com/gorilla/websocket`. Retain exact selected,
inherited, unloaded v1.4.2 without dependency metadata changes under the
bounded security and upstream-test exceptions below. Do not implement a
dependency change, audit another group, or begin P8 in this decision-recording
move.

# Authorized Roadmap

P2A-P6 and all earlier P7 outcomes remain final. P7 is active only for this
Gorilla product decision; P8 remains queued. This session may record one user
choice and prepare one bounded follow-up, but may not implement that choice or
combine another dependency group.

# Authorized Product Decision

On 2026-09-14 the user explicitly selected option 1 with the recommended
bounds: retain exact selected `github.com/gorilla/websocket v1.4.2` as an
inherited, unloaded selection without changing `go.mod` or `go.sum`. Accept
only the already documented GO-2026-6278 weak `math/rand` WebSocket client-mask
behavior and upstream full-source Go 1.26 cross-test lifecycle race. This does
not accept any new or independently discovered defect.

The exception is non-transferable and valid only while exact v1.4.2 and its
sole incoming mvn-pom-mutator v0.2.3 edge remain unchanged, the complete
project load contains zero Gorilla WebSocket packages, the module remains
runtime-unreachable, and no new advisory or independent disqualifier appears.
Revalidate and record those guards. Direct import or loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent disqualifier expires the exception and requires a fresh dependency
and product decision before merge.

Do not add a direct target edge, select v1.5.3 or another release merely to
silence an advisory range, change mvn-pom-mutator or GoConvey, reopen the final
GopherJS edge guard, raise the Go floor, authorize a parent/removal evaluation,
or select a patch, fork, replacement, unreleased commit, alternate path,
unrelated-module move, or dependency implementation commit. Existing GopherJS,
Enterprise Certificate Proxy, and GAX exceptions remain separate. Do not stop
or ask for this same Gorilla WebSocket decision again while all guards hold.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves the
canonical public active unarchived non-fork BSD-2-Clause repository. Serious
stable releases are v1.4.2, v1.5.0, v1.5.1, and v1.5.3. GitHub marks v1.5.2 a
prerelease. Its proxy identity is original commit `1bddf2e0`, but the forge tag
was overwritten to `9ec25ca`; unreleased main later retracts v1.5.2 for that
reason.

V1.4.2, v1.5.0, and v1.5.3 declare Go 1.12 and have standard-library-only
complete minimal source/test closures under exact Go 1.26.7 and Go 1.18.10.
V1.5.1 and v1.5.2 declare Go 1.20 and have `x/net` closures; their module zips
also have broken default vendor mode. V1.5.2 imports Go-1.20-only
`http.NewResponseController` and fails Go 1.18 compilation. V1.5.3 is the
highest stable floor-preserving release, but it is not security-qualified.

Fresh primary GO-2026-6278/GHSA-w67g-5rqw-f597 evidence describes weak PRNG
use for WebSocket client masks, marks versions before v1.5.3 affected, and
labels v1.5.3 fixed. Release source disproves that boundary: v1.5.3 explicitly
uses `math/rand`, and a seed-controlled independent fixture proves repeatable
mask keys. Security commit
`d67f41855da42d7bccd9ef050c49f7e54e783b95` changes production to
`crypto/rand` only after v1.5.3 on unreleased main. The original proxy v1.5.2
tree contains the fix despite the advisory range, but is a prerelease and
floor-ineligible. Selecting v1.5.3 merely makes the current range green without
fixing the behavior and is not an offered solution.

All packages, examples, tests, benchmarks, generated/build-tag files, exported
APIs, and relevant connection behavior were inspected. Independent fixtures
cover handshake success/failure, headers/auth/cookies, origin/subprotocol,
compression, fragmentation/control interleaving, malformed/oversized frames,
close/deadline/error behavior, prepared messages, pools, deterministic
masking, supported concurrency, and cleanup under both SDKs. A repeated full-
source Go 1.26 race selection finds an upstream cross-test lifecycle defect in
v1.4.2 and v1.5.3: a test handler logs after its owning test returns. The
independent supported-use fixture is race-clean.

Selected remains exact v1.4.2 solely through direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3`; `go mod why -m` is
negative, source imports are zero, and zero target packages load. Exact
v1.4.2/v1.5.0/v1.5.3 projections add only a redundant target requirement and
tidy byte-identically to base; v1.5.1/v1.5.2 additionally alter unrelated
`x/*` selections. No candidate was implemented.

The unchanged project has 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed packages, 41 loaded modules, 1,067
`go.sum` lines, and a 432-line tidy projection. `go.mod` and `go.sum` remain
byte-identical at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 project and full-preflight gates pass; applicable Go 1.18.10
gates retain only the two accepted shell wording failures. Accepted quality
remains 27/27 Q0-Q2 PASS at L2 with scorecard
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

# Role And Boundaries

This is a decision-recording move, not a renewed technical audit or dependency
implementation. The user has explicitly selected and bounded option 1; do not
ask for that decision again or broaden it into direct use, another version,
parent/removal, GopherJS-edge, floor, patch, fork, replacement, unreleased-
commit, alternate-path, or unrelated-module authorization.

## Guarded Decisions

P2A-P6 and all earlier P7 outcomes are final. The user's 2026-09-13 GopherJS
option 1 decision remains valid only for exact
`v0.0.0-20181017120253-0766667cb4d1`, its sole GoConvey v1.6.4 edge, zero
loaded packages, runtime unreachability, and no new advisory or independent
disqualifier. Enterprise Certificate Proxy v0.2.1 retains only its separate
exception for its sole Viper v1.15.0 edge under the same zero-load and
unreachable guards. GAX v2.7.0 retains only its two exceptions while zero GAX
packages load. Revalidate these guards before recording a choice and stop for
a fresh owning decision if any expired.

Option 1 authorizes only a decision record and no dependency edit. Do not
transfer another dependency's exception or reopen an earlier decision.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, guarded exceptions, and
`./codex-dev-start.sh --check`. Read this archive, the answered Gorilla
evaluation, rolling handover, roadmap, `go.mod`, `go.sum`, and referenced
release, quality, compatibility, evidence, and lifecycle contracts. Earlier
outcomes are final.

# Three Moves

First, revalidate only the unchanged exact selection, sole incoming edge,
negative why result, zero Gorilla/GopherJS/ECP/GAX package loads, module hashes,
runtime unreachability, and current advisory state. Reuse the completed audit;
do not repeat or broaden it. Second, record the exact option 1 exception,
accepted findings, guards, expiration triggers, and no-change result in the
roadmap and rolling handover. Third, answer this archive and prepare one
reciprocal NEXT archive for the next bounded P7 mission without executing it.

# Automatic Handoff

Run the launcher contract and applicable no-change gates, create the required
local `docs: prepare next agent session` commit, and stop. Do not launch the
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Recorded the user's 2026-09-14 option 1 selection and retained exact-path
`github.com/gorilla/websocket v1.4.2` as an inherited, unloaded module without
changing `go.mod` or `go.sum` and without manufacturing a dependency or
implementation commit.

The decision accepts only the two findings already established in the answered
dependency evaluation: GO-2026-6278's weak `math/rand` WebSocket client-mask
behavior and the upstream full-source Go 1.26 cross-test lifecycle race in
which `cstHandler` logs after its owning test returns. No new or independently
discovered defect is accepted.

Every guard was revalidated from clean decision HEAD `778d578` with a freshly
unpacked exact Go 1.26.7 distribution whose archive and binary SHA-256 values
are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
It ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.
Project selection remains Gorilla WebSocket v1.4.2, mvn-pom-mutator v0.2.3,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, GoConvey v1.6.4,
Enterprise Certificate Proxy v0.2.1, Viper v1.15.0, and GAX v2.7.0. The
unchanged Gorilla graph path is main -> direct mvn-pom-mutator v0.2.3 ->
selected v1.4.2, and the mvn-pom-mutator edge remains the target's sole
incoming edge. `go mod why -m` remains negative, and repository Go source has
zero direct target imports. Together with the zero package load below, this
revalidates that Gorilla WebSocket remains runtime-unreachable in the current
project.

The complete project test load remains 429 entries, including 197 module-backed
packages across 41 loaded modules, with exactly zero Gorilla WebSocket,
GopherJS, Enterprise Certificate Proxy, or GAX packages. Their exact selected
versions remain unchanged; GopherJS retains only its sole GoConvey v1.6.4 edge,
and Enterprise Certificate Proxy retains only its sole Viper v1.15.0 edge.
All four `go mod why -m` results remain negative, repository Go source has
zero direct imports of all four targets, and the project still selects 234
modules with 3,599 graph edges.

`go.mod` and `go.sum` remain byte-identical to clean decision HEAD at
SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
`go.sum` remains 1,067 lines. Exact Go 1.26.7 module verification, build,
full count-1 tests, full race tests, vet, and the 62-control launcher lifecycle
suite pass. No changed-selection scorecard run applies; the accepted 27/27
Q0-Q2 PASS at L2 remains unchanged.

Fresh primary vulnerability data still contains 1,398 module records, has
index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and was last modified 2026-09-10T16:28:28Z. The exact module entry contains
only GO-2020-0019, whose fixed boundary v1.4.1 excludes selected v1.4.2, and
selected-affecting GO-2026-6278. The latter remains unwithdrawn at record
SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`
with the same pre-v1.5.3 published range and GHSA alias. No new advisory or
independently observed disqualifier arose during this guard-only revalidation;
the completed technical audit was not repeated or broadened.

The non-transferable Gorilla WebSocket exception remains valid only while exact
selected v1.4.2 and its sole incoming mvn-pom-mutator v0.2.3 edge remain
unchanged, the complete project load contains zero target packages, the module
stays runtime-unreachable, and no new advisory or independent disqualifier
appears. Direct import or loading, runtime reachability, a target version or
incoming-edge change, or a new advisory or independent disqualifier expires
the exception and requires a fresh dependency and product decision before
merge. The GopherJS, Enterprise Certificate Proxy, and GAX exceptions remain
separate under their own guards. Do not request this same Gorilla WebSocket
decision again while all guards hold.
