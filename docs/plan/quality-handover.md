# Quality Upgrade Handover

Generated: 2026-09-10T05:48:41+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Google UUID v1.4.0 commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Google
  pprof and Gogo Protobuf remain retained without dependency edits. All
  earlier P7 decisions are final.
- The answered UUID archive and sole NEXT Googleapis GAX Go v2 archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof
`v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW remain retained. P8 remains queued. Do
not reopen earlier groups or combine another dependency group.

On 2026-09-12 the user explicitly authorized the bounded inherited Go-floor
exception for unloaded GAX v2.7.0. Keep the main-module Go 1.18 floor, retain
GAX v2.7.0 without a dependency edit if its remaining qualification finds no
separate disqualifier, and do not broaden into a GAX downgrade, Viper or
parent-graph work, or a Go-floor increase. Do not ask for this same product
decision again.

The user also authorized the recommended bounded behavior exception for the
known v2.7.0 `apierror.ParseError(err, false)` panic. It applies only while the
complete project load contains zero GAX packages. Revalidate and record that
invariant in the final decision and handover. If any GAX package becomes loaded
or directly imported later, the exception expires and the owning checkpoint
must stop for a fresh dependency and product decision before merge. Upstream
first fixed the defect in v2.12.5, which declares Go 1.20. Do not patch or
replace GAX, add a root edge, change Viper or another parent, raise the Go
floor, or ask for this same behavior decision again.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Google UUID Decision And Identity

Upgrade exact-path `github.com/google/uuid` from v1.1.2 to the highest
qualified stable release v1.4.0. Exact Go 1.26.7
`go get github.com/google/uuid@v1.4.0` created dependency-only commit
`cf53bc64`.

The exact proxy exposes thirteen stable releases from v1.0.0 through v1.6.0.
Proxy, sumdb, `go-import`, Git tags, and GitHub Releases agree on the public,
active, unarchived, non-fork BSD-3-Clause repository
`https://github.com/google/uuid`. There is no deprecation, retraction,
redirect, qualified fork, or alternate path. Malformed/non-version tags,
master, and the unmerged release-please 1.7.0 branch do not qualify.

Selected v1.1.2 is commit
`0e4e31197428a347842d152773b4cace4645ca25`, tree
`6dfe4fade2a1ba2cfd25aef72398b1581fa4b8ba`, parent
`cb32006e483f2a23230e24209cf185c65b477dbf`, at
2020-07-02T18:56:42Z. V1.4.0 is commit
`8de8764e294f072b7a2f1a209e88fdcdb1ebc875`, tree
`5c05b7240eb99c0f069b47ddbc7d0eefbf43cf80`, parent
`7c22e97ff7647f3b21c3e0870ab335c3889de467`, at
2023-10-26T15:24:04Z. Latest v1.6.0 is commit
`0f11ee6918f41a04c201eceeadf612a377bc7fbc`, tree
`42ba8f689f0586db861c6fdef4f0042efc62c958`, parent
`16939dafc37a38d2743810a8bdf60fdad6a0f3a3`, at
2024-01-23T18:54:04Z. Selected is an ancestor of every higher release.

V1.4.0 source/mod sums are
`h1:MtMxsa51/r9yyhkyLsVeVt0B+BGQZzpQiTQ4eHZ8bc4=` /
`h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=`. Proxy/Git manifests
match all releases. Selected, v1.4.0, and v1.6.0 manifests hash to
`8274103d1b27ca97b7411af582aaef13839a7852aca2bf36fad6f01dd6adaae1`,
`bc11a7d4aaa7a131c6beb7811d0e1b09b97334b4c46dc2898bd24ed38eff9d00`,
and `d5645957b5ba37ad5b05344c88bf7a36610a090be01a2e4c91dfc39559a26cf7`.

## Closure, API, Behavior, And Qualifications

All thirteen stable releases are one module and one package with standard
library-only imports; their module files have no Go directive or
requirements. Under exact Go 1.26.7 and contained Go 1.18.10, all verify,
list complete production/test closures, pass native count-1 tests, ten
fresh-process repeats, race, and vet. Complete-test populations total 141 and
97 entries respectively. All 156 production cross-builds and 24 selected/
latest test cross-builds pass; v1.4.0's three fuzz targets pass both SDKs.

Old upstream suites leak the finite random reader from `TestSetRand` between
same-process count iterations through v1.4.0, causing a later EOF panic.
Fresh-process repeats, independent repeats, fixtures, and races pass; this is
an upstream test-order leak, not a production defect. Apidiff from v1.1.2 to
v1.4.0 reports compatible additions only: `NewString`, random-pool control,
invalid-length classification, `NullUUID`, and `UUIDs`.

A 336-line common fixture and 91-line v1.4 fixture cover parsing/formatting,
versions/variants, deterministic MD5/SHA1 names, v1/v2/v4 generators,
entropy errors/short reads, clock sequence/node state, global reuse,
concurrency, serialization, SQL scanning, errors, panics, pools, and new API.
Their SHA-256 values are
`4d38ef9f924e00a583b8129824bb45bf743141cfcad803fefd8871c16c6f0a2f`
and `e8d0a20b664707dc402b64ab8a45609ff40af133891c5f958faa34f2aca5785f`.

`UUID` is a comparable `[16]byte`; Nil formats as the zero UUID. Parsing
accepts canonical, raw, case-insensitive URN, and permissive Microsoft-wrapper
forms; it rejects whitespace and malformed lengths/content but intentionally
does not validate semantic version/variant bits. V4 uses `crypto/rand` with
`io.ReadFull`; error-returning APIs preserve errors while `New`,
`NewString`, and `Must` panic at documented boundaries. Generation is
concurrency-safe absent global `SetRand`/pool reconfiguration.

