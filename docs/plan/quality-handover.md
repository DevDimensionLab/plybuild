# Quality Upgrade Handover

Generated: 2026-09-14T12:29:52+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`. The gRPC Prometheus
  evaluation started from clean HEAD
  `f910bf879bbe4d79913bd620d7f8f689370bba23`, tree
  `558e9b44cc788170f21dac35a18cf88c68fb5f12`.
- The latest dependency implementation remains exact Google UUID v1.4.0
  commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. gRPC Prometheus and the
  retained modules since Google UUID have no dependency implementation or
  metadata commit.
- The gRPC Prometheus evaluation archive is answered. The sole NEXT archive is
  `docs/plan/agent-sessions/2026-09-14T122952+0200-decide-grpc-ecosystem-go-grpc-prometheus-product-direction.md`.
  It may obtain and record only one bounded product choice. Do not infer that
  choice, implement it, evaluate another dependency group, or begin P8.
- No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish, release,
  stash, revert, bypass cleanup, or remove the worktree.

## Lifecycle And Retained Decisions

P2A-P6 are complete. P7 remains active at the gRPC Prometheus product stop;
P8 is queued. Exact Go 1.26.7, every accepted dependency move through Google
UUID v1.4.0, and all earlier retained-module decisions remain final. Do not
combine another dependency group or begin P8.

The user's 2026-09-14 gRPC middleware option 1 decision is final. Exact
`github.com/grpc-ecosystem/go-grpc-middleware v1.0.0` retains only its recorded
missing-release-metadata/deterministic-closure, TLS test, retry cancellation/
timer-resource, and related qualification exceptions while exact v1.0.0 and
its sole mvn-pom-mutator v0.2.3 edge remain unchanged, zero packages load, it
remains runtime-unreachable, and no new advisory or independent disqualifier
appears. Direct import/loading, runtime reachability, a version or incoming-
edge change, or a new advisory/defect expires that exception.

Gorilla WebSocket v1.4.2 retains only GO-2026-6278's weak `math/rand` client-
mask behavior and its documented upstream full-source Go 1.26 cross-test race,
under the same exact-version/sole-mvn-pom-mutator-edge/negative-why/zero-load/
runtime-unreachable/no-new-finding guards. GopherJS
`v0.0.0-20181017120253-0766667cb4d1` retains only its separate option 1
exception through sole GoConvey v1.6.4. Enterprise Certificate Proxy v0.2.1
and GAX v2.7.0 retain only their separate Viper v1.15.0-edge exceptions. All
guards revalidated unchanged during this evaluation. Do not change
mvn-pom-mutator or GoConvey, reopen an earlier decision, or transfer any
exception to gRPC Prometheus.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or bypass launcher cleanup.

## gRPC Prometheus Bounded Stop

No exact-path stable release qualifies. The project was left byte-identical
and no target exception was created. The next session must obtain an explicit
product choice before retention, graph removal, replacement, or patch work.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
`https://github.com/grpc-ecosystem/go-grpc-prometheus.git`, public, archived,
enabled, non-fork, and Apache-2.0. Proxy and sumdb expose only stable v1.2.0;
v1.0/v1.1 are non-semantic tags. Master is an unreleased 27-commit-later
pseudo-version and deprecates the root project in favor of the different
maintained `github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus`
module path. The `draft-v2.0.0` branch is also an unreleased root module. No
retraction or formal module deprecation directive exists.

Exact v1.2.0 is signed annotated tag object
`502116f1a0a0c1140aab04fd3787489209b357d3`, peeled commit
`c225b8c3b01faf2899099b768856a9e916e5087b`, parents
`39de4380c2e0353a115b80b1c730719c79bfb771` /
`c451f7210e5f6bb49fe4ee5ccad1829104e2c447`, tree
`5548a81e1bd9d5476b5b6d07ba43b0dca6f9c435`, released
2018-06-04T12:28:56Z with a GitHub-verified signature. Its proxy archive is
byte-identical to Git and Apache-2.0 license SHA-256 is
`332c4f0a1a657a8937ad0caf7a335a31ec72343821481e364d8894f568c8b010`.
The release has 25 files, 14 Go files, five packages, two native tests, two
generated protobuf files, no benchmarks/fuzz/testdata/build-tag branches/cgo/
embed, and no `go.mod` or `go.sum`; the proxy synthesizes one module line.

## Closure, API, And Behavior

No deterministic release-native closure qualifies:

- Current source-time resolution selects 72 modules/135 edges, reaches external
  Go 1.26, and cannot compile native tests because modern Prometheus removed
  `prometheus.Handler`. Contained Go 1.18.10 cannot compile that closure.
