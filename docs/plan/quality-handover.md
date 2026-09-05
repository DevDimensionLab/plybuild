# Quality Upgrade Handover

Generated: 2026-09-05T17:58:06+02:00

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
- The incoming x/net documentation session is commit
  `e5ea7ae4ddb5bc4d58085ff3fcfe5b472e6951ce`, exact parent
  `9e9ef26ff356f9c1c2d051992eb92a10948ec807`, tree
  `897b89bbf7a8094272e58a13f73c2c174cf844e3`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered x/net archive and sole NEXT `golang.org/x/text` archive link
  reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, and x/net v0.25.0. Go-cmp v0.7.0,
Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0,
fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, jwalterweatherman, gotenv, x/image, and x/net
changes remain rejected or unnecessary for their recorded floor,
loaded-behavior, release-qualification, self-test, exact-latest, or closure
decisions. Do not revisit them. P8 remains queued.

## Rejected Golang X Net Group

Fresh Go proxy, checksum-database, authoritative Go Git, GitHub mirror, and
primary Go vulnerability evidence enumerate exactly 58 stable semantic
releases, v0.1.0 through v0.58.0, with no prereleases, retractions, formal
deprecation marker, fork, archived or disabled repository, nested or later
major module path, or GitHub Release objects. Proxy `@latest` and exact
`go list` resolve v0.58.0 at commit
`acc78e0d2b2c855c0c4fbdcfe5f42a9e3d0f9778`, published
2026-08-12T17:41:32Z and declaring Go 1.25.0. Master commit
`c23af1b9b8cc40d792e7ffd20aa02bfb8984a1f5` is 22 commits later, maps to
pseudo-version `v0.58.1-0.20260904215552-c23af1b9b8cc`, declares Go 1.26.0,
and is unreleased.

V0.1.0-v0.17.0 declare Go 1.17, v0.18.0-v0.35.0 declare Go 1.18,
v0.36.0-v0.43.0 declare Go 1.23.0, v0.44.0-v0.50.0 declare Go 1.24.0, and
v0.51.0-v0.58.0 declare Go 1.25.0. The direct declaration alone is not the
floor decision: v0.26.0 is the first release whose x/text requirement selects
x/tools pseudo-version `v0.21.1-0.20240508182429-e35e4ccd0d2d`, which
declares Go 1.19. Every later Go-1.18 x/net release through v0.35.0 retains
that above-floor closure. V0.25.0 instead requires x/text v0.15.0, whose
x/tools v0.6.0 and x/mod v0.8.0 requirements declare Go 1.18 or lower.
V0.25.0 is therefore the highest stable candidate whose complete changed
closure preserves the retained Go 1.18 floor.

Candidate v0.25.0 is unsigned lightweight tag commit
`d27919b57fa8dd03198f85ca9e675e1a09babd7d` at
2024-05-06T16:24:48Z. Its checksum pair is
`h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=` /
`h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=`. All 778 proxy files
match the exact authoritative tag; ZIP SHA-256 is
`7fd8464681c3011736f2c75beb20f88fff553a17f4f574325bce5ca5dc1fcf83`.
Selected v0.7.0 is unsigned lightweight tag commit
`8e2b117aee74f6b86c207a808b0255de45c0a18a` at
2023-02-14T17:04:22Z and declares Go 1.17. Its checksum pair is
`h1:rJrUqqhjsgNp7KqAIc25s9pZnjU7TUcSY7HcVZjdn1g=` /
`h1:2Tu9+aMcznHK/AK1HMvgo6xiTLG5rD5rZLDS+rp2Bjs=`; all 667 proxy files
match and ZIP SHA-256 is
`060552064526a90ac9b0bdce2ba0ab34592decc9428d1441c3c9722f853cd290`.
The Go repository and GitHub mirror have identical master and all 58 tag refs.
Every tag is a lightweight commit ref without an independent tag signature;
local Git and GitHub both report relevant commits unsigned.

