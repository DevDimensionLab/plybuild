# Quality Upgrade Handover

Generated: 2026-09-08T15:30:02+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains exact XXHash v2.3.0 commit
  `e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
  `a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
  `5213ba55981d77d7c8061312915e237e80d29af8`. It changes only `go.mod`
  and `go.sum` with three insertions.
- Imaging, Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt were retained
  without implementation commits. Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor, as do
  Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`.
- The answered Imaging archive and sole NEXT ansimage archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, and Imaging selections. P8 remains queued. All
earlier acceptances, rejections, no-change decisions, evidence corrections,
and lifecycle ancestry are final; do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve the launcher's
scratch cleanup and reciprocal archive contract.

## Retained Imaging Group

Retain exact-path `github.com/disintegration/imaging v1.6.2` without editing
dependency metadata. A fresh proxy list contains exactly 15 stable versions,
v1.0.0 through v1.6.2, and exact `@latest` is selected v1.6.2 at
2019-11-16T20:43:25Z. Its module declares the exact path, no Go directive, and
`golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8`; it has no
deprecation or retraction. Sumdb source/module hashes are
`h1:w1LecBlG2Lnp8B3jk5zSuNqd7b4DXhcjwek1ei82L+c=` and
`h1:44/5580QXChDfwIclfc/PCwrr44amcmDAg8hxG0Ewe4=`.

All 15 proxy ZIPs match their exact Git tags. Selected ZIP SHA-256 is
`2934e7bace3c8c0b1b4a07144197e8720b9ffbe922600e3a3c764f77792ac7c4`;
the normalized 61-file source-manifest SHA-256 is
`71fbee7ec7a3b983f75bd27e1cc91d95c67096e01492897f5f715f8644965d51`.
Go-import metadata resolves to the public, enabled, unarchived, non-fork
GitHub repository on `master`.

Stable GitHub Release and lightweight tag v1.6.2 resolve to commit
`acabd8315e63bfcaac97d52d68a7a0b88d2eea93`, tree
`6584cbb2a26e4d38bfec8f2587633f234500810e`, with parents
`9aab30e6aa535fe3337b489b76759ef97dfaf362` and
`675e3c209ff3e9bbee22db0bffe990d3abace4ce`. The tag has no tag signature;
the commit's embedded GitHub web-flow signature independently verifies with
fingerprint `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`. Later master
`d40f48ce0f098c53ab1fcd6e0e402da682262da5`, tree
`cfae2d9af62546482831388535eb24023ceac298`, resolves only as unreleased
`v1.6.3-0.20201218193011-d40f48ce0f09`; the later Dependabot branch is also
an unreleased pseudo-version. Exact v1.6.3, v1.6.3-rc.1, and `/v2` do not
exist. `gopkg.in/disintegration/imaging.v1@v1.5.0` is an alternate path. No
higher exact stable release exists.

The complete standalone graph is Imaging, its x/image pseudo-version (Go
1.12), and directive-free x/text v0.3.0. Loaded tests use Imaging plus
x/image/bmp, ccitt, tiff/lzw, and tiff; x/text is graph-only. Pristine proxy
and exact-Git source pass verify, listing, count-1, two count-10 passes, race,
and vet under exact Go 1.26.7 and contained Go 1.18.10. Windows/amd64,
Linux/amd64, Linux/arm64, FreeBSD/amd64, and js/wasm test builds pass under
both SDKs, proving the complete Go 1.18 floor.

Imaging is pure Go with no Cgo, assembly, generated files, build constraints,
or platform-specific implementation. Its suite exposes 59 tests and 27
benchmarks; one documentation example has no `Output` assertion and is not an
executable example test. There are no fuzz or property tests. Native goldens
cover exported I/O and transforms, five formats, JPEG EXIF orientations 0-8,
filters, crop/composite, adjustments, and convolution. Independent fixtures
add all five round trips, malformed and huge-header input, extension errors,
alpha correctness, non-zero origins, invalid dimensions, and 100 seeded
transform properties. All fixture repeat/race/vet checks pass under both SDKs.
Imaging has no configurable decoded-pixel, byte, or reader limit; callers must
bound input.

Ply loads Imaging in production through `plybuild/cmd ->
go-term-markdown -> pixterm/pkg/ansimage -> Imaging`. Markdown passes local or
HTTP readers to ansimage; ansimage decodes through registered formats and
calls Imaging `Resize`, `Fit`, or `Fill` with Lanczos. Project MVS selects
x/image v0.5.0. Keep that selection, Imaging's older declared edge, the root
indirect requirement, and loaded production behavior distinct. This group did
not reopen x/image.

Ansimage has no native `_test.go` files, so its native commands establish only
compilation. The historical Markdown suite passes under both SDKs with a PTY
and `NO_COLOR` unset; it has expected terminal-dependent ANSI golden failures
otherwise and unrelated project-MVS Chroma golden drift. Independent
41-module consumer fixtures have a highest Go 1.18 declaration and pass
count-1, two count-10 passes, race, vet, and Windows compile under both SDKs.
They exercise the exact production path for JPEG/PNG/GIF/TIFF/BMP, malformed
input, and all three scale modes. They preserve one selected ansimage defect:
no-dither `RenderExt` emits nothing for a two-pixel-high scaled image while
Markdown reports it rendered; four-pixel fixtures render normally. This is
not an Imaging defect or a basis for changing Imaging.

## Projection, Quality, And Vulnerability Measurements

Accepted state remains 234 selected modules, 3,583 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 33 checksum lines and
removes zero. Main-module Go remains 1.18 and toolchain remains 1.26.7.

Exact selected Imaging `go get` is entirely inert: it changes no requirement,
selection, edge, checksum, tidy projection, or status. It was not applied.
Non-mutating `go mod tidy -diff` reproduces the accepted 381-line unrelated
cleanup while retaining Imaging and x/image selections.

The unchanged project passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical root/status/upgrade/build
help, API/CLI compatibility and reports, and empty-HOME count-2. API and CLI
report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
No changed stable selection exists, so the changed-selection-only full P7 gate
was not invoked; the accepted XXHash quality result remains authoritative.

Fresh vulnerability data has 1,392 module records and no Imaging record.
Canonicalized current/no-op sets are identical: 20 IDs/22 reachable traces in
Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin
module IDs/findings. No Imaging module, package, symbol, or trace appears.

Imaging evidence has 1,713 verified entries; manifest SHA-256 is
`046746e0c4004d62ebac4838dac739ce37a0d4576a0fae3e5d1db987e4d47308`.
Decision-summary SHA-256 is
`507f403c1289ff6d698beffb31eea6c3a3c07835efcf5609d6bde475bb4cdc5e`.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive SHA-256 remains
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  Retained GoReleaser and apidiff receipts remain final; a scratch rebuild of
  apidiff was used only for verified zero-diff function.
- Retain accepted XXHash v2.3.0 and Speakeasy v0.2.0 moves and every earlier
  exact decision. Answered archives remain authoritative detail.
- Authoritative accepted Q0-Q2 scorecard SHA-256 remains
  `579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

