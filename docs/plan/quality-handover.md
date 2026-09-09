# Quality Upgrade Handover

Generated: 2026-09-09T08:22:09+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is exact Go Stack v1.8.1 commit
  `647d4fd71b226fbb1e916b4238b4c7ad87cd7975`, parent
  `171ffd27ac9502c6018c31b2962a6163b2433a03`, tree
  `7c54d59e484d197e7290bc0eec1281c75f614fcd`. It changes only `go.mod` and
  `go.sum`, with three insertions and no deletions.
- Go Logfmt v0.6.0 `3d4cfdbae0a67e757d37022be7eeedaf32c72772`,
  Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0 `e5d6252`, and Speakeasy
  v0.2.0 `41f9561` remain ancestors. All accepted or retained P7 decisions are
  final and documented by their answered archives.
- The answered Go Stack archive and sole NEXT Godbus D-Bus v5 archive must
  link reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish, release,
  stash, revert, launch a successor, bypass cleanup, or remove the worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0, and
Go Stack v1.8.1 moves. Crypt, OpenCensus Proto, Logex, Readline, Fnmatch,
Imaging, ansimage, Fsnotify, Ghodss YAML, and historical root GLFW remain
retained. P8 remains queued. Do not reopen earlier groups or combine another
dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, run `go mod download all` in a
measured worktree, or bypass launcher scratch cleanup.

## Go Stack Decision

Upgrade exact-path `github.com/go-stack/stack` from v1.8.0 to highest stable
v1.8.1. It declares Go 1.17, has no module requirements, preserves the full Go
1.18 floor, and passes every applicable contract. No redirect, fork, alternate
path, floor-ineligible version, or unreleased commit was promoted.

The proxy exposes exactly ten stable versions, v1.4.0 through v1.8.1, and no
prereleases. Go-import maps directly to the public, enabled, unarchived,
non-fork `github.com/go-stack/stack` repository, whose default `master` is
exactly v1.8.1. There are no retractions, deprecation notices, redirects, or
nested exact-path modules. Historical `gopkg.in/stack.v0`,
`gopkg.in/stack.v1`, and the `github.com/zhiyunliu/stack` fork remain distinct.

Selected v1.8.0 sums are
`h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=` /
`h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=`. Qualified v1.8.1 sums are
`h1:ntEHSVwIt7PNXNpgPmVfMrNhLtgjlmnZha2kOpuRiDw=` /
`h1:dcoOX6HbPZSZptuspn9bctJ+N/CnF5gGygcUP3XYfe4=` and its proxy ZIP SHA-256
is `944a204de02272c5718a6819f1f4f6d433a0ef50ca9e737154fceae742694477`.
The incoming ZIP hash was wrong and was not accepted.

V1.8.1's annotated unsigned tag object
`065fe02de4413f66fc553f43129a8e4372e9c54b` peels to commit
`93c7c7e3550c72bc91dead1452a0020142e2a902`, tree
`970d18b6c7c8b790ab52ddd5d82cbaee98261044`, with parents
`2fee6af1a9795aafbe0253a0cfbdf668e1fb8a9a` and
`473edce91b111d1f6c7b946691b24af71f5a5b25`, at
2021-08-18T18:48:21Z. Its release is non-draft and non-prerelease. Proxy ZIPs
for selected and candidate byte-match their peeled Git trees and sumdb confirms
their checksum pairs.

## Closure, API, Behavior, And Gates

Both releases contain one pure-Go package and no external module requirement.
Their complete package/test closures are standard-library-only: 79 entries on
Go 1.18.10 and 125 on Go 1.26.7. Production, tests, examples, and exported API
are byte-identical; v1.8.1 changes only CI declarations and adds `go 1.17`.
The surface is `ErrNoFunc`, `Call`, `CallStack`, `Caller`, `Trace`, and their
frame, PC, formatting, string, marshal, and trim methods.

