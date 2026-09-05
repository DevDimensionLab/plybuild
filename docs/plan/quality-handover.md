# Quality Upgrade Handover

Generated: 2026-09-05T09:04:35+02:00

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
- Pelletier go-toml v1 required no implementation commit. The current
  documentation handoff commit must have exact parent
  `09cc40bcab6f1b03273d726df7dd5ce832f92af6`, whose exact parent is
  `f974b581baf36f40bb8fa31868371eb0db712faf`, whose exact parent is ca19dcd.
  Relative to accepted go-cmp commit c314bcb, dependency metadata changes only
  go-colorful v1.2.0 -> v1.4.1 and go-runewidth v0.0.14 -> v0.0.17, adding
  exactly those two selected versions' checksum pairs.
- The answered go-toml v1 archive and the go-toml `/v2` NEXT archive must link
  reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted
`github.com/mattn/go-runewidth v0.0.17`, plus the no-change go-toml v1
evaluation. The HCL, mousetrap, properties, mapstructure, and go-toml v1
no-change/rejection decisions are closed. Go-cmp v0.7.0, Viper v1.16.0,
Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap,
properties, mapstructure, or go-toml v1 selection changes remain rejected or
unnecessary for their recorded floor, loaded-behavior, release-qualification,
self-test, exact-latest, or closure decisions. Do not revisit them. P8 remains
queued.

## No-Change Pelletier Go-Toml V1 Group

Fresh Go proxy, checksum-database, GitHub, and upstream Git evidence establish
`github.com/pelletier/go-toml v1.9.5` as canonical latest stable and the
highest stable release compatible with the retained Go 1.18 floor. The proxy
lists exactly 28 stable releases v0.1.0-v1.9.5. Proxy `@latest`, exact
`go list -m -versions -retracted`, and `go list -m -u` agree on v1.9.5; the
update response has no `Update` field. The 28 upstream stable tags match the
proxy list, 26 have stable GitHub Release objects, and there are no v1
prereleases or retractions.

Versions v0.1.0-v1.2.0 have no `go` directive. Versions v1.3.0-v1.9.5
declare Go 1.12. Every stable release is floor-compatible, so canonical latest
needs no fallback. The complete release table records every proxy time, Go
declaration, requirement count, checksum pair, exact tag commit/time, tag
type, tag-signature status, and floor decision.

Selected v1.9.5 is lightweight tag commit
`fed1464066413075eac02cd4dc368b5221845541`, committed
2022-01-05T14:17:32Z. Its stable GitHub Release was published
2022-04-21T23:21:51Z. The lightweight ref has no tag object or tag signature.
GitHub reports the embedded commit signature valid and verified at
2024-11-12T16:10:24Z; local commit verification is unavailable because `gpg`
is absent. Its checksum pair is
`h1:4yBQzkHv+7BHq2PQUZF3Mx0IYxG7LsP222s7Agd3ve8=` /
`h1:u1nR/EPcESfeI/szUZKdtJ0xRNbUoANCkoOuaOx1Y+c=`.

All 65 files eligible for the v1.9.5 module ZIP match the exact tag commit
byte-for-byte. Go correctly excludes six files under the nested `benchmark`
module. The proxy ZIP SHA-256 is
`de3dcda660cc800cd86d03273a25956d67f416e8fcbe4d2001a2cb4a01e6ac60`;
the eligible file-manifest SHA-256 is
`70106b3cd0fd7c623c4768f809217108f6869de674c4cb0fd7554ec8d6c6b692`.
The v1.9.4 -> v1.9.5 history is four commits and seven files, adding uint64
parsing, fixing the LoadBytes invalid type assertion, adding SECURITY.md, and
updating v2 guidance.

The repository is active and unarchived for its default `v2` branch. The v1
module has no formal `Deprecated` marker, but `master` states that v1 will
receive no updates and strongly recommends v2. Its single commit after v1.9.5
only adds that README notice; no later v1 tag qualifies it as a release.
`github.com/pelletier/go-toml/v2` is a distinct module and import path, not an
in-place v1 selection. Its observed lineage data must not substitute for an
independent `/v2` evaluation.

Exact `go get github.com/pelletier/go-toml@v1.9.5` under the pinned Go 1.26.7
environment exits 0 with zero stdout/stderr and zero go.mod/go.sum diff. Old
and candidate states remain byte-identical at 234 selected modules, 3,557
graph edges, 429 native complete-test packages, three v1 checksum lines, and a
313-line unapplied `go mod tidy -diff` projection. V1.9.5 has no requirements,
so its exact minimal selection, edge, checksum, package, and metadata closure
is empty.

No v1 package loads, there is no direct v1 consumer, and `go mod why -m` says
the main module does not need it. No focused loaded behavior therefore
applies. Five packages from the distinct `/v2` module do load; Viper's TOML
codec imports `/v2` directly.

The selected v1 module's complete tests, count-10 tests, and race pass across
six packages. Standalone `go vet ./...` reports 132 legacy unkeyed `Position`
literal findings in `query` test files. This is retained as a historical
upstream diagnostic, not a candidate regression: the exact selection is
unchanged, default complete tests pass, and repository vet passes.

Repository build, complete tests/race/vet, Windows build, pinned
golangci-lint 2.12.2, CLI surface, empty-HOME count-2, and byte-identical
help/API/CLI reports pass. The help/API/CLI SHA-256 values remain
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`,
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`,
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.

The final complete preflight passes all repository, launcher, Make,
compatibility, snapshot/Docker meta, 80/80 mutation, and 15 audit-meta
controls. Two earlier preflight attempts retain the established launcher
signal-fixture timing diagnostic. The isolated launcher suite then passes
62/62 both directly and through Make, and the final complete preflight passes
without apparatus changes. No post-implementation quality run was required
because the selection did not change.

Pinned govulncheck v1.7.0 independently preserves exact and identical
Darwin-symbol/Darwin-module/Windows-symbol ID populations 20/30/20. The fresh
primary Go vulnerability module index, last modified
2026-09-02T19:32:20Z, has no go-toml v1 or v2 entry. Dependency metadata
remained byte-identical and no implementation commit was manufactured.

## Evidence And Tool Identity

- Go-toml v1 decision root
  `/private/tmp/ply-p7-go-toml-selection.09cc40b.jGVQz6` fully verifies
  37,719 entries at manifest SHA-256
  `e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`;
  selection-summary SHA-256 is
  `85deac527b69ebbfd2608a970684b3f915563d48a538644f3da357961adc268e`.
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
`github.com/pelletier/go-toml/v2 v2.0.7` as exactly one bounded P7 module
group. Resolve canonical latest and every potentially compatible release from
fresh Go proxy, checksum-database, upstream, and primary Go vulnerability
evidence. Treat latest, Go floor, source identity, release qualification,
closure, loaded path, and behavior as unknown. Do not reuse the completed v1
decision as `/v2` evidence.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused Viper TOML behavior, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata. Do
not manufacture a dependency commit if v2.0.7 is already the selected
floor-compatible decision.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
