# Quality Upgrade Handover

Generated: 2026-09-04T23:15:13+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is dependency-only Google btree commit
  `2e2f8e09918d1105cdc24ac34214690daf8d0cff`, exact parent
  `6f3f39db3d2f449937f0e0bb35c3a3e919d627ed`, clean tree
  `b70ada27a9db80515551c07c2c738f3fe0336498`.
- The documentation handoff must have exact parent `2e2f8e0`; ordinary and
  ignored status must be empty afterward. The answered btree archive and the
  Google go-cmp NEXT archive link reciprocally.
- Only `go.mod` and `go.sum` changed in the implementation. Production Go,
  API/CLI, language/toolchain declarations, quality apparatus, Docker/release
  inputs, packaging, publishers, and P8 code are unchanged.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, successor launch, retained
  evidence/image deletion, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through accepted
`github.com/google/btree v1.1.3`. Viper v1.16.0, Emoji v2.2.14, ini v1.67.3,
fatih/color v1.19.0, fsnotify v1.10.1, and latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` remain rejected for their recorded floor
or loaded-behavior failures. The accepted gomarkdown fallback remains
`v0.0.0-20260824151336-45814d58469f`. P8 remains queued.

Only the launcher's mutable header and prompt regions may change during a
handoff. Keep caches, projections, reports, artifacts, build contexts, and
audit evidence outside the worktree.

## Accepted Google Btree Group

Selected indirect `github.com/google/btree` moved from v1.0.0 to latest
v1.1.3. Independent proxy, checksum-database, and repository evidence binds
v1.1.3 to lightweight unsigned tag commit
`aeba20f7a1e1315badec4eca4fdc9f754f5f880a`, published and committed at
2024-08-21T16:26:17Z, with module declaration `go 1.18`. Its checksum pair is
`h1:CVpQJjYgC4VbzxeGVHfvZrv1ctoYCAI8vbl07Fcxlyg=` /
`h1:qOPhT0dTNdNzV6Z/lhRX0YXUafgPLFUh+gZMl761Gm4=`. All nine proxy source
files are byte-identical to the tag-commit archive. The Go vulnerability
module index has no btree record.

The 15-commit history since v1.0.0 adds module metadata and the generic
`BTreeG` API, requires Go 1.18, repairs `Has`, and in v1.1.3 fixes legacy
`BTree` clone copy-on-write behavior after child removal. The tag is a commit
object rather than an annotated tag; tag verification therefore fails as
non-tag, and both the commit object and GitHub verification record are
unsigned.

Old and candidate states both select 234 modules and the same 429 complete
test packages. The only selected-module change is btree v1.0.0 -> v1.1.3.
Graph edges change 3,556 -> 3,557 solely because exact selection adds the
main-module edge to btree v1.1.3. The historical dependency edges selecting
v1.0.0 remain, so no transitive selection changes. No btree package is loaded
in the complete main-module test population, and `go mod why -m` says the main
module does not need it.

Exact `go get github.com/google/btree@v1.1.3` adds one indirect requirement
and the v1.1.3 checksum pair; no prior checksum is removed. Tidy remains an
unapplied projection: 285 -> 299 lines. Its only candidate-specific changes
are removal of the explicit unloaded btree requirement and the two newly
selected checksum lines; all other lines are pre-existing tidy debt.

The candidate module's complete tests, race, vet, and focused clone tests at
count 10 pass. Repository build, complete tests/race/vet, pinned lint,
byte-identical public help, identical API/CLI reports, CLI surface, exact
quality, empty-HOME count-2, and clean-tree gates pass. Vulnerability ID sets
remain exactly 20 Darwin symbol, 30 Darwin module, and 20 Windows symbol.

## Accepted Evidence

- Selection root `/private/tmp/ply-p7-btree-selection.6f3f39d.oPbcTr`:
  verified 52,933-entry manifest SHA-256
  `bc86252747c48e4389fd1b60c8c4843f8003c7342a4f7cf79b66b7f1bf4d8006`;
  selection-summary SHA-256
  `45d485b4cf7d1d6d2db91f2b320b4b3fe2ef42fff36a6803619d576c8dd3b1e8`.
- Commit-bound schema-2 review
  `/private/tmp/ply-p7-btree-quality-review.2e2f8e0.kUNRFX`: verified
  9,482-entry manifest SHA-256
  `52aa4d4743e2cc70fa9f1d17a2e068941d7dbb67f90ee25e8b2dc6c9652c5ee1`;
  manual-evidence SHA-256
  `e7625cadbab2cd817c9e932e28b83248a7b000573adc4f24fb00f40a8380d03d`.
  All six focused receipts pass, and all 112 governed evidence files / 258
  subjects are unchanged from the parent.
- Exact `make quality` root
  `/private/tmp/ply-p7-btree-quality-parent.2e2f8e0-final.6v4pBz`: verified
  252,397-entry manifest SHA-256
  `2395eb672221bac15a70b564dbce09bf28546d93876111dd22fff6500b1fc7cf`;
  Q0-Q2 scorecard SHA-256
  `980dec57493b17f6d1954dc247e59f1e3be81754ccf230efa60f059a8d273395`.
  The exact 21-stage ledger exits 0 with 27/27 rows at L2, 80/80 mutations,
  8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
  not-comparable, or dirty counts.
- Independent regression root
  `/private/tmp/ply-p7-btree-regression-gate.2e2f8e0.FhWpnk`: verified
  24,197-entry manifest SHA-256
  `9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`;
  regression-summary SHA-256
  `f348bd0b9e0f68c698c43f466904b66c1c3db6eab5e783b14364aaf40b6b1ea2`.
  It independently binds graph/package/tidy populations, module and
  repository gates, byte-identical help/API/CLI, exact 20/30/20
  vulnerabilities, empty-HOME count-2, and the expected full-audit exit 1
  only for Q3.1, Q3.3, Q3.4, and Q3.7.

The prior accepted gomarkdown final evidence remains sealed at decision,
review, exact-quality, and regression roots recorded in the answered
gomarkdown archive. All accepted manifests used here were completely
reverified after their final use.

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

Independently evaluate selected indirect `github.com/google/go-cmp v0.5.9`
as exactly one bounded P7 module group. Preliminary proxy evidence identifies
latest v0.7.0 at commit `9b12f366a942ebc7254abc7f32ca05068b455fb7`,
published 2025-01-14T18:15:44Z, with a Go 1.21 declaration that appears to
exceed the retained floor. If independently confirmed, reject latest and
evaluate highest floor-compatible v0.6.0 in the same module group; preliminary
proxy evidence binds v0.6.0 to commit
`c3ad8435e7bef96af35732bc0789e5a2278c6d5f`, published
2023-08-31T17:32:40Z, with Go 1.13. Treat all of these as candidate inputs.

The complete current test population loads five go-cmp packages through
`plybuild/cmd -> mvn-pom-mutator/pkg/pom -> go-cmp/cmp`, so compatibility must
be proven through focused consumer behavior as well as module self-tests.
Resolve both canonical identities, checksums, source/signature status,
history, vulnerability data, exact selection/edge/package/checksum closure,
loaded path, explicit exact-get diff, and tidy projection. Implement exactly
v0.6.0 only if latest is rejected for the floor, the fallback is the highest
compatible release, closure is minimal and explained, and every focused and
full quality gate passes.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source change, language/toolchain or quality-policy change,
packaging, publication, or P8. Do not push, merge, publish, release, delete
evidence, stash, revert, launch a successor, or remove the worktree.
