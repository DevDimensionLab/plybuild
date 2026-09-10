# Quality Upgrade Handover

Generated: 2026-09-10T01:53:35+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Google Renameio v1.0.1 commit
  `394ec36ac5cb6712640024f0a47af28bf3631906`, parent
  `29b36a6574e2faabedd55784e2576f051d83f655`, tree
  `a55a902a1f491d4972dc2118752b8426d754d1b0`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Google Martian v3.3.2 `4644476`, Golang Snappy v1.0.0 `372f8e9`, Golang
  Protobuf v1.5.3 `6870e02`, Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1
  `647d4fd`, Go Logfmt v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`,
  XXHash v2.3.0 `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain
  ancestors. Google pprof and Gogo Protobuf remain retained without
  dependency edits. All earlier P7 decisions are final.
- The answered Renameio archive and sole NEXT Google UUID archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, and Google Renameio v1.0.1 moves. Google
pprof `v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW remain retained. P8 remains queued. Do
not reopen earlier groups or combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Google Renameio Decision And Identity

Upgrade exact-path `github.com/google/renameio` from v0.1.0 to the highest
qualified stable release v1.0.1. Exact Go 1.26.7
`go get github.com/google/renameio@v1.0.1` created dependency-only commit
`394ec36`. The root proxy exposes only v0.1.0, v1.0.0, and v1.0.1. V2 tags
declare `github.com/google/renameio/v2` and do not version this module.

Proxy, sumdb, `go-import`, and Git agree on the public, active, unarchived,
non-fork Apache-2.0 repository `https://github.com/google/renameio`. No
deprecation, retraction, redirect, qualified fork, alternate root path, or
GitHub Release exists. Stable Git tags qualify the releases; no branch head
or arbitrary commit was promoted.

Selected v0.1.0 is commit `f0e32980c006571efd537032e5f9cd8c1a92819e`,
tree `84b7bc5fe672067d091f8461617877d6286e55c3`, parent
`8bac85ca74197884789b6cd1c9d10170d8b3a194`, at
2019-01-09T16:53:11Z. V1.0.0 is commit
`ad9e5e50f5274088511afd58a920bf584f0a6b3b`. V1.0.1 is commit
`81588dbe0453c6d02ccf518db46a3586e2920084`, tree
`e450d36b9bdc3bbc5e9520da375d61f88e35ff49`, parent
`650fcb4eeb71a5dbcaa1ec3592f48d369f55489f`, at
2021-04-06T14:11:08Z. Selected precedes v1.0.1 by nine commits. The v1
commits have valid GitHub signatures; annotated tag objects are unsigned.

Selected source/mod sums are
`h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=` /
`h1:KWCgfxg9yswjAJkECMjeO8J8rahYeXnNhOm40UhjYkI=`. V1.0.1 sums are
`h1:Lh/jXZmvZxb0BBeSY5VKEfidcbcbenKjZFzM/q0fSeU=` /
`h1:t/HQoYBZSsWSNK35C6CO/TpPLDVWvxOHboWUAweKUpk=`. Normalized proxy/Git
manifests match; v0.1.0, v1.0.0, and v1.0.1 manifest hashes are
`dc9b1550d412d9a381527ba737e6a85ec5f3f9a98018dd65d9fc5161b71b8fff`,
`f9f7c804d3162a5bab090d8738d497a5908305fe888258e47b3cb3ab3c8de207`,
and `130e3c033d66f007ae309e70233209af9359ce47486585e33415df81e200f956`.

## Closure, API, Behavior, And Qualifications

Every serious root release has a one-module, standard-library-only closure.
V0.1.0 has no Go directive; both v1 releases declare Go 1.13. Exact Go
1.26.7 and contained Go 1.18.10 pass verification, complete-test loading,
native tests, two repeats, race, and vet. The root exported API remains
`TempDir`, `TempFile`, embedded-file `PendingFile`, `Cleanup`,
`CloseAtomicallyReplace`, `Symlink`, and `WriteFile`. V1 adds only compatible
package `renameio/maybe`; apidiff otherwise is empty.

Forty-eight production cross-builds cover Linux, Windows, FreeBSD, OpenBSD,
NetBSD, DragonFly, Solaris, and js/wasm under both SDKs. Selected and v1.0.1
pass all. V1.0.0 fails Windows because `maybe/maybe_unix.go` references the
Windows-excluded root `WriteFile`; v1.0.1 repairs this with OS-specific files
and is the qualified candidate.

A 309-line independent fixture covers atomic replacement and old-or-whole-new
visibility, exact modes, 0600 temporary files, hard links, symlinks, temp
placement, missing/read-only parents, write/sync/rename failures, cleanup,
reuse, and 96 concurrent writers/readers. Race count-10 passes under both
SDKs; fixture SHA-256 is
`44f371058ec5076187b73f313e6855dacef26c9bfe19b6ecffc4a2cb8c5a1a53`.

