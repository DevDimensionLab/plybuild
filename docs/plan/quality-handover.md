# Quality Upgrade Handover

Generated: 2026-09-07T13:33:22+02:00

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
- Operator-authorized lifecycle repair commit
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`, parent `6f4d02e`, tree
  `81ca3d7a6537ce2b9c26337c2bf3d4d67c28223a`, remains intentionally between
  that implementation and Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  Preserve it. The Repr handoff is the required direct child of f9f0f76.
- The answered Kong archive and sole NEXT Chroma archive link reciprocally.
  Only launcher mutable regions change during handoff. Ordinary and ignored
  status must end empty.
- Kong was rejected before implementation. Relative to accepted go-cmp commit
  `c314bcb`, accepted dependency metadata still adds exactly 27 checksum lines;
  no Kong requirement or checksum was manufactured.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, external source, report, evidence tree,
and build context must remain beneath `$CODEX_SESSION_SCRATCH_ROOT`. The
launcher removes per-turn scratch on success, failure, and interruption and
removes supervisor logs after successful completion. Preserve the operator-
authorized lifecycle controls and scratch-only mutation, snapshot, Docker,
and acceptance cleanup.

The historical `/private/tmp/ply-*` evidence and tool roots were deliberately
purged after their receipts were recorded. Treat those digests as historical
records, not live roots. P2A-P6 are complete. P7 remains active after exact Go
1.26.7 and bounded dependency decisions through rejected Kong v1.3.0. All
prior accepted, rejected, and no-change decisions and evidence corrections
remain final. P8 remains queued.

## Rejected Alecthomas Kong Group

The exact `github.com/alecthomas/kong` proxy path lists 75 stable semantic
tags from v0.1.0 through v1.16.1, with no prerelease or retraction. Proxy
`@latest` and exact Go `@latest`/`@v1` resolve canonical-latest stable v1.16.1
at 2026-08-09T07:05:31Z. It is a lightweight tag at commit
`0678fd30af8be8bae6dc9f9c6f143cc549450be2`, tree
`250e0703356ebc561a781b385d4baf7de410e4d6`; GitHub validates the embedded
commit signature. It declares Go 1.20 and exceeds the retained Go 1.18 floor.
Current master is newer unreleased pseudo-version
`v1.16.2-0.20260828071222-a60008c6dae2`, commit
`a60008c6dae2b98a1fbec12d1aebd43ac5ab76e0` at
2026-08-28T07:12:22Z. It also declares Go 1.20 and is neither a stable release
nor the exact-path latest selection.

V1.3.0 is the highest stable release whose declaration and complete changed
closure preserve Go 1.18. It is a lightweight tag at unsigned commit
`7bbb0b76ada1610f18cf71c54cca74209da88bd8`, tree
`96a61d79e4745b05738a01906c43a7b8e5820b5a`, parent
`373692af87b177d48898c89ad53b6054f5b339bf`, published
2024-11-01T01:25:41Z. V1.4.0 is the immediately newer release and first to
declare Go 1.20. The selected pseudo-version is unsigned commit
`0548c6b1afaea7616098cec708d9bbf49eef2cd6`, tree
`15a2da506ba84c893a3d6cae5e4cdf02f67dc5ac`, parent
`33b8d2b19cff931184e23eff38c3a371948f527d`, at
2019-07-08T04:11:08Z.

The enabled, unarchived, undisabled, non-fork repository uses default branch
`master`, has no GitHub Release objects, alternate major module path, or
formal deprecation marker, and all 75 tags are lightweight. Relevant history
spans 245 commits and 67 paths from the selected pseudo-version to v1.3.0,
including the stable v1.0 line and v1.3.0 recursive provider-parameter
injection; the next release adds hooks and validation work but raises Go to
1.20.

Old checksum pair:
`h1:C4Q9m+oXOxcSWwYk9XzzafY2xAVAaeubZbUHJkw3PlY=` /
`h1:+inYUSluD+p4L8KdviBSgzcqEjUQOfC5fQDRFuc36lI=`. Candidate pair:
`h1:YJKuU6/TV2XOBtymafSeuzDvLAFR8cYMZiXVNLhAO6g=` /
`h1:IDc8HyiouDdpdiEiY81iaEJM8rSIW6LzX8On4FCO0bE=`. Independent checksum-
database lookup agrees. Old/candidate proxy ZIP SHA-256 values are
`4292d9b6903d67f060d3bd57ffca0a4ebca359824ce2d32a512ac1b963fa3dc0`
and `446c6e411b1291658281e3413476787fda86d9c8d8586d7e1f851d0b4f8943f3`.
All 64 candidate proxy regular files match the tag commit. Five tracked Hermit
symlinks and seven regular files in nested module `_examples/server` are
correctly absent from the proxy module ZIP.

Candidate Kong declares Go 1.18 and requires assert/v2 v2.11.0, Repr v0.4.0,
and gotextdiff v1.0.3. Assert/v2 declares Go 1.18 and requests Repr v0.4.0 and
gotextdiff v1.0.3; gotextdiff declares Go 1.16. Existing accepted Repr v0.5.4
at Go 1.18 dominates v0.4.0. The complete closure therefore preserves the
retained floor by declaration.

Exact external replay projects 234 -> 236 modules, 3,580 -> 3,584 graph
edges, the byte-identical 429 complete packages, 1,043 -> 1,045 go.sum lines,
and a 356 -> 359-line unapplied tidy diff. Three selections change: Kong
advances and assert/v2 plus gotextdiff enter. The four new edges are main ->
Kong and Kong's three requirements; the historical Chroma v0.7.1 -> old-Kong
node and its five edges remain in the pruned graph. Only candidate Kong's
checksum pair enters go.sum. Tidy would remove the unused pin and pair and was
not used as implementation.

Kong appears in zero old or candidate complete packages. Repository Go source
has no Kong import, and `go mod why -m` says the main module does not need it.
It is unloaded historical MVS graph debt, including dependency tests.

The standalone candidate suite selects four modules across six edges. Writable
proxy and exact-tag sources independently verify, list one package, and pass
count-1, count-10, and race tests without source mutation. Both mandatory vet
runs fail with byte-identical 77-line malformed-struct-tag diagnostics: 54 in
`help_test.go`, eight in `kong_test.go`, one in `resolver_test.go`, and 14 in
`tag_test.go`. Diagnostic SHA-256 is
`d9cfb57f6329a7a4034e6799b29830db482a77302ba4dc68f4759a047cf0c41e`.
The explicit stop rule rejects v1.3.0 before downstream repository, snapshot/
Docker, exact-quality, and audit gates. No dependency metadata changed.

Fresh primary vulnerability data was updated 2026-09-02T19:12:04Z and has
1,392 module records with no Kong record. Old/candidate IDs and normalized
traces are byte-identical at exact 20 Darwin symbol, 30 Darwin module, and 20
Windows symbol findings. Normalized trace SHA-256 values are respectively
`bbade0271bfcd4b612bfb9f760b8f7f2045d320264b7ef10523c70fc341e96bc`,
`17ac7fe0974efb5dad0636d85c1a19ae7f2d67a4d23dafea026ad52a20ecf8c7`,
and the same `bbade027...` Windows value.

## Kong Evidence And Tools

- Decision evidence beneath the active scratch root has 211 verified entries.
  Evidence-manifest SHA-256 is
  `6148c409545f453a78ffdc3a0a10934b75e082c1681f2805c29c8253a23d01ae`;
  decision-summary SHA-256 is
  `028c0fd0025a5f5de433b65d63a5258e18c6e959b440446e9834b03925f98b36`.
- Exact Go 1.26.7 was recreated beneath scratch. Binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  It was first in PATH with GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and no
  ambient GOFLAGS.
- Fresh govulncheck v1.7.0 reports Go 1.26.7 and the database identity above.
  Its current SHA-256 is
  `ac8675f248a7686f926eaf2ad000d0731f449b5da4ab7155862fa15bdcda0601`.
  Preserve earlier nonportable build receipts
  `0c536fb25db0d42aff2e5f4496d37a2a89b9f6c449466e2e493975db2652bb2e`
  and `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`
  as history; version and functionality, not byte identity, are required for
  fresh builds.
- Preserve the official golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 binary receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
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
`github.com/alecthomas/chroma v0.10.0` as the next single P7 group. A fresh
post-Kong survey says proxy `@latest` remains stable v0.10.0 at
2022-01-12T10:49:38Z and declares Go 1.13, while exact `@master` resolves
unreleased pseudo-version `v0.10.1-0.20220126230913-d491f1b5c1d2` at
2022-01-26T23:09:13Z. Treat canonical qualification, source/signature and
release history, alternate `/v2` lineage, exact closure, self-tests, and
vulnerability effect as unknown until independently resolved.

Chroma is directly pinned only as an indirect requirement. Thirty-three
Chroma packages load through `plybuild/cmd -> go-term-markdown -> chroma`,
while repository Go source has no direct Chroma import. Treat those as starting
observations and identify the actual consumer packages and used symbols.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Change only Chroma's exact selected metadata and explained minimal MVS
closure if a qualified changed selection preserves the floor and passes every
applicable gate. Do not combine Kong, another group, or P8.
