# Agent Session: Decide Kr Fs Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-22T081929+0200-decide-kr-fs-product-direction`
Created: `2026-09-22T08:19:29+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `81553d61fc720f499781599c3c3d138786094bf3702827e391ad98de5c68bb38`
Previous: [2026-09-22T071843+0200-evaluate-kr-fs-dependency.md](2026-09-22T071843+0200-evaluate-kr-fs-dependency.md)
Next: [2026-09-22T084421+0200-evaluate-kr-logfmt-dependency.md](2026-09-22T084421+0200-evaluate-kr-logfmt-dependency.md)
Outcome: Selected option 1: explicitly retained exact selected, inherited, unloaded kr/fs v0.1.0 as unqualified under a kr/fs-specific non-transferable exception, left product source and dependency metadata unchanged, and prepared one bounded kr/logfmt evaluation.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with one bounded product decision for exact selected, inherited,
unloaded `github.com/kr/fs v0.1.0`. Choose one authorized direction from the
completed evaluation: explicitly retain exact v0.1.0 under a new target-
specific exception, authorize exactly one later measurement-only owner/request
study, or stop P7 unresolved. Do not repeat the completed behavior evaluation,
implement a parent or workaround, combine another dependency group, or begin
P8.

# Defensive Scope

This is an ordinary dependency product decision. Use public metadata, static
records, project graph/build commands, and public advisory evidence. Do not
fuzz, stress, probe resource exhaustion, generate oversized or deeply nested
trees, generate adversarial malformed paths or filesystem states, reproduce a
security issue, or perform security or exploitability analysis. The completed
source and ordinary filesystem behavior evaluation is final and must not be
repeated.

Every disposable cache, tool, archive, report, project copy, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact selected,
inherited, unloaded go-windows-terminal-sequences v1.0.1. Every earlier
outcome and exception is final under its own guards. P8 remains queued.

The go-windows-terminal-sequences option-1 decision retains exact v1.0.1 only
while both historical Logrus requests and genuine owner routes, negative why,
zero import/load/runtime state, graph/module/tidy/Go-floor facts, every earlier
guard, and no-new-finding/release/owner conditions remain exact. No go-windows-
terminal-sequences, gotool, errcheck, httprouter, jtolds/gls,
go-junit-report, json-iterator, clockwork, demangle, strcase, memberlist, or
other exception transfers to kr/fs.

# Measurements At Start

The completed evaluation began from clean HEAD
`f12458f3e753abb2fe6b0420d345e3ef39b3af7f`, parent
`f5d1d5de0197fead75749e749a9b9bf24b66f0bb`, tree
`7d6cb2df56d735966ed27574c97d00cfceca218f`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. The evaluation made no source or
dependency-metadata change and prepared this decision-only handoff commit.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The project remains 234 selected modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sums, and the recorded 432-line tidy projection. Accepted
quality remains 27/27 Q0-Q2 PASS at L2. Verify the new handoff HEAD, parent,
tree, exact changed set, and clean ordinary and ignored status at start rather
than assuming their values.

# Completed Evaluation

No exact-path stable kr/fs release qualifies. The public proxy exposes only
stable v0.1.0 and resolves `@latest` to it. Its exact module file declares only
`module "github.com/kr/fs"`; it has no Go directive, requirements,
retractions, deprecation, or replacement. Sumdb records exact source/module
sums. Exact `go-import` metadata resolves without redirect to the public,
enabled, unarchived, non-fork BSD-3-Clause `kr/fs` repository.

The sole annotated tag is unsigned and peels to unsigned commit
`1455def202f6e05b95cc7bfc7e8ae67ae5141eba`; `master` is the same commit.
There are no GitHub Releases. Default `main` is one unsigned README-only
commit later and resolves as pseudo-version
`v0.1.1-0.20210218185759-c64b65e7619f`; all Go source is identical. Its
recommendation of `kr.dev/walk` names a distinct module/repository and is an
unauthorized alternate path. Exact `/v2` and `/v3` lines do not exist. Proxy
and tagged Git regular files agree byte for byte.

V0.1.0 is one standard-library-only package with no declared Go floor and
preserves Go 1.18. It exports `FileSystem`, `Walker`, `Walk`, `WalkFS`, and the
Walker access/control methods. Exact Go 1.26.7 and contained Go 1.18.10 both
pass upstream build, tests, repeated race tests, vet, supported cross-builds,
and all non-disqualifying bounded behavior coverage. There is no cgo,
generated source, global mutable state, network/subprocess boundary, or
library-owned caller resource.

The only stable nevertheless fails an exported ordinary behavior contract.
`Walker` promises lexical-order traversal unconditionally. Internal `Walk`
uses sorted `ioutil.ReadDir`, but `WalkFS` accepts any valid `FileSystem` and
never sorts the returned entries. Under both exact SDKs, a bounded in-memory
filesystem returning ordinary root entries `[b, a]` yields
`[root, root/b, root/a]`, not promised `[root, root/a, root/b]`. The
README-only main pseudo-version has identical source and the same defect. The
completed source/test closure, API, error/zero-value behavior, ownership,
mutation, determinism, concurrency, global-state, lifecycle, platform, and
external-boundary evaluation is final.

MVS selects exact v0.1.0 through exactly two requests: SFTP v1.13.1 and
v1.10.1. The current route is main -> direct Afero v1.9.4 -> SFTP v1.13.1 ->
target. The historical route is main -> direct mvn-pom-mutator v0.2.3 ->
Viper v1.10.1 -> Afero v1.6.0 -> SFTP v1.10.1 -> target. Both SFTP releases'
production clients implement `kr/fs.FileSystem`, call `fs.WalkFS`, and retain
server reply order rather than sorting. The failed contract is therefore
relevant to the actual requester even though neither SFTP version loads in the
project.

The target has negative why, zero repository imports, zero production or
complete-test loads, and no runtime reachability. Current Afero loads only its
root, `internal/common`, and `mem` packages, not its SFTP adapter. These facts
bound present exposure; they are not qualification or implicit acceptance.

Disposable exact-v0.1.0 and main-pseudo direct roots manufactured a main edge
and sums but preserved all project loads. Tidy removed both roots, restored
the common 52-line/948-line tidy state, and restored inherited v0.1.0 in the
pseudo case because no genuine owner requests it. No projection was retained.
Product source and dependency metadata remain unchanged.

Fresh OSV and GitHub exact-version results are empty for both candidates, the
repository advisory endpoint is empty, and pinned isolated govulncheck scans
have no target result. Base and disposable project populations are identical
with no target trace. The 518,501-byte, 1,402-record Go index remains exact;
guard OSV retains only the recorded Gorilla/retryablehttp pairs, x/mod retains
GO-2026-6179 and GO-2026-6180, and the PUBLISHED memberlist CNA response
remains exact. Advisory absence does not qualify the ordinary behavior.

The unchanged project remains 234 modules, 3,599 graph edges, 355 production
and 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, exact module hashes, and the recorded 432-line tidy
projection. All 34 guarded selections and 223 incoming edges remain exact at
sorted SHA-256
`8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`;
all guarded why results are negative and imports/loads zero. Accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Required Product Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, unloaded v0.1.0 without
   source or dependency-metadata changes under a kr/fs-specific,
   non-transferable exception. The decision must call it unqualified, accept
   only the completed lexical-order finding plus the completed related
   evaluation, and define exact expiry guards for path/version, both SFTP
   requests and their genuine Afero/Viper/mvn-pom owner routes, root/import/
   load/runtime, graph/module/tidy/Go-floor, every earlier guard, advisories/
   findings/releases/owners, and any compatible genuine route to a qualified
   exact-path release.
2. Authorize exactly one later bounded measurement-only owner/request study.
   Name the existing current Afero v1.9.4 -> SFTP v1.13.1 request route and
   the exact question to measure: whether a genuine supported tidy-stable
   owner change removes kr/fs or selects a qualified exact-path stable release
   while preserving the Go floor and every contract. Do not run it, implement
   a parent change, add a direct target root, transfer an exception, or imply
   approval of graph/source/metadata changes.
3. Stop P7 unresolved, record the blocker, and prepare no dependency or P8
   implementation.

The decision must not call physical MVS selection, negative why, zero loading,
or advisory absence qualification. It must not select the pseudo-version
directly without a genuine owner, promote `kr.dev/walk`, select a fork,
branch, replacement, patch, wrapper, or workaround; add a direct edge; change
Afero, SFTP, Viper, mvn-pom-mutator, the Go floor, product source, dependency
metadata, or another selection; reopen an earlier decision; or begin P8.

# Role And Boundaries

This is a decision-recording session, not a renewed target audit or an
implementation. Reuse the completed evaluation and make exactly one choice.
Do not broaden it into direct use, a pseudo-version without a genuine owner,
alternate-path promotion, parent removal, patching, forking, wrapping,
replacement, a Go-floor change, unrelated-module authorization, another
dependency group, or P8.

# Required Reading

Read this decision archive, the answered target evaluation, go-windows-
terminal-sequences, gotool, and errcheck decisions/evaluations, direct Logrus
and relevant Afero/Viper/mvn-pom retained-module records, rolling handover,
roadmap, `go.mod`, and `go.sum`. Verify the new handoff HEAD/parent/tree and
exact changed set, reciprocal archive chain, latest Google UUID implementation
ancestry, exact Go 1.26.7 identity, launcher check, module hashes/counts/tidy
projection, both requests and genuine routes, target/owner why-import-load
state, all 34 guarded selections and 223-edge snapshot, and fresh target/guard
advisory identities. Treat the completed behavior and disposable project
evaluation as final.

# Three Moves

First, revalidate only the decision guards and public metadata needed to know
whether the evaluation remains current. Second, record exactly one authorized
product direction without changing source or dependency metadata or executing
an owner study. Third, update the roadmap and rolling handover, answer this
archive, prepare exactly one reciprocal successor only if the chosen direction
requires one, verify containment, and commit the documentation handoff without
executing the successor.

# Automatic Handoff

Make the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, repeat completed behavior, implement a parent or
workaround, combine another dependency group, reopen earlier work, write
outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 was selected on 2026-09-22. Exact selected
`github.com/kr/fs v0.1.0` is explicitly retained as inherited and unloaded
without changing product source, `go.mod`, or `go.sum`. It remains unqualified
and is not described as secure. Physical MVS selection, the negative
`go mod why` result, zero loading, and advisory absence bound present exposure
but did not qualify or silently authorize it. No exact-path stable release
satisfies the exported lexical-order traversal contract.

The kr/fs-specific, non-transferable exception accepts only the completed
lexical-order failure plus the completed repository, archive, release, module,
Go-floor, source, exported API, ordinary filesystem and error/zero-value
behavior, caller ownership, mutation, determinism, concurrency/global-state,
lifecycle, platform, external-boundary, MVS, loading, project, vulnerability,
and related evaluation findings. It accepts no uncharacterized behavior, new
advisory, or independently discovered defect. It promotes no `kr.dev/walk`,
alternate path, fork, branch, prerelease, pseudo-version, replacement, direct
root, patch, wrapper, or workaround and transfers no earlier dependency
exception.

This repository owns only the explicit decision to tolerate exact inherited
and unloaded v0.1.0 in the verified graph. The exception remains valid only
while every one of these facts remains exact:

- selected exact path and version `github.com/kr/fs v0.1.0`;
- both selected-version requests and requester identities:
  `github.com/pkg/sftp@v1.13.1 -> github.com/kr/fs@v0.1.0` and
  `github.com/pkg/sftp@v1.10.1 -> github.com/kr/fs@v0.1.0`;
- the genuine current route from the main module through direct
  `github.com/spf13/afero v1.9.4`, SFTP v1.13.1, and the target;
- the genuine historical route through direct
  `github.com/devdimensionlab/mvn-pom-mutator v0.2.3`, historical
  `github.com/spf13/viper v1.10.1`, Afero v1.6.0, SFTP v1.10.1, and the
  target;
- both SFTP clients continuing to implement `kr/fs.FileSystem`, call
  `fs.WalkFS`, and preserve server reply order, while neither SFTP version
  loads and selected Afero loads only its root, `internal/common`, and `mem`
  packages rather than its SFTP adapter;
- no direct target root, a negative target why result, zero repository imports,
  zero production or complete-test package loads, and no target runtime
  reachability;
- 234 selected modules, 3,599 graph edges, 355 production entries, 429
  complete-test entries, 197 module-backed complete-test entries across 41
  loaded modules, 1,067 sum lines, and the exact 432-line tidy projection at
  SHA-256
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
- `go.mod` and `go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and the common tidy-applied 52/948-line hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`;
- the declared Go 1.18 floor and exact Go 1.26.7 archive/binary identities;
- every earlier target-specific guard, with all 34 earlier guarded selections
  exact, all 34 why results negative, guarded repository imports and loads
  zero, and the 223 sorted incoming edges at SHA-256
  `8cb329e95dcd81165af873c9bb4f11c85f61f7d68f39794394ab16a49f5cddc6`;
  including kr/fs yields 35 guarded selections and 225 incoming edges at
  SHA-256
  `5dae96771b3994a9ce1999f8d0487f152d94adb3b2d0bb467b18054be39dc284`;
  and
- no new kr/fs or requester-closure advisory, independent finding, repository
  owner or status change, exact-path stable release, requester or owner,
  qualified release, genuine supported tidy-stable owner, or compatible
  genuine route to a qualified exact-path release.

Any target version/path, request, SFTP, Afero, Viper, mvn-pom-mutator or other
owner identity/route, target interface/use/order behavior, root, import, load,
runtime, graph, module-hash, tidy, Go-floor, earlier-guard, advisory,
independent finding, repository/release/owner, qualification, or compatible-
route change expires the exception and requires a fresh kr/fs dependency and
product decision before merge. The decision authorizes no owner study,
workaround, direct target root, parent change, alternate path, product-source
or dependency-metadata change, unrelated selection, or implementation.

Guard-only revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at handoff HEAD
`8e576d81de9521eb81aacdcc540d3f4b014baa11`, parent
`f12458f3e753abb2fe6b0420d345e3ef39b3af7f`, tree
`4c8edcd542d43b432dac579bdb819df8f6b42e9d`. That handoff changes exactly
`codex-dev-start.sh`, the answered kr/fs evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. Its reciprocal 252-archive
chain, sole NEXT state, launcher/archive prompt mirror, exact changed set, and
launcher check passed.

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum`; it remains an ancestor of the handoff. A freshly unpacked official
Go 1.26.7 archive and binary retained SHA-256 values
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
All guard commands used that binary with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, no ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and
`umask 022`.

