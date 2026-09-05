# Quality Upgrade Handover

Generated: 2026-09-05T08:20:54+02:00

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
- Mapstructure required no implementation commit. The current documentation
  handoff commit must have exact parent
  `f974b581baf36f40bb8fa31868371eb0db712faf`, whose exact parent is ca19dcd.
  Relative to accepted go-cmp commit c314bcb, dependency metadata changes only
  go-colorful v1.2.0 -> v1.4.1 and go-runewidth v0.0.14 -> v0.0.17, adding
  exactly those two selected versions' checksum pairs.
- The answered mapstructure archive and the pelletier/go-toml NEXT archive
  must link reciprocally. Only the launcher's mutable header and prompt
  regions may change during handoff. Ordinary and ignored status must end
  empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/mattn/go-runewidth v0.0.17`. The HCL, mousetrap, properties, and
mapstructure no-change/rejection decisions are closed. Go-cmp v0.7.0, Viper
v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1,
latest gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap,
properties, or mapstructure selection changes remain rejected or unnecessary
for their recorded floor, loaded-behavior, release-qualification, self-test,
exact-latest, or closure decisions. Do not revisit them. P8 remains queued.

## No-Change Mapstructure Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
`github.com/mitchellh/mapstructure v1.5.0` as canonical latest stable and the
highest stable release compatible with the retained Go 1.18 floor. The proxy
lists exactly 17 stable releases v1.0.0-v1.5.0. Proxy `@latest`, exact
`go list -m -versions -retracted`, and `go list -m -u` agree on v1.5.0; the
update response has no `Update` field. There are no prereleases, retractions,
original-repository v2 tags, or GitHub Release objects.

Versions v1.0.0-v1.1.2 have no `go` directive. Versions v1.2.0-v1.5.0
declare Go 1.14. Every stable release has zero module requirements, and all 17
are floor-compatible. Canonical latest therefore needs no fallback selection.
The complete release table records each proxy time, Go declaration, checksum
pair, exact tag commit/time, and lightweight-tag identity.

The selected v1.5.0 tag is lightweight commit
`ab69d8d93410fce4361f4912bb1ff88110a81311`,
2022-04-20T22:31:31Z. The lightweight ref has no tag object or tag signature.
GitHub reports the embedded commit signature valid and verified at
2024-11-07T16:42:32Z; local commit verification is unavailable because `gpg`
is absent. Its checksum pair is
`h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=` /
`h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=`. All 13 files in the proxy
archive match the exact tag-commit archive byte-for-byte; the proxy ZIP
SHA-256 is `118d5b2cb65c50dba967fb6d708f450a9caf93f321f8fc99080675b2ee374199`.

The original repository is archived, but the module has no formal
`Deprecated` marker. Its default branch has 12 untagged commits after v1.5.0,
ending at verified commit `8508981c8b6c964e6986dd8aa85490e70ce3c2e2`,
2023-12-16T20:14:59Z. The changelog names unreleased 1.5.1 fixes, but no
original tag qualifies that source as a stable release.

Repository/module-path lineage is deliberately separate:

- `github.com/go-viper/mapstructure v1.6.0` comes from the maintained fork but
  declares `module github.com/mitchellh/mapstructure`; adoption requires an
  explicit `replace`. The original-path proxy returns 404/unknown revision for
  v1.6.0, so it is not a canonical exact update.
- `github.com/go-viper/mapstructure/v2 v2.5.0` is a distinct maintained module
  and import path. It declares Go 1.18 but requires a source migration outside
  this bounded dependency-only group.

Exact `go get github.com/mitchellh/mapstructure@v1.5.0` under the pinned Go
1.26.7 environment exits 0 with zero stdout/stderr and zero go.mod/go.sum
diff. Old and candidate states remain byte-identical at 234 selected modules,
3,557 graph edges, 429 native complete-test packages, five mapstructure
checksum lines, and a 313-line unapplied `go mod tidy -diff` projection. With
no module requirements, the exact minimal selection, edge, checksum, package,
and metadata closure is empty.

One mapstructure package loads through `plybuild/cmd -> spf13/viper ->
mapstructure`. Viper v1.15.0 directly uses `DecoderConfig`,
`ComposeDecodeHookFunc`, `StringToTimeDurationHookFunc`,
`StringToSliceHookFunc`, `NewDecoder`, and `Decode`. Ten fresh-process focused
Viper unmarshalling runs pass in both old and candidate states, the Viper
source manifests remain unchanged, and the three Ply consumer packages pass
at count 10. A single-process Viper `-count=10` diagnostic was discarded
because Viper package-global state contaminates later repetitions; it is not a
candidate regression.

The selected module's complete tests, count-10 tests, race, and vet pass.
Repository build, complete tests/race/vet, Windows build, pinned
golangci-lint 2.12.2, CLI surface, and byte-identical help/API/CLI reports
pass. The help/API/CLI SHA-256 values remain
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`,
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`,
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

The final complete preflight passes all repository, launcher, Make,
compatibility, snapshot/Docker meta, 80/80 mutation, and 15 audit-meta
controls. Two earlier preflight attempts are retained as environment-only
diagnostics: command-line `GOLANGCI_LINT` leaked through `MAKEFLAGS`, then two
ignored standalone compatibility reports made the disposable candidate clone
non-clean. Environment-scoped tools, empty `MAKEOVERRIDES`, and a clean clone
produce the passing result without apparatus changes. No post-implementation
quality run was required because the selection did not change.

Pinned govulncheck v1.7.0 independently preserves exact and identical
Darwin-symbol/Darwin-module/Windows-symbol ID populations 20/30/20. The fresh
primary Go vulnerability index, updated 2026-09-02T19:12:04Z, has no entry for
the selected original path. GO-2025-3787 and GO-2025-3900 are indexed under
the maintained fork paths, not the selected original module path.

## Evidence And Tool Identity

- Mapstructure decision root
  `/private/tmp/ply-p7-mapstructure-selection.f974b58.GbAMTJ` fully verifies
  46,352 entries at manifest SHA-256
  `f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`;
  selection-summary SHA-256 is
  `8910174c8f8fa387a32cd49871dfbd37668ff4a2fbf1c51c2dae6fdc2612eff3`.
- Accepted runewidth selection/review/exact-quality/regression roots remain
  fully verified at 28,870/15/236,329/4,888 entries and manifest SHA-256
  values `2d6986055d24d4937dcbbaa8fa1e46da6a7bcc1e94f481125343284c96687193`,
  `f885b995dc6192cb7ef2d8c9241f93cb44a885deda7faec7e0415f68ce21d764`,
  `28fbbde41059dd806dfea5653f273a390b6b1fa2f60ebeec2e3d654cc792e2ed`,
  and `b83091c1c8a52ec2176c4c3fd7a4c3c5e456adf7a61ee2c6b6f7f97706eb4894`.
- Recovery root `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` retains its
  fully verified two-entry manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Preserve the go-colorful correction: the selection manifest retains its
  digest and 47,633 entries, with one recorded mutable telemetry mismatch;
  the other 47,632 entries and stable summary remain valid, and the later
  review/exact-quality/regression manifests fully verify.
- Preserve the btree correction: its selection, review, and exact-quality
  manifests fully verify. Its 24,197-entry regression root reproduces exactly
  163 recorded mutable cache/HOME mismatches; fresh later evidence supersedes
  those mutable entries.
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
`github.com/pelletier/go-toml v1.9.5` as exactly one bounded P7 module group.
Resolve canonical latest and every potentially compatible release from fresh
Go proxy, checksum-database, upstream, and primary Go vulnerability evidence.
Treat latest, Go floor, source identity, repository/module-path lineage,
closure, loaded status, and behavior as unknown. A distinct
`github.com/pelletier/go-toml/v2` module is not an in-place v1 selection.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused behavior when loaded, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata. Do
not manufacture a dependency commit if v1.9.5 is already canonical latest.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
