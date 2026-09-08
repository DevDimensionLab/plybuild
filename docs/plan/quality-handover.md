# Quality Upgrade Handover

Generated: 2026-09-08T12:07:54+02:00

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
- Logex, OpenCensus Proto, and Crypt were retained without implementation
  commits. Speakeasy v0.2.0 implementation
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remains an ancestor, as do
  Circbuf implementation `3be2183ee310ccdc358ce4ed372c0785de25b88b`
  and operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`.
- The answered Logex archive and sole NEXT Readline archive link reciprocally.
  No `.agent-task/current.md` or repository `.quality/manual-evidence.json`
  was created. No push, merge, publication, release, stash, revert, successor
  launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
and Logex selections. P8 remains queued. All earlier acceptances, rejections,
no-change decisions, evidence corrections, and lifecycle ancestry are final;
do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve the launcher's
scratch cleanup and reciprocal archive contract.

## Retained Logex Group

Retain exact-path `github.com/chzyer/logex v1.2.1` without editing dependency
metadata. The fresh exact proxy lists v1.1.1 through v1.1.10 plus v1.2.0 and
v1.2.1; selected is already canonical exact stable latest at
2022-04-24T13:13:51Z. It declares the exact path, Go 1.15, no requirements,
no retractions, and no deprecation. Short lightweight Git tags v1.0/v1.1
predate modules and are not canonical three-part proxy releases. Historical
`gopkg.in/logex.v1` is a distinct path through v1.1.10; no `/v2` path exists.

Go-import identifies `https://github.com/chzyer/logex.git`. The public
repository is enabled, unarchived, non-fork, and defaults to `master`. Selected
v1.2.1 is unsigned lightweight tag/commit
`2f95bdde8c3c97bfbf6d016fcc410669a895b9e7`, tree
`36fcd9ac7d56d659872b2a6576ca6f66385c938a`, parent
`a21c317abc1e9a4f23ed3455107a4d20375735cc`. Its GitHub Release was
published at 2022-04-24T13:15:45Z. Master
`5a7e37d2e8a8bbe3ef54984ab949eebaa948b8b4`, tree
`f9cc17fbf471a8558b15bdc07ee3e1a9dba4631d`, is a verified GitHub
web-flow-signed merge. Its five post-release commits touch only tests, CI,
module/build input, and are not a release. Exact `@master` resolves unreleased
`v1.2.2-0.20240402154933-5a7e37d2e8a8` and declares Go 1.21.

All twelve proxy ZIPs byte-match their exact Git tags. V1.2.1 ZIP SHA-256 is
`8bc36e064d4f53348c25a5745bd3a9030e3710c7083407da632905114d878bae`;
normalized source-manifest SHA-256 is
`cf48dc5a2062f0aa5877e2e7df8e95153d3f251dd69c03ce0d12c59c72316bb3`.
Its sumdb source/module pair is
`h1:XHDu3E6q+gdHgsdTPH6ImJMIp436vR6MPtH8gP05QzM=` and
`h1:JLbx6lG2kDbNRFnfkgvh4eRJRPX1QCoOIWomwysCBrQ=`.

V1.2.1 is one pure-Go package with a standard-library-only closure and no
Cgo, assembly, generated, build-tagged, or OS/architecture-specific files.
Its closure preserves Go 1.18. Proxy and Git native `TestLogex` fail under Go
1.18.10 and 1.26.7 because two assertions hard-code stale caller line numbers;
vet passes. Unreleased `3e09012` repairs exactly those expectations. This is a
selected-release test gap, not source divergence or candidate regression, and
is an additional reason not to select the branch head.

Independent fixtures pass count-1, two count-10 runs, race, vet, API checks,
and relevant cross-builds under both SDKs. They cover the constructors, Logger
methods, formatting/trace helpers, stack/code/error operations, interfaces,
standard global rebinding, and child-process panic/fatal/debug behavior.
V1.2.0 to v1.2.1 changes only README text and the Go directive from 1.17 to
1.15; apidiff reports no exported API change.

