# Agent Session: Decide Mwitkow Go-Conntrack Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T001118+0200-decide-mwitkow-go-conntrack-product-direction`
Created: `2026-09-23T00:11:18+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `2a28c3fcb8c7cbed504e3b1449753ff449ee550bd911bdbae13d9983ca868591`
Previous: [2026-09-22T231926+0200-evaluate-mwitkow-go-conntrack-dependency.md](2026-09-22T231926+0200-evaluate-mwitkow-go-conntrack-dependency.md)
Next: [2026-09-23T003641+0200-evaluate-oklog-ulid-dependency.md](2026-09-23T003641+0200-evaluate-oklog-ulid-dependency.md)
Outcome: Option 1 selected. Exact inherited, unloaded go-conntrack pseudo-version is retained under a target-specific, non-transferable unqualified exception; source and dependency metadata remain unchanged.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one product direction for graph-selected,
transitive, unloaded, unqualified exact
`github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223`.
Choose only one of the three authorized options below. This is one bounded
decision-recording group, not dependency implementation, owner research, a
new evaluation, another dependency group, or P8.

# Authorized Evaluation Result

The canonical exact-path stable release line is empty. Exact go-import maps to
public enabled unarchived non-fork Apache-2.0 repository
`mwitkow/go-conntrack`, but it has zero tags and zero GitHub Releases, and the
Go proxy stable list is empty. Selected
`v0.0.0-20161129095857-cc309e4a2223` and proxy `@latest`
`v0.0.0-20190716064945-2f068394615f` are pseudo-versions, not stables. Do not
promote either pseudo-version, `master`, a pull-request ref, fork, replacement,
or alternate path as a stable. No canonical stable exists or qualifies, so no
selection or dependency metadata changed.

Selected commit `cc309e4a22231782e8893f3c35ced0967807a33e`, tree
`8e653f1c675ad96001b221496975aa5858f0df09`, is an unsigned ancestor of
signed current master/latest commit
`2f068394615f73e460c2f3d2c158b0ad9321cadb`, tree
`a8074d12ad49ac88f1078e3942c1a95c8bfbc5f2`. Both use synthetic exact-path
module files without a Go directive, requirements, replacement, retraction,
or deprecation. There is no `/v2` or `/v3`. The selected proxy archive exactly
matches its 18-file Git tree. Preserve all recorded owner, license, commit,
tree, signature, ancestry, proxy, sumdb, archive, source, and module identities.

The selected source's historical Prometheus closure builds and passes count-
one tests and vet under exact Go 1.18.10 and Go 1.26.7, but repeated tests fail
under both SDKs because global metric state accumulates. Go 1.26.7 race testing
also finds concurrent historical `x/net/trace` event-map access. With the
project's actual Prometheus client v1.4.0/common v0.9.1 closure, full upstream
build/tests/race/vet fail under both SDKs at removed `prometheus.Handler` use;
the two library packages and all nine supported cross-build targets compile.
Preserve the recorded exported API, documentation, ownership/mutation,
determinism, concurrency, lifecycle, cleanup, errors, globals, platform,
network/TLS/file, cgo/build-tag/generated/embed, and Go-floor findings.

Exactly two graph edges request the selected pseudo-version: Prometheus common
v0.9.1 and historical common v0.4.1. Both genuinely import and use the target.
The current route is main -> direct mvn-pom-mutator v0.2.3 -> historical Viper
v1.10.1 -> go-metrics v0.3.10 -> common v0.9.1 -> target, with a second branch
through Prometheus client v1.4.0. The historical route continues common v0.9.1
-> Prometheus client v1.0.0 -> common v0.4.1 -> target. Target and requester
why results are negative, repository target imports are zero, and target
production, complete-test, module-backed load, and runtime relevance are zero.

Disposable exact gets are not retainable. The selected get keeps 234 modules
but promotes the target and seven already-selected package dependencies as
main roots, producing 82/1,076 metadata lines and 3,609 graph edges. The latest
pseudo get additionally selects `jpillora/backoff v1.0.0`, producing 235
modules, 83/1,079 metadata lines, and 3,610 edges. Both remain why-negative and
unloaded and tidy back to the exact common 52/948-line projection. Real state
remains 234/3,599/355/429/197/41/1,067 with the recorded module/tidy hashes and
Go 1.18 floor.

All 45 earlier guarded selections and 273 incoming edges remain exact at
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Every reflect2, concurrent, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, Cast, Viper, closed-exception, owner/request,
why/import/load/runtime, graph/module/tidy/Go-floor, advisory, qualification,
and compatible-route fact remains final under its own recorded guard.

Exact target OSV, GitHub global, repository advisory, and pinned govulncheck
target findings are empty. The unchanged project retains exact 30/22/20/20
non-stdlib govulncheck populations. Guard OSV retains only the recorded Gorilla
WebSocket and go-retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
GO-2026-6180. The 518,501-byte/1,402-record Go index and PUBLISHED 2,807-byte
memberlist CNA identities remain exact. Advisory absence is not qualification.
Final exact-Go verify/build/count-one/race/vet passes and accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, unloaded
   `github.com/mwitkow/go-conntrack
   v0.0.0-20161129095857-cc309e4a2223` without product-source or dependency-
   metadata changes under a go-conntrack-specific, non-transferable exception.
   Call it unqualified: the exact-path stable line is empty, the selection is
   a pseudo-version, and its ordinary upstream closure gates are not all green.
   Accept only the completed evaluation and define the exact expiry guards
   below.
2. Authorize exactly one later bounded measurement-only study named
   **Prometheus Common Go-Conntrack Ownership Study**. Its sole question is
   whether a later genuine supported tidy-stable selection along the existing
   Prometheus common owner/request route removes go-conntrack or requests a new
   qualified Go-1.18-compatible exact-path stable while preserving every
   project and earlier-decision guard. Do not run the study, evaluate or change
   Prometheus common/client/go-metrics/Viper/mvn-pom-mutator, add a root, change
   a selection, or reopen an earlier decision in this move; prepare exactly one
   reciprocal measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap, or
   P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical MVS selection,
transitive presence, negative why, zero loading, an active repository, or
advisory absence is not qualification. Do not call either pseudo-version
stable or qualified, select latest/master, add a direct target root, change a
requester, transfer an exception, or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/pseudo-version and both common v0.9.1/v0.4.1 requests;
- both genuine requester imports, complete current/historical routes, negative
  why/import/load/runtime facts, zero target root, and no supported stable
  owner/request route;
- the empty exact-path stable line, absent major lines, zero tags/releases,
  exact repository/owner/status/license/default branch, pseudo-version/commit/
  tree/signature/ancestry/archive/sumdb/module identities, and no replacement,
  retraction, deprecation, eligible alternate owner, or later stable;
- every selected-source API/behavior/closure/upstream/repeated/race/vet/cross-
  build/Go-floor finding, including actual-project Prometheus incompatibility;
- exact baseline and disposable-projection counts/hashes, tidy restoration,
  Go 1.18 floor, exact Go 1.26.7 identity, and all source/API/CLI/help/launcher/
  Make/quality contracts;
- all 45 earlier selections/273 incoming edges, all closed decisions and their
  separate expiry boundaries, and accepted 27/27 Q0-Q2 PASS at L2; and
- no new advisory, independent defect, release, owner, qualified exact-path
  stable, supported tidy-stable requester, or compatible genuine route.

Any path/version, request, requester import, owner route, root/why/import/load/
runtime, source/behavior/closure, graph/module/tidy/Go-floor, earlier guard,
advisory/finding, repository/release/owner, qualification, supported-owner, or
compatible-route change expires the exception and requires a fresh target
evaluation and product decision before merge. Option 1 authorizes no owner
study, direct root, dependency edit, workaround, or implementation.

# Role And Boundaries

This session records one decision only. Revalidate only the minimum continuity,
target/graph/advisory identities, launcher, and unchanged-project exact-Go
checks needed to record the chosen option safely. Do not repeat network/TLS
fixtures, upstream behavior work, or race reproduction. Do not launch a study
or successor. Preserve a clean worktree except for the bounded documentation/
launcher handoff, then make the required local handoff commit.

# Required Reading

Read this decision archive, the answered go-conntrack evaluation, answered
reflect2/concurrent evaluations, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, and kr/pty decisions/evaluations, kr/pretty and Cast owner records,
rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Verify the handoff
HEAD/parent/tree, exact changed set and ancestry, reciprocal archive chain,
clean ordinary/ignored state, launcher check, exact Go identities, module
hashes/counts/tidy projection, target requests/routes/why/import/load facts,
all 45 earlier guards/273 edges, closed-exception boundaries, and fresh
advisory identities.

# Three Moves

First, ask for or apply exactly one explicit option without reopening the
evaluation. Second, record only that option and its exact guards; change no
product source or dependency metadata. Third, update roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal successor only
if the chosen option requires one, verify containment, and make the required
local handoff commit without executing a successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen any earlier
decision, evaluate another dependency group, write outside the managed scratch
root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, unloaded
`github.com/mwitkow/go-conntrack
v0.0.0-20161129095857-cc309e4a2223` is explicitly retained without changing
product source, `go.mod`, or `go.sum` under a go-conntrack-specific,
non-transferable exception. It remains unqualified and is not described as
stable or secure. The canonical exact-path stable line is empty, the selected
identity is a pseudo-version, and its completed ordinary upstream closure
gates are not all green. Physical MVS selection, transitive presence, negative
why, zero loading, an active repository, and advisory absence are not
qualification or implicit acceptance.

The exception accepts only the completed repository/release/source/module,
API and ordinary behavior, closure/upstream/repeated/race/vet/cross-build,
owner/request/route, why/import/load/runtime, projection/graph/tidy/Go-floor,
earlier-guard/closed-exception, and advisory findings recorded by the answered
evaluation. It accepts no uncharacterized behavior, new advisory, independent
defect, future route, or direct/loaded use. It transfers or broadens no earlier
exception. Option 2's **Prometheus Common Go-Conntrack Ownership Study** was
neither authorized nor run, and option 3 was not selected.

### Exact Retention, Request, And Route Guards

Retention requires exact selected path and pseudo-version
`github.com/mwitkow/go-conntrack
v0.0.0-20161129095857-cc309e4a2223`. Main must retain no direct target root.
The only two graph requests must remain selected Prometheus common v0.9.1 and
historical Prometheus common v0.4.1, both requesting that exact pseudo-version.
Both requester source trees must continue genuinely importing and using the
target from `config/http_config.go` to construct their traced, named dial
context functions; neither edge may become metadata-only or acquire a
different request.

The complete current route remains main -> direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` -> historical Viper
v1.10.1 -> `github.com/armon/go-metrics v0.3.10` -> Prometheus common v0.9.1
-> target, with the parallel go-metrics -> Prometheus client v1.4.0 -> common
v0.9.1 branch. The historical route remains common v0.9.1 -> Prometheus
client v1.0.0 -> common v0.4.1 -> target. Target, both common vertices,
go-metrics, both Prometheus clients, Consul API, and Serf remain why-negative
and unloaded. Repository target imports, target production loads, complete-
test loads, module-backed loads, and runtime relevance remain zero. There is
no supported stable owner/request route.