V1/v2 clock state is mutex-protected and increments the 14-bit sequence for
equal/backward time. Node identity caches the first interface address with at
least six bytes, falls back to random bytes, and copies `SetNodeID` input.
Unix DCE uses uid/gid and Windows uses site-defined behavior. Text/JSON use
canonical strings; binary forms copy 16 bytes. JSON null and nil/empty SQL
scans deliberately do not clear an existing UUID; invalid scans do not
assign. `NullUUID` supplies explicit validity semantics.

Reject v1.5.0 and v1.6.0 because their UUIDv6 generator and `Time()`
extractor do not round trip. The independent 25-line fixture hashes to
`a9400dcb57e716c821d3f22551adc109b4bad060d728c0ab05240dce4098e9f9`
and reproduces this under both SDKs. V1.5.0 also duplicates V7 values for
fixed time/random input; v1.6.0 fixes V7 monotonicity but not V6. Master
`53dda83ebe99c23d0e66c2472abdbf178097c3b8` fixes V6 but is unreleased.

## Project, Vulnerability, And Quality Measurements

Fourteen historical gRPC vertices from v1.33.1 through v1.43.0 declare UUID
v1.1.2. One shortest path is main -> Google Martian/v3 v3.3.2 -> gRPC
v1.37.0 -> UUID. MVS retains requirements from encountered historical
vertices, but package imports never reach UUID; `go mod why -m` is therefore
negative and zero UUID packages load.

The accepted graph exactly matches the v1.4.0 projection: 234 modules, 3,599
edges, 429 complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,067 sum lines, and a 432-line unapplied tidy projection. Only the
exact UUID root edge and authorized sums changed. Every unrelated module and
loaded package is identical. Relative to go-cmp commit `c314bcb`, sums are
+51/-0. The main module remains Go 1.18 with toolchain Go 1.26.7.

Fresh vulnerability data has 1,393 records, Last-Modified
2026-09-09T18:12:50Z, and no UUID advisory. Direct selected, v1.4.0, v1.5.0,
and v1.6.0 scans are zero. Base/candidate project populations are identical
at 30 module findings, 22 Darwin and 23 Windows package findings, and 20
IDs/22 reachable traces per symbol platform. None involves UUID.

Exact Go 1.26.7 project verification/load/build, two repeats, race, vet,
pinned lint, empty-HOME count-2, Linux/Windows builds, compatibility,
launcher, host, snapshot, and Docker acceptance pass. The Go 1.18.10
projection removes only the toolchain line and passes verification/load/
build/vet/cross-builds, 26-package repeats/race, and all 31 compatible shell
tests; full count-1 has only the two accepted closed-file wording differences.

The authoritative 21-stage quality gate kills 80/80 mutants and records all
27 Q0-Q2 rows PASS at L2 with seven improved and zero held, regressed,
non-comparable, or dirty counts. Scorecard SHA-256 is
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The 3,206-entry evidence manifest and decision summary hash to
`0796b997d528e8b9d25d35cbd90a04cc93bf89fbbb3cdacc27fcd46546d00c6d`
and `522b3fc0387507a9612daab87242625d2bfde7aabd74a45ebd24b058d883c31b`.

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
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`;
  apidiff archive SHA-256 is
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
- The final gate used a scratch-only bare-`mktemp` shim, Docker Desktop
  Buildx v0.33.0-desktop.1 copied into empty scratch HOME, and real Python
  3.14 ahead of the system interpreter because Docker emitted 8–9 digit
  fractional seconds. One classified run used `/usr/bin/python3` and stopped
  at Docker timestamp validation; the fresh corrected full run passed.

## Next Objective

Independently evaluate exact-path
`github.com/googleapis/gax-go/v2 v2.7.0` as the next single P7 dependency
group. Do not combine Viper, Cloud Go, Google API, Crypt, or another group.

The shortest current path is main -> Viper v1.15.0 -> GAX v2.7.0.
Historical Cloud Go and Google API family vertices retain lower GAX
requirements. `go mod why -m` is negative and no GAX package loads.

The initial proxy list contains 42 stable versions through v2.24.1. Selected
v2.7.0 is commit `2592e2286ac291a2d9e7c06e1f31b3a6d772131c`, dated
2022-11-02T20:02:53Z, and declares Go 1.19. V2.5.1 appears to be the last
nearby Go-1.18 release, but cannot defeat Viper's v2.7.0 requirement through
an exact root downgrade. Latest v2.24.1 is commit
`269185f57eafcc619f159ffe8857eee6073e955e`, dated
2026-09-03T20:21:21Z, and declares Go 1.25. Resolve the selected-floor
conflict without changing or independently auditing Viper. If no authorized
exact-path selection can satisfy the floor, stop for product/roadmap
direction rather than broadening scope.

That direction is now supplied: use the 2026-09-12 bounded inherited-floor
exception for unloaded v2.7.0. Finish the remaining repository, release,
closure, API, behavior, concurrency, platform, vulnerability, and applicable
project qualification. Preserve Go 1.18, Viper v1.15.0, the selected graph,
and dependency metadata. The documented `ParseError(err, false)` panic is
also accepted only while the complete project load contains zero GAX packages;
revalidate and record that invariant. If GAX becomes loaded or directly
imported, stop for a fresh decision before merge. Any other independent
disqualifying evidence still requires a new stop.
