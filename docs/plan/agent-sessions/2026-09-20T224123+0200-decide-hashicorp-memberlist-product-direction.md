# Agent Session: Decide Hashicorp Memberlist Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T224123+0200-decide-hashicorp-memberlist-product-direction`
Created: `2026-09-20T22:41:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `0f109c3fe885a097e911e934c55bb6de7624d1a2605c49075e36c5750b5218b3`
Previous: [2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency.md](2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency.md)
Next: [2026-09-20T230357+0200-study-hashicorp-serf-memberlist-removal.md](2026-09-20T230357+0200-study-hashicorp-serf-memberlist-removal.md)
Outcome: Authorized option 3 after every guard passed: a separate measurement-only Serf/owning-graph removal study; memberlist, Serf, product source, the Go floor, and dependency metadata remain unchanged.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one fresh bounded product decision for selected
exact-path `github.com/hashicorp/memberlist v0.3.0`. Its defensive evaluation
is complete: the selected release is affected by HCSEC-2026-18 /
CVE-2026-14362, v0.6.0 is first fixed, and every fixed stable release requires
Go 1.25. Choose exactly one option below, record its precise authorization and
expiry bounds, and prepare one coherent successor without executing it. Do not
repeat or extend the security investigation, reproduce the advisory, infer an
exception from zero loading, implement before authorization, combine another
dependency group, or begin P8.

# Defensive Decision Scope

This is an authorized defensive software-supply-chain product decision for the
user's local repository. Use only the completed public vendor/CNA facts,
static remediation identity, release and module metadata, and project graph
measurements recorded below. Do not create, reproduce, simulate,
operationalize, or optimize an exploit or proof of concept. Do not craft
packets or attack traffic, scan or contact hosts, exercise production or non-
public systems, inspect credentials or key material, test a security boundary,
or recreate the advisory's failure condition. No further exploitability or
protocol-security analysis is needed for this choice.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and every target-specific retained-module decision through exact inherited,
unloaded mdns v1.0.4. Every earlier outcome and lifecycle ancestor is final.
No earlier exception transfers to memberlist. Preserve every guarded selection
and stop for its owning decision if any guard differs. P8 remains queued.

Selected memberlist v0.3.0 is inherited only through the exact request from
Serf v0.9.6; historical Serf v0.8.2 requests memberlist v0.1.3. Memberlist has
no direct root or repository Go import, its why result is negative, production
and complete-test closures load zero target packages, and it is runtime-
unreachable in this project. These facts bound present exposure but do not
qualify the affected release or authorize retention.

# Measurements At Start

The defensive evaluation began from clean ordinary and ignored state at HEAD
`8a292c9a10f25aeae5be9e30693a7ed27f91912d`, parent
`a1dee3a81a71e0e6656de09fef603b94cc135bab`, tree
`017eaba0cc14d1442f7bea5bb7648ca290438df2`. That commit changed exactly the
launcher, memberlist evaluation archive, rolling handover, and roadmap. Its
reciprocal 225-archive chain, latest Google UUID implementation ancestry, and
launcher check passed. Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Verify the new handoff rather than assuming these facts.

The unchanged project has 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, and the recorded 432-line tidy projection.
`go.mod`/`go.sum` SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No dependency implementation or metadata commit exists; accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
public active unarchived non-fork MPL-2.0 repository. The exact path has 24
stable releases v0.1.0-v0.7.0; v0.7.0 is latest. The suffix tag
v0.3.1-metrics-labels is a prerelease. There is no retraction, module
deprecation, redirect, alternate major line, or qualifying branch/fork.
Inspected proxy archives match Git and sumdb verifies them. Serious tags are
lightweight and unsigned while their commits are GitHub-verified.

Selected v0.3.0 is verified commit
`923f1b205dd4653f2ea35e1c9531088e52053aa0`, parent
`123f3fbfeacd70bbe467775ab563a2b87e8c5cd8`, tree
`5c0926d642da064035733e5a591214269c223805`, time
2021-11-12T16:15:55-06:00. Its proxy ZIP SHA-256 is
`77e9fd374266825d7875a21a099cd80675bf9c8d469290f79c0a41bccba9df30`;
Git/proxy contents match and sumdb verifies it. It declares Go 1.12 and its
complete imported closure preserves Go 1.18. Its Go-1.18 source/test and race
compilations plus vet pass, but exact Go 1.26.7 cannot link its root test
binary because old x/net references removed `syscall.recvmsg`.

HashiCorp bulletin HCSEC-2026-18 states that memberlist through v0.5.4 is
affected and v0.6.0 fixes CVE-2026-14362. The published HashiCorp CNA record
also bounds affected versions below v0.6.0. Exact OSV and GitHub advisory
feeds are still empty and do not override the primary vendor evidence. Static
inspection identifies fix commit/v0.6.0
`371698b4dc493ecf2fc94c4383c75b2c9358f6c9`, parent
`0000b77c906dc53889a10de7a12fe4d146e84c50`, tree
`9009580c109364895d764539ede0498fc086254b`. It changes only `net.go` and
`net_test.go`; patch SHA-256 is
`f40e2a2034c07db2df2b2b76e76f09e3f67cca1289fe66a986d124a057cdb020`.
No advisory failure was dynamically reproduced.

First fixed v0.6.0 and latest fixed v0.7.0 both declare Go 1.25.0. Their
complete source/test closures load dependencies declaring Go 1.24-Go 1.25,
and Go 1.18.10 rejects their directives before loading. Under exact Go 1.26.7
both fixed releases pass complete package listing, source/test and race
compilation without executing test bodies, vet, and applicable cross-
compilation. Their existing security-regression test bodies were deliberately
not executed.

MVS selects v0.3.0 only through Serf v0.9.6, with the historical Serf
v0.8.2-to-v0.1.3 request also exact. There is no direct root/import/load. A
redundant v0.3.0 root changes no selection and tidy returns to the common base
projection. Direct fixed roots are not isolated target moves: v0.6.0 raises
the main Go directive to 1.25.0, yields 259 modules/3,743 edges/1,103 sums,
and changes four guarded selections; v0.7.0 yields 257/3,709/1,095 and changes
five guarded selections. Both preserve zero target load but make project
tests/vet fail existing non-constant format calls under Go 1.25. Their tidy
projections retain Go 1.25, dozens of selection-diff lines, and non-base
hashes. No projection was applied.

All 23 earlier guarded selections and 164 incoming edges remain exact at
snapshot SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
All guarded why/import/load results remain negative or zero. Fresh primary
vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
with only the recorded Gorilla and go-retryablehttp exact-version pairs among
earlier guards. No earlier guard expired.

# Role And Boundaries

Choose exactly one:

1. Authorize a separate integrated Go-1.25/fixed-memberlist migration. The
   latest fixed exact-path candidate is v0.7.0; v0.6.0 is first fixed. The
   follow-up must choose and fully qualify one exact release, explicitly own
   the Go-floor change, repair or decide the project's stricter-format
   failures, obtain fresh owning decisions for every changed guarded module,
   preserve API/CLI and runtime behavior, and pass complete source/test,
   platform, MVS, vulnerability, and project gates before implementation can
   be retained. This decision session must not edit the floor or metadata.
2. Retain exact selected, inherited, unloaded memberlist v0.3.0 without
   changing product source, `go.mod`, or `go.sum`, under a new target-specific,
   non-transferable risk decision. This is the smallest current graph change
   but explicitly accepts the primary vendor advisory and exact-Go test-link
   failure only within the verified zero-load state. It must expire if exact
   v0.3.0, either Serf request, direct root/import/load/runtime reachability,
   an earlier guard, or the no-new-advisory-or-defect condition changes. It
   must not be described as a qualified or secure release.
3. Authorize a separate owning-parent/graph-removal study for exact Serf
   v0.8.2/v0.9.6 and the request population that preserves memberlist. The
   follow-up must independently measure every Serf upgrade, replacement, or
   removal and all transitive API, behavior, MVS, guarded-edge, vulnerability,
   and project effects. It may not change Serf, memberlist, or another guarded
   module in this decision session.

Do not combine options, manufacture a direct root, silently raise the floor,
select v0.6.0/v0.7.0 without the integrated qualification, retain the affected
release without an explicit product choice, patch source, promote a fork,
branch, prerelease, or alternate path, change Serf or another parent, move an
unrelated selection, transfer another exception, or implement anything before
the choice is explicitly authorized. If none is acceptable, stop P7 explicitly
without source or metadata changes.

# Required Reading

Read this archive, its answered memberlist evaluation, the answered mdns,
golang-lru, and go.net decisions/evaluations, rolling handover, roadmap,
`go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored status,
handoff HEAD/parent/tree and exact changed-file set, reciprocal archive chain,
latest Google UUID implementation ancestry, exact Go identity, both
memberlist/Serf requests and zero-load state, every earlier guarded selection/
request/why/import/load condition, module hashes, vulnerability guards, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, revalidate only the narrow decision guards and primary advisory record;
do not repeat the evaluation or run security-focused tests. Second, choose and
record exactly one product option with precise bounds; do not implement it in
this decision session. Third, update the roadmap and rolling handover, answer
this archive, prepare exactly one reciprocal NEXT mission matching the
authorized choice, and commit the documentation-only handoff. Do not execute
the successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit. Do not
create a dependency implementation commit in this decision session. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 3 is explicitly authorized. Prepare and perform one separate,
measurement-only owning-parent and graph-removal study for the exact historical
`github.com/hashicorp/serf v0.8.2` and `v0.9.6` graph vertices and the exact
request population that preserves memberlist. The study must independently
measure every serious Serf upgrade, replacement, and removal route and all
transitive API, CLI, runtime-behavior, MVS, guarded-edge, Go-floor,
vulnerability, build, test, platform, and project effects. It may recommend
one later product direction, but it may not implement a Serf, memberlist,
parent, guarded-module, product-source, `go.mod`, or `go.sum` change.

