# Agent Session: Evaluate Google PProf Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T214949+0200-evaluate-google-pprof-dependency`
Created: `2026-09-09T21:49:49+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `908d0f711523a60fe3c06e9926ce542281f3128fde082a0881d90dc0a8ed5397`
Previous: [2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md](2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md)
Next: [2026-09-09T232057+0200-evaluate-google-renameio-dependency.md](2026-09-09T232057+0200-evaluate-google-renameio-dependency.md)
Outcome: Retained exact-path Google pprof at the selected 2021 pseudo-version without dependency metadata edits. No semver release exists; proxy latest is an unreleased protected-main head, requires Go 1.25.0, and is API-incompatible. Complete closure, behavior, MVS/loading, project, vulnerability, and applicable quality gates are resolved.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1` as one bounded
dependency group. Resolve its complete repository and pseudo-version identity,
full Go-floor closure, package behavior and public API, actual project loading,
exact MVS effects, and every applicable quality contract. Retain or select only
an exact-path version whose complete minimal closure preserves Go 1.18 and
whose relevant behavior passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, and Google Martian v3.3.2 moves. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW remain retained. All earlier decisions
and lifecycle ancestry are final. Do not revisit them or combine another
dependency group. P8 remains queued.

Project MVS selects Google pprof
v0.0.0-20210720184732-4bb14d4b1be1 through the declared edge
`cloud.google.com/go v0.90.0 -> github.com/google/pprof` even though the
project selects Cloud Go v0.105.0. `go mod why -m github.com/google/pprof` is
negative, and the 429-entry complete project package load contains no pprof
package. Resolve this ancestry precisely. Do not combine, upgrade, remove, or
independently audit Cloud Go, Demangle, Readline, Logex, X/Sys, or another
dependency group. Candidate MVS movements caused by the exact pprof edge
remain in scope to measure, but not to broaden into independent audits.

A minimal post-Martian survey finds no tagged versions in the Go proxy list.
Selected is a 2021-07-20T18:47:32Z pseudo-version with source/mod sums
`h1:K6RDEckDVWvDI9JAJYCmNdQXq6neHJOYx3V6jnqNEec=` /
`h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=`. It declares Go 1.14
and requires historical Readline, Logex, Test, Demangle, and X/Sys versions.
Proxy `@latest` is
`v0.0.0-20260906184651-6331bc6350fe`, dated 2026-09-06T18:46:51Z at
origin commit `6331bc6350fe55a6fec2957299e0581dd7510e36`, with source/mod sums
`h1:QAinXoAFJdGQYztXn3VpFey7KCwpedbZ/EkzbplQ0cY=` /
`h1:jl5iWTm0/hd5PjEYEOuwAJ57L/CibdZfrqZ5XA5GrCk=` and a Go 1.25.0
directive, so it is floor-ineligible. Treat every incoming fact only as a
survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path pseudo-version candidate. Do not silently promote an
unreleased commit, arbitrary branch head, redirect, fork, alternate path,
floor-ineligible pseudo-version, prerelease, or tag that does not version this
module. Absence of semver tags is not permission to choose an arbitrary commit.

# Measurements At Start

The latest dependency implementation is exact Google Martian v3.3.2 commit
`4644476aef78e24d0c7e6a1140baba4c3226a1cc`, parent
`8508a9d54f066985fb81ebd35fa31f45fa36da40`, and tree
`f20d9b3817fe29af4876745d09c14e9cbcac9152`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gogo Protobuf v1.3.2 remains
retained without a dependency commit or metadata edit.

