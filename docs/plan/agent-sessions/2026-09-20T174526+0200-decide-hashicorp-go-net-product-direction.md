# Agent Session: Decide Hashicorp Go Net Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T174526+0200-decide-hashicorp-go-net-product-direction`
Created: `2026-09-20T17:45:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `33c0f113bc92666202ece5160fc62b0dc8d0ee24cc9ff576a5fc23a158f1965d`
Previous: [2026-09-20T171035+0200-evaluate-hashicorp-go-net-dependency.md](2026-09-20T171035+0200-evaluate-hashicorp-go-net-dependency.md)
Next: [2026-09-20T180252+0200-evaluate-hashicorp-golang-lru-dependency.md](2026-09-20T180252+0200-evaluate-hashicorp-golang-lru-dependency.md)
Outcome: Authorized option 1 for exact inherited, unloaded go.net v0.0.1 after every guard passed; product source and dependency metadata stayed unchanged, and one bounded golang-lru evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one fresh bounded product decision for selected
exact-path `github.com/hashicorp/go.net v0.0.1`. Its independent evaluation is
complete and found no qualifying exact-path stable release. Choose exactly one
of the options below, record the authorized direction and its precise bounds,
and prepare one coherent follow-up without executing it. Do not repeat the
audit, silently accept a defect, implement before authorization, combine
another dependency group, change the Go floor, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go-uuid v1.0.1. Every earlier outcome and lifecycle ancestor is
final. No earlier exception transfers to go.net. Preserve every guarded
selection and stop for its owning decision if any guard differs. P8 remains
queued.

Selected `github.com/hashicorp/go.net v0.0.1` exists only through the exact
request from `github.com/hashicorp/mdns v1.0.0`. It has no direct main-module
root or repository Go import; `go mod why -m` is negative; production and
complete-test loads contain zero target packages; and it is runtime-
unreachable in this project. These facts bound present project exposure but do
not qualify the release or authorize an exception.

# Measurements At Start

The evaluation handoff began from clean ordinary and ignored state at HEAD
`70867a4a7701f69cbc4551de3d4db53fd7d86fe2`, parent
`2774d4e54971b833c3eae3342ba158f21c1b139a`, tree
`b97eb952c930dd8baaa0c5e9313c811a6c52be52`. Its exact five-file delta,
reciprocal archive chain, latest Google UUID implementation ancestry, and
launcher check passed. Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The project starts this decision unchanged at 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 1,067 sum lines, and the
recorded 432-line tidy projection. Verify the new handoff rather than assuming
these facts.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
public active unarchived non-fork BSD-3-Clause repository. V0.0.1 is the only
proxy version, only exact-path stable tag, and proxy latest. It is a lightweight
tag at GitHub-verified commit
`afc3cb3a421746fc66dd55b09a270c750cf536ce`, merge parents
`104dcad90073cd8d1e6828b2af19185b60cf3e29` and
`b9a611abb799b370ad8aed695239fd086ae5a1fa`, tree
`0a681728dc14ccfa08d532e5a20e85280c2fdb04`, time
2019-01-18T17:37:16Z. Proxy and Git files match byte-for-byte and sumdb
verifies the release. There is no GitHub Release, prerelease, retraction,
module deprecation, redirect, alternate exact path, or `/v2` line. Branches,
forks, and current master are not stable candidates. Master
`b69938bb7c99c1b9deeece78cc99714a94ec933b` is four compliance-only commits
beyond v0.0.1, unreleased, and retains the release defects.

The release contains 273 regular files, 204 Go files, 46 test files, and 16
packages: context, dict, html, html/atom, html/charset, idna, internal/iana,
internal/icmp, internal/nettest, ipv4, ipv6, netutil, proxy, publicsuffix,
spdy, and websocket. It also contains build-ignored generators, generated
tables and platform syscall files, examples, benchmarks, and testdata; it has
no release CLI or modern fuzz target. Exported APIs cover the legacy context,
HTML/charset/IDNA/public-suffix, dictionary, IPv4/IPv6 socket, listener-limit,
proxy/SOCKS5, SPDY v3, and WebSocket surfaces.

V0.0.1 has only a module directive and no requirements. Its complete source
and test closure does not resolve under exact Go 1.26.7 or contained Go
1.18.10. Production html and html/charset still import obsolete
`code.google.com/p/go.net` and `code.google.com/p/go.text` paths; context,
html, ipv4, ipv6, and websocket tests retain obsolete self-imports. Therefore
complete package listing, build, count-one, two count-ten repeats, race, vet,
and test-symbol vulnerability scanning fail under both SDKs. The mdns-used
ipv4/ipv6 production packages build and ordinary mdns tests pass on Darwin,
but their own complete tests do not resolve. Those packages cross-build on
the tested amd64 systems and Darwin arm64, but fail Linux arm64 and `js/wasm`
under both SDKs because required platform helpers are absent. No complete
minimal source/test closure preserves Go 1.18, so the only stable release is
disqualified before behavior risk could qualify it.

Source/API inspection and independent fixture SHA-256
`46328fe4873c5266f6f8de9435dd96d105760401b20b54732dcbba858308d3d6`
record additional exact behavior. `LimitListener` panics for negative limits
and blocks all accepts at zero; WebSocket server handling panics when the
response writer lacks `http.Hijacker`; SOCKS5 string-wraps I/O errors and
loses `errors.Is` identity; `PerHost` retains caller-owned IP/network storage;
and proxy dialer registration mutates an unsynchronized process-global map.
WebSocket dialing has no context or timeout, entropy failures panic, and SPDY
allocates from peer-declared lengths. Other packages retain their documented
caller-owned connection/deadline/resource, buffer-aliasing, static-table,
allocation, and non-concurrent-use boundaries. The fixture passes count-ten
and race under both SDKs while reproducing the negative-limit panic,
non-Hijacker panic, IP aliasing, process-global registration, and SOCKS5 error-
identity loss. These findings do not replace the closure disqualification.

MVS selects v0.0.1 through the sole mdns v1.0.0 request. A disposable direct
v0.0.1 root changes no selected version or loaded package; it adds only one
root graph edge and the content checksum. Tidy removes the manufactured root.
The direct-root project projection passes exact-Go verification, build,
count-one, two count-ten repeats, race, vet, API/CLI compatibility, host
acceptance, and supported production cross-builds because no target package
loads. The project remains 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, and the recorded 432-line base tidy projection.
Base `go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No projection was applied.

Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and has no target entry. Exact-version OSV and both GitHub advisory feeds are
empty. Exact-Go production module/package/symbol scans of the mdns-used
packages are empty; complete test-symbol scans cannot load the broken tests.
Old-Go findings belong only to Go 1.18.10's standard library. Base and direct-
root project scans are identical at 30 module, 22 package, and 20 production/
test-symbol IDs, with no target assignment or trace. All 20 earlier guarded
selections, requests, negative why results, zero Go-source imports, zero
loads, runtime-unreachability conditions, and advisory states remain exact.

