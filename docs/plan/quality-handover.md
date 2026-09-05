# Quality Upgrade Handover

Generated: 2026-09-05T04:24:16+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only go-colorful commit
  `dc4f27c0319c6f023b5750df2148219a4cb32e60`, exact parent
  `f40b0322974eee388d5c064596826898e63f13ca`, clean tree
  `a5f3282ed71f9e5557718d63e5ea8d23020c44e2`. It changes only `go.mod` and
  `go.sum`.
- The go-colorful documentation handoff commit must have exact parent dc4f27c.
  Ordinary and ignored status must be empty. Relative to accepted go-cmp
  commit c314bcb, dependency metadata changes only go-colorful v1.2.0 ->
  v1.4.1 and adds its exact checksum pair.
- The answered go-colorful archive and the magiconair/properties NEXT archive
  must link reciprocally. Only the launcher's mutable header and prompt regions
  may change during handoff.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted `github.com/lucasb-eyer/go-colorful
v1.4.1`. The HCL and mousetrap no-change/rejection decisions remain closed.
Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color
v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and HCL or mousetrap selection changes
remain rejected for their recorded floor, loaded-behavior,
release-qualification, self-test, or exact-latest decisions. P8 remains queued.

## Accepted Go-Colorful Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
`github.com/lucasb-eyer/go-colorful v1.4.1` as both canonical latest stable and
the highest release compatible with the retained Go 1.18 floor. The proxy
lists exactly v1.0.1, v1.0.2, v1.0.3, v1.1.0, v1.2.0, v1.3.0, v1.4.0, and
v1.4.1. `@latest`, `go list -m -u`, and the default branch all resolve v1.4.1,
with zero later commits. There are no prerelease tags, retractions, v2 tags, or
major-path changes. Noncanonical v0.9/v1.0 tags are historical only; all nine
GitHub Release objects are stable, with no draft or prerelease objects.

| Version | Commit / UTC time | Go | Tag and commit signature | Proxy source and sumdb pair |
| --- | --- | --- | --- | --- |
| `v1.2.0` | `d2b05a0d83cca9d610425691c3253d5f36d0ad06`, 2021-01-28T03:22:51Z | 1.12 | lightweight tag; commit unsigned | 38 files identical; `h1:1nnpGOrhyZZuNyfu1QjKiUICQ74+3FNCN69Aj6K7nkY=` / `h1:R4dSotOR9KMtayYi1e77YzuveK+i7ruzyGqttikkLy0=` |
| `v1.3.0` | `680f8257cbbd7f283eaf717de5ce105a6a741bd6`, 2025-09-08T14:15:45Z | 1.12 | lightweight tag; GitHub validates commit signature, local `gpg` unavailable | 43 files identical; `h1:2/yBRLdWBZKrf7gB40FoiKfAWYQ0lqNcbuQwVHXptag=` / same mod sum |
| `v1.4.0` | `960803eeca7760b91ead14a54fabac75e3cfa5d8`, 2026-03-28T13:15:33Z | 1.12 | lightweight tag; commit unsigned | 45 files identical; `h1:UtrWVfLdarDgc44HcS7pYloGHJUjHV/4FwW4TvVgFr4=` / same mod sum |
| `v1.4.1` | `315b48282c63bac7b48ba128d0c87b7f827b2285`, 2026-08-02T08:53:53Z | 1.12 | lightweight tag; commit unsigned | 45 files identical; `h1:1EO+WB73+EH8EVbzlrG3KLAfEypQWVHIBqlTf+2hNss=` / same mod sum |

The 26 commits from v1.2.0 to v1.3.0 add color spaces, blends, distances,
sorting, YAML support, HSV/HCL gray fixes, and faster strict Hex handling. The
five commits to v1.4.0 add CSS Color 4 wide-gamut/D50 support and Stringer. The
two commits to v1.4.1 correct the D50-to-D65 matrix. Each release has no module
requirements, so exact `go get ...@v1.4.1` changes exactly one selection, the
matching main edge, and two additive go.sum lines while retaining the old
v1.2.0 pair. No other selected module, edge, or checksum changes.

