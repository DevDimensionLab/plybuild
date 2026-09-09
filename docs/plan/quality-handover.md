# Quality Upgrade Handover

Generated: 2026-09-09T23:20:57+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Google Martian v3.3.2 commit
  `4644476aef78e24d0c7e6a1140baba4c3226a1cc`, parent
  `8508a9d54f066985fb81ebd35fa31f45fa36da40`, tree
  `f20d9b3817fe29af4876745d09c14e9cbcac9152`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Google
  pprof and Gogo Protobuf remain retained without dependency edits. All
  earlier P7 decisions are final.
- The answered Google pprof archive and sole NEXT Google Renameio archive
  link reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, and Google Martian v3.3.2 moves. Google pprof
`v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW remain retained. P8 remains queued. Do
not reopen earlier groups or combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Google Pprof Decision And Identity

Retain exact-path `github.com/google/pprof` at
`v0.0.0-20210720184732-4bb14d4b1be1`. The proxy list, Git tags, and GitHub
releases are all empty. The only branch is protected `main`, so no later
commit has release qualification. Proxy latest
`v0.0.0-20260906184651-6331bc6350fe` is that unreleased head, declares Go
1.25.0, and is API-incompatible. No dependency edit or implementation commit
was created.

Proxy, sumdb, `go-import`, Git, and GitHub agree on the public, active,
unarchived, non-fork Apache-2.0 repository
`https://github.com/google/pprof`. No deprecation, retraction, redirect,
qualified fork, or alternate root path exists. The nested current
`github.com/google/pprof/browsertests` module is only a separate browser-test
module. GitHub reports valid signatures for selected and latest commits.

Selected is commit `4bb14d4b1be14417e47d0bbaf2bd4e188eda647f`, tree
`94d5b610d24eee4637e1bb02a9598326db7edbd5`, parent
`86eeefc3e4714e359ff19990b871874910636148`, at
2021-07-20T18:47:32Z. Latest is commit
`6331bc6350fe55a6fec2957299e0581dd7510e36`, tree
`70fc553307c8471558c9d4b8830efea94ab68735`, parent
`d6c3cb2f37ec22719bbaf5eb031d9a46635cb5b2`, at
2026-09-06T18:46:51Z. Selected precedes latest by 237 commits. Proxy and Git
manifests match: selected's 214 files hash to
`4344614ca6a1575f5ab8d45982fbd02534733fa73087dbcffa1b6637a22e3b06`;
latest's 245 root-module files hash to
`2fea39c1fed5aaecef93284693cbc208a27864c6046c7fd521a3b62687cb5ce9`.

Selected source/mod sums are
`h1:K6RDEckDVWvDI9JAJYCmNdQXq6neHJOYx3V6jnqNEec=` /
`h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=`. Latest sums are
`h1:QAinXoAFJdGQYztXn3VpFey7KCwpedbZ/EkzbplQ0cY=` /
`h1:jl5iWTm0/hd5PjEYEOuwAJ57L/CibdZfrqZ5XA5GrCk=`.

## Closure, API, Behavior, And Qualifications

Selected declares Go 1.14 and has a six-module standalone graph containing
historical Readline, Logex, Test, Demangle, and X/Sys. Complete test loading
has 260 entries under Go 1.26.7 and 196 under Go 1.18.10, with 40
module-backed entries and 38 pprof variants. Only pprof, Readline, and
Demangle load. Actual build/test execution under contained Go 1.18.10 proves
the complete floor.

The 18 package directories cover the command, injectable driver API, profile
model and handwritten protobuf wire codec, legacy go-fuzz adapter, object and
binary inspection, graph/measurement/report generation, filtering,
local/remote fetch, symbolization, transport, web UI, and embedded assets.
Exported profile behavior includes parse/validation/write/copy, merge,
aggregate/normalize/scale, mapping/location/function and sample/value
identity, tags/labels/units, pruning, filtering, and proc-map parsing.

A 153-line independent profile fixture covers compressed/uncompressed round
trips, identity tables, relocated mapping merges, sample sums, copy/filter,
malformed/truncated/gzip input, validation, deterministic concurrent writes,
and failing output. Race count-10 passes under both SDKs; fixture SHA-256 is
`3c15f63b822c49eb4cd58789410607467ab26a3cb0bcc90162c02b274a9521ab`.
Local and loopback-HTTP command runs pass. Top and DOT output is byte-identical
across repeats and SDKs, hashing to
`0ee6392de3c786ced86fc0908b75c6ae34a75604e4375e10e99e3d9ea55d2afc`
and `6e1b0492964398dfcf1f1a72d2abcb40ada803b1355f05d43ea2fd89b39e46e1`.

