# Agent Session: Study Hashicorp Serf Memberlist Removal

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T230357+0200-study-hashicorp-serf-memberlist-removal`
Created: `2026-09-20T23:03:57+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fda8dbf5b0ab9bba2d734cfb760f937469a1fe3d7519534e5da4aa00efcebfe6`
Previous: [2026-09-20T224123+0200-decide-hashicorp-memberlist-product-direction.md](2026-09-20T224123+0200-decide-hashicorp-memberlist-product-direction.md)
Next: [2026-09-21T000630+0200-decide-hashicorp-memberlist-final-direction.md](2026-09-21T000630+0200-decide-hashicorp-memberlist-final-direction.md)
Outcome: Neither owning-parent removal nor compatible Serf/parent modernization is a bounded viable remediation; recommend a fresh decision favoring the separately authorized fixed-memberlist integration over explicit affected-version risk acceptance.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only with the authorized measurement-only owning-parent and graph-
removal study for the exact historical `github.com/hashicorp/serf v0.8.2` and
`v0.9.6` graph vertices and the request population that preserves affected
selected `github.com/hashicorp/memberlist v0.3.0`. Independently measure every
serious Serf upgrade, replacement, and removal route and all transitive API,
CLI, runtime-behavior, MVS, guarded-edge, Go-floor, vulnerability, build, test,
platform, and project effects. Prepare one bounded product decision; do not
implement a route, change Serf/memberlist/a parent/another guarded module, or
begin P8.

# Authorized Roadmap

The 2026-09-20 memberlist product decision explicitly authorized option 3
only. It did not risk-accept, qualify, or describe memberlist v0.3.0 as secure;
it did not authorize the Go-1.25/fixed-memberlist migration; and it did not
authorize any source or dependency metadata change. Zero loading bounds the
current exposure and authorizes this study, but is not itself an exception.

This study may use fresh public repository, release, module, advisory, API,
and project-graph evidence plus ordinary non-adversarial build/test tools in
disposable trees. It must independently establish the exact Serf repository,
stable release population, selected/historical graph identities, request
population, complete source/test closures, APIs, ordinary behavior, and every
candidate projection needed to decide whether the memberlist-preserving graph
can be removed or replaced. It may recommend exactly one later direction, but
it may not retain an implementation or edit product source, `go.mod`, or
`go.sum`.

Do not repeat or extend the memberlist security investigation, reproduce
HCSEC-2026-18 / CVE-2026-14362, execute its regression, craft packets or
attack traffic, scan or contact hosts, exercise non-public systems, inspect
credentials or key material, test a security boundary, or recreate the
advisory failure condition. Use the completed affected/fixed range and static
remediation identity only as decision constraints. Ordinary Serf behavior
evaluation must use non-adversarial inputs and must stop if it would require
security-focused testing.

## Exact Authorization And Expiry Guards

The authorization is valid only while selected memberlist remains exact
v0.3.0; the exact Serf-v0.9.6-to-memberlist-v0.3.0 and historical
Serf-v0.8.2-to-memberlist-v0.1.3 requests remain unchanged; selected Serf
remains exact v0.10.1; and the five incoming Serf requests remain exact:
Consul API v1.1.0 -> Serf v0.8.2; Consul API v1.12.0,
`sagikazarmark/crypt v0.4.0`, and Viper v1.10.1 -> Serf v0.9.6; and Viper
v1.15.0 -> Serf v0.10.1.

It also requires no direct memberlist root or repository Go import, a negative
memberlist why result, zero memberlist production and complete-test loads,
runtime unreachability, all 23 earlier guarded selections and their 164
incoming edges remaining exact, unchanged module hashes, and no new memberlist
advisory, independently observed defect, affected/fixed boundary, or Go-1.18-
compatible fixed stable release. Any target, Serf, or incoming-request change;
direct root/import/load or runtime reachability; earlier owning-guard or module-
hash change; or new advisory, independent defect, affected/fixed-range change,
or floor-compatible fixed stable release expires this authorization. Stop for
the owning fresh decision if any guard differs.

The study authorization itself expires after one coherent documented study
and reciprocal product-decision handoff. It grants no later implementation
authority. Do not manufacture a direct root, patch source, promote a fork,
branch, prerelease, pseudo-version, redirect, alternate module path, or
different major line; silently raise the Go floor; transfer another exception;
or combine an unrelated dependency group.

# Measurements At Start

Decision revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at HEAD
`636e6b4f886ee15947c954deee70365c9eab1873`, parent
`8a292c9a10f25aeae5be9e30693a7ed27f91912d`, tree
`43dbc70674569af72990e06d1b16c38cda6f6dd7`. That handoff changed exactly
the launcher, memberlist evaluation and decision archives, rolling handover,
and roadmap. Its reciprocal 226-archive chain, latest Google UUID v1.4.0
implementation ancestry, and launcher check passed. Verify the new handoff
rather than assuming these facts.

Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The unchanged project has 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed complete-test entries
across 41 loaded modules, 1,067 sum lines, and a 432-line tidy projection.
`go.mod`/`go.sum` SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No dependency implementation or metadata commit exists; accepted quality
remains 27/27 Q0-Q2 PASS at L2.

All 23 earlier guarded selections and 164 incoming edges remain exact at
snapshot SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
All 24 target-plus-earlier why results are negative, repository imports are
zero, and production and complete-test closures load zero guarded packages.
MVS selects memberlist v0.3.0 only through the Serf v0.9.6 request; historical
Serf v0.8.2 still requests memberlist v0.1.3. Memberlist has no direct
root/import/load.

HashiCorp bulletin HCSEC-2026-18 and the PUBLISHED HashiCorp CNA record still
bound affected memberlist below v0.6.0 and identify v0.6.0 as first fixed.
The CNA response SHA-256 is
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Primary Go vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV results remain empty for
memberlist and every guard except the recorded Gorilla and
go-retryablehttp pairs. These facts constrain the study; do not reopen the
memberlist evaluation.

# Role And Boundaries

Resolve the complete exact-path Serf stable-release and repository identity:
proxy/sumdb/`go-import`/Git/forge agreement, tags, commits, parents, trees,
times, signatures, ancestry, repository status, license, retractions,
deprecation, redirects, module lines, and serious candidates. Distinguish
selected Serf v0.10.1 from historical v0.8.2/v0.9.6 graph vertices and explain
why the latter retain memberlist requests despite MVS selecting v0.10.1.

Independently inventory and verify the five incoming Serf requests and every
root-to-request path that preserves them. For every serious exact Serf upgrade,
replacement, or removal route, use disposable projections to measure selected
modules, graph edges, sums, tidy behavior, Go directives and complete closure
floors, loaded packages, direct/indirect roots, why/import state, and every
changed guarded selection or incoming edge. Do not treat a redundant direct
root as remediation. Any route requiring a Consul API, Viper, crypt, or other
parent change must identify that ownership and all collateral changes; it may
not silently apply them.

Inspect Serf packages, commands, exported API, actual project consumers,
ordinary lifecycle and behavior, platform/build-tag branches, tests, and
resource/concurrency boundaries only as needed to compare viable graph routes.
Run source verification, complete package listing, ordinary upstream tests,
repeats, race, vet, meaningful cross-builds, API/CLI comparison, project
tests, vulnerability checks, and applicable quality gates for serious
candidates under exact Go 1.26.7 and contained Go 1.18.10. Classify failures
precisely. Do not execute memberlist security-regression tests or add security-
focused fixtures.

The output must compare the measured routes and prepare exactly one fresh
bounded product decision. At minimum, distinguish: removal of the request
population that preserves historical Serf/memberlist vertices; a compatible
Serf/parent modernization route; and, only if neither is viable, the precise
remaining choice between separately authorized fixed-memberlist integration
and explicit affected-version risk acceptance. Do not choose or implement
that later product decision in this study session.

# Required Reading

Read this archive, the answered memberlist decision and evaluation, the
answered mdns, golang-lru, and go.net decisions/evaluations, the relevant
answered Consul API and owning-parent records, rolling handover, roadmap,
`go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored status,
handoff HEAD/parent/tree and exact changed-file set, reciprocal archive chain,
latest Google UUID implementation ancestry, exact Go identity, all memberlist
and Serf selections/requests, every earlier guarded selection/request/why/
import/load condition, module hashes, vulnerability guards, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, revalidate only the exact expiry guards and map the Serf/request
population without changing the worktree. Second, perform the bounded owning-
graph study in disposable trees and leave product source and dependency
metadata unchanged. Third, answer this archive, update the roadmap and rolling
handover, prepare exactly one reciprocal NEXT product-decision mission for the
measured outcome, and commit the documentation-only handoff. Do not execute
the successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit. Do not
create a dependency implementation commit, launch a successor, push, merge,
publish, release, stash, revert, bypass cleanup, remove the worktree, combine
another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Neither removal of the historical request population nor compatible Serf or
parent modernization is a bounded viable remediation. The sole owner,
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3`, is already its latest
stable release and its default-branch head. Removing it does remove the
historical Serf/memberlist graph, but the project imports its POM model in 21
source/test files and no inspected alternative is API-compatible. The
disposable removal therefore fails package loading and compilation before
behavior can be preserved. Every fixed Serf stable line, meanwhile, requires
Go 1.25 or newer, and a direct Serf root is removed by tidy, so it cannot
rewrite the owning historical requests durably.

The study recommends exactly one later direction: make a fresh bounded product
decision favoring the separately authorized integrated Go-1.25/fixed-
memberlist migration over explicit risk acceptance for affected memberlist
v0.3.0. This study does not make that product decision and grants no
implementation authority. Product source, `go.mod`, `go.sum`, Serf,
memberlist, every parent, every guarded module, and the Go floor remain
unchanged.

The study began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at HEAD
`36bf3a3f5a05594213e0ccd6f5c45a2e2f3bbf3a`, parent
`636e6b4f886ee15947c954deee70365c9eab1873`, tree
`76cc264732cd570e599d3620985d1492adb71b74`. That handoff changed exactly the
launcher, memberlist decision archive, this then-NEXT archive, rolling
handover, and roadmap. Its reciprocal 227-archive chain, latest Google UUID
v1.4.0 implementation ancestry, and launcher check passed. Exact Go 1.26.7
binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
and contained official Go 1.18.10 archive/binary SHA-256 values
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
were independently verified.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence agrees on
`github.com/hashicorp/serf`, HashiCorp's public active unarchived non-fork
MPL-2.0 repository at `https://github.com/hashicorp/serf.git`, default branch
`master`. Proxy and Git expose 37 stable releases from v0.1.0 through current
v0.11.0 plus prerelease `v0.9.5-metrics-labels`. There is no retraction,
module deprecation, redirect, alternate major module path, or repository-
supported replacement. All serious candidates retain the exact module line.
Proxy archives match Git byte-for-byte except that the proxy correctly omits
the vendored tree in v0.8.2/v0.8.6; sumdb verifies every candidate. V0.8.2 is
an annotated PGP-signed tag whose GitHub verification reports an unknown key;
the other serious tags are lightweight and their commits are GitHub-verified.