- A historical 28-module reconstruction from post-release master preserves the
  floor and passes build/vet/race/Linux/Windows under both SDKs, but uses
  unreleased metadata. Both SDKs' two independent count-10 native repeats fail
  because server-stream tests inspect metrics before asynchronous server-
  handler updates.
- The project-selected closure builds production under both SDKs but reaches a
  Go 1.19 Genproto dependency and fails the strict Go 1.18 floor; native tests
  also fail on the removed handler API.

The root API exports default client/server collectors and interceptors,
constructors/options, histogram enablement, server registration/initialization,
and the four RPC type constants. Every package, generated file, example, and
test was inspected.

The independent 379-line behavior fixture SHA-256 is
`d6d2366bc2add1d9e7fac686de0b1b556a767841698599e061371ddab377058f`.
It passes count-1, two count-10 repeats, race, vet, focused cancellation/
deadline repeats, and 128-way supported concurrency under both SDKs. It proves:

- unary client success omits received-message counting while errors increment
  it; upstream fixed this only on unreleased master;
- repeated terminal client-stream receives double-count handled metrics;
- raw canceled/deadline errors retain identity but classify as Unknown;
- abandoned streams, send failures, and panics can leave started without
  handled observations;
- init globally registers eight counter vectors; histogram enablement ignores
  registration errors and can silently retain stale collectors; global
  enablement mutation is unsynchronized; and
- nil callback/options/descriptors/info/streamer boundaries can panic. Example
  servers bind fixed ports without graceful shutdown, but no example or target
  package is loaded by this project.

## Loading, MVS, Vulnerabilities, And Gates

Selected v1.2.0 exists only through main -> direct mvn-pom-mutator v0.2.3 ->
target, its sole incoming edge. `go mod why -m` is negative, source imports are
zero, and the 429-entry complete-test load contains zero target packages. It is
runtime-unreachable.

Exact disposable `go get target@v1.2.0` adds a direct requirement and grows
234 modules/3,599 edges/1,067 sum lines to 346/3,755/1,080. It also upgrades
four unrelated existing selections: BigQuery v1.8.0 -> v1.44.0, Datastore
v1.1.0 -> v1.10.0, Pub/Sub v1.3.1 -> v1.27.1, and Envoy control-plane v0.10.1
-> `v0.10.2-0.20220325020618-49ff273808a1`. Loading remains 429 entries,
197 module-backed packages, 41 modules, and zero target packages. Candidate and
base tidy projections converge byte-identically and remove the direct root.
No MVS effect was applied.

Fresh primary vulnerability data remains 1,398 records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
Last-Modified 2026-09-10T16:28:28Z, with no exact target record. The exact OSV
query is empty. Project scans remain 30 module findings, 22 package findings,
and 20 called IDs/22 symbol traces with zero target occurrence. Historical
target-closure results inherit old Prometheus/gRPC/x/net findings rather than
an exact target advisory.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 full preflight, verify, full tests, two count-10 repeats, race,
vet, pinned lint, empty-HOME count-2, Linux/Windows builds, API/CLI, 62
lifecycle controls, 80 mutation checks, and 15 audit-meta controls pass. The
contained Go 1.18.10 projection passes compatible verify/load/build/vet/cross-
build/repeats/race and retains only the two accepted `pkg/shell` closed-file
wording failures in the full count-1 suite. Accepted quality stays 27/27 Q0-Q2
PASS at L2, scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The lifecycle suite passed all 62 controls before the handoff edit; one final
rerun hit the already documented nested source-archive signal-retention timing
flake at control 51, and an immediate unchanged rerun passed all 62. This is a
launcher test-design/timing classification, not target behavior.
The 263-entry evidence manifest and decision-summary hashes are
`b0fede88c14009f0932c1cd6698cda3c2da1fd96acd4e14c7c10729543a83f61`
and `c5dfbf19a58b61f6ce0fba615bf0c6ba517dc3475eab524ddf2bc0fe3ae028be`.

## Next Bounded Objective

Present three choices and stop until the user explicitly selects one:

1. Recommended: retain exact inherited unloaded v1.2.0 under a new target-
   specific, non-transferable exception for only the documented closure, test,
   metric, lifecycle, archived, and global-state findings, guarded by exact
   version/sole edge/zero load/runtime unreachability/no new finding.
2. Authorize a fresh parent/removal evaluation centered on mvn-pom-mutator
   v0.2.3, owning all effects on middleware, Gorilla, Prometheus, and the graph.
3. Authorize a maintained replacement/patch program with an exact owned source
   identity, Go 1.18 floor, corrected behavior, and full compatibility gates.

Do not record or implement a choice until the user supplies it. Do not execute
another dependency group or begin P8.
