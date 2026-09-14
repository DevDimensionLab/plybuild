# Agent Session: Decide Hashicorp Consul SDK Product Direction

Status: NEXT
Session ID: `2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction`
Created: `2026-09-14T21:55:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `50c9a575444608a5f2c8b106409d038777527bb61272c19f41a1693b1ef040d2`
Previous: [2026-09-14T210212+0200-evaluate-hashicorp-consul-sdk-dependency.md](2026-09-14T210212+0200-evaluate-hashicorp-consul-sdk-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Resolve the fresh bounded P7 product decision for selected exact-path
`github.com/hashicorp/consul/sdk v0.8.0` after its completed independent
evaluation found no qualified stable release that preserves the Go 1.18 floor.
Present or record one explicit option with exact bounds. Do not repeat the
technical audit, implement a choice, change dependency metadata, combine
another dependency group, or begin P8 in this decision session.

# Authorized Roadmap

P2A-P6 are complete. P7 is active but stopped for this Consul SDK decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, and all retained-module decisions through Consul API v1.18.0. P8 and
further dependency groups remain queued. Earlier outcomes and lifecycle
ancestry are final.

The Consul SDK evaluation found no exact-path stable release that satisfies
both the current Go 1.18 floor and behavior contracts. Selected v0.8.0 remains
inherited and unloaded without an exception. Do not infer retention from zero
load or transfer the separate Consul API exception.

The available bounded choices are:

1. Recommended: retain exact selected, inherited, unloaded v0.8.0 without
   changing `go.mod` or `go.sum`, under a new target-specific exception that
   accepts only the completed SDK findings. Keep it valid only while exact
   v0.8.0 and the sole selected-version incoming edge from the historical
   Consul API v1.12.0 vertex remain unchanged, zero SDK packages load, the
   module remains runtime-unreachable, and no new advisory or independently
   discovered disqualifier appears. Direct import/loading, runtime
   reachability, a version or incoming-edge change, or a new advisory or defect
   must expire the exception and require a fresh SDK dependency/product
   decision before merge.
2. Authorize a new bounded parent/graph-removal investigation to eliminate the
   historical Consul API v1.12.0 vertex and thereby the SDK selection. This
   reopens owning ancestry and may affect other modules. Record the authority
   and prepare that investigation, but do not implement it in this decision
   session.
3. Authorize a new bounded Go-floor and dependency-modernization plan around a
   newer SDK line. V0.13.0 still fails behavior qualification; v0.13.1 and all
   later releases exceed Go 1.18, while latest v0.18.2 declares Go 1.26.7 and
   moves unrelated selections. Record the intended floor/architecture scope,
   but do not implement it in this decision session.

If the user has not explicitly selected an option, stop and request that
choice. Do not manufacture authorization. If the user selects an option,
record exactly that choice and its guards; do not broaden it.

The user's 2026-09-14 Consul API option 1 decision retains exact inherited,
unloaded `github.com/hashicorp/consul/api v1.18.0` without metadata changes. It
accepts only its documented completed findings and is valid only while exact
v1.18.0, the sole Viper v1.15.0 incoming edge, zero package load, runtime
unreachability, and no-new-finding guards hold. Direct import/loading, runtime
reachability, a version or edge change, or a new advisory/independent defect
requires a fresh owning decision. Do not change Viper or transfer this
exception to SDK.

Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC middleware v1.0.0, Gorilla
WebSocket v1.4.2, GopherJS `v0.0.0-20181017120253-0766667cb4d1`, Enterprise
Certificate Proxy v0.2.1, and GAX v2.7.0 decisions remain final under their
exact-selection, selected-version incoming-edge, zero-load, runtime-
unreachable, and no-new-finding guards. Do not change their owning parents,
mvn-pom-mutator, GoConvey, or Viper.

# Measurements At Start

The evaluation began from clean HEAD
`19e0f4a2e36c654c9e1810b8f660cec2a29f22fb`, parent
`474909583bfe6efbc5ae49b9ff1d68ec7402ee41`, and tree
`e1d077544884d461318861fb1d6aad72e8caa202`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

Fresh identity evidence resolves canonical public active unarchived non-fork
`https://github.com/hashicorp/consul.git`, `sdk` module subdirectory, and
MPL-2.0 candidate source. The proxy lists 30 stable releases and five RCs, all
with exact `sdk/` Git tags. Serious tags are lightweight and have no GitHub
Release object. Selected v0.8.0 is commit
`b1ee900870b9f04377d64936d4af452a5799b72c`, parent
`a4a43460e51d66bc562fbea476844b6c13418543`, repository tree
`8e328b47681316d07b7abbd89dff743dfbf1293d`, and SDK subtree
`dcbc915786cb1a668b9ca0f10bfd129b4526b4a8`.