Serious release identities are:

| Release | Commit | Parent | Tree | Commit time |
| --- | --- | --- | --- | --- |
| v0.8.2 | `b89a09eb2bd13faaa5331bdfaa26e86c590ce2d5` | `4e336ccd7405828251a9321c5b1835f40c2a1ca2` | `3e14f1a99d78201f8aa45b23fa1369d387046bcd` | 2019-01-24T16:20:00-05:00 |
| v0.8.6 | `e3039418ff39cad48000f266ece8e608e6693b08` | `1d3fdf866a95236ca9766cc8916a7279b1919317` | `ee522e31328b480b92fa20d98d71c4815321e7bc` | 2020-07-20T17:06:43-05:00 |
| v0.9.6 | `d223c691d21225d7d3434753306e89c62a9e2440` | `707607cfccc8cb896a0a67223b04f66aea19c0da` | `a4a1b819da3ba5f1afb984064748ece9242f45cb` | 2021-11-12T16:27:52-06:00 |
| v0.9.8 | `a2bba5671b3e50db2010ff184c18172517d6d699` | `9b0b7a3a61a7afb66b8935005114496396556182` | `cc605a663dc25b78b04c8bc7bdc83a8ff68646c8` | 2022-05-12T16:02:22-04:00 |
| v0.10.1 | `e853b5653f183e1681a3d10da6177d7ce1bb6fe3` | `830be12c4440b8d8c3762ce3ab6a1a450d0544df` | `e6fe01f1d3da922733c21bc5a8037ca9f8a794fd` | 2022-10-04T12:07:28-06:00 |
| v0.10.2 | `93c69c0756fc2a2bc26fbcd53feb117250663952` | `0a77f05544f6f5ef8e7f5acff270a9d6ab10fe2b` | `42e6e034bdc0a71e404d63caa1dba2872ac82ef2` | 2025-01-14T09:20:06-05:00 |
| v0.10.3 | `b2a6337cfdefd1ba347343aa25dbd05f3b791ac9` | `a24f620f3bf817cfb0a8981ce1a8a4b8950ed4e1` | `6a63ea31fde105cd2fe93f449e250b0a7d6a58b6` | 2026-07-07T09:16:54-04:00 |
| v0.10.4 | `5eba0bee9e10994a8844d63f24d1f0d48a17075d` | `b2a6337cfdefd1ba347343aa25dbd05f3b791ac9` | `9550fe48034688f60838f30af78f532898b2eb63` | 2026-07-07T09:52:03-04:00 |
| v0.11.0 | `fe2acb9e29aa82fa66b17238f9a829f39c56d0ab` | `1b5eabc9532ec6c970a087b9cff6082b250c1548` | `e508c2921cef3b1554eae247fdcaa500f9847be2` | 2026-09-16T11:09:21-04:00 |

