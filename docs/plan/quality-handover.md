# Quality Upgrade Handover

Generated: 2026-09-08T11:09:36+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is exact XXHash v2.3.0 commit
  `e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
  `a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
  `5213ba55981d77d7c8061312915e237e80d29af8`. It changes only `go.mod`
  and `go.sum` with three insertions.
- OpenCensus Proto and Crypt were retained without implementation commits.
  Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor.
  Earlier Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` also remain ancestors.
- The answered XXHash archive and sole NEXT Logex archive must link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt and OpenCensus
Proto selections. P8 remains queued. All earlier acceptances, rejections,
no-change decisions, evidence corrections, and lifecycle ancestry are final;
do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*`
roots, never run `go mod download all` in a measured tree, and preserve the
launcher's scratch cleanup and reciprocal archive contract.

## Accepted XXHash V2 Group

Accept exact-path `github.com/cespare/xxhash/v2 v2.3.0`, canonical latest,
over historically selected v2.1.2. The exact proxy lists only stable v2.0.0,
v2.1.0, v2.1.1, v2.1.2, v2.2.0, and v2.3.0. Exact latest is dated
2024-04-04T20:00:10Z. No version is retracted and no module deprecation
directive exists. Selected and latest declare the exact v2 path and Go 1.11.
The incompatible v1 module and nested benchmark module are distinct paths.

Go-import identifies `https://github.com/cespare/xxhash.git`. The public
repository is enabled, unarchived, non-fork, and defaults to `main`. All six
v2 tags are annotated and unsigned; their commits are unsigned; there are no
GitHub Release objects.

Selected tag object `7ae26c41ed6fb1f8a6c21e05eeff4d91b5e401c4`
dereferences commit `e7a6b52374f7e2abfb8abb27249d53a1997b09a7`, tree
`78bbde5aae5fc0324314789c69f39da97d995f60`, parent
`3b9a65d476075116e7baf9ca3c37294f637ea033`. Candidate tag object
`7438b35f14d771ee32d8bbcd9527d32a336e7dad` dereferences commit
`998dce232f17418a7a5721ecf87ca714025a3243`, tree
`9a415e99332a4ef14d2d590e9ed748c99ad4bc70`, parent
`21fc82c137a876186c8acb0349e941ddb280bf03`. Default-branch head
`ab37246f3853501fb3e16d199556315b50889ad2` has only one later CI-only
commit and is not a release.

Every proxy ZIP matches exact Git source after excluding the nested benchmark
module. Normalized manifest SHA-256 values for v2.0.0, v2.1.0, v2.1.1,
v2.1.2, v2.2.0, and v2.3.0 are respectively
`04cbb9e84ae2b86cd24a46b500e0708e159c3e006fc70133c3d6cd6698d58032`,
`9667f1743cafabff992f020509d2c1538573540bb075b0dc80d0b83d97e965d1`,
`daab817c6c64cab7ae3126847d6d132fbb5aad7b7ea4ad3ef763174b182f84bc`,
`69a8ba25b194eb9b2c195ac76f8f81ab03595cecebf0699a59b244611591352f`,
`2b04a8ff379c136d984046f4db44a3dd2c057595869646190b009023620a796a`,
and `b604c7ab37f70e52a009950fbde4a6708f7372751450549f6e63aff0e7b8bd09`.
The answered archive records all checksum-database pairs.

Selected and candidate have no root requirements and preserve the Go 1.18
floor. Their three packages pass verify/list/count-1/two count-10/vet and
pure-Go tests under exact Go 1.26.7 and contained Go 1.18.10. Root race tests
pass. Dynamic normal repeated tests pass; the invariant non-race plugin built
by TestMain cannot load into a race host and is not a candidate regression.
The nested benchmark module has four modules, tops out at Go 1.13, and passes
its count/race/vet closure under both SDKs.

Independent known-vector and streaming fixtures pass all constructors,
checksum helpers, digest methods, reset, size/block size, Hash64, and binary
marshal interfaces; v2.3.0 seeded behavior also passes. Pinned apidiff reports
only compatible `NewWithSeed` and `(*Digest).ResetWithSeed` additions.
V2.3.0 adds Darwin-arm64 assembly; both versions retain Windows/Linux amd64
assembly and pure-Go/appengine/non-gc fallbacks. Cross-builds and static
Linux-amd64 execution pass under both SDKs.

Ply reaches Viper v1.10.1 through mvn-pom-mutator v0.2.3. Viper and Crypt
v0.4.0 request XXHash v2.1.2; grpc v1.43.0 and Prometheus client v1.4.0 request
lower v2.1.1. Ply imports and loads no XXHash package. Viper/Crypt tests do not
load it. The actual historical consumer is grpc xDS ringhash/resolver using
`Sum64String`; focused count-1/count-10/race/vet pass against selected and
candidate under both SDKs.

