# Quality Upgrade Handover

Generated: 2026-09-05T14:46:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only cast commit
  `cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
  `17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
  `485690c010cef3b02385d8f2471d3ef53611abe9`. It changes only `go.mod` and
  `go.sum`.
- This documentation handoff must be the direct child of `cf4fd493`. The next
  dependency implementation, if any, must use the resulting documentation
  commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves are go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Minimal closures
  select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty v0.3.1,
  and rogpeppe/go-internal v1.9.0. Those four accepted groups add exactly 15
  checksum lines.
- The answered cast archive and sole NEXT `spf13/jwalterweatherman` archive
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted
`github.com/spf13/cast v1.5.1`. HCL, mousetrap, properties, mapstructure,
go-toml v1, and Afero decisions are closed. Go-cmp v0.7.0, Viper v1.16.0,
Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and changes to any closed
group remain rejected or unnecessary for their recorded floor,
loaded-behavior, release-qualification, self-test, exact-latest, or closure
decisions. Do not revisit them. P8 remains queued.

## Accepted Spf13 Cast Group

Fresh Go proxy, checksum-database, GitHub, upstream Git, and primary Go
vulnerability evidence establishes v1.10.0 as canonical latest stable and
v1.5.1 as the highest stable release compatible with the retained Go 1.18
floor. The proxy lists exactly 17 stable releases, with no prereleases,
retractions, formal deprecation marker, or later major module path. Releases
v1.0.0-v1.4.1 have no Go directive; v1.5.0-v1.5.1 declare Go 1.18;
v1.6.0-v1.8.0 declare Go 1.19; and v1.9.0-v1.10.0 declare Go 1.21.0.

Canonical latest v1.10.0 is annotated tag object
`3efe057bdf41dab273f8e65893fbcb89c1a4dc3b`, peeled commit
`fc73346bfc4e6597bc520fb6eea04360299e77d2` at
2025-09-08T16:45:31Z. GitHub verifies the tag and commit signatures. Candidate
v1.5.1 is unsigned lightweight tag commit
`bbed5559a59db54120dda79106f85f9ec0afb038` at
2023-05-15T08:28:46Z and was published as a stable GitHub Release at
2023-05-15T08:31:31Z. Selected v1.5.0 and first-incompatible v1.6.0 are also
unsigned lightweight commits. Local signature verification is unavailable
because `gpg` is absent.

The active, unarchived, undisabled upstream is not a fork. Master commit
`6fd6afdc9df068c662e923a7c3ad04fda4647121` at
2026-04-06T18:43:42Z is 32 commits after v1.10.0; those commits are unreleased
and excluded from stable selection. V1.10.0 has a nested generator module;
the Go proxy correctly excludes its three files.

Candidate checksum pair is
`h1:R+kOtfhWQE6TVQzY+4D7wJLBgkdVasCEFxSUBYBYIlA=` /
`h1:b9PdjNptOpzXr7Rq1q9gJML/2cdGQAo69NKzQ10KN48=`. All 13 candidate proxy
files match the exact tag commit byte-for-byte; ZIP SHA-256 is
`15f13bcba5bcfe815edcc71baccecd2ac202f55e2817a48676972ea9502833d8`
and file-manifest SHA-256 is
`417dcd9f6a303102262cfb6ebe037f52553f33279ba5f75c077bcd6badde580b`.
Selected v1.5.0, incompatible-boundary v1.6.0, and latest v1.10.0 also have
exact proxy/tag source identity. The selected-to-candidate range is 12 commits
across nine files, +137/-59, adding number/duration Boolean conversions,
prioritized date layouts, and maintenance/test-dependency changes.

Exact `go get github.com/spf13/cast@v1.5.1` changes four selected versions:
cast v1.5.0 -> v1.5.1, quicktest v1.14.3 -> v1.14.4, kr/pretty v0.3.0 ->
v0.3.1, and rogpeppe/go-internal v1.6.1 -> v1.9.0. Cast directly raises
quicktest, whose closure raises pretty and go-internal. Their declarations are
Go 1.13, Go 1.12, and Go 1.17. Candidate go-cmp v0.5.9 loses to selected
v0.6.0; kr/text stays v0.2.0; pkg/diff does not move. This is the minimal MVS
closure.

Both states select 234 modules and the byte-identical 429-package complete
population. Graph edges change 3,565 -> 3,564: candidate cast has five rather
than six outbound edges and drops x/xerrors. Go.mod changes only cast. Go.sum
adds the cast checksum pair plus quicktest, pretty, and go-internal content
checksums: exactly five additions and zero removals. The unapplied tidy
projection changes 321 -> 332 lines; its exact SHA-256 is
`20291af321d20d3e1b032e1f8bc300af51dbca0376af1f7e964354746f86898c`.
Tidy was not used as implementation.

Exactly one cast package loads through
`plybuild/cmd -> spf13/viper -> spf13/cast`. Real consumers are Viper root
and its dotenv, ini, and javaproperties encoders, using 17 cast symbols. Ten
fresh focused Viper root
runs, each encoder suite at count 10, and Ply `./cmd` at count 10 pass in both
selected and candidate states. One naive Viper root `-count=10` process
exposed Viper's pre-existing global-state leak in TestDefault; the canonical
ten-fresh-process run passes both versions.

