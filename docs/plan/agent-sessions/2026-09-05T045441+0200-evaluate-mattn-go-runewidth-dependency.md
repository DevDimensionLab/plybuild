# Agent Session: Evaluate Mattn Go-Runewidth Dependency

Status: NEXT
Session ID: `2026-09-05T045441+0200-evaluate-mattn-go-runewidth-dependency`
Created: `2026-09-05T04:54:41+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5118293dee286d29bbf67a67e3679607775983484919bcfcfc6a3bb9ba5d47cf`
Previous: [2026-09-05T042416+0200-evaluate-magiconair-properties-dependency.md](2026-09-05T042416+0200-evaluate-magiconair-properties-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/mattn/go-runewidth v0.0.14` as one bounded dependency group.
Resolve the canonical latest release and highest floor-compatible candidate
from primary evidence. Implement one exact changed selection only if it
preserves the retained Go 1.18 floor, has an explained minimal closure, and
passes every quality contract; do not manufacture a dependency commit when
the exact selected release is already canonical latest.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted go-colorful v1.4.1, and the
rejected/no-change HCL, mousetrap, and properties evaluations. Go-cmp v0.7.0,
Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify
v1.10.1, latest gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and HCL,
mousetrap, or properties selection changes remain rejected for their recorded
floor, loaded-behavior, release-qualification, self-test, exact-latest, or
closure decisions. Do not revisit them or combine another module group. P8
remains queued.

The current build list selects mattn/go-runewidth v0.0.14. Treat the target
version, latest release, Go declaration, closure, loaded package population,
source history, and test status as unknown until independently resolved. This
session may change only go-runewidth's exact required go.mod/go.sum metadata
and the roadmap/handoff record. Do not change production Go, another
dependency, language/toolchain declarations, quality apparatus,
Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only go-colorful commit
`dc4f27c0319c6f023b5750df2148219a4cb32e60`, exact parent
`f40b0322974eee388d5c064596826898e63f13ca`, clean tree
`a5f3282ed71f9e5557718d63e5ea8d23020c44e2`. Properties evaluation began at
documentation commit `469049fdb547850d4874579f937e92e0418c7e6d`, exact
parent dc4f27c; its documentation handoff commit must have exact parent
469049f. Relative to c314bcb, go.mod/go.sum still change only go-colorful
v1.2.0 -> v1.4.1 and add its exact two checksum lines. Ordinary and ignored
status must be empty.

Current dependency measurements remain 234 selected modules, 3,557 graph
edges, 429 native complete-test packages, a 310-line tidy projection, and
exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
20/30/20. The retained main module declares Go 1.18 and prefers toolchain Go
1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` has verified two-entry manifest
SHA-256
`1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

Properties decision evidence is
`/private/tmp/ply-p7-properties-selection.469049f.iT6qcB`, with fully verified
17,240-entry manifest SHA-256
`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`
and decision-summary SHA-256
`23b8556153b7ca9d606d5156d50d7d7791b4978ec2368ac3171ef6370b4f2c2d`.
It proves canonical latest v1.18.11 and every release above selected v1.8.7
declare Go 1.19, while highest Go-1.18-compatible v1.8.6 exact-get requires
21 selection changes including direct Viper v1.15.0 -> v1.14.0 and fails its
complete race self-test twice. Properties v1.8.7 therefore remains unchanged
as documented unresolved floor debt; no implementation commit or repository
quality run was made.

Accepted go-colorful review, exact-quality, and regression evidence remains
fully verified at
`/private/tmp/ply-p7-go-colorful-quality-review.dc4f27c.gs9tGv`,
`/private/tmp/ply-p7-go-colorful-quality-final.dc4f27c.B7eOLZ`, and
`/private/tmp/ply-p7-go-colorful-regression-gate.dc4f27c.X2Y4Aj`, with
13/236,146/19,011 entries and manifest SHA-256 values
`6444e55370fb37a6c2c06c2c1c280b8fcbd34fc1f94f110a70bc7d067b05673d`,
`924500cb4512fa9451cf218b119ce8969e32e7f02a67558807d067e3f41017a6`,
and `2cf393025cf7768ef90fbfa8c0faf378284c3521cb118220dd74ba6fc4735d0e`.
Its selection manifest retains the recorded 47,633-entry digest, but one
mutable Go telemetry counter no longer verifies; the other 47,632 entries and
stable selection summary SHA-256
`07467cdfd752a8ee9c64b1b6632370c8e0d75b6e3715090931ec8831b2080658`
remain valid. Preserve this correction and the documented btree
regression-manifest correction.

No-change mousetrap evidence remains
`/private/tmp/ply-p7-mousetrap-selection.516ae39.PzdXjl`, verified manifest
29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`.
Rejected HCL evidence remains
`/private/tmp/ply-p7-hcl-selection.99f508d.CFvyqD`, verified manifest
21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.

# Role And Boundaries

From fresh external archives and caches, resolve go-runewidth releases through
the Go proxy, checksum database, and upstream repository. Record exact tag
commits/times, module Go declarations, checksum pairs, source identity, tag
and commit signature status, relevant release history, and primary Go
vulnerability data. Prove which stable release is canonical latest and which
release is highest compatible with Go 1.18. Distinguish stable releases,
prereleases, retractions, and major-path changes explicitly; do not infer
compatibility from a successful modern-toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. If go-runewidth is loaded, run focused consumer behavior in addition
to the candidate module's complete tests, repository build, complete
tests/race/vet, pinned lint, byte-identical public help, and identical API/CLI
reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected release is already canonical latest and exact get is a no-op, record
that no-change decision without forcing an implementation commit. Measure the
exact vulnerability population; do not assume parity.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, properties decision evidence, accepted
go-colorful evidence, mousetrap and HCL evidence, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve both manifest corrections.

# Three Moves

Only if all decision evidence passes and the selection changes, use exact Go
1.26.7 and exact
`go get github.com/mattn/go-runewidth@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy
as implementation. Preserve every retained selection, including properties
v1.8.7, go-colorful v1.4.1, mousetrap v1.1.0, HCL v1.0.0, go-cmp v0.6.0,
btree v1.1.3, gomarkdown fallback, regexp2 v1.12.0, Cobra v1.10.2, YAML
v3.0.5, go-md2man v2.0.7, Blackfriday v2.1.0, pflag v1.0.10, Uniseg v0.4.7,
Colorable v0.1.15, and go-isatty v0.0.20.

After a changed selection, re-run the complete P7 dependency gate: focused
behavior when loaded, graph/path, tests/race/vet, pinned lint, help/API/CLI,
launcher and Make contracts, preflight, host plus fresh snapshot/Docker meta
and acceptance, audit meta, focused and exact Q0-Q2 audits, separate full
audit, vulnerability comparison, empty-HOME count-2, and final
ordinary/ignored cleanliness. Exact make quality must exit 0 with all 27 rows
at L2 and zero held, regressed, not-comparable, or dirty counts. The full
audit may exit 1 only for established queued L3 rows, never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not
implement that next group, launch a successor, push, merge, publish, release,
stash, revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
