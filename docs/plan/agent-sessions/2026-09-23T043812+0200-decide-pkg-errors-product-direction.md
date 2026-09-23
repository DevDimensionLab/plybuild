# Agent Session: Decide Pkg Errors Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T043812+0200-decide-pkg-errors-product-direction`
Created: `2026-09-23T04:38:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e4041bf05918ff4f357335b1cbe3e6866fd8c394bf8351656765bae2d21656d0`
Previous: [2026-09-23T035008+0200-evaluate-pkg-errors-dependency.md](2026-09-23T035008+0200-evaluate-pkg-errors-dependency.md)
Next: [2026-09-23T051316+0200-evaluate-pkg-sftp-dependency.md](2026-09-23T051316+0200-evaluate-pkg-sftp-dependency.md)
Outcome: Option 1 selected; exact selected, inherited, indirect, unloaded `github.com/pkg/errors v0.9.1` is retained under a pkg/errors-specific unqualified, non-transferable exception, and only the bounded `github.com/pkg/sftp v1.13.1` evaluation is prepared next.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by making exactly one product decision for exact selected,
inherited, indirect, unloaded `github.com/pkg/errors v0.9.1` from the completed
evaluation. Choose only one direction offered below, preserve every earlier
target-specific guard, update the roadmap and rolling handover, answer this
archive, prepare exactly one reciprocal successor required by the chosen
direction, and make the local handoff commit. Do not reevaluate pkg/errors,
implement a dependency change, execute the named study, evaluate another
dependency group, or begin P8.

# Defensive Scope

This is an ordinary product decision and guard-only revalidation. Use public
metadata, static repository records, and ordinary project graph/build
commands. Do not fuzz, stress, probe resource exhaustion, construct oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, response, report, or project copy must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Verify
containment and remove task-owned scratch evidence before handoff.

# Completed Evaluation Is Final

No genuine supported upgrade, replacement, removal, or exact-retention
direction qualifies. Canonical proxy-latest v0.9.1 is the end of the exact
13-stable line; there is no prerelease, `/v2`, `/v3`, replacement, retraction,
deprecation, redirect, or eligible alternate path. Do not promote master, a
branch, fork, alternate path, pseudo-version, patch, vendor copy, or the
standard library as a release or behavior-equivalent replacement.

Selected v0.9.1 preserves the Go 1.18 floor and passes exact Go 1.18.10
build/count-one/count-ten/race/vet plus all required cgo-disabled cross gates
under both exact SDKs. It does not pass the complete Go 1.26.7 contract:
ordinary tests/race first fail because vet rejects two non-constant Wrapf and
WithMessagef test calls; standalone vet fails the same two sites; and with vet
disabled count-one, every count-ten repetition, and race fail TestStackTrace
because Go 1.26.7 reports a changed nested-function name. Build and cross-
compile still pass. A partial pass does not qualify exact retention.

MVS selects v0.9.1 from ten current/historical requests. Selected SFTP v1.13.1
genuinely imports pkg/errors and requests v0.9.1 through selected Afero
v1.9.4; selected direct Viper v1.15.0 also requests v0.9.1 but is metadata-
only. Kong, both Consul SDK versions, historical SFTP, both Prometheus Common
versions, TSDB, and Zap request v0.8.0/v0.8.1 and genuinely import the target
in the recorded production/test-support boundaries. The completed evaluation
records every exact edge and shortest route through direct go-term-markdown,
Afero, Viper, and mvn-pom-mutator plus historical Chroma, crypt, Consul API,
go-metrics, Prometheus Client, etcd, and their selected requesters.

The target and every requester except selected Viper have negative `go mod
why`. The repository has no target import; target production, complete-test,
and module-backed loads and runtime relevance are zero. There is no main root
and repository history contains no go.mod root. Physical selection and a
genuine SFTP owner do not override the complete upstream failure.

An exact disposable v0.9.1 get adds only an unused indirect root, its source
sum, and one main edge, then tidy returns the common projection with v0.9.1
still selected. Exact removal instead removes direct mvn-pom-mutator,
downgrades direct go-term-markdown and Viper plus Afero, collapses to 109
modules/430 edges, and cannot load the product. Tidy reintroduces exact
v0.9.1 but leaves a non-equivalent 228-module/3,465-edge project with lower
direct Viper, markdown, and Afero selections and changed loads. It violates
earlier direct-owner and selection guards. Neither projection was retained.

The repository/release/tag/commit/tree/signature/ancestry/proxy/sumdb/module/
license/archive identities; selected API and legacy Cause versus Is/As/Unwrap
behavior; empty module dependency closure; build constraints; cgo/generated/
embed boundaries; ordinary behavior; exact SDK gate and cross-build results;
all projections; and advisory results are final. Do not repeat them unless a
narrow guard-only check proves an input changed; stop for a fresh owning
evaluation if it did.

