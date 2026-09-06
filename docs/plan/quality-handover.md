# Quality Upgrade Handover

Generated: 2026-09-06T20:17:39+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest dependency implementation is Repr commit
  `6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
  `53475076e1c79f6d2181877e3238a1a5d389c246`, tree
  `82e7b1f5c503659082207481b8339e8113e38c10`. It changes only `go.mod` and
  `go.sum`. The operator-authorized lifecycle repair is layered on top before
  the Repr decision handoff; preserve it and account for that ancestry when
  finalizing the active session.
- The answered Assert archive and sole NEXT `github.com/alecthomas/repr`
  archive link reciprocally. Only launcher mutable regions change during
  handoff. Ordinary and ignored status must end empty.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves are go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, cast v1.5.0 -> v1.5.1, Units' 2019 -> 2024
  pseudo-version, Assert's old pseudo-version -> v1.0.0, colour's 2016
  pseudo-version -> v0.1.0, repr's 2018 -> 2021 pseudo-version, and go-diff
  v1.0.0 -> v1.2.0. Their recorded minimal closures add exactly 25 checksum
  lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Operator Lifecycle Recovery

On 2026-09-06, accumulated reproducible agent scratch, audit trees, Go caches,
and Ply Docker layers had filled the data volume. The operator authorized their
deletion and a lifecycle repair. Cleanup freed approximately 763 GiB, moving
the data volume from 17 GiB free (99% used) to about 780 GiB free (13% used).
Historical digests and decisions below remain records, but referenced external
`/private/tmp/ply-*` evidence and tool roots were intentionally purged.

Every successor must use `$CODEX_SESSION_SCRATCH_ROOT` for disposable state.
The launcher now removes per-turn scratch on success, failure, and interruption,
and removes supervisor logs after successful completion. Mutation workspaces,
snapshot caches, Docker tags, run-labeled builder images, and acceptance
containers clean themselves on exit. Preserve these operator-authorized
lifecycle controls while completing the active Repr handoff.

P2A-P6 are complete. P7 remains active after Go 1.26.7 and bounded groups
through accepted Assert. All earlier accepted, rejected, and no-change
decisions and every evidence correction remain final. P8 remains queued.

## Accepted Alecthomas Assert Group

The exact `github.com/alecthomas/assert` path has one proxy-listed stable tag,
v1.0.0, and no exact-path prerelease. Proxy `@latest` and exact Go `@latest`,
`@v1`, and `@master` all resolve v1.0.0 at
2021-12-01T05:52:01Z. It is therefore canonical latest and the highest
qualified stable exact-path release. Selected old pseudo-version
`v0.0.0-20170929043011-405dbfeb8e38` maps to commit
`405dbfeb8e38effee6e723317226e93fff912d06`, tree
`188d96b077666f551bf61f7d184486a7fa416e70`, at
2017-09-29T04:30:11Z. Candidate tag v1.0.0 is lightweight commit
`73444aca37e09619baa21040dad4527f857f9abb`, tree
`d09edea32b1be9b1e73428797782d491949b50b7`, at the reported proxy time.
Both commits and the tag are unsigned.

The authoritative repository is enabled, unarchived, and not a fork, with
default branch `master`. It has no GitHub Release objects, v1 retractions, or
module deprecation marker. Stable v2.11.0 at 2024-09-17T06:12:02Z and master
pseudo-version `v2.11.1-0.20251014091832-afff49140ce5` belong to distinct
module path `github.com/alecthomas/assert/v2`; exact v1-path queries for v2
tags, commits, and prereleases fail the semantic-import-version check. They are
not v1 candidates. Relative to the selected commit, candidate v1.0.0 adds a
license filename change and module/Hermit metadata but no production Go source
change.

Old checksum pair:
`h1:smF2tmSOzy2Mm+0dGI2AIUHY+w0BUc+4tn40djz7+6U=` /
`h1:r7bzyVFMNntcxPZXK3/+KdruV1H5KSlyVY0gc+NgInI=`. Candidate pair:
`h1:3XmGh/PSuLzDbK3W2gUbRXwgW5lqPkuqvRgeQ30FI5o=` /
`h1:va/d2JC+M7F6s+80kl/R3G7FUiW6JzUO+hPhLyJ36ZY=`. Proxy ZIP SHA-256 values
are `873d257170b1363142cbf5e16b49c6a21cccb3e4aaceb9d370c3b78b051a5663`
and `be346d2847db5cfc7e40babb0b3b7062fa20e76af366652eb4f65710ea7d5fdc`.
All ten old files match upstream. All 17 candidate regular files match after
normalizing the proxy's documented omission of three tracked Hermit symlinks;
the equal normalized manifest SHA-256 is
`3e4aa8146850de3d14cd63ac1640afb3f7da28d5a806cead0895521a3f8ca08d`.

Candidate declares Go 1.17 and requires colour v0.1.0, repr pseudo-version
`v0.0.0-20210801044451-80ca428c5142`, and go-diff v1.2.0, plus lower
isatty/x/sys versions. Exact get advances those four selections. Colour has no
Go declaration, repr declares Go 1.15, and go-diff declares Go 1.12. Existing
selected requirement targets declare no more than Go 1.18, so the complete
closure preserves the retained floor by declaration.

The closure's exact source identities are lightweight unsigned colour tag
commit `a1c6bd85eba7190e4d2959ecd15831d0a25b37b9`, unsigned repr pseudo-version
commit `80ca428c51421b9f0ceedd9218af5e1068cd8153`, and annotated unsigned go-diff
tag object `b292a3123758b064eeaa3a5aa86df1adca0f4401` peeling commit
`0a651d56613f9de4bed8b9c4769b776ef168bfca`. Their proxy ZIP SHA-256 values
are `74d51002731fa104943b62ee11fb61b14c517e75a4a3983bfb03976b6c75349b`,
`e92498fa15fbef295ef530c6ae96d17f446e8f3403c583dd721c73b8da24174d`,
and `da1accb73e9ac304a805eb59fba2c50d0089f9206a2574b8812af7e75e8ec105`.
All regular files match exact upstream commits; repr normalization accounts
for five tracked symlinks omitted by the proxy.

Exact get keeps 234 selected modules and 429 complete packages. Graph edges
move 3,566 -> 3,580 through four main requirements, five Assert requirements,
and five go-diff requirements, with no removal. Every requirement target was
already selected, so there is no fifth changed selection or new module. Go.sum
moves 1,033 -> 1,041 with exactly four checksum pairs; the unapplied tidy
projection moves 341 -> 354 lines and was never used as implementation.

Assert, colour, repr, and go-diff load in zero old/new complete packages and
repository source imports none of them. `go mod why -m` says the main module
does not need Assert. Its historical graph path is main -> go-term-markdown
v0.1.4 -> Chroma v0.7.1 -> the old Assert pseudo-version. This is unloaded
MVS graph debt, not a main-package or dependency-test consumer, so there are
no actually used Assert packages or symbols to focus-test.

Candidate proxy and exact upstream source independently verify, list one
native package, and pass count-1, count-10, race, and vet without mutation.
That package has no native test files; `_example/example_test.go` is excluded
by the Go `./...` rule. No separate test-only requirement apparatus is needed.
Old, candidate, and committed repository verification, build, complete tests/
race/vet, Windows build, pinned lint, and help/API/CLI contracts pass; public
help streams and reports are byte-identical. Empty-HOME count-2 passes all 27
packages. Primary govulncheck v1.7.0 data updated 2026-09-02T19:12:04Z has
1,392 records and no Assert or closure record. Old/candidate exact IDs and
normalized traces are byte-identical, preserving 20/30/20 Darwin-symbol/
Darwin-module/Windows-symbol populations.

Exact `go get github.com/alecthomas/assert@v1.0.0` produced the implementation
commit recorded above with 12 insertions. Exact `make quality` exits 0 at
`/private/tmp/ply-p7-assert-quality.0781fd6.q1`; its 21-stage ledger attains L2
for all 27 Q0-Q2 rows with 80/80 mutations killed, 8/8 mutation and 4/4
acceptance populations, valid manual evidence, passing host/snapshot/Docker
acceptance, and zero held, regressed, not-comparable, or dirty counts. Q0-Q2
scorecard SHA-256 is
`3904a989cb3857a3c07e9c55f0815865231dfef1393753ab906bc852ed366fb2`.
Snapshot and Docker report SHA-256 values are
`56b3d9b35d65e551ac705cf261da82363a65c8cc08b34bafa8f34287a0cdbe1d`
and `f7961c8e122df69118deb368571ae34c9227da02eaf8f51de23a3744916d0d00`.
Standalone audit meta and the six-row focused audit pass. The separate full
audit exits expected 1, never 2, only for established queued Q3.1, Q3.3, Q3.4,
and Q3.7; its scorecard SHA-256 is
`ca4a8411d0ae57daed70c9a0dbddb73992e35eac10606c615a8604b79da3de56`.

The first external API comparison lacked its independently warmed exact
v1.0.1 compatibility base; the corrected fresh base archive and external
cache passed old, candidate, and committed gates. The first supplemental
audit-meta and empty-HOME invocations named an external TMPDIR before creating
it; their setup-only diagnostics are retained, while corrected runs pass.
Neither correction changed the repository or quality apparatus.

## Assert Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-assert-selection.8f24778.fEzhH8`. Its fully verified
  46,700-entry manifest SHA-256 is
  `072c8a42889b85e1fa86d0e5701db7f31f8840c255f06de5a7633403db621782`;
  decision-summary SHA-256 is
  `bb13fa01a123a32399806db46f7e61813e6549b7ca1c5a42c25a09254684d17f`.
