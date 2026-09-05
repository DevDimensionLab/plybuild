# Quality Upgrade Handover

Generated: 2026-09-05T15:41:29+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only cast commit
  `cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
  `17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
  `485690c010cef3b02385d8f2471d3ef53611abe9`. It changes only `go.mod` and
  `go.sum`.
- The current jwalterweatherman documentation session is commit
  `2dd8fdfc61d54d9327effad90242f1f05060fa42`, exact parent `cf4fd493`, tree
  `df410b8bc43402d93793830e912107d400a64f81`. This handoff must be its direct
  child. The next dependency implementation, if any, must use the resulting
  documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Minimal closures
  select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty v0.3.1,
  and rogpeppe/go-internal v1.9.0. Those four accepted groups add exactly 15
  checksum lines.
- The answered jwalterweatherman archive and sole NEXT `subosito/gotenv`
  archive link reciprocally. Only the launcher's mutable header and prompt
  regions may change during handoff. Ordinary and ignored status must end
  empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1, the
rejected afero v1.10.0 candidate, and retained jwalterweatherman v1.1.0.
Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color
v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, and jwalterweatherman changes remain rejected
or unnecessary for their recorded floor, loaded-behavior,
release-qualification, self-test, exact-latest, or closure decisions. Do not
revisit them. P8 remains queued.

## Retained Spf13 Jwalterweatherman Group

Fresh Go proxy, checksum-database, GitHub, upstream Git, and primary Go
vulnerability evidence establishes selected v1.1.0 as both canonical latest
stable and highest stable release compatible with the retained Go 1.18 floor.
The proxy lists exactly v1.0.0 and v1.1.0, both without a Go directive. There
are no prereleases, retractions, formal deprecation marker, fork, or later
major module path. The upstream repository is active, unarchived, and
undisabled.

V1.1.0 is lightweight tag commit
`94f6ae3ed3bceceafa716478c5fbf8d29ca601a1` at
2018-10-28T14:53:47Z. GitHub verifies the commit signature; the lightweight
tag has no separate tag object or signature. The sole GitHub Release is stable
and was published at 2019-03-02T09:22:47Z. V1.0.0 is unsigned annotated tag
object `251b5ffc375703044127b2929d4c820504c3ab7a`, peeled to GitHub-verified
commit `4a4406e478ca629068e7768fc33f3f044173c0a6`. Local signature verification
remains unavailable because `gpg` is absent.

Upstream master `91990e8269243f4f5d234f9c24e0a73c4a0df47f` is three commits after
v1.1.0. Its proxy pseudo-version
`v1.1.1-0.20230519082352-91990e826924` declares Go 1.20, but is unreleased and
excluded from stable selection. The v1.0.0 -> v1.1.0 range is two commits,
seven files, +141/-73; it adds LogListener, Counter, SetStdoutOutput, and test
requirements.

Candidate checksum pair is
`h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=` /
`h1:aNWZUN0dPAAO/Ljvb5BEdw96iTZ0EXowPYD95IqWIGo=`. All nine proxy files
match the exact tag commit byte-for-byte; ZIP SHA-256 is
`43cc5f056caf66dc8225dca36637bfc18509521b103a69ca76fbc2b6519194a3`
and file-manifest SHA-256 is
`dd5824d06389cec49ebdfed59947ae4d96bf1e9ef0ced0bab30d7088820451ee`.
The checksum pair also verifies against sum.golang.org.

Exact `go get github.com/spf13/jwalterweatherman@v1.1.0` changes zero
go.mod/go.sum bytes. Both projections remain byte-identical at 234 modules,
3,564 graph edges, 429 complete packages, two jwalterweatherman checksum lines,
and a 332-line unapplied tidy projection. V1.1.0 requires go-spew v1.1.1,
go-difflib v1.0.0, and testify v1.2.2. The first two already remain selected;
the testify edge loses to selected v1.9.0. Its minimal selection, edge,
checksum, and metadata closure is empty.

Exactly one jwalterweatherman package loads through
`github.com/devdimensionlab/plybuild/cmd -> github.com/spf13/viper ->
github.com/spf13/jwalterweatherman`. Viper root is the sole real consumer and
maps five methods to TRACE, DEBUG, INFO, WARN, and ERROR. Ten fresh runs of
focused Viper read/merge/override/alias behavior and Ply `./cmd` count-10 pass
in both states. Candidate module verification, complete tests at count 1 and
count 10, race, and vet pass. Because the historical tag omitted go.sum, its
canonical self-test used a writable external source, downloaded only its three
specific requirements, restored the exact proxy go.mod, and retained the
generated go.sum outside the worktree.

