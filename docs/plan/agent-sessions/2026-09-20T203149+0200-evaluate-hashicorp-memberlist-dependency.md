# Agent Session: Evaluate Hashicorp Memberlist Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency`
Created: `2026-09-20T20:31:49+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e75f031049356b10eba461bd03e6bc1e21c183155940c4e60e8b2147d89f26b5`
Previous: [2026-09-20T201536+0200-decide-hashicorp-mdns-product-direction.md](2026-09-20T201536+0200-decide-hashicorp-mdns-product-direction.md)
Next: [2026-09-20T224123+0200-decide-hashicorp-memberlist-product-direction.md](2026-09-20T224123+0200-decide-hashicorp-memberlist-product-direction.md)
Outcome: No fixed exact-path stable release preserves Go 1.18; selected v0.3.0 is affected and its exact-Go test closure also fails to link, so metadata stayed unchanged and one bounded memberlist product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/memberlist v0.3.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Defensive Retry Scope

This is an authorized defensive software-supply-chain review of the user's
local repository. A first attempt from clean HEAD
`a1dee3a81a71e0e6656de09fef603b94cc135bab` was automatically stopped by a
cybersecurity policy check after collecting public advisory and release facts;
it made no tracked worktree change. Resume this same bounded session under the
narrower rules below rather than restarting the broad investigation.

Use only public release metadata, the public vendor advisory, static source and
remediation-diff inspection, the project module graph, and ordinary existing
build or test commands with non-adversarial inputs in disposable local trees.
Do not create, reproduce, simulate, operationalize, or optimize an exploit or
proof of concept. Do not craft malformed packets, generate custom attack
traffic, scan or contact hosts, exercise production or non-public systems,
inspect credentials or key material, or test a security boundary. Do not
recreate the advisory's failure condition. If completing the decision would
require any such activity, stop and record that limitation instead.

The interrupted attempt retained these provisional facts for narrow
confirmation: public vendor advisory `HCSEC-2026-18` / `CVE-2026-14362` reports
all memberlist releases through v0.5.4 as affected; v0.6.0 is the first fixed
release; and v0.6.0 declares Go 1.25. Confirm only the affected/fixed range,
the candidate Go directive and source/test closure, and the static remediation
identity from public primary material. Do not investigate exploitability. If
confirmed, no fixed memberlist release preserves the project's Go 1.18 floor,
so leave dependency metadata unchanged and prepare the bounded product
decision. Treat missing temporary output from the interrupted attempt as a
reason to recheck these few facts, not to repeat the broad evaluation.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded mdns v1.0.4. Every earlier outcome and lifecycle ancestor is final.
Evaluate only Hashicorp memberlist in this session; do not reopen or combine
another dependency group. P8 remains queued.