# Authorized Roadmap And Guards

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the nine final target-specific option-1 decisions through Goe, and qualified
no-selection-change modern-go/concurrent and modern-go/reflect2. P8 remains
queued.

Exact Goe v0.1.0, ULID, go-conntrack, mapstructure, go-homedir, Promptui,
emoji/v2, kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-
sequences, gotool, errcheck, httprouter, GLS, go-junit-report, json-iterator,
and clockwork exceptions are final, target-specific, unqualified, and non-
transferable. Every concurrent, reflect2, Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen, broaden, or transfer any decision.

The unchanged project is 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The full graph SHA-256 is
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 pre-Goe guarded selections and 276 sorted incoming edges remain at
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
The separate complete Goe guard is unchanged. Forty-one guarded why results
are negative; only kr/pretty, kr/text, emoji/v2, Promptui, go-homedir, and
mapstructure are positive. Promptui and go-homedir remain the only guarded
repository imports; emoji/v2, Promptui, go-homedir, and mapstructure remain
the only loaded guarded modules.

Exact target OSV, GitHub global, repository advisory, and isolated pinned
govulncheck results are empty. Project populations remain exactly 30/22/20/20
with no target trace. Guard OSV remains limited to the recorded Gorilla
WebSocket and go-retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
GO-2026-6180. The Go vulnerability index remains 518,501 bytes/1,402 records
at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

# Measurements At Start

