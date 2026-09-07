# Quality Upgrade Handover

Generated: 2026-09-07T12:55:45+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest dependency implementation is Repr commit
  `6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
  `53475076e1c79f6d2181877e3238a1a5d389c246`, tree
  `82e7b1f5c503659082207481b8339e8113e38c10`. Exact Go 1.26.7 `go get`
  changed only `go.mod` and `go.sum`.
- Operator-authorized lifecycle repair commit
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`, parent `6f4d02e`, tree
  `81ca3d7a6537ce2b9c26337c2bf3d4d67c28223a`, is intentionally between the
  dependency commit and this handoff. Preserve it. The Repr documentation
  handoff must be its direct child.
- The answered Repr archive and sole NEXT Kong archive link reciprocally.
  Only launcher mutable regions change during handoff. Ordinary and ignored
  status must end empty.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  changes now add exactly 27 checksum lines. The previously recorded groups
  remain final; Repr additionally advances its 2021 pseudo-version to v0.5.4
  and adds exactly its checksum pair.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, external source, report, evidence tree,
and build context must remain beneath `$CODEX_SESSION_SCRATCH_ROOT`. The
launcher removes per-turn scratch on success, failure, and interruption and
removes supervisor logs after successful completion. Preserve the operator-
authorized lifecycle controls and scratch-only mutation, snapshot, Docker,
and acceptance cleanup.

The historical `/private/tmp/ply-*` evidence and tool roots were deliberately
purged after their receipts were recorded. Treat those digests as historical
records, not live roots. P2A-P6 are complete. P7 remains active after exact Go
1.26.7 and bounded dependency groups through accepted Repr v0.5.4. All prior
accepted, rejected, and no-change decisions and evidence corrections remain
final. P8 remains queued.

## Accepted Alecthomas Repr Group

The exact `github.com/alecthomas/repr` proxy path lists ten stable semantic
tags, v0.1.0 through v0.5.4, and no prerelease. Proxy `@latest` and exact Go
`@latest`, `@v0`, and `@master` all resolve stable v0.5.4 at
2026-07-15T12:04:01Z. It is a lightweight tag at unsigned commit
`9b3680b9bb4e172c4fe539347ac5d0ddf5de0aa9`, tree
`d5c923ba26338b41955d8fc02e88948c460777a0`, parent
`7c84ded4f40f0246ceaf9d40eea3161b5b354a57`; author time is
2026-07-15T12:03:31Z. This is canonical latest and the highest qualified
stable exact-path release.

The selected pseudo-version
`v0.0.0-20210801044451-80ca428c5142` maps to unsigned commit
`80ca428c51421b9f0ceedd9218af5e1068cd8153`, tree
`55841124de15dc614a54c0ebe5a22274c43f4218`, parent
`1a2716c6f5633eecd218817f6b3f3c6189b43b24`, at
2021-08-01T04:44:51Z. The enabled, unarchived, non-fork repository has default
branch `master`, no GitHub Release objects, retractions, deprecation marker,
alternate module path, or v2 tags. All ten tags are lightweight commit refs.
A newer Renovate branch commit `4e78e327210ff565537c3f04eb0a7c24e6f754da`
resolves only as unreleased pseudo-version
`v0.5.5-0.20260902080120-4e78e327210f`; it is neither the default-branch head
nor a qualified release and is not the exact Go latest decision.

Old checksum pair:
`h1:8Uy0oSf5co/NZXje7U1z8Mpep++QJOldL2hs/sBQf48=` /
`h1:2kn6fqh/zIyPLmm3ugklbEi5hg5wS435eygvNfaDQL8=`. Candidate pair:
`h1:OVP7JEcuzU9CCDsT6STCr3rg17oQfWILtPWd2EG0uN4=` /
`h1:Fr0507jx4eOXV7AlPV6AVZLYrLIuIeSOWtW57eE/O/4=`. Independent checksum-
database lookup agrees. Old and candidate proxy ZIP SHA-256 values are
`e92498fa15fbef295ef530c6ae96d17f446e8f3403c583dd721c73b8da24174d`
and `54d756f1b363c943c3c2981f48e199405ef3c1a554715260c886114c990a5492`.
All proxy regular files match their exact upstream commits after accounting
for the proxy's omission of five tracked symlinks. Equal normalized manifest
SHA-256 values are
`006d6ed03c0bb9b15f59c303512738655daec6d051f4ede7065884a9a0e731f8`
old and
`d2f7117e74b6e8f9f4da49eb81773292efdbabfbfce135ec4abb375251c0bf6c`
candidate.

Candidate v0.5.4 declares Go 1.18 and has no requirements. The selected old
version declares Go 1.15 and also has no requirements; v0.1.0 declares Go
1.15 and every later stable tag declares Go 1.18. The complete changed closure
is therefore exactly Repr itself and preserves the retained Go 1.18 floor by
declaration. Relative to the selected commit, v0.5.4 spans 27 commits, 11
paths, 521 insertions, and 111 deletions. Relevant release history includes
the v0.5.3 visited-set recursion repair and v0.5.4 repeated-scalar
preservation.

Exact replay keeps 234 selected modules, 3,580 graph edges, and 429 native
complete-test packages. Exactly one selection and one existing graph line
change: the main requirement is relabeled from the old Repr version to v0.5.4.
There are no added modules or dependency edges. Go.sum moves 1,041 -> 1,043
lines by adding only the candidate checksum pair and retaining the old pair.
The unapplied `go mod tidy -diff` projection moves 354 -> 356 lines; tidy was
not used as implementation. The replayed `go.mod` and `go.sum` are byte-
identical to dependency commit `6f4d02e`.

