# Agent Session: Decide Sergi Go-Diff Product Direction

Status: NEXT
Session ID: `2026-09-24T011213+0200-decide-sergi-go-diff-product-direction`
Created: `2026-09-24T01:12:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f2816b5fd3f85b7fcd804f60be6fa8461da99002f75dfac0239556ee0fec152c`
Previous: [2026-09-24T002515+0200-evaluate-sergi-go-diff-dependency.md](2026-09-24T002515+0200-evaluate-sergi-go-diff-dependency.md)
Next: none
Outcome: pending

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