`WriteFile` creates a random dot-prefixed temp, chmods before one write, then
orders file Sync, Close, and Rename. Cleanup removes failed temps and becomes
a no-op after success. Automatic TempDir probes same-filesystem rename and
falls back to the destination directory; explicit placement remains
caller-controlled. Local rename is atomic, independent writers are
last-writer-wins, and readers see old or whole new content. Replacement
breaks hard links, replaces a symlink rather than its referent, creates a
caller-owned inode, and does not preserve prior ownership.

The parent directory is not fsynced, so the directory entry lacks a complete
crash-durability guarantee. NFS multi-client atomicity is filesystem-specific.
Sharing a `PendingFile` concurrently is unsupported. Missing/read-only
parents, directory targets, cross-device explicit temps, and
write/sync/close/rename failures propagate and clean up. On Windows,
v1.0.1 `maybe.WriteFile` deliberately falls back to non-atomic
`ioutil.WriteFile`. No commands, generators, benchmarks, fuzz targets,
property tests, or testdata exist; examples compile without Output checks.

## Project, Vulnerability, And Quality Measurements

Three historical HTools vertices declare v0.1.0. One exact path is main ->
mvn-pom-mutator v0.2.3 -> Crypt -> Firestore v1.1.0 -> Cloud Go v0.46.3 ->
HTools v0.0.1-2019.2.3 -> Renameio. MVS retains requirements from encountered
historical vertices, but no loaded package import reaches Renameio;
`go mod why -m` is negative and zero Renameio packages load.

The accepted graph exactly matches the v1.0.1 projection: 234 modules, 3,598
edges, 429 complete-test entries, 41 loaded modules, 197 loaded module-backed
entries, 1,065 sum lines, and a 428-line unapplied tidy projection. Only the
exact root edge and v1.0.1 source/mod sums are new. Every unrelated selection
and package is identical. Relative to go-cmp commit `c314bcb`, sums are +49/-0.
The main module remains Go 1.18 with toolchain Go 1.26.7.

Fresh vulnerability data has 1,393 records, Last-Modified
2026-09-09T18:12:50Z, and no Renameio record. Direct selected and v1.0.1
scans are zero. V1.0.0 module/package and Darwin symbol scans are zero; its
Windows symbol load is blocked by the confirmed compile defect. Base/v1.0.1
project populations are identical at 30 module findings, 22 Darwin and 23
Windows package findings, and 20 IDs/22 reachable traces per symbol platform.
None involves Renameio.

Exact Go 1.26.7 project verification/load/build, two independent count-10
repeats, race, vet, pinned lint, empty-HOME count-2, Linux/Windows builds,
API/CLI and CLI-surface compatibility, launcher checks, host acceptance,
snapshot acceptance, and Docker acceptance pass. The Go 1.18.10 projection
removes only the toolchain line and passes verification/load/build/vet,
cross-builds, 26-package repeat/race, and all 31 compatible `pkg/shell`
tests. Full count-1 retains only the two accepted closed-file wording
differences.

The authoritative quality gate has an ordered 21-stage ledger, kills 80/80
mutants, and records all 27 Q0-Q2 rows PASS at L2 with zero held, regressed,
non-comparable, or dirty counts. Scorecard SHA-256 is
`6be51aec069011e99a80c2b5b5097bb79a67abe1c0078589ccd60198ab75c2ce`.
The 7,798-entry evidence manifest and decision summary hash to
`4d4fd72b84786270f069b18934eb831d31ecc801e1d348e6e912f756fc114f84`
and `07fcaace013997f5d4513de49c5a59b3359b7b270c4e6634b39729c8a45a6096`.

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
  v1.7.0.
- The final quality run needed a scratch-only bare-`mktemp` shim for the
  macOS sandbox alias and Docker Desktop Buildx v0.33.0-desktop.1 copied into
  the empty scratch HOME. These are external execution controls, not project
  changes. Earlier attempts are preserved and classified in the evidence.

## Next Objective

Independently evaluate exact-path `github.com/google/uuid v1.1.2` as the next
single P7 dependency group. Do not combine gRPC, Martian, or another
dependency.

Fourteen historical `google.golang.org/grpc` vertices spanning v1.33.1 through
v1.43.0 declare UUID v1.1.2. One shortest current graph path is main ->
Google Martian/v3 v3.3.2 -> gRPC v1.37.0 -> UUID v1.1.2. `go mod why -m` is
negative and no UUID package loads. This is graph ancestry only; do not
reopen Martian or independently audit gRPC.

The initial exact proxy list has 13 stable versions: v1.0.0, v1.1.0 through
v1.1.5, v1.2.0, v1.3.0, v1.3.1, v1.4.0, v1.5.0, and v1.6.0. Selected is
dated 2020-07-02T18:56:42Z with source/mod sums
`h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=` /
`h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=`. Proxy latest v1.6.0 is
dated 2024-01-23T18:54:04Z at origin commit
`0f11ee6918f41a04c201eceeadf612a377bc7fbc`, with source/mod sums
`h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=` /
`h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=`. Both module files have
no Go directive or requirements. Treat these only as incoming survey facts
to verify.
