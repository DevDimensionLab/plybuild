# Agent Session: Evaluate Golang Snappy Dependency

Status: NEXT
Session ID: `2026-09-09T162826+0200-evaluate-golang-snappy-dependency`
Created: `2026-09-09T16:28:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `51ecc1f9b8e42db0ceee9ca8fe20251268cc72bb8d291b17dffeba90aac59a21`
Previous: [2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md](2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/golang/snappy v0.0.3` as one bounded dependency group. Resolve its
complete release and repository identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, and Golang Protobuf v1.5.3 moves.
Gogo Protobuf v1.3.2, Crypt, OpenCensus Proto, Logex, Readline, Fnmatch,
Imaging, ansimage, Fsnotify, Ghodss YAML, and historical root GLFW remain
retained. All earlier decisions and lifecycle ancestry are final. Do not
revisit them or combine another dependency group. P8 remains queued.

Project MVS selects Golang Snappy v0.0.3 through the declared edge
`github.com/google/martian/v3 v3.2.1 -> github.com/golang/snappy v0.0.3`.
`go mod why -m github.com/golang/snappy` says the main module does not need it,
and the complete project package load contains no Snappy package. Do not
combine, upgrade, remove, or independently audit Google Martian or another
dependency group. Candidate MVS movements caused by the exact Snappy edge
remain in scope to measure, but not to broaden into independent audits.

A minimal post-Golang-Protobuf survey finds five proxy versions: v0.0.1,
v0.0.2, selected v0.0.3, v0.0.4, and latest v1.0.0, with no prerelease.
Selected v0.0.3 is dated 2020-11-03T22:46:00Z at verified commit
`674baa8c7fc30da5df3074a459494a7e6b427dff`, tree
`5bbf3ddf1162e96e96a3019f84811da0b2416ee8`, and has source/mod sums
`h1:fHPg5GQYlCeLIPB9BZqMVR5nR9A+IM5zcgeTdjMYmLA=` /
`h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=`. Latest v1.0.0 is
dated 2023-12-25T22:57:46Z at verified commit
`43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, tree
`c6f19681f2ae79c2e55b6c6113e9b5ea183e3bcb`, and has source/mod sums
`h1:Oy607GVXHs7RtbggtPBnr2RmDArIsAefDwvrdWvRhGs=` /
`h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=`. Both module files
contain only the exact module declaration, with no Go directive or
requirements. Treat every incoming fact only as a survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path candidate. The repository has newer unreleased master
commits; do not silently promote an unreleased commit, redirect, fork,
alternate path, floor-ineligible release, prerelease, or tag that does not
version this module.

# Measurements At Start

The latest dependency implementation is exact Golang Protobuf v1.5.3 commit
`6870e029474b30ecd149d33b6c344c52cced384c`, parent
`d1e076ba35b267331286ca6264354f66c1b6941e`, and tree
`f74214fe65e77543abf43f1ff459636572ce282d`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gogo Protobuf v1.3.2 remains
retained without a dependency commit or metadata edit.

Accepted project measurements are 234 selected modules, 3,589 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,059 `go.sum` lines, and a 404-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 43
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Golang
Protobuf has no finding or trace. Do not attribute inherited findings to
Snappy without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256
`65d833a10bb1cd6b48147d3446f8bc21ef3145cc9b9fc265d9ea2618d02a8807`.
Golang Protobuf decision-summary SHA-256 is
`fec25073702fd5dd3d9ce18fd4aa9bf5b3d9a642bed32bdd690f421778bb324e`;
its 327-entry selected-evidence manifest SHA-256 is
`5d116bb51bfea923778db9b2d6c03bf6c7df2078d966398c5c71a81a68b60dc7`.

Read the answered Golang Protobuf archive and rolling handover for its complete
release, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen Golang Protobuf, Google Martian, Google Protobuf, Go CMP,
or earlier groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
inject no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools
beneath scratch as required. Portable receipts remain golangci-lint 2.12.2
archive `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module
and package/test closure for selected and every serious candidate under exact
Go 1.26.7 and contained Go 1.18.10. Inspect imported source and test
dependencies rather than treating an absent Go directive as floor proof. Keep
isolated source-time resolution separate from the project's selected graph.

Inspect every package and exported API. Characterize raw Snappy block and
framed stream formats, literal and copy encoding/decoding, CRC-32C masking and
validation, maximum decoded length and integer overflow, corrupt/truncated/
unsupported input, reader and writer construction/reset/flush/close behavior,
short reads/writes and error propagation, buffering, allocation and resource
limits, determinism, concurrency and reuse, architecture-specific assembly,
pure-Go fallbacks, build tags, platform behavior, command tools, examples,
benchmarks, testdata, fuzz/property coverage, and upstream CI. Distinguish the
runtime package from commands, generated assembly, and test fixtures.

Add independent fixtures where useful for raw and framed golden vectors,
round trips, incompressible and highly compressible data, checksums, skippable
and unskippable chunks, malformed/truncated/oversized input, short or failing
I/O, flush/close/reset semantics, buffer reuse, concurrency, and compatibility
with the selected project graph. Run source verification, package listing,
native complete tests, two independent repeats, race, vet, and meaningful
cross-builds under both SDKs. Classify every assembly, generator, toolchain,
platform, resource, timing, or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious candidate in disposable trees. Explain why Snappy exists in
MVS while no package is loaded, and preserve every unrelated module selection.
Any change outside the exact Snappy edge and its necessary authorized MVS
projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered Golang Protobuf archive, rolling handover, roadmap,
`go.mod`, `go.sum`, and every referenced quality, compatibility, release,
runner, evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release,
prerelease, or unreleased commit. If no higher release qualifies, retain
selected without hand-editing metadata or manufacturing a dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run applicable dependency, closure, project, compatibility,
vulnerability, empty-HOME, and cleanliness gates without manufacturing
activity. Full changed-selection quality must preserve 27/27 Q0-Q2 PASS at
L2 with zero held, regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never create `.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Golang Snappy decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
