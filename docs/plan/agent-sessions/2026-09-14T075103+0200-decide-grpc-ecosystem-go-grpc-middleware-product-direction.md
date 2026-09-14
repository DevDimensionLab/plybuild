# Agent Session: Decide gRPC Middleware Product Direction

Status: NEXT
Session ID: `2026-09-14T075103+0200-decide-grpc-ecosystem-go-grpc-middleware-product-direction`
Created: `2026-09-14T07:51:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f04ddea44f6ea5a20c99170f0eb070732e64c9b1e20f2b144d3c191e09dde8e5`
Previous: [2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md](2026-09-14T062852+0200-evaluate-grpc-ecosystem-go-grpc-middleware-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and recording a fresh bounded product decision
for exact-path `github.com/grpc-ecosystem/go-grpc-middleware`. The completed
evaluation found no exact-path stable release that passes every existing
source, test, vet, behavior, API, and MVS contract. Present the three bounded
options below and ask the user to choose one. Do not infer a choice, implement
a dependency change, reopen an earlier decision, audit another dependency
group, or begin P8.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 is active only for this
gRPC middleware product decision; P8 remains queued. This session may record
one explicit user choice and prepare one bounded follow-up, but may not
implement that choice or combine another dependency group.

## Decision Required

Option 1 (recommended): retain exact selected, inherited, unloaded v1.0.0
without changing `go.mod` or `go.sum`, under a new target-specific exception.
Accept only the already documented release-closure, TLS-test, retry
cancellation/timer-resource, and related source/test qualification findings.
The exception must remain non-transferable and valid only while the exact
version and sole mvn-pom-mutator v0.2.3 edge remain unchanged, zero target
packages load, the module remains runtime-unreachable, and no new advisory or
independent disqualifier appears. Direct import/loading, runtime reachability,
a version or incoming-edge change, or a new advisory/defect expires it.

Option 2: authorize a separate bounded patch/replacement design. That new
scope must choose an exact source base, preserve Go 1.18, repair canceled-
context retry behavior and context timer cleanup, and, if based on v1.3.0 or
v1.4.0, repair settable logger variadic forwarding and relevant test lifecycle
defects. It must explicitly own fork/replacement maintenance, repository and
release identity, API compatibility, vulnerability, and full project gates.
Do not begin that design in this decision turn.

Option 3: authorize a separate parent/removal redesign centered on the sole
incoming `github.com/devdimensionlab/mvn-pom-mutator v0.2.3` edge, with any
parent, GoConvey, or related graph effects treated as a new bounded dependency
group. Do not change the parent or graph in this decision turn.

Do not offer an unpatched v1.4.0 upgrade as a qualified resolution: it retains
the universal retry defects, adds the settable logger defect and repeat-suite
lifecycle panic, introduces incompatible go-kit logger API changes, and exact
`go get` changes four unrelated existing selections plus 116 graph modules.

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

This is a decision turn, not a renewed technical audit or dependency
implementation. Present the three bounded options exactly enough for the user
to choose; do not infer a choice, manufacture activity, or broaden the selected
option beyond its documented consequences.

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
`go.sum`, and referenced contracts. This is a product-choice turn: present the
options and stop until the user chooses.

# Three Moves

First, revalidate only the unchanged selection, sole incoming edge, negative
why result, zero guarded package loads, runtime unreachability, module hashes,
and current advisory state. Reuse the completed audit; do not repeat or broaden
it. Second, present the three options and stop until the user makes an explicit
choice. Third, after that choice arrives, record only the choice and its exact
bounds, answer this archive, update the roadmap and rolling handover, and
prepare one reciprocal NEXT mission without executing it.

# Automatic Handoff

After an explicit choice, run the applicable no-change lifecycle gates and
create the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
