# Quality Upgrade Handover

Generated: 2026-09-05T01:13:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only go-cmp commit
  `c314bcb440b5f94871d71a249bca7ae7f87d5543`, exact parent
  `cbdb0a915d7d30bb6633cfbb2ccf0323aea1b77a`, clean tree
  `9a212b377cb08e7d6e3044fb6dae071058b6242b`.
- The documentation handoff must have exact parent `c314bcb`; ordinary and
  ignored status must be empty afterward. The answered go-cmp archive and the
  Hashicorp HCL NEXT archive must link reciprocally.
- Only `go.mod` and `go.sum` changed in the implementation. Production Go,
  API/CLI, language/toolchain declarations, quality apparatus, Docker/release
  inputs, packaging, publishers, and P8 code are unchanged.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, successor launch, retained
  evidence/image deletion, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through accepted `github.com/google/go-cmp
v0.6.0`. Latest go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3,
fatih/color v1.19.0, fsnotify v1.10.1, and latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` remain rejected for their recorded floor
or loaded-behavior failures. The accepted gomarkdown fallback remains
`v0.0.0-20260824151336-45814d58469f`. P8 remains queued.

Only the launcher's mutable header and prompt regions may change during a
handoff. Keep caches, projections, reports, artifacts, build contexts, and
audit evidence outside the worktree.

## Accepted Google Go-Cmp Group

Latest v0.7.0 is rejected because its independently reproduced `go 1.21`
declaration exceeds the retained Go 1.18 compatibility floor. Primary proxy,
sumdb, and upstream evidence binds v0.7.0 to lightweight tag commit
`9b12f366a942ebc7254abc7f32ca05068b455fb7`, published
2025-01-14T18:15:44Z, with checksum pair
`h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=` /
`h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=`.

Selected indirect go-cmp moved from v0.5.9 to v0.6.0, the highest compatible
release because the only higher listed release is rejected v0.7.0. Version
v0.6.0 resolves to lightweight tag commit
`c3ad8435e7bef96af35732bc0789e5a2278c6d5f`, published
2023-08-31T17:32:40Z, declares Go 1.13, and has checksum pair
`h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=` /
`h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=`. Both release archives
contain 48 files byte-identical to their tag commits. Both tags are commit
objects rather than annotated tags. GitHub records valid commit signatures;
local verification could not run because `gpg` is absent. The primary Go
vulnerability module index contains no go-cmp entry.

The six-commit v0.5.9 -> v0.6.0 history removes purego fallbacks, adds Go 1.20
testing, pins workflow inputs, uses identifier links, and adds
`cmpopts.EquateComparable`. Both states select 234 modules, have 3,557 graph
edges and 429 complete test packages, and load the same five go-cmp packages.
Only the go-cmp selection and matching main edge change. The path is
`plybuild/cmd -> mvn-pom-mutator/pkg/pom -> go-cmp/cmp`.

Exact `go get github.com/google/go-cmp@v0.6.0` replaces the indirect
requirement and adds the v0.6.0 checksum pair while retaining v0.5.9 sums. The
299 -> 308-line tidy result remains an unapplied projection; comparing fully
tidied projections isolates only the version and checksum pair. Candidate
module complete tests/race, old and new Ply consumers at count 10, repository
build/tests/race/vet, pinned lint, byte-identical help, identical API/CLI,
launcher/Make/preflight, host/snapshot/Docker, audit, empty-HOME count-2, and
clean-tree gates pass. Vulnerability ID sets remain exact 20/30/20 for Darwin
symbol, Darwin module, and Windows symbol.

The candidate module's optional `go vet ./...` emits seven diagnostics only in
upstream test files (two deliberately malformed struct tags and five unkeyed
internal test-proto/test-struct literals); its required complete module tests
and race pass. The packaged mvn-pom-mutator v0.2.3 test has a pre-existing
`Marshall` arity mismatch in both old and new states; the actual six-package
Ply consumer population passes in both.

## Accepted Evidence

- Selection root `/private/tmp/ply-p7-go-cmp-selection.cbdb0a9.kbzc4C`:
  verified 29,812-entry manifest SHA-256
  `dcd18baaa45eb8ad59b268fa7c356164dabc5d4c9836423c3fd1c3aa6d0d1797`;
  selection-summary SHA-256
  `5f76248b77f1451bfb3304b0857ba9464c3a5f75921ba6df26b750080c361ee1`.
- Commit-bound schema-2 review
  `/private/tmp/ply-p7-go-cmp-quality-review.c314bcb.XeWg4T`: verified
  9,487-entry manifest SHA-256
  `24ede5a473b9ac9f47839b121d25e9a8d921369c5ba149bf42abd902b5271727`;
  manual-evidence SHA-256
  `6f7a469bf4fa8da00a44570d4b7c669f64c5a31dbc116b4304c0604c505ba1ac`.
  All six focused receipts pass; all 112 governed files and 258 subjects are
  unchanged from the parent.
- Exact `make quality` root
  `/private/tmp/ply-p7-go-cmp-quality-parent.c314bcb.zxFKyw`: verified
  252,077-entry manifest SHA-256
  `6d8104b43fdcb838ec8c0955870798be38a3b4a1f9a7e73f22e683ee566a9f69`;
  Q0-Q2 scorecard SHA-256
  `6531ade31d4c504c304cd707cbab158e9433d487ec8fa96d05a8277d9dc149e1`.
  Its exact 21-stage ledger exits 0 with 27/27 rows at L2, 80/80 mutations,
  8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
  not-comparable, or dirty counts.
- Independent regression root
  `/private/tmp/ply-p7-go-cmp-regression-gate.c314bcb-final.rMW0l5`: verified
  30,255-entry manifest SHA-256
  `00b736b53494a0c4f27ee7503ecef1c7c7bf74f00bdaab85f7afa9eab19ed3c5`;
  regression-summary SHA-256
  `a27e5f4d00fa0f60d9246309e1485353cdb2aacecd4a3d80e0d3bf64efe2323e`.
  It independently binds closure, module/consumer/repository gates, lint,
  help/API/CLI, preflight, host acceptance, focused and exact audits,
  20/30/20 vulnerabilities, empty-HOME count-2, and cleanliness. Full audit
  exits expected 1, never 2, only for Q3.1, Q3.3, Q3.4, and Q3.7.

The first exact-quality attempt is retained at
`/private/tmp/ply-p7-go-cmp-quality-parent.c314bcb.n2kYFd`. Midnight external
tmp cleanup removed the SDK's standard-library source files while leaving its
binary and hash intact, so the run stopped during mutation meta-testing. The
official Go 1.26.7 archive was re-fetched at exact SHA-256
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`
and overlaid at the mandated path. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` has verified two-entry manifest
SHA-256 `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

