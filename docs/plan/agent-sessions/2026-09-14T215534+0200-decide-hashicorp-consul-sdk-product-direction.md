# Agent Session: Decide Hashicorp Consul SDK Product Direction

Status: NEXT
Session ID: `2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction`
Created: `2026-09-14T21:55:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ba29a645f8adef37d9be586b6e94ff1e3c17846c56e61bf625738ce4214e965d`
Previous: [2026-09-14T210212+0200-evaluate-hashicorp-consul-sdk-dependency.md](2026-09-14T210212+0200-evaluate-hashicorp-consul-sdk-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-14 explicit selection of
option 1 for exact-path `github.com/hashicorp/consul/sdk`. Retain exact
selected, inherited, unloaded v0.8.0 without dependency metadata changes under
the bounded retry, freeport, iptables, TempFile, test-server, API, closure-
vulnerability, and qualification exceptions below. Do not repeat the technical
audit, implement a dependency change, evaluate another dependency group, or
begin P8 in this decision-recording session.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active only for this Consul SDK product
decision after exact Go 1.26.7, every accepted dependency move through Google
UUID v1.4.0, and all retained-module decisions through Consul API v1.18.0. The
user has explicitly chosen to retain exact v0.8.0 through the sole historical
Consul API v1.12.0 graph edge while zero SDK packages load and the module
remains runtime-unreachable. This session may record that one choice and
prepare one bounded follow-up, but may not implement the choice or combine
another dependency group. P8 remains queued.

The completed Consul SDK evaluation found no exact-path stable release that
satisfies both the current Go 1.18 floor and behavior contracts. The user's
choice grants only the target-specific exception defined below; it does not
transfer the Consul API exception or authorize a direct SDK edge, another SDK
version, parent removal, a Go-floor change, or dependency modernization.

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

# Authorized Product Decision

On 2026-09-14 the user explicitly selected option 1 with the recommended
bounds: retain exact selected `github.com/hashicorp/consul/sdk v0.8.0` as an
inherited, unloaded selection without changing `go.mod` or `go.sum`. Accept
only the completed target findings: retry nontermination for negative count,
post-deadline retry, nil-Stop panic and absent cancellation; freeport duplicate,
zero-block, and process-global-state behavior; malformed and non-atomic
iptables handling with lost wrapped error identity; writable descriptors after
TempFile cleanup; the documented test-server timeout, context, deadline,
response-body, path, log-fatal, and double-Wait defects; the recorded API
incompatibilities; inherited closure-only x/sys findings; and related completed
qualification findings. This does not accept a new or independently discovered
defect.

The exception is SDK-specific and non-transferable. It is valid only while
exact v0.8.0 and its sole selected-version incoming edge from the historical
Consul API v1.12.0 graph vertex remain unchanged, the complete project load
contains zero SDK packages, the module remains runtime-unreachable, and no new
advisory or independent disqualifier appears. Revalidate and record those
guards. Direct import or loading, runtime reachability, a target version or
incoming-edge change, or a new advisory or independent defect expires the
exception and requires a fresh SDK dependency and product decision before
merge.

Do not add a direct SDK edge, select v0.13.0 or a later release, change or
remove the historical parent, raise the Go floor, authorize dependency
modernization, alter Viper, move unrelated modules, or manufacture a dependency
implementation commit. The Consul API decision and every earlier exception
remain separate. Do not stop or ask for this same Consul SDK decision again
while all guards hold.

# Role And Boundaries

This is a decision-recording session, not a renewed audit or implementation.
The user has explicitly selected and bounded option 1. Reuse the answered
evaluation; do not ask for the decision again or broaden it into direct use,
another SDK version, parent removal, Go-floor or dependency modernization,
Viper changes, unrelated-module authorization, another dependency group, or
P8.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status,
reciprocal archive history, latest dependency implementation identity, P7/P8
state, all guarded invariants, and `./codex-dev-start.sh --check`. Read this
archive, the answered SDK evaluation, Consul API decision, rolling handover,
roadmap, `go.mod`, `go.sum`, and referenced contracts. Do not repeat the audit.

# Three Moves

First, revalidate only exact selection, incoming edge, negative why, zero
target and guarded package loads, module hashes, runtime unreachability, and
current advisory state; reuse the completed audit and do not broaden it.
Second, record the exact option 1 exception, accepted findings, guards, expiry
triggers, and no-change result in the roadmap and rolling handover. Third,
answer this archive and prepare exactly one reciprocal NEXT mission for the
next bounded P7 group without executing it.

# Automatic Handoff

After recording an explicit choice, run applicable no-change lifecycle gates
and make the required local `docs: prepare next agent session` commit. Do not
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