- Exact quality evidence at `/private/tmp/ply-p7-assert-quality.0781fd6.q1`
  has 237,950 entries and verified manifest SHA-256
  `04b202f6084740762589a042c8847c2411da5a54ee6572d3fa54790c4a3301e1`.
- Schema-2 manual evidence SHA-256 is
  `fa33df10f95a2982137186339c7f5b586b3ec2a38d9f7500dcf8a8464004a5d2`.
  Its fresh governed-source digest covers 112 identical files and has SHA-256
  `793dba73deb0c45bd9d2b6a8d29a4a19b1cc71993b7243229a056243224bd33b`.
- The inherited 37-root verification table is clean and has SHA-256
  `8669ff369a82ee4dd93a56f005df16a1b9327bc3bd477ef1f10101872c32b2b8`.
  Preserve the recorded go-colorful mutable telemetry, btree mutable cache/
  HOME regression, cast NUL-delimited whitespace-path handling, source-
  archive normalization, Check history table, Errgo preflight runner, Resty
  committed projection/signal fixture, Kingpin external-report correction,
  Units launcher timing correction, and Assert external-cache/TMPDIR setup
  corrections.

## Retained Decisions And Tools

- Retain exact selections `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Preserve all other prior bounded decisions.
- Kingpin, Resty, Errgo, Check, YAML v3, YAML v2, x/text, x/net, x/image,
  gotenv, jwalterweatherman, Units, and Assert records remain verified at their
  exact counts and manifest hashes in their answered archives. Their external
  roots were intentionally purged; do not reinterpret the recorded decisions.
- Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, SHA-256
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  Put its directory first in PATH, set GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Its former recovery root
  was purged after verification at two entries and manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Recreate the following tools beneath `$CODEX_SESSION_SCRATCH_ROOT` as needed
  and verify their recorded hashes: golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Resume and independently re-verify the bounded `github.com/alecthomas/repr`
decision underlying dependency-only commit `6f4d02e`. The prior disposable
evidence was purged, so recreate only the evidence needed for this decision
beneath `$CODEX_SESSION_SCRATCH_ROOT`. Treat canonical latest, release
qualification, declarations, complete floor-compatible closure, source
identity, signatures, repository state, loaded population, real consumers,
self-tests, and vulnerability data as unknown until independently proved from
primary evidence, then complete or reject the Repr handoff.

Current measurements are 234 selected modules, 3,580 graph edges, 429 native
complete-test packages, 1,041 go.sum lines, a 354-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7.

Retain the changed selection only if it is floor-compatible, has an explained
minimal closure, and passes every applicable gate. Do not combine Repr, Assert,
Units, another dependency, or P8.