All 30 stable declarations and serious complete source/test closures were
inspected. V0.13.0 is the highest closure-eligible release and declares Go
1.12 with imported closure no higher than Go 1.17. V0.13.1 is the first Go
1.19 release; latest v0.18.2 declares Go 1.26.7.

Every floor-eligible stable v0.1.0-v0.13.0 has an exported retry `Counter`
whose negative `Count` never terminates because its stop condition is equality-
only. Selected v0.8.0 and highest-candidate v0.13.0 additionally reproduce
retry after a timer deadline, nil `Stop` panic, duplicate and zero-block
freeport defects, malformed iptables input acceptance, and a descriptor that
remains writable after `TempFile` cleanup. Selected test-server helpers ignore
their ready-timeout setting, omit HTTP/subprocess contexts and deadlines, leak
successful service/check response bodies, build unescaped paths, may terminate
the caller through `log.Fatal`, and include a v0.8.0 double-`Wait` cleanup
defect. These exact target findings are independent disqualifiers.

Exact source verifies, builds, vets, cross-builds, and—with only a scratch
overlay for the managed sandbox's denied Darwin sysctl probe—passes count-one,
two count-ten repeats, and race under exact Go 1.26.7 and contained Go 1.18.10.
Independent fixtures reproduce the defects under both Go lines. Pinned API
diff records incompatible v0.8.0-to-v0.13.0 and latest-line changes.

Selected v0.8.0 exists only through one selected-version edge from the
historical Consul API v1.12.0 graph vertex. Its why result is negative,
repository imports are zero, zero SDK packages occur in the complete project
load, and it is runtime-unreachable. Exact candidate gets retain zero load but
manufacture a direct root; v0.13.0 also adds go-version, and later choices
violate the floor or move unrelated selections. No projection was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no SDK record. Exact candidate OSV
queries are empty. Isolated old closures have only inherited x/sys findings;
the project has zero SDK package, symbol, test-symbol, or reachable-trace
occurrence.

The unchanged project remains 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed packages across 41 loaded modules,
1,067 `go.sum` lines, and a 432-line unapplied tidy projection. `go.mod` and
`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Applicable exact-Go gates pass, all earlier exception guards remain valid, and
accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Role And Boundaries

This is a decision-recording session, not a renewed audit or implementation.
Reuse the answered evaluation. Do not add a direct SDK edge, select v0.13.0 or
a later release, change the historical parent, raise the Go floor, alter Viper,
move unrelated modules, implement another dependency group, or begin P8 unless
the selected option explicitly authorizes a later bounded mission—and even
then only prepare that mission here.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status,
reciprocal archive history, latest dependency implementation identity, P7/P8
state, all guarded invariants, and `./codex-dev-start.sh --check`. Read this
archive, the answered SDK evaluation, Consul API decision, rolling handover,
roadmap, `go.mod`, `go.sum`, and referenced contracts. Do not repeat the audit.

# Three Moves

First, if no explicit user choice is present, present the three bounded options
and stop for one. Second, after an explicit choice, revalidate only exact
selection, incoming edge, negative why, zero target/guarded loads, module
hashes, runtime unreachability, and current advisory state. Third, record the
choice and exact expiry guards in roadmap/handover, answer this archive, and
prepare exactly one reciprocal NEXT mission authorized by that choice without
executing it.

# Automatic Handoff

After recording an explicit choice, run applicable no-change lifecycle gates
and make the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
