# Quality Upgrade Handover

Generated: 2026-09-07T16:41:19+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest dependency implementation remains Repr commit
  `6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
  `53475076e1c79f6d2181877e3238a1a5d389c246`, tree
  `82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
  `go.sum`.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  that implementation and direct-child Repr handoff
  `342c7ece82eae3cd726aa658462089cfb876fbe0`. Preserve it.
- The answered Chroma archive and sole NEXT Colour archive link reciprocally.
  Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  still adds exactly 27 checksum lines. Chroma changed no metadata.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, external source, report, evidence tree,
and build context must remain beneath `$CODEX_SESSION_SCRATCH_ROOT`. Preserve
the launcher lifecycle controls and scratch-only mutation, snapshot, Docker,
and acceptance cleanup.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Chroma v0.10.0. All earlier accepted,
rejected, and no-change decisions and evidence corrections remain final. P8
remains queued.

## Retained Alecthomas Chroma Group

The exact `github.com/alecthomas/chroma` proxy path lists 30 stable tags from
v0.1.0 through v0.10.0, with no prerelease or retraction. Proxy `@latest` and
exact Go `@latest`/`@v0` resolve canonical-latest stable v0.10.0 at
2022-01-12T10:49:38Z. Its lightweight tag is unsigned commit
`36bdd4b98823bd1d7be96767cde3dd575e60b406`, tree
`d27058989b845352d49b45e7538f7b0004c8d651`, parent
`b01c8fcab6f3d79534dadbc8fabd4c54cabaf4d8`. GitHub has a non-draft,
non-prerelease Release object for v0.10.0; 26 of 30 v0 tags have Release
objects. The enabled, unarchived, undisabled, non-fork repository has no v0
deprecation marker.

Exact v0-path `@master` resolves the later unreleased pseudo-version
`v0.10.1-0.20220126230913-d491f1b5c1d2`, unsigned commit
`d491f1b5c1d2e20b85309ccf34fc778d950dc3f2`, tree
`e5bcdbd2d1db14d16da0df023888dd0ba849f8fe`, at
2022-01-26T23:09:13Z. Three post-release commits add a style, aliases, and
helper scripts. The next commit changes the module path to the distinct
`github.com/alecthomas/chroma/v2`. That line has 38 stable releases and four
v2.0.0 alphas; latest v2.27.0 declares Go 1.25. Current master declares the
distinct `/v3` module and has only alpha tags. No pseudo-version or alternate
module path is an in-place stable v0 update.

Selected v0.10.0 declares Go 1.13 and requires regexp2 v1.4.0 plus test-only
testify v1.7.0 and indirect go-spew v1.1.1. Its standalone complete test graph
selects eight modules across nine edges. Only Chroma and testify declare Go
versions, both 1.13; all others have no Go directive. The closure preserves
the retained Go 1.18 floor. Project MVS already selects v0.10.0 and later
regexp2/testify versions, so the changed-selection closure is empty.

The selected checksum pair is
`h1:7XDcGkCQopCNKjZHfYrNLraA+M7e0fMiJ/Mfikbfjek=` /
`h1:jtJATyUxlIORhUOFNA9NZDWGAQ8wpxQQqNSB4rjA/1s=`; the checksum database
agrees. Proxy ZIP SHA-256 is
`beb07b996ee33bc052fe039c93d1c0726e61bcc4819ca39f7bf63304f2ae8c49`.
All 696 proxy regular files match the tag; eight nested-module regular files
and seven tracked Hermit symlinks are correctly omitted. The pseudo-version
pair is `h1:Gg09t2u+C08At6TYucNrD3Cbaq97SUHax84BzQwRTgU=` / the same go.mod
sum. Its ZIP SHA-256 is
`033bd0de6588cb5c9508b9ced1eb0a0620eb9861fc43a11843d416ca3e245d79`;
all 699 regular proxy files match the exact commit with the same justified
omissions.

