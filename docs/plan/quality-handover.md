# Quality Upgrade Handover

Generated: 2026-09-05T02:09:20+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only go-cmp commit
  `c314bcb440b5f94871d71a249bca7ae7f87d5543`, exact parent
  `cbdb0a915d7d30bb6633cfbb2ccf0323aea1b77a`, clean tree
  `9a212b377cb08e7d6e3044fb6dae071058b6242b`.
- The mousetrap evaluation made no dependency or production change. Its final
  documentation handoff must have exact parent
  `516ae39b7d8efd36e20d7fe07bc3c6a8d22dbaae`; `go.mod` and `go.sum` remain
  byte-identical to `c314bcb`. Ordinary and ignored status must be empty.
- The answered mousetrap archive and the lucasb-eyer/go-colorful NEXT archive
  must link reciprocally. Only the launcher's mutable header and prompt regions
  may change during handoff.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted `github.com/google/go-cmp
v0.6.0`, and the rejected/no-change HCL and mousetrap evaluations. Latest
go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3,
fatih/color v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and HCL or mousetrap selection changes
remain rejected for their recorded floor, loaded-behavior,
release-qualification, self-test, or exact-latest decisions. P8 remains
queued.

## Accepted No-Change Mousetrap Group

Fresh Go proxy, checksum-database, and upstream evidence agree that selected
indirect `github.com/inconshreveable/mousetrap v1.1.0` is already the
canonical latest stable module release and the highest release compatible with
the retained Go 1.18 floor. The proxy lists exactly v1.0.0, v1.0.1, and
v1.1.0; `@latest` and `go list -m -u` select v1.1.0 with no `Update` field.
There are no listed prereleases, retract directives, v2 tags/module path, or
commits after v1.1.0. Upstream's noncanonical v1.0 and v1.1 alias tags point to
the exact v1.0.0 and v1.1.0 commits. GitHub Releases contains only an old v1.0
alias object, so stable semver tags plus the Go proxy are the canonical module
resolution.

| Version | Commit / UTC time | Go | Tag and commit signature | Proxy source and sumdb pair |
| --- | --- | --- | --- | --- |
| `v1.0.0` | `76626ae9c91c4f2a10f34cad8ce83ea42c93bb75`, 2014-10-17T20:07:13Z | none | lightweight tag; commit unsigned | 5 files identical; `h1:Z8tu5sraLXCXIcARxBp/8cbvlwVa7Z1NHg9XEKhtSvM=` / `h1:PxqpIevigyE2G7u3NXJIT2ANytuPF1OarO4DADm73n8=` |
| `v1.0.1` | `3a66f25f8779fad208598f21472174ef7b35c3ec`, 2022-08-07T15:49:23Z | 1.18 | lightweight tag; commit unsigned | 6 files identical; `h1:U3uMjPSQEBMNp1lFxmllqCPM6P5u/Xq7Pgzkat/bFNc=` / `h1:vpF70FUmC8bwa3OWnCshd2FqLfsEA9PFc4w1p2J65bw=` |
| `v1.1.0` | `4e8053ee7ef85a6bd26368364a6d27f1641c1d21`, 2022-11-27T22:01:53Z | 1.18 | lightweight tag; GitHub validates the commit signature, local verification lacks `gpg` | 5 files identical; `h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=` / `h1:vpF70FUmC8bwa3OWnCshd2FqLfsEA9PFc4w1p2J65bw=` |

The three commits from v1.0.0 to v1.0.1 only expand the license and add the
Go-1.18 module declaration. The four commits from v1.0.1 to v1.1.0 remove
Go-1.3/1.4 legacy files, modernize build tags, and consolidate the Windows
implementation; there are zero later default-branch commits.

Exact `go get github.com/inconshreveable/mousetrap@v1.1.0` in a fresh
candidate archive emits no output and changes zero `go.mod` or `go.sum` bytes.
Old and candidate states retain byte-identical selections and populations: 234
modules, 3,557 graph edges, 429 native complete-test packages, 433 Windows
complete-test packages, the same four mousetrap checksum lines, and the same
308-line unapplied tidy projection. Mousetrap v1.1.0 has no module
requirements, so the exact minimal closure is empty: no selection, edge,
checksum, or metadata line changes.

Mousetrap is not loaded on Darwin. Windows loads its one package directly from
Cobra through `plybuild/cmd -> github.com/spf13/cobra ->
github.com/inconshreveable/mousetrap`. Old and candidate Ply `cmd` consumers
pass at count 10, both states cross-build all packages for windows/amd64, and
an exact Cobra v1.10.2 source copy passes all packages at count 10 and compiles
its Windows test binary. A diagnostic Cobra run directly from the read-only
module cache fails because `TestDeadcodeElimination` ignores a failed source-
directory `Mkdir`; the writable byte-identical source replay passes and leaves
its manifest unchanged, so this is not a loaded-behavior failure.

Mousetrap's complete one-package `go test`, race, and vet commands pass, though
the module contains no test files. Repository build, complete tests, race, vet,
pinned golangci-lint 2.12.2, CLI surface, byte-identical public help, and
identical API/CLI reports pass. Help/API/CLI SHA-256 values remain respectively
`ea32c45fa1b86fbe46c8a0cc244f37157201610ba6cb10a9bd80dfb11e9d3d47`,
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`,
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
No dependency commit or post-implementation quality run was appropriate for
the exact no-op selection.

Independent old/candidate govulncheck v1.7.0 scans preserve exact identical 20
Darwin-symbol, 30 Darwin-module, and 20 Windows-symbol ID populations. The
fresh primary Go vulnerability module index contains no mousetrap entry.

## Evidence And Tool Identity

- No-change mousetrap evidence root
  `/private/tmp/ply-p7-mousetrap-selection.516ae39.PzdXjl` has a fully verified
  29,283-entry manifest SHA-256
  `71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`;
  selection-summary SHA-256 is
  `64371a02034a3265fa879330e05bcdb477c5c2560de16fa9e57c62e4a5c2e08a`.
- All accepted prerequisites were freshly verified entry-by-entry before the
  decision: two-entry toolchain recovery `1b30193f...e83b`, 21,427-entry HCL
  selection `ad2954c3...acff`, 29,812-entry go-cmp selection
  `dcd18baa...1797`, 9,487-entry schema-2 review `24ede5a4...1727`,
  252,077-entry exact quality `6d8104b4...f69`, and 30,255-entry regression
  `00b736b5...3c5`.
- Preserve the btree correction: its regression manifest retains the recorded
  digest and 24,197 entries, but 163 mutable cache/HOME entries no longer
  verify after later enumeration. Its stable regression summary is unchanged;
  the other three btree manifests verify, and fresh go-cmp evidence supersedes
  current regression claims.
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
`github.com/lucasb-eyer/go-colorful v1.2.0` as exactly one bounded P7 module
group. Resolve the canonical latest release and every potentially compatible
release from fresh Go proxy, checksum-database, upstream, and primary Go
vulnerability evidence. Do not assume its latest version, Go floor, source
identity, closure, loaded status, or behavior.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused behavior when loaded, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata. Do
not manufacture a dependency commit if v1.2.0 is already canonical latest.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
