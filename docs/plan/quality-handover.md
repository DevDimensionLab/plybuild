# Quality Upgrade Handover

Generated: 2026-09-08T13:21:10+02:00

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
- Readline, Logex, OpenCensus Proto, and Crypt were retained without
  implementation commits. Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor, as do
  Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`.
- The answered Readline archive and sole NEXT Fnmatch archive link
  reciprocally.
  No `.agent-task/current.md` or repository `.quality/manual-evidence.json`
  was created. No push, merge, publication, release, stash, revert, successor
  launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, and Readline selections. P8 remains queued. All earlier acceptances,
rejections, no-change decisions, evidence corrections, and lifecycle ancestry
are final; do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve the launcher's
scratch cleanup and reciprocal archive contract.

## Retained Readline Group

Retain exact-path `github.com/chzyer/readline v1.5.1` without editing
dependency metadata. The fresh exact proxy lists only v1.5.0 and v1.5.1;
selected is already canonical stable latest at 2022-07-15T12:48:48Z. Both
declare the exact path and Go 1.15 with no deprecation or retraction. Selected
sumdb source/module hashes are
`h1:upd/6fQk4src78LMRzh5vItIt361/o4uq553V8B5sGI=` and
`h1:Eh+b79XXUwfKfcPLepksvw2tcLE/Ct21YObkaSkeBlk=`. Its ZIP SHA-256 is
`ce25854a8beae5c20bdde840d5142e6fbd1f86f0e58442705b8fb21dfce48501`;
its corrected 48-file normalized source-manifest SHA-256 is
`983bc675215c194e1d126f58b44d3a6211f39eb3af5ce2ae1d06613768f501c5`.
Proxy and exact Git content match.

Go-import identifies `https://github.com/chzyer/readline.git`. The public
repository is enabled, unarchived, non-fork, and defaults to `main`. Selected
v1.5.1 is unsigned annotated tag object
`704f339125f222987e1fde71641f3185f6eda206` targeting unsigned commit
`7f93d88cd5ffa0e805d58d2f9fc3191be15ec668`, tree
`d842017d1ed9d9fd529cce8e199c3a3a69e68e0c`, parent
`8e4bd417b9169c9482a55f3faaeef208b5bf7eb4`. V1.5.0 is an unsigned
lightweight tag with a GitHub Release; v1.5.1 has no GitHub Release but is an
exact proxy/tag release.

V1.0 through v1.4 are historical GitHub Releases/tags lacking a patch
component and `go.mod`. They and synthesized v1.N.0 spellings are proxy-absent;
their commits resolve only as exact-path v0 pseudo-versions. Historical
`gopkg.in/readline.v1` and `gopkg.in/chzyer/readline.v1` are distinct vanity
paths, and no `/v2` module exists.

Main `9dfc369f8652ba9013dadffd2d2efeada64fe44d`, tree
`c0ed5f5684075d6df7c6e1eb34e15e567e11d3a2`, is three first-parent commits
after selected and resolves only as unreleased
`v1.5.2-0.20250620033330-9dfc369f8652`. Its primary-evidence exact parent is
`fcb4d7d9a9f653462a7adf557fb1f931f00391f2`, correcting the incoming
near-match. The dev_v2 branch is also only an unreleased v1.5.2 pseudo-version.
Neither is selectable as a stable release.

The complete closure has four modules and tops out at Go 1.17: Readline,
chzyer/test v1.0.0, and retained Logex v1.2.1 declare Go 1.15; selected x/sys
declares Go 1.17. Readline is pure Go with no Cgo, assembly, or generated
files. Legacy build tags split Windows kernel32 console/syscall behavior from
Unix raw-mode handling, Linux/BSD ioctls, and the AIX/Solaris x/sys/unix path.
Windows, Linux, Darwin, BSD, Solaris, and AIX cross-builds pass under Go
1.18.10 and Go 1.26.7.

