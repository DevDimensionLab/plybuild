# Quality Upgrade Handover

Generated: 2026-09-05T16:14:21+02:00

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
- The gotenv documentation session is commit
  `3a8cadc24758c7a24fa6a429297d26ca29c311a9`, exact parent
  `2dd8fdfc61d54d9327effad90242f1f05060fa42`, tree
  `b6b500f83966d95fc91fa9b6c14e2ced4175cad0`. This handoff must be its direct
  child. The next dependency implementation, if any, must use the resulting
  documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Minimal closures
  select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty v0.3.1,
  and rogpeppe/go-internal v1.9.0. Those four accepted groups add exactly 15
  checksum lines.
- The answered gotenv archive and sole NEXT `golang.org/x/image` archive link
  reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
and rejected gotenv v1.6.0. Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini
v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, jwalterweatherman, and gotenv changes remain
rejected or unnecessary for their recorded floor, loaded-behavior,
release-qualification, self-test, exact-latest, or closure decisions. Do not
revisit them. P8 remains queued.

## Rejected Subosito Gotenv Group

Fresh Go proxy, checksum-database, GitHub, upstream Git, and primary Go
vulnerability evidence lists exactly eleven stable releases, v0.1.0 through
v1.6.0, with no prereleases, retractions, formal deprecation, fork, archive,
disabled repository, or later major module path. The GitHub Releases API has
no release objects, so the semantic tags provide release qualification.
Releases through v1.2.0 have no Go directive; v1.3.0 through v1.6.0 declare
Go 1.18. Proxy `@latest` and
`go list -m ...@latest` both resolve v1.6.0, so it is canonical latest stable
and the highest stable release compatible with the retained Go 1.18 floor.

V1.6.0 is annotated tag object
`b1e4a564ab83bcd77e00d5687055911884986d98`, peeled to commit
`14a05352a5cf0f66fd7cbce114374f56065891f0` at
2023-08-15T12:05:45Z; the tag timestamp is 2023-08-15T12:06:00Z. GitHub
verifies both tag and commit signatures. Local verification is unavailable
because `gpg` is absent. Its checksum pair is
`h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=` /
`h1:Dk4QP5c2W3ibzajGcXpNraDfq2IrhjMIvMSWPKKo0FU=`. All 21 proxy files match
the exact tag commit; ZIP SHA-256 is
`142db3dd2328e744c157e85cf3291d027013b79f92a45984f860fe38bc0f1f8d`
and file-manifest SHA-256 is
`a0044f44e956b351d4482487488db8e95e5e9d3210a42524a50b1274b5bfc7c9`.

Selected v1.4.2 is unsigned lightweight tag/commit
`d2e64e6317e94db9c598bb9934e75fb80c804bca` at
2023-01-11T21:45:46Z. Its checksum pair is
`h1:X1TuBLAMDFbaTAChgCBLu3DU3UPyELpnF2jjJ2cz/S8=` /
`h1:ayKnFf/c6rvx/2iiLrJUk1e6plDbT3edrFNGqEflhK0=`. All 19 proxy files match;
ZIP SHA-256 is
`5baaaaa7d88a44a5795c7a8ed8e6ffffb8d7fb27fa9c1467eca544c16136b561`.
Its file-manifest SHA-256 is
`7777e4ee03e0b63b5a69576e28fe7c5940c4681dfb39389e091691ea6d6f0610`.
The v1.4.2 -> v1.6.0 range is seven commits, nine files, +138/-20, chiefly
reader/scanner error propagation and UTF-16 BOM support. The changelog has two
headings labeled 1.5.0; the signed semantic tag and proxy determine v1.6.0's
release identity. Six commits after v1.6.0 are unreleased; current master
`d24eb16ed8bfac3d8c3ac0312b930dae9e1bd1f1` maps to a pseudo-version and
declares Go 1.22, so it is not a stable floor-compatible candidate.

The exact v1.6.0 projection changes gotenv v1.4.2 -> v1.6.0, x/text v0.7.0 ->
v0.12.0, x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 -> v0.8.0, and
x/tools v0.1.12 -> v0.6.0. Their declarations are Go 1.17 or 1.18;
candidate edges to x/sys v0.5.0, x/net v0.6.0, testify v1.7.5, goldmark
v1.4.13, and x/sync v0.1.0 lose to existing selections. Modules remain 234,
graph edges rise 3,564 -> 3,568, complete packages rise 429 -> 434 through
five x/text packages, and go.sum gains exactly four lines: the gotenv and
x/text checksum pairs. The unapplied tidy projection changes 332 -> 335
lines; the edge set removes six and adds ten. The new x/text pair is
`h1:k+n5B8goJNdU7hSvEtMUz3d1Q6D/XW4COJSJR6fN0mc=` /
`h1:TvPlkZtksWOMsz7fbANvkp4WM8x/WCo/om8BMLbz+aE=`. Sumdb also verifies
x/mod v0.8.0 as
`h1:LUYupSeNrTNCGzR/hVBk2NHZO4hXcVaW1k4Qx7rjPx8=` /
`h1:iBbtSCu2XBx23ZKBPSOrRkjjQPZFPuis4dIYUhu/chs=` and x/tools v0.6.0 as
`h1:BOw41kyTf3PuCW1pVQf8+Cyg8pMlkYB1oo9iJ6D/lKM=` /
`h1:Xwgl3UAJ/d3gWutnCtw505GrjyAbvKui8lOU390QaIU=`; those pairs are not new
go.sum lines. This is the explained minimal MVS selection/edge/checksum
closure.

