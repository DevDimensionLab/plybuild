# Quality Upgrade Handover

Generated: 2026-09-06T15:03:39+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only Units commit
  `1c874fda104baa36c3a32aaef2d2d2be9689a942`, exact parent
  `ae5e731121c17338d08dba4c74a99a1c95a36c49`, clean tree
  `744d61deae7d7c3dc25fe1fd1e13a1fc94798981`. It changes only `go.mod` and
  `go.sum`. This documentation handoff must be its direct child; a next
  implementation, if any, must use that documentation commit as exact parent.
- The answered Units archive and sole NEXT `github.com/alecthomas/assert`
  archive link reciprocally. Only launcher mutable regions change during
  handoff. Ordinary and ignored status must end empty.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves are go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, cast v1.5.0 -> v1.5.1, and Units' 2019 ->
  2024 pseudo-version. Their recorded minimal closures select testify
  v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty v0.3.1, and
  rogpeppe/go-internal v1.9.0, and add exactly 17 checksum lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after Go 1.26.7 and bounded groups
through accepted Units. All earlier accepted, rejected, and no-change
decisions and every evidence correction remain final. P8 remains queued.

## Accepted Alecthomas Units Group

The Go proxy has no stable semantic-version or prerelease tags for
`github.com/alecthomas/units`. Proxy `@latest`, exact Go `@latest`, and exact
Go `@master` all resolve
`v0.0.0-20240927000941-0f3dac36c52b`, commit
`0f3dac36c52b29c22285af9a6e6593035dadd74c` at
2024-09-27T00:09:41Z. This is the canonical latest unreleased default-branch
pseudo-version, not a stable release. The selected old pseudo-version is
`v0.0.0-20190717042225-c3de453c63f4`, commit
`c3de453c63f4bdb4dadffab9805ec00426c505f7` at
2019-07-17T04:22:25Z.

The authoritative repository is enabled, unarchived, and not a fork; default
branch is `master`. It has no tags, GitHub releases, retractions, alternate
module paths, or module deprecation. A newer `renovate/all-minor-patch` head
resolves a 2026 pseudo-version but is an unreleased non-default branch; pull
request head `7355547fc901` is unmerged and not proxy-resolvable. Neither is a
qualified exact-path candidate.

The old module declares only its module path. Candidate `go.mod` declares Go
1.15 and requires `github.com/stretchr/testify v1.9.0`; the project already
selects that module and it declares Go 1.17. The complete changed closure is
therefore compatible with the retained Go 1.18 floor by declaration.

Old checksum pair:
`h1:Hs82Z41s6SdL1CELW+XaDYmOH4hkBN4/N9og/AsOv7E=` /
`h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=`. Candidate pair:
`h1:mimo19zliBX/vSQ6PWWSL9lK8qwHozUj03+zLoEB8O0=` /
`h1:fvzegU4vN3H1qMT+8wDmzjAcDONcgo2/SZ/TyfdUOFs=`. Proxy ZIP SHA-256 values
are `5f9f0ba0037b25179cf0a4bf52ab2c9981dbb366125faf8ed44beaf6a654dc5e`
and `9a275dbb1454d52d2b868b990be424c7d31c0acf2c4a99850faf18ecb9f48b91`.
All old eight and candidate ten proxy files match their exact upstream
commits; source manifests are
`1c1119e80d757ff028a0f7768dcfa9df181a31930519ac6f8e4f9c629f4c517b`
and `a2579e87929e9f8da86788fcb05806a88204cb7084082ac345280dbb7f14e1c3`.
GitHub reports valid signatures for both commits. There are no tag objects;
local `%G?` is `N` and local signature verification cannot use the sandboxed
OpenPGP home.

External exact get changes exactly Units. Modules stay 234, graph edges move
3,564 -> 3,566 through main -> candidate and candidate -> already-selected
testify v1.9.0, complete packages stay 429, and go.sum moves 1,031 -> 1,033
with exactly the candidate checksum pair. The unapplied tidy projection moves
332 -> 341 lines and was never used as implementation. Units loads in zero
old/new complete packages; repository Go has zero imports and `go mod why`
says the main module does not need it. It is historical unloaded MVS graph
debt, not a main-package or dependency-test consumer.

Candidate proxy and upstream sources pass verify/list/count-1/count-10/race/
vet and remain source-clean. The old source omits test requirements, so its
exact external test-only apparatus adds assert through
`go get github.com/stretchr/testify/assert@v1.9.0`, selecting testify v1.9.0,
go-spew v1.1.1, go-difflib v1.0.0, and YAML v3 v3.0.1. Both old sources then
pass the same suite and remain clean. This apparatus is not project MVS
closure.

