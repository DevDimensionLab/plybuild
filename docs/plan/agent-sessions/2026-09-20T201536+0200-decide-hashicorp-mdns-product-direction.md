# Agent Session: Decide Hashicorp Mdns Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T201536+0200-decide-hashicorp-mdns-product-direction`
Created: `2026-09-20T20:15:36+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c1627b13c98f7673dc5492bdce48d98d4135340321225886521241bb8499f409`
Previous: [2026-09-20T192912+0200-evaluate-hashicorp-mdns-dependency.md](2026-09-20T192912+0200-evaluate-hashicorp-mdns-dependency.md)
Next: [2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency.md](2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency.md)
Outcome: Authorized option 1 for exact inherited, unloaded mdns v1.0.4 after every guard passed; product source and dependency metadata stayed unchanged, and one bounded memberlist v0.3.0 evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one fresh bounded product decision for selected
exact-path `github.com/hashicorp/mdns v1.0.4`. Its independent evaluation is
complete and found no qualifying exact-path stable release. Choose exactly one
option below, record its precise authorization and expiry bounds, and prepare
one coherent follow-up without executing it. Do not repeat the audit, infer an
exception from zero loading, implement before authorization, combine another
dependency group, change the Go floor, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded golang-lru v0.5.4. Every earlier outcome and lifecycle ancestor is
final. No earlier exception transfers to mdns. Preserve every guarded
selection and stop for its owning decision if any guard differs. P8 remains
queued.

Selected mdns v1.0.4 is inherited through the exact request from Serf v0.9.6;
historical Serf v0.8.2 requests mdns v1.0.0. The latter is also the sole graph
request preserving guarded go.net v0.0.1. Mdns has no direct main-module root
or repository Go import, its why result is negative, production and complete-
test closures load zero target packages, and it is runtime-unreachable in this
project. These facts bound current exposure but do not qualify the release or
authorize an exception.

# Measurements At Start

The evaluation began from clean ordinary and ignored state at HEAD
`b641b55ec01f18a79ecf57b8e5e866310dd44562`, parent
`961c48582ca568f0867d91e7781c6fa0b970dcea`, tree
`a7d492e85e9b8897d3d058af97b4801781f04bac`. Its exact five-file delta,
reciprocal 223-archive chain, latest Google UUID implementation ancestry, and
launcher check passed. Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Verify the new handoff rather than assuming these facts.

The unchanged project has 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, and a 432-line tidy projection. `go.mod`/`go.sum`
SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No dependency implementation or metadata commit exists; accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Role And Boundaries

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
public active unarchived non-fork MIT repository. Exact-path stable releases
are v1.0.0-v1.0.7; v1.0.7 is latest. There is no retraction, deprecation,
redirect, alternate major module line, or qualifying branch/fork substitute.
Every proxy archive matches Git and sumdb verifies it. All tags are
lightweight and unsigned, while their commits are GitHub-verified.

Selected v1.0.4 is commit
`f504e0b3c01cf50e077ab6805546fa95ed9f2933`, parent
`34dc81282184e8c94b4ed4312182d67b99827a5e`, tree
`50afb5c00de8ee437c229be31e0d0c5c0671f713`, time
2021-04-14T13:17:39-05:00. It is an ancestor of v1.0.5-v1.0.7 and main.
Its proxy ZIP SHA-256 is
`f7d04f484c5b398018385ad954a922b0d94d0cb26315780bf995bedede166704`;
Git/proxy contents match and sumdb verifies it.

Selected contains one library package, three production files, two tests, and
no command, example, benchmark, fuzz target, testdata, generated file, or build
constraint. V1.0.0-v1.0.4 share the same API; v1.0.5 compatibly adds address-
family controls; v1.0.6/v1.0.7 add context/logger/IPv6-zone surface.

No release passes the complete gate. V1.0.0 has a vulnerable DNS closure,
real shutdown races, and a vet defect. V1.0.1 fixes the races but retains that
vet and vulnerability failure. V1.0.2-v1.0.4 preserve Go 1.18 but exact Go
1.26.7 cannot link their Darwin tests because old x/net references removed
`syscall.recvmsg`; v1.0.3/v1.0.4 also fail native Go 1.18.10 tests when IPv6
multicast routing is unavailable. V1.0.5 fixes that floor-SDK condition but
retains the exact-Go Darwin test-link failure. V1.0.6 tests cleanly but its
loaded miekg/dns closure requires Go 1.19. V1.0.7 declares Go 1.25, loads
Go-1.24 dependencies, and cannot compile under Go 1.18.

