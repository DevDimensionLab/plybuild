# Quality Upgrade Handover

Generated: 2026-09-05T07:31:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only go-runewidth commit
  `ca19dcd4da9320112e6c9f0c5db507ccbe8b88e5`, exact parent
  `5d0fa3ae0fff98784c349cc23e001018eb9637e1`, clean tree
  `273c9b635e988ddd686d6c83846f7320bb9e4e33`. It changes only `go.mod` and
  `go.sum` from v0.0.14 to v0.0.17 and adds the candidate's exact checksum
  pair.
- The documentation handoff commit must have exact parent ca19dcd. Relative to
  accepted go-cmp commit c314bcb, dependency metadata changes only go-colorful
  v1.2.0 -> v1.4.1 and go-runewidth v0.0.14 -> v0.0.17, adding exactly the two
  checksum pairs for those selected versions.
- The answered go-runewidth archive and the mitchellh/mapstructure NEXT archive
  must link reciprocally. Only the launcher's mutable header and prompt regions
  may change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/mattn/go-runewidth v0.0.17`. The HCL, mousetrap, and properties
no-change/rejection decisions are closed. Go-cmp v0.7.0, Viper v1.16.0, Emoji
v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap, or properties
selection changes remain rejected for their recorded floor, loaded-behavior,
release-qualification, self-test, exact-latest, or closure decisions. Do not
revisit them. P8 remains queued.

## Accepted Go-Runewidth Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
`github.com/mattn/go-runewidth v0.0.29` as canonical latest stable at commit
`218f489f6718aaf2c68285d71dfc363e13387075`,
2026-09-02T01:31:19Z. Proxy `@latest`, `go list -m -u`, the v0.0.29 tag, and
upstream `master` agree, with zero later commits. The proxy lists 29 stable
versions v0.0.1-v0.0.29. There are no prerelease tags, retractions, v1 or
higher tags, major-path changes, or GitHub Release objects.

Versions v0.0.1-v0.0.4 have no Go directive, v0.0.5-v0.0.17 declare Go 1.9,
v0.0.18-v0.0.25 declare Go 1.20, and v0.0.26-v0.0.29 declare Go 1.23. Thus
v0.0.17 is the highest stable release compatible with the retained Go 1.18
floor. Canonical latest and all 11 other higher releases exceed that floor.
Compatibility follows the declared Go version, not a modern-toolchain build.

Relevant release identities are:

| Version | Commit / commit UTC | Go | Source and primary signature status |
| --- | --- | --- | --- |
| `v0.0.14` | `2c6a438f68cfe01255a90824599da41fdf76d1e2` / 2022-09-20T12:35:16Z | 1.9 | 16 proxy/tag files identical; lightweight tag; GitHub commit valid |
| `v0.0.15` | `44b7c5b4d67df8ca22917b6800c158a6d3be3560` / 2023-07-23T16:42:41Z | 1.9 | 16 files identical; lightweight tag; GitHub commit valid |
| `v0.0.16` | `6ceadc68530e7bfea8cba17d6523bed32912d4fa` / 2024-07-22T12:40:34Z | 1.9 | 16 files identical; lightweight tag; GitHub commit valid |
| `v0.0.17` | `94c0db1df07f7755a1e3e962becc01cb5fe8a086` / 2025-09-25T15:58:59Z | 1.9 | 16 files identical; lightweight tag; GitHub commit valid |
| `v0.0.18` | `61f04f380c801eced5965598b97580446ca7c931` / 2025-09-29T13:19:08Z | 1.20 | 19 files identical; lightweight tag; GitHub commit valid |
| `v0.0.29` | `218f489f6718aaf2c68285d71dfc363e13387075` / 2026-09-02T01:31:19Z | 1.23 | 20 files identical; lightweight tag; GitHub commit valid |

Local commit-signature verification is unavailable because `gpg` is absent.
Only v0.0.28 uses an annotated tag; it is unsigned and points to a
GitHub-verified signed commit. The other 28 stable tags are lightweight. The
accepted v0.0.17 checksum pair is
`h1:78v8ZlW0bP43XfmAfPsdXcoNCelfMHsDmd/pkENfrjQ=` /
`h1:Jdepj2loyihRzMpdS35Xk/zdY8IAYHsh153qUoGf23w=`.

The v0.0.14 -> v0.0.17 history refreshes CI, updates width tables to Unicode
15.1, fixes benchmark expectations, and treats Windows Terminal `WT_SESSION`
as non-East-Asian ambiguous width. Version v0.0.18 replaces uniseg with a
uax29 implementation while raising the Go floor to 1.20. Later incompatible
releases add Unicode 17 tables, cluster-width fixes, performance work,
`IsCombiningWidth`, `TruncatePrefix`, and a Go 1.23 floor.

Exact `go get github.com/mattn/go-runewidth@v0.0.17` has the authorized minimal
closure. Selected modules remain 234 and graph edges remain 3,557. Exactly one
selection changes; two edges are relabeled: main -> runewidth and runewidth ->
uniseg v0.2.0. The latter requirement is unchanged from v0.0.14 and MVS retains
the already-selected uniseg v0.4.7. The native complete-test population stays
byte-identical at 429 packages. Go.sum adds only the v0.0.17 module and go.mod
sums while retaining the historical v0.0.14 pair. Tidy stays unapplied and
projects 310 -> 313 lines.

One package loads through `plybuild/cmd -> go-term-markdown -> go-term-text ->
go-runewidth`. The direct consumer calls `StringWidth`, `Truncate`, `FillRight`,
and `RuneWidth`. Candidate and selected module tests/race/vet pass; candidate
Windows tests compile. Old/candidate go-term-text complete and focused suites
and the three Ply consumer packages pass at count 10. The refreshed Unicode
tables and Windows ambiguity behavior produce no focused regression.

Repository build, complete tests/race/vet, pinned golangci-lint 2.12.2,
Windows build, CLI surface, and byte-identical help/API/CLI reports pass. The
help SHA-256 remains
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`.
Exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
remain 20/30/20, and the primary Go vulnerability module index has no
runewidth entry.