All older-to-newer ancestry checks pass. Proxy ZIP SHA-256 values for the same
sequence are `0f4316583d26fdc959204e58edacb6c17870f1e52984986613b569dcda7edcee`,
`15538cee17a14041037a01dc1b1a7832942fc82dc8309304648a136420ec33ed`,
`60a00a6dcb5428f4127e89801d715af6e6f4ab488f066a9b1b31ea941de4e79f`,
`58203f7378aa0246df6e82332a84cbf023129c22b8935f61c6727e358370e8bb`,
`661b6ad561f8b9f297e2b9cc07c1bd24932328b8809e7215ef496f8087e80d6e`,
`117a872d7878f74adaafbb7028721aa61679432945478af7dc5b2d5a38b52efe`,
`7e6d991e02cd6ec05bd0d59d9c548b69404fbb686a2512fe490ba5cb560a327a`,
`f00920188bf50dff870ece582021c383658833fcfc95a13486b2d213dd365501`,
and `fefaa045605983488772226604649120a920441cfad72e00af8bfa65a80a4579`.
Corresponding sumdb module/go.mod pairs are:

| Release | Module sum | go.mod sum |
| --- | --- | --- |
| v0.8.2 | `h1:YZ7UKsJv+hKjqGVUUbtE3HNj79Eln2oQ75tniF6iPt0=` | `h1:6hOLApaqBFA1NXqRQAsxw9QxuDEvNxSQRwA/JwenrHc=` |
| v0.8.6 | `h1:w2ZEHuK1297elT/WbZjUojVzpZA3BuPUusa9vdXXTjc=` | `h1:P/AVgr4UHsUYqVHG1y9eFhz8S35pqhGhLZaDpfGKIMo=` |
| v0.9.6 | `h1:uuEX1kLR6aoda1TBttmJQKDLZE1Ob7KN0NPdE7EtCDc=` | `h1:TXZNMjZQijwlDvp+r0b63xZ45H7JmCmgg4gpTwn9UV4=` |
| v0.9.8 | `h1:JGklO/2Drf1QGa312EieQN3zhxQ+aJg6pG+aC3MFaVo=` | `h1:TXZNMjZQijwlDvp+r0b63xZ45H7JmCmgg4gpTwn9UV4=` |
| v0.10.1 | `h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=` | `h1:yL2t6BqATOLGc5HF7qbFkTfXoPIY0WZdWHfEvMqbG+4=` |
| v0.10.2 | `h1:m5IORhuNSjaxeljg5DeQVDlQyVkhRIjJDimbkCa8aAc=` | `h1:T1CmSGfSeGfnfNy/w0odXQUR1rfECGd2Qdsp84DjOiY=` |
| v0.10.3 | `h1:gKKs51YMnqRwv0cunl9hgaKwfrDqduYZ++s739qgGJw=` | `h1:sN/NigiHpcLenTce71rIMVEvkRwNi2m9tszYt8i8dcY=` |
| v0.10.4 | `h1:TCQOrJXHZ1Xf80c4WBhMM9OwUFgDaIP0R+YvoQUKadI=` | `h1:l+s5Q1OSPWU6b9l9m7ODJzTp7mLevSaVzAI03Nka2F0=` |
| v0.11.0 | `h1:8PbIr0pQOHs5hpBAkbX0teVLckuLvG3qIZm7hZJO0gc=` | `h1:k7zXXBvdNnDT5F1xrC0uMJB/0FWf6erElwZNQq/SsBA=` |