A prerequisite replay also corrects the prior handover: the btree regression
manifest retains its recorded digest and 24,197 entries, but 163 mutable
cache/HOME entries no longer verify after preliminary version enumeration
mutated 161 `@v/list` files, one sumdb latest record, and one telemetry count.
Stable regression-summary SHA-256
`f348bd0b9e0f68c698c43f466904b66c1c3db6eab5e783b14364aaf40b6b1ea2`
is unchanged. The other three btree manifests completely verify, and current
regression claims are superseded by the fresh go-cmp evidence above.

Tool identities remain Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
golangci-lint 2.12.2
`3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`;
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`;
apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`;
and govulncheck v1.7.0
`0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected indirect `github.com/hashicorp/hcl v1.0.0`
as exactly one bounded P7 module group. Resolve the current latest release and
every potentially compatible release from fresh Go proxy, checksum-database,
and upstream evidence. Do not assume a target version, Go floor, closure, or
loaded status from secondary sources.

Measure exact old/candidate selections, graph edges, complete package
population, loaded packages and dependency path, checksums, exact-get diff,
tidy projection, focused behavior if loaded, candidate module tests, repository
quality, help/API/CLI identity, and vulnerability populations. Implement one
exact selection only if it preserves the retained Go 1.18 floor, has a minimal
explained closure, and passes every contract. Otherwise record rejection with
no dependency metadata edit.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source change, language/toolchain or quality-policy change,
packaging, publication, or P8. Do not push, merge, publish, release, delete
evidence, stash, revert, launch a successor, or remove the worktree.
