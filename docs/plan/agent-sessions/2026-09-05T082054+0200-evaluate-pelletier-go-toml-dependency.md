# Agent Session: Evaluate Pelletier Go-Toml Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-05T082054+0200-evaluate-pelletier-go-toml-dependency`
Created: `2026-09-05T08:20:54+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `4c798d1f43ba0a47c0b607e84d65b1f029d6b65e1d50f2f8888de3212dd4d643`
Previous: [2026-09-05T073150+0200-evaluate-mitchellh-mapstructure-dependency.md](2026-09-05T073150+0200-evaluate-mitchellh-mapstructure-dependency.md)
Next: [2026-09-05T090435+0200-evaluate-pelletier-go-toml-v2-dependency.md](2026-09-05T090435+0200-evaluate-pelletier-go-toml-v2-dependency.md)
Outcome: Retained canonical latest and highest Go-1.18-compatible go-toml
  v1.9.5 unchanged. Exact get was a zero-byte no-op with an empty closure;
  the v1 module is not loaded, its complete tests and every repository
  decision contract passed, and no dependency implementation commit was
  manufactured.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/pelletier/go-toml v1.9.5` as one bounded dependency group. Resolve
the canonical latest release and highest floor-compatible candidate from
primary evidence. Implement one exact changed selection only if it preserves
the retained Go 1.18 floor, has an explained minimal closure, and passes every
quality contract; do not manufacture a dependency commit when the exact
selected release is already canonical latest.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted go-runewidth v0.0.17, and the
rejected/no-change HCL, mousetrap, properties, and mapstructure evaluations.
Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0,
fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap, properties, or
mapstructure selection changes remain rejected or unnecessary for their
recorded floor, loaded-behavior, release-qualification, self-test,
exact-latest, or closure decisions. Do not revisit them or combine another
module group. P8 remains queued.

The current build list selects pelletier/go-toml v1.9.5. Treat the target
version, latest release, Go declaration, closure, loaded package population,
source history, and test status as unknown until independently resolved. This
session may change only go-toml's exact required go.mod/go.sum metadata and the
roadmap/handoff record. Do not change production Go, another dependency,
language/toolchain declarations, quality apparatus, Docker/release inputs,
packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation is dependency-only go-runewidth commit
`ca19dcd4da9320112e6c9f0c5db507ccbe8b88e5`, exact parent
`5d0fa3ae0fff98784c349cc23e001018eb9637e1`, clean tree
`273c9b635e988ddd686d6c83846f7320bb9e4e33`. Mapstructure required no
implementation commit; the current documentation handoff must have exact
parent `f974b581baf36f40bb8fa31868371eb0db712faf`, whose exact parent is
ca19dcd. Relative to accepted go-cmp commit c314bcb, go.mod/go.sum change only
go-colorful v1.2.0 -> v1.4.1 and go-runewidth v0.0.14 -> v0.0.17, adding
exactly those two selected versions' checksum pairs. Ordinary and ignored
status must be empty.

Current dependency measurements are 234 selected modules, 3,557 graph edges,
429 native complete-test packages, a 313-line tidy projection, and exact
Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations 20/30/20.
The retained main module declares Go 1.18 and prefers toolchain Go 1.26.7.

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

No-change mapstructure evidence is fully verified at
`/private/tmp/ply-p7-mapstructure-selection.f974b58.GbAMTJ`, with 46,352
entries and manifest SHA-256
`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
Its selection-summary SHA-256 is
`8910174c8f8fa387a32cd49871dfbd37668ff4a2fbf1c51c2dae6fdc2612eff3`.
It proves the original path's canonical latest and highest Go-1.18-compatible
release are both selected v1.5.0, with an empty exact-get closure, passing
loaded behavior and repository decision gates, and no implementation commit.
Preserve mapstructure v1.5.0.

Accepted go-runewidth selection, review, exact-quality, and regression
evidence remains fully verified at
`/private/tmp/ply-p7-go-runewidth-selection.5d0fa3.dmbKoB`,
`/private/tmp/ply-p7-go-runewidth-quality-review.ca19dcd.UgdZ2N`,
`/private/tmp/ply-p7-go-runewidth-quality-final.ca19dcd.t6LExP`, and
`/private/tmp/ply-p7-go-runewidth-regression-gate.ca19dcd.O6Uxyq`, with
28,870/15/236,329/4,888 entries and manifest SHA-256 values
`2d6986055d24d4937dcbbaa8fa1e46da6a7bcc1e94f481125343284c96687193`,
`f885b995dc6192cb7ef2d8c9241f93cb44a885deda7faec7e0415f68ce21d764`,
`28fbbde41059dd806dfea5653f273a390b6b1fa2f60ebeec2e3d654cc792e2ed`,
and `b83091c1c8a52ec2176c4c3fd7a4c3c5e456adf7a61ee2c6b6f7f97706eb4894`.
Preserve go-runewidth v0.0.17 and the recorded external Python/buildx
accommodations when a complete changed-selection gate requires them.

Properties decision evidence remains fully verified at
`/private/tmp/ply-p7-properties-selection.469049f.iT6qcB`,
17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`.
No-change mousetrap and rejected HCL evidence remain fully verified at
29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`
and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
Preserve the recorded go-colorful mutable telemetry correction and btree
regression-manifest correction.

# Role And Boundaries

From fresh external archives and caches, resolve go-toml releases through the
Go proxy, checksum database, and upstream repository. Record exact tag
commits/times, module Go declarations, checksum pairs, source identity, tag
and commit signature status, relevant release history, and primary Go
vulnerability data. Prove which stable release is canonical latest and which
release is highest compatible with Go 1.18. Distinguish stable releases,
prereleases, retractions, module-path or repository lineage changes, and
archived/deprecated status explicitly. In particular, do not treat a distinct
`github.com/pelletier/go-toml/v2` module as an in-place v1 selection. Do not
infer compatibility from a successful modern-toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. If go-toml is loaded, run focused consumer behavior in addition to
the candidate module's complete tests, repository build, complete
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
the roadmap, go.mod/go.sum, mapstructure decision evidence, accepted
runewidth evidence, properties decision evidence, accepted go-colorful
evidence, mousetrap and HCL evidence, and the toolchain, compatibility,
snapshot/Docker, quality, baseline-reproduction, and audit contracts. Preserve
both manifest corrections.

# Three Moves

Only if all decision evidence passes and the selection changes, use exact Go
1.26.7 and exact
`go get github.com/pelletier/go-toml@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy
as implementation. Preserve every retained selection, including mapstructure
v1.5.0, runewidth v0.0.17, properties v1.8.7, go-colorful v1.4.1, mousetrap
v1.1.0, HCL v1.0.0, go-cmp v0.6.0, btree v1.1.3, gomarkdown fallback,
regexp2 v1.12.0, Cobra v1.10.2, YAML v3.0.5, go-md2man v2.0.7, Blackfriday
v2.1.0, pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15, and go-isatty
v0.0.20.

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
