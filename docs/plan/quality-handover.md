# Quality Upgrade Handover

Generated: 2026-09-08T02:34:21+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation remains Circbuf commit
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum`. Consul API, Go Metrics, Go Radix, and Perks made no dependency
  edit or implementation commit.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  Repr implementation `6f4d02eb9c86ec2df8a488a85ed973aababe1f38` and
  direct-child Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  Preserve this ancestry.
- The answered Perks archive and sole NEXT Speakeasy archive link
  reciprocally. Relative to accepted go-cmp commit `c314bcb`, accepted
  dependency metadata still adds exactly 29 checksum lines.
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
dependency decisions through retained Perks v1.0.1. All earlier accepted,
rejected, and no-change decisions and evidence corrections remain final. P8
remains queued.

## Retained Beorn7 Perks Group

Retain exact-path `github.com/beorn7/perks v1.0.1` without dependency
metadata. The exact proxy contains only stable v1.0.0 and v1.0.1 at
2019-04-14T22:11:40Z and 2019-07-31T12:00:54Z. Exact `@latest`, `@v1`,
`@master`, default-branch `master`, and tag v1.0.1 all resolve selected commit
`37c8de3658fcb183f997c4e13e8337516ab753e6`. There are no prereleases,
retractions, deprecation declarations, `/v2` module, or later exact-path
default-branch commit.

Go-import metadata maps `github.com/beorn7/perks` to the enabled, unarchived
fork `https://github.com/beorn7/perks.git`. GitHub parent/source
`bmizerany/perks` records repository ancestry, not module identity. Parent
go-import declares distinct path `github.com/bmizerany/perks`; its untagged
master `03f9df79da1edead2cdf5f8b4cf4d4f831d6e2d1` is divergent, lacks
`go.mod`, and resolves only parent pseudo-version
`v0.0.0-20230307044200-03f9df79da1e`.

The exact fork has five branches, two signed annotated tags, and zero GitHub
Releases. Its old non-default heads are v0-era refs, not newer v1 candidates.
The parent has three branches, zero tags, and zero Releases. Stable tags are
release-qualified here, but a Git tag is not a GitHub Release object.

Tag v1.0.0 object `4ded152d4a3e...` targets unsigned commit
`4b2b341e8d7715fae06375aa633dbb6e91b3fb46`, tree
`de22b44fded2eabaabd0ac413a3b04b44c8b3490`. Its signature verifies locally
with fingerprint `5C69F212D616C4340FA8DD8504ABA6153ADA0C25`; GitHub reports
`unknown_key`. Tag v1.0.1 object `c49ff274687...` targets unsigned commit
`37c8de3658fcb183f997c4e13e8337516ab753e6`, tree
`77aeb432cabb3c9b94297378a2fe7697c29a5a88`. It verifies locally with
fingerprint `A100A34F34DEC17EE5EEF14C851C3DA17D748D03`; GitHub reports valid.
V1.0.1 changes only `go.mod`, lowering Go 1.12 to Go 1.11.

The checksum pairs for v1.0.0 and v1.0.1 are respectively
`h1:HWo1m869IqiPhD389kmkxeTalrjNbbJTC8LXupb+sl0=` /
`h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=` and
`h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=` /
`h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=`; sumdb agrees. Proxy ZIP
SHA-256 values are
`a7ec6164e31ea8e10c601abb9793753ec43cb218283b226800c134fb23cea409`
and `25bd9e2d94aca770e6dbc1f53725f84f6af4432f631d35dd2c46f96ef0512f1a`.
All 15 regular files in each archive match the exact corresponding tag.

Both modules have no requirements. Their complete closures are one module and
three packages, with explicit Go 1.12/1.11 declarations; both preserve the
retained Go 1.18 floor. Final proxy and exact-tag matrices pass verify/list,
count-1, two count-10 runs, race, vet, and source immutability under exact Go
1.26.7 and Go 1.18.10. V1.0.1 additionally passes focused topk count-100 in
both forms/toolchains.

An earlier v1.0.0 proxy Go-1.26.7 count-10 run failed once in
`topk.TestTopK` (`want "9", got "b"`). Four final matrices and count-100
stress did not reproduce it. The test sorts only by count, making equal-count/
map-order nondeterminism plausible. Preserve the observation; do not relabel
it as a v1.0.1 source change because v1.0.1 changes only `go.mod`.

## Consumer, Projection, And Quality Scope