# Role And Boundaries

Choose exactly one:

1. Retain exact selected, inherited, unloaded go.net v0.0.1 without changing
   product source, `go.mod`, or `go.sum`, under a new target-specific,
   non-transferable exception. This is the recommended bounded option for the
   current zero-load graph. It must explicitly accept only the completed
   unresolved production/test closure, obsolete imports, Linux-arm64/wasm
   gaps, panic/error-identity/aliasing/global-state/timeout/resource and
   related recorded findings. It must expire if v0.0.1 or the sole mdns
   v1.0.0 request changes, a direct root/import/load or runtime reachability
   appears, an earlier guard changes, or a new advisory or independent defect
   appears.
2. Authorize a separate Go-1.18-compatible exact-path patch/fork/replacement
   design. That follow-up must establish ownership and provenance, repair the
   entire production/test closure and relevant behavior, preserve required API
   and consumer semantics, qualify platforms and vulnerabilities, and measure
   all MVS/project effects before any implementation may be retained. It may
   not silently select master, a fork, alternate module path, or a different
   major line.
3. Authorize a separate parent/graph-removal decision for mdns v1.0.0. That
   follow-up must independently qualify the mdns upgrade, replacement, or
   removal and every transitive selection/API/behavior effect. It may not
   change the guarded parent or remove go.net in this decision session.

Do not infer option 1 merely from physical selection or zero reachability.
Do not combine options, manufacture a direct root, select unreleased master,
patch source, change mdns, raise the Go floor, move an unrelated selection,
transfer another exception, or implement anything before the choice is
explicitly authorized. If none is acceptable, stop P7 explicitly without a
dependency or source change.

