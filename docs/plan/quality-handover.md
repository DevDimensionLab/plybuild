# Quality Upgrade Handover

Generated: 2026-09-09T16:28:26+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Golang Protobuf v1.5.3 commit
  `6870e029474b30ecd149d33b6c344c52cced384c`, parent
  `d1e076ba35b267331286ca6264354f66c1b6941e`, tree
  `f74214fe65e77543abf43f1ff459636572ce282d`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Gogo
  Protobuf v1.3.2 remains retained without a dependency edit. All earlier
  P7 decisions are final.
- The answered Golang Protobuf archive and sole NEXT Golang Snappy archive
  link reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, and Golang Protobuf v1.5.3 moves.
Gogo Protobuf v1.3.2, Crypt, OpenCensus Proto, Logex, Readline, Fnmatch,
Imaging, ansimage, Fsnotify, Ghodss YAML, and historical root GLFW remain
retained. P8 remains queued. Do not reopen earlier groups or combine another
dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Golang Protobuf Decision And Identity

Upgrade exact-path `github.com/golang/protobuf` from v1.5.2 to v1.5.3.
V1.5.3 is the highest qualified stable exact-path release. V1.5.4 was
rejected because it removes public symbol
`protoc-gen-go/descriptor.Default_FileOptions_PhpGenericServices`.

The proxy exposes 22 versions: 18 stable v1.0.0-v1.5.4 and four v1.4.0
release candidates. Go-import, proxy, sumdb, Git, and GitHub agree on the
public, enabled, unarchived, non-fork repository and exact module path. There
are no retractions, redirects, or mirrors. The module is deprecated in favor
of distinct successor `google.golang.org/protobuf`; that path, forks,
alternates, prereleases, and unreleased commits were not promoted.

Candidate tags are lightweight, so no tag-object signatures exist. GitHub
reports the signed commits verified. The candidates and ancestry are:

- V1.5.2: commit `ae97035608a719c7a1c1c41bed0ae0744bdb0c6f`, tree
  `d5f27d9813f7fd9318cd2c502ffe08bed57c451b`, parent
  `af940030a2b77f37337632168f18403433e21e61`, at
  2021-03-29T18:20:59Z.
- V1.5.3: commit `5d5e8c018a13017f9d5b8bf4fad64aaa42a87308`, tree
  `65ba39c63b7751b84b386c74f84d2f71c1b78a41`, parents `ae970356` and
  `37828f96`, at 2021-09-16T00:37:10Z.
- V1.5.4: commit `75de7c059e36b64f01d0dd234ff2fff404ec3374`, tree
  `49890db482e1da888cf22e4502cd742c36c1e92a`, parents `5d5e8c01` and
  `b7697bb6`, at 2024-03-06T06:45:40Z. It is the reviewed master head.

V1.5.3 source/mod sums are
`h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=` /
`h1:XVQd3VNwM+JqD3oG2Ue2ip4fOMUkwXdXDdiuN0vRsmY=`. It retains Go 1.9,
Go CMP v0.5.5, and Google Protobuf v1.26.0. V1.5.4 declares Go 1.17 and
Google Protobuf v1.33.0. Proxy/tagged-Git manifest identities for v1.5.2,
v1.5.3, and v1.5.4 are `af276680...`, `c6c6cbfe...`, and `d7fce623...`.

## Closure, API, Behavior, And Qualifications

All candidates expose 21 packages, including two commands and five packages
with tests. V1.5.2/v1.5.3 resolve four modules. Complete test loading has
219 entries/75 module-backed under Go 1.26.7 and 171/75 under Go 1.18.10.
V1.5.4 has 226/81 and 178/81. Complete closure floors are at most Go 1.11
for v1.5.2/v1.5.3 and Go 1.17 for v1.5.4.

Module verification, package listing, two native count-10 repeats, race,
production vet, both commands, and Linux, Windows, and js/wasm compile-only
cross-builds pass for every candidate under both SDKs. Full upstream vet
reports only historical unkeyed composite literals in test fixtures: 447 for
v1.5.2 and 449 for v1.5.3/v1.5.4. Production vet passes. Initial attempts to
execute cross-built binaries produced expected host exec-format errors; the
corrected compile-only population passes.

Apidiff finds no public change from v1.5.2 to v1.5.3. V1.5.3's focused
behavior change accepts JSON `null` for `NullValue`. V1.5.4 adds compatible
descriptor/generator APIs but has the disqualifying public removal above.

