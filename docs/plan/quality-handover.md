# Quality Upgrade Handover

Generated: 2026-09-05T11:05:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only Pelletier go-toml `/v2` commit
  `8cce879d285f08a3f8710a1ea35358fb057a9bfc`, exact parent
  `34f1314af51de0548fdb8ac3214dbd4481aabb24`, clean tree
  `8a88c3ff919b5275599d8461c1a176457997d2e7`. It changes only `go.mod` and
  `go.sum`.
- The documentation handoff commit created after this file must have exact
  parent `8cce879d285f08a3f8710a1ea35358fb057a9bfc`. Its predecessor is the
  answered go-toml v1 handoff `34f1314`, whose exact parent is `09cc40bc`,
  then `f974b581`, then accepted go-runewidth implementation `ca19dcd4`.
- Relative to accepted go-cmp commit `c314bcb`, dependency metadata changes
  only go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17, and
  go-toml/v2 v2.0.7 -> v2.2.2. The last move also selects testify v1.9.0 and
  objx v0.5.2 through minimal MVS closure and adds six checksum lines; the
  three accepted groups add ten checksum lines in total.
- The answered go-toml `/v2` archive and the sole NEXT `spf13/afero` archive
  must link reciprocally. Only the launcher's mutable header and prompt
  regions may change during handoff. Ordinary and ignored status must end
  empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/pelletier/go-toml/v2 v2.2.2`. The HCL, mousetrap, properties,
mapstructure, and go-toml v1 no-change/rejection decisions are closed.
Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color
v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and changes to any closed group remain
rejected or unnecessary for their recorded floor, loaded-behavior,
release-qualification, self-test, exact-latest, or closure decisions. Do not
revisit them. P8 remains queued.

## Accepted Pelletier Go-Toml V2 Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
v2.4.3 as canonical latest stable and v2.2.2 as the highest stable release
compatible with the retained Go 1.18 floor. The proxy lists 33 semantic
versions: 23 stable releases and ten alpha/beta prereleases. Proxy `@latest`,
exact `go list -m -versions -retracted`, and `go list -m -u` agree on v2.4.3;
there are no retractions.

All 15 stable releases v2.0.0-v2.2.2 declare Go 1.16. Every later stable
release v2.2.3-v2.4.3 declares Go 1.21.0, so successful execution under Go
1.26.7 does not make canonical latest compatible with the retained floor. The
ten prereleases declare Go 1.15 or Go 1.16 but do not displace a stable
candidate. V1 is a distinct module path; the proxy has no v3 versions, and a
non-semver GitHub prerelease named `latest` is not a Go module version.

Candidate v2.2.2 is lightweight tag commit
`a3d5a0bb530b5206c728eed9cb57323061922bcb`, committed and timestamped by
the proxy at 2024-04-29T10:02:54Z. Its stable GitHub Release was published
2024-05-01T15:13:03Z. Canonical latest v2.4.3 is lightweight tag commit
`071a36c2a57244f2e70369bfc69889fda2a1f60f` at
2026-07-05T02:25:11Z. Every one of the 33 tags is a lightweight commit ref,
so no tag object or tag signature exists; GitHub reports every embedded commit
signature present, valid, and verified.

The candidate checksum pair is
`h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=` /
`h1:1t835xjRzz80PqgE6HHgN2JOsmgYu/h4qDAS4n929Rs=`. All 97 files in the
candidate module ZIP match the exact tag commit byte-for-byte; its ZIP SHA-256 is
`8d724e35b485503810f866bca278d518e731713441e380634f6b33c27aefdf3e`
and file-manifest SHA-256 is
`5435b918796c4e60e035a46daeb137bbff8db131334d385c8a920f612bda6370`.
Selected v2.0.7 matches 93/93 files and latest v2.4.3 matches 113/113. The
repository is active, unarchived, and undisabled on default branch `v2`, has
no module `Deprecated` marker, and has seven untagged commits after v2.4.3;
none qualifies as a release.

The v2.0.7 -> v2.2.2 history is 46 commits across 35 files, +3,338/-723. The
seven-commit v2.2.2 -> v2.2.3 range includes the Go declaration jump and is
the first incompatible stable boundary.

Exact pinned `go get github.com/pelletier/go-toml/v2@v2.2.2` changes the
target selection plus testify v1.8.1 -> v1.9.0 and objx v0.5.0 -> v0.5.2.
V2.2.2 directly requires testify v1.9.0, which requires objx v0.5.2; objx's
testify v1.8.4 requirement adds only that version's go.mod checksum while MVS
retains v1.9.0. This is the explained minimal closure. Selected modules remain
234, complete package population remains byte-identical at 429, and graph
edges change 3,557 -> 3,565 through exactly six removed and 14 added edges.
Go.mod changes one target line; go.sum adds exactly six lines and removes
none. Historical v2.0.7 checksums remain. The unapplied `go mod tidy -diff`
projection changes 313 -> 321 lines and correctly exits 1 because a diff is
present; tidy was not used as implementation.

Exactly five `/v2` packages load:
`internal/characters`, `internal/danger`, `unstable`, `internal/tracker`, and
the module root. The path is `plybuild/cmd -> spf13/viper ->
viper/internal/encoding/toml -> pelletier/go-toml/v2`. Viper's real TOML codec
calls `/v2` Marshal and Unmarshal; old and candidate encode/decode tests pass
at count 10.

Candidate complete tests at count 1 and count 10 and race pass across all 16
packages. Standalone `go vet ./...` reports one malformed test-tag diagnostic
in `marshaler_test.go`. Normalized output is identical in selected v2.0.7 and
candidate v2.2.2, and blame locates it at pre-selection commit `0d20a845` from
2021. This is a historical upstream test diagnostic, not a regression:
default module tests/race and repository vet pass.

Repository mod verification, build, complete tests/race/vet, Windows build,
pinned golangci-lint 2.12.2, CLI surface, launcher check, and explicit
empty-HOME count-2 all pass. Help, API, and CLI reports remain byte-identical
at SHA-256 values
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`,
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`,
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

Fresh exact `make quality` passes all 21 ordered stages: preflight, eight
mutation meta-stages, 80/80 live mutants, host acceptance, fresh snapshot and
Docker acceptance, and authoritative Q0-Q2 audit. All 27 rows pass at L2;
manual evidence is valid with six receipts; ratchets have seven improved and
zero held, regressed, not-comparable, or dirty counts. A discarded first exact
run passed 20 stages then correctly failed closed because refreshed receipt
command strings still had their old payload digests. All six digests were
recalculated, a focused six-row audit independently passed, and the complete
quality run was restarted in a fresh root. No apparatus or source changed.

During final handoff, two Make launcher-suite attempts reproduced the
established nested signal-fixture timing diagnostic: the interrupted child did
not retain its partial raw log before inspection. The isolated 62-control
suite then passed, and a final Make entry-point run passed the same 62/62
controls. Prompt/archive identity and every graph control passed throughout.

The separate full audit exits expected 1, never 2, with only established
queued Q3.1, Q3.3, Q3.4, and Q3.7 non-passing. Govulncheck v1.7.0 preserves
exact and identical Darwin-symbol/Darwin-module/Windows-symbol populations
20/30/20. The fresh primary Go vulnerability module index, last modified
2026-09-02T19:32:20Z, has no go-toml v1, v2, or v3 record.

## Evidence And Tool Identity

- Selection root `/private/tmp/ply-p7-go-toml-v2-selection-final.34f1314.gn2ns2`
  fully verifies 17,093 entries at manifest SHA-256
  `ce1a47d85cd2fa0ff3bde19c1953c0dc73914017145abe2f4c5a27c4cf685ff1`;
  selection-summary SHA-256 is
  `a04e47128ea2acd2745e2d15010d39ae1090375652a00ff5d78803f7b4ffdfaf`.
- Source/manual review root
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/quality-review-final`
  fully verifies 8 entries at manifest SHA-256
  `328672c0410e2bfffa058fdf5ec389f4c8640bd73caccde22603b0ad57c9b070`;
  review-summary SHA-256 is
  `99b0520bec8d1796654f32916e95860bae358cb80de0313692e712dfc0ca0fde`.
