# Quality Upgrade Handover

Generated: 2026-09-04T20:27:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only gomarkdown fallback commit
  `d4f053808e389cf163665087f78690cca05999ce`, exact parent
  `18d67bb9293cef56b13b8bcb9e3dfbe22fd23045`, clean tree
  `cbf28809e55d95589b86a0d9e57eeaba2401542d`.
- The documentation handoff must have exact parent `d4f0538`; ordinary and
  ignored status must be empty afterward. The answered gomarkdown archive and
  the Google btree NEXT archive link reciprocally.
- Only `go.mod` and `go.sum` changed in the implementation. Production Go,
  consumer rendering, API/CLI, language/toolchain declarations, quality
  apparatus, Docker/release inputs, packaging, publishers, and P8 code are
  unchanged.
- No push, merge, publication, release, stash, successor launch, retained
  evidence/image deletion, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through regexp2 v1.12.0 plus accepted gomarkdown
fallback `v0.0.0-20260824151336-45814d58469f`. Viper v1.16.0, Emoji v2.2.14,
ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, and latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` remain rejected for their recorded floor
or loaded-behavior failures. P8 remains queued.

Only the launcher's mutable header and prompt regions may change during a
handoff. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`; keep caches, projections, reports, artifacts,
build contexts, and audit evidence external.

## Accepted Gomarkdown Fallback

Selected indirect `github.com/gomarkdown/markdown` moved from
`v0.0.0-20221013030248-663e2500819c` to canonical proxy pseudoversion
`v0.0.0-20260824151336-45814d58469f`. Primary proxy and repository evidence
bind it to unsigned commit `45814d58469f9a462c899674ee58ca440b432dbf` at
2026-08-24T15:13:36Z with no proxy tags and Go 1.12 declaration. Its checksum
pair is `h1:21LNG7BIMF2dePpJYpuNzDCaHakI3TYHtOmNRkrNlzA=` /
`h1:JDGcbDT52eL4fju3sZ4TeHGsQwhG9nbDV21aMyhwPoA=`. All 185 proxy source files
match the commit archive.

The fallback is an ancestor of rejected latest and the immediate parent of
`d18ffa2a47071de47ab9ad4dd61f0a01158ff97a`, so it excludes the new
`ReferenceDefinition` AST behavior. Fixed commits for GO-2023-2074,
GO-2024-3205, and GO-2026-5208 are all ancestors of the fallback according to
the primary Go vulnerability records and repository history.

Old and fallback states both select exactly 234 modules, 3,556 graph edges,
429 complete test packages, and four loaded gomarkdown packages through
`plybuild/cmd -> go-term-markdown -> gomarkdown/markdown`. Only gomarkdown and
the main-module edge change. Exact get replaces the indirect requirement and
adds the fallback checksum pair while retaining the old pair. Tidy projects
282 -> 285 lines and was recorded, not applied.

All 34 go-term-markdown v0.1.4 renderer cases ran independently with
`NO_COLOR` unset. Both selections pass the same 30 and fail the same four
fixture-expectation cases; there are zero old-pass/fallback-fail regressions.
The five reference-definition cases that panic on rejected latest all pass.
The focused repository contract and gomarkdown self-tests also pass.

Go vulnerability populations improve exactly 22/33/22 -> 20/30/20 with no
additions. GO-2023-2074 and GO-2024-3205 leave both symbol scans; those two plus
GO-2026-5208 leave the module scan. Build, complete tests/race/vet, pinned
lint, byte-identical help, API/CLI reports, and every distribution and quality
contract pass.

## Accepted Evidence

- Decision replay
  `/private/tmp/ply-p7-gomarkdown-fallback-replay.18d67bb.U5kpjt`: final
  post-use verified 42,127-entry manifest SHA-256
  `b438197999acb18ad6955d66c5d38e0af92df8b0693ecba7445f5a960b1b1fa8`;
  selection-summary SHA-256
  `697d2467a20850656b1213494c9c2fe4557df8a61d1ec9eacff94cdf41d2b28a`.
  The original 42,120-entry seal
  `dfdce8f20a3dbd4257a0a1dc22034d2c903f000a20f935d12c0509cfff303b4c`
  now differs only in four Go telemetry counters touched by the first quality
  cache use; the final seal was taken after all use and verifies completely.
- Commit-bound schema-2 review
  `/private/tmp/ply-p7-gomarkdown-fallback-quality-review.d4f0538.sBeoKD`:
  verified 2,304-entry manifest SHA-256
  `cdd7691693373d726fdec2e954aacdc11bcae0bbffde3a471edd56c34de78ab6`;
  `manual-evidence-schema-2.json` SHA-256
  `ffdee0e8a622a8de34220a7773bf32b87d5368698cc25d18a28c2b49b711cc17`;
  all six focused receipts pass.
- Exact `make quality` root
  `/private/tmp/ply-p7-gomarkdown-fallback-quality-parent.d4f0538-final.2kKFpS`:
  verified 252,314-entry manifest SHA-256
  `70ef46136c13b7e5276ac1ef1af62514013902e93c9e115551d863f0f3a8e40b`;
  Q0-Q2 scorecard SHA-256
  `a2298bb0457549c71c90a9069d3f1ffe219f7f21b62d7b47684d16dc7fdfe870`.
  The exact 21-stage ledger exits 0 with 27/27 rows at L2, 80/80 mutations,
  8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
  not-comparable, or dirty counts.
- Independent regression root
  `/private/tmp/ply-p7-gomarkdown-fallback-regression-gate.d4f0538.6YaNVS`:
  final verified 57,472-entry manifest SHA-256
  `daba32728fe7cf75d594371d5e04f0be74101f3a38cd007aaf4e6a51210f8e4f`;
  regression-summary SHA-256
  `4add3dd6e08503abaf553113b72e5677e97cad09b9553cbf6f9a7216ae4cc145`;
  renderer-results SHA-256
  `f614be24b822a30abdb1a70268c0658005c5db0ed0ea54c5ea29dc8c686252e5`.
  Full scorecard SHA-256
  `7380fecaf35addf6a476f60483143f343396a7c460ae27630bba207be637b383`
  exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

All four accepted final manifests were reverified completely after their last
use with zero missing, mismatched, or malformed entries and empty verification
stderr.

Tool identities remain Go 1.26.7
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
golangci-lint 2.12.2
`3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
and govulncheck v1.7.0
`0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected indirect `github.com/google/btree v1.0.0`
against latest v1.1.3 as exactly one bounded P7 group. Preliminary primary
proxy evidence resolves latest to 2024-08-21T16:26:17Z commit
`aeba20f7a1e1315badec4eca4fdc9f754f5f880a`, declaring Go 1.18. The current
selection appears only in `go.sum`; `go mod why -m` reports that the main
module does not need it, and no btree package is loaded. Treat those only as
starting observations and independently reproduce them.

Resolve tag, commit, time, checksums, signature status, source identity, Go
floor, repository delta, vulnerability data, exact selected-module/graph/
package/checksum closure, loaded path, explicit exact-get diff, and tidy
projection. Implement one exact selection only if the result is explained,
minimal, floor-compatible, and every focused and full quality gate passes;
otherwise record rejection without dependency metadata edits.

Warm caches from separate Git archives. Never use `go mod download all` in a
measured tree or tidy as implementation. Stop before another dependency,
behavior/source change, language/toolchain or quality-policy change, packaging,
publication, or P8. Do not push, merge, publish, release, delete evidence,
stash, revert, launch a successor, or remove the worktree.