Independent fixtures cover depth/skip, frame and PC resolution, formatting
verbs and flags, marshal behavior, zero/nil/invalid values, trim semantics,
determinism, concurrent read-only reuse, the 511-frame effective cap, and
no-inline compilation. They pass both SDKs, repeats, and race. Native tests,
two count-10 repeats, vet, source verification, package listing, and 24
platform/architecture cross-test builds also pass both releases and SDKs.

Both releases retain one path-sensitive `go test -trimpath` limitation:
relative runtime paths cannot match the package's empty captured runtime
prefix, so `TrimRuntime` tests fail identically. This is unchanged
compiler/path coupling in byte-identical, unloaded code, not a candidate
regression. Values are immutable for concurrent reads, although callers can
mutate the exported `CallStack` slice.

Project module verification, build, native test, repeats, race, vet, offline
listing, Windows build, pinned lint, API/CLI compatibility, launcher,
empty-HOME, and contained-Go-1.18 gates pass. The Go 1.18 projection removes
only the later toolchain line; both selections retain the same two inherited
`pkg/shell` closed-file wording assertions.

Exact 21-stage `make quality` exits 0, including 80/80 killed mutations, host,
snapshot, Docker, and audit acceptance. All 27 Q0-Q2 rows PASS at L2 with
seven improved and zero held, regressed, not-comparable, or dirty counts.
Scorecard SHA-256 is
`4e1e1b50c2d16666d18278efd0b9ded86df668503e32184c49a1024053dd227c`.
Full audit exits expected 1 only for queued Q3.1/Q3.3/Q3.4/Q3.7; its scorecard
SHA-256 is
`ea43aab34d28e513c92ba9210eb753071e46ecbc983671b7455bfa3bb3a74e0c`.

## Project And Vulnerability Measurements

Prometheus TSDB v0.7.1 and Prometheus Common v0.4.1 each requested v1.8.0.
No loaded package imports Go Stack and `go mod why -m` says the main module
does not need it, but MVS retains those declared edges. The new exact root
v1.8.1 edge changes no unrelated selection and loads no package.

Current measurements are 234 selected modules, 3,585 graph edges, 429
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,055 `go.sum` lines, and a 396-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, sums are +39/-0. Main Go 1.18 and toolchain
Go 1.26.7 remain unchanged.

Fresh primary vulnerability data has 1,392 module records. Selected and
candidate normalize identically to 30 Darwin module findings, 22 Darwin
package findings, and 20 IDs/22 reachable traces for Darwin and Windows
symbol scans. Go Stack has no record, finding, symbol, or reachable trace.

Decision-summary SHA-256 is
`26fca7f6b5b683f4d16fb6b18a9aa88e412948ffde04c7d71764a1a94c077b72`;
the 817-entry selected-evidence manifest SHA-256 is
`e3edb91f49616d36c93ea08e85eeba8ce5730d3f61dea8b668b4da49405a8421`;
manual-evidence SHA-256 is
`98ecb560fb335de466dc47a4f676fcaf918c431caaf2533d2275905b714dffa5`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive, GoReleaser 2.17.1 binary, and apidiff
  receipts remain
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

## Next Objective

Independently evaluate exact-path `github.com/godbus/dbus/v5 v5.0.4` as the
next single P7 group. Do not combine `github.com/coreos/go-systemd/v22` or any
other dependency group.

MVS selects v5.0.4 through go-systemd v22.3.2; `go mod why -m` says the main
module does not need it. The initial proxy survey exposes eleven stable
versions through v5.2.2 and no prereleases. V5.2.0-v5.2.2 declare Go 1.20 and
are floor-ineligible. Highest initial serious candidate v5.1.0 declares Go
1.12, has no requirements, and has source/mod sums
`h1:4KLkAxT3aOY8Li4FRJe/KvhoNFFxo0m6fNuFUO8QJUk=` /
`h1:xhWf0FNVPg57R7Z0UbKHbJfkEywrmjJnf7w5xrFpKfA=`. Treat all of these only as
incoming survey facts to verify.