Proxy and Git native suites pass count-1, two independent count-10 runs, and
race under both SDKs. Vet consistently reports only the historical
nonstandard `WriteTo(io.Writer) (int, error)` and `ReadRune() rune` method
signatures. They are selected public APIs, not test failures, source divergence,
or a higher-candidate regression. V1.5.0 to v1.5.1 has three compatible API
additions: `CaptureExitSignal`, `(*Instance).CaptureExitSignal`, and `CharO`.

Ply's production path is `plybuild/cmd -> Promptui v0.9.0 -> Readline v1.5.1`.
Promptui uses Config initialization, NewEx, listener-driven Readline editing,
Write/Close, masking, cursor/screen output, Vim mode, default completion, and
cancelable input. It disables persistent history with `HistoryLimit=-1` and no
history file. Ply constructs Promptui `Prompt` twice, not `Select`. Fresh
Promptui historical-consumer projections pass count-1, two count-10 runs,
race, vet, and Windows build under both SDKs. Readline's chzyer/test -> Logex
edge is dependency-test-only; Readline production imports neither.

## Projection, Quality, And Vulnerability Measurements

Accepted state remains 234 selected modules, 3,583 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 33 checksum lines and
removes zero. Main-module Go remains 1.18 and toolchain remains 1.26.7.

Ply reaches selected Readline through its existing main indirect requirement
and Promptui v0.9.0 production import. Exact selected-version get is entirely
inert: it changes no requirement, selected module, edge, checksum, tidy
projection, or status. It was not applied.

The unchanged project passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical help, API/CLI compatibility
and reports, and empty-HOME count-2. Current and no-op projections have
identical API report `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and CLI report `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
No changed stable selection exists, so the changed-selection-only full P7 gate
was not invoked; the accepted XXHash quality result remains authoritative.

Fresh vulnerability data has 1,392 module records and no Readline record.
Current/no-op results are identical: 20 IDs/22 reachable traces for Darwin
and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin module
IDs/findings. No target module, package, symbol, or trace appears.

Readline evidence has 745 verified entries; manifest SHA-256 is
`8ab036d5f96a8d92ebf682ed0fef1d5f116621dff18d670bb200172581158a00`.
Decision-summary SHA-256 is
`435eafc93ae6df466970eb1af57a8127e5f8c389863c82fffadec236ca6bbf8a`.

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
  apidiff was used only for its verified zero-diff function, not as a new
  portable binary receipt.
- Retain accepted XXHash v2.3.0 and Speakeasy v0.2.0 moves and every earlier
  exact decision, including retained OpenCensus Proto and Crypt selections.
  The answered archives remain the authoritative detail; do not revisit them.
- Authoritative accepted Q0-Q2 scorecard SHA-256 remains
  `579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

## Next Objective

Independently evaluate selected exact-path `github.com/danwakefield/fnmatch
v0.0.0-20160403171240-cbb64ac3d964` as the next single P7 group. The exact
proxy version list is empty, while `@latest` is the already-selected
pseudo-version at 2016-04-03T17:12:40Z. Its proxy-synthesized module file has
the exact path but no Go directive or requirements.

The main module has an explicit indirect Fnmatch requirement, yet
`go mod why -m` says it is not needed and no Fnmatch package is loaded. The
graph retains a historical Chroma v0.7.1 -> Fnmatch edge while MVS selects the
already-final Chroma v0.10.0. Keep graph history, selected production/test
packages, and tidy's unrelated stale-requirement projection distinct.

The public enabled/unarchived/non-fork repository defaults to `master`, has no
tags or GitHub Releases, and master equals selected unsigned commit
`cbb64ac3d964b81592e64f957ad53df015803288`, parent
`eb9738ef552dd59a56a5953a4de6216f70564908`. Independently establish
canonical/upstream identity, complete floor, matching semantics, native tests,
historical consumers, exact no-op or candidate projection, and vulnerability
effect. Do not treat an untagged commit or repository metadata activity as a
stable release, remove the root requirement through tidy, or combine another
group.