### Repository, Release, Source, And Qualification Guards

The canonical exact-path stable line must remain empty. Public, enabled,
unarchived, non-fork Apache-2.0 repository `mwitkow/go-conntrack`, ID 72483369,
must remain owned by `mwitkow` and default to `master`, with zero tags, zero
GitHub Releases, and an empty Go proxy stable list. Exact `/v2` and `/v3`
lines remain absent. There must be no replacement, retraction, deprecation,
eligible alternate owner, later exact-path stable, or supported tidy-stable
requester route.

Every completed identity remains a guard. Selected pseudo-version commit
`cc309e4a22231782e8893f3c35ced0967807a33e`, tree
`8e653f1c675ad96001b221496975aa5858f0df09`, remains unsigned and an ancestor
of signed current master/proxy-latest pseudo-version commit
`2f068394615f73e460c2f3d2c158b0ad9321cadb`, tree
`a8074d12ad49ac88f1078e3942c1a95c8bfbc5f2`. Neither pseudo-version,
`master`, a pull-request ref, fork, replacement, or alternate path is a stable.
Both synthetic module files retain only the exact module directive, without a
Go directive, requirement, replacement, retraction, or deprecation. The
recorded license, signature, ancestry, proxy/sumdb, source/module, ZIP, exact
18-file archive-to-Git, symlink, and submodule identities remain exact.

