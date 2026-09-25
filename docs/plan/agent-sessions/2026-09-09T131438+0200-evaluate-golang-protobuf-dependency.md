# Agent Session: Evaluate Golang Protobuf Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T131438+0200-evaluate-golang-protobuf-dependency`
Created: `2026-09-09T13:14:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fca0510206da51b0f3e66d08ba3828b42f08cd73587c14510eb8ec113c9170e3`
Previous: [2026-09-09T111920+0200-evaluate-gogo-protobuf-dependency.md](2026-09-09T111920+0200-evaluate-gogo-protobuf-dependency.md)
Next: [2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md](2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md)
Outcome: Upgraded exact-path Golang Protobuf from v1.5.2 to highest qualified stable v1.5.3 at dependency-only commit `6870e02`; v1.5.4 was rejected for a public descriptor API removal, while complete identity, Go-floor, behavior, MVS/loading, vulnerability, and exact quality contracts pass for v1.5.3.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/golang/protobuf v1.5.2` as one bounded dependency group. Resolve
its complete release and repository identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, and Godbus D-Bus v5.1.0 moves. Gogo Protobuf v1.3.2,
Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage,
Fsnotify, Ghodss YAML, and historical root GLFW remain retained. All earlier
decisions and lifecycle ancestry are final. Do not revisit them or combine
another dependency group. P8 remains queued.

Project MVS selects Golang Protobuf v1.5.2. Selected Viper v1.15.0, Google
Martian v3.2.1, and protoc-gen-star v0.5.3 declare v1.5.2; thirteen other
selected modules declare versions from v1.2.0 through v1.5.0. `go mod why -m
github.com/golang/protobuf` says the main module does not need it, and the
complete project package load contains no Golang Protobuf package. Do not
combine, upgrade, remove, or independently audit any declaring module,
`google.golang.org/protobuf`, Go CMP, or another dependency group. Candidate
MVS movements caused by the exact Golang Protobuf edge remain in scope to
measure, but not to broaden into independent audits.

A minimal post-Gogo survey finds 22 proxy versions: eighteen stable releases
from v1.0.0 through latest v1.5.4 and four v1.4.0 release candidates.
Selected v1.5.2 has source/mod sums
`h1:ROPKBNFfQgOUMifHyP+KYbvpjbdoFNs+aK7DXlji0Tw=` /
`h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=`, declares Go 1.9,
and requires `github.com/google/go-cmp v0.5.5` plus
`google.golang.org/protobuf v1.26.0`. Latest v1.5.4 is dated
2024-03-06T06:45:40Z at commit
`75de7c059e36b64f01d0dd234ff2fff404ec3374`, has source/mod sums
`h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=` /
`h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=`, declares Go 1.17,
requires Go CMP v0.5.5 and Google Protobuf v1.33.0, and marks the module
deprecated in favor of `google.golang.org/protobuf`. Treat every incoming fact
only as a survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, successor/alternate module paths,
and every serious exact-path candidate. Do not silently promote a redirect,
fork, alternate path, floor-ineligible release, prerelease, or tag that does
not version this module.

# Measurements At Start

The latest dependency implementation remains exact Godbus D-Bus v5.1.0
commit `6472dce617eb80484ed022ae8a53cc350c8be6fe`, parent
`262d97da7a6a50ecc1170bc1a2f33f02732c093e`, and tree
`250d55d4c374958fc9fab703fa7faa931ea4d4f8`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gogo Protobuf v1.3.2 was
retained without a dependency commit or metadata edit.

Accepted project measurements remain 234 selected modules, 3,586 graph
edges, 429 native complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,057 `go.sum` lines, and a 400-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, metadata adds
exactly 41 checksum lines and removes zero. The main module retains Go 1.18
and toolchain Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Gogo
Protobuf v1.3.2 has no finding or trace; v1.3.1 is affected by GO-2021-0053.
Do not attribute inherited findings to Golang Protobuf without exact evidence.

The accepted quality result remains all 27 Q0-Q2 rows PASS at L2, with
scorecard SHA-256
`1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.
Gogo Protobuf decision-summary SHA-256 is
`964fa19e8793454bc1e1d2db71da02cf3a5cc7c4b706e0555ce46dc9b6a12deb`;
its selected-evidence manifest identity is recorded in the answered Gogo
archive and rolling handover.

Read the answered Gogo Protobuf archive and rolling handover for its complete
release, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen Gogo Protobuf, Viper, other declaring modules, Google
Protobuf, Go CMP, or earlier groups.

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
dependencies rather than treating a Go directive alone as floor proof. Keep
isolated source-time resolution separate from the project's selected graph.