The dependency-only implementation is ca19dcd. Exact `make quality` passes all
21 stages, all 27 Q0-Q2 rows at L2, 80/80 mutations, host/snapshot/Docker
acceptance, and zero held, regressed, not-comparable, or dirty counts. The
independent regression gate passes module and consumer behavior, repository
and launcher contracts, preflight, vulnerability comparison, empty-HOME
count-2, and cleanliness. The separate full audit exits expected 1, never 2,
only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Several quality attempts were discarded before the successful run because
background execution changed launcher signal behavior, the launcher parser
raced under load, an empty external HOME lacked Docker buildx, and globally
using system Python broke Docker-helper nanosecond timestamp parsing. The
successful external runtime uses a native dispatcher that selects
`/usr/bin/python3 -S` only for generated `validate-events.py` and Homebrew
Python otherwise, plus an empty-HOME symlink to the verified Docker Desktop
buildx plugin. The dispatcher binary/source SHA-256 values are
`d2f649ec9d50c5e23bf424de4ef27319550359c066707ac31ac2538b15d02b97` and
`bd32f8501fea390f5c8c9c4a8c8ae063033b84611505d0c8a201be9c8f67dba2`.
The buildx plugin was separately hash-verified. These are external execution
accommodations, not repository or quality-apparatus changes. The successful
quality root records all discarded attempts.

## Evidence And Tool Identity

- Selection root `/private/tmp/ply-p7-go-runewidth-selection.5d0fa3.dmbKoB`
  fully verifies 28,870 entries at manifest SHA-256
  `2d6986055d24d4937dcbbaa8fa1e46da6a7bcc1e94f481125343284c96687193`;
  selection-summary SHA-256 is
  `919632a255cbca3a4e0899fc3d0226bcff5c1afd484e22771f894534ace61362`.
- Review root `/private/tmp/ply-p7-go-runewidth-quality-review.ca19dcd.UgdZ2N`
  fully verifies 15 entries at manifest SHA-256
  `f885b995dc6192cb7ef2d8c9241f93cb44a885deda7faec7e0415f68ce21d764`;
  review-summary SHA-256 is
  `13c16c8736779c320b78fe9b5ec4ef404b07ba013bd387652f7273e08ac3a74b`.
- Exact-quality root
  `/private/tmp/ply-p7-go-runewidth-quality-final.ca19dcd.t6LExP` fully
  verifies 236,329 entries at manifest SHA-256
  `28fbbde41059dd806dfea5653f273a390b6b1fa2f60ebeec2e3d654cc792e2ed`;
  quality-summary SHA-256 is
  `b5eebc7e4c89c3993453cc6456bb9587c47fd3aa1d30810527cc7fb9888c9159`.
- Regression root
  `/private/tmp/ply-p7-go-runewidth-regression-gate.ca19dcd.O6Uxyq` fully
  verifies 4,888 entries at manifest SHA-256
  `b83091c1c8a52ec2176c4c3fd7a4c3c5e456adf7a61ee2c6b6f7f97706eb4894`;
  regression-summary SHA-256 is
  `ec2ca15ba9624472480abacace094b73dff0348effe86ee900b764121c6b72fb`.
- Recovery root `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` retains its
  fully verified two-entry manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Preserve the go-colorful correction: its 47,633-entry selection manifest
  retains its digest, but one mutable Go telemetry counter no longer verifies;
  the other 47,632 entries and stable summary
  `07467cdfd752a8ee9c64b1b6632370c8e0d75b6e3715090931ec8831b2080658`
  remain valid. Its review, exact-quality, and regression manifests still
  fully verify.
- Preserve the btree correction: 163 mutable cache/HOME entries in its
  regression manifest no longer verify; its other three manifests verify and
  fresh later regression evidence supersedes that root.
- Rejected properties evidence remains fully verified at
  `/private/tmp/ply-p7-properties-selection.469049f.iT6qcB`,
  17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`.
  No-change mousetrap and rejected HCL remain fully verified at
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`
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

Independently evaluate selected indirect
`github.com/mitchellh/mapstructure v1.5.0` as exactly one bounded P7 module
group. Resolve canonical latest and every potentially compatible release from
fresh Go proxy, checksum-database, upstream, and primary Go vulnerability
evidence. Treat latest, Go floor, source identity, closure, loaded status, and
behavior as unknown.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused behavior when loaded, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata. Do
not manufacture a dependency commit if v1.5.0 is already canonical latest.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
