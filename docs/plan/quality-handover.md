# Quality Upgrade Handover

Generated: 2026-09-09T18:31:55+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Golang Snappy v1.0.0 commit
  `372f8e988d9bb86a4f426ae6f953d2e8226c05fe`, parent
  `838fcd3dff204535d42b3b4078041688ad6de40d`, tree
  `044149546152b4ce44eb1189721b88f508ea9477`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Golang Protobuf v1.5.3 `6870e02`, Godbus D-Bus v5.1.0 `6472dce`, Go
  Stack v1.8.1 `647d4fd`, Go Logfmt v0.6.0 `3d4cfd`, Fatih Color v1.15.0
  `6ca672e`, XXHash v2.3.0 `e5d6252`, and Speakeasy v0.2.0 `41f9561`
  remain ancestors. Gogo Protobuf v1.3.2 remains retained without a
  dependency edit. All earlier P7 decisions are final.
- The answered Golang Snappy archive and sole NEXT Google Martian v3 archive
  link reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, and Golang
Snappy v1.0.0 moves. Gogo Protobuf v1.3.2, Crypt, OpenCensus Proto, Logex,
Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, and historical
root GLFW remain retained. P8 remains queued. Do not reopen earlier groups or
combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Golang Snappy Decision And Identity

Upgrade exact-path `github.com/golang/snappy` from v0.0.3 to v1.0.0.
V1.0.0 is the highest qualified stable exact-path release. The proxy exposes
exactly five stable versions, v0.0.1-v0.0.4 and v1.0.0, with no prerelease.
Go-import, proxy, sumdb, tagged Git, and GitHub agree on the exact public,
enabled, unarchived, non-fork repository. There are no retractions,
deprecation, redirects, or qualified alternate module paths. Forks,
alternates, and the two commits newer than v1.0.0 on master were not promoted.
V1.0.0 is the sole GitHub Release and is neither draft nor prerelease.

Candidate tags are lightweight, so no tag-object signatures exist. GitHub
reports valid commit signatures for v0.0.3 and v1.0.0; v0.0.4 is unsigned.
The serious candidates and ancestry are:

- V0.0.3: commit `674baa8c7fc30da5df3074a459494a7e6b427dff`, tree
  `5bbf3ddf1162e96e96a3019f84811da0b2416ee8`, parents
  `196ae77b8a26000fa30caa8b2b541e09674dbc43` and
  `f81760ec4c9208e6e6866cf18e1888544fd2dbc9`, at 2020-11-03T22:46:00Z.
- V0.0.4: commit `544b4180ac705b7605231d4a4550a1acb22a19fe`, tree
  `e5a4c4cd6beb5307dd3c32ef56d78a2f6913fbcd`, parent
  `0eaccd47634261995ecbb48b545add281213ff88`, at 2021-06-08T04:05:37Z.
- V1.0.0: commit `43d5d4cd4e0e3390b0b645d5c3ef1187642403d8`, tree
  `c6f19681f2ae79c2e55b6c6113e9b5ea183e3bcb`, parents
  `fa5810519dcbfaa60974183cbe11c6c0d050a23d` and
  `470b9ed42b44d6e57a1303611bc363992cfed65d`, at 2023-12-25T22:57:46Z.
- Master: commit `9ae09f520e93143c782aa6088d18a5e300c69dd0`, tree
  `33e272e7ac600aeb47984696311ee0f03902e0f1`, parents v1.0.0 and
  `b6a66aa0e1517e1e9d255d6d6cc750d1a975b543`, at
  2026-07-16T11:44:14Z, two commits after v1.0.0.

V1.0.0 source/mod sums are
`h1:Oy607GVXHs7RtbggtPBnr2RmDArIsAefDwvrdWvRhGs=` /
`h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=`. All candidates'
module files contain only the exact module declaration, with no Go directive
or requirements. Proxy and tagged-Git manifests agree byte-for-byte. Their
SHA-256 identities for v0.0.3, v0.0.4, and v1.0.0 are
`4573785efd0252fca4ec8e6c7fd03693ac733e86252ccd3cfd88884f2f221333`,
`d9f61d5d1bbded9ff0b6a07879223b0b610268e5f58301af486d20d14732b311`,
and `a35a0f797c8ef03c088161b93f2d0f4fb4c289909d30f10becca375ff5c9e5c7`.

## Closure, API, Behavior, And Qualifications

The complete standalone closure for every serious candidate is one module
using only the standard library. It exposes the runtime package and
`cmd/snappytool`. Complete test loading has 209 entries/four module-backed
packages under Go 1.26.7 and 146/four under Go 1.18.10, so the actual closure
preserves the Go 1.18 floor.

Apidiff reports only the compatible `Reader.ReadByte` addition from v0.0.3 to
v0.0.4, no public change from v0.0.4 to v1.0.0, and no change from v1.0.0 to
unreleased master. V0.0.4 fixes ARM64 assembler syntax and adds golden
testdata. V1.0.0 changes README and ARM64 framing details without public API
drift.

The runtime implements raw Snappy blocks and checksum-protected framed
streams: literal/copy tags, masked CRC-32C, skippable chunks, rejection of
unskippable and corrupt chunks, buffered read/write, reset, flush, close, and
deterministic encoding. Framed reader resources are bounded to encoded and
decoded buffers of roughly 76,494 and 65,536 bytes. Short/failing I/O errors
propagate and become sticky where specified. Reader and Writer objects are
mutable and not safe for concurrent reuse; independent instances are safe.
AMD64/ARM64 use GC assembler unless `appengine` or `noasm` selects the pure-Go
fallback.