Independent fixture SHA-256
`73739d0dcb546dcb5f7edcdb645d9c964444543296f5f1d8e74f9da22011967f`
proves selected behavior. Accepting v1.0.4 would accept unchecked wrapped
ports; weak trailing-dot-only name validation; caller IP/TXT and returned TXT
aliasing; inconsistent records after exported-field mutation; nil receiver,
nil config, and nil Zone panics; dropping all but the last complete answer in
a multi-entry packet; panic on a closed result channel; silent drops for nil,
unbuffered, or slow channels; direct requester responses even on the multicast
path; caller QueryParam mutation; acceptance of unsolicited answers; timeout
without context/cancellation; no goroutine lifecycle join; partial socket-
setup/resource and interface/route dependencies; unsynchronized caller-owned
config/zone/service data; process-global logging and hostname/DNS interaction;
the internal UDP-address panic boundary; eight record allocations; and all
recorded API, error, concurrency, platform, and ownership findings.

The selected graph has only Serf v0.9.6 -> mdns v1.0.4 and historical Serf
v0.8.2 -> mdns v1.0.0. Direct v1.0.0-v1.0.3 projections broadly downgrade or
remove guarded selections and fail project loading. Direct v1.0.4/v1.0.5
roots preserve loaded project packages but manufacture mdns/DNS roots; tidy
removes them back to the common base. V1.0.6 moves unrelated x modules and
has a non-base tidy projection. V1.0.7 additionally raises the main Go
directive to 1.25 and makes project test/race/vet fail on stricter format
checks. No projection was applied.

Fresh 1,402-record primary vulnerability data has SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and no mdns entry. Every stable version's exact OSV result and both selected
GitHub advisory feeds are empty. V1.0.0/v1.0.1 have a reachable DNS advisory;
later target scans have no called target/closure ID. Base/v1.0.4/v1.0.5
project scans remain 30 module, 22 vulnerable-package, and 20 called-symbol/
test-symbol IDs with no mdns assignment or trace. V1.0.6/v1.0.7 reduce counts
only by moving unrelated x modules.

All 22 earlier guarded selections and their 162 incoming edges remain exact;
all why/import/load results remain negative or zero. Exact OSV remains empty
except for the recorded Gorilla and go-retryablehttp pairs. No earlier guard
expired and no new advisory or independent guarded defect appeared.

# Product Decision

Choose exactly one:

1. Retain exact selected, inherited, unloaded mdns v1.0.4 without changing
   product source, `go.mod`, or `go.sum`, under a new target-specific, non-
   transferable exception. This is the recommended bounded option for the
   current zero-load graph because it avoids a manufactured direct root,
   unrelated graph moves, and a Go-floor change. It must explicitly accept
   only the complete linker/test, IPv6, protocol, packet/channel, timeout/
   cancellation, lifecycle/resource, nil/panic, mutation/aliasing, concurrency,
   environment, allocation, API, MVS, and vulnerability findings above. It
   must expire if v1.0.4, either Serf request, the mdns-v1.0.0 -> go.net request,
   direct root/import/load/runtime reachability, an earlier guard, or the no-
   new-advisory-or-defect condition changes.
2. Authorize a separate Go-1.18-compatible exact-path patch/fork/replacement
   design. That follow-up must establish ownership/provenance; repair the
   complete Darwin/source/test closure, IPv6 behavior, packet/result delivery,
   multicast semantics, cancellation, lifecycle/resource cleanup, panic and
   ownership boundaries while preserving required exported API and consumer
   behavior; and qualify every platform, vulnerability, MVS, and project effect
   before implementation may be retained. It may not silently select v1.0.5,
   v1.0.6, v1.0.7, main, a fork, or an alternate path.
3. Authorize a separate parent/graph-removal decision for exact Serf
   v0.8.2/v0.9.6 and their request population. That follow-up must independently
   qualify every Serf upgrade, replacement, or removal and all transitive API,
   behavior, MVS, go.net-guard, and project effects. It may not change Serf,
   mdns, go.net, or another guarded parent in this decision session.

Do not combine options, manufacture a direct root, accept a floor-ineligible
release, patch source, change Serf or another parent, raise the Go floor, move
an unrelated selection, transfer another exception, or implement anything
before the choice is explicitly authorized. If none is acceptable, stop P7
explicitly without source or metadata changes.

# Required Reading

Read this archive, its answered mdns evaluation, the answered golang-lru and
go.net decisions/evaluations, rolling handover, roadmap, `go.mod`, and
`go.sum`. Verify branch, clean ordinary and ignored status, handoff HEAD/
parent/tree and exact changed-file set, reciprocal archive chain, latest
Google UUID implementation ancestry, exact Go identity, every guarded
selection/request/why/import/load/advisory condition, unchanged project hashes
and counts, mdns/Serf/go.net graph edges, and `./codex-dev-start.sh --check`.
Earlier outcomes are final.

# Three Moves

First, choose exactly one product option without repeating the completed
evaluation. Second, record only that authorization and its guards and prepare
the one matching successor without implementing it. Third, verify the
reciprocal handoff and make the required documentation-only commit.

# Automatic Handoff

