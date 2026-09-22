# Agent Session: Evaluate Mwitkow Go-Conntrack Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-22T231926+0200-evaluate-mwitkow-go-conntrack-dependency`
Created: `2026-09-22T23:19:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d9cf3e35fa26406305f6287ab332a3dd71bec645dd77b4a7d549bc1d16a61b1a`
Previous: [2026-09-22T222053+0200-evaluate-modern-go-reflect2-dependency.md](2026-09-22T222053+0200-evaluate-modern-go-reflect2-dependency.md)
Next: [2026-09-23T001118+0200-decide-mwitkow-go-conntrack-product-direction.md](2026-09-23T001118+0200-decide-mwitkow-go-conntrack-product-direction.md)
Outcome: no canonical exact-path stable exists or qualifies; the selected unloaded pseudo-version and all dependency metadata remain unchanged pending a target-specific product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating the next unanswered selected queue
item, graph-selected transitive exact
`github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223`, as
exactly one bounded dependency group. Resolve its canonical exact-path release
line and highest qualified Go-1.18-compatible stable from primary evidence.
Implement one exact dependency-only changed selection only if the candidate,
its complete minimal closure, and every earlier target-specific guard remain
exact. Do not combine another dependency group or begin P8.

# Defensive Scope

This is an ordinary dependency-quality evaluation. Use public metadata,
static source/repository records, project graph/build commands, upstream tests,
and only small bounded ordinary fixtures required by documented behavior. Do
not fuzz, stress, probe resource exhaustion, create oversized, deeply nested,
cyclic, malformed, adversarial, or escape-sequence payloads, reproduce a
security issue, or perform security or exploitability analysis.

