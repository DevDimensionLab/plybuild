# Agent Session: Evaluate Google Btree Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-04T202725+0200-evaluate-google-btree-dependency`
Created: `2026-09-04T20:27:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `cf5ece372228bf7a701796fa565abd122938dbb450718e11d06f4202e0fc043a`
Previous: [2026-09-04T171914+0200-evaluate-gomarkdown-compatible-fallback.md](2026-09-04T171914+0200-evaluate-gomarkdown-compatible-fallback.md)
Next: [2026-09-04T231513+0200-evaluate-google-go-cmp-dependency.md](2026-09-04T231513+0200-evaluate-google-go-cmp-dependency.md)
Outcome: Upgraded only selected indirect `github.com/google/btree` from v1.0.0 to latest v1.1.3 in dependency-only commit `2e2f8e0`; primary identity, exact closure, module self-tests, repository behavior, exact quality, vulnerability, and clean-tree gates all passed, and P7 remains active for the measured Google go-cmp group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/google/btree v1.0.0` against latest `v1.1.3` as one bounded
dependency group. Resolve the candidate from primary evidence and implement it
with one exact Go selection only if it preserves the retained Go 1.18 floor,
has an explained minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through accepted gomarkdown fallback
`v0.0.0-20260824151336-45814d58469f`. Viper v1.16.0, Emoji v2.2.14, ini
v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, and latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` remain rejected for their recorded floor
or loaded-behavior failures. Do not revisit those decisions or combine another
group. P8 remains queued.

The current build list selects google/btree v1.0.0 through the historical
module graph, while the complete main-module test population loads no
google/btree package and `go mod why -m` reports that the main module does not
need it. Preliminary primary proxy evidence identifies latest v1.1.3 at tag
commit `aeba20f7a1e1315badec4eca4fdc9f754f5f880a`, published
2024-08-21T16:26:17Z, with a Go 1.18 declaration. Treat every one of those
facts as a candidate input until independently reproduced.

This session may change only google/btree's exact required go.mod/go.sum
metadata and the roadmap/handoff record. Do not change production Go, another
dependency, the language/toolchain declaration, quality apparatus,
Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Accepted gomarkdown implementation is commit
`d4f053808e389cf163665087f78690cca05999ce`, exact parent
`18d67bb9293cef56b13b8bcb9e3dfbe22fd23045`, clean tree
`cbf28809e55d95589b86a0d9e57eeaba2401542d`. After the documentation handoff,
continuity HEAD must have exact parent `d4f0538`, and ordinary and ignored
status must be empty.

Current dependency measurements are 234 selected modules, 3,556 graph edges,
429 complete test packages, four loaded gomarkdown packages, a 285-line tidy
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

The accepted gomarkdown decision replay is
`/private/tmp/ply-p7-gomarkdown-fallback-replay.18d67bb.U5kpjt`. Its original
42,120-entry seal is retained at SHA-256
`dfdce8f20a3dbd4257a0a1dc22034d2c903f000a20f935d12c0509cfff303b4c`;
four later-used Go telemetry counters drifted, and its verified final
42,127-entry post-use manifest SHA-256 is
`b438197999acb18ad6955d66c5d38e0af92df8b0693ecba7445f5a960b1b1fa8`.
Selection-summary SHA-256 is
`697d2467a20850656b1213494c9c2fe4557df8a61d1ec9eacff94cdf41d2b28a`.

Commit-bound schema-2 evidence is
`/private/tmp/ply-p7-gomarkdown-fallback-quality-review.d4f0538.sBeoKD`, with
verified 2,304-entry manifest SHA-256
`cdd7691693373d726fdec2e954aacdc11bcae0bbffde3a471edd56c34de78ab6`
and manual-evidence SHA-256
`ffdee0e8a622a8de34220a7773bf32b87d5368698cc25d18a28c2b49b711cc17`.

Exact accepted quality evidence is
`/private/tmp/ply-p7-gomarkdown-fallback-quality-parent.d4f0538-final.2kKFpS`.
Its verified 252,314-entry manifest SHA-256 is
`70ef46136c13b7e5276ac1ef1af62514013902e93c9e115551d863f0f3a8e40b`;
Q0-Q2 scorecard SHA-256 is
`a2298bb0457549c71c90a9069d3f1ffe219f7f21b62d7b47684d16dc7fdfe870`.
Exact make quality passes all 27 rows at L2, 80/80 mutations, 8/8 mutation and
4/4 acceptance populations, with held, regressed, not-comparable, and dirty
counts zero.

The independent regression root is
`/private/tmp/ply-p7-gomarkdown-fallback-regression-gate.d4f0538.6YaNVS`.
Its verified 57,472-entry manifest SHA-256 is
`daba32728fe7cf75d594371d5e04f0be74101f3a38cd007aaf4e6a51210f8e4f`;
regression-summary SHA-256 is
`4add3dd6e08503abaf553113b72e5677e97cad09b9553cbf6f9a7216ae4cc145`.
It binds the exact graph/package/tidy populations, 30 retained renderer
passes across 34 individual cases with zero regressions, exact 20/30/20
vulnerabilities, empty-HOME count-2, and the expected full-audit L3-only exit.

# Role And Boundaries

From fresh external archives and caches, resolve google/btree v1.1.3 through
the Go proxy, checksum database, and upstream repository. Record the exact tag
commit/time, module Go declaration, checksum pair, source identity, tag and
commit signature status, relevant history since v1.0.0, and primary Go
vulnerability data. Do not infer compatibility from the preliminary proxy
check.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require exactly one changed module selection and explain
every main-edge or checksum change. Because google/btree is currently not
loaded, prove that status independently and run the candidate module's own
complete tests in addition to repository build, complete tests/race/vet,
pinned lint, byte-identical public help, and identical API/CLI reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity,
module self-tests, or any repository behavior/quality contract fails. Measure
the exact vulnerability population; do not assume parity merely because the
module is unloaded.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the accepted gomarkdown evidence, and the
toolchain, compatibility, snapshot/Docker, quality, baseline-reproduction,
and audit contracts.

# Three Moves

Only if all decision evidence passes, use exact Go 1.26.7 and exact
`go get github.com/google/btree@v1.1.3` for one dependency-only commit. Do not
hand-edit module metadata and do not use tidy as the implementation command.
Preserve every retained selection, including gomarkdown fallback, regexp2
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
