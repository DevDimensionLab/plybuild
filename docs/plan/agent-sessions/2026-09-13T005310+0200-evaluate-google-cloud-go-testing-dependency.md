# Agent Session: Evaluate Google Cloud Go Testing Dependency

Status: NEXT
Session ID: `2026-09-13T005310+0200-evaluate-google-cloud-go-testing-dependency`
Created: `2026-09-13T00:53:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5170e6ae91da384631db67b2dc5281bc244146a255b43dae4d43f220421d3dbf`
Previous: [2026-09-10T054841+0200-evaluate-googleapis-gax-go-v2-dependency.md](2026-09-10T054841+0200-evaluate-googleapis-gax-go-v2-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/googleapis/google-cloud-go-testing
v0.0.0-20200911160855-bcd43fbb19e8` as one bounded dependency group.
Resolve its repository and release identity, complete Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
and every applicable quality contract. Retain or select only an exact-path
release whose complete minimal source/test closure preserves Go 1.18 and whose
relevant behavior passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof remains retained at
`v0.0.0-20210720184732-4bb14d4b1be1`. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, and Googleapis GAX Go v2 v2.7.0 also
remain retained. All earlier decisions and lifecycle ancestry are final. Do
not revisit them or combine another dependency group. P8 remains queued.

GAX v2.7.0 is retained without dependency metadata changes under the user's
two explicit bounded decisions. Its inherited Go-1.19 floor exception and
known `apierror.ParseError(err, false)` panic exception remain valid only
while the complete project load contains zero GAX packages. If a future
checkpoint loads or directly imports GAX, that checkpoint must stop for a
fresh dependency and product decision before merge. Do not reopen GAX.

Project MVS selects Google Cloud Go Testing through
`github.com/spf13/afero v1.9.4`. The shortest current graph path is main ->
Afero v1.9.4 -> Google Cloud Go Testing selected. Historical Afero v1.8.2
declares the same selected pseudo-version. `go mod why -m` is negative, and
the complete project load contains no package from this module. Resolve this
ancestry precisely. Do not combine, upgrade, downgrade, remove, or
independently audit Afero, Cloud Go, BigQuery, Datastore, Pub/Sub, Storage,
Google API, or another dependency group. Candidate MVS movements caused by
the exact Google Cloud Go Testing edge remain in scope to measure, but not to
broaden into independent audits.

A fresh minimal survey found an empty stable proxy version list. Selected is
the pseudo-version above, dated 2020-09-11T16:08:55Z at commit
`bcd43fbb19e8`, with source/mod sums
`h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`. Proxy latest is
`v0.0.0-20210719221736-1c9a4c676720`, dated 2021-07-19T22:17:36Z at commit
`1c9a4c676720`, with source/mod sums
`h1:zC34cGQu69FG7qzJ3WiKW244WfhDC3xxYMeNOX2gtUQ=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`. Both module files
declare Go 1.11 and the same Cloud Go v0.44.3, BigQuery v1.0.1, Datastore
v1.0.0, Google API v0.9.0, and tool requirements. The minimal archive survey
observed only README changes and a new SECURITY.md in latest; verify rather
than trust every incoming fact.

An empty proxy version list does not qualify an arbitrary branch head or
pseudo-version as a release. Resolve repository tags, releases, branches,
archive/deprecation status, commits, ancestry, and the exact origin of both
pseudo-versions. Do not silently promote proxy latest, an unreleased commit,
arbitrary branch head, redirect, fork, alternate path, floor-ineligible
version, prerelease, or tag that does not version this module. If the module
has no qualified exact-path release newer than selected, retain selected
unless an independent disqualifier requires product direction; do not create
a dependency commit for an inert decision.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod`
and `go.sum` with three insertions and no deletions. GAX, Google pprof, and
Gogo Protobuf remain retained without dependency commits or metadata edits.

Accepted project measurements remain 234 selected modules, 3,599 graph
edges, 429 native complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,067 `go.sum` lines, and a 432-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, metadata adds
exactly 51 checksum lines and removes zero. The main module retains Go 1.18
and toolchain Go 1.26.7.

Fresh GAX-session primary vulnerability data has 1,398 module records, index
Last-Modified 2026-09-10T16:28:28Z, and scanner DB update
2026-09-10T14:48:42Z. Accepted project populations remain 30 module
findings, 22 Darwin package findings, 23 Windows package findings, and 20
IDs/22 reachable traces on both Darwin and Windows symbol scans. No GAX
module, package, or trace enters the project. Do not attribute inherited
findings to Google Cloud Go Testing without exact evidence.

The accepted quality result remains all 27 Q0-Q2 rows PASS at L2, with seven
improved and zero held, regressed, non-comparable, or dirty counts. Scorecard
SHA-256 is
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
GAX decision-summary SHA-256 is
`26a2a6a7ad4d06a2c059959d7eaf8df64ee9baeb1688c03927b08b1fd650bfe1`;
its 863-entry selected-evidence manifest SHA-256 is
`8b8ea0d9b842b8b3eecc1d7c7b8f8754c24662a74dd21d9d1cef0304bbd18b5b`.

Read the answered GAX archive and rolling handover for its complete identity,
closure, API, behavior, MVS, vulnerability, quality, exceptions, and evidence
record. Do not reopen GAX, UUID, Renameio, pprof, Martian, Afero, or earlier
groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
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

Inspect every package and exported API. The initial source tree contains the
root documentation package plus BigQuery `bqiface`, Datastore `dsiface`,
Pub/Sub `psiface`, and Storage `stiface` adapters/interfaces. Characterize
all wrappers, forwarding, embedded interfaces, constructors, return values,
nil/zero inputs, nil embedded clients, typed nils, aliasing, iterator and
stream behavior, callbacks, context propagation and cancellation, errors and
panics, concurrent reuse, generated/manual content, examples, platform
differences, and compatibility with their pinned Cloud client APIs.
Distinguish runtime packages, tools-only requirements, internal helpers,
commands, examples, benchmarks, testdata, fuzz/property coverage, and upstream
CI. Do not infer implementation guarantees from interface declarations alone.

Add independent fixtures where useful for exact method delegation, argument
and context identity, return/error preservation, interface satisfaction,
nil/panic boundaries, iterator/page behavior, stream receive/send behavior,
callback invocation, mock substitution, concurrency, deterministic inputs,
and compatibility with the selected project graph. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every clock, timer,
randomness, network, emulator, credential, platform, resource, or test-design
failure precisely. Do not require live cloud services for unit qualification;
classify integration prerequisites and keep them contained.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why Google Cloud Go
Testing exists in MVS while no package is loaded, and preserve every unrelated
module selection. Any change outside the exact target edge and its necessary
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject, retain, or stop if canonical
identity, release qualification, complete floor, behavior, concurrency,
tests, API, loading, MVS, or any applicable quality contract fails. Do not
treat absence from the loaded package set as permission to ignore source-time
floor, closure, behavior, or vulnerability evidence.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered GAX archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path release. Do not promote a
redirect, fork, alternate path, floor-ineligible version, prerelease,
non-versioning tag, or unreleased commit. In particular, proxy `@latest` does
not itself qualify the latest pseudo-version as a release. If no higher
release qualifies, retain selected without hand-editing metadata or
manufacturing a dependency commit.

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

After the Google Cloud Go Testing decision, rewrite the rolling handover and
roadmap, answer this archive, create exactly one reciprocal NEXT archive for
the next single P7 group, replace only launcher mutable regions, run
launcher/handoff contracts, and make the normal
`docs: prepare next agent session` commit. Do not implement the next group,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
