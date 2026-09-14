# Agent Session: Decide gRPC Gateway Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T151903+0200-decide-grpc-ecosystem-grpc-gateway-product-direction`
Created: `2026-09-14T15:19:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1f395972cede9eeaf3535b491e6f7b77d0493d0ec093263ea61d8f75809be6de`
Previous: [2026-09-14T134759+0200-evaluate-grpc-ecosystem-grpc-gateway-dependency.md](2026-09-14T134759+0200-evaluate-grpc-ecosystem-grpc-gateway-dependency.md)
Next: [2026-09-14T180109+0200-evaluate-hashicorp-consul-api-dependency.md](2026-09-14T180109+0200-evaluate-hashicorp-consul-api-dependency.md)
Outcome: Recorded the user's bounded option 1 retention of exact inherited, unloaded Gateway v1.16.0; revalidated every guard without dependency metadata changes; and prepared the next bounded P7 dependency evaluation.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-14 explicit selection of
option 1 for exact-path `github.com/grpc-ecosystem/grpc-gateway`. Retain exact
selected, inherited, unloaded v1.16.0 without dependency metadata changes
under the bounded resource, native-test, inherited-vulnerability, generator-
diagnostic, and qualification exceptions below. Do not implement a dependency
change, evaluate another dependency group, or begin P8 in this decision-
recording move.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 remains active only for this
Gateway product decision; P8 remains queued. This session may record the one
explicit user choice and prepare one bounded follow-up, but may not implement
the choice or combine another dependency group.

# Authorized Product Decision

On 2026-09-14 the user explicitly selected option 1 with its recommended
bounds: retain exact selected
`github.com/grpc-ecosystem/grpc-gateway v1.16.0` as an inherited, unloaded
selection without changing `go.mod` or `go.sum`. Accept only the documented
universal discarded-cancel resource defect, native repeated-suite global and
test-state failures, eight inherited called/reachable vulnerability IDs in the
release-native closure, malformed-generator diagnostic nondeterminism, and
related recorded qualification findings. This does not accept a new or
independently discovered defect.

The exception is target-specific and non-transferable. It is valid only while
exact v1.16.0 and both selected-version incoming edges from etcd/api/v3 v3.5.1
and OTLP v0.7.0 remain unchanged, the complete project load contains zero
Gateway packages, the module remains runtime-unreachable, and no new advisory
or independent disqualifier appears. Revalidate and record those guards.
Direct import or loading, runtime reachability, a target version or incoming-
edge change, or a new advisory or independent defect expires the exception and
requires a fresh dependency and product decision before merge.

Do not add a direct target edge, select v1.15.2 or another v1 release, change
either parent, move to the different `/v2` module path, raise the Go floor,
authorize parent-removal or architecture work, move unrelated selections, or
manufacture a dependency implementation commit. Existing gRPC Prometheus,
gRPC middleware, Gorilla WebSocket, GopherJS, Enterprise Certificate Proxy,
and GAX exceptions remain separate. Do not stop or ask for this same Gateway
decision again while all guards hold.

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

This is a decision-recording turn, not a renewed dependency audit or
implementation. The user has explicitly selected and bounded option 1. Reuse
the completed evidence; do not ask for the decision again or broaden it into
direct use, another version, parent ownership/removal, `/v2`, a floor change,
or unrelated-module authorization.

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

First, revalidate only the unchanged exact selection, both incoming edges,
negative why result, zero Gateway and guarded-target package loads, module
hashes, runtime unreachability, and current advisory state. Reuse the completed
audit; do not repeat or broaden it. Second, record the exact option 1 exception,
accepted findings, guards, expiration triggers, and no-change result in the
roadmap and rolling handover. Third, answer this archive and prepare one
reciprocal NEXT mission for the next bounded P7 group without executing it.

# Automatic Handoff

After recording the decision, run the applicable no-change lifecycle gates and
create the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Recorded the user's 2026-09-14 option 1 selection and retained exact-path
`github.com/grpc-ecosystem/grpc-gateway v1.16.0` as an inherited, unloaded
module without changing `go.mod` or `go.sum` and without manufacturing a
dependency implementation commit.

The target-specific, non-transferable decision accepts only the findings
already established in the answered evaluation: the universal production
discarded-cancel timer/context resource defect, native repeated-suite consumed-
buffer and package-global timeout/flag state failures, the eight inherited
called/reachable vulnerability IDs in the release-native closure, malformed-
generator diagnostic nondeterminism, and related recorded qualification
findings. No new or independently discovered defect is accepted.

Every guard was revalidated from clean decision HEAD
`b9e4b64d00d5593bb1a9d45e4506e55fd6eeb1f4`, tree
`a5a2d5d6d7dafd058b6df74142463fae112561f3`, with a freshly unpacked exact
Go 1.26.7 distribution. Its archive and binary SHA-256 values are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
It ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, `LC_ALL=C`, `LANG=C`, and no ambient `GOFLAGS`.

Project selection remains Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC
middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0. Gateway retains exactly two selected-version incoming edges:
etcd/api/v3 v3.5.1 -> Gateway v1.16.0 and OTLP v0.7.0 -> Gateway v1.16.0.
The earlier guarded targets retain their exact selected versions and sole
recorded incoming edges from mvn-pom-mutator v0.2.3, GoConvey v1.6.4, or Viper
v1.15.0. All seven `go mod why -m` results remain negative, and repository Go
source has zero imports of every guarded target.

The complete project test load remains 429 entries, including 197 module-
backed packages across 41 loaded modules, with exactly zero packages from
Gateway or any earlier guarded target. Together with the unchanged selected-
version edges, negative why results, and zero source imports, this revalidates
that all seven targets remain runtime-unreachable. The project still selects
234 modules with 3,599 graph edges.

`go.mod` and `go.sum` remain byte-identical to clean decision HEAD at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
`go.sum` remains 1,067 lines. The applicable exact-Go no-change module, build,
test, race, vet, and launcher lifecycle gates pass. No changed-selection
scorecard run applies; accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains byte-identical at 1,398 module
records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact Gateway, gRPC
Prometheus, gRPC middleware, GopherJS, Enterprise Certificate Proxy, or GAX
record. Exact OSV target/version queries for those six modules remain empty.
Gorilla retains only GO-2020-0019 and the selected-affecting GO-2026-6278; the
latter remains unwithdrawn at record SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`
with its existing GHSA alias. No new advisory or independently observed
disqualifier appeared during this guard-only revalidation; the completed
technical audit was not repeated or broadened.

The Gateway exception remains valid only while exact v1.16.0 and both incoming
etcd/api/v3 v3.5.1 and OTLP v0.7.0 edges remain unchanged, zero Gateway
packages load, the module remains runtime-unreachable, and no new advisory or
independent disqualifier appears. Direct import or loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent defect expires the exception and requires a fresh Gateway
dependency and product decision before merge. Every earlier exception remains
separate under its own guards. Do not request this same Gateway decision again
while all guards hold.