The baseline counts, module hashes, tidy projection, both target requests, and
both genuine owner routes reproduce. Target and SFTP why results remain
negative; repository target imports and target/SFTP production/complete-test
loads remain zero; and Afero retains only the three recorded loaded packages.
All 34 earlier guarded selections remain exact with negative why and zero
imports/loads; both the 223-edge earlier snapshot and 225-edge snapshot
including kr/fs reproduce.

Fresh proxy metadata still exposes only stable v0.1.0 and resolves `@latest`
to it. GitHub still reports the exact public, enabled, unarchived, non-fork
BSD-3-Clause repository with one tag, no Releases, and `main` at
`c64b65e7619f39d51006be99f6699903dba2c890`. Exact target OSV and GitHub
global/repository advisory results remain empty. Exact OSV across all 34
earlier guarded selections retains only Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go vulnerability index remains 518,501 bytes and 1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
response remains byte-exact at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

The completed kr/fs behavior evaluation and disposable project projections
were not repeated, option 2 was not run, and no source, dependency metadata,
parent, toolchain declaration, or earlier guard changed. No changed-selection
scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at L2. Final
unchanged-project exact-Go module verification, build, count-one tests, race
count-one tests, and vet pass. The reciprocal 253-archive chain, single NEXT
state, launcher/archive prompt mirror, exact documentation-only changed set,
diff checks, and launcher check pass after the handoff edit. Every task-owned
tool, cache, report, archive, and advisory snapshot was contained beneath the
managed session scratch root and removed; only the pre-existing launcher-owned
Node compile cache remains there.

P7 continues only with the reciprocal bounded evaluation of the next selected
queue item, `github.com/kr/logfmt
v0.0.0-20140226030751-b84e30acd515`, linked above. That successor was prepared
but not executed. P8 remains queued.
