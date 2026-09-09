# Agent Session: Evaluate Gogo Protobuf Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T111920+0200-evaluate-gogo-protobuf-dependency`
Created: `2026-09-09T11:19:20+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `03a6e67e92d42850e7f3d6f32eefcd204a2a5a7de4360f27cb290d871b4746f6`
Previous: [2026-09-09T082209+0200-evaluate-godbus-dbus-v5-dependency.md](2026-09-09T082209+0200-evaluate-godbus-dbus-v5-dependency.md)
Next: [2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md](2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md)
Outcome: Retained canonical latest exact-path Gogo Protobuf v1.3.2 without a dependency edit; its complete Go floor, API/behavior, MVS/loading, vulnerability, and applicable project contracts were resolved, with precise upstream generator, randomized-test, vet, and mixed-shared-message race qualifications recorded.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/gogo/protobuf v1.3.2` as one bounded dependency group. Resolve
its complete release and repository identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, and Godbus D-Bus v5.1.0 moves. Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, and
historical root GLFW remain retained. All earlier decisions and lifecycle
ancestry are final. Do not revisit them or combine another dependency group.
P8 remains queued.

Project MVS selects Gogo Protobuf v1.3.2 through declared edges from
`github.com/spf13/viper v1.15.0` and `go.etcd.io/etcd/api/v3 v3.5.1`.
Prometheus TSDB v0.7.1 and Prometheus Common v0.4.1 also declare v1.1.1.
`go mod why -m github.com/gogo/protobuf` says the main module does not need
it. Do not combine, upgrade, remove, or independently audit Viper, etcd,
Prometheus, or any other dependency group.

A minimal post-Godbus survey finds eight proxy versions: v1.0.0, v1.1.0,
v1.1.1, v1.2.0, v1.2.1, v1.3.0, v1.3.1, and selected/latest v1.3.2,
with no prereleases. Selected's source/mod checksum pair is
`h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=` /
`h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=`. It declares Go 1.15
and requires indirect `github.com/kisielk/errcheck v1.5.0`,
`github.com/kisielk/gotool v1.0.0`, and
`golang.org/x/tools v0.0.0-20210106214847-113979e3529a`. Treat every
incoming fact only as a survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, successor/alternate module paths,
and every serious exact-path candidate. Do not silently promote a redirect,
fork, alternate path, floor-ineligible release, or tag that does not version
this module.

# Measurements At Start

The latest dependency implementation is exact Godbus D-Bus v5.1.0 commit
`6472dce617eb80484ed022ae8a53cc350c8be6fe`, parent
`262d97da7a6a50ecc1170bc1a2f33f02732c093e`, and tree
`250d55d4c374958fc9fab703fa7faa931ea4d4f8`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

Accepted project measurements are 234 selected modules, 3,586 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,057 `go.sum` lines, and a 400-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 41
checksum lines and removes zero. The main module retains Go 1.18 and
toolchain Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Godbus has
no record, finding, or trace. Do not attribute inherited findings to Gogo
Protobuf without exact evidence.

