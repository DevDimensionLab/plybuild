# Agent Session: Decide Hashicorp Errwrap Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-14T234716+0200-decide-hashicorp-errwrap-product-direction`
Created: `2026-09-14T23:47:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3948644c2d85be123cf06ea7167366b78c21aeae3674bb1e09bf1d027db4e3fa`
Previous: [2026-09-14T224709+0200-evaluate-hashicorp-errwrap-dependency.md](2026-09-14T224709+0200-evaluate-hashicorp-errwrap-dependency.md)
Next: [2026-09-15T212106+0200-evaluate-hashicorp-go-cleanhttp-dependency.md](2026-09-15T212106+0200-evaluate-hashicorp-go-cleanhttp-dependency.md)
Outcome: Recorded the user's bounded option 1 decision to retain exact inherited, unloaded Errwrap v1.0.0 under target-specific guards without dependency metadata changes, then prepared the next bounded P7 mission.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-15 explicit selection of
option 1 for exact-path `github.com/hashicorp/errwrap`. Retain exact selected,
inherited, unloaded v1.0.0 without dependency metadata changes under the
Errwrap-specific exception below. Revalidate only the decision guards and
record exactly this choice. Do not repeat the audit, implement a dependency
change, evaluate another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is recording the user's bounded Errwrap option 1
decision after exact Go 1.26.7, accepted dependency moves through Google UUID
v1.4.0, and all recorded retained-module decisions through Consul SDK v0.8.0.
All earlier outcomes and lifecycle ancestry are final. The user has explicitly
chosen to retain exact inherited, unloaded Errwrap v1.0.0 through its sole
selected-version go-multierror v1.1.0 incoming edge. This session may record
that choice and prepare one bounded follow-up, but may not implement a
dependency change or combine another group. P8 remains queued.

Every prior exception remains target-specific. Consul SDK v0.8.0, Consul API
v1.18.0, Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC middleware v1.0.0,
Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 remain final only under their exact-selection, recorded
selected-version incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards before recording a choice and stop for
the owning decision if any expires. Do not change either Gateway parent,
mvn-pom-mutator, GoConvey, Viper, the historical Consul API parent, or transfer
an earlier exception to Errwrap.

# Measurements At Start

The evaluation began from clean HEAD
`94a4ad292427754755374f887c766f6bfa72021e`, parent
`2346fcecd1f9f464697d51c3fa223b6709e84972`, tree
`ed6adadcc242ab335dbb60f37844f378a1a02f89`. The latest dependency
implementation remains Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`. No later retained group has a
dependency implementation or metadata commit.

Canonical proxy, sumdb, `go-import`, Git, and GitHub evidence resolves public,
active, unarchived, non-fork `https://github.com/hashicorp/errwrap`, MPL-2.0,
no GitHub Release objects, no retractions or module deprecation, and no `/v2`
line. The proxy exposes only stable v1.0.0 and v1.1.0. Both have stdlib-only
complete source/test closures that preserve Go 1.18, so v1.1.0 is the highest
floor-eligible stable candidate. Unreleased master declares Go 1.24 and was
not promoted.

Selected v1.0.0 is lightweight tag/commit
`8a6fb523712970c966eefc6b39ed2c5e74880354`, parent
`d6c0cd88035724dd42e0f335ae30161c20575ecc`, tree
`9863613ad8fe960290d1631d84ede660af5f0746`, dated
2018-08-24T00:39:10Z. V1.1.0 is lightweight tag/verified merge commit
`7b00e5db719c64d14dd0caaacbd13e76254d02c0`, parents v1.0.0 and
`96a78ad11c51762df122b738e6f4f30f58e03d8d`, tree
`aefd62cf8e9549e154a65afb8b4524b7b6ed5f2e`, dated
2020-07-14T15:51:01Z. Proxy and Git source match, strict Git verification
passes, and v1.0.0 is an ancestor of v1.1.0.