Exact projected `go get golang.org/x/net@v0.25.0` changes x/net v0.7.0 ->
v0.25.0 and x/text v0.7.0 -> v0.15.0 in go.mod and adds exactly those two
checksum pairs, four lines. MVS also moves x/crypto to v0.23.0, x/mod to
v0.8.0, and x/tools to v0.6.0; existing x/sys v0.30.0 and x/term v0.29.0
dominate the candidate's lower requirements. Every changed selection declares
Go 1.18 or lower. Modules stay 234, complete packages stay byte-identical at
429, graph edges change 3,564 -> 3,568 through five removals and nine
additions, and the unapplied tidy projection changes 332 -> 335 lines. This is
the explained minimal MVS selection, edge, checksum, and metadata closure.

Two x/net packages load: html and html/atom. The path is `plybuild/cmd ->
go-term-markdown -> x/net/html`; go-term-markdown calls `html.Parse` and
walks `html.Node` values using node types and attributes. Both loaded
packages and the Ply markdown rendering contract pass at count 10 in both
projections with normalized-identical output. The upstream consumer's inline
HTML test reproduces the same pre-existing ANSI expected-text mismatch in both
states; it exposes no x/net behavior delta.

Candidate module verification and vet pass and its 778-file source remains
unchanged, but mandatory complete tests fail under exact Go 1.26.7 at count 1,
count 10, and race. `route.TestRouteMessage` cannot create an AF_ROUTE raw
socket in the measured Darwin sandbox and reports `operation not permitted`;
count 10 produces exactly ten such diagnostics. Selected v0.7.0 has the same
route failure, and the test's fatal socket behavior remains unchanged on
unreleased master. The explicit self-test stop rule still rejects v0.25.0,
retains v0.7.0, forbids dependency metadata edits, and makes repository-wide
candidate and post-implementation quality gates inapplicable. No dependency or
implementation commit was made.

Govulncheck v1.7.0 preserves exact Darwin and Windows symbol populations at
20/20, including the same nine symbol-reachable x/net IDs and traces. The
Darwin module population improves 30 -> 27 by removing GO-2023-1988,
GO-2023-2102, and GO-2024-2687. The fresh 1,392-entry primary index, last
modified 2026-09-02T19:32:21Z, contains 30 x/net records. That module-only
improvement cannot override mandatory module test failure.

## X Net Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-x-net-selection.e5ea7ae.cP57X8`. Its fully verified
  71,887-entry manifest SHA-256 is
  `88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`;
  decision-summary SHA-256 is
  `6e32698abcde67604e159488ba21f6e780fd61dfba2ed9299644b5f011940dff`.
- The manifest covers fresh proxy/sumdb/Git/GitHub history, all 58 releases and
  declarations, upstream/mirror refs, exact source identity, release and
  signature qualification, complete closure boundary probes, exact-get and
  tidy diffs, selection/edge/package/loaded populations, consumer paths and
  focused behavior, selected and candidate module tests, source-mutation
  checks, and all vulnerability scans. Every entry was recomputed
  successfully after sealing.

## Inherited Evidence And Tool Identity

- All 27 complete inherited manifest roots were independently verified at
  session start, including the immediate x/image, gotenv, jwalterweatherman,
  cast, go-toml/v2, go-toml v1, mapstructure, runewidth, properties,
  mousetrap, HCL, go-colorful, btree, and toolchain roots recorded previously.
- X/image rejection evidence remains fully verified at
  `/private/tmp/ply-p7-x-image-selection.9e9ef26.9J9x1F`,
  52,533/`84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`;
  decision summary is
  `2b70234a63a199ab6cf11050212645162a78c4836aa9aafa4cb00dad764eb936`.
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

Independently evaluate selected indirect `golang.org/x/text v0.7.0` as
exactly one bounded P7 module group. Resolve canonical latest, every stable
release that could satisfy the retained Go 1.18 floor through its complete
changed closure, source identity, release qualification, requirements, loaded
package population and behavior, module self-tests, and primary vulnerability
data from fresh evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths and symbols, focused
consumer behavior, repository quality, help/API/CLI identity, and exact Darwin
and Windows vulnerability populations. Implement only an exact
floor-compatible release with explained minimal closure and every gate
passing; otherwise record rejection or no-change without dependency metadata
edits. Do not combine x/net, x/image, gotenv, jwalterweatherman, cast, afero,
another dependency group, or P8.