Accepted project measurements are 234 selected modules, 3,597 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,063 `go.sum` lines, and a 418-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 47
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted project
populations are 30 Darwin module findings, 22 Darwin package findings, 23
Windows package findings, and 20 IDs/22 reachable traces for both Darwin and
Windows symbol scans. Martian has no exact advisory record; its direct closure
findings are inherited from its historical dependencies, while the project
loads no Martian package or trace. Do not attribute inherited findings to pprof
without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256 `094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.
Google Martian decision-summary SHA-256 is
`374c235a5ebeeb7a97254e8eab57b3dc60224fc7f103611beceafff701f5756e`;
its 509-entry selected-evidence manifest SHA-256 is
`614921201697ce6f3accdf0c6e619c86ddb3502f23c3b03ee6314a2347a9a649`.

Read the answered Google Martian archive and rolling handover for its complete
release, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen Martian, Cloud Go, Snappy, Protobuf, gRPC, X/Net, Go CMP,
or earlier groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
inject no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools
beneath scratch as required. Portable receipts remain golangci-lint 2.12.2
archive `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module and
package/test closure for selected and every serious candidate under exact Go
1.26.7 and contained Go 1.18.10. Inspect imported source and test dependencies
rather than treating a module Go directive as full closure proof. Keep isolated
source-time resolution separate from the project's selected graph.

Inspect every package and exported API. Characterize profile parsing,
validation, serialization, merging, normalization, mapping and sample/value
semantics; generated protocol code; driver configuration; local and remote
profile fetching; gzip and malformed/truncated input; symbolization,
demangling, disassembly, graph/report generation, filtering, tags, units,
interactive commands, web UI and HTTP handlers; temporary files, subprocesses,
browser invocation, network and filesystem boundaries; resource limits,
determinism, concurrency and reuse; build tags and platform behavior; commands,
examples, benchmarks, testdata, fuzz/property coverage, and upstream CI.
Distinguish runtime libraries, generated code, commands, embedded web assets,
and test fixtures.

Add independent fixtures where useful for profile/protocol round trips and
merges, malformed or truncated profiles, gzip streams, mapping and location
identity, sample/value types, labels/tags/units, graph/report determinism,
filters, local/HTTP fetching, symbolization boundaries, short/failing I/O,
resource limits, reuse, concurrency, and compatibility with the selected
project graph. Run source verification, package listing, native complete tests,
two independent repeats, race, vet, and meaningful cross-builds under both
SDKs. Classify every generator, network, subprocess, timing, platform,
resource, or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why pprof exists in MVS
while no package is loaded, and preserve every unrelated module selection. Any
change outside the exact pprof edge and its necessary authorized MVS projection
is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
version qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status, current
ancestry, latest dependency implementation identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
the answered Martian archive, rolling handover, roadmap, `go.mod`, `go.sum`,
and every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path pseudo-version. Do not
promote an arbitrary commit, redirect, fork, alternate path, floor-ineligible
version, prerelease, or non-versioning tag. If no higher pseudo-version has a
defensible release qualification and complete Go-1.18 closure, retain selected
without hand-editing metadata or manufacturing a dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete P7
dependency gate. For an inert/retained selection, prove the no-change effect
and run applicable dependency, closure, project, compatibility, vulnerability,
empty-HOME, and cleanliness gates without manufacturing activity. Full
changed-selection quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held,
regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
run `go mod download all` in a measured worktree or create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Google pprof decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Retain exact-path `github.com/google/pprof` at
`v0.0.0-20210720184732-4bb14d4b1be1` without changing dependency metadata or
manufacturing an implementation commit. The module proxy version list is
empty, the repository has no tags or releases, and no later commit has a
defensible release qualification. Proxy `@latest`,
`v0.0.0-20260906184651-6331bc6350fe`, is merely the protected `main` branch
head, declares Go 1.25.0, and incompatibly changes selected public API. It is
therefore not selectable for the retained Go 1.18 floor.

Proxy, sumdb, `go-import` metadata, Git, and GitHub resolve the module to the
public, active, unarchived, non-fork Apache-2.0 repository
`https://github.com/google/pprof`. Its only branch is protected `main`; no
redirect, retraction, deprecation, alternate root module path, or qualified
fork exists. The current repository also contains the separate nested
`github.com/google/pprof/browsertests` test module, which is not an alternate
path for the root module. GitHub reports valid commit signatures for selected
and latest.

