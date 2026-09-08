# Quality Upgrade Handover

Generated: 2026-09-08T17:07:36+02:00

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
- Ansimage, Imaging, Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt were
  retained without implementation commits. Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor, as do
  Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`.
- The answered ansimage archive and sole NEXT Fatih Color archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, Imaging, and ansimage selections. P8 remains queued.
All earlier acceptances, rejections, no-change decisions, evidence corrections,
and lifecycle ancestry are final; do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*`
roots, never run `go mod download all` in a measured tree, and preserve the
launcher's scratch cleanup and reciprocal archive contract.

## Retained Ansimage Group

Retain exact nested module
`github.com/eliukblau/pixterm/pkg/ansimage
v0.0.0-20191210081756-9fb6cf8c2f75` without editing dependency metadata. Its
exact proxy list is empty; exact `@latest`, `@v0`, and `@v1` return 404
with “no matching versions.” The selected pseudo-version remains fetchable
with source/module hashes
`h1:vbix8DDQ/rfatfFr/8cf/sJfIL69i4BcZfjrVOxsMqk=` and
`h1:0gZuvTO1ikSA5LtTI6E13LEOdWQNjIo5MTQOvrV0eFg=`. It declares the exact
path, Go 1.13, Imaging v1.6.2, go-colorful v1.0.3, and x/image
`v0.0.0-20191206065243-da761ea9ff43`; there is no deprecation or retraction.

Selected unsigned commit
`9fb6cf8c2f75275ebcd4ac8c30a0e26930d497f7`, tree
`bdcdecab7b23ba3a6d18efd3ced9800826668031`, parent
`be34e524a7d8fbf6ab827fd06f669ad4e50943b0`, dated
2019-12-10T08:17:56Z, is also the target of root-project stable Release and
unsigned annotated tag v1.3.0. Because that tag is not prefixed
`pkg/ansimage/`, it does not version the nested module; the exact path
correctly resolves the commit only as the selected pseudo-version.

The incoming statement that selected was the only repository commit containing
`pkg/ansimage/go.mod` was false. It is the only commit that adds the file.
Six later side-branch or merge commits still contain it and are fetchable as
exact nested pseudo-versions:

- `v0.0.0-20191216152442-66ce7d8b90d7`, which changes five HTTP URLs in
  comments to HTTPS and no behavior;
- `v0.0.0-20191216153813-7fa2a5d29053`;
- `v0.0.0-20191216175341-64e5dd02853a`;
- `v0.0.0-20191221043440-fdd01950d74d`;
- `v0.0.0-20191221043740-0b8ce46c8e11`; and
- latest fetchable `v0.0.0-20191221044037-630511e42559`.

The latter five preserve the selected ansimage subtree byte-for-byte. All six
are unreleased pseudo-versions, not higher exact stable releases. Immediate
successor `9f095995d66abbf03a06cea8e8e9b7cfd679e06c`, tree
`ed6919823270b5ac3737c516ea8900c20bb4911e`, removes the nested module in
favor of repository-root module `github.com/eliukblau/pixterm`. Later root
releases v1.3.1-v1.3.3, current master
`24a1aedad1a99b2177808bfd72b23e27318989f6`, its Go 1.25.0 root module,
branches, forks, alternate paths, and redirects do not provide a stable release
of the exact nested identity. Root v1.3.0-v1.3.2 tag objects and commits are
unsigned; v1.3.3 is a lightweight unsigned tag. The public upstream is enabled,
unarchived, non-fork, and defaults to protected `master`.

Selected proxy and exact-Git source are byte-identical. The normalized
three-file manifest SHA-256 is
`0d7ed0f139e2f4397ba2028248dd08d2cb4749a29e44ab5d6b39535f771906f0`.
The selected ZIP contains only `ansimage.go`, `go.mod`, and `go.sum`;
there is no native test, example, fuzz/property suite, testdata, generated
file, build constraint, Cgo, assembly, or OS/architecture implementation.

The complete standalone closure is five modules: ansimage, Imaging,
go-colorful, x/image, and graph-only x/text v0.3.0. Its highest declaration is
Go 1.13. Eleven non-standard packages load: ansimage, Imaging, go-colorful,
and x/image's BMP, CCITT, TIFF/LZW, TIFF, RIFF, VP8, VP8L, and WebP packages.
Pristine source and independent fixtures pass verify/list, count-1, two
count-10 passes, race, and vet under exact Go 1.26.7 and contained Go 1.18.10.
Windows/amd64, Linux/amd64, Linux/arm64, FreeBSD/amd64, and js/wasm test
compilation passes under both SDKs, proving the complete retained floor.

Independent fixtures cover exported constructors from images, readers, files,
and URLs; Resize/Fill/Fit; GIF/JPEG/PNG/BMP/TIFF/WebP; getters/setters;
alpha/background compositing; malformed and small images; non-zero bounds;
Go-code/background-disabled rendering; Draw/ClearTerminal; and URL status
handling. Ansimage emits ANSI 24-bit true color only, has no 256-color mode,
and does not discover terminal size. Ply passes the line width through
go-term-markdown.