The five incoming requests and every preserving root path are exact. The main
module reaches mvn-pom-mutator v0.2.3, which reaches the old bketelsen crypt
pseudo-version and Consul API v1.1.0 -> Serf v0.8.2. It separately reaches
Viper v1.10.1 -> Serf v0.9.6, Viper v1.10.1 -> Consul API v1.12.0 -> Serf
v0.9.6, and Viper v1.10.1 -> sagikazarmark/crypt v0.4.0 -> both Serf v0.9.6
and Consul API v1.12.0 -> Serf v0.9.6. The direct Viper v1.15.0 root requests
selected Serf v0.10.1. MVS selects the maximum version but retains historical
vertices and their requirements. Because Viper v1.15.0's pruned module graph
records Serf v0.10.1 without its descendants, the historical unpruned
v0.9.6 requirement remains the selecting request for memberlist v0.3.0;
v0.8.2 still records memberlist v0.1.3.

Release requests and floors close the compatible-upgrade question. V0.8.2-
v0.8.5 have no `go` directive and request memberlist v0.1.3/mdns v1.0.0;
v0.8.6 adds `go 1.12` without changing those requests. V0.9.6-v0.9.8 declare
Go 1.12 and request memberlist v0.3.0/mdns v1.0.4. Selected v0.10.1 declares
Go 1.12 and requests affected memberlist v0.5.0; its complete test closure
tops out at Go 1.17. V0.10.2 declares Go 1.19, requests affected memberlist
v0.5.2, and its closure reaches Go 1.20. First fixed tags v0.10.3/v0.10.4
declare Go 1.25.0, request memberlist v0.6.0, and their closures reach Go
1.25. Latest v0.11.0 declares Go 1.26.0, requests memberlist v0.7.0, and its
closure reaches Go 1.26. No fixed stable Serf release preserves Go 1.18.