The completed exported API and documentation; caller/package ownership and
mutation; determinism, concurrency, lifecycle, cleanup, error, and global-
state behavior; platform, network/TLS/file, cgo, build-tag, generated, embed,
and Go-floor boundaries all remain guards. With the historical Prometheus
closure, both exact SDKs must retain the recorded passing build/count-one/vet
results and repeated-test global-metric accumulation failure; Go 1.26.7 must
retain the recorded historical `x/net/trace` event-map race result. With the
project's actual Prometheus client v1.4.0/common v0.9.1 closure, both SDKs must
retain the full-upstream `prometheus.Handler` incompatibility while the two
library packages and all nine supported cross-build targets continue to
compile. No selected-source behavior, closure, upstream-gate, or independent
finding may change.

### Project, Toolchain, Earlier-Decision, And Advisory Guards

The real project remains exactly 234 selected modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` remain
74/1,067 lines at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`,
and applied tidy must remain the common 52/948-line state at
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

The selected-pseudo disposable get must retain its recorded
234-module/3,609-edge/82-line-go.mod/1,076-line-go.sum root-promotion result.
The latest-pseudo disposable get must retain its recorded
235-module/3,610-edge/83-line-go.mod/1,079-line-go.sum result and sole new
`github.com/jpillora/backoff v1.0.0` selection. Both must remain target-why-
negative, preserve exact 355/429/197/41 loads and the Go 1.18 floor, and tidy
back to the common 52/948-line projection. Neither projection is retained.