The pkg/errors evaluation began from clean branch `codex/upgrade-quality` at
handoff HEAD `bf71bc950a601e38233d20e33a50c45f44f0e9c9`, parent
`413fb1a154a5c9a4acfe151dfb9518714a701ed1`, tree
`ab1d295120cac85d6f7c4d633c54d5e919cbafa2`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
evaluation changed no product source or dependency metadata. Independently
verify the new handoff HEAD, parent, tree, exact changed set, ancestry, clean
ordinary and ignored status, reciprocal archive chain, launcher check, exact
Go identities, fixed project/earlier/Goe/target guards, advisories, and final
gates rather than assuming them.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/pkg/errors v0.9.1` without product-source or dependency-metadata
   changes under a pkg/errors-specific, unqualified, non-transferable
   exception. Record that v0.9.1 is latest and has a genuine selected SFTP
   owner but is not qualified because it fails the complete Go 1.26.7 test,
   race, and vet gates. Guard all ten requests, requester import/metadata-only
   boundaries, every route, no main root/history root, why/import/load/runtime
   facts, repository/release/source/behavior/closure identities, both
   projections, all earlier decisions, and no new advisory or independent
   defect. If selected, prepare but do not execute the next bounded P7 queue
   item named by the roadmap after pkg/errors.
2. Authorize exactly one later, measurement-only **Pkg Errors Owner-Route
   Removal Study**. It may determine whether genuine supported versions on
   the existing direct Afero/SFTP, direct Viper, direct go-term-markdown, and
   direct mvn-pom-mutator owner routes can eliminate pkg/errors while
   preserving Go 1.18 and every earlier guard. It may inspect the recorded
   historical requester routes only as evidence. It may not change product
   source, dependency metadata, target roots, Afero, SFTP, Viper, markdown,
   mvn-pom-mutator, any other requester, or an earlier decision. Any future
   implementation requires its own owning product decisions. If selected,
   prepare that named study as the sole successor but do not execute it.
3. Stop P7 unresolved, leave source/dependency metadata and all decisions
   unchanged, prepare no dependency evaluation or study, keep P8 queued, and
   mark the launcher COMPLETE only for the authorized roadmap.

Do not invent a fourth option, combine options, silently retain v0.9.1,
describe it as qualified, safe, fixed, or behavior-equivalent to standard
errors, transfer an earlier exception, or use advisory absence as acceptance.
This decision authorizes no dependency implementation commit.

# Required Reading And Revalidation

Read this archive, the answered pkg/errors evaluation, answered Goe decision/
evaluation, answered ULID and go-conntrack decisions/evaluations, reflect2 and
concurrent evaluations, mapstructure, go-homedir, Promptui, emoji/v2, kr/text,
and kr/pty decisions/evaluations, rolling handover, P7/P8 roadmap, `go.mod`,
and `go.sum`. Earlier evaluation facts are final. Refresh only the narrow
decision inputs: selection and ten requests, routes and requester boundaries,
root/history/why/import/load/runtime facts, project hashes/counts/tidy state,
earlier and Goe guards, repository release status, exact target and guarded
advisories, Go index/memberlist identities, and final exact-Go module
verification, build, count-one tests, race count-one tests, and vet. Stop for
a fresh owning evaluation if a guard changed.

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
metadata, add a root, change an owner/requester, transfer an exception,
reevaluate pkg/errors or an earlier group, evaluate another dependency, write
outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, indirect, unloaded
`github.com/pkg/errors v0.9.1` is explicitly retained without product-source
or dependency-metadata changes under a pkg/errors-specific, unqualified,
non-transferable exception. V0.9.1 is the canonical latest stable and has a
genuine selected SFTP owner, but it is not qualified, safe, or fixed: it fails
the complete required Go 1.26.7 test, race, and vet gates. Physical selection,
transitive presence, a genuine owner, zero loading, and advisory absence do
not qualify or implicitly accept it.

The exception accepts only the completed repository/release/source, API and
ordinary behavior, closure/upstream/vet/race/cross-build, owner/request/route,
why/import/load/runtime, projection/graph/tidy/Go-floor, earlier-guard, and
advisory findings recorded by the answered evaluation. It accepts no
uncharacterized behavior, new advisory, independent defect, future owner
route, direct root, or loaded use. It transfers or broadens no earlier
exception. Option 2's **Pkg Errors Owner-Route Removal Study** was neither
authorized nor run, and option 3 was not selected.

### Exact Selection, Request, And Route Guards

Retention requires exact path/version `github.com/pkg/errors v0.9.1`, with no
direct or indirect main-module target root and no repository-history `go.mod`
root. The ten exact current/historical requests remain:

1. Kong pseudo-version `0548c6b1afae` -> v0.8.1.
2. Consul SDK v0.1.1 -> v0.8.1.
3. Consul SDK v0.8.0 -> v0.8.1.
4. SFTP v1.10.1 -> v0.8.1.
5. Selected SFTP v1.13.1 -> v0.9.1.
6. Prometheus Common v0.4.1 -> v0.8.0.
7. Prometheus Common v0.9.1 -> v0.8.1.
8. Prometheus TSDB v0.7.1 -> v0.8.0.
9. Selected Viper v1.15.0 -> v0.9.1.
10. Zap v1.17.0 -> v0.8.1.

Kong, both Consul SDK `testutil` packages, both SFTP versions, both
Prometheus Common versions, TSDB, and Zap must retain their recorded genuine
production/test-support target imports. Selected Viper's request must remain
metadata-only. Selected SFTP v1.13.1 through selected Afero v1.9.4 remains
the genuine source owner that supplies the maximum v0.9.1 request; Viper's
equal request does not create a source import.

Every recorded shortest route remains required: direct go-term-markdown ->
historical Chroma -> Kong; direct mvn-pom-mutator -> bketelsen/crypt -> Consul
API v1.1.0 -> Consul SDK v0.1.1; mvn-pom-mutator -> historical Viper ->
Consul API v1.12.0 -> Consul SDK v0.8.0; historical Viper -> Afero v1.6.0 ->
SFTP v1.10.1; selected Afero v1.9.4 -> selected SFTP v1.13.1; historical
Viper -> go-metrics -> Prometheus Common v0.9.1 -> Prometheus Client v1.0.0
-> Common v0.4.1; historical Viper -> go-metrics -> Common v0.9.1; direct
mvn-pom-mutator -> TSDB v0.7.1; direct selected Viper v1.15.0; and historical
Viper -> etcd client/pkg v3.5.1 -> Zap v1.17.0.

Only selected Viper remains why-positive. Target and every other requester
remain why-negative. Repository imports, target production entries, complete-
test entries, module-backed entries, and runtime relevance remain zero. A
changed request, requester boundary, owner route, root, why/import/load/
runtime fact, or new supported owner expires this exception.

### Repository, Source, Qualification, And Projection Guards

The canonical exact-path line remains exactly 13 stables through proxy-latest
v0.9.1, with no prerelease, `/v2`, `/v3`, replacement, retraction,
deprecation, redirect, or eligible alternate. Public, enabled, unarchived,
non-fork BSD-2-Clause repository `pkg/errors`, ID 48643510, remains owned by
`pkg`, defaults to `master`, and retains its recorded maintenance-mode status,
13 tags, and 12 GitHub Release objects. Current signed master remains
`87f8819acf6dc28bf5d3c14b334268236d686f48`, tree
`60652f0e917d39e5d310641579b61c4682d64164`.

Every completed tag, release, commit, tree, signature, ancestry, synthetic
module, proxy/sumdb, license, archive-to-Git, symlink, and submodule identity
remains a guard. Selected v0.9.1 remains signed commit
`614d223910a179a466c1767a985424175c39b465`, tree
`6dd01fd9b7f97a850cc87788579cfc01fd6431fd`, with the recorded source/module
sums, proxy ZIP, and normalized 16-file manifest identities.

The completed exported API, documentation, legacy Cause versus standard
Is/As/Unwrap behavior, stack and format semantics, empty module dependency
closure, build constraints, cgo/generated/embed boundaries, and ordinary
behavior remain exact. V0.9.1 must retain its exact Go 1.18.10 green build,
count-one/count-ten, race, and vet results and all recorded cross gates. It
must also retain the complete Go 1.26.7 failure: vet rejects the two
non-constant Wrapf/WithMessagef test calls, while test and race with vet
disabled still fail TestStackTrace on the changed nested-function name. A
partial pass cannot qualify retention, and standard `errors` remains
behaviorally non-equivalent to the guarded stack/Cause/format API.

The selected-version disposable get must continue to manufacture only an
unused indirect root, selected source sum, and one main edge before tidy
reaches the common projection with v0.9.1 still selected. Exact removal must
retain its recorded direct mvn-pom-mutator removal, direct go-term-markdown,
Viper, and Afero downgrades, unloadable 109-module/430-edge raw state, and
non-equivalent 228-module/3,465-edge tidy state that reintroduces v0.9.1.
Neither projection is retained.

### Project, Earlier-Decision, And Advisory Guards

The real project remains exactly 234 selected modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` remain
74/1,067 lines at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The full graph remains SHA-256
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
The 432-line tidy diff remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`,
and applied tidy remains the common 52/948-line state at
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

The Go 1.18 language/compatibility floor and official Darwin arm64 identities
remain exact: Go 1.26.7 archive/binary SHA-256
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
and Go 1.18.10 archive/binary SHA-256
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.

All 47 pre-Goe guarded selections and 276 incoming edges remain exact at
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
The separate Goe guard retains exact v0.1.0, its four requests, go-metrics and
Consul API test-import boundaries, memberlist metadata-only boundaries,
negative why/import/load/runtime state, and every recorded route. All 41
guarded why-negative results, the six positive results, two guarded repository
imports, and four loaded guarded modules remain exact. Every earlier
qualification and target-specific exception remains separate, final, and
non-transferable.

Exact target OSV, GitHub global, repository advisory, and isolated pinned
govulncheck target results remain empty. Project govulncheck populations
remain exactly 30 module, 22 package, 20 symbol, and 20 test-symbol OSVs with
no pkg/errors trace. Guard OSV remains limited to Gorilla WebSocket
`GO-2026-6278` / `GHSA-w67g-5rqw-f597` and go-retryablehttp
`GO-2024-2947` / `GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains
`GO-2026-6179` and `GO-2026-6180`. The Go vulnerability index remains
518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

