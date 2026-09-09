# Agent Session: Evaluate Golang Snappy Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T162826+0200-evaluate-golang-snappy-dependency`
Created: `2026-09-09T16:28:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `51ecc1f9b8e42db0ceee9ca8fe20251268cc72bb8d291b17dffeba90aac59a21`
Previous: [2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md](2026-09-09T131438+0200-evaluate-golang-protobuf-dependency.md)
Next: [2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md](2026-09-09T183155+0200-evaluate-google-martian-v3-dependency.md)
Outcome: Upgraded exact-path Golang Snappy from v0.0.3 to highest qualified stable v1.0.0 at dependency-only commit `372f8e9`; complete identity, Go-floor closure, API and behavior, MVS/loading, vulnerability, project, and exact quality contracts pass, with the hostile raw-decoder allocation limit documented.

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

## Answer

Upgrade exact-path `github.com/golang/snappy` from v0.0.3 to v1.0.0.
V1.0.0 is the highest stable exact-path release whose complete minimal closure
preserves Go 1.18 and whose applicable API, behavior, project, vulnerability,
and quality contracts pass. Exact Go 1.26.7 `go get
github.com/golang/snappy@v1.0.0` produced dependency-only commit
`372f8e988d9bb86a4f426ae6f953d2e8226c05fe`, parent `838fcd3`, tree
`044149546152b4ce44eb1189721b88f508ea9477`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

The proxy exposes exactly five stable releases, v0.0.1-v0.0.4 and v1.0.0,
with no prerelease. Proxy, sumdb, go-import metadata, tagged Git, and GitHub
agree on the exact public, enabled, unarchived, non-fork repository. There are
no retractions, deprecation, redirects, or qualified alternate paths. Forks,
alternate paths, and the two commits newer than v1.0.0 on master were not
promoted. V1.0.0 is the repository's only GitHub Release; it is neither draft
nor prerelease. Candidate tags are lightweight, so they have no signed tag
objects. GitHub reports valid commit signatures for v0.0.3 and v1.0.0;
v0.0.4's commit is unsigned. The v1.0.0 release was published
2025-03-07T05:42:29Z.

The serious candidates and ancestry are v0.0.3 commit
`674baa8c7fc30da5df3074a459494a7e6b427dff`, tree
`5bbf3ddf1162e96e96a3019f84811da0b2416ee8`, parents
`196ae77b8a26000fa30caa8b2b541e09674dbc43` and
`f81760ec4c9208e6e6866cf18e1888544fd2dbc9`, at 2020-11-03T22:46:00Z;
v0.0.4 commit `544b4180ac705b7605231d4a4550a1acb22a19fe`, tree
`e5a4c4cd6beb5307dd3c32ef56d78a2f6913fbcd`, parent
`0eaccd47634261995ecbb48b545add281213ff88`, at 2021-06-08T04:05:37Z; and
v1.0.0
commit `43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, tree
`c6f19681f2ae79c2e55b6c6113e9b5ea183e3bcb`, parents
`fa5810519dcbfaa60974183cbe11c6c0d050a23d` and
`470b9ed42b44d6e57a1303611bc363992cfed65d`, at 2023-12-25T22:57:46Z.
Master is commit `9ae09f520e93143c782aa6088d18a5e300c69dd0`, tree
`33e272e7ac600aeb47984696311ee0f03902e0f1`, with parents v1.0.0 and
`b6a66aa0e1517e1e9d255d6d6cc750d1a975b543`, two commits after v1.0.0.
Normalized proxy and tagged-Git manifests match byte-for-byte; their SHA-256
identities are `4573785e...`, `d9f61d5d...`, and `a35a0f79...`.

V1.0.0 source/mod sums are
`h1:Oy607GVXHs7RtbggtPBnr2RmDArIsAefDwvrdWvRhGs=` and
`h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=`. All serious candidates'
module files contain only the exact module declaration, with no Go directive
or requirements. Their complete standalone closure is one module using only
the standard library. Complete test loading is 209 entries/four
module-backed packages under Go 1.26.7 and 146/four under Go 1.18.10, proving
the full candidate closure remains within the Go 1.18 floor.

The module exposes runtime package `github.com/golang/snappy` and command
`cmd/snappytool`. Apidiff reports the compatible addition of `Reader.ReadByte`
from v0.0.3 to v0.0.4 and no public change from v0.0.4 to v1.0.0 or from
v1.0.0 to unreleased master. V0.0.4 also fixes ARM64 assembler syntax and adds
golden testdata; v1.0.0 changes README and ARM64 framing details without API
drift.