Exactly one gotenv package loads through `plybuild/cmd -> spf13/viper ->
spf13/viper/internal/encoding/dotenv -> subosito/gotenv`; Viper's codec calls
`gotenv.StrictParse`. Ten-run Viper internal dotenv codec, public dotenv
read/write, and Ply `./cmd` focused behavior pass with normalized-identical
old/candidate output.

Candidate verification and vet pass, but complete module `go test ./...`
fails at count 1 and count 10, and race fails for the same reason. In
`TestScanner`, trailing LF, CR, and CRLF fixtures expect four tokens but
produce three under Go 1.22 and later, including exact Go 1.26.7. Selected
v1.4.2 normalizes to the same failure. Upstream commit
`490d1d0ccb14e51f1425e3598d6468c33ae0b1fb` fixes the scanner expectation,
but it exists only on the unreleased Go-1.22 master line. A first attempt
omitted creation of its external GOTMPDIR; the canonical replay created it
and exposed the real test failure. The mandatory candidate-module self-test
stop rule therefore rejects v1.6.0, retains v1.4.2, and prevents dependency
metadata edits or downstream candidate repository gates.

Govulncheck v1.7.0 preserves exact and identical old/candidate
Darwin-symbol/Darwin-module/Windows-symbol populations 20/30/20. The fresh
1,392-entry primary vulnerability index, last modified
2026-09-02T19:32:21Z, has no gotenv record; it has four x/text and two x/mod
records, none of which changes the exact scans. No implementation or
dependency commit was made.

## Gotenv Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-gotenv-selection.3a8cadc.Dai0IG`. Its fully verified
  48,662-entry manifest SHA-256 is
  `ab2bd356b863d529085cff9b918d004aa2a03a121a7ea360a66767a2820cc88c`;
  decision-summary SHA-256 is
  `014c050a3ccefd6fa501f4c5d04aded2ee93d02ad4039a9935d6f4d478c21e11`.
- The manifest covers fresh proxy/sumdb/Git/GitHub history, exact source
  identity, tag and commit status, all releases and declarations, both module
  projections, exact-get and tidy diffs, minimal closure, loaded population
  and paths, focused behavior, candidate module tests, and vulnerability
  scans. Its external manifest-verification log exits 0.

## Inherited Evidence And Tool Identity

- All 26 complete inherited manifest roots were independently verified at
  session start.
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
- Accepted go-toml/v2 selection/review/exact-quality/regression evidence
  remains fully verified at 17,093/8/236,599/3,743 entries with hashes
  `ce1a47d85cd2fa0ff3bde19c1953c0dc73914017145abe2f4c5a27c4cf685ff1`,
  `328672c0410e2bfffa058fdf5ec389f4c8640bd73caccde22603b0ad57c9b070`,
  `85f24a00baaf63d62bb6a3a07dc97bffcaa98b8fed1417b7d6f8425ba8d98d3d`,
  and `9601ce867fe389eb1efca59b2ae578052ba9b48670d1112cef3dc66298026ec7`.
- Go-toml v1 and mapstructure evidence remains fully verified at
  37,719/`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`
  and 46,352/`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
  Properties, mousetrap, and HCL remain fully verified at
  17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
  29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
  and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
  All four runewidth roots remain fully verified with their roadmap-recorded
  hashes.
- Preserve the go-colorful correction: 47,632 stable entries verify with one
  recorded mutable telemetry mismatch. Preserve the btree correction: its
  selection/review/exact-quality roots verify; the regression root retains
  exactly 163 recorded mutable cache/HOME mismatches superseded by later
  evidence.
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

Independently evaluate selected indirect `golang.org/x/image v0.5.0` as
exactly one bounded P7 module group. Resolve canonical latest, every stable
release that could satisfy the retained Go 1.18 floor, source identity,
release qualification, requirements, closure, loaded package population and
behavior, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

Measure exact old/candidate modules, graph edges, complete packages, checksums,
exact-get diff, unapplied tidy projection, loaded paths, focused image/render
behavior if loaded, repository quality, help/API/CLI identity, and exact
vulnerability populations. Implement only an exact floor-compatible release
with explained minimal closure and every gate passing; otherwise record
rejection or no-change without dependency metadata edits. Do not combine
gotenv, jwalterweatherman, cast, afero, another dependency group, or P8.