Exact external `go get github.com/alecthomas/chroma@v0.10.0` exits zero and
changes no byte. Old and replayed states are identical at 234 selected
modules, 3,580 graph edges, 429 complete packages, 1,043 go.sum lines, and a
356-line unapplied tidy projection. No dependency implementation was
authorized or created.

Exactly 33 Chroma packages load through
`plybuild/cmd -> go-term-markdown -> chroma`; repository source has no direct
Chroma import. Go-term-markdown uses the Chroma root, formatters, lexers, and
styles packages, including lexer lookup/analysis/fallback, tokenization,
coalescing, TTY8 formatting, and Pygments styling.

The project's Markdown reference-definition contract passes ten times. The
direct consumer's focused `TestRender/codeblock` golden fails all ten
project-MVS repetitions: rendering completes, but ANSI intensity/color
sequences and one-space code padding differ from its v0.7.1-era golden. The
613-line report SHA-256 is
`b6b6f24dcd9c1b97c9280e01b4277aa3614b5513e35bd8ad46e6cd7c9cea0d3c`.

Writable proxy and exact-tag sources independently verify, list 34 packages,
and pass complete count-1, count-10, and race tests without source mutation.
Mandatory vet fails both with the same normalized 7,175 unkeyed
`github.com/alecthomas/chroma.Rule` literal diagnostics across 206 lexer
files. Sorted diagnostic SHA-256 is
`cbfbdc99bce5e357f83b20eec1a08644ddb86a16ed90f3a2c176cb501a68604e`.
The loaded-behavior and vet failures are independent stop conditions; no
downstream changed-selection gates apply.

Fresh primary vulnerability data was updated 2026-09-02T19:12:04Z and has
1,392 module records with no Chroma record. Old/replay IDs and normalized
traces are identical at 20 unique Darwin symbol IDs, 30 Darwin module IDs,
and 20 unique Windows symbol IDs. No trace contains Chroma.

## Chroma Evidence And Tools

- Decision evidence beneath the active scratch root has 4,335 verified
  entries. Evidence-manifest SHA-256 is
  `bbd6a41e0cff3248bacc2fef6630ed3679a3299484cb21446fc2841a39a742d5`;
  decision-summary SHA-256 is
  `476f926b06cd8fbfce8ac22c9e0ed1102c8f5c8f8e2070bcb1636a825225927f`.
- Exact Go 1.26.7 was recreated beneath scratch. Binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  It was first in PATH with GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and no
  ambient GOFLAGS.
- Fresh govulncheck v1.7.0 reports Go 1.26.7 and the database identity above.
  Its nonportable build SHA-256 is
  `30e6f51fc7d968d10f0d789c82dce56ff0da880c7dfed7c82e08e6feaa554f24`;
  version and functionality are the durable proof.
- Preserve the portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections, including go-colorful mutable telemetry,
  btree regression, cast whitespace-path, source-archive normalization, Check
  history table, Errgo preflight runner, Resty committed projection/signal
  fixture, Kingpin external report, Units launcher timing, Assert external
  cache/TMPDIR, and Repr cache/mktemp/Docker/manifest corrections.
- The deliberately purged former two-entry recovery manifest remains
  historical SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path
`github.com/alecthomas/colour v0.1.0` as the next single P7 group. A fresh
post-Chroma survey finds only stable v0.1.0 in the exact-path proxy list.
Proxy `@latest` and exact Go `@latest`, `@v0`, and `@master` all resolve
v0.1.0 at 2019-11-01T02:47:59Z; its module has no Go directive. Treat release
qualification, repository/source identity, signature and history, complete
closure, self-tests, loaded consumers, actual symbols, and vulnerability
effect as unknown until independently resolved.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Change only Colour's exact selected metadata and explained minimal
MVS closure if a qualified changed selection preserves the floor and passes
every applicable gate. Do not combine Chroma, another group, or P8.
