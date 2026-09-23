# Agent Session: Evaluate Pascaldekloe Goe Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T021646+0200-evaluate-pascaldekloe-goe-dependency`
Created: `2026-09-23T02:16:46+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `72e55abb617ebbafee8d346dd44f0ba483d35c5ae50ccec0bf96ff270be77d46`
Previous: [2026-09-23T014437+0200-decide-oklog-ulid-product-direction.md](2026-09-23T014437+0200-decide-oklog-ulid-product-direction.md)
Next: [2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md](2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md)
Outcome: No exact-path stable qualifies. Selected v0.1.0 and latest v0.1.1 both fail the complete upstream race gate; v0.1.0 also fails an ordinary Go 1.26.7 example, and v0.1.1 has no genuine supported tidy-stable project owner. Source and dependency metadata remain unchanged, and one reciprocal product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating the next unanswered selected queue
item, graph-selected exact `github.com/pascaldekloe/goe v0.1.0`, as exactly one
bounded dependency group. Resolve its canonical exact-path release line and
highest qualified Go-1.18-compatible stable from primary evidence. Implement
one exact dependency-only changed selection only if the candidate, its
complete minimal closure, and every earlier target-specific guard remain
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

# Authorized Roadmap And Guards

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the eight final target-specific option-1 decisions through ULID, and qualified
no-selection-change modern-go/concurrent and modern-go/reflect2. P8 remains
queued.

Exact selected, inherited, unloaded `github.com/oklog/ulid v1.3.1` is retained
only under its own explicit, unqualified, target-specific, non-transferable
option-1 exception. V0.3.0 is the highest behavior-qualified Go-1.18-
compatible exact-path stable but has no genuine supported tidy-stable project
owner. Selected v1.3.1 and every stable from v1.0.0 through v1.3.1 have an
incomplete standalone module closure because `cmd/ulid` imports undeclared
`github.com/pborman/getopt/v2`. Preserve the exact selected version, no target
root, sole TSDB v0.7.1 request and genuine production/test imports, main ->
direct mvn-pom-mutator v0.2.3 -> metadata-only TSDB -> target route, negative
target/requester why, zero repository import/load/runtime boundary, exact
repository/release/archive/source/behavior/closure identities, both disposable
projections, and every expiry condition. Do not run the rejected Prometheus
TSDB ULID Ownership Study, select v0.3.0 or `/v2`, add a root, change TSDB or
mvn-pom-mutator, patch/vendor upstream, or transfer the exception.

Exact go-conntrack, mapstructure, go-homedir, Promptui, emoji/v2, kr/text,
kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-sequences, gotool,
errcheck, httprouter, GLS, go-junit-report, json-iterator, and clockwork
exceptions remain final, target-specific, unqualified, and non-transferable.
Canonical modern-go/reflect2 v1.0.2 remains the highest qualified exact-path
stable with a genuine supported owner and is already selected. Modern-go/
concurrent remains qualified at owner release `1.0.3`, whose exact commit and
source are selected through pseudo-version
`v0.0.0-20180306012644-bacd9c7ef1dd`. Every Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen or transfer any decision.

# Measurements At Start

The ULID decision began from clean branch `codex/upgrade-quality` at
evaluation-handoff HEAD `27046bec05003b317329a1918fd69baf5ffe4eb0`, parent
`c0819eb1524b5a5bf6aabbefebaf185d30d36e88`, tree
`15f3c216faad8c67df78a5e6eaf73ce567940653`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The decision changed no product source or dependency metadata and
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

All 46 selections earlier than ULID and their 275 incoming edges remain exact
at SHA-256
`32bd1b6893f9a7964567c1b4bf772d60d690e85509f2ec85fd299d537086e544`.
Including ULID gives 47 guarded selections and 276 incoming edges at SHA-256
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
Forty-one guarded why results are negative; only kr/pretty, kr/text, emoji/v2,
Promptui, go-homedir, and mapstructure are positive. Promptui and go-homedir
are the only guarded repository imports; emoji/v2, Promptui, go-homedir, and
mapstructure are the only production/complete-test loaded guarded modules.
Every target-specific exception, concurrent, reflect2, go-conntrack, and ULID
guard remains exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.

The queue identifies graph-selected exact
`github.com/pascaldekloe/goe v0.1.0`. Treat physical selection and that version
as starting observations only. Independently reproduce every target request,
current or historical route, requester import or metadata-only boundary,
target and requester why, repository imports, production/complete-test/module-
backed loads, runtime relevance, and root state before choosing a candidate.
Do not infer qualification, ownership, or project relevance from transitive
graph presence.

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
hashes, tidy projection, all 47 guarded selections, and the 276-edge snapshot.
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
`go get github.com/pascaldekloe/goe@<selected-version>` for one dependency-
only commit; do not hand-edit metadata and do not use tidy as the
implementation. Explain and verify the minimal exact transitive closure.

If a candidate changes any ULID, go-conntrack, reflect2, concurrent,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, or kr/pty path/version,
request, requester import, owner route, root, why/import/load/runtime fact,
graph/module/tidy/Go-floor state, advisory/release/owner identity,
qualification, or compatible route, retain no projection and stop for the
corresponding fresh decision before merge. This evaluation cannot silently
expire, replace, broaden, or transfer any closed exception or earlier
decision.

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

Read this archive, the answered ULID decision/evaluation, answered go-
conntrack decision/evaluation, reflect2 and concurrent evaluations,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/
evaluations, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Verify
branch, ancestry, clean ordinary/ignored state, reciprocal archive chain,
launcher check, exact Go identities, module hashes/counts/tidy projection, all
47 guards and the 276-edge snapshot, every closed-exception boundary, and
fresh advisory identities before any implementation.

# Three Moves

First, independently evaluate only pascaldekloe/goe and choose the exact
qualified stable or bounded no-qualified result. Second, make at most the one
authorized dependency-only selection change and verify it, or leave source and
metadata unchanged; stop before implementation if any earlier guard would
expire. Third, update roadmap and rolling handover, answer this archive,
prepare exactly one reciprocal successor required by the result, verify
containment, and make the required local handoff commit without executing the
successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen ULID, go-
conntrack, reflect2, concurrent, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, or an earlier decision, evaluate another dependency group,
write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No genuine exact-path stable `github.com/pascaldekloe/goe` release qualifies.
Selected v0.1.0 and latest v0.1.1 are Go-1.18-compatible, but both fail the
complete upstream race gate: each release's ordinary `metrics` tests read a
caller-owned `bytes.Buffer` while the permanent `NewStatsD` goroutine writes
to it, and the API exposes no close, wait, flush, or other synchronization
contract. Selected v0.1.0 additionally fails its ordinary `el.ExampleInt`
under exact Go 1.26.7 with `panic: lookup fail` after standard-library TLS
internals changed. V0.1.1 deletes that failing example but does not repair the
metrics lifecycle or race. Physical selection, transitive presence, empty
advisory results, and otherwise-green build/test gates do not qualify either
release.

No projection was retained. Product source, `go.mod`, and `go.sum` remain
byte-identical, and no dependency implementation commit exists. P7 is stopped
only for the linked reciprocal product decision: retain selected v0.1.0 under
a Goe-specific unqualified exception, authorize exactly one later
measurement-only **Armon Go-Metrics Goe Ownership Study**, or stop P7
unresolved. This evaluation chooses none of those directions. P8 remains
queued.

### Canonical Identity And Release Line

The exact go-import metadata maps `github.com/pascaldekloe/goe` to
`https://github.com/pascaldekloe/goe.git`. Repository ID 35051646 is public,
enabled, unarchived, non-fork, CC0-1.0 licensed, owned by `pascaldekloe`, and
defaults to `master`. It has no GitHub Release objects. The repository remains
maintained, including a merged 2026 vet-format repair after the latest tag.