Every disposable cache, tool, archive, report, project copy, fixture, or
advisory response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`.
Never write to `/private/tmp`, `/tmp`, a sibling of the managed root, or
another external root. Verify containment and remove task-owned scratch
evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the six final target-specific option-1 decisions through mapstructure v1.5.0,
qualified no-selection-change modern-go/concurrent, and qualified no-
selection-change modern-go/reflect2. P8 remains queued.

Canonical `github.com/modern-go/reflect2 v1.0.2` is the highest qualified
Go-1.18-compatible exact-path stable with genuine supported unarchived owner
`modern-go/reflect2`. It is already selected, so no selection or closure
changed. A disposable exact get only adds an otherwise absent main indirect
root, source sum, and one graph edge; tidy removes that promotion, so it was
not retained. Preserve its exact four-release line, owner/release/commit/tree/
signature/archive/source/module identities, Go-1.18 behavior fix, unsafe and
platform boundaries, nine request and requester-import boundaries, negative
why/import/load/runtime facts, projection, graph, advisory, qualification, and
supported-owner facts as guards.

Modern-go/concurrent remains qualified at owner release `1.0.3`, whose exact
commit/source is already selected through pseudo-version
`v0.0.0-20180306012644-bacd9c7ef1dd`. Mapstructure, go-homedir, Promptui,
emoji/v2, kr/text, and kr/pty option-1 decisions separately remain exact,
unqualified, target-specific, non-transferable, and final under every recorded
expiry guard. Do not transfer an exception or reopen reflect2, concurrent,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast,
Viper, or an earlier decision.

# Measurements At Start

The modern-go/reflect2 evaluation began from clean branch
`codex/upgrade-quality` at handoff HEAD
`1c047f8ded007382fd26df384a0f998ba15156f0`, parent
`7a4222d5d3d9e63931ac7d0f270904cfbc509f58`, tree
`74fb7c88eed37b863708a76d0b3a1be6d44d9187`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The evaluation changed no product source or dependency metadata and
prepared this evaluation-only handoff. Verify the new handoff HEAD, parent,
tree, exact changed set, ancestry, reciprocal archive chain, and clean
ordinary and ignored status rather than assuming them.

The unchanged project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The unchanged 432-line tidy projection remains SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

The 44 selections earlier than reflect2 remain exact with 264 sorted incoming
edges at SHA-256
`4a83d4e3a015f93da4b4d35bd590137071b01c40315ec9bb9c987688004f9b09`.
Including qualified closed reflect2 gives 45 guarded selections and 273
incoming edges at SHA-256
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Thirty-nine why results are negative; only closed kr/pretty, kr/text,
emoji/v2, Promptui, go-homedir, and mapstructure are positive. Promptui and
go-homedir are the only guarded repository imports; emoji/v2, Promptui,
go-homedir, and mapstructure are the only production/complete-test loaded
guarded modules. Every closed-exception, concurrent, and reflect2 guard
remains exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.

The queue selects exact
`github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223`
without a main `go.mod` request. Two graph edges request the target: selected
`github.com/prometheus/common v0.9.1` and historical Prometheus common v0.4.1
both request that exact pseudo-version. The target has only its module checksum
in current `go.sum`. Current target why is negative, current repository source
has no target import, and no target package is production or complete-test
loaded. Treat physical selection, transitive graph presence, negative why,
non-loading, pseudo-version status, or apparent release status as observations
rather than qualification. Independently reproduce both requests, every
current and historical route, requester import or metadata-only boundary, and
why/import/load/runtime fact before choosing a candidate.

# Evaluation Contract

Resolve exact go-import and module-path identity; repository owner/status/
license/default branch; version/tag/release, commit/tree/signature/ancestry,
proxy/sumdb/archive-to-Git identity; module directives and requirements;
retractions, deprecation, replacements, and exact-path major lines. Consider
only genuine exact-path stable releases. Do not promote a fork, branch,
pseudo-version, prerelease, replacement, alternate module path, or ownerless
candidate as a stable.

For selected and every serious stable candidate, inspect the complete module
and test closure, exported API and documentation, Go-floor compatibility,
platform/build-tag/cgo/generated/embed boundaries, globals, ownership and
mutation, determinism, concurrency, lifecycle, cleanup, and error behavior.
Exercise only small bounded ordinary values needed to verify documented
behavior. Run upstream build, tests, repeated tests, race, vet, and supported
cross-builds under exact Go 1.26.7 and a contained Go 1.18 toolchain. A release
qualifies only if all applicable ordinary documented contracts and every
project guard pass.

Map every target MVS request and genuine current or historical route.
Reproduce target and requester why, repository imports, production and
complete-test loads, module-backed entries, runtime relevance, graph counts,
hashes, tidy projection, all 45 guarded selections, and the 273-edge snapshot.
Physical selection, transitive graph presence, a why result, loading, or
advisory absence is not qualification.

Use disposable project copies beneath the managed scratch root to measure
exact candidate projections. Never add or alter a target root in the real
project outside the one final exact dependency implementation authorized
below. For each projection record exact selection, closure, graph, imports/
loads, sums, tidy result, Go-floor effect, genuine supported tidy-stable
ownership, and whether every earlier guard remains exact. Do not retain a
projection unless the candidate qualifies and the normal dependency
implementation contract authorizes it.

Refresh exact-version OSV and GitHub advisory evidence, repository advisories,
the guarded advisory population, x/mod guard, Go vulnerability index identity,
memberlist CNA identity, and a pinned isolated govulncheck comparison.
Advisory absence cannot override ordinary behavior, an upstream gate,
ownership, or an earlier-guard failure. Stay within the defensive scope.

# Decision And Implementation Boundary

Select only the highest qualified Go-1.18-compatible exact-path stable with a
genuine supported project owner. If that exact selection changes and every
earlier guard remains exact, use exact Go 1.26.7 and exact
`go get github.com/mwitkow/go-conntrack@<selected-version>` for one
dependency-only commit; do not hand-edit metadata and do not use tidy as the
implementation. Explain and verify the minimal exact transitive closure.

If a candidate changes any reflect2, concurrent, mapstructure, go-homedir,
Promptui, emoji/v2, kr/text, or kr/pty path/version, request, requester import,
owner route, root, why/import/load/runtime fact, graph/module/tidy/Go-floor
state, advisory/release/owner identity, qualification, or compatible route,
retain no projection and stop for the corresponding fresh decision before
merge. This evaluation cannot silently expire, replace, broaden, or transfer
any closed exception or earlier decision.

If the current selection is already the highest qualified exact decision and
exact get is a byte no-op, record the no-change decision without forcing a
dependency commit. If no stable qualifies, or no qualified candidate has a
genuine supported tidy-stable owner, retain no projection and prepare one
reciprocal product-decision archive offering only a target-specific
unqualified exception, exactly one named later measurement-only genuine
owner/request study, or stopping P7 unresolved. Do not choose that product
direction during the evaluation.

After any changed selection, run the complete P7 dependency gate required by
the roadmap. For an unchanged result, run focused target/graph/advisory guards
and final exact-Go module verification, build, count-one tests, race count-one
tests, and vet. Preserve Go 1.18, source/API/CLI/help/launcher/Make/quality
contracts, every earlier decision, and accepted 27/27 Q0-Q2 PASS at L2.

# Required Reading

Read this archive, the answered reflect2 and concurrent evaluations, answered
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/
evaluations, kr/pretty and Cast owner records, rolling handover, P7/P8 roadmap,
`go.mod`, and `go.sum`. Verify branch, ancestry, clean ordinary/ignored state,
reciprocal archive chain, launcher check, exact Go identities, module hashes/
counts/tidy projection, all 45 guards and the 273-edge snapshot, every closed-
exception boundary, and fresh advisory identities before any implementation.

# Three Moves

First, independently evaluate only mwitkow/go-conntrack and choose the exact
qualified stable or bounded no-qualified result. Second, make at most the one
authorized dependency-only selection change and verify it, or leave source
and metadata unchanged; stop before implementation if any earlier guard would
expire. Third, update roadmap and rolling handover, answer this archive,
prepare exactly one reciprocal successor required by the result, verify
containment, and make the required local handoff commit without executing the
successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen reflect2,
concurrent, mapstructure, go-homedir, Promptui, emoji/v2, kr/text, kr/pty,
kr/pretty, Cast, Viper, or an earlier decision, evaluate another dependency
group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No canonical exact-path stable exists for `github.com/mwitkow/go-conntrack`,
so there is no qualified stable to select. The graph-selected
`v0.0.0-20161129095857-cc309e4a2223` remains an unqualified pseudo-version.
Product source, `go.mod`, and `go.sum` remain unchanged, no projection was
retained, and no dependency implementation commit exists.

### Canonical Identity And Empty Stable Line

Exact go-import metadata maps `github.com/mwitkow/go-conntrack` to
`https://github.com/mwitkow/go-conntrack.git`. The public enabled, unarchived,
non-fork Apache-2.0 repository, ID 72483369, remains owned by `mwitkow` and
defaults to `master`. It has zero tags and zero GitHub Releases, and the Go
proxy version list is empty. Thus the canonical exact-path stable line is
empty: neither a branch nor either observed pseudo-version can be promoted as
a stable release.

