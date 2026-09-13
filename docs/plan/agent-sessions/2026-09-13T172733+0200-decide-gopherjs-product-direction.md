# Agent Session: Decide GopherJS Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-13T172733+0200-decide-gopherjs-product-direction`
Created: `2026-09-13T17:27:33+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `45181d18e878ac425cb259204c50565c11ed47664d2cb20c8574d66c0c9db06d`
Previous: [2026-09-13T140912+0200-evaluate-gopherjs-dependency.md](2026-09-13T140912+0200-evaluate-gopherjs-dependency.md)
Next: [2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md](2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md)
Outcome: Recorded the user's bounded option 1 decision, retained the exact inherited and unloaded GopherJS pseudo-version without metadata changes, revalidated every exception guard, and prepared the next bounded P7 group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-13 explicit selection of
option 1 for exact-path
`github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1`. Retain the
exact selected, inherited, unloaded pseudo-version without dependency metadata
changes under the bounded release, closure, compiler, command, and runtime
exceptions below. Do not implement a dependency change, audit another group,
or begin P8 in this decision-recording move.

# Authorized Roadmap

P2A-P6 and all earlier P7 dependency decisions remain final. Existing
Enterprise Certificate Proxy and GAX exceptions remain separate under their
own zero-load guards and do not transfer to GopherJS.

# Authorized Product Decision

On 2026-09-13 the user explicitly selected option 1 with the recommended
bounds: retain exact selected
`gopherjs v0.0.0-20181017120253-0766667cb4d1` as an inherited, unloaded
selection. Accept only the already documented unqualified pseudo-release,
unsigned and untagged identity, proxy-synthesized module file, undeclared and
incomplete modern source/test closure, missing vendored test package,
post-floor modern resolution, and known compiler, command, and generated-
runtime incompatibilities. This does not accept a new or independently
discovered defect.

Keep `go.mod` and `go.sum` unchanged. Do not add a direct GopherJS edge, select
another target version, change mvn-pom-mutator or GoConvey, authorize a broader
parent/removal evaluation, raise the Go floor, or select a prerelease, fork,
replacement, patch, alternate module path, unrelated-module move, or dependency
implementation commit.

The exception is non-transferable and valid only while the complete project
load contains zero GopherJS packages, the module remains runtime-unreachable,
the exact selected pseudo-version and its sole incoming GoConvey v1.6.4 edge
remain unchanged, and no new advisory or independent disqualifier appears.
Revalidate and record those guards. Direct import or loading, runtime
reachability, a version or incoming-edge change, or a new advisory or
independent disqualifier expires the exception and requires a fresh dependency
and product decision before merge. Do not stop or ask for this same GopherJS
decision again while all guards hold.

# Measurements At Start

Fresh proxy, sumdb, go-import, strict Git, and GitHub evidence resolves the
public active non-fork BSD-2-Clause repository. The proxy lists thirteen
versions. Selected is unsigned 2018 commit
`0766667cb4d1cfb8d5fde1fe210ae41ead3cf589`, not a tagged release, and its
module file is proxy-synthesized. Its dependency/test closure is undeclared and
incomplete; modern resolution exceeds the Go floor, a vendored test package is
absent, and compilation deliberately fails outside Go 1.11. V1.12.80 requires
Go 1.12.

Stable v1.17.2 is the highest release whose complete imported closure preserves
Go 1.18. It is commit `fcf8e05a6f4fe7573b43c6bc65c2d1166fbd48cd`,
declares Go 1.17, and has a 157-module graph whose actual test closure uses 29
external packages across 14 modules. Under exact Go 1.18.10 host, contained Go
1.17.9 target, Node 12.22.12, and the native syscall addon, its compiler suite
reproducibly records 421 pass, 62 known failures, and seven unexpected
Darwin/arm64 generated-program crashes through unimplemented
`internal/abi.FuncPCABI0`. A trivial generated program independently crashes
at first output. Required Go 1.26.7 also exposes generic-AST and VERSION-parser
failures. Partial native/race/vet and deterministic Linux output passes do not
qualify it.

V1.18.0-beta3 preserves a Go 1.18 closure and passes its Go 1.18.10 suite, but
it is a prerelease and cannot link under Go 1.26.7 because its goembed closure
references private `go/build.parseGoEmbed`. Stable v1.20.2 and v1.21.0 declare
Go 1.20 and Go 1.21. Pinned apidiff records incompatible build/compiler API
changes across these lines. No stable, prerelease, or selected pseudo-version
passes the complete existing contract.

The unchanged graph path is main -> direct mvn-pom-mutator v0.2.3 -> GoConvey
v1.6.4 -> selected, with GoConvey supplying the sole target edge. `go mod why
-m` is negative. The complete 429-entry project test load contains 197
module-backed packages across 41 loaded modules with zero GopherJS packages,
zero Enterprise Certificate Proxy packages, and zero GAX packages. Project
selection remains 234 modules and 3,599 edges; `go.sum` remains 1,067 lines.
The `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Disposable exact v1.17.2/beta3 requests add target closure modules but retain
zero target loading and tidy byte-identically back to the base projection.
V1.20.2/v1.21.0 move unrelated `x/*` selections, and v1.21.0 raises the main
Go directive. A direct target edge has no selection or runtime purpose.