Inspect every package and exported API. Characterize the legacy protobuf wire
and generated-code contracts, bridging to the newer Google Protobuf runtime,
descriptors and reflection, registration/global state, extensions, unknown
fields, maps, oneofs, nullable values, well-known types, JSON/text behavior,
deterministic serialization, size/merge/clone/equality/discard behavior,
errors and invalid input, recursion/resource limits, concurrency and reuse,
build tags, platform behavior, generators/plugins, examples, testdata,
fuzz/property coverage, and upstream CI. Distinguish runtime libraries,
compatibility wrappers, code-generation commands, and generated fixtures.

Add independent fixtures where useful for old/new runtime interoperability,
golden wire compatibility, round-trips, deterministic maps, extensions and
oneofs, descriptors and well-known types, malformed/truncated input, unknown
fields, nil/zero values, concurrency, and compatibility with the selected
project graph. Run source verification, package listing, native complete
tests, two independent repeats, race, vet, and meaningful cross-builds under
both SDKs. Classify every generator, toolchain, platform, resource, timing,
or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious candidate in disposable trees. Explain why Golang Protobuf
exists in MVS while no package is loaded, and preserve every unrelated module
selection. Any change outside the exact Golang Protobuf edge and its necessary
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered Gogo archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

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

After the Golang Protobuf decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Upgrade exact-path `github.com/golang/protobuf` from v1.5.2 to v1.5.3.
V1.5.3 is the highest stable exact-path release whose complete minimal closure
preserves Go 1.18 and whose applicable behavior and public API pass. Exact Go
1.26.7 `go get github.com/golang/protobuf@v1.5.3` produced dependency-only
commit `6870e029474b30ecd149d33b6c344c52cced384c`, parent `d1e076ba`, tree
`f74214fe65e77543abf43f1ff459636572ce282d`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

The proxy lists exactly 22 versions: 18 stable releases from v1.0.0 through
v1.5.4 and four v1.4.0 release candidates. Go-import metadata, proxy, sumdb,
Git, and GitHub agree on the exact public, enabled, unarchived, non-fork
repository. There are no retractions or redirects. The exact module and README
are deprecated in favor of distinct successor path
`google.golang.org/protobuf`; that successor, forks, alternate paths,
prereleases, and unreleased master commits were not promoted.

Candidate tags are lightweight and therefore have no signed tag objects;
GitHub verifies the signed commits. The release ancestry is v1.5.2 commit
`ae97035608a719c7a1c1c41bed0ae0744bdb0c6f`, tree
`d5f27d9813f7fd9318cd2c502ffe08bed57c451b`, parent `af940030`, at
2021-03-29T18:20:59Z; v1.5.3 commit
`5d5e8c018a13017f9d5b8bf4fad64aaa42a87308`, tree
`65ba39c63b7751b84b386c74f84d2f71c1b78a41`, parents `ae970356` and
`37828f7c`, at 2021-09-16T00:37:10Z; and v1.5.4 commit
`75de7c059e36b64f01d0dd234ff2fff404ec3374`, tree
`49890db482e1da888cf22e4502cd742c36c1e92a`, parents `5d5e8c01` and
`b7696c3f`, at 2024-03-06T06:45:40Z. Proxy and tagged-Git manifests agree,
with SHA-256 identities `af276680...`, `c6c6cbfe...`, and `d7fce623...`.

V1.5.3 retains Go 1.9, Go CMP v0.5.5, and Google Protobuf v1.26.0. Its
source/mod sums are
`h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=` and
`h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=`. V1.5.4 declares Go
1.17 and moves Google Protobuf to v1.33.0, but apidiff finds the incompatible
public removal
`protoc-gen-go/descriptor.Default_FileOptions_PhpGenericServices`.
V1.5.4 is therefore rejected. V1.5.3 has no public API change from v1.5.2;
its focused behavior change accepts JSON `null` for `NullValue`.

Each candidate exposes 21 packages, including two commands and five packages
with tests. V1.5.2/v1.5.3 resolve four modules; complete test loading has
219 entries/75 module-backed under Go 1.26.7 and 171/75 under Go 1.18.10.
V1.5.4 has 226/81 and 178/81 respectively. All complete module floors are at
most Go 1.11 for v1.5.2/v1.5.3 and Go 1.17 for v1.5.4. Module verification,
package listing, two native count-10 repeats, race, production vet, command
builds, and Linux, Windows, and js/wasm compilation pass for every candidate
under both SDKs. Full upstream vet reports only historical unkeyed composite
literals in test fixtures: 447 findings for v1.5.2 and 449 for v1.5.3/v1.5.4.
Production-only vet passes. Initial cross-build execution produced the
expected host exec-format error; corrected compile-only cross-builds pass.