The authorized 2026-09-20 mdns option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/mdns v1.0.4` without product source
or dependency metadata changes. Its target-specific, non-transferable
exception accepts only the completed exact-Go linker/test, Go-1.18 IPv6,
protocol, packet/channel, timeout/cancellation, lifecycle/resource, nil/panic,
mutation/aliasing, concurrency, environment, allocation, API, MVS,
vulnerability, and related recorded findings. It remains valid only while
exact v1.0.4, the Serf-v0.9.6-to-v1.0.4 and historical
Serf-v0.8.2-to-v1.0.0 requests, the mdns-v1.0.0-to-go.net-v0.0.1 request, no
direct root or repository import, zero target loads, runtime unreachability,
every earlier guard, and the no-new-advisory-or-defect condition remain exact.
Any change requires the owning fresh decision.

The golang-lru v0.5.4, go.net v0.0.1, go-uuid v1.0.1, go-syslog v1.0.0,
go-sockaddr v1.0.0, go-rootcerts v1.0.2, go-retryablehttp v0.5.3,
go-multierror v1.1.0, go-msgpack v0.5.3, go-immutable-radix v1.3.1,
go-hclog v1.2.0, Errwrap v1.0.0, qualified go-cleanhttp v0.5.2, and every
other recorded exception or qualification remain separate under their exact
selection, incoming-request, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards and stop for the owning decision if
any expires. No earlier exception transfers to memberlist. Do not change a
guarded parent, the Go floor, or an unrelated module.

Selected `github.com/hashicorp/memberlist v0.3.0` is inherited through the
exact request from `github.com/hashicorp/serf v0.9.6`; historical Serf v0.8.2
requests memberlist v0.1.3. The current repository import search and
production and complete-test loads contain zero target packages, its why
result is negative, and there is no direct main-module root. These queue
observations and the physical MVS selection are not proof of repository
identity, release qualification, ancestry, floor, behavior, vulnerability
state, or suitability. Resolve them independently and do not add a direct
edge merely to alter MVS.

# Measurements At Start

The mdns decision recording began from clean ordinary and ignored state at
HEAD `9edb38c51e0067672e90224d8665bae580f92a05`, parent
`b641b55ec01f18a79ecf57b8e5e866310dd44562`, tree
`63a43a0d26649532e760560a9e114fea47feb8d1`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive chain, latest Google UUID implementation ancestry, and
`./codex-dev-start.sh --check` rather than assuming them.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 22
earlier guarded selections and 162 incoming graph edges plus exact mdns
v1.0.4 and its two incoming Serf requests; the mdns-v1.0.0-to-go.net request
also remains exact. All guarded why results remain negative, repository
imports are zero, and production and complete-test loads contain zero guarded
packages. The project remains 234 modules, 3,599 graph
edges, 355 production entries, 429 complete-test entries, 197 module-backed
entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
unapplied tidy projection. Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
mdns and every guarded target except the recorded Gorilla
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp
GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Mdns's exact GitHub global and
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
launcher signal-retention timing race; none is memberlist evidence.

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

Inspect public source, packages, exported API, platform/build-tag branches, and
ordinary runtime lifecycle only as needed for the dependency decision. Reuse
the interrupted attempt's conclusions where recorded. Run existing upstream
tests and non-adversarial compatibility checks without modifying them, plus
source verification, package listing, race, vet, and meaningful cross-builds
under both SDKs where still necessary. Do not add security-focused fixtures or
custom network, packet, cryptographic, credential, malformed-input, fuzzing, or
failure-condition tests. Do not expand ordinary functional review into a
protocol-security assessment. Classify any ordinary build or test failure
precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. For the advisory, record only public primary affected/fixed
version facts, the fixed release's Go floor, the project's zero-load and
runtime-unreachable state, and a static remediation-diff identity. Do not
perform exploitability analysis, dynamic reproduction, or custom reachability
testing.

# Required Reading

At start read this archive, the answered mdns decision and evaluation, the
answered golang-lru decision and evaluation, the answered go.net decision and
evaluation, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify the
recorded handoff identity, changed-file set, clean ordinary and ignored
status, exact toolchain identity, reciprocal archive history, every guarded
selection/request/load result, and `./codex-dev-start.sh --check`. Earlier
outcomes are final.

# Three Moves

First, narrowly confirm the retained public advisory range, first fixed
release, Go-floor incompatibility, repository identity, project zero-load, and
exact MVS facts without changing the worktree or reproducing the advisory.
Second, if those provisional facts hold, leave metadata unchanged and prepare
one bounded product decision comparing at least: raising the Go floor and
selecting the fixed line; retaining exact unloaded v0.3.0 under an explicit
target-specific risk decision; and removing or changing the owning parent
edge in a separately measured study. Do not choose or implement an option in
this session. Only if primary evidence disproves the provisional facts may a
Go-1.18-compatible fixed exact-path stable candidate proceed through ordinary
qualification. Third, update the roadmap and rolling handover, answer this
archive, prepare exactly one reciprocal NEXT decision mission, and commit the
handoff. Do not execute the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable `github.com/hashicorp/memberlist` release qualifies under
the bounded defensive decision. HashiCorp's public `HCSEC-2026-18` bulletin
states that memberlist through v0.5.4 is affected by CVE-2026-14362 and that
v0.6.0 is the first fixed release. Selected v0.3.0 is therefore affected.
Both fixed stable releases, v0.6.0 and current latest v0.7.0, declare Go
1.25.0 and have imported source/test closures reaching Go 1.24 and Go 1.25.
Neither preserves the project's Go 1.18 floor. Product source, `go.mod`,
`go.sum`, the Serf parents, the Go floor, and every unrelated selection remain
unchanged; no dependency implementation commit was created.

The retry began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at HEAD
`8a292c9a10f25aeae5be9e30693a7ed27f91912d`, parent
`a1dee3a81a71e0e6656de09fef603b94cc135bab`, tree
`017eaba0cc14d1442f7bea5bb7648ca290438df2`. That defensive-scope commit
changed exactly the launcher, this then-NEXT archive, rolling handover, and
roadmap. The reciprocal 225-archive chain, latest Google UUID v1.4.0
implementation ancestry, and launcher check pass. Exact Go 1.26.7 binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`. Contained official
Go 1.18.10 archive/binary SHA-256 values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the public,
active, unarchived, non-fork MPL-2.0 repository
`https://github.com/hashicorp/memberlist.git`, default branch `master`. The
proxy and Git expose 24 exact-path stable releases from v0.1.0 through v0.7.0;
`v0.3.1-metrics-labels` is a prerelease and was not promoted. V0.7.0 is proxy
latest. There is no retraction, module deprecation, redirect, alternate major
module line, or qualifying branch/fork substitute. Serious tags are
lightweight and therefore unsigned; their commits are GitHub-verified. Every
inspected selected/last-affected/fixed proxy archive matches its Git tag
byte-for-byte and sumdb verifies it.