The canonical proxy and Git repository expose exactly two genuine exact-path
stable tags, v0.1.0 and v0.1.1. `/v2` and `/v3` are absent. There is no
prerelease, replacement, retraction, deprecation, redirect, alternate module
path, or eligible fork/branch/pseudo-version stable.

- V0.1.0 is lightweight unsigned tag/commit
  `57f6aae5913c64c9bcae5dbdffd33365b5a7f138`, parent
  `9dc2dab6f887587f166043c220ce33ff015e8d52`, tree
  `752c419e5ca27e3893074365e3fa87fc22811f3f`, timestamp
  2018-06-27T14:32:12Z. Git has no module file; the proxy supplies a synthetic
  exact-path module file with no Go directive or requirements. Source/module
  sums are `h1:cBOtyMzM9HTpWjXfbbunk26uA6nG3a8n06Wieeh0MwY=` /
  `h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=`. The proxy ZIP has SHA-256
  `37b73886f1eec9b093143e7b03f547b90ab55d8d5c9aa3966e90f9df2d07353c`;
  its normalized 17-file manifest is
  `909528a40180eb5283b8b08e63310fda9fdb439235534d5f9559d4594fe4530a`.
- V0.1.1 is lightweight unsigned tag/commit
  `ad498856776794021b60e2b5ca6752ac7be7e588`, parent
  `10514f506b0c6985b96f8f4d4577d7648a8b4642`, tree
  `e198f3136240ed195fe4b1981dced867e7a90ab6`, timestamp
  2023-03-26T17:29:27Z. Its exact-path module declares Go 1.16 and no
  requirements. Source/module sums are
  `h1:Ah6WQ56rZONR3RW3qWa2NCZ6JAVvSpUcoLBaOmYFt9Q=` /
  `h1:KSyfaxQOh0HZPjDP1FL/kFtbqYqrALJTaMafFUIccqU=`. The proxy ZIP has SHA-256
  `250bf987eb4c8aac6570c278c0a4a1db76e2e4e14f6b01e234d9b4d3ed2e328a`;
  its normalized 18-file manifest is
  `e8e31bfe4ee257b5c65e7a0034d9debeb7363963e5e02fec73c3bf3ae6e4c2fe`.

