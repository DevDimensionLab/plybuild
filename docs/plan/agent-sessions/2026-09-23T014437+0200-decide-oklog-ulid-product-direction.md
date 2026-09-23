# Agent Session: Decide Oklog ULID Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T014437+0200-decide-oklog-ulid-product-direction`
Created: `2026-09-23T01:44:37+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f2d0113d352624faadaab7125608511e63ca40beff7cc932bb4ffbe8bde1df53`
Previous: [2026-09-23T003641+0200-evaluate-oklog-ulid-dependency.md](2026-09-23T003641+0200-evaluate-oklog-ulid-dependency.md)
Next: [2026-09-23T021646+0200-evaluate-pascaldekloe-goe-dependency.md](2026-09-23T021646+0200-evaluate-pascaldekloe-goe-dependency.md)
Outcome: Option 1 selected. Exact selected, inherited, unloaded ULID v1.3.1 is retained under a target-specific, unqualified, non-transferable exception; source and dependency metadata remain unchanged.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by making exactly one product decision for selected, inherited,
unloaded `github.com/oklog/ulid v1.3.1` from the completed evaluation. Choose
only one direction offered below, preserve every earlier target-specific guard,
update the roadmap and rolling handover, answer this archive, prepare exactly
one reciprocal successor required by the chosen direction, and make the local
handoff commit. Do not reevaluate ULID, evaluate another dependency group,
execute the named study, or begin P8.

# Defensive Scope

This is an ordinary product decision and guard-only revalidation. Use public
metadata, static repository records, and ordinary project graph/build commands.
Do not fuzz, stress, probe resource exhaustion, create oversized, deeply
nested, cyclic, malformed, adversarial, or escape-sequence payloads, reproduce
a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, response, report, or project copy must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Verify
containment and remove task-owned scratch evidence before handoff.

# Completed Evaluation Is Final

The exact `github.com/oklog/ulid` path has eight genuine stable releases:
v0.1.0, v0.2.0, v0.3.0, v1.0.0, v1.1.0, v1.2.0, v1.3.0, and v1.3.1.
V0.3.0 is the highest behavior-qualified Go-1.18-compatible exact-path stable,
but no genuine supported tidy-stable project owner requests it. The sole
genuine requester, Prometheus TSDB v0.7.1, requests selected v1.3.1. Every
release from v1.0.0 through v1.3.1 has an incomplete standalone module closure:
`cmd/ulid` imports undeclared `github.com/pborman/getopt/v2`, and complete
build/test/vet fail under exact Go 1.18.10 and Go 1.26.7. Therefore no stable
satisfies both complete ordinary closure qualification and supported ownership.

The separate `/v2` module through v2.1.2 is not an exact-path candidate, and
`/v3` is absent. Do not promote the v2 path, a branch, fork, replacement,
pseudo-version, prerelease, or alternate owner as the v1 stable decision. Do
not patch upstream, add the missing CLI requirement, vendor code, add a target
root, or change TSDB or mvn-pom-mutator in this decision.

The exact current route is main -> direct mvn-pom-mutator v0.2.3 -> metadata-
only Prometheus TSDB v0.7.1 -> exact ULID v1.3.1. TSDB genuinely imports ULID
in production and tests. Target and TSDB why are negative; repository imports
and target production/complete-test/module-backed loads are zero. The selected
module is not runtime-relevant to the project.

A disposable selected get only adds an unused main root, source sum, and one
graph edge; tidy restores the common projection. A disposable v0.3.0 get
removes direct mvn-pom-mutator, TSDB, go-conntrack, and 24 guarded selections,
yields 157 modules/2,219 edges, and makes project loading fail. Tidy restores
the exact common projection and reselects v1.3.1. Neither projection was
retained. Do not repeat the release, source, API, behavior, closure, race,
cross-build, archive, projection, or govulncheck evaluation unless a guard-only
check proves an input changed; stop for a fresh ULID evaluation if it did.

# Authorized Roadmap And Guards

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the seven final target-specific option-1 decisions through go-conntrack, and
qualified no-selection-change modern-go/concurrent and modern-go/reflect2.
P8 remains queued.

Exact go-conntrack, mapstructure, go-homedir, Promptui, emoji/v2, kr/text,
kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-sequences, gotool,
errcheck, httprouter, GLS, go-junit-report, json-iterator, and clockwork
exceptions are final, target-specific, unqualified, and non-transferable.
Every concurrent, reflect2, Cast, Viper, memberlist, earlier selection, owner,
route, why/import/load/runtime, graph/module/tidy/Go-floor, advisory, source,
behavior, closure, qualification, and expiry guard remains final. Do not reopen
or transfer any decision.