Native tests with package vet disabled, two independent repeats, race, and
Darwin/Linux/Windows/FreeBSD builds pass both SDKs. Default test/vet exposes
one inherited production redundant-newline finding; Go 1.26 adds two
test-only loop-variable captures. Same-process count-2 exposes global flag
redefinition in `TestSymbolzAfterMerge`; independent repeats pass. js/wasm
fails only because historical Readline lacks JS implementations. These are
classified legacy vet, test-design/global-state, and unsupported-platform
results rather than behavior regressions.

Retained boundaries are explicit. Profile and HTTP inputs are whole-buffered
without application size limits. Remote fetch concurrency is capped at 64,
but memory remains profile-sized. The unauthenticated web server has no
read/write timeouts. Saved remote files are not closed, gzip finalization
errors are discarded, and global driver flag/config state is not reentrant.
Callers must bound hostile input, constrain web binding, and coordinate
concurrent mutation and external browser/viewer/Graphviz/object-tool
subprocesses.

Latest fixes the saved-file and gzip-close errors, but Go 1.18.10 rejects its
Go 1.25.0 directive. Pinned apidiff additionally finds incompatible
`Profile.Aggregate` and `driver.ObjTool.Open` signatures, removed d3 and
d3flamegraph packages, and `svgpan.JSSource` changing from const to var.
There is no qualified intermediate pseudo-version because no later commit is
tagged or released.

## Project, Vulnerability, And Quality Measurements

Pprof remains in MVS through main -> mvn-pom-mutator v0.2.3 -> historical
Viper v1.10.1 -> Firestore v1.6.1 -> Cloud Go v0.97.0 -> GAX v2.1.0 ->
Google API v0.54.0 -> Cloud Go v0.90.0 -> selected pprof. Requirements on
encountered historical vertices remain graph edges despite selected Cloud Go
v0.105.0. Import loading explains negative `go mod why -m` and zero loaded
pprof packages.

The project remains at 234 modules, 3,597 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, 1,063 sum
lines, and a 418-line unapplied tidy projection. Relative to accepted go-cmp
commit `c314bcb`, sums remain +47/-0. Latest would move only pprof, Demangle,
and X/Sys, force the main Go line to 1.25.0, add six sums, and produce 3,605
edges and a 440-line tidy diff without changing the loaded population.

Fresh vulnerability data has 1,393 records, modified
2026-09-09T17:56:34Z, and no pprof record. Direct selected/latest module,
package, and symbol scans are zero. Base/latest project results are identical:
30 module IDs, 22 Darwin package IDs, 23 Windows package IDs, and 20 IDs/22
reachable traces on both symbol platforms, none involving pprof.

Exact Go 1.26.7 project verification/load/build/tests, independent repeats,
race, vet, pinned lint, empty-HOME count-2, Linux/Windows builds, API/CLI and
CLI-surface compatibility, launcher checks, and complete contained preflight
pass. The Go 1.18.10 projection removes only the toolchain line and passes
verification, load, build, vet, cross-builds, the 26-package repeat/race
population, and all 31 compatible shell tests; full count-1 retains the two
accepted closed-file wording differences. Changed-selection quality is
inapplicable. The accepted 27/27 Q0-Q2 L2 scorecard remains
`094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.

The 188-entry selected-evidence manifest SHA-256 is
`e6090a35f4f669bfff0eccf3b179def466dcbc18645ec61ed5d2e81f92e4f0f7`;
decision-summary SHA-256 is
`20804b312d90e5390df74dde4c1b252e68d8b1183897e7e562f2c8a133ac5d08`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Official golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 binary SHA-256 remains
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
  Rebuilt apidiff has exact required module/version identity; govulncheck is
  v1.7.0.

## Next Objective

Independently evaluate exact-path `github.com/google/renameio v0.1.0` as the
next single P7 group. Do not combine HTools, Crypt, Firestore, Cloud Go, or
another dependency.

Three historical HTools vertices declare v0.1.0. One shortest declaration
path is main -> mvn-pom-mutator v0.2.3 -> Crypt -> Firestore v1.1.0 -> Cloud
Go v0.46.3 -> HTools v0.0.1-2019.2.3 -> Renameio. `go mod why -m` is
negative and no Renameio package loads.

The initial proxy list is exactly v0.1.0, v1.0.0, and v1.0.1. Selected is
dated 2019-01-09T16:53:11Z, contains no Go directive/requirements, and has
source/mod sums `h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=` /
`h1:KWCgfxg9yswjAJkECMjeO8J8rahYeXnNhOm40UhjYkI=`. Latest v1.0.1 is
dated 2021-04-06T14:11:08Z, declares Go 1.13 with no requirements, and has
sums `h1:Lh/jXZmvZxb0BBeSY5VKEfidcbcbenKjZFzM/q0fSeU=` /
`h1:t/HQoYBZSsWSNK35C6CO/TpPLDVWvxOHboWUAweKUpk=`. A minimal exact
projection changes only Renameio, adds one graph edge and two sum lines, and
does not alter module/package populations. Treat all survey facts as incoming
evidence to verify.