Both proxy archives contain only byte-identical regular files from their
exact Git trees, with no symlink or submodule. SumDB records match. V0.1.0 is
an ancestor of v0.1.1, and both are ancestors of signed current `master`
commit `460052cb77c265139461ae461f395abd001b47e0`, tree
`f05b17c821c648dee0eae052ccd03f67d62308b1`. The only unreleased source delta
after v0.1.1 is six vet-format repairs; it was not promoted as a stable.

### Source, Closure, And Upstream Qualification

Both releases expose the standard-library-only packages `el`, `metrics`,
`rest`, and `verify`. They have no cgo, build tags, generated directives, or
embed boundary. The production boundary includes reflection and mutable
caller data in `el`; global registration and mutable `StatsDPackMax` in
`metrics`; JSON and `net/http` callbacks in `rest`; and testing/reporting in
`verify`. `NewStatsD` owns a permanent goroutine, channel, and buffer pool,
writes to a caller-provided writer, discards write errors, and exposes no
lifecycle or cleanup operation. Prefix and caller-owned state are not made
concurrency-safe by the package.

Under exact Go 1.26.7 the complete closure contains 194 production packages,
228 complete-test packages, and 14 module-backed entries across one module.
Under exact Go 1.18.10 it contains 132/164/14 across one module. Both releases
pass module verification, build, and vet under both SDKs, and pass cgo-disabled
build plus test compilation for Darwin amd64/arm64, Linux amd64/arm64/386,
Windows amd64/386, FreeBSD amd64, and js/wasm.

V0.1.0 passes count-one and count-ten upstream tests under Go 1.18.10 but
fails both under Go 1.26.7 first at `el.ExampleInt`. V0.1.1 passes count-one
and count-ten tests under both SDKs. Both releases fail `go test -race ./...
-count=1` under both SDKs in the ordinary upstream `metrics` tests, at the
caller buffer read versus the `NewStatsD` worker write. A clean Go 1.18.10
rerun with the direct Xcode compiler and SDK reproduces the same target race.
No adversarial fixture, security reproduction, fuzzing, stress, or resource
probe was used. Therefore neither stable passes the complete ordinary
qualification contract.

### Exact Requests, Routes, And Project Boundary

MVS selects v0.1.0 from four requests:

1. `github.com/armon/go-metrics@v0.3.10 -> github.com/pascaldekloe/goe@v0.1.0`
2. `github.com/hashicorp/consul/api@v1.1.0 -> github.com/pascaldekloe/goe@v0.0.0-20180627143212-57f6aae5913c`
3. `github.com/hashicorp/memberlist@v0.1.3 -> github.com/pascaldekloe/goe@v0.0.0-20180627143212-57f6aae5913c`
4. `github.com/hashicorp/memberlist@v0.3.0 -> github.com/pascaldekloe/goe@v0.0.0-20180627143212-57f6aae5913c`

Go-metrics v0.3.10 genuinely imports `goe/verify` in tests. Consul API v1.1.0
genuinely imports it in `txn_test.go` and `session_test.go`. Both memberlist
requests are metadata-only. The current and historical routes include main ->
direct mvn-pom-mutator v0.2.3 -> historical Viper v1.10.1 -> go-metrics
v0.3.10 -> target, plus Viper -> sagikazarmark/crypt v0.4.0 -> go-metrics;
and main -> mvn-pom-mutator -> bketelsen/crypt -> Consul API v1.1.0 -> target,
with Consul/Serf and Viper/crypt/Serf routes through memberlist v0.1.3 or
v0.3.0.