Repository build, complete tests/race/vet, pinned golangci-lint 2.12.2, API
compatibility, CLI compatibility and surface, complete preflight, and
empty-start-HOME count-2 pass. Root/status/upgrade/build help is byte-identical;
API report SHA-256 remains
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and CLI report SHA-256 remains
`955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
The complete preflight retry passed after creating its external TMPDIR; the
first attempt was setup failure before API comparison, not a product failure.

Govulncheck v1.7.0 preserves exact and identical old/candidate
Darwin-symbol/Darwin-module/Windows-symbol populations 20/30/20. The fresh
1,392-entry primary vulnerability module index, last modified
2026-09-02T19:32:21Z, has no jwalterweatherman record. No implementation or
dependency-metadata commit was made, so the post-implementation exact-quality
and review sequence was not applicable.

## Jwalterweatherman Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-jwalterweatherman-selection.2dd8fdf.h4EKeQ`. Its fully
  verified 76,374-entry manifest SHA-256 is
  `54040aefec38df7415f56da66af04d48c8256dd46ae6de8e06a7c4bbe5ce4d8e`;
  decision-summary SHA-256 is
  `342330b73509caf414a406dda3a14c026bfc59512d09bf6393f710127974dc52`.
- The manifest covers fresh proxy/sumdb/Git/GitHub/release history, exact source
  identity, both module projections, exact-get and tidy diffs, loaded
  population and paths, focused behavior, module tests, repository gates,
  help/API/CLI comparisons, preflight, empty-HOME tests, and vulnerability
  scans. The external manifest-verification log exits 0.
- A first module-mode vulnerability invocation incorrectly supplied a package
  pattern and exited 2; the retained canonical no-argument `-scan=module` run
  exited 3 only for the expected findings and recorded the exact 30-module ID
  set. A first API comparison lacked the old base module in its fresh cache;
  warming from a separate external v1.0.1 archive fixed setup, and both states
  passed.

## Inherited Evidence And Tool Identity

- Cast selection/review/exact-quality/regression evidence remains fully
  verified at 82,586/12/236,782/607 entries with manifest SHA-256 values
  `dfa699f75d0b222d7d64e8fd7a09d05306de98ef4ac5cbf6d253fd33b2dddac0`,
  `850d40f319d13024275b35909e8ff84ce2f9c5ce5959cf8d707461a3b5e5d9ce`,
  `b3d410a70d89e2176daf1a013f97d53e704ef3887f02d67bd6c2564ab0a6d02a`,
  and `97cafb6719f470ab7d057a603b522f0d10418deda2d997e6d5723883efd2e24a`.
  Preserve its corrected NUL-delimited selection manifest for paths containing
  whitespace. Exact quality passed 21 stages, 27/27 Q0-Q2 rows at L2, 80/80
  mutation and meta-mutation kills, host/snapshot/Docker acceptance, and zero
  publication.
- Afero rejection evidence remains fully verified at
  `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`,
  79,866/`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`;
  candidate v1.10.0 fails required complete count-10 module self-tests, so
  v1.9.4 remains selected.
- Accepted go-toml/v2 selection/review/exact-quality/regression evidence
  remains fully verified at 17,093/8/236,599/3,743 entries with manifest
  SHA-256 values
  `ce1a47d85cd2fa0ff3bde19c1953c0dc73914017145abe2f4c5a27c4cf685ff1`,
  `328672c0410e2bfffa058fdf5ec389f4c8640bd73caccde22603b0ad57c9b070`,
  `85f24a00baaf63d62bb6a3a07dc97bffcaa98b8fed1417b7d6f8425ba8d98d3d`,
  and `9601ce867fe389eb1efca59b2ae578052ba9b48670d1112cef3dc66298026ec7`.
- Go-toml v1 and mapstructure decisions remain fully verified at
  37,719/`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`
  and 46,352/`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
  Accepted runewidth selection/review/exact-quality/regression evidence remains
  fully verified at 28,870/15/236,329/4,888 entries with its roadmap-recorded
  hashes.
- Properties, mousetrap, and HCL evidence remains fully verified at
  17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
  and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
- Preserve the go-colorful correction: its 47,633-entry selection manifest has
  exactly one recorded mutable telemetry mismatch; 47,632 stable entries and
  its later review/exact-quality/regression manifests verify. Preserve the
  btree correction: its selection/review/exact-quality manifests verify, while
  the 24,197-entry regression root has exactly 163 recorded mutable cache/HOME
  mismatches superseded by later evidence.
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

Independently evaluate selected indirect `github.com/subosito/gotenv v1.4.2`
as exactly one bounded P7 module group. Resolve canonical latest, every stable
release that could satisfy the retained Go 1.18 floor, source identity,
release qualification, requirements, closure, loaded package population and
behavior, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths, focused environment
and configuration behavior, repository quality, help/API/CLI identity, and
exact vulnerability populations. Implement only an exact floor-compatible
release with explained minimal closure and every gate passing; otherwise
record rejection or no-change without dependency metadata edits. Do not
combine jwalterweatherman, cast, afero, another dependency group, or P8.