# Required Reading

This is a product-decision session, not a renewed audit. At start read this
archive, its answered go.net evaluation, the answered go-uuid decision and
evaluation, the rolling handover, roadmap, `go.mod`, and `go.sum`. Verify
branch, clean ordinary and ignored status, handoff ancestry and exact changed-
file set, reciprocal archive chain, exact Go 1.26.7 identity, the sole target
request and zero-load state, all 20 earlier guarded selection/request/why/
import/load conditions, module hashes, vulnerability guards, and
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
`github.com/hashicorp/go.net v0.0.1` as an inherited, unloaded module without
changing product source, `go.mod`, or `go.sum`. The new go.net-specific,
non-transferable exception accepts only the completed unresolved production
and test closure; obsolete `code.google.com/p/go.net` and
`code.google.com/p/go.text` imports; Linux-arm64 and `js/wasm` platform gaps;
negative-limit and non-Hijacker panics; zero-limit blocking; SOCKS5 error-
identity loss; `PerHost` caller-storage aliasing; unsynchronized process-
global dialer registration; absent WebSocket context and timeout control;
entropy-failure panic; peer-sized SPDY allocation; and the documented caller-
owned connection, deadline, and resource, buffer-aliasing, static-table,
allocation, non-concurrent-use, API, MVS, vulnerability, and related completed
findings. It accepts no new or independently discovered defect.

The exception remains valid only while exact v0.0.1 and its sole exact request
from `github.com/hashicorp/mdns v1.0.0` remain unchanged, no direct main-module
root or repository Go import is added, production and complete-test target
loads remain zero, the module remains runtime-unreachable, all 20 earlier
guards remain intact, and no new advisory or independent defect appears. A
target version or mdns request change, direct root/import/load, runtime
reachability, an earlier owning-guard change, or a new advisory or independent
defect expires the exception and requires the owning fresh dependency and
product decision before merge. The patch/fork/replacement and mdns parent/
graph-removal alternatives are not authorized. No exception transfers to
another target.

Guard-only revalidation began from clean ordinary and ignored worktree state
on branch `codex/upgrade-quality` at HEAD
`fa9af5e9b168520f024da6f2049065ad186735b8`, parent
`70867a4a7701f69cbc4551de3d4db53fd7d86fe2`, tree
`53de66fc8ea39ffdfa68611677564af70f487f8e`. That handoff changed exactly the
launcher, answered go.net evaluation archive, this then-NEXT decision archive,
rolling handover, and roadmap. The evaluation handoff `70867a4`, parent
`2774d4e`, tree `b97eb952`, retains its exact five-file delta. The reciprocal
220-archive chain, latest Google UUID implementation ancestry, and launcher
check pass.

Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, public proxy and sumdb, `LC_ALL=C`, `LANG=C`, and
`umask 022`. Go.net remains exact v0.0.1 through only
`github.com/hashicorp/mdns@v1.0.0 github.com/hashicorp/go.net@v0.0.1`.
There is no direct root. Its `go mod why -m` result remains negative,
repository Go imports remain zero, and the 355-entry production and 429-entry
complete-test loads contain zero target packages.

All 20 earlier guarded modules retain their exact selections and recorded
requests. All 20 why results remain negative, repository Go imports remain
zero, and production and complete-test loads contain zero guarded packages.
The complete-test load retains 197 module-backed entries across 41 loaded
modules, so go.net and every earlier guarded target remain runtime-
unreachable. The project remains 234 modules, 3,599 graph edges, 1,067 sum
lines, and the recorded 432-line tidy projection. `go.mod` and `go.sum` retain
SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact module verification passes.

Fresh primary vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
with no go.net record. Exact-version OSV results remain empty for go.net and
every guarded target except the recorded Gorilla WebSocket
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp
GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Both exact go.net GitHub advisory
feeds remain empty. No new advisory or independently observed defect appeared.

No dependency implementation or metadata commit was created. No changed-
selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS at
L2. The sole reciprocal successor is the bounded P7 evaluation of selected
exact-path `github.com/hashicorp/golang-lru v0.5.4`; it was prepared but not
executed.
