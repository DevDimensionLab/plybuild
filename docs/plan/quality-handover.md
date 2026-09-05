# Quality Upgrade Handover

Generated: 2026-09-05T17:05:19+02:00

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
- The incoming x/image documentation session is commit
  `9e9ef26ff356f9c1c2d051992eb92a10948ec807`, exact parent
  `3a8cadc24758c7a24fa6a429297d26ca29c311a9`, tree
  `1c0f03d50eb215f8224acbaaae304366261a5c03`. This handoff must be its direct
  child. The next dependency implementation, if any, must use the resulting
  documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered x/image archive and sole NEXT `golang.org/x/net` archive link
  reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, and rejected x/image v0.16.0. Go-cmp v0.7.0, Viper
v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1,
latest gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap,
properties, mapstructure, go-toml v1, afero, jwalterweatherman, gotenv, and
x/image changes remain rejected or unnecessary for their recorded floor,
loaded-behavior, release-qualification, self-test, exact-latest, or closure
decisions. Do not revisit them. P8 remains queued.

## Rejected Golang X Image Group

Fresh Go proxy, checksum-database, authoritative Go Git, GitHub mirror, and
primary Go vulnerability evidence enumerate exactly 45 stable semantic
releases, v0.1.0 through v0.45.0, with no prereleases, retractions, formal
deprecation marker, fork, archived or disabled repository, later major module
path, or GitHub Release objects. Proxy `@latest` and exact `go list` both
resolve v0.45.0 at commit
`3ebddc7c54bd879f8d84d11db82892726f5192fd`, published
2026-08-11T17:56:44Z and declaring Go 1.25.0. One later master commit
`09b0b4f7d91066e434edc5239b0d2268560dc80d` maps to pseudo-version
`v0.45.1-0.20260819171614-09b0b4f7d910`, declares Go 1.26.0, and is
unreleased.

V0.1.0-v0.12.0 declare Go 1.12, v0.13.0-v0.24.0 declare Go 1.18,
v0.25.0-v0.30.0 declare Go 1.23.0, v0.31.0-v0.36.0 declare Go 1.24.0, and
v0.37.0-v0.45.0 declare Go 1.25.0. The direct declaration alone is not the
floor decision: v0.17.0 is the first release whose x/text requirement selects
x/tools pseudo-version `v0.21.1-0.20240508182429-e35e4ccd0d2d`, which
declares Go 1.19. Every later Go-1.18 x/image release through v0.24.0 retains
that above-floor closure. V0.16.0 requires x/text v0.15.0, whose x/tools
v0.6.0 and x/mod v0.8.0 requirements declare Go 1.18. V0.16.0 is therefore
the highest stable candidate whose complete changed closure preserves the
retained Go 1.18 floor.

Candidate v0.16.0 is unsigned lightweight tag commit
`55c4ab6bd625a2e8433671ec9f9b6c46daddf2cf` at
2024-05-05T12:58:27Z. Its checksum pair is
`h1:9kloLAKhUufZhA12l5fwnx2NZW39/we1UhBesW433jw=` /
`h1:ugSZItdV4nOxyqp56HmXwH0Ry0nBCpjnZdpDaIHdoPs=`. All 253 proxy files
match the exact authoritative tag; ZIP SHA-256 is
`2ccf3619375bf633937c1b72ccd9a426953aead2ec9b2ccb0ad39e089ea6ac8a`.
Selected v0.5.0 is unsigned lightweight tag commit
`e6c2a4cdd539b91fd11131f9eecf9bb5087ab55f` at
2023-02-14T17:44:59Z and declares Go 1.12. Its checksum pair is
`h1:5JMiNunQeQw++mMOz48/ISeNu3Iweh/JaZU8ZLqHRrI=` /
`h1:FVC7BI/5Ym8R25iw5OLsgshdUBbT1h5jZTpA+mvAdZ4=`; all 250 proxy files
match and ZIP SHA-256 is
`300661d9c1e114914d6f70b14b50768ce3e93dbbede56cd801f611ff4edda220`.
The Go repository and GitHub mirror have identical master and all 45 tag refs.
Every tag is a lightweight commit ref without an independent tag signature;
local Git and GitHub both report relevant commits unsigned.

Exact projected `go get golang.org/x/image@v0.16.0` changes x/image v0.5.0 ->
v0.16.0 and x/text v0.7.0 -> v0.15.0 in go.mod and adds exactly those two
checksum pairs, four lines. It selects x/mod v0.8.0 and x/tools v0.6.0;
existing x/sys v0.30.0 dominates v0.5.0. Because x/tools is unloaded and its
Go-1.18 graph is pruned, removal of old x/tools v0.1.12 requirements exposes
the retained legacy-graph goldmark v1.3.5 selection instead of v1.4.13.
Modules remain 234, complete packages remain byte-identical at 429, graph
edges change 3,564 -> 3,548 through 22 removals and six additions, and the
unapplied tidy projection changes 332 -> 381 lines. No module or package is
added. This is the explained minimal MVS selection/edge/checksum closure.

Eight x/image packages load: bmp, ccitt, tiff/lzw, tiff, riff, vp8, vp8l, and
webp. The path is `plybuild/cmd -> go-term-markdown -> pixterm/pkg/ansimage ->
x/image/bmp`; ansimage blank-imports bmp, tiff, and webp, imaging imports bmp
and tiff, and go-term-markdown calls `ansimage.NewScaledFromReader` and
`Render`. All eight loaded decoder packages pass count-10 in both projections
with normalized-identical output. The go-term-markdown image fixture produces
normalized-identical output in both states while reproducing the same
pre-existing expected-text mismatch from the retained gomarkdown fallback, so
there is no x/image behavior delta.