The APIv1 runtime wrappers bridge to APIv2 for legacy generated wire code,
descriptors/reflection, registrations, extensions, oneofs, maps, nullable and
unknown fields, JSON/text and well-known types, deterministic marshal, size,
merge, clone, equality, and discard. Global caches are synchronized; global
registration mutation is generated-init behavior. Buffers and mutated
messages are not safe for concurrent reuse. V1.5.2/v1.5.3 with Google
Protobuf v1.26.0 have no general binary recursion cap, so callers must bound
hostile deeply nested messages. V1.5.4's runtime adds a 10,000-depth cap.
Upstream CI covers Go 1.11-1.16 on Linux/macOS; no build tags or native
fuzz/property corpus exist.

Independent fixtures cover golden wire, APIv1/APIv2 interoperability,
round trips, deterministic maps, extensions, oneofs, descriptors/reflection/
registration, well-known types, JSON/text, unknown retention/discard,
nil/zero and malformed/truncated/overflow input, buffer reuse, and
independent-message concurrency. Fixture SHA-256 is
`51f4b0343a10c5f65b7fff613b5df896f98a7698d06acdcdaf5576f930429718`.
Race count-20 passes for all candidates under both SDKs. V1.5.3 also passes
against project-selected Google Protobuf v1.28.1.

## Project, Vulnerability, And Quality Measurements

Sixteen selected modules declare Golang Protobuf v1.2.0-v1.5.2. MVS reads
requirements of selected modules even when no loaded package imports the
dependency; therefore v1.5.2 was selected while `go mod why -m` was negative
and the package population was zero. The accepted exact root edge selects
v1.5.3 only. Google Protobuf remains v1.28.1 and every unrelated module
selection is unchanged. The rejected v1.5.4 projection additionally selected
Google Protobuf v1.33.0.

Current measurements are 234 modules, 3,589 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, zero Golang
Protobuf packages, 1,059 sum lines, and a 404-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, sums are +43/-0. Main Go 1.18
and preferred toolchain Go 1.26.7 remain unchanged.

Fresh vulnerability data contains 1,392 module records and no exact-module
record. Direct v1.5.2/v1.5.3/v1.5.4 scans are zero at module, package, and
symbol levels. Project results are identical across base and candidates: 30
Darwin module findings, 22 Darwin package findings, and 20 IDs/22 reachable
traces on Darwin and Windows, with no Golang Protobuf finding or trace.

Exact Go 1.26.7 verification, build, test, two repeats, race, vet, pinned
lint, offline package load, Linux/Windows builds, API/CLI compatibility, CLI
surface, empty-HOME, and host/snapshot/Docker acceptance pass. The Go 1.18.10
projection removes only the toolchain line; build/vet/cross-build and
applicable repeat/race gates pass. Full count-1 has only the two accepted
`pkg/shell` closed-file wording assertions; the other 26 packages and 31
applicable shell tests pass repeats and race.

Exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2,
80/80 mutations killed, seven improvements, and zero held, regressed,
not-comparable, or dirty counts. Scorecard SHA-256 is
`65d833a10bb1cd6b48147d3446f8bc21ef3145cc9b9fc265d9ea2618d02a8807`.
A superseded Docker run exposed only the macOS Python 3.9 inability to parse
Docker 29's valid nanosecond timestamp. Physical Python 3.14 accepts it;
isolated and authoritative Docker gates then pass without a repository edit.

The 327-entry selected-evidence manifest SHA-256 is
`5d116bb51bfea923778db9b2d6c03bf6c7df2078d966398c5c71a81a68b60dc7`;
decision-summary SHA-256 is
`fec25073702fd5dd3d9ce18fd4aa9bf5b3d9a642bed32bdd690f421778bb324e`.

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

Independently evaluate exact-path `github.com/golang/snappy v0.0.3` as the
next single P7 group. Do not combine Google Martian or another dependency.

The initial proxy survey exposes five stable versions v0.0.1-v0.0.4 and
v1.0.0, with no prerelease. Selected v0.0.3 has source/mod sums
`h1:fHPg5GQYlCeLIPB9BZqMVR5nR9A+IM5zcgeTdjMYmLA=` /
`h1:/XxbfmMg8lxefKM7IXC3fBNl/7bRcc72aCRzEWrmP2Q=`. Latest v1.0.0 has
source/mod sums `h1:Oy607GVXHs7RtbggtPBnr2RmDArIsAefDwvrdWvRhGs=` /
the same mod sum. Both module files contain only the exact module declaration,
with no Go directive or requirements. Selected is declared only by Google
Martian v3.2.1; `go mod why -m` is negative and no Snappy package is loaded.
Treat all of these only as incoming survey facts to verify.