Each release has one package, one production file, one test file, and no
commands, examples, benchmarks, fuzz targets, testdata, generated files,
build-tag/platform branches, cgo, embeds, go:generate directives, symlinks, or
external dependencies. Both verify, build, pass native count-one, two
count-ten repeats, race, vet, and broad production/test cross-compilation under
exact Go 1.26.7 and contained Go 1.18.10. Their exported declaration sets are
compatible; v1.1.0 adds standard single-error unwrapping and deprecates
`Wrapf`.

Neither stable release qualifies. Both implement exported concrete-type
matching with `reflect.Type.String()`, so distinct types in different import
paths with the same package/type spelling falsely match. Both fail to traverse
standard `Unwrap() []error` graphs such as `errors.Join` under Go 1.26.7;
v1.0.0 additionally lacks standard `Unwrap() error` interoperability.
Independent fixtures reproduce the defects and characterize deterministic
format/order and identity, nil lookup, explicit nil outer-error and callback
panics, aliasing, result-slice isolation, immutable concurrent reads,
allocation, absence of global state/resources, and unguarded recursion.

Selected v1.0.0 has exactly one selected-version incoming edge from
`github.com/hashicorp/go-multierror v1.1.0`; a historical multierror v1.0.0
edge also requests it. Its shortest graph path is main -> mvn-pom-mutator
v0.2.3 -> historical Viper v1.10.1 -> Serf v0.9.6 -> go-multierror v1.1.0 ->
Errwrap v1.0.0. Its why result is negative, repository imports are zero, zero
target packages load, and it is runtime-unreachable. Disposable exact gets
manufacture only a direct root, one graph edge, and target checksums while
preserving zero load and every unrelated selection; tidy returns both to
selected v1.0.0 and the base projection. Nothing was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no Errwrap record. Exact v1.0.0
and v1.1.0 OSV responses are empty. Direct source and project scans have zero
Errwrap module, package, symbol, test-symbol, or reachable-trace finding; Go
1.18 source scans report only that legacy SDK's standard-library age.