Candidate module verification passes and its 253-file source remains
unchanged, but mandatory complete tests fail under exact Go 1.26.7 at count 1,
count 10, and race. `draw.TestScaleDown` rejects four checked-in golden images
because Go 1.26 changed image/jpeg decoding. Vet also fails with 80 unkeyed
`basicfont.Range` diagnostics in generated Inconsolata sources. Selected
v0.5.0 has the same failure classes. Upstream commit
`261d2777537237975bcb057e6209a6109a002288` explicitly accepts Go 1.26 JPEG
output but first ships in v0.32.0; commit
`9032ff7c7b86f42b9bebdf6133191648224aecc0` clears vet but first ships in
v0.34.0. Both releases declare Go 1.24.0 and cannot repair v0.16.0 within the
retained floor.

Govulncheck v1.7.0 measures a projected improvement from exact Darwin-symbol/
Darwin-module/Windows-symbol populations 20/30/20 to 18/28/18, removing
GO-2023-1989 and GO-2023-1990. The fresh 1,392-entry primary index, last
modified 2026-09-02T19:32:21Z, contains 13 x/image records. The vulnerability
improvement cannot override mandatory module test and vet failures. The stop
rule therefore rejects v0.16.0, retains v0.5.0, forbids dependency metadata
edits, and makes repository-wide candidate and post-implementation quality
gates inapplicable. No dependency or implementation commit was made.

## X Image Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-x-image-selection.9e9ef26.9J9x1F`. Its fully verified
  52,533-entry manifest SHA-256 is
  `84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`;
  decision-summary SHA-256 is
  `2b70234a63a199ab6cf11050212645162a78c4836aa9aafa4cb00dad764eb936`.
- The manifest covers fresh proxy/sumdb/Git/GitHub history, all 45 releases and
  declarations, upstream/mirror refs, exact source identity, release and
  signature qualification, closure boundary probes, exact-get and tidy diffs,
  selection/edge/package/loaded populations, real consumer paths, focused
  behavior, selected and candidate module tests, source-mutation checks, and
  all vulnerability scans. Every entry was recomputed successfully after
  sealing.

## Inherited Evidence And Tool Identity

- All 27 complete inherited manifest roots were independently verified at
  session start, including the immediate gotenv and jwalterweatherman roots
  and all cast, go-toml/v2, go-toml v1, mapstructure, runewidth, properties,
  mousetrap, HCL, go-colorful, btree, and toolchain roots recorded previously.
- Gotenv rejection evidence remains fully verified at
  `/private/tmp/ply-p7-gotenv-selection.3a8cadc.Dai0IG`,
  48,662/`ab2bd356b863d529085cff9b918d004aa2a03a121a7ea360a66767a2820cc88c`;
  decision summary is
  `014c050a3ccefd6fa501f4c5d04aded2ee93d02ad4039a9935d6f4d478c21e11`.
- Jwalterweatherman no-change evidence remains fully verified at
  `/private/tmp/ply-p7-jwalterweatherman-selection.2dd8fdf.h4EKeQ`,
  76,374/`54040aefec38df7415f56da66af04d48c8256dd46ae6de8e06a7c4bbe5ce4d8e`;
  decision summary is
  `342330b73509caf414a406dda3a14c026bfc59512d09bf6393f710127974dc52`.
- Cast selection/review/exact-quality/regression evidence remains fully
  verified at 82,586/12/236,782/607 entries with manifest SHA-256 values
  `dfa699f75d0b222d7d64e8fd7a09d05306de98ef4ac5cbf6d253fd33b2dddac0`,
  `850d40f319d13024275b35909e8ff84ce2f9c5ce5959cf8d707461a3b5e5d9ce`,
  `b3d410a70d89e2176daf1a013f97d53e704ef3887f02d67bd6c2564ab0a6d02a`,
  and `97cafb6719f470ab7d057a603b522f0d10418deda2d997e6d5723883efd2e24a`.
  Preserve its NUL-delimited manifest correction for whitespace paths.
- Afero rejection evidence remains fully verified at
  `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`,
  79,866/`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`.
- Preserve the go-colorful correction: its manifest retains 47,633 entries and
  SHA-256
  `a47facc21231accc5c2ed73980c9e3a7542538c6a41d7778de5cfbb410312adc`;
  47,632 stable entries verify and exactly one mutable telemetry counter still
  differs. Preserve the btree correction: its 24,197-entry regression manifest
  retains SHA-256
  `9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`
  and exactly 163 previously recorded mutable cache/HOME mismatches still
  reproduce; later evidence supersedes them.
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

Independently evaluate selected indirect `golang.org/x/net v0.7.0` as exactly
one bounded P7 module group. Resolve canonical latest, every stable release
that could satisfy the retained Go 1.18 floor through its complete changed
closure, source identity, release qualification, requirements, loaded package
population and behavior, module self-tests, and primary vulnerability data
from fresh evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths and symbols, focused
consumer behavior, repository quality, help/API/CLI identity, and exact Darwin
and Windows vulnerability populations. Implement only an exact floor-compatible
release with explained minimal closure and every gate passing; otherwise record
rejection or no-change without dependency metadata edits. Do not combine
x/image, gotenv, jwalterweatherman, cast, afero, another dependency group, or
P8.