The runtime implements raw Snappy blocks and checksum-protected framed
streams, including literal and copy tags, masked CRC-32C validation,
skippable-chunk handling, rejection of unskippable/invalid chunks, buffered
read/write, reset, flush, close, and deterministic encoding. Framed chunks
bound encoded and decoded buffers to approximately 76,494 and 65,536 bytes.
Short or failing I/O errors propagate and become sticky where appropriate.
Reader and Writer instances are mutable and not safe for concurrent use;
independent instances are safe. AMD64 and ARM64 use GC assembler unless
`appengine` or `noasm` selects pure Go fallbacks.

All released candidates accept a noncanonical overlong eight-byte raw decoded
length varint; unreleased master tightens this to five bytes. More importantly,
direct `Decode(nil, src)` allocates the declared raw decoded size before full
body validation and can request up to `0xffffffff` bytes on 64-bit systems.
This is a retained, documented resource contract: callers decoding hostile raw
blocks must inspect `DecodedLen` and enforce an application limit first.
Framed decoding is bounded as described above. No candidate otherwise failed
malformed, truncated, checksum, unsupported-input, overflow, I/O, reset,
flush, close, reuse, or concurrency contracts.

The independent 354-line fixture SHA-256 is
`78196ceb64c56ed9c7fb34a0d61de434848ab1b2abe18ed9929a799935358e80`.
It covers raw and framed golden vectors, round trips, incompressible and highly
compressible inputs, checksum failures, skippable and unskippable chunks,
malformed/truncated/oversized inputs, short and failing I/O, flush/close/reset,
buffer reuse, resource behavior, determinism, and independent-instance
concurrency. Race count-20 passes for every candidate under both SDKs.

Module verification, package listing, native count-1, `noasm`, `appengine`,
vet, two independent count-10 repeats, race, and command behavior all pass for
v0.0.3, v0.0.4, and v1.0.0 under exact Go 1.26.7 and contained Go 1.18.10.
All 108 compile-only cross-builds pass across Darwin AMD64/ARM64, Linux
AMD64/ARM64/386, Windows AMD64, FreeBSD AMD64, and js/wasm, including relevant
fallback variants. V1.0.0 contains 31 tests, 38 benchmarks, no examples, and
no native fuzz target. The release tree contains no repository-owned test CI;
GitHub currently exposes only managed dependency-graph and CodeQL workflows.

Google Martian v3.2.1's declared v0.0.3 requirement explained the prior MVS
selection even though `go mod why -m` was negative and no Snappy package was
loaded. Exact v0.0.4 and v1.0.0 projections preserve all 234 module selections
except Snappy and add only the main root edge plus two checksums; Martian's
declaration remains v0.0.3. The accepted project has 234 modules, 3,590 graph
edges, 429 complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, zero Snappy packages, 1,061 sum lines, and a 407-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, sums are +45/-0.
Main Go 1.18 and preferred toolchain Go 1.26.7 remain unchanged.

Fresh vulnerability data contains 1,392 module records and no exact Snappy
record. Direct v0.0.3/v0.0.4/v1.0.0 module, package, and symbol scans are all
zero. Project results are identical across base and candidates: 30 Darwin
module findings, 22 Darwin package findings, 23 Windows package findings, and
20 IDs/22 reachable traces on Darwin and Windows, with no Snappy finding or
trace.

Exact Go 1.26.7 project verification, listing, build, count-1, two independent
count-10 repeats, race, vet, empty-HOME replay, and Linux/Windows builds pass.
The Go 1.18.10 projection removes only the toolchain line; verification,
listing, build, vet, cross-builds, and applicable repeat/race populations
pass. Full count-1 has only the two accepted `pkg/shell` closed-file wording
assertions; the other 26 packages and 31 applicable shell tests pass repeated
and race runs.

Canonical `make quality` passes its ordered 21-stage ledger, kills 80/80
mutants, and records all 27 Q0-Q2 rows PASS at L2 with seven improvements and
zero held, regressed, not-comparable, or dirty counts. Scorecard SHA-256 is
`7320f432090c578973bcfe4abc8736ef68a83f492a5bd057922ec150ce2a03c5`.
Superseded attempts exposed an offline API-base cache miss and macOS
`mktemp` sandbox placement; both were corrected only in session-contained
caches/tool shims. An accidental diagnostic `go mod download all` sum expansion
was reversed byte-exact before the authoritative run; no unrelated repository
change remains.

The 677-entry selected-evidence manifest SHA-256 is
`e4c81eed9a21f737f3a8427360908835b7f195f9981bd1be9a5bdbd953091449`;
decision-summary SHA-256 is
`3d606c0a6ca3f0cd523c4603de9ded354311d1ac6135c9f679ee58371d3ebc79`.
Exact Go 1.26.7 and Go 1.18.10 binary/archive receipts are
`9da68c65...`/`020a1e82...` and `f96ea900...`/`718b32cb...`.