The unchanged project has 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed packages across 41 loaded modules,
1,067 sum lines, and the recorded 432-line unapplied tidy projection. Its
`go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Applicable gates pass, all earlier guards remain valid, and accepted quality
remains 27/27 Q0-Q2 PASS at L2. No dependency implementation was created.

# Authorized Product Decision

On 2026-09-15 the user explicitly selected option 1 with the recommended
bounds: retain exact selected `github.com/hashicorp/errwrap v1.0.0` as an
inherited, unloaded selection without changing `go.mod` or `go.sum`. Accept
only the completed concrete-type collision, absent standard single/multi-error
traversal, nil/panic, aliasing, allocation, recursion, API, and related
qualification findings. This does not accept a new or independently discovered
defect.

The exception is Errwrap-specific and non-transferable. It is valid only while
exact v1.0.0 and its sole selected-version incoming edge from
`github.com/hashicorp/go-multierror v1.1.0` remain unchanged, the complete
project load contains zero Errwrap packages, the module remains runtime-
unreachable, and no new advisory or independent defect appears. Revalidate
and record these guards. Direct import/loading, runtime reachability, a target
version or incoming-edge change, or a new advisory or independent defect
expires the exception and requires a fresh Errwrap dependency and product
decision before merge.

Do not add a direct Errwrap edge, select v1.1.0 or another version, remove or
change the historical parent chain, change go-multierror, Serf, Viper or
mvn-pom-mutator, raise the Go floor, authorize replacement or modernization,
move unrelated selections, or manufacture a dependency implementation commit.
Do not stop or ask for this same Errwrap decision again while all guards hold.
Preserve the Consul SDK decision and every earlier separate exception.

# Role And Boundaries

This is a decision-recording session, not a renewed audit or implementation.
The user has explicitly selected and bounded option 1. Reuse the completed
evaluation; do not broaden the choice into direct use, another version,
parent removal, replacement, modernization, a Go-floor change, unrelated-
module authorization, another dependency group, or P8.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status,
reciprocal archive history, latest dependency implementation identity, P7/P8
state, every earlier guarded invariant, current Errwrap exact selection and
incoming edge, negative why, zero imports/load, runtime unreachability, current
advisory state, project hashes, and `./codex-dev-start.sh --check`. Read this
archive, the answered Errwrap evaluation, rolling handover, roadmap, go.mod,
go.sum, and referenced lifecycle contracts. Reuse the completed audit; do not
repeat or broaden it.

# Three Moves

First, revalidate only exact selection, incoming edge, negative why, zero
target and earlier guarded package loads, project hashes, runtime
unreachability, and current advisory state; reuse the completed audit and do
not broaden it. Stop for the owning decision if any guard expires. Second,
record exactly the user's option 1 exception, accepted findings, expiry
triggers, and no-change decision result in the roadmap and rolling handover.
Third, answer this archive and prepare exactly one reciprocal NEXT mission
for the next bounded P7 group without executing it.

# Automatic Handoff

After recording an explicit choice, run applicable no-change lifecycle gates
and make the required local `docs: prepare next agent session` commit. Do not
implement an option, launch a successor, push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, combine another dependency group,
or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Session Outcome

On 2026-09-15 the user explicitly selected option 1 with the recommended
bounds. Exact selected `github.com/hashicorp/errwrap v1.0.0` remains inherited
and unloaded without changing `go.mod` or `go.sum`. The Errwrap-specific,
non-transferable exception accepts only the completed concrete-type collision,
absent standard single/multi-error traversal, nil/panic, aliasing, allocation,
recursion, API, and related qualification findings. It does not accept any new
or independently discovered defect.

Guard-only revalidation ran from clean decision HEAD
`f77cb8a75c6813b28656c85b654d17d7df464731`, parent
`a2f1bc98dc6b978d802d42f7c36bae96627bda33`, tree
`9dd5c04d0cc3cdebf482b831faef8809ea25643e`, under exact Go 1.26.7 binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Errwrap remains exact v1.0.0 through its sole selected-version incoming edge
from `github.com/hashicorp/go-multierror v1.1.0`; the historical
go-multierror v1.0.0 graph vertex also continues to request it. All ten
guarded versions and recorded incoming edges remain unchanged. All ten
`go mod why -m` results remain negative, repository Go source contains zero
guarded-path occurrences, and production and complete-test loads contain zero
guarded packages. The unchanged complete-test load has 429 entries, including
197 module-backed entries across 41 loaded modules. All guarded modules remain
runtime-unreachable.

The project remains 234 selected modules, 3,599 graph edges, and 1,067
`go.sum` lines. `go.mod` and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`. No dependency implementation
or metadata change was created.

Fresh primary vulnerability data now has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. It has no Errwrap or other new
guarded-target record. Exact Errwrap v1.0.0 and v1.1.0 OSV responses remain
empty. Gorilla retains only GO-2020-0019 and unwithdrawn GO-2026-6278; the
latter's record SHA-256 remains
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.
No new advisory or independently observed defect appeared.

The exception remains valid only while exact Errwrap v1.0.0 and its sole
selected-version go-multierror v1.1.0 incoming edge remain unchanged, zero
Errwrap packages load, the module remains runtime-unreachable, and no new
advisory or independent defect appears. Direct import/loading, runtime
reachability, a target version or incoming-edge change, or a new advisory or
independent defect expires the exception and requires a fresh Errwrap
dependency and product decision before merge. Every earlier exception remains
separate. The reciprocal successor is the bounded Hashicorp go-cleanhttp
evaluation; it was prepared but not executed.

Exact Go 1.26.7 module verification, build, count-one tests, race tests, and
vet pass. The launcher lifecycle check recognizes the reciprocal successor.
No changed-selection scorecard applies, and accepted quality remains 27/27
Q0-Q2 PASS at L2.
