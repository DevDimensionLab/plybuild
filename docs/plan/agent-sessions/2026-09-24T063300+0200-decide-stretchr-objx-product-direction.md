# Agent Session: Decide Stretchr Objx Product Direction

Status: NEXT
Session ID: `2026-09-24T063300+0200-decide-stretchr-objx-product-direction`
Created: `2026-09-24T06:33:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ecbee243021abd3d216927f0132a061a15d690485e494a513a9cf0882b106e12`
Previous: [2026-09-24T055739+0200-evaluate-stretchr-objx-dependency.md](2026-09-24T055739+0200-evaluate-stretchr-objx-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one explicit product-direction decision for exact
selected indirect `github.com/stretchr/objx v0.5.2`. Choose one of the three
authorized moves below. This is a decision-only session: do not rerun the
completed evaluation, execute a study, change a dependency, add a root,
combine another dependency group, relax the Go 1.18 floor, or begin P8.

# Closed Evaluation

The fresh bounded Objx evaluation is complete. Exact go-import metadata maps
the path to the public, enabled, unarchived, non-fork MIT
`stretchr/objx` repository on `master`. The canonical exact-path line has
exactly nine stable releases v0.1.0, v0.1.1, v0.2.0, v0.3.0, v0.4.0, v0.5.0,
v0.5.1, selected v0.5.2, and latest v0.5.3, with no prerelease, replacement,
retraction, deprecation, alternate-major line, or non-default branch. All
tag/commit/tree/signature/ancestry, GitHub Release, proxy/sumdb/archive-to-Git,
module/license, package, and source/build identities are resolved. V0.5.1 is
the highest canonical stable compatible with Go 1.18; v0.5.2 and v0.5.3 both
declare Go 1.20.

No compatible stable fully qualifies. Every compatible release v0.1.0-v0.5.1
has a mandatory upstream test that intentionally supplies malformed JSON;
v0.1.1-v0.5.1 also exercise an invalid percent escape. Complete count-one,
repeated, and race rows were therefore STOPPED rather than violate the
defensive scope. V0.1.0-v0.1.1 separately fail test compilation and vet
because synthesized module metadata omits their Testify test dependency.
V0.2.0-v0.5.1 pass safe module verification, production build, no-run native
test compilation, vet, and every supported cgo-disabled production/test-
compilation row under exact Go 1.18.10 and Go 1.26.7. Passing safe evidence is
partial and cannot replace the stopped mandatory rows or qualify a release.

Objx's API exposes map/value wrappers, parsing and signed encoding helpers,
selector access and mutation, typed conversions, collection transforms, and
map helpers. Constructors and accessors may alias caller maps/slices; `Set`,
`MergeHere`, and JSON cleaning can mutate receiver backing data; shallow copies
can retain nested aliases; map callback order is nondeterministic; and callers
must synchronize shared maps and slices. Global URL-slice-suffix mutation is
unsynchronized. There is no retained external-resource lifecycle. Error forms
return parse/signature errors, `Must` forms panic, and invalid selectors may
produce a missing value or no-op. Exact API/behavior/caller/concurrency/
lifecycle/error, build-tag/cgo/generate/embed/codegen, symlink/submodule, and
archive-entry facts are closed guards.

The current graph has exactly 11 Objx requests. Logrus v1.2.0 and v1.4.2
request v0.1.1 without importing Objx, so those edges are metadata-only.
Testify v1.3.0-v1.7.1 request v0.1.0, v1.8.0 requests v0.4.0, v1.8.4 requests
v0.5.0, and selected v1.9.0 requests selected v0.5.2. Every listed Testify
release genuinely imports Objx in production `mock/mock.go`. The selected
route is main -> allecthomas/units `0f3dac36c52b` -> Testify v1.9.0 -> Objx
v0.5.2. Objx why is negative; Testify why is positive only through a
go-term-markdown dependency-test route. Neither is imported by project source
or present in production, complete-test, or module-backed load populations.
Objx is indirect, unloaded, runtime-irrelevant, and not a main root; Testify is
also not a root.

All 84 distinct historical module checkpoints graph successfully. Objx and
Testify are never roots. Selection is v0.1.1 in 38 checkpoints, v0.4.0 in four,
v0.5.0 in 25, and v0.5.2 in 17. All 19 route epochs end at the same genuine-
Testify-production-import or metadata-only-Logrus boundary. No current or
historical requester owns v0.5.1.

A disposable v0.5.1 get manufactures a main target root and forces unrelated
downgrades, including Testify v1.9.0 -> v1.8.4, Units, and go-toml/v2. Ordinary
tidy removes v0.5.1 and selects v0.5.0 through Testify v1.8.4 while retaining
an unrelated go-toml downgrade. It does not restore the no-op baseline and is
not an authorized dependency-only selection. No genuine supported tidy-stable
requester owns v0.5.1, and no projection is retained. Normal no-op tidy still
reproduces the exact protected common projection.

Exact-version OSV results for all nine releases, narrow GitHub global results,
and the repository advisory endpoint are empty. Pinned govulncheck v1.8.0
focal results are 0/0/0/0 for v0.5.1 and v0.5.2. Advisory absence is not a
safety or qualification claim. The corrected index/CNA identities, project
30/22/20/20 population without an Objx or protected-target trace, and the
client_golang GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698 identity
remain exact.

# Protected State And Closed Guards

Start only from the clean Objx evaluation handoff on
`codex/upgrade-quality`. Verify its HEAD, parent, tree, exact changed set,
branch and Google UUID ancestry, reciprocal archive chain, sole NEXT state,
launcher mirror/check, and ordinary and ignored cleanliness. Stop for a fresh
owning decision if a protected input changed. Do not rerun completed release,
source, behavior, closure, test, projection, history, or network research when
guard-only evidence remains exact.

The unchanged real project remains exactly 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
over 41 loaded modules, and 1,067 `go.sum` lines. Protected `go.mod`, `go.sum`,
graph, normal-tidy diff, common 52/948-line 234-module/3,557-edge projection,
Go 1.18 floor, 27/27 Q0-Q2 PASS at L2, and exact Go 1.18.10/Go 1.26.7 archive
and binary identities remain exact. No projection is retained.

Preserve all 47 pre-Goe selections and 276 incoming edges; every separate
selected version and request count; every qualified result; and every target-
specific, unqualified, non-transferable exception through Murmur3.
Rogpeppe/go-internal remains closed as Cast v1.5.1's minimal test closure.
Do not run a rejected study, reopen a completed module, select a rejected
candidate, transfer an exception, or relax the Go floor.

Preserve Objx's exact owner/repository/release chronology, Go floors, GitHub
Releases, tag/commit/tree/signature/ancestry, proxy/sumdb/archive-to-Git,
module/license, package/source/build, selected and candidate sums, public API/
behavior/caller ownership and mutation, determinism/concurrency/lifecycle/
cleanup/error, build-tag/cgo/generate/embed/codegen, symlink/submodule, closure,
native vet/test/race/cross, projection, and advisory identities. Preserve all
current requests and route/import boundaries; 84-checkpoint history; target/
requester why, imports, loads, ownership, runtime, and root facts; and every
earlier guard.

The corrected Go vulnerability index remains exactly 518,501 bytes and 1,402
records at its protected SHA-256. The PUBLISHED CNA response remains exactly
2,807 bytes at its protected SHA-256. The project advisory guard remains
30/22/20/20 without an Objx or named protected trace, while client_golang
retains its recorded identifiers. Empty advisory evidence cannot establish
qualification or safety.

# Authorized Product Moves

Choose exactly one:

1. **Retain selected v0.5.2 under an Objx-specific exception
   (recommended).** Keep exact selected indirect
   `github.com/stretchr/objx v0.5.2` unchanged under a new unqualified,
   non-transferable exception. State precisely that v0.5.1 is the highest
   Go-1.18-compatible stable but no compatible stable completes the mandatory
   qualification gates, while selected v0.5.2 declares Go 1.20 and is not a
   qualified Go-1.18 release. The bounded acceptance may rely on the exact
   genuine Testify production-import owner, unloaded/runtime-irrelevant target,
   unchanged history, passing safe evidence, current successful exact-Go-1.18
   consumption, and advisory boundary, but is neither release qualification
   nor a safety claim. Expire it on any owner, request, route, import/load/
   runtime, root, source/API/behavior, declared or effective Go floor, closure,
   defensive-scope, vet/test/race, projection, or advisory change. Make no
   dependency, source, or root change.
2. **Authorize a later Testify Objx Ownership Study.** Define one separate,
   measurement-only study that may determine whether a supported Testify line
   can genuinely own a fully qualified Go-1.18-compatible Objx selection while
   preserving every project guard. Do not run the study now, pre-authorize a
   dependency change, add an Objx root, combine Testify with this group, waive
   the stopped mandatory rows or Go floor, or imply that v0.5.1, v0.5.2, or
   v0.5.3 qualifies. P7 remains blocked on the later study and another explicit
   owning decision.
3. **Stop P7 unresolved.** Make no dependency, source, root, exception, study,
   other-group, or P8 change. Record Objx as the unresolved blocker.

# Decision Contract

Make exactly one choice. Preserve selected v0.5.2 and every earlier guard
unless option 1 explicitly grants only the target-specific retention above.
Do not transfer another exception. A study authorization is not permission to
run it or to change Testify, Objx, Units, go-term-markdown, or another module.

Answer this archive, update the roadmap and rolling handover, verify contained
scratch cleanup, run final exact-Go-1.26.7 project module verification, build,
count-one tests, race count-one tests, and vet, and make one local handoff
commit. Do not rerun the completed Objx evaluation unless a guard-only check
proves a protected input changed.

# Required Reading

Read the answered Objx evaluation and all incorporated Murmur3, Cmux,
Goconvey, Assertions, sanitized_anchor_name, go-diff, Blackfriday, fastuuid,
TSDB, and earlier guards before deciding. Do not infer qualification from a
physical selection, genuine test-dependency ownership, passing safe rows,
successful old-toolchain consumption, or empty advisory evidence.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change a dependency or product
source, add a target root, transfer or reopen an exception, combine another
dependency group, or begin P8. If option 1 closes Objx, prepare but do not
launch exactly one next bounded P7 queue item. If option 2 is chosen, prepare
but do not launch exactly the named study. If option 3 is chosen, prepare no
successor.
<!-- CODEX_SESSION_PROMPT_END -->
