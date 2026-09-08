# Quality Upgrade Handover

Generated: 2026-09-08T14:17:43+02:00

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
- Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt were retained without
  implementation commits. Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor, as do
  Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`.
- The answered Fnmatch archive and sole NEXT Imaging archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, and Fnmatch selections. P8 remains queued. All earlier
acceptances, rejections, no-change decisions, evidence corrections, and
lifecycle ancestry are final; do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve the launcher's
scratch cleanup and reciprocal archive contract.

## Retained Fnmatch Group

Retain exact-path `github.com/danwakefield/fnmatch
v0.0.0-20160403171240-cbb64ac3d964` without editing dependency metadata. A
fresh exact proxy version list is empty, and exact `@latest` is the selected
pseudo-version at 2016-04-03T17:12:40Z. Its proxy-synthesized module file
declares only exact path `github.com/danwakefield/fnmatch`, without a Go
directive, requirements, deprecation, or retractions.

Sumdb records selected source/module hashes
`h1:y5HC9v93H5EPKqaS1UYVg1uYah5Xf51mBfIoWehClUQ=` and
`h1:Xd9hchkHSWYkEqJwUGisez3G1QY8Ryz0sdWrLPMGjLk=`. The proxy ZIP SHA-256 is
`f601e8d25a43ed32e00851e1686a93b0175dadea8f4e32c8af2f1533f20736bc`;
proxy and exact Git files match, with normalized five-file source-manifest
SHA-256 `8bedd8645805cd07f06541f956bbd16e7975fefbc23e58da308cc3e5a69420bb`.

Go-import metadata identifies `https://github.com/danwakefield/fnmatch.git`.
The public repository is enabled, unarchived, non-fork, and defaults to
`master`. Selected and master are the same unsigned commit
`cbb64ac3d964b81592e64f957ad53df015803288`, tree
`e31339f278164c2b9c1c4c45d08fb964b3fcae0f`, parent
`eb9738ef552dd59a56a5953a4de6216f70564908`, in a four-commit history. The
repository has no tags or GitHub Releases. Its 2023 `pushed_at` date comes from
an unmerged pull-request ref, not a commit beyond 2016 master.

Open and closed PR heads are unmerged fork commits and return unknown revision
through the exact parent path. Gandarez fork tags v0.1.0/v0.1.1 declare
`github.com/gandarez/fnmatch`; Slashid master declares
`github.com/slashid/fnmatch` and Go 1.20; other forks and inspired modules are
also alternate paths. The kballard/lilyball gist and Daniel Wakefield's gist
fork establish source ancestry, not a release identity. No exact stable or
prerelease tag, higher qualified stable candidate, redirect, or alternate
identity changes the exact selection.

The complete selected module closure is Fnmatch alone; its source and tests
use only the standard library. The package is pure Go, with no Cgo, assembly,
generated files, build tags, or platform-specific files. Proxy and exact Git
native suites pass count-1, two independent count-10 runs, and race under
exact Go 1.26.7 and contained Go 1.18.10. Independent behavior fixtures and
Windows/Linux/FreeBSD/js-wasm cross-builds also pass under both SDKs. This
complete execution, not the directive-free synthesized module file, proves
the retained Go 1.18 floor.

The sole function `Match(pattern, s string, flags int) bool` and flag constants
implement rune-aware wildcard, bracket/range/negation, escape, slash, Unicode
case-fold, leading-directory, malformed-pattern, and period behavior. `/` is
the separator on all platforms. Independent fixtures preserve two documented
BSD-derived period quirks and prove the selected implementation's actual panic
for `Match("*", "", FNM_PERIOD)`. Vet consistently reports one unreachable
statement at `fnmatch.go:91`; upstream PR #1 removes it, but the unmerged commit
is not an exact-path version. These are historical selected-release gaps, not
a qualified candidate.

The graph retains only the root indirect requirement and historical
`github.com/alecthomas/chroma v0.7.1 -> Fnmatch` edge. Chroma v0.7.1 called
Fnmatch for lexer filename globs, and focused consumer fixtures pass under both
SDKs. Project MVS selects final Chroma v0.10.0, which uses `filepath.Match`,
requires and imports no Fnmatch, and loads no Fnmatch package. `go mod why -m`
correctly says the main module does not need it. Do not reopen Chroma or use
tidy to remove the root requirement.

## Projection, Quality, And Vulnerability Measurements

Accepted state remains 234 selected modules, 3,583 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 33 checksum lines and
removes zero. Main-module Go remains 1.18 and toolchain remains 1.26.7.

Exact selected Fnmatch `go get` is entirely inert: it changes no requirement,
selection, edge, checksum, tidy projection, or status. It was not applied. The
unchanged project passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical root/status/upgrade/build
help, API/CLI compatibility and reports, and empty-HOME count-2. API and CLI
report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
No changed stable selection exists, so the changed-selection-only full P7 gate
was not invoked; the accepted XXHash quality result remains authoritative.

Fresh vulnerability data has 1,392 module records and no Fnmatch record.
Canonicalized current/no-op sets are identical: 20 IDs/22 reachable traces in
Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin
module IDs/findings. No target module, package, symbol, or trace appears.

Fnmatch evidence has 72 verified entries; manifest SHA-256 is
`97130774da6c5f883f1c8fd0afa8190490be6aa3f8c998af4deaa9b7c79d44f1`.
Decision-summary SHA-256 is
`96f4a6d499ac460e99d3c4b7b9395b2ed826374d79fca4602e89c359b676eae1`.

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
  apidiff was used only for verified zero-diff function, not as a replacement
  portable receipt.
- Retain accepted XXHash v2.3.0 and Speakeasy v0.2.0 moves and every earlier
  exact decision, including retained Readline, Logex, OpenCensus Proto, and
  Crypt selections. Answered archives remain authoritative detail.
- Authoritative accepted Q0-Q2 scorecard SHA-256 remains
  `579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

## Next Objective

Independently evaluate selected exact-path
`github.com/disintegration/imaging v1.6.2` as the next single P7 group. A
minimal post-Fnmatch survey finds 15 exact proxy versions from v1.0.0 through
v1.6.2; exact `@latest` is the selected stable v1.6.2 at
2019-11-16T20:43:25Z. Its module file declares the exact path, no Go directive,
and `golang.org/x/image` pseudo-version
`v0.0.0-20191009234506-e7c1f5e7dbb8`.

The public enabled/unarchived/non-fork repository defaults to `master`.
Selected v1.6.2 is a stable GitHub Release and lightweight tag at verified
commit `acabd8315e63bfcaac97d52d68a7a0b88d2eea93`, tree
`6584cbb2a26e4d38bfec8f2587633f234500810e`, with parents `9aab30e...` and
`675e3c2...`. Master is the later verified but unreleased commit
`d40f48ce0f098c53ab1fcd6e0e402da682262da5`, resolving only as
`v1.6.3-0.20201218193011-d40f48ce0f09`; do not select it as stable.

Ply loads Imaging in production through `plybuild/cmd ->
github.com/MichaelMure/go-term-markdown ->
github.com/eliukblau/pixterm/pkg/ansimage -> Imaging`. The graph also has the
main indirect pin and Imaging -> x/image edge; project MVS selects x/image
v0.5.0. Resolve release identity, complete closure floor, image/native
behavior, production consumers, exact no-op or candidate projection, and
vulnerability effect without combining the x/image group or an unrelated tidy
cleanup.