Selected `github.com/prometheus/client_golang v1.4.0` is the actual graph
consumer. Its Summary imports `perks/quantile` and uses `Stream`,
`NewTargeted`, `Insert`, `Count`, `Query`, and `Reset`. Focused Summary tests
pass count-1/count-10/race under exact Go 1.18.10. The separate standalone
consumer closure is 41 modules and 167 test packages with no declared Go
version above 1.12. Six legacy int-to-string diagnostics in consumer test
formatting block default test vet; execution with implicit test vet isolated
passes. Perks and project vet pass normally.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
complete-test packages, 41 loaded modules, zero loaded Perks packages, 1,045
`go.sum` lines, and a 361-line tidy projection. Exact selected-version get
changes no selection and projects only a redundant indirect requirement, one
main edge, and the full v1.0.1 checksum: 234/3,582/429/41/0, 1,046 sum lines,
and 370 tidy lines. Tidy removes those two metadata lines and restores
inherited selection. No projection was applied.

The inherited requests are TSDB v0.7.1 -> old Perks pseudo-version,
client_golang v1.4.0 -> v1.0.1, client_golang v1.0.0 -> v1.0.0, and
Prometheus common v0.4.1 -> old pseudo-version. Perks has no outgoing graph
edge, loads in no Ply package, and `go mod why -m` says the main module does
not need it.

Baseline/projected trees pass verify, build, complete count-1/count-10/race,
vet, Windows-amd64 build, pinned golangci-lint 2.12.2, public help, API/CLI,
and empty-HOME count-2. API and CLI reports remain byte-identical at
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Corrected preflight passes 62 launcher controls, every Make/distribution/
acceptance meta-contract, 80/80 mutation controls, and 15 audit controls.

Fresh govulncheck v1.7.0 uses primary data updated
2026-09-02T19:12:04Z. The 1,392-record module index contains neither exact nor
parent Perks. Baseline/projected results are identical: 20 IDs/22 traces for
Darwin and Windows reachable symbols and 30 Darwin module IDs. Normalized
Darwin-symbol, Windows-symbol, and module hashes are
`fcf9d449a455a00562126bca8b863b882ede61d7c087a967ac9694c30978c859`,
`b7d80a548ee7476342d7a0bd2b0efd87524a5efa0c0da960e802c4e186fb3e8f`,
and `390f7bf685e19fc5d10602da7e35394c97f8fd460fe40f92c4e6726c00627d29`.

Canonical qualification produced no changed selection. Changed-selection-
only host/snapshot/Docker acceptance, exact `make quality`, focused/Q0-Q2/
full audits, and dependency commit were inapplicable and are not claimed.
Perks evidence contains 477 verified entries; evidence-manifest SHA-256 is
`aa77d2ab9cf676ecfb7ef544c1c5db76ecda86f580ff8ff8f03b8c78038677f0`;
decision-summary SHA-256 is
`6db4ca7260bde6b4affa48c62adb381bda20687f39f013c4b1cdc25397418754`.

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
- Scratch-only corrections are retained explicitly: initial environment and
  source-prefix fixes, removed empty slash-branch query, vulnerability-order
  normalization, exact historical compatibility graph warming, golangci
  cache relocation, BSD `mktemp` adaptation, cleared recursive make
  overrides, and a 200 ms signal-fixture scheduling shim after two raw-log
  readiness races. None changed a repository file or concealed a dependency
  failure.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`, accepted Circbuf, Consul API, Go
  Metrics v0.4.0, Go Radix v1.0.0, Perks v1.0.1,
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

Independently evaluate selected exact-path
`github.com/bgentry/speakeasy v0.1.0` at 2017-04-17T20:07:03Z as the next
single P7 group. Project MVS selects it through Mitchellh CLI without an
explicit `go.mod` requirement. The exact proxy lists stable v0.1.0 and
v0.2.0; exact `@latest`, `@v0`, and `@master` all resolve v0.2.0 at
2022-09-10T01:20:23Z, commit
`760eaf8b681647364e7a400b856e0921248728a5`. Both proxy module files expose
only exact path `github.com/bgentry/speakeasy`; do not infer floor
compatibility from their missing Go directives and requirements.

The public repository is enabled, unarchived, undisabled, and non-fork, with
default `master`, one branch, two tags, and one non-draft, non-prerelease
GitHub Release. Master and v0.2.0 identify
`760eaf8b681647364e7a400b856e0921248728a5`; v0.1.0 identifies
`4aabc24848ce5fd31929f7d1e4ea74d3709c14cd`. The v0.2.0 Release was
published 2024-06-27T20:45:36Z for the 2022 tag. Independently establish
source/signature/release history, complete closure floor, tests, real consumer
terminal behavior, project projection, quality, and vulnerability effect.

Current accepted measurements remain 234 modules, 3,581 graph edges, 429
complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, a 361-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Change only Speakeasy's exact
selected version and explained minimal closure if v0.2.0 is qualified,
preserves Go 1.18, and passes every applicable contract. Do not combine
another dependency group or P8.