The selected no-dither render loop begins at terminal row 1. A two-pixel-high
image therefore renders empty; a four-pixel image renders only its second
pixel pair. Markdown still reports the image rendered and emits only its
title/destination in the two-pixel case. This is invariant ansimage behavior,
not Imaging behavior, and the latest fetchable pseudo-version is
source-identical. Dither rendering omits its final aggregate row; it groups 8x4
source blocks and rejects sizes that produce fewer than two aggregate rows or
columns. Positive non-zero image origins inflate dimensions by using
`Bounds.Max` rather than `Dx/Dy`; negative origins fail bounds checks.
`SetMaxProcs` accepts zero and negative values without validation. Unknown
scale and dither modes panic.

Ansimage has no configurable input-byte or decoded-pixel bound. Its own URL
constructors use the default client without a timeout and do not close non-200
bodies. Markdown uses a five-second HTTP client, but `renderImage` does not
close successful local/HTTP readers and also leaves non-200 bodies open.
Production-consumer fixtures nevertheless pass both SDKs, repeats, race, vet,
and five cross-target builds while exercising the exact local and HTTP path,
all six formats, all three dither modes, malformed input, and small images.

## Projection, Quality, And Vulnerability Measurements

Exact selected ansimage `go get` is byte-inert and was not applied. Current
and no-op states remain 234 selected modules, 3,583 graph edges, 429 native
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, `go.sum` adds exactly 33 lines and
removes zero. The main module remains Go 1.18 with toolchain Go 1.26.7.

The unchanged project passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned golangci-lint 2.12.2 with zero findings,
byte-identical root/status/upgrade/build help, API/CLI compatibility,
empty-HOME count-2, all four acceptance verifiers, preflight, 80/80 mutation
controls, and 15/15 quality-audit controls. API and CLI report SHA-256 values
remain `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
No qualified changed selection exists, so the changed-selection-only full P7
quality/snapshot/Docker gate was inapplicable; the accepted XXHash scorecard
remains authoritative.

Fresh govulncheck v1.7.0 used the 1,392-record primary database updated
2026-09-02T19:12:04Z. Current/no-op sets are identical at 20 IDs/22 reachable
traces for Darwin and Windows symbol scans, 22 Darwin package IDs/findings,
and 30 Darwin module IDs/findings. Ansimage has no direct primary record. It
is correctly present as a call frame in 12 existing traces across 11 x/image
IDs because project MVS selects x/image v0.5.0; the exact nested module's
declared 2019 x/image closure has zero findings. This session did not reopen
or change x/image.

Ansimage evidence has 73 verified entries; manifest SHA-256 is
`95d6ca0ac0d767c94f78b3a0a5c32f1008cd2e7b05a3ffc5f639b9af021dd112`.
Decision-summary SHA-256 is
`3f39b1aa5f16d3a1b398ce329e4a73621931b09d3b191dd7c0bbadc5743de0ed`.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 binary SHA-256 is
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
- Pinned golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  Retained GoReleaser and apidiff receipts remain final; the scratch apidiff
  executable was independently verified by embedded module metadata.
- Retain accepted XXHash v2.3.0 and Speakeasy v0.2.0 moves and every earlier
  exact decision. Answered archives remain authoritative detail.
- Authoritative accepted Q0-Q2 scorecard SHA-256 remains
  `579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

## Next Objective

Independently evaluate exact-path `github.com/fatih/color v1.14.1` as the
next single P7 group. A minimal survey finds 22 proxy versions and stable
`@latest` v1.19.0. Selected v1.14.1 declares Go 1.17 and requires
go-colorable v0.1.13, go-isatty v0.0.17, and x/sys v0.3.0. V1.15.0-v1.18.0
continue to declare Go 1.17; v1.18.0 requires colorable v0.1.13, isatty
v0.0.20, and x/sys v0.25.0. Latest v1.19.0 declares Go 1.25.0 and therefore
does not preserve the retained Go 1.18 floor.

The enabled, unarchived, non-fork upstream defaults to `main`. Stable
GitHub Releases continue through v1.19.0. Selected tag commit is
`3d5097c6b003cf3a784e670ddb79710cf46e9a07`; v1.18.0 is
`1c8d8706604ee5fb9a464e5097ba113101828a75`; v1.19.0 is
`ca25f6e17f118a5a259f3c2c0d395949d1103a5a`. Project MVS already selects
go-colorable v0.1.15, go-isatty v0.0.20, and x/sys v0.30.0. Keep those
already-reviewed selections distinct rather than turning the group into their
upgrade review.

Ply loads Fatih Color through `plybuild/cmd -> go-term-markdown ->
github.com/fatih/color`. Resolve release identity, complete floor closure,
actual ANSI/TTY/environment behavior, concurrency and global state, exact MVS
effects, API compatibility, consumers, and vulnerability parity. V1.18.0 is
the highest immediately visible stable release whose own directive preserves
Go 1.18, but its full closure and every contract still require independent
qualification. Do not select v1.19.0, combine terminal-support dependencies,
or apply unrelated tidy cleanup.