The runtime packages are APIv1 compatibility wrappers over APIv2 behavior,
covering legacy generated-code wire contracts, descriptors and reflection,
registrations, extensions, oneofs, maps, nullable and unknown fields,
JSON/text and well-known types, deterministic marshaling, size, merge, clone,
equality, and unknown discard. Global caches are synchronized, while global
registration mutation is intended for generated initialization. Buffers and
mutated messages are not safe for concurrent reuse. V1.5.2/v1.5.3 with Google
Protobuf v1.26.0 have no general binary-decode recursion cap, so callers must
bound hostile deeply nested messages; v1.5.4's v1.33.0 runtime adds a
10,000-depth limit. Upstream CI covers Go 1.11-1.16 on Linux and macOS; there
are no build tags or native fuzz/property corpus.

The independent fixture's SHA-256 is
`51f4b0343a10c5f65b7fff613b5df896f98a7698d06acdcdaf5576f930429718`.
It covers golden wire bytes, APIv1/APIv2 bridging, round trips, size/merge/
clone/equality, deterministic maps, extensions, oneofs, descriptors,
reflection and global registration, timestamps and Any, JSON/text, unknown
retention/discard, nil/zero values, malformed/truncated/overflow input,
buffer reuse, and independent-message concurrency. Race count-20 passes under
both SDKs for all candidates. V1.5.3 also passes under both SDKs against the
project-selected Google Protobuf v1.28.1.

Sixteen selected modules declare Golang Protobuf versions from v1.2.0 through
v1.5.2. MVS considers those requirements even though no loaded package imports
the module, which is why `go mod why -m` remains negative while the module is
selected. The v1.5.3 projection changes only this selection; project-selected
Google Protobuf remains v1.28.1 and every unrelated module is identical.
V1.5.4 would additionally select Google Protobuf v1.33.0. The accepted project
has 234 modules, 3,589 graph edges, 429 complete-test entries, 41 loaded
modules, 197 loaded module-backed packages, zero Golang Protobuf packages,
1,059 sum lines, and a 404-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, sums are +43/-0. The main module remains Go
1.18 with preferred toolchain Go 1.26.7.

Fresh vulnerability data contains 1,392 module records and no record for this
exact module. Direct selected/candidate module, package, and symbol scans all
produce zero findings. Base, v1.5.3, and v1.5.4 project results are identical:
30 Darwin module findings, 22 Darwin package findings, and 20 IDs/22 reachable
traces for Darwin and Windows, with no Golang Protobuf finding or trace.

Exact Go 1.26.7 project verification, build, count-1, two count-10 repeats,
race, vet, pinned lint, offline loading, Linux/Windows builds, API/CLI
compatibility, CLI surface, empty-HOME replay, host acceptance, and fresh
snapshot/Docker acceptance pass. The Go 1.18.10 projection removes only the
toolchain line; verification, build, vet, cross-builds, and applicable repeat/
race populations pass. Full count-1 has only the two accepted `pkg/shell`
closed-file wording assertions; the other 26 packages and 31 applicable shell
tests pass repeated and race runs.

Exact `make quality` passes its ordered 21-stage ledger, kills 80/80 mutants,
and records all 27 Q0-Q2 rows PASS at L2 with zero held, regressed,
not-comparable, or dirty counts. Scorecard SHA-256 is
`65d833a10bb1cd6b48147d3446f8bc21ef3145cc9b9fc265d9ea2618d02a8807`.
A Docker 29 nanosecond timestamp exposed the macOS Python 3.9
`fromisoformat` limitation in one superseded run; the same metadata validates
under physical Python 3.14, and isolated plus authoritative Docker contracts
then pass completely. No repository workaround was introduced.

The 327-entry selected-evidence manifest SHA-256 is
`5d116bb51bfea923778db9b2d6c03bf6c7df2078d966398c5c71a81a68b60dc7`;
decision-summary SHA-256 is
`fec25073702fd5dd3d9ce18fd4aa9bf5b3d9a642bed32bdd690f421778bb324e`.
Exact Go 1.26.7 and Go 1.18.10 binary/archive receipts are
`9da68c65...`/`020a1e82...` and `f96ea900...`/`718b32cb...`.
