# Quality Upgrade Handover

Generated: 2026-09-08T00:58:41+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation remains Circbuf commit
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum`. Consul API, Go Metrics, and Go Radix made no dependency edit or
  implementation commit.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  Repr implementation `6f4d02eb9c86ec2df8a488a85ed973aababe1f38` and
  direct-child Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  Preserve this ancestry.
- The answered Go Radix archive and sole NEXT Perks archive link reciprocally.
  Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  still adds exactly 29 checksum lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Preserve launcher lifecycle and scratch-only
mutation, compatibility, snapshot, Docker, acceptance, and audit rules. Never
run `go mod download all` in a measured tree.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Go Radix v1.0.0. All earlier accepted,
rejected, and no-change decisions and evidence corrections remain final. P8
remains queued.

## Retained Armon Go Radix Group

Retain exact-path `github.com/armon/go-radix v1.0.0` without dependency
metadata. The exact proxy lists one semantic version, stable v1.0.0. Exact
`@latest` and `@v1` select it at 2018-08-24T02:57:28Z. Exact `@master` selects
unreleased pseudo-version `v1.0.1-0.20221118154546-54df44f2176c` at
2022-11-18T15:45:46Z. Master is neither a tag nor a GitHub Release and was not
selected merely because it is newer.

Go-import metadata identifies `https://github.com/armon/go-radix.git` as the
canonical exact-path source. It is public, enabled, unarchived, and non-fork,
with default branch `master`, one branch, one tag, 36 default-branch commits,
and zero GitHub Releases. Stable tag v1.0.0 is release-qualified, but it is a
lightweight Git tag rather than a GitHub Release object.