The unchanged project is 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 45 earlier guards and 273 sorted incoming edges have SHA-256
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Including go-conntrack gives 46 guards/275 edges at SHA-256
`32bd1b6893f9a7964567c1b4bf772d60d690e85509f2ec85fd299d537086e544`.
Forty why results are negative; only kr/pretty, kr/text, emoji/v2, Promptui,
go-homedir, and mapstructure are positive. Promptui and go-homedir are the only
guarded repository imports. Emoji/v2, Promptui, go-homedir, and mapstructure
are the only production/complete-test loaded guarded modules.

Exact target OSV, GitHub global, repository advisory, and isolated pinned
govulncheck findings are empty. Project govulncheck populations are exactly
30/22/20/20 with no target trace. Guard OSV remains limited to the recorded
Gorilla WebSocket and go-retryablehttp pairs; x/mod v0.14.0 retains
GO-2026-6179 and GO-2026-6180. The Go vulnerability index is 518,501 bytes/
1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response is 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

# Measurements At Start

The ULID evaluation began from clean branch `codex/upgrade-quality` at
decision-handoff HEAD `c0819eb1524b5a5bf6aabbefebaf185d30d36e88`, parent
`a586649d5ecaef9a9c9af5e8e33dbe631321c439`, tree
`e2f0a3da140de4e35a1030f97510800ba1ddf990`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
evaluation changed no product source or dependency metadata. Independently
verify the new handoff HEAD, parent, tree, exact changed set, ancestry, clean
ordinary and ignored status, reciprocal archive chain, launcher check, exact
Go identities, module hashes/counts/tidy projection, route/why/import/load
facts, earlier guards, advisories, and final gates rather than assuming them.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, unloaded
   `github.com/oklog/ulid v1.3.1` without product-source or dependency-metadata
   changes under a ULID-specific, unqualified, non-transferable exception.
   Record that v0.3.0 is the highest behavior-qualified stable but has no
   genuine supported tidy-stable project owner, while selected v1.3.1 and all
   later v1 stables fail complete module closure. Guard the exact target
   version, sole TSDB request and genuine imports, complete route, no main
   root, negative why/import/load/runtime facts, repository/release/archive/
   source/behavior/closure identities, projections, all earlier decisions,
   and no new advisory or independent defect. If selected, prepare but do not
   execute the next bounded P7 evaluation of graph-selected exact
   `github.com/pascaldekloe/goe v0.1.0`.
2. Authorize exactly one later, measurement-only **Prometheus TSDB ULID
   Ownership Study**. It may determine whether a genuine supported tidy-stable
   TSDB owner/request route can remove ULID or request a qualified Go-1.18-
   compatible exact-path stable while preserving every earlier guard. It may
   not change product source, dependency metadata, target roots, TSDB,
   mvn-pom-mutator, another requester, or an earlier decision. If selected,
   prepare that named study as the sole successor but do not execute it.
3. Stop P7 unresolved, leave source/dependency metadata and all decisions
   unchanged, prepare no dependency evaluation or study, keep P8 queued, and
   mark the launcher COMPLETE only for the authorized roadmap.

Do not invent a fourth option, combine options, silently choose v0.3.0,
describe v1.3.1 as qualified, stable-supported, safe, or fixed, transfer an
earlier exception, or use advisory absence as acceptance. This decision
authorizes no dependency implementation commit.

# Required Reading And Revalidation

Read this archive, the answered ULID evaluation, answered go-conntrack
decision/evaluation, reflect2 and concurrent evaluations, mapstructure,
go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/evaluations,
rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Earlier evaluation
facts are final. Refresh only the narrow decision inputs: selection and sole
request, route and requester import boundary, target/root/why/import/load/
runtime facts, project hashes/counts/tidy state, earlier selection/edge
snapshot, repository release status, exact target and guarded advisories, Go
index/memberlist identities, and final exact-Go module verification, build,
count-one tests, race count-one tests, and vet. Stop for a fresh owning
evaluation if a guard changed.

# Three Moves

First, guard-only revalidate the completed result and choose exactly one
offered direction. Second, record only that direction without dependency or
product implementation. Third, update roadmap and rolling handover, answer
this archive, prepare exactly one reciprocal successor required by the choice
or mark the authorized roadmap COMPLETE for option 3, verify scratch
containment, and make the local handoff commit without executing a successor.

