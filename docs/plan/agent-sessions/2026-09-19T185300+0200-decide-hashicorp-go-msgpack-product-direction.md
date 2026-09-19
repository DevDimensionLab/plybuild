# Agent Session: Decide Hashicorp Go Msgpack Product Direction

Status: NEXT
Session ID: `2026-09-19T185300+0200-decide-hashicorp-go-msgpack-product-direction`
Created: `2026-09-19T18:53:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `0d57e78496ee1264d1c8062357735b90ef371493f2e91c95600a56a0a6ca0629`
Previous: [2026-09-16T232707+0200-evaluate-hashicorp-go-msgpack-dependency.md](2026-09-16T232707+0200-evaluate-hashicorp-go-msgpack-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/hashicorp/go-msgpack`. The completed independent evaluation found
no exact-path stable release that passes every qualification contract. Obtain
or apply one explicit authorized choice from the options below, record its
exact accepted findings and expiry guards, and stop. Do not repeat the audit,
silently accept an exception, implement a dependency change, evaluate another
dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on the go-msgpack product decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2, and every recorded retained-module
decision through go-immutable-radix v1.3.1. All earlier outcomes and lifecycle
ancestors are final. Preserve every existing target-specific qualification or
exception; none transfers to go-msgpack. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-msgpack v0.5.3` remains inherited and unloaded. This
physical MVS selection is not qualification or risk acceptance. Do not add a
direct edge merely to alter MVS, change memberlist, Serf, Viper,
mvn-pom-mutator, another guarded parent, the Go floor, or any unrelated module.

# Measurements At Start

The evaluation began from clean handoff HEAD
`13f63a2c7639e373b261ca74176955b2be6461b2`, parent
`7c6d298b71e723a335dc295a1dc42a9b3ead3c93`, tree
`f41cd4a30e079f22d691477800e24c6b2e206366`. The latest dependency
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
1,067 sum lines, and the recorded 432-line tidy projection. Accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh exact-path evidence resolves the module to the active MIT Hashicorp
repository, itself a fork of `ugorji/go`. Selected v0.5.3 is commit
`be3a5be7ee2202386d02936a19ae4fbde1c77800`, tree
`92376aea28f7d767e816d4f24cc830bba62de0d0`. Highest non-retracted exact-path
stable v0.5.5 is commit `ad60660ecf9c5a1eae0ca32182ed72bab5807961`,
tree `51aaf54c47643758233939e0ef0f010610605846`. V0.5.0-v0.5.2 expose no
root package. V1.1.5 is retracted for breaking compatibility; v1.1.6 is
retracted and requires Go 1.19. The maintained `/v2` line is a distinct module
path and requires at least Go 1.24.

V0.5.3-v0.5.5 have a standard-library-only complete minimal production/test
closure, preserve Go 1.18, and expose byte-identical APIs. They build and
cross-compile under exact Go 1.26.7 and Go 1.18.10. None qualifies because
every native, repeated, and race test run under both SDKs exits during test
initialization: upstream `z_helper_test.go` calls `flag.Parse()` before the
testing package registers `-test.paniconexit0`. V0.5.3/v0.5.4 additionally
fail vet on a malformed legacy build comment; v0.5.5 fixes that tag and passes
vet but retains the early parse.

The independent behavior fixture SHA-256 is
`a788002acaf4eb5902c81e0c10d7c577ea1dd1b5ac873d2ce2f45e04969db864`.
Selected v0.5.3 panics when decoding into an existing map value and incorrectly
omits a non-nil pointer to `false`. V0.5.4 fixes the map panic. V0.5.5 fixes
both findings and passes the independent native/repeat/race/vet fixture under
both SDKs, but its upstream complete suite still fails before test execution.
Map wire order remains nondeterministic and decoder container length/recursion
limits remain caller responsibilities. Nil, aliasing, allocation, concurrency,
global-cache, malformed-input, error-identity, and RPC cleanup boundaries are
fully characterized in the answered evaluation.

Selected v0.5.3 has exactly four target requests from memberlist v0.1.3 and
v0.3.0 plus Serf v0.8.2 and v0.9.6. Its why result is negative, repository Go
imports are zero, production and complete-test target loads are zero, and it
is runtime-unreachable. Disposable exact v0.5.3/v0.5.4/v0.5.5 gets preserve
all 234 selections and every unrelated version with zero target load. They
manufacture only one direct indirect target root/graph edge; v0.5.4/v0.5.5
move only the target. Tidy removes the root and returns every copy to the same
base projection. No projection was applied.

Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no target record. Exact OSV results
for every serious release are empty. Go 1.26.7 isolated target scans are empty;
Go 1.18.10 reports only old-SDK standard-library advisories. Base and v0.5.5
project findings are identical with zero target package, symbol, test-symbol,
or reachable trace. All 13 earlier guarded modules retain their exact guards
and prior advisory state.

Every applicable no-change project gate passes: exact Go 1.26.7 verify, build,
count-one, two count-ten repeats, race, vet, pinned lint, API/CLI, empty-HOME,
four cross-builds, host/snapshot/Docker acceptance, full preflight, 80/80
mutation kills, and all 15 audit controls. No changed-selection scorecard
applies. The 521-entry selected-evidence manifest SHA-256 is
`100fa819df480ac176571f5d4fe672afd582554c9feb1fa8d9acc0fb50dc6a76`.

# Required Product Decision

Choose exactly one; option 1 is recommended because the target is unloaded and
it avoids manufacturing a dependency edge or accepting an unqualified upgrade:

1. Retain exact selected, inherited, unloaded v0.5.3 without metadata changes
   under a go-msgpack-specific exception. Accept only the completed map-update
   panic, pointer-`false` omission, malformed-build-comment/vet failure,
   upstream early-test-flag failure, and the fully characterized API, nil,
   aliasing, allocation, concurrency, global-state, resource, deterministic,
   malformed-input, and vulnerability findings. Guard exact v0.5.3, all four
   incoming requests, zero load, runtime unreachability, and no new advisory or
   independent defect.
2. Explicitly authorize a separately implemented exact v0.5.5 upgrade and a
   go-msgpack-specific exception for its upstream early-test-flag failure and
   characterized residual boundaries. This would fix both selected behavior
   defects but manufacture a direct indirect root and checksum lines. It needs
   a later dependency-only implementation commit and full changed-selection
   quality gate; this decision session must not implement it.
3. Keep P7 blocked and authorize a separately scoped remediation, fork,
   replacement, parent-removal, or architecture study. No such work is
   authorized in this decision session.

Do not infer acceptance from the current physical selection. If no explicit
authorized choice is available, report the blocker and preserve this decision
as the sole next boundary without changing source or dependency metadata.

# Role And Boundaries

This is a decision session, not a renewed audit or implementation. Reuse the
completed evaluation. Revalidate only selection, incoming requests, negative
why, zero target and guarded package loads, project hashes, runtime
unreachability, and current advisory state. Stop for the owning decision if a
guard changed. Do not broaden a choice into direct use, parent changes,
patching, forking, replacement, a Go-floor change, unrelated-module changes,
another dependency group, or P8.

# Required Reading

At start read this archive, the answered go-msgpack evaluation, rolling
handover, roadmap, `go.mod`, `go.sum`, the go-immutable-radix and go-hclog
decisions, and referenced lifecycle contracts. Verify feature branch, clean
ordinary and ignored status, ancestry, reciprocal archive history, latest
Google UUID implementation identity, P7/P8 state, all guarded invariants, and
`./codex-dev-start.sh --check`. Reuse the completed audit.

# Three Moves

First, revalidate only the decision guards and reuse the completed audit.
Second, obtain or apply exactly one explicit authorized option and record its
accepted findings, non-transferable bounds, expiry triggers, and implementation
status in the roadmap and rolling handover. Third, answer this archive and
prepare at most one reciprocal NEXT mission required by that choice. Do not
execute it.

# Automatic Handoff

Only after an explicit choice, run applicable no-change lifecycle gates and
make the required local `docs: prepare next agent session` commit. Without an
explicit choice, do not manufacture another handoff commit. Do not launch a
successor, push, merge, publish, release, stash, revert, bypass cleanup, remove
the worktree, combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
