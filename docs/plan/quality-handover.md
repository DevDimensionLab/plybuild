# Quality Upgrade Handover

Generated: 2026-09-07T22:37:40+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation remains Circbuf commit
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum`. The Consul API evaluation made no dependency edit or implementation
  commit.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  Repr implementation `6f4d02eb9c86ec2df8a488a85ed973aababe1f38` and
  direct-child Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  All remain ancestors.
- The answered Consul API archive and sole NEXT Go Metrics archive link
  reciprocally. Relative to accepted go-cmp commit `c314bcb`, accepted
  dependency metadata still adds exactly 29 checksum lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Preserve the launcher lifecycle and scratch-
only mutation, compatibility, snapshot, Docker, acceptance, and audit rules.
Never run `go mod download all` in a measured tree.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained canonical-latest Consul API pseudo-
version `v0.0.0-20180202201655-eb2c6b5be1b6`. All earlier accepted, rejected,
and no-change decisions and evidence corrections remain final. P8 remains
queued.

## Retained Armon Consul API Group

Retain exact-path `github.com/armon/consul-api` at selected pseudo-version
`v0.0.0-20180202201655-eb2c6b5be1b6` without changing dependency metadata.
The exact stable proxy list is empty; proxy and exact Go `@latest` and
`@master` resolve the already-selected version, while exact `@v0` has no
matching version. It is canonical latest and the highest declaration-
compatible exact-path candidate, but it is an unreleased pseudo-version, not a
stable release.

The authoritative repository is public, enabled, operationally unarchived,
and non-fork, with default branch `master`, one branch, 46 commits, zero tags,
and zero GitHub Releases. Selected/master-head merge commit
`eb2c6b5be1b66bab83016e0b05f01b8d5496ffbd`, tree
`aeb2299aaf107d0823ce91f057798119b821e81b`, is dated
2018-02-02T20:16:55Z. Its embedded GitHub web-flow signature is
cryptographically valid for fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the key is currently expired.
There is no tag or tag signature.

The repository is operationally unarchived but its README explicitly
deprecates this source in favor of distinct module/import path
`github.com/hashicorp/consul/api`. That successor is not an alternate version
on the exact path; its current latest v1.34.4 declares Go 1.26 and is an out-of-
scope migration.

All 21 regular proxy files match the exact commit archive and Git tree.
Normalized proxy/upstream source-manifest SHA-256 is
`358c1cfd7c63c682fc02e069ef0ab20f494bc44d9789e62ca39e62fc4026cbee`.
Proxy ZIP and GitHub archive SHA-256 values are
`091b79667f16ae245785956c490fe05ee26970a89f8ecdbe858ae3510d725088`
and `ba5e02e82cb70d6d7cb7970a8f81535df3266cfeae2b2d3b9e38de7ca2d76b4a`.
The selected checksum pair is
`h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=` /
`h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=`; sum.golang.org agrees.

The synthetic module declares no Go version and no requirements. Its complete
declared minimal closure is Consul API alone and preserves Go 1.18. Exact Go
1.26.7 proxy and exact-commit forms independently verify/list but identically
fail complete count-1/count-10/race tests with 33/330/33 refused connections
to the required external Consul agent at 127.0.0.1:8500. Vet fails on two
`testing.T.Fatalf` calls from non-test goroutines. Exact Go 1.18.10 compiles
the package and reproduces the test and vet failures. These are dependency
stop-rule failures.

Project MVS receives the selected version from
`github.com/devdimensionlab/mvn-pom-mutator@v0.2.3`, whose source does not
import the module. Consul API loads in zero Ply packages, Ply has no source
import, and `go mod why -m` says the main module does not need it. Historical
Crypt consumer source and an external Go-1.18 HTTP fixture cover
`DefaultConfig`, `Config.Address`, `NewClient`, `Client.KV`, `KV.Get`,
`KV.List`, `KV.Put`, `KVPair` key/value fields, `QueryOptions.WaitIndex`, and
`QueryMeta.LastIndex`. The fixture passes count-10, race, and vet; it is
consumer evidence, not Ply execution.

Exact selected-version get changes no selection. It projects 234 -> 234
selected modules, 3,581 -> 3,582 graph edges, 429 -> 429 complete packages,
41 -> 41 loaded modules, zero -> zero loaded Consul packages, 1,045 -> 1,046
sum lines, and 361 -> 363 tidy-diff lines. The only changes are a redundant
indirect requirement, main graph edge, and full checksum; tidy removes all
three while leaving inherited debt and the existing go.mod checksum unchanged.
Therefore no metadata or dependency implementation commit was created.

## Consul API Quality And Evidence

- Govulncheck v1.7.0 used primary data updated
  2026-09-02T19:12:04Z. Its 1,392-record module index has no Consul API record
  and no trace contains the module. Original/exact-get projections are
  byte-identical at 20 IDs/22 traces for Darwin reachable symbols, 30 Darwin
  module IDs, and 20 IDs/22 traces for Windows reachable symbols. Reachable
  and module normalized SHA-256 values are
  `5bebff017082945899b65922b3262a8b91affcda7abfe2fade54dee439749dfe`
  and `f1cc7393d1a1d82da88379fceb830d926f8c1f0fe553e4af339e76ac6f7a1bc3`.
- Downstream repository, snapshot/Docker, acceptance, and Q0-Q2 gates were not
  run or claimed: no selected version changed, exact get manufactured only
  tidy-removable metadata, and the dependency module test/vet stop rule failed.
  The accepted Circbuf quality baseline remains unchanged.
- Consul API decision evidence contains 315 verified entries. Evidence-
  manifest SHA-256 is
  `1dbe7063e72dade7bc31e9c8967da78a60ef97859b68562ffa1a07b75a0b3b0b`;
  decision-summary SHA-256 is
  `0f5413813a179949c2dbce2a029becfceb9f650f7a603ac8b19755c3f4d48733`.

## Tools And Corrections

- Exact Go 1.26.7 binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 archive/binary SHA-256 values are
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`
  and `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
  Go 1.26.7 remained first in PATH with GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and no ambient GOFLAGS.
- Preserve portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Govulncheck v1.7.0 and rebuilt-tool hashes are nonportable receipts; versions
  and functionality were proved.
- Supersede the first Go-1.18 CGO setup, which attempted an Xcode resolver cache
  outside scratch, with final direct SDK/clang and scratch TMPDIR/GOTMPDIR
  runs. Supersede invalid patterned govulncheck module-mode calls with correct
  no-pattern module scans. Preserve concatenated-JSON decoding and trace-depth
  filtering. No correction mutated the repository or hid a candidate failure.
- Preserve all previously recorded compatibility-cache, source/signature,
  checksum-delta, vulnerability-order, Docker-path, scratch-local BSD-mktemp,
  preflight, locale, and exact-Go-1.18 final-rerun corrections.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`, accepted Circbuf candidate, retained
  Consul API pseudo-version, `gopkg.in/alecthomas/kingpin.v2 v2.2.6`,
  `gopkg.in/resty.v1 v1.12.0`, `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections from prior dependency groups. The
  deliberately purged former two-entry recovery manifest remains historical
  SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path
`github.com/armon/go-metrics v0.4.0` at 2022-05-25T15:01:32Z as the next single
P7 group. Project MVS selects it without an explicit `go.mod` requirement. A
minimal survey finds 22 exact stable proxy versions; exact `@latest` and
`@v0` resolve v0.6.1 at 2026-07-29T13:06:13Z with `go 1.25.0`, above the
retained floor. Exact `@master` resolves unreleased pseudo-version
`v0.6.2-0.20260907064447-465585286d74` at 2026-09-07T06:44:47Z, also with Go
1.25.

The GitHub request for `armon/go-metrics` resolves repository metadata for
`hashicorp/go-metrics`, reporting public, enabled, unarchived, non-fork status,
default branch `master`, 20 branches, 21 tags, and 15 Releases. Treat redirect
and exact module source identity, release/tag/signature history, retractions,
all stable version declarations, complete closure, module tests, loaded
packages, real consumers, actual symbols, and vulnerability effect as unknown
until independently proved. Resolve the highest qualified complete closure that
preserves Go 1.18; do not select the Go-1.25 latest or guess a candidate.

Current measurements are 234 selected modules, 3,581 graph edges, 429 native
complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, a 361-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Change only Go Metrics' exact
selected version and explained minimal MVS closure if a qualified changed
selection preserves the floor and passes every applicable gate. If v0.4.0 is
the highest qualified version, do not manufacture redundant metadata. Do not
combine Consul API, another module group, or P8.