The accepted quality result has all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256 `1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.
Godbus decision-summary SHA-256 is
`8973eda3a30a7a1ec4311178f4ee4c701ba114d53ff7454ab8e947f8fcd754d2`;
its 576-entry selected-evidence manifest SHA-256 is
`a665eda7b6763d2e4b0815fc19269ea3d1424eee65e19761ea897a13dd075ac9`.

Read the answered Godbus archive and rolling handover for its complete
release, closure, API, MVS, vulnerability, quality, and evidence record. Do
not reopen Godbus, Viper, etcd, Prometheus, or earlier groups.

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

Inspect every package and exported API. Characterize protobuf wire
encoding/decoding, size and merge behavior, generated-code contracts,
extensions, unknown fields, custom types, maps, oneofs, nullable values,
JSON/text handling, deterministic serialization, registration and global
state, clone/equality/discard behavior, errors and invalid input, recursion
and resource limits, concurrency and reuse, build tags, platform behavior,
generators/plugins, examples, testdata, fuzz/property coverage, and upstream
CI. Distinguish runtime libraries from code-generation commands and generated
fixtures.

Add independent fixtures where useful for golden wire compatibility,
round-trips, deterministic maps, extensions/oneofs/custom types, malformed
and truncated input, unknown fields, nil/zero values, concurrency, and
compatibility with the selected project graph. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify any generator,
toolchain, platform, resource, timing, or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious candidate in disposable trees. Explain why Gogo Protobuf
exists in MVS while no package may be loaded, and preserve every unrelated
module selection. Any change outside the exact Gogo Protobuf edge and its
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Godbus implementation identity, reciprocal archive history,
P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive, the
answered Godbus archive, rolling handover, roadmap, `go.mod`, `go.sum`, and
every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release, or
unreleased commit. If selected v1.3.2 is already the highest qualified
release, retain it without hand-editing metadata or manufacturing a dependency
commit.

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

After the Gogo Protobuf decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit.
Do not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Retain exact-path `github.com/gogo/protobuf v1.3.2` without changing
dependency metadata. It is the canonical latest and highest qualified stable
release, its complete closure preserves Go 1.18, and the project loads no Gogo
package. Exact selected `go get` would add only a redundant indirect root
requirement and source checksum without changing any selected version, so no
dependency commit was manufactured. The latest implementation therefore
remains Godbus D-Bus v5.1.0 commit `6472dce617eb80484ed022ae8a53cc350c8be6fe`.

The proxy contains exactly eight stable versions, v1.0.0 through v1.3.2, and
no prereleases. Sumdb confirms selected source/mod sums
`h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=` /
`h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=`. V1.3.2 is a
non-draft/non-prerelease lightweight unsigned tag at unsigned commit
`b03c65ea87cdc3521ede29f62fe3ce239267c1bc`, tree
`56a9801e7f0e1b7da577e89c8be36d978d243a49`, parent
`550e88954e617545f49920b752c154d72abf1d8d`, dated
2021-01-10T08:01:47Z. Proxy and tagged Git archives contain the same 770
byte-identical paths and share manifest SHA-256
`a259a75ae50454043cfbf66080ceda30be04fe2041bff5de3713eda48574dd5b`.

Go-import, proxy, sumdb, Git, and GitHub agree on the public, enabled,
unarchived, non-fork `github.com/gogo/protobuf` repository, default branch
`master`. There are no retractions, redirects, or exact-path v2 releases.
Release-era documentation sought new ownership and current master explicitly
marks the project deprecated. Three postrelease master commits are unreleased.
Cosmos gogoproto and PlanetScale vtprotobuf are distinct paths with floors
above Go 1.18, not promotable exact-path candidates.

V1.3.2 declares Go 1.15 and resolves 12 modules under both SDKs; every
transitive Go directive is at most 1.15. It exposes 168 packages, including 12
commands and 110 test packages. Complete test listing has 530 entries under
Go 1.26.7 and 485 under Go 1.18.10. All 394 module-backed package entries are
Gogo itself: the historical errcheck, gotool, and x/tools graph loads no
external package. All generator commands and applicable Linux, Windows,
js/wasm, and purego builds pass under both SDKs. Public API across all packages
is unchanged from v1.3.1.

The runtime and fixtures cover generated and reflective wire paths, size,
merge, clone/equality, unknown retention/discard, extensions, custom values,
maps, oneofs, nullable values, JSONPB/text, deterministic ordering,
registration, well-known types, malformed/truncated/overflow input, nil/zero
values, delimited size limits, reuse, and concurrency. Independent fixtures
pass `-race -count=20` under both SDKs. The decoder has integer and size guards
but no general recursion-depth limit, and global registration is intended for
generated initialization rather than concurrent mutation.

Four upstream qualifications are retained explicitly. Generator golden tests
expect three ignored files absent from the release and observe modern gofmt
drift under Go 1.26. Random GoString generation can render nil `*time.Time` as
the invalid expression `<nil>`; seed 11 reproduces it. Vet finds a test-only
`T.Fatal` call from a worker goroutine. Finally, `proto.TestRace` reproducibly
finds an inherited race under both SDKs when binary sizing writes
`XXX_sizecache` while JSON reflection reads the same mutable generated
message. The remainder of the race suite and independent-message concurrency
pass. None is a v1.3.2 regression or loaded Ply behavior.

V1.3.2 is the first fixed release for GO-2021-0053/CVE-2021-3121: its exact
tag commit adds the negative/overflowed generated-unmarshal skip bounds check.
Direct v1.3.1 module/package/symbol scans each find the issue and the symbol
trace reaches `(*unmarshal).Generate`; all three v1.3.2 levels are clean.
Project populations remain 30 Darwin module findings, 22 Darwin package
findings, and 20 IDs/22 reachable traces on Darwin and Windows, with no Gogo
finding or trace.

MVS retains v1.3.2 through Viper v1.15.0 and etcd/api v3.5.1 declarations;
Prometheus TSDB/Common declare v1.1.1. `go mod why -m` remains negative and
zero Gogo packages load. No unrelated selection moved. Current measurements
remain 234 modules, 3,586 graph edges, 429 complete-test entries, 41 loaded
modules, 197 loaded module-backed packages, 1,057 sum lines, +41/-0 sums since
accepted go-cmp, and the accepted 400-line unapplied tidy projection.

Exact Go 1.26.7 project verification, build, count-1, two count-10 repeats,
race, vet, pinned lint, compatibility, offline load, Linux/Windows build,
empty-HOME, and launcher/Make contract components pass. The read-only Go
1.18.10 projection builds and passes applicable repeats/race/vet/cross-builds;
only the two already accepted `pkg/shell` closed-file wording assertions fail.
No changed-selection-only quality activity was manufactured. The accepted
27/27 Q0-Q2 L2 scorecard remains
`1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.

The 2,163-entry selected-evidence manifest SHA-256 is
`b92185ef4d63011df510369fb75a8f35a07d02f11c464efeedc6b0739d38f433`;
decision-summary SHA-256 is
`964fa19e8793454bc1e1d2db71da02cf3a5cc7c4b706e0555ce46dc9b6a12deb`.
