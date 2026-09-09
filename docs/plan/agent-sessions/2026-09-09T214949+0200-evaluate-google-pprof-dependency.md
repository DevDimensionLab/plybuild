# Agent Session: Evaluate Google PProf Dependency

Status: NEXT
Session ID: `2026-09-09T214949+0200-evaluate-google-pprof-dependency`
Created: `2026-09-09T21:49:49+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `908d0f711523a60fe3c06e9926ce542281f3128fde082a0881d90dc0a8ed5397`
Previous: [2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md](2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md)
Next: none
Outcome: pending

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