# Automatic Handoff

Do not launch a successor, run the study, push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change source or dependency
metadata, add a root, change a requester, transfer an exception, reevaluate
ULID or an earlier group, evaluate pascaldekloe/goe or another dependency,
write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, unloaded
`github.com/oklog/ulid v1.3.1` is explicitly retained without changing
product source, `go.mod`, or `go.sum` under a ULID-specific, unqualified,
non-transferable exception. It is not described as qualified, stable-
supported, safe, or fixed. V0.3.0 is the highest behavior-qualified Go-1.18-
compatible exact-path stable, but no genuine supported tidy-stable project
owner requests it. Selected v1.3.1 and every stable from v1.0.0 through
v1.3.1 fail complete standalone module closure. Physical MVS selection,
transitive presence, negative why, zero loading, and advisory absence do not
qualify or implicitly accept v1.3.1.

The exception accepts only the completed repository/release/module/source,
API and ordinary behavior, closure/upstream/race/vet/cross-build, owner/
request/route, why/import/load/runtime, projection/graph/tidy/Go-floor,
earlier-guard/closed-exception, and advisory findings recorded by the answered
evaluation. It accepts no uncharacterized behavior, new advisory, independent
defect, future route, direct root, or loaded use. It transfers or broadens no
earlier exception. Option 2's **Prometheus TSDB ULID Ownership Study** was
neither authorized nor run, and option 3 was not selected.

### Exact Retention, Request, And Route Guards

Retention requires exact selected path/version
`github.com/oklog/ulid v1.3.1`, with no direct or indirect main-module target
root. The sole graph request must remain
`github.com/prometheus/tsdb v0.7.1 -> github.com/oklog/ulid v1.3.1`.
TSDB must continue genuinely importing the target in production `block.go`,
`compact.go`, and `db.go` and test `db_test.go`; the request may not become
metadata-only or select a different version.

