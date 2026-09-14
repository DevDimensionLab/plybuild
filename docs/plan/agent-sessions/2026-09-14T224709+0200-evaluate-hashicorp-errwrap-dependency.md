# Agent Session: Evaluate Hashicorp Errwrap Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T224709+0200-evaluate-hashicorp-errwrap-dependency`
Created: `2026-09-14T22:47:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1846cd70f1df8adf9cc72f76931128475271c1168dee69b3736fdaeab8b5b491`
Previous: [2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction.md](2026-09-14T215534+0200-decide-hashicorp-consul-sdk-product-direction.md)
Next: [2026-09-14T234716+0200-decide-hashicorp-errwrap-product-direction.md](2026-09-14T234716+0200-decide-hashicorp-errwrap-product-direction.md)
Outcome: Completed exact-path Hashicorp Errwrap qualification; neither stable release satisfies the behavior contracts, so the unchanged project stops for a fresh bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/errwrap v1.0.0` as one bounded dependency group. Resolve
its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, and all recorded retained-module
decisions through Consul SDK v0.8.0. All earlier outcomes and lifecycle
ancestry are final. Evaluate only Hashicorp Errwrap in this session; do not
reopen or combine another dependency group. P8 remains queued.

The user's 2026-09-14 Consul SDK option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/consul/sdk v0.8.0` without dependency
metadata changes. It accepts only the completed retry, freeport, iptables,
TempFile, test-server, API, inherited closure-only x/sys, and related
qualification findings. Its exception is target-specific and valid only while
exact v0.8.0 and the sole selected-version historical Consul API v1.12.0 edge
remain unchanged, zero SDK packages load, the module remains runtime-
unreachable, and no new advisory or independent disqualifier appears. Direct
import/loading, runtime reachability, a version or incoming-edge change, or a
new advisory or independent defect requires a fresh Consul SDK dependency and
product decision before merge. Revalidate these guards before work and stop for
that owning decision if any fails. Do not transfer its exception to Errwrap or
change the historical parent.

The Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC
middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain final and separate under their exact-selection,
selected-version incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards and stop for the fresh owning decision
if one expires. Do not change either Gateway parent, mvn-pom-mutator, GoConvey,
Viper, or transfer an earlier exception.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Consul SDK, Consul API, and
every retained group since Google UUID have no dependency implementation or
metadata commit.

The clean Consul SDK decision retained exact v0.8.0 through exactly the sole
selected-version incoming edge from the historical Consul API v1.12.0 graph
vertex. Its `go mod why -m` result remains negative, repository source imports
are zero, and zero SDK packages occur in the complete project load. The
unchanged project has 234 selected modules, 3,599 graph edges, 429 complete-
test entries, 197 module-backed packages, 41 loaded modules, 1,067 `go.sum`
lines, and a 432-line unapplied tidy projection. Exactly zero SDK, Consul API,
Gateway, gRPC Prometheus, gRPC middleware, Gorilla WebSocket, GopherJS,
Enterprise Certificate Proxy, or GAX packages load. `go.mod` and `go.sum`
SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains byte-identical at 1,398 module
records, index SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`,
and Last-Modified 2026-09-10T16:28:28Z. It has no exact SDK, Consul API,
Gateway, gRPC Prometheus, gRPC middleware, GopherJS, Enterprise Certificate
Proxy, or GAX record. Their exact target/version OSV responses remain empty.
Gorilla retains only its recorded entries, including unwithdrawn GO-2026-6278
at record SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.

A bounded queue survey identifies selected Hashicorp Errwrap v1.0.0 in the
current build list. That selection is not proof of repository identity,
release qualification, ancestry, floor, package loading, behavior,
vulnerability state, or suitability. Resolve those facts independently and do
not add a direct edge merely to alter MVS.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`
and GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
Preserve the known apidiff archive reproducibility discrepancy, Python 3.14
Docker timestamp control, and nested launcher signal-retention timing race;
none is Hashicorp Errwrap evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, major
module lines, and every serious exact-path stable candidate. Do not silently
promote a redirect, fork, alternate path, different major module path,
prerelease, non-versioning tag, unreleased branch head, or floor-ineligible
release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph and from every
already-final guarded dependency decision.

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize wrapping and unwrapping semantics, formatting,
error identity and interoperability, nil/panic behavior, mutation and aliasing,
allocation, concurrency, global state, resource cleanup, and malformed inputs.
Distinguish tools, examples, optional packages, and test-only helpers from
behavior actually loaded by this project.

Add independent fixtures where useful for deterministic wrapping, formatting,
error identity, nil/panic boundaries, malformed inputs, mutation, aliasing,
supported concurrency, and selected-project consumers. Run source
verification, package listing, native complete tests, two independent repeats,
race, vet, and meaningful cross-builds under both SDKs. Classify every failure
precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent implementation.
Compare primary vulnerability results at module, package, symbol, test-symbol,
and reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, every unchanged exception guard, and
`./codex-dev-start.sh --check`. Read this archive, the answered Consul SDK
decision and evaluation, the Consul API decision, the earlier guarded-decision
archives, rolling handover, roadmap, `go.mod`, `go.sum`, and every referenced
quality, compatibility, release, runner, evidence, and lifecycle contract.
Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose
complete minimal source/test closure preserves Go 1.18. Do not promote an
unqualified or floor-ineligible identity. If no candidate satisfies the
existing contracts, preserve the evidence and stop for a fresh bounded product
decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, record the exact decision and evidence in the roadmap and rolling
handover, answer this archive, and follow the lifecycle contract for one next
bounded authorized mission. Do not execute the successor in this turn.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the Hashicorp Errwrap decision or bounded product stop, rewrite the
rolling handover and roadmap, answer this archive, and follow the repository
lifecycle contract. Do not implement another dependency group, launch a
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable release qualifies. The only stable releases are v1.0.0
and v1.1.0; both have complete stdlib-only source/test closures that preserve
Go 1.18, making v1.1.0 the highest floor-eligible candidate. Both nevertheless
violate the package's exported concrete-type matching promise because they
compare `reflect.Type.String()`: different types in different import paths
with the same package/type spelling falsely match. Both also fail to traverse
standard `Unwrap() []error` graphs such as `errors.Join` under Go 1.26.7, and
selected v1.0.0 additionally lacks standard single-error unwrapping. These
independent target defects reproduce in contained fixtures; no earlier product
exception transfers to Errwrap.

