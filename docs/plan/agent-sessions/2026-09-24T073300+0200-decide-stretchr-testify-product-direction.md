# Agent Session: Decide Stretchr Testify Product Direction

Status: NEXT
Session ID: `2026-09-24T073300+0200-decide-stretchr-testify-product-direction`
Created: `2026-09-24T07:33:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3a1023f02f623717d9d4d3c213a6700f6fc05382314db4590e3853e5712b1da7`
Previous: [2026-09-24T064425+0200-evaluate-stretchr-testify-dependency.md](2026-09-24T064425+0200-evaluate-stretchr-testify-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one explicit product-direction decision for exact
selected indirect `github.com/stretchr/testify v1.9.0`. Choose one of the
three authorized moves below. This is a decision-only session: do not rerun
the completed evaluation, execute a study, change a dependency, add a root,
combine another dependency group, relax the Go 1.18 floor, or begin P8.

# Closed Evaluation

The fresh bounded Testify evaluation is complete. Exact go-import metadata
maps the path to the public, enabled, unarchived, non-fork MIT
`stretchr/testify` repository on `master`. The proxy has 31 stable artifacts
and no prerelease; v1.2.3 is a misspelled-path, self-retracted artifact and not
a canonical exact-path release, leaving 30 canonical v1 stables. There is no
alternate major, replacement, canonical retraction, or deprecation. All
release chronology, GitHub Release, tag/commit/tree/signature/ancestry,
proxy/sumdb/archive-to-Git, module/license, package, and source/build identities
are resolved.

No canonical stable fully qualifies. V1.12.1 is the highest stable whose
target module declares a Go version compatible with Go 1.18, but its complete
closure reaches Go 1.20 through Objx v0.5.3. V1.8.2 is the highest stable whose
complete closure remains within Go 1.18. Every release's complete suite
necessarily runs assert and require tests that intentionally provide invalid
JSON. Count-one, repeated, and race rows were therefore STOPPED rather than
cross the defensive scope. Safe verification, build, no-run test compilation,
vet, and supported cgo-disabled cross-compilation evidence is partial and
cannot replace stopped mandatory rows or qualify a release.

Testify exposes assert/require helpers, mock expectations and callbacks, suite
lifecycle hooks, and HTTP/filesystem/JSON/YAML/equality/error/collection/timing
helpers. Caller maps, slices, interfaces, pointers, callbacks, channels, test
objects, and suite instances remain caller-owned; mock arguments can alias,
callbacks can mutate, pointer changes affect comparisons, and callers must
synchronize shared state. Require uses `FailNow` on the test goroutine;
timing/channel helpers can block and are scheduling-dependent. Exact API,
behavior, ownership/mutation, determinism, concurrency, lifecycle/cleanup,
error, build-tag/cgo/generate/embed/unsafe, symlink/submodule, and archive facts
are closed guards.

The current graph has 51 requests: 43 test-file imports, three Logrus internal
test-support imports, and five metadata-only edges. Selected Units
`0f3dac36c52b` and go-toml/v2 v2.2.2 request selected v1.9.0 only for tests.
The selected shortest route is main -> Units -> Testify v1.9.0; target why is
positive only through go-term-markdown's dependency tests. Testify is not
imported by project source, is absent from production/complete-test/module-
backed load populations, is indirect, unloaded, runtime-irrelevant, and not a
root. It genuinely imports excepted Objx v0.5.2 in production `mock/mock.go`,
but that does not establish project ownership or qualification.

All 84 historical checkpoints graph successfully and Testify is never a root.
Selections are v1.3.0/1, v1.4.0/29, v1.7.0/7, v1.7.1/1, v1.8.0/3,
v1.8.1/26, and v1.9.0/17. Eighteen shortest routes span 20 consecutive
epochs. The complete historical union has 94 requester-version edges: 82
test-file imports, six Logrus internal test-support imports, and six metadata-
only edges. Every selected route ends at a test boundary and no requester owns
v1.12.1.

A disposable v1.12.1 get manufactures a main Testify root. Ordinary tidy
retains it and selects Testify v1.12.1, Objx v0.5.3, and go.yaml/v3 v3.0.5.
That is three changed selections, invalidates the exact Objx exception, and
has no genuine supported owner. No projection is retained; normal no-op tidy
reproduces the protected common projection.

Exact-version OSV, narrow GitHub, and pinned selected/candidate focal
govulncheck results are empty. Advisory absence does not imply safety or
qualification. The corrected index/CNA identities, project 30/22/20/20
population without Testify/Objx/protected traces, and client_golang advisory
identity remain exact.

# Protected State And Closed Guards

Start only from the clean Testify evaluation handoff on
`codex/upgrade-quality`. Verify its HEAD, parent, tree, exact changed set,
branch and Google UUID ancestry, reciprocal archive chain, sole NEXT state,
launcher mirror/check, and ordinary/ignored cleanliness. Stop for a fresh
owning decision if a protected input changed. Do not rerun completed release,
source, behavior, closure, test, projection, history, or network research when
guard-only evidence remains exact.

The unchanged project remains exactly 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, and 1,067 `go.sum` lines. Protected `go.mod`, `go.sum`,
graph, normal-tidy diff, common 52/948-line 234-module/3,557-edge projection,
Go 1.18 floor, 27/27 Q0-Q2 PASS at L2, and exact Go 1.18.10/Go 1.26.7 archive
and binary identities remain exact. No projection is retained.

Preserve every earlier qualified result and every target-specific,
unqualified, non-transferable exception through Objx. Objx v0.5.2 retains only
its exact genuine-Testify production-import owner and all recorded expiry
guards. Do not run the rejected **Testify Objx Ownership Study**, transfer an
exception, reopen a completed module, select a rejected candidate, or relax
the Go floor.

Preserve Testify's exact owner/repository/release chronology; Go floors;
GitHub Releases; tag/commit/tree/signature/ancestry; proxy/sumdb/archive-to-
Git; module/license/package/source identities; selected/candidate sums; public
API and behavior; caller ownership/mutation, determinism, concurrency,
lifecycle/cleanup/error boundaries; build-tag/cgo/generate/embed/unsafe facts;
closure/native/vet/test/race/cross results; projection effects; advisories;
all 51 current and 94 historical-union requests and boundaries; all 84
checkpoints, selection counts, routes/epochs; why/import/load/runtime/root
facts; and every earlier guard.

# Authorized Product Moves

Choose exactly one:

1. **Retain selected v1.9.0 under a Testify-specific exception
   (recommended).** Keep exact selected indirect
   `github.com/stretchr/testify v1.9.0` unchanged under a new unqualified,
   non-transferable exception. State precisely that no canonical stable fully
   qualifies, selected v1.9.0's closure reaches Go 1.20 through separately
   excepted Objx v0.5.2, and v1.12.1 is not an authorized candidate. The
   bounded acceptance may rely only on the exact test-only requesters/why
   route, unloaded/runtime-irrelevant target, unchanged history, passing safe
   evidence, current successful exact-Go-1.18 consumption, genuine selected
   Objx import, and advisory boundary. It is neither release qualification nor
   a safety claim. Expire it on any owner, request, route, import/load/runtime,
   root, source/API/behavior, Go-floor, closure, defensive-scope, vet/test/race,
   projection, Objx-exception, or advisory change. Make no dependency, source,
   or root change.
2. **Authorize a later Units/Go-TOML Testify Ownership Study.** Define one
   separate, measurement-only study that may determine whether supported
   genuine test owners can own a fully qualified Go-1.18-compatible Testify
   selection while preserving Objx and every project guard. Do not run it now,
   pre-authorize a dependency change, add a Testify root, waive stopped rows or
   a Go floor, transfer the Objx exception, or imply any release qualifies.
   P7 remains blocked on that later study and another explicit owning decision.
3. **Stop P7 unresolved.** Make no dependency, source, root, exception, study,
   other-group, or P8 change. Record Testify as the unresolved blocker.

# Decision Contract

Make exactly one choice. Preserve selected v1.9.0 and every earlier guard
unless option 1 explicitly grants only the target-specific retention above.
Do not transfer another exception. A study authorization is not permission to
run it or to change Testify, Objx, Units, go-toml/v2, go-term-markdown, or
another module.

Answer this archive, update the roadmap and rolling handover, verify contained
scratch cleanup, run final exact-Go-1.26.7 project module verification, build,
count-one tests, race count-one tests, and vet, and make one local handoff
commit. Do not rerun the completed Testify evaluation unless a guard-only
check proves a protected input changed.

# Required Reading

Read the answered Testify evaluation and incorporated Objx, Murmur3, Cmux,
Goconvey, Assertions, sanitized_anchor_name, go-diff, Blackfriday, fastuuid,
TSDB, and earlier guards before deciding. Do not infer qualification from a
physical selection, test-only ownership, passing safe rows, successful old-
toolchain consumption, or empty advisory evidence.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change a dependency or product
source, add a target root, transfer or reopen an exception, combine another
dependency group, or begin P8. If option 1 closes Testify, prepare but do not
launch exactly one next bounded P7 queue item. If option 2 is chosen, prepare
but do not launch exactly the named study. If option 3 is chosen, prepare no
successor.
<!-- CODEX_SESSION_PROMPT_END -->