Selected resolves to commit
`4bb14d4b1be14417e47d0bbaf2bd4e188eda647f`, tree
`94d5b610d24eee4637e1bb02a9598326db7edbd5`, parent
`86eeefc3e4714e359ff19990b871874910636148`, at
2021-07-20T18:47:32Z. Latest resolves to commit
`6331bc6350fe55a6fec2957299e0581dd7510e36`, tree
`70fc553307c8471558c9d4b8830efea94ab68735`, parent
`d6c3cb2f37ec22719bbaf5eb031d9a46635cb5b2`, at
2026-09-06T18:46:51Z. Selected is an ancestor of latest, with 237 intervening
commits. Normalized proxy and Git manifests match exactly: selected has 214
files and SHA-256
`4344614ca6a1575f5ab8d45982fbd02534733fa73087dbcffa1b6637a22e3b06`;
latest's 245 root-module files hash to
`2fea39c1fed5aaecef93284693cbc208a27864c6046c7fd521a3b62687cb5ce9`.

Selected source/mod sums are
`h1:K6RDEckDVWvDI9JAJYCmNdQXq6neHJOYx3V6jnqNEec=` and
`h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=`. It declares Go 1.14
and requires the historical Readline, Logex, Test, Demangle, and X/Sys
versions stated in the prompt. Its complete standalone graph contains six
modules. Complete test loading has 260 entries under Go 1.26.7 and 196 under
Go 1.18.10, with 40 module-backed entries and 38 pprof package/test variants
under both SDKs. Only pprof, Readline, and Demangle load; Logex, Test, and
X/Sys remain graph-only. Building and exercising imported sources and tests
under contained Go 1.18.10 proves the closure floor.

The 18 package directories comprise the pprof command, injectable `driver`
API, `profile` model and codec, a legacy go-fuzz adapter, internal binary and
object tooling, graph/measurement/report generation, local/remote fetching,
symbolization, transport and web UI, plus embedded JavaScript/CSS packages.
`profile/proto.go` is a handwritten protobuf-wire codec for
`proto/profile.proto`, not generated `.pb.go` code. Exported profile behavior
covers validation, gzip/uncompressed parsing and writing, copying,
aggregation, normalization, scaling, merging, filtering, labels/tags/units,
mapping/location/function identity, pruning, and proc-map parsing. Driver
exports injectable writers, flags, fetchers, symbolizers, object tools/files,
UI, mappings, frames, instructions, symbols, and HTTP-server arguments.

The independent 153-line fixture validates compressed and uncompressed
round trips, table pointer identity, mappings, locations, functions,
sample/value types, values, string/numeric labels and units, relocated-map
merge identity and summed samples, deep copying, filters, malformed/truncated
profiles and gzip, validation failures, deterministic concurrent writes, and
failing short output. It passes race count-10 under both SDKs and has SHA-256
`3c15f63b822c49eb4cd58789410607467ab26a3cb0bcc90162c02b274a9521ab`.
Built commands produce byte-identical top and DOT reports twice and across
both SDKs, with SHA-256 values
`0ee6392de3c786ced86fc0908b75c6ae34a75604e4375e10e99e3d9ea55d2afc`
and `6e1b0492964398dfcf1f1a72d2abcb40ada803b1355f05d43ea2fd89b39e46e1`.
Actual loopback HTTP fetching, gzip storage, local input, and isolated-HOME
execution pass. Native object-tool fixture coverage substitutes for absent
Graphviz and llvm-symbolizer host binaries.

The selected source has 75 Go files, 26 test files, 116 `Test` functions,
114 testdata files, one legacy go-fuzz corpus, and no native example,
benchmark, or Go fuzz target. Complete native tests pass with package vet
disabled, two independent process repeats pass, and race passes under both
SDKs. Default `go test` and standalone vet expose the same inherited
production `fmt.Fprintln` redundant-newline finding; Go 1.26 additionally
finds two test-only loop-variable captures. Same-process count-2 exposes a
test-design panic because `TestSymbolzAfterMerge` redefines global
`flag.CommandLine`; independent repeats remain green. CGO-disabled Darwin,
Linux, Windows, and FreeBSD builds pass. js/wasm alone fails in the historical
Readline closure because that selected dependency has no JS implementation;
it is classified as an unsupported platform, not a pprof regression.

