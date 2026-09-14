# Agent Session: Decide gRPC Prometheus Product Direction

Status: NEXT
Session ID: `2026-09-14T122952+0200-decide-grpc-ecosystem-go-grpc-prometheus-product-direction`
Created: `2026-09-14T12:29:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `cd8227a9670cda181f61ea5d4426426bbf7b41f83dbd5eca98bd3ec89dd7a856`
Previous: [2026-09-14T112708+0200-evaluate-grpc-ecosystem-go-grpc-prometheus-dependency.md](2026-09-14T112708+0200-evaluate-grpc-ecosystem-go-grpc-prometheus-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and recording a fresh bounded product decision
for exact-path `github.com/grpc-ecosystem/go-grpc-prometheus`. The completed
evaluation found no exact-path stable release whose complete minimal source
and test closure both preserves Go 1.18 and passes every existing behavior,
test, API, MVS, and project contract. Present the three bounded options below
and ask the user to choose one. Do not infer a choice, implement a dependency
change, reopen an earlier decision, evaluate another dependency group, or
begin P8.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 is active only for this gRPC
Prometheus product decision; P8 remains queued. This session may record one
explicit user choice and prepare one bounded follow-up, but may not implement
that choice or combine another dependency group.

## Decision Required

Option 1 (recommended): retain exact selected, inherited, unloaded v1.2.0
without changing `go.mod` or `go.sum`, under a new target-specific exception.
Accept only the documented missing release module metadata and deterministic
complete closure, asynchronous native stream-test failures, unary client
received-message inversion, repeated terminal stream-receive counting, raw
cancellation/deadline classification as Unknown, and related archived/global-
state and lifecycle qualification findings. The exception must be non-
transferable and valid only while exact v1.2.0 and its sole mvn-pom-mutator
v0.2.3 edge remain unchanged, zero target packages load, the module remains
runtime-unreachable, and no new advisory or independent disqualifier appears.
Direct import/loading, runtime reachability, a version or incoming-edge change,
or a new advisory/defect expires it.

Option 2: authorize a separate bounded parent/removal evaluation centered on
the sole incoming `github.com/devdimensionlab/mvn-pom-mutator v0.2.3` edge.
That scope must independently own every resulting gRPC middleware, Gorilla
WebSocket, gRPC Prometheus, parent, graph, compatibility, and project effect;
it may not reuse or broaden any existing exception. Do not change the parent
or graph in this decision turn.

Option 3: authorize a separate maintained replacement or patch program. That
scope must choose and own an exact source/release identity, preserve Go 1.18,
repair the documented metric and stream-lifecycle defects, control registry
and global state, and prove API, vulnerability, MVS, and full project gates.
The maintained alternate module path, a fork, or a patched source tree is a
different dependency/architecture decision and must not be silently promoted.
Do not begin that design in this decision turn.

Do not offer the unreleased master pseudo-version, the `draft-v2.0.0` branch,
non-versioning tags v1.0/v1.1, or the different maintained
`github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus` path as a
qualified exact-path stable release.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves the
canonical public archived non-fork Apache-2.0 repository. The proxy exposes
only stable v1.2.0. Tags v1.0 and v1.1 are not semantic module versions;
master is an unreleased 27-commit-later pseudo-version and documents the root
module as deprecated in favor of a different module path. There is no
retraction or formal module deprecation directive.

V1.2.0 is a signed annotated tag at object
`502116f1a0a0c1140aab04fd3787489209b357d3`, peeled commit
`c225b8c3b01faf2899099b768856a9e916e5087b`, tree
`5548a81e1bd9d5476b5b6d07ba43b0dca6f9c435`, released
2018-06-04T12:28:56Z. Its proxy archive is byte-identical to Git and its
Apache-2.0 license SHA-256 is
`332c4f0a1a657a8937ad0caf7a335a31ec72343821481e364d8894f568c8b010`.
The release has 25 files, 14 Go files, five packages, two native tests, two
generated protobuf files, no examples as test functions, benchmarks, fuzz
targets, testdata, build tags, cgo, or embeds, and no `go.mod` or `go.sum`.
The proxy synthesizes only a module declaration.

No deterministic release-native closure qualifies. Current source-time
resolution selects 72 modules/135 edges and reaches Go 1.26 dependencies;
native tests do not compile because `prometheus.Handler` was removed.
Contained Go 1.18.10 cannot compile that modern closure. A reconstruction from
post-release master pins 28 historical modules and peaks at Go 1.9, but that
metadata was never released and both SDKs' two count-10 native repeats fail
because server-stream tests inspect counters and histograms before the server
handler's asynchronous update. The project-selected closure builds production
packages under both SDKs but includes a Go 1.19 dependency, so it does not
preserve the strict Go 1.18 floor; its native tests also fail to compile on the
removed handler API.

Independent fixtures characterize collectors, names, labels, histogram
options, initialization, duplicate registration, unary and stream success and
failure, cancellation/deadlines, error identity, panic/nil boundaries, and
128-way supported concurrency. The 379-line fixture SHA-256 is
`d6d2366bc2add1d9e7fac686de0b1b556a767841698599e061371ddab377058f`
and passes count-1, two count-10 repeats, race, and vet under both SDKs in the
historical and project-selected closures.

The audit confirms independent behavior defects. Unary client success omits
received-message counting while an error increments it. Repeated terminal
client-stream `RecvMsg` calls increment handled metrics repeatedly. Raw
`context.Canceled` and `context.DeadlineExceeded` retain error identity but
are classified as Unknown. Abandoned streams, send failures, and panics can
leave started RPCs without handled metrics. Init registers eight counter
vectors globally, histogram enablement ignores registration errors, repeated
enablement can silently retain stale collectors, and nil callback/options/
descriptors/info/streamer boundaries panic. Ordinary metric updates support
concurrency; global enablement mutation is not synchronized. The unary
received-count defect was fixed only after v1.2.0 on unreleased master.

Selected v1.2.0 exists only through main -> direct mvn-pom-mutator v0.2.3 ->
v1.2.0. This is the sole incoming target edge. `go mod why -m` is negative,
repository source imports are zero, and the 429-entry complete test load has
zero target packages, so the target is runtime-unreachable.

An exact disposable `go get` of the already selected target is not a no-op: it
adds a direct requirement, grows 234 modules/3,599 edges/1,067 sum lines to
346/3,755/1,080, and changes four unrelated existing selections: BigQuery
v1.8.0 to v1.44.0, Datastore v1.1.0 to v1.10.0, Pub/Sub v1.3.1 to v1.27.1,
and Envoy control-plane v0.10.1 to
`v0.10.2-0.20220325020618-49ff273808a1`. Loading remains unchanged and zero
for the target. Tidy converges byte-identically to the base tidy projection
and removes the target root. No graph change was authorized or applied.

Fresh primary vulnerability data has 1,398 module records, SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
Last-Modified 2026-09-10T16:28:28Z, and no exact target record. The OSV exact
target/version query is empty. Base project scans report 30 module findings,
22 package findings, and 20 called IDs/22 symbol traces with zero target
occurrence. Historical target-closure scans inherit old Prometheus, gRPC, and
x/net findings through examples/tests; none is an exact target advisory.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod` and `go.sum` remain at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 full preflight, repeats, race, vet, pinned lint, hermetic,
cross-build, API, CLI, lifecycle, mutation, and audit-meta gates pass. The Go
1.18.10 projection retains only the two accepted `pkg/shell` closed-file
wording failures. Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Role And Boundaries

This is a decision turn, not a renewed technical audit or dependency
implementation. Present the three bounded options exactly enough for the user
to choose; do not infer a choice, manufacture activity, or broaden the selected
option beyond its documented consequences.

## Guarded Decisions

The user's gRPC middleware option 1 decision and the Gorilla WebSocket,
GopherJS, Enterprise Certificate Proxy, and GAX exceptions remain final and
separate. Revalidate only their exact-selection, incoming-edge, negative-why,
zero-load, runtime-unreachable, and no-new-finding guards before recording a
choice. Stop for the owning decision if any expired. Do not change
mvn-pom-mutator or GoConvey or transfer an exception.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, reciprocal archive history, latest dependency implementation, all
guarded invariants, and `./codex-dev-start.sh --check`. Read this archive, the
answered gRPC Prometheus evaluation, rolling handover, roadmap, `go.mod`,
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