## Projection, Quality, And Vulnerability Measurements

The accepted project is 234 selected modules, 3,583 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines and removes zero.

Exact selected-version get changes no selection and adds only a redundant
root requirement/main edge/full source checksum; it was not applied. Exact
v2.3.0 get changes only XXHash, adds only main -> v2.3.0, and adds only its
checksum pair. Tidy would remove the pin/pair and restore inherited v2.1.2
while preserving pre-existing tidy debt.

The projected and committed candidate pass mod verify, fresh online then fully
offline package-list replay, build, count-1/count-10/race/vet, Windows-amd64
build, pinned lint, byte-identical help, API/CLI compatibility, CLI surface,
and empty-HOME count-2.

Exact `make quality` passes all 21 ordered stages. All 27 Q0-Q2 criteria pass
at L2, 80/80 mutants are killed, acceptance is 4/4, six manual receipts are
valid, and held/regressed/not-comparable/dirty counts are zero. Its scorecard
SHA-256 is
`579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.
Snapshot/Docker report hashes are
`1c1a34fa6fb82d518daa834e326ceaf56bf19c12011db9b75985c83baf3270d6`
and `4112e3de3237295070697e3a3ee55267ddfc6389bc9c5d4b4c1086e5b8e13ecf`.
The separate audit meta-suite passes 15 controls. Full audit exits 1 only for
queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7, attains L2, and never exits 2.

Fresh vulnerability data has 1,392 module records and no XXHash record.
Baseline/candidate results are identical: 20 IDs/22 reachable traces for
Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin
module IDs. No target module, package, symbol, or trace appears.

XXHash evidence has 608 verified entries; manifest SHA-256 is
`6f57be6cbf177c6617ebf62e0eea7b38f6b3b41eb00873060e160ddf02f1cf34`.
Decision-summary SHA-256 is
`4e2f0559dbefd98f22e8f1efe34a1e538ecb8e1dd2da0749af8478c88c0a1a20`.

## Tools And Corrections

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Portable receipts remain golangci-lint 2.12.2 archive `a9c54498…`,
  GoReleaser 2.17.1 binary `f5f08a77…`, and apidiff `0c55d9e3…`.
- Scratch-only corrections supplied an explicit BSD `mktemp` template,
  contained Xcode SDK/compiler resolution, and Python 3.14 before Apple Python
  3.9 for Docker 29's eight-digit fractional timestamp. A standalone fresh
  Docker pass proved the correction before the full gate reran. Superseded
  attempts are retained but are not proof. No measured source was changed.
  Evaluation images were removed and can be reproduced from retained scratch
  inputs.

## Retained Decisions

Retain accepted XXHash v2.3.0 and Speakeasy v0.2.0 moves and all earlier exact
selections/outcomes, including Repr v0.5.4, Assert v1.0.0, Units
`0f3dac36c52b`, Chroma v0.10.0, Colour v0.1.0, Template
`fb15b899a751`, Optional v1.0.0, Circbuf `5111143e8da2`, Consul API
`eb2c6b5be1b6`, Go Metrics v0.4.0, Go Radix v1.0.0, Perks v1.0.1, Kingpin
v2.2.6, Resty v1.12.0, Errgo v2.1.0, Check `41f04d3bba15`, YAML
v2.4.0/v3.0.1, and `go.yaml.in/yaml/v3 v3.0.5`. Preserve accepted Kong and
all roadmap/archive decisions.

## Next Objective

Independently evaluate selected exact-path `github.com/chzyer/logex v1.2.1`
as the next single P7 group. The exact proxy lists v1.1.1 through v1.1.10,
v1.2.0, and v1.2.1; selected is already canonical latest at
2022-04-24T13:13:51Z and declares Go 1.15. Repository tags v1.0/v1.1 are absent
from the proxy list and need precise qualification.

Readline v1.5.1 and chzyer/test v1.0.0 request selected; Promptui v0.9.0 and
historical pprof versions request lower v1.1.10. Ply reaches Readline through
Promptui, but Logex is in Readline's test closure rather than Ply's loaded
packages.

The public `github.com/chzyer/logex` repository is enabled, unarchived,
non-fork, and defaults to `master`. Selected is lightweight tag/commit
`2f95bdde8c3c97bfbf6d016fcc410669a895b9e7`, tree
`36fcd9ac7d56d659872b2a6576ca6f66385c938a`, parent
`a21c317abc1e9a4f23ed3455107a4d20375735cc`. Master
`5a7e37d2e8a8bbe3ef54984ab949eebaa948b8b4` contains five later test/CI
commits and is not a release. Independently establish release/signature
identity, complete closure floor, native behavior, historical consumers, exact
projection, and vulnerability effect. Do not implement or combine another
group until that bounded decision is complete.
