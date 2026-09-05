# Quality Upgrade Handover

Generated: 2026-09-05T11:59:15+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only Pelletier go-toml `/v2` commit
  `8cce879d285f08a3f8710a1ea35358fb057a9bfc`, exact parent
  `34f1314af51de0548fdb8ac3214dbd4481aabb24`, clean tree
  `8a88c3ff919b5275599d8461c1a176457997d2e7`. It changes only `go.mod` and
  `go.sum`.
- Current Afero documentation handoff is
  `5a4653999ac6a569a95529034a43d457a41dab1b`, exact parent `8cce879d`.
  The documentation handoff created after this file must have exact parent
  `5a4653999ac6a569a95529034a43d457a41dab1b`; no dependency implementation
  commit exists for the rejected Afero group.
- Relative to accepted go-cmp commit `c314bcb`, dependency metadata changes
  only go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17, and
  go-toml/v2 v2.0.7 -> v2.2.2. The last accepted move also selects testify
  v1.9.0 and objx v0.5.2 through minimal MVS closure. Those three accepted
  groups add exactly ten checksum lines in total.
- The answered Afero archive and sole NEXT `spf13/cast` archive must link
  reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/pelletier/go-toml/v2 v2.2.2`, plus rejected Afero v1.10.0.
The HCL, mousetrap, properties, mapstructure, go-toml v1, and Afero decisions
are closed. Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3,
fatih/color v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and changes to any closed group remain
rejected or unnecessary for their recorded floor, loaded-behavior,
release-qualification, self-test, exact-latest, or closure decisions. Do not
revisit them. P8 remains queued.

## Rejected Spf13 Afero Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establishes
v1.15.0 as canonical latest stable and v1.10.0 as the highest stable release
compatible with the retained Go 1.18 floor. The proxy lists 37 main-module
versions, all stable, with zero prereleases, retractions, or formal deprecation
markers. Releases v1.0.0-v1.2.2 have no Go directive; v1.3.0-v1.8.2 declare
Go 1.13; v1.9.0-v1.10.0 declare Go 1.16; v1.11.0 declares Go 1.19;
v1.12.0 declares Go 1.21; and v1.13.0-v1.15.0 declare Go 1.23.0.

Canonical latest v1.15.0 is annotated tag object
`18d690e34969d06817fa791ccf69194ebd4a5e8d`, peeled commit
`399bb34ad9fd8a252ad1d8bfaef96279b66dc774` at
2025-09-08T16:25:29Z. GitHub verifies its tag and commit signatures.
Candidate v1.10.0 is lightweight tag commit
`ee6eef77ef4a6c73b07a4bc070a9a2f076fd121e` at
2023-09-22T14:18:35Z; the tag has no object or signature and GitHub verifies
the commit signature. Selected v1.9.4 is unsigned lightweight commit
`cf95922e71986c0116204b6eeb3b345a01ffd842`. Local signature verification
is unavailable because `gpg` is absent.

The active, unarchived, undisabled upstream is not a fork. Four additional
`gcsfs/` and `sftpfs/` v1.14/v1.15 tags belong to nested modules; there is no
main-module v2 path. Master has 77 commits after v1.15.0, with current
unreleased commit `768f1fb0e5535b77d90e44c531aacd652aabd96a`; none is a
released main-module version. The selected-to-candidate range is seven commits
across nine files, +159/-26, covering SFTP ReadAt, z/OS flags, and MemMapFs
descendant-rename behavior.

Candidate checksum pair is
`h1:EaGW2JJh15aKOejeuJ+wpFSHnbd7GE6Wvp3TsNhb6LY=` /
`h1:UBogFpq8E9Hx+xc5CNTTEpTnuHVmXDwZcZcE1eb/UhQ=`. All 66 candidate proxy
files match the exact tag commit byte-for-byte; ZIP SHA-256 is
`96f549d577731c5f75d5a1b2bf8b937ff4300fab318197776af4f335c0bc0ff6`
and file-manifest SHA-256 is
`93f868aa14b79c81fc82369f4d48c4989aced10c76d12bb3df37539df0e7c28a`.
Selected v1.9.4 and incompatible-boundary v1.11.0 also match 66/66 files.
Latest v1.15.0 matches all 59 proxy-eligible files; 22 Git files belong to
its two nested modules.

Projected exact `go get github.com/spf13/afero@v1.10.0` changes Afero plus
`golang.org/x/crypto` from
`v0.0.0-20220525230936-793ad666bf5e` to
`v0.0.0-20220722155217-630584e8d5aa`. Afero directly raises that requirement;
the candidate x/crypto declares Go 1.17. Its x/net requirement and Afero's
x/text requirement remain below already selected versions, so neither
selection moves. This is the minimal MVS closure. The projection retains 234
modules and the byte-identical 429-package complete population; graph edges
change 3,565 -> 3,572. Go.mod changes only Afero, while go.sum adds the
Afero checksum pair plus x/crypto and x/net historical go.mod checksums.
The unapplied tidy projection changes 321 -> 323 lines; tidy was not used as
implementation.

Three Afero packages load: the root, `internal/common`, and `mem`. The real
path is `plybuild/cmd -> spf13/viper -> spf13/afero`; Viper uses its
filesystem interface, OS/IO adapters, file reads, existence checks, and
write/sync operations. Focused Viper filesystem/configuration behavior and Ply
`./cmd` tests pass at count 10 in both selected and projected states, and the
candidate root package passes at count 10.

Candidate complete tests, race, and vet pass at count 1 across all seven
packages. The required complete `go test ./... -count=10` does not: tarfs
`TestRead` emits exactly 27 instances of `got 1 read bytes, expected 8`.
Selected v1.9.4 reproduces the identical diagnostics, the test source is
byte-identical at SHA-256
`ef093f4a40042918f70cbf1e6fa0c93a4f3d575140245880921b69a77c1a398d`,
and blame dates it to 2020 commit `a4ea980f`. Ten fresh count-1 tarfs
processes pass; the repeated-run defect comes from TestMain reusing one
stateful tar reader. It is historical rather than a candidate regression, but
the mission explicitly requires stopping on any module self-test failure.
V1.10.0 is therefore rejected and v1.9.4 retained without metadata edits or
an implementation commit. Repository-wide candidate gates after that stop
were intentionally not run.

Govulncheck v1.7.0 preserves exact and identical old/candidate
Darwin-symbol/Darwin-module/Windows-symbol ID populations 20/30/20. The
primary vulnerability module index, last modified 2026-09-02T19:32:21Z, has
no Afero entry.

## Evidence And Tool Identity

- Afero decision root `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`
  fully verifies 79,866 entries at manifest SHA-256
  `f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`;
  decision-summary SHA-256 is
  `c3301a1c2b13f36956551ed1eb09b1a83bf2ac51cf236c6614f46eaf70ea65de`.
- Accepted go-toml/v2 selection/review/exact-quality/regression roots remain
  fully verified at 17,093/`ce1a47d85cd2fa0ff3bde19c1953c0dc73914017145abe2f4c5a27c4cf685ff1`,
  8/`328672c0410e2bfffa058fdf5ec389f4c8640bd73caccde22603b0ad57c9b070`,
  236,599/`85f24a00baaf63d62bb6a3a07dc97bffcaa98b8fed1417b7d6f8425ba8d98d3d`,
  and 3,743/`9601ce867fe389eb1efca59b2ae578052ba9b48670d1112cef3dc66298026ec7`.
- Go-toml v1 and mapstructure decisions remain fully verified at
  37,719/`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`
  and 46,352/`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
  Accepted runewidth selection/review/exact-quality/regression evidence remains
  fully verified at 28,870/15/236,329/4,888 entries and its recorded hashes.
- Properties, mousetrap, and HCL evidence remains fully verified at
  17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
  and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
- Preserve the go-colorful correction: its 47,633-entry selection manifest
  has exactly one recorded mutable telemetry mismatch; 47,632 entries and its
  later review/exact-quality/regression manifests verify. Preserve the btree
  correction: selection/review/exact-quality verify and its 24,197-entry
  regression root has exactly 163 recorded mutable cache/HOME mismatches
  superseded by later evidence.
- Recovery `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` remains fully
  verified at two entries and manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
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

Independently evaluate selected indirect `github.com/spf13/cast v1.5.0` as
exactly one bounded P7 module group. Resolve canonical latest, every stable
release that could satisfy the retained Go 1.18 floor, source identity,
release qualification, requirements, closure, loaded package population and
behavior, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths, focused conversion/
configuration behavior, repository quality, help/API/CLI identity, and exact
vulnerability populations. Implement only an exact floor-compatible release
with explained minimal closure and every gate passing; otherwise record
rejection without changing dependency metadata. Do not combine
jwalterweatherman, gotenv, afero, x/*, another dependency group, or P8.
