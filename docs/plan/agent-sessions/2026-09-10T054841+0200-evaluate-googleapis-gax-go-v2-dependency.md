# Agent Session: Evaluate Googleapis GAX Go V2 Dependency

Status: NEXT
Session ID: `2026-09-10T054841+0200-evaluate-googleapis-gax-go-v2-dependency`
Created: `2026-09-10T05:48:41+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b97c1fcee29dd1b83b23ada7c9db928847138d22b625f4cd44a858c31567e9d1`
Previous: [2026-09-10T015335+0200-evaluate-google-uuid-dependency.md](2026-09-10T015335+0200-evaluate-google-uuid-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/googleapis/gax-go/v2 v2.7.0` as one bounded dependency group.
Resolve its complete repository and release identity, full Go-floor closure,
package behavior and public API, actual project loading, exact MVS effects,
and every applicable quality contract. Retain or select only an exact-path
version whose complete minimal closure preserves Go 1.18 and whose relevant
behavior passes every contract, except for the explicit bounded inherited
GAX v2.7.0 floor exception authorized below.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof remains retained at
`v0.0.0-20210720184732-4bb14d4b1be1`. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW also remain retained. All earlier
decisions and lifecycle ancestry are final. Do not revisit them or combine
another dependency group. P8 remains queued.

Project MVS selects GAX v2.7.0 through
`github.com/spf13/viper v1.15.0`. The shortest current graph path is
main -> Viper v1.15.0 -> GAX v2.7.0. Additional historical Cloud Go,
Google API, Firestore, Storage, BigQuery, Pub/Sub, Viper, and Crypt vertices
declare v2.0.4, v2.0.5, v2.1.0, or v2.1.1. `go mod why -m
github.com/googleapis/gax-go/v2` is negative, and the complete project load
contains no GAX package. Resolve this ancestry precisely. Do not combine,
upgrade, downgrade, remove, or independently audit Viper, Cloud Go, Google
API, Crypt, or another dependency group. Candidate MVS movements caused by
the exact GAX edge remain in scope to measure, but not to broaden into
independent audits.

A minimal post-UUID survey finds 42 stable proxy versions: v2.0.0 through
v2.0.5; v2.1.0, v2.1.1; v2.2.0 through v2.5.0; v2.5.1 through v2.9.0;
v2.9.1 through v2.12.0; v2.12.1 through v2.12.5; v2.13.0; v2.14.0 through
v2.14.2; and v2.15.0 through v2.24.1. Selected v2.7.0 is dated
2022-11-02T20:02:53Z at origin commit
`2592e2286ac291a2d9e7c06e1f31b3a6d772131c`, with source/mod sums
`h1:IcsPKeInNvYi7eqSaDjiZqDDKu5rsmunY0Y1YupQSSQ=` /
`h1:TEop28CZZQ2y+c0VxMUmu1lV+fQx57QpBWsYpwqHJx8=`. Its module file
declares Go 1.19 and requires go-cmp v0.5.9, Google API v0.102.0, a Genproto
pseudo-version, gRPC v1.50.1, and Protobuf v1.28.1 plus indirect requirements.

V2.5.1 is the last nearby surveyed release declaring Go 1.18; it is dated
2022-08-04 at origin commit
`dfa3344`. V2.6.0 and selected v2.7.0 already declare Go 1.19. Proxy latest
v2.24.1 is dated 2026-09-03T20:21:21Z at origin commit
`269185f57eafcc619f159ffe8857eee6073e955e`, has source/mod sums
`h1:AtqTN21IXMMWo99LiEVAiBfNNQmO40d8xUfZI640mc0=` /
`h1:bWeBei0NVwaNZKb2y1HUBS7gLXIF3/Tu3pq7j8D2Tb0=`, and declares Go
1.25 with a materially newer closure. Treat every incoming fact only as a
survey to verify.

The selected Go 1.19 directive conflicts with the repository's Go 1.18 floor,
while the apparent v2.5.1 floor-compatible release cannot defeat Viper's
v2.7.0 MVS requirement by an exact root downgrade. Resolve this boundary
explicitly. Do not hand-edit or exclude/replace the module, and do not audit
or change Viper to force a lower selection. Determine whether any exact-path
stable upgrade has a fully Go-1.18-compatible source/test closure; if none
does, prove whether the existing unloaded v2.7.0 selection must be retained
as an inherited floor exception or whether the roadmap contract requires a
stop for product direction. Do not manufacture a dependency commit for an
inert decision.

# Authorized Product Decision

On 2026-09-12 the user explicitly selected option 1: retain the existing,
unloaded GAX v2.7.0 selection as a bounded inherited Go-floor exception.
Keep the main module's Go 1.18 floor and do not downgrade GAX, change or
independently audit Viper or another parent, raise the Go floor, add an
explicit root GAX requirement, or manufacture a dependency commit. The prior
runs already proved that v2.5.1 is the highest Go-1.18-declaring release,
that every v2.6.0-or-newer stable release declares at least Go 1.19, that an
exact v2.5.1 request downgrades Viper and multiple unrelated graph
selections, and that no GAX package loads in the project. Treat that
specific floor conflict as resolved by product direction; do not stop or ask
for the same decision again. Finish the remaining repository, release,
closure, API, behavior, concurrency, platform, vulnerability, and applicable
project qualification. Retain v2.7.0 without metadata changes if no separate
disqualifying evidence appears. All other stop conditions remain in force.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path stable candidate. Do not silently promote an unreleased
commit, arbitrary branch head, redirect, fork, alternate path,
floor-ineligible version, prerelease, or tag that does not version this
module.

# Measurements At Start

The latest dependency implementation is exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod`
and `go.sum` with three insertions and no deletions. Google pprof and Gogo
Protobuf remain retained without dependency commits or metadata edits.

Accepted project measurements are 234 selected modules, 3,599 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,067 `go.sum` lines, and a 432-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 51
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,393 module records, modified
2026-09-09T18:12:50Z. Accepted project populations remain 30 module findings,
22 Darwin package findings, 23 Windows package findings, and 20 IDs/22
reachable traces for both Darwin and Windows symbol scans. UUID has no exact
advisory record, package finding, or trace. Do not attribute inherited
findings to GAX without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with seven
improved and zero held, regressed, non-comparable, or dirty counts. Scorecard
SHA-256 is
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
Google UUID decision-summary SHA-256 is
`522b3fc0387507a9612daab87242625d2bfde7aabd74a45ebd24b058d883c31b`;
its 3,206-entry selected-evidence manifest SHA-256 is
`0796b997d528e8b9d25d35cbd90a04cc93bf89fbbb3cdacc27fcd46546d00c6d`.

Read the answered Google UUID archive and rolling handover for its complete
identity, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen UUID, Renameio, pprof, Martian, Viper, or earlier
groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify
binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
inject no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools
beneath scratch as required. Portable receipts remain golangci-lint 2.12.2
archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Docker timestamps with more than six fractional digits require the real
Python 3.14 interpreter ahead of `/usr/bin/python3`; preserve this as an
external execution control.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module and
package/test closure for selected and every serious candidate under exact Go
1.26.7 and contained Go 1.18.10. Inspect imported source and test dependencies
rather than treating module Go directives as full closure proof. Separate
isolated source-time resolution from the project's selected graph, and report
both each module's declared floor and the floor of every dependency in its
minimal production and test closure.

Inspect every package and exported API. Characterize Invoke, CallOption,
retry/backoff, timeout, deadline, and cancellation behavior; gRPC code
matching and HTTP status/error adaptation; APIError wrapping, unwrapping,
Details, Reason, Domain, metadata, and nil/malformed inputs; header insertion
and metadata merging; sleep and clock behavior; deterministic delay
calculation, overflow, negative durations, jitter/randomness, and retry
termination; context propagation; concurrent reuse; zero values; and
platform differences. Distinguish runtime libraries, internal helpers,
commands, examples, benchmarks, testdata, generated content, fuzz/property
coverage, and upstream CI.

Add independent fixtures where useful for retry sequences, no-retry paths,
context cancellation and deadline precedence, timeout composition, backoff
growth/capping/jitter, gRPC and HTTP error conversion, wrapped status/detail
round trips, metadata/header behavior, option ordering, nil functions and
panic boundaries, deterministic inputs, concurrent invocations, and
compatibility with the selected project graph. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every clock, timer,
randomness, network, platform, resource, or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why GAX exists in MVS
while no package is loaded, and preserve every unrelated module selection.
Any change outside the exact GAX edge and its necessary authorized MVS
projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject, retain, or stop if canonical
identity, release qualification, complete floor, behavior, concurrency,
timers/randomness, tests, API, loading, MVS, or any applicable quality
contract fails. Do not treat absence from the loaded package set as permission
to ignore source-time floor or closure evidence.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered UUID archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release,
prerelease, non-versioning tag, or unreleased commit. If no higher release
qualifies, retain selected without hand-editing metadata or manufacturing a
dependency commit. The GAX floor conflict and unavailable bounded downgrade
are already proven; apply the authorized inherited exception above rather
than stopping or requesting product direction again.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run applicable dependency, closure, project, compatibility,
vulnerability, empty-HOME, and cleanliness gates without manufacturing
activity. Full changed-selection quality must preserve 27/27 Q0-Q2 PASS at L2
with zero held, regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Googleapis GAX Go v2 decision, rewrite the rolling handover and
roadmap, answer this archive, create exactly one reciprocal NEXT archive for
the next single P7 group, replace only launcher mutable regions, run
launcher/handoff contracts, and make the normal
`docs: prepare next agent session` commit. Do not implement the next group,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