The selected pseudo-version is commit
`cc309e4a22231782e8893f3c35ced0967807a33e`, tree
`8e653f1c675ad96001b221496975aa5858f0df09`, and is unsigned. Current master
and proxy `@latest` are the selected commit's descendant
`2f068394615f73e460c2f3d2c158b0ad9321cadb`, tree
`a8074d12ad49ac88f1078e3942c1a95c8bfbc5f2`, exposed only as
`v0.0.0-20190716064945-2f068394615f`; GitHub verifies that merge commit's
signature. Later 2020/2021 objects are confined to non-default pull-request or
feature refs and are not releases.

Both proxy module files are synthetic single-directive files for the exact
module path, with no Go directive, requirements, replacement, retraction, or
deprecation. No `/v2` or `/v3` line or eligible alternate exact-path owner
exists. Selected source/module sums are
`h1:F9x/1yl3T2AeKLr2AMdilSD8+f9bvMnNN8VS5iDtovc=` and
`h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=`; latest source uses
`h1:KUppIJq7/+SVif2QVs3tOP0zanoHgBEVAwHxUSIzRqU=` with the same module sum.
The selected proxy ZIP SHA-256 is
`9b3c27b462c9526d7851db3655de380727abec8d62a34076960eaac16a0caf90`.
Its 18 regular files byte-match the exact Git tree, whose normalized manifest
hash is `dcd17777cd72c25ac71b7c34641e95b716e034cd357ff2aa2ab930475929608a`;
there are no symlinks or submodules. The Apache-2.0 license hash is
`b40930bb246614c2e022741116d891b3e286412dfe40fc2011e75509306333e1`.

### Source And Ordinary Behavior

The selected source exports TLS certificate helpers, Prometheus dialer
preregistration, dial/listener option functions, context name helpers, dialer
constructors, and `NewListener`. It wraps caller-owned connections/listeners,
serializes trace finishing, closes the underlying connection on `Close`, and
updates package-global Prometheus collectors registered from `init`. Callers
retain their contexts, TLS configuration, option inputs, and wrapped network
resources. Map/trace/metric side effects are global; listener TCP keepalive
errors are ignored. There is no cgo, generated source, embed input, or build
tag. Network/TLS behavior and the certificate-file helpers are its ordinary
documented external boundaries; it retains no package-owned background
goroutine.