Canonical proxy, sumdb, `go-import`, Git, and GitHub evidence resolves public,
active, unarchived, non-fork `https://github.com/hashicorp/errwrap`, MPL-2.0,
no GitHub Release objects, no retractions or module deprecation, and no `/v2`
line. Selected v1.0.0 is lightweight tag/commit
`8a6fb523712970c966eefc6b39ed2c5e74880354`, parent
`d6c0cd88035724dd42e0f335ae30161c20575ecc`, tree
`9863613ad8fe960290d1631d84ede660af5f0746`, dated
2018-08-24T00:39:10Z. V1.1.0 is lightweight tag/verified merge commit
`7b00e5db719c64d14dd0caaacbd13e76254d02c0`, parents v1.0.0 and
`96a78ad11c51762df122b738e6f4f30f58e03d8d`, tree
`aefd62cf8e9549e154a65afb8b4524b7b6ed5f2e`, dated
2020-07-14T15:51:01Z. Proxy and Git bytes match, strict Git verification
passes, and v1.0.0 is an ancestor of v1.1.0. Unreleased master declares Go
1.24 and was not promoted.

Each release has one package, one production file, one test file, and no
commands, examples, benchmarks, fuzz targets, testdata, generated files,
build-tag/platform variants, cgo, embeds, go:generate directives, symlinks, or
non-stdlib dependencies. Both verify, build, pass native count-one, two
count-ten repeats, race, vet, and production/test cross-compilation for Darwin
amd64, Linux amd64/arm64, Windows amd64, and js/wasm under exact Go 1.26.7 and
contained Go 1.18.10. Pinned API diff finds the same nine exported functions
plus `WalkFunc` and `Wrapper` in both releases with no incompatible declaration
change. V1.1.0 adds standard single-error unwrapping and deprecates `Wrapf`.

The independent fixtures also characterize deterministic wrapping,
formatting, outer-to-inner order, lookup identity, nil lookup, explicit nil
outer-error and nil callback panics, error aliasing, isolated returned slices,
custom child order, supported immutable concurrent reads, and allocation. The
module has no global state or resources; recursion has no cycle guard.

Selected v1.0.0 exists through exactly one selected-version incoming edge from
`github.com/hashicorp/go-multierror v1.1.0`; a historical multierror v1.0.0
edge also requests it. Its shortest graph path is main -> mvn-pom-mutator
v0.2.3 -> historical Viper v1.10.1 -> Serf v0.9.6 -> go-multierror v1.1.0 ->
Errwrap v1.0.0. Its why result is negative, repository imports are zero, and
zero target packages occur in production or complete-test loads, so it is
runtime-unreachable. Disposable exact gets preserve 234 modules, 429 complete-
test entries, zero target load, and every unrelated selection; they manufacture
only a direct root/one graph edge and respectively one or two checksum lines.
All tidy projections converge to selected v1.0.0 and the base projection. No
projection was applied.

The disposable exact v1.0.0 get produces `go.mod`/`go.sum` SHA-256
`bb556bf7d253eb791848db7d863669f60bc0bfe4965b33315f31f78f4f312702` /
`fac428d90355fe794c00ed7cbd9140c6d1fc67b67f6b4c9fdda0b637a40c2a18`;
v1.1.0 produces
`e7af709fbeaadd1e32d7d75aae6c68c5a79bce7580296c8a0ddee609e9b15531` /
`d3c5f9d9162f4fb81575902117b92ef7459be79c95d2af338005004367d8151c`.
All three tidy projections have hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Exact Go 1.26.7 isolated closure measurement is one selected module, one
synthetic graph edge, 44 production and 124 complete-test entries; Go 1.18.10
is one module, zero graph edges, 27 production and 79 complete-test entries.
Both have zero external module-backed packages.

Fresh primary data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no Errwrap record. Exact candidate
OSV responses are empty. Go 1.26.7 source scans have no module, package,
symbol, or test-symbol findings; Go 1.18.10 reports only its old standard-
library findings. Base and v1.1.0 project scans are byte-identical in every
mode with zero target package, symbol, test-symbol, or reachable trace.

The independent fixture file-list receipt SHA-256 is
`6eb0df6dedb99aaf6d00c9e3b9b9771543ed09f394aa8f509c0184eafc884765`.
The 322-entry disposable evidence manifest SHA-256 is
`181ef08a5ce2f8132d9bab9b080bb5cdfb53aedef14f0e87f0ad78e152853cc8`.

The project remains unchanged at 234 modules, 3,599 graph edges, 429 complete-
test entries, 197 module-backed packages across 41 modules, 1,067 sum lines,
and the recorded 432-line tidy projection. `go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Every earlier guard remains valid. Exact-Go applicable gates pass; contained
Go 1.18 retains only the two accepted Darwin shell wording failures. One
preflight launcher invocation reproduced the known signal/log-retention timing
race, and its immediate independent rerun passed all 62 controls. Accepted
quality remains 27/27 Q0-Q2 PASS at L2. No dependency implementation or
metadata commit was created.