This authorization is deliberately not option 2: it does not qualify, secure,
or risk-accept selected memberlist v0.3.0, and zero loading is only the bounded
current-exposure premise for a non-mutating study. It is also not option 1: it
does not select v0.6.0 or v0.7.0, raise the Go floor, repair the project's
stricter-format failures, or authorize any integrated migration. The study is
the sole authorized successor and expires when its one coherent decision
handoff is committed; any implementation still requires a fresh explicit
owning decision.

The authorization remains valid only while selected memberlist is exact
v0.3.0; the exact Serf-v0.9.6-to-memberlist-v0.3.0 and historical
Serf-v0.8.2-to-memberlist-v0.1.3 requests remain unchanged; selected Serf is
exact v0.10.1; and the five recorded incoming Serf requests remain exact:
Consul API v1.1.0 -> Serf v0.8.2; Consul API v1.12.0,
`sagikazarmark/crypt v0.4.0`, and Viper v1.10.1 -> Serf v0.9.6; and Viper
v1.15.0 -> Serf v0.10.1. It also requires no direct memberlist root or
repository import, negative memberlist why, zero memberlist production and
complete-test loads, runtime unreachability, all 23 earlier guarded selections
and 164 incoming edges remaining exact, unchanged module hashes, and no new
memberlist advisory, independently observed defect, fixed-version boundary, or
Go-floor-compatible fixed stable release. Any target, Serf, or incoming-request
change; direct root/import/load or runtime reachability; earlier owning-guard or
module-hash change; or new advisory, independent defect, affected/fixed-range
change, or floor-compatible fixed stable release expires this authorization
and requires the owning fresh decision before further work.