Repr loads in zero old and candidate complete packages. Repository Go source
has no Repr import, `go mod why -m` says the main module does not need it, and
an explicit `go list -deps -test ./...` still returns 429 packages with no
Repr package. It is unloaded historical MVS graph debt, not a main or
dependency-test consumer, so there are no actually used Repr symbols to
focus-test.

Candidate proxy and exact upstream source independently pass module verify,
list, count-1, count-10, race, and vet without source mutation. There is one
native package and no separate test-only module apparatus. Repository module
verify, build, complete tests/race/vet, Windows amd64 build, pinned lint,
empty-HOME count-2, launcher/Make/preflight contracts, and byte-identical
public help/API/CLI reports all pass.

Fresh primary Go vulnerability data was updated
2026-09-02T19:12:04Z and contains 1,392 records with no Repr record. Exact old
and candidate IDs and normalized traces are byte-identical: 20 Darwin symbol,
30 Darwin module, and 20 Windows symbol findings.

Exact `make quality` exits 0 under the scratch-only final evidence tree. Its
21-stage ledger has all 27 Q0-Q2 rows at L2, 80/80 mutations killed across
8/8 subjects, 4/4 acceptance flows, fresh host/snapshot/Docker acceptance,
valid schema-2 manual evidence, and zero held, regressed, not-comparable, or
dirty counts. Q0-Q2 scorecard SHA-256 is
`229878b617f3c1161a0864989aeb929c99c661facce4e87cb18b348399363e84`;
snapshot and Docker report SHA-256 values are
`05639217539e23340d8eb8421f3636f5300ab8b9e95aa545bf371edf06a0add2`
and `9faad0c5ad16551be6edc6267ff41f433e0f28503ebab6b6a6fca2388fb685bc`.
The separate full audit exits expected 1, never 2, only for established queued
Q3.1, Q3.3, Q3.4, and Q3.7; its scorecard SHA-256 is
`ab3ecb65fc362efc8fe37a86942e3e213fc469efbfd3194972757a8db9b747ed`.

Fresh quality setup required separately warming the API base and current
module caches, a scratch-local adapter because BSD bare `mktemp -d` ignores
TMPDIR, and an uninherited Docker configuration so Docker could find its
bundled buildx plugin. Three earlier setup diagnostics and one environment-
setup diagnostic are retained; the corrected clean run passes without any
repository mutation. A first manifest loop reused zsh's special `path`
variable and therefore cleared PATH; it produced no usable receipt. The
corrected neutral loop produced and verified the receipt below. A later
verification invocation began in the repository rather than the manifest's
scratch-relative root; rerunning from the sealed root verified all entries.

## Repr Evidence And Tools

- Decision evidence beneath the active scratch root has 447 entries. The
  fully verified evidence-manifest SHA-256 is
  `aba1a6ac45d6b8a51a664e0d11f40ccccf82017240c91beaa59ec36866e26711`;
  decision-summary SHA-256 is
  `900e96c7001a89857ca7c16e1f2572ea0b90db9f2ab4d4e4e6a1c3153808bbcf`.
- Schema-2 manual evidence SHA-256 is
  `b89a970369d7867799a33a069a41c7869ed4361b81d6d8ef3744e8f455dcf737`.
  Its governed-source digest covers 303 paths and has SHA-256
  `9f90139973eb2efe51e372af638ebcc857e1320a8430f672ed64a392382ec882`.
  Q1.9 is an honest aggregate receipt with 168 iterated collections and 168
  empty-population assertions; the audit accepts it.
- Exact Go 1.26.7 was recreated beneath scratch. Binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  downloaded archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- The official golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`;
  its Darwin arm64 binary SHA-256 is
  `691b9100ce968ff0009b6b7757ef6a585e31ae9ab11dfe0340ebb6e8e21fdc3d`.
  Exact reproduced GoReleaser 2.17.1 and apidiff binary SHA-256 values are
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Rebuilt govulncheck v1.7.0 SHA-256 is
  `0c536fb25db0d42aff2e5f4496d37a2a89b9f6c449466e2e493975db2652bb2e`.
- The historical golangci-lint receipt
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`
  and govulncheck receipt
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`
  are nonportable binary-build receipts and did not reproduce. Their pinned
  versions were independently obtained/rebuilt as above and passed every
  gate. Preserve the historical receipts, but do not require byte identity
  from a fresh build.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Preserve every other prior bounded decision.
- Prior exact counts, manifests, decisions, and corrections remain recorded
  in their answered archives. Preserve the go-colorful mutable telemetry,
  btree regression, cast whitespace-path, source-archive normalization, Check
  history table, Errgo preflight runner, Resty committed projection/signal
  fixture, Kingpin external-report, Units launcher-timing, Assert external-
  cache/TMPDIR, and Repr cache/mktemp/Docker/manifest corrections.
- Put recreated Go first in PATH, keep GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and inject no ambient GOFLAGS. The deliberately purged
  former two-entry recovery manifest remains historical SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected `github.com/alecthomas/kong`
`v0.2.1-0.20190708041108-0548c6b1afae` as the next single P7 dependency
group. A fresh post-Repr survey reports v1.16.1 at
2026-08-09T07:05:31Z, but canonical latest, release qualification, Go
declarations, complete closure, source identity, signatures, repository
state, loaded population, consumers, self-tests, and vulnerability effect are
unknown until proved from fresh primary evidence. Kong is not directly
required in current `go.mod`; the selected old version is transitive, go.sum
contains its go.mod checksum, repository source has no Kong import, and
`go mod why -m` says the main module does not need it. Treat those only as
starting observations.

Current measurements are 234 selected modules, 3,580 graph edges, 429 native
complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7.

Change only Kong's exact selected metadata and its explained minimal MVS
closure if the highest qualified candidate preserves the Go 1.18 floor and
passes every applicable gate. Do not combine Repr, Assert, Units, another
dependency group, or P8.