Old and selected states both contain 234 selected modules, 3,557 graph edges,
and 429 complete-test packages. Tidy is an unapplied 308 -> 310-line
projection. One go-colorful package loads through `plybuild/cmd ->
go-term-markdown -> ansimage -> go-colorful`; ansimage uses `MakeColor` and
`Color.Hsv`. Parent/current cmd and ansimage consumers plus focused
MakeColor/HSV tests pass at count 10.

Go-colorful's seven-package complete tests and race pass under exact Go
1.26.7. Standalone `go vet ./...` exits 1 only for five legacy unkeyed struct
literals in `doc/colordist` and `doc/palettegens`; default test vet passes, as
do repository build, complete tests/race/vet, and pinned golangci-lint 2.12.2.
Public help and API/CLI reports remain byte-identical at SHA-256 values
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`,
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`,
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

Exact `make quality` passes its ordered 21-stage ledger, all 27 Q0-Q2 rows at
L2, 80/80 mutation kills, host 4/4, fresh snapshot and Docker acceptance, and
valid six-receipt schema-2 evidence. Held, regressed, not-comparable, and dirty
counts are zero. Independent regression repeats preflight, launcher/Make and
audit meta-contracts, snapshot/Docker meta and acceptance, compatibility,
empty-HOME count-2, and clean-tree gates. Its full audit exits expected 1,
never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

Independent govulncheck v1.7.0 scans preserve byte-identical ID populations:
20 Darwin symbol, 30 Darwin module, and 20 Windows symbol. The fresh primary
Go vulnerability module index contains no go-colorful entry.

## Evidence And Tool Identity

- Selection root `/private/tmp/ply-p7-go-colorful-selection.f40b032.KHS4f6`
  has verified 47,633-entry manifest SHA-256
  `a47facc21231accc5c2ed73980c9e3a7542538c6a41d7778de5cfbb410312adc`;
  selection-summary SHA-256 is
  `07467cdfd752a8ee9c64b1b6632370c8e0d75b6e3715090931ec8831b2080658`.
- Commit-bound review root
  `/private/tmp/ply-p7-go-colorful-quality-review.dc4f27c.gs9tGv` has verified
  13-entry manifest SHA-256
  `6444e55370fb37a6c2c06c2c1c280b8fcbd34fc1f94f110a70bc7d067b05673d`;
  manual-evidence SHA-256 is
  `21863ed3334d86e4b32a2bb97be47787cfda8726ac87f4a698bce35bf5911013`.
- Exact-quality root `/private/tmp/ply-p7-go-colorful-quality-final.dc4f27c.B7eOLZ`
  has verified 236,146-entry manifest SHA-256
  `924500cb4512fa9451cf218b119ce8969e32e7f02a67558807d067e3f41017a6`;
  Q0-Q2 scorecard SHA-256 is
  `7ab6fc62901a6963244c46097155d95c72f4a909fbfbdf3b005cd6b0713a5fc1`.
- Regression root
  `/private/tmp/ply-p7-go-colorful-regression-gate.dc4f27c.X2Y4Aj` has verified
  19,011-entry manifest SHA-256
  `2cf393025cf7768ef90fbfa8c0faf378284c3521cb118220dd74ba6fc4735d0e`;
  regression-summary SHA-256 is
  `ac26683176a7ec40bf1da1a9324cbf648ebcd08708e0b343c5dbb008fa86a894`.
- All inherited manifests were verified before the decision. Preserve the
  btree correction: its regression manifest retains its recorded digest and
  24,197 entries, but 163 mutable cache/HOME entries no longer verify; its
  other three manifests verify, and fresh later evidence supersedes it.
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
`github.com/magiconair/properties v1.8.7` as exactly one bounded P7 module
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
not manufacture a dependency commit if v1.8.7 is already canonical latest.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