Guard-only revalidation began from clean ordinary and ignored worktree state
on branch `codex/upgrade-quality` at HEAD
`636e6b4f886ee15947c954deee70365c9eab1873`, parent
`8a292c9a10f25aeae5be9e30693a7ed27f91912d`, tree
`43dbc70674569af72990e06d1b16c38cda6f6dd7`. That handoff changed exactly
the launcher, answered memberlist evaluation archive, this then-NEXT decision
archive, rolling handover, and roadmap. The evaluation handoff `8a292c9`,
parent `a1dee3a`, tree `017eaba`, retains its exact four-file defensive-scope
delta. The reciprocal 226-archive chain, latest Google UUID v1.4.0
implementation ancestry, and launcher check pass.

Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
ran first in `PATH` with `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, no
ambient `GOFLAGS`, `LC_ALL=C`, and `LANG=C`. MVS selects memberlist v0.3.0
only through the Serf v0.9.6 request; historical Serf v0.8.2 still requests
v0.1.3. There is no direct root, its `go mod why -m` result remains negative,
repository Go imports remain zero, and the 355-entry production and 429-entry
complete-test closures contain zero target packages. The project remains 234
modules, 3,599 graph edges, 197 module-backed complete-test entries across 41
loaded modules, 1,067 sum lines, and the 432-line tidy projection. Exact module
verification passes.

All 23 earlier guarded selections and their 164 incoming edges remain exact;
the sorted incoming-edge snapshot retains SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
All 24 target-plus-earlier why results are negative, repository imports are
zero, and production and complete-test loads contain zero guarded packages.
`go.mod` and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Fresh primary revalidation preserves the decision premise. HashiCorp bulletin
HCSEC-2026-18 still identifies memberlist through v0.5.4 as affected and
v0.6.0 as fixed. The published HashiCorp CNA record remains PUBLISHED, was last
updated 2026-07-08T19:40:16.119Z, encodes an affected semver range below
v0.6.0, and has response SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
The Go vulnerability module index remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV results remain empty for
memberlist and every guarded target except the recorded Gorilla WebSocket
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp
GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. The lagging empty memberlist result
does not override the primary records. No earlier guard expired.

No product source, dependency metadata, parent, guarded selection, direct
root, Go floor, or unrelated module changed; no dependency implementation
commit was created. No changed-selection scorecard applies, and accepted
quality remains 27/27 Q0-Q2 PASS at L2. The sole reciprocal successor is the
authorized bounded Serf/owning-graph removal study; it was prepared but not
executed. Final exact-Go module verification, build, canonical-`umask 022`
count-one tests, race tests, vet, launcher check, and all 62 launcher lifecycle
controls pass. Initial parallel test/race runs inherited `umask 077` and
reproduced only the already-recorded fixture-mode assertions; the canonical
reruns supersede them.