API comparison across public `client`, `coordinate`, and `serf` packages finds
no change from v0.9.6 to v0.9.8 or v0.10.2 through v0.10.4. V0.8.2->v0.8.6
changes exported `UserEventSizeLimit`, adds its config field and prune APIs,
and changes CLI encryption-key guidance. V0.8.6->v0.9.6 makes compatible
query, validation, reconnect, and key-response additions. V0.9.8->v0.10.1
adds metrics labels and makes coordinate `Config` non-comparable. V0.10.2 adds
the msgpack time-format field and three mDNS CLI flags. V0.11.0 incompatibly
changes both exported metrics-label fields from the compatibility label type
to `go-metrics.Label`. Top-level CLI help and the 15-command population are
stable; all candidate CLIs build under an applicable SDK.

Static ordinary-behavior inspection confirms that Serf is a networked library
and CLI: `Create` opens memberlist transports and starts coalescing, query,
reconnect, reap, snapshot, and event goroutines; `Join`, `Leave`, and
`Shutdown` have distinct synchronization and teardown contracts; user events
and queries use size/queue limits; callbacks, channels, returned node data,
snapshots, RPC, mDNS, syslog, subprocess event handlers, signals, and
caller-owned configuration create explicit concurrency/resource boundaries.
The project imports and loads no Serf package and invokes no Serf command, so
none of that API, CLI, lifecycle, network, filesystem, subprocess, or resource
surface is reachable in the product.

Complete package listing resolves 8 packages for v0.8.6, 10 for v0.9.6 and
later, and the same usable 8-package set for v0.8.2 after excluding its broken
`depreqs` tool-import package. Exact Go 1.26.7 compilation fails v0.8.2-
v0.10.1 on old Darwin x/net syscall linkage; v0.8.2 also has the invalid
command-package import. V0.10.2-v0.11.0 compile. Go 1.18.10 compiles v0.8.6-
v0.10.2, while v0.10.3+ is rejected at its directive. Repeated coordinate and
test-support tests and race-instrumented compilation pass for every serious
candidate. Vet and four-platform cross-builds pass from v0.9.6 onward;
v0.8.x vet reports legacy test-goroutine/stringer findings and cross-builds
fail its old go.net or `depreqs` closure.

Full upstream tests were attempted only with ordinary inputs. Exact-Go native
Darwin runs could not acquire loopback listeners in the host sandbox. A local
Linux/arm64 exact-Go container independently passed the tested core Serf,
coordinate, client, command, and retry suites across the v0.9.6-v0.11.0
candidate line; full aggregates were limited by absent container syslog, IPv6
multicast/mDNS facilities, and, when packages ran in parallel, fixed-address
listener collisions. V0.8.x also
shows old-go.net build failures and legacy timing-sensitive tests. These are
platform/environment or obsolete-closure failures, not advisory reproduction.
No memberlist security regression, crafted packet, attack input, host scan, or
security-boundary test was run.

Disposable project projections establish the route effects before tidy:

| Serf root | Go | Modules / edges / sums | Selected memberlist | Guard changes | Project gates |
| --- | --- | --- | --- | --- | --- |
| v0.10.1 | 1.18 | 234 / 3,629 / 1,071 | v0.5.0 | 0 selections; +9 edges | verify/build/test/race/vet pass |
| v0.10.2 | 1.18 | 244 / 3,645 / 1,073 | v0.5.2 | 5 selections; +8 edges | pass, but closure floor is Go 1.20 and memberlist remains affected |
| v0.10.3 | 1.25.0 | 246 / 3,652 / 1,079 | v0.6.0 | 5 selections; +8 edges | fixed, but floor-incompatible |
| v0.10.4 | 1.25.0 | 246 / 3,652 / 1,079 | v0.6.0 | 5 selections; +8 edges | build passes; tests/race/vet fail stricter format checks |
| v0.11.0 | 1.26.0 | 246 / 3,648 / 1,079 | v0.7.0 | 5 selections; +8 edges | same project failures plus API break |

The five changed guards from v0.10.2 onward are errwrap, go-multierror,
go-sockaddr, golang-lru, and mdns. Tidy removes every manufactured Serf root
and restores selected Serf v0.10.1/memberlist v0.3.0. V0.10.1 returns exactly
to the common 234-module/3,557-edge/948-sum tidy projection at hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Later `go get` projections leave incidental indirect upgrades and, on the
fixed lines, the raised Go directive after the direct root disappears. A
redundant direct root is therefore not remediation.

Version replacements are also rejected. Replacing only selected v0.10.1 with
v0.10.4 leaves historical memberlist v0.3.0 selected, because its owning
requirements remain. Replacing both historical vertices with v0.10.4 selects
memberlist v0.6.0 but is version masquerading, raises Go to 1.25, and expands
the observed graph to 449 modules/20,732 edges after full materialization. It
is unsupported, non-tidy, and expressly forbidden.

Dropping the actual owning parent is the only projection that removes the
preserving population. It yields 156 modules and 2,218 edges, removes
memberlist and historical Serf v0.8.2/v0.9.6 while leaving selected Serf
v0.10.1 through direct Viper v1.15.0, removes 15 guarded selections and 91 of
164 guarded incoming edges, and leaves Go 1.18. Package listing and compilation
fail on the missing `mvn-pom-mutator/pkg/pom` provider; tidy immediately
re-adds v0.2.3. The project uses `pom.Model`, dependencies, plugins, profiles,
properties, parsing, mutation, merge, sort, and writing across 21 files.
Inspected `gopom`, `mvnparser`, and a local encoding/xml implementation are
non-drop-in product migrations; copying or replacing the package locally
would be a forbidden fork. This is a separate POM subsystem redesign, not a
bounded dependency remediation.

Fresh vulnerability guards remain unchanged. The PUBLISHED CNA response is
byte-identical at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
and the 1,402-record Go index remains byte-identical at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Exact OSV results are empty for every
Serf candidate and memberlist v0.3.0. Earlier guards retain only the recorded
Gorilla and go-retryablehttp pairs. Source-mode govulncheck finds no Serf or
memberlist symbol in base or direct-root project projections, consistent with
zero loading; feed emptiness does not override the primary memberlist record.

At study close, all 24 memberlist-plus-earlier why results remain negative,
repository imports and production/complete-test guarded loads remain zero,
and runtime unreachability holds. All 23 earlier selections and 164 incoming
edges retain snapshot SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
The unchanged project remains 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed test entries across 41
loaded modules, 1,067 sum lines, and the 432-line tidy projection. `go.mod` /
`go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No changed-selection scorecard applies; accepted quality remains 27/27 Q0-Q2
PASS at L2. The sole reciprocal successor is the bounded final memberlist
direction decision; it was prepared but not executed. Final exact-Go module
verification, build, count-one tests, race tests, and vet pass. All 15 quality-
audit meta-controls and all 62 launcher lifecycle controls pass. A focused raw
Q0-Q2 audit without the required external manual-evidence receipt reports the
expected 21 automated PASS and six manual-evidence-bound UNMEASURABLE rows,
with zero ratchet regressions; it does not replace the accepted clean-tree
manual-evidence-adjusted 27/27 L2 scorecard. Two preliminary meta-suite runs
were superseded: one used the noncanonical host Go binary and one selected the
Apple Python shim, whose cache write was sandbox-denied. The canonical exact-
Go/Homebrew-Python rerun in writable scratch passed.