Any path/version, request/requester boundary, owner/route, root/why/import/
load/runtime, repository/release/source/behavior/closure, projection/graph/
module/tidy/Go-floor, earlier guard, advisory/finding, independent defect,
qualification, supported-owner, or compatible-route change expires this
retention and requires a fresh pkg/errors evaluation and product decision
before merge. No owner-route study, direct root, patch, vendor copy,
standard-library replacement, requester change, dependency edit,
implementation, or transferred exception is authorized.

### Guard-Only Revalidation And Handoff

Guard-only revalidation began clean on branch `codex/upgrade-quality` at
evaluation-handoff HEAD
`7b078f6d42a02c57feb7a45dd9a821e02d4b0f14`, parent
`bf71bc950a601e38233d20e33a50c45f44f0e9c9`, tree
`f1bc36ceef014e57a2ca12c0ccb37dcc2e0dc8fe`. That handoff changes exactly
the launcher, answered pkg/errors evaluation archive, this then-NEXT decision
archive, rolling handover, and roadmap. Google UUID implementation commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.

The reciprocal archive chain, sole NEXT state, launcher/archive prompt mirror,
exact changed set, branch, ordinary/ignored cleanliness, official Go
identities, project counts and hashes, full graph, tidy state, all ten target
requests/routes/import boundaries, root/history/why/import/load/runtime facts,
pre-Goe and separate Goe guards, repository/release status, target/guard
advisories, Go-index, memberlist-CNA, and exact-Go final gates reproduce with
no guard drift. Completed pkg/errors evaluation work was not repeated, no
independent defect was found, and no owner-route study ran.

A module-download guard initially populated missing source sums in the real
`go.sum`; those task-created additions were immediately removed with the patch
editor, restoring the exact 1,067-line hash above. All subsequent mutable Go
commands ran only in disposable project copies. Final exact Go 1.26.7 module
verification, build, count-one tests, race count-one tests, and vet pass under
canonical `umask 022`. No product source, dependency metadata, target root,
owner/requester, or earlier decision changed.

P7 continues only with the prepared bounded evaluation of the first
unanswered selected queue item after pkg/errors, exact
`github.com/pkg/sftp v1.13.1`. That successor was prepared but not executed.
It may not reopen or transfer this exception or any earlier decision. P8
remains queued.
