# Agent Session: Decide Sergi Go-Diff Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T011213+0200-decide-sergi-go-diff-product-direction`
Created: `2026-09-24T01:12:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f2816b5fd3f85b7fcd804f60be6fa8461da99002f75dfac0239556ee0fec152c`
Previous: [2026-09-24T002515+0200-evaluate-sergi-go-diff-dependency.md](2026-09-24T002515+0200-evaluate-sergi-go-diff-dependency.md)
Next: [2026-09-24T014102+0200-evaluate-shurcool-sanitized-anchor-name-dependency.md](2026-09-24T014102+0200-evaluate-shurcool-sanitized-anchor-name-dependency.md)
Outcome: Option 1 is final. Exact selected indirect unloaded go-diff v1.2.0 remains unchanged only under its own unqualified, non-transferable exception; no study, dependency or source change, transferred exception, other dependency group, or P8 work was authorized.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product-direction decision for exact
selected indirect `github.com/sergi/go-diff v1.2.0`. The completed fresh
evaluation found no fully qualified canonical Go-1.18-compatible stable and no
changed selection with a genuine supported tidy-stable owner. Choose and record
exactly one of the three authorized directions below. Do not repeat the
evaluation, implement a dependency or source change, combine another group,
launch a study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only completed
public release/source/build/graph/projection/advisory evidence and bounded
read-only guard checks. Do not fuzz, stress, probe resource exhaustion, create
oversized, deeply nested, cyclic, malformed, adversarial, or escape-sequence
payloads, reproduce a security issue, or perform security or exploitability
analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Set every tool
temp/cache root explicitly, verify containment, and remove task-owned scratch
evidence before handoff.

# Completed Evaluation

Exact go-import metadata maps the target to the public, enabled, unarchived,
non-fork MIT `sergi/go-diff` repository, ID 7148036, on default branch
`master`. Default HEAD is highest stable v1.4.0. The exact-path canonical line
has exactly six stables: v1.0.0, v1.1.0, selected v1.2.0, v1.3.0, v1.3.1, and
v1.4.0. There is no prerelease, replacement, retraction, deprecation, `/v2`,
or `/v3` line. V1.0.0 uses a proxy-synthesized module file; v1.1.0-v1.3.1
declare Go 1.12 and v1.4.0 declares Go 1.13, so every stable preserves the Go
1.18 floor.

All exact tag/commit/tree/signature/ancestry, proxy/sumdb/archive-to-Git,
module/license, package, regular-archive-entry, generated/source/build-boundary,
and release chronology identities are recorded in the answered evaluation.
Selected v1.2.0 remains unsigned commit
`0a651d56613f9de4bed8b9c4769b776ef168bfca`, tree
`fabcd3c10c75b0ab23b87194198b5b7f708a26ae`, with source/module sums
`h1:XU+rvMAioB0UC3q1MFrIQy4Vo5/4VsRDQQXHsEya6xQ=` /
`h1:STckp+ISIX8hZLjrqAeVduY0gWCT9IjLuqbuNXdaHfM=`. Highest v1.4.0 remains
signed commit `57c41f4cb9849a2e83cdbd7644b31e6d7a7e2586`, tree
`3d545e152d9f1386c187dcdc77bc83803741164f`, with source/module sums
`h1:n/SP9D5ad1fORl+llWyN+D6qoUETXNZARKjyY2/KVCw=` /
`h1:A0bzQcvG0E7Rwjx0REVgAGH58e96+X0MeOfepqsbeW4=`.

The package exposes configurable diff-match-patch text/rune diffing, cleanup,
fuzzy matching, patch construction/application, delta encoding, and
unified-diff-like rolling-context patch text. Caller ownership, mutation,
determinism/timeout, concurrency, error/panic, and no-cleanup-lifecycle facts
are recorded. V1.4.0 adds Unicode constants and removes selected v1.2.0's
exported `IndexSeparator`, so it is not API-identical.

V1.0.0 has an incomplete test closure because its synthesized module metadata
omits imported Testify. V1.1.0-v1.4.0 have ten-module declared build lists and
standard-library-only production imports. Under exact verified Go 1.18.10 and
Go 1.26.7 all six pass module verification, production build, and all 120
cgo-disabled production compilation rows across the ten recorded targets.
V1.0.0 cannot vet or compile tests; v1.1.0 test packages compile with vet
disabled but mandatory vet fails under both SDKs; v1.2.0-v1.4.0 pass safe vet.

