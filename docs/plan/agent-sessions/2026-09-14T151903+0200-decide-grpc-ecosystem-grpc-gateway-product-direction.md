# Agent Session: Decide gRPC Gateway Product Direction

Status: NEXT
Session ID: `2026-09-14T151903+0200-decide-grpc-ecosystem-grpc-gateway-product-direction`
Created: `2026-09-14T15:19:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fc7639d69aab6b0322e914f805364dcaf73d2a79dceb6835f50333c8683392fa`
Previous: [2026-09-14T134759+0200-evaluate-grpc-ecosystem-grpc-gateway-dependency.md](2026-09-14T134759+0200-evaluate-grpc-ecosystem-grpc-gateway-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by resolving the fresh bounded product decision for selected
exact-path `github.com/grpc-ecosystem/grpc-gateway v1.16.0`. No exact-path
stable release satisfies the existing qualification contracts. Present the
three bounded options below, obtain an explicit user choice, and stop without
implementing that choice. Do not evaluate another dependency group or begin
P8.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 remains active only for this
Gateway product choice; P8 remains queued. This session may request one
explicit choice but may not record or implement it, change another dependency,
or prepare a successor.

# Decision Required

Choose exactly one direction:

1. Retain exact inherited, unloaded v1.16.0 without dependency metadata
   changes under a new target-specific exception. The exception would accept
   only the documented universal discarded-cancel resource defect, native
   repeated-suite global/test-state failures, eight inherited reachable
   vulnerability IDs in the release-native closure, malformed-generator
   diagnostic nondeterminism, and related recorded qualification findings.
   This preserves the present graph and Go 1.18 floor but deliberately accepts
   an unqualified dependency while it remains unloaded and unreachable.
2. Authorize a separate bounded parent-ownership decision to remove or change
   the two incoming edges from etcd/api/v3 v3.5.1 and OTLP v0.7.0. This does
   not authorize either parent change now; it opens an independent graph and
   product evaluation because either move may alter unrelated selections or
   behavior.
3. Authorize a separate architecture and toolchain-floor decision for the
   different `github.com/grpc-ecosystem/grpc-gateway/v2` module path. This is
   not an in-place upgrade: current `/v2` main requires Go 1.26 and conflicts
   with the preserved Go 1.18 floor. No migration or floor change is authorized
   in this decision turn.

Recommend option 1 only if the user explicitly accepts its narrow exception
and guards. Otherwise recommend option 2 as the route to eliminate the unused
v1 selection without silently adopting a different major module path. Do not
infer a choice from prior retained-module decisions; no earlier exception
transfers to Gateway.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
public active unarchived non-fork BSD-3-Clause repository
`https://github.com/grpc-ecosystem/grpc-gateway.git`. The proxy exposes 55
versions: 53 stable and two prereleases. V1.16.0 is the latest exact-path
stable release and is lightweight tag/commit
`094a6fe78b3ca888297d090185cdf30f0e42e157`, parent
`e0a026aeb20c3eb0f32cbfdc537c9966634895b9`, tree
`2283b306f85e2b97858d8575954fa0db23db4b5d`, dated
2020-10-28T10:29:51Z with a GitHub-verified commit signature. Proxy and Git
source are byte-identical; sumdb agrees. The unreleased v1 branch and the
distinct maintained `/v2` line are not eligible exact-path releases.

V1.16.0 declares Go 1.14, and its complete isolated source/test closure
preserves Go 1.18 under both exact Go 1.26.7 and contained Go 1.18.10. It has
33 modules, 29 packages, 12 loaded external modules, and 342/279 complete-test
entries under the two SDKs. Native count-1, race, and broad cross-builds pass.
Both independent count-10 repeats fail under both SDKs because tests retain a
consumed buffer and leak package-global timeout/flag state. Vet fails because
production `runtime/context.go` discards the cancel from
`context.WithTimeout`; every one of the 53 stable exact-path releases contains
that timer/context resource defect. Adjacent v1.15.2 reproduces the repeat and
vet failures. Its API delta to v1.16.0 has zero incompatible and nine
compatible Swagger-option additions.

An independent runtime fixture covers routing and path parameters, query
handling, marshaling, metadata, unary/streaming/status mapping, cancellation
and deadlines, malformed input, cleanup, concurrency, global state, and
nil/panic boundaries. It passes count-10 and race under both SDKs while the
source vet defect remains. Gateway and Swagger protocol fixtures are
byte-deterministic within and across SDKs; malformed generator input exits 255
with timestamped and stack-bearing nondeterministic diagnostics.

The selected module has exactly two incoming edges, from etcd/api/v3 v3.5.1
and OTLP v0.7.0. `go mod why -m` is negative, repository imports are zero,
zero Gateway packages occur in the 429-entry complete project load, and it is
runtime-unreachable. Exact disposable `go get v1.16.0` merely adds a direct
root/source checksum, grows graph edges 3,599 -> 3,600, and leaves loading at
zero; tidy removes the root. Exact v1.15.2 instead causes an unrelated
234 -> 160 module downgrade cascade and makes the project unloadable. No
dependency edit was applied.

The fresh primary vulnerability index remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z with no exact Gateway record; its exact
OSV response is empty. Govulncheck v1.8.0 nevertheless reports 33 inherited
module IDs, 17 imported-package IDs, and eight called/reachable IDs in the
native v1.16.0 production and test closure: GO-2020-0036, GO-2022-0956,
GO-2023-1571, GO-2023-2153, GO-2024-2687, GO-2025-3372, GO-2026-4762, and
GO-2026-6061. Normal project scans contain zero Gateway occurrence because no
target package loads.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 project gates pass; the Go 1.18.10 projection retains only the
two accepted `pkg/shell` closed-file wording assertions. Focused quality
remains 27/27 Q0-Q2 PASS at L2 with zero held, regressed, or non-comparable
rows. Docker acceptance is unavailable because this host's Docker CLI lacks
buildx. The full quality wrapper also retains the known nested launcher
signal-retention timing race; every applicable stage passes independently.

# Role And Boundaries

This is a product-choice turn, not a renewed dependency audit, decision record,
or implementation. Reuse the completed evidence. Do not infer authorization
from prior retained-module decisions and do not change the repository.

# Guarded Decisions

The user's gRPC Prometheus v1.2.0 and gRPC middleware v1.0.0 option 1
decisions, plus the Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 exceptions, remain final and separate. Their exact selections,
incoming edges, negative why results, zero package loads, runtime
unreachability, and no-new-finding guards all revalidated unchanged. Stop for
the owning decision if a guard expires. Do not change mvn-pom-mutator,
GoConvey, either Gateway parent, or any unrelated module.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status,
reciprocal archive history, latest Google UUID dependency implementation,
P7/P8 state, guarded invariants, and `./codex-dev-start.sh --check`. Read this
archive, the answered Gateway evaluation, rolling handover, roadmap,
`go.mod`, `go.sum`, and referenced contracts. Reuse the completed technical
audit; do not repeat or broaden it.

# Three Moves

First, present the three bounded options with the stated tradeoffs and
recommendation. Second, request exactly one explicit user choice. Third, stop
without recording or implementing a choice or preparing another session.

# Automatic Handoff

This is a product-choice turn. First present the three options and request one
explicit choice. Then stop. Do not record an exception, change dependency
metadata, authorize work by implication, prepare a successor, or create a
commit until the user supplies the decision in a later bounded instruction.
Do not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