## Next Objective

Independently evaluate selected exact-path
`github.com/eliukblau/pixterm/pkg/ansimage
v0.0.0-20191210081756-9fb6cf8c2f75` as the next single P7 group. The exact
proxy list is empty and exact `@latest` returns 404, while the selected
pseudo-version remains fetchable. It declares the exact nested path, Go 1.13,
Imaging v1.6.2, go-colorful v1.0.3, and an x/image pseudo-version.

Selected unsigned commit `9fb6cf8c2f75275ebcd4ac8c30a0e26930d497f7`, tree
`bdcdecab7b23ba3a6d18efd3ced9800826668031`, is also root-project Release/tag
v1.3.0, but the unprefixed root tag does not version the nested module. It is
the only commit containing `pkg/ansimage/go.mod`. Immediate successor
`9f095995d66abbf03a06cea8e8e9b7cfd679e06c` consolidates to a root module.
Root releases continue through v1.3.3 at master
`24a1aedad1a99b2177808bfd72b23e27318989f6`, whose root go.mod declares Go
1.25.0. Prove exact nested-module identity rather than treating later root
tags or source directories as releases of the selected path.

Preserve the actual Markdown production chain and characterize constructors,
formats, scaling/dimension logic, dithering, ANSI colors, terminal sizing,
alpha/background, malformed and small images, and the two-pixel no-dither gap.
Do not combine Imaging, go-colorful, x/image, terminal behavior, root-module
migration, or unrelated tidy cleanup.