Every stable test suite contains explicit invalid URL-escape, invalid UTF-8,
invalid patch, and timeout-amplification inputs. V1.2.0+ adds a 2,488,955-byte
text fixture and 4,680,695-byte generated Go fixture; v1.4.0 also has an
exhaustive Unicode-range loop. Static review found those boundaries before
execution. Complete count-one, repeated, and race tests were stopped for every
release. V1.2.0+ test-compilation rows were also stopped before ingesting the
oversized fixture. Partial safe evidence is not qualification. No canonical
stable fully qualifies, including highest compatible v1.4.0.

The current graph has exactly three target requests and routes: direct main ->
v1.2.0; direct Assert v1.0.0 -> v1.2.0; and main -> direct go-term-markdown
v0.1.4 -> historical Chroma v0.7.1 -> v1.0.0. Assert genuinely imports the
target in production; Chroma v0.7.1 has zero target source imports and its
request is metadata-only. Go-term-markdown genuinely imports Chroma. Target
and Assert why are negative; go-term-markdown and Chroma why are positive.
Target and Assert have zero repository imports, production/complete-test/
module-backed loads, and runtime relevance. Selected Chroma has 33 production
and 33 complete-test entries and runtime relevance through go-term-markdown,
but that does not make its historical target metadata edge genuine ownership.

All 84 historical module checkpoints and route epochs reproduce. Chroma first
becomes a root at `4ea30e28145509bf937137d8e35dbcb980fad323`; target and Assert
first become main roots together at
`0781fd6cdb625623c6d726743bf58352113d0ccb`. Exact selected get changes no byte
or selection. Ordinary tidy removes the target and Assert roots, selects the
2017 Assert pseudo-version and go-diff v1.0.0 through Chroma's metadata-only
route, and restores the common projection.

A v1.4.0 get changes only the target selection and adds its two sums while
preserving the Go floor and project loads. Ordinary tidy discards v1.4.0 and
again reaches the common projection at protected hashes, selecting v1.0.0.
No requester asks for v1.3.0, v1.3.1, or v1.4.0. Thus no changed selection has
a genuine supported tidy-stable project owner. No projection, source, root, or
dependency change was retained.

Exact-version target OSV and narrow GitHub global/repository responses are
empty without implying safety or qualification. Pinned govulncheck v1.8.0 at
database timestamp 2026-09-16T18:00:43Z has no non-standard-library production
finding for the loadable focal closures, but incomplete/stopped test closures
prevent a complete qualification claim. The project remains 30/22/20/20 with
no go-diff, Blackfriday, fastuuid, TSDB, or Procfs trace. Client_golang v1.4.0
retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. Corrected Go
index/CNA identities, exact SDKs, real 234/3,599/355/429/197/41/1,067 state,
module/graph/common-tidy hashes, Go floor, all 47 pre-Goe selections/276 edges,
separate later requests/exceptions, every earlier decision, and 27/27 Q0-Q2
PASS at L2 remain exact.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected, indirect, unloaded
   `github.com/sergi/go-diff v1.2.0` unchanged under a go-diff-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   safe, or fixed and lacks a genuine supported tidy-stable owner. Preserve
   the exact three requests/routes, Assert's genuine import, Chroma's metadata-
   only boundary, target/requester why and load/runtime facts, all 84 historical
   checkpoints, root history, selected sums, release/source/API/behavior/
   closure/native/test-scope/cross/projection/advisory identities, and every
   earlier guard as expiry conditions. Do not transfer another exception.
2. Authorize exactly one later measurement-only **Go-Diff Request/Ownership
   Study** to determine whether redundant target/Assert roots can be removed,
   a supported genuine requester can own a qualified release, or the target can
   be eliminated while preserving product behavior, direct roots, the Go 1.18
   floor, and every earlier guard. Do not run the study, grant an exception,
   change dependency/source files, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the exact selection is unloaded and runtime-
irrelevant, its direct Assert requester genuinely imports it and remains in an
active public repository, and every qualification/ownership boundary can be
made explicit. This is a product acceptance decision, not a qualification or
safety claim.

# Required Reading And Handoff

Start only from the clean go-diff evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch and Google UUID ancestry, reciprocal archive
chain, sole NEXT state, launcher mirror/check, ordinary and ignored
cleanliness, official SDK identities, real project counts/hashes, common tidy
projection, Go 1.18 floor, all target/requester request/import/route/relevance/
root facts, target/closure/project advisory identities, and every earlier
guard. Stop for a fresh owning decision if any protected input changed.

Preserve exact selected Blackfriday v2.1.0, TSDB v0.7.1, Procfs v0.0.8,
Common v0.9.1, client_model v0.2.0, and client_golang v1.4.0 only under their
own separate unqualified, non-transferable exceptions. Preserve Complete,
go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and every earlier
qualified or excepted result under its exact guards. Preserve fully qualified
fastuuid v1.2.0 without adding a root. Do not run rejected studies, select
rejected candidates, reopen a completed module, or transfer an exception. P8
remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, projection, native, cross, or advisory-matrix work.
Record the chosen product direction, answer this archive, update the roadmap
and rolling handover, verify containment and cleanup, run final exact-Go-1.26.7
project module verification, build, count-one tests, race count-one tests, and
vet, and make one local handoff commit.