Selected v0.3.0 is verified commit
`923f1b205dd4653f2ea35e1c9531088e52053aa0`, parent
`123f3fbfeacd70bbe467775ab563a2b87e8c5cd8`, tree
`5c0926d642da064035733e5a591214269c223805`, time
2021-11-12T16:15:55-06:00. Its proxy ZIP SHA-256 is
`77e9fd374266825d7875a21a099cd80675bf9c8d469290f79c0a41bccba9df30`;
sumdb records module/go.mod sums
`h1:8+567mCcFDnS5ADl7lrpxPMWiFCElyUEeW0gtj34fMA=` /
`h1:MS2lj3INKhZjWNqd3N0m3J+Jxf3DAOnAH9VT3Sh9MUE=`. Last affected v0.5.4
is verified commit `21a632a23b186477c158742fbc1b131999f3ab3f`, tree
`ba009916101d9e93a03baa2720a3711d93820587`. First fixed v0.6.0 is
verified commit `371698b4dc493ecf2fc94c4383c75b2c9358f6c9`, parent
`0000b77c906dc53889a10de7a12fe4d146e84c50`, tree
`9009580c109364895d764539ede0498fc086254b`, time
2026-07-07T09:08:50-04:00. Its proxy ZIP SHA-256 and module/go.mod sums are
`b2c7ebd03f36b0dd22604e65b73b1e87ee1ce75eb36a333b1b6666b1f54497ff`,
`h1:hhVDLQUzWkLaitLLSrxLLqSD2l2+qiOz1DMr5zb9EQQ=`, and
`h1:a2lqh8KICpm8JibWOmuld7DaA+9QU1YcUtTTTMAtt/M=`. Latest fixed v0.7.0
is verified commit `4e3e17fc20e38ac48982ab6f65fec243fd139ef9`, parent
`44cfdfb1d2ba6b459fd131579b86ebc0930b62b9`, tree
`1e2d57c8a19d11776ed294cd4692e91c6e43aac6`, time
2026-09-16T10:40:44-04:00. Its proxy ZIP SHA-256 and module/go.mod sums are
`a77ea0b0c1bc0120cdcaf5ba34b1e95971b566d05f9eba7ac36191aedd83ff4d`,
`h1:JfqTDFUIAzDEYKMhSc3Gpwe05zvSU3/cYtiZ3yW59TM=`, and
`h1:Qar5D5CgaQAb74gk8Ph/jVcATn4epSDOHOvbSKOLHwg=`. Selected and v0.5.4
are ancestors of both fixed releases.