The complete route remains main -> direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` -> metadata-only
Prometheus TSDB v0.7.1 -> target. Main must continue genuinely importing
mvn-pom-mutator while mvn-pom-mutator's TSDB edge remains metadata-only.
Target and TSDB why results remain negative. Repository target imports,
target production loads, complete-test loads, module-backed loads, and runtime
relevance remain zero. No second requester, historical target selection, or
supported owner route may appear.

### Repository, Release, Source, And Qualification Guards

The canonical exact-path line remains exactly the eight genuine stable
releases v0.1.0, v0.2.0, v0.3.0, v1.0.0, v1.1.0, v1.2.0, v1.3.0, and
selected/latest v1.3.1. The separate `github.com/oklog/ulid/v2` line through
v2.1.2 remains ineligible and `/v3` remains absent. Public, enabled,
unarchived, non-fork Apache-2.0 repository `oklog/ulid`, ID 75744278, must
remain owned by `oklog` and default to `main`, with no v0/v1 GitHub Release
objects, prerelease, replacement, retraction, deprecation, eligible alternate
owner, later exact-path stable, or supported tidy-stable request for v0.3.0.

Every completed tag, commit, tree, signature, ancestry, proxy/sumdb,
synthetic-module, archive-to-Git, Git-vendor-exclusion, and license identity
remains a guard. In particular, selected v1.3.1 remains commit
`02a8604050d8466dd915307496174adb9be4593a`, tree
`a88c6a184c68833b0fd6542e490c6fe8f1ccd1d0`, with source/module sums
`h1:EGfNDEx6MqHz8B3uNV6QAib1UR2Lm97sHi3ocA6ESJ4=` /
`h1:CirwcVhetQ6Lv90oh/F+FBtV6XMibvdAFo93nm5qn4U=`, proxy ZIP SHA-256
`40e502c064a922d5eb7f2bc2cda9c6a2a929ec0fc76c9aae4db54fb7b6b611ae`,
and normalized 12-file manifest SHA-256
`7dc92518b194bbf7c62effbe85dde1d79e479e9e055af700339818906b64424c`.

The completed exported API, documentation, caller ownership/mutation,
determinism, concurrency, lifecycle, cleanup, error, global-state, platform,
build-tag, cgo/generated/embed, and Go-floor boundaries remain exact.
V0.3.0 must remain the highest behavior-qualified Go-1.18-compatible stable.
V1.0.0-v1.3.1 must retain their incomplete complete-module closure because
`cmd/ulid` imports undeclared `github.com/pborman/getopt/v2`; their root-
library success cannot override failed complete build/test/vet. No completed
source, behavior, closure, race, vet, or cross-build finding may change.

### Projection, Project, Toolchain, And Earlier-Decision Guards

The selected-version disposable get must retain its recorded 234-module,
3,600-edge, 75/1,068-line root-promotion result with unchanged loads and
negative target why, then tidy to the common 52/948-line projection. The
v0.3.0 disposable get must retain its removal of direct mvn-pom-mutator,
TSDB, go-conntrack, and 24 guarded selections, its 157-module/2,219-edge
unloadable result, and tidy restoration of the common projection with v1.3.1
reselected. Neither projection is retained.

The real project remains exactly 234 selected modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` remain
74/1,067 lines at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`,
and applied tidy remains the common 52/948-line state at
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

The Go 1.18 language/compatibility floor and exact official Darwin arm64 SDK
identities remain: Go 1.26.7 archive/binary SHA-256
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
and Go 1.18.10 archive/binary SHA-256
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Every source/API/CLI/help/launcher/Make/quality contract remains unchanged,
and accepted quality remains 27/27 Q0-Q2 PASS at L2.

All 45 selections earlier than go-conntrack and their 273 sorted incoming
edges remain exact at SHA-256
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Including go-conntrack remains 46 selections/275 edges at SHA-256
`32bd1b6893f9a7964567c1b4bf772d60d690e85509f2ec85fd299d537086e544`.
Forty guarded why results remain negative; only kr/pretty, kr/text, emoji/v2,
Promptui, go-homedir, and mapstructure are positive. Promptui and go-homedir
remain the only guarded repository imports. Emoji/v2, Promptui, go-homedir,
and mapstructure remain the only production/complete-test loaded guarded
modules. Every earlier selection, qualification, owner/request route,
closed-exception boundary, and expiry guard remains separate and exact.

Exact target OSV, GitHub global, and repository advisory results remain
empty. The completed isolated pinned govulncheck target result remains empty,
and unchanged-project 30/22/20/20 non-stdlib populations retain no target
trace. Guard OSV remains limited to Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go vulnerability index remains 518,501 bytes/1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z. Advisory absence is not qualification.

Any target path/version, request, requester import, owner identity or route,
root, why/import/load/runtime fact, source/behavior/closure result, projection,
graph/module/tidy/Go-floor state, earlier or closed-exception guard, advisory/
finding, independent defect, repository/release/owner fact, qualification,
supported-owner result, or compatible genuine route change expires this
exception and requires a fresh ULID dependency evaluation and product decision
before merge. Any earlier guard change separately requires its own fresh
decision. This decision authorizes no ownership study, direct root, v0.3.0 or
`/v2` selection, requester change, dependency edit, patch, vendoring,
workaround, implementation, or transferred exception.

### Guard-Only Revalidation And Handoff

Guard-only revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at evaluation-handoff HEAD
`27046bec05003b317329a1918fd69baf5ffe4eb0`, parent
`c0819eb1524b5a5bf6aabbefebaf185d30d36e88`, tree
`15f3c216faad8c67df78a5e6eaf73ce567940653`. That handoff changes exactly
`codex-dev-start.sh`, the answered ULID evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor.

The reciprocal archive chain, sole NEXT state, launcher/archive prompt mirror,
exact changed set, official Go identities, project counts/hashes/tidy state,
sole target request/route and requester imports, negative why/import/load/
runtime boundary, both earlier edge snapshots, guarded why/import/load sets,
repository/release status, target/guard advisory identities, Go index, and
memberlist CNA identity all reproduce without drift. Completed ULID release,
source, API, behavior, closure, race, cross-build, archive, projection, and
govulncheck work was not repeated, and no ownership study ran.

Final exact Go 1.26.7 module verification, build, count-one tests, race count-
one tests, and vet pass under canonical `umask 022`. No product source,
dependency metadata, target root, requester, or earlier decision changed; no
dependency implementation commit exists. Every task-owned SDK, cache, report,
project copy, and advisory response remained beneath the managed session
scratch root and was removed before handoff.

P7 continues only with the linked reciprocal bounded evaluation of the next
unanswered selected queue item, graph-selected exact
`github.com/pascaldekloe/goe v0.1.0`. That successor was prepared but not
executed and may not reopen or transfer this ULID exception or any earlier
decision. P8 remains queued.