Fresh primary vulnerability data has 1,398 records and no exact GopherJS
record. Direct candidate package scans report inherited Logrus and `x/sys`
advisories but zero reachable symbol or test-symbol trace. Unchanged project
populations remain 30 module findings, 22 Darwin and 23 Windows package
findings, and 20 IDs/22 reachable traces on both symbol platforms, with no
GopherJS occurrence.

Exact Go 1.26.7 project verification/load/build, independent tests, race, vet,
pinned lint, empty-HOME, cross-builds, API/CLI compatibility, and full preflight
pass. Applicable Go 1.18.10 gates retain only the two accepted shell wording
differences. No dependency metadata or implementation commit was created;
accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Role And Boundaries

P2A-P6 and all earlier P7 dependency decisions are final. The 2026-09-13
Enterprise Certificate Proxy option 1 decision remains valid only for exact
v0.2.1 under its own zero-load/version/Viper-edge/advisory guards. The two GAX
v2.7.0 exceptions remain valid only under their own zero-load guard. Both
targets remain at zero loaded packages. Do not reopen those decisions or
transfer their exceptions to GopherJS.

This is a decision-recording move, not a renewed technical audit or dependency
implementation. The user has explicitly selected and bounded option 1; do not
ask for that decision again or broaden it into parent, removal, direct-use,
future-version, floor, prerelease, fork, replacement, patch, alternate-path, or
unrelated-module authorization.

# Required Reading

Read the answered GopherJS evaluation archive, rolling handover, roadmap, the
answered Enterprise Certificate Proxy decision/evaluation, and the answered
GAX archive. Verify the feature branch, reciprocal lifecycle graph, clean
ordinary and ignored status, current ancestry, latest dependency implementation
identity, unchanged module hashes, P7/P8 state, and
`./codex-dev-start.sh --check` before recording a decision.

Use exact Go 1.26.7 first in PATH with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, no ambient `GOFLAGS`, and `LC_ALL=C LANG=C`. Put every
disposable artifact beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never run
`go mod download all` in a measured worktree or create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Three Moves

First, revalidate only the decision guards: exact selected GopherJS and sole
GoConvey edge, negative why result, zero GopherJS/ECP/GAX loaded packages,
unchanged module hashes, and current primary advisory state. Reuse the completed
technical audit; do not repeat or broaden it.

Second, record the exact option 1 exception, accepted known findings, guards,
expiration triggers, and no-change result in the roadmap and rolling handover.
Retain the selected pseudo-version without a dependency or implementation
commit.

Third, answer this archive and follow the lifecycle contract for exactly one
next bounded P7 mission under the existing queue. Do not execute the successor
mission in this turn.

# Automatic Handoff

After recording the decision, follow the repository lifecycle contract. Do not
implement a dependency strategy, audit another dependency, launch a successor,
push, merge, publish, release, stash, revert, bypass cleanup, or remove the
worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Recorded the user's 2026-09-13 option 1 selection and retained exact-path
`github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1` as an
inherited, unloaded module without changing `go.mod` or `go.sum` and without
manufacturing a dependency or implementation commit.

The decision accepts only the findings already established in the answered
dependency evaluation: the unqualified pseudo-release; unsigned, untagged
identity; proxy-synthesized module file; undeclared and incomplete modern
source/test closure; absent vendored test package; post-floor modern
resolution; and the recorded compiler, command, and generated-runtime
incompatibilities. No new or independently discovered defect is accepted.

Every guard was revalidated from clean decision HEAD `ae2a9d3` with a freshly
unpacked exact Go 1.26.7 distribution whose archive and binary SHA-256 values
are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
It ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.
Project selection remains mvn-pom-mutator v0.2.3, GoConvey v1.6.4,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate
Proxy v0.2.1, and GAX v2.7.0. The unchanged graph path is main -> direct
mvn-pom-mutator v0.2.3 -> GoConvey v1.6.4 -> selected GopherJS, and the
GoConvey edge is the target's sole incoming edge. `go mod why -m` remains
negative, and repository Go source has zero direct target imports. Together
with the zero package load below, this revalidates that GopherJS remains
runtime-unreachable in the current project.

The complete project test load remains 429 entries, including 197 module-
backed packages across 41 loaded modules, with exactly zero GopherJS,
Enterprise Certificate Proxy, or GAX packages. The project still selects 234
modules with 3,599 graph edges. `go.mod` and `go.sum` remain byte-identical to
clean decision HEAD at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
`go.sum` remains 1,067 lines.

Fresh primary vulnerability data still contains 1,398 module records, was
last modified 2026-09-10T16:28:28Z, and contains no exact GopherJS record.
No new advisory or independently observed disqualifier arose during this
guard-only revalidation. The completed technical audit was not repeated or
broadened.

The non-transferable GopherJS exception remains valid only while the complete
project load contains zero target packages, the module stays runtime-
unreachable, exact selected
`v0.0.0-20181017120253-0766667cb4d1` and its sole incoming GoConvey v1.6.4
edge remain unchanged, and no new advisory or independent disqualifier
appears. Direct import or loading, runtime reachability, a version or incoming-
edge change, or a new advisory or independent disqualifier expires the
exception and requires a fresh dependency and product decision before merge.
The Enterprise Certificate Proxy and GAX exceptions remain separate under
their own zero-load guards. Do not request this same GopherJS decision again
while all guards hold.