# Three Moves

Choose only one numbered direction. Option 1 is the recommended bounded
acceptance; option 2 authorizes only a later measurement study; option 3 stops
unresolved. None authorizes work in another dependency group or P8.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8. If and only if the chosen direction closes go-diff, prepare but do not
launch the next bounded P7 evaluation of exact selected
`github.com/shurcooL/sanitized_anchor_name v1.0.0`.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, indirect, unloaded
`github.com/sergi/go-diff v1.2.0` is explicitly retained unchanged under a
go-diff-specific, unqualified, non-transferable exception. It is not
qualified, safe, or fixed, and it lacks a genuine supported tidy-stable owner.
Its partial passing build evidence and empty advisory results do not alter
those boundaries.

Option 2's **Go-Diff Request/Ownership Study** is not authorized or run.
Option 3 is not selected because the exact selection is unloaded and runtime-
irrelevant, its direct Assert requester genuinely imports it and remains in
an active public repository, and retention can be bounded by exact expiry
guards. No dependency or source change, target or requester root change,
projection, implementation pre-authorization, transferred exception, other
dependency group, or P8 work is included.

### Exact exception and expiry boundary

Retention requires the exact three requests and complete routes to remain
unchanged:

1. main -> go-diff v1.2.0;
2. main -> Assert v1.0.0 -> go-diff v1.2.0; and
3. main -> go-term-markdown v0.1.4 -> Chroma v0.7.1 -> go-diff v1.0.0.

Assert v1.0.0 must continue genuinely importing go-diff in production.
Chroma v0.7.1 must retain zero target source imports and its metadata-only
v1.0.0 request; go-term-markdown v0.1.4 must continue genuinely importing
Chroma and requesting v0.7.1. The public Assert, Chroma, and go-term-markdown
repositories must remain enabled, unarchived, non-fork, and supported. Assert's
genuine import does not own a tidy-stable selection while ordinary tidy
removes its direct root, and Chroma's metadata-only edge does not establish
genuine ownership.

Target and Assert `go mod why -m` must remain negative; selected Chroma and
go-term-markdown why must remain positive. Repository imports of target and
Assert, their production/complete-test/module-backed loads, and their runtime
relevance must remain zero. Selected Chroma must retain exactly 33 production
and 33 complete-test entries and runtime relevance only through go-term-
markdown. Main must continue first rooting target and Assert together only at
`0781fd6cdb625623c6d726743bf58352113d0ccb`; Chroma must first become a root
only at `4ea30e28145509bf937137d8e35dbcb980fad323`. All 84 historical module
checkpoints and their recorded route epochs remain incorporated by reference.
Any request, route, requester import or metadata boundary, repository support,
why/import/load/runtime fact, root history, or project sum change expires this
exception and requires a fresh owning evaluation and explicit product decision
before merge.

Selected v1.2.0 remains unsigned commit
`0a651d56613f9de4bed8b9c4769b776ef168bfca`, tree
`fabcd3c10c75b0ab23b87194198b5b7f708a26ae`, with source and module sums
`h1:XU+rvMAioB0UC3q1MFrIQy4Vo5/4VsRDQQXHsEya6xQ=` /
`h1:STckp+ISIX8hZLjrqAeVduY0gWCT9IjLuqbuNXdaHfM=`. The protected `go.sum`
must continue containing exactly those selected lines. Any identity or sum
change expires retention.

Every canonical owner, release, and source identity recorded in the answered
evaluation remains an expiry guard: exact go-import metadata; public, enabled,
unarchived, non-fork MIT `sergi/go-diff` repository ID 7148036 on `master`;
exact v1.0.0, v1.1.0, v1.2.0, v1.3.0, v1.3.1, and v1.4.0 stable line; absent
prerelease, replacement, retraction, deprecation, `/v2`, and `/v3` lines; and
the preserved Go 1.18 floor. Exact tag/commit/tree/signature/ancestry,
proxy/sumdb/archive-to-Git, module/license, package, regular-archive-entry,
generated/source/build-boundary, release chronology, exported API, and
documented ordinary behavior evidence remains incorporated by reference and
non-transferable. Highest v1.4.0 remains signed commit
`57c41f4cb9849a2e83cdbd7644b31e6d7a7e2586`, tree
`3d545e152d9f1386c187dcdc77bc83803741164f`, and is not API-identical to
selected v1.2.0 because it adds Unicode constants and removes exported
`IndexSeparator`.