The Go 1.18 language/compatibility floor and exact Go 1.26.7 Darwin arm64
official archive/binary SHA-256 identities remain
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Every source/API/CLI/help/launcher/Make/quality contract remains unchanged,
and accepted quality remains 27/27 Q0-Q2 PASS at L2.

All 45 earlier guarded selections and their 273 sorted incoming edges remain
exact at SHA-256
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Every reflect2, concurrent, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, Cast, Viper, closed-exception, owner/request,
why/import/load/runtime, graph/module/tidy/Go-floor, advisory, qualification,
supported-owner, and compatible-route fact remains final under its own
separate expiry boundary. Including go-conntrack adds only its exact selection
guard and the two recorded incoming request edges; it broadens no earlier
decision.

Exact selected/latest target OSV, selected target GitHub global, repository
advisory, and pinned isolated govulncheck findings remain empty. Prometheus
common v0.9.1/v0.4.1 OSV results remain empty. The unchanged project retains
the exact 30/22/20/20 non-stdlib govulncheck populations with no target trace.
Guard OSV remains limited to Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go vulnerability index remains 518,501 bytes/1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z. Advisory absence is not qualification.

Any target path/version, request, requester import, owner identity or route,
root, why/import/load/runtime fact, source/behavior/closure result, graph,
module hash, tidy state, Go floor, earlier or closed-exception guard,
advisory/finding, independent defect, repository/release/owner fact,
qualification, supported-owner result, or compatible genuine route change
expires this exception and requires a fresh go-conntrack dependency evaluation
and product decision before merge. Any earlier guard change separately
requires its own fresh decision. This decision authorizes no owner study,
direct root, latest/master selection, requester change, dependency edit,
workaround, implementation, or transferred exception.

### Guard-Only Revalidation And Handoff

Guard-only revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at evaluation-handoff HEAD
`a586649d5ecaef9a9c9af5e8e33dbe631321c439`, parent
`ab5d27df6969e7ebc6766c73adf3af9c94e2cbe6`, tree
`8028ed1f41444ba77b66f6763684add385bd59fc`. That handoff changes exactly
`codex-dev-start.sh`, the answered go-conntrack evaluation archive, this
then-NEXT decision archive, rolling handover, and roadmap. Exact Google UUID
v1.4.0 dependency commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
reciprocal archive chain, sole NEXT state, launcher/archive prompt mirror,
exact changed set, and launcher check passed.

Exact Go identity, baseline counts and hashes, tidy projection and applied-
tidy identities, both target requests, every recorded route edge, negative
target/requester why results, zero repository import/load state, and the exact
unchanged graph reproduce. Because `go.mod`, `go.sum`, and the 3,599-edge graph
are byte-exact under the same toolchain, the 45-selection/273-edge earlier-
guard snapshot and all closed-exception boundaries retain their recorded
identity without drift.

Fresh primary metadata still reports repository ID 72483369 with the exact
owner/status/license/default-branch identity, zero tags/releases/repository
advisories, an empty global selected-version advisory result, and an empty
proxy stable list. Selected/latest target and both common requester OSV
results are empty. The four recorded guard advisory identities, Go index, and
memberlist CNA record reproduce exactly. The completed upstream/source/
behavior/closure/race/cross-build/projection/archive/govulncheck evaluation was
not repeated, and no owner study ran.

No product source, dependency metadata, parent, toolchain declaration, or
earlier decision changed. No changed-selection scorecard applies. Final exact-
Go module verification, build, count-one tests, race count-one tests, and vet
pass under canonical `umask 022`; accepted quality remains 27/27 Q0-Q2 PASS at
L2. Every task-owned cache, report, response, and project copy was contained
beneath the managed session scratch root and removed before handoff.

P7 continues only with the linked reciprocal bounded evaluation of the next
unanswered selected queue item, graph-selected transitive exact
`github.com/oklog/ulid v1.3.1`. That successor was prepared but not executed
and may not reopen or transfer this go-conntrack exception or any earlier
decision. P8 remains queued.