Every released candidate accepts a noncanonical overlong eight-byte raw length
varint. Unreleased master reduces this to five bytes. Direct `Decode(nil,
src)` also allocates the declared raw decoded size before validating the full
body and can request up to `0xffffffff` bytes on 64-bit systems. Callers must
therefore inspect `DecodedLen` and prebound hostile raw blocks. Framed decoding
is bounded. This is the one retained resource qualification; no candidate
otherwise fails corrupt/truncated/unsupported input, integer overflow,
checksum, I/O, reset/flush/close, reuse, determinism, or concurrency behavior.

The 354-line independent fixture covers raw/framed golden vectors, round
trips, incompressible and highly compressible data, checksum failures,
skippable/unskippable chunks, malformed/truncated/oversized input,
short/failing I/O, flush/close/reset, buffer reuse, resources, determinism,
and independent-instance concurrency. Fixture SHA-256 is
`78196ceb64c56ed9c7fb34a0d61de434848ab1b2abe18ed9929a799935358e80`.
Race count-20 passes every candidate under both SDKs.

Module verification, package listing, native count-1, `noasm`, `appengine`,
vet, two independent count-10 repeats, race, and command behavior pass for all
three candidates under exact Go 1.26.7 and Go 1.18.10. All 108 compile-only
cross-builds pass for Darwin AMD64/ARM64, Linux AMD64/ARM64/386, Windows
AMD64, FreeBSD AMD64, and js/wasm, including fallback variants. V1.0.0 has 31
tests, 38 benchmarks, no examples, and no native fuzz target. The release tree
has no repository-owned test CI; GitHub currently exposes only managed
dependency-graph and CodeQL workflows.

## Project, Vulnerability, And Quality Measurements

Google Martian v3.2.1 declared Snappy v0.0.3. MVS considered that requirement
even though no loaded package imported Snappy, explaining the selected module
alongside a negative `go mod why -m`. Exact v0.0.4 and v1.0.0 projections
change only the Snappy selection and add a root edge plus two checksums;
Martian's declared requirement remains v0.0.3. Every unrelated module
selection is identical.

Current measurements are 234 modules, 3,590 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, zero Snappy
packages, 1,061 sum lines, and a 407-line unapplied tidy projection. Relative
to accepted go-cmp commit `c314bcb`, sums are +45/-0. Main Go 1.18 and
preferred toolchain Go 1.26.7 remain unchanged.

Fresh vulnerability data contains 1,392 module records and no exact Snappy
record. Direct v0.0.3/v0.0.4/v1.0.0 scans are zero at module, package, and
symbol levels. Project results are identical across base and candidates: 30
Darwin module findings, 22 Darwin package findings, 23 Windows package
findings, and 20 IDs/22 reachable traces on Darwin and Windows, with no
Snappy finding or trace.

Exact Go 1.26.7 verification, package load, build, tests, two repeats, race,
vet, empty-HOME, and Linux/Windows builds pass. The Go 1.18.10 projection
removes only the toolchain line; verification, build, vet, cross-builds, and
applicable repeats/race pass. Full count-1 retains only two accepted
`pkg/shell` closed-file wording assertions; the other 26 packages and 31
applicable shell tests pass repeated and race runs.

Exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2, 80/80
mutations killed, seven improvements, and zero held, regressed,
not-comparable, or dirty counts. Scorecard SHA-256 is
`7320f432090c578973bcfe4abc8736ef68a83f492a5bd057922ec150ce2a03c5`.
Superseded quality attempts exposed only an offline API-base cache miss and
macOS `mktemp` sandbox placement; session-contained caches and a scratch-only
shim corrected them. One accidental diagnostic `go mod download all` expanded
`go.sum`; the expansion was reversed byte-exact before authoritative gates.

The 677-entry selected-evidence manifest SHA-256 is
`e4c81eed9a21f737f3a8427360908835b7f195f9981bd1be9a5bdbd953091449`;
decision-summary SHA-256 is
`3d606c0a6ca3f0cd523c4603de9ded354311d1ac6135c9f679ee58371d3ebc79`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Official golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 binary SHA-256 is
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
  Rebuilt apidiff has exact required module/version identity; govulncheck is
  v1.7.0 with database update 2026-09-02T19:12:04Z.

## Next Objective

Independently evaluate exact-path `github.com/google/martian/v3 v3.2.1` as
the next single P7 group. Do not combine Cloud Go, Golang Snappy, Protobuf,
gRPC, X/Net, or another dependency.

The initial proxy survey exposes five stable versions: v3.0.0, v3.1.0,
v3.2.1, v3.3.2, and v3.3.3. Selected v3.2.1 is dated
2021-05-19T22:06:43Z, declares Go 1.11, and has source/mod sums
`h1:d8MncMlErDFTwQGBK1xhv026j9kqhvw1Qv9IbWT1VLQ=` /
`h1:oBOf6HBosgwRXnUGWUB05QECsc6uvmMiJ3+6W4l/CUk=`. Latest v3.3.3 is
dated 2022-08-16T15:12:57Z, declares Go 1.18, and has source/mod sums
`h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=` /
`h1:iEPrYcgCF7jA9OtScMFQyAlZZ4YXTKEtJ1E6RWzmBA0=`. Treat these only as
incoming survey facts to verify.

The module graph contains v3.2.1 declarations from historical Cloud Go
vertices v0.83.0, v0.84.0, v0.87.0, v0.90.0, v0.93.3, v0.94.1, v0.97.0,
and v0.99.0, while the project selects Cloud Go v0.105.0. `go mod why -m` is
negative and no Martian package is loaded. Resolve this MVS ancestry exactly,
but do not independently audit or change Cloud Go.