The qualification boundary remains exact. V1.0.0 retains an incomplete test
closure because its synthesized module metadata omits imported Testify;
v1.1.0 retains mandatory vet failures under both SDKs. V1.2.0-v1.4.0 pass
safe vet, and all six stables retain passing module verification, production
build, and 120 cgo-disabled production compilation rows. Every stable's tests
still include explicit invalid URL-escape, invalid UTF-8, invalid patch, and
timeout-amplification inputs; v1.2.0+ retains multi-megabyte fixtures and
v1.4.0 an exhaustive Unicode-range loop. Complete count-one, repeated, race,
and affected test-compilation rows remain stopped at the defensive boundary.
No upstream test was run in this decision session, partial results are not
qualification, and any newly fully qualified canonical stable expires the
exception.

The exact projection boundary remains an expiry guard. Exact selected get
must continue changing no selection or project byte. Ordinary tidy must
continue removing target and Assert direct roots, selecting the 2017 Assert
pseudo-version and go-diff v1.0.0 through Chroma's metadata-only route, and
restoring the common 52/948-line, 234-module/3,557-edge projection. A v1.4.0
selection must continue changing only target and its two sums before ordinary
tidy discards it. No requester asks for v1.3.0, v1.3.1, or v1.4.0. No
projection is authorized or retained.

Exact empty target OSV and narrow GitHub global/repository responses remain
guards, and their absence does not imply safety or qualification. Pinned
govulncheck v1.8.0 at database timestamp 2026-09-16T18:00:43Z retains no
non-standard-library production finding for the safely loadable focal
closures; incomplete or stopped test closures prevent a complete
qualification claim. The unchanged project remains 30/22/20/20 without a
go-diff, Blackfriday, fastuuid, TSDB, or Procfs trace. Client_golang v1.4.0
retains its separate GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698 identity and exception.

This exception does not transfer the Blackfriday v2.1.0, TSDB v0.7.1, Procfs
v0.0.8, Common v0.9.1, client_model v0.2.0, or client_golang v1.4.0
unqualified, non-transferable exceptions. It does not transfer, reopen, or
alter Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack,
fastuuid, or any earlier qualified or excepted result. Each remains governed
only by its own exact guards.

### Guard-only verification and handoff

- Work began clean on `codex/upgrade-quality` at go-diff evaluation handoff
  HEAD `fccb8361de745a870bc2602d198945d88c0dc54d`, parent
  `0710ac3a6d65daf77e28c06d52c44c93a16e3300`, tree
  `5a3131b89a34a739de166a1af0bbb75f4232edef`. It changes exactly the
  launcher, answered go-diff evaluation, this then-NEXT decision, rolling
  handover, and roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains ancestral at its exact
  parent and tree.
- The reciprocal 301-archive chain, sole NEXT state, launcher/archive prompt
  mirror and check, branch, exact changed set, and ordinary and ignored
  cleanliness reproduce. Completed release, source, behavior, API, closure,
  native, test-scope, cross, projection-candidate, and advisory-matrix work
  was not repeated.
- Fresh official Darwin arm64 SDKs reproduce Go 1.18.10 archive/binary
  SHA-256
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and Go 1.26.7
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- The unchanged project reproduces 234 modules, 3,599 graph edges, 355
  production entries, 429 complete-test entries, 197 module-backed entries
  across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`, and graph
  hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  The Go floor remains 1.18. Contained ordinary tidy reaches the exact common
  projection at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  with 234 modules and 3,557 edges. No projection is retained.
- The exact three requests/routes, requester imports and metadata boundary,
  why results, import/load/runtime relevance, first-root commits, 84-checkpoint
  history, and selected sums reproduce. The byte-exact graph preserves all 47
  pre-Goe selections/276 incoming edges and every separate later selection/
  request count through go-diff v1.2.0/three. Every earlier decision remains
  closed, separate, and untransferred; accepted quality remains 27/27 Q0-Q2
  PASS at L2.
- Fresh narrow target and project advisory guards reproduce. The Go module
  index remains 518,501 bytes/1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  the PUBLISHED CVE-2026-14362 CNA response remains 2,807 bytes at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under ordinary umask 022. Product
source, `go.mod`, and `go.sum` remain byte-exact. Every task-owned SDK, cache,
response, report, tool, and project copy was contained beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`; all 23,891 entries were verified beneath
the exact managed root with no symlink or outside-path entry, then removed
before handoff. The reciprocal archive graph now has 302 records and exactly
one NEXT successor.

P7 remains active with one prepared, unlaunched bounded evaluation of the next
unevaluated alphabetical module, exact
`github.com/shurcooL/sanitized_anchor_name v1.0.0`. No ownership study, other
dependency evaluation, or P8 work was run; P8 remains queued.