Ply -> Promptui v0.9.0 -> Readline v1.5.1 -> Readline test -> chzyer/test
v1.0.0 -> Logex is the exact dependency-test chain. Readline production does
not import Logex; chzyer/test init calls `Define`, and its comparison/error
paths use `Equal` and `DecodeError`. Focused fixtures pass. Readline full and
focused count/race/cross-build tests pass under both SDKs; its two vet warnings
are invariant historical signatures. Chzyer/test's native `TestMemDisk`
independently panics on `customEqual(nil,nil)` under both SDKs. Neither debt is
caused by Logex. Consumer closures top out at Go 1.17.

## Projection, Quality, And Vulnerability Measurements

Accepted state remains 234 selected modules, 3,583 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded module-backed packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 33 checksum lines and
removes zero. Main-module Go remains 1.18 and toolchain remains 1.26.7.

Readline v1.5.1 and chzyer/test v1.0.0 request Logex v1.2.1; Promptui v0.9.0
and thirteen historical pprof releases request v1.1.10. Ply loads no Logex
package. Exact selected-version get changes no selection or checksum and adds
only a redundant root requirement/main edge. It would produce 3,584 edges and
a 382-line tidy projection that removes that pin, so it was not applied.

The unchanged project passes mod verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical help, API/CLI compatibility
and reports, and empty-HOME count-2. Current and redundant projections have
identical API report `ce39e6fda3f642b10075b07c431ed1a070b65eaeb073b913c28e16cd53fae282`
and CLI report `955f1de3cf3a14a0581af38b527beef6e5a9b5281d5e1cd3330647e821ed7d29`.
No changed stable selection exists, so the changed-selection-only full P7 gate
was not invoked; the accepted XXHash quality result remains authoritative.

Fresh vulnerability data has 1,392 module records and no Logex record.
Current/redundant results are identical: 20 IDs/22 reachable traces for Darwin
and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin module
IDs/findings. No target module, package, symbol, or trace appears.

Logex evidence has 636 verified entries; manifest SHA-256 is
`aa33eeb8ea36b17f6830ae7f063fbe1c39ec9ac5d3c9e2417fa2fe844ceb8516`.
Decision-summary SHA-256 is
`6720115c3316e14b89a3bedcb40a59faba4a4f898727051177a682bb1435b49d`.

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

Independently evaluate selected exact-path `github.com/chzyer/readline v1.5.1`
as the next single P7 group. The exact proxy lists only v1.5.0 and v1.5.1;
selected is already `@latest` at 2022-07-15T12:48:48Z and declares Go 1.15.
Short repository tags v1.0 through v1.4 are proxy-absent and need precise
module/release qualification.

Ply reaches Readline in production through `github.com/manifoldco/promptui
v0.9.0`. Readline v1.5.1 requires chzyer/test v1.0.0, Logex v1.2.1, and its
historical x/sys pseudo-version. Preserve the now-final Logex decision while
separating Readline native tests, Promptui production behavior, and the
dependency-test-only chzyer/test/Logex path.

The public `github.com/chzyer/readline` repository is enabled, unarchived,
non-fork, and defaults to `main`. Selected v1.5.1 has an annotated tag that
needs exact tag-object, target-commit, and signature verification. Main head
`9dfc369f8652ba9013dadffd2d2efeada64fe44d`, tree
`c0ed5f5684075d6df7c6e1eb34e15e567e11d3a2`, parent
`fcb4d79af3fbe295b4cb6360e14b8c0b8337353f`, dated
2025-06-20T03:33:30Z, has three post-release commits and is not a release.
Independently establish canonical identity, tag/release history, complete
floor, terminal/OS behavior, exact Promptui consumers, projection, and
vulnerability effect. Do not implement or combine another group until that
bounded decision is complete.
