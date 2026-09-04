# Agent Session: Evaluate Hashicorp HCL Dependency

Status: NEXT
Session ID: `2026-09-05T011325+0200-evaluate-hashicorp-hcl-dependency`
Created: `2026-09-05T01:13:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `85aeb4b4c992fe99dbd78df28e99f1b3fa10390d5e592251074ea261e375a8bd`
Previous: [2026-09-04T231513+0200-evaluate-google-go-cmp-dependency.md](2026-09-04T231513+0200-evaluate-google-go-cmp-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/hashicorp/hcl v1.0.0` as one bounded dependency group. Resolve the
latest release and highest floor-compatible candidate from primary evidence.
Implement one exact selection only if it preserves the retained Go 1.18 floor,
has an explained minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through accepted Google go-cmp v0.6.0. Go-cmp
v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0,
fsnotify v1.10.1, and latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` remain rejected for their recorded floor
or loaded-behavior failures. Do not revisit those decisions or combine another
module group. P8 remains queued.

The current build list selects Hashicorp HCL v1.0.0. Treat the target version,
latest release, Go declaration, closure, loaded package population, and source
history as unknown until independently resolved. This session may change only
HCL's exact required go.mod/go.sum metadata and the roadmap/handoff record. Do
not change production Go, another dependency, the language/toolchain
declaration, quality apparatus, Docker/release inputs, packaging, publishers,
or P8 code.

# Measurements At Start

Accepted go-cmp implementation is commit
`c314bcb440b5f94871d71a249bca7ae7f87d5543`, exact parent
`cbdb0a915d7d30bb6633cfbb2ccf0323aea1b77a`, clean tree
`9a212b377cb08e7d6e3044fb6dae071058b6242b`. After the documentation handoff,
continuity HEAD must have exact parent `c314bcb`, and ordinary and ignored
status must be empty.

Current dependency measurements are 234 selected modules, 3,557 graph edges,
429 complete test packages, a 308-line tidy projection, and exact
Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations 20/30/20.
The retained main module declares Go 1.18 and prefers toolchain Go 1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. The SDK was restored from the
official archive after external tmp cleanup; recovery evidence is sealed at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs`, verified two-entry manifest
SHA-256 `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

Accepted go-cmp selection evidence is
`/private/tmp/ply-p7-go-cmp-selection.cbdb0a9.kbzc4C`, with verified
29,812-entry manifest SHA-256
`dcd18baaa45eb8ad59b268fa7c356164dabc5d4c9836423c3fd1c3aa6d0d1797`
and selection-summary SHA-256
`5f76248b77f1451bfb3304b0857ba9464c3a5f75921ba6df26b750080c361ee1`.

Commit-bound schema-2 evidence is
`/private/tmp/ply-p7-go-cmp-quality-review.c314bcb.XeWg4T`, with verified
9,487-entry manifest SHA-256
`24ede5a473b9ac9f47839b121d25e9a8d921369c5ba149bf42abd902b5271727`
and manual-evidence SHA-256
`6f7a469bf4fa8da00a44570d4b7c669f64c5a31dbc116b4304c0604c505ba1ac`.

Exact accepted quality evidence is
`/private/tmp/ply-p7-go-cmp-quality-parent.c314bcb.zxFKyw`. Its verified
252,077-entry manifest SHA-256 is
`6d8104b43fdcb838ec8c0955870798be38a3b4a1f9a7e73f22e683ee566a9f69`;
Q0-Q2 scorecard SHA-256 is
`6531ade31d4c504c304cd707cbab158e9433d487ec8fa96d05a8277d9dc149e1`.
Exact make quality passes all 21 stages and 27 rows at L2, 80/80 mutations,
8/8 mutation and 4/4 acceptance populations, with held, regressed,
not-comparable, and dirty counts zero.

The independent regression root is
`/private/tmp/ply-p7-go-cmp-regression-gate.c314bcb-final.rMW0l5`. Its verified
30,255-entry manifest SHA-256 is
`00b736b53494a0c4f27ee7503ecef1c7c7bf74f00bdaab85f7afa9eab19ed3c5`;
regression-summary SHA-256 is
`a27e5f4d00fa0f60d9246309e1485353cdb2aacecd4a3d80e0d3bf64efe2323e`.
It binds exact graph/package/tidy populations, module and repository gates,
identical help/API/CLI, exact 20/30/20 vulnerabilities, empty-HOME count-2,
and the expected full-audit L3-only exit.

# Role And Boundaries

From fresh external archives and caches, resolve HCL releases through the Go
proxy, checksum database, and upstream repository. Record exact tag
commits/times, module Go declarations, checksum pairs, source identity, tag and
commit signature status, relevant history since v1.0.0, and primary Go
vulnerability data. Prove which release is latest and which is the highest
compatible with Go 1.18; do not infer compatibility from a successful modern
toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. If HCL is loaded, run focused consumer behavior in addition to the
candidate module's complete tests, repository build, complete tests/race/vet,
pinned lint, byte-identical public help, and identical API/CLI reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. Measure
the exact vulnerability population; do not assume parity.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, accepted go-cmp evidence, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve the documented btree regression-manifest correction.

# Three Moves

Only if all decision evidence passes, use exact Go 1.26.7 and exact
`go get github.com/hashicorp/hcl@<selected-version>` for one dependency-only
commit. Do not hand-edit module metadata and do not use tidy as implementation.
Preserve every retained selection, including go-cmp v0.6.0, btree v1.1.3,
gomarkdown fallback, regexp2 v1.12.0, Cobra v1.10.2, YAML v3.0.5, go-md2man
v2.0.7, Blackfriday v2.1.0, pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15,
and go-isatty v0.0.20.

Re-run the complete P7 dependency gate: focused behavior when loaded,
graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and Make
contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
make quality must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. The full audit may exit 1 only for established
queued L3 rows, never 2.

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