With its historical Prometheus common v0.4.1/client v0.9.1 closure, the full
upstream build, count-one tests, and vet pass under exact Go 1.18.10 and Go
1.26.7; Go 1.18.10 race count one also passes. Repeated tests fail under both
SDKs because the package-global connection-refused counter reaches ten while
the test still requires one. Go 1.26.7 race count one additionally reports a
concurrent read/write in the historical `x/net/trace` event-family map between
trace rendering and the target's listener tracker. With the project's actual
selected Prometheus client v1.4.0/common v0.9.1 closure, full upstream build,
tests, race, and vet fail under both SDKs because the upstream example and
test helper still call removed `prometheus.Handler`; the target and
`connhelpers` library packages themselves build. Nine bounded cgo-disabled
Darwin, Linux, Windows, FreeBSD, and js/wasm cross-builds pass under both
SDKs. These ordinary failures independently leave the selected pseudo-version
unqualified, although the empty stable line is already decisive.

### Requests, Routes, And Project Relevance

Exactly two graph edges request the target, both at the selected pseudo-
version: selected `github.com/prometheus/common v0.9.1` and historical common
v0.4.1. Both requester sources genuinely import the package in
`config/http_config.go` and construct a traced, named dial context function;
neither edge is metadata-only.

The current route is main -> direct `mvn-pom-mutator v0.2.3` -> historical
Viper v1.10.1 -> `armon/go-metrics v0.3.10` -> common v0.9.1 -> target;
go-metrics also reaches common v0.9.1 through Prometheus client v1.4.0. The
historical route continues common v0.9.1 -> Prometheus client v1.0.0 -> common
v0.4.1 -> target. Both common vertices, go-metrics, both Prometheus clients,
Consul API, and Serf are why-negative and unloaded. Target why is negative,
the repository has zero target imports, and no target package occurs in the
355 production, 429 complete-test, or 197 module-backed entries across 41
loaded modules. There is no current runtime route.

### Disposable Projections And Earlier Guards

The real project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

A disposable exact selected-version get keeps 234 selections but promotes the
target and seven already-selected package dependencies to main indirect roots,
adds their needed source sums, and exposes two pre-existing protobuf edges.
It produces 82/1,076 metadata lines and 3,609 graph edges. A disposable latest-
pseudo get instead produces 235 modules, 83/1,079 metadata lines, and 3,610
edges because it additionally selects and roots `github.com/jpillora/backoff
v1.0.0`. Both retain exact 355/429/197/41 loads, negative target why, zero
target load, and Go 1.18. Tidy removes every manufactured root and latest-only
selection and restores the exact common 52/948-line projection. Neither
projection is a qualified stable selection, so neither was retained.

All 45 earlier guarded selections and their 273 sorted incoming edges remain
byte-exact at
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
The 39 negative and six positive why results, two guarded repository imports,
four guarded loaded modules, every closed exception, both modern-go
qualifications, every owner/request route, and every expiry boundary remain
unchanged.

### Advisories, Verification, And Handoff

Selected and latest exact-version OSV and GitHub global results are empty;
the repository advisory endpoint is also empty. Pinned govulncheck v1.8.0,
built with exact Go 1.26.7, reports no target module/package/symbol/test-symbol
finding. The unchanged project retains identical 30 module, 22 imported-
package, 20 reachable-symbol, and 20 test-symbol non-stdlib OSV populations;
the target remains unloaded. Guard OSV retains only Gorilla WebSocket
`GO-2026-6278` / `GHSA-w67g-5rqw-f597` and go-retryablehttp
`GO-2024-2947` / `GHSA-v6v8-xj6m-xwqh`; the x/mod v0.14.0 guard retains
`GO-2026-6179` and `GO-2026-6180`.

The Go vulnerability module index is 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z. Advisory absence was not used as
qualification.

Official Go 1.26.7 archive/binary SHA-256 identities are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
Go 1.18.10 identities are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Final required-umask exact-Go module verification, build, count-one tests,
race count-one tests, and vet pass. Accepted quality remains 27/27 Q0-Q2 PASS
at L2.

P7 remains unresolved only for the linked target-specific product direction:
retain the exact selected pseudo-version under an explicit unqualified
exception, authorize the single named later Prometheus Common Go-Conntrack
Ownership Study, or stop P7 unresolved. That successor was prepared but not
executed. P8 remains queued.
