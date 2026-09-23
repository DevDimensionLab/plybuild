# Agent Session: Decide Russross Blackfriday V2 Product Direction

Status: NEXT
Session ID: `2026-09-23T234834+0200-decide-russross-blackfriday-v2-product-direction`
Created: `2026-09-23T23:48:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8d49a85e173c8b6ce37de6073998e7aafe8c10361e03913df61303f2dd96a8a6`
Previous: [2026-09-23T224910+0200-evaluate-russross-blackfriday-v2-dependency.md](2026-09-23T224910+0200-evaluate-russross-blackfriday-v2-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product-direction decision for exact
selected indirect `github.com/russross/blackfriday/v2 v2.1.0`. The completed
fresh evaluation found no fully qualified canonical Go-1.18-compatible stable,
although v2.1.0 is highest compatible and has a genuine supported tidy-stable
requester. Choose and record exactly one of the three authorized directions
below. Do not repeat the evaluation, implement a dependency or source change,
combine another group, launch a study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only completed
public release/source/build/graph/projection/advisory evidence and bounded
read-only guard checks. Do not fuzz, stress, probe resource exhaustion, create
oversized, deeply nested, cyclic, malformed, adversarial, or escape-sequence
payloads, reproduce a security issue, or perform security or exploitability
analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Set every
tool temp/cache root explicitly, verify containment, and remove task-owned
scratch evidence before handoff.

# Completed Evaluation

Parent-path primary go-import metadata maps the exact owner to the public,
enabled, unarchived, non-fork `russross/blackfriday` repository. Its default
`master` branch remains on the separate v1.6.0 line; exact branch `v2` is at
selected v2.1.0. The repository license file contains Simplified BSD
two-clause terms despite GitHub's unclassified license result.

The exact-path proxy exposes v2.0.0, v2.0.1, v2.1.0-pre.1, and v2.1.0.
V2.0.0 is not a usable canonical `/v2` release because its tree lacks an
applicable `go.mod`; v2.1.0-pre.1 is a prerelease at the same commit as the
stable. Thus v2.0.1 and v2.1.0 are the two canonical fetchable stables, with
v2.1.0 highest. No replacement, retraction, deprecation, or alternate-major
exact-path line exists. Neither module declares a Go version, so both preserve
the project Go 1.18 floor. Exact tags, commits, trees, signatures, ancestry,
proxy/sumdb/archive-to-Git identity, module/license files, packages, source
boundaries, and regular-only archive entries are recorded in the answered
evaluation.

V2.0.1 has an incomplete declared standalone closure: production imports
`github.com/shurcooL/sanitized_anchor_name`, tests also import go-difflib, and
its `go.mod` declares neither. Native readonly build/vet/test compilation and
all 40 cgo-disabled production/test compilation results therefore fail across
both exact SDKs and ten supported targets.

V2.1.0 has a one-module, one-package standard-library closure. Under exact Go
1.18.10 and Go 1.26.7 it passes module verification, build, vet, safe
test-package compilation, and all 40 cgo-disabled production/test-compilation
rows across Darwin, Linux, Windows, FreeBSD, Plan 9, and js/wasm. These are
partial qualification results only.

Both stable ordinary test suites invoke a helper explicitly designed to test
every input substring to stress bounds checking and contain a test explicitly
using malformed input. Static review found the boundary before execution, so
complete count-one, repeated, and race tests were stopped and not run under
either SDK. The governing defensive scope forbids those payloads and forbids
using partial evidence as qualification. Consequently no canonical stable
fully qualifies; no advisory absence or passing compile row changes that
result.

The current graph has exactly four v2.1.0 requests: main and go-md2man/v2
v2.0.1, v2.0.6, and selected v2.0.7. Each go-md2man release genuinely imports
Blackfriday in production and tests and requests exactly v2.1.0. Selected
Cobra v1.10.2 genuinely imports go-md2man in its production documentation
package. The public go-md2man repository is enabled, unarchived, non-fork, and
active, establishing a genuine supported tidy-stable transitive requester
boundary even though the project does not load it.

The four current complete routes are the direct main edge; main -> selected
go-md2man v2.0.7; main -> Cobra v1.10.2 -> go-md2man v2.0.6; and main ->
direct mvn-pom-mutator v0.2.3 -> Cobra v1.4.0 -> go-md2man v2.0.1. All 84
historical module checkpoints reproduce the route evolution from the initial
2020 Cobra v1.0.0/go-md2man v2.0.0/Blackfriday v2.0.1 route through the
current v2.1.0 route families. Main first requests target v2.1.0 and
go-md2man v2.0.7 at commit `55dc69dda20c4d8f96b6dbcdd70cfef76467cc85`;
neither is an earlier root.

Target and requester why are negative; repository imports, project production
and complete-test loads, module-backed entries, and runtime relevance are all
zero for both. The unchanged project populations remain 355/429/197 across 41
loaded modules. Exact selected get changes no byte or selection. Ordinary tidy
removes the redundant direct target/requester roots while retaining v2.1.0
through genuine go-md2man routes and restores the established common
projection. No projection was retained.

Exact-version OSV and narrow GitHub/repository responses are empty without
implying safety or qualification. Pinned govulncheck v1.8.0 reports 0/0/0/0
for v2.1.0's complete focal closure; v2.0.1 cannot be loaded completely and
has no usable empty result. The project remains 30/22/20/20 without a
Blackfriday, fastuuid, TSDB, or Procfs trace; client_golang v1.4.0 retains
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. No exploitability claim
was made.

The real project remains exactly 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, and 1,067 sum lines at the protected hashes and Go 1.18
floor. Official SDK identities, the established common tidy projection,
corrected Go-index/CNA identities, all 47 pre-Goe selections and 276 incoming
edges, every separate later selection/request count and exception, all earlier
decisions, and accepted 27/27 Q0-Q2 PASS at L2 remain exact. No dependency,
source, root, projection, exception, or study was retained by the evaluation.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected, indirect, unloaded
   `github.com/russross/blackfriday/v2 v2.1.0` unchanged under a
   Blackfriday-specific, unqualified, non-transferable exception. State that it
   is highest compatible but is not qualified, safe, or fixed. Preserve the
   exact four requests and current routes, all historical route families,
   genuine go-md2man requester/import boundary, negative why, zero project
   import/load/runtime facts, root history, sums, release/source/API/behavior/
   closure/native/test-scope/cross/projection/advisory identities, and every
   earlier guard as expiry conditions. Do not transfer another exception.
2. Authorize exactly one later measurement-only **Go-Md2man Blackfriday
   Elimination Study** to determine whether a supported requester update or
   replacement can eliminate Blackfriday while preserving product behavior,
   direct roots, the Go 1.18 floor, and every earlier guard. Do not run the
   study, grant an exception, change a dependency or source file, or
   pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the selected highest compatible release is
unloaded and runtime-irrelevant, the genuine transitive requester is supported,
and the qualification boundary is explicit. This is a product acceptance
decision, not a qualification or safety claim.

# Required Reading And Handoff

Start only from the clean Blackfriday evaluation handoff. Verify its HEAD,
parent, tree, exact changed set, branch and Google UUID ancestry, reciprocal
archive chain, sole NEXT state, launcher mirror/check, ordinary and ignored
cleanliness, official SDK identities, real project counts and hashes, common
tidy projection, Go 1.18 floor, all target/requester request/import/route/
relevance/root facts, target/closure/project advisory identities, and every
earlier guard. Stop for a fresh owning decision if any protected input changed.

Preserve exact selected TSDB v0.7.1, Procfs v0.0.8, Common v0.9.1,
client_model v0.2.0, and client_golang v1.4.0 only under their own separate
unqualified, non-transferable exceptions. Preserve Complete, go-difflib, SFTP,
pkg/errors, Goe, ULID, go-conntrack, and every earlier qualified or excepted
result under its exact guards. Preserve fully qualified selected fastuuid
v1.2.0 without adding a root. Do not run rejected studies, select rejected
candidates, reopen a completed module, or transfer an exception. P8 remains
queued.

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
P8.
<!-- CODEX_SESSION_PROMPT_END -->