The tag identifies commit `1a2de0c21c94309923825da3df33a4381872c795`, tree
`8c6d01daaee6076244d5f41247608c75a8ad4224`, parent
`7fddfc383310abc091d79a27f116d30cf0424032`. The tag has no tag-object
signature. Its commit OpenPGP signature verifies locally with fingerprint
`7A01BBD67E7E8ADD50E00714744E147AA52F5B0A`, though GitHub currently reports
`unknown_key`. Master is merge commit
`54df44f2176c4a553657a4f0dbe6fdb108288be3`, tree
`87f38e748e5fc5c602ba3296793b25a782fbf02d`; GitHub reports its signature
valid and local verification identifies expired web-flow fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`.

The selected tag is an ancestor of master. Ten later commits change only
`.travis.yml`, `radix.go`, and `radix_test.go` with 86 insertions and 13
deletions. The history adds faster sorted insertion, documentation, Linux
power support, and WalkPrefix/delete-during-walk changes, but publishes no new
release. There are no proxy prereleases, alternate major paths, renamed
authoritative paths, retractions, or deprecation declarations.

Selected and master checksum pairs are respectively
`h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=` /
`h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=` and
`h1:651/eoCRnQ7YtSjAnSzRucrJz+3iGEFt+ysraELS81M=` with the same `go.mod`
checksum. Sumdb agrees. Proxy ZIP SHA-256 values are
`df93c816505baf12c3efe61328dc6f8fa42438f68f80b0b3725cae957d021c90`
and `f261a1141112f65564cec8f652ef6abf1d654228275e4fec6220b98337667c13`.
Each seven-file proxy archive matches its exact Git source. Normalized source
manifest SHA-256 values are
`f7801af436dc78ae44f1611f9160c6c37c3355ca3729c236c91f0fbd8834c518`
and `18854160a6271e82a14e65bcc15d142d83b640f78b087186b095176c162a1895`.

Neither form declares a Go version or requirements. Each complete minimal
module closure contains only Go Radix and one package. Missing declarations
were not compatibility evidence: writable proxy and exact-Git forms pass
verification/listing and complete count-1/count-10/race/vet under exact Go
1.26.7, while proxy forms also pass under exact Go 1.18.10. Direct execution
therefore proves that both one-module closures preserve the Go 1.18 floor.

## Consumers, Projection, And Quality Scope

Selected `github.com/mitchellh/cli v1.1.0` is the actual graph consumer. It
uses `New`, `Insert`, `Get`, `Walk`, `WalkPrefix`, `LongestPrefix`, `Tree`, and
`WalkFn`; selected and historical Serf versions only declare Radix indirectly
through CLI. Focused nested-command, help, autocomplete, and subcommand tests
pass count-1, count-10, race, and vet against v1.0.0 and master. Radix loads in
zero Ply packages, Ply has no import, and `go mod why -m` says the main module
does not need it.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
complete packages, 41 loaded modules, zero loaded Radix packages, 1,045
`go.sum` lines, and 361 tidy-diff lines. Exact selected-version get changes no
selection and projects only a redundant indirect requirement, one main edge,
and one full checksum: 234/3,582/429/41/0, 1,046 sum lines, and 367 tidy lines.
Tidy removes that metadata.

Exact master get changes only the selected Radix version but projects an
explicit main edge and the candidate checksum pair: 234/3,582/429/41/0,
1,047 sum lines, and 369 tidy lines. Tidy removes the pin and pair and restores
inherited v1.0.0. Because release qualification rejects master and exact
stable get is a selection no-op, neither projection was applied.

Baseline and master projections pass repository module verification, build,
complete count-1/count-10/race tests, vet, Windows-amd64 build, and pinned
golangci-lint 2.12.2. Root, status, upgrade, and build help are byte-identical.
API and CLI reports are byte-identical at
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Corrected baseline preflight passes all 62 launcher controls, every
distribution/lint/install/toolchain and script meta-contract, and 15 audit
controls. Empty-HOME count-2 passes.

Fresh govulncheck v1.7.0 uses primary data updated
2026-09-02T19:12:04Z. Its 1,392-record module index has no Radix record or
trace. Baseline/master projections are byte-identical: 20 IDs/22 traces for
Darwin and Windows reachable symbols and 30 Darwin module IDs.

Canonical qualification produced no changed selection. The changed-selection-
only host/snapshot/Docker acceptance, exact `make quality`, focused/Q0-Q2/full
audits, and dependency commit were inapplicable and are not claimed. Go Radix
evidence contains 390 verified entries; evidence-manifest SHA-256 is
`6be6b857e77c7097b402ea5f4fe0ce849e15382b7167a004448cd85d82926eaf`;
decision-summary SHA-256 is
`493e832f2d0e6456bb64462104e1bd5b80f4b6eb32d401b9b098acdd0dd95e6a`.

## Tools And Corrections

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
  Go 1.26.7 remained first in PATH with GOENV/GOWORK off, local toolchain, and
  no ambient GOFLAGS.
- Preserve portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Rebuilt tool binaries are nonportable; version/build metadata and functional
  use are the proof.
- Superseded runner-only corrections stayed beneath session scratch:
  source-prefix normalization, exact historical compatibility-cache warming,
  reachable vulnerability trace-depth filtering, and scratch-local BSD
  `mktemp` plus cleared recursive make overrides. None changed repository
  files or concealed a failure.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`, accepted Circbuf, Consul API, Go
  Metrics v0.4.0, Go Radix v1.0.0,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections from prior dependency groups. The
  deliberately purged former two-entry recovery manifest remains historical
  SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path `github.com/beorn7/perks v1.0.1`
at 2019-07-31T12:00:54Z as the next single P7 group. Project MVS selects it
without an explicit `go.mod` requirement. The exact proxy lists v1.0.0 and
v1.0.1; exact `@latest`, `@v1`, and `@master` all resolve selected v1.0.1.
Its proxy module declares Go 1.11, but the complete closure floor remains
unproved.

The public exact-path repository is enabled and unarchived but is a fork. It
identifies `bmizerany/perks` as parent and source, reports default `master`,
five branches, two tags, and zero GitHub Releases. Exact-path master and tag
v1.0.1 identify `37c8de3658fcb183f997c4e13e8337516ab753e6`; tag v1.0.0 identifies
`4b2b341e8d7715fae06375aa633dbb6e91b3fb46`. Independently resolve canonical
published identity across that fork relationship; do not assume repository
ancestry changes the exact module path or release qualification.

Current accepted measurements remain 234 modules, 3,581 graph edges, 429
complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, a 361-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Change only Perks's exact selected
version and explained minimal closure if a qualified changed selection
preserves Go 1.18 and passes every applicable contract. If v1.0.1 is already
the highest qualified selection, do not manufacture redundant metadata. Do
not combine another dependency group or P8.