Static inspection identifies the remediation exactly, without reproducing the
advisory. V0.6.0 itself is fix commit
`371698b4dc493ecf2fc94c4383c75b2c9358f6c9`, subject `limit remote state
header values (#357)`. It changes only `net.go` and `net_test.go`, with 67
insertions and one deletion; its email/binary patch SHA-256 is
`f40e2a2034c07db2df2b2b76e76f09e3f67cca1289fe66a986d124a057cdb020`.
The commit is absent from v0.5.4 and is an ancestor of v0.6.0 and v0.7.0.
No exploitability analysis, custom packet, malformed-input fixture, attack
traffic, dynamic failure reproduction, scan, credential access, or security-
boundary test was performed. Existing security-regression test bodies were
not executed.

The complete directive history independently closes the floor question.
V0.1.0-v0.1.5 have no `go` directive; v0.1.6-v0.5.0 declare Go 1.12;
v0.5.1-v0.5.3 declare Go 1.20; v0.5.4 declares Go 1.24.0; and fixed
v0.6.0/v0.7.0 declare Go 1.25.0. V0.6.0's loaded test closure includes
go-metrics/x-sys at Go 1.25 and go-msgpack/v2 plus miekg/dns at Go
1.24 or later. V0.7.0's closure similarly includes go-metrics, miekg/dns, and
x-sys at Go 1.25 plus Go-1.24 go-msgpack/v2. Contained Go 1.18.10 rejects both
fixed modules' `go 1.25.0` directives before package loading.

Under exact Go 1.26.7, v0.6.0 and v0.7.0 each expose the root library and
`internal/retry`; their complete source/test closures list successfully.
Source/test compilation without executing tests, race-instrumented
compilation, vet, and Linux amd64/arm64, Windows amd64, and FreeBSD amd64 test
cross-compilation pass. V0.6.0 resolves 221 production and 241 complete-test
dependency entries; v0.7.0 resolves 220 and 240. Full native test execution
was deliberately not run because it would execute the advisory regression.
Selected v0.3.0 declares Go 1.12 and its imported closure declares no more
than Go 1.12; its source/test and race-instrumented compilations plus vet pass
under Go 1.18.10. Exact Go 1.26.7 cannot link its root test binary because its
old x/net closure references removed `syscall.recvmsg`; `internal/retry`
still compiles. This ordinary closure incompatibility independently prevents
qualification but is not advisory reproduction.

Selected v0.3.0 is a library-only module with no `package main` and no Go
source build-tag branches. Its public root-package surface centers on
`Config` and the LAN/WAN/local defaults, `Create`, the `Memberlist` join,
update, membership, send, leave, health, protocol, and shutdown methods,
`Node`/`NodeStateType`, transport/address/packet abstractions and
`NetTransport`, delegate/event hooks, `Keyring`, the bounded-transmit broadcast
queue and interfaces, and label/logging helpers. `internal/retry` is internal
test support, not an externally importable API. `Create` opens listeners and
starts background goroutines/tickers, after which callers must not mutate the
supplied config; `Join` performs state synchronization; delegate callbacks
carry documented concurrency, blocking, and borrowed-data constraints;
`Leave` is repeat-safe and broadcasts without stopping listeners but panics
after shutdown; and repeat-safe `Shutdown` tears down the transport and
background work without broadcasting leave. Returned member/event node data
has explicit no-mutation or copy-first contracts, and oversized delegate
metadata makes `UpdateNode` panic. The remediation commit changes no exported
declaration. Because the project imports and loads no target package, none of
this library surface or lifecycle enters the product API, CLI, or runtime.

