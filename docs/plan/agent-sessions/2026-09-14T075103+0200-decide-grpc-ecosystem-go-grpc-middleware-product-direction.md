# Agent Session: Decide gRPC Middleware Product Direction

Status: NEXT
Session ID: `2026-09-14T075103+0200-decide-grpc-ecosystem-go-grpc-middleware-product-direction`
Created: `2026-09-14T07:51:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `08d30c32431325fa1511c059a5f0f3da291618b6356fa222af495b77c3c21003`
Previous: [2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md](2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-14 explicit selection of
option 1 for exact-path `github.com/grpc-ecosystem/go-grpc-middleware`.
Retain exact selected, inherited, unloaded v1.0.0 without dependency metadata
changes under the bounded release-closure, source, test, and retry exceptions
below. Do not implement a dependency change, reopen an earlier decision, audit
another dependency group, or begin P8 in this decision-recording move.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 is active only for this
gRPC middleware product decision; P8 remains queued. This session may record
one explicit user choice and prepare one bounded follow-up, but may not
implement that choice or combine another dependency group.

## Authorized Product Decision

On 2026-09-14 the user explicitly selected option 1 with the recommended
bounds: retain exact selected
`github.com/grpc-ecosystem/go-grpc-middleware v1.0.0` as an inherited,
unloaded selection without changing `go.mod` or `go.sum`. Accept only the
already documented missing release module metadata and deterministic complete
closure, TLS test-certificate failures in later candidates, universal retry
cancellation and discarded-timer-cancel resource defects, and related recorded
source/test qualification findings. This does not accept a new or independently
discovered defect.

The exception is target-specific and non-transferable. It is valid only while
exact v1.0.0 and its sole incoming mvn-pom-mutator v0.2.3 edge remain unchanged,
the complete project load contains zero middleware packages, the module remains
runtime-unreachable, and no new advisory or independent disqualifier appears.
Revalidate and record those guards. Direct import or loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent defect expires the exception and requires a fresh dependency and
product decision before merge.

Do not add a direct target edge, select v1.4.0 or another release, change
mvn-pom-mutator or GoConvey, raise the Go floor, authorize a patch/replacement
or parent/removal design, change unrelated selections, or manufacture a
dependency implementation commit. Unpatched v1.4.0 remains disqualified by
the universal retry defects, logger and repeat-suite defects, incompatible
logger API changes, and broad MVS effects. Existing Gorilla WebSocket,
GopherJS, Enterprise Certificate Proxy, and GAX exceptions remain separate.
Do not stop or ask for this same gRPC middleware decision again while all
guards hold.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves the
canonical public active unarchived non-fork Apache-2.0 repository. The exact
path has seven stable releases, v1.0.0 through v1.4.0, and no retraction or
formal module deprecation. V1.4.0 is proxy latest and the highest stable
release whose complete minimal source/test closure preserves Go 1.18, but it
is not qualified. Current main is the different
`github.com/grpc-ecosystem/go-grpc-middleware/v2` module requiring Go 1.24;
the unreleased v1 branch and v2 line are not exact-path candidates.

Every exact v1 release reproduces the same production cancellation defect on
both exact Go 1.26.7 and contained Go 1.18.10: with retry enabled and zero
backoff, an already-canceled parent still reaches the invoker three times and
returns Unavailable rather than Canceled. Every release source also discards
the cancel function from `context.WithTimeout` in per-call retry contexts;
v1.1.0-v1.4.0 fail modern vet on that resource leak. Selected v1.0.0 has no
release `go.mod` and cannot supply a complete deterministic source/test
closure. V1.1.0-v1.2.1 default tests fail because their CN-only certificate
has no SAN. V1.2.2 has a Go 1.26 server-restart state race. V1.4.0 fails both
SDKs' two full count-10 repeats when unclosed gRPC work calls a subtest-bound
global logger after the subtest ends, and `logging/settable` nests variadic
arguments instead of forwarding them. V1.3.0 passes repeats, race, and cross-
builds but retains the retry cancellation and timer defects.

Pinned apidiff records generated test-protobuf comparability and JSON
marshaller incompatibilities from selected to later releases. V1.3.0 to
v1.4.0 additionally changes every logging/kit public logger type from
`github.com/go-kit/kit/log` to `github.com/go-kit/log`. Positive independent
fixtures for unary/stream ordering, short-circuiting, context/metadata, auth,
rate limiting, recovery, validation, retry attempts/headers, stream limits,
tags, repeats, and supported race behavior pass under both SDKs; the two
production defects above fail identically.

Selected exists only through main -> direct mvn-pom-mutator v0.2.3 -> v1.0.0.
`go mod why -m` is negative, repository source imports are zero, and zero
target packages load. Exact v1.4.0 `go get` keeps target loading at zero but
changes target plus four unrelated existing selections, adds 116 graph
modules, and changes the project from 234/3,599/1,067 modules/edges/sum lines
to 350/3,784/1,084. Tidy returns byte-identically to the base tidy projection
and restores inherited v1.0.0. A manual target-only requirement would avoid
the four existing-version moves but violates the mandated exact-`go get`
workflow and is not authorized.

Fresh primary vulnerability data has 1,398 module records, SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
Last-Modified 2026-09-10T16:28:28Z, and no exact target record. Reviewed base
and v1.4.0 project scans are identical at 30 module findings, 22 package
findings, and 20 IDs/22 symbol traces, with zero target occurrence. The
isolated v1.4.0 old gRPC/x/net closure has inherited findings; none is an exact
middleware advisory.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod` and `go.sum` remain at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 full preflight, repeats, race, vet, pinned lint, hermetic,
cross-build, API, and CLI gates pass. The Go 1.18.10 projection retains only
the two accepted `pkg/shell` closed-file wording failures. Accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Role And Boundaries

This is a decision-recording turn, not a renewed technical audit or dependency
implementation. The user has explicitly selected and bounded option 1; do not
ask for that decision again, manufacture activity, or broaden it into direct
use, another version, parent/removal, patch/replacement, floor, API, or
unrelated-module authorization.

## Guarded Decisions

P2A-P6 and every earlier P7 outcome are final. Gorilla WebSocket v1.4.2,
GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy
v0.2.1, and GAX v2.7.0 retain only their separate documented exceptions under
their exact version, incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards before recording a choice; stop for
the owning decision if any expired. Do not transfer their exceptions to gRPC
middleware, change mvn-pom-mutator or GoConvey, or combine another group.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, reciprocal archive history, latest dependency implementation, all
guarded invariants, and `./codex-dev-start.sh --check`. Read this archive, the
answered gRPC middleware evaluation, rolling handover, roadmap, `go.mod`,
`go.sum`, and referenced contracts. This is a decision-recording turn; reuse
the completed technical audit and do not ask the user to choose again.

# Three Moves

First, revalidate only the unchanged selection, sole incoming edge, negative
why result, zero guarded package loads, runtime unreachability, module hashes,
and current advisory state. Reuse the completed audit; do not repeat or broaden
it. Second, record the exact option 1 exception, accepted known findings,
guards, expiration triggers, and no-change result in the roadmap and rolling
handover. Third, answer this archive and prepare one reciprocal NEXT mission
for the next bounded P7 group without executing it.

# Automatic Handoff

After recording the decision, run the applicable no-change lifecycle gates and
create the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