Selected retains explicit boundary qualifications. Parsers and non-200 HTTP
error handling use whole-buffer reads without an application size limit, so
callers must bound hostile compressed, protobuf, legacy, and response input.
Concurrent remote chunk fetches are capped at 64, but aggregate memory remains
profile-sized. The web server can bind a caller-supplied host, has no
authentication or server read/write timeout, and must not be exposed as an
untrusted service. Saved remote-profile files are not closed; `Profile.Write`
discards deferred gzip-finalization errors. Driver flags, commands, and mode
state are global and not reentrant. Profile encoding state is mutex-protected,
but concurrent caller mutation is unsupported. Optional graph, browser,
viewer, perf conversion, disassembly, and symbolization paths cross explicit
filesystem, network, temporary-file, and subprocess boundaries.

Latest source/mod sums are
`h1:QAinXoAFJdGQYztXn3VpFey7KCwpedbZ/EkzbplQ0cY=` and
`h1:jl5iWTm0/hd5PjEYEOuwAJ57L/CibdZfrqZ5XA5GrCk=`. Its Go 1.25.0
directive cannot even be parsed by Go 1.18.10. Under Go 1.26.7 its six-module,
16-package closure verifies and passes native tests and vet, and it fixes the
selected gzip-close error and saved-file descriptor. Those improvements do
not overcome its missing release qualification or floor. Pinned apidiff also
reports incompatible `Profile.Aggregate` and `driver.ObjTool.Open` signature
changes, removal of the exported d3 and d3flamegraph packages, and conversion
of `svgpan.JSSource` from constant to variable, alongside compatible profile,
mapping, line, frame, and label additions.

Pprof remains selected in project MVS through this exact declaration chain:
main -> `mvn-pom-mutator v0.2.3` -> `viper v1.10.1` -> Firestore v1.6.1 ->
Cloud Go v0.97.0 -> GAX v2.1.0 -> Google API v0.54.0 -> Cloud Go v0.90.0 ->
selected pprof. The graph retains requirements from encountered historical
module vertices even though another path selects Cloud Go v0.105.0. Package
loading traverses imports rather than every module-graph edge, explaining the
negative `go mod why -m` and zero loaded pprof packages.

The retained project stays at 234 modules, 3,597 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed entries, 1,063 sum lines,
and a 418-line unapplied tidy projection. Exact latest would move only pprof,
Demangle, and X/Sys selections, force the main Go line from 1.18 to 1.25.0,
add six sum lines, and produce 234 modules, 3,605 edges, the same loaded
population, and a 440-line tidy projection. This projection was discarded.

Fresh vulnerability data has 1,393 module records, modified
2026-09-09T17:56:34Z, and no exact pprof record. Direct selected and latest
module, package, and symbol scans are all zero. Project base/latest results
are identical: 30 module IDs, 22 Darwin package IDs, 23 Windows package IDs,
and 20 IDs/22 reachable traces on both symbol platforms, with no pprof
finding or trace.

Exact Go 1.26.7 project verification, complete load, build, full tests, two
independent repeats, race, vet, pinned lint, empty-HOME count-2, Linux/Windows
builds, API/CLI compatibility, CLI surface, launcher contracts, and complete
contained preflight pass. The Go 1.18.10 projection removes only the toolchain
line; verification, loading, build, vet, cross-builds, 26-package repeats and
race, and all 31 compatible shell tests pass. Full count-1 retains only the
two accepted closed-file error-wording assertions. Because selection is
inert, changed-selection quality is inapplicable; the accepted 27/27 Q0-Q2 L2
scorecard remains
`094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.

The 188-entry selected-evidence manifest SHA-256 is
`e6090a35f4f669bfff0eccf3b179def466dcbc18645ec61ed5d2e81f92e4f0f7`;
decision-summary SHA-256 is
`20804b312d90e5390df74dde4c1b252e68d8b1183897e7e562f2c8a133ac5d08`.