Target, go-metrics, Consul API, and both memberlist requester why results are
negative. The repository has no Goe import. Its 355 production entries, 429
complete-test entries, and 197 module-backed entries across 41 loaded modules
contain no target package. Goe has no main root and no project runtime
relevance. The genuine requester imports exist only in unloaded historical
dependency tests; physical graph presence is not supported project ownership.

### Disposable Projections And Earlier Guards

An exact disposable v0.1.0 get is output-empty and adds only one unused root,
one source sum, and one main graph edge. It produces 75/1,068-line module
files with SHA-256
`8f93778370f9deca4a62e6b6e57279c79d35db7367ea28565f290b71ce69394a` /
`23ac827f15e13cb3574f61f704224530ea74e1343408029b5dcfed51f6ae4ee9`,
234 modules, 3,600 edges, and unchanged 355/429/197/41 loads. Its 436-line
tidy diff has SHA-256
`db6145b4d0f454da79bc8ad293b8f79849ed4d97f5e2dc03434b235b5660ce12`;
applied tidy restores the common 52/948-line state and inherited v0.1.0.

An exact disposable v0.1.1 get adds the same unused root/edge and two sums,
changing only the selected target version. It produces 75/1,069-line module
files with SHA-256
`ffab908e260747604813fa491447741a701fb44c2c5fda382261c7dd1ef15b92` /
`38d75c9287fccd450406c1d675fc4c13c089baf90d6723a6bc0a1faea134a41d`,
234 modules, 3,600 edges, and unchanged loads. Its 436-line tidy diff has
SHA-256
`74c7b07ae0773e72c877b2d24c5325c3630ad60ea74462657470a1372fcbd531`;
applied tidy again restores the common 52/948-line state and reselects
v0.1.0. Thus v0.1.1 has no genuine supported tidy-stable project owner, in
addition to failing the race gate. Neither projection was retained.

The real project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. Real 74/1,067-line `go.mod` /
`go.sum` SHA-256 remains
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection and common applied 52/948-line state retain
SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708` /
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 earlier guarded selections and their 276 sorted incoming edges remain
exact at SHA-256
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`
in the base and both projections. Forty-one guarded why results remain
negative; only kr/pretty, kr/text, emoji/v2, Promptui, go-homedir, and
mapstructure remain positive. Promptui and go-homedir remain the only guarded
repository imports; emoji/v2, Promptui, go-homedir, and mapstructure remain
the only loaded guarded modules. Every earlier exception, qualified decision,
owner/request route, projection, advisory, and expiry boundary remains exact.

### Advisory, Continuity, And Final Gate

Fresh exact-version OSV and GitHub global results are empty for v0.1.0 and
v0.1.1, and the repository advisory endpoint is empty. Pinned isolated
govulncheck v1.8.0 module/package/symbol/test-symbol scans are empty for both
releases. Exact-Go base and v0.1.1 project populations are identical at
30 module, 22 vulnerable-package, 20 called-symbol, and 20 test-symbol OSVs,
with no Goe trace.

Guard OSV remains limited to Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go module vulnerability index remains 518,501 bytes and
1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z. Advisory absence is not qualification.

Evaluation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at handoff HEAD
`536df39770daedc4f841a6898e5c69b8605eed6f`, parent
`27046bec05003b317329a1918fd69baf5ffe4eb0`, tree
`bf97e1053d7f50c7df6e55e1800ed6d751d77915`. That handoff changes exactly
the launcher, answered ULID decision archive, this then-NEXT evaluation
archive, rolling handover, and roadmap. Google UUID implementation commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.

The official Darwin arm64 SDK identities used were Go 1.26.7 archive/binary
SHA-256
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
and Go 1.18.10 archive/binary SHA-256
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Final unchanged-project exact Go 1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under `umask 022`. The Go 1.18,
source/API/CLI/help/launcher/Make/quality contracts and accepted 27/27 Q0-Q2
PASS at L2 remain unchanged.

Every task-owned SDK, cache, archive, report, repository clone, project copy,
tool, and advisory response remained beneath the managed session scratch root
and was removed before handoff; only the pre-existing launcher-owned Node
compile cache remains. The reciprocal successor was prepared but not
executed. No successor was launched, no other dependency was evaluated, and
P8 was not begun.
