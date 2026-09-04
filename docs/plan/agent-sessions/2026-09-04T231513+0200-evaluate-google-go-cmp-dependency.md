# Agent Session: Evaluate Google Go-Cmp Dependency

Status: NEXT
Session ID: `2026-09-04T231513+0200-evaluate-google-go-cmp-dependency`
Created: `2026-09-04T23:15:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `503b852b2cbbaeb482ee1194002e441a5eb6e3472ed89cccf1e7baf45dd3789f`
Previous: [2026-09-04T202725+0200-evaluate-google-btree-dependency.md](2026-09-04T202725+0200-evaluate-google-btree-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/google/go-cmp v0.5.9` as one bounded dependency group. Resolve and
reject latest v0.7.0 if its preliminary Go 1.21 floor is confirmed, then
evaluate highest floor-compatible v0.6.0 in the same module group. Implement
one exact selection only if it preserves the retained Go 1.18 floor, has an
explained minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through accepted Google btree v1.1.3. Viper
v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1,
and latest gomarkdown `v0.0.0-20260824154242-13c5cf49db8d` remain rejected for
their recorded floor or loaded-behavior failures. Do not revisit those
decisions or combine another module group. P8 remains queued.

The current build list selects go-cmp v0.5.9. The complete test population
loads five go-cmp packages through
`plybuild/cmd -> mvn-pom-mutator/pkg/pom -> go-cmp/cmp`. Preliminary primary
proxy evidence identifies latest v0.7.0 at tag commit
`9b12f366a942ebc7254abc7f32ca05068b455fb7`, published
2025-01-14T18:15:44Z, with Go 1.21. Preliminary proxy evidence identifies
v0.6.0 at tag commit `c3ad8435e7bef96af35732bc0789e5a2278c6d5f`,
published 2023-08-31T17:32:40Z, with Go 1.13. Treat every one of those facts as
a candidate input until independently reproduced.

This session may change only go-cmp's exact required go.mod/go.sum metadata
and the roadmap/handoff record. Do not change production Go, another
dependency, the language/toolchain declaration, quality apparatus,
Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Accepted btree implementation is commit
`2e2f8e09918d1105cdc24ac34214690daf8d0cff`, exact parent
`6f3f39db3d2f449937f0e0bb35c3a3e919d627ed`, clean tree
`b70ada27a9db80515551c07c2c738f3fe0336498`. After the documentation handoff,
continuity HEAD must have exact parent `2e2f8e0`, and ordinary and ignored
status must be empty.

Current dependency measurements are 234 selected modules, 3,557 graph edges,
429 complete test packages, five loaded go-cmp packages, a 299-line tidy
projection, and exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability
populations 20/30/20. The retained main module declares Go 1.18 and prefers
toolchain Go 1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Retain the verified
golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and govulncheck v1.7.0
binaries and hashes recorded in the handover.

Accepted btree selection evidence is
`/private/tmp/ply-p7-btree-selection.6f3f39d.oPbcTr`, with verified
52,933-entry manifest SHA-256
`bc86252747c48e4389fd1b60c8c4843f8003c7342a4f7cf79b66b7f1bf4d8006`
and selection-summary SHA-256
`45d485b4cf7d1d6d2db91f2b320b4b3fe2ef42fff36a6803619d576c8dd3b1e8`.

Commit-bound schema-2 evidence is
`/private/tmp/ply-p7-btree-quality-review.2e2f8e0.kUNRFX`, with verified
9,482-entry manifest SHA-256
`52aa4d4743e2cc70fa9f1d17a2e068941d7dbb67f90ee25e8b2dc6c9652c5ee1`
and manual-evidence SHA-256
`e7625cadbab2cd817c9e932e28b83248a7b000573adc4f24fb00f40a8380d03d`.

Exact accepted quality evidence is
`/private/tmp/ply-p7-btree-quality-parent.2e2f8e0-final.6v4pBz`. Its verified
252,397-entry manifest SHA-256 is
`2395eb672221bac15a70b564dbce09bf28546d93876111dd22fff6500b1fc7cf`;
Q0-Q2 scorecard SHA-256 is
`980dec57493b17f6d1954dc247e59f1e3be81754ccf230efa60f059a8d273395`.
Exact make quality passes all 21 stages and 27 rows at L2, 80/80 mutations,
8/8 mutation and 4/4 acceptance populations, with held, regressed,
not-comparable, and dirty counts zero.

The independent regression root is
`/private/tmp/ply-p7-btree-regression-gate.2e2f8e0.FhWpnk`. Its verified
24,197-entry manifest SHA-256 is
`9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`;
regression-summary SHA-256 is
`f348bd0b9e0f68c698c43f466904b66c1c3db6eab5e783b14364aaf40b6b1ea2`.
It binds exact graph/package/tidy populations, module and repository gates,
identical help/API/CLI, exact 20/30/20 vulnerabilities, empty-HOME count-2,
and the expected full-audit L3-only exit.

# Role And Boundaries

From fresh external archives and caches, resolve go-cmp v0.7.0 and v0.6.0
through the Go proxy, checksum database, and upstream repository. Record exact
tag commits/times, module Go declarations, checksum pairs, source identity,
tag and commit signature status, relevant history since v0.5.9, and primary
Go vulnerability data. Do not infer compatibility from the preliminary proxy
checks.

First resolve latest v0.7.0. If its Go 1.21 declaration is confirmed, record
its rejection without selecting it. Then prove that v0.6.0 is the highest
release compatible with Go 1.18 before evaluating its closure and behavior.
Keep latest rejection and fallback acceptance in this one go-cmp group; do not
evaluate another module.

Measure old versus fallback selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require exactly one changed module selection and explain
every main-edge or checksum change. Because go-cmp is loaded, run focused
mvn-pom-mutator/consumer behavior and the candidate module's complete tests in
addition to repository build, complete tests/race/vet, pinned lint,
byte-identical public help, and identical API/CLI reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. Measure
the exact vulnerability population; do not assume parity.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the accepted btree evidence, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts.

# Three Moves

Only if all decision evidence passes, use exact Go 1.26.7 and exact
`go get github.com/google/go-cmp@v0.6.0` for one dependency-only commit. Do not
hand-edit module metadata and do not use tidy as implementation. Preserve every
retained selection, including btree v1.1.3, gomarkdown fallback, regexp2
v1.12.0, Cobra v1.10.2, YAML v3.0.5, go-md2man v2.0.7, Blackfriday v2.1.0,
pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15, and go-isatty v0.0.20.

Re-run the complete P7 dependency gate: focused behavior, graph/path,
tests/race/vet, pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact make quality
must exit 0 with all 27 rows at L2 and zero held, regressed, not-comparable, or
dirty counts. The full audit may exit 1 only for established queued L3 rows,
never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
that next group, launch a successor, push, merge, publish, release, stash,
revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