Old, projected, and committed repository verification, build, complete tests/
race/vet, Windows build, pinned lint, and API/CLI/help contracts pass; public
help streams and reports are byte-identical. Empty-HOME count-2 passes. Primary
govulncheck v1.7.0 data updated 2026-09-02T19:12:04Z contains 1,392 records and
no Units record. Old/candidate exact IDs and normalized traces are identical,
preserving 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol populations.

Exact `go get` produced the implementation commit recorded above with three
insertions. Exact `make quality` exits 0 at `/private/tmp/ply-p7-uq11`; its
21-stage ledger attains L2 for all 27 Q0-Q2 rows with 80/80 mutations killed,
8/8 mutation and 4/4 acceptance populations, valid manual evidence, passing
host/snapshot/Docker acceptance, and zero held, regressed, not-comparable, or
dirty counts. Q0-Q2 scorecard SHA-256 is
`ec538e05ab0c0a202ae53458bb16b507c3c8877e3d6a99fe3aaa60b794df1a1b`.
The fresh snapshot and Docker report SHA-256 values are
`2da1b429b760563cca641133745c2453390363e9328654a216904c7c8fe53666`
and `439cadfaa5a43310edb81b8d07310801a788d5f68ef7a7de9bae7d20bd934620`.
The separate full audit exits expected 1, never 2, only for established queued
Q3.1, Q3.3, Q3.4, and Q3.7; its scorecard SHA-256 is
`860ed0e2077736e1eae14a38660bf11a0c178f847284eb12c31056e770b729bf`.

Several retained aggregate attempts exposed the established launcher signal-
retention timing fixture before preflight. No apparatus changed. Putting the
already verified external Python 3.14 earlier in PATH made the report parser
finish before that race; exact Go remained first among Go installations.
Docker buildx discovery required operator HOME, followed by acceptance with a
fresh empty Docker configuration.

## Units Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-units-selection.ae5e731.smI4xp`. Its fully verified
  72,950-entry manifest SHA-256 is
  `2cdb64e74a40304c86ef35f6b8e479c49190d4c2d0e6678770bee4e6486493ec`;
  decision-summary SHA-256 is
  `e30405aeac0eecf6356570f8a49677de1715ad9a7aee15e8efdc9739f55c4cd0`.
- Exact quality evidence at `/private/tmp/ply-p7-uq11` has 237,860 entries and
  verified manifest SHA-256
  `e8b4fdc59ec12c42afa21286939a1a5ca28a4fef29c88655f73890b580564328`.
- Schema-2 manual evidence SHA-256 is
  `3bf267b1fe55ffba41e5d4024ceaa9a7fc8a2c8d713d00dcfc08af19f4743978`.
  Its fresh governed-source digest covers 112 identical files and has SHA-256
  `a716b4ddc8918db9bde0f912c317df394ed68f15d8396ec7c7b9ef3b2e89a0e6`.
- The inherited 37-root verification table is clean and has SHA-256
  `8669ff369a82ee4dd93a56f005df16a1b9327bc3bd477ef1f10101872c32b2b8`.
  Preserve the recorded go-colorful mutable telemetry, btree mutable cache/
  HOME regression, cast NUL-delimited whitespace-path handling, source-
  archive normalization, Check history table, Errgo preflight runner, Resty
  committed projection/signal fixture, and Kingpin external-report correction.

## Retained Decisions And Tools

- Retain exact selections `gopkg.in/alecthomas/kingpin.v2 v2.2.6`,
  `gopkg.in/resty.v1 v1.12.0`, `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Preserve all other prior bounded decisions.
- Kingpin, Resty, Errgo, Check, YAML v3, YAML v2, x/text, x/net, x/image,
  gotenv, and jwalterweatherman evidence roots remain verified at their exact
  counts and manifest hashes recorded in their answered archives. Cast's
  selection/review/quality/regression roots remain immutable. Do not revisit
  them or reinterpret their corrections.
- Exact Go 1.26.7 remains at
  `/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  Put its directory first in PATH, set GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery root
  `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` remains verified at two
  entries and manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Retain verified golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected `github.com/alecthomas/assert`
`v0.0.0-20170929043011-405dbfeb8e38` as exactly one bounded P7 module group.
The fresh post-Units survey reports v1.0.0 at 2021-12-01T05:52:01Z, but treat
canonical latest, release qualification, declarations, complete floor-
compatible closure, source identity, signatures, repository state, loaded
population, real consumers, self-tests, and vulnerability data as unknown
until independently proved from primary evidence.

Current measurements are 234 selected modules, 3,566 graph edges, 429 native
complete-test packages, 1,033 go.sum lines, a 341-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. If selected is already the
decision, do not manufacture a requirement or commit. Do not combine Assert,
Units, Kingpin, another dependency, or P8.
