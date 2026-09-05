# Quality Upgrade Handover

Generated: 2026-09-05T18:42:49+02:00

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
- The incoming x/text documentation session is commit
  `6ab5945de5fef2ee6f7deff10f4bc51e94596b30`, exact parent
  `e5ea7ae4ddb5bc4d58085ff3fcfe5b472e6951ce`, tree
  `0b6b23569b9296c50085a084737f74debaf95e62`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered x/text archive and sole NEXT `gopkg.in/yaml.v2` archive must
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, and x/text v0.15.0.
Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color
v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, jwalterweatherman, gotenv, x/image, x/net, and
x/text changes remain rejected or unnecessary for their recorded floor,
loaded-behavior, release-qualification, self-test, exact-latest, or closure
decisions. Do not revisit them. P8 remains queued.

## Rejected Golang X Text Group

Fresh Go proxy, checksum-database, authoritative Go Git, GitHub mirror, and
primary Go vulnerability evidence enumerate exactly 49 stable semantic
releases, v0.1.0 through v0.41.0, with no prereleases, retractions, formal
deprecation marker, fork, archived or disabled repository, or nested/later
major module path. Proxy `@latest` and exact `go list` resolve v0.41.0 at
commit `acdba6655fd45cdb5ab73c9d6a8981333bd65a39`, published
2026-08-11T15:22:47Z and declaring Go 1.25.0. The three historical GitHub
Release objects are stable; later releases are proxy-listed authoritative
tags. Master `f53c31601f90c1b840c0703c5537c3ce1e4b6f5c` is 12 commits later,
declares Go 1.26.0, and is unreleased.

V0.1.0-v0.3.2 have no Go declaration, v0.3.3-v0.3.6 declare Go 1.11,
v0.3.7-v0.13.0 declare Go 1.17, v0.14.0-v0.22.0 declare Go 1.18,
v0.23.0-v0.28.0 declare Go 1.23.0, v0.29.0-v0.34.0 declare Go 1.24.0, and
v0.35.0-v0.41.0 declare Go 1.25.0. Direct declarations are insufficient:
v0.16.0 first requires x/tools pseudo-version
`v0.21.1-0.20240508182429-e35e4ccd0d2d`, which declares Go 1.19, and that
requirement remains through v0.22.0. V0.15.0 instead requires x/tools v0.6.0,
x/mod v0.8.0, and x/sys v0.5.0, all at Go 1.18 or lower. V0.15.0 is therefore
the highest stable complete-closure-compatible candidate.

Candidate v0.15.0 is unsigned lightweight tag commit
`8d533a0c40adec778a7d09ac6c8aa640d3c883f4` at
2024-04-15T18:14:38Z. Its checksum pair is
`h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=` /
`h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=`. All 542 proxy files
match the exact authoritative tag; ZIP SHA-256 is
`13faee7e46c8a18c8a28f3eceebf15db6d724b9a108c3c0482a6d2e58ba73a73`.
Selected v0.7.0 is unsigned lightweight tag commit
`71a9c9afc4cd710b9412f7f99f0d8e35b10e488a` at
2023-01-31T16:01:06Z. Its checksum pair is
`h1:4BRB4x83lYWy72KwLD/qYDuTu7q9PjSagHvijDw7cLo=` /
`h1:mrYo+phRRbMaCq/xk9113O4dZlRixOauAjOtrjsXDZ8=`; all 530 proxy files
match, with ZIP SHA-256
`4d017493c58addadf3c753056b921b47ae386a4cfd10eab2d90ed1252c6ba0e4`.
The authoritative repository and GitHub mirror have identical master and all
tag refs; relevant tags are lightweight and relevant commits are unsigned.

Exact projected `go get golang.org/x/text@v0.15.0` changes only x/text in
go.mod and adds exactly its checksum pair. MVS also moves x/mod from
`v0.6.0-dev.0.20220419223038-86c51ed26bb4` to v0.8.0 and x/tools v0.1.12 to
v0.6.0; x/sys v0.30.0 dominates the candidate's v0.5.0 requirement and x/sync
remains v0.1.0. Modules stay 234, complete packages stay byte-identical at
429, graph edges change 3,564 -> 3,567 through one removal and four additions,
and the unapplied tidy projection changes 332 -> 333 lines. This is the
explained minimal closure.

Three x/text packages load: runes, transform, and unicode/norm. The real path
is `plybuild/cmd -> spf13/viper -> spf13/afero -> x/text/runes`; Afero uses
`transform.Chain`, normalization forms, rune removal, and `transform.String`.
Loaded-package, Afero, and Ply cmd focused tests pass at count 10 in both
projections with normalized-identical output.

Candidate source remains byte-identical and verification passes, but complete
tests at count 1/count 10, race, and vet all fail under exact Go 1.26.7. Three
stale Example identifiers fail vet, and message/pipeline panics in x/tools
v0.6.0's SSA builder on a Go 1.26 range-over-function construct; repeated
tests also expose `cases.TestShortBuffersAndOverflow` persistence failures.
Standalone vet additionally reports 20 unkeyed literals, four unreachable
statements, and one unused `currency.Unit.String` result. Selected v0.7.0
reproduces the failure classes. Example fixes first ship in v0.18.0, whose
closure already declares Go 1.19. The mandatory module stop rule rejects
v0.15.0, retains v0.7.0, forbids metadata edits, and makes downstream
repository gates inapplicable. No implementation commit was made.

Govulncheck v1.7.0 preserves exact Darwin-symbol/Darwin-module/Windows-symbol
populations at 20/30/20. Both states have the same sole x/text module finding,
GO-2026-5970, fixed in v0.39.0, and neither has an x/text symbol-reachable
finding or trace. The fresh 1,392-entry primary index contains four x/text
records; the other three were fixed before v0.7.0.

## X Text Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-x-text-selection.6ab5945.3qn5Bj`. Its fully verified
  82,629-entry manifest SHA-256 is
  `48c872da4ab796ef1115003fcf0a226a07bf5791a163d36bc9c946c9d9f92ed6`;
  decision-summary SHA-256 is
  `50d6bee90ea7cea060e23400e820b1537742ec2167bee04f188070e2e5eb39c9`.
- The manifest covers fresh proxy/sumdb/Git/GitHub history, all 49 releases and
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

Independently evaluate selected indirect `gopkg.in/yaml.v2 v2.4.0` as
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
edits. Do not combine x/text, x/net, x/image, gotenv, jwalterweatherman, cast,
afero, yaml.v3, another dependency group, or P8.