Candidate module verification, complete tests at count 1 and count 10, race,
and vet pass for its sole package. Repository build, complete tests, race,
vet, Windows build, pinned lint, CLI surface, and launcher contract pass.
Root/status/upgrade/build help is byte-identical, as are API report SHA-256
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and CLI report SHA-256
`955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

Govulncheck v1.7.0 preserves exact and identical old/candidate and committed
Darwin-symbol/Darwin-module/Windows-symbol populations 20/30/20. The primary
vulnerability module index was last modified 2026-09-02T19:32:21Z, contains
1,392 entries, and has no cast record.

## Cast Quality And Evidence

- Selection evidence `/private/tmp/ply-p7-cast-selection.17f7277.9nRKpS`
  fully verifies 82,586 entries at manifest SHA-256
  `dfa699f75d0b222d7d64e8fd7a09d05306de98ef4ac5cbf6d253fd33b2dddac0`;
  decision-summary SHA-256 is
  `13c7ab63019a891d66508638b6ade40155ffb340548f36a17e761b204da89c87`.
  The first top-level manifest attempt split paths containing `Library/Application
  Support`; it was detected before sealing. The retained correction explains
  the superseded attempt, and the final NUL-delimited manifest covers and
  verifies every file.
- Commit-bound schema-2 review evidence is its `quality-review-final` child.
  Its 12-entry manifest SHA-256 is
  `850d40f319d13024275b35909e8ff84ce2f9c5ce5959cf8d707461a3b5e5d9ce`;
  review-summary SHA-256 is
  `8ff415dd7d4cdb0153937a61a9076b1a112c33f8048c923b948d6706520d31d7`
  and manual-evidence SHA-256 is
  `b38814b133a4f30d4a6a113fa7c6581627d33505a7fa8535b3f7c9c4bb710d68`.
  All 112 governed source files and 258 subjects match the parent, all six
  criterion receipts validate, and the focused audit passes 6/6.
- Exact `make quality` exits 0 at
  `/private/tmp/ply-p7-cast-exact-quality-final.cf4fd49.wlsG4F/output`. Its
  fully verified 236,782-entry manifest SHA-256 is
  `b3d410a70d89e2176daf1a013f97d53e704ef3887f02d67bd6c2564ab0a6d02a`;
  quality-summary SHA-256 is
  `c8c2bc454c816656297ab1622da50b534c84b090c0a2ddac6366f065cca9d89d`.
  All 21 stages pass; all 27 Q0-Q2 rows pass at L2 with zero held, regressed,
  not-comparable, or dirty counts; both meta and actual mutation populations
  kill 80/80; four host features, fresh snapshot, and Docker acceptance pass.
  Docker builds once, runs ten containers, records 54 calls, and publishes
  zero times.
- Two exact-quality attempts stopped safely after 19 stages: the first fresh
  HOME hid Docker's buildx plugin; the second no-cache Docker build reproduced
  the established nested launcher signal-fixture timing diagnostic. They are
  retained at `/private/tmp/ply-p7-cast-exact-quality.cf4fd49.9fqPO0/output`
  and
  `/private/tmp/ply-p7-cast-exact-quality-retry.cf4fd49.v0wyQ6/output`.
  Isolated exact Docker confirmation then passed at
  `/private/tmp/ply-p7-cast-docker-isolated.cf4fd49.xhymf9`, followed by the
  complete passing quality run.
- Independent regression evidence at
  `/private/tmp/ply-p7-cast-regression-gate.cf4fd49.sAkB56` fully verifies 607
  entries at manifest SHA-256
  `97cafb6719f470ab7d057a603b522f0d10418deda2d997e6d5723883efd2e24a`;
  regression-summary SHA-256 is
  `533ccdf2704153e7d4f467e3ceb75f932202a1367ec38dc1c3f8df51ec3a458f`.
  Exact selection/graph/packages/checksums/path, module self-tests, repository
  gates, byte-identical help/API/CLI, empty-start-HOME count-2, and 20/30/20
  vulnerability identity pass. The fresh focused audit is 27/27 PASS at L2;
  the full audit exits expected 1, never 2, only for Q3.1, Q3.3, Q3.4, and
  Q3.7. The final ordinary and ignored status is empty.

## Inherited Evidence And Tool Identity

- Afero rejection evidence remains fully verified at
  `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`,
  79,866/`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`;
  its candidate fails the required complete count-10 module self-test, so
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
  Accepted runewidth evidence remains fully verified at
  28,870/15/236,329/4,888 entries and its roadmap-recorded hashes.
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

Independently evaluate selected indirect
`github.com/spf13/jwalterweatherman v1.1.0` as exactly one bounded P7 module
group. Resolve canonical latest, every stable release that could satisfy the
retained Go 1.18 floor, source identity, release qualification, requirements,
closure, loaded package population and behavior, module self-tests, and
primary vulnerability data from fresh evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths, focused logging and
configuration behavior, repository quality, help/API/CLI identity, and exact
vulnerability populations. Implement only an exact floor-compatible release
with explained minimal closure and every gate passing; otherwise record
rejection or no-change without dependency metadata edits. Do not combine
gotenv, cast, afero, another dependency group, or P8.