Record exactly the chosen option and precise bounds in the roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal NEXT mission for
that authorized direction, and make the required local
`docs: prepare next agent session` commit. Do not create a dependency commit
in this decision session. Do not execute the successor, push, merge, publish,
release, stash, revert, bypass cleanup, remove the worktree, combine another
dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is explicitly authorized. Retain exact selected
`github.com/hashicorp/mdns v1.0.4` as an inherited, unloaded module without
changing product source, `go.mod`, or `go.sum`. This is an explicit bounded
product decision for the verified zero-load graph, not an inference from
physical selection. It avoids manufacturing a direct mdns/DNS root, moving
unrelated graph selections, or changing the Go floor.

The new mdns-specific, non-transferable exception accepts only the completed
exact-Go Darwin test-link failure through old x/net's removed
`syscall.recvmsg`; the Go-1.18 IPv6 multicast-routing test failure; unchecked
wrapped ports; trailing-dot-only name validation; caller IP/TXT and returned
TXT aliasing; inconsistent records after exported-field mutation; nil service
receiver, nil config, nil Zone, and internal UDP-address assertion panics;
last-only delivery from multi-entry packets; closed-result-channel panic;
silent nil, unbuffered, or slow result-channel drops; requester-directed
responses on the multicast path; caller `QueryParam` mutation; acceptance of
unsolicited answers; timeout without context or cancellation; absent goroutine
lifecycle joins; partial socket setup and resource cleanup; interface, route,
hostname, and DNS dependencies; unsynchronized caller-owned config, Zone, and
service state; process-global logging; eight record allocations; and the
recorded API, error, protocol, packet/channel, concurrency, platform,
ownership, environment, MVS, and vulnerability findings. The vulnerability
bounds include the empty exact v1.0.4 and project mdns results and the recorded
reachable DNS advisory in the historical v1.0.0/v1.0.1 closure. The exception
accepts no new advisory or independently discovered defect.

The exception remains valid only while exact selected v1.0.4, the exact
`github.com/hashicorp/serf@v0.9.6 -> github.com/hashicorp/mdns@v1.0.4`
request, the historical
`github.com/hashicorp/serf@v0.8.2 -> github.com/hashicorp/mdns@v1.0.0`
request, and the exact
`github.com/hashicorp/mdns@v1.0.0 -> github.com/hashicorp/go.net@v0.0.1`
request remain unchanged. It also requires no direct main-module root or
repository Go import, zero production and complete-test target loads, runtime
unreachability, every earlier owning guard remaining intact, and no new
advisory or independent defect. Any target selection or recorded request
change, direct root/import/load, runtime reachability, earlier-guard change, or
new advisory or independent defect expires the exception and requires the
owning fresh dependency and product decision before merge.

The exact-path remediation and Serf/graph-removal alternatives are not
authorized. V1.0.5, v1.0.6, v1.0.7, main, forks, alternate paths, direct
roots, patches, parent changes, unrelated selection moves, and a Go-floor
change remain unauthorized. No exception transfers to another target.

Guard-only revalidation began from clean ordinary and ignored worktree state
on branch `codex/upgrade-quality` at HEAD
`9edb38c51e0067672e90224d8665bae580f92a05`, parent
`b641b55ec01f18a79ecf57b8e5e866310dd44562`, tree
`63a43a0d26649532e760560a9e114fea47feb8d1`. That handoff changed exactly the
launcher, answered mdns evaluation archive, this then-NEXT decision archive,
rolling handover, and roadmap. The evaluation handoff `b641b55`, parent
`961c485`, tree `a7d492e`, retains its exact five-file delta. The reciprocal
224-archive chain, latest Google UUID v1.4.0 implementation ancestry, and
launcher check pass.

Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`. Mdns remains exact
v1.0.4 through the recorded Serf requests, and the mdns-v1.0.0-to-go.net edge
remains exact. It has no direct root, its `go mod why -m` result remains
negative, repository Go imports remain zero, and the 355-entry production and
429-entry complete-test closures contain zero target packages.

All 22 earlier guarded selections and their 162 incoming edges remain exact;
the sorted edge snapshot retains SHA-256
`73b2f342ab8b1f5334afedead3c49cee97ee0177bf2413ea13d10f5d45995c63`.
All 22 why results remain negative, repository Go imports remain zero, and
production and complete-test loads contain zero guarded packages. The
complete-test closure retains 197 module-backed entries across 41 loaded
modules, so mdns and every earlier guarded target remain runtime-unreachable.
The project remains 234 modules, 3,599 graph edges, 1,067 sum lines, and the
432-line tidy projection. `go.mod` and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and exact module verification passes.

Fresh primary vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no mdns record. Exact-version OSV
results remain empty for mdns and every guarded target except the recorded
Gorilla WebSocket GO-2026-6278/GHSA-w67g-5rqw-f597 and
go-retryablehttp GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Both exact mdns
GitHub advisory feeds remain empty. No new advisory or independently observed
defect appeared.

No dependency implementation or metadata commit was created. No changed-
selection scorecard applies, and accepted quality remains 27/27 Q0-Q2 PASS at
L2. The sole reciprocal successor is the bounded P7 evaluation of selected
exact-path `github.com/hashicorp/memberlist v0.3.0`; it was prepared but not
executed.