- Exact-quality root
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/exact-quality-final`
  fully verifies 236,599 entries at manifest SHA-256
  `85f24a00baaf63d62bb6a3a07dc97bffcaa98b8fed1417b7d6f8425ba8d98d3d`;
  quality-summary SHA-256 is
  `a97c28ff45a8b85dfd63aaa58744afcc4b68ca63a640c00245b86062df3a3fd8`.
- Regression root
  `/private/tmp/ply-p7-go-toml-v2-selection.34f1314.s0l6bA/regression-gate-final`
  fully verifies 3,743 entries at manifest SHA-256
  `9601ce867fe389eb1efca59b2ae578052ba9b48670d1112cef3dc66298026ec7`;
  regression-summary SHA-256 is
  `4ea9f46a17f122f660ecca42b5601f59d4ccef4054d0a9de9ff7ab3e69e49e02`.
- Go-toml v1 and mapstructure decisions remain fully verified at
  37,719/`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`
  and 46,352/`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
  Accepted runewidth selection/review/exact-quality/regression evidence remains
  fully verified at 28,870/15/236,329/4,888 entries and its recorded hashes.
- Recovery remains fully verified at
  `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs`, two entries and manifest
  SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Preserve the go-colorful correction: 47,633 selection entries have exactly
  one recorded mutable telemetry mismatch and 47,632 valid entries; later
  review/exact-quality/regression manifests fully verify. Preserve the btree
  correction: its selection/review/exact-quality manifests fully verify and
  its 24,197-entry regression root has exactly 163 recorded mutable cache/HOME
  mismatches superseded by later evidence.
- Properties, mousetrap, and HCL evidence remains fully verified at
  17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
  and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
- Tool identities remain Go 1.26.7
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
  golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected indirect `github.com/spf13/afero v1.9.4` as
exactly one bounded P7 module group. Resolve canonical latest, every stable
release that could satisfy the retained Go 1.18 floor, source identity,
release qualification, requirements, closure, loaded package population and
behavior, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths, focused behavior,
repository quality, help/API/CLI identity, and exact vulnerability populations.
Implement only an exact floor-compatible release with explained minimal
closure and every gate passing; otherwise record rejection without changing
dependency metadata. Do not combine cast, jwalterweatherman, gotenv, x/*,
another dependency group, or P8.