Project MVS selects v0.3.0 only because
`github.com/hashicorp/serf@v0.9.6` requests it; historical Serf v0.8.2
requests memberlist v0.1.3. The main module has no direct target root,
`go mod why -m` is negative, repository Go imports are zero, and the 355-entry
production and 429-entry complete-test closures load zero memberlist packages.
The unchanged graph remains 234 modules, 3,599 edges, 197 module-backed
complete-test entries across 41 loaded modules, and 1,067 sum lines. Exact
module verification passes.

A disposable direct v0.3.0 root changes no selection, preserves zero target
load and 234 modules, and produces 3,620 graph edges and 1,078 sum lines;
tidy removes the manufactured root and returns the recorded common projection
hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection was applied.

A disposable v0.6.0 root raises the main `go` directive to 1.25.0, produces
259 modules, 3,743 graph edges, and 1,103 sum lines while preserving 355/429
project loads and zero target packages. It changes guarded Errwrap v1.0.0 to
v1.1.0, go-multierror v1.1.0 to v1.1.1, go-retryablehttp v0.5.3 to v0.7.7,
and go-sockaddr v1.0.0 to v1.0.7, plus many unrelated modules. Its tidy
projection removes the direct root but retains Go 1.25.0, 234 modules, 3,574
edges, 956 sum lines, and 60 selection-diff lines at go.mod/go.sum SHA-256
`5bdcb30d8d8d3e12ac64491282601159882f37419cfd301bbfc2a486c3b531ec` /
`94162600c940f03235b2b9724241786832f0735aed053a95baa1e38a688432c2`.

A disposable v0.7.0 root likewise raises the floor, produces 257 modules,
3,709 edges, and 1,095 sum lines, and keeps 355/429 loads with zero target
packages. It changes the same four guarded selections plus golang-lru v0.5.4
to v1.0.2. Its tidy projection retains Go 1.25.0, 234 modules, 3,572 edges,
957 sum lines, and 60 selection-diff lines at hashes
`a7674f7203a85fd401e62f4094e695843fce94393f333a41f4b5cc32fbbfb7ad` /
`719db7f7bd37a4405dd211bc82d08c135e8f7637064e2c592ffd67fc61bb4c06`.
Both direct fixed-line projections verify, but ordinary project `go test` and
vet fail existing non-constant format calls after the Go-1.25 directive
activates stricter checking. These are project/floor integration failures, not
memberlist runtime failures.

Fresh primary Go vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Exact OSV and both exact GitHub advisory
queries are empty for v0.3.0; v0.5.4-v0.7.0 exact OSV results are also empty.
Those lagging empty feeds do not override HashiCorp's primary bulletin or its
published CNA record, which encodes an affected range below v0.6.0. Every
earlier exact OSV guard remains unchanged: only Gorilla WebSocket retains
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp retains
GO-2024-2947/GHSA-v6v8-xj6m-xwqh.

All 23 earlier guarded selections and their 164 incoming edges remain exact;
the sorted snapshot retains SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
All 23 why results remain negative, repository imports are zero, and
production/complete-test loads contain zero guarded packages. Base `go.mod`
and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No earlier guard expired and no new independent guarded defect appeared.

Because no source or dependency metadata changed, no changed-selection
scorecard applies and accepted quality remains 27/27 Q0-Q2 PASS at L2. P7
stops on the sole reciprocal memberlist product decision comparing a Go-1.25
fixed-line migration, explicit guarded retention of affected but unloaded
v0.3.0, and a separately measured Serf/graph-removal study. It was prepared
but not executed. Final unchanged-project exact-Go module verification, build,
count-one tests, race tests, and vet pass; the byte-exact reciprocal launcher
check and all 62 launcher lifecycle contract checks pass.
