# Agent Session: Evaluate Google UUID Dependency

Status: NEXT
Session ID: `2026-09-10T015335+0200-evaluate-google-uuid-dependency`
Created: `2026-09-10T01:53:35+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `2f1357e45eb9ff64c6557ba56a8a5db443b3fea973f2415277a93a23285446f5`
Previous: [2026-09-09T232057+0200-evaluate-google-renameio-dependency.md](2026-09-09T232057+0200-evaluate-google-renameio-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/google/uuid v1.1.2` as one bounded dependency group. Resolve its
complete repository and release identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, and Google Renameio v1.0.1 moves. Google
pprof remains retained at `v0.0.0-20210720184732-4bb14d4b1be1`. Gogo
Protobuf v1.3.2, Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging,
ansimage, Fsnotify, Ghodss YAML, and historical root GLFW also remain
retained. All earlier decisions and lifecycle ancestry are final. Do not
revisit them or combine another dependency group. P8 remains queued.

Project MVS selects UUID v1.1.2 through fourteen historical
`google.golang.org/grpc` vertices spanning v1.33.1 through v1.43.0. One shortest
current graph path is main -> Google Martian/v3 v3.3.2 -> gRPC v1.37.0 ->
UUID v1.1.2. `go mod why -m github.com/google/uuid` is negative, and the
complete project load contains no UUID package. Resolve this ancestry
precisely. Do not combine, upgrade, remove, or independently audit gRPC,
Martian, or another dependency group. Candidate MVS movements caused by the
exact UUID edge remain in scope to measure, but not to broaden into
independent audits.

A minimal post-Renameio survey finds 13 stable proxy versions: v1.0.0,
v1.1.0 through v1.1.5, v1.2.0, v1.3.0, v1.3.1, v1.4.0, v1.5.0, and
v1.6.0. Selected is dated 2020-07-02T18:56:42Z, has source/mod sums
`h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=` /
`h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=`, and its module file has
no Go directive or requirements. Proxy latest v1.6.0 is dated
2024-01-23T18:54:04Z at origin commit
`0f11ee6918f41a04c201eceeadf612a377bc7fbc`, has source/mod sums
`h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=` /
`h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=`, and likewise has no Go
directive or requirements. Treat every incoming fact only as a survey to
verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path stable candidate. Do not silently promote an unreleased
commit, arbitrary branch head, redirect, fork, alternate path,
floor-ineligible version, prerelease, or tag that does not version this
module.

# Measurements At Start

The latest dependency implementation is exact Google Renameio v1.0.1 commit
`394ec36ac5cb6712640024f0a47af28bf3631906`, parent
`29b36a6574e2faabedd55784e2576f051d83f655`, and tree
`a55a902a1f491d4972dc2118752b8426d754d1b0`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Google pprof and Gogo
Protobuf remain retained without dependency commits or metadata edits.

Accepted project measurements are 234 selected modules, 3,598 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,065 `go.sum` lines, and a 428-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 49
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,393 module records, modified
2026-09-09T18:12:50Z. Accepted project populations remain 30 module findings,
22 Darwin package findings, 23 Windows package findings, and 20 IDs/22
reachable traces for both Darwin and Windows symbol scans. Renameio has no
exact advisory record, package finding, or trace. Do not attribute inherited
findings to UUID without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256 `6be51aec069011e99a80c2b5b5097bb79a67abe1c0078589ccd60198ab75c2ce`.
Google Renameio decision-summary SHA-256 is
`07fcaace013997f5d4513de49c5a59b3359b7b270c4e6634b39729c8a45a6096`;
its 7,798-entry selected-evidence manifest SHA-256 is
`4d4fd72b84786270f069b18934eb831d31ecc801e1d348e6e912f756fc114f84`.

Read the answered Google Renameio archive and rolling handover for its
complete identity, closure, API, behavior, MVS, vulnerability, quality, and
evidence record. Do not reopen Renameio, pprof, Martian, gRPC, or earlier
groups.

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
rather than treating module Go directives as full closure proof. Keep
isolated source-time resolution separate from the project's selected graph.

Inspect every package and exported API. Characterize UUID parsing and
formatting for canonical, raw, URN, Microsoft, invalid, short, long, mixed
case, and whitespace inputs; version/variant validation; random, time,
name-hash, DCE, and future-version generation; entropy, clock sequence, node
identity, hardware-interface discovery, fallback, state mutation, and
concurrent generation; nil/zero values, Must/panic boundaries, comparison,
JSON/text/binary/database marshal and scan behavior, aliasing, deterministic
outputs, and Windows/Unix differences. Distinguish runtime libraries,
internal helpers, commands, examples, benchmarks, testdata, generated
content, fuzz/property coverage, and upstream CI.

Add independent fixtures where useful for parse/format round trips and
rejections, versions and variants, deterministic MD5/SHA1 name UUIDs, random
source failures and short reads, timestamp/clock-sequence rollover, node ID
selection and override, global-state reuse, concurrent generators,
serialization and SQL scanning, error and panic contracts, platform behavior,
and compatibility with the selected project graph. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every generator,
randomness, clock, network-interface, platform, resource, or test-design
failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why UUID exists in MVS
while no package is loaded, and preserve every unrelated module selection.
Any change outside the exact UUID edge and its necessary authorized MVS
projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, randomness,
tests, API, loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered Renameio archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release,
prerelease, non-versioning tag, or unreleased commit. If no higher release
qualifies, retain selected without hand-editing metadata or manufacturing a
dependency commit.

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

After the Google UUID decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
