# Agent Session: Decide Hashicorp Go Multierror Product Direction

Status: NEXT
Session ID: `2026-09-19T203710+0200-decide-hashicorp-go-multierror-product-direction`
Created: `2026-09-19T20:37:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c7f749ee1d9725c10846fea7aafbc338ebad24effa01764e6621255b227708f1`
Previous: [2026-09-19T190917+0200-evaluate-hashicorp-go-multierror-dependency.md](2026-09-19T190917+0200-evaluate-hashicorp-go-multierror-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/hashicorp/go-multierror`. The completed independent evaluation
found no exact-path stable release that passes every qualification contract.
Obtain or apply one explicit authorized choice from the options below, record
its exact accepted findings and expiry guards, and stop. Do not repeat the
audit, silently accept an exception, implement a dependency change, evaluate
another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on the go-multierror product decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2, and all recorded retained-module
decisions through go-msgpack v0.5.3. Every earlier outcome and lifecycle
ancestor is final. Preserve every existing target-specific qualification or
exception; none transfers to go-multierror. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-multierror v1.1.0` remains inherited and unloaded.
This physical MVS selection is not qualification or risk acceptance. Do not
add a direct edge merely to alter MVS, change memberlist, Serf, mitchellh/cli,
posener/complete, Viper, mvn-pom-mutator, another guarded parent, the Go floor,
or any unrelated module.

The go-msgpack v0.5.3, go-immutable-radix v1.3.1, and go-hclog v1.2.0
exceptions, qualified go-cleanhttp v0.5.2 result, and the Errwrap v1.0.0,
Consul SDK v0.8.0, Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus
v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain separate under their recorded guards. Stop for
the owning decision if any guard changes. In particular, an added main-module
Errwrap root or a new go-multierror-v1.1.1-to-Errwrap graph edge expires the
separate Errwrap guard and cannot be accepted implicitly here.

# Measurements At Start

The evaluation recording began from clean handoff HEAD
`ea7bc704e1160eca51b3b35d00411db28babcac8`, parent
`69017090043510b4754e8e81feab762a8011fe90`, tree
`c2b8da74a6a6b637ea5c51f4a5b8b871ef0c1019`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive chain, and `./codex-dev-start.sh --check` rather than
assuming them.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The project remains 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries across 41 loaded modules,
1,067 sum lines, and the recorded 432-line unapplied tidy projection.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to Hashicorp's public, active, unarchived, non-fork MPL-2.0 repository.
The proxy exposes only v1.0.0, v1.1.0, and v1.1.1; v1.1.1 is `@latest`.
There are no retractions, module deprecations, redirects, alternate major
module lines, or GitHub Release objects. Main is 73 commits beyond v1.1.1 but
is unreleased and was not promoted.

Selected v1.1.0 is lightweight tag/verified commit
`2004d9dba6b07a5b8d133209244f376680f9d472`, parent
`a5e98147459549195562243aaf8caf5cedbb4864`, tree
`fcc486d753842d176b67384af19ca7c51157c43c`. Latest v1.1.1 is lightweight
tag/verified merge commit `9974e9ec57696378079ecc3accd3d6f29401b3a0`,
parents `ab6846acb1cf402787cde041fcff5f01c0b8e362` and
`0023bb0ab1225a10509dca08627b42bed25d74a5`, tree
`93989e7fff83edc30d172b8a10d3dcb9c67b8cc1`. Proxy and Git source bytes
agree for all three releases.

All three releases have the same minimal module closure: target plus exact
`github.com/hashicorp/errwrap v1.0.0`. Complete imported production and test
source preserves Go 1.18. Each release has one package and no command,
example, benchmark, fuzz target, testdata, generated file, platform/build-tag
branch, cgo, embed, generator, or symlink. Every release passes source
verification, native count-one, two independent count-ten repeats, race, vet,
and production/test cross-compilation for Darwin AMD64, Linux AMD64/ARM64,
Windows AMD64, and js/wasm under exact Go 1.26.7 and Go 1.18.10.

V1.0.0 lacks the v1.1 `Group` and `(*Error).Unwrap` API. Those are compatible
additions in v1.1.0; v1.1.0 and v1.1.1 have no API diff. The independent
passing behavior fixture SHA-256 is
`e5fd467b66758b1d2c79d3a5b02f8eaea0c28e4c688e234e13d45fd1c87806c3`.
It characterizes formatting/order, Append and Flatten mutation/aliasing,
ErrorOrNil allocation, nil/panic boundaries, ordinary Is/As/Unwrap traversal,
Group concurrency/lifecycle, absence of package resource/global-state
ownership, and deterministic versus completion-order behavior under both
SDKs. V1.0.0 and selected v1.1.0 panic on `WrappedErrors` through a typed-nil
receiver; v1.1.1 fixes that panic.

No stable release qualifies because `Prefix` delegates to Errwrap v1.0.0 and
destroys standard error identity. The blocker fixture SHA-256 is
`1620577fe5db5a1e76e50d68498165924d54513848b8093f94c94518d87872ac`.
Under both SDKs, v1.1.0 and v1.1.1 fail all three assertions: `errors.Is`
cannot traverse Prefix on a plain error or a multierror member, and
`errors.As` cannot traverse a prefixed multierror member. V1.0.0 has the same
Prefix design. Upstream issue 56 and unmerged PR 155 head
`8224dca4d7fdce17091ec99615804e283e17bf13` corroborate the defect, but the
fix is unreleased and cannot be promoted.

Selected v1.1.0 has six requests: Serf v0.9.6 requests v1.1.0; memberlist
v0.1.3/v0.3.0, mitchellh/cli v1.0.0/v1.1.0, and posener/complete v1.2.3
request v1.0.0. The shortest path is main -> mvn-pom-mutator v0.2.3 ->
historical Viper v1.10.1 -> Serf v0.9.6 -> go-multierror. Its why result and
repository import search are negative and production/complete-test loads
contain zero target packages, so it is runtime-unreachable.

Disposable exact v1.1.0/v1.1.1 gets preserve all 234 selections and the
355/429 package populations with zero target load. Exact v1.1.0 manufactures
main-module indirect roots for both Multierror and Errwrap plus two graph
edges and source checksums without changing a selection. Exact v1.1.1 makes
those same roots and moves only Multierror to v1.1.1, adding a new
v1.1.1-to-Errwrap graph edge; that expires the separate Errwrap guard.
V1.0.0 removes mvn-pom-mutator v0.2.3 and makes the project unloadable.
Tidy removes the v1.1 roots and returns both viable projections byte-identical
to the base tidy projection, selecting v1.1.0 again. No projection was applied.

Fresh primary vulnerability data has 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no target record. Exact OSV results
for all three releases are empty. Go 1.26.7 isolated scans are empty; Go
1.18.10 reports only its old standard-library advisory population. Base and
v1.1.1 project module/package/symbol/test-symbol scan outputs are byte-
identical with zero target advisory or trace. All 14 earlier guards remain
intact.

Every applicable no-change project gate passes: exact Go 1.26.7 verify,
build, count-one, two count-ten repeats, race, vet, pinned lint, API/CLI,
empty-HOME, four cross-builds, full preflight, 17 script/meta stages, 80/80
mutation kills, host/snapshot/Docker acceptance, and all 15 audit controls.
The authoritative clean-tree scorecard is 27/27 Q0-Q2 PASS at L2. The
469-entry selected-evidence manifest SHA-256 is
`7db323255b0c8361371237d1e3e7405e697e3ad7f53167678f756535dc66b7f3`.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target is unloaded,
it preserves every parent and the separate Errwrap guard, and no stable
upgrade repairs the blocking Prefix defect:

1. Retain exact selected, inherited, unloaded v1.1.0 without metadata changes
   under a go-multierror-specific exception. Accept only the completed Prefix
   `errors.Is`/`errors.As` identity loss, typed-nil `WrappedErrors` panic, and
   fully characterized API, formatting/order, nil/panic, mutation, aliasing,
   allocation, concurrency, Group lifecycle, resource/global-state, MVS,
   vulnerability, and related findings. Guard exact v1.1.0, all six incoming
   requests, zero target load, runtime unreachability, the unchanged Errwrap
   guard, and no new advisory or independent defect.
2. Explicitly authorize a separately implemented exact v1.1.1 selection and
   go-multierror-specific Prefix exception. This fixes only the typed-nil
   `WrappedErrors` panic, requires new indirect Multierror and Errwrap roots,
   and expires the separate Errwrap guard. It therefore also requires a fresh,
   separately sequenced Errwrap product decision before merge plus a later
   dependency-only implementation and full changed-selection gate. This
   decision session must not implement it or silently decide Errwrap.
3. Keep P7 blocked and authorize a separately scoped remediation, fork,
   replacement, parent-removal, or architecture study. No such work is
   authorized in this decision session.

Do not infer acceptance from the current physical selection. If no explicit
authorized choice is available, report the blocker and preserve this decision
as the sole next boundary without changing source or dependency metadata.

# Role And Boundaries

This is a decision session, not a renewed audit or implementation. Reuse the
completed evaluation. Revalidate only selection, all six incoming requests,
negative why, zero target and guarded package loads, project hashes, runtime
unreachability, separate Errwrap guard, and current advisory state. Stop for
the owning decision if a guard changed. Do not broaden a choice into direct
use, parent changes, patching, forking, replacement, a Go-floor change,
unrelated-module changes, another dependency group, or P8.

# Required Reading

At start read this archive, its answered go-multierror evaluation, the
answered Errwrap evaluation/decision, go-msgpack decision, rolling handover,
roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored
status, handoff ancestry and changed-file set, exact toolchain identity,
reciprocal archive history, all guarded selections/requests/load results, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, revalidate the narrow guards and stop for the owning decision if any
changed. Second, obtain or apply exactly one explicit choice above without
implementation, record the accepted findings and precise expiry conditions,
and preserve every unrelated selection. Third, update the roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal NEXT mission for
the authorized bounded continuation, and commit the handoff. Do not execute
the successor.

# Automatic Handoff

After a coherent explicit decision, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
